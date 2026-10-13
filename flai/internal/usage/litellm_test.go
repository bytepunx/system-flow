package usage

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"
)

// The fixtures under testdata/litellm are /spend/logs answers in the shape
// of LiteLLM 1.104.2's spend log table, with the fields S-0356's check read
// (flai/internal/serve/gateway_test.go): spend-logs.json holds the unified
// route's rows, session_id and end_user both Claude Code's session ID, for
// the run (three calls between 09:00 and 09:10 summing to $0.06), for the
// same key's runs before and after it, and for another story's key inside
// it; spend-logs-passthrough.json holds the pass-through's, whose session_id
// is the proxy's own and whose end_user is Claude Code's; unauthorized.json
// is the proxy's refusal of a virtual key.

const (
	// testMaster is the proxy's master key, the only one that reads the log.
	testMaster = "sk-litellm-master-0123456789abcdef"
	// testAgentKey is the virtual key the run's calls used; its SHA-256 is
	// the api_key of the run's rows in the fixtures.
	testAgentKey = "sk-agent-S-0357-virtual-0123456789"
	// theRun is the run's Claude Code session ID in spend-logs.json.
	theRun = "9a1e5f02-the-run"
)

// fakeLiteLLM answers /spend/logs as a LiteLLM proxy does, from a fixture.
type fakeLiteLLM struct {
	fixture string
	// truncated marks the answer as the proxy's cap of rows.
	truncated bool
	// hang holds every answer until the request is cancelled.
	hang bool

	mu       sync.Mutex
	requests int
	query    url.Values
}

func (f *fakeLiteLLM) serve(t *testing.T) *httptest.Server {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		f.mu.Lock()
		f.requests++
		f.query = r.URL.Query()
		f.mu.Unlock()
		switch {
		case f.hang:
			<-r.Context().Done()
		case r.Method != http.MethodGet || r.URL.Path != "/spend/logs":
			w.WriteHeader(http.StatusNotFound)
		case r.Header.Get("Authorization") != "Bearer "+testMaster:
			answerLiteLLM(t, w, http.StatusUnauthorized, "unauthorized.json")
		default:
			if f.truncated {
				w.Header().Set("x-litellm-spend-logs-truncated", "true")
			}
			answerLiteLLM(t, w, http.StatusOK, f.fixture)
		}
	}))
	t.Cleanup(srv.Close)
	return srv
}

// asked is how many requests the fake answered and the last one's query.
func (f *fakeLiteLLM) asked() (int, url.Values) {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.requests, f.query
}

