package cmd

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/cobra"

	"github.com/bytepunx/system-flow/flai/internal/buildinfo"
	"github.com/bytepunx/system-flow/flai/internal/config"
	"github.com/bytepunx/system-flow/flai/internal/guard"
	"github.com/bytepunx/system-flow/flai/internal/harness"
	"github.com/bytepunx/system-flow/flai/internal/mcphttp"
	"github.com/bytepunx/system-flow/flai/internal/mcpserver"
	"github.com/bytepunx/system-flow/flai/internal/serve"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// flai mcp over HTTP (ADR-0030): a process of its own on the host, one per
// project, for agents that cannot start flai mcp themselves. Its token,
// state, and log are files under the main checkout's .flai-cache.

const (
	mcpDefaultAddr  = "127.0.0.1:4243"
	mcpPath         = "/mcp"
	mcpTokenFile    = serve.MCPTokenFile
	mcpStateFile    = serve.MCPStateFile
	mcpAddrFile     = "mcp-http.addr"
	mcpLogFile      = "mcp-http.log"
	mcpStateEvery   = 3 * time.Second
	mcpStaleAfter   = serve.MCPStaleAfter
	mcpShutdownWait = 5 * time.Second
)

// mcpState is what a running server says about itself, rewritten while it
// runs so that a file left by a killed process is known to be stale. Addr
// and URL name its loopback listener when it has one, else its first, so
// that a reader of one address finds it as before; Addrs is every address
// it listens on, as bound.
type mcpState struct {
	PID     int      `json:"pid"`
	Addr    string   `json:"addr"`
	URL     string   `json:"url"`
	Addrs   []string `json:"addrs,omitempty"`
	Started string   `json:"started"`
	Updated string   `json:"updated"`
	Version string   `json:"version"`
	// ExitWith is the process it goes with (--exit-with), flai serve's when
	// flai serve started it: a restart with a new token keeps it (S-0105).
	ExitWith int `json:"exit_with,omitempty"`
}

func mcpFile(repo *workitem.Repo, name string) string { return filepath.Join(repo.CacheDir(), name) }

func readMCPState(repo *workitem.Repo, now time.Time) (mcpState, bool) {
	var st mcpState
	data, err := os.ReadFile(mcpFile(repo, mcpStateFile))
	if err != nil || json.Unmarshal(data, &st) != nil {
		return mcpState{}, false
	}
	updated, err := time.Parse(time.RFC3339, st.Updated)
	if err != nil || now.Sub(updated) > mcpStaleAfter {
		return st, false
	}
	return st, serve.Alive(st.PID)
}

func writeMCPState(repo *workitem.Repo, st mcpState) error {
	data, err := json.MarshalIndent(st, "", "  ")
	if err != nil {
		return err
	}
	path := mcpFile(repo, mcpStateFile)
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, append(data, '\n'), 0o600); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}

// mcpToken returns the server's token, creating it when missing or asked to.
func mcpToken(repo *workitem.Repo, rotate bool) (token string, created bool, err error) {
	path := mcpFile(repo, mcpTokenFile)
	if !rotate {
		if data, err := os.ReadFile(path); err == nil && strings.TrimSpace(string(data)) != "" {
			return strings.TrimSpace(string(data)), false, nil
		}
	}
	buf := make([]byte, tokenBytes)
	if _, err := rand.Read(buf); err != nil {
		return "", false, err
	}
	token = base64.RawURLEncoding.EncodeToString(buf)
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return "", false, err
	}
	if err := os.WriteFile(path, []byte(token+"\n"), tokenFileMode); err != nil {
		return "", false, err
	}
	return token, true, os.Chmod(path, tokenFileMode)
}

func loopback(addr string) bool {
	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		return false
	}
	if host == "localhost" {
		return true
	}
	ip := net.ParseIP(host)
	return ip != nil && ip.IsLoopback()
}

// mcpURL is where a client on this machine reaches a listener.
func mcpURL(addr net.Addr) string {
	host, port, err := net.SplitHostPort(addr.String())
	if err != nil {
		return "http://" + addr.String() + mcpPath
	}
	if ip := net.ParseIP(host); ip != nil && ip.IsUnspecified() {
		host = "127.0.0.1"
	}
	return "http://" + net.JoinHostPort(host, port) + mcpPath
}

