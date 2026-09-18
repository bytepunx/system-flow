package check

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/gitver"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

func TestGitCompat(t *testing.T) {
	const plain = "[core]\n\trepositoryformatversion = 0\n"
	const extended = "[core]\n\trepositoryformatversion = 1\n[extensions]\n\trelativeWorktrees = true\n"
	oldGit, newGit := gitver.Version{Major: 2, Minor: 47}, gitver.Version{Major: 2, Minor: 54}
	for _, c := range []struct {
		name, config string
		git          gitver.Version
		line         int
	}{
		{"extension absent, old git", plain, oldGit, 0},
		{"extension set, new git", extended, newGit, 0},
		{"extension set, old git", extended, oldGit, 4},
		{"extension off, old git", "[extensions]\n\trelativeWorktrees = false\n", oldGit, 0},
		{"same key in another section", "[other]\n\trelativeWorktrees = true\n", oldGit, 0},
	} {
		t.Run(c.name, func(t *testing.T) {
			root := t.TempDir()
			_ = os.MkdirAll(filepath.Join(root, ".git"), 0o755)
			_ = os.WriteFile(filepath.Join(root, ".git", "config"), []byte(c.config), 0o644)
			res := &Result{}
			GitCompat(res, &workitem.Repo{Root: root}, c.git)
			if c.line == 0 {
				if len(res.Findings) != 0 || res.Warnings != 0 {
					t.Fatalf("unexpected finding: %+v", res.Findings)
				}
				return
			}
			if len(res.Findings) != 1 || res.Warnings != 1 || res.Errors != 0 {
				t.Fatalf("want one warning, got %+v", res)
			}
			f := res.Findings[0]
			if f.Rule != "git.relative-worktrees" || f.Level != Warning || f.Line != c.line || f.Path != filepath.Join(".git", "config") {
				t.Errorf("finding: %+v", f)
			}
			for _, want := range []string{"2.47", "2.48", "git worktree repair", "worktrees.relative_paths"} {
				if !strings.Contains(f.Message, want) {
					t.Errorf("message should mention %q: %s", want, f.Message)
				}
			}
		})
	}
	// no repository at all: nothing to say
	res := &Result{}
	GitCompat(res, &workitem.Repo{Root: t.TempDir()}, oldGit)
	if len(res.Findings) != 0 {
		t.Errorf("no .git: %+v", res.Findings)
	}
}
