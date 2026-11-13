package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"slices"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/config"
	"github.com/bytepunx/system-flow/flai/internal/execx"
	"github.com/bytepunx/system-flow/flai/internal/harness"
	"github.com/bytepunx/system-flow/flai/internal/serve"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// S-0096, S-0106: flai host starts a server for every project flai serve
// serves, so they cannot all take the one default port: each takes the first
// free one from it that no other project was given, and keeps it once used.
func TestMCPAddrsForAProjectFlaiServeStarts(t *testing.T) {
	t.Parallel()
	useConfig(t, filepath.Join(t.TempDir(), "cfg.json"))
	root := tempProject(t)
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	// the default range is probed for real, with a port in it the test holds,
	// so that what freePort does with the operating system stays covered
	mcpAddrFor := func(taken, avoid map[string]bool) (string, error) {
		return oneAddr(mcpAddrsFor(repo, []string{"127.0.0.1"}, taken, avoid, freePort))
	}
	// a remembered port is one the test says is free, not one it closed and
	// hopes no parallel test takes (I-0372)
	var ports heldPorts
	remembering := func(taken, avoid map[string]bool) (string, error) {
		return oneAddr(mcpAddrsFor(repo, []string{"127.0.0.1"}, taken, avoid, ports.free))
	}
	inRange := func(addr string) bool {
		host, port, err := net.SplitHostPort(addr)
		p, _ := strconv.Atoi(port)
		return err == nil && host == "127.0.0.1" && p >= 4243 && p < 4243+mcpPorts
	}
	// a port in the range taken, as by another project's server. The test
	// holds it before it asks, not after: the tests that run alongside probe
	// the same range by listening and closing, so a port mcpAddrsFor found
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
	first, err := mcpAddrFor(nil, nil)
	if err != nil || !inRange(first) || first == taken {
		t.Fatalf("first free, with %s taken: %q %v", taken, first, err)
	}
	// given to another project and not yet listening: not given again
	if other, err := mcpAddrFor(map[string]bool{first: true}, nil); err != nil || other == first || other == taken || !inRange(other) {
		t.Errorf("with %s given to another and %s taken: %q %v", first, taken, other, err)
	}
	// once a project has used an address, it keeps it
	if err := os.MkdirAll(repo.CacheDir(), 0o700); err != nil {
		t.Fatal(err)
	}
	mine := "127.0.0.1:" + rememberedPort
	if err := os.WriteFile(mcpFile(repo, mcpAddrFile), []byte(mine+"\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if got, _ := remembering(nil, nil); got != mine {
		t.Errorf("remembered: %q", got)
	}
	// S-0110: a remembered address another project already has is not given twice
	moved, err := remembering(map[string]bool{mine: true}, nil)
	if err != nil || moved == mine || !inRange(moved) {
		t.Errorf("remembered, and given to another: %q %v", moved, err)
	}
	// and the new one is none that another project remembers
	if again, err := remembering(map[string]bool{mine: true}, map[string]bool{moved: true}); err != nil || again == moved || again == mine {
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
	var ports heldPorts
	l := &hostLauncher{a: &app{configPath: configOf(t), log: slog.New(slog.NewTextHandler(&logged, nil)), runner: desktopDocker}, taken: map[string]string{}, free: ports.free}
	shared := "127.0.0.1:" + rememberedPort
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
	first, err := oneAddr(l.mcpAddrs(repos[0]))
	if err != nil || first != shared {
		t.Fatalf("the first keeps what it remembered: %q %v", first, err)
	}
	second, err := oneAddr(l.mcpAddrs(repos[1]))
	if err != nil || second == first || second == "127.0.0.1:4244" {
		t.Fatalf("the second gets another, and not one the third remembers: %q %v", second, err)
	}
	if again, _ := oneAddr(l.mcpAddrs(repos[1])); again != second {
		t.Errorf("and keeps it while the host runs: %q", again)
	}
	if !strings.Contains(logged.String(), "mcp address moved") || !strings.Contains(logged.String(), "to="+second) {
		t.Errorf("the log says so:\n%s", logged.String())
	}
}

// S-0506: mcpAddrsFor puts the project's port on the host's mcp.bind, the
// remembered port included, and tells ports apart whatever host the taken and
// avoided addresses name, so a change of mcp.bind moves no project off its
// port and gives no project another's.
func TestMCPAddrsForListensOnTheBindAddress(t *testing.T) {
	t.Parallel()
	useConfig(t, filepath.Join(t.TempDir(), "cfg.json"))
	repo, err := workitem.Open(tempProject(t))
	if err != nil {
		t.Fatal(err)
	}
	var ports heldPorts
	mcpAddrFor := func(bind string, taken, avoid map[string]bool) (string, error) {
		return oneAddr(mcpAddrsFor(repo, []string{bind}, taken, avoid, ports.free))
	}
	first, err := mcpAddrFor("0.0.0.0", nil, nil)
	if host, _, _ := net.SplitHostPort(first); err != nil || host != "0.0.0.0" {
		t.Fatalf("a new port, on the bind address: %q %v", first, err)
	}
	if err := os.MkdirAll(repo.CacheDir(), 0o700); err != nil {
		t.Fatal(err)
	}
	port := rememberedPort
	if err := os.WriteFile(mcpFile(repo, mcpAddrFile), []byte("127.0.0.1:"+port+"\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if got, err := mcpAddrFor("0.0.0.0", nil, nil); err != nil || got != "0.0.0.0:"+port {
		t.Errorf("the remembered port, on the bind address: %q %v", got, err)
	}
	if got, err := mcpAddrFor("::1", nil, nil); err != nil || got != "[::1]:"+port {
		t.Errorf("an IPv6 bind address is bracketed: %q %v", got, err)
	}
	moved, err := mcpAddrFor("0.0.0.0", map[string]bool{"127.0.0.1:" + port: true}, nil)
	if err != nil || moved == "0.0.0.0:"+port {
		t.Errorf("a port another project has on another host is not given twice: %q %v", moved, err)
	}
	if again, err := mcpAddrFor("0.0.0.0", map[string]bool{"127.0.0.1:" + port: true}, map[string]bool{"127.0.0.1:" + portOf(moved): true}); err != nil || again == moved {
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
	var ports heldPorts
	l := &hostLauncher{a: &app{configPath: cfgPath, log: slog.New(slog.NewTextHandler(&logged, nil)), runner: desktopDocker}, taken: map[string]string{}, free: ports.free}
	repo, err := workitem.Open(tempProject(t))
	if err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(repo.CacheDir(), 0o700); err != nil {
		t.Fatal(err)
	}
	port := rememberedPort
	if err := os.WriteFile(mcpFile(repo, mcpAddrFile), []byte("127.0.0.1:"+port+"\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if got, err := l.mcpAddrs(repo); err != nil || !slices.Equal(got, []string{"127.0.0.1:" + port}) {
		t.Fatalf("unset, loopback under Docker Desktop: %q %v", got, err)
	}
	setMCPBind(t, cfgPath, "0.0.0.0")
	if got, err := l.mcpAddrs(repo); err != nil || !slices.Equal(got, []string{"0.0.0.0:" + port}) {
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
	if got, err := l.mcpAddrs(other); err != nil || len(got) != 1 || portOf(got[0]) == port {
		t.Errorf("another project that remembers that port gets another: %q %v", got, err)
	}
	setMCPBind(t, cfgPath)
	if got, err := l.mcpAddrs(repo); err != nil || !slices.Equal(got, []string{"127.0.0.1:" + port}) {
		t.Errorf("cleared, loopback again: %q %v", got, err)
	}
	if strings.Contains(logged.String(), "root="+repo.MainRoot) {
		t.Errorf("a change of the bind address is no move:\n%s", logged.String())
	}
}

// oneAddr is the one address of addrs, for a test on one bind address,
// and an error when there is not one.
func oneAddr(addrs []string, err error) (string, error) {
	if err != nil {
		return "", err
	}
	if len(addrs) != 1 {
		return "", fmt.Errorf("want one address, got %q", addrs)
	}
	return addrs[0], nil
}

// setMCPBind sets mcp.bind in the configuration at cfgPath to addrs, or
// unsets it when there are none.
func setMCPBind(t *testing.T, cfgPath string, addrs ...string) {
	t.Helper()
	cfg, _, err := config.Load(cfgPath)
	if err != nil {
		t.Fatal(err)
	}
	if len(addrs) == 0 {
		cfg = cfg.WithoutMCPBind()
	} else if cfg, err = cfg.WithMCPBind(addrs...); err != nil {
		t.Fatal(err)
	}
	if err := config.Save(cfgPath, cfg); err != nil {
		t.Fatal(err)
	}
}

// hostLog is one line of what the host logs, the fields the tests pin.
type hostLog struct {
	Level     string `json:"level"`
	Component string `json:"component"`
	Root      string `json:"root"`
	Addr      string `json:"addr"`
	From      string `json:"from"`
	To        string `json:"to"`
	Err       string `json:"err"`
}

// hostLogs is the lines of component host a JSON handler wrote to logged:
// loading the configuration logs lines of its own.
func hostLogs(t *testing.T, logged *bytes.Buffer) []hostLog {
	t.Helper()
	var lines []hostLog
	for _, line := range strings.Split(strings.TrimSpace(logged.String()), "\n") {
		if line == "" {
			continue
		}
		var e hostLog
		if err := json.Unmarshal([]byte(line), &e); err != nil {
			t.Fatalf("not a JSON line: %s", line)
		}
		if e.Component == "host" {
			lines = append(lines, e)
		}
	}
	return lines
}

// specAddrs is the addresses a spec's flai mcp http is given, in order.
func specAddrs(args []string) []string {
	var addrs []string
	for i, arg := range args {
		if arg == "--addr" && i+1 < len(args) {
			addrs = append(addrs, args[i+1])
		}
	}
	return addrs
}

// needLoopback2 skips a test that listens on 127.0.0.2 where this machine
// cannot, as macOS without an alias cannot.
func needLoopback2(t *testing.T) {
	t.Helper()
	if err := listenErr("127.0.0.2:0"); err != nil {
		t.Skipf("cannot listen on 127.0.0.2: %v", err)
	}
}

// mcpHostFor is a host launcher on the configuration at cfgPath, with
// docker, logging as JSON to logged, asking ports which addresses are free,
// and a project that remembers rememberedPort, which it answers.
func mcpHostFor(t *testing.T, cfgPath string, docker execx.Runner, logged *bytes.Buffer, ports *heldPorts) (*hostLauncher, *workitem.Repo, string) {
	t.Helper()
	l := &hostLauncher{a: &app{configPath: cfgPath, log: slog.New(slog.NewJSONHandler(logged, nil)), runner: docker}, taken: map[string]string{}, free: ports.free}
	repo, err := workitem.Open(tempProject(t))
	if err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(repo.CacheDir(), 0o700); err != nil {
		t.Fatal(err)
	}
	port := rememberedPort
	if err := os.WriteFile(mcpFile(repo, mcpAddrFile), []byte("127.0.0.1:"+port+"\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	return l, repo, port
}

// S-0732: with mcp.bind holding two addresses, the host starts the project's
// server on both, on one port; a port another process holds on the second
// alone moves the project to one free on both, and the log says from and to.
func TestTheHostStartsEachProjectsServerOnEveryBindAddress(t *testing.T) {
	t.Parallel()
	needLoopback2(t)
	cfgPath := filepath.Join(t.TempDir(), "cfg.json")
	useConfig(t, cfgPath)
	setMCPBind(t, cfgPath, "127.0.0.1", "127.0.0.2")
	var logged bytes.Buffer
	var ports heldPorts
	l, repo, port := mcpHostFor(t, cfgPath, desktopDocker, &logged, &ports)
	spec, err := l.mcp(repo.MainRoot)
	if want := []string{"127.0.0.1:" + port, "127.0.0.2:" + port}; err != nil || !slices.Equal(specAddrs(spec.Args), want) {
		t.Fatalf("the spec listens on %q (%v), want %q: %q", specAddrs(spec.Args), err, want, spec.Args)
	}
	ports.hold("127.0.0.2:" + port)
	spec, err = l.mcp(repo.MainRoot)
	moved := specAddrs(spec.Args)
	if err != nil || len(moved) != 2 || portOf(moved[0]) == port || portOf(moved[0]) != portOf(moved[1]) ||
		moved[0] != "127.0.0.1:"+portOf(moved[0]) || moved[1] != "127.0.0.2:"+portOf(moved[0]) {
		t.Fatalf("with %s held on 127.0.0.2 alone: %q %v", port, moved, err)
	}
	lines := hostLogs(t, &logged)
	if len(lines) != 1 || lines[0].Level != "WARN" || lines[0].Root != repo.MainRoot ||
		lines[0].From != "127.0.0.1:"+port || lines[0].To != moved[0] {
		t.Errorf("the log says the project moved, from and to: %+v", lines)
	}
}

// S-0732: auto, resolved at each start, listens on 127.0.0.1 and, under a
// native Docker engine, on the container network's gateway, here 127.0.0.2
// so that the test can listen on it; when docker cannot say, on 127.0.0.1
// alone, and the host warns with what docker could not say.
func TestTheHostResolvesAutoAtEachStart(t *testing.T) {
	t.Parallel()
	needLoopback2(t)
	cfgPath := filepath.Join(t.TempDir(), "cfg.json")
	useConfig(t, cfgPath)
	var logged bytes.Buffer
	l, repo, port := mcpHostFor(t, cfgPath, mcpDocker{engine: "Ubuntu 24.04 LTS", gateway: "127.0.0.2"}, &logged, &heldPorts{})
	spec, err := l.mcp(repo.MainRoot)
	if want := []string{"127.0.0.1:" + port, "127.0.0.2:" + port}; err != nil || !slices.Equal(specAddrs(spec.Args), want) {
		t.Fatalf("auto on a native engine listens on %q (%v), want %q", specAddrs(spec.Args), err, want)
	}
	if lines := hostLogs(t, &logged); len(lines) != 0 {
		t.Errorf("nothing to warn of: %+v", lines)
	}
	l.a.runner = mcpDocker{}
	spec, err = l.mcp(repo.MainRoot)
	if want := []string{"127.0.0.1:" + port}; err != nil || !slices.Equal(specAddrs(spec.Args), want) {
		t.Fatalf("auto without docker listens on %q (%v), want %q", specAddrs(spec.Args), err, want)
	}
	lines := hostLogs(t, &logged)
	if len(lines) != 1 || lines[0].Level != "WARN" || lines[0].Root != repo.MainRoot ||
		!strings.Contains(lines[0].Err, "Cannot connect to the Docker daemon") || lines[0].Addr != "" {
		t.Errorf("the host warns with what docker could not say: %+v", lines)
	}
}

// S-0732: an address this machine cannot listen on at all, such as a
// gateway whose bridge is down, is dropped with a warning naming it while
// another remains, and fails the start, saying how to set another, when it
// is the only one.
func TestTheHostDropsAnAddressThisMachineCannotListenOn(t *testing.T) {
	t.Parallel()
	const away = "192.0.2.1" // TEST-NET-1, assigned to no machine
	if err := listenErr(away + ":0"); err == nil {
		t.Skipf("this machine can listen on %s", away)
	}
	cfgPath := filepath.Join(t.TempDir(), "cfg.json")
	useConfig(t, cfgPath)
	setMCPBind(t, cfgPath, "127.0.0.1", away)
	var logged bytes.Buffer
	l, repo, port := mcpHostFor(t, cfgPath, desktopDocker, &logged, &heldPorts{})
	spec, err := l.mcp(repo.MainRoot)
	if want := []string{"127.0.0.1:" + port}; err != nil || !slices.Equal(specAddrs(spec.Args), want) {
		t.Fatalf("with %s beside loopback, the spec listens on %q (%v), want %q", away, specAddrs(spec.Args), err, want)
	}
	lines := hostLogs(t, &logged)
	if len(lines) != 1 || lines[0].Level != "WARN" || lines[0].Root != repo.MainRoot ||
		lines[0].Addr != away || lines[0].Err == "" {
		t.Errorf("the host warns of the address it dropped: %+v", lines)
	}
	setMCPBind(t, cfgPath, away)
	if _, err := l.mcp(repo.MainRoot); err == nil || !strings.Contains(err.Error(), away) || !strings.Contains(err.Error(), "flai serve mcp bind") {
		t.Errorf("with %s alone, the start fails saying how to set another: %v", away, err)
	}
}

// rememberedPort is the port a test's project remembers, outside the default
// range. The test's heldPorts says whether it is free, so no real port is
// raced for: a port the kernel gave and the test closed again was taken by a
// parallel test before flai probed it (I-0372).
const rememberedPort = "4300"

// heldPorts is a free-port check a test controls: every address is free but
// those the test holds, by host and port as the address names them.
type heldPorts struct{ held map[string]bool }

// free says whether addr is free, as freePort does of a real one.
func (p *heldPorts) free(addr string) bool { return !p.held[addr] }

// hold makes addr held, as by another process.
func (p *heldPorts) hold(addr string) {
	if p.held == nil {
		p.held = map[string]bool{}
	}
	p.held[addr] = true
}

// freePort asks the operating system: an address the test holds a listener
// on is not free while it holds it, and a launcher with no check of its own
// asks the same way.
func TestFreePortSaysWhetherAnAddressIsHeld(t *testing.T) {
	t.Parallel()
	held, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = held.Close() }()
	addr := held.Addr().String()
	if freePort(addr) {
		t.Errorf("%s is held by the test, but freePort says it is free", addr)
	}
	if (&hostLauncher{}).freeCheck()(addr) {
		t.Errorf("a launcher with no check of its own says the held %s is free", addr)
	}
	if !freePort("127.0.0.1:0") {
		t.Error("freePort says no port the kernel picks is free")
	}
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
	var ports heldPorts
	l := &hostLauncher{a: &app{configPath: configOf(t), log: slog.New(slog.NewTextHandler(&logged, nil)), runner: desktopDocker}, taken: map[string]string{}, free: ports.free}
	repo, err := workitem.Open(tempProject(t))
	if err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(repo.CacheDir(), 0o700); err != nil {
		t.Fatal(err)
	}
	mine := "127.0.0.1:" + rememberedPort
	if err := os.WriteFile(mcpFile(repo, mcpAddrFile), []byte(mine+"\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if got, err := oneAddr(l.mcpAddrs(repo)); err != nil || got != mine {
		t.Fatalf("the first start, on the remembered port: %q %v", got, err)
	}
	// the server exited, and another process took its port before the restart
	ports.hold(mine)
	moved, err := oneAddr(l.mcpAddrs(repo))
	if err != nil || moved == mine || !ports.free(moved) {
		t.Fatalf("the restart, with %s held: %q %v", mine, moved, err)
	}
	if !strings.Contains(logged.String(), "mcp address moved: another process holds it") ||
		!strings.Contains(logged.String(), "from="+mine) || !strings.Contains(logged.String(), "to="+moved) {
		t.Errorf("the log names the held port and the new one:\n%s", logged.String())
	}
	if again, err := oneAddr(l.mcpAddrs(repo)); err != nil || again != moved {
		t.Errorf("and the project keeps the new one while it is free: %q %v", again, err)
	}
	// a host started again, while the port the server remembers is still held
	fresh := &hostLauncher{a: l.a, taken: map[string]string{}, free: ports.free}
	if got, err := oneAddr(fresh.mcpAddrs(repo)); err != nil || got == mine {
		t.Errorf("a new host, with the remembered %s held: %q %v", mine, got, err)
	}
	if got, err := oneAddr(mcpAddrsFor(repo, []string{"127.0.0.1"}, nil, nil, ports.free)); err != nil || got == mine {
		t.Errorf("mcpAddrsFor, with the remembered %s held: %q %v", mine, got, err)
	}
}

// mcpDocker is a docker for what resolves mcp.bind's auto: docker info names
// engine, and fails when it is empty, as without a running Docker; network
// inspect names gateway for any network, 172.17.0.1 when it is empty. calls
// records each call's arguments, so a test can see which network it asked
// about.
type mcpDocker struct {
	engine  string
	gateway string
	calls   *[]string
}

// desktopDocker is a Docker Desktop engine, under which auto is loopback
// alone, for a test that shows the settings without looking at Docker.
var desktopDocker = mcpDocker{engine: serve.DockerDesktop}

func (d mcpDocker) Run(_, _ string, args ...string) (string, error) {
	if d.calls != nil {
		*d.calls = append(*d.calls, strings.Join(args, " "))
	}
	switch {
	case len(args) > 0 && args[0] == "info":
		if d.engine == "" {
			return "Cannot connect to the Docker daemon", errors.New("exit status 1")
		}
		return d.engine + "\n", nil
	case len(args) > 1 && args[0] == "network" && args[1] == "inspect":
		if d.gateway != "" {
			return d.gateway + " \n", nil
		}
		return "172.17.0.1 \n", nil
	}
	return "", fmt.Errorf("mcpDocker: no answer for docker %s", strings.Join(args, " "))
}

func (d mcpDocker) RunInput(dir, name, _ string, args ...string) (string, error) {
	return d.Run(dir, name, args...)
}

func (d mcpDocker) LookPath(name string) (string, error) { return name, nil }

// runServeMCP runs flai serve mcp with args on the configuration at cfgPath
// and docker, as JSON when asJSON.
func runServeMCP(t *testing.T, cfgPath string, docker execx.Runner, asJSON bool, args ...string) (string, error) {
	t.Helper()
	var out bytes.Buffer
	a := &app{configPath: cfgPath, out: &out, errOut: &bytes.Buffer{}, cwd: t.TempDir(), jsonOut: asJSON, runner: docker}
	c := newServeMCPCmd(a)
	c.SetArgs(args)
	c.SetOut(&out)
	c.SetErr(&bytes.Buffer{})
	c.SilenceUsage, c.SilenceErrors = true, true
	err := c.Execute()
	return out.String(), err
}

// showMCPJSON is flai serve mcp show --json read back.
func showMCPJSON(t *testing.T, cfgPath string, docker execx.Runner) mcpBindView {
	t.Helper()
	js, err := runServeMCP(t, cfgPath, docker, true, "show")
	if err != nil {
		t.Fatalf("show --json: %v", err)
	}
	var v mcpBindView
	if err := json.Unmarshal([]byte(js), &v); err != nil {
		t.Fatalf("show --json is not JSON: %v\n%s", err, js)
	}
	return v
}

// S-0732: flai serve mcp bind takes several addresses and saves them as a
// list, show prints them as written, --json gives them as binds and listen,
// and an address it refuses among several leaves mcp.bind as it was.
func TestServeMCPBindTakesSeveralAddresses(t *testing.T) {
	t.Parallel()
	cfgPath := filepath.Join(t.TempDir(), "cfg.json")
	useConfig(t, cfgPath)
	out, err := runServeMCP(t, cfgPath, desktopDocker, false, "bind", "127.0.0.1", "172.17.0.1")
	if err != nil || !strings.HasPrefix(out, "bind: 127.0.0.1 172.17.0.1: ") || strings.Contains(out, "listen:") {
		t.Fatalf("bind two: %v %q", err, out)
	}
	data, err := os.ReadFile(cfgPath)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), `"bind": [`) {
		t.Errorf("two addresses are saved as a list: %s", data)
	}
	want := mcpBindView{Bind: "127.0.0.1", BindSet: true, Binds: []string{"127.0.0.1", "172.17.0.1"}, Listen: []string{"127.0.0.1", "172.17.0.1"}}
	if got := showMCPJSON(t, cfgPath, desktopDocker); !reflect.DeepEqual(got, want) {
		t.Errorf("--json: %+v, want %+v", got, want)
	}
	if _, err := runServeMCP(t, cfgPath, desktopDocker, false, "bind", "10.0.0.5", "0.0.0.0:4243"); err == nil || !strings.Contains(err.Error(), "names a port") {
		t.Errorf("an address with a port among several is refused: %v", err)
	}
	if _, err := runServeMCP(t, cfgPath, desktopDocker, false, "bind"); err == nil || !strings.Contains(err.Error(), "requires at least 1 arg") {
		t.Errorf("no address is refused: %v", err)
	}
	if got := showMCPJSON(t, cfgPath, desktopDocker); !reflect.DeepEqual(got.Binds, want.Binds) {
		t.Errorf("a refusal changes nothing: %v", got.Binds)
	}
}

// S-0732: auto, set or by default, resolves under a native Docker engine to
// 127.0.0.1 and the gateway of the network the container harness runs on,
// which show prints beside the addresses as written; under Docker Desktop to
// 127.0.0.1 alone; and, when docker cannot name its engine, to 127.0.0.1
// alone with a note saying what it could not say.
func TestServeMCPShowSaysWhatAutoResolvesTo(t *testing.T) {
	t.Parallel()
	cfgPath := filepath.Join(t.TempDir(), "cfg.json")
	useConfig(t, cfgPath)
	native := mcpDocker{engine: "Ubuntu 24.04 LTS"}
	out, err := runServeMCP(t, cfgPath, native, false, "show")
	if err != nil || !strings.Contains(out, "bind: auto (default): ") || !strings.Contains(out, "\nlisten: 127.0.0.1 172.17.0.1: ") || strings.Contains(out, "note:") {
		t.Errorf("unset, native engine: %v %q", err, out)
	}
	if out, err := runServeMCP(t, cfgPath, native, false, "bind", "auto", "10.0.0.5"); err != nil || !strings.HasPrefix(out, "bind: auto 10.0.0.5: ") || !strings.Contains(out, "\nlisten: 127.0.0.1 172.17.0.1 10.0.0.5: ") {
		t.Errorf("auto beside an address: %v %q", err, out)
	}
	want := mcpBindView{Bind: "auto", BindSet: true, Binds: []string{"auto", "10.0.0.5"}, Listen: []string{"127.0.0.1", "172.17.0.1", "10.0.0.5"}}
	if got := showMCPJSON(t, cfgPath, native); !reflect.DeepEqual(got, want) {
		t.Errorf("--json: %+v, want %+v", got, want)
	}
	if out, err := runServeMCP(t, cfgPath, desktopDocker, false, "bind", "auto"); err != nil || !strings.HasPrefix(out, "bind: auto: ") || !strings.Contains(out, "\nlisten: 127.0.0.1: ") {
		t.Errorf("auto under Docker Desktop: %v %q", err, out)
	}
	if data, _ := os.ReadFile(cfgPath); !strings.Contains(string(data), `"bind": "auto"`) {
		t.Errorf("one address is saved as a string, as before: %s", data)
	}
	out, err = runServeMCP(t, cfgPath, mcpDocker{}, false, "show")
	if err != nil || !strings.Contains(out, "\nlisten: 127.0.0.1: ") || !strings.Contains(out, "\nnote: auto listens on 127.0.0.1 alone: ") || !strings.Contains(out, "Cannot connect to the Docker daemon") {
		t.Errorf("docker cannot say: %v %q", err, out)
	}
	if got := showMCPJSON(t, cfgPath, mcpDocker{}); !strings.Contains(got.ListenError, "could not name the Docker engine") || !reflect.DeepEqual(got.Listen, []string{"127.0.0.1"}) {
		t.Errorf("--json, docker cannot say: %+v", got)
	}
}

// S-0732: mcpListen asks docker about the network the operator gave the
// container harness, and on a configuration it cannot load it is 127.0.0.1
// with the error.
func TestMCPListen(t *testing.T) {
	t.Parallel()
	cfgPath := filepath.Join(t.TempDir(), "cfg.json")
	useConfig(t, cfgPath)
	cfg, _, err := config.Load(cfgPath)
	if err != nil {
		t.Fatal(err)
	}
	cfg.Agent.Harnesses = map[string]config.HarnessHost{harness.Container: {Network: "flai-net"}}
	if err := config.Save(cfgPath, cfg); err != nil {
		t.Fatal(err)
	}
	var calls []string
	a := &app{configPath: cfgPath, errOut: &bytes.Buffer{}, runner: mcpDocker{engine: "Ubuntu 24.04 LTS", calls: &calls}}
	if got, err := a.mcpListen(); err != nil || !reflect.DeepEqual(got, []string{"127.0.0.1", "172.17.0.1"}) {
		t.Errorf("auto: %v %v", got, err)
	}
	if !slices.ContainsFunc(calls, func(c string) bool { return strings.HasPrefix(c, "network inspect flai-net ") }) {
		t.Errorf("the container harness's network is asked about: %q", calls)
	}
	if err := os.WriteFile(cfgPath, []byte("{not json"), 0o600); err != nil {
		t.Fatal(err)
	}
	if got, err := a.mcpListen(); err == nil || !reflect.DeepEqual(got, []string{config.DefaultMCPBind}) {
		t.Errorf("an unreadable configuration: %v %v", got, err)
	}
}
