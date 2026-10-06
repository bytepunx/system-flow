package storygit

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/bytepunx/system-flow/flai/internal/execx"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// What flai stream sync reads of a story's worktree before and after it
// rebases (ADR-0069): agents never rebase by hand, so sync refuses a worktree
// it could leave in a state that loses work, and names what stopped it.

// Uncommitted lists the paths in the worktree at dir with uncommitted
// changes, tracked and untracked, as git status --porcelain names them (an
// untracked folder is one path).
func Uncommitted(r execx.Runner, dir string) ([]string, error) {
	st, err := r.Run(dir, "git", "status", "--porcelain")
	if err != nil {
		return nil, err
	}
	return workitem.PorcelainPaths(st), nil
}

// RebaseInProgress reports whether a rebase has stopped in the worktree at
// dir and waits to be continued or aborted.
func RebaseInProgress(r execx.Runner, dir string) bool {
	for _, state := range []string{"rebase-merge", "rebase-apply"} {
		p, err := r.Run(dir, "git", "rev-parse", "--git-path", state)
		if err != nil || p == "" {
			continue
		}
		if !filepath.IsAbs(p) {
			p = filepath.Join(dir, p)
		}
		if _, err := os.Stat(p); err == nil {
			return true
		}
	}
	return false
}

// ContinueRebase continues the rebase stopped in the worktree at dir once
// its conflicts are resolved and staged, keeping each commit's message with
// no editor (ADR-0098). It returns git's error when the rebase stops again,
// which RebaseInProgress then reports waiting.
func ContinueRebase(r execx.Runner, dir string) error {
	_, err := r.Run(dir, "git", "-c", "core.editor=true", "rebase", "--continue")
	return err
}

// Conflicts lists the paths left unmerged in the worktree at dir, which a
// stopped rebase waits for an agent to resolve and git add.
func Conflicts(r execx.Runner, dir string) []string {
	out, err := r.Run(dir, "git", "diff", "--name-only", "--diff-filter=U")
	if err != nil {
		return nil
	}
	var paths []string
	for _, f := range strings.Split(out, "\n") {
		if f = strings.TrimSpace(f); f != "" {
			paths = append(paths, f)
		}
	}
	return paths
}
