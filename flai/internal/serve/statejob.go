package serve

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/closeout"
	"github.com/bytepunx/system-flow/flai/internal/harness"
	"github.com/bytepunx/system-flow/flai/internal/hostapi"
	"github.com/bytepunx/system-flow/flai/internal/kluster"
	"github.com/bytepunx/system-flow/flai/internal/sendback"
	"github.com/bytepunx/system-flow/flai/internal/storygit"
	"github.com/bytepunx/system-flow/flai/internal/workflow"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// The state job (ADR-0312 §4, S-0810). When a story enters a state whose
// runs name something (workflow.StateRuns), flai serve's verify pass starts
// one job of kind closeout.KindState for that entry, keyed by the time of the
// transition that entered the state: once per entry, and again on a later
// entry. Its steps are the rebase of the story's branch onto the main branch,
// as the Rebase and resolve job rebases it (branchRebase), then each suite
// the state names, run as flai verify S-nnnn --all --suite <suite>, in the
// story's container when its agent runs in one (ADR-0274), then the cluster's
// steps, deploy then checks, each of which waits for flai serve's kluster
// pass (kluster.go) to record that step at the head the job left, as passed
// or failed: the pass deploys a story in a state whose runs name the deploy,
// once the job has reached its cluster steps, and keeps its own thread and
// send back. A cluster step the pass would not run, as klusterSkip says why
// (ADR-0312 §5), is skipped with that reason and never holds the story; a
// skipped deploy skips the checks. A step with nothing to run is recorded
// skipped with its reason, and when no step after the rebase runs, the
// rebase is skipped too and nothing runs at all: on the baseline, whose
// review runs only the cluster's steps, a project with no cluster rebases
// nothing on entering review and sends nothing back. The job's state is
// S-nnnn.state.json and its output S-nnnn.state.log, apart from the
// close-out's, and its progress lines are the state job's own and flai
// verify's. A job that passes leaves the story in the state; a rebase that
// stops or a suite that fails sends it back as flai by the state's
// send-back transition, with the send back's thread
// (sendback.StateJobTitle) that flai serve starts its agent on, unless the
// story left the state or entered it again meanwhile, or the state has no
// send-back transition, when the failure is recorded and nothing moves. A
// cluster step that fails sends nothing back from the job: the kluster pass
// sends the story back, once, with its own thread.

const (
	// stateEvery is how often the verify pass looks for a story that entered
	// a state whose runs name something.
	stateEvery = 10 * time.Second
	// stateMethod is the journal method of a state job that ran a step,
	// which flai stats does not count among the close-outs.
	stateMethod = "state.job"
	// stateOutputLines is how many lines of a failed suite's output the send
	// back's thread quotes.
	stateOutputLines = 20
	// allSuites is the name of the suite step of a state whose runs name
	// every suite.
	allSuites = "all"
	// klusterPollLeast is the shortest a cluster step waits between two looks
	// at what the kluster pass recorded; it looks as often as the verify
	// pass looks for states, stateEvery, and no more often than this.
	klusterPollLeast = 100 * time.Millisecond
)

// lookStates starts a state job for each story, not archived, that is in a
// state whose runs name something and whose entry into it has had no job,
// or one that flai serve's restart interrupted; a story a job of which runs
// is left to it and looked at again at a later look. While an upgrade waits
// on the drain it starts none: the entries run once it lifts (ADR-0320).
func (p *verifyPass) lookStates(ctx context.Context) {
	if _, holding := p.drain.holding(); holding {
		return
	}
	repo, err := openRepo(p.o, p.e.Root)
	if err != nil {
		return
	}
	flow := repo.Flow()
	if !slices.ContainsFunc(flow.States, func(s workflow.State) bool { return flow.StateRuns(s.Name) != nil }) {
		return
	}
	items, err := repo.List(false)
	if err != nil {
		p.o.Logger.Warn("state jobs not looked for", "component", "serve", "root", p.e.Root, "err", err.Error())
		return
	}
	for _, it := range items {
		if ctx.Err() != nil {
			return
		}
		if it.Type != workitem.Story {
			continue
		}
		runs := flow.StateRuns(it.Status)
		if runs == nil {
			continue
		}
		entry := it.EnteredAt().UTC()
		if p.busy(it.ID) || p.entryHadJob(it.ID, it.Status, entry) {
			continue
		}
		p.startState(ctx, repo, it, runs, entry)
	}
}

