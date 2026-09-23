package cmd

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"sort"
	"syscall"
	"time"

	"github.com/spf13/cobra"

	"github.com/bytepunx/system-flow/flai/internal/buildinfo"
	"github.com/bytepunx/system-flow/flai/internal/config"
	"github.com/bytepunx/system-flow/flai/internal/serve"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// flai serve is the process on the host that dashboards are reached through
// (ADR-0029). flai dashboard registers projects with it and starts it; this
// command runs it, and says how it is.

func (a *app) serveDir() serve.Dir {
	path := config.ResolvePath(a.configPath)
	if abs, err := filepath.Abs(path); err == nil {
		path = abs
	}
	return serve.DirFor(path)
}

func newServeCmd(a *app) *cobra.Command {
	c := &cobra.Command{
		Use:   "serve",
		Short: "Run flai on the host for the dashboards: it dials each registered project's dashboard and answers it",
		Long: `One flai serve per user serves every project registered with it. For each it
opens a WebSocket to that project's dashboard and keeps it open; the
dashboard asks, flai answers (ADR-0029). The dashboard never connects to the
host, and can ask only for the methods flai offers.

flai dashboard registers the project and starts flai serve in the background
when it is not running, so this command is for watching it work in a
terminal, and for start, stop, and status. Its registry, state, and log live
in a folder named serve beside flai's config file.

Started in a folder that is not a project, ~/git say, it also serves every
system-flow project below the folder, for as long as it runs, and offers the
folder's other git repositories for import on the board, as flai dashboard
there does.`,
		Example: `  flai serve            # in the foreground; Ctrl-C stops it
  flai serve start      # in the background
  flai serve status
  flai serve stop`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx, stop := signal.NotifyContext(cmd.Context(), os.Interrupt, syscall.SIGTERM)
			defer stop()
			// Started in a folder that is not a project (S-0102), it serves the
			// projects below it and offers the folder's other repositories for
			// import, as flai dashboard there does (ADR-0036).
			// The flai it runs for a dashboard reads and writes the configuration
			// this one does, --config included (S-0105): a setting changed from
			// the dashboard lands where flai serve reads it.
			if _, path, err := a.loadConfig(); err == nil && path != "" {
				_ = os.Setenv(config.EnvVar, path)
			}
			folder := ""
			if repo, err := a.projectOrNone(); err == nil && repo == nil {
				folder, _ = a.workingDir()
			}
			return serve.Run(ctx, serve.Options{Folder: folder, Dir: a.serveDir(), Version: buildinfo.Version, Logger: a.logger(), Now: a.now, Host: a.host(), Agent: a.agentConfig, MCP: &serveMCP{a: a}, ImportRoots: a.importRoots})
		},
	}
	c.AddCommand(
		&cobra.Command{
			Use:   "start",
			Short: "Start flai serve in the background, if it is not running",
			Args:  cobra.NoArgs,
			RunE: func(*cobra.Command, []string) error {
				st, started, err := a.ensureServe()
				if err != nil {
					return err
				}
				if a.jsonOut {
					return a.printJSON(map[string]any{"started": started, "pid": st.PID, "log": a.serveDir().Log()})
				}
				if started {
					fmt.Fprintf(a.out, "flai serve started (pid %d); log: %s\n", st.PID, a.serveDir().Log())
				} else {
					fmt.Fprintf(a.out, "flai serve is already running (pid %d, since %s)\n", st.PID, st.Started)
				}
				return nil
			},
		},
		&cobra.Command{
			Use:   "stop",
			Short: "Stop the running flai serve",
			Args:  cobra.NoArgs,
			RunE: func(*cobra.Command, []string) error {
				st, alive := a.serveDir().ReadStatus(time.Now())
				if !alive {
					fmt.Fprintln(a.out, "flai serve is not running")
					return nil
				}
				if err := serve.Terminate(st.PID); err != nil {
					return err
				}
				fmt.Fprintf(a.out, "flai serve stopped (pid %d)\n", st.PID)
				return nil
			},
		},
		newServeEnableCmd(a, true), newServeEnableCmd(a, false), newServeActionsCmd(a), newServeImportCmd(a), newServeJournalCmd(a), newServeAgentCmd(a), newServeChecksCmd(a),
		&cobra.Command{
			Use:   "status",
			Short: "Whether flai serve runs, which projects it serves, and which dashboards have it connected",
			Args:  cobra.NoArgs,
			RunE: func(*cobra.Command, []string) error {
				return a.printServeStatus()
			},
		},
	)
	return c
}

