//go:build !windows

package serve

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/hostapi"
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

// planHere has the lab's launcher start the planner for id and wait for it,
// as a serving flai would.
func (lab *agentLab) planHere(id string) bool {
	lab.l.mu.Lock()
	defer lab.l.mu.Unlock()
	return lab.l.plan(context.Background(), lab.cfg, id, nil)
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
// run, and recorded by item, apart from the stories' runs; a second run on
// the item is refused while the first runs.
func TestThePlannerIsStartedForAnItem(t *testing.T) {
	t.Setenv("FLAI_STORY", "S-0999") // the session that ran flai plan
	lab := planLab(t)
	lab.hold()
	id := lab.epic.ID
	run, err := lab.plan(id)
	if err != nil {
		t.Fatal(err)
	}
	if run.Item != id || run.Story != "" || run.Agent != "planner-"+id || run.PID == 0 || run.Session == "" || run.Harness != "command" || run.Command != "plan-agent" {
		t.Errorf("run = %+v", run)
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
// document from its log.
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
	if e := doc.Entries[0]; e.Summary != "Drafted three stories for "+id+"." || e.Cost != 0.5 || e.Seconds != 61 {
		t.Errorf("activity = %+v, want the run's summary, its 0.5 USD, and its 61 s", e)
	}
	// ADR-0083: and what it spent is charged to the epic it planned, once
	// the entry is in
	waitFor(t, "its cost charged", func() bool { it, err := lab.repo.Get(id); return err == nil && it.Usage != nil })
	it, _ := lab.repo.Get(id)
	if u := it.Usage; u == nil || len(u.Strategic) != 1 || u.Strategic[0].Kind != workitem.ActivityPlanner || u.Strategic[0].Cost() != 0.5 || u.Strategic[0].Seconds != 61 {
		t.Errorf("epic's usage = %+v, want the planner's 0.5 USD over 61 s", u)
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
	lab.l.planEnded(&AgentRun{Item: lab.epic.ID, Agent: "planner-" + lab.epic.ID, Started: runStart.Format(time.RFC3339), Session: "s"}, &zero)
	doc, err := lab.repo.Activity(workitem.ActivityPlanner)
	if err != nil || len(doc.Entries) != 1 {
		t.Fatalf("activity = %+v (%v), want one entry", doc, err)
	}
	want := strings.Join([]string{lab.epic.ID, made.ID, task.ID, edited.ID}, ",")
	if e := doc.Entries[0]; strings.Join(e.Items, ",") != want || e.Cost != 0.5 {
		t.Errorf("entry = %+v, want items %s (not %s or %s) and the run's 0.5 USD", e, want, left.ID, elsewhere.ID)
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
