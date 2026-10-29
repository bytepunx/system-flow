//go:build !windows

package serve

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/conventions"
	"github.com/bytepunx/system-flow/flai/internal/harness"
	"github.com/bytepunx/system-flow/flai/internal/hostapi"
	"github.com/bytepunx/system-flow/flai/internal/inbox"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/threads"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// planLab is an agentLab whose command is a planner stub: it writes what it
// was given to plan-<item>.txt, prints stream-<item> as its output when
// there is one, waits while held until released, and exits with the code in
// exit-<item> when there is one. The plan action is on.
func planLab(t *testing.T) *agentLab {
	t.Helper()
	lab := newAgentLab(t)
	stub := filepath.Join(lab.outDir, "plan-agent")
	d := lab.outDir
	script := "#!/bin/sh\n" +
		"{ echo \"args: $*\"; echo \"dir: $(pwd)\"; echo \"agent: $FLAI_AGENT\"; echo \"story: ${FLAI_STORY-unset}\"; echo \"role: $FLAI_ROLE\"; echo \"item: $FLAI_ITEM\"; echo \"session: $FLAI_SESSION\"; echo \"by: $FLAI_STARTED_BY\"; } > \"" + d + "/plan-$FLAI_ITEM.txt\"\n" +
		"[ -f \"" + d + "/stream-$FLAI_ITEM\" ] && cat \"" + d + "/stream-$FLAI_ITEM\"\n" +
		"while [ -f \"" + d + "/hold\" ] && [ ! -f \"" + d + "/release-$FLAI_ITEM\" ]; do sleep 0.05; done\n" +
		"touch \"" + d + "/ended-$$\"\n" +
		"[ -f \"" + d + "/exit-$FLAI_ITEM\" ] && exit \"$(cat \"" + d + "/exit-$FLAI_ITEM\")\"\n" +
		"exit 0\n"
	if err := os.WriteFile(stub, []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	lab.cfg.Command, lab.cfg.Plan = []string{stub, "story={story}"}, true
	return lab
}

func (lab *agentLab) plan(id string) (*AgentRun, error) {
	return Plan(context.Background(), lab.o, Entry{Key: "t", Name: "t", Root: lab.root}, id)
}

// planHere has the lab's launcher start the planner for id, as the operator
// asked, and wait for it, as a serving flai would.
func (lab *agentLab) planHere(id string) bool {
	lab.l.mu.Lock()
	defer lab.l.mu.Unlock()
	started, _ := lab.l.plan(context.Background(), lab.cfg, id, nil, TriggerAsked)
	return started
}

func (lab *agentLab) planRun(id string) *AgentRun { return lab.state().Plans[id] }

// given is what the planner stub was given for item, once it has written it.
func (lab *agentLab) given(item string) string {
	lab.t.Helper()
	out := filepath.Join(lab.outDir, "plan-"+item+".txt")
	waitFor(lab.t, "the planner stub ran for "+item, func() bool {
		data, _ := os.ReadFile(out)
		return strings.Contains(string(data), "by: ")
	})
	data, _ := os.ReadFile(out)
	return string(data)
}

// S-0208: the planner is started for an epic in the project's main
// checkout, as the planner and for its item alone, logged as the planner's
// run, and recorded by item, apart from the stories' runs, as asked
// (ADR-0084); a second run on the item is refused while the first runs.
func TestThePlannerIsStartedForAnItem(t *testing.T) {
	t.Setenv("FLAI_STORY", "S-0999") // the session that ran flai plan
	lab := planLab(t)
	lab.hold()
	id := lab.epic.ID
	run, err := lab.plan(id)
	if err != nil {
		t.Fatal(err)
	}
	if run.Item != id || run.Story != "" || run.Agent != "planner-"+id || run.PID == 0 || run.Session == "" || run.Harness != "command" || run.Command != "plan-agent" || run.Trigger != "asked" {
		t.Errorf("run = %+v", run)
	}
	if data, err := os.ReadFile(lab.o.Dir.agents()); err != nil || !strings.Contains(string(data), `"trigger": "asked"`) {
		t.Errorf("serve/agents.json does not record the run as asked (%v):\n%s", err, data)
	}
	logs, err := lab.o.Dir.ActivityLogs("t", workitem.ActivityPlanner)
	if err != nil || len(logs) != 1 || logs[0] != run.Log {
		t.Errorf("planner logs = %v (%v), want the run's %s", logs, err, run.Log)
	}
	got := lab.given(id)
	real, _ := filepath.EvalSymlinks(lab.root)
	for _, want := range []string{"args: story=\n", "dir: " + real + "\n", "agent: planner-" + id + "\n", "story: unset\n", "role: plan\n", "item: " + id + "\n", "session: 2", "by: flai-serve\n"} {
		if !strings.Contains(got, want) {
			t.Errorf("the planner was given\n%s\nwithout %q", got, want)
		}
	}
	if st := lab.state(); st.Stories[id] != nil || st.Running != nil || !st.Plans[id].same(run) {
		t.Errorf("state = %+v, want the run by item alone", st)
	}
	j := lab.entries()
	if len(j) != 1 || j[0].Action != hostapi.ActionPlan || j[0].Method != "serve.plan" || j[0].Outcome != "done" || !strings.Contains(j[0].Detail, "started plan-agent (command) to plan "+id+" as planner-"+id+" (pid ") {
		t.Errorf("journal: %+v", j)
	}
	_, err = lab.plan(id)
	var no *Refused
	if !errors.As(err, &no) || !strings.Contains(no.Why, "the planner is already running for "+id) || !strings.Contains(no.Why, "started "+run.Started) {
		t.Errorf("a second run: %v", err)
	}
	// another item's planner, started in the same second, logs apart, and
	// neither holds a story's agent back
	story := lab.ready("Planned")
	other, err := lab.plan(story)
	if err != nil {
		t.Fatalf("another item: %v", err)
	}
	if logs, _ := lab.o.Dir.ActivityLogs("t", workitem.ActivityPlanner); len(logs) != 2 || other.Log == run.Log || !strings.Contains(lab.given(story), "item: "+story+"\n") {
		t.Errorf("planner logs = %v, want one for each run", logs)
	}
	lab.l.look(context.Background(), false)
	waitFor(t, "the story's agent started beside the planners", func() bool { return lab.run(story) != nil })
	lab.release(id)
	lab.release(story)
}

// S-0208: a planner run that ends is recorded as a story's agent's is, with
// its exit and how it went, and its activity is logged in the planner's
// document from its log, with what started it (ADR-0084).
func TestAPlannerRunThatEndsIsRecordedAndItsActivityLogged(t *testing.T) {
	lab := planLab(t)
	id := lab.epic.ID
	stream := strings.Join([]string{streamCall("s", "m1", runStart, 999), streamCall("s", "m2", runStart.Add(time.Minute), 999), streamFinal("Drafted three stories for "+id+".", 2000, 0.5)}, "\n") + "\n"
	if err := os.WriteFile(filepath.Join(lab.outDir, "stream-"+id), []byte(stream), 0o644); err != nil {
		t.Fatal(err)
	}
	if !lab.planHere(id) {
		t.Fatalf("not started: %+v", lab.planRun(id))
	}
	waitFor(t, "the planner run ended", func() bool { r := lab.planRun(id); return r != nil && r.Ended != "" })
	run := lab.planRun(id)
	if run.Outcome != OutcomeWorked || run.Exit == nil || *run.Exit != 0 || run.Session == "" || run.Log == "" {
		t.Errorf("run = %+v, want worked with exit 0", run)
	}
	waitFor(t, "its activity logged", func() bool {
		doc, err := lab.repo.Activity(workitem.ActivityPlanner)
		return err == nil && len(doc.Entries) == 1
	})
	doc, _ := lab.repo.Activity(workitem.ActivityPlanner)
	if e := doc.Entries[0]; e.Summary != "Drafted three stories for "+id+"." || e.Trigger != "asked" || e.Cost != 0.5 || e.Seconds != 61 {
		t.Errorf("activity = %+v, want the run's summary, asked, its 0.5 USD, and its 61 s", e)
	}
	if data, _ := os.ReadFile(doc.Path); !strings.Contains(string(data), "\n- Summary: Drafted three stories for "+id+".\n- Trigger: asked\n") {
		t.Errorf("the entry does not say it was asked:\n%s", data)
	}
	// ADR-0083: and what it spent is charged to the epic it planned, once
	// the entry is in
	waitFor(t, "its cost charged", func() bool { it, err := lab.repo.Get(id); return err == nil && it.Usage != nil })
	it, _ := lab.repo.Get(id)
	if u := it.Usage; u == nil || len(u.Strategic) != 1 || u.Strategic[0].Kind != workitem.ActivityPlanner || u.Strategic[0].Cost() != 0.5 || u.Strategic[0].Seconds != 61 {
		t.Errorf("epic's usage = %+v, want the planner's 0.5 USD over 61 s", u)
	}
}

// S-0219: a planner run the orchestrator asked for records orchestrator, in
// place of asked, as what started it, and its activity entry says so.
func TestAPlannerRunTheOrchestratorAskedForSaysSo(t *testing.T) {
	lab := planLab(t)
	lab.hold()
	id := lab.epic.ID
	run, err := PlanForOrchestrator(context.Background(), lab.o, Entry{Key: "t", Name: "t", Root: lab.root}, id)
	if err != nil {
		t.Fatal(err)
	}
	if run.Trigger != TriggerOrchestrator || run.Agent != "planner-"+id {
		t.Errorf("run = %+v, want the trigger orchestrator", run)
	}
	if data, err := os.ReadFile(lab.o.Dir.agents()); err != nil || !strings.Contains(string(data), `"trigger": "orchestrator"`) {
		t.Errorf("serve/agents.json does not record the run as the orchestrator's (%v):\n%s", err, data)
	}
	_, err = PlanForOrchestrator(context.Background(), lab.o, Entry{Key: "t", Name: "t", Root: lab.root}, id)
	var no *Refused
	if !errors.As(err, &no) || !strings.Contains(no.Why, "the planner is already running for "+id) {
		t.Errorf("refused as Plan is: %v", err)
	}
	lab.release(id)

	// the run that ends logs its activity with the orchestrator as its trigger
	story := lab.backlog("Planned", nil)
	stream := strings.Join([]string{streamCall("s", "m1", runStart, 999), streamFinal("Enriched "+story+".", 2000, 0.25)}, "\n") + "\n"
	if err := os.WriteFile(filepath.Join(lab.outDir, "stream-"+story), []byte(stream), 0o644); err != nil {
		t.Fatal(err)
	}
	lab.release(story)
	lab.l.mu.Lock()
	started, _ := lab.l.plan(context.Background(), lab.cfg, story, nil, TriggerOrchestrator)
	lab.l.mu.Unlock()
	if !started {
		t.Fatalf("not started: %+v", lab.planRun(story))
	}
	waitFor(t, "its activity logged", func() bool {
		doc, err := lab.repo.Activity(workitem.ActivityPlanner)
		return err == nil && len(doc.Entries) == 1
	})
	doc, _ := lab.repo.Activity(workitem.ActivityPlanner)
	if e := doc.Entries[0]; e.Trigger != "orchestrator" || e.Summary != "Enriched "+story+"." {
		t.Errorf("activity = %+v, want the orchestrator as its trigger", e)
	}
	if data, _ := os.ReadFile(doc.Path); !strings.Contains(string(data), "\n- Trigger: orchestrator\n") {
		t.Errorf("the entry does not say the orchestrator asked:\n%s", data)
	}
}

// S-0209: a planner run's activity names its item, then what it created
// under the item, then what it changed there, and not what it left alone or
// what is under another item.
func TestAPlannerRunsActivityNamesWhatItPlanned(t *testing.T) {
	lab := planLab(t)
	before, during := runStart.Add(-time.Hour), runStart.Add(30*time.Second)
	story := func(parent, title string, at time.Time) *workitem.Item {
		t.Helper()
		it, err := lab.repo.Create(workitem.NewOptions{Type: workitem.Story, Title: title, Parent: parent, Owner: "alex", Now: at})
		if err != nil {
			t.Fatal(err)
		}
		return it
	}
	left := story(lab.epic.ID, "Left alone", before)
	edited := story(lab.epic.ID, "Edited", before)
	edited.Updated = runStart.Format(workitem.TimeFormat) // the second the run started
	if err := lab.repo.Save(edited); err != nil {
		t.Fatal(err)
	}
	made := story(lab.epic.ID, "Drafted", during)
	task, err := lab.repo.Create(workitem.NewOptions{Type: workitem.Task, Title: "Work", Parent: made.ID, Owner: "alex", Now: during})
	if err != nil {
		t.Fatal(err)
	}
	other, err := lab.repo.Create(workitem.NewOptions{Type: workitem.Epic, Title: "Other", Owner: "alex", Now: before})
	if err != nil {
		t.Fatal(err)
	}
	elsewhere := story(other.ID, "Elsewhere", during)
	logs := filepath.Join(string(lab.o.Dir), "agents")
	if err := os.MkdirAll(logs, 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(logs, "t-planner-20261003T100000Z.log"), []byte(strings.Join([]string{
		streamCall("s", "m1", runStart, 999), streamFinal("Drafted "+made.ID+".", 1000, 0.5)}, "\n")+"\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	zero := 0
	// coalesced triggers (ADR-0084) are logged as the run recorded them
	trigger := "edited goal by alex; schedule daily"
	lab.l.planEnded(&AgentRun{Item: lab.epic.ID, Agent: "planner-" + lab.epic.ID, Started: runStart.Format(time.RFC3339), Session: "s", Trigger: trigger}, &zero)
	doc, err := lab.repo.Activity(workitem.ActivityPlanner)
	if err != nil || len(doc.Entries) != 1 {
		t.Fatalf("activity = %+v (%v), want one entry", doc, err)
	}
	want := strings.Join([]string{lab.epic.ID, made.ID, task.ID, edited.ID}, ",")
	if e := doc.Entries[0]; strings.Join(e.Items, ",") != want || e.Cost != 0.5 || e.Trigger != trigger {
		t.Errorf("entry = %+v, want items %s (not %s or %s), the run's 0.5 USD, and trigger %q", e, want, left.ID, elsewhere.ID, trigger)
	}
	// planned on a story, the run names the story and its tasks
	if got, err := plannedItems(lab.root, made.ID, runStart); err != nil || strings.Join(got, ",") != made.ID+","+task.ID {
		t.Errorf("items planned on %s = %v (%v), want it and %s", made.ID, got, err, task.ID)
	}
}

// S-0255: a planner run on a story, which drafts its tasks or revisits the
// ones it has, names the story in its activity, then the tasks it created
// under it, then those it changed there, and not a task it left alone or a
// task of another story.
func TestAPlannerRunOnAStoryNamesItsTasks(t *testing.T) {
	lab := planLab(t)
	before, during := runStart.Add(-time.Hour), runStart.Add(30*time.Second)
	item := func(typ, parent, title string, at time.Time) *workitem.Item {
		t.Helper()
		it, err := lab.repo.Create(workitem.NewOptions{Type: typ, Title: title, Parent: parent, Owner: "alex", Now: at})
		if err != nil {
			t.Fatal(err)
		}
		return it
	}
	story := item(workitem.Story, lab.epic.ID, "Planned", before)
	left := item(workitem.Task, story.ID, "Left alone", before)
	edited := item(workitem.Task, story.ID, "Edited", before)
	edited.Updated = runStart.Format(workitem.TimeFormat) // the second the run started
	if err := lab.repo.Save(edited); err != nil {
		t.Fatal(err)
	}
	made := item(workitem.Task, story.ID, "Drafted", during)
	other := item(workitem.Story, lab.epic.ID, "Other", before)
	elsewhere := item(workitem.Task, other.ID, "Elsewhere", during)
	logs := filepath.Join(string(lab.o.Dir), "agents")
	if err := os.MkdirAll(logs, 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(logs, "t-planner-20261003T100000Z.log"), []byte(strings.Join([]string{
		streamCall("s", "m1", runStart, 999), streamFinal("Drafted "+made.ID+" for "+story.ID+".", 1000, 0.5)}, "\n")+"\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	zero := 0
	lab.l.planEnded(&AgentRun{Item: story.ID, Agent: "planner-" + story.ID, Started: runStart.Format(time.RFC3339), Session: "s"}, &zero)
	doc, err := lab.repo.Activity(workitem.ActivityPlanner)
	if err != nil || len(doc.Entries) != 1 {
		t.Fatalf("activity = %+v (%v), want one entry", doc, err)
	}
	want := strings.Join([]string{story.ID, made.ID, edited.ID}, ",")
	if e := doc.Entries[0]; strings.Join(e.Items, ",") != want || e.Cost != 0.5 {
		t.Errorf("entry = %+v, want items %s (not %s or %s) and the run's 0.5 USD", e, want, left.ID, elsewhere.ID)
	}
}

// S-0535: a planner run that wrote a story's tasks appends one planned event
// naming the story, the tasks it wrote, and the planner; a run that wrote
// nothing beyond its item appends none.
func TestAPlannerRunThatWroteItemsLogsAPlannedEvent(t *testing.T) {
	lab := planLab(t)
	before, during := runStart.Add(-time.Hour), runStart.Add(30*time.Second)
	item := func(typ, parent, title string, at time.Time) *workitem.Item {
		t.Helper()
		it, err := lab.repo.Create(workitem.NewOptions{Type: typ, Title: title, Parent: parent, Owner: "alex", Now: at})
		if err != nil {
			t.Fatal(err)
		}
		return it
	}
	idle := item(workitem.Story, lab.epic.ID, "Left as it was", before)
	item(workitem.Task, idle.ID, "Untouched", before)
	zero := 0
	end := func(id string) {
		lab.l.planEnded(&AgentRun{Item: id, Agent: "planner-" + id, Started: runStart.Format(time.RFC3339), Session: "s"}, &zero)
	}
	// the lab's looks may log a ready-empty event: only planned ones count
	planned := func() []inbox.LoggedEvent {
		var out []inbox.LoggedEvent
		for _, e := range inbox.LoggedEvents(lab.repo) {
			if e.Kind == inbox.Planned {
				out = append(out, e)
			}
		}
		return out
	}
	end(idle.ID)
	if got := planned(); len(got) != 0 {
		t.Fatalf("events after a run that wrote nothing = %+v, want none", got)
	}
	story := item(workitem.Story, lab.epic.ID, "Planned", before)
	first := item(workitem.Task, story.ID, "First", during)
	second := item(workitem.Task, story.ID, "Second", during)
	end(story.ID)
	got := planned()
	if len(got) != 1 {
		t.Fatalf("events = %+v, want one planned event", got)
	}
	e := got[0]
	if e.Kind != inbox.Planned || e.ID != story.ID || e.Type != workitem.Story || e.Title != "Planned" ||
		e.To != first.ID+","+second.ID || e.By != "planner-"+story.ID || e.At == "" {
		t.Errorf("event = %+v, want planned on %s by planner-%s writing %s,%s, stamped", e, story.ID, story.ID, first.ID, second.ID)
	}
}

// S-0208: a planner that ended with its question on its item open is asked;
// one that exited with a failure, failed.
func TestAPlannerRunIsJudgedByItsQuestionAndItsExit(t *testing.T) {
	lab := planLab(t)
	asks, fails := lab.backlog("Asks", nil), lab.backlog("Fails", nil)
	th, err := threads.New(lab.repo, threads.NewOptions{Title: "What is it worth?", On: asks, Author: "planner-" + asks, Text: "The value is missing.", Now: lab.now})
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(lab.outDir, "exit-"+fails), []byte("3"), 0o644); err != nil {
		t.Fatal(err)
	}
	for _, id := range []string{asks, fails} {
		if !lab.planHere(id) {
			t.Fatalf("%s not started: %+v", id, lab.planRun(id))
		}
		waitFor(t, id+"'s planner run ended", func() bool { r := lab.planRun(id); return r != nil && r.Ended != "" })
	}
	if r := lab.planRun(asks); r.Outcome != OutcomeAsked || r.Thread != th.ID || !strings.Contains(r.Why, "waiting for an answer to "+th.ID) {
		t.Errorf("asked = %+v", r)
	}
	if r := lab.planRun(fails); r.Outcome != OutcomeFailed || r.Exit == nil || *r.Exit != 3 || r.Why != "ended (exit 3) planning "+fails {
		t.Errorf("failed = %+v", r)
	}
}

// S-0208: a planner run whose process is gone while no launcher waits for
// it, as one flai plan started is, is settled at the next look.
func TestAPlannerRunGoneIsSettledAtALook(t *testing.T) {
	lab := planLab(t)
	gone := exec.Command("true")
	if err := gone.Run(); err != nil {
		t.Fatal(err)
	}
	id := lab.epic.ID
	lab.o.Dir.updateAgent(lab.root, func(s *AgentState) {
		s.put(&AgentRun{Item: id, Agent: "planner-" + id, PID: gone.Process.Pid, Started: "2026-10-03T10:00:00Z", Session: "s"})
	})
	lab.l.look(context.Background(), false)
	if r := lab.planRun(id); r.Ended == "" || r.Outcome != OutcomeWorked || r.Exit != nil {
		t.Errorf("settled = %+v, want worked with no exit seen", r)
	}
	lab.logs.mu.Lock()
	logged := lab.logs.buf.String()
	lab.logs.mu.Unlock()
	if !strings.Contains(logged, `msg="planner ended" component=serve project=t item=`+id) {
		t.Errorf("the end is not logged:\n%s", logged)
	}
}

// S-0208: the planner is refused while the plan action is off, for a task,
// for what is no epic's or story's ID, for an item that does not exist, is
// done, or is cancelled, and when nothing can start it.
func TestThePlannerIsRefusedWithWhy(t *testing.T) {
	lab := planLab(t)
	done := lab.backlog("Done", nil)
	it, err := lab.repo.Get(done)
	if err != nil {
		t.Fatal(err)
	}
	data, _ := os.ReadFile(it.Path) // accepted, as its file says
	if err := os.WriteFile(it.Path, []byte(strings.Replace(string(data), "status: backlog", "status: done", 1)), 0o644); err != nil {
		t.Fatal(err)
	}
	cancelled, err := lab.repo.Create(workitem.NewOptions{Type: workitem.Epic, Title: "Dropped", Owner: "alex", Now: lab.now})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := lab.repo.Transition(cancelled, workitem.Cancelled, "alex", "not needed", lab.now); err != nil {
		t.Fatal(err)
	}
	refusedFor := func(id, want string) {
		t.Helper()
		_, err := lab.plan(id)
		var no *Refused
		if !errors.As(err, &no) || !strings.Contains(no.Why, want) {
			t.Errorf("plan %s: %v, want refused for %q", id, err, want)
		}
	}
	refusedFor("T-0001", "T-0001 is a task; the planner plans an epic or a story")
	refusedFor("X-0001", `"X-0001" is not an epic's or a story's ID`)
	refusedFor(done, done+" is done; there is nothing left to plan")
	refusedFor(cancelled.ID, cancelled.ID+" is cancelled")
	if _, err := lab.plan("S-0999"); err == nil || !strings.Contains(err.Error(), "S-0999 not found") {
		t.Errorf("an item that does not exist: %v", err)
	}
	lab.cfg.Command = nil
	refusedFor(lab.epic.ID, "the planner's agent names no harness")
	lab.cfg.Plan = false
	refusedFor(lab.epic.ID, "the plan host action is off for this project: flai serve enable plan")
	if st := lab.state(); len(st.Plans) != 0 || len(lab.entries()) != 0 {
		t.Errorf("a refusal was recorded: %+v %+v", st.Plans, lab.entries())
	}
}

// S-0208: agents.json written before planner runs still reads, and one with
// them keeps the stories' runs apart.
func TestAgentStatesReadWithAndWithoutPlans(t *testing.T) {
	d := Dir(t.TempDir())
	if err := os.WriteFile(d.agents(), []byte(`{"/p":{"stories":{"S-0001":{"story":"S-0001","command":"claude","agent":"agent-S-0001","started":"2026-10-01T00:00:00Z"}}}}`), 0o600); err != nil {
		t.Fatal(err)
	}
	st := d.AgentStates()["/p"]
	if st.Stories["S-0001"] == nil || st.Plans != nil {
		t.Fatalf("an older file = %+v", st)
	}
	d.updateAgent("/p", func(s *AgentState) {
		s.put(&AgentRun{Item: "E-0001", Agent: "planner-E-0001", Started: "2026-10-03T00:00:00Z", Ended: "2026-10-03T00:01:00Z", Outcome: OutcomeWorked, StoryAgent: &manifest.Agent{}})
	})
	st = d.AgentStates()["/p"]
	if st.Plans["E-0001"] == nil || len(st.Stories) != 1 || st.Last != nil {
		t.Errorf("with a planner run = %+v, want it by item and the story's run untouched", st)
	}
}

// agentEventsAre says that the agent events on id are one per change, to and
// why as want, in turn, each of type typ, titled title, by by, and of role
// (S-0532). The launcher appends a run's end just after it records it, so
// they are waited for.
func (lab *agentLab) agentEventsAre(id, typ, title, by, role string, want ...[2]string) {
	lab.t.Helper()
	var got []inbox.LoggedEvent
	waitFor(lab.t, fmt.Sprintf("%d agent events on %s", len(want), id), func() bool {
		got = inbox.AgentEventsOf(lab.repo, id, time.Time{})
		return len(got) >= len(want)
	})
	if len(got) != len(want) {
		lab.t.Fatalf("%d agent events on %s, want %d: %+v", len(got), id, len(want), got)
	}
	for i, e := range got {
		if e.ID != id || e.Type != typ || e.Title != title || e.By != by || e.Fields[inbox.FieldRole] != role ||
			e.To != want[i][0] || e.Fields[inbox.FieldWhy] != want[i][1] {
			lab.t.Errorf("event %d: %+v, want %s %s %q by %s of role %s for %q", i+1, e, typ, id, title, by, role, want[i])
		}
	}
}

// S-0532: a planner run appends one agent event on the item it plans at its
// start, saying what started it, and one at its end, ended for a run that
// worked or asked and failed for one that did not, for the words judgePlan
// gives it; a start that could not be made is failed, for its error, and a
// run settled again appends nothing more.
func TestAPlannerRunAppendsAnAgentEventAtItsStartAndItsEnd(t *testing.T) {
	asked := [2]string{inbox.AgentStarted, "the operator asked for it"}
	t.Run("worked is ended", func(t *testing.T) {
		lab := planLab(t)
		id := lab.epic.ID
		if !lab.planHere(id) {
			t.Fatalf("not started: %+v", lab.planRun(id))
		}
		waitFor(t, "the planner run ended", func() bool { r := lab.planRun(id); return r != nil && r.Ended != "" })
		want := [][2]string{asked, {inbox.AgentEnded, "ended with its work done"}}
		lab.agentEventsAre(id, workitem.Epic, lab.epic.Title, "planner-"+id, conventions.RolePlan, want...)
		again := *lab.planRun(id)
		lab.l.planEnded(&again, again.Exit)
		lab.agentEventsAre(id, workitem.Epic, lab.epic.Title, "planner-"+id, conventions.RolePlan, want...)
	})
	t.Run("a failed exit is failed", func(t *testing.T) {
		lab := planLab(t)
		id := lab.backlog("Fails", nil)
		if err := os.WriteFile(filepath.Join(lab.outDir, "exit-"+id), []byte("3"), 0o644); err != nil {
			t.Fatal(err)
		}
		if !lab.planHere(id) {
			t.Fatalf("not started: %+v", lab.planRun(id))
		}
		lab.agentEventsAre(id, workitem.Story, "Fails", "planner-"+id, conventions.RolePlan, asked, [2]string{inbox.AgentFailed, "ended (exit 3) planning " + id})
	})
	t.Run("a start that could not be made is failed", func(t *testing.T) {
		lab := planLab(t)
		lab.cfg.Command = []string{filepath.Join(lab.outDir, "no-such-program")}
		id := lab.epic.ID
		if lab.planHere(id) {
			t.Fatalf("started: %+v", lab.planRun(id))
		}
		r := lab.planRun(id)
		if r == nil || !strings.HasPrefix(r.Why, "could not be started: ") {
			t.Fatalf("not recorded as failed to start: %+v", r)
		}
		lab.agentEventsAre(id, workitem.Epic, lab.epic.Title, "planner-"+id, conventions.RolePlan, [2]string{inbox.AgentFailed, r.Why})
	})
}

// S-0532: a strategic role's started event says what started it, in a short
// form of its journal entry's words.
func TestAStartedEventSaysWhatStartedTheRole(t *testing.T) {
	for trigger, want := range map[string]string{
		TriggerAsked:           "the operator asked for it",
		TriggerOrchestrator:    "the orchestrator asked for it",
		"schedule daily":       "started by schedule daily",
		"accepted S-0001":      "started by accepted S-0001",
		"edited title by alex": "started by edited title by alex",
		"":                     "trigger not recorded",
	} {
		if got := triggerWhy(&AgentRun{Trigger: trigger}); got != want {
			t.Errorf("%q: %q, want %q", trigger, got, want)
		}
	}
}

// S-0534: a planner run whose log ends on the recorded usage-limit refusal
// with no reset time, settled with its exit unseen, is failed on the limit
// and sets the budget state for the first backoff; the probe that starts
// once the backoff has passed and ends clean clears it.
func TestAPlannerRunOnTheUsageLimitSetsTheBudgetStateForTheBackoff(t *testing.T) {
	lab := newAgentLab(t)
	at := limitReset.Add(-time.Hour)
	lab.budgetClock(&at)
	item := lab.epic.ID
	run := &AgentRun{Item: item, Agent: "planner-" + item, Harness: harness.ClaudeCode, Trigger: TriggerAsked, Started: at.Format(time.RFC3339), Log: limitResult(t, false)}
	lab.l.planEnded(run, nil)
	lab.heldUntil(at.Add(budgetBackoffFirst), false, "5m0s")
	want := "ended on the harness's usage limit (usage_limit_reached, 429); it named no reset time"
	if r := lab.planRun(item); r.Outcome != OutcomeFailed || !r.UsageLimited || r.ResetsAt != "" || r.Why != want {
		t.Errorf("planner run = %+v; want failed on the limit, %q", r, want)
	}

	at = at.Add(budgetBackoffFirst)
	probe := &AgentRun{Item: item, Agent: "planner-" + item, Harness: harness.ClaudeCode, Trigger: TriggerAsked, Started: at.Format(time.RFC3339), Log: cleanResult(t)}
	lab.probing(probe)
	zero := 0
	lab.l.planEnded(probe, &zero)
	if b, held := lab.budget(); held || b.Provider != "" {
		t.Errorf("budget state after the probe's clean end = %+v, held %v; want none", b, held)
	}
	if r := lab.planRun(item); r.Outcome != OutcomeWorked || r.UsageLimited {
		t.Errorf("probe = %+v; want worked, on no limit", r)
	}
}

// budgetHeldOn sets the budget state of key, held until until, the reset
// the API reported, as a run that ended on the usage limit sets it.
func (lab *agentLab) budgetHeldOn(key string, until time.Time) {
	lab.t.Helper()
	all, err := lab.o.Dir.budgets()
	if err != nil {
		lab.t.Fatal(err)
	}
	until = until.UTC()
	all[key] = Budget{Provider: key, Limit: "usage_limit_reached", Status: 429, Since: until.Add(-time.Hour), Until: until, Reset: until, Backoff: budgetBackoffFirst.String()}
	if err := lab.o.Dir.write(lab.o.Dir.budgetFile(), all); err != nil {
		lab.t.Fatal(err)
	}
}

// budgetCleared clears the budget state of the stub's harness, command, as
// its probe's clean end does.
func (lab *agentLab) budgetCleared() {
	lab.t.Helper()
	all, err := lab.o.Dir.budgets()
	if err != nil {
		lab.t.Fatal(err)
	}
	delete(all, harness.Command)
	if err := lab.o.Dir.write(lab.o.Dir.budgetFile(), all); err != nil {
		lab.t.Fatal(err)
	}
}

// heldFor is the refusal of a start on the stub's harness, command, while
// budgetHeldOn holds it until until.
func heldFor(until time.Time) string {
	return "held: command usage limit (usage_limit_reached) until " + until.UTC().Format(time.RFC3339) + "; flai serve starts it when the limit resets"
}

// anHour is a time an hour from now, to the second, which a budget state
// set until it holds a start at the launcher's clock.
func anHour() time.Time { return time.Now().Add(time.Hour).UTC().Truncate(time.Second) }

// S-0534: flai plan, and the orchestrator's asking, is refused while the
// budget state of the planner's provider holds, naming the provider and
// until, with no run recorded and nothing journalled; a state on another
// provider holds nothing, and once until has passed the start goes as the
// provider's probe, recorded with its process.
func TestAPlanIsRefusedWhileTheBudgetStateHoldsItsProvider(t *testing.T) {
	lab := planLab(t)
	lab.hold()
	until := anHour()
	lab.budgetHeldOn(harness.Command, until)
	id := lab.epic.ID
	for name, start := range map[string]func(context.Context, Options, Entry, string) (*AgentRun, error){"flai plan": Plan, "the orchestrator's asking": PlanForOrchestrator} {
		if _, err := start(context.Background(), lab.o, Entry{Key: "t", Name: "t", Root: lab.root}, id); refusal(err) != heldFor(until) {
			t.Errorf("%s: %v; want refused with %q", name, err, heldFor(until))
		}
	}
	if r := lab.planRun(id); r != nil || len(lab.entries()) != 0 || lab.startedCount() != 0 {
		t.Fatalf("a held start recorded %+v, journalled %+v", r, lab.entries())
	}

	lab.budgetCleared()
	lab.budgetHeldOn("gateway", until)
	if _, err := lab.plan(id); err != nil {
		t.Fatalf("held by another provider's state: %v", err)
	}
	lab.release(id)

	lab.budgetHeldOn(harness.Command, time.Now().Add(-time.Minute))
	story := lab.ready("Probed")
	run, err := lab.plan(story)
	if err != nil {
		t.Fatalf("once until has passed: %v", err)
	}
	b, held, err := lab.o.Dir.budgetHeld(harness.Command, time.Now())
	if err != nil || !held || b.Probe == nil || b.Probe.Work != story || b.Probe.PID != run.PID || b.Probe.Started != run.Started {
		t.Errorf("budget state %+v, held %v (%v); want %s's run the probe, with its process", b, held, err, story)
	}
	lab.release(story)
}

// S-0534: a start flai serve makes on its own that the budget state of its
// provider holds is not made and records nothing: the launcher says it was
// held, and by what.
func TestAPlannerStartTheBudgetStateHoldsRecordsNothing(t *testing.T) {
	lab := planLab(t)
	until := anHour()
	lab.budgetHeldOn(harness.Command, until)
	lab.l.mu.Lock()
	started, held := lab.l.plan(context.Background(), lab.cfg, lab.epic.ID, nil, TriggerAsked)
	lab.l.mu.Unlock()
	if started || held == nil || held.Provider != harness.Command || !held.Until.Equal(until) {
		t.Errorf("started %v, held %+v; want held by command's state until %s", started, held, until)
	}
	if r := lab.planRun(lab.epic.ID); r != nil || len(lab.entries()) != 0 || lab.startedCount() != 0 {
		t.Errorf("a held start recorded %+v, journalled %+v", r, lab.entries())
	}
}
