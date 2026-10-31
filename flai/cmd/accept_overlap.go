package cmd

import (
	"fmt"
	"strings"

	"github.com/bytepunx/system-flow/flai/internal/itemedit"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// overlapsOf is what accepting a story tells the stories still open (S-0132,
// ADR-0046): each open story itemedit.Covering finds the accepted story's
// change reaches, as an overlap notice. A changed path inside a pattern of
// shared is left out first (S-0550): many stories change those paths in
// separate sections or new files, so their change is nobody's notice, as it
// holds nobody (ADR-0096), and a change wholly inside them tells nobody.
func overlapsOf(items []*workitem.Item, projects []manifest.Project, shared manifest.Claims, accepted string, changed []string) []itemedit.Overlap {
	var told []string
	for _, p := range changed {
		if _, ok := shared.Match(p); !ok {
			told = append(told, p)
		}
	}
	var out []itemedit.Overlap
	for _, c := range itemedit.Covering(items, projects, accepted, told) {
		out = append(out, itemedit.Overlap{ID: c.ID, Title: c.Title, Accepted: accepted, Paths: c.Paths})
	}
	return out
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
// the accepted story changed outside the shared paths, as system-flow.yaml
// lists them now, for the MCP inbox and wait_for_events to report to its
// agent. It is best effort after the acceptance has been committed: a
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
	told := overlapsOf(items, repo.Manifest.Projects, repo.SharedClaims(), it.ID, changed)
	at := a.now().UTC().Format(workitem.TimeFormat)
	for i := range told {
		told[i].At, told[i].By = at, by
		itemedit.RecordOverlap(repo, told[i])
	}
	return told
}
