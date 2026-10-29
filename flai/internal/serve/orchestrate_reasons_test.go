//go:build !windows

package serve

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/cron"
	"github.com/bytepunx/system-flow/flai/internal/inbox"
	"github.com/bytepunx/system-flow/flai/internal/itemedit"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/roles"
	"github.com/bytepunx/system-flow/flai/internal/threads"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// reasonsGiven are the lines of FLAI_REASONS the orchestrator stub was given
// in the run with pid.
func (lab *orchestrateLab) reasonsGiven(pid int) []string {
	lab.t.Helper()
	given := lab.given(pid)
	_, after, ok := strings.Cut(given, "reasons: ")
	if !ok {
		lab.t.Fatalf("the stub wrote no reasons:\n%s", given)
	}
	reasons, _, _ := strings.Cut(after, "\nby: ")
	return strings.Split(reasons, "\n")
}

// history makes the lab's items from here on an hour old, so that what is
// made of them before a look is history to its catch-up.
func (lab *orchestrateLab) history() { lab.now = time.Now().Add(-time.Hour) }

// inProgress makes a story with a task and moves it to in progress, as
// history.
func (lab *orchestrateLab) inProgress(title string) (story, task string) {
	lab.t.Helper()
	story = lab.backlog(title, nil)
	lab.toReady(story)
	lab.moveAt(story, workitem.InProgress, lab.now)
	tk, err := lab.repo.Create(workitem.NewOptions{Type: workitem.Task, Title: "Work", Parent: story, Owner: "builder-" + story, Now: lab.now})
	if err != nil {
		lab.t.Fatal(err)
	}
	return story, tk.ID
}

// moveAt moves item id to state to at, as its agent would.
func (lab *orchestrateLab) moveAt(id, to string, at time.Time) {
	lab.t.Helper()
	it, err := lab.repo.Get(id)
	if err != nil {
		lab.t.Fatal(err)
	}
	if _, err := lab.repo.Transition(it, to, "builder", "", at); err != nil {
		lab.t.Fatal(err)
	}
}

// askOperator opens a thread on story by its agent at, which awaits the
// operator.
func (lab *orchestrateLab) askOperator(story, title string, at time.Time) string {
	lab.t.Helper()
	th, err := threads.New(lab.repo, threads.NewOptions{Title: title, On: story, Author: "agent-" + story, Text: "Which way?", Now: at})
	if err != nil {
		lab.t.Fatal(err)
	}
	return th.ID
}

// started are the journal's entries for the orchestrator's starts.
func (lab *orchestrateLab) started() []string {
	var out []string
	for _, e := range lab.orchestrateEntries() {
		if strings.HasPrefix(e.Detail, "started ") {
			out = append(out, e.Detail)
		}
	}
	return out
}

// firstRun has the lab's orchestrator turned on at base, its first run, for
// asked, started and ended.
func (lab *orchestrateLab) firstRun(base time.Time) *AgentRun {
	lab.t.Helper()
	lab.at = base
	lab.cfg.Orchestrate = true
	lab.look()
	first := lab.orchestrator()
	if first == nil {
		lab.t.Fatal("the orchestrator did not start when turned on")
	}
	return lab.ended(first.Session)
}

