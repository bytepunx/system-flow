package workitem

import (
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/workflow"
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

func TestAStoryGivenAnotherEpicMovesBoth(t *testing.T) {
	r := followProject(t, InProgress, InProgress, Review)
	other := mustCreate(t, r, Epic, "Other", "")
	all, err := r.List(true)
	if err != nil {
		t.Fatal(err)
	}
	story, err := r.Get("S-0001")
	if err != nil {
		t.Fatal(err)
	}
	story.Parent = other.ID
	epics, walks, err := r.FollowParent(all, story, "E-0001", "alex", t1)
	if err != nil {
		t.Fatal(err)
	}
	if len(walks) != 2 {
		t.Fatalf("walks %+v", walks)
	}
	// E-0001 keeps only S-0002, in review: it follows it there
	if w := walks[0]; w.ID != "E-0001" || w.From != InProgress || w.To != Review || w.Story != "S-0001" || epics[0].Status != Review {
		t.Errorf("the epic it left: %+v", w)
	}
	// E-0002 gains S-0001, in progress, from backlog
	if w := walks[1]; w.ID != other.ID || w.From != Backlog || w.To != InProgress || epics[1].Status != InProgress {
		t.Errorf("the epic it joined: %+v", w)
	}
	if !strings.Contains(epics[1].Body, "follows S-0001, which joined it from E-0001") {
		t.Errorf("the note names the story:\n%s", epics[1].Body)
	}
}

func TestAStoryInBacklogGivenAnotherEpicMovesNeither(t *testing.T) {
	r := followProject(t, Ready, Backlog, Ready)
	other := mustCreate(t, r, Epic, "Other", "")
	all, _ := r.List(true)
	story, _ := r.Get("S-0001")
	story.Parent = other.ID
	if _, walks, err := r.FollowParent(all, story, "E-0001", "alex", t1); err != nil || len(walks) != 0 {
		t.Errorf("walks %+v err %v", walks, err)
	}
}

// followKindFlow is one workflow the epic-following cases run under: how to
// open a repository that runs it, and its name for each baseline state.
type followKindFlow struct {
	name  string
	open  func(*testing.T) *Repo
	state func(string) string
}

// followKindFlows are the baseline and the fixture, with in-progress renamed
// doing and testing added between it and review, for the cases that must
// give the same result under each wherever their kinds agree (S-0739).
func followKindFlows() []followKindFlow {
	return []followKindFlow{
		{"baseline", newProject, func(s string) string { return s }},
		{"fixture", func(t *testing.T) *Repo { return withTestFlow(t, newProject(t)) }, inTestFlow},
	}
}

// followKindWalk is the path a new item takes to state, a baseline name or
// fixtureTesting, in f's names: walkTo's, with testing reached through doing.
func followKindWalk(f followKindFlow, state string) []string {
	if state == fixtureTesting {
		return []string{f.state(Ready), fixtureDoing, fixtureTesting}
	}
	var out []string
	for _, s := range walkTo[state] {
		out = append(out, f.state(s))
	}
	return out
}

// followKindProject is followProject under f: an epic E-0001 in epicState
// with one story under it per entry of states, S-0001 onwards, each with a
// task, moved there without the epic following. States are baseline names,
// or fixtureTesting under the fixture.
func followKindProject(t *testing.T, f followKindFlow, epicState string, states ...string) *Repo {
	t.Helper()
	r := f.open(t)
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
				mustMove(t, r, task, f.state(to), "")
			}
		}
		for _, to := range followKindWalk(f, state) {
			mustMove(t, r, s, to, "setup")
		}
	}
	for _, to := range followKindWalk(f, epicState) {
		mustMove(t, r, e, to, "setup")
	}
	return r
}

// followKindSteps are the states the epic entered at t1, in order: its walk
// for the moves under test.
func followKindSteps(t *testing.T, r *Repo) []string {
	t.Helper()
	e, err := r.Get("E-0001")
	if err != nil {
		t.Fatal(err)
	}
	var out []string
	for _, tr := range e.Transitions {
		if tr.At == t1.Format(TimeFormat) {
			out = append(out, tr.To)
		}
	}
	return out
}

