package cmd

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/bytepunx/system-flow/flai/internal/config"
	"github.com/bytepunx/system-flow/flai/internal/harness"
	"github.com/bytepunx/system-flow/flai/internal/hostapi"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/serve"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// Host actions (S-0078, ADR-0029): what flai serve may do on this host, with
// the operator's credentials, because a dashboard asked. None is enabled
// until the operator enables it by name, here, in a shell on the host. The
// setting lives in the host's flai configuration; the journal, and the
// records that answer a repeated detached write (S-0109), lie beside flai
// serve's state. No method a dashboard can ask for touches either.

// host is the operator's say, read from the configuration at every question,
// so that enabling and disabling take effect without a restart.
func (a *app) host() hostapi.Host {
	return hostapi.Host{
		Enabled: func(action, root string) bool {
			cfg, _, err := a.loadConfig()
			return err == nil && cfg.ActionEnabled(action, root)
		},
		Agent: func(root string) any {
			cfg, _, _ := a.loadConfig()
			st := a.serveDir().AgentStates()[root]
			command := ""
			if len(cfg.Agent.Command) > 0 {
				command = filepath.Base(cfg.Agent.Command[0]) // its name, never its arguments
			}
			return map[string]any{"command": command, "running": st.Running, "last": st.Last, "waiting": st.Waiting, "stories": serve.Activity(root, st)}
		},
		Settings: a.hostSettings,
		Requests: a.serveDir().Requests,
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
		Short: "What flai serve starts when a story becomes ready: each story's harness, or your command",
		Long: `flai serve can start an agent for each story that becomes ready, while the
in-progress limit leaves room and no agent of yours is attending the project.
It is a host action, off until you enable it (flai serve enable agent).

A story names its agent (flai agent, flai edit --harness): a harness, a
model, and options. flai starts a harness it knows through its adapter
(claude-code runs claude -p with the story's model, and flai's MCP server),
with the program and the arguments you set for it here with flai serve agent
harness, which say what the agent may do. A story names only tunables the
adapter checks, never a program, a flag, or a permission.

A story that names no harness, or names "command", is started with your
command, when one is set: an argument list, run as it stands in the
project's directory, as you, never through a shell. In an argument {story},
{root}, {model}, and {harness} are replaced; nothing else is interpreted, and
the story's options are in FLAI_AGENT_CONFIG as JSON.

Every session's environment carries FLAI_AGENT (your --name and the story,
such as agent-S-0104), FLAI_STORY, and FLAI_SESSION. Its output goes to a log
beside flai serve's state, and every start and failure is in flai serve
journal.`,
		Example: `  flai serve agent harness claude-code
  flai serve agent harness claude-code --program /opt/claude/bin/claude -- --permission-mode acceptEdits --allowedTools Bash,mcp__flai
  flai serve agent set -- /home/me/bin/start-agent {story} --model {model}
  flai serve agent show
  flai serve enable agent
  flai serve agent clear`,
		Args: cobra.NoArgs,
		RunE: func(*cobra.Command, []string) error { return a.showAgentCommand() },
	}
	var name string
	var attended int
	set := &cobra.Command{
		Use:   "set [--name] [-- <program> [args...]]",
		Short: "Set the command, as an argument list after --, or only the name",
		Args:  cobra.ArbitraryArgs,
		RunE: func(_ *cobra.Command, args []string) error {
			if attended > 0 {
				fmt.Fprintln(a.errOut, "--attended-minutes is retired: nobody attending holds a ready story back any more (ADR-0043), so it does nothing")
				if len(args) == 0 && name == "" {
					return nil
				}
			}
			if len(args) == 0 && name == "" {
				return fmt.Errorf("nothing to set: give -- <program> [args...] or --name")
			}
			if len(args) > 0 && strings.TrimSpace(args[0]) == "" {
				return fmt.Errorf("the program is empty")
			}
			cfg, path, err := a.loadConfig()
			if err != nil {
				return err
			}
			if len(args) > 0 {
				cfg.Agent.Command = args
			}
			if name != "" {
				cfg.Agent.Name = name
			}
			if err := config.Save(path, cfg); err != nil {
				return err
			}
			return a.showAgentCommand()
		},
	}
	set.Flags().StringVar(&name, "name", "", "the FLAI_AGENT the session works under (default agent)")
	// Retired (S-0116, ADR-0043): accepted so that a script that sets it still runs.
	set.Flags().IntVar(&attended, "attended-minutes", 0, "retired: nobody attending holds a ready story back any more")
	_ = set.Flags().MarkHidden("attended-minutes")
	c.AddCommand(set, newServeAgentHarnessCmd(a), newServeAgentStartCmd(a), newServeAgentRestartCmd(a),
		&cobra.Command{Use: "show", Short: "Print the command and whether the action is enabled here", Args: cobra.NoArgs,
			RunE: func(*cobra.Command, []string) error { return a.showAgentCommand() }},
		&cobra.Command{Use: "clear", Short: "Remove the command; a story with a harness is still started with it", Args: cobra.NoArgs,
			RunE: func(*cobra.Command, []string) error {
				cfg, path, err := a.loadConfig()
				if err != nil {
					return err
				}
				cfg.Agent = config.AgentStart{Harnesses: cfg.Agent.Harnesses}
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
		hosts := map[string]harness.Host{}
		for _, name := range settable() {
			hosts[name] = a.harnessHost(cfg, name)
		}
		return a.printJSON(map[string]any{"command": cmd, "name": cfg.Agent.Name, "harnesses": hosts, "enabled_here": enabled})
	}
	if len(cfg.Agent.Command) == 0 {
		fmt.Fprintln(a.out, "no command is set, so a story that names no harness is not started; flai serve agent set -- <program> [args...]")
	} else {
		fmt.Fprintf(a.out, "command: %s\n  run as it stands, in the project's directory, never through a shell\n", quoteArgs(cfg.Agent.Command))
	}
	for _, name := range settable() {
		h := a.harnessHost(cfg, name)
		fmt.Fprintf(a.out, "harness %s: %s\n", name, quoteArgs(append([]string{h.Program}, h.Args...)))
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

// checksConfig is the operator's say about running checks for a story in
// review, for repo: the host's own named commands when it names any, else
// the manifest's, read fresh at every question (S-0082).
func (a *app) checksConfig(repo *workitem.Repo) (serve.ChecksConfig, error) {
	cfg, _, err := a.loadConfig()
	if err != nil {
		return serve.ChecksConfig{}, err
	}
	host := make([]manifest.NamedCommand, len(cfg.Checks.Commands))
	for i, nc := range cfg.Checks.Commands {
		host[i] = manifest.NamedCommand{Name: nc.Name, Command: nc.Command}
	}
	timeout := serve.DefaultChecksTimeout
	if cfg.Checks.TimeoutMinutes > 0 {
		timeout = time.Duration(cfg.Checks.TimeoutMinutes) * time.Minute
	}
	return serve.ChecksConfig{
		Enabled:  cfg.ActionEnabled(hostapi.ActionChecks, mainRootOf(repo)),
		Commands: serve.ResolveChecks(host, repo.Manifest.Checks),
		Timeout:  timeout,
	}, nil
}

// agentConfig is the operator's say about starting agents for a project,
// read from the configuration at every look (S-0079).
func (a *app) agentConfig(root string) serve.AgentConfig {
	cfg, _, err := a.loadConfig()
	if err != nil {
		return serve.AgentConfig{}
	}
	hosts := map[string]harness.Host{}
	for _, name := range settable() {
		hosts[name] = a.harnessHost(cfg, name)
	}
	if len(cfg.Agent.Command) > 0 {
		hosts[harness.Command] = harness.Host{Program: cfg.Agent.Command[0], Args: cfg.Agent.Command[1:]}
	}
	self, _ := os.Executable()
	return serve.AgentConfig{
		Flai:      self,
		Enabled:   cfg.ActionEnabled(hostapi.ActionAgent, root),
		Command:   cfg.Agent.Command,
		Harnesses: hosts,
		Name:      cfg.Agent.Name,
	}
}

// flai serve checks: the named commands a story in review is checked with
// (S-0082). Unlike agent, there can be several, so set names one at a time
// and clear takes an optional name; show lists all of them.
func newServeChecksCmd(a *app) *cobra.Command {
	c := &cobra.Command{
		Use:   "checks",
		Short: "The named commands a story in review is checked with, on this host",
		Long: `flai serve can run checks for a story in review, in the story's worktree, and
keep the outcome with it until it is accepted. It is a host action, off
until you enable it (flai serve enable checks), and by default nothing is
named here: the manifest's own checks: is used instead. Naming any command
here, on this host, uses this list instead of the manifest's, for this host
only.

Each command is an argument list, run as it stands in the worktree, as you,
never through a shell. In an argument {story} is replaced by the story's ID
and {root} by the worktree's directory; nothing else is interpreted. Every
named command runs in order; the first to fail stops the rest. One run per
story at a time; every run and every refusal is in flai serve journal.`,
		Example: `  flai serve checks set --name flai -- scripts/flai-test.sh
  flai serve checks set --name flaiover -- bash -c "cd flaiover && pnpm run check && pnpm run test:unit -- run && pnpm run lint"
  flai serve checks show
  flai serve enable checks
  flai serve checks clear flaiover
  flai serve checks timeout 20`,
		Args: cobra.NoArgs,
		RunE: func(*cobra.Command, []string) error { return a.showChecksCommands() },
	}
	var name string
	set := &cobra.Command{
		Use:   "set --name <name> -- <program> [args...]",
		Short: "Add or replace one named command, as an argument list after --",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if strings.TrimSpace(name) == "" {
				return fmt.Errorf("--name is required")
			}
			if strings.TrimSpace(args[0]) == "" {
				return fmt.Errorf("the program is empty")
			}
			cfg, path, err := a.loadConfig()
			if err != nil {
				return err
			}
			replaced := false
			for i, nc := range cfg.Checks.Commands {
				if nc.Name == name {
					cfg.Checks.Commands[i].Command = args
					replaced = true
					break
				}
			}
			if !replaced {
				cfg.Checks.Commands = append(cfg.Checks.Commands, config.NamedCommand{Name: name, Command: args})
			}
			if err := config.Save(path, cfg); err != nil {
				return err
			}
			return a.showChecksCommands()
		},
	}
	set.Flags().StringVar(&name, "name", "", "the check's name (required)")
	timeout := &cobra.Command{
		Use:   "timeout <minutes>",
		Short: "Bound how long one run of every command together may take (default 15)",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			minutes, err := strconv.Atoi(args[0])
			if err != nil || minutes <= 0 {
				return fmt.Errorf("timeout must be a positive number of minutes, got %q", args[0])
			}
			cfg, path, err := a.loadConfig()
			if err != nil {
				return err
			}
			cfg.Checks.TimeoutMinutes = minutes
			if err := config.Save(path, cfg); err != nil {
				return err
			}
			return a.showChecksCommands()
		},
	}
	clear := &cobra.Command{
		Use:   "clear [name]",
		Short: "Remove one named command, or every one when no name is given",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			cfg, path, err := a.loadConfig()
			if err != nil {
				return err
			}
			if len(args) == 0 {
				cfg.Checks.Commands = nil
			} else {
				kept := cfg.Checks.Commands[:0]
				for _, nc := range cfg.Checks.Commands {
					if nc.Name != args[0] {
						kept = append(kept, nc)
					}
				}
				cfg.Checks.Commands = kept
			}
			if err := config.Save(path, cfg); err != nil {
				return err
			}
			return a.showChecksCommands()
		},
	}
	c.AddCommand(set, timeout, clear,
		&cobra.Command{Use: "show", Short: "Print the named commands and whether the action is enabled here", Args: cobra.NoArgs,
			RunE: func(*cobra.Command, []string) error { return a.showChecksCommands() }},
	)
	return c
}

func (a *app) showChecksCommands() error {
	cfg, _, err := a.loadConfig()
	if err != nil {
		return err
	}
	here, enabled := "", false
	if repo, err := a.project(); err == nil {
		here = mainRootOf(repo)
		enabled = cfg.ActionEnabled(hostapi.ActionChecks, here)
	}
	timeout := cfg.Checks.TimeoutMinutes
	if timeout <= 0 {
		timeout = 15
	}
	if a.jsonOut {
		cmds := cfg.Checks.Commands
		if cmds == nil {
			cmds = []config.NamedCommand{}
		}
		return a.printJSON(map[string]any{"commands": cmds, "timeout_minutes": timeout, "enabled_here": enabled})
	}
	if len(cfg.Checks.Commands) == 0 {
		fmt.Fprintln(a.out, "no command is named here; the manifest's own checks: is used instead")
	} else {
		for _, nc := range cfg.Checks.Commands {
			quoted := make([]string, len(nc.Command))
			for i, arg := range nc.Command {
				quoted[i] = fmt.Sprintf("%q", arg)
			}
			fmt.Fprintf(a.out, "%s: %s\n", nc.Name, strings.Join(quoted, " "))
		}
	}
	fmt.Fprintf(a.out, "timeout: %d minute%s\n", timeout, map[bool]string{true: "", false: "s"}[timeout == 1])
	if here != "" {
		if enabled {
			fmt.Fprintln(a.out, "the checks action is on for this project; flai serve disable checks turns it off")
		} else {
			fmt.Fprintln(a.out, "the checks action is off for this project; flai serve enable checks turns it on")
		}
	}
	return nil
}

// harnessHost is the program and arguments a harness runs with here: the
// operator's where set, the adapter's defaults where not (S-0104).
func (a *app) harnessHost(cfg config.Config, name string) harness.Host {
	h := harness.Adapters[name].DefaultHost()
	set := cfg.Agent.Harnesses[name]
	if set.Program != "" {
		h.Program = set.Program
	}
	if set.Args != nil {
		h.Args = *set.Args
	}
	return h
}

func quoteArgs(args []string) string {
	quoted := make([]string, len(args))
	for i, arg := range args {
		quoted[i] = fmt.Sprintf("%q", arg)
	}
	return strings.Join(quoted, " ")
}

// flai serve agent harness: the operator's program and arguments for a
// harness a story may name (S-0104). Arguments after -- replace the
// adapter's defaults, even with none; --reset goes back to the defaults.
func newServeAgentHarnessCmd(a *app) *cobra.Command {
	var program string
	var reset bool
	c := &cobra.Command{
		Use:   "harness [<name>] [--program <path>] [--reset] [-- args...]",
		Short: "Set the program a harness is and the arguments that say what its agent may do",
		Long: `A story names its harness, such as claude-code; you say here what that is on
this host. --program is what is run, and the arguments after -- replace the
adapter's defaults, which say what the agent may do (for claude-code:
--permission-mode acceptEdits --allowedTools Bash,mcp__flai). "--" with
nothing after it runs it with none. --reset goes back to the defaults. With
no name, every harness is shown.`,
		Example: `  flai serve agent harness
  flai serve agent harness claude-code --program /opt/claude/bin/claude
  flai serve agent harness claude-code -- --permission-mode acceptEdits --allowedTools Bash,Edit,mcp__flai
  flai serve agent harness claude-code --reset`,
		RunE: func(cmd *cobra.Command, args []string) error {
			dash := cmd.ArgsLenAtDash()
			names, rest := args, []string(nil)
			if dash >= 0 {
				names, rest = args[:dash], args[dash:]
			}
			if len(names) == 0 {
				if dash >= 0 || program != "" || reset {
					return fmt.Errorf("name the harness: %s", strings.Join(settable(), ", "))
				}
				return a.showAgentCommand()
			}
			if len(names) > 1 {
				return fmt.Errorf("one harness at a time; arguments for it go after --")
			}
			name := names[0]
			if !slices.Contains(settable(), name) {
				return fmt.Errorf("there is no harness %q to set; there are: %s (the command is set with flai serve agent set)", name, strings.Join(settable(), ", "))
			}
			cfg, path, err := a.loadConfig()
			if err != nil {
				return err
			}
			if cfg.Agent.Harnesses == nil {
				cfg.Agent.Harnesses = map[string]config.HarnessHost{}
			}
			h := cfg.Agent.Harnesses[name]
			if reset {
				if program != "" || dash >= 0 {
					return fmt.Errorf("--reset takes neither --program nor arguments")
				}
				h = config.HarnessHost{}
			}
			if program != "" {
				h.Program = program
			}
			if dash >= 0 {
				list := append([]string{}, rest...)
				h.Args = &list
			}
			if h.Program == "" && h.Args == nil {
				delete(cfg.Agent.Harnesses, name)
			} else {
				cfg.Agent.Harnesses[name] = h
			}
			if err := config.Save(path, cfg); err != nil {
				return err
			}
			return a.showAgentCommand()
		},
	}
	c.Flags().StringVar(&program, "program", "", "the program the harness is on this host")
	c.Flags().BoolVar(&reset, "reset", false, "go back to the adapter's program and arguments")
	return c
}

// settable are the harnesses the operator sets with flai serve agent harness:
// all but the command, which flai serve agent set sets.
func settable() []string {
	var out []string
	for _, n := range harness.Names() {
		if n != harness.Command {
			out = append(out, n)
		}
	}
	return out
}

// hostSettings is the host's configuration as it applies to the project at
// root, for the dashboard's settings page (S-0105). The tokens are not in
// it: a rotation answers with the new one, once.
func (a *app) hostSettings(root string) any {
	cfg, _, err := a.loadConfig()
	if err != nil {
		return map[string]any{"error": err.Error()}
	}
	names := make([]string, 0, len(hostapi.Actions))
	for n := range hostapi.Actions {
		names = append(names, n)
	}
	sort.Strings(names)
	actions := make([]map[string]any, 0, len(names))
	for _, n := range names {
		actions = append(actions, map[string]any{"name": n, "means": hostapi.Actions[n],
			"here": cfg.ActionEnabled(n, root), "everywhere": cfg.ActionEnabled(n, config.AllProjects)})
	}
	harnesses := map[string]any{}
	for _, n := range settable() {
		h := a.harnessHost(cfg, n)
		_, set := cfg.Agent.Harnesses[n]
		harnesses[n] = map[string]any{"program": h.Program, "args": h.Args, "set": set}
	}
	command := cfg.Agent.Command
	if command == nil {
		command = []string{}
	}
	checks := cfg.Checks.Commands
	if checks == nil {
		checks = []config.NamedCommand{}
	}
	timeout := cfg.Checks.TimeoutMinutes
	if timeout <= 0 {
		timeout = int(serve.DefaultChecksTimeout / time.Minute)
	}
	roots := cfg.ImportRoots
	if roots == nil {
		roots = []string{}
	}
	out := map[string]any{
		"actions": actions,
		"agent": map[string]any{"name": cfg.Agent.Name, "command": command,
			"harnesses": harnesses},
		"checks":       map[string]any{"commands": checks, "timeout_minutes": timeout},
		"import_roots": roots,
	}
	if repo, err := workitem.Open(root); err == nil {
		out["default_agent"] = repo.Manifest.Agent
		out["manifest_checks"] = repo.Manifest.Checks
		if st, running := readMCPState(repo, a.now()); running {
			out["mcp"] = map[string]any{"running": true, "url": st.URL, "pid": st.PID}
		} else {
			out["mcp"] = map[string]any{"running": false}
		}
	}
	return out
}

func newServeAgentRestartCmd(a *app) *cobra.Command {
	return &cobra.Command{
		Use:   "restart <story-id>",
		Short: "Start a new agent for a story in ready or in progress whose agent dropped or failed",
		Long: `Starts a new agent, in a new session, for a story in ready or
in-progress whose last agent flai serve started has ended or dropped, the
way flai serve starts one when a story enters ready: the story's harness,
model, and options, or the host's command. The agent is told how its last
one ended and to go on from the story's narrative. The run is recorded
where the serving flai tracks it: the dot on the card, the outcome, the
restart on an answer (S-0116, ADR-0043).

It refuses, and says why, while the agent action is off for the project,
when the story is in another state, when flai serve has started no agent for
it, while its agent runs or waits for an answer, when nothing can start it,
and for a story in ready while the in-progress limit is full. The story
page's Restart agent button runs this.`,
		Args: cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			return a.agentNow(args[0], serve.Restart)
		},
	}
}

