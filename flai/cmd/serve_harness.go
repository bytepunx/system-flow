package cmd

import (
	"context"
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/spf13/cobra"

	"github.com/bytepunx/system-flow/flai/internal/config"
	"github.com/bytepunx/system-flow/flai/internal/harness"
	"github.com/bytepunx/system-flow/flai/internal/serve"
)

// harnessHost is the program and arguments a harness runs with under agent,
// the settings that apply to a project (ADR-0261): the operator's where set,
// the adapter's defaults where not (S-0104); the
// container harness's images, network, and login as the operator sets them,
// whose defaults the adapter applies (S-0508); and the harness's outside
// folders (ADR-0289).
func (a *app) harnessHost(agent config.AgentStart, name string) harness.Host {
	h := harness.Adapters[name].DefaultHost()
	set := agent.Harnesses[name]
	if set.Program != "" {
		h.Program = set.Program
	}
	if set.Args != nil {
		h.Args = *set.Args
	}
	h.Images, h.Network, h.Login = maps.Clone(set.Images), set.Network, set.Login
	h.Outside = slices.Clone(set.Outside)
	return h
}

// harnessIdentity is what a start on the harness name runs as under agent,
// as its budget state is keyed (serve.BudgetKey): the command's program and
// arguments for the command harness, else harnessHost's.
func (a *app) harnessIdentity(agent config.AgentStart, name string) harness.Host {
	if name != harness.Command {
		return a.harnessHost(agent, name)
	}
	if len(agent.Command) == 0 {
		return harness.Host{}
	}
	return harness.Host{Program: agent.Command[0], Args: agent.Command[1:]}
}

// identityChanges says what of the program, the arguments, and the login
// changed from was to now, such as "program from rmhclod to meclod"; none
// when none did.
func identityChanges(was, now harness.Host) []string {
	var out []string
	if was.Program != now.Program {
		out = append(out, fmt.Sprintf("program from %s to %s", orNone(was.Program), orNone(now.Program)))
	}
	if !slices.Equal(was.Args, now.Args) {
		out = append(out, fmt.Sprintf("arguments from %s to %s", orNone(quoteArgs(was.Args)), orNone(quoteArgs(now.Args))))
	}
	if was.Login != now.Login {
		out = append(out, fmt.Sprintf("login from %s to %s", orNone(was.Login), orNone(now.Login)))
	}
	return out
}

// budgetLift is a budget state a change of a harness's identity cleared, and
// why, as the command that made the change prints it.
type budgetLift struct {
	Key string `json:"key"`
	Why string `json:"why"`
}

