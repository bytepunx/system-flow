package cmd

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/execx"
)

func TestTemplatePushCommand(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	r := execx.System{}
	bare := filepath.Join(t.TempDir(), "tpl.git")
	if _, err := r.Run("", "git", "init", "-q", "--bare", "-b", "main", bare); err != nil {
		t.Fatal(err)
	}
	dir := t.TempDir()
	_ = filepath.WalkDir(miniTemplate, func(path string, d os.DirEntry, err error) error {
		rel, _ := filepath.Rel(miniTemplate, path)
		if d.IsDir() {
			return os.MkdirAll(filepath.Join(dir, rel), 0o755)
		}
		b, _ := os.ReadFile(path)
		return os.WriteFile(filepath.Join(dir, rel), b, 0o644)
	})
	y, _ := os.ReadFile(filepath.Join(dir, "template.yaml"))
	_ = os.WriteFile(filepath.Join(dir, "template.yaml"), []byte(string(y)+"publish:\n  repo: "+bare+"\n  ref: main\n"), 0o644)

	out, errOut, code := runCLI(t, "template", "push", dir, "--dry-run")
	if code != 0 || !strings.Contains(out, "would push") || !strings.Contains(out, "dry run: remote untouched") {
		t.Fatalf("dry run: %d %s %s", code, out, errOut)
	}
	out, errOut, code = runCLI(t, "template", "push", dir, "--tag")
	if code != 0 || !strings.Contains(out, "pushed 9.9.9 to "+bare+" main") || !strings.Contains(out, "tag v9.9.9") || !strings.Contains(out, "branch created") {
		t.Fatalf("push: %d %s %s", code, out, errOut)
	}
	out, _, code = runCLI(t, "template", "push", dir)
	if code != 0 || !strings.Contains(out, "nothing to push") {
		t.Errorf("second push: %d %s", code, out)
	}
	_, errOut, code = runCLI(t, "template", "push", dir, "--tag")
	if code == 0 || !strings.Contains(errOut, "already exists") {
		t.Errorf("tag twice must fail: %d %s", code, errOut)
	}
	_, errOut, code = runCLI(t, "template", "push")
	if code == 0 || !strings.Contains(errOut, "not a local directory") {
		t.Errorf("config remote template must ask for a dir: %d %s", code, errOut)
	}
}
