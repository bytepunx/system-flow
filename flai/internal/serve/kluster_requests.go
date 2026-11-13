package serve

import (
	"cmp"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/execx"
	"github.com/bytepunx/system-flow/flai/internal/inbox"
	"github.com/bytepunx/system-flow/flai/internal/kluster"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// A story's agent asking flai serve's kluster pass to deploy its own story's
// namespace, run the project's checks against it, or tear it down, while the
// story is in progress (S-0660, S-0661, ADR-0292 §1 to §5). The agent's MCP
// tool, or flai kluster deploy --job or flai checks run --job, writes a
// request under .flai-cache/kluster/requests (package kluster, which both
// sides import); every second, beside the minute's looks, the pass takes each
// and answers it at once: refused, saying why and what to do, with nothing
// queued; with the job of the story that waits or runs already; from the
// namespace's record, for a story in review deployed at its head; or by
// queuing a job in the pass's one queue, answering its id and the jobs ahead
// of it. The job's state is written when it is queued, when it starts, and
// when it ends, and its end appends a kluster event to the project's event
// log, which the agent waits for. A job left waiting or running by an earlier
// flai serve ends interrupted, with its event, when the pass starts again.

// klusterRequestEvery is how often the pass looks for a story's request.
const klusterRequestEvery = time.Second

// klusterRequestRoot is the folder whose .flai-cache holds the requests: the
// served root, the project's main checkout, as for the close-out's.
func (p *klusterPass) klusterRequestRoot() string { return p.e.Root }

// takeRequests ends, at its first, the asked-for jobs an earlier flai serve
// left, then takes each request waiting and answers it; a request withdrawn
// first is left. It wakes the worker when it queued a job.
func (p *klusterPass) takeRequests(ctx context.Context) {
	root := p.klusterRequestRoot()
	if !p.swept {
		p.swept = p.sweepAsked()
	}
	pending, err := kluster.Pending(root)
	if err != nil {
		p.o.Logger.Warn("kluster requests not taken", "component", "serve", "root", p.e.Root, "err", err.Error())
		return
	}
	if len(pending) == 0 {
		return
	}
	repo, err := workitem.Open(p.e.Root)
	if err != nil {
		// left where they are: the asker withdraws them and says flai serve is not taking them
		p.o.Logger.Warn("kluster requests not taken", "component", "serve", "root", p.e.Root, "err", err.Error())
		return
	}
	woke := false
	for _, file := range pending {
		if ctx.Err() != nil {
			return
		}
		req, ok, err := kluster.Take(root, file.Story, file.Kind)
		if err == nil && !ok {
			continue // withdrawn
		}
		var a kluster.Answer
		queued := false
		if err != nil {
			a = p.refuse(file, err.Error())
		} else {
			a, queued = p.answerRequest(repo, file, req)
		}
		if err := kluster.WriteAnswer(root, a); err != nil {
			p.o.Logger.Warn("kluster answer not written", "component", "serve", "root", p.e.Root, "story", file.Story, "kind", file.Kind, "err", err.Error())
		}
		woke = woke || queued
	}
	if woke {
		p.signal()
	}
}

// refuse is the answer to the request in file, which queues nothing, for why.
func (p *klusterPass) refuse(file kluster.Request, why string) kluster.Answer {
	p.o.Logger.Info("kluster request refused", "component", "serve", "root", p.e.Root, "story", file.Story, "kind", file.Kind, "by", file.By, "reason", why)
	return kluster.Answer{Story: file.Story, Kind: file.Kind, Reason: why}
}

// answerRequest answers req, taken from the file named by file's story and
// kind, and reports whether it queued a job.
func (p *klusterPass) answerRequest(repo *workitem.Repo, file, req kluster.Request) (kluster.Answer, bool) {
	refuse := func(format string, args ...any) (kluster.Answer, bool) {
		return p.refuse(file, fmt.Sprintf(format, args...)), false
	}
	if workitem.CanonicalID(req.Story) != file.Story || req.Kind != file.Kind {
		return refuse("the request at %s names %s's %s, not %s's %s, so nothing is queued: ask again with the MCP tool or flai kluster %s %s --job", kluster.RequestFile(p.klusterRequestRoot(), file.Story, file.Kind, kluster.TakenSuffix), req.Story, req.Kind, file.Story, file.Kind, file.Kind, file.Story)
	}
	story, kind, m := file.Story, file.Kind, repo.Manifest
	if req.By == "" {
		return refuse("the request for %s's %s names no agent that asks, so nothing is queued: ask again from the story's session, whose FLAI_AGENT names you", story, kind)
	}
	switch {
	case m.Kluster == nil:
		return refuse("the manifest has no kluster: block, so flai serve runs no kluster job for %s: turning the kluster host action on with flai serve enable kluster, which writes one, is the operator's. Say in %s's Notes which criteria you could not verify against a running system, and why", story, story)
	case !p.on(m):
		return refuse("the kluster host action is off for this project, so flai serve runs no kluster job for %s: turning it on with flai serve enable kluster is the operator's. Say in %s's Notes which criteria you could not verify against a running system, and why", story, story)
	case kind != kluster.KindDeploy && kind != kluster.KindChecks && kind != kluster.KindTeardown:
		return refuse("flai serve does not run a kluster %s job for a story's agent; ask for %s, %s, or %s", kind, kluster.KindDeploy, kluster.KindChecks, kluster.KindTeardown)
	case kind != kluster.KindTeardown && !klusterDeployable(m):
		return refuse("deploys are not set up yet for this project: its kluster: block names no deploy command, so flai serve deploys no story and runs no checks against a namespace. Setting kluster.deploy is the operator's; say in %s's Notes which criteria you could not verify against a running system, and why", story)
	}
	if kind == kluster.KindChecks {
		if _, why := p.checksConfig(repo); why != "" {
			return refuse("%s, so flai serve runs no checks against %s's namespace. Say in %s's Notes which criteria you could not verify against a running system, and why", why, story, story)
		}
	}
	it, err := repo.Get(story)
	if err != nil || it.Type != workitem.Story {
		return refuse("%s is not a story of this project that flai serve could read, so nothing is queued: ask for your own story's namespace", story)
	}
	root := mainCheckout(repo)
	ns, _ := kluster.StoryNamespace(story)
	record, recorded, err := kluster.ReadRecord(root, ns)
	if err != nil {
		return refuse("flai serve could not read %s's namespace's record (%v), so nothing is queued: ask again, and say so on a thread on %s if it goes on", story, err, story)
	}
	switch {
	case it.Archived || (it.Status != workitem.InProgress && it.Status != workitem.Review):
		return refuse("%s is %s, and a story's agent asks for its namespace only while the story is in progress: move it to in-progress first", story, it.Status)
	case it.Status == workitem.Review && kind == kluster.KindTeardown:
		return refuse("%s is in review, and its namespace is the review's evidence for its acceptance: flai serve tears it down once the story is accepted, cancelled, or moved back to ready or backlog", story)
	case it.Status == workitem.Review:
		return p.answerReview(repo, file, record, recorded), false
	case kind == kluster.KindDeploy && m.KlusterAgentNamespaces() == 0:
		return refuse("deploys before review are off for this project (kluster.agent_namespaces is 0), so flai serve deploys %s only once it is in review. Say in its Notes which criteria you could not verify against a running system, and why: the review's deploy checks them", story)
	}
	worktree := repo.WorktreePath(story)
	if !isDir(worktree) {
		return refuse("%s has no worktree at %s for flai serve to build in, so nothing is queued: open it with flai stream open %s, and ask again", story, worktree, story)
	}
	head := ""
	if kind != kluster.KindTeardown {
		if head, err = kluster.StoryHead(gitOf(p.o), root, story); err != nil {
			return refuse("flai serve could not resolve %s's head for its %s (%v), so nothing is queued", story, kind, err)
		}
	}
	// one job per story: the one that runs or waits is answered, but for a
	// deploy at a newer head, which takes the waiting one's place, as an
	// asked-for job takes a look's that would not run for a story in progress
	waiting, ahead, running, found := p.jobOf(story)
	if found && (kind == kluster.KindChecks || running || waiting.kind != kind || (waiting.asked() && head == waiting.head)) {
		return p.found(file, waiting, ahead), false
	}
	dirty, err := dirtyFiles(gitOf(p.o), worktree)
	switch {
	case err != nil:
		return refuse("flai serve could not read whether %s's worktree %s has uncommitted changes (%v), so nothing is queued: ask again", story, worktree, err)
	case len(dirty) > 0:
		return refuse("%s's worktree has uncommitted changes to tracked files, and a kluster job builds the committed head only: %s. Commit them, then ask again", story, namedFiles(dirty))
	case kind == kluster.KindTeardown && !recorded:
		return refuse("%s's namespace %s has no record, so there is nothing to tear down", story, ns)
	}
	if kind == kluster.KindChecks {
		// not counted against kluster.agent_namespaces: it needs a record
		// that passed, which a namespace held already has
		_, env, why, err := klusterEnvWhy(repo, story, gitOf(p.o), p.o.Host.Enabled)
		switch {
		case err != nil:
			return refuse("flai serve could not read whether %s's namespace is deployed at its head (%v), so nothing is queued: ask again, and say so on a thread on %s if it goes on", story, err, story)
		case env == nil:
			return refuse("%s, so %s's checks have nothing deployed at its head to run against and nothing is queued: %s", why, story, deployFirst(story))
		}
	}
	if kind == kluster.KindDeploy && !recorded {
		limit := m.KlusterAgentNamespaces()
		if holders := p.holders(repo, root, story); len(holders) >= limit {
			return refuse("%s has no namespace yet, and %s in progress hold one already, the most kluster.agent_namespaces allows (%d): %s. A namespace frees only when its story is accepted, cancelled, moved back to ready or backlog, or torn down, so nothing is queued: go on with the work that needs no namespace and ask again later, or say in %s's Notes which criteria you could not verify against a running system and why; the review's deploy checks them",
				story, countStories(len(holders)), limit, strings.Join(holders, "; "), story)
		}
	}
	j := klusterJob{kind: kind, story: story, head: head, by: req.By}
	j, ahead = p.queueAsked(j)
	if err := kluster.WriteAskedJob(p.klusterRequestRoot(), askedState(j)); err != nil {
		p.o.Logger.Warn("kluster job's state not written", "component", "serve", "root", p.e.Root, "story", story, "job_id", j.id, "err", err.Error())
	}
	p.o.Logger.Info("kluster job queued for a story's agent", "component", "serve", "root", p.e.Root, "story", story, "kind", kind, "job_id", j.id, "by", j.by, "commit", head, "ahead", ahead)
	return kluster.Answer{Job: j.id, Kind: kind, Story: story, Ahead: ahead, Started: true, Log: p.logOf(kind, root, ns)}, true
}

// logOf is where a job of kind for namespace ns writes its output: the
// folder of the checks runs' logs for a checks run, the namespace's log
// folder in root for the rest.
func (p *klusterPass) logOf(kind, root, ns string) string {
	if kind == kluster.KindChecks {
		return p.o.Dir.checksDir()
	}
	return kluster.LogDir(root, ns)
}

// answerReview answers the deploy or checks request in file for a story in
// review: with the job of the story that waits or runs, from record when its
// deploy passed at its branch's head, its checks too while they are on, or
// saying that the review deploys it; it queues nothing (ADR-0292 §1). A
// checks request answered from the record has its checks run's commit,
// outcome, and log.
func (p *klusterPass) answerReview(repo *workitem.Repo, file kluster.Request, record kluster.Record, recorded bool) kluster.Answer {
	story := file.Story
	if j, ahead, _, found := p.jobOf(story); found {
		return p.found(file, j, ahead)
	}
	root := mainCheckout(repo)
	ns, _ := kluster.StoryNamespace(story)
	head, _ := kluster.StoryHead(gitOf(p.o), root, story)
	checks := file.Kind == kluster.KindChecks
	if recorded && head != "" && passedAt(record, head, checks || p.checksOn(repo)) {
		p.o.Logger.Info("kluster request answered from the record", "component", "serve", "root", p.e.Root, "story", story, "kind", file.Kind, "by", file.By, "commit", record.Commit)
		a := kluster.Answer{Kind: file.Kind, Story: story, Recorded: true, Commit: record.Commit, Outcome: record.Outcome, Log: kluster.LogDir(root, ns)}
		if checks {
			a.Commit, a.Outcome, a.Log = record.ChecksCommit, record.ChecksOutcome, record.ChecksLog
		}
		return a
	}
	return p.refuse(file, fmt.Sprintf("%s is in review, where flai serve deploys it itself at its branch's head on its next look, and again at each new head, and runs its checks once the deploy passes, so a request for it queues nothing: flai kluster deployments shows its namespace's record, and a deploy or checks run that fails there sends it back to in-progress with a thread", story))
}

// found answers the request in file with j, the job of its story that waits
// with ahead jobs before it, or runs.
func (p *klusterPass) found(file kluster.Request, j klusterJob, ahead int) kluster.Answer {
	p.o.Logger.Info("kluster request answered with the job of its story", "component", "serve", "root", p.e.Root, "story", file.Story, "kind", file.Kind, "by", file.By, "job_id", j.id, "job_kind", j.kind, "ahead", ahead)
	root := p.klusterRequestRoot()
	ns, _ := kluster.StoryNamespace(file.Story)
	return kluster.Answer{Job: j.id, Kind: file.Kind, Story: file.Story, Ahead: ahead, Log: kluster.LogDir(root, ns)}
}

// jobOf is the job of story that runs, or else the first that waits, with
// the jobs ahead of it, and whether it runs; false when none does.
func (p *klusterPass) jobOf(story string) (j klusterJob, ahead int, running, found bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.running != nil && p.running.story == story {
		return *p.running, 0, true, true
	}
	for i, q := range p.queue {
		if q.story == story {
			return q, i + p.runningCount(), false, true
		}
	}
	return klusterJob{}, 0, false, false
}

// runningCount is 1 while a job runs; the caller holds mu.
func (p *klusterPass) runningCount() int {
	if p.running != nil {
		return 1
	}
	return 0
}

// queueAsked queues j, a job a story's agent asked for, with a new id: in the
// place of the job of the same kind of its story that waits, at an older
// head, or at the queue's end. It answers j and the jobs ahead of it.
func (p *klusterPass) queueAsked(j klusterJob) (klusterJob, int) {
	p.mu.Lock()
	defer p.mu.Unlock()
	j.queued = p.o.Now().UTC()
	j.id = kluster.JobID(j.story, j.kind, j.queued)
	for i, q := range p.queue {
		if q.story == j.story && q.kind == j.kind {
			p.queue[i] = j
			return j, i + p.runningCount()
		}
	}
	p.queue = append(p.queue, j)
	return j, len(p.queue) - 1 + p.runningCount()
}

// holders names each story in progress but story that holds a namespace, its
// record's commit and when its deploy started, or a deploy its agent asked
// for that waits or runs, which will make one (ADR-0292 §4). Stories in
// review never count.
func (p *klusterPass) holders(repo *workitem.Repo, root, story string) []string {
	items, err := repo.List(false)
	if err != nil {
		p.o.Logger.Warn("stories holding a kluster namespace not counted", "component", "serve", "root", p.e.Root, "err", err.Error())
		return nil
	}
	asked := p.askedDeploys()
	var out []string
	for _, it := range items {
		if it.Type != workitem.Story || it.Status != workitem.InProgress || it.Archived || it.ID == story {
			continue
		}
		ns, err := kluster.StoryNamespace(it.ID)
		if err != nil {
			continue
		}
		if r, ok, _ := kluster.ReadRecord(root, ns); ok {
			out = append(out, fmt.Sprintf("%s, deployed at %s, started %s", it.ID, short(r.Commit), r.StartedAt.UTC().Format(time.RFC3339)))
		} else if j, ok := asked[it.ID]; ok {
			out = append(out, fmt.Sprintf("%s, its deploy at %s asked for by %s and queued %s", it.ID, short(j.head), j.by, j.queued.Format(time.RFC3339)))
		}
	}
	return out
}

// askedDeploys are the deploys stories' agents asked for that wait or run,
// by story.
func (p *klusterPass) askedDeploys() map[string]klusterJob {
	p.mu.Lock()
	defer p.mu.Unlock()
	out := map[string]klusterJob{}
	jobs := p.queue
	if p.running != nil {
		jobs = append([]klusterJob{*p.running}, jobs...)
	}
	for _, j := range jobs {
		if j.asked() && j.kind == klusterDeploy {
			out[j.story] = j
		}
	}
	return out
}

// dirtyFiles are the tracked files with uncommitted changes in worktree, as
// git status --porcelain names them.
func dirtyFiles(git execx.Runner, worktree string) ([]string, error) {
	out, err := git.Run(worktree, "git", "status", "--porcelain", "--untracked-files=no")
	if err != nil {
		return nil, err
	}
	var files []string
	for line := range strings.Lines(out) {
		line = strings.TrimRight(line, "\r\n")
		if len(line) > 3 {
			files = append(files, strings.TrimSpace(line[3:]))
		}
	}
	return files, nil
}

// namedFiles names files, at most ten of them, saying how many more there
// are.
func namedFiles(files []string) string {
	const most = 10
	if len(files) <= most {
		return strings.Join(files, ", ")
	}
	return fmt.Sprintf("%s, and %d more", strings.Join(files[:most], ", "), len(files)-most)
}

// countStories is "1 story" or "n stories".
func countStories(n int) string {
	if n == 1 {
		return "1 story"
	}
	return fmt.Sprintf("%d stories", n)
}

// isDir reports whether path names a folder.
func isDir(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

// askedState is the state of j, a job a story's agent asked for, as it is
// queued.
func askedState(j klusterJob) kluster.AskedJob {
	ns, _ := kluster.StoryNamespace(j.story)
	return kluster.AskedJob{ID: j.id, Kind: j.kind, Story: j.story, By: j.by, Namespace: ns, QueuedAt: j.queued}
}

// startAsked writes that j, when a story's agent asked for it, has started.
func (p *klusterPass) startAsked(j klusterJob) {
	if !j.asked() {
		return
	}
	s := askedState(j)
	s.StartedAt = p.o.Now().UTC()
	if err := kluster.WriteAskedJob(p.klusterRequestRoot(), s); err != nil {
		p.o.Logger.Warn("kluster job's start not written", "component", "serve", "root", p.e.Root, "story", j.story, "job_id", j.id, "err", err.Error())
	}
}

// endAsked writes how j, a job a story's agent asked for and the one
// running, ended, marks it no longer running, and then appends its kluster
// event to repo's event log, when repo is not nil: so that a request its
// reader makes on reading it queues a new job rather than being answered this
// one. A job that ended with no outcome, as flai serve stops under it, is left
// for the next start to end interrupted.
func (p *klusterPass) endAsked(repo *workitem.Repo, j klusterJob, end klusterEnd) {
	defer p.done()
	if end.outcome == "" {
		return
	}
	root := p.klusterRequestRoot()
	s, ok, err := kluster.ReadAskedJob(root, j.story)
	if err != nil || !ok || s.ID != j.id {
		s = askedState(j)
	}
	s.EndedAt, s.Commit, s.Outcome = p.o.Now().UTC(), cmp.Or(end.commit, j.head), end.outcome
	s.FailedStep, s.Error, s.Steps, s.Log, s.ExitCode, s.Checks = end.fail.step, end.fail.err, end.steps, end.log, nil, end.checks
	if end.fail.step != "" {
		code := end.fail.exit
		s.ExitCode = &code
	}
	if err := kluster.WriteAskedJob(root, s); err != nil {
		p.o.Logger.Warn("kluster job's end not written", "component", "serve", "root", p.e.Root, "story", j.story, "job_id", j.id, "err", err.Error())
	}
	if end.outcome == kluster.OutcomeRefused && s.StartedAt.IsZero() {
		p.journalJob(j, "failed", fmt.Sprintf("%s of %s ran nothing: %s", j.kind, j.story, end.fail.err))
	}
	p.o.Logger.Info("kluster job asked for by a story's agent ended", "component", "serve", "root", p.e.Root, "story", j.story, "job_id", j.id, "kind", j.kind, "by", j.by, "outcome", s.Outcome, "commit", s.Commit, "failed_step", s.FailedStep, "log", s.Log)
	p.done()
	if repo != nil {
		p.announceKluster(repo, s)
	}
}

// sweepAsked ends interrupted each asked-for job an earlier flai serve left
// waiting or running, with its event, and reports whether it is done: false
// when there were some and the project could not be opened, to try again at
// the next look.
func (p *klusterPass) sweepAsked() bool {
	root := p.klusterRequestRoot()
	paths, _ := filepath.Glob(filepath.Join(kluster.RequestsDir(root), "*"+kluster.JobSuffix))
	var left []kluster.AskedJob
	for _, path := range paths {
		story := strings.TrimSuffix(filepath.Base(path), kluster.JobSuffix)
		j, ok, err := kluster.ReadAskedJob(root, story)
		if err != nil {
			p.o.Logger.Warn("kluster job not read", "component", "serve", "root", p.e.Root, "path", path, "err", err.Error())
			continue
		}
		if ok && !j.Ended() {
			left = append(left, j)
		}
	}
	if len(left) == 0 {
		return true
	}
	repo, err := workitem.Open(p.e.Root)
	if err != nil {
		p.o.Logger.Warn("interrupted kluster jobs not ended", "component", "serve", "root", p.e.Root, "jobs", len(left), "err", err.Error())
		return false
	}
	now := p.o.Now().UTC()
	for _, s := range left {
		s.EndedAt, s.Outcome = now, kluster.OutcomeInterrupted
		s.Error = "flai serve stopped while the job waited or ran; ask for it again"
		if err := kluster.WriteAskedJob(root, s); err != nil {
			p.o.Logger.Warn("kluster job's end not written", "component", "serve", "root", p.e.Root, "story", s.Story, "job_id", s.ID, "err", err.Error())
		}
		p.journalBy(s.By, "failed", fmt.Sprintf("%s of %s job %s: %s, flai serve stopped while it waited or ran, asked by %s", s.Kind, s.Story, s.ID, kluster.OutcomeInterrupted, s.By))
		p.o.Logger.Warn("kluster job interrupted", "component", "serve", "root", p.e.Root, "story", s.Story, "job_id", s.ID, "kind", s.Kind, "by", s.By)
		p.announceKluster(repo, s)
	}
	return true
}

// announceKluster appends the kluster event of s, a job that ended, to the
// project's event log; a failure loses the news and is logged.
func (p *klusterPass) announceKluster(repo *workitem.Repo, s kluster.AskedJob) {
	title := ""
	if it, err := repo.Get(s.Story); err == nil {
		title = it.Title
	}
	ev := inbox.KlusterEvent(s.Story, title, forgeAuthor, inbox.KlusterOutcome{Job: s.ID, Kind: s.Kind, Commit: s.Commit, Outcome: s.Outcome,
		FailedStep: s.FailedStep, ExitCode: s.ExitCode, Log: s.Log})
	if err := inbox.AppendEvent(repo, ev); err != nil {
		p.o.Logger.Warn("kluster event not logged", "component", "serve", "root", p.e.Root, "story", s.Story, "job_id", s.ID, "err", err.Error())
	}
}
