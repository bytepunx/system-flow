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

// mcpAddrFor is the address the project's server listens on, on bind, the
// host's mcp.bind: the port it used last, so an agent's configuration
// survives a restart, else the first free port from mcpDefaultAddr's that is
// not among taken, since flai host starts one for every project flai serve
// serves and they cannot all have the one default (S-0096, S-0106). A
// remembered port that is among taken, given to another project already, is
// not given twice: two projects that remembered the same port both got it,
// and one failed to listen on every start (S-0110). A new port also leaves
// out avoid, the ports other projects remember while their servers are not
// running yet, so that one move does not push the next project off its port.
// taken and avoid hold addresses and are compared by port alone, whatever
// host each names, so that a change of mcp.bind moves no project off its
// port.
func mcpAddrFor(repo *workitem.Repo, bind string, taken, avoid map[string]bool) (string, error) {
	takenPorts, avoidPorts := portsOf(taken), portsOf(avoid)
	if port := portOf(rememberedMCPAddr(repo)); port != "" && !takenPorts[port] {
		return net.JoinHostPort(bind, port), nil
	}
	first, err := strconv.Atoi(portOf(mcpDefaultAddr))
	if err != nil {
		return "", err
	}
	for p := first; p < first+mcpPorts; p++ {
		port := strconv.Itoa(p)
		addr := net.JoinHostPort(bind, port)
		if !takenPorts[port] && !avoidPorts[port] && freePort(addr) {
			return addr, nil
		}
	}
	return "", fmt.Errorf("no port between %d and %d is free to listen on at %s for flai mcp http; give the project one with flai mcp start --addr, or, if this machine has no such address, set another with flai serve mcp bind", first, first+mcpPorts-1, bind)
}

// portOf is the port of addr, a host and a port, or "" when it names none.
func portOf(addr string) string {
	_, port, err := net.SplitHostPort(addr)
	if err != nil {
		return ""
	}
	return port
}

// portsOf is the ports of addrs' addresses.
func portsOf(addrs map[string]bool) map[string]bool {
	ports := map[string]bool{}
	for addr, in := range addrs {
		if port := portOf(addr); in && port != "" {
			ports[port] = true
		}
	}
	return ports
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
