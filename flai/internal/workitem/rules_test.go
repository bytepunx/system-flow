package workitem

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/verifyrecord"
	"github.com/bytepunx/system-flow/flai/internal/workflow"
)

// storyInProgress makes S-0001 with criteria as its acceptance criteria and
// one task, and moves it to in-progress.
func storyInProgress(t *testing.T, r *Repo, criteria string) *Item {
	t.Helper()
	mustCreate(t, r, Epic, "E", "")
	s := mustCreate(t, r, Story, "S", "E-0001")
	s.Body = strings.Replace(s.Body, "## Acceptance criteria\n- [ ]\n", "## Acceptance criteria\n"+criteria, 1)
	if err := r.Save(s); err != nil {
		t.Fatal(err)
	}
	mustMove(t, r, s, Ready, "")
	mustMove(t, r, s, InProgress, "")
	mustCreate(t, r, Task, "T", "S-0001")
	s, err := r.Get("S-0001")
	if err != nil {
		t.Fatal(err)
	}
	return s
}

func criteriaWarning(warnings []string) string {
	for _, w := range warnings {
		if strings.Contains(w, "unticked acceptance criteria") {
			return w
		}
	}
	return ""
}

// A story goes to review with a box unticked, but is told how many and what
// to do (S-0282): only done refuses it.
func TestMoveToReviewWarnsOfUntickedCriteria(t *testing.T) {
	r := newProject(t)
	s := storyInProgress(t, r, "- [ ] one\n- [x] two\n- [ ] three\n")
	items, _ := r.List(false)
	warnings, err := r.Move(s, Review, MoveOptions{Now: t0, Items: items})
	if err != nil {
		t.Fatalf("an unticked criterion must not refuse review: %v", err)
	}
	want := "S-0001 goes to review with unticked acceptance criteria (2 of 3): tick each one verified with flai criteria tick S-0001 <n>, and say in its notes why any other is left unticked"
	if got := criteriaWarning(warnings); got != want {
		t.Errorf("warning:\n got %q\nwant %q", got, want)
	}
}

func TestMoveToReviewWithEveryCriterionTickedDoesNotWarn(t *testing.T) {
	r := newProject(t)
	s := storyInProgress(t, r, "- [x] one\n- [X] two\n")
	items, _ := r.List(false)
	warnings, err := r.Move(s, Review, MoveOptions{Now: t0, Items: items})
	if err != nil {
		t.Fatal(err)
	}
	if w := criteriaWarning(warnings); w != "" {
		t.Errorf("every box ticked, yet: %q", w)
	}
}

func TestMoveTaskToReviewDoesNotWarnOfCriteria(t *testing.T) {
	r := newProject(t)
	storyInProgress(t, r, "- [ ] one\n")
	task, err := r.Get("T-0001")
	if err != nil {
		t.Fatal(err)
	}
	task.Body += "\n## Acceptance criteria\n- [ ] one\n"
	mustMove(t, r, task, Ready, "")
	mustMove(t, r, task, InProgress, "")
	items, _ := r.List(false)
	warnings, err := r.Move(task, Review, MoveOptions{Now: t0, Items: items})
	if err != nil {
		t.Fatal(err)
	}
	if w := criteriaWarning(warnings); w != "" {
		t.Errorf("a task has no criteria to tick, yet: %q", w)
	}
}

// I-0347: S-0678 was blocked, went to review, and was accepted with the
// interval open, so every close-out after it reported item.blocked-closed. A
// move to done or cancelled closes the interval at the move and warns.
func TestAClosingMoveClosesTheOpenBlockedInterval(t *testing.T) {
	for _, to := range []string{Done, Cancelled} {
		t.Run(to, func(t *testing.T) {
			r := newProject(t)
			s := storyInProgress(t, r, "- [x] one\n")
			task, err := r.Get("T-0001")
			if err != nil {
				t.Fatal(err)
			}
			mustMove(t, r, task, Ready, "")
			mustMove(t, r, task, InProgress, "")
			mustMove(t, r, task, Done, "")
			if err := BlockItem(s, "pause for urgent fix", t0); err != nil {
				t.Fatal(err)
			}
			items, _ := r.List(false)
			if to == Done {
				if _, err := r.Move(s, Review, MoveOptions{Now: t0.Add(time.Minute), Items: items}); err != nil {
					t.Fatal(err)
				}
				if !s.IsBlocked() {
					t.Fatal("a move to review must leave the block open")
				}
			}
			at := t0.Add(2 * time.Minute)
			warnings, err := r.Move(s, to, MoveOptions{Now: at, Items: items, Reason: "dropped"})
			if err != nil {
				t.Fatal(err)
			}
			if s.IsBlocked() {
				t.Fatalf("%s left the blocked interval open: %+v", to, s.Blocked)
			}
			if got, want := s.Blocked[0].Until, at.UTC().Format(TimeFormat); got != want {
				t.Errorf("until = %q, want the move's time %q", got, want)
			}
			want := "S-0001 was blocked (pause for urgent fix); the move to " + to + " closed its blocked interval"
			if !slices.Contains(warnings, want) {
				t.Errorf("warnings %q lack %q", warnings, want)
			}
		})
	}
}

// A closing move of an item that is not blocked says nothing of blocks.
func TestAClosingMoveOfAnUnblockedItemDoesNotWarnOfABlock(t *testing.T) {
	r := newProject(t)
	s := storyInProgress(t, r, "- [ ] one\n")
	items, _ := r.List(false)
	warnings, err := r.Move(s, Cancelled, MoveOptions{Now: t0, Items: items, Reason: "dropped"})
	if err != nil {
		t.Fatal(err)
	}
	for _, w := range warnings {
		if strings.Contains(w, "blocked") {
			t.Errorf("unexpected warning %q", w)
		}
	}
}

