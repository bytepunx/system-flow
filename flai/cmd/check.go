package cmd

import (
	"fmt"
	"maps"
	"path/filepath"
	"slices"
	"strings"

	"github.com/spf13/cobra"

	"github.com/bytepunx/system-flow/flai/internal/check"
	"github.com/bytepunx/system-flow/flai/internal/issues"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

func newCheckCmd(a *app) *cobra.Command {
	var strict, record bool
	var story string
	c := &cobra.Command{
		Use:   "check [dir]",
		Short: "Validate the repository against the system-flow standard",
		Long: `Check the manifest, layout, every work item in kanban and archive, the
narratives and their index, the board, and documentation front matter.

Findings print as path:line: level: rule: message. Errors exit 1; with
--strict warnings do too, except two that only the operator clears: the
review column over its limit, which acceptance clears, and an epic behind
its stories (epic.lags-stories), which moving or accepting the epic clears.

With --story S-nnnn, a finding is inside the story when it is on the
story's item file or one of its tasks', its narrative, a thread anchored on
the story or one of its tasks, or a path its branch changes against the main
branch (as flai stream diff reads it), or that is uncommitted in its
worktree. Every other finding is outside it. A wip.overlap that names the
story is outside it too, and one between two other stories is left out, as
is an item.archive outside the story, which only flai archive in the main
checkout clears, a board.wip-limit outside the story, which counts the main
checkout's columns and only the operator's acceptances, cancellations, or a
raised limit clear, and a markdown or narrative.state finding on the narrative
of another story that is neither done nor cancelled, which that story's own
close-out finds. A finding outside keeps its level and is printed with
"(outside S-nnnn)"; it is a note that neither an error nor --strict fails
on, and the summary counts it, as does outside in --json. Without --story
every finding counts, as the main branch's check needs.

With --record-issues as well, each rule with findings outside the story is
recorded in design/issues, in the checkout the run reads, which is the
story's worktree at close-out. A wip.overlap is not recorded: the pull hold
and the other story's agent clear it. Each rule is recorded in the open
issue whose title names it, as "flai check finds ` + "`threads.archived`" + ` outside
the story at close-out" does, opened with class efficiency when none is,
with an instance naming the story and each finding's path and message. An
instance for the same story and the same findings is written once, so
running the check again does not count it again; another story, or other
findings, bump the count.
summary.md is regenerated, each issue is printed as opened, bumped, or
already recorded, and --json lists them in recorded. The run records
whether or not it passes; a failure to record is an error.`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if record && story == "" {
				return fmt.Errorf("--record-issues records the findings outside a story: give --story S-nnnn")
			}
			var repo *workitem.Repo
			var err error
			if len(args) == 1 {
				repo, err = workitem.Open(args[0])
			} else {
				repo, err = a.project()
			}
			if err != nil {
				return err
			}
			if story != "" {
				it, err := repo.Get(story)
				if err != nil {
					return err
				}
				story = it.ID
			}
			res, err := check.RunScoped(repo, a.runner, a.now(), story)
			if err != nil {
				return err
			}
			var recorded []recordedIssue
			if record {
				if recorded, err = a.recordOutside(repo, res, story); err != nil {
					return err
				}
			}
			if a.jsonOut {
				var v any = res
				if record {
					v = checkRecorded{Result: res, Recorded: recorded}
				}
				if err := a.printJSON(v); err != nil {
					return err
				}
			} else {
				for _, f := range res.Findings {
					fmt.Fprintf(a.out, "%s:%d: %s: %s: %s", f.Path, f.Line, f.Level, f.Rule, f.Message)
					if f.Outside {
						fmt.Fprintf(a.out, " (outside %s)", story)
					}
					fmt.Fprintln(a.out)
				}
				fmt.Fprintf(a.out, "%d items checked, %d errors, %d warnings", res.Items, res.Errors, res.Warnings)
				if res.Advisory > 0 {
					// S-0243: say which warnings --strict does not fail on.
					fmt.Fprintf(a.out, " (%d that --strict passes over: review over its limit and an epic behind its stories, which only the operator clears, or a narrative not yet written, which flai stream state writes)", res.Advisory)
				}
				if res.Outside > 0 {
					// S-0249: say how many findings were notes outside the story.
					fmt.Fprintf(a.out, "; %d outside %s, notes the run passes over", res.Outside, story)
				}
				fmt.Fprintln(a.out)
				for _, r := range recorded {
					fmt.Fprintf(a.out, "recorded %s outside %s in %s (%s)\n", r.Rule, story, r.Issue, r.describe())
				}
			}
			if !res.OK(strict) {
				return &exitError{code: 1}
			}
			return nil
		},
	}
	c.Flags().BoolVar(&strict, "strict", false, "treat warnings as failures, except the review column over its limit and an epic behind its stories")
	c.Flags().StringVar(&story, "story", "", "scope the run to a story (S-nnnn): findings outside it are notes that do not fail it")
	c.Flags().BoolVar(&record, "record-issues", false, "with --story, record each rule's findings outside the story in an issue, opening or bumping it once per story and findings")
	return c
}

