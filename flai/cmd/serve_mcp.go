package cmd

import (
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/bytepunx/system-flow/flai/internal/config"
	"github.com/bytepunx/system-flow/flai/internal/harness"
	"github.com/bytepunx/system-flow/flai/internal/serve"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// mcpPorts is how far past mcpDefaultAddr's port a project with no
// remembered address looks for a free one.
const mcpPorts = 20

// mcpAddrsFor is the addresses the project's server listens on, one on each
// of hosts, what the host's mcp.bind resolves to, all on one port: the port
// it used last, so an agent's configuration survives a restart, else the
// first port from mcpDefaultAddr's that is free on every host and not among
// taken, since flai host starts one for every project flai serve serves and
// they cannot all have the one default (S-0096, S-0106, S-0732). A
// remembered port that is among taken, given to another project already, is
// not given twice: two projects that remembered the same port both got it,
// and one failed to listen on every start (S-0110). A new port also leaves
// out avoid, the ports other projects remember while their servers are not
// running yet, so that one move does not push the next project off its port.
// Nor is a remembered port that cannot be listened on on every host: flai
// host asks only when the project's own server is not running, so another
// process holds it, and the server would fail to listen on it at every start
// (I-0271). taken and avoid hold addresses and are compared by port alone,
// whatever host each names, so that a change of mcp.bind moves no project
// off its port.
func mcpAddrsFor(repo *workitem.Repo, hosts []string, taken, avoid map[string]bool) ([]string, error) {
	takenPorts, avoidPorts := portsOf(taken), portsOf(avoid)
	if port := portOf(rememberedMCPAddr(repo)); port != "" && !takenPorts[port] && freeOnAll(hosts, port) {
		return onPort(hosts, port), nil
	}
	first, err := strconv.Atoi(portOf(mcpDefaultAddr))
	if err != nil {
		return nil, err
	}
	for p := first; p < first+mcpPorts; p++ {
		port := strconv.Itoa(p)
		if !takenPorts[port] && !avoidPorts[port] && freeOnAll(hosts, port) {
			return onPort(hosts, port), nil
		}
	}
	return nil, fmt.Errorf("no port between %d and %d is free to listen on at %s for flai mcp http; give the project one with flai mcp start --addr, or, if this machine has no such address, set another with flai serve mcp bind", first, first+mcpPorts-1, strings.Join(hosts, ", "))
}

// onPort is hosts' addresses, each on port.
func onPort(hosts []string, port string) []string {
	addrs := make([]string, len(hosts))
	for i, h := range hosts {
		addrs[i] = net.JoinHostPort(h, port)
	}
	return addrs
}

// freeOnAll says whether port can be listened on now on every one of hosts.
func freeOnAll(hosts []string, port string) bool {
	for _, addr := range onPort(hosts, port) {
		if !freePort(addr) {
			return false
		}
	}
	return true
}

// unlistenable is a host this machine cannot listen on at any port, and why.
type unlistenable struct {
	host string
	err  error
}

// listenableHosts is hosts split into those this machine can listen on, on
// a port the kernel picks, and those it cannot at all, such as the gateway
// of a Docker network whose bridge is down. It fails, naming them, when
// none can be.
func listenableHosts(hosts []string) ([]string, []unlistenable, error) {
	var kept []string
	var dropped []unlistenable
	for _, h := range hosts {
		if err := listenErr(net.JoinHostPort(h, "0")); err != nil {
			dropped = append(dropped, unlistenable{host: h, err: err})
			continue
		}
		kept = append(kept, h)
	}
	if len(kept) > 0 {
		return kept, dropped, nil
	}
	why := make([]string, len(dropped))
	for i, d := range dropped {
		why[i] = d.err.Error()
	}
	return nil, dropped, fmt.Errorf("flai mcp http cannot listen on %s on this machine (%s); set another address with flai serve mcp bind", strings.Join(hosts, ", "), strings.Join(why, "; "))
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

// mcpListen is the addresses the host's mcp.bind resolves to now, each auto
// expanded for the container harness's Docker program and network as the
// operator set them, for flai host to listen on at each server start. On a
// configuration it cannot load it is config.DefaultMCPBind alone, with the
// error; when docker cannot say what auto needs, it is the addresses
// serve.ResolveMCPBinds still gives, with what docker could not say.
func (a *app) mcpListen() ([]string, error) {
	cfg, _, err := a.loadConfig()
	if err != nil {
		return []string{config.DefaultMCPBind}, err
	}
	return a.mcpListenOf(cfg)
}

// mcpListenOf is mcpListen for cfg. Docker is asked from the working
// directory, as flai's other calls to docker are: the host's configuration
// belongs to no project.
func (a *app) mcpListenOf(cfg config.Config) ([]string, error) {
	h := a.harnessHost(cfg.Agent, harness.Container)
	return serve.ResolveMCPBinds(a.runner, "", h.Program, h.Network, cfg.MCPBinds())
}

// mcpBindView is mcp.bind as what shows the host's settings gives it: Bind,
// its first address as written, and BindSet, whether it is set, as they were
// while it held one address (S-0506); Binds, its addresses as written, or
// config.DefaultMCPBinds when unset; Listen, what they resolve to now for
// flai serve mcp show, or what the running server listens on for the
// settings, which ask docker nothing; and ListenError, what docker could
// not say for auto, when it could not.
type mcpBindView struct {
	Bind        string   `json:"bind"`
	BindSet     bool     `json:"bind_set"`
	Binds       []string `json:"binds"`
	Listen      []string `json:"listen,omitempty"`
	ListenError string   `json:"listen_error,omitempty"`
}

// mcpBindOf is cfg's mcp.bind as written, as mcpBindView gives it, without
// Listen: what shows every setting asks docker nothing.
func mcpBindOf(cfg config.Config) mcpBindView {
	binds := cfg.MCPBinds()
	return mcpBindView{Bind: binds[0], BindSet: len(cfg.MCP.Bind) > 0, Binds: binds}
}

// mcpBindResolved is mcpBindOf with Listen what its addresses resolve to
// now, for flai serve mcp show.
func (a *app) mcpBindResolved(cfg config.Config) mcpBindView {
	v := mcpBindOf(cfg)
	listen, err := a.mcpListenOf(cfg)
	v.Listen = nonNil(listen)
	if err != nil {
		v.ListenError = err.Error()
	}
	return v
}

// holdsAuto is whether v's addresses hold config.MCPBindAuto, so that what
// they resolve to now is worth saying.
func (v mcpBindView) holdsAuto() bool {
	return slices.Contains(v.Binds, config.MCPBindAuto)
}

// mcpDefaultBinds is config.DefaultMCPBinds as the help and the summaries
// write it.
var mcpDefaultBinds = strings.Join(config.DefaultMCPBinds(), " ")

// flai serve mcp: the host's say about the HTTP MCP server flai host starts
// for each project it serves (S-0506), on one address or several, or auto
// (S-0732).
func newServeMCPCmd(a *app) *cobra.Command {
	c := &cobra.Command{
		Use:   "mcp",
		Short: "The addresses flai host starts each project's HTTP MCP server on, on this host",
		Long: `flai host starts a flai mcp http for every project flai serve serves, each on
a port of its own that it keeps across restarts, listening on every address
mcp.bind holds, all on that one port.

mcp.bind holds one address or several, each an IP address or a host name
without a port, or ` + config.MCPBindAuto + `; unset is ` + mcpDefaultBinds + `. ` + config.MCPBindAuto + ` listens on ` + config.DefaultMCPBind + ` and, where the
host's Docker engine routes host.docker.internal to a network's gateway
rather than to loopback (a native engine, not Docker Desktop), also on the
gateway of the network the container harness runs on (flai serve agent
harness container --network), so that an agent in a container reaches the
server. flai resolves ` + config.MCPBindAuto + ` at each server start: flai host restart mcp
resolves it again, and flai serve mcp show says what it resolves to now.

flai serve mcp bind sets the addresses, such as 0.0.0.0, or Docker's bridge
gateway, 172.17.0.1, beside ` + config.DefaultMCPBind + `; flai serve mcp clear goes back to
` + mcpDefaultBinds + `. It is for every project on this host, and each server listens on
it from its next start: flai host restart mcp starts them again now.

Wherever it listens, the server refuses a request without the bearer token
flai mcp token prints and one with an Origin header, and it warns in its log
of an address beyond this machine that is not a container network's
gateway, since flai does not encrypt the token.`,
		Example: `  flai serve mcp bind auto
  flai serve mcp bind 127.0.0.1 172.17.0.1
  flai host restart mcp
  flai serve mcp show --json
  flai serve mcp clear`,
		Args: cobra.NoArgs,
		RunE: func(*cobra.Command, []string) error { return a.showMCPSettings() },
	}
	bind := &cobra.Command{
		Use:   "bind <address> [<address>...]",
		Short: "Start each project's MCP server on these addresses, each an IP address or a host name without a port, or " + config.MCPBindAuto + " (default " + mcpDefaultBinds + ")",
		Long: `Sets mcp.bind to the addresses given, each an IP address or a host name
without a port, or ` + config.MCPBindAuto + `, which listens on ` + config.DefaultMCPBind + ` and, under a native
Docker engine, on the gateway of the network the container harness runs on.
Every address listens on the project's one port, from each server's next
start: flai host restart mcp starts them again now. An address that is
refused, with a port or that is no IP address or host name, leaves mcp.bind
as it was, whichever of the addresses it is.`,
		Example: `  flai serve mcp bind auto
  flai serve mcp bind 127.0.0.1 172.17.0.1
  flai serve mcp bind ` + config.MCPBindAuto + ` 10.0.0.5`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			cfg, path, err := a.loadConfig()
			if err != nil {
				return err
			}
			if cfg, err = cfg.WithMCPBind(args...); err != nil {
				return err
			}
			if err := config.Save(path, cfg); err != nil {
				return err
			}
			return a.showMCPSettings()
		},
	}
	clear := &cobra.Command{
		Use:     "clear",
		Short:   "Unset mcp.bind, so that each project's MCP server listens on " + mcpDefaultBinds + " again, the default",
		Example: `  flai serve mcp clear`,
		Args:    cobra.NoArgs,
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
	show := &cobra.Command{
		Use:   "show",
		Short: "Print the addresses mcp.bind holds and, for " + config.MCPBindAuto + ", what it resolves to now",
		Example: `  flai serve mcp show
  flai serve mcp show --json`,
		Args: cobra.NoArgs,
		RunE: func(*cobra.Command, []string) error { return a.showMCPSettings() },
	}
	c.AddCommand(bind, clear, show)
	return c
}

// showMCPSettings prints mcp.bind as flai host uses it, whether it is set,
// and what it resolves to now.
func (a *app) showMCPSettings() error {
	cfg, _, err := a.loadConfig()
	if err != nil {
		return err
	}
	v := a.mcpBindResolved(cfg)
	if a.jsonOut {
		return a.printJSON(v)
	}
	writeMCPBind(a.out, v)
	return nil
}

// writeMCPBind writes the summary of v: the addresses as written, what they
// resolve to now when they hold auto, and what docker could not say.
func writeMCPBind(w io.Writer, v mcpBindView) {
	binds := strings.Join(v.Binds, " ")
	if v.BindSet {
		fmt.Fprintf(w, "bind: %s: flai host starts each project's MCP server on these addresses, on its one port, from its next start (flai host restart mcp); flai serve mcp clear goes back to %s\n", binds, mcpDefaultBinds)
	} else {
		fmt.Fprintf(w, "bind: %s (default): flai host starts each project's MCP server on what %s resolves to; flai serve mcp bind <address> [<address>...] sets others\n", binds, config.MCPBindAuto)
	}
	if v.holdsAuto() {
		fmt.Fprintf(w, "listen: %s: what they resolve to now; flai host resolves %s again at each server start\n", strings.Join(v.Listen, " "), config.MCPBindAuto)
	}
	if v.ListenError != "" {
		fmt.Fprintf(w, "note: %s\n", v.ListenError)
	}
}
