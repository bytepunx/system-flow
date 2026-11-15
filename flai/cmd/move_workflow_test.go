package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/guard"
	"github.com/bytepunx/system-flow/flai/internal/workflow"
)

// renamedWorkflow is a board unlike the baseline (ADR-0301 §6): review is
// renamed acceptance-review, keeping review as its alias, testing sits
// between in-progress and it, done is shipped, and cancelled is abandoned.
// Only the operator moves a story from testing to acceptance-review.
const renamedWorkflow = `---
title: Board
updated: 2026-11-15
status: active
states:
  - {name: backlog, kind: queue, initial: true, order: true}
  - {name: ready, kind: queue, order: true, limit: true}
  - {name: in-progress, kind: active, limit: true}
  - {name: testing, kind: active, limit: true}
  - {name: acceptance-review, kind: waiting, limit: true, aliases: [review]}
  - {name: shipped, kind: accepted}
  - {name: abandoned, kind: cancelled}
transitions:
  - {from: backlog, to: ready, effects: [commit]}
  - {from: ready, to: in-progress, effects: [pull]}
  - {from: in-progress, to: testing}
  - {from: testing, to: acceptance-review, types: [epic, story], by: [operator]}
  - {from: testing, to: acceptance-review, types: [task]}
  - {from: acceptance-review, to: shipped, types: [epic, story], by: [acceptance], effects: [accept]}
  - {from: acceptance-review, to: shipped, types: [task]}
  - {from: [in-progress, testing], to: shipped, types: [task]}
  - {from: acceptance-review, to: in-progress, effects: [send-back]}
  - {from: in-progress, to: ready}
  - {from: ready, to: backlog}
  - {from: [backlog, ready, in-progress, testing], to: abandoned, effects: [cancel]}
  - {from: acceptance-review, to: abandoned, by: [cascade], effects: [cancel]}
  - {from: abandoned, to: backlog, effects: [reopen]}
---
`

// useWorkflow makes every project the commands open in this test run the
// workflow file text in place of the baseline.
func useWorkflow(t *testing.T, text string) {
	t.Helper()
	w, p := workflow.Parse("design/workflows/board.md", text)
	if p != nil {
		t.Fatalf("the test's workflow does not parse: %s", p)
	}
	testWorkflow = &w
	t.Cleanup(func() { testWorkflow = nil })
}

// asOperator clears what makes this flai's caller a story's agent or a role.
func asOperator(t *testing.T) {
	t.Helper()
	for _, k := range []string{"FLAI_ROLE", "FLAI_STORY", guard.StartedByEnv} {
		t.Setenv(k, "")
	}
}

// renamedProject is a project on renamedWorkflow with an epic and a story
// whose criterion is written, both in backlog.
func renamedProject(t *testing.T) string {
	t.Helper()
	useConfig(t, filepath.Join(t.TempDir(), "cfg.json"))
	useWorkflow(t, renamedWorkflow)
	asOperator(t)
	root := tempProject(t)
	for _, args := range [][]string{{"config", "set", "author", "olive"}, {"epic", "new", "Epic"}, {"story", "new", "Moved", "--epic", "E-0001"}} {
		if _, errOut, code := runIn(t, root, args...); code != 0 {
			t.Fatal(errOut)
		}
	}
	file := filepath.Join(root, "wip/kanban/stories/S-0001-moved.md")
	s, err := os.ReadFile(file)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(file, []byte(strings.Replace(string(s), "## Acceptance criteria\n- [ ]\n", "## Acceptance criteria\n- [ ] ok\n", 1)), 0o644); err != nil {
		t.Fatal(err)
	}
	return root
}

// mustRun runs flai in root and fails the test on a non-zero exit.
func mustRun(t *testing.T, root string, args ...string) string {
	t.Helper()
	out, errOut, code := runIn(t, root, args...)
	if code != 0 {
		t.Fatalf("flai %s: exit %d\n%s%s", strings.Join(args, " "), code, out, errOut)
	}
	return out
}

// stateOf is the status in id's file.
func stateOf(t *testing.T, root, id string) string {
	t.Helper()
	dir := map[byte]string{'E': "epics", 'S': "stories", 'T': "tasks"}[id[0]]
	return statusOf(t, root, dir+"/"+id+"-*.md")
}

// S-0744 criteria 2 and 6: on a board with a state renamed and one added
// between in-progress and review, flai move takes a story through the added
// state, refuses a story's agent the move only the operator may make, naming
// the transition, who may make it, and the command that does, and lets the
// operator make it, naming the state by its alias.
func TestMoveRefusesATransitionTheMoverMayNotMake(t *testing.T) {
	root := renamedProject(t)
	for _, to := range []string{"ready", "in-progress", "testing"} {
		mustRun(t, root, "move", "S-0001", to)
	}
	t.Setenv("FLAI_STORY", "S-0001")
	_, errOut, code := runIn(t, root, "move", "S-0001", "review")
	for _, want := range []string{"rule: S-0001: only the operator may move", "from testing to acceptance-review", "flai move S-0001 acceptance-review"} {
		if code == 0 || !strings.Contains(errOut, want) {
			t.Errorf("an agent's move to review must be refused with %q: exit %d\n%s", want, code, errOut)
		}
	}
	if got := stateOf(t, root, "S-0001"); got != "testing" {
		t.Errorf("a refused move changes nothing: status %q, want testing", got)
	}
	t.Setenv("FLAI_STORY", "")
	if out := mustRun(t, root, "move", "S-0001", "review"); !strings.Contains(out, "S-0001 → acceptance-review") {
		t.Errorf("the operator's move names the state's current name: %s", out)
	}
}

// S-0744 criteria 2 and 6: an epic's acceptance stays the acceptance's on
// the added board too: an agent moving an epic in acceptance-review to
// shipped is not refused by who may move it, as flai move lets an epic be
// accepted today.
func TestMoveLetsAnEpicTakeTheAcceptanceTransition(t *testing.T) {
	root := renamedProject(t)
	for _, to := range []string{"ready", "in-progress", "testing", "acceptance-review"} {
		mustRun(t, root, "move", "E-0001", to)
	}
	t.Setenv("FLAI_STORY", "S-0002")
	mustRun(t, root, "move", "E-0001", "shipped")
	if got := stateOf(t, root, "E-0001"); got != "shipped" {
		t.Errorf("the epic's status is %q, want shipped", got)
	}
}

// S-0744 criterion 6: a cancellation on the added board moves the item and
// everything open under it to the cancelled kind's state, abandoned, and
// --dry-run names the waiting state an item in it was in.
func TestMoveCancelsAlongTheWorkflowsTransitions(t *testing.T) {
	root := renamedProject(t)
	mustRun(t, root, "task", "new", "--story", "S-0001", "Step")
	for _, to := range []string{"ready", "in-progress", "testing"} {
		mustRun(t, root, "move", "S-0001", to)
	}
	mustRun(t, root, "move", "T-0001", "ready")
	mustRun(t, root, "move", "T-0001", "in-progress")
	out := mustRun(t, root, "move", "E-0001", "abandoned", "--reason", "a different route", "--dry-run")
	for _, want := range []string{"would also cancel 2 items", "S-0001", "T-0001"} {
		if !strings.Contains(out, want) {
			t.Errorf("the dry run must name %q:\n%s", want, out)
		}
	}
	mustRun(t, root, "move", "E-0001", "abandoned", "--reason", "a different route", "--yes")
	for _, id := range []string{"E-0001", "S-0001", "T-0001"} {
		if got := stateOf(t, root, id); got != "abandoned" {
			t.Errorf("%s is %q, want abandoned", id, got)
		}
	}
}