// samePort refuses addresses that are not host:port or whose ports differ:
// the server listens on one port, at every address it is given.
func samePort(addrs []string) error {
	port := ""
	for _, addr := range addrs {
		_, p, err := net.SplitHostPort(addr)
		if err != nil {
			return fmt.Errorf("cannot listen on %s: %w; give each --addr as host:port, such as %s", addr, err, mcpDefaultAddr)
		}
		if port != "" && p != port {
			return fmt.Errorf("the addresses %s name different ports (%s and %s): the server listens on one port at every address, so give each --addr the same port", strings.Join(addrs, ", "), port, p)
		}
		port = p
	}
	return nil
}

// primaryAddr is the address mcpState's Addr names among addrs: the first
// loopback one, else the first.
func primaryAddr(addrs []string) int {
	if i := slices.IndexFunc(addrs, loopback); i >= 0 {
		return i
	}
	return 0
}

// listenMCP listens on every address, on one port: with port 0 the first
// listener picks it and the others listen on the same. A failure on any
// address closes those opened and names the address.
func listenMCP(addrs []string) ([]net.Listener, error) {
	if len(addrs) == 0 {
		return nil, errors.New("no address to listen on; give one with --addr")
	}
	if err := samePort(addrs); err != nil {
		return nil, err
	}
	var lns []net.Listener
	port := ""
	for _, addr := range addrs {
		if host, p, _ := net.SplitHostPort(addr); port != "" && p == "0" {
			addr = net.JoinHostPort(host, port)
		}
		ln, err := net.Listen("tcp", addr)
		if err != nil {
			closeListeners(lns)
			return nil, fmt.Errorf("cannot listen on %s: %w; another project's server may hold it, so choose another with --addr", addr, err)
		}
		lns = append(lns, ln)
		if port == "" {
			_, port, _ = net.SplitHostPort(ln.Addr().String())
		}
	}
	return lns, nil
}

func closeListeners(lns []net.Listener) {
	for _, ln := range lns {
		_ = ln.Close()
	}
}

// logReach says, for each address that is not loopback, whether it reaches
// beyond this machine, warned about, or is a Docker network's gateway, where
// only the host and its containers reach it. Docker is asked for the
// gateways, from dir, only when some address is neither loopback nor
// unspecified; when it cannot answer, every such address is warned about.
func (a *app) logReach(log *slog.Logger, dir string, addrs []string) {
	var gateways []string
	for _, addr := range addrs {
		if !loopback(addr) && !unspecified(addr) {
			gateways = a.containerGateways(log, dir)
			break
		}
	}
	for _, addr := range addrs {
		switch {
		case loopback(addr):
		case !unspecified(addr) && isGateway(addr, gateways):
			log.Info("listening on a container network's gateway: the host and its containers reach it", "addr", addr)
		default:
			log.Warn("listening beyond this machine: the token travels in every request and flai does not encrypt it, so put a proxy or tunnel that terminates TLS in front", "addr", addr)
		}
	}
}

// containerGateways is every Docker network's gateway, as the container
// harness's docker client names them; none, logged at debug, when it cannot.
func (a *app) containerGateways(log *slog.Logger, dir string) []string {
	program := a.harnessHost(config.AgentStart{}, harness.Container).Program
	if cfg, _, err := a.loadConfig(); err == nil {
		program = a.harnessHost(cfg.Agent, harness.Container).Program
	} else {
		log.Debug("could not read the config for the container harness's docker client", "program", program, "err", err.Error())
	}
	gateways, err := serve.ContainerGateways(a.runner, dir, program)
	if err != nil {
		log.Debug("could not ask docker for its networks' gateways", "program", program, "err", err.Error())
	}
	return gateways
}

func unspecified(addr string) bool {
	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		return false
	}
	ip := net.ParseIP(host)
	return host == "" || ip != nil && ip.IsUnspecified()
}

func isGateway(addr string, gateways []string) bool {
	host, _, err := net.SplitHostPort(addr)
	ip := net.ParseIP(host)
	if err != nil || ip == nil {
		return false
	}
	return slices.ContainsFunc(gateways, func(g string) bool { return ip.Equal(net.ParseIP(g)) })
}

