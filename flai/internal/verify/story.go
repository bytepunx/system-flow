package verify

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/atomicfile"
	"github.com/bytepunx/system-flow/flai/internal/check"
	"github.com/bytepunx/system-flow/flai/internal/closeout"
	"github.com/bytepunx/system-flow/flai/internal/execx"
	"github.com/bytepunx/system-flow/flai/internal/gitver"
	"github.com/bytepunx/system-flow/flai/internal/issues"
	"github.com/bytepunx/system-flow/flai/internal/kluster"
	"github.com/bytepunx/system-flow/flai/internal/storygit"
	"github.com/bytepunx/system-flow/flai/internal/verifyrecord"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// The steps Verify runs before the tiers, in order, cheapest first.
const (
	StepRebase    = "rebase"
	StepSync      = "sync"
	StepNarrative = "narrative"
	StepCheck     = "check"
	// StepKluster reads the story's kluster namespace's record against its
	// head and runs nothing; it is a step only while the main checkout's
	// manifest has a kluster: block (ADR-0292 §6).
	StepKluster = "kluster"
)

// StoryEnv names the variable a story's tiers run with set to the story, so
// that the scripts and tests that check the project scope themselves to it
// (ADR-0085).
const StoryEnv = "CLOSE_OUT_STORY"

// MergeQueueStoriesEnv names the variable a merge queue batch's tiers run
// with set to the stories the batch verifies, comma-separated in batch
// order, so that the scripts and tests that check the project scope
// themselves to them; unlike StoryEnv it is no close-out, and leads them to
// record no issue.
const MergeQueueStoriesEnv = "MERGE_QUEUE_STORIES"

// narrativeSections are the narrative's sections that must be written
// before a story goes to review.
var narrativeSections = []string{"Current state", "Next steps"}

// StoryOptions are what Verify needs to answer for a story.
type StoryOptions struct {
	// Story is the story's ID, as typed: S-12 is S-0012.
	Story string
	// Project is the project, opened at any of its checkouts.
	Project *workitem.Repo
	// Worktree is the story's worktree; "" is where flai stream open puts it.
	Worktree string
	// Base is the main branch the story's branch must contain and its
	// changes are counted against; "" is the branch the main checkout has.
	Base string
	// Tiers are the project's tiers, cheapest first, as the worktree's
	// manifest declares them.
	Tiers []Tier
	// Git runs git; nil is execx.System.
	Git execx.Runner
	// Fresh runs every selected tier, whatever the story's last report
	// would let the run reuse.
	Fresh bool
	// LeaveToQueue leaves each selected tier marked Queue to the merge
	// queue's batch, which runs it: it is not run but reported Queued, after
	// the tiers that run, unless its When matches the story's claim or a path
	// its branch changed, when it runs with them (ADR-0233). Set it while the
	// merge queue can batch for the project and the caller did not ask for
	// every tier (ADR-0174).
	LeaveToQueue bool
	// QueueOff is why the merge queue cannot batch for the project, when
	// that alone keeps the caller from leaving tiers to it: the report's
	// QueueOff then names the selected tiers marked Queue that the run did
	// not leave to the queue, with it. "" when the queue can batch, or the
	// caller asked for every tier (S-0466).
	QueueOff string
	// Suites are the suites whose tiers the run selects, by each tier's
	// SuiteOrName, as StorySuites gives those the story's next moves need
	// (ADR-0312 §2): nil selects every suite's, as flai verify --all does; a
	// list that is not nil selects only the tiers of its suites, which the
	// report's Suites names, and a tier of it that only a tier of another
	// suite covers runs itself, as SelectStorySuites selects them.
	Suites []string
	// KlusterOn and ChecksOn say whether the kluster and checks host actions
	// are on for the project, as the kluster step reads the story's
	// namespace's record under them (ADR-0292 §6).
	KlusterOn, ChecksOn bool
	// StepStarted is called with the name of each step before the tiers as
	// it starts, and Started with each tier's as its command starts, as
	// TierOptions has it; Ended is called with each step and tier as it ends,
	// a tier's with Tier set. A tier that runs ends with its name and its
	// state as its command comes back, as RunOptions' TierEnded tells it,
	// which Verify sets, with TierStarted, from these. A tier that ends
	// without running, reused, skipped for the tier that covers it, or left
	// to the merge queue, ends with its whole step and no start; a step or
	// tier not reached calls neither. Parallel tiers run beside each other,
	// and the hooks are called one at a time. Nil calls nothing.
	StepStarted func(step string)
	Started     func(tier string)
	Ended       func(Step)
	// NoStore stores the report nowhere: neither as the story's last
	// (SaveReport), which the acceptance and the merge queue read as its
	// close-out, nor in the record (KeepAt), so that a run of some suites
	// never stands for the whole close-out. A passing verdict the record
	// holds at the head is still answered, when its fingerprint is this
	// run's, as it is for the same tiers.
	NoStore bool
	// RunOptions are the tiers' run's; FS is the worktree, read for the
	// tiers' packages, and Now times every step and stamps the report.
	RunOptions
}

// Report is the answer to Verify: pass, or the step that failed and why.
type Report struct {
	Story string `json:"story"`
	// Commit is the worktree's HEAD when it was verified.
	Commit string `json:"commit"`
	// Base is the main branch the story was verified against, and
	// BaseCommit the commit it had then; "" when git could not read it.
	Base       string `json:"base"`
	BaseCommit string `json:"base_commit,omitempty"`
	// MergeBase is where the branch left Base, git merge-base HEAD Base;
	// "" when git could not tell.
	MergeBase string `json:"merge_base,omitempty"`
	// Patch is the stable patch id of the branch's own change, what it
	// changed since MergeBase but for the files a sync's rebase writes again
	// from others, the issues summary (ADR-0098); "" when the worktree had
	// anything uncommitted, the branch changed nothing else, or git could
	// not tell (ADR-0156).
	Patch string `json:"patch,omitempty"`
	// Fingerprint hashes what selects and runs the tiers: the project's
	// tiers and the selected ones' commands, those run and those left to the
	// merge queue alike, so that a run that leaves tiers to the queue and one
	// that runs them have the same (ADR-0174, I-0203). Each of SuiteRuns has
	// the fingerprint of the project's tiers and its own selected ones: a
	// later run with the same reuses the tiers of that suite this one passed
	// when it is at the same Commit and BaseCommit (S-0341), or, rebased, has
	// the same Patch, but for the tiers that what Base changed between the
	// two MergeBases sends back (ADR-0156, ADR-0312 §7).
	Fingerprint string    `json:"fingerprint,omitempty"`
	RanAt       time.Time `json:"ran_at"`
	// DurationMS and Duration are how long the whole run took.
	DurationMS int64  `json:"duration_ms"`
	Duration   string `json:"duration,omitempty"`
	Passed     bool   `json:"passed"`
	// StoppedAt names the step that failed; "" when every step passed.
	StoppedAt string `json:"stopped_at,omitempty"`
	// Paths are the root-relative files the branch changed against Base,
	// which the tiers were selected for.
	Paths []string `json:"paths,omitempty"`
	// Suites are the suites Tiers ran the tiers of, as TierOptions.Suites
	// named them, so that a caller such as the merge queue's batch can name
	// them (ADR-0312), or that Verify selected the tiers of, as
	// StoryOptions.Suites named them; nil for a run of every suite. The JSON
	// leaves out an empty list as it leaves out nil.
	Suites []string `json:"suites,omitempty"`
	Steps  []Step   `json:"steps"`
	// SuiteRuns are what the run recorded of each suite whose tiers it
	// selected: its tiers, how each ended, and the head and base they passed
	// at (ADR-0312 §7). The story's last report, and its verdict at a commit,
	// keep too the suites an earlier run recorded that a later one did not
	// run (SaveReport, KeepAt). None in a report stored before suites were.
	SuiteRuns []verifyrecord.Suite `json:"suite_runs,omitempty"`
	// Notes are the check's findings outside the story, which do not fail
	// it (S-0249); flai verify --record-issues records them.
	Notes []Note `json:"notes,omitempty"`
	// QueueOff says why the tiers marked queue that the run selected ran
	// here rather than being left to the merge queue: it names them and the
	// condition that keeps the queue from batching, as "tiers integration
	// and smoke ran here, not left to the merge queue: <why>". "" when no
	// tier ran here for that reason: the run left them to the queue,
	// selected none, or was asked for every tier (S-0466).
	QueueOff string `json:"queue_off,omitempty"`
	// Clean says the worktree had nothing uncommitted, untracked files
	// included, when the run began, so that it verified exactly Commit and
	// the record keeps it under Commit (ADR-0229).
	Clean bool `json:"clean,omitempty"`
	// FromRecord says the report is the verdict the record holds at Commit,
	// answered without running any step: RanAt and the durations are those of
	// the run the record kept (ADR-0229).
	FromRecord bool `json:"from_record,omitempty"`
	// Kluster is what the kluster step read of the story's namespace: its
	// last deploy and checks run, their commits and outcomes, and whether
	// each is at the head (ADR-0292 §6); nil while none applies, the step was
	// not reached, or the project has no kluster: block, and in a report
	// stored before it was kept.
	Kluster *kluster.Evidence `json:"kluster,omitempty"`
	// Moved says the report's Commit is not the head of the story's branch
	// that LastReportAt was given: the report is from before the head moved,
	// as by a sync (I-0384). It is not stored; only LastReportAt sets it.
	Moved bool `json:"moved,omitempty"`
	// Running is the story's running close-out job, its id, start, head, and
	// step, when the job runs at another commit than the report's, a job that
	// has not written its head taken as running at the branch's head (I-0384).
	// It is not stored; only LastReportAt sets it.
	Running *closeout.Progress `json:"running,omitempty"`
}