// busy reports whether a job of story runs.
func (p *verifyPass) busy(story string) bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	_, ok := p.running[story]
	return ok
}

// entryHadJob reports whether story's entry into state at entry had a state
// job that was not interrupted; a state job that cannot be read is taken as
// none, so that the next one writes over it.
func (p *verifyPass) entryHadJob(story, state string, entry time.Time) bool {
	j, ok, err := closeout.ReadStateJob(filepath.Join(p.e.Root, ".flai-cache"), story)
	if err != nil {
		p.o.Logger.Warn("state job not read", "component", "serve", "root", p.e.Root, "story", story, "err", err.Error())
		return false
	}
	return ok && j.State == state && j.Entry.Equal(entry) && j.Verdict != closeout.VerdictInterrupted
}

// stateStep is a step of a state job with how it runs.
type stateStep struct {
	closeout.StateStep
	// run runs the step, its output to out, and answers how it ended,
	// closeout.StepPassed, StepFailed, or StepSkipped, and why; nil for a
	// step skipped from the start, whose Reason the plan gives.
	run func(ctx context.Context, out io.Writer) (outcome, reason string)
}

// stepEnd is how a step that passed or failed ended, with why.
func stepEnd(passed bool, reason string) (outcome, why string) {
	if passed {
		return closeout.StepPassed, reason
	}
	return closeout.StepFailed, reason
}

// stateJob is one state job of a story: its rebase, its steps in the order
// they run, and what its failed suite printed, which only the job's own
// goroutine reads and writes.
type stateJob struct {
	branchRebase
	id    string
	stamp string
	state string
	// entry is when the story entered the state, the entry the job runs for.
	entry time.Time
	// agent is the story's agent when it runs on the container harness,
	// whose suites run in its container; nil runs them on the host.
	agent *containerAgent
	steps []stateStep
	// containers counts the containers the job's suites ran in, which name
	// them apart.
	containers int
	// said is the last lines of the output of the suite that failed.
	said string
}

// newStateJob is the job of story it's entry, at entry, into its state,
// whose runs are runs, started at stamp, with its steps planned.
func newStateJob(p *verifyPass, repo *workitem.Repo, it *workitem.Item, entry time.Time, stamp string, runs *workflow.Runs) *stateJob {
	story := it.ID
	j := &stateJob{
		branchRebase: branchRebase{p: p, story: story, worktree: repo.WorktreePath(story), by: storygit.SentBackAuthor, source: closeout.SourceState},
		id:           story + "-" + stamp, stamp: stamp, state: it.Status, entry: entry, agent: p.containerAgentOf(repo, story),
	}
	j.container = j.agent != nil
	j.plan(repo, it, runs)
	return j
}

// plan sets the job's steps from runs, for story it in repo: the rebase,
// each suite, or one for every suite, then the cluster's steps; the rebase
// is skipped when no step after it runs.
func (j *stateJob) plan(repo *workitem.Repo, it *workitem.Item, runs *workflow.Runs) {
	names, all := runs.SuiteNames()
	if all {
		names = []string{allSuites}
	}
	var after []stateStep
	for _, name := range names {
		after = append(after, stateStep{StateStep: closeout.StateStep{Kind: closeout.StateStepSuite, Name: name},
			run: func(ctx context.Context, out io.Writer) (string, string) { return stepEnd(j.runSuite(ctx, out, name)) }})
	}
	after = append(after, j.klusterSteps(repo, it, runs)...)
	rebase := stateStep{StateStep: closeout.StateStep{Kind: closeout.StateStepRebase, Name: closeout.StepRebase},
		run: func(ctx context.Context, out io.Writer) (string, string) { return stepEnd(j.runRebase(ctx, out)) }}
	if !slices.ContainsFunc(after, func(s stateStep) bool { return s.run != nil }) {
		rebase.run, rebase.Reason = nil, "no step after it runs, so the branch is not rebased"
	}
	j.steps = append([]stateStep{rebase}, after...)
}

