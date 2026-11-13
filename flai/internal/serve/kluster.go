package serve

import (
	"cmp"
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
// fails, or a checks run against the namespace that does not pass, opens one
// thread on the story naming the job, the step, its exit code, the commit,
// and the log, or adds an entry to the one open, and sends the story back
// from review to in-progress, as flai, with a reason naming the job, the step,
// and the commit (ADR-0254 §5a), as a failing required check does (ADR-0167):
// unless the action went off meanwhile, or the story is no longer in review
// at that head. The send back is kept with what the pass tried, naming the
// thread, so that the story's agent is started on it once
// (startKlusterSentBack), and the story keeps its namespace while it stands
// in progress where the send back put it. Back in review, at a new head or
// the same one, it is deployed again in place; a deploy that passes resolves
// the deploy's thread, and a checks run that passes the checks'. A head tried
// once, passed or failed, is not tried again until the story returns to
// review after a send back, so a failure is not retried every minute.
//
// A story's agent asks the pass, too, to deploy its own story's namespace,
// run the project's checks against it as deployed, or tear it down, while
// the story is in progress (ADR-0292, kluster_requests.go): those jobs wait
// in the same queue, run whatever the tried-once rule says, and end with a
// kluster event rather than a thread. A checks run an agent asks for is the
// review's, through the same code, and keeps its outcome in the namespace's
// record as the review's does. A story in review whose
// namespace's record passed at its branch's head, its checks too while the
// checks host action is on, is not deployed again, whoever asked for it.
//
// A namespace stays while its story is in progress, whoever moved it there
// (ADR-0292 §3). A story that moves to ready or backlog is torn down when its
// namespace has a record: teardown, then its namespace deleted (ADR-0254 §5);
// cancelled or accepted, it is torn down all the same, and its failure
// threads resolved. A story accepted while the
// pass looks is torn down too, and main is
// then promoted into the default namespace; the acceptances seen at one look,
// or while a promotion waits, promote once, at main's newest head. An
// acceptance from before the pass started, or while the action was off,
// promotes nothing. A teardown that fails is not tried again until the story
// moves again, and a promotion that fails is not tried again until the next
// acceptance. A failed teardown opens one thread on the story, or on the
// manifest once the story is archived; a failed promotion opens one on the
// manifest, since it has no one story. Neither moves anything: neither has a
// story in review to send back.

// klusterEvery is how often the pass looks at the stories in review.
const klusterEvery = time.Minute

// klusterMethod is the method the pass's entries in the host's journal name.
const klusterMethod = "serve.kluster"

// The kinds of job the pass runs.
const (
	klusterDeploy   = "deploy"
	klusterChecks   = "checks"
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
	// Deploy deploys story, whose status is status, in progress or in review,
	// from its worktree dir at commit into its namespace, for by, who asked
	// for it, and answers the namespace's record.
	Deploy(ctx context.Context, root string, m manifest.Manifest, by, story, status, dir, commit string) (kluster.Record, error)
	// Teardown tears story's namespace down from dir, its worktree, or root
	// when dir is empty, for by, who asked for it, and answers the record it
	// removed, or kept marked failed; kluster.ErrNotDeployed when it has none.
	Teardown(ctx context.Context, root string, m manifest.Manifest, by, story, dir string) (kluster.Record, error)
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
func (SystemKluster) Deploy(ctx context.Context, root string, m manifest.Manifest, by, story, status, dir, commit string) (kluster.Record, error) {
	c, err := kluster.New(root, m, kluster.System{})
	if err != nil {
		return kluster.Record{}, err
	}
	return c.DeployBy(ctx, m, by, story, status, dir, commit, io.Discard)
}

// Teardown tears story down as flai kluster teardown does.
func (SystemKluster) Teardown(ctx context.Context, root string, m manifest.Manifest, by, story, dir string) (kluster.Record, error) {
	c, err := kluster.New(root, m, kluster.System{})
	if err != nil {
		return kluster.Record{}, err
	}
	return c.TeardownBy(ctx, m, by, story, dir, io.Discard)
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
// (none for a promotion), at which head (a deploy's only). id names a story's
// job to a request that finds it waiting or running; by is the agent that
// asked for it, empty for a job the looks queued; queued is when it was.
type klusterJob struct {
	kind   string
	story  string
	head   string
	id     string
	by     string
	queued time.Time
}

// asked reports whether a story's agent asked for j (ADR-0292 §1), rather
// than a look queuing it.
func (j klusterJob) asked() bool { return j.by != "" }

// asker is who asked for j, as the record's AskedBy and the journal name it:
// the agent, or flai serve for a job the looks queued.
func (j klusterJob) asker() string {
	if j.asked() {
		return j.by
	}
	return forgeAuthor
}

// refuseOverLimit refuses the deploy request in file for kluster.agent_namespaces,
// limit, with the reason format and args make, as refuse does, and keeps the
// refusal under the project's .flai-cache/kluster, so that flai kluster
// deployments and the Deployments page list it (S-0663). A refusal it cannot
// keep is logged, and refused all the same.
func (p *klusterPass) refuseOverLimit(file kluster.Request, limit int, format string, args ...any) (kluster.Answer, bool) {
	why := fmt.Sprintf(format, args...)
	r := kluster.Refusal{Story: file.Story, By: file.By, At: p.o.Now().UTC(), Limit: limit, Reason: why}
	if err := kluster.KeepRefusal(p.klusterRequestRoot(), r); err != nil {
		p.o.Logger.Warn("kluster refusal for the limit not kept", "component", "serve", "root", p.e.Root, "story", file.Story, "by", file.By, "err", err.Error())
	}
	return p.refuse(file, why), false
}

// klusterTried is what the pass last did for a story's namespace, kept in
// klusterTriedPath, so that a head is tried once whether a record was written
// or the deploy was refused before it, and a failed teardown is tried again
// only once the story has moved to another column than Teardown names. It
// keeps too the send back a failure at Head made of the story, if one did.
type klusterTried struct {
	Story    string           `json:"story"`
	Head     string           `json:"head,omitempty"`
	Outcome  string           `json:"outcome,omitempty"`
	Teardown string           `json:"teardown,omitempty"` // the story's column when a teardown failed
	SentBack *klusterSentBack `json:"sent_back,omitempty"`
	At       time.Time        `json:"at"`
}

// klusterSentBack is the pass's send back of a story from review to
// in-progress after its deploy or its checks run failed (ADR-0254 §5a): the
// head it failed at, the thread it opened or added to, the move's reason,
// and the move's time, to the second, as the story's transition keeps it.
// The next deploy the story's return to review queues writes what the pass
// tried afresh, without it.
type klusterSentBack struct {
	Head   string    `json:"head"`
	Thread string    `json:"thread"`
	Reason string    `json:"reason"`
	At     time.Time `json:"at"`
}

// klusterKept is whether it stands where the send back t records put it:
// in progress, not archived, its newest move a move to in-progress by flai at
// or after the send back. Such a story's agent is started on the send back's
// thread; once it moves on, to review, cancelled, or anywhere else, or
// someone else moves it back to in-progress, it is not. Its namespace stays
// either way while it is in progress (ADR-0292 §3).
func klusterKept(it *workitem.Item, t klusterTried) bool {
	if t.SentBack == nil || it.Status != workitem.InProgress || it.Archived || len(it.Transitions) == 0 {
		return false
	}
	last := it.Transitions[len(it.Transitions)-1]
	at, err := time.Parse(workitem.TimeFormat, last.At)
	return err == nil && last.To == workitem.InProgress && last.By == sendBackBy && !at.Before(t.SentBack.At)
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

	// swept is whether the asked-for jobs an earlier flai serve left waiting
	// or running were ended; only the requests' loop reads and writes it.
	swept bool
}

// newKlusterPass is the pass for project e. Acceptances before now promote
// nothing.
func newKlusterPass(o Options, e Entry, wg *sync.WaitGroup) *klusterPass {
	return &klusterPass{o: o, e: e, wg: wg, jobs: klusterOf(o), wake: make(chan struct{}, 1),
		since: o.Now().UTC().Truncate(time.Second), seen: map[string]time.Time{}, said: map[string]bool{}}
}

// start looks now and then every minute, takes the story agents' requests
// every second, and runs the jobs the looks and the requests queue on a
// worker of its own, all until ctx ends: a deploy takes minutes, and neither
// the project's other looks nor this pass's must wait for it.
func (p *klusterPass) start(ctx context.Context) {
	p.wg.Add(3)
	go func() {
		defer p.wg.Done()
		for {
			p.takeRequests(ctx)
			select {
			case <-ctx.Done():
				return
			case <-time.After(klusterRequestEvery):
			}
		}
	}()
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
		queued = p.lookReview(ctx, items, root, repo.Manifest, p.checksOn(repo))
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

// checksOn is whether a passing deploy of the project's stories is followed
// by a checks run against the namespace: the checks host action is on and
// names checks.
func (p *klusterPass) checksOn(repo *workitem.Repo) bool {
	_, why := p.checksConfig(repo)
	return why == ""
}

// passedAt reports whether record r says its namespace was deployed at head
// and passed, and, while checksOn, that a checks run against it passed there
// too, at head exactly: what a request for a story in review is answered
// from the record for (ADR-0292 §1).
func passedAt(r kluster.Record, head string, checksOn bool) bool {
	if r.Commit != head || r.Outcome != kluster.OutcomePassed {
		return false
	}
	return !checksOn || r.ChecksCommit == head && r.ChecksOutcome == kluster.OutcomePassed
}

// lookReview queues a deploy of each story in review whose namespace's record
// did not pass at its branch's head, its checks too while checksOn, and that
// the pass did not try at that head, or that is back in review after a send
// back, at whatever head; and names them. A deploy the agent asked for that
// passed at the head stands as the review's. At the head is as the
// acceptance gate judges it with m's layout (ADR-0292 §3): the head, or a
// commit followed only by changes under the design, docs, and wip folders,
// so the issues a close-out commits do not deploy the story again. A record
// or a diff that cannot be read is not at the head: it is logged, and the
// story deployed.
func (p *klusterPass) lookReview(ctx context.Context, items []*workitem.Item, root string, m manifest.Manifest, checksOn bool) []string {
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
		head, err := kluster.StoryHead(git, root, it.ID)
		if err != nil {
			p.o.Logger.Debug("story in review has no branch to deploy", "component", "serve", "root", p.e.Root, "story", it.ID)
			continue
		}
		passed, _, err := kluster.AcceptGate(git, root, m, it.ID, head, checksOn)
		if err != nil {
			p.o.Logger.Warn("kluster record not judged at the story's head, so it is deployed again", "component", "serve", "root", p.e.Root, "story", it.ID, "head", head, "err", err.Error())
		}
		if passed {
			continue
		}
		// back from a send back, it is deployed again in place, at a new head or
		// the same one: its deploy writes what was tried afresh
		if t, tried := readKlusterTried(root, ns); tried && t.SentBack == nil && t.Head == head {
			continue
		}
		if p.enqueue(klusterJob{kind: klusterDeploy, story: it.ID, head: head}) {
			queued = append(queued, klusterDeploy+" "+it.ID)
		}
	}
	return queued
}

// lookLeft queues a teardown of each story namespace with a record whose
// story is neither in review nor in progress, whoever moved it there, unless
// a teardown already failed while the story stood where it stands, and
// forgets what the pass tried for one with no record whose story is neither;
// it names the teardowns queued.
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
		if klusterHolds(it) {
			if tried && t.Teardown != "" {
				// back in review or in progress: a teardown that failed before is
				// tried again once it leaves
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

// klusterHolds is whether story it keeps its namespace: it is in progress or
// in review, not archived (ADR-0292 §3).
func klusterHolds(it *workitem.Item) bool {
	return !it.Archived && (it.Status == workitem.InProgress || it.Status == workitem.Review)
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

// enqueue adds j, a job a look queues, to the queue unless the same job runs
// or waits already, whoever asked for it; a deploy waiting for the same story
// at another head takes the newer head in its place, keeping its id and who
// asked. A promotion running does not keep another from waiting, since it may
// have resolved main's head before the acceptance that queues it. It reports
// whether the queue changed.
func (p *klusterPass) enqueue(j klusterJob) bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	if r := p.running; r != nil && j.kind != klusterPromote && r.kind == j.kind && r.story == j.story && r.head == j.head {
		return false
	}
	for i, q := range p.queue {
		if q.kind == j.kind && q.story == j.story {
			if q.head == j.head {
				return false
			}
			p.queue[i].head = j.head
			return true
		}
	}
	j.queued = p.o.Now().UTC()
	if j.story != "" {
		j.id = kluster.JobID(j.story, j.kind, j.queued)
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
// since it was queued, or, for a deploy or a promotion, its deploy command;
// a job a story's agent asked for then ends refused. Its end is kept in the
// job's state, and told in a kluster event, when the agent asked for it.
func (p *klusterPass) run(ctx context.Context, j klusterJob) {
	repo, err := workitem.Open(p.e.Root)
	if err != nil {
		p.o.Logger.Warn("kluster job not run", "component", "serve", "root", p.e.Root, "kind", j.kind, "story", j.story, "err", err.Error())
		if j.asked() {
			p.endAsked(nil, j, refusedEnd(fmt.Sprintf("flai serve could not read the project to run it (%v); ask again", err)))
		}
		return
	}
	end := p.runIn(ctx, repo, j)
	if j.asked() {
		p.endAsked(repo, j, end)
	}
}

// runIn runs j in the project repo and answers how it ended.
func (p *klusterPass) runIn(ctx context.Context, repo *workitem.Repo, j klusterJob) klusterEnd {
	switch {
	case repo.Manifest.Kluster == nil:
		return refusedEnd("the manifest's kluster: block was removed while the job waited; asking again is refused until flai serve enable kluster, the operator's, writes one")
	case !p.on(repo.Manifest):
		return refusedEnd("the kluster host action was turned off while the job waited; turning it on again with flai serve enable kluster is the operator's")
	case j.kind != klusterTeardown && !klusterDeployable(repo.Manifest):
		return refusedEnd("the manifest's kluster: block names no deploy command now, so deploys are not set up yet")
	}
	switch j.kind {
	case klusterDeploy:
		it, err := repo.Get(j.story)
		if err != nil {
			return refusedEnd(fmt.Sprintf("flai serve could not read %s to deploy it (%v)", j.story, err))
		}
		if !klusterHolds(it) || (!j.asked() && it.Status != workitem.Review) {
			// left review since it was queued: the next look tears it down
			return refusedEnd(fmt.Sprintf("%s is %s now, and flai serve deploys a story only while it is in progress or in review", j.story, it.Status))
		}
		return p.deploy(ctx, repo, j, it.Status)
	case klusterChecks:
		it, err := repo.Get(j.story)
		if err != nil {
			return refusedEnd(fmt.Sprintf("flai serve could not read %s to run its checks (%v)", j.story, err))
		}
		if !klusterHolds(it) {
			return refusedEnd(fmt.Sprintf("%s is %s now, and flai serve runs a story's checks against its namespace only while it is in progress or in review", j.story, it.Status))
		}
		return p.checksJob(ctx, repo, j)
	case klusterTeardown:
		return p.teardown(ctx, repo, j)
	case klusterPromote:
		p.promote(ctx, repo)
	}
	return klusterEnd{}
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

// klusterEnd is how a story's job ended, for its state and its kluster event:
// its outcome, empty when flai serve stopped under it and the next start ends
// it interrupted; the commit it built or found; why it failed or was refused;
// the steps and log of the namespace's record; and the checks run it ran,
// nil when none did.
type klusterEnd struct {
	outcome string
	commit  string
	fail    klusterFailure
	steps   []kluster.Step
	log     string
	checks  *kluster.ChecksResult
}

// refusedEnd is a job that ran nothing, for why.
func refusedEnd(why string) klusterEnd {
	return klusterEnd{outcome: kluster.OutcomeRefused, fail: klusterFailure{exit: -1, err: why}}
}

// deploy deploys j's story, whose status is status, at j.head, bringing the
// cluster up first when it is absent, then runs its checks; it settles the
// story's failure thread for a job the looks queued, and opens none for one
// its agent asked for (ADR-0292 §5).
func (p *klusterPass) deploy(ctx context.Context, repo *workitem.Repo, j klusterJob, status string) klusterEnd {
	root := mainCheckout(repo)
	ns, _ := kluster.StoryNamespace(j.story)
	p.startAsked(j)
	p.journalJob(j, "started", fmt.Sprintf("deploy of %s into %s at %s", j.story, ns, short(j.head)))
	p.o.Logger.Info("kluster deploy started", "component", "serve", "root", p.e.Root, "story", j.story, "namespace", ns, "commit", j.head, "by", j.asker())
	r, fail, err := p.deployOnce(ctx, repo, j, status)
	end := klusterEnd{outcome: kluster.OutcomePassed, commit: j.head, fail: klusterFailure{exit: -1}, steps: r.Steps, log: kluster.LogDir(root, ns)}
	switch {
	case err != nil && ctx.Err() != nil:
		// flai serve stops: the next look queues it again, and the next start
		// ends an asked-for one interrupted
		p.journalJob(j, "failed", fmt.Sprintf("deploy of %s into %s at %s not run, as flai serve stops: %v", j.story, ns, short(j.head), err))
		p.o.Logger.Info("kluster deploy not run", "component", "serve", "root", p.e.Root, "story", j.story, "namespace", ns, "commit", j.head, "err", err.Error())
		return klusterEnd{}
	case errors.Is(err, kluster.ErrBusy):
		// another job holds the lock: the next look queues it again, and the
		// agent asks again
		p.journalJob(j, "failed", fmt.Sprintf("deploy of %s into %s at %s not run, another job holding the lock: %v", j.story, ns, short(j.head), err))
		p.o.Logger.Info("kluster deploy not run", "component", "serve", "root", p.e.Root, "story", j.story, "namespace", ns, "commit", j.head, "err", err.Error())
		return refusedEnd(err.Error())
	}
	if !j.asked() {
		tried := klusterTried{Story: j.story, Head: j.head, Outcome: kluster.OutcomePassed, At: p.o.Now().UTC()}
		if err != nil {
			tried.Outcome = kluster.OutcomeFailed
		}
		if werr := writeKlusterTried(root, ns, tried); werr != nil {
			p.o.Logger.Warn("kluster deploy not recorded", "component", "serve", "root", p.e.Root, "story", j.story, "err", werr.Error())
		}
	}
	if err != nil {
		p.journalJob(j, "failed", fmt.Sprintf("deploy of %s into %s at %s failed %s: %s", j.story, ns, short(j.head), fail.where(), fail.err))
		p.o.Logger.Warn("kluster deploy failed", "component", "serve", "root", p.e.Root, "story", j.story, "namespace", ns, "commit", j.head, "step", fail.step, "exit_code", fail.exit, "log", fail.log, "by", j.asker())
		if !j.asked() {
			p.failed(repo, j, ns, fail)
		}
		end.outcome, end.fail, end.log = kluster.OutcomeFailed, fail, cmp.Or(fail.log, end.log)
		return end
	}
	p.journalJob(j, "done", fmt.Sprintf("deployed %s into %s at %s", j.story, ns, short(j.head)))
	p.o.Logger.Info("kluster deploy passed", "component", "serve", "root", p.e.Root, "story", j.story, "namespace", ns, "commit", j.head, "by", j.asker())
	p.passed(repo, j, ns)
	run := p.checks(ctx, repo, j, ns)
	if run == nil {
		return end
	}
	end.checks = checksResult(*run, j.head)
	if run.Outcome != kluster.OutcomePassed {
		end.outcome, end.log = kluster.OutcomeFailed, run.Log
		end.fail = klusterFailure{step: "checks", exit: failedExit(*run), log: run.Log, err: checksEnded(*run)}
	}
	return end
}

// checksJob runs the checks a story's agent asked for against j's story's
// namespace as it is deployed (ADR-0292 §1, §2): the review's checks run,
// with the namespace's variables and flai's whole kubeconfig. It runs nothing
// and ends refused when the checks host action went off, no checks are named,
// or the namespace is no longer deployed at the story branch's head with a
// deploy that passed. A run that does not pass ends the job failed at the
// check it stopped at, with that check's exit code, and opens no thread.
func (p *klusterPass) checksJob(ctx context.Context, repo *workitem.Repo, j klusterJob) klusterEnd {
	ns, _ := kluster.StoryNamespace(j.story)
	cfg, env, why := p.checksReady(repo, j.story)
	if why != "" {
		return refusedEnd(why + ", so flai serve ran no checks for " + j.story + ": " + deployFirst(j.story))
	}
	p.startAsked(j)
	run, err := p.runChecks(ctx, repo, j, ns, cfg, env)
	switch {
	case ctx.Err() != nil && (err != nil || run.Outcome == "cancelled"):
		return klusterEnd{} // flai serve stops: the next start ends it interrupted
	case err != nil:
		return refusedEnd(fmt.Sprintf("flai serve could not start the checks of %s against its namespace %s: %v. Ask again once that is put right", j.story, ns, err))
	}
	commit := cmp.Or(run.Commit, j.head)
	end := klusterEnd{outcome: kluster.OutcomePassed, commit: commit, fail: klusterFailure{exit: -1}, log: run.Log, checks: checksResult(*run, commit)}
	if run.Outcome != kluster.OutcomePassed {
		step := "checks"
		if n := len(run.Steps); n > 0 {
			step = run.Steps[n-1].Name
		}
		end.outcome = kluster.OutcomeFailed
		end.fail = klusterFailure{step: step, exit: failedExit(*run), log: run.Log, err: checksEnded(*run)}
	}
	return end
}

// deployFirst says what a story's agent does when its namespace is not
// deployed at its head with a deploy that passed.
func deployFirst(story string) string {
	return fmt.Sprintf("deploy it first with flai kluster deploy %s --job, or the tool kluster_deploy, which runs the checks once the deploy passes, and ask for a checks run again once its kluster event says it passed", story)
}

// checksEnded says how a checks run that did not pass ended, and where.
func checksEnded(run ChecksRun) string {
	return fmt.Sprintf("the checks run against the namespace ended %s, at %s", run.Outcome, failedCheck(run))
}

// failedExit is the exit code of the check a run that did not pass stopped
// at, -1 when it gave none.
func failedExit(run ChecksRun) int {
	if n := len(run.Steps); n > 0 && run.Steps[n-1].Exit != nil {
		return *run.Steps[n-1].Exit
	}
	return -1
}

// checksResult is run, a checks run against a namespace deployed at commit,
// as a job's state keeps it: each check that ran, passed but for the last of
// a run that did not pass, which ended as the run did.
func checksResult(run ChecksRun, commit string) *kluster.ChecksResult {
	out := &kluster.ChecksResult{Outcome: run.Outcome, Commit: commit, Log: run.Log}
	for i, step := range run.Steps {
		c := kluster.Check{Name: step.Name, Outcome: kluster.OutcomePassed, ExitCode: step.Exit, Log: run.Log}
		if i == len(run.Steps)-1 && run.Outcome != kluster.OutcomePassed {
			c.Outcome = run.Outcome
		}
		out.Checks = append(out.Checks, c)
	}
	return out
}

// KlusterEnv is the environment a checks run for story gets (ADR-0254 §8):
// the variables of its namespace, as Cluster.Env gives them, while the
// kluster host action is on for the project, its manifest has a kluster:
// block, and the namespace's record passed at the head of the story's
// branch, which git reads in the main checkout. Otherwise it is nil, and so
// is ns, the namespace. Only what cannot be read is an error.
func KlusterEnv(repo *workitem.Repo, story string, git execx.Runner, enabled func(action, root string) bool) (ns string, env []string, err error) {
	ns, env, _, err = klusterEnvWhy(repo, story, git, enabled)
	return ns, env, err
}

// klusterEnvWhy is KlusterEnv with why, when it gives no variables, in a
// refusal's words: the action off, no block, no record, a record that did
// not pass, no branch, or a record at another commit. The review's checks
// run passes over a story with none silently; a checks run its agent asks
// for is refused with why.
func klusterEnvWhy(repo *workitem.Repo, story string, git execx.Runner, enabled func(action, root string) bool) (ns string, env []string, why string, err error) {
	root, m := mainCheckout(repo), repo.Manifest
	switch {
	case m.Kluster == nil:
		return "", nil, "the manifest has no kluster: block", nil
	case enabled == nil || !enabled(hostapi.ActionKluster, root):
		return "", nil, "the kluster host action is off for this project", nil
	}
	story = workitem.CanonicalID(story)
	ns, err = kluster.StoryNamespace(story)
	if err != nil {
		return "", nil, "", err
	}
	r, ok, err := kluster.ReadRecord(root, ns)
	if err != nil {
		return "", nil, "", fmt.Errorf("read the kluster record of %s's namespace %s, to give its checks the namespace: %w", story, ns, err)
	}
	head, _ := kluster.StoryHead(git, root, story)
	switch {
	case !ok:
		return "", nil, fmt.Sprintf("%s's namespace %s was never deployed: it has no record", story, ns), nil
	case r.Outcome != kluster.OutcomePassed:
		return "", nil, fmt.Sprintf("the last deploy of %s's namespace %s, at %s, ended %s, not passed", story, ns, orUnresolved(short(r.Commit)), cmp.Or(r.Outcome, "with no outcome")), nil
	case head == "":
		return "", nil, fmt.Sprintf("%s has no branch story/%s whose head its namespace %s could be deployed at", story, story, ns), nil
	case head != r.Commit:
		return "", nil, fmt.Sprintf("%s's namespace %s was deployed at %s, not at the story branch's head %s", story, ns, short(r.Commit), short(head)), nil
	}
	c, err := kluster.New(root, m, kluster.System{})
	if err != nil {
		return "", nil, "", err
	}
	t, err := kluster.StoryTarget(story, r.Checkout, r.Commit)
	if err != nil {
		return "", nil, "", err
	}
	return ns, c.Env(t), "", nil
}

// checks runs j's story's checks against its namespace once its deploy has
// passed, in the same job, while the checks host action is on and the
// project names checks. A run that ends keeps its outcome at j.head in the
// namespace's record, which an acceptance reads (ADR-0254 §5a); one that
// does not pass opens one thread on the story and sends it back to
// in-progress, as a failed deploy does, unless its agent asked for the job
// (ADR-0292 §5), and one that passes resolves it. It answers the run, or nil
// when none ended.
func (p *klusterPass) checks(ctx context.Context, repo *workitem.Repo, j klusterJob, ns string) *ChecksRun {
	cfg, env, why := p.checksReady(repo, j.story)
	if why != "" {
		// the checks off or none named, or the head moved since the deploy:
		// the next look deploys it again
		p.o.Logger.Debug("kluster checks not run", "component", "serve", "root", p.e.Root, "story", j.story, "why", why)
		return nil
	}
	run, _ := p.runChecks(ctx, repo, j, ns, cfg, env)
	return run
}

// checksReady is the checks configuration and the namespace's variables a
// checks run for story gets, or why it gets none, in a refusal's words that
// say what to do: the checks host action off, no checks named, a
// configuration or a record that cannot be read, or the namespace not
// deployed at the story branch's head with a deploy that passed.
func (p *klusterPass) checksReady(repo *workitem.Repo, story string) (cfg ChecksConfig, env []string, why string) {
	if cfg, why = p.checksConfig(repo); why != "" {
		return ChecksConfig{}, nil, why
	}
	_, env, why, err := klusterEnvWhy(repo, story, gitOf(p.o), p.o.Host.Enabled)
	switch {
	case err != nil:
		p.o.Logger.Warn("kluster checks not run", "component", "serve", "root", p.e.Root, "story", story, "err", err.Error())
		return ChecksConfig{}, nil, fmt.Sprintf("flai serve could not give the checks %s's namespace (%v)", story, err)
	case env == nil:
		return ChecksConfig{}, nil, why
	}
	return cfg, env, ""
}

// checksConfig is the project's checks configuration in repo, or why flai
// serve runs no checks for it, in a refusal's words naming what the operator
// does: the checks host action off, no checks named, or checks that cannot
// be read.
func (p *klusterPass) checksConfig(repo *workitem.Repo) (ChecksConfig, string) {
	const off = "the checks host action is off for this project: turning it on with flai serve enable checks is the operator's"
	if p.o.Checks == nil {
		return ChecksConfig{}, off
	}
	cfg, err := p.o.Checks(repo)
	switch {
	case err != nil:
		return ChecksConfig{}, fmt.Sprintf("flai serve could not read the project's checks (%v)", err)
	case !cfg.Enabled:
		return ChecksConfig{}, off
	case len(cfg.Commands) == 0:
		return ChecksConfig{}, "no checks are named for this project: naming them with flai serve checks set, or checks: in the manifest, is the operator's"
	}
	return cfg, ""
}

// runChecks runs cfg's checks for j's story with env, its namespace ns's
// variables, for who asked for j, in the story's container when its agent
// runs in one. It keeps the outcome and the log in the namespace's record at
// the commit the variables name, resolves the story's checks thread on a
// pass, and on a failure sends the story back with a thread unless its agent
// asked for j. It answers the run, or why it could not start one.
func (p *klusterPass) runChecks(ctx context.Context, repo *workitem.Repo, j klusterJob, ns string, cfg ChecksConfig, env []string) (*ChecksRun, error) {
	commit := cmp.Or(envValue(env, "FLAI_KLUSTER_COMMIT"), j.head)
	p.journalJob(j, "started", fmt.Sprintf("checks of %s against %s at %s", j.story, ns, short(commit)))
	p.o.Logger.Info("kluster checks started", "component", "serve", "root", p.e.Root, "story", j.story, "namespace", ns, "commit", commit, "by", j.asker())
	agents := AgentConfig{}
	if p.o.Agent != nil {
		agents = p.o.Agent(p.e.Root)
	}
	in := checksContainerOf(repo, j.story, agents, dockerOf(p.o), p.o.imageJobs, p.e.Key)
	run, err := p.o.Dir.RunChecksAs(ctx, j.asker(), j.story, mainCheckout(repo), repo.WorktreePath(j.story), cfg, env, in, p.o.Now)
	if err != nil {
		p.journalJob(j, "failed", fmt.Sprintf("checks of %s against %s at %s not run: %v", j.story, ns, short(commit), err))
		p.o.Logger.Info("kluster checks not run", "component", "serve", "root", p.e.Root, "story", j.story, "namespace", ns, "err", err.Error())
		return nil, err
	}
	if err := kluster.RecordChecks(mainCheckout(repo), ns, commit, run.Outcome, run.Log, run.AskedBy); err != nil {
		p.o.Logger.Warn("kluster checks outcome not recorded", "component", "serve", "root", p.e.Root, "story", j.story, "namespace", ns, "commit", commit, "outcome", run.Outcome, "err", err.Error())
	}
	if run.Outcome == kluster.OutcomePassed {
		p.journalJob(j, "done", fmt.Sprintf("checks of %s against %s at %s passed", j.story, ns, short(commit)))
		p.o.Logger.Info("kluster checks passed", "component", "serve", "root", p.e.Root, "story", j.story, "namespace", ns, "commit", commit, "by", j.asker())
		p.clear(repo, j.story, checksFailureTitle(j.story), fmt.Sprintf("flai serve's checks of %s against its namespace %s passed at %s.", j.story, ns, short(commit)))
		return &run, nil
	}
	check := failedCheck(run)
	p.journalJob(j, "failed", fmt.Sprintf("checks of %s against %s at %s %s at %s", j.story, ns, short(commit), run.Outcome, check))
	p.o.Logger.Warn("kluster checks failed", "component", "serve", "root", p.e.Root, "story", j.story, "namespace", ns, "commit", commit, "outcome", run.Outcome, "check", check, "log", run.Log, "by", j.asker())
	if !j.asked() {
		reason := fmt.Sprintf("flai serve's checks of %s against its kluster namespace ended %s, at %s, at commit %s", j.story, run.Outcome, check, short(j.head))
		p.sendBack(repo, j, ns, checksFailureTitle(j.story), reason, func(unsent string) string { return checksFailureText(j, ns, run, unsent) })
	}
	return &run, nil
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
// the commit, and the log, and says what flai serve did with the story, as
// heldText does for unsent, and what happens next.
func checksFailureText(j klusterJob, ns string, run ChecksRun, unsent string) string {
	log := "no log"
	if run.Log != "" {
		log = "`" + run.Log + "`"
	}
	return fmt.Sprintf("flai serve's checks of %s against its namespace `%s` at %s ended %s, at %s. Their log is %s.\n\n%s flai checks status %s shows the run.",
		j.story, ns, j.head, run.Outcome, failedCheck(run), log, heldText(j.story, unsent, "a checks run against it passes"), j.story)
}

// heldText says what flai serve did with story after a failure in its
// namespace: sent it back to in-progress, keeping its namespace, or, with
// unsent, why it did not; and when it resolves the thread, once what passes
// says passes.
func heldText(story, unsent, passes string) string {
	if unsent != "" {
		return fmt.Sprintf("flai serve did not send %[1]s back to in-progress: %[2]s. It deploys again at the story's next head in review, and resolves this thread once %[3]s.", story, unsent, passes)
	}
	return fmt.Sprintf("flai serve sent %[1]s back to in-progress, and keeps its namespace as the failure left it, so that its agent and the operator can look at what failed. The agent fixes the change or, when the cause is outside the story, says so on this thread. Once %[1]s is back in review, at a new head or the same one, flai serve deploys it again in place, and resolves this thread once %[2]s. Cancelling or accepting %[1]s tears the namespace down.", story, passes)
}

// deployOnce brings the cluster up when it is absent, then deploys j's
// story, whose status is status, from its worktree for who asked, and
// answers the namespace's record; on a failure it says which step failed and
// where its log is.
func (p *klusterPass) deployOnce(ctx context.Context, repo *workitem.Repo, j klusterJob, status string) (kluster.Record, klusterFailure, error) {
	root, m := mainCheckout(repo), repo.Manifest
	dir := repo.WorktreePath(j.story)
	if info, err := os.Stat(dir); err != nil || !info.IsDir() {
		err := fmt.Errorf("%s has no worktree at %s to build and deploy from; open it with flai stream open %s", j.story, dir, j.story)
		return kluster.Record{}, klusterFailure{exit: -1, err: err.Error()}, err
	}
	absent, err := p.jobs.Absent(ctx, root, m)
	if err != nil {
		return kluster.Record{}, klusterFailure{step: "up", exit: -1, err: err.Error()}, err
	}
	if absent {
		if err := p.up(ctx, root, m); err != nil {
			return kluster.Record{}, klusterFailure{step: "up", exit: kluster.ExitCode(err), log: klusterUpLog(root), err: err.Error()}, err
		}
	}
	r, err := p.jobs.Deploy(ctx, root, m, j.asker(), j.story, status, dir, j.head)
	if err == nil {
		return r, klusterFailure{exit: -1}, nil
	}
	return r, failureOf(r, err), err
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
// the one open, and sends the story back to in-progress.
func (p *klusterPass) failed(repo *workitem.Repo, j klusterJob, ns string, f klusterFailure) {
	reason := fmt.Sprintf("flai serve's kluster %s of %s failed %s, at commit %s", j.kind, j.story, f.where(), short(j.head))
	p.sendBack(repo, j, ns, klusterFailureTitle(j.story), reason, func(unsent string) string { return klusterFailureText(j, ns, f, unsent) })
}

// sendBack opens the thread titled title on j's story, or adds to the one
// open, with what text says for why the story was not sent back, empty when
// it was; then, once the thread is open, moves the story from review to
// in-progress as flai for reason, and keeps the send back with what the
// pass tried for namespace ns (ADR-0254 §5a). It moves nothing when the
// kluster action went off meanwhile, or the story is no longer in review at
// j.head, and the thread says so.
func (p *klusterPass) sendBack(repo *workitem.Repo, j klusterJob, ns, title, reason string, text func(unsent string) string) {
	fresh, it, unsent := p.inReviewAt(j)
	th := p.report(repo, j.story, title, text(unsent))
	if unsent != "" {
		p.o.Logger.Info("story not sent back to in-progress after its kluster failure", "component", "serve", "root", p.e.Root, "story", j.story, "commit", j.head, "why", unsent)
		return
	}
	if th == nil {
		return // no thread to send it back on: report warned of it
	}
	now := p.o.Now()
	if _, err := fresh.Transition(it, workitem.InProgress, sendBackBy, reason, now); err != nil {
		p.o.Logger.Warn("story not sent back to in-progress after its kluster failure", "component", "serve", "root", p.e.Root, "story", j.story, "thread", th.ID, "err", err.Error())
		return
	}
	root := mainCheckout(fresh)
	t, _ := readKlusterTried(root, ns)
	t.Story = j.story
	t.SentBack = &klusterSentBack{Head: j.head, Thread: th.ID, Reason: reason, At: now.UTC().Truncate(time.Second)}
	if err := writeKlusterTried(root, ns, t); err != nil {
		p.o.Logger.Warn("kluster send back not recorded, so the story's namespace is not kept and its agent not started for it", "component", "serve", "root", p.e.Root, "story", j.story, "thread", th.ID, "err", err.Error())
	}
	p.journal("done", fmt.Sprintf("sent %s back to in-progress on %s: %s", j.story, th.ID, reason))
	p.o.Logger.Info("story sent back to in-progress after its kluster failure", "component", "serve", "root", p.e.Root, "story", j.story, "thread", th.ID, "reason", reason)
}

// inReviewAt reads the project and j's story again, and answers them with an
// empty why when the story may be sent back: the kluster action is still on,
// and the story still in review at j.head. Otherwise why says, in a thread's
// words, why it is not sent back.
func (p *klusterPass) inReviewAt(j klusterJob) (repo *workitem.Repo, it *workitem.Item, why string) {
	repo, err := workitem.Open(p.e.Root)
	if err != nil {
		return nil, nil, "flai serve could not read the project again to move it (" + err.Error() + ")"
	}
	if !p.on(repo.Manifest) {
		return nil, nil, "the kluster host action was turned off meanwhile"
	}
	if it, err = repo.Get(j.story); err != nil {
		return nil, nil, "flai serve could not read the story again to move it (" + err.Error() + ")"
	}
	if it.Status != workitem.Review || it.Archived {
		return nil, nil, fmt.Sprintf("it is %s now, no longer in review", it.Status)
	}
	if head, err := kluster.StoryHead(gitOf(p.o), mainCheckout(repo), j.story); err != nil || head != j.head {
		return nil, nil, fmt.Sprintf("its branch is no longer at %s", short(j.head))
	}
	return repo, it, ""
}

// passed resolves the story's deploy failure thread, if one is open.
func (p *klusterPass) passed(repo *workitem.Repo, j klusterJob, ns string) {
	p.clear(repo, j.story, klusterFailureTitle(j.story), fmt.Sprintf("flai serve deployed %s into its namespace %s at %s.", j.story, ns, short(j.head)))
}

// teardown tears j's story down, from its worktree or the main checkout when
// it is gone, unless it is back in review, or back in progress for a
// teardown a look queued, and settles its teardown failure thread; a story
// done or cancelled has its deploy and checks failure threads resolved too. A
// failure is kept with the story's column, so that it is tried again only
// once the story moves; one its agent asked for opens no thread (ADR-0292
// §5).
func (p *klusterPass) teardown(ctx context.Context, repo *workitem.Repo, j klusterJob) klusterEnd {
	it, err := repo.Get(j.story)
	if err != nil {
		p.o.Logger.Warn("kluster teardown not run", "component", "serve", "root", p.e.Root, "story", j.story, "err", err.Error())
		return refusedEnd(fmt.Sprintf("flai serve could not read %s to tear it down (%v)", j.story, err))
	}
	if it.Status == workitem.Review && !it.Archived {
		// back in review: its deploy updates the namespace in place
		return refusedEnd(fmt.Sprintf("%s moved to review while the teardown waited, and a namespace in review is the review's: flai serve tears it down once the story leaves review for ready or backlog, is cancelled, or is accepted", j.story))
	}
	if !j.asked() && klusterHolds(it) {
		return klusterEnd{} // back in progress since it was queued: it keeps its namespace
	}
	root := mainCheckout(repo)
	ns, _ := kluster.StoryNamespace(j.story)
	dir := ""
	if info, err := os.Stat(repo.WorktreePath(j.story)); err == nil && info.IsDir() {
		dir = repo.WorktreePath(j.story)
	}
	p.startAsked(j)
	p.journalJob(j, "started", fmt.Sprintf("teardown of %s's namespace %s, %s", j.story, ns, it.Status))
	p.o.Logger.Info("kluster teardown started", "component", "serve", "root", p.e.Root, "story", j.story, "namespace", ns, "status", it.Status, "by", j.asker())
	r, err := p.jobs.Teardown(ctx, root, repo.Manifest, j.asker(), j.story, dir)
	end := klusterEnd{outcome: kluster.OutcomePassed, commit: r.Commit, fail: klusterFailure{exit: -1}, steps: r.Steps}
	switch {
	case errors.Is(err, kluster.ErrNotDeployed):
		p.forget(root, ns)
		p.journalJob(j, "done", fmt.Sprintf("teardown of %s's namespace %s: nothing was deployed", j.story, ns))
		p.o.Logger.Info("kluster teardown found nothing deployed", "component", "serve", "root", p.e.Root, "story", j.story, "namespace", ns)
	case err != nil && ctx.Err() != nil:
		p.journalJob(j, "failed", fmt.Sprintf("teardown of %s's namespace %s not run, as flai serve stops: %v", j.story, ns, err))
		p.o.Logger.Info("kluster teardown not run", "component", "serve", "root", p.e.Root, "story", j.story, "namespace", ns, "err", err.Error())
		return klusterEnd{}
	case errors.Is(err, kluster.ErrBusy):
		p.journalJob(j, "failed", fmt.Sprintf("teardown of %s's namespace %s not run, another job holding the lock: %v", j.story, ns, err))
		p.o.Logger.Info("kluster teardown not run", "component", "serve", "root", p.e.Root, "story", j.story, "namespace", ns, "err", err.Error())
		return refusedEnd(err.Error())
	case err != nil:
		f := failureOf(r, err)
		t, _ := readKlusterTried(root, ns)
		t.Story, t.Teardown, t.At = j.story, it.Status, p.o.Now().UTC()
		if werr := writeKlusterTried(root, ns, t); werr != nil {
			p.o.Logger.Warn("kluster teardown not recorded", "component", "serve", "root", p.e.Root, "story", j.story, "err", werr.Error())
		}
		p.journalJob(j, "failed", fmt.Sprintf("teardown of %s's namespace %s failed %s: %s", j.story, ns, f.where(), f.err))
		p.o.Logger.Warn("kluster teardown failed", "component", "serve", "root", p.e.Root, "story", j.story, "namespace", ns, "step", f.step, "exit_code", f.exit, "log", f.log, "by", j.asker())
		if !j.asked() {
			p.report(repo, teardownAnchor(it), teardownFailureTitle(j.story), teardownFailureText(it, ns, r.Commit, f))
		}
		end.outcome, end.fail, end.log = kluster.OutcomeFailed, f, f.log
	default:
		p.forget(root, ns)
		p.journalJob(j, "done", fmt.Sprintf("tore down %s's namespace %s", j.story, ns))
		p.o.Logger.Info("kluster teardown passed", "component", "serve", "root", p.e.Root, "story", j.story, "namespace", ns, "by", j.asker())
		why := fmt.Sprintf("flai serve tore down %s's namespace %s.", j.story, ns)
		p.clear(repo, j.story, teardownFailureTitle(j.story), why)
		p.clear(repo, manifest.File, teardownFailureTitle(j.story), why)
		if it.Status == workitem.Done || it.Status == workitem.Cancelled {
			why = fmt.Sprintf("%s is %s, and flai serve tore down its namespace %s.", j.story, it.Status, ns)
			p.clear(repo, j.story, klusterFailureTitle(j.story), why)
			p.clear(repo, j.story, checksFailureTitle(j.story), why)
		}
	}
	return end
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
// to the one open, and answers the thread; nil when none could be opened.
func (p *klusterPass) report(repo *workitem.Repo, on, title, text string) *threads.Thread {
	now := p.o.Now()
	if th := klusterThread(repo, on, title); th != nil {
		if _, err := threads.Reply(repo, th.ID, forgeAuthor, text, now); err != nil {
			p.o.Logger.Warn("kluster failure not written on its thread", "component", "serve", "root", p.e.Root, "on", on, "thread", th.ID, "err", err.Error())
		}
		return th
	}
	th, err := threads.New(repo, threads.NewOptions{Title: title, On: on, Author: forgeAuthor, Text: text, Now: now})
	if err != nil {
		p.o.Logger.Warn("kluster failure thread not opened", "component", "serve", "root", p.e.Root, "on", on, "err", err.Error())
		return nil
	}
	return th
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

// journal keeps one entry of the pass's in the host's journal, by flai serve.
func (p *klusterPass) journal(outcome, detail string) {
	p.journalBy(forgeAuthor, outcome, detail)
}

// journalJob keeps one entry of job j's in the host's journal: by the agent
// that asked for it, its detail naming the agent, or by flai serve for a job
// a look queued.
func (p *klusterPass) journalJob(j klusterJob, outcome, detail string) {
	if j.asked() {
		detail += ", asked by " + j.by
	}
	p.journalBy(j.asker(), outcome, detail)
}

func (p *klusterPass) journalBy(by, outcome, detail string) {
	if p.o.Host.Record == nil {
		return
	}
	p.o.Host.Record(hostapi.Entry{At: p.o.Now().UTC().Format(time.RFC3339), Action: hostapi.ActionKluster, Method: klusterMethod, Project: p.e.Key, Root: p.e.Root,
		By: by, Outcome: outcome, Detail: detail})
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
// log, and the error, and says what flai serve did with the story, as
// heldText does for unsent, and what happens next.
func klusterFailureText(j klusterJob, ns string, f klusterFailure, unsent string) string {
	exit, log := f.said()
	return fmt.Sprintf("flai serve's %s of %s into the namespace `%s` at %s failed %s, with %s. Its log is %s.\n\n```text\n%s\n```\n\n%s flai kluster deployments shows the namespace's record.",
		j.kind, j.story, ns, j.head, f.where(), exit, log, f.err, heldText(j.story, unsent, "a deploy passes"))
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
