package workitem

import (
	"strings"
	"testing"
	"time"
)

// t1 is when the moves under test happen, after the setup's moves at t0.
var t1 = t0.Add(time.Hour)

// walkTo is the path a new item takes to each state, as mustMove makes it.
var walkTo = map[string][]string{
	Backlog:    nil,
	Ready:      {Ready},
	InProgress: {Ready, InProgress},
	Review:     {Ready, InProgress, Review},
	Done:       {Ready, InProgress, Review, Done},
	Cancelled:  {Cancelled},
}

// followProject is an epic E-0001 in epicState with one story under it per
// entry of states, S-0001 onwards, each with a task, moved there without the
// epic following.
func followProject(t *testing.T, epicState string, states ...string) *Repo {
	t.Helper()
	r := newProject(t)
	e := mustCreate(t, r, Epic, "E", "")
	for _, state := range states {
		s := mustCreate(t, r, Story, "in "+state, "E-0001")
		s.Body = strings.Replace(s.Body, "## Acceptance criteria\n- [ ]\n", "## Acceptance criteria\n- [x] works\n", 1)
		if err := r.Save(s); err != nil {
			t.Fatal(err)
		}
		task := mustCreate(t, r, Task, "under "+state, s.ID)
		if state == Done {
			for _, to := range []string{Ready, InProgress, Done} {
				mustMove(t, r, task, to, "")
			}
		}
		for _, to := range walkTo[state] {
			mustMove(t, r, s, to, "setup")
		}
	}
	for _, to := range walkTo[epicState] {
		mustMove(t, r, e, to, "setup")
	}
	return r
}

// move moves id through TransitionAll as alex at t1.
func move(t *testing.T, r *Repo, id, to, reason string) *MoveResult {
	t.Helper()
	it, err := r.Get(id)
	if err != nil {
		t.Fatal(err)
	}
	res, err := r.TransitionAll(it, to, "alex", reason, t1, false)
	if err != nil {
		t.Fatalf("move %s to %s: %v", id, to, err)
	}
	return res
}

func epicStatus(t *testing.T, r *Repo) string {
	t.Helper()
	e, err := r.Get("E-0001")
	if err != nil {
		t.Fatal(err)
	}
	return e.Status
}

func TestAnEpicGoesToReadyWithItsFirstStory(t *testing.T) {
	r := followProject(t, Backlog, Backlog, Backlog)
	res := move(t, r, "S-0001", Ready, "")
	want := Followed{ID: "E-0001", Type: Epic, Title: "E", From: Backlog, To: Ready, Story: "S-0001"}
	if res.Followed == nil || *res.Followed != want {
		t.Fatalf("followed %+v, want %+v", res.Followed, want)
	}
	e, _ := r.Get("E-0001")
	last := e.Transitions[len(e.Transitions)-1]
	if e.Status != Ready || last.To != Ready || last.By != "alex" || last.At != t1.Format(TimeFormat) {
		t.Errorf("the saved epic is %s with %+v, want ready by alex at %s", e.Status, last, t1.Format(TimeFormat))
	}
	if !strings.Contains(e.Body, "moved to ready: follows S-0001, which moved to ready") {
		t.Errorf("the epic's note does not name the story:\n%s", e.Body)
	}
	if res := move(t, r, "S-0002", Ready, ""); res.Followed != nil || epicStatus(t, r) != Ready {
		t.Errorf("the second story moved the epic: %+v, %s", res.Followed, epicStatus(t, r))
	}
}

func TestAnEpicStartsWithItsFirstStory(t *testing.T) {
	r := followProject(t, Ready, Ready, Backlog)
	if res := move(t, r, "S-0001", InProgress, ""); res.Followed == nil || res.Followed.From != Ready || res.Followed.To != InProgress {
		t.Errorf("followed %+v", res.Followed)
	}
	if got := epicStatus(t, r); got != InProgress {
		t.Errorf("epic is %s, want in-progress", got)
	}
}

