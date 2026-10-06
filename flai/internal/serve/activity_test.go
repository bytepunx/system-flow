package serve

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/issues"
	"github.com/bytepunx/system-flow/flai/internal/usage"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// activityLab is a project named t and the serve folder its strategic
// agents' logs are kept in.
type activityLab struct {
	t    *testing.T
	root string
	dir  Dir
	repo *workitem.Repo
}

func newActivityLab(t *testing.T) *activityLab {
	t.Helper()
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "system-flow.yaml"), []byte("version: 1\nname: t\nkey: t\nlayout:\n  design: design\n  docs: docs\n  wip: wip\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	for _, d := range []string{"wip/kanban/epics", "wip/kanban/stories", "wip/kanban/tasks", "wip/agents"} {
		if err := os.MkdirAll(filepath.Join(root, d), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	dir := Dir(filepath.Join(t.TempDir(), "serve"))
	if err := os.MkdirAll(filepath.Join(string(dir), "agents"), 0o700); err != nil {
		t.Fatal(err)
	}
	return &activityLab{t: t, root: root, dir: dir, repo: repo}
}

// log writes one run's log under the serve folder.
func (lab *activityLab) log(name string, lines ...string) {
	lab.t.Helper()
	if err := os.WriteFile(filepath.Join(string(lab.dir), "agents", name), []byte(strings.Join(lines, "\n")+"\n"), 0o600); err != nil {
		lab.t.Fatal(err)
	}
}

// doc reads kind's activity document from disk.
func (lab *activityLab) doc(kind string) *workitem.Activity {
	lab.t.Helper()
	a, err := lab.repo.Activity(kind)
	if err != nil {
		lab.t.Fatal(err)
	}
	return a
}

// streamUser is a tool result handed back to session at at.
func streamUser(session string, at time.Time) string {
	return fmt.Sprintf(`{"type":"user","timestamp":%q,"session_id":%q,"message":{"content":[{"type":"tool_result","content":"ok"}]}}`, at.UTC().Format(time.RFC3339Nano), session)
}

// streamFinal is the end of a run of session s that answers text: its totals.
func streamFinal(text string, read int, cost float64) string {
	return fmt.Sprintf(`{"type":"result","subtype":"success","session_id":"s","result":%q,"modelUsage":{"claude-opus-5-5":{"inputTokens":2,"outputTokens":300,"cacheReadInputTokens":%d,"cacheCreationInputTokens":0,"costUSD":%v}}}`, text, read, cost)
}

// toTheCentsHundredth is c as the log keeps it.
func toTheCentsHundredth(c float64) float64 { return math.Round(c*1e4) / 1e4 }

var runStart = time.Date(2026, 10, 3, 10, 0, 0, 0, time.UTC)

// A strategic agent's logs are its kind's, named for the time each run
// started, apart from the stories' logs and another project's.
func TestActivityLogsAreTheKindsOwn(t *testing.T) {
	lab := newActivityLab(t)
	for _, name := range []string{
		"t-planner-20261003T110000Z.log", "t-planner-20261003T100000Z.log", "t-orchestrator-20261003T100000Z.log",
		"t-S-0001-20261003T100000Z.log", "other-planner-20261003T100000Z.log", "t-planner-S-0002-20261003T100000Z.log",
	} {
		lab.log(name, streamCall("s-"+name, "m-"+name, runStart, 1))
	}
	got, err := lab.dir.ActivityLogs("t", workitem.ActivityPlanner)
	if err != nil {
		t.Fatal(err)
	}
	var names []string
	for _, p := range got {
		names = append(names, filepath.Base(p))
	}
	if strings.Join(names, ",") != "t-planner-20261003T100000Z.log,t-planner-20261003T110000Z.log" {
		t.Errorf("planner logs = %v, want the project's two, oldest first", names)
	}
	if stories, _ := lab.dir.LoggedStories("t"); strings.Join(stories, ",") != "S-0001" {
		t.Errorf("logged stories = %v, want the story's alone", stories)
	}
}

// An activity is its run's wall-clock time since the run began and the
// share of the run's reported cost its calls in that time weigh.
func TestAnActivityIsMeasuredFromItsRun(t *testing.T) {
	lab := newActivityLab(t)
	// calls weigh 1000, 1000, and 2000 of the session's 4000
	lab.log("t-orchestrator-20261003T100000Z.log",
		streamCall("s", "m1", runStart, 999),
		streamCall("s", "m2", runStart.Add(time.Minute), 999),
		streamCall("s", "m3", runStart.Add(2*time.Minute), 1999),
		streamResult("s", 4000, 1.2345))
	got, err := LogActivity(lab.dir, lab.root, "t", workitem.ActivityOrchestrator, "Ordered the board", []string{"S-0001", "S-0002"}, runStart.Add(90*time.Second))
	if err != nil {
		t.Fatal(err)
	}
	e := got.Entry
	if e.Seconds != 90 || e.Cost != toTheCentsHundredth(1.2345*2000/4000) || !e.Estimated || !e.At.Equal(runStart.Add(90*time.Second)) {
		t.Fatalf("entry = %+v, want 90 s and half of 1.2345 USD, estimated", e)
	}
	if strings.Join(e.Items, ",") != "S-0001,S-0002" || e.Summary != "Ordered the board" || len(got.Logs) != 1 {
		t.Errorf("logged = %+v", got)
	}
	doc := lab.doc(workitem.ActivityOrchestrator)
	if doc.AccruedCost != e.Cost || doc.AccruedSeconds != 90 || doc.TasksCompleted != 1 || doc.LastRun != "2026-10-03T10:01:30Z" || len(doc.Entries) != 1 {
		t.Errorf("document = %+v, want the entry's totals", doc)
	}
}

// Two activities in one run: the second begins where the first ended, and
// the totals are their sum.
func TestTheNextActivityBeginsWhereTheLastEnded(t *testing.T) {
	lab := newActivityLab(t)
	lab.log("t-orchestrator-20261003T100000Z.log",
		streamCall("s", "m1", runStart, 999),
		streamCall("s", "m2", runStart.Add(time.Minute), 999),
		streamCall("s", "m3", runStart.Add(2*time.Minute), 1999),
		// the tool result the second activity_log call answers
		streamUser("s", runStart.Add(150*time.Second)),
		streamResult("s", 4000, 1.2345))
	first, err := LogActivity(lab.dir, lab.root, "t", workitem.ActivityOrchestrator, "first", nil, runStart.Add(90*time.Second))
	if err != nil {
		t.Fatal(err)
	}
	second, err := LogActivity(lab.dir, lab.root, "t", workitem.ActivityOrchestrator, "second", nil, runStart.Add(150*time.Second))
	if err != nil {
		t.Fatal(err)
	}
	if e := second.Entry; e.Seconds != 60 || e.Cost != toTheCentsHundredth(1.2345*2000/4000) {
		t.Fatalf("second = %+v, want 60 s and m3's half", e)
	}
	doc := lab.doc(workitem.ActivityOrchestrator)
	if doc.AccruedSeconds != 150 || doc.AccruedCost != toTheCentsHundredth(first.Entry.Cost+second.Entry.Cost) || doc.TasksCompleted != 2 {
		t.Errorf("document = %+v, want the two entries' sum", doc)
	}
}

// A session resumed in a later run reports cumulative totals; an activity in
// the later run is charged only its own calls' share of them.
func TestAnActivityInAResumedRunIsChargedItsShare(t *testing.T) {
	lab := newActivityLab(t)
	lab.log("t-orchestrator-20261003T100000Z.log",
		streamCall("s", "m1", runStart, 999),
		streamCall("s", "m2", runStart.Add(time.Minute), 999),
		streamFinal("Waiting for an answer.", 2000, 2.0))
	ended, err := LogRunEnd(lab.dir, lab.root, "t", workitem.ActivityOrchestrator, nil)
	if err != nil || ended == nil {
		t.Fatalf("the first run's end = %+v (%v), want logged", ended, err)
	}
	if e := ended.Entry; e.Cost != 2.0 || e.Seconds != 61 || e.Summary != "Waiting for an answer." {
		t.Errorf("first run = %+v, want its whole 2.0 USD over 61 s", e)
	}
	later := runStart.Add(time.Hour)
	lab.log("t-orchestrator-20261003T110000Z.log",
		streamCall("s", "m3", later, 1999),
		streamUser("s", later.Add(30*time.Second)),
		streamFinal("Done.", 4000, 4.0))
	got, err := LogActivity(lab.dir, lab.root, "t", workitem.ActivityOrchestrator, "answered", nil, later.Add(30*time.Second))
	if err != nil {
		t.Fatal(err)
	}
	if e := got.Entry; e.Cost != 2.0 || e.Seconds != 30 {
		t.Fatalf("resumed = %+v, want m3's 2000 of 4000 weight of 4.0 USD over 30 s", e)
	}
	if doc := lab.doc(workitem.ActivityOrchestrator); doc.AccruedCost != 4.0 {
		t.Errorf("accrued = %v, want the session's 4.0", doc.AccruedCost)
	}
}

// A run that ends after its agent logged its last activity, with nothing
// spent since, logs nothing more.
func TestARunEndWithNothingSpentSinceLogsNothing(t *testing.T) {
	lab := newActivityLab(t)
	lab.log("t-planner-20261003T100000Z.log",
		streamCall("s", "m1", runStart, 999),
		streamCall("s", "m2", runStart.Add(time.Minute), 999),
		streamUser("s", runStart.Add(95*time.Second)),
		streamFinal("Planned S-0001.", 2000, 2.0))
	if _, err := LogActivity(lab.dir, lab.root, "t", workitem.ActivityPlanner, "planned", []string{"S-0001"}, runStart.Add(90*time.Second)); err != nil {
		t.Fatal(err)
	}
	got, err := LogRunEnd(lab.dir, lab.root, "t", workitem.ActivityPlanner, nil)
	if err != nil || got != nil {
		t.Fatalf("run end = %+v (%v), want nothing logged", got, err)
	}
	if doc := lab.doc(workitem.ActivityPlanner); doc.TasksCompleted != 1 {
		t.Errorf("document = %+v, want the one entry", doc)
	}
}

// A run that ends with calls since its last activity logs them as one, with
// the first line of the run's final text as its summary, and no trigger when
// none is given (ADR-0084).
func TestARunEndLogsWhatWasSpentSinceTheLastActivity(t *testing.T) {
	lab := newActivityLab(t)
	lab.log("t-planner-20261003T100000Z.log",
		streamCall("s", "m1", runStart, 999),
		streamCall("s", "m2", runStart.Add(time.Minute), 2999),
		streamFinal("\n  Planned S-0001 and S-0002.  \nThe rest of the answer.", 4000, 1.0))
	if _, err := LogActivity(lab.dir, lab.root, "t", workitem.ActivityPlanner, "read the board", nil, runStart.Add(30*time.Second)); err != nil {
		t.Fatal(err)
	}
	got, err := LogRunEnd(lab.dir, lab.root, "t", workitem.ActivityPlanner, nil)
	if err != nil || got == nil {
		t.Fatalf("run end = %+v (%v), want logged", got, err)
	}
	e := got.Entry
	if e.Summary != "Planned S-0001 and S-0002." || e.Trigger != "" || len(e.Items) != 0 || e.Seconds != 31 || e.Cost != 0.75 || !e.At.Equal(runStart.Add(61*time.Second)) {
		t.Errorf("entry = %+v, want the result's first line, no trigger, 31 s, and m2's 0.75 USD", e)
	}
	if doc := lab.doc(workitem.ActivityPlanner); doc.TasksCompleted != 2 || doc.AccruedCost != 1.0 {
		t.Errorf("document = %+v, want both entries' totals", doc)
	}
}

// A run whose final text says nothing ends with a summary that says so.
func TestARunEndWithNoTextIsSaidToHaveEnded(t *testing.T) {
	lab := newActivityLab(t)
	lab.log("t-analyzer-20261003T100000Z.log", streamCall("s", "m1", runStart, 999), streamResult("s", 1000, 0.5))
	got, err := LogRunEnd(lab.dir, lab.root, "t", workitem.ActivityAnalyzer, nil)
	if err != nil || got == nil || got.Entry.Summary != "run ended" || got.Entry.Cost != 0.5 {
		t.Fatalf("run end = %+v (%v), want run ended at 0.5 USD", got, err)
	}
}

// An activity outside any run flai serve logged is still recorded, with no
// seconds and no cost; a run end with no run logs nothing.
func TestAnActivityWithNoRunIsLoggedWithNothingMeasured(t *testing.T) {
	lab := newActivityLab(t)
	if got, err := LogRunEnd(lab.dir, lab.root, "t", workitem.ActivityPlanner, nil); err != nil || got != nil {
		t.Fatalf("run end with no run = %+v (%v)", got, err)
	}
	got, err := LogActivity(lab.dir, lab.root, "t", workitem.ActivityPlanner, "planned by hand", []string{"S-0001"}, runStart.Add(-time.Hour))
	if err != nil {
		t.Fatal(err)
	}
	if e := got.Entry; e.Seconds != 0 || e.Cost != 0 || e.Estimated || len(got.Logs) != 0 {
		t.Errorf("entry = %+v, want nothing measured", got)
	}
	// a run logged to its end, then an activity by hand an hour later
	lab.log("t-planner-20261003T100000Z.log", streamCall("s", "m1", runStart, 999), streamResult("s", 1000, 0.5))
	if _, err := LogRunEnd(lab.dir, lab.root, "t", workitem.ActivityPlanner, nil); err != nil {
		t.Fatal(err)
	}
	got, err = LogActivity(lab.dir, lab.root, "t", workitem.ActivityPlanner, "planned by hand again", nil, runStart.Add(time.Hour))
	if err != nil {
		t.Fatal(err)
	}
	if e := got.Entry; e.Seconds != 0 || e.Cost != 0 {
		t.Errorf("entry after a logged run = %+v, want nothing measured", e)
	}
	if doc := lab.doc(workitem.ActivityPlanner); doc.TasksCompleted != 3 || doc.AccruedCost != 0.5 {
		t.Errorf("document = %+v", doc)
	}
}

// Only a strategic agent's kind is measured or logged.
func TestAnUnknownKindIsRefused(t *testing.T) {
	lab := newActivityLab(t)
	if _, err := lab.dir.ActivityLogs("t", "builder"); err == nil || !strings.Contains(err.Error(), "planner, orchestrator, analyzer") {
		t.Errorf("logs of builder: %v, want refused naming the kinds", err)
	}
	if _, err := LogActivity(lab.dir, lab.root, "t", "builder", "built", nil, runStart); err == nil {
		t.Error("logging builder's activity was not refused")
	}
	if _, err := LogRunEnd(lab.dir, lab.root, "t", "../stories", nil); err == nil {
		t.Error("ending a run of ../stories was not refused")
	}
	if m, _ := filepath.Glob(filepath.Join(lab.root, "wip", "agents", "*.md")); len(m) != 0 {
		t.Errorf("wrote %v", m)
	}
}

// planned makes an epic and a story under it, and returns their IDs.
func (lab *activityLab) planned() (epic, story string) {
	lab.t.Helper()
	e, err := lab.repo.Create(workitem.NewOptions{Type: workitem.Epic, Title: "Epic", Owner: "alex", Now: runStart})
	if err != nil {
		lab.t.Fatal(err)
	}
	s, err := lab.repo.Create(workitem.NewOptions{Type: workitem.Story, Title: "Story", Parent: e.ID, Owner: "alex", Now: runStart})
	if err != nil {
		lab.t.Fatal(err)
	}
	return e.ID, s.ID
}

// planRun records in serve/agents.json, as flai serve does, that the newest
// planner run for item kept its log in name.
func (lab *activityLab) planRun(item, name string) {
	lab.dir.updateAgent(lab.root, func(s *AgentState) {
		s.put(&AgentRun{Item: item, Agent: "planner-" + item, Started: runStart.Format(time.RFC3339), Log: filepath.Join(string(lab.dir), "agents", name)})
	})
}

// strategic is what the planner spent on id, as its usage says; nil when
// nothing.
func (lab *activityLab) strategic(id string) *usage.Strategic {
	lab.t.Helper()
	return lab.strategicOf(id, workitem.ActivityPlanner)
}

// strategicOf is what the strategic agent kind spent on id, as its usage
// says; nil when nothing.
func (lab *activityLab) strategicOf(id, kind string) *usage.Strategic {
	lab.t.Helper()
	it, err := lab.repo.Get(id)
	if err != nil {
		lab.t.Fatal(err)
	}
	if it.Usage == nil {
		return nil
	}
	for i, s := range it.Usage.Strategic {
		if s.Kind == kind {
			return &it.Usage.Strategic[i]
		}
	}
	return nil
}

// ADR-0083: a planner run for a story ends; what its activity spent, the
// share its entry's cost is, is charged to the story and to its epic, apart
// from the agents' figures.
func TestAPlannerRunEndIsChargedToTheItemItPlannedAndAbove(t *testing.T) {
	lab := newActivityLab(t)
	epic, story := lab.planned()
	name := "t-planner-20261003T100000Z.log"
	lab.log(name,
		streamCall("s", "m1", runStart, 999),
		streamCall("s", "m2", runStart.Add(time.Minute), 999),
		streamFinal("Drafted the story's tasks.", 2000, 0.5))
	lab.planRun(story, name)
	got, err := LogRunEnd(lab.dir, lab.root, "t", workitem.ActivityPlanner, nil)
	if err != nil || got == nil {
		t.Fatalf("run end = %+v (%v), want logged", got, err)
	}
	if got.Planned != story || strings.Join(got.Charged, ",") != story+","+epic {
		t.Errorf("charged %s: %v, want %s and its epic", got.Planned, got.Charged, story)
	}
	e := got.Entry
	for _, id := range []string{story, epic} {
		s := lab.strategic(id)
		if s == nil || s.Cost() != e.Cost || e.Cost != 0.5 || s.Seconds != e.Seconds || !s.Estimated || s.Tokens() != 2+300+2000 {
			t.Errorf("%s's planner usage = %+v, want the entry's %v USD over %d s and the run's 2302 tokens, estimated", id, s, e.Cost, e.Seconds)
		}
		if u := lab.usageOf(id); !u.Empty() || u.Source != usage.SourceSum {
			t.Errorf("%s's agents' figures = %+v, want none, summed from nothing", id, u)
		}
	}
}

// ADR-0083: two activities in one run each charge only their span's share,
// so that the item ends with their entries' sum.
func TestEachPlannerActivityChargesItsShare(t *testing.T) {
	lab := newActivityLab(t)
	epic, _ := lab.planned()
	name := "t-planner-20261003T100000Z.log"
	lab.log(name,
		streamCall("s", "m1", runStart, 999),
		streamCall("s", "m2", runStart.Add(time.Minute), 2999),
		streamFinal("Drafted two stories.", 4000, 1.0))
	lab.planRun(epic, name)
	first, err := LogActivity(lab.dir, lab.root, "t", workitem.ActivityPlanner, "read the epic", nil, runStart.Add(30*time.Second))
	if err != nil {
		t.Fatal(err)
	}
	if s := lab.strategic(epic); s == nil || s.Cost() != first.Entry.Cost || first.Entry.Cost != 0.25 || s.Seconds != 30 {
		t.Fatalf("after the first activity = %+v, want its 0.25 USD over 30 s", s)
	}
	second, err := LogRunEnd(lab.dir, lab.root, "t", workitem.ActivityPlanner, nil)
	if err != nil || second == nil {
		t.Fatalf("run end = %+v (%v)", second, err)
	}
	s := lab.strategic(epic)
	if s == nil || s.Cost() != toTheCentsHundredth(first.Entry.Cost+second.Entry.Cost) || s.Cost() != 1.0 || s.Seconds != first.Entry.Seconds+second.Entry.Seconds {
		t.Errorf("epic's planner usage = %+v, want the two entries' sum, 1.0 USD", s)
	}
	if n := s.Tokens(); n < 4301 || n > 4303 {
		t.Errorf("tokens = %d, want the run's 4302 split between the two", n)
	}
}

// ADR-0083: a session resumed in a later run reports cumulative totals; the
// item is charged the later run's share alone, not the first run's again.
func TestAResumedPlannerSessionIsNotChargedTwice(t *testing.T) {
	lab := newActivityLab(t)
	_, story := lab.planned()
	lab.log("t-planner-20261003T100000Z.log",
		streamCall("s", "m1", runStart, 999),
		streamFinal("Waiting for an answer.", 1000, 2.0))
	lab.planRun(story, "t-planner-20261003T100000Z.log")
	if _, err := LogRunEnd(lab.dir, lab.root, "t", workitem.ActivityPlanner, nil); err != nil {
		t.Fatal(err)
	}
	later := runStart.Add(time.Hour)
	lab.log("t-planner-20261003T110000Z.log",
		streamCall("s", "m2", later, 999),
		streamFinal("Done.", 2000, 4.0))
	lab.planRun(story, "t-planner-20261003T110000Z.log")
	got, err := LogRunEnd(lab.dir, lab.root, "t", workitem.ActivityPlanner, nil)
	if err != nil || got == nil || got.Entry.Cost != 2.0 {
		t.Fatalf("resumed run end = %+v (%v), want its 2.0 USD share", got, err)
	}
	if s := lab.strategic(story); s == nil || s.Cost() != 4.0 {
		t.Errorf("story's planner usage = %+v, want the session's 4.0 USD once", s)
	}
}

// ADR-0083: a run no item's newest planner run kept its log in, such as one
// a person ran by hand or one replaced by a newer run for its item, charges
// nothing; nor does an activity outside any run, nor another kind's.
func TestAnActivityOfNoPlannedRunChargesNothing(t *testing.T) {
	lab := newActivityLab(t)
	epic, story := lab.planned()
	lab.log("t-planner-20261003T100000Z.log",
		streamCall("s", "m1", runStart, 999),
		streamFinal("Planned by hand.", 1000, 0.5))
	lab.planRun(story, "t-planner-20261003T090000Z.log")
	got, err := LogRunEnd(lab.dir, lab.root, "t", workitem.ActivityPlanner, nil)
	if err != nil || got == nil || got.Entry.Cost != 0.5 {
		t.Fatalf("run end = %+v (%v), want logged", got, err)
	}
	if got.Planned != "" || got.Charged != nil {
		t.Errorf("charged %s: %v, want nothing", got.Planned, got.Charged)
	}
	// the run logged to its end, then an activity outside any run
	lab.planRun(story, "t-planner-20261003T100000Z.log")
	got, err = LogActivity(lab.dir, lab.root, "t", workitem.ActivityPlanner, "planned by hand again", nil, runStart.Add(time.Hour))
	if err != nil || got.Charged != nil {
		t.Fatalf("outside a run = %+v (%v), want nothing charged", got, err)
	}
	// an orchestrator's run is no plan, whatever its log is called
	lab.log("t-orchestrator-20261003T100000Z.log", streamCall("o", "m1", runStart, 999), streamResult("o", 1000, 0.5))
	lab.planRun(epic, "t-orchestrator-20261003T100000Z.log")
	if got, err = LogRunEnd(lab.dir, lab.root, "t", workitem.ActivityOrchestrator, nil); err != nil || got == nil || got.Charged != nil {
		t.Fatalf("orchestrator = %+v (%v), want logged and nothing charged", got, err)
	}
	for _, id := range []string{story, epic} {
		if u := lab.usageOf(id); u != nil {
			t.Errorf("%s's usage = %+v, want none", id, u)
		}
	}
}

// usageOf is the item's usage as it stands.
func (lab *activityLab) usageOf(id string) *usage.Usage {
	lab.t.Helper()
	it, err := lab.repo.Get(id)
	if err != nil {
		lab.t.Fatal(err)
	}
	return it.Usage
}

// A run that died without its end logged charges a later activity only up to
// its last event, not the hours since.
func TestAnActivityEndsNoLaterThanItsRun(t *testing.T) {
	lab := newActivityLab(t)
	lab.log("t-orchestrator-20261003T100000Z.log",
		streamCall("s", "m1", runStart, 999),
		streamCall("s", "m2", runStart.Add(time.Minute), 999))
	got, err := LogActivity(lab.dir, lab.root, "t", workitem.ActivityOrchestrator, "Ordered the board", nil, runStart.Add(3*time.Hour))
	if err != nil {
		t.Fatal(err)
	}
	if e := got.Entry; e.Seconds != 61 || !e.Estimated || !e.At.Equal(runStart.Add(3*time.Hour)) {
		t.Errorf("entry = %+v, want 61 s to the second after the run's last event, ending when reported", e)
	}
}

// orchestrated makes two epics, a story under each, and a task under the
// first story, and archives the second story; it returns their IDs.
func (lab *activityLab) orchestrated() (epic, story, task, otherEpic, archived string) {
	lab.t.Helper()
	epic, story = lab.planned()
	otherEpic, archived = lab.planned()
	tk, err := lab.repo.Create(workitem.NewOptions{Type: workitem.Task, Title: "Task", Parent: story, Owner: "alex", Now: runStart})
	if err != nil {
		lab.t.Fatal(err)
	}
	it, err := lab.repo.Get(archived)
	if err != nil {
		lab.t.Fatal(err)
	}
	dir := lab.repo.ItemDir(workitem.Story, true)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		lab.t.Fatal(err)
	}
	if err := os.Rename(it.Path, filepath.Join(dir, filepath.Base(it.Path))); err != nil {
		lab.t.Fatal(err)
	}
	return epic, story, tk.ID, otherEpic, archived
}