// klusterSteps are the cluster's steps runs names, for story it in repo, in
// the order they run: each waits for what flai serve's kluster pass records
// of it at the job's head (awaitKluster), or is skipped from the start with
// why the pass runs nothing of it (klusterSkip, ADR-0312 §5). A skipped
// deploy skips the checks.
func (j *stateJob) klusterSteps(repo *workitem.Repo, it *workitem.Item, runs *workflow.Runs) []stateStep {
	var out []stateStep
	deploys, skipped := runs.RunsStep(workflow.KlusterDeploy), ""
	for _, name := range workflow.KlusterSteps {
		if !runs.RunsStep(name) {
			continue
		}
		s := stateStep{StateStep: closeout.StateStep{Kind: closeout.StateStepKluster, Name: name}}
		why := klusterSkip(j.p.o, j.p.e.Root, repo, it, name)
		if name == workflow.KlusterChecks && skipped != "" {
			why = "kluster deploy is skipped before it: " + skipped
		}
		if why == "" {
			s.run = func(ctx context.Context, out io.Writer) (string, string) {
				return j.awaitKluster(ctx, out, name, deploys)
			}
		} else {
			s.Reason = why
		}
		if name == workflow.KlusterDeploy {
			skipped = why
		}
		out = append(out, s)
	}
	return out
}

// runsAny reports whether a step of the job runs.
func (j *stateJob) runsAny() bool {
	return slices.ContainsFunc(j.steps, func(s stateStep) bool { return s.run != nil })
}

// records are the job's steps as its state keeps them.
func (j *stateJob) records() []closeout.StateStep {
	out := make([]closeout.StateStep, len(j.steps))
	for i, s := range j.steps {
		out[i] = s.StateStep
	}
	return out
}

// startState starts the job of it's entry, at entry, into its state, whose
// runs are runs. A job that runs nothing is recorded ended at once, every
// step skipped with its reason, nothing rebased, journaled, or announced;
// any other is recorded and journaled as it starts and runs in the
// background, the story's one job meanwhile.
func (p *verifyPass) startState(ctx context.Context, repo *workitem.Repo, it *workitem.Item, runs *workflow.Runs, entry time.Time) {
	dir := p.dir()
	started := p.o.Now().UTC()
	j := newStateJob(p, repo, it, entry, started.Format("20060102T150405Z"), runs)
	job := closeout.Job{ID: j.id, Story: it.ID, Kind: closeout.KindState, By: j.by, Started: started, State: j.state, Entry: entry}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		p.o.Logger.Warn("state job not started", "component", "serve", "root", p.e.Root, "story", it.ID, "state", j.state, "err", err.Error())
		return
	}
	if !j.runsAny() {
		for i := range j.steps {
			j.steps[i].Ended, j.steps[i].Outcome = started, closeout.StepSkipped
		}
		code := 0
		job.StateSteps, job.Ended, job.ExitCode, job.Verdict, job.Outcome = j.records(), started, &code, closeout.VerdictPassed, closeout.OutcomeStayedInState
		job.Reason = fmt.Sprintf("no step of %s's runs has anything to run, so nothing ran and %s stays in %s", j.state, it.ID, j.state)
		if err := closeout.WriteJob(dir, job); err != nil {
			p.o.Logger.Warn("state job not recorded", "component", "serve", "root", p.e.Root, "story", it.ID, "state", j.state, "err", err.Error())
			return
		}
		p.o.Logger.Info("state job ran nothing", "component", "serve", "root", p.e.Root, "story", it.ID, "state", j.state, "job_id", job.ID)
		return
	}
	job.StateSteps = j.records()
	// marked running before it is recorded, so that an entry the drain
	// holds back is not recorded as a job that runs
	if !p.track(it.ID, job.ID, j.state) {
		return
	}
	if err := closeout.WriteJob(dir, job); err != nil {
		p.forgetState(it.ID)
		p.o.Logger.Warn("state job not started", "component", "serve", "root", p.e.Root, "story", it.ID, "state", j.state, "err", err.Error())
		return
	}
	p.journal(hostapi.Entry{At: started.Format(time.RFC3339), Method: stateMethod, By: job.By, Outcome: "started", Detail: fmt.Sprintf("%s job %s on entering %s", it.ID, job.ID, j.state)})
	p.o.Logger.Info("state job started", "component", "serve", "root", p.e.Root, "story", it.ID, "state", j.state, "job_id", job.ID, "entry", entry.Format(time.RFC3339))
	p.wg.Add(1)
	go func() {
		defer p.wg.Done()
		p.runState(ctx, dir, j, job)
	}()
}

