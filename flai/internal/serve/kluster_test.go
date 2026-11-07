package serve

import (
	"context"
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

// fakeKluster is a KlusterRunner that runs nothing: it records each call,
// answers Absent with absent, fails Up with upErr, and fails each deploy with
// what fail answers for its story and head. While gate is not nil, each
// deploy says it started on started and waits for gate before it ends.
type fakeKluster struct {
	mu       sync.Mutex
	calls    []string // "absent", "up", or "deploy S-nnnn <head>"
	absent   bool
	upErr    error
	fail     func(story, head string) error
	gate     chan struct{}
	started  chan string
	inFlight int
	most     int
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

func (f *fakeKluster) Deploy(_ context.Context, _ string, _ manifest.Manifest, story, _, commit string) (kluster.Record, error) {
	f.note("deploy " + story + " " + commit)
	f.mu.Lock()
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
	if fail != nil {
		if err := fail(story, commit); err != nil {
			return kluster.Record{}, err
		}
	}
	return kluster.Record{Story: story, Commit: commit, Outcome: kluster.OutcomePassed}, nil
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

// klusterHeads answers git rev-parse of a story's branch from heads, by story.
type klusterHeads struct {
	mu    sync.Mutex
	heads map[string]string
}

func (g *klusterHeads) set(story, head string) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.heads[story] = head
}

func (g *klusterHeads) Run(_, name string, args ...string) (string, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
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
	lab.pass = newKlusterPass(o, Entry{Key: "t", Name: "t", Root: root}, &sync.WaitGroup{})
	return lab
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
	lab.t.Helper()
	all, err := threads.For(lab.repo, story)
	if err != nil {
		lab.t.Fatal(err)
	}
	var out []*threads.Thread
	for _, th := range all {
		if th.Opener() == forgeAuthor && th.Title == workitem.CleanTitle(klusterFailureTitle(story)) {
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

func TestAFailedDeployOpensOneThreadMovesNothingAndAPassResolvesIt(t *testing.T) {
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
	if len(ths) != 1 {
		t.Fatalf("threads after a failed deploy = %d, want 1", len(ths))
	}
	for _, want := range []string{"deploy step", "exit code 3", log, klHeadA, ns} {
		if !strings.Contains(ths[0].Body, want) {
			t.Errorf("the failure thread does not name %q:\n%s", want, ths[0].Body)
		}
	}
	if got := lab.status(story); got != workitem.Review {
		t.Fatalf("the story is %s after a failed deploy, want it left in review", got)
	}
	if journal := lab.entries(); journal[len(journal)-1].Outcome != "failed" {
		t.Errorf("the journal's last entry is %+v, want the failure", journal[len(journal)-1])
	}

	lab.round()
	if got := lab.kluster.deploys(); len(got) != 1 {
		t.Fatalf("a failed head was tried again: %v", got)
	}
	if ths := lab.threadsOn(story); len(ths) != 1 || len(ths[0].Entries()) != 1 {
		t.Fatalf("a second look at the failed head changed the threads: %d", len(ths))
	}

	lab.git.set(story, klHeadB)
	lab.round()
	ths = lab.threadsOn(story)
	if len(ths) != 1 || len(ths[0].Entries()) != 2 || !ths[0].Open() {
		t.Fatalf("a failure at a new head should add an entry to the one open thread; threads %d", len(ths))
	}

	lab.git.set(story, klHeadC)
	lab.round()
	ths = lab.threadsOn(story)
	if len(ths) != 1 || ths[0].Open() {
		t.Fatalf("a passing deploy should resolve the failure thread; threads %d", len(ths))
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
