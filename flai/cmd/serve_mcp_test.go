package cmd

import (
	"bytes"
	"context"
	"log/slog"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/config"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// S-0096, S-0106: flai host starts a server for every project flai serve
// serves, so they cannot all take the one default port: each takes the first
// free one from it that no other project was given, and keeps it once used.
func TestMCPAddrForAProjectFlaiServeStarts(t *testing.T) {
	t.Parallel()
	useConfig(t, filepath.Join(t.TempDir(), "cfg.json"))
	root := tempProject(t)
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	inRange := func(addr string) bool {
		host, port, err := net.SplitHostPort(addr)
		p, _ := strconv.Atoi(port)
		return err == nil && host == "127.0.0.1" && p >= 4243 && p < 4243+mcpPorts
	}
	// a port in the range taken, as by another project's server. The test
	// holds it before it asks, not after: the tests that run alongside probe
	// the same range by listening and closing, so a port mcpAddrFor found
	// free is not one the test can count on taking a moment later (CI, macOS).
	var held net.Listener
	for p := 4243; p < 4243+mcpPorts && held == nil; p++ {
		held, _ = net.Listen("tcp", net.JoinHostPort("127.0.0.1", strconv.Itoa(p)))
	}
	if held == nil {
		t.Fatalf("no port between 4243 and %d to hold", 4243+mcpPorts-1)
	}
	defer func() { _ = held.Close() }()
	taken := held.Addr().String()
	// the first free one, past the taken port
	first, err := mcpAddrFor(repo, "127.0.0.1", nil, nil)
	if err != nil || !inRange(first) || first == taken {
		t.Fatalf("first free, with %s taken: %q %v", taken, first, err)
	}
	// given to another project and not yet listening: not given again
	if other, err := mcpAddrFor(repo, "127.0.0.1", map[string]bool{first: true}, nil); err != nil || other == first || other == taken || !inRange(other) {
		t.Errorf("with %s given to another and %s taken: %q %v", first, taken, other, err)
	}
	// once a project has used an address, it keeps it
	if err := os.MkdirAll(repo.CacheDir(), 0o700); err != nil {
		t.Fatal(err)
	}
	mine := "127.0.0.1:" + unusedPort(t)
	if err := os.WriteFile(mcpFile(repo, mcpAddrFile), []byte(mine+"\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if got, _ := mcpAddrFor(repo, "127.0.0.1", nil, nil); got != mine {
		t.Errorf("remembered: %q", got)
	}
	// S-0110: a remembered address another project already has is not given twice
	moved, err := mcpAddrFor(repo, "127.0.0.1", map[string]bool{mine: true}, nil)
	if err != nil || moved == mine || !inRange(moved) {
		t.Errorf("remembered, and given to another: %q %v", moved, err)
	}
	// and the new one is none that another project remembers
	if again, err := mcpAddrFor(repo, "127.0.0.1", map[string]bool{mine: true}, map[string]bool{moved: true}); err != nil || again == moved || again == mine {
		t.Errorf("remembered by another: %q %v", again, err)
	}
}

// flai mcp http --exit-with stops the server once that process is gone: a
// flai host that was killed never stops its children itself.
func TestMCPHTTPExitsWithTheProcessItWasGiven(t *testing.T) {
	t.Setenv("FLAI_AGENT", "tester")
	useConfig(t, filepath.Join(t.TempDir(), "cfg.json"))
	inEnv(t, "FLAI_AGENT", "tester")
	old := exitWithProcessEvery
	exitWithProcessEvery = 20 * time.Millisecond
	t.Cleanup(func() { exitWithProcessEvery = old })
	root := tempProject(t)

	parent := exec.Command("sleep", "60")
	if err := parent.Start(); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = parent.Process.Kill(); _ = parent.Wait() })

	var out, errOut bytes.Buffer
	a := &app{configPath: configOf(t), out: &out, errOut: &errOut, cwd: root}
	server := newRootCmdWith(a)
	server.SetArgs([]string{"mcp", "http", "--addr", "127.0.0.1:0", "--exit-with", strconv.Itoa(parent.Process.Pid)})
	served := make(chan error, 1)
	go func() { served <- server.ExecuteContext(context.Background()) }()

	repo, _ := workitem.Open(root)
	deadline := time.Now().Add(5 * time.Second)
	for {
		if _, alive := readMCPState(repo, time.Now()); alive {
			break
		}
		if time.Now().After(deadline) {
			t.Fatalf("the server did not come up: %s %s", out.String(), errOut.String())
		}
		time.Sleep(20 * time.Millisecond)
	}
	select {
	case err := <-served:
		t.Fatalf("stopped while its process still ran: %v", err)
	case <-time.After(200 * time.Millisecond):
	}

	_ = parent.Process.Kill()
	_ = parent.Wait()
	select {
	case err := <-served:
		if err != nil {
			t.Errorf("stopping: %v", err)
		}
	case <-time.After(10 * time.Second):
		t.Fatal("the server outlived the process it was told to exit with")
	}
}

