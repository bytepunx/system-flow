//go:build !windows

package serve

import (
	"context"
	"encoding/json"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/config"
	"github.com/bytepunx/system-flow/flai/internal/harness"
	"github.com/bytepunx/system-flow/flai/internal/hostapi"
	"github.com/bytepunx/system-flow/flai/internal/threads"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// claudeLab is a stand-in for claude, a flai serve state folder, and the
// projects it serves. The stand-in prints a version for --version; run
// with -p, it keeps its arguments and what the scratch project held, and
// either writes the file asked for in each worktree, or prints the events
// and the error of a refused write and fails.
type claudeLab struct {
	t       *testing.T
	stub    string
	out     string // what the stand-in kept
	tmp     string // the temporary folder scratch projects are made in
	logs    syncBuffer
	checker *claudeChecker
	o       Options
	served  []Entry
}

const (
	claudeLabVersion = "9.8.7"
	claudeLabError   = "Permission prompt tool returned an invalid result: expected behavior"
)

func newClaudeLab(t *testing.T, fail bool) *claudeLab {
	t.Helper()
	lab := &claudeLab{t: t, out: t.TempDir(), tmp: t.TempDir()}
	t.Setenv("TMPDIR", lab.tmp)
	lab.stub = filepath.Join(t.TempDir(), "claude")
	do := "for d in \"$PWD\"/.flai-cache/worktrees/*/; do mkdir -p \"$d.claude\"; printf '%s\\n' '" + claudeCheckContent(claudeLabVersion) + "' > \"${d}.claude/flai-check.md\"; done\n"
	if fail {
		do = "echo '{\"type\":\"system\",\"subtype\":\"init\",\"model\":\"claude-haiku\"}'\n" +
			"echo '{\"type\":\"assistant\",\"message\":{\"content\":[{\"type\":\"tool_use\",\"name\":\"Write\",\"input\":{\"file_path\":\"/scratch/.claude/flai-check.md\"}}]}}'\n" +
			"echo '{\"type\":\"user\",\"message\":{\"content\":[{\"type\":\"tool_result\",\"content\":\"" + claudeLabError + "\",\"is_error\":true}]}}'\n" +
			"echo '{\"type\":\"result\",\"subtype\":\"success\",\"is_error\":false,\"result\":\"I could not write the file.\"}'\n" +
			"echo 'stand-in: the write was refused' >&2\n" +
			"exit 1\n"
	}
	script := "#!/bin/sh\n" +
		"if [ \"$1\" = --version ]; then echo version >> \"" + lab.out + "/calls\"; echo '" + claudeLabVersion + " (Claude Code)'; exit 0; fi\n" +
		"echo p >> \"" + lab.out + "/calls\"\n" +
		"printf '%s\\n' \"$@\" > \"" + lab.out + "/argv\"\n" +
		"cp ../config.json \"" + lab.out + "/config.json\"\n" +
		"cp -r wip/kanban/stories \"" + lab.out + "/stories\"\n" +
		"pwd > \"" + lab.out + "/pwd\"\n" +
		do
	if err := os.WriteFile(lab.stub, []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	lab.o = Options{Dir: Dir(filepath.Join(t.TempDir(), "serve")), Logger: slog.New(slog.NewTextHandler(&lab.logs, nil)), Now: time.Now}
	lab.checker = newClaudeChecker(lab.o, func() []Entry { return lab.served })
	return lab
}

// config is the operator's say with the agent action on and the stand-in as
// claude-code's program.
func (lab *claudeLab) config() AgentConfig {
	return AgentConfig{Enabled: true, Flai: "/usr/local/bin/flai", Harnesses: map[string]harness.Host{harness.ClaudeCode: {Program: lab.stub}}}
}

func (lab *claudeLab) check() {
	lab.t.Helper()
	lab.checker.check(context.Background(), harness.ClaudeCodeHost(lab.config().host(harness.ClaudeCode)), "/usr/local/bin/flai")
}

// calls are the stand-in's runs, version for --version and p for -p.
func (lab *claudeLab) calls() []string {
	data, _ := os.ReadFile(filepath.Join(lab.out, "calls"))
	return strings.Fields(string(data))
}

func (lab *claudeLab) records() map[string]ClaudeCheck {
	lab.t.Helper()
	all, err := lab.o.Dir.claudeChecks()
	if err != nil {
		lab.t.Fatal(err)
	}
	return all
}

// serve adds a project to those served, its agents on harness when it names one.
func (lab *claudeLab) serve(key, harnessName string) *workitem.Repo {
	lab.t.Helper()
	root := lab.t.TempDir()
	m := "version: 1\nname: " + key + "\nkey: " + key + "\nlayout:\n  design: design\n  docs: docs\n  wip: wip\n"
	if harnessName != "" {
		m += "agent:\n  harness: " + harnessName + "\n"
	}
	if err := os.WriteFile(filepath.Join(root, "system-flow.yaml"), []byte(m), 0o644); err != nil {
		lab.t.Fatal(err)
	}
	for _, d := range []string{"wip/kanban/epics", "wip/kanban/stories", "wip/kanban/tasks", "wip/agents"} {
		if err := os.MkdirAll(filepath.Join(root, d), 0o755); err != nil {
			lab.t.Fatal(err)
		}
	}
	lab.served = append(lab.served, Entry{Key: key, Name: key, Root: root})
	repo, err := workitem.Open(root)
	if err != nil {
		lab.t.Fatal(err)
	}
	return repo
}

func threadsOf(t *testing.T, repo *workitem.Repo) []*threads.Thread {
	t.Helper()
	all, err := threads.List(repo)
	if err != nil {
		t.Fatal(err)
	}
	return all
}

// S-0286: a Claude Code version this host has not checked is checked once
// with a write under a scratch story's worktree's .claude/ folder, through
// permission_prompt, with the adapter's permission arguments, auto-approve
// on for the scratch project alone, and the cheapest model; it passes, is
// recorded, and the scratch project is removed.
func TestAnUncheckedClaudeCodeIsCheckedAndRecorded(t *testing.T) {
	lab := newClaudeLab(t, false)
	repo := lab.serve("harbour", harness.ClaudeCode)
	lab.check()

	if got := strings.Join(lab.calls(), " "); got != "version p" {
		t.Fatalf("the stand-in ran as %q, want its version read and one check; flai serve logged:\n%s", got, lab.logs.String())
	}
	rec, ok := lab.records()[claudeLabVersion]
	if !ok || rec.Outcome != ClaudeCheckPassed || rec.Error != "" || rec.Checked == "" || rec.Program != lab.stub {
		t.Fatalf("record %+v (found %v), want %s passed with when and the program", rec, ok, claudeLabVersion)
	}
	argv, _ := os.ReadFile(filepath.Join(lab.out, "argv"))
	args := strings.Split(strings.TrimSpace(string(argv)), "\n")
	joined := strings.Join(args, "\n")
	for _, want := range []string{"-p", "--model\nhaiku", "--permission-prompt-tool\n" + harness.PermissionPromptTool, "--permission-mode\nacceptEdits", "--strict-mcp-config"} {
		if !strings.Contains(joined, want) {
			t.Errorf("the check's arguments lack %q:\n%s", want, joined)
		}
	}
	var mcp struct {
		MCPServers map[string]struct {
			Command string   `json:"command"`
			Args    []string `json:"args"`
		} `json:"mcpServers"`
	}
	for i, a := range args {
		if a == "--mcp-config" && i+1 < len(args) {
			if err := json.Unmarshal([]byte(args[i+1]), &mcp); err != nil {
				t.Fatal(err)
			}
		}
	}
	server := mcp.MCPServers["flai"]
	if server.Command != "/usr/local/bin/flai" || len(server.Args) < 4 || server.Args[0] != "--config" || server.Args[2] != "mcp" {
		t.Errorf("the MCP server is %+v, want this flai's mcp with the scratch configuration", server)
	}
	scratch, _ := os.ReadFile(filepath.Join(lab.out, "pwd"))
	root := strings.TrimSpace(string(scratch))
	if server.Args[1] != filepath.Join(filepath.Dir(root), "config.json") {
		t.Errorf("the MCP server's configuration is %s, want the scratch one beside %s", server.Args[1], root)
	}
	cfg, _, err := config.Load(filepath.Join(lab.out, "config.json"))
	if err != nil {
		t.Fatal(err)
	}
	if !cfg.ActionEnabled(hostapi.ActionAutoApprove, root) {
		t.Errorf("auto-approve is not on in the scratch configuration: %+v", cfg.HostActions)
	}
	stories, _ := filepath.Glob(filepath.Join(lab.out, "stories", "S-*.md"))
	if len(stories) != 1 {
		t.Fatalf("the scratch project holds stories %v, want one", stories)
	}
	if data, _ := os.ReadFile(stories[0]); !strings.Contains(string(data), "status: in-progress") {
		t.Errorf("the scratch story is not in progress:\n%s", data)
	}
	if left, _ := os.ReadDir(lab.tmp); len(left) != 0 {
		t.Errorf("the scratch folder is left behind: %v", left)
	}
	if n := len(threadsOf(t, repo)); n != 0 {
		t.Errorf("a passed check opened %d threads", n)
	}
	logs := lab.logs.String()
	for _, want := range []string{`msg="claude code check started"`, `msg="claude code check ended"`, "version=" + claudeLabVersion, "outcome=passed"} {
		if !strings.Contains(logs, want) {
			t.Errorf("the log lacks %s:\n%s", want, logs)
		}
	}
}

// S-0286: a version with a record is not checked again, whatever its outcome.
func TestARecordedClaudeCodeIsNotCheckedAgain(t *testing.T) {
	for _, outcome := range []string{ClaudeCheckPassed, ClaudeCheckFailed} {
		t.Run(outcome, func(t *testing.T) {
			lab := newClaudeLab(t, false)
			repo := lab.serve("harbour", harness.ClaudeCode)
			before := ClaudeCheck{Checked: "2026-10-01T00:00:00Z", Program: lab.stub, Outcome: outcome}
			if err := lab.o.Dir.recordClaudeCheck(claudeLabVersion, before); err != nil {
				t.Fatal(err)
			}
			lab.check()
			if got := strings.Join(lab.calls(), " "); got != "version" {
				t.Errorf("the stand-in ran as %q, want its version read and nothing else", got)
			}
			if got := lab.records()[claudeLabVersion]; got != before {
				t.Errorf("the record became %+v", got)
			}
			if n := len(threadsOf(t, repo)); n != 0 {
				t.Errorf("%d threads opened", n)
			}
		})
	}
}

// S-0286: a failed check is recorded with Claude Code's error, and told on
// one thread on the manifest of each served project whose agents run
// claude-code, quoting it and naming the version; a project on another
// harness hears nothing.
func TestAFailedClaudeCodeCheckOpensOneThreadQuotingItsError(t *testing.T) {
	lab := newClaudeLab(t, true)
	claude := lab.serve("harbour", harness.ClaudeCode)
	other := lab.serve("quay", harness.Command)
	lab.check()

	rec := lab.records()[claudeLabVersion]
	if rec.Outcome != ClaudeCheckFailed || !strings.Contains(rec.Error, claudeLabError) || !strings.Contains(rec.Error, "stand-in: the write was refused") || !strings.Contains(rec.Error, "was not written") {
		t.Fatalf("record %+v, want failed with the stand-in's error", rec)
	}
	all := threadsOf(t, claude)
	if len(all) != 1 {
		t.Fatalf("%d threads on the claude-code project, want one", len(all))
	}
	th := all[0]
	if th.Anchor.Path != "system-flow.yaml" || !th.Open() || th.Opener() != claudeCheckAuthor || !strings.Contains(th.Title, claudeLabVersion) {
		t.Errorf("thread %+v, want one open on system-flow.yaml by %s naming %s", th, claudeCheckAuthor, claudeLabVersion)
	}
	for _, want := range []string{"```text\n", claudeLabError, "Claude Code " + claudeLabVersion, ".mcp.json", "will not go through permission_prompt"} {
		if !strings.Contains(th.Body, want) {
			t.Errorf("the thread lacks %q:\n%s", want, th.Body)
		}
	}
	if n := len(threadsOf(t, other)); n != 0 {
		t.Errorf("%d threads on the project on another harness", n)
	}
	logs := lab.logs.String()
	if !strings.Contains(logs, `level=WARN msg="claude code check ended"`) || !strings.Contains(logs, "outcome=failed") {
		t.Errorf("the failure is not logged at warn:\n%s", logs)
	}

	lab.check() // checked: no second check, no second thread
	if got := strings.Join(lab.calls(), " "); got != "version p version" {
		t.Errorf("the stand-in ran as %q", got)
	}
	if n := len(threadsOf(t, claude)); n != 1 {
		t.Errorf("%d threads after a second look, want one", n)
	}
}

// S-0286: a look that may start a claude-code agent checks the claude it
// would run in the background, once for the same program; one with the
// agent actions off, or for a project on another harness, checks nothing.
func TestALookChecksTheClaudeItWouldRunOnce(t *testing.T) {
	lab := newClaudeLab(t, false)
	claude := lab.serve("harbour", harness.ClaudeCode)
	other := lab.serve("quay", harness.Command)
	ctx := context.Background()

	off := lab.config()
	off.Enabled = false
	lab.checker.consider(ctx, claude.Root, off)
	lab.checker.consider(ctx, other.Root, lab.config())
	lab.checker.wait()
	if got := lab.calls(); len(got) != 0 {
		t.Fatalf("the stand-in ran as %v with nothing to start on claude-code", got)
	}

	lab.checker.consider(ctx, claude.Root, lab.config())
	lab.checker.wait()
	lab.checker.consider(ctx, claude.Root, lab.config())
	lab.checker.wait()
	if got := strings.Join(lab.calls(), " "); got != "version p" {
		t.Fatalf("the stand-in ran as %q, want one check; flai serve logged:\n%s", got, lab.logs.String())
	}
	if rec := lab.records()[claudeLabVersion]; rec.Outcome != ClaudeCheckPassed {
		t.Fatalf("record %+v", rec)
	}

	// another claude, the same version: read, and not checked again
	later := time.Now().Add(time.Minute)
	if err := os.Chtimes(lab.stub, later, later); err != nil {
		t.Fatal(err)
	}
	lab.checker.consider(ctx, claude.Root, lab.config())
	lab.checker.wait()
	if got := strings.Join(lab.calls(), " "); got != "version p version" {
		t.Errorf("the stand-in ran as %q, want its version read again and no check", got)
	}
}
