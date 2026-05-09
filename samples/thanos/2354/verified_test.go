// Copyright (c) The Thanos Authors.
// Licensed under the Apache License 2.0.

package block

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"sync"
	"testing"

	"github.com/go-kit/kit/log"
	"github.com/oklog/ulid"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/thanos-io/thanos/pkg/objstore/inmem"
)

// TestRace_2354 triggers the race condition on MetaFetcher.cached map.
// In the BUG state, MetaFetcher.cached is a map[ulid.ULID]*metadata.Meta that is:
//   - READ by loadMeta() on each worker goroutine within Fetch()
//   - WRITTEN (replaced) by Fetch() at the end: s.cached = cached
//
// Multiple concurrent calls to Fetch() cause a concurrent map read/write data race.
// The FIX introduces BaseFetcher with singleflight.Group to serialize cached updates.
func TestRace_2354(t *testing.T) {
	bkt := inmem.NewBucket()

	// Populate bucket with block directories containing meta.json files.
	// Minimal meta.json format with version=1 to pass validation.
	numBlocks := 30
	var blockIDs []ulid.ULID
	for i := 0; i < numBlocks; i++ {
		id := ulid.MustNew(uint64(1000+i), nil)
		blockIDs = append(blockIDs, id)

		metaJSON := fmt.Sprintf(
			`{"ulid":"%s","minTime":0,"maxTime":1000,"version":1,"thanos":{"source":"test"}}`,
			id.String(),
		)
		bkt.Objects()[path.Join(id.String(), MetaFilename)] = []byte(metaJSON)
	}

	_ = blockIDs

	dir, err := ioutil.TempDir("", "test-race-2354")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	r := prometheus.NewRegistry()
	fetcher, err := NewMetaFetcher(log.NewNopLogger(), 4, bkt, dir, r, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Trigger race: 50+ goroutines calling Fetch concurrently.
	// Each goroutine reads s.cached via loadMeta while another
	// writes s.cached = cached, causing concurrent map R/W.
	var wg sync.WaitGroup
	iterations := 200

	for g := 0; g < 50; g++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < iterations; i++ {
				ctx := context.Background()
				_, _, _ = fetcher.Fetch(ctx)
			}
		}()
	}

	wg.Wait()
	t.Logf("Completed %d goroutines x %d iterations without deadlock", 50, iterations)
}
