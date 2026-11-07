package serve

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/atomicfile"
	"github.com/bytepunx/system-flow/flai/internal/hostapi"
	"github.com/bytepunx/system-flow/flai/internal/kluster"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/threads"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// Deploying the stories in review into the project's kluster (S-0579,
// ADR-0254 §5). Once a minute, and at once when a project is first served,
// for each served project whose manifest has a kluster: block while its
// kluster host action is on, flai serve queues a deploy of each story in
// review whose namespace was not deployed, or tried, at its branch's head:
// once when the story enters review, and again at each new head. With the
// action off or no block, it runs no kluster command at all.
//
// The jobs of a project run one at a time, in the order their causes came,
// on one worker of the project's own, so that a look never waits for a job.
// A deploy that finds the cluster absent brings it up first. A job refused
// because another holds the project's kluster lock, a flai kluster run from a
// shell say, is no failure: the next look queues it again. A deploy that
// fails opens one thread on the story naming the step, its exit code, the
// commit, and the log, adds an entry to that thread when a later head fails
// too, and moves nothing; a deploy that passes resolves it. A head tried
// once, passed or failed, is not tried again, so a failure is not retried
// every minute.

// klusterEvery is how often the pass looks at the stories in review.
const klusterEvery = time.Minute

// klusterMethod is the method the pass's entries in the host's journal name.
const klusterMethod = "serve.kluster"

// The kinds of job the pass runs.
const (
	klusterDeploy = "deploy"
)

// KlusterRunner runs a project's kluster jobs for flai serve's kluster pass,
// from root, the project's main checkout, with the cluster m's kluster: block
// names. SystemKluster is the real one; tests give their own.
type KlusterRunner interface {
	// Absent says whether the cluster must be brought up before a deploy.
	Absent(ctx context.Context, root string, m manifest.Manifest) (bool, error)
	// Up brings the cluster up and keeps its kubeconfig, its output to out.
	Up(ctx context.Context, root string, m manifest.Manifest, out io.Writer) error
	// Deploy deploys story, in review, from its worktree dir at commit into
	// its namespace, and answers the namespace's record.
	Deploy(ctx context.Context, root string, m manifest.Manifest, story, dir, commit string) (kluster.Record, error)
}

// SystemKluster is the KlusterRunner with this host's kluster, docker, and
// kubectl, through flai/internal/kluster, as flai kluster runs them.
type SystemKluster struct{}

// Absent is true when the cluster has no kubeconfig at
// kluster.KubeconfigPath, or kluster status does not list it running.
func (SystemKluster) Absent(ctx context.Context, root string, m manifest.Manifest) (bool, error) {
	if _, err := os.Stat(kluster.KubeconfigPath(root)); errors.Is(err, os.ErrNotExist) {
		return true, nil
	} else if err != nil {
		return false, fmt.Errorf("read the kubeconfig %s: %w", kluster.KubeconfigPath(root), err)
	}
	c, err := kluster.New(root, m, kluster.System{})
	if err != nil {
		return false, err
	}
	st, err := c.Status(ctx)
	if err != nil {
		return false, err
	}
	return !st.Exists || !st.Running, nil
}

// Up brings the cluster up as flai kluster up does.
func (SystemKluster) Up(ctx context.Context, root string, m manifest.Manifest, out io.Writer) error {
	c, err := kluster.New(root, m, kluster.System{})
	if err != nil {
		return err
	}
	return c.Up(ctx, out)
}

// Deploy deploys story as flai kluster deploy does, its progress lines left
// out: each step's output is in its log.
func (SystemKluster) Deploy(ctx context.Context, root string, m manifest.Manifest, story, dir, commit string) (kluster.Record, error) {
	c, err := kluster.New(root, m, kluster.System{})
	if err != nil {
		return kluster.Record{}, err
	}
	return c.Deploy(ctx, m, story, workitem.Review, dir, commit, io.Discard)
}

// klusterOf is the runner the options name, SystemKluster when none.
func klusterOf(o Options) KlusterRunner {
	if o.Kluster == nil {
		return SystemKluster{}
	}
	return o.Kluster
}

