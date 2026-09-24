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

	"github.com/bytepunx/system-flow/flai/internal/threads"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// S-0115: the operator starts a ready story's agent now, whatever the
// launcher's own rules say about when, and is told why when it will not.
func TestAReadyStorysAgentIsStartedOnTheOperatorsWord(t *testing.T) {
	ctx := context.Background()
	start := func(lab *agentLab, id string) (*AgentRun, error) {
		return Start(ctx, lab.o, Entry{Key: "t", Name: "t", Root: lab.root}, id)
	}
	refusedFor := func(t *testing.T, lab *agentLab, id, want string) {
		t.Helper()
		_, err := start(lab, id)
		var no *Refused
		if !errors.As(err, &no) || !strings.Contains(no.Why, want) {
			t.Errorf("start %s: %v, want refused for %q", id, err, want)
		}
	}
	t.Run("a story the launcher will not start again", func(t *testing.T) {
		lab := newAgentLab(t)
		id := lab.ready("A")
		lab.l.look(ctx, false)
		waitFor(t, "its agent ended", func() bool { r := lab.run(id); return r != nil && r.Ended != "" })
		first := lab.run(id)
		lab.l.look(ctx, false)
		if r := lab.run(id); r.Started != first.Started || r.Session != first.Session {
			t.Fatalf("the launcher started it again: %+v", r)
		}
		run, err := start(lab, id)
		if err != nil {
			t.Fatal(err)
		}
		if run.Session == first.Session || run.Agent != "builder-"+id || run.PID == 0 {
			t.Errorf("a new session for the same agent: %+v, was %+v", run, first)
		}
		waitFor(t, "it ran", func() bool {
			data, _ := os.ReadFile(filepath.Join(lab.outDir, id+".txt"))
			return strings.Contains(string(data), "agent: builder-"+id)
		})
		if j := lab.entries(); len(j) != 2 || !strings.Contains(j[1].Detail, "started stub-agent (command) for "+id) {
			t.Errorf("journal: %+v", j)
		}
	})
	t.Run("a story no launcher has looked at, with someone attending", func(t *testing.T) {
		lab := newAgentLab(t)
		lab.attend()
		id := lab.ready("A") // no look: flai serve is not running, or has not seen it
		run, err := start(lab, id)
		if err != nil {
			t.Fatal(err)
		}
		if got := lab.run(id); got == nil || got.Started != run.Started || got.PID != run.PID {
			t.Errorf("the run is not recorded where flai serve tracks it: %+v", got)
		}
	})
	t.Run("past a full in-progress limit, with a warning", func(t *testing.T) {
		lab := newAgentLab(t)
		lab.hold()
		lab.limit(1)
		busy := lab.ready("Busy")
		lab.move(busy, workitem.InProgress)
		id := lab.ready("Waits")
		lab.l.look(ctx, false)
		if r := lab.run(id); r != nil {
			t.Fatalf("the launcher started it with the limit full: %+v", r)
		}
		run, err := start(lab, id)
		if err != nil {
			t.Fatal(err)
		}
		if !run.live() {
			t.Errorf("not started: %+v", run)
		}
		lab.logs.mu.Lock()
		logged := lab.logs.buf.String()
		lab.logs.mu.Unlock()
		if !strings.Contains(logged, `level=WARN msg="agent started past the in-progress limit"`) || !strings.Contains(logged, "story="+id) {
			t.Errorf("no warning: %s", logged)
		}
		lab.release(id)
	})
	t.Run("the serving flai settles it and starts it again on an answer", func(t *testing.T) {
		lab := newAgentLab(t)
		lab.hold()
		id := lab.ready("Asks")
		run, err := start(lab, id)
		if err != nil {
			t.Fatal(err)
		}
		if lab.l.waiting[run.PID] {
			t.Fatal("the serving launcher waits for a process it did not start")
		}
		lab.move(id, workitem.InProgress) // the agent pulled it
		th, err := threads.New(lab.repo, threads.NewOptions{Title: "Which port?", On: id, Author: run.Agent, Text: "Eight or nine?", Now: time.Now()})
		if err != nil {
			t.Fatal(err)
		}
		lab.release(id)
		// The command that started it has exited, and the process is gone:
		// in this test the process stays this test's unreaped child, so a
		// reaped one stands in for it.
		gone := exec.Command("true")
		_ = gone.Run()
		lab.l.dir.updateAgent(lab.root, func(s *AgentState) {
			r := *s.Stories[id]
			r.PID = gone.Process.Pid
			s.put(&r)
		})
		_ = os.Remove(filepath.Join(lab.outDir, "release-"+id))
		lab.l.look(ctx, false)
		if r := lab.run(id); r.Ended == "" || r.Outcome != OutcomeAsked || r.Thread != th.ID {
			t.Fatalf("not settled as asking: %+v", r)
		}
		if _, err := threads.Reply(lab.repo, th.ID, "alex", "Nine.", time.Now()); err != nil {
			t.Fatal(err)
		}
		lab.l.look(ctx, false)
		waitFor(t, "it runs again", func() bool { return lab.run(id).live() })
		if again := lab.run(id); again.Agent != run.Agent || again.Session != run.Session || again.Answered != th.ID {
			t.Errorf("the same agent, in its session, for the answer: %+v, first %+v", again, run)
		}
		lab.release(id)
	})
	t.Run("refusals", func(t *testing.T) {
		lab := newAgentLab(t)
		lab.hold()
		backlog := lab.backlog("Backlog", nil)
		refusedFor(t, lab, backlog, "is in backlog; only a story in ready is started")
		// in progress: restart is the way
		busy := lab.ready("Busy")
		lab.move(busy, workitem.InProgress)
		refusedFor(t, lab, busy, "flai serve agent restart starts a new agent for one in progress")
		// running
		running := lab.ready("Running")
		lab.l.look(ctx, false)
		waitFor(t, "it runs", func() bool { return lab.run(running).live() })
		refusedFor(t, lab, running, "agent is running")
		other := lab.ready("Other")
		lab.release(running)
		waitFor(t, "it ends", func() bool { return !lab.run(running).live() })
		// asked
		lab.limit(5)
		lab.l.dir.updateAgent(lab.root, func(s *AgentState) {
			r := *s.Stories[running]
			r.Outcome, r.Thread = OutcomeAsked, "TH-0001"
			s.put(&r)
		})
		refusedFor(t, lab, running, "waiting for an answer to TH-0001")
		// nothing to start it with
		lab.cfg.Command = nil
		refusedFor(t, lab, other, "names no harness, and no command is set")
		refusedFor(t, lab, "T-0001", "is not a story's ID")
		// the action off
		lab.cfg.Enabled = false
		refusedFor(t, lab, other, "agent host action is off")
		if r := lab.run(other); r != nil {
			t.Errorf("a refused start left a run: %+v", r)
		}
	})
}
