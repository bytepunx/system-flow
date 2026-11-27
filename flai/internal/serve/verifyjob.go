package serve

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/atomicfile"
	"github.com/bytepunx/system-flow/flai/internal/buildinfo"
	"github.com/bytepunx/system-flow/flai/internal/closeout"
	"github.com/bytepunx/system-flow/flai/internal/config"
	"github.com/bytepunx/system-flow/flai/internal/host"
	"github.com/bytepunx/system-flow/flai/internal/hostapi"
	"github.com/bytepunx/system-flow/flai/internal/inbox"
	"github.com/bytepunx/system-flow/flai/internal/verify"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// The close-out as a job of flai serve (S-0546, ADR-0229). Claude Code ends
// a tool call at 600 s, and a close-out often runs longer, so a story's agent
// asks flai serve for it instead: the MCP tool verify, or flai verify --job,
// writes a request under .flai-cache/verify/jobs (package closeout, which
// both sides import), which flai serve's verify pass takes by renaming it, as the merge pass takes a merge request, and
// answers at once with the job's id, or with the verdict the verify record
// holds for the commit. The job runs the project's close-out in the story's
// worktree with no model; its end is journaled and announced by a verified
// event in the project's event log, which the agent waits for. While flai
// serve's upgrade waits for the jobs to end, a request that would start one
// is held until the drain lifts, after the restart (ADR-0320).

const (
	// verifyEvery is how often the verify pass looks for a request.
	verifyEvery = time.Second
	// verifyJobLimit is the longest a close-out job runs before it is stopped.
	verifyJobLimit = 2 * time.Hour
	// closeOutScript is the project's close-out, in the story's worktree.
	closeOutScript = "scripts/close-out.sh"
	// The journal's action, and its method for a close-out job and for a
	// rebase and resolve job.
	closeOutAction = "close-out"
	closeOutMethod = "verify.job"
	resolveMethod  = "resolve.job"
	// fromTheRecord begins the journal detail of a close-out answered from
	// the verify record, which flai stats counts apart (metrics.md).
	fromTheRecord = "from the record"
	// verifyServeAuthor is who tells of a close-out's end in the event log.
	verifyServeAuthor = "flai serve"
)

// closeOutRunner runs argv in dir with env, its output to out, until it ends
// or ctx does, and answers its exit code; an error says it could not start.
type closeOutRunner func(ctx context.Context, dir string, env []string, out io.Writer, argv []string) (int, error)

// runCloseOut is the closeOutRunner that runs the process in its own process
// group, so that ctx ending stops it and everything it started.
func runCloseOut(ctx context.Context, dir string, env []string, out io.Writer, argv []string) (int, error) {
	cmd := exec.Command(argv[0], argv[1:]...)
	cmd.Dir, cmd.Env, cmd.Stdout, cmd.Stderr = dir, env, out, out
	Detach(cmd)
	if err := cmd.Start(); err != nil {
		return -1, err
	}
	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()
	select {
	case err := <-done:
		return *exitCodeOf(err), nil
	case <-ctx.Done():
		terminateThenKill(cmd.Process.Pid)
		return *exitCodeOf(<-done), nil
	}
}

// The variables that ask the worktree's close-out script for --fresh and
// --all (S-0671).
const (
	closeOutFreshEnv = "CLOSE_OUT_FRESH"
	closeOutAllEnv   = "CLOSE_OUT_ALL"
)

// closeOutEnv is env as a close-out or resolve job has it: without the host's
// address and token, the variables that mark an agent's run or session, the
// forges' default token variables, as an agent's, and the close-out script's
// own fresh and all, which only the request sets, and with FLAI_AGENT naming
// the agent that asked and closeout.ProgressEnv set to 1, so that flai verify
// writes its steps and tiers to the job's log as each starts and ends.
func closeOutEnv(env []string, by string) []string {
	drop := map[string]bool{"FLAI_AGENT": true, "FLAI_SESSION": true, host.URLEnv: true, host.TokenEnv: true, closeOutFreshEnv: true, closeOutAllEnv: true, closeout.ProgressEnv: true}
	for _, name := range runMarks {
		drop[name] = true
	}
	for _, f := range config.KnownForges() {
		drop[config.DefaultForgeTokenEnv(f)] = true
	}
	out := slices.DeleteFunc(slices.Clone(env), func(kv string) bool {
		name, _, _ := strings.Cut(kv, "=")
		return drop[name]
	})
	return append(out, "FLAI_AGENT="+by, closeout.ProgressEnv+"=1")
}

