package serve

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/workflow"
)

// kindsFlow is the board the serve tests run against a workflow other than
// the baseline (S-0742), the one S-0743's tests use: in-progress is doing
// and review is awaiting, neither with an alias, and testing, of the active
// kind, sits between them. A pass that reads a state's name in place of its
// kind or a transition's effect starts, restarts, counts, or reasons about
// the wrong story under it.
const kindsFlow = `---
title: Board
updated: 2026-11-15
status: active
states:
  - {name: backlog, kind: queue, initial: true, order: true}
  - {name: ready, kind: queue, order: true, limit: true}
  - {name: doing, kind: active, limit: true}
  - {name: testing, kind: active, limit: true}
  - {name: awaiting, kind: waiting, limit: true}
  - {name: done, kind: accepted}
  - {name: cancelled, kind: cancelled}
transitions:
  - {from: backlog, to: ready, effects: [commit]}
  - {from: ready, to: doing, effects: [pull]}
  - {from: doing, to: testing}
  - {from: testing, to: awaiting}
  - {from: awaiting, to: done, types: [epic, story], by: [acceptance], effects: [accept]}
  - {from: awaiting, to: done, types: [task]}
  - {from: doing, to: done, types: [task]}
  - {from: awaiting, to: doing, effects: [send-back]}
  - {from: testing, to: doing}
  - {from: doing, to: ready}
  - {from: ready, to: backlog}
  - {from: [backlog, ready, doing, testing], to: cancelled, effects: [cancel]}
  - {from: awaiting, to: cancelled, by: [cascade], effects: [cancel]}
  - {from: cancelled, to: backlog, effects: [reopen]}
---
`

// kindsWorkflow is kindsFlow parsed, for Options.Workflow.
func kindsWorkflow(t *testing.T) *workflow.Workflow {
	t.Helper()
	w, p := workflow.Parse("design/workflows/board.md", kindsFlow)
	if p != nil {
		t.Fatalf("the test's workflow does not parse: %s", p)
	}
	return &w
}

// S-0742: the helpers read the pull, the kinds, and the active-work room
// from the workflow, so that in-progress renamed doing and testing added
// after it answer as the baseline's in-progress does.
func TestKindsReadTheWorkflowNotAStatesName(t *testing.T) {
	base := workflow.Baseline()
	renamed := *kindsWorkflow(t)
	cases := []struct {
		name     string
		wf       workflow.Workflow
		queue    string
		into     string
		actives  []string
		waits    string
		limits   map[string]int
		counts   map[string]int
		room     int
		limitSum int
	}{
		{"baseline", base, "ready", "in-progress", []string{"in-progress"}, "review",
			map[string]int{"ready": 5, "in-progress": 3, "review": 3}, map[string]int{"in-progress": 2, "review": 4}, 1, 3},
		{"renamed and added", renamed, "ready", "doing", []string{"doing", "testing"}, "awaiting",
			map[string]int{"ready": 5, "doing": 2, "testing": 1, "awaiting": 3}, map[string]int{"doing": 1, "testing": 1, "awaiting": 4}, 1, 3},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := pullQueue(c.wf); got != c.queue {
				t.Errorf("pullQueue = %q, want %q", got, c.queue)
			}
			if !isPullable(c.wf, c.queue) || isPullable(c.wf, "backlog") || isPullable(c.wf, c.into) {
				t.Errorf("isPullable does not take from %s alone", c.queue)
			}
			if !isPulledInto(c.wf, c.into) || isPulledInto(c.wf, c.queue) {
				t.Errorf("isPulledInto does not name %s", c.into)
			}
			for _, s := range c.actives {
				if !isActive(c.wf, s) || !isUnderway(c.wf, s) || isWaiting(c.wf, s) || isQueued(c.wf, s) || isClosed(c.wf, s) {
					t.Errorf("%s does not read as active work", s)
				}
			}
			if isActive(c.wf, c.waits) || !isWaiting(c.wf, c.waits) || !isUnderway(c.wf, c.waits) {
				t.Errorf("%s does not read as waiting", c.waits)
			}
			if !isQueued(c.wf, "backlog") || !isQueued(c.wf, c.queue) || isUnderway(c.wf, c.queue) {
				t.Error("the queue states do not read as queued")
			}
			if !isAccepted(c.wf, "done") || !isClosed(c.wf, "done") || isAccepted(c.wf, "cancelled") || !isClosed(c.wf, "cancelled") {
				t.Error("done and cancelled do not read as closed")
			}
			room, limit := activeRoom(c.wf, c.limits, c.counts)
			if room != c.room || limit != c.limitSum {
				t.Errorf("activeRoom = %d of %d, want %d of %d", room, limit, c.room, c.limitSum)
			}
		})
	}
	if room, limit := activeRoom(renamed, map[string]int{"ready": 5}, map[string]int{"doing": 4}); room != -1 || limit != 0 {
		t.Errorf("with no limit on an active state, activeRoom = %d of %d, want -1 of 0", room, limit)
	}
	if room, _ := activeRoom(renamed, map[string]int{"doing": 1}, map[string]int{"doing": 1, "testing": 2}); room != 0 {
		t.Errorf("a full active-work limit leaves room %d, want 0", room)
	}
	// The name in-progress no longer names a state of the renamed board, so
	// a pass that reads it finds nothing to work.
	if isActive(renamed, "in-progress") || isPulledInto(renamed, "in-progress") {
		t.Error("in-progress reads as a state of the renamed board")
	}
}

// S-0742: the launcher, openRepo, and Options.flow run the options'
// workflow, and the baseline when the options give none.
func TestTheLauncherOpensTheOptionsWorkflow(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "system-flow.yaml"), []byte("version: 1\nname: t\nkey: t\nlayout:\n  design: design\n  docs: docs\n  wip: wip\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	for _, o := range []Options{{}, {Workflow: kindsWorkflow(t)}} {
		into := "in-progress"
		if o.Workflow != nil {
			into = "doing"
		}
		l := &launcher{flow: o.Workflow}
		repo, err := l.open(root)
		if err != nil {
			t.Fatal(err)
		}
		other, err := openRepo(o, root)
		if err != nil {
			t.Fatal(err)
		}
		for name, wf := range map[string]workflow.Workflow{"the launcher's repository": repo.Flow(), "openRepo's repository": other.Flow(), "Options.flow": o.flow()} {
			if !isPulledInto(wf, into) {
				t.Errorf("%s does not pull into %s", name, into)
			}
		}
	}
	if _, err := (&launcher{}).open(t.TempDir()); err == nil {
		t.Error("the launcher opened a folder with no manifest")
	}
}
