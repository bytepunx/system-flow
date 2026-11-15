package itemedit

import (
	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/workflow"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// Covered is an open story whose claim covers paths another story changed,
// and those paths.
type Covered struct {
	ID    string   `json:"id" jsonschema:"the open story"`
	Title string   `json:"title" jsonschema:"its title"`
	Paths []string `json:"paths" jsonschema:"the changed paths its claim covers"`
}

// Covering is CoveringIn under the baseline workflow, for a caller with no
// repository's workflow at hand.
func Covering(items []*workitem.Item, projects []manifest.Project, changer string, changed []string) []Covered {
	return CoveringIn(workflow.Baseline(), items, projects, changer, changed)
}

// CoveringIn is who a story's change reaches (S-0132, S-0333): each story in
// a state of flow of kind active or waiting, in progress or in review on the
// baseline, other than changer whose claim covers one of the changed paths,
// with the paths it covers, in the items' order. The shared paths are no
// exception: they keep two claims from holding each other, not a change from
// being told. A story whose claim is empty may change anything, so every
// changed path is its business. Stories whose claim covers none of the paths
// are left out.
func CoveringIn(flow workflow.Workflow, items []*workitem.Item, projects []manifest.Project, changer string, changed []string) []Covered {
	if len(changed) == 0 {
		return nil
	}
	changer = workitem.CanonicalID(changer)
	holds := workitem.NewHolds(items, projects)
	var out []Covered
	for _, it := range items {
		if it.Archived || it.Type != workitem.Story || it.ID == changer || !underWay(flow, it.Status) {
			continue
		}
		claim := holds.Claim(it)
		var paths []string
		for _, p := range changed {
			if len(claim) == 0 || claimCovers(claim, p) {
				paths = append(paths, p)
			}
		}
		if len(paths) > 0 {
			out = append(out, Covered{ID: it.ID, Title: it.Title, Paths: paths})
		}
	}
	return out
}

// underWay says whether status is a state of flow of kind active or waiting:
// a story worked on or awaiting acceptance, whatever the board names the
// state (ADR-0301).
func underWay(flow workflow.Workflow, status string) bool {
	switch flow.Kind(status) {
	case workflow.KindActive, workflow.KindWaiting:
		return true
	}
	return false
}

// claimCovers reports whether a path of the claim overlaps path, shared or
// not.
func claimCovers(claim []string, path string) bool {
	for _, c := range claim {
		if workitem.PathsOverlap(path, c) {
			return true
		}
	}
	return false
}
