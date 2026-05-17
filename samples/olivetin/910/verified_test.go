package otoauth2

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	config "github.com/OliveTin/OliveTin/internal/config"
	"golang.org/x/oauth2"
)

// TestRace_910_OAuth2RegisteredStates reproduces PR #910:
// HandleOAuthLogin writes to h.registeredStates while checkOAuthCallbackCookie reads
// from it; bare map without lock → race.
func TestRace_910_OAuth2RegisteredStates(t *testing.T) {
	cfg := &config.Config{
		Security: config.SecurityConfig{ForceSecureCookies: false},
	}
	h := &OAuth2Handler{
		cfg:              cfg,
		registeredStates: map[string]*oauth2State{},
		registeredProviders: map[string]*oauth2.Config{
			"test-provider": {
				ClientID: "x",
				Endpoint: oauth2.Endpoint{
					AuthURL: "https://example.com/auth", TokenURL: "https://example.com/tok",
				},
				RedirectURL: "https://example.com/cb",
			},
		},
	}

	// Pre-seed some states for readers to find.
	for i := 0; i < 50; i++ {
		h.registeredStates[fmt.Sprintf("seed-%d", i)] = &oauth2State{providerName: "test-provider"}
	}

	const N = 20
	const ITERS = 50
	var wg sync.WaitGroup
	// Writers via HandleOAuthLogin (writes to map)
	for i := 0; i < N; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < ITERS; j++ {
				w := httptest.NewRecorder()
				r := httptest.NewRequest("GET", "/?provider=test-provider", nil)
				h.HandleOAuthLogin(w, r)
			}
		}()
	}
	// Readers via checkOAuthCallbackCookie (reads from map)
	for i := 0; i < N; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			for j := 0; j < ITERS; j++ {
				st := fmt.Sprintf("seed-%d", j%50)
				w := httptest.NewRecorder()
				r := httptest.NewRequest("GET", "/?state="+st, nil)
				r.AddCookie(&http.Cookie{Name: "olivetin-sid-oauth", Value: st})
				_, _, _ = h.checkOAuthCallbackCookie(w, r)
			}
		}(i)
	}
	wg.Wait()
}
