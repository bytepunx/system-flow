package workitem

import (
	"strings"
	"testing"
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
