package cmd

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/execx"
)

func TestAcceptEndToEnd(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	root := tempProject(t)
	_ = os.MkdirAll(filepath.Join(root, "design", "conventions"), 0o755)
	_ = os.WriteFile(filepath.Join(root, "design", "conventions", "README.md"), []byte("# c\n"), 0o644)
	_ = os.WriteFile(filepath.Join(root, "system-flow.yaml"), []byte("version: 1\nname: t\nkey: t\nlayout:\n  design: design\n  docs: docs\n  wip: wip\nprojects:\n  - name: cli\n    path: cli\n    kind: go\n    tags: [command]\n"), 0o644)
	_ = os.MkdirAll(filepath.Join(root, "cli"), 0o755)
	_ = os.WriteFile(filepath.Join(root, "cli", "main.go"), []byte("package main\n"), 0o644)
	r := execx.System{}
	for _, args := range [][]string{{"init", "-q", "-b", "main"}, {"config", "user.email", "t@t"}, {"config", "user.name", "t"}, {"add", "-A"}, {"commit", "-q", "-m", "init"}} {
		if _, err := r.Run(root, "git", args...); err != nil {
			t.Fatal(err)
		}
	}
	for _, step := range [][]string{
		{"epic", "new", "E"}, {"story", "new", "Ship it", "--epic", "E-001", "--tag", "command"}, {"task", "new", "T", "--story", "S-001"},
	} {
		if _, errOut, code := runIn(t, root, step...); code != 0 {
			t.Fatalf("%v: %s", step, errOut)
		}
	}
	sp := filepath.Join(root, "wip", "kanban", "stories", "S-001-ship-it.md")
	b, _ := os.ReadFile(sp)
	_ = os.WriteFile(sp, []byte(strings.Replace(string(b), "- [ ]\n", "- [x] done\n", 1)), 0o644)
	for _, step := range [][]string{{"move", "S-001", "ready"}, {"move", "S-001", "in-progress"}, {"move", "T-001", "ready"}, {"move", "T-001", "in-progress"}, {"move", "T-001", "done"}, {"move", "S-001", "review"}} {
		if _, errOut, code := runIn(t, root, step...); code != 0 {
			t.Fatalf("%v: %s", step, errOut)
		}
	}
	_ = os.WriteFile(filepath.Join(root, "cli", "feature.go"), []byte("package main\n"), 0o644)
	for _, args := range [][]string{{"add", "-A"}, {"commit", "-q", "-m", "feat: [S-001] ship it"}} {
		if _, err := r.Run(root, "git", args...); err != nil {
			t.Fatal(err)
		}
	}
	out, errOut, code := runIn(t, root, "release", "S-001", "--dry-run")
	if code != 0 || !strings.Contains(out, "cli        0.0.0 -> 0.1.0  minor  delivered") {
		t.Fatalf("release dry run: %d %s %s", code, out, errOut)
	}
	out, errOut, code = runIn(t, root, "accept", "S-001", "--by", "alex", "--no-push", "--trailer", "Co-Authored-By: t <t@t>")
	if code != 0 {
		t.Fatalf("accept: %s\n%s", errOut, out)
	}
	if !strings.Contains(out, "accepted S-001: done, 2 items archived, committed, tagged cli/v0.1.0") {
		t.Errorf("summary: %s", out)
	}
	if _, err := os.Stat(filepath.Join(root, "wip", "archive", "kanban", "stories", "S-001-ship-it.md")); err != nil {
		t.Error("story not archived")
	}
	log, _ := r.Run(root, "git", "log", "-1", "--format=%B")
	if !strings.HasPrefix(log, "chore: [S-001] accept and archive; release cli 0.1.0") || !strings.Contains(log, "Co-Authored-By: t <t@t>") {
		t.Errorf("commit message:\n%s", log)
	}
	st, _ := r.Run(root, "git", "status", "--porcelain")
	if strings.TrimSpace(st) != "" {
		t.Errorf("working tree dirty after accept:\n%s", st)
	}
	tag, _ := r.Run(root, "git", "tag", "-l", "cli/*")
	if strings.TrimSpace(tag) != "cli/v0.1.0" {
		t.Errorf("tag: %q", tag)
	}
	_, errOut, code = runIn(t, root, "accept", "S-001", "--by", "alex", "--no-push")
	if code == 0 || !strings.Contains(errOut, "already done") {
		t.Errorf("second accept must refuse: %s", errOut)
	}
	// the epic is accepted straight from backlog: walked to done, major release
	out, errOut, code = runIn(t, root, "accept", "E-001", "--by", "alex", "--no-push", "--deliver", "cli")
	if code != 0 || !strings.Contains(out, "tagged cli/v1.0.0") {
		t.Errorf("epic accept from backlog: %d %s %s", code, out, errOut)
	}
}
