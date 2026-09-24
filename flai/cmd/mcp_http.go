package cmd

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/cobra"

	"github.com/bytepunx/system-flow/flai/internal/buildinfo"
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
	mcpTokenFile    = "mcp.token"
	mcpStateFile    = "mcp-http.json"
	mcpAddrFile     = "mcp-http.addr"
	mcpLogFile      = "mcp-http.log"
	mcpStateEvery   = 3 * time.Second
	mcpStaleAfter   = 10 * time.Second
	mcpShutdownWait = 5 * time.Second
)

// mcpState is what a running server says about itself, rewritten while it
// runs so that a file left by a killed process is known to be stale.
type mcpState struct {
	PID     int    `json:"pid"`
	Addr    string `json:"addr"`
	URL     string `json:"url"`
	Started string `json:"started"`
	Updated string `json:"updated"`
	Version string `json:"version"`
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

// serveMCPHTTP runs the server until ctx ends.
func (a *app) serveMCPHTTP(ctx context.Context, repo *workitem.Repo, addr string, maxSessions int, idle time.Duration, exitWith int) error {
	if st, alive := readMCPState(repo, time.Now()); alive {
		return fmt.Errorf("flai mcp is already serving this project at %s (pid %d); flai mcp stop ends it", st.URL, st.PID)
	}
	token, _, err := mcpToken(repo, false)
	if err != nil {
		return err
	}
	log := a.logger().With("component", "mcp")
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("cannot listen on %s: %w; another project's server may hold it, so choose another with --addr", addr, err)
	}
	if !loopback(addr) {
		log.Warn("listening beyond this machine: the token travels in every request and flai does not encrypt it, so put a proxy or tunnel that terminates TLS in front", "addr", addr)
	}
	closing := make(chan struct{})
	handler := mcphttp.Handler(mcphttp.Options{
		Token: token, ProjectKey: repo.Manifest.Key, ProjectName: repo.Manifest.Name,
		MaxSessions: maxSessions, Idle: idle, Logger: log,
		NewServer: func(agent string) *mcp.Server {
			return mcpserver.New(mcpserver.Options{Repo: repo, Agent: agent, Version: buildinfo.Version, Now: a.now, Runner: a.runner, Closing: closing})
		},
	})
	mux := http.NewServeMux()
	mux.Handle(mcpPath, handler)
	// No write timeout: wait_for_events holds a request for up to five minutes.
	srv := &http.Server{Handler: mux, ReadHeaderTimeout: 10 * time.Second}

	st := mcpState{PID: os.Getpid(), Addr: ln.Addr().String(), URL: mcpURL(ln.Addr()), Started: time.Now().UTC().Format(time.RFC3339), Version: buildinfo.Version, ExitWith: exitWith}
	beat := func() error {
		st.Updated = time.Now().UTC().Format(time.RFC3339)
		return writeMCPState(repo, st)
	}
	if err := beat(); err != nil {
		_ = ln.Close()
		return err
	}
	defer func() { _ = os.Remove(mcpFile(repo, mcpStateFile)) }()
	log.Info("mcp server started", "url", st.URL, "root", repo.Root, "project", repo.Manifest.Key)
	fmt.Fprintf(a.errOut, "flai mcp serving %s at %s\n  token: %s (flai mcp token prints it)\n", repo.Manifest.Name, st.URL, mcpFile(repo, mcpTokenFile))

	failed := make(chan error, 1)
	go func() { failed <- srv.Serve(ln) }()
	tick := time.NewTicker(mcpStateEvery)
	defer tick.Stop()
	for {
		select {
		case err := <-failed:
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

func mcpClientConfig(url string) string {
	return fmt.Sprintf(`{ "mcpServers": { "flai": { "type": "http", "url": %q,
    "headers": { "Authorization": "Bearer <flai mcp token>", "X-Flai-Agent": "<the agent's name>" } } } }`, url)
}

func newMCPHTTPCmds(a *app) []*cobra.Command {
	var addr string
	var maxSessions int
	var idle time.Duration
	var exitWith int
	flags := func(c *cobra.Command) {
		c.Flags().StringVar(&addr, "addr", "", "address to listen on (default the last one used here, else "+mcpDefaultAddr+")")
		c.Flags().IntVar(&maxSessions, "max-sessions", 16, "MCP sessions at once")
		c.Flags().DurationVar(&idle, "idle", 30*time.Minute, "a session with no request for this long ends")
	}
	// The address sticks, so an agent's configuration keeps working across restarts.
	resolve := func(repo *workitem.Repo) string {
		if addr != "" {
			return addr
		}
		if data, err := os.ReadFile(mcpFile(repo, mcpAddrFile)); err == nil && strings.TrimSpace(string(data)) != "" {
			return strings.TrimSpace(string(data))
		}
		return mcpDefaultAddr
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
otherwise, and that address is remembered for this project. Every request
carries the token in .flai-cache/mcp.token as a bearer header; a request from
a browser (one with an Origin header) is refused. The agent is named by the
X-Flai-Agent header, else by the client's own name.

An address beyond this machine is allowed and warned about: the token travels
in every request, and TLS is a proxy's or a tunnel's job, not flai's.`,
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
			listen := resolve(repo)
			remember(repo, listen)
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
				listen := resolve(repo)
				args := []string{"mcp", "http", "--addr", listen, "--max-sessions", fmt.Sprint(maxSessions), "--idle", idle.String()}
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
writes a new one and restarts a running server, which ends every session.`,
		Args: cobra.NoArgs,
		RunE: func(*cobra.Command, []string) error {
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
	token.Flags().BoolVar(&rotate, "rotate", false, "replace the token")
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

// restartArgs start again, with a new token, the server that st was. One
// flai serve started goes on going with it, so that it stays flai serve's
// to supervise and to stop (S-0105: a rotation from the dashboard left one
// running after flai serve had gone).
func restartArgs(st mcpState, configPath string) []string {
	args := []string{"mcp", "http", "--addr", st.Addr}
	if st.ExitWith > 0 && serve.Alive(st.ExitWith) {
		args = append(args, "--exit-with", strconv.Itoa(st.ExitWith))
	}
	if configPath != "" {
		args = append(args, "--config", configPath)
	}
	return args
}