// runState runs j's steps in order, each between the progress lines that say
// it started and how it ended, for at most the pass's limit, its output in
// the job's log, writing the job's state as each starts and ends; once one
// fails, those after it are skipped. Then it settles how the job ended. When
// ctx ends first, the job is left running in its state, for the next flai
// serve to end interrupted and run again, and nothing is settled.
func (p *verifyPass) runState(ctx context.Context, dir string, j *stateJob, job closeout.Job) {
	story := job.Story
	logPath := closeout.File(dir, story, closeout.StateLogSuffix)
	out := io.Discard
	if f, err := os.Create(logPath); err != nil {
		p.o.Logger.Warn("state job's log not made", "component", "serve", "root", p.e.Root, "story", story, "job_id", job.ID, "err", err.Error())
	} else {
		defer func() { _ = f.Close() }()
		out = f
	}
	line := func(s string) { fmt.Fprintln(out, s) }
	write := func() {
		job.StateSteps = j.records()
		if err := closeout.WriteJob(dir, job); err != nil {
			p.o.Logger.Warn("state job's state not written", "component", "serve", "root", p.e.Root, "story", story, "job_id", job.ID, "err", err.Error())
		}
	}
	job.Head = p.head(j.worktree)
	jctx, cancel := context.WithTimeout(ctx, p.limit)
	defer cancel()
	failed := -1
	for i := range j.steps {
		s := &j.steps[i]
		if s.run == nil || failed >= 0 {
			if s.run != nil {
				s.Reason = j.steps[failed].Label() + " failed before it"
			}
			s.Ended, s.Outcome = p.o.Now().UTC(), closeout.StepSkipped
			line(closeout.EndedLine(closeout.SourceState, story, closeout.LineStep, closeout.StepSkipped, s.Label(), s.Ended))
			continue
		}
		s.Started = p.o.Now().UTC()
		line(closeout.StartedLine(closeout.SourceState, story, closeout.LineStep, s.Label(), s.Started))
		write()
		outcome, reason := s.run(jctx, out)
		if ctx.Err() != nil {
			p.forgetState(story)
			p.o.Logger.Info("state job stopped with flai serve", "component", "serve", "root", p.e.Root, "story", story, "job_id", job.ID)
			return
		}
		if errors.Is(jctx.Err(), context.DeadlineExceeded) {
			outcome, reason = closeout.StepFailed, fmt.Sprintf("%s; stopped at the %s limit", reason, p.limit)
		}
		s.Ended, s.Outcome, s.Reason = p.o.Now().UTC(), outcome, reason
		if outcome == closeout.StepFailed {
			failed = i
		}
		line(closeout.EndedLine(closeout.SourceState, story, closeout.LineStep, s.Outcome, s.Label(), s.Ended))
		write()
	}
	ended := p.o.Now().UTC()
	code := 0
	job.Verdict = closeout.VerdictPassed
	if failed >= 0 {
		code, job.Verdict, job.StoppedAt = 1, closeout.VerdictStopped, j.steps[failed].Label()
	}
	job.Ended, job.ExitCode, job.DurationMS, job.Commit = ended, &code, ended.Sub(job.Started).Milliseconds(), p.head(j.worktree)
	// a project that cannot be opened is nil, which settle sends nothing back for
	repo, _ := openRepo(p.o, p.e.Root)
	j.settle(repo, &job, failed)
	write()
	p.forgetState(story)
	duration := time.Duration(job.DurationMS) * time.Millisecond
	detail := fmt.Sprintf("%s job %s on entering %s: %s%s in %s; %s", story, job.ID, job.State, job.Verdict, atStep(job.StoppedAt), duration.Round(time.Second), job.Outcome)
	if job.Thread != "" {
		detail += " on " + job.Thread
	}
	p.journal(hostapi.Entry{At: ended.Format(time.RFC3339), Method: stateMethod, By: job.By, Outcome: closeOutOutcome(job.Verdict), Verdict: job.Verdict, DurationMS: job.DurationMS, Detail: detail})
	p.o.Logger.Info("state job ended", "component", "serve", "root", p.e.Root, "story", story, "state", job.State, "job_id", job.ID,
		"verdict", job.Verdict, "stopped_at", job.StoppedAt, "duration_ms", job.DurationMS, "outcome", job.Outcome, "thread", job.Thread)
}

// forgetState marks story's state job as no longer running.
func (p *verifyPass) forgetState(story string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.untrack(story)
	delete(p.states, story)
}

