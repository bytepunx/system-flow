package channel

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"sync/atomic"
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
	// keys verify what signer signs: the stamp it answers with by default.
	keys   []*ecdsa.PublicKey
	signer *ecdsa.PrivateKey
	// release is the stamp it answers hello with, nil for none; refuse,
	// when set, is the fault it closes hello with instead; hellos are the
	// hellos it read.
	release *Release
	refuse  string
	hellos  []helloParams
}

// sign makes a release stamp of statement with the dashboard's key.
func (d *fakeDashboard) sign(statement string) *Release {
	d.t.Helper()
	sum := sha256.Sum256([]byte(statement))
	sig, err := ecdsa.SignASN1(rand.Reader, d.signer, sum[:])
	if err != nil {
		d.t.Fatal(err)
	}
	return &Release{Statement: statement, Signature: base64.StdEncoding.EncodeToString(sig)}
}

func (d *fakeDashboard) answer(r *Release, refuse string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.release, d.refuse = r, refuse
}

func newFakeDashboard(t *testing.T, key string) *fakeDashboard {
	t.Helper()
	signer, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	d := &fakeDashboard{t: t, key: []byte(key), proven: make(chan *websocket.Conn, 8), failed: make(chan string, 8),
		keys: []*ecdsa.PublicKey{&signer.PublicKey}, signer: signer}
	d.release = d.sign("flaiover 2.3.0 1a2b3c4")
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
		d.mu.Lock()
		d.hellos = append(d.hellos, hp)
		release, refuse := d.release, d.refuse
		d.mu.Unlock()
		if refuse != "" {
			_ = c.Close(CloseUnsigned, UnsignedReason+": "+refuse)
			return
		}
		mine := nonce()
		res, _ := json.Marshal(helloResult{Nonce: mine, Proof: Proof(d.key, "dashboard", hp.Nonce, mine), Dashboard: "test", Release: release})
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

func startClient(t *testing.T, url, key string, p Project, keys []*ecdsa.PublicKey) *Client {
	t.Helper()
	c := &Client{URL: url, Key: []byte(key), Project: p, Methods: testMethods(), Version: "test", ReleaseKeys: keys,
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
	c := startClient(t, d.srv.URL, "s3cret", project(t), d.keys)
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
	c := startClient(t, d.srv.URL, "s3cret", project(t), d.keys)
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
	c := startClient(t, d.srv.URL, "s3cret", project(t), d.keys)
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
	c := &Client{URL: d.srv.URL, Key: []byte("s3cret"), ReleaseKeys: d.keys, Project: project(t), Methods: testMethods(), Version: "test",
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
	c := &Client{URL: d.srv.URL, Key: []byte("s3cret"), ReleaseKeys: d.keys, Project: p, Version: "test", MaxBytes: 4096,
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
	c := &Client{URL: d.srv.URL, Key: []byte("s3cret"), ReleaseKeys: d.keys, Project: project(t), Version: "test",
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
	c := &Client{URL: d.srv.URL, Key: []byte("s3cret"), ReleaseKeys: d.keys, Project: project(t), Version: "test", Slow: time.Nanosecond,
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

// refusals are the "connection refused as unsigned" events in logs.
func refusals(t *testing.T, logs *lockedBuffer) []map[string]any {
	t.Helper()
	var out []map[string]any
	for _, line := range strings.Split(logs.String(), "\n") {
		var ev map[string]any
		if json.Unmarshal([]byte(line), &ev) == nil && ev["msg"] == "connection refused as unsigned" {
			out = append(out, ev)
		}
	}
	return out
}

func waitState(t *testing.T, c *Client, what string, ok func(State) bool) State {
	t.Helper()
	deadline := time.Now().Add(3 * time.Second)
	for time.Now().Before(deadline) {
		if st := c.State(); ok(st) {
			return st
		}
		time.Sleep(5 * time.Millisecond)
	}
	t.Fatalf("%s: state %+v", what, c.State())
	return State{}
}

func runClient(t *testing.T, c *Client) {
	t.Helper()
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() { c.Run(ctx); close(done) }()
	t.Cleanup(func() { cancel(); <-done })
}

// S-0237: hello carries flai's stamp under protocol 2, and a dashboard whose
// stamp verifies and names flaiover is served, its version and that it is
// signed in the state.
func TestASignedDashboardConnects(t *testing.T) {
	d := newFakeDashboard(t, "s3cret")
	own := Release{Statement: "flai 1.43.0 a1825d8", Signature: "c2lnbmF0dXJl"}
	c := &Client{URL: d.srv.URL, Key: []byte("s3cret"), Project: project(t), Methods: testMethods(), Version: "1.43.0",
		Release: own, ReleaseKeys: d.keys, PingEvery: 50 * time.Millisecond, MinBackoff: 10 * time.Millisecond, MaxBackoff: 40 * time.Millisecond}
	runClient(t, c)
	waitProven(t, d)
	d.mu.Lock()
	hello := d.hellos[0]
	d.mu.Unlock()
	if hello.Protocol != 2 || hello.Release == nil || *hello.Release != own {
		t.Errorf("hello said protocol %d and release %+v, want 2 and %+v", hello.Protocol, hello.Release, own)
	}
	st := waitState(t, c, "connected", func(s State) bool { return s.Connected })
	if !st.PeerSigned || st.PeerVersion != "2.3.0" || st.RefusedBy != "" || st.RefusedReason != "" {
		t.Errorf("state: %+v", st)
	}
}

// S-0237: a dashboard whose stamp is missing, as a protocol 1 dashboard's
// answer is, does not verify, or names another component, is closed with
// 4403 before flai proves itself; the state says flai refused it and why,
// and one error says so with the dashboard's version.
func TestAnUnsignedDashboardIsRefusedBeforeTheProof(t *testing.T) {
	for _, c := range []struct {
		name    string
		release func(d *fakeDashboard) *Release
		reason  string
		version string
	}{
		{"an answer without release, as protocol 1 sends", func(*fakeDashboard) *Release { return nil }, "missing stamp", "test"},
		{"signed by another key", func(d *fakeDashboard) *Release {
			other, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
			if err != nil {
				t.Fatal(err)
			}
			return (&fakeDashboard{t: d.t, signer: other}).sign("flaiover 2.3.0 1a2b3c4")
		}, "invalid signature", "2.3.0"},
		{"not base64", func(*fakeDashboard) *Release {
			return &Release{Statement: "flaiover 2.3.0 1a2b3c4", Signature: "not base64!"}
		}, "invalid signature", "2.3.0"},
		{"wrong component", func(d *fakeDashboard) *Release { return d.sign("flai 2.3.0 1a2b3c4") }, "wrong component flai", "2.3.0"},
	} {
		t.Run(c.name, func(t *testing.T) {
			d := newFakeDashboard(t, "s3cret")
			d.answer(c.release(d), "")
			var logs lockedBuffer
			cl := &Client{URL: d.srv.URL, Key: []byte("s3cret"), Project: project(t), Methods: testMethods(), Version: "test", ReleaseKeys: d.keys,
				Logger: slog.New(slog.NewJSONHandler(&logs, nil)), MinBackoff: 10 * time.Millisecond, MaxBackoff: 40 * time.Millisecond}
			runClient(t, cl)
			select {
			case msg := <-d.failed:
				if !strings.Contains(msg, "no proof") || !strings.Contains(msg, "4403") || !strings.Contains(msg, "unsigned peer: "+c.reason) {
					t.Errorf("the dashboard read %q, want a close with 4403 and the reason before any proof", msg)
				}
			case <-d.proven:
				t.Fatal("flai proved itself to an unsigned dashboard")
			case <-time.After(5 * time.Second):
				t.Fatal("nothing happened")
			}
			st := waitState(t, cl, "refused", func(s State) bool { return s.RefusedBy != "" })
			if st.Connected || st.RefusedBy != RefusedByFlai || st.RefusedReason != c.reason || st.PeerVersion != c.version || st.PeerSigned ||
				!strings.Contains(st.LastError, "flai refused the dashboard: unsigned peer: "+c.reason) || !strings.Contains(st.LastError, "trying again in 1m0s") {
				t.Errorf("state: %+v", st)
			}
			said := refusals(t, &logs)
			if len(said) != 1 || said[0]["level"] != "ERROR" || said[0]["refused_by"] != "flai" || said[0]["reason"] != c.reason ||
				said[0]["peer_version"] != c.version || said[0]["component"] != "channel" {
				t.Errorf("logged %v, want one error with the reason and the peer's version", said)
			}
			d.mu.Lock()
			defer d.mu.Unlock()
			if len(d.hellos) != 1 || d.hellos[0].Release != nil {
				t.Errorf("hellos %+v, want one, and no stamp from a development build", d.hellos)
			}
		})
	}
}

// S-0237: a dashboard that closes hello with 4403 refused this flai: flai
// waits the held back-off before it dials again, says so once however often
// it is refused, and clears the refusal once it is served.
func TestADashboardThatRefusesThisFlaiIsDialledAgainAfterTheHeldBackoff(t *testing.T) {
	d := newFakeDashboard(t, "s3cret")
	signed := d.sign("flaiover 2.3.0 1a2b3c4")
	d.answer(signed, "missing stamp")
	var logs lockedBuffer
	c := &Client{URL: d.srv.URL, Key: []byte("s3cret"), Project: project(t), Methods: testMethods(), Version: "test", ReleaseKeys: d.keys,
		Logger: slog.New(slog.NewJSONHandler(&logs, nil)), MinBackoff: 10 * time.Millisecond, MaxBackoff: 40 * time.Millisecond, HeldBackoff: 600 * time.Millisecond}
	runClient(t, c)
	st := waitState(t, c, "refused", func(s State) bool { return s.RefusedBy != "" })
	refused := time.Now()
	if st.Connected || st.RefusedBy != RefusedByDashboard || st.RefusedReason != "missing stamp" || st.PeerVersion != "" ||
		!strings.Contains(st.LastError, "the dashboard refused this flai: unsigned peer: missing stamp") || !strings.Contains(st.LastError, "trying again in 600ms") {
		t.Errorf("state: %+v", st)
	}
	hellos := func() int {
		d.mu.Lock()
		defer d.mu.Unlock()
		return len(d.hellos)
	}
	waitFor := func(n int) {
		deadline := time.Now().Add(5 * time.Second)
		for hellos() < n {
			if time.Now().After(deadline) {
				t.Fatalf("%d hellos, want %d", hellos(), n)
			}
			time.Sleep(5 * time.Millisecond)
		}
	}
	waitFor(2)
	// jitter spreads the wait over [HeldBackoff/2, HeldBackoff)
	if waited := time.Since(refused); waited < 250*time.Millisecond {
		t.Errorf("dialled again after %v, not after the held back-off", waited)
	}
	waitFor(3)
	if said := refusals(t, &logs); len(said) != 1 || said[0]["level"] != "ERROR" || said[0]["refused_by"] != "dashboard" || said[0]["reason"] != "missing stamp" {
		t.Errorf("logged %v, want one error for the same refusal met %d times", said, hellos())
	}

	d.answer(signed, "")
	waitProven(t, d)
	st = waitState(t, c, "served", func(s State) bool { return s.Connected })
	if st.RefusedBy != "" || st.RefusedReason != "" || !st.PeerSigned || st.PeerVersion != "2.3.0" {
		t.Errorf("the refusal outlived the connection: %+v", st)
	}
}

// A refusal is the last attempt's: one that ends otherwise, such as a
// dashboard that is gone, clears it from State, so status says why now.
func TestARefusalIsClearedByAnAttemptThatEndsOtherwise(t *testing.T) {
	d := newFakeDashboard(t, "s3cret")
	d.answer(d.sign("flaiover 2.3.0 1a2b3c4"), "missing stamp")
	c := &Client{URL: d.srv.URL, Key: []byte("s3cret"), Project: project(t), Methods: testMethods(), Version: "test", ReleaseKeys: d.keys,
		MinBackoff: 10 * time.Millisecond, MaxBackoff: 40 * time.Millisecond, HeldBackoff: 50 * time.Millisecond}
	runClient(t, c)
	waitState(t, c, "refused", func(s State) bool { return s.RefusedBy != "" })
	d.srv.Close()
	st := waitState(t, c, "cleared", func(s State) bool { return s.RefusedBy == "" && s.Attempts > 1 })
	if st.Connected || st.RefusedReason != "" || strings.Contains(st.LastError, UnsignedReason) {
		t.Errorf("state: %+v", st)
	}
}

// A close with 4403 that gives no fault still reads as the dashboard's refusal.
func TestDashboardRefusal(t *testing.T) {
	for reason, want := range map[string]string{"unsigned peer: invalid signature": "invalid signature", "unsigned peer": "unsigned", "": "unsigned"} {
		got := dashboardRefusal(websocket.CloseError{Code: CloseUnsigned, Reason: reason})
		if got == nil || got.By != RefusedByDashboard || got.Reason != want || !errors.Is(got, ErrUnsigned) {
			t.Errorf("close reason %q: %+v, want the dashboard's refusal with %q", reason, got, want)
		}
	}
	if got := dashboardRefusal(websocket.CloseError{Code: CloseHeld, Reason: "unsigned peer: missing stamp"}); got != nil {
		t.Errorf("a close with 4409 read as an unsigned refusal: %+v", got)
	}
}

// The gate: AllowUnsigned admits a dashboard without a stamp, never one whose
// stamp does not verify or names another component.
func TestAllowUnsignedAdmitsOnlyAMissingStamp(t *testing.T) {
	d := newFakeDashboard(t, "s3cret")
	d.answer(nil, "")
	c := &Client{URL: d.srv.URL, Key: []byte("s3cret"), Project: project(t), Methods: testMethods(), Version: "test", ReleaseKeys: d.keys, AllowUnsigned: allowed,
		MinBackoff: 10 * time.Millisecond, MaxBackoff: 40 * time.Millisecond}
	runClient(t, c)
	waitProven(t, d)
	if st := waitState(t, c, "connected", func(s State) bool { return s.Connected }); st.PeerSigned || st.RefusedBy != "" {
		t.Errorf("state: %+v", st)
	}
	for _, r := range []*Release{d.sign("flai 2.3.0 1a2b3c4"), {Statement: "flaiover 2.3.0 1a2b3c4", Signature: "not base64!"}} {
		if _, refused := c.checkDashboard(helloResult{Release: r}, true); refused == nil {
			t.Errorf("AllowUnsigned admitted %+v", r)
		}
	}
}

func allowed() bool { return true }

// provenPastRefusals waits for a proven connection, passing over the
// handshakes that flai refused before it.
func provenPastRefusals(t *testing.T, d *fakeDashboard) *websocket.Conn {
	t.Helper()
	for {
		select {
		case c := <-d.proven:
			return c
		case <-d.failed:
		case <-time.After(5 * time.Second):
			t.Fatal("no connection")
		}
	}
}

// signedFlai is a stamp for this flai's side: the fake dashboard does not
// check it, so any that is not the zero value stands for a release build.
var signedFlai = Release{Statement: "flai 1.43.0 a1825d8", Signature: "c2lnbmF0dXJl"}

// allowedUnsigned are the warn events that name the unsigned sides of a connection.
func allowedUnsigned(logs *lockedBuffer) []map[string]any {
	var out []map[string]any
	for _, line := range strings.Split(logs.String(), "\n") {
		var ev map[string]any
		if json.Unmarshal([]byte(line), &ev) == nil && ev["level"] == "WARN" && ev["unsigned"] != nil {
			out = append(out, ev)
		}
	}
	return out
}

// S-0239: with the allowance a connection records the sides it admitted
// without a stamp, flai then the dashboard, and warns once naming them; a
// signed pair records and says nothing.
func TestTheAllowanceRecordsAndWarnsOfTheUnsignedSides(t *testing.T) {
	for _, tc := range []struct {
		name      string
		own       Release
		dashboard bool // the dashboard answers with its signed stamp
		want      []string
	}{
		{"a dashboard without a stamp", signedFlai, false, []string{UnsignedDashboard}},
		{"a development build of flai", Release{}, true, []string{UnsignedFlai}},
		{"both", Release{}, false, []string{UnsignedFlai, UnsignedDashboard}},
		{"a signed pair", signedFlai, true, nil},
	} {
		t.Run(tc.name, func(t *testing.T) {
			d := newFakeDashboard(t, "s3cret")
			if !tc.dashboard {
				d.answer(nil, "")
			}
			var logs lockedBuffer
			c := &Client{URL: d.srv.URL, Key: []byte("s3cret"), Project: project(t), Methods: testMethods(), Version: "test", ReleaseKeys: d.keys,
				Release: tc.own, AllowUnsigned: allowed, Logger: slog.New(slog.NewJSONHandler(&logs, nil)),
				PingEvery: time.Hour, MinBackoff: 10 * time.Millisecond, MaxBackoff: 40 * time.Millisecond}
			runClient(t, c)
			waitProven(t, d)
			st := waitState(t, c, "connected", func(s State) bool { return s.Connected })
			if !slices.Equal(st.Unsigned, tc.want) || st.PeerSigned != tc.dashboard {
				t.Errorf("state %+v, want unsigned %v", st, tc.want)
			}
			said := allowedUnsigned(&logs)
			if tc.want == nil {
				if len(said) != 0 {
					t.Errorf("a signed pair logged %v", said)
				}
				return
			}
			if len(said) != 1 || said[0]["component"] != "channel" || said[0]["project"] != "harbour" || fmt.Sprint(said[0]["unsigned"]) != fmt.Sprint(tc.want) {
				t.Errorf("logged %v, want one warn naming %v", said, tc.want)
			}
		})
	}
}

// S-0239: without the allowance a dashboard without a stamp is closed with
// 4403 before the proof; the allowance is asked at each attempt, so once it
// is on the next attempt is admitted. Each connection proven warns once, and
// the unsigned side is cleared when it ends.
func TestTheAllowanceIsAskedAtEachAttempt(t *testing.T) {
	d := newFakeDashboard(t, "s3cret")
	d.answer(nil, "")
	var on atomic.Bool
	var logs lockedBuffer
	c := &Client{URL: d.srv.URL, Key: []byte("s3cret"), Project: project(t), Methods: testMethods(), Version: "test", ReleaseKeys: d.keys,
		Release: signedFlai, AllowUnsigned: on.Load, Logger: slog.New(slog.NewJSONHandler(&logs, nil)),
		PingEvery: time.Hour, MinBackoff: 10 * time.Millisecond, MaxBackoff: 40 * time.Millisecond, HeldBackoff: 40 * time.Millisecond}
	runClient(t, c)
	select {
	case msg := <-d.failed:
		if !strings.Contains(msg, "4403") || !strings.Contains(msg, "unsigned peer: missing stamp") {
			t.Errorf("the dashboard read %q, want a close with 4403 before any proof", msg)
		}
	case <-d.proven:
		t.Fatal("flai admitted a dashboard without a stamp without the allowance")
	case <-time.After(5 * time.Second):
		t.Fatal("nothing happened")
	}
	if st := waitState(t, c, "refused", func(s State) bool { return s.RefusedBy != "" }); st.RefusedBy != RefusedByFlai || st.Unsigned != nil {
		t.Errorf("state: %+v", st)
	}

	on.Store(true)
	first := provenPastRefusals(t, d)
	if st := waitState(t, c, "connected", func(s State) bool { return s.Connected }); !slices.Equal(st.Unsigned, []string{UnsignedDashboard}) || st.RefusedBy != "" {
		t.Errorf("state: %+v", st)
	}

	d.answer(nil, "missing stamp") // the dashboard now refuses this flai, so no connection follows
	_ = first.Close(websocket.StatusGoingAway, "restart")
	if st := waitState(t, c, "dropped", func(s State) bool { return !s.Connected && s.RefusedBy == RefusedByDashboard }); st.Unsigned != nil {
		t.Errorf("the unsigned side outlived the connection: %+v", st)
	}
	d.answer(nil, "")
	provenPastRefusals(t, d)
	waitState(t, c, "connected again", func(s State) bool { return s.Connected })
	if said := allowedUnsigned(&logs); len(said) != 2 {
		t.Errorf("logged %d warns over two connections, want one each: %v", len(said), said)
	}
}

// fakeGate is a BeforeDial that refuses with refuse while it is set, and
// records when it was asked.
type fakeGate struct {
	mu     sync.Mutex
	refuse error
	asked  []time.Time
}

func (g *fakeGate) before(context.Context) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.asked = append(g.asked, time.Now())
	return g.refuse
}

func (g *fakeGate) set(refuse error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.refuse = refuse
}

func (g *fakeGate) times() []time.Time {
	g.mu.Lock()
	defer g.mu.Unlock()
	return append([]time.Time(nil), g.asked...)
}

func waitUntil(t *testing.T, what string, ok func() bool) {
	t.Helper()
	deadline := time.Now().Add(3 * time.Second)
	for !ok() {
		if time.Now().After(deadline) {
			t.Fatalf("never %s", what)
		}
		time.Sleep(5 * time.Millisecond)
	}
}

func hellos(d *fakeDashboard) int {
	d.mu.Lock()
	defer d.mu.Unlock()
	return len(d.hellos)
}

// gateRefusals are the error events that give reason as theirs.
func gateRefusals(logs *lockedBuffer, reason string) []map[string]any {
	var out []map[string]any
	for _, line := range strings.Split(logs.String(), "\n") {
		var ev map[string]any
		if json.Unmarshal([]byte(line), &ev) == nil && ev["level"] == "ERROR" && ev["reason"] == reason {
			out = append(out, ev)
		}
	}
	return out
}

func gatedClient(d *fakeDashboard, g *fakeGate, logs *lockedBuffer, held time.Duration) *Client {
	return &Client{URL: d.srv.URL, Key: []byte("s3cret"), Project: Project{Key: "harbour", Name: "Harbour"}, Methods: testMethods(), Version: "test",
		ReleaseKeys: d.keys, BeforeDial: g.before, Logger: slog.New(slog.NewJSONHandler(logs, nil)),
		PingEvery: 50 * time.Millisecond, MinBackoff: 10 * time.Millisecond, MaxBackoff: 40 * time.Millisecond, HeldBackoff: held}
}

const unsignedImage = "the dashboard runs an unsigned image ghcr.io/bytepunx/flaiover:2.3.0 sha256:0b1d"

// S-0238: a gate that passes is asked once and the dashboard is dialled.
func TestAGateThatPassesDials(t *testing.T) {
	d := newFakeDashboard(t, "s3cret")
	g := &fakeGate{}
	var logs lockedBuffer
	c := gatedClient(d, g, &logs, 60*time.Millisecond)
	runClient(t, c)
	waitProven(t, d)
	st := waitState(t, c, "connected", func(s State) bool { return s.Connected })
	if n := len(g.times()); n != 1 || st.RefusedBy != "" || st.LastError != "" {
		t.Errorf("gate asked %d times, state %+v; want once and connected", n, st)
	}
}

// S-0238: a gate that refuses is not dialled; the state says flai refused
// the dashboard and why, as a refused stamp does; the gate is asked again
// after the held back-off; and one error says so however often it refuses.
func TestAGateThatRefusesIsNotDialled(t *testing.T) {
	d := newFakeDashboard(t, "s3cret")
	g := &fakeGate{refuse: errors.New(unsignedImage)}
	var logs lockedBuffer
	c := gatedClient(d, g, &logs, 60*time.Millisecond)
	runClient(t, c)
	waitUntil(t, "asked the gate four times", func() bool { return len(g.times()) >= 4 })

	if n := hellos(d); n != 0 {
		t.Errorf("dialled %d times past a gate that refused", n)
	}
	st := c.State()
	if st.Connected || st.RefusedBy != RefusedByFlai || st.RefusedReason != unsignedImage || st.PeerVersion != "" || st.PeerSigned || st.Attempts < 3 ||
		st.LastError != "flai refused the dashboard: "+unsignedImage+"; trying again in 60ms" {
		t.Errorf("state: %+v", st)
	}
	asked := g.times()
	for i := 1; i < len(asked); i++ {
		// jitter spreads the wait over [HeldBackoff/2, HeldBackoff)
		if gap := asked[i].Sub(asked[i-1]); gap < 25*time.Millisecond {
			t.Errorf("asked again after %v, not after the held back-off", gap)
		}
	}
	said := gateRefusals(&logs, unsignedImage)
	if len(said) != 1 || said[0]["component"] != "channel" || said[0]["url"] != d.srv.URL || said[0]["project"] != "harbour" || said[0]["retry_after"] != "60ms" {
		t.Errorf("logged %v, want one error for the same refusal met %d times", said, len(asked))
	}
}

// S-0238: the gate is asked on every attempt: a refusal is asked again and
// dialled once it passes, which clears the refusal; a reconnection after the
// dashboard ends the connection asks it again, and a refusal then is said
// again, since an attempt passed in between.
func TestTheGateIsAskedOnEveryAttempt(t *testing.T) {
	d := newFakeDashboard(t, "s3cret")
	g := &fakeGate{refuse: errors.New(unsignedImage)}
	var logs lockedBuffer
	c := gatedClient(d, g, &logs, 60*time.Millisecond)
	runClient(t, c)
	waitState(t, c, "refused", func(s State) bool { return s.RefusedBy != "" })

	g.set(nil)
	first := waitProven(t, d)
	st := waitState(t, c, "connected", func(s State) bool { return s.Connected })
	if st.RefusedBy != "" || st.RefusedReason != "" || st.LastError != "" || !st.PeerSigned {
		t.Errorf("the refusal outlived the connection: %+v", st)
	}
	before := len(g.times())

	g.set(errors.New(unsignedImage))
	_ = first.Close(websocket.StatusGoingAway, "restart")
	st = waitState(t, c, "refused again", func(s State) bool { return s.RefusedBy == RefusedByFlai })
	if len(g.times()) <= before || hellos(d) != 1 || st.RefusedReason != unsignedImage {
		t.Errorf("after the dashboard ended the connection: gate asked %d times (%d before), %d hellos, state %+v", len(g.times()), before, hellos(d), st)
	}
	if said := gateRefusals(&logs, unsignedImage); len(said) != 2 {
		t.Errorf("logged %d refusals, want the second said again after a connection", len(said))
	}

	g.set(nil)
	waitProven(t, d)
}

// S-0238: a client the gate refuses still ends as soon as its context does,
// however long the held back-off.
func TestARefusedGateEndsWithTheContext(t *testing.T) {
	d := newFakeDashboard(t, "s3cret")
	g := &fakeGate{refuse: errors.New(unsignedImage)}
	var logs lockedBuffer
	c := gatedClient(d, g, &logs, 0)
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() { c.Run(ctx); close(done) }()
	waitState(t, c, "refused", func(s State) bool { return strings.Contains(s.LastError, "trying again in 1m0s") })
	cancel()
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("Run outlived its context")
	}
}
