package serve

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

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
	ended, err := LogRunEnd(lab.dir, lab.root, "t", workitem.ActivityOrchestrator)
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
	got, err := LogRunEnd(lab.dir, lab.root, "t", workitem.ActivityPlanner)
	if err != nil || got != nil {
		t.Fatalf("run end = %+v (%v), want nothing logged", got, err)
	}
	if doc := lab.doc(workitem.ActivityPlanner); doc.TasksCompleted != 1 {
		t.Errorf("document = %+v, want the one entry", doc)
	}
}

// A run that ends with calls since its last activity logs them as one, with
// the first line of the run's final text as its summary.
func TestARunEndLogsWhatWasSpentSinceTheLastActivity(t *testing.T) {
	lab := newActivityLab(t)
	lab.log("t-planner-20261003T100000Z.log",
		streamCall("s", "m1", runStart, 999),
		streamCall("s", "m2", runStart.Add(time.Minute), 2999),
		streamFinal("\n  Planned S-0001 and S-0002.  \nThe rest of the answer.", 4000, 1.0))
	if _, err := LogActivity(lab.dir, lab.root, "t", workitem.ActivityPlanner, "read the board", nil, runStart.Add(30*time.Second)); err != nil {
		t.Fatal(err)
	}
	got, err := LogRunEnd(lab.dir, lab.root, "t", workitem.ActivityPlanner)
	if err != nil || got == nil {
		t.Fatalf("run end = %+v (%v), want logged", got, err)
	}
	e := got.Entry
	if e.Summary != "Planned S-0001 and S-0002." || len(e.Items) != 0 || e.Seconds != 31 || e.Cost != 0.75 || !e.At.Equal(runStart.Add(61*time.Second)) {
		t.Errorf("entry = %+v, want the result's first line, 31 s, and m2's 0.75 USD", e)
	}
	if doc := lab.doc(workitem.ActivityPlanner); doc.TasksCompleted != 2 || doc.AccruedCost != 1.0 {
		t.Errorf("document = %+v, want both entries' totals", doc)
	}
}

// A run whose final text says nothing ends with a summary that says so.
func TestARunEndWithNoTextIsSaidToHaveEnded(t *testing.T) {
	lab := newActivityLab(t)
	lab.log("t-analyzer-20261003T100000Z.log", streamCall("s", "m1", runStart, 999), streamResult("s", 1000, 0.5))
	got, err := LogRunEnd(lab.dir, lab.root, "t", workitem.ActivityAnalyzer)
	if err != nil || got == nil || got.Entry.Summary != "run ended" || got.Entry.Cost != 0.5 {
		t.Fatalf("run end = %+v (%v), want run ended at 0.5 USD", got, err)
	}
}

// An activity outside any run flai serve logged is still recorded, with no
// seconds and no cost; a run end with no run logs nothing.
func TestAnActivityWithNoRunIsLoggedWithNothingMeasured(t *testing.T) {
	lab := newActivityLab(t)
	if got, err := LogRunEnd(lab.dir, lab.root, "t", workitem.ActivityPlanner); err != nil || got != nil {
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
	if _, err := LogRunEnd(lab.dir, lab.root, "t", workitem.ActivityPlanner); err != nil {
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
	if _, err := LogRunEnd(lab.dir, lab.root, "t", "../stories"); err == nil {
		t.Error("ending a run of ../stories was not refused")
	}
	if m, _ := filepath.Glob(filepath.Join(lab.root, "wip", "agents", "*.md")); len(m) != 0 {
		t.Errorf("wrote %v", m)
	}
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