func TestAnEpicInBacklogWalksToInProgress(t *testing.T) {
	r := followProject(t, Backlog, Ready)
	res := move(t, r, "S-0001", InProgress, "")
	if res.Followed == nil || res.Followed.From != Backlog || res.Followed.To != InProgress {
		t.Fatalf("followed %+v", res.Followed)
	}
	e, _ := r.Get("E-0001")
	n := len(e.Transitions)
	at := t1.Format(TimeFormat)
	if n != 2 || e.Transitions[0] != (Transition{To: Ready, At: at, By: "alex"}) || e.Transitions[1] != (Transition{To: InProgress, At: at, By: "alex"}) {
		t.Errorf("the epic did not walk one state at a time: %+v", e.Transitions)
	}
}

func TestAnEpicGoesToReviewWithItsLastOpenStory(t *testing.T) {
	r := followProject(t, InProgress, Review, InProgress, Done)
	if res := move(t, r, "S-0002", Review, ""); res.Followed == nil || res.Followed.To != Review {
		t.Errorf("followed %+v", res.Followed)
	}
	if got := epicStatus(t, r); got != Review {
		t.Errorf("epic is %s, want review", got)
	}
}

func TestAnEpicStaysWhileAnotherStoryIsOpen(t *testing.T) {
	r := followProject(t, InProgress, InProgress, InProgress, Backlog)
	if res := move(t, r, "S-0001", Review, ""); res.Followed != nil {
		t.Errorf("followed %+v", res.Followed)
	}
	if got := epicStatus(t, r); got != InProgress {
		t.Errorf("epic is %s, want in-progress", got)
	}
}

func TestAForwardMoveDoesNotPullAnEpicBack(t *testing.T) {
	// Moved ahead of its stories by hand.
	r := followProject(t, Review, Backlog, Backlog)
	if res := move(t, r, "S-0001", Ready, ""); res.Followed != nil || epicStatus(t, r) != Review {
		t.Errorf("followed %+v, epic %s", res.Followed, epicStatus(t, r))
	}
}

func TestAnEpicGoesBackFromReviewWithAStory(t *testing.T) {
	r := followProject(t, Review, Review, Review)
	res := move(t, r, "S-0001", InProgress, "found a gap")
	if res.Followed == nil || res.Followed.From != Review || res.Followed.To != InProgress {
		t.Fatalf("followed %+v", res.Followed)
	}
	e, _ := r.Get("E-0001")
	if e.Status != InProgress || !strings.Contains(e.Body, "moved to in-progress: follows S-0001, which moved to in-progress") {
		t.Errorf("epic %s:\n%s", e.Status, e.Body)
	}
}

func TestAnEpicGoesBackToReadyOnlyWhenNoStoryHoldsIt(t *testing.T) {
	r := followProject(t, InProgress, InProgress, InProgress)
	if res := move(t, r, "S-0001", Ready, ""); res.Followed != nil || epicStatus(t, r) != InProgress {
		t.Errorf("another story in progress holds the epic: %+v, %s", res.Followed, epicStatus(t, r))
	}
	if res := move(t, r, "S-0002", Ready, ""); res.Followed == nil || epicStatus(t, r) != Ready {
		t.Errorf("no story holds the epic: %+v, %s", res.Followed, epicStatus(t, r))
	}
}

func TestADoneStoryHoldsItsEpicInProgress(t *testing.T) {
	r := followProject(t, InProgress, InProgress, Done)
	if res := move(t, r, "S-0001", Ready, ""); res.Followed != nil || epicStatus(t, r) != InProgress {
		t.Errorf("followed %+v, epic %s", res.Followed, epicStatus(t, r))
	}
}