// checkRecorded is flai check --json with --record-issues: the result, and
// the issues the findings outside the story were recorded in.
type checkRecorded struct {
	*check.Result
	Recorded []recordedIssue `json:"recorded"`
}

// recordedIssue is the issue one rule's findings outside a story were
// recorded in, and what recording them did.
type recordedIssue struct {
	Rule    string         `json:"rule"`
	Issue   string         `json:"issue"`
	Path    string         `json:"path"`
	Outcome issues.Outcome `json:"outcome"`
	Count   int            `json:"count"`
}

func (r recordedIssue) describe() string {
	if r.Outcome == issues.Opened {
		return string(r.Outcome)
	}
	return fmt.Sprintf("%s, count %d", r.Outcome, r.Count)
}

// recordOutside records each rule's findings outside story in its open issue
// (issues.RecordOnce), in rule order, and regenerates summary.md when one was
// opened or bumped. A wip.overlap is never recorded (ADR-0115 §3): the pull
// hold, the overlapped change, the report on a growing claim, and the trial
// merge report it where it can be acted on, and an issue would only count
// how often stories ran side by side (I-0076). An item.archive outside the
// story never reaches it: the scoped check leaves it out (ADR-0122), as it
// does a board.wip-limit outside the story (ADR-0133) and a markdown or
// narrative.state finding on another open story's narrative (ADR-0123,
// ADR-0125).
func (a *app) recordOutside(repo *workitem.Repo, res *check.Result, story string) ([]recordedIssue, error) {
	byRule := map[string][]string{}
	for _, f := range res.Findings {
		if !f.Outside || f.Rule == "wip.overlap" {
			continue
		}
		line := findingText(f.Message)
		if f.Path != "" {
			line = "`" + projectPath(repo, f.Path) + "`: " + line
		}
		if !slices.Contains(byRule[f.Rule], line) {
			byRule[f.Rule] = append(byRule[f.Rule], line)
		}
	}
	rules := slices.Sorted(maps.Keys(byRule))
	out := make([]recordedIssue, 0, len(rules))
	changed := false
	for _, rule := range rules {
		lines := byRule[rule]
		slices.Sort(lines)
		is, outcome, err := issues.RecordOnce(repo, issues.NewOptions{
			Title: fmt.Sprintf("flai check finds `%s` outside the story at close-out", rule),
			Class: "efficiency", Story: story, Now: a.now(), Runner: a.runner,
			Note: "flai check found outside the story:\n" + strings.Join(lines, "\n"),
		})
		if err != nil {
			return nil, fmt.Errorf("record the %s findings outside %s as an issue: %w; fix the issue file or run without --record-issues", rule, story, err)
		}
		changed = changed || outcome != issues.Already
		out = append(out, recordedIssue{Rule: rule, Issue: is.ID, Path: relPath(repo.Root, is.Path), Outcome: outcome, Count: is.Count})
	}
	if changed {
		if _, err := issues.WriteSummary(repo, a.now()); err != nil {
			return nil, fmt.Errorf("regenerate %s after recording the findings outside %s: %w", issues.SummaryFile, story, err)
		}
	}
	return out, nil
}

// findingText is a finding's message as one line of an issue's markdown:
// its whitespace collapsed, and its backticks made quotes, since a message
// that quotes a code span, such as a markdown finding's context, would
// otherwise open one the markdown lint refuses.
func findingText(message string) string {
	return strings.ReplaceAll(strings.Join(strings.Fields(message), " "), "`", "'")
}

// projectPath is p as the project names it, from its main checkout where
// it has one: run in a story's worktree, the findings on wip/, which lives in
// the main checkout, are named from the worktree with ../ (S-0249).
func projectPath(repo *workitem.Repo, p string) string {
	abs := p
	if !filepath.IsAbs(abs) {
		abs = filepath.Join(repo.Root, p)
	}
	for _, root := range []string{repo.Root, repo.MainRoot} {
		if root == "" {
			continue
		}
		if rel, err := filepath.Rel(root, abs); err == nil && rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
			return filepath.ToSlash(rel)
		}
	}
	return filepath.ToSlash(p)
}

// exitError signals a non-zero exit, with an optional message the boundary
// logs as the fatal event.
type exitError struct {
	code int
	msg  string
}

func (e *exitError) Error() string { return e.msg }
