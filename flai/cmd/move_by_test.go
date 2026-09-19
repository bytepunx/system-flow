package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// S-0058: an agent's moves carry its name, so MCP does not report them back
// to it as the designer's.
func TestMoveRecordsWhoMovedIt(t *testing.T) {
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	t.Setenv("FLAI_AGENT", "")
	root := tempProject(t)
	for _, args := range [][]string{{"config", "set", "author", "olive"}, {"epic", "new", "Epic"}, {"story", "new", "Moved", "--epic", "E-0001"}} {
		if _, errOut, code := runIn(t, root, args...); code != 0 {
			t.Fatal(errOut)
		}
	}
	file := filepath.Join(root, "wip/kanban/stories/S-0001-moved.md")
	s, _ := os.ReadFile(file)
	_ = os.WriteFile(file, []byte(strings.Replace(string(s), "## Acceptance criteria\n- [ ]\n", "## Acceptance criteria\n- [ ] ok\n", 1)), 0o644)
	lastBy := func() string {
		data, _ := os.ReadFile(file)
		lines := strings.Split(string(data), "\n")
		by := ""
		for _, l := range lines {
			if strings.HasPrefix(strings.TrimSpace(l), "by: ") {
				by = strings.TrimPrefix(strings.TrimSpace(l), "by: ")
			}
		}
		return by
	}
	if _, errOut, code := runIn(t, root, "move", "S-0001", "ready"); code != 0 {
		t.Fatal(errOut)
	}
	if got := lastBy(); got != "olive" {
		t.Errorf("without FLAI_AGENT the config author moved it, got %q", got)
	}
	t.Setenv("FLAI_AGENT", "claude")
	if _, errOut, code := runIn(t, root, "move", "S-0001", "in-progress"); code != 0 {
		t.Fatal(errOut)
	}
	if got := lastBy(); got != "claude" {
		t.Errorf("with FLAI_AGENT the agent moved it, got %q", got)
	}
	if _, errOut, code := runIn(t, root, "move", "S-0001", "cancelled", "--reason", "x", "--by", "dana"); code != 0 {
		t.Fatal(errOut)
	}
	if got := lastBy(); got != "dana" {
		t.Errorf("--by wins, got %q", got)
	}
}
