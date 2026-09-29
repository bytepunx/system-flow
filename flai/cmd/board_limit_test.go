package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// S-0167: flai board limit sets a column's WIP limit in board.md, and flai
// board reads it from there.
func TestBoardLimitCommand(t *testing.T) {
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	t.Setenv("FLAI_AGENT", "")
	root := tempProject(t)
	run := func(args ...string) string {
		t.Helper()
		out, errOut, code := runIn(t, root, args...)
		if code != 0 {
			t.Fatalf("flai %v: %s", args, errOut)
		}
		return out
	}
	limits := func() map[string]int {
		var v struct {
			WIPLimits map[string]int `json:"wip_limits"`
		}
		if err := json.Unmarshal([]byte(run("board", "--json")), &v); err != nil {
			t.Fatal(err)
		}
		return v.WIPLimits
	}

	if out := run("board", "limit", "in-progress", "3"); !strings.Contains(out, "in-progress: WIP limit 3") {
		t.Errorf("the command says what it set, got %q", out)
	}
	if got := limits(); got["in-progress"] != 3 || got["ready"] != 5 || got["review"] != 3 {
		t.Errorf("limits after setting in-progress: %v", got)
	}
	data, err := os.ReadFile(filepath.Join(root, "wip", "kanban", "board.md"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), "  in-progress: 3\n") {
		t.Errorf("board.md does not carry the limit:\n%s", data)
	}

	run("board", "limit", "review", "0")
	if got := limits(); got["review"] != 0 {
		t.Errorf("0 removes the limit: %v", got)
	}
	if out := run("board", "--json", "limit", "ready", "7"); !strings.Contains(out, `"limit": 7`) {
		t.Errorf("--json answers the limit set, got %q", out)
	}

	for _, args := range [][]string{
		{"board", "limit", "backlog", "3"},
		{"board", "limit", "done", "1"},
		{"board", "limit", "--", "ready", "-1"},
		{"board", "limit", "ready", "many"},
	} {
		if _, errOut, code := runIn(t, root, args...); code == 0 {
			t.Errorf("flai %v was not refused", args)
		} else if !strings.Contains(errOut, "WIP limit") {
			t.Errorf("flai %v: the refusal does not say why: %s", args, errOut)
		}
	}
}