// serveStatus is flai serve status as data; flai dashboard status shows
// the part about its own project.
type serveStatus struct {
	Running  bool          `json:"running"`
	Status   *serve.Status `json:"status,omitempty"`
	Projects []serve.Entry `json:"projects"`
	Dir      string        `json:"dir"`
	// MCP is each project's HTTP MCP server that is running, by root: the
	// one flai serve keeps (S-0096, ADR-0034), or one started by hand.
	MCP map[string]mcpState `json:"mcp,omitempty"`
}

func (a *app) readServeStatus() (serveStatus, error) {
	dir := a.serveDir()
	projects, err := dir.Projects()
	if err != nil {
		return serveStatus{}, err
	}
	if projects == nil {
		projects = []serve.Entry{}
	}
	out := serveStatus{Projects: projects, Dir: string(dir)}
	if st, alive := dir.ReadStatus(time.Now()); alive {
		out.Running, out.Status = true, &st
	}
	for _, p := range projects {
		repo, err := workitem.Open(p.Root)
		if err != nil {
			continue
		}
		if st, alive := readMCPState(repo, time.Now()); alive {
			if out.MCP == nil {
				out.MCP = map[string]mcpState{}
			}
			out.MCP[p.Root] = st
		}
	}
	return out, nil
}

func (a *app) printServeStatus() error {
	st, err := a.readServeStatus()
	if err != nil {
		return err
	}
	if a.jsonOut {
		return a.printJSON(st)
	}
	if !st.Running {
		fmt.Fprintln(a.out, "flai serve is not running (flai serve start, or flai dashboard in a project)")
	} else {
		fmt.Fprintf(a.out, "flai serve %s running since %s (pid %d)\n", st.Status.Version, st.Status.Started, st.Status.PID)
	}
	sort.Slice(st.Projects, func(i, j int) bool { return st.Projects[i].Key < st.Projects[j].Key })
	for _, p := range st.Projects {
		line := "not connected"
		if st.Running {
			if c, ok := st.Status.Connections[p.Root]; ok && c.Connected {
				line = "connected since " + c.Since
			} else if ok && c.LastError != "" {
				line = "not connected: " + c.LastError
			}
		}
		fmt.Fprintf(a.out, "  %s  %s  %s\n    %s\n", p.Key, p.URL, line, p.Root)
		if m, ok := st.MCP[p.Root]; ok {
			fmt.Fprintf(a.out, "    mcp: %s (pid %d); flai mcp status prints an agent's configuration\n", m.URL, m.PID)
		}
	}
	if len(st.Projects) == 0 {
		fmt.Fprintln(a.out, "  no projects registered; flai dashboard in a project registers it")
	}
	if st.Running && st.Status.Folder != "" {
		fmt.Fprintf(a.out, "started in %s: serves the projects below it, and offers its other repositories for import\n", st.Status.Folder)
		for _, p := range st.Status.FolderProjects {
			line := "not connected"
			if c, ok := st.Status.Connections[p.Root]; ok && c.Connected {
				line = "connected since " + c.Since
			}
			fmt.Fprintf(a.out, "  %s  %s  %s\n    %s\n", p.Key, p.URL, line, p.Root)
		}
	}
	if st.Running && len(st.Status.Offered) > 0 {
		fmt.Fprintln(a.out, "offered for import (flai serve import):")
		for _, c := range st.Status.Offered {
			fmt.Fprintf(a.out, "  %s  %s\n", c.Name, c.Root)
		}
	}
	return nil
}

// ensureServe starts flai serve detached unless one is running. It needs no
// root: the process is this user's, and it stops with flai serve stop.
func (a *app) ensureServe() (serve.Status, bool, error) {
	dir := a.serveDir()
	if st, alive := dir.ReadStatus(time.Now()); alive {
		return st, false, nil
	}
	exe, err := os.Executable()
	if err != nil {
		return serve.Status{}, false, err
	}
	if err := os.MkdirAll(string(dir), 0o700); err != nil {
		return serve.Status{}, false, err
	}
	logFile, err := os.OpenFile(dir.Log(), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o600)
	if err != nil {
		return serve.Status{}, false, err
	}
	defer func() { _ = logFile.Close() }()
	args := []string{"serve", "--config", config.ResolvePath(a.configPath)}
	cmd := exec.CommandContext(context.WithoutCancel(context.Background()), exe, args...)
	cmd.Stdout, cmd.Stderr = logFile, logFile
	cmd.Env = append(os.Environ(), "LOG_FORMAT=json")
	serve.Detach(cmd)
	if err := cmd.Start(); err != nil {
		return serve.Status{}, false, err
	}
	pid := cmd.Process.Pid
	_ = cmd.Process.Release()
	// It has started when it has written its first status.
	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		if st, alive := dir.ReadStatus(time.Now()); alive && st.PID == pid {
			return st, true, nil
		}
		time.Sleep(100 * time.Millisecond)
	}
	return serve.Status{}, false, fmt.Errorf("flai serve did not come up; see %s", dir.Log())
}
