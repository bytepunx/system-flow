package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/spf13/cobra"

	"github.com/bytepunx/system-flow/flai/internal/buildinfo"
	"github.com/bytepunx/system-flow/flai/internal/config"
	"github.com/bytepunx/system-flow/flai/internal/host"
	"github.com/bytepunx/system-flow/flai/internal/serve"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// flai host is the one process per machine that runs flai serve and each
// served project's MCP server as its children (S-0106). flai dashboard
// starts it; this command runs it in a terminal, starts and stops it, and
// asks it to start, stop, or restart its children, and to upgrade flai.

func (a *app) hostDir() host.Dir {
	path := config.ResolvePath(a.configPath)
	if abs, err := filepath.Abs(path); err == nil {
		path = abs
	}
	return host.DirFor(path)
}

func newHostCmd(a *app) *cobra.Command {
	c := &cobra.Command{
		Use:   "host",
		Short: "Run flai host: the one process per machine that keeps flai serve and the MCP servers running",
		Long: `flai host runs flai serve, and for every project flai serve serves that
project's HTTP MCP server (flai mcp http), as its own children. It starts a
child again when it ends, with a wait that grows while it keeps failing, and
stops every child when it stops, or within seconds of being killed, since
each child goes with it (--exit-with).

One host runs per machine: it holds ` + host.DefaultAddr + ` (or $` + host.AddrEnv + `), a
control API that answers its token alone and no web page. flai serve reaches
it there to say which projects need an MCP server, and the dashboard reaches
it through flai serve for status, restarts, and upgrades (flai serve enable
host). Its state, token, and log live in a folder named host beside flai's
config file.

flai dashboard starts it when it is not running, so this command is for
watching it work in a terminal, and for start, stop, status, and the rest.`,
		Example: `  flai host                  # in the foreground; Ctrl-C stops it and its children
  flai host start            # in the background
  flai host status
  flai host restart serve    # serve, mcp, or all
  flai host stop mcp         # stop the MCP servers until flai host start mcp
  flai host check            # is a newer flai published?
  flai host upgrade          # install it and restart everything on it
  flai host versions         # the published flai releases, newest first
  flai host upgrade --version 1.2.0   # install that one instead, an earlier one included
  flai host stop`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx, stop := signal.NotifyContext(cmd.Context(), os.Interrupt, syscall.SIGTERM)
			defer stop()
			return a.runHost(ctx)
		},
	}
	process := func(action, short string) *cobra.Command {
		return &cobra.Command{
			Use:       action + " <serve|mcp|all>",
			Short:     short,
			Args:      cobra.ExactArgs(1),
			ValidArgs: []string{host.Serve, host.MCP, host.All},
			RunE: func(cmd *cobra.Command, args []string) error {
				return a.hostAct(cmd.Context(), args[0], action)
			},
		}
	}
	start := &cobra.Command{
		Use:   "start [serve|mcp|all]",
		Short: "Start flai host in the background, if it is not running; with a process, ask the host to start it",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 1 {
				return a.hostAct(cmd.Context(), args[0], "start")
			}
			st, started, err := a.ensureHost()
			if err != nil {
				return err
			}
			if a.jsonOut {
				return a.printJSON(map[string]any{"started": started, "pid": st.PID, "addr": st.Addr, "log": a.hostDir().Log()})
			}
			if started {
				fmt.Fprintf(a.out, "flai host started (pid %d, %s); it runs flai serve; log: %s\n", st.PID, st.Addr, a.hostDir().Log())
			} else {
				fmt.Fprintf(a.out, "flai host is already running (pid %d, since %s)\n", st.PID, st.Started)
			}
			return nil
		},
	}
	stop := &cobra.Command{
		Use:   "stop [serve|mcp|all]",
		Short: "Stop flai host and every process it runs; with a process, ask the host to stop that one only",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 1 {
				return a.hostAct(cmd.Context(), args[0], "stop")
			}
			return a.stopHost()
		},
	}
	status := &cobra.Command{
		Use:   "status",
		Short: "Whether flai host runs, and each process it keeps: its state, PID, version, and restarts",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return a.printHostStatus(cmd.Context())
		},
	}
	check := &cobra.Command{
		Use:   "check",
		Short: "Ask the host whether a newer flai is published than the one it runs",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			c, _, err := a.hostDir().Client(time.Now())
			if err != nil {
				return err
			}
			out, err := c.Check(cmd.Context())
			if err != nil {
				return err
			}
			return a.printRaw(out, func(m map[string]any) string {
				state := "upgrade available"
				if m["up_to_date"] == true {
					state = "up to date"
				}
				return fmt.Sprintf("flai %v runs; latest is %v (%s)\n", m["current"], m["latest"], state)
			})
		},
	}
	versions := &cobra.Command{
		Use:   "versions",
		Short: "Ask the host for the published flai releases, newest first, marking the installed one and the newest",
		Long: `Asks the host for the published flai releases, newest first, as flai
self-upgrade --list prints them from the host's flai: each with its date,
marking the installed one (the flai flai host upgrade would replace), the
newest, and any below the flai.minimum of a project flai serve serves. With
--json it prints the list as the host answers it: version, tag, published,
installed, latest, and below_minimum (each project and minimum the release is
below, left out when none). It installs nothing.`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			c, st, err := a.hostDir().Client(time.Now())
			if err != nil {
				return err
			}
			out, err := c.Versions(cmd.Context())
			if err != nil {
				return err
			}
			return a.printHostVersions(out, st.Version)
		},
	}
	var version string
	upgrade := &cobra.Command{
		Use:   "upgrade",
		Short: "Have the host install the newest flai, or --version a published one, and restart itself and every process it runs on it",
		Long: `Has the host install the newest flai, as flai self-upgrade does, and
restart itself and every process it runs on what it installed.

--version installs that published release instead, an earlier one included,
the same way: the host restarts on it with its processes. A version that is
not X.Y.Z is refused before anything runs, and one that is not a published
release is refused, naming the published ones (flai host versions), before
anything is downloaded or replaced. A release below a served project's
flai.minimum is installed, not refused: flai host versions marks it.`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			c, _, err := a.hostDir().Client(time.Now())
			if err != nil {
				return err
			}
			out, err := c.Upgrade(cmd.Context(), version)
			if err != nil {
				return err
			}
			return a.printRaw(out, func(m map[string]any) string {
				up, _ := m["upgrade"].(map[string]any)
				if m["restarting"] == true {
					return fmt.Sprintf("installed flai %v (was %v); the host restarts on it with its processes\n", up["installed"], up["previous"])
				}
				return fmt.Sprintf("flai %v is already the latest release\n", up["current"])
			})
		},
	}
	upgrade.Flags().StringVar(&version, "version", "", "install this published release instead of the newest, e.g. 1.2.0")
	c.AddCommand(start, stop, status, process("restart", "Ask the host to restart serve, the MCP servers, or all of them"), check, versions, upgrade)
	return c
}