// S-0739, ADR-0301 §2: an epic goes to the pull's source with its first
// story there, under either workflow.
func TestFollowKindTheFirstStoryReadyMakesTheEpicReady(t *testing.T) {
	for _, f := range followKindFlows() {
		t.Run(f.name, func(t *testing.T) {
			r := followKindProject(t, f, Backlog, Backlog, Backlog)
			res := move(t, r, "S-0001", f.state(Ready), "")
			want := Followed{ID: "E-0001", Type: Epic, Title: "E", From: f.state(Backlog), To: f.state(Ready), Story: "S-0001"}
			if res.Followed == nil || *res.Followed != want {
				t.Fatalf("followed %+v, want %+v", res.Followed, want)
			}
			if res := move(t, r, "S-0002", f.state(Ready), ""); res.Followed != nil || epicStatus(t, r) != f.state(Ready) {
				t.Errorf("the second story moved the epic: %+v, %s", res.Followed, epicStatus(t, r))
			}
		})
	}
}

// S-0739, ADR-0301 §2: an epic goes to the pull's target with its first
// story started, walking from the initial state one transition at a time.
func TestFollowKindTheFirstStoryStartedStartsTheEpic(t *testing.T) {
	for _, f := range followKindFlows() {
		t.Run(f.name, func(t *testing.T) {
			r := followKindProject(t, f, Backlog, Ready)
			res := move(t, r, "S-0001", f.state(InProgress), "")
			if res.Followed == nil || res.Followed.From != f.state(Backlog) || res.Followed.To != f.state(InProgress) {
				t.Fatalf("followed %+v", res.Followed)
			}
			if got, want := followKindSteps(t, r), []string{f.state(Ready), f.state(InProgress)}; !slices.Equal(got, want) {
				t.Errorf("the epic walked %v, want %v", got, want)
			}
		})
	}
}

// S-0739: under the fixture, a story in the added active state has started
// too, so it starts its epic and holds it there; a move between the two
// active states leaves the epic where it is.
func TestFollowKindAStoryInTheAddedActiveStateHasStarted(t *testing.T) {
	f := followKindFlows()[1]
	if got, ok := EpicFollowsIn(*testFlow(t), inMemory(Backlog, Backlog, fixtureTesting), "E-0001"); got != fixtureDoing || !ok {
		t.Errorf("a story in testing puts the epic in %q %v, want doing", got, ok)
	}
	r := followKindProject(t, f, Ready, Ready)
	move(t, r, "S-0001", fixtureDoing, "")
	if res := move(t, r, "S-0001", fixtureTesting, ""); res.Followed != nil || epicStatus(t, r) != fixtureDoing {
		t.Errorf("a move from doing to testing moved the epic: %+v, %s", res.Followed, epicStatus(t, r))
	}
	r = followKindProject(t, f, InProgress, fixtureTesting, InProgress)
	if res := move(t, r, "S-0002", Ready, ""); res.Followed != nil || epicStatus(t, r) != fixtureDoing {
		t.Errorf("a story in testing does not hold the epic: %+v, %s", res.Followed, epicStatus(t, r))
	}
}

// S-0739, ADR-0301 §2: an epic goes to the acceptance's source with its last
// open story, and under the fixture walks there through testing, one
// transition at a time; a story in testing moving on takes it the same way.
func TestFollowKindTheLastOpenStoryInReviewTakesTheEpicToReview(t *testing.T) {
	for _, f := range followKindFlows() {
		for _, last := range []string{InProgress, fixtureTesting} {
			if last == fixtureTesting && f.name == "baseline" {
				continue
			}
			t.Run(f.name+" from "+last, func(t *testing.T) {
				r := followKindProject(t, f, InProgress, Review, last, Done)
				res := move(t, r, "S-0002", Review, "")
				if res.Followed == nil || res.Followed.From != f.state(InProgress) || res.Followed.To != Review {
					t.Fatalf("followed %+v", res.Followed)
				}
				want := []string{Review}
				if f.name == "fixture" {
					want = []string{fixtureTesting, Review}
				}
				if got := followKindSteps(t, r); !slices.Equal(got, want) {
					t.Errorf("the epic walked %v, want %v", got, want)
				}
			})
		}
	}
}

