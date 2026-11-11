// Package storygit reads a project's story branches (ADR-0019) with git:
// whether git is in use, the main branch, whether a story has a branch, and
// what that branch changes. Each story is worked on story/S-nnnn in a
// worktree under .flai-cache/worktrees; wip/ stays in the main checkout.
// The commands and flai serve's reads share it (S-0159).
package storygit

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/bytepunx/system-flow/flai/internal/execx"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
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

// Heads is the commit each local branch of a checkout is at, by branch name,
// and the branch the checkout has checked out, "" when HEAD is at none of
// them, as one git for-each-ref reads them.
type Heads struct {
	Commits map[string]string
	Checked string
}

// ReadHeads reads every local branch's head and the branch checked out at
// root with one git for-each-ref, in place of one rev-parse a branch
// (S-0698). A branch at an annotated tag is at the tag's commit, and one at
// no commit is left out, as rev-parse's ^{commit} finds none for it.
func ReadHeads(r execx.Runner, root string) (Heads, error) {
	read, err := readRefs(r, root)
	if err != nil {
		return Heads{}, err
	}
	return Heads{Commits: read.heads, Checked: read.main}, nil
}

// Main is the branch checked out and the commit it is at, as MainBranch and
// a rev-parse of it read them, or an error saying why there is none.
func (h Heads) Main() (branch, head string, err error) {
	if h.Checked == "" {
		return "", "", fmt.Errorf("the main checkout has a detached HEAD or no branch; check out a branch first")
	}
	head, ok := h.Commits[h.Checked]
	if !ok {
		return "", "", fmt.Errorf("git names no commit for %s; check that git can read the main branch", h.Checked)
	}
	return h.Checked, head, nil
}

// BranchExists reports whether the local branch exists.
func BranchExists(r execx.Runner, root, branch string) bool {
	_, err := r.Run(root, "git", "rev-parse", "--verify", "--quiet", "refs/heads/"+branch)
	return err == nil
}

// brokenStoryRef is the git directory root's worktrees share and, relative
// to it, the loose ref file of story id's branch when that is there but
// resolves to no commit, as a crash while git wrote it leaves it (I-0137).
// ref is "" when the branch resolves or has no loose file, as a packed ref or
// a reftable repository has none, and there is nothing to repair. Such a ref
// holds no commits, so removing it loses no work.
func brokenStoryRef(r execx.Runner, mainRoot, id string) (common, ref string, err error) {
	branch := Branch(id)
	if BranchExists(r, mainRoot, branch) {
		return "", "", nil
	}
	if common, err = commonDir(r, mainRoot); err != nil {
		return "", "", err
	}
	ref = filepath.Join("refs", "heads", filepath.FromSlash(branch))
	if info, statErr := os.Lstat(filepath.Join(common, ref)); statErr == nil && info.Mode().IsRegular() {
		return common, ref, nil
	}
	// no such file; a reftable repository's refs/heads is a file, which
	// fails the stat too
	return common, "", nil
}

// commonDir is the git directory root's worktrees share, absolute.
func commonDir(r execx.Runner, root string) (string, error) {
	out, err := r.Run(root, "git", "rev-parse", "--git-common-dir")
	if err != nil {
		return "", fmt.Errorf("could not find the git directory of %s: %w", root, err)
	}
	if !filepath.IsAbs(out) {
		out = filepath.Join(root, out)
	}
	return filepath.Clean(out), nil
}

// IsWorktree reports whether git knows dir as the top of a work tree. A
// folder with no .git of its own is inside the main checkout's instead, and
// one whose .git file is empty, as a crash leaves it (I-0137), or names a
// registration git no longer has, is not a repository at all.
func IsWorktree(r execx.Runner, dir string) bool {
	top, err := r.Run(dir, "git", "rev-parse", "--show-toplevel")
	if err != nil {
		return false
	}
	want, err := filepath.EvalSymlinks(dir)
	return err == nil && filepath.Clean(filepath.FromSlash(top)) == want
}

// Remote is the remote a story's branch is fetched from: origin when the
// clone has it, else its only remote, else "".
func Remote(r execx.Runner, root string) string {
	out, err := r.Run(root, "git", "remote")
	if err != nil {
		return ""
	}
	remotes := strings.Fields(out)
	switch {
	case slices.Contains(remotes, "origin"):
		return "origin"
	case len(remotes) == 1:
		return remotes[0]
	}
	return ""
}