// runRebase is the job's rebase step: this flai's flai stream sync of the
// story in its worktree, a rebase that stops on conflicts aborted, as the
// Rebase and resolve job runs it.
func (j *stateJob) runRebase(ctx context.Context, out io.Writer) (bool, string) {
	if info, err := os.Stat(j.worktree); err != nil || !info.IsDir() {
		return false, fmt.Sprintf("%s has no worktree at %s, so its branch cannot be rebased nor its suites run: open it with flai stream open %s", j.story, j.worktree, j.story)
	}
	branch, base := storygit.Branch(j.story), j.base()
	code, err := j.rebase(ctx, out)
	switch {
	case err != nil:
		return false, fmt.Sprintf("flai stream sync %s could not rebase %s onto %s: %v", j.story, branch, base, err)
	case code == 0:
		return true, fmt.Sprintf("rebased %s onto %s", branch, base)
	case len(j.conflicts) > 0 || j.waits != "":
		state := "which flai aborted, so the branch is as it was"
		if j.waits != "" {
			state = "which still waits in the worktree: " + j.waits
		}
		on := ""
		if len(j.conflicts) > 0 {
			on = " on conflicts in " + workitem.Shorten(j.conflicts, 10)
		}
		return false, fmt.Sprintf("the rebase of %s onto %s stopped%s, %s", branch, base, on, state)
	}
	return false, fmt.Sprintf("flai stream sync %s exited %d and left no rebase waiting, as when the worktree has uncommitted changes or git failed", j.story, code)
}

// klusterThreadTitle is the title of the thread the kluster pass opens on
// story when its step, workflow.KlusterDeploy or KlusterChecks, fails.
func klusterThreadTitle(step, story string) string {
	if step == workflow.KlusterChecks {
		return checksFailureTitle(story)
	}
	return klusterFailureTitle(story)
}

// awaitKluster is the job's cluster step step, workflow.KlusterDeploy or
// KlusterChecks, in a state whose runs name the deploy when deploys: it
// waits, looking as often as the verify pass looks for states, until what
// flai serve's kluster pass recorded of the story at the head the job left
// in its worktree says how the step ended (klusterSaid), for as long as ctx,
// bounded by the job's limit, lets it. Checks after a deploy that was
// skipped are skipped.
func (j *stateJob) awaitKluster(ctx context.Context, out io.Writer, step string, deploys bool) (string, string) {
	p := j.p
	if step == workflow.KlusterChecks {
		i := slices.IndexFunc(j.steps, func(s stateStep) bool { return s.Kind == closeout.StateStepKluster && s.Name == workflow.KlusterDeploy })
		if i >= 0 && j.steps[i].Outcome == closeout.StepSkipped {
			return closeout.StepSkipped, "kluster deploy was skipped before it: " + j.steps[i].Reason
		}
	}
	head := p.head(j.worktree)
	if head == "" {
		return closeout.StepFailed, fmt.Sprintf("git could not say the head of %s's worktree %s, so what flai serve's kluster pass recorded cannot be judged at it: check that git can read the worktree", j.story, j.worktree)
	}
	fmt.Fprintf(out, "state: waiting for flai serve's kluster pass to record the %s of %s at %s\n", step, j.story, short(head))
	poll := max(p.stateEvery, klusterPollLeast)
	for {
		if outcome, reason := j.klusterSaid(step, head, deploys); outcome != "" {
			fmt.Fprintf(out, "state: kluster %s %s: %s\n", step, outcome, reason)
			return outcome, reason
		}
		select {
		case <-ctx.Done():
			return closeout.StepFailed, fmt.Sprintf("flai serve's kluster pass recorded no %s of %s at %s", step, j.story, short(head))
		case <-time.After(poll):
		}
	}
}