// liftBudgets lifts the hold the old identity of the harness name, was, puts
// on its starts, once by, such as flai serve agent harness, has saved a
// change to now for the project at root, or for every project with root ""
// (S-0720): when the program, the arguments, or the login changed, it clears
// the state kept under was's key and one kept under the harness's name alone
// from before S-0720, journals each it removed under serve.budget, and sends
// its notice through notify when a channel is set, a failure of either
// warned of, as flai serve budget clear does. A change of nothing of those
// clears nothing. A clear that fails is an error that says the setting was
// saved and how to clear by hand.
func (a *app) liftBudgets(ctx context.Context, notify config.Notify, by, name, root string, was, now harness.Host) ([]budgetLift, error) {
	changes := identityChanges(was, now)
	if len(changes) == 0 {
		return nil, nil
	}
	// never no keys: ClearBudgets clears every state with none
	keys := []string{serve.BudgetKey(name, was)}
	if keys[0] != name {
		keys = append(keys, name)
	}
	dir := a.serveDir()
	cleared, err := dir.ClearBudgets(keys)
	if err != nil {
		hand := make([]string, len(keys))
		for i, k := range keys {
			hand[i] = serve.BudgetClearCommand(k)
		}
		return nil, fmt.Errorf("the setting was saved, but the budget state of the %s harness's old program and login was not cleared: %w; clear it by hand with %s",
			name, err, strings.Join(hand, " and "))
	}
	if len(cleared.Removed) == 0 {
		return nil, nil
	}
	scope, project := "every project", ""
	if root != "" {
		scope = root
		registry, err := dir.Projects()
		if err != nil {
			a.logger().Warn("budget clear names no project: the registry is not read", "component", "serve", "err", err.Error())
		}
		for _, e := range registry {
			if e.Root == root {
				project = e.Key
			}
		}
	}
	why := fmt.Sprintf("%s changed the %s harness's %s for %s", by, name, strings.Join(changes, " and "), scope)
	a.announceBudgetClears(ctx, cleared.Removed, budgetClearing{why: why, by: by, project: project, root: root}, notify)
	out := make([]budgetLift, 0, len(cleared.Removed))
	for _, b := range cleared.Removed {
		out = append(out, budgetLift{Key: b.Provider, Why: why})
	}
	return out, nil
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
// --deny-protected, and --reset: its program is flai serve agent set's. Only
// the container harness takes --image, --network, and --login (S-0508).
func newServeAgentHarnessCmd(a *app) *cobra.Command {
	var program, guard, network, login, project string
	var images, outside []string
	var reset, denyProtected, all bool
	c := &cobra.Command{
		Use:   "harness [<name>] [--program <path>] [--guard none|\"\"] [--deny-protected=true|false] [--image <inner>=<image>]... [--network <name>] [--login <path>] [--outside <dir>]... [--reset] [--project <root> | --all-projects] [-- args...]",
		Short: "Set the program a harness is, the arguments that say what its agent may do, and what you lift on it",
		Long: `A story names its harness, such as claude-code; you say here what that is on
this host. --program is what is run, and the arguments after -- replace the
adapter's defaults, which say what the agent may do (for claude-code:
--permission-mode acceptEdits --allowedTools Bash,Monitor,mcp__flai). "--" with
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

The opencode harness runs opencode run with no arguments of flai's: it can
neither run flai guard nor ask, so it starts only with the auto-approve host
action on or --deny-protected, and --deny-protected has flai give OpenCode a
deny rule for each protected path (ADR-0260).

The container harness runs an inner harness, claude-code, in a Docker
container, and only it takes --image, --network, and --login.
--image <inner>=<image> is the image the inner harness runs in,
flai-agent-<inner>:local (make agent-images) unless set; repeat it for each
inner harness. --network is the Docker network the container joins, Docker's
default bridge with the host reached through Docker's host gateway unless
set; --network "" goes back to it. --login is the folder holding your Claude
login, mounted read-only, an existing folder given as an absolute path or
under ~/; --login "" clears it.

--outside grants a folder outside the project that a story's agent and its
sub-agents may read and write, on every harness: the claude-code harness
passes it as --add-dir and the container harness mounts it at its own path
(ADR-0289). Repeat it once per folder; the folders given replace the list,
and --outside "" clears it. Each must be an existing folder, given absolute
or under ~/, and is kept absolute. It refuses your home folder, ~/.flai and
any folder under or holding it, and any folder that holds the project or
that the project holds; with --all-projects, every project this host knows.

The command harness takes only --guard, --deny-protected, --outside, and
--reset: its program and arguments are set with flai serve agent set.

A change of the program, the arguments, or the login a harness runs with
lifts the hold of the old one's budget state: the state kept under the
harness and its old program and login, and one kept under the harness's name
alone from before S-0720, are cleared, journaled, and sent as a notice like
any other clear, so that the next start goes. A change of anything else
clears nothing.

A harness's entry is the project's in the working directory, or of
--project <root>, and every project's, the host-wide fallback a project that
sets none of its own reads, only with --all-projects (ADR-0261). A project's
entry replaces the host-wide entry of that harness whole: the first change
for a project starts from the entry it reads, and --reset for a project
removes its entry, so that it reads every project's again.`,
		Example: `  flai serve agent harness
  flai serve agent harness claude-code --program /opt/claude/bin/claude
  flai serve agent harness claude-code -- --permission-mode acceptEdits --allowedTools Bash,Edit,mcp__flai
  flai serve agent harness claude-code --reset
  flai serve agent harness command --guard none
  flai serve agent harness command --guard ""
  flai serve agent harness command --deny-protected
  flai serve agent harness command --deny-protected=false
  flai serve agent harness opencode --program /opt/opencode/bin/opencode --deny-protected
  flai serve agent harness container --image claude-code=flai-agent-claude-code:local --login ~/.claude
  flai serve agent harness container --network flai-agents
  flai serve agent harness claude-code --outside ~/src/other --outside /srv/data
  flai serve agent harness claude-code --outside ""
  flai serve agent harness claude-code --all-projects --program /usr/local/bin/claude
  flai serve agent harness claude-code --project /home/me/src/other --reset`,
		RunE: func(cmd *cobra.Command, args []string) error {
			guardSet, denySet := cmd.Flags().Changed("guard"), cmd.Flags().Changed("deny-protected")
			imageSet, networkSet, loginSet := cmd.Flags().Changed("image"), cmd.Flags().Changed("network"), cmd.Flags().Changed("login")
			containerSet := imageSet || networkSet || loginSet
			outsideSet := cmd.Flags().Changed("outside")
			dash := cmd.ArgsLenAtDash()
			names, rest := args, []string(nil)
			if dash >= 0 {
				names, rest = args[:dash], args[dash:]
			}
			known := append(settable(), harness.Command)
			if len(names) == 0 {
				if dash >= 0 || program != "" || reset || guardSet || denySet || containerSet || outsideSet {
					return fmt.Errorf("name the harness: %s", strings.Join(known, ", "))
				}
				return a.showAgentCommand("")
			}
			if len(names) > 1 {
				return fmt.Errorf("one harness at a time; arguments for it go after --")
			}
			name := names[0]
			if !slices.Contains(known, name) {
				return fmt.Errorf("there is no harness %q to set; there are: %s", name, strings.Join(known, ", "))
			}
			if name == harness.Command && (program != "" || dash >= 0) {
				return fmt.Errorf("the command's program and arguments are set with flai serve agent set -- <program> [args...]; flai serve agent harness command takes only --guard, --deny-protected, --outside, and --reset")
			}
			if containerSet && name != config.ContainerHarness {
				return fmt.Errorf("--image, --network, and --login are the %s harness's, which %s is not; set them with flai serve agent harness %s", config.ContainerHarness, name, config.ContainerHarness)
			}
			if reset && (program != "" || dash >= 0 || guardSet || denySet || containerSet || outsideSet) {
				return fmt.Errorf("--reset takes neither --program, --guard, --deny-protected, --image, --network, --login, --outside, nor arguments")
			}
			setImages, err := parseImages(images)
			if err != nil {
				return err
			}
			if loginSet && login != "" {
				if login, err = loginFolder(login); err != nil {
					return err
				}
			}
			root, err := a.agentScope(project, all)
			if err != nil {
				return err
			}
			cfg, path, err := a.loadConfig()
			if err != nil {
				return err
			}
			was := a.harnessIdentity(cfg.Agent.ForProject(root), name)
			// a project's first change starts from the entry it reads
			h := cfg.Agent.ForProject(root).Harnesses[name]
			if reset {
				h = config.HarnessHost{}
			}
			if imageSet {
				h.Images = maps.Clone(h.Images)
				if h.Images == nil {
					h.Images = map[string]string{}
				}
				maps.Copy(h.Images, setImages)
			}
			if networkSet {
				h.Network = network
			}
			if loginSet {
				h.Login = login
			}
			if outsideSet {
				if h.Outside, err = a.outsideFolders(name, outside, root, cfg); err != nil {
					return err
				}
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
			if root == "" {
				cfg.Agent, err = cfg.Agent.WithHarness(name, h)
			} else {
				p := cfg.Agent.Project(root)
				p.Harnesses = maps.Clone(p.Harnesses)
				if p.Harnesses == nil {
					p.Harnesses = map[string]config.HarnessHost{}
				}
				p.Harnesses[name] = h // an entry that sets nothing is removed
				cfg.Agent, err = cfg.Agent.WithProject(root, p)
			}
			if err != nil {
				return err
			}
			if err := config.Save(path, cfg); err != nil {
				return err
			}
			lifted, err := a.liftBudgets(cmd.Context(), cfg.Notify, "flai serve agent harness", name, root, was, a.harnessIdentity(cfg.Agent.ForProject(root), name))
			if err != nil {
				return err
			}
			return a.showAgentCommand(root, lifted...)
		},
	}
	c.Flags().StringVar(&project, "project", "", "the main checkout of the project whose harness this is, not the one in the working directory")
	c.Flags().BoolVar(&all, "all-projects", false, "every project's harness, the host-wide fallback a project that sets none of its own reads")
	c.Flags().StringVar(&program, "program", "", "the program the harness is on this host")
	c.Flags().StringVar(&guard, "guard", "", `none lets a story's sub-agents run unguarded on a harness that cannot run flai guard; "" puts the refusal back`)
	c.Flags().BoolVar(&denyProtected, "deny-protected", false, "your arguments for the harness deny the protected paths, so flai starts the agent on a harness that cannot ask before a protected write")
	c.Flags().StringArrayVar(&images, "image", nil, "the container harness's image for an inner harness, as <inner>=<image>; repeat for each (flai-agent-<inner>:local unless set)")
	c.Flags().StringVar(&network, "network", "", `the Docker network the container harness's agents join; "" is Docker's default bridge`)
	c.Flags().StringVar(&login, "login", "", `the folder holding your Claude login, which the container harness mounts read-only; "" clears it`)
	c.Flags().StringArrayVar(&outside, "outside", nil, `a folder outside the project a story's agent and its sub-agents may read and write, as --add-dir or a container mount (ADR-0289); repeat for each, "" clears them`)
	c.Flags().BoolVar(&reset, "reset", false, "go back to the adapter's program and arguments, and lift nothing")
	return c
}

// parseImages reads the --image flags, each <inner>=<image>, into the images
// they set; config refuses an inner harness or an image it does not take.
func parseImages(flags []string) (map[string]string, error) {
	out := map[string]string{}
	for _, f := range flags {
		inner, ref, ok := strings.Cut(f, "=")
		if !ok || inner == "" {
			return nil, fmt.Errorf("--image %q is not <inner>=<image>: give the inner harness and its image, as --image claude-code=%s", f, config.DefaultImage("claude-code"))
		}
		out[inner] = ref
	}
	return out, nil
}

// loginFolder is the --login folder as an absolute path, ~/ expanded,
// refusing one that is not an existing folder: Docker would create a missing
// one, empty, for the mount.
func loginFolder(p string) (string, error) {
	expanded := config.ExpandHome(p)
	if strings.HasPrefix(expanded, "~") {
		return "", fmt.Errorf("--login %q: give the folder holding your Claude login as an absolute path or under ~/, as --login ~/.claude", p)
	}
	abs, err := filepath.Abs(expanded)
	if err != nil {
		return "", fmt.Errorf("--login %q: resolve it to an absolute path: %w", p, err)
	}
	if info, err := os.Stat(abs); err != nil || !info.IsDir() {
		return "", fmt.Errorf("--login %q: %s is not a folder; give the folder holding your Claude login, as --login ~/.claude", p, abs)
	}
	return abs, nil
}

// outsideFolders is the --outside folders of the harness name, checked by
// config.OutsideFolders against the home folder and the project at root, or
// with root "" (--all-projects) every project this host knows: those flai
// serve serves, the one in the working directory, and those with settings in
// cfg. "" alone, or none at all, clears the list (ADR-0289).
func (a *app) outsideFolders(name string, flags []string, root string, cfg config.Config) ([]string, error) {
	dirs := slices.DeleteFunc(slices.Clone(flags), func(d string) bool { return d == "" })
	if len(dirs) == 0 {
		return nil, nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("--outside: find your home folder, which no outside folder may be: %w; set HOME", err)
	}
	roots := []string{root}
	if root == "" {
		roots = append(a.servedRoots(), slices.Collect(maps.Keys(cfg.Agent.Projects))...)
		for _, list := range cfg.HostActions {
			roots = append(roots, slices.DeleteFunc(slices.Clone(list), func(r string) bool { return !filepath.IsAbs(r) })...)
		}
	}
	return config.OutsideFolders(name, dirs, home, roots)
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
