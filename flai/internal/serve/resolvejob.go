package serve

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/bytepunx/system-flow/flai/internal/closeout"
	"github.com/bytepunx/system-flow/flai/internal/sendback"
	"github.com/bytepunx/system-flow/flai/internal/storygit"
	"github.com/bytepunx/system-flow/flai/internal/verify"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// The Rebase and resolve job (S-0804): the operator's flai accept --resolve,
// or the review page's button, asks flai serve's verify pass for a job of
// kind closeout.KindResolve on a story in review. With no agent, the job
// rebases the story's branch onto the main branch with this flai's flai
// stream sync, then runs the story's close-out at the rebased head exactly as
// a close-out job does, its output and its own progress lines in the job's
// log. A clean rebase and a passing close-out leave the story in review.
// Only a rebase that stops on conflicts, which the job aborts, or a step that
// fails sends the story back to in-progress, with the send back's thread by
// flai (package sendback) naming the paths, or the step, its findings, and
// the log, which flai serve starts the story's agent on (ADR-0248).

// resolveFindings is how many of the stopped step's findings the send back's
// thread lists; it counts the rest.
const resolveFindings = 10

// resolveMessage is the commit message the close-out of story's rebase and
// resolve commits what is outstanding with, when the request gives none.
func resolveMessage(story string) string {
	return fmt.Sprintf("chore: [%s] close out again after rebase and resolve", story)
}

// resolveRefusal is why story's rebase and resolve does not run: it cannot
// be read, or it is not in review; "" when it may run.
func resolveRefusal(repo *workitem.Repo, story string) string {
	it, err := repo.Get(story)
	if err != nil {
		return fmt.Sprintf("flai serve could not read %s (%v), so its rebase and resolve does not run: run flai check, and ask again", story, err)
	}
	if it.Status != workitem.Review {
		return fmt.Sprintf("%s is in %s, not in review: rebase and resolve runs only for a story in review, so nothing runs", story, it.Status)
	}
	return ""
}

// resolveJob is one rebase and resolve job of a story: its run, the job's
// body, and its settle, which runJob calls once the job ended, from the same
// goroutine, so what the run found needs no lock.
type resolveJob struct {
	p        *verifyPass
	story    string
	worktree string
	by       string
	// container is whether the story's agent runs in a container, whose
	// worktree's git files the host checks before each of its own steps
	// there (ADR-0274).
	container bool
	// closeOut is the story's close-out, as a close-out job runs it.
	closeOut jobBody

	// rebased is whether the sync rebased the branch cleanly.
	rebased bool
	// conflicts are the paths the sync's rebase stopped on.
	conflicts []string
	// waits says why the rebase that stopped still waits in the worktree:
	// one the sync found there, which the job leaves as it found it, or one
	// it could not abort; "" when it aborted it, or none stopped.
	waits string
}

// newResolveJob is the rebase and resolve of story in worktree, which by
// asked for, running closeOut once the rebase is clean.
func newResolveJob(p *verifyPass, story, worktree, by string, container bool, closeOut jobBody) *resolveJob {
	return &resolveJob{p: p, story: story, worktree: worktree, by: by, container: container, closeOut: closeOut}
}

// run is the job's body: the rebase, then, when it is clean, the close-out,
// each between the progress lines that say it started and how it ended. It
// answers the sync's exit code when the rebase did not end clean, else the
// close-out's.
func (r *resolveJob) run(ctx context.Context, out io.Writer) (int, error) {
	line := func(s string) { fmt.Fprintln(out, s) }
	now := r.p.o.Now
	line(closeout.StartedLine(closeout.SourceResolve, r.story, closeout.LineStep, closeout.StepRebase, now()))
	code, err := r.rebase(ctx, out)
	if err != nil || code != 0 {
		line(closeout.EndedLine(closeout.SourceResolve, r.story, closeout.LineStep, "failed", closeout.StepRebase, now()))
		return code, err
	}
	r.rebased = true
	line(closeout.EndedLine(closeout.SourceResolve, r.story, closeout.LineStep, "passed", closeout.StepRebase, now()))
	line(closeout.StartedLine(closeout.SourceResolve, r.story, closeout.LineStep, closeout.StepCloseOut, now()))
	code, err = r.closeOut(ctx, out)
	outcome := "passed"
	if err != nil || code != 0 {
		outcome = "failed"
	}
	line(closeout.EndedLine(closeout.SourceResolve, r.story, closeout.LineStep, outcome, closeout.StepCloseOut, now()))
	return code, err
}

