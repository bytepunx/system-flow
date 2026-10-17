package serve

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/components"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
)

// TestHelperComponentsVite is the fake vite the supervisor's tests start:
// the test binary run again by the script fakeSubProject writes, with
// FLAI_SERVE_FAKE_VITE naming what it does. Run as a test, it does nothing.
func TestHelperComponentsVite(t *testing.T) {
	mode := os.Getenv("FLAI_SERVE_FAKE_VITE")
	if mode == "" {
		return
	}
	args := os.Args
	if i := slices.Index(args, "--"); i >= 0 {
		args = args[i+1:]
	}
	if log := os.Getenv("FLAI_SERVE_FAKE_VITE_LOG"); log != "" {
		f, err := os.OpenFile(log, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o600)
		if err == nil {
			_, _ = fmt.Fprintln(f, os.Getpid())
			_ = f.Close()
		}
	}
	port := ""
	if i := slices.Index(args, "--port"); i >= 0 && i+1 < len(args) {
		port = args[i+1]
	}
	if mode == "crash-later" {
		go func() {
			time.Sleep(300 * time.Millisecond)
			fmt.Fprintln(os.Stderr, "crashed")
			os.Exit(2)
		}()
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/__flai/components/", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = io.WriteString(w, "harness")
	})
	err := http.ListenAndServe("127.0.0.1:"+port, mux)
	fmt.Fprintln(os.Stderr, "cannot serve:", err)
	os.Exit(1)
}

// fakeSubProject is a project root holding SvelteKit sub-projects of the
// names given, each with a node_modules/.bin/vite that runs the test binary
// as TestHelperComponentsVite in mode, and the file each fake vite writes
// its pid to, a line each.
func fakeSubProject(t *testing.T, mode string, names ...string) (string, []manifest.Project, string) {
	t.Helper()
	if runtime.GOOS == "windows" {
		t.Skip("the fake vite is a POSIX shell script")
	}
	exe, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	root := t.TempDir()
	var subs []manifest.Project
	for _, name := range names {
		bin := filepath.Join(root, name, "node_modules", ".bin")
		if err := os.MkdirAll(bin, 0o755); err != nil {
			t.Fatal(err)
		}
		script := fmt.Sprintf("#!/bin/sh\nexec '%s' '-test.run=^TestHelperComponentsVite$' -- \"$@\"\n", exe)
		if err := os.WriteFile(filepath.Join(bin, "vite"), []byte(script), 0o755); err != nil {
			t.Fatal(err)
		}
		subs = append(subs, manifest.Project{Name: name, Path: name, Kind: manifest.ProjectKindSvelteKit})
	}
	log := filepath.Join(t.TempDir(), "pids")
	t.Setenv("FLAI_SERVE_FAKE_VITE", mode)
	t.Setenv("FLAI_SERVE_FAKE_VITE_LOG", log)
	return root, subs, log
}

// clock is a time a test moves on by hand.
type clock struct {
	mu sync.Mutex
	at time.Time
}

func (c *clock) now() time.Time {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.at
}

func (c *clock) add(d time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.at = c.at.Add(d)
}

// supervisor is a Components on clk, idle after idle, closed when the test
// ends.
func supervisor(t *testing.T, clk *clock, idle time.Duration) *Components {
	t.Helper()
	c := NewComponents(ComponentsOptions{Idle: func() time.Duration { return idle }, Now: clk.now, Timeout: 10 * time.Second})
	t.Cleanup(func() { _ = c.Close() })
	return c
}

// pidsIn are the pids the fake vites wrote to log.
func pidsIn(t *testing.T, log string) []string {
	t.Helper()
	data, err := os.ReadFile(log)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		t.Fatal(err)
	}
	return strings.Fields(string(data))
}

// waitGone fails the test unless pid is gone within a few seconds.
func waitGone(t *testing.T, pid int) {
	t.Helper()
	for range 100 {
		if !Alive(pid) {
			return
		}
		time.Sleep(50 * time.Millisecond)
	}
	t.Errorf("pid %d still runs", pid)
}