// klusterSaid is how the job's cluster step step ended, by what the kluster
// pass recorded of the story at head, and why, or "" while nothing tells
// yet. It passed when the namespace's record passed at the head, as the
// acceptance gate judges it (kluster.AtHead), and for the checks step when
// a checks run passed there too. It failed when the pass tried the head in
// the job's state, in this entry or in an earlier one it sent nothing back
// for, which it does not try again, and the deploy failed, or the checks run
// at the head did not pass; or when the story left the entry. It is skipped
// when the pass no longer runs the step (klusterSkip), or, for the checks of
// a state whose runs name no deploy, when no deploy that passed is at the
// head for them to run against. A project or a record that cannot be read
// tells nothing yet, and is read again at the next look.
func (j *stateJob) klusterSaid(step, head string, deploys bool) (outcome, reason string) {
	p, story := j.p, j.story
	repo, err := openRepo(p.o, p.e.Root)
	if err != nil {
		return "", ""
	}
	it, err := repo.Get(story)
	if err != nil {
		p.o.Logger.Debug("state job's story not read for its cluster step", "component", "serve", "root", p.e.Root, "story", story, "step", step, "err", err.Error())
		return "", ""
	}
	root := mainCheckout(repo)
	ns, err := kluster.StoryNamespace(story)
	if err != nil {
		return closeout.StepSkipped, fmt.Sprintf("flai serve's kluster pass deploys no namespace for %s: %v", story, err)
	}
	r, recorded, err := kluster.ReadRecord(root, ns)
	if err != nil {
		p.o.Logger.Debug("kluster record not read for a state job's cluster step", "component", "serve", "root", p.e.Root, "story", story, "step", step, "err", err.Error())
		return "", ""
	}
	at := func(commit string) bool {
		judged, err := kluster.AtHead(gitOf(p.o), root, repo.Manifest, commit, head)
		return err == nil && judged.AtHead
	}
	t, _ := readKlusterTried(root, ns)
	triedHere := t.Head == head && t.triedIn(j.state) && (!t.At.Before(j.entry) || t.SentBack == nil)
	deployed := recorded && r.Outcome == kluster.OutcomePassed && at(r.Commit)
	checks := step == workflow.KlusterChecks
	switch {
	case !checks && deployed:
		return closeout.StepPassed, fmt.Sprintf("flai serve's kluster pass deployed %s into its namespace %s at %s, at the head %s", story, ns, short(r.Commit), short(head))
	case !checks && triedHere && t.Outcome == kluster.OutcomeFailed:
		return closeout.StepFailed, fmt.Sprintf("flai serve's kluster deploy of %s into its namespace %s at %s failed: flai kluster deployments shows its record", story, ns, short(head))
	case checks && deployed && r.ChecksOutcome == kluster.OutcomePassed && at(r.ChecksCommit):
		return closeout.StepPassed, fmt.Sprintf("flai serve's checks of %s against its namespace %s passed at %s, at the head %s", story, ns, short(r.ChecksCommit), short(head))
	case checks && triedHere && r.ChecksCommit == head && r.ChecksOutcome != "" && r.ChecksOutcome != kluster.OutcomePassed:
		return closeout.StepFailed, fmt.Sprintf("flai serve's checks of %s against its namespace %s at %s ended %s: flai checks status %s shows the run", story, ns, short(head), r.ChecksOutcome, story)
	case checks && !deploys && !deployed:
		return closeout.StepSkipped, fmt.Sprintf("%s's runs name no deploy, and %s's namespace %s is not deployed at its head %s with a deploy that passed, so its checks have nothing to run against", j.state, story, ns, short(head))
	case it.Status != j.state || !it.EnteredAt().UTC().Equal(j.entry):
		return closeout.StepFailed, fmt.Sprintf("%s left %s for %s before flai serve's kluster pass recorded its %s at %s", story, j.state, it.Status, step, short(head))
	}
	if why := klusterSkip(p.o, p.e.Root, repo, it, step); why != "" {
		return closeout.StepSkipped, why
	}
	return "", ""
}

// verifyArgv is flai's verify of the job's story for suite, with no tier left
// to the merge queue, as the executable flai names it.
func (j *stateJob) verifyArgv(flai, suite string) []string {
	argv := []string{flai, "verify", j.story, "--all"}
	if suite != allSuites {
		argv = append(argv, "--suite", suite)
	}
	return argv
}

