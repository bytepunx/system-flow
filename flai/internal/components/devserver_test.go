package components

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"slices"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"testing"
	"time"
)

// TestHelperVite is the fake vite the dev server tests start: the test binary
// run again by the script newFakeVite writes, with FLAI_FAKE_VITE naming
// what it does. Run as a test, it does nothing.
func TestHelperVite(t *testing.T) {
	mode := os.Getenv("FLAI_FAKE_VITE")
	if mode == "" {
		return
	}
	runFakeVite(mode)
}

// runFakeVite behaves as mode says and never returns. Before anything else
// it writes its pid, its arguments after --, and its folder into
// FLAI_FAKE_VITE_DIR, when that is set.
func runFakeVite(mode string) {
	args := os.Args
	if i := slices.Index(args, "--"); i >= 0 {
		args = args[i+1:]
	}
	if dir := os.Getenv("FLAI_FAKE_VITE_DIR"); dir != "" {
		wd, _ := os.Getwd()
		_ = os.WriteFile(filepath.Join(dir, "args"), []byte(strings.Join(args, "\n")), 0o600)
		_ = os.WriteFile(filepath.Join(dir, "cwd"), []byte(wd), 0o600)
		_ = os.WriteFile(filepath.Join(dir, "pid"), []byte(strconv.Itoa(os.Getpid())), 0o600)
	}
	port := ""
	if i := slices.Index(args, "--port"); i >= 0 && i+1 < len(args) {
		port = args[i+1]
	}
	status := http.StatusOK
	switch mode {
	case "sleep", "never-answer":
		fmt.Println("waiting")
		time.Sleep(time.Hour)
		os.Exit(0)
	case "exit-early":
		for i := 1; i <= 25; i++ {
			fmt.Println("line", i)
		}
		fmt.Println()
		fmt.Fprintln(os.Stderr, "boom")
		os.Exit(1)
	case "spawn-child":
		exe, _ := os.Executable()
		child := exec.Command(exe, "-test.run=^TestHelperVite$")
		child.Env = append(os.Environ(), "FLAI_FAKE_VITE=sleep", "FLAI_FAKE_VITE_DIR=")
		if err := child.Start(); err != nil {
			fmt.Fprintln(os.Stderr, "cannot start the child:", err)
			os.Exit(3)
		}
		_ = os.WriteFile(filepath.Join(os.Getenv("FLAI_FAKE_VITE_DIR"), "child"), []byte(strconv.Itoa(child.Process.Pid)), 0o600)
	case "ignore-term":
		signal.Ignore(syscall.SIGTERM)
	case "crash-later":
		go func() {
			time.Sleep(300 * time.Millisecond)
			fmt.Fprintln(os.Stderr, "crashed")
			os.Exit(2)
		}()
	case "not-found":
		status = http.StatusNotFound
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/__flai/components/", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(status)
		_, _ = io.WriteString(w, "harness")
	})
	fmt.Println("VITE ready")
	err := http.ListenAndServe("127.0.0.1:"+port, mux)
	fmt.Fprintln(os.Stderr, "cannot serve:", err)
	os.Exit(1)
}

