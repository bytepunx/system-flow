package workitem

import (
	"errors"
	"fmt"
	"os"
	"slices"
	"strings"
	"sync"
	"testing"
	"time"
)

// readyStory creates S-0001 with a criterion and moves it to ready.
func readyStory(t *testing.T, r *Repo) {
	t.Helper()
	s := mustCreate(t, r, Story, "Planned while started", "")
	s.Body = strings.Replace(s.Body, "## Acceptance criteria\n- [ ]\n", "## Acceptance criteria\n- [ ] works\n", 1)
	if err := r.Save(s); err != nil {
		t.Fatal(err)
	}
	mustMove(t, r, s, Ready, "")
}

// copies reads n copies of S-0001 through the parse cache, every one before
// the caller writes any, and makes sure each is the file as it is.
func copies(t *testing.T, r *Repo, n int) []*Item {
	t.Helper()
	out := make([]*Item, n)
	for i := range out {
		it, err := r.Get("S-0001")
		if err != nil {
			t.Fatal(err)
		}
		if it.Marshal() != onDisk(t, it.Path).Marshal() {
			t.Fatalf("copy %d is not the file as it is", i)
		}
		out[i] = it
	}
	return out
}

// onDisk reads the item from its file, past the parse cache.
func onDisk(t *testing.T, path string) *Item {
	t.Helper()
	it, err := ReadItem(path)
	if err != nil {
		t.Fatal(err)
	}
	return it
}

func stamp(d time.Duration) string { return t0.Add(d).UTC().Format(TimeFormat) }

// I-0175: the replanner's forecast write read S-0463 while it was ready,
// flai story start moved it to in-progress meanwhile, and the forecast write
// then saved its stale copy over the move.
func TestSaveKeepsAMoveSavedSinceTheCopyWasRead(t *testing.T) {
	r := newProject(t)
	readyStory(t, r)
	c := copies(t, r, 2)
	starter, planner := c[0], c[1]
	if _, err := r.TransitionAll(starter, InProgress, "agent", "", t0.Add(time.Minute), false); err != nil {
		t.Fatal(err)
	}
	planner.Forecast = &Forecast{Duration: "4h", Basis: "three stories like it", By: "planner", At: stamp(2 * time.Minute)}
	planner.Updated = stamp(2 * time.Minute)
	if err := r.Save(planner); err != nil {
		t.Fatalf("the forecast write was refused: %v", err)
	}
	got := onDisk(t, planner.Path)
	if got.Status != InProgress || got.Transitions[len(got.Transitions)-1].To != InProgress {
		t.Errorf("the forecast write lost the move: status %s, transitions %+v", got.Status, got.Transitions)
	}
	if got.Forecast.IsZero() || got.Forecast.Duration != "4h" {
		t.Errorf("the forecast was not written: %+v", got.Forecast)
	}
	if got.Updated != stamp(2*time.Minute) {
		t.Errorf("updated %s, want the later write's %s", got.Updated, stamp(2*time.Minute))
	}
	if planner.Status != InProgress || planner.Marshal() != got.Marshal() {
		t.Errorf("the saved copy is not what was written: status %s", planner.Status)
	}
}

// assertUpdatedInOrder fails unless the file's updated is at or after its
// created and its last transition's at, as flai check's item.chronology asks.
func assertUpdatedInOrder(t *testing.T, it *Item, want string) {
	t.Helper()
	got := onDisk(t, it.Path)
	if got.Updated != want {
		t.Errorf("updated %s, want %s", got.Updated, want)
	}
	if got.Updated < got.Created {
		t.Errorf("updated %s is before created %s", got.Updated, got.Created)
	}
	if n := len(got.Transitions); n > 0 && got.Updated < got.Transitions[n-1].At {
		t.Errorf("updated %s is before the last transition's %s", got.Updated, got.Transitions[n-1].At)
	}
	if it.Marshal() != got.Marshal() {
		t.Errorf("the saved copy is not what was written: updated %s", it.Updated)
	}
}

// I-0183: a writer whose clock was read before a move it then read the item
// after saved an updated earlier than the move's transition.
func TestSaveKeepsUpdatedAtTheLastTransitionWhenTheWriterStampIsEarlier(t *testing.T) {
	r := newProject(t)
	readyStory(t, r)
	mover := copies(t, r, 1)[0]
	if _, err := r.TransitionAll(mover, InProgress, "agent", "", t0.Add(time.Minute), false); err != nil {
		t.Fatal(err)
	}
	planner := copies(t, r, 1)[0]
	planner.Forecast = &Forecast{Duration: "4h", Basis: "three stories like it", By: "planner", At: stamp(30 * time.Second)}
	planner.Updated = stamp(30 * time.Second)
	if err := r.Save(planner); err != nil {
		t.Fatalf("the forecast write was refused: %v", err)
	}
	assertUpdatedInOrder(t, planner, stamp(time.Minute))
}

