//go:build !windows

package serve

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/hostapi"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// agentLab is a project with an epic, and a launcher whose command is a stub
// that writes what it was given and, when told to, stays until released.
type agentLab struct {
	t       *testing.T
	root    string
	repo    *workitem.Repo
	epic    *workitem.Item
	l       *launcher
	cfg     AgentConfig
	mu      sync.Mutex
	journal []hostapi.Entry
	stub    string
	outDir  string
	now     time.Time
}

func newAgentLab(t *testing.T) *agentLab {
	t.Helper()
	root := t.TempDir()
	_ = os.WriteFile(filepath.Join(root, "system-flow.yaml"), []byte("version: 1\nname: t\nkey: t\nlayout:\n  design: design\n  docs: docs\n  wip: wip\n"), 0o644)
	for _, d := range []string{"wip/kanban/epics", "wip/kanban/stories", "wip/kanban/tasks", "wip/agents"} {
		_ = os.MkdirAll(filepath.Join(root, d), 0o755)
	}
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	lab := &agentLab{t: t, root: root, repo: repo, outDir: t.TempDir(), now: time.Now()}
	lab.epic, err = repo.Create(workitem.NewOptions{Type: workitem.Epic, Title: "Epic", Owner: "alex", Now: lab.now})
	if err != nil {
		t.Fatal(err)
	}
	// The stub is a program, not a shell command line: flai runs it as it
	// stands. It records its arguments, directory, and environment, then
	// waits for a file named after its story to appear, if asked to.
	lab.stub = filepath.Join(lab.outDir, "stub-agent")
	script := "#!/bin/sh\n" +
		"out=\"" + lab.outDir + "/$FLAI_STORY.txt\"\n" +
		"{ echo \"args: $*\"; echo \"dir: $(pwd)\"; echo \"agent: $FLAI_AGENT\"; echo \"story: $FLAI_STORY\"; echo \"session: $FLAI_SESSION\"; } > \"$out\"\n" +
		"while [ -f \"" + lab.outDir + "/hold\" ] && [ ! -f \"" + lab.outDir + "/release-$FLAI_STORY\" ]; do sleep 0.05; done\n"
	if err := os.WriteFile(lab.stub, []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	lab.cfg = AgentConfig{Enabled: true, Command: []string{lab.stub, "work on {story}", "--root={root}", "$(echo not a shell)"}, Name: "builder"}
	o := Options{Dir: Dir(filepath.Join(t.TempDir(), "serve")), Logger: slog.New(slog.DiscardHandler), Now: time.Now,
		Agent: func(string) AgentConfig { return lab.cfg },
		Host: hostapi.Host{Record: func(e hostapi.Entry) {
			lab.mu.Lock()
			defer lab.mu.Unlock()
			lab.journal = append(lab.journal, e)
		}}}
	_ = os.MkdirAll(string(o.Dir), 0o700)
	lab.l = newLauncher(o, Entry{Key: "t", Name: "t", Root: root})
	lab.l.look(context.Background(), false) // what flai serve does when it begins to serve a project
	return lab
}

func (lab *agentLab) ready(title string) string {
	lab.t.Helper()
	st, err := lab.repo.Create(workitem.NewOptions{Type: workitem.Story, Title: title, Parent: lab.epic.ID, Owner: "alex", Now: lab.now})
	if err != nil {
		lab.t.Fatal(err)
	}
	data, _ := os.ReadFile(st.Path)
	_ = os.WriteFile(st.Path, []byte(strings.Replace(string(data), "## Acceptance criteria\n", "## Acceptance criteria\n- [ ] works\n", 1)), 0o644)
	st, _ = lab.repo.Get(st.ID)
	if _, err := lab.repo.Transition(st, workitem.Ready, "alex", "", lab.now); err != nil {
		lab.t.Fatal(err)
	}
	return st.ID
}

func (lab *agentLab) state() AgentState { return lab.l.dir.AgentStates()[lab.root] }

func (lab *agentLab) entries() []hostapi.Entry {
	lab.mu.Lock()
	defer lab.mu.Unlock()
	return append([]hostapi.Entry{}, lab.journal...)
}

func waitFor(t *testing.T, what string, ok func() bool) {
	t.Helper()
	deadline := time.Now().Add(5 * time.Second)
	for !ok() {
		if time.Now().After(deadline) {
			t.Fatalf("never happened: %s", what)
		}
		time.Sleep(20 * time.Millisecond)
	}
}

func TestAStoryEnteringReadyStartsTheOperatorsCommand(t *testing.T) {
	lab := newAgentLab(t)
	ctx := context.Background()
	id := lab.ready("First")
	lab.l.look(ctx, false)
	out := filepath.Join(lab.outDir, id+".txt")
	waitFor(t, "the stub ran", func() bool { _, err := os.Stat(out); return err == nil })
	waitFor(t, "the stub ended", func() bool {
		return lab.state().Running == nil && lab.state().Last != nil && lab.state().Last.Ended != ""
	})
	got, _ := os.ReadFile(out)
	real, _ := filepath.EvalSymlinks(lab.root)
	for _, want := range []string{
		"args: work on " + id + " --root=" + lab.root + " $(echo not a shell)", // replaced, and nothing else interpreted
		"dir: " + real, "agent: builder", "story: " + id, "session: 2",
	} {
		if !strings.Contains(string(got), want) {
			t.Errorf("the stub was given\n%s\nwithout %q", got, want)
		}
	}
	last := lab.state().Last
	if last.Story != id || last.Command != "stub-agent" || last.Agent != "builder" || last.PID == 0 || last.Exit == nil || *last.Exit != 0 || last.Log == "" {
		t.Errorf("state: %+v", last)
	}
	j := lab.entries()
	if len(j) != 1 || j[0].Action != hostapi.ActionAgent || j[0].Outcome != "done" || !strings.Contains(j[0].Detail, "started stub-agent for "+id) || j[0].Project != "t" {
		t.Errorf("journal: %+v", j)
	}
	// nothing new entered ready: another look starts nothing
	lab.l.look(ctx, false)
	if len(lab.entries()) != 1 {
		t.Errorf("started twice: %+v", lab.entries())
	}
}

func TestOneAgentAtATimeAndTheNextWhenItEnds(t *testing.T) {
	lab := newAgentLab(t)
	ctx := context.Background()
	_ = os.WriteFile(filepath.Join(lab.outDir, "hold"), nil, 0o644)
	first := lab.ready("First")
	lab.l.look(ctx, false)
	waitFor(t, "the first runs", func() bool { return lab.state().Running != nil })
	second := lab.ready("Second")
	lab.l.look(ctx, false)
	if st := lab.state(); st.Running == nil || st.Running.Story != first || !strings.Contains(st.Waiting, "still running") {
		t.Fatalf("the second waits for the first: %+v", st)
	}
	if len(lab.entries()) != 1 {
		t.Fatalf("one start so far: %+v", lab.entries())
	}
	// the first session pulls its story, as an agent does, and ends
	st, _ := lab.repo.Get(first)
	if _, err := lab.repo.Transition(st, workitem.InProgress, "builder", "", lab.now); err != nil {
		t.Fatal(err)
	}
	_ = os.WriteFile(filepath.Join(lab.outDir, "release-"+first), nil, 0o644)
	select {
	case <-lab.l.again:
	case <-time.After(5 * time.Second):
		t.Fatal("the launcher was not told the agent ended")
	}
	// the narrative the first agent would have written would count as attending;
	// it wrote none here, so the next is started
	lab.l.look(ctx, true)
	waitFor(t, "the second runs", func() bool { r := lab.state().Running; return r != nil && r.Story == second })
	_ = os.WriteFile(filepath.Join(lab.outDir, "release-"+second), nil, 0o644)
	waitFor(t, "the second ends", func() bool { return lab.state().Running == nil })
}

func TestNothingIsStartedWhenItShouldNotBe(t *testing.T) {
	ctx := context.Background()
	t.Run("what was ready before flai serve began is not news", func(t *testing.T) {
		lab := newAgentLab(t)
		lab.ready("Was ready")
		fresh := newLauncher(Options{Dir: lab.l.dir, Logger: slog.New(slog.DiscardHandler), Now: time.Now, Agent: func(string) AgentConfig { return lab.cfg }}, lab.l.entry)
		fresh.look(ctx, false)
		fresh.look(ctx, false)
		if st := lab.state(); st.Running != nil || st.Last != nil {
			t.Errorf("a restart started a session: %+v", st)
		}
	})
	t.Run("off, or no command", func(t *testing.T) {
		lab := newAgentLab(t)
		lab.cfg.Enabled = false
		lab.ready("A")
		lab.l.look(ctx, false)
		lab.cfg.Enabled, lab.cfg.Command = true, nil
		lab.ready("B")
		lab.l.look(ctx, false)
		if st := lab.state(); st.Running != nil || st.Last != nil || len(lab.entries()) != 0 {
			t.Errorf("started: %+v %+v", st, lab.entries())
		}
	})
	t.Run("someone is attending", func(t *testing.T) {
		lab := newAgentLab(t)
		_ = os.MkdirAll(filepath.Join(lab.root, ".flai-cache", "mcp"), 0o755)
		_ = os.WriteFile(filepath.Join(lab.root, ".flai-cache", "mcp", "claude.json"), []byte(`{}`), 0o600)
		lab.ready("A")
		lab.l.look(ctx, false)
		if st := lab.state(); st.Running != nil || st.Last != nil || !strings.Contains(st.Waiting, "attending") || !strings.Contains(st.Waiting, ".flai-cache/mcp/claude.json") {
			t.Errorf("attended: %+v", st)
		}
		// the index flai writes itself is no sign of an agent; an old cursor is none either
		old := time.Now().Add(-time.Hour)
		_ = os.Chtimes(filepath.Join(lab.root, ".flai-cache", "mcp", "claude.json"), old, old)
		_ = os.WriteFile(filepath.Join(lab.root, "wip", "agents", "index.md"), []byte("# index\n"), 0o644)
		if yes, _ := attended(lab.root, 6*time.Minute, time.Now()); yes {
			t.Error("an hour-old cursor and a fresh index are nobody attending")
		}
	})
	t.Run("the limit leaves no room", func(t *testing.T) {
		lab := newAgentLab(t)
		_ = os.WriteFile(filepath.Join(lab.root, "wip/kanban/board.md"), []byte("---\ntitle: Board\nstatus: active\nwip_limits:\n  ready: 5\n  in-progress: 1\n  review: 3\n---\n\n# Board\n"), 0o644)
		busy := lab.ready("Busy")
		st, _ := lab.repo.Get(busy)
		if _, err := lab.repo.Transition(st, workitem.InProgress, "alex", "", lab.now); err != nil {
			t.Fatal(err)
		}
		lab.l.look(ctx, false)
		lab.ready("Waiting")
		lab.l.look(ctx, false)
		if st := lab.state(); st.Running != nil || st.Last != nil || !strings.Contains(st.Waiting, "limit") {
			t.Errorf("over the limit: %+v", st)
		}
	})
	t.Run("a command that cannot start is reported and journalled", func(t *testing.T) {
		lab := newAgentLab(t)
		lab.cfg.Command = []string{filepath.Join(lab.outDir, "no-such-program"), "{story}"}
		id := lab.ready("A")
		lab.l.look(ctx, false)
		st := lab.state()
		if st.Running != nil || st.Last == nil || st.Last.Story != id || st.Last.Error == "" {
			t.Fatalf("state: %+v", st)
		}
		j := lab.entries()
		if len(j) != 1 || j[0].Outcome != "failed" || !strings.Contains(j[0].Detail, "could not be started") {
			t.Errorf("journal: %+v", j)
		}
	})
}
