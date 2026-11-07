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
	t.Parallel()
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
	t.Parallel()
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
	t.Parallel()
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
	t.Parallel()
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

// ADR-0244: a server with MaxWait HTTPLongestWait, as flai mcp http builds
// it, holds a wait_for_events or wait_for_work asked for 600 s, or for
// nothing, at most HTTPLongestWait and then answers timed_out, and both
// descriptions say so; without MaxWait, as flai mcp builds it over stdio, the
// same call is held the 600 s it asked for.
func TestWaitsHoldAtMostTheServersBound(t *testing.T) {
	t.Parallel()
	for _, c := range []struct {
		name    string
		maxWait time.Duration
		asked   time.Duration // what timeout_seconds 600 is held for
		events  time.Duration // what wait_for_events with none is held for
		work    time.Duration // what wait_for_work with none is held for
		says    string
	}{
		{"stdio", 0, 600 * time.Second, eventsWait, workWait, "at most 1800 (30 minutes)"},
		{"HTTP", HTTPLongestWait, HTTPLongestWait, HTTPLongestWait, HTTPLongestWait, "at most 45 (45 seconds), the most this server holds one call (ADR-0244)"},
	} {
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()
			d := newDeadlines()
			f := setupWith(t, func(o *Options) { o.MaxWait, o.After = c.maxWait, d.after })
			f.toReview(t)
			if _, failed := f.call(t, "inbox", map[string]any{}); failed != "" {
				t.Fatal(failed)
			}
			*f.clock = t0.Add(10 * time.Minute)
			tools, err := f.cs.ListTools(context.Background(), nil)
			if err != nil {
				t.Fatal(err)
			}
			described := map[string]bool{}
			for _, tool := range tools.Tools {
				if (tool.Name == "wait_for_events" || tool.Name == "wait_for_work") && strings.Contains(tool.Description, c.says) {
					described[tool.Name] = true
				}
			}
			for _, w := range []struct {
				tool string
				def  time.Duration
			}{{"wait_for_events", c.events}, {"wait_for_work", c.work}} {
				if !described[w.tool] {
					t.Errorf("%s's description states the bound in force, %q", w.tool, c.says)
				}
				for _, h := range []struct {
					args map[string]any
					want time.Duration
				}{{map[string]any{"timeout_seconds": 600}, c.asked}, {map[string]any{}, w.def}} {
					done := make(chan map[string]any, 1)
					go func() {
						out, failed := f.call(t, w.tool, h.args)
						if failed != "" {
							t.Errorf("%s: %s", w.tool, failed)
						}
						done <- out
					}()
					select {
					case got := <-d.armed:
						if got != h.want {
							t.Errorf("%s %v is held %s, want %s", w.tool, h.args, got, h.want)
						}
					case <-time.After(5 * time.Second):
						t.Fatalf("%s %v never armed a deadline", w.tool, h.args)
					}
					d.advance(h.want)
					select {
					case out := <-done:
						if out["timed_out"] != true {
							t.Errorf("%s %v answers timed_out once its bound passes: %v", w.tool, h.args, out)
						}
					case <-time.After(5 * time.Second):
						t.Fatalf("%s %v never answered once its bound passed", w.tool, h.args)
					}
				}
			}
		})
	}
}

// pausingClock is a server's Now that, once armed, stops the first look that
// reads it: it says so on fired and goes on when the test closes resume. A
// wait reads Now in its first look, after it has listed the items, so a test
// changes the project there, between the wait's first look and its hold.
type pausingClock struct {
	mu     sync.Mutex
	at     time.Time
	fired  chan struct{}
	resume chan struct{}
}

func (c *pausingClock) now() time.Time {
	c.mu.Lock()
	fired, resume := c.fired, c.resume
	c.fired, c.resume = nil, nil
	at := c.at
	c.mu.Unlock()
	if fired != nil {
		close(fired)
		<-resume
	}
	return at
}

