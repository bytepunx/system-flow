package cmd

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/spf13/cobra"

	"github.com/bytepunx/system-flow/flai/internal/buildinfo"
	"github.com/bytepunx/system-flow/flai/internal/check"
	"github.com/bytepunx/system-flow/flai/internal/closeout"
	"github.com/bytepunx/system-flow/flai/internal/config"
	"github.com/bytepunx/system-flow/flai/internal/guard"
	"github.com/bytepunx/system-flow/flai/internal/hostapi"
	"github.com/bytepunx/system-flow/flai/internal/mergequeue"
	"github.com/bytepunx/system-flow/flai/internal/verify"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// exitVerifyUnusable is flai verify's exit status when it could not answer:
// a usage error, a story with no worktree, tiers the manifest does not
// declare validly, or a run stopped before it finished, kept apart from a
// failing step's 1, as flai test keeps them (S-0270).
const exitVerifyUnusable = 2

// verifyJobWait is how long flai verify --job waits for flai serve to take
// and answer its request: closeout's defaults, ten seconds each; a test
// shortens it.
var verifyJobWait = closeout.Wait{}

// newVerifyCmd runs a story's close-out steps and the tiers its branch
// selects in its worktree and answers one report (S-0270).
func newVerifyCmd(a *app) *cobra.Command {
	var record, last, syncOnly, fresh, all, job bool
	var maxFindings int
	var message string
	c := &cobra.Command{
		Use:   "verify <story>",
		Short: "Run a story's close-out steps and the tiers its branch selects, in its worktree, and answer one result",
		Long: `Verify a story before it goes to review, in its worktree under
.flai-cache/worktrees, from any checkout of the project. The steps run in
order, cheapest first, and the run stops at the first that fails; the steps
after it are not reached:

  rebase     no rebase is left unfinished in the worktree
  sync       the story's branch contains the main branch, but for commits
             that change only paths under the wip folder that the branch
             does not change, which a note under the step names
  narrative  the narrative's Current state and Next steps are written
  check      flai check --strict, scoped to the story, passes
  <tier>     each test and lint tier of the worktree's system-flow.yaml that
             what the branch changed against the main branch selects, run
             with CLOSE_OUT_STORY set to the story and FLAI_ROLE=verify,
             whatever role ran flai verify, so a test that runs flai's
             writes is not refused as the orchestrator (S-0311); a tier
             that another selected tier covers is not run but listed as
             "skipped tier <name>: covered by <tier>", even after a step
             before it failed (S-0342)

While the merge queue can batch for the project, a tier the manifest marks
queue is left to the queue's batch, which runs it, and not run here: the
merge-queue host action is on for the project in the configuration this flai
reads, orchestration.permissions.accept_reviews is on, and acceptance is
local. Each such tier the branch selects is listed after the rest as
"queued tier <name>: left to the merge queue", the run passes without it, and
a tier that only it covers runs. A tier marked queue whose when matches one
of the story's touches, as its claim holds them (its own and its open
tasks'), or a path its branch changed runs here all the same, and the batch
runs it again: its line names the pattern and the touch or path that chose
it, as "passed tier <name> (<duration>): ran here, its when <pattern> matches
<touch or path>", and --json carries that as its step's ran_for (ADR-0233).
--all leaves no tier to the queue: every tier the branch selects runs
(ADR-0174). When the branch selects a tier marked queue but the queue cannot
batch, a line before the notes and the outcome names those tiers and the
condition that is off, such as the host action off for the project in the
configuration this flai reads, and --json carries it as queue_off (S-0466).

When the story's last run was at the same head and base commit, and the
tiers, their commands, and those the branch selects are the same, the steps
before the tiers run again but the tiers before the first the last run showed
failed or not reached do not: each is reported reused, with the time of the
run that ran it, and the run goes on from the tier that failed (S-0341).
After a clean rebase onto a newer main branch the head and base commit
differ, but when the branch's own change has the same patch id
(git patch-id --stable of its diff from the merge base, recorded only from a
worktree with nothing uncommitted; it leaves out the files flai writes again
at a sync, design/issues/summary.md, which differ after a rebase though the
branch's own change does not), the tiers are reused all the same, up to
the first that the paths the main branch brought in between the two merge
bases select, but for wip paths the branch does not change; that tier runs
again with every tier after it. A note under the first reused tier, or under
the first tier run again, says why (ADR-0156). Any other change makes the
run full; --fresh forces one.

flai verify commits nothing. The answer is a line for each step, its state
and duration, the failing step's first findings under it, at most --max
across the run, then the check's findings outside the story, which are notes
that do not fail it, and last a line as the close-out ends:
"verify: S-nnnn passed every step", or "verify: S-nnnn stopped at <step>
(<exit status or state>)", either followed by ", <tiers> left to the merge
queue" when tiers were left to it, and then by " (reused N tiers from
<time>)" when tiers were reused, the earliest run they come from. With
--json the answer is the report, with the story, the commit and base it was verified at
and the base's commit, the merge base and the branch's patch id, when it ran,
how long it took, passed, stopped_at, the
paths the branch changed, and each step's state, duration, and findings, a
reused tier's reused_from, a skipped tier's covered_by, a tier left to the
merge queue as queued, the ran_for of a tier marked queue that ran here,
and the notes.

Each run's report is stored in the project's .flai-cache/verify, whether it
passes or not; --last prints the story's stored report and runs nothing, or
says there is none (null with --json), and exits 0 whatever it holds.

--sync-only runs the rebase and sync steps alone, prints and exits as a full
run does, and stores no report, so the stored one stays the last full run's.
The close-out's last check runs it after its commit. It is refused with
--last and with --record-issues.

--record-issues records the notes in design/issues of the story's worktree,
as flai check --story --record-issues does: each rule's notes in the open
issue whose title names the rule, once per story and notes. The close-out
commits them with the story.

Each run on a worktree with nothing uncommitted is kept in the record under
its commit as well. Unless --fresh, a run at a commit whose record holds a
passing verdict, with the worktree clean and the main branch at the
record's base commit or ahead of it only by commits the sync step passes
over, answers that verdict and runs nothing: its last line ends
"(answered from the record of <time> at <commit>; ran nothing)", and
--record-issues records nothing for it (ADR-0229).

--job runs nothing here: it asks flai serve for the story's close-out as a
job in its worktree, scripts/close-out.sh S-nnnn -m <message> where the
worktree has the script and flai verify --record-issues where it has none,
and prints the job's id as soon as flai serve takes the request; wait for
the job's verified event (wait_for_events), then read the result with
flai verify S-nnnn --last. A job of the story already running is answered
with its id, and a commit the record holds a passing verdict for is
answered from the record at once, starting no job. -m/--message is the
commit message the close-out commits with, and --all and --fresh go to the
job. A request flai serve does not take within ten seconds is withdrawn:
flai serve is not running, so run the close-out script, or flai verify where
the project has none, yourself (ADR-0229).

The exit status is 0 when every step passed, 1 when a step failed, and 2 when
flai verify could not answer: the story has no worktree, the manifest's tiers
are not valid, the notes could not be recorded, or the run was stopped before
it finished.`,
		Example: `  flai verify S-0270
  flai verify S-0270 --job -m "feat: [S-0270] what the story did"
  flai verify S-0270 --json
  flai verify S-0270 --record-issues
  flai verify S-0270 --fresh
  flai verify S-0270 --all
  flai verify S-0270 --last --json
  flai verify S-0270 --sync-only`,
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
			if syncOnly && last {
				return &exitError{code: exitVerifyUnusable, msg: "--sync-only runs the rebase and sync steps and --last runs nothing; drop one of them"}
			}
			if syncOnly && record {
				return &exitError{code: exitVerifyUnusable, msg: "--sync-only runs no check, so it has no notes to record; drop --record-issues, or --sync-only to run every step"}
			}
			if msg := jobFlagsRefusal(job, last, syncOnly, record, message); msg != "" {
				return &exitError{code: exitVerifyUnusable, msg: msg}
			}
			repo, err := a.project()
			if err != nil {
				return &exitError{code: exitVerifyUnusable, msg: err.Error()}
			}
			// the story's worktree, its stored result, and wip/ are the main
			// checkout's; a worktree mounted alone has none of them (I-0219)
			if err := repo.RequireMain("flai verify"); err != nil {
				return &exitError{code: exitVerifyUnusable, msg: err.Error()}
			}
			if last {
				return a.printLastVerify(repo, args[0], maxFindings)
			}
			if job {
				return a.askCloseOut(cmd, repo, args[0], closeout.Request{Message: message, All: all, Fresh: fresh})
			}
			// the orchestrator's container reads no host action and lacks
			// the tiers' tools, so it asks flai serve (I-0256, ADR-0278)
			if !syncOnly {
				if cfg, ok := a.orchestratorInRoleContainer(repo); ok {
					return &exitError{code: exitVerifyUnusable, msg: roleContainerVerifyRefusal(cfg, args[0])}
				}
			}
			return a.verifyStory(cmd, repo, args[0], verifyOptions{record: record, syncOnly: syncOnly, fresh: fresh, all: all, max: maxFindings})
		},
	}
	c.SetFlagErrorFunc(func(_ *cobra.Command, err error) error {
		return &exitError{code: exitVerifyUnusable, msg: err.Error()}
	})
	c.Flags().BoolVar(&record, "record-issues", false, "record the check's findings outside the story in design/issues of its worktree, once per story and findings")
	c.Flags().BoolVar(&fresh, "fresh", false, "force a full run: run every selected tier even when the last run at the same head, base, and tiers passed some")
	c.Flags().BoolVar(&all, "all", false, "leave no tier to the merge queue: run every tier the branch's changes select")
	c.Flags().BoolVar(&last, "last", false, "print the story's stored last result and run nothing")
	c.Flags().BoolVar(&syncOnly, "sync-only", false, "run only the rebase and sync steps, and store no result")
	c.Flags().IntVar(&maxFindings, "max", verify.DefaultMax, "the most findings to report across the run")
	c.Flags().BoolVar(&job, "job", false, "ask flai serve for the story's close-out as a job and print its id at once; wait for its verified event")
	c.Flags().StringVarP(&message, "message", "m", "", "with --job, the commit message the close-out commits with")
	return c
}

