//go:build !windows

package serve

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// told collects the runs the lab's launcher tells the dashboard of, as flai
// serve's changed does.
type told struct {
	mu   sync.Mutex
	runs []AgentRun
}

func (tl *told) of(lab *orchestrateLab) *told {
	lab.l.changed = func(run *AgentRun) {
		tl.mu.Lock()
		defer tl.mu.Unlock()
		tl.runs = append(tl.runs, *run)
	}
	return tl
}

// since are the runs told from the nth on.
func (tl *told) since(n int) []AgentRun {
	tl.mu.Lock()
	defer tl.mu.Unlock()
	return append([]AgentRun{}, tl.runs[min(n, len(tl.runs)):]...)
}

func (tl *told) count() int {
	tl.mu.Lock()
	defer tl.mu.Unlock()
	return len(tl.runs)
}

// orchestrating has the lab's stub print usage over a minute, so that a run's
// end logs its activity, turns the action on, and starts the orchestrator,
// whose run it returns once the stub has printed its usage.
func (lab *orchestrateLab) orchestrating() *AgentRun {
	lab.t.Helper()
	lab.hold() // the stub waits until released, or asked to stop
	stream := strings.Join([]string{streamCall("s", "m1", runStart, 999), streamCall("s", "m2", runStart.Add(time.Minute), 999)}, "\n") + "\n"
	if err := os.WriteFile(filepath.Join(lab.outDir, "stream-orchestrator"), []byte(stream), 0o644); err != nil {
		lab.t.Fatal(err)
	}
	lab.cfg.Orchestrate = true
	lab.look()
	run := lab.orchestrator()
	if !run.live() {
		lab.t.Fatalf("the orchestrator did not start: %+v", run)
	}
	lab.given(run.PID)
	return run
}

// reap waits for pid, a process the lab started and let go, as the system
// would for a process an earlier flai serve or a command left.
func reap(pid int) {
	go func() {
		if p, err := os.FindProcess(pid); err == nil {
			_, _ = p.Wait()
		}
	}()
}

// activity is the orchestrator's activity document's entries once there are
// n of them.
func (lab *orchestrateLab) activity(n int) []workitem.ActivityEntry {
	lab.t.Helper()
	waitFor(lab.t, "the orchestrator's activity logged", func() bool {
		doc, err := lab.repo.Activity(workitem.ActivityOrchestrator)
		return err == nil && len(doc.Entries) == n
	})
	doc, _ := lab.repo.Activity(workitem.ActivityOrchestrator)
	return doc.Entries
}

func refusal(err error) string {
	var no *Refused
	if errors.As(err, &no) {
		return no.Why
	}
	return ""
}