// arm has the next read of the clock pause; the test closes resume once it
// has changed the project.
func (c *pausingClock) arm() (fired <-chan struct{}, resume chan<- struct{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	f, r := make(chan struct{}), make(chan struct{})
	c.fired, c.resume = f, r
	return f, r
}

// heldAcross calls a wait whose first look the clock pauses, makes change
// there, and returns what the wait answered.
func heldAcross(t *testing.T, clock *pausingClock, call func() toolAnswer, change func()) toolAnswer {
	t.Helper()
	fired, resume := clock.arm()
	done := make(chan toolAnswer, 1)
	go func() { done <- call() }()
	select {
	case <-fired:
	case a := <-done:
		t.Fatalf("the wait answered before its first look read the clock: %v %s", a.out, a.failed)
	case <-time.After(heldWait):
		t.Fatal("the wait's first look never read the clock")
	}
	change()
	close(resume)
	select {
	case a := <-done:
		return a
	case <-time.After(heldWait + 10*time.Second):
		t.Fatal("the wait never answered")
	}
	return toolAnswer{}
}

// I-0246: a change made after a wait's first look and before its hold
// fingerprints the project wakes the wait. Each wait once took its
// fingerprint after its first look, so a change in between was in the
// fingerprint and unseen by the look, and the wait held to its deadline; a
// loaded host, whose first look runs trial merges, widened the window.
func TestAChangeDuringAWaitsFirstLookWakesIt(t *testing.T) {
	t.Parallel()
	t.Run("wait_for_work", func(t *testing.T) {
		t.Parallel()
		clock := &pausingClock{at: t0.Add(time.Minute)}
		f := setupWith(t, func(o *Options) { o.Agent, o.Now, o.MaxWait = "idle", clock.now, heldWait })
		if quiet, failed := f.call(t, "wait_for_work", map[string]any{"timeout_seconds": 1}); failed != "" || quiet["timed_out"] != true {
			t.Fatalf("nothing ready: %v %s", quiet, failed)
		}
		var story *workitem.Item
		a := heldAcross(t, clock, func() toolAnswer {
			res, err := f.cs.CallTool(context.Background(), &mcp.CallToolParams{Name: "wait_for_work", Arguments: map[string]any{"timeout_seconds": int(heldWait / time.Second)}})
			return answerOf("wait_for_work", res, err)
		}, func() {
			story = readyStoryIn(t, f.repo, "Made ready meanwhile", t0.Add(2*time.Minute), "docs/users/guide.md")
		})
		if s, _ := a.out["story"].(map[string]any); a.failed != "" || a.out["timed_out"] == true || a.out["reason"] != "pull" || s["id"] != story.ID {
			t.Errorf("a story made ready during the first look: %v %s", a.out, a.failed)
		}
	})
	t.Run("a folder's wait_for_work", func(t *testing.T) {
		t.Parallel()
		root := t.TempDir()
		beta := makeProject(t, filepath.Join(root, "beta"), "beta")
		clock := &pausingClock{at: t0.Add(time.Minute)}
		f := folderSetupWith(t, root, func(o *Options) { o.Now, o.MaxWait = clock.now, heldWait })
		if quiet, failed := f.call(t, "wait_for_work", map[string]any{"timeout_seconds": 1}); failed != "" || quiet["timed_out"] != true {
			t.Fatalf("nothing ready: %v %s", quiet, failed)
		}
		var story *workitem.Item
		a := heldAcross(t, clock, func() toolAnswer {
			res, err := f.cs.CallTool(context.Background(), &mcp.CallToolParams{Name: "wait_for_work", Arguments: map[string]any{"timeout_seconds": int(heldWait / time.Second)}})
			return answerOf("wait_for_work", res, err)
		}, func() { story = readyStoryIn(t, beta, "Made ready meanwhile", t0.Add(2*time.Minute)) })
		if s, _ := a.out["story"].(map[string]any); a.failed != "" || a.out["timed_out"] == true || a.out["reason"] != "pull" || a.out["project"] != "beta" || s["id"] != story.ID {
			t.Errorf("a story made ready during the first look: %v %s", a.out, a.failed)
		}
	})
	t.Run("a folder's wait_for_events", func(t *testing.T) {
		t.Parallel()
		root := t.TempDir()
		alpha := makeProject(t, filepath.Join(root, "alpha"), "alpha")
		clock := &pausingClock{at: t0.Add(time.Minute)}
		f := folderSetupWith(t, root, func(o *Options) { o.Now, o.MaxWait = clock.now, heldWait })
		// the first look reports the epic makeProject wrote; the next holds
		for i := 0; ; i++ {
			out, failed := f.call(t, "wait_for_events", map[string]any{"timeout_seconds": 1})
			if failed != "" {
				t.Fatal(failed)
			}
			if out["timed_out"] == true {
				break
			}
			if i == 2 {
				t.Fatalf("wait_for_events never held with nothing new: %v", out)
			}
		}
		var story *workitem.Item
		a := heldAcross(t, clock, func() toolAnswer {
			res, err := f.cs.CallTool(context.Background(), &mcp.CallToolParams{Name: "wait_for_events", Arguments: map[string]any{"timeout_seconds": int(heldWait / time.Second)}})
			return answerOf("wait_for_events", res, err)
		}, func() { story = backlogStoryIn(t, alpha, "Written meanwhile", t0.Add(2*time.Minute)) })
		changed, _ := a.out["changed"].([]any)
		events, _ := a.out["events"].([]any)
		if a.failed != "" || a.out["timed_out"] == true || len(changed)+len(events) == 0 {
			t.Errorf("a story written during the first look (%s): %v %s", story.ID, a.out, a.failed)
		}
	})
}
