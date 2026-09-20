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

	if _, errOut, code := runIn(t, root, "serve", "enable", "pull"); code == 0 || !strings.Contains(errOut, "there is: push") {
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