// klusterUpLog is the file a bring-up's output goes to, the one the
// dashboard's kluster up writes, so that the Klusters page's log is the last
// bring-up's whoever ran it.
func klusterUpLog(root string) string { return filepath.Join(kluster.Dir(root), "up.log") }

// klusterJob is one job of the pass's queue: what it does, to which story,
// at which head.
type klusterJob struct {
	kind  string
	story string
	head  string
}

// klusterTried is what the pass last did for a story's namespace, kept in
// klusterTriedPath, so that a head is tried once whether a record was written
// or the deploy was refused before it.
type klusterTried struct {
	Story   string    `json:"story"`
	Head    string    `json:"head"`
	Outcome string    `json:"outcome"`
	At      time.Time `json:"at"`
}

// klusterTriedPath is where the pass keeps what it last did for namespace ns.
func klusterTriedPath(root, ns string) string {
	return filepath.Join(kluster.Dir(root), "serve", ns+".json")
}

func readKlusterTried(root, ns string) (klusterTried, bool) {
	data, err := os.ReadFile(klusterTriedPath(root, ns))
	if err != nil {
		return klusterTried{}, false
	}
	var t klusterTried
	if json.Unmarshal(data, &t) != nil {
		return klusterTried{}, false
	}
	return t, true
}

func writeKlusterTried(root, ns string, t klusterTried) error {
	path := klusterTriedPath(root, ns)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("make the folder of flai serve's kluster jobs: %w", err)
	}
	data, err := json.MarshalIndent(t, "", "  ")
	if err != nil {
		return fmt.Errorf("encode what flai serve did for %s: %w", ns, err)
	}
	if err := atomicfile.WriteFile(path, append(data, '\n'), 0o644); err != nil {
		return fmt.Errorf("write what flai serve did for %s: %w", ns, err)
	}
	return nil
}

// klusterPass deploys the stories in review of one served project into its
// kluster: the looks queue jobs, and one worker runs them one at a time.
type klusterPass struct {
	o    Options
	e    Entry
	wg   *sync.WaitGroup
	jobs KlusterRunner

	mu      sync.Mutex
	queue   []klusterJob
	running *klusterJob
	wake    chan struct{}
}

func newKlusterPass(o Options, e Entry, wg *sync.WaitGroup) *klusterPass {
	return &klusterPass{o: o, e: e, wg: wg, jobs: klusterOf(o), wake: make(chan struct{}, 1)}
}

// start looks now and then every minute, and runs the jobs the looks queue on
// a worker of its own, both until ctx ends: a deploy takes minutes, and
// neither the project's other looks nor this pass's must wait for it.
func (p *klusterPass) start(ctx context.Context) {
	p.wg.Add(2)
	go func() {
		defer p.wg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			case <-p.wake:
				p.drain(ctx)
			}
		}
	}()
	go func() {
		defer p.wg.Done()
		for {
			p.look(ctx)
			select {
			case <-ctx.Done():
				return
			case <-time.After(klusterEvery):
			}
		}
	}()
}

// on is whether the pass may run anything for the project: its kluster host
// action is on and its manifest m has a kluster: block.
func (p *klusterPass) on(m manifest.Manifest) bool {
	return m.Kluster != nil && p.o.Host.Enabled != nil && p.o.Host.Enabled(hostapi.ActionKluster, p.e.Root)
}

