package usage

import (
	"context"
	"errors"
	"math"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"
)

// The fixtures under testdata/openrouter are OpenRouter's answers in the
// shapes S-0351's check read them (flai/internal/serve/gateway_test.go):
// /api/v1/key's data.usage, from a current copy and one that lags it by the
// $0.2791 against $0.3384 the check saw, and /api/v1/generation's
// data.total_cost for two generations that sum to the $0.060022 the check's
// key rose by. The 400 is what it answered for Claude Code's message IDs.

const testKey = "sk-or-v1-test-0123456789abcdef"

// fakeOpenRouter answers as OpenRouter does, from the fixtures.
type fakeOpenRouter struct {
	// keys are the fixtures answered to successive reads of /v1/key, the
	// last repeated.
	keys []string
	// generations are the fixtures answered by generation ID; an ID in
	// pending is answered 404, any other 400.
	generations map[string]string
	pending     map[string]bool
	// hang holds every answer until the request is cancelled.
	hang bool

	mu       sync.Mutex
	keyReads int
	requests int
}

func (f *fakeOpenRouter) serve(t *testing.T) *httptest.Server {
	t.Helper()
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v1/key", func(w http.ResponseWriter, r *http.Request) {
		f.mu.Lock()
		i := min(f.keyReads, len(f.keys)-1)
		f.keyReads++
		f.mu.Unlock()
		answer(t, w, http.StatusOK, f.keys[i])
	})
	mux.HandleFunc("GET /api/v1/generation", func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		switch {
		case f.generations[id] != "":
			answer(t, w, http.StatusOK, f.generations[id])
		case f.pending[id]:
			answer(t, w, http.StatusNotFound, "generation-not-found.json")
		default:
			answer(t, w, http.StatusBadRequest, "invalid-generation.json")
		}
	})
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		f.mu.Lock()
		f.requests++
		f.mu.Unlock()
		if f.hang {
			<-r.Context().Done()
			return
		}
		if r.Header.Get("Authorization") != "Bearer "+testKey {
			answer(t, w, http.StatusUnauthorized, "unauthorized.json")
			return
		}
		mux.ServeHTTP(w, r)
	}))
	t.Cleanup(srv.Close)
	return srv
}

// counts are the reads of /v1/key and the requests the fake has answered.
func (f *fakeOpenRouter) counts() (keyReads, requests int) {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.keyReads, f.requests
}