// Step is what one step of Verify did.
type Step struct {
	Name string `json:"name"`
	// Tier says the step ran a test or lint tier, named Name.
	Tier       bool      `json:"tier,omitempty"`
	State      State     `json:"state"`
	DurationMS int64     `json:"duration_ms"`
	Duration   string    `json:"duration,omitempty"`
	Findings   []Finding `json:"findings,omitempty"`
	// Note is what a step says about why it ran as it did: what a step
	// that passed passed over, as sync names the commits of Base it did not
	// need (ADR-0135), or, after a rebase, that a tier was reused across it
	// or that the paths Base brought in sent it back (ADR-0156). It is the
	// step's own, not a check note, and nothing records it as an issue.
	Note string `json:"note,omitempty"`
	// Omitted counts the findings past the run's cap.
	Omitted int `json:"omitted,omitempty"`
	// Command, Dir, and ExitCode are a tier's, as TierResult has them.
	Command  []string `json:"command,omitempty"`
	Dir      string   `json:"dir,omitempty"`
	ExitCode *int     `json:"exit_code,omitempty"`
	// ReusedFrom is when the run that ran a Reused tier began.
	ReusedFrom *time.Time `json:"reused_from,omitempty"`
	// CoveredBy names the tier that covers a skipped tier, as TierResult
	// has it (S-0342).
	CoveredBy string `json:"covered_by,omitempty"`
	// RanFor says why a tier marked Queue that ran or was reused was not
	// left to the merge queue, as Selected has it: "its when <pattern>
	// matches <touch or path>" (ADR-0233).
	RanFor string `json:"ran_for,omitempty"`
	// Log is the file holding a tier's whole output, as TierResult has it.
	Log string `json:"log,omitempty"`
}

// Reused counts the report's reused tiers and answers the earliest run they
// come from; 0 and the zero time when it ran every tier it reached.
func (r Report) Reused() (int, time.Time) {
	n, from := 0, time.Time{}
	for _, s := range r.Steps {
		if s.State != Reused {
			continue
		}
		n++
		if s.ReusedFrom != nil && (from.IsZero() || s.ReusedFrom.Before(from)) {
			from = *s.ReusedFrom
		}
	}
	return n, from
}

// Note is a check finding outside the story.
type Note struct {
	Rule  string `json:"rule"`
	Level string `json:"level"`
	// Path is relative to the worktree, or for wip/ to the main checkout.
	Path    string `json:"path,omitempty"`
	Line    int    `json:"line,omitempty"`
	Message string `json:"message"`
}

// Verify runs a story's close-out checks in its worktree and answers one
// report: no rebase is unfinished, the branch contains Base but for commits
// that change only wip paths the branch does not change, the narrative's
// Current state and Next steps are written, flai check --strict scoped to the
// story passes, while the main checkout's manifest has a kluster: block the
// story's namespace's deploy and checks run did not fail at its head (the
// kluster step, which runs nothing and keeps what it read as the report's
// Kluster), and then each tier of Suites, nil for every suite, that the
// branch's changes select (SelectStorySuites) passes, run with
// CLOSE_OUT_STORY set to the story. It stops at the first
// step that fails; the steps after it are not reached, but for a tier that
// another selected tier covers, which is skipped either way. With
// LeaveToQueue, the selected tiers marked Queue are not run but reported
// Queued after the rest, and the run passes without them (ADR-0174), but for
// those whose When matches the story's claim, as Holds.Claim reads it from
// the project, or a path its branch changed, which run with the rest, their
// steps saying why in RanFor (ADR-0233); a claim the project cannot give is
// taken as empty. When the story's
// last report records a tier's suite at the same head and base commit, with
// the same fingerprint and tiers, and Fresh is not set, the tiers before the
// first it shows failed, not reached, or left to the merge queue are not run
// again but reported Reused, with the time of the run that ran them; the
// steps before the tiers always run (S-0341, ADR-0312 §7).
// LeaveToQueue does not change the fingerprint, so a run without it reuses
// the tiers a run with it passed and runs those it left to the queue, and
// the other way round (I-0203). When the head or base commit
// differ, as after a clean rebase, they are reused only when the branch's own
// change has the same patch id, and only up to the first tier that the paths
// the main branch brought in select, which runs again with every tier after
// it (ADR-0156). The report's SuiteRuns record what the run did of each suite.
// It commits nothing and records
// no issue. The report, whole, is stored as the story's last (SaveReport)
// whether it passes or not, and, when the worktree was clean as it began and
// ctx did not end, kept in the record as the verdict at its head (KeepAt);
// with NoStore it is stored nowhere.
// Unless Fresh is set, a passing verdict the record holds at the head is
// answered first, as Recorded answers it, marked FromRecord: then no step
// runs and nothing is stored (ADR-0229). It returns an error, with no
// report, when the story, its worktree, or its changes cannot be read, and,
// with the report so far, when ctx ends during a tier or the report cannot
// be stored.
func Verify(ctx context.Context, opts StoryOptions) (Report, error) {
	v, err := newStoryRun(opts)
	if err != nil {
		return Report{}, err
	}
	if !opts.Fresh {
		// a record that cannot be read is answered by running, which keeps
		// the run in its place
		if rep, ok, _ := v.recorded(); ok {
			return rep, nil
		}
	}
	rep, runErr := v.run(ctx)
	if opts.NoStore {
		return rep, runErr
	}
	if err := SaveReport(opts.Project, rep); err != nil {
		return rep, errors.Join(runErr, err)
	}
	if rep.Clean && runErr == nil {
		if err := KeepAt(opts.Project, rep.Commit, rep); err != nil {
			return rep, err
		}
	}
	return rep, runErr
}

// SyncOnly runs only Verify's rebase and sync steps for a story, in its
// worktree, and answers their report, which it does not store, so the
// story's last report stays the last full run's: what flai verify
// --sync-only answers (ADR-0135). It returns an error, with no report, when
// the story, its worktree, or its changes cannot be read.
func SyncOnly(opts StoryOptions) (Report, error) {
	v, err := newStoryRun(opts)
	if err != nil {
		return Report{}, err
	}
	start := v.opts.Now()
	rep := v.newReport(start)
	v.runSteps(&rep, v.syncSteps())
	rep.DurationMS, rep.Duration = took(v.opts.Now().Sub(start))
	return rep, nil
}

