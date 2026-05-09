/*
 * Teleport
 * Copyright (C) 2023  Gravitational, Inc.
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package gcp

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"cloud.google.com/go/iam/credentials/apiv1/credentialspb"
	"github.com/googleapis/gax-go/v2"
	"github.com/jonboulle/clockwork"
	"github.com/stretchr/testify/require"
)

// TestRace_65910 triggers the data race on tokenResult and errorResult
// shared variables in getToken().
//
// BUG: getToken() declares tokenResult and errorResult as local variables,
// then starts a goroutine that writes to them via FnCacheGet. The parent
// goroutine reads these variables after selecting on cancelCtx.Done().
// When the context is pre-cancelled, cancelCtx.Done() is immediately closed,
// and the parent reads the (still zero) shared variables while the inner
// goroutine has not yet run or is in the process of writing to them.
// This concurrent read/write without proper synchronization is a data race.
//
// FIX: Replaces shared mutable state with a channel-based result passing
// pattern (result struct sent over buffered channel).
//
// WARNING DATA RACE evidence (from PR #65910 body):
//   Write at handler.go:275 by inner goroutine (tokenResult = FnCacheGet)
//   Read at handler.go:285 by parent (return tokenResult, errorResult)
func TestRace_65910(t *testing.T) {
	// Use a pre-cancelled context so that cancelCtx.Done() is immediately
	// closed when getToken() creates cancelCtx from it. This causes the
	// parent goroutine's select to immediately take the cancelCtx.Done()
	// path and read tokenResult/errorResult while the inner goroutine
	// is concurrently writing to them.
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	fwd, err := newGCPHandler(context.Background(), HandlerConfig{
		Clock: clockwork.NewRealClock(),
		cloudClientGCP: makeTestCloudClient(&testIAMCredentialsClient{
			generateAccessToken: func(ctx context.Context, req *credentialspb.GenerateAccessTokenRequest, opts ...gax.CallOption) (*credentialspb.GenerateAccessTokenResponse, error) {
				// The race is between the goroutine writing tokenResult/errorResult
				// and the parent reading them. Adding a small delay widens the
				// race window to make detection more reliable.
				time.Sleep(time.Microsecond)
				return &credentialspb.GenerateAccessTokenResponse{AccessToken: "ok"}, nil
			},
		}),
	})
	require.NoError(t, err)

	numGoroutines := 50
	iterations := 200

	var wg sync.WaitGroup
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				// Unique service account per iteration to avoid FnCache hits.
				// Each call creates its own inner goroutine that writes to
				// tokenResult/errorResult concurrently with the parent read.
				sa := fmt.Sprintf("race-65910-g%d-i%d-%d", id, j, time.Now().UnixNano())
				_, _ = fwd.getToken(ctx, sa)
			}
		}(i)
	}
	wg.Wait()
}
