package pilot

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/conventions"
	"github.com/bytepunx/system-flow/flai/internal/harness"
	"github.com/bytepunx/system-flow/flai/internal/roles"
)

// fakeClaude is a claude for pilot to start: it records its arguments, one
// after another each ended by a NUL, and the variables pilot gives it, but
// of the token only whether it is set; writes a line to its standard error;
// and then, by PILOT_TEST_MODE, prints the file out and exits with
// PILOT_TEST_EXIT, traps SIGTERM and exits 143 once it comes, or is killed by
// the signal pilot forwards. It writes ready once it waits for a signal.
const fakeClaude = `#!/bin/sh
d="$PILOT_TEST_DIR"
: > "$d/argv"
for a in "$@"; do printf '%s\000' "$a" >> "$d/argv"; done
{
  printf 'CLAUDE_CONFIG_DIR=%s\n' "$CLAUDE_CONFIG_DIR"
  printf 'FLAI_AGENT=%s\n' "$FLAI_AGENT"
  printf 'FLAI_SESSION=%s\n' "$FLAI_SESSION"
  printf 'FLAI_STORY=%s\n' "$FLAI_STORY"
  printf 'FLAI_STORY is %s\n' "${FLAI_STORY+set}"
  printf 'FLAI_ROLE=%s\n' "$FLAI_ROLE"
  printf 'FLAI_ITEM=%s\n' "$FLAI_ITEM"
  if [ -n "$FLAI_MCP_TOKEN" ]; then echo 'token=set'; else echo 'token=unset'; fi
} > "$d/env"
echo 'the harness on its standard error' >&2
case "$PILOT_TEST_MODE" in
trap)
  trap 'kill "$pid" 2>/dev/null; echo trapped > "$d/trapped"; exit 143' TERM
  sleep 30 > /dev/null 2>&1 &
  pid=$!
  echo ready > "$d/ready"
  wait "$pid"
  exit 1 ;;
die)
  echo ready > "$d/ready"
  exec sleep 30 > /dev/null ;;
esac
cat "$d/out"
exit "${PILOT_TEST_EXIT:-0}"
`

// printed is what the fake claude prints: a line, a line longer than any
// line buffer's default, and a last line without a newline.
var printed = "{\"type\":\"system\"}\n" + strings.Repeat("x", 200<<10) + "\n" + `{"type":"result"} without a newline`

// lockedBuffer is a bytes.Buffer that pilot and the harness's standard
// error copier may write to at once.
type lockedBuffer struct {
	mu  sync.Mutex
	buf bytes.Buffer
}

func (b *lockedBuffer) Write(p []byte) (int, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.buf.Write(p)
}

func (b *lockedBuffer) String() string {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.buf.String()
}

// fake puts the fake claude first on PATH in a new folder, with printed as
// its output, and gives the folder.
func fake(t *testing.T) string {
	t.Helper()
	if runtime.GOOS == "windows" {
		t.Skip("the fake harness is a POSIX shell script")
	}
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "claude"), []byte(fakeClaude), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "out"), []byte(printed), 0o600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", dir+string(filepath.ListSeparator)+os.Getenv("PATH"))
	return dir
}

// environ is pilot's environment for the fake claude in dir: this test's,
// with the MCP server's variables and extra.
func environ(dir string, extra ...string) []string {
	return append(os.Environ(), append([]string{
		"FLAI_MCP_URL=" + urlValue,
		"FLAI_MCP_TOKEN=" + tokenValue,
		"FLAI_MCP_HOST=" + hostValue,
		"PILOT_TEST_DIR=" + dir,
	}, extra...)...)
}

// result is what one run of pilot gave.
type result struct {
	code           int
	stdout, stderr string
}

