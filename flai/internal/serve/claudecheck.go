package serve

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"
	"unicode"

	"github.com/bytepunx/system-flow/flai/internal/config"
	"github.com/bytepunx/system-flow/flai/internal/execx"
	"github.com/bytepunx/system-flow/flai/internal/harness"
	"github.com/bytepunx/system-flow/flai/internal/hostapi"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/threads"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// Checking each new Claude Code against permission_prompt (S-0286,
// ADR-0106). A Claude Code release once changed the answer it accepts from
// the permission prompt tool, and the first sign was a story whose agent
// hung. So when flai serve starts, and at each look that may start a
// claude-code agent, it reads the version of the claude it runs; a version
// this host has no record of is checked once, beside the agents, never
// before them: one headless claude -p, with the adapter's own permission
// arguments on the cheapest model, is asked to Write a file in the .claude/
// folder of an in-progress story's worktree in a scratch project whose
// auto-approve is on. The file holding what was asked is a pass; anything
// else is a failure, recorded with what Claude Code said, and told on a
// thread to the operator of every served project whose agents run
// claude-code. A version with a record is not checked again, whatever its
// outcome.

// Outcomes of a check of a Claude Code version.
const (
	ClaudeCheckPassed = "passed"
	ClaudeCheckFailed = "failed"
)

// ClaudeCheck is this host's record of checking one Claude Code version.
type ClaudeCheck struct {
	Checked string `json:"checked"`         // when, RFC 3339
	Program string `json:"program"`         // the claude that was run, symbolic links followed
	Outcome string `json:"outcome"`         // ClaudeCheckPassed or ClaudeCheckFailed
	Error   string `json:"error,omitempty"` // what went wrong and what Claude Code said, on a failure
}

const (
	// claudeCheckModel is the cheapest model Claude Code names.
	claudeCheckModel = "haiku"
	// claudeCheckTimeout bounds the headless run: an answer the prompt
	// cannot give leaves it waiting.
	claudeCheckTimeout = 5 * time.Minute
	// claudeVersionTimeout bounds claude --version.
	claudeVersionTimeout = 30 * time.Second
	// claudeCheckAgent is who the scratch project's MCP server serves.
	claudeCheckAgent = "flai-check"
	// claudeCheckFile is what is written, below the scratch story's worktree.
	claudeCheckFile = ".claude/flai-check.md"
	// claudeCheckAuthor opens the thread a failure is told on.
	claudeCheckAuthor = "flai serve"
	// claudeCheckLines is how many lines of Claude Code's output a failure quotes.
	claudeCheckLines = 20
)

func (d Dir) claudeChecksFile() string { return filepath.Join(string(d), "claude-checks.json") }

// claudeChecks reads the record of each Claude Code version checked, by
// version; none is an empty map.
func (d Dir) claudeChecks() (map[string]ClaudeCheck, error) {
	out := map[string]ClaudeCheck{}
	data, err := os.ReadFile(d.claudeChecksFile())
	if errors.Is(err, os.ErrNotExist) {
		return out, nil
	}
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data, &out); err != nil {
		return nil, fmt.Errorf("%s: %w", d.claudeChecksFile(), err)
	}
	return out, nil
}

// recordClaudeCheck adds or replaces the record of version.
func (d Dir) recordClaudeCheck(version string, c ClaudeCheck) error {
	all, err := d.claudeChecks()
	if err != nil {
		return err
	}
	all[version] = c
	return d.write(d.claudeChecksFile(), all)
}

// claudeChecker runs the check for flai serve, one at a time.
type claudeChecker struct {
	o      Options
	served func() []Entry // the projects flai serve serves now

	mu      sync.Mutex
	running bool
	seen    string // programIdentity of the claude last checked or found checked
	wg      sync.WaitGroup
}

func newClaudeChecker(o Options, served func() []Entry) *claudeChecker {
	return &claudeChecker{o: o, served: served}
}

// consider starts a check in the background when the project at root may
// have flai serve start a claude-code agent and the claude it would run is
// not the one last looked at. It never waits for the check, and starts none
// while one runs.
func (c *claudeChecker) consider(ctx context.Context, root string, cfg AgentConfig) {
	if !cfg.Enabled && !cfg.Plan && !cfg.Orchestrate && !cfg.Analyze {
		return
	}
	host := harness.ClaudeCodeHost(cfg.host(harness.ClaudeCode))
	id, ok := programIdentity(host.Program)
	if !ok {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.running || id == c.seen || !usesClaudeCode(root) {
		return
	}
	c.running, c.seen = true, id
	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		defer func() {
			c.mu.Lock()
			c.running = false
			c.mu.Unlock()
		}()
		c.check(ctx, host, cfg.Flai)
	}()
}