// ADR-0095: an orchestrator activity's usage is split evenly between the
// work items it named that exist, archived ones included, each once, and
// each share is summed up to its epic; a thread, an issue, and an ID no item
// has take no share, and a story named with its own task carries both
// shares.
func TestAnOrchestratorActivityIsChargedEvenlyToTheWorkItemsItNamed(t *testing.T) {
	lab := newActivityLab(t)
	epic, story, task, otherEpic, archived := lab.orchestrated()
	lab.log("t-orchestrator-20261003T100000Z.log",
		streamCall("s", "m1", runStart, 999),
		streamCall("s", "m2", runStart.Add(time.Minute), 999),
		streamFinal("Accepted the stories.", 2000, 0.5))
	named := []string{task, archived, "TH-0186", "I-0007", "S-0099", story, strings.ToLower(task)}
	got, err := LogRunEnd(lab.dir, lab.root, "t", workitem.ActivityOrchestrator, named)
	if err != nil || got == nil {
		t.Fatalf("run end = %+v (%v), want logged", got, err)
	}
	if strings.Join(got.Entry.Items, ",") != strings.Join(named, ",") {
		t.Errorf("entry's items = %v, want all it named", got.Entry.Items)
	}
	if want := []string{task, archived, story}; strings.Join(got.Shared, ",") != strings.Join(want, ",") {
		t.Errorf("shared = %v, want %v, the work items named that exist, once each, in order", got.Shared, want)
	}
	if want := []string{task, story, epic, archived, otherEpic}; strings.Join(got.Charged, ",") != strings.Join(want, ",") {
		t.Errorf("charged = %v, want %v", got.Charged, want)
	}
	if got.Planned != "" {
		t.Errorf("planned = %q, want none for the orchestrator", got.Planned)
	}
	// 2 input, 300 output, and 2000 cache read tokens, and 61 seconds, in
	// three shares, the remainder to the first named
	if got.Entry.Seconds != 61 {
		t.Fatalf("entry's seconds = %d, want 61 to the second after the run's last event", got.Entry.Seconds)
	}
	secs := []int64{21, 20, 20}
	shareCost := toTheCentsHundredth(0.5 / 3)
	want := map[string]struct {
		tokens, seconds int64
		cost            float64
	}{
		task:      {1 + 100 + 667, secs[0], shareCost},
		archived:  {1 + 100 + 667, secs[1], shareCost},
		otherEpic: {1 + 100 + 667, secs[1], shareCost},
		story:     {(1 + 100 + 667) + (0 + 100 + 666), secs[0] + secs[2], toTheCentsHundredth(shareCost + 0.5/3)},
		epic:      {(1 + 100 + 667) + (0 + 100 + 666), secs[0] + secs[2], toTheCentsHundredth(shareCost + 0.5/3)},
	}
	for id, w := range want {
		s := lab.strategicOf(id, workitem.ActivityOrchestrator)
		if s == nil || s.Tokens() != w.tokens || s.Seconds != w.seconds || s.Cost() != w.cost || !s.Estimated {
			t.Errorf("%s's orchestrator usage = %+v, want %d tokens over %d s, %v USD, estimated", id, s, w.tokens, w.seconds, w.cost)
		}
		if p := lab.strategic(id); p != nil {
			t.Errorf("%s's planner usage = %+v, want none", id, p)
		}
	}
}