// run runs Main on req, given as JSON, failing t when it has not returned
// within the time limit, or when it printed the token.
func run(t *testing.T, ctx context.Context, req any, env []string, signals <-chan os.Signal) result {
	t.Helper()
	data, err := json.Marshal(req)
	if err != nil {
		t.Fatal(err)
	}
	var stdout, stderr lockedBuffer
	codes := make(chan int, 1)
	go func() {
		codes <- Main(ctx, bytes.NewReader(data), &stdout, &stderr, env, signals)
	}()
	select {
	case code := <-codes:
		r := result{code: code, stdout: stdout.String(), stderr: stderr.String()}
		if strings.Contains(r.stdout, tokenValue) || strings.Contains(r.stderr, tokenValue) {
			t.Errorf("pilot printed the token's value; stderr: %s", r.stderr)
		}
		return r
	case <-time.After(30 * time.Second):
		t.Fatalf("pilot has not returned after 30s; stderr: %s", stderr.String())
		return result{}
	}
}

// pilotLines are stderr's lines that are pilot's own.
func pilotLines(stderr string) []string {
	var out []string
	for _, l := range strings.Split(strings.TrimSuffix(stderr, "\n"), "\n") {
		if strings.HasPrefix(l, "pilot: ") {
			out = append(out, l)
		}
	}
	return out
}

// argv is what the fake claude in dir was last started with.
func argv(t *testing.T, dir string) []string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(dir, "argv"))
	if err != nil {
		t.Fatal(err)
	}
	return strings.Split(strings.TrimSuffix(string(data), "\x00"), "\x00")
}

// follows is whether args has flag followed by value.
func follows(args []string, flag, value string) bool {
	i := slices.Index(args, flag)
	return i >= 0 && i+1 < len(args) && args[i+1] == value
}