// storyRun is one Verify call with its options resolved.
type storyRun struct {
	opts     StoryOptions
	story    string
	worktree string
	base     string
	max      int
	// commit and paths are the worktree's HEAD and what it changed, and
	// baseCommit the commit base has, "" when git cannot read it.
	commit     string
	baseCommit string
	paths      []string
	// mergeBase and patch are where the branch left base and the patch id
	// of its own change, as Report has them; Verify reads them.
	mergeBase string
	patch     string
	// clean says the worktree had nothing uncommitted when it was read.
	clean bool
	// note is what the step running says about what it passed over, which
	// runSteps puts on its result.
	note string
}

// newStoryRun resolves opts: the story's canonical ID, its worktree, the
// main branch, and the defaults; and reads the commit the worktree has and
// the paths it changed against the main branch.
func newStoryRun(opts StoryOptions) (*storyRun, error) {
	if opts.Project == nil {
		return nil, errors.New("verify a story: no project given; open the project first")
	}
	it, err := opts.Project.Get(opts.Story)
	if err != nil {
		return nil, fmt.Errorf("verify %s: %w; name a story of this project", opts.Story, err)
	}
	if it.Type != workitem.Story {
		return nil, fmt.Errorf("verify %s: it is a %s; name a story (S-nnnn)", it.ID, it.Type)
	}
	if opts.Git == nil {
		opts.Git = execx.System{}
	}
	if opts.Now == nil {
		opts.Now = time.Now
	}
	v := &storyRun{opts: opts, story: it.ID, worktree: opts.Worktree, base: opts.Base, max: opts.Max}
	if v.worktree == "" {
		v.worktree = opts.Project.WorktreePath(it.ID)
	}
	if st, err := os.Stat(v.worktree); err != nil || !st.IsDir() {
		return nil, fmt.Errorf("verify %s: it has no worktree at %s; open it with flai stream open %s", it.ID, v.worktree, it.ID)
	}
	if v.base == "" {
		if v.base, err = storygit.MainBranch(opts.Git, mainRoot(opts.Project)); err != nil {
			return nil, fmt.Errorf("verify %s: find the main branch its branch must contain: %w", it.ID, err)
		}
	}
	if v.max <= 0 {
		v.max = DefaultMax
	}
	if v.opts.FS == nil {
		v.opts.FS = os.DirFS(v.worktree)
	}
	commit, err := opts.Git.Run(v.worktree, "git", "rev-parse", "HEAD")
	if err != nil {
		return nil, fmt.Errorf("verify %s: read the commit its worktree %s has checked out: %w", v.story, v.worktree, err)
	}
	v.commit = strings.TrimSpace(commit)
	// a base git cannot read fails the sync step, and leaves nothing to resume
	if base, err := opts.Git.Run(v.worktree, "git", "rev-parse", "--verify", "--quiet", v.base+"^{commit}"); err == nil {
		v.baseCommit = strings.TrimSpace(base)
	}
	if v.paths, err = ChangedPaths(opts.Git, v.worktree, v.base); err != nil {
		return nil, fmt.Errorf("verify %s: %w", v.story, err)
	}
	if v.clean, err = v.worktreeClean(); err != nil {
		return nil, fmt.Errorf("verify %s: %w", v.story, err)
	}
	return v, nil
}

// storyStep is a step of Verify before the tiers: its name, and what finds
// why it fails, none when it passes.
type storyStep struct {
	name string
	run  func() []Finding
}

// syncSteps are the steps that say whether the branch is ready to verify:
// no rebase left unfinished, and the main branch contained.
func (v *storyRun) syncSteps() []storyStep {
	return []storyStep{{StepRebase, v.rebase}, {StepSync, v.sync}}
}

// newReport is a report begun at start that has passed so far.
func (v *storyRun) newReport(start time.Time) Report {
	return Report{
		Story: v.story, Commit: v.commit, Base: v.base, BaseCommit: v.baseCommit, MergeBase: v.mergeBase, Patch: v.patch,
		RanAt: start.UTC(), Passed: true, Paths: v.paths, Suites: slices.Clone(v.opts.Suites), Steps: []Step{}, Clean: v.clean,
	}
}

// runSteps runs steps in order onto rep; once one fails, those after it are
// not reached.
func (v *storyRun) runSteps(rep *Report, steps []storyStep) {
	for _, s := range steps {
		if !rep.Passed {
			rep.Steps = append(rep.Steps, Step{Name: s.name, State: NotReached})
			continue
		}
		if v.opts.StepStarted != nil {
			v.opts.StepStarted(s.name)
		}
		began := v.opts.Now()
		v.note = ""
		found := s.run()
		step := Step{Name: s.name, State: Passed, Note: v.note}
		step.DurationMS, step.Duration = took(v.opts.Now().Sub(began))
		if len(found) > 0 {
			step.State, step.Note = Failed, ""
			n := min(v.max, len(found))
			step.Findings, step.Omitted = found[:n], len(found)-n
			rep.Passed, rep.StoppedAt = false, s.name
		}
		rep.Steps = append(rep.Steps, step)
		v.ended(step)
	}
}

// ended calls the Ended hook, if any, with step.
func (v *storyRun) ended(step Step) {
	if v.opts.Ended != nil {
		v.opts.Ended(step)
	}
}

// endedUnrun calls the Ended hook with each of steps that ended without
// running: reused, skipped, or queued; one not reached did not end, and one
// that ran was told as its command came back.
func (v *storyRun) endedUnrun(steps []Step) {
	for _, step := range steps {
		if step.State == Reused || step.State == Skipped || step.State == Queued {
			v.ended(step)
		}
	}
}

// generated are the paths of the files a sync's rebase writes again from
// others (ADR-0098, ADR-0228), which the branch's patch id leaves out: the
// issues summary and the ADR index, written from the issue and ADR files,
// which stay in the patch, and main's, whose paths select the tiers main
// sends back (ADR-0156).
func (v *storyRun) generated() []string {
	var paths []string
	for _, g := range issues.Generated(v.opts.Project, v.story, v.opts.Now) {
		paths = append(paths, g.Path)
	}
	return paths
}

// run runs the steps in order and stops at the first that fails. It returns
// ctx's error, with the report, when ctx ends during a tier.
func (v *storyRun) run(ctx context.Context) (Report, error) {
	start := v.opts.Now()
	// what git cannot tell is left "", which reuses nothing across a
	// rebase; the steps below find what is wrong with the branch
	v.mergeBase, v.patch, _ = BranchPatch(v.opts.Git, v.worktree, v.base, v.generated())
	rep := v.newReport(start)
	steps := append(v.syncSteps(),
		storyStep{StepNarrative, v.narrative},
		storyStep{StepCheck, func() []Finding { return v.check(&rep) }},
	)
	if s, ok := v.klusterStep(&rep); ok {
		steps = append(steps, s)
	}
	v.runSteps(&rep, steps)
	selected, queued, fp := v.selection()
	listed := inListOrder(v.opts.Tiers, selected, queued)
	if !v.opts.LeaveToQueue {
		rep.QueueOff = queueOffText(selected, v.opts.QueueOff)
	}
	rep.Fingerprint = fp
	var runErr error
	if rep.Passed {
		reused, sentBack := v.reusable(selected, listed)
		for i := range reused {
			// why this run, not the last, did not leave a queue tier out
			reused[i].RanFor = ""
			if reused[i].State == Reused {
				reused[i].RanFor = selected[i].RanFor
			}
		}
		rep.Steps = append(rep.Steps, reused...)
		v.endedUnrun(reused)
		ro := v.opts.RunOptions
		// the check step passed, as every step before the tiers did, and the
		// reused tiers passed in the run they come from (S-0610)
		passed := []string{PassedStepCheck}
		for i, step := range reused {
			if step.State == Reused {
				passed = append(passed, PassedWord(selected[i]))
			}
		}
		ro.Env = append(append([]string{}, ro.Env...), StoryEnv+"="+v.story, PassedEnv+"="+strings.Join(passed, " "))
		ro.TierStarted, ro.TierEnded = v.opts.Started, nil
		if v.opts.Ended != nil {
			ro.TierEnded = func(tier string, state State) { v.opts.Ended(Step{Name: tier, Tier: true, State: state}) }
		}
		var res Result
		res, runErr = Run(ctx, v.worktree, selected[len(reused):], ro)
		for i, tr := range res.Tiers {
			step := tierStep(tr)
			if i == 0 {
				step.Note = sentBack
			}
			if tr.State == Passed || tr.State == Failed {
				step.RanFor = selected[len(reused)+i].RanFor
			}
			rep.Steps = append(rep.Steps, step)
			v.endedUnrun([]Step{step})
			if tr.State == Failed && rep.Passed {
				rep.Passed, rep.StoppedAt = false, tr.Name
			}
		}
	} else {
		for _, s := range selected {
			step := unrunStep(s)
			rep.Steps = append(rep.Steps, step)
			v.endedUnrun([]Step{step})
		}
	}
	for _, s := range queued {
		step := queuedStep(s)
		rep.Steps = append(rep.Steps, step)
		v.ended(step)
	}
	rep.SuiteRuns = v.suiteRecords(rep, listed)
	rep.DurationMS, rep.Duration = took(v.opts.Now().Sub(start))
	return rep, runErr
}