// verifyPass takes the close-out requests of one served project every second
// and answers each: from the verify record, with the id of the story's job
// that runs, by starting one, or, while an upgrade waits on the drain, by
// holding it. Jobs of different stories run side by side.
type verifyPass struct {
	o  Options
	e  Entry
	wg *sync.WaitGroup
	// drain counts the pass's jobs for flai serve's upgrader and says
	// whether one may start (ADR-0320); nil holds nothing.
	drain *jobDrain
	every time.Duration
	limit time.Duration
	run   closeOutRunner
	// exe is this flai's own executable, which runs flai verify where the
	// worktree has no close-out script, and the host's steps of a container
	// story's close-out (ADR-0274).
	exe func() (string, error)
	// recorded is the verdict the verify record holds for the story's
	// worktree as it stands, and true, or false when it holds none; nil
	// holds none. keepAt keeps a passing run's report under the commit the
	// close-out made after it; nil keeps nothing.
	recorded func(repo *workitem.Repo, story, worktree string, all bool) (verify.Report, bool, error)
	keepAt   func(repo *workitem.Repo, commit string, rep verify.Report) error

	// swept is whether the jobs an earlier flai serve left running were
	// ended; only the pass's own loop reads and writes it.
	swept bool
	// stateEvery is how often the pass looks for a story that entered a
	// state whose runs name something, and statesAt when it last did; only
	// the pass's own loop reads and writes statesAt.
	stateEvery time.Duration
	statesAt   time.Time

	mu sync.Mutex
	// running is counted in drain: track adds to it and forget and
	// forgetState delete from it.
	running map[string]string // story → id of its job that runs
	states  map[string]string // story → the state its state job that runs is for
}

func newVerifyPass(o Options, e Entry, wg *sync.WaitGroup, drain *jobDrain) *verifyPass {
	return &verifyPass{o: o, e: e, wg: wg, drain: drain, every: verifyEvery, limit: verifyJobLimit, run: runCloseOut, exe: os.Executable, stateEvery: stateEvery,
		running: map[string]string{}, states: map[string]string{}, recorded: recordedFor(o, e), keepAt: verify.KeepAt}
}

// recordedFor asks the verify record of project e as flai verify would at
// the worktree: with the tiers its branch's manifest declares, of the suites
// the moves out of the story's state need, leaving those marked queue to the
// merge queue while it batches, unless the request asks for all, which is
// every suite with none left to the queue (ADR-0229 §5, ADR-0312 §2). The
// record answers only when it holds a passing verdict at the head for every
// suite of those that the branch selects tiers of.
func recordedFor(o Options, e Entry) func(*workitem.Repo, string, string, bool) (verify.Report, bool, error) {
	return func(repo *workitem.Repo, story, worktree string, all bool) (verify.Report, bool, error) {
		git := gitOf(o)
		tiers, err := verify.CheckoutTiersFor(worktree, git, buildinfo.Version)
		if err != nil {
			return verify.Report{}, false, err
		}
		opts := verify.StoryOptions{Story: story, Project: repo, Worktree: worktree, Tiers: tiers, Git: git}
		if !all {
			opts.LeaveToQueue = queueBatches(o, e, repo.Manifest)
			if opts.Suites, err = verify.StorySuites(repo, story); err != nil {
				return verify.Report{}, false, err
			}
		}
		return verify.Recorded(opts)
	}
}

// start looks now and then every second in the background until ctx ends;
// a job running is stopped with ctx and waited for through the pass's group.
func (p *verifyPass) start(ctx context.Context) {
	p.wg.Add(1)
	go func() {
		defer p.wg.Done()
		for {
			p.look(ctx)
			select {
			case <-ctx.Done():
				return
			case <-time.After(p.every):
			}
		}
	}()
}

// dir is the project's verify jobs folder: the served root is the main
// checkout, whose .flai-cache is the project's cache.
func (p *verifyPass) dir() string { return closeout.Dir(filepath.Join(p.e.Root, ".flai-cache")) }

