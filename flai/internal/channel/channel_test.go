package channel

import (
	"context"
	"crypto/hmac"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/coder/websocket"

	"github.com/bytepunx/system-flow/flai/internal/perf"
)

// fakeDashboard is the other end: it answers hello, checks the proof, and
// then lets a test send requests down each proven connection.
type fakeDashboard struct {
	t      *testing.T
	key    []byte
	srv    *httptest.Server
	mu     sync.Mutex
	conns  []*websocket.Conn
	proven chan *websocket.Conn
	failed chan string
}

func newFakeDashboard(t *testing.T, key string) *fakeDashboard {
	t.Helper()
	d := &fakeDashboard{t: t, key: []byte(key), proven: make(chan *websocket.Conn, 8), failed: make(chan string, 8)}
	d.srv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != Path {
			http.NotFound(w, r)
			return
		}
		c, err := websocket.Accept(w, r, nil)
		if err != nil {
			return
		}
		ctx := r.Context()
		_, data, err := c.Read(ctx)
		if err != nil {
			return
		}
		var hello message
		var hp helloParams
		_ = json.Unmarshal(data, &hello)
		_ = json.Unmarshal(hello.Params, &hp)
		mine := nonce()
		res, _ := json.Marshal(helloResult{Nonce: mine, Proof: Proof(d.key, "dashboard", hp.Nonce, mine), Dashboard: "test"})
		out, _ := json.Marshal(message{JSONRPC: "2.0", ID: hello.ID, Result: res})
		_ = c.Write(ctx, websocket.MessageText, out)
		_, data, err = c.Read(ctx)
		if err != nil {
			d.failed <- "no proof: " + err.Error()
			return
		}
		var prove message
		var pp struct{ Proof string }
		_ = json.Unmarshal(data, &prove)
		_ = json.Unmarshal(prove.Params, &pp)
		if prove.Method != "hello.prove" || !hmac.Equal([]byte(pp.Proof), []byte(Proof(d.key, "flai", mine, hp.Nonce))) {
			d.failed <- "bad proof"
			_ = c.Close(websocket.StatusPolicyViolation, "bad proof")
			return
		}
		d.mu.Lock()
		d.conns = append(d.conns, c)
		d.mu.Unlock()
		d.proven <- c
		<-ctx.Done()
	}))
	t.Cleanup(d.srv.Close)
	return d
}

func (d *fakeDashboard) ask(c *websocket.Conn, id, method, params string) message {
	d.t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	req, _ := json.Marshal(message{JSONRPC: "2.0", ID: json.RawMessage(id), Method: method, Params: json.RawMessage(params)})
	if err := c.Write(ctx, websocket.MessageText, req); err != nil {
		d.t.Fatal(err)
	}
	_, data, err := c.Read(ctx)
	if err != nil {
		d.t.Fatal(err)
	}
	var m message
	_ = json.Unmarshal(data, &m)
	return m
}

func project(t *testing.T) Project {
	t.Helper()
	root := t.TempDir()
	manifest := "version: 1\nname: Harbour\nkey: harbour\nowner: olive\nlayout:\n  design: design\n  docs: docs\n  wip: wip\n"
	if err := os.WriteFile(filepath.Join(root, "system-flow.yaml"), []byte(manifest), 0o644); err != nil {
		t.Fatal(err)
	}
	return Project{Key: "harbour", Name: "Harbour", Root: root}
}

// ProjectInfo and testMethods stand in for the table flai serve offers
// (internal/hostapi), which this package must not import.
type ProjectInfo struct {
	Name   string            `json:"name"`
	Owner  string            `json:"owner"`
	Layout map[string]string `json:"layout"`
	Flai   string            `json:"flai"`
}

func testMethods() map[string]Method {
	return map[string]Method{"project.info": func(_ context.Context, p Project, _ json.RawMessage) (any, *Error) {
		return ProjectInfo{Name: p.Name, Owner: "olive", Layout: map[string]string{"wip": "wip"}, Flai: "test"}, nil
	}}
}

func startClient(t *testing.T, url, key string, p Project) *Client {
	t.Helper()
	c := &Client{URL: url, Key: []byte(key), Project: p, Methods: testMethods(), Version: "test",
		PingEvery: 50 * time.Millisecond, MinBackoff: 10 * time.Millisecond, MaxBackoff: 40 * time.Millisecond}
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() { c.Run(ctx); close(done) }()
	t.Cleanup(func() { cancel(); <-done })
	return c
}

func waitProven(t *testing.T, d *fakeDashboard) *websocket.Conn {
	t.Helper()
	select {
	case c := <-d.proven:
		return c
	case msg := <-d.failed:
		t.Fatalf("handshake failed: %s", msg)
	case <-time.After(5 * time.Second):
		t.Fatal("no connection")
	}
	return nil
}