// look queues a deploy of each story in review whose namespace was neither
// deployed nor tried at its branch's head, and wakes the worker.
func (p *klusterPass) look(ctx context.Context) {
	repo, err := workitem.Open(p.e.Root)
	if err != nil || !p.on(repo.Manifest) {
		return
	}
	items, err := repo.List(false)
	if err != nil {
		p.o.Logger.Warn("stories in review not deployed", "component", "serve", "root", p.e.Root, "err", err.Error())
		return
	}
	root := mainCheckout(repo)
	git := gitOf(p.o)
	var queued []string
	for _, it := range items {
		if ctx.Err() != nil {
			return
		}
		if it.Type != workitem.Story || it.Status != workitem.Review {
			continue
		}
		ns, err := kluster.StoryNamespace(it.ID)
		if err != nil {
			continue
		}
		out, err := git.Run(root, "git", "rev-parse", "--verify", "--quiet", "refs/heads/story/"+it.ID+"^{commit}")
		head := strings.TrimSpace(out)
		if err != nil || head == "" {
			p.o.Logger.Debug("story in review has no branch to deploy", "component", "serve", "root", p.e.Root, "story", it.ID)
			continue
		}
		if r, ok, _ := kluster.ReadRecord(root, ns); ok && r.Commit == head {
			continue
		}
		if t, ok := readKlusterTried(root, ns); ok && t.Head == head {
			continue
		}
		if p.enqueue(klusterJob{kind: klusterDeploy, story: it.ID, head: head}) {
			queued = append(queued, it.ID)
		}
	}
	if len(queued) > 0 {
		p.o.Logger.Info("stories in review queued for deploy", "component", "serve", "root", p.e.Root, "stories", strings.Join(queued, ","))
		p.signal()
	}
}

// enqueue adds j to the queue unless the same job runs or waits already; a
// deploy waiting for the same story at another head takes the newer head in
// its place. It reports whether the queue changed.
func (p *klusterPass) enqueue(j klusterJob) bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.running != nil && *p.running == j {
		return false
	}
	for i, q := range p.queue {
		if q.kind == j.kind && q.story == j.story {
			if q == j {
				return false
			}
			p.queue[i] = j
			return true
		}
	}
	p.queue = append(p.queue, j)
	return true
}

// signal wakes the worker without waiting for it.
func (p *klusterPass) signal() {
	select {
	case p.wake <- struct{}{}:
	default:
	}
}

// next takes the first job off the queue and marks it running; false when
// the queue is empty.
func (p *klusterPass) next() (klusterJob, bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if len(p.queue) == 0 {
		return klusterJob{}, false
	}
	j := p.queue[0]
	p.queue = p.queue[1:]
	p.running = &j
	return j, true
}

func (p *klusterPass) done() {
	p.mu.Lock()
	p.running = nil
	p.mu.Unlock()
}

// drain runs the queued jobs one at a time, in order, until the queue is
// empty or ctx ends.
func (p *klusterPass) drain(ctx context.Context) {
	for ctx.Err() == nil {
		j, ok := p.next()
		if !ok {
			return
		}
		p.run(ctx, j)
		p.done()
	}
}

// run runs one job, unless the action was turned off or the block removed
// since it was queued.
func (p *klusterPass) run(ctx context.Context, j klusterJob) {
	repo, err := workitem.Open(p.e.Root)
	if err != nil {
		p.o.Logger.Warn("kluster job not run", "component", "serve", "root", p.e.Root, "kind", j.kind, "story", j.story, "err", err.Error())
		return
	}
	if !p.on(repo.Manifest) {
		return
	}
	if j.kind == klusterDeploy {
		p.deploy(ctx, repo, j)
	}
}

// klusterFailure is why a job failed: the step, empty when the job was
// refused before its first step, its exit code (-1 for none), its log, and
// the error.
type klusterFailure struct {
	step string
	exit int
	log  string
	err  string
}

