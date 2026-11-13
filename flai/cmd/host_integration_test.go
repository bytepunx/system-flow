//go:build !windows

package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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
// The host and the project's server are given ports the test took from the
// kernel and closed, since they are other processes; another process may
// take one first (I-0372), so a try that loses one is tried again on new
// ports, and only a test that loses every try fails on that.
func TestHostRunsServeAndTheMCPServersAndTakesThemWithIt(t *testing.T) {
	t.Parallel()
	if testing.Short() {
		t.Skip("builds flai and runs real processes")
	}
	exe := filepath.Join(t.TempDir(), "flai")
	build := exec.Command("go", "build", "-o", exe, "..")
	if out, err := build.CombinedOutput(); err != nil {
		t.Fatalf("go build: %v\n%s", err, out)
	}
	const tries = 5
	var lost []string
	for try := 0; try < tries; try++ {
		why, ok := hostRunsServeAndTheMCPServers(t, exe)
		if ok {
			return
		}
		lost = append(lost, why)
	}
	t.Fatalf("every one of %d tries lost a port it gave away:\n%s", tries, strings.Join(lost, "\n"))
}

// hostRunsServeAndTheMCPServers is one try of
// TestHostRunsServeAndTheMCPServersAndTakesThemWithIt on fresh ports: ok when
// it showed everything, else why the try lost a port to another process.
// Any other failure fails the test.
func hostRunsServeAndTheMCPServers(t *testing.T, exe string) (string, bool) {
	t.Helper()
	tmp := t.TempDir()
	cfg := filepath.Join(tmp, "config", "config.json")
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	addr := ln.Addr().String()
	_ = ln.Close()
	// The project's MCP server gets a port of its own too, remembered in its
	// cache as flai mcp http remembers one: in the default range it races
	// the host's own MCP servers and this test running in other worktrees
	// for the first free port, and the loser exits and is restarted on the
	// same port until the winner lets it go (I-0247).
	mcpLn, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	mcpAddr := mcpLn.Addr().String()
	_ = mcpLn.Close()
	env := append(os.Environ(), "FLAI_CONFIG="+cfg, host.AddrEnv+"="+addr, "FLAI_AGENT=tester", "LOG_FORMAT=json")
	env = withoutEnv(env, host.URLEnv, host.TokenEnv)

	root := tempProject(t)
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(repo.CacheDir(), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(mcpFile(repo, mcpAddrFile), []byte(mcpAddr+"\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	keyFile := filepath.Join(tmp, "agent-key")
	_ = os.WriteFile(keyFile, []byte("k\n"), 0o600)
	sdir := serve.DirFor(cfg)
	if err := sdir.Register(serve.Entry{Key: "t", Name: "T", Root: root, URL: "http://127.0.0.1:1", KeyFile: keyFile}); err != nil {
		t.Fatal(err)
	}
	run := func(args ...string) ([]byte, []byte, error) {
		c := exec.Command(exe, append(args, "--config", cfg)...)
		c.Env, c.Dir = env, root
		out, err := c.Output()
		var stderr []byte
		var ee *exec.ExitError
		if errors.As(err, &ee) {
			stderr = ee.Stderr
		}
		return out, stderr, err
	}
	flai := func(args ...string) []byte {
		t.Helper()
		out, stderr, err := run(args...)
		if err != nil {
			t.Fatalf("flai %v: %v\n%s\n%s", args, err, out, stderr)
		}
		return out
	}
	hdir := host.DirFor(cfg)
	// startHost starts the host, answering why when another process holds
	// its port, as another host or as anything the host cannot listen beside
	startHost := func() string {
		t.Helper()
		out, stderr, err := run("host", "start")
		if err == nil {
			return ""
		}
		if strings.Contains(string(stderr), "already runs on this machine") || strings.Contains(readFile(hdir.Log()), "address already in use") {
			return fmt.Sprintf("the host's %s was taken: %s", addr, strings.TrimSpace(string(stderr)))
		}
		t.Fatalf("flai host start: %v\n%s\n%s\nhost log:\n%s", err, out, stderr, readFile(hdir.Log()))
		return ""
	}
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

	if why := startHost(); why != "" {
		return why, false
	}
	st := until("serve and the project's MCP server running", running)
	servePID, mcpPID := childOf(st, host.Serve).PID, childOf(st, host.MCP).PID
	if childOf(st, host.MCP).Root != root {
		t.Errorf("the MCP server is for %q, want %q", childOf(st, host.MCP).Root, root)
	}
	var mcpSt mcpState
	var mcpOK bool
	start := time.Now()
	waitUntilOr(t, "the MCP server's own state", func() bool {
		mcpSt, mcpOK = readMCPState(repo, time.Now())
		return mcpOK && mcpSt.PID == mcpPID
	}, func() string {
		return fmt.Sprintf("after %s: last state %+v (fresh and alive: %v), want pid %d; host status %+v\nmcp log:\n%s", time.Since(start).Round(time.Millisecond), mcpSt, mcpOK, mcpPID, status(), readFile(mcpFile(repo, mcpLogFile)))
	})
	if mcpSt.Addr != mcpAddr {
		// the host moves a project off a port another process holds (I-0271)
		if strings.Contains(readFile(hdir.Log()), "mcp address moved: another process holds it") {
			flai("host", "stop")
			return fmt.Sprintf("the project's %s was taken, and the host moved it to %s", mcpAddr, mcpSt.Addr), false
		}
		t.Errorf("the MCP server listens at %s, want the port the project remembers, %s", mcpSt.Addr, mcpAddr)
	}

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

	// the host's port was free between the kill and this start
	if why := startHost(); why != "" {
		return why, false
	}
	st = until("everything back", running)
	servePID, mcpPID = childOf(st, host.Serve).PID, childOf(st, host.MCP).PID
	flai("host", "stop")
	if serve.Alive(servePID) || serve.Alive(mcpPID) || serve.Alive(st.PID) {
		t.Errorf("after flai host stop: host %v, serve %v, mcp %v", serve.Alive(st.PID), serve.Alive(servePID), serve.Alive(mcpPID))
	}
	return "", true
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
	waitUntilOr(t, what, ok, func() string { return "" })
}

// waitUntilOr is waitUntil that adds why to its failure.
func waitUntilOr(t *testing.T, what string, ok func() bool, why func() string) {
	t.Helper()
	deadline := time.Now().Add(20 * time.Second)
	for time.Now().Before(deadline) {
		if ok() {
			return
		}
		time.Sleep(100 * time.Millisecond)
	}
	t.Fatalf("never %s: %s", what, why())
}