// S-0739, ADR-0301 §2: the last story accepted takes the epic to the
// accepted state only with the acceptance, else to the acceptance's source.
func TestFollowKindTheLastStoryAcceptedAcceptsTheEpic(t *testing.T) {
	for _, f := range followKindFlows() {
		t.Run(f.name, func(t *testing.T) {
			r := f.open(t)
			items := inMemory(f.state(InProgress), Done, Done)
			without, _, err := r.Follow(items, items[2], Review, "alex", t1, false)
			if err != nil || without == nil || without.Status != Review {
				t.Fatalf("without accept: epic %+v err %v", without, err)
			}
			items = inMemory(f.state(InProgress), Done, Done)
			epic, followed, err := r.Follow(items, items[2], Review, "alex", t1, true)
			if err != nil {
				t.Fatal(err)
			}
			if epic.Status != Done || followed.From != f.state(InProgress) || followed.To != Done {
				t.Errorf("epic %s followed %+v", epic.Status, followed)
			}
			want := []string{Review, Done}
			if f.name == "fixture" {
				want = []string{fixtureTesting, Review, Done}
			}
			var got []string
			for _, tr := range epic.Transitions {
				got = append(got, tr.To)
			}
			if !slices.Equal(got, want) {
				t.Errorf("the epic walked %v, want %v", got, want)
			}
		})
	}
}

// S-0739, ADR-0301 §2: a move back takes the epic back only when no story
// holds it, along a transition that leads back: from review straight to the
// pull's target under the fixture, which has none into testing.
func TestFollowKindBackMovesTakeTheEpicBack(t *testing.T) {
	for _, f := range followKindFlows() {
		t.Run(f.name, func(t *testing.T) {
			r := followKindProject(t, f, Review, Review, Review)
			res := move(t, r, "S-0001", f.state(InProgress), "found a gap")
			if res.Followed == nil || res.Followed.From != Review || res.Followed.To != f.state(InProgress) {
				t.Fatalf("followed %+v", res.Followed)
			}
			if got, want := followKindSteps(t, r), []string{f.state(InProgress)}; !slices.Equal(got, want) {
				t.Errorf("the epic walked %v, want %v", got, want)
			}
			r = followKindProject(t, f, InProgress, InProgress, InProgress)
			if res := move(t, r, "S-0001", Ready, ""); res.Followed != nil || epicStatus(t, r) != f.state(InProgress) {
				t.Errorf("another started story holds the epic: %+v, %s", res.Followed, epicStatus(t, r))
			}
			if res := move(t, r, "S-0002", Ready, ""); res.Followed == nil || epicStatus(t, r) != Ready {
				t.Errorf("no story holds the epic: %+v, %s", res.Followed, epicStatus(t, r))
			}
			r = followKindProject(t, f, Ready, Ready, Ready)
			if res := move(t, r, "S-0001", Backlog, ""); res.Followed != nil || epicStatus(t, r) != Ready {
				t.Errorf("another ready story holds the epic: %+v, %s", res.Followed, epicStatus(t, r))
			}
			if res := move(t, r, "S-0002", Backlog, ""); res.Followed == nil || epicStatus(t, r) != Backlog {
				t.Errorf("no story holds the epic: %+v, %s", res.Followed, epicStatus(t, r))
			}
		})
	}
}