// I-0347: an interval left open on a closed item, as S-0678's was, closes at
// the item's close, not at now, and never before the interval began.
func TestUnblockOnAClosedItemClosesAtItsClose(t *testing.T) {
	closed := t0.Add(time.Hour)
	for _, c := range []struct {
		name, from, want string
	}{
		{"before the close", t0.Format(TimeFormat), closed.Format(TimeFormat)},
		{"after the close", closed.Add(time.Minute).Format(TimeFormat), closed.Add(time.Minute).Format(TimeFormat)},
	} {
		t.Run(c.name, func(t *testing.T) {
			it := &Item{ID: "S-0001", Status: Done,
				Transitions: []Transition{{To: Review, At: t0.Format(TimeFormat)}, {To: Done, At: closed.Format(TimeFormat)}},
				Blocked:     []Block{{From: c.from, Reason: "pause for urgent fix"}}}
			later := t0.Add(48 * time.Hour)
			if err := UnblockItem(it, later); err != nil {
				t.Fatal(err)
			}
			if it.IsBlocked() || it.Blocked[0].Until != c.want {
				t.Errorf("until = %q, want %q", it.Blocked[0].Until, c.want)
			}
			if it.Updated != later.Format(TimeFormat) {
				t.Errorf("updated = %q, want the write's time", it.Updated)
			}
		})
	}
}

func TestUncheckedCount(t *testing.T) {
	for _, c := range []struct {
		body string
		want int
	}{
		{"", 0},
		{"## Acceptance criteria\n- [ ]\n", 0},
		{"## Acceptance criteria\n- [ ] one\n- [x] two\n  - [ ] nested\n", 2},
		{"- [ ] before\n## Acceptance criteria\n- [X] one\n\n## Notes\n- [ ] after\n", 0},
	} {
		if got := uncheckedCount(c.body); got != c.want {
			t.Errorf("uncheckedCount(%q) = %d, want %d", c.body, got, c.want)
		}
	}
}

// afterWarning is the warning a move to ready gives of an after: hold.
func afterWarning(warnings []string) string {
	for _, w := range warnings {
		if strings.Contains(w, "goes to ready held (after)") {
			return w
		}
	}
	return ""
}

// readyCandidate makes S-0002 under E-0001 with after and touches, and moves
// it to ready with the items as they are, returning the move's warnings.
func readyCandidate(t *testing.T, r *Repo, after, touches []string) []string {
	t.Helper()
	s := mustCreate(t, r, Story, "Waits", "E-0001")
	s.After, s.Touches = after, touches
	s.Body = strings.Replace(s.Body, "## Acceptance criteria\n- [ ]\n", "## Acceptance criteria\n- [ ] one\n", 1)
	if err := r.Save(s); err != nil {
		t.Fatal(err)
	}
	items, _ := r.List(false)
	warnings, err := r.Move(s, Ready, MoveOptions{Now: t0, Items: items})
	if err != nil {
		t.Fatalf("a story with after: must not be refused ready: %v", err)
	}
	return warnings
}

// I-0285: a story moved to ready by hand while a story it names in after: was
// not done sat there for days, and the move said nothing (ADR-0279).
func TestMoveToReadyWarnsOfAfterHold(t *testing.T) {
	r := newProject(t)
	s := storyInProgress(t, r, "- [ ] one\n")
	s.Touches = []string{"a.go"}
	if err := r.Save(s); err != nil {
		t.Fatal(err)
	}
	got := afterWarning(readyCandidate(t, r, []string{"S-0001"}, []string{"b.go"}))
	want := "S-0002 goes to ready held (after): waits for S-0001 (in progress); starts when S-0001 is done; flai serve starts no agent for it until then"
	if got != want {
		t.Errorf("warning:\n got %q\nwant %q", got, want)
	}
}

func TestMoveToReadyAfterADoneStoryDoesNotWarn(t *testing.T) {
	r := newProject(t)
	s := storyInProgress(t, r, "- [x] one\n")
	s.Status = Done
	if err := r.Save(s); err != nil {
		t.Fatal(err)
	}
	if w := afterWarning(readyCandidate(t, r, []string{"S-0001"}, []string{"b.go"})); w != "" {
		t.Errorf("S-0001 is done, yet: %q", w)
	}
}

// An overlap alone does not warn on the move to ready: a trial merge may
// clear it at the next look (ADR-0230), and only the move to in-progress says
// so.
func TestMoveToReadyOnOverlapAloneDoesNotWarn(t *testing.T) {
	r := newProject(t)
	s := storyInProgress(t, r, "- [ ] one\n")
	s.Touches = []string{"a.go"}
	if err := r.Save(s); err != nil {
		t.Fatal(err)
	}
	for _, w := range readyCandidate(t, r, nil, []string{"a.go"}) {
		if strings.Contains(w, "goes to ready") {
			t.Errorf("an overlap alone warned: %q", w)
		}
	}
}

// otherWorkflow is a board unlike the baseline: review is renamed
// acceptance-review, keeping review as its alias, and testing sits between
// in-progress and it (ADR-0301 §6).
const otherWorkflow = `---
title: Board
updated: 2026-11-14
status: active
states:
  - {name: backlog, kind: queue, initial: true, order: true}
  - {name: ready, kind: queue, order: true, limit: true}
  - {name: in-progress, kind: active, limit: true}
  - {name: testing, kind: active, limit: true}
  - {name: acceptance-review, kind: waiting, limit: true, aliases: [review]}
  - {name: done, kind: accepted}
  - {name: cancelled, kind: cancelled}
transitions:
  - {from: backlog, to: ready, effects: [commit]}
  - {from: ready, to: in-progress, effects: [pull]}
  - {from: in-progress, to: testing}
  - {from: testing, to: acceptance-review}
  - {from: acceptance-review, to: done, types: [epic, story], by: [acceptance], effects: [accept]}
  - {from: acceptance-review, to: done, types: [task]}
  - {from: in-progress, to: done, types: [task]}
  - {from: acceptance-review, to: in-progress, effects: [send-back]}
  - {from: in-progress, to: ready}
  - {from: ready, to: backlog}
  - {from: [backlog, ready, in-progress, testing], to: cancelled, effects: [cancel]}
  - {from: acceptance-review, to: cancelled, by: [cascade], effects: [cancel]}
  - {from: cancelled, to: backlog, effects: [reopen]}
---
`

// repoWithWorkflow is a new project that runs the workflow file text in place
// of the baseline, for tests of what a board with renamed and added states
// does (ADR-0301).
func repoWithWorkflow(t *testing.T, text string) *Repo {
	t.Helper()
	w, p := workflow.Parse("design/workflows/board.md", text)
	if p != nil {
		t.Fatalf("the test's workflow does not parse: %s", p)
	}
	r := newProject(t)
	r.Workflow = &w
	return r
}

