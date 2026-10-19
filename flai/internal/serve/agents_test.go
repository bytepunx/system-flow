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
	"github.com/bytepunx/system-flow/flai/internal/forge"
	"github.com/bytepunx/system-flow/flai/internal/harness"
	"github.com/bytepunx/system-flow/flai/internal/host"
	"github.com/bytepunx/system-flow/flai/internal/hostapi"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/mergequeue"
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
		"{ echo \"args: $*\"; echo \"dir: $(pwd)\"; echo \"agent: $FLAI_AGENT\"; echo \"story: $FLAI_STORY\"; echo \"session: $FLAI_SESSION\"; echo \"answered: $FLAI_ANSWERED\"; echo \"host: ${FLAI_HOST_URL-unset} ${FLAI_HOST_TOKEN-unset}\"; echo \"config: ${FLAI_CONFIG-unset}\"; echo \"gateway: ${ANTHROPIC_BASE_URL-unset} ${ANTHROPIC_AUTH_TOKEN-unset}\"; echo \"forge: ${GITHUB_TOKEN-unset} ${S0393_FORGE_TOKEN-unset} ${S0393_KEPT-unset}\"; } > \"$out\"\n" +
		"while [ -f \"" + lab.outDir + "/hold\" ] && [ ! -f \"" + lab.outDir + "/release-$FLAI_STORY\" ]; do sleep 0.05; done\n" +
		"touch \"" + lab.outDir + "/ended-$$\"\n"
	if err := os.WriteFile(lab.stub, []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	// The command cannot ask before a protected write: the operator says
	// their arguments deny the protected paths, so it starts (ADR-0131).
	lab.cfg = AgentConfig{Enabled: true, Command: []string{lab.stub, "work on {story}", "--root={root}", "$(echo not a shell)"}, Name: "builder",
		Lifts: map[string]HarnessLifts{harness.Command: {DenyProtected: true}}}
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

// S-0393: an agent gets no forge's token, neither in the forge's default
// variable nor in the one the host names, whatever the project's acceptance
// mode, and keeps the rest of its environment.
func TestAnAgentIsStartedWithoutAForgesToken(t *testing.T) {
	t.Setenv(config.DefaultGitHubTokenEnv, "ghp_default-never-given")
	t.Setenv("S0393_FORGE_TOKEN", "ghp_named-never-given")
	t.Setenv("S0393_KEPT", "kept")
	lab := newAgentLab(t)
	lab.cfg.ForgeTokenEnvs = []string{"S0393_FORGE_TOKEN"}
	id := lab.ready("Tokenless")
	lab.l.look(context.Background(), false)
	waitFor(t, "the stub ended", func() bool { r := lab.run(id); return r != nil && r.Ended != "" })
	got, _ := os.ReadFile(filepath.Join(lab.outDir, id+".txt"))
	if !strings.Contains(string(got), "forge: unset unset kept\n") || strings.Contains(string(got), "ghp_") {
		t.Errorf("the stub was given\n%s\nwant neither forge token and the rest kept", got)
	}
}

// The forges' tokens are left out of an environment whether or not the host
// names a variable for one, and with no say from the host at all.
func TestWithoutForgeTokens(t *testing.T) {
	env := []string{"PATH=/bin", "GITHUB_TOKEN=d", "GH_PAT=n", "GITHUB_TOKEN_X=x", "HOME=/h"}
	lab := newAgentLab(t)
	lab.cfg.ForgeTokenEnvs = []string{"GH_PAT"}
	if got := lab.l.withoutForgeTokens(slices.Clone(env)); !slices.Equal(got, []string{"PATH=/bin", "GITHUB_TOKEN_X=x", "HOME=/h"}) {
		t.Errorf("named: env = %v, want both tokens left out and the rest kept", got)
	}
	lab.cfg.ForgeTokenEnvs = nil
	if got := lab.l.withoutForgeTokens(slices.Clone(env)); !slices.Equal(got, []string{"PATH=/bin", "GH_PAT=n", "GITHUB_TOKEN_X=x", "HOME=/h"}) {
		t.Errorf("default: env = %v, want GITHUB_TOKEN left out", got)
	}
	bare := &launcher{}
	if got := bare.withoutForgeTokens(slices.Clone(env)); !slices.Equal(got, []string{"PATH=/bin", "GH_PAT=n", "GITHUB_TOKEN_X=x", "HOME=/h"}) {
		t.Errorf("no config: env = %v, want GITHUB_TOKEN left out", got)
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

// settled writes story's settled reply on conversation id, by its agent.
func (lab *agentLab) settled(id, story, text string) {
	lab.t.Helper()
	if _, err := messages.ReplyWith(lab.repo, id, story, "builder-"+story, text, time.Now(), messages.Marks{Settled: true}); err != nil {
		lab.t.Fatal(err)
	}
}

// asksOn opens a question by story's agent on story, as its run ended
// asking, and records it on the run.
func (lab *agentLab) asksOn(story string, at time.Time) *threads.Thread {
	lab.t.Helper()
	th, err := threads.New(lab.repo, threads.NewOptions{Title: "Which port?", On: story, Author: "builder-" + story, Text: "Eight or nine?", Now: at})
	if err != nil {
		lab.t.Fatal(err)
	}
	lab.l.dir.updateAgent(lab.root, func(s *AgentState) { s.Stories[story].Thread = th.ID })
	return th
}

// S-0372, I-0125: a settled reply to a story whose agent ended waiting on the
// operator's answer does not start it again, as the agreement on MS-0109
// started S-0380's; a message that awaits its reply still does.
func TestASettledReplyDoesNotStartAnAgentThatEndedForAThread(t *testing.T) {
	lab := newAgentLab(t)
	ctx := context.Background()
	lab.hold()
	other := lab.inProgress("Other")
	id := lab.putEnded("Waits on the operator", AgentRun{Outcome: OutcomeAsked})
	before := time.Now().Add(-2 * time.Minute)
	lab.asksOn(id, before)
	c, err := messages.Send(lab.repo, messages.SendOptions{From: other, To: id, Author: "builder-" + other, Text: "I take docs/, you take flai/.", Now: before})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := messages.ReplyWith(lab.repo, c.ID, id, "builder-"+id, "Agreed; nothing to change.", before, messages.Marks{Settled: true}); err != nil {
		t.Fatal(err)
	}
	lab.endedAgo(id, time.Minute)

	lab.settled(c.ID, other, "Agreed.")
	lab.l.look(ctx, false)
	if r := lab.run(id); r.live() || r.Session != "earlier" || len(lab.entries()) != 0 {
		t.Fatalf("started again for a settled reply: %+v, journal %+v", r, lab.entries())
	}

	if _, err := messages.Reply(lab.repo, c.ID, other, "builder-"+other, "One more: I renamed docs/ports.md.", time.Now()); err != nil {
		t.Fatal(err)
	}
	lab.l.look(ctx, false)
	waitFor(t, "it runs again for the message", func() bool { return lab.run(id).Answered == c.ID })
	lab.release(id)
	waitFor(t, "it ends", func() bool { return !lab.run(id).live() })
}

// S-0372, I-0125: an agent that ended waiting on another story's agent's
// reply is started again once when that reply comes settled, and not again
// when it answers settled in turn; nor is the other story's agent, which
// ended waiting on the operator, started for that answer. A run that ends
// with only settled conversations is not asked on them.
func TestASettledReplyStartsTheAgentWaitingOnItOnce(t *testing.T) {
	lab := newAgentLab(t)
	ctx := context.Background()
	lab.hold()
	other := lab.putEnded("Other", AgentRun{Outcome: OutcomeAsked})
	lab.asksOn(other, time.Now())
	id := lab.putEnded("Messages and waits", AgentRun{Outcome: OutcomeAsked})
	run := lab.run(id)
	started, _ := time.Parse(time.RFC3339, run.Started)
	th := lab.asksOn(id, started)
	c, err := messages.Send(lab.repo, messages.SendOptions{From: id, To: other, Author: run.Agent, Text: "Do you change flai/cmd?", Now: started})
	if err != nil {
		t.Fatal(err)
	}
	if judgeRun(lab.root, run, new(int)); run.Outcome != OutcomeAsked || run.Thread != th.ID || !slices.Equal(run.Conversations, []string{c.ID}) {
		t.Fatalf("asked on both: %+v", run)
	}
	lab.l.dir.updateAgent(lab.root, func(s *AgentState) { s.put(run) })
	lab.l.look(ctx, false)
	if len(lab.entries()) != 0 {
		t.Fatalf("started before the reply: %+v", lab.entries())
	}

	lab.settled(c.ID, other, "No; nothing to change.")
	lab.l.look(ctx, false)
	waitFor(t, "it runs again for the settled reply", func() bool { return lab.run(id).Answered == c.ID })
	if again := lab.run(id); again.Session != "earlier" {
		t.Errorf("the same session: %+v", again)
	}
	lab.settled(c.ID, id, "Settled; no reply is needed.")
	lab.release(id)
	waitFor(t, "it ends", func() bool { return !lab.run(id).live() })
	if r := lab.run(id); r.Outcome != OutcomeAsked || r.Thread != th.ID || len(r.Conversations) != 0 {
		t.Fatalf("asked on the settled conversation: %+v", r)
	}
	_ = os.Remove(filepath.Join(lab.outDir, "release-"+id))
	lab.l.look(ctx, false)
	if r, o := lab.run(id), lab.run(other); r.live() || o.live() || len(lab.entries()) != 1 {
		t.Fatalf("started again for a settled answer: %+v, %+v, journal %+v", r, o, lab.entries())
	}

	// with no question of its own open, it ends neither asking nor waiting
	alone := AgentRun{Story: id, Agent: "someone-else", Started: started.Format(time.RFC3339)}
	if judgeRun(lab.root, &alone, new(int)); alone.Outcome != OutcomeFailed || len(alone.Conversations) != 0 {
		t.Errorf("judged on a settled conversation: %+v", alone)
	}
}

// S-0372, I-0125: while its story is blocked, an agent that ended asking is
// not started again for a new message to its story, as S-0337's was twelve
// times, nor for the conversation it waited on closing; the answer to the
// question it asked still starts it.
func TestABlockedStorysAgentIsStartedOnlyForTheAnswerToItsQuestion(t *testing.T) {
	lab := newAgentLab(t)
	ctx := context.Background()
	lab.hold()
	other := lab.inProgress("Other")
	id := lab.putEnded("Blocked and waits", AgentRun{Outcome: OutcomeAsked})
	before := time.Now().Add(-2 * time.Minute)
	th := lab.asksOn(id, before)
	c, err := messages.Send(lab.repo, messages.SendOptions{From: id, To: other, Author: "builder-" + id, Text: "Do you change flai/cmd?", Now: before})
	if err != nil {
		t.Fatal(err)
	}
	lab.l.dir.updateAgent(lab.root, func(s *AgentState) {
		s.Stories[id].Conversations, s.Stories[id].WaitsOn = []string{c.ID}, []string{other}
	})
	lab.endedAgo(id, 30*time.Second)
	st, _ := lab.repo.Get(id)
	if err := workitem.BlockItem(st, "waiting on "+th.ID, time.Now()); err != nil {
		t.Fatal(err)
	}
	if err := lab.repo.Save(st); err != nil {
		t.Fatal(err)
	}
	notStarted := func(what string) {
		t.Helper()
		lab.l.look(ctx, false)
		if r := lab.run(id); r.live() || r.Session != "earlier" || len(lab.entries()) != 0 {
			t.Fatalf("started while blocked for %s: %+v, journal %+v", what, r, lab.entries())
		}
	}

	if _, err := messages.Send(lab.repo, messages.SendOptions{From: other, To: id, Author: "builder-" + other, Text: "I changed flai/cmd/prime.go.", Now: time.Now()}); err != nil {
		t.Fatal(err)
	}
	notStarted("a new message")
	if _, err := messages.Close(lab.repo, c.ID, "builder-"+other, "nothing to change", time.Now()); err != nil {
		t.Fatal(err)
	}
	notStarted("a closed conversation")

	if _, err := threads.Reply(lab.repo, th.ID, "alex", "Nine.", time.Now()); err != nil {
		t.Fatal(err)
	}
	lab.l.look(ctx, false)
	waitFor(t, "it runs again for the answer", func() bool { return lab.run(id).Answered == th.ID })
	lab.release(id)
	waitFor(t, "it ends", func() bool { return !lab.run(id).live() })
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

// asks counts, by holding story, the messages flai wrote on held's
// conversations asking about its hold (ADR-0134).
func (lab *agentLab) asks(held string) map[string]int {
	lab.t.Helper()
	cs, err := messages.For(lab.repo, held)
	if err != nil {
		lab.t.Fatal(err)
	}
	out := map[string]int{}
	for _, c := range cs {
		for _, e := range c.Entries() {
			if e.Author == messages.AskAuthor {
				out[c.Other(held)]++
			}
		}
	}
	return out
}

// heldNow is why story is held, as the board reads it.
func (lab *agentLab) heldNow(story string) *workitem.Hold {
	lab.t.Helper()
	items, err := lab.repo.List(false)
	if err != nil {
		lab.t.Fatal(err)
	}
	it, err := lab.repo.Get(story)
	if err != nil {
		lab.t.Fatal(err)
	}
	return lab.repo.Holds(items).Of(it)
}

// logged counts the lines logged with msg.
func (lab *agentLab) logged(msg string) int { return strings.Count(lab.logText(), `msg="`+msg+`"`) }

// ADR-0134: at each look flai serve asks the agent of each story in progress
// that holds a ready story on overlap alone about the hold, once per pair
// while the held story stays in ready, the in-progress limit full or not:
// a second look writes nothing new.
func TestAHoldOnOverlapAloneIsAskedAboutOncePerPair(t *testing.T) {
	lab := newAgentLab(t)
	ctx := context.Background()
	lab.limit(1)
	holder := lab.readyTouching("Holder", "flai/cmd")
	lab.move(holder, workitem.InProgress)
	held := lab.readyTouching("Held", "flai/cmd/serve.go")
	lab.l.look(ctx, false)
	cs, err := messages.For(lab.repo, held)
	if err != nil {
		t.Fatal(err)
	}
	if len(cs) != 1 || cs[0].From != held || cs[0].To != holder || !slices.Equal(cs[0].About, []string{"flai/cmd/serve.go"}) {
		t.Fatalf("one conversation from %s to %s about the overlap: %+v", held, holder, cs)
	}
	if got := lab.asks(held); len(got) != 1 || got[holder] != 1 {
		t.Fatalf("asks: %v", got)
	}
	if n := lab.logged("hold asked about"); n != 1 || !strings.Contains(lab.logText(), `msg="hold asked about" component=serve project=t held=`+held+" holder="+holder+" conversation="+cs[0].ID) {
		t.Errorf("logged %d asks:\n%s", n, lab.logText())
	}
	lab.l.look(ctx, false)
	if got := lab.asks(held); len(got) != 1 || got[holder] != 1 {
		t.Errorf("a second look asked again: %v", got)
	}
	if again, _ := messages.Get(lab.repo, cs[0].ID); len(again.Entries()) != len(cs[0].Entries()) {
		t.Errorf("a second look wrote on %s: %d entries, then %d", cs[0].ID, len(cs[0].Entries()), len(again.Entries()))
	}
	if n := lab.logged("hold asked about"); n != 1 {
		t.Errorf("logged %d asks after two looks", n)
	}
	if lab.run(held) != nil {
		t.Errorf("a held story was started: %+v", lab.run(held))
	}
}

// ADR-0134: a story held by after:, or by a story with no touches, or one
// whose own claim is empty, is held on more than overlap and is not asked
// about; nor is a story held only by one whose agent was started and is
// still in ready.
func TestNoAskForAHoldThatIsNotOnOverlapAlone(t *testing.T) {
	ctx := context.Background()
	t.Run("after", func(t *testing.T) {
		lab := newAgentLab(t)
		holder := lab.readyTouching("Holder", "flai/cmd")
		lab.move(holder, workitem.InProgress)
		st, err := lab.repo.Get(lab.backlog("Waits", nil, "flai/cmd/serve.go"))
		if err != nil {
			t.Fatal(err)
		}
		st.After = []string{holder}
		if err := lab.repo.Save(st); err != nil {
			t.Fatal(err)
		}
		lab.toReady(st.ID)
		if h := lab.heldNow(st.ID); h == nil || h.Code != workitem.HoldAfter {
			t.Fatalf("held by after: %+v", h)
		}
		lab.l.look(ctx, false)
		if got := lab.asks(st.ID); len(got) != 0 {
			t.Errorf("asked about a story held by after: %v", got)
		}
	})
	t.Run("no touches", func(t *testing.T) {
		lab := newAgentLab(t)
		bare := lab.readyTouching("Bare", []string{}...)
		lab.move(bare, workitem.InProgress)
		held := lab.readyTouching("Held", "flai/cmd/serve.go")
		if h := lab.heldNow(held); h == nil || h.Code != workitem.HoldNoTouches {
			t.Fatalf("held by a story with no touches: %+v", h)
		}
		lab.l.look(ctx, false)
		if got := lab.asks(held); len(got) != 0 {
			t.Errorf("asked about a story held by one with no touches: %v", got)
		}
	})
	t.Run("a holder still ready", func(t *testing.T) {
		lab := newAgentLab(t)
		lab.limit(3)
		lab.hold()
		first := lab.readyTouching("First", "flai/cmd")
		held := lab.readyTouching("Held", "flai/cmd/serve.go")
		lab.l.look(ctx, false)
		waitFor(t, "the first story's agent runs", func() bool { return lab.run(first).live() })
		lab.l.look(ctx, false)
		if lab.run(held) != nil {
			t.Fatalf("started a story held by one whose agent started: %+v", lab.run(held))
		}
		if got := lab.asks(held); len(got) != 0 {
			t.Errorf("asked a story still in ready: %v", got)
		}
		if n := lab.logged("hold not asked about"); n != 0 {
			t.Errorf("warned %d times", n)
		}
		lab.release(first)
		waitFor(t, "it ends", func() bool { return !lab.run(first).live() })
	})
}

// ADR-0134: an ask that cannot be written is warned of with the held and
// holding stories and the error, once until the error changes, and the look
// goes on.
func TestAnAskThatCannotBeWrittenIsWarnedOfOnce(t *testing.T) {
	lab := newAgentLab(t)
	ctx := context.Background()
	lab.limit(3)
	lab.hold()
	holder := lab.readyTouching("Holder", "flai/cmd")
	lab.move(holder, workitem.InProgress)
	held := lab.readyTouching("Held", "flai/cmd/serve.go")
	clear := lab.readyTouching("Clear", "docs/clear")
	// a file where the conversations' folder goes: no conversation can be written
	if err := os.WriteFile(filepath.Join(lab.root, "wip", messages.Folder), nil, 0o644); err != nil {
		t.Fatal(err)
	}
	lab.l.look(ctx, false)
	lab.l.look(ctx, false)
	if n := lab.logged("hold not asked about"); n != 1 || !strings.Contains(lab.logText(), `msg="hold not asked about" component=serve project=t held=`+held+" holder="+holder+" err=") {
		t.Errorf("warned %d times:\n%s", n, lab.logText())
	}
	waitFor(t, "the look went on to the clear story", func() bool { return lab.run(clear).live() })
	lab.release(clear)
	waitFor(t, "it ends", func() bool { return !lab.run(clear).live() })
}

// ADR-0134: a story held on overlap alone by two stories in progress is asked
// about of each, with the paths each holds it on.
func TestAHoldByTwoIsAskedAboutOfEach(t *testing.T) {
	lab := newAgentLab(t)
	ctx := context.Background()
	lab.limit(3)
	a := lab.readyTouching("A", "flai/cmd/a.go")
	lab.move(a, workitem.InProgress)
	b := lab.readyTouching("B", "docs/guide.md")
	lab.move(b, workitem.InProgress)
	held := lab.readyTouching("Held", "flai/cmd", "docs/guide.md")
	lab.l.look(ctx, false)
	if got := lab.asks(held); len(got) != 2 || got[a] != 1 || got[b] != 1 {
		t.Fatalf("asks: %v", got)
	}
	cs, err := messages.For(lab.repo, held)
	if err != nil {
		t.Fatal(err)
	}
	about := map[string][]string{}
	for _, c := range cs {
		about[c.Other(held)] = c.About
	}
	if !slices.Equal(about[a], []string{"flai/cmd/a.go"}) || !slices.Equal(about[b], []string{"docs/guide.md"}) {
		t.Errorf("about: %v", about)
	}
}

// ADR-0134: once the holding story's agent shares the paths, the held story
// is started, and its agent's prompt names the conversation, the holder, the
// paths, and the split; a story started on no share is told of none.
func TestAStoryStartedOnAShareIsToldTheSplit(t *testing.T) {
	lab := newAgentLab(t)
	ctx := context.Background()
	lab.limit(3)
	lab.cfg.Command = nil
	lab.cfg.Harnesses = map[string]harness.Host{harness.ClaudeCode: {Program: lab.stub}}
	claude := &manifest.Agent{Harness: harness.ClaudeCode}
	holder := lab.readyTouching("Holder", "flai/cmd")
	lab.move(holder, workitem.InProgress)
	id := lab.backlog("Held", claude, "flai/cmd/serve.go")
	lab.toReady(id)
	lab.l.look(ctx, false)
	cs, err := messages.For(lab.repo, id)
	if err != nil || len(cs) != 1 {
		t.Fatalf("asked: %v %+v", err, cs)
	}
	if lab.run(id) != nil {
		t.Fatalf("started before the share: %+v", lab.run(id))
	}
	split := "Holder changes the flags; Held adds the serve command below them."
	if _, err := messages.Share(lab.repo, messages.ShareOptions{ID: cs[0].ID, Story: holder, Author: "builder-" + holder, Paths: []string{"flai/cmd/serve.go"}, Split: split, Now: lab.now.Add(time.Minute)}); err != nil {
		t.Fatal(err)
	}
	lab.l.look(ctx, false)
	waitFor(t, "the held story's agent ran", func() bool { r := lab.run(id); return r != nil && !r.live() })
	got, _ := os.ReadFile(filepath.Join(lab.outDir, id+".txt"))
	for _, want := range []string{
		id + " works on a share (ADR-0134)",
		"On " + cs[0].ID + ", " + holder + "'s agent shared `flai/cmd/serve.go` with " + id + `, split so: "` + split + `".`,
		"Keep to the split",
		"Your first inbox lists " + cs[0].ID + " under messages",
		"message_reply (flai message reply " + cs[0].ID + " on the host)",
	} {
		if !strings.Contains(string(got), want) {
			t.Errorf("the prompt lacks %q:\n%s", want, got)
		}
	}

	other := lab.backlog("Unshared", claude, "docs/unshared.md")
	lab.toReady(other)
	lab.l.look(ctx, false)
	waitFor(t, "the other story's agent ran", func() bool { r := lab.run(other); return r != nil && !r.live() })
	if got, _ := os.ReadFile(filepath.Join(lab.outDir, other+".txt")); !strings.Contains(string(got), "You are builder-"+other) || strings.Contains(string(got), "works on a share") {
		t.Errorf("a story started on no share is told of one:\n%s", got)
	}
}

// providerLab is an agentLab whose stories' agents run claude-code, its
// program the stub, with the host's provider gateway, whose key is in
// S0350_GATEWAY_KEY (ADR-0129).
func providerLab(t *testing.T) *agentLab {
	t.Helper()
	lab := newAgentLab(t)
	lab.cfg.Command = nil
	lab.cfg.Harnesses = map[string]harness.Host{harness.ClaudeCode: {Program: lab.stub}}
	lab.cfg.Providers = map[string]manifest.Provider{
		"gateway": {API: manifest.APIAnthropicMessages, BaseURL: "https://gateway.example.com", KeyEnv: "S0350_GATEWAY_KEY"},
		"keyless": {API: manifest.APIAnthropicMessages, BaseURL: "https://keyless.example.com", KeyEnv: "S0350_UNSET_KEY"},
		"empty":   {API: manifest.APIAnthropicMessages, BaseURL: "https://empty.example.com", KeyEnv: "S0350_EMPTY_KEY"},
		"openai":  {API: manifest.APIOpenAIChat, BaseURL: "https://openai.example.com", KeyEnv: "S0350_GATEWAY_KEY"},
	}
	return lab
}

// S-0350, ADR-0129: a story's agent that names a provider is started with
// the provider's base URL, the project's override laid over the host's, and
// the key copied from flai serve's environment into the child's alone: the
// run records the provider and its key's variable, and the key is in no run
// log, record, journal entry, or line flai serve logs.
func TestAStorysProviderKeyReachesItsAgentAlone(t *testing.T) {
	const key = "sk-sentinel-S0350-never-logged"
	t.Setenv("S0350_GATEWAY_KEY", key)
	lab := providerLab(t)
	ctx := context.Background()
	id := lab.readyWith("Through the gateway", &manifest.Agent{Harness: harness.ClaudeCode, Provider: "gateway"})
	lab.l.look(ctx, false)
	waitFor(t, "it ended", func() bool { r := lab.run(id); return r != nil && !r.live() })
	if got, _ := os.ReadFile(filepath.Join(lab.outDir, id+".txt")); !strings.Contains(string(got), "gateway: https://gateway.example.com "+key+"\n") {
		t.Errorf("the agent was not given the provider's URL and key:\n%s", got)
	}
	r := lab.run(id)
	if r.Provider != "gateway" || r.KeyEnv != "S0350_GATEWAY_KEY" || r.PID == 0 || r.Error != "" {
		t.Errorf("run: %+v", r)
	}
	record, _ := os.ReadFile(lab.o.Dir.agents())
	if !strings.Contains(string(record), `"provider": "gateway"`) || !strings.Contains(string(record), `"key_env": "S0350_GATEWAY_KEY"`) {
		t.Errorf("the start record names no provider:\n%s", record)
	}
	runLog, err := os.ReadFile(r.Log)
	if err != nil {
		t.Fatal(err)
	}
	for what, text := range map[string]string{"the run log": string(runLog), "serve/agents.json": string(record), "the journal": fmt.Sprintf("%+v", lab.entries()), "flai serve's log": lab.logText()} {
		if strings.Contains(text, key) {
			t.Errorf("%s holds the key:\n%s", what, text)
		}
	}

	// the project's manifest overrides the host's base URL
	m, _ := os.ReadFile(filepath.Join(lab.root, manifest.File))
	if err := os.WriteFile(filepath.Join(lab.root, manifest.File), append(m, "providers:\n  gateway:\n    base_url: https://project.example.com\n"...), 0o644); err != nil {
		t.Fatal(err)
	}
	over := lab.readyWith("Through the project's gateway", &manifest.Agent{Harness: harness.ClaudeCode, Provider: "gateway"})
	lab.l.look(ctx, false)
	waitFor(t, "it ended", func() bool { r := lab.run(over); return r != nil && !r.live() })
	if got, _ := os.ReadFile(filepath.Join(lab.outDir, over+".txt")); !strings.Contains(string(got), "gateway: https://project.example.com "+key+"\n") {
		t.Errorf("the project's override was not laid over the host's entry:\n%s", got)
	}
}

// S-0350, ADR-0129: a story's agent whose provider the host has no entry for,
// whose key's variable is unset or empty in flai serve's environment, or
// whose api its harness cannot speak, is not started: the run fails as an
// adapter's refusal does, saying what to set, and no process is started.
func TestAStorysAgentIsRefusedAProviderItCannotHave(t *testing.T) {
	t.Setenv("S0350_GATEWAY_KEY", "sk-sentinel")
	t.Setenv("S0350_EMPTY_KEY", "")
	lab := providerLab(t)
	ctx := context.Background()
	for _, c := range []struct {
		provider string
		want     []string
	}{
		{"keyless", []string{"S0350_UNSET_KEY", "unset or empty in flai serve's environment", "flai serve agent provider keyless --key-env <variable>"}},
		{"empty", []string{"S0350_EMPTY_KEY", "flai serve agent provider empty --key-env <variable>"}},
		{"nowhere", []string{`provider "nowhere" has no entry on this host`, "flai serve agent provider nowhere"}},
		{"openai", []string{"claude-code speaks anthropic-messages alone", "flai serve agent provider openai --api anthropic-messages"}},
	} {
		id := lab.readyWith("Through "+c.provider, &manifest.Agent{Harness: harness.ClaudeCode, Provider: c.provider})
		lab.l.look(ctx, false)
		r := lab.run(id)
		if r == nil || r.Outcome != OutcomeFailed || r.PID != 0 || r.Provider != c.provider || !strings.HasPrefix(r.Why, "could not be started: ") {
			t.Fatalf("%s: run %+v", c.provider, r)
		}
		for _, want := range c.want {
			if !strings.Contains(r.Why, want) {
				t.Errorf("%s: %q does not say %q", c.provider, r.Why, want)
			}
		}
		if _, err := os.Stat(filepath.Join(lab.outDir, id+".txt")); err == nil {
			t.Errorf("%s: the agent was started", c.provider)
		}
		if j := lab.entries(); j[len(j)-1].Outcome != "failed" || !strings.Contains(j[len(j)-1].Detail, id+" could not be started") {
			t.Errorf("%s: journal %+v", c.provider, j[len(j)-1])
		}
	}
	if strings.Contains(lab.logText(), "sk-sentinel") {
		t.Errorf("flai serve's log holds a key:\n%s", lab.logText())
	}
}

// S-0350, ADR-0129: a strategic agent's provider is resolved, and refused,
// as a story's agent's is: the planner, through the operator's command,
// records its provider, and is refused while the key's variable is unset.
func TestThePlannersProviderIsResolvedAsAStorysAgentsIs(t *testing.T) {
	lab := planLab(t)
	lab.cfg.Providers = map[string]manifest.Provider{"gateway": {API: manifest.APIOpenAIChat, BaseURL: "https://gateway.example.com", KeyEnv: "S0350_PLANNER_KEY"}}
	agent := &manifest.Agent{Provider: "gateway"}
	start := func(id string) bool {
		lab.l.mu.Lock()
		defer lab.l.mu.Unlock()
		return lab.l.plan(context.Background(), lab.cfg, id, agent, TriggerAsked)
	}
	id := lab.epic.ID
	if start(id) {
		t.Fatal("the planner was started without its provider's key")
	}
	r := lab.planRun(id)
	if r == nil || r.Outcome != OutcomeFailed || r.PID != 0 || r.Provider != "gateway" || !strings.Contains(r.Why, "S0350_PLANNER_KEY") || !strings.Contains(r.Why, "flai serve agent provider gateway --key-env") {
		t.Fatalf("refused: %+v", r)
	}
	if j := lab.entries(); len(j) != 1 || j[0].Outcome != "failed" || !strings.Contains(j[0].Detail, "the planner for "+id+" could not be started") {
		t.Errorf("journal: %+v", j)
	}
	if _, err := os.Stat(filepath.Join(lab.outDir, "plan-"+id+".txt")); err == nil {
		t.Error("the planner was started")
	}

	t.Setenv("S0350_PLANNER_KEY", "sk-sentinel")
	if !start(id) {
		t.Fatalf("not started with the key set: %+v", lab.planRun(id))
	}
	if r := lab.planRun(id); r.Provider != "gateway" || r.KeyEnv != "S0350_PLANNER_KEY" || r.PID == 0 {
		t.Errorf("run: %+v", r)
	}
	lab.given(id)
}

// S-0355, ADR-0131: what a harness cannot do refuses a story's agent unless
// the operator lifts it: roles without the guard unless guard none is set,
// which the start then says; a start without asking unless auto-approve is
// on or deny_protected is set; both refusals at once when both hold, each
// naming what lifts it. A harness that can do both is never refused.
func TestStartAllowedWeighsWhatAHarnessCannotDo(t *testing.T) {
	none := harness.Adapters[harness.Command].Capabilities()
	claude := harness.Adapters[harness.ClaudeCode].Capabilities()
	roles := &manifest.Agent{Roles: map[string]manifest.Role{manifest.RoleVerify: {Model: "sonnet"}, manifest.RoleExplore: {Model: "haiku"}}}
	lift := func(l HarnessLifts) map[string]HarnessLifts { return map[string]HarnessLifts{harness.Command: l} }
	guardRefusal := "harness command cannot run flai guard, so the story's roles (explore, verify) cannot run; flai serve agent harness command --guard none runs them unguarded (ADR-0131)"
	askRefusal := "harness command cannot ask before a protected write, so flai serve does not start a story's agent on it; turn on the project's auto-approve host action (flai serve enable auto-approve), or say your arguments for it deny the protected paths with flai serve agent harness command --deny-protected (ADR-0131)"
	unguarded := "its roles (explore, verify) run without flai guard (guard: none on harness command, ADR-0131)"
	for _, c := range []struct {
		name      string
		harness   string
		caps      harness.Capabilities
		agent     *manifest.Agent
		cfg       AgentConfig
		unguarded string
		refused   string
	}{
		{"roles without the guard", harness.Command, none, roles, AgentConfig{Lifts: lift(HarnessLifts{DenyProtected: true})}, "", guardRefusal},
		{"roles with guard none", harness.Command, none, roles, AgentConfig{Lifts: lift(HarnessLifts{GuardNone: true, DenyProtected: true})}, unguarded, ""},
		{"no roles needs no guard", harness.Command, none, &manifest.Agent{Model: "m"}, AgentConfig{AutoApprove: true}, "", ""},
		{"no ask", harness.Command, none, nil, AgentConfig{}, "", askRefusal},
		{"no ask, auto-approve on", harness.Command, none, nil, AgentConfig{AutoApprove: true}, "", ""},
		{"no ask, deny_protected", harness.Command, none, nil, AgentConfig{Lifts: lift(HarnessLifts{DenyProtected: true})}, "", ""},
		{"another harness's lifts lift nothing", harness.Command, none, nil, AgentConfig{Lifts: map[string]HarnessLifts{harness.ClaudeCode: {DenyProtected: true}}}, "", askRefusal},
		{"neither", harness.Command, none, roles, AgentConfig{}, "", guardRefusal + "; and " + askRefusal},
		{"claude-code with roles", harness.ClaudeCode, claude, roles, AgentConfig{}, "", ""},
	} {
		t.Run(c.name, func(t *testing.T) {
			got, err := startAllowed(c.harness, c.caps, c.agent, c.cfg)
			if c.refused != "" {
				if err == nil || err.Error() != c.refused || got != "" {
					t.Errorf("startAllowed = %q, %v; want refused with\n%s", got, err, c.refused)
				}
				return
			}
			if err != nil || got != c.unguarded {
				t.Errorf("startAllowed = %q, %v; want %q", got, err, c.unguarded)
			}
		})
	}
}

// S-0355, ADR-0131: a story on the command harness, which can neither run
// flai guard nor ask, is refused before its agent starts unless the
// operator lifts what it lacks: the refusal is the run's failure, in the
// journal and in the story's agent status as the dashboard's story page
// reads it, and no process is started. Lifted, it starts, and a start
// whose roles run unguarded says so in the journal and the log.
func TestACommandStartIsRefusedWhatItsHarnessCannotDo(t *testing.T) {
	ctx := context.Background()
	verify := &manifest.Agent{Roles: map[string]manifest.Role{manifest.RoleVerify: {Model: "sonnet"}}}
	refused := func(t *testing.T, lab *agentLab, id string, want ...string) {
		t.Helper()
		r := lab.run(id)
		if r == nil || r.Outcome != OutcomeFailed || r.PID != 0 || r.Harness != harness.Command || !strings.HasPrefix(r.Why, "could not be started: harness command cannot ") {
			t.Fatalf("run %+v", r)
		}
		a := Activity(lab.root, lab.state())[id]
		if a.State != ActivityFailed {
			t.Errorf("the story's agent status: %+v", a)
		}
		for _, w := range want {
			if !strings.Contains(a.Why, w) {
				t.Errorf("the story's agent status %q does not say %q", a.Why, w)
			}
		}
		if j := lab.entries(); len(j) == 0 || j[len(j)-1].Outcome != "failed" || !strings.Contains(j[len(j)-1].Detail, "an agent for "+id+" could not be started: harness command cannot ") {
			t.Errorf("journal: %+v", j)
		}
		if _, err := os.Stat(filepath.Join(lab.outDir, id+".txt")); err == nil {
			t.Error("the agent was started")
		}
	}
	started := func(t *testing.T, lab *agentLab, id string) hostapi.Entry {
		t.Helper()
		waitFor(t, "it ended", func() bool { r := lab.run(id); return r != nil && r.PID != 0 && !r.live() })
		j := lab.entries()
		if len(j) == 0 || j[len(j)-1].Outcome != "done" || !strings.Contains(j[len(j)-1].Detail, "started stub-agent (command) for "+id) {
			t.Fatalf("journal: %+v", j)
		}
		return j[len(j)-1]
	}
	t.Run("roles without the guard are refused", func(t *testing.T) {
		lab := newAgentLab(t)
		id := lab.readyWith("Roles", verify)
		lab.l.look(ctx, false)
		refused(t, lab, id, "harness command cannot run flai guard", "roles (verify)", "flai serve agent harness command --guard none")
	})
	t.Run("roles with guard none start unguarded, and say so", func(t *testing.T) {
		lab := newAgentLab(t)
		lab.cfg.Lifts = map[string]HarnessLifts{harness.Command: {GuardNone: true, DenyProtected: true}}
		id := lab.readyWith("Roles unguarded", verify)
		lab.l.look(ctx, false)
		e := started(t, lab, id)
		said := "its roles (verify) run without flai guard (guard: none on harness command, ADR-0131)"
		if !strings.Contains(e.Detail, said) {
			t.Errorf("the journal does not say the roles run unguarded: %q", e.Detail)
		}
		if !strings.Contains(lab.logText(), `msg="agent started"`) || !strings.Contains(lab.logText(), `roles="`+said+`"`) {
			t.Errorf("the log does not say the roles run unguarded:\n%s", lab.logText())
		}
	})
	t.Run("without asking it is refused", func(t *testing.T) {
		lab := newAgentLab(t)
		lab.cfg.Lifts = nil
		id := lab.ready("No ask")
		lab.l.look(ctx, false)
		refused(t, lab, id, "harness command cannot ask before a protected write", "flai serve enable auto-approve", "flai serve agent harness command --deny-protected")
	})
	t.Run("auto-approve on starts it", func(t *testing.T) {
		lab := newAgentLab(t)
		lab.cfg.Lifts, lab.cfg.AutoApprove = nil, true
		id := lab.ready("Auto-approved")
		lab.l.look(ctx, false)
		if e := started(t, lab, id); strings.Contains(e.Detail, "without flai guard") {
			t.Errorf("a start with no roles says they run unguarded: %q", e.Detail)
		}
	})
	t.Run("deny_protected starts it", func(t *testing.T) {
		lab := newAgentLab(t)
		lab.cfg.Lifts = map[string]HarnessLifts{harness.Command: {DenyProtected: true}}
		id := lab.ready("Denied")
		lab.l.look(ctx, false)
		started(t, lab, id)
	})
}

// S-0355, ADR-0131: a claude-code start, whose harness runs flai guard and
// asks, is unchanged with nothing lifted: a story with roles starts, and
// nothing says they run unguarded.
func TestAClaudeCodeStartIsUnchangedByWhatOtherHarnessesCannotDo(t *testing.T) {
	lab := newAgentLab(t)
	lab.cfg.Command, lab.cfg.Lifts, lab.cfg.AutoApprove = nil, nil, false
	lab.cfg.Flai = "/usr/local/bin/flai"
	lab.cfg.Harnesses = map[string]harness.Host{harness.ClaudeCode: {Program: lab.stub}}
	id := lab.readyWith("Worked by claude", &manifest.Agent{Harness: harness.ClaudeCode, Roles: map[string]manifest.Role{manifest.RoleVerify: {Harness: harness.ClaudeCode}}})
	lab.l.look(context.Background(), false)
	waitFor(t, "it ended", func() bool { r := lab.run(id); return r != nil && !r.live() })
	if r := lab.run(id); r.PID == 0 || r.Error != "" || r.Harness != harness.ClaudeCode {
		t.Fatalf("run: %+v", r)
	}
	j := lab.entries()
	if len(j) != 1 || j[0].Outcome != "done" || strings.Contains(j[0].Detail, "without flai guard") {
		t.Errorf("journal: %+v", j)
	}
	if strings.Contains(lab.logText(), "without flai guard") {
		t.Errorf("the log says the roles run unguarded:\n%s", lab.logText())
	}
}

// sendBack does what flai serve's forge pass does when a required check
// fails at commit on story's pull request (ADR-0167): it opens a thread on the
// story by flai serve and keeps the send back, made at at, in the story's
// pull request status. The story's move is the test's.
func (lab *agentLab) sendBack(story, commit string, at time.Time) *threads.Thread {
	lab.t.Helper()
	th, err := threads.New(lab.repo, threads.NewOptions{Title: "Pull request #12 of " + story + ": the required check ci failed", On: story, Author: forgeAuthor, Text: "The required check ci failed.", Now: at})
	if err != nil {
		lab.t.Fatal(err)
	}
	st, _, err := forge.ReadStatus(lab.repo.CacheDir(), story)
	if err != nil {
		lab.t.Fatal(err)
	}
	st.Story, st.PullRequest, st.HeadSHA = story, 12, commit
	st.SentBack = append(st.SentBack, forge.SendBack{Key: forge.CheckKey("ci", commit), Kind: forge.EventCheck, Name: "ci", Commit: commit, Thread: th.ID, At: at.UTC().Format(time.RFC3339)})
	if err := forge.WriteStatus(lab.repo.CacheDir(), st); err != nil {
		lab.t.Fatal(err)
	}
	return th
}

// backToInProgress moves story from review to in-progress, as flai does when
// it sends the story back.
func (lab *agentLab) backToInProgress(story string) {
	lab.t.Helper()
	st, _ := lab.repo.Get(story)
	if _, err := lab.repo.Transition(st, workitem.InProgress, sendBackBy, "the required check ci failed", lab.now); err != nil {
		lab.t.Fatal(err)
	}
}

// startedCount counts the agents the journal says were started.
func (lab *agentLab) startedCount() int {
	n := 0
	for _, e := range lab.entries() {
		if strings.HasPrefix(e.Detail, "started ") {
			n++
		}
	}
	return n
}

// S-0394, ADR-0167: a story flai sent back from review gets its agent started
// once, in its last session, told the thread, whatever room the in-progress
// limit leaves and without spending an automatic restart; a second look
// starts no second agent, and a later send back starts it again.
func TestAStorySentBackFromReviewGetsItsAgentOnce(t *testing.T) {
	lab := newAgentLab(t)
	ctx := context.Background()
	lab.cfg.Command = nil
	lab.cfg.Harnesses = map[string]harness.Host{harness.ClaudeCode: {Program: lab.stub}}
	lab.limit(1)
	lab.hold()
	id := lab.readyWith("Sent back", &manifest.Agent{Harness: harness.ClaudeCode})
	lab.l.look(ctx, false)
	waitFor(t, "it runs", func() bool { return lab.run(id).live() })
	first := lab.run(id)
	lab.move(id, workitem.InProgress)
	lab.move(id, workitem.Review)
	lab.release(id)
	waitFor(t, "it ends in review", func() bool { r := lab.run(id); return !r.live() && r.Outcome == OutcomeWorked })
	lab.endedAgo(id, 3*time.Minute)
	_ = os.Remove(filepath.Join(lab.outDir, "release-"+id))

	lab.backToInProgress(id)
	th := lab.sendBack(id, "abc1234def5678", time.Now().Add(-2*time.Minute))
	lab.l.look(ctx, false)
	waitFor(t, "it runs for the send back", func() bool { r := lab.run(id); return r.live() && r.SentBack == th.ID })
	again := lab.run(id)
	if again.Agent != first.Agent || again.Session != first.Session || again.AutoRestarts != 0 || again.Answered != "" {
		t.Errorf("the same agent, in its session, spending no restart: %+v, first %+v", again, first)
	}
	out := filepath.Join(lab.outDir, id+".txt")
	waitFor(t, "the stub was told the thread", func() bool { got, _ := os.ReadFile(out); return strings.Contains(string(got), th.ID) })
	got, _ := os.ReadFile(out)
	for _, want := range []string{
		"flai sent " + id + " back from review to in-progress, on thread " + th.ID + ": the required check ci failed on abc1234def56 of pull request #12.",
		"thread_get (flai thread show " + th.ID + " on the host); your inbox lists it, awaiting you",
		"move " + id + " to review again; flai serve pushes the fix",
		"--resume " + first.Session,
	} {
		if !strings.Contains(string(got), want) {
			t.Errorf("the agent was given\n%s\nwithout %q", got, want)
		}
	}
	if j := lab.entries(); len(j) != 2 || !strings.Contains(j[1].Detail, "again for "+id+" as "+first.Agent+" in its session, flai having sent it back from review on "+th.ID) {
		t.Errorf("journal: %+v", j)
	}

	// a second look starts no second agent for the same send back
	lab.l.look(ctx, false)
	if r := lab.run(id); !r.same(again) || lab.startedCount() != 2 {
		t.Fatalf("started twice for one send back: %+v, journal %+v", r, lab.entries())
	}
	// the agent fixes the branch and moves the story to review again
	lab.move(id, workitem.Review)
	lab.release(id)
	waitFor(t, "it ends in review again", func() bool { r := lab.run(id); return !r.live() && r.Outcome == OutcomeWorked })
	lab.l.look(ctx, false)
	if lab.startedCount() != 2 {
		t.Fatalf("started for a story in review: %+v", lab.entries())
	}

	// a later send back starts it again
	lab.endedAgo(id, 90*time.Second)
	lab.backToInProgress(id)
	later := lab.sendBack(id, "0123456789abcdef", time.Now().Add(-time.Minute))
	lab.l.look(ctx, false)
	waitFor(t, "it is started for the later send back", func() bool { return lab.run(id).SentBack == later.ID })
	if r := lab.run(id); r.Session != first.Session || lab.startedCount() != 3 {
		t.Errorf("the later send back: %+v, journal %+v", r, lab.entries())
	}
	waitFor(t, "it ends", func() bool { return !lab.run(id).live() })
}

// S-0394, ADR-0167: a story sent back with no run of its agent here is
// started in a new session, told to take the story up; one whose run failed
// after the send back is started for the send back, in its session, before
// flai serve would restart it on its own, and spends no automatic restart.
func TestAStorySentBackIsStartedAfreshOrBeforeARestart(t *testing.T) {
	ctx := context.Background()
	t.Run("no run, then a run without a session", func(t *testing.T) {
		lab := newAgentLab(t)
		lab.cfg.Command = nil
		lab.cfg.Harnesses = map[string]harness.Host{harness.ClaudeCode: {Program: lab.stub}}
		id := lab.backlog("Begun elsewhere", &manifest.Agent{Harness: harness.ClaudeCode})
		lab.toReady(id)
		lab.move(id, workitem.InProgress)
		th := lab.sendBack(id, "abc1234", time.Now().Add(-time.Minute))
		// worked from a session of its own, which flai serve did not start:
		// nothing is started beside it
		lab.l.look(ctx, false)
		if r := lab.run(id); r != nil {
			t.Fatalf("a story flai serve never ran an agent for got one: %+v", r)
		}
		// a run from before sessions were kept, before the send back
		ago := time.Now().Add(-2 * time.Minute).UTC().Format(time.RFC3339)
		lab.l.dir.updateAgent(lab.root, func(s *AgentState) {
			s.put(&AgentRun{Story: id, Agent: "builder-" + id, Started: ago, Ended: ago, Outcome: OutcomeWorked})
		})
		lab.l.look(ctx, false)
		waitFor(t, "it ends", func() bool { r := lab.run(id); return r != nil && r.SentBack != "" && !r.live() })
		r := lab.run(id)
		if r.SentBack != th.ID || r.Session == "" || r.Agent != "builder-"+id || r.AutoRestarts != 0 {
			t.Fatalf("run: %+v", r)
		}
		got, _ := os.ReadFile(filepath.Join(lab.outDir, id+".txt"))
		for _, want := range []string{
			"because flai sent it back from review to in-progress, on thread " + th.ID + ": the required check ci failed on abc1234 of pull request #12.",
			"This host has no session of its last agent to go on in.",
			"your first inbox lists it, awaiting you",
			"flai stream open " + id,
			"--session-id " + r.Session,
		} {
			if !strings.Contains(string(got), want) {
				t.Errorf("the agent was given\n%s\nwithout %q", got, want)
			}
		}
		if strings.Contains(string(got), "--resume") || strings.Contains(string(got), "story_start") {
			t.Errorf("a new session resumed, or was told to start the story:\n%s", got)
		}
		if j := lab.entries(); len(j) != 1 || !strings.Contains(j[0].Detail, "in a new session, flai having sent it back from review on "+th.ID) {
			t.Errorf("journal: %+v", j)
		}
	})
	t.Run("a failed run from before the send back", func(t *testing.T) {
		lab := newAgentLab(t)
		lab.cfg.AutoRestarts = 2
		id := lab.putEnded("Failed before", AgentRun{Outcome: OutcomeFailed, AutoRestarts: 1})
		lab.endedAgo(id, 2*time.Minute)
		th := lab.sendBack(id, "abc1234", time.Now().Add(-time.Minute))
		lab.l.look(ctx, false)
		waitFor(t, "it is started for the send back", func() bool { return lab.run(id).SentBack == th.ID })
		if r := lab.run(id); r.Session != "earlier" || r.AutoRestarts != 0 || lab.startedCount() != 1 {
			t.Errorf("run %+v, journal %+v", r, lab.entries())
		}
		waitFor(t, "it ends", func() bool { return !lab.run(id).live() })
	})
}

// S-0394, ADR-0167: what a send back does not start an agent for.
func TestASentBackStorysAgentIsNotStartedWhenItShouldNotBe(t *testing.T) {
	ctx := context.Background()
	for _, c := range []struct {
		name string
		run  AgentRun
		ago  time.Duration // how long before now flai sent it back; its run started 2 minutes ago
		set  func(lab *agentLab, id string, th *threads.Thread)
	}{
		{"its thread is resolved", AgentRun{Outcome: OutcomeWorked}, time.Minute, func(lab *agentLab, _ string, th *threads.Thread) {
			if _, err := threads.Resolve(lab.repo, th.ID, "alex", "fixed by hand", time.Now()); err != nil {
				lab.t.Fatal(err)
			}
		}},
		{"its story is blocked", AgentRun{Outcome: OutcomeWorked}, time.Minute, func(lab *agentLab, id string, _ *threads.Thread) {
			st, _ := lab.repo.Get(id)
			if err := workitem.BlockItem(st, "the CI runner", time.Now()); err != nil {
				lab.t.Fatal(err)
			}
			if err := lab.repo.Save(st); err != nil {
				lab.t.Fatal(err)
			}
		}},
		{"its story is still in review", AgentRun{Outcome: OutcomeWorked}, time.Minute, func(lab *agentLab, id string, _ *threads.Thread) {
			lab.move(id, workitem.Review)
		}},
		{"its run started after the send back", AgentRun{Outcome: OutcomeWorked}, 5 * time.Minute, nil},
		{"the operator stopped its run", AgentRun{Outcome: OutcomeStopped, Stopped: "2026-10-06T00:00:00Z"}, time.Minute, nil},
		{"the agent action is off", AgentRun{Outcome: OutcomeWorked}, time.Minute, func(lab *agentLab, _ string, _ *threads.Thread) { lab.cfg.Enabled = false }},
		{"nothing can start it", AgentRun{Outcome: OutcomeWorked}, time.Minute, func(lab *agentLab, _ string, _ *threads.Thread) { lab.cfg.Command = nil }},
	} {
		t.Run(c.name, func(t *testing.T) {
			lab := newAgentLab(t)
			id := lab.putEnded("Sent back", c.run)
			lab.endedAgo(id, 2*time.Minute)
			th := lab.sendBack(id, "abc1234", time.Now().Add(-c.ago))
			if c.set != nil {
				c.set(lab, id, th)
			}
			lab.l.look(ctx, false)
			if r := lab.run(id); r.Session != "earlier" || r.live() {
				t.Errorf("started: %+v", r)
			}
			if j := lab.entries(); len(j) != 0 {
				t.Errorf("journal: %+v", j)
			}
		})
	}
}

// queueBatch is the merge queue batch the lab's queue send backs fail.
const queueBatch = "20261019-220000"

// queueSendBack does what the merge queue's send back does when batch
// queueBatch fails with story at unit (ADR-0176): it opens a thread on the
// story by flai serve and records the story's outcome, failing at at, at a
// head of its branch. No pull request status is written. The story's move
// is the test's.
func (lab *agentLab) queueSendBack(story string, at time.Time) *threads.Thread {
	lab.t.Helper()
	th, err := threads.New(lab.repo, threads.NewOptions{Title: "Merge queue batch " + queueBatch + " fails with " + story + " at unit", On: story, Author: "flai serve",
		Text: "Merge queue batch " + queueBatch + " fails with " + story + " at unit.", Now: at})
	if err != nil {
		lab.t.Fatal(err)
	}
	lab.queueOutcome(story, mergequeue.OutcomeFailing, "sent back to in-progress: merge queue batch "+queueBatch+" fails with it at unit", at)
	return th
}

// queueOutcome records outcome, with detail, at at, as story's last in the
// merge queue's state.
func (lab *agentLab) queueOutcome(story, outcome, detail string, at time.Time) {
	lab.t.Helper()
	if _, err := mergequeue.Update(lab.repo.CacheDir(), func(s *mergequeue.State) error {
		s.RecordOutcome(mergequeue.Outcome{Story: story, Outcome: outcome, At: at.UTC(), Detail: detail, Base: "ba5e0000000000000000", Head: "fa11fa11fa11fa11"})
		return nil
	}); err != nil {
		lab.t.Fatal(err)
	}
}

// queueBackToInProgress moves story from review to in-progress, as flai does
// when the merge queue sends it back.
func (lab *agentLab) queueBackToInProgress(story string) {
	lab.t.Helper()
	st, _ := lab.repo.Get(story)
	if _, err := lab.repo.Transition(st, workitem.InProgress, sendBackBy, "merge queue batch "+queueBatch+" fails with "+story+" at 4e4e4e4e4e4e: unit failed", lab.now); err != nil {
		lab.t.Fatal(err)
	}
}

// S-0400, ADR-0176: a story the merge queue sent back gets its agent started
// once, in its last session, told the thread and that the queue batches it
// again once it is enqueued, with no pull request status at all; a second
// look starts no second agent, nor does one after that agent ended.
func TestAStoryTheMergeQueueSentBackGetsItsAgentOnce(t *testing.T) {
	lab := newAgentLab(t)
	ctx := context.Background()
	lab.cfg.Command = nil
	lab.cfg.Harnesses = map[string]harness.Host{harness.ClaudeCode: {Program: lab.stub}}
	lab.hold()
	id := lab.readyWith("Sent back by the queue", &manifest.Agent{Harness: harness.ClaudeCode})
	lab.l.look(ctx, false)
	waitFor(t, "it runs", func() bool { return lab.run(id).live() })
	first := lab.run(id)
	lab.move(id, workitem.InProgress)
	lab.move(id, workitem.Review)
	lab.release(id)
	waitFor(t, "it ends in review", func() bool { r := lab.run(id); return !r.live() && r.Outcome == OutcomeWorked })
	lab.endedAgo(id, 3*time.Minute)
	_ = os.Remove(filepath.Join(lab.outDir, "release-"+id))

	lab.queueBackToInProgress(id)
	th := lab.queueSendBack(id, time.Now().Add(-2*time.Minute))
	if files, _ := filepath.Glob(forge.StatusPath(lab.repo.CacheDir(), "*")); len(files) != 0 {
		t.Fatalf("pull request statuses %v, want none", files)
	}
	lab.l.look(ctx, false)
	waitFor(t, "it runs for the send back", func() bool { r := lab.run(id); return r.live() && r.SentBack == th.ID })
	again := lab.run(id)
	if again.Agent != first.Agent || again.Session != first.Session || again.AutoRestarts != 0 || again.Answered != "" {
		t.Errorf("the same agent, in its session, spending no restart: %+v, first %+v", again, first)
	}
	out := filepath.Join(lab.outDir, id+".txt")
	waitFor(t, "the stub was told the thread", func() bool { got, _ := os.ReadFile(out); return strings.Contains(string(got), th.ID) })
	got, _ := os.ReadFile(out)
	for _, want := range []string{
		"flai sent " + id + " back from review to in-progress, on thread " + th.ID + ": merge queue batch " + queueBatch + " fails with it at unit.",
		"thread_get (flai thread show " + th.ID + " on the host); your inbox lists it, awaiting you",
		"move " + id + " to review again; flai serve resolves " + th.ID + " once " + id + " is back in review at a new commit, and the merge queue batches it again once it is enqueued.",
		"--resume " + first.Session,
	} {
		if !strings.Contains(string(got), want) {
			t.Errorf("the agent was given\n%s\nwithout %q", got, want)
		}
	}
	if strings.Contains(string(got), "pushes the fix") || strings.Contains(string(got), "sent back to in-progress: ") {
		t.Errorf("the agent was told of a push, or the outcome's prefix:\n%s", got)
	}
	if j := lab.entries(); len(j) != 2 || !strings.Contains(j[1].Detail, "again for "+id+" as "+first.Agent+" in its session, flai having sent it back from review on "+th.ID) {
		t.Errorf("journal: %+v", j)
	}

	// a second look starts no second agent for the same send back
	lab.l.look(ctx, false)
	if r := lab.run(id); !r.same(again) || lab.startedCount() != 2 {
		t.Fatalf("started twice for one send back: %+v, journal %+v", r, lab.entries())
	}
	// nor does one once that agent ended, the story still in progress
	lab.release(id)
	waitFor(t, "it ends", func() bool { return !lab.run(id).live() })
	lab.l.look(ctx, false)
	if lab.startedCount() != 2 {
		t.Fatalf("started again for the send back its agent was started for: %+v", lab.entries())
	}
}

// S-0400, ADR-0176: what a merge queue send back does not start an agent for.
func TestAStoryTheMergeQueueSentBackIsNotStartedWhenItShouldNotBe(t *testing.T) {
	ctx := context.Background()
	for _, c := range []struct {
		name string
		ago  time.Duration // how long before now the queue sent it back; its run started 2 minutes ago
		set  func(lab *agentLab, id string, th *threads.Thread)
	}{
		{"its thread is resolved", time.Minute, func(lab *agentLab, _ string, th *threads.Thread) {
			if _, err := threads.Resolve(lab.repo, th.ID, "alex", "fixed by hand", time.Now()); err != nil {
				lab.t.Fatal(err)
			}
		}},
		{"its thread is not flai serve's", time.Minute, func(lab *agentLab, id string, th *threads.Thread) {
			if _, err := threads.Resolve(lab.repo, th.ID, "alex", "fixed by hand", time.Now()); err != nil {
				lab.t.Fatal(err)
			}
			if _, err := threads.New(lab.repo, threads.NewOptions{Title: th.Title, On: id, Author: "alex", Text: "Is it?", Now: time.Now()}); err != nil {
				lab.t.Fatal(err)
			}
		}},
		{"its outcome is not failing", time.Minute, func(lab *agentLab, id string, _ *threads.Thread) {
			lab.queueOutcome(id, mergequeue.OutcomeLanded, "landed by merge queue batch "+queueBatch+" as alex", time.Now().Add(-30*time.Second))
		}},
		{"its story is blocked", time.Minute, func(lab *agentLab, id string, _ *threads.Thread) {
			st, _ := lab.repo.Get(id)
			if err := workitem.BlockItem(st, "the CI runner", time.Now()); err != nil {
				lab.t.Fatal(err)
			}
			if err := lab.repo.Save(st); err != nil {
				lab.t.Fatal(err)
			}
		}},
		{"its story is in review again", time.Minute, func(lab *agentLab, id string, _ *threads.Thread) {
			lab.move(id, workitem.Review)
		}},
		{"its run started after the send back", 5 * time.Minute, nil},
	} {
		t.Run(c.name, func(t *testing.T) {
			lab := newAgentLab(t)
			id := lab.putEnded("Sent back by the queue", AgentRun{Outcome: OutcomeWorked})
			lab.endedAgo(id, 2*time.Minute)
			th := lab.queueSendBack(id, time.Now().Add(-c.ago))
			if c.set != nil {
				c.set(lab, id, th)
			}
			lab.l.look(ctx, false)
			if r := lab.run(id); r.Session != "earlier" || r.live() || r.SentBack != "" {
				t.Errorf("started: %+v", r)
			}
			if j := lab.entries(); len(j) != 0 {
				t.Errorf("journal: %+v", j)
			}
		})
	}
	// the control: the same story, sent back a minute ago, is started
	lab := newAgentLab(t)
	id := lab.putEnded("Sent back by the queue", AgentRun{Outcome: OutcomeWorked})
	lab.endedAgo(id, 2*time.Minute)
	th := lab.queueSendBack(id, time.Now().Add(-time.Minute))
	lab.l.look(ctx, false)
	waitFor(t, "it is started for the send back", func() bool { return lab.run(id).SentBack == th.ID })
	waitFor(t, "it ends", func() bool { return !lab.run(id).live() })
}