func TestAnEpicGoesBackToBacklogOnlyWhenNoStoryHoldsIt(t *testing.T) {
	r := followProject(t, Ready, Ready, Ready)
	if res := move(t, r, "S-0001", Backlog, ""); res.Followed != nil || epicStatus(t, r) != Ready {
		t.Errorf("another ready story holds the epic: %+v, %s", res.Followed, epicStatus(t, r))
	}
	if res := move(t, r, "S-0002", Backlog, ""); res.Followed == nil || epicStatus(t, r) != Backlog {
		t.Errorf("no story holds the epic: %+v, %s", res.Followed, epicStatus(t, r))
	}
}

func TestCancellingTheLastOpenStoryTakesTheEpicToReview(t *testing.T) {
	r := followProject(t, InProgress, Review, Done, InProgress)
	if res := move(t, r, "S-0003", Cancelled, "not needed"); res.Followed == nil || res.Followed.To != Review {
		t.Errorf("followed %+v", res.Followed)
	}
	if got := epicStatus(t, r); got != Review {
		t.Errorf("epic is %s, want review", got)
	}
}

func TestCancellingTheLastOpenStoryNeverTakesTheEpicToDone(t *testing.T) {
	r := followProject(t, InProgress, Done, Done, InProgress)
	move(t, r, "S-0003", Cancelled, "not needed")
	if got := epicStatus(t, r); got != Review {
		t.Errorf("epic is %s, want review: done is acceptance", got)
	}
}

func TestACancelledStoryDoesNotCount(t *testing.T) {
	r := followProject(t, InProgress, InProgress, Cancelled)
	if res := move(t, r, "S-0001", Review, ""); res.Followed == nil || epicStatus(t, r) != Review {
		t.Errorf("followed %+v, epic %s", res.Followed, epicStatus(t, r))
	}
}

func TestAStoryBackFromCancelledTakesTheEpicBack(t *testing.T) {
	r := followProject(t, Review, Review, Cancelled)
	res := move(t, r, "S-0002", Backlog, "")
	if res.Followed == nil || res.Followed.From != Review || res.Followed.To != InProgress {
		t.Errorf("followed %+v", res.Followed)
	}
}

func TestCancellingAnEpicDoesNotMakeItFollowItsStories(t *testing.T) {
	r := followProject(t, InProgress, InProgress, Ready)
	res := move(t, r, "E-0001", Cancelled, "a different route")
	if res.Followed != nil || len(res.Cancelled) != 4 {
		t.Errorf("followed %+v, cancelled %+v", res.Followed, res.Cancelled)
	}
	e, _ := r.Get("E-0001")
	if n := len(e.Transitions); e.Status != Cancelled || e.Transitions[n-2].To != InProgress {
		t.Errorf("the epic moved more than once: %+v", e.Transitions)
	}
}

func TestAStoryWithoutAnEpicMovesAlone(t *testing.T) {
	r := followProject(t, Backlog, Backlog)
	it, _ := r.Get("S-0001")
	if _, followed, err := r.Follow(nil, it, Backlog, "alex", t1, false); followed != nil || err != nil {
		t.Errorf("followed %+v, err %v", followed, err)
	}
}

// inMemory is an epic in epicState and stories in states, S-0001 onwards.
func inMemory(epicState string, states ...string) []*Item {
	items := []*Item{{ID: "E-0001", Type: Epic, Title: "E", Status: epicState}}
	for i, st := range states {
		items = append(items, &Item{ID: "S-000" + string(rune('1'+i)), Type: Story, Parent: "E-0001", Status: st})
	}
	return items
}

func TestADoneCancelledOrArchivedEpicIsLeftAlone(t *testing.T) {
	r := newProject(t)
	for _, tc := range []struct {
		name     string
		state    string
		archived bool
	}{{"done", Done, false}, {"cancelled", Cancelled, false}, {"archived", Review, true}} {
		items := inMemory(tc.state, Backlog, Backlog)
		items[0].Archived = tc.archived
		story := *items[1]
		story.Status = InProgress
		epic, followed, err := r.Follow(items, &story, Backlog, "alex", t1, true)
		if epic != nil || followed != nil || err != nil || items[0].Status != tc.state {
			t.Errorf("%s: epic %+v followed %+v err %v", tc.name, epic, followed, err)
		}
	}
}