// look ends, at its first, the jobs an earlier flai serve left running,
// starts, every stateEvery once they are, the state jobs of the stories that
// entered a state whose runs name something (lookStates), puts back the
// requests held while no upgrade waits on the drain, then takes each request
// waiting and answers it; a request another look took first is left to it.
func (p *verifyPass) look(ctx context.Context) {
	dir := p.dir()
	if !p.swept {
		p.swept = p.sweep(dir)
	}
	if now := p.o.Now(); p.swept && now.Sub(p.statesAt) >= p.stateEvery {
		p.statesAt = now
		p.lookStates(ctx)
	}
	if _, holding := p.drain.holding(); !holding {
		p.putBack(dir)
	}
	requests, _ := filepath.Glob(filepath.Join(dir, "*"+closeout.RequestSuffix))
	if len(requests) == 0 {
		return
	}
	repo, err := workitem.Open(p.e.Root)
	if err != nil {
		p.o.Logger.Warn("close-out requests not taken", "component", "serve", "root", p.e.Root, "err", err.Error())
		return
	}
	for _, request := range requests {
		if ctx.Err() != nil {
			return
		}
		story := strings.TrimSuffix(filepath.Base(request), closeout.RequestSuffix)
		taken := closeout.File(dir, story, closeout.TakenSuffix)
		if err := os.Rename(request, taken); err != nil {
			continue
		}
		a := p.answer(ctx, repo, dir, story, taken)
		if err := closeout.WriteAnswer(dir, a); err != nil {
			p.o.Logger.Warn("close-out answer not written", "component", "serve", "root", p.e.Root, "story", story, "err", err.Error())
		}
	}
}

// answer answers the request at taken for story's close-out, starting its
// job when no job of the story runs and the record holds no verdict, or for
// its rebase and resolve, starting that job when the story is in review and
// no job of it runs, never from the record (S-0804); a job it would start
// while an upgrade waits on the drain is held instead (ADR-0320).
func (p *verifyPass) answer(ctx context.Context, repo *workitem.Repo, dir, story, taken string) closeout.Answer {
	refused := func(why string) closeout.Answer {
		p.o.Logger.Warn("close-out request refused", "component", "serve", "root", p.e.Root, "story", story, "reason", why)
		return closeout.Answer{Story: story, Reason: why}
	}
	var req closeout.Request
	data, err := os.ReadFile(taken)
	if err == nil {
		err = json.Unmarshal(data, &req)
	}
	if err != nil {
		return refused(fmt.Sprintf("flai serve could not read the request for %s's close-out at %s (%v), so nothing runs: ask again", story, taken, err))
	}
	if workitem.CanonicalID(req.Story) != story {
		return refused(fmt.Sprintf("the request at %s names %s, not %s, so nothing runs: ask again with the story's ID zero padded", taken, req.Story, story))
	}
	switch req.Kind {
	case "":
	case closeout.KindResolve:
		if why := resolveRefusal(repo, story); why != "" {
			return refused(why)
		}
		if req.Message == "" {
			req.Message = resolveMessage(story)
		}
	default:
		return refused(fmt.Sprintf("the request at %s asks for a job of kind %q, which this flai serve does not run, so nothing runs: ask for a close-out or a rebase and resolve, or upgrade flai serve", taken, req.Kind))
	}
	worktree := repo.WorktreePath(story)
	if info, err := os.Stat(worktree); err != nil || !info.IsDir() {
		return refused(fmt.Sprintf("%s has no worktree at %s, so its close-out cannot run: open it with flai stream open %s, or start the story with flai story start %s, and ask again", story, worktree, story, story))
	}
	// A job of the story that runs is answered before the record is asked:
	// it may be committing in the worktree, and its end tells the verdict.
	p.mu.Lock()
	id, running := p.running[story]
	state := p.states[story]
	p.mu.Unlock()
	if running && state != "" {
		return refused(fmt.Sprintf("flai serve's job %s for %s's entry into %s runs, rebasing its branch and running what %s runs in its worktree, so no %s starts beside it: ask again once it ended; its output is in %s",
			id, story, state, state, jobName(req.Kind), closeout.File(dir, story, closeout.StateLogSuffix)))
	}
	if running {
		return closeout.Answer{Job: id, Story: story, Log: closeout.File(dir, story, closeout.LogSuffix)}
	}
	// A story whose agent runs in a container is closed out in it (ADR-0274):
	// the record is not asked through git files the container may have
	// changed to point elsewhere.
	ca := p.containerAgentOf(repo, story)
	if ca != nil {
		if err := gitPointers(p.e.Root, story, worktree); err != nil {
			return refused(err.Error())
		}
	}
	// a resolve job always rebases first, so the record cannot answer it
	if req.Kind == "" && !req.Fresh && p.recorded != nil {
		rep, ok, err := p.recorded(repo, story, worktree, req.All)
		switch {
		case err != nil:
			p.o.Logger.Warn("verify record not read", "component", "serve", "root", p.e.Root, "story", story, "err", err.Error())
		case ok:
			return p.fromRecord(repo, req, rep)
		}
	}
	if version, holding := p.drain.holding(); holding {
		return p.hold(dir, req, version, refused)
	}
	return p.startJob(ctx, dir, worktree, req, ca, refused)
}