func newServeAgentStartCmd(a *app) *cobra.Command {
	return &cobra.Command{
		Use:   "start <story-id>",
		Short: "Start a ready story's agent now, whatever flai serve's own rules say about when",
		Long: `Starts the agent of a story in ready now, the way flai serve starts one
when a story enters ready: the story's harness, model, and options, or the
host's command. It starts whether or not the story was ready before flai
serve began, whether or not it has had an agent since it entered ready, and
whoever is attending the project. The run is recorded where the serving
flai tracks it: the dot on the card, the outcome when it ends, the start
again on an answer (S-0115).

It refuses, and says why, while the agent action is off for the project,
when the story is not in ready, while its agent runs or waits for an
answer, and when nothing can start it (no harness and no command). A full
in-progress limit does not stop it: it starts past the limit, with a
warning, as a move does. A story in progress whose agent dropped or failed
gets a new one from flai serve agent restart. The story page's Start agent
button runs this.`,
		Example: `  flai serve agent start S-0115`,
		Args:    cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			return a.agentNow(args[0], serve.Start)
		},
	}
}

// agentNow starts story's agent on the operator's word, with start or
// restart deciding whether it may, and prints the run.
func (a *app) agentNow(story string, how func(context.Context, serve.Options, serve.Entry, string) (*serve.AgentRun, error)) error {
	repo, err := a.project()
	if err != nil {
		return err
	}
	o := serve.Options{Dir: a.serveDir(), Logger: a.logger(), Now: a.now, Agent: a.agentConfig, Host: a.host()}
	e := serve.Entry{Key: repo.Manifest.Key, Name: repo.Manifest.Name, Root: mainRootOf(repo)}
	run, err := how(context.Background(), o, e, workitem.CanonicalID(story))
	var no *serve.Refused
	if errors.As(err, &no) {
		return fmt.Errorf("rule: %s", no.Why)
	}
	if err != nil {
		return err
	}
	if a.jsonOut {
		return a.printJSON(map[string]any{"story": run.Story, "agent": run.Agent, "harness": run.Harness, "command": run.Command, "pid": run.PID, "log": run.Log, "started": run.Started})
	}
	fmt.Fprintf(a.out, "started %s (%s) for %s as %s (pid %d); log %s\n", run.Command, run.Harness, run.Story, run.Agent, run.PID, run.Log)
	return nil
}
