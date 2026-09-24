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

	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// S-0096, S-0106: flai host starts a server for every project flai serve
// serves, so they cannot all take the one default port: each takes the first
// free one from it that no other project was given, and keeps it once used.
func TestMCPAddrForAProjectFlaiServeStarts(t *testing.T) {
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
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
	first, err := mcpAddrFor(repo, nil, nil)
	if err != nil || !inRange(first) {
		t.Fatalf("first free: %q %v", first, err)
	}
	// given to another project and not yet listening: not given again
	if other, err := mcpAddrFor(repo, map[string]bool{first: true}, nil); err != nil || other == first || !inRange(other) {
		t.Errorf("with %s given to another: %q %v", first, other, err)
	}
	// that port taken (another project's server): the next free one
	held, err := net.Listen("tcp", first)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = held.Close() }()
	second, err := mcpAddrFor(repo, nil, nil)
	if err != nil || !inRange(second) || second == first {
		t.Errorf("with %s taken: %q %v", first, second, err)
	}
	// once a project has used an address, it keeps it
	if err := os.MkdirAll(repo.CacheDir(), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(mcpFile(repo, mcpAddrFile), []byte("127.0.0.1:5999\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if got, _ := mcpAddrFor(repo, nil, nil); got != "127.0.0.1:5999" {
		t.Errorf("remembered: %q", got)
	}
	// S-0110: a remembered address another project already has is not given twice
	moved, err := mcpAddrFor(repo, map[string]bool{"127.0.0.1:5999": true}, nil)
	if err != nil || moved == "127.0.0.1:5999" || !inRange(moved) {
		t.Errorf("remembered, and given to another: %q %v", moved, err)
	}
	// and the new one is none that another project remembers
	if again, err := mcpAddrFor(repo, map[string]bool{"127.0.0.1:5999": true}, map[string]bool{moved: true}); err != nil || again == moved || again == "127.0.0.1:5999" {
		t.Errorf("remembered by another: %q %v", again, err)
	}
}

// flai mcp http --exit-with stops the server once that process is gone: a
// flai host that was killed never stops its children itself.
func TestMCPHTTPExitsWithTheProcessItWasGiven(t *testing.T) {
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	t.Setenv("FLAI_AGENT", "tester")
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
	a := &app{out: &out, errOut: &errOut, cwd: root}
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
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	var logged bytes.Buffer
	l := &hostLauncher{a: &app{log: slog.New(slog.NewTextHandler(&logged, nil))}, taken: map[string]string{}}
	var repos []*workitem.Repo
	for i := 0; i < 2; i++ {
		repo, err := workitem.Open(tempProject(t))
		if err != nil {
			t.Fatal(err)
		}
		if err := os.MkdirAll(repo.CacheDir(), 0o700); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(mcpFile(repo, mcpAddrFile), []byte("127.0.0.1:4243\n"), 0o600); err != nil {
			t.Fatal(err)
		}
		repos = append(repos, repo)
	}
	// a third project flai serve serves, not running, remembers 4244: a move does not take it
	third := tempProject(t)
	thirdRepo, _ := workitem.Open(third)
	_ = os.MkdirAll(thirdRepo.CacheDir(), 0o700)
	_ = os.WriteFile(mcpFile(thirdRepo, mcpAddrFile), []byte("127.0.0.1:4244\n"), 0o600)
	serveDir := filepath.Join(filepath.Dir(os.Getenv("FLAI_CONFIG")), "serve")
	_ = os.MkdirAll(serveDir, 0o700)
	_ = os.WriteFile(filepath.Join(serveDir, "projects.json"), []byte(`[{"key":"c","name":"c","root":"`+third+`","url":"http://127.0.0.1:1"}]`), 0o600)
	first, err := l.mcpAddr(repos[0])
	if err != nil || first != "127.0.0.1:4243" {
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
