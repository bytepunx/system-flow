package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/channel"
	"github.com/bytepunx/system-flow/flai/internal/config"
	"github.com/bytepunx/system-flow/flai/internal/host"
	"github.com/bytepunx/system-flow/flai/internal/hostapi"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/serve"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// inProcess runs the flai commands behind the host's write methods in this
// process: under go test the executable is the test binary, not flai.
func inProcess(t *testing.T) {
	t.Helper()
	was := hostapi.Commands
	t.Cleanup(func() { hostapi.Commands = was })
	t.Setenv("LOG_FORMAT", "json")
	hostapi.Commands = func(_ context.Context, r hostapi.Run) (hostapi.Ran, error) {
		out, errOut, code := runIn(t, r.Dir, r.Args...)
		ran := hostapi.Ran{Stdout: []byte(out), Exit: code}
		for _, line := range strings.Split(errOut, "\n") {
			var ev map[string]any
			if json.Unmarshal([]byte(line), &ev) == nil {
				ran.Events = append(ran.Events, ev)
			}
		}
		return ran, nil
	}
}

const publishRequest = `{"request_id":"3f0c1a52-7d3b-4f0e-9a51-0c2d4e6f8a10"}`

// S-0078: a host action is off until the operator enables it by name on the
// host; asked for meanwhile it is refused and says what enables it; enabled,
// it runs; and every request is in the journal, refused ones too. The push
// action gates Publish, the one way accepted work reaches the remote
// (ADR-0067).
func TestHostActionPush(t *testing.T) {
	root, remote := researchProject(t, "feature", true)
	inProcess(t)
	head := func(dir string) string { return strings.TrimSpace(gitIn(t, dir, "rev-parse", "main")) }
	before := head(remote)
	// S-0087: acceptance itself never pushes or tags, whatever the push
	// action is set to; it only merges, archives, and commits locally.
	// Releasing waits for Publish, so acceptances batch into one release.
	out, _, code := runIn(t, root, "hostapi", "accept.run", `{"id":"S-0001","include_uncommitted":true,"request_id":"3f0c1a52-7d3b-4f0e-9a51-0c2d4e6f8a09"}`)
	if code != 0 || head(remote) != before || head(root) == before {
		t.Fatalf("an acceptance with the action off: accepted here, nothing pushed: %d %s", code, out)
	}

	out, _, _ = runIn(t, root, "serve", "actions")
	if !strings.Contains(out, "push: off everywhere") || !strings.Contains(out, "off for this project; flai serve enable push") || !strings.Contains(out, "publish everything accepted") {
		t.Errorf("actions, before: %s", out)
	}
	if j, _, _ := runIn(t, root, "serve", "journal"); !strings.Contains(j, "no host action has been asked for") {
		t.Errorf("an empty journal: %s", j)
	}
	info, _, _ := runIn(t, root, "hostapi", "project.info")
	if !strings.Contains(info, `"host_actions"`) || !strings.Contains(info, `"push":false`) {
		t.Errorf("the dashboard is told the action is off: %s", info)
	}

	// off: refused, nothing pushed
	out, _, code = runIn(t, root, "hostapi", "publish.run", publishRequest)
	if code == 0 || !strings.Contains(out, "is not enabled for this project") || !strings.Contains(out, "flai serve enable push") || head(remote) != before {
		t.Fatalf("disabled: %d %s", code, out)
	}

	if _, errOut, code := runIn(t, root, "serve", "enable", "pull"); code == 0 || !strings.Contains(errOut, "there are: agent, analyze, auto-approve, auto-publish, auto-upgrade, checks, components, curate, dashboard, forge, host, merge-queue, orchestrate, plan, push") {
		t.Errorf("an action there is not: %d %s", code, errOut)
	}
	out, _, code = runIn(t, root, "serve", "enable", "push")
	if code != 0 || !strings.Contains(out, "push enabled for t") || !strings.Contains(out, "flai serve journal") || !strings.Contains(out, "flai serve disable push") {
		t.Fatalf("enable: %d %s", code, out)
	}
	cfg, _ := os.ReadFile(os.Getenv("FLAI_CONFIG"))
	if !strings.Contains(string(cfg), `"host_actions"`) || !strings.Contains(string(cfg), root) {
		t.Errorf("it is written to the host's configuration, for this project: %s", cfg)
	}
	// another project on the same host is not enabled by it
	other := tempProject(t)
	if info, _, _ := runIn(t, other, "hostapi", "project.info"); !strings.Contains(info, `"push":false`) {
		t.Errorf("per project: %s", info)
	}

	// on: published, the acceptance tagged and pushed with its tag
	out, _, code = runIn(t, root, "hostapi", "publish.run", publishRequest)
	if code != 0 || !strings.Contains(out, `"pushed":true`) || !strings.Contains(out, "cli/v1.1.0") || head(remote) != head(root) {
		t.Fatalf("enabled: %d %s", code, out)
	}
	if tags := gitIn(t, remote, "tag", "--list"); !strings.Contains(tags, "cli/v1.1.0") {
		t.Errorf("the release tag reached the remote: %s", tags)
	}
	out, _, _ = runIn(t, root, "hostapi", "publish.run", `{"request_id":"3f0c1a52-7d3b-4f0e-9a51-0c2d4e6f8a11"}`)
	if !strings.Contains(out, "nothing pending") {
		t.Errorf("nothing left: %s", out)
	}

	// the journal has all three, in order, and says for whom
	js, _, _ := runIn(t, root, "serve", "journal", "--json")
	var entries []struct{ Action, Method, Outcome, Detail, By, Root, At string }
	if err := json.Unmarshal([]byte(js), &entries); err != nil || len(entries) != 3 {
		t.Fatalf("journal: %v %s", err, js)
	}
	for i, want := range []struct{ outcome, detail string }{{"disabled", ""}, {"done", "pushed with tags cli/v1.1.0"}, {"done", "nothing pushed: nothing pending"}} {
		e := entries[i]
		if e.Outcome != want.outcome || e.Detail != want.detail || e.Action != "push" || e.Method != "publish.run" || e.Root != root || e.By == "" || e.At == "" {
			t.Errorf("entry %d: %+v", i, e)
		}
	}
	text, _, _ := runIn(t, root, "serve", "journal", "-n", "1")
	if !strings.Contains(text, "nothing pushed") || strings.Contains(text, "disabled") {
		t.Errorf("journal -n 1: %s", text)
	}
	if st, err := os.Stat(filepath.Join(filepath.Dir(os.Getenv("FLAI_CONFIG")), "serve", "journal.jsonl")); err != nil || st.Mode().Perm() != 0o600 {
		t.Errorf("the journal is the operator's alone: %v %v", st, err)
	}

	// off again, at once
	if out, _, _ := runIn(t, root, "serve", "disable", "push"); !strings.Contains(out, "push disabled for t") {
		t.Errorf("disable: %s", out)
	}
	if out, _, code := runIn(t, root, "hostapi", "publish.run", `{"request_id":"3f0c1a52-7d3b-4f0e-9a51-0c2d4e6f8a12"}`); code == 0 || !strings.Contains(out, "not enabled") {
		t.Errorf("disabled again: %d %s", code, out)
	}
	// every project, and one project's disable under it says so
	runIn(t, root, "serve", "enable", "push", "--all-projects")
	if out, _, _ := runIn(t, root, "serve", "disable", "push"); !strings.Contains(out, "still enabled") || !strings.Contains(out, "--all-projects") {
		t.Errorf("disable under all projects: %s", out)
	}
	if out, _, _ := runIn(t, root, "serve", "disable", "push", "--all-projects"); !strings.Contains(out, "push disabled for every project") {
		t.Errorf("disable everywhere: %s", out)
	}
	if out, _, _ := runIn(t, root, "serve", "actions"); !strings.Contains(out, "push: off everywhere") {
		t.Errorf("actions, after: %s", out)
	}
}

// S-0079: the agent's command is the operator's, an argument list with no
// default, managed on the host and nowhere else.
func TestServeAgentCommand(t *testing.T) {
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	root := tempProject(t)
	out, _, _ := runIn(t, root, "serve", "agent")
	if !strings.Contains(out, "no command is set, so a story that names no harness is not started") || !strings.Contains(out, "off for this project") {
		t.Errorf("by default: %s", out)
	}
	if out, _, _ := runIn(t, root, "serve", "actions"); !strings.Contains(out, "agent: off everywhere") || !strings.Contains(out, "whoever can move a story to ready") {
		t.Errorf("actions names it and what it means: %s", out)
	}
	out, errOut, code := runIn(t, root, "serve", "agent", "set", "--name", "builder", "--", "claude", "-p", "work on {story}; echo $HOME")
	if code != 0 || !strings.Contains(out, `command: "claude" "-p" "work on {story}; echo $HOME"`) || !strings.Contains(out, "never through a shell") {
		t.Fatalf("set: %d %s %s", code, out, errOut)
	}
	cfg, _ := os.ReadFile(os.Getenv("FLAI_CONFIG"))
	if !strings.Contains(string(cfg), `"work on {story}; echo $HOME"`) || !strings.Contains(string(cfg), `"name": "builder"`) {
		t.Errorf("kept as an argument list, as written: %s", cfg)
	}
	if _, errOut, code := runIn(t, root, "config", "set", "agent.command", "rm -rf /"); code == 0 {
		t.Errorf("flai config set does not reach it: %s", errOut)
	}
	if _, _, code := runIn(t, root, "serve", "agent", "set"); code == 0 {
		t.Error("set needs a program")
	}
	runIn(t, root, "serve", "enable", "agent")
	js, _, _ := runIn(t, root, "serve", "agent", "show", "--json")
	if !strings.Contains(js, `"enabled_here": true`) || !strings.Contains(js, `"claude"`) {
		t.Errorf("show --json: %s", js)
	}
	got := (&app{}).agentConfig(root)
	if !got.Enabled || len(got.Command) != 3 || got.Name != "builder" {
		t.Errorf("what flai serve is given: %+v", got)
	}
	if out, _, _ := runIn(t, root, "serve", "agent", "clear"); !strings.Contains(out, "no command is set") {
		t.Errorf("clear: %s", out)
	}
	if got := (&app{}).agentConfig(root); len(got.Command) != 0 {
		t.Errorf("cleared: %+v", got)
	}
}

// S-0218: orchestrate is off until the operator enables it, flai serve is
// told at every look whether it is on, and agent.status carries the
// project's orchestrator run, null before the first.
func TestServeOrchestrateAction(t *testing.T) {
	cfgPath := filepath.Join(t.TempDir(), "cfg.json")
	t.Setenv("FLAI_CONFIG", cfgPath)
	root := tempProject(t)
	if out, _, _ := runIn(t, root, "serve", "actions"); !strings.Contains(out, "orchestrate: off everywhere") || !strings.Contains(out, "orchestration.permissions") {
		t.Errorf("actions names it and what it means: %s", out)
	}
	if (&app{}).agentConfig(root).Orchestrate {
		t.Error("orchestrate is on by default")
	}
	status := func() string {
		t.Helper()
		got, err := json.Marshal((&app{}).host().Agent(root))
		if err != nil {
			t.Fatal(err)
		}
		return string(got)
	}
	if got := status(); !strings.Contains(got, `"orchestrator":null`) {
		t.Errorf("agent.status before any run: %s", got)
	}
	if out, errOut, code := runIn(t, root, "serve", "enable", "orchestrate"); code != 0 || !strings.Contains(out, "orchestrate enabled for t") {
		t.Fatalf("enable: %d %s %s", code, out, errOut)
	}
	if !(&app{}).agentConfig(root).Orchestrate {
		t.Error("enabled, flai serve is not told")
	}
	run := `{"` + root + `":{"orchestrator":{"story":"","harness":"command","command":"run-agent","agent":"orchestrator","pid":42,"started":"2026-10-05T10:00:00Z","session":"s","log":"/l/t-orchestrator-20261005T100000Z.log"}}}`
	dir := string(serve.DirFor(cfgPath))
	if err := os.MkdirAll(dir, 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "agents.json"), []byte(run), 0o600); err != nil {
		t.Fatal(err)
	}
	if got := status(); !strings.Contains(got, `"orchestrator":{"story":"","harness":"command","command":"run-agent","agent":"orchestrator","pid":42,"started":"2026-10-05T10:00:00Z"`) || !strings.Contains(got, `"running":null`) {
		t.Errorf("agent.status with a run: %s", got)
	}
	runIn(t, root, "serve", "disable", "orchestrate")
	if (&app{}).agentConfig(root).Orchestrate {
		t.Error("disabled, flai serve is still told it is on")
	}
}

// S-0399, ADR-0173: merge-queue is off until the operator enables it, flai
// serve actions names it and what it means, the dashboard's settings page
// lists it, and enable and disable turn it on and off for the project.
func TestServeMergeQueueAction(t *testing.T) {
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	root := tempProject(t)
	enabled := func() bool {
		t.Helper()
		cfg, _, err := (&app{}).loadConfig()
		if err != nil {
			t.Fatal(err)
		}
		return cfg.ActionEnabled(hostapi.ActionMergeQueue, root)
	}
	listed := func() map[string]any {
		t.Helper()
		got, err := json.Marshal((&app{}).hostSettings(root))
		if err != nil {
			t.Fatal(err)
		}
		var settings struct {
			Actions []map[string]any `json:"actions"`
		}
		if err := json.Unmarshal(got, &settings); err != nil {
			t.Fatal(err)
		}
		for _, a := range settings.Actions {
			if a["name"] == hostapi.ActionMergeQueue {
				return a
			}
		}
		return nil
	}
	if out, _, _ := runIn(t, root, "serve", "actions"); !strings.Contains(out, "merge-queue: off everywhere") || !strings.Contains(out, "merge_queue.batch_max") || !strings.Contains(out, "orchestration.permissions.accept_reviews") {
		t.Errorf("actions names it and what it means: %s", out)
	}
	if enabled() {
		t.Error("merge-queue is on by default")
	}
	if a := listed(); a == nil || a["here"] != false || a["everywhere"] != false {
		t.Errorf("the settings page lists it, off: %v", a)
	}
	if out, errOut, code := runIn(t, root, "serve", "enable", "merge-queue"); code != 0 || !strings.Contains(out, "merge-queue enabled for t") || !strings.Contains(out, "flai serve disable merge-queue") {
		t.Fatalf("enable: %d %s %s", code, out, errOut)
	}
	if !enabled() {
		t.Error("enabled, the host's configuration does not say so")
	}
	if a := listed(); a == nil || a["here"] != true {
		t.Errorf("the settings page lists it, on: %v", a)
	}
	if out, errOut, code := runIn(t, root, "serve", "disable", "merge-queue"); code != 0 || !strings.Contains(out, "merge-queue disabled for t") {
		t.Fatalf("disable: %d %s %s", code, out, errOut)
	}
	if enabled() {
		t.Error("disabled, the host's configuration still has it on")
	}
}

