package mcpserver

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/threads"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// held calls wait_for_work in the background and returns its answer.
func (f *fixture) held(t *testing.T, timeout int) <-chan map[string]any {
	t.Helper()
	done := make(chan map[string]any, 1)
	go func() {
		out, failed := f.call(t, "wait_for_work", map[string]any{"timeout_seconds": timeout})
		if failed != "" {
			t.Errorf("wait_for_work: %s", failed)
		}
		done <- out
	}()
	return done
}

func answered(t *testing.T, done <-chan map[string]any) map[string]any {
	t.Helper()
	select {
	case out := <-done:
		return out
	case <-time.After(5 * time.Second):
		t.Fatal("wait_for_work never answered")
		return nil
	}
}

func storyOf(out map[string]any) string {
	if s, ok := out["story"].(map[string]any); ok {
		return s["id"].(string)
	}
	return ""
}

// toReview moves the fixture's own story on, so its agent is idle.
func (f *fixture) toReview(t *testing.T) {
	t.Helper()
	s, _ := f.repo.Get(f.story.ID)
	if _, err := f.repo.Transition(s, workitem.Review, "claude", "", t0); err != nil {
		t.Fatal(err)
	}
}

func (f *fixture) limitInProgress(t *testing.T, n int) {
	t.Helper()
	board := "---\ntitle: Board\nwip_limits:\n  in-progress: " + string(rune('0'+n)) + "\norder: []\n---\n\n# Board\n"
	if err := os.WriteFile(filepath.Join(f.repo.KanbanDir(), "board.md"), []byte(board), 0o644); err != nil {
		t.Fatal(err)
	}
}

// S-0097: an agent whose own story is still in progress is sent back to it,
// not to a new one.
func TestWaitForWorkResumesYourOwnStory(t *testing.T) {
	f := setup(t)
	f.readyStory(t, "Next", t0)
	out := answered(t, f.held(t, 3))
	if out["reason"] != "resume" || storyOf(out) != f.story.ID {
		t.Errorf("with my story in progress: %v", out)
	}
}

// Criterion 3: nothing ready, so it waits; a story made ready is returned to
// pull within a second.
func TestWaitForWorkWaitsForAStoryToBeReady(t *testing.T) {
	f := setup(t)
	f.toReview(t)
	quiet, _ := f.call(t, "wait_for_work", map[string]any{"timeout_seconds": 1})
	if quiet["timed_out"] != true || quiet["waiting_for"] != "ready" || quiet["reason"] != "" {
		t.Fatalf("nothing ready: %v", quiet)
	}
	done := f.held(t, 3)
	time.Sleep(150 * time.Millisecond)
	next := f.readyStory(t, "Next", t0.Add(2*time.Minute))
	start := time.Now()
	out := answered(t, done)
	if out["reason"] != "pull" || storyOf(out) != next.ID || out["can_pull"] != true {
		t.Errorf("once one is ready: %v", out)
	}
	if time.Since(start) > 1500*time.Millisecond {
		t.Errorf("heard of it late: %v", time.Since(start))
	}
}

// Criteria 1 and 2: with room it answers at once with the first story in
// pull order; with the in-progress limit full it waits for room, and answers
// once another agent's story leaves in-progress.
func TestWaitForWorkWaitsForRoomUnderTheLimit(t *testing.T) {
	f := setup(t)
	f.toReview(t)
	first := f.readyStory(t, "First", t0)
	f.readyStory(t, "Second", t0)
	if out := answered(t, f.held(t, 1)); out["reason"] != "pull" || storyOf(out) != first.ID {
		t.Fatalf("with room: %v", out)
	}

	// another agent works a story; the limit is one
	theirs := f.readyStory(t, "Theirs", t0)
	theirs, _ = f.repo.Get(theirs.ID)
	if _, err := f.repo.Transition(theirs, workitem.InProgress, "codex", "", t0); err != nil {
		t.Fatal(err)
	}
	if _, err := f.repo.OpenStream(theirs, workitem.StreamOptions{Agent: "codex", Now: t0}); err != nil {
		t.Fatal(err)
	}
	f.limitInProgress(t, 1)
	quiet := answered(t, f.held(t, 1))
	if quiet["timed_out"] != true || quiet["waiting_for"] != "room" || quiet["can_pull"] != false || len(quiet["ready"].([]any)) != 2 {
		t.Fatalf("limit full: %v", quiet)
	}

	done := f.held(t, 3)
	time.Sleep(150 * time.Millisecond)
	if _, err := f.repo.Create(workitem.NewOptions{Type: workitem.Task, Title: "T", Parent: theirs.ID, Owner: "alex", Now: t0}); err != nil {
		t.Fatal(err)
	}
	theirs, _ = f.repo.Get(theirs.ID)
	if _, err := f.repo.Transition(theirs, workitem.Review, "codex", "", t0); err != nil {
		t.Fatal(err)
	}
	if out := answered(t, done); out["reason"] != "pull" || storyOf(out) != first.ID {
		t.Errorf("once there is room: %v", out)
	}
}