// rebase runs this flai's flai stream sync of the story in its worktree, its
// output to out, and answers its exit code. When the sync stops with a
// rebase waiting, it reads the paths in conflict and aborts the rebase, but
// one it found waiting before the sync, which it leaves as it was. When ctx
// ends during the sync, it leaves the worktree as the sync left it.
func (r *resolveJob) rebase(ctx context.Context, out io.Writer) (int, error) {
	p, root, git := r.p, r.p.e.Root, gitOf(r.p.o)
	flai, err := p.exe()
	if err != nil {
		return -1, fmt.Errorf("this flai's executable is not known, so %s's branch cannot be rebased: run flai stream sync %s in its worktree yourself: %w", r.story, r.story, err)
	}
	env := closeOutEnv(os.Environ(), r.by)
	if r.container {
		if err := gitPointers(root, r.story, r.worktree); err != nil {
			return -1, err
		}
		env = containerHostEnv(root, r.by)
	}
	found := storygit.RebaseInProgress(git, r.worktree)
	code, err := p.run(ctx, r.worktree, env, out, []string{flai, "stream", "sync", r.story})
	if err != nil || code == 0 || ctx.Err() != nil || !storygit.RebaseInProgress(git, r.worktree) {
		return code, err
	}
	if r.container {
		if err := gitPointers(root, r.story, r.worktree); err != nil {
			r.waits = err.Error()
			fmt.Fprintf(out, "resolve: the rebase was not aborted: %v\n", err)
			return code, nil
		}
	}
	r.conflicts = storygit.Conflicts(git, r.worktree)
	branch := storygit.Branch(r.story)
	if found {
		r.waits = fmt.Sprintf("a rebase of %s was already in progress in %s, which flai left as it found it", branch, r.worktree)
		fmt.Fprintf(out, "resolve: %s\n", r.waits)
		return code, nil
	}
	if err := storygit.AbortRebase(git, r.worktree, branch); err != nil {
		r.waits = fmt.Sprintf("flai could not abort the rebase: %v", err)
		fmt.Fprintf(out, "resolve: %s\n", r.waits)
		return code, nil
	}
	fmt.Fprintf(out, "resolve: flai aborted the rebase of %s, so the branch is as it was\n", branch)
	return code, nil
}

// settle sets how job ended once runJob gave it its verdict: a rebase that
// did not end clean is stopped at closeout.StepRebase; a passing close-out
// leaves the story in review; anything else sends it back, naming the paths
// in conflict, or the step, its findings in rep, the story's verify report,
// when newest, the job's own, and the job's log.
func (r *resolveJob) settle(repo *workitem.Repo, job *closeout.Job, rep verify.Report, newest bool) {
	branch, base := storygit.Branch(r.story), r.base()
	job.Conflicts = r.conflicts
	if !r.rebased {
		job.Verdict, job.StoppedAt = closeout.VerdictStopped, closeout.StepRebase
	}
	if job.Verdict == closeout.VerdictPassed {
		job.Outcome = closeout.OutcomeStayed
		job.Reason = fmt.Sprintf("the rebase of %s onto %s was clean and its close-out passed at %s, so %s stays in review, ready to accept", branch, base, short(job.Commit), r.story)
		return
	}
	log := closeout.File(r.p.dir(), r.story, closeout.LogSuffix)
	var reason, title string
	var text func(back string) string
	switch {
	case len(r.conflicts) > 0 || r.waits != "":
		state := "which flai aborted, so the branch is as it was"
		if r.waits != "" {
			state = "which still waits in the worktree: " + r.waits
		}
		on := ""
		if len(r.conflicts) > 0 {
			on = " on conflicts in " + workitem.Shorten(r.conflicts, 10)
		}
		reason = fmt.Sprintf("rebase and resolve job %s: the rebase of %s onto %s stopped%s, %s; its agent settles it in the worktree and syncs it with flai stream sync %s",
			job.ID, branch, base, on, state, r.story)
		title = storygit.SentBackTitle(branch, base, r.conflicts)
		text = func(back string) string { return r.conflictText(job, branch, base, log, back) }
	case !r.rebased:
		reason = fmt.Sprintf("rebase and resolve job %s: flai stream sync %s could not rebase %s onto %s, and its close-out did not run; its output is in %s; its agent syncs it, closes out again, and moves it to review",
			job.ID, r.story, branch, base, log)
		title = storygit.SentBackPrefix + fmt.Sprintf("%s could not be rebased onto %s by rebase and resolve", branch, base)
		text = func(back string) string { return r.stepText(job, branch, base, log, back, nil, 0) }
	default:
		findings, more := stepFindings(rep, newest, job.StoppedAt)
		reason = fmt.Sprintf("rebase and resolve job %s: %s's close-out after its clean rebase onto %s %s; its output is in %s; its agent fixes it, closes out again, and moves it to review",
			job.ID, branch, base, stoppedWhere(job.StoppedAt, len(findings)+more), log)
		title = storygit.SentBackPrefix + fmt.Sprintf("%s's close-out after its rebase onto %s %s", branch, base, stoppedWhere(job.StoppedAt, -1))
		text = func(back string) string { return r.stepText(job, branch, base, log, back, findings, more) }
	}
	r.sendBack(repo, job, reason, title, text)
}