// runSuite is the job's step for suite: flai verify of the story limited to
// it, or for every suite, in the worktree on the host, or in the story's
// container when its agent runs in one, which nothing from the worktree runs
// beside on the host (ADR-0274). It keeps the last lines of what a failing
// run printed, for the send back's thread.
func (j *stateJob) runSuite(ctx context.Context, out io.Writer, suite string) (bool, string) {
	p := j.p
	words := strings.Join(j.verifyArgv("flai", suite), " ")
	var run, env []string
	if j.agent != nil {
		c := &containerCloseOut{p: p, agent: *j.agent, story: j.story, worktree: j.worktree, req: closeout.Request{Story: j.story, By: j.by}, id: j.stamp}
		host, image, err := c.image(ctx, out)
		if err == nil && image == "" {
			err = fmt.Errorf("the container harness has no image for the inner harness %s of %s's agent: name one with flai serve agent harness container --image %s=<image>", innerOf(j.agent.agent), j.story, innerOf(j.agent.agent))
		}
		if err != nil {
			return false, fmt.Sprintf("%s's agent runs in a container, so %s runs there and nothing runs on the host: %v", j.story, words, err)
		}
		j.containers++
		cmd, err := harness.ContainerCommand(harness.CommandRun{Root: p.e.Root, Story: j.story, Purpose: "state", ID: fmt.Sprintf("%s-%d", j.stamp, j.containers), Image: image,
			Env: []string{"FLAI_AGENT", closeout.ProgressEnv}, Argv: j.verifyArgv(harness.ImageFlai, suite)}, host)
		if err != nil {
			return false, fmt.Sprintf("%s could not run in %s's container: %v", words, j.story, err)
		}
		run, env = cmd.Argv, append(closeOutEnv(withoutGitEnv(os.Environ()), j.by), cmd.Env...)
		fmt.Fprintf(out, "state: %s in the container %s\n", words, cmd.Name)
	} else {
		flai, err := p.exe()
		if err != nil {
			return false, fmt.Sprintf("this flai's executable is not known, so %s cannot run: run it in the story's worktree yourself: %v", words, err)
		}
		run, env = j.verifyArgv(flai, suite), closeOutEnv(os.Environ(), j.by)
	}
	var said bytes.Buffer
	code, err := p.run(ctx, j.worktree, env, io.MultiWriter(out, &said), run)
	switch {
	case err != nil:
		return false, fmt.Sprintf("%s could not run: %v", words, err)
	case code != 0:
		j.said = outputTail(said.String(), j.story, stateOutputLines)
		return false, fmt.Sprintf("%s exited %d: %s", words, code, lastLine(said.String()))
	}
	return true, words + " passed"
}

// outputTail is the last n lines of what a flai verify of story printed,
// but for its progress lines, which name no finding.
func outputTail(said, story string, n int) string {
	var lines []string
	for line := range strings.Lines(said) {
		line = strings.TrimRight(line, "\r\n")
		if strings.HasPrefix(line, closeout.SourceVerify+": "+story+" ") && (strings.Contains(line, " step started ") || strings.Contains(line, " step ended ") ||
			strings.Contains(line, " tier started ") || strings.Contains(line, " tier ended ")) {
			continue
		}
		lines = append(lines, line)
	}
	return strings.Join(lines[max(0, len(lines)-n):], "\n")
}

