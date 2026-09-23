package serve

import (
	"context"
	"os/exec"
	"sync"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/channel"
	"github.com/bytepunx/system-flow/flai/internal/hostapi"
)

// fakeMCP stands in for cmd's launcher; each "server" is a real process, so
// that liveness, and stopping it, are real.
type fakeMCP struct {
	t       *testing.T
	mu      sync.Mutex
	running map[string]*exec.Cmd // by root
	stopped []int
	foreign map[string]int // a server someone else started, by root
}

func newFakeMCP(t *testing.T) *fakeMCP {
	f := &fakeMCP{t: t, running: map[string]*exec.Cmd{}, foreign: map[string]int{}}
	t.Cleanup(func() {
		f.mu.Lock()
		defer f.mu.Unlock()
		for _, c := range f.running {
			_ = c.Process.Kill()
			_ = c.Wait()
		}
	})
	return f
}

func (f *fakeMCP) Ensure(root string) (int, bool, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if pid, ok := f.foreign[root]; ok {
		return pid, false, nil
	}
	if c, ok := f.running[root]; ok && Alive(c.Process.Pid) {
		return c.Process.Pid, false, nil
	}
	c := exec.Command("sleep", "60")
	if err := c.Start(); err != nil {
		f.t.Fatal(err)
	}
	f.running[root] = c
	return c.Process.Pid, true, nil
}

func (f *fakeMCP) Stop(pid int) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.stopped = append(f.stopped, pid)
	for root, c := range f.running {
		if c.Process.Pid == pid {
			_ = c.Process.Kill()
			_ = c.Wait()
			delete(f.running, root)
		}
	}
}

func (f *fakeMCP) pid(root string) int {
	f.mu.Lock()
	defer f.mu.Unlock()
	if c, ok := f.running[root]; ok {
		return c.Process.Pid
	}
	return 0
}

func (f *fakeMCP) stops() []int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return append([]int(nil), f.stopped...)
}

func runWithMCP(t *testing.T, dir Dir, m MCP) (stop func()) {
	t.Helper()
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() {
		done <- Run(ctx, Options{Dir: dir, Version: "test", Every: 20 * time.Millisecond, WatchEvery: 20 * time.Millisecond,
			MCP: m, MCPEvery: 30 * time.Millisecond,
			NewClient: func(e Entry, key []byte) *channel.Client {
				return &channel.Client{URL: e.URL, Key: key, Project: channel.Project{Key: e.Key, Name: e.Name, Root: e.Root},
					Methods: hostapi.Methods("test", nil), Version: "test", MinBackoff: 10 * time.Millisecond, MaxBackoff: 40 * time.Millisecond}
			}})
	}()
	var once sync.Once
	stop = func() {
		once.Do(func() {
			cancel()
			if err := <-done; err != nil {
				t.Errorf("Run: %v", err)
			}
		})
	}
	t.Cleanup(stop)
	return stop
}

// S-0096: flai serve starts each served project's MCP server, starts it
// again when it stops, and stops it when flai serve stops.
func TestServeKeepsEachProjectsMCPServerAndStopsItWithItself(t *testing.T) {
	dir := DirFor(t.TempDir() + "/config.json")
	alpha, alphaKey := scratchProject(t, "alpha")
	beta, betaKey := scratchProject(t, "beta")
	for _, e := range []Entry{
		{Key: "alpha", Name: "Alpha", Root: alpha, URL: "http://127.0.0.1:1", KeyFile: alphaKey},
		{Key: "beta", Name: "Beta", Root: beta, URL: "http://127.0.0.1:1", KeyFile: betaKey},
	} {
		if err := dir.Register(e); err != nil {
			t.Fatal(err)
		}
	}
	m := newFakeMCP(t)
	stop := runWithMCP(t, dir, m)

	waitFor(t, "a server per project", func() bool { return m.pid(alpha) != 0 && m.pid(beta) != 0 })
	first := m.pid(alpha)

	// one that stops (crashed, or stopped by hand) is started again
	m.mu.Lock()
	_ = m.running[alpha].Process.Kill()
	_ = m.running[alpha].Wait()
	m.mu.Unlock()
	waitFor(t, "alpha's server started again", func() bool { p := m.pid(alpha); return p != 0 && p != first && Alive(p) })

	a, b := m.pid(alpha), m.pid(beta)
	stop()
	got := m.stops()
	both := len(got) == 2 && (got[0] == a && got[1] == b || got[0] == b && got[1] == a)
	if !both {
		t.Errorf("stopped with flai serve: %v, want %d and %d", got, a, b)
	}
	if Alive(a) || Alive(b) {
		t.Error("a server outlived flai serve")
	}
}

// A project that is no longer served has its server stopped with it; the
// other project's keeps running.
func TestADroppedProjectsMCPServerStopsWithIt(t *testing.T) {
	dir := DirFor(t.TempDir() + "/config.json")
	alpha, alphaKey := scratchProject(t, "alpha")
	beta, betaKey := scratchProject(t, "beta")
	_ = dir.Register(Entry{Key: "alpha", Name: "Alpha", Root: alpha, URL: "http://127.0.0.1:1", KeyFile: alphaKey})
	_ = dir.Register(Entry{Key: "beta", Name: "Beta", Root: beta, URL: "http://127.0.0.1:1", KeyFile: betaKey})
	m := newFakeMCP(t)
	runWithMCP(t, dir, m)
	waitFor(t, "a server per project", func() bool { return m.pid(alpha) != 0 && m.pid(beta) != 0 })
	a, b := m.pid(alpha), m.pid(beta)

	_ = dir.Unregister(alpha)
	waitFor(t, "alpha's server stopped", func() bool { return !Alive(a) })
	if got := m.stops(); len(got) != 1 || got[0] != a {
		t.Errorf("stopped: %v, want only %d", got, a)
	}
	if !Alive(b) {
		t.Error("beta's server stopped with alpha")
	}
}

// A server someone else started (flai mcp start) is used and never stopped:
// it is not flai serve's.
func TestAnMCPServerFlaiServeDidNotStartIsLeftAlone(t *testing.T) {
	dir := DirFor(t.TempDir() + "/config.json")
	alpha, alphaKey := scratchProject(t, "alpha")
	_ = dir.Register(Entry{Key: "alpha", Name: "Alpha", Root: alpha, URL: "http://127.0.0.1:1", KeyFile: alphaKey})
	theirs := exec.Command("sleep", "60")
	if err := theirs.Start(); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = theirs.Process.Kill(); _ = theirs.Wait() })
	m := newFakeMCP(t)
	m.foreign[alpha] = theirs.Process.Pid
	stop := runWithMCP(t, dir, m)
	time.Sleep(100 * time.Millisecond) // several looks
	stop()
	if got := m.stops(); len(got) != 0 {
		t.Errorf("stopped a server it did not start: %v", got)
	}
	if !Alive(theirs.Process.Pid) {
		t.Error("their server is gone")
	}
}