// S-0228: the operator's stop ends the run flai serve waits for as its stop
// of a story's agent does, holds the orchestrator stopped, logs the run's
// activity as stopped from the dashboard, journals the stop, and the
// dashboard is told of the end; no look starts it again while the action
// stays on, and a second stop is refused.
func TestTheOperatorStopsTheOrchestratorAndItStaysStopped(t *testing.T) {
	lab := newOrchestrateLab(t)
	tl := (&told{}).of(lab)
	if _, err := OrchestratorStop(lab.o, lab.entry()); !strings.Contains(refusal(err), "the orchestrate host action is off for this project, so flai serve runs no orchestrator to stop: flai serve enable orchestrate") {
		t.Fatalf("off: %v", err)
	}
	lab.cfg.Orchestrate = true
	if _, err := OrchestratorStop(lab.o, lab.entry()); refusal(err) != "flai serve has started no orchestrator for this project, so there is none to stop" {
		t.Fatalf("no run: %v", err)
	}
	run := lab.orchestrating()
	before := tl.count()
	got, err := OrchestratorStop(lab.o, lab.entry())
	if err != nil {
		t.Fatal(err)
	}
	if Alive(run.PID) {
		t.Errorf("pid %d still runs", run.PID)
	}
	if got.Session != run.Session || !got.Held || got.Outcome != OutcomeStopped || got.Stopped == "" || got.live() {
		t.Errorf("answered %+v, want the run ended stopped and held", got)
	}
	if r := lab.orchestrator(); !r.Held || r.Outcome != OutcomeStopped {
		t.Errorf("recorded %+v, want it held and stopped", r)
	}
	if e := lab.activity(1)[0]; e.Summary != "stopped: stopped from the dashboard" || e.Seconds != 61 {
		t.Errorf("activity = %+v, want %q over the run's 61 s", e, orchestrateHeld)
	}
	j := lab.orchestrateEntries()
	if last := j[len(j)-1]; last.Method != "serve.orchestrate" || last.Outcome != "done" ||
		last.Detail != "stopped the orchestrator orchestrator on the operator's word, and held it stopped while orchestrate stays on: ended its process group (pid "+strconv.Itoa(run.PID)+")" {
		t.Errorf("journal: %+v", j)
	}
	if ends := tl.since(before); len(ends) == 0 || ends[len(ends)-1].Session != run.Session || ends[len(ends)-1].Ended == "" || !ends[len(ends)-1].Held {
		t.Errorf("the dashboard was told %+v, want the held run's end", ends)
	}
	for _, shift := range []time.Duration{0, 2 * orchestrateRetry} {
		lab.shift = shift
		lab.look()
		if r := lab.orchestrator(); r.Session != run.Session {
			t.Fatalf("started again while held: %+v", r)
		}
	}
	if _, err := OrchestratorStop(lab.o, lab.entry()); refusal(err) != "the orchestrator is already held stopped: flai serve orchestrate start starts it again" {
		t.Errorf("a second stop: %v", err)
	}
}

// S-0228: the operator's start lifts the hold and starts the orchestrator as
// turning the action on does, handed over to the serving flai, which starts
// no other beside it and tells the dashboard of it at its next look; it is
// refused while the action is off and while the orchestrator is not held.
func TestTheOperatorStartsTheHeldOrchestrator(t *testing.T) {
	lab := newOrchestrateLab(t)
	tl := (&told{}).of(lab)
	if _, err := OrchestratorStart(context.Background(), lab.o, lab.entry()); !strings.Contains(refusal(err), "the orchestrate host action is off for this project, so flai serve runs no orchestrator to start") {
		t.Fatalf("off: %v", err)
	}
	run := lab.orchestrating()
	if _, err := OrchestratorStart(context.Background(), lab.o, lab.entry()); !strings.HasPrefix(refusal(err), "the orchestrator is not held stopped: flai serve runs it") {
		t.Fatalf("not held: %v", err)
	}
	if _, err := OrchestratorStop(lab.o, lab.entry()); err != nil {
		t.Fatal(err)
	}
	lab.look() // the serving flai has seen it held
	before := tl.count()
	got, err := OrchestratorStart(context.Background(), lab.o, lab.entry())
	if err != nil {
		t.Fatal(err)
	}
	reap(got.PID)
	if got.Session == run.Session || got.Held || !got.live() || got.Agent != "orchestrator" || got.Command != "orch-agent" {
		t.Fatalf("started %+v, want a new run, not held", got)
	}
	if r := lab.orchestrator(); r.Session != got.Session || r.Held {
		t.Errorf("recorded %+v, want the new run, not held", r)
	}
	if given := lab.given(got.PID); !strings.Contains(given, "role: orchestrate\n") || !strings.Contains(given, "agent: orchestrator\n") {
		t.Errorf("the orchestrator was given\n%s", given)
	}
	lab.l.mu.Lock()
	waited := lab.l.waiting[got.PID]
	lab.l.mu.Unlock()
	if waited {
		t.Error("the serving flai waits for a run a command started; it settles it once gone")
	}
	lab.look()
	if r := lab.orchestrator(); r.Session != got.Session {
		t.Errorf("another run was started beside the operator's: %+v", r)
	}
	if news := tl.since(before); len(news) != 1 || news[0].Session != got.Session || news[0].Held {
		t.Errorf("the dashboard was told %+v, want the started run once", news)
	}
	j := lab.orchestrateEntries()
	if last := j[len(j)-1]; last.Outcome != "done" || !strings.Contains(last.Detail, "started orch-agent (command) as orchestrator (pid "+strconv.Itoa(got.PID)+")") {
		t.Errorf("journal: %+v", j)
	}
}