func TestFollowReadsTheMovedStoryNotTheListsCopy(t *testing.T) {
	r := newProject(t)
	items := inMemory(Backlog, Backlog)
	story := *items[1]
	story.Status = Ready
	epic, followed, err := r.Follow(items, &story, Backlog, "alex", t1, false)
	if err != nil || followed == nil || epic.Status != Ready {
		t.Errorf("epic %+v followed %+v err %v", epic, followed, err)
	}
}

func TestAcceptingTheLastStoryAcceptsTheEpic(t *testing.T) {
	r := newProject(t)
	items := inMemory(InProgress, Done, Done)
	story := items[2]
	without, _, err := r.Follow(items, story, Review, "alex", t1, false)
	if err != nil || without == nil || without.Status != Review {
		t.Fatalf("without accept: epic %+v err %v", without, err)
	}
	items = inMemory(InProgress, Done, Done)
	epic, followed, err := r.Follow(items, items[2], Review, "alex", t1, true)
	if err != nil {
		t.Fatal(err)
	}
	if epic.Status != Done || followed.From != InProgress || followed.To != Done {
		t.Errorf("epic %s followed %+v", epic.Status, followed)
	}
	if n := len(epic.Transitions); n != 2 || epic.Transitions[0].To != Review || epic.Transitions[1].To != Done {
		t.Errorf("transitions %+v", epic.Transitions)
	}
}

func TestEpicFollowsItsStories(t *testing.T) {
	for _, tc := range []struct {
		states []string
		want   string
		ok     bool
	}{
		{nil, "", false},
		{[]string{Cancelled}, "", false},
		{[]string{Backlog, Cancelled}, Backlog, true},
		{[]string{Backlog, Ready}, Ready, true},
		{[]string{Ready, Done}, InProgress, true},
		{[]string{Review, Done, Cancelled}, Review, true},
		{[]string{Done, Done, Cancelled}, Done, true},
	} {
		got, ok := EpicFollows(inMemory(Backlog, tc.states...), "E-0001")
		if got != tc.want || ok != tc.ok {
			t.Errorf("%v: got %q %v, want %q %v", tc.states, got, ok, tc.want, tc.ok)
		}
	}
	archived := inMemory(Backlog, Backlog)
	archived = append(archived, &Item{ID: "S-0009", Type: Story, Parent: "E-0001", Status: Done, Archived: true})
	if got, _ := EpicFollows(archived, "E-0001"); got != InProgress {
		t.Errorf("an archived done story counts: got %s", got)
	}
}

func TestChangesNameTheStoryAnEpicFollowed(t *testing.T) {
	r := followProject(t, Backlog, Backlog)
	move(t, r, "S-0001", Ready, "")
	e, _ := r.Get("E-0001")
	board, _ := r.LoadBoard()
	items, _ := r.List(true)
	later := t1.Add(time.Minute)
	if _, err := r.Move(e, Backlog, MoveOptions{By: "alex", Now: later, Items: items, Board: board}); err != nil {
		t.Fatal(err)
	}
	if err := r.Save(e); err != nil {
		t.Fatal(err)
	}
	items, _ = r.List(true)
	follows := map[string]string{}
	for _, c := range Changes(items, t0, "", nil) {
		if c.ID == "E-0001" {
			follows[c.To+" "+c.At] = c.Follows
		}
	}
	if got := follows[Ready+" "+t1.Format(TimeFormat)]; got != "S-0001" {
		t.Errorf("the epic's move to ready follows %q, want S-0001 (%v)", got, follows)
	}
	if got := follows[Backlog+" "+later.Format(TimeFormat)]; got != "" {
		t.Errorf("a move by hand follows %q", got)
	}
}
