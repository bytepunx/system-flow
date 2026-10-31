package check

import (
	"fmt"
	"maps"
	"slices"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/issues"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// Fixer rewrites the generated file a rule's findings are on from the files
// it is generated from, and returns that file's path from the project root.
type Fixer func(repo *workitem.Repo, now time.Time) (string, error)

// fixers maps a rule to the Fixer of its generated file. Only a rule on a
// generated file has one: a hand-written file is never rewritten.
var fixers = map[string]Fixer{
	"issues.summary": fixIssuesSummary,
}

// Fixed is one file flai check --fix rewrote, and the rule that reported it.
type Fixed struct {
	Rule string `json:"rule"`
	Path string `json:"path"`
}

// Fix runs, once each and in rule order, the Fixer of every rule res reports
// inside the run, passing over the findings outside the story it is scoped
// to, and returns the files they rewrote.
func Fix(repo *workitem.Repo, res *Result, now time.Time) ([]Fixed, error) {
	reported := map[string]bool{}
	for _, f := range res.Findings {
		if _, ok := fixers[f.Rule]; ok && !f.Outside {
			reported[f.Rule] = true
		}
	}
	out := []Fixed{}
	for _, rule := range slices.Sorted(maps.Keys(reported)) {
		path, err := fixers[rule](repo, now)
		if err != nil {
			return out, fmt.Errorf("fix %s: %w", rule, err)
		}
		out = append(out, Fixed{Rule: rule, Path: path})
	}
	return out, nil
}

// fixIssuesSummary rewrites design/issues/summary.md as flai issue summary
// does.
func fixIssuesSummary(repo *workitem.Repo, now time.Time) (string, error) {
	path := issues.SummaryPath(repo)
	if _, err := issues.WriteSummary(repo, now); err != nil {
		return path, fmt.Errorf("rewrite %s from the issue files: %w; fix the issue file it names, then run flai issue summary", path, err)
	}
	return path, nil
}
