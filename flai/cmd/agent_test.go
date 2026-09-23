package cmd

import (
	"os"
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
