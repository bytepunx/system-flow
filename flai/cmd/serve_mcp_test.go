package cmd

import (
	"bytes"
	"context"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// S-0096: flai serve starts a server for every project it serves, so they
// cannot all take the one default port: each takes the first free one from it,
// and keeps it once used.
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
	first, err := mcpAddrFor(repo)
	if err != nil || !inRange(first) {
		t.Fatalf("first free: %q %v", first, err)
	}
	// that port taken (another project's server): the next free one
	held, err := net.Listen("tcp", first)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = held.Close() }()
	second, err := mcpAddrFor(repo)
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
	if got, _ := mcpAddrFor(repo); got != "127.0.0.1:5999" {
		t.Errorf("remembered: %q", got)
	}
}

// flai mcp http --exit-with stops the server once that process is gone: a
// flai serve that was killed never stops its children itself.
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