// deploy deploys j's story at j.head, bringing the cluster up first when it
// is absent, and settles the story's failure thread.
func (p *klusterPass) deploy(ctx context.Context, repo *workitem.Repo, j klusterJob) {
	root := mainCheckout(repo)
	ns, _ := kluster.StoryNamespace(j.story)
	p.journal("started", fmt.Sprintf("deploy of %s into %s at %s", j.story, ns, short(j.head)))
	p.o.Logger.Info("kluster deploy started", "component", "serve", "root", p.e.Root, "story", j.story, "namespace", ns, "commit", j.head)
	fail, err := p.deployOnce(ctx, repo, j)
	if errors.Is(err, kluster.ErrBusy) || (err != nil && ctx.Err() != nil) {
		// another job holds the lock, or flai serve stops: the next look queues it again
		p.journal("failed", fmt.Sprintf("deploy of %s into %s at %s not run, and queued again at the next look: %v", j.story, ns, short(j.head), err))
		p.o.Logger.Info("kluster deploy not run", "component", "serve", "root", p.e.Root, "story", j.story, "namespace", ns, "commit", j.head, "err", err.Error())
		return
	}
	tried := klusterTried{Story: j.story, Head: j.head, Outcome: kluster.OutcomePassed, At: p.o.Now().UTC()}
	if err != nil {
		tried.Outcome = kluster.OutcomeFailed
	}
	if werr := writeKlusterTried(root, ns, tried); werr != nil {
		p.o.Logger.Warn("kluster deploy not recorded", "component", "serve", "root", p.e.Root, "story", j.story, "err", werr.Error())
	}
	if err != nil {
		p.journal("failed", fmt.Sprintf("deploy of %s into %s at %s failed %s: %s", j.story, ns, short(j.head), fail.where(), fail.err))
		p.o.Logger.Warn("kluster deploy failed", "component", "serve", "root", p.e.Root, "story", j.story, "namespace", ns, "commit", j.head, "step", fail.step, "exit_code", fail.exit, "log", fail.log)
		p.failed(repo, j, ns, fail)
		return
	}
	p.journal("done", fmt.Sprintf("deployed %s into %s at %s", j.story, ns, short(j.head)))
	p.o.Logger.Info("kluster deploy passed", "component", "serve", "root", p.e.Root, "story", j.story, "namespace", ns, "commit", j.head)
	p.passed(repo, j, ns)
}

// deployOnce brings the cluster up when it is absent, then deploys j's story
// from its worktree; on a failure it says which step failed and where its
// log is.
func (p *klusterPass) deployOnce(ctx context.Context, repo *workitem.Repo, j klusterJob) (klusterFailure, error) {
	root, m := mainCheckout(repo), repo.Manifest
	dir := repo.WorktreePath(j.story)
	if info, err := os.Stat(dir); err != nil || !info.IsDir() {
		err := fmt.Errorf("%s has no worktree at %s to build and deploy from; open it with flai stream open %s", j.story, dir, j.story)
		return klusterFailure{exit: -1, err: err.Error()}, err
	}
	absent, err := p.jobs.Absent(ctx, root, m)
	if err != nil {
		return klusterFailure{step: "up", exit: -1, err: err.Error()}, err
	}
	if absent {
		if err := p.up(ctx, root, m); err != nil {
			return klusterFailure{step: "up", exit: kluster.ExitCode(err), log: klusterUpLog(root), err: err.Error()}, err
		}
	}
	r, err := p.jobs.Deploy(ctx, root, m, j.story, dir, j.head)
	if err == nil {
		return klusterFailure{}, nil
	}
	f := klusterFailure{exit: -1, err: err.Error()}
	var se *kluster.StepError
	switch {
	case errors.As(err, &se):
		f.step, f.exit, f.log = se.Step, se.ExitCode, se.Log
	case r.FailedStep != "":
		f.step, f.exit = r.FailedStep, r.ExitCode
		for _, s := range r.Steps {
			if s.Name == r.FailedStep {
				f.log = s.Log
			}
		}
	}
	return f, err
}

// up brings the cluster up within kluster.timeout, its output to the
// bring-up log.
func (p *klusterPass) up(ctx context.Context, root string, m manifest.Manifest) error {
	if p.o.Host.Kluster != nil {
		if running, _ := p.o.Host.Kluster.State(root); running != nil {
			return fmt.Errorf("%w: the dashboard's kluster %s of the cluster %s runs since %s", kluster.ErrBusy, running.Action, running.Cluster, running.StartedAt)
		}
	}
	if err := os.MkdirAll(kluster.Dir(root), 0o755); err != nil {
		return fmt.Errorf("make the folder for the kluster's bring-up log: %w", err)
	}
	f, err := os.Create(klusterUpLog(root))
	if err != nil {
		return fmt.Errorf("open the kluster's bring-up log %s: %w", klusterUpLog(root), err)
	}
	ctx, cancel := context.WithTimeout(ctx, m.KlusterTimeout())
	defer cancel()
	p.o.Logger.Info("kluster absent, brought up before a deploy", "component", "serve", "root", p.e.Root, "log", klusterUpLog(root))
	err = p.jobs.Up(ctx, root, m, f)
	if cerr := f.Close(); err == nil && cerr != nil {
		err = fmt.Errorf("close the kluster's bring-up log %s: %w", klusterUpLog(root), cerr)
	}
	return err
}