// wait waits for a check that runs.
func (c *claudeChecker) wait() { c.wg.Wait() }

// check reads the version of the claude host runs and, when this host has no
// record of it, checks it, records the outcome, and tells a failure. flai is
// the flai the scratch project's MCP server is; the running one when empty.
func (c *claudeChecker) check(ctx context.Context, host harness.Host, flai string) {
	program, err := exec.LookPath(host.Program)
	if err != nil {
		c.o.Logger.Warn("claude code not found", "component", "serve", "program", host.Program, "err", err.Error())
		return
	}
	if real, err := filepath.EvalSymlinks(program); err == nil {
		program = real
	}
	version, err := claudeVersion(ctx, host.Program)
	if err != nil {
		c.o.Logger.Warn("claude code version unreadable", "component", "serve", "program", program, "err", err.Error())
		return
	}
	checks, err := c.o.Dir.claudeChecks()
	if err != nil {
		c.o.Logger.Warn("claude code checks unreadable", "component", "serve", "version", version, "err", err.Error())
		return
	}
	if _, done := checks[version]; done {
		return
	}
	if flai == "" {
		if flai, err = os.Executable(); err != nil {
			c.o.Logger.Warn("claude code not checked", "component", "serve", "version", version, "err", fmt.Sprintf("this flai's executable is not known: %v", err))
			return
		}
	}
	c.o.Logger.Info("claude code check started", "component", "serve", "version", version, "program", program)
	why, err := c.run(ctx, host, flai, version)
	if ctx.Err() != nil {
		// flai serve is stopping: that is no outcome, and the next one checks it
		c.o.Logger.Info("claude code check stopped", "component", "serve", "version", version)
		return
	}
	if err != nil {
		c.o.Logger.Warn("claude code not checked", "component", "serve", "version", version, "err", err.Error())
		return
	}
	rec := ClaudeCheck{Checked: c.o.Now().UTC().Format(time.RFC3339), Program: program, Outcome: ClaudeCheckPassed}
	if why != "" {
		rec.Outcome, rec.Error = ClaudeCheckFailed, why
	}
	if err := c.o.Dir.recordClaudeCheck(version, rec); err != nil {
		c.o.Logger.Warn("claude code check not recorded", "component", "serve", "version", version, "file", c.o.Dir.claudeChecksFile(), "err", err.Error())
	}
	if rec.Outcome == ClaudeCheckPassed {
		c.o.Logger.Info("claude code check ended", "component", "serve", "version", version, "outcome", rec.Outcome)
		return
	}
	c.o.Logger.Warn("claude code check ended", "component", "serve", "version", version, "outcome", rec.Outcome, "err", rec.Error)
	c.tell(version, program, rec.Error)
}

// claudeVersion is the version claude --version prints, such as 2.1.290 of
// "2.1.290 (Claude Code)".
func claudeVersion(ctx context.Context, program string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, claudeVersionTimeout)
	defer cancel()
	cmd := exec.CommandContext(ctx, program, "--version")
	cmd.Env = agentEnv(os.Environ())
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("%s --version: %w", program, err)
	}
	fields := strings.Fields(string(out))
	if len(fields) == 0 {
		return "", fmt.Errorf("%s --version printed nothing", program)
	}
	return fields[0], nil
}

// run makes the check in a scratch project it removes afterwards. It says
// why the check failed, "" when it passed, and errs only when the check
// could not be made.
func (c *claudeChecker) run(ctx context.Context, host harness.Host, flai, version string) (string, error) {
	s, err := newClaudeScratch(gitOf(c.o), c.o.Now())
	if err != nil {
		return "", fmt.Errorf("the scratch project could not be made: %w", err)
	}
	defer func() { _ = os.RemoveAll(s.dir) }()
	want := claudeCheckContent(version)
	prompt := fmt.Sprintf("This is flai's check that Claude Code's permission prompt tool works. Use the Write tool once to create the file %s with exactly this content, on one line: %s\nUse no other tool, and do not read the file first.", s.target, want)
	argv, err := harness.ClaudeCodeCheck(host, flai, []string{"--config", s.config, "mcp", "--agent", claudeCheckAgent}, claudeCheckModel, prompt)
	if err != nil {
		return "", err
	}
	out, runErr := runClaude(ctx, s.root, argv)
	rel := filepath.Join(".flai-cache", "worktrees", s.story, filepath.FromSlash(claudeCheckFile))
	got, readErr := os.ReadFile(s.target)
	var verdict string
	switch {
	case readErr == nil && strings.TrimSpace(string(got)) == want:
		return "", nil
	case readErr == nil:
		verdict = fmt.Sprintf("%s holds %q, not %q", rel, strings.TrimSpace(string(got)), want)
	default:
		verdict = fmt.Sprintf("%s was not written", rel)
	}
	if runErr != nil {
		verdict += "; claude " + runErr.Error()
	}
	if said := claudeSaid(out); said != "" {
		verdict += "\n" + said
	}
	return verdict, nil
}