// newFakeVite is the options for a sub-project whose node_modules/.bin/vite
// runs the test binary as TestHelperVite in mode, and the folder the fake
// writes what it was given into.
func newFakeVite(t *testing.T, mode string) (DevServerOptions, string) {
	t.Helper()
	if runtime.GOOS == "windows" {
		t.Skip("the fake vite is a POSIX shell script")
	}
	exe, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	sub := t.TempDir()
	bin := filepath.Join(sub, "node_modules", ".bin")
	if err := os.MkdirAll(bin, 0o755); err != nil {
		t.Fatal(err)
	}
	script := fmt.Sprintf("#!/bin/sh\nexec '%s' '-test.run=^TestHelperVite$' -- \"$@\"\n", exe)
	if err := os.WriteFile(filepath.Join(bin, "vite"), []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	state := t.TempDir()
	t.Setenv("FLAI_FAKE_VITE", mode)
	t.Setenv("FLAI_FAKE_VITE_DIR", state)
	return DevServerOptions{
		Dir:     sub,
		Config:  filepath.Join(state, "vite.config.mjs"),
		Path:    "/__flai/components/",
		Timeout: 10 * time.Second,
	}, state
}

// readPID is the pid the fake wrote to name in state.
func readPID(t *testing.T, state, name string) int {
	t.Helper()
	b, err := os.ReadFile(filepath.Join(state, name))
	if err != nil {
		t.Fatalf("the fake vite wrote no %s: %v", name, err)
	}
	pid, err := strconv.Atoi(string(b))
	if err != nil {
		t.Fatal(err)
	}
	return pid
}

// waitGone fails the test when pid is still running after a few seconds,
// and then kills it, by its PID.
func waitGone(t *testing.T, pid int) {
	t.Helper()
	deadline := time.Now().Add(5 * time.Second)
	for !gone(pid) {
		if time.Now().After(deadline) {
			if p, err := os.FindProcess(pid); err == nil {
				_ = p.Kill()
			}
			t.Fatalf("process %d is still running", pid)
		}
		time.Sleep(20 * time.Millisecond)
	}
}

// gone reports whether pid has ended: no process answers a signal 0, or it
// is a zombie, ended and waiting only for its parent to read its status.
func gone(pid int) bool {
	if stat, err := os.ReadFile(fmt.Sprintf("/proc/%d/stat", pid)); err == nil {
		if i := bytes.LastIndexByte(stat, ')'); i >= 0 && bytes.HasPrefix(stat[i+1:], []byte(" Z")) {
			return true
		}
	}
	p, err := os.FindProcess(pid)
	if err != nil {
		return true
	}
	defer func() { _ = p.Release() }()
	return p.Signal(syscall.Signal(0)) != nil
}

// syncBuffer is a bytes.Buffer safe to write from vite's output while the
// test reads it.
type syncBuffer struct {
	mu sync.Mutex
	b  bytes.Buffer
}

func (s *syncBuffer) Write(p []byte) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.b.Write(p)
}

func (s *syncBuffer) String() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.b.String()
}

func closed(c <-chan struct{}) bool {
	select {
	case <-c:
		return true
	default:
		return false
	}
}

func TestStartDevServerRunsTheSubProjectsViteAndWaitsForItToAnswer(t *testing.T) {
	opts, state := newFakeVite(t, "serve")
	var out syncBuffer
	opts.Output = &out
	s, err := StartDevServer(context.Background(), opts)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = s.Stop() })
	if s.Port() <= 0 || s.URL() != fmt.Sprintf("http://127.0.0.1:%d/__flai/components/", s.Port()) {
		t.Errorf("port %d, URL %s", s.Port(), s.URL())
	}
	if pid := readPID(t, state, "pid"); s.PID() != pid {
		t.Errorf("PID %d, the fake is %d", s.PID(), pid)
	}
	if s.Started().IsZero() || time.Since(s.Started()) > time.Minute {
		t.Errorf("started %v", s.Started())
	}
	args, _ := os.ReadFile(filepath.Join(state, "args"))
	want := strings.Join([]string{"--config", opts.Config, "--host", "127.0.0.1", "--port", strconv.Itoa(s.Port()), "--strictPort"}, "\n")
	if string(args) != want {
		t.Errorf("vite ran with\n%s\nwant\n%s", args, want)
	}
	cwd, _ := os.ReadFile(filepath.Join(state, "cwd"))
	if got, _ := filepath.EvalSymlinks(string(cwd)); got != mustEval(t, opts.Dir) {
		t.Errorf("vite ran in %s, want %s", cwd, opts.Dir)
	}
	resp, err := http.Get(s.URL())
	if err != nil {
		t.Fatal(err)
	}
	body, _ := io.ReadAll(resp.Body)
	_ = resp.Body.Close()
	if resp.StatusCode != http.StatusOK || string(body) != "harness" {
		t.Errorf("GET %s: %d %q", s.URL(), resp.StatusCode, body)
	}
	if !strings.Contains(out.String(), "VITE ready") {
		t.Errorf("Output got %q, want vite's output", out.String())
	}
	if closed(s.Done()) || s.Err() != nil {
		t.Errorf("a running server is done (%v) or has an error: %v", closed(s.Done()), s.Err())
	}
	if err := s.Stop(); err != nil {
		t.Fatal(err)
	}
	if !closed(s.Done()) || s.Err() != nil {
		t.Errorf("after Stop, done %v and error %v; want done and no error", closed(s.Done()), s.Err())
	}
	waitGone(t, s.PID())
	if err := s.Stop(); err != nil {
		t.Errorf("a second Stop: %v", err)
	}
}