// sendBack moves the story, read again, back to the state the workflow's
// send back takes it to, as storygit.SentBackAuthor with reason, and opens
// the send back's thread titled title and saying text, or writes on the one
// still open, recording both in job. A story no longer in review is left as
// it is, and job's reason says so; so is one in a project that could not be
// opened, repo nil.
func (r *resolveJob) sendBack(repo *workitem.Repo, job *closeout.Job, reason, title string, text func(back string) string) {
	p, now := r.p, r.p.o.Now()
	warn := func(msg string, err error) {
		p.o.Logger.Warn(msg, "component", "serve", "root", p.e.Root, "story", r.story, "job_id", job.ID, "err", err.Error())
	}
	if repo == nil {
		job.Reason = fmt.Sprintf("%s; but flai serve could not open the project at %s, so nothing was sent back: tell its agent with flai thread new --on %s", reason, p.e.Root, r.story)
		return
	}
	it, err := repo.Get(r.story)
	if err != nil {
		job.Reason = fmt.Sprintf("%s; but flai serve could not read %s (%v), so nothing was sent back: tell its agent with flai thread new --on %s", reason, r.story, err, r.story)
		warn("rebase and resolve's story not read", err)
		return
	}
	if it.Status != workitem.Review {
		job.Reason = fmt.Sprintf("%s is in %s, no longer in review, so flai left it as it is: %s", r.story, it.Status, reason)
		return
	}
	back := sendback.Target(repo.Flow(), it)
	if _, err := repo.Transition(it, back, storygit.SentBackAuthor, reason, now); err != nil {
		job.Reason = fmt.Sprintf("%s; but flai could not move %s back to %s (%v): the operator moves it with flai move %s %s", reason, r.story, back, err, r.story, back)
		warn("rebase and resolve's story not sent back", err)
		return
	}
	job.Outcome, job.Reason = closeout.OutcomeSentBack, reason
	th, err := sendback.Thread(repo, r.story, title, text(back), now)
	if err != nil {
		job.Reason += fmt.Sprintf("; but its thread could not be opened (%v): tell its agent with flai thread new --on %s", err, r.story)
		warn("rebase and resolve's send back thread not opened", err)
		return
	}
	job.Thread = th.ID
	p.o.Logger.Info("story sent back from rebase and resolve", "component", "serve", "root", p.e.Root, "story", r.story, "job_id", job.ID,
		"to", back, "stopped_at", job.StoppedAt, "paths", r.conflicts, "thread", th.ID)
}

// base is the main branch the story's branch was rebased onto, or words for
// it when git cannot say.
func (r *resolveJob) base() string {
	base, err := storygit.MainBranch(gitOf(r.p.o), r.p.e.Root)
	if err != nil || base == "" {
		return "the main branch"
	}
	return base
}

// askedBy is who asked for job, as the thread names them.
func askedBy(job *closeout.Job) string {
	if job.By == "" {
		return ""
	}
	return ", which " + job.By + " asked for,"
}

