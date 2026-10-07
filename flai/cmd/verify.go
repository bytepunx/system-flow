package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/spf13/cobra"

	"github.com/bytepunx/system-flow/flai/internal/check"
	"github.com/bytepunx/system-flow/flai/internal/verify"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// exitVerifyUnusable is flai verify's exit status when it could not answer:
// a usage error, a story with no worktree, tiers the manifest does not
// declare validly, or a run stopped before it finished, kept apart from a
// failing step's 1, as flai test keeps them (S-0270).
const exitVerifyUnusable = 2

// newVerifyCmd runs a story's close-out steps and the tiers its branch
// selects in its worktree and answers one report (S-0270).
func newVerifyCmd(a *app) *cobra.Command {
	var record, last bool
	var maxFindings int
	c := &cobra.Command{
		Use:   "verify <story>",
		Short: "Run a story's close-out steps and the tiers its branch selects, in its worktree, and answer one result",
		Long: `Verify a story before it goes to review, in its worktree under
.flai-cache/worktrees, from any checkout of the project. The steps run in
order, cheapest first, and the run stops at the first that fails; the steps
after it are not reached:

  rebase     no rebase is left unfinished in the worktree
  sync       the story's branch contains the main branch
  narrative  the narrative's Current state and Next steps are written
  check      flai check --strict, scoped to the story, passes
  <tier>     each test and lint tier of the worktree's system-flow.yaml that
             what the branch changed against the main branch selects, run
             with CLOSE_OUT_STORY set to the story and FLAI_ROLE=verify,
             whatever role ran flai verify, so a test that runs flai's
             writes is not refused as the orchestrator (S-0311)

flai verify commits nothing. The answer is a line for each step, its state
and duration, the failing step's first findings under it, at most --max
across the run, then the check's findings outside the story, which are notes
that do not fail it, and last a line as the close-out ends:
"verify: S-nnnn passed every step", or "verify: S-nnnn stopped at <step>
(<exit status or state>)". With --json the answer is the report, with the
story, the commit and base it was verified at, when it ran, how long it took,
passed, stopped_at, the paths the branch changed, and each step's state,
duration, and findings, and the notes.

Each run's report is stored in the project's .flai-cache/verify, whether it
passes or not; --last prints the story's stored report and runs nothing, or
says there is none (null with --json), and exits 0 whatever it holds.

--record-issues records the notes in design/issues of the story's worktree,
as flai check --story --record-issues does: each rule's notes in the open
issue whose title names the rule, once per story and notes. The close-out
commits them with the story.

The exit status is 0 when every step passed, 1 when a step failed, and 2 when
flai verify could not answer: the story has no worktree, the manifest's tiers
are not valid, the notes could not be recorded, or the run was stopped before
it finished.`,
		Example: `  flai verify S-0270
  flai verify S-0270 --json
  flai verify S-0270 --record-issues
  flai verify S-0270 --last --json`,
		Args: func(_ *cobra.Command, args []string) error {
			if len(args) != 1 {
				return &exitError{code: exitVerifyUnusable, msg: "name the story to verify: flai verify S-nnnn"}
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if maxFindings < 1 {
				return &exitError{code: exitVerifyUnusable, msg: fmt.Sprintf("--max %d is not a count of findings; give 1 or more", maxFindings)}
			}
			if last && record {
				return &exitError{code: exitVerifyUnusable, msg: "--last prints the stored result and runs nothing, so it records nothing; drop --record-issues, or --last to run the story again"}
			}
			repo, err := a.project()
			if err != nil {
				return &exitError{code: exitVerifyUnusable, msg: err.Error()}
			}
			if last {
				return a.printLastVerify(repo, args[0], maxFindings)
			}
			return a.verifyStory(cmd, repo, args[0], verifyOptions{record: record, max: maxFindings})
		},
	}
	c.SetFlagErrorFunc(func(_ *cobra.Command, err error) error {
		return &exitError{code: exitVerifyUnusable, msg: err.Error()}
	})
	c.Flags().BoolVar(&record, "record-issues", false, "record the check's findings outside the story in design/issues of its worktree, once per story and findings")
	c.Flags().BoolVar(&last, "last", false, "print the story's stored last result and run nothing")
	c.Flags().IntVar(&maxFindings, "max", verify.DefaultMax, "the most findings to report across the run")
	return c
}

// verifyOptions are flai verify's flags for a run.
type verifyOptions struct {
	record bool
	max    int
}

// verifyRecorded is flai verify --json with --record-issues: the report,
// and the issues its notes were recorded in.
type verifyRecorded struct {
	verify.Report
	Recorded []recordedIssue `json:"recorded"`
}

// verifyStory runs verify.Verify for the story in its worktree, records its
// notes when asked, prints the report, and answers the exit status.
func (a *app) verifyStory(cmd *cobra.Command, repo *workitem.Repo, story string, o verifyOptions) error {
	id, worktree, err := storyWorktree(repo, story)
	if err != nil {
		return &exitError{code: exitVerifyUnusable, msg: err.Error()}
	}
	tiers, err := verify.CheckoutTiers(worktree)
	if err != nil {
		return &exitError{code: exitVerifyUnusable, msg: fmt.Sprintf("verify %s: %v", id, err)}
	}
	ctx, stop := signal.NotifyContext(cmd.Context(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	rep, err := verify.Verify(ctx, verify.StoryOptions{
		Story: id, Project: repo, Worktree: worktree, Tiers: tiers, Git: a.runner,
		RunOptions: verify.RunOptions{Max: o.max, Now: a.now},
	})
	if err != nil && rep.Story == "" {
		return &exitError{code: exitVerifyUnusable, msg: err.Error()}
	}
	stopped := err != nil && ctx.Err() != nil
	a.logger().Debug("story verified", "component", "cmd", "story", rep.Story, "steps", len(rep.Steps), "notes", len(rep.Notes), "passed", rep.Passed)
	var recorded []recordedIssue
	if o.record {
		if recorded, err = a.recordNotes(worktree, rep); err != nil {
			return &exitError{code: exitVerifyUnusable, msg: err.Error()}
		}
	}
	if a.jsonOut {
		var v any = rep
		if o.record {
			v = verifyRecorded{Report: rep, Recorded: recorded}
		}
		if err := a.printJSON(v); err != nil {
			return err
		}
	} else {
		fmt.Fprint(a.out, verifyText(rep, recorded, stopped, o.max))
	}
	switch {
	case stopped:
		return &exitError{code: exitVerifyUnusable, msg: "the run was stopped before it finished: " + err.Error()}
	case err != nil:
		return &exitError{code: exitVerifyUnusable, msg: err.Error()}
	case !rep.Passed:
		return &exitError{code: 1}
	}
	return nil
}

// storyWorktree is the story's canonical ID and its worktree under
// .flai-cache/worktrees, refusing an item that is not a story and a story
// that has none.
func storyWorktree(repo *workitem.Repo, story string) (id, worktree string, err error) {
	it, err := repo.Get(story)
	if err != nil {
		return "", "", fmt.Errorf("verify %s: %w; name a story of this project", story, err)
	}
	if it.Type != workitem.Story {
		return "", "", fmt.Errorf("verify %s: it is a %s; name a story (S-nnnn)", it.ID, it.Type)
	}
	worktree = repo.WorktreePath(it.ID)
	if st, err := os.Stat(worktree); err != nil || !st.IsDir() {
		return "", "", fmt.Errorf("verify %s: it has no worktree at %s; open it with flai stream open %s, then verify it", it.ID, worktree, it.ID)
	}
	return it.ID, worktree, nil
}

// recordNotes records the report's notes in design/issues of the story's
// worktree through flai check's recorder, as the check's findings outside
// the story they are.
func (a *app) recordNotes(worktree string, rep verify.Report) ([]recordedIssue, error) {
	wt, err := workitem.Open(worktree)
	if err != nil {
		return nil, fmt.Errorf("record the notes outside %s: open its worktree %s: %w", rep.Story, worktree, err)
	}
	res := &check.Result{Findings: make([]check.Finding, 0, len(rep.Notes))}
	for _, n := range rep.Notes {
		res.Findings = append(res.Findings, check.Finding{Level: n.Level, Rule: n.Rule, Path: n.Path, Line: n.Line, Message: n.Message, Outside: true})
	}
	return a.recordOutside(wt, res, rep.Story)
}

// printLastVerify prints the story's stored last report, or that it has
// none, and answers success whatever the report holds.
func (a *app) printLastVerify(repo *workitem.Repo, story string, maxNotes int) error {
	rep, ok, err := verify.LastReport(repo, story)
	if err != nil {
		return &exitError{code: exitVerifyUnusable, msg: err.Error()}
	}
	if a.jsonOut {
		if !ok {
			return a.printJSON(nil)
		}
		return a.printJSON(rep)
	}
	if !ok {
		id := story
		if it, err := repo.Get(story); err == nil {
			id = it.ID
		}
		fmt.Fprintf(a.out, "verify: %s has no stored result; run flai verify %s\n", id, id)
		return nil
	}
	fmt.Fprintf(a.out, "last verified %s at %s against %s\n", rep.RanAt.Format(time.RFC3339), shortCommit(rep.Commit), rep.Base)
	fmt.Fprint(a.out, verifyText(rep, nil, false, maxNotes))
	return nil
}

// verifyText is the report as flai verify prints it: a line for each step,
// its state, name, and duration, a tier's name after "tier", the failing
// step's findings indented under it, then at most maxNotes of the notes
// outside the story, the issues they were recorded in, and last the
// outcome as the close-out's last line has it.
func verifyText(rep verify.Report, recorded []recordedIssue, stopped bool, maxNotes int) string {
	var b strings.Builder
	for _, s := range rep.Steps {
		name := s.Name
		if s.Tier {
			name = "tier " + s.Name
		}
		if s.State == verify.NotReached {
			fmt.Fprintf(&b, "%s %s\n", s.State, name)
		} else {
			fmt.Fprintf(&b, "%s %s (%s)\n", s.State, name, s.Duration)
		}
		b.WriteString(stepFindings(s))
	}
	if len(rep.Notes) > 0 {
		fmt.Fprintf(&b, "outside %s, notes that do not fail it:\n", rep.Story)
		for _, n := range rep.Notes[:min(maxNotes, len(rep.Notes))] {
			b.WriteString("    " + noteText(n) + "\n")
		}
		if left := len(rep.Notes) - maxNotes; left > 0 {
			fmt.Fprintf(&b, "    … %d more notes left out\n", left)
		}
	}
	for _, r := range recorded {
		fmt.Fprintf(&b, "recorded %s outside %s in %s (%s)\n", r.Rule, rep.Story, r.Issue, r.describe())
	}
	b.WriteString(verifyOutcome(rep, stopped) + "\n")
	return b.String()
}

// stepFindings are the step's findings and how many were left out, as flai
// test prints a tier's, each line indented under the step.
func stepFindings(s verify.Step) string {
	var b strings.Builder
	for _, f := range s.Findings {
		for _, line := range strings.SplitAfter(f.Text(), "\n") {
			if line != "" {
				b.WriteString("    " + line)
			}
		}
	}
	if s.Omitted > 0 {
		fmt.Fprintf(&b, "    … %d more findings left out\n", s.Omitted)
	}
	return b.String()
}

// noteText is a note as flai check prints a finding: path:line: level:
// rule: message, the message on one line.
func noteText(n verify.Note) string {
	loc := ""
	if n.Path != "" {
		loc = n.Path + ":" + strconv.Itoa(n.Line) + ": "
	}
	return loc + n.Level + ": " + n.Rule + ": " + strings.Join(strings.Fields(n.Message), " ")
}

// verifyOutcome is the report's last line, in the close-out's form: the
// story passed every step, or the step it stopped at and why.
func verifyOutcome(rep verify.Report, stopped bool) string {
	if rep.Passed && !stopped {
		return fmt.Sprintf("verify: %s passed every step", rep.Story)
	}
	step, reason := rep.StoppedAt, "failed"
	for _, s := range rep.Steps {
		if s.Name != rep.StoppedAt {
			continue
		}
		if s.ExitCode != nil && *s.ExitCode != 0 {
			reason = "exit " + strconv.Itoa(*s.ExitCode)
		}
		break
	}
	if stopped {
		reason = "interrupted"
	}
	if step == "" {
		step = "an unfinished step"
	}
	return fmt.Sprintf("verify: %s stopped at %s (%s)", rep.Story, step, reason)
}

// shortCommit is a commit's first seven characters, as git abbreviates it.
func shortCommit(commit string) string {
	if len(commit) > 7 {
		return commit[:7]
	}
	return commit
}
