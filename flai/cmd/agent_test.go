package cmd

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// S-0103: flai agent sets the project's default agent in system-flow.yaml,
// and every story created while it is set gets a copy.
func TestFlaiAgentSetsTheDefaultNewStoriesGet(t *testing.T) {
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	root := tempProject(t)
	if out, _, _ := runIn(t, root, "agent"); !strings.Contains(out, "no default agent") {
		t.Errorf("none yet: %s", out)
	}
	if _, errOut, code := runIn(t, root, "agent", "set"); code == 0 || !strings.Contains(errOut, "nothing to set") {
		t.Errorf("empty set: %d %s", code, errOut)
	}
	out, errOut, code := runIn(t, root, "agent", "set", "--harness", "claude-code", "--model", "claude-opus-5-5", "--config", "effort=high", "--config", "max_turns=50")
	if code != 0 || !strings.Contains(out, "default agent: claude-code, claude-opus-5-5, effort=high, max_turns=50") {
		t.Fatalf("set: %d %s %s", code, out, errOut)
	}
	out, _, _ = runIn(t, root, "agent", "set", "--model", "claude-sonnet-5", "--unset", "max_turns")
	if !strings.Contains(out, "default agent: claude-code, claude-sonnet-5, effort=high") || strings.Contains(out, "max_turns") {
		t.Errorf("changed one part: %s", out)
	}
	if _, errOut, code := runIn(t, root, "agent", "set", "--harness", "Not Valid"); code == 0 || !strings.Contains(errOut, "not a name such as claude-code") {
		t.Errorf("an invalid harness: %d %s", code, errOut)
	}

	runIn(t, root, "epic", "new", "E")
	if _, errOut, code := runIn(t, root, "story", "new", "S", "--epic", "E-0001"); code != 0 {
		t.Fatalf("story new: %s", errOut)
	}
	matches, _ := filepath.Glob(filepath.Join(root, "wip/kanban/stories/S-0001-*.md"))
	data, _ := os.ReadFile(matches[0])
	if !strings.Contains(string(data), "agent:\n  harness: claude-code\n  model: claude-sonnet-5\n  config:\n    effort: high\n") {
		t.Errorf("the story's front matter:\n%s", data)
	}

	if out, _, _ := runIn(t, root, "agent", "clear"); !strings.Contains(out, "no default agent") {
		t.Errorf("clear: %s", out)
	}
	manifestData, _ := os.ReadFile(filepath.Join(root, "system-flow.yaml"))
	if strings.Contains(string(manifestData), "agent:") {
		t.Errorf("cleared manifest:\n%s", manifestData)
	}
	// a story's agent and a cleared default raise nothing in check
	if out, _, _ := runIn(t, root, "check"); strings.Contains(out, "agent") {
		t.Errorf("check: %s", out)
	}
}

// S-0103: the story's own agent, given when it is created and changed with
// flai edit: merged into what it has, a key given with no value removed,
// --clear-agent alone removing it and with other flags replacing it.
func TestAStorysAgentIsSetAndEdited(t *testing.T) {
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	root := tempProject(t)
	runIn(t, root, "agent", "set", "--harness", "claude-code", "--model", "claude-opus-5-5", "--config", "effort=high")
	runIn(t, root, "epic", "new", "E")
	if _, errOut, code := runIn(t, root, "story", "new", "S", "--epic", "E-0001", "--model", "claude-sonnet-5", "--agent-config", "max_turns=30"); code != 0 {
		t.Fatal(errOut)
	}
	agentLine := func() string {
		out, _, _ := runIn(t, root, "show", "S-0001")
		for _, l := range strings.Split(out, "\n") {
			if strings.HasPrefix(strings.TrimSpace(l), "agent:") {
				return strings.TrimSpace(l)
			}
		}
		return ""
	}
	if got := agentLine(); got != "agent: claude-code, claude-sonnet-5, effort=high, max_turns=30" {
		t.Errorf("created: %q", got)
	}
	edit := func(args ...string) {
		t.Helper()
		if _, errOut, code := runIn(t, root, append([]string{"edit", "S-0001"}, args...)...); code != 0 {
			t.Fatalf("edit %v: %s", args, errOut)
		}
	}
	edit("--model", "claude-opus-5-5", "--agent-config", "effort=")
	if got := agentLine(); got != "agent: claude-code, claude-opus-5-5, max_turns=30" {
		t.Errorf("merged, effort removed: %q", got)
	}
	edit("--clear-agent", "--harness", "codex", "--model", "gpt-5.1")
	if got := agentLine(); got != "agent: codex, gpt-5.1" {
		t.Errorf("replaced: %q", got)
	}
	if out, _, _ := runIn(t, root, "edit", "S-0001", "--show"); !strings.Contains(out, "agent: codex, gpt-5.1 (project default: claude-code, claude-opus-5-5, effort=high)") {
		t.Errorf("--show: %s", out)
	}
	edit("--clear-agent")
	if got := agentLine(); got != "" {
		t.Errorf("cleared: %q", got)
	}
	if _, errOut, code := runIn(t, root, "edit", "E-0001", "--model", "x"); code == 0 || !strings.Contains(errOut, "only a story carries an agent") {
		t.Errorf("an epic's agent: %d %s", code, errOut)
	}
	if _, errOut, code := runIn(t, root, "story", "new", "Bad", "--harness", "Bad Harness"); code == 0 || !strings.Contains(errOut, "not a name such as claude-code") {
		t.Errorf("an invalid harness: %d %s", code, errOut)
	}
}

// S-0105: the dashboard replaces the default whole and commits it, as its
// other writes are committed.
func TestTheDefaultAgentIsReplacedAndCommitted(t *testing.T) {
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	root := tempProject(t)
	for _, c := range [][]string{{"init", "-q", "-b", "main"}, {"config", "user.email", "t@example.invalid"}, {"config", "user.name", "T"}, {"add", "-A"}, {"commit", "-qm", "setup"}} {
		if out, err := exec.Command("git", append([]string{"-C", root}, c...)...).CombinedOutput(); err != nil {
			t.Fatalf("git %v: %s", c, out)
		}
	}
	runIn(t, root, "agent", "set", "--harness", "claude-code", "--model", "claude-opus-5-5", "--config", "effort=high")
	out, errOut, code := runIn(t, root, "agent", "set", "--replace", "--model", "claude-sonnet-5", "--autocommit", "--trailer", "Co-Authored-By: flaiover <flaiover@localhost>", "--json")
	if code != 0 || !strings.Contains(out, `"model": "claude-sonnet-5"`) || strings.Contains(out, "claude-code") || strings.Contains(out, "effort") || !strings.Contains(out, `"commit"`) {
		t.Fatalf("replace: %d %s %s", code, out, errOut)
	}
	log, _ := exec.Command("git", "-C", root, "log", "-1", "--format=%s%n%b", "--name-only").CombinedOutput()
	if !strings.Contains(string(log), "chore: set the project's default agent") || !strings.Contains(string(log), "Co-Authored-By: flaiover") || !strings.Contains(string(log), "system-flow.yaml") {
		t.Errorf("commit: %s", log)
	}
	if _, _, code := runIn(t, root, "agent", "clear", "--autocommit"); code != 0 {
		t.Fatal("clear")
	}
	log, _ = exec.Command("git", "-C", root, "log", "-1", "--format=%s").CombinedOutput()
	if !strings.Contains(string(log), "clear the project's default agent") {
		t.Errorf("clear commit: %s", log)
	}
}
