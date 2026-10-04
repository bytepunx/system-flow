//go:build !windows

package serve

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/cron"
	"github.com/bytepunx/system-flow/flai/internal/docedit"
	"github.com/bytepunx/system-flow/flai/internal/execx"
	"github.com/bytepunx/system-flow/flai/internal/itemedit"
	"github.com/bytepunx/system-flow/flai/internal/planning"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// S-0211, ADR-0084: an edit of a planned backlog or ready story's goal,
// criteria, or touches by someone other than a planner, after its forecast
// or its cost of delay value, is a trigger, naming the fields and who; so is
// an edit of its cost of delay's inputs after the value; nothing else is.
func TestAnEditOfAPlannedStoryIsATrigger(t *testing.T) {
	value := 1000.0
	story := func(change func(*workitem.Item)) *workitem.Item {
		it := &workitem.Item{ID: "S-0001", Type: workitem.Story, Status: workitem.Ready,
			Forecast: &workitem.Forecast{Duration: "2h", Delivery: "2026-10-05T00:00:00Z", By: "planner-S-0001", At: "2026-10-04T09:00:00Z"}}
		if change != nil {
			change(it)
		}
		return it
	}
	notice := func(by string, changed ...string) itemedit.Notice {
		return itemedit.Notice{At: "2026-10-04T10:00:00Z", By: by, ID: "S-0001", Type: workitem.Story, Changed: changed}
	}
	unplanned := func(it *workitem.Item) { it.Forecast = nil }
	for _, c := range []struct {
		name string
		it   *workitem.Item
		n    itemedit.Notice
		want string
	}{
		{"goal", story(nil), notice("alex", "goal"), "edited goal by alex"},
		{"criteria", story(nil), notice("alex", "criteria"), "edited criteria by alex"},
		{"touches", story(nil), notice("alex", "touches"), "edited touches by alex"},
		{"several, in the trigger's order", story(nil), notice("alex", "touches", "goal", "notes"), "edited goal, touches by alex"},
		{"no one named", story(nil), notice("", "goal"), "edited goal by someone"},
		{"in backlog", story(func(it *workitem.Item) { it.Status = workitem.Backlog }), notice("alex", "goal"), "edited goal by alex"},
		{"a value older than the edit, no forecast", story(func(it *workitem.Item) {
			it.Forecast, it.CostOfDelay = nil, &workitem.CostOfDelay{Value: &value, By: "alex", At: "2026-10-04T09:00:00Z"}
		}), notice("alex", "criteria"), "edited criteria by alex"},
		{"inputs changed after the value", story(func(it *workitem.Item) {
			it.CostOfDelay = &workitem.CostOfDelay{Inputs: &workitem.CostInputs{RevenuePerWeek: &value, By: "alex", At: "2026-10-04T10:00:00Z"}, Value: &value, By: "planner-S-0001", At: "2026-10-04T09:00:00Z"}
		}), notice("alex", "cost_of_delay"), "edited cost_of_delay by alex"},
		{"a value written alone", story(func(it *workitem.Item) {
			it.CostOfDelay = &workitem.CostOfDelay{Inputs: &workitem.CostInputs{RevenuePerWeek: &value, By: "alex", At: "2026-10-04T08:00:00Z"}, Value: &value, By: "alex", At: "2026-10-04T10:00:00Z"}
		}), notice("alex", "cost_of_delay"), ""},
		{"other fields", story(nil), notice("alex", "title", "notes", "body", "forecast"), ""},
		{"the planner's own edit", story(nil), notice("planner-S-0001", "goal"), ""},
		{"never planned", story(unplanned), notice("alex", "goal"), ""},
		{"never planned, inputs only", story(func(it *workitem.Item) {
			it.Forecast, it.CostOfDelay = nil, &workitem.CostOfDelay{Inputs: &workitem.CostInputs{RevenuePerWeek: &value, By: "alex", At: "2026-10-04T10:00:00Z"}}
		}), notice("alex", "cost_of_delay"), ""},
		{"a forecast newer than the edit", story(func(it *workitem.Item) { it.Forecast.At = "2026-10-04T11:00:00Z" }), notice("alex", "goal"), ""},
		{"a forecast in the edit's second", story(func(it *workitem.Item) { it.Forecast.At = "2026-10-04T10:00:00Z" }), notice("alex", "goal"), ""},
		{"a forecast with no stamp", story(func(it *workitem.Item) { it.Forecast.At = "" }), notice("alex", "goal"), ""},
		{"in progress", story(func(it *workitem.Item) { it.Status = workitem.InProgress }), notice("alex", "goal"), ""},
		{"archived", story(func(it *workitem.Item) { it.Archived = true }), notice("alex", "goal"), ""},
		{"a task", story(func(it *workitem.Item) { it.Type = workitem.Task }), notice("alex", "touches"), ""},
		{"an item not there", nil, notice("alex", "goal"), ""},
	} {
		if got := editTrigger(c.it, c.n); got != c.want {
			t.Errorf("%s: trigger %q, want %q", c.name, got, c.want)
		}
	}
}

