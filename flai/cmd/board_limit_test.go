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
	t.Parallel()
	useConfig(t, filepath.Join(t.TempDir(), "cfg.json"))
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

// boardLimitBaseline is the baseline workflow file flai carries, as the
// tests in this package read it from the source tree.
const boardLimitBaseline = "../internal/workflow/baseline/board.md"

// boardLimitTestingWorkflow is the baseline workflow file with a state
// testing, of the active kind and with limit: true, added between
// in-progress and review, the moves into and out of it, and its cancel.
func boardLimitTestingWorkflow(t *testing.T) string {
	t.Helper()
	data, err := os.ReadFile(boardLimitBaseline)
	if err != nil {
		t.Fatal(err)
	}
	text := string(data)
	for _, r := range [][2]string{
		{"  - {name: in-progress, kind: active, limit: true}\n",
			"  - {name: in-progress, kind: active, limit: true}\n  - {name: testing, kind: active, limit: true}\n"},
		{"transitions:\n",
			"transitions:\n  - {from: in-progress, to: testing}\n  - {from: testing, to: review}\n"},
		{"{from: [backlog, ready, in-progress], to: cancelled,", "{from: [backlog, ready, in-progress, testing], to: cancelled,"},
	} {
		if !strings.Contains(text, r[0]) {
			t.Fatalf("the baseline workflow file no longer has %q", r[0])
		}
		text = strings.Replace(text, r[0], r[1], 1)
	}
	return text
}

// I-0380 instance 4: in a project whose workflow file adds a state testing
// with limit: true, flai board limit sets testing's limit, a later write of
// the board by flai order keeps it, and a state without limit: true, or one
// the workflow does not have, is refused with the workflow's limited states.
func TestBoardLimitTakesAStateTheWorkflowAdds(t *testing.T) {
	t.Parallel()
	useConfig(t, filepath.Join(t.TempDir(), "cfg.json"))
	root := tempProject(t)
	file := filepath.Join(root, "design", "workflows", "board.md")
	if err := os.MkdirAll(filepath.Dir(file), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(file, []byte(boardLimitTestingWorkflow(t)), 0o644); err != nil {
		t.Fatal(err)
	}
	run := func(args ...string) string {
		t.Helper()
		out, errOut, code := runIn(t, root, args...)
		if code != 0 {
			t.Fatalf("flai %v: %s", args, errOut)
		}
		return out
	}
	boardFile := filepath.Join(root, "wip", "kanban", "board.md")
	limitIn := func(when string) {
		t.Helper()
		data, err := os.ReadFile(boardFile)
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(string(data), "  testing: 2\n") {
			t.Errorf("%s, board.md does not carry testing's limit:\n%s", when, data)
		}
	}

	if out := run("board", "limit", "testing", "2"); !strings.Contains(out, "testing: WIP limit 2") {
		t.Errorf("the command says what it set, got %q", out)
	}
	limitIn("after flai board limit")

	run("story", "new", "One")
	run("story", "new", "Two")
	run("order", "S-0002", "--top")
	limitIn("after flai order wrote the board")
	data, err := os.ReadFile(boardFile)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), "order:\n  - S-0002\n") {
		t.Errorf("flai order did not write the board:\n%s", data)
	}

	for _, state := range []string{"backlog", "done", "staging"} {
		_, errOut, code := runIn(t, root, "board", "limit", state, "1")
		if code == 0 {
			t.Errorf("flai board limit %s was not refused", state)
		} else if want := state + " has no WIP limit; only ready, in-progress, testing, review do"; !strings.Contains(errOut, want) {
			t.Errorf("flai board limit %s: the refusal %q does not say %q", state, errOut, want)
		}
	}
	limitIn("after the refusals")
}
