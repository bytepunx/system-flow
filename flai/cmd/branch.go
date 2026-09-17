package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// Story branches (ADR-0019): each story is worked on story/S-nnnn in a git
// worktree under .flai-cache/worktrees; wip/ stays in the main checkout so
// the board is live. stream open creates the branch, stream sync rebases it
// onto the main branch, accept merges it and removes the worktree.

const storyBranchPrefix = "story/"

func storyBranch(id string) string { return storyBranchPrefix + id }

// inGitWorkTree reports whether root is inside a git work tree.
func (a *app) inGitWorkTree(root string) bool {
	if _, err := a.runner.LookPath("git"); err != nil {
		return false
	}
	_, err := a.runner.Run(root, "git", "rev-parse", "--is-inside-work-tree")
	return err == nil
}

// mainBranch is the branch checked out in the main checkout.
func (a *app) mainBranch(mainRoot string) (string, error) {
	out, err := a.runner.Run(mainRoot, "git", "rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		return "", err
	}
	if out == "HEAD" {
		return "", fmt.Errorf("the main checkout has a detached HEAD; check out a branch first")
	}
	return out, nil
}

func (a *app) branchExists(root, branch string) bool {
	_, err := a.runner.Run(root, "git", "rev-parse", "--verify", "--quiet", "refs/heads/"+branch)
	return err == nil
}

// openStoryBranch creates story/<id> from the main branch and checks it out
// in a worktree. Returns the branch and path, or "" when git is not in use.
func (a *app) openStoryBranch(repo *workitem.Repo, id string) (branch, path string, created bool, err error) {
	if !a.inGitWorkTree(repo.MainRoot) {
		return "", "", false, nil
	}
	branch, path = storyBranch(id), repo.WorktreePath(id)
	if _, err := os.Stat(path); err == nil {
		return branch, path, false, nil
	}
	base, err := a.mainBranch(repo.MainRoot)
	if err != nil {
		return "", "", false, err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return "", "", false, err
	}
	if a.branchExists(repo.MainRoot, branch) {
		_, err = a.runner.Run(repo.MainRoot, "git", "worktree", "add", "--quiet", path, branch)
	} else {
		_, err = a.runner.Run(repo.MainRoot, "git", "worktree", "add", "--quiet", "-b", branch, path, base)
	}
	if err != nil {
		return "", "", false, err
	}
	a.logger().Info("story branch opened", "component", "git", "branch", branch, "worktree", relPath(repo.MainRoot, path), "base", base)
	return branch, path, true, nil
}

// syncStoryBranch rebases the story branch onto the main branch inside its
// worktree. On conflicts the rebase is left in progress for the agent to
// resolve, and the conflicting files are returned with the error.
func (a *app) syncStoryBranch(repo *workitem.Repo, id string) (base string, conflicts []string, err error) {
	path := repo.WorktreePath(id)
	if _, err := os.Stat(path); err != nil {
		return "", nil, fmt.Errorf("%s has no worktree at %s; open one with flai stream open %s", id, relPath(repo.MainRoot, path), id)
	}
	base, err = a.mainBranch(repo.MainRoot)
	if err != nil {
		return "", nil, err
	}
	if _, err := a.runner.Run(path, "git", "rebase", "--autostash", base); err != nil {
		out, _ := a.runner.Run(path, "git", "diff", "--name-only", "--diff-filter=U")
		for _, f := range strings.Split(out, "\n") {
			if f = strings.TrimSpace(f); f != "" {
				conflicts = append(conflicts, f)
			}
		}
		if len(conflicts) > 0 {
			return base, conflicts, fmt.Errorf("rebase of %s onto %s stopped on conflicts in %s; resolve them in %s, then `git rebase --continue` there and run flai stream sync %s again (or `git rebase --abort`)", storyBranch(id), base, strings.Join(conflicts, ", "), relPath(repo.MainRoot, path), id)
		}
		return base, nil, err
	}
	return base, nil, nil
}

// mergeStoryBranch brings a story's branch into the main branch: the
// worktree must be clean, the branch is rebased, fast-forwarded into main,
// and the worktree and branch are removed. Returns false when the story has
// no branch.
func (a *app) mergeStoryBranch(repo *workitem.Repo, id string) (bool, error) {
	if !a.inGitWorkTree(repo.MainRoot) || !a.branchExists(repo.MainRoot, storyBranch(id)) {
		return false, nil
	}
	branch, path := storyBranch(id), repo.WorktreePath(id)
	hasWorktree := false
	if _, err := os.Stat(path); err == nil {
		hasWorktree = true
		if st, _ := a.runner.Run(path, "git", "status", "--porcelain"); strings.TrimSpace(st) != "" {
			return false, fmt.Errorf("the worktree %s has uncommitted changes; commit them on %s (or discard them) before accepting", relPath(repo.MainRoot, path), branch)
		}
		if _, _, err := a.syncStoryBranch(repo, id); err != nil {
			return false, err
		}
	}
	if _, err := a.runner.Run(repo.MainRoot, "git", "merge", "--ff-only", "--quiet", branch); err != nil {
		return false, fmt.Errorf("merge %s into the main branch: %w", branch, err)
	}
	if hasWorktree {
		if _, err := a.runner.Run(repo.MainRoot, "git", "worktree", "remove", path); err != nil {
			return false, err
		}
	}
	if _, err := a.runner.Run(repo.MainRoot, "git", "branch", "-d", branch); err != nil {
		return false, err
	}
	a.logger().Info("story branch merged", "component", "git", "branch", branch)
	return true, nil
}

// dirtyOutsideWip lists uncommitted paths that are not under the wip
// folder, which accumulates transitions between story landings by design.
func (a *app) dirtyOutsideWip(repo *workitem.Repo) []string {
	st, _ := a.runner.Run(repo.MainRoot, "git", "status", "--porcelain")
	wip := strings.TrimSuffix(repo.Manifest.Layout["wip"], "/") + "/"
	var out []string
	for _, line := range strings.Split(st, "\n") {
		if len(line) < 4 {
			continue
		}
		p := strings.TrimSpace(line[3:])
		if i := strings.LastIndex(p, " -> "); i >= 0 {
			p = p[i+4:]
		}
		if !strings.HasPrefix(p, wip) {
			out = append(out, p)
		}
	}
	return out
}
