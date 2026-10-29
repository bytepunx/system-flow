package cmd

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/bytepunx/system-flow/flai/internal/config"
	"github.com/bytepunx/system-flow/flai/internal/execx"
	"github.com/bytepunx/system-flow/flai/internal/harness"
	"github.com/bytepunx/system-flow/flai/internal/hostapi"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/mcpserver"
	"github.com/bytepunx/system-flow/flai/internal/selfupgrade"
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
			plans := st.Plans // the newest planner run for each item, by its ID (S-0208)
			if plans == nil {
				plans = map[string]*serve.AgentRun{}
			}
			// the newest run of each role the project defines, by its name and
			// then its item, "" for none (S-0417)
			defined := st.Roles
			if defined == nil {
				defined = map[string]map[string]*serve.AgentRun{}
			}
			// the newest orchestrator run, null before the first (S-0218), with
			// held while the operator holds it stopped (S-0228), and the newest
			// analyzer run, its focus and report among it (S-0223), and the
			// newest curator run, its since and entries among it (S-0431)
			out := map[string]any{"command": command, "running": st.Running, "last": st.Last, "waiting": st.Waiting, "stories": serve.Activity(root, st), "plans": plans,
				"orchestrator": st.Orchestrator, "analyzer": st.Analyzer, "curator": st.Curator, "roles": defined}
			if budget := a.hostBudget(); len(budget) > 0 {
				out["budget"] = budget // the host's, the same for every project (S-0534)
			}
			return out
		},
		AgentStream: func(root, story string, after int64) (any, error) {
			return serve.Stream(a.serveDir().AgentStates()[root], story, after)
		},
		PlanStream: func(root, item string, after int64) (any, error) {
			return serve.PlanStream(a.serveDir().AgentStates()[root], item, after)
		},
		RoleStream: func(root, role string, after int64) (any, error) {
			return serve.RoleStream(a.serveDir().AgentStates()[root], role, after)
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

// hostBudget is the host's budget state (S-0534, ADR-0209 §5), read from
// budget.json in the serve folder: one entry per provider whose state is
// set, by provider, none when no file is there. A file that cannot be read
// is warned of and taken as none, so that the rest of the answer stands.
func (a *app) hostBudget() []hostapi.Budget {
	states, err := a.serveDir().Budgets()
	if err != nil {
		a.logger().Warn("the host's budget state not read; taken as none", "component", "serve", "err", err.Error())
		return nil
	}
	out := make([]hostapi.Budget, 0, len(states))
	for _, b := range states {
		out = append(out, hostapi.Budget{Provider: b.Provider, Limit: b.Limit, Status: b.Status,
			Since: b.Since.UTC().Format(time.RFC3339), Until: b.Until.UTC().Format(time.RFC3339), Probing: b.Probe != nil})
	}
	return out
}

// budgetLine is how flai serve status names a provider's budget state at
// now: held, with the limit the harness named, until when; and, once until
// has passed, that a probe runs, or that the next start probes.
func budgetLine(b hostapi.Budget, now time.Time) string {
	line := "budget: " + b.Provider + " held"
	switch {
	case b.Limit != "":
		line += " (" + b.Limit + ")"
	case b.Status != 0:
		line += " (HTTP " + strconv.Itoa(b.Status) + ")"
	}
	line += " until " + b.Until
	if b.Probing {
		return line + ", probing"
	}
	if until, err := time.Parse(time.RFC3339, b.Until); err == nil && !now.Before(until) {
		return line + ", passed; the next start probes"
	}
	return line
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

// hostUpgrade is how flai serve asks for a newer flai (ADR-0158): of the
// host its MCP servers are kept by, under the same check (S-0183), and of
// none when it keeps none.
func hostUpgrade(m serve.MCP) serve.HostUpgrade {
	h, ok := m.(hostMCP)
	if !ok {
		return nil
	}
	return h.upgrade
}

// upgrade asks the host to install the flai release version, as flai host
// upgrade --version does, and answers whether it restarts on it.
func (m hostMCP) upgrade(ctx context.Context, version string) (bool, error) {
	out, err := m.c.Upgrade(ctx, version)
	if err != nil {
		return false, err
	}
	var said struct {
		Restarting bool `json:"restarting"`
	}
	if err := json.Unmarshal(out, &said); err != nil {
		return false, fmt.Errorf("flai host answered the upgrade to flai %s with %s: %w", version, out, err)
	}
	return said.Restarting, nil
}

// releasePublished is how flai serve learns whether a flai release is
// published before it names it on a thread (ADR-0161): among the releases
// flai self-upgrade --list reads, from FLAI_RELEASES_API when set.
func releasePublished(r execx.Runner) serve.ReleasePublished {
	return func(ctx context.Context, version string) (bool, error) {
		published, err := selfupgrade.List(ctx, selfupgrade.Options{APIBase: strings.TrimSpace(os.Getenv("FLAI_RELEASES_API")), Token: selfupgrade.Token(r)}, selfupgrade.TagPrefix)
		if err != nil {
			return false, err
		}
		return slices.ContainsFunc(published, func(rel selfupgrade.Release) bool { return rel.Version == version }), nil
	}
}

// hostActions are the host actions there are for the projects at roots,
// each with what it lets a dashboard do or has flai do: hostapi.Actions, and
// those the roles each project defines name (S-0417).
func hostActions(roots ...string) map[string]string {
	out := maps.Clone(hostapi.Actions)
	for _, root := range roots {
		for n, means := range hostapi.DefinedActions(root) {
			out[n] = means
		}
	}
	return out
}

func knownAction(name string, actions map[string]string) error {
	if _, ok := actions[name]; ok {
		return nil
	}
	names := slices.Sorted(maps.Keys(actions))
	return fmt.Errorf("there is no host action %q; there %s: %s", name, oneOrMany(len(names), "is", "are"), strings.Join(names, ", "))
}

// actionMeans says what the host action n of actions does once enabled:
// what it lets a dashboard do, or, for one no dashboard sees, what it has
// flai do.
func actionMeans(n string, actions map[string]string) string {
	if dashboardSees(n) {
		return "lets a dashboard " + actions[n]
	}
	return actions[n]
}

// dashboardSees reports whether a dashboard is told of the host action n:
// one of hostapi.Actions it sees, or one a role names, which role.run asks
// for (S-0417).
func dashboardSees(n string) bool {
	_, fixed := hostapi.Actions[n]
	return !fixed || hostapi.DashboardSees(n)
}

// servedRoots are the projects flai serve serves on this host, by their
// main checkouts, and the one in the working directory when there is one,
// whose role files name the host actions --all-projects may turn on
// (S-0417). A registry that cannot be read names none.
func (a *app) servedRoots() []string {
	var roots []string
	if repo, err := a.project(); err == nil {
		roots = append(roots, mainRootOf(repo))
	}
	served, _ := a.serveDir().Projects()
	for _, e := range served {
		roots = append(roots, e.Root)
	}
	return roots
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
			root, where := config.AllProjects, "every project"
			var roots []string
			if all {
				roots = a.servedRoots()
			} else {
				repo, err := a.project()
				if err != nil {
					if unknown := knownAction(action, hostActions()); unknown != nil {
						return unknown
					}
					return fmt.Errorf("%w; or name every project with --all-projects", err)
				}
				root, where = mainRootOf(repo), repo.Manifest.Name
				roots = []string{root}
			}
			cfg, path, err := a.loadConfig()
			if err != nil {
				return err
			}
			// an action a role named before its file went is still turned off
			actions := hostActions(roots...)
			if err := knownAction(action, actions); err != nil && (on || len(cfg.HostActions[action]) == 0) {
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
				fmt.Fprintf(a.out, "%s enabled for %s\n  it %s\n", action, where, actionMeans(action, actions))
				if dashboardSees(action) {
					fmt.Fprintln(a.out, "  every use is journalled: flai serve journal")
				}
				fmt.Fprintf(a.out, "  turn it off with: flai serve disable %s\n", action)
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
			actions := hostActions()
			if here != "" {
				actions = hostActions(here) // and those the roles it defines name (S-0417)
			}
			names := slices.Sorted(maps.Keys(actions))
			if a.jsonOut {
				out := []map[string]any{}
				for _, n := range names {
					roots := cfg.HostActions[n]
					if roots == nil {
						roots = []string{}
					}
					out = append(out, map[string]any{"action": n, "means": actions[n], "enabled_for": roots, "enabled_here": here != "" && cfg.ActionEnabled(n, here)})
				}
				return a.printJSON(out)
			}
			for _, n := range names {
				state := "off everywhere"
				if roots := cfg.HostActions[n]; len(roots) > 0 {
					state = "on for " + strings.ReplaceAll(strings.Join(roots, ", "), config.AllProjects, "every project")
				}
				fmt.Fprintf(a.out, "%s: %s\n  %s\n", n, state, actionMeans(n, actions))
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
in-progress limit leaves room, review is under its limit, and no agent of yours is attending the project.
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
	var attended, restarts int
	set := &cobra.Command{
		Use:   "set [--name] [--auto-restarts <n>] [-- <program> [args...]]",
		Short: "Set the command, as an argument list after --, the name, or how many times a story's agent is restarted on its own",
		Long: `Sets what you give and leaves the rest as it was: the command, as an
argument list after --; the name the sessions work under; and, with
--auto-restarts, how many times flai serve restarts a story's agent on its own
after it ends with its story in progress, 2 when unset and 0 for never
(ADR-0108). Past that, flai serve opens a thread on the story for you.`,
		Args: cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			restartsSet := cmd.Flags().Changed("auto-restarts")
			if attended > 0 {
				fmt.Fprintln(a.errOut, "--attended-minutes is retired: nobody attending holds a ready story back any more (ADR-0043), so it does nothing")
				if len(args) == 0 && name == "" && !restartsSet {
					return nil
				}
			}
			if len(args) == 0 && name == "" && !restartsSet {
				return fmt.Errorf("nothing to set: give -- <program> [args...], --name, or --auto-restarts")
			}
			if len(args) > 0 && strings.TrimSpace(args[0]) == "" {
				return fmt.Errorf("the program is empty")
			}
			cfg, path, err := a.loadConfig()
			if err != nil {
				return err
			}
			if restartsSet {
				if cfg.Agent, err = cfg.Agent.WithAutoRestarts(restarts); err != nil {
					return err
				}
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
	set.Flags().IntVar(&restarts, "auto-restarts", config.DefaultAutoRestarts, "how many times flai serve restarts a story's agent on its own after it ends with its story in progress; 0 for never")
	// Retired (S-0116, ADR-0043): accepted so that a script that sets it still runs.
	set.Flags().IntVar(&attended, "attended-minutes", 0, "retired: nobody attending holds a ready story back any more")
	_ = set.Flags().MarkHidden("attended-minutes")
	c.AddCommand(set, newServeAgentHarnessCmd(a), newServeAgentProviderCmd(a), newServeAgentStartCmd(a), newServeAgentRestartCmd(a), newServeAgentCommitCmd(a), newServeAgentStopCmd(a), newServeAgentStreamCmd(a), newServeAgentUsageCmd(a),
		&cobra.Command{Use: "show", Short: "Print the command and whether the action is enabled here", Args: cobra.NoArgs,
			RunE: func(*cobra.Command, []string) error { return a.showAgentCommand() }},
		&cobra.Command{Use: "clear", Short: "Remove the command and the name; the harnesses, the providers, and the restart limit stay", Args: cobra.NoArgs,
			RunE: func(*cobra.Command, []string) error {
				cfg, path, err := a.loadConfig()
				if err != nil {
					return err
				}
				cfg.Agent = config.AgentStart{Harnesses: cfg.Agent.Harnesses, Providers: cfg.Agent.Providers, AutoRestarts: cfg.Agent.AutoRestarts}
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
		type shown struct {
			harness.Host
			Guard         string `json:"guard,omitempty"`
			DenyProtected bool   `json:"deny_protected,omitempty"`
		}
		hosts := map[string]shown{}
		for _, name := range settable() {
			set := cfg.Agent.Harnesses[name]
			hosts[name] = shown{a.harnessHost(cfg, name), set.Guard, set.DenyProtected}
		}
		if set, ok := cfg.Agent.Harnesses[harness.Command]; ok {
			var h harness.Host
			if len(cfg.Agent.Command) > 0 {
				h = harness.Host{Program: cfg.Agent.Command[0], Args: cfg.Agent.Command[1:]}
			}
			hosts[harness.Command] = shown{h, set.Guard, set.DenyProtected}
		}
		return a.printJSON(map[string]any{"command": cmd, "name": cfg.Agent.Name, "harnesses": hosts, "auto_restarts": cfg.Agent.AutoRestartLimit(), "enabled_here": enabled})
	}
	if len(cfg.Agent.Command) == 0 {
		fmt.Fprintln(a.out, "no command is set, so a story that names no harness is not started; flai serve agent set -- <program> [args...]")
	} else {
		fmt.Fprintf(a.out, "command: %s\n  run as it stands, in the project's directory, never through a shell\n", quoteArgs(cfg.Agent.Command))
	}
	if set, ok := cfg.Agent.Harnesses[harness.Command]; ok {
		fmt.Fprintf(a.out, "harness %s: the command above\n", harness.Command)
		printLifts(a, set)
	}
	for _, name := range settable() {
		h := a.harnessHost(cfg, name)
		fmt.Fprintf(a.out, "harness %s: %s\n", name, quoteArgs(append([]string{h.Program}, h.Args...)))
		printLifts(a, cfg.Agent.Harnesses[name])
	}
	fmt.Fprintf(a.out, "automatic restarts: %s\n  of a story's agent that ends with its story in progress; set with --auto-restarts\n", autoRestartsSaid(cfg.Agent))
	if here != "" {
		if enabled {
			fmt.Fprintln(a.out, "the agent action is on for this project; flai serve disable agent turns it off")
		} else {
			fmt.Fprintln(a.out, "the agent action is off for this project; flai serve enable agent turns it on")
		}
	}
	return nil
}

// printLifts prints what the operator lifts on a harness, under its line in
// show: guard none and deny_protected, each with what it gives up (ADR-0131).
func printLifts(a *app, set config.HarnessHost) {
	if set.Guard == config.GuardNone {
		fmt.Fprintln(a.out, "  guard: none, so a story's sub-agents run without flai guard where the harness cannot run it, and could move items, write threads, or commit")
	}
	if set.DenyProtected {
		fmt.Fprintln(a.out, "  deny_protected: your arguments deny the protected paths, so flai starts the agent where the harness cannot ask before a protected write")
	}
}

// autoRestartsSaid is the restart limit as show prints it: the count, and
// whether it is the default or off.
func autoRestartsSaid(agent config.AgentStart) string {
	n := agent.AutoRestartLimit()
	switch {
	case agent.AutoRestarts == nil:
		return fmt.Sprintf("%d (default)", n)
	case n == 0:
		return "0 (off)"
	}
	return strconv.Itoa(n)
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
	lifts := map[string]serve.HarnessLifts{}
	for name, set := range cfg.Agent.Harnesses {
		lifts[name] = serve.HarnessLifts{GuardNone: set.Guard == config.GuardNone, DenyProtected: set.DenyProtected}
	}
	var forgeTokens []string
	for _, f := range config.KnownForges() {
		forgeTokens = append(forgeTokens, cfg.ForgeTokenEnv(f))
	}
	// whether each host action a role the project defines names is on (S-0417)
	defined := map[string]bool{}
	for action := range hostapi.DefinedActions(root) {
		defined[action] = cfg.ActionEnabled(action, root)
	}
	self, _ := os.Executable()
	return serve.AgentConfig{
		Flai:           self,
		Enabled:        cfg.ActionEnabled(hostapi.ActionAgent, root),
		Plan:           cfg.ActionEnabled(hostapi.ActionPlan, root),
		Orchestrate:    cfg.ActionEnabled(hostapi.ActionOrchestrate, root),
		Analyze:        cfg.ActionEnabled(hostapi.ActionAnalyze, root),
		Curate:         cfg.ActionEnabled(hostapi.ActionCurate, root),
		Command:        cfg.Agent.Command,
		Harnesses:      hosts,
		Providers:      cfg.Agent.Providers,
		Name:           cfg.Agent.Name,
		AutoRestarts:   cfg.Agent.AutoRestartLimit(),
		Lifts:          lifts,
		AutoApprove:    cfg.ActionEnabled(hostapi.ActionAutoApprove, root),
		ForgeTokenEnvs: forgeTokens,
		RoleActions:    defined,
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

// serveHost is the host flai serve itself answers with: a.host() and the
// supervisor of the component dev servers (S-0438, ADR-0163), which runs at
// most components.max_servers per project and stops each one idle for
// components.idle, and each of a story's worktree once the story leaves
// review (S-0442), until ctx ends, and then every one.
// stop ends it as ctx ending does and waits until every dev server has
// stopped; flai serve calls it as it stops.
func (a *app) serveHost(ctx context.Context) (host hostapi.Host, stop func()) {
	comps := serve.NewComponents(serve.ComponentsOptions{Idle: a.componentsIdle, MaxServers: a.componentsMaxServers, Now: a.now, Logger: a.logger()})
	runCtx, cancel := context.WithCancel(ctx)
	done := make(chan struct{})
	go func() {
		defer close(done)
		if err := comps.Run(runCtx); err != nil {
			a.logger().Warn("component dev servers not all stopped", "component", "components", "err", err.Error())
		}
	}()
	host = a.host()
	host.Components = comps
	return host, func() {
		cancel()
		<-done
	}
}

// componentsIdle is the operator's components.idle, read from the
// configuration at every look, so that a change takes effect without a
// restart; serve.DefaultComponentsIdle when it is not set or not readable.
func (a *app) componentsIdle() time.Duration {
	cfg, _, err := a.loadConfig()
	if err != nil {
		return serve.DefaultComponentsIdle
	}
	idle, _ := componentsIdleOf(cfg)
	return idle
}

// componentsIdleOf is cfg's components.idle, and whether it is set; the
// default when it is not, or when it is not a positive duration.
func componentsIdleOf(cfg config.Config) (time.Duration, bool) {
	d, err := time.ParseDuration(strings.TrimSpace(cfg.Components.Idle))
	if err != nil || d <= 0 {
		return serve.DefaultComponentsIdle, false
	}
	return d, true
}

// componentsMaxServers is the operator's components.max_servers, read from
// the configuration at every start, so that a change takes effect without a
// restart; serve.DefaultComponentsMaxServers when it is not set or not
// readable.
func (a *app) componentsMaxServers() int {
	cfg, _, err := a.loadConfig()
	if err != nil {
		return serve.DefaultComponentsMaxServers
	}
	n, _ := componentsMaxServersOf(cfg)
	return n
}

// componentsMaxServersOf is cfg's components.max_servers, and whether it is
// set; the default when it is not, or when it is not positive.
func componentsMaxServersOf(cfg config.Config) (int, bool) {
	if cfg.Components.MaxServers <= 0 {
		return serve.DefaultComponentsMaxServers, false
	}
	return cfg.Components.MaxServers, true
}

// flai serve components: the host's say about the component dev servers
// flai serve runs under the components host action (S-0438).
func newServeComponentsCmd(a *app) *cobra.Command {
	c := &cobra.Command{
		Use:   "components",
		Short: "How long flai serve keeps a component dev server nobody looks at, and how many it runs, on this host",
		Long: `flai serve runs a Svelte sub-project's own Vite dev server, with flai's
component harness, when the dashboard's Components page asks it to
(components.start), on 127.0.0.1 and a port flai picks. It runs the
project's code on this host, so it is a host action, off until you enable it
(flai serve enable components). Each start and stop is in flai serve journal.

flai serve stops a dev server that had no components.status, which the page
asks while it is open, for components.idle, 30 minutes by default, and every
dev server when it stops. flai serve components idle sets it, as a Go
duration such as 45m or 2h, for every project on this host; it takes effect
at once, with no restart.

The page also runs a dev server in the worktree of a story in progress or in
review, beside the main checkout's, on a port of its own (S-0442); flai serve
stops it once the story leaves review, to done, cancelled, or back to in
progress, or is no longer in progress or in review at all. At most
components.max_servers dev servers run or start at once per project, main's
and every story's together, 3 by default; a start beyond it is refused,
naming the ones running. flai serve components max-servers sets it, for
every project on this host, at once.`,
		Example: `  flai serve enable components
  flai serve components idle 45m
  flai serve components max-servers 4
  flai serve components show --json`,
		Args: cobra.NoArgs,
		RunE: func(*cobra.Command, []string) error { return a.showComponentsSettings() },
	}
	idle := &cobra.Command{
		Use:   "idle <duration>",
		Short: "Stop a dev server after this long without a components.status (default 30m)",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			given := strings.TrimSpace(args[0])
			d, err := time.ParseDuration(given)
			if err != nil || d <= 0 {
				return fmt.Errorf("components.idle must be a positive duration, such as 30m or 2h, not %q", args[0])
			}
			cfg, path, err := a.loadConfig()
			if err != nil {
				return err
			}
			cfg.Components.Idle = given
			if err := config.Save(path, cfg); err != nil {
				return err
			}
			return a.showComponentsSettings()
		},
	}
	maxServers := &cobra.Command{
		Use:   "max-servers <n>",
		Short: "Run at most this many component dev servers at once per project, main's and every story's (default 3)",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			n, err := strconv.Atoi(strings.TrimSpace(args[0]))
			if err != nil || n <= 0 {
				return fmt.Errorf("components.max_servers must be a positive whole number, such as 3, not %q", args[0])
			}
			cfg, path, err := a.loadConfig()
			if err != nil {
				return err
			}
			cfg.Components.MaxServers = n
			if err := config.Save(path, cfg); err != nil {
				return err
			}
			return a.showComponentsSettings()
		},
	}
	c.AddCommand(idle, maxServers,
		&cobra.Command{Use: "show", Short: "Print components.idle, components.max_servers, and whether the action is enabled here", Args: cobra.NoArgs,
			RunE: func(*cobra.Command, []string) error { return a.showComponentsSettings() }},
	)
	return c
}

// showComponentsSettings prints components.idle, components.max_servers, and
// whether the components action is on for the project in the working
// directory.
func (a *app) showComponentsSettings() error {
	cfg, _, err := a.loadConfig()
	if err != nil {
		return err
	}
	here, enabled := "", false
	if repo, err := a.project(); err == nil {
		here = mainRootOf(repo)
		enabled = cfg.ActionEnabled(hostapi.ActionComponents, here)
	}
	idle, set := componentsIdleOf(cfg)
	most, mostSet := componentsMaxServersOf(cfg)
	if a.jsonOut {
		return a.printJSON(map[string]any{"idle": shortDuration(idle), "idle_seconds": int(idle / time.Second), "idle_set": set, "max_servers": most, "max_servers_set": mostSet, "enabled_here": enabled})
	}
	byDefault := func(set bool) string {
		if set {
			return ""
		}
		return " (default)"
	}
	fmt.Fprintf(a.out, "idle: %s%s: a dev server without a components.status for this long is stopped\n", shortDuration(idle), byDefault(set))
	fmt.Fprintf(a.out, "max_servers: %d%s: at most this many dev servers run at once per project, main's and every story's together\n", most, byDefault(mostSet))
	if here != "" {
		if enabled {
			fmt.Fprintln(a.out, "the components action is on for this project; flai serve disable components turns it off")
		} else {
			fmt.Fprintln(a.out, "the components action is off for this project; flai serve enable components turns it on")
		}
	}
	return nil
}

// shortDuration is d as Go writes it, without the zero units it ends with:
// 30m rather than 30m0s.
func shortDuration(d time.Duration) string {
	s := d.String()
	if strings.HasSuffix(s, "m0s") {
		s = strings.TrimSuffix(s, "0s")
	}
	if strings.HasSuffix(s, "h0m") {
		s = strings.TrimSuffix(s, "0m")
	}
	return s
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
// harness a story may name (S-0104), and what they lift on it (ADR-0131).
// Arguments after -- replace the adapter's defaults, even with none; --guard
// and --deny-protected change only themselves; --reset goes back to the
// defaults and lifts nothing. The command takes only --guard,
// --deny-protected, and --reset: its program is flai serve agent set's.
func newServeAgentHarnessCmd(a *app) *cobra.Command {
	var program, guard string
	var reset, denyProtected bool
	c := &cobra.Command{
		Use:   "harness [<name>] [--program <path>] [--guard none|\"\"] [--deny-protected=true|false] [--reset] [-- args...]",
		Short: "Set the program a harness is, the arguments that say what its agent may do, and what you lift on it",
		Long: `A story names its harness, such as claude-code; you say here what that is on
this host. --program is what is run, and the arguments after -- replace the
adapter's defaults, which say what the agent may do (for claude-code:
--permission-mode acceptEdits --allowedTools Bash,mcp__flai). "--" with
nothing after it runs it with none. --reset goes back to the defaults and
clears --guard and --deny-protected. With no name, every harness is shown.

On a harness that cannot run flai guard, a story whose agent has roles is
refused (ADR-0131). --guard none lifts that: the story's sub-agents run
unguarded, so a sub-agent could move items, write threads, or commit, which
only the story's agent should. --guard "" puts the refusal back; no other
value is taken.

On a harness that cannot ask before a protected write, a start is refused
unless the project's auto-approve host action is on. --deny-protected says
that your arguments for the harness already deny the protected paths in its
own vocabulary, which flai cannot read, so flai starts the agent on a harness
that cannot ask before a protected write; --deny-protected=false takes it
back.

The command harness takes only --guard, --deny-protected, and --reset: its
program and arguments are set with flai serve agent set.`,
		Example: `  flai serve agent harness
  flai serve agent harness claude-code --program /opt/claude/bin/claude
  flai serve agent harness claude-code -- --permission-mode acceptEdits --allowedTools Bash,Edit,mcp__flai
  flai serve agent harness claude-code --reset
  flai serve agent harness command --guard none
  flai serve agent harness command --guard ""
  flai serve agent harness command --deny-protected
  flai serve agent harness command --deny-protected=false`,
		RunE: func(cmd *cobra.Command, args []string) error {
			guardSet, denySet := cmd.Flags().Changed("guard"), cmd.Flags().Changed("deny-protected")
			dash := cmd.ArgsLenAtDash()
			names, rest := args, []string(nil)
			if dash >= 0 {
				names, rest = args[:dash], args[dash:]
			}
			known := append(settable(), harness.Command)
			if len(names) == 0 {
				if dash >= 0 || program != "" || reset || guardSet || denySet {
					return fmt.Errorf("name the harness: %s", strings.Join(known, ", "))
				}
				return a.showAgentCommand()
			}
			if len(names) > 1 {
				return fmt.Errorf("one harness at a time; arguments for it go after --")
			}
			name := names[0]
			if !slices.Contains(known, name) {
				return fmt.Errorf("there is no harness %q to set; there are: %s", name, strings.Join(known, ", "))
			}
			if name == harness.Command && (program != "" || dash >= 0) {
				return fmt.Errorf("the command's program and arguments are set with flai serve agent set -- <program> [args...]; flai serve agent harness command takes only --guard, --deny-protected, and --reset")
			}
			if reset && (program != "" || dash >= 0 || guardSet || denySet) {
				return fmt.Errorf("--reset takes neither --program, --guard, --deny-protected, nor arguments")
			}
			cfg, path, err := a.loadConfig()
			if err != nil {
				return err
			}
			h := cfg.Agent.Harnesses[name]
			if reset {
				h = config.HarnessHost{}
			}
			if program != "" {
				h.Program = program
			}
			if dash >= 0 {
				list := append([]string{}, rest...)
				h.Args = &list
			}
			if guardSet {
				h.Guard = guard
			}
			if denySet {
				h.DenyProtected = denyProtected
			}
			if cfg.Agent, err = cfg.Agent.WithHarness(name, h); err != nil {
				return err
			}
			if err := config.Save(path, cfg); err != nil {
				return err
			}
			return a.showAgentCommand()
		},
	}
	c.Flags().StringVar(&program, "program", "", "the program the harness is on this host")
	c.Flags().StringVar(&guard, "guard", "", `none lets a story's sub-agents run unguarded on a harness that cannot run flai guard; "" puts the refusal back`)
	c.Flags().BoolVar(&denyProtected, "deny-protected", false, "your arguments for the harness deny the protected paths, so flai starts the agent on a harness that cannot ask before a protected write")
	c.Flags().BoolVar(&reset, "reset", false, "go back to the adapter's program and arguments, and lift nothing")
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

// flai serve agent provider: where the model calls of a story's agent that
// names a provider go, on this host (ADR-0129). The flags given change those
// fields of the entry, --model given at all replaces its aliases, and --reset
// removes it. No flag takes a key, and nothing here reads the variable that
// holds it.
func newServeAgentProviderCmd(a *app) *cobra.Command {
	var p manifest.Provider
	var models []string
	var reset bool
	c := &cobra.Command{
		Use:   "provider [<name>] [--api <api>] [--base-url <url>] [--key-env <variable>] [--spend-key-env <variable>] [--model <alias>=<name>]... [--reset]",
		Short: "Set where the model calls of a story's agent that names a provider go: its API, its URL, and the variable that holds its key",
		Long: `A story's agent may name a provider, such as litellm or openrouter, beside
its harness and model; you say here what that is on this host. --api is what
the provider speaks (` + strings.Join(manifest.ProviderAPIs, ", ") + `), --base-url
where it answers, and --key-env the name of the environment variable that
holds its key. No flag takes a key: the key stays in that variable, which
flai serve reads when it starts the agent, and is never written to the
configuration. --spend-key-env, optional, names the variable of a key that
can read the gateway's spend log where the agent's own key cannot, such as
LiteLLM's master key; an empty value clears it. --model <alias>=<name>,
repeated, maps the harness's model aliases to names on the provider; given
at all, it replaces the entry's aliases.

A new entry needs --api, --base-url, and --key-env; on an entry that is
set, the flags given change those fields and leave the rest. --reset removes
the entry. With a name and no flag, the entry is shown; with no name, every
entry. A project's manifest may override an entry's api, base_url, and
models under providers.<name>, never its key_env or spend_key_env.`,
		Example: `  flai serve agent provider
  flai serve agent provider litellm --api anthropic-messages --base-url http://127.0.0.1:4000 --key-env LITELLM_API_KEY
  flai serve agent provider openrouter --api openai-chat --base-url https://openrouter.ai/api/v1 --key-env OPENROUTER_API_KEY --model opus=anthropic/claude-opus-5-5 --model sonnet=anthropic/claude-sonnet-5
  flai serve agent provider litellm --base-url https://litellm.example.com
  flai serve agent provider litellm --spend-key-env LITELLM_MASTER_KEY
  flai serve agent provider litellm --reset`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			changed := func(flags ...string) bool {
				return slices.ContainsFunc(flags, cmd.Flags().Changed)
			}
			given := changed("api", "base-url", "key-env", "spend-key-env", "model")
			if len(args) == 0 {
				if given || reset {
					return fmt.Errorf("name the provider, such as litellm: flai serve agent provider <name> --api <api> --base-url <url> --key-env <variable>")
				}
				return a.showProviders("")
			}
			name := args[0]
			if !given && !reset {
				return a.showProviders(name)
			}
			if reset && given {
				return fmt.Errorf("--reset takes no other flag: it removes provider %s whole", name)
			}
			cfg, path, err := a.loadConfig()
			if err != nil {
				return err
			}
			if reset {
				cfg.Agent = cfg.Agent.WithoutProvider(name)
			} else {
				entry := cfg.Agent.Providers[name]
				if changed("api") {
					entry.API = p.API
				}
				if changed("base-url") {
					entry.BaseURL = p.BaseURL
				}
				if changed("key-env") {
					entry.KeyEnv = p.KeyEnv
				}
				if changed("spend-key-env") {
					entry.SpendKeyEnv = p.SpendKeyEnv
				}
				if changed("model") {
					if entry.Models, err = providerModels(models); err != nil {
						return err
					}
				}
				if cfg.Agent, err = cfg.Agent.WithProvider(name, entry); err != nil {
					return err
				}
			}
			if err := config.Save(path, cfg); err != nil {
				return err
			}
			return a.showProviders(name)
		},
	}
	c.Flags().StringVar(&p.API, "api", "", "what the provider speaks: "+strings.Join(manifest.ProviderAPIs, ", "))
	c.Flags().StringVar(&p.BaseURL, "base-url", "", "the provider's absolute http or https URL, with no user or password in it")
	c.Flags().StringVar(&p.KeyEnv, "key-env", "", "the name of the environment variable that holds the provider's key, such as LITELLM_API_KEY; never the key")
	c.Flags().StringVar(&p.SpendKeyEnv, "spend-key-env", "", "the name of the environment variable that holds a key that can read the provider's spend log where the agent's key cannot, such as LITELLM_MASTER_KEY; never the key; empty clears it")
	c.Flags().StringArrayVar(&models, "model", nil, "a model alias and its name on the provider, alias=name (repeatable; replaces the entry's aliases)")
	c.Flags().BoolVar(&reset, "reset", false, "remove the provider's entry")
	return c
}

// providerModels are the model aliases --model gives, alias=name each; the
// names themselves are checked with the rest of the entry.
func providerModels(pairs []string) (map[string]string, error) {
	out := map[string]string{}
	for _, pair := range pairs {
		alias, model, ok := strings.Cut(pair, "=")
		alias = strings.TrimSpace(alias)
		if !ok || alias == "" {
			return nil, fmt.Errorf("--model %q is not alias=name; write one such as opus=anthropic/claude-opus-5-5", pair)
		}
		out[alias] = strings.TrimSpace(model)
	}
	return out, nil
}

// showProviders prints the host's provider name, or every one when name is
// empty: what each speaks, where, and the name of the variable that holds its
// key and, where set, of the one that holds its spend log's key, never a key.
func (a *app) showProviders(name string) error {
	cfg, _, err := a.loadConfig()
	if err != nil {
		return err
	}
	providers := cfg.Agent.Providers
	if name != "" {
		providers = map[string]manifest.Provider{}
		if p, ok := cfg.Agent.Providers[name]; ok {
			providers[name] = p
		}
	}
	if a.jsonOut {
		if providers == nil {
			providers = map[string]manifest.Provider{}
		}
		return a.printJSON(map[string]any{"providers": providers})
	}
	const add = "--api <api> --base-url <url> --key-env <variable>"
	switch {
	case len(providers) > 0:
	case name == "":
		fmt.Fprintf(a.out, "no provider is set on this host; add one with flai serve agent provider <name> %s\n", add)
		return nil
	default:
		fmt.Fprintf(a.out, "provider %s has no entry on this host; add it with flai serve agent provider %s %s\n", name, name, add)
		return nil
	}
	for _, n := range slices.Sorted(maps.Keys(providers)) {
		p := providers[n]
		fmt.Fprintf(a.out, "provider %s: %s at %s, its key in the variable %s\n", n, p.API, p.BaseURL, p.KeyEnv)
		if p.SpendKeyEnv != "" {
			fmt.Fprintf(a.out, "  spend log read with the key in the variable %s\n", p.SpendKeyEnv)
		}
		for _, alias := range slices.Sorted(maps.Keys(p.Models)) {
			fmt.Fprintf(a.out, "  model %s: %s\n", alias, p.Models[alias])
		}
	}
	return nil
}

// flai serve forge: the name of the environment variable a forge's token is
// read from, on this host (S-0392), kept beside the providers as their keys'
// variables are (ADR-0129). No flag takes a token, and the command reads the
// variable only to say whether it is set, never printing its value.
func newServeForgeCmd(a *app) *cobra.Command {
	var tokenEnv string
	var reset bool
	c := &cobra.Command{
		Use:   "forge [<forge>] [--token-env <NAME>] [--reset]",
		Short: "Set the environment variable a forge's token is read from, such as GITHUB_TOKEN for github",
		Long: `A project that accepts by pull request reaches its forge, such as github,
with a token; you say here which environment variable holds it on this host.
--token-env names the variable, ` + config.DefaultGitHubTokenEnv + ` for github when you name none, and
--reset goes back to that default. No flag takes a token: it stays in the
variable, which flai reads only when it calls the forge, and is never written
to the configuration, a log, a work item, or a start record.

With no forge, every forge flai calls is shown; with a forge and no flag,
that forge: the variable its token is read from, and whether that variable is
set in this environment, never its value. flai serve reads its own
environment, so export the variable where flai serve starts too.`,
		Example: `  flai serve forge
  flai serve forge github --token-env FLAI_GITHUB_TOKEN
  flai serve forge github --reset
  flai serve forge --json`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			given := cmd.Flags().Changed("token-env")
			if len(args) == 0 {
				if given || reset {
					return fmt.Errorf("name the forge, such as github: flai serve forge <forge> --token-env <NAME>")
				}
				return a.showForges("")
			}
			name := args[0]
			if !slices.Contains(config.KnownForges(), name) {
				return fmt.Errorf("flai calls no forge %q; it calls %s", name, strings.Join(config.KnownForges(), ", "))
			}
			if !given && !reset {
				return a.showForges(name)
			}
			if reset && given {
				return fmt.Errorf("--reset takes no --token-env: it goes back to the variable %s", config.DefaultForgeTokenEnv(name))
			}
			if given && strings.TrimSpace(tokenEnv) == "" {
				return fmt.Errorf("--token-env needs the name of a variable, such as %s; --reset goes back to it", config.DefaultForgeTokenEnv(name))
			}
			cfg, path, err := a.loadConfig()
			if err != nil {
				return err
			}
			if reset {
				cfg = cfg.WithoutForge(name)
			} else if cfg, err = cfg.WithForge(name, config.ForgeHost{TokenEnv: tokenEnv}); err != nil {
				return err
			}
			if err := config.Save(path, cfg); err != nil {
				return err
			}
			return a.showForges(name)
		},
	}
	c.Flags().StringVar(&tokenEnv, "token-env", "", "the name of the environment variable that holds the forge's token, such as "+config.DefaultGitHubTokenEnv+"; never the token")
	c.Flags().BoolVar(&reset, "reset", false, "read the forge's token from its default variable again")
	return c
}

// showForges prints forge, or every forge flai calls when forge is empty: the
// variable its token is read from, whether that is the default, and whether
// it is set and not empty in this environment, never its value.
func (a *app) showForges(forge string) error {
	cfg, _, err := a.loadConfig()
	if err != nil {
		return err
	}
	names := config.KnownForges()
	if forge != "" {
		names = []string{forge}
	}
	out := map[string]any{}
	for _, n := range names {
		env := cfg.ForgeTokenEnv(n)
		value, ok := os.LookupEnv(env)
		set := ok && strings.TrimSpace(value) != ""
		isDefault := env == config.DefaultForgeTokenEnv(n)
		if a.jsonOut {
			out[n] = map[string]any{"token_env": env, "default": isDefault, "token_set": set}
			continue
		}
		note := ""
		if isDefault {
			note = " (default)"
		}
		state := "set in this environment"
		if !set {
			state = "not set in this environment"
		}
		fmt.Fprintf(a.out, "forge %s: its token in the variable %s%s, %s\n", n, env, note, state)
	}
	if a.jsonOut {
		return a.printJSON(map[string]any{"forges": out})
	}
	return nil
}

// flai serve notify: the channel flai serve reminds the operator through of
// a thread that has awaited them past the manifest's threads.remind_after
// (S-0498, ADR-0205), a webhook, a shell command, or both, kept on the host
// beside the providers and forges.
func newServeNotifyCmd(a *app) *cobra.Command {
	c := &cobra.Command{
		Use:   "notify",
		Short: "The channel flai serve reminds you through of a thread that awaits you, on this host",
		Long: `flai serve reminds you of each thread that has awaited you longer than a
project's threads.remind_after, and again every threads.remind_every, through
the channel you set here, for every project it serves on this host. It writes
nothing on the thread. With no channel set it only logs the reminder.

--url is a webhook, http or https, that each reminder is POSTed to as JSON:
one attempt, a five second timeout, an answer other than 2xx a failure.
--command is run through the host's shell, sh -c, with the reminder's JSON on
its standard input and a 30 second timeout. Set either or both; set changes
only what you give it. The JSON holds the project, the thread, its title and
item, how long it has awaited you and since when, a pending recommendation's
author, first line, and source, and the dashboard's link to it, plus the
entry the dashboard posts to dashboard.notify_url, so one receiver can take
both. No token and no file content is sent.`,
		Example: `  flai serve notify set --url https://ntfy.sh/<topic>
  flai serve notify set --command 'curl -s -H "Title: flai" -d @- https://ntfy.sh/<topic>'
  flai serve notify show --json
  flai serve notify clear`,
		Args: cobra.NoArgs,
		RunE: func(*cobra.Command, []string) error { return a.showNotify() },
	}
	var target, command string
	set := &cobra.Command{
		Use:   "set [--url <url>] [--command <command>]",
		Short: "Set the webhook, the shell command, or both, that reminders are sent to",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			urlGiven, commandGiven := cmd.Flags().Changed("url"), cmd.Flags().Changed("command")
			if !urlGiven && !commandGiven {
				return fmt.Errorf("give --url <webhook>, --command <shell command>, or both; flai serve notify clear removes the channel")
			}
			if urlGiven && strings.TrimSpace(target) == "" {
				return fmt.Errorf("--url needs an http or https webhook, such as https://ntfy.sh/<topic>; flai serve notify clear removes the channel")
			}
			if commandGiven && strings.TrimSpace(command) == "" {
				return fmt.Errorf("--command needs a shell command that reads the reminder's JSON on its standard input; flai serve notify clear removes the channel")
			}
			cfg, path, err := a.loadConfig()
			if err != nil {
				return err
			}
			next := cfg.Notify
			if urlGiven {
				next.URL = strings.TrimSpace(target)
			}
			if commandGiven {
				next.Command = command
			}
			if cfg, err = cfg.WithNotify(next); err != nil {
				return err
			}
			if err := config.Save(path, cfg); err != nil {
				return err
			}
			return a.showNotify()
		},
	}
	set.Flags().StringVar(&target, "url", "", "an http or https webhook each reminder is POSTed to as JSON")
	set.Flags().StringVar(&command, "command", "", "a shell command run with each reminder's JSON on its standard input")
	clear := &cobra.Command{
		Use:   "clear",
		Short: "Remove the channel, so that flai serve only logs a reminder",
		Args:  cobra.NoArgs,
		RunE: func(*cobra.Command, []string) error {
			cfg, path, err := a.loadConfig()
			if err != nil {
				return err
			}
			if err := config.Save(path, cfg.WithoutNotify()); err != nil {
				return err
			}
			return a.showNotify()
		},
	}
	c.AddCommand(set, clear,
		&cobra.Command{Use: "show", Short: "Print the webhook and the shell command reminders are sent to", Args: cobra.NoArgs,
			RunE: func(*cobra.Command, []string) error { return a.showNotify() }},
	)
	return c
}

// showNotify prints the host's reminder channel: the webhook and the shell
// command, each or none. The operator set both in their own shell, so the
// URL is printed whole, secret and all, as they wrote it.
func (a *app) showNotify() error {
	cfg, _, err := a.loadConfig()
	if err != nil {
		return err
	}
	n := cfg.Notify
	if a.jsonOut {
		return a.printJSON(map[string]any{"url": n.URL, "command": n.Command})
	}
	if n == (config.Notify{}) {
		fmt.Fprintln(a.out, "no reminder channel is set; flai serve only logs a reminder (flai serve notify set --url <url> or --command <command>)")
		return nil
	}
	if n.URL != "" {
		fmt.Fprintf(a.out, "url: %s\n", n.URL)
	}
	if n.Command != "" {
		fmt.Fprintf(a.out, "command: %s\n", n.Command)
	}
	return nil
}

// hostProviders are the host's providers for the dashboard's settings page
// (S-0359): each entry's fields and whether the variables its key_env and
// spend_key_env name are set and not empty in lookup's environment, never a
// variable's value.
func hostProviders(set map[string]manifest.Provider, lookup func(string) (string, bool)) map[string]any {
	isSet := func(name string) bool {
		if name == "" {
			return false
		}
		value, ok := lookup(name)
		return ok && value != ""
	}
	out := map[string]any{}
	for name, p := range set {
		models := p.Models
		if models == nil {
			models = map[string]string{}
		}
		out[name] = map[string]any{"api": p.API, "base_url": p.BaseURL, "key_env": p.KeyEnv, "models": models, "key_set": isSet(p.KeyEnv),
			"spend_key_env": p.SpendKeyEnv, "spend_key_set": isSet(p.SpendKeyEnv)}
	}
	return out
}

// hostHarnesses are the host's harnesses for the settings page: each one's
// program and arguments, with set saying whether the host gives either rather
// than the adapter's defaults; its guard and deny_protected (ADR-0131); and
// whether its adapter runs flai guard and can ask before a protected write,
// where neither lifts anything. The command's program and arguments are flai
// serve agent set's command, never set here.
func (a *app) hostHarnesses(cfg config.Config) map[string]any {
	out := map[string]any{}
	for _, n := range append(settable(), harness.Command) {
		var h harness.Host
		if n != harness.Command {
			h = a.harnessHost(cfg, n)
		} else if len(cfg.Agent.Command) > 0 {
			h = harness.Host{Program: cfg.Agent.Command[0], Args: cfg.Agent.Command[1:]}
		}
		args := h.Args
		if args == nil {
			args = []string{}
		}
		set := cfg.Agent.Harnesses[n]
		caps := harness.Adapters[n].Capabilities()
		out[n] = map[string]any{"program": h.Program, "args": args, "set": set.Program != "" || set.Args != nil,
			"guard": set.Guard, "deny_protected": set.DenyProtected, "guard_capable": caps.Guard, "ask_capable": caps.Ask}
	}
	return out
}

// hostSettings is the host's configuration as it applies to the project at
// root, for the dashboard's settings page (S-0105) and its Host page
// (S-0489). The tokens are not in it: a rotation answers with the new one,
// once.
func (a *app) hostSettings(root string) any {
	cfg, _, err := a.loadConfig()
	if err != nil {
		return map[string]any{"error": err.Error()}
	}
	all := hostActions(root) // those the roles the project defines name among them (S-0417)
	names := make([]string, 0, len(all))
	for n := range all {
		if dashboardSees(n) {
			names = append(names, n)
		}
	}
	sort.Strings(names)
	actions := make([]map[string]any, 0, len(names))
	for _, n := range names {
		actions = append(actions, map[string]any{"name": n, "means": all[n],
			"here": cfg.ActionEnabled(n, root), "everywhere": cfg.ActionEnabled(n, config.AllProjects)})
	}
	providers := hostProviders(cfg.Agent.Providers, os.LookupEnv)
	command := cfg.Agent.Command
	if command == nil {
		command = []string{}
	}
	harnesses := a.hostHarnesses(cfg)
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
	// allowed_hosts is the hosts the dashboard is to accept beside its own
	// address, as the configuration names them now, which the container's next
	// start gives it (S-0489); flai dashboard status says what it runs with.
	out := map[string]any{
		"actions": actions,
		"agent": map[string]any{"name": cfg.Agent.Name, "command": command,
			"harnesses": harnesses, "providers": providers},
		"allowed_hosts": config.NormalizeAllowedHosts(cfg.Dashboard.AllowedHosts),
		"checks":        map[string]any{"commands": checks, "timeout_minutes": timeout},
		"import_roots":  roots,
		"projects":      a.servedView(),
	}
	if repo, err := workitem.Open(root); err == nil {
		out["default_agent"] = repo.Manifest.Agent
		// the shared paths as the manifest has them now (S-0295, ADR-0096)
		out["shared"] = nonNil(repo.Manifest.Claims.Shared)
		out["manifest_checks"] = repo.Manifest.Checks
		out["planning"] = planningTriggers(repo.Manifest.Planning, cfg.ActionEnabled(hostapi.ActionPlan, root), a.now())
		out["strategic"] = strategicSettings(repo.Manifest, cfg.ActionEnabled(hostapi.ActionSettings, root))
		if st, running := readMCPState(repo, a.now()); running {
			out["mcp"] = map[string]any{"running": true, "url": st.URL, "pid": st.PID}
		} else {
			out["mcp"] = map[string]any{"running": false}
		}
	}
	return out
}

// planningTriggers is when flai serve plans again on its own, for the
// settings page (S-0211, ADR-0084): plan, whether the plan host action is on
// for the project, which every trigger needs and under which an edit to a
// planned story queues the planner; replan, planning.replan or its default,
// with replan_set saying whether the manifest sets it; and schedule as
// written, with next, its next run in UTC. A value flai cannot read is
// given with its error instead.
func planningTriggers(p manifest.Planning, plan bool, now time.Time) map[string]any {
	out := map[string]any{"plan": plan, "replan_set": strings.TrimSpace(p.Replan) != "",
		"schedule": strings.TrimSpace(p.Schedule)}
	if policy, err := p.ReplanPolicy(); err != nil {
		out["replan"] = p.Replan
		out["replan_error"] = err.Error()
	} else {
		out["replan"] = policy
	}
	if sched, err := p.PlanSchedule(); err != nil {
		out["schedule_error"] = err.Error()
	} else if sched != nil {
		out["next"] = sched.Next(now).UTC().Format(time.RFC3339)
	}
	return out
}

// strategicSettings are the strategic agents' settings for their panels
// (S-0229): each key flai manifest set writes, with its value as the
// manifest has it and whether it is set, its default, kind, values, meaning,
// and for a permission its risk; editable, whether the settings host action
// is on for the project, which settings.manifest needs; and enable, the
// command that turns it on, run in a shell on the host.
func strategicSettings(m manifest.Manifest, editable bool) map[string]any {
	return map[string]any{"editable": editable, "enable": hostapi.EnableCommand(hostapi.ActionSettings),
		"settings": m.SettingValues()}
}

// servedView is what the settings page shows of the projects flai serve
// serves, and of those below the folders named for import that it does not,
// with why (S-0122): what flai serve project list says, without the
// dashboard's address or the credential's path.
func (a *app) servedView() map[string]any {
	l, err := a.listServedProjects()
	if err != nil {
		return map[string]any{"error": err.Error()}
	}
	cfg, _, err := a.loadConfig()
	if err != nil {
		return map[string]any{"error": err.Error()}
	}
	type project struct {
		Key  string `json:"key"`
		Name string `json:"name"`
		Root string `json:"root"`
		From string `json:"from,omitempty"`
		// Below is the folder it is below that flai serve serves from, when
		// it is: removing it puts it on the list of removed projects.
		Below     string `json:"below,omitempty"`
		State     string `json:"state,omitempty"`
		Since     string `json:"since,omitempty"`
		LastError string `json:"last_error,omitempty"`
		Reason    string `json:"reason,omitempty"`
		// Removed is whether it is not served because the operator removed
		// it (S-0123).
		Removed bool `json:"removed,omitempty"`
		// Settings is whether the dashboard may serve or remove it: the
		// settings action is on for it.
		Settings bool `json:"settings"`
	}
	served := make([]project, 0, len(l.Served))
	for _, p := range l.Served {
		// a registered one below such a folder too: removing it lists it (S-0123)
		below := a.servedBelow(p.Root)
		served = append(served, project{Key: p.Key, Name: p.Name, Root: p.Root, From: p.From, Below: below,
			State: p.State, Since: p.Since, LastError: p.LastError, Reason: p.Reason, Settings: cfg.ActionEnabled(hostapi.ActionSettings, p.Root)})
	}
	unserved := make([]project, 0, len(l.Unserved))
	for _, f := range l.Unserved {
		unserved = append(unserved, project{Key: f.Key, Name: f.Name, Root: f.Root, Reason: f.Reason, Removed: f.Removed, Settings: cfg.ActionEnabled(hostapi.ActionSettings, f.Root)})
	}
	return map[string]any{"running": l.Running, "served": served, "unserved": unserved}
}

func newServeAgentRestartCmd(a *app) *cobra.Command {
	return &cobra.Command{
		Use:   "restart <story-id>",
		Short: "Start a new agent for a story in ready or in progress whose agent dropped or failed, or that has none here",
		Long: `Starts a new agent, in a new session, for a story in ready or
in-progress whose last agent flai serve started has ended or dropped, the
way flai serve starts one when a story enters ready: the story's harness,
model, and options, or the host's command. The agent is told how its last
one ended and to go on from the story's narrative. The run is recorded
where the serving flai tracks it: the dot on the card, the outcome, the
restart on an answer (S-0116, ADR-0043).

A story that this host's flai serve has started no agent for gets one too
(S-0177, ADR-0064). One in progress, begun on another host, gets an agent
told who began it, when, and on which host, which threads have been written
to since, and to reconcile with what is committed rather than start over.

For a story in ready while the in-progress limit is full, review is full,
or a claim holds it, it queues the new agent instead: flai serve starts it as soon as the
story can start, as it starts a story that enters ready, and until then the
story's agent reads as waiting, queued (S-0118).

It refuses, and says why, while the agent action is off for the project,
when the story is in another state, while its agent runs or waits for an
answer, when one is already queued, and when nothing can start it. The story
page's Retry button runs this, and its Start agent button for a story in
progress that has no agent here.`,
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
in-progress limit or a full review does not stop it: it starts past them,
with a warning, as a move does. A story in progress whose agent dropped or failed,
or that has had no agent here, gets one from flai serve agent restart. The story page's Start agent
button runs this.`,
		Example: `  flai serve agent start S-0115`,
		Args:    cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			return a.agentNow(args[0], serve.Start)
		},
	}
}

func newServeAgentCommitCmd(a *app) *cobra.Command {
	return &cobra.Command{
		Use:   "commit <story-id>",
		Short: "Start an agent to commit what a story in review left uncommitted in its worktree",
		Long: `Starts the story's agent, in a new session, to commit what its worktree
under .flai-cache/worktrees holds and nothing else: a story in review whose
worktree has uncommitted changes cannot be accepted, because acceptance
merges the branch as it is committed and removes the worktree. The agent is
told to commit the changes on the story branch the way the project's
conventions say, to discard only what plainly does not belong in the
repository, to log what it did in the narrative, and to leave the story in
review. The run is recorded where the serving flai tracks it, like a start
(S-0140).

It refuses, and says why, while the agent action is off for the project,
when the story is not in review, while its agent runs or waits for an
answer, when the story has no worktree or the worktree has nothing
uncommitted, and when nothing can start it. The Have an agent commit them
button in the dashboard's acceptance confirmation runs this.`,
		Example: `  flai serve agent commit S-0140`,
		Args:    cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			return a.agentNow(args[0], serve.Commit)
		},
	}
}

func newServeAgentStopCmd(a *app) *cobra.Command {
	return &cobra.Command{
		Use:   "stop <story-id>",
		Short: "Stop a story's agent: end its process and everything it started",
		Long: `Stops the agent flai serve started for a story (S-0170). An agent that
runs is sent SIGTERM with every process it started (its process group), and
SIGKILL if it has not ended ten seconds later. An agent that ended waiting
for an answer is not started again when the answer comes. The run is
recorded as stopped by the operator: the story stays where it is, its
worktree holds what the agent left, committed or not, and it gets no new
agent until it is retried (flai serve agent restart) or moved back to ready.

A process is signalled only while it is still the agent flai started: the
leader of its own session, which on Linux started at the clock tick since
boot recorded with the run. A PID the system has given to another process
since, after a reboot, is left alone, and the run is only recorded as
stopped; flai serve settles such a run by itself at its next look.

It refuses, and says why, when flai serve has started no agent for the story
and when its agent is neither running nor waiting for an answer. It works
whether or not the agent action is on: the dashboard's Stop, on the
activity page, asks for it only while it is.`,
		Example: `  flai serve agent stop S-0170`,
		Args:    cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			return a.agentStop(args[0])
		},
	}
}

// flai serve orchestrate: the orchestrator stopped and held stopped, and
// started again, on the operator's word (S-0228).
func newServeOrchestrateCmd(a *app) *cobra.Command {
	c := &cobra.Command{
		Use:   "orchestrate",
		Short: "Stop the project's orchestrator and hold it stopped, or start it again",
		Long: `While the orchestrate host action is on, flai serve runs the project's
orchestrator, and starts it again whenever its run ends (S-0218). stop ends
the run and holds the orchestrator stopped, so that flai serve does not start
it again while the action stays on; start lifts the hold and starts it, as
turning the action on does. Turning the action off and on lifts the hold too.
The dashboard's orchestrator page runs these; both are refused while the
action is off.`,
		Example: `  flai serve orchestrate stop
  flai serve orchestrate start`,
		Args: cobra.NoArgs,
	}
	c.AddCommand(
		&cobra.Command{
			Use:   "stop",
			Short: "End the orchestrator's run and hold it stopped while orchestrate stays on",
			Long: `Ends the orchestrator's run as flai serve agent stop ends a story's agent:
its process group is sent SIGTERM, and SIGKILL if it has not ended ten
seconds later, while it is still the orchestrator flai started. The run is
recorded as stopped by the operator, its activity is logged in
wip/agents/orchestrator.md as stopped from the dashboard, and the
orchestrator is held stopped: flai serve does not start it again until flai
serve orchestrate start, or until orchestrate is turned off and on. A run
that failed and waits to be started again is held as it is.

It refuses, and says why, while the orchestrate action is off for the
project, when flai serve has started no orchestrator for it, and when the
orchestrator is already held stopped.`,
			Args: cobra.NoArgs,
			RunE: func(*cobra.Command, []string) error { return a.orchestrateNow(false) },
		},
		&cobra.Command{
			Use:   "start",
			Short: "Start the orchestrator now, lifting a hold on it",
			Long: `Starts an orchestrator run at once, for the reason asked, with any reasons
flai serve has queued for it (ADR-0209): in the project's main checkout, with
its orchestration agent and the harnesses and the command set with flai serve
agent. It lifts the hold flai serve orchestrate stop put on the orchestrator,
if there is one. The run acts on its reasons and ends, and flai serve starts
the next when another reason comes. A start that fails lifts the hold as
well, keeps the reasons, and flai serve tries again a minute later.

It refuses, and says why, while the orchestrate action is off for the
project, and while an orchestrator run is going or still being stopped.`,
			Args: cobra.NoArgs,
			RunE: func(*cobra.Command, []string) error { return a.orchestrateNow(true) },
		},
	)
	return c
}

// orchestrateNow stops and holds, or starts, the orchestrator of the project
// in the working directory, and prints the run as recorded.
func (a *app) orchestrateNow(start bool) error {
	repo, err := a.project()
	if err != nil {
		return err
	}
	o := serve.Options{Dir: a.serveDir(), Logger: a.logger(), Now: a.now, Agent: a.agentConfig, Host: a.host()}
	e := serve.Entry{Key: repo.Manifest.Key, Name: repo.Manifest.Name, Root: mainRootOf(repo)}
	var run *serve.AgentRun
	if start {
		run, err = serve.OrchestratorStart(context.Background(), o, e)
	} else {
		run, err = serve.OrchestratorStop(o, e)
	}
	var no *serve.Refused
	if errors.As(err, &no) {
		return fmt.Errorf("rule: %s", no.Why)
	}
	if err != nil {
		return err
	}
	if a.jsonOut {
		return a.printJSON(map[string]any{"held": run.Held, "orchestrator": run})
	}
	if start {
		fmt.Fprintf(a.out, "started %s (%s) as %s (pid %d); log %s\n", run.Command, run.Harness, run.Agent, run.PID, run.Log)
		return nil
	}
	fmt.Fprintf(a.out, "stopped the orchestrator %s (pid %d), and held it stopped; flai serve orchestrate start starts it again\n", run.Agent, run.PID)
	return nil
}

// agentStop stops story's agent and prints the run as it ended.
func (a *app) agentStop(story string) error {
	repo, err := a.project()
	if err != nil {
		return err
	}
	o := serve.Options{Dir: a.serveDir(), Logger: a.logger(), Now: a.now, Agent: a.agentConfig, Host: a.host()}
	e := serve.Entry{Key: repo.Manifest.Key, Name: repo.Manifest.Name, Root: mainRootOf(repo)}
	run, err := serve.Stop(o, e, workitem.CanonicalID(story))
	var no *serve.Refused
	if errors.As(err, &no) {
		return fmt.Errorf("rule: %s", no.Why)
	}
	if err != nil {
		return err
	}
	if a.jsonOut {
		return a.printJSON(map[string]any{"story": run.Story, "agent": run.Agent, "pid": run.PID, "stopped": run.Stopped, "ended": run.Ended, "why": run.Why})
	}
	fmt.Fprintf(a.out, "stopped %s's agent %s (pid %d)\n", run.Story, run.Agent, run.PID)
	return nil
}

func newServeAgentStreamCmd(a *app) *cobra.Command {
	var from int64
	var follow bool
	c := &cobra.Command{
		Use:   "stream <story-id>",
		Short: "Print what a story's agent said and did, from the log flai serve gave it",
		Long: `Prints the stream of the newest agent flai serve started for a story:
what it said, the tools it called and what they answered, the background
tasks it ran, and how its session ended, one entry to a line, each cut to a
few hundred characters. It is read from the log flai serve writes the
agent's output to; Claude Code's stream-json is turned into entries, and any
other output is printed as it is. The dashboard's activity page shows the
same stream (S-0142).

Without --from it starts at the last 256 KiB of the log. --follow goes on
printing what the agent writes until it ends. --json prints each read as
one JSON object on a line: the entries, whether the agent runs, and next,
the offset to read from again with --from.`,
		Example: `  flai serve agent stream S-0142
  flai serve agent stream S-0142 --follow
  flai serve agent stream S-0142 --from 0 --json`,
		Args: cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			return a.agentStream(workitem.CanonicalID(args[0]), from, follow)
		},
	}
	c.Flags().Int64Var(&from, "from", -1, "the byte offset of the log to read from (default: its last 256 KiB)")
	c.Flags().BoolVarP(&follow, "follow", "f", false, "go on printing until the agent ends")
	return c
}

// agentStream prints a story's agent's stream, once or until it ends.
func (a *app) agentStream(story string, from int64, follow bool) error {
	repo, err := a.project()
	if err != nil {
		return err
	}
	root := mainRootOf(repo)
	started := ""
	for {
		got, err := serve.Stream(a.serveDir().AgentStates()[root], story, from)
		if errors.Is(err, hostapi.ErrNoAgent) {
			return fmt.Errorf("rule: flai serve has started no agent for %s", story)
		}
		if err != nil {
			return err
		}
		if started != "" && got.Started != started {
			// another agent was started for the story: its log is new
			from = 0
			started = got.Started
			continue
		}
		if a.jsonOut {
			line, err := json.Marshal(got)
			if err != nil {
				return err
			}
			fmt.Fprintln(a.out, string(line))
		} else {
			if started == "" {
				state := "ended " + got.Ended
				if got.Running {
					state = "running"
				}
				fmt.Fprintf(a.out, "%s: %s, started %s, %s\n", story, got.Agent, got.Started, state)
			}
			printStream(a, got)
		}
		started, from = got.Started, got.Next
		if !follow || (!got.Running && !got.More) {
			return nil
		}
		if !got.More {
			a.sleep(time.Second)
		}
	}
}

// printStream prints the entries of one read, one to a line, the lines of
// an entry's text after its first indented under it.
func printStream(a *app, got *serve.StreamRead) {
	if got.Skipped > 0 {
		fmt.Fprintf(a.out, "(%d earlier %s left out)\n", got.Skipped, oneOrMany(got.Skipped, "entry", "entries"))
	}
	for _, e := range got.Entries {
		text := e.Text
		if e.Tool != "" {
			text = e.Tool + ": " + text
		}
		kind := e.Kind
		if e.Error {
			kind += "!"
		}
		fmt.Fprintf(a.out, "%-9s %s\n", kind, strings.ReplaceAll(text, "\n", "\n          "))
	}
}

// mcpAgents starts or restarts a story's agent for flai mcp's tools
// agent_start and agent_restart (S-0177, ADR-0064), as flai serve agent
// start and restart do, under the same agent host action, and journals who
// asked and what came of it. They are for the operator's own agent: one that
// flai serve started, whose flai mcp has FLAI_STARTED_BY, is refused, since
// a start goes past holds, the in-progress limit, and a full review on the
// operator's word.
func (a *app) mcpAgents(ctx context.Context, verb, root, story, by string) (mcpserver.AgentStarted, error) {
	if os.Getenv("FLAI_STARTED_BY") == "flai-serve" {
		return mcpserver.AgentStarted{}, fmt.Errorf("an agent flai serve started does not start agents: %s's agent is the operator's to %s, from the story's page or with flai serve agent %s %s on the host", story, verb, verb, story)
	}
	how := serve.Start
	if verb == "restart" {
		how = serve.Restart
	}
	repo, err := workitem.Open(root)
	if err != nil {
		return mcpserver.AgentStarted{}, err
	}
	host := a.host()
	o := serve.Options{Dir: a.serveDir(), Logger: a.logger(), Now: a.now, Agent: a.agentConfig, Host: host}
	e := serve.Entry{Key: repo.Manifest.Key, Name: repo.Manifest.Name, Root: root}
	run, err := how(ctx, o, e, story)
	entry := hostapi.Entry{At: a.now().UTC().Format(time.RFC3339), Action: hostapi.ActionAgent, Method: "mcp.agent_" + verb, Project: e.Key, Root: root, By: by, Outcome: "done"}
	var no *serve.Refused
	switch {
	case errors.As(err, &no) && !o.Agent(root).Enabled:
		entry.Outcome, entry.Detail = "disabled", fmt.Sprintf("%s asked to %s %s's agent: %s", by, verb, story, no.Why)
	case err != nil:
		entry.Outcome, entry.Detail = "failed", fmt.Sprintf("%s asked to %s %s's agent: %s", by, verb, story, err)
	case run.Queued != "":
		entry.Detail = fmt.Sprintf("%s asked to %s %s's agent: queued, and flai serve starts it when it can", by, verb, story)
	default:
		entry.Detail = fmt.Sprintf("%s asked to %s %s's agent: started %s as %s (pid %d)", by, verb, story, run.Command, run.Agent, run.PID)
	}
	host.Record(entry)
	if err != nil {
		return mcpserver.AgentStarted{}, err
	}
	return mcpserver.AgentStarted{Story: run.Story, Agent: run.Agent, Harness: run.Harness, Command: run.Command, PID: run.PID, Log: run.Log, Started: run.Started, Queued: run.Queued}, nil
}

// agentNow starts story's agent on the operator's word, with start or
// restart deciding whether it may, and prints the run, or that it was
// queued.
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
	if run.Queued != "" {
		if a.jsonOut {
			return a.printJSON(map[string]any{"story": run.Story, "agent": run.Agent, "queued": run.Queued})
		}
		fmt.Fprintf(a.out, "queued another agent for %s: the in-progress limit, a full review, or a hold leaves no room, and flai serve starts it when there is\n", run.Story)
		return nil
	}
	if a.jsonOut {
		return a.printJSON(map[string]any{"story": run.Story, "agent": run.Agent, "harness": run.Harness, "command": run.Command, "pid": run.PID, "log": run.Log, "started": run.Started})
	}
	fmt.Fprintf(a.out, "started %s (%s) for %s as %s (pid %d); log %s\n", run.Command, run.Harness, run.Story, run.Agent, run.PID, run.Log)
	return nil
}
