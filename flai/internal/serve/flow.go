package serve

import (
	"errors"
	"log/slog"
	"sync"

	"github.com/bytepunx/system-flow/flai/internal/workflow"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// flow is the workflow the options give the served projects, Workflow, else
// the baseline flai carries (ADR-0301). A pass with a repository at hand
// reads the repository's own, Repo.Flow, which is Workflow when it is set and
// the project's workflow file otherwise (ADR-0307); flow is for a pass that
// could not open one.
func (o Options) flow() workflow.Workflow {
	if o.Workflow != nil {
		return *o.Workflow
	}
	return baselineFlow
}

// baselineFlow is the baseline workflow, read once; no pass changes it.
var baselineFlow = workflow.Baseline()

// openRepo opens the project at root as workitem.Open does, running the
// workflow o gives when it gives one, and the project's workflow file
// otherwise, so that its moves and the passes that read it follow one board.
// A workflow file flai cannot run refuses the open, and the refusal is logged
// once for each distinct finding at the root (ADR-0307 §3).
func openRepo(o Options, root string) (*workitem.Repo, error) {
	repo, err := workitem.Open(root)
	if err != nil {
		refusedFlows.note(o.Logger, root, err)
		return nil, err
	}
	refusedFlows.note(o.Logger, root, nil)
	repo.Workflow = o.Workflow
	return repo, nil
}

// refusedFlows is what openRepo last logged of each root's workflow file
// refusal: every pass opens its project afresh at each look, and a file that
// stays broken is told once, not at every look of every pass.
var refusedFlows = &flowRefusals{said: map[string]string{}}

// flowRefusals is the workflow file refusal last logged, by project root.
type flowRefusals struct {
	mu   sync.Mutex
	said map[string]string
}

// note logs err at warn level when it is a workflow file refusal whose
// finding differs from the one last logged for root, and forgets root's when
// err is nil, an open that succeeded, so that the file breaking again is told
// again. Other errors are the passes' own to tell.
func (f *flowRefusals) note(logger *slog.Logger, root string, err error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if err == nil {
		delete(f.said, root)
		return
	}
	var le *workflow.LoadError
	if !errors.As(err, &le) {
		return
	}
	finding := le.Error()
	if f.said[root] == finding {
		return
	}
	f.said[root] = finding
	if logger != nil {
		logger.Warn("project's workflow file cannot be run; flai serve's passes leave the project alone until it is fixed", "component", "serve", "root", root, "err", err.Error())
	}
}
