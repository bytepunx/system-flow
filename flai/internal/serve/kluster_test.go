package serve

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/hostapi"
	"github.com/bytepunx/system-flow/flai/internal/kluster"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/threads"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// fakeKluster is a KlusterRunner that runs no cluster: it records each call,
// answers Absent with absent, fails Up with upErr, and fails each deploy with
// what fail answers for its story and head. A deploy writes the namespace's
// record, as the kluster package does, unless it was refused before its
// first step; a teardown removes it, or fails with teardownErr, and answers
// kluster.ErrNotDeployed without one; a promotion fails with promoteErr.
// While gate is not nil, each deploy says it started on started and waits for
// gate before it ends.
type fakeKluster struct {
	mu          sync.Mutex
	calls       []string // "absent", "up", "deploy S-nnnn <head>", "teardown S-nnnn", or "promote <branch>"
	absent      bool
	upErr       error
	fail        func(story, head string) error
	teardownErr error
	promoteErr  error
	gate        chan struct{}
	started     chan string
	inFlight    int
	most        int
	by          []string // who each deploy and teardown was for, in order
}

func (f *fakeKluster) note(call string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.calls = append(f.calls, call)
}

func (f *fakeKluster) Absent(context.Context, string, manifest.Manifest) (bool, error) {
	f.note("absent")
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.absent, nil
}

func (f *fakeKluster) Up(_ context.Context, _ string, _ manifest.Manifest, out io.Writer) error {
	f.note("up")
	fmt.Fprintln(out, "kluster up")
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.upErr == nil {
		f.absent = false
	}
	return f.upErr
}

func (f *fakeKluster) Deploy(_ context.Context, root string, _ manifest.Manifest, by, story, _, _, commit string) (kluster.Record, error) {
	f.note("deploy " + story + " " + commit)
	f.mu.Lock()
	f.by = append(f.by, by)
	f.inFlight++
	f.most = max(f.most, f.inFlight)
	gate, started, fail := f.gate, f.started, f.fail
	f.mu.Unlock()
	if started != nil {
		started <- story
	}
	if gate != nil {
		<-gate
	}
	f.mu.Lock()
	f.inFlight--
	f.mu.Unlock()
	r := kluster.Record{Story: story, Commit: commit, Outcome: kluster.OutcomePassed, AskedBy: by, StartedAt: time.Now().UTC()}
	var err error
	if fail != nil {
		err = fail(story, commit)
	}
	var se *kluster.StepError
	if err != nil && !errors.As(err, &se) {
		return kluster.Record{}, err // refused before its first step: no record
	}
	if err != nil {
		r.Outcome, r.FailedStep, r.ExitCode = kluster.OutcomeFailed, se.Step, se.ExitCode
	}
	writeFakeRecord(root, r)
	return r, err
}

// writeFakeRecord writes r as its namespace's record, as a deploy does.
func writeFakeRecord(root string, r kluster.Record) {
	ns, _ := kluster.StoryNamespace(r.Story)
	r.Namespace = ns
	data, _ := json.Marshal(r)
	_ = os.MkdirAll(filepath.Dir(kluster.RecordPath(root, ns)), 0o755)
	_ = os.WriteFile(kluster.RecordPath(root, ns), data, 0o644)
}

func (f *fakeKluster) Teardown(_ context.Context, root string, _ manifest.Manifest, by, story, _ string) (kluster.Record, error) {
	f.note("teardown " + story)
	f.mu.Lock()
	f.by = append(f.by, by)
	f.mu.Unlock()
	ns, _ := kluster.StoryNamespace(story)
	r, ok, err := kluster.ReadRecord(root, ns)
	if err != nil {
		return kluster.Record{}, err
	}
	if !ok {
		return kluster.Record{}, fmt.Errorf("%s: %w", story, kluster.ErrNotDeployed)
	}
	f.mu.Lock()
	terr := f.teardownErr
	f.mu.Unlock()
	if terr != nil {
		return r, terr
	}
	_ = os.Remove(kluster.RecordPath(root, ns))
	return r, nil
}

func (f *fakeKluster) Promote(_ context.Context, _ string, _ manifest.Manifest, branch string) (kluster.Record, error) {
	f.note("promote " + branch)
	f.mu.Lock()
	defer f.mu.Unlock()
	r := kluster.Record{Namespace: "t", Commit: klHeadC, Outcome: kluster.OutcomePassed}
	if f.promoteErr != nil {
		r.Outcome = kluster.OutcomeFailed
		return r, f.promoteErr
	}
	return r, nil
}

// since are the calls made after the first n.
func (f *fakeKluster) since(n int) []string {
	return f.all()[n:]
}

// deploys are the deploy calls made so far.
func (f *fakeKluster) deploys() []string {
	f.mu.Lock()
	defer f.mu.Unlock()
	var out []string
	for _, c := range f.calls {
		if strings.HasPrefix(c, "deploy ") {
			out = append(out, c)
		}
	}
	return out
}

func (f *fakeKluster) all() []string {
	f.mu.Lock()
	defer f.mu.Unlock()
	return append([]string(nil), f.calls...)
}

// klusterHeads answers git rev-parse of a story's branch from heads, by
// story, and git status in a folder from dirty, by folder: clean when it
// names none.
type klusterHeads struct {
	mu    sync.Mutex
	heads map[string]string
	dirty map[string]string
}

// setDirty makes git status in dir answer out.
func (g *klusterHeads) setDirty(dir, out string) {
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.dirty == nil {
		g.dirty = map[string]string{}
	}
	g.dirty[dir] = out
}

func (g *klusterHeads) set(story, head string) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.heads[story] = head
}

func (g *klusterHeads) Run(dir, name string, args ...string) (string, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	if name == "git" && len(args) > 0 && args[0] == "status" {
		return g.dirty[dir], nil
	}
	if name == "git" && len(args) == 3 && args[0] == "rev-parse" && args[1] == "--abbrev-ref" && args[2] == "HEAD" {
		return "main", nil
	}
	if name == "git" && len(args) > 0 && args[0] == "rev-parse" {
		ref := args[len(args)-1]
		story := strings.TrimSuffix(strings.TrimPrefix(ref, "refs/heads/story/"), "^{commit}")
		if head, ok := g.heads[story]; ok {
			return head + "\n", nil
		}
		return "", errors.New("exit status 1")
	}
	return "", fmt.Errorf("unexpected %s %v", name, args)
}

