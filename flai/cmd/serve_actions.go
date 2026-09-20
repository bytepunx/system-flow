package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/bytepunx/system-flow/flai/internal/config"
	"github.com/bytepunx/system-flow/flai/internal/hostapi"
	"github.com/bytepunx/system-flow/flai/internal/serve"
)

// Host actions (S-0078, ADR-0029): what flai serve may do on this host, with
// the operator's credentials, because a dashboard asked. None is enabled
// until the operator enables it by name, here, in a shell on the host. The
// setting lives in the host's flai configuration; the journal lies beside
// flai serve's state. No method a dashboard can ask for touches either.

// host is the operator's say, read from the configuration at every question,
// so that enabling and disabling take effect without a restart.
func (a *app) host() hostapi.Host {
	return hostapi.Host{
		Enabled: func(action, root string) bool {
			cfg, _, err := a.loadConfig()
			return err == nil && cfg.ActionEnabled(action, root)
		},
		Record: func(e hostapi.Entry) {
			if err := a.journal(e); err != nil {
				a.logger().Warn("host action not journalled", "component", "serve", "action", e.Action, "err", err.Error())
			}
			a.logger().Info("host action", "component", "serve", "action", e.Action, "method", e.Method, "project", e.Project, "by", e.By, "outcome", e.Outcome, "detail", e.Detail)
		},
	}
}

func (a *app) journal(e hostapi.Entry) error {
	dir := a.serveDir()
	if err := os.MkdirAll(string(dir), 0o700); err != nil {
		return err
	}
	line, err := json.Marshal(e)
	if err != nil {
		return err
	}
	f, err := os.OpenFile(dir.Journal(), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o600)
	if err != nil {
		return err
	}
	defer func() { _ = f.Close() }()
	_, err = f.Write(append(line, '\n'))
	return err
}

func knownAction(name string) error {
	if _, ok := hostapi.Actions[name]; ok {
		return nil
	}
	names := make([]string, 0, len(hostapi.Actions))
	for n := range hostapi.Actions {
		names = append(names, n)
	}
	sort.Strings(names)
	return fmt.Errorf("there is no host action %q; there %s: %s", name, oneOrMany(len(names), "is", "are"), strings.Join(names, ", "))
}

func oneOrMany(n int, one, many string) string {
	if n == 1 {
		return one
	}
	return many
}

func newServeEnableCmd(a *app, on bool) *cobra.Command {
	var all bool
	use, short := "enable <action>", "Let flai serve perform a host action for this project when its dashboard asks"
	if !on {
		use, short = "disable <action>", "Stop flai serve performing a host action for this project"
	}
	c := &cobra.Command{
		Use:   use,
		Short: short,
		Long: `A host action is something flai serve does on this host, as you and with your
credentials, because a dashboard asked: flai serve actions lists them and
says what each one means. None is enabled until you enable it by name, in a
shell on the host; nothing a dashboard can ask for changes that. It applies
to the project in the working directory, or to every project with
--all-projects, and takes effect at once, with no restart. disable
--all-projects removes it everywhere.`,
		Example: `  flai serve enable push
  flai serve actions
  flai serve disable push --all-projects`,
		Args: cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			action := args[0]
			if err := knownAction(action); err != nil {
				return err
			}
			root, where := config.AllProjects, "every project"
			if !all {
				repo, err := a.project()
				if err != nil {
					return fmt.Errorf("%w; or name every project with --all-projects", err)
				}
				root, where = mainRootOf(repo), repo.Manifest.Name
			}
			cfg, path, err := a.loadConfig()
			if err != nil {
				return err
			}
			next := cfg.WithAction(action, root, on)
			if err := config.Save(path, next); err != nil {
				return err
			}
			still := !on && !all && next.ActionEnabled(action, root)
			if a.jsonOut {
				return a.printJSON(map[string]any{"action": action, "enabled": next.ActionEnabled(action, root), "for": root, "config": path})
			}
			switch {
			case on:
				fmt.Fprintf(a.out, "%s enabled for %s\n  it lets a dashboard %s\n  every use is journalled: flai serve journal\n  turn it off with: flai serve disable %s\n", action, where, hostapi.Actions[action], action)
			case still:
				fmt.Fprintf(a.out, "%s is still enabled for %s, because it is enabled for every project; flai serve disable %s --all-projects turns it off everywhere\n", action, where, action)
			default:
				fmt.Fprintf(a.out, "%s disabled for %s\n", action, where)
			}
			return nil
		},
	}
	c.Flags().BoolVar(&all, "all-projects", false, "every project on this host, not only the one in the working directory")
	return c
}

