package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestIssueCommands(t *testing.T) {
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	root := tempProject(t)
	_ = os.MkdirAll(filepath.Join(root, "design", "conventions"), 0o755)
	_ = os.WriteFile(filepath.Join(root, "design", "conventions", "README.md"), []byte("# c\n"), 0o644)
	out, errOut, code := runIn(t, root, "issue", "new", "Go not on PATH", "--class", "blocker", "--cost", "3m")
	if code != 0 || !strings.HasPrefix(out, "I-0001 Go not on PATH") {
		t.Fatalf("new: %d %s %s", code, out, errOut)
	}
	_, errOut, code = runIn(t, root, "issue", "new", "no class")
	if code == 0 || !strings.Contains(errOut, "class") {
		t.Fatalf("class required: %s", errOut)
	}
	out, _, code = runIn(t, root, "issue", "bump", "I-0001", "--cost", "5m", "--note", "again")
	if code != 0 || !strings.Contains(out, "count 2, avg cost 4m") {
		t.Fatalf("bump: %s", out)
	}
	sum, _ := os.ReadFile(filepath.Join(root, "design", "issues", "summary.md"))
	if !strings.Contains(string(sum), "| [I-0001](I-0001-go-not-on-path.md) | blocker | Go not on PATH | 2 | 4m | 8m |") {
		t.Errorf("summary after bump:\n%s", sum)
	}
	out, _, _ = runIn(t, root, "prime")
	if !strings.HasSuffix(strings.TrimSpace(out), "design/issues/summary.md") {
		t.Errorf("prime should end with the summary when issues are open:\n%s", out)
	}
	out, _, _ = runIn(t, root, "prime", "--cat")
	if !strings.Contains(out, "open issues\n===") || !strings.Contains(out, "[I-0001]") {
		t.Errorf("prime --cat should include the table:\n%s", out)
	}
	out, _, code = runIn(t, root, "issue", "list", "--json")
	var list []map[string]any
	if code != 0 || json.Unmarshal([]byte(out), &list) != nil || len(list) != 1 || list[0]["count"].(float64) != 2 {
		t.Fatalf("list json: %s", out)
	}
	out, _, code = runIn(t, root, "issue", "close", "I-0001", "--reason", "fixed")
	if code != 0 || !strings.Contains(out, "I-0001 closed") {
		t.Fatalf("close: %s", out)
	}
	out, _, _ = runIn(t, root, "issue", "summary")
	if strings.TrimSpace(out) != "no open issues" {
		t.Errorf("summary after close: %s", out)
	}
	out, _, _ = runIn(t, root, "prime")
	if strings.Contains(out, "summary.md") {
		t.Error("prime should not list the summary when nothing is open")
	}
	out, _, _ = runIn(t, root, "check")
	for _, line := range strings.Split(out, "\n") {
		if strings.Contains(line, "issues.") {
			t.Errorf("check reports an issues finding after the commands: %s", line)
		}
	}
}