// S-0737: a repository given another workflow moves a story through the
// states that workflow has and the baseline lacks, records each move under
// the state's current name when the move names an alias, and refuses a move
// that workflow has no transition for, naming the moves it has; the
// baseline refuses the added state as no state at all.
func TestMoveFollowsTheRepositorysWorkflow(t *testing.T) {
	r := repoWithWorkflow(t, otherWorkflow)
	s := storyInProgress(t, r, "- [x] one\n")
	mustMove(t, r, s, "testing", "")
	mustMove(t, r, s, Review, "")
	got, err := r.Get(s.ID)
	if err != nil {
		t.Fatal(err)
	}
	var tos []string
	for _, tr := range got.Transitions {
		tos = append(tos, tr.To)
	}
	if want := []string{Ready, InProgress, "testing", "acceptance-review"}; got.Status != "acceptance-review" || !slices.Equal(tos, want) {
		t.Errorf("status %q, transitions %v; want acceptance-review and %v", got.Status, tos, want)
	}

	other := mustCreate(t, r, Story, "S2", "E-0001")
	other.Body = strings.Replace(other.Body, "## Acceptance criteria\n- [ ]\n", "## Acceptance criteria\n- [ ] one\n", 1)
	if err := r.Save(other); err != nil {
		t.Fatal(err)
	}
	mustMove(t, r, other, Ready, "")
	mustMove(t, r, other, InProgress, "")
	items, _ := r.List(false)
	_, err = r.Move(other, Review, MoveOptions{Now: t0, Items: items})
	if want := "rule: S-0002 cannot go from in-progress to acceptance-review (allowed: testing, ready, cancelled)"; err == nil || err.Error() != want {
		t.Errorf("in-progress to review under the other workflow: %v; want %q", err, want)
	}

	base := newProject(t)
	b := storyInProgress(t, base, "- [x] one\n")
	items, _ = base.List(false)
	if _, err := base.Move(b, "testing", MoveOptions{Now: t0, Items: items}); err == nil || !strings.Contains(err.Error(), `"testing" is not a state; use one of backlog, ready, in-progress, review, done, cancelled`) {
		t.Errorf("the baseline moved a story to testing: %v", err)
	}
}

// S-0737: Follows reads the baseline, a cascade-only move included: review
// to cancelled is no move of its own, only a cancellation's.
func TestFollowsReadsTheBaseline(t *testing.T) {
	for _, tc := range []struct {
		typ, from, to string
		want          bool
	}{
		{Story, Backlog, Ready, true},
		{Story, InProgress, Done, false},
		{Task, InProgress, Done, true},
		{Story, Review, Done, true},
		{Task, Review, Cancelled, false},
		{Story, Done, Review, false},
		{Story, InProgress, "testing", false},
	} {
		if got := Follows(tc.typ, tc.from, tc.to); got != tc.want {
			t.Errorf("Follows(%s, %s, %s) = %v, want %v", tc.typ, tc.from, tc.to, got, tc.want)
		}
	}
}

// S-0738: the tests below hold Repo.Move to the rules its workflow lists,
// in their order, each with today's words.

// statusRunner answers every git command with out and err, as git status in
// a story's worktree would.
type statusRunner struct {
	out string
	err error
}

func (s statusRunner) Run(string, string, ...string) (string, error) { return s.out, s.err }

func (s statusRunner) RunInput(string, string, string, ...string) (string, error) {
	return s.out, s.err
}

func (s statusRunner) LookPath(name string) (string, error) { return name, nil }

// dirtyWorktree gives the story a worktree whose git status answers out and
// err.
func dirtyWorktree(t *testing.T, r *Repo, id, out string, err error) {
	t.Helper()
	if e := os.MkdirAll(r.WorktreePath(id), 0o755); e != nil {
		t.Fatal(e)
	}
	r.Git = statusRunner{out: out, err: err}
}

// ruleStory makes a story under E-0001 with criteria as its acceptance
// criteria and touches, in backlog.
func ruleStory(t *testing.T, r *Repo, criteria string, touches ...string) *Item {
	t.Helper()
	s := mustCreate(t, r, Story, "S", "E-0001")
	s.Body = strings.Replace(s.Body, "## Acceptance criteria\n- [ ]\n", "## Acceptance criteria\n"+criteria, 1)
	s.Touches = touches
	if err := r.Save(s); err != nil {
		t.Fatal(err)
	}
	return s
}

// moveWith is r.Move of it to to at t0, with the repository's items as they
// are and opt's other fields.
func moveWith(t *testing.T, r *Repo, it *Item, to string, opt MoveOptions) ([]string, error) {
	t.Helper()
	items, err := r.List(false)
	if err != nil {
		t.Fatal(err)
	}
	opt.Now, opt.Items = t0, items
	return r.Move(it, to, opt)
}

// refused fails the test unless err is the refusal want, word for word.
func refused(t *testing.T, err error, want string) {
	t.Helper()
	if err == nil || err.Error() != want {
		t.Errorf("refusal:\n got %v\nwant %s", err, want)
	}
}

// letThrough fails the test unless the move went through with the warnings
// want, in their order.
func letThrough(t *testing.T, warnings []string, err error, want ...string) {
	t.Helper()
	if err != nil {
		t.Fatalf("the move was refused: %v", err)
	}
	if !slices.Equal(warnings, want) {
		t.Errorf("warnings:\n got %q\nwant %q", warnings, want)
	}
}

// writeNarrative writes the story's narrative with body below its front
// matter.
func writeNarrative(t *testing.T, r *Repo, id, body string) {
	t.Helper()
	text := "---\nstream: " + id + "\ntitle: S\n---\n\n" + body
	if err := os.WriteFile(r.NarrativePath(id), []byte(text), 0o644); err != nil {
		t.Fatal(err)
	}
}

// headUnread is the baseline's warning, from verified, of a move to review
// whose story's head git did not read, saying why.
func headUnread(why string) string {
	return "moving S-0001 to review needs its last flai verify to have passed every suite at its branch's head, but its head was not read in .flai-cache/worktrees/S-0001: " + why + "; run flai verify S-0001, then move it again"
}

const reviewWithUnticked = "S-0001 goes to review with unticked acceptance criteria (1 of 2): tick each one verified with flai criteria tick S-0001 <n>, and say in its notes why any other is left unticked"