// hold keeps req as its story's held request, for the look after the drain
// lifts to put back, and answers that it is held, saying until when; the
// upgrade waits for version, "" when it is not known.
func (p *verifyPass) hold(dir string, req closeout.Request, version string, refused func(string) closeout.Answer) closeout.Answer {
	story, name := req.Story, jobName(req.Kind)
	held := closeout.File(dir, story, closeout.HeldSuffix)
	data, err := json.MarshalIndent(req, "", "  ")
	if err == nil {
		err = atomicfile.WriteFile(held, append(data, '\n'), 0o644)
	}
	if err != nil {
		return refused(fmt.Sprintf("flai serve waits to upgrade flai until the jobs it runs end, so it starts no %s for %s now, and could not keep the request to start it after the restart (%v): ask again once flai serve has restarted", name, story, err))
	}
	release := "a newer flai release"
	if version != "" {
		release = "flai " + version
	}
	p.o.Logger.Info("close-out request held", "component", "serve", "root", p.e.Root, "story", story, "kind", req.Kind, "by", req.By, "newest", version)
	return closeout.Answer{Story: story, Held: true, Reason: fmt.Sprintf("flai serve upgrades to %s once the jobs it runs end, at most 45 minutes from now, so it holds %s's %s and starts it once it has restarted, ending with a verified event as any job does: hold wait_for_events until it reports verified for %s", release, story, name, story)}
}

// putBack renames each request held in dir back to its story's request, for
// this look to take, unless the story has a request waiting already, which
// is newer and wins: the held one is removed.
func (p *verifyPass) putBack(dir string) {
	held, _ := filepath.Glob(filepath.Join(dir, "*"+closeout.HeldSuffix))
	for _, path := range held {
		story := strings.TrimSuffix(filepath.Base(path), closeout.HeldSuffix)
		request := closeout.File(dir, story, closeout.RequestSuffix)
		// a link never writes over a request the story's agent wrote meanwhile
		err := os.Link(path, request)
		if err != nil && !errors.Is(err, fs.ErrExist) {
			if _, serr := os.Stat(request); errors.Is(serr, fs.ErrNotExist) {
				err = os.Rename(path, request)
			}
		}
		switch {
		case errors.Is(err, fs.ErrExist):
			p.o.Logger.Info("held close-out request dropped for a newer one", "component", "serve", "root", p.e.Root, "story", story)
		case err != nil:
			p.o.Logger.Warn("held close-out request not put back", "component", "serve", "root", p.e.Root, "story", story, "err", err.Error())
			continue
		default:
			p.o.Logger.Info("held close-out request put back", "component", "serve", "root", p.e.Root, "story", story)
		}
		if err := os.Remove(path); err != nil && !errors.Is(err, fs.ErrNotExist) {
			p.o.Logger.Warn("held close-out request not removed", "component", "serve", "root", p.e.Root, "story", story, "err", err.Error())
		}
	}
}