// runClaude runs argv in dir, as an agent's environment has it, for at most
// claudeCheckTimeout, and returns what it printed.
func runClaude(ctx context.Context, dir string, argv []string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, claudeCheckTimeout)
	defer cancel()
	cmd := exec.CommandContext(ctx, argv[0], argv[1:]...)
	cmd.Dir = dir
	cmd.Env = append(agentEnv(os.Environ()), "FLAI_AGENT="+claudeCheckAgent, "FLAI_STARTED_BY=flai-serve")
	var out bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &out
	cmd.WaitDelay = 5 * time.Second
	err := cmd.Run()
	if errors.Is(ctx.Err(), context.DeadlineExceeded) {
		err = fmt.Errorf("gave no answer within %s", claudeCheckTimeout)
	} else if err != nil {
		err = fmt.Errorf("ended: %w", err)
	}
	return out.String(), err
}

// claudeCheckContent is what the check asks to be written for version.
func claudeCheckContent(version string) string {
	return fmt.Sprintf("flai checked Claude Code %s through permission_prompt.", version)
}

// claudeSaid is the part of Claude Code's output a failure quotes: the tool
// calls, the tool errors, how the session ended, and every line that is not
// stream-json, such as an error printed before the session began; the last
// claudeCheckLines of them.
func claudeSaid(out string) string {
	var lines []string
	for line := range strings.Lines(out) {
		for _, e := range streamEntries([]byte(line)) {
			switch {
			case e.Kind == StreamTool:
				lines = append(lines, e.Tool+": "+e.Text)
			case e.Kind == StreamResult && e.Error:
				lines = append(lines, "error: "+e.Text)
			case e.Kind == StreamEnd, e.Kind == StreamOutput:
				lines = append(lines, e.Text)
			}
		}
	}
	var kept []string
	for _, l := range lines {
		for part := range strings.Lines(l) {
			if part = quotable(part); part != "" {
				kept = append(kept, part)
			}
		}
	}
	return strings.Join(kept[max(0, len(kept)-claudeCheckLines):], "\n")
}

// quotable is a line with its control characters dropped and no space at
// its end, as a fenced block in a thread takes it.
func quotable(line string) string {
	line = strings.Map(func(r rune) rune {
		if r == '\t' {
			return ' '
		}
		if unicode.IsControl(r) {
			return -1
		}
		return r
	}, line)
	return strings.TrimRightFunc(line, unicode.IsSpace)
}

// claudeScratch is the scratch project a check runs in.
type claudeScratch struct {
	dir    string // the temporary folder that holds it all
	root   string // the project
	config string // flai's configuration for it, auto-approve on
	story  string // the story in progress
	target string // the file to be written, in the story's worktree's .claude/
}