// settle sets how job ended once its steps did, failed the index of the step
// that failed or -1: none failed leaves the story in the state; one that
// failed sends the story, read again, back by the state's send-back
// transition, as storygit.SentBackAuthor, with the send back's thread naming
// the state, the step, its findings or the paths in conflict, and the job's
// log, recording both in job. A story that left the state, or entered it
// again, since the entry the job ran for is left as it is, as is one in a
// state with no send-back transition or a project that could not be opened,
// repo nil, and job's reason says so. A cluster step that failed moves
// nothing: the kluster pass sends the story back itself, with its own
// thread, so that it is sent back once.
func (j *stateJob) settle(repo *workitem.Repo, job *closeout.Job, failed int) {
	p, story, state := j.p, j.story, j.state
	if failed < 0 {
		job.Outcome = closeout.OutcomeStayedInState
		job.Reason = fmt.Sprintf("every step of %s's job that ran passed at %s, so %s stays in %s", state, short(job.Commit), story, state)
		return
	}
	job.Conflicts = j.conflicts
	s := j.steps[failed].StateStep
	log := closeout.File(p.dir(), story, closeout.StateLogSuffix)
	what := s.Label() + " failed"
	if s.Kind == closeout.StateStepRebase && len(j.conflicts) > 0 {
		what = "the rebase stopped on conflicts in " + workitem.Shorten(j.conflicts, 3)
	}
	reason := fmt.Sprintf("state job %s on entering %s: %s: %s; its output is in %s", job.ID, state, what, s.Reason, log)
	job.Outcome, job.Reason = closeout.OutcomeNotSentBack, reason
	if s.Kind == closeout.StateStepKluster {
		// the kluster pass sends the story back after its own failure, once,
		// with its own thread, whether or not it has done so yet
		job.Reason = fmt.Sprintf("%s; a failed cluster step is flai serve's kluster pass's to send %s back for, by %s's send-back transition with its thread %q, so the job moved nothing",
			reason, story, state, klusterThreadTitle(s.Name, story))
		return
	}
	warn := func(msg string, err error) {
		p.o.Logger.Warn(msg, "component", "serve", "root", p.e.Root, "story", story, "job_id", job.ID, "err", err.Error())
	}
	if repo == nil {
		job.Reason = fmt.Sprintf("%s; but flai serve could not open the project at %s, so nothing was sent back: tell its agent with flai thread new --on %s", reason, p.e.Root, story)
		return
	}
	it, err := repo.Get(story)
	if err != nil {
		job.Reason = fmt.Sprintf("%s; but flai serve could not read %s (%v), so nothing was sent back: tell its agent with flai thread new --on %s", reason, story, err, story)
		warn("state job's story not read", err)
		return
	}
	if it.Status != state || !it.EnteredAt().Equal(job.Entry) {
		job.Reason = fmt.Sprintf("%s is in %s since %s, no longer in the entry into %s of %s that the job ran for, so flai left it as it is: %s",
			story, it.Status, it.EnteredAt().UTC().Format(time.RFC3339), state, job.Entry.Format(time.RFC3339), reason)
		return
	}
	t, ok := repo.Flow().SendBack(it.Type, state)
	if !ok {
		job.Reason = fmt.Sprintf("%s; the workflow has no send-back transition out of %s, so flai moved nothing and %s stays in %s", reason, state, story, state)
		return
	}
	now := p.o.Now()
	if _, err := repo.Transition(it, t.To, storygit.SentBackAuthor, reason, now); err != nil {
		job.Reason = fmt.Sprintf("%s; but flai could not move %s back to %s (%v): the operator moves it with flai move %s %s", reason, story, t.To, err, story, t.To)
		warn("state job's story not sent back", err)
		return
	}
	job.Outcome = closeout.OutcomeSentBack
	th, err := sendback.Thread(repo, story, sendback.StateJobTitle(story, state, what), j.threadText(job, s, log, t.To), now)
	if err != nil {
		job.Reason += fmt.Sprintf("; but its thread could not be opened (%v): tell its agent with flai thread new --on %s", err, story)
		warn("state job's send back thread not opened", err)
		return
	}
	job.Thread = th.ID
	p.o.Logger.Info("story sent back from its state job", "component", "serve", "root", p.e.Root, "story", story, "job_id", job.ID,
		"state", state, "to", t.To, "stopped_at", job.StoppedAt, "paths", j.conflicts, "thread", th.ID)
}

// threadText is what the send back's thread tells the story's agent when
// job's step s failed: the state, the step, the paths in conflict or the
// last lines of what the suite printed, where the job's output is, and how
// the agent fixes it and brings the story back to the state, where the job
// runs again.
func (j *stateJob) threadText(job *closeout.Job, s closeout.StateStep, log, back string) string {
	var b strings.Builder
	story, branch := j.story, storygit.Branch(j.story)
	fmt.Fprintf(&b, "flai serve's job %s, which runs what the workflow's state %s runs each time a story enters it, stopped at its step %s: %s.\n\n", job.ID, j.state, s.Label(), s.Reason)
	switch {
	case len(j.conflicts) > 0:
		fmt.Fprintf(&b, "The rebase of `%s` onto `%s` stopped on conflicts in:\n\n", branch, j.base())
		for _, path := range j.conflicts {
			fmt.Fprintf(&b, "- `%s`\n", path)
		}
		b.WriteString("\n")
	case j.said != "":
		b.WriteString("The last lines it printed:\n\n```\n" + j.said + "\n```\n\n")
	}
	fmt.Fprintf(&b, "Its whole output is in `%s`.\n\n", log)
	fmt.Fprintf(&b, "flai moved %s back to %s. "+
		"This thread is for %s's agent, which fixes it; the operator does not. "+
		"In its worktree, run `flai stream sync %s`, resolve each path it stops on or fix what failed, close out again, with the MCP tool `verify` or `flai verify %s --job`, and move %s on until it is in %s again. "+
		"flai serve resolves this thread once %s is back in %s, and runs the state's job again.",
		story, back, story, story, story, story, j.state, story, j.state)
	return b.String()
}