// S-0211: a story accepted or cancelled is work ahead completing; a task's
// move, an epic's, and a story's move elsewhere are not.
func TestAStoryAcceptedOrCancelledIsATrigger(t *testing.T) {
	changes := []workitem.Change{
		{ID: "S-0001", Type: workitem.Story, Kind: workitem.Moved, To: workitem.Done},
		{ID: "T-0001", Type: workitem.Task, Kind: workitem.Moved, To: workitem.Done},
		{ID: "S-0002", Type: workitem.Story, Kind: workitem.Moved, To: workitem.Ready},
		{ID: "E-0001", Type: workitem.Epic, Kind: workitem.Moved, To: workitem.Done},
		{ID: "S-0003", Type: workitem.Story, Kind: workitem.WasBlocked},
		{ID: "S-0004", Type: workitem.Story, Kind: workitem.Moved, To: workitem.Cancelled},
	}
	if got := moveTriggers(changes); !slices.Equal(got, []string{"accepted S-0001", "cancelled S-0004"}) {
		t.Errorf("triggers = %v", got)
	}
}

// S-0211: a move or an edit is news once: from the replanner's start, those
// in its very second included, and not again at a later look, while one
// stamped in the second of the last look and not yet seen still is.
func TestChangesAndEditsAreNewsOnce(t *testing.T) {
	start := time.Date(2026, 10, 4, 10, 0, 0, 0, time.UTC)
	stamp := func(d time.Duration) string { return start.Add(d).Format(workitem.TimeFormat) }
	items := []*workitem.Item{{ID: "S-0001", Type: workitem.Story, Transitions: []workitem.Transition{
		{To: workitem.Ready, By: "alex", At: stamp(-time.Second)},
		{To: workitem.InProgress, By: "agent", At: stamp(0)},
		{To: workitem.Review, By: "agent", At: stamp(5 * time.Second)},
	}}}
	notices := []itemedit.Notice{{At: stamp(-time.Second), By: "alex", ID: "S-0001", Changed: []string{"goal"}}, {At: stamp(5 * time.Second), By: "alex", ID: "S-0001", Changed: []string{"touches"}}}
	changes, edits, since, seen := news(items, notices, start, map[string]time.Time{}, start.Add(5*time.Second+300*time.Millisecond))
	if len(changes) != 2 || changes[0].To != workitem.InProgress || changes[1].To != workitem.Review || len(edits) != 1 || edits[0].Changed[0] != "touches" {
		t.Fatalf("first look: %+v %+v", changes, edits)
	}
	// the same second again: nothing new, until a move stamped in it appears
	if changes, edits, _, _ := news(items, notices, since, seen, start.Add(5*time.Second+600*time.Millisecond)); len(changes) != 0 || len(edits) != 0 {
		t.Errorf("seen again: %+v %+v", changes, edits)
	}
	items[0].Transitions = append(items[0].Transitions, workitem.Transition{To: workitem.Done, By: "alex", At: stamp(5 * time.Second)})
	changes, _, since, seen = news(items, notices, since, seen, start.Add(6*time.Second))
	if len(changes) != 1 || changes[0].To != workitem.Done {
		t.Fatalf("a move in the last look's second: %+v", changes)
	}
	if changes, edits, _, _ := news(items, notices, since, seen, start.Add(time.Minute)); len(changes) != 0 || len(edits) != 0 {
		t.Errorf("a later look: %+v %+v", changes, edits)
	}
}