// S-0223: analyze is off until the operator enables it, flai serve is told
// whether it is on, and agent.status carries the project's newest analyzer
// run, with its focus, trigger, outcome, report, and log, null before the
// first, apart from the orchestrator's.
func TestServeAnalyzeAction(t *testing.T) {
	cfgPath := filepath.Join(t.TempDir(), "cfg.json")
	t.Setenv("FLAI_CONFIG", cfgPath)
	root := tempProject(t)
	if out, _, _ := runIn(t, root, "serve", "actions"); !strings.Contains(out, "analyze: off everywhere") || !strings.Contains(out, "analysis.agent") || !strings.Contains(out, "design/analysis") {
		t.Errorf("actions names it and what it means: %s", out)
	}
	if (&app{}).agentConfig(root).Analyze {
		t.Error("analyze is on by default")
	}
	status := func() string {
		t.Helper()
		got, err := json.Marshal((&app{}).host().Agent(root))
		if err != nil {
			t.Fatal(err)
		}
		return string(got)
	}
	if got := status(); !strings.Contains(got, `"analyzer":null`) {
		t.Errorf("agent.status before any run: %s", got)
	}
	if out, errOut, code := runIn(t, root, "serve", "enable", "analyze"); code != 0 || !strings.Contains(out, "analyze enabled for t") {
		t.Fatalf("enable: %d %s %s", code, out, errOut)
	}
	if !(&app{}).agentConfig(root).Analyze {
		t.Error("enabled, flai serve is not told")
	}
	run := `{"` + root + `":{"analyzer":{"story":"","harness":"command","command":"run-agent","agent":"analyzer","pid":42,"started":"2026-10-06T10:00:00Z","ended":"2026-10-06T10:05:00Z","log":"/l/t-analyzer-20261006T100000Z.log","trigger":"asked","focus":"risk","report":"design/analysis/2026-10-06-risk.md","outcome":"worked"}}}`
	dir := string(serve.DirFor(cfgPath))
	if err := os.MkdirAll(dir, 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "agents.json"), []byte(run), 0o600); err != nil {
		t.Fatal(err)
	}
	got := status()
	for _, want := range []string{`"analyzer":{"story":"","harness":"command","command":"run-agent","agent":"analyzer","pid":42,"started":"2026-10-06T10:00:00Z","ended":"2026-10-06T10:05:00Z"`,
		`"log":"/l/t-analyzer-20261006T100000Z.log"`, `"trigger":"asked","focus":"risk","report":"design/analysis/2026-10-06-risk.md","outcome":"worked"`, `"orchestrator":null`, `"running":null`} {
		if !strings.Contains(got, want) {
			t.Errorf("agent.status with a run lacks %s: %s", want, got)
		}
	}
	runIn(t, root, "serve", "disable", "analyze")
	if (&app{}).agentConfig(root).Analyze {
		t.Error("disabled, flai serve is still told it is on")
	}
}

// S-0431: curate is off until the operator enables it, flai serve is told
// whether it is on, and agent.status carries the project's newest curator
// run, with its since, trigger, outcome, entries, and log, null before the
// first, apart from the analyzer's.
func TestServeCurateAction(t *testing.T) {
	cfgPath := filepath.Join(t.TempDir(), "cfg.json")
	t.Setenv("FLAI_CONFIG", cfgPath)
	root := tempProject(t)
	if out, _, _ := runIn(t, root, "serve", "actions"); !strings.Contains(out, "curate: off everywhere") || !strings.Contains(out, "digest.agent") || !strings.Contains(out, "flai digest new and edit") {
		t.Errorf("actions names it and what it means: %s", out)
	}
	if (&app{}).agentConfig(root).Curate {
		t.Error("curate is on by default")
	}
	status := func() string {
		t.Helper()
		got, err := json.Marshal((&app{}).host().Agent(root))
		if err != nil {
			t.Fatal(err)
		}
		return string(got)
	}
	if got := status(); !strings.Contains(got, `"curator":null`) {
		t.Errorf("agent.status before any run: %s", got)
	}
	if out, errOut, code := runIn(t, root, "serve", "enable", "curate"); code != 0 || !strings.Contains(out, "curate enabled for t") {
		t.Fatalf("enable: %d %s %s", code, out, errOut)
	}
	if !(&app{}).agentConfig(root).Curate {
		t.Error("enabled, flai serve is not told")
	}
	run := `{"` + root + `":{"curator":{"story":"","harness":"command","command":"run-agent","agent":"curator","role":"curate","pid":42,"started":"2026-10-23T10:00:00Z","ended":"2026-10-23T10:05:00Z","log":"/l/t-curator-20261023T100000Z.log","trigger":"asked","since":"2026-10-01","entries":["D-0007"],"outcome":"worked"}}}`
	dir := string(serve.DirFor(cfgPath))
	if err := os.MkdirAll(dir, 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "agents.json"), []byte(run), 0o600); err != nil {
		t.Fatal(err)
	}
	got := status()
	for _, want := range []string{`"curator":{"story":"","harness":"command","command":"run-agent","agent":"curator"`, `"pid":42`, `"log":"/l/t-curator-20261023T100000Z.log"`,
		`"trigger":"asked"`, `"since":"2026-10-01"`, `"entries":["D-0007"]`, `"outcome":"worked"`, `"analyzer":null`, `"orchestrator":null`, `"running":null`} {
		if !strings.Contains(got, want) {
			t.Errorf("agent.status with a run lacks %s: %s", want, got)
		}
	}
	runIn(t, root, "serve", "disable", "curate")
	if (&app{}).agentConfig(root).Curate {
		t.Error("disabled, flai serve is still told it is on")
	}
}

// S-0104: what each harness a story may name is on this host, and what its
// agent may do, are the operator's, set here and nowhere else.
func TestServeAgentHarness(t *testing.T) {
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	root := tempProject(t)
	out, _, _ := runIn(t, root, "serve", "agent", "harness")
	if !strings.Contains(out, `harness claude-code: "claude" "--permission-mode" "acceptEdits" "--allowedTools" "Bash,mcp__flai"`) {
		t.Errorf("the defaults: %s", out)
	}
	out, errOut, code := runIn(t, root, "serve", "agent", "harness", "claude-code", "--program", "/opt/claude", "--", "--permission-mode", "bypassPermissions")
	if code != 0 || !strings.Contains(out, `harness claude-code: "/opt/claude" "--permission-mode" "bypassPermissions"`) {
		t.Fatalf("set: %d %s %s", code, out, errOut)
	}
	got := (&app{}).agentConfig(root).Harnesses["claude-code"]
	if got.Program != "/opt/claude" || strings.Join(got.Args, " ") != "--permission-mode bypassPermissions" {
		t.Errorf("what flai serve is given: %+v", got)
	}
	// the program alone keeps the arguments; "--" alone means none
	runIn(t, root, "serve", "agent", "harness", "claude-code", "--program", "/usr/bin/claude")
	if got := (&app{}).agentConfig(root).Harnesses["claude-code"]; got.Program != "/usr/bin/claude" || len(got.Args) != 2 {
		t.Errorf("program only: %+v", got)
	}
	runIn(t, root, "serve", "agent", "harness", "claude-code", "--")
	if got := (&app{}).agentConfig(root).Harnesses["claude-code"]; len(got.Args) != 0 {
		t.Errorf("no arguments: %+v", got)
	}
	// setting the command keeps the harnesses, and it is a harness of its own
	runIn(t, root, "serve", "agent", "set", "--", "run-agent", "{story}")
	cfg := (&app{}).agentConfig(root)
	if cfg.Harnesses["command"].Program != "run-agent" || cfg.Harnesses["claude-code"].Program != "/usr/bin/claude" {
		t.Errorf("both: %+v", cfg.Harnesses)
	}
	runIn(t, root, "serve", "agent", "clear")
	if got := (&app{}).agentConfig(root).Harnesses["claude-code"]; got.Program != "/usr/bin/claude" {
		t.Errorf("clear removes only the command: %+v", got)
	}
	out, _, _ = runIn(t, root, "serve", "agent", "harness", "claude-code", "--reset")
	if !strings.Contains(out, `harness claude-code: "claude" "--permission-mode" "acceptEdits"`) {
		t.Errorf("reset: %s", out)
	}
	for _, bad := range [][]string{{"command", "--program", "x"}, {"cursor"}, {"--program", "x"}, {"claude-code", "--reset", "--program", "x"}} {
		if _, _, code := runIn(t, root, append([]string{"serve", "agent", "harness"}, bad...)...); code == 0 {
			t.Errorf("%q was taken", bad)
		}
	}
}

// S-0355 (ADR-0131): a harness's guard none and deny_protected are set and
// cleared on the host, the command's among them, each flag changing only its
// own field; --reset clears both; the command's program stays flai serve agent
// set's; and a guard other than none is refused, naming what sets it.
func TestServeAgentHarnessGuardAndDenyProtected(t *testing.T) {
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	root := tempProject(t)
	harnesses := func() map[string]config.HarnessHost {
		cfg, _, err := (&app{}).loadConfig()
		if err != nil {
			t.Fatal(err)
		}
		return cfg.Agent.Harnesses
	}
	runIn(t, root, "serve", "agent", "set", "--", "run-agent", "{story}")
	out, errOut, code := runIn(t, root, "serve", "agent", "harness", "command", "--guard", "none")
	if code != 0 || !strings.Contains(out, "harness command: the command above\n  guard: none, so a story's sub-agents run without flai guard") {
		t.Fatalf("guard none on the command: %d %s %s", code, out, errOut)
	}
	out, _, _ = runIn(t, root, "serve", "agent", "harness", "command", "--deny-protected")
	if got := harnesses()["command"]; got != (config.HarnessHost{Guard: config.GuardNone, DenyProtected: true}) {
		t.Errorf("--deny-protected keeps the guard: %+v", got)
	}
	if !strings.Contains(out, "  deny_protected: your arguments deny the protected paths") {
		t.Errorf("show says deny_protected: %s", out)
	}
	if got := (&app{}).agentConfig(root).Harnesses["command"]; got.Program != "run-agent" || strings.Join(got.Args, " ") != "{story}" {
		t.Errorf("the command's program is still flai serve agent set's: %+v", got)
	}
	out, _, _ = runIn(t, root, "serve", "agent", "show", "--json")
	var shown struct {
		Harnesses map[string]struct {
			Program       string
			Args          []string
			Guard         string `json:"guard"`
			DenyProtected bool   `json:"deny_protected"`
		} `json:"harnesses"`
	}
	if err := json.Unmarshal([]byte(out), &shown); err != nil {
		t.Fatalf("show --json: %v %s", err, out)
	}
	if got := shown.Harnesses["command"]; got.Program != "run-agent" || strings.Join(got.Args, " ") != "{story}" || got.Guard != "none" || !got.DenyProtected {
		t.Errorf("show --json gives the command's guard and deny_protected: %s", out)
	}
	runIn(t, root, "serve", "agent", "harness", "command", "--guard", "")
	if got := harnesses()["command"]; got != (config.HarnessHost{DenyProtected: true}) {
		t.Errorf(`--guard "" clears only the guard: %+v`, got)
	}
	runIn(t, root, "serve", "agent", "harness", "command", "--deny-protected=false")
	if _, ok := harnesses()["command"]; ok {
		t.Errorf("an entry that lifts nothing is removed: %+v", harnesses())
	}
	// on claude-code, beside its program; --reset clears both
	runIn(t, root, "serve", "agent", "harness", "claude-code", "--program", "/opt/claude")
	runIn(t, root, "serve", "agent", "harness", "claude-code", "--guard", "none", "--deny-protected")
	if got := harnesses()["claude-code"]; got.Program != "/opt/claude" || got.Guard != config.GuardNone || !got.DenyProtected {
		t.Errorf("guard and deny_protected keep the program: %+v", got)
	}
	runIn(t, root, "serve", "agent", "harness", "command", "--guard", "none")
	runIn(t, root, "serve", "agent", "harness", "command", "--reset")
	runIn(t, root, "serve", "agent", "harness", "claude-code", "--reset")
	if got := harnesses(); got != nil {
		t.Errorf("--reset clears guard and deny_protected: %+v", got)
	}
	if _, errOut, code := runIn(t, root, "serve", "agent", "harness", "claude-code", "--guard", "off"); code == 0 || !strings.Contains(errOut, `flai serve agent harness claude-code --guard none, or clear it with --guard`) {
		t.Errorf("a guard other than none is refused, naming what sets it: %d %s", code, errOut)
	}
	for _, bad := range [][]string{{"command", "--program", "x"}, {"command", "--", "x"}, {"command", "--reset", "--guard", "none"}, {"--guard", "none"}, {"--deny-protected"}} {
		if _, _, code := runIn(t, root, append([]string{"serve", "agent", "harness"}, bad...)...); code == 0 {
			t.Errorf("%q was taken", bad)
		}
	}
	if _, errOut, _ := runIn(t, root, "serve", "agent", "harness", "command", "--program", "x"); !strings.Contains(errOut, "flai serve agent set") {
		t.Errorf("the command's program points at flai serve agent set: %s", errOut)
	}
	if got := harnesses(); got != nil {
		t.Errorf("a refusal changes nothing: %+v", got)
	}
}