// fromRecord answers req with the verdict the record holds in rep, journals
// it as one end entry, and announces it: nothing runs.
func (p *verifyPass) fromRecord(repo *workitem.Repo, req closeout.Request, rep verify.Report) closeout.Answer {
	verdict, stoppedAt := closeout.VerdictFailed, rep.StoppedAt
	switch {
	case rep.Passed:
		verdict, stoppedAt = closeout.VerdictPassed, ""
	case stoppedAt != "":
		verdict = closeout.VerdictStopped
	}
	p.journal(hostapi.Entry{By: req.By, Outcome: closeOutOutcome(verdict), Verdict: verdict,
		Detail: fmt.Sprintf("%s: %s's close-out %s%s at %s, run %s", fromTheRecord, req.Story, verdict, atStep(stoppedAt), short(rep.Commit), rep.RanAt.UTC().Format(time.RFC3339))})
	p.announce(repo, req.Story, inbox.Verification{Commit: rep.Commit, Verdict: verdict, StoppedAt: stoppedAt, Recorded: true})
	p.o.Logger.Info("close-out answered from the record", "component", "serve", "root", p.e.Root, "story", req.Story, "verdict", verdict, "commit", rep.Commit)
	return closeout.Answer{Story: req.Story, Recorded: true, Verdict: verdict, StoppedAt: stoppedAt, Commit: rep.Commit, RanAt: rep.RanAt}
}

// jobBody is a job's run, its output to out, answering the close-out's exit
// code; an error says a step could not run.
type jobBody func(ctx context.Context, out io.Writer) (int, error)

// startJob starts req's close-out, or its rebase and resolve, in worktree as
// a job, records and journals its start, and answers its id. The close-out of
// a story whose agent runs in a container, ca, runs its verify there and its
// commit and sync on the host (ADR-0274); any other story's runs command's
// argv on the host. A rebase and resolve runs this flai's flai stream sync on
// the host before that close-out (resolveJob).
func (p *verifyPass) startJob(ctx context.Context, dir, worktree string, req closeout.Request, ca *containerAgent, refused func(string) closeout.Answer) closeout.Answer {
	story := req.Story
	started := p.o.Now().UTC()
	stamp := started.Format("20060102T150405Z")
	var body jobBody
	command := ""
	if ca != nil {
		c := &containerCloseOut{p: p, agent: *ca, story: story, worktree: worktree, req: req, id: stamp}
		body, command = c.run, c.String()
	} else {
		argv, extra, err := p.command(worktree, req)
		if err != nil {
			return refused(err.Error())
		}
		body = func(jctx context.Context, out io.Writer) (int, error) {
			return p.run(jctx, worktree, append(closeOutEnv(os.Environ(), req.By), extra...), out, argv)
		}
		command = strings.Join(append(slices.Clone(extra), argv...), " ")
	}
	var settle func(*workitem.Repo, *closeout.Job, verify.Report, bool)
	if req.Kind == closeout.KindResolve {
		r := newResolveJob(p, story, worktree, req.By, ca != nil, body)
		body, settle = r.run, r.settle
		command = "flai stream sync " + story + ", then " + command
	}
	job := closeout.Job{ID: story + "-" + stamp, Story: story, Kind: req.Kind, By: req.By, Message: req.Message, All: req.All, Fresh: req.Fresh, Started: started}
	if !p.track(story, job.ID, "") {
		// the upgrade began to wait on the drain since answer looked
		version, _ := p.drain.holding()
		return p.hold(dir, req, version, refused)
	}
	if err := closeout.WriteJob(dir, job); err != nil {
		p.forget(story)
		return refused(fmt.Sprintf("flai serve could not start %s's %s: %v", story, jobName(req.Kind), err))
	}
	p.journal(hostapi.Entry{At: started.Format(time.RFC3339), Method: jobMethod(job.Kind), By: req.By, Outcome: "started", Detail: story + " job " + job.ID})
	p.o.Logger.Info("close-out job started", "component", "serve", "root", p.e.Root, "story", story, "job_id", job.ID, "kind", job.Kind, "by", req.By, "command", command)
	p.wg.Add(1)
	go func() {
		defer p.wg.Done()
		p.runJob(ctx, dir, worktree, job, body, settle)
	}()
	return closeout.Answer{Job: job.ID, Story: story, Started: true, Log: closeout.File(dir, story, closeout.LogSuffix)}
}

