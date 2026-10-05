package mcpserver

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// lockedBuffer is a log destination a test reads while the server writes.
type lockedBuffer struct {
	mu sync.Mutex
	b  bytes.Buffer
}

func (l *lockedBuffer) Write(p []byte) (int, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.b.Write(p)
}

// answered are the "request answered" events logged, by method.
func (l *lockedBuffer) answered(t *testing.T) map[string]map[string]any {
	t.Helper()
	l.mu.Lock()
	defer l.mu.Unlock()
	out := map[string]map[string]any{}
	for _, line := range strings.Split(l.b.String(), "\n") {
		if !strings.Contains(line, `"request answered"`) {
			continue
		}
		var e map[string]any
		if err := json.Unmarshal([]byte(line), &e); err != nil {
			t.Fatal(err)
		}
		out[e["method"].(string)] = e
	}
	return out
}

func TestEachToolCallIsTimedBeneathTheTransport(t *testing.T) {
	root := t.TempDir()
	_ = os.WriteFile(filepath.Join(root, "system-flow.yaml"), []byte("version: 1\nname: t\nkey: t\nlayout:\n  design: design\n  docs: docs\n  wip: wip\n"), 0o644)
	for _, d := range []string{"wip/kanban/epics", "wip/kanban/stories", "wip/kanban/tasks", "wip/agents", "wip/archive"} {
		_ = os.MkdirAll(filepath.Join(root, d), 0o755)
	}
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	var logs lockedBuffer
	srv := New(Options{Repo: repo, Agent: "claude", Version: "test", MaxWait: 50 * time.Millisecond, Poll: 10 * time.Millisecond,
		Logger: slog.New(slog.NewJSONHandler(&logs, &slog.HandlerOptions{Level: slog.LevelDebug})), Slow: time.Nanosecond})
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	ct, st := mcp.NewInMemoryTransports()
	if _, err := srv.Connect(ctx, st, nil); err != nil {
		t.Fatal(err)
	}
	cs, err := mcp.NewClient(&mcp.Implementation{Name: "test-agent", Version: "0"}, nil).Connect(ctx, ct, nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = cs.Close() })
	for name, args := range map[string]any{"who_touches": map[string]any{"path": "flai"}, "wait_for_events": map[string]any{"timeout_seconds": 1}} {
		if res, err := cs.CallTool(ctx, &mcp.CallToolParams{Name: name, Arguments: args}); err != nil || res.IsError {
			t.Fatalf("%s: %v %+v", name, err, res)
		}
	}
	got := logs.answered(t)
	who := got["who_touches"]
	if who == nil || who["level"] != "INFO" || who["transport"] != "mcp" || who["duration_ms"] == nil || who["bytes"] == nil || who["err"] != nil {
		t.Errorf("who_touches: %v", who)
	}
	// a wait is long by design and is never logged as slow
	if w := got["wait_for_events"]; w == nil || w["level"] != "DEBUG" {
		t.Errorf("wait_for_events: %v", w)
	}
	// a request that is not a tool call is named by its method
	if got["tools/call"] != nil || len(got) < 3 {
		t.Errorf("methods timed: %v", got)
	}
}

