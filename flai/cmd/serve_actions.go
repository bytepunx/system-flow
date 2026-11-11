package cmd

import (
	"bufio"
	"context"
	"encoding/json"
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
	"github.com/bytepunx/system-flow/flai/internal/execx"
	"github.com/bytepunx/system-flow/flai/internal/hostapi"
	"github.com/bytepunx/system-flow/flai/internal/kluster"
	"github.com/bytepunx/system-flow/flai/internal/perf"
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
		// Agent times each part of the state it builds as a phase of the
		// request, inside hostapi's agent.state (S-0603). Without finished it
		// leaves out the runs of finished work and every past run (S-0611).
		Agent: func(ctx context.Context, root string, finished bool) any {
			cfg, _, cfgErr := a.loadConfig()
			done := perf.Track(ctx, "agent.states")
			st := a.serveDir().AgentStates()[root]
			done()
			command := ""
			if agent := cfg.Agent.ForProject(root); len(agent.Command) > 0 {
				command = filepath.Base(agent.Command[0]) // its name, never its arguments
			}
			plans := st.Plans // the newest planner run for each item, by its ID (S-0208)
			if plans == nil {
				plans = map[string]*serve.AgentRun{}
			}
			// every role the project runs, the built-in ones included, by its
			// name: its host action's state, its run under way, its past runs,
			// and held (S-0420)
			done = perf.Track(ctx, "agent.roles")
			defined := serve.RoleStates(st, hostapi.ProjectRoles(root), func(action string) bool {
				return cfgErr == nil && cfg.ActionEnabled(action, root)
			})
			done()
			// the orchestrator's and the analyzer's runs before their newest,
			// newest first (S-0228), which the dashboard's pages list
			pastOrchestrators, pastAnalyzers := st.PastOrchestrators, st.PastAnalyzers
			if pastOrchestrators == nil {
				pastOrchestrators = []*serve.AgentRun{}
			}
			if pastAnalyzers == nil {
				pastAnalyzers = []*serve.AgentRun{}
			}
			// the newest orchestrator run, null before the first (S-0218), with
			// held while the operator holds it stopped (S-0228), and the newest
			// analyzer run, its focus and report among it (S-0223), and the
			// newest curator run, its since and entries among it (S-0431)
			done = perf.Track(ctx, "agent.activity")
			stories := serve.Activity(ctx, root, st, command != "")
			done()
			if !finished {
				// the open and live runs alone (S-0611): the board, the activity
				// page, and a story's agent panel show no more, and the pages
				// that list past runs ask with finished
				done = perf.Track(ctx, "agent.unfinished")
				stories, plans = serve.Unfinished(root, stories, plans)
				done()
				pastOrchestrators, pastAnalyzers = []*serve.AgentRun{}, []*serve.AgentRun{}
				for name, r := range defined {
					r.Runs = []*serve.AgentRun{}
					defined[name] = r
				}
			}
			out := map[string]any{"command": command, "running": st.Running, "last": st.Last, "waiting": st.Waiting, "stories": stories, "plans": plans,
				"orchestrator": st.Orchestrator, "analyzer": st.Analyzer, "curator": st.Curator, "roles": defined,
				"past_orchestrators": pastOrchestrators, "past_analyzers": pastAnalyzers}
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
	var trailers []string
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
--all-projects removes it everywhere.

Enabling kluster for a project whose system-flow.yaml has no kluster: block
writes one, naming only config: kluster.yaml, and writes kluster.yaml at the
repository root when it is absent, naming the cluster, the provider k3d, and
the profile spire; it commits those two files alone, with each --trailer,
unless the project sets dashboard.autocommit: false. The block names no
deploy, so flai deploys nothing until kluster.deploy is set. --all-projects
writes no block.`,
		Example: `  flai serve enable push
  flai serve enable kluster
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
			// a project with no kluster: block is given one before the action
			// is on, so that it has a cluster to run (S-0586)
			var scaffold *klusterScaffold
			if on && !all && action == hostapi.ActionKluster {
				if scaffold, err = a.scaffoldKluster(root, trailers); err != nil {
					return err
				}
			}
			next := cfg.WithAction(action, root, on)
			if err := config.Save(path, next); err != nil {
				return err
			}
			still := !on && !all && next.ActionEnabled(action, root)
			if a.jsonOut {
				out := map[string]any{"action": action, "enabled": next.ActionEnabled(action, root), "for": root, "config": path}
				if scaffold != nil {
					out["kluster"] = scaffold
				}
				return a.printJSON(out)
			}
			switch {
			case on:
				if scaffold != nil {
					scaffold.print(a)
				}
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
	if on {
		c.Flags().StringArrayVar(&trailers, "trailer", nil, "a trailer line for the commit of the kluster: block enabling kluster writes (repeatable)")
	}
	return c
}

// klusterScaffold is what enabling kluster wrote for a project with no
// kluster: block, and the commit that holds it.
type klusterScaffold struct {
	kluster.Scaffold
	// Commit is the short hash of the commit of Wrote; empty when nothing was
	// written or the project does not commit what flai writes.
	Commit string `json:"commit,omitempty"`
}

// print says what enabling kluster wrote and the commit.
func (s *klusterScaffold) print(a *app) {
	if len(s.Wrote) == 0 {
		return
	}
	fmt.Fprintf(a.out, "wrote a kluster: block naming %s, with no deploy, so nothing is deployed until kluster.deploy is set: %s\n", s.Config, strings.Join(s.Wrote, ", "))
	if s.Kept {
		fmt.Fprintf(a.out, "  kept the %s that was there\n", s.Config)
	}
	if s.Commit != "" {
		fmt.Fprintf(a.out, "  committed as %s\n", s.Commit)
	} else {
		fmt.Fprintln(a.out, "  not committed: dashboard.autocommit is false, or the commit failed; commit them by hand")
	}
}

// scaffoldKluster gives the project at root a kluster: block and its config
// when its manifest has none, and commits what it wrote; nil when it had a
// block already (ADR-0254, S-0586).
func (a *app) scaffoldKluster(root string, trailers []string) (*klusterScaffold, error) {
	repo, err := workitem.Open(root)
	if err != nil {
		return nil, err
	}
	if repo.Manifest.Kluster != nil {
		return nil, nil
	}
	s, err := kluster.Init(root, repo.Manifest)
	if err != nil {
		return nil, err
	}
	out := &klusterScaffold{Scaffold: s}
	if len(s.Wrote) > 0 {
		out.Commit = a.commitScaffold(repo, trailers, s.Wrote)
	}
	return out, nil
}

// commitScaffold commits the paths the kluster scaffold wrote, and nothing
// else, unless the project sets dashboard.autocommit: false; it says the
// commit, or nothing, warning on standard error when the commit fails.
func (a *app) commitScaffold(repo *workitem.Repo, trailers, paths []string) string {
	if !repo.Manifest.Autocommit() {
		return ""
	}
	msg := "chore: give the project a kluster: block and its config"
	if len(trailers) > 0 {
		msg += "\n\n" + strings.Join(trailers, "\n")
	}
	fail := func(err error) string {
		a.logger().Warn("kluster scaffold written but not committed", "component", "cmd", "paths", strings.Join(paths, ","), "error", err, "fix", "commit the paths by hand")
		return ""
	}
	if _, err := a.runner.Run(repo.Root, "git", append([]string{"add", "--"}, paths...)...); err != nil {
		return fail(err)
	}
	if _, err := a.runner.Run(repo.Root, "git", append([]string{"commit", "-q", "-m", msg, "--"}, paths...)...); err != nil {
		return fail(err)
	}
	sha, err := a.runner.Run(repo.Root, "git", "rev-parse", "--short", "HEAD")
	if err != nil {
		return fail(err)
	}
	return strings.TrimSpace(sha)
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
				if e.Verdict != "" {
					fmt.Fprintf(a.out, "  verdict %s", e.Verdict)
				}
				if e.DurationMS > 0 {
					fmt.Fprintf(a.out, " · %s", (time.Duration(e.DurationMS) * time.Millisecond).Round(time.Second))
				}
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

// serveHost is the host flai serve itself answers with: a.host() and the
// supervisor of the component dev servers (S-0438, ADR-0163), which runs at
// most components.max_servers per project and stops each one idle for
// components.idle, and each of a story's worktree once the story leaves
// review (S-0442), until ctx ends, and then every one.
// It also runs each project's kluster up and kluster down that the
// dashboard asks for (S-0580, ADR-0254), until ctx ends.
// stop ends it as ctx ending does and waits until every dev server has
// stopped and every kluster job has ended; flai serve calls it as it stops.
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
	jobs := hostapi.NewKlusterJobs(runCtx, hostapi.KlusterJobsOptions{Now: a.now, Logger: a.logger()})
	host = a.host()
	host.Components = comps
	host.Kluster = jobs
	return host, func() {
		cancel()
		<-done
		jobs.Wait()
	}
}