// command is what req's close-out runs in worktree, and the variables it adds
// to the job's environment: the worktree's close-out script with the commit
// message, and CLOSE_OUT_FRESH and CLOSE_OUT_ALL set to 1 as asked, or, where
// it has none, this flai's verify with --all and --fresh as asked. The script
// is given no flag for them, since one that predates them would take a
// leading flag for the story; it ignores variables it does not read (S-0671).
func (p *verifyPass) command(worktree string, req closeout.Request) ([]string, []string, error) {
	if info, err := os.Stat(filepath.Join(worktree, closeOutScript)); err == nil && !info.IsDir() {
		argv := []string{"sh", closeOutScript, req.Story}
		if req.Message != "" {
			argv = append(argv, "-m", req.Message)
		}
		var env []string
		if req.Fresh {
			env = append(env, closeOutFreshEnv+"=1")
		}
		if req.All {
			env = append(env, closeOutAllEnv+"=1")
		}
		return argv, env, nil
	}
	flai, err := p.exe()
	if err != nil {
		return nil, nil, fmt.Errorf("%s's worktree has no %s and this flai's executable is not known, so its close-out cannot run: run flai verify %s --record-issues in the worktree yourself: %w", req.Story, closeOutScript, req.Story, err)
	}
	argv := []string{flai, "verify", req.Story, "--record-issues"}
	if req.All {
		argv = append(argv, "--all")
	}
	if req.Fresh {
		argv = append(argv, "--fresh")
	}
	return argv, nil, nil
}

// runJob writes the head job runs at in its state, runs its close-out, body,
// in worktree for at most the pass's limit, its output in its log, then ends
// it: settle, when not nil, sets how a resolve job ended and sends its story
// back when it must, from the project, the job with its verdict, and the
// story's verify report, newest when the job's own run wrote it. When ctx
// ends first, the process group is stopped and the job left running in its
// state, for the next flai serve to end interrupted, and nothing is settled.
func (p *verifyPass) runJob(ctx context.Context, dir, worktree string, job closeout.Job, body jobBody, settle func(*workitem.Repo, *closeout.Job, verify.Report, bool)) {
	before := p.head(worktree)
	// the head it runs at, so that a reader can tell a running job at a new
	// head from the report of one before it (I-0384)
	job.Head = before
	if err := closeout.WriteJob(dir, job); err != nil {
		p.o.Logger.Warn("close-out job's head not written", "component", "serve", "root", p.e.Root, "story", job.Story, "job_id", job.ID, "err", err.Error())
	}
	jctx, cancel := context.WithTimeout(ctx, p.limit)
	defer cancel()
	code := -1
	var runErr error
	logPath := closeout.File(dir, job.Story, closeout.LogSuffix)
	if out, err := os.Create(logPath); err != nil {
		runErr = fmt.Errorf("create its log %s: %w", logPath, err)
	} else {
		code, runErr = body(jctx, out)
		_ = out.Close()
	}
	if ctx.Err() != nil {
		p.forget(job.Story)
		p.o.Logger.Info("close-out job stopped with flai serve", "component", "serve", "root", p.e.Root, "story", job.Story, "job_id", job.ID)
		return
	}
	timedOut := errors.Is(jctx.Err(), context.DeadlineExceeded)
	ended := p.o.Now().UTC()
	job.Ended, job.ExitCode, job.DurationMS = ended, &code, ended.Sub(job.Started).Milliseconds()
	job.Commit = p.head(worktree)
	repo, err := workitem.Open(p.e.Root)
	if err != nil {
		p.o.Logger.Warn("close-out job's project not opened", "component", "serve", "root", p.e.Root, "story", job.Story, "job_id", job.ID, "err", err.Error())
	}
	var rep verify.Report
	newest := false
	if repo != nil {
		if r, ok, err := verify.LastReport(repo, job.Story); err != nil {
			p.o.Logger.Warn("verify report not read", "component", "serve", "root", p.e.Root, "story", job.Story, "job_id", job.ID, "err", err.Error())
		} else if ok && !r.RanAt.Before(job.Started.Truncate(time.Second)) {
			rep, newest = r, true // the job's own run, not one before it
		}
	}
	job.Verdict, job.StoppedAt = closeout.VerdictFailed, ""
	switch {
	case code == 0 && runErr == nil:
		job.Verdict = closeout.VerdictPassed
	case newest && !rep.Passed && rep.StoppedAt != "":
		job.Verdict, job.StoppedAt = closeout.VerdictStopped, rep.StoppedAt
	}
	if job.Verdict == closeout.VerdictPassed && newest && rep.Passed && job.Commit != "" && job.Commit != before && rep.Commit != job.Commit {
		p.keep(job, rep)
	}
	if settle != nil {
		settle(repo, &job, rep, newest)
	}
	if err := closeout.WriteJob(dir, job); err != nil {
		p.o.Logger.Warn("close-out job's end not written", "component", "serve", "root", p.e.Root, "story", job.Story, "job_id", job.ID, "err", err.Error())
	}
	// forgotten before the event, so that a request its reader makes on
	// reading it starts a new job rather than being answered this one's id
	p.forget(job.Story)
	duration := time.Duration(job.DurationMS) * time.Millisecond
	detail := fmt.Sprintf("%s job %s: %s%s in %s", job.Story, job.ID, job.Verdict, atStep(job.StoppedAt), duration.Round(time.Second))
	switch {
	case runErr != nil:
		detail += fmt.Sprintf(" (it could not run: %v)", runErr)
	case timedOut:
		detail += fmt.Sprintf(" (stopped at the %s limit)", p.limit)
	}
	if job.Outcome != "" {
		detail += "; " + job.Outcome
		if job.Thread != "" {
			detail += " on " + job.Thread
		}
	}
	p.journal(hostapi.Entry{At: ended.Format(time.RFC3339), Method: jobMethod(job.Kind), By: job.By, Outcome: closeOutOutcome(job.Verdict), Verdict: job.Verdict, DurationMS: job.DurationMS, Detail: detail})
	if repo != nil {
		p.announce(repo, job.Story, inbox.Verification{Job: job.ID, Commit: job.Commit, Verdict: job.Verdict, StoppedAt: job.StoppedAt, Duration: duration})
	}
	level, attrs := p.o.Logger.Info, []any{"component", "serve", "root", p.e.Root, "story", job.Story, "job_id", job.ID, "kind", job.Kind,
		"verdict", job.Verdict, "stopped_at", job.StoppedAt, "exit_code", code, "duration_ms", job.DurationMS, "timed_out", timedOut,
		"outcome", job.Outcome, "thread", job.Thread}
	if runErr != nil {
		level, attrs = p.o.Logger.Warn, append(attrs, "err", runErr.Error())
	}
	level("close-out job ended", attrs...)
}