// S-0228: turning the action off lifts the hold and tells the dashboard, so
// that turned on again the orchestrator starts.
func TestTurningOrchestrateOffAndOnLiftsTheHold(t *testing.T) {
	lab := newOrchestrateLab(t)
	tl := (&told{}).of(lab)
	run := lab.orchestrating()
	if _, err := OrchestratorStop(lab.o, lab.entry()); err != nil {
		t.Fatal(err)
	}
	lab.look()
	before := tl.count()
	lab.cfg.Orchestrate = false
	lab.look()
	r := lab.orchestrator()
	if r.Held || r.Session != run.Session {
		t.Fatalf("turned off: %+v, want the run no longer held", r)
	}
	if news := tl.since(before); len(news) != 1 || news[0].Held {
		t.Errorf("the dashboard was told %+v, want the hold lifted once", news)
	}
	lab.look()
	if n := tl.count(); n != before+1 {
		t.Errorf("told again at a look that changed nothing: %+v", tl.since(before))
	}
	lab.cfg.Orchestrate = true
	lab.look()
	if again := lab.orchestrator(); again.Session == run.Session || !again.live() {
		t.Errorf("turned on again: %+v, want a new run", again)
	}
}

// S-0228: a run no flai serve waits for, one a command started, is recorded
// as ended by the stop itself, with its activity logged as stopped from the
// dashboard; and a run that failed, waiting to be started again, is held as
// it is, and not started again a minute later.
func TestTheOperatorStopsARunNobodyWaitsForAndAFailedOne(t *testing.T) {
	was := orchestrateSettle
	orchestrateSettle = 200 * time.Millisecond
	t.Cleanup(func() { orchestrateSettle = was })
	lab := newOrchestrateLab(t)
	lab.l.handOver = true
	run := lab.orchestrating()
	reap(run.PID)
	got, err := OrchestratorStop(lab.o, lab.entry())
	if err != nil {
		t.Fatal(err)
	}
	if got.Session != run.Session || got.Outcome != OutcomeStopped || !got.Held || got.Ended == "" {
		t.Errorf("answered %+v, want the run ended stopped and held", got)
	}
	if e := lab.activity(1)[0]; e.Summary != orchestrateHeld {
		t.Errorf("activity = %+v, want %q", e, orchestrateHeld)
	}

	lab = newOrchestrateLab(t)
	lab.cfg.Command, lab.cfg.Orchestrate = nil, true
	lab.look()
	failed := lab.orchestrator()
	if failed == nil || failed.Outcome != OutcomeFailed {
		t.Fatalf("no harness: %+v", failed)
	}
	got, err = OrchestratorStop(lab.o, lab.entry())
	if err != nil {
		t.Fatal(err)
	}
	if !got.Held || got.Session != failed.Session || got.Outcome != OutcomeFailed || got.Stopped != "" {
		t.Errorf("answered %+v, want the failed run held as it is", got)
	}
	j := lab.orchestrateEntries()
	if last := j[len(j)-1]; !strings.HasSuffix(last.Detail, "held it stopped while orchestrate stays on: it was not running: it ended "+failed.Ended+", failed: "+failed.Why) {
		t.Errorf("journal: %+v", last)
	}
	lab.shift = 2 * orchestrateRetry
	lab.look()
	if r := lab.orchestrator(); r.Session != failed.Session {
		t.Errorf("started again while held: %+v", r)
	}
}

func (lab *orchestrateLab) entry() Entry { return Entry{Key: "t", Name: "t", Root: lab.root} }