func mustEval(t *testing.T, p string) string {
	t.Helper()
	got, err := filepath.EvalSymlinks(p)
	if err != nil {
		t.Fatal(err)
	}
	return got
}

func TestStopEndsWhatTheServerStarted(t *testing.T) {
	opts, state := newFakeVite(t, "spawn-child")
	s, err := StartDevServer(context.Background(), opts)
	if err != nil {
		t.Fatal(err)
	}
	child := readPID(t, state, "child")
	if gone(child) {
		t.Fatalf("the fake's child %d is not running before Stop", child)
	}
	if err := s.Stop(); err != nil {
		t.Fatal(err)
	}
	waitGone(t, s.PID())
	waitGone(t, child)
}

func TestStopKillsAServerThatIgnoresTerm(t *testing.T) {
	opts, _ := newFakeVite(t, "ignore-term")
	defer func(g time.Duration) { stopGrace = g }(stopGrace)
	stopGrace = 300 * time.Millisecond
	s, err := StartDevServer(context.Background(), opts)
	if err != nil {
		t.Fatal(err)
	}
	began := time.Now()
	if err := s.Stop(); err != nil {
		t.Fatal(err)
	}
	if took := time.Since(began); took < stopGrace {
		t.Errorf("Stop took %s, under the grace of %s: the fake did not ignore TERM", took, stopGrace)
	}
	waitGone(t, s.PID())
}

func TestStopFromManyGoroutinesStopsOnce(t *testing.T) {
	opts, _ := newFakeVite(t, "serve")
	s, err := StartDevServer(context.Background(), opts)
	if err != nil {
		t.Fatal(err)
	}
	var wg sync.WaitGroup
	for range 4 {
		wg.Go(func() {
			if err := s.Stop(); err != nil {
				t.Error(err)
			}
			if !closed(s.Done()) {
				t.Error("Stop returned before the server was done")
			}
		})
	}
	wg.Wait()
	waitGone(t, s.PID())
}

func TestStartDevServerReturnsTheLastLinesOfAViteThatExitsEarly(t *testing.T) {
	opts, state := newFakeVite(t, "exit-early")
	_, err := StartDevServer(context.Background(), opts)
	if err == nil {
		t.Fatal("a vite that exited started")
	}
	msg := err.Error()
	for _, w := range []string{"exited before it answered at http://127.0.0.1:", "exit status 1", "line 7\n", "line 25\nboom"} {
		if !strings.Contains(msg, w) {
			t.Errorf("error lacks %q:\n%s", w, msg)
		}
	}
	if strings.Contains(msg, "line 6") {
		t.Errorf("error keeps more than the last %d lines:\n%s", tailLines, msg)
	}
	waitGone(t, readPID(t, state, "pid"))
}

func TestStartDevServerStopsAViteThatDoesNotAnswerInTime(t *testing.T) {
	opts, state := newFakeVite(t, "never-answer")
	opts.Timeout = 500 * time.Millisecond
	_, err := StartDevServer(context.Background(), opts)
	if err == nil {
		t.Fatal("a vite that never answered started")
	}
	for _, w := range []string{"did not answer at http://127.0.0.1:", "within 500ms", "flai stopped it", "waiting"} {
		if !strings.Contains(err.Error(), w) {
			t.Errorf("error lacks %q:\n%v", w, err)
		}
	}
	waitGone(t, readPID(t, state, "pid"))
}

