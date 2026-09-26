package workitem

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// Uncommitted lists the paths with uncommitted changes, untracked files
// among them, in a story's worktree (ADR-0019). It lists none when the story
// has no worktree or the repo has no Git to ask (S-0140).
func (r *Repo) Uncommitted(storyID string) ([]string, error) {
	if r.Git == nil {
		return nil, nil
	}
	path := r.WorktreePath(storyID)
	if _, err := os.Stat(path); err != nil {
		return nil, nil
	}
	out, err := r.Git.Run(path, "git", "status", "--porcelain")
	if err != nil {
		return nil, fmt.Errorf("git status in %s: %w", r.RelToMain(path), err)
	}
	return PorcelainPaths(out), nil
}

// RelToMain is path relative to the main checkout, for messages; path itself
// when it is not below it.
func (r *Repo) RelToMain(path string) string {
	if rel, err := filepath.Rel(r.mainRoot(), path); err == nil && !strings.HasPrefix(rel, "..") {
		return filepath.ToSlash(rel)
	}
	return path
}

// UncommittedRule is the refusal for a story whose worktree holds
// uncommitted paths, saying what to do and when: before is the step it
// stops, such as "it goes to review".
func (r *Repo) UncommittedRule(storyID string, paths []string, before string) string {
	return fmt.Sprintf("the worktree %s has uncommitted changes (%s); commit them on story/%s, or discard them, before %s",
		r.RelToMain(r.WorktreePath(storyID)), Shorten(paths, 10), storyID, before)
}

// Shorten joins paths with commas, naming at most max and counting the rest.
func Shorten(paths []string, max int) string {
	if len(paths) <= max {
		return strings.Join(paths, ", ")
	}
	return fmt.Sprintf("%s and %d more", strings.Join(paths[:max], ", "), len(paths)-max)
}

var porcelainLine = regexp.MustCompile(`^\s*\S{1,2}\s+(.+)$`)

// PorcelainPaths extracts the paths from `git status --porcelain` output.
// The runner trims the whole output, so the first line may have lost the
// leading space of an unstaged status (" M path"); parse by shape, not by
// column. Renames report the new path.
func PorcelainPaths(out string) []string {
	var paths []string
	for _, line := range strings.Split(out, "\n") {
		m := porcelainLine.FindStringSubmatch(line)
		if m == nil {
			continue
		}
		p := strings.Trim(m[1], `"`)
		if i := strings.LastIndex(p, " -> "); i >= 0 {
			p = p[i+4:]
		}
		paths = append(paths, p)
	}
	return paths
}
