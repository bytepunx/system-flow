package hostapi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/channel"
	"github.com/bytepunx/system-flow/flai/internal/kluster"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// fakeCluster is a cluster that runs no kluster: Up and Down write a line to
// their output, wait for release when it is set, and return their error.
type fakeCluster struct {
	name      string
	status    kluster.Status
	statusErr error
	upErr     error
	downErr   error
	readyErr  error
	deployErr error
	services  []KlusterService
	release   chan struct{}
	mu        sync.Mutex
	ran       []string
	// calls are the deploys, teardowns, and promotions run, with what each
	// was given.
	calls       []string
	listedTimes int
}

func (f *fakeCluster) Named() kluster.Status {
	return kluster.Status{Name: f.name, Profile: "full", Provider: "k3d", Context: "k3d-" + f.name}
}

func (f *fakeCluster) job(ctx context.Context, out io.Writer, action string, err error) error {
	f.mu.Lock()
	f.ran = append(f.ran, action)
	f.mu.Unlock()
	_, _ = fmt.Fprintf(out, "kluster %s of %s\n", action, f.name)
	if f.release != nil {
		select {
		case <-f.release:
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	return err
}

func (f *fakeCluster) Up(ctx context.Context, out io.Writer) error {
	return f.job(ctx, out, KlusterUp, f.upErr)
}

func (f *fakeCluster) Down(ctx context.Context, out io.Writer) error {
	return f.job(ctx, out, KlusterDown, f.downErr)
}

func (f *fakeCluster) Status(context.Context) (kluster.Status, error) {
	if f.statusErr != nil {
		return kluster.Status{}, f.statusErr
	}
	return f.status, nil
}

func (f *fakeCluster) Ready(manifest.Manifest, bool) error { return f.readyErr }

func (f *fakeCluster) Deploy(ctx context.Context, _ manifest.Manifest, story, status, dir, commit string, out io.Writer) (kluster.Record, error) {
	f.mu.Lock()
	f.calls = append(f.calls, fmt.Sprintf("deploy %s %s %s %s", story, status, dir, commit))
	f.mu.Unlock()
	return kluster.Record{}, f.job(ctx, out, KlusterDeploy, f.deployErr)
}

func (f *fakeCluster) Teardown(ctx context.Context, _ manifest.Manifest, story, dir string, out io.Writer) (kluster.Record, error) {
	f.mu.Lock()
	f.calls = append(f.calls, fmt.Sprintf("teardown %s %s", story, dir))
	f.mu.Unlock()
	return kluster.Record{}, f.job(ctx, out, KlusterTeardown, f.deployErr)
}

func (f *fakeCluster) Promote(ctx context.Context, _ manifest.Manifest, branch string, out io.Writer) (kluster.Record, error) {
	f.mu.Lock()
	f.calls = append(f.calls, "promote "+branch)
	f.mu.Unlock()
	return kluster.Record{}, f.job(ctx, out, KlusterPromote, f.deployErr)
}

func (f *fakeCluster) Services(context.Context) ([]KlusterService, error) {
	f.mu.Lock()
	f.listedTimes++
	f.mu.Unlock()
	return f.services, nil
}

func (f *fakeCluster) called() []string {
	f.mu.Lock()
	defer f.mu.Unlock()
	return append([]string{}, f.calls...)
}

func (f *fakeCluster) actions() []string {
	f.mu.Lock()
	defer f.mu.Unlock()
	return append([]string{}, f.ran...)
}

// quayside is a project whose manifest has a kluster: block when block says
// so, with c as its cluster for the test's length.
func quayside(t *testing.T, block bool, c *fakeCluster) channel.Project {
	t.Helper()
	root := t.TempDir()
	m := "version: 1\nname: Harbour\nkey: harbour\nowner: olive\nlayout:\n  design: design\n  docs: docs\n  wip: wip\n"
	if block {
		m += "kluster:\n  config: integration/kluster.yaml\n  deploy: [scripts/kluster-deploy.sh]\n"
	}
	if err := os.WriteFile(filepath.Join(root, manifest.File), []byte(m), 0o644); err != nil {
		t.Fatal(err)
	}
	was := openKluster
	t.Cleanup(func() { openKluster = was })
	openKluster = func(string, manifest.Manifest) (KlusterCluster, error) { return c, nil }
	return channel.Project{Key: "harbour", Name: "Harbour", Root: root}
}

// klusterHost is a host whose kluster action is on when on says so,
// journalling into journal, with jobs running its ups and downs.
func klusterHost(on *bool, journal *[]Entry, jobs *KlusterJobs) Host {
	return Host{
		Enabled: func(action, _ string) bool { return *on && action == ActionKluster },
		Record:  func(e Entry) { *journal = append(*journal, e) },
		Kluster: jobs,
	}
}

// reqID is a request ID of its own for the n-th write of a test.
func reqID(n int) string { return fmt.Sprintf(`{"request_id":"req-%08d"}`, n) }

func runningStatus() kluster.Status {
	return kluster.Status{Name: "harbour", Profile: "full", Provider: "k3d", Context: "k3d-harbour", Exists: true, Running: true, Age: "2h50m39s", AgeSeconds: 10239}
}

func TestKlusterIsAHostActionTheDashboardSees(t *testing.T) {
	if _, ok := Actions[ActionKluster]; !ok || !DashboardSees(ActionKluster) {
		t.Errorf("kluster is a host action the dashboard sees and the operator enables with %s", EnableCommand(ActionKluster))
	}
	for _, name := range []string{"kluster.up", "kluster.down"} {
		if sp := klusterSpecs()[name]; sp.action != ActionKluster || sp.reads || sp.do == nil {
			t.Errorf("%s is a write under the kluster action, answered in flai serve: %+v", name, sp)
		}
	}
	if sp := klusterSpecs()["kluster.status"]; sp.action != "" || !sp.reads || sp.do == nil {
		t.Errorf("kluster.status is a read, answered whatever the action says: %+v", sp)
	}
}

func TestKlusterStatusReportsTheClusterWithTheActionOnOrOff(t *testing.T) {
	c := &fakeCluster{name: "harbour", status: runningStatus()}
	p := quayside(t, true, c)
	for _, on := range []bool{true, false} {
		journal := []Entry{}
		var got KlusterStatus
		if e := ask(t, klusterHost(&on, &journal, NewKlusterJobs(context.Background(), KlusterJobsOptions{})), p, "kluster.status", `{}`, &got); e != nil {
			t.Fatalf("on %v: %+v", on, e)
		}
		if got.Enabled != on || got.Enable != "flai serve enable kluster" || !got.Block || got.NoBlock != "" || got.StatusError != "" {
			t.Errorf("on %v: %+v", on, got)
		}
		if got.Cluster == nil || *got.Cluster != runningStatus() {
			t.Errorf("on %v: the cluster is kluster status's: %+v", on, got.Cluster)
		}
		if got.Running != nil || got.Last != nil {
			t.Errorf("on %v: no job ran: %+v %+v", on, got.Running, got.Last)
		}
		if len(journal) != 0 {
			t.Errorf("a read is not journalled: %v", journal)
		}
	}
}

func TestKlusterStatusWithNoBlockNamesTheBlockToAdd(t *testing.T) {
	on, journal := true, []Entry{}
	p := quayside(t, false, &fakeCluster{name: "harbour"})
	var got KlusterStatus
	if e := ask(t, klusterHost(&on, &journal, nil), p, "kluster.status", `{}`, &got); e != nil {
		t.Fatal(e)
	}
	if got.Block || got.NoBlock != kluster.ErrNoBlock.Error() || got.Cluster != nil {
		t.Errorf("no block: %+v", got)
	}
	if !strings.Contains(got.NoBlock, "kluster:") || !strings.Contains(got.NoBlock, "turn the kluster host action on") {
		t.Errorf("the sentence names the block and how to get one: %q", got.NoBlock)
	}
}

func TestKlusterStatusKeepsTheConfigsNamesWhenKlusterStatusFails(t *testing.T) {
	on, journal := true, []Entry{}
	p := quayside(t, true, &fakeCluster{name: "harbour", statusErr: kluster.ErrNotInstalled})
	var got KlusterStatus
	if e := ask(t, klusterHost(&on, &journal, nil), p, "kluster.status", `{}`, &got); e != nil {
		t.Fatal(e)
	}
	if got.Cluster == nil || got.Cluster.Name != "harbour" || got.Cluster.Profile != "full" || got.Cluster.Exists || got.StatusError != kluster.ErrNotInstalled.Error() {
		t.Errorf("kluster missing: %+v %+v", got, got.Cluster)
	}
}

func TestKlusterWritesAreRefusedWhileTheActionIsOff(t *testing.T) {
	on, journal := false, []Entry{}
	c := &fakeCluster{name: "harbour"}
	p := quayside(t, true, c)
	jobs := NewKlusterJobs(context.Background(), KlusterJobsOptions{})
	host := klusterHost(&on, &journal, jobs)
	for i, method := range []string{"kluster.up", "kluster.down"} {
		e := ask(t, host, p, method, reqID(i+1), nil)
		if e == nil || e.Code != Disabled || !strings.Contains(e.Message, "flai serve enable kluster") {
			t.Errorf("%s with the action off: %+v", method, e)
		}
	}
	jobs.Wait()
	if got := c.actions(); len(got) != 0 {
		t.Errorf("nothing runs while the action is off: %v", got)
	}
	if len(journal) != 2 || journal[0].Outcome != "disabled" || journal[0].Action != ActionKluster {
		t.Errorf("each refusal is journalled as disabled: %+v", journal)
	}
}

func TestKlusterWritesAreRefusedWithNoBlock(t *testing.T) {
	on, journal := true, []Entry{}
	c := &fakeCluster{name: "harbour"}
	p := quayside(t, false, c)
	jobs := NewKlusterJobs(context.Background(), KlusterJobsOptions{})
	for i, method := range []string{"kluster.up", "kluster.down"} {
		e := ask(t, klusterHost(&on, &journal, jobs), p, method, reqID(i+1), nil)
		if e == nil || e.Code != Rule || e.Message != kluster.ErrNoBlock.Error() {
			t.Errorf("%s with no block: %+v", method, e)
		}
	}
	jobs.Wait()
	if got := c.actions(); len(got) != 0 {
		t.Errorf("nothing runs with no block: %v", got)
	}
}

func TestKlusterWritesAreRefusedByAFlaiThatRunsNoJobs(t *testing.T) {
	on, journal := true, []Entry{}
	p := quayside(t, true, &fakeCluster{name: "harbour"})
	e := ask(t, klusterHost(&on, &journal, nil), p, "kluster.up", reqID(1), nil)
	if e == nil || e.Code != Rule || !strings.Contains(e.Message, "flai kluster up") {
		t.Errorf("no jobs here: %+v", e)
	}
}

func TestKlusterUpStartsAtOnceAndStatusSeesItRunThenEnd(t *testing.T) {
	on, journal := true, []Entry{}
	c := &fakeCluster{name: "harbour", status: runningStatus(), release: make(chan struct{})}
	p := quayside(t, true, c)
	jobs := NewKlusterJobs(context.Background(), KlusterJobsOptions{Now: func() time.Time { return t0 }})
	host := klusterHost(&on, &journal, jobs)

	var w struct {
		Data KlusterJob `json:"data"`
	}
	if e := ask(t, host, p, "kluster.up", reqID(1), &w); e != nil {
		t.Fatal(e)
	}
	started := w.Data
	want := KlusterJob{Action: "up", Cluster: "harbour", StartedAt: "2026-09-20T08:00:00Z", Log: ".flai-cache/kluster/up.log"}
	if started != want {
		t.Errorf("kluster.up answers the job it started: %+v, want %+v", started, want)
	}
	if len(journal) != 1 || journal[0].Outcome != "done" || !strings.Contains(journal[0].Detail, "started kluster up of the cluster harbour") {
		t.Errorf("the start is journalled: %+v", journal)
	}

	// while it runs, status says so, with its output so far
	var mid KlusterStatus
	for range 200 {
		if e := ask(t, host, p, "kluster.status", `{}`, &mid); e != nil {
			t.Fatal(e)
		}
		if mid.Running != nil && mid.Running.LogTail != "" {
			break
		}
		time.Sleep(5 * time.Millisecond)
	}
	if mid.Running == nil || mid.Running.Action != "up" || mid.Running.Outcome != "" || mid.Running.LogTail != "kluster up of harbour\n" || mid.Last != nil {
		t.Errorf("status while up runs: %+v %+v", mid.Running, mid.Last)
	}

	// a second job is refused while one runs
	e := ask(t, host, p, "kluster.down", reqID(2), nil)
	if e == nil || e.Code != Rule || !strings.Contains(e.Message, "kluster up of the cluster harbour runs since 2026-09-20T08:00:00Z") {
		t.Errorf("down while up runs: %+v", e)
	}

	close(c.release)
	jobs.Wait()
	var end KlusterStatus
	if e := ask(t, host, p, "kluster.status", `{}`, &end); e != nil {
		t.Fatal(e)
	}
	if end.Running != nil || end.Last == nil || end.Last.Outcome != "done" || end.Last.Error != "" || end.Last.EndedAt == "" || end.Last.LogTail != "kluster up of harbour\n" {
		t.Errorf("status once up ended: %+v %+v", end.Running, end.Last)
	}
	if got := c.actions(); len(got) != 1 || got[0] != "up" {
		t.Errorf("only the up ran: %v", got)
	}
}

func TestKlusterDownThatFailsIsSeenByStatusWithItsError(t *testing.T) {
	on, journal := true, []Entry{}
	c := &fakeCluster{name: "harbour", status: runningStatus()}
	p := quayside(t, true, c)
	c.downErr = errors.New("kluster down of the cluster harbour, from integration/kluster.yaml in " + p.Root + ", failed: exit status 1")
	jobs := NewKlusterJobs(context.Background(), KlusterJobsOptions{})
	host := klusterHost(&on, &journal, jobs)
	var w struct {
		Data KlusterJob `json:"data"`
	}
	if e := ask(t, host, p, "kluster.down", reqID(1), &w); e != nil {
		t.Fatal(e)
	}
	started := w.Data
	if started.Action != "down" || started.Log != ".flai-cache/kluster/down.log" {
		t.Errorf("kluster.down answers the job it started: %+v", started)
	}
	jobs.Wait()
	var got KlusterStatus
	if e := ask(t, host, p, "kluster.status", `{}`, &got); e != nil {
		t.Fatal(e)
	}
	if got.Last == nil || got.Last.Outcome != "failed" || got.Last.Error != "kluster down of the cluster harbour, from integration/kluster.yaml in the project's main checkout, failed: exit status 1" {
		t.Errorf("the failed down: %+v", got.Last)
	}
	// a job may start again once the last ended
	if e := ask(t, host, p, "kluster.up", reqID(2), nil); e != nil {
		t.Errorf("up after the down ended: %+v", e)
	}
	jobs.Wait()
}

func TestKlusterJobsEndWhenTheirContextEnds(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	c := &fakeCluster{name: "harbour", release: make(chan struct{})}
	root := t.TempDir()
	jobs := NewKlusterJobs(ctx, KlusterJobsOptions{})
	if _, err := jobs.Start(root, KlusterUp, c, time.Hour); err != nil {
		t.Fatal(err)
	}
	cancel()
	jobs.Wait()
	if _, last := jobs.State(root); last == nil || last.Outcome != "failed" || !strings.Contains(last.Error, "context canceled") {
		t.Errorf("a job cut short by flai serve's stop: %+v", last)
	}
}

func TestTailKeepsWholeLinesFromTheEnd(t *testing.T) {
	path := filepath.Join(t.TempDir(), "up.log")
	if err := os.WriteFile(path, []byte("one\ntwo\nthree\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if got := tail(path, 7); got != "three\n" {
		t.Errorf("tail 7: %q", got)
	}
	if got := tail(path, 100); got != "one\ntwo\nthree\n" {
		t.Errorf("tail 100: %q", got)
	}
	if got := tail(filepath.Join(t.TempDir(), "none"), 10); got != "" {
		t.Errorf("no file: %q", got)
	}
}

// berth is quayside with work items: S-0001 Gantry and S-0002 Bollards in
// review, S-0003 Fenders in backlog; S-0001 has a worktree. The story
// branches' and main's heads are those heads names.
func berth(t *testing.T, c *fakeCluster) (channel.Project, map[string]string) {
	t.Helper()
	p := quayside(t, true, c)
	for _, d := range []string{"wip/kanban/epics", "wip/kanban/stories", "wip/kanban/tasks", "wip/agents", "wip/archive", "wip/threads"} {
		if err := os.MkdirAll(filepath.Join(p.Root, d), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	repo, err := workitem.Open(p.Root)
	if err != nil {
		t.Fatal(err)
	}
	mk := func(typ, title, parent string) *workitem.Item {
		it, err := repo.Create(workitem.NewOptions{Type: typ, Title: title, Parent: parent, Owner: "olive", Now: t0})
		if err != nil {
			t.Fatal(err)
		}
		if typ == workitem.Story {
			data, _ := os.ReadFile(it.Path)
			_ = os.WriteFile(it.Path, []byte(strings.Replace(string(data), "## Acceptance criteria\n- [ ]\n", "## Acceptance criteria\n- [x] works\n", 1)), 0o644)
			it, _ = repo.Get(it.ID)
		}
		return it
	}
	move := func(it *workitem.Item, states ...string) {
		for i, to := range states {
			if _, err := repo.Transition(it, to, "olive", "", t0.Add(time.Duration(i+1)*time.Minute)); err != nil {
				t.Fatalf("%s to %s: %v", it.ID, to, err)
			}
		}
	}
	mk(workitem.Epic, "Quay", "")
	for _, title := range []string{"Gantry", "Bollards"} {
		s := mk(workitem.Story, title, "E-0001")
		task := mk(workitem.Task, "Rig "+title, s.ID)
		move(task, workitem.Ready, workitem.InProgress, workitem.Done)
		move(s, workitem.Ready, workitem.InProgress, workitem.Review)
	}
	mk(workitem.Story, "Fenders", "E-0001")
	if err := os.MkdirAll(repo.WorktreePath("S-0001"), 0o755); err != nil {
		t.Fatal(err)
	}
	heads := map[string]string{
		"refs/heads/main":           strings.Repeat("a", 40),
		"refs/heads/story/S-0001":   strings.Repeat("b", 40),
		"refs/heads/story/S-0002":   strings.Repeat("c", 40),
		repo.WorktreePath("S-0001"): strings.Repeat("b", 40),
	}
	wasHead, wasBranch := klusterHead, klusterMainBranch
	t.Cleanup(func() { klusterHead, klusterMainBranch = wasHead, wasBranch })
	klusterHead = func(_ context.Context, dir, rev string) (string, error) {
		if rev == "HEAD" {
			rev = dir
		}
		if h, ok := heads[rev]; ok {
			return h, nil
		}
		return "", fmt.Errorf("git rev-parse %s: exit status 1", rev)
	}
	klusterMainBranch = func(context.Context, string) (string, error) { return "main", nil }
	return p, heads
}

// recordAt writes r as its namespace's record under root, with a log for
// each step holding its name.
func recordAt(t *testing.T, root string, r kluster.Record) {
	t.Helper()
	for i, st := range r.Steps {
		path := kluster.LogPath(root, r.Namespace, st.Name)
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(st.Name+" of "+r.Namespace+"\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		r.Steps[i].Log = path
	}
	data, err := json.Marshal(r)
	if err != nil {
		t.Fatal(err)
	}
	path := kluster.RecordPath(root, r.Namespace)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatal(err)
	}
}

// kubeconfigAt writes a kubeconfig for the cluster run from root.
func kubeconfigAt(t *testing.T, root string) {
	t.Helper()
	path := kluster.KubeconfigPath(root)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte("apiVersion: v1\n"), 0o600); err != nil {
		t.Fatal(err)
	}
}

func passedSteps(names ...string) []kluster.Step {
	out := []kluster.Step{}
	for _, n := range names {
		out = append(out, kluster.Step{Name: n, StartedAt: t0, EndedAt: t0.Add(time.Minute), Outcome: kluster.OutcomePassed})
	}
	return out
}

func TestKlusterDeploymentMethodsAreAReadAndWritesUnderTheAction(t *testing.T) {
	for _, name := range []string{"kluster.deploy", "kluster.teardown", "kluster.promote"} {
		if sp := klusterSpecs()[name]; sp.action != ActionKluster || sp.reads || sp.do == nil || sp.describe == nil {
			t.Errorf("%s is a write under the kluster action, answered in flai serve: %+v", name, sp)
		}
	}
	if sp := klusterSpecs()["kluster.deployments"]; sp.action != "" || !sp.reads || sp.do == nil {
		t.Errorf("kluster.deployments is a read, answered whatever the action says: %+v", sp)
	}
}

func TestKlusterDeploymentsListsMainFirstThenEachStoryWithItsLogsAndAccess(t *testing.T) {
	c := &fakeCluster{name: "harbour", services: []KlusterService{
		{Namespace: "s-0001", Name: "web", Ports: []KlusterPort{{Name: "http", Port: 8080, Protocol: "TCP"}, {Name: "admin", Port: 9090, Protocol: "TCP"}}},
		{Namespace: "harbour", Name: "web", Ports: []KlusterPort{{Port: 8080, Protocol: "TCP"}}},
		{Namespace: "kube-system", Name: "kube-dns", Ports: []KlusterPort{{Port: 53, Protocol: "UDP"}}},
	}}
	p, heads := berth(t, c)
	kubeconfigAt(t, p.Root)
	recordAt(t, p.Root, kluster.Record{Story: "S-0001", Namespace: "s-0001", Commit: strings.Repeat("9", 40), Tag: "s-0001-999999999999", Checkout: filepath.Join(p.Root, ".flai-cache", "worktrees", "S-0001"), Images: []string{"harbour/web:s-0001-999999999999"}, StartedAt: t0, EndedAt: t0.Add(3 * time.Minute), Outcome: kluster.OutcomePassed, Steps: passedSteps(kluster.StepBuild, kluster.StepLoad, kluster.StepDeploy)})
	recordAt(t, p.Root, kluster.Record{Namespace: "harbour", Commit: heads["refs/heads/main"], Tag: "main-aaaaaaaaaaaa", Checkout: kluster.MainWorktree(p.Root), Images: []string{}, StartedAt: t0, EndedAt: t0.Add(2 * time.Minute), Outcome: kluster.OutcomeFailed, FailedStep: kluster.StepDeploy, ExitCode: 2, Error: "the deploy step for the namespace harbour failed; read its log " + kluster.LogPath(p.Root, "harbour", "deploy"), Steps: passedSteps(kluster.StepBuild, kluster.StepLoad, kluster.StepDeploy)})
	for _, on := range []bool{true, false} {
		journal := []Entry{}
		var got KlusterDeployments
		if e := ask(t, klusterHost(&on, &journal, NewKlusterJobs(context.Background(), KlusterJobsOptions{})), p, "kluster.deployments", `{}`, &got); e != nil {
			t.Fatalf("on %v: %+v", on, e)
		}
		if got.Enabled != on || got.Enable != "flai serve enable kluster" || !got.Block || got.Cluster != "harbour" || got.DefaultNamespace != "harbour" || got.Running != nil || len(journal) != 0 {
			t.Errorf("on %v: %+v", on, got)
		}
		if len(got.Deployments) != 2 {
			t.Fatalf("two deployments: %+v", got.Deployments)
		}
		main, story := got.Deployments[0], got.Deployments[1]
		if !main.Default || main.Namespace != "harbour" || main.Head != heads["refs/heads/main"] || main.Outcome != kluster.OutcomeFailed || main.FailedStep != kluster.StepDeploy || main.Checkout != ".flai-cache/kluster/main" || main.Error != "the deploy step for the namespace harbour failed; read its log .flai-cache/kluster/logs/harbour/deploy.log" {
			t.Errorf("main's namespace comes first, at main's head: %+v", main)
		}
		if story.Default || story.Story != "S-0001" || story.Title != "Gantry" || story.Status != workitem.Review || story.Head != heads["refs/heads/story/S-0001"] || story.Commit != strings.Repeat("9", 40) || story.Checkout != ".flai-cache/worktrees/S-0001" || story.HeadError != "" {
			t.Errorf("the story's namespace, deployed at an older commit than its head: %+v", story)
		}
		if len(story.Steps) != 3 || story.Steps[2].Name != "deploy" || story.Steps[2].Log != ".flai-cache/kluster/logs/s-0001/deploy.log" || story.Steps[2].LogTail != "deploy of s-0001\n" {
			t.Errorf("each step with its log and its end: %+v", story.Steps)
		}
		a := story.Access
		base := "kubectl --kubeconfig .flai-cache/kluster/kubeconfig --context k3d-harbour --namespace s-0001"
		if a.Kubeconfig != ".flai-cache/kluster/kubeconfig" || a.Context != "k3d-harbour" || a.Namespace != "s-0001" || a.List != base+" get services" || a.ServicesError != "" {
			t.Errorf("the story's access: %+v", a)
		}
		if len(a.Services) != 1 || a.Services[0].Name != "web" || a.Services[0].Namespace != "" || a.Services[0].Host != "web.s-0001.svc.cluster.local" || a.Services[0].PortForward != base+" port-forward service/web 8080:8080 9090:9090" || len(a.Services[0].Ports) != 2 {
			t.Errorf("the story's own services, each with its port-forward: %+v", a.Services)
		}
		if len(main.Access.Services) != 1 || main.Access.Services[0].Host != "web.harbour.svc.cluster.local" {
			t.Errorf("main's own services: %+v", main.Access.Services)
		}
		if len(got.Review) != 1 || got.Review[0] != (KlusterReview{Story: "S-0002", Title: "Bollards", Head: heads["refs/heads/story/S-0002"]}) {
			t.Errorf("the story in review with no deployment: %+v", got.Review)
		}
	}
	if c.listedTimes != 2 {
		t.Errorf("one kubectl call per list, for every namespace at once: %d", c.listedTimes)
	}
}

func TestKlusterDeploymentsOfAnEmptyClusterListsTheStoriesInReview(t *testing.T) {
	c := &fakeCluster{name: "harbour"}
	p, heads := berth(t, c)
	on, journal := true, []Entry{}
	var got KlusterDeployments
	if e := ask(t, klusterHost(&on, &journal, nil), p, "kluster.deployments", `{}`, &got); e != nil {
		t.Fatal(e)
	}
	if got.Deployments == nil || len(got.Deployments) != 0 || got.Running != nil || got.Last != nil {
		t.Errorf("no deployments: %+v", got)
	}
	want := []KlusterReview{{Story: "S-0001", Title: "Gantry", Head: heads["refs/heads/story/S-0001"]}, {Story: "S-0002", Title: "Bollards", Head: heads["refs/heads/story/S-0002"]}}
	if len(got.Review) != 2 || got.Review[0] != want[0] || got.Review[1] != want[1] {
		t.Errorf("both stories in review can be deployed: %+v", got.Review)
	}
	if c.listedTimes != 0 {
		t.Errorf("no namespace, no kubectl: %d", c.listedTimes)
	}
}

func TestKlusterDeploymentsSaysWhyServicesCannotBeReadWithNoKubeconfig(t *testing.T) {
	c := &fakeCluster{name: "harbour"}
	p, _ := berth(t, c)
	recordAt(t, p.Root, kluster.Record{Story: "S-0001", Namespace: "s-0001", Commit: strings.Repeat("b", 40), StartedAt: t0, Outcome: kluster.OutcomePassed, Steps: passedSteps(kluster.StepDeploy)})
	on, journal := true, []Entry{}
	var got KlusterDeployments
	if e := ask(t, klusterHost(&on, &journal, nil), p, "kluster.deployments", `{}`, &got); e != nil {
		t.Fatal(e)
	}
	if len(got.Deployments) != 1 || !strings.Contains(got.Deployments[0].Access.ServicesError, "flai kluster up") || got.Deployments[0].Access.List == "" || len(got.Deployments[0].Access.Services) != 0 || c.listedTimes != 0 {
		t.Errorf("no kubeconfig: %+v", got.Deployments)
	}
}

func TestKlusterDeploymentsWithNoBlockNamesTheBlockToAdd(t *testing.T) {
	on, journal := true, []Entry{}
	p := quayside(t, false, &fakeCluster{name: "harbour"})
	for _, d := range []string{"wip/kanban/epics", "wip/kanban/stories", "wip/kanban/tasks"} {
		_ = os.MkdirAll(filepath.Join(p.Root, d), 0o755)
	}
	var got KlusterDeployments
	if e := ask(t, klusterHost(&on, &journal, nil), p, "kluster.deployments", `{}`, &got); e != nil {
		t.Fatal(e)
	}
	if got.Block || got.NoBlock != kluster.ErrNoBlock.Error() || len(got.Deployments) != 0 {
		t.Errorf("no block: %+v", got)
	}
}

func TestKlusterDeploymentWritesAreRefusedWhileTheActionIsOff(t *testing.T) {
	on, journal := false, []Entry{}
	c := &fakeCluster{name: "harbour"}
	p, _ := berth(t, c)
	jobs := NewKlusterJobs(context.Background(), KlusterJobsOptions{})
	host := klusterHost(&on, &journal, jobs)
	for i, call := range []struct{ method, params string }{
		{"kluster.deploy", `{"story":"S-0001","request_id":"req-00000001"}`},
		{"kluster.teardown", `{"story":"S-0001","request_id":"req-00000002"}`},
		{"kluster.promote", reqID(3)},
	} {
		e := ask(t, host, p, call.method, call.params, nil)
		if e == nil || e.Code != Disabled || !strings.Contains(e.Message, "flai serve enable kluster") {
			t.Errorf("%d %s with the action off: %+v", i, call.method, e)
		}
	}
	jobs.Wait()
	if got := c.called(); len(got) != 0 {
		t.Errorf("nothing runs while the action is off: %v", got)
	}
	if len(journal) != 3 || journal[0].Outcome != "disabled" || journal[0].Action != ActionKluster {
		t.Errorf("each refusal is journalled as disabled: %+v", journal)
	}
}

func TestKlusterDeployStartsTheStorysDeployFromItsWorktreeAtItsHead(t *testing.T) {
	on, journal := true, []Entry{}
	c := &fakeCluster{name: "harbour", release: make(chan struct{})}
	p, heads := berth(t, c)
	jobs := NewKlusterJobs(context.Background(), KlusterJobsOptions{Now: func() time.Time { return t0 }})
	host := klusterHost(&on, &journal, jobs)

	var w struct {
		Data KlusterJob `json:"data"`
	}
	if e := ask(t, host, p, "kluster.deploy", `{"story":"s-1","request_id":"req-00000001"}`, &w); e != nil {
		t.Fatal(e)
	}
	want := KlusterJob{Action: "deploy", Cluster: "harbour", Story: "S-0001", Namespace: "s-0001", StartedAt: "2026-09-20T08:00:00Z", Log: ".flai-cache/kluster/deploy.log", Logs: ".flai-cache/kluster/logs/s-0001"}
	if w.Data != want {
		t.Errorf("kluster.deploy answers the job it started: %+v, want %+v", w.Data, want)
	}
	if len(journal) != 1 || journal[0].Outcome != "done" || !strings.Contains(journal[0].Detail, "started kluster deploy of S-0001 in the namespace s-0001") {
		t.Errorf("the start is journalled: %+v", journal)
	}

	// while it runs, the list shows it, with its progress so far
	var mid KlusterDeployments
	for range 200 {
		if e := ask(t, host, p, "kluster.deployments", `{}`, &mid); e != nil {
			t.Fatal(e)
		}
		if mid.Running != nil && mid.Running.LogTail != "" {
			break
		}
		time.Sleep(5 * time.Millisecond)
	}
	if mid.Running == nil || mid.Running.Action != "deploy" || mid.Running.Story != "S-0001" || mid.Running.LogTail != "kluster deploy of harbour\n" || mid.Last != nil {
		t.Errorf("the list while the deploy runs: %+v %+v", mid.Running, mid.Last)
	}
	var st KlusterStatus
	if e := ask(t, host, p, "kluster.status", `{}`, &st); e != nil {
		t.Fatal(e)
	}
	if st.Last != nil {
		t.Errorf("a deploy is not the last up or down: %+v", st.Last)
	}

	// every other job is refused while it runs
	recordAt(t, p.Root, kluster.Record{Story: "S-0001", Namespace: "s-0001", Commit: strings.Repeat("b", 40), StartedAt: t0, Outcome: kluster.OutcomeRunning, Steps: []kluster.Step{}})
	for i, call := range []struct{ method, params string }{
		{"kluster.teardown", `{"story":"S-0001","request_id":"req-00000002"}`},
		{"kluster.promote", reqID(3)},
		{"kluster.up", reqID(4)},
	} {
		e := ask(t, host, p, call.method, call.params, nil)
		if e == nil || e.Code != Rule || !strings.Contains(e.Message, "kluster deploy of the cluster harbour runs since 2026-09-20T08:00:00Z") {
			t.Errorf("%d %s while the deploy runs: %+v", i, call.method, e)
		}
	}

	close(c.release)
	jobs.Wait()
	if got := c.called(); len(got) != 1 || got[0] != "deploy S-0001 review "+filepath.Join(p.Root, ".flai-cache", "worktrees", "S-0001")+" "+heads["refs/heads/story/S-0001"] {
		t.Errorf("the deploy ran as flai kluster deploy runs it: %v", got)
	}
	var end KlusterDeployments
	if e := ask(t, host, p, "kluster.deployments", `{}`, &end); e != nil {
		t.Fatal(e)
	}
	if end.Running != nil || end.Last == nil || end.Last.Outcome != "done" || end.Last.Action != "deploy" {
		t.Errorf("the list once the deploy ended: %+v %+v", end.Running, end.Last)
	}
}

func TestKlusterDeployRefusesWhatFlaiKlusterDeployRefuses(t *testing.T) {
	on, journal := true, []Entry{}
	c := &fakeCluster{name: "harbour"}
	p, _ := berth(t, c)
	jobs := NewKlusterJobs(context.Background(), KlusterJobsOptions{})
	host := klusterHost(&on, &journal, jobs)
	cases := []struct{ params, want string }{
		{`{"story":"S-0003","request_id":"req-00000001"}`, "S-0003 is backlog, and flai deploys only a story in progress or in review"},
		{`{"story":"S-0002","request_id":"req-00000002"}`, "S-0002 has no worktree at .flai-cache/worktrees/S-0002 to build and deploy from; open it with flai stream open S-0002"},
	}
	for _, tc := range cases {
		e := ask(t, host, p, "kluster.deploy", tc.params, nil)
		if e == nil || e.Code != Rule || !strings.Contains(e.Message, tc.want) {
			t.Errorf("%s: %+v", tc.params, e)
		}
	}
	if e := ask(t, host, p, "kluster.deploy", `{"story":"E-0001","request_id":"req-00000003"}`, nil); e == nil || e.Code != channel.CodeInvalidParams {
		t.Errorf("an epic's ID: %+v", e)
	}
	c.readyErr = errors.New("the cluster harbour has no kubeconfig at " + kluster.KubeconfigPath(p.Root) + "; bring it up with flai kluster up, then run this again")
	if e := ask(t, host, p, "kluster.deploy", `{"story":"S-0001","request_id":"req-00000004"}`, nil); e == nil || e.Code != Rule || e.Message != "the cluster harbour has no kubeconfig at .flai-cache/kluster/kubeconfig; bring it up with flai kluster up, then run this again" {
		t.Errorf("no kubeconfig: %+v", e)
	}
	jobs.Wait()
	if got := c.called(); len(got) != 0 {
		t.Errorf("nothing ran: %v", got)
	}
}

func TestKlusterDeploymentWritesAreRefusedByAFlaiThatRunsNoJobs(t *testing.T) {
	on, journal := true, []Entry{}
	p, _ := berth(t, &fakeCluster{name: "harbour"})
	e := ask(t, klusterHost(&on, &journal, nil), p, "kluster.promote", reqID(1), nil)
	if e == nil || e.Code != Rule || !strings.Contains(e.Message, "run flai kluster promote in a shell") {
		t.Errorf("no jobs here: %+v", e)
	}
	e = ask(t, klusterHost(&on, &journal, nil), p, "kluster.deploy", `{"story":"S-0001","request_id":"req-00000002"}`, nil)
	if e == nil || e.Code != Rule || !strings.Contains(e.Message, "run flai kluster deploy S-0001 in a shell") {
		t.Errorf("no jobs here: %+v", e)
	}
}

func TestKlusterTeardownStartsOnlyForARecordedDeployment(t *testing.T) {
	on, journal := true, []Entry{}
	c := &fakeCluster{name: "harbour"}
	p, _ := berth(t, c)
	jobs := NewKlusterJobs(context.Background(), KlusterJobsOptions{})
	host := klusterHost(&on, &journal, jobs)
	if e := ask(t, host, p, "kluster.teardown", `{"story":"S-0002","request_id":"req-00000001"}`, nil); e == nil || e.Code != Rule || e.Message != "S-0002 has no deployment recorded; nothing to tear down" {
		t.Errorf("no record: %+v", e)
	}
	recordAt(t, p.Root, kluster.Record{Story: "S-0001", Namespace: "s-0001", Commit: strings.Repeat("b", 40), StartedAt: t0, Outcome: kluster.OutcomePassed, Steps: passedSteps(kluster.StepDeploy)})
	var w struct {
		Data KlusterJob `json:"data"`
	}
	if e := ask(t, host, p, "kluster.teardown", `{"story":"S-0001","request_id":"req-00000002"}`, &w); e != nil {
		t.Fatal(e)
	}
	if w.Data.Action != "teardown" || w.Data.Story != "S-0001" || w.Data.Namespace != "s-0001" || w.Data.Log != ".flai-cache/kluster/teardown.log" {
		t.Errorf("kluster.teardown answers the job it started: %+v", w.Data)
	}
	jobs.Wait()
	if got := c.called(); len(got) != 1 || got[0] != "teardown S-0001 "+filepath.Join(p.Root, ".flai-cache", "worktrees", "S-0001") {
		t.Errorf("the teardown ran in the story's worktree: %v", got)
	}
}

func TestKlusterPromoteStartsMainsDeployIntoTheDefaultNamespace(t *testing.T) {
	on, journal := true, []Entry{}
	c := &fakeCluster{name: "harbour"}
	p, _ := berth(t, c)
	c.deployErr = errors.New("the deploy step for the namespace harbour failed with exit code 2; read its log " + kluster.LogPath(p.Root, "harbour", "deploy"))
	jobs := NewKlusterJobs(context.Background(), KlusterJobsOptions{})
	host := klusterHost(&on, &journal, jobs)
	var w struct {
		Data KlusterJob `json:"data"`
	}
	if e := ask(t, host, p, "kluster.promote", reqID(1), &w); e != nil {
		t.Fatal(e)
	}
	if w.Data.Action != "promote" || w.Data.Story != "" || w.Data.Namespace != "harbour" || w.Data.Logs != ".flai-cache/kluster/logs/harbour" {
		t.Errorf("kluster.promote answers the job it started: %+v", w.Data)
	}
	jobs.Wait()
	if got := c.called(); len(got) != 1 || got[0] != "promote main" {
		t.Errorf("main was promoted: %v", got)
	}
	var end KlusterDeployments
	if e := ask(t, host, p, "kluster.deployments", `{}`, &end); e != nil {
		t.Fatal(e)
	}
	if end.Last == nil || end.Last.Outcome != "failed" || end.Last.Error != "the deploy step for the namespace harbour failed with exit code 2; read its log .flai-cache/kluster/logs/harbour/deploy.log" {
		t.Errorf("the failed promotion, its error from the root: %+v", end.Last)
	}
}

func TestParseServicesReadsKubectlsJSON(t *testing.T) {
	got, err := parseServices([]byte(`{"items":[{"metadata":{"name":"web","namespace":"s-0001"},"spec":{"ports":[{"name":"http","port":8080,"protocol":"TCP","targetPort":3000}]}},{"metadata":{"name":"headless","namespace":"s-0001"},"spec":{}}]}`))
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 2 || got[0].Name != "web" || got[0].Namespace != "s-0001" || len(got[0].Ports) != 1 || got[0].Ports[0] != (KlusterPort{Name: "http", Port: 8080, Protocol: "TCP"}) || got[1].Ports == nil {
		t.Errorf("services: %+v", got)
	}
	a := klusterAccess("k3d-harbour", "s-0001", got, "")
	if len(a.Services) != 2 || a.Services[1].PortForward != "" || a.Services[0].PortForward == "" {
		t.Errorf("a service with no port has no port-forward: %+v", a.Services)
	}
	if _, err := parseServices([]byte("no")); err == nil {
		t.Error("JSON that is not kubectl's is an error")
	}
}
