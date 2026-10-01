package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/bytepunx/system-flow/flai/internal/gitver"
	"github.com/bytepunx/system-flow/flai/internal/storygit"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// Story branches (ADR-0019): each story is worked on story/S-nnnn in a git
// worktree under .flai-cache/worktrees; wip/ stays in the main checkout so
// the board is live. stream open creates the branch, stream sync rebases it
// onto the main branch, accept merges it and removes the worktree. What is
// only read of them is in storygit, which flai serve's reads share (S-0159);
// these run it with the command's runner.

func storyBranch(id string) string { return storygit.Branch(id) }

func (a *app) inGitWorkTree(root string) bool { return storygit.InWorkTree(a.runner, root) }

func (a *app) mainBranch(mainRoot string) (string, error) {
	return storygit.MainBranch(a.runner, mainRoot)
}

func (a *app) branchExists(root, branch string) bool {
	return storygit.BranchExists(a.runner, root, branch)
}

// Where openStoryBranch found a story's branch, besides a remote's name.
const (
	branchLocal = "local"
	branchMain  = "main"
)

// openStoryBranch checks story/<id> out in a worktree: this clone's branch
// when it has one, else, when fromRemote, the remote's, fetched (ADR-0064),
// else a new one from the main branch. It returns the branch, the path, and
// where the branch came from (branchLocal, a remote's name, or branchMain;
// "" when the worktree was there already), or "" for all when git is not in
// use.
func (a *app) openStoryBranch(repo *workitem.Repo, id string, fromRemote bool) (branch, path, from string, err error) {
	if !a.inGitWorkTree(repo.MainRoot) {
		return "", "", "", nil
	}
	branch, path = storyBranch(id), repo.WorktreePath(id)
	if _, err := os.Stat(path); err == nil {
		return branch, path, "", nil
	}
	base, err := a.mainBranch(repo.MainRoot)
	if err != nil {
		return "", "", "", err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return "", "", "", err
	}
	from = branchLocal
	if !a.branchExists(repo.MainRoot, branch) && fromRemote {
		if remote := storygit.Remote(a.runner, repo.MainRoot); remote != "" {
			fetched, err := storygit.FetchBranch(a.runner, repo.MainRoot, remote, branch)
			if err != nil {
				return "", "", "", fmt.Errorf("%w; run flai stream open %s again when %s answers, or create the branch from %s with git branch %s %s first", err, id, remote, base, branch, base)
			}
			if fetched {
				from = remote
			}
		}
	}
	args := []string{"worktree", "add", "--quiet"}
	if a.relativeWorktrees() {
		args = append(args, "--relative-paths")
	}
	if a.branchExists(repo.MainRoot, branch) {
		args = append(args, path, branch)
	} else {
		from = branchMain
		args = append(args, "-b", branch, path, base)
	}
	_, err = a.runner.Run(repo.MainRoot, "git", args...)
	if err != nil {
		return "", "", "", err
	}
	a.logger().Info("story branch opened", "component", "git", "branch", branch, "worktree", relPath(repo.MainRoot, path), "base", base, "from", from)
	return branch, path, from, nil
}

// relativeWorktrees reports whether to link a new worktree with relative
// paths: only when the operator set worktrees.relative_paths, and only with
// a git that understands them. The git version never turns it on by itself
// (ADR-0022); with the key on and an older git an ordinary worktree is made,
// which the dashboard's host-path mount covers.
func (a *app) relativeWorktrees() bool {
	cfg, _, err := a.loadConfig()
	if err != nil || !cfg.Worktrees.RelativePaths {
		return false
	}
	v, err := gitver.Installed(a.runner)
	if err != nil {
		a.logger().Warn("git version unreadable, creating an ordinary worktree", "component", "git", "setting", "worktrees.relative_paths", "err", err)
		return false
	}
	if !v.AtLeast(gitver.RelativeWorktrees) {
		a.logger().Warn("git is too old for relative worktree paths, creating an ordinary worktree", "component", "git", "setting", "worktrees.relative_paths", "git", v.String(), "needs", gitver.RelativeWorktrees.String())
		return false
	}
	return true
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