func answerLiteLLM(t *testing.T, w http.ResponseWriter, status int, fixture string) {
	body, err := os.ReadFile(filepath.Join("testdata", "litellm", fixture))
	if err != nil {
		t.Errorf("read fixture %s: %v", fixture, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = w.Write(body)
}

// keyHash is the SHA-256 in hex of key, as the proxy files rows under.
func keyHash(key string) string {
	sum := sha256.Sum256([]byte(key))
	return hex.EncodeToString(sum[:])
}

// runWindow is the run's in spend-logs.json, 09:00 to 09:10 on 2026-10-13.
func runWindow() (from, to time.Time) {
	from = time.Date(2026, 10, 13, 9, 0, 0, 0, time.UTC)
	return from, from.Add(10 * time.Minute)
}

func TestLiteLLMJoinsBySession(t *testing.T) {
	f := &fakeLiteLLM{fixture: "spend-logs.json"}
	srv := f.serve(t)
	from, to := runWindow()
	spent, err := NewLiteLLM(srv.URL, testMaster, srv.Client()).Run(context.Background(), RunRef{Sessions: []string{"e07b9d33-unrelated", theRun}, From: from, To: to})
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if !near(spent.Cost, 0.06) || spent.Calls != 3 || spent.By != SpentBySession {
		t.Errorf("Run = %+v; want $0.06 over 3 calls by session", spent)
	}
	_, q := f.asked()
	want := url.Values{"start_date": {"2026-10-13"}, "end_date": {"2026-10-14"}, "summarize": {"false"}}
	if q.Encode() != want.Encode() {
		t.Errorf("query = %s; want %s", q.Encode(), want.Encode())
	}
}

// A run that resumes a session keeps its ID, so the rows of the session that
// started before the run are an earlier run's.
func TestLiteLLMJoinsBySessionWithinTheRunsWindow(t *testing.T) {
	f := &fakeLiteLLM{fixture: "spend-logs.json"}
	srv := f.serve(t)
	from, to := runWindow()
	from = from.Add(2 * time.Minute)
	spent, err := NewLiteLLM(srv.URL, testMaster, srv.Client()).Run(context.Background(), RunRef{Sessions: []string{theRun}, From: from, To: to})
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if !near(spent.Cost, 0.05) || spent.Calls != 2 || spent.By != SpentBySession {
		t.Errorf("Run = %+v; want $0.05 over the 2 calls from 09:02 by session", spent)
	}
}

// On the pass-through the proxy makes a session_id of its own, and Claude
// Code's session ID is only in end_user (S-0356).
func TestLiteLLMJoinsByEndUser(t *testing.T) {
	f := &fakeLiteLLM{fixture: "spend-logs-passthrough.json"}
	srv := f.serve(t)
	from := time.Date(2026, 10, 12, 9, 0, 0, 0, time.UTC)
	spent, err := NewLiteLLM(srv.URL, testMaster, srv.Client()).Run(context.Background(), RunRef{Sessions: []string{"b5a0c4d9-pass-through-run"}, From: from, To: from.Add(5 * time.Minute)})
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if !near(spent.Cost, 0.75) || spent.Calls != 2 || spent.By != SpentBySession {
		t.Errorf("Run = %+v; want $0.75 over 2 calls by session", spent)
	}
}

// The window join takes only the key's rows inside the window: not the same
// key's runs before and after it, and not another key's rows inside it,
// which the reader drops itself even though the proxy filters by the key.
func TestLiteLLMJoinsByWindowAndKeyHash(t *testing.T) {
	from, to := runWindow()
	for _, c := range []struct {
		name     string
		sessions []string
	}{
		{"no session", nil},
		{"no row for the session", []string{"f00dfeed-not-logged"}},
	} {
		f := &fakeLiteLLM{fixture: "spend-logs.json"}
		srv := f.serve(t)
		spent, err := NewLiteLLM(srv.URL+"/v1/", testMaster, srv.Client()).Run(context.Background(), RunRef{Sessions: c.sessions, KeyHash: keyHash(testAgentKey), From: from, To: to})
		if err != nil {
			t.Fatalf("%s: Run: %v", c.name, err)
		}
		if !near(spent.Cost, 0.06) || spent.Calls != 3 || spent.By != SpentByWindow {
			t.Errorf("%s: Run = %+v; want $0.06 over 3 calls by window", c.name, spent)
		}
		if _, q := f.asked(); q.Get("api_key") != keyHash(testAgentKey) {
			t.Errorf("%s: api_key = %q; want the agent key's hash", c.name, q.Get("api_key"))
		}
	}
}

func TestLiteLLMHoldsNoSpendForTheRun(t *testing.T) {
	from, to := runWindow()
	for _, c := range []struct {
		name    string
		fixture string
		run     RunRef
	}{
		{"no row at all", "spend-logs-empty.json", RunRef{Sessions: []string{theRun}, KeyHash: keyHash(testAgentKey), From: from, To: to}},
		{"no row for the session", "spend-logs.json", RunRef{Sessions: []string{"f00dfeed-not-logged"}, From: from, To: to}},
		{"no row of the key in the window", "spend-logs.json", RunRef{KeyHash: keyHash(testAgentKey), From: to.Add(time.Hour), To: to.Add(2 * time.Hour)}},
		{"no row of the key", "spend-logs.json", RunRef{KeyHash: keyHash("sk-unused"), From: from, To: to}},
	} {
		f := &fakeLiteLLM{fixture: c.fixture}
		srv := f.serve(t)
		spent, err := NewLiteLLM(srv.URL, testMaster, srv.Client()).Run(context.Background(), c.run)
		if !errors.Is(err, ErrNoSpend) || spent != nil {
			t.Errorf("%s: Run = %+v, %v; want ErrNoSpend", c.name, spent, err)
		}
	}
}

// Without a session or the key's hash the window alone could take another
// story's rows, so the reader asks nothing.
func TestLiteLLMWithNoSessionAndNoKeyHashAsksNothing(t *testing.T) {
	f := &fakeLiteLLM{fixture: "spend-logs.json"}
	srv := f.serve(t)
	from, to := runWindow()
	_, err := NewLiteLLM(srv.URL, testMaster, srv.Client()).Run(context.Background(), RunRef{From: from, To: to})
	if requests, _ := f.asked(); !errors.Is(err, ErrNoSpend) || requests != 0 {
		t.Errorf("Run: %v after %d requests; want ErrNoSpend after none", err, requests)
	}
}

// A virtual key gets 401 from /spend/logs (S-0356): an error naming the
// endpoint, the status, and the setting to fix, never the key.
func TestLiteLLMVirtualKeyIsRefusedWithoutTheKey(t *testing.T) {
	f := &fakeLiteLLM{fixture: "spend-logs.json"}
	srv := f.serve(t)
	_, err := NewLiteLLM(srv.URL, testAgentKey, srv.Client()).Run(context.Background(), RunRef{Sessions: []string{theRun}})
	if err == nil || errors.Is(err, ErrNoSpend) {
		t.Fatalf("Run: %v; want an error that is not ErrNoSpend", err)
	}
	for _, want := range []string{srv.URL + "/spend/logs", "401", "spend_key_env"} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("error %q does not name %q", err, want)
		}
	}
	if strings.Contains(err.Error(), testAgentKey) || strings.Contains(err.Error(), testMaster) {
		t.Errorf("error %q carries a key", err)
	}
}

// The proxy answers at most its cap of the most recent rows; a sum over
// them could miss the run's, so it is an error, not a figure.
func TestLiteLLMTruncatedLogIsAnError(t *testing.T) {
	f := &fakeLiteLLM{fixture: "spend-logs.json", truncated: true}
	srv := f.serve(t)
	spent, err := NewLiteLLM(srv.URL, testMaster, srv.Client()).Run(context.Background(), RunRef{Sessions: []string{theRun}})
	if err == nil || errors.Is(err, ErrNoSpend) || spent != nil {
		t.Errorf("Run = %+v, %v; want an error that is not ErrNoSpend", spent, err)
	}
}

func TestLiteLLMStopsAtTheContextsDeadline(t *testing.T) {
	f := &fakeLiteLLM{hang: true}
	srv := f.serve(t)
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	_, err := NewLiteLLM(srv.URL, testMaster, srv.Client()).Run(ctx, RunRef{Sessions: []string{theRun}})
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("Run: %v; want the context's deadline", err)
	}
}