// keep keeps rep, the passing run of job, under the commit the close-out
// made after it, which holds exactly the files it verified (ADR-0229 §5).
func (p *verifyPass) keep(job closeout.Job, rep verify.Report) {
	if p.keepAt == nil {
		return
	}
	repo, err := workitem.Open(p.e.Root)
	if err == nil {
		err = p.keepAt(repo, job.Commit, rep)
	}
	if err != nil {
		p.o.Logger.Warn("verify record not kept at the close-out's commit", "component", "serve", "root", p.e.Root, "story", job.Story, "job_id", job.ID, "commit", job.Commit, "err", err.Error())
	}
}

// track marks story's job id as running, for its entry into state when it
// is a state job, and counts it in the drain; it reports false and marks
// nothing while an upgrade waits on the drain.
func (p *verifyPass) track(story, id, state string) bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	if _, ok := p.running[story]; !ok && !p.drain.begin() {
		return false
	}
	p.running[story] = id
	if state != "" {
		p.states[story] = state
	}
	return true
}

// forget marks story's job as no longer running.
func (p *verifyPass) forget(story string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.untrack(story)
}

// untrack deletes story from running, counting its job's end in the drain;
// p.mu is held.
func (p *verifyPass) untrack(story string) {
	if _, ok := p.running[story]; ok {
		delete(p.running, story)
		p.drain.end()
	}
}

