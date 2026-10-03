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

// storyID is a story's ID as a commit subject names it in brackets, in any
// zero padding.
var storyID = regexp.MustCompile(`^S-\d+$`)

// Committed is, per story by its canonical ID, the files changed by the
// commits on the main branch and the story branches whose subject names it in
// brackets ([S-nnnn]), merges left out, and the files under the wip folder
// left out, since flai writes them (S-0205). Paths are relative to the
// project's root, sorted and without repeats; a story none of whose commits
// changed a file outside wip has no entry. It reads git once, in the main
// checkout, and fails when git or the repository cannot be read.
func Committed(r execx.Runner, repo *workitem.Repo) (map[string][]string, error) {
	root := repo.MainRoot
	if root == "" {
		root = repo.Root
	}
	if _, err := r.LookPath("git"); err != nil {
		return nil, fmt.Errorf("cannot read the stories' commits: git is not on PATH (%w); install git to compare touches with commits", err)
	}
	main, err := MainBranch(r, root)
	if err != nil {
		return nil, fmt.Errorf("cannot read the stories' commits in %s: %w; run flai in a git checkout with a branch and a commit on it", root, err)
	}
	// --relative and the pathspec keep the files under the project's root and
	// name them from it; --full-history keeps every commit that touched them,
	// whichever side of a merge it is on.
	out, err := r.Run(root, "git", "-c", "core.quotePath=false", "log", "--no-merges", "--no-renames", "--full-history", "--relative",
		"--name-only", "--format=%x1e%s", main, "--branches="+Prefix+"*", "--", ".")
	if err != nil {
		return nil, fmt.Errorf("cannot read the commits of %s and the story branches in %s: %w", main, root, err)
	}
	wip := strings.TrimSuffix(repo.Manifest.Layout["wip"], "/") + "/"
	files := map[string]map[string]bool{}
	for _, rec := range strings.Split(out, "\x1e") {
		lines := strings.Split(rec, "\n")
		ids := subjectStories(lines[0])
		if len(ids) == 0 {
			continue
		}
		for _, f := range lines[1:] {
			if f == "" || strings.HasPrefix(f, wip) {
				continue
			}
			for _, id := range ids {
				if files[id] == nil {
					files[id] = map[string]bool{}
				}
				files[id][f] = true
			}
		}
	}
	committed := make(map[string][]string, len(files))
	for id, set := range files {
		committed[id] = slices.Sorted(maps.Keys(set))
	}
	return committed, nil
}

// subjectStories is every story a commit subject names in brackets, by its
// canonical ID, read as release reads the IDs of a message.
func subjectStories(subject string) []string {
	var ids []string
	for i := 0; i < len(subject); i++ {
		if subject[i] != '[' {
			continue
		}
		end := strings.IndexByte(subject[i+1:], ']')
		if end < 0 {
			break
		}
		if tok := subject[i+1 : i+1+end]; storyID.MatchString(tok) {
			if id := workitem.CanonicalID(tok); !slices.Contains(ids, id) {
				ids = append(ids, id)
			}
		}
	}
	return ids
}