// S-0739, ADR-0301 §2: a story in the cancelled state does not count, a
// cancellation takes the epic forward short of the accepted state, and a
// story back from the cancelled state takes it back.
func TestFollowKindCancellationsUnderEitherWorkflow(t *testing.T) {
	for _, f := range followKindFlows() {
		t.Run(f.name, func(t *testing.T) {
			r := followKindProject(t, f, InProgress, Done, Done, InProgress)
			if res := move(t, r, "S-0003", Cancelled, "not needed"); res.Followed == nil || res.Followed.To != Review {
				t.Errorf("followed %+v", res.Followed)
			}
			want := []string{Review}
			if f.name == "fixture" {
				want = []string{fixtureTesting, Review}
			}
			if got := followKindSteps(t, r); !slices.Equal(got, want) || epicStatus(t, r) != Review {
				t.Errorf("the epic walked %v to %s, want %v to review", got, epicStatus(t, r), want)
			}
			r = followKindProject(t, f, InProgress, InProgress, Cancelled)
			if res := move(t, r, "S-0001", Review, ""); res.Followed == nil || epicStatus(t, r) != Review {
				t.Errorf("followed %+v, epic %s", res.Followed, epicStatus(t, r))
			}
			r = followKindProject(t, f, Review, Review, Cancelled)
			res := move(t, r, "S-0002", Backlog, "")
			if res.Followed == nil || res.Followed.From != Review || res.Followed.To != f.state(InProgress) {
				t.Errorf("followed %+v", res.Followed)
			}
			for _, tc := range []struct {
				states []string
				want   string
				ok     bool
			}{
				{[]string{Cancelled}, "", false},
				{[]string{Backlog, Cancelled}, Backlog, true},
				{[]string{Review, Done, Cancelled}, Review, true},
				{[]string{Done, Done, Cancelled}, Done, true},
			} {
				got, ok := EpicFollowsIn(r.Flow(), inMemory(Backlog, tc.states...), "E-0001")
				if got != tc.want || ok != tc.ok {
					t.Errorf("%v: got %q %v, want %q %v", tc.states, got, ok, tc.want, tc.ok)
				}
			}
		})
	}
}

// S-0739, ADR-0301 §2: a story given another epic takes the one it left as
// though cancelled out of it and the one it joined as though from the
// initial state, under either workflow.
func TestFollowKindAStoryGivenAnotherEpicMovesBoth(t *testing.T) {
	for _, f := range followKindFlows() {
		t.Run(f.name, func(t *testing.T) {
			r := followKindProject(t, f, InProgress, InProgress, Review)
			other := mustCreate(t, r, Epic, "Other", "")
			all, err := r.List(true)
			if err != nil {
				t.Fatal(err)
			}
			story, err := r.Get("S-0001")
			if err != nil {
				t.Fatal(err)
			}
			story.Parent = other.ID
			epics, walks, err := r.FollowParent(all, story, "E-0001", "alex", t1)
			if err != nil {
				t.Fatal(err)
			}
			if len(walks) != 2 {
				t.Fatalf("walks %+v", walks)
			}
			if w := walks[0]; w.ID != "E-0001" || w.From != f.state(InProgress) || w.To != Review || epics[0].Status != Review {
				t.Errorf("the epic it left: %+v", w)
			}
			if w := walks[1]; w.ID != other.ID || w.From != Backlog || w.To != f.state(InProgress) || epics[1].Status != f.state(InProgress) {
				t.Errorf("the epic it joined: %+v", w)
			}
		})
	}
}

// S-0739, ADR-0301 §2: the baseline's EpicFollows and EpicFollowsIn under
// either workflow agree on every mix of stories, in each one's names.
func TestFollowKindEpicFollowsAgreesUnderEitherWorkflow(t *testing.T) {
	for _, tc := range []struct {
		states []string
		want   string
	}{
		{[]string{Backlog}, Backlog},
		{[]string{Backlog, Ready}, Ready},
		{[]string{Ready, InProgress}, InProgress},
		{[]string{Ready, Done}, InProgress},
		{[]string{Review, InProgress}, InProgress},
		{[]string{Review, Done}, Review},
		{[]string{Done}, Done},
	} {
		for _, f := range []struct {
			name string
			w    workflow.Workflow
			st   func(string) string
		}{{"baseline", workflow.Baseline(), func(s string) string { return s }}, {"fixture", *testFlow(t), inTestFlow}} {
			var states []string
			for _, s := range tc.states {
				states = append(states, f.st(s))
			}
			got, ok := EpicFollowsIn(f.w, inMemory(f.st(Backlog), states...), "E-0001")
			if !ok || got != f.st(tc.want) {
				t.Errorf("%s %v: got %q %v, want %q", f.name, states, got, ok, f.st(tc.want))
			}
			if f.name == "baseline" {
				if base, _ := EpicFollows(inMemory(Backlog, states...), "E-0001"); base != got {
					t.Errorf("%v: EpicFollows %q, EpicFollowsIn %q", states, base, got)
				}
			}
		}
	}
}
