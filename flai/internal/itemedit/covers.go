package itemedit

import (
	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// Covered is an open story whose claim covers paths another story changed,
// and those paths.
type Covered struct {
	ID    string   `json:"id" jsonschema:"the open story"`
	Title string   `json:"title" jsonschema:"its title"`
	Paths []string `json:"paths" jsonschema:"the changed paths its claim covers"`
}

// Covering is who a story's change reaches (S-0132, S-0333): each story in
// progress or in review other than changer whose claim covers one of the
// changed paths, with the paths it covers, in the items' order. The shared
// paths are no exception: they keep two claims from holding each other, not
// a change from being told. A story whose claim is empty may change
// anything, so every changed path is its business. Stories whose claim
// covers none of the paths are left out.
func Covering(items []*workitem.Item, projects []manifest.Project, changer string, changed []string) []Covered {
	if len(changed) == 0 {
		return nil
	}
	changer = workitem.CanonicalID(changer)
	holds := workitem.NewHolds(items, projects)
	var out []Covered
	for _, it := range items {
		if it.Archived || it.Type != workitem.Story || it.ID == changer || (it.Status != workitem.InProgress && it.Status != workitem.Review) {
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