// Each rule the baseline refuses with refuses its move, in today's words.
func TestBaselineRulesRefuse(t *testing.T) {
	t.Run("criteria-present at backlog to ready", func(t *testing.T) {
		r := newProject(t)
		mustCreate(t, r, Epic, "E", "")
		s := mustCreate(t, r, Story, "S", "E-0001")
		_, err := moveWith(t, r, s, Ready, MoveOptions{})
		refused(t, err, "rule: a story needs an '## Acceptance criteria' section with at least one checkbox before it is ready")
	})
	t.Run("criteria-present at in-progress to ready", func(t *testing.T) {
		r := newProject(t)
		s := storyInProgress(t, r, "- [ ] one\n")
		s.Body = strings.Replace(s.Body, "- [ ] one\n", "", 1)
		_, err := moveWith(t, r, s, Ready, MoveOptions{})
		refused(t, err, "rule: a story needs an '## Acceptance criteria' section with at least one checkbox before it is ready")
	})
	t.Run("not-draft, and Finalize lets it through", func(t *testing.T) {
		r := newProject(t)
		mustCreate(t, r, Epic, "E", "")
		s := ruleStory(t, r, "- [ ] one\n")
		s.Draft = true
		_, err := moveWith(t, r, s, Ready, MoveOptions{})
		refused(t, err, "rule: S-0001 is a draft: finalize it first (flai edit S-0001 --no-draft, or flai move S-0001 ready --yes)")
		w, err := moveWith(t, r, s, Ready, MoveOptions{Finalize: true, By: "operator"})
		letThrough(t, w, err)
		if s.Draft || s.Status != Ready || s.Finalized == nil || s.Finalized.By != "operator" {
			t.Errorf("finalized move: draft %v, status %s, finalized %+v", s.Draft, s.Status, s.Finalized)
		}
	})
	t.Run("has-task", func(t *testing.T) {
		r := newProject(t)
		mustCreate(t, r, Epic, "E", "")
		s := ruleStory(t, r, "- [x] one\n")
		mustMove(t, r, s, Ready, "")
		mustMove(t, r, s, InProgress, "")
		_, err := moveWith(t, r, s, Review, MoveOptions{})
		refused(t, err, `rule: a story needs at least one task before it goes to review (flai task new --story S-0001 "...")`)
	})
	t.Run("no-open-question", func(t *testing.T) {
		r := newProject(t)
		s := storyInProgress(t, r, "- [x] one\n")
		writeNarrative(t, r, s.ID, "## Open questions\n- Which way?\n")
		_, err := moveWith(t, r, s, Review, MoveOptions{})
		refused(t, err, `rule: S-0001 has an open question in its narrative (`+r.NarrativePath(s.ID)+`): "Which way?"; it needs an answer before review (the operator answers it from the dashboard's inbox, or: flai stream answer S-0001 "Which way?" "<answer>")`)
	})
	t.Run("clean-worktree", func(t *testing.T) {
		r := newProject(t)
		s := storyInProgress(t, r, "- [x] one\n")
		dirtyWorktree(t, r, s.ID, " M flai/x.go\n?? docs/y.md\n", nil)
		_, err := moveWith(t, r, s, Review, MoveOptions{})
		want := "rule: the worktree .flai-cache/worktrees/S-0001 has uncommitted changes (flai/x.go, docs/y.md); commit them on story/S-0001, or discard them, before it goes to review"
		refused(t, err, want)
		if got := "rule: " + r.UncommittedRule(s.ID, []string{"flai/x.go", "docs/y.md"}, "it goes to review"); got != want {
			t.Errorf("the refusal is not UncommittedRule's text: %q", got)
		}
	})
	t.Run("criteria-ticked at done", func(t *testing.T) {
		r := newProject(t)
		s := storyInProgress(t, r, "- [ ] one\n- [x] two\n")
		mustMove(t, r, s, Review, "")
		tk, _ := r.Get("T-0001")
		mustMove(t, r, tk, Cancelled, "not needed")
		_, err := moveWith(t, r, s, Done, MoveOptions{})
		refused(t, err, "rule: S-0001 has unchecked acceptance criteria")
	})
	t.Run("children-closed, a story's open task", func(t *testing.T) {
		r := newProject(t)
		s := storyInProgress(t, r, "- [x] one\n")
		mustMove(t, r, s, Review, "")
		_, err := moveWith(t, r, s, Done, MoveOptions{})
		refused(t, err, "rule: S-0001 cannot be done while T-0001 is backlog")
	})
	t.Run("children-closed, an epic's open story", func(t *testing.T) {
		r := newProject(t)
		e := mustCreate(t, r, Epic, "E", "")
		s := ruleStory(t, r, "- [x] one\n")
		done := ruleStory(t, r, "- [x] one\n")
		mustMove(t, r, s, Ready, "")
		mustMove(t, r, s, InProgress, "")
		done.Status = Done
		if err := r.Save(done); err != nil {
			t.Fatal(err)
		}
		mustMove(t, r, e, Ready, "")
		mustMove(t, r, e, InProgress, "")
		mustMove(t, r, e, Review, "")
		_, err := moveWith(t, r, e, Done, MoveOptions{})
		refused(t, err, "rule: E-0001 cannot be done while S-0001 is in-progress")
		mustMove(t, r, s, Cancelled, "dropped")
		w, err := moveWith(t, r, e, Done, MoveOptions{})
		letThrough(t, w, err) // a done story and a cancelled one are closed
	})
	t.Run("reason-given", func(t *testing.T) {
		r := newProject(t)
		s := storyInProgress(t, r, "- [x] one\n")
		tk, _ := r.Get("T-0001")
		_, err := moveWith(t, r, tk, Cancelled, MoveOptions{})
		refused(t, err, "rule: moving T-0001 to cancelled needs --reason")
		_, err = moveWith(t, r, s, Cancelled, MoveOptions{Reason: "  "})
		refused(t, err, "rule: moving S-0001 to cancelled needs --reason")
		mustMove(t, r, s, Review, "")
		_, err = moveWith(t, r, s, InProgress, MoveOptions{})
		refused(t, err, "rule: moving S-0001 to in-progress needs --reason")
	})
	t.Run("parent-open", func(t *testing.T) {
		r := newProject(t)
		s := storyInProgress(t, r, "- [x] one\n")
		tk, _ := r.Get("T-0001")
		mustMove(t, r, tk, Cancelled, "gone")
		mustMove(t, r, s, Cancelled, "gone")
		_, err := moveWith(t, r, tk, Backlog, MoveOptions{})
		refused(t, err, "rule: T-0001 cannot go back to backlog while its parent S-0001 is cancelled; move S-0001 back first")
	})
}