// waitFor waits for path to exist, failing t, without stopping it, after
// the time limit: it runs beside the test, not in it.
func waitFor(t *testing.T, path string) {
	deadline := time.Now().Add(10 * time.Second)
	for {
		if _, err := os.Stat(path); err == nil {
			return
		}
		if time.Now().After(deadline) {
			t.Errorf("%s was not written within 10s", path)
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
}

func TestRunRelaysOutputAndExitCode(t *testing.T) {
	for _, code := range []int{0, 7} {
		dir := fake(t)
		r := run(t, t.Context(), valid(t.TempDir()), environ(dir, "PILOT_TEST_EXIT="+strconv.Itoa(code)), nil)
		if r.code != code {
			t.Errorf("pilot exited %d, want the harness's %d; stderr: %s", r.code, code, r.stderr)
		}
		if r.stdout != printed {
			t.Errorf("pilot relayed %d bytes, want the %d the harness printed, unchanged", len(r.stdout), len(printed))
		}
		if !strings.Contains(r.stderr, "the harness on its standard error\n") {
			t.Errorf("the harness's standard error did not reach pilot's: %s", r.stderr)
		}
		lines := pilotLines(r.stderr)
		if len(lines) != 2 || !strings.Contains(lines[0], "starting claude-code") || !strings.HasSuffix(lines[1], "exited with code "+strconv.Itoa(code)) {
			t.Errorf("pilot's own lines are %q, want a start line and an end line with the code", lines)
		}
	}
}

func TestRunStartsNewSessionThenResumesIt(t *testing.T) {
	dir := fake(t)
	cache := t.TempDir()
	req := valid(cache)
	mcp := filepath.Join(cache, "mcp.json")

	first := run(t, t.Context(), req, environ(dir), nil)
	args := argv(t, dir)
	if first.code != 0 || !follows(args, "--session-id", req.Session) || slices.Contains(args, "--resume") {
		t.Errorf("the first run exited %d with --session-id %s and no --resume? argv %q", first.code, req.Session, args)
	}
	if !follows(args, "--mcp-config", mcp) || !slices.Contains(args, "--strict-mcp-config") {
		t.Errorf("argv does not name %s in --mcp-config with --strict-mcp-config: %q", mcp, args)
	}
	if !strings.Contains(first.stderr, "in a new session") {
		t.Errorf("the first run's start line does not say the session is new: %s", first.stderr)
	}
	if strings.Contains(first.stderr, req.Session) {
		t.Errorf("pilot's lines name the session: %s", first.stderr)
	}

	second := run(t, t.Context(), req, environ(dir), nil)
	args = argv(t, dir)
	if second.code != 0 || !follows(args, "--resume", req.Session) || slices.Contains(args, "--session-id") {
		t.Errorf("the second run exited %d with --resume %s and no --session-id? argv %q", second.code, req.Session, args)
	}
	if !strings.Contains(second.stderr, "resuming its session") {
		t.Errorf("the second run's start line does not say it resumes: %s", second.stderr)
	}
}

func TestRunWritesTheMCPFileWithoutTheToken(t *testing.T) {
	dir := fake(t)
	cache := t.TempDir()
	mcp := filepath.Join(cache, "mcp.json")
	// A file left from before, open to all, is rewritten and closed.
	if err := os.WriteFile(mcp, []byte("{}"), 0o644); err != nil {
		t.Fatal(err)
	}
	if r := run(t, t.Context(), valid(cache), environ(dir), nil); r.code != 0 {
		t.Fatalf("pilot exited %d; stderr: %s", r.code, r.stderr)
	}
	info, err := os.Stat(mcp)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0o600 {
		t.Errorf("the MCP file's mode is %v, want 0600", info.Mode().Perm())
	}
	data, err := os.ReadFile(mcp)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(data), tokenValue) || !strings.Contains(string(data), "${FLAI_MCP_TOKEN}") {
		t.Errorf("the MCP file holds the token's value, or does not name its variable: %s", data)
	}
	// ADR-0258: the file names the per-server timeout, so a slow task_done or test outlives Claude Code's 60 seconds
	var file struct {
		MCPServers map[string]struct {
			Timeout int64 `json:"timeout"`
		} `json:"mcpServers"`
	}
	if err := json.Unmarshal(data, &file); err != nil {
		t.Fatalf("the MCP file is not JSON: %v: %s", err, data)
	}
	if got := file.MCPServers["flai"].Timeout; got != 1800000 {
		t.Errorf("the MCP file's timeout for flai is %d ms, want 1800000: %s", got, data)
	}
}

func TestRunGivesTheHarnessItsEnvironment(t *testing.T) {
	dir := fake(t)
	cache := t.TempDir()
	req := valid(cache)
	if r := run(t, t.Context(), req, environ(dir), nil); r.code != 0 {
		t.Fatalf("pilot exited %d; stderr: %s", r.code, r.stderr)
	}
	data, err := os.ReadFile(filepath.Join(dir, "env"))
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		"CLAUDE_CONFIG_DIR=" + filepath.Join(cache, "claude"),
		"FLAI_AGENT=" + req.Name,
		"FLAI_SESSION=" + req.Session,
		"FLAI_STORY=" + req.Story,
		"token=set",
	} {
		if !strings.Contains(string(data), want+"\n") {
			t.Errorf("the harness's environment lacks %s:\n%s", want, data)
		}
	}
}

