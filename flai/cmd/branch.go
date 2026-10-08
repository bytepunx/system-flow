package cmd

import (
	"errors"
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/bytepunx/system-flow/flai/internal/gitver"
	"github.com/bytepunx/system-flow/flai/internal/issues"
	"github.com/bytepunx/system-flow/flai/internal/preview"
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

// openStoryBranch checks story/<id> out in a worktree with storygit.Open,
// with the command's runner, logger, and worktrees.relative_paths.
func (a *app) openStoryBranch(repo *workitem.Repo, id string, fromRemote bool) (storygit.Opened, error) {
	return storygit.Open(storygit.OpenOptions{Runner: a.runner, Repo: repo, Story: id, FromRemote: fromRemote, RelativePaths: a.relativeWorktrees, Log: a.logger()})
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

// syncStopped is a sync that did not rebase the story branch (ADR-0069):
// refused over uncommitted changes or a rebase already in progress, or
// stopped on conflicts. It names each path and how to continue or abort, so
// an agent never has to work out a rebase by hand and no work is lost.
type syncStopped struct {
	Story    string
	Base     string
	Worktree string // relative to the main checkout
	// Uncommitted is set when the sync was refused over uncommitted changes;
	// nothing was touched.
	Uncommitted []string
	// RebaseInProgress is set when a rebase waits in the worktree, whether
	// this sync stopped it or found it.
	RebaseInProgress bool
	// Found is set when the rebase was in progress before this sync, which
	// refused to touch it.
	Found     bool
	Conflicts []string
	Continue  string
	Abort     string // "" when there is nothing to abort
}

// headline says why the branch was not synced, with what names the paths:
// the paths themselves for the one-line error, a count for the report.
func (s *syncStopped) headline(paths string) string {
	b := storyBranch(s.Story)
	switch {
	case s.Found && paths != "":
		return fmt.Sprintf("%s was not synced: a rebase is already in progress in %s, with conflicts in %s", b, s.Worktree, paths)
	case s.Found:
		return fmt.Sprintf("%s was not synced: a rebase is already in progress in %s", b, s.Worktree)
	case s.RebaseInProgress && paths != "":
		return fmt.Sprintf("%s was not synced: the rebase onto %s stopped on conflicts in %s", b, s.Base, paths)
	case s.RebaseInProgress:
		return fmt.Sprintf("%s was not synced: the rebase onto %s stopped", b, s.Base)
	}
	return fmt.Sprintf("%s was not synced: its worktree %s has uncommitted changes in %s", b, s.Worktree, paths)
}

// Error is one short line, for the error a command exits with: why, the
// paths (at most ten), and the next step.
func (s *syncStopped) Error() string {
	msg := s.headline(workitem.Shorten(s.paths(), 10))
	if s.RebaseInProgress {
		return msg + fmt.Sprintf("; resolve them and git rebase --continue, or git rebase --abort, in %s", s.Worktree)
	}
	return msg + fmt.Sprintf("; commit them on %s, then sync again", storyBranch(s.Story))
}

// Report is why, each path on a line of its own, and the steps to continue
// and abort, for flai stream sync to print.
func (s *syncStopped) Report() string {
	var b strings.Builder
	paths := s.paths()
	if len(paths) > 0 {
		b.WriteString(s.headline(plural(len(paths), "path")) + ":\n")
	} else {
		b.WriteString(s.headline("") + "\n")
	}
	for _, p := range paths {
		b.WriteString("  " + p + "\n")
	}
	b.WriteString("To continue: " + s.Continue + "\n")
	if s.Abort != "" {
		b.WriteString("To abort: " + s.Abort + "\n")
	}
	return b.String()
}

// paths is every path the stop names: uncommitted or conflicting.
func (s *syncStopped) paths() []string { return slices.Concat(s.Uncommitted, s.Conflicts) }

// nonNil is l, or an empty list for JSON's [] in place of null.
func nonNil(l []string) []string {
	if l == nil {
		return []string{}
	}
	return l
}

// syncStoryBranch rebases the story branch onto the main branch inside its
// worktree with storygit.Rebase, the rebase of flai stream sync without its
// checks (ADR-0069, ADR-0098). A refusal or a rebase left for the agent is
// a *syncStopped, with the conflicting paths when there are any.
func (a *app) syncStoryBranch(repo *workitem.Repo, id string) (base string, conflicts []string, err error) {
	// the rebase reads only the story's ID
	res, err := storygit.Rebase(storygit.SyncOptions{Runner: a.runner, Repo: repo, Story: &workitem.Item{ID: id}, Now: a.now(), Files: issues.SyncFiles(repo, id, a.runner, a.now), Log: a.logger()})
	if err != nil {
		return res.Base, nil, err
	}
	if !res.Synced {
		return res.Base, res.Conflicts, syncStoppedFrom(res)
	}
	return res.Base, nil, nil
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
	// a conflict git-added unresolved must not reach main (S-0253, I-0066)
	if err := a.refuseConflictMarkers(repo, id); err != nil {
		return false, err
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

// refuseConflictMarkers refuses a story branch any of whose files added or
// changed against the main branch carries a merge conflict marker, naming
// each path and line, by the check the acceptance preview reports by
// (S-0253, S-0276, I-0066). It runs after the sync, so it reads the branch
// as it will be merged.
func (a *app) refuseConflictMarkers(repo *workitem.Repo, id string) error {
	msg, err := preview.ConflictMarkers(a.runner, repo, id)
	if err != nil || msg == "" {
		return err
	}
	return errors.New(msg)
}