// Each rule the baseline warns with lets its move through and says why.
func TestBaselineRulesWarn(t *testing.T) {
	t.Run("clean-worktree unread", func(t *testing.T) {
		r := newProject(t)
		s := storyInProgress(t, r, "- [x] one\n")
		dirtyWorktree(t, r, s.ID, "", errors.New("exit status 128"))
		w, err := moveWith(t, r, s, Review, MoveOptions{})
		letThrough(t, w, err, "S-0001's worktree was not checked for uncommitted changes: git status in .flai-cache/worktrees/S-0001: exit status 128",
			headUnread("exit status 128"))
	})
	t.Run("not-blocked at ready to in-progress", func(t *testing.T) {
		r := newProject(t)
		mustCreate(t, r, Epic, "E", "")
		s := ruleStory(t, r, "- [ ] one\n")
		mustMove(t, r, s, Ready, "")
		if err := BlockItem(s, "waiting", t0); err != nil {
			t.Fatal(err)
		}
		w, err := moveWith(t, r, s, InProgress, MoveOptions{})
		letThrough(t, w, err, "S-0001 has an open blocked interval")
	})
	t.Run("not-blocked at review to in-progress", func(t *testing.T) {
		r := newProject(t)
		s := storyInProgress(t, r, "- [x] one\n")
		mustMove(t, r, s, Review, "")
		if err := BlockItem(s, "waiting", t0); err != nil {
			t.Fatal(err)
		}
		w, err := moveWith(t, r, s, InProgress, MoveOptions{Reason: "redo"})
		letThrough(t, w, err, "S-0001 has an open blocked interval")
	})
	t.Run("not-blocked on a task", func(t *testing.T) {
		r := newProject(t)
		storyInProgress(t, r, "- [x] one\n")
		tk, _ := r.Get("T-0001")
		mustMove(t, r, tk, Ready, "")
		if err := BlockItem(tk, "waiting", t0); err != nil {
			t.Fatal(err)
		}
		w, err := moveWith(t, r, tk, InProgress, MoveOptions{})
		letThrough(t, w, err, "T-0001 has an open blocked interval")
	})
	t.Run("not-held", func(t *testing.T) {
		r := newProject(t)
		first := storyInProgress(t, r, "- [ ] one\n")
		first.Touches = []string{"a.go"}
		if err := r.Save(first); err != nil {
			t.Fatal(err)
		}
		s := ruleStory(t, r, "- [ ] one\n", "a.go")
		mustMove(t, r, s, Ready, "")
		w, err := moveWith(t, r, s, InProgress, MoveOptions{})
		letThrough(t, w, err, "S-0002 is held (overlap): touches a.go, which S-0001 (in progress) touches too; starts when S-0001 moves to review, is cancelled, or is sent back")
	})
	t.Run("after-done at in-progress to ready", func(t *testing.T) {
		r := newProject(t)
		first := storyInProgress(t, r, "- [ ] one\n")
		first.Touches = []string{"a.go"}
		if err := r.Save(first); err != nil {
			t.Fatal(err)
		}
		s := ruleStory(t, r, "- [ ] one\n", "b.go")
		mustMove(t, r, s, Ready, "")
		mustMove(t, r, s, InProgress, "")
		s.After = []string{"S-0001"}
		if err := r.Save(s); err != nil {
			t.Fatal(err)
		}
		w, err := moveWith(t, r, s, Ready, MoveOptions{})
		letThrough(t, w, err, "S-0002 goes to ready held (after): waits for S-0001 (in progress); starts when S-0001 is done; flai serve starts no agent for it until then")
	})
	t.Run("task-after", func(t *testing.T) {
		r := newProject(t)
		storyInProgress(t, r, "- [ ] one\n")
		tk := mustCreate(t, r, Task, "T2", "S-0001")
		tk.After = []string{"T-0001"}
		if err := r.Save(tk); err != nil {
			t.Fatal(err)
		}
		mustMove(t, r, tk, Ready, "")
		w, err := moveWith(t, r, tk, InProgress, MoveOptions{})
		letThrough(t, w, err, "T-0002 is waiting (after): waits for T-0001 (ready); ready to start when T-0001 is done or cancelled")
	})
}

// The rules run in the order the transition lists them: the first that
// refuses stops the move, and warnings come in the listed order, the WIP
// limit's last.
func TestMoveRunsTheRulesInTheirOrder(t *testing.T) {
	t.Run("refusals at in-progress to review", func(t *testing.T) {
		r := newProject(t)
		mustCreate(t, r, Epic, "E", "")
		s := ruleStory(t, r, "- [x] one\n")
		mustMove(t, r, s, Ready, "")
		mustMove(t, r, s, InProgress, "")
		writeNarrative(t, r, s.ID, "## Open questions\n- Which way?\n")
		dirtyWorktree(t, r, s.ID, " M a.go\n", nil)
		_, err := moveWith(t, r, s, Review, MoveOptions{})
		refused(t, err, `rule: a story needs at least one task before it goes to review (flai task new --story S-0001 "...")`)
		mustCreate(t, r, Task, "T", s.ID)
		_, err = moveWith(t, r, s, Review, MoveOptions{})
		if err == nil || !strings.HasPrefix(err.Error(), "rule: S-0001 has an open question in its narrative") {
			t.Errorf("no-open-question comes before clean-worktree: %v", err)
		}
		writeNarrative(t, r, s.ID, "## Open questions\n")
		_, err = moveWith(t, r, s, Review, MoveOptions{})
		refused(t, err, "rule: the worktree .flai-cache/worktrees/S-0001 has uncommitted changes (a.go); commit them on story/S-0001, or discard them, before it goes to review")
		r.Git = statusRunner{}
		w, err := moveWith(t, r, s, Review, MoveOptions{})
		letThrough(t, w, err, headUnread("git named no commit"))
	})
	t.Run("refusals at backlog to ready", func(t *testing.T) {
		r := newProject(t)
		mustCreate(t, r, Epic, "E", "")
		s := mustCreate(t, r, Story, "S", "E-0001")
		s.Draft = true
		_, err := moveWith(t, r, s, Ready, MoveOptions{})
		refused(t, err, "rule: a story needs an '## Acceptance criteria' section with at least one checkbox before it is ready")
	})
	t.Run("refusals at review to done", func(t *testing.T) {
		r := newProject(t)
		s := storyInProgress(t, r, "- [ ] one\n")
		mustMove(t, r, s, Review, "")
		_, err := moveWith(t, r, s, Done, MoveOptions{})
		refused(t, err, "rule: S-0001 cannot be done while T-0001 is backlog")
	})
	t.Run("warnings at ready to in-progress, the WIP limit last", func(t *testing.T) {
		r := newProject(t)
		first := storyInProgress(t, r, "- [ ] one\n")
		first.Touches = []string{"a.go"}
		if err := r.Save(first); err != nil {
			t.Fatal(err)
		}
		s := ruleStory(t, r, "- [ ] one\n", "a.go")
		mustMove(t, r, s, Ready, "")
		if err := BlockItem(s, "waiting", t0); err != nil {
			t.Fatal(err)
		}
		board, err := r.LoadBoard()
		if err != nil {
			t.Fatal(err)
		}
		board.WIPLimits[InProgress] = 1
		w, err := moveWith(t, r, s, InProgress, MoveOptions{Board: board})
		letThrough(t, w, err,
			"S-0002 has an open blocked interval",
			"S-0002 is held (overlap): touches a.go, which S-0001 (in progress) touches too; starts when S-0001 moves to review, is cancelled, or is sent back",
			"WIP limit for in-progress is 1, this makes 2")
	})
	t.Run("warnings at in-progress to review", func(t *testing.T) {
		r := newProject(t)
		s := storyInProgress(t, r, "- [ ] one\n- [x] two\n")
		dirtyWorktree(t, r, s.ID, "", errors.New("exit status 128"))
		w, err := moveWith(t, r, s, Review, MoveOptions{})
		letThrough(t, w, err,
			"S-0001's worktree was not checked for uncommitted changes: git status in .flai-cache/worktrees/S-0001: exit status 128",
			reviewWithUnticked,
			headUnread("exit status 128"))
	})
}

