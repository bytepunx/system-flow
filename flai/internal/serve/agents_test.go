//go:build !windows

package serve

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/config"
	"github.com/bytepunx/system-flow/flai/internal/harness"
	"github.com/bytepunx/system-flow/flai/internal/host"
	"github.com/bytepunx/system-flow/flai/internal/hostapi"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/messages"
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

// attend leaves the sign of someone attending: an MCP cursor written now.
func (lab *agentLab) attend() {
	lab.t.Helper()
	cursor := filepath.Join(lab.root, ".flai-cache", "mcp", "claude.json")
	_ = os.MkdirAll(filepath.Dir(cursor), 0o755)
	if err := os.WriteFile(cursor, []byte(`{}`), 0o600); err != nil {
		lab.t.Fatal(err)
	}
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
	// Asked to stop (S-0170), it marks its exit as it would have ended.
	script := "#!/bin/sh\n" +
		"trap 'touch \"" + lab.outDir + "/ended-$$\"; exit 143' TERM\n" +
		"out=\"" + lab.outDir + "/$FLAI_STORY.txt\"\n" +
		"{ echo \"args: $*\"; echo \"dir: $(pwd)\"; echo \"agent: $FLAI_AGENT\"; echo \"story: $FLAI_STORY\"; echo \"session: $FLAI_SESSION\"; echo \"answered: $FLAI_ANSWERED\"; echo \"host: ${FLAI_HOST_URL-unset} ${FLAI_HOST_TOKEN-unset}\"; echo \"config: ${FLAI_CONFIG-unset}\"; } > \"$out\"\n" +
		"while [ -f \"" + lab.outDir + "/hold\" ] && [ ! -f \"" + lab.outDir + "/release-$FLAI_STORY\" ]; do sleep 0.05; done\n" +
		"touch \"" + lab.outDir + "/ended-$$\"\n"
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
	// Registered after the temp folders, so it runs before they are removed:
	// a stub still writing, or the launcher recording that one ended, would
	// otherwise put a file in a folder as testing removes it.
	t.Cleanup(lab.settle)
	return lab
}

// startedPID finds the process ID in a journal entry for a started agent.
var startedPID = regexp.MustCompile(`\(pid (\d+)\)`)

// settle ends every stub still held and waits, for a few seconds at most,
// until each stub the journal says was started has marked its exit and the
// launcher has recorded as ended each agent it waits for.
func (lab *agentLab) settle() {
	_ = os.Remove(filepath.Join(lab.outDir, "hold"))
	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		lab.l.mu.Lock()
		waiting := len(lab.l.waiting)
		lab.l.mu.Unlock()
		ended := waiting == 0
		for _, e := range lab.entries() {
			if m := startedPID.FindStringSubmatch(e.Detail); m != nil {
				if _, err := os.Stat(filepath.Join(lab.outDir, "ended-"+m[1])); err != nil {
					ended = false
				}
			}
		}
		if ended {
			return
		}
		time.Sleep(20 * time.Millisecond)
	}
}

func (lab *agentLab) ready(title string) string { return lab.readyWith(title, nil) }

// readyWith makes a story with an agent and moves it to ready.
func (lab *agentLab) readyWith(title string, a *manifest.Agent) string {
	lab.t.Helper()
	id := lab.backlog(title, a)
	lab.toReady(id)
	return id
}

// backlog makes a story with an agent and a criterion, in backlog. It
// touches a path of its own, so that no other story holds it (S-0128).
func (lab *agentLab) backlog(title string, a *manifest.Agent, touches ...string) string {
	lab.t.Helper()
	if touches == nil {
		touches = []string{"docs/" + strings.ToLower(strings.ReplaceAll(title, " ", "-"))}
	}
	st, err := lab.repo.Create(workitem.NewOptions{Type: workitem.Story, Title: title, Parent: lab.epic.ID, Owner: "alex", Agent: a, Touches: touches, Now: lab.now})
	if err != nil {
		lab.t.Fatal(err)
	}
	data, _ := os.ReadFile(st.Path)
	_ = os.WriteFile(st.Path, []byte(strings.Replace(string(data), "## Acceptance criteria\n", "## Acceptance criteria\n- [ ] works\n", 1)), 0o644)
	return st.ID
}

// toReady moves a story to ready, as the designer does from the board.
func (lab *agentLab) toReady(id string) {
	lab.t.Helper()
	st, err := lab.repo.Get(id)
	if err != nil {
		lab.t.Fatal(err)
	}
	if _, err := lab.repo.Transition(st, workitem.Ready, "alex", "", lab.now); err != nil {
		lab.t.Fatal(err)
	}
}

// S-0116: a story moved from backlog to ready gets its agent at the next look
// while the in-progress limit has room, and once there is room when it has none.
func TestAStoryMovedFromBacklogToReadyGetsItsAgentWhenThereIsRoom(t *testing.T) {
	lab := newAgentLab(t)
	ctx := context.Background()
	lab.limit(1)
	lab.hold()
	first, second := lab.backlog("First", nil), lab.backlog("Second", nil)
	lab.l.look(ctx, false)
	if len(lab.entries()) != 0 {
		t.Fatalf("a story in backlog was started: %+v", lab.entries())
	}
	lab.toReady(first)
	lab.l.look(ctx, false)
	waitFor(t, "the first story's agent runs", func() bool { return lab.run(first).live() })
	// its agent pulls it; the limit is full when the second is moved to ready
	lab.move(first, workitem.InProgress)
	lab.toReady(second)
	lab.l.look(ctx, false)
	if lab.run(second) != nil || !strings.Contains(lab.state().Waiting, "limit leaves no room for "+second) {
		t.Fatalf("started over the limit: %+v", lab.state())
	}
	// the first reaches review: there is room, and the second is started
	lab.move(first, workitem.Review)
	lab.l.look(ctx, false)
	waitFor(t, "the second story's agent runs", func() bool { return lab.run(second).live() })
	lab.release(first)
	lab.release(second)
	waitFor(t, "both end", func() bool { return !lab.run(first).live() && !lab.run(second).live() })
}