func TestHandshakeAndProjectInfo(t *testing.T) {
	d := newFakeDashboard(t, "s3cret")
	c := startClient(t, d.srv.URL, "s3cret", project(t))
	conn := waitProven(t, d)

	got := d.ask(conn, "1", "project.info", `{"project":"harbour"}`)
	var info ProjectInfo
	if got.Error != nil || json.Unmarshal(got.Result, &info) != nil || info.Name != "Harbour" || info.Owner != "olive" || info.Layout["wip"] != "wip" || info.Flai != "test" {
		t.Errorf("project.info: %+v %s", got.Error, got.Result)
	}
	if st := c.State(); !st.Connected || st.Since == "" || st.LastError != "" {
		t.Errorf("state: %+v", st)
	}
	if got := d.ask(conn, "2", "shell.run", `{"project":"harbour","argv":["id"]}`); got.Error == nil || got.Error.Code != CodeMethodNotFound {
		t.Errorf("a method flai does not offer: %+v", got)
	}
	if got := d.ask(conn, "3", "project.info", `{"project":"lighthouse"}`); got.Error == nil || got.Error.Code != CodeUnknownProject {
		t.Errorf("a project this connection does not serve: %+v", got)
	}
	if got := d.ask(conn, "4", "project.info", `[]`); got.Error == nil || got.Error.Code != CodeInvalidParams {
		t.Errorf("params that are not an object: %+v", got)
	}
}

func TestAnImpostorGetsNothingAndIsNamed(t *testing.T) {
	d := newFakeDashboard(t, "not-the-key")
	c := startClient(t, d.srv.URL, "s3cret", project(t))
	select {
	case msg := <-d.failed: // the client hung up instead of proving itself
		if !strings.Contains(msg, "no proof") {
			t.Errorf("the client answered an impostor: %s", msg)
		}
	case <-d.proven:
		t.Fatal("the client proved itself to a server that did not hold the key")
	case <-time.After(5 * time.Second):
		t.Fatal("nothing happened")
	}
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) && !strings.Contains(c.State().LastError, "did not prove") {
		time.Sleep(10 * time.Millisecond)
	}
	if st := c.State(); st.Connected || !strings.Contains(st.LastError, "did not prove") {
		t.Errorf("state: %+v", st)
	}
}

func TestReconnectsAfterTheDashboardGoesAway(t *testing.T) {
	d := newFakeDashboard(t, "s3cret")
	c := startClient(t, d.srv.URL, "s3cret", project(t))
	first := waitProven(t, d)
	_ = first.Close(websocket.StatusGoingAway, "restart")
	second := waitProven(t, d)
	if got := d.ask(second, "1", "project.info", `{"project":"harbour"}`); got.Error != nil {
		t.Errorf("after the reconnect: %+v", got.Error)
	}
	// Attempts counts failures since the last proven connection.
	if st := c.State(); !st.Connected || st.Attempts != 0 {
		t.Errorf("state: %+v", st)
	}
}

// S-0184: a dashboard that refuses this flai because another holds the
// project is not dialled again at once, and the state says why.
func TestBacksOffWhenTheDashboardSaysTheProjectIsHeld(t *testing.T) {
	d := newFakeDashboard(t, "s3cret")
	c := &Client{URL: d.srv.URL, Key: []byte("s3cret"), Project: project(t), Methods: testMethods(), Version: "test",
		PingEvery: 50 * time.Millisecond, MinBackoff: 10 * time.Millisecond, MaxBackoff: 40 * time.Millisecond, HeldBackoff: 600 * time.Millisecond}
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() { c.Run(ctx); close(done) }()
	t.Cleanup(func() { cancel(); <-done })

	first := waitProven(t, d)
	refused := time.Now()
	_ = first.Close(CloseHeld, "the project is held by a flai that answers")
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) && !strings.Contains(c.State().LastError, "another flai") {
		time.Sleep(5 * time.Millisecond)
	}
	if st := c.State(); st.Connected || !strings.Contains(st.LastError, "another flai serves the project") || !strings.Contains(st.LastError, "trying again in 600ms") {
		t.Errorf("state: %+v", st)
	}
	select {
	case <-d.proven:
		// jitter spreads the wait over [HeldBackoff/2, HeldBackoff)
		if waited := time.Since(refused); waited < 250*time.Millisecond {
			t.Errorf("dialled again after %v, not after the held back-off", waited)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("never dialled again")
	}

	// any other close is retried at once, as before
	second := lastProven(t, d)
	again := time.Now()
	_ = second.Close(websocket.StatusGoingAway, "restart")
	waitProven(t, d)
	if waited := time.Since(again); waited > 250*time.Millisecond {
		t.Errorf("an ordinary close waited %v", waited)
	}
}

