package host

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"
)

// The children in these tests are this test binary, told by childEnv what to
// be: one that runs until it is asked to stop, one that ends at once, and
// one that will not stop when asked.
const childEnv = "FLAI_HOST_TEST_CHILD"

func TestMain(m *testing.M) {
	switch os.Getenv(childEnv) {
	case "":
		os.Exit(m.Run())
	case "run":
		if f := os.Getenv("FLAI_HOST_TEST_ENVFILE"); f != "" {
			_ = os.WriteFile(f, []byte(os.Getenv(URLEnv)+"\n"+os.Getenv(TokenEnv)), 0o600)
		}
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGTERM, os.Interrupt)
		select {
		case <-c:
		case <-time.After(time.Minute):
		}
		os.Exit(0)
	case "crash":
		time.Sleep(50 * time.Millisecond)
		os.Exit(3)
	case "stubborn":
		signal.Ignore(syscall.SIGTERM)
		time.Sleep(time.Minute)
		os.Exit(0)
	}
}

func spec(kind string, env ...string) func() (Spec, error) {
	return func() (Spec, error) {
		return Spec{Path: os.Args[0], Env: append([]string{childEnv + "=" + kind}, env...), Version: "test"}, nil
	}
}

type running struct {
	t      *testing.T
	o      Options
	cancel context.CancelFunc
	done   chan error
	client *Client
}

// start runs a host on a free loopback port and waits for its first state.
func start(t *testing.T, o Options) *running {
	t.Helper()
	if o.Dir == "" {
		o.Dir = Dir(filepath.Join(t.TempDir(), "host"))
	}
	if o.Addr == "" {
		o.Addr = "127.0.0.1:0"
	}
	if o.Serve == nil {
		o.Serve = spec("run")
	}
	if o.Every == 0 {
		o.Every = 50 * time.Millisecond
	}
	if o.Backoff == 0 {
		o.Backoff = 50 * time.Millisecond
	}
	if o.Grace == 0 {
		o.Grace = 2 * time.Second
	}
	ctx, cancel := context.WithCancel(context.Background())
	r := &running{t: t, o: o, cancel: cancel, done: make(chan error, 1)}
	go func() { r.done <- Run(ctx, o) }()
	t.Cleanup(func() { _ = r.stop() })
	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		if c, _, err := o.Dir.Client(time.Now()); err == nil {
			r.client = c
			return r
		}
		select {
		case err := <-r.done:
			t.Fatalf("host ended at once: %v", err)
		case <-time.After(20 * time.Millisecond):
		}
	}
	t.Fatal("host wrote no state")
	return nil
}

func (r *running) stop() error {
	r.cancel()
	select {
	case err, ok := <-r.done:
		if ok {
			close(r.done)
		}
		return err
	case <-time.After(10 * time.Second):
		r.t.Fatal("host did not stop")
	}
	return nil
}

// until polls the host's status until ok says yes.
func (r *running) until(what string, ok func(Status) bool) Status {
	r.t.Helper()
	deadline := time.Now().Add(5 * time.Second)
	var st Status
	for time.Now().Before(deadline) {
		var err error
		st, err = r.client.Status(context.Background())
		if err == nil && ok(st) {
			return st
		}
		time.Sleep(20 * time.Millisecond)
	}
	r.t.Fatalf("never %s: %+v", what, st)
	return st
}

func childOf(st Status, name, root string) Child {
	for _, c := range st.Children {
		if c.Name == name && c.Root == root {
			return c
		}
	}
	return Child{}
}

func isRunning(name, root string) func(Status) bool {
	return func(st Status) bool { return childOf(st, name, root).State == "running" }
}

func TestServeRunsWithTheWayBackAndGoesWithTheHost(t *testing.T) {
	envFile := filepath.Join(t.TempDir(), "env")
	r := start(t, Options{Serve: spec("run", "FLAI_HOST_TEST_ENVFILE="+envFile), Version: "9.9.9"})
	st := r.until("serve running", isRunning(Serve, ""))
	pid := childOf(st, Serve, "").PID
	if st.Version != "9.9.9" || childOf(st, Serve, "").Version != "test" {
		t.Errorf("versions: host %q, serve %q", st.Version, childOf(st, Serve, "").Version)
	}
	deadline := time.Now().Add(5 * time.Second)
	var env []byte
	for time.Now().Before(deadline) {
		if env, _ = os.ReadFile(envFile); len(env) > 0 {
			break
		}
		time.Sleep(20 * time.Millisecond)
	}
	if want := r.client.URL + "\n" + strings.TrimSpace(r.client.Token); string(env) != want {
		t.Errorf("the child was given %q, want %q", env, want)
	}
	if err := r.stop(); err != nil {
		t.Fatal(err)
	}
	if alive(pid) {
		t.Errorf("serve (pid %d) outlived the host", pid)
	}
	if _, ok := r.o.Dir.ReadStatus(time.Now()); ok {
		t.Error("the state of a stopped host still reads as running")
	}
	if _, err := os.Stat(r.o.Dir.TokenFile()); !errors.Is(err, os.ErrNotExist) {
		t.Errorf("the token outlived the host: %v", err)
	}
}

