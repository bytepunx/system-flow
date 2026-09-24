//go:build !windows

package serve

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/harness"
	"github.com/bytepunx/system-flow/flai/internal/hostapi"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/threads"
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
	o       Options
	logs    lockedBuffer
}

// lockedBuffer is a log the stub's goroutines and the test share.
type lockedBuffer struct {
	mu  sync.Mutex
	buf bytes.Buffer
}

func (b *lockedBuffer) Write(p []byte) (int, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.buf.Write(p)
}

// said counts the lines logged as agent not started for story.
func (lab *agentLab) said(story string) int {
	lab.logs.mu.Lock()
	defer lab.logs.mu.Unlock()
	n := 0
	for line := range strings.Lines(lab.logs.buf.String()) {
		if strings.Contains(line, `msg="agent not started"`) && strings.Contains(line, "story="+story+" ") {
			n++
		}
	}
	return n
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
		"{ echo \"args: $*\"; echo \"dir: $(pwd)\"; echo \"agent: $FLAI_AGENT\"; echo \"story: $FLAI_STORY\"; echo \"session: $FLAI_SESSION\"; echo \"answered: $FLAI_ANSWERED\"; } > \"$out\"\n" +
		"while [ -f \"" + lab.outDir + "/hold\" ] && [ ! -f \"" + lab.outDir + "/release-$FLAI_STORY\" ]; do sleep 0.05; done\n"
	if err := os.WriteFile(lab.stub, []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	lab.cfg = AgentConfig{Enabled: true, Command: []string{lab.stub, "work on {story}", "--root={root}", "$(echo not a shell)"}, Name: "builder"}
	o := Options{Dir: Dir(filepath.Join(t.TempDir(), "serve")), Logger: slog.New(slog.NewTextHandler(&lab.logs, nil)), Now: time.Now,
		Agent: func(string) AgentConfig { return lab.cfg },
		Host: hostapi.Host{Record: func(e hostapi.Entry) {
			lab.mu.Lock()
			defer lab.mu.Unlock()
			lab.journal = append(lab.journal, e)
		}}}
	_ = os.MkdirAll(string(o.Dir), 0o700)
	lab.o = o
	lab.l = newLauncher(o, Entry{Key: "t", Name: "t", Root: root})
	lab.l.look(context.Background(), false) // what flai serve does when it begins to serve a project
	return lab
}

func (lab *agentLab) ready(title string) string { return lab.readyWith(title, nil) }

// readyWith makes a story with an agent and moves it to ready.
func (lab *agentLab) readyWith(title string, a *manifest.Agent) string {
	lab.t.Helper()
	st, err := lab.repo.Create(workitem.NewOptions{Type: workitem.Story, Title: title, Parent: lab.epic.ID, Owner: "alex", Agent: a, Now: lab.now})
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

func (lab *agentLab) run(id string) *AgentRun { return lab.state().Stories[id] }

// move puts a story in a state, as its agent would.
func (lab *agentLab) move(id, to string) {
	lab.t.Helper()
	if to == workitem.Review {
		// a story goes to review with its tasks written
		if _, err := lab.repo.Create(workitem.NewOptions{Type: workitem.Task, Title: "Work", Parent: id, Owner: "builder-" + id, Now: lab.now}); err != nil {
			lab.t.Fatal(err)
		}
	}
	st, _ := lab.repo.Get(id)
	if _, err := lab.repo.Transition(st, to, "builder-"+id, "", lab.now); err != nil {
		lab.t.Fatal(err)
	}
}

func (lab *agentLab) limit(inProgress int) {
	_ = os.WriteFile(filepath.Join(lab.root, "wip/kanban/board.md"), []byte(fmt.Sprintf("---\ntitle: Board\nstatus: active\nwip_limits:\n  ready: 5\n  in-progress: %d\n  review: 3\n---\n\n# Board\n", inProgress)), 0o644)
}

func (lab *agentLab) release(id string) {
	_ = os.WriteFile(filepath.Join(lab.outDir, "release-"+id), nil, 0o644)
}

func (lab *agentLab) hold() { _ = os.WriteFile(filepath.Join(lab.outDir, "hold"), nil, 0o644) }

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
	waitFor(t, "the stub ended", func() bool { r := lab.run(id); return r != nil && r.Ended != "" })
	got, _ := os.ReadFile(out)
	real, _ := filepath.EvalSymlinks(lab.root)
	for _, want := range []string{
		"args: work on " + id + " --root=" + lab.root + " $(echo not a shell)", // replaced, and nothing else interpreted
		"dir: " + real, "agent: builder-" + id, "story: " + id, "session: 2",
	} {
		if !strings.Contains(string(got), want) {
			t.Errorf("the stub was given\n%s\nwithout %q", got, want)
		}
	}
	last := lab.state().Last
	if last.Story != id || last.Command != "stub-agent" || last.Harness != harness.Command || last.Agent != "builder-"+id || last.PID == 0 || last.Exit == nil || *last.Exit != 0 || last.Log == "" {
		t.Errorf("state: %+v", last)
	}
	// it ended with its story still in ready: that is no work done
	if last.Outcome != OutcomeFailed || !strings.Contains(last.Why, "ended (exit 0) with "+id+" in ready") {
		t.Errorf("outcome: %+v", last)
	}
	j := lab.entries()
	if len(j) != 1 || j[0].Action != hostapi.ActionAgent || j[0].Outcome != "done" || !strings.Contains(j[0].Detail, "started stub-agent (command) for "+id) || j[0].Project != "t" {
		t.Errorf("journal: %+v", j)
	}
	// tried since it entered ready: another look starts nothing, until it
	// enters ready again, and says so once
	lab.l.look(ctx, false)
	lab.l.look(ctx, false)
	if len(lab.entries()) != 1 {
		t.Errorf("started twice: %+v", lab.entries())
	}
	if w := lab.state().Waiting; !strings.Contains(w, id+" has had its agent since it entered ready") || !strings.Contains(w, "failed: ended (exit 0)") {
		t.Errorf("waiting: %q", w)
	}
	if n := lab.said(id); n != 1 {
		t.Errorf("logged %d times why %s waits", n, id)
	}
}

// S-0104: one agent per story, as many as the in-progress limit leaves room
// for, counting those started for stories still in ready.
func TestAnAgentForEveryReadyStoryWithinTheLimit(t *testing.T) {
	lab := newAgentLab(t)
	ctx := context.Background()
	lab.limit(2)
	lab.hold()
	a, b, c := lab.ready("A"), lab.ready("B"), lab.ready("C")
	lab.l.look(ctx, false)
	waitFor(t, "two run", func() bool { return lab.run(a).live() && lab.run(b).live() })
	if lab.run(c) != nil || !strings.Contains(lab.state().Waiting, "leaves no room for "+c) {
		t.Fatalf("the third waits: %+v", lab.state())
	}
	// A's agent pulls its story and finishes it; C still waits while B is in ready with an agent
	lab.move(a, workitem.InProgress)
	lab.l.look(ctx, false)
	if lab.run(c) != nil {
		t.Fatalf("A in progress and B reserved leave no room: %+v", lab.run(c))
	}
	lab.move(a, workitem.Review)
	lab.release(a)
	select {
	case <-lab.l.again:
	case <-time.After(5 * time.Second):
		t.Fatal("the launcher was not told the agent ended")
	}
	if r := lab.run(a); r.Outcome != OutcomeWorked || r.Why != "" {
		t.Errorf("A reached review: %+v", r)
	}
	lab.l.look(ctx, true)
	waitFor(t, "C runs", func() bool { return lab.run(c).live() })
	lab.release(b)
	lab.release(c)
	waitFor(t, "all end", func() bool { return !lab.run(b).live() && !lab.run(c).live() })
}

// S-0104: a story's harness is started through its adapter with the story's
// model and options and the operator's program for it.
func TestAStorysHarnessIsStartedWithItsModel(t *testing.T) {
	lab := newAgentLab(t)
	ctx := context.Background()
	lab.cfg.Command = nil
	lab.cfg.Flai = "/usr/local/bin/flai"
	lab.cfg.Harnesses = map[string]harness.Host{harness.ClaudeCode: {Program: lab.stub, Args: []string{"--permission-mode", "acceptEdits"}}}
	id := lab.readyWith("Worked by claude", &manifest.Agent{Harness: harness.ClaudeCode, Model: "claude-haiku-4-5", Config: map[string]string{"effort": "low"}})
	lab.l.look(ctx, false)
	out := filepath.Join(lab.outDir, id+".txt")
	waitFor(t, "it ended", func() bool { r := lab.run(id); return r != nil && !r.live() })
	got, _ := os.ReadFile(out)
	for _, want := range []string{"-p You are builder-" + id, "--model claude-haiku-4-5", "--effort low", "--permission-mode acceptEdits", `"command":"/usr/local/bin/flai"`, "agent: builder-" + id} {
		if !strings.Contains(string(got), want) {
			t.Errorf("the harness was given\n%s\nwithout %q", got, want)
		}
	}
	if r := lab.run(id); r.Harness != harness.ClaudeCode || r.Model != "claude-haiku-4-5" || r.Command != "stub-agent" {
		t.Errorf("run: %+v", r)
	}

	// what the story says that the harness does not take is refused, and said
	bad := lab.readyWith("Asks for too much", &manifest.Agent{Harness: harness.ClaudeCode, Config: map[string]string{"permission_mode": "bypassPermissions"}})
	lab.l.look(ctx, false)
	if r := lab.run(bad); r == nil || r.Outcome != OutcomeFailed || !strings.Contains(r.Why, `takes no config key "permission_mode"`) || r.PID != 0 {
		t.Fatalf("refused: %+v", r)
	}
	unknown := lab.readyWith("Some other harness", &manifest.Agent{Harness: "cursor"})
	lab.l.look(ctx, false)
	if r := lab.run(unknown); r == nil || !strings.Contains(r.Why, `cannot start harness "cursor"`) {
		t.Fatalf("unknown harness: %+v", r)
	}
	if j := lab.entries(); len(j) != 3 || j[1].Outcome != "failed" || j[2].Outcome != "failed" {
		t.Errorf("journal: %+v", j)
	}
}

func TestAnAgentsOwnSignsAreNotSomeoneAttending(t *testing.T) {
	lab := newAgentLab(t)
	ctx := context.Background()
	lab.hold()
	first := lab.ready("First")
	lab.l.look(ctx, false)
	waitFor(t, "the first runs", func() bool { return lab.run(first).live() })
	// it connects over MCP and writes its narrative, as an agent at work does
	_ = os.MkdirAll(filepath.Join(lab.root, ".flai-cache", "mcp"), 0o755)
	_ = os.WriteFile(filepath.Join(lab.root, ".flai-cache", "mcp", "builder-"+first+".json"), []byte(`{}`), 0o600)
	_ = os.WriteFile(filepath.Join(lab.root, "wip", "agents", first+".md"), []byte("---\nstream: "+first+"\ntitle: First\nupdated: 2026-09-23T00:00:00Z\nagent: builder-"+first+"\n---\n\n# "+first+" First\n"), 0o644)
	second := lab.ready("Second")
	lab.l.look(ctx, false)
	waitFor(t, "the second runs too", func() bool { return lab.run(second).live() })
	lab.release(first)
	lab.release(second)
	waitFor(t, "both end", func() bool { return !lab.run(first).live() && !lab.run(second).live() })
}

func TestAnAgentThatOutlivedFlaiServeIsSettled(t *testing.T) {
	lab := newAgentLab(t)
	id := lab.ready("Left running")
	lab.move(id, workitem.InProgress)
	lab.move(id, workitem.Review)
	// a process that has gone, recorded as running by an earlier flai serve
	gone := exec.Command("true")
	_ = gone.Run()
	lab.l.dir.updateAgent(lab.root, func(s *AgentState) {
		s.put(&AgentRun{Story: id, Agent: "builder-" + id, PID: gone.Process.Pid, Started: time.Now().UTC().Format(time.RFC3339)})
	})
	lab.l.look(context.Background(), false)
	if r := lab.run(id); r.live() || r.Outcome != OutcomeWorked || r.Exit != nil {
		t.Fatalf("settled: %+v", r)
	}
	if lab.state().Running != nil || lab.state().Last.Story != id {
		t.Errorf("running and last: %+v", lab.state())
	}
}

// S-0112: flai host restarts flai serve on every upgrade, host restart, and
// crash. The first look after one starts what has had no agent since it
// entered ready, and nothing that has.
func TestARestartStartsWhatIsReadyAndNothingTwice(t *testing.T) {
	lab := newAgentLab(t)
	ctx := context.Background()
	// tried: its agent ran and ended with the story still in ready
	tried := lab.ready("Tried")
	lab.l.look(ctx, false)
	waitFor(t, "the first agent ended", func() bool { r := lab.run(tried); return r != nil && r.Ended != "" })
	// running: its agent outlives the flai serve that started it
	lab.hold()
	running := lab.ready("Running")
	lab.l.look(ctx, false)
	waitFor(t, "the second agent runs", func() bool { return lab.run(running).live() })
	// missed: it entered ready while no flai serve ran
	missed := lab.ready("Missed")

	fresh := newLauncher(lab.o, lab.l.entry)
	fresh.look(ctx, false)
	fresh.look(ctx, false)
	waitFor(t, "the missed story's agent runs", func() bool { return lab.run(missed).live() })
	started := map[string]int{}
	for _, e := range lab.entries() {
		for _, id := range []string{tried, running, missed} {
			if strings.Contains(e.Detail, "for "+id) {
				started[id]++
			}
		}
	}
	if started[tried] != 1 || started[running] != 1 || started[missed] != 1 {
		t.Errorf("each story started once in all: %v, journal %+v", started, lab.entries())
	}
	if w := lab.state().Waiting; !strings.Contains(w, tried+" has had its agent since it entered ready") || strings.Contains(w, running) || strings.Contains(w, missed) {
		t.Errorf("waiting: %q", w)
	}
	lab.release(running)
	lab.release(missed)
	waitFor(t, "both end", func() bool { return !lab.run(running).live() && !lab.run(missed).live() })
}

func TestNothingIsStartedWhenItShouldNotBe(t *testing.T) {
	ctx := context.Background()
	t.Run("off, or nothing to start it with", func(t *testing.T) {
		lab := newAgentLab(t)
		lab.cfg.Enabled = false
		a := lab.ready("A")
		lab.l.look(ctx, false)
		if st := lab.state(); !strings.Contains(st.Waiting, "the agent host action is off") || lab.said(a) != 1 {
			t.Errorf("off: %q, logged %d", st.Waiting, lab.said(a))
		}
		lab.cfg.Enabled, lab.cfg.Command = true, nil
		b := lab.ready("B")
		lab.l.look(ctx, false)
		lab.l.look(ctx, false)
		if st := lab.state(); st.Running != nil || st.Last != nil || len(lab.entries()) != 0 {
			t.Errorf("started: %+v %+v", st, lab.entries())
		}
		if st := lab.state(); !strings.Contains(st.Waiting, a+", "+b+" name no harness, and no command is set") {
			t.Errorf("waiting: %q", st.Waiting)
		}
		if lab.said(a) != 2 || lab.said(b) != 1 {
			t.Errorf("logged %d for %s and %d for %s: once for each reason", lab.said(a), a, lab.said(b), b)
		}
	})
	t.Run("someone is attending", func(t *testing.T) {
		lab := newAgentLab(t)
		_ = os.MkdirAll(filepath.Join(lab.root, ".flai-cache", "mcp"), 0o755)
		_ = os.WriteFile(filepath.Join(lab.root, ".flai-cache", "mcp", "claude.json"), []byte(`{}`), 0o600)
		a := lab.ready("A")
		lab.l.look(ctx, false)
		if st := lab.state(); st.Running != nil || st.Last != nil || !strings.Contains(st.Waiting, "attending") || !strings.Contains(st.Waiting, ".flai-cache/mcp/claude.json") || !strings.Contains(st.Waiting, a) || lab.said(a) != 1 {
			t.Errorf("attended: %+v", st)
		}
		// the index flai writes itself is no sign of an agent; an old cursor is none either
		old := time.Now().Add(-time.Hour)
		_ = os.Chtimes(filepath.Join(lab.root, ".flai-cache", "mcp", "claude.json"), old, old)
		_ = os.WriteFile(filepath.Join(lab.root, "wip", "agents", "index.md"), []byte("# index\n"), 0o644)
		// a new project's README there is no narrative either (found in S-0104's trial)
		_ = os.WriteFile(filepath.Join(lab.root, "wip", "agents", "README.md"), []byte("# agents\n"), 0o644)
		if yes, sign := attended(lab.root, 6*time.Minute, time.Now(), nil); yes {
			t.Errorf("an hour-old cursor, a fresh index, and a README are nobody attending: %s", sign)
		}
		_ = os.WriteFile(filepath.Join(lab.root, "wip", "agents", "S-0009.md"), []byte("# S-0009\n"), 0o644)
		if yes, _ := attended(lab.root, 6*time.Minute, time.Now(), nil); !yes {
			t.Error("a narrative written now is someone at work")
		}
	})
	t.Run("the limit leaves no room", func(t *testing.T) {
		lab := newAgentLab(t)
		lab.limit(1)
		busy := lab.ready("Busy")
		st, _ := lab.repo.Get(busy)
		if _, err := lab.repo.Transition(st, workitem.InProgress, "alex", "", lab.now); err != nil {
			t.Fatal(err)
		}
		lab.l.look(ctx, false)
		waiting := lab.ready("Waiting")
		lab.l.look(ctx, false)
		if st := lab.state(); st.Running != nil || st.Last != nil || !strings.Contains(st.Waiting, "limit leaves no room for "+waiting) || lab.said(waiting) != 1 {
			t.Errorf("over the limit: %+v", st)
		}
	})
	t.Run("a command that cannot start is reported and journalled", func(t *testing.T) {
		lab := newAgentLab(t)
		lab.cfg.Command = []string{filepath.Join(lab.outDir, "no-such-program"), "{story}"}
		id := lab.ready("A")
		lab.l.look(ctx, false)
		st := lab.state()
		if st.Running != nil || st.Last == nil || st.Last.Story != id || st.Last.Error == "" || st.Last.Outcome != OutcomeFailed {
			t.Fatalf("state: %+v", st)
		}
		j := lab.entries()
		if len(j) != 1 || j[0].Outcome != "failed" || !strings.Contains(j[0].Detail, "could not be started") {
			t.Errorf("journal: %+v", j)
		}
	})
}

// S-0104: what each story's agent is doing, as the board's dots show it.
func TestWhatEachStorysAgentIsDoing(t *testing.T) {
	lab := newAgentLab(t)
	ctx := context.Background()
	lab.limit(3)
	lab.hold()
	asks, blocked, fails := lab.ready("Asks"), lab.ready("Blocked"), lab.ready("Fails")
	lab.l.look(ctx, false)
	waitFor(t, "three run", func() bool { return lab.run(asks).live() && lab.run(blocked).live() && lab.run(fails).live() })
	act := func() map[string]StoryActivity { return Activity(lab.root, lab.state()) }
	for _, id := range []string{asks, blocked, fails} {
		if a := act()[id]; a.State != ActivityWorking || a.Run.Story != id {
			t.Fatalf("%s at work: %+v", id, a)
		}
	}
	// the agent asks the designer: it waits for the answer
	th, err := threads.New(lab.repo, threads.NewOptions{Title: "Which port?", On: asks, Author: "builder-" + asks, Text: "Eight or nine?", Now: time.Now()})
	if err != nil {
		t.Fatal(err)
	}
	if a := act()[asks]; a.State != ActivityWaiting || a.Thread != th.ID || !strings.Contains(a.Why, "Which port?") {
		t.Fatalf("asked: %+v", a)
	}
	if _, err := threads.Reply(lab.repo, th.ID, "alex", "Nine.", time.Now()); err != nil {
		t.Fatal(err)
	}
	if a := act()[asks]; a.State != ActivityWorking {
		t.Fatalf("answered: %+v", a)
	}
	// a blocked story's agent is waiting too
	st, _ := lab.repo.Get(blocked)
	if err := workitem.BlockItem(st, "the API key", time.Now()); err != nil {
		t.Fatal(err)
	}
	if err := lab.repo.Save(st); err != nil {
		t.Fatal(err)
	}
	if a := act()[blocked]; a.State != ActivityWaiting || a.Why != "blocked: the API key" {
		t.Fatalf("blocked: %+v", a)
	}
	// one that ends with its story still in ready has failed; one that reached review has worked
	lab.release(fails)
	lab.move(asks, workitem.InProgress)
	lab.move(asks, workitem.Review)
	lab.release(asks)
	waitFor(t, "both end", func() bool { return !lab.run(fails).live() && !lab.run(asks).live() })
	if a := act()[fails]; a.State != ActivityFailed || !strings.Contains(a.Why, "in ready") {
		t.Errorf("failed: %+v", a)
	}
	if a := act()[asks]; a.State != ActivityWorked {
		t.Errorf("worked: %+v", a)
	}
	lab.release(blocked)
	waitFor(t, "the last ends", func() bool { return !lab.run(blocked).live() })
}

// S-0104: an agent that ends with a question of its own open on its story is
// waiting, not failed, and is started again, as itself, once it is answered.
func TestAnAgentThatEndedAskingIsStartedAgainWhenAnswered(t *testing.T) {
	lab := newAgentLab(t)
	ctx := context.Background()
	lab.hold()
	id := lab.ready("Asks and ends")
	lab.l.look(ctx, false)
	waitFor(t, "it runs", func() bool { return lab.run(id).live() })
	first := lab.run(id)
	lab.move(id, workitem.InProgress)
	th, err := threads.New(lab.repo, threads.NewOptions{Title: "Which port?", On: id, Author: first.Agent, Text: "Eight or nine?", Now: time.Now()})
	if err != nil {
		t.Fatal(err)
	}
	lab.release(id)
	waitFor(t, "it ends", func() bool { return !lab.run(id).live() })
	r := lab.run(id)
	if r.Outcome != OutcomeAsked || r.Thread != th.ID || !strings.Contains(r.Why, "waiting for an answer to "+th.ID) {
		t.Fatalf("asked: %+v", r)
	}
	if a := Activity(lab.root, lab.state())[id]; a.State != ActivityWaiting || a.Thread != th.ID {
		t.Fatalf("activity: %+v", a)
	}
	// nothing is started again while the question is open
	_ = os.Remove(filepath.Join(lab.outDir, "release-"+id))
	lab.l.look(ctx, false)
	if lab.run(id).live() {
		t.Fatal("started again before the answer")
	}
	if _, err := threads.Reply(lab.repo, th.ID, "alex", "Nine.", time.Now()); err != nil {
		t.Fatal(err)
	}
	lab.l.look(ctx, false)
	waitFor(t, "it runs again", func() bool { return lab.run(id).live() })
	again := lab.run(id)
	if again.Agent != first.Agent || again.Session != first.Session || again.Answered != th.ID {
		t.Errorf("the same agent, in its session, for the answer: %+v, first %+v", again, first)
	}
	lab.release(id)
	waitFor(t, "it ends again", func() bool { return !lab.run(id).live() })
	got, _ := os.ReadFile(filepath.Join(lab.outDir, id+".txt"))
	if !strings.Contains(string(got), "answered: "+th.ID) {
		t.Errorf("the command was told what was answered:\n%s", got)
	}
	// this time it ended with no question open and its story in progress: failed
	if r := lab.run(id); r.Outcome != OutcomeFailed {
		t.Errorf("after the answer: %+v", r)
	}
	if j := lab.entries(); len(j) != 2 || !strings.Contains(j[1].Detail, "again for "+id) {
		t.Errorf("journal: %+v", j)
	}
}