// lastProven is the connection the dashboard proved last: a select that
// took it from d.proven left nothing there to wait for.
func lastProven(t *testing.T, d *fakeDashboard) *websocket.Conn {
	t.Helper()
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.conns[len(d.conns)-1]
}

func TestAnAnswerOverTheCapIsAnError(t *testing.T) {
	d := newFakeDashboard(t, "s3cret")
	p := project(t)
	c := &Client{URL: d.srv.URL, Key: []byte("s3cret"), Project: p, Version: "test", MaxBytes: 4096,
		Methods: map[string]Method{"big": func(context.Context, Project, json.RawMessage) (any, *Error) {
			return strings.Repeat("x", 8192), nil
		}}}
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() { c.Run(ctx); close(done) }()
	t.Cleanup(func() { cancel(); <-done })
	conn := waitProven(t, d)
	if got := d.ask(conn, "1", "big", `{"project":"harbour"}`); got.Error == nil || !strings.Contains(got.Error.Message, "cap") {
		t.Errorf("an answer over the cap: %+v", got)
	}
}

func TestACancelledRequestStops(t *testing.T) {
	d := newFakeDashboard(t, "s3cret")
	stopped := make(chan struct{})
	c := &Client{URL: d.srv.URL, Key: []byte("s3cret"), Project: project(t), Version: "test",
		Methods: map[string]Method{"slow": func(ctx context.Context, _ Project, _ json.RawMessage) (any, *Error) {
			<-ctx.Done()
			close(stopped)
			return nil, &Error{Code: CodeInternal, Message: "cancelled"}
		}}}
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() { c.Run(ctx); close(done) }()
	t.Cleanup(func() { cancel(); <-done })
	conn := waitProven(t, d)
	wctx, wdone := context.WithTimeout(context.Background(), 5*time.Second)
	defer wdone()
	_ = conn.Write(wctx, websocket.MessageText, []byte(`{"jsonrpc":"2.0","id":7,"method":"slow","params":{"project":"harbour"}}`))
	time.Sleep(50 * time.Millisecond)
	_ = conn.Write(wctx, websocket.MessageText, []byte(`{"jsonrpc":"2.0","method":"$/cancel","params":{"id":7}}`))
	select {
	case <-stopped:
	case <-time.After(3 * time.Second):
		t.Fatal("the method was not cancelled")
	}
}

// lockedBuffer is a log destination a test reads while the client writes.
type lockedBuffer struct {
	mu sync.Mutex
	b  strings.Builder
}

func (l *lockedBuffer) Write(p []byte) (int, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.b.Write(p)
}

func (l *lockedBuffer) String() string {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.b.String()
}

func TestEachRequestIsTimedBeneathTheTransport(t *testing.T) {
	d := newFakeDashboard(t, "s3cret")
	var logs lockedBuffer
	c := &Client{URL: d.srv.URL, Key: []byte("s3cret"), Project: project(t), Version: "test", Slow: time.Nanosecond,
		Logger: slog.New(slog.NewJSONHandler(&logs, nil)),
		Methods: map[string]Method{"board.get": func(ctx context.Context, _ Project, _ json.RawMessage) (any, *Error) {
			defer perf.Track(ctx, "items.list")()
			return map[string]string{"board": "ok"}, nil
		}}}
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() { c.Run(ctx); close(done) }()
	t.Cleanup(func() { cancel(); <-done })
	conn := waitProven(t, d)
	if got := d.ask(conn, "1", "board.get", `{"project":"harbour"}`); got.Error != nil {
		t.Fatalf("board.get: %+v", got.Error)
	}
	var ev map[string]any
	for _, line := range strings.Split(logs.String(), "\n") {
		if strings.Contains(line, `"request answered"`) {
			if err := json.Unmarshal([]byte(line), &ev); err != nil {
				t.Fatal(err)
			}
		}
	}
	if ev == nil {
		t.Fatalf("no request answered event in %s", logs.String())
	}
	phases, _ := ev["phases"].(string)
	if ev["level"] != "INFO" || ev["transport"] != "channel" || ev["method"] != "board.get" || ev["project"] != "harbour" ||
		ev["bytes"] != float64(len(`{"board":"ok"}`)) || ev["duration_ms"] == nil ||
		!strings.Contains(phases, "items.list=") || !strings.Contains(phases, "encode=") {
		t.Errorf("event %v", ev)
	}
}