// I-0224: the start request scripts/agent-container.sh wrote for a trial of
// the image, a story's agent under the story's own name, started claude with
// the story's prompt, which begins S-0507 with story_start, and FLAI_STORY;
// the same request marked probe, under a name of its own, starts it with the
// probe's prompt and no FLAI_STORY, even one the container inherited.
func TestRunStartsAProbeWithoutTheStory(t *testing.T) {
	for _, c := range []struct {
		name   string
		probe  bool
		agent  string
		prompt string // what claude's prompt holds
		not    string // what it does not
		story  string // the child's FLAI_STORY lines
	}{
		{"the story's agent", false, "agent-S-0507", "story_start with S-0507", "trial of the agent image", "FLAI_STORY=S-0507\nFLAI_STORY is set\n"},
		{"a probe", true, "probe-S-0507", "trial of the agent image", "story_start with S-0507", "FLAI_STORY=\nFLAI_STORY is \n"},
	} {
		t.Run(c.name, func(t *testing.T) {
			dir := fake(t)
			req := valid(t.TempDir())
			req.Name, req.Probe = c.agent, c.probe
			if r := run(t, t.Context(), req, environ(dir, "FLAI_STORY=S-0507"), nil); r.code != 0 {
				t.Fatalf("pilot exited %d; stderr: %s", r.code, r.stderr)
			}
			args := argv(t, dir)
			i := slices.Index(args, "-p")
			if i < 0 || i+1 >= len(args) {
				t.Fatalf("argv has no -p prompt: %q", args)
			}
			if prompt := args[i+1]; !strings.Contains(prompt, c.prompt) || strings.Contains(prompt, c.not) {
				t.Errorf("claude's prompt does not hold %q, or holds %q:\n%s", c.prompt, c.not, prompt)
			}
			env, err := os.ReadFile(filepath.Join(dir, "env"))
			if err != nil {
				t.Fatal(err)
			}
			for _, want := range []string{"FLAI_AGENT=" + c.agent + "\n", "FLAI_SESSION=" + req.Session + "\n", c.story} {
				if !strings.Contains(string(env), want) {
					t.Errorf("the harness's environment lacks %q:\n%s", want, env)
				}
			}
		})
	}
}

// checkout is a project's main checkout as a strategic agent's container
// mounts it: its manifest, and the definitions of the four built-in roles'
// sessions, each with a prompt of its own that a role file's body replaces.
func checkout(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	files := map[string]string{"system-flow.yaml": "version: 1\nname: flow\nlayout:\n  design: design\n  docs: docs\n  wip: wip\n"}
	for _, name := range []string{"planner", "orchestrator", "analyzer", "curator"} {
		files[filepath.Join(".claude", "agents", name+".md")] = "---\nname: " + name + "\ndescription: The " + name + ".\n---\nThe definition's own prompt.\n"
	}
	for rel, data := range files {
		p := filepath.Join(root, rel)
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte(data), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	return root
}

// agentPrompt is the prompt of the definition name in the --agents JSON of
// args.
func agentPrompt(t *testing.T, args []string, name string) string {
	t.Helper()
	i := slices.Index(args, "--agents")
	if i < 0 || i+1 >= len(args) {
		t.Fatalf("argv has no --agents: %q", args)
	}
	var agents map[string]map[string]any
	if err := json.Unmarshal([]byte(args[i+1]), &agents); err != nil {
		t.Fatalf("--agents %s: %v", args[i+1], err)
	}
	prompt, _ := agents[name]["prompt"].(string)
	return prompt
}

// S-0509, ADR-0247 §3: a built-in role's start request reaches the inner
// harness's start as its definition, with --agent and FLAI_ROLE, and with
// the role file pilot reads from the mounted checkout, the project's own, or
// the baseline flai carries when the project has none, as its prompt.
func TestRunStartsAStrategicAgentAsItsDefinition(t *testing.T) {
	baseline, ok := roles.BaselineFor("design", conventions.RolePlan)
	if !ok {
		t.Fatal("flai carries no baseline planner")
	}
	firstLine := strings.TrimSpace(strings.SplitN(strings.TrimSpace(baseline.Baseline), "\n", 2)[0])
	for _, c := range []struct {
		name     string
		roleFile string // the project's planner role file; none when empty
		want     string // what the planner's prompt holds
	}{
		{"the project's role file", baselineFile(t) + "\nPlan in this project's own way.\n", "Plan in this project's own way."},
		{"the baseline", "", firstLine},
	} {
		t.Run(c.name, func(t *testing.T) {
			dir := fake(t)
			root := checkout(t)
			if c.roleFile != "" {
				p := filepath.Join(root, "design", "workflows", "roles", "planner.md")
				if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(p, []byte(c.roleFile), 0o600); err != nil {
					t.Fatal(err)
				}
			}
			req := valid(t.TempDir())
			req.Story, req.Root, req.Role, req.Item = "", root, conventions.RolePlan, "E-0029"
			if r := run(t, t.Context(), req, environ(dir), nil); r.code != 0 {
				t.Fatalf("pilot exited %d; stderr: %s", r.code, r.stderr)
			}
			args := argv(t, dir)
			if !follows(args, "--agent", "planner") {
				t.Errorf("the planner's argv has no --agent planner: %q", args)
			}
			if prompt := agentPrompt(t, args, "planner"); !strings.Contains(prompt, c.want) || strings.Contains(prompt, "The definition's own prompt.") {
				t.Errorf("the planner's prompt is not the role file's body with %q:\n%s", c.want, prompt)
			}
			env, err := os.ReadFile(filepath.Join(dir, "env"))
			if err != nil {
				t.Fatal(err)
			}
			for _, want := range []string{"FLAI_ROLE=plan\n", "FLAI_ITEM=E-0029\n"} {
				if !strings.Contains(string(env), want) {
					t.Errorf("the planner's environment lacks %q:\n%s", want, env)
				}
			}
		})
	}
}

// baselineFile is the baseline planner's role file, as a project made from
// the template holds it.
func baselineFile(t *testing.T) string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join("..", "roles", "baseline", "planner.md"))
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}