// suiteRecords are what the run rep recorded of each suite of the tiers it
// listed, run and left to the merge queue alike, as verifyrecord.Group groups
// its tier steps, each with the fingerprint of the project's tiers and the
// suite's listed ones (ADR-0312 §7).
func (v *storyRun) suiteRecords(rep Report, listed []Selected) []verifyrecord.Suite {
	suites := verifyrecord.Group(rep.record(), suiteOf(v.opts.Tiers))
	for i := range suites {
		suites[i].Fingerprint = fingerprint(v.opts.Tiers, ofSuite(listed, suites[i].Name))
	}
	return suites
}

// record is what verifyrecord reads of the report, as it reads it stored.
func (r Report) record() verifyrecord.Record {
	rec := verifyrecord.Record{Story: r.Story, Commit: r.Commit, Base: r.Base, BaseCommit: r.BaseCommit, MergeBase: r.MergeBase,
		Patch: r.Patch, RanAt: r.RanAt, Passed: r.Passed, Clean: r.Clean, Steps: make([]verifyrecord.Step, 0, len(r.Steps)), SuiteRuns: r.SuiteRuns}
	for _, s := range r.Steps {
		rec.Steps = append(rec.Steps, verifyrecord.Step{Name: s.Name, Tier: s.Tier, State: verifyrecord.State(s.State), CoveredBy: s.CoveredBy, ReusedFrom: s.ReusedFrom})
	}
	return rec
}

// PassedSuites says whether the report holds, for every suite of names, an
// entry at commit that passed: what a verdict answered from the record
// covers (ADR-0312 §7). No names is true.
func (r Report) PassedSuites(commit string, names []string) bool {
	return verifyrecord.AllPassed(r.SuiteRuns, commit, names)
}

// suiteOf gives each of tiers' suite by its name, SuiteOrName; "" for a
// name none of them has.
func suiteOf(tiers []Tier) func(string) string {
	return func(name string) string {
		for _, t := range tiers {
			if t.Name == name {
				return t.SuiteOrName()
			}
		}
		return ""
	}
}

// ofSuite are the tiers of selection in the suite, in its order.
func ofSuite(selection []Selected, suite string) []Selected {
	var out []Selected
	for _, s := range selection {
		if s.Tier.SuiteOrName() == suite {
			out = append(out, s)
		}
	}
	return out
}

// suiteNames are the suites of selection's tiers, each once, in the order of
// the first tier of each.
func suiteNames(selection []Selected) []string {
	var out []string
	for _, s := range selection {
		if name := s.Tier.SuiteOrName(); !slices.Contains(out, name) {
			out = append(out, name)
		}
	}
	return out
}

// tierNames are the names of selection's tiers, in its order.
func tierNames(selection []Selected) []string {
	names := make([]string, 0, len(selection))
	for _, s := range selection {
		names = append(names, s.Tier.Name)
	}
	return names
}

// selection is the tiers of the run's Suites that the branch's paths select:
// those the run runs, and, with LeaveToQueue, those it leaves to the merge
// queue, and the fingerprint of both.
func (v *storyRun) selection() (selected, queued []Selected, fp string) {
	if v.opts.LeaveToQueue {
		selected, queued = SelectStoryLeavingQueue(v.opts.FS, v.opts.Tiers, v.paths, v.claim(), v.opts.Suites)
	} else {
		selected = SelectStorySuites(v.opts.FS, v.opts.Tiers, v.paths, v.opts.Suites)
	}
	return selected, queued, fingerprint(v.opts.Tiers, inListOrder(v.opts.Tiers, selected, queued))
}

// claim is the story's claim, as Holds.Claim reads it from the project's
// items (ADR-0096): its touches, each folder narrowed to its tasks' touches
// inside it, and its open tasks' touches. It is none when the items cannot
// be read or do not hold the story, so that the changed paths alone choose
// the queue tiers a run does not leave to the merge queue (ADR-0233).
func (v *storyRun) claim() []string {
	items, err := v.opts.Project.List(false)
	if err != nil {
		return nil
	}
	for _, it := range items {
		if it.Type == workitem.Story && it.ID == v.story {
			return workitem.NewHolds(items, v.opts.Project.Manifest.Projects).Claim(it)
		}
	}
	return nil
}

// queueOffText is Report's QueueOff for a run that selected these tiers and
// left none to the merge queue, which cannot batch because of why: the tiers
// marked Queue among them, named in list order, and why; "" when why is ""
// or none is marked Queue (S-0466).
func queueOffText(selected []Selected, why string) string {
	if why == "" {
		return ""
	}
	var names []string
	for _, s := range selected {
		if s.Tier.Queue {
			names = append(names, s.Tier.Name)
		}
	}
	if len(names) == 0 {
		return ""
	}
	return tierList(names) + " ran here, not left to the merge queue: " + why
}

// tierList names tiers as a sentence does: "tier a", or "tiers a, b and c".
func tierList(names []string) string {
	if len(names) == 1 {
		return "tier " + names[0]
	}
	return "tiers " + strings.Join(names[:len(names)-1], ", ") + " and " + names[len(names)-1]
}

// inListOrder is the tiers a run runs and those it leaves to the merge queue
// as one selection, in the order tiers lists them: what SelectStory chooses,
// but for which are marked covered, which the fingerprint does not hash
// (I-0203). A tier marked Queue that the run runs (ADR-0233) takes its place
// in the list as one left to the queue would.
func inListOrder(tiers []Tier, run, queued []Selected) []Selected {
	if len(queued) == 0 {
		return run
	}
	out := make([]Selected, 0, len(run)+len(queued))
	for _, t := range tiers {
		switch {
		case len(run) > 0 && run[0].Tier.Queue == t.Queue && run[0].Tier.Name == t.Name:
			out, run = append(out, run[0]), run[1:]
		case len(queued) > 0 && t.Queue && queued[0].Tier.Name == t.Name:
			out, queued = append(out, queued[0]), queued[1:]
		}
	}
	// both come from tiers in its order, so nothing is left; if anything
	// were, it is still hashed
	return append(append(out, run...), queued...)
}

