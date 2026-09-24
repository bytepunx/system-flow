package serve

import (
	"context"
	"errors"
	"slices"
	"sync"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/channel"
	"github.com/bytepunx/system-flow/flai/internal/hostapi"
)

// fakeMCP stands in for cmd's host client: it remembers what it was told.
type fakeMCP struct {
	mu   sync.Mutex
	told [][]string
	fail error
}

func (f *fakeMCP) Keep(roots []string) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.fail != nil {
		return f.fail
	}
	f.told = append(f.told, append([]string(nil), roots...))
	return nil
}

func (f *fakeMCP) calls() [][]string {
	f.mu.Lock()
	defer f.mu.Unlock()
	return append([][]string(nil), f.told...)
}

func (f *fakeMCP) last() []string {
	c := f.calls()
	if len(c) == 0 {
		return nil
	}
	return c[len(c)-1]
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

// S-0106: flai serve tells the host every project it serves, at once when
// that changes and again every so often, and a project it no longer serves
// leaves the list. It starts no MCP server itself.
func TestServeTellsTheHostWhichProjectsNeedAnMCPServer(t *testing.T) {
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
	m := &fakeMCP{}
	stop := runWithMCP(t, dir, m)
	both := []string{alpha, beta}
	slices.Sort(both)
	waitFor(t, "both projects told", func() bool { return slices.Equal(m.last(), both) })
	n := len(m.calls())
	waitFor(t, "told again unchanged", func() bool { return len(m.calls()) > n })

	_ = dir.Unregister(alpha)
	waitFor(t, "alpha dropped", func() bool { return slices.Equal(m.last(), []string{beta}) })

	stop()
	if got := m.last(); !slices.Equal(got, []string{beta}) {
		t.Errorf("flai serve stopping told the host %v; the host keeps what it has until told otherwise", got)
	}
}

// A host that cannot be reached is told again at the next look, not after
// the long wait.
func TestAHostThatDidNotHearIsToldAgain(t *testing.T) {
	dir := DirFor(t.TempDir() + "/config.json")
	alpha, alphaKey := scratchProject(t, "alpha")
	_ = dir.Register(Entry{Key: "alpha", Name: "Alpha", Root: alpha, URL: "http://127.0.0.1:1", KeyFile: alphaKey})
	m := &fakeMCP{fail: errors.New("flai host restarting")}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	done := make(chan error, 1)
	go func() {
		done <- Run(ctx, Options{Dir: dir, Version: "test", Every: 20 * time.Millisecond, MCP: m, MCPEvery: time.Hour,
			NewClient: func(e Entry, key []byte) *channel.Client {
				return &channel.Client{URL: e.URL, Key: key, Project: channel.Project{Key: e.Key, Root: e.Root}, Methods: hostapi.Methods("test", nil)}
			}})
	}()
	time.Sleep(60 * time.Millisecond)
	m.mu.Lock()
	m.fail = nil
	m.mu.Unlock()
	waitFor(t, "told once the host answers", func() bool { return slices.Equal(m.last(), []string{alpha}) })
	cancel()
	if err := <-done; err != nil {
		t.Fatal(err)
	}
}
