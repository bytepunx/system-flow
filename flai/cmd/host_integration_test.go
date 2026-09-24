//go:build !windows

package cmd

import (
	"encoding/json"
	"errors"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/host"
	"github.com/bytepunx/system-flow/flai/internal/serve"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// S-0106, end to end with a real flai built from this tree: flai host starts
// flai serve, keeps the MCP server of the project flai serve serves, restarts
// serve when asked, and when it is killed outright, its children go with it.
func TestHostRunsServeAndTheMCPServersAndTakesThemWithIt(t *testing.T) {
	if testing.Short() {
		t.Skip("builds flai and runs real processes")
	}
	tmp := t.TempDir()
	exe := filepath.Join(tmp, "flai")
	build := exec.Command("go", "build", "-o", exe, "..")
	if out, err := build.CombinedOutput(); err != nil {
		t.Fatalf("go build: %v\n%s", err, out)
	}
	cfg := filepath.Join(tmp, "config", "config.json")
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	addr := ln.Addr().String()
	_ = ln.Close()
	env := append(os.Environ(), "FLAI_CONFIG="+cfg, host.AddrEnv+"="+addr, "FLAI_AGENT=tester", "LOG_FORMAT=json")
	env = withoutEnv(env, host.URLEnv, host.TokenEnv)

	root := tempProject(t)
	keyFile := filepath.Join(tmp, "agent-key")
	_ = os.WriteFile(keyFile, []byte("k\n"), 0o600)
	sdir := serve.DirFor(cfg)
	if err := sdir.Register(serve.Entry{Key: "t", Name: "T", Root: root, URL: "http://127.0.0.1:1", KeyFile: keyFile}); err != nil {
		t.Fatal(err)
	}
	flai := func(args ...string) []byte {
		t.Helper()
		c := exec.Command(exe, append(args, "--config", cfg)...)
		c.Env, c.Dir = env, root
		out, err := c.Output()
		if err != nil {
			var stderr []byte
			var ee *exec.ExitError
			if errors.As(err, &ee) {
				stderr = ee.Stderr
			}
			t.Fatalf("flai %v: %v\n%s\n%s", args, err, out, stderr)
		}
		return out
	}
	hdir := host.DirFor(cfg)
	t.Cleanup(func() {
		if st, ok := hdir.ReadStatus(time.Now()); ok {
			_ = syscall.Kill(st.PID, syscall.SIGKILL)
		}
	})
	status := func() host.Status {
		var st hostStatus
		if err := json.Unmarshal(flai("host", "status", "--json"), &st); err != nil || st.Status == nil {
			return host.Status{}
		}
		return *st.Status
	}
	until := func(what string, ok func(host.Status) bool) host.Status {
		t.Helper()
		deadline := time.Now().Add(20 * time.Second)
		var st host.Status
		for time.Now().Before(deadline) {
			if st = status(); ok(st) {
				return st
			}
			time.Sleep(200 * time.Millisecond)
		}
		t.Fatalf("never %s: %+v\nhost log:\n%s", what, st, readFile(hdir.Log()))
		return st
	}
	childOf := func(st host.Status, name string) host.Child {
		for _, c := range st.Children {
			if c.Name == name {
				return c
			}
		}
		return host.Child{}
	}
	running := func(st host.Status) bool {
		return childOf(st, host.Serve).State == "running" && childOf(st, host.MCP).State == "running"
	}

	flai("host", "start")
	st := until("serve and the project's MCP server running", running)
	servePID, mcpPID := childOf(st, host.Serve).PID, childOf(st, host.MCP).PID
	if childOf(st, host.MCP).Root != root {
		t.Errorf("the MCP server is for %q, want %q", childOf(st, host.MCP).Root, root)
	}
	repo, _ := workitem.Open(root)
	waitUntil(t, "the MCP server's own state", func() bool { s, ok := readMCPState(repo, time.Now()); return ok && s.PID == mcpPID })

	flai("host", "restart", "serve")
	st = until("a new serve", func(st host.Status) bool {
		c := childOf(st, host.Serve)
		return c.State == "running" && c.PID != servePID
	})
	if childOf(st, host.MCP).PID != mcpPID {
		t.Errorf("restarting serve restarted the MCP server: %d, was %d", childOf(st, host.MCP).PID, mcpPID)
	}
	servePID = childOf(st, host.Serve).PID

	// killed outright, the host stops nothing itself: its children go by themselves
	hostPID := st.PID
	if err := syscall.Kill(hostPID, syscall.SIGKILL); err != nil {
		t.Fatal(err)
	}
	waitUntil(t, "serve and the MCP server gone with the host", func() bool { return !serve.Alive(servePID) && !serve.Alive(mcpPID) })

	flai("host", "start")
	st = until("everything back", running)
	servePID, mcpPID = childOf(st, host.Serve).PID, childOf(st, host.MCP).PID
	flai("host", "stop")
	if serve.Alive(servePID) || serve.Alive(mcpPID) || serve.Alive(st.PID) {
		t.Errorf("after flai host stop: host %v, serve %v, mcp %v", serve.Alive(st.PID), serve.Alive(servePID), serve.Alive(mcpPID))
	}
}

func withoutEnv(env []string, names ...string) []string {
	out := env[:0:0]
next:
	for _, kv := range env {
		for _, n := range names {
			if len(kv) > len(n) && kv[:len(n)+1] == n+"=" {
				continue next
			}
		}
		out = append(out, kv)
	}
	return out
}

func readFile(p string) string {
	data, _ := os.ReadFile(p)
	return string(data)
}

func waitUntil(t *testing.T, what string, ok func() bool) {
	t.Helper()
	deadline := time.Now().Add(20 * time.Second)
	for time.Now().Before(deadline) {
		if ok() {
			return
		}
		time.Sleep(100 * time.Millisecond)
	}
	t.Fatalf("never %s", what)
}