// QueueTiersNotRun answers, for a story's last report (ok false when it has
// none) and the head its branch has, why the tiers marked Queue in tiers
// cannot be taken as run at that head, as a sentence fragment a caller can
// embed, "" when they ran; and, for a report at head, the names of those the
// report selected but did not run, in its order: each step left Queued to
// the merge queue, and each step of a Queue tier that did not pass, was not
// reused, and was not skipped for a tier that covers it. A Queue tier the
// report did not select is not named. A run with every selected tier, or
// with the queue off, that passed at head answers nothing (ADR-0174, I-0203).
// A report not at head, or none, whose Running LastReportAt set names the
// close-out job that runs, where, and at which of its steps (I-0384).
func QueueTiersNotRun(rep Report, ok bool, head string, tiers []Tier) (names []string, why string) {
	if !ok {
		if rep.Running != nil {
			return nil, "no flai verify has run on it, and " + jobText(*rep.Running, head, "at its head "+shortHash(head))
		}
		return nil, "no flai verify has run on it"
	}
	if rep.Commit != head {
		why := fmt.Sprintf("its last flai verify ran at %s, not at its head %s", shortHash(rep.Commit), shortHash(head))
		if p := rep.Running; p != nil {
			if p.Head == "" || p.Head == head {
				why += ", where " + jobText(*p, head, "")
			} else {
				why += ", and " + jobText(*p, head, "")
			}
		}
		return nil, why
	}
	queue := map[string]bool{}
	for _, t := range tiers {
		if t.Queue {
			queue[t.Name] = true
		}
	}
	for _, s := range rep.Steps {
		if !s.Tier || slices.Contains(names, s.Name) {
			continue
		}
		ran := s.State == Passed || s.State == Reused || s.State == Skipped
		if s.State == Queued || (queue[s.Name] && !ran) {
			names = append(names, s.Name)
		}
	}
	switch {
	case !rep.Passed:
		return names, "its last flai verify stopped at " + rep.StoppedAt
	case len(names) > 0:
		return names, "its last flai verify did not run " + tierList(names) + ", marked queue"
	}
	return nil, ""
}

// jobText says that close-out job p runs, as "close-out job <id> runs, at its
// step <step>": at its own head when that is not head, else with atHead, ""
// or where head is; the step is left out when p has none.
func jobText(p closeout.Progress, head, atHead string) string {
	s := "close-out job " + p.ID + " runs"
	switch {
	case p.Head != "" && p.Head != head:
		s += " at " + shortHash(p.Head)
	case atHead != "":
		s += " " + atHead
	}
	if p.Step != "" {
		s += ", at its step " + p.Step
	}
	return s
}

// MainPathsSince are the paths the main branch, at the commit or branch main,
// brought in since a story's report rep ran: what changed from the report's
// merge base to main, read in the checkout at dir, but for the paths under
// wip, the wip folder with a trailing slash, that the report's branch does not
// change, which no tier's result reads (ADR-0135, ADR-0291). A report that
// recorded no merge base is taken from its base commit, else from where its
// commit and main meet now. A report that recorded none of the three, and a
// diff git cannot make, are errors.
func MainPathsSince(r execx.Runner, dir string, rep Report, main, wip string) ([]string, error) {
	from := rep.MergeBase
	if from == "" {
		from = rep.BaseCommit
	}
	if from == "" && rep.Commit != "" {
		out, err := r.Run(dir, "git", "merge-base", rep.Commit, main)
		if err != nil {
			return nil, fmt.Errorf("find where %s's last flai verify, at %s, left %s, since its report records no merge base: %w", rep.Story, shortHash(rep.Commit), main, err)
		}
		from = strings.TrimSpace(out)
	}
	if from == "" {
		return nil, fmt.Errorf("tell what %s brought in since %s's last flai verify: its report records no merge base, base commit, or commit", main, rep.Story)
	}
	paths, err := BetweenPaths(r, dir, from, main)
	if err != nil {
		return nil, err
	}
	var out []string
	for _, p := range paths {
		if !strings.HasPrefix(p, wip) || slices.Contains(rep.Paths, p) {
			out = append(out, p)
		}
	}
	return out, nil
}

// StaleTiers are the tiers a story's report rep ran, passed or reused, that
// mainPaths, the paths the main branch brought in since it ran
// (MainPathsSince), send back, in the report's order: each tier the paths
// select among tiers as a story's own paths select its tiers (SelectStory,
// {importers} included), and each tier the report ran in place of one they
// select, as the tier that covered it. A tier the report skipped, left to the
// merge queue, or did not reach is not named: it ran nothing to be stale
// (ADR-0156, ADR-0291).
func StaleTiers(fsys fs.FS, tiers []Tier, rep Report, mainPaths []string) []string {
	if len(mainPaths) == 0 {
		return nil
	}
	back := map[string]bool{}
	for _, s := range SelectStory(fsys, tiers, mainPaths) {
		back[s.Tier.Name] = true
	}
	for _, s := range rep.Steps {
		if s.Tier && s.State == Skipped && s.CoveredBy != "" && back[s.Name] {
			back[s.CoveredBy] = true
		}
	}
	var names []string
	for _, s := range rep.Steps {
		if s.Tier && (s.State == Passed || s.State == Reused) && back[s.Name] && !slices.Contains(names, s.Name) {
			names = append(names, s.Name)
		}
	}
	return names
}

// queuedStep is a selected tier left to the merge queue's batch as a step:
// queued, whatever the steps before it did, since the batch runs it
// (ADR-0174).
func queuedStep(s Selected) Step {
	return Step{Name: s.Tier.Name, Tier: true, State: Queued, Command: s.Argv, Dir: s.Tier.Dir}
}

// unrunStep is a selected tier that does not run as a step: skipped, naming
// the tier that covers it, or else not reached.
func unrunStep(s Selected) Step {
	step := Step{Name: s.Tier.Name, Tier: true, State: NotReached, Command: s.Argv, Dir: s.Tier.Dir}
	if s.CoveredBy != "" {
		step.State, step.CoveredBy = Skipped, s.CoveredBy
	}
	return step
}

// TierOptions are what Tiers needs to verify a checkout with the tiers
// its changes select.
type TierOptions struct {
	// Dir is the checkout whose HEAD is verified.
	Dir string
	// Base is the branch or commit the checkout's changes are counted
	// against.
	Base string
	// Tiers are the project's tiers, cheapest first, as the checkout's
	// manifest declares them.
	Tiers []Tier
	// Git runs git; nil is execx.System.
	Git execx.Runner
	// Started is called with a tier's name as it starts, and Ended with its
	// step as it ends; nil calls nothing. Neither is called for a tier that
	// is skipped or not reached.
	Started func(tier string)
	Ended   func(Step)
	// Suites are the suites whose tiers run, by each tier's SuiteOrName, as
	// the batch runs those its workflow's verify_at.batch lists (ADR-0312):
	// nil runs every tier the paths select, as the word all does; a list
	// that is not nil runs only those of its suites; an empty one runs none.
	Suites []string
	// RunOptions are the tiers' run's; FS is the checkout, read for the
	// tiers' packages, and Now times every tier and stamps the report.
	RunOptions
}