// S-0243, I-0007: while review is full no ready story's agent is started,
// whatever room the in-progress limit leaves, and the state says why; an
// agent already started goes on; once a story leaves review, the waiting
// one is started.
func TestNoAgentStartsWhileReviewIsFull(t *testing.T) {
	lab := newAgentLab(t)
	ctx := context.Background()
	lab.limits(3, 1)
	lab.hold()
	started := lab.ready("Started")
	lab.l.look(ctx, false)
	waitFor(t, "the started story's agent runs", func() bool { return lab.run(started).live() })
	reviewed := lab.ready("Reviewed")
	lab.move(reviewed, workitem.InProgress)
	lab.move(reviewed, workitem.Review)
	waiting := lab.ready("Waiting")
	lab.l.look(ctx, false)
	full := "review is full (1 of 1): accept or send back a story; it holds " + waiting
	if lab.run(waiting) != nil || !strings.Contains(lab.state().Waiting, full) || lab.said(waiting) != 1 {
		t.Fatalf("started with review full: %+v", lab.state())
	}
	if !lab.run(started).live() {
		t.Errorf("a full review stopped an agent already started: %+v", lab.run(started))
	}
	// the reviewed story is sent back: review has room, and the waiting one is started
	back, _ := lab.repo.Get(reviewed)
	if _, err := lab.repo.Transition(back, workitem.InProgress, "alex", "not yet", lab.now); err != nil {
		t.Fatal(err)
	}
	lab.l.look(ctx, false)
	waitFor(t, "the waiting story's agent runs", func() bool { return lab.run(waiting).live() })
	lab.release(started)
	lab.release(waiting)
	waitFor(t, "both end", func() bool { return !lab.run(started).live() && !lab.run(waiting).live() })
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

func (lab *agentLab) limit(inProgress int) { lab.limits(inProgress, 3) }

// limits sets the in-progress and review limits (S-0243).
func (lab *agentLab) limits(inProgress, review int) {
	_ = os.WriteFile(filepath.Join(lab.root, "wip/kanban/board.md"), []byte(fmt.Sprintf("---\ntitle: Board\nstatus: active\nwip_limits:\n  ready: 5\n  in-progress: %d\n  review: %d\n---\n\n# Board\n", inProgress, review)), 0o644)
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

// S-0183: an agent gets no way to the operator's host and no config of flai
// serve's (I-0044), and keeps the rest of its environment.
func TestAnAgentIsStartedWithoutTheHostsAddressTokenOrConfig(t *testing.T) {
	t.Setenv(host.URLEnv, "http://127.0.0.1:4241")
	t.Setenv(host.TokenEnv, "host-secret")
	t.Setenv(config.EnvVar, filepath.Join(t.TempDir(), "config.json"))
	lab := newAgentLab(t)
	id := lab.ready("Clean")
	lab.l.look(context.Background(), false)
	out := filepath.Join(lab.outDir, id+".txt")
	waitFor(t, "the stub ended", func() bool { r := lab.run(id); return r != nil && r.Ended != "" })
	got, _ := os.ReadFile(out)
	for _, want := range []string{"host: unset unset\n", "config: unset\n", "story: " + id + "\n"} {
		if !strings.Contains(string(got), want) {
			t.Errorf("the stub was given\n%s\nwithout %q", got, want)
		}
	}
}

// setAgent changes a story's agent, as the dashboard or flai edit would.
func (lab *agentLab) setAgent(id string, a *manifest.Agent) {
	lab.t.Helper()
	it, err := lab.repo.Get(id)
	if err != nil {
		lab.t.Fatal(err)
	}
	it.Agent = a
	if err := lab.repo.Save(it); err != nil {
		lab.t.Fatal(err)
	}
}

// S-0116: a story in ready whose agent failed gets another once its agent is
// changed, within the in-progress limit, and not while its agent runs.
func TestAChangedAgentStartsAReadyStoryAgain(t *testing.T) {
	lab := newAgentLab(t)
	ctx := context.Background()
	lab.limit(2)
	started := func() int { return len(lab.entries()) }
	ended := func(id string) func() bool {
		return func() bool { r := lab.run(id); return r != nil && r.Ended != "" }
	}
	// a harness flai cannot start: the run fails, and records what it was started with
	id := lab.readyWith("A", &manifest.Agent{Harness: "no-such-harness"})
	lab.l.look(ctx, false)
	if r := lab.run(id); r == nil || r.Outcome != OutcomeFailed || r.StoryAgent == nil || r.StoryAgent.Harness != "no-such-harness" {
		t.Fatalf("the failed start: %+v", r)
	}
	lab.l.look(ctx, false)
	if started() != 1 || !strings.Contains(lab.state().Waiting, "it gets another when it enters ready again, its agent is changed, or it is restarted") {
		t.Fatalf("started again with nothing changed: %d, %q", started(), lab.state().Waiting)
	}
	// the designer changes its agent while it stays in ready
	lab.setAgent(id, &manifest.Agent{Model: "m1"})
	lab.l.look(ctx, false)
	if started() != 2 {
		t.Fatalf("a changed harness started nothing: %+v", lab.entries())
	}
	waitFor(t, "the second run ended", ended(id))
	if r := lab.run(id); r.StoryAgent == nil || r.StoryAgent.Model != "m1" || r.Harness != harness.Command {
		t.Errorf("the second run: %+v", r)
	}
	lab.l.look(ctx, false)
	if started() != 2 {
		t.Fatalf("started again with nothing changed since: %+v", lab.entries())
	}
	// a config value is a change; the in-progress limit still holds it back
	busy := lab.ready("Busy")
	lab.move(busy, workitem.InProgress)
	other := lab.ready("Other")
	lab.move(other, workitem.InProgress)
	lab.setAgent(id, &manifest.Agent{Model: "m1", Config: map[string]string{"effort": "high"}})
	lab.l.look(ctx, false)
	if started() != 2 || !strings.Contains(lab.state().Waiting, "limit leaves no room for "+id) {
		t.Fatalf("started over the limit: %d, %q", started(), lab.state().Waiting)
	}
	lab.move(busy, workitem.Review)
	lab.hold()
	lab.l.look(ctx, false)
	waitFor(t, "the third run", func() bool { return lab.run(id).live() })
	if started() != 3 || lab.run(id).StoryAgent.Config["effort"] != "high" {
		t.Fatalf("room again: %+v", lab.run(id))
	}
	// a change while its agent runs starts nothing more
	lab.setAgent(id, &manifest.Agent{Model: "m2"})
	lab.l.look(ctx, false)
	if started() != 3 {
		t.Errorf("a change while its agent ran started another: %+v", lab.entries())
	}
	lab.release(id)
	waitFor(t, "the third run ended", ended(id))
}

// A run recorded before S-0116 says nothing of the story's agent: it counts as
// unchanged, so upgrading flai starts no story again on its own.
func TestARunWithoutItsStorysAgentCountsAsUnchanged(t *testing.T) {
	for _, c := range []struct {
		run  *manifest.Agent
		now  *manifest.Agent
		want bool
	}{
		{nil, &manifest.Agent{Model: "m"}, false},
		{&manifest.Agent{}, nil, false},
		{&manifest.Agent{}, &manifest.Agent{Harness: "claude-code"}, true},
		{&manifest.Agent{Model: "m", Config: map[string]string{"effort": "high"}}, &manifest.Agent{Model: "m", Config: map[string]string{"effort": "high"}}, false},
		{&manifest.Agent{Model: "m", Config: map[string]string{"effort": "high"}}, &manifest.Agent{Model: "m", Config: map[string]string{"effort": "low"}}, true},
	} {
		if got := agentChanged(&AgentRun{StoryAgent: c.run}, c.now); got != c.want {
			t.Errorf("run %+v, story %+v: changed %v, want %v", c.run, c.now, got, c.want)
		}
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

// S-0116, ADR-0043: someone attending holds nothing back. A fresh MCP cursor
// and a narrative written now do not keep a ready story from its agent while
// the in-progress limit has room.
func TestSomeoneAttendingHoldsNothingBack(t *testing.T) {
	lab := newAgentLab(t)
	ctx := context.Background()
	lab.hold()
	lab.attend()
	a := lab.ready("A")
	// an agent at work on another story writes its narrative
	_ = os.WriteFile(filepath.Join(lab.root, "wip", "agents", "S-0009.md"), []byte("---\nstream: S-0009\ntitle: Other\nupdated: 2026-09-24T00:00:00Z\nagent: someone\n---\n\n# S-0009 Other\n"), 0o644)
	lab.l.look(ctx, false)
	waitFor(t, "A runs", func() bool { return lab.run(a).live() })
	if w := lab.state().Waiting; w != "" {
		t.Errorf("nothing waits: %q", w)
	}
	lab.release(a)
	waitFor(t, "A ends", func() bool { return !lab.run(a).live() })
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

// S-0182: an agent that ended asking is started again on the answer only
// while its story is open. Sent back to ready meanwhile, the story waits for
// the answer and then for its hold, as any ready story does, and its agent
// is started again in its session once nothing holds it.
func TestAnAnsweredAgentWhoseStoryIsBackInReadyWaitsForItsHold(t *testing.T) {
	lab := newAgentLab(t)
	ctx := context.Background()
	lab.limit(3)
	lab.hold()
	id := lab.readyTouching("Asks and goes back", "flai/cmd")
	lab.l.look(ctx, false)
	waitFor(t, "it runs", func() bool { return lab.run(id).live() })
	first := lab.run(id)
	lab.move(id, workitem.InProgress)
	th, err := threads.New(lab.repo, threads.NewOptions{Title: "Which port?", On: id, Author: first.Agent, Text: "Eight or nine?", Now: time.Now()})
	if err != nil {
		t.Fatal(err)
	}
	lab.release(id)
	waitFor(t, "it ends asking", func() bool { r := lab.run(id); return !r.live() && r.Outcome == OutcomeAsked })
	_ = os.Remove(filepath.Join(lab.outDir, "release-"+id))

	// the designer sends it back to ready, and another story claims a file in its folder
	st, _ := lab.repo.Get(id)
	if _, err := lab.repo.Transition(st, workitem.Ready, "alex", "not yet", lab.now); err != nil {
		t.Fatal(err)
	}
	other := lab.backlog("Other", nil, "flai/cmd/prime.go")
	lab.toReady(other)
	lab.move(other, workitem.InProgress)
	lab.l.look(ctx, false)
	if r := lab.run(id); r.live() || !strings.Contains(lab.state().Waiting, id+"'s agent is waiting for an answer to "+th.ID) {
		t.Fatalf("before the answer: %+v, waiting %q", r, lab.state().Waiting)
	}

	// answered while held: nothing starts, and the board says what holds it
	if _, err := threads.Reply(lab.repo, th.ID, "alex", "Nine.", time.Now()); err != nil {
		t.Fatal(err)
	}
	lab.l.look(ctx, false)
	why := "held (overlap): touches flai/cmd, which holds flai/cmd/prime.go that " + other + " (in progress) touches"
	if r := lab.run(id); r.live() || !strings.Contains(lab.state().Waiting, id+" "+why) {
		t.Fatalf("started past its hold: %+v, waiting %q", r, lab.state().Waiting)
	}
	if a := Activity(lab.root, lab.state())[id]; a.Hold == nil || !strings.HasPrefix(a.Why, why) {
		t.Errorf("activity: %+v", a)
	}

	// the other story is cancelled: started again, in its session, for the answer
	o, _ := lab.repo.Get(other)
	if _, err := lab.repo.Transition(o, workitem.Cancelled, "alex", "not now", lab.now); err != nil {
		t.Fatal(err)
	}
	lab.l.look(ctx, false)
	waitFor(t, "it runs again", func() bool { return lab.run(id).live() })
	if again := lab.run(id); again.Agent != first.Agent || again.Session != first.Session || again.Answered != th.ID {
		t.Errorf("the same agent, in its session, for the answer: %+v, first %+v", again, first)
	}
	lab.release(id)
	waitFor(t, "it ends again", func() bool { return !lab.run(id).live() })
}

// S-0182: an answered agent whose story went back to backlog is not started.
func TestAnAnsweredAgentWhoseStoryIsInBacklogIsNotStarted(t *testing.T) {
	lab := newAgentLab(t)
	ctx := context.Background()
	lab.hold()
	id := lab.ready("Asks and is shelved")
	lab.l.look(ctx, false)
	waitFor(t, "it runs", func() bool { return lab.run(id).live() })
	first := lab.run(id)
	th, err := threads.New(lab.repo, threads.NewOptions{Title: "Which port?", On: id, Author: first.Agent, Text: "Eight or nine?", Now: time.Now()})
	if err != nil {
		t.Fatal(err)
	}
	lab.release(id)
	waitFor(t, "it ends asking", func() bool { r := lab.run(id); return !r.live() && r.Outcome == OutcomeAsked })
	st, _ := lab.repo.Get(id)
	if _, err := lab.repo.Transition(st, workitem.Backlog, "alex", "later", lab.now); err != nil {
		t.Fatal(err)
	}
	if _, err := threads.Reply(lab.repo, th.ID, "alex", "Nine.", time.Now()); err != nil {
		t.Fatal(err)
	}
	lab.l.look(ctx, false)
	if r := lab.run(id); r.live() || r.Answered != "" {
		t.Fatalf("started for a story in backlog: %+v", r)
	}
}

// inProgress makes a story in progress that no agent here works, for a
// story's agent to converse with.
func (lab *agentLab) inProgress(title string) string {
	lab.t.Helper()
	id := lab.backlog(title, nil)
	lab.toReady(id)
	lab.move(id, workitem.InProgress)
	return id
}

// endedAgo makes story's newest run, which has ended, read as one that
// started and ended d earlier than it did.
func (lab *agentLab) endedAgo(story string, d time.Duration) {
	lab.t.Helper()
	lab.l.dir.updateAgent(lab.root, func(s *AgentState) {
		r := *s.Stories[story]
		for _, at := range []*string{&r.Started, &r.Ended} {
			t, err := time.Parse(time.RFC3339, *at)
			if err != nil {
				lab.t.Fatal(err)
			}
			*at = t.Add(-d).Format(time.RFC3339)
		}
		s.put(&r)
	})
}

// S-0335: an agent that ends waiting on another story's agent's reply ended
// asking, and says it waits on that story rather than on the operator, while
// it runs and once it has ended. Its story sent back to ready waits for the
// reply, and once the other story's agent replies its agent is started again,
// as itself, in its session, for the conversation.
func TestAnAgentThatEndedOnAConversationIsStartedAgainOnTheReply(t *testing.T) {
	lab := newAgentLab(t)
	ctx := context.Background()
	lab.hold()
	other := lab.inProgress("Other")
	id := lab.ready("Messages and ends")
	lab.l.look(ctx, false)
	waitFor(t, "it runs", func() bool { return lab.run(id).live() })
	first := lab.run(id)
	lab.move(id, workitem.InProgress)
	c, err := messages.Send(lab.repo, messages.SendOptions{From: id, To: other, Author: first.Agent, Text: "Which port do you bind?", Now: time.Now()})
	if err != nil {
		t.Fatal(err)
	}
	waits := other + "'s agent to reply on " + c.ID + ": " + c.Title
	if a := Activity(lab.root, lab.state())[id]; a.State != ActivityWaiting || a.Thread != "" || !slices.Equal(a.WaitsOn, []string{other}) ||
		!slices.Equal(a.Conversations, []string{c.ID}) || a.Why != "waiting for "+waits {
		t.Fatalf("running, it waits on the other story's agent: %+v", a)
	}
	lab.release(id)
	waitFor(t, "it ends", func() bool { return !lab.run(id).live() })
	r := lab.run(id)
	if r.Outcome != OutcomeAsked || r.Thread != "" || !slices.Equal(r.Conversations, []string{c.ID}) || !slices.Equal(r.WaitsOn, []string{other}) || r.Why != "waiting for "+waits {
		t.Fatalf("asked: %+v", r)
	}
	if a := Activity(lab.root, lab.state())[id]; a.State != ActivityWaiting || a.Thread != "" || !slices.Equal(a.WaitsOn, []string{other}) ||
		!slices.Equal(a.Conversations, []string{c.ID}) || a.Why != r.Why {
		t.Fatalf("ended, it waits on the other story's agent: %+v", a)
	}

	// sent back to ready, it waits for the reply, and the board says so
	_ = os.Remove(filepath.Join(lab.outDir, "release-"+id))
	st, _ := lab.repo.Get(id)
	if _, err := lab.repo.Transition(st, workitem.Ready, "alex", "not yet", lab.now); err != nil {
		t.Fatal(err)
	}
	lab.l.look(ctx, false)
	if r := lab.run(id); r.live() || !strings.Contains(lab.state().Waiting, id+"'s agent is waiting for "+other+"'s agent to reply on "+c.ID) {
		t.Fatalf("before the reply: %+v, waiting %q", r, lab.state().Waiting)
	}

	if _, err := messages.Reply(lab.repo, c.ID, other, "builder-"+other, "Nine.", time.Now()); err != nil {
		t.Fatal(err)
	}
	lab.l.look(ctx, false)
	waitFor(t, "it runs again", func() bool { return lab.run(id).live() })
	if again := lab.run(id); again.Agent != first.Agent || again.Session != first.Session || again.Answered != c.ID {
		t.Errorf("the same agent, in its session, for the reply: %+v, first %+v", again, first)
	}
	lab.release(id)
	waitFor(t, "it ends again", func() bool { return !lab.run(id).live() })
	got, _ := os.ReadFile(filepath.Join(lab.outDir, id+".txt"))
	if !strings.Contains(string(got), "answered: "+c.ID) {
		t.Errorf("the command was told what was answered:\n%s", got)
	}
	if j := lab.entries(); len(j) != 2 || !strings.Contains(j[1].Detail, "again for "+id+" as "+first.Agent+", a message on "+c.ID) {
		t.Errorf("journal: %+v", j)
	}
}

// S-0335: an agent that ended asking the operator and waiting on another
// story's agent records both. A new message to its story starts it again, in
// its session, while the question to the operator is still open, and the
// message it was started for does not start it again when it ends once more.
func TestAnAgentThatEndedAskingIsStartedAgainOnANewMessageToItsStory(t *testing.T) {
	lab := newAgentLab(t)
	ctx := context.Background()
	lab.hold()
	other := lab.inProgress("Other")
	id := lab.ready("Asks both and ends")
	lab.l.look(ctx, false)
	waitFor(t, "it runs", func() bool { return lab.run(id).live() })
	first := lab.run(id)
	lab.move(id, workitem.InProgress)
	th, err := threads.New(lab.repo, threads.NewOptions{Title: "Which port?", On: id, Author: first.Agent, Text: "Eight or nine?", Now: time.Now()})
	if err != nil {
		t.Fatal(err)
	}
	c, err := messages.Send(lab.repo, messages.SendOptions{From: id, To: other, Author: first.Agent, Text: "Do you change flai/cmd?", Now: time.Now()})
	if err != nil {
		t.Fatal(err)
	}
	both := "waiting for an answer to " + th.ID + ": " + th.Title + "; and for " + other + "'s agent to reply on " + c.ID + ": " + c.Title
	if a := Activity(lab.root, lab.state())[id]; a.State != ActivityWaiting || a.Thread != th.ID || !slices.Equal(a.WaitsOn, []string{other}) || a.Why != both {
		t.Fatalf("running, the operator's question is its thread, and the other story is named too: %+v", a)
	}
	lab.release(id)
	waitFor(t, "it ends", func() bool { return !lab.run(id).live() })
	if r := lab.run(id); r.Outcome != OutcomeAsked || r.Thread != th.ID || !slices.Equal(r.Conversations, []string{c.ID}) || !slices.Equal(r.WaitsOn, []string{other}) || r.Why != both {
		t.Fatalf("asked on both: %+v", r)
	}
	if a := Activity(lab.root, lab.state())[id]; a.State != ActivityWaiting || a.Thread != th.ID || !slices.Equal(a.WaitsOn, []string{other}) {
		t.Fatalf("ended: %+v", a)
	}
	_ = os.Remove(filepath.Join(lab.outDir, "release-"+id))
	lab.l.look(ctx, false)
	if lab.run(id).live() {
		t.Fatal("started again before an answer, a reply, or a message")
	}

	// a minute on, the other story's agent writes to it about something else
	lab.endedAgo(id, time.Minute)
	c2, err := messages.Send(lab.repo, messages.SendOptions{From: other, To: id, Author: "builder-" + other, Text: "I am renaming flai/cmd/prime.go.", Now: time.Now()})
	if err != nil {
		t.Fatal(err)
	}
	lab.l.look(ctx, false)
	waitFor(t, "it runs again", func() bool { return lab.run(id).live() })
	if again := lab.run(id); again.Agent != first.Agent || again.Session != first.Session || again.Answered != c2.ID {
		t.Errorf("the same agent, in its session, for the message: %+v, first %+v", again, first)
	}
	if open, _ := threads.Get(lab.repo, th.ID); !awaitsAnswer(open, first.Agent) {
		t.Errorf("the question to the operator was answered: %+v", open)
	}

	// it ends with its question still open, and the message it read starts nothing
	lab.release(id)
	waitFor(t, "it ends again", func() bool { return !lab.run(id).live() })
	if r := lab.run(id); r.Outcome != OutcomeAsked || r.Thread != th.ID || !slices.Equal(r.Conversations, []string{c.ID}) {
		t.Fatalf("asked again: %+v", r)
	}
	_ = os.Remove(filepath.Join(lab.outDir, "release-"+id))
	lab.l.look(ctx, false)
	if r := lab.run(id); r.live() || len(lab.entries()) != 2 {
		t.Fatalf("started again for the message it was started for: %+v, journal %+v", r, lab.entries())
	}
}

// S-0335, I-0095: a reply that comes before the end of a run is judged, after
// the run started, leaves its agent with a message unanswered: it ended
// asking on that conversation, which answers it at once, and the next look
// starts it again in its session. A message that came before the run started
// is no such reply.
func TestAReplyBeforeTheEndIsJudgedStartsTheAgentAgain(t *testing.T) {
	lab := newAgentLab(t)
	ctx := context.Background()
	other, id := lab.inProgress("Other"), lab.inProgress("Messages")
	agent := "builder-" + id
	started := time.Now().UTC().Add(-time.Minute)
	c, err := messages.Send(lab.repo, messages.SendOptions{From: id, To: other, Author: agent, Text: "Which port do you bind?", Now: started.Add(10 * time.Second)})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := messages.Reply(lab.repo, c.ID, other, "builder-"+other, "Nine.", time.Now()); err != nil {
		t.Fatal(err)
	}
	later := AgentRun{Story: id, Agent: agent, Started: time.Now().UTC().Add(time.Minute).Format(time.RFC3339)}
	if judgeRun(lab.root, &later, new(int)); later.Outcome != OutcomeFailed || len(later.Conversations) != 0 {
		t.Errorf("a run started after the reply: %+v", later)
	}

	run := AgentRun{Story: id, Agent: agent, Started: started.Format(time.RFC3339), Ended: time.Now().UTC().Format(time.RFC3339), Session: "earlier"}
	judgeRun(lab.root, &run, new(int))
	if run.Outcome != OutcomeAsked || run.Thread != "" || !slices.Equal(run.Conversations, []string{c.ID}) || len(run.WaitsOn) != 0 ||
		run.Why != "ended with a message from "+other+"'s agent on "+c.ID+" unanswered: "+c.Title {
		t.Fatalf("judged: %+v", run)
	}
	if got := answered(lab.repo, &run); got != c.ID {
		t.Errorf("answered at once by %q", got)
	}
	lab.l.dir.updateAgent(lab.root, func(s *AgentState) { s.put(&run) })
	if a := Activity(lab.root, lab.state())[id]; a.State != ActivityWaiting || len(a.WaitsOn) != 0 || !slices.Equal(a.Conversations, []string{c.ID}) {
		t.Errorf("activity: %+v", a)
	}
	lab.l.look(ctx, false)
	waitFor(t, "it runs again", func() bool { r := lab.run(id); return r.Answered == c.ID })
	if again := lab.run(id); again.Agent != agent || again.Session != "earlier" {
		t.Errorf("the same agent, in its session: %+v", again)
	}
	waitFor(t, "it ends", func() bool { return !lab.run(id).live() })
}

// S-0317, I-0095: an answer to the question an agent ended waiting on that
// comes before the end of its run is judged leaves no question open: the run
// ended asking on that thread, which answers it at once, and the next look
// starts it again in its session. A question asked before the run started
// is no such question, and the run started again for its answer, writing on
// it no more, ends failed rather than asking on it once more.
func TestAnAnswerBeforeTheEndIsJudgedStartsTheAgentAgain(t *testing.T) {
	lab := newAgentLab(t)
	ctx := context.Background()
	id := lab.inProgress("Asks")
	agent := "builder-" + id
	started := time.Now().UTC().Add(-time.Minute)
	th, err := threads.New(lab.repo, threads.NewOptions{Title: "Which port?", On: id, Author: agent, Text: "Eight or nine?", Now: started.Add(10 * time.Second)})
	if err != nil {
		t.Fatal(err)
	}
	ended := time.Now().UTC()
	if _, err := threads.Reply(lab.repo, th.ID, "alex", "Nine.", ended); err != nil {
		t.Fatal(err)
	}
	later := AgentRun{Story: id, Agent: agent, Started: ended.Add(time.Minute).Format(time.RFC3339)}
	if judgeRun(lab.root, &later, new(int)); later.Outcome != OutcomeFailed || later.Thread != "" {
		t.Errorf("a run started after the question: %+v", later)
	}

	run := AgentRun{Story: id, Agent: agent, Started: started.Format(time.RFC3339), Ended: ended.Format(time.RFC3339), Session: "earlier"}
	judgeRun(lab.root, &run, new(int))
	if run.Outcome != OutcomeAsked || run.Thread != th.ID || len(run.Conversations) != 0 ||
		run.Why != "ended asking on "+th.ID+", answered before the end was judged: "+th.Title {
		t.Fatalf("judged: %+v", run)
	}
	if got := answered(lab.repo, &run); got != th.ID {
		t.Errorf("answered at once by %q", got)
	}
	lab.l.dir.updateAgent(lab.root, func(s *AgentState) { s.put(&run) })
	if a := Activity(lab.root, lab.state())[id]; a.State != ActivityWaiting || a.Thread != th.ID {
		t.Errorf("activity: %+v", a)
	}
	lab.l.look(ctx, false)
	waitFor(t, "it runs again", func() bool { r := lab.run(id); return r.Answered == th.ID })
	if again := lab.run(id); again.Agent != agent || again.Session != "earlier" {
		t.Errorf("the same agent, in its session: %+v", again)
	}
	waitFor(t, "it ends", func() bool { return !lab.run(id).live() })
	if r := lab.run(id); r.Outcome != OutcomeFailed || r.Thread != "" {
		t.Errorf("asked again on the answer it was started for: %+v", r)
	}
}

// S-0317, I-0095: which thread on its story answers the question a run's
// agent wrote after the second the run started.
func TestAskedAnswered(t *testing.T) {
	lab := newAgentLab(t)
	id, other := lab.inProgress("Asks"), lab.inProgress("Other")
	agent := "builder-" + id
	started := time.Now().UTC().Add(-time.Hour).Truncate(time.Second)
	run := AgentRun{Story: id, Agent: agent, Started: started.Format(time.RFC3339)}
	ask := func(on string, at time.Time) *threads.Thread {
		t.Helper()
		th, err := threads.New(lab.repo, threads.NewOptions{Title: "Which port?", On: on, Author: agent, Text: "Eight or nine?", Now: at})
		if err != nil {
			t.Fatal(err)
		}
		return th
	}
	check := func(what string, want *threads.Thread) {
		t.Helper()
		got := askedAnswered(lab.repo, &run)
		if want == nil && got != nil || want != nil && (got == nil || got.ID != want.ID) {
			t.Errorf("%s: answered by %+v, want %+v", what, got, want)
		}
	}
	during, then := started.Add(time.Second), started.Add(time.Minute)

	// in the second the run started, the last run's question, maybe
	before := ask(id, started)
	if _, err := threads.Reply(lab.repo, before.ID, "alex", "Nine.", then); err != nil {
		t.Fatal(err)
	}
	elsewhere := ask(other, during)
	if _, err := threads.Reply(lab.repo, elsewhere.ID, "alex", "Nine.", then); err != nil {
		t.Fatal(err)
	}
	unanswered := ask(id, during)
	recommended := ask(id, during)
	if _, err := threads.ReplyWith(lab.repo, recommended.ID, "planner", "Nine.", then, threads.Marks{Recommendation: true}); err != nil {
		t.Fatal(err)
	}
	selfResolved := ask(id, during)
	if _, err := threads.Resolve(lab.repo, selfResolved.ID, agent, "found it", then); err != nil {
		t.Fatal(err)
	}
	check("asked as the run started, on another story, unanswered, recommended, or resolved by itself", nil)

	if _, err := threads.Confirm(lab.repo, recommended.ID, "alex", then); err != nil {
		t.Fatal(err)
	}
	check("a recommendation confirmed", recommended)
	if _, err := threads.Resolve(lab.repo, recommended.ID, agent, "done", then); err != nil {
		t.Fatal(err)
	}

	if _, err := threads.Resolve(lab.repo, unanswered.ID, "alex", "use nine", then); err != nil {
		t.Fatal(err)
	}
	check("resolved by the operator", unanswered)
}

// S-0335: a message written in the second a run ended, after it was judged,
// is new to it; one its story answered in that second, or one from before,
// is not.
func TestAMessageInTheSecondARunEndedIsNew(t *testing.T) {
	lab := newAgentLab(t)
	other, id := lab.inProgress("Other"), lab.inProgress("Asks")
	ended := time.Now().UTC().Truncate(time.Second)
	run := AgentRun{Story: id, Agent: "builder-" + id, Started: ended.Add(-time.Hour).Format(time.RFC3339), Ended: ended.Format(time.RFC3339), Outcome: OutcomeAsked, Thread: "TH-0001"}
	before, err := messages.Send(lab.repo, messages.SendOptions{From: other, To: id, Author: "builder-" + other, Text: "Before.", Now: ended.Add(-time.Second)})
	if err != nil {
		t.Fatal(err)
	}
	answeredThen, err := messages.Send(lab.repo, messages.SendOptions{From: other, To: id, Author: "builder-" + other, Text: "Answered then.", Now: ended})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := messages.Reply(lab.repo, answeredThen.ID, id, run.Agent, "Noted.", ended); err != nil {
		t.Fatal(err)
	}
	if got := messaged(lab.repo, &run); got != "" {
		t.Fatalf("%s and %s are not new to it, yet %q is", before.ID, answeredThen.ID, got)
	}
	then, err := messages.Send(lab.repo, messages.SendOptions{From: other, To: id, Author: "builder-" + other, Text: "Then.", Now: ended})
	if err != nil {
		t.Fatal(err)
	}
	if got := answered(lab.repo, &run); got != then.ID {
		t.Errorf("answered by %q, not %s", got, then.ID)
	}
}

// S-0116, ADR-0043: the operator restarts the agent of a story in ready or in
// progress whose agent dropped or failed, and is told why when it cannot be.
func TestAStoryWhoseAgentDroppedOrFailedIsRestarted(t *testing.T) {
	ctx := context.Background()
	restart := func(lab *agentLab, id string) (*AgentRun, error) {
		return Restart(ctx, lab.o, Entry{Key: "t", Name: "t", Root: lab.root}, id)
	}
	refusedFor := func(t *testing.T, lab *agentLab, id, want string) {
		t.Helper()
		_, err := restart(lab, id)
		var no *Refused
		if !errors.As(err, &no) || !strings.Contains(no.Why, want) {
			t.Errorf("restart %s: %v, want refused for %q", id, err, want)
		}
	}
	t.Run("a failed agent of a story in progress", func(t *testing.T) {
		lab := newAgentLab(t)
		id := lab.ready("A")
		lab.l.look(ctx, false)
		waitFor(t, "its agent ended", func() bool { r := lab.run(id); return r != nil && r.Ended != "" })
		lab.move(id, workitem.InProgress) // the agent pulled it, then failed
		first := lab.run(id)
		run, err := restart(lab, id)
		if err != nil {
			t.Fatal(err)
		}
		if run.Session == first.Session || run.Agent != "builder-"+id || run.PID == 0 {
			t.Errorf("a new session for the same agent: %+v, was %+v", run, first)
		}
		waitFor(t, "the new agent ran", func() bool { r := lab.run(id); return r.Started != first.Started || r.Session != first.Session })
		if j := lab.entries(); len(j) != 2 || !strings.Contains(j[1].Detail, "started stub-agent (command) for "+id) {
			t.Errorf("journal: %+v", j)
		}
	})
	t.Run("an agent whose process is gone", func(t *testing.T) {
		lab := newAgentLab(t)
		id := lab.ready("A")
		lab.move(id, workitem.InProgress)
		lab.l.dir.updateAgent(lab.root, func(s *AgentState) {
			s.put(&AgentRun{Story: id, Agent: "builder-" + id, PID: 999999999, Started: time.Now().UTC().Format(time.RFC3339), Command: "stub-agent"})
		})
		if _, err := restart(lab, id); err != nil {
			t.Fatal(err)
		}
	})
	// S-0118: a retry for a story in ready while the limit is full waits for
	// room, and flai serve starts it as it starts a story entering ready
	t.Run("a story in ready while the in-progress limit is full is queued", func(t *testing.T) {
		lab := newAgentLab(t)
		lab.hold()
		lab.limit(1)
		id := lab.ready("A")
		lab.l.look(ctx, false)
		waitFor(t, "it runs", func() bool { return lab.run(id).live() })
		lab.release(id)
		waitFor(t, "it ends in ready", func() bool { return !lab.run(id).live() })
		busy := lab.ready("Busy")
		lab.move(busy, workitem.InProgress)
		first := lab.run(id)

		run, err := restart(lab, id)
		if err != nil {
			t.Fatal(err)
		}
		if run.Queued == "" || run.Started != first.Started || run.PID != first.PID {
			t.Errorf("queued, not started: %+v, was %+v", run, first)
		}
		if r := lab.run(id); r.Queued == "" || r.Started != first.Started {
			t.Errorf("recorded: %+v", r)
		}
		a := Activity(lab.root, lab.state())[id]
		if a.State != ActivityWaiting || !strings.Contains(a.Why, "queued") {
			t.Errorf("a queued agent waits: %+v", a)
		}
		if j := lab.entries(); len(j) != 2 || !strings.Contains(j[1].Detail, "queued another agent for "+id) {
			t.Errorf("journal: %+v", j)
		}
		refusedFor(t, lab, id, "already queued")

		// no room: the launcher leaves it queued
		lab.l.look(ctx, false)
		if r := lab.run(id); r.Started != first.Started || r.Queued == "" {
			t.Errorf("started with no room: %+v", r)
		}
		if w := lab.state().Waiting; !strings.Contains(w, "in-progress limit leaves no room for "+id) {
			t.Errorf("waiting: %q", w)
		}
		// room: it starts, and the new run is no longer queued
		lab.move(busy, workitem.Review)
		lab.l.look(ctx, false)
		waitFor(t, "the queued agent started", func() bool { return lab.run(id).Started != first.Started || lab.run(id).Session != first.Session })
		if r := lab.run(id); r.Queued != "" || r.PID == 0 {
			t.Errorf("the new run: %+v", r)
		}
		lab.release(id)
		waitFor(t, "it ends", func() bool { return !lab.run(id).live() })
		// ended again, in ready, not queued: the launcher holds it as before
		if a := Activity(lab.root, lab.state())[id]; a.State != ActivityFailed {
			t.Errorf("after the queued agent failed: %+v", a)
		}
	})
	t.Run("a queued story that leaves ready is not waiting", func(t *testing.T) {
		lab := newAgentLab(t)
		id := lab.ready("A")
		lab.move(id, workitem.InProgress)
		lab.l.dir.updateAgent(lab.root, func(s *AgentState) {
			s.put(&AgentRun{Story: id, Agent: "builder-" + id, Started: "2026-09-26T00:00:00Z", Ended: "2026-09-26T00:01:00Z",
				Outcome: OutcomeFailed, Why: "ended (exit 1)", Queued: "2026-09-26T00:02:00Z"})
		})
		if a := Activity(lab.root, lab.state())[id]; a.State != ActivityFailed {
			t.Errorf("in progress, a queued retry is moot: %+v", a)
		}
		// in progress, a retry starts at once, queued or not
		run, err := restart(lab, id)
		if err != nil {
			t.Fatal(err)
		}
		if run.Queued != "" || run.PID == 0 {
			t.Errorf("started: %+v", run)
		}
	})
	// S-0177, ADR-0064: a story begun on another host, which this host has
	// had no agent for, gets one told where, when, and by whom
	t.Run("a story in progress begun on another host", func(t *testing.T) {
		lab := newAgentLab(t)
		lab.hold()
		id := lab.ready("Elsewhere")
		lab.move(id, workitem.InProgress) // by its agent, on the other host
		there := &workitem.Narrative{Stream: id, Title: "Elsewhere", Updated: lab.now.UTC().Format(workitem.TimeFormat),
			Agent: "builder-" + id, Session: "there", Host: "far-away", Path: lab.repo.NarrativePath(id), Body: "\n# " + id + " Elsewhere\n"}
		if err := there.Save(); err != nil {
			t.Fatal(err)
		}
		asked, err := threads.New(lab.repo, threads.NewOptions{Title: "Which port?", On: id, Author: "builder-" + id, Text: "Eight or nine?", Now: time.Now()})
		if err != nil {
			t.Fatal(err)
		}
		if _, err := threads.Reply(lab.repo, asked.ID, "alex", "Nine.", time.Now()); err != nil {
			t.Fatal(err)
		}
		if _, err := threads.New(lab.repo, threads.NewOptions{Title: "Unanswered", On: id, Author: "builder-" + id, Text: "And then?", Now: time.Now()}); err != nil {
			t.Fatal(err)
		}
		it, _ := lab.repo.Get(id)
		b := begun(lab.repo, it)
		if b.By != "builder-"+id || b.At == "" || b.Agent != "builder-"+id || b.Host != "far-away" || len(b.Threads) != 1 || b.Threads[0] != asked.ID {
			t.Errorf("where it was begun, and the thread answered since: %+v", b)
		}

		// the dashboard is told it was begun elsewhere, and has no agent here
		a := Activity(lab.root, lab.state())[id]
		if a.State != ActivityWaiting || a.Elsewhere == nil || a.Elsewhere.Host != "far-away" || a.Elsewhere.Here || a.Run == nil || a.Run.Started != "" ||
			a.Why != "begun by builder-"+id+" on far-away at "+b.At+"; no agent here" {
			t.Errorf("a story begun elsewhere: %+v %+v", a, a.Elsewhere)
		}
		there.Host = workitem.ThisHost()
		_ = there.Save()
		if a := Activity(lab.root, lab.state())[id]; a.Elsewhere == nil || !a.Elsewhere.Here || !strings.Contains(a.Why, "on this host, outside flai serve") {
			t.Errorf("a story begun on this host outside flai serve: %+v", a)
		}

		run, err := restart(lab, id)
		if err != nil {
			t.Fatal(err)
		}
		// with an agent here it is this host's; the question its namesake asked
		// there is its own
		if a := Activity(lab.root, lab.state())[id]; a.Elsewhere != nil || a.State != ActivityWaiting || a.Thread == "" {
			t.Errorf("with an agent here it is this host's: %+v", a)
		}
		if run.PID == 0 || run.Agent != "builder-"+id || run.Queued != "" {
			t.Errorf("started at once: %+v", run)
		}
		if j := lab.entries(); len(j) != 1 || !strings.Contains(j[0].Detail, "begun by builder-"+id) || !strings.Contains(j[0].Detail, "on the host "+workitem.ThisHost()) {
			t.Errorf("journal: %+v", j)
		}
		waitFor(t, "it runs", func() bool { return lab.run(id).running() })
		// this host's agent now runs for it, and is not started twice
		refusedFor(t, lab, id, "agent is running")
	})
	t.Run("refusals", func(t *testing.T) {
		lab := newAgentLab(t)
		lab.hold()
		none := lab.backlog("Never started", nil)
		refusedFor(t, lab, none, "is in backlog")
		lab.toReady(none)
		lab.limit(1)
		busy := lab.ready("Busy")
		lab.move(busy, workitem.InProgress)
		// one that never had an agent here is queued, as a retry is (S-0177)
		if run, err := restart(lab, none); err != nil || run.Queued == "" || run.Started != "" {
			t.Errorf("a story in ready with no agent here and no room is queued: %+v %v", run, err)
		}
		refusedFor(t, lab, none, "already queued")
		// running
		lab.limit(5)
		running := lab.ready("Running")
		lab.l.look(ctx, false)
		waitFor(t, "it runs", func() bool { return lab.run(running).live() })
		refusedFor(t, lab, running, "agent is running")
		lab.release(running)
		waitFor(t, "it ends", func() bool { return !lab.run(running).live() })
		lab.limit(1)
		// asked
		lab.l.dir.updateAgent(lab.root, func(s *AgentState) {
			r := *s.Stories[running]
			r.Outcome, r.Thread = OutcomeAsked, "TH-0001"
			s.put(&r)
		})
		refusedFor(t, lab, running, "waiting for an answer to TH-0001")
		// nothing to start it with
		lab.limit(5)
		lab.l.dir.updateAgent(lab.root, func(s *AgentState) {
			r := *s.Stories[running]
			r.Outcome = OutcomeFailed
			s.put(&r)
		})
		lab.cfg.Command = nil
		refusedFor(t, lab, running, "names no harness, and no command is set")
		// the action off
		lab.cfg.Enabled = false
		refusedFor(t, lab, running, "agent host action is off")
	})
}

// S-0220: a recommendation on the question an agent ended asking is no
// answer until the operator confirms it (ADR-0090): the agent is not started
// again on it, but on the confirmation; and an answer citing a source, as the
// orchestrator gives with answer_threads autonomous, starts it again too.
func TestAnAgentThatEndedAskingIsStartedAgainOnAConfirmationNotARecommendation(t *testing.T) {
	lab := newAgentLab(t)
	ctx := context.Background()
	lab.hold()
	id := lab.ready("Asks and is recommended")
	lab.l.look(ctx, false)
	waitFor(t, "it runs", func() bool { return lab.run(id).live() })
	agent := lab.run(id).Agent
	lab.move(id, workitem.InProgress)
	th, err := threads.New(lab.repo, threads.NewOptions{Title: "Which port?", On: id, Author: agent, Text: "Eight or nine?", Now: time.Now()})
	if err != nil {
		t.Fatal(err)
	}
	// recommended while it still runs: it is waiting all the same
	if _, err := threads.ReplyWith(lab.repo, th.ID, workitem.ActivityOrchestrator, "Nine.", time.Now(), threads.Marks{Recommendation: true, Source: threads.Source{Path: id}}); err != nil {
		t.Fatal(err)
	}
	if a := Activity(lab.root, lab.state())[id]; a.State != ActivityWaiting || a.Thread != th.ID {
		t.Fatalf("a recommendation leaves the agent waiting: %+v", a)
	}
	lab.release(id)
	waitFor(t, "it ends", func() bool { return !lab.run(id).live() })
	if r := lab.run(id); r.Outcome != OutcomeAsked || r.Thread != th.ID {
		t.Fatalf("it ended asking, with a recommendation pending: %+v", r)
	}
	_ = os.Remove(filepath.Join(lab.outDir, "release-"+id))
	lab.l.look(ctx, false)
	if lab.run(id).live() {
		t.Fatal("started again on a recommendation")
	}
	if _, err := threads.Confirm(lab.repo, th.ID, "alex", time.Now()); err != nil {
		t.Fatal(err)
	}
	lab.l.look(ctx, false)
	waitFor(t, "it runs again on the confirmation", func() bool { return lab.run(id).live() })
	if r := lab.run(id); r.Answered != th.ID {
		t.Errorf("started again for the confirmed recommendation: %+v", r)
	}

	// asked again, it is answered by the orchestrator with a source
	again, err := threads.New(lab.repo, threads.NewOptions{Title: "Which host?", On: id, Author: agent, Text: "Here or there?", Now: time.Now()})
	if err != nil {
		t.Fatal(err)
	}
	lab.release(id)
	waitFor(t, "it ends again", func() bool { return !lab.run(id).live() })
	if r := lab.run(id); r.Outcome != OutcomeAsked || r.Thread != again.ID {
		t.Fatalf("it ended asking again: %+v", r)
	}
	_ = os.Remove(filepath.Join(lab.outDir, "release-"+id))
	if _, err := threads.ReplyWith(lab.repo, again.ID, workitem.ActivityOrchestrator, "Here.", time.Now(), threads.Marks{Source: threads.Source{Path: id}}); err != nil {
		t.Fatal(err)
	}
	lab.l.look(ctx, false)
	waitFor(t, "it runs again on the orchestrator's answer", func() bool { return lab.run(id).live() })
	if r := lab.run(id); r.Answered != again.ID {
		t.Errorf("started again for the orchestrator's answer: %+v", r)
	}
	lab.release(id)
	waitFor(t, "it ends at last", func() bool { return !lab.run(id).live() })
}

func TestTheOrchestratorsAndTheAnalyzersEarlierRunsAreKeptNewestFirst(t *testing.T) {
	var s AgentState
	for i := 1; i <= pastRuns+2; i++ {
		started := fmt.Sprintf("2026-10-07T00:%02d:00Z", i)
		s.put(&AgentRun{Agent: "orchestrator", Started: started, PID: i})
		s.put(&AgentRun{Agent: "analyzer", Focus: "risk", Started: started, PID: i})
	}
	// The newest run recorded again, as when it ends, is no earlier run.
	ended := *s.Orchestrator
	ended.Ended = "2026-10-07T01:00:00Z"
	s.put(&ended)
	for name, past := range map[string][]*AgentRun{"orchestrator": s.PastOrchestrators, "analyzer": s.PastAnalyzers} {
		if len(past) != pastRuns {
			t.Fatalf("%s: kept %d earlier runs, want %d", name, len(past), pastRuns)
		}
		if past[0].PID != pastRuns+1 || past[pastRuns-1].PID != 2 {
			t.Errorf("%s: earlier runs from pid %d to %d, want %d newest first to 2", name, past[0].PID, past[pastRuns-1].PID, pastRuns+1)
		}
	}
	if s.Orchestrator.PID != pastRuns+2 || s.Orchestrator.Ended == "" {
		t.Errorf("newest orchestrator run = %+v, want pid %d ended", s.Orchestrator, pastRuns+2)
	}
}

// S-0294, ADR-0108, I-0084: a story's agent that ends with its story in
// progress, unblocked and asking nothing, is restarted at the next look, in
// a new session, as the same agent, told that flai serve restarted it on its
// own.
func TestAgentRestartedAfterEndingInProgress(t *testing.T) {
	lab := newAgentLab(t)
	ctx := context.Background()
	lab.cfg.AutoRestarts = 2
	lab.cfg.Command = nil
	lab.cfg.Harnesses = map[string]harness.Host{harness.ClaudeCode: {Program: lab.stub}}
	lab.hold()
	id := lab.readyWith("Ends early", &manifest.Agent{Harness: harness.ClaudeCode})
	lab.l.look(ctx, false)
	waitFor(t, "it runs", func() bool { return lab.run(id).live() })
	lab.move(id, workitem.InProgress) // its agent pulled it
	first := lab.run(id)
	// claude -p ends it while its story is in progress (I-0084)
	lab.release(id)
	waitFor(t, "it ends in progress", func() bool { r := lab.run(id); return !r.live() && r.Outcome == OutcomeFailed })

	lab.l.look(ctx, false)
	waitFor(t, "another agent is started", func() bool { return lab.run(id).Session != first.Session })
	again := lab.run(id)
	if again.Agent != first.Agent || again.AutoRestarts != 1 || again.Answered != "" || again.PID == 0 || first.AutoRestarts != 0 {
		t.Errorf("the same agent, in a new session, restart 1: %+v, first %+v", again, first)
	}
	if j := lab.entries(); len(j) != 2 || !strings.Contains(j[1].Detail, "again for "+id+" as "+first.Agent+" on its own, automatic restart 1 of 2") {
		t.Errorf("journal: %+v", j)
	}
	waitFor(t, "it ends again", func() bool { return !lab.run(id).live() })
	got, _ := os.ReadFile(filepath.Join(lab.outDir, id+".txt"))
	for _, want := range []string{"its last agent ended (exit 0) with " + id + " in in-progress", "flai serve started it again on its own, automatic restart 1 of 2"} {
		if !strings.Contains(string(got), want) {
			t.Errorf("the agent was told\n%s\nwithout %q", got, want)
		}
	}
	if strings.Contains(string(got), "the operator restarted it") {
		t.Errorf("the agent was told the operator restarted it:\n%s", got)
	}
}

// limitThreads are the threads flai serve opened on story at the limit of
// automatic restarts.
func (lab *agentLab) limitThreads(story string) []*threads.Thread {
	lab.t.Helper()
	all, err := threads.List(lab.repo)
	if err != nil {
		lab.t.Fatal(err)
	}
	var out []*threads.Thread
	for _, th := range all {
		if th.Opener() == limitAuthor && threads.StoryOf(lab.repo, th) == story {
			out = append(out, th)
		}
	}
	return out
}

// S-0294, ADR-0108: flai serve restarts a story's agent up to the limit,
// then opens one thread to the operator, and restarts it no more; the
// operator's restart starts the count again from 0.
func TestAgentRestartsStopAtTheLimitWithOneThread(t *testing.T) {
	lab := newAgentLab(t)
	ctx := context.Background()
	lab.cfg.AutoRestarts = 2
	lab.hold()
	id := lab.ready("Ends every time")
	lab.l.look(ctx, false)
	waitFor(t, "it runs", func() bool { return lab.run(id).live() })
	lab.move(id, workitem.InProgress)
	lab.release(id) // every agent for it ends at once
	for n := 1; n <= 2; n++ {
		waitFor(t, "it ends", func() bool { r := lab.run(id); return !r.live() && r.Outcome == OutcomeFailed })
		was := lab.run(id).Session
		lab.l.look(ctx, false)
		waitFor(t, "it is restarted", func() bool { return lab.run(id).Session != was })
		if r := lab.run(id); r.AutoRestarts != n {
			t.Fatalf("restart %d: %+v", n, r)
		}
	}
	waitFor(t, "it ends a third time", func() bool { r := lab.run(id); return !r.live() && r.Outcome == OutcomeFailed })
	last := lab.run(id)
	lab.l.look(ctx, false)
	lab.l.look(ctx, false)
	r := lab.run(id)
	if r.Session != last.Session || r.live() || r.AutoRestarts != 2 {
		t.Fatalf("restarted past the limit: %+v", r)
	}
	ths := lab.limitThreads(id)
	if len(ths) != 1 {
		t.Fatalf("%d threads at the limit, want 1", len(ths))
	}
	th := ths[0]
	if r.LimitThread != th.ID || !strings.Contains(th.Title, id+"'s agent ended 3 times") {
		t.Errorf("thread %s %q, run %+v", th.ID, th.Title, r)
	}
	text := th.Entries()[0].Text
	for _, want := range []string{"ended 3 times in a row", "restarted it on its own 2 times", "Its last run ended (exit 0) with " + id + " in in-progress", "flai serve agent restart " + id, "Retry", "--auto-restarts"} {
		if !strings.Contains(text, want) {
			t.Errorf("the thread says\n%s\nwithout %q", text, want)
		}
	}
	started := 0
	for _, e := range lab.entries() {
		if strings.HasPrefix(e.Detail, "started ") {
			started++
		}
	}
	if started != 3 {
		t.Errorf("%d agents started, want 3: %+v", started, lab.entries())
	}
	if a := Activity(lab.root, lab.state())[id]; a.State != ActivityFailed {
		t.Errorf("at the limit it reads failed: %+v", a)
	}

	// the operator restarts it: a run of the count 0, with no thread
	run, err := Restart(ctx, lab.o, Entry{Key: "t", Name: "t", Root: lab.root}, id)
	if err != nil {
		t.Fatal(err)
	}
	if run.AutoRestarts != 0 || run.LimitThread != "" || run.Session == last.Session {
		t.Fatalf("the operator's restart: %+v", run)
	}
	// it ends in progress, as seen at a look: restarted on its own again, as the first of 2
	ended := *run
	ended.PID, ended.Ended, ended.Outcome, ended.Why = 0, time.Now().UTC().Format(time.RFC3339), OutcomeFailed, "ended (exit 0) with "+id+" in in-progress"
	lab.l.dir.updateAgent(lab.root, func(s *AgentState) { s.put(&ended) })
	lab.l.look(ctx, false)
	waitFor(t, "it is restarted on its own", func() bool { return lab.run(id).Session != run.Session })
	if r := lab.run(id); r.AutoRestarts != 1 {
		t.Errorf("the count starts again from 0: %+v", r)
	}
	waitFor(t, "it ends", func() bool { return !lab.run(id).live() })
}

// putEnded makes a story in progress whose agent's run ended as run says,
// with its story, agent, and times filled in.
func (lab *agentLab) putEnded(title string, run AgentRun) string {
	lab.t.Helper()
	id := lab.ready(title)
	lab.move(id, workitem.InProgress)
	now := time.Now().UTC().Format(time.RFC3339)
	run.Story, run.Agent, run.Started, run.Ended, run.Session = id, "builder-"+id, now, now, "earlier"
	if run.Why == "" && run.Outcome == OutcomeFailed {
		run.Why = "ended (exit 0) with " + id + " in in-progress"
	}
	lab.l.dir.updateAgent(lab.root, func(s *AgentState) { s.put(&run) })
	return id
}

// S-0294, ADR-0108: with agent.auto_restarts 0, an agent that ends in
// progress is not restarted, and the thread is opened at once, once.
func TestAgentRestartsOffOpensTheThreadAtTheFirstEnd(t *testing.T) {
	lab := newAgentLab(t)
	ctx := context.Background()
	lab.cfg.AutoRestarts = 0
	id := lab.putEnded("Ends with restarts off", AgentRun{Outcome: OutcomeFailed})
	lab.l.look(ctx, false)
	lab.l.look(ctx, false)
	r := lab.run(id)
	ths := lab.limitThreads(id)
	if r.Session != "earlier" || len(ths) != 1 || r.LimitThread != ths[0].ID {
		t.Fatalf("run %+v, threads %d", r, len(ths))
	}
	if text := ths[0].Entries()[0].Text; !strings.Contains(text, "ended once in a row") || !strings.Contains(text, "agent.auto_restarts is 0") {
		t.Errorf("the thread says\n%s", text)
	}
	if j := lab.entries(); len(j) != 1 || j[0].Outcome != "done" || !strings.Contains(j[0].Detail, "opened "+ths[0].ID) {
		t.Errorf("journal: %+v", j)
	}
}

// S-0294, ADR-0108: what flai serve does not restart on its own, and opens
// no thread for.
func TestAnEndedAgentIsNotRestartedOnItsOwnWhenItShouldNotBe(t *testing.T) {
	ctx := context.Background()
	for _, c := range []struct {
		name string
		run  AgentRun
		set  func(lab *agentLab, id string)
	}{
		{"its story is blocked", AgentRun{Outcome: OutcomeFailed}, func(lab *agentLab, id string) {
			st, _ := lab.repo.Get(id)
			if err := workitem.BlockItem(st, "the API key", time.Now()); err != nil {
				lab.t.Fatal(err)
			}
			if err := lab.repo.Save(st); err != nil {
				lab.t.Fatal(err)
			}
		}},
		{"the operator stopped it", AgentRun{Outcome: OutcomeStopped, Stopped: "2026-10-06T00:00:00Z"}, nil},
		{"it could not be started", AgentRun{Outcome: OutcomeFailed, Error: "no such program"}, nil},
		{"it ended asking, unanswered", AgentRun{Outcome: OutcomeAsked}, func(lab *agentLab, id string) {
			th, err := threads.New(lab.repo, threads.NewOptions{Title: "Which port?", On: id, Author: "builder-" + id, Text: "Eight or nine?", Now: time.Now()})
			if err != nil {
				lab.t.Fatal(err)
			}
			lab.l.dir.updateAgent(lab.root, func(s *AgentState) { s.Stories[id].Thread = th.ID })
		}},
		{"it ended in ready", AgentRun{Outcome: OutcomeFailed}, func(lab *agentLab, id string) {
			st, _ := lab.repo.Get(id)
			if _, err := lab.repo.Transition(st, workitem.Ready, "alex", "not yet", lab.now); err != nil {
				lab.t.Fatal(err)
			}
		}},
		{"the agent action is off", AgentRun{Outcome: OutcomeFailed}, func(lab *agentLab, _ string) { lab.cfg.Enabled = false }},
		{"its thread at the limit is open", AgentRun{Outcome: OutcomeFailed, AutoRestarts: 2, LimitThread: "TH-0001"}, nil},
	} {
		t.Run(c.name, func(t *testing.T) {
			lab := newAgentLab(t)
			lab.cfg.AutoRestarts = 2
			id := lab.putEnded("Ended", c.run)
			if c.set != nil {
				c.set(lab, id)
			}
			lab.l.look(ctx, false)
			if r := lab.run(id); r.Session != "earlier" || r.live() {
				t.Errorf("restarted: %+v", r)
			}
			if j := lab.entries(); len(j) != 0 {
				t.Errorf("journal: %+v", j)
			}
			if ths := lab.limitThreads(id); len(ths) != 0 {
				t.Errorf("%d threads opened", len(ths))
			}
		})
	}
}