// printHostVersions prints the host's list of published flai releases as it
// is, or as flai self-upgrade --list prints it, under the version the host runs.
func (a *app) printHostVersions(out json.RawMessage, running string) error {
	var listed []listedFlai
	if err := json.Unmarshal(out, &listed); err != nil {
		return fmt.Errorf("flai host answered a list of releases that is not one: %w", err)
	}
	if a.jsonOut {
		var v []any
		if err := json.Unmarshal(out, &v); err != nil {
			return err
		}
		return a.printJSON(v)
	}
	if len(listed) == 0 {
		fmt.Fprintf(a.out, "no published flai release (the host runs flai %s)\n", running)
		return nil
	}
	fmt.Fprintf(a.out, "published flai releases, newest first (the host runs flai %s):\n", running)
	if err := writeFlaiReleases(a.out, listed); err != nil {
		return err
	}
	found := false
	for _, r := range listed {
		found = found || r.Installed
	}
	if !found {
		fmt.Fprintln(a.out, "the flai the host would replace is not a published release")
	}
	return nil
}

// printRaw prints the host's JSON answer as it is, or as a line.
func (a *app) printRaw(out json.RawMessage, line func(map[string]any) string) error {
	if a.jsonOut {
		var v any
		if err := json.Unmarshal(out, &v); err != nil {
			return err
		}
		return a.printJSON(v)
	}
	var m map[string]any
	_ = json.Unmarshal(out, &m)
	fmt.Fprint(a.out, line(m))
	return nil
}