// S-0536 criterion 2, ADR-0209 §2: three events inside one coalescing
// window, a story moved to review, a thread come to await the operator, and
// a planner's run that wrote a story, start exactly one run, once the window
// has passed and not before, with all three reasons in its FLAI_REASONS and
// in its journal entry; the run takes them off the queue.
func TestThreeEventsInOneWindowStartOneRun(t *testing.T) {
	lab := newOrchestrateLab(t)
	lab.history()
	story, _ := lab.inProgress("Reviewed")
	base := time.Now().UTC().Truncate(time.Second)
	first := lab.firstRun(base)
	lab.moveAt(story, workitem.Review, base.Add(time.Second))
	thread := lab.askOperator(story, "Which way", base.Add(2*time.Second))
	planned := inbox.PlannedEvent(lab.epic.ID, workitem.Epic, "Epic", "planner", []string{story})
	planned.At = base.Add(3 * time.Second).Format(workitem.TimeFormat)
	if err := inbox.AppendEvent(lab.repo, planned); err != nil {
		t.Fatal(err)
	}
	queued := base.Add(5 * time.Second)
	for _, at := range []time.Time{queued, queued.Add(manifest.DefaultCoalesce - time.Second)} {
		lab.at = at
		lab.look()
		if r := lab.orchestrator(); r.Session != first.Session {
			t.Fatalf("started at %s, inside the window from %s: %+v", at, queued, r)
		}
	}
	if q := lab.o.Dir.reasonQueue(lab.root); len(q.Reasons) != 3 || q.First != queued.Format(time.RFC3339) {
		t.Fatalf("queue = %+v, want three reasons first queued at %s", q, queued)
	}
	lab.at = queued.Add(manifest.DefaultCoalesce)
	lab.look()
	run := lab.orchestrator()
	if run.Session == first.Session {
		t.Fatalf("not started once the window passed: %+v", run)
	}
	want := []string{
		"moved " + story + " to review: Reviewed",
		"thread " + thread + " awaits the operator: Which way",
		"planned " + lab.epic.ID + ", wrote " + story + ": Epic",
	}
	if got := lab.reasonsGiven(run.PID); strings.Join(got, "\n") != strings.Join(want, "\n") {
		t.Errorf("FLAI_REASONS =\n%s\nwant\n%s", strings.Join(got, "\n"), strings.Join(want, "\n"))
	}
	lab.ended(run.Session)
	starts := lab.started()
	if len(starts) != 2 || !strings.Contains(starts[1], " for 3 reasons: "+strings.Join(want, "; ")+"; log ") {
		t.Errorf("the journal's starts: %q, want the second for the three reasons", starts)
	}
	if q := lab.o.Dir.reasonQueue(lab.root); len(q.Reasons) != 0 || q.First != "" {
		t.Errorf("the run left its reasons queued: %+v", q)
	}
	lab.at = queued.Add(2 * manifest.DefaultCoalesce)
	lab.look()
	if r := lab.orchestrator(); r.Session != run.Session {
		t.Errorf("started again with nothing queued: %+v", r)
	}
	lab.logs.mu.Lock()
	logged := lab.logs.buf.String()
	lab.logs.mu.Unlock()
	if !strings.Contains(logged, `msg="orchestrator started"`) || !strings.Contains(logged, "reasons=3") {
		t.Errorf("the start was not logged with its count of reasons:\n%s", logged)
	}
}

// ADR-0209 §2: a reason that comes while a run goes waits for it, past its
// window, and starts the next run at the first look after the run ends.
func TestAReasonWhileARunGoesStartsTheNextRunWhenItEnds(t *testing.T) {
	lab := newOrchestrateLab(t)
	lab.history()
	story, _ := lab.inProgress("Reviewed")
	lab.hold()
	base := time.Now().UTC().Truncate(time.Second)
	lab.at = base
	lab.cfg.Orchestrate = true
	lab.look()
	first := lab.orchestrator()
	lab.moveAt(story, workitem.Review, base.Add(time.Second))
	for _, at := range []time.Time{base.Add(2 * time.Second), base.Add(2*time.Second + 2*manifest.DefaultCoalesce)} {
		lab.at = at
		lab.look()
		if r := lab.orchestrator(); r.Session != first.Session || !r.live() {
			t.Fatalf("at %s, a run started beside the one going: %+v", at, r)
		}
	}
	if q := lab.o.Dir.reasonQueue(lab.root); len(q.Reasons) != 1 || q.Reasons[0].Kind != workitem.Moved {
		t.Fatalf("queue = %+v, want the move", q)
	}
	lab.release("orchestrator")
	lab.ended(first.Session)
	if err := os.Remove(filepath.Join(lab.outDir, "release-orchestrator")); err != nil {
		t.Fatal(err)
	}
	lab.at = base.Add(3*time.Second + 2*manifest.DefaultCoalesce)
	lab.look()
	next := lab.orchestrator()
	if next.Session == first.Session || !next.live() {
		t.Fatalf("not started at the first look after the run ended: %+v", next)
	}
	if got := lab.reasonsGiven(next.PID); len(got) != 1 || got[0] != "moved "+story+" to review: Reviewed" {
		t.Errorf("FLAI_REASONS = %q, want the move", got)
	}
	lab.release("orchestrator")
	lab.ended(next.Session)
}

