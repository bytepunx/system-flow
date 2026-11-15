package hostapi

import (
	"github.com/bytepunx/system-flow/flai/internal/workflow"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// testFlow, when set, is the workflow every repository openRepo opens runs:
// a test sets it, with withWorkflow, to try the methods on a board with
// renamed and added states (ADR-0301). Nil leaves each repository the
// workflow it runs.
var testFlow *workflow.Workflow

// openRepo opens the project above root for work item operations, running
// the workflow the project runs (ADR-0301).
func openRepo(root string) (*workitem.Repo, error) {
	repo, err := workitem.Open(root)
	if err != nil {
		return nil, err
	}
	if testFlow != nil {
		repo.Workflow = testFlow
	}
	return repo, nil
}
