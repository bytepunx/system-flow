package cmd

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestMigrateIDs(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	t.Setenv("FLAI_AGENT", "tester")
	root := tempProject(t)
	git := func(args ...string) {
		t.Helper()
		c := exec.Command("git", args...)
		c.Dir = root
		c.Env = append(os.Environ(), "GIT_AUTHOR_NAME=t", "GIT_AUTHOR_EMAIL=t@t", "GIT_COMMITTER_NAME=t", "GIT_COMMITTER_EMAIL=t@t")
		if out, err := c.CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v %s", args, err, out)
		}
	}
	git("init", "-q")
	legacy := "---\nid: S-003\ntype: story\nnature: feature\ntitle: Legacy\nstatus: backlog\nowner: a\ncreated: 2026-09-15T21:00:00Z\nupdated: 2026-09-15T21:00:00Z\ntransitions: []\ntags: []\n---\n\n# S-003 Legacy\n\n## Goal\nx\n"
	_ = os.WriteFile(filepath.Join(root, "wip/kanban/stories/S-003-legacy.md"), []byte(legacy), 0o644)
	_ = os.WriteFile(filepath.Join(root, "README.md"), []byte("Start with S-003.\n"), 0o644)
	git("add", "-A")
	git("commit", "-q", "-m", "legacy")
	untracked := strings.ReplaceAll(strings.ReplaceAll(legacy, "S-003", "S-004"), "Legacy", "Untracked")
	_ = os.WriteFile(filepath.Join(root, "wip/kanban/stories/S-004-untracked.md"), []byte(untracked), 0o644)

	out, errOut, code := runIn(t, root, "migrate", "ids", "--dry-run")
	if code != 0 || !strings.Contains(out, "would migrate 2 item(s), 2 file move(s), 3 file(s) rewritten") || !strings.Contains(out, "mv wip/kanban/stories/S-003-legacy.md -> wip/kanban/stories/S-0003-legacy.md") {
		t.Fatalf("dry run: %d %s %s", code, out, errOut)
	}
	if _, err := os.Stat(filepath.Join(root, "wip/kanban/stories/S-003-legacy.md")); err != nil {
		t.Fatal("dry run must not move files")
	}
	out, errOut, code = runIn(t, root, "migrate", "ids")
	if code != 0 || !strings.Contains(out, "migrated 2 item(s)") {
		t.Fatalf("migrate: %d %s %s", code, out, errOut)
	}
	status := exec.Command("git", "status", "--porcelain")
	status.Dir = root
	st, _ := status.Output()
	if !strings.Contains(string(st), "wip/kanban/stories/S-003-legacy.md -> wip/kanban/stories/S-0003-legacy.md") {
		t.Errorf("expected a git rename, got:\n%s", st)
	}
	if _, err := os.Stat(filepath.Join(root, "wip/kanban/stories/S-0004-untracked.md")); err != nil {
		t.Error("untracked file should be renamed without git")
	}
	readme, _ := os.ReadFile(filepath.Join(root, "README.md"))
	if string(readme) != "Start with S-0003.\n" {
		t.Errorf("readme: %s", readme)
	}
	out, _, code = runIn(t, root, "show", "S-3")
	if code != 0 || !strings.Contains(out, "S-0003 Legacy") {
		t.Errorf("show by short id: %d %s", code, out)
	}
	out, _, _ = runIn(t, root, "migrate", "ids")
	if !strings.Contains(out, "nothing to migrate") {
		t.Errorf("second run: %s", out)
	}
}
