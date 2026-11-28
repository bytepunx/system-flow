package workitem

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"
)

// readyStories makes n stories with a ticked criterion, outside any epic so
// that no move takes an epic along, and moves each to ready, which names it
// in the pull order. The in-progress WIP limit is lifted so that no move of
// them is refused.
func readyStories(t *testing.T, r *Repo, n int) []*Item {
	t.Helper()
	var out []*Item
	for i := range n {
		s := mustCreate(t, r, Story, fmt.Sprintf("S%d", i), "")
		s.Body = strings.Replace(s.Body, "## Acceptance criteria\n- [ ]\n", "## Acceptance criteria\n- [x] works\n", 1)
		if err := r.Save(s); err != nil {
			t.Fatal(err)
		}
		mustMove(t, r, s, Ready, "")
		out = append(out, s)
	}
	board, err := r.LoadBoard()
	if err != nil {
		t.Fatal(err)
	}
	for _, col := range stateNames(limitedStates(board.Flow())) {
		if err := board.SetLimit(col, 0); err != nil {
			t.Fatal(err)
		}
	}
	if err := board.Save("2026-09-15"); err != nil {
		t.Fatal(err)
	}
	return out
}

// I-0375: stories moved to in-progress at once each read the board, took
// themselves out of the order, and wrote it back, one over another, so a
// story stayed in the order while in progress. Under the board lock each
// move reads the board the one before it wrote (ADR-0306).
func TestStoriesStartedAtOnceAllLeaveTheOrder(t *testing.T) {
	r := newProject(t)
	stories := readyStories(t, r, 6)
	board, err := r.LoadBoard()
	if err != nil {
		t.Fatal(err)
	}
	if len(board.Order) != len(stories) {
		t.Fatalf("order before the moves: %v, want all %d stories", board.Order, len(stories))
	}

	// Each item save holds the items lock a while, as a slow disk would,
	// so that without the board lock every move reads the board before any
	// writes it, and the race I-0375 lost is lost every run.
	orig := lockFile
	t.Cleanup(func() { lockFile = orig })
	lockFile = func(path string, wait time.Duration) (func(), error) {
		unlock, err := orig(path, wait)
		if err == nil && filepath.Base(path) == "items.lock" {
			time.Sleep(20 * time.Millisecond)
		}
		return unlock, err
	}

	start := make(chan struct{})
	errs := make(chan error, len(stories))
	var wg sync.WaitGroup
	for _, s := range stories {
		it, err := r.Get(s.ID)
		if err != nil {
			t.Fatal(err)
		}
		wg.Go(func() {
			<-start
			if _, err := r.TransitionAll(it, InProgress, "test", "", t0.Add(time.Hour), false); err != nil {
				errs <- fmt.Errorf("move %s: %w", it.ID, err)
			}
		})
	}
	close(start)
	wg.Wait()
	close(errs)
	for err := range errs {
		t.Fatal(err)
	}

	board, err = r.LoadBoard()
	if err != nil {
		t.Fatal(err)
	}
	if len(board.Order) != 0 {
		t.Errorf("order after every story moved to in-progress at once: %v, want none of them", board.Order)
	}
	for id, st := range statuses(t, r) {
		if st != InProgress {
			t.Errorf("%s is %s, want %s", id, st, InProgress)
		}
	}
}

// A write through UpdateBoard drops from the order an in-progress story an
// earlier write left in it and a name with no item, with their placements,
// keeps a story whose file does not parse, and keeps the ready and backlog
// names in their order.
func TestUpdateBoardDropsStaleNamesFromTheOrder(t *testing.T) {
	r := newProject(t)
	stories := readyStories(t, r, 3) // S-0001, S-0002, S-0003, ready
	mustMove(t, r, stories[2], InProgress, "")
	backlog := mustCreate(t, r, Story, "Backlog", "") // S-0004
	broken := mustCreate(t, r, Story, "Broken", "")   // S-0005
	if err := os.WriteFile(broken.Path, []byte("PLACEHOLDER"), 0o644); err != nil {
		t.Fatal(err)
	}

	board, err := r.LoadBoard()
	if err != nil {
		t.Fatal(err)
	}
	board.Order = []string{"S-0002", "S-0003", "S-0001", "S-0099", "S-0005", backlog.ID}
	board.RecordPlacement("S-0003", "alex", t0)
	board.RecordPlacement("S-0099", "alex", t0)
	board.RecordPlacement("S-0098", "alex", t0) // placed, not in the order
	board.RecordPlacement("S-0001", "alex", t0)
	if err := board.Save("2026-09-15"); err != nil {
		t.Fatal(err)
	}

	var gotItems int
	if err := r.UpdateBoard("2026-09-16", func(b *Board, items []*Item) error {
		gotItems = len(items)
		return b.SetLimit(Ready, 4)
	}); err != nil {
		t.Fatal(err)
	}
	if gotItems != 4 {
		t.Errorf("change was given %d items, want the 4 active ones that parse", gotItems)
	}
	board, err = r.LoadBoard()
	if err != nil {
		t.Fatal(err)
	}
	if want := []string{"S-0002", "S-0001", "S-0005", backlog.ID}; !reflect.DeepEqual(board.Order, want) {
		t.Errorf("order: got %v want %v", board.Order, want)
	}
	if _, ok := board.Placed["S-0003"]; ok {
		t.Errorf("the in-progress story's placement stayed: %v", board.Placed)
	}
	if _, ok := board.Placed["S-0099"]; ok {
		t.Errorf("the placement of a name with no item stayed: %v", board.Placed)
	}
	if _, ok := board.Placed["S-0098"]; ok {
		t.Errorf("a placement of a name the order no longer holds stayed: %v", board.Placed)
	}
	if _, ok := board.Placed["S-0001"]; !ok {
		t.Errorf("the ready story's placement went: %v", board.Placed)
	}
	if board.WIPLimits[Ready] != 4 || board.Updated != "2026-09-16" {
		t.Errorf("the change was not written: limits %v, updated %s", board.WIPLimits, board.Updated)
	}
}

// An error from the change writes nothing and comes back as it is.
func TestUpdateBoardWritesNothingWhenTheChangeFails(t *testing.T) {
	r := newProject(t)
	readyStories(t, r, 1)
	path := filepath.Join(r.KanbanDir(), BoardFile)
	before := readFile(t, path)
	refusal := errors.New("rule: no")
	err := r.UpdateBoard("2026-09-16", func(b *Board, _ []*Item) error {
		b.Order = nil
		return refusal
	})
	if !errors.Is(err, refusal) || err.Error() != refusal.Error() {
		t.Errorf("err = %v, want the change's own", err)
	}
	if after := readFile(t, path); after != before {
		t.Errorf("board written on a failed change:\n%s", after)
	}
}
