package check

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"

	"github.com/bytepunx/system-flow/flai/internal/gitver"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// GitCompat adds git.relative-worktrees when the clone carries
// extensions.relativeWorktrees and the installed git is too old to open it
// (ADR-0022). The extension is read from .git/config directly, because that
// git refuses the repository. The caller supplies the installed version so
// this package runs no commands.
func GitCompat(res *Result, repo *workitem.Repo, installed gitver.Version) {
	root := repo.MainRoot
	if root == "" {
		root = repo.Root
	}
	path := filepath.Join(root, ".git", "config")
	line := relativeWorktreesLine(path)
	if line == 0 || installed.AtLeast(gitver.RelativeWorktrees) {
		return
	}
	rel := path
	if r, err := filepath.Rel(repo.Root, path); err == nil {
		rel = r
	}
	res.Findings = append(res.Findings, Finding{Level: Warning, Rule: "git.relative-worktrees", Path: rel, Line: line,
		Message: "this clone has extensions.relativeWorktrees set (worktrees.relative_paths was on) and git " + installed.String() +
			" on PATH refuses such a repository; git " + gitver.RelativeWorktrees.String() + " or newer is needed. " +
			"Upgrade git, or go back: with a newer git run `git worktree repair --no-relative-paths <worktree>` for each worktree and `git config --unset extensions.relativeWorktrees`; " +
			"with only this git, delete the relativeWorktrees line from .git/config and run `git worktree repair <worktree>` for each worktree before anything prunes them. " +
			"Then `flai config set worktrees.relative_paths false`."})
	res.Warnings++
}

// relativeWorktreesLine returns the 1-based line in a git config file where
// extensions.relativeWorktrees is set to a true value, or 0.
func relativeWorktreesLine(path string) int {
	f, err := os.Open(path)
	if err != nil {
		return 0
	}
	defer func() { _ = f.Close() }()
	sc := bufio.NewScanner(f)
	section, n := "", 0
	for sc.Scan() {
		n++
		l := strings.TrimSpace(sc.Text())
		if strings.HasPrefix(l, "[") {
			section = strings.ToLower(strings.Trim(l, "[] \t"))
			continue
		}
		key, value, _ := strings.Cut(l, "=")
		if section != "extensions" || !strings.EqualFold(strings.TrimSpace(key), "relativeworktrees") {
			continue
		}
		switch strings.ToLower(strings.TrimSpace(value)) {
		case "", "true", "yes", "on", "1": // a bare key is true in git config
			return n
		}
	}
	return 0
}
