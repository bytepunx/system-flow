package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/guard"
	"github.com/bytepunx/system-flow/flai/internal/verifyrecord"
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

// verifyWorkflow is a board whose move to review needs the suite static,
// with the severity VERIFY_SEVERITY, passed at the story's head (ADR-0312 §2).
const verifyWorkflow = `---
title: Board
updated: 2026-11-18
status: active
states:
  - {name: backlog, kind: queue, initial: true, order: true}
  - {name: ready, kind: queue, order: true, limit: true}
  - {name: in-progress, kind: active, limit: true}
  - {name: review, kind: waiting, limit: true}
  - {name: done, kind: accepted}
  - {name: cancelled, kind: cancelled}
transitions:
  - {from: backlog, to: ready, effects: [commit]}
  - {from: ready, to: in-progress, effects: [pull]}
  - {from: in-progress, to: review, rules: [has-task], verify: {suites: [static], severity: VERIFY_SEVERITY}}
  - {from: review, to: done, types: [epic, story], by: [acceptance], effects: [accept], verify: all}
  - {from: review, to: done, types: [task]}
  - {from: in-progress, to: done, types: [task]}
  - {from: review, to: in-progress, effects: [send-back]}
  - {from: in-progress, to: ready}
  - {from: ready, to: backlog}
  - {from: [backlog, ready, in-progress], to: cancelled, effects: [cancel]}
  - {from: review, to: cancelled, by: [cascade], effects: [cancel]}
  - {from: cancelled, to: backlog, effects: [reopen]}
---
`

// verifyProject is a project on verifyWorkflow at severity with S-0001 in
// progress, a task under it, and a worktree whose one commit it answers.
func verifyProject(t *testing.T, severity string) (root, head string) {
	t.Helper()
	useConfig(t, filepath.Join(t.TempDir(), "cfg.json"))
	useWorkflow(t, strings.Replace(verifyWorkflow, "VERIFY_SEVERITY", severity, 1))
	asOperator(t)
	root = tempProject(t)
	for _, args := range [][]string{{"config", "set", "author", "olive"}, {"epic", "new", "Epic"}, {"story", "new", "Verified", "--epic", "E-0001"}, {"task", "new", "--story", "S-0001", "Step"}} {
		mustRun(t, root, args...)
	}
	for _, to := range []string{"ready", "in-progress"} {
		mustRun(t, root, "move", "S-0001", to)
	}
	wt := filepath.Join(root, ".flai-cache", "worktrees", "S-0001")
	if err := os.MkdirAll(wt, 0o755); err != nil {
		t.Fatal(err)
	}
	gitIn(t, wt, "init", "-q")
	if err := os.WriteFile(filepath.Join(wt, "a.go"), []byte("package a\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	gitIn(t, wt, "add", "a.go")
	gitIn(t, wt, "commit", "-q", "-m", "a")
	return root, gitIn(t, wt, "rev-parse", "HEAD")
}

// keepVerdict writes the verify record's verdict of S-0001 at commit: a
// clean run whose suite static passed when passed.
func keepVerdict(t *testing.T, root, commit string, passed bool) {
	t.Helper()
	state := verifyrecord.Passed
	if !passed {
		state = verifyrecord.Failed
	}
	rec := verifyrecord.Record{Story: "S-0001", Commit: commit, Passed: passed, Clean: true, SuiteRuns: []verifyrecord.Suite{
		{Name: "static", Passed: passed, Commit: commit, Clean: true, Tiers: []verifyrecord.Tier{{Name: "gofmt", State: state}}}}}
	data, err := json.Marshal(rec)
	if err != nil {
		t.Fatal(err)
	}
	path := verifyrecord.CommitPath(filepath.Join(root, ".flai-cache"), "S-0001", commit)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatal(err)
	}
}

// S-0807 criterion 1: through a workflow file whose move to review gives
// verify, flai move refuses the move until the record holds the suite passed
// at the story's head, naming the suite and flai verify, and lets it through
// once it does.
func TestMoveRefusesAMoveItsVerifyHasNotPassed(t *testing.T) {
	root, head := verifyProject(t, "refuse")
	keepVerdict(t, root, head, false)
	_, errOut, code := runIn(t, root, "move", "S-0001", "review")
	for _, want := range []string{"rule: moving S-0001 to review needs its last flai verify to have passed suite static at its branch's head", "it did not pass static (0 of 1 tiers passed)", "run flai verify S-0001"} {
		if code == 0 || !strings.Contains(errOut, want) {
			t.Errorf("the move to review must be refused with %q: exit %d\n%s", want, code, errOut)
		}
	}
	if got := stateOf(t, root, "S-0001"); got != "in-progress" {
		t.Errorf("a refused move changes nothing: status %q, want in-progress", got)
	}
	keepVerdict(t, root, head, true)
	if out := mustRun(t, root, "move", "S-0001", "review"); !strings.Contains(out, "S-0001 → review") {
		t.Errorf("the move once static passed at the head: %s", out)
	}
}

// S-0807 criterion 1: at warn, flai move lets the move through and warns,
// naming the suite and flai verify.
func TestMoveWarnsOfAVerifyNotPassedAtWarn(t *testing.T) {
	root, _ := verifyProject(t, "warn")
	out, errOut, code := runIn(t, root, "move", "S-0001", "review")
	if code != 0 || !strings.Contains(out, "S-0001 → review") {
		t.Fatalf("the move at warn must go through: exit %d\n%s%s", code, out, errOut)
	}
	for _, want := range []string{"workflow policy warning", "needs its last flai verify to have passed suite static", "flai verify recorded no run of it at its head", "run flai verify S-0001"} {
		if !strings.Contains(errOut, want) {
			t.Errorf("the warning must say %q:\n%s", want, errOut)
		}
	}
}
