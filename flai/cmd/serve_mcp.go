package cmd

import (
	"context"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/bytepunx/system-flow/flai/internal/config"
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
// Nor is a remembered port that cannot be listened on: flai host asks only
// when the project's own server is not running, so another process holds it,
// and the server would fail to listen on it at every start (I-0271).
// taken and avoid hold addresses and are compared by port alone, whatever
// host each names, so that a change of mcp.bind moves no project off its
// port.
func mcpAddrFor(repo *workitem.Repo, bind string, taken, avoid map[string]bool) (string, error) {
	takenPorts, avoidPorts := portsOf(taken), portsOf(avoid)
	if port := portOf(rememberedMCPAddr(repo)); port != "" && !takenPorts[port] {
		if addr := net.JoinHostPort(bind, port); freePort(addr) {
			return addr, nil
		}
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

// mcpBind is the host's mcp.bind, the address flai host starts each
// project's MCP server on, read from the configuration at every look;
// config.DefaultMCPBind when it is unset or the configuration is not
// readable.
func (a *app) mcpBind() string {
	cfg, _, err := a.loadConfig()
	if err != nil {
		return config.DefaultMCPBind
	}
	return cfg.MCPBind()
}

// mcpBindOf is cfg's mcp.bind as flai host uses it, and whether it is set,
// for what shows the host's settings.
func mcpBindOf(cfg config.Config) (string, bool) {
	return cfg.MCPBind(), len(cfg.MCP.Bind) > 0
}

// flai serve mcp: the host's say about the HTTP MCP server flai host starts
// for each project it serves (S-0506).
func newServeMCPCmd(a *app) *cobra.Command {
	c := &cobra.Command{
		Use:   "mcp",
		Short: "The address flai host starts each project's HTTP MCP server on, on this host",
		Long: `flai host starts a flai mcp http for every project flai serve serves, each on
a port of its own that it keeps across restarts, and on mcp.bind, ` + config.DefaultMCPBind + ` by
default, so that only this machine reaches it. flai serve mcp bind sets
another address, an IP address or a host name without a port, such as
0.0.0.0 or Docker's bridge gateway, 172.17.0.1, so that an agent in a
container reaches the server through the host gateway; flai serve mcp clear
goes back to ` + config.DefaultMCPBind + `. It is for every project on this host, and each
server listens on it from its next start: flai host restart mcp starts them
again now.

Wherever it listens, the server refuses a request without the bearer token
flai mcp token prints and one with an Origin header, and it warns in its log
when it listens beyond this machine, since flai does not encrypt the token.`,
		Example: `  flai serve mcp bind 0.0.0.0
  flai host restart mcp
  flai serve mcp show --json
  flai serve mcp clear`,
		Args: cobra.NoArgs,
		RunE: func(*cobra.Command, []string) error { return a.showMCPSettings() },
	}
	bind := &cobra.Command{
		Use:   "bind <address>",
		Short: "Start each project's MCP server on this address, an IP address or a host name without a port (default " + config.DefaultMCPBind + ")",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			cfg, path, err := a.loadConfig()
			if err != nil {
				return err
			}
			if cfg, err = cfg.WithMCPBind(args[0]); err != nil {
				return err
			}
			if err := config.Save(path, cfg); err != nil {
				return err
			}
			return a.showMCPSettings()
		},
	}
	clear := &cobra.Command{
		Use:   "clear",
		Short: "Unset mcp.bind, so that each project's MCP server listens on " + config.DefaultMCPBind + " again",
		Args:  cobra.NoArgs,
		RunE: func(*cobra.Command, []string) error {
			cfg, path, err := a.loadConfig()
			if err != nil {
				return err
			}
			if err := config.Save(path, cfg.WithoutMCPBind()); err != nil {
				return err
			}
			return a.showMCPSettings()
		},
	}
	c.AddCommand(bind, clear,
		&cobra.Command{Use: "show", Short: "Print the address flai host starts each project's MCP server on", Args: cobra.NoArgs,
			RunE: func(*cobra.Command, []string) error { return a.showMCPSettings() }},
	)
	return c
}

// showMCPSettings prints mcp.bind as flai host uses it, and whether it is set.
func (a *app) showMCPSettings() error {
	cfg, _, err := a.loadConfig()
	if err != nil {
		return err
	}
	bind, set := mcpBindOf(cfg)
	if a.jsonOut {
		return a.printJSON(map[string]any{"bind": bind, "bind_set": set})
	}
	if set {
		fmt.Fprintf(a.out, "bind: %s: flai host starts each project's MCP server on this address, from its next start (flai host restart mcp); flai serve mcp clear goes back to %s\n", bind, config.DefaultMCPBind)
	} else {
		fmt.Fprintf(a.out, "bind: %s (default): flai host starts each project's MCP server on this machine only; flai serve mcp bind <address> sets another\n", bind)
	}
	return nil
}
