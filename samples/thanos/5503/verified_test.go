// Copyright (c) The Thanos Authors.
// Licensed under the Apache License 2.0.

package compact

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/tsdb"

	"github.com/thanos-io/thanos/pkg/block"
	"github.com/thanos-io/thanos/pkg/block/metadata"
	"github.com/thanos-io/thanos/pkg/objstore"
	"github.com/thanos-io/thanos/pkg/testutil"
	"github.com/thanos-io/thanos/pkg/testutil/e2eutil"
)

// TestRace_5503 triggers the data race on the named return err variable
// in compact() that is shared across errgroup goroutines.
//
// BUG: compact() has a named return (shouldRerun bool, compID ulid.ULID, err error).
// Inside compact(), multiple g.Go() goroutine closures capture and write to the
// shared err variable concurrently:
//   - err = block.Download(...) at line ~1033
//   - stats, err = block.GatherIndexHealthStats(...) at line ~1043
//
// When multiple errgroup goroutines run concurrently, the Go race detector
// catches concurrent writes to the same err variable.
//
// FIX: The patch changes DoInSpanWithErr to return error, and modifies all
// closure sites to use locally-scoped error variables (:= instead of =).
// Each goroutine now captures its own independent err via named return
// in the inner closure, eliminating the shared mutable state.
//
// WARNING DATA RACE evidence in PR body (Level 1 artifact).
func TestRace_5503(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()

	dir, err := ioutil.TempDir("", "test-race-5503")
	testutil.Ok(t, err)
	defer func() { testutil.Ok(t, os.RemoveAll(dir)) }()

	bkt := objstore.NewInMemBucket()

	// Create many small blocks that can be compacted together.
	numBlocks := 50

	prepareDir, err := ioutil.TempDir("", "test-race-5503-prepare")
	testutil.Ok(t, err)
	defer func() { testutil.Ok(t, os.RemoveAll(prepareDir)) }()

	baseTime := int64(1000000)
	blockSpan := int64(1000)
	for i := 0; i < numBlocks; i++ {
		// Each block gets a unique series name with contiguous, non-overlapping time ranges.
		// Contiguous ranges at the same level enable levelled compaction.
		id, err := e2eutil.CreateBlock(
			ctx, prepareDir,
			[]labels.Labels{{{Name: "a", Value: fmt.Sprintf("%d", i)}}},
			100,
			baseTime + blockSpan*int64(i),     // mint
			baseTime + blockSpan*int64(i+1),   // maxt (contiguous)
			labels.Labels{{Name: "e1", Value: "1"}},
			124,
			metadata.NoneFunc,
		)
		testutil.Ok(t, err)

		testutil.Ok(t, block.Upload(ctx, log.NewNopLogger(), bkt, filepath.Join(prepareDir, id.String()), metadata.NoneFunc))
	}

	// Set up compaction infrastructure with high errgroup concurrency.
	logger := log.NewNopLogger()
	reg := prometheus.NewRegistry()

	insBkt := objstore.WithNoopInstr(bkt)
	ignoreDeletionMarkFilter := block.NewIgnoreDeletionMarkFilter(logger, insBkt, 48*time.Hour, 32)
	duplicateBlocksFilter := block.NewDeduplicateFilter(32)
	noCompactMarkerFilter := NewGatherNoCompactionMarkFilter(logger, insBkt, 2)
	metaFetcher, err := block.NewMetaFetcher(nil, 32, insBkt, "", nil, []block.MetadataFilter{
		ignoreDeletionMarkFilter,
		duplicateBlocksFilter,
		noCompactMarkerFilter,
	})
	testutil.Ok(t, err)

	blocksMarkedForDeletion := promauto.With(nil).NewCounter(prometheus.CounterOpts{})
	blocksMaredForNoCompact := promauto.With(nil).NewCounter(prometheus.CounterOpts{})
	garbageCollectedBlocks := promauto.With(nil).NewCounter(prometheus.CounterOpts{})

	sy, err := NewMetaSyncer(nil, nil, bkt, metaFetcher, duplicateBlocksFilter, ignoreDeletionMarkFilter, blocksMarkedForDeletion, garbageCollectedBlocks)
	testutil.Ok(t, err)

	comp, err := tsdb.NewLeveledCompactor(ctx, reg, logger, []int64{1000, 3000}, nil, nil)
	testutil.Ok(t, err)

	planner := NewPlanner(logger, []int64{1000, 3000}, noCompactMarkerFilter)

	grouper := NewDefaultGrouper(logger, bkt, false, false, reg,
		blocksMarkedForDeletion, garbageCollectedBlocks, blocksMaredForNoCompact,
		metadata.NoneFunc,
		2,  // blockFilesConcurrency
		10, // compactBlocksFetchConcurrency
	)

	bComp, err := NewBucketCompactor(logger, sy, grouper, planner, comp, dir, bkt, 2, true)
	testutil.Ok(t, err)

	// Sync and verify blocks are picked up.
	testutil.Ok(t, sy.SyncMetas(ctx))
	metas := sy.Metas()
	t.Logf("Synced %d metas", len(metas))
	if len(metas) == 0 {
		t.Fatal("No metas synced - blocks not recognized")
	}

	// Check what groups are formed.
	groups, err := grouper.Groups(metas)
	testutil.Ok(t, err)
	t.Logf("Groups formed: %d", len(groups))
	for _, grp := range groups {
		t.Logf("  Group %s: %d IDs", grp.Key(), len(grp.IDs()))
	}

	// Run compaction — the 50 blocks are created with overlapping time ranges
	// and same ext labels, so they will be grouped together. With compactBlocksFetchConcurrency=10,
	// up to 10 errgroup goroutines in compact() will run concurrently,
	// all writing to the shared named return err variable.
	t.Logf("Running compaction...")
	err = bComp.Compact(ctx)
	t.Logf("Compact result err=%v", err)
}
