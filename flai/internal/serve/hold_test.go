//go:build !windows

package serve

import (
	"context"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// readyTouching makes a story that touches paths and moves it to ready.
func (lab *agentLab) readyTouching(title string, touches ...string) string {
	lab.t.Helper()
	id := lab.backlog(title, nil, touches...)
	lab.toReady(id)
	return id
}

// S-0128, ADR-0046: the launcher skips a ready story whose claim overlaps an
// open story's, starts the next that is clear within the limit, says why the
// held one waits, and starts it first once the open story is gone.
func TestAHeldStoryIsSkippedAndStartedFirstOnceClear(t *testing.T) {
	lab := newAgentLab(t)
	ctx := context.Background()
	lab.limit(2)
	lab.hold()
	open := lab.readyTouching("Open", "flai/cmd")
	lab.move(open, workitem.InProgress)
	held := lab.readyTouching("Held", "flai/cmd/serve")
	clear := lab.readyTouching("Clear", "docs/clear")
	behind := lab.readyTouching("Behind", "docs/behind")
	lab.l.look(ctx, false)
	waitFor(t, "the clear story's agent runs", func() bool { return lab.run(clear).live() })
	why := "held (overlap): touches flai/cmd/serve, inside flai/cmd which " + open + " (in progress) touches; starts when " + open + " is accepted, cancelled, or sent back"
	if lab.run(held) != nil || lab.run(behind) != nil {
		t.Fatalf("started a held story, or one past the limit: %+v", lab.state().Stories)
	}
	if w := lab.state().Waiting; !strings.Contains(w, held+" "+why) || !strings.Contains(w, "limit leaves no room for "+behind) {
		t.Errorf("waiting: %q", w)
	}
	// agent.status: waiting, with the reason, though it never had an agent
	a := Activity(lab.root, lab.state())[held]
	if a.State != ActivityWaiting || a.Why != why || a.Hold == nil || a.Hold.Code != workitem.HoldOverlap || a.Run == nil || a.Run.Story != held || a.Run.PID != 0 {
		t.Errorf("activity: %+v", a)
	}
	if n := lab.said(held); n != 1 {
		t.Errorf("logged %d times why %s waits", n, held)
	}

	// the open story is sent back: the held story keeps its place and starts first
	st, _ := lab.repo.Get(open)
	if _, err := lab.repo.Transition(st, workitem.Cancelled, "alex", "not now", lab.now); err != nil {
		t.Fatal(err)
	}
	lab.l.look(ctx, false)
	waitFor(t, "the held story's agent runs", func() bool { return lab.run(held).live() })
	if lab.run(behind) != nil || !strings.Contains(lab.state().Waiting, "limit leaves no room for "+behind) {
		t.Errorf("the story behind it took its place: %+v", lab.state())
	}
	if _, ok := Activity(lab.root, lab.state())[held]; !ok || Activity(lab.root, lab.state())[held].Hold != nil {
		t.Errorf("still held once started: %+v", Activity(lab.root, lab.state())[held])
	}
	for _, id := range []string{held, clear} {
		lab.release(id)
	}
	waitFor(t, "both end", func() bool { return !lab.run(held).live() && !lab.run(clear).live() })
}

// One look never starts two stories that overlap: the first one started
// claims its paths before its agent pulls it. A story with no touches is
// held by any story open, and holds every other while it is open itself.
func TestOneLookStartsNoTwoThatOverlap(t *testing.T) {
	lab := newAgentLab(t)
	ctx := context.Background()
	lab.limit(3)
	lab.hold()
	first := lab.readyTouching("First", "flai")
	second := lab.readyTouching("Second", "flai/internal/serve")
	bare := lab.readyTouching("Bare", []string{}...)
	lab.l.look(ctx, false)
	waitFor(t, "the first runs", func() bool { return lab.run(first).live() })
	if lab.run(second) != nil || lab.run(bare) != nil {
		t.Fatalf("started beside an overlapping story: %+v", lab.state().Stories)
	}
	act := Activity(lab.root, lab.state())
	if a := act[second]; a.State != ActivityWaiting || !strings.Contains(a.Why, "inside flai which "+first+" (its agent started) touches") {
		t.Errorf("second: %+v", a)
	}
	if a := act[bare]; a.Hold == nil || a.Hold.Code != workitem.HoldNoTouches || !strings.Contains(a.Why, "declares no touches") {
		t.Errorf("bare: %+v", a)
	}
	lab.release(first)
	waitFor(t, "the first ends", func() bool { return !lab.run(first).live() })
}

// The operator's word starts a held story's agent, with a warning.
func TestAHeldStoryStartsOnTheOperatorsWord(t *testing.T) {
	lab := newAgentLab(t)
	ctx := context.Background()
	lab.hold()
	open := lab.readyTouching("Open", "docs")
	lab.move(open, workitem.InProgress)
	id := lab.readyTouching("Held", "docs/guide.md")
	lab.l.look(ctx, false)
	if lab.run(id) != nil {
		t.Fatal("the launcher started a held story")
	}
	run, err := Start(ctx, lab.o, Entry{Key: "t", Name: "t", Root: lab.root}, id)
	if err != nil || !run.live() {
		t.Fatalf("start: %+v %v", run, err)
	}
	lab.logs.mu.Lock()
	logged := lab.logs.buf.String()
	lab.logs.mu.Unlock()
	if !strings.Contains(logged, `level=WARN msg="agent started for a held story"`) || !strings.Contains(logged, "held (overlap): touches docs/guide.md, inside docs which "+open) {
		t.Errorf("no warning: %s", logged)
	}
	lab.release(id)
}
