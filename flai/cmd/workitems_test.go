package cmd

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// runIn runs the CLI with the working directory and clock fixed.
func runIn(t *testing.T, dir string, args ...string) (string, string, int) {
	t.Helper()
	var out, errOut bytes.Buffer
	a := &app{out: &out, errOut: &errOut, cwd: dir, clock: func() time.Time { return time.Date(2026, 9, 15, 21, 0, 0, 0, time.UTC) }}
	root := newRootCmdWith(a)
	root.SetArgs(args)
	code := 0
	if err := root.Execute(); err != nil {
		errOut.WriteString("flai: " + err.Error() + "\n")
		code = 1
	}
	return out.String(), errOut.String(), code
}

func tempProject(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	_ = os.WriteFile(filepath.Join(root, "system-flow.yaml"), []byte("version: 1\nname: t\nkey: t\nlayout:\n  design: design\n  docs: docs\n  wip: wip\n"), 0o644)
	for _, d := range []string{"wip/kanban/epics", "wip/kanban/stories", "wip/kanban/tasks", "wip/agents"} {
		_ = os.MkdirAll(filepath.Join(root, d), 0o755)
	}
	return root
}

func TestWorkItemLifecycle(t *testing.T) {
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	t.Setenv("FLAI_AGENT", "tester")
	root := tempProject(t)
	nested := filepath.Join(root, "wip", "kanban")

	out, errOut, code := runIn(t, nested, "epic", "new", "Big thing", "--nature", "research")
	if code != 0 || !strings.HasPrefix(out, "E-001 Big thing") {
		t.Fatalf("epic new: %d %s %s", code, out, errOut)
	}
	_, errOut, code = runIn(t, nested, "story", "new", "Slice")
	if code == 0 || !strings.Contains(errOut, "epic") {
		t.Fatalf("story without --epic should fail: %s", errOut)
	}
	out, _, code = runIn(t, nested, "story", "new", "Slice", "--epic", "E-001", "--tag", "a,b")
	if code != 0 || !strings.HasPrefix(out, "S-001 Slice") {
		t.Fatalf("story new: %s", out)
	}
	out, _, _ = runIn(t, nested, "task", "new", "Piece", "--story", "S-001", "--json")
	var created map[string]any
	if err := json.Unmarshal([]byte(out), &created); err != nil || created["id"] != "T-001" {
		t.Fatalf("task new --json: %v %s", err, out)
	}

	// story needs real criteria before ready
	storyPath := filepath.Join(root, "wip", "kanban", "stories", "S-001-slice.md")
	b, _ := os.ReadFile(storyPath)
	_ = os.WriteFile(storyPath, bytes.Replace(b, []byte("- [ ]\n"), []byte("- [ ] it works\n"), 1), 0o644)

	for _, step := range [][]string{
		{"move", "S-001", "ready"}, {"move", "S-001", "in-progress"},
		{"stream", "open", "S-001"}, {"stream", "log", "S-001", "started work"},
		{"move", "T-001", "ready"}, {"move", "T-001", "in-progress"}, {"block", "T-001", "--reason", "waiting"},
		{"unblock", "T-001"}, {"move", "T-001", "done"}, {"move", "S-001", "review"},
	} {
		if _, errOut, code := runIn(t, nested, step...); code != 0 {
			t.Fatalf("%v: %s", step, errOut)
		}
	}
	_, errOut, code = runIn(t, nested, "move", "S-001", "done")
	if code == 0 || !strings.Contains(errOut, "unchecked") {
		t.Fatalf("done with unchecked criteria should fail: %s", errOut)
	}
	b, _ = os.ReadFile(storyPath)
	_ = os.WriteFile(storyPath, bytes.Replace(b, []byte("- [ ] it works"), []byte("- [x] it works"), 1), 0o644)
	if _, errOut, code = runIn(t, nested, "move", "S-001", "done", "--by", "alex"); code != 0 {
		t.Fatalf("done: %s", errOut)
	}

	out, _, _ = runIn(t, nested, "show", "S-001")
	if !strings.Contains(out, "story · feature · done") || !strings.Contains(out, "T-001  done") || !strings.Contains(out, "by alex") {
		t.Errorf("show:\n%s", out)
	}
	out, _, _ = runIn(t, nested, "board", "--all")
	if !strings.Contains(out, "done") || !strings.Contains(out, "S-001") || !strings.Contains(out, "E-001") {
		t.Errorf("board:\n%s", out)
	}
	idx, _ := os.ReadFile(filepath.Join(root, "wip", "agents", "index.md"))
	if !strings.Contains(string(idx), "| [S-001](S-001.md) | Slice | done | tester |") {
		t.Errorf("index:\n%s", idx)
	}

	out, _, code = runIn(t, nested, "archive", "--dry-run")
	if code != 0 || !strings.Contains(out, "would archive S-001") || !strings.Contains(out, "would archive T-001") || !strings.Contains(out, "narrative") {
		t.Fatalf("archive dry run: %s", out)
	}
	if _, err := os.Stat(storyPath); err != nil {
		t.Fatal("dry run moved files")
	}
	out, _, code = runIn(t, nested, "archive")
	if code != 0 || !strings.Contains(out, "archived S-001") {
		t.Fatalf("archive: %s", out)
	}
	if _, err := os.Stat(filepath.Join(root, "wip", "archive", "agents", "S-001.md")); err != nil {
		t.Error("narrative not archived")
	}
	idx, _ = os.ReadFile(filepath.Join(root, "wip", "agents", "index.md"))
	if strings.Contains(string(idx), "S-001") {
		t.Error("index still lists archived stream")
	}
	out, _, _ = runIn(t, nested, "show", "S-001")
	if !strings.Contains(out, "archived") {
		t.Errorf("show archived:\n%s", out)
	}
}

func TestCommandsOutsideProject(t *testing.T) {
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	_, errOut, code := runIn(t, t.TempDir(), "board")
	if code == 0 || !strings.Contains(errOut, "system-flow.yaml") {
		t.Fatalf("expected manifest error: %s", errOut)
	}
}