// S-0355 (ADR-0131): flai serve is told at every look what the operator
// lifts on each harness, the command's among them, and whether the
// project's auto-approve host action is on.
func TestServeAgentConfigCarriesTheLiftsAndAutoApprove(t *testing.T) {
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	root := tempProject(t)
	runIn(t, root, "serve", "agent", "set", "--", "run-agent", "{story}")
	if got := (&app{}).agentConfig(root); len(got.Lifts) != 0 || got.AutoApprove {
		t.Errorf("nothing lifted: %+v %v", got.Lifts, got.AutoApprove)
	}
	runIn(t, root, "serve", "agent", "harness", "command", "--guard", "none")
	runIn(t, root, "serve", "agent", "harness", "claude-code", "--deny-protected")
	got := (&app{}).agentConfig(root).Lifts
	if got["command"] != (serve.HarnessLifts{GuardNone: true}) || got["claude-code"] != (serve.HarnessLifts{DenyProtected: true}) {
		t.Errorf("lifts: %+v", got)
	}
	if _, errOut, code := runIn(t, root, "serve", "enable", "auto-approve"); code != 0 {
		t.Fatal(errOut)
	}
	if !(&app{}).agentConfig(root).AutoApprove {
		t.Error("auto-approve is on for the project")
	}
}

// S-0349 (ADR-0129): a provider's entry is the operator's, set, shown,
// changed field by field, and removed on the host; a new entry is whole, a
// refused value changes nothing, clear keeps the providers, and nothing the
// command prints or writes holds the key, which no flag takes and the
// command never reads.
func TestServeAgentProvider(t *testing.T) {
	cfgPath := filepath.Join(t.TempDir(), "cfg.json")
	t.Setenv("FLAI_CONFIG", cfgPath)
	const key = "sk-never-shown-0123456789"
	t.Setenv("LITELLM_API_KEY", key)
	const masterKey = "sk-master-never-shown-0123456789"
	t.Setenv("LITELLM_MASTER_KEY", masterKey)
	root := tempProject(t)
	var seen strings.Builder
	run := func(args ...string) (string, string, int) {
		t.Helper()
		out, errOut, code := runIn(t, root, append([]string{"serve", "agent", "provider"}, args...)...)
		seen.WriteString(out + errOut)
		return out, errOut, code
	}
	saved := func() map[string]manifest.Provider {
		t.Helper()
		cfg, _, err := config.Load(cfgPath)
		if err != nil {
			t.Fatal(err)
		}
		return cfg.Agent.Providers
	}
	if out, _, _ := run(); !strings.Contains(out, "no provider is set on this host; add one with flai serve agent provider <name> --api <api> --base-url <url> --key-env <variable>") {
		t.Errorf("none set: %s", out)
	}
	if out, _, _ := run("litellm"); !strings.Contains(out, "provider litellm has no entry on this host; add it with flai serve agent provider litellm --api") {
		t.Errorf("a name with no entry: %s", out)
	}
	if _, errOut, code := run("litellm", "--api", "openai-chat"); code == 0 || !strings.Contains(errOut, "has no base_url; set it with flai serve agent provider litellm --base-url") || saved() != nil {
		t.Errorf("a new entry must be whole: %d %s %+v", code, errOut, saved())
	}
	out, errOut, code := run("litellm", "--api", "openai-chat", "--base-url", "http://127.0.0.1:4000", "--key-env", "LITELLM_API_KEY", "--model", "opus=anthropic/claude-opus-5-5", "--model", "sonnet=anthropic/claude-sonnet-5")
	if code != 0 || !strings.Contains(out, "provider litellm: openai-chat at http://127.0.0.1:4000, its key in the variable LITELLM_API_KEY\n  model opus: anthropic/claude-opus-5-5\n  model sonnet: anthropic/claude-sonnet-5\n") {
		t.Fatalf("set: %d %s %s", code, out, errOut)
	}
	if out, _, _ := run("litellm"); !strings.Contains(out, "provider litellm: openai-chat at http://127.0.0.1:4000") {
		t.Errorf("print: %s", out)
	}
	// the flags given change those fields; --model given replaces the aliases
	run("litellm", "--base-url", "https://litellm.example.com", "--model", "haiku=anthropic/claude-haiku-4-5")
	want := manifest.Provider{API: "openai-chat", BaseURL: "https://litellm.example.com", KeyEnv: "LITELLM_API_KEY", Models: map[string]string{"haiku": "anthropic/claude-haiku-4-5"}}
	if got := saved(); !reflect.DeepEqual(got, map[string]manifest.Provider{"litellm": want}) {
		t.Errorf("update: %+v", got)
	}
	run("litellm", "--api", "anthropic-messages")
	want.API = "anthropic-messages"
	if got := saved()["litellm"]; !reflect.DeepEqual(got, want) {
		t.Errorf("--api alone keeps the rest: %+v", got)
	}
	if js, _, _ := run("litellm", "--json"); !strings.Contains(js, `"litellm": {`) || !strings.Contains(js, `"key_env": "LITELLM_API_KEY"`) || strings.Contains(js, "spend_key_env") {
		t.Errorf("--json: %s", js)
	}
	// S-0357: --spend-key-env names the variable of the key that reads the
	// spend log, shown where set, and an empty value clears it
	out, errOut, code = run("litellm", "--spend-key-env", "LITELLM_MASTER_KEY")
	if code != 0 || !strings.Contains(out, "its key in the variable LITELLM_API_KEY\n  spend log read with the key in the variable LITELLM_MASTER_KEY\n  model haiku:") {
		t.Errorf("--spend-key-env: %d %s %s", code, out, errOut)
	}
	if got, w := saved()["litellm"], (manifest.Provider{API: want.API, BaseURL: want.BaseURL, KeyEnv: want.KeyEnv, SpendKeyEnv: "LITELLM_MASTER_KEY", Models: want.Models}); !reflect.DeepEqual(got, w) {
		t.Errorf("--spend-key-env keeps the rest: %+v", got)
	}
	if js, _, _ := run("litellm", "--json"); !strings.Contains(js, `"spend_key_env": "LITELLM_MASTER_KEY"`) {
		t.Errorf("--json with spend_key_env: %s", js)
	}
	if out, _, _ := run("litellm", "--spend-key-env", ""); strings.Contains(out, "spend log") || !reflect.DeepEqual(saved()["litellm"], want) {
		t.Errorf("--spend-key-env \"\" clears it: %s %+v", out, saved()["litellm"])
	}
	for _, bad := range []struct {
		args []string
		want string
	}{
		{[]string{"litellm", "--api", "bedrock"}, "is not an API flai speaks"},
		{[]string{"litellm", "--base-url", "litellm.example.com"}, "is not an absolute http or https URL"},
		{[]string{"litellm", "--key-env", "litellm-key"}, "is not the name of an environment variable"},
		{[]string{"litellm", "--spend-key-env", masterKey}, "spend_key_env is not the name of an environment variable"},
		{[]string{"litellm", "--model", "opus"}, "is not alias=name"},
		{[]string{"--api", "openai-chat"}, "name the provider"},
		{[]string{"litellm", "--reset", "--api", "openai-chat"}, "--reset takes no other flag"},
		{[]string{"lite llm", "--api", "openai-chat", "--base-url", "https://x.example.com", "--key-env", "X"}, "is not a provider name"},
	} {
		if _, errOut, code := run(bad.args...); code == 0 || !strings.Contains(errOut, bad.want) {
			t.Errorf("%q is refused, saying why: %d %s", bad.args, code, errOut)
		}
	}
	if got := saved()["litellm"]; !reflect.DeepEqual(got, want) {
		t.Errorf("a refused value changes nothing: %+v", got)
	}
	runIn(t, root, "serve", "agent", "clear")
	if got := saved(); len(got) != 1 {
		t.Errorf("clear keeps the providers: %+v", got)
	}
	if out, _, _ := run("litellm", "--reset"); !strings.Contains(out, "provider litellm has no entry on this host") || saved() != nil {
		t.Errorf("reset: %s %+v", out, saved())
	}
	data, err := os.ReadFile(cfgPath)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(seen.String()+string(data), key) || strings.Contains(seen.String()+string(data), masterKey) {
		t.Errorf("the key is in what flai printed or wrote:\n%s\n%s", seen.String(), data)
	}
}

// S-0392: flai serve forge shows the variable each forge's token is read
// from, GITHUB_TOKEN by default, and whether it is set here; --token-env
// names another and --reset goes back; an unknown forge and a name that is
// not a variable's are refused, changing nothing; and nothing the command
// prints or writes holds the token, which no flag takes.
func TestServeForge(t *testing.T) {
	cfgPath := filepath.Join(t.TempDir(), "cfg.json")
	t.Setenv("FLAI_CONFIG", cfgPath)
	const token = "ghp_never-shown-0123456789"
	t.Setenv("GITHUB_TOKEN", "")
	t.Setenv("FLAI_GH_TOKEN", token)
	root := tempProject(t)
	var seen strings.Builder
	run := func(args ...string) (string, string, int) {
		t.Helper()
		out, errOut, code := runIn(t, root, append([]string{"serve", "forge"}, args...)...)
		seen.WriteString(out + errOut)
		return out, errOut, code
	}
	saved := func() map[string]config.ForgeHost {
		t.Helper()
		cfg, _, err := config.Load(cfgPath)
		if err != nil {
			t.Fatal(err)
		}
		return cfg.Forges
	}
	if out, _, code := run(); code != 0 || out != "forge github: its token in the variable GITHUB_TOKEN (default), not set in this environment\n" {
		t.Errorf("unset: %d %q", code, out)
	}
	out, errOut, code := run("github", "--token-env", "FLAI_GH_TOKEN")
	if code != 0 || out != "forge github: its token in the variable FLAI_GH_TOKEN, set in this environment\n" {
		t.Fatalf("set: %d %s %s", code, out, errOut)
	}
	if got := saved(); !reflect.DeepEqual(got, map[string]config.ForgeHost{"github": {TokenEnv: "FLAI_GH_TOKEN"}}) {
		t.Errorf("saved: %+v", got)
	}
	if out, _, _ := run("github"); !strings.Contains(out, "forge github: its token in the variable FLAI_GH_TOKEN, set") {
		t.Errorf("show one: %s", out)
	}
	if js, _, _ := run("--json"); !strings.Contains(js, `"github": {`) || !strings.Contains(js, `"token_env": "FLAI_GH_TOKEN"`) || !strings.Contains(js, `"token_set": true`) || !strings.Contains(js, `"default": false`) {
		t.Errorf("--json: %s", js)
	}
	for _, bad := range []struct {
		args []string
		want string
	}{
		{[]string{"gitlab", "--token-env", "GITLAB_TOKEN"}, "gitlab\\\"; it calls github"},
		{[]string{"gitlab"}, "flai calls no forge"},
		{[]string{"github", "--token-env", "gh-token"}, "forges.github.token_env is not the name of an environment variable"},
		{[]string{"github", "--token-env", token}, "never the token"},
		{[]string{"github", "--token-env", ""}, "--token-env needs the name of a variable, such as GITHUB_TOKEN"},
		{[]string{"--token-env", "X"}, "name the forge"},
		{[]string{"github", "--reset", "--token-env", "X"}, "--reset takes no --token-env"},
	} {
		if _, errOut, code := run(bad.args...); code == 0 || !strings.Contains(errOut, bad.want) {
			t.Errorf("%q is refused, saying why: %d %s", bad.args, code, errOut)
		}
	}
	if got := saved()["github"]; got.TokenEnv != "FLAI_GH_TOKEN" {
		t.Errorf("a refused value changes nothing: %+v", got)
	}
	// S-0393: flai serve is told the variable it keeps from every agent
	if got := (&app{}).agentConfig(root).ForgeTokenEnvs; len(got) != 1 || got[0] != "FLAI_GH_TOKEN" {
		t.Errorf("the forge tokens kept from the agents: %q, want FLAI_GH_TOKEN", got)
	}
	data, err := os.ReadFile(cfgPath)
	if err != nil {
		t.Fatal(err)
	}
	if out, _, _ := run("github", "--reset"); out != "forge github: its token in the variable GITHUB_TOKEN (default), not set in this environment\n" || saved() != nil {
		t.Errorf("reset: %q %+v", out, saved())
	}
	if got := (&app{}).agentConfig(root).ForgeTokenEnvs; len(got) != 1 || got[0] != "GITHUB_TOKEN" {
		t.Errorf("the forge tokens kept from the agents once reset: %q, want GITHUB_TOKEN", got)
	}
	if strings.Contains(seen.String()+string(data), token) {
		t.Errorf("the token is in what flai printed or wrote:\n%s\n%s", seen.String(), data)
	}
}

// S-0294 (ADR-0108): --auto-restarts sets how many times flai serve restarts
// a story's agent on its own, alone without touching the command, the name,
// or the harnesses; show says it and whether it is the default; a negative
// count is refused; and clear keeps it.
func TestServeAgentAutoRestarts(t *testing.T) {
	cfgPath := filepath.Join(t.TempDir(), "cfg.json")
	t.Setenv("FLAI_CONFIG", cfgPath)
	root := tempProject(t)
	if out, _, _ := runIn(t, root, "serve", "agent", "show"); !strings.Contains(out, "automatic restarts: 2 (default)") {
		t.Errorf("unset, show says the default: %s", out)
	}
	runIn(t, root, "serve", "agent", "set", "--name", "builder", "--", "run-agent", "{story}")
	runIn(t, root, "serve", "agent", "harness", "claude-code", "--program", "/opt/claude")
	saved := func() config.AgentStart {
		t.Helper()
		cfg, _, err := config.Load(cfgPath)
		if err != nil {
			t.Fatal(err)
		}
		return cfg.Agent
	}
	for _, c := range []struct {
		n, said string
		want    int
	}{{"0", "automatic restarts: 0 (off)", 0}, {"3", "automatic restarts: 3\n", 3}} {
		out, errOut, code := runIn(t, root, "serve", "agent", "set", "--auto-restarts", c.n)
		if code != 0 || !strings.Contains(out, c.said) {
			t.Fatalf("--auto-restarts %s: %d %s %s", c.n, code, out, errOut)
		}
		got := saved()
		if got.AutoRestarts == nil || *got.AutoRestarts != c.want {
			t.Errorf("--auto-restarts %s is written: %v", c.n, got.AutoRestarts)
		}
		if strings.Join(got.Command, " ") != "run-agent {story}" || got.Name != "builder" || got.Harnesses["claude-code"].Program != "/opt/claude" {
			t.Errorf("--auto-restarts %s alone keeps the command, the name, and the harnesses: %+v", c.n, got)
		}
	}
	if js, _, _ := runIn(t, root, "serve", "agent", "show", "--json"); !strings.Contains(js, `"auto_restarts": 3`) {
		t.Errorf("show --json: %s", js)
	}
	_, errOut, code := runIn(t, root, "serve", "agent", "set", "--auto-restarts", "-1")
	if code == 0 || !strings.Contains(errOut, "0 or more") {
		t.Errorf("a negative count is refused, saying what to give: %d %s", code, errOut)
	}
	if got := saved(); got.AutoRestarts == nil || *got.AutoRestarts != 3 {
		t.Errorf("a refused count changes nothing: %v", got.AutoRestarts)
	}
	runIn(t, root, "serve", "agent", "clear")
	if got := saved(); len(got.Command) != 0 || got.AutoRestarts == nil || *got.AutoRestarts != 3 {
		t.Errorf("clear removes the command and keeps the limit: %+v", got)
	}
}

