package mcpserver

import (
	"fmt"
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

// I-0059: wait_for_work holds for the timeout asked, up to 30 minutes, and
// for five minutes when none is asked.
func TestWaitForWorkHoldsUpToTheLongestWait(t *testing.T) {
	d := newDeadlines()
	f := setupWith(t, func(o *Options) { o.MaxWait, o.After = 0, d.after })
	f.toReview(t)
	checkHolds(t, d, func(args map[string]any) map[string]any {
		out, failed := f.call(t, "wait_for_work", args)
		if failed != "" {
			t.Errorf("wait_for_work: %s", failed)
		}
		return out
	}, 5*time.Minute)
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

func (f *fixture) limits(t *testing.T, inProgress, review int) {
	t.Helper()
	board := fmt.Sprintf("---\ntitle: Board\nwip_limits:\n  in-progress: %d\n  review: %d\norder: []\n---\n\n# Board\n", inProgress, review)
	if err := os.WriteFile(filepath.Join(f.repo.KanbanDir(), "board.md"), []byte(board), 0o644); err != nil {
		t.Fatal(err)
	}
}

// S-0243, I-0007: while review is at its limit no story is pulled, whatever
// room the in-progress limit leaves: wait_for_work waits for review and says
// why, inbox says the same, and a story is offered once one leaves review.
func TestWaitForWorkWaitsWhileReviewIsFull(t *testing.T) {
	f := setup(t)
	f.toReview(t) // the fixture's story is the one in review
	next := f.readyStory(t, "Next", t0)
	f.limits(t, 3, 1)
	full := "review is full (1 of 1): accept or send back a story"
	quiet := answered(t, f.held(t, 1))
	if quiet["timed_out"] != true || quiet["waiting_for"] != "review" || quiet["can_pull"] != false || quiet["pull_hold"] != full || len(quiet["ready"].([]any)) != 1 {
		t.Fatalf("review full: %v", quiet)
	}
	if in, _ := f.call(t, "inbox", map[string]any{}); in["can_pull"] != false || in["pull_hold"] != full {
		t.Errorf("inbox with review full: %v %v", in["can_pull"], in["pull_hold"])
	}

	// both full: the in-progress limit is named first, and it waits for room
	theirs := f.readyStory(t, "Theirs", t0)
	if _, err := f.repo.Transition(theirs, workitem.InProgress, "codex", "", t0); err != nil {
		t.Fatal(err)
	}
	f.limits(t, 1, 1)
	quiet = answered(t, f.held(t, 1))
	if quiet["waiting_for"] != "room" || quiet["pull_hold"] != "the in-progress limit is full (1 of 1)" {
		t.Fatalf("both full: %v", quiet)
	}

	// the in-progress limit has room again, and the story in review is
	// accepted, its task dropped: review has room, and the ready story is offered
	f.limits(t, 3, 1)
	items, _ := f.repo.List(false)
	for _, it := range items {
		if it.Parent == f.story.ID && it.Type == workitem.Task {
			if _, err := f.repo.Transition(it, workitem.Cancelled, "alex", "not needed", t0); err != nil {
				t.Fatal(err)
			}
		}
	}
	done := f.held(t, 3)
	time.Sleep(150 * time.Millisecond)
	mine, _ := f.repo.Get(f.story.ID)
	if _, err := f.repo.Transition(mine, workitem.Done, "alex", "", t0); err != nil {
		t.Fatal(err)
	}
	out := answered(t, done)
	if out["reason"] != "pull" || storyOf(out) != next.ID || out["can_pull"] != true || out["pull_hold"] != nil {
		t.Errorf("once review has room: %v", out)
	}
}

// S-0128: a ready story whose claim overlaps that of a story in progress is
// passed over for the next clear one; when every ready story is held it
// waits, says so, and offers the held one as soon as it is clear. Another
// agent's story, in progress, touches flai/internal/mcpserver.
func TestWaitForWorkPassesOverAHeldStory(t *testing.T) {
	f := setup(t)
	f.toReview(t)
	theirs := f.readyStory(t, "Theirs", t0, "flai/internal/mcpserver")
	if _, err := f.repo.Transition(theirs, workitem.InProgress, "gemini", "", t0); err != nil {
		t.Fatal(err)
	}
	if _, err := f.repo.OpenStream(theirs, workitem.StreamOptions{Agent: "gemini", Now: t0}); err != nil {
		t.Fatal(err)
	}
	held := f.readyStory(t, "Held", t0, "flai/internal")
	clear := f.readyStory(t, "Clear", t0)
	f.limitInProgress(t, 3)
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
	quiet := answered(t, f.held(t, 1))
	ready := quiet["ready"].([]any)
	if quiet["timed_out"] != true || quiet["waiting_for"] != "held" || quiet["can_pull"] != true || len(ready) != 1 {
		t.Fatalf("every ready story held: %v", quiet)
	}
	if h, _ := ready[0].(map[string]any)["held"].(map[string]any); h["code"] != "overlap" || !strings.Contains(h["reason"].(string), "starts when "+theirs.ID+" moves to review, is cancelled, or is sent back") {
		t.Errorf("held: %v", ready[0])
	}

	// the story in progress narrows its claim: the held story is clear, and offered
	done := f.held(t, 3)
	time.Sleep(150 * time.Millisecond)
	s, _ := f.repo.Get(theirs.ID)
	s.Touches = []string{"design/system"}
	if err := f.repo.Save(s); err != nil {
		t.Fatal(err)
	}
	if out := answered(t, done); out["reason"] != "pull" || storyOf(out) != held.ID {
		t.Errorf("once clear: %v", out)
	}
}

// S-0295, I-0087, ADR-0096: a story in review holds nothing. A ready story
// whose claim overlaps only that of the fixture's story, in review and
// touching flai/internal/mcpserver, is not held, and wait_for_work offers it
// first.
func TestWaitForWorkOffersAStoryOverlappingOnlyOneInReview(t *testing.T) {
	f := setup(t)
	f.toReview(t)
	over := f.readyStory(t, "Over", t0, "flai/internal")
	f.readyStory(t, "Behind", t0.Add(time.Second))
	out := answered(t, f.held(t, 1))
	if out["reason"] != "pull" || storyOf(out) != over.ID {
		t.Fatalf("the story overlapping one in review: %v", out)
	}
	for _, c := range out["ready"].([]any) {
		if card := c.(map[string]any); card["held"] != nil {
			t.Errorf("%v is held with only a story in review open", card["id"])
		}
	}
}

// S-0130: a ready story that names, in after:, a story not yet done is
// passed over, and the ready list says what it waits for.
func TestWaitForWorkPassesOverAStoryWaitingForAnother(t *testing.T) {
	f := setup(t)
	f.toReview(t)
	waits := f.readyStory(t, "Waits", t0)
	waits.After = []string{f.story.ID}
	if err := f.repo.Save(waits); err != nil {
		t.Fatal(err)
	}
	clear := f.readyStory(t, "Clear", t0.Add(time.Second))
	out := answered(t, f.held(t, 1))
	if out["reason"] != "pull" || storyOf(out) != clear.ID {
		t.Fatalf("the clear story: %v", out)
	}
	want := "held (after): waits for " + f.story.ID + " (in review); starts when " + f.story.ID + " is done"
	var held map[string]any
	for _, c := range out["ready"].([]any) {
		if card := c.(map[string]any); card["id"] == waits.ID {
			held, _ = card["held"].(map[string]any)
		}
	}
	if held["code"] != "after" || held["reason"] != want {
		t.Errorf("held: %v in %v", held, out["ready"])
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

// S-0220: a recommendation on a thread the agent opened is no answer to it
// until the operator confirms it (ADR-0090): the thread awaits someone else
// in its inbox, and wait_for_work is not woken by it; the confirmation is the
// answer, and it is.
func TestARecommendationIsNoAnswerToTheAgentThatAskedUntilConfirmed(t *testing.T) {
	f := setup(t)
	f.toReview(t)
	th, err := threads.New(f.repo, threads.NewOptions{Title: "Which port?", On: f.story.ID, Author: "claude", Text: "Eight or nine?", Now: t0})
	if err != nil {
		t.Fatal(err)
	}
	if quiet, _ := f.call(t, "wait_for_work", map[string]any{"timeout_seconds": 1}); quiet["reason"] != "" {
		t.Fatalf("nothing to do yet: %v", quiet)
	}
	done := f.held(t, 3)
	time.Sleep(150 * time.Millisecond)
	source := threads.Source{Path: "design/system/plan.md", Heading: "Shape"}
	if _, err := threads.ReplyWith(f.repo, th.ID, workitem.ActivityOrchestrator, "Nine.", t0.Add(2*time.Minute), threads.Marks{Recommendation: true, Source: source}); err != nil {
		t.Fatal(err)
	}
	select {
	case out := <-done:
		t.Fatalf("a recommendation woke the agent that asked: %v", out)
	case <-time.After(600 * time.Millisecond):
	}
	out, _ := f.call(t, "inbox", map[string]any{"all": true})
	sum := out["threads"].([]any)[0].(map[string]any)
	pending, _ := sum["pending_recommendation"].(map[string]any)
	if out["awaiting_you"].(float64) != 0 || sum["awaiting"] != "other" || sum["last_by"] != workitem.ActivityOrchestrator || pending == nil || pending["recommendation"] != true {
		t.Fatalf("a pending recommendation awaits the operator, not the agent that asked: %v", out)
	}
	if out, _ := f.call(t, "inbox", map[string]any{}); len(out["threads"].([]any)) != 0 {
		t.Errorf("a thread awaiting the operator's confirmation is not in the asker's default inbox: %v", out)
	}
	if _, err := threads.Confirm(f.repo, th.ID, "alex", t0.Add(3*time.Minute)); err != nil {
		t.Fatal(err)
	}
	woke := answered(t, done)
	if ths := woke["threads"].([]any); woke["reason"] != "thread" || len(ths) != 1 || ths[0].(map[string]any)["id"] != th.ID {
		t.Errorf("the confirmation wakes the agent that asked: %v", woke)
	}
	out, _ = f.call(t, "inbox", map[string]any{})
	if out["awaiting_you"].(float64) != 1 {
		t.Fatalf("the confirmed recommendation is the answer: %v", out)
	}
	if sum := out["threads"].([]any)[0].(map[string]any); sum["awaiting"] != "you" || sum["last_by"] != "alex" || sum["pending_recommendation"] != nil {
		t.Errorf("confirmed: %v", sum)
	}
}

// S-0220: a plain answer citing a source, as the orchestrator gives with
// answer_threads autonomous, is an answer: the agent that asked goes on.
func TestAnAnswerWithASourceIsAnAnswerToTheAgentThatAsked(t *testing.T) {
	f := setup(t)
	th, err := threads.New(f.repo, threads.NewOptions{Title: "Which port?", On: f.story.ID, Author: "claude", Text: "Eight or nine?", Now: t0})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := threads.ReplyWith(f.repo, th.ID, workitem.ActivityOrchestrator, "Nine.", t0.Add(time.Minute), threads.Marks{Source: threads.Source{Path: "design/system/plan.md"}}); err != nil {
		t.Fatal(err)
	}
	out, _ := f.call(t, "inbox", map[string]any{})
	if out["awaiting_you"].(float64) != 1 {
		t.Fatalf("an answer awaits the agent that asked: %v", out)
	}
	if sum := out["threads"].([]any)[0].(map[string]any); sum["awaiting"] != "you" || sum["status"] != "answered" || sum["pending_recommendation"] != nil {
		t.Errorf("answered: %v", sum)
	}
}