// ADR-0209 §1: a task's move, and the replan's forecast edit by flai, even
// of a draft, are no reasons and start nothing; another's edit of a draft
// is one.
func TestATaskMoveAndAForecastEditByFlaiStartNothing(t *testing.T) {
	lab := newOrchestrateLab(t)
	lab.history()
	_, task := lab.inProgress("Worked")
	draft, err := lab.repo.Create(workitem.NewOptions{Type: workitem.Story, Title: "Drafted", Parent: lab.epic.ID, Owner: "alex", Draft: true, Now: lab.now})
	if err != nil {
		t.Fatal(err)
	}
	base := time.Now().UTC().Truncate(time.Second)
	first := lab.firstRun(base)
	lab.moveAt(task, workitem.Ready, base.Add(time.Second))
	lab.moveAt(task, workitem.InProgress, base.Add(time.Second))
	edit := func(ch itemedit.Change, by string, at time.Time) {
		t.Helper()
		if _, err := itemedit.Apply(lab.repo, gitOf(lab.o), draft.ID, ch, itemedit.Options{By: by, NoCommit: true, Now: at}); err != nil {
			t.Fatal(err)
		}
	}
	duration, delivery, basis := "6h", base.Add(time.Hour).Format(workitem.TimeFormat), "three tasks like S-0185's"
	edit(itemedit.Change{Forecast: &itemedit.ForecastEdit{Duration: &duration, Delivery: &delivery, Basis: &basis}}, "flai", base.Add(time.Second))
	for _, at := range []time.Time{base.Add(2 * time.Second), base.Add(time.Hour)} {
		lab.at = at
		lab.look()
	}
	if r := lab.orchestrator(); r.Session != first.Session {
		t.Fatalf("started for a task's move or flai's forecast edit: %+v", r)
	}
	if q := lab.o.Dir.reasonQueue(lab.root); len(q.Reasons) != 0 {
		t.Fatalf("queued %+v, want nothing", q.Reasons)
	}
	touches := []string{"docs/drafted"}
	edit(itemedit.Change{Touches: &touches}, "planner", base.Add(time.Hour+time.Second))
	lab.at = base.Add(time.Hour + 2*time.Second)
	lab.look()
	if q := lab.o.Dir.reasonQueue(lab.root); len(q.Reasons) != 1 || q.Reasons[0].line() != "edited draft "+draft.ID+" (touches): Drafted" {
		t.Errorf("queued %+v, want the planner's edit of the draft", q.Reasons)
	}
}

// ADR-0209 §2: the first look with orchestrate on starts one run for asked,
// and none of the catch-up's history; a run that ends with nothing queued
// is not started again.
func TestTheFirstLookStartsOneRunForAskedAlone(t *testing.T) {
	lab := newOrchestrateLab(t)
	lab.history()
	story, _ := lab.inProgress("Reviewed")
	lab.moveAt(story, workitem.Review, lab.now)
	lab.hold()
	lab.cfg.Orchestrate = true
	lab.look()
	run := lab.orchestrator()
	if got := lab.reasonsGiven(run.PID); len(got) != 1 || got[0] != "asked: orchestrate turned on" {
		t.Errorf("FLAI_REASONS = %q, want asked alone, not the hour-old move", got)
	}
	lab.release("orchestrator")
	lab.ended(run.Session)
	for _, shift := range []time.Duration{0, time.Hour} {
		lab.shift = shift
		lab.look()
		if r := lab.orchestrator(); r.Session != run.Session {
			t.Fatalf("started again with nothing queued: %+v", r)
		}
	}
	if starts := lab.started(); len(starts) != 1 {
		t.Errorf("the journal's starts: %q, want one", starts)
	}
}

// ADR-0209 §2: the queue is kept in the serve folder, so that a flai serve
// started again, a new orchestrator, starts the run for what the last one
// queued, and does not take itself for turned on.
func TestTheQueueSurvivesARestartOfFlaiServe(t *testing.T) {
	lab := newOrchestrateLab(t)
	lab.history()
	story, _ := lab.inProgress("Reviewed")
	base := time.Now().UTC().Truncate(time.Second)
	first := lab.firstRun(base)
	lab.moveAt(story, workitem.Review, base.Add(time.Second))
	lab.at = base.Add(2 * time.Second)
	lab.look()
	if r := lab.orchestrator(); r.Session != first.Session {
		t.Fatalf("started inside the window: %+v", r)
	}
	lab.orch = newOrchestrator(lab.orch.o, lab.entry(), lab.l)
	lab.at = base.Add(2*time.Second + manifest.DefaultCoalesce)
	lab.look()
	run := lab.orchestrator()
	if run.Session == first.Session {
		t.Fatalf("the new orchestrator did not start for the reason queued before it: %+v", run)
	}
	if got := lab.reasonsGiven(run.PID); len(got) != 1 || got[0] != "moved "+story+" to review: Reviewed" {
		t.Errorf("FLAI_REASONS = %q, want the move alone", got)
	}
	lab.ended(run.Session)
}

