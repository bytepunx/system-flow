package storygit

import (
	"fmt"
	"maps"
	"regexp"
	"slices"
	"strings"

	"github.com/bytepunx/system-flow/flai/internal/execx"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// CommitFiles is, per commit on the main branch, newest first, the files it
// changed: sorted, without repeats, merges and flai's bookkeeping commits
// left out, and the files under the wip folder left out, since flai writes
// them. A commit with no file left has no entry. Story branches are not read,
// so that a commit is not counted once on its branch and again after its
// acceptance. Paths are relative to the project's root. It reads git once,
// in the main checkout, and fails when git or the repository cannot be read.
func CommitFiles(r execx.Runner, repo *workitem.Repo) ([][]string, error) {
	root := repo.MainRoot
	if root == "" {
		root = repo.Root
	}
	if _, err := r.LookPath("git"); err != nil {
		return nil, fmt.Errorf("cannot read the project's commits: git is not on PATH (%w); install git to suggest touches from history", err)
	}
	main, err := MainBranch(r, root)
	if err != nil {
		return nil, fmt.Errorf("cannot read the project's commits in %s: %w; run flai in a git checkout with a branch and a commit on it", root, err)
	}
	// --relative and the pathspec keep the files under the project's root and
	// name them from it.
	out, err := r.Run(root, "git", "-c", "core.quotePath=false", "log", "--no-merges", "--no-renames", "--relative",
		"--name-only", "--format=%x1e%s", main, "--", ".")
	if err != nil {
		return nil, fmt.Errorf("cannot read the commits of %s in %s: %w", main, root, err)
	}
	wip := strings.TrimSuffix(repo.Manifest.Layout["wip"], "/") + "/"
	var commits [][]string
	for _, rec := range strings.Split(out, "\x1e") {
		lines := strings.Split(rec, "\n")
		if bookkeeping(lines[0]) {
			continue
		}
		set := map[string]bool{}
		for _, f := range lines[1:] {
			if f != "" && !strings.HasPrefix(f, wip) {
				set[f] = true
			}
		}
		if len(set) > 0 {
			commits = append(commits, slices.Sorted(maps.Keys(set)))
		}
	}
	return commits, nil
}

// bookkeepingSubject matches the subjects of the commits flai writes itself:
// an acceptance, an item's creation, an edit of an item named by its fields,
// a publish, the project's default agent set or cleared, and an import.
// Releases are among them because every release changes the changelogs and
// template.yaml, which would then seem to change together with everything;
// an older acceptance released in the same commit ("accept and archive;
// release ...") for the same reason.
var bookkeepingSubject = regexp.MustCompile(`^chore: (` +
	`\[[A-Z]-\d+\] (accept and archive|create |edit )` +
	`|publish ` +
	`|(set|clear) the project's default agent$` +
	`|bring .+ under system-flow \(flai import\)$)`)

// bookkeeping reports whether a commit subject is one flai writes itself.
func bookkeeping(subject string) bool {
	return bookkeepingSubject.MatchString(subject)
}