// mcpHTTPOptions are the options of the server an agent reaches over HTTP:
// permission_prompt holds a write at most mcpserver.HTTPPermissionHold, below
// the 60 seconds after which Claude Code drops a call there (ADR-0215), and
// wait_for_events and wait_for_work hold a call at most
// mcpserver.HTTPLongestWait, the same bound (ADR-0244). It
// allows a protected write without asking only under the auto-approve host
// action and only for the agent flai serve runs for the story (ADR-0243).
func (a *app) mcpHTTPOptions(repo *workitem.Repo, agent string, closing <-chan struct{}, installed string) mcpserver.Options {
	return mcpserver.Options{Repo: repo, Agent: agent, Version: buildinfo.Version, Now: a.now, Runner: a.runner, Closing: closing, Logger: a.logger(), Agents: a.mcpAgents, Plans: a.mcpPlan, Analyses: a.mcpAnalyze, Curations: a.mcpCurate, Roles: a.mcpRoleRun, Activities: a.mcpActivity, AutoApprove: a.mcpAutoApprove, ServedAgent: a.mcpServedAgent, Publish: a.mcpPublisher(), Accept: a.mcpAccepter(), RelativePaths: a.relativeWorktrees, Installed: installed, PermissionHold: mcpserver.HTTPPermissionHold, MaxWait: mcpserver.HTTPLongestWait, Budgets: a.mcpBudgets(), HostAction: a.hostActionOn}.WithSession(mcpserver.SessionFromEnv())
}

// mcpServedAgent reports whether agent is the name of the run flai serve
// records, in its folder beside the configuration, as started for story in
// the project at root and still running: the process it started is alive and
// is the one it started (ADR-0243). A state it cannot read names no run.
func (a *app) mcpServedAgent(root, agent, story string) bool {
	run := a.serveDir().AgentStates()[root].Stories[story]
	return run != nil && run.Agent == agent && run.Ended == "" && run.Error == "" && run.PID > 0 && serve.Owns(run.PID, run.Start)
}

// mcpServedOrchestrator reports whether agent is the name of the
// orchestrator's run flai serve records for the project at root and that is
// still running: the process it started is alive and is the one it started.
// The server over HTTP is a child of flai host, whose FLAI_ROLE says nothing
// of its callers, so it knows the orchestrator by this alone (ADR-0278): a
// name is only a claim, which every holder of the token can make. A state
// it cannot read names no run.
func (a *app) mcpServedOrchestrator(root, agent string) bool {
	run := a.serveDir().AgentStates()[root].Orchestrator
	return agent != "" && run != nil && run.Agent == agent && run.Ended == "" && run.Error == "" && run.PID > 0 && serve.Owns(run.PID, run.Start)
}

// mcpHTTPRole is the role of the caller named agent of the server over HTTP
// for the project at root, read when its session starts: the orchestrator's
// for the orchestrator's run flai serve records, else none, which keeps the
// role SessionFromEnv read from FLAI_ROLE (ADR-0278).
func (a *app) mcpHTTPRole(root, agent string) string {
	if a.mcpServedOrchestrator(root, agent) {
		return guard.RoleOrchestrate
	}
	return ""
}