// S-0211: the schedule comes round when its next time after the mark has
// passed, and several missed times come round once.
func TestTheScheduleComesRoundOnce(t *testing.T) {
	daily, err := cron.Parse("daily")
	if err != nil {
		t.Fatal(err)
	}
	mark := time.Date(2026, 10, 4, 10, 0, 0, 0, time.UTC)
	if scheduleDue(daily, mark, time.Date(2026, 10, 4, 23, 59, 59, 0, time.UTC)) {
		t.Error("due before midnight")
	}
	if !scheduleDue(daily, mark, time.Date(2026, 10, 5, 0, 0, 0, 0, time.UTC)) {
		t.Error("not due at midnight")
	}
	late := time.Date(2026, 10, 8, 12, 0, 0, 0, time.UTC) // four midnights missed
	if !scheduleDue(daily, mark, late) {
		t.Error("not due after several missed times")
	}
	if scheduleDue(daily, late, late.Add(time.Hour)) {
		t.Error("due again once it came round")
	}
}

// S-0211: the queue holds a story once, with each trigger once, in the order
// they came; a planner is started for the first entry whose item has none
// running.
func TestTheQueueCoalescesAndPicksAStoryWithNoPlannerRunning(t *testing.T) {
	var q []queued
	for _, e := range [][2]string{{"S-0001", "edited goal by alex"}, {"S-0002", "reordered"}, {"S-0001", "accepted S-0009"}, {"S-0001", "edited goal by alex"}} {
		q = enqueue(q, e[0], e[1])
	}
	if len(q) != 2 || q[0].item != "S-0001" || !slices.Equal(q[0].triggers, []string{"edited goal by alex", "accepted S-0009"}) || q[1].item != "S-0002" {
		t.Fatalf("queue = %+v", q)
	}
	running := func(ids ...string) func(string) bool { return func(id string) bool { return slices.Contains(ids, id) } }
	if i := pick(q, running()); i != 0 {
		t.Errorf("pick = %d, want the first", i)
	}
	if i := pick(q, running("S-0001")); i != 1 {
		t.Errorf("pick = %d, want the second while the first's planner runs", i)
	}
	if i := pick(q, running("S-0001", "S-0002")); i != -1 {
		t.Errorf("pick = %d, want none", i)
	}
	if i := pick(nil, running()); i != -1 {
		t.Errorf("pick of nothing = %d", i)
	}
}

// replanLab is a planLab with a replanner whose clock the test sets, and
// whose git is recorded.
type replanLab struct {
	*agentLab
	r     *replanner
	clock time.Time
	git   *recordedGit
}

// recordedGit records the git commands run, and answers each as done.
type recordedGit struct {
	execx.System
	mu    sync.Mutex
	calls [][]string
}

func (g *recordedGit) Run(dir, name string, args ...string) (string, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.calls = append(g.calls, append([]string{name}, args...))
	return "true", nil
}

// commits are the git commit commands run.
func (g *recordedGit) commits() [][]string {
	g.mu.Lock()
	defer g.mu.Unlock()
	var out [][]string
	for _, c := range g.calls {
		if len(c) > 1 && c[1] == "commit" {
			out = append(out, c)
		}
	}
	return out
}

