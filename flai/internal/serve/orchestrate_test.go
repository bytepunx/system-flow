//go:build !windows

package serve

import (
	"context"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/hostapi"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// orchestrateLab is an agentLab whose command is an orchestrator stub, with
// the project's orchestrator on the lab's launcher and a clock of its own,
// later than the launcher's by shift. The stub prints stream-orchestrator as
// its output when there is one, then writes what it was given to
// orchestrator-<pid>.txt, so that once given has read it, the run's usage is
// in its log and a stop cannot come before it (I-0090). It then waits while
// held until released, marks its exit, as it would have ended, when asked to
// stop, and exits with the code in exit-orchestrator when there is one. The
// orchestrate action is off.
type orchestrateLab struct {
	*agentLab
	orch  *orchestrator
	shift time.Duration
}

func newOrchestrateLab(t *testing.T) *orchestrateLab {
	t.Helper()
	lab := &orchestrateLab{agentLab: newAgentLab(t)}
	d := lab.outDir
	stub := filepath.Join(d, "orch-agent")
	script := "#!/bin/sh\n" +
		"trap 'touch \"" + d + "/ended-$$\"; exit 143' TERM\n" +
		"[ -f \"" + d + "/stream-orchestrator\" ] && cat \"" + d + "/stream-orchestrator\"\n" +
		"{ echo \"args: $*\"; echo \"dir: $(pwd)\"; echo \"agent: $FLAI_AGENT\"; echo \"story: ${FLAI_STORY-unset}\"; echo \"role: $FLAI_ROLE\"; echo \"item: ${FLAI_ITEM-unset}\"; echo \"session: $FLAI_SESSION\"; echo \"by: $FLAI_STARTED_BY\"; } > \"" + d + "/orchestrator-$$.txt\"\n" +
		"while [ -f \"" + d + "/hold\" ] && [ ! -f \"" + d + "/release-orchestrator\" ]; do sleep 0.05; done\n" +
		"touch \"" + d + "/ended-$$\"\n" +
		"[ -f \"" + d + "/exit-orchestrator\" ] && exit \"$(cat \"" + d + "/exit-orchestrator\")\"\n" +
		"exit 0\n"
	if err := os.WriteFile(stub, []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	lab.cfg.Command = []string{stub, "story={story}"}
	o := lab.o
	o.Now = func() time.Time { return time.Now().Add(lab.shift) }
	lab.orch = newOrchestrator(o, Entry{Key: "t", Name: "t", Root: lab.root}, lab.l)
	return lab
}

func (lab *orchestrateLab) look() { lab.orch.look(context.Background()) }

func (lab *orchestrateLab) orchestrator() *AgentRun { return lab.state().Orchestrator }

// given is what the stub was given in the run with pid, once it has written
// it.
func (lab *orchestrateLab) given(pid int) string {
	lab.t.Helper()
	out := filepath.Join(lab.outDir, "orchestrator-"+strconv.Itoa(pid)+".txt")
	waitFor(lab.t, "the orchestrator stub ran", func() bool {
		data, _ := os.ReadFile(out)
		return strings.Contains(string(data), "by: ")
	})
	data, _ := os.ReadFile(out)
	return string(data)
}

// ended waits until the orchestrator run with session has ended, and
// returns it.
func (lab *orchestrateLab) ended(session string) *AgentRun {
	lab.t.Helper()
	waitFor(lab.t, "the orchestrator run ended", func() bool {
		r := lab.orchestrator()
		return r != nil && r.Session == session && r.Ended != ""
	})
	return lab.orchestrator()
}

// orchestrateEntries are the journal's entries for the orchestrate action.
func (lab *orchestrateLab) orchestrateEntries() []hostapi.Entry {
	var out []hostapi.Entry
	for _, e := range lab.entries() {
		if e.Action == hostapi.ActionOrchestrate {
			out = append(out, e)
		}
	}
	return out
}

// S-0218: with orchestrate off nothing starts; turned on, one orchestrator
// run starts for the project, in its main checkout, as the orchestrator and
// for no story or item, logged as the orchestrator's and recorded under
// orchestrator, apart from the stories' runs; while it runs, no other is
// started, and it holds no story's agent back.
func TestTheOrchestratorRunsWhileTheActionIsOn(t *testing.T) {
	t.Setenv("FLAI_STORY", "S-0999") // the session flai serve was started from
	lab := newOrchestrateLab(t)
	lab.hold()
	lab.look()
	if st := lab.state(); st.Orchestrator != nil || len(lab.orchestrateEntries()) != 0 {
		t.Fatalf("started while orchestrate is off: %+v", st.Orchestrator)
	}
	lab.cfg.Orchestrate = true
	lab.look()
	run := lab.orchestrator()
	if run == nil || !run.live() || run.Story != "" || run.Item != "" || run.Agent != "orchestrator" || run.Harness != "command" || run.Command != "orch-agent" || run.Session == "" {
		t.Fatalf("run = %+v", run)
	}
	logs, err := lab.o.Dir.ActivityLogs("t", workitem.ActivityOrchestrator)
	if err != nil || len(logs) != 1 || logs[0] != run.Log || !strings.HasPrefix(filepath.Base(run.Log), "t-orchestrator-") {
		t.Errorf("orchestrator logs = %v (%v), want the run's %s", logs, err, run.Log)
	}
	got := lab.given(run.PID)
	real, _ := filepath.EvalSymlinks(lab.root)
	for _, want := range []string{"args: story=\n", "dir: " + real + "\n", "agent: orchestrator\n", "story: unset\n", "role: orchestrate\n", "item: unset\n", "session: 2", "by: flai-serve\n"} {
		if !strings.Contains(got, want) {
			t.Errorf("the orchestrator was given\n%s\nwithout %q", got, want)
		}
	}
	if data, err := os.ReadFile(lab.o.Dir.agents()); err != nil || !strings.Contains(string(data), `"orchestrator": {`) {
		t.Errorf("serve/agents.json does not record the run under orchestrator (%v):\n%s", err, data)
	}
	if st := lab.state(); st.Running != nil || st.Last != nil || len(st.Stories) != 0 || len(st.Plans) != 0 {
		t.Errorf("state = %+v, want the run under orchestrator alone", st)
	}
	j := lab.orchestrateEntries()
	if len(j) != 1 || j[0].Method != "serve.orchestrate" || j[0].Outcome != "done" || !strings.Contains(j[0].Detail, "started orch-agent (command) as orchestrator (pid ") {
		t.Errorf("journal: %+v", j)
	}
	lab.look()
	if again := lab.orchestrator(); !again.same(run) || len(lab.orchestrateEntries()) != 1 {
		t.Errorf("a second run was started beside the first: %+v", again)
	}
	// the in-progress limit does not count it
	lab.limit(1)
	story := lab.ready("Beside")
	lab.l.look(context.Background(), false)
	waitFor(t, "the story's agent started beside the orchestrator", func() bool { return lab.run(story).live() })
	lab.release(story)
	lab.release("orchestrator")
}

// S-0218: a run that ends while the action is on is started again at the
// next look, and, when it failed, no sooner than a minute after it ended;
// each end is journalled.
func TestTheOrchestratorIsStartedAgainWhenItEnds(t *testing.T) {
	lab := newOrchestrateLab(t)
	lab.hold()
	lab.cfg.Orchestrate = true
	lab.look()
	first := lab.orchestrator()
	lab.release("orchestrator")
	if r := lab.ended(first.Session); r.Outcome != OutcomeWorked || r.Exit == nil || *r.Exit != 0 {
		t.Fatalf("first = %+v, want worked with exit 0", r)
	}
	if err := os.WriteFile(filepath.Join(lab.outDir, "exit-orchestrator"), []byte("3"), 0o644); err != nil {
		t.Fatal(err)
	}
	lab.look()
	second := lab.orchestrator()
	if second.Session == first.Session || second.Started == "" {
		t.Fatalf("not started again after it worked: %+v", second)
	}
	failed := lab.ended(second.Session)
	if failed.Outcome != OutcomeFailed || failed.Exit == nil || *failed.Exit != 3 || failed.Why != "ended (exit 3)" {
		t.Fatalf("second = %+v, want failed with exit 3", failed)
	}
	lab.look()
	// well short of the retry: the lab's clock is the host's plus the shift, and a loaded host
	// adds seconds of its own between the failure and this look (I-0102)
	lab.shift = orchestrateRetry - 20*time.Second
	lab.look()
	if r := lab.orchestrator(); r.Session != second.Session {
		t.Fatalf("started again within a minute of a failure: %+v", r)
	}
	lab.shift = orchestrateRetry + time.Second
	lab.look()
	third := lab.orchestrator()
	if third.Session == second.Session {
		t.Fatalf("not started again a minute after a failure: %+v", third)
	}
	lab.ended(third.Session)
	var ends []hostapi.Entry
	for _, e := range lab.orchestrateEntries() {
		if strings.Contains(e.Detail, ") ended, ") {
			ends = append(ends, e)
		}
	}
	if len(ends) != 3 || ends[0].Outcome != "done" || !strings.HasSuffix(ends[0].Detail, "ended, worked") ||
		ends[1].Outcome != "failed" || !strings.HasSuffix(ends[1].Detail, "ended, failed: ended (exit 3)") {
		t.Errorf("ends journalled: %+v", ends)
	}
}

// S-0218: turning the action off stops the run going, as the operator's stop
// of a story's agent does, whoever waits for it, journals the stop, and logs
// the run's activity in orchestrator.md as stopped for the action.
func TestTheOrchestratorIsStoppedWhenTheActionIsTurnedOff(t *testing.T) {
	for name, handOver := range map[string]bool{"a run this flai serve waits for": false, "a run an earlier flai serve started": true} {
		t.Run(name, func(t *testing.T) {
			lab := newOrchestrateLab(t)
			lab.hold()
			stream := strings.Join([]string{streamCall("s", "m1", runStart, 999), streamCall("s", "m2", runStart.Add(time.Minute), 999)}, "\n") + "\n"
			if err := os.WriteFile(filepath.Join(lab.outDir, "stream-orchestrator"), []byte(stream), 0o644); err != nil {
				t.Fatal(err)
			}
			lab.l.handOver = handOver
			lab.cfg.Orchestrate = true
			lab.look()
			run := lab.orchestrator()
			if handOver {
				// the process an earlier flai serve left is reaped by the
				// system once it ends, not by this one
				go func() {
					if p, err := os.FindProcess(run.PID); err == nil {
						_, _ = p.Wait()
					}
				}()
			}
			lab.given(run.PID) // it runs, has set its trap, and has printed its usage
			lab.cfg.Orchestrate = false
			lab.look()
			if Alive(run.PID) {
				t.Errorf("pid %d still runs", run.PID)
			}
			stopped := lab.ended(run.Session)
			if stopped.Outcome != OutcomeStopped || stopped.Stopped == "" || stopped.live() {
				t.Errorf("not recorded as stopped: %+v", stopped)
			}
			waitFor(t, "its activity logged", func() bool {
				doc, err := lab.repo.Activity(workitem.ActivityOrchestrator)
				return err == nil && len(doc.Entries) == 1
			})
			doc, _ := lab.repo.Activity(workitem.ActivityOrchestrator)
			if e := doc.Entries[0]; e.Summary != orchestrateOff || len(e.Items) != 0 || e.Seconds != 61 {
				t.Errorf("activity = %+v, want %q over the run's 61 s and no items", e, orchestrateOff)
			}
			j := lab.orchestrateEntries()
			if last := j[len(j)-1]; last.Outcome != "done" || !strings.Contains(last.Detail, "stopped the orchestrator orchestrator, since orchestrate was turned off: ended its process group (pid "+strconv.Itoa(run.PID)+")") {
				t.Errorf("journal: %+v", j)
			}
			// off, it is not started again
			lab.look()
			if r := lab.orchestrator(); r.Session != run.Session {
				t.Errorf("started while off: %+v", r)
			}
		})
	}
}

// S-0218: the orchestrator is refused, saying why, when its agent names no
// harness and no command is set, and when the harness refuses it, as
// claude-code does without orchestrator.md; each refusal is recorded and
// journalled once, not at every look.
func TestTheOrchestratorIsRefusedOnce(t *testing.T) {
	lab := newOrchestrateLab(t)
	lab.cfg.Command, lab.cfg.Orchestrate = nil, true
	lab.look()
	run := lab.orchestrator()
	if run == nil || run.Outcome != OutcomeFailed || !strings.Contains(run.Error, "the orchestrator's agent names no harness (orchestration.agent or agent in system-flow.yaml), and no command is set on the host") {
		t.Fatalf("no harness: %+v", run)
	}
	for _, shift := range []time.Duration{0, 2 * orchestrateRetry} {
		lab.shift = shift
		lab.look()
	}
	if j := lab.orchestrateEntries(); len(j) != 1 || j[0].Outcome != "failed" || !strings.Contains(j[0].Detail, "the orchestrator could not be started: the orchestrator's agent names no harness") {
		t.Errorf("journal: %+v", j)
	}
	if r := lab.orchestrator(); r.Started != run.Started {
		t.Errorf("the refusal was recorded again: %+v", r)
	}
	manifest := "version: 1\nname: t\nkey: t\nlayout:\n  design: design\n  docs: docs\n  wip: wip\norchestration:\n  agent:\n    harness: claude-code\n"
	if err := os.WriteFile(filepath.Join(lab.root, "system-flow.yaml"), []byte(manifest), 0o644); err != nil {
		t.Fatal(err)
	}
	lab.shift = 4 * orchestrateRetry
	lab.look()
	if r := lab.orchestrator(); r.Harness != "claude-code" || !strings.Contains(r.Error, "the orchestrator's session runs as the agent orchestrator, and its definition could not be read") {
		t.Errorf("no orchestrator.md: %+v", r)
	}
	lab.shift = 6 * orchestrateRetry
	lab.look()
	if j := lab.orchestrateEntries(); len(j) != 2 {
		t.Errorf("journal: %+v", j)
	}
	lab.logs.mu.Lock()
	logged := lab.logs.buf.String()
	lab.logs.mu.Unlock()
	if n := strings.Count(logged, `msg="orchestrator could not be started"`); n != 2 {
		t.Errorf("the refusals were logged %d times, want once each:\n%s", n, logged)
	}
}
