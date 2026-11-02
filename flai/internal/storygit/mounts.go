package storygit

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Mount is a host folder a story's container mounts at the same path.
type Mount struct {
	Path     string // absolute and cleaned
	ReadOnly bool
}

// GitMounts is what a container running in worktree needs mounted beside it
// for git to work there (ADR-0246): the repository's common dir read-only,
// then the worktree's own registration read-write over it, in that order.
// The worktree's .git file names the registration, and the registration's
// commondir the common dir, each absolute or relative to the folder that
// holds it, as worktrees.relative_paths links them (ADR-0022). A worktree
// whose .git is a folder carries its repository and needs no mount. A
// missing .git, registration, or commondir is an error naming the path, so a
// start fails with the reason rather than a container without git (I-0215).
// It reads files only, and runs no git.
func GitMounts(worktree string) ([]Mount, error) {
	worktree, err := filepath.Abs(worktree)
	if err != nil {
		return nil, err
	}
	dotGit := filepath.Join(worktree, ".git")
	info, err := os.Stat(dotGit)
	if errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("%s is not there, so %s is not a git worktree: flai stream open makes a story's worktree", dotGit, worktree)
	}
	if err != nil {
		return nil, err
	}
	if info.IsDir() {
		return nil, nil
	}
	reg, err := registration(dotGit)
	if err != nil {
		return nil, err
	}
	common, err := linkedCommonDir(worktree, reg)
	if err != nil {
		return nil, err
	}
	return []Mount{{Path: common, ReadOnly: true}, {Path: reg}}, nil
}

// registration is the folder the worktree's .git file names after gitdir:,
// resolved against the worktree, which must be there.
func registration(dotGit string) (string, error) {
	b, err := os.ReadFile(dotGit)
	if err != nil {
		return "", err
	}
	line, _, _ := strings.Cut(string(b), "\n")
	gitdir, ok := strings.CutPrefix(strings.TrimSpace(line), "gitdir:")
	gitdir = strings.TrimSpace(gitdir)
	if !ok || gitdir == "" {
		return "", fmt.Errorf("%s names no gitdir, so its worktree is not linked to a repository: run git worktree repair %s from the main checkout", dotGit, filepath.Dir(dotGit))
	}
	reg := resolve(filepath.Dir(dotGit), gitdir)
	if info, err := os.Stat(reg); err != nil || !info.IsDir() {
		return "", fmt.Errorf("%s names %s as its registration, which is not there: run git worktree repair %s from the main checkout if the repository moved, or git worktree prune and flai stream open to make the worktree again", dotGit, reg, filepath.Dir(dotGit))
	}
	return reg, nil
}

// linkedCommonDir is the folder the registration's commondir file names,
// resolved against the registration, which must be there.
func linkedCommonDir(worktree, reg string) (string, error) {
	file := filepath.Join(reg, "commondir")
	b, err := os.ReadFile(file)
	if errors.Is(err, os.ErrNotExist) {
		return "", fmt.Errorf("%s is not there, so %s is not a linked worktree's registration: run git worktree repair %s from the main checkout", file, reg, worktree)
	}
	if err != nil {
		return "", err
	}
	named := strings.TrimSpace(string(b))
	if named == "" {
		return "", fmt.Errorf("%s is empty, so %s names no repository: run git worktree repair %s from the main checkout", file, reg, worktree)
	}
	common := resolve(reg, named)
	if info, err := os.Stat(common); err != nil || !info.IsDir() {
		return "", fmt.Errorf("%s names %s as the repository, which is not there: run git worktree repair %s from the main checkout", file, common, worktree)
	}
	return common, nil
}

// resolve is p cleaned, joined to dir when relative.
func resolve(dir, p string) string {
	if filepath.IsAbs(p) {
		return filepath.Clean(p)
	}
	return filepath.Join(dir, p)
}