func (g *klusterHeads) RunInput(dir, name, _ string, args ...string) (string, error) {
	return g.Run(dir, name, args...)
}

func (g *klusterHeads) LookPath(name string) (string, error) { return name, nil }

// klusterLab is a project with a kluster: block, flai serve's kluster pass
// over it with the kluster host action on, a fake kluster, and fake heads.
type klusterLab struct {
	t       *testing.T
	root    string
	repo    *workitem.Repo
	on      bool
	git     *klusterHeads
	kluster *fakeKluster
	pass    *klusterPass
	opts    Options
	mu      sync.Mutex
	journal []hostapi.Entry
	clock   time.Time // a minute later at each round, as flai serve's looks are
}

const klusterBlock = "kluster:\n  config: integration/kluster.yaml\n  deploy: [scripts/kluster-deploy.sh]\n"

func newKlusterLab(t *testing.T, block string) *klusterLab {
	t.Helper()
	root := t.TempDir()
	m := "version: 1\nname: t\nkey: t\nlayout:\n  design: design\n  docs: docs\n  wip: wip\n" + block
	if err := os.WriteFile(filepath.Join(root, "system-flow.yaml"), []byte(m), 0o644); err != nil {
		t.Fatal(err)
	}
	for _, d := range []string{"wip/kanban/epics", "wip/kanban/stories", "wip/kanban/tasks", "wip/agents"} {
		if err := os.MkdirAll(filepath.Join(root, d), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	lab := &klusterLab{t: t, root: root, repo: repo, on: true, git: &klusterHeads{heads: map[string]string{}}, kluster: &fakeKluster{}}
	lab.clock = time.Now()
	o := Options{Logger: slog.New(slog.DiscardHandler), Git: lab.git, Kluster: lab.kluster,
		Now: func() time.Time {
			lab.mu.Lock()
			defer lab.mu.Unlock()
			return lab.clock
		},
		Host: hostapi.Host{
			Enabled: func(action, r string) bool {
				lab.mu.Lock()
				defer lab.mu.Unlock()
				return lab.on && action == hostapi.ActionKluster && r == root
			},
			Record: func(e hostapi.Entry) {
				lab.mu.Lock()
				defer lab.mu.Unlock()
				lab.journal = append(lab.journal, e)
			},
		},
	}
	lab.opts = o
	lab.pass = newKlusterPass(o, Entry{Key: "t", Name: "t", Root: root}, &sync.WaitGroup{})
	return lab
}

// restart is flai serve started again: a new pass, which counts no
// acceptance from before it.
func (lab *klusterLab) restart() {
	lab.mu.Lock()
	lab.clock = lab.clock.Add(time.Hour)
	lab.mu.Unlock()
	lab.pass = newKlusterPass(lab.opts, Entry{Key: "t", Name: "t", Root: lab.root}, &sync.WaitGroup{})
}

// accept moves a story to done a second after the lab's clock, so that the
// pass's next look sees it as an acceptance since its last.
func (lab *klusterLab) accept(id string) {
	lab.t.Helper()
	lab.mu.Lock()
	at := lab.clock.Add(time.Second)
	lab.mu.Unlock()
	items, err := lab.repo.List(false)
	if err != nil {
		lab.t.Fatal(err)
	}
	for _, task := range items {
		if task.Type != workitem.Task || task.Parent != id {
			continue
		}
		for _, to := range []string{workitem.Ready, workitem.InProgress, workitem.Done} {
			if _, err := lab.repo.Transition(task, to, "alex", "the test finishes it", at); err != nil {
				lab.t.Fatal(err)
			}
		}
	}
	it, err := lab.repo.Get(id)
	if err != nil {
		lab.t.Fatal(err)
	}
	data, _ := os.ReadFile(it.Path)
	if err := os.WriteFile(it.Path, []byte(strings.Replace(string(data), "- [ ] deploys", "- [x] deploys", 1)), 0o644); err != nil {
		lab.t.Fatal(err)
	}
	if it, err = lab.repo.Get(id); err != nil {
		lab.t.Fatal(err)
	}
	if _, err := lab.repo.Transition(it, workitem.Done, "alex", "accepted", at); err != nil {
		lab.t.Fatal(err)
	}
}

func (lab *klusterLab) setOn(on bool) {
	lab.mu.Lock()
	defer lab.mu.Unlock()
	lab.on = on
}

// inReview makes a story with a worktree folder, its branch at head, and
// moves it to review.
func (lab *klusterLab) inReview(title, head string) string {
	lab.t.Helper()
	epic, err := lab.repo.Create(workitem.NewOptions{Type: workitem.Epic, Title: "Epic " + title, Owner: "alex", Now: time.Now()})
	if err != nil {
		lab.t.Fatal(err)
	}
	st, err := lab.repo.Create(workitem.NewOptions{Type: workitem.Story, Title: title, Parent: epic.ID, Owner: "alex", Now: time.Now()})
	if err != nil {
		lab.t.Fatal(err)
	}
	data, _ := os.ReadFile(st.Path)
	if err := os.WriteFile(st.Path, []byte(strings.Replace(string(data), "## Acceptance criteria\n", "## Acceptance criteria\n- [ ] deploys\n", 1)), 0o644); err != nil {
		lab.t.Fatal(err)
	}
	lab.move(st.ID, workitem.Ready)
	lab.move(st.ID, workitem.InProgress)
	if _, err := lab.repo.Create(workitem.NewOptions{Type: workitem.Task, Title: "Work", Parent: st.ID, Owner: "builder", Now: time.Now()}); err != nil {
		lab.t.Fatal(err)
	}
	lab.move(st.ID, workitem.Review)
	if err := os.MkdirAll(lab.repo.WorktreePath(st.ID), 0o755); err != nil {
		lab.t.Fatal(err)
	}
	lab.git.set(st.ID, head)
	return st.ID
}

func (lab *klusterLab) move(id, to string) {
	lab.t.Helper()
	it, err := lab.repo.Get(id)
	if err != nil {
		lab.t.Fatal(err)
	}
	if _, err := lab.repo.Transition(it, to, "alex", "the test moves it", time.Now()); err != nil {
		lab.t.Fatal(err)
	}
}

// round looks once and runs what the look queued, in this goroutine.
func (lab *klusterLab) round() {
	lab.t.Helper()
	lab.mu.Lock()
	lab.clock = lab.clock.Add(klusterEvery)
	lab.mu.Unlock()
	ctx := context.Background()
	lab.pass.look(ctx)
	lab.pass.drain(ctx)
}

// threadsOn are the deploy failure threads flai serve opened on story, open
// or not.
func (lab *klusterLab) threadsOn(story string) []*threads.Thread {
	return lab.titled(story, klusterFailureTitle(story))
}

// titled are the threads titled title that flai serve opened on on, open or
// not.
func (lab *klusterLab) titled(on, title string) []*threads.Thread {
	lab.t.Helper()
	all, err := threads.For(lab.repo, on)
	if err != nil {
		lab.t.Fatal(err)
	}
	var out []*threads.Thread
	for _, th := range all {
		if th.Opener() == forgeAuthor && th.Title == workitem.CleanTitle(title) {
			out = append(out, th)
		}
	}
	return out
}

func (lab *klusterLab) entries() []hostapi.Entry {
	lab.mu.Lock()
	defer lab.mu.Unlock()
	return append([]hostapi.Entry(nil), lab.journal...)
}

func (lab *klusterLab) status(id string) string {
	lab.t.Helper()
	it, err := lab.repo.Get(id)
	if err != nil {
		lab.t.Fatal(err)
	}
	return it.Status
}

const (
	klHeadA = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	klHeadB = "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"
	klHeadC = "cccccccccccccccccccccccccccccccccccccccc"
)

func TestAStoryEnteringReviewIsDeployedOnceAndAgainAtANewHead(t *testing.T) {
	lab := newKlusterLab(t, klusterBlock)
	story := lab.inReview("Deployed", klHeadA)

	lab.round()
	want := []string{"deploy " + story + " " + klHeadA}
	if got := lab.kluster.deploys(); strings.Join(got, ";") != strings.Join(want, ";") {
		t.Fatalf("deploys after entering review = %v, want %v", got, want)
	}
	journal := lab.entries()
	if len(journal) != 2 || journal[0].Outcome != "started" || journal[1].Outcome != "done" {
		t.Fatalf("journal = %+v, want a started entry and a done entry", journal)
	}
	for _, e := range journal {
		if e.Action != hostapi.ActionKluster || e.Method != klusterMethod || !strings.Contains(e.Detail, story) || !strings.Contains(e.Detail, "s-") {
			t.Errorf("journal entry %+v does not name the kluster action, its method, the story, and its namespace", e)
		}
	}

	lab.round()
	if got := lab.kluster.deploys(); len(got) != 1 {
		t.Fatalf("a second look at the same head deployed again: %v", got)
	}

	lab.git.set(story, klHeadB)
	lab.round()
	want = append(want, "deploy "+story+" "+klHeadB)
	if got := lab.kluster.deploys(); strings.Join(got, ";") != strings.Join(want, ";") {
		t.Fatalf("deploys after a new head = %v, want %v", got, want)
	}
}

func TestAStoryDeployedFromAShellAtItsHeadIsNotDeployedAgain(t *testing.T) {
	lab := newKlusterLab(t, klusterBlock)
	story := lab.inReview("Shell", klHeadA)
	ns, _ := kluster.StoryNamespace(story)
	record := filepath.Join(kluster.Dir(lab.root), "deployments", ns+".json")
	if err := os.MkdirAll(filepath.Dir(record), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(record, []byte(`{"story":"`+story+`","namespace":"`+ns+`","commit":"`+klHeadA+`","outcome":"passed"}`), 0o644); err != nil {
		t.Fatal(err)
	}
	lab.round()
	if got := lab.kluster.all(); len(got) != 0 {
		t.Fatalf("a story whose record is at its head was deployed again: %v", got)
	}
}

func TestAProjectsDeploysRunOneAtATimeInReviewOrder(t *testing.T) {
	lab := newKlusterLab(t, klusterBlock)
	first := lab.inReview("First", klHeadA)
	second := lab.inReview("Second", klHeadB)
	lab.kluster.gate = make(chan struct{})
	lab.kluster.started = make(chan string, 2)

	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	lab.pass.wg = &wg
	lab.pass.start(ctx)
	defer func() {
		cancel()
		wg.Wait()
	}()

	if got := <-lab.kluster.started; got != first {
		t.Fatalf("first deploy was %s, want %s", got, first)
	}
	// a look while the first runs queues nothing twice
	lab.pass.look(ctx)
	lab.kluster.gate <- struct{}{}
	if got := <-lab.kluster.started; got != second {
		t.Fatalf("second deploy was %s, want %s", got, second)
	}
	lab.kluster.gate <- struct{}{}
	lab.kluster.mu.Lock()
	most := lab.kluster.most
	lab.kluster.mu.Unlock()
	if most != 1 {
		t.Fatalf("%d deploys ran at once, want one at a time", most)
	}
	cancel()
	wg.Wait()
	if got := lab.kluster.deploys(); len(got) != 2 {
		t.Fatalf("deploys = %v, want one for each story", got)
	}
}

// ADR-0254 §5a: a failed deploy opens one thread, sends the story back to
// in-progress as flai with a reason naming the job, the step, and the
// commit, and keeps the send back for the agent's start; a look while it is
// in progress tries nothing again. A failure at its next head in review adds
// to the one thread and sends it back again, and a pass resolves it.
func TestAFailedDeployOpensOneThreadSendsTheStoryBackAndAPassResolvesIt(t *testing.T) {
	lab := newKlusterLab(t, klusterBlock)
	story := lab.inReview("Broken", klHeadA)
	ns, _ := kluster.StoryNamespace(story)
	log := filepath.Join(kluster.LogDir(lab.root, ns), "deploy.log")
	lab.kluster.fail = func(_, head string) error {
		if head == klHeadC {
			return nil
		}
		return &kluster.StepError{Namespace: ns, Step: kluster.StepDeploy, ExitCode: 3, Log: log, Err: errors.New("exit status 3")}
	}

	lab.round()
	ths := lab.threadsOn(story)
	if len(ths) != 1 || !ths[0].Open() {
		t.Fatalf("threads after a failed deploy = %d, want 1 open", len(ths))
	}
	for _, want := range []string{"deploy step", "exit code 3", log, klHeadA, ns, "sent " + story + " back to in-progress", "keeps its namespace"} {
		if !strings.Contains(ths[0].Body, want) {
			t.Errorf("the failure thread does not name %q:\n%s", want, ths[0].Body)
		}
	}
	lab.sentBackBy(story, ths[0].ID, "flai serve's kluster deploy of "+story+" failed at the deploy step, at commit "+short(klHeadA))
	if journal := lab.entries(); journal[len(journal)-1].Outcome != "done" || !strings.Contains(journal[len(journal)-1].Detail, "sent "+story+" back to in-progress on "+ths[0].ID) {
		t.Errorf("the journal's last entry is %+v, want the send back", journal[len(journal)-1])
	}

	n := len(lab.kluster.all())
	lab.round()
	if got := lab.kluster.since(n); len(got) != 0 {
		t.Fatalf("a look at the sent back story ran %v, want nothing: no deploy, and no teardown", got)
	}
	if ths := lab.threadsOn(story); len(ths) != 1 || len(ths[0].Entries()) != 1 {
		t.Fatalf("a second look at the failed head changed the threads: %d", len(ths))
	}

	lab.git.set(story, klHeadB)
	lab.move(story, workitem.Review)
	lab.round()
	ths = lab.threadsOn(story)
	if len(ths) != 1 || len(ths[0].Entries()) != 2 || !ths[0].Open() {
		t.Fatalf("a failure at a new head should add an entry to the one open thread; threads %d", len(ths))
	}
	lab.sentBackBy(story, ths[0].ID, "flai serve's kluster deploy of "+story+" failed at the deploy step, at commit "+short(klHeadB))

	lab.git.set(story, klHeadC)
	lab.move(story, workitem.Review)
	lab.round()
	ths = lab.threadsOn(story)
	if len(ths) != 1 || ths[0].Open() {
		t.Fatalf("a passing deploy should resolve the failure thread; threads %d", len(ths))
	}
	if got := lab.status(story); got != workitem.Review {
		t.Fatalf("the story is %s after its deploy passed, want it left in review", got)
	}
	if tr, _ := readKlusterTried(lab.root, ns); tr.SentBack != nil {
		t.Errorf("the passing deploy kept the old send back: %+v", tr.SentBack)
	}
}

// sentBackBy fails the test unless story stands in progress, moved there by
// flai as its newest move, and what the pass tried for its namespace keeps
// the send back on thread for reason.
func (lab *klusterLab) sentBackBy(story, thread, reason string) {
	lab.t.Helper()
	it, err := lab.repo.Get(story)
	if err != nil {
		lab.t.Fatal(err)
	}
	if last := it.Transitions[len(it.Transitions)-1]; it.Status != workitem.InProgress || last.By != sendBackBy {
		lab.t.Fatalf("%s is %s, last moved by %s; want it sent back to in-progress by %s", story, it.Status, last.By, sendBackBy)
	}
	ns, _ := kluster.StoryNamespace(story)
	tr, ok := readKlusterTried(lab.root, ns)
	if !ok || tr.SentBack == nil || tr.SentBack.Thread != thread || tr.SentBack.Reason != reason || !klusterKept(it, tr) {
		lab.t.Fatalf("the send back kept for %s is %+v, want one on %s for %q that keeps the story", ns, tr.SentBack, thread, reason)
	}
}

// ADR-0254 §5a: a story sent back after a failed deploy keeps its
// namespace while it stands in progress, however many looks pass, and is torn
// down once it is cancelled, its failure thread resolved.
func TestASentBackStoryKeepsItsNamespaceUntilItIsCancelled(t *testing.T) {
	lab := newKlusterLab(t, klusterBlock)
	story := lab.inReview("Kept", klHeadA)
	ns, _ := kluster.StoryNamespace(story)
	lab.kluster.fail = func(string, string) error {
		return &kluster.StepError{Namespace: ns, Step: kluster.StepDeploy, ExitCode: 1, Err: errors.New("exit status 1")}
	}
	lab.round()
	if got := lab.status(story); got != workitem.InProgress {
		t.Fatalf("the story is %s after a failed deploy, want it sent back to in-progress", got)
	}
	n := len(lab.kluster.all())
	lab.round()
	lab.round()
	if got := lab.kluster.since(n); len(got) != 0 {
		t.Fatalf("looks at a sent back story ran %v, want its namespace kept", got)
	}
	if _, ok, _ := kluster.ReadRecord(lab.root, ns); !ok {
		t.Fatalf("the sent back story's namespace %s has no record", ns)
	}

	lab.move(story, workitem.Cancelled)
	lab.round()
	if got := lab.kluster.since(n); strings.Join(got, ";") != "teardown "+story {
		t.Fatalf("calls once the sent back story was cancelled = %v, want one teardown", got)
	}
	if _, ok := readKlusterTried(lab.root, ns); ok {
		t.Errorf("what flai serve tried for %s is kept after its teardown passed", ns)
	}
	if ths := lab.threadsOn(story); len(ths) != 1 || ths[0].Open() {
		t.Errorf("the cancelled story's deploy failure thread should be resolved by its teardown")
	}
}

// ADR-0254 §5a: a sent back story that returns to review is deployed again
// in place, at the same head or a new one, with no teardown between, and a
// pass resolves the thread. Once it has passed, a move back to in-progress
// by its agent keeps its namespace too (ADR-0292 §3), and a move to ready
// tears it down.
func TestASentBackStorysReturnToReviewDeploysAgainInPlace(t *testing.T) {
	for _, c := range []struct{ name, head string }{{"same head", klHeadA}, {"new head", klHeadB}} {
		t.Run(c.name, func(t *testing.T) {
			lab := newKlusterLab(t, klusterBlock)
			story := lab.inReview("Returns", klHeadA)
			ns, _ := kluster.StoryNamespace(story)
			failing := true
			lab.kluster.fail = func(string, string) error {
				if failing {
					return &kluster.StepError{Namespace: ns, Step: kluster.StepBuild, ExitCode: 2, Err: errors.New("exit status 2")}
				}
				return nil
			}
			lab.round()
			if got := lab.status(story); got != workitem.InProgress {
				t.Fatalf("the story is %s after a failed deploy, want it sent back", got)
			}
			failing = false
			n := len(lab.kluster.all())
			lab.git.set(story, c.head)
			lab.move(story, workitem.Review)
			lab.round()
			if got := lab.kluster.since(n); strings.Join(got, ";") != "absent;deploy "+story+" "+c.head {
				t.Fatalf("calls once the story was back in review = %v, want one deploy at %s and no teardown", got, short(c.head))
			}
			if ths := lab.threadsOn(story); len(ths) != 1 || ths[0].Open() {
				t.Fatalf("the deploy that passed should resolve the failure thread")
			}
			if got := lab.status(story); got != workitem.Review {
				t.Fatalf("the story is %s after its deploy passed, want it in review", got)
			}
			n = len(lab.kluster.all())
			lab.round()
			if got := lab.kluster.since(n); len(got) != 0 {
				t.Fatalf("a look after the deploy in place ran %v", got)
			}

			lab.move(story, workitem.InProgress)
			lab.round()
			if got := lab.kluster.since(n); len(got) != 0 {
				t.Fatalf("calls once its agent moved it back to in-progress = %v, want its namespace kept", got)
			}
			lab.move(story, workitem.Ready)
			lab.round()
			if got := lab.kluster.since(n); strings.Join(got, ";") != "teardown "+story {
				t.Fatalf("calls once it moved to ready = %v, want one teardown", got)
			}
		})
	}
}

func TestNoKlusterCommandRunsWithTheActionOffOrNoBlock(t *testing.T) {
	off := newKlusterLab(t, klusterBlock)
	off.inReview("Off", klHeadA)
	off.setOn(false)
	off.round()
	if got := off.kluster.all(); len(got) != 0 {
		t.Fatalf("the action off ran %v", got)
	}

	none := newKlusterLab(t, "")
	none.inReview("None", klHeadA)
	none.round()
	if got := none.kluster.all(); len(got) != 0 {
		t.Fatalf("a project with no kluster: block ran %v", got)
	}

	late := newKlusterLab(t, klusterBlock)
	late.inReview("Late", klHeadA)
	late.pass.look(context.Background())
	late.setOn(false)
	late.pass.drain(context.Background())
	if got := late.kluster.all(); len(got) != 0 {
		t.Fatalf("a job queued before the action was turned off ran %v", got)
	}
}

// klusterNoDeploy is the block enabling the action writes (S-0586): a config
// and no deploy.
const klusterNoDeploy = "kluster:\n  config: kluster.yaml\n"

// setBlock rewrites the manifest's kluster: block, as a save from the
// dashboard does between two looks.
func (lab *klusterLab) setBlock(block string) {
	lab.t.Helper()
	m := "version: 1\nname: t\nkey: t\nlayout:\n  design: design\n  docs: docs\n  wip: wip\n" + block
	if err := os.WriteFile(filepath.Join(lab.root, "system-flow.yaml"), []byte(m), 0o644); err != nil {
		lab.t.Fatal(err)
	}
}

func TestABlockWithNoDeployDeploysNothingUntilOneIsSet(t *testing.T) {
	lab := newKlusterLab(t, klusterNoDeploy)
	lab.kluster.absent = true
	accepted := lab.inReview("Accepted", klHeadA)
	review := lab.inReview("Review", klHeadB)
	lab.accept(accepted)
	lab.round()
	if got := lab.kluster.all(); len(got) != 0 {
		t.Fatalf("a block with no deploy ran %v, want no job, no bring-up, no promotion", got)
	}
	for _, story := range []string{accepted, review} {
		if ths := lab.threadsOn(story); len(ths) != 0 {
			t.Fatalf("a block with no deploy opened %d threads on %s", len(ths), story)
		}
	}
	if ths := lab.titled(manifest.File, promoteFailureTitle()); len(ths) != 0 {
		t.Fatalf("a block with no deploy opened %d promotion threads", len(ths))
	}
	ns, _ := kluster.StoryNamespace(review)
	if _, ok := readKlusterTried(lab.root, ns); ok {
		t.Fatalf("a block with no deploy recorded a try for %s, which would keep it from deploying once deploy is set", ns)
	}

	lab.setBlock("kluster:\n  config: kluster.yaml\n  deploy: [\"\"]\n")
	lab.round()
	if got := lab.kluster.all(); len(got) != 0 {
		t.Fatalf("a deploy whose program is empty ran %v", got)
	}

	lab.setBlock("kluster:\n  config: kluster.yaml\n  deploy: [scripts/kluster-deploy.sh]\n")
	lab.round()
	want := "absent;up;deploy " + review + " " + klHeadB
	if got := strings.Join(lab.kluster.all(), ";"); got != want {
		t.Fatalf("calls once deploy is set = %s, want %s: the story in review deployed, and the acceptance before it not promoted", got, want)
	}
}

func TestAJobQueuedBeforeDeployWasRemovedDoesNotRun(t *testing.T) {
	lab := newKlusterLab(t, klusterBlock)
	lab.inReview("Removed", klHeadA)
	lab.pass.look(context.Background())
	lab.setBlock(klusterNoDeploy)
	lab.pass.drain(context.Background())
	if got := lab.kluster.all(); len(got) != 0 {
		t.Fatalf("a deploy queued before deploy was removed ran %v", got)
	}
}

func TestADeployedStoryThatLeftReviewIsTornDownWithNoDeploySet(t *testing.T) {
	lab := newKlusterLab(t, klusterBlock)
	story := lab.inReview("Leaves", klHeadA)
	lab.round()
	lab.setBlock(klusterNoDeploy)
	n := len(lab.kluster.all())
	lab.move(story, workitem.InProgress)
	lab.move(story, workitem.Ready)
	lab.round()
	if got := lab.kluster.since(n); strings.Join(got, ";") != "teardown "+story {
		t.Fatalf("calls after %s left review for ready with no deploy set = %v, want one teardown", story, got)
	}
}

func TestABusyClusterOpensNoThreadAndTheNextLookTriesAgain(t *testing.T) {
	lab := newKlusterLab(t, klusterBlock)
	story := lab.inReview("Busy", klHeadA)
	busy := true
	lab.kluster.fail = func(string, string) error {
		if busy {
			busy = false
			return fmt.Errorf("%w: a deploy of the namespace s-0001, started by process 1", kluster.ErrBusy)
		}
		return nil
	}
	lab.round()
	if ths := lab.threadsOn(story); len(ths) != 0 {
		t.Fatalf("a busy cluster opened %d threads", len(ths))
	}
	lab.round()
	if got := lab.kluster.deploys(); len(got) != 2 {
		t.Fatalf("deploys = %v, want the busy one tried again at the next look", got)
	}
}

func TestAnAbsentClusterIsBroughtUpBeforeTheDeploy(t *testing.T) {
	lab := newKlusterLab(t, klusterBlock)
	story := lab.inReview("Up", klHeadA)
	lab.kluster.absent = true
	lab.round()
	want := "absent;up;deploy " + story + " " + klHeadA
	if got := strings.Join(lab.kluster.all(), ";"); got != want {
		t.Fatalf("calls = %s, want %s", got, want)
	}
	if data, err := os.ReadFile(klusterUpLog(lab.root)); err != nil || !strings.Contains(string(data), "kluster up") {
		t.Fatalf("the bring-up's output is not in %s: %q, %v", klusterUpLog(lab.root), data, err)
	}

	failing := newKlusterLab(t, klusterBlock)
	other := failing.inReview("Down", klHeadA)
	failing.kluster.absent = true
	failing.kluster.upErr = errors.New("kluster up of the cluster c failed: exit status 1")
	failing.round()
	if got := failing.kluster.deploys(); len(got) != 0 {
		t.Fatalf("a failed bring-up still deployed: %v", got)
	}
	ths := failing.threadsOn(other)
	if len(ths) != 1 || !strings.Contains(ths[0].Body, "up step") || !strings.Contains(ths[0].Body, klusterUpLog(failing.root)) {
		t.Fatalf("a failed bring-up should open one thread naming the up step and its log; got %d", len(ths))
	}
}

// ADR-0292 §3: a story that leaves review keeps its namespace while it is in
// progress, whoever moved it there, and is torn down once, when it moves on
// to ready or is cancelled.
func TestAStoryThatLeavesReviewIsTornDownOnce(t *testing.T) {
	for _, to := range []string{workitem.Ready, workitem.Cancelled} {
		t.Run(to, func(t *testing.T) {
			lab := newKlusterLab(t, klusterBlock)
			story := lab.inReview("Leaves", klHeadA)
			lab.round()
			n := len(lab.kluster.all())
			lab.move(story, workitem.InProgress)
			lab.round()
			if got := lab.kluster.since(n); len(got) != 0 {
				t.Fatalf("calls after %s left review for in-progress = %v, want its namespace kept", story, got)
			}
			lab.move(story, to)
			lab.round()
			if got := lab.kluster.since(n); strings.Join(got, ";") != "teardown "+story {
				t.Fatalf("calls after %s left review for %s = %v, want one teardown", story, to, got)
			}
			ns, _ := kluster.StoryNamespace(story)
			if _, ok := readKlusterTried(lab.root, ns); ok {
				t.Errorf("what flai serve tried for %s is kept after its teardown passed", ns)
			}
			journal := lab.entries()
			if last := journal[len(journal)-1]; last.Outcome != "done" || !strings.Contains(last.Detail, "tore down") {
				t.Errorf("the journal's last entry is %+v, want the teardown done", last)
			}
			n = len(lab.kluster.all())
			lab.round()
			if got := lab.kluster.since(n); len(got) != 0 {
				t.Fatalf("a second look after the teardown ran %v", got)
			}
		})
	}
}

func TestAStoryNeverDeployedIsNotTornDown(t *testing.T) {
	lab := newKlusterLab(t, klusterBlock)
	never := lab.inReview("Never", klHeadA)
	lab.setOn(false)
	lab.round()
	lab.move(never, workitem.InProgress)
	lab.setOn(true)

	refused := lab.inReview("Refused", klHeadB)
	lab.kluster.fail = func(story, _ string) error {
		if story == refused {
			return errors.New("the cluster has no kubeconfig")
		}
		return nil
	}
	lab.round()
	ns, _ := kluster.StoryNamespace(refused)
	if _, ok := readKlusterTried(lab.root, ns); !ok {
		t.Fatalf("a refused deploy left no record of what flai serve tried")
	}
	if got := lab.status(refused); got != workitem.InProgress {
		t.Fatalf("the refused story is %s, want it sent back to in-progress", got)
	}
	n := len(lab.kluster.all())
	lab.round()
	if got := lab.kluster.since(n); len(got) != 0 {
		t.Fatalf("stories with no deployment were torn down: %v", got)
	}
	if tr, ok := readKlusterTried(lab.root, ns); !ok || tr.SentBack == nil {
		t.Errorf("what flai serve tried for %s is not kept while its send back holds it: %+v", ns, tr)
	}
	lab.move(refused, workitem.Cancelled)
	lab.round()
	if got := lab.kluster.since(n); len(got) != 0 {
		t.Fatalf("a cancelled story with no deployment was torn down: %v", got)
	}
	if _, ok := readKlusterTried(lab.root, ns); ok {
		t.Errorf("what flai serve tried for %s is kept though it was cancelled with nothing deployed", ns)
	}
}

func TestAnAcceptedStoryIsTornDownThenMainPromoted(t *testing.T) {
	lab := newKlusterLab(t, klusterBlock)
	story := lab.inReview("Accepted", klHeadA)
	lab.round()
	n := len(lab.kluster.all())
	lab.accept(story)
	lab.round()
	want := "teardown " + story + ";promote main"
	if got := strings.Join(lab.kluster.since(n), ";"); got != want {
		t.Fatalf("calls after the acceptance = %s, want %s", got, want)
	}
	journal := lab.entries()
	if last := journal[len(journal)-1]; last.Outcome != "done" || !strings.Contains(last.Detail, "promoted main") {
		t.Errorf("the journal's last entry is %+v, want the promotion done", last)
	}
	n = len(lab.kluster.all())
	lab.round()
	if got := lab.kluster.since(n); len(got) != 0 {
		t.Fatalf("a second look after the promotion ran %v", got)
	}
}

func TestTwoAcceptancesAtOneLookPromoteOnce(t *testing.T) {
	lab := newKlusterLab(t, klusterBlock)
	first := lab.inReview("First", klHeadA)
	second := lab.inReview("Second", klHeadB)
	lab.round()
	n := len(lab.kluster.all())
	lab.accept(first)
	lab.accept(second)
	lab.round()
	want := "teardown " + first + ";teardown " + second + ";promote main"
	if got := strings.Join(lab.kluster.since(n), ";"); got != want {
		t.Fatalf("calls after two acceptances = %s, want %s", got, want)
	}
}

func TestAnAcceptanceBeforeThePassStartedPromotesNothing(t *testing.T) {
	lab := newKlusterLab(t, klusterBlock)
	story := lab.inReview("Earlier", klHeadA)
	lab.round()
	lab.accept(story)
	lab.restart()
	n := len(lab.kluster.all())
	lab.round()
	if got := strings.Join(lab.kluster.since(n), ";"); got != "teardown "+story {
		t.Fatalf("calls after a restart = %s, want only the teardown of the story that left review", got)
	}
}

func TestAFailedTeardownOpensOneThreadAndIsNotTriedAgainUntilTheStoryMoves(t *testing.T) {
	lab := newKlusterLab(t, klusterBlock)
	story := lab.inReview("Stuck", klHeadA)
	lab.round()
	ns, _ := kluster.StoryNamespace(story)
	log := kluster.LogPath(lab.root, ns, kluster.StepTeardown)
	lab.kluster.teardownErr = &kluster.StepError{Namespace: ns, Step: kluster.StepTeardown, ExitCode: 2, Log: log, Err: errors.New("exit status 2")}
	lab.move(story, workitem.InProgress)
	lab.move(story, workitem.Ready)
	lab.round()
	ths := lab.titled(story, teardownFailureTitle(story))
	if len(ths) != 1 {
		t.Fatalf("threads after a failed teardown = %d, want 1", len(ths))
	}
	for _, want := range []string{"teardown step", "exit code 2", log, klHeadA, ns} {
		if !strings.Contains(ths[0].Body, want) {
			t.Errorf("the teardown failure thread does not name %q:\n%s", want, ths[0].Body)
		}
	}
	if got := lab.status(story); got != workitem.Ready {
		t.Fatalf("the story is %s after a failed teardown, want it left ready", got)
	}
	n := len(lab.kluster.all())
	lab.round()
	if got := lab.kluster.since(n); len(got) != 0 {
		t.Fatalf("a failed teardown was tried again at the next look: %v", got)
	}

	lab.kluster.teardownErr = nil
	lab.move(story, workitem.Cancelled)
	lab.round()
	if got := lab.kluster.since(n); strings.Join(got, ";") != "teardown "+story {
		t.Fatalf("calls once the story moved again = %v, want the teardown tried again", got)
	}
	if ths := lab.titled(story, teardownFailureTitle(story)); len(ths) != 1 || ths[0].Open() {
		t.Fatalf("a passing teardown should resolve the one failure thread")
	}
}

func TestAFailedPromotionOpensOneThreadOnTheManifestAndAPassResolvesIt(t *testing.T) {
	lab := newKlusterLab(t, klusterBlock)
	first := lab.inReview("First", klHeadA)
	second := lab.inReview("Second", klHeadB)
	lab.round()
	log := kluster.LogPath(lab.root, "t", kluster.StepDeploy)
	lab.kluster.promoteErr = &kluster.StepError{Namespace: "t", Step: kluster.StepDeploy, ExitCode: 4, Log: log, Err: errors.New("exit status 4")}
	lab.accept(first)
	lab.round()
	ths := lab.titled(manifest.File, promoteFailureTitle())
	if len(ths) != 1 {
		t.Fatalf("threads on %s after a failed promotion = %d, want 1", manifest.File, len(ths))
	}
	for _, want := range []string{"deploy step", "exit code 4", log, klHeadC, "`t`"} {
		if !strings.Contains(ths[0].Body, want) {
			t.Errorf("the promotion failure thread does not name %q:\n%s", want, ths[0].Body)
		}
	}
	n := len(lab.kluster.all())
	lab.round()
	if got := lab.kluster.since(n); len(got) != 0 {
		t.Fatalf("a failed promotion was tried again with no new acceptance: %v", got)
	}
	if ths := lab.titled(manifest.File, promoteFailureTitle()); len(ths) != 1 || len(ths[0].Entries()) != 1 {
		t.Fatalf("a second look changed the promotion's threads")
	}

	lab.kluster.promoteErr = nil
	lab.accept(second)
	lab.round()
	if got := strings.Join(lab.kluster.since(n), ";"); got != "teardown "+second+";promote main" {
		t.Fatalf("calls after the next acceptance = %s, want its teardown and a promotion", got)
	}
	if ths := lab.titled(manifest.File, promoteFailureTitle()); len(ths) != 1 || ths[0].Open() {
		t.Fatalf("a passing promotion should resolve the one failure thread")
	}
}

func TestNoTeardownOrPromotionRunsWithTheActionOff(t *testing.T) {
	lab := newKlusterLab(t, klusterBlock)
	left := lab.inReview("Left", klHeadA)
	accepted := lab.inReview("Accepted", klHeadB)
	lab.round()
	n := len(lab.kluster.all())
	lab.setOn(false)
	lab.move(left, workitem.InProgress)
	lab.accept(accepted)
	lab.round()
	if got := lab.kluster.since(n); len(got) != 0 {
		t.Fatalf("the action off ran %v", got)
	}
	lab.setOn(true)
	lab.round()
	for _, c := range lab.kluster.since(n) {
		if strings.HasPrefix(c, "promote") {
			t.Fatalf("an acceptance while the action was off promoted main: %v", lab.kluster.since(n))
		}
	}
}

// withConfig writes the kluster config file the lab's block names, so that
// the cluster can be named, as KlusterEnv needs.
func (lab *klusterLab) withConfig() {
	lab.t.Helper()
	path := filepath.Join(lab.root, "integration", "kluster.yaml")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		lab.t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte("name: t-int\nprovider: k3d\n"), 0o644); err != nil {
		lab.t.Fatal(err)
	}
}

// envOf is KlusterEnv for story in the lab, failing the test on an error.
func (lab *klusterLab) envOf(story string) (string, []string) {
	lab.t.Helper()
	ns, env, err := KlusterEnv(lab.repo, story, lab.git, lab.opts.Host.Enabled)
	if err != nil {
		lab.t.Fatal(err)
	}
	return ns, env
}

func TestAStoryDeployedAtItsHeadGetsItsNamespacesVariablesForItsChecks(t *testing.T) {
	lab := newKlusterLab(t, klusterBlock)
	lab.withConfig()
	story := lab.inReview("Checked", klHeadA)
	if _, env := lab.envOf(story); env != nil {
		t.Fatalf("a story not deployed yet got %v", env)
	}

	lab.round()
	ns, env := lab.envOf(story)
	want := map[string]string{
		"FLAI_STORY":             story,
		"FLAI_KLUSTER_NAMESPACE": ns,
		"FLAI_KLUSTER_NAME":      "t-int",
		"FLAI_KLUSTER_COMMIT":    klHeadA,
		"FLAI_KLUSTER_CONTEXT":   "k3d-t-int",
		"KUBECONFIG":             kluster.KubeconfigPath(lab.root),
	}
	if ns != strings.ToLower(story) {
		t.Errorf("namespace = %q, want %q", ns, strings.ToLower(story))
	}
	for name, value := range want {
		if got := envValue(env, name); got != value {
			t.Errorf("%s = %q, want %q (env %v)", name, got, value, env)
		}
	}

	lab.git.set(story, klHeadB) // a new head, not deployed yet
	if _, env := lab.envOf(story); env != nil {
		t.Errorf("a story at a head its namespace was not deployed at got %v", env)
	}
	lab.git.set(story, klHeadA)
	lab.setOn(false)
	if _, env := lab.envOf(story); env != nil {
		t.Errorf("with the kluster action off a story got %v", env)
	}
}

func TestAStoryWhoseDeployFailedGetsNoNamespaceForItsChecks(t *testing.T) {
	lab := newKlusterLab(t, klusterBlock)
	lab.withConfig()
	lab.kluster.fail = func(string, string) error {
		return &kluster.StepError{Step: "deploy", ExitCode: 2, Err: errors.New("exit status 2")}
	}
	story := lab.inReview("Broken", klHeadA)
	lab.round()
	if _, env := lab.envOf(story); env != nil {
		t.Errorf("a story whose deploy failed got %v", env)
	}
}

func TestAChecksRunAgainstANamespaceKeepsItsOutcomeAndCommitInTheRecord(t *testing.T) {
	lab := newKlusterLab(t, klusterBlock)
	failing, _ := nsCheck(t, 1)
	lab.withChecks(failing)
	story := lab.inReview("Checked", klHeadA)
	ns := strings.ToLower(story)

	lab.round()
	r, ok, err := kluster.ReadRecord(lab.root, ns)
	if err != nil || !ok {
		t.Fatal(ok, err)
	}
	if r.Commit != klHeadA || r.Outcome != kluster.OutcomePassed || r.ChecksCommit != klHeadA || r.ChecksOutcome != kluster.OutcomeFailed {
		t.Fatalf("a failed checks run at the deployed head is kept in the record: %+v", r)
	}
	if got := lab.titled(story, checksFailureTitle(story)); len(got) != 1 {
		t.Errorf("the failed run opens one thread, got %d", len(got))
	}
	if ok, reason, err := kluster.AcceptGate(lab.root, story, klHeadA, true); ok || err != nil || !strings.Contains(reason, "ended failed") {
		t.Errorf("the gate after the failed run: %v %q %v", ok, reason, err)
	}

	passing, _ := nsCheck(t, 0)
	lab.withChecks(passing)
	lab.git.set(story, klHeadB)
	lab.move(story, workitem.Review) // the failed run sent it back (ADR-0254 §5a)
	lab.round()
	if r, _, _ = kluster.ReadRecord(lab.root, ns); r.Commit != klHeadB || r.ChecksCommit != klHeadB || r.ChecksOutcome != kluster.OutcomePassed {
		t.Fatalf("the passed run at the new head is kept in the record: %+v", r)
	}
	if ok, reason, err := kluster.AcceptGate(lab.root, story, klHeadB, true); !ok || err != nil {
		t.Errorf("the gate after the passed run: %v %q %v", ok, reason, err)
	}
}

// ADR-0254 §5a, §8: a checks run against the namespace that fails sends the
// story back as a failed deploy does, and keeps its namespace; back in review
// at the same head, it is deployed and checked again in place, and a passing
// run resolves the checks' thread.
func TestFailingChecksAgainstTheNamespaceSendTheStoryBack(t *testing.T) {
	lab := newKlusterLab(t, klusterBlock)
	failing, _ := nsCheck(t, 3)
	lab.withChecks(failing)
	story := lab.inReview("Failing", klHeadA)

	lab.round()
	ths := lab.titled(story, checksFailureTitle(story))
	if len(ths) != 1 || !ths[0].Open() {
		t.Fatalf("checks failure threads = %d, want one open", len(ths))
	}
	for _, want := range []string{"ended failed", "golden", klHeadA, "sent " + story + " back to in-progress"} {
		if !strings.Contains(ths[0].Body, want) {
			t.Errorf("the thread does not name %q:\n%s", want, ths[0].Body)
		}
	}
	lab.sentBackBy(story, ths[0].ID, "flai serve's checks of "+story+" against its kluster namespace ended failed, at the check golden, at commit "+short(klHeadA))

	n := len(lab.kluster.all())
	lab.round()
	if got := lab.kluster.since(n); len(got) != 0 {
		t.Fatalf("a look at the story the checks sent back ran %v, want its namespace kept", got)
	}

	passing, _ := nsCheck(t, 0)
	lab.withChecks(passing)
	lab.move(story, workitem.Review)
	lab.round()
	if got := lab.kluster.since(n); strings.Join(got, ";") != "absent;deploy "+story+" "+klHeadA {
		t.Fatalf("calls once the story was back in review at the same head = %v, want one deploy in place", got)
	}
	if ths := lab.titled(story, checksFailureTitle(story)); len(ths) != 1 || ths[0].Open() {
		t.Fatalf("the passing checks run should resolve the checks failure thread")
	}
	if got := lab.status(story); got != workitem.Review {
		t.Errorf("the story is %s after its checks passed, want it in review", got)
	}
}