// ownWorkflow is a board of the operator's own: not-blocked refuses where
// the baseline warns, has-task and criteria-ticked warn where it refuses,
// ready to backlog lists rules for stories that a task's move skips, and the
// moves the baseline gives reason-given, parent-open, and children-closed
// list none.
const ownWorkflow = `---
title: Board
updated: 2026-11-14
status: active
states:
  - {name: backlog, kind: queue, initial: true, order: true}
  - {name: ready, kind: queue, order: true, limit: true}
  - {name: in-progress, kind: active, limit: true}
  - {name: review, kind: waiting, limit: true}
  - {name: done, kind: accepted}
  - {name: cancelled, kind: cancelled}
transitions:
  - {from: backlog, to: ready, effects: [commit], rules: [criteria-present]}
  - {from: ready, to: in-progress, effects: [pull], rules: [{rule: not-blocked, severity: refuse}]}
  - {from: in-progress, to: review, rules: [{rule: has-task, severity: warn}]}
  - {from: review, to: done, types: [epic, story], by: [acceptance], effects: [accept], rules: [{rule: criteria-ticked, severity: warn}]}
  - {from: review, to: done, types: [task]}
  - {from: in-progress, to: done, types: [task]}
  - {from: review, to: in-progress, effects: [send-back]}
  - {from: in-progress, to: ready}
  - {from: ready, to: backlog, rules: [has-task, criteria-present]}
  - {from: [backlog, ready, in-progress], to: cancelled, effects: [cancel]}
  - {from: review, to: cancelled, by: [cascade], effects: [cancel]}
  - {from: cancelled, to: backlog, effects: [reopen]}
---
`

// S-0738: a workflow of the operator's own sets each rule's severity, skips
// a rule for a type it does not apply to, and runs no rule a transition
// does not list.
func TestMoveRunsTheRepositorysOwnRules(t *testing.T) {
	t.Run("a warning rule set to refuse refuses with its warning", func(t *testing.T) {
		r := repoWithWorkflow(t, ownWorkflow)
		mustCreate(t, r, Epic, "E", "")
		s := ruleStory(t, r, "- [ ] one\n")
		mustMove(t, r, s, Ready, "")
		if err := BlockItem(s, "waiting", t0); err != nil {
			t.Fatal(err)
		}
		_, err := moveWith(t, r, s, InProgress, MoveOptions{})
		refused(t, err, "rule: S-0001 has an open blocked interval")
	})
	t.Run("a refusing rule set to warn lets the move through", func(t *testing.T) {
		r := repoWithWorkflow(t, ownWorkflow)
		mustCreate(t, r, Epic, "E", "")
		s := ruleStory(t, r, "- [ ] one\n- [x] two\n")
		mustMove(t, r, s, Ready, "")
		mustMove(t, r, s, InProgress, "")
		w, err := moveWith(t, r, s, Review, MoveOptions{})
		letThrough(t, w, err, `a story needs at least one task before it goes to review (flai task new --story S-0001 "...")`)
		w, err = moveWith(t, r, s, Done, MoveOptions{})
		letThrough(t, w, err, reviewWithUnticked)
	})
	t.Run("a rule is skipped for a type it does not apply to", func(t *testing.T) {
		r := repoWithWorkflow(t, ownWorkflow)
		mustCreate(t, r, Epic, "E", "")
		s := ruleStory(t, r, "- [ ] one\n")
		tk := mustCreate(t, r, Task, "T", s.ID)
		mustMove(t, r, tk, Ready, "")
		w, err := moveWith(t, r, tk, Backlog, MoveOptions{})
		letThrough(t, w, err) // has-task and criteria-present check stories alone
		lone := ruleStory(t, r, "- [ ] one\n")
		mustMove(t, r, lone, Ready, "")
		_, err = moveWith(t, r, lone, Backlog, MoveOptions{})
		refused(t, err, `rule: a story needs at least one task before it goes to review (flai task new --story S-0002 "...")`)
	})
	t.Run("a transition with no rules lets through what the baseline refuses", func(t *testing.T) {
		r := repoWithWorkflow(t, ownWorkflow)
		mustCreate(t, r, Epic, "E", "")
		s := ruleStory(t, r, "- [x] one\n")
		mustCreate(t, r, Task, "T", s.ID)
		mustMove(t, r, s, Ready, "")
		mustMove(t, r, s, InProgress, "")
		mustMove(t, r, s, Review, "")
		w, err := moveWith(t, r, s, InProgress, MoveOptions{})
		letThrough(t, w, err) // reason-given
		mustMove(t, r, s, Review, "")
		w, err = moveWith(t, r, s, Done, MoveOptions{})
		letThrough(t, w, err) // children-closed: T-0001 is in backlog
		other := ruleStory(t, r, "- [ ] one\n")
		tk := mustCreate(t, r, Task, "T2", other.ID)
		w, err = moveWith(t, r, tk, Cancelled, MoveOptions{})
		letThrough(t, w, err) // reason-given
		if err := r.Save(tk); err != nil {
			t.Fatal(err)
		}
		mustMove(t, r, other, Cancelled, "")
		w, err = moveWith(t, r, tk, Backlog, MoveOptions{})
		letThrough(t, w, err) // parent-open: S-0002 is cancelled
	})
}