// S-0082: the checks commands are the operator's, named argument lists with
// no default, several at once, managed on the host and nowhere else.
func TestServeChecksCommands(t *testing.T) {
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	root := tempProject(t)
	out, _, _ := runIn(t, root, "serve", "checks")
	if !strings.Contains(out, "no command is named here; the manifest's own checks: is used instead") ||
		!strings.Contains(out, "timeout: 15 minutes") || !strings.Contains(out, "off for this project") {
		t.Errorf("by default: %s", out)
	}
	if out, _, _ := runIn(t, root, "serve", "actions"); !strings.Contains(out, "checks: off everywhere") || !strings.Contains(out, "in a story's worktree") {
		t.Errorf("actions names it and what it means: %s", out)
	}
	if _, _, code := runIn(t, root, "serve", "checks", "set", "--", "go", "test", "./..."); code == 0 {
		t.Error("set needs --name")
	}
	if _, _, code := runIn(t, root, "serve", "checks", "set", "--name", "flai"); code == 0 {
		t.Error("set needs a program")
	}
	out, errOut, code := runIn(t, root, "serve", "checks", "set", "--name", "flai", "--", "scripts/flai-test.sh")
	if code != 0 || !strings.Contains(out, `flai: "scripts/flai-test.sh"`) {
		t.Fatalf("set: %d %s %s", code, out, errOut)
	}
	out, _, code = runIn(t, root, "serve", "checks", "set", "--name", "flaiover", "--", "bash", "-c", "cd flaiover; pnpm test")
	if code != 0 || !strings.Contains(out, `flai: "scripts/flai-test.sh"`) || !strings.Contains(out, `flaiover: "bash" "-c" "cd flaiover; pnpm test"`) {
		t.Fatalf("a second name adds, does not replace: %d %s", code, out)
	}
	out, _, code = runIn(t, root, "serve", "checks", "set", "--name", "flai", "--", "scripts/flai-test.sh", "--fast")
	if code != 0 || !strings.Contains(out, `flai: "scripts/flai-test.sh" "--fast"`) || strings.Count(out, "flai:") != 1 {
		t.Fatalf("the same name replaces: %d %s", code, out)
	}
	cfg, _ := os.ReadFile(os.Getenv("FLAI_CONFIG"))
	if !strings.Contains(string(cfg), `"cd flaiover; pnpm test"`) || !strings.Contains(string(cfg), `"name": "flaiover"`) {
		t.Errorf("kept as an argument list, as written: %s", cfg)
	}
	if _, errOut, code := runIn(t, root, "config", "set", "checks.commands", "rm -rf /"); code == 0 {
		t.Errorf("flai config set does not reach it: %s", errOut)
	}
	out, _, code = runIn(t, root, "serve", "checks", "timeout", "20")
	if code != 0 || !strings.Contains(out, "timeout: 20 minutes") {
		t.Fatalf("timeout: %d %s", code, out)
	}
	if _, _, code := runIn(t, root, "serve", "checks", "timeout", "0"); code == 0 {
		t.Error("timeout must be positive")
	}
	runIn(t, root, "serve", "enable", "checks")
	js, _, _ := runIn(t, root, "serve", "checks", "show", "--json")
	if !strings.Contains(js, `"enabled_here": true`) || !strings.Contains(js, `"timeout_minutes": 20`) || !strings.Contains(js, `"flai"`) || !strings.Contains(js, `"flaiover"`) {
		t.Errorf("show --json: %s", js)
	}
	if out, _, _ := runIn(t, root, "serve", "checks", "clear", "flai"); strings.Contains(out, `flai:`) || !strings.Contains(out, `flaiover:`) {
		t.Errorf("clear one name: %s", out)
	}
	if out, _, _ := runIn(t, root, "serve", "checks", "clear"); !strings.Contains(out, "no command is named here") {
		t.Errorf("clear with no name clears all: %s", out)
	}
}

// S-0438: components is a host action the operator enables as checks is, and
// components.idle a host setting flai serve components idle keeps, which
// flai serve's supervisor of the dev servers reads.
func TestServeComponentsActionAndIdle(t *testing.T) {
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	root := tempProject(t)
	out, _, _ := runIn(t, root, "serve", "components")
	if !strings.Contains(out, "idle: 30m (default)") || !strings.Contains(out, "max_servers: 3 (default): at most this many dev servers run at once per project") ||
		!strings.Contains(out, "components action is off for this project; flai serve enable components") {
		t.Errorf("by default: %s", out)
	}
	if out, _, _ := runIn(t, root, "serve", "actions"); !strings.Contains(out, "components: off everywhere") || !strings.Contains(out, "Vite dev server") {
		t.Errorf("actions names it and what it means: %s", out)
	}
	for _, bad := range []string{"0", "-5m", "soon", "30", ""} {
		if _, errOut, code := runIn(t, root, "serve", "components", "idle", "--", bad); code == 0 || !strings.Contains(errOut, "positive duration") {
			t.Errorf("idle %q is refused: %d %s", bad, code, errOut)
		}
	}
	out, errOut, code := runIn(t, root, "serve", "components", "idle", "45m")
	if code != 0 || !strings.Contains(out, "idle: 45m: a dev server") {
		t.Fatalf("idle 45m: %d %s %s", code, out, errOut)
	}
	// S-0442: components.max_servers, a positive whole number, 3 by default
	for _, bad := range []string{"0", "-1", "two", "2.5", ""} {
		if _, errOut, code := runIn(t, root, "serve", "components", "max-servers", "--", bad); code == 0 || !strings.Contains(errOut, "positive whole number") {
			t.Errorf("max-servers %q is refused: %d %s", bad, code, errOut)
		}
	}
	out, errOut, code = runIn(t, root, "serve", "components", "max-servers", "5")
	if code != 0 || !strings.Contains(out, "max_servers: 5: at most") || !strings.Contains(out, "idle: 45m: a dev server") {
		t.Fatalf("max-servers 5: %d %s %s", code, out, errOut)
	}
	cfg, _ := os.ReadFile(os.Getenv("FLAI_CONFIG"))
	if !strings.Contains(string(cfg), `"components"`) || !strings.Contains(string(cfg), `"idle": "45m"`) || !strings.Contains(string(cfg), `"max_servers": 5`) {
		t.Errorf("kept in the host's configuration: %s", cfg)
	}
	if out, _, code := runIn(t, root, "serve", "enable", "components"); code != 0 || !strings.Contains(out, "components enabled for") {
		t.Fatalf("enable components: %d %s", code, out)
	}
	js, _, _ := runIn(t, root, "serve", "components", "show", "--json")
	var shown struct {
		Idle        string `json:"idle"`
		IdleSeconds int    `json:"idle_seconds"`
		IdleSet     bool   `json:"idle_set"`
		MaxServers  int    `json:"max_servers"`
		MaxSet      bool   `json:"max_servers_set"`
		EnabledHere bool   `json:"enabled_here"`
	}
	if err := json.Unmarshal([]byte(js), &shown); err != nil || shown.Idle != "45m" || shown.IdleSeconds != 2700 || !shown.IdleSet || shown.MaxServers != 5 || !shown.MaxSet || !shown.EnabledHere {
		t.Errorf("show --json: %v %s", err, js)
	}
	if st, _, _ := runIn(t, root, "hostapi", "components.status"); !strings.Contains(st, `"enabled":true`) || !strings.Contains(st, `"enable":"flai serve enable components"`) {
		t.Errorf("the dashboard is told the action is on: %s", st)
	}

	// flai serve's own host has the supervisor, on components.idle and
	// components.max_servers as set, and stops it when asked
	var sink bytes.Buffer
	a := &app{out: &sink, errOut: &sink, cwd: root}
	host, stop := a.serveHost(t.Context())
	if host.Components == nil || host.Components.Idle() != 45*time.Minute || !host.Enabled(hostapi.ActionComponents, root) {
		t.Errorf("serveHost: %+v", host)
	}
	if most, ok := host.Components.(interface{ MaxServers() int }); !ok || most.MaxServers() != 5 {
		t.Errorf("serveHost's supervisor reads components.max_servers: %+v", host.Components)
	}
	stopped := make(chan struct{})
	go func() { stop(); close(stopped) }()
	select {
	case <-stopped:
	case <-time.After(30 * time.Second):
		t.Fatal("serveHost's stop did not return")
	}
	if _, err := host.Components.Start(t.Context(), root, "t", manifest.Project{Name: "ui", Path: "ui", Kind: manifest.ProjectKindSvelteKit}, "", ""); err == nil || !strings.Contains(err.Error(), "flai serve is stopping") {
		t.Errorf("once stopped, it starts no dev server: %v", err)
	}
	if out, _, _ := runIn(t, root, "serve", "disable", "components"); !strings.Contains(out, "components disabled for") {
		t.Errorf("disable components: %s", out)
	}
}

// S-0082: the host's own named commands are used when there are any; the
// manifest's checks: only when the host names none. Never a merge.
func TestChecksConfigResolvesHostOverManifest(t *testing.T) {
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	root := tempProject(t)
	runIn(t, root, "serve", "checks", "show") // creates the config file, as flai itself would
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	a := &app{}

	// nothing named anywhere
	cc, err := a.checksConfig(repo)
	if err != nil || len(cc.Commands) != 0 || cc.Enabled || cc.Timeout != serve.DefaultChecksTimeout {
		t.Fatalf("nothing named: %+v %v", cc, err)
	}

	// the manifest names checks:, host config names none: the manifest's are used
	manifestYAML := "version: 1\nname: t\nkey: t\nlayout:\n  design: design\n  docs: docs\n  wip: wip\n" +
		"checks:\n  - name: flai\n    command: [scripts/flai-test.sh]\n"
	if err := os.WriteFile(filepath.Join(root, "system-flow.yaml"), []byte(manifestYAML), 0o644); err != nil {
		t.Fatal(err)
	}
	repo, err = workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	cc, err = a.checksConfig(repo)
	if err != nil || len(cc.Commands) != 1 || cc.Commands[0].Name != "flai" {
		t.Fatalf("falls back to the manifest: %+v %v", cc, err)
	}

	// the host names a different one: the host's is used, not a merge
	runIn(t, root, "serve", "checks", "set", "--name", "flaiover", "--", "true")
	cc, err = a.checksConfig(repo)
	if err != nil || len(cc.Commands) != 1 || cc.Commands[0].Name != "flaiover" {
		t.Fatalf("the host's own commands, not merged with the manifest's: %+v %v", cc, err)
	}
	if cc.Enabled {
		t.Error("still off until enabled")
	}
	runIn(t, root, "serve", "enable", "checks")
	if cc, err := a.checksConfig(repo); err != nil || !cc.Enabled {
		t.Errorf("enabled: %+v %v", cc, err)
	}
}

// S-0105: the settings host action, the agent's name without a command, and
// what the dashboard's settings page is told.
func TestServeSettingsAreShownToTheDashboard(t *testing.T) {
	cfgPath := filepath.Join(t.TempDir(), "cfg.json")
	t.Setenv("FLAI_CONFIG", cfgPath)
	root := tempProject(t)
	// S-0122: the served projects, without the dashboard's address or the credential's path
	if err := serve.DirFor(cfgPath).Register(serve.Entry{Key: "harbour", Name: "Harbour", Root: "/p/harbour", URL: "http://127.0.0.1:4242", KeyFile: "/p/harbour/k"}); err != nil {
		t.Fatal(err)
	}
	if out, _, _ := runIn(t, root, "serve", "actions"); !strings.Contains(out, "settings: off everywhere") || !strings.Contains(out, "only a shell turns this off") {
		t.Errorf("actions names settings: %s", out)
	}
	// --attended-minutes is retired (S-0116): still accepted, it says so and does nothing
	if _, errOut, code := runIn(t, root, "serve", "agent", "set", "--name", "builder", "--attended-minutes", "9"); code != 0 || !strings.Contains(errOut, "attended-minutes") || !strings.Contains(errOut, "does nothing") {
		t.Fatalf("name without a command: %d %s", code, errOut)
	}
	cfg := (&app{}).agentConfig(root)
	if cfg.Name != "builder" || len(cfg.Command) != 0 {
		t.Errorf("name: %+v", cfg)
	}
	if loaded, _, _ := (&app{}).loadConfig(); loaded.Agent.AttendedMinutes != 0 {
		t.Errorf("the retired minutes were stored: %+v", loaded.Agent)
	}
	runIn(t, root, "serve", "enable", "settings")
	runIn(t, root, "serve", "checks", "set", "--name", "unit", "--", "go", "test", "./...")
	// S-0488: a temp directory's name may hold 4242, so the import root holds
	// it always, and only the registered entry's URL is refused below
	importRoot := filepath.Join(t.TempDir(), "harbour-4242")
	if err := os.Mkdir(importRoot, 0o755); err != nil {
		t.Fatal(err)
	}
	runIn(t, root, "serve", "import", "add", importRoot)
	runIn(t, root, "agent", "set", "--harness", "claude-code", "--model", "claude-haiku-4-5")
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	got, _ := json.Marshal((&app{}).hostSettings(mainRootOf(repo)))
	for _, want := range []string{`"here":true,"means":"change this project's host settings`, `"here":false,"means":"publish when you press Publish`, `"name":"builder"`,
		`"claude-code":{"args":["--permission-mode","acceptEdits","--allowedTools","Bash,mcp__flai"],"ask_capable":true,"deny_protected":false,"guard":"","guard_capable":true,"program":"claude","set":false}`,
		`"command":{"args":[],"ask_capable":false,"deny_protected":false,"guard":"","guard_capable":false,"program":"","set":false}`,
		`{"name":"unit","command":["go","test","./..."]}`, `"timeout_minutes":15`, `"import_roots":["` + importRoot + `"]`, `"default_agent":{"harness":"claude-code","model":"claude-haiku-4-5"}`, `"mcp":{"running":false}`,
		`"projects":{"running":false,"served":[{"key":"harbour","name":"Harbour","root":"/p/harbour","from":"registry","state":"unavailable","reason":"`, `"unserved":[]}`} {
		if !strings.Contains(string(got), want) {
			t.Errorf("settings lack %s:\n%s", want, got)
		}
	}
	if strings.Contains(string(got), `"token"`) || strings.Contains(string(got), "127.0.0.1:4242") || strings.Contains(string(got), "/p/harbour/k") {
		t.Errorf("no token in the settings: %s", got)
	}
}