// S-0110: flai host gives two projects that remembered the same MCP address
// one each: the first keeps it, the second gets a free one, and the log says so.
func TestTheHostGivesTwoProjectsThatRememberOnePortOneEach(t *testing.T) {
	t.Parallel()
	useConfig(t, filepath.Join(t.TempDir(), "cfg.json"))
	var logged bytes.Buffer
	l := &hostLauncher{a: &app{configPath: configOf(t), log: slog.New(slog.NewTextHandler(&logged, nil))}, taken: map[string]string{}}
	// a port no other test probes: the default range is the live host's and
	// the parallel tests'
	shared := "127.0.0.1:" + unusedPort(t)
	var repos []*workitem.Repo
	for i := 0; i < 2; i++ {
		repo, err := workitem.Open(tempProject(t))
		if err != nil {
			t.Fatal(err)
		}
		if err := os.MkdirAll(repo.CacheDir(), 0o700); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(mcpFile(repo, mcpAddrFile), []byte(shared+"\n"), 0o600); err != nil {
			t.Fatal(err)
		}
		repos = append(repos, repo)
	}
	// a third project flai serve serves, not running, remembers 4244: a move does not take it
	third := tempProject(t)
	thirdRepo, _ := workitem.Open(third)
	_ = os.MkdirAll(thirdRepo.CacheDir(), 0o700)
	_ = os.WriteFile(mcpFile(thirdRepo, mcpAddrFile), []byte("127.0.0.1:4244\n"), 0o600)
	serveDir := filepath.Join(filepath.Dir(configOf(t)), "serve")
	_ = os.MkdirAll(serveDir, 0o700)
	_ = os.WriteFile(filepath.Join(serveDir, "projects.json"), []byte(`[{"key":"c","name":"c","root":"`+third+`","url":"http://127.0.0.1:1"}]`), 0o600)
	first, err := l.mcpAddr(repos[0])
	if err != nil || first != shared {
		t.Fatalf("the first keeps what it remembered: %q %v", first, err)
	}
	second, err := l.mcpAddr(repos[1])
	if err != nil || second == first || second == "127.0.0.1:4244" {
		t.Fatalf("the second gets another, and not one the third remembers: %q %v", second, err)
	}
	if again, _ := l.mcpAddr(repos[1]); again != second {
		t.Errorf("and keeps it while the host runs: %q", again)
	}
	if !strings.Contains(logged.String(), "mcp address moved") || !strings.Contains(logged.String(), "to="+second) {
		t.Errorf("the log says so:\n%s", logged.String())
	}
}

