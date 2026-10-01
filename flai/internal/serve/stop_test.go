//go:build !windows

package serve

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/threads"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// trapped waits until the story's stub runs and has set its trap for TERM,
// which it says by writing its output file. A process that is only running
// may not have read its trap yet, and dies of the signal instead of handling
// it: on macOS, whose /bin/sh is bash, the gap is wide enough that the test
// stopped every stub before it had (I-0045).
func (lab *agentLab) trapped(t *testing.T, id string) {
	t.Helper()
	waitFor(t, "it runs and has set its trap", func() bool {
		_, err := os.Stat(filepath.Join(lab.outDir, id+".txt"))
		return lab.run(id).live() && err == nil
	})
}

// S-0170: the operator stops a story's agent. What runs is ended with what it
// started, a process that is no longer the agent is left alone, an agent
// waiting for an answer is not started again, and the run ends stopped
// whoever records its end.
func TestTheOperatorStopsAStorysAgent(t *testing.T) {
	ctx := context.Background()
	stop := func(lab *agentLab, id string) (*AgentRun, error) {
		return Stop(lab.o, Entry{Key: "t", Name: "t", Root: lab.root}, id)
	}
	stoppedAs := func(t *testing.T, r *AgentRun) {
		t.Helper()
		if r == nil || r.Ended == "" || r.Stopped == "" || r.Outcome != OutcomeStopped || !strings.Contains(r.Why, "stopped by the operator at "+r.Stopped) || r.live() {
			t.Errorf("not recorded as stopped: %+v", r)
		}
	}
	t.Run("a running agent is ended, and flai serve keeps the stop when it sees the end", func(t *testing.T) {
		lab := newAgentLab(t)
		lab.hold()
		id := lab.ready("Runs")
		lab.l.look(ctx, false)
		lab.trapped(t, id)
		started := lab.run(id)
		lab.move(id, workitem.InProgress) // the agent pulled it
		run, err := stop(lab, id)
		if err != nil {
			t.Fatal(err)
		}
		stoppedAs(t, run)
		if Alive(started.PID) {
			t.Errorf("pid %d still runs", started.PID)
		}
		// the launcher that started it records the end it saw, as a stop
		waitFor(t, "the launcher saw it end", func() bool {
			lab.l.mu.Lock()
			defer lab.l.mu.Unlock()
			return !lab.l.waiting[started.PID]
		})
		got := lab.run(id)
		stoppedAs(t, got)
		if got.Exit == nil || *got.Exit != 143 {
			t.Errorf("the exit the launcher saw is kept: %+v", got)
		}
		if a := Activity(lab.root, lab.state())[id]; a.State != ActivityFailed || !strings.Contains(a.Why, "stopped by the operator") {
			t.Errorf("activity: %+v", a)
		}
		j := lab.entries()
		if last := j[len(j)-1]; !strings.Contains(last.Detail, "stopped "+id+"'s agent builder-"+id+": ended its process group (pid "+strconv.Itoa(started.PID)+")") {
			t.Errorf("journal: %+v", j)
		}
		// it is not started again by itself, and a retry starts a new one
		lab.l.look(ctx, false)
		if r := lab.run(id); r.Started != started.Started {
			t.Errorf("started again after a stop: %+v", r)
		}
		lab.release(id)
		again, err := Restart(ctx, lab.o, Entry{Key: "t", Name: "t", Root: lab.root}, id)
		if err != nil || !again.live() {
			t.Errorf("retry after a stop: %+v, %v", again, err)
		}
	})
	t.Run("an agent that does not end when asked is killed", func(t *testing.T) {
		lab := newAgentLab(t)
		stubborn := filepath.Join(lab.outDir, "stubborn")
		if err := os.WriteFile(stubborn, []byte("#!/bin/sh\ntrap '' TERM\ntouch \""+lab.outDir+"/$FLAI_STORY.txt\"\nwhile :; do sleep 0.05; done\n"), 0o755); err != nil {
			t.Fatal(err)
		}
		lab.cfg.Command = []string{stubborn}
		lab.o.StopGrace = 200 * time.Millisecond
		id := lab.ready("Stubborn")
		lab.l.look(ctx, false)
		lab.trapped(t, id)
		pid := lab.run(id).PID
		// what lab.settle waits for, since a killed process marks nothing
		t.Cleanup(func() { _ = os.WriteFile(filepath.Join(lab.outDir, "ended-"+strconv.Itoa(pid)), nil, 0o644) })
		run, err := stop(lab, id)
		if err != nil {
			t.Fatal(err)
		}
		stoppedAs(t, run)
		if Alive(pid) {
			t.Errorf("pid %d still runs", pid)
		}
		j := lab.entries()
		if last := j[len(j)-1]; !strings.Contains(last.Detail, "killed its process group (pid "+strconv.Itoa(pid)+")") {
			t.Errorf("journal: %+v", j)
		}
	})
	t.Run("a PID that is no longer the agent is left alone", func(t *testing.T) {
		lab := newAgentLab(t)
		lab.hold()
		id := lab.ready("Reused")
		lab.l.look(ctx, false)
		waitFor(t, "it runs", func() bool { return lab.run(id).live() })
		lab.release(id)
		waitFor(t, "it ends", func() bool { return !lab.run(id).live() })
		// after a reboot the PID is another process's: this test's own child,
		// which leads no session of its own
		other := exec.Command("sleep", "30")
		if err := other.Start(); err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() { _ = other.Process.Kill(); _ = other.Wait() })
		lab.l.dir.updateAgent(lab.root, func(s *AgentState) {
			r := *s.Stories[id]
			r.Ended, r.Exit, r.Outcome, r.Why, r.PID = "", nil, "", "", other.Process.Pid
			s.put(&r)
		})
		run, err := stop(lab, id)
		if err != nil {
			t.Fatal(err)
		}
		stoppedAs(t, run)
		if !Alive(other.Process.Pid) {
			t.Error("a process that is not the agent was signalled")
		}
		j := lab.entries()
		if last := j[len(j)-1]; !strings.Contains(last.Detail, "is no longer the agent flai started") {
			t.Errorf("journal: %+v", j)
		}
	})
	t.Run("a session leader that is not the process started is left alone, and settled", func(t *testing.T) {
		if runtime.GOOS != "linux" {
			t.Skip("when a process started is read from /proc")
		}
		other := exec.Command("sleep", "30")
		Detach(other)
		if err := other.Start(); err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() { _ = other.Process.Kill(); _ = other.Wait() })
		pid, start := other.Process.Pid, Started(other.Process.Pid)
		if start == 0 || !Owns(pid, start) || !Owns(pid, 0) {
			t.Fatalf("pid %d is not taken for the process it is (start %d)", pid, start)
		}
		if Owns(pid, start+1) {
			t.Errorf("pid %d is taken for a process that started at another time", pid)
		}
		// after a reboot: the run's PID leads another session, started at another time
		lab := newAgentLab(t)
		lab.hold()
		id := lab.ready("Rebooted")
		lab.l.look(ctx, false)
		waitFor(t, "it runs", func() bool { return lab.run(id).live() })
		if r := lab.run(id); r.Start == 0 {
			t.Errorf("the run does not record when its process started: %+v", r)
		}
		lab.release(id)
		waitFor(t, "it ends", func() bool { return !lab.run(id).live() })
		lab.l.dir.updateAgent(lab.root, func(s *AgentState) {
			r := *s.Stories[id]
			r.Ended, r.Exit, r.Outcome, r.Why, r.PID, r.Start = "", nil, "", "", pid, start+1
			s.put(&r)
		})
		lab.l.look(ctx, false) // what flai serve does after the reboot
		if r := lab.run(id); r.live() || r.Outcome != OutcomeFailed {
			t.Errorf("a run whose PID another process has is not settled: %+v", r)
		}
		if !Alive(pid) {
			t.Error("settling the run signalled the process that has its PID")
		}
	})
	t.Run("an agent waiting for an answer is not started again", func(t *testing.T) {
		lab := newAgentLab(t)
		lab.hold()
		id := lab.ready("Asks")
		lab.l.look(ctx, false)
		waitFor(t, "it runs", func() bool { return lab.run(id).live() })
		lab.move(id, workitem.InProgress)
		agent := lab.run(id).Agent
		th, err := threads.New(lab.repo, threads.NewOptions{Title: "Which port?", On: id, Author: agent, Text: "Eight or nine?", Now: time.Now()})
		if err != nil {
			t.Fatal(err)
		}
		lab.release(id)
		waitFor(t, "it ends asking", func() bool { return lab.run(id).Outcome == OutcomeAsked })
		_ = os.Remove(filepath.Join(lab.outDir, "release-"+id))
		ended := lab.run(id)
		run, err := stop(lab, id)
		if err != nil {
			t.Fatal(err)
		}
		stoppedAs(t, run)
		if run.Ended != ended.Ended || run.Thread != "" {
			t.Errorf("the end it had is kept, and it waits on nothing: %+v", run)
		}
		if _, err := threads.Reply(lab.repo, th.ID, "alex", "Nine.", time.Now()); err != nil {
			t.Fatal(err)
		}
		lab.l.look(ctx, false)
		if r := lab.run(id); r.Started != ended.Started {
			t.Errorf("an answer started a stopped agent again: %+v", r)
		}
		j := lab.entries()
		if last := j[len(j)-1]; !strings.Contains(last.Detail, "waiting for an answer to "+th.ID) {
			t.Errorf("journal: %+v", j)
		}
	})
	t.Run("refusals", func(t *testing.T) {
		lab := newAgentLab(t)
		refusedFor := func(id, want string) {
			t.Helper()
			_, err := stop(lab, id)
			var no *Refused
			if !errors.As(err, &no) || !strings.Contains(no.Why, want) {
				t.Errorf("stop %s: %v, want refused for %q", id, err, want)
			}
		}
		none := lab.backlog("None", nil)
		refusedFor(none, "flai serve has started no agent for "+none)
		done := lab.ready("Done")
		lab.l.look(ctx, false)
		waitFor(t, "it ends", func() bool { r := lab.run(done); return r != nil && r.Ended != "" })
		refusedFor(done, done+"'s agent is not running: it ended")
		refusedFor("T-0001", "is not a story's ID")
		n := len(lab.entries())
		if _, err := stop(lab, done); err == nil || len(lab.entries()) != n {
			t.Errorf("a refused stop was journalled: %+v", lab.entries())
		}
	})
}