// Every rule of the vocabulary has its check, and nothing else has one.
func TestEveryRuleHasACheck(t *testing.T) {
	for _, rule := range workflow.Vocabulary {
		if ruleChecks[rule.Name] == nil {
			t.Errorf("rule %s has no check in ruleChecks", rule.Name)
		}
	}
	for name := range ruleChecks {
		if _, ok := workflow.LookupRule(name); !ok {
			t.Errorf("ruleChecks has %s, which workflow.Vocabulary lacks", name)
		}
	}
	if len(ruleChecks) != len(workflow.Vocabulary) {
		t.Errorf("ruleChecks has %d checks for %d rules", len(ruleChecks), len(workflow.Vocabulary))
	}
}

// S-0807: the tests below hold the verified rule (ADR-0312 §2) to the verify
// record at the story's branch's head.

// headRunner answers git rev-parse with head and every other git command,
// git status among them, with nothing, as a clean worktree at head would.
type headRunner struct{ head string }

func (h headRunner) Run(_, _ string, args ...string) (string, error) {
	if len(args) > 0 && args[0] == "rev-parse" {
		return h.head, nil
	}
	return "", nil
}

func (h headRunner) RunInput(dir, name, _ string, args ...string) (string, error) {
	return h.Run(dir, name, args...)
}

func (headRunner) LookPath(name string) (string, error) { return name, nil }

const (
	verifiedHead = "1111111aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	olderHead    = "2222222bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"
)

// verifyStory is S-0001 in progress with its criteria ticked and a task, in
// a project whose move to review gives v, nil for the baseline's, with a
// worktree at verifiedHead.
func verifyStory(t *testing.T, v *workflow.Verify) (*Repo, *Item) {
	t.Helper()
	r := newProject(t)
	if v != nil {
		w := workflow.Baseline()
		for i, tr := range w.Transitions {
			if tr.To == Review && slices.Contains(tr.From, InProgress) {
				w.Transitions[i].Verify = v
			}
		}
		r.Workflow = &w
	}
	s := storyInProgress(t, r, "- [x] one\n")
	if err := os.MkdirAll(r.WorktreePath(s.ID), 0o755); err != nil {
		t.Fatal(err)
	}
	r.Git = headRunner{head: verifiedHead}
	return r, s
}

// suite is a suite's entry at commit, clean, whose tiers ended in the states
// given as name, state pairs; it passed when every tier passed.
func suite(name, commit string, tiers ...string) verifyrecord.Suite {
	s := verifyrecord.Suite{Name: name, Passed: true, Commit: commit, Clean: true}
	for i := 0; i+1 < len(tiers); i += 2 {
		state := verifyrecord.State(tiers[i+1])
		s.Tiers = append(s.Tiers, verifyrecord.Tier{Name: tiers[i], State: state})
		s.Passed = s.Passed && state == verifyrecord.Passed
	}
	return s
}

// keepRecord writes the record's verdict of S-0001 at commit: a clean run
// that passed when passed, holding suites.
func keepRecord(t *testing.T, r *Repo, commit string, passed bool, suites ...verifyrecord.Suite) {
	t.Helper()
	rec := verifyrecord.Record{Story: "S-0001", Commit: commit, Passed: passed, Clean: true, SuiteRuns: suites}
	data, err := json.Marshal(rec)
	if err != nil {
		t.Fatal(err)
	}
	path := verifyrecord.CommitPath(r.CacheDir(), "S-0001", commit)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatal(err)
	}
}

// notVerified is verified's text for S-0001's move to review needing what
// and failing for why.
func notVerified(what, why string) string {
	return "moving S-0001 to review needs its last flai verify to have passed " + what + " at its branch's head, but " + why + "; run flai verify S-0001, then move it again"
}

var (
	staticBehavior = &workflow.Verify{Suites: workflow.Suites{Names: []string{"static", "behavior"}}}
	everySuite     = &workflow.Verify{Suites: workflow.Suites{All: true}}
)

func TestVerifiedPassesWhenTheSuitesPassedAtTheHead(t *testing.T) {
	r, s := verifyStory(t, staticBehavior)
	keepRecord(t, r, verifiedHead, true,
		suite("static", verifiedHead, "gofmt", "passed", "vet", "passed"),
		suite("behavior", verifiedHead, "go-test", "passed"),
		suite("smoke", verifiedHead, "smoke", "failed"))
	w, err := moveWith(t, r, s, Review, MoveOptions{})
	letThrough(t, w, err) // smoke is not a suite the move needs
}

func TestVerifiedRefusesASuiteThatDidNotPass(t *testing.T) {
	r, s := verifyStory(t, staticBehavior)
	keepRecord(t, r, verifiedHead, false,
		suite("static", verifiedHead, "gofmt", "passed", "vet", "passed"),
		suite("behavior", verifiedHead, "go-test", "failed", "go-race", "not-reached"))
	_, err := moveWith(t, r, s, Review, MoveOptions{})
	refused(t, err, "rule: "+notVerified("suites static, behavior", "at 1111111 it did not pass behavior (0 of 2 tiers passed)"))
}

// A tier left to the merge queue runs in its batch, and the close-out passes
// without it (ADR-0174), so the move needs no more of it.
func TestVerifiedPassesASuiteLeftToTheMergeQueue(t *testing.T) {
	r, s := verifyStory(t, staticBehavior)
	keepRecord(t, r, verifiedHead, true,
		suite("static", verifiedHead, "gofmt", "passed"),
		suite("behavior", verifiedHead, "go-test", "passed", "integration", "queued"))
	w, err := moveWith(t, r, s, Review, MoveOptions{})
	letThrough(t, w, err)
}