// newClaudeScratch makes a project in a temporary folder, a git repository
// with one story in progress and its worktree's folder, and a flai
// configuration of its own beside it with auto-approve on. The operator's
// configuration is not read or changed.
func newClaudeScratch(git execx.Runner, now time.Time) (s claudeScratch, err error) {
	tmp, err := os.MkdirTemp("", "flai-claude-check-")
	if err != nil {
		return s, err
	}
	defer func() {
		if err != nil {
			_ = os.RemoveAll(tmp)
		}
	}()
	// permission_prompt follows symbolic links to tell a worktree's paths:
	// the folder is named as it is, not through a link such as macOS's /var.
	if s.dir, err = filepath.EvalSymlinks(tmp); err != nil {
		return s, err
	}
	s.root, s.config = filepath.Join(s.dir, "project"), filepath.Join(s.dir, "config.json")
	for _, d := range []string{"wip/kanban/epics", "wip/kanban/stories", "wip/kanban/tasks", "wip/agents"} {
		if err := os.MkdirAll(filepath.Join(s.root, filepath.FromSlash(d)), 0o755); err != nil {
			return s, err
		}
	}
	m := "version: 1\nname: flai-check\nkey: flai-check\nlayout:\n  design: design\n  docs: docs\n  wip: wip\n"
	if err := os.WriteFile(filepath.Join(s.root, manifest.File), []byte(m), 0o644); err != nil {
		return s, err
	}
	if _, err := git.Run(s.root, "git", "init", "-q"); err != nil {
		return s, err
	}
	repo, err := workitem.Open(s.root)
	if err != nil {
		return s, err
	}
	epic, err := repo.Create(workitem.NewOptions{Type: workitem.Epic, Title: "Check Claude Code", Owner: claudeCheckAgent, Now: now})
	if err != nil {
		return s, err
	}
	story, err := repo.Create(workitem.NewOptions{Type: workitem.Story, Title: "Write a file under .claude", Parent: epic.ID, Owner: claudeCheckAgent, Touches: []string{claudeCheckFile}, Now: now,
		Body: "## Goal\n\nClaude Code writes a file under .claude through permission_prompt.\n\n## Acceptance criteria\n\n- [ ] " + claudeCheckFile + " holds what was asked.\n"})
	if err != nil {
		return s, err
	}
	for _, to := range []string{workitem.Ready, workitem.InProgress} {
		if _, err := repo.Transition(story, to, claudeCheckAgent, "", now); err != nil {
			return s, err
		}
	}
	s.story = story.ID
	worktree := repo.WorktreePath(story.ID)
	if err := os.MkdirAll(worktree, 0o755); err != nil {
		return s, err
	}
	s.target = filepath.Join(worktree, filepath.FromSlash(claudeCheckFile))
	// every project of this configuration's is the scratch one
	cfg := config.Default().WithAction(hostapi.ActionAutoApprove, config.AllProjects, true)
	return s, config.Save(s.config, cfg)
}

// tell opens a thread to the operator on the manifest of every served
// project whose agents run claude-code, quoting why the check failed.
func (c *claudeChecker) tell(version, program, why string) {
	fence := "```"
	for strings.Contains(why, fence) {
		fence += "`"
	}
	text := fmt.Sprintf("flai serve checked Claude Code %s (%s) against permission_prompt (ADR-0106), and it failed. A headless claude -p in a scratch project, with auto-approve on and the permission arguments a story's agent gets, was asked to Write %s in an in-progress story's worktree. What went wrong, and what Claude Code said:\n\n%stext\n%s\n%s\n\n"+
		"Until this is fixed, writes to the paths Claude Code protects (.claude/, .mcp.json, and the rest ADR-0106 lists) will not go through permission_prompt: an agent's write to one is refused. The agents still run, and a story that writes no protected path is not held up. flai does not check %s again; to check it again, remove its entry from %s.",
		version, program, claudeCheckFile, fence, why, fence, version, c.o.Dir.claudeChecksFile())
	for _, e := range c.served() {
		if !usesClaudeCode(e.Root) {
			continue
		}
		repo, err := workitem.Open(e.Root)
		if err != nil {
			c.o.Logger.Warn("claude code check not told", "component", "serve", "project", e.Key, "version", version, "err", err.Error())
			continue
		}
		th, err := threads.New(repo, threads.NewOptions{Title: fmt.Sprintf("Claude Code %s fails flai's permission_prompt check", version), On: manifest.File, Author: claudeCheckAuthor, Text: text, Now: c.o.Now()})
		if err != nil {
			c.o.Logger.Warn("claude code check not told", "component", "serve", "project", e.Key, "version", version, "err", err.Error())
			continue
		}
		c.o.Logger.Info("claude code check told", "component", "serve", "project", e.Key, "version", version, "thread", th.ID)
	}
}

// usesClaudeCode says whether the project's manifest has its agents, a
// story's or a strategic one's, run claude-code.
func usesClaudeCode(root string) bool {
	m, err := manifest.Load(filepath.Join(root, manifest.File))
	if err != nil {
		return false
	}
	for _, a := range []*manifest.Agent{m.Agent, m.PlanningAgent(), m.AnalysisAgent(), m.OrchestrationAgent()} {
		if a != nil && a.Harness == harness.ClaudeCode {
			return true
		}
	}
	return false
}

// programIdentity names the file program is, symbolic links followed, with
// its size and time: a Claude Code update changes it, so a look need not run
// claude --version to see that nothing changed.
func programIdentity(program string) (string, bool) {
	path, err := exec.LookPath(program)
	if err != nil {
		return "", false
	}
	if path, err = filepath.EvalSymlinks(path); err != nil {
		return "", false
	}
	fi, err := os.Stat(path)
	if err != nil {
		return "", false
	}
	return fmt.Sprintf("%s %d %d", path, fi.Size(), fi.ModTime().UnixNano()), true
}
