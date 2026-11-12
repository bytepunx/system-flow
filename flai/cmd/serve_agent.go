package cmd

import (
	"errors"
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"slices"
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
(claude-code runs claude -p with the story's model, and flai's MCP server;
opencode runs opencode run with its provider/model and flai's MCP server),
with the program and the arguments you set for it here with flai serve agent
harness, which say what the agent may do. A story names only tunables the
adapter checks, never a program, a flag, or a permission.

A story that names no harness, or names "command", is started with your
command, when one is set: an argument list, run as it stands in the
project's directory, as you, never through a shell. In an argument {story},
{root}, {model}, and {harness} are replaced; nothing else is interpreted, and
the story's options are in FLAI_AGENT_CONFIG as JSON.

A story whose agent and project name no harness, on a project with no
command, is started with flai's built-in default agent, claude-code with the
harness's own default model (ADR-0261).

The command, the name, and each harness's program and arguments are kept per
project (ADR-0261): set, clear, and harness write those of the project in the
working directory, or of --project <root>, and those of every project, the
host-wide fallback a project that sets none of its own reads, only with
--all-projects. A harness a project sets replaces the host-wide entry of that
harness whole. The providers and the restart limit are host-wide.

Every session's environment carries FLAI_AGENT (your --name and the story,
such as agent-S-0104), FLAI_STORY, and FLAI_SESSION. Its output goes to a log
beside flai serve's state, and every start and failure is in flai serve
journal.`,
		Example: `  flai serve agent harness claude-code
  flai serve agent harness claude-code --program /opt/claude/bin/claude -- --permission-mode acceptEdits --allowedTools Bash,Monitor,mcp__flai
  flai serve agent set -- /home/me/bin/start-agent {story} --model {model}
  flai serve agent set --all-projects --name builder
  flai serve agent harness claude-code --project /home/me/src/other --program /opt/claude/bin/claude
  flai serve agent show
  flai serve enable agent
  flai serve agent clear`,
		Args: cobra.NoArgs,
		RunE: func(*cobra.Command, []string) error { return a.showAgentCommand("") },
	}
	var name, setProject, clearProject string
	var attended, restarts int
	var setAll, clearAll bool
	set := &cobra.Command{
		Use:   "set [--name] [--auto-restarts <n>] [--project <root> | --all-projects] [-- <program> [args...]]",
		Short: "Set the command, as an argument list after --, the name, or how many times a story's agent is restarted on its own",
		Long: `Sets what you give and leaves the rest as it was: the command, as an
argument list after --; the name the sessions work under; and, with
--auto-restarts, how many times flai serve restarts a story's agent on its own
after it ends with its story in progress, 2 when unset and 0 for never
(ADR-0108). Past that, flai serve opens a thread on the story for you.

The command and the name are the project's in the working directory, or of
--project <root>; with --all-projects they are every project's, the host-wide
fallback a project that sets none of its own reads (ADR-0261). The restart
limit is host-wide whichever you give.

A change of the command's program or arguments lifts the hold of the old
one's budget state, as flai serve agent harness does for a harness's.`,
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
			root := ""
			if len(args) > 0 || name != "" {
				var err error
				if root, err = a.agentScope(setProject, setAll); err != nil {
					return err
				}
			}
			cfg, path, err := a.loadConfig()
			if err != nil {
				return err
			}
			was := a.harnessIdentity(cfg.Agent.ForProject(root), harness.Command)
			if restartsSet {
				if cfg.Agent, err = cfg.Agent.WithAutoRestarts(restarts); err != nil {
					return err
				}
			}
			if root == "" {
				if len(args) > 0 {
					cfg.Agent.Command = args
				}
				if name != "" {
					cfg.Agent.Name = name
				}
			} else if len(args) > 0 || name != "" {
				p := cfg.Agent.Project(root)
				if len(args) > 0 {
					p.Command = args
				}
				if name != "" {
					p.Name = name
				}
				if cfg.Agent, err = cfg.Agent.WithProject(root, p); err != nil {
					return err
				}
			}
			if err := config.Save(path, cfg); err != nil {
				return err
			}
			var lifted []budgetLift
			if len(args) > 0 {
				if lifted, err = a.liftBudgets(cmd.Context(), cfg.Notify, "flai serve agent set", harness.Command, root, was, a.harnessIdentity(cfg.Agent.ForProject(root), harness.Command)); err != nil {
					return err
				}
			}
			return a.showAgentCommand(root, lifted...)
		},
	}
	set.Flags().StringVar(&name, "name", "", "the FLAI_AGENT the session works under (default agent)")
	set.Flags().IntVar(&restarts, "auto-restarts", config.DefaultAutoRestarts, "how many times flai serve restarts a story's agent on its own after it ends with its story in progress; 0 for never; host-wide")
	set.Flags().StringVar(&setProject, "project", "", "the main checkout of the project whose command and name these are, not the one in the working directory")
	set.Flags().BoolVar(&setAll, "all-projects", false, "every project's command and name, the host-wide fallback a project that sets none of its own reads")
	// Retired (S-0116, ADR-0043): accepted so that a script that sets it still runs.
	set.Flags().IntVar(&attended, "attended-minutes", 0, "retired: nobody attending holds a ready story back any more")
	_ = set.Flags().MarkHidden("attended-minutes")
	clear := &cobra.Command{Use: "clear [--project <root> | --all-projects]", Short: "Remove the command and the name; the harnesses, the providers, and the restart limit stay",
		Long: `Removes the command and the name of the project in the working directory,
or of --project <root>, so that it reads every project's again; with
--all-projects, every project's own, the host-wide fallback (ADR-0261). The
harnesses, the providers, and the restart limit stay.`,
		Args: cobra.NoArgs,
		RunE: func(*cobra.Command, []string) error {
			root, err := a.agentScope(clearProject, clearAll)
			if err != nil {
				return err
			}
			cfg, path, err := a.loadConfig()
			if err != nil {
				return err
			}
			if root == "" {
				cfg.Agent.Command, cfg.Agent.Name, cfg.Agent.AttendedMinutes = nil, "", 0
			} else {
				p := cfg.Agent.Project(root)
				p.Command, p.Name = nil, ""
				if cfg.Agent, err = cfg.Agent.WithProject(root, p); err != nil {
					return err
				}
			}
			if err := config.Save(path, cfg); err != nil {
				return err
			}
			return a.showAgentCommand(root)
		}}
	clear.Flags().StringVar(&clearProject, "project", "", "the main checkout of the project whose command and name to remove, not the one in the working directory")
	clear.Flags().BoolVar(&clearAll, "all-projects", false, "every project's command and name, the host-wide fallback")
	c.AddCommand(set, newServeAgentHarnessCmd(a), newServeAgentProviderCmd(a), newServeAgentStartCmd(a), newServeAgentRestartCmd(a), newServeAgentCommitCmd(a), newServeAgentStopCmd(a), newServeAgentStreamCmd(a), newServeAgentUsageCmd(a),
		&cobra.Command{Use: "show", Short: "Print the command, the harnesses, and whether the action is enabled, for the project here", Args: cobra.NoArgs,
			RunE: func(*cobra.Command, []string) error { return a.showAgentCommand("") }},
		clear,
	)
	return c
}

// agentScope is whose command, name, and harnesses flai serve agent set,
// clear, and harness write (ADR-0261): the project of --project, or of the
// working directory, by its main checkout; or, with --all-projects, every
// project's, the host-wide fallback, as root "". With neither and no project
// in the working directory, it refuses, naming both.
func (a *app) agentScope(project string, all bool) (string, error) {
	switch {
	case all && project != "":
		return "", errors.New("give --project <root> or --all-projects, not both")
	case all:
		return "", nil
	case project != "":
		abs, err := filepath.Abs(config.ExpandHome(project))
		if err != nil {
			return "", fmt.Errorf("--project %q: resolve it to an absolute path: %w", project, err)
		}
		repo, err := workitem.Open(abs)
		if err != nil {
			return "", fmt.Errorf("--project %q: %w; name the main checkout of a system-flow project", project, err)
		}
		return mainRootOf(repo), nil
	}
	repo, err := a.project()
	if err != nil {
		return "", fmt.Errorf("%w; or name a project with --project <root>, or every project with --all-projects", err)
	}
	return mainRootOf(repo), nil
}

// agentFrom says where a project's agent setting comes from (ADR-0261):
// fromProject when the project sets it, fromAllProjects when the host-wide
// fallback does, and "" when neither does.
func agentFrom(project, everywhere bool) string {
	switch {
	case project:
		return fromProject
	case everywhere:
		return fromAllProjects
	}
	return ""
}

// Where a project's agent setting comes from, as show and the settings page
// say it (ADR-0261).
const (
	fromProject     = "project"
	fromAllProjects = "all-projects"
)

// fromSaid is agentFrom's answer as show prints it after a value.
func fromSaid(from string) string {
	switch from {
	case fromProject:
		return " (this project)"
	case fromAllProjects:
		return " (every project)"
	}
	return ""
}

// showAgentCommand prints the agent settings that apply to the project at
// root, or, with root "", to the one in the working directory, each saying
// whether it is the project's own or every project's (ADR-0261); outside a
// project, the host-wide ones; and the budget states the change just made
// lifted, when it lifted any (S-0720).
func (a *app) showAgentCommand(root string, lifted ...budgetLift) error {
	cfg, _, err := a.loadConfig()
	if err != nil {
		return err
	}
	here, enabled := root, false
	if here == "" {
		if repo, err := a.project(); err == nil {
			here = mainRootOf(repo)
		}
	}
	if here != "" {
		enabled = cfg.ActionEnabled(hostapi.ActionAgent, here)
	}
	agent, own := cfg.Agent.ForProject(here), cfg.Agent.Project(here)
	commandFrom := agentFrom(len(own.Command) > 0, len(cfg.Agent.Command) > 0)
	nameFrom := agentFrom(own.Name != "", cfg.Agent.Name != "")
	harnessFrom := func(name string) string {
		_, project := own.Harnesses[name]
		_, everywhere := cfg.Agent.Harnesses[name]
		return agentFrom(project, everywhere)
	}
	if a.jsonOut {
		cmd := agent.Command
		if cmd == nil {
			cmd = []string{}
		}
		// the program and arguments keep the names they had while shown
		// embedded harness.Host
		type shown struct {
			Program       string            `json:"Program"`
			Args          []string          `json:"Args"`
			Guard         string            `json:"guard,omitempty"`
			DenyProtected bool              `json:"deny_protected,omitempty"`
			Images        map[string]string `json:"images,omitempty"`
			Network       string            `json:"network,omitempty"`
			Login         string            `json:"login,omitempty"`
			Outside       []string          `json:"outside,omitempty"`
			From          string            `json:"from,omitempty"`
		}
		hosts := map[string]shown{}
		for _, name := range settable() {
			set, h := agent.Harnesses[name], a.harnessHost(agent, name)
			s := shown{Program: h.Program, Args: h.Args, Guard: set.Guard, DenyProtected: set.DenyProtected, Outside: set.Outside, From: harnessFrom(name)}
			if name == config.ContainerHarness {
				s.Images, s.Network, s.Login = set.ShownImages(), set.Network, set.Login
			}
			hosts[name] = s
		}
		if set, ok := agent.Harnesses[harness.Command]; ok {
			s := shown{Guard: set.Guard, DenyProtected: set.DenyProtected, Outside: set.Outside, From: harnessFrom(harness.Command)}
			if len(agent.Command) > 0 {
				s.Program, s.Args = agent.Command[0], agent.Command[1:]
			}
			hosts[harness.Command] = s
		}
		shownAll := map[string]any{"command": cmd, "command_from": commandFrom, "name": agent.Name, "name_from": nameFrom,
			"harnesses": hosts, "auto_restarts": cfg.Agent.AutoRestartLimit(), "enabled_here": enabled, "project": here}
		if len(lifted) > 0 {
			shownAll["budget_cleared"] = lifted
		}
		return a.printJSON(shownAll)
	}
	if len(agent.Command) == 0 {
		fmt.Fprintln(a.out, "no command is set, so a story whose agent and project name no harness starts on flai's built-in default, claude-code (ADR-0261); flai serve agent set -- <program> [args...] sets one")
	} else {
		fmt.Fprintf(a.out, "command: %s%s\n  run as it stands, in the project's directory, never through a shell\n", quoteArgs(agent.Command), fromSaid(commandFrom))
	}
	if agent.Name != "" {
		fmt.Fprintf(a.out, "name: %s%s\n", agent.Name, fromSaid(nameFrom))
	}
	if set, ok := agent.Harnesses[harness.Command]; ok {
		fmt.Fprintf(a.out, "harness %s: the command above%s\n", harness.Command, fromSaid(harnessFrom(harness.Command)))
		printLifts(a, set)
	}
	for _, name := range settable() {
		h := a.harnessHost(agent, name)
		fmt.Fprintf(a.out, "harness %s: %s%s\n", name, quoteArgs(append([]string{h.Program}, h.Args...)), fromSaid(harnessFrom(name)))
		printLifts(a, agent.Harnesses[name])
		if name == config.ContainerHarness {
			printContainer(a, agent.Harnesses[name])
		}
	}
	fmt.Fprintf(a.out, "automatic restarts: %s\n  of a story's agent that ends with its story in progress, on every project; set with --auto-restarts\n", autoRestartsSaid(cfg.Agent))
	if here != "" {
		if enabled {
			fmt.Fprintln(a.out, "the agent action is on for this project; flai serve disable agent turns it off")
		} else {
			fmt.Fprintln(a.out, "the agent action is off for this project; flai serve enable agent turns it on")
		}
	}
	for _, l := range lifted {
		fmt.Fprintf(a.out, "cleared the budget state of %s: %s; the next start on it is not held\n", l.Key, l.Why)
	}
	return nil
}

// printLifts prints what the operator lifts on a harness, under its line in
// show: guard none and deny_protected, each with what it gives up (ADR-0131),
// and the outside folders it grants (ADR-0289).
func printLifts(a *app, set config.HarnessHost) {
	if set.Guard == config.GuardNone {
		fmt.Fprintln(a.out, "  guard: none, so a story's sub-agents run without flai guard where the harness cannot run it, and could move items, write threads, or commit")
	}
	if set.DenyProtected {
		fmt.Fprintln(a.out, "  deny_protected: your arguments deny the protected paths, so flai starts the agent where the harness cannot ask before a protected write")
	}
	if len(set.Outside) > 0 {
		fmt.Fprintf(a.out, "  outside: %s, which a story's agent and its sub-agents may read and write; --outside \"\" clears them\n", strings.Join(set.Outside, ", "))
	}
}

// printContainer prints the container harness's settings under its line in
// show: the image for each inner harness, the network, and the login folder,
// each saying whether it is the default.
func printContainer(a *app, set config.HarnessHost) {
	images := set.ShownImages()
	for _, inner := range slices.Sorted(maps.Keys(images)) {
		from := "set"
		if set.Images[inner] == "" {
			from = "default, built by make agent-images"
		}
		fmt.Fprintf(a.out, "  image %s: %s (%s)\n", inner, images[inner], from)
	}
	if set.Network == "" {
		fmt.Fprintln(a.out, "  network: Docker's default bridge, the host reached through Docker's host gateway")
	} else {
		fmt.Fprintf(a.out, "  network: %s\n", set.Network)
	}
	if set.Login == "" {
		fmt.Fprintln(a.out, "  login: not set; set the folder holding your Claude login with --login ~/.claude")
	} else {
		fmt.Fprintf(a.out, "  login: %s, mounted read-only\n", set.Login)
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
	// the project's own command, name, and harnesses over the host-wide ones
	// (ADR-0261)
	agent := cfg.Agent.ForProject(root)
	hosts := map[string]harness.Host{}
	for _, name := range settable() {
		hosts[name] = a.harnessHost(agent, name)
	}
	if len(agent.Command) > 0 {
		// the command's outside folders are its harness entry's (ADR-0289)
		hosts[harness.Command] = harness.Host{Program: agent.Command[0], Args: agent.Command[1:], Outside: slices.Clone(agent.Harnesses[harness.Command].Outside)}
	}
	lifts := map[string]serve.HarnessLifts{}
	for name, set := range agent.Harnesses {
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
		Command:        agent.Command,
		Harnesses:      hosts,
		Providers:      cfg.Agent.Providers,
		Name:           agent.Name,
		AutoRestarts:   cfg.Agent.AutoRestartLimit(),
		Lifts:          lifts,
		AutoApprove:    cfg.ActionEnabled(hostapi.ActionAutoApprove, root),
		ForgeTokenEnvs: forgeTokens,
		RoleActions:    defined,
	}
}
