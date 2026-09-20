package workitem

import (
	"os"
	"strings"
	"testing"
)

// cascadeProject is an epic with a story in every state a story can be in,
// tasks under them, and a second epic that must not be touched.
func cascadeProject(t *testing.T) *Repo {
	t.Helper()
	r := newProject(t)
	mustCreate(t, r, Epic, "E", "")
	mustCreate(t, r, Epic, "Other", "")
	for _, state := range []string{Backlog, Ready, InProgress, Review, Done, Cancelled} {
		s := mustCreate(t, r, Story, "in "+state, "E-0001")
		s.Body = strings.Replace(s.Body, "## Acceptance criteria\n- [ ]\n", "## Acceptance criteria\n- [x] works\n", 1)
		if err := r.Save(s); err != nil {
			t.Fatal(err)
		}
		open := mustCreate(t, r, Task, "open under "+state, s.ID)
		closed := mustCreate(t, r, Task, "done under "+state, s.ID)
		for _, to := range []string{Ready, InProgress, Done} {
			mustMove(t, r, closed, to, "")
		}
		path := []string{}
		switch state {
		case Ready:
			path = []string{Ready}
		case InProgress:
			path = []string{Ready, InProgress}
		case Review, Done:
			path = []string{Ready, InProgress, Review}
		}
		for _, to := range path {
			mustMove(t, r, s, to, "")
		}
		switch state {
		case Done:
			for _, to := range []string{Ready, InProgress, Done} {
				mustMove(t, r, open, to, "")
			}
			mustMove(t, r, s, Done, "")
		case Cancelled:
			// As an older flai left it: the story cancelled, its task still open.
			mustMove(t, r, s, Cancelled, "its own reason")
		}
	}
	other := mustCreate(t, r, Story, "elsewhere", "E-0002")
	mustCreate(t, r, Task, "elsewhere", other.ID)
	return r
}

func statuses(t *testing.T, r *Repo) map[string]string {
	t.Helper()
	items, err := r.List(false)
	if err != nil {
		t.Fatal(err)
	}
	out := map[string]string{}
	for _, it := range items {
		out[it.ID] = it.Status
	}
	return out
}

func TestCancellingAnEpicCancelsEverythingOpenUnderIt(t *testing.T) {
	r := cascadeProject(t)
	before := statuses(t, r)
	e, _ := r.Get("E-0001")
	res, err := r.TransitionAll(e, Cancelled, "alex", "a different route", t0)
	if err != nil {
		t.Fatal(err)
	}
	after := statuses(t, r)
	// Stories S-0001..S-0006 are backlog, ready, in-progress, review, done, cancelled;
	// each has an open task (odd T) and a done task (even T), except the done story.
	want := map[string]string{
		"E-0001": Cancelled, "E-0002": Backlog,
		"S-0001": Cancelled, "S-0002": Cancelled, "S-0003": Cancelled, "S-0004": Cancelled,
		"S-0005": Done, "S-0006": Cancelled, "S-0007": Backlog,
		"T-0001": Cancelled, "T-0003": Cancelled, "T-0005": Cancelled, "T-0007": Cancelled,
		"T-0009": Done, "T-0011": Cancelled, "T-0013": Backlog,
		"T-0002": Done, "T-0004": Done, "T-0006": Done, "T-0008": Done, "T-0010": Done, "T-0012": Done,
	}
	for id, st := range want {
		if after[id] != st {
			t.Errorf("%s is %s, want %s (was %s)", id, after[id], st, before[id])
		}
	}
	var got []string
	for _, c := range res.Cancelled {
		got = append(got, c.ID+":"+c.From)
	}
	wantPlan := "S-0001:backlog T-0001:backlog S-0002:ready T-0003:backlog S-0003:in-progress T-0005:backlog S-0004:review T-0007:backlog T-0011:backlog"
	if strings.Join(got, " ") != wantPlan {
		t.Errorf("cancelled with it:\n got %s\nwant %s", strings.Join(got, " "), wantPlan)
	}

	task, _ := r.Get("T-0005")
	last := task.Transitions[len(task.Transitions)-1]
	if last.To != Cancelled || last.By != "alex" || last.At != t0.UTC().Format(TimeFormat) {
		t.Errorf("child transition: %+v", last)
	}
	if !strings.Contains(task.Body, "moved to cancelled: E-0001 cancelled: a different route") {
		t.Errorf("child note does not name the cause:\n%s", task.Body)
	}
	already, _ := r.Get("S-0006")
	if n := len(already.Transitions); already.Transitions[n-1].By != "test" || strings.Contains(already.Body, "E-0001 cancelled") {
		t.Errorf("a story that was already cancelled keeps its own history: %+v", already.Transitions)
	}
	board, _ := r.LoadBoard()
	if len(board.Order) != 0 {
		t.Errorf("cancelled stories stay in the pull order: %v", board.Order)
	}
}

func TestCancellingAStoryCancelsItsOpenTasks(t *testing.T) {
	r := cascadeProject(t)
	s, _ := r.Get("S-0003")
	res, err := r.TransitionAll(s, Cancelled, "alex", "not needed", t0)
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Cancelled) != 1 || res.Cancelled[0].ID != "T-0005" {
		t.Errorf("cancelled with it: %+v", res.Cancelled)
	}
	after := statuses(t, r)
	if after["T-0005"] != Cancelled || after["T-0006"] != Done || after["S-0002"] != Ready || after["E-0001"] != Backlog {
		t.Errorf("after: %v", after)
	}
}

func TestReviewIsCancelledOnlyThroughAParent(t *testing.T) {
	r := cascadeProject(t)
	s, _ := r.Get("S-0004")
	if _, err := r.TransitionAll(s, Cancelled, "alex", "directly", t0); err == nil || !strings.Contains(err.Error(), "cannot go from review to cancelled") {
		t.Errorf("a story in review cancelled directly: %v", err)
	}
}

func TestARefusedCancellationChangesNothing(t *testing.T) {
	r := cascadeProject(t)
	read := func() map[string]string {
		out := map[string]string{}
		items, _ := r.List(false)
		for _, it := range items {
			b, _ := os.ReadFile(it.Path)
			out[it.ID] = string(b)
		}
		return out
	}
	before := read()
	e, _ := r.Get("E-0001")
	if _, err := r.TransitionAll(e, Cancelled, "alex", "  ", t0); err == nil {
		t.Fatal("cancelling without a reason must be refused")
	}
	for id, body := range read() {
		if body != before[id] {
			t.Errorf("%s changed by a refused cancellation", id)
		}
	}
}

func TestNothingOpenMeansNothingExtra(t *testing.T) {
	r := newProject(t)
	mustCreate(t, r, Epic, "E", "")
	e, _ := r.Get("E-0001")
	res, err := r.TransitionAll(e, Cancelled, "alex", "empty", t0)
	if err != nil || len(res.Cancelled) != 0 || res.Cancelled == nil {
		t.Errorf("res %+v err %v", res, err)
	}
}