func (a *app) hostAct(ctx context.Context, process, action string) error {
	c, _, err := a.hostDir().Client(time.Now())
	if err != nil {
		return err
	}
	st, err := c.Act(ctx, process, action)
	if err != nil {
		return err
	}
	if a.jsonOut {
		return a.printJSON(st)
	}
	fmt.Fprintf(a.out, "flai host: %s %s asked\n", action, process)
	a.describeHost(st)
	return nil
}

// hostStatus is flai host status as data.
type hostStatus struct {
	Running bool         `json:"running"`
	Status  *host.Status `json:"status,omitempty"`
	Dir     string       `json:"dir"`
	// Elsewhere is a host that holds the machine's address for another
	// config file.
	Elsewhere *host.Health `json:"elsewhere,omitempty"`
}

func (a *app) readHostStatus(ctx context.Context) hostStatus {
	dir := a.hostDir()
	out := hostStatus{Dir: string(dir)}
	if c, _, err := dir.Client(time.Now()); err == nil {
		if st, err := c.Status(ctx); err == nil {
			out.Running, out.Status = true, &st
			return out
		}
	}
	if st, alive := dir.ReadStatus(time.Now()); alive {
		out.Running, out.Status = true, &st
		return out
	}
	if h, ok := host.Probe(host.Addr()); ok {
		out.Elsewhere = &h
	}
	return out
}

func (a *app) printHostStatus(ctx context.Context) error {
	st := a.readHostStatus(ctx)
	if a.jsonOut {
		return a.printJSON(st)
	}
	switch {
	case st.Elsewhere != nil:
		fmt.Fprintf(a.out, "flai host is not running for this config; the machine's host (pid %d, flai %s) runs for %s\n", st.Elsewhere.PID, st.Elsewhere.Version, st.Elsewhere.Config)
	case !st.Running:
		fmt.Fprintln(a.out, "flai host is not running (flai host start, or flai dashboard)")
	default:
		fmt.Fprintf(a.out, "flai host %s running since %s (pid %d, %s)\n", st.Status.Version, st.Status.Started, st.Status.PID, st.Status.Addr)
		a.describeHost(*st.Status)
	}
	return nil
}

func (a *app) describeHost(st host.Status) {
	for _, c := range st.Children {
		name := c.Name
		if c.Root != "" {
			name += " " + c.Root
		}
		line := c.State
		switch c.State {
		case "running":
			line = fmt.Sprintf("running since %s (pid %d, flai %s)", c.Since, c.PID, c.Version)
		case "external":
			line = fmt.Sprintf("run by a process the host did not start (pid %d); stop it for the host to take over", c.PID)
		case "stopping":
			line = fmt.Sprintf("stopping (pid %d): asked to end, and given the grace period", c.PID)
		}
		if c.Restarts > 0 {
			line += fmt.Sprintf(", %d restart(s)", c.Restarts)
		}
		fmt.Fprintf(a.out, "  %s: %s\n", name, line)
		if c.LastError != "" {
			fmt.Fprintf(a.out, "    last error: %s\n", c.LastError)
		} else if c.State != "running" && c.LastExit != "" {
			fmt.Fprintf(a.out, "    last exit: %s\n", c.LastExit)
		}
	}
	if d := st.Dashboard; d != nil {
		line := d.State
		switch d.State {
		case host.DashboardOff:
			line = "not watched (none started, stopped with flai dashboard stop, or dashboard.no_restart set)"
		case host.DashboardNotAnswering:
			line = "not answering"
		}
		if d.Restarts > 0 {
			line += fmt.Sprintf(", %d restart(s), the last at %s because it was %s", d.Restarts, d.LastRestart, d.LastReason)
		}
		fmt.Fprintf(a.out, "  dashboard: %s\n", line)
		if d.LastError != "" {
			fmt.Fprintf(a.out, "    last error: %s\n", d.LastError)
		}
	}
}