// inbox lists the items once for the board and the changes (S-0156).
func TestInboxListsItemsOnce(t *testing.T) {
	root := t.TempDir()
	_ = os.WriteFile(filepath.Join(root, "system-flow.yaml"), []byte("version: 1\nname: t\nkey: t\nlayout:\n  design: design\n  docs: docs\n  wip: wip\n"), 0o644)
	for _, d := range []string{"wip/kanban/epics", "wip/kanban/stories", "wip/kanban/tasks", "wip/agents", "wip/archive"} {
		_ = os.MkdirAll(filepath.Join(root, d), 0o755)
	}
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	var logs lockedBuffer
	srv := New(Options{Repo: repo, Agent: "claude", Version: "test",
		Logger: slog.New(slog.NewJSONHandler(&logs, &slog.HandlerOptions{Level: slog.LevelDebug})), Slow: time.Nanosecond})
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	ct, st := mcp.NewInMemoryTransports()
	if _, err := srv.Connect(ctx, st, nil); err != nil {
		t.Fatal(err)
	}
	cs, err := mcp.NewClient(&mcp.Implementation{Name: "test-agent", Version: "0"}, nil).Connect(ctx, ct, nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = cs.Close() })
	if res, err := cs.CallTool(ctx, &mcp.CallToolParams{Name: "inbox", Arguments: map[string]any{}}); err != nil || res.IsError {
		t.Fatalf("inbox: %v %+v", err, res)
	}
	phases, _ := logs.answered(t)["inbox"]["phases"].(string)
	lists := 0
	for _, p := range strings.Fields(phases) {
		if name, took, _ := strings.Cut(p, "="); name == "repo.list" {
			lists = 1
			if _, n, ok := strings.Cut(took, "x"); ok {
				lists, _ = strconv.Atoi(n)
			}
		}
	}
	if lists != 1 {
		t.Errorf("inbox listed the items %d times: %s", lists, phases)
	}
}

// I-0059: a wait holds for what it asked, up to LongestWait, and for its
// default when it did not ask.
func TestHoldFor(t *testing.T) {
	for _, c := range []struct {
		seconds int
		def     time.Duration
		longest time.Duration
		want    time.Duration
	}{
		{0, time.Minute, LongestWait, time.Minute},
		{-5, 5 * time.Minute, LongestWait, 5 * time.Minute},
		{1800, time.Minute, LongestWait, 30 * time.Minute},
		{3600, time.Minute, LongestWait, LongestWait},
		{1 << 62, time.Minute, LongestWait, LongestWait}, // too many seconds for a Duration
		{0, 5 * time.Minute, 3 * time.Second, 3 * time.Second},
		{2, time.Minute, 50 * time.Millisecond, 50 * time.Millisecond},
	} {
		if got := holdFor(c.seconds, c.def, c.longest); got != c.want {
			t.Errorf("holdFor(%d, %s, %s) = %s, want %s", c.seconds, c.def, c.longest, got, c.want)
		}
	}
}

// deadlines stand in for time.After in a held wait: a deadline armed for d
// passes once the test has moved the clock on by d.
type deadlines struct {
	mu    sync.Mutex
	now   time.Duration
	set   []deadline
	armed chan time.Duration // every duration armed, in order
}

type deadline struct {
	at time.Duration
	c  chan time.Time
}

func newDeadlines() *deadlines {
	return &deadlines{armed: make(chan time.Duration, 16)}
}

func (d *deadlines) after(wait time.Duration) <-chan time.Time {
	d.mu.Lock()
	defer d.mu.Unlock()
	c := make(chan time.Time, 1)
	d.set = append(d.set, deadline{at: d.now + wait, c: c})
	d.armed <- wait
	return c
}

// advance moves the clock on and passes every deadline it reaches.
func (d *deadlines) advance(by time.Duration) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.now += by
	left := d.set[:0]
	for _, dl := range d.set {
		if dl.at <= d.now {
			dl.c <- time.Time{}
			continue
		}
		left = append(left, dl)
	}
	d.set = left
}