func TestComponentsStartsOneDevServerAndAnswersItAgain(t *testing.T) {
	root, subs, log := fakeSubProject(t, "serve", "ui")
	clk := &clock{at: time.Date(2026, 10, 16, 8, 0, 0, 0, time.UTC)}
	c := supervisor(t, clk, time.Hour)
	first, err := c.Start(context.Background(), root, "Shipyard", subs[0])
	if err != nil {
		t.Fatal(err)
	}
	if !first.Running || first.AlreadyRunning || first.Port == 0 || first.PID == 0 || first.Project != "ui" ||
		first.URL != "http://127.0.0.1:"+strconv.Itoa(first.Port)+"/__flai/components/" {
		t.Errorf("the first start: %+v", first)
	}
	if _, err := time.Parse(time.RFC3339, first.Since); err != nil {
		t.Errorf("since is RFC 3339: %q", first.Since)
	}
	if _, err := os.Stat(filepath.Join(root, ".flai-cache", "components", "Shipyard", "ui", components.HarnessConfig)); err != nil {
		t.Errorf("the harness is written under the project's name: %v", err)
	}
	st := c.Status(root, "ui")
	if !st.Running || st.URL != first.URL || st.Port != first.Port || st.PID != first.PID || st.Since != first.Since || st.LastError != "" {
		t.Errorf("the status of the running server: %+v", st)
	}
	again, err := c.Start(context.Background(), root, "Shipyard", subs[0])
	if err != nil || !again.AlreadyRunning || again.PID != first.PID {
		t.Errorf("a second start answers the running server: %+v %v", again, err)
	}
	if got := pidsIn(t, log); len(got) != 1 {
		t.Errorf("one vite ran: %v", got)
	}
	stopped, err := c.Stop(root, "ui")
	if err != nil || !stopped.Stopped || stopped.PID != first.PID {
		t.Fatalf("stop: %+v %v", stopped, err)
	}
	waitGone(t, first.PID)
	if st := c.Status(root, "ui"); st.Running || st.LastError != "" {
		t.Errorf("a stopped server is not running, and was no error: %+v", st)
	}
	if stopped, err := c.Stop(root, "ui"); err != nil || stopped.Stopped {
		t.Errorf("a stop with none running: %+v %v", stopped, err)
	}
}

func TestTwoStartsAtOnceStartOneDevServer(t *testing.T) {
	root, subs, log := fakeSubProject(t, "serve", "ui")
	c := supervisor(t, &clock{at: time.Now()}, time.Hour)
	var wg sync.WaitGroup
	pids := make([]int, 4)
	errs := make([]error, 4)
	for i := range pids {
		wg.Go(func() {
			s, err := c.Start(context.Background(), root, "Shipyard", subs[0])
			pids[i], errs[i] = s.PID, err
		})
	}
	wg.Wait()
	for i := range pids {
		if errs[i] != nil || pids[i] == 0 || pids[i] != pids[0] {
			t.Errorf("start %d: pid %d, %v; want pid %d", i, pids[i], errs[i], pids[0])
		}
	}
	if got := pidsIn(t, log); len(got) != 1 {
		t.Errorf("one vite ran: %v", got)
	}
}

func TestAnIdleDevServerIsStoppedAndAStatusKeepsItAlive(t *testing.T) {
	root, subs, _ := fakeSubProject(t, "serve", "ui")
	clk := &clock{at: time.Date(2026, 10, 16, 8, 0, 0, 0, time.UTC)}
	c := supervisor(t, clk, 10*time.Minute)
	s, err := c.Start(context.Background(), root, "Shipyard", subs[0])
	if err != nil {
		t.Fatal(err)
	}
	clk.add(9 * time.Minute)
	c.stopIdle()
	if !c.Status(root, "ui").Running { // and the status marks it used
		t.Fatal("a server used within components.idle runs on")
	}
	clk.add(9 * time.Minute) // 18 minutes since the start, 9 since the status
	c.stopIdle()
	if !c.Status(root, "ui").Running {
		t.Fatal("a status keeps the server alive")
	}
	clk.add(10 * time.Minute)
	c.stopIdle()
	waitGone(t, s.PID)
	if st := c.Status(root, "ui"); st.Running {
		t.Errorf("a server without a status for components.idle is stopped: %+v", st)
	}
	if _, err := c.Start(context.Background(), root, "Shipyard", subs[0]); err != nil {
		t.Errorf("an idle server that was stopped starts again: %v", err)
	}
}