// serveMCPHTTP runs the server on every address in addrs, one port, until
// ctx ends.
func (a *app) serveMCPHTTP(ctx context.Context, repo *workitem.Repo, addrs []string, maxSessions int, idle time.Duration, exitWith int) error {
	if st, alive := readMCPState(repo, time.Now()); alive {
		return fmt.Errorf("flai mcp is already serving this project at %s (pid %d); flai mcp stop ends it", st.URL, st.PID)
	}
	token, _, err := mcpToken(repo, false)
	if err != nil {
		return err
	}
	log := a.logger().With("component", "mcp")
	lns, err := listenMCP(addrs)
	if err != nil {
		return err
	}
	bound := make([]string, len(lns))
	for i, ln := range lns {
		bound[i] = ln.Addr().String()
	}
	a.logReach(log, repo.MainRoot, bound)
	primary := lns[primaryAddr(bound)]
	closing := make(chan struct{})
	installed := a.mcpInstalled()
	handler := mcphttp.Handler(mcphttp.Options{
		Token: token, ProjectKey: repo.Manifest.Key, ProjectName: repo.Manifest.Name,
		MaxSessions: maxSessions, Idle: idle, Logger: log,
		NewServer: func(agent string) *mcp.Server {
			opt := a.mcpHTTPOptions(repo, agent, closing, installed)
			if role := a.mcpHTTPRole(mainRootOf(repo), agent); role != "" {
				opt.Role = role
			}
			return mcpserver.New(opt)
		},
	})
	mux := http.NewServeMux()
	mux.Handle(mcpPath, handler)
	// No write timeout: a call such as task_done may run past any bound, and
	// the SDK pauses a session's --idle timer while one of its requests is
	// open; wait_for_events and wait_for_work hold one for up to
	// mcpserver.HTTPLongestWait (ADR-0244).
	srv := &http.Server{Handler: mux, ReadHeaderTimeout: 10 * time.Second}

	st := mcpState{PID: os.Getpid(), Addr: primary.Addr().String(), URL: mcpURL(primary.Addr()), Addrs: bound, Started: time.Now().UTC().Format(time.RFC3339), Version: buildinfo.Version, ExitWith: exitWith}
	beat := func() error {
		st.Updated = time.Now().UTC().Format(time.RFC3339)
		return writeMCPState(repo, st)
	}
	if err := beat(); err != nil {
		closeListeners(lns)
		return err
	}
	defer func() { _ = os.Remove(mcpFile(repo, mcpStateFile)) }()
	log.Info("mcp server started", "url", st.URL, "addrs", bound, "root", repo.Root, "project", repo.Manifest.Key)
	fmt.Fprintf(a.errOut, "flai mcp serving %s at %s\n  token: %s (flai mcp token prints it)\n", repo.Manifest.Name, st.URL, mcpFile(repo, mcpTokenFile))

	failed := make(chan error, len(lns))
	for _, ln := range lns {
		go func() { failed <- srv.Serve(ln) }()
	}
	tick := time.NewTicker(mcpStateEvery)
	defer tick.Stop()
	for {
		select {
		case err := <-failed:
			_ = srv.Close() // the other listeners with it
			return err
		case <-tick.C:
			if err := beat(); err != nil {
				log.Warn("could not write state", "error", err.Error())
			}
		case <-ctx.Done():
			// An idle agent holds wait_for_events for minutes: those calls are
			// ended as if their time had passed, and anything else gets a moment.
			close(closing)
			wait, cancel := context.WithTimeout(context.Background(), mcpShutdownWait)
			defer cancel()
			if err := srv.Shutdown(wait); err != nil {
				_ = srv.Close()
			}
			log.Info("mcp server stopped")
			return nil
		}
	}
}

// startMCPHTTP starts flai mcp http detached and waits for it to listen.
func (a *app) startMCPHTTP(repo *workitem.Repo, args []string) (mcpState, error) {
	exe, err := os.Executable()
	if err != nil {
		return mcpState{}, err
	}
	if err := os.MkdirAll(repo.CacheDir(), 0o700); err != nil {
		return mcpState{}, err
	}
	logPath := mcpFile(repo, mcpLogFile)
	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o600)
	if err != nil {
		return mcpState{}, err
	}
	defer func() { _ = logFile.Close() }()
	cmd := exec.CommandContext(context.WithoutCancel(context.Background()), exe, args...)
	cmd.Dir = repo.MainRoot
	cmd.Stdout, cmd.Stderr = logFile, logFile
	cmd.Env = append(os.Environ(), "LOG_FORMAT=json")
	serve.Detach(cmd)
	if err := cmd.Start(); err != nil {
		return mcpState{}, err
	}
	pid := cmd.Process.Pid
	exited := make(chan struct{})
	go func() { _ = cmd.Wait(); close(exited) }()
	deadline := time.After(5 * time.Second)
	for {
		if st, alive := readMCPState(repo, time.Now()); alive && st.PID == pid {
			return st, nil
		}
		select {
		case <-exited:
			return mcpState{}, fmt.Errorf("flai mcp http stopped at once; see %s", logPath)
		case <-deadline:
			return mcpState{}, fmt.Errorf("flai mcp http did not come up; see %s", logPath)
		case <-time.After(100 * time.Millisecond):
		}
	}
}