// S-0355, ADR-0131: settings.get gives each harness its guard and
// deny_protected; an entry holding only those keeps set false, so the page
// still says the program and arguments are the defaults; and the command is
// there with flai serve agent set's program and arguments.
func TestServeSettingsShowEachHarnessGuardAndDenyProtected(t *testing.T) {
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	root := tempProject(t)
	runIn(t, root, "serve", "agent", "set", "--", "run-agent", "{story}")
	for _, args := range [][]string{{"claude-code", "--guard", "none"}, {"command", "--guard", "none", "--deny-protected"}} {
		if _, errOut, code := runIn(t, root, append([]string{"serve", "agent", "harness"}, args...)...); code != 0 {
			t.Fatalf("%q: %s", args, errOut)
		}
	}
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	got, _ := json.Marshal((&app{}).hostSettings(mainRootOf(repo)))
	for _, want := range []string{
		`"claude-code":{"args":["--permission-mode","acceptEdits","--allowedTools","Bash,mcp__flai"],"ask_capable":true,"deny_protected":false,"guard":"none","guard_capable":true,"program":"claude","set":false}`,
		`"command":{"args":["{story}"],"ask_capable":false,"deny_protected":true,"guard":"none","guard_capable":false,"program":"run-agent","set":false}`,
	} {
		if !strings.Contains(string(got), want) {
			t.Errorf("settings lack %s:\n%s", want, got)
		}
	}
	runIn(t, root, "serve", "agent", "harness", "claude-code", "--program", "/opt/claude")
	got, _ = json.Marshal((&app{}).hostSettings(mainRootOf(repo)))
	if want := `"guard":"none","guard_capable":true,"program":"/opt/claude","set":true}`; !strings.Contains(string(got), want) {
		t.Errorf("a program beside the guard is set: want %s in\n%s", want, got)
	}
}

// S-0359, ADR-0129: settings.get's host carries the host's providers, each
// entry's fields and whether the variables its key_env and spend_key_env
// (S-0357) name are set and not empty in flai serve's environment, never a
// variable's value.
func TestServeSettingsShowTheProvidersNeverTheirKeys(t *testing.T) {
	cfgPath := filepath.Join(t.TempDir(), "cfg.json")
	t.Setenv("FLAI_CONFIG", cfgPath)
	const key = "sk-never-shown-0123456789"
	t.Setenv("LITELLM_API_KEY", key)
	t.Setenv("OPENROUTER_API_KEY", "")
	const masterKey = "sk-master-never-shown-0123456789"
	t.Setenv("LITELLM_MASTER_KEY", masterKey)
	t.Setenv("FLAI_TEST_UNSET_GATEWAY_KEY", "") // restored after the test, unset in it
	if err := os.Unsetenv("FLAI_TEST_UNSET_GATEWAY_KEY"); err != nil {
		t.Fatal(err)
	}
	root := tempProject(t)
	runIn(t, root, "serve", "agent", "provider") // creates the config, which an app with no logger cannot
	agent := func() map[string]any {
		t.Helper()
		repo, err := workitem.Open(root)
		if err != nil {
			t.Fatal(err)
		}
		got, _ := json.Marshal((&app{}).hostSettings(mainRootOf(repo)))
		if strings.Contains(string(got), key) || strings.Contains(string(got), masterKey) {
			t.Errorf("a key is in the settings: %s", got)
		}
		var out struct {
			Agent map[string]any `json:"agent"`
		}
		if err := json.Unmarshal(got, &out); err != nil {
			t.Fatal(err)
		}
		return out.Agent
	}
	if got, _ := json.Marshal(agent()["providers"]); string(got) != `{}` {
		t.Errorf("no provider set is an empty object: %s", got)
	}
	for _, args := range [][]string{
		{"litellm", "--api", "anthropic-messages", "--base-url", "http://127.0.0.1:4000", "--key-env", "LITELLM_API_KEY", "--spend-key-env", "LITELLM_MASTER_KEY", "--model", "opus=anthropic/claude-opus-5-5"},
		{"openrouter", "--api", "openai-chat", "--base-url", "https://openrouter.ai/api/v1", "--key-env", "OPENROUTER_API_KEY"},
		{"gateway", "--api", "openai-responses", "--base-url", "https://gw.example.com", "--key-env", "FLAI_TEST_UNSET_GATEWAY_KEY", "--spend-key-env", "FLAI_TEST_UNSET_GATEWAY_KEY"},
	} {
		if _, errOut, code := runIn(t, root, append([]string{"serve", "agent", "provider"}, args...)...); code != 0 {
			t.Fatalf("set %s: %s", args[0], errOut)
		}
	}
	got, _ := json.Marshal(agent()["providers"])
	want := `{"gateway":{"api":"openai-responses","base_url":"https://gw.example.com","key_env":"FLAI_TEST_UNSET_GATEWAY_KEY","key_set":false,"models":{},"spend_key_env":"FLAI_TEST_UNSET_GATEWAY_KEY","spend_key_set":false},` +
		`"litellm":{"api":"anthropic-messages","base_url":"http://127.0.0.1:4000","key_env":"LITELLM_API_KEY","key_set":true,"models":{"opus":"anthropic/claude-opus-5-5"},"spend_key_env":"LITELLM_MASTER_KEY","spend_key_set":true},` +
		`"openrouter":{"api":"openai-chat","base_url":"https://openrouter.ai/api/v1","key_env":"OPENROUTER_API_KEY","key_set":false,"models":{},"spend_key_env":"","spend_key_set":false}}`
	if string(got) != want {
		t.Errorf("providers:\n%s\nwant\n%s", got, want)
	}
}

// S-0295, ADR-0096: settings.get carries the shared paths as the manifest
// has them now; settings.shared changes them through flai shared under the
// settings action, its refusal a rule with flai's reason; and
// settings.shared_check answers what flai shared check does, for anyone.
func TestSettingsSharedThroughTheHostAPI(t *testing.T) {
	root := sharedProject(t)
	inProcess(t)
	hostapiCall := func(method, params string) (string, int) {
		t.Helper()
		out, _, code := runIn(t, root, "hostapi", method, params)
		return out, code
	}
	shared := func() string {
		t.Helper()
		out, code := hostapiCall("settings.get", `{}`)
		var got struct {
			Host struct {
				Shared json.RawMessage `json:"shared"`
			} `json:"host"`
		}
		if err := json.Unmarshal([]byte(out), &got); code != 0 || err != nil {
			t.Fatalf("settings.get: %d %v %s", code, err, out)
		}
		return string(got.Host.Shared)
	}
	if got := shared(); got != `[]` {
		t.Errorf("settings.get with no shared paths: %s, want []", got)
	}

	before := readManifestShared(t, root)
	if out, code := hostapiCall("settings.shared", `{"action":"add","pattern":"design/adrs","request_id":"req-00000051"}`); code == 0 || !strings.Contains(out, "flai serve enable settings") || !strings.Contains(out, "-32012") {
		t.Errorf("with the settings action off: %d %s", code, out)
	}
	if got := readManifestShared(t, root); !reflect.DeepEqual(got, before) {
		t.Errorf("a refused change wrote the manifest: %v", got)
	}

	runIn(t, root, "serve", "enable", "settings")
	for i, c := range []struct{ params, want string }{
		{`{"action":"add","pattern":"design/adrs"`, `"data":{"added":["design/adrs"],"shared":["design/adrs"]}`},
		{`{"action":"add","pattern":"docs/users/*.md"`, `"data":{"added":["docs/users/*.md"],"shared":["design/adrs","docs/users/*.md"]}`},
		{`{"action":"add","pattern":"template"`, `"data":{"added":["template"],"shared":["design/adrs","docs/users/*.md","template"]}`},
		{`{"action":"remove","pattern":"template"`, `"data":{"removed":["template"],"shared":["design/adrs","docs/users/*.md"]}`},
	} {
		out, code := hostapiCall("settings.shared", c.params+`,"request_id":"req-0000006`+strconv.Itoa(i)+`"}`)
		if code != 0 || !strings.Contains(out, c.want) {
			t.Errorf("settings.shared %s: %d %s, want %s", c.params, code, out, c.want)
		}
	}
	if got := shared(); got != `["design/adrs","docs/users/*.md"]` {
		t.Errorf("settings.get after the changes: %s", got)
	}

	for i, c := range []struct{ params, why string }{
		{`{"action":"add","pattern":"design/adrs"`, "in the list already"},
		{`{"action":"add","pattern":"/design"`, "absolute"},
		{`{"action":"remove","pattern":"docs"`, "not in the list"},
	} {
		out, code := hostapiCall("settings.shared", c.params+`,"request_id":"req-0000007`+strconv.Itoa(i)+`"}`)
		if code == 0 || !strings.Contains(out, "-32011") || !strings.Contains(out, c.why) {
			t.Errorf("settings.shared %s: %d %s, want a rule saying %q", c.params, code, out, c.why)
		}
	}
	if got := readManifestShared(t, root); !reflect.DeepEqual(got, []string{"design/adrs", "docs/users/*.md"}) {
		t.Errorf("the refusals changed the manifest: %v", got)
	}

	// a read for anyone: with the settings action off again it still answers
	runIn(t, root, "serve", "disable", "settings")
	out, code := hostapiCall("settings.shared_check", `{"paths":["design/adrs/0001-x.md","docs/users","cli"],"story":"S-0001"}`)
	var got struct {
		Data []sharedCheck `json:"data"`
	}
	if err := json.Unmarshal([]byte(out), &got); code != 0 || err != nil {
		t.Fatalf("settings.shared_check: %d %v %s", code, err, out)
	}
	want := []sharedCheck{
		{Entry: "design/adrs/0001-x.md", Path: "design/adrs/0001-x.md", Shared: true, Pattern: "design/adrs"},
		{Entry: "docs/users", Path: "docs/users"},
		{Entry: "cli", Path: "flai"},
		{Entry: "docs/users/flai.md", Path: "docs/users/flai.md", Shared: true, Pattern: "docs/users/*.md", Story: "S-0001"},
		{Entry: "flai", Path: "flai", Story: "S-0001"},
		{Entry: "design/adrs", Path: "design/adrs", Shared: true, Pattern: "design/adrs", Story: "S-0001"},
		{Entry: "template/CHANGELOG.md", Path: "template/CHANGELOG.md", Story: "S-0001"},
	}
	if !reflect.DeepEqual(got.Data, want) {
		t.Errorf("settings.shared_check =\n%+v\nwant\n%+v", got.Data, want)
	}
	if out, code := hostapiCall("settings.shared_check", `{"story":"S-0099"}`); code == 0 || !strings.Contains(out, "-32004") {
		t.Errorf("a story that is not there: %d %s", code, out)
	}

	js, _, _ := runIn(t, root, "serve", "journal")
	if !strings.Contains(js, "flai shared add --autocommit -- design/adrs") || !strings.Contains(js, "flai shared remove --autocommit -- template") || strings.Contains(js, "shared_check") {
		t.Errorf("the journal names each change and no check:\n%s", js)
	}
}

