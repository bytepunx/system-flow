package serve

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/atomicfile"
	"github.com/bytepunx/system-flow/flai/internal/execx"
	"github.com/bytepunx/system-flow/flai/internal/hostapi"
	"github.com/bytepunx/system-flow/flai/internal/kluster"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/storygit"
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
//
// A story that leaves review, for any column but review, is torn down when
// its namespace has a record: teardown, then its namespace deleted (ADR-0254
// §5). A story accepted while the pass looks is torn down too, and main is
// then promoted into the default namespace; the acceptances seen at one look,
// or while a promotion waits, promote once, at main's newest head. An
// acceptance from before the pass started, or while the action was off,
// promotes nothing. A teardown that fails is not tried again until the story
// moves again, and a promotion that fails is not tried again until the next
// acceptance. A failed teardown opens one thread on the story, or on the
// manifest once the story is archived; a failed promotion opens one on the
// manifest, since it has no one story. Neither moves anything.

// klusterEvery is how often the pass looks at the stories in review.
const klusterEvery = time.Minute

// klusterMethod is the method the pass's entries in the host's journal name.
const klusterMethod = "serve.kluster"

// The kinds of job the pass runs.
const (
	klusterDeploy   = "deploy"
	klusterTeardown = "teardown"
	klusterPromote  = "promote"
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
	// Teardown tears story's namespace down from dir, its worktree, or root
	// when dir is empty, and answers the record it removed, or kept marked
	// failed; kluster.ErrNotDeployed when it has none.
	Teardown(ctx context.Context, root string, m manifest.Manifest, story, dir string) (kluster.Record, error)
	// Promote deploys the head of the main branch branch into the default
	// namespace, and answers its record.
	Promote(ctx context.Context, root string, m manifest.Manifest, branch string) (kluster.Record, error)
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

// Teardown tears story down as flai kluster teardown does.
func (SystemKluster) Teardown(ctx context.Context, root string, m manifest.Manifest, story, dir string) (kluster.Record, error) {
	c, err := kluster.New(root, m, kluster.System{})
	if err != nil {
		return kluster.Record{}, err
	}
	return c.Teardown(ctx, m, story, dir, io.Discard)
}

// Promote promotes main as flai kluster promote does.
func (SystemKluster) Promote(ctx context.Context, root string, m manifest.Manifest, branch string) (kluster.Record, error) {
	c, err := kluster.New(root, m, kluster.System{})
	if err != nil {
		return kluster.Record{}, err
	}
	return c.Promote(ctx, m, branch, io.Discard)
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

// klusterJob is one job of the pass's queue: what it does, to which story
// (none for a promotion), at which head (a deploy's only).
type klusterJob struct {
	kind  string
	story string
	head  string
}

// klusterTried is what the pass last did for a story's namespace, kept in
// klusterTriedPath, so that a head is tried once whether a record was written
// or the deploy was refused before it, and a failed teardown is tried again
// only once the story has moved to another column than Teardown names.
type klusterTried struct {
	Story    string    `json:"story"`
	Head     string    `json:"head,omitempty"`
	Outcome  string    `json:"outcome,omitempty"`
	Teardown string    `json:"teardown,omitempty"` // the story's column when a teardown failed
	At       time.Time `json:"at"`
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

// klusterTriedNamespaces are the story namespaces the pass has tried, by
// their tried files, with the story each names.
func klusterTriedNamespaces(root string) map[string]string {
	files, _ := filepath.Glob(filepath.Join(kluster.Dir(root), "serve", "*.json"))
	out := map[string]string{}
	for _, f := range files {
		ns := strings.TrimSuffix(filepath.Base(f), ".json")
		if t, ok := readKlusterTried(root, ns); ok && t.Story != "" {
			out[ns] = t.Story
		}
	}
	return out
}

// klusterPass deploys the stories in review of one served project into its
// kluster, tears down those that left, and promotes main on acceptance: the
// looks queue jobs, and one worker runs them one at a time.
type klusterPass struct {
	o    Options
	e    Entry
	wg   *sync.WaitGroup
	jobs KlusterRunner

	mu        sync.Mutex
	queue     []klusterJob
	running   *klusterJob
	repromote bool // a promotion found the lock held: the next look queues it again
	wake      chan struct{}

	// lookMu keeps one look at a time; since and seen are the acceptances
	// already seen, as the curate scheduler keeps them, and said the
	// namespaces whose story was not found, logged once.
	lookMu sync.Mutex
	since  time.Time
	seen   map[string]time.Time
	said   map[string]bool
}

// newKlusterPass is the pass for project e. Acceptances before now promote
// nothing.
func newKlusterPass(o Options, e Entry, wg *sync.WaitGroup) *klusterPass {
	return &klusterPass{o: o, e: e, wg: wg, jobs: klusterOf(o), wake: make(chan struct{}, 1),
		since: o.Now().UTC().Truncate(time.Second), seen: map[string]time.Time{}, said: map[string]bool{}}
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

// klusterDeployable is whether manifest m's kluster: block names a deploy
// command. A block with none, as enabling the action writes it (S-0586), is a
// cluster whose deploys are not set up yet: the pass deploys and promotes
// nothing for it, and still tears down what was deployed.
func klusterDeployable(m manifest.Manifest) bool {
	return m.Kluster != nil && len(m.Kluster.Deploy) > 0 && strings.TrimSpace(m.Kluster.Deploy[0]) != ""
}

// look queues a deploy of each story in review whose namespace was neither
// deployed nor tried at its branch's head, a teardown of each deployed story
// that left review, and one promotion when stories were accepted since the
// last look, and wakes the worker. With the action off it queues nothing, and
// the acceptances meanwhile promote nothing. With a block that names no
// deploy it queues only teardowns, and the acceptances meanwhile promote
// nothing; it reads the block again at every look, so a deploy set later is
// picked up at the next.
func (p *klusterPass) look(ctx context.Context) {
	p.lookMu.Lock()
	defer p.lookMu.Unlock()
	now := p.o.Now().UTC()
	repo, err := workitem.Open(p.e.Root)
	if err != nil {
		return
	}
	if !p.on(repo.Manifest) {
		p.since, p.seen = now.Truncate(time.Second), map[string]time.Time{}
		return
	}
	items, err := repo.List(true)
	if err != nil {
		p.o.Logger.Warn("stories not deployed or torn down", "component", "serve", "root", p.e.Root, "err", err.Error())
		return
	}
	var changes []workitem.Change
	changes, _, p.since, p.seen = news(items, nil, p.since, p.seen, now)
	root := mainCheckout(repo)
	deployable := klusterDeployable(repo.Manifest)
	var queued []string
	if deployable {
		queued = p.lookReview(ctx, items, root)
	}
	if ctx.Err() != nil {
		return
	}
	queued = append(queued, p.lookLeft(items, root)...)
	if !deployable {
		p.o.Logger.Debug("kluster block names no deploy, so no story is deployed and main is not promoted", "component", "serve", "root", p.e.Root)
	} else if acceptances := accepted(nil, changes); len(acceptances) > 0 || p.takeRepromote() {
		if p.enqueue(klusterJob{kind: klusterPromote}) {
			queued = append(queued, klusterPromote)
		}
	}
	if len(queued) > 0 {
		p.o.Logger.Info("kluster jobs queued", "component", "serve", "root", p.e.Root, "jobs", strings.Join(queued, ","))
		p.signal()
	}
}

// lookReview queues a deploy of each story in review whose namespace was
// neither deployed nor tried at its branch's head, and names them.
func (p *klusterPass) lookReview(ctx context.Context, items []*workitem.Item, root string) []string {
	git := gitOf(p.o)
	var queued []string
	for _, it := range items {
		if ctx.Err() != nil {
			return queued
		}
		if it.Type != workitem.Story || it.Status != workitem.Review || it.Archived {
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
			queued = append(queued, klusterDeploy+" "+it.ID)
		}
	}
	return queued
}

// lookLeft queues a teardown of each story namespace with a record whose
// story is no longer in review, unless a teardown already failed while the
// story stood where it stands, and forgets what the pass tried for one with
// no record; it names the teardowns queued.
func (p *klusterPass) lookLeft(items []*workitem.Item, root string) []string {
	byID := map[string]*workitem.Item{}
	for _, it := range items {
		if it.Type == workitem.Story {
			byID[it.ID] = it
		}
	}
	namespaces := klusterTriedNamespaces(root)
	if list, err := kluster.Deployments(root, ""); err == nil {
		for _, d := range list {
			if d.Story != "" {
				namespaces[d.Namespace] = d.Story
			}
		}
	} else {
		p.o.Logger.Warn("kluster deployments not read", "component", "serve", "root", p.e.Root, "err", err.Error())
	}
	var queued []string
	for _, ns := range slices.Sorted(maps.Keys(namespaces)) {
		story := namespaces[ns]
		it, ok := byID[workitem.CanonicalID(story)]
		if !ok {
			if !p.said[ns] {
				p.said[ns] = true
				p.o.Logger.Debug("kluster namespace's story not found, left as it is", "component", "serve", "root", p.e.Root, "namespace", ns, "story", story)
			}
			continue
		}
		t, tried := readKlusterTried(root, ns)
		if it.Status == workitem.Review && !it.Archived {
			if tried && t.Teardown != "" {
				// back in review: a teardown that failed before is tried again once it leaves
				t.Teardown = ""
				if err := writeKlusterTried(root, ns, t); err != nil {
					p.o.Logger.Warn("kluster teardown not recorded", "component", "serve", "root", p.e.Root, "story", it.ID, "err", err.Error())
				}
			}
			continue
		}
		_, deployed, _ := kluster.ReadRecord(root, ns)
		if !deployed {
			if tried {
				p.forget(root, ns) // nothing deployed, so nothing to tear down
			}
			continue
		}
		if tried && t.Teardown == it.Status {
			continue // failed while the story stood here: tried again once it moves
		}
		if p.enqueue(klusterJob{kind: klusterTeardown, story: it.ID}) {
			queued = append(queued, klusterTeardown+" "+it.ID)
		}
	}
	return queued
}

// forget removes what the pass tried for namespace ns.
func (p *klusterPass) forget(root, ns string) {
	if err := os.Remove(klusterTriedPath(root, ns)); err != nil && !errors.Is(err, os.ErrNotExist) {
		p.o.Logger.Warn("kluster job's record of what flai serve tried not removed", "component", "serve", "root", p.e.Root, "namespace", ns, "err", err.Error())
	}
}

// takeRepromote reports whether a promotion found the lock held since the
// last look, and clears it.
func (p *klusterPass) takeRepromote() bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	again := p.repromote
	p.repromote = false
	return again
}

// enqueue adds j to the queue unless the same job runs or waits already; a
// deploy waiting for the same story at another head takes the newer head in
// its place. A promotion running does not keep another from waiting, since it
// may have resolved main's head before the acceptance that queues it. It
// reports whether the queue changed.
func (p *klusterPass) enqueue(j klusterJob) bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.running != nil && *p.running == j && j.kind != klusterPromote {
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
// since it was queued, or, for a deploy or a promotion, its deploy command.
func (p *klusterPass) run(ctx context.Context, j klusterJob) {
	repo, err := workitem.Open(p.e.Root)
	if err != nil {
		p.o.Logger.Warn("kluster job not run", "component", "serve", "root", p.e.Root, "kind", j.kind, "story", j.story, "err", err.Error())
		return
	}
	if !p.on(repo.Manifest) {
		return
	}
	if j.kind != klusterTeardown && !klusterDeployable(repo.Manifest) {
		return
	}
	switch j.kind {
	case klusterDeploy:
		if it, err := repo.Get(j.story); err != nil || it.Status != workitem.Review || it.Archived {
			return // left review since it was queued: the next look tears it down
		}
		p.deploy(ctx, repo, j)
	case klusterTeardown:
		p.teardown(ctx, repo, j)
	case klusterPromote:
		p.promote(ctx, repo)
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
	p.checks(ctx, repo, j, ns)
}

// KlusterEnv is the environment a checks run for story gets (ADR-0254 §8):
// the variables of its namespace, as Cluster.Env gives them, while the
// kluster host action is on for the project, its manifest has a kluster:
// block, and the namespace's record passed at the head of the story's
// branch, which git reads in the main checkout. Otherwise it is nil, and so
// is ns, the namespace. Only what cannot be read is an error.
func KlusterEnv(repo *workitem.Repo, story string, git execx.Runner, enabled func(action, root string) bool) (ns string, env []string, err error) {
	root, m := mainCheckout(repo), repo.Manifest
	if m.Kluster == nil || enabled == nil || !enabled(hostapi.ActionKluster, root) {
		return "", nil, nil
	}
	story = workitem.CanonicalID(story)
	ns, err = kluster.StoryNamespace(story)
	if err != nil {
		return "", nil, err
	}
	r, ok, err := kluster.ReadRecord(root, ns)
	if err != nil {
		return "", nil, fmt.Errorf("read the kluster record of %s's namespace %s, to give its checks the namespace: %w", story, ns, err)
	}
	if !ok || r.Outcome != kluster.OutcomePassed {
		return "", nil, nil
	}
	head := ""
	if out, gerr := git.Run(root, "git", "rev-parse", "--verify", "--quiet", "refs/heads/story/"+story+"^{commit}"); gerr == nil {
		head = strings.TrimSpace(out)
	}
	if head == "" || head != r.Commit {
		return "", nil, nil // no branch, or deployed at another head
	}
	c, err := kluster.New(root, m, kluster.System{})
	if err != nil {
		return "", nil, err
	}
	t, err := kluster.StoryTarget(story, r.Checkout, r.Commit)
	if err != nil {
		return "", nil, err
	}
	return ns, c.Env(t), nil
}

// checks runs j's story's checks against its namespace once its deploy has
// passed, in the same job, while the checks host action is on and the
// project names checks; a run that does not pass opens one thread on the
// story and moves nothing, and one that passes resolves it.
func (p *klusterPass) checks(ctx context.Context, repo *workitem.Repo, j klusterJob, ns string) {
	if p.o.Checks == nil {
		return
	}
	cfg, err := p.o.Checks(repo)
	if err != nil {
		p.o.Logger.Warn("kluster checks not run", "component", "serve", "root", p.e.Root, "story", j.story, "err", err.Error())
		return
	}
	if !cfg.Enabled || len(cfg.Commands) == 0 {
		return
	}
	_, env, err := KlusterEnv(repo, j.story, gitOf(p.o), p.o.Host.Enabled)
	if err != nil {
		p.o.Logger.Warn("kluster checks not run", "component", "serve", "root", p.e.Root, "story", j.story, "err", err.Error())
		return
	}
	if env == nil {
		return // the head moved since the deploy: the next look deploys it again
	}
	p.journal("started", fmt.Sprintf("checks of %s against %s at %s", j.story, ns, short(j.head)))
	p.o.Logger.Info("kluster checks started", "component", "serve", "root", p.e.Root, "story", j.story, "namespace", ns, "commit", j.head)
	agents := AgentConfig{}
	if p.o.Agent != nil {
		agents = p.o.Agent(p.e.Root)
	}
	in := ChecksContainer(repo, j.story, agents, dockerOf(p.o))
	run, err := p.o.Dir.RunChecksIn(ctx, j.story, mainCheckout(repo), repo.WorktreePath(j.story), cfg, env, in, p.o.Now)
	if err != nil {
		p.journal("failed", fmt.Sprintf("checks of %s against %s at %s not run: %v", j.story, ns, short(j.head), err))
		p.o.Logger.Info("kluster checks not run", "component", "serve", "root", p.e.Root, "story", j.story, "namespace", ns, "err", err.Error())
		return
	}
	if run.Outcome == "passed" {
		p.journal("done", fmt.Sprintf("checks of %s against %s at %s passed", j.story, ns, short(j.head)))
		p.o.Logger.Info("kluster checks passed", "component", "serve", "root", p.e.Root, "story", j.story, "namespace", ns, "commit", j.head)
		p.clear(repo, j.story, checksFailureTitle(j.story), fmt.Sprintf("flai serve's checks of %s against its namespace %s passed at %s.", j.story, ns, short(j.head)))
		return
	}
	check := failedCheck(run)
	p.journal("failed", fmt.Sprintf("checks of %s against %s at %s %s at %s", j.story, ns, short(j.head), run.Outcome, check))
	p.o.Logger.Warn("kluster checks failed", "component", "serve", "root", p.e.Root, "story", j.story, "namespace", ns, "commit", j.head, "outcome", run.Outcome, "check", check, "log", run.Log)
	p.report(repo, j.story, checksFailureTitle(j.story), checksFailureText(j, ns, run))
}

// failedCheck names the check a run that did not pass ended at.
func failedCheck(run ChecksRun) string {
	if len(run.Steps) == 0 {
		return "no check"
	}
	return "the check " + run.Steps[len(run.Steps)-1].Name
}

func checksFailureTitle(story string) string {
	return fmt.Sprintf("flai serve's checks against %s's kluster namespace failed", story)
}

// checksFailureText names the outcome, the check it ended at, the namespace,
// the commit, and the log, and says what happens next.
func checksFailureText(j klusterJob, ns string, run ChecksRun) string {
	log := "no log"
	if run.Log != "" {
		log = "`" + run.Log + "`"
	}
	return fmt.Sprintf("flai serve's checks of %s against its namespace `%s` at %s ended %s, at %s. Their log is %s.\n\n"+
		"%s stays in review: flai serve moves nothing. It deploys and checks again at the story's next head, and resolves this thread once a checks run there passes. flai checks status %s shows the run.",
		j.story, ns, j.head, run.Outcome, failedCheck(run), log, j.story, j.story)
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
	return failureOf(r, err), err
}

// failureOf says where a job that answered r and err failed: the step, its
// exit code, and its log, from the step's error or else from the record.
func failureOf(r kluster.Record, err error) klusterFailure {
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
	return f
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
	p.report(repo, j.story, klusterFailureTitle(j.story), klusterFailureText(j, ns, f))
}

// passed resolves the story's deploy failure thread, if one is open.
func (p *klusterPass) passed(repo *workitem.Repo, j klusterJob, ns string) {
	p.clear(repo, j.story, klusterFailureTitle(j.story), fmt.Sprintf("flai serve deployed %s into its namespace %s at %s.", j.story, ns, short(j.head)))
}

// teardown tears j's story down, from its worktree or the main checkout when
// it is gone, unless it is back in review, and settles its teardown failure
// thread. A failure is kept with the story's column, so that it is tried
// again only once the story moves.
func (p *klusterPass) teardown(ctx context.Context, repo *workitem.Repo, j klusterJob) {
	it, err := repo.Get(j.story)
	if err != nil {
		p.o.Logger.Warn("kluster teardown not run", "component", "serve", "root", p.e.Root, "story", j.story, "err", err.Error())
		return
	}
	if it.Status == workitem.Review && !it.Archived {
		return // back in review: its deploy updates the namespace in place
	}
	root := mainCheckout(repo)
	ns, _ := kluster.StoryNamespace(j.story)
	dir := ""
	if info, err := os.Stat(repo.WorktreePath(j.story)); err == nil && info.IsDir() {
		dir = repo.WorktreePath(j.story)
	}
	p.journal("started", fmt.Sprintf("teardown of %s's namespace %s, %s", j.story, ns, it.Status))
	p.o.Logger.Info("kluster teardown started", "component", "serve", "root", p.e.Root, "story", j.story, "namespace", ns, "status", it.Status)
	r, err := p.jobs.Teardown(ctx, root, repo.Manifest, j.story, dir)
	switch {
	case errors.Is(err, kluster.ErrNotDeployed):
		p.forget(root, ns)
		p.journal("done", fmt.Sprintf("teardown of %s's namespace %s: nothing was deployed", j.story, ns))
		p.o.Logger.Info("kluster teardown found nothing deployed", "component", "serve", "root", p.e.Root, "story", j.story, "namespace", ns)
	case errors.Is(err, kluster.ErrBusy) || (err != nil && ctx.Err() != nil):
		p.journal("failed", fmt.Sprintf("teardown of %s's namespace %s not run, and queued again at the next look: %v", j.story, ns, err))
		p.o.Logger.Info("kluster teardown not run", "component", "serve", "root", p.e.Root, "story", j.story, "namespace", ns, "err", err.Error())
	case err != nil:
		f := failureOf(r, err)
		t, _ := readKlusterTried(root, ns)
		t.Story, t.Teardown, t.At = j.story, it.Status, p.o.Now().UTC()
		if werr := writeKlusterTried(root, ns, t); werr != nil {
			p.o.Logger.Warn("kluster teardown not recorded", "component", "serve", "root", p.e.Root, "story", j.story, "err", werr.Error())
		}
		p.journal("failed", fmt.Sprintf("teardown of %s's namespace %s failed %s: %s", j.story, ns, f.where(), f.err))
		p.o.Logger.Warn("kluster teardown failed", "component", "serve", "root", p.e.Root, "story", j.story, "namespace", ns, "step", f.step, "exit_code", f.exit, "log", f.log)
		p.report(repo, teardownAnchor(it), teardownFailureTitle(j.story), teardownFailureText(it, ns, r.Commit, f))
	default:
		p.forget(root, ns)
		p.journal("done", fmt.Sprintf("tore down %s's namespace %s", j.story, ns))
		p.o.Logger.Info("kluster teardown passed", "component", "serve", "root", p.e.Root, "story", j.story, "namespace", ns)
		why := fmt.Sprintf("flai serve tore down %s's namespace %s.", j.story, ns)
		p.clear(repo, j.story, teardownFailureTitle(j.story), why)
		p.clear(repo, manifest.File, teardownFailureTitle(j.story), why)
	}
}

// teardownAnchor is where a story's teardown failure thread goes: on the
// story, or on the manifest once the story is archived, since a thread on an
// archived item is a finding of its own.
func teardownAnchor(it *workitem.Item) string {
	if it.Archived {
		return manifest.File
	}
	return it.ID
}

// promote promotes the main branch's head into the default namespace, and
// settles the promotion failure thread. One that finds the lock held is
// queued again at the next look; one that fails waits for the next
// acceptance.
func (p *klusterPass) promote(ctx context.Context, repo *workitem.Repo) {
	root, m := mainCheckout(repo), repo.Manifest
	ns := m.KlusterNamespace()
	p.journal("started", fmt.Sprintf("promotion of main into %s", ns))
	p.o.Logger.Info("kluster promotion started", "component", "serve", "root", p.e.Root, "namespace", ns)
	branch, err := storygit.MainBranch(gitOf(p.o), root)
	var r kluster.Record
	if err != nil {
		err = fmt.Errorf("find the main branch of %s to promote: %w", root, err)
	} else {
		r, err = p.jobs.Promote(ctx, root, m, branch)
	}
	switch {
	case errors.Is(err, kluster.ErrBusy) || (err != nil && ctx.Err() != nil):
		p.mu.Lock()
		p.repromote = true
		p.mu.Unlock()
		p.journal("failed", fmt.Sprintf("promotion of main into %s not run, and queued again at the next look: %v", ns, err))
		p.o.Logger.Info("kluster promotion not run", "component", "serve", "root", p.e.Root, "namespace", ns, "err", err.Error())
	case err != nil:
		f := failureOf(r, err)
		p.journal("failed", fmt.Sprintf("promotion of main into %s at %s failed %s: %s", ns, orUnresolved(r.Commit), f.where(), f.err))
		p.o.Logger.Warn("kluster promotion failed", "component", "serve", "root", p.e.Root, "namespace", ns, "commit", r.Commit, "step", f.step, "exit_code", f.exit, "log", f.log)
		p.report(repo, manifest.File, promoteFailureTitle(), promoteFailureText(ns, r.Commit, f))
	default:
		p.journal("done", fmt.Sprintf("promoted main into %s at %s", ns, short(r.Commit)))
		p.o.Logger.Info("kluster promotion passed", "component", "serve", "root", p.e.Root, "namespace", ns, "commit", r.Commit)
		p.clear(repo, manifest.File, promoteFailureTitle(), fmt.Sprintf("flai serve promoted main into the default namespace %s at %s.", ns, short(r.Commit)))
	}
}

// report opens the thread titled title on on, an item or a path, or adds text
// to the one open.
func (p *klusterPass) report(repo *workitem.Repo, on, title, text string) {
	now := p.o.Now()
	if th := klusterThread(repo, on, title); th != nil {
		if _, err := threads.Reply(repo, th.ID, forgeAuthor, text, now); err != nil {
			p.o.Logger.Warn("kluster failure not written on its thread", "component", "serve", "root", p.e.Root, "on", on, "thread", th.ID, "err", err.Error())
		}
		return
	}
	if _, err := threads.New(repo, threads.NewOptions{Title: title, On: on, Author: forgeAuthor, Text: text, Now: now}); err != nil {
		p.o.Logger.Warn("kluster failure thread not opened", "component", "serve", "root", p.e.Root, "on", on, "err", err.Error())
	}
}

// clear resolves the thread titled title on on, if one is open, saying why.
func (p *klusterPass) clear(repo *workitem.Repo, on, title, why string) {
	th := klusterThread(repo, on, title)
	if th == nil {
		return
	}
	if _, err := threads.Resolve(repo, th.ID, forgeAuthor, why, p.o.Now()); err != nil {
		p.o.Logger.Warn("kluster failure thread not resolved", "component", "serve", "root", p.e.Root, "on", on, "thread", th.ID, "err", err.Error())
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

// klusterThread is the open thread titled title that flai serve opened on
// on, an item or a path, or nil.
func klusterThread(repo *workitem.Repo, on, title string) *threads.Thread {
	all, err := threads.For(repo, on)
	if err != nil {
		return nil
	}
	title = workitem.CleanTitle(title)
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
	exit, log := f.said()
	return fmt.Sprintf("flai serve's %s of %s into the namespace `%s` at %s failed %s, with %s. Its log is %s.\n\n```text\n%s\n```\n\n"+
		"%s stays in review: flai serve moves nothing. It deploys again at the story's next head, and resolves this thread once a deploy passes. flai kluster deployments shows the namespace's record.",
		j.kind, j.story, ns, j.head, f.where(), exit, log, f.err, j.story)
}

func teardownFailureTitle(story string) string {
	return fmt.Sprintf("flai serve could not tear down %s's kluster namespace", story)
}

// teardownFailureText names the story and where it stands, the namespace,
// the commit deployed there, the step, its exit code, the log, and the error,
// and says what happens next.
func teardownFailureText(it *workitem.Item, ns, commit string, f klusterFailure) string {
	exit, log := f.said()
	return fmt.Sprintf("flai serve's teardown of %s's namespace `%s`, deployed at %s, failed %s, with %s, once %s was %s. Its log is %s.\n\n```text\n%s\n```\n\n"+
		"flai serve moves nothing, keeps the namespace's record, and tries again only once %s moves again. flai kluster teardown %s tears it down from a shell; this thread is resolved once a teardown passes.",
		it.ID, ns, orUnresolved(commit), f.where(), exit, it.ID, it.Status, log, f.err, it.ID, it.ID)
}

func promoteFailureTitle() string {
	return "flai serve could not promote main into the kluster's default namespace"
}

// promoteFailureText names the namespace, main's commit, the step, its exit
// code, the log, and the error, and says what happens next.
func promoteFailureText(ns, commit string, f klusterFailure) string {
	exit, log := f.said()
	return fmt.Sprintf("flai serve's promotion of main into the default namespace `%s` at %s failed %s, with %s. Its log is %s.\n\n```text\n%s\n```\n\n"+
		"The default namespace keeps what the failed promotion left, and flai serve moves nothing. The next acceptance promotes again, and flai kluster promote does from a shell; this thread is resolved once a promotion passes.",
		ns, orUnresolved(commit), f.where(), exit, log, f.err)
}

// orUnresolved is commit, or says it was not resolved.
func orUnresolved(commit string) string {
	if commit == "" {
		return "a commit not resolved"
	}
	return commit
}

// said names the failure's exit code and its log, in a thread's words.
func (f klusterFailure) said() (exit, log string) {
	exit, log = "no exit code", "no log"
	if f.exit >= 0 {
		exit = "exit code " + strconv.Itoa(f.exit)
	}
	if f.log != "" {
		log = "`" + f.log + "`"
	}
	return exit, log
}

// where says where the job failed: at its step, or before its first.
func (f klusterFailure) where() string {
	if f.step == "" {
		return "before its first step"
	}
	return "at the " + f.step + " step"
}