// A move in the same second as the copy's base leaves updated as base had it,
// so the merge takes the writer's earlier stamp, the only side that changed it.
func TestSaveKeepsUpdatedAtTheLastTransitionAfterAMergeOfASameSecondMove(t *testing.T) {
	r := newProject(t)
	readyStory(t, r)
	s := copies(t, r, 1)[0]
	s.Touches = []string{"a.go"}
	s.Updated = stamp(time.Minute)
	if err := r.Save(s); err != nil {
		t.Fatal(err)
	}
	c := copies(t, r, 2)
	mover, planner := c[0], c[1]
	if _, err := r.TransitionAll(mover, InProgress, "agent", "", t0.Add(time.Minute+500*time.Millisecond), false); err != nil {
		t.Fatal(err)
	}
	if got := onDisk(t, mover.Path).Updated; got != stamp(time.Minute) {
		t.Fatalf("the move's updated %s, want the base's %s", got, stamp(time.Minute))
	}
	planner.Forecast = &Forecast{Duration: "4h", Basis: "three stories like it", By: "planner", At: stamp(30 * time.Second)}
	planner.Updated = stamp(30 * time.Second)
	if err := r.Save(planner); err != nil {
		t.Fatalf("the forecast write was refused: %v", err)
	}
	if got := onDisk(t, planner.Path); got.Status != InProgress || got.Forecast.IsZero() {
		t.Errorf("the merge lost a change: status %s, forecast %+v", got.Status, got.Forecast)
	}
	assertUpdatedInOrder(t, planner, stamp(time.Minute))
}

func TestSaveKeepsUpdatedAtCreated(t *testing.T) {
	r := newProject(t)
	s := mustCreate(t, r, Story, "Stamped early", "")
	if len(s.Transitions) != 0 {
		t.Fatalf("a new story has transitions %+v", s.Transitions)
	}
	s.Touches = []string{"a.go"}
	s.Updated = stamp(-time.Minute)
	if err := r.Save(s); err != nil {
		t.Fatal(err)
	}
	assertUpdatedInOrder(t, s, stamp(0))
}

func TestSaveRefusesAKeyBothChanged(t *testing.T) {
	r := newProject(t)
	readyStory(t, r)
	c := copies(t, r, 2)
	starter, canceller := c[0], c[1]
	if _, err := r.TransitionAll(starter, InProgress, "agent", "", t0.Add(time.Minute), false); err != nil {
		t.Fatal(err)
	}
	written := onDisk(t, starter.Path).Marshal()
	items, _ := r.List(false)
	board, _ := r.LoadBoard()
	if _, err := r.Move(canceller, Cancelled, MoveOptions{Now: t0.Add(2 * time.Minute), Items: items, Board: board, Reason: "not wanted", By: "test"}); err != nil {
		t.Fatal(err)
	}
	err := r.Save(canceller)
	var conflict *ConflictError
	if !errors.As(err, &conflict) {
		t.Fatalf("saving a second move of a copy read before the first: %v, want a *ConflictError", err)
	}
	if conflict.ID != "S-0001" || !slices.Contains(conflict.Keys, "status") || !slices.Contains(conflict.Keys, "transitions") {
		t.Errorf("conflict %+v, want S-0001 with status and transitions", conflict)
	}
	for _, want := range []string{"S-0001", "another writer", "status", "nothing was written", "read S-0001 again"} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("error %q does not say %q", err, want)
		}
	}
	if got := onDisk(t, starter.Path).Marshal(); got != written {
		t.Errorf("a refused save changed the file:\n%s\nwant\n%s", got, written)
	}
}

// Writers that each read their own copy before any writes, then change a
// key of their own and save at once, all keep their change.
func TestSavesAtOnceOfCopiesReadTogetherKeepEveryChange(t *testing.T) {
	r := newProject(t)
	readyStory(t, r)
	changes := []func(it *Item){
		func(it *Item) { it.Tags = []string{"planned"} },
		func(it *Item) { it.Touches = []string{"flai/internal/workitem"} },
		func(it *Item) { it.Topics = []string{"locking"} },
		func(it *Item) { it.After = []string{"S-0002"} },
		func(it *Item) { it.Body = strings.Replace(it.Body, "## Goal\n", "## Goal\n\nWritten at once.\n", 1) },
		func(it *Item) { it.Forecast = &Forecast{Duration: "2h", By: "planner", At: stamp(time.Hour)} },
		func(it *Item) { it.Estimate = "3h" },
		func(it *Item) { it.Stream = "locking" },
	}
	c := copies(t, r, len(changes))
	var wg sync.WaitGroup
	errs := make(chan error, len(changes))
	for i, change := range changes {
		wg.Go(func() {
			change(c[i])
			c[i].Updated = stamp(time.Duration(i+1) * time.Minute)
			if err := r.Save(c[i]); err != nil {
				errs <- fmt.Errorf("writer %d: %w", i, err)
			}
		})
	}
	wg.Wait()
	close(errs)
	for err := range errs {
		t.Error(err)
	}
	got := onDisk(t, c[0].Path)
	checks := map[string]bool{
		"tags":     slices.Equal(got.Tags, []string{"planned"}),
		"touches":  slices.Equal(got.Touches, []string{"flai/internal/workitem"}),
		"topics":   slices.Equal(got.Topics, []string{"locking"}),
		"after":    slices.Equal(got.After, []string{"S-0002"}),
		"body":     strings.Contains(got.Body, "Written at once."),
		"forecast": !got.Forecast.IsZero() && got.Forecast.Duration == "2h",
		"estimate": got.Estimate == "3h",
		"stream":   got.Stream == "locking",
		"updated":  got.Updated == stamp(time.Duration(len(changes))*time.Minute),
		"status":   got.Status == Ready,
	}
	for key, ok := range checks {
		if !ok {
			t.Errorf("%s was lost:\n%s", key, got.Marshal())
		}
	}
}

