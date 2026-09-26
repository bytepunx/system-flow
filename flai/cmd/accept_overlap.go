package cmd

import (
	"fmt"
	"strings"

	"github.com/bytepunx/system-flow/flai/internal/itemedit"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// overlapsOf is what accepting a story tells the stories still open (S-0132,
// ADR-0046): for each story in progress or in review other than accepted,
// the changed paths its claim covers. A story whose claim is empty may change
// anything, so every changed path is its business. Stories whose claim covers
// none of the paths are left out.
func overlapsOf(items []*workitem.Item, projects []manifest.Project, accepted string, changed []string) []itemedit.Overlap {
	if len(changed) == 0 {
		return nil
	}
	holds := workitem.NewHolds(items, projects)
	var out []itemedit.Overlap
	for _, it := range items {
		if it.Archived || it.Type != workitem.Story || it.ID == accepted || (it.Status != workitem.InProgress && it.Status != workitem.Review) {
			continue
		}
		claim := holds.Claim(it)
		var paths []string
		for _, p := range changed {
			if len(claim) == 0 || covered(p, claim) {
				paths = append(paths, p)
			}
		}
		if len(paths) > 0 {
			out = append(out, itemedit.Overlap{ID: it.ID, Title: it.Title, Accepted: accepted, Paths: paths})
		}
	}
	return out
}

func covered(path string, claim []string) bool {
	for _, c := range claim {
		if workitem.PathsOverlap(path, c) {
			return true
		}
	}
	return false
}

// headOf is the commit the main checkout is at, or "" when git cannot say.
func (a *app) headOf(root string) string {
	out, err := a.runner.Run(root, "git", "rev-parse", "HEAD")
	if err != nil {
		return ""
	}
	return strings.TrimSpace(out)
}

// changedSince lists the paths that differ between from and the head of the
// main checkout: what a story's merge brought into the main branch. A rename
// lists both of its paths, since a story may claim either.
func (a *app) changedSince(root, from string) ([]string, error) {
	out, err := a.runner.Run(root, "git", "diff", "--name-only", "--no-renames", from, "HEAD")
	if err != nil {
		return nil, fmt.Errorf("list the paths the merge changed: %w", err)
	}
	var paths []string
	for _, p := range strings.Split(out, "\n") {
		if p = strings.TrimSpace(p); p != "" {
			paths = append(paths, p)
		}
	}
	return paths, nil
}

// tellOverlaps records a notice for each open story whose claim covers what
// the accepted story changed, for the MCP inbox and wait_for_events to report
// to its agent. It is best effort after the acceptance has been committed: a
// notice that cannot be worked out is logged, and the acceptance stands.
func (a *app) tellOverlaps(repo *workitem.Repo, it *workitem.Item, by string, changed []string) []itemedit.Overlap {
	if len(changed) == 0 {
		return nil
	}
	items, err := repo.List(false)
	if err != nil {
		a.logger().Warn("overlap notices not sent", "component", "accept", "item", it.ID, "err", err)
		return nil
	}
	told := overlapsOf(items, repo.Manifest.Projects, it.ID, changed)
	at := a.now().UTC().Format(workitem.TimeFormat)
	for i := range told {
		told[i].At, told[i].By = at, by
		itemedit.RecordOverlap(repo, told[i])
	}
	return told
}