// ADR-0095: an orchestrator activity that names no work item charges no
// item; its cost is left to the orchestrator's project strategic total.
func TestAnOrchestratorActivityNamingNoWorkItemChargesNothing(t *testing.T) {
	lab := newActivityLab(t)
	epic, story := lab.planned()
	lab.log("t-orchestrator-20261003T100000Z.log",
		streamCall("s", "m1", runStart, 999),
		streamCall("s", "m2", runStart.Add(time.Minute), 999),
		streamFinal("Answered a thread.", 2000, 0.5))
	got, err := LogActivity(lab.dir, lab.root, "t", workitem.ActivityOrchestrator, "Answered TH-0186", []string{"TH-0186", "S-0099", "design/system/metrics.md"}, runStart.Add(30*time.Second))
	if err != nil {
		t.Fatal(err)
	}
	if got.Entry.Cost == 0 {
		t.Fatalf("entry = %+v, want a cost to leave uncharged", got.Entry)
	}
	if got.Shared != nil || got.Charged != nil || got.Planned != "" {
		t.Errorf("charged %v, %v, %q, want nothing", got.Shared, got.Charged, got.Planned)
	}
	for _, id := range []string{story, epic} {
		if u := lab.usageOf(id); u != nil {
			t.Errorf("%s's usage = %+v, want none", id, u)
		}
	}
}