// S-0506: mcpAddrFor puts the project's port on the host's mcp.bind, the
// remembered port included, and tells ports apart whatever host the taken and
// avoided addresses name, so a change of mcp.bind moves no project off its
// port and gives no project another's.
func TestMCPAddrForListensOnTheBindAddress(t *testing.T) {
	t.Parallel()
	useConfig(t, filepath.Join(t.TempDir(), "cfg.json"))
	repo, err := workitem.Open(tempProject(t))
	if err != nil {
		t.Fatal(err)
	}
	first, err := mcpAddrFor(repo, "0.0.0.0", nil, nil)
	if host, _, _ := net.SplitHostPort(first); err != nil || host != "0.0.0.0" {
		t.Fatalf("a new port, on the bind address: %q %v", first, err)
	}
	if err := os.MkdirAll(repo.CacheDir(), 0o700); err != nil {
		t.Fatal(err)
	}
	port := unusedPort(t)
	if err := os.WriteFile(mcpFile(repo, mcpAddrFile), []byte("127.0.0.1:"+port+"\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if got, err := mcpAddrFor(repo, "0.0.0.0", nil, nil); err != nil || got != "0.0.0.0:"+port {
		t.Errorf("the remembered port, on the bind address: %q %v", got, err)
	}
	// the remembered port must be free to listen on, so only where this
	// machine has IPv6 loopback
	if ln, err := net.Listen("tcp", "[::1]:0"); err == nil {
		_ = ln.Close()
		if got, err := mcpAddrFor(repo, "::1", nil, nil); err != nil || got != "[::1]:"+port {
			t.Errorf("an IPv6 bind address is bracketed: %q %v", got, err)
		}
	}
	moved, err := mcpAddrFor(repo, "0.0.0.0", map[string]bool{"127.0.0.1:" + port: true}, nil)
	if err != nil || moved == "0.0.0.0:"+port {
		t.Errorf("a port another project has on another host is not given twice: %q %v", moved, err)
	}
	if again, err := mcpAddrFor(repo, "0.0.0.0", map[string]bool{"127.0.0.1:" + port: true}, map[string]bool{"127.0.0.1:" + portOf(moved): true}); err != nil || again == moved {
		t.Errorf("nor one another project remembers on another host: %q %v", again, err)
	}
}

// S-0506: flai host starts each project's server on mcp.bind, read at every
// start: a project keeps its port when the setting changes, the log says
// nothing moved, and clearing it goes back to loopback.
func TestTheHostStartsEachProjectsServerOnTheBindAddress(t *testing.T) {
	t.Parallel()
	cfgPath := filepath.Join(t.TempDir(), "cfg.json")
	useConfig(t, cfgPath)
	var logged bytes.Buffer
	l := &hostLauncher{a: &app{configPath: cfgPath, log: slog.New(slog.NewTextHandler(&logged, nil))}, taken: map[string]string{}}
	repo, err := workitem.Open(tempProject(t))
	if err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(repo.CacheDir(), 0o700); err != nil {
		t.Fatal(err)
	}
	port := unusedPort(t)
	if err := os.WriteFile(mcpFile(repo, mcpAddrFile), []byte("127.0.0.1:"+port+"\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if got, err := l.mcpAddr(repo); err != nil || got != "127.0.0.1:"+port {
		t.Fatalf("unset, loopback: %q %v", got, err)
	}
	bindTo := func(addr string) {
		t.Helper()
		cfg, _, err := config.Load(cfgPath)
		if err != nil {
			t.Fatal(err)
		}
		if addr == "" {
			cfg = cfg.WithoutMCPBind()
		} else if cfg, err = cfg.WithMCPBind(addr); err != nil {
			t.Fatal(err)
		}
		if err := config.Save(cfgPath, cfg); err != nil {
			t.Fatal(err)
		}
	}
	bindTo("0.0.0.0")
	if got, err := l.mcpAddr(repo); err != nil || got != "0.0.0.0:"+port {
		t.Errorf("set, the same port on the bind address: %q %v", got, err)
	}
	other, err := workitem.Open(tempProject(t))
	if err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(other.CacheDir(), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(mcpFile(other, mcpAddrFile), []byte("0.0.0.0:"+port+"\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if got, err := l.mcpAddr(other); err != nil || portOf(got) == port {
		t.Errorf("another project that remembers that port gets another: %q %v", got, err)
	}
	bindTo("")
	if got, err := l.mcpAddr(repo); err != nil || got != "127.0.0.1:"+port {
		t.Errorf("cleared, loopback again: %q %v", got, err)
	}
	if strings.Contains(logged.String(), "root="+repo.MainRoot) {
		t.Errorf("a change of the bind address is no move:\n%s", logged.String())
	}
}

// unusedPort is a port the kernel just gave and took back, outside the
// default range that the live host and the parallel tests probe, so a test
// can remember it and count on its being free to listen on.
func unusedPort(t *testing.T) string {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	port := portOf(ln.Addr().String())
	_ = ln.Close()
	return port
}

// I-0271: another process takes the port a project's server had, or
// remembered, while the server is not running. The host gave the project
// that port again at every start, and the server failed to listen on it
// until the other process let go. Now it moves the project to a free port,
// as a new project would get, keeps that one, and the log says so.
func TestTheHostMovesAProjectOffAPortAnotherProcessHolds(t *testing.T) {
	t.Parallel()
	useConfig(t, filepath.Join(t.TempDir(), "cfg.json"))
	var logged bytes.Buffer
	l := &hostLauncher{a: &app{configPath: configOf(t), log: slog.New(slog.NewTextHandler(&logged, nil))}, taken: map[string]string{}}
	repo, err := workitem.Open(tempProject(t))
	if err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(repo.CacheDir(), 0o700); err != nil {
		t.Fatal(err)
	}
	mine := "127.0.0.1:" + unusedPort(t)
	if err := os.WriteFile(mcpFile(repo, mcpAddrFile), []byte(mine+"\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if got, err := l.mcpAddr(repo); err != nil || got != mine {
		t.Fatalf("the first start, on the remembered port: %q %v", got, err)
	}
	// the server exited, and another process took its port before the restart
	holder, err := net.Listen("tcp", mine)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = holder.Close() }()
	moved, err := l.mcpAddr(repo)
	if err != nil || moved == mine || freePort(mine) {
		t.Fatalf("the restart, with %s held: %q %v", mine, moved, err)
	}
	if !strings.Contains(logged.String(), "mcp address moved: another process holds it") ||
		!strings.Contains(logged.String(), "from="+mine) || !strings.Contains(logged.String(), "to="+moved) {
		t.Errorf("the log names the held port and the new one:\n%s", logged.String())
	}
	if again, err := l.mcpAddr(repo); err != nil || again != moved {
		t.Errorf("and the project keeps the new one while it is free: %q %v", again, err)
	}
	// a host started again, while the port the server remembers is still held
	fresh := &hostLauncher{a: l.a, taken: map[string]string{}}
	if got, err := fresh.mcpAddr(repo); err != nil || got == mine {
		t.Errorf("a new host, with the remembered %s held: %q %v", mine, got, err)
	}
	if got, err := mcpAddrFor(repo, "127.0.0.1", nil, nil); err != nil || got == mine {
		t.Errorf("mcpAddrFor, with the remembered %s held: %q %v", mine, got, err)
	}
}