func (a *app) stopMCPHTTP(repo *workitem.Repo) (mcpState, bool, error) {
	st, alive := readMCPState(repo, time.Now())
	if !alive {
		_ = os.Remove(mcpFile(repo, mcpStateFile)) // what a killed process left
		return st, false, nil
	}
	if err := serve.Terminate(st.PID); err != nil {
		return st, false, err
	}
	deadline := time.Now().Add(mcpShutdownWait + 2*time.Second)
	for time.Now().Before(deadline) {
		if !serve.Alive(st.PID) {
			break
		}
		if _, err := os.Stat(mcpFile(repo, mcpStateFile)); errors.Is(err, os.ErrNotExist) {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}
	return st, true, nil
}

// mcpClientConfig is the sample client configuration flai mcp status prints,
// naming the per-server timeout flai's own MCP file names (ADR-0258).
func mcpClientConfig(url string) string {
	return fmt.Sprintf(`{ "mcpServers": { "flai": { "type": "http", "url": %q,
    "headers": { "Authorization": "Bearer <flai mcp token>", "X-Flai-Agent": "<the agent's name>" },
    "timeout": %d } } }`, url, mcphttp.ClientTimeout.Milliseconds())
}

func newMCPHTTPCmds(a *app) []*cobra.Command {
	var addrs []string
	var maxSessions int
	var idle time.Duration
	var exitWith int
	flags := func(c *cobra.Command) {
		c.Flags().StringArrayVar(&addrs, "addr", nil, "address to listen on, repeatable, every one on the same port (default the last one used here, else "+mcpDefaultAddr+")")
		c.Flags().IntVar(&maxSessions, "max-sessions", 16, "MCP sessions at once")
		c.Flags().DurationVar(&idle, "idle", 30*time.Minute, "a session with no request for this long ends")
	}
	// The address sticks, so an agent's configuration keeps working across
	// restarts: the one mcpState's Addr names, when several are given.
	resolve := func(repo *workitem.Repo) ([]string, error) {
		if len(addrs) > 0 {
			return addrs, samePort(addrs)
		}
		if data, err := os.ReadFile(mcpFile(repo, mcpAddrFile)); err == nil && strings.TrimSpace(string(data)) != "" {
			return []string{strings.TrimSpace(string(data))}, nil
		}
		return []string{mcpDefaultAddr}, nil
	}
	remember := func(repo *workitem.Repo, addr string) {
		_ = os.MkdirAll(repo.CacheDir(), 0o700)
		_ = os.WriteFile(mcpFile(repo, mcpAddrFile), []byte(addr+"\n"), 0o600)
	}

	httpCmd := &cobra.Command{
		Use:   "http",
		Short: "Serve MCP over Streamable HTTP in the foreground, on this machine only unless --addr says otherwise",
		Long: `The same server as flai mcp, over HTTP at /mcp, for agents that cannot start
a process here (ADR-0030). It listens on ` + mcpDefaultAddr + ` unless --addr says
otherwise. --addr may be given more than once, every address on the same port,
such as loopback for clients here and a Docker network's gateway for agents in
containers; the loopback one (else the first) is remembered for this project
and named as the server's URL. Every request
carries the token in .flai-cache/mcp.token as a bearer header; a request from
a browser (one with an Origin header) is refused. The token is the gate,
whatever Host a request names, so an agent in a container may dial
host.docker.internal (ADR-0216). The agent is named by the X-Flai-Agent
header, else by the client's own name.

An address beyond this machine is allowed and warned about: the token travels
in every request, and TLS is a proxy's or a tunnel's job, not flai's. A Docker
network's gateway, which only the host and its containers reach, is not.`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			repo, err := a.httpMCPProject()
			if err != nil {
				return err
			}
			ctx, stop := signal.NotifyContext(cmd.Context(), os.Interrupt, syscall.SIGTERM)
			defer stop()
			if exitWith > 0 {
				ctx = exitWithProcess(ctx, exitWith)
			}
			listen, err := resolve(repo)
			if err != nil {
				return err
			}
			remember(repo, listen[primaryAddr(listen)])
			return a.serveMCPHTTP(ctx, repo, listen, maxSessions, idle, exitWith)
		},
	}
	flags(httpCmd)
	// flai serve starts it with its own PID (S-0096): a flai serve that was
	// killed never stops its children, so each goes by itself.
	httpCmd.Flags().IntVar(&exitWith, "exit-with", 0, "stop when the process with this PID is gone")
	_ = httpCmd.Flags().MarkHidden("exit-with")

	start := &cobra.Command{
		Use:   "start",
		Short: "Start the HTTP server in the background, if it is not running",
		Args:  cobra.NoArgs,
		RunE: func(*cobra.Command, []string) error {
			repo, err := a.httpMCPProject()
			if err != nil {
				return err
			}
			st, alive := readMCPState(repo, time.Now())
			started := false
			if !alive {
				listen, err := resolve(repo)
				if err != nil {
					return err
				}
				args := append(addrArgs(listen), "--max-sessions", fmt.Sprint(maxSessions), "--idle", idle.String())
				if a.configPath != "" {
					args = append(args, "--config", a.configPath)
				}
				if st, err = a.startMCPHTTP(repo, args); err != nil {
					return err
				}
				started = true
			}
			if a.jsonOut {
				return a.printJSON(map[string]any{"started": started, "state": st, "log": mcpFile(repo, mcpLogFile), "token_file": mcpFile(repo, mcpTokenFile)})
			}
			if started {
				fmt.Fprintf(a.out, "flai mcp started at %s (pid %d)\n", st.URL, st.PID)
			} else {
				fmt.Fprintf(a.out, "flai mcp is already running at %s (pid %d, since %s)\n", st.URL, st.PID, st.Started)
			}
			fmt.Fprintf(a.out, "  token: %s (flai mcp token prints it)\n  log: %s\n  stop with: flai mcp stop\n", mcpFile(repo, mcpTokenFile), mcpFile(repo, mcpLogFile))
			return nil
		},
	}
	flags(start)

	stop := &cobra.Command{
		Use:   "stop",
		Short: "Stop the HTTP server; every session ends",
		Args:  cobra.NoArgs,
		RunE: func(*cobra.Command, []string) error {
			repo, err := a.httpMCPProject()
			if err != nil {
				return err
			}
			st, stopped, err := a.stopMCPHTTP(repo)
			if err != nil {
				return err
			}
			if a.jsonOut {
				return a.printJSON(map[string]any{"stopped": stopped})
			}
			if stopped {
				fmt.Fprintf(a.out, "flai mcp stopped (pid %d)\n", st.PID)
			} else {
				fmt.Fprintln(a.out, "flai mcp is not running over HTTP here")
			}
			return nil
		},
	}

	status := &cobra.Command{
		Use:   "status",
		Short: "Whether the HTTP server runs for this project, where, and how an agent connects",
		Args:  cobra.NoArgs,
		RunE: func(*cobra.Command, []string) error {
			repo, err := a.httpMCPProject()
			if err != nil {
				return err
			}
			st, alive := readMCPState(repo, time.Now())
			if a.jsonOut {
				out := map[string]any{"running": alive, "token_file": mcpFile(repo, mcpTokenFile), "log": mcpFile(repo, mcpLogFile)}
				if alive {
					out["state"] = st
				}
				return a.printJSON(out)
			}
			if !alive {
				fmt.Fprintln(a.out, "flai mcp is not running over HTTP here (flai mcp start); over stdio it needs no server: .mcp.json starts it")
				return nil
			}
			fmt.Fprintf(a.out, "flai mcp %s serving %s at %s since %s (pid %d)\n  token: %s (flai mcp token prints it)\n  log: %s\n  an agent's configuration:\n%s\n",
				st.Version, repo.Manifest.Name, st.URL, st.Started, st.PID, mcpFile(repo, mcpTokenFile), mcpFile(repo, mcpLogFile), mcpClientConfig(st.URL))
			return nil
		},
	}

	var rotate bool
	token := &cobra.Command{
		Use:   "token",
		Short: "Print the HTTP server's token; --rotate replaces it and restarts a running server",
		Long: `The token lives at .flai-cache/mcp.token (mode 0600, git-ignored) and is
created the first time it is needed. It is not the dashboard's token: that
one opens the dashboard, this one lets an agent act on the project. --rotate
writes a new one and restarts a running server, which ends every session.

In a run flai serve started for a harness that reaches the server over HTTP,
such as one in a container, the token is handed to the run as ` + serve.MCPTokenEnv + `,
with the server's URL as ` + serve.MCPURLEnv + ` (S-0506): while ` + serve.MCPTokenEnv + ` is set,
flai mcp token prints it, the token in force, and reads no file, so that it
answers where there is no .flai-cache. --rotate ignores ` + serve.MCPTokenEnv + ` and
rotates the project's file.`,
		Args: cobra.NoArgs,
		RunE: func(*cobra.Command, []string) error {
			if tok := strings.TrimSpace(os.Getenv(serve.MCPTokenEnv)); tok != "" && !rotate {
				if a.jsonOut {
					return a.printJSON(map[string]any{"token": tok, "created": false, "rotated": false, "restarted": false, "env": serve.MCPTokenEnv})
				}
				fmt.Fprintln(a.out, tok)
				return nil
			}
			repo, err := a.httpMCPProject()
			if err != nil {
				return err
			}
			tok, created, err := mcpToken(repo, rotate)
			if err != nil {
				return err
			}
			restarted := false
			if rotate {
				if st, stopped, err := a.stopMCPHTTP(repo); err != nil {
					return err
				} else if stopped && a.hostKeeps(st) {
					// the host that started it starts it again, on the new token (S-0106)
					if err := a.awaitMCPRestart(repo, st.PID); err != nil {
						return err
					}
					restarted = true
				} else if stopped {
					if _, err := a.startMCPHTTP(repo, restartArgs(st, a.configPath)); err != nil {
						return err
					}
					restarted = true
				}
			}
			if a.jsonOut {
				return a.printJSON(map[string]any{"token": tok, "created": created, "rotated": rotate, "restarted": restarted, "file": mcpFile(repo, mcpTokenFile)})
			}
			fmt.Fprintln(a.out, tok)
			if restarted {
				fmt.Fprintln(a.errOut, "the running server was restarted with the new token; every session ended")
			}
			return nil
		},
	}
	token.Flags().BoolVar(&rotate, "rotate", false, "replace the token in the project's file, whatever "+serve.MCPTokenEnv+" holds")
	return []*cobra.Command{httpCmd, start, stop, status, token}
}

