package check

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/bytepunx/system-flow/flai/internal/workflow"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// The codes flai check reports a workflow under besides the parser's,
// workflow.file, workflow.rule, and workflow.timer (ADR-0301 §2, §3, §6):
// the workflow package's, which judges the workflow; workflow.item, a
// finding on an item; and workflow.suite and workflow.suite-unnamed, which
// hold the suites the file names against the manifest's test tiers
// (ADR-0312 §9). The last two are flai check's alone: no workflow.Finding
// carries them, so neither refuses the file (ADR-0307 §2).
const (
	codeState        = workflow.CodeState
	codeUnreachable  = workflow.CodeUnreachable
	codeDeadEnd      = workflow.CodeDeadEnd
	codeNoPull       = workflow.CodeNoPull
	codeNoAccept     = workflow.CodeNoAccept
	codeItem         = "workflow.item"
	codeSuite        = "workflow.suite"
	codeSuiteUnnamed = "workflow.suite-unnamed"
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
	refused := c.flowRefused()
	for _, f := range w.Findings() {
		c.add(flowLevel(f), f.Code, abs, f.LineIn(string(data)), "%s", f.Message)
		refused = refused || f.Level == workflow.LevelError
	}
	if !refused {
		c.suites(w, abs)
	}
}

// suites holds the suites the project's workflow file w, at abs, names
// against the suites of the manifest's test tiers, each tier's
// SuiteOrDefault (ADR-0312 §9). A suite the file names that no tier has is
// workflow.suite, an error on the line it is named on, once for each place
// that names it; all is no suite's name and never one. A suite of a tier
// that is not manual, which no transition, state, or point names, is
// workflow.suite-unnamed, a warning once for the suite, on verify_at's line
// or else the transitions', naming its tiers that are not manual: no step
// runs it. A file that names all anywhere runs every suite and gets no
// warning. A project without the file runs the baseline, which names only
// all, and workflow does not call suites for it; nor while the file is
// refused, as the items are not judged against it, nor while the tests key
// is wrong, which manifest.tests reports.
func (c *checker) suites(w workflow.Workflow, abs string) {
	tiers, err := c.repo.Manifest.TestTiers(os.DirFS(c.repo.Root))
	if err != nil {
		return
	}
	var have []string
	for _, t := range tiers {
		if s := t.SuiteOrDefault(); !slices.Contains(have, s) {
			have = append(have, s)
		}
	}
	named, all := w.NamedSuites()
	for _, n := range named {
		if slices.Contains(have, n.Suite) {
			continue
		}
		tiersHave := "system-flow.yaml declares no test tiers"
		if len(have) > 0 {
			tiersHave = "the suites of system-flow.yaml's test tiers are " + strings.Join(have, ", ")
		}
		c.add(Error, codeSuite, abs, n.Line, "%s names suite %s, which no test tier has; %s. Name a suite a tier in system-flow.yaml's tests gives, add suite: %s to the tier it belongs to, or give all for every suite (ADR-0312 §9)", namedAt(n), n.Suite, tiersHave, n.Suite)
	}
	if all {
		return
	}
	line := keyLine(abs, "verify_at")
	if line == 1 {
		line = keyLine(abs, "transitions")
	}
	for _, suite := range have {
		if slices.ContainsFunc(named, func(n workflow.Named) bool { return n.Suite == suite }) {
			continue
		}
		var run []string
		for _, t := range tiers {
			if !t.Manual && t.SuiteOrDefault() == suite {
				run = append(run, strings.TrimSpace(t.Name))
			}
		}
		if len(run) == 0 {
			continue
		}
		of := "tier " + run[0]
		if len(run) > 1 {
			of = "tiers " + strings.Join(run, ", ")
		}
		c.add(Warning, codeSuiteUnnamed, abs, line, "suite %s, of %s, is named by no transition's verify, state's runs, or point of verify_at, so no step runs it; name it in one of them, give all for every suite, or mark %s manual: true in system-flow.yaml's tests (ADR-0312 §9)", suite, of, of)
	}
}

// namedAt is where the workflow file names n's suite, in words: transition
// in-progress → review's verify, state review's runs, or verify_at.batch.
func namedAt(n workflow.Named) string {
	switch n.Where {
	case workflow.WhereTransition:
		return "transition " + n.Name + "'s verify"
	case workflow.WhereState:
		return "state " + n.Name + "'s runs"
	}
	return "verify_at." + n.Name
}

// flowRefused reports whether the project's workflow file is one flai
// cannot run, as c.repo.FlowErr() keeps it (ADR-0307 §2): flai check opens
// the repository without refusing it, and the repository runs the baseline
// in its place. While it is, the judgments of an item, and of the board's
// order and limits, against the workflow's states, kinds, and transitions
// are left out, workflow.item and item.sequence among them, since they would
// judge the items against a board the project does not run; workflow reports
// the file's own findings, which say what to fix. A project without the file
// runs the baseline and is judged against it, with no advisory (ADR-0307 §5).
func (c *checker) flowRefused() bool {
	return c.repo.FlowErr() != nil
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
// status line when the status names it, else on the transitions line. It
// reports nothing while the project's workflow file is refused (flowRefused).
func (c *checker) itemStates(it *workitem.Item) {
	if c.flowRefused() {
		return
	}
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