// S-0128: a ready story whose claim overlaps an open story's is passed over
// for the next clear one; when every ready story is held it waits, says so,
// and offers the held one as soon as it is clear. The fixture's story, in
// review, touches flai/internal/mcpserver.
func TestWaitForWorkPassesOverAHeldStory(t *testing.T) {
	f := setup(t)
	f.toReview(t)
	held := f.readyStory(t, "Held", t0, "flai/internal")
	clear := f.readyStory(t, "Clear", t0)
	if out := answered(t, f.held(t, 1)); out["reason"] != "pull" || storyOf(out) != clear.ID {
		t.Fatalf("the clear story: %v", out)
	}
	c, _ := f.repo.Get(clear.ID)
	if _, err := f.repo.Transition(c, workitem.InProgress, "codex", "", t0); err != nil {
		t.Fatal(err)
	}
	if _, err := f.repo.OpenStream(c, workitem.StreamOptions{Agent: "codex", Now: t0}); err != nil {
		t.Fatal(err)
	}
	f.limitInProgress(t, 3)
	quiet := answered(t, f.held(t, 1))
	ready := quiet["ready"].([]any)
	if quiet["timed_out"] != true || quiet["waiting_for"] != "held" || quiet["can_pull"] != true || len(ready) != 1 {
		t.Fatalf("every ready story held: %v", quiet)
	}
	if h, _ := ready[0].(map[string]any)["held"].(map[string]any); h["code"] != "overlap" || !strings.Contains(h["reason"].(string), "starts when "+f.story.ID+" is accepted, cancelled, or sent back") {
		t.Errorf("held: %v", ready[0])
	}

	// the open story's claim narrows: the held story is clear, and offered
	done := f.held(t, 3)
	time.Sleep(150 * time.Millisecond)
	s, _ := f.repo.Get(f.story.ID)
	s.Touches = []string{"design/system"}
	if err := f.repo.Save(s); err != nil {
		t.Fatal(err)
	}
	if out := answered(t, done); out["reason"] != "pull" || storyOf(out) != held.ID {
		t.Errorf("once clear: %v", out)
	}
}

// S-0128: item_move warns on a held story and moves it all the same.
func TestItemMoveWarnsOnAHeldStory(t *testing.T) {
	f := setup(t)
	held := f.readyStory(t, "Held", t0, "flai/internal/mcpserver/work.go")
	out, failed := f.call(t, "item_move", map[string]any{"id": held.ID, "to": "in-progress"})
	if failed != "" || out["status"] != "in-progress" {
		t.Fatalf("move: %v %s", out, failed)
	}
	want := held.ID + " is held (overlap): touches flai/internal/mcpserver/work.go, inside flai/internal/mcpserver which " + f.story.ID + " (in progress) touches"
	if w := out["warnings"].([]any); len(w) != 1 || !strings.HasPrefix(w[0].(string), want) {
		t.Errorf("warnings: %v", w)
	}
}

// A thread written to while it waits is work too; one that was already
// awaiting the agent before does not wake it again and again.
func TestWaitForWorkWakesForAThreadWrittenWhileWaiting(t *testing.T) {
	f := setup(t)
	f.toReview(t)
	if _, err := threads.New(f.repo, threads.NewOptions{Title: "Old", On: f.story.ID, Author: "alex", Text: "Long ago.", Now: t0}); err != nil {
		t.Fatal(err)
	}
	quiet, _ := f.call(t, "wait_for_work", map[string]any{"timeout_seconds": 1})
	if quiet["reason"] != "" || len(quiet["threads"].([]any)) != 0 {
		t.Fatalf("an old thread woke it: %v", quiet)
	}
	done := f.held(t, 3)
	time.Sleep(150 * time.Millisecond)
	th, err := threads.New(f.repo, threads.NewOptions{Title: "Question", On: f.story.ID, Author: "alex", Text: "Which port?", Now: t0.Add(2 * time.Minute)})
	if err != nil {
		t.Fatal(err)
	}
	out := answered(t, done)
	ths := out["threads"].([]any)
	if out["reason"] != "thread" || len(ths) != 1 || ths[0].(map[string]any)["id"] != th.ID {
		t.Errorf("a thread written while waiting: %v", out)
	}
}