// verifyOptions are flai verify's flags for a run.
type verifyOptions struct {
	record bool
	// syncOnly runs only the rebase and sync steps and stores no report
	// (ADR-0135).
	syncOnly bool
	// fresh runs every selected tier, reusing none from the last run.
	fresh bool
	// all leaves no tier to the merge queue (ADR-0174).
	all bool
	max int
}

// verifyRecorded is flai verify --json with --record-issues: the report,
// and the issues its notes were recorded in.
type verifyRecorded struct {
	verify.Report
	Recorded []recordedIssue `json:"recorded"`
}

// verifyStory runs verify.Verify for the story in its worktree, or with
// --sync-only verify.SyncOnly, records its notes when asked, prints the
// report, and answers the exit status.
func (a *app) verifyStory(cmd *cobra.Command, repo *workitem.Repo, story string, o verifyOptions) error {
	id, worktree, err := storyWorktree(repo, story)
	if err != nil {
		return &exitError{code: exitVerifyUnusable, msg: err.Error()}
	}
	opts := verify.StoryOptions{
		Story: id, Project: repo, Worktree: worktree, Git: a.runner, Fresh: o.fresh,
		RunOptions: verify.RunOptions{Max: o.max, Now: a.now},
	}
	if !o.all && !o.syncOnly {
		// the report says why, when the queue cannot batch, so that a
		// configuration without the host action is seen (S-0466)
		opts.QueueOff = a.queueOff(mainRootOf(repo))
		opts.LeaveToQueue = opts.QueueOff == ""
	}
	var rep verify.Report
	stopped := false
	if o.syncOnly {
		rep, err = verify.SyncOnly(opts)
	} else {
		if opts.Tiers, err = verify.CheckoutTiersFor(worktree, a.runner, buildinfo.Version); err != nil {
			return &exitError{code: exitVerifyUnusable, msg: fmt.Sprintf("verify %s: %v", id, err)}
		}
		ctx, stop := signal.NotifyContext(cmd.Context(), os.Interrupt, syscall.SIGTERM)
		defer stop()
		rep, err = verify.Verify(ctx, opts)
		stopped = err != nil && ctx.Err() != nil
	}
	if err != nil && rep.Story == "" {
		return &exitError{code: exitVerifyUnusable, msg: err.Error()}
	}
	a.logger().Debug("story verified", "component", "cmd", "story", rep.Story, "steps", len(rep.Steps), "notes", len(rep.Notes), "passed", rep.Passed)
	var recorded []recordedIssue
	// a verdict from the record ran nothing now, and its notes were recorded
	// when it ran, by the close-out that ran it
	if o.record && !rep.FromRecord {
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
// its state, name, and duration, a tier's name after "tier", a reused tier's
// time of the run it comes from in place of its duration, a skipped tier's
// line naming the tier that covers it in place of a duration, after the
// duration of a tier marked queue that ran here why it was not left to the
// merge queue (ADR-0233), a passing step's
// note and the failing step's findings indented under it, then why the tiers
// marked queue ran here when the merge queue cannot batch (S-0466), at most
// maxNotes of the notes outside the story, the issues they were recorded in,
// and last the outcome as the close-out's last line has it.
func verifyText(rep verify.Report, recorded []recordedIssue, stopped bool, maxNotes int) string {
	var b strings.Builder
	for _, s := range rep.Steps {
		name := s.Name
		if s.Tier {
			name = "tier " + s.Name
		}
		ranFor := ""
		if s.RanFor != "" {
			ranFor = ": ran here, " + s.RanFor
		}
		switch {
		case s.State == verify.NotReached:
			fmt.Fprintf(&b, "%s %s\n", s.State, name)
		case s.State == verify.Skipped:
			fmt.Fprintf(&b, "%s %s: covered by %s\n", s.State, name, s.CoveredBy)
		case s.State == verify.Queued:
			fmt.Fprintf(&b, "%s %s: left to the merge queue\n", s.State, name)
		case s.State == verify.Reused && s.ReusedFrom != nil:
			fmt.Fprintf(&b, "%s %s (from %s)%s\n", s.State, name, s.ReusedFrom.UTC().Format(time.RFC3339), ranFor)
		default:
			fmt.Fprintf(&b, "%s %s (%s)%s\n", s.State, name, s.Duration, ranFor)
		}
		if s.Note != "" {
			b.WriteString("    " + strings.Join(strings.Fields(s.Note), " ") + "\n")
		}
		b.WriteString(verify.StepFindings(s, "    "))
	}
	if rep.QueueOff != "" {
		b.WriteString(strings.Join(strings.Fields(rep.QueueOff), " ") + "\n")
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
// story passed every step, or the step it stopped at and why, then the tiers
// it left to the merge queue, if any, the tiers it reused and the earliest
// run they come from, if any, which the close-out reads from the line's end,
// and last, for a verdict from the record, that it came from the record, when
// that run was, and that nothing ran (ADR-0229).
func verifyOutcome(rep verify.Report, stopped bool) string {
	line := verifyResult(rep, stopped) + queuedText(rep) + reusedText(rep)
	if text := rep.RecordText(); text != "" {
		line += " (" + text + ")"
	}
	return line
}

// queuedText is ", <tiers> left to the merge queue" for a report that left
// tiers to the queue, named in the order it lists them (ADR-0174); "" for
// one that left none.
func queuedText(rep verify.Report) string {
	var names []string
	for _, s := range rep.Steps {
		if s.State == verify.Queued {
			names = append(names, s.Name)
		}
	}
	switch len(names) {
	case 0:
		return ""
	case 1:
		return ", " + names[0] + " left to the merge queue"
	}
	return ", " + strings.Join(names[:len(names)-1], ", ") + " and " + names[len(names)-1] + " left to the merge queue"
}

// queueBatches reports whether the merge queue can batch for the project
// whose main checkout is root, as flai serve builds a batch (ADR-0174): the
// merge-queue host action is on for it in the configuration this flai reads,
// and the main checkout's manifest lets the queue run. A configuration that
// is missing turns nothing on and is not created; one or a manifest that
// cannot be read is warned of and turns nothing on, so every tier runs.
func (a *app) queueBatches(root string) bool {
	return a.queueOff(root) == ""
}

// queueOff is why the merge queue cannot batch for the project whose main
// checkout is root, the condition queueBatches finds off and what turns it
// on, so that flai verify can say why it runs the tiers marked queue
// (S-0466); "" while the queue can batch.
func (a *app) queueOff(root string) string {
	path := config.ResolvePath(a.configPath)
	on, err := a.hostActionOn(hostapi.ActionMergeQueue, root)
	if err != nil {
		a.logger().Warn("configuration unreadable, so no tier is left to the merge queue", "component", "verify", "err", err.Error())
		return fmt.Sprintf("the configuration %s cannot be read: %v", path, err)
	}
	if !on {
		return fmt.Sprintf("the %s host action is off for %s in the configuration %s; turn it on there with %s in a shell on the host, or run flai with --config naming the configuration flai serve reads", hostapi.ActionMergeQueue, root, path, hostapi.EnableCommand(hostapi.ActionMergeQueue))
	}
	repo, err := workitem.Open(root)
	if err != nil {
		a.logger().Warn("project unreadable, so no tier is left to the merge queue", "component", "verify", "root", root, "err", err.Error())
		return fmt.Sprintf("the project at %s cannot be read: %v", root, err)
	}
	switch {
	case repo.Manifest.Acceptance.PullRequest():
		return "acceptance.mode is pull-request, under which the merge queue does not run (ADR-0173)"
	case !mergequeue.Batches(true, repo.Manifest):
		return "orchestration.permissions.accept_reviews is off in the main checkout's system-flow.yaml"
	}
	return ""
}

// orchestratorInRoleContainer reports whether this flai runs as the
// orchestrator with the configuration it reads, --config or FLAI_CONFIG,
// inside the project's role cache folders, <main checkout>/.flai-cache/agents,
// as a role's container gives it (ADR-0247), and that configuration's path.
// No host action can be read there: hostActionOn would find every one off,
// so the orchestrator's acceptance and a verify that runs tiers refuse in
// place of acting on that (I-0256, ADR-0278).
func (a *app) orchestratorInRoleContainer(repo *workitem.Repo) (string, bool) {
	if os.Getenv("FLAI_ROLE") != guard.RoleOrchestrate {
		return "", false
	}
	return a.configInAgentCache(repo, "")
}

// storyInContainer reports whether this flai runs in a story's agent
// container: FLAI_STORY names the story, and the configuration it reads lies
// in that story's cache folder, <main checkout>/.flai-cache/agents/<story>, as
// the container gives it (harness.ContainerCache), and that story. The main
// checkout's .git is mounted read-only there (ADR-0246), though its
// system-flow.yaml is mounted too (ADR-0253), so Repo.RequireMain lets
// commands through that write to it (I-0272).
func (a *app) storyInContainer(repo *workitem.Repo) (string, bool) {
	story := os.Getenv("FLAI_STORY")
	if story == "" {
		return "", false
	}
	_, ok := a.configInAgentCache(repo, story)
	return story, ok
}

// configInAgentCache returns the configuration this flai reads, --config or
// FLAI_CONFIG, and whether it lies in the cache folder that the container of
// subject, a story's ID or a role, mounts under the main checkout's
// .flai-cache/agents, or in any of those folders when subject is empty.
func (a *app) configInAgentCache(repo *workitem.Repo, subject string) (string, bool) {
	path := config.ResolvePath(a.configPath)
	return path, pathWithin(filepath.Join(repo.CacheDir(), "agents", subject), path)
}

// refuseSyncInStoryContainer is the refusal of what, flai stream sync or
// flai task done, in a story's agent container, before any git call: its
// rebase and trial merge write to the main checkout's .git, which is
// read-only there, and the MCP tools task_done and verify sync on the host
// instead (I-0272, ADR-0274). Nil outside such a container.
func (a *app) refuseSyncInStoryContainer(repo *workitem.Repo, what string) error {
	story, ok := a.storyInContainer(repo)
	if !ok {
		return nil
	}
	return fmt.Errorf("rule: %s cannot run in %s's agent container: the main checkout's .git is mounted read-only there (ADR-0246), so no rebase or trial merge can be written; close the task with the MCP tool task_done, which closes a task and syncs on the host, or run the story's close-out with the MCP tool verify, which syncs on the host (ADR-0274)", what, story)
}

// pathWithin reports whether path lies below dir, both made absolute and
// clean, so that neither a sibling sharing dir's prefix nor a path that
// climbs out with .. counts.
func pathWithin(dir, path string) bool {
	d, err := filepath.Abs(dir)
	if err != nil {
		return false
	}
	p, err := filepath.Abs(path)
	if err != nil {
		return false
	}
	rel, err := filepath.Rel(d, p)
	if err != nil {
		return false
	}
	return rel != "." && rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator)) && !filepath.IsAbs(rel)
}

// roleContainerVerifyRefusal is flai verify's refusal to run story's tiers
// as the orchestrator in a role's container, whose configuration is cfg.
func roleContainerVerifyRefusal(cfg, story string) string {
	return fmt.Sprintf("the host's actions cannot be read in a role's container (its configuration is %s), so flai verify %s cannot tell which tiers the merge queue runs, and the container lacks the tiers' tools: ask flai serve for the close-out with the MCP tool verify, or flai verify %s --job, read the verdict with flai verify %s --last, and accept with the MCP tool accept, which runs on the host (ADR-0278)", cfg, story, story, story)
}

// hostActionOn reports whether a host action is on for the project at root
// in the configuration this flai reads, without creating it: a missing one
// turns nothing on.
func (a *app) hostActionOn(action, root string) (bool, error) {
	path := config.ResolvePath(a.configPath)
	data, err := os.ReadFile(path)
	if errors.Is(err, fs.ErrNotExist) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("read the configuration %s: %w", path, err)
	}
	cfg, err := config.Parse(data)
	if err != nil {
		return false, fmt.Errorf("parse the configuration %s: %w; fix it, or run with --config naming another", path, err)
	}
	return cfg.ActionEnabled(action, root), nil
}

// verifyResult is the last line's outcome: the story passed every step, or
// the step it stopped at and why.
func verifyResult(rep verify.Report, stopped bool) string {
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

// reusedText is " (reused N tiers from <time>)" for a report that reused
// tiers, the time the earliest run they come from began, in UTC; "" for one
// that ran every tier it reached.
func reusedText(rep verify.Report) string {
	n, from := rep.Reused()
	if n == 0 {
		return ""
	}
	tiers := "tiers"
	if n == 1 {
		tiers = "tier"
	}
	if from.IsZero() {
		return fmt.Sprintf(" (reused %d %s)", n, tiers)
	}
	return fmt.Sprintf(" (reused %d %s from %s)", n, tiers, from.UTC().Format(time.RFC3339))
}

// shortCommit is a commit's first seven characters, as git abbreviates it.
func shortCommit(commit string) string {
	if len(commit) > 7 {
		return commit[:7]
	}
	return commit
}

// jobFlagsRefusal is why flai verify's flags cannot go together with --job,
// or -m/--message without it; "" when they can.
func jobFlagsRefusal(job, last, syncOnly, record bool, message string) string {
	switch {
	case !job && message != "":
		return "--message is the commit message of the close-out job, so it goes with --job; add --job, or drop --message"
	case job && last:
		return "--job asks flai serve for a close-out and --last prints the stored result; drop one of them"
	case job && syncOnly:
		return "--job asks flai serve for the whole close-out and --sync-only runs only the rebase and sync steps here; drop one of them"
	case job && record:
		return "--job's close-out records the notes outside the story itself; drop --record-issues"
	}
	return ""
}

// askCloseOut asks flai serve for the story's close-out as a job (ADR-0229)
// and prints its answer: the job it started, or the one of the story that
// runs, with where to wait and read; the verdict the record holds, as a
// verify that passed or stopped answers it; or why it ran nothing, as an
// unusable verify.
func (a *app) askCloseOut(cmd *cobra.Command, repo *workitem.Repo, story string, req closeout.Request) error {
	id, _, err := storyWorktree(repo, story)
	if err != nil {
		return &exitError{code: exitVerifyUnusable, msg: err.Error()}
	}
	req.Story, req.By = id, a.writer()
	ctx, stop := signal.NotifyContext(cmd.Context(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	ans, err := closeout.Ask(ctx, repo.CacheDir(), req, verifyJobWait)
	if err != nil {
		if errors.Is(err, context.Canceled) || errors.Is(err, closeout.ErrServeNotTaking) {
			return &exitError{code: exitVerifyUnusable, msg: err.Error()}
		}
		return &exitError{code: exitVerifyUnusable, msg: fmt.Sprintf("ask flai serve for %s's close-out: %v", id, err)}
	}
	a.logger().Debug("close-out asked", "component", "cmd", "story", id, "job_id", ans.Job, "started", ans.Started, "recorded", ans.Recorded, "refused", ans.Refused())
	if a.jsonOut {
		if err := a.printJSON(ans); err != nil {
			return err
		}
	} else if !ans.Refused() {
		fmt.Fprintln(a.out, closeOutAnswerText(ans))
	}
	switch {
	case ans.Refused():
		msg := ans.Reason
		if a.jsonOut {
			msg = ""
		}
		return &exitError{code: exitVerifyUnusable, msg: msg}
	case ans.Recorded && ans.Verdict != closeout.VerdictPassed:
		return &exitError{code: 1}
	}
	return nil
}

// closeOutAnswerText is flai serve's answer to a close-out request as flai
// verify --job prints it: the job it started or found running, what to wait
// for, and where its output is; or the verdict the record holds.
func closeOutAnswerText(ans closeout.Answer) string {
	if ans.Recorded {
		verdict := ans.Verdict
		if ans.StoppedAt != "" {
			verdict += " at " + ans.StoppedAt
		}
		return fmt.Sprintf("%s %s from the record of %s at %s; ran nothing", ans.Story, verdict, ans.RanAt.UTC().Format(time.RFC3339), shortCommit(ans.Commit))
	}
	state := "started"
	if ans.Running() {
		state = "is already running"
	}
	return fmt.Sprintf("close-out job %s %s for %s: wait for its verified event (wait_for_events), then read it with flai verify %s --last; its output is in %s", ans.Job, state, ans.Story, ans.Story, ans.Log)
}