func TestAChildThatEndsIsStartedAgain(t *testing.T) {
	r := start(t, Options{Serve: spec("crash")})
	st := r.until("serve restarted twice", func(st Status) bool { return childOf(st, Serve, "").Restarts >= 2 })
	if c := childOf(st, Serve, ""); !strings.Contains(c.LastExit, "exit status 3") {
		t.Errorf("last exit %q", c.LastExit)
	}
}

func TestOneHostHoldsTheAddress(t *testing.T) {
	r := start(t, Options{Version: "1.0.0", Config: "/a/config.json"})
	st, _ := r.o.Dir.ReadStatus(time.Now())
	err := Run(context.Background(), Options{Dir: Dir(t.TempDir()), Addr: st.Addr, Serve: spec("run")})
	var already *AlreadyRunningError
	if !errors.As(err, &already) {
		t.Fatalf("a second host on %s: %v", st.Addr, err)
	}
	if already.Health.PID != os.Getpid() || already.Health.Config != "/a/config.json" || !strings.Contains(err.Error(), "one host serves the machine") {
		t.Errorf("said %v with %+v", err, already.Health)
	}
}

func TestTheAPIWantsTheTokenAndNoOrigin(t *testing.T) {
	r := start(t, Options{})
	for _, tc := range []struct {
		name, auth, origin string
		want               int
	}{
		{"no token", "", "", http.StatusUnauthorized},
		{"wrong token", "Bearer nope", "", http.StatusUnauthorized},
		{"a web page", "Bearer " + r.client.Token, "http://localhost:4242", http.StatusForbidden},
		{"the token", "Bearer " + r.client.Token, "", http.StatusOK},
	} {
		req, _ := http.NewRequest(http.MethodGet, r.client.URL+"/status", nil)
		if tc.auth != "" {
			req.Header.Set("Authorization", tc.auth)
		}
		if tc.origin != "" {
			req.Header.Set("Origin", tc.origin)
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		_ = resp.Body.Close()
		if resp.StatusCode != tc.want {
			t.Errorf("%s: %d, want %d", tc.name, resp.StatusCode, tc.want)
		}
	}
	if h, ok := Probe(strings.TrimPrefix(r.client.URL, "http://")); !ok || h.PID != os.Getpid() {
		t.Errorf("health without the token: %+v %v", h, ok)
	}
}

func TestTheHostKeepsTheMCPServersServeAsksFor(t *testing.T) {
	r := start(t, Options{MCP: func(root string) (Spec, error) { return spec("run")() }})
	ctx := context.Background()
	if _, err := r.client.KeepMCP(ctx, []string{"/a", "/b"}); err != nil {
		t.Fatal(err)
	}
	st := r.until("both running", func(st Status) bool { return isRunning(MCP, "/a")(st) && isRunning(MCP, "/b")(st) })
	b := childOf(st, MCP, "/b").PID
	if _, err := r.client.KeepMCP(ctx, []string{"/a"}); err != nil {
		t.Fatal(err)
	}
	st = r.until("b released", func(st Status) bool { return childOf(st, MCP, "/b").Name == "" })
	if alive(b) {
		t.Errorf("the released server (pid %d) still runs", b)
	}
	a := childOf(st, MCP, "/a").PID
	if _, err := r.client.Act(ctx, MCP, "stop"); err != nil {
		t.Fatal(err)
	}
	r.until("a stopped", func(st Status) bool { return childOf(st, MCP, "/a").State == "stopped" })
	if _, err := r.client.KeepMCP(ctx, []string{"/a", "/c"}); err != nil {
		t.Fatal(err)
	}
	st = r.until("c kept stopped", func(st Status) bool { return childOf(st, MCP, "/c").State == "stopped" })
	if alive(a) {
		t.Errorf("a stopped server (pid %d) still runs", a)
	}
	if _, err := r.client.Act(ctx, MCP, "start"); err != nil {
		t.Fatal(err)
	}
	r.until("a and c running", func(st Status) bool { return isRunning(MCP, "/a")(st) && isRunning(MCP, "/c")(st) })
	if !isRunning(Serve, "")(st) {
		t.Errorf("serve was touched: %+v", childOf(st, Serve, ""))
	}
}

func TestServeIsStoppedStartedAndRestarted(t *testing.T) {
	r := start(t, Options{})
	ctx := context.Background()
	first := childOf(r.until("serve", isRunning(Serve, "")), Serve, "").PID
	if _, err := r.client.Act(ctx, Serve, "restart"); err != nil {
		t.Fatal(err)
	}
	second := childOf(r.until("a new serve", func(st Status) bool {
		c := childOf(st, Serve, "")
		return c.State == "running" && c.PID != first
	}), Serve, "").PID
	if alive(first) {
		t.Errorf("the serve before the restart (pid %d) still runs", first)
	}
	st, err := r.client.Act(ctx, Serve, "stop")
	if err != nil {
		t.Fatal(err)
	}
	if c := childOf(st, Serve, ""); c.State != "stopped" || alive(second) {
		t.Errorf("after stop: %+v, pid %d alive %v", c, second, alive(second))
	}
	if _, err := r.client.Act(ctx, Serve, "start"); err != nil {
		t.Fatal(err)
	}
	r.until("serve again", isRunning(Serve, ""))
	if _, err := r.client.Act(ctx, "dashboard", "stop"); err == nil || !strings.Contains(err.Error(), "not a process") {
		t.Errorf("an unknown process: %v", err)
	}
}

func TestAProcessSomeoneElseStartedIsLeftAlone(t *testing.T) {
	r := start(t, Options{Serve: func() (Spec, error) { return Spec{External: 4242}, nil }})
	st := r.until("serve external", func(st Status) bool { return childOf(st, Serve, "").State == "external" })
	if c := childOf(st, Serve, ""); c.PID != 4242 {
		t.Errorf("external %+v", c)
	}
}

func TestAChildThatWillNotStopIsKilled(t *testing.T) {
	r := start(t, Options{Serve: spec("stubborn"), Grace: 200 * time.Millisecond})
	pid := childOf(r.until("serve", isRunning(Serve, "")), Serve, "").PID
	time.Sleep(100 * time.Millisecond) // it has ignored SIGTERM by now
	if err := r.stop(); err != nil {
		t.Fatal(err)
	}
	if alive(pid) {
		t.Errorf("pid %d survived the host", pid)
	}
}

// S-0108: a child asked to stop is stopping, with its PID, for as long as
// its process runs, and stopped only once it has gone. It said stopped at
// once, and a slow runner caught the process still alive (CI after S-0106).
func TestAChildIsStoppedOnlyOnceItsProcessHasEnded(t *testing.T) {
	r := start(t, Options{MCP: func(string) (Spec, error) { return spec("stubborn")() }, Grace: time.Second})
	ctx := context.Background()
	if _, err := r.client.KeepMCP(ctx, []string{"/a"}); err != nil {
		t.Fatal(err)
	}
	pid := childOf(r.until("a running", isRunning(MCP, "/a")), MCP, "/a").PID
	time.Sleep(100 * time.Millisecond) // it has ignored SIGTERM by now
	if _, err := r.client.Act(ctx, MCP, "stop"); err != nil {
		t.Fatal(err)
	}
	st := r.until("a stopping", func(st Status) bool { return childOf(st, MCP, "/a").State == "stopping" })
	if c := childOf(st, MCP, "/a"); c.PID != pid || !alive(pid) {
		t.Errorf("stopping names the process that is still running: %+v, alive %v", c, alive(pid))
	}
	r.until("a stopped", func(st Status) bool { return childOf(st, MCP, "/a").State == "stopped" })
	if alive(pid) {
		t.Errorf("stopped, and pid %d still runs", pid)
	}
}

func TestAnUpgradeThatInstalledRestartsTheHost(t *testing.T) {
	installed := false
	r := start(t, Options{Upgrade: func(context.Context) (any, bool, error) {
		return map[string]string{"installed": "2.0.0"}, !installed, nil
	}})
	pid := childOf(r.until("serve", isRunning(Serve, "")), Serve, "").PID
	out, err := r.client.Upgrade(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(out), `"restarting":true`) {
		t.Errorf("answered %s", out)
	}
	select {
	case err := <-r.done:
		close(r.done)
		if !errors.Is(err, ErrRestart) {
			t.Fatalf("Run returned %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("the host did not end to restart")
	}
	if alive(pid) {
		t.Errorf("serve (pid %d) outlived the host's restart", pid)
	}
}