// filed records an open issue for each title, as the analyzer files a
// finding, and returns their IDs.
func (lab *activityLab) filed(titles ...string) []string {
	lab.t.Helper()
	var ids []string
	for _, title := range titles {
		is, err := issues.New(lab.repo, issues.NewOptions{Title: title, Class: "defect", Now: runStart})
		if err != nil {
			lab.t.Fatal(err)
		}
		ids = append(ids, is.ID)
	}
	return ids
}

// issueUsage is the issue's usage as it stands.
func (lab *activityLab) issueUsage(id string) *usage.Usage {
	lab.t.Helper()
	is, err := issues.Get(lab.repo, id)
	if err != nil {
		lab.t.Fatal(err)
	}
	return is.Usage
}

// analyzed is what the analyzer spent on the issue, as its usage says; nil
// when nothing.
func (lab *activityLab) analyzed(id string) *usage.Strategic {
	lab.t.Helper()
	u := lab.issueUsage(id)
	if u == nil {
		return nil
	}
	for i, s := range u.Strategic {
		if s.Kind == workitem.ActivityAnalyzer {
			return &u.Strategic[i]
		}
	}
	return nil
}

// analyzerRun logs an analyzer run of two calls a minute apart that spent
// 0.5 USD over 2 input, 300 output, and 2000 cache read tokens.
func (lab *activityLab) analyzerRun() {
	lab.log("t-analyzer-20261003T100000Z.log",
		streamCall("s", "m1", runStart, 999),
		streamCall("s", "m2", runStart.Add(time.Minute), 999),
		streamFinal("Wrote design/analysis/2026-10-03-all.md.", 2000, 0.5))
}

