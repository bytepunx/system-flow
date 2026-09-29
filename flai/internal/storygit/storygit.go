// Package storygit reads a project's story branches (ADR-0019) with git:
// whether git is in use, the main branch, whether a story has a branch, and
// what that branch changes. Each story is worked on story/S-nnnn in a
// worktree under .flai-cache/worktrees; wip/ stays in the main checkout.
// The commands and flai serve's reads share it (S-0159).
package storygit

import (
	"fmt"
	"strings"

	"github.com/bytepunx/system-flow/flai/internal/execx"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// Prefix begins every story branch's name.
const Prefix = "story/"

// Branch is the name of a story's branch.
func Branch(id string) string { return Prefix + id }

// InWorkTree reports whether root is inside a git work tree.
func InWorkTree(r execx.Runner, root string) bool {
	if _, err := r.LookPath("git"); err != nil {
		return false
	}
	_, err := r.Run(root, "git", "rev-parse", "--is-inside-work-tree")
	return err == nil
}

// MainBranch is the branch checked out in the main checkout.
func MainBranch(r execx.Runner, mainRoot string) (string, error) {
	out, err := r.Run(mainRoot, "git", "rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		return "", err
	}
	if out == "HEAD" {
		return "", fmt.Errorf("the main checkout has a detached HEAD; check out a branch first")
	}
	return out, nil
}

// BranchExists reports whether the local branch exists.
func BranchExists(r execx.Runner, root, branch string) bool {
	_, err := r.Run(root, "git", "rev-parse", "--verify", "--quiet", "refs/heads/"+branch)
	return err == nil
}

// DirtyOutsideWip lists uncommitted paths that are not under the wip
// folder, which accumulates transitions between story landings by design.
func DirtyOutsideWip(r execx.Runner, repo *workitem.Repo) []string {
	st, err := r.Run(repo.MainRoot, "git", "status", "--porcelain")
	if err != nil {
		return nil // not a repository: nothing to be dirty
	}
	wip := strings.TrimSuffix(repo.Manifest.Layout["wip"], "/") + "/"
	var out []string
	for _, p := range workitem.PorcelainPaths(st) {
		if !strings.HasPrefix(p, wip) {
			out = append(out, p)
		}
	}
	return out
}