// conflictText is what the send back's thread tells the story's agent when
// job's rebase of branch onto base stopped on conflicts: each path, what flai
// did with the rebase, how the agent settles them, closes out, and comes back
// to review, and where the job's output is.
func (r *resolveJob) conflictText(job *closeout.Job, branch, base, log, back string) string {
	var b strings.Builder
	fmt.Fprintf(&b, "flai serve's rebase and resolve job %s%s rebased `%s` onto `%s`, and the rebase stopped on conflicts that the merge of the issue files, the ADR index, the issue summary, and the shared markdown files did not settle", job.ID, askedBy(job), branch, base)
	if len(r.conflicts) > 0 {
		b.WriteString(", in:\n\n")
		for _, path := range r.conflicts {
			fmt.Fprintf(&b, "- `%s`\n", path)
		}
		b.WriteString("\n")
	} else {
		b.WriteString(".\n\n")
	}
	story := r.story
	if r.waits != "" {
		fmt.Fprintf(&b, "The rebase still waits in `%s`: %s. flai moved %s back to %s. "+
			"This thread is for %s's agent, which settles the conflicts; the operator does not. "+
			"In its worktree, resolve each path that still conflicts, `git add` it, and run `flai stream sync %s`, which continues the rebase, or run `git rebase --abort` and then `flai stream sync %s`; ",
			r.worktree, r.waits, story, back, story, story, story)
	} else {
		fmt.Fprintf(&b, "flai aborted the rebase, so `%s` and `%s` are as they were, and moved %s back to %s. "+
			"This thread is for %s's agent, which settles the conflicts; the operator does not. "+
			"In its worktree, run `flai stream sync %s`, resolve each path it stops on, `git add` it, and run `flai stream sync %s` again, which continues the rebase; ",
			branch, base, story, back, story, story, story)
	}
	fmt.Fprintf(&b, "then close out again, with the MCP tool `verify` or `flai verify %s --job`, and move %s to review again. "+
		"flai serve resolves this thread once %s is back in review. The job's whole output is in `%s`.",
		story, story, story, log)
	return b.String()
}

// stepText is what the send back's thread tells the story's agent when job's
// rebase of branch onto base could not be made, or its close-out after it
// stopped or failed: the step, its findings, at most resolveFindings with
// more counted, where the job's output is, and how the agent fixes it,
// closes out, and comes back to review.
func (r *resolveJob) stepText(job *closeout.Job, branch, base, log, back string, findings []verify.Finding, more int) string {
	var b strings.Builder
	story := r.story
	if !r.rebased {
		fmt.Fprintf(&b, "flai serve's rebase and resolve job %s%s could not rebase `%s` onto `%s`: `flai stream sync %s` stopped on no path in conflict, as when the worktree has uncommitted changes or git failed, and left no rebase waiting, so the close-out did not run.\n\n", job.ID, askedBy(job), branch, base, story)
	} else {
		fmt.Fprintf(&b, "flai serve's rebase and resolve job %s%s rebased `%s` onto `%s` cleanly, and the close-out after it %s.\n\n", job.ID, askedBy(job), branch, base, stoppedWhere(job.StoppedAt, -1))
	}
	if len(findings) > 0 {
		b.WriteString("Its findings, from the story's last verify report:\n\n")
		for _, f := range findings {
			fmt.Fprintf(&b, "- %s\n", strings.ReplaceAll(strings.TrimRight(f.Text(), "\n"), "\n", "\n  "))
		}
		if more > 0 {
			fmt.Fprintf(&b, "- and %d more\n", more)
		}
		b.WriteString("\n")
	}
	fmt.Fprintf(&b, "Its whole output is in `%s`.\n\n", log)
	fmt.Fprintf(&b, "flai moved %s back to %s. "+
		"This thread is for %s's agent, which fixes it; the operator does not. "+
		"In its worktree, run `flai stream sync %s`, fix what stopped it, close out again, with the MCP tool `verify` or `flai verify %s --job`, and move %s to review again. "+
		"flai serve resolves this thread once %s is back in review.",
		story, back, story, story, story, story, story)
	return b.String()
}

// stepFindings are the findings rep, the story's verify report, holds for
// step, at most resolveFindings, and how many more there are; none when rep
// is not the job's own run, newest, or names no such step.
func stepFindings(rep verify.Report, newest bool, step string) ([]verify.Finding, int) {
	if !newest || step == "" {
		return nil, 0
	}
	for _, s := range rep.Steps {
		if s.Name != step {
			continue
		}
		if len(s.Findings) <= resolveFindings {
			return s.Findings, s.Omitted
		}
		return s.Findings[:resolveFindings], len(s.Findings) - resolveFindings + s.Omitted
	}
	return nil, 0
}

// stoppedWhere says where a close-out ended that did not pass: "stopped at
// <step>", with n findings when n is not negative, or "failed at no step its
// verify report names" for step "".
func stoppedWhere(step string, n int) string {
	if step == "" {
		return "failed at no step its verify report names"
	}
	s := "stopped at " + step
	switch {
	case n == 1:
		s += " with 1 finding"
	case n > 1:
		s += fmt.Sprintf(" with %d findings", n)
	}
	return s
}