// S-0227: an analyzer activity's usage is split evenly between the issues
// it named, each once at whatever padding, and charged to each under its
// analyzer entry; the work items and the thread it named take no share.
func TestAnAnalyzerActivityIsChargedEvenlyToTheIssuesItNamed(t *testing.T) {
	lab := newActivityLab(t)
	epic, story := lab.planned()
	ids := lab.filed("Slow review", "Flaky test")
	lab.analyzerRun()
	named := []string{story, "I-1", "TH-0001", epic, "i-0002", ids[0]}
	got, err := LogRunEnd(lab.dir, lab.root, "t", workitem.ActivityAnalyzer, named)
	if err != nil || got == nil {
		t.Fatalf("run end = %+v (%v), want logged", got, err)
	}
	if strings.Join(got.Shared, ",") != strings.Join(ids, ",") || strings.Join(got.Charged, ",") != strings.Join(ids, ",") || got.Planned != "" {
		t.Errorf("shared %v, charged %v, planned %q; want %v, the issues named, once each, in order", got.Shared, got.Charged, got.Planned, ids)
	}
	if got.Entry.Seconds != 61 || got.Entry.Cost != 0.5 {
		t.Fatalf("entry = %+v, want 61 s and 0.5 USD", got.Entry)
	}
	// 2 input, 300 output, and 2000 cache read tokens, and 61 seconds, in
	// two shares, the remainder to the first named
	for i, w := range []struct{ tokens, seconds int64 }{{1 + 150 + 1000, 31}, {1 + 150 + 1000, 30}} {
		s := lab.analyzed(ids[i])
		if s == nil || s.Tokens() != w.tokens || s.Seconds != w.seconds || s.Cost() != 0.25 || !s.Estimated {
			t.Errorf("%s's analyzer usage = %+v, want %d tokens over %d s, 0.25 USD, estimated", ids[i], s, w.tokens, w.seconds)
		}
		if u := lab.issueUsage(ids[i]); len(u.Strategic) != 1 || !u.Empty() {
			t.Errorf("%s's usage = %+v, want the analyzer's entry alone", ids[i], u)
		}
	}
	for _, id := range []string{story, epic} {
		if u := lab.usageOf(id); u != nil {
			t.Errorf("%s's usage = %+v, want none", id, u)
		}
	}
}

