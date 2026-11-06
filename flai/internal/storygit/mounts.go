package storygit

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/bytepunx/system-flow/flai/internal/manifest"
)

// Mount is a host folder, or a host file, a story's container mounts at the
// same path.
type Mount struct {
	Path     string // absolute and cleaned
	ReadOnly bool
	// Create is set on a folder the host makes before the container starts
	// when it is absent, since Docker would make a missing source as root and
	// the container's user could not write to it (ADR-0253).
	Create bool
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
	reg, common, err := linked(worktree)
	if err != nil || reg == "" {
		return nil, err
	}
	return []Mount{{Path: common, ReadOnly: true}, {Path: reg}}, nil
}

// MainMounts is what a container running in worktree needs mounted of the
// main checkout beside GitMounts' for flai verify to run there against the
// main checkout's state (ADR-0253): its wip folder, as its manifest's
// layout names it, read-only, its
// system-flow.yaml read-only, and its .flai-cache/verify folder read-write,
// in that order, none of them over another or over the .git. The main
// checkout is the folder that holds the common dir the worktree's .git file
// and its registration's commondir name, as GitMounts resolves them. A
// worktree whose .git is a folder carries its repository and has no main
// checkout to mount. A common dir that is not a main checkout's .git, and a
// main checkout with no system-flow.yaml or no wip folder, are errors naming
// the path. The verify folder may be absent, and carries Create, so the
// caller makes it on the host. It reads files only, and runs no git.
func MainMounts(worktree string) ([]Mount, error) {
	reg, common, err := linked(worktree)
	if err != nil || reg == "" {
		return nil, err
	}
	if filepath.Base(common) != ".git" {
		return nil, fmt.Errorf("the repository %s of %s is not a main checkout's .git, so the main checkout whose wip and system-flow.yaml flai verify reads cannot be found: open the story's worktree with flai stream open from the project's main checkout", common, worktree)
	}
	main := filepath.Dir(common)
	manifestPath := filepath.Join(main, "system-flow.yaml")
	if info, err := os.Stat(manifestPath); err != nil || info.IsDir() {
		return nil, fmt.Errorf("%s is not there, so %s, the main checkout of %s, is not a system-flow project's: open the story's worktree from the project's main checkout", manifestPath, main, worktree)
	}
	m, err := manifest.Load(manifestPath)
	if err != nil {
		return nil, fmt.Errorf("reading the main checkout's manifest for its wip folder: %w", err)
	}
	wip := m.Dir(main, "wip")
	if info, err := os.Stat(wip); err != nil || !info.IsDir() {
		return nil, fmt.Errorf("%s is not there, so the main checkout %s has no work items for flai verify to read: restore it with git checkout in %s", wip, main, main)
	}
	return []Mount{
		{Path: wip, ReadOnly: true},
		{Path: manifestPath, ReadOnly: true},
		{Path: filepath.Join(main, ".flai-cache", "verify"), Create: true},
	}, nil
}

// linked is worktree's registration and the common dir it names, both empty
// when worktree's .git is a folder, which carries its repository.
func linked(worktree string) (reg, common string, err error) {
	worktree, err = filepath.Abs(worktree)
	if err != nil {
		return "", "", err
	}
	dotGit := filepath.Join(worktree, ".git")
	info, err := os.Stat(dotGit)
	if errors.Is(err, os.ErrNotExist) {
		return "", "", fmt.Errorf("%s is not there, so %s is not a git worktree: flai stream open makes a story's worktree", dotGit, worktree)
	}
	if err != nil {
		return "", "", err
	}
	if info.IsDir() {
		return "", "", nil
	}
	if reg, err = registration(dotGit); err != nil {
		return "", "", err
	}
	if common, err = linkedCommonDir(worktree, reg); err != nil {
		return "", "", err
	}
	return reg, common, nil
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
