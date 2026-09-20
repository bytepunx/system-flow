package cmd

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/hostapi"
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
	// I-0028: an acceptance a dashboard asks for runs here as the operator,
	// with a remote it can push to, and pushes nothing, because nobody said
	// it may. From S-0075 until this story it pushed.
	out, _, code := runIn(t, root, "hostapi", "accept.run", `{"id":"S-0001","include_uncommitted":true,"request_id":"3f0c1a52-7d3b-4f0e-9a51-0c2d4e6f8a09"}`)
	if code != 0 || !strings.Contains(out, `"pushed":false`) || head(remote) != before || head(root) == before {
		t.Fatalf("an acceptance with the action off: accepted here, nothing pushed: %d %s", code, out)
	}
	if strings.Contains(gitIn(t, remote, "tag", "--list"), "cli/v1.1.0") {
		t.Fatal("nor its release tag")
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

	if _, errOut, code := runIn(t, root, "serve", "enable", "pull"); code == 0 || !strings.Contains(errOut, "there are: agent, push") {
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

	// on: pushed, with the tag
	out, _, code = runIn(t, root, "hostapi", "push.run", pushRequest)
	if code != 0 || !strings.Contains(out, `"pushed":true`) || head(remote) != head(root) || !strings.Contains(gitIn(t, remote, "tag", "--list"), "cli/v1.1.0") {
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
	if !strings.Contains(out, "no command is set, so nothing is started") || !strings.Contains(out, "off for this project") {
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
