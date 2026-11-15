package check

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/bytepunx/system-flow/flai/internal/workflow"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// The codes flai check reports a workflow under besides the parser's,
// workflow.file, workflow.rule, and workflow.timer (ADR-0301 §2, §3, §6):
// the workflow package's, which judges the workflow, and workflow.item, a
// finding on an item.
const (
	codeState       = workflow.CodeState
	codeUnreachable = workflow.CodeUnreachable
	codeDeadEnd     = workflow.CodeDeadEnd
	codeNoPull      = workflow.CodeNoPull
	codeNoAccept    = workflow.CodeNoAccept
	codeItem        = "workflow.item"
)

// workflowRel is the path, relative to the root, of the project's workflow
// file under the design folder, whether or not the project has it.
func (c *checker) workflowRel() string {
	return workflow.Rel(c.repo.Manifest.Layout["design"])
}

// workflow validates the workflow the project runs (ADR-0301): the design
// folder's workflow file when the project has one, else c.repo.Flow(), the
// baseline flai carries or the one a test gives the repository, reported on
// the path the project's file would have, on line 1. A file that does not
// parse is one error under its problem's code, workflow.file, workflow.rule,
// or workflow.timer, and nothing else is checked on it; a workflow that
// parses is held to what flai needs to run it, as workflow.Findings judges
// it, each finding on its entry's or its key's line in the file.
func (c *checker) workflow() {
	rel := c.workflowRel()
	abs := filepath.Join(c.repo.Root, filepath.FromSlash(rel))
	data, err := os.ReadFile(abs)
	switch {
	case errors.Is(err, os.ErrNotExist):
		for _, f := range c.repo.Flow().Findings() {
			c.add(flowLevel(f), f.Code, abs, 1, "%s", f.Message)
		}
		return
	case err != nil:
		c.add(Error, workflow.CodeFile, abs, 1, "%s cannot be read: %v; make it readable, or remove it to run the baseline workflow flai carries", rel, err)
		return
	}
	w, prob := workflow.Parse(rel, string(data))
	if prob != nil {
		msg := prob.Message
		if prob.Field != "" {
			msg = prob.Field + ": " + msg
		}
		c.add(Error, prob.Code, filepath.Join(c.repo.Root, filepath.FromSlash(prob.Path)), prob.Line, "%s", msg)
		return
	}
	for _, f := range w.Findings() {
		c.add(flowLevel(f), f.Code, abs, f.LineIn(string(data)), "%s", f.Message)
	}
}

// flowLevel is the check level of a workflow finding.
func flowLevel(f workflow.Finding) string {
	if f.Level == workflow.LevelWarning {
		return Warning
	}
	return Error
}

// itemStates reports workflow.item, an error, once for each name item's
// status or a transition of its history gives that the workflow the
// repository runs neither has as a state nor lists as an alias (ADR-0301
// §6): flai cannot tell which state the item is in, or was. It is on the
// status line when the status names it, else on the transitions line.
func (c *checker) itemStates(it *workitem.Item) {
	flow := c.repo.Flow()
	seen := map[string]bool{}
	report := func(name, key, where string) {
		if _, ok := flow.State(name); ok || name == "" || seen[name] {
			return
		}
		seen[name] = true
		c.add(Error, codeItem, it.Path, keyLine(it.Path, key), "%s's %s names state %s, which the workflow neither has nor lists as an alias; add the state to %s, or add %s to the aliases of the state its items go to (ADR-0301 §6)", it.ID, where, name, c.workflowRel(), name)
	}
	report(it.Status, "status", "status")
	for i, tr := range it.Transitions {
		report(tr.To, "transitions", fmt.Sprintf("transitions[%d]", i))
	}
}
