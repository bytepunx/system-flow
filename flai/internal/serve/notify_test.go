package serve

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/config"
	"github.com/bytepunx/system-flow/flai/internal/threads"
)

func testReminder() Reminder {
	return Reminder{
		Project:         "demo",
		Thread:          "TH-0324",
		Title:           "Which key signs the release?",
		Item:            "S-0232",
		AwaitingSeconds: 7200,
		Since:           time.Date(2026, 10, 20, 9, 30, 0, 0, time.FixedZone("CEST", 2*3600)),
		Recommendation:  &ReminderRecommendation{Author: "orchestrator", FirstLine: "Use the cosign key.", Source: &threads.Source{Path: "design/adrs/0070-x.md", Heading: "Decision"}},
		Href:            ReminderHref("S-0232", "TH-0324"),
	}
}

// wantBody is testReminder's JSON as a receiver decodes it: the ADR's fields
// and flaiover's entry.
var wantBody = map[string]any{
	"project":          "demo",
	"thread":           "TH-0324",
	"title":            "Which key signs the release?",
	"item":             "S-0232",
	"awaiting_seconds": float64(7200),
	"since":            "2026-10-20T07:30:00Z",
	"recommendation": map[string]any{
		"author":     "orchestrator",
		"first_line": "Use the cosign key.",
		"source":     map[string]any{"path": "design/adrs/0070-x.md", "heading": "Decision"},
	},
	"href": "/items/S-0232?thread=TH-0324",
	"entry": map[string]any{
		"key":   "thread:TH-0324",
		"kind":  "thread",
		"title": "Which key signs the release?",
		"href":  "/items/S-0232?thread=TH-0324",
		"at":    "2026-10-20T07:30:00Z",
	},
}

// posted is what a test webhook was sent.
type posted struct {
	method, contentType, userAgent string
	body                           map[string]any
}

// webhook is a test server that records each POST and answers status.
func webhook(t *testing.T, status int) (*httptest.Server, func() []posted) {
	t.Helper()
	var mu sync.Mutex
	var got []posted
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, _ := io.ReadAll(r.Body)
		var body map[string]any
		if err := json.Unmarshal(data, &body); err != nil {
			t.Errorf("the webhook's body is not JSON: %v\n%s", err, data)
		}
		mu.Lock()
		got = append(got, posted{r.Method, r.Header.Get("Content-Type"), r.Header.Get("User-Agent"), body})
		mu.Unlock()
		w.WriteHeader(status)
	}))
	t.Cleanup(srv.Close)
	return srv, func() []posted {
		mu.Lock()
		defer mu.Unlock()
		return append([]posted(nil), got...)
	}
}

// stdinTo is a shell command that writes its standard input to path.
func stdinTo(t *testing.T, path string) string {
	t.Helper()
	if runtime.GOOS == "windows" {
		t.Skip("the command is a POSIX shell's")
	}
	return "cat > '" + path + "'"
}

func readJSON(t *testing.T, path string) map[string]any {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("the command wrote nothing: %v", err)
	}
	if !strings.HasSuffix(string(data), "}\n") {
		t.Errorf("the command's input ends in a newline: %q", data)
	}
	var got map[string]any
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("the command's input is not JSON: %v\n%s", err, data)
	}
	return got
}

// S-0498 (ADR-0205): the webhook gets one POST of the reminder's JSON, as
// application/json from flai, and an answer of 500 is an error naming its
// host alone.
func TestSendPostsTheReminderOnce(t *testing.T) {
	srv, got := webhook(t, http.StatusNoContent)
	if err := Send(context.Background(), config.Notify{URL: srv.URL + "/flai?auth=secret"}, testReminder()); err != nil {
		t.Fatalf("Send: %v", err)
	}
	posts := got()
	if len(posts) != 1 {
		t.Fatalf("one POST, got %d", len(posts))
	}
	p := posts[0]
	if p.method != http.MethodPost || p.contentType != "application/json" || p.userAgent != "flai" {
		t.Errorf("method %s, content-type %q, user-agent %q", p.method, p.contentType, p.userAgent)
	}
	if !reflect.DeepEqual(p.body, wantBody) {
		t.Errorf("body:\n got %v\nwant %v", p.body, wantBody)
	}

	failing, tried := webhook(t, http.StatusInternalServerError)
	err := Send(context.Background(), config.Notify{URL: failing.URL + "/flai?auth=secret"}, testReminder())
	if err == nil || !strings.Contains(err.Error(), "500 Internal Server Error") || !strings.Contains(err.Error(), "TH-0324") {
		t.Errorf("a 500 is an error naming the thread and the answer: %v", err)
	}
	if strings.Contains(err.Error(), "secret") || !strings.Contains(err.Error(), strings.TrimPrefix(failing.URL, "http://")) {
		t.Errorf("the error names the host, never the URL: %v", err)
	}
	if n := len(tried()); n != 1 {
		t.Errorf("one attempt, no retry: %d", n)
	}
}

