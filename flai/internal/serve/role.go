package serve

import (
	"context"
	"errors"
	"fmt"
	"path"
	"slices"
	"strings"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/harness"
	"github.com/bytepunx/system-flow/flai/internal/hostapi"
	"github.com/bytepunx/system-flow/flai/internal/inbox"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/roles"
	"github.com/bytepunx/system-flow/flai/internal/threads"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// A role a project defines, started on the operator's word (S-0417,
// ADR-0207).
//
// flai role run, the MCP tool role_run, and role.run through the host API
// start a role that a file under the design folder's workflows/roles
// defines through RunRole, behind the host action its file names: in the
// project's main checkout, with the role's agent (roles.<name>.agent over
// agent) and the harness's request for its file, on the item it is given
// when its file takes one. Like the planner it runs in a process of its own
// and records the run in serve/agents.json, under roles by its name and its
// item, with its log beside the others' under its name (openLog); the
// serving flai settles the run once its process is gone, judges how it went
// as a planner run's is, journals its end, and logs its activity in
// wip/agents/<name>.md, which charges nothing. It is no story's agent: it is
// not counted against the in-progress limit, and holds no story back.

// MethodRole is the method a defined role's start and end are journalled
// under, with the host action its file names as the entry's action.
const MethodRole = "serve.role"

// builtinRoles are the built-in roles, by name, each with how it is started
// instead of by RunRole: they keep their own commands, and run as their own
// definitions rather than from their role files until S-0419.
var builtinRoles = map[string]string{
	workitem.ActivityPlanner:      "flai plan <item>",
	workitem.ActivityOrchestrator: "flai serve enable orchestrate",
	workitem.ActivityAnalyzer:     "flai analyze",
	workitem.ActivityCurator:      "flai curate",
}

// RunRole starts name, a role project e defines in its role file, on item
// when its file takes one, with inputs, the run's inputs in the order given,
// and returns its run as recorded, its trigger asked. It is refused, as a
// *Refused, for a built-in role's name, a name no role file gives, or one
// whose file does not load; while the host action the file names is off for
// the project; for an item when the file takes none, without one when it
// takes some, and for an item of a kind it does not take or that is
// archived, done, or cancelled; while a run of the role on the same item, or
// on none, runs; when its agent names no harness and no command is set on
// the host; and while the budget state of its agent's provider holds its
// start, naming the provider and until, with nothing recorded (S-0534). It
// does not wait for the run: the serving flai settles it once the process is
// gone.
func RunRole(ctx context.Context, o Options, e Entry, name, item string, inputs []harness.Input) (*AgentRun, error) {
	rs, err := roleCheck(o, e, name, item)
	if err != nil {
		return nil, err
	}
	l := newLauncher(o, e)
	l.handOver = true
	if err := l.budgetRefusal(rs.cfg, rs.agent); err != nil {
		return nil, err
	}
	if _, held := l.role(ctx, rs.cfg, rs.file, rs.item, rs.agent, inputs, TriggerAsked); held != nil {
		return nil, l.budgetRefused(*held)
	}
	run := o.Dir.AgentStates()[e.Root].Roles[rs.file.Name][rs.item]
	if run == nil {
		return nil, fmt.Errorf("the run of the role %s was not recorded in %s", rs.file.Name, o.Dir.agents())
	}
	if run.Error != "" {
		return run, errors.New(run.Why)
	}
	return run, nil
}

// startRole starts name, a role project e defines, on item, or on none,
// with trigger, what started it, through l, the project's launcher, whose
// lock the caller holds: the start an event or a schedule takes (S-0418).
// roleCheck's refusal is a *Refused, as for RunRole. A start the budget state
// of its provider holds is not made and nothing is recorded: held is that
// state, which the caller keeps the start for, to try again at a later look,
// rather than refuse it as RunRole does a start asked by hand (D-0421).
func startRole(ctx context.Context, o Options, e Entry, l *launcher, name, item, trigger string) (started bool, held *Budget, err error) {
	rs, err := roleCheck(o, e, name, item)
	if err != nil {
		return false, nil, err
	}
	started, held = l.role(ctx, rs.cfg, rs.file, rs.item, rs.agent, nil, trigger)
	return started, held, nil
}

// roleStart is what a defined role is started with once the checks pass:
// the operator's say, its file, its item's canonical ID or "", and its
// agent.
type roleStart struct {
	cfg   AgentConfig
	file  roles.Role
	item  string
	agent *manifest.Agent
}

// roleCheck makes RunRole's checks before name is started on item in
// project e, and returns what it is started with. A refusal is a *Refused;
// a project, a role file folder, or an item that cannot be read is an error.
func roleCheck(o Options, e Entry, name, item string) (roleStart, error) {
	repo, err := workitem.Open(e.Root)
	if err != nil {
		return roleStart{}, err
	}
	f, err := definedRole(repo, name)
	if err != nil {
		return roleStart{}, err
	}
	var cfg AgentConfig
	if o.Agent != nil {
		cfg = o.Agent(e.Root)
	}
	if !cfg.RoleActions[f.Action] {
		return roleStart{}, refused("the %s host action, which the role %s's file %s names, is off for this project: flai serve enable %s", f.Action, f.Name, f.Path, f.Action)
	}
	if item, err = roleItem(repo, f, item); err != nil {
		return roleStart{}, err
	}
	if run := o.Dir.AgentStates()[e.Root].Roles[f.Name][item]; run.running() {
		return roleStart{}, refused("the role %s is already running %s (pid %d, started %s); a role has one run %s at a time", f.Name, onItem(item), run.PID, run.Started, onItem(item))
	}
	agent := repo.Manifest.RoleAgent(f.Name)
	if (agent == nil || agent.Harness == "") && cfg.host(harness.Command).Program == "" {
		return roleStart{}, refused("the role %s's agent names no harness (roles.%s.agent or agent in system-flow.yaml), and no command is set on the host (flai serve agent set -- <program> [args...])", f.Name, f.Name)
	}
	return roleStart{cfg: cfg, file: f, item: item, agent: agent}, nil
}

// onItem says where a run is: on its item, or on no item.
func onItem(item string) string {
	if item == "" {
		return "on no item"
	}
	return "on " + item
}

// definedRole is the role file of repo's project that gives the role name.
// A built-in role's name, a name no file gives, and a file of that name that
// does not load are refused, saying why and what to do.
func definedRole(repo *workitem.Repo, name string) (roles.Role, error) {
	if how, builtin := builtinRoles[name]; builtin {
		return roles.Role{}, refused("%s is a built-in role, which flai starts with %s, not as a role the project defines", name, how)
	}
	design := repo.Manifest.Layout["design"]
	all, problems, err := roles.Load(repo.Root, design)
	if err != nil {
		return roles.Role{}, err
	}
	var names []string
	for _, f := range all {
		if f.Name == name {
			return f, nil
		}
		if _, builtin := builtinRoles[f.Name]; !builtin {
			names = append(names, f.Name)
		}
	}
	for _, p := range problems {
		if strings.TrimSuffix(path.Base(p.Path), ".md") == name {
			return roles.Role{}, refused("the role file %s does not load (%s: %s); fix it, as flai check reports it", p.Path, p.Code, p.Message)
		}
	}
	if len(names) == 0 {
		return roles.Role{}, refused("the project defines no role %s, nor any other: give a role its file under %s, as ADR-0207 says", name, roles.Dir(design))
	}
	return roles.Role{}, refused("the project defines no role %s; the roles it defines are %s", name, strings.Join(names, ", "))
}

// roleItem is item, as f's run takes it, by its canonical ID, or "" for a
// run on no item. A file that takes no item refuses one; a file that takes
// items refuses a run on none, since every run of such a role is on one of
// them (ADR-0207 §2, §3), and an item of a kind it does not take, or that is
// archived, done, or cancelled.
func roleItem(repo *workitem.Repo, f roles.Role, item string) (string, error) {
	switch {
	case item == "" && len(f.Item) == 0:
		return "", nil
	case item == "":
		return "", refused("the role %s runs on an item, a %s, as its file %s says; give it one", f.Name, strings.Join(f.Item, " or "), f.Path)
	case len(f.Item) == 0:
		return "", refused("the role %s takes no item, as its file %s says, and was given %s; run it on none", f.Name, f.Path, item)
	}
	id := workitem.CanonicalID(item)
	if kind := workitem.TypeOfID(id); !slices.Contains(f.Item, kind) {
		return "", refused("the role %s runs on a %s, as its file %s says, and %q is none", f.Name, strings.Join(f.Item, " or "), f.Path, item)
	}
	it, err := repo.Get(id)
	if err != nil {
		return "", err
	}
	switch {
	case it.Archived:
		return "", refused("%s is archived; the role %s runs on an item still on the board", it.ID, f.Name)
	case it.Status == workitem.Done || it.Status == workitem.Cancelled:
		return "", refused("%s is %s; the role %s runs on an item still open", it.ID, it.Status, f.Name)
	}
	return it.ID, nil
}

// role starts the role f on item with agent and inputs, records the run
// with trigger, what started it (asked, "<event> <ID>", or "schedule
// <spec>", joined by semicolons when coalesced), journals it under f's
// action, and says whether it started. A start
// the budget state of its provider holds is not made and nothing is
// recorded: held is that state, nil when it was not held (spawnGated,
// S-0534).
func (l *launcher) role(ctx context.Context, cfg AgentConfig, f roles.Role, item string, agent *manifest.Agent, inputs []harness.Input, trigger string) (started bool, held *Budget) {
	now := l.now().UTC()
	run := &AgentRun{Item: item, Agent: harness.RoleAgent(f.Name, item), Role: f.Role, RoleName: f.Name, Action: f.Action,
		Started: now.Format(time.RFC3339), Session: newSession(), Trigger: trigger}
	if agent != nil {
		run.Model = agent.Model
	}
	entry := hostapi.Entry{At: run.Started, Action: f.Action, Method: MethodRole, Project: l.entry.Key, Root: l.entry.Root, By: "flai serve"}
	fail := func(err error) (bool, *Budget) {
		run.Error, run.Ended, run.Outcome, run.Why = err.Error(), run.Started, OutcomeFailed, "could not be started: "+err.Error()
		l.told(run)
		entry.Outcome, entry.Detail = "failed", fmt.Sprintf("the role %s %s could not be started: %s", f.Name, onItem(item), run.Error)
		if l.record != nil {
			l.record(entry)
		}
		l.log("role could not be started", "role", f.Name, "item", item, "trigger", run.Trigger, "harness", run.Harness, "err", run.Error)
		appendAgentEvent(l.entry.Root, run, inbox.AgentFailed, run.Why, l.warn)
		return false, nil
	}
	name, adapter, err := harness.For(agent, cfg.host(harness.Command).Program != "")
	run.Harness = name
	if err != nil {
		return fail(err)
	}
	spec, keys, err := l.startOf(cfg, run, adapter, name, harness.Request{Role: f.Role, RoleFile: &f, Item: item, Inputs: inputs, Root: l.entry.Root, Project: l.entry.Key, Agent: agent,
		Name: run.Agent, Flai: cfg.Flai, Session: run.Session})
	if err != nil {
		return fail(err)
	}
	cmd, out, held, err := l.spawnGated(ctx, run, spec.Argv, f.Name, now, append(spec.Env, keys...))
	if held != nil {
		return false, held
	}
	if err != nil {
		return fail(err)
	}
	l.told(run)
	entry.Outcome, entry.Detail = "done", fmt.Sprintf("started %s (%s) to run the role %s %s as %s (pid %d), started by %s; log %s", run.Command, run.Harness, f.Name, onItem(item), run.Agent, run.PID, run.Trigger, run.Log)
	if l.record != nil {
		l.record(entry)
	}
	l.log("role started", "role", f.Name, "item", item, "trigger", run.Trigger, "harness", run.Harness, "command", run.Command, "pid", run.PID)
	appendAgentEvent(l.entry.Root, run, inbox.AgentStarted, triggerWhy(run), l.warn)
	l.await(cmd, out, run.PID, func(code int) { l.roleEnded(run, &code) })
	return true, nil
}

// roleEnded records a defined role's run as ended, with its exit code when
// it was seen, as judgeRole judges it, journals the end under the role's
// action, and logs its activity in the role's activity document,
// wip/agents/<name>.md, naming its item when it ran on one and what started
// it, which charges nothing. An activity that cannot be logged is warned of,
// and the run stays recorded as it ended. The end that recorded the run as
// ended first appends its agent event, as a planner run's does (S-0532). A
// run whose log ends on the harness's usage limit and that asked nothing is
// failed on it, with why naming the limit, whether its exit was seen or not,
// and marked so (limitEnded); every end feeds its provider's budget state
// once (budgetEnd, S-0534), as a planner run's does.
func (l *launcher) roleEnded(run *AgentRun, exit *int) {
	ended := *run
	ended.Ended, ended.Exit = l.now().UTC().Format(time.RFC3339), exit
	ended.Outcome, ended.Why, ended.Thread = judgeRole(l.entry.Root, run, exit)
	lim, _ := runLimit(run, l.warn)
	ended.limitEnded(lim)
	if lim != nil && ended.Outcome != OutcomeAsked {
		ended.Outcome, ended.Why = OutcomeFailed, limitWhy(lim)
	}
	first := l.ended(&ended)
	if first {
		appendAgentEnd(l.entry.Root, &ended, l.warn)
	}
	l.budgetEnd(&ended, lim, first)
	args := []any{"role", run.RoleName, "item", run.Item, "pid", run.PID, "outcome", ended.Outcome}
	if exit != nil {
		args = append(args, "exit", *exit)
	}
	l.log("role ended", args...)
	if l.record != nil {
		e := hostapi.Entry{At: ended.Ended, Action: run.Action, Method: MethodRole, Project: l.entry.Key, Root: l.entry.Root, By: "flai serve",
			Outcome: "done", Detail: fmt.Sprintf("the role %s %s, as %s (pid %d), ended, %s", run.RoleName, onItem(run.Item), run.Agent, run.PID, ended.Outcome)}
		switch ended.Outcome {
		case OutcomeFailed:
			e.Outcome, e.Detail = "failed", e.Detail+": "+ended.Why
		case OutcomeAsked:
			e.Detail += ": " + ended.Why
		}
		l.record(e)
	}
	var items []string
	if run.Item != "" {
		items = []string{run.Item}
	}
	if _, err := logRunEnd(l.dir, l.entry.Root, l.entry.Key, run.RoleName, run.Trigger, "", items); err != nil {
		l.warn("role activity not logged", "role", run.RoleName, "item", run.Item, "err", err)
	}
}

// judgeRole is how a defined role's run that has ended went, as a planner
// run's is (judgePlan): asked when a question of its own is open and
// unanswered, on its item when it ran on one and anywhere otherwise; failed
// on a failed exit; and worked otherwise, an exit nobody saw included.
func judgeRole(root string, run *AgentRun, exit *int) (outcome, why, thread string) {
	code := "an exit code nobody saw"
	if exit != nil {
		code = fmt.Sprintf("exit %d", *exit)
	}
	repo, err := workitem.Open(root)
	if err != nil {
		return OutcomeFailed, fmt.Sprintf("ended (%s); the project could not be read: %v", code, err), ""
	}
	var on []*threads.Thread
	if run.Item != "" {
		on, err = threads.For(repo, run.Item)
	} else {
		on, err = threads.List(repo)
	}
	if err != nil {
		return OutcomeFailed, fmt.Sprintf("ended (%s); the threads the role %s may have asked on could not be read: %v", code, run.RoleName, err), ""
	}
	for _, th := range on {
		if awaitsAnswer(th, run.Agent) {
			return OutcomeAsked, "waiting for an answer to " + th.ID + ": " + th.Title, th.ID
		}
	}
	if exit != nil && *exit != 0 {
		return OutcomeFailed, fmt.Sprintf("ended (%s) running the role %s %s", code, run.RoleName, onItem(run.Item)), ""
	}
	return OutcomeWorked, "", ""
}