// Tiers runs, in the checkout at Dir, the tiers flai verify selects for
// the paths its HEAD changed against Base (SelectStory), or of them those of
// Suites when it is not nil (SelectStorySuites), cheapest first, and stops at
// the first that fails, as Verify runs them, and answers them as a report
// with no story and no step before the tiers, naming the Suites it was given.
// It verifies a checkout that is no story's, as the merge queue's batch: it
// reuses no earlier run, sets no CLOSE_OUT_STORY, and stores no report. The
// caller may name the stories a batch verifies to its tiers by setting
// MergeQueueStoriesEnv in RunOptions.Env. It returns an error, with no
// report, when the checkout's HEAD or changes cannot be read, and, with the
// report so far, when ctx ends.
func Tiers(ctx context.Context, opts TierOptions) (Report, error) {
	if opts.Git == nil {
		opts.Git = execx.System{}
	}
	if opts.Now == nil {
		opts.Now = time.Now
	}
	if opts.FS == nil {
		opts.FS = os.DirFS(opts.Dir)
	}
	start := opts.Now()
	commit, err := opts.Git.Run(opts.Dir, "git", "rev-parse", "HEAD")
	if err != nil {
		return Report{}, fmt.Errorf("verify the checkout %s: read the commit it has checked out: %w; check that it is a git worktree", opts.Dir, err)
	}
	paths, err := ChangedPaths(opts.Git, opts.Dir, opts.Base)
	if err != nil {
		return Report{}, fmt.Errorf("verify the checkout %s: %w", opts.Dir, err)
	}
	rep := Report{Commit: strings.TrimSpace(commit), Base: opts.Base, RanAt: start.UTC(), Passed: true, Paths: paths, Steps: []Step{}}
	if base, err := opts.Git.Run(opts.Dir, "git", "rev-parse", "--verify", "--quiet", opts.Base+"^{commit}"); err == nil {
		rep.BaseCommit = strings.TrimSpace(base)
	}
	if opts.Suites != nil {
		rep.Suites = slices.Clone(opts.Suites)
	}
	selected := SelectStorySuites(opts.FS, opts.Tiers, paths, opts.Suites)
	rep.Fingerprint = fingerprint(opts.Tiers, selected)
	var runErr error
	// what passed earlier in the batch, which each tier is told as one Run
	// of every tier tells the tiers after them (S-0610)
	var passed []string
	for _, s := range selected {
		if s.CoveredBy != "" || !rep.Passed || ctx.Err() != nil {
			rep.Steps = append(rep.Steps, unrunStep(s))
			continue
		}
		if opts.Started != nil {
			opts.Started(s.Tier.Name)
		}
		ro := opts.RunOptions
		ro.Env = append(append([]string{}, ro.Env...), PassedEnv+"="+strings.Join(passed, " "))
		res, err := Run(ctx, opts.Dir, []Selected{s}, ro)
		step := tierStep(res.Tiers[0])
		rep.Steps = append(rep.Steps, step)
		if step.State == Failed {
			rep.Passed, rep.StoppedAt = false, step.Name
		}
		if step.State == Passed {
			passed = append(passed, PassedWord(s))
		}
		if opts.Ended != nil {
			opts.Ended(step)
		}
		if err != nil {
			runErr = err
		}
	}
	if runErr == nil {
		runErr = ctx.Err()
	}
	rep.DurationMS, rep.Duration = took(opts.Now().Sub(start))
	return rep, runErr
}

// rebase finds a rebase left unfinished in the worktree.
func (v *storyRun) rebase() []Finding {
	if !storygit.RebaseInProgress(v.opts.Git, v.worktree) {
		return nil
	}
	return []Finding{{Name: StepRebase, Message: fmt.Sprintf("a rebase is in progress; resolve each conflicting path, git add it, and finish it with flai stream sync %s, which continues it, or undo it with git rebase --abort, then verify again", v.story)}}
}

// sync finds that the story's branch lacks commits of the main branch,
// unless what they changed since the branch left it is only paths under the
// manifest's wip folder that the branch does not change; then it passes
// over them, and its note names them (ADR-0135).
func (v *storyRun) sync() []Finding {
	paths, commits, err := BaseChanges(v.opts.Git, v.worktree, v.base)
	if err != nil {
		return []Finding{{Name: StepSync, Message: fmt.Sprintf("cannot tell whether the branch contains %s: %v; check that %s is a branch of the repository", v.base, err, v.base)}}
	}
	if len(commits) == 0 {
		return nil
	}
	wip, err := v.wipPrefix()
	if err != nil {
		return []Finding{{Name: StepSync, Message: err.Error()}}
	}
	var outside, shared []string
	for _, p := range paths {
		switch {
		case !strings.HasPrefix(p, wip):
			outside = append(outside, p)
		case slices.Contains(v.paths, p):
			shared = append(shared, p)
		}
	}
	behind := fmt.Sprintf("the branch does not contain %s, %s behind it", v.base, plural(len(commits), "commit"))
	fix := fmt.Sprintf("run flai stream sync %s, resolve what it reports, and verify again", v.story)
	change := "change"
	if len(commits) == 1 {
		change = "changes"
	}
	switch {
	case len(outside) > 0:
		return []Finding{{Name: StepSync, Message: fmt.Sprintf("%s, which %s paths outside %s: %s; %s", behind, change, wip, few(outside, 3), fix)}}
	case len(shared) > 0:
		return []Finding{{Name: StepSync, Message: fmt.Sprintf("%s, which %s %s paths the branch changes too: %s; %s", behind, change, wip, few(shared, 3), fix)}}
	}
	v.note = fmt.Sprintf("passed over %s of %s that %s only %s paths the branch does not change: %s", plural(len(commits), "commit"), v.base, change, wip, few(commits, 10))
	return nil
}

// wipPrefix is the manifest's wip folder relative to the project's root,
// with a trailing slash, as git names the paths under it.
func (v *storyRun) wipPrefix() (string, error) {
	wip, err := filepath.Rel(mainRoot(v.opts.Project), v.opts.Project.WipDir())
	if err != nil {
		return "", fmt.Errorf("cannot tell which paths are under the wip folder %s: %w; check the layout in the project's manifest", v.opts.Project.WipDir(), err)
	}
	return filepath.ToSlash(wip) + "/", nil
}

// plural is n and the noun, plural but for one.
func plural(n int, noun string) string {
	if n == 1 {
		return "1 " + noun
	}
	return strconv.Itoa(n) + " " + noun + "s"
}

// few lists the first limit of items and says how many more there are.
func few(items []string, limit int) string {
	if len(items) <= limit {
		return strings.Join(items, ", ")
	}
	return fmt.Sprintf("%s and %d more", strings.Join(items[:limit], ", "), len(items)-limit)
}

// narrative finds the narrative's sections that are not written, in the
// main checkout, where the narrative lives.
func (v *storyRun) narrative() []Finding {
	path := v.opts.Project.NarrativePath(v.story)
	rel := projectPath(v.opts.Project, v.worktree, path)
	data, err := os.ReadFile(path)
	if err != nil {
		return []Finding{{Name: StepNarrative, Path: rel, Message: fmt.Sprintf("no narrative for %s: %v; open it with flai stream open %s", v.story, err, v.story)}}
	}
	found := Unwritten(string(data), narrativeSections)
	for i := range found {
		found[i].Path = rel
	}
	return found
}

// check runs flai check --strict scoped to the story, in process, on the
// worktree: a finding inside the story at error level, or at warning level
// but those --strict passes over, fails it, and one outside is put in rep's
// notes (S-0249).
func (v *storyRun) check(rep *Report) []Finding {
	repo, err := workitem.Open(v.worktree)
	if err != nil {
		return []Finding{{Name: StepCheck, Message: fmt.Sprintf("open the project at the worktree %s: %v", v.worktree, err)}}
	}
	repo.Git, repo.TemplateDir = v.opts.Git, v.opts.Project.TemplateDir
	res, err := check.Run(repo, v.opts.Now())
	if err != nil {
		return []Finding{{Name: StepCheck, Message: fmt.Sprintf("flai check could not run: %v", err)}}
	}
	// The clone may need a newer git than the one installed (ADR-0022).
	if installed, err := gitver.Installed(v.opts.Git); err == nil {
		check.GitCompat(res, repo, installed)
	}
	// The project's agent image layer, when it names one (S-0510).
	check.AgentImage(res, repo, check.DockerImage(v.opts.Git))
	changed, err := storygit.StoryChanges(v.opts.Git, repo, v.story)
	if err != nil {
		return []Finding{{Name: StepCheck, Message: fmt.Sprintf("scope the check to %s: %v", v.story, err)}}
	}
	if err := check.ScopeToStory(res, repo, v.story, changed); err != nil {
		return []Finding{{Name: StepCheck, Message: err.Error()}}
	}
	var inside []Finding
	for _, f := range res.Findings {
		p := ""
		if f.Path != "" {
			p = projectPath(repo, v.worktree, f.Path)
		}
		if f.Outside {
			rep.Notes = append(rep.Notes, Note{Rule: f.Rule, Level: f.Level, Path: p, Line: f.Line, Message: f.Message})
			continue
		}
		inside = append(inside, Finding{Name: f.Rule, Path: p, Line: f.Line, Message: f.Level + ": " + f.Message})
	}
	if res.OK(true) {
		return nil
	}
	return inside
}

