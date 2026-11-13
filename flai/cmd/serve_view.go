package cmd

import (
	"context"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/config"
	"github.com/bytepunx/system-flow/flai/internal/harness"
	"github.com/bytepunx/system-flow/flai/internal/hostapi"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/perf"
	"github.com/bytepunx/system-flow/flai/internal/serve"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

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

// hostHarnesses are the harnesses of agent, the settings that apply to a
// project or the host-wide ones, for the settings page: each one's program
// and arguments, with set saying whether the settings give either rather
// than the adapter's defaults; its guard and deny_protected (ADR-0131); and
// whether its adapter runs flai guard and can ask before a protected write,
// where neither lifts anything; and the container harness's image for each
// inner harness, the defaults filled in, its network, and its login folder
// (S-0508). from, when given, says for each where its entry comes from
// (ADR-0261). The command's program and arguments are flai serve agent set's
// command, never set here. A harness with outside folders, the folders
// outside the project its agents reach, gives them (ADR-0289).
func (a *app) hostHarnesses(agent config.AgentStart, from func(name string) string) map[string]any {
	out := map[string]any{}
	for _, n := range append(settable(), harness.Command) {
		var h harness.Host
		if n != harness.Command {
			h = a.harnessHost(agent, n)
		} else if len(agent.Command) > 0 {
			h = harness.Host{Program: agent.Command[0], Args: agent.Command[1:]}
		}
		args := h.Args
		if args == nil {
			args = []string{}
		}
		set := agent.Harnesses[n]
		caps := harness.Adapters[n].Capabilities()
		entry := map[string]any{"program": h.Program, "args": args, "set": set.Program != "" || set.Args != nil,
			"guard": set.Guard, "deny_protected": set.DenyProtected, "guard_capable": caps.Guard, "ask_capable": caps.Ask}
		if len(set.Outside) > 0 {
			entry["outside"] = set.Outside
		}
		if from != nil {
			entry["from"] = from(n)
		}
		if n == config.ContainerHarness {
			entry["images"], entry["network"], entry["login"] = set.ShownImages(), set.Network, set.Login
		}
		out[n] = entry
	}
	return out
}

// hostSettings is the host's configuration as it applies to the project at
// root, for the dashboard's settings page (S-0105) and its Host page
// (S-0489). The tokens are not in it: a rotation answers with the new one,
// once. It records three phases on ctx (S-0605): settings.config, the host
// configuration and what it says of the project; settings.served, the
// projects flai serve serves; and settings.project, the project's own
// manifest and MCP server.
func (a *app) hostSettings(ctx context.Context, root string) any {
	done := perf.Track(ctx, "settings.config")
	cfg, _, err := a.loadConfig()
	if err != nil {
		done()
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
	// the project's command, name, and harnesses over every project's, each
	// saying where it comes from, and every project's own beside them, the
	// host-wide fallback (ADR-0261)
	agent, own := cfg.Agent.ForProject(root), cfg.Agent.Project(root)
	command := nonNil(agent.Command)
	harnesses := a.hostHarnesses(agent, func(name string) string {
		_, project := own.Harnesses[name]
		_, everywhere := cfg.Agent.Harnesses[name]
		return agentFrom(project, everywhere)
	})
	everywhere := map[string]any{"name": cfg.Agent.Name, "command": nonNil(cfg.Agent.Command),
		"harnesses": a.hostHarnesses(cfg.Agent.ForProject(""), nil)}
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
	done()
	done = perf.Track(ctx, "settings.served")
	projects := a.servedView()
	done()
	out := map[string]any{
		"actions": actions,
		"agent": map[string]any{"name": agent.Name, "command": command,
			"name_from": agentFrom(own.Name != "", cfg.Agent.Name != ""), "command_from": agentFrom(len(own.Command) > 0, len(cfg.Agent.Command) > 0),
			"harnesses": harnesses, "providers": providers, "everywhere": everywhere,
			"image": agentImageShown(cfg.Agent, root)},
		"allowed_hosts": config.NormalizeAllowedHosts(cfg.Dashboard.AllowedHosts),
		"checks":        map[string]any{"commands": checks, "timeout_minutes": timeout},
		"import_roots":  roots,
		"projects":      projects,
	}
	defer perf.Track(ctx, "settings.project")()
	if repo, err := workitem.Open(root); err == nil {
		out["default_agent"] = repo.Manifest.Agent
		// the shared paths as the manifest has them now (S-0295, ADR-0096)
		out["shared"] = nonNil(repo.Manifest.Claims.Shared)
		out["manifest_checks"] = repo.Manifest.Checks
		out["planning"] = planningTriggers(repo.Manifest.Planning, cfg.ActionEnabled(hostapi.ActionPlan, root), a.now())
		out["strategic"] = strategicSettings(root, repo.Manifest, cfg.ActionEnabled(hostapi.ActionSettings, root))
		// bind is mcp.bind's first address as written, with bind_set saying
		// whether it is set (S-0506); binds is its addresses as written, or
		// the default, and listen what the running server listens on, from
		// its state, so that the settings ask docker nothing (S-0732)
		b := mcpBindOf(cfg)
		mcp := map[string]any{"running": false, "bind": b.Bind, "bind_set": b.BindSet, "binds": b.Binds}
		if st, running := readMCPState(repo, a.now()); running {
			mcp["running"], mcp["url"], mcp["pid"], mcp["listen"] = true, st.URL, st.PID, st.listening()
		}
		out["mcp"] = mcp
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
// is on for the project, which settings.manifest needs; enable, the
// command that turns it on, run in a shell on the host; and roles, each role
// the project at root runs that is not built-in with its agent, schedule,
// and declared permissions, which settings.manifest writes as it writes the
// others (S-0420).
func strategicSettings(root string, m manifest.Manifest, editable bool) map[string]any {
	return map[string]any{"editable": editable, "enable": hostapi.EnableCommand(hostapi.ActionSettings),
		"settings": m.SettingValues(), "roles": hostapi.DefinedRoleSettings(root, m)}
}

// servedView is what the settings page shows of the projects flai serve
// serves, and of those below the folders named for import that it does not,
// with why (S-0122): what flai serve project list says, without the
// dashboard's address or the credential's path. It reads the configuration
// once and walks no folder for candidates, which it does not show; while
// flai serve runs it walks none at all, and works out each project's below
// from its path (S-0605).
func (a *app) servedView() map[string]any {
	cfg, _, err := a.loadConfig()
	if err != nil {
		return map[string]any{"error": err.Error()}
	}
	l, err := a.servedListing(cfg.ImportRoots, false)
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
		below := servedBelowIn(cfg.ImportRoots, l.folder, p.Root)
		served = append(served, project{Key: p.Key, Name: p.Name, Root: p.Root, From: p.From, Below: below,
			State: p.State, Since: p.Since, LastError: p.LastError, Reason: p.Reason, Settings: cfg.ActionEnabled(hostapi.ActionSettings, p.Root)})
	}
	unserved := make([]project, 0, len(l.Unserved))
	for _, f := range l.Unserved {
		unserved = append(unserved, project{Key: f.Key, Name: f.Name, Root: f.Root, Reason: f.Reason, Removed: f.Removed, Settings: cfg.ActionEnabled(hostapi.ActionSettings, f.Root)})
	}
	return map[string]any{"running": l.Running, "served": served, "unserved": unserved}
}
