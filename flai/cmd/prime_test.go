package cmd

import (
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"
)

func TestPrime(t *testing.T) {
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	dir := "../internal/metrics/testdata/good"
	out, errOut, code := runIn(t, dir, "prime")
	if code != 0 {
		t.Fatalf("prime: %s", errOut)
	}
	want := "design/conventions/README.md\ndesign/conventions/session-start.md\ndesign/conventions/git.md\n"
	if out != want {
		t.Errorf("paths:\n%s", out)
	}
	out, _, _ = runIn(t, dir, "prime", "--cat")
	if !strings.Contains(out, "design/conventions/git.md\n===") || !strings.Contains(out, "- Commit at landing.") || strings.Index(out, "Session start") > strings.Index(out, "# Git") {
		t.Errorf("cat:\n%s", out)
	}
	out, _, _ = runIn(t, dir, "prime", "--json")
	var v struct {
		Files []struct {
			Name  string `json:"name"`
			Order int    `json:"order"`
		} `json:"files"`
	}
	if err := json.Unmarshal([]byte(out), &v); err != nil || len(v.Files) != 2 || v.Files[1].Order != 20 {
		t.Errorf("json: %v %s", err, out)
	}
	_, errOut, code = runIn(t, tempProject(t), "prime")
	if code == 0 || !strings.Contains(errOut, "no conventions folder") {
		t.Errorf("missing folder should fail: %d %s", code, errOut)
	}
}
