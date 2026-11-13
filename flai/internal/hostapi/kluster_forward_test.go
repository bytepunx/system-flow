package hostapi

import (
	"context"
	"os/exec"
	"slices"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/channel"
	"github.com/bytepunx/system-flow/flai/internal/kluster"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// fakeKubectl stands in for kubectl port-forward: each start runs script in
// sh, its arguments recorded, on the local port 41234.
type fakeKubectl struct {
	mu   sync.Mutex
	args [][]string
}

func (f *fakeKubectl) started() [][]string {
	f.mu.Lock()
	defer f.mu.Unlock()
	return append([][]string{}, f.args...)
}

// forwarding is a kubectl that says it forwards, then waits to be stopped.
const forwarding = `echo "Forwarding from 0.0.0.0:41234 -> 8080"; exec sleep 60`

func fakeForwards(t *testing.T, script string) *fakeKubectl {
	t.Helper()
	f := &fakeKubectl{}
	wasCmd, wasPort := klusterForwardCommand, klusterFreePort
	t.Cleanup(func() { klusterForwardCommand, klusterFreePort = wasCmd, wasPort })
	klusterFreePort = func() (int, error) { return 41234, nil }
	klusterForwardCommand = func(ctx context.Context, root string, args []string) *exec.Cmd {
		f.mu.Lock()
		f.args = append(f.args, args)
		f.mu.Unlock()
		cmd := exec.CommandContext(ctx, "sh", "-c", script)
		cmd.Dir = root
		cmd.WaitDelay = time.Second
		return cmd
	}
	return f
}

// forwardBerth is berth's project with a kubeconfig, a deployment of S-0001
// recorded, and the service web on 8080 and 9090 in its namespace.
func forwardBerth(t *testing.T) (channel.Project, *fakeCluster) {
	t.Helper()
	c := &fakeCluster{name: "harbour", services: []KlusterService{
		{Namespace: "s-0001", Name: "web", Ports: []KlusterPort{{Port: 8080}, {Port: 9090}}},
	}}
	p, _ := berth(t, c)
	kubeconfigAt(t, p.Root)
	recordAt(t, p.Root, kluster.Record{Story: "S-0001", Namespace: "s-0001", Commit: strings.Repeat("b", 40), StartedAt: t0, Outcome: kluster.OutcomePassed, Steps: passedSteps(kluster.StepDeploy)})
	return p, c
}

func TestKlusterForwardMethodsAreWritesUnderTheAction(t *testing.T) {
	for _, name := range []string{"kluster.forward.start", "kluster.forward.stop"} {
		if sp := klusterSpecs()[name]; sp.action != ActionKluster || sp.reads || sp.do == nil || sp.describe == nil {
			t.Errorf("%s is a write under the kluster action, answered in flai serve: %+v", name, sp)
		}
	}
}

func TestAForwardStartsOn0000AndIsListedThenStopped(t *testing.T) {
	kc := fakeForwards(t, forwarding)
	p, _ := forwardBerth(t)
	on, journal := true, []Entry{}
	jobs := NewKlusterJobs(context.Background(), KlusterJobsOptions{Now: func() time.Time { return t0 }})
	host := klusterHost(&on, &journal, jobs)
	var w struct {
		Data KlusterForward `json:"data"`
	}
	if e := ask(t, host, p, "kluster.forward.start", `{"namespace":"s-0001","service":"web","port":8080,"request_id":"req-00000001"}`, &w); e != nil {
		t.Fatal(e)
	}
	if w.Data != (KlusterForward{Namespace: "s-0001", Service: "web", Port: 8080, LocalPort: 41234, Address: "0.0.0.0", StartedAt: t0.UTC().Format(time.RFC3339), Running: true}) {
		t.Errorf("the forward, on the free port, bound to 0.0.0.0: %+v", w.Data)
	}
	got := kc.started()
	if len(got) != 1 {
		t.Fatalf("one kubectl: %v", got)
	}
	args := got[0]
	want := []string{"--context", "k3d-harbour", "--namespace", "s-0001", "port-forward", "--address", "0.0.0.0", "service/web", "41234:8080"}
	if len(args) != 2+len(want) || args[0] != "--kubeconfig" || args[1] != kluster.KubeconfigPath(p.Root) || !slices.Equal(args[2:], want) {
		t.Errorf("kubectl port-forward with flai's kubeconfig, bound to 0.0.0.0: %v", args)
	}

	// a second start answers the forward that runs, and starts no kubectl
	if e := ask(t, host, p, "kluster.forward.start", `{"namespace":"s-0001","service":"web","port":8080,"request_id":"req-00000002"}`, &w); e != nil {
		t.Fatal(e)
	}
	if !w.Data.Running || w.Data.LocalPort != 41234 || len(kc.started()) != 1 {
		t.Errorf("the running forward, again: %+v, %d starts", w.Data, len(kc.started()))
	}

	var d KlusterDeployments
	if e := ask(t, host, p, "kluster.deployments", `{}`, &d); e != nil {
		t.Fatal(e)
	}
	var fws []KlusterForward
	for _, row := range d.Deployments {
		if row.Namespace == "s-0001" {
			fws = row.Access.Forwards
		}
	}
	if len(fws) != 1 || fws[0].Service != "web" || fws[0].Port != 8080 || fws[0].LocalPort != 41234 || !fws[0].Running {
		t.Errorf("kluster.deployments lists the forward under its namespace: %+v", fws)
	}

	if e := ask(t, host, p, "kluster.forward.stop", `{"namespace":"s-0001","service":"web","port":8080,"request_id":"req-00000003"}`, &w); e != nil {
		t.Fatal(e)
	}
	if w.Data.Running || w.Data.EndedAt == "" {
		t.Errorf("the forward, stopped: %+v", w.Data)
	}
	if left := jobs.Forwards(p.Root, "s-0001"); len(left) != 0 {
		t.Errorf("a stopped forward is forgotten: %+v", left)
	}
	jobs.Wait()
	if len(journal) != 3 || journal[0].Action != ActionKluster || journal[0].Outcome != "done" || !strings.Contains(journal[0].Detail, "0.0.0.0:41234") {
		t.Errorf("each start and stop is journalled: %+v", journal)
	}
}

func TestASecondStartWaitsUntilTheFirstForwards(t *testing.T) {
	kc := fakeForwards(t, `sleep 0.4; echo "Forwarding from 0.0.0.0:41234 -> 8080"; exec sleep 60`)
	p, _ := forwardBerth(t)
	jobs := NewKlusterJobs(context.Background(), KlusterJobsOptions{})
	t.Cleanup(func() { jobs.StopForwards(p.Root, "s-0001"); jobs.Wait() })
	first := make(chan error, 1)
	go func() {
		_, err := jobs.StartForward(p.Root, "k3d-harbour", "s-0001", "web", 8080)
		first <- err
	}()
	for deadline := time.Now().Add(5 * time.Second); len(kc.started()) == 0; {
		if time.Now().After(deadline) {
			t.Fatal("the first start ran no kubectl")
		}
		time.Sleep(10 * time.Millisecond)
	}
	began := time.Now()
	got, err := jobs.StartForward(p.Root, "k3d-harbour", "s-0001", "web", 8080)
	if err != nil {
		t.Fatal(err)
	}
	if waited := time.Since(began); waited < 200*time.Millisecond {
		t.Errorf("the second start answered after %s, before kubectl said it forwards", waited)
	}
	if !got.Running || got.LocalPort != 41234 || len(kc.started()) != 1 {
		t.Errorf("the second start answers the first's forward and runs no kubectl: %+v, %d starts", got, len(kc.started()))
	}
	if err := <-first; err != nil {
		t.Fatal(err)
	}
}

func TestAForwardWhoseKubectlExitsIsKeptEndedWithItsOutput(t *testing.T) {
	fakeForwards(t, `echo "Forwarding from 0.0.0.0:41234 -> 8080"; sleep 0.2; echo "lost connection to pod"; exit 1`)
	p, _ := forwardBerth(t)
	on, journal := true, []Entry{}
	jobs := NewKlusterJobs(context.Background(), KlusterJobsOptions{})
	host := klusterHost(&on, &journal, jobs)
	if e := ask(t, host, p, "kluster.forward.start", `{"namespace":"s-0001","service":"web","port":8080,"request_id":"req-00000001"}`, nil); e != nil {
		t.Fatal(e)
	}
	jobs.Wait()
	fws := jobs.Forwards(p.Root, "s-0001")
	if len(fws) != 1 || fws[0].Running || fws[0].EndedAt == "" || !strings.Contains(fws[0].Output, "lost connection to pod") {
		t.Errorf("an ended forward is kept with the end of its output: %+v", fws)
	}
}

func TestAForwardThatExitsBeforeForwardingFailsWithItsOutput(t *testing.T) {
	fakeForwards(t, `echo "error: services \"web\" not found" >&2; exit 1`)
	p, _ := forwardBerth(t)
	on, journal := true, []Entry{}
	jobs := NewKlusterJobs(context.Background(), KlusterJobsOptions{})
	e := ask(t, klusterHost(&on, &journal, jobs), p, "kluster.forward.start", `{"namespace":"s-0001","service":"web","port":8080,"request_id":"req-00000001"}`, nil)
	if e == nil || !strings.Contains(e.Message, "exited before it forwarded") || !strings.Contains(e.Message, `services "web" not found`) {
		t.Errorf("kubectl's own words: %+v", e)
	}
	jobs.Wait()
}

func TestAForwardIsRefusedWhatItCannotReach(t *testing.T) {
	kc := fakeForwards(t, forwarding)
	p, _ := forwardBerth(t)
	on, journal := true, []Entry{}
	jobs := NewKlusterJobs(context.Background(), KlusterJobsOptions{})
	host := klusterHost(&on, &journal, jobs)
	for i, c := range []struct {
		method, params string
		code           int
		says           string
	}{
		{"kluster.forward.start", `{"namespace":"S_0001","service":"web","port":8080,"request_id":"req-00000001"}`, channel.CodeInvalidParams, "is not a namespace"},
		{"kluster.forward.start", `{"namespace":"s-0001","service":"web","port":0,"request_id":"req-00000002"}`, channel.CodeInvalidParams, "is not a port"},
		{"kluster.forward.start", `{"namespace":"s-0001","service":"api","port":8080,"request_id":"req-00000003"}`, Rule, "lists no service api with the port 8080"},
		{"kluster.forward.start", `{"namespace":"s-0001","service":"web","port":7000,"request_id":"req-00000004"}`, Rule, "lists no service web with the port 7000"},
		{"kluster.forward.stop", `{"namespace":"s-0001","service":"web","port":8080,"request_id":"req-00000005"}`, Rule, "nothing to stop"},
	} {
		e := ask(t, host, p, c.method, c.params, nil)
		if e == nil || e.Code != c.code || !strings.Contains(e.Message, c.says) {
			t.Errorf("%d %s %s: %+v", i, c.method, c.params, e)
		}
	}
	if len(kc.started()) != 0 {
		t.Errorf("no kubectl for a refused forward: %v", kc.started())
	}
}

func TestAForwardIsRefusedWithNoKubeconfigOrWithTheActionOff(t *testing.T) {
	kc := fakeForwards(t, forwarding)
	c := &fakeCluster{name: "harbour", services: []KlusterService{{Namespace: "s-0001", Name: "web", Ports: []KlusterPort{{Port: 8080}}}}}
	p, _ := berth(t, c)
	on, journal := true, []Entry{}
	jobs := NewKlusterJobs(context.Background(), KlusterJobsOptions{})
	host := klusterHost(&on, &journal, jobs)
	e := ask(t, host, p, "kluster.forward.start", `{"namespace":"s-0001","service":"web","port":8080,"request_id":"req-00000001"}`, nil)
	if e == nil || e.Code != Rule || !strings.Contains(e.Message, "no kubeconfig yet") {
		t.Errorf("no kubeconfig: %+v", e)
	}
	on = false
	for i, method := range []string{"kluster.forward.start", "kluster.forward.stop"} {
		e := ask(t, host, p, method, `{"namespace":"s-0001","service":"web","port":8080,"request_id":"req-0000001`+string(rune('2'+i))+`"}`, nil)
		if e == nil || e.Code != Disabled || !strings.Contains(e.Message, "flai serve enable kluster") {
			t.Errorf("%s with the action off: %+v", method, e)
		}
	}
	if e := ask(t, klusterHost(&on, &journal, nil), p, "kluster.forward.stop", `{"namespace":"s-0001","service":"web","port":8080,"request_id":"req-00000009"}`, nil); e == nil {
		t.Error("a flai that runs no jobs runs no forwards")
	}
	if len(kc.started()) != 0 {
		t.Errorf("no kubectl: %v", kc.started())
	}
}

func TestATeardownStopsItsNamespacesForwards(t *testing.T) {
	fakeForwards(t, forwarding)
	p, c := forwardBerth(t)
	on, journal := true, []Entry{}
	jobs := NewKlusterJobs(context.Background(), KlusterJobsOptions{})
	host := klusterHost(&on, &journal, jobs)
	for i, port := range []string{"8080", "9090"} {
		if e := ask(t, host, p, "kluster.forward.start", `{"namespace":"s-0001","service":"web","port":`+port+`,"request_id":"req-0000000`+string(rune('1'+i))+`"}`, nil); e != nil {
			t.Fatal(e)
		}
	}
	if n := len(jobs.Forwards(p.Root, "s-0001")); n != 2 {
		t.Fatalf("two forwards: %d", n)
	}
	if e := ask(t, host, p, "kluster.teardown", `{"story":"S-0001","request_id":"req-00000003"}`, nil); e != nil {
		t.Fatal(e)
	}
	jobs.Wait()
	if left := jobs.Forwards(p.Root, "s-0001"); len(left) != 0 {
		t.Errorf("the teardown stopped the namespace's forwards: %+v", left)
	}
	if got := c.called(); len(got) != 1 || !strings.HasPrefix(got[0], "teardown S-0001") {
		t.Errorf("the teardown ran: %v", got)
	}
}

// S-0663: a forward into a story in progress's namespace starts and stops as
// one into a story in review's does: neither looks at the story's column.
func TestAForwardIntoAStoryInProgressStartsAndStopsAsOneInReview(t *testing.T) {
	kc := fakeForwards(t, forwarding)
	p, c := forwardBerth(t)
	repo, err := workitem.Open(p.Root)
	if err != nil {
		t.Fatal(err)
	}
	// the heads do not matter to a forward
	id, head := capstanAt(t, p, map[string]string{})
	recordAt(t, p.Root, kluster.Record{Story: id, Namespace: "s-0004", Commit: head, StartedAt: t0, Outcome: kluster.OutcomePassed, Steps: passedSteps(kluster.StepDeploy)})
	c.services = append(c.services, KlusterService{Namespace: "s-0004", Name: "web", Ports: []KlusterPort{{Port: 8080}}})
	on, journal := true, []Entry{}
	jobs := NewKlusterJobs(context.Background(), KlusterJobsOptions{Now: func() time.Time { return t0 }})
	host := klusterHost(&on, &journal, jobs)
	for i, tc := range []struct{ story, ns, column string }{{"S-0001", "s-0001", workitem.Review}, {id, "s-0004", workitem.InProgress}} {
		if it, err := repo.Get(tc.story); err != nil || it.Status != tc.column {
			t.Fatalf("%s is %s: %v", tc.story, tc.column, err)
		}
		params := `{"namespace":"` + tc.ns + `","service":"web","port":8080,"request_id":"req-0000000`
		var w struct {
			Data KlusterForward `json:"data"`
		}
		if e := ask(t, host, p, "kluster.forward.start", params+string(rune('1'+2*i))+`"}`, &w); e != nil {
			t.Fatalf("%s: %+v", tc.column, e)
		}
		if w.Data != (KlusterForward{Namespace: tc.ns, Service: "web", Port: 8080, LocalPort: 41234, Address: "0.0.0.0", StartedAt: t0.UTC().Format(time.RFC3339), Running: true}) {
			t.Errorf("%s: the forward, on the free port, bound to 0.0.0.0: %+v", tc.column, w.Data)
		}
		if got := kc.started(); len(got) != i+1 || !slices.Contains(got[i], tc.ns) || got[i][len(got[i])-1] != "41234:8080" {
			t.Errorf("%s: kubectl port-forward into %s: %v", tc.column, tc.ns, got)
		}
		var d KlusterDeployments
		if e := ask(t, host, p, "kluster.deployments", `{}`, &d); e != nil {
			t.Fatal(e)
		}
		listed := false
		for _, row := range d.Deployments {
			if row.Namespace == tc.ns {
				listed = row.Column == tc.column && len(row.Access.Forwards) == 1 && row.Access.Forwards[0].Running
			}
		}
		if !listed {
			t.Errorf("%s: kluster.deployments lists the forward under its namespace: %+v", tc.column, d.Deployments)
		}
		if e := ask(t, host, p, "kluster.forward.stop", params+string(rune('2'+2*i))+`"}`, &w); e != nil {
			t.Fatalf("%s: %+v", tc.column, e)
		}
		if w.Data.Running || w.Data.EndedAt == "" {
			t.Errorf("%s: the forward, stopped: %+v", tc.column, w.Data)
		}
		if left := jobs.Forwards(p.Root, tc.ns); len(left) != 0 {
			t.Errorf("%s: a stopped forward is forgotten: %+v", tc.column, left)
		}
	}
	jobs.Wait()
	if len(journal) != 4 {
		t.Errorf("each start and stop is journalled: %+v", journal)
	}
}

func TestForwardsEndWhenTheirContextEnds(t *testing.T) {
	fakeForwards(t, forwarding)
	p, _ := forwardBerth(t)
	ctx, cancel := context.WithCancel(context.Background())
	jobs := NewKlusterJobs(ctx, KlusterJobsOptions{})
	if _, err := jobs.StartForward(p.Root, "k3d-harbour", "s-0001", "web", 8080); err != nil {
		t.Fatal(err)
	}
	cancel()
	jobs.Wait()
	if fws := jobs.Forwards(p.Root, "s-0001"); len(fws) != 1 || fws[0].Running {
		t.Errorf("flai serve stopping stops its forwards: %+v", fws)
	}
}
