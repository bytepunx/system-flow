//go:build !windows

package serve

import (
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"testing"
	"time"

	ctxpack "github.com/bytepunx/system-flow/flai/internal/context"
	"github.com/bytepunx/system-flow/flai/internal/cron"
	"github.com/bytepunx/system-flow/flai/internal/inbox"
	"github.com/bytepunx/system-flow/flai/internal/itemedit"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/roles"
	"github.com/bytepunx/system-flow/flai/internal/threads"
	"github.com/bytepunx/system-flow/flai/internal/workflow"
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
// ADR-0209 §5: the reasons a run took wait for it to end. A run that fails,
// on the usage limit or otherwise, gives them back as they were queued, so
// that the next run takes them; any other end drops them, and the end of a
// run that took none, or of another session, settles nothing.
func TestARunThatFailsGivesItsReasonsBack(t *testing.T) {
	t0 := time.Date(2030, 1, 1, 9, 0, 0, 0, time.UTC)
	review := orchestrateReason{Kind: workitem.Moved, ID: "S-0001", To: workitem.Review}
	planned := orchestrateReason{Kind: "planned", ID: "E-0001"}
	later := orchestrateReason{Kind: workitem.Moved, ID: "S-0002", To: workitem.Review}
	for _, c := range []struct {
		name    string
		session string
		failed  bool
		given   int
		left    []string
	}{
		{"failed", "s1", true, 2, []string{"S-0001", "E-0001", "S-0002"}},
		{"worked", "s1", false, 0, []string{"S-0002"}},
		{"another session", "s2", true, 0, []string{"S-0002"}},
	} {
		var q reasonQueue
		q.add(t0, review, planned)
		took := slices.Clone(q.Reasons)
		q.take("s1", took)
		q.add(t0.Add(time.Minute), later)
		if len(q.Reasons) != 1 || len(q.Taken) != 2 || q.Session != "s1" {
			t.Fatalf("%s: after the start %+v, want two taken and the later one queued", c.name, q)
		}
		if n := q.settle(c.session, c.failed); n != c.given {
			t.Errorf("%s: gave back %d, want %d", c.name, n, c.given)
		}
		var left []string
		for _, r := range q.Reasons {
			left = append(left, r.ID)
		}
		if !slices.Equal(left, c.left) {
			t.Errorf("%s: queued %v, want %v", c.name, left, c.left)
		}
		if c.session == "s1" && (len(q.Taken) != 0 || q.Session != "") {
			t.Errorf("%s: still taken %+v", c.name, q)
		}
		if c.failed && c.session == "s1" && q.First != t0.Format(time.RFC3339) {
			t.Errorf("%s: first = %s, want the given back ones' queued time, so that they are due at once", c.name, q.First)
		}
	}
	var q reasonQueue
	if n := q.settle("", true); n != 0 {
		t.Errorf("a run that took nothing gave back %d", n)
	}
}

func TestTheRoleFilesTriggersChooseTheReasons(t *testing.T) {
	file := roles.Role{Name: roles.Orchestrator, Trigger: []roles.Trigger{{Kind: roles.TriggerEvent, Event: workitem.Moved, When: roles.CondTo, To: workitem.Review}}}
	if got := orchestratorTriggers(workflow.Baseline(), nil); len(got) != len(builtinTriggers) {
		t.Errorf("no role file: %+v", got)
	}
	continuous := roles.Role{Name: roles.Orchestrator, Trigger: []roles.Trigger{{Kind: roles.TriggerContinuous}}}
	if got := orchestratorTriggers(workflow.Baseline(), []roles.Role{continuous}); len(got) != len(builtinTriggers) {
		t.Errorf("continuous: %+v", got)
	}
	if got := orchestratorTriggers(workflow.Baseline(), []roles.Role{{Name: "planner"}, file}); len(got) != 1 {
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
	for _, r := range reasonsOf(workflow.Baseline(), events, builtinTriggers, nil) {
		ids = append(ids, r.ID)
	}
	if got := strings.Join(ids, " "); got != "S-0001 S-0002 TH-0001 TH-0002 board" {
		t.Errorf("built-in reasons: %s", got)
	}
	ids = nil
	for _, r := range reasonsOf(workflow.Baseline(), events, file.Trigger, nil) {
		ids = append(ids, r.ID)
		if r.Title != "T wo" {
			t.Errorf("title %q, want it on one line", r.Title)
		}
	}
	if got := strings.Join(ids, " "); got != "S-0001" {
		t.Errorf("the file's reasons: %s", got)
	}
	legacy := []roles.Trigger{{Kind: roles.TriggerEvent, Event: "review"}}
	if got := reasonsOf(workflow.Baseline(), events, legacy, nil); len(got) != 1 || got[0].ID != "S-0001" {
		t.Errorf("{event: review}: %+v", got)
	}
}

// verifyingBoard is the baseline board's front matter with a state of the
// project's own, verifying, between in-progress and review.
const verifyingBoard = `---
title: Board
updated: 2026-11-15
status: active
states:
  - {name: backlog, kind: queue, initial: true, order: true}
  - {name: ready, kind: queue, order: true, limit: true}
  - {name: in-progress, kind: active, limit: true}
  - {name: verifying, kind: waiting}
  - {name: review, kind: waiting, limit: true}
  - {name: done, kind: accepted}
  - {name: cancelled, kind: cancelled}
transitions:
  - {from: backlog, to: ready, effects: [commit], rules: [criteria-present, not-draft, after-done]}
  - {from: ready, to: in-progress, effects: [pull], rules: [not-blocked, not-held, task-after]}
  - {from: in-progress, to: verifying, rules: [has-task, clean-worktree]}
  - {from: verifying, to: review, rules: [no-open-question, kluster-verified]}
  - {from: verifying, to: in-progress, effects: [send-back], rules: [reason-given, not-blocked, task-after]}
  - {from: review, to: done, types: [epic, story], by: [acceptance], effects: [accept], rules: [children-closed, criteria-ticked]}
  - {from: review, to: done, types: [task], rules: [children-closed]}
  - {from: in-progress, to: done, types: [task], rules: [children-closed]}
  - {from: review, to: in-progress, effects: [send-back], rules: [reason-given, not-blocked, task-after]}
  - {from: ready, to: backlog}
  - {from: [backlog, ready, in-progress, verifying], to: cancelled, effects: [cancel], rules: [reason-given]}
  - {from: review, to: cancelled, by: [cascade], effects: [cancel], rules: [reason-given]}
  - {from: cancelled, to: backlog, effects: [reopen], rules: [parent-open]}
---
`

// gateRole is a role a project defines that starts on a story's move into
// verifying, a state of the project's own workflow (ADR-0304 §3).
const gateRole = `---
title: Gate
updated: 2026-11-15
status: active
name: gate
role: gate
definition:
  description: Checks a story as it enters verifying, for flai serve.
  tools: [Read, Grep]
  model: inherit
action: gate
trigger:
  - {event: moved, when: {to: verifying}}
item: [story]
topics: [planning]
works_story: false
charge: none
tools: []
commands: []
files: []
permissions: []
---

You are the gate: check the story as it enters verifying.

<!-- system-flow:end-of-baseline -->

## Project additions
`

// S-0748, ADR-0304 §3: a role whose trigger names a state the project's
// workflow adds between in-progress and review is valid against that
// workflow and starts on a move into that state, and not on a move into
// review or any other.
func TestARoleStartsOnAMoveIntoAStateOfTheProjectsOwn(t *testing.T) {
	wf, wp := workflow.Parse("design/workflows/board.md", verifyingBoard)
	if wp != nil {
		t.Fatalf("workflow.Parse: %s", wp)
	}
	r, p := roles.Parse("design/workflows/roles/gate.md", gateRole)
	if p != nil {
		t.Fatalf("roles.Parse: %s", p)
	}
	if ps := roles.Validate([]roles.Role{r}, "design", roles.Vocabulary{Workflow: &wf}); len(ps) != 0 {
		t.Fatalf("roles.Validate against the project's workflow = %v, want none", ps)
	}
	if ps := roles.Validate([]roles.Role{r}, "design", roles.Vocabulary{}); len(ps) != 1 || !strings.Contains(ps[0].Message, `state "verifying" is none of the workflow's states`) {
		t.Errorf("roles.Validate against the baseline = %v, want verifying refused", ps)
	}
	for _, c := range []struct {
		to   string
		want bool
	}{{"verifying", true}, {workitem.Review, false}, {workitem.InProgress, false}, {workitem.Backlog, false}} {
		got := matching(r.Trigger, roles.Occurrence{Event: workitem.Moved, To: c.to})
		if (len(got) == 1) != c.want {
			t.Errorf("a move into %s: matching = %+v, want it to start: %t", c.to, got, c.want)
		}
	}
	events := []inbox.Event{
		{Change: workitem.Change{Kind: workitem.Moved, Type: workitem.Story, ID: "S-0001", Title: "Verify", To: "verifying"}, Summary: "s"},
		{Change: workitem.Change{Kind: workitem.Moved, Type: workitem.Story, ID: "S-0002", Title: "Review", To: workitem.Review}, Summary: "s"},
		{Change: workitem.Change{Kind: inbox.Edited, Type: workitem.Story, ID: "S-0003", Title: "Edit", To: "verifying"}, Summary: "s"},
	}
	if got := reasonsOf(wf, events, r.Trigger, nil); len(got) != 1 || got[0].ID != "S-0001" || got[0].To != "verifying" {
		t.Errorf("reasons = %+v, want the move of S-0001 into verifying alone", got)
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

// ADR-0217 §6: each reason flai serve starts the orchestrator for has the
// form prime takes, and asked, schedule, and an event with no form, none;
// every form serve writes is one flai prime --reason takes.
func TestEachReasonHasTheFormPrimeTakes(t *testing.T) {
	sched, err := cron.Parse("daily")
	if err != nil {
		t.Fatal(err)
	}
	for _, c := range []struct {
		r    orchestrateReason
		want string
	}{
		{orchestrateReason{Kind: workitem.Moved, ID: "S-0123", Type: workitem.Story, To: workitem.Review}, "review:S-0123"},
		{orchestrateReason{Kind: workitem.Moved, ID: "S-0123", Type: workitem.Story, To: workitem.Backlog}, "draft:S-0123"},
		{orchestrateReason{Kind: inbox.Edited, ID: "S-0123", Type: workitem.Story, To: "criteria,touches"}, "draft:S-0123"},
		{orchestrateReason{Kind: inbox.Planned, ID: "E-0016", Type: workitem.Epic, To: "S-0200,S-0201"}, "planned:E-0016"},
		{orchestrateReason{Kind: inbox.Planned, ID: "S-0208", Type: workitem.Story}, "planned:S-0208"},
		{orchestrateReason{Kind: inbox.Thread, ID: "TH-0007", Type: inbox.Thread, To: "you"}, "thread:TH-0007"},
		{orchestrateReason{Kind: inbox.Thread, ID: "TH-0008", Type: inbox.Thread, To: "operator"}, "thread:TH-0008"},
		{orchestrateReason{Kind: inbox.ReleaseDue, ID: "S-0123", Type: workitem.Story}, "release"},
		{orchestrateReason{Kind: inbox.ReadyEmpty, ID: "board", Type: "board", To: "1"}, "ready"},
		{orchestrateReason{Kind: reordered, ID: "board"}, "ready"},
		{askedReason(askedOperator, time.Now()), ""},
		{scheduleReason(&sched, time.Now()), ""},
		{orchestrateReason{Kind: workitem.Moved, ID: "S-0123", Type: workitem.Story, To: workitem.Done}, ""},
		{orchestrateReason{Kind: workitem.Moved, ID: "E-0016", Type: workitem.Epic, To: workitem.Review}, ""},
		{orchestrateReason{Kind: inbox.Edited, ID: "E-0016", Type: workitem.Epic, To: "title"}, ""},
		{orchestrateReason{Kind: "agent", ID: "S-0123"}, ""},
	} {
		got := c.r.prime(workflow.Baseline())
		if got != c.want {
			t.Errorf("%s: prime %q, want %q", c.r.line(), got, c.want)
		}
		if got == "" {
			continue
		}
		parsed, err := ctxpack.ParseReasons([]string{got})
		if err != nil || len(parsed) != 1 || parsed[0].String() != got {
			t.Errorf("%s: flai prime --reason %s reads %+v, %v", c.r.line(), got, parsed, err)
		}
	}
}

// ADR-0217 §6: a run's forms for prime are one for each reason, in turn,
// a form given twice kept for its line, and none when any reason has no
// form, so that the run primes with the general pack.
func TestPrimeReasonsAreNoneWhenAnyReasonHasNoForm(t *testing.T) {
	review := orchestrateReason{Kind: workitem.Moved, ID: "S-0123", Type: workitem.Story, To: workitem.Review}
	empty := orchestrateReason{Kind: inbox.ReadyEmpty, ID: "board", Type: "board", To: "1"}
	order := orchestrateReason{Kind: reordered, ID: "board"}
	if got, want := primeReasons(workflow.Baseline(), []orchestrateReason{review, empty, order}), []string{"review:S-0123", "ready", "ready"}; !slices.Equal(got, want) {
		t.Errorf("primeReasons %q, want %q", got, want)
	}
	if got := primeReasons(workflow.Baseline(), []orchestrateReason{review, askedReason(askedTurnedOn, time.Now())}); got != nil {
		t.Errorf("with asked, primeReasons %q, want none", got)
	}
	if got := primeReasons(workflow.Baseline(), nil); got != nil {
		t.Errorf("with no reasons, primeReasons %q, want none", got)
	}
}

// t3216Moves are a story's moves, S-0001 to the first of states, S-0002 to
// the second, and so on, as the inbox gives them.
func t3216Moves(states ...string) []inbox.Event {
	var out []inbox.Event
	for i, to := range states {
		out = append(out, inbox.Event{Change: workitem.Change{Kind: workitem.Moved, Type: workitem.Story, ID: "S-000" + strconv.Itoa(i+1), Title: "T", To: to}, Summary: "s"})
	}
	return out
}

// t3216Reasons are the reasons triggers find among events read in wf, each
// as its ID and the form prime takes it in, joined by spaces.
func t3216Reasons(wf workflow.Workflow, events []inbox.Event, triggers []roles.Trigger) string {
	var out []string
	for _, r := range reasonsOf(wf, events, triggers, nil) {
		out = append(out, r.ID+"="+r.prime(wf))
	}
	return strings.Join(out, " ")
}

// S-0742 criterion 4: the orchestrator's reasons read a moved story's state
// by its kind. On a board where review is awaiting and in-progress doing,
// with testing, of the active kind, between them, a story's move to
// awaiting is the reason a move to review is on the baseline, with the same
// form for prime; a move to doing or testing is none that a move to review
// would be; and the moves to backlog, ready, done, and cancelled are what
// they are on the baseline, under the built-in triggers read by kind, a role
// file's to condition naming the board's own state, and each event name
// ADR-0207 gives a move.
func TestTheReasonsReadAMovedStatesKind(t *testing.T) {
	kinds := *kindsWorkflow(t)
	base := workflow.Baseline()
	baseMoves := t3216Moves(workitem.Review, workitem.InProgress, workitem.InProgress, workitem.Backlog, workitem.Ready, workitem.Done, workitem.Cancelled)
	kindsMoves := t3216Moves("awaiting", "doing", "testing", "backlog", "ready", "done", "cancelled")
	ev := func(event, when string) []roles.Trigger {
		return []roles.Trigger{{Kind: roles.TriggerEvent, Event: event, When: when}}
	}
	to := func(state string) []roles.Trigger {
		return []roles.Trigger{{Kind: roles.TriggerEvent, Event: workitem.Moved, When: roles.CondTo, To: state}}
	}
	// The built-in triggers are read on each board by kind; a role file
	// names the board's own states (ADR-0304), so its condition names
	// review on the baseline and awaiting on the other board.
	for _, c := range []struct {
		name        string
		base, other []roles.Trigger
		want        string
	}{
		{"built-in", builtinTriggersIn(base), builtinTriggersIn(kinds), "S-0001=review:S-0001 S-0004=draft:S-0004"},
		{"to the waiting state", to(workitem.Review), to("awaiting"), "S-0001=review:S-0001"},
		{"to the initial state", to(workitem.Backlog), to("backlog"), "S-0004=draft:S-0004"},
		{"review", ev(eventReview, ""), ev(eventReview, ""), "S-0001=review:S-0001"},
		{"ready", ev(eventReady, ""), ev(eventReady, ""), "S-0005="},
		{"accepted", ev(eventAccepted, ""), ev(eventAccepted, ""), "S-0006="},
		{"cancelled", ev(eventCancelled, ""), ev(eventCancelled, ""), "S-0007="},
	} {
		if got := t3216Reasons(base, baseMoves, c.base); got != c.want {
			t.Errorf("%s on the baseline: %q, want %q", c.name, got, c.want)
		}
		if got := t3216Reasons(kinds, kindsMoves, c.other); got != c.want {
			t.Errorf("%s on the other board: %q, want %q, as on the baseline", c.name, got, c.want)
		}
	}
	// A role file that names review on the other board names a state it
	// has not, so no move there starts it.
	if got := t3216Reasons(kinds, kindsMoves, to(workitem.Review)); got != "" {
		t.Errorf("a trigger on review reads the other board's moves as %q, want none", got)
	}
	// Read in the baseline, a move to awaiting is a move to a state the
	// baseline has not: it is no review.
	if got := t3216Reasons(base, t3216Moves("awaiting"), builtinTriggersIn(base)); got != "" {
		t.Errorf("the baseline reads a move to awaiting as %q, want no reason", got)
	}
	for _, to := range []string{"doing", "testing"} {
		r := orchestrateReason{Kind: workitem.Moved, ID: "S-0123", Type: workitem.Story, To: to}
		if got := r.prime(kinds); got != "" {
			t.Errorf("a move to %s has the form %q, want none", to, got)
		}
		if names := eventNames(kinds, t3216Moves(to)[0]); !slices.Equal(names, []string{workitem.Moved}) {
			t.Errorf("a move to %s is the events %v, want moved alone", to, names)
		}
	}
	// An epic's move to awaiting is a review event with no form for prime,
	// and its move to done no acceptance, as on the baseline.
	epic := t3216Moves("awaiting")[0]
	epic.Type, epic.ID = workitem.Epic, "E-0001"
	if names := eventNames(kinds, epic); !slices.Equal(names, []string{workitem.Moved, eventReview}) {
		t.Errorf("an epic's move to awaiting is the events %v, want moved and review", names)
	}
	if got := (orchestrateReason{Kind: workitem.Moved, ID: "E-0001", Type: workitem.Epic, To: "awaiting"}).prime(kinds); got != "" {
		t.Errorf("an epic's move to awaiting has the form %q, want none", got)
	}
	epic.To = "done"
	if names := eventNames(kinds, epic); !slices.Equal(names, []string{workitem.Moved}) {
		t.Errorf("an epic's move to done is the events %v, want moved alone", names)
	}
}

// S-0742: the curator's acceptances are the stories' moves into a state of
// the accepted kind, on the other board as on the baseline: a move to
// awaiting, or to cancelled, is none, and a task's or an epic's move to done
// is none.
func TestTheCuratorsAcceptancesReadTheAcceptedKind(t *testing.T) {
	changes := []workitem.Change{
		{ID: "S-0003", Type: workitem.Story, Kind: workitem.Moved, To: "done"},
		{ID: "S-0001", Type: workitem.Story, Kind: workitem.Moved, To: "awaiting"},
		{ID: "S-0002", Type: workitem.Story, Kind: workitem.Moved, To: "cancelled"},
		{ID: "T-0001", Type: workitem.Task, Kind: workitem.Moved, To: "done"},
		{ID: "E-0001", Type: workitem.Epic, Kind: workitem.Moved, To: "done"},
	}
	if got := acceptedIn(*kindsWorkflow(t), nil, changes); !slices.Equal(got, []string{"S-0003"}) {
		t.Errorf("acceptedIn = %v, want S-0003 alone", got)
	}
	if o := acceptance(); o.Event != eventAccepted || o.To != workitem.Done {
		t.Errorf("an acceptance is %+v, want accepted into the baseline's done", o)
	}
}
