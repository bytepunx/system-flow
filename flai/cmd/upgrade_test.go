package cmd

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/execx"
)

func TestUpgradeCommand(t *testing.T) {
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	dest := filepath.Join(t.TempDir(), "proj")
	if _, errOut, code := runIn(t, ".", "new", dest, "--template", miniTemplate, "--defaults", "--no-git"); code != 0 {
		t.Fatalf("new: %s", errOut)
	}
	if _, err := os.Stat(filepath.Join(dest, "system-flow.lock.yaml")); err != nil {
		t.Fatal("flai new must write the lock")
	}
	// a newer template
	nt := t.TempDir()
	_ = filepath.WalkDir(miniTemplate, func(path string, d os.DirEntry, err error) error {
		rel, _ := filepath.Rel(miniTemplate, path)
		if d.IsDir() {
			return os.MkdirAll(filepath.Join(nt, rel), 0o755)
		}
		b, _ := os.ReadFile(path)
		return os.WriteFile(filepath.Join(nt, rel), b, 0o644)
	})
	y, _ := os.ReadFile(filepath.Join(nt, "template.yaml"))
	_ = os.WriteFile(filepath.Join(nt, "template.yaml"), []byte(strings.Replace(string(y), "version: 9.9.9", "version: 10.0.0", 1)), 0o644)
	_ = os.WriteFile(filepath.Join(nt, "root", "plain.txt"), []byte("verbatim v2\n"), 0o644)
	_ = os.WriteFile(filepath.Join(nt, "root", "new.txt"), []byte("new\n"), 0o644)
	_ = os.WriteFile(filepath.Join(dest, "plain.txt"), []byte("edited\n"), 0o644)

	out, _, code := runIn(t, dest, "upgrade", "--template", miniTemplate)
	if code != 0 || !strings.Contains(out, "already at template 9.9.9") {
		t.Fatalf("same version: %d %s", code, out)
	}
	out, _, code = runIn(t, dest, "upgrade", "--template", nt, "--dry-run")
	if code != 0 || !strings.Contains(out, "9.9.9 -> 10.0.0") || !strings.Contains(out, "conflict  plain.txt") || !strings.Contains(out, "add       new.txt") || !strings.Contains(out, "dry run") {
		t.Fatalf("dry run: %d\n%s", code, out)
	}
	if _, err := os.Stat(filepath.Join(dest, "new.txt")); err == nil {
		t.Fatal("dry run wrote files")
	}
	_, errOut, code := runIn(t, dest, "upgrade", "--template", nt)
	if code != 1 || !strings.Contains(errOut, "conflict(s) need a decision") {
		t.Fatalf("non-interactive without policy must refuse: %d %s", code, errOut)
	}
	if _, err := os.Stat(filepath.Join(dest, "new.txt")); err == nil {
		t.Fatal("refused upgrade wrote files")
	}
	out, errOut, code = runIn(t, dest, "upgrade", "--template", nt, "--keep-all")
	if code != 0 || !strings.Contains(out, "1 conflicts kept") {
		t.Fatalf("keep-all: %d %s %s", code, out, errOut)
	}
	if b, _ := os.ReadFile(filepath.Join(dest, "plain.txt")); string(b) != "edited\n" {
		t.Error("kept file overwritten")
	}
	if _, err := os.Stat(filepath.Join(dest, "new.txt")); err != nil {
		t.Error("new file not added")
	}
	mf, _ := os.ReadFile(filepath.Join(dest, "system-flow.yaml"))
	if !strings.Contains(string(mf), `version: "10.0.0"`) {
		t.Errorf("manifest version not updated:\n%s", mf)
	}
	out, _, code = runIn(t, dest, "upgrade", "--template", nt)
	if code != 0 || !strings.Contains(out, "already at template 10.0.0") {
		t.Errorf("after upgrade: %d %s", code, out)
	}
	// dirty tree guard with git
	if _, err := exec.LookPath("git"); err == nil {
		r := execx.System{}
		for _, args := range [][]string{{"init", "-q", "-b", "main"}, {"config", "user.email", "t@t"}, {"config", "user.name", "t"}, {"add", "-A"}, {"commit", "-q", "-m", "x"}} {
			_, _ = r.Run(dest, "git", args...)
		}
		_ = os.WriteFile(filepath.Join(dest, "dirty.txt"), []byte("x"), 0o644)
		_, errOut, code = runIn(t, dest, "upgrade", "--template", nt, "--force", "--keep-all")
		if code != 0 {
			t.Errorf("--force should allow a dirty tree: %s", errOut)
		}
		_, _ = r.Run(dest, "git", "add", "-A")
		_, _ = r.Run(dest, "git", "commit", "-q", "-m", "y")
		_ = os.WriteFile(filepath.Join(dest, "dirty2.txt"), []byte("x"), 0o644)
		_ = os.WriteFile(filepath.Join(nt, "template.yaml"), []byte(strings.Replace(string(y), "version: 9.9.9", "version: 11.0.0", 1)), 0o644)
		_, errOut, code = runIn(t, dest, "upgrade", "--template", nt, "--keep-all")
		if code == 0 || !strings.Contains(errOut, "uncommitted changes") {
			t.Errorf("dirty tree must be refused: %d %s", code, errOut)
		}
	}
	// relock on a project without a lock
	_ = os.Remove(filepath.Join(dest, "system-flow.lock.yaml"))
	out, _, code = runIn(t, dest, "upgrade", "--template", nt, "--relock", "--force")
	if code != 0 || !strings.Contains(out, "relocked") {
		t.Errorf("relock: %d %s", code, out)
	}
}
