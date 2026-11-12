package workitem

import (
	"slices"
	"strings"
	"testing"
	"time"
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
