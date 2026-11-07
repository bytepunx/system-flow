package hostapi

import (
	"context"
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
)

// fakeCluster is a cluster that runs no kluster: Up and Down write a line to
// their output, wait for release when it is set, and return their error.
type fakeCluster struct {
	name      string
	status    kluster.Status
	statusErr error
	upErr     error
	downErr   error
	release   chan struct{}
	mu        sync.Mutex
	ran       []string
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
	if !strings.Contains(got.NoBlock, "kluster:") || !strings.Contains(got.NoBlock, "design/system/kluster.md") {
		t.Errorf("the sentence names the block and where it is described: %q", got.NoBlock)
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