func newReplanLab(t *testing.T, planningYAML string) *replanLab {
	t.Helper()
	lab := &replanLab{agentLab: planLab(t), git: &recordedGit{}}
	if planningYAML != "" {
		if err := os.WriteFile(filepath.Join(lab.root, "system-flow.yaml"), []byte("version: 1\nname: t\nkey: t\nlayout:\n  design: design\n  docs: docs\n  wip: wip\nplanning:\n"+planningYAML), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	lab.clock = time.Now().UTC()
	o := lab.o
	o.Now = func() time.Time { return lab.clock }
	o.Git = lab.git
	lab.r = newReplanner(o, Entry{Key: "t", Name: "t", Root: lab.root}, lab.l)
	return lab
}

// edit edits story as by would from the dashboard, at at.
func (lab *replanLab) edit(id, by string, at time.Time, ch itemedit.Change) {
	lab.t.Helper()
	if _, err := itemedit.Apply(lab.repo, lab.git, id, ch, itemedit.Options{By: by, NoCommit: true, Now: at}); err != nil {
		var no *docedit.RefusedError
		if errors.As(err, &no) {
			lab.t.Fatalf("edit %s: %v %+v", id, err, no.Findings)
		}
		lab.t.Fatalf("edit %s: %v", id, err)
	}
}

// forecast gives story a forecast as its planner would, at the lab's clock.
func (lab *replanLab) forecast(id, duration, delivery string) {
	lab.t.Helper()
	basis := "From history."
	lab.edit(id, "planner-"+id, lab.clock, itemedit.Change{Forecast: &itemedit.ForecastEdit{Duration: &duration, Delivery: &delivery, Basis: &basis}})
}

// accept moves a story through to done at the lab's clock.
func (lab *replanLab) accept(id string) {
	lab.t.Helper()
	task, err := lab.repo.Create(workitem.NewOptions{Type: workitem.Task, Title: "Work", Parent: id, Owner: "builder-" + id, Now: lab.clock})
	if err != nil {
		lab.t.Fatal(err)
	}
	for _, m := range [][2]string{{id, workitem.InProgress}, {task.ID, workitem.Ready}, {task.ID, workitem.InProgress}, {task.ID, workitem.Done}, {id, workitem.Review}, {id, workitem.Done}} {
		it := lab.item(m[0])
		if m == [2]string{id, workitem.Done} { // its criterion met
			data, _ := os.ReadFile(it.Path)
			if err := os.WriteFile(it.Path, []byte(strings.Replace(string(data), "- [ ] works", "- [x] works", 1)), 0o644); err != nil {
				lab.t.Fatal(err)
			}
			it = lab.item(id)
		}
		if _, err := lab.repo.Transition(it, m[1], "alex", "", lab.clock); err != nil {
			lab.t.Fatalf("move %s to %s: %v", m[0], m[1], err)
		}
	}
}

func (lab *replanLab) item(id string) *workitem.Item {
	lab.t.Helper()
	it, err := lab.repo.Get(id)
	if err != nil {
		lab.t.Fatal(err)
	}
	return it
}

// newGoal is a change of story's goal.
func (lab *replanLab) newGoal(id string) itemedit.Change {
	lab.t.Helper()
	v, err := itemedit.Show(lab.repo, id)
	if err != nil {
		lab.t.Fatal(err)
	}
	body := strings.Replace(v.Body, "## Goal\n", "## Goal\n\nShip it sooner.\n", 1)
	if body == v.Body {
		lab.t.Fatalf("%s has no goal section:\n%s", id, v.Body)
	}
	return itemedit.Change{Body: &body}
}

// S-0211, ADR-0084: with plan on, an edit of a planned ready story's goal
// starts the planner for it, and the run records what started it.
func TestAnEditOfAPlannedStoryStartsThePlanner(t *testing.T) {
	lab := newReplanLab(t, "")
	ctx := context.Background()
	id := lab.ready("Planned")
	unplanned := lab.ready("Unplanned")
	lab.forecast(id, "2h", "2026-12-01T00:00:00Z")
	stream := strings.Join([]string{streamCall("s", "m1", runStart, 999), streamFinal("Replanned "+id+".", 1000, 0.5)}, "\n") + "\n"
	if err := os.WriteFile(filepath.Join(lab.outDir, "stream-"+id), []byte(stream), 0o644); err != nil {
		t.Fatal(err)
	}
	lab.r.look(ctx)
	if run := lab.planRun(id); run != nil {
		t.Fatalf("the planner's own forecast started a planner: %+v", run)
	}
	lab.clock = lab.clock.Add(time.Minute)
	lab.edit(id, "alex", lab.clock, lab.newGoal(id))
	lab.edit(unplanned, "alex", lab.clock, lab.newGoal(unplanned))
	lab.r.look(ctx)
	waitFor(t, "the planner started", func() bool { return lab.planRun(id) != nil })
	if run := lab.planRun(id); run.Trigger != "edited goal by alex" || run.Error != "" {
		t.Errorf("run = %+v, want started for the edit", run)
	}
	if run := lab.planRun(unplanned); run != nil {
		t.Errorf("a story never planned was planned on an edit: %+v", run)
	}
	waitFor(t, "the planner ended", func() bool { r := lab.planRun(id); return r.Ended != "" })
	waitFor(t, "its activity logged with the trigger", func() bool {
		doc, err := lab.repo.Activity(workitem.ActivityPlanner)
		return err == nil && len(doc.Entries) == 1 && doc.Entries[0].Trigger == "edited goal by alex"
	})
	// a later look starts nothing again
	lab.clock = lab.clock.Add(time.Minute)
	started := lab.planRun(id).Started
	lab.r.look(ctx)
	if run := lab.planRun(id); run.Started != started {
		t.Errorf("the edit started the planner twice: %+v", run)
	}
}

// S-0211: with plan off, an edit starts nothing, and turning plan on does
// not act on it afterwards.
func TestNothingIsPlannedOnItsOwnWhilePlanIsOff(t *testing.T) {
	lab := newReplanLab(t, "")
	ctx := context.Background()
	id := lab.ready("Planned")
	lab.forecast(id, "2h", "2026-12-01T00:00:00Z")
	lab.cfg.Plan = false
	lab.clock = lab.clock.Add(time.Minute)
	lab.edit(id, "alex", lab.clock, lab.newGoal(id))
	lab.r.look(ctx)
	lab.cfg.Plan = true
	lab.clock = lab.clock.Add(time.Minute)
	lab.r.look(ctx)
	if run := lab.planRun(id); run != nil || len(lab.r.queue) != 0 {
		t.Errorf("planned while off: %+v, queue %+v", run, lab.r.queue)
	}
}

// replayed is what story's forecast replays as now, from the files.
func (lab *replanLab) replayed(id string) planning.Result {
	lab.t.Helper()
	repo, err := workitem.Open(lab.root)
	if err != nil {
		lab.t.Fatal(err)
	}
	items, err := repo.List(true)
	if err != nil {
		lab.t.Fatal(err)
	}
	board, err := repo.LoadBoard()
	if err != nil {
		lab.t.Fatal(err)
	}
	res, err := planning.Replay(items, board, repo.Manifest.Planning, id, lab.clock)
	if err != nil {
		lab.t.Fatal(err)
	}
	return res
}

// S-0211, ADR-0084: under deterministic, the default, a story ahead accepted
// moves the forecast delivery of the stories behind it, keeping their
// durations, stamped by flai, in one commit, and starts no planner.
func TestAStoryAcceptedReplansTheForecastsBehindIt(t *testing.T) {
	lab := newReplanLab(t, "")
	ctx := context.Background()
	lab.limit(1)
	ahead, later, backlog := lab.ready("Ahead"), lab.ready("Later"), lab.backlog("Behind", nil)
	stale := "2026-01-01T00:00:00Z"
	lab.forecast(ahead, "8h", stale)
	lab.forecast(later, "1h", stale)
	lab.forecast(backlog, "3h", stale)
	lab.accept(ahead)
	lab.clock = lab.clock.Add(time.Minute)
	lab.r.look(ctx)
	for _, id := range []string{later, backlog} {
		it, want := lab.item(id), lab.replayed(id)
		if f := it.Forecast; f.Delivery != want.Delivery || f.Basis != want.Basis || f.By != "flai" || f.Delivery == stale || (id == later && f.Duration != "1h") {
			t.Errorf("%s forecast = %+v, want delivery %s and basis %q by flai", id, f, want.Delivery, want.Basis)
		}
	}
	commits := lab.git.commits()
	if len(commits) != 1 {
		t.Fatalf("commits = %v, want one", commits)
	}
	c := strings.Join(commits[0], " ")
	if !strings.Contains(c, "-m chore: replan forecasts after accepted "+ahead+" --") || !strings.Contains(c, filepath.Base(lab.item(later).Path)) || !strings.Contains(c, filepath.Base(lab.item(backlog).Path)) {
		t.Errorf("commit = %s", c)
	}
	if st := lab.state(); len(st.Plans) != 0 {
		t.Errorf("a planner was started under deterministic: %+v", st.Plans)
	}
	if !strings.Contains(lab.logText(), `msg="forecasts replanned" component=serve project=t policy=deterministic moved=2 triggers="accepted `+ahead+`"`) {
		t.Errorf("the replan is not logged:\n%s", lab.logText())
	}
}

// S-0211: a change of the pull order moves the forecasts it changes, in a
// commit that says so.
func TestAReorderingReplansTheForecasts(t *testing.T) {
	lab := newReplanLab(t, "")
	ctx := context.Background()
	lab.limit(1)
	first, second := lab.ready("First"), lab.ready("Second")
	lab.forecast(first, "2h", "2026-01-01T00:00:00Z")
	lab.forecast(second, "3h", "2026-01-01T00:00:00Z")
	lab.board(first, second)
	lab.r.look(ctx)
	before := lab.item(first).Forecast.Delivery
	commits := len(lab.git.commits())
	lab.board(second, first)
	lab.clock = lab.clock.Add(time.Minute)
	lab.r.look(ctx)
	if got := lab.git.commits(); len(got) != commits+1 || !strings.Contains(strings.Join(got[len(got)-1], " "), "-m chore: replan forecasts after reordered --") {
		t.Errorf("commits = %v, want one after reordered", got)
	}
	if f := lab.item(first).Forecast; f.Delivery == before || f.Delivery != lab.replayed(first).Delivery {
		t.Errorf("%s forecast = %+v, want it behind %s", first, f, second)
	}
}

// board sets the pull order.
func (lab *replanLab) board(order ...string) {
	lab.t.Helper()
	b, err := lab.repo.LoadBoard()
	if err != nil {
		lab.t.Fatal(err)
	}
	b.Order = order
	if err := b.Save(lab.clock.Format("2006-01-02")); err != nil {
		lab.t.Fatal(err)
	}
}

// S-0211: under never, a story accepted moves no forecast.
func TestReplanNeverMovesNoForecast(t *testing.T) {
	lab := newReplanLab(t, "  replan: never\n")
	ctx := context.Background()
	ahead, later := lab.ready("Ahead"), lab.ready("Later")
	stale := "2026-01-01T00:00:00Z"
	lab.forecast(later, "1h", stale)
	lab.accept(ahead)
	lab.clock = lab.clock.Add(time.Minute)
	lab.r.look(ctx)
	if f := lab.item(later).Forecast; f.Delivery != stale || len(lab.git.calls) != 0 {
		t.Errorf("forecast = %+v, git %v: want it left", f, lab.git.calls)
	}
}

// S-0211: under agent, a story accepted moves the forecasts behind it and
// queues the planner for each story moved, one run at a time.
func TestReplanAgentQueuesThePlannerForEachStoryMoved(t *testing.T) {
	lab := newReplanLab(t, "  replan: agent\n")
	ctx := context.Background()
	lab.hold()
	ahead, later, behind := lab.ready("Ahead"), lab.ready("Later"), lab.ready("Behind")
	stale := "2026-01-01T00:00:00Z"
	lab.forecast(later, "1h", stale)
	lab.forecast(behind, "1h", stale)
	lab.accept(ahead)
	lab.clock = lab.clock.Add(time.Minute)
	lab.r.look(ctx)
	if f := lab.item(later).Forecast; f.Delivery == stale || f.By != "flai" || len(lab.git.commits()) != 1 {
		t.Errorf("forecast = %+v, commits %v: want it moved and committed", f, lab.git.commits())
	}
	trigger := "accepted " + ahead
	waitFor(t, "the first planner started", func() bool { return lab.planRun(later) != nil })
	if run := lab.planRun(later); run.Trigger != trigger {
		t.Errorf("run = %+v, want %q", run, trigger)
	}
	lab.r.look(ctx)
	if run := lab.planRun(behind); run != nil {
		t.Fatalf("a second planner started while the first runs: %+v", run)
	}
	lab.release(later)
	waitFor(t, "the first planner ended", func() bool { return lab.planRun(later).Ended != "" })
	lab.r.look(ctx)
	waitFor(t, "the second planner started", func() bool { return lab.planRun(behind) != nil })
	if run := lab.planRun(behind); run.Trigger != trigger {
		t.Errorf("run = %+v, want %q", run, trigger)
	}
	lab.release(behind)
}

// S-0211: the queue skips a story whose planner the operator started, keeps
// it, and drops an entry the checks refuse.
func TestTheQueueSkipsARunningPlannerAndDropsARefusal(t *testing.T) {
	lab := newReplanLab(t, "")
	ctx := context.Background()
	lab.hold()
	asked, next := lab.ready("Asked"), lab.ready("Next")
	done := lab.ready("Done")
	lab.accept(done)
	if !lab.planHere(asked) {
		t.Fatal("the operator's planner did not start")
	}
	lab.r.queue = []queued{{item: done, triggers: []string{"reordered"}}, {item: asked, triggers: []string{"reordered"}}, {item: next, triggers: []string{"reordered", "schedule daily"}}}
	lab.r.drain(ctx)
	waitFor(t, "the next planner started", func() bool { return lab.planRun(next) != nil })
	if run := lab.planRun(next); run.Trigger != "reordered; schedule daily" {
		t.Errorf("run = %+v", run)
	}
	if len(lab.r.queue) != 1 || lab.r.queue[0].item != asked {
		t.Errorf("queue = %+v, want only %s kept", lab.r.queue, asked)
	}
	if lab.planRun(done) != nil || !strings.Contains(lab.logText(), `msg="queued planner dropped" component=serve project=t item=`+done) {
		t.Errorf("the refused entry was not dropped and logged:\n%s", lab.logText())
	}
	lab.release(asked)
	lab.release(next)
}

// S-0211: planning.schedule queues the planner for every ready story when it
// comes round, once however many times were missed.
func TestTheScheduleQueuesTheReadyStories(t *testing.T) {
	lab := newReplanLab(t, "  schedule: daily\n")
	ctx := context.Background()
	lab.hold()
	first, second := lab.ready("First"), lab.ready("Second")
	lab.backlog("Waiting", nil)
	lab.r.look(ctx)
	if len(lab.r.queue) != 0 || lab.planRun(first) != nil {
		t.Fatalf("queued before the schedule came round: %+v", lab.r.queue)
	}
	lab.clock = lab.clock.Add(72 * time.Hour)
	lab.r.look(ctx)
	waitFor(t, "the first planner started", func() bool { return lab.planRun(first) != nil })
	if run := lab.planRun(first); run.Trigger != "schedule daily" {
		t.Errorf("run = %+v", run)
	}
	if len(lab.r.queue) != 1 || lab.r.queue[0].item != second {
		t.Errorf("queue = %+v, want %s alone left", lab.r.queue, second)
	}
	lab.release(first)
	lab.release(second)
}