// S-0509: each built-in role's start request, as the container adapter
// writes it, reaches the inner harness as the role's definition with
// FLAI_ROLE.
func TestRunStartsEachBuiltinRole(t *testing.T) {
	for _, c := range []struct {
		role, agent string
		change      func(*harness.Request)
	}{
		{conventions.RolePlan, "planner", func(r *harness.Request) { r.Item = "S-0509" }},
		{conventions.RoleOrchestrate, "orchestrator", func(r *harness.Request) { r.Reasons = []string{"moved S-0509 to review: a title"} }},
		{conventions.RoleAnalyze, "analyzer", func(r *harness.Request) { r.Focus = "risk" }},
		{conventions.RoleCurate, "curator", func(r *harness.Request) { r.Since = "2026-10-01" }},
	} {
		t.Run(c.role, func(t *testing.T) {
			dir := fake(t)
			req := valid(t.TempDir())
			req.Story, req.Root, req.Role = "", checkout(t), c.role
			c.change(&req)
			if r := run(t, t.Context(), req, environ(dir), nil); r.code != 0 {
				t.Fatalf("pilot exited %d; stderr: %s", r.code, r.stderr)
			}
			args := argv(t, dir)
			if !follows(args, "--agent", c.agent) {
				t.Errorf("argv has no --agent %s: %q", c.agent, args)
			}
			if prompt := agentPrompt(t, args, c.agent); prompt == "" || strings.Contains(prompt, "The definition's own prompt.") {
				t.Errorf("the %s's prompt is not its baseline role file's body:\n%s", c.agent, prompt)
			}
			env, err := os.ReadFile(filepath.Join(dir, "env"))
			if err != nil {
				t.Fatal(err)
			}
			if !strings.Contains(string(env), "FLAI_ROLE="+c.role+"\n") {
				t.Errorf("the %s's environment lacks FLAI_ROLE=%s:\n%s", c.agent, c.role, env)
			}
		})
	}
}

// S-0509: a built-in role's start in a folder that holds no project's
// manifest is refused, as flai serve refuses it on the host, before the inner
// harness starts.
func TestRunRefusesARoleWithoutTheProjectsManifest(t *testing.T) {
	dir := fake(t)
	root := checkout(t)
	if err := os.Remove(filepath.Join(root, "system-flow.yaml")); err != nil {
		t.Fatal(err)
	}
	req := valid(t.TempDir())
	req.Story, req.Root, req.Role, req.Item = "", root, conventions.RolePlan, "E-0029"
	r := run(t, t.Context(), req, environ(dir), nil)
	if r.code != ExitUsage || !strings.Contains(r.stderr, "mount the project's main checkout at the start request's root") {
		t.Errorf("pilot exited %d, want %d with what to do; stderr: %s", r.code, ExitUsage, r.stderr)
	}
	if _, err := os.Stat(filepath.Join(dir, "argv")); err == nil {
		t.Error("the harness was started")
	}
}

