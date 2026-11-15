package preview

import (
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/workflow"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// renamedBoard is a board unlike the baseline (ADR-0301 §6), as cmd's
// move_workflow_test.go has it: review is renamed acceptance-review, keeping
// review as its alias, testing sits between in-progress and it, done is
// shipped, and cancelled is abandoned. The acceptance lists the rules {rules}
// stands for.
const renamedBoard = `---
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
  - {from: acceptance-review, to: shipped, types: [epic, story], by: [acceptance], effects: [accept], rules: [{rules}]}
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

// boardWith is renamedBoard parsed, its acceptance listing rules.
func boardWith(t *testing.T, rules string) workflow.Workflow {
	t.Helper()
	w, p := workflow.Parse("design/workflows/board.md", strings.Replace(renamedBoard, "{rules}", rules, 1))
	if p != nil {
		t.Fatalf("the test's workflow does not parse: %s", p)
	}
	return w
}

// renamedRepo is a project, not a git repository, on w, with story S-0001
// whose one criterion is unticked, walked to status, and its task T-0001 in
// progress.
func renamedRepo(t *testing.T, w workflow.Workflow, status string, at time.Time) (*workitem.Repo, *workitem.Item, *workitem.Item) {
	t.Helper()
	root := t.TempDir()
	manifest := "version: 1\nname: t\nkey: t\nlayout:\n  design: design\n  docs: docs\n  wip: wip\n"
	if err := os.WriteFile(filepath.Join(root, "system-flow.yaml"), []byte(manifest), 0o644); err != nil {
		t.Fatal(err)
	}
	for _, d := range []string{"wip/kanban/epics", "wip/kanban/stories", "wip/kanban/tasks", "wip/agents"} {
		if err := os.MkdirAll(filepath.Join(root, d), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	repo.Workflow = &w
	s, err := repo.Create(workitem.NewOptions{Type: workitem.Story, Title: "Gate", Nature: "feature", Owner: "alex", Now: at})
	if err != nil {
		t.Fatal(err)
	}
	s.Body = strings.Replace(s.Body, "## Acceptance criteria\n", "## Acceptance criteria\n- [ ] gated\n", 1)
	if err := repo.Save(s); err != nil {
		t.Fatal(err)
	}
	task, err := repo.Create(workitem.NewOptions{Type: workitem.Task, Title: "Build", Parent: s.ID, Owner: "alex", Now: at})
	if err != nil {
		t.Fatal(err)
	}
	walk(t, repo, task, at, "ready", "in-progress")
	path := []string{"ready", "in-progress", "testing", "acceptance-review"}
	walk(t, repo, s, at, path[:slices.Index(path, status)+1]...)
	return repo, s, task
}

// walk moves it through each state in turn.
func walk(t *testing.T, repo *workitem.Repo, it *workitem.Item, at time.Time, states ...string) {
	t.Helper()
	for _, st := range states {
		if _, err := repo.Transition(it, st, "alex", "", at); err != nil {
			t.Fatal(err)
		}
	}
}

// tick ticks the story's criterion.
func tick(t *testing.T, repo *workitem.Repo, s *workitem.Item) {
	t.Helper()
	s.Body = strings.Replace(s.Body, "- [ ] gated", "- [x] gated", 1)
	if err := repo.Save(s); err != nil {
		t.Fatal(err)
	}
}

// S-0744 criteria 5 and 6: on a board with a state renamed and one added
// between in-progress and review, the acceptance preview of a story in the
// state the acceptance leaves from names the blockers the acceptance's rules
// give, in today's words with the state the acceptance reaches: an open task
// for children-closed, then, once the task is shipped, the unticked criterion
// for criteria-ticked; with both ticked and closed, nothing blocks.
func TestTheAcceptancePreviewListsTheBlockersOfTheAcceptancesRules(t *testing.T) {
	t.Parallel()
	at := time.Date(2026, 11, 14, 3, 0, 0, 0, time.UTC)
	log := slog.New(slog.NewTextHandler(io.Discard, nil))
	repo, s, task := renamedRepo(t, boardWith(t, "children-closed, criteria-ticked"), "acceptance-review", at)

	res, err := AcceptWith(noGit{}, repo, s, "alex", at, log, AcceptOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if want := "S-0001 cannot be shipped while T-0001 is in-progress"; !slices.Contains(res.Blockers, want) {
		t.Errorf("children-closed must block with %q: %q", want, res.Blockers)
	}

	walk(t, repo, task, at, "shipped")
	if res, err = AcceptWith(noGit{}, repo, s, "alex", at, log, AcceptOptions{}); err != nil {
		t.Fatal(err)
	}
	if want := "S-0001 has unchecked acceptance criteria"; !slices.Contains(res.Blockers, want) {
		t.Errorf("criteria-ticked must block with %q: %q", want, res.Blockers)
	}

	tick(t, repo, s)
	if res, err = AcceptWith(noGit{}, repo, s, "alex", at, log, AcceptOptions{}); err != nil {
		t.Fatal(err)
	}
	if len(res.Blockers) != 0 {
		t.Errorf("a story whose task is shipped and criterion ticked is accepted: %q", res.Blockers)
	}
}

// S-0744 criteria 5 and 6: the blockers come from the rules the acceptance
// lists: one listing neither has none for an open task or an unticked
// criterion, and one listing a rule as a warning is not blocked by it.
func TestTheAcceptancePreviewBlocksOnlyByTheRulesTheAcceptanceLists(t *testing.T) {
	t.Parallel()
	at := time.Date(2026, 11, 14, 3, 0, 0, 0, time.UTC)
	log := slog.New(slog.NewTextHandler(io.Discard, nil))
	for rules, want := range map[string][]string{
		"not-blocked": nil,
		"{rule: children-closed, severity: warn}, criteria-ticked": {"S-0001 has unchecked acceptance criteria"},
	} {
		repo, s, _ := renamedRepo(t, boardWith(t, rules), "acceptance-review", at)
		res, err := AcceptWith(noGit{}, repo, s, "alex", at, log, AcceptOptions{})
		if err != nil {
			t.Fatal(err)
		}
		if !slices.Equal(res.Blockers, want) {
			t.Errorf("with rules [%s] the blockers are %q, want %q", rules, res.Blockers, want)
		}
	}
}

// S-0744 criteria 3 and 6: a story in the added state, testing, is not in
// the state the acceptance leaves from, so its preview refuses it in today's
// words, naming testing and acceptance-review.
func TestTheAcceptancePreviewRefusesAStoryOutsideTheAcceptancesSource(t *testing.T) {
	t.Parallel()
	at := time.Date(2026, 11, 14, 3, 0, 0, 0, time.UTC)
	log := slog.New(slog.NewTextHandler(io.Discard, nil))
	repo, s, _ := renamedRepo(t, boardWith(t, "children-closed, criteria-ticked"), "testing", at)
	_, err := AcceptWith(noGit{}, repo, s, "alex", at, log, AcceptOptions{})
	if want := "S-0001 is testing; a story is accepted from acceptance-review"; err == nil || err.Error() != want {
		t.Errorf("the refusal is %v, want %q", err, want)
	}
}

// S-0744 criterion 6: the walk to the acceptance runs through the workflow's
// states, the added one among them; on the baseline it is today's walk.
func TestStepsToAcceptWalksTheWorkflowsStates(t *testing.T) {
	t.Parallel()
	w := boardWith(t, "children-closed")
	for _, c := range []struct {
		typ, from string
		want      []string
	}{
		{workitem.Story, "backlog", []string{"ready", "in-progress", "testing", "acceptance-review", "shipped"}},
		{workitem.Story, "in-progress", []string{"testing", "acceptance-review", "shipped"}},
		{workitem.Story, "review", []string{"shipped"}},
		{workitem.Epic, "testing", []string{"acceptance-review", "shipped"}},
		{workitem.Story, "abandoned", []string{"shipped"}},
	} {
		if got := StepsToAccept(w, c.typ, c.from); !slices.Equal(got, c.want) {
			t.Errorf("%s from %s: %q, want %q", c.typ, c.from, got, c.want)
		}
	}
	for from, want := range map[string][]string{
		"backlog":     {"ready", "in-progress", "review", "done"},
		"ready":       {"in-progress", "review", "done"},
		"in-progress": {"review", "done"},
		"review":      {"done"},
		"done":        {"done"},
		"cancelled":   {"done"},
	} {
		if got := StepsToDone(from); !slices.Equal(got, want) {
			t.Errorf("baseline from %s: %q, want %q", from, got, want)
		}
	}
}
