package cmd

import (
	"context"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/serve"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// mcpPorts is how far past mcpDefaultAddr's port a project with no
// remembered address looks for a free one.
const mcpPorts = 20

// mcpAddrFor is the address the project's server listens on: the one it used
// last, so an agent's configuration survives a restart, else the first free
// port from mcpDefaultAddr's that is not among taken, since flai host starts
// one for every project flai serve serves and they cannot all have the one
// default (S-0096, S-0106). A remembered address that is among taken, given
// to another project already, is not given twice: two projects that
// remembered the same port both got it, and one failed to listen on every
// start (S-0110). A new port also leaves out avoid, the ports other projects
// remember while their servers are not running yet, so that one move does
// not push the next project off its port.
func mcpAddrFor(repo *workitem.Repo, taken, avoid map[string]bool) (string, error) {
	if addr := rememberedMCPAddr(repo); addr != "" && !taken[addr] {
		return addr, nil
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
		if !taken[addr] && !avoid[addr] && freePort(addr) {
			return addr, nil
		}
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

// rememberedMCPAddr is the address the project's server last listened on,
// or "" when it has none.
func rememberedMCPAddr(repo *workitem.Repo) string {
	data, err := os.ReadFile(mcpFile(repo, mcpAddrFile))
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(data))
}