func TestRunForwardsASignalToTheHarness(t *testing.T) {
	dir := fake(t)
	signals := make(chan os.Signal, 1)
	go func() {
		waitFor(t, filepath.Join(dir, "ready"))
		signals <- syscall.SIGTERM
	}()
	r := run(t, t.Context(), valid(t.TempDir()), environ(dir, "PILOT_TEST_MODE=trap"), signals)
	if _, err := os.Stat(filepath.Join(dir, "trapped")); err != nil {
		t.Errorf("the harness did not get the SIGTERM sent to pilot: %v", err)
	}
	if r.code != 143 {
		t.Errorf("pilot exited %d, want the harness's 143; stderr: %s", r.code, r.stderr)
	}
}

func TestRunEndsTheHarnessWhenItsContextIsDone(t *testing.T) {
	dir := fake(t)
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()
	go func() {
		waitFor(t, filepath.Join(dir, "ready"))
		cancel()
	}()
	r := run(t, ctx, valid(t.TempDir()), environ(dir, "PILOT_TEST_MODE=die"), nil)
	if want := 128 + int(syscall.SIGTERM); r.code != want {
		t.Errorf("pilot exited %d for a harness SIGTERM killed, want %d; stderr: %s", r.code, want, r.stderr)
	}
}

func TestMainRefusesInOneLine(t *testing.T) {
	cases := []struct {
		name  string
		req   func(cache string) any
		empty bool // PATH holds no claude
	}{
		{"empty object", func(string) any { return map[string]any{} }, false},
		{"unsafe session", func(cache string) any { r := valid(cache); r.Session = "../S-0507"; return r }, false},
		{"no program", func(cache string) any { return valid(cache) }, true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			dir := fake(t)
			if c.empty {
				t.Setenv("PATH", t.TempDir())
			}
			r := run(t, t.Context(), c.req(t.TempDir()), environ(dir), nil)
			if r.code != ExitUsage {
				t.Errorf("pilot exited %d, want %d", r.code, ExitUsage)
			}
			if lines := strings.Split(strings.TrimSuffix(r.stderr, "\n"), "\n"); len(lines) != 1 || !strings.HasPrefix(lines[0], "pilot: ") {
				t.Errorf("pilot's standard error is not one pilot: line: %q", r.stderr)
			}
			if r.stdout != "" {
				t.Errorf("pilot wrote %q on standard output", r.stdout)
			}
			if _, err := os.Stat(filepath.Join(dir, "argv")); err == nil {
				t.Error("the harness was started")
			}
		})
	}
}

// triagerFile is a role the project defines in a role file of its own, the
// triager: it runs in role triage on a story.
const triagerFile = `---
title: Triager
updated: 2026-11-11
status: active
name: triager
role: triage
definition:
  description: Triages a story for flai serve.
  tools: [Read, Grep, mcp__flai__item_get]
  model: inherit
action: triage
trigger:
  - asked
item: [story]
topics: [planning]
works_story: false
charge: none
tools: []
commands: []
files: []
permissions: []
---

You are the triager: read the story and say where it belongs.
`

// withRoleFile is checkout with data as its role file name.md under the
// design folder's workflows/roles.
func withRoleFile(t *testing.T, name, data string) string {
	t.Helper()
	root := checkout(t)
	p := filepath.Join(root, "design", "workflows", "roles", name+".md")
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(p, []byte(data), 0o600); err != nil {
		t.Fatal(err)
	}
	return root
}

// definedStart is the start request the container adapter writes for the
// triager on S-0042 with two inputs, in the checkout root: its name in
// role_name, never its file.
func definedStart(cache, root string) harness.Request {
	req := valid(cache)
	req.Story, req.Root, req.Role, req.RoleName, req.Item, req.Name = "", root, "triage", "triager", "S-0042", harness.RoleAgent("triager", "S-0042")
	req.Inputs = []harness.Input{{Name: "depth", Value: "thorough"}, {Name: "scope", Value: "touches"}}
	return req
}