// S-0211, ADR-0084: the settings page is told when flai serve plans again on
// its own: whether the plan host action is on, the replan policy and whether
// the manifest sets it, and the schedule as written with its next run; a
// value flai cannot read comes with its error.
func TestSettingsShowThePlanningTriggers(t *testing.T) {
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	root := tempProject(t)
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	a := &app{errOut: &bytes.Buffer{}, clock: func() time.Time { return time.Date(2026, 10, 4, 9, 30, 0, 0, time.UTC) }}
	planning := func() string {
		t.Helper()
		got, _ := json.Marshal(a.hostSettings(mainRootOf(repo)).(map[string]any)["planning"])
		return string(got)
	}
	if got, want := planning(), `{"plan":false,"replan":"deterministic","replan_set":false,"schedule":""}`; got != want {
		t.Errorf("defaults:\n got %s\nwant %s", got, want)
	}
	m := filepath.Join(root, "system-flow.yaml")
	data, err := os.ReadFile(m)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(m, append(data, []byte("planning:\n  replan: agent\n  schedule: daily\n")...), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, errOut, code := runIn(t, root, "serve", "enable", "plan"); code != 0 {
		t.Fatalf("enable plan: %d %s", code, errOut)
	}
	if got, want := planning(), `{"next":"2026-10-05T00:00:00Z","plan":true,"replan":"agent","replan_set":true,"schedule":"daily"}`; got != want {
		t.Errorf("set:\n got %s\nwant %s", got, want)
	}
	if err := os.WriteFile(m, append(data, []byte("planning:\n  replan: sometimes\n  schedule: \"61 * * * *\"\n")...), 0o644); err != nil {
		t.Fatal(err)
	}
	got := planning()
	for _, want := range []string{`"replan":"sometimes"`, `"replan_error":"planning.replan \"sometimes\" is not a replan policy`, `"schedule":"61 * * * *"`, `"schedule_error":"planning.schedule`} {
		if !strings.Contains(got, want) {
			t.Errorf("errors lack %s:\n%s", want, got)
		}
	}
	if strings.Contains(got, `"next"`) {
		t.Errorf("a schedule flai cannot read has no next run: %s", got)
	}
}

// S-0229: settings.get gives the strategic agents' settings, each with its
// value as the manifest has it or unset, its default, and what it means,
// with whether the settings action lets the dashboard change them and what
// turns it on; settings.manifest changes them through flai manifest set only
// with that action on, and flai's refusal reaches the caller as Refused with
// the field and the reason, the manifest left as it was.
func TestSettingsShowAndSaveTheStrategicAgentsSettings(t *testing.T) {
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	root := tempProject(t)
	file := filepath.Join(root, manifest.File)
	data, err := os.ReadFile(file)
	if err != nil {
		t.Fatal(err)
	}
	fixture := "orchestration:\n  policy: wsjf\n  permissions:\n    promote_to_ready: true\n  release:\n    policy: threshold\n    value: 500\nplanning:\n  agent:\n    model: claude-sonnet-5\n"
	if err := os.WriteFile(file, append(data, fixture...), 0o644); err != nil {
		t.Fatal(err)
	}
	inProcess(t)
	type setting struct {
		Key     string          `json:"key"`
		Kind    string          `json:"kind"`
		Values  []string        `json:"values"`
		Default string          `json:"default"`
		Meaning string          `json:"meaning"`
		Risk    string          `json:"risk"`
		Whole   bool            `json:"whole"`
		Value   json.RawMessage `json:"value"`
		Set     bool            `json:"set"`
	}
	type strategic struct {
		Editable bool      `json:"editable"`
		Enable   string    `json:"enable"`
		Settings []setting `json:"settings"`
	}
	read := func() (strategic, map[string]setting) {
		t.Helper()
		out, _, code := runIn(t, root, "hostapi", "settings.get", `{}`)
		var got struct {
			Host struct {
				Strategic strategic `json:"strategic"`
			} `json:"host"`
		}
		if err := json.Unmarshal([]byte(out), &got); code != 0 || err != nil {
			t.Fatalf("settings.get: %d %v %s", code, err, out)
		}
		byKey := map[string]setting{}
		for _, s := range got.Host.Strategic.Settings {
			byKey[s.Key] = s
		}
		return got.Host.Strategic, byKey
	}

	st, byKey := read()
	if st.Editable || st.Enable != "flai serve enable settings" || len(st.Settings) != len(manifest.Settings()) {
		t.Errorf("strategic: editable %v, enable %q, %d settings, want false, flai serve enable settings, %d", st.Editable, st.Enable, len(st.Settings), len(manifest.Settings()))
	}
	for _, c := range []struct {
		key, kind, value, def string
		set                   bool
	}{
		{"orchestration.permissions.promote_to_ready", "boolean", `true`, "false", true},
		{"orchestration.permissions.accept_reviews", "boolean", ``, "false", false},
		{"orchestration.permissions.answer_threads", "choice", ``, "off", false},
		{"orchestration.policy", "choice", `"wsjf"`, "fifo", true},
		{"orchestration.release.policy", "choice", `"threshold"`, "judgement", true},
		{"orchestration.release.value", "number", `500`, "", true},
		{"orchestration.release.count", "number", ``, "", false},
		{"planning.agent", "agent", `{"model":"claude-sonnet-5"}`, "", true},
		{"planning.replan", "choice", ``, "deterministic", false},
		{"analysis.schedule", "cron", ``, "", false},
	} {
		s, ok := byKey[c.key]
		if !ok || s.Kind != c.kind || string(s.Value) != c.value || s.Default != c.def || s.Set != c.set || s.Meaning == "" {
			t.Errorf("%s: %+v (value %s), want kind %s, value %q, default %q, set %v, and a meaning", c.key, s, s.Value, c.kind, c.value, c.def, c.set)
		}
	}
	if s := byKey["orchestration.permissions.promote_to_ready"]; s.Risk == "" || !reflect.DeepEqual(s.Values, []string{"true", "false"}) {
		t.Errorf("a permission says its risk and its values: %+v", s)
	}
	if s := byKey["orchestration.release.count"]; !s.Whole {
		t.Errorf("a count is a whole number: %+v", s)
	}
	if s := byKey["orchestration.policy"]; !reflect.DeepEqual(s.Values, manifest.OrderPolicies) {
		t.Errorf("a choice lists its values: %+v", s)
	}

	save := func(params string) (string, int) {
		t.Helper()
		out, _, code := runIn(t, root, "hostapi", "settings.manifest", params)
		return out, code
	}
	written := func() string {
		t.Helper()
		b, err := os.ReadFile(file)
		if err != nil {
			t.Fatal(err)
		}
		return string(b)
	}
	before := written()
	if out, code := save(`{"set":{"planning.replan":"agent"},"request_id":"req-00000091"}`); code == 0 || !strings.Contains(out, "-32012") || !strings.Contains(out, "flai serve enable settings") {
		t.Errorf("with the settings action off: %d %s", code, out)
	}
	runIn(t, root, "serve", "enable", "settings")
	if out, code := save(`{"set":{"planning.replan":"agent","orchestration.release.value":-1},"request_id":"req-00000092"}`); code == 0 || !strings.Contains(out, "-32010") ||
		!strings.Contains(out, `"field":"orchestration.release.value"`) || !strings.Contains(out, `"reason":"`) {
		t.Errorf("flai's refusal reaches the caller with the field and the reason: %d %s", code, out)
	}
	if got := written(); got != before {
		t.Errorf("a refused save wrote the manifest:\n%s", got)
	}

	if out, code := save(`{"set":{"planning.replan":"agent"},"unset":["orchestration.policy"],"request_id":"req-00000093"}`); code != 0 || !strings.Contains(out, `{"key":"planning.replan","value":"agent"}`) {
		t.Fatalf("a save: %d %s", code, out)
	}
	st, byKey = read()
	if !st.Editable {
		t.Error("with the settings action on, strategic says editable")
	}
	if s := byKey["planning.replan"]; string(s.Value) != `"agent"` || !s.Set {
		t.Errorf("planning.replan after the save: %+v (value %s)", s, s.Value)
	}
	if s := byKey["orchestration.policy"]; s.Set || len(s.Value) != 0 || s.Default != "fifo" {
		t.Errorf("orchestration.policy after it was unset: %+v (value %s)", s, s.Value)
	}
	if js, _, _ := runIn(t, root, "serve", "journal"); !strings.Contains(js, "flai manifest set --autocommit --unset=orchestration.policy -- planning.replan=agent") {
		t.Errorf("the journal names the change as the command it ran:\n%s", js)
	}
}

// ADR-0067: auto-publish is the operator's shell tool, outside the workflow.
// The serve actions listing names it and serve enable turns it on, saying
// what it has flai push do; the dashboard is neither told of it nor can turn
// it on or off.
func TestAutoPublishIsTheShellsAlone(t *testing.T) {
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	root := tempProject(t)
	out, _, _ := runIn(t, root, "serve", "actions")
	if !strings.Contains(out, "auto-publish: off everywhere\n  has flai push --pending, run in a shell") || !strings.Contains(out, "push: off everywhere\n  lets a dashboard publish when you press Publish") {
		t.Errorf("actions lists auto-publish for the operator: %s", out)
	}
	out, errOut, code := runIn(t, root, "serve", "enable", "auto-publish")
	if code != 0 || !strings.Contains(out, "auto-publish enabled for t\n  it has flai push --pending") || strings.Contains(out, "journal") || !strings.Contains(out, "flai serve disable auto-publish") {
		t.Fatalf("enable: %d %s %s", code, out, errOut)
	}
	if info, _, _ := runIn(t, root, "hostapi", "project.info"); strings.Contains(info, "auto-publish") || !strings.Contains(info, `"push":false`) {
		t.Errorf("project.info names auto-publish: %s", info)
	}
	runIn(t, root, "serve", "enable", "settings")
	out, _, code = runIn(t, root, "hostapi", "settings.action", `{"action":"auto-publish","on":false,"request_id":"3f0c1a52-7d3b-4f0e-9a51-0c2d4e6f8a20"}`)
	if code == 0 || !strings.Contains(out, "flai serve disable auto-publish") || !strings.Contains(out, "in a shell on the host") {
		t.Errorf("settings.action refuses auto-publish, naming what to run: %d %s", code, out)
	}
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	if cfg, _, _ := (&app{}).loadConfig(); !cfg.ActionEnabled("auto-publish", mainRootOf(repo)) {
		t.Error("the refusal left auto-publish as the shell set it")
	}
	got, _ := json.Marshal((&app{}).hostSettings(mainRootOf(repo)))
	if strings.Contains(string(got), "auto-publish") || !strings.Contains(string(got), `"name":"push"`) {
		t.Errorf("the settings page is told of auto-publish: %s", got)
	}
}

// S-0387, ADR-0158: auto-upgrade, which has flai serve upgrade the flai on
// the host, is the operator's shell tool: serve enable turns it on for the
// project, saying what it has flai serve do, and the dashboard is neither
// told of it nor can turn it on or off.
func TestAutoUpgradeIsTheShellsAlone(t *testing.T) {
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	root := tempProject(t)
	out, errOut, code := runIn(t, root, "serve", "enable", "auto-upgrade")
	if code != 0 || !strings.Contains(out, "auto-upgrade enabled for t\n  it has flai serve ask flai host") || strings.Contains(out, "journal") || !strings.Contains(out, "flai serve disable auto-upgrade") {
		t.Fatalf("enable: %d %s %s", code, out, errOut)
	}
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	if !(&app{}).host().Enabled(hostapi.ActionAutoUpgrade, mainRootOf(repo)) {
		t.Error("flai serve sees auto-upgrade on for the project")
	}
	if info, _, _ := runIn(t, root, "hostapi", "project.info"); strings.Contains(info, "auto-upgrade") {
		t.Errorf("project.info names auto-upgrade: %s", info)
	}
	runIn(t, root, "serve", "enable", "settings")
	out, _, code = runIn(t, root, "hostapi", "settings.action", `{"action":"auto-upgrade","on":false,"request_id":"5d2e3b61-8c4a-4e1f-9b62-1d3e5f7a9b31"}`)
	if code == 0 || !strings.Contains(out, "flai serve disable auto-upgrade") {
		t.Errorf("settings.action refuses auto-upgrade, naming what to run: %d %s", code, out)
	}
	got, _ := json.Marshal((&app{}).hostSettings(mainRootOf(repo)))
	if strings.Contains(string(got), "auto-upgrade") {
		t.Errorf("the settings page is told of auto-upgrade: %s", got)
	}
}

// ADR-0158: flai serve asks the host it runs under for the project's release
// by version, as flai host upgrade --version does, and hears whether the
// host restarts on it or why it refused; a serve under no host asks none.
func TestServeAsksItsHostForTheRelease(t *testing.T) {
	var asked []string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		asked = append(asked, r.Method+" "+r.URL.Path+" "+string(body)+" "+r.Header.Get("Authorization"))
		if strings.Contains(string(body), "9.9.9") {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"error":"flai 9.9.9 is not published; published are 1.50.0"}`))
			return
		}
		_, _ = w.Write([]byte(`{"upgrade":{"installed":"1.50.0"},"restarting":true}`))
	}))
	t.Cleanup(srv.Close)
	upgrade := hostUpgrade(hostMCP{c: &host.Client{URL: srv.URL, Token: "host-token"}})
	if upgrade == nil {
		t.Fatal("a serve under a host has no way to ask it for an upgrade")
	}
	restarting, err := upgrade(context.Background(), "1.50.0")
	if err != nil || !restarting {
		t.Errorf("an installed release: %v %v", restarting, err)
	}
	if len(asked) != 1 || asked[0] != `POST /upgrade {"version":"1.50.0"} Bearer host-token` {
		t.Errorf("the host was asked %q", asked)
	}
	if restarting, err := upgrade(context.Background(), "9.9.9"); err == nil || restarting || !strings.Contains(err.Error(), "not published") {
		t.Errorf("a refused release: %v %v", restarting, err)
	}
	if hostUpgrade(nil) != nil {
		t.Error("a serve under no host asks one for an upgrade")
	}
}

// S-0116: flai serve agent restart says why it will not start an agent.
func TestServeAgentRestartSaysWhyItRefuses(t *testing.T) {
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	root := tempProject(t)
	if _, errOut, code := runIn(t, root, "serve", "agent", "restart", "S-0001"); code == 0 || !strings.Contains(errOut, "rule: the agent host action is off for this project") {
		t.Errorf("action off: %d %s", code, errOut)
	}
	runIn(t, root, "serve", "enable", "agent")
	if _, errOut, code := runIn(t, root, "serve", "agent", "restart", "S-1"); code == 0 || !strings.Contains(errOut, "S-0001") {
		t.Errorf("no such story: %d %s", code, errOut)
	}
}

// S-0118: flai serve agent restart queues the agent of a story in ready
// while the in-progress limit is full, and says so.
func TestServeAgentRestartQueuesWhenTheLimitIsFull(t *testing.T) {
	cfg := filepath.Join(t.TempDir(), "cfg.json")
	t.Setenv("FLAI_CONFIG", cfg)
	root := tempProject(t)
	runIn(t, root, "epic", "new", "Epic")
	for _, s := range []struct{ id, file string }{{"S-0001", "S-0001-slice.md"}, {"S-0002", "S-0002-busy.md"}} {
		runIn(t, root, "story", "new", strings.TrimSuffix(s.file[7:], ".md"), "--epic", "E-0001")
		file := filepath.Join(root, "wip/kanban/stories", s.file)
		body, _ := os.ReadFile(file)
		_ = os.WriteFile(file, []byte(strings.Replace(string(body), "## Acceptance criteria\n- [ ]\n", "## Acceptance criteria\n- [ ] ok\n", 1)), 0o644)
		if _, errOut, code := runIn(t, root, "move", s.id, "ready"); code != 0 {
			t.Fatal(errOut)
		}
	}
	runIn(t, root, "serve", "enable", "agent")
	runIn(t, root, "serve", "agent", "set", "--", "true", "{story}")
	runIn(t, root, "serve", "agent", "harness", "command", "--deny-protected") // the command cannot ask (ADR-0131)
	if out, errOut, code := runIn(t, root, "serve", "agent", "start", "S-1"); code != 0 {
		t.Fatalf("start: %d %s %s", code, out, errOut)
	}
	// its agent failed with the story in ready; another story takes the only slot
	agents := filepath.Join(string(serve.DirFor(cfg)), "agents.json")
	var all map[string]serve.AgentState
	data, _ := os.ReadFile(agents)
	if err := json.Unmarshal(data, &all); err != nil {
		t.Fatal(err)
	}
	for key, st := range all {
		r := st.Stories["S-0001"]
		r.PID, r.Ended, r.Outcome, r.Why = 0, r.Started, serve.OutcomeFailed, "ended (exit 1) with S-0001 in ready"
		all[key] = st
	}
	data, _ = json.Marshal(all)
	_ = os.WriteFile(agents, data, 0o600)
	runIn(t, root, "move", "S-0002", "in-progress")
	_ = os.WriteFile(filepath.Join(root, "wip/kanban/board.md"), []byte("---\ntitle: Board\nstatus: active\nwip_limits:\n  ready: 5\n  in-progress: 1\n  review: 3\n---\n\n# Board\n"), 0o644)

	out, errOut, code := runIn(t, root, "serve", "agent", "restart", "S-1")
	if code != 0 || !strings.HasPrefix(out, "queued another agent for S-0001: the in-progress limit, a full review, or a hold leaves no room") {
		t.Fatalf("restart: %d %s %s", code, out, errOut)
	}
	if _, errOut, code := runIn(t, root, "serve", "agent", "restart", "S-1"); code == 0 || !strings.Contains(errOut, "already queued") {
		t.Errorf("twice: %d %s", code, errOut)
	}
	js, _, _ := runIn(t, root, "serve", "journal", "--json")
	if !strings.Contains(js, "queued another agent for S-0001") {
		t.Errorf("the queue is not journalled: %s", js)
	}
}

// S-0115: flai serve agent start starts a ready story's agent now, recorded
// where flai serve tracks it, and says why when it will not.
func TestServeAgentStartStartsAReadyStorysAgent(t *testing.T) {
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	root := tempProject(t)
	runIn(t, root, "epic", "new", "Epic")
	runIn(t, root, "story", "new", "Slice", "--epic", "E-0001")
	file := filepath.Join(root, "wip/kanban/stories/S-0001-slice.md")
	s, _ := os.ReadFile(file)
	_ = os.WriteFile(file, []byte(strings.Replace(string(s), "## Acceptance criteria\n- [ ]\n", "## Acceptance criteria\n- [ ] ok\n", 1)), 0o644)
	if _, errOut, code := runIn(t, root, "serve", "agent", "start", "S-1"); code == 0 || !strings.Contains(errOut, "rule: the agent host action is off for this project") {
		t.Errorf("action off: %d %s", code, errOut)
	}
	runIn(t, root, "serve", "enable", "agent")
	if _, errOut, code := runIn(t, root, "serve", "agent", "start", "S-1"); code == 0 || !strings.Contains(errOut, "rule: S-0001 is in backlog") {
		t.Errorf("in backlog: %d %s", code, errOut)
	}
	if _, errOut, code := runIn(t, root, "move", "S-0001", "ready"); code != 0 {
		t.Fatal(errOut)
	}
	if _, errOut, code := runIn(t, root, "serve", "agent", "start", "S-1"); code == 0 || !strings.Contains(errOut, "rule: S-0001 names no harness, and no command is set") {
		t.Errorf("nothing to start it with: %d %s", code, errOut)
	}
	runIn(t, root, "serve", "agent", "set", "--", "true", "{story}")
	// S-0355 (ADR-0131): the command cannot ask before a protected write, so
	// it starts only once the operator lifts that
	if _, errOut, code := runIn(t, root, "serve", "agent", "start", "S-1"); code == 0 || !strings.Contains(errOut, "harness command cannot ask before a protected write") || !strings.Contains(errOut, "flai serve agent harness command --deny-protected") {
		t.Errorf("started where nothing lifts asking: %d %s", code, errOut)
	}
	runIn(t, root, "serve", "agent", "harness", "command", "--deny-protected")
	out, errOut, code := runIn(t, root, "serve", "agent", "start", "S-1")
	if code != 0 || !strings.HasPrefix(out, "started true (command) for S-0001 as agent-S-0001") {
		t.Fatalf("start: %d %s %s", code, out, errOut)
	}
	js, _, _ := runIn(t, root, "serve", "journal", "--json")
	if !strings.Contains(js, "started true (command) for S-0001") {
		t.Errorf("the start is not journalled: %s", js)
	}
}

// S-0140: flai serve agent commit starts an agent to commit what a story in
// review left uncommitted in its worktree, and says why when it will not.
func TestServeAgentCommitStartsAnAgentForAnUncommittedWorktree(t *testing.T) {
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	root := tempProject(t)
	runIn(t, root, "epic", "new", "Epic")
	runIn(t, root, "story", "new", "Slice", "--epic", "E-0001")
	file := filepath.Join(root, "wip/kanban/stories/S-0001-slice.md")
	s, _ := os.ReadFile(file)
	_ = os.WriteFile(file, []byte(strings.Replace(string(s), "## Acceptance criteria\n- [ ]\n", "## Acceptance criteria\n- [ ] ok\n", 1)), 0o644)
	if _, errOut, code := runIn(t, root, "serve", "agent", "commit", "S-1"); code == 0 || !strings.Contains(errOut, "rule: the agent host action is off for this project") {
		t.Errorf("action off: %d %s", code, errOut)
	}
	runIn(t, root, "serve", "enable", "agent")
	runIn(t, root, "serve", "agent", "set", "--", "true", "{story}")
	for _, args := range [][]string{{"move", "S-0001", "ready"}, {"move", "S-0001", "in-progress"}, {"task", "new", "Do", "--story", "S-0001"}, {"move", "S-0001", "review"}} {
		if _, errOut, code := runIn(t, root, args...); code != 0 {
			t.Fatalf("%v: %s", args, errOut)
		}
	}
	if _, errOut, code := runIn(t, root, "serve", "agent", "commit", "S-1"); code == 0 || !strings.Contains(errOut, "rule: S-0001 has no worktree at .flai-cache/worktrees/S-0001") {
		t.Errorf("no worktree: %d %s", code, errOut)
	}
}

// S-0142: flai serve agent stream prints what a story's agent said and did,
// from the log flai serve gave it, and the host answers agent.stream from it.
func TestServeAgentStreamPrintsTheLog(t *testing.T) {
	cfg := filepath.Join(t.TempDir(), "cfg.json")
	t.Setenv("FLAI_CONFIG", cfg)
	root := tempProject(t)
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	log := filepath.Join(t.TempDir(), "sf-S-0001.log")
	lines := `{"type":"system","subtype":"init","model":"m"}` + "\n" +
		`{"type":"assistant","message":{"content":[{"type":"text","text":"Reading.\nThen writing."},{"type":"tool_use","name":"Bash","input":{"description":"List files"}}]}}` + "\n" +
		`{"type":"user","message":{"content":[{"type":"tool_result","content":"nope","is_error":true}]}}` + "\n"
	if err := os.WriteFile(log, []byte(lines), 0o600); err != nil {
		t.Fatal(err)
	}
	dir := serve.DirFor(cfg)
	_ = os.MkdirAll(string(dir), 0o700)
	state := map[string]serve.AgentState{mainRootOf(repo): {Stories: map[string]*serve.AgentRun{
		"S-0001": {Story: "S-0001", Agent: "agent-S-0001", Started: "2026-09-29T05:00:00Z", Ended: "2026-09-29T05:01:00Z", Log: log, Outcome: serve.OutcomeWorked},
	}}}
	data, _ := json.Marshal(state)
	if err := os.WriteFile(filepath.Join(string(dir), "agents.json"), data, 0o600); err != nil {
		t.Fatal(err)
	}

	out, errOut, code := runIn(t, root, "serve", "agent", "stream", "S-1")
	want := "S-0001: agent-S-0001, started 2026-09-29T05:00:00Z, ended 2026-09-29T05:01:00Z\n" +
		"session   session started (m)\n" +
		"text      Reading.\n          Then writing.\n" +
		"tool      Bash: List files\n" +
		"result!   nope\n"
	if code != 0 || out != want {
		t.Errorf("stream: %d %s\n%q\nwant\n%q", code, errOut, out, want)
	}
	// --follow ends with the agent, which has ended
	if out, _, code := runIn(t, root, "serve", "agent", "stream", "S-0001", "--follow", "--from", "0"); code != 0 || out != want {
		t.Errorf("follow an ended agent: %d %q", code, out)
	}
	first := strconv.Itoa(strings.Index(lines, "\n") + 1)
	js, _, _ := runIn(t, root, "serve", "agent", "stream", "S-0001", "--json", "--from", first)
	var got serve.StreamRead
	if err := json.Unmarshal([]byte(js), &got); err != nil || strconv.FormatInt(got.From, 10) != first || got.Next != int64(len(lines)) || len(got.Entries) != 3 || got.Running {
		t.Errorf("--json from %s: %v %s", first, err, js)
	}
	if _, errOut, code := runIn(t, root, "serve", "agent", "stream", "S-0002"); code == 0 || !strings.Contains(errOut, "rule: flai serve has started no agent for S-0002") {
		t.Errorf("no agent: %d %s", code, errOut)
	}

	// the host answers agent.stream from the same state
	a := &app{}
	res, e := hostapi.MethodsFor("test", nil, a.host())["agent.stream"](context.Background(), channel.Project{Key: "t", Root: mainRootOf(repo)}, json.RawMessage(`{"story":"S-0001","after":`+first+`}`))
	if e != nil {
		t.Fatal(e)
	}
	if r := res.(*serve.StreamRead); strconv.FormatInt(r.From, 10) != first || len(r.Entries) != 3 {
		t.Errorf("agent.stream: %+v", r)
	}
}

// S-0177, ADR-0064: flai mcp's agent tools start a story's agent as flai
// serve agent restart does, under the same agent host action, and journal
// the agent that asked.
func TestTheMCPServerStartsAStorysAgentUnderTheAgentAction(t *testing.T) {
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	t.Setenv("FLAI_STARTED_BY", "") // the operator's own agent, whoever runs the tests
	root := tempProject(t)
	runIn(t, root, "epic", "new", "Epic")
	runIn(t, root, "story", "new", "Slice", "--epic", "E-0001")
	file := filepath.Join(root, "wip/kanban/stories/S-0001-slice.md")
	s, _ := os.ReadFile(file)
	_ = os.WriteFile(file, []byte(strings.Replace(string(s), "## Acceptance criteria\n- [ ]\n", "## Acceptance criteria\n- [ ] ok\n", 1)), 0o644)
	for _, to := range []string{"ready", "in-progress"} {
		if _, errOut, code := runIn(t, root, "move", "S-0001", to); code != 0 {
			t.Fatal(errOut)
		}
	}
	var out, errOut bytes.Buffer
	a := &app{out: &out, errOut: &errOut, cwd: root, clock: func() time.Time { return time.Date(2026, 10, 1, 10, 0, 0, 0, time.UTC) }}
	ctx := context.Background()
	if _, err := a.mcpAgents(ctx, "restart", root, "S-0001", "agent-ops"); err == nil || !strings.Contains(err.Error(), "the agent host action is off") {
		t.Errorf("action off: %v", err)
	}
	runIn(t, root, "serve", "enable", "agent")
	runIn(t, root, "serve", "agent", "set", "--", "true", "{story}")
	runIn(t, root, "serve", "agent", "harness", "command", "--deny-protected") // the command cannot ask (ADR-0131)
	got, err := a.mcpAgents(ctx, "restart", root, "S-0001", "agent-ops")
	if err != nil || got.Story != "S-0001" || got.PID == 0 || got.Command != "true" {
		t.Fatalf("a story in progress with no agent here: %+v %v", got, err)
	}
	// an agent flai serve started does not start others
	t.Setenv("FLAI_STARTED_BY", "flai-serve")
	if _, err := a.mcpAgents(ctx, "start", root, "S-0001", "agent-S-0002"); err == nil || !strings.Contains(err.Error(), "an agent flai serve started does not start agents") {
		t.Errorf("an agent flai serve started: %v", err)
	}
	js, _, _ := runIn(t, root, "serve", "journal", "--json")
	for _, want := range []string{`"method": "mcp.agent_restart"`, `"by": "agent-ops"`, `"outcome": "disabled"`, "agent-ops asked to restart S-0001's agent: started true"} {
		if !strings.Contains(js, want) {
			t.Errorf("journal lacks %q: %s", want, js)
		}
	}
}

// orchestratorState writes the project's orchestrator run, and the
// analyzer's when there is one, as flai serve records them.
func orchestratorState(t *testing.T, cfg, root string, orchestrator, analyzer *serve.AgentRun) {
	t.Helper()
	dir := serve.DirFor(cfg)
	_ = os.MkdirAll(string(dir), 0o700)
	data, _ := json.Marshal(map[string]serve.AgentState{root: {Orchestrator: orchestrator, Analyzer: analyzer}})
	if err := os.WriteFile(filepath.Join(string(dir), "agents.json"), data, 0o600); err != nil {
		t.Fatal(err)
	}
}

// S-0228: the stop holds the orchestrator stopped, here a run that failed
// and waits to be started again, and the start lifts the hold and starts
// it, each refused while the orchestrate action is off; the host's
// orchestrate.stop and orchestrate.start run them, and agent.status says
// the orchestrator is held.
func TestServeOrchestrateHoldsTheOrchestratorStopped(t *testing.T) {
	cfg := filepath.Join(t.TempDir(), "cfg.json")
	t.Setenv("FLAI_CONFIG", cfg)
	root := tempProject(t)
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	main := mainRootOf(repo)
	for _, verb := range []string{"stop", "start"} {
		if _, errOut, code := runIn(t, root, "serve", "orchestrate", verb); code == 0 || !strings.Contains(errOut, "rule: the orchestrate host action is off for this project, so flai serve runs no orchestrator to "+verb) {
			t.Errorf("%s while off: %d %s", verb, code, errOut)
		}
	}
	if _, errOut, code := runIn(t, root, "serve", "enable", "orchestrate"); code != 0 {
		t.Fatal(errOut)
	}
	if _, errOut, code := runIn(t, root, "serve", "orchestrate", "start"); code == 0 || !strings.Contains(errOut, "rule: the orchestrator is not held stopped") {
		t.Errorf("start with no hold: %d %s", code, errOut)
	}
	failed := &serve.AgentRun{Agent: "orchestrator", Harness: "command", Started: "2026-09-15T20:59:00Z", Ended: "2026-09-15T20:59:00Z", Session: "s1",
		Error: "no harness", Outcome: serve.OutcomeFailed, Why: "could not be started: no harness"}
	orchestratorState(t, cfg, main, failed, nil)
	out, errOut, code := runIn(t, root, "serve", "orchestrate", "stop", "--json")
	var said struct {
		Held         bool           `json:"held"`
		Orchestrator serve.AgentRun `json:"orchestrator"`
	}
	if err := json.Unmarshal([]byte(out), &said); code != 0 || err != nil || !said.Held || !said.Orchestrator.Held || said.Orchestrator.Session != "s1" || said.Orchestrator.Outcome != serve.OutcomeFailed {
		t.Fatalf("stop: %d %s %s", code, out, errOut)
	}
	if _, errOut, code := runIn(t, root, "serve", "orchestrate", "stop"); code == 0 || !strings.Contains(errOut, "rule: the orchestrator is already held stopped") {
		t.Errorf("stop twice: %d %s", code, errOut)
	}
	a := &app{out: &bytes.Buffer{}, errOut: &bytes.Buffer{}}
	p := channel.Project{Key: "t", Root: main}
	res, e := hostapi.MethodsFor("test", nil, a.host())["agent.status"](context.Background(), p, json.RawMessage(`{}`))
	if b, _ := json.Marshal(res); e != nil || !strings.Contains(string(b), `"orchestrator":{`) || !strings.Contains(string(b), `"session":"s1","outcome":"failed","why":"could not be started: no harness","held":true}`) {
		t.Errorf("agent.status: %+v %s", e, b)
	}

	// through the host: a start that fails, as with no harness and no
	// command, lifts the hold as well, and a stop holds the failed run
	inProcess(t)
	m := hostapi.MethodsFor("test", nil, a.host())
	if _, e := m["orchestrate.start"](context.Background(), p, json.RawMessage(`{"request_id":"req-00000001"}`)); e == nil || !strings.Contains(e.Message, "the orchestrator's agent names no harness") {
		t.Fatalf("orchestrate.start: %+v", e)
	}
	run := a.serveDir().AgentStates()[main].Orchestrator
	if run == nil || run.Held || run.Session == "s1" || run.Outcome != serve.OutcomeFailed {
		t.Fatalf("after the start: %+v, want a new failed run, not held", run)
	}
	res, e = m["orchestrate.stop"](context.Background(), p, json.RawMessage(`{"request_id":"req-00000002"}`))
	w, _ := res.(hostapi.Written)
	if err := json.Unmarshal(w.Data, &said); e != nil || err != nil || !said.Held || !said.Orchestrator.Held || said.Orchestrator.Session != run.Session {
		t.Fatalf("orchestrate.stop: %s %+v", w.Data, e)
	}
	if r := a.serveDir().AgentStates()[main].Orchestrator; !r.Held || r.Session != run.Session {
		t.Errorf("after the stop: %+v, want the failed run held", r)
	}
	js, _, _ := runIn(t, root, "serve", "journal", "--json")
	for _, want := range []string{`"method": "orchestrate.start"`, `"method": "orchestrate.stop"`,
		"stopped the orchestrator orchestrator on the operator's word, and held it stopped while orchestrate stays on: it was not running"} {
		if !strings.Contains(js, want) {
			t.Errorf("journal lacks %q: %s", want, js)
		}
	}
}

// S-0228: the host answers agent.stream with role from the analyzer's newest
// run, as flai serve recorded it, and is not found with no run of the role.
func TestTheHostAnswersARolesStream(t *testing.T) {
	cfg := filepath.Join(t.TempDir(), "cfg.json")
	t.Setenv("FLAI_CONFIG", cfg)
	root := tempProject(t)
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	log := filepath.Join(t.TempDir(), "t-analyzer.log")
	if err := os.WriteFile(log, []byte(`{"type":"assistant","message":{"content":[{"type":"text","text":"Reading the metrics."}]}}`+"\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	orchestratorState(t, cfg, mainRootOf(repo), nil, &serve.AgentRun{Agent: "analyzer", Focus: "risk", Started: "2026-09-15T20:00:00Z", Ended: "2026-09-15T20:05:00Z", Outcome: serve.OutcomeWorked, Log: log})
	a := &app{}
	m := hostapi.MethodsFor("test", nil, a.host())["agent.stream"]
	p := channel.Project{Key: "t", Root: mainRootOf(repo)}
	res, e := m(context.Background(), p, json.RawMessage(`{"role":"analyze"}`))
	if e != nil {
		t.Fatal(e)
	}
	if r := res.(*serve.StreamRead); r.Role != "analyze" || r.Agent != "analyzer" || r.Running || len(r.Entries) != 1 || r.Entries[0].Text != "Reading the metrics." {
		t.Errorf("agent.stream: %+v", r)
	}
	if _, e := m(context.Background(), p, json.RawMessage(`{"role":"orchestrate"}`)); e == nil || e.Code != hostapi.NotFound || !strings.Contains(e.Message, "flai serve has started no orchestrator for this project") {
		t.Errorf("no orchestrator run: %+v", e)
	}
}

// S-0489: settings.get's host block carries dashboard.allowed_hosts as the
// configuration names it, normalized as the container is given it: empty
// until it is set.
func TestServeSettingsShowTheAllowedHosts(t *testing.T) {
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	root := tempProject(t)
	allowed := func() any {
		t.Helper()
		got, err := json.Marshal((&app{}).hostSettings(root))
		if err != nil {
			t.Fatal(err)
		}
		var settings map[string]any
		if err := json.Unmarshal(got, &settings); err != nil {
			t.Fatal(err)
		}
		v, ok := settings["allowed_hosts"]
		if !ok {
			t.Fatalf("settings lack allowed_hosts: %s", got)
		}
		return v
	}
	if _, errOut, code := runIn(t, root, "config", "get", "dashboard.allowed_hosts"); code != 0 { // creates the config
		t.Fatalf("config get: %s", errOut)
	}
	if got := allowed(); got != "" {
		t.Errorf("allowed_hosts before it is set = %v, want empty", got)
	}
	if _, errOut, code := runIn(t, root, "config", "set", "dashboard.allowed_hosts", " Flow.Example.com, ,board.example.net:8443"); code != 0 {
		t.Fatalf("config set: %s", errOut)
	}
	if got := allowed(); got != "flow.example.com,board.example.net:8443" {
		t.Errorf("allowed_hosts = %v, want the normalized list", got)
	}
}

// S-0498 (ADR-0205): flai serve notify set keeps a webhook and a shell
// command on the host, changing only what it is given, show prints them,
// with --json too, clear removes them, and a URL that is not http or https,
// an empty command, and a set with neither are refused, changing nothing.
func TestServeNotify(t *testing.T) {
	cfgPath := filepath.Join(t.TempDir(), "cfg.json")
	t.Setenv("FLAI_CONFIG", cfgPath)
	root := tempProject(t)
	run := func(args ...string) (string, string, int) {
		t.Helper()
		return runIn(t, root, append([]string{"serve", "notify"}, args...)...)
	}
	saved := func() config.Notify {
		t.Helper()
		cfg, _, err := config.Load(cfgPath)
		if err != nil {
			t.Fatal(err)
		}
		return cfg.Notify
	}
	if out, _, code := run("show"); code != 0 || !strings.Contains(out, "no reminder channel is set") {
		t.Errorf("unset: %d %q", code, out)
	}
	if out, errOut, code := run("set", "--url", "https://ntfy.example/flai"); code != 0 || out != "url: https://ntfy.example/flai\n" {
		t.Fatalf("set --url: %d %q %s", code, out, errOut)
	}
	if out, _, code := run("show"); code != 0 || out != "url: https://ntfy.example/flai\n" {
		t.Errorf("show: %d %q", code, out)
	}
	const command = `curl -s -H "Title: flai" -d @- https://ntfy.sh/flai`
	if _, errOut, code := run("set", "--command", command); code != 0 {
		t.Fatalf("set --command: %d %s", code, errOut)
	}
	if got := saved(); got != (config.Notify{URL: "https://ntfy.example/flai", Command: command}) {
		t.Errorf("set changes only what it is given: %+v", got)
	}
	if js, _, _ := run("show", "--json"); !strings.Contains(js, `"url": "https://ntfy.example/flai"`) || !strings.Contains(js, `"command": "curl -s -H \"Title: flai\" -d @- https://ntfy.sh/flai"`) {
		t.Errorf("--json: %s", js)
	}
	for _, bad := range []struct {
		args []string
		want string
	}{
		{[]string{"set", "--url", "ftp://x"}, "notify.url's scheme is ftp, and flai posts reminders over http or https alone"},
		{[]string{"set", "--url", ""}, "--url needs an http or https webhook"},
		{[]string{"set", "--command", " "}, "--command needs a shell command"},
		{[]string{"set"}, "give --url <webhook>, --command <shell command>, or both"},
	} {
		if _, errOut, code := run(bad.args...); code == 0 || !strings.Contains(errOut, bad.want) {
			t.Errorf("%q is refused, saying why: %d %s", bad.args, code, errOut)
		}
	}
	if got := saved(); got.URL != "https://ntfy.example/flai" || got.Command != command {
		t.Errorf("a refused value changes nothing: %+v", got)
	}
	if out, _, code := run("clear"); code != 0 || !strings.Contains(out, "no reminder channel is set") || saved() != (config.Notify{}) {
		t.Errorf("clear: %d %q %+v", code, out, saved())
	}
	if data, _ := os.ReadFile(cfgPath); strings.Contains(string(data), "notify") {
		t.Errorf("a cleared channel is not written: %s", data)
	}
}

// S-0417: the host action a role the project defines names is one flai
// serve enable and disable take for the project, flai serve actions and the
// settings page list it, flai serve is told whether it is on, and
// agent.status carries the role's runs under roles, by name and item.
func TestServeEnableTakesADefinedRolesAction(t *testing.T) {
	root := rolesProject(t, map[string]string{"reviewer.md": reviewerRole})
	cfgPath := os.Getenv("FLAI_CONFIG")
	if out, _, _ := runIn(t, root, "serve", "actions"); !strings.Contains(out, "review: off everywhere\n  lets a dashboard start the role reviewer, which design/workflows/roles/reviewer.md defines") {
		t.Errorf("actions names it and what it means: %s", out)
	}
	if (&app{}).agentConfig(root).RoleActions["review"] {
		t.Error("review is on by default")
	}
	status := func() string {
		t.Helper()
		got, err := json.Marshal((&app{}).host().Agent(root))
		if err != nil {
			t.Fatal(err)
		}
		return string(got)
	}
	if got := status(); !strings.Contains(got, `"roles":{}`) {
		t.Errorf("agent.status before any run: %s", got)
	}
	if _, errOut, code := runIn(t, root, "serve", "enable", "audit"); code == 0 || !strings.Contains(errOut, `there is no host action \"audit\"`) || !strings.Contains(errOut, ", review, ") {
		t.Errorf("an action no role names: %d %s", code, errOut)
	}
	if out, errOut, code := runIn(t, root, "serve", "enable", "review"); code != 0 || !strings.Contains(out, "review enabled for t\n  it lets a dashboard start the role reviewer") || !strings.Contains(out, "every use is journalled") {
		t.Fatalf("enable: %d %s %s", code, out, errOut)
	}
	if !(&app{}).agentConfig(root).RoleActions["review"] {
		t.Error("enabled, flai serve is not told")
	}
	settings, err := json.Marshal((&app{}).hostSettings(root))
	if err != nil || !strings.Contains(string(settings), `"here":true,"means":"start the role reviewer`) || !strings.Contains(string(settings), `"name":"review"`) {
		t.Errorf("the settings page does not list review on: %s %v", settings, err)
	}
	run := `{"` + root + `":{"roles":{"reviewer":{"S-0001":{"story":"","item":"S-0001","role":"review","role_name":"reviewer","action":"review","harness":"command","command":"run-agent","agent":"reviewer-S-0001","pid":42,"started":"2026-10-28T10:00:00Z","session":"s","log":"/l/t-reviewer-20261028T100000Z.log"}}}}}`
	dir := string(serve.DirFor(cfgPath))
	if err := os.MkdirAll(dir, 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "agents.json"), []byte(run), 0o600); err != nil {
		t.Fatal(err)
	}
	if got := status(); !strings.Contains(got, `"roles":{"reviewer":{"S-0001":{`) || !strings.Contains(got, `"role_name":"reviewer"`) || !strings.Contains(got, `"agent":"reviewer-S-0001"`) {
		t.Errorf("agent.status with a run: %s", got)
	}
	// a role file gone, its action is still turned off
	if err := os.Remove(filepath.Join(root, "design", "workflows", "roles", "reviewer.md")); err != nil {
		t.Fatal(err)
	}
	if out, errOut, code := runIn(t, root, "serve", "disable", "review"); code != 0 || !strings.Contains(out, "review disabled for t") {
		t.Errorf("disable once its role is gone: %d %s %s", code, out, errOut)
	}
	if _, errOut, code := runIn(t, root, "serve", "enable", "review"); code == 0 || !strings.Contains(errOut, `there is no host action \"review\"`) {
		t.Errorf("enable once its role is gone: %d %s", code, errOut)
	}
}