// sweep ends interrupted each job in dir that an earlier flai serve left
// running, its state job too, and reports whether it is done: false when
// there were some and the project could not be opened, to try again at the
// next look. Nothing is sent back; a close-out or resolve job's end is
// announced, and an interrupted state job's entry runs again (lookStates).
func (p *verifyPass) sweep(dir string) bool {
	paths, _ := filepath.Glob(filepath.Join(dir, "*"+closeout.JobSuffix))
	states, _ := filepath.Glob(filepath.Join(dir, "*"+closeout.StateSuffix))
	paths = append(paths, states...)
	var left []closeout.Job
	for _, path := range paths {
		j, ok, err := closeout.ReadJobAt(path)
		if err != nil {
			p.o.Logger.Warn("close-out job not read", "component", "serve", "root", p.e.Root, "path", path, "err", err.Error())
			continue
		}
		if ok && j.Running() {
			left = append(left, j)
		}
	}
	if len(left) == 0 {
		return true
	}
	repo, err := workitem.Open(p.e.Root)
	if err != nil {
		p.o.Logger.Warn("interrupted close-out jobs not ended", "component", "serve", "root", p.e.Root, "jobs", len(left), "err", err.Error())
		return false
	}
	now := p.o.Now().UTC()
	stories := make([]string, 0, len(left))
	for _, j := range left {
		code := -1
		j.Ended, j.ExitCode, j.Verdict = now, &code, closeout.VerdictInterrupted
		j.Commit = p.head(repo.WorktreePath(j.Story))
		if err := closeout.WriteJob(dir, j); err != nil {
			p.o.Logger.Warn("close-out job's end not written", "component", "serve", "root", p.e.Root, "story", j.Story, "job_id", j.ID, "err", err.Error())
		}
		p.journal(hostapi.Entry{At: now.Format(time.RFC3339), Method: jobMethod(j.Kind), By: j.By, Outcome: closeOutOutcome(closeout.VerdictInterrupted), Verdict: closeout.VerdictInterrupted,
			Detail: fmt.Sprintf("%s job %s: %s, flai serve stopped while it ran", j.Story, j.ID, closeout.VerdictInterrupted)})
		if j.Kind != closeout.KindState {
			p.announce(repo, j.Story, inbox.Verification{Job: j.ID, Commit: j.Commit, Verdict: closeout.VerdictInterrupted})
		}
		stories = append(stories, j.Story)
	}
	p.o.Logger.Warn("close-out jobs interrupted", "component", "serve", "root", p.e.Root, "jobs", len(left), "stories", strings.Join(stories, ","))
	return true
}

// head is the worktree's head commit; "" when git cannot say.
func (p *verifyPass) head(worktree string) string {
	out, err := gitOf(p.o).Run(worktree, "git", "rev-parse", "HEAD")
	if err != nil {
		return ""
	}
	return strings.TrimSpace(out)
}

// journal records a close-out entry in the host journal, stamped now when
// it has no time, of a close-out job when it names no method.
func (p *verifyPass) journal(e hostapi.Entry) {
	if p.o.Host.Record == nil {
		return
	}
	if e.At == "" {
		e.At = p.o.Now().UTC().Format(time.RFC3339)
	}
	if e.Method == "" {
		e.Method = closeOutMethod
	}
	e.Action, e.Project, e.Root = closeOutAction, p.e.Key, p.e.Root
	p.o.Host.Record(e)
}

// jobMethod is the journal method of a job of kind: resolveMethod for a
// rebase and resolve and stateMethod for a state job, which flai stats does
// not count among the close-outs (metrics.md § Close-outs), closeOutMethod
// for a close-out.
func jobMethod(kind string) string {
	switch kind {
	case closeout.KindResolve:
		return resolveMethod
	case closeout.KindState:
		return stateMethod
	}
	return closeOutMethod
}

// jobName is what a job of kind is called in a refusal.
func jobName(kind string) string {
	if kind == closeout.KindResolve {
		return "rebase and resolve"
	}
	return "close-out"
}

// announce appends the verified event of story's close-out to the project's
// event log; a failure loses the news and is logged.
func (p *verifyPass) announce(repo *workitem.Repo, story string, v inbox.Verification) {
	title := ""
	if it, err := repo.Get(story); err == nil {
		title = it.Title
	}
	if err := inbox.AppendEvent(repo, inbox.VerifiedEvent(story, title, verifyServeAuthor, v)); err != nil {
		p.o.Logger.Warn("verified event not logged", "component", "serve", "root", p.e.Root, "story", story, "err", err.Error())
	}
}

// closeOutOutcome is the journal outcome of a close-out's end: done for a pass,
// failed otherwise.
func closeOutOutcome(verdict string) string {
	if verdict == closeout.VerdictPassed {
		return "done"
	}
	return "failed"
}

// atStep is " at step", or "" for no step.
func atStep(step string) string {
	if step == "" {
		return ""
	}
	return " at " + step
}