// A note one copy adds and a child line another writer added since, each in
// a section of its own, are both kept; two notes are refused.
func TestSaveMergesTheBodyBySection(t *testing.T) {
	r := newProject(t)
	readyStory(t, r)
	c := copies(t, r, 2)
	mustCreate(t, r, Task, "Added since", "S-0001")
	items, _ := r.List(false)
	board, _ := r.LoadBoard()
	if _, err := r.Move(c[0], InProgress, MoveOptions{Now: t0.Add(time.Minute), Items: items, Board: board, Reason: "started", By: "test"}); err != nil {
		t.Fatal(err)
	}
	if err := r.Save(c[0]); err != nil {
		t.Fatalf("a note while a child was added: %v", err)
	}
	got := onDisk(t, c[0].Path)
	if !strings.Contains(got.Body, "## Tasks\n- T-0001 Added since\n") || !strings.Contains(got.Body, "moved to in-progress: started") {
		t.Errorf("the child line or the note was lost:\n%s", got.Body)
	}
	c[1].Body = strings.Replace(c[1].Body, "## Notes\n", "## Notes\n- another note\n", 1)
	err := r.Save(c[1])
	var conflict *ConflictError
	if !errors.As(err, &conflict) || !slices.Equal(conflict.Keys, []string{"body section ## Notes"}) {
		t.Errorf("two notes in one section: %v, want a conflict on body section ## Notes", err)
	}
}

// A copy saved again after another change is not refused for its own first
// write, even when it changes the same key again.
func TestSaveAgainOfTheSameCopy(t *testing.T) {
	r := newProject(t)
	readyStory(t, r)
	it := copies(t, r, 1)[0]
	it.Tags = []string{"first"}
	if err := r.Save(it); err != nil {
		t.Fatal(err)
	}
	it.Tags = []string{"second"}
	it.Touches = []string{"flai/cmd"}
	if err := r.Save(it); err != nil {
		t.Fatalf("a second save of one copy: %v", err)
	}
	got := onDisk(t, it.Path)
	if !slices.Equal(got.Tags, []string{"second"}) || !slices.Equal(got.Touches, []string{"flai/cmd"}) {
		t.Errorf("tags %v, touches %v, want [second] and [flai/cmd]", got.Tags, got.Touches)
	}
}

// Restore refuses to take out a change another writer changed again since,
// writing nothing, and gives a file that is gone back as it was (S-0479).
func TestRestoreRefusesAKeyChangedAgainSince(t *testing.T) {
	r := newProject(t)
	readyStory(t, r)
	it := copies(t, r, 1)[0]
	before, err := os.ReadFile(it.Path)
	if err != nil {
		t.Fatal(err)
	}
	it.Tags = []string{"edited"}
	if err := r.Save(it); err != nil {
		t.Fatal(err)
	}
	wrote := []byte(it.Marshal())
	again := copies(t, r, 1)[0]
	again.Tags = []string{"changed again"}
	if err := r.Save(again); err != nil {
		t.Fatal(err)
	}
	written := onDisk(t, it.Path).Marshal()
	err = r.Restore(it.Path, wrote, before)
	var conflict *ConflictError
	if !errors.As(err, &conflict) || conflict.ID != "S-0001" || !slices.Equal(conflict.Keys, []string{"tags"}) {
		t.Fatalf("restore over tags changed again: %v, want a conflict on S-0001's tags", err)
	}
	if got := onDisk(t, it.Path).Marshal(); got != written {
		t.Errorf("a refused restore changed the file:\n%s\nwant\n%s", got, written)
	}
	if err := os.Remove(it.Path); err != nil {
		t.Fatal(err)
	}
	if err := r.Restore(it.Path, wrote, before); err != nil {
		t.Fatalf("restore of a file that is gone: %v", err)
	}
	if got, _ := os.ReadFile(it.Path); string(got) != string(before) {
		t.Errorf("restore of a file that is gone wrote\n%s\nwant\n%s", got, before)
	}
}