// A webhook that does not answer within the timeout is an error, and its URL
// is not in it.
func TestSendGivesUpOnASlowWebhook(t *testing.T) {
	release := make(chan struct{})
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		select {
		case <-r.Context().Done():
		case <-release:
		}
	}))
	t.Cleanup(srv.Close)
	t.Cleanup(func() { close(release) })
	n := notifier{client: http.DefaultClient, shell: runShell, urlTimeout: 50 * time.Millisecond, commandTimeout: time.Second}
	err := n.send(context.Background(), config.Notify{URL: srv.URL + "/?token=secret"}, testReminder())
	if err == nil || !strings.Contains(err.Error(), "it gave no answer within 50ms") || strings.Contains(err.Error(), "secret") {
		t.Errorf("a slow webhook is an error without its URL: %v", err)
	}
}

// The command gets the reminder's JSON on its standard input, and one that
// exits non-zero is an error that quotes what it said.
func TestSendPipesTheReminderToTheCommand(t *testing.T) {
	out := filepath.Join(t.TempDir(), "reminder.json")
	if err := Send(context.Background(), config.Notify{Command: stdinTo(t, out)}, testReminder()); err != nil {
		t.Fatalf("Send: %v", err)
	}
	if got := readJSON(t, out); !reflect.DeepEqual(got, wantBody) {
		t.Errorf("stdin:\n got %v\nwant %v", got, wantBody)
	}

	err := Send(context.Background(), config.Notify{Command: "echo no such topic >&2; exit 3"}, testReminder())
	if err == nil || !strings.Contains(err.Error(), "exit status 3: no such topic") || !strings.Contains(err.Error(), "TH-0324") {
		t.Errorf("a command that exits non-zero is an error quoting it: %v", err)
	}
	if strings.Contains(err.Error(), "echo") {
		t.Errorf("the error does not repeat the command, which may hold a secret: %v", err)
	}
}

// A command that runs past the timeout is stopped and is an error.
func TestSendStopsASlowCommand(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("the command is a POSIX shell's")
	}
	n := notifier{client: http.DefaultClient, shell: runShell, urlTimeout: time.Second, commandTimeout: 100 * time.Millisecond}
	start := time.Now()
	err := n.send(context.Background(), config.Notify{Command: "sleep 10"}, testReminder())
	if err == nil || !strings.Contains(err.Error(), "it did not end within 100ms") {
		t.Errorf("a slow command is an error: %v", err)
	}
	if took := time.Since(start); took > 5*time.Second {
		t.Errorf("the command was not stopped: %s", took)
	}
}

// With both set, both are sent, the command though the webhook failed, and
// both errors are returned; with neither, nothing is sent.
func TestSendBothAndNeither(t *testing.T) {
	ok, got := webhook(t, http.StatusOK)
	out := filepath.Join(t.TempDir(), "reminder.json")
	if err := Send(context.Background(), config.Notify{URL: ok.URL, Command: stdinTo(t, out)}, testReminder()); err != nil {
		t.Fatalf("Send: %v", err)
	}
	if len(got()) != 1 || !reflect.DeepEqual(readJSON(t, out), wantBody) {
		t.Errorf("both are sent: %d posts", len(got()))
	}

	failing, _ := webhook(t, http.StatusBadGateway)
	_ = os.Remove(out)
	err := Send(context.Background(), config.Notify{URL: failing.URL, Command: stdinTo(t, out) + "; exit 1"}, testReminder())
	if err == nil || !strings.Contains(err.Error(), "502 Bad Gateway") || !strings.Contains(err.Error(), "run the notify command") {
		t.Errorf("both errors are returned: %v", err)
	}
	if _, statErr := os.Stat(out); statErr != nil {
		t.Errorf("the command runs though the webhook failed: %v", statErr)
	}

	called := false
	n := notifier{client: http.DefaultClient, shell: func(context.Context, string, []byte) ([]byte, error) {
		called = true
		return nil, nil
	}, urlTimeout: time.Second, commandTimeout: time.Second}
	if err := n.send(context.Background(), config.Notify{}, testReminder()); err != nil || called {
		t.Errorf("with no channel nothing is sent: %v, command run %v", err, called)
	}
}

// Without a pending recommendation the reminder carries none; a thread on a
// document links to the threads page.
func TestReminderWithoutRecommendation(t *testing.T) {
	r := testReminder()
	r.Recommendation, r.Item, r.Href = nil, "", ReminderHref("", "TH-0001")
	data, err := json.Marshal(r)
	if err != nil {
		t.Fatal(err)
	}
	var got map[string]any
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatal(err)
	}
	if _, has := got["recommendation"]; has {
		t.Errorf("no recommendation is pending: %s", data)
	}
	if got["href"] != "/threads?thread=TH-0001" || got["item"] != "" {
		t.Errorf("a thread on a document: %s", data)
	}
}