// ADR-0209 §1: orchestration.schedule coming round queues a schedule reason,
// once for the times missed between two looks; one set comes round first at
// its next time from the look that saw it.
func TestTheScheduleComingRoundIsAReason(t *testing.T) {
	lab := newOrchestrateLab(t)
	base := time.Date(2030, 1, 1, 9, 0, 0, 0, time.UTC)
	first := lab.firstRun(base)
	m := "version: 1\nname: t\nkey: t\nlayout:\n  design: design\n  docs: docs\n  wip: wip\norchestration:\n  schedule: \"0 * * * *\"\n"
	if err := os.WriteFile(filepath.Join(lab.root, "system-flow.yaml"), []byte(m), 0o644); err != nil {
		t.Fatal(err)
	}
	lab.at = base.Add(30 * time.Minute) // the schedule is seen, set
	lab.look()
	if q := lab.o.Dir.reasonQueue(lab.root); len(q.Reasons) != 0 {
		t.Fatalf("queued when set: %+v", q.Reasons)
	}
	lab.at = base.Add(3*time.Hour + time.Minute) // three times missed come round once
	lab.look()
	q := lab.o.Dir.reasonQueue(lab.root)
	if len(q.Reasons) != 1 || q.Reasons[0].line() != "schedule 0 * * * * came round" {
		t.Fatalf("queue = %+v, want the schedule once", q.Reasons)
	}
	if r := lab.orchestrator(); r.Session != first.Session {
		t.Fatalf("started inside the window: %+v", r)
	}
	lab.at = lab.at.Add(manifest.DefaultCoalesce)
	lab.look()
	if r := lab.orchestrator(); r.Session == first.Session {
		t.Fatalf("not started for the schedule: %+v", r)
	}
	lab.ended(lab.orchestrator().Session)
}

// The queue holds a reason once, by kind, ID, and to; is due for an asked
// reason at once, and otherwise once coalesce has passed since the first
// was queued; and a start takes out the reasons it took, and no others.
func TestTheReasonQueue(t *testing.T) {
	t0 := time.Date(2030, 1, 1, 9, 0, 0, 0, time.UTC)
	var q reasonQueue
	if ready, at := q.due(t0, time.Minute); ready || !at.IsZero() {
		t.Errorf("an empty queue is due: %v %v", ready, at)
	}
	moved := orchestrateReason{Kind: workitem.Moved, ID: "S-0001", To: workitem.Review}
	if n := q.add(t0, moved, moved); n != 1 {
		t.Errorf("queued %d of a reason given twice, want 1", n)
	}
	backlog := moved
	backlog.To = workitem.Backlog
	if n := q.add(t0.Add(10*time.Second), moved, backlog); n != 1 || len(q.Reasons) != 2 || q.First != t0.Format(time.RFC3339) {
		t.Errorf("queue = %+v, want the move to backlog beside the move to review", q)
	}
	if ready, at := q.due(t0.Add(59*time.Second), time.Minute); ready || !at.Equal(t0.Add(time.Minute)) {
		t.Errorf("due inside the window: %v %v", ready, at)
	}
	if ready, _ := q.due(t0.Add(time.Minute), time.Minute); !ready {
		t.Error("not due once the window passed")
	}
	took := append([]orchestrateReason{}, q.Reasons[0])
	q.remove(took)
	if len(q.Reasons) != 1 || q.Reasons[0].To != workitem.Backlog || q.First != t0.Add(10*time.Second).Format(time.RFC3339) {
		t.Errorf("after the start: %+v, want the move to backlog, first queued when it was", q)
	}
	q.add(t0.Add(20*time.Second), askedReason(askedOperator, t0))
	if ready, _ := q.due(t0.Add(20*time.Second), time.Minute); !ready {
		t.Error("an asked reason is not due at once")
	}
}

