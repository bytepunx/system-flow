package storygit

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/execx"
)

// Integration: these run real git.
func gitRepo(t *testing.T) (execx.Runner, string) {
	t.Helper()
	if testing.Short() {
		t.Skip("integration: runs real git")
	}
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}
	r, dir := execx.System{}, t.TempDir()
	for _, args := range [][]string{
		{"init", "-q", "-b", "main"},
		{"config", "user.email", "t@t"},
		{"config", "user.name", "t"},
	} {
		if out, err := r.Run(dir, "git", args...); err != nil {
			t.Fatal(out, err)
		}
	}
	write(t, dir, "a.md", "one\n")
	git(t, r, dir, "add", "-A")
	git(t, r, dir, "commit", "-q", "-m", "init")
	return r, dir
}

func write(t *testing.T, dir, rel, content string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, rel), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func git(t *testing.T, r execx.Runner, dir string, args ...string) {
	t.Helper()
	if out, err := r.Run(dir, "git", args...); err != nil {
		t.Fatal(out, err)
	}
}

func TestUncommittedRebaseInProgressAndConflicts(t *testing.T) {
	r, dir := gitRepo(t)
	if got, err := Uncommitted(r, dir); err != nil || len(got) != 0 {
		t.Fatalf("clean: %v %v", got, err)
	}
	if RebaseInProgress(r, dir) || len(Conflicts(r, dir)) != 0 {
		t.Fatal("clean repository reports a rebase")
	}

	write(t, dir, "a.md", "changed\n")
	write(t, dir, "new.md", "new\n")
	if got, err := Uncommitted(r, dir); err != nil || strings.Join(got, ",") != "a.md,new.md" {
		t.Fatalf("uncommitted: %v %v", got, err)
	}
	git(t, r, dir, "checkout", "-q", "--", "a.md")
	_ = os.Remove(filepath.Join(dir, "new.md"))

	// a branch and main change the same line, so the rebase stops
	git(t, r, dir, "checkout", "-q", "-b", "story/S-0001")
	write(t, dir, "a.md", "story\n")
	git(t, r, dir, "commit", "-q", "-am", "story")
	git(t, r, dir, "checkout", "-q", "main")
	write(t, dir, "a.md", "main\n")
	git(t, r, dir, "commit", "-q", "-am", "main")
	git(t, r, dir, "checkout", "-q", "story/S-0001")
	if _, err := r.Run(dir, "git", "rebase", "main"); err == nil {
		t.Fatal("rebase did not stop")
	}
	if !RebaseInProgress(r, dir) {
		t.Error("stopped rebase not reported")
	}
	if got := Conflicts(r, dir); strings.Join(got, ",") != "a.md" {
		t.Errorf("conflicts: %v", got)
	}
	git(t, r, dir, "rebase", "--abort")
	if RebaseInProgress(r, dir) {
		t.Error("aborted rebase still reported")
	}
}
