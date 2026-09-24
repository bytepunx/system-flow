//go:build !windows

package serve

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

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
		// the limit full: Busy is in progress, and Running's agent holds a place
		lab.limit(2)
		other := lab.ready("Other")
		refusedFor(t, lab, other, "in-progress limit leaves no room for "+other)
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