func TestVerifiedWarnsAtWarn(t *testing.T) {
	r, s := verifyStory(t, nil) // the baseline: every suite, at warn
	keepRecord(t, r, verifiedHead, false, suite("static", verifiedHead, "gofmt", "failed"))
	w, err := moveWith(t, r, s, Review, MoveOptions{})
	letThrough(t, w, err, notVerified("every suite", "at 1111111 it did not pass static (0 of 1 tiers passed)"))
	if s.Status != Review {
		t.Errorf("the move at warn went through to %s, want review", s.Status)
	}
	other, o := verifyStory(t, nil)
	keepRecord(t, other, verifiedHead, true, suite("static", verifiedHead, "gofmt", "passed"))
	w, err = moveWith(t, other, o, Review, MoveOptions{})
	letThrough(t, w, err) // the warning is gone once every suite passed at the head
}

func TestVerifiedRefusesARecordAtAnOlderCommit(t *testing.T) {
	r, s := verifyStory(t, staticBehavior)
	keepRecord(t, r, olderHead, true,
		suite("static", olderHead, "gofmt", "passed"),
		suite("behavior", olderHead, "go-test", "passed"))
	_, err := moveWith(t, r, s, Review, MoveOptions{})
	refused(t, err, "rule: "+notVerified("suites static, behavior", "flai verify recorded no run of it at its head 1111111"))
}

func TestVerifiedRefusesASuiteNotVerifiedCleanAtTheHead(t *testing.T) {
	r, s := verifyStory(t, staticBehavior)
	dirty := suite("behavior", verifiedHead, "go-test", "passed")
	dirty.Clean = false
	keepRecord(t, r, verifiedHead, true, suite("static", verifiedHead, "gofmt", "passed"), dirty)
	_, err := moveWith(t, r, s, Review, MoveOptions{})
	refused(t, err, "rule: "+notVerified("suites static, behavior", "at 1111111 it did not pass behavior (its tiers passed, but not in a clean run at the head)"))
}

// A suite the record holds no entry of is one the story's paths select no
// tier of when the run passed, and is not verified when the run did not.
func TestVerifiedReadsAMissingSuiteByTheRun(t *testing.T) {
	r, s := verifyStory(t, staticBehavior)
	keepRecord(t, r, verifiedHead, true, suite("static", verifiedHead, "gofmt", "passed"))
	w, err := moveWith(t, r, s, Review, MoveOptions{})
	letThrough(t, w, err)

	r, s = verifyStory(t, staticBehavior)
	keepRecord(t, r, verifiedHead, false, suite("static", verifiedHead, "gofmt", "passed"))
	_, err = moveWith(t, r, s, Review, MoveOptions{})
	refused(t, err, "rule: "+notVerified("suites static, behavior", "at 1111111 it did not pass behavior (not in the record, whose run did not pass)"))
}

func TestVerifiedAllNeedsEverySuiteTheRecordHolds(t *testing.T) {
	r, s := verifyStory(t, everySuite)
	keepRecord(t, r, verifiedHead, false,
		suite("static", verifiedHead, "gofmt", "passed"),
		suite("smoke", verifiedHead, "smoke", "failed"))
	_, err := moveWith(t, r, s, Review, MoveOptions{})
	refused(t, err, "rule: "+notVerified("every suite", "at 1111111 it did not pass smoke (0 of 1 tiers passed)"))

	r, s = verifyStory(t, everySuite)
	keepRecord(t, r, verifiedHead, false)
	_, err = moveWith(t, r, s, Review, MoveOptions{})
	refused(t, err, "rule: "+notVerified("every suite", "its run at 1111111 did not pass"))

	r, s = verifyStory(t, everySuite)
	keepRecord(t, r, verifiedHead, true, suite("static", verifiedHead, "gofmt", "passed"), suite("smoke", verifiedHead, "smoke", "passed"))
	w, err := moveWith(t, r, s, Review, MoveOptions{})
	letThrough(t, w, err)
}

func TestVerifiedRefusesAStoryWithNoWorktree(t *testing.T) {
	r, s := verifyStory(t, everySuite)
	if err := os.Remove(r.WorktreePath(s.ID)); err != nil {
		t.Fatal(err)
	}
	_, err := moveWith(t, r, s, Review, MoveOptions{})
	refused(t, err, "rule: "+notVerified("every suite", "it has no worktree at .flai-cache/worktrees/S-0001"))
}

func TestVerifiedIsSkippedWithoutGit(t *testing.T) {
	r, s := verifyStory(t, everySuite)
	r.Git = nil
	w, err := moveWith(t, r, s, Review, MoveOptions{})
	letThrough(t, w, err)
}

// A record written before suites were recorded groups its tiers by the
// manifest's suites.
func TestVerifiedGroupsAnOldRecordByTheManifestsSuites(t *testing.T) {
	r, s := verifyStory(t, &workflow.Verify{Suites: workflow.Suites{Names: []string{"static"}}})
	r.Manifest.Tests = []manifest.TestTier{{Name: "gofmt", Suite: "static"}, {Name: "vet", Suite: "static"}}
	rec := verifyrecord.Record{Story: "S-0001", Commit: verifiedHead, Passed: false, Clean: true, Steps: []verifyrecord.Step{
		{Name: "gofmt", Tier: true, State: verifyrecord.Passed}, {Name: "vet", Tier: true, State: verifyrecord.Failed}}}
	data, err := json.Marshal(rec)
	if err != nil {
		t.Fatal(err)
	}
	path := verifyrecord.CommitPath(r.CacheDir(), "S-0001", verifiedHead)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatal(err)
	}
	_, err = moveWith(t, r, s, Review, MoveOptions{})
	refused(t, err, "rule: "+notVerified("suite static", "at 1111111 it did not pass static (1 of 2 tiers passed)"))
}

// The acceptance is not held by verified: S-0808 makes its gates read its
// verify (ADR-0312 §3).
func TestVerifiedDoesNotHoldTheAcceptance(t *testing.T) {
	r, s := verifyStory(t, nil)
	keepRecord(t, r, verifiedHead, true, suite("static", verifiedHead, "gofmt", "passed"))
	mustMove(t, r, s, Review, "")
	if err := os.RemoveAll(r.WorktreePath(s.ID)); err != nil {
		t.Fatal(err)
	}
	tk, _ := r.Get("T-0001")
	mustMove(t, r, tk, Cancelled, "not needed")
	w, err := moveWith(t, r, s, Done, MoveOptions{})
	letThrough(t, w, err)
}