// tierStep is a tier's result as a step of Verify.
func tierStep(tr TierResult) Step {
	return Step{
		Name: tr.Name, Tier: true, State: tr.State, DurationMS: tr.DurationMS, Duration: tr.Duration,
		Findings: tr.Findings, Omitted: tr.Omitted, Command: tr.Command, Dir: tr.Dir, ExitCode: tr.ExitCode,
		CoveredBy: tr.CoveredBy, Log: tr.Log,
	}
}

// reusable is the steps of the selected tiers that the story's last report
// lets this run reuse, in order, always the first of selected: each tier is
// paired with its entry in the report's record of its suite, and those
// before the first it shows failed, not reached, or left Queued to the merge
// queue are each Reused with the time of the run that ran it, but for a tier
// another selected tier covers, which stays Skipped. Pairing by name lets a
// run that leaves tiers to the queue reuse what one that ran them passed, and
// the other way round, since the fingerprint is the same (ADR-0174, I-0203).
// listed is the run's tiers, run and left to the queue, as the fingerprint
// hashes them. It is none with Fresh or with no last report; every tier, when
// it passed them all. It stops at the first tier whose suite the report
// holds no entry of, or an entry with another fingerprint or other tiers in
// another order. At the head and base commit of the suite's entry that is all
// (S-0341). At another, it reuses the suite only when the branch's own
// change, its patch id, is the same, and stops at the first tier that the
// paths the main branch brought in between the two merge bases select, but
// for wip paths the branch does not change (ADR-0135), as the branch's own
// select its tiers, or that covers one they select; then it answers, as the
// note for that tier, why it runs again, and the first tier it reuses across
// the main branch's new commits says so (ADR-0156, ADR-0312 §7).
func (v *storyRun) reusable(selected, listed []Selected) (steps []Step, sentBackNote string) {
	if v.opts.Fresh || v.baseCommit == "" {
		return nil, ""
	}
	last, ok, err := LastReport(v.opts.Project, v.story)
	if err != nil || !ok {
		return nil, ""
	}
	// the last run's tier steps by name, whose command, duration, and log a
	// tier it ran or reused keeps when reused again
	ran := map[string]Step{}
	for _, s := range last.Steps {
		if _, seen := ran[s.Name]; s.Tier && !seen {
			ran[s.Name] = s
		}
	}
	gates := map[string]suiteGate{}
	var out []Step
	noted := false
	for i, sel := range selected {
		name := sel.Tier.SuiteOrName()
		g, seen := gates[name]
		if !seen {
			g = v.gate(last.SuiteRuns, name, selected, listed)
			gates[name] = g
		}
		if !g.ok {
			break
		}
		t, ok := g.suite.Tier(sel.Tier.Name)
		if !ok {
			break
		}
		if g.sentBack[t.Name] {
			sentBackNote = v.sentBackNote(selected[i:], g.sentBack, g.mainPaths)
			break
		}
		s, lastRan := reusedStep(ran, sel)
		// A tier another selected tier covers ran nothing, so it neither
		// stops the reuse nor is reused: it stays Skipped (S-0342).
		if t.State == verifyrecord.Skipped && sel.CoveredBy != "" {
			s.State, s.CoveredBy, s.ReusedFrom = Skipped, sel.CoveredBy, nil
			out = append(out, s)
			continue
		}
		if t.State != verifyrecord.Passed && t.State != verifyrecord.Reused {
			break
		}
		from := g.suite.RanAt.UTC()
		if lastRan {
			from = last.RanAt.UTC()
		}
		if t.State == verifyrecord.Reused && t.ReusedFrom != nil {
			from = t.ReusedFrom.UTC()
		}
		s.State, s.ReusedFrom = Reused, &from
		if g.moved && !noted {
			s.Note, noted = v.reusedNote(g.suite, g.mainPaths), true
		}
		out = append(out, s)
	}
	return out, sentBackNote
}

// suiteGate is what a suite's entry in the last report lets a run reuse: ok
// when its tiers may be reused at all, moved when the entry is at another
// head or base commit, the paths the main branch brought in since its merge
// base, and the tiers they send back.
type suiteGate struct {
	ok, moved bool
	suite     verifyrecord.Suite
	mainPaths []string
	sentBack  map[string]bool
}

// gate is the suiteGate of the suite name among the last report's suites:
// not ok when none is recorded, its fingerprint or its tiers and their order
// differ from this run's listed ones of it, or, at another head or base
// commit, the branch's own change may differ or what main brought in since
// cannot be told (ADR-0156).
func (v *storyRun) gate(suites []verifyrecord.Suite, name string, selected, listed []Selected) suiteGate {
	i := slices.IndexFunc(suites, func(s verifyrecord.Suite) bool { return s.Name == name })
	if i < 0 {
		return suiteGate{}
	}
	rec, of := suites[i], ofSuite(listed, name)
	if rec.Fingerprint == "" || rec.Fingerprint != fingerprint(v.opts.Tiers, of) || !slices.Equal(rec.TierNames(), tierNames(of)) {
		return suiteGate{}
	}
	g := suiteGate{ok: true, suite: rec, moved: rec.Commit != v.commit || rec.BaseCommit != v.baseCommit}
	if !g.moved {
		return g
	}
	if v.patch == "" || rec.Patch != v.patch || v.mergeBase == "" || rec.MergeBase == "" {
		return suiteGate{}
	}
	if rec.MergeBase != v.mergeBase {
		// what main brought in cannot be told, so nothing is reused
		paths, err := BetweenPaths(v.opts.Git, v.worktree, rec.MergeBase, v.mergeBase)
		if err != nil {
			return suiteGate{}
		}
		if g.mainPaths, err = v.outsideWip(paths); err != nil {
			return suiteGate{}
		}
	}
	g.sentBack = v.sentBack(selected, g.mainPaths)
	return g
}

// reusedStep is the step a reused tier begins from: the last run's of its
// name, with no note, and lastRan true, when that run ran, reused, or skipped
// it; else one with the selected tier's command, since the entry of its
// suite is an earlier run's.
func reusedStep(ran map[string]Step, sel Selected) (step Step, lastRan bool) {
	if s, ok := ran[sel.Tier.Name]; ok && (s.State == Passed || s.State == Reused || s.State == Skipped) {
		s.Note = ""
		return s, true
	}
	return Step{Name: sel.Tier.Name, Tier: true, Command: sel.Argv, Dir: sel.Tier.Dir}, false
}

// outsideWip leaves out of paths the main branch brought in those under the
// wip folder that the branch does not change: no tier's result reads them
// (ADR-0135), so they send no tier back, as they hold no sync back.
func (v *storyRun) outsideWip(paths []string) ([]string, error) {
	wip, err := v.wipPrefix()
	if err != nil {
		return nil, err
	}
	var out []string
	for _, p := range paths {
		if !strings.HasPrefix(p, wip) || slices.Contains(v.paths, p) {
			out = append(out, p)
		}
	}
	return out, nil
}

// sentBack are the names of the selected tiers that the paths the main
// branch brought in reach: those the paths select as the branch's own select
// its tiers (SelectStory, {importers} included), and the tier that covers
// each of them in this run, which runs in its place (ADR-0156).
func (v *storyRun) sentBack(selected []Selected, mainPaths []string) map[string]bool {
	out := map[string]bool{}
	if len(mainPaths) == 0 {
		return out
	}
	for _, s := range SelectStory(v.opts.FS, v.opts.Tiers, mainPaths) {
		out[s.Tier.Name] = true
	}
	for _, s := range selected {
		if out[s.Tier.Name] && s.CoveredBy != "" {
			out[s.CoveredBy] = true
		}
	}
	return out
}

// sentBackNote says that the tiers from rest on run again, and which of
// them the main branch's new paths select.
func (v *storyRun) sentBackNote(rest []Selected, sentBack map[string]bool, mainPaths []string) string {
	var names []string
	for _, s := range rest {
		if sentBack[s.Tier.Name] {
			names = append(names, s.Tier.Name)
		}
	}
	return fmt.Sprintf("run again from here: the %s %s brought in since the last run %s %s: %s",
		plural(len(mainPaths), "path"), v.base, selectVerb(len(mainPaths)), strings.Join(names, ", "), few(mainPaths, 3))
}