func TestRunStopsIdleServersAndEveryServerWhenFlaiServeStops(t *testing.T) {
	root, subs, _ := fakeSubProject(t, "serve", "ui", "kit")
	other, otherSubs, _ := fakeSubProject(t, "serve", "ui")
	clk := &clock{at: time.Date(2026, 10, 16, 8, 0, 0, 0, time.UTC)}
	c := NewComponents(ComponentsOptions{Idle: func() time.Duration { return time.Hour }, Every: 20 * time.Millisecond, Now: clk.now, Timeout: 10 * time.Second})
	ctx, stop := context.WithCancel(context.Background())
	ran := make(chan error, 1)
	go func() { ran <- c.Run(ctx) }()
	var pids []int
	for _, start := range []struct {
		root string
		sub  manifest.Project
	}{{root, subs[0]}, {root, subs[1]}, {other, otherSubs[0]}} {
		s, err := c.Start(context.Background(), start.root, "Shipyard", start.sub)
		if err != nil {
			t.Fatal(err)
		}
		pids = append(pids, s.PID)
	}
	if pids[0] == pids[2] {
		t.Errorf("one dev server per project and sub-project: %v", pids)
	}
	clk.add(2 * time.Hour)
	waitGone(t, pids[0]) // Run's look stops what is idle
	s, err := c.Start(context.Background(), root, "Shipyard", subs[0])
	if err != nil {
		t.Fatal(err)
	}
	stop()
	select {
	case err := <-ran:
		if err != nil {
			t.Errorf("Run: %v", err)
		}
	case <-time.After(30 * time.Second):
		t.Fatal("Run did not return once its context ended")
	}
	waitGone(t, s.PID)
	if _, err := c.Start(context.Background(), root, "Shipyard", subs[0]); err == nil || !strings.Contains(err.Error(), "flai serve is stopping") {
		t.Errorf("no start once flai serve stops: %v", err)
	}
	if err := c.Close(); err != nil {
		t.Errorf("a second Close: %v", err)
	}
}

func TestADevServerThatExitsOnItsOwnIsForgottenWithItsError(t *testing.T) {
	root, subs, _ := fakeSubProject(t, "crash-later", "ui")
	c := supervisor(t, &clock{at: time.Now()}, time.Hour)
	s, err := c.Start(context.Background(), root, "Shipyard", subs[0])
	if err != nil {
		t.Fatal(err)
	}
	waitGone(t, s.PID)
	var st = c.Status(root, "ui")
	for range 100 {
		if st = c.Status(root, "ui"); !st.Running {
			break
		}
		time.Sleep(50 * time.Millisecond)
	}
	if st.Running || !strings.Contains(st.LastError, "exited") || !strings.Contains(st.LastError, "crashed") {
		t.Errorf("a server that exited is not running, and says why: %+v", st)
	}
	t.Setenv("FLAI_SERVE_FAKE_VITE", "serve")
	again, err := c.Start(context.Background(), root, "Shipyard", subs[0])
	if err != nil || again.PID == s.PID || again.LastError != "" {
		t.Errorf("it starts again, and its error is gone: %+v %v", again, err)
	}
}

func TestAStartThatFailsKeepsWhy(t *testing.T) {
	root, subs, _ := fakeSubProject(t, "serve", "ui")
	if err := os.RemoveAll(filepath.Join(root, "ui", "node_modules")); err != nil {
		t.Fatal(err)
	}
	c := supervisor(t, &clock{at: time.Now()}, time.Hour)
	_, err := c.Start(context.Background(), root, "Shipyard", subs[0])
	if !errors.Is(err, components.ErrViteMissing) {
		t.Fatalf("a sub-project without vite: %v", err)
	}
	if st := c.Status(root, "ui"); st.Running || st.Starting || !strings.Contains(st.LastError, "vite is not installed") {
		t.Errorf("its status says why: %+v", st)
	}
	if c.Idle() != time.Hour || NewComponents(ComponentsOptions{}).Idle() != DefaultComponentsIdle || DefaultComponentsIdle != 30*time.Minute {
		t.Errorf("components.idle is the operator's, else 30 minutes")
	}
}