func newServeActionsCmd(a *app) *cobra.Command {
	return &cobra.Command{
		Use:   "actions",
		Short: "The host actions there are, what each lets a dashboard do, and where each is enabled",
		Args:  cobra.NoArgs,
		RunE: func(*cobra.Command, []string) error {
			cfg, _, err := a.loadConfig()
			if err != nil {
				return err
			}
			here := ""
			if repo, err := a.project(); err == nil {
				here = mainRootOf(repo)
			}
			names := make([]string, 0, len(hostapi.Actions))
			for n := range hostapi.Actions {
				names = append(names, n)
			}
			sort.Strings(names)
			if a.jsonOut {
				out := []map[string]any{}
				for _, n := range names {
					roots := cfg.HostActions[n]
					if roots == nil {
						roots = []string{}
					}
					out = append(out, map[string]any{"action": n, "means": hostapi.Actions[n], "enabled_for": roots, "enabled_here": here != "" && cfg.ActionEnabled(n, here)})
				}
				return a.printJSON(out)
			}
			for _, n := range names {
				state := "off everywhere"
				if roots := cfg.HostActions[n]; len(roots) > 0 {
					state = "on for " + strings.ReplaceAll(strings.Join(roots, ", "), config.AllProjects, "every project")
				}
				fmt.Fprintf(a.out, "%s: %s\n  lets a dashboard %s\n", n, state, hostapi.Actions[n])
				if here != "" {
					if cfg.ActionEnabled(n, here) {
						fmt.Fprintf(a.out, "  on for this project; flai serve disable %s turns it off\n", n)
					} else {
						fmt.Fprintf(a.out, "  off for this project; flai serve enable %s turns it on\n", n)
					}
				}
			}
			return nil
		},
	}
}

func newServeJournalCmd(a *app) *cobra.Command {
	var last int
	c := &cobra.Command{
		Use:   "journal",
		Short: "Every host action a dashboard asked for: when, what, which project, for whom, and what became of it",
		Long: `The journal is written by flai serve on this host, one line per host action
asked for, refused ones included, in journal.jsonl beside its state. It is
append-only and yours: nothing a dashboard can ask for reads or changes it.`,
		Args: cobra.NoArgs,
		RunE: func(*cobra.Command, []string) error {
			path := a.serveDir().Journal()
			f, err := os.Open(path)
			if err != nil {
				if os.IsNotExist(err) {
					if a.jsonOut {
						return a.printJSON([]hostapi.Entry{})
					}
					fmt.Fprintln(a.out, "no host action has been asked for on this host")
					return nil
				}
				return err
			}
			defer func() { _ = f.Close() }()
			var entries []hostapi.Entry
			sc := bufio.NewScanner(f)
			sc.Buffer(make([]byte, 0, 64*1024), 1024*1024)
			for sc.Scan() {
				var e hostapi.Entry
				if json.Unmarshal(sc.Bytes(), &e) == nil && e.At != "" {
					entries = append(entries, e)
				}
			}
			if err := sc.Err(); err != nil {
				return err
			}
			if last > 0 && len(entries) > last {
				entries = entries[len(entries)-last:]
			}
			if a.jsonOut {
				if entries == nil {
					entries = []hostapi.Entry{}
				}
				return a.printJSON(entries)
			}
			for _, e := range entries {
				fmt.Fprintf(a.out, "%s  %-8s %-10s %s  for %s  %s", e.At, e.Outcome, e.Action, e.Project, e.By, e.Method)
				if e.Detail != "" {
					fmt.Fprintf(a.out, "\n    %s", e.Detail)
				}
				fmt.Fprintln(a.out)
			}
			fmt.Fprintf(a.errOut, "%d entr%s from %s\n", len(entries), oneOrMany(len(entries), "y", "ies"), path)
			return nil
		},
	}
	c.Flags().IntVarP(&last, "last", "n", 0, "only the last n entries")
	return c
}