// The orchestrator's triggers are its role file's, unless it has none, or
// it still says continuous, when they are ADR-0209 §1's; only the events a
// trigger matches are reasons, a task's never.
func TestTheRoleFilesTriggersChooseTheReasons(t *testing.T) {
	file := roles.Role{Name: roles.Orchestrator, Trigger: []roles.Trigger{{Kind: roles.TriggerEvent, Event: workitem.Moved, When: "to-review"}}}
	if got := orchestratorTriggers(nil); len(got) != len(builtinTriggers) {
		t.Errorf("no role file: %+v", got)
	}
	continuous := roles.Role{Name: roles.Orchestrator, Trigger: []roles.Trigger{{Kind: roles.TriggerContinuous}}}
	if got := orchestratorTriggers([]roles.Role{continuous}); len(got) != len(builtinTriggers) {
		t.Errorf("continuous: %+v", got)
	}
	if got := orchestratorTriggers([]roles.Role{{Name: "planner"}, file}); len(got) != 1 {
		t.Errorf("the file's: %+v", got)
	}
	if scheduled(file.Trigger) || !scheduled(builtinTriggers) {
		t.Error("scheduled reads the schedule trigger wrong")
	}
	ev := func(kind, typ, id, to string, fields map[string]string) inbox.Event {
		return inbox.Event{Change: workitem.Change{Kind: kind, Type: typ, ID: id, Title: "T\nwo", To: to}, Summary: "s", Fields: fields}
	}
	events := []inbox.Event{
		ev(workitem.Moved, workitem.Story, "S-0001", workitem.Review, nil),
		ev(workitem.Moved, workitem.Task, "T-0001", workitem.Backlog, nil),
		ev(workitem.Moved, workitem.Story, "S-0002", workitem.Backlog, nil),
		ev(workitem.Moved, workitem.Story, "S-0003", workitem.InProgress, nil),
		ev(inbox.Thread, inbox.Thread, "TH-0001", "you", map[string]string{inbox.FieldAwaits: "orchestrator"}),
		ev(inbox.Thread, inbox.Thread, "TH-0002", "operator", map[string]string{inbox.FieldAwaits: "operator"}),
		ev(inbox.Agent, inbox.RoleType, "orchestrator", inbox.AgentEnded, nil),
		ev(inbox.Overlapped, workitem.Story, "S-0004", "docs", nil),
		ev(reordered, "board", "board", "", nil),
	}
	var ids []string
	for _, r := range reasonsOf(events, builtinTriggers, nil) {
		ids = append(ids, r.ID)
	}
	if got := strings.Join(ids, " "); got != "S-0001 S-0002 TH-0001 TH-0002 board" {
		t.Errorf("built-in reasons: %s", got)
	}
	ids = nil
	for _, r := range reasonsOf(events, file.Trigger, nil) {
		ids = append(ids, r.ID)
		if r.Title != "T wo" {
			t.Errorf("title %q, want it on one line", r.Title)
		}
	}
	if got := strings.Join(ids, " "); got != "S-0001" {
		t.Errorf("the file's reasons: %s", got)
	}
	legacy := []roles.Trigger{{Kind: roles.TriggerEvent, Event: "review"}}
	if got := reasonsOf(events, legacy, nil); len(got) != 1 || got[0].ID != "S-0001" {
		t.Errorf("{event: review}: %+v", got)
	}
}

// A run's reasons are named in a journal entry or a started event's why up
// to reasonsNamed, each cut to reasonSaidMax, and the rest counted.
func TestReasonsSaidCountsTheRest(t *testing.T) {
	var rs []orchestrateReason
	for i := range reasonsNamed + 2 {
		rs = append(rs, orchestrateReason{Kind: workitem.Moved, ID: "S-000" + strconv.Itoa(i), To: workitem.Review, Title: strings.Repeat("long ", 40)})
	}
	said := reasonsSaid(rs)
	if !strings.HasPrefix(said, strconv.Itoa(reasonsNamed+2)+" reasons: moved S-0000 to review: long") || !strings.HasSuffix(said, "; and 2 more") {
		t.Errorf("said %q", said)
	}
	for part := range strings.SplitSeq(said, "; ") {
		if len([]rune(part)) > reasonSaidMax+len("7 reasons: ")+1 {
			t.Errorf("a reason was not cut: %q", part)
		}
	}
	sched, err := cron.Parse("daily")
	if err != nil {
		t.Fatal(err)
	}
	if got := reasonsSaid([]orchestrateReason{scheduleReason(&sched, time.Now())}); got != "1 reason: schedule daily came round" {
		t.Errorf("said %q", got)
	}
}
