package cmd

import (
	"context"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/serve"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// mcpPorts is how far past mcpDefaultAddr's port a project with no
// remembered address looks for a free one.
const mcpPorts = 20

// serveMCP is how flai serve keeps each project's HTTP MCP server running
// (S-0096, ADR-0034): the same `flai mcp http` that `flai mcp start` runs,
// started with --exit-with so that it goes when this flai serve does, even
// one that was killed.
type serveMCP struct {
	a *app
	// one start at a time: a project with no remembered address takes the
	// first free port, and the one started just before must already hold its own
	mu sync.Mutex
}

// Ensure implements serve.MCP: one already running is used, else one is started.
func (m *serveMCP) Ensure(root string) (int, bool, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	repo, err := workitem.Open(root)
	if err != nil {
		return 0, false, err
	}
	if st, alive := readMCPState(repo, time.Now()); alive {
		return st.PID, false, nil
	}
	addr, err := mcpAddrFor(repo)
	if err != nil {
		return 0, false, err
	}
	args := []string{"mcp", "http", "--addr", addr, "--exit-with", strconv.Itoa(os.Getpid())}
	if m.a.configPath != "" {
		args = append(args, "--config", m.a.configPath)
	}
	st, err := m.a.startMCPHTTP(repo, args)
	if err != nil {
		return 0, false, err
	}
	return st.PID, true, nil
}

// Stop implements serve.MCP: SIGTERM, then a short wait for it to go.
func (m *serveMCP) Stop(pid int) {
	if err := serve.Terminate(pid); err != nil {
		return
	}
	deadline := time.Now().Add(mcpShutdownWait + 2*time.Second)
	for time.Now().Before(deadline) && serve.Alive(pid) {
		time.Sleep(100 * time.Millisecond)
	}
}

// mcpAddrFor is the address the project's server listens on: the one it used
// last, so an agent's configuration survives a restart, else the first free
// port from mcpDefaultAddr's, since flai serve starts one for every project it
// serves and they cannot all have the one default.
func mcpAddrFor(repo *workitem.Repo) (string, error) {
	if data, err := os.ReadFile(mcpFile(repo, mcpAddrFile)); err == nil && strings.TrimSpace(string(data)) != "" {
		return strings.TrimSpace(string(data)), nil
	}
	host, port, err := net.SplitHostPort(mcpDefaultAddr)
	if err != nil {
		return "", err
	}
	first, err := strconv.Atoi(port)
	if err != nil {
		return "", err
	}
	for p := first; p < first+mcpPorts; p++ {
		addr := net.JoinHostPort(host, strconv.Itoa(p))
		ln, err := net.Listen("tcp", addr)
		if err != nil {
			continue
		}
		_ = ln.Close()
		return addr, nil
	}
	return "", fmt.Errorf("no free port for flai mcp http between %d and %d; give the project one with flai mcp start --addr", first, first+mcpPorts-1)
}

// exitWithProcessEvery is how often --exit-with looks for its process.
var exitWithProcessEvery = 2 * time.Second

// exitWithProcess is ctx, ended once the process pid is gone.
func exitWithProcess(ctx context.Context, pid int) context.Context {
	ctx, cancel := context.WithCancel(ctx)
	go func() {
		defer cancel()
		tick := time.NewTicker(exitWithProcessEvery)
		defer tick.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-tick.C:
				if !serve.Alive(pid) {
					return
				}
			}
		}
	}()
	return ctx
}