// I-0258, S-0705: a role the project defines, named by role_name in the
// start request, starts the inner harness as that role, from its file in
// the mounted checkout, as flai serve starts it on the host: with --agent
// and its definition, FLAI_AGENT, FLAI_ROLE, FLAI_ITEM, and a start prompt
// that holds the run's inputs.
func TestRunStartsADefinedRoleFromTheFileItsNameNames(t *testing.T) {
	dir := fake(t)
	req := definedStart(t.TempDir(), withRoleFile(t, "triager", triagerFile))
	if r := run(t, t.Context(), req, environ(dir), nil); r.code != 0 {
		t.Fatalf("pilot exited %d; stderr: %s", r.code, r.stderr)
	}
	args := argv(t, dir)
	if !follows(args, "--agent", "triager") {
		t.Errorf("the triager's argv has no --agent triager: %q", args)
	}
	if prompt := agentPrompt(t, args, "triager"); !strings.Contains(prompt, "You are the triager") {
		t.Errorf("the triager's definition's prompt is not its role file's body:\n%s", prompt)
	}
	i := slices.Index(args, "-p")
	if i < 0 || i+1 >= len(args) {
		t.Fatalf("argv has no -p prompt: %q", args)
	}
	for _, want := range []string{"- depth: thorough", "- scope: touches", "S-0042"} {
		if !strings.Contains(args[i+1], want) {
			t.Errorf("the triager's start prompt lacks %q:\n%s", want, args[i+1])
		}
	}
	env, err := os.ReadFile(filepath.Join(dir, "env"))
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"FLAI_AGENT=triager-S-0042\n", "FLAI_ROLE=triage\n", "FLAI_ITEM=S-0042\n"} {
		if !strings.Contains(string(env), want) {
			t.Errorf("the triager's environment lacks %q:\n%s", want, env)
		}
	}
}

// S-0705: pilot refuses a defined role's start, before the inner harness
// starts, when the checkout holds no file of its role_name, when its file
// does not parse or names another role, and when its file's verb is not the
// request's Role, each in one line that says what to do.
func TestRunRefusesADefinedRoleItCannotLoad(t *testing.T) {
	for _, c := range []struct {
		name   string
		file   string // the role file triager.md; none when empty
		change func(*harness.Request)
		says   string
	}{
		{"no file", "", func(*harness.Request) {}, "names the role file design/workflows/roles/triager.md, which is not in the checkout mounted at"},
		{"a file that does not parse", "---\nname: [\n---\n", func(*harness.Request) {}, "the role file design/workflows/roles/triager.md that the start request's role_name names does not load (roles.file: "},
		{"a file of another role", strings.Replace(triagerFile, "name: triager", "name: sorter", 1), func(*harness.Request) {}, `names its role "sorter", not the start request's role_name triager`},
		{"another verb", triagerFile, func(r *harness.Request) { r.Role = "review" }, `the role triager runs in role triage, as its file design/workflows/roles/triager.md says, and the start request's Role is "review"; set Role to triage`},
	} {
		t.Run(c.name, func(t *testing.T) {
			dir := fake(t)
			root := checkout(t)
			if c.file != "" {
				root = withRoleFile(t, "triager", c.file)
			}
			req := definedStart(t.TempDir(), root)
			c.change(&req)
			r := run(t, t.Context(), req, environ(dir), nil)
			if r.code != ExitUsage || !strings.Contains(r.stderr, c.says) {
				t.Errorf("pilot exited %d, want %d saying %q; stderr: %s", r.code, ExitUsage, c.says, r.stderr)
			}
			if lines := strings.Split(strings.TrimSuffix(r.stderr, "\n"), "\n"); len(lines) != 1 {
				t.Errorf("pilot's refusal is not one line: %q", r.stderr)
			}
			if _, err := os.Stat(filepath.Join(dir, "argv")); err == nil {
				t.Error("the harness was started")
			}
		})
	}
}
