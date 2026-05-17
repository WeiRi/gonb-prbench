// Race test for victoriametrics-8258 — AlertManager.send alertsToSend aliases shared alerts
// BUG: alertsToSend := alerts[:0] aliases the input slice; concurrent Send overwrites alerts[i]
// FIX: alertsToSend := make([]Alert, 0, len(alerts)) — independent backing array
package notifier

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/VictoriaMetrics/VictoriaMetrics/lib/promauth"
)

func TestRace_8258_AlertsSlice(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	gen := func(a Alert) string { return "http://gen/" + a.Name }
	am, err := NewAlertManager(srv.URL+alertManagerPath, gen, promauth.HTTPClientConfig{}, nil, 0)
	if err != nil {
		t.Fatalf("NewAlertManager: %v", err)
	}

	// SHARED alerts (non-empty Labels so relabel filter doesn't drop them)
	alerts := []Alert{
		{Name: "a1", Labels: map[string]string{"job": "j1"}},
		{Name: "a2", Labels: map[string]string{"job": "j2"}},
		{Name: "a3", Labels: map[string]string{"job": "j3"}},
	}

	var wg sync.WaitGroup
	const N = 30
	for i := 0; i < N; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 50; j++ {
				_ = am.Send(context.Background(), alerts, nil)
			}
		}()
	}
	wg.Wait()
}
