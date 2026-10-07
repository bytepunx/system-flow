package itemedit

import (
	"reflect"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// Who a story's change reaches (S-0333): each other open story whose claim
// covers a changed path, by the prefix rule, a component by its path, a
// task's touches as its story's, a shared path like any other, and an empty
// claim as everything.
func TestCovering(t *testing.T) {
	story := func(id, status string, touches ...string) *workitem.Item {
		return &workitem.Item{ID: id, Type: workitem.Story, Status: status, Title: "title " + id, Touches: touches}
	}
	archived := story("S-0009", workitem.InProgress, "flai/cmd")
	archived.Archived = true
	items := []*workitem.Item{
		story("S-0001", workitem.InProgress, "flai/cmd"),           // the changer
		story("S-0002", workitem.InProgress, "flai/cmd/accept.go"), // a file it changed
		story("S-0003", workitem.Review, "cli"),                    // the component, by tag
		story("S-0004", workitem.InProgress, "flaiover"),           // elsewhere
		story("S-0005", workitem.InProgress, "flaiover"),           // its task claims design/system
		story("S-0006", workitem.Review),                           // no touches: may change anything
		story("S-0007", workitem.Ready, "flai/cmd"),                // not open
		story("S-0008", workitem.InProgress, "CHANGELOG.md"),       // told whatever the shared paths hold
		archived,
		{ID: "T-0001", Type: workitem.Task, Status: workitem.InProgress, Parent: "S-0005", Touches: []string{"design/system"}},
	}
	projects := []manifest.Project{{Name: "flai", Path: "flai", Tags: []string{"cli"}}, {Name: "flaiover", Path: "flaiover"}}
	changed := []string{"flai/cmd/accept.go", "flai/cmd/move.go", "design/system/workflow.md", "CHANGELOG.md"}
	got := Covering(items, projects, "S-1", changed)
	want := []Covered{
		{ID: "S-0002", Title: "title S-0002", Paths: []string{"flai/cmd/accept.go"}},
		{ID: "S-0003", Title: "title S-0003", Paths: []string{"flai/cmd/accept.go", "flai/cmd/move.go"}},
		{ID: "S-0005", Title: "title S-0005", Paths: []string{"design/system/workflow.md"}},
		{ID: "S-0006", Title: "title S-0006", Paths: changed},
		{ID: "S-0008", Title: "title S-0008", Paths: []string{"CHANGELOG.md"}},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("covering:\n got %+v\nwant %+v", got, want)
	}
	if got := Covering(items, projects, "S-0001", nil); got != nil {
		t.Errorf("a change of nothing reaches nobody, got %+v", got)
	}
	if got := Covering(items, projects, "S-0001", []string{"README.md"}); got[0].ID != "S-0006" || len(got) != 1 {
		t.Errorf("a path no claim covers reaches only the empty claim, got %+v", got)
	}
}