func TestStartDevServerSaysWhatAViteThatDoesNotAnswer200Answered(t *testing.T) {
	opts, state := newFakeVite(t, "not-found")
	opts.Timeout = 600 * time.Millisecond
	_, err := StartDevServer(context.Background(), opts)
	if err == nil || !strings.Contains(err.Error(), "it answered 404, and the harness answers 200") {
		t.Errorf("error %v, want the 404 it answered", err)
	}
	waitGone(t, readPID(t, state, "pid"))
}

func TestStartDevServerStopsViteWhenItsContextEnds(t *testing.T) {
	opts, state := newFakeVite(t, "never-answer")
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()
	_, err := StartDevServer(ctx, opts)
	if !errors.Is(err, context.DeadlineExceeded) || !strings.Contains(err.Error(), "stopped the Vite dev server") {
		t.Errorf("error %v, want the context's", err)
	}
	waitGone(t, readPID(t, state, "pid"))
}

func TestErrSaysWhyAServerExitedOnItsOwn(t *testing.T) {
	opts, _ := newFakeVite(t, "crash-later")
	s, err := StartDevServer(context.Background(), opts)
	if err != nil {
		t.Fatal(err)
	}
	select {
	case <-s.Done():
	case <-time.After(5 * time.Second):
		_ = s.Stop()
		t.Fatal("the server did not exit on its own")
	}
	for _, check := range []string{"before Stop", "after Stop"} {
		err := s.Err()
		if err == nil || !strings.Contains(err.Error(), "exit status 2") || !strings.Contains(err.Error(), "crashed") {
			t.Errorf("%s, Err %v; want its exit status and last lines", check, err)
		}
		if err := s.Stop(); err != nil {
			t.Errorf("Stop of an exited server: %v", err)
		}
	}
}

func TestStartDevServerRefusesASubProjectWithoutVite(t *testing.T) {
	dir := t.TempDir()
	_, err := StartDevServer(context.Background(), DevServerOptions{Dir: dir})
	if !errors.Is(err, ErrViteMissing) || !strings.Contains(err.Error(), "npm install or pnpm install in "+dir) {
		t.Errorf("error %v, want ErrViteMissing with the fix", err)
	}
}

func TestStartDevServerRefusesARelativeFolder(t *testing.T) {
	if _, err := StartDevServer(context.Background(), DevServerOptions{Dir: "web"}); err == nil || !strings.Contains(err.Error(), "absolute path") {
		t.Errorf("error %v, want one asking for an absolute path", err)
	}
}

func TestTailKeepsTheLastLinesThatAreNotBlank(t *testing.T) {
	var copied bytes.Buffer
	tl := &tail{out: &copied}
	var in strings.Builder
	for i := 1; i <= tailLines+2; i++ {
		fmt.Fprintf(&in, "line %d\r\n\n", i)
	}
	in.WriteString(strings.Repeat("x", maxLine+10))
	in.WriteString("\nhalf")
	for _, chunk := range strings.SplitAfter(in.String(), "e") {
		if _, err := tl.Write([]byte(chunk)); err != nil {
			t.Fatal(err)
		}
	}
	if copied.String() != in.String() {
		t.Errorf("out got %d bytes, want all %d written", copied.Len(), in.Len())
	}
	got := strings.Split(string(tl.last()), "\n")
	want := []string{}
	for i := 4; i <= tailLines+2; i++ {
		want = append(want, fmt.Sprintf("line %d", i))
	}
	want = append(want, strings.Repeat("x", maxLine), "half")
	if !slices.Equal(got, want) {
		t.Errorf("kept\n%q\nwant\n%q", got, want)
	}
}