// ensureHost starts flai host detached unless one runs for this config. A
// host that holds the machine's address for another config is an error: it
// runs another flai serve, which does not serve this config's projects.
func (a *app) ensureHost() (host.Status, bool, error) {
	if a.hostStarter != nil {
		return a.hostStarter()
	}
	dir := a.hostDir()
	if st, alive := dir.ReadStatus(time.Now()); alive {
		return st, false, nil
	}
	if h, ok := host.Probe(host.Addr()); ok {
		return host.Status{}, false, fmt.Errorf("flai host already runs on this machine for %s (pid %d), not for %s; use that config, or stop it with flai host stop --config %s", h.Config, h.PID, config.ResolvePath(a.configPath), h.Config)
	}
	exe, err := os.Executable()
	if err != nil {
		return host.Status{}, false, err
	}
	if err := os.MkdirAll(string(dir), 0o700); err != nil {
		return host.Status{}, false, err
	}
	logFile, err := os.OpenFile(dir.Log(), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o600)
	if err != nil {
		return host.Status{}, false, err
	}
	defer func() { _ = logFile.Close() }()
	cmd := exec.CommandContext(context.WithoutCancel(context.Background()), exe, "host", "--config", config.ResolvePath(a.configPath))
	// the working directory goes on to flai serve: a folder that is not a
	// project is served whole (S-0102)
	cmd.Dir, _ = a.workingDir()
	cmd.Stdout, cmd.Stderr = logFile, logFile
	cmd.Env = append(os.Environ(), "LOG_FORMAT=json")
	serve.Detach(cmd)
	if err := cmd.Start(); err != nil {
		return host.Status{}, false, err
	}
	pid := cmd.Process.Pid
	exited := make(chan struct{})
	go func() { _ = cmd.Wait(); close(exited) }()
	deadline := time.After(5 * time.Second)
	for {
		if st, alive := dir.ReadStatus(time.Now()); alive && st.PID == pid {
			return st, true, nil
		}
		select {
		case <-exited:
			return host.Status{}, false, fmt.Errorf("flai host stopped at once; see %s", dir.Log())
		case <-deadline:
			return host.Status{}, false, fmt.Errorf("flai host did not come up; see %s", dir.Log())
		case <-time.After(100 * time.Millisecond):
		}
	}
}

func (a *app) stopHost() error {
	dir := a.hostDir()
	st, alive := dir.ReadStatus(time.Now())
	if !alive {
		fmt.Fprintln(a.out, "flai host is not running")
		return nil
	}
	if err := serve.Terminate(st.PID); err != nil {
		return err
	}
	// it stops its children first, each given its grace
	deadline := time.Now().Add(hostGrace + 5*time.Second)
	for time.Now().Before(deadline) && serve.Alive(st.PID) {
		time.Sleep(100 * time.Millisecond)
	}
	if a.jsonOut {
		return a.printJSON(map[string]any{"stopped": true, "pid": st.PID})
	}
	fmt.Fprintf(a.out, "flai host stopped (pid %d) with every process it ran\n", st.PID)
	return nil
}

// hostGrace is how long the host gives a child to stop before it kills it:
// flai mcp http takes up to mcpShutdownWait to let its agents go.
const hostGrace = mcpShutdownWait + 5*time.Second