// flai serve agent: the command flai serve starts an agent with (S-0079).
func newServeAgentCmd(a *app) *cobra.Command {
	c := &cobra.Command{
		Use:   "agent",
		Short: "The command flai serve starts when a story becomes ready and nobody is attending",
		Long: `flai serve can start an agent session for you when a story becomes ready and
no agent is attending the project. It is a host action, off until you enable
it (flai serve enable agent), and it has no default command: you write one.

The command is an argument list, run as it stands in the project's directory,
as you, never through a shell. In an argument {story} is replaced by the
story's ID and {root} by the project's directory; nothing else is interpreted.
The session's environment carries FLAI_AGENT, FLAI_STORY, and FLAI_SESSION.
One agent is started per project at a time; its output goes to a log beside
flai serve's state, and every start and failure is in flai serve journal.`,
		Example: `  flai serve agent set -- claude -p "Work on story {story} as the conventions say"
  flai serve agent set --name builder -- /home/me/bin/start-agent {story}
  flai serve agent show
  flai serve enable agent
  flai serve agent clear`,
		Args: cobra.NoArgs,
		RunE: func(*cobra.Command, []string) error { return a.showAgentCommand() },
	}
	var name string
	var attended int
	set := &cobra.Command{
		Use:   "set -- <program> [args...]",
		Short: "Set the command, as an argument list after --",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if strings.TrimSpace(args[0]) == "" {
				return fmt.Errorf("the program is empty")
			}
			cfg, path, err := a.loadConfig()
			if err != nil {
				return err
			}
			cfg.Agent.Command = args
			if name != "" {
				cfg.Agent.Name = name
			}
			if attended > 0 {
				cfg.Agent.AttendedMinutes = attended
			}
			if err := config.Save(path, cfg); err != nil {
				return err
			}
			return a.showAgentCommand()
		},
	}
	set.Flags().StringVar(&name, "name", "", "the FLAI_AGENT the session works under (default agent)")
	set.Flags().IntVar(&attended, "attended-minutes", 0, "how recent a sign of an agent counts as attending (default 6)")
	c.AddCommand(set,
		&cobra.Command{Use: "show", Short: "Print the command and whether the action is enabled here", Args: cobra.NoArgs,
			RunE: func(*cobra.Command, []string) error { return a.showAgentCommand() }},
		&cobra.Command{Use: "clear", Short: "Remove the command; nothing is started without one", Args: cobra.NoArgs,
			RunE: func(*cobra.Command, []string) error {
				cfg, path, err := a.loadConfig()
				if err != nil {
					return err
				}
				cfg.Agent = config.AgentStart{}
				if err := config.Save(path, cfg); err != nil {
					return err
				}
				return a.showAgentCommand()
			}},
	)
	return c
}

func (a *app) showAgentCommand() error {
	cfg, _, err := a.loadConfig()
	if err != nil {
		return err
	}
	here, enabled := "", false
	if repo, err := a.project(); err == nil {
		here = mainRootOf(repo)
		enabled = cfg.ActionEnabled(hostapi.ActionAgent, here)
	}
	if a.jsonOut {
		cmd := cfg.Agent.Command
		if cmd == nil {
			cmd = []string{}
		}
		return a.printJSON(map[string]any{"command": cmd, "name": cfg.Agent.Name, "attended_minutes": cfg.Agent.AttendedMinutes, "enabled_here": enabled})
	}
	if len(cfg.Agent.Command) == 0 {
		fmt.Fprintln(a.out, "no command is set, so nothing is started; flai serve agent set -- <program> [args...]")
	} else {
		quoted := make([]string, len(cfg.Agent.Command))
		for i, arg := range cfg.Agent.Command {
			quoted[i] = fmt.Sprintf("%q", arg)
		}
		fmt.Fprintf(a.out, "command: %s\n  run as it stands, in the project's directory, never through a shell\n", strings.Join(quoted, " "))
	}
	if here != "" {
		if enabled {
			fmt.Fprintln(a.out, "the agent action is on for this project; flai serve disable agent turns it off")
		} else {
			fmt.Fprintln(a.out, "the agent action is off for this project; flai serve enable agent turns it on")
		}
	}
	return nil
}

// agentConfig is the operator's say about starting agents for a project,
// read from the configuration at every look (S-0079).
func (a *app) agentConfig(root string) serve.AgentConfig {
	cfg, _, err := a.loadConfig()
	if err != nil {
		return serve.AgentConfig{}
	}
	return serve.AgentConfig{
		Enabled:  cfg.ActionEnabled(hostapi.ActionAgent, root),
		Command:  cfg.Agent.Command,
		Name:     cfg.Agent.Name,
		Attended: time.Duration(cfg.Agent.AttendedMinutes) * time.Minute,
	}
}
