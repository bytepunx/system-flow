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
	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/threads"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// claudeLab is a stand-in for claude, a flai serve state folder, and the
// projects it serves. The stand-in prints a version for --version; run
// with -p, it keeps its arguments, its environment, and what the scratch
// project held, and either writes the file asked for in each worktree, or
// prints the events and the error of a refused write and fails.
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
	// claudeLabGateway is the provider the lab's host has whole, its key in
	// claudeLabKeyEnv; claudeLabChat speaks an API claude does not.
	claudeLabGateway = "openrouter"
	claudeLabChat    = "chat"
	claudeLabKeyEnv  = "FLAI_TEST_CLAUDE_CHECK_KEY"
	claudeLabBaseURL = "https://gateway.example/api"
	claudeLabHaiku   = "vendor/haiku-on-host"
)

func newClaudeLab(t *testing.T, fail bool) *claudeLab {
	t.Helper()
	lab := &claudeLab{t: t, out: t.TempDir(), tmp: t.TempDir()}
	t.Setenv("TMPDIR", lab.tmp)
	// as the record has it, symbolic links followed: a macOS temp dir is a link under /private
	stubDir := t.TempDir()
	if real, err := filepath.EvalSymlinks(stubDir); err == nil {
		stubDir = real
	}
	lab.stub = filepath.Join(stubDir, "claude")
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
		"env > \"" + lab.out + "/env\"\n" +
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

// config is the operator's say with the agent action on, the stand-in as
// claude-code's program, and the host's providers claudeLabGateway and
// claudeLabChat.
func (lab *claudeLab) config() AgentConfig {
	return AgentConfig{Enabled: true, Flai: "/usr/local/bin/flai", Harnesses: map[string]harness.Host{harness.ClaudeCode: {Program: lab.stub}},
		Providers: map[string]manifest.Provider{
			claudeLabGateway: {API: manifest.APIAnthropicMessages, BaseURL: claudeLabBaseURL, KeyEnv: claudeLabKeyEnv, Models: map[string]string{"haiku": claudeLabHaiku}},
			claudeLabChat:    {API: manifest.APIOpenAIChat, BaseURL: "https://chat.example/v1", KeyEnv: claudeLabKeyEnv},
		}}
}

// check checks the stand-in through no provider.
func (lab *claudeLab) check() {
	lab.t.Helper()
	lab.checkThrough("", "")
}

// checkThrough checks the stand-in through provider, as the default agent of
// the project at root names it.
func (lab *claudeLab) checkThrough(root, provider string) {
	lab.t.Helper()
	lab.checker.check(context.Background(), harness.ClaudeCodeHost(lab.config().host(harness.ClaudeCode)), lab.config(), root, provider)
}

// env is the stand-in's environment at its last -p run, by variable.
func (lab *claudeLab) env() map[string]string {
	lab.t.Helper()
	data, err := os.ReadFile(filepath.Join(lab.out, "env"))
	if err != nil {
		lab.t.Fatalf("the stand-in kept no environment: %v", err)
	}
	out := map[string]string{}
	for line := range strings.Lines(string(data)) {
		if name, value, ok := strings.Cut(strings.TrimSuffix(line, "\n"), "="); ok {
			out[name] = value
		}
	}
	return out
}

// calls are the stand-in's runs, version for --version and p for -p.
func (lab *claudeLab) calls() []string {
	data, _ := os.ReadFile(filepath.Join(lab.out, "calls"))
	return strings.Fields(string(data))
}

// records are the host's records, by version and then by provider.
func (lab *claudeLab) records() map[string]map[string]ClaudeCheck {
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
	return lab.serveOn(key, harnessName, "", "")
}

// serveOn adds a project to those served, its agents on harness when it
// names one, its default agent naming provider when it names one, and more
// written at the end of its manifest.
func (lab *claudeLab) serveOn(key, harnessName, provider, more string) *workitem.Repo {
	lab.t.Helper()
	root := lab.t.TempDir()
	m := "version: 1\nname: " + key + "\nkey: " + key + "\nlayout:\n  design: design\n  docs: docs\n  wip: wip\n"
	if harnessName != "" {
		m += "agent:\n  harness: " + harnessName + "\n"
	}
	if provider != "" {
		m += "  provider: " + provider + "\n"
	}
	m += more
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
	rec, ok := lab.records()[claudeLabVersion][ClaudeCheckAnthropic]
	if !ok || rec.Outcome != ClaudeCheckPassed || rec.Error != "" || rec.Checked == "" || rec.Program != lab.stub {
		t.Fatalf("record %+v (found %v), want %s passed through %s with when and the program", rec, ok, claudeLabVersion, ClaudeCheckAnthropic)
	}
	// S-0360: through no provider, the environment is an agent's, with no gateway's
	env := lab.env()
	if env["ANTHROPIC_BASE_URL"] == claudeLabBaseURL || env["ANTHROPIC_DEFAULT_HAIKU_MODEL"] == claudeLabHaiku || env["FLAI_AGENT"] != claudeCheckAgent {
		t.Errorf("the check through no provider ran with ANTHROPIC_BASE_URL %q, ANTHROPIC_DEFAULT_HAIKU_MODEL %q, FLAI_AGENT %q; want no gateway's and %s", env["ANTHROPIC_BASE_URL"], env["ANTHROPIC_DEFAULT_HAIKU_MODEL"], env["FLAI_AGENT"], claudeCheckAgent)
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
	for _, want := range []string{`msg="claude code check started"`, `msg="claude code check ended"`, "version=" + claudeLabVersion, "provider=" + ClaudeCheckAnthropic, "outcome=passed"} {
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
			if err := lab.o.Dir.recordClaudeCheck(claudeLabVersion, ClaudeCheckAnthropic, before); err != nil {
				t.Fatal(err)
			}
			lab.check()
			if got := strings.Join(lab.calls(), " "); got != "version" {
				t.Errorf("the stand-in ran as %q, want its version read and nothing else", got)
			}
			if got := lab.records()[claudeLabVersion][ClaudeCheckAnthropic]; got != before {
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
// harness hears nothing. S-0360: nor does one whose default agent names
// another provider than the check ran through, and the thread names it.
func TestAFailedClaudeCodeCheckOpensOneThreadQuotingItsError(t *testing.T) {
	lab := newClaudeLab(t, true)
	claude := lab.serve("harbour", harness.ClaudeCode)
	other := lab.serve("quay", harness.Command)
	gateway := lab.serveOn("pier", harness.ClaudeCode, claudeLabGateway, "")
	lab.check()

	rec := lab.records()[claudeLabVersion][ClaudeCheckAnthropic]
	if rec.Outcome != ClaudeCheckFailed || !strings.Contains(rec.Error, claudeLabError) || !strings.Contains(rec.Error, "stand-in: the write was refused") || !strings.Contains(rec.Error, "was not written") {
		t.Fatalf("record %+v, want failed with the stand-in's error", rec)
	}
	all := threadsOf(t, claude)
	if len(all) != 1 {
		t.Fatalf("%d threads on the claude-code project, want one", len(all))
	}
	th := all[0]
	if th.Anchor.Path != "system-flow.yaml" || !th.Open() || th.Opener() != claudeCheckAuthor || !strings.Contains(th.Title, claudeLabVersion+" through "+ClaudeCheckAnthropic) {
		t.Errorf("thread %+v, want one open on system-flow.yaml by %s naming %s through %s", th, claudeCheckAuthor, claudeLabVersion, ClaudeCheckAnthropic)
	}
	for _, want := range []string{"```text\n", claudeLabError, "Claude Code " + claudeLabVersion, "through " + ClaudeCheckAnthropic, ".mcp.json", "will not go through permission_prompt", "remove the " + ClaudeCheckAnthropic + " entry under " + claudeLabVersion} {
		if !strings.Contains(th.Body, want) {
			t.Errorf("the thread lacks %q:\n%s", want, th.Body)
		}
	}
	if n := len(threadsOf(t, other)); n != 0 {
		t.Errorf("%d threads on the project on another harness", n)
	}
	if n := len(threadsOf(t, gateway)); n != 0 {
		t.Errorf("%d threads on the project whose agent names %s", n, claudeLabGateway)
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
	if rec := lab.records()[claudeLabVersion][ClaudeCheckAnthropic]; rec.Outcome != ClaudeCheckPassed {
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

// S-0360: a project whose default agent names a provider has the check run
// through it: the adapter's provider environment, from the host's entry with
// the project's overrides laid over, the haiku alias resolved through its
// models, and the key copied from its variable; it is recorded under the
// provider, and the key is in no log line or record.
func TestAClaudeCodeCheckRunsThroughTheProviderTheDefaultAgentNames(t *testing.T) {
	lab := newClaudeLab(t, false)
	const key = "sk-test-not-a-real-key"
	t.Setenv(claudeLabKeyEnv, key)
	t.Setenv("ANTHROPIC_API_KEY", "sk-ant-inherited")
	repo := lab.serveOn("harbour", harness.ClaudeCode, claudeLabGateway, "providers:\n  "+claudeLabGateway+":\n    models:\n      sonnet: vendor/sonnet-in-project\n")
	lab.checkThrough(repo.Root, claudeLabGateway)

	if got := strings.Join(lab.calls(), " "); got != "version p" {
		t.Fatalf("the stand-in ran as %q, want its version read and one check; flai serve logged:\n%s", got, lab.logs.String())
	}
	env := lab.env()
	for name, want := range map[string]string{
		"ANTHROPIC_BASE_URL":             claudeLabBaseURL,
		"ANTHROPIC_DEFAULT_HAIKU_MODEL":  claudeLabHaiku,
		"ANTHROPIC_DEFAULT_SONNET_MODEL": "vendor/sonnet-in-project",
		"ANTHROPIC_AUTH_TOKEN":           key,
		"ANTHROPIC_API_KEY":              "",
		"FLAI_AGENT":                     claudeCheckAgent,
	} {
		if got, ok := env[name]; !ok || got != want {
			t.Errorf("the check ran with %s %q (set %v), want %q", name, got, ok, want)
		}
	}
	argv, _ := os.ReadFile(filepath.Join(lab.out, "argv"))
	if !strings.Contains(string(argv), "--model\nhaiku\n") {
		t.Errorf("the check's arguments lack --model haiku:\n%s", argv)
	}
	all := lab.records()[claudeLabVersion]
	if rec, ok := all[claudeLabGateway]; !ok || rec.Outcome != ClaudeCheckPassed {
		t.Errorf("records %+v, want %s passed through %s", all, claudeLabVersion, claudeLabGateway)
	}
	if _, ok := all[ClaudeCheckAnthropic]; ok {
		t.Errorf("records %+v: the check through %s is recorded as %s too", all, claudeLabGateway, ClaudeCheckAnthropic)
	}
	logs := lab.logs.String()
	if !strings.Contains(logs, `msg="claude code check ended"`) || !strings.Contains(logs, "provider="+claudeLabGateway) {
		t.Errorf("the log does not name the provider the check ran through:\n%s", logs)
	}
	file, _ := os.ReadFile(lab.o.Dir.claudeChecksFile())
	if strings.Contains(logs, key) || strings.Contains(string(file), key) {
		t.Errorf("the key's value reached the log or the record")
	}
}

// S-0360: a version that passed through Anthropic is checked again through a
// provider the first time a served project's default agent names it, and not
// again after that, by this flai serve or by the next one.
func TestAVersionPassedThroughOneProviderIsCheckedThroughAnother(t *testing.T) {
	lab := newClaudeLab(t, false)
	t.Setenv(claudeLabKeyEnv, "sk-test-not-a-real-key")
	direct := lab.serve("harbour", harness.ClaudeCode)
	gateway := lab.serveOn("pier", harness.ClaudeCode, claudeLabGateway, "")
	ctx := context.Background()

	lab.checker.consider(ctx, direct.Root, lab.config())
	lab.checker.wait()
	if got := strings.Join(lab.calls(), " "); got != "version p" {
		t.Fatalf("the stand-in ran as %q, want one check through %s; flai serve logged:\n%s", got, ClaudeCheckAnthropic, lab.logs.String())
	}
	lab.checker.consider(ctx, gateway.Root, lab.config())
	lab.checker.wait()
	if got := strings.Join(lab.calls(), " "); got != "version p version p" {
		t.Fatalf("the stand-in ran as %q, want a second check through %s; flai serve logged:\n%s", got, claudeLabGateway, lab.logs.String())
	}
	if env := lab.env(); env["ANTHROPIC_BASE_URL"] != claudeLabBaseURL {
		t.Errorf("the second check ran with ANTHROPIC_BASE_URL %q, want %s", env["ANTHROPIC_BASE_URL"], claudeLabBaseURL)
	}
	all := lab.records()[claudeLabVersion]
	if all[ClaudeCheckAnthropic].Outcome != ClaudeCheckPassed || all[claudeLabGateway].Outcome != ClaudeCheckPassed {
		t.Fatalf("records %+v, want %s passed through %s and %s", all, claudeLabVersion, ClaudeCheckAnthropic, claudeLabGateway)
	}

	lab.checker.consider(ctx, direct.Root, lab.config())
	lab.checker.wait()
	lab.checker.consider(ctx, gateway.Root, lab.config())
	lab.checker.wait()
	if got := strings.Join(lab.calls(), " "); got != "version p version p" {
		t.Errorf("the stand-in ran as %q after more looks, want nothing more", got)
	}
	next := newClaudeChecker(lab.o, func() []Entry { return lab.served })
	next.consider(ctx, gateway.Root, lab.config())
	next.wait()
	if got := strings.Join(lab.calls(), " "); got != "version p version p version" {
		t.Errorf("the stand-in ran as %q for the next flai serve, want its version read and no check", got)
	}
}

// S-0360: a provider whose key's variable is unset or empty is not checked,
// nothing is recorded, no thread is opened, and the log says the check was
// skipped and names the variable; one claude cannot call, of another API, is
// not checked either, and the log says why.
func TestAProviderWithoutItsKeyIsNotChecked(t *testing.T) {
	for _, c := range []struct{ provider, msg, why string }{
		{claudeLabGateway, `level=WARN msg="claude code check skipped"`, claudeLabKeyEnv + ", which is unset or empty"},
		{claudeLabChat, `level=WARN msg="claude code not checked"`, "claude-code speaks " + manifest.APIAnthropicMessages + " alone"},
	} {
		t.Run(c.provider, func(t *testing.T) {
			lab := newClaudeLab(t, true)
			t.Setenv(claudeLabKeyEnv, "")
			repo := lab.serveOn("harbour", harness.ClaudeCode, c.provider, "")
			lab.checkThrough(repo.Root, c.provider)

			if got := strings.Join(lab.calls(), " "); got != "version" {
				t.Errorf("the stand-in ran as %q, want its version read and no check", got)
			}
			if all := lab.records(); len(all) != 0 {
				t.Errorf("records %+v, want none", all)
			}
			if n := len(threadsOf(t, repo)); n != 0 {
				t.Errorf("%d threads opened", n)
			}
			logs := lab.logs.String()
			for _, want := range []string{c.msg, "provider=" + c.provider, "version=" + claudeLabVersion, c.why} {
				if !strings.Contains(logs, want) {
					t.Errorf("the log lacks %s:\n%s", want, logs)
				}
			}
		})
	}
}

// S-0360: a record written before checks were made per provider, one record
// for a version, reads as checked through anthropic, and the next record
// writes the file by version and then by provider.
func TestAnOlderClaudeChecksFileReadsAsAnthropic(t *testing.T) {
	lab := newClaudeLab(t, false)
	lab.serve("harbour", harness.ClaudeCode)
	if err := os.MkdirAll(string(lab.o.Dir), 0o755); err != nil {
		t.Fatal(err)
	}
	old := `{"` + claudeLabVersion + `": {"checked": "2026-10-01T00:00:00Z", "program": "/opt/claude", "outcome": "failed", "error": "refused"}}`
	if err := os.WriteFile(lab.o.Dir.claudeChecksFile(), []byte(old), 0o644); err != nil {
		t.Fatal(err)
	}
	want := ClaudeCheck{Checked: "2026-10-01T00:00:00Z", Program: "/opt/claude", Outcome: ClaudeCheckFailed, Error: "refused"}
	if got := lab.records()[claudeLabVersion]; len(got) != 1 || got[ClaudeCheckAnthropic] != want {
		t.Fatalf("the older record reads as %+v, want %+v under %s alone", got, want, ClaudeCheckAnthropic)
	}
	lab.check()
	if got := strings.Join(lab.calls(), " "); got != "version" {
		t.Errorf("the stand-in ran as %q, want its version read and no check", got)
	}

	if err := lab.o.Dir.recordClaudeCheck(claudeLabVersion, claudeLabGateway, ClaudeCheck{Checked: "2026-10-02T00:00:00Z", Program: "/opt/claude", Outcome: ClaudeCheckPassed}); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(lab.o.Dir.claudeChecksFile())
	if err != nil {
		t.Fatal(err)
	}
	var written map[string]map[string]ClaudeCheck
	if err := json.Unmarshal(data, &written); err != nil {
		t.Fatalf("the file is not by version and then by provider: %v\n%s", err, data)
	}
	if got := written[claudeLabVersion]; got[ClaudeCheckAnthropic] != want || got[claudeLabGateway].Outcome != ClaudeCheckPassed {
		t.Errorf("the file holds %+v, want the older record under %s and the new one under %s", got, ClaudeCheckAnthropic, claudeLabGateway)
	}
}