func answer(t *testing.T, w http.ResponseWriter, status int, fixture string) {
	body, err := os.ReadFile(filepath.Join("testdata", "openrouter", fixture))
	if err != nil {
		t.Errorf("read fixture %s: %v", fixture, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = w.Write(body)
}

// reader is OpenRouter's reader of srv at base, with key, reading the key
// without a pause.
func reader(srv *httptest.Server, base, key string) *OpenRouter {
	o := NewOpenRouter(srv.URL+base, key, srv.Client())
	o.pause = 0
	return o
}

func near(a, b float64) bool { return math.Abs(a-b) < 1e-9 }

func TestOpenRouterJoinsByGeneration(t *testing.T) {
	for _, base := range []string{"/api", "/api/", "/api/v1", "/api/v1/"} {
		t.Run(base, func(t *testing.T) {
			f := &fakeOpenRouter{generations: map[string]string{"gen-a": "generation-a.json", "gen-b": "generation-b.json"}}
			srv := f.serve(t)
			got, err := reader(srv, base, testKey).Run(context.Background(), RunRef{Generations: []string{"gen-a", "gen-b"}})
			if err != nil {
				t.Fatalf("Run: %v", err)
			}
			if !near(got.Cost, 0.060022) || got.Calls != 2 || got.By != SpentByGeneration {
				t.Errorf("Run = %+v; want $0.060022 over 2 calls by generation", got)
			}
		})
	}
}

// S-0351: /api/v1/key answers from copies that lag one another, so a
// reading is the highest of three reads, whichever copy answers first.
func TestOpenRouterJoinsByKeyTakingTheHighestRead(t *testing.T) {
	for _, keys := range [][]string{
		{"key-lagging.json", "key-current.json", "key-lagging.json"},
		{"key-current.json", "key-lagging.json", "key-lagging.json"},
		{"key-lagging.json", "key-lagging.json", "key-current.json"},
	} {
		f := &fakeOpenRouter{keys: keys}
		srv := f.serve(t)
		before := 0.278378
		got, err := reader(srv, "/api", testKey).Run(context.Background(), RunRef{KeyBefore: &before})
		if err != nil {
			t.Fatalf("%v: Run: %v", keys, err)
		}
		if !near(got.Cost, 0.060022) || got.Calls != 0 || got.By != SpentByKey {
			t.Errorf("%v: Run = %+v; want $0.060022 by key, calls not counted", keys, got)
		}
		if reads, _ := f.counts(); reads != keyReads {
			t.Errorf("%v: /v1/key read %d times; want %d", keys, reads, keyReads)
		}
	}
}

func TestOpenRouterKeyUsageIsTheHighestOfThreeReads(t *testing.T) {
	f := &fakeOpenRouter{keys: []string{"key-current.json", "key-lagging.json"}}
	srv := f.serve(t)
	got, err := reader(srv, "/api/v1", testKey).KeyUsage(context.Background())
	if reads, _ := f.counts(); err != nil || got != 0.3384 || reads != keyReads {
		t.Errorf("KeyUsage = %v, %v after %d reads; want 0.3384 after %d", got, err, reads, keyReads)
	}
}

func TestOpenRouterHoldsNoSpendForTheRun(t *testing.T) {
	current := 0.3384
	for _, c := range []struct {
		name string
		run  RunRef
		want string
	}{
		{"no generation IDs and no reading", RunRef{Sessions: []string{"s-1"}, From: time.Now().Add(-time.Hour), To: time.Now()}, "no generation IDs and no reading"},
		{"the key has not risen", RunRef{KeyBefore: &current}, "has not risen"},
		{"a generation not recorded yet", RunRef{Generations: []string{"gen-a", "gen-later"}}, "1 of the run's 2 generations"},
	} {
		t.Run(c.name, func(t *testing.T) {
			f := &fakeOpenRouter{
				keys:        []string{"key-lagging.json", "key-current.json"},
				generations: map[string]string{"gen-a": "generation-a.json"},
				pending:     map[string]bool{"gen-later": true},
			}
			srv := f.serve(t)
			got, err := reader(srv, "/api", testKey).Run(context.Background(), c.run)
			if !errors.Is(err, ErrNoSpend) || got != nil {
				t.Fatalf("Run = %+v, %v; want ErrNoSpend", got, err)
			}
			if !strings.Contains(err.Error(), c.want) {
				t.Errorf("error %q does not say %q", err, c.want)
			}
		})
	}
}

func TestOpenRouterWithNoReadingAsksNothing(t *testing.T) {
	f := &fakeOpenRouter{}
	srv := f.serve(t)
	if _, err := reader(srv, "/api", testKey).Run(context.Background(), RunRef{}); !errors.Is(err, ErrNoSpend) {
		t.Fatalf("Run: %v; want ErrNoSpend", err)
	}
	if _, requests := f.counts(); requests != 0 {
		t.Errorf("%d requests; want none", requests)
	}
}

// A refused key is an error naming the endpoint and the status, not
// ErrNoSpend, and never carries the key.
func TestOpenRouterRefusedKeyIsAnErrorWithoutTheKey(t *testing.T) {
	const wrong = "sk-or-v1-wrong-fedcba9876543210"
	before := 0.1
	f := &fakeOpenRouter{keys: []string{"key-current.json"}, generations: map[string]string{"gen-a": "generation-a.json"}}
	srv := f.serve(t)
	o := reader(srv, "/api", wrong)
	for _, c := range []struct {
		name     string
		run      RunRef
		endpoint string
	}{
		{"by key", RunRef{KeyBefore: &before}, "/api/v1/key"},
		{"by generation", RunRef{Generations: []string{"gen-a"}}, "/api/v1/generation?id=gen-a"},
	} {
		_, err := o.Run(context.Background(), c.run)
		if err == nil || errors.Is(err, ErrNoSpend) {
			t.Fatalf("%s: Run: %v; want an error that is not ErrNoSpend", c.name, err)
		}
		for _, want := range []string{c.endpoint, "401"} {
			if !strings.Contains(err.Error(), want) {
				t.Errorf("%s: error %q does not name %q", c.name, err, want)
			}
		}
		if strings.Contains(err.Error(), wrong) || strings.Contains(err.Error(), testKey) {
			t.Errorf("%s: error %q carries a key", c.name, err)
		}
	}
}

// Claude Code's message IDs are not OpenRouter's generation IDs, and
// OpenRouter answers 400 for them (S-0351).
func TestOpenRouterInvalidGenerationIsAnError(t *testing.T) {
	f := &fakeOpenRouter{}
	srv := f.serve(t)
	_, err := reader(srv, "/api", testKey).Run(context.Background(), RunRef{Generations: []string{"msg_bdrk_01"}})
	if err == nil || errors.Is(err, ErrNoSpend) || !strings.Contains(err.Error(), "/api/v1/generation?id=msg_bdrk_01 answered 400") {
		t.Errorf("Run: %v; want an error naming the endpoint and 400", err)
	}
}

func TestOpenRouterStopsAtTheContextsDeadline(t *testing.T) {
	before := 0.1
	t.Run("a request", func(t *testing.T) {
		f := &fakeOpenRouter{hang: true}
		srv := f.serve(t)
		ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
		defer cancel()
		_, err := reader(srv, "/api", testKey).Run(ctx, RunRef{KeyBefore: &before})
		if !errors.Is(err, context.DeadlineExceeded) {
			t.Errorf("Run: %v; want the context's deadline", err)
		}
	})
	t.Run("the pause between reads", func(t *testing.T) {
		f := &fakeOpenRouter{keys: []string{"key-current.json"}}
		srv := f.serve(t)
		o := reader(srv, "/api", testKey)
		o.pause = time.Hour
		ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
		defer cancel()
		_, err := o.KeyUsage(ctx)
		if reads, _ := f.counts(); !errors.Is(err, context.DeadlineExceeded) || reads != 1 {
			t.Errorf("KeyUsage: %v after %d reads; want the context's deadline after 1", err, reads)
		}
	})
}