// runHost runs the host in this process until ctx ends, and when an upgrade
// installed a new flai, becomes that flai.
func (a *app) runHost(ctx context.Context) error {
	_, path, err := a.loadConfig()
	if err != nil {
		return err
	}
	if abs, err := filepath.Abs(path); err == nil {
		path = abs
	}
	exe, err := os.Executable()
	if err != nil {
		return err
	}
	wd, _ := a.workingDir()
	l := &hostLauncher{a: a, exe: exe, config: path, dir: wd, taken: map[string]string{}}
	err = host.Run(ctx, host.Options{
		Dir: a.hostDir(), Addr: host.Addr(), Version: buildinfo.Version, Config: path, Logger: a.logger(),
		Serve: l.serve, MCP: l.mcp, Check: l.check, Versions: l.versions, Upgrade: l.upgradeTo, Grace: hostGrace,
		Dashboard: a.dashboardWatch(),
	})
	if errors.Is(err, host.ErrRestart) {
		// the upgrade may have installed elsewhere than this binary: a flai
		// run from a checkout installs under ~/.flai/bin (S-0111)
		if to := l.installedTo(); to != "" {
			exe = to
		}
		a.logger().Info("host restarting on the upgraded flai", "component", "host", "exe", exe)
		return host.Reexec(exe)
	}
	return err
}

// hostLauncher says how the host starts its children: this flai's own
// commands, with this config, each going with the host.
type hostLauncher struct {
	a      *app
	exe    string
	config string
	dir    string

	mu        sync.Mutex
	taken     map[string]string // the MCP address given to each root, so two never share one
	installed string            // where the last upgrade installed flai
	ver       string
	verAt     time.Time
}

// version is the flai the children start from: the binary on disk, which an
// upgrade or an install may have replaced since the host started.
func (l *hostLauncher) version() string {
	l.mu.Lock()
	defer l.mu.Unlock()
	fi, err := os.Stat(l.exe)
	if err == nil && l.ver != "" && fi.ModTime().Equal(l.verAt) {
		return l.ver
	}
	out, err := exec.Command(l.exe, "version", "--json").Output()
	var info struct {
		Version string `json:"version"`
	}
	if err != nil || json.Unmarshal(out, &info) != nil {
		return buildinfo.Version
	}
	if fi != nil {
		l.ver, l.verAt = info.Version, fi.ModTime()
	}
	return info.Version
}

func (l *hostLauncher) withConfig(args ...string) []string {
	return append(args, "--config", l.config, "--exit-with", strconv.Itoa(os.Getpid()))
}

func (l *hostLauncher) serve() (host.Spec, error) {
	dir := serve.DirFor(l.config)
	if st, alive := dir.ReadStatus(time.Now()); alive && st.PID != os.Getpid() {
		return host.Spec{External: st.PID}, nil
	}
	if err := os.MkdirAll(string(dir), 0o700); err != nil {
		return host.Spec{}, err
	}
	return host.Spec{Path: l.exe, Args: l.withConfig("serve"), Dir: l.dir, Log: dir.Log(), Env: []string{"LOG_FORMAT=json"}, Version: l.version()}, nil
}

func (l *hostLauncher) mcp(root string) (host.Spec, error) {
	repo, err := workitem.Open(root)
	if err != nil {
		return host.Spec{}, err
	}
	if st, alive := readMCPState(repo, time.Now()); alive {
		return host.Spec{External: st.PID}, nil
	}
	addr, err := l.mcpAddr(repo)
	if err != nil {
		return host.Spec{}, err
	}
	if err := os.MkdirAll(repo.CacheDir(), 0o700); err != nil {
		return host.Spec{}, err
	}
	return host.Spec{Path: l.exe, Args: l.withConfig("mcp", "http", "--addr", addr), Dir: repo.MainRoot,
		Log: mcpFile(repo, mcpLogFile), Env: []string{"LOG_FORMAT=json"}, Version: l.version()}, nil
}