// reusedNote says that a tier was reused across the main branch's new
// commits since its suite's entry last, and why it could be.
func (v *storyRun) reusedNote(last verifyrecord.Suite, mainPaths []string) string {
	if last.MergeBase == v.mergeBase {
		return fmt.Sprintf("reused though %s moved on: the branch's own change is the same as at %s, and it was not rebased", v.base, shortHash(last.Commit))
	}
	if len(mainPaths) == 0 {
		return fmt.Sprintf("reused across the rebase onto %s: the branch's own change is the same as at %s, and %s brought in only wip paths the branch does not change",
			v.base, shortHash(last.Commit), v.base)
	}
	return fmt.Sprintf("reused across the rebase onto %s: the branch's own change is the same as at %s, and the %s %s brought in %s none of the tiers reused",
		v.base, shortHash(last.Commit), plural(len(mainPaths), "path"), v.base, selectVerb(len(mainPaths)))
}

// selectVerb is "select" agreeing with n paths.
func selectVerb(n int) string {
	if n == 1 {
		return "selects"
	}
	return "select"
}

// shortHash is a commit's first seven characters, as git abbreviates it.
func shortHash(commit string) string {
	if len(commit) > 7 {
		return commit[:7]
	}
	return commit
}

// fingerprint hashes the project's tiers as declared and the selection's
// names, directories, and commands: a change to any changes what a story's
// tiers would run. The selection is every tier selected, those left to the
// merge queue as well as those run, so whether a run leaves tiers to the
// queue does not change it (ADR-0174, I-0203).
func fingerprint(tiers []Tier, selection []Selected) string {
	type sel struct {
		Name string   `json:"name"`
		Dir  string   `json:"dir"`
		Argv []string `json:"argv"`
	}
	in := struct {
		Tiers    []Tier `json:"tiers"`
		Selected []sel  `json:"selected"`
	}{Tiers: tiers, Selected: []sel{}}
	for _, s := range selection {
		in.Selected = append(in.Selected, sel{Name: s.Tier.Name, Dir: s.Tier.Dir, Argv: s.Argv})
	}
	data, err := json.Marshal(in)
	if err != nil {
		// nothing to compare, so nothing is reused
		return ""
	}
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}

// took is d as a report gives it: whole milliseconds, and as a person
// reads it.
func took(d time.Duration) (int64, string) {
	return d.Round(time.Millisecond).Milliseconds(), human(d)
}

// listMarker is what a line of a section may begin with and still be
// empty: blanks, and a bullet or a number such as "1.", as close-out.sh
// reads it.
var listMarker = regexp.MustCompile(`^\s*(?:[-*]|[0-9]+\.)?\s*`)

// Unwritten finds each of the sections, by heading text, that a narrative
// lacks or leaves empty: under its "## " heading, up to the next, no line
// holds more than blanks and a bare list marker. A finding's line is the
// heading's, 0 when it is missing; its path is left for the caller.
func Unwritten(narrative string, sections []string) []Finding {
	lines := strings.Split(strings.ReplaceAll(narrative, "\r\n", "\n"), "\n")
	var out []Finding
	for _, sec := range sections {
		heading := "## " + sec
		at, written := -1, false
		for i, line := range lines {
			if at < 0 {
				if line == heading {
					at = i
				}
				continue
			}
			if strings.HasPrefix(line, "## ") {
				break
			}
			if listMarker.ReplaceAllString(line, "") != "" {
				written = true
				break
			}
		}
		switch {
		case at < 0:
			out = append(out, Finding{Name: StepNarrative, Message: fmt.Sprintf("the narrative has no %s section; add it and write it, then verify again", heading)})
		case !written:
			out = append(out, Finding{Name: StepNarrative, Line: at + 1, Message: fmt.Sprintf("%s is empty; write it, then verify again", heading)})
		}
	}
	return out
}

// projectPath is p relative to the worktree, or, for what lives in the main
// checkout such as wip/, to it; p itself when it is below neither.
func projectPath(repo *workitem.Repo, worktree, p string) string {
	if !filepath.IsAbs(p) {
		p = filepath.Join(worktree, p)
	}
	for _, root := range []string{worktree, mainRoot(repo)} {
		if rel, err := filepath.Rel(root, p); err == nil && rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
			return filepath.ToSlash(rel)
		}
	}
	return filepath.ToSlash(p)
}

// mainRoot is the project's main checkout.
func mainRoot(repo *workitem.Repo) string {
	if repo.MainRoot != "" {
		return repo.MainRoot
	}
	return repo.Root
}

// ReportPath is where a story's last report is stored: in the project's
// cache, .flai-cache/verify/<story>.json.
func ReportPath(project *workitem.Repo, story string) string {
	return verifyrecord.ReportPath(project.CacheDir(), story)
}

// SaveReport stores rep as its story's last report, replacing the one
// before it in one step, so that a reader sees one whole report or the
// other. Its SuiteRuns keep those of the report before it that rep did not run,
// as verifyrecord.Merge merges them, so that the last report holds the last
// run of each suite (ADR-0312 §7); a report before it that does not read
// keeps none.
func SaveReport(project *workitem.Repo, rep Report) error {
	// what LastReportAt reads against the head is not the run's to store
	rep.Moved, rep.Running = false, nil
	if before, ok, err := LastReport(project, rep.Story); err == nil && ok {
		rep.SuiteRuns, _ = verifyrecord.Merge(before.SuiteRuns, rep.SuiteRuns)
	}
	path := ReportPath(project, rep.Story)
	data, err := json.MarshalIndent(rep, "", "  ")
	if err != nil {
		return fmt.Errorf("store the verify report of %s: %w", rep.Story, err)
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("store the verify report of %s: make %s: %w; check that the project's .flai-cache is writable", rep.Story, filepath.Dir(path), err)
	}
	if err := atomicfile.WriteFile(path, append(data, '\n'), 0o644); err != nil {
		return fmt.Errorf("store the verify report of %s at %s: %w; check that the project's .flai-cache is writable", rep.Story, path, err)
	}
	return nil
}

// LastReport is the story's last report and true, or false with no error
// when the story has none. The story may be given as typed, S-12 for
// S-0012, when the project has it.
func LastReport(project *workitem.Repo, story string) (Report, bool, error) {
	if it, err := project.Get(story); err == nil {
		story = it.ID
	}
	path := ReportPath(project, story)
	data, err := os.ReadFile(path)
	if errors.Is(err, fs.ErrNotExist) {
		return Report{}, false, nil
	}
	if err != nil {
		return Report{}, false, fmt.Errorf("read the last verify report of %s at %s: %w", story, path, err)
	}
	var rep Report
	if err := json.Unmarshal(data, &rep); err != nil {
		return Report{}, false, fmt.Errorf("read the last verify report of %s at %s: %w; delete it and verify again", story, path, err)
	}
	rep.Moved, rep.Running = false, nil
	return rep, true, nil
}

// LastReportAt is LastReport read against head, the head of the story's
// branch: Moved says the report is not at head, and Running is the story's
// running close-out job when it runs at another commit than the report's, a
// job that has not written its head taken as running at head (I-0384). With
// no report, ok is false and the report holds only Story and Running. An
// empty head, from a caller that has none, leaves Moved false.
func LastReportAt(project *workitem.Repo, story, head string) (Report, bool, error) {
	if it, err := project.Get(story); err == nil {
		story = it.ID
	}
	rep, ok, err := LastReport(project, story)
	if err != nil {
		return Report{}, false, err
	}
	p, running, err := closeout.Running(project.CacheDir(), story)
	if err != nil {
		return Report{}, false, err
	}
	if !ok {
		rep.Story = story
	}
	rep.Moved = ok && head != "" && rep.Commit != head
	at := p.Head
	if at == "" {
		at = head
	}
	if running && (!ok || at != rep.Commit) {
		rep.Running = &p
	}
	return rep, ok, nil
}