// failed opens the story's deploy failure thread, or adds the new failure to
// the one open.
func (p *klusterPass) failed(repo *workitem.Repo, j klusterJob, ns string, f klusterFailure) {
	text := klusterFailureText(j, ns, f)
	now := p.o.Now()
	if th := klusterFailureThread(repo, j.story); th != nil {
		if _, err := threads.Reply(repo, th.ID, forgeAuthor, text, now); err != nil {
			p.o.Logger.Warn("kluster deploy failure not written on its thread", "component", "serve", "root", p.e.Root, "story", j.story, "thread", th.ID, "err", err.Error())
		}
		return
	}
	if _, err := threads.New(repo, threads.NewOptions{Title: klusterFailureTitle(j.story), On: j.story, Author: forgeAuthor, Text: text, Now: now}); err != nil {
		p.o.Logger.Warn("kluster deploy failure thread not opened", "component", "serve", "root", p.e.Root, "story", j.story, "err", err.Error())
	}
}

// passed resolves the story's deploy failure thread, if one is open.
func (p *klusterPass) passed(repo *workitem.Repo, j klusterJob, ns string) {
	th := klusterFailureThread(repo, j.story)
	if th == nil {
		return
	}
	why := fmt.Sprintf("flai serve deployed %s into its namespace %s at %s.", j.story, ns, short(j.head))
	if _, err := threads.Resolve(repo, th.ID, forgeAuthor, why, p.o.Now()); err != nil {
		p.o.Logger.Warn("kluster deploy failure thread not resolved", "component", "serve", "root", p.e.Root, "story", j.story, "thread", th.ID, "err", err.Error())
	}
}

// journal keeps one entry of the pass's in the host's journal.
func (p *klusterPass) journal(outcome, detail string) {
	if p.o.Host.Record == nil {
		return
	}
	p.o.Host.Record(hostapi.Entry{At: p.o.Now().UTC().Format(time.RFC3339), Action: hostapi.ActionKluster, Method: klusterMethod, Project: p.e.Key, Root: p.e.Root,
		By: forgeAuthor, Outcome: outcome, Detail: detail})
}

// klusterFailureThread is the open thread flai serve opened on story for a
// failed deploy, or nil.
func klusterFailureThread(repo *workitem.Repo, story string) *threads.Thread {
	all, err := threads.For(repo, story)
	if err != nil {
		return nil
	}
	title := workitem.CleanTitle(klusterFailureTitle(story))
	for _, th := range all {
		if th.Open() && th.Opener() == forgeAuthor && th.Title == title {
			return th
		}
	}
	return nil
}

func klusterFailureTitle(story string) string {
	return fmt.Sprintf("flai serve could not deploy %s into its kluster namespace", story)
}

// klusterFailureText names the job, the step, its exit code, the commit, the
// log, and the error, and says what happens next.
func klusterFailureText(j klusterJob, ns string, f klusterFailure) string {
	exit := "no exit code"
	if f.exit >= 0 {
		exit = "exit code " + strconv.Itoa(f.exit)
	}
	log := "no log"
	if f.log != "" {
		log = "`" + f.log + "`"
	}
	return fmt.Sprintf("flai serve's %s of %s into the namespace `%s` at %s failed %s, with %s. Its log is %s.\n\n```text\n%s\n```\n\n"+
		"%s stays in review: flai serve moves nothing. It deploys again at the story's next head, and resolves this thread once a deploy passes. flai kluster deployments shows the namespace's record.",
		j.kind, j.story, ns, j.head, f.where(), exit, log, f.err, j.story)
}

// where says where the job failed: at its step, or before its first.
func (f klusterFailure) where() string {
	if f.step == "" {
		return "before its first step"
	}
	return "at the " + f.step + " step"
}