// mcpAddr is mcpAddrFor, with the addresses already given to other projects
// left out: two starting at once must not pick the same free port.
func (l *hostLauncher) mcpAddr(repo *workitem.Repo) (string, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if addr, ok := l.taken[repo.MainRoot]; ok {
		return addr, nil
	}
	others := map[string]bool{}
	for _, addr := range l.taken {
		others[addr] = true
	}
	// the ports the other projects flai serve serves remember, so that a new
	// one is not one of theirs (S-0110)
	remembered := map[string]bool{}
	if entries, err := l.a.serveDir().Projects(); err == nil {
		for _, e := range entries {
			if e.Root == repo.MainRoot {
				continue
			}
			if other, err := workitem.Open(e.Root); err == nil {
				if addr := rememberedMCPAddr(other); addr != "" {
					remembered[addr] = true
				}
			}
		}
	}
	addr, err := mcpAddrFor(repo, others, remembered)
	if err != nil {
		return "", err
	}
	if was := rememberedMCPAddr(repo); was != "" && was != addr {
		// the server remembers the new one when it starts; an agent given the
		// old one must be given this
		l.a.logger().Warn("mcp address moved: another project the host runs has it", "component", "host", "root", repo.MainRoot, "from", was, "to", addr)
	}
	l.taken[repo.MainRoot] = addr
	return addr, nil
}

// flaiRaw runs this flai with --json and answers the JSON it prints.
func (l *hostLauncher) flaiRaw(ctx context.Context, args ...string) (json.RawMessage, error) {
	cmd := exec.CommandContext(ctx, l.exe, append(args, "--json", "--config", l.config)...)
	cmd.Env = append(os.Environ(), "LOG_FORMAT=json")
	out, err := cmd.Output()
	if err != nil {
		var ee *exec.ExitError
		if errors.As(err, &ee) {
			return nil, fmt.Errorf("flai %s: %s", args[0], fatalOf(ee.Stderr))
		}
		return nil, err
	}
	if !json.Valid(out) {
		return nil, fmt.Errorf("flai %s printed what is not JSON: %.200q", args[0], out)
	}
	return out, nil
}

// flaiJSON runs this flai with --json and reads the object it prints.
func (l *hostLauncher) flaiJSON(ctx context.Context, args ...string) (map[string]any, error) {
	out, err := l.flaiRaw(ctx, args...)
	if err != nil {
		return nil, err
	}
	var m map[string]any
	if err := json.Unmarshal(out, &m); err != nil {
		return nil, fmt.Errorf("flai %s printed what is not a JSON object: %w", args[0], err)
	}
	return m, nil
}

// fatalOf is the error a failed flai logged, else its last line.
func fatalOf(stderr []byte) string {
	last := ""
	for _, line := range strings.Split(string(stderr), "\n") {
		var ev map[string]any
		if json.Unmarshal([]byte(line), &ev) == nil {
			if e, ok := ev["err"].(string); ok && ev["level"] == "FATAL" {
				return e
			}
		}
		if line != "" {
			last = line
		}
	}
	return last
}

func (l *hostLauncher) check(ctx context.Context) (any, error) {
	return l.flaiJSON(ctx, "self-upgrade", "--check")
}

// versions lists the published flai releases as self-upgrade --list does.
func (l *hostLauncher) versions(ctx context.Context) (any, error) {
	return l.flaiRaw(ctx, "self-upgrade", "--list")
}

// upgradeTo installs the newest flai, or the published release version names,
// which self-upgrade refuses, naming the published ones, when it is not one.
func (l *hostLauncher) upgradeTo(ctx context.Context, version string) (any, bool, error) {
	args := []string{"self-upgrade"}
	if version != "" {
		args = append(args, "--version", version)
	}
	m, err := l.flaiJSON(ctx, args...)
	if err != nil {
		return nil, false, err
	}
	_, installed := m["installed"]
	if path, ok := m["path"].(string); ok && installed {
		l.mu.Lock()
		l.installed = path
		l.mu.Unlock()
	}
	return m, installed, nil
}

// installedTo is where the last upgrade installed flai, "" when none did.
func (l *hostLauncher) installedTo() string {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.installed
}

// hostMCP is flai serve's way to have its projects' MCP servers kept: it
// tells the host that started it.
type hostMCP struct{ c *host.Client }

// Keep implements serve.MCP.
func (m hostMCP) Keep(roots []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := m.c.KeepMCP(ctx, roots)
	return err
}

// freePort says whether addr can be listened on now.
func freePort(addr string) bool {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return false
	}
	_ = ln.Close()
	return true
}
