// Copyright (c) The Thanos Authors.
// Licensed under the Apache License 2.0.

package store

import (
	"context"
	"path/filepath"
	"sync"
	"testing"

	"github.com/go-kit/log"
	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/tsdb"

	"github.com/efficientgo/core/testutil"
	"github.com/thanos-io/objstore/providers/filesystem"
	"github.com/thanos-io/thanos/pkg/block/indexheader"
	"github.com/thanos-io/thanos/pkg/block/metadata"
	storetestutil "github.com/thanos-io/thanos/pkg/store/storepb/testutil"
)

// TestRace_6575 triggers the data race on sort.Slice mutating shared
// matchers slice in ExpandedPostings().
//
// BUG: bucketIndexReader.ExpandedPostings() calls sort.Slice(ms, ...) on
// the shared matchers slice. When multiple goroutines call ExpandedPostings()
// concurrently with the same matchers slice (e.g. from Series() which passes
// blockMatchers to multiple blocks), sort.Slice mutates the slice elements
// concurrently, causing a data race.
//
// FIX: sort.Slice is moved OUT of ExpandedPostings() to the outermost call
// sites (Series/LabelNames/LabelValues). A new sortedMatchers type wraps
// pre-sorted matchers, and ExpandedPostings() accepts sortedMatchers instead
// of []*labels.Matcher, removing the in-place sort.
//
// Issue #6545. Regression test TestExpandedPostingsRace added in the fix PR.
func TestRace_6575(t *testing.T) {
	tmpDir := t.TempDir()

	bkt, err := filesystem.NewBucket(filepath.Join(tmpDir, "bkt"))
	testutil.Ok(t, err)
	defer func() { testutil.Ok(t, bkt.Close()) }()

	id := uploadTestBlock(t, tmpDir, bkt, 100)

	r, err := indexheader.NewBinaryReader(context.Background(), log.NewNopLogger(), bkt, tmpDir, id, DefaultPostingOffsetInMemorySampling)
	testutil.Ok(t, err)

	b := &bucketBlock{
		logger:            log.NewNopLogger(),
		metrics:           newBucketStoreMetrics(nil),
		indexHeaderReader: r,
		indexCache:        noopCache{},
		bkt:               bkt,
		meta:              &metadata.Meta{BlockMeta: tsdb.BlockMeta{ULID: id}},
		partitioner:       NewGapBasedPartitioner(PartitionerMaxGapSize),
	}

	indexr := newBucketIndexReader(b)

	// Create a shared matchers slice — all goroutines will pass the same slice.
	// sort.Slice in ExpandedPostings mutates this slice concurrently.
	sharedMatchers := []*labels.Matcher{
		labels.MustNewMatcher(labels.MatchEqual, "n", "1"+storetestutil.LabelLongSuffix),
		labels.MustNewMatcher(labels.MatchEqual, "j", "foo"),
	}

	numGoroutines := 50
	iterations := 200

	var wg sync.WaitGroup
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				// Pass the SHARED matchers slice. Each call to ExpandedPostings
				// will call sort.Slice on it, causing concurrent mutations.
				_, _ = indexr.ExpandedPostings(context.Background(), sharedMatchers, NewBytesLimiterFactory(0)(nil))
			}
		}(i)
	}
	wg.Wait()
}
