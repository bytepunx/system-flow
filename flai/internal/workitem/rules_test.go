package workitem

import (
	"slices"
	"strings"
	"testing"
	"time"

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
