package inbox

import (
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/threads"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

const promotableBody = "## Goal\n\nDo it.\n\n## Acceptance criteria\n- [ ] it works\n\n## Notes\n"

// backlogStory is a backlog story under E-0001 that is a promotion candidate
// unless the test changes it.
func backlogStory(id, title string, after ...string) *workitem.Item {
	value := 100.0
	return &workitem.Item{
		ID: id, Type: workitem.Story, Status: workitem.Backlog, Parent: "E-0001", Title: title,
		Created: "2026-11-01T10:00:00Z", Body: promotableBody, Touches: []string{"docs/" + id}, After: after,
		CostOfDelay: &workitem.CostOfDelay{Value: &value}, Forecast: &workitem.Forecast{Duration: "2h"},
	}
}

// starvedOf judges items in order, as PromotionCandidates does, and says what
// each refused story waits on.
func starvedOf(t *testing.T, items []*workitem.Item, order []string, sums []ThreadSummary) Starved {
	t.Helper()
	c, err := workitem.PromotionCandidates(items, order, workitem.NewHolds(items, nil), workitem.PolicyFIFO, 0)
	if err != nil {
		t.Fatal(err)
	}
	return StarvedFrom(items, c, sums)
}

// I-0288: from 2026-11-03 to 2026-11-05 every backlog story was a draft or
// held by after behind one, and the drafts waited on threads awaiting the
// operator. The list names both drafts with their threads, a resolved thread
// and one awaiting another left out, and gives each of the five held stories
// the draft at the root of its chain, a done story on the chain counting as
// done.
func TestTheStarvedBoardOfI0288NamesEachDecision(t *testing.T) {
	epic := &workitem.Item{ID: "E-0001", Type: workitem.Epic, Status: workitem.InProgress, Title: "Epic"}
	shipped := backlogStory("S-0400", "Shipped")
	shipped.Status = workitem.Done
	first := backlogStory("S-0406", "Count the turns")
	first.Draft = true
	second := backlogStory("S-0412", "Price the tokens")
	second.Draft = true
	a := backlogStory("S-0413", "Chart the turns", "S-0406")
	b := backlogStory("S-0414", "Alert on turns", "S-0413")
	c := backlogStory("S-0415", "Cap the turns", "S-0414")
	d := backlogStory("S-0416", "Report the cap", "S-0415")
	e := backlogStory("S-0417", "Tune the cap", "S-0400", "S-0416")
	items := []*workitem.Item{epic, shipped, first, second, a, b, c, d, e}
	order := []string{"S-0406", "S-0413", "S-0414", "S-0412", "S-0415", "S-0416", "S-0417"}
	sums := []ThreadSummary{
		{ID: "TH-0530", Title: "Settled already", Status: "resolved", Story: "S-0406", Awaiting: "operator"},
		{ID: "TH-0531", Title: "Asked the agent", Status: "open", Story: "S-0406", Awaiting: "other"},
		{ID: "TH-0532", Title: "Which turns count?", Status: "open", Story: "S-406", Awaiting: "operator"},
		{ID: "TH-0533", Title: "Which price list", Status: "answered", Story: "S-0412", Awaiting: "operator"},
	}

	got := starvedOf(t, items, order, sums)
	if len(got.Candidates) != 0 {
		t.Fatalf("no story is a candidate: %+v", got.Candidates)
	}
	want := "- S-0406 Count the turns: a draft; finalize it. It waits on TH-0532 Which turns count?\n" +
		"- S-0413 Chart the turns: held by `after`; the root of its chain is S-0406 Count the turns (backlog).\n" +
		"- S-0414 Alert on turns: held by `after`; the root of its chain is S-0406 Count the turns (backlog).\n" +
		"- S-0412 Price the tokens: a draft; finalize it. It waits on TH-0533 Which price list.\n" +
		"- S-0415 Cap the turns: held by `after`; the root of its chain is S-0406 Count the turns (backlog).\n" +
		"- S-0416 Report the cap: held by `after`; the root of its chain is S-0406 Count the turns (backlog).\n" +
		"- S-0417 Tune the cap: held by `after`; the root of its chain is S-0406 Count the turns (backlog).\n"
	if text := got.Text(); text != want {
		t.Errorf("text:\n%s\nwant:\n%s", text, want)
	}
	for _, st := range got.Stories {
		switch st.ID {
		case "S-0406", "S-0412":
			if st.Waits != WaitsDraft || len(st.Threads) != 1 || len(st.Reasons) != 1 || st.Reasons[0] != "draft: finalize it first" {
				t.Errorf("%s waits on its draft and one thread, its reasons as the candidates give them: %+v", st.ID, st)
			}
		default:
			if st.Waits != WaitsAfter || st.Loops || !reflect.DeepEqual(st.Roots, []Ref{{ID: "S-0406", Title: "Count the turns", Status: workitem.Backlog}}) || !strings.HasPrefix(st.Held, "held (after)") {
				t.Errorf("%s waits on S-0406 at the root of its chain: %+v", st.ID, st)
			}
		}
	}
	if again := starvedOf(t, items, order, sums).Text(); again != got.Text() {
		t.Errorf("the same list gives the same text:\n%s\n%s", again, got.Text())
	}
}

// A draft with no thread still waits on finalizing; a blocked story waits on
// its block reason and the thread it names; a chain that forks has a root on
// each branch, one in progress; one that loops says so; one held behind a
// cancelled story says its hold; and any other gives its reasons as they are.
func TestEachStarvedDecisionInWords(t *testing.T) {
	epic := &workitem.Item{ID: "E-0001", Type: workitem.Epic, Status: workitem.InProgress, Title: "Epic"}
	lone := backlogStory("S-0501", "Lone draft")
	lone.Draft = true
	blocked := backlogStory("S-0502", "Blocked one")
	blocked.Forecast = nil
	blocked.Blocked = []workitem.Block{
		{From: "2026-11-01T10:00:00Z", Until: "2026-11-02T10:00:00Z", Reason: "an old block"},
		{From: "2026-11-03T10:00:00Z", Reason: "waits on th-540 and TH-0541"},
	}
	working := backlogStory("S-0503", "Working")
	working.Status = workitem.InProgress
	forks := backlogStory("S-0504", "Forks", "S-0501", "S-0503")
	loopA := backlogStory("S-0505", "Loop A", "S-0506")
	loopB := backlogStory("S-0506", "Loop B", "S-0505")
	gone := backlogStory("S-0507", "Gone")
	gone.Status = workitem.Cancelled
	behindGone := backlogStory("S-0508", "Behind gone", "S-0507")
	unpriced := backlogStory("S-0509", "Unpriced")
	unpriced.CostOfDelay, unpriced.Forecast = nil, nil
	ready := backlogStory("S-0510", "Ready to promote")
	items := []*workitem.Item{epic, lone, blocked, working, forks, loopA, loopB, gone, behindGone, unpriced, ready}
	sums := []ThreadSummary{{ID: "TH-0540", Title: "Unblock it?", Status: "open", Awaiting: "operator"}}

	got := starvedOf(t, items, nil, sums)
	if len(got.Candidates) != 1 || got.Candidates[0].ID != "S-0510" {
		t.Errorf("the candidates come through: %+v", got.Candidates)
	}
	want := "- S-0501 Lone draft: a draft; finalize it.\n" +
		"- S-0502 Blocked one: blocked: waits on th-540 and TH-0541. It waits on TH-0540 Unblock it? and TH-0541.\n" +
		"- S-0504 Forks: held by `after`; the roots of its chain are S-0501 Lone draft (backlog) and S-0503 Working (in-progress).\n" +
		"- S-0505 Loop A: held by `after` in a chain that loops back on itself; fix `after`.\n" +
		"- S-0506 Loop B: held by `after` in a chain that loops back on itself; fix `after`.\n" +
		"- S-0508 Behind gone: held (after): waits for S-0507 (cancelled); starts when S-0507 is done; S-0507 was cancelled, so drop it from after: if S-0508 no longer needs it.\n" +
		"- S-0509 Unpriced: no forecast duration; no cost of delay value.\n"
	if text := got.Text(); text != want {
		t.Errorf("text:\n%s\nwant:\n%s", text, want)
	}
	if st := got.Stories[1]; st.Waits != WaitsBlocked || st.Block != "waits on th-540 and TH-0541" || !reflect.DeepEqual(st.Threads, []Ref{{ID: "TH-0540", Title: "Unblock it?"}, {ID: "TH-0541"}}) {
		t.Errorf("the blocked story waits on its open block and the threads it names: %+v", st)
	}
	if st := got.Stories[3]; st.Waits != WaitsAfter || !st.Loops || len(st.Roots) != 0 {
		t.Errorf("a loop has no root: %+v", st)
	}
	if st := got.Stories[6]; st.Waits != WaitsOther || len(st.Reasons) != 2 {
		t.Errorf("any other waits on its reasons: %+v", st)
	}
	if text := (Starved{}).Text(); text != "" {
		t.Errorf("no stories, no text: %q", text)
	}
}

// ReadStarved reads the backlog and the threads from a repository: a draft's
// thread awaiting the operator, on the story or on one of its tasks, is
// named, and a resolved one is left out.
func TestReadStarvedReadsTheBacklogAndTheThreads(t *testing.T) {
	f := setup(t)
	draft, err := f.repo.Create(workitem.NewOptions{Type: workitem.Story, Title: "Draft", Parent: f.epic.ID, Owner: "alex", Touches: []string{"docs/draft"}, Draft: true, Now: t0})
	if err != nil {
		t.Fatal(err)
	}
	task, err := f.repo.Create(workitem.NewOptions{Type: workitem.Task, Title: "Its task", Parent: draft.ID, Owner: "alex", Now: t0})
	if err != nil {
		t.Fatal(err)
	}
	planner := "planner-" + draft.ID
	asked := f.newThread(t, "Split it?", draft.ID, planner, "Operator, split it?")
	onTask := f.newThread(t, "Which file?", task.ID, planner, "Operator, which file?")
	settled := f.newThread(t, "Settled", draft.ID, planner, "Operator, settled?")
	if _, err := threads.Resolve(f.repo, settled.ID, "alex", "settled", t0.Add(time.Minute)); err != nil {
		t.Fatal(err)
	}
	held := f.storyIn(t, "Held", []string{"docs/held"})
	held.After = []string{draft.ID}
	if err := f.repo.Save(held); err != nil {
		t.Fatal(err)
	}
	items, err := f.repo.List(false)
	if err != nil {
		t.Fatal(err)
	}

	got, err := ReadStarved(f.repo, items, Reader{Agent: "flai serve"})
	if err != nil {
		t.Fatal(err)
	}
	if len(got.Candidates) != 0 || len(got.Stories) != 2 {
		t.Fatalf("two stories cannot be promoted, and none can: %+v", got)
	}
	byID := map[string]StarvedStory{}
	for _, st := range got.Stories {
		byID[st.ID] = st
	}
	if st := byID[draft.ID]; st.Waits != WaitsDraft || !reflect.DeepEqual(st.Threads, []Ref{{ID: asked.ID, Title: asked.Title}, {ID: onTask.ID, Title: onTask.Title}}) {
		t.Errorf("the draft names its open threads, on it and its task, the resolved one left out: %+v", st)
	}
	if st := byID[held.ID]; st.Waits != WaitsAfter || len(st.Roots) != 1 || st.Roots[0].ID != draft.ID {
		t.Errorf("the held story waits on the draft: %+v", st)
	}
}
