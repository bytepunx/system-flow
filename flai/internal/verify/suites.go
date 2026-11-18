package verify

import (
	"fmt"
	"slices"

	"github.com/bytepunx/system-flow/flai/internal/workflow"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// NeededSuites are the suites a story in state needs its close-out to run:
// those the verify of each transition w lets a story take out of state
// lists, each once, in file order (ADR-0312 §2). Nil is every suite: when
// one of those transitions gives verify: all, and when none gives verify,
// so that a workflow that names no suite for them runs what it runs today.
// A workflow file that gives none of ADR-0312's keys reads as the
// baseline's (Workflow.UsesSuiteKeys, D-0669), whose move to review gives
// verify: all.
func NeededSuites(w workflow.Workflow, state string) []string {
	if !w.UsesSuiteKeys() {
		w = workflow.Baseline()
	}
	var names []string
	for _, t := range w.From(workitem.Story, state) {
		if t.Verify == nil {
			continue
		}
		if t.Verify.Suites.All {
			return nil
		}
		for _, name := range t.Verify.Suites.Names {
			if !slices.Contains(names, name) {
				names = append(names, name)
			}
		}
	}
	return names
}

// StorySuites are the suites the close-out of story runs: NeededSuites of
// the project's workflow at the state the story is in; nil is every suite.
// It refuses an ID that names no story of the project.
func StorySuites(project *workitem.Repo, story string) ([]string, error) {
	it, err := project.Get(story)
	if err != nil {
		return nil, fmt.Errorf("find the suites %s's close-out runs: %w; name a story of this project", story, err)
	}
	if it.Type != workitem.Story {
		return nil, fmt.Errorf("find the suites %s's close-out runs: it is a %s; name a story (S-nnnn)", it.ID, it.Type)
	}
	return NeededSuites(project.Flow(), it.Status), nil
}
