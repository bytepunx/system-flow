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
	"github.com/bytepunx/system-flow/flai/internal/hostapi"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
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
	first, err := c.Start(context.Background(), root, "Shipyard", subs[0], "", "")
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
	st := c.Status(root, "ui", "")
	if !st.Running || st.URL != first.URL || st.Port != first.Port || st.PID != first.PID || st.Since != first.Since || st.LastError != "" {
		t.Errorf("the status of the running server: %+v", st)
	}
	again, err := c.Start(context.Background(), root, "Shipyard", subs[0], "", "")
	if err != nil || !again.AlreadyRunning || again.PID != first.PID {
		t.Errorf("a second start answers the running server: %+v %v", again, err)
	}
	if got := pidsIn(t, log); len(got) != 1 {
		t.Errorf("one vite ran: %v", got)
	}
	stopped, err := c.Stop(root, "ui", "")
	if err != nil || !stopped.Stopped || stopped.PID != first.PID {
		t.Fatalf("stop: %+v %v", stopped, err)
	}
	waitGone(t, first.PID)
	if st := c.Status(root, "ui", ""); st.Running || st.LastError != "" {
		t.Errorf("a stopped server is not running, and was no error: %+v", st)
	}
	if stopped, err := c.Stop(root, "ui", ""); err != nil || stopped.Stopped {
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
			s, err := c.Start(context.Background(), root, "Shipyard", subs[0], "", "")
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
	s, err := c.Start(context.Background(), root, "Shipyard", subs[0], "", "")
	if err != nil {
		t.Fatal(err)
	}
	clk.add(9 * time.Minute)
	c.stopIdle()
	if !c.Status(root, "ui", "").Running { // and the status marks it used
		t.Fatal("a server used within components.idle runs on")
	}
	clk.add(9 * time.Minute) // 18 minutes since the start, 9 since the status
	c.stopIdle()
	if !c.Status(root, "ui", "").Running {
		t.Fatal("a status keeps the server alive")
	}
	clk.add(10 * time.Minute)
	c.stopIdle()
	waitGone(t, s.PID)
	if st := c.Status(root, "ui", ""); st.Running {
		t.Errorf("a server without a status for components.idle is stopped: %+v", st)
	}
	if _, err := c.Start(context.Background(), root, "Shipyard", subs[0], "", ""); err != nil {
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
		s, err := c.Start(context.Background(), start.root, "Shipyard", start.sub, "", "")
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
	s, err := c.Start(context.Background(), root, "Shipyard", subs[0], "", "")
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
	if _, err := c.Start(context.Background(), root, "Shipyard", subs[0], "", ""); err == nil || !strings.Contains(err.Error(), "flai serve is stopping") {
		t.Errorf("no start once flai serve stops: %v", err)
	}
	if err := c.Close(); err != nil {
		t.Errorf("a second Close: %v", err)
	}
}

func TestADevServerThatExitsOnItsOwnIsForgottenWithItsError(t *testing.T) {
	root, subs, _ := fakeSubProject(t, "crash-later", "ui")
	c := supervisor(t, &clock{at: time.Now()}, time.Hour)
	s, err := c.Start(context.Background(), root, "Shipyard", subs[0], "", "")
	if err != nil {
		t.Fatal(err)
	}
	waitGone(t, s.PID)
	var st = c.Status(root, "ui", "")
	for range 100 {
		if st = c.Status(root, "ui", ""); !st.Running {
			break
		}
		time.Sleep(50 * time.Millisecond)
	}
	if st.Running || !strings.Contains(st.LastError, "exited") || !strings.Contains(st.LastError, "crashed") {
		t.Errorf("a server that exited is not running, and says why: %+v", st)
	}
	t.Setenv("FLAI_SERVE_FAKE_VITE", "serve")
	again, err := c.Start(context.Background(), root, "Shipyard", subs[0], "", "")
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
	_, err := c.Start(context.Background(), root, "Shipyard", subs[0], "", "")
	if !errors.Is(err, components.ErrViteMissing) {
		t.Fatalf("a sub-project without vite: %v", err)
	}
	if st := c.Status(root, "ui", ""); st.Running || st.Starting || !strings.Contains(st.LastError, "vite is not installed") {
		t.Errorf("its status says why: %+v", st)
	}
	if c.Idle() != time.Hour || NewComponents(ComponentsOptions{}).Idle() != DefaultComponentsIdle || DefaultComponentsIdle != 30*time.Minute {
		t.Errorf("components.idle is the operator's, else 30 minutes")
	}
}

// fakeInstall is a package manager that installs a sub-project's
// dependencies by copying vite, the script at vite, into the folder it runs
// in, and remembers each folder it ran in.
type fakeInstall struct {
	vite string
	mu   sync.Mutex
	dirs []string
}

func (f *fakeInstall) LookPath(name string) (string, error) { return "/usr/bin/" + name, nil }

func (f *fakeInstall) Run(_ context.Context, dir, _, _ string, _ ...string) ([]byte, []byte, error) {
	f.mu.Lock()
	f.dirs = append(f.dirs, dir)
	f.mu.Unlock()
	script, err := os.ReadFile(f.vite)
	if err != nil {
		return nil, nil, err
	}
	bin := filepath.Join(dir, "node_modules", ".bin")
	if err := os.MkdirAll(bin, 0o755); err != nil {
		return nil, nil, err
	}
	return []byte("installed\n"), nil, os.WriteFile(filepath.Join(bin, "vite"), script, 0o755)
}

// worktreeOf is a story's worktree beside root: its sub-project ui with a
// pnpm lock file and no dependencies installed.
func worktreeOf(t *testing.T) string {
	t.Helper()
	wt := t.TempDir()
	if err := os.MkdirAll(filepath.Join(wt, "ui"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(wt, "ui", "pnpm-lock.yaml"), []byte("lockfileVersion: '9.0'\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	return wt
}

// S-0442: a story's dev server runs in its worktree, after the install its
// missing dependencies need, on a port of its own, beside main's; each is
// answered, stopped, and keyed by its own checkout.
func TestAStoryDevServerRunsInItsWorktreeBesideMain(t *testing.T) {
	root, subs, _ := fakeSubProject(t, "serve", "ui")
	wt := worktreeOf(t)
	install := &fakeInstall{vite: filepath.Join(root, "ui", "node_modules", ".bin", "vite")}
	clk := &clock{at: time.Date(2026, 10, 16, 8, 0, 0, 0, time.UTC)}
	c := NewComponents(ComponentsOptions{Idle: func() time.Duration { return time.Hour }, Now: clk.now, Timeout: 10 * time.Second, Install: install})
	t.Cleanup(func() { _ = c.Close() })
	main, err := c.Start(context.Background(), root, "Shipyard", subs[0], "", "")
	if err != nil {
		t.Fatal(err)
	}
	story, err := c.Start(context.Background(), root, "Shipyard", subs[0], "S-0001", wt)
	if err != nil {
		t.Fatal(err)
	}
	if !story.Running || story.Story != "S-0001" || story.AlreadyRunning || story.PID == main.PID || story.Port == main.Port || main.Story != "" {
		t.Errorf("the story's server is its own, on its own port: main %+v, story %+v", main, story)
	}
	if len(install.dirs) != 1 || install.dirs[0] != filepath.Join(wt, "ui") {
		t.Errorf("the worktree's dependencies are installed in it by its lock file: %v", install.dirs)
	}
	if _, err := os.Stat(filepath.Join(wt, ".flai-cache", "components", "Shipyard", "ui", components.HarnessConfig)); err != nil {
		t.Errorf("the story's harness is written in its worktree: %v", err)
	}
	if st := c.Status(root, "ui", "S-0001"); !st.Running || st.PID != story.PID || st.Story != "S-0001" {
		t.Errorf("the story's status: %+v", st)
	}
	if st := c.Status(root, "ui", "S-0002"); st.Running || st.Story != "S-0002" {
		t.Errorf("another story's status: %+v", st)
	}
	again, err := c.Start(context.Background(), root, "Shipyard", subs[0], "S-0001", wt)
	if err != nil || !again.AlreadyRunning || again.PID != story.PID || len(install.dirs) != 1 {
		t.Errorf("a second start of the story's answers it, installing nothing: %+v %v %v", again, err, install.dirs)
	}
	stopped, err := c.Stop(root, "ui", "S-0001")
	if err != nil || !stopped.Stopped || stopped.PID != story.PID || stopped.Story != "S-0001" {
		t.Fatalf("stop the story's: %+v %v", stopped, err)
	}
	waitGone(t, story.PID)
	if st := c.Status(root, "ui", ""); !st.Running || st.PID != main.PID {
		t.Errorf("main's runs on: %+v", st)
	}
}

// S-0442: a story's start whose install fails is the start's error, kept as
// its last error, and starts no vite.
func TestAStoryStartWithoutALockFileKeepsWhy(t *testing.T) {
	root, subs, log := fakeSubProject(t, "serve", "ui")
	wt := t.TempDir()
	if err := os.MkdirAll(filepath.Join(wt, "ui"), 0o755); err != nil {
		t.Fatal(err)
	}
	c := supervisor(t, &clock{at: time.Now()}, time.Hour)
	_, err := c.Start(context.Background(), root, "Shipyard", subs[0], "S-0001", wt)
	if !errors.Is(err, components.ErrNoLockFile) {
		t.Fatalf("a worktree with no lock file: %v", err)
	}
	if st := c.Status(root, "ui", "S-0001"); st.Running || !strings.Contains(st.LastError, "no lock file") {
		t.Errorf("its status says why: %+v", st)
	}
	if got := pidsIn(t, log); len(got) != 0 {
		t.Errorf("no vite ran: %v", got)
	}
}

// S-0442: components.max_servers caps the dev servers running or starting
// per project, main's and every story's together; a start beyond it is
// refused naming them, a start answered by one running is not, and another
// project's do not count. 3 is the default.
func TestComponentsMaxServersCapsAProjectsDevServers(t *testing.T) {
	root, subs, _ := fakeSubProject(t, "serve", "ui")
	other, otherSubs, _ := fakeSubProject(t, "serve", "ui")
	wt1, _, _ := fakeSubProject(t, "serve", "ui")
	wt2, _, _ := fakeSubProject(t, "serve", "ui")
	clk := &clock{at: time.Date(2026, 10, 16, 8, 0, 0, 0, time.UTC)}
	c := NewComponents(ComponentsOptions{Idle: func() time.Duration { return time.Hour }, MaxServers: func() int { return 2 }, Now: clk.now, Timeout: 10 * time.Second})
	t.Cleanup(func() { _ = c.Close() })
	if _, err := c.Start(context.Background(), root, "Shipyard", subs[0], "", ""); err != nil {
		t.Fatal(err)
	}
	clk.add(time.Second)
	if _, err := c.Start(context.Background(), root, "Shipyard", subs[0], "S-0441", wt1); err != nil {
		t.Fatal(err)
	}
	_, err := c.Start(context.Background(), root, "Shipyard", subs[0], "S-0440", wt2)
	want := "cannot start S-0440's ui component dev server: flai runs 2 component dev servers for this project, components.max_servers: main's ui, S-0441's ui; stop one first"
	if !errors.Is(err, hostapi.ErrComponentServersFull) || !strings.HasPrefix(err.Error(), want) {
		t.Errorf("a third start: %v\nwant %s", err, want)
	}
	if st := c.Status(root, "ui", "S-0440"); st.Running || st.Starting || st.LastError != "" {
		t.Errorf("a refused start leaves nothing behind: %+v", st)
	}
	if again, err := c.Start(context.Background(), root, "Shipyard", subs[0], "S-0441", wt1); err != nil || !again.AlreadyRunning {
		t.Errorf("a start answered by the running server is not refused: %+v %v", again, err)
	}
	if _, err := c.Start(context.Background(), other, "Shipyard", otherSubs[0], "", ""); err != nil {
		t.Errorf("another project's servers are its own: %v", err)
	}
	if _, err := c.Stop(root, "ui", ""); err != nil {
		t.Fatal(err)
	}
	if _, err := c.Start(context.Background(), root, "Shipyard", subs[0], "S-0440", wt2); err != nil {
		t.Errorf("once one stops, another starts: %v", err)
	}
	if c.MaxServers() != 2 || NewComponents(ComponentsOptions{}).MaxServers() != 3 || DefaultComponentsMaxServers != 3 ||
		NewComponents(ComponentsOptions{MaxServers: func() int { return 0 }}).MaxServers() != 3 {
		t.Errorf("components.max_servers is the operator's, else 3")
	}
}

// at is a transition to the state to, at t.
func at(to string, t time.Time) workitem.Transition {
	return workitem.Transition{To: to, At: t.UTC().Format(time.RFC3339), By: "olive"}
}

// S-0442: a story's dev server stops when its story leaves review, to any
// state, or is neither in progress nor in review, or is gone; one in
// progress, or moved from in progress to review, keeps it, as does one that
// left review before its server started.
func TestStoryLeft(t *testing.T) {
	start := time.Date(2026, 10, 16, 8, 0, 0, 500_000_000, time.UTC)
	before, after := start.Add(-time.Hour), start.Add(time.Minute)
	story := func(status string, ts ...workitem.Transition) *workitem.Item {
		return &workitem.Item{ID: "S-0001", Type: workitem.Story, Status: status, Transitions: ts}
	}
	open := []workitem.Transition{at(workitem.Ready, before), at(workitem.InProgress, before)}
	for _, c := range []struct {
		name string
		it   *workitem.Item
		left bool
	}{
		{"in progress", story(workitem.InProgress, open...), false},
		{"in progress to review", story(workitem.Review, append(open, at(workitem.Review, after))...), false},
		{"in review since before", story(workitem.Review, append(open, at(workitem.Review, before))...), false},
		{"back to in progress before the start", story(workitem.InProgress, append(open, at(workitem.Review, before), at(workitem.InProgress, before))...), false},
		{"review to done", story(workitem.Done, append(open, at(workitem.Review, before), at(workitem.Done, after))...), true},
		{"review to in progress", story(workitem.InProgress, append(open, at(workitem.Review, before), at(workitem.InProgress, after))...), true},
		{"review to in progress and back", story(workitem.Review, append(open, at(workitem.Review, before), at(workitem.InProgress, after), at(workitem.Review, after))...), true},
		{"review to in progress in the start's second", story(workitem.InProgress, append(open, at(workitem.Review, before), at(workitem.InProgress, start))...), true},
		{"cancelled from in progress", story(workitem.Cancelled, append(open, at(workitem.Cancelled, after))...), true},
		{"ready", story(workitem.Ready, at(workitem.Ready, before)), true},
		{"archived", &workitem.Item{Status: workitem.Done, Archived: true}, true},
		{"gone", nil, true},
	} {
		if got := StoryLeft(c.it, start); got != c.left {
			t.Errorf("%s: StoryLeft = %v, want %v", c.name, got, c.left)
		}
	}
}

// stories is a Story for ComponentsOptions answering from a map the test
// changes, an ID it lacks not found, and err when set.
type stories struct {
	mu  sync.Mutex
	by  map[string]*workitem.Item
	err error
}

func (s *stories) get(_, id string) (*workitem.Item, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.err != nil {
		return nil, s.err
	}
	it, ok := s.by[id]
	if !ok {
		return nil, fmt.Errorf("%s %w", id, workitem.ErrNotFound)
	}
	return it, nil
}

func (s *stories) set(id string, it *workitem.Item) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.by[id] = it
}

// S-0442: Run's look stops a story's dev server once the story leaves
// review, and keeps it while it moves from in progress to review, while its
// story cannot be read, and main's always; it stops one whose story is gone.
func TestAStorysDevServerStopsWhenItsStoryLeavesReview(t *testing.T) {
	root, subs, _ := fakeSubProject(t, "serve", "ui")
	wt1, _, _ := fakeSubProject(t, "serve", "ui")
	wt2, _, _ := fakeSubProject(t, "serve", "ui")
	clk := &clock{at: time.Date(2026, 10, 16, 8, 0, 0, 0, time.UTC)}
	early := clk.now().Add(-time.Hour)
	board := &stories{by: map[string]*workitem.Item{
		"S-0001": {ID: "S-0001", Status: workitem.InProgress, Transitions: []workitem.Transition{at(workitem.Ready, early), at(workitem.InProgress, early)}},
		"S-0002": {ID: "S-0002", Status: workitem.Review, Transitions: []workitem.Transition{at(workitem.Ready, early), at(workitem.InProgress, early), at(workitem.Review, early)}},
	}}
	c := NewComponents(ComponentsOptions{Idle: func() time.Duration { return time.Hour }, StoryEvery: 20 * time.Millisecond, Story: board.get, Now: clk.now, Timeout: 10 * time.Second})
	ctx, stop := context.WithCancel(context.Background())
	ran := make(chan error, 1)
	go func() { ran <- c.Run(ctx) }()
	t.Cleanup(func() {
		stop()
		<-ran
	})
	main, err := c.Start(context.Background(), root, "Shipyard", subs[0], "", "")
	if err != nil {
		t.Fatal(err)
	}
	one, err := c.Start(context.Background(), root, "Shipyard", subs[0], "S-0001", wt1)
	if err != nil {
		t.Fatal(err)
	}
	two, err := c.Start(context.Background(), root, "Shipyard", subs[0], "S-0002", wt2)
	if err != nil {
		t.Fatal(err)
	}
	later := clk.now().Add(time.Minute)
	board.set("S-0001", &workitem.Item{ID: "S-0001", Status: workitem.Review, Transitions: append(board.by["S-0001"].Transitions, at(workitem.Review, later))})
	board.set("S-0002", &workitem.Item{ID: "S-0002", Status: workitem.InProgress, Transitions: append(board.by["S-0002"].Transitions, at(workitem.InProgress, later))})
	waitGone(t, two.PID) // review to in progress
	if st := c.Status(root, "ui", "S-0002"); st.Running || st.LastError != "" {
		t.Errorf("a story that left review has no server, and no error: %+v", st)
	}
	if st := c.Status(root, "ui", "S-0001"); !st.Running || st.PID != one.PID {
		t.Errorf("in progress to review keeps the server: %+v", st)
	}
	board.mu.Lock()
	board.err = errors.New("permission denied")
	board.mu.Unlock()
	time.Sleep(100 * time.Millisecond)
	if !c.Status(root, "ui", "S-0001").Running {
		t.Error("a story that cannot be read keeps its server")
	}
	board.mu.Lock()
	board.err = nil
	delete(board.by, "S-0001") // archived away, say
	board.mu.Unlock()
	waitGone(t, one.PID)
	if st := c.Status(root, "ui", ""); !st.Running || st.PID != main.PID {
		t.Errorf("main's server runs on: %+v", st)
	}
}