// checkHolds calls a wait three times, with nothing to answer: asked for
// 1800 s it is still held past 300 s and times out at 1800 s; asked for
// 3600 s it is cut to LongestWait; not asked, it holds for def.
func checkHolds(t *testing.T, d *deadlines, call func(args map[string]any) map[string]any, def time.Duration) {
	t.Helper()
	hold := func(args map[string]any, want time.Duration) <-chan map[string]any {
		t.Helper()
		done := make(chan map[string]any, 1)
		go func() { done <- call(args) }()
		select {
		case got := <-d.armed:
			if got != want {
				t.Fatalf("%v armed a deadline of %s, want %s", args, got, want)
			}
		case <-time.After(5 * time.Second):
			t.Fatalf("%v never armed a deadline", args)
		}
		return done
	}
	stillHeld := func(done <-chan map[string]any, args map[string]any) {
		t.Helper()
		select {
		case out := <-done:
			t.Fatalf("%v answered before its time: %v", args, out)
		case <-time.After(200 * time.Millisecond):
		}
	}
	timedOut := func(done <-chan map[string]any, args map[string]any) {
		t.Helper()
		select {
		case out := <-done:
			if out["timed_out"] != true {
				t.Errorf("%v: %v", args, out)
			}
		case <-time.After(5 * time.Second):
			t.Fatalf("%v never answered once its time passed", args)
		}
	}

	long := map[string]any{"timeout_seconds": 1800}
	done := hold(long, 30*time.Minute)
	d.advance(301 * time.Second)
	stillHeld(done, long)
	d.advance(30*time.Minute - 301*time.Second)
	timedOut(done, long)

	over := map[string]any{"timeout_seconds": 3600}
	done = hold(over, LongestWait)
	d.advance(LongestWait - time.Second)
	stillHeld(done, over)
	d.advance(time.Second)
	timedOut(done, over)

	unasked := map[string]any{}
	done = hold(unasked, def)
	d.advance(def - time.Second)
	stillHeld(done, unasked)
	d.advance(time.Second)
	timedOut(done, unasked)
}

// I-0059: a held wait sends a progress notification every heartbeat to a
// call that carries a progress token, so that a client that drops a silent
// call (Claude Code: 30 minutes over stdio, 5 over HTTP) keeps it; a call
// without one gets none.
func TestHeldWaitsSendProgress(t *testing.T) {
	var opt Options
	setupWith(t, func(o *Options) {
		o.Heartbeat, o.MaxWait = 20*time.Millisecond, 300*time.Millisecond
		opt = *o
	})
	opt.Agent = "idle" // with no story in progress, so wait_for_work holds
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	var mu sync.Mutex
	got := map[any]int{}
	ct, st := mcp.NewInMemoryTransports()
	if _, err := New(opt).Connect(ctx, st, nil); err != nil {
		t.Fatal(err)
	}
	cs, err := mcp.NewClient(&mcp.Implementation{Name: "test-agent", Version: "0"}, &mcp.ClientOptions{
		ProgressNotificationHandler: func(_ context.Context, req *mcp.ProgressNotificationClientRequest) {
			mu.Lock()
			defer mu.Unlock()
			got[req.Params.ProgressToken]++
		},
	}).Connect(ctx, ct, nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = cs.Close() })
	if _, err := cs.CallTool(ctx, &mcp.CallToolParams{Name: "inbox", Arguments: map[string]any{}}); err != nil {
		t.Fatal(err)
	}
	for _, tool := range []string{"wait_for_events", "wait_for_work"} {
		with := &mcp.CallToolParams{Name: tool, Arguments: map[string]any{"timeout_seconds": 1}}
		with.SetProgressToken(tool)
		if _, err := cs.CallTool(ctx, with); err != nil {
			t.Fatalf("%s: %v", tool, err)
		}
		if _, err := cs.CallTool(ctx, &mcp.CallToolParams{Name: tool, Arguments: map[string]any{"timeout_seconds": 1}}); err != nil {
			t.Fatalf("%s: %v", tool, err)
		}
	}
	mu.Lock()
	defer mu.Unlock()
	for _, tool := range []string{"wait_for_events", "wait_for_work"} {
		if got[tool] < 2 {
			t.Errorf("%s held for 300ms sent %d progress notifications, want one every 20ms", tool, got[tool])
		}
		delete(got, tool)
	}
	if len(got) > 0 {
		t.Errorf("a call without a progress token got progress: %v", got)
	}
}