// S-0227, ADR-0095: an analyzer activity that names no issue charges no
// issue and no work item; its cost is left to the analyzer's project
// strategic total.
func TestAnAnalyzerActivityNamingNoIssueChargesNothing(t *testing.T) {
	lab := newActivityLab(t)
	epic, story := lab.planned()
	ids := lab.filed("Slow review")
	lab.analyzerRun()
	got, err := LogActivity(lab.dir, lab.root, "t", workitem.ActivityAnalyzer, "Read the metrics", []string{story, epic, "I-0009", "design/analysis/2026-10-03-all.md"}, runStart.Add(30*time.Second))
	if err != nil {
		t.Fatal(err)
	}
	if got.Entry.Cost == 0 {
		t.Fatalf("entry = %+v, want a cost to leave uncharged", got.Entry)
	}
	if got.Shared != nil || got.Charged != nil || got.Planned != "" {
		t.Errorf("charged %v, %v, %q, want nothing", got.Shared, got.Charged, got.Planned)
	}
	if u := lab.issueUsage(ids[0]); u != nil {
		t.Errorf("%s's usage = %+v, want none", ids[0], u)
	}
	for _, id := range []string{story, epic} {
		if u := lab.usageOf(id); u != nil {
			t.Errorf("%s's usage = %+v, want none", id, u)
		}
	}
}

// S-0227: an ID no issue has takes no share, so the issue named beside it
// is charged the whole.
func TestAnUnknownIssueTakesNoShareOfAnAnalyzerActivity(t *testing.T) {
	lab := newActivityLab(t)
	ids := lab.filed("Slow review")
	lab.analyzerRun()
	got, err := LogRunEnd(lab.dir, lab.root, "t", workitem.ActivityAnalyzer, []string{"I-0042", ids[0]})
	if err != nil || got == nil {
		t.Fatalf("run end = %+v (%v), want logged", got, err)
	}
	if strings.Join(got.Shared, ",") != ids[0] || strings.Join(got.Charged, ",") != ids[0] {
		t.Errorf("shared %v, charged %v, want %s alone", got.Shared, got.Charged, ids[0])
	}
	if s := lab.analyzed(ids[0]); s == nil || s.Tokens() != 2+300+2000 || s.Seconds != 61 || s.Cost() != 0.5 {
		t.Errorf("%s's analyzer usage = %+v, want the whole 2302 tokens over 61 s, 0.5 USD", ids[0], s)
	}
}