// FetchBranch fetches branch from remote into this clone's branch of the
// same name, which must not exist. It reports false when the remote has no
// such branch, and an error when the remote could not be asked.
func FetchBranch(r execx.Runner, root, remote, branch string) (bool, error) {
	ref := "refs/heads/" + branch
	out, err := r.Run(root, "git", "ls-remote", "--heads", remote, ref)
	if err != nil {
		return false, fmt.Errorf("could not ask %s whether it has %s: %w", remote, branch, err)
	}
	if strings.TrimSpace(out) == "" {
		return false, nil
	}
	if _, err := r.Run(root, "git", "fetch", "--quiet", remote, ref+":"+ref); err != nil {
		return false, fmt.Errorf("could not fetch %s from %s: %w", branch, remote, err)
	}
	return true, nil
}

// Dirty is the main checkout's uncommitted and untracked files, each file
// of an untracked folder named, split by the shared set (ADR-0248). Shared
// are those in the shared set: under the manifest's wip folder, under its
// design folder, or under or at a path in claims.shared, flai's own work
// between landings, which every acceptance commits. Rest are the others, the
// operator's, which every acceptance leaves as they are.
type Dirty struct {
	Shared []string
	Rest   []string
}

// DirtyMain reads the uncommitted and untracked files of repo's main
// checkout through r, as DirtyFiles reads them, and splits them by the
// shared set of repo's manifest, as InSharedSet reads it. An error says what
// could not be read.
func DirtyMain(r execx.Runner, repo *workitem.Repo) (Dirty, error) {
	files, err := DirtyFiles(r, repo.MainRoot)
	if err != nil {
		return Dirty{}, err
	}
	var d Dirty
	for _, f := range files {
		if InSharedSet(repo.Manifest, f) {
			d.Shared = append(d.Shared, f)
		} else {
			d.Rest = append(d.Rest, f)
		}
	}
	return d, nil
}

// InSharedSet reports whether p, a file's path relative to the repository
// root, is in the shared set of ADR-0248: under the folder m's layout names
// wip or design, or under or at a path a pattern of its claims.shared
// matches, as Claims.Match reads one.
func InSharedSet(m manifest.Manifest, p string) bool {
	for _, key := range []string{"wip", "design"} {
		if dir := strings.TrimSuffix(m.Layout[key], "/"); dir != "" && strings.HasPrefix(p, dir+"/") {
			return true
		}
	}
	_, ok := m.Claims.Match(p)
	return ok
}

// DirtyFiles lists the uncommitted and untracked files of the checkout at
// root, each file of an untracked folder named, read through r.
func DirtyFiles(r execx.Runner, root string) ([]string, error) {
	st, err := r.Run(root, "git", "-c", "core.quotePath=false", "status", "--porcelain", "--no-renames", "--untracked-files=all")
	if err != nil {
		return nil, fmt.Errorf("read the uncommitted files of %s: %w; check that git can read it", root, err)
	}
	return workitem.PorcelainPaths(st), nil
}

// ChangedAmong lists, sorted, the paths of dirty that the change from commit
// from to commit to changes, from the commit they share: of the main
// checkout's uncommitted files, those git refuses to overwrite as it moves
// the checkout's branch from one to the other (ADR-0178). It runs no git when
// dirty is empty.
func ChangedAmong(r execx.Runner, mainRoot, from, to string, dirty []string) ([]string, error) {
	if len(dirty) == 0 {
		return nil, nil
	}
	diff, err := r.Run(mainRoot, "git", "-c", "core.quotePath=false", "diff", "--name-only", "--no-renames", from+"..."+to)
	if err != nil {
		return nil, fmt.Errorf("read the files the change from %s to %s changes: %w; check that both commits are in the repository", from, to, err)
	}
	changed := strings.Split(diff, "\n")
	var out []string
	for _, p := range dirty {
		if slices.Contains(changed, p) && !slices.Contains(out, p) {
			out = append(out, p)
		}
	}
	slices.Sort(out)
	return out, nil
}
