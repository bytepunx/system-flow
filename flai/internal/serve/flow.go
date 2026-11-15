package serve

import (
	"github.com/bytepunx/system-flow/flai/internal/workflow"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// flow is the workflow the served projects run: Workflow, else the baseline
// flai carries (ADR-0301).
func (o Options) flow() workflow.Workflow {
	if o.Workflow != nil {
		return *o.Workflow
	}
	return baselineFlow
}

// baselineFlow is the baseline workflow, read once; no pass changes it.
var baselineFlow = workflow.Baseline()

// openRepo opens the project at root as workitem.Open does, running the
// workflow o gives, so that its moves and the passes that read it follow
// one board.
func openRepo(o Options, root string) (*workitem.Repo, error) {
	repo, err := workitem.Open(root)
	if err != nil {
		return nil, err
	}
	repo.Workflow = o.Workflow
	return repo, nil
}
