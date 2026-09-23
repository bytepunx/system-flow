package cmd

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/hostapi"
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

const pushRequest = `{"request_id":"3f0c1a52-7d3b-4f0e-9a51-0c2d4e6f8a10"}`

// S-0078: a host action is off until the operator enables it by name on the
// host; asked for meanwhile it is refused and says what enables it; enabled,
// it runs; and every request is in the journal, refused ones too.
func TestHostActionPush(t *testing.T) {
	root, remote := researchProject(t, "feature", true)
	inProcess(t)
	head := func(dir string) string { return strings.TrimSpace(gitIn(t, dir, "rev-parse", "main")) }
	before := head(remote)
	// S-0087: acceptance itself never pushes or tags, whatever the push
	// action is set to; it only merges, archives, and commits locally.
	// S-0094: the tag is computed and created at push instead, not a
	// separate step, so it is not missing here for anyone to notice.
	out, _, code := runIn(t, root, "hostapi", "accept.run", `{"id":"S-0001","include_uncommitted":true,"request_id":"3f0c1a52-7d3b-4f0e-9a51-0c2d4e6f8a09"}`)
	if code != 0 || head(remote) != before || head(root) == before {
		t.Fatalf("an acceptance with the action off: accepted here, nothing pushed: %d %s", code, out)
	}

	out, _, _ = runIn(t, root, "serve", "actions")
	if !strings.Contains(out, "push: off everywhere") || !strings.Contains(out, "off for this project; flai serve enable push") || !strings.Contains(out, "publish any story") {
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
	out, _, code = runIn(t, root, "hostapi", "push.run", pushRequest)
	if code == 0 || !strings.Contains(out, "is not enabled for this project") || !strings.Contains(out, "flai serve enable push") || head(remote) != before {
		t.Fatalf("disabled: %d %s", code, out)
	}

	if _, errOut, code := runIn(t, root, "serve", "enable", "pull"); code == 0 || !strings.Contains(errOut, "there are: agent, checks, dashboard, push") {
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

	// on: pushed, tagging what accept left unreleased first (S-0094)
	out, _, code = runIn(t, root, "hostapi", "push.run", pushRequest)
	if code != 0 || !strings.Contains(out, `"pushed":true`) || !strings.Contains(out, `"tags":["cli/v1.1.0"]`) || head(remote) != head(root) {
		t.Fatalf("enabled: %d %s", code, out)
	}
	out, _, _ = runIn(t, root, "hostapi", "push.run", `{"request_id":"3f0c1a52-7d3b-4f0e-9a51-0c2d4e6f8a11"}`)
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
		if e.Outcome != want.outcome || e.Detail != want.detail || e.Action != "push" || e.Method != "push.run" || e.Root != root || e.By == "" || e.At == "" {
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
	if out, _, code := runIn(t, root, "hostapi", "push.run", `{"request_id":"3f0c1a52-7d3b-4f0e-9a51-0c2d4e6f8a12"}`); code == 0 || !strings.Contains(out, "not enabled") {
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