// httpMCPProject is the project an HTTP MCP server serves: one, always
// (ADR-0030). Outside any project it says what to run instead (S-0101).
func (a *app) httpMCPProject() (*workitem.Repo, error) {
	repo, err := a.projectOrNone()
	if err == nil && repo == nil {
		return nil, errors.New("flai mcp over HTTP serves one project, and there is none here or above; run it in a project (flai serve keeps one running for each project it serves), or run flai mcp on stdio here, which serves every project below this folder")
	}
	return repo, err
}

// hostKeeps says whether the server st was is one flai host keeps: it goes
// with the host that runs now, which starts it again once it has stopped.
func (a *app) hostKeeps(st mcpState) bool {
	h, alive := a.hostDir().ReadStatus(time.Now())
	return alive && st.ExitWith == h.PID
}

// awaitMCPRestart waits for the host to start the project's server again
// after the one with pid stopped.
func (a *app) awaitMCPRestart(repo *workitem.Repo, pid int) error {
	deadline := time.Now().Add(15 * time.Second)
	for time.Now().Before(deadline) {
		if st, alive := readMCPState(repo, time.Now()); alive && st.PID != pid {
			return nil
		}
		time.Sleep(100 * time.Millisecond)
	}
	return fmt.Errorf("flai host did not start the MCP server again; flai host status says why")
}

// listening is every address the server st was listens on: Addrs, or, in a
// state written before S-0732, Addr alone.
func (st mcpState) listening() []string {
	if len(st.Addrs) == 0 {
		return []string{st.Addr}
	}
	return st.Addrs
}

// restartArgs start again, with a new token, the server that st was. One
// flai serve started goes on going with it, so that it stays flai serve's
// to supervise and to stop (S-0105: a rotation from the dashboard left one
// running after flai serve had gone).
func restartArgs(st mcpState, configPath string) []string {
	args := addrArgs(st.listening())
	if st.ExitWith > 0 && serve.Alive(st.ExitWith) {
		args = append(args, "--exit-with", strconv.Itoa(st.ExitWith))
	}
	if configPath != "" {
		args = append(args, "--config", configPath)
	}
	return args
}

// addrArgs start flai mcp http on every address in addrs.
func addrArgs(addrs []string) []string {
	args := []string{"mcp", "http"}
	for _, addr := range addrs {
		args = append(args, "--addr", addr)
	}
	return args
}
