package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/bytepunx/system-flow/flai/internal/atomicfile"
	"github.com/bytepunx/system-flow/flai/internal/formulas"
	"github.com/bytepunx/system-flow/flai/internal/guard"
	"github.com/bytepunx/system-flow/flai/internal/harness"
	"github.com/bytepunx/system-flow/flai/internal/hostapi"
	"github.com/bytepunx/system-flow/flai/internal/mcpserver"
	"github.com/bytepunx/system-flow/flai/internal/roles"
	"github.com/bytepunx/system-flow/flai/internal/serve"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

func newRoleCmd(a *app) *cobra.Command {
	c := &cobra.Command{
		Use:   "role",
		Short: "Read the project's strategic roles: one file each under design/workflows/roles (ADR-0207)",
		Long: `A strategic role is one markdown file, <name>.md, in the design folder's
workflows/roles (ADR-0207). Its front matter names the role, its harness
definition, its host action, its triggers, the items it takes, its topics, how
its usage is charged, and the calls flai guard allows it, each plain or under
a named condition, with each call only the operator may allow inside a
permissions entry naming the operator's setting. Its body is its harness
definition's prompt. flai check reports a file that falls short as
roles.file, roles.grant, or roles.definition. flai role run starts a role
the project defines, behind the host action its file names, and flai serve
starts it on the events and the schedules its trigger names (S-0418).`,
	}
	c.AddCommand(newRoleListCmd(a), newRoleShowCmd(a), newRoleDefinitionCmd(a), newRoleRunCmd(a))
	return c
}

// flai role definition: a role's harness definition, generated from its
// file, printed or written into .flai-cache/agents (ADR-0236).
func newRoleDefinitionCmd(a *app) *cobra.Command {
	var write bool
	c := &cobra.Command{
		Use:   "definition <name> | --write [<name>...]",
		Short: "Print a role's harness definition generated from its role file, or write it into .flai-cache/agents",
		Long: `Print the harness definition flai generates from the role file of the role
named <name> (ADR-0207 §4, ADR-0236), as a Claude Code agent file: front
matter with name, description, tools on one comma-separated line, and model,
then the prompt, the file's body without the marker's line. The model is the
file's definition.model, inherit when it names none; flai serve puts the
role's agent's model over it at start. A built-in role (planner,
orchestrator, analyzer, curator) is generated from the project's role file,
or from the baseline flai carries when the project has none; a role the
project defines is generated from its file under the design folder's
workflows/roles. An unknown name is refused, naming the roles there are, and
so is a role file that lacks a description, tools, or a body.

With --write it writes the definition of each named role, or of every
built-in role when none is named, to .flai-cache/agents/<name>.md in the
checkout it runs in, a story's worktree included, which is git-ignored, and
prints each path it wrote. Run it again after a role file changes: the
generated file is never edited by hand. A session started by hand reads a
built-in role's definition through .claude/agents/<name>.md, a committed
symbolic link to ../../.flai-cache/agents/<name>.md. For each built-in role
written whose .claude/agents/<name>.md is not that link, it prints the
command that makes it, for you to run once in the checkout's root and commit;
flai writes nothing under .claude/. flai serve does not read the written
file: it generates the definition at every start.

--json gives the definition, its name, description, tools, model, and prompt;
with --write, the paths written and the link commands.`,
		Example: `  flai role definition planner
  flai role definition --write
  flai role definition --write reviewer`,
		RunE: func(_ *cobra.Command, args []string) error {
			if write {
				return a.roleDefinitionWrite(args)
			}
			if len(args) != 1 {
				return fmt.Errorf("give one role's name, or --write: usage: flai role definition <name> | --write [<name>...], which writes every built-in role's definition when no name is given")
			}
			return a.roleDefinitionPrint(args[0])
		},
	}
	c.Flags().BoolVar(&write, "write", false, "write each named role's definition, or every built-in role's, to .flai-cache/agents/<name>.md")
	return c
}

// agentsDir is where flai role definition --write puts the generated
// definitions, relative to the checkout's root (ADR-0236).
const agentsDir = ".flai-cache/agents"

// roleDefinitionPrint prints the definition generated for the role name.
func (a *app) roleDefinitionPrint(name string) error {
	repo, err := a.project()
	if err != nil {
		return err
	}
	g, err := a.roleDefinition(repo, name)
	if err != nil {
		return err
	}
	if a.jsonOut {
		return a.printJSON(map[string]any{"name": g.Name, "description": g.Description, "tools": g.Tools, "model": g.Model, "prompt": g.Prompt})
	}
	_, err = a.out.Write(g.File())
	return err
}

// roleDefinitionWrite writes the definitions of names, or of every built-in
// role when names is empty, into .flai-cache/agents, prints each path, and
// then the command that links each built-in one's .claude/agents/<name>.md
// that is not the link already.
func (a *app) roleDefinitionWrite(names []string) error {
	repo, err := a.project()
	if err != nil {
		return err
	}
	design := repo.Manifest.Layout["design"]
	builtin := map[string]bool{}
	var all []string
	for _, r := range roles.Baseline(design) {
		builtin[r.Name] = true
		all = append(all, r.Name)
	}
	if len(names) == 0 {
		names = all
	}
	// every name first, so that one that fails writes nothing
	defs := make([]roles.Generated, 0, len(names))
	for _, name := range names {
		g, err := a.roleDefinition(repo, name)
		if err != nil {
			return err
		}
		defs = append(defs, g)
	}
	dir := filepath.Join(repo.Root, filepath.FromSlash(agentsDir))
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("making %s: %w", agentsDir, err)
	}
	var written, links []string
	mkdir := false
	for _, g := range defs {
		path := filepath.Join(dir, g.Name+".md")
		if err := atomicfile.WriteFile(path, g.File(), 0o644); err != nil {
			return fmt.Errorf("writing %s/%s.md: %w", agentsDir, g.Name, err)
		}
		written = append(written, path)
		if !builtin[g.Name] {
			continue
		}
		link := "../../" + agentsDir + "/" + g.Name + ".md"
		at := filepath.Join(repo.Root, ".claude", "agents", g.Name+".md")
		if target, err := os.Readlink(at); err == nil && target == link {
			continue
		}
		if _, err := os.Stat(filepath.Dir(at)); errors.Is(err, os.ErrNotExist) {
			mkdir = true
		}
		links = append(links, fmt.Sprintf("ln -sf %s .claude/agents/%s.md", link, g.Name))
	}
	if mkdir {
		links = append([]string{"mkdir -p .claude/agents"}, links...)
	}
	if a.jsonOut {
		return a.printJSON(map[string]any{"written": written, "root": repo.Root, "links": orEmpty(links)})
	}
	for _, p := range written {
		fmt.Fprintf(a.out, "wrote %s\n", p)
	}
	if len(links) > 0 {
		fmt.Fprintf(a.out, "\nA session started by hand reads a built-in role through .claude/agents/<name>.md, a link to its generated definition, and flai writes nothing under .claude/: run this once in %s and commit the links:\n\n", repo.Root)
		for _, l := range links {
			fmt.Fprintf(a.out, "  %s\n", l)
		}
	}
	return nil
}

// orEmpty is list, or an empty list for none, so that JSON gives [].
func orEmpty(list []string) []string {
	if list == nil {
		return []string{}
	}
	return list
}

// roleDefinition is the definition flai generates for the role name in repo:
// a built-in role's from the project's file or the baseline flai carries, a
// defined role's from its file. An unknown name is refused, naming the roles
// there are.
func (a *app) roleDefinition(repo *workitem.Repo, name string) (roles.Generated, error) {
	design := repo.Manifest.Layout["design"]
	var known []string
	for _, r := range roles.Baseline(design) {
		known = append(known, r.Name)
		if r.Name != name {
			continue
		}
		r, fromBaseline, err := roles.Builtin(repo.Root, design, name)
		if err != nil {
			return roles.Generated{}, err
		}
		if _, err := os.Stat(filepath.Join(repo.Root, filepath.FromSlash(r.Path))); fromBaseline && err == nil {
			fmt.Fprintf(a.errOut, "%s does not parse, so flai runs the %s from the baseline it carries, which this is; flai check says why\n", r.Path, name)
		}
		return roles.Generate(r)
	}
	list, problems, err := roles.Load(repo.Root, design)
	if err != nil {
		return roles.Generated{}, err
	}
	for _, r := range list {
		if r.Name == name {
			return roles.Generate(r)
		}
		if !slices.Contains(known, r.Name) {
			known = append(known, r.Name)
		}
	}
	for _, p := range problems {
		fmt.Fprintf(a.errOut, "%s (left out)\n", p)
	}
	return roles.Generated{}, fmt.Errorf("no role %s; the roles are %s", name, strings.Join(known, ", "))
}

func newRoleListCmd(a *app) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List the project's roles with their action, triggers, items, topics, charge, calls, and permissions",
		Long: `List the roles in this checkout's design folder's workflows/roles, in file
order: each one's name, role verb, and file, and below it its action,
triggers, item kinds, topics, charge, the MCP tools, commands, and folders it
may use, and each permission's setting with the calls it allows. A call under
a condition is shown as call (when condition). A role file that does not load
is left out and said on stderr; flai check reports it. --json gives an array
of the roles, [] when none.`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			list, err := a.loadRoles()
			if err != nil {
				return err
			}
			if a.jsonOut {
				return a.printJSON(list)
			}
			if len(list) == 0 {
				fmt.Fprintln(a.out, "no roles")
				return nil
			}
			for i, r := range list {
				if i > 0 {
					fmt.Fprintln(a.out)
				}
				a.printRole(r)
			}
			return nil
		},
	}
}

func newRoleShowCmd(a *app) *cobra.Command {
	return &cobra.Command{
		Use:   "show <name>",
		Short: "Print a role: its front matter, its harness definition, and its body",
		Long: `Print one role by its name: what flai role list prints of it, then its
title, status, updated, works_story, its harness definition's description,
tools, and model, and the body, its definition's prompt. --json gives the
role.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			list, err := a.loadRoles()
			if err != nil {
				return err
			}
			var names []string
			for _, r := range list {
				if r.Name != args[0] {
					names = append(names, r.Name)
					continue
				}
				if a.jsonOut {
					return a.printJSON(r)
				}
				a.printRole(r)
				for _, f := range [][2]string{
					{"title", r.Title}, {"status", r.Status}, {"updated", r.Updated}, {"works_story", fmt.Sprint(r.WorksStory)},
					{"description", r.Definition.Description}, {"tools", digestList(r.Definition.Tools)}, {"model", r.Definition.Model},
				} {
					fmt.Fprintf(a.out, "  %-12s %s\n", f[0]+":", f[1])
				}
				if body := strings.TrimSpace(r.Body); body != "" {
					fmt.Fprintf(a.out, "\n%s\n", body)
				}
				return nil
			}
			if len(names) == 0 {
				return fmt.Errorf("no role %s: the project has no role that loads; flai role list and flai check say why", args[0])
			}
			return fmt.Errorf("no role %s; the roles are %s", args[0], strings.Join(names, ", "))
		},
	}
}

// flai role run: a role the project defines, on the operator's word
// (S-0417, ADR-0207).
func newRoleRunCmd(a *app) *cobra.Command {
	var inputs []string
	c := &cobra.Command{
		Use:   "run <name> [<item>]",
		Short: "Start a role the project defines, on an item its file takes or on none, behind the host action its file names",
		Long: `Starts a role the project defines in its role file under the design
folder's workflows/roles, now, on this host and as you (S-0417, ADR-0207).
It runs in the project's main checkout with the role's agent:
roles.<name>.agent in system-flow.yaml over the project's agent, started with
the harnesses and the command you set with flai serve agent, as the role
file's harness definition, its body the prompt. A role whose file's item
names epic or story runs on one such item, given as <item>; one whose file
takes none runs on none. Each --input, name=value, is one of the run's
inputs, in the order given, which its prompt names.

The run is recorded where flai serve tracks agents, under roles by the
role's name and item, and its output goes to a log beside flai serve's
state. Once it ends, the serving flai records how (worked, asked, or
failed), journals it under the role's action, and logs its activity in
wip/agents/<name>.md, its cost charged under usage.strategic by the rule
its file names in charge (ADR-0249): item, to <item> and the items above
it; named, to the work items the entry names; issues, to the issues it
names; none, to nothing. It is no story's agent: the in-progress limit does
not count it, and it holds no story back.

It is behind the host action its file names, off until you enable it (flai
serve enable <action>). It refuses, and says why, for a built-in role
(flai plan, flai analyze, and flai curate start those) and a name no role
file gives, while the action is off for the project, for an item when its
file takes none, without one when it takes some, for an item of a kind it
does not take or that is archived, done, or cancelled, while a run of the
role on the same item, or on none, runs, and when its agent names no harness
and no command is set on the host. The host API's role.run runs this.

This is the role's asked trigger, which its run records as asked. While
the action is on, flai serve also starts the role on the other triggers its
file names (S-0418): on an event, {event: <event>, when: <condition>}, once
orchestration.coalesce has passed since the first, with every event seen
meanwhile, one run per item for a role that takes items; and on a schedule,
{schedule: roles.<name>.schedule}, a cron expression in UTC or daily in
system-flow.yaml, on no item. The run records what started it: asked,
<event> <ID>, or schedule <spec>, joined by semicolons when coalesced.

An agent flai serve started is refused, as running a role is the
operator's to ask for, save one (S-0422, ADR-0239): the agent of a
formula's instance, for a role a role step of its formula names, on the
instance as <item>. Its story is FLAI_STORY, else the S-nnnn of its name
agent-S-nnnn. Its run records the trigger formula <name> <S-nnnn>, and its
end is an agent event on the instance, which that agent's wait_for_events
reports.`,
		Example: `  flai role list
  flai serve enable review
  flai role run reviewer S-0042
  flai role run auditor --input focus=docs --json`,
		Args: cobra.RangeArgs(1, 2),
		RunE: func(_ *cobra.Command, args []string) error {
			given, err := roleInputs(inputs)
			if err != nil {
				return err
			}
			item := ""
			if len(args) == 2 {
				item = args[1]
			}
			return a.roleRun(args[0], item, given)
		},
	}
	c.Flags().StringArrayVar(&inputs, "input", nil, "one of the run's inputs, as name=value, which its prompt names; give it once for each, in order")
	return c
}

// roleInputs are the run's inputs as --input gives them, name=value each,
// in order; one with no = or no name is refused.
func roleInputs(given []string) ([]harness.Input, error) {
	out := make([]harness.Input, 0, len(given))
	for _, g := range given {
		name, value, ok := strings.Cut(g, "=")
		if name = strings.TrimSpace(name); !ok || name == "" {
			return nil, fmt.Errorf("--input %q is not name=value: give each input as its name, =, and its value, such as --input focus=docs", g)
		}
		out = append(out, harness.Input{Name: name, Value: value})
	}
	return out, nil
}

// roleRun starts the role name the project defines on item, or on none,
// with inputs, prints the run, and says on standard error when no flai
// serve is running to settle it, as flai plan does. An agent flai serve
// started is refused, save a formula instance's agent for a role step of its
// formula (runTrigger).
func (a *app) roleRun(name, item string, inputs []harness.Input) error {
	repo, err := a.project()
	if err != nil {
		return err
	}
	trigger, err := runTrigger(repo.Root, name, item, os.Getenv("FLAI_AGENT"))
	var no *serve.Refused
	if errors.As(err, &no) {
		return fmt.Errorf("rule: %s", no.Why)
	}
	if err != nil {
		return err
	}
	o, e := a.planOn(repo)
	run, err := serve.RunRole(context.Background(), o, e, name, item, inputs, trigger)
	if errors.As(err, &no) {
		return fmt.Errorf("rule: %s", no.Why)
	}
	if err != nil {
		return err
	}
	if _, running := o.Dir.ReadStatus(a.now()); !running {
		fmt.Fprintln(a.errOut, "flai serve is not running (flai serve start, or flai host start): the role runs, and its end is recorded and its activity logged once flai serve is")
	}
	role, on, action := roleStartedAs(name, run)
	if a.jsonOut {
		return a.printJSON(map[string]any{"role": role, "item": on, "action": action, "agent": run.Agent, "harness": run.Harness, "command": run.Command,
			"pid": run.PID, "log": run.Log, "session": run.Session, "started": run.Started, "trigger": run.Trigger})
	}
	fmt.Fprintf(a.out, "started %s (%s) to run the role %s %s as %s (pid %d); log %s\n", run.Command, run.Harness, role, onItem(on), run.Agent, run.PID, run.Log)
	return nil
}

// roleStartedAs is the role, the item, and the host action of run, which
// serve.RunRole started when asked for the role name: a defined role's, as
// its run records them, or, for the analyzer a formula instance's role step
// started (S-0424), name, the instance it ran for, and the analyze action.
func roleStartedAs(name string, run *serve.AgentRun) (role, item, action string) {
	if run.RoleName != "" {
		return run.RoleName, run.Item, run.Action
	}
	return name, run.Instance, hostapi.ActionAnalyze
}

// mcpRoleRun starts the role name the project defines on item, or on none,
// with inputs, for flai mcp's tool role_run (S-0417), as flai role run does,
// under the host action its file names, and journals who asked and what came
// of it under that action. No agent flai serve started may ask, as a role is
// the operator's to start, save the agent of a formula's instance for a role
// step of its formula, on the instance (runTrigger, S-0422); a refusal of
// any other is journalled nowhere. A name the project defines no role by, a
// built-in role's among them, names no action, and its refusal is journalled
// nowhere, as role.run's is not, save the analyzer a formula instance's role
// step asks for, which runs and is journalled under the analyze action
// (S-0424).
func (a *app) mcpRoleRun(ctx context.Context, root, name, item string, inputs []harness.Input, by string) (mcpserver.RoleStarted, error) {
	trigger, err := runTrigger(root, name, item, by)
	if err != nil {
		return mcpserver.RoleStarted{}, err
	}
	repo, err := workitem.Open(root)
	if err != nil {
		return mcpserver.RoleStarted{}, err
	}
	o, e := a.planOn(repo)
	run, err := serve.RunRole(ctx, o, e, name, item, inputs, trigger)
	asked := fmt.Sprintf("%s asked to run the role %s", by, name)
	if item != "" {
		asked += " on " + item
	}
	action := definedRoleAction(e.Root, name)
	if action == "" && name == workitem.ActivityAnalyzer && trigger != serve.TriggerAsked {
		action = hostapi.ActionAnalyze // a formula's role step starts the analyzer under its own action
	}
	if action != "" {
		entry := hostapi.Entry{At: a.now().UTC().Format(time.RFC3339), Action: action, Method: "mcp.role_run", Project: e.Key, Root: e.Root, By: by, Outcome: "done"}
		cfg := o.Agent(e.Root)
		on := cfg.RoleActions[action]
		if action == hostapi.ActionAnalyze {
			on = cfg.Analyze
		}
		var no *serve.Refused
		switch {
		case errors.As(err, &no) && !on:
			entry.Outcome, entry.Detail = "disabled", fmt.Sprintf("%s: %s", asked, no.Why)
		case err != nil:
			entry.Outcome, entry.Detail = "failed", fmt.Sprintf("%s: %s", asked, err)
		default:
			entry.Detail = fmt.Sprintf("%s: started %s as %s (pid %d)", asked, run.Command, run.Agent, run.PID)
		}
		o.Host.Record(entry)
	}
	if err != nil {
		return mcpserver.RoleStarted{}, err
	}
	role, on, started := roleStartedAs(name, run)
	return mcpserver.RoleStarted{Role: role, Item: on, Action: started, Agent: run.Agent, Harness: run.Harness, Command: run.Command, PID: run.PID,
		Log: run.Log, Session: run.Session, Started: run.Started, Trigger: run.Trigger}, nil
}

// runTrigger is the trigger of the run of the role name on item that the
// session asks for, by being its agent's name: asked, for a session flai
// serve did not start (the operator's own), and, for one it did, formula
// <name> <S-nnnn> when the session is the agent of a formula's instance
// (S-0421), item is that instance, and a role step of its formula names the
// role (S-0422). The session's story is FLAI_STORY, else the S-nnnn of by
// when by is agent-S-nnnn. Any other session flai serve started is refused,
// as a *serve.Refused that says why; a story or a formula that cannot be read
// is an error. A role step that names the built-in analyzer gets its trigger
// too: serve.RunRole starts the analyzer for it through the analyzer's own
// start, and refuses the other built-in roles, as it refuses every built-in
// role to the operator's asking, naming its command (S-0424).
func runTrigger(root, name, item, by string) (string, error) {
	if os.Getenv(guard.StartedByEnv) != guard.StartedByServe {
		return serve.TriggerAsked, nil
	}
	operators := fmt.Sprintf("running the role %s is the operator's to ask for, from the dashboard or with flai role run %s on the host", name, name)
	story := recordingStory("", func(key string) string {
		if key == "FLAI_AGENT" {
			return by
		}
		return os.Getenv(key)
	}, func() string { return "" })
	if story == "" {
		return "", &serve.Refused{Why: "an agent flai serve started does not start a role, save the agent of a formula's instance for a role step of its formula, and this one works no story (no FLAI_STORY, and no agent name of the form agent-S-nnnn): " + operators}
	}
	repo, err := workitem.Open(root)
	if err != nil {
		return "", err
	}
	it, err := repo.Get(story)
	if err != nil {
		return "", fmt.Errorf("reading %s, the story of the agent that asked to run the role %s: %w", story, name, err)
	}
	if it.Formula == "" {
		return "", &serve.Refused{Why: fmt.Sprintf("an agent flai serve started does not start a role, save the agent of a formula's instance for a role step of its formula, and %s is no formula's instance (it names no formula): %s", it.ID, operators)}
	}
	if workitem.CanonicalID(item) != it.ID {
		return "", &serve.Refused{Why: fmt.Sprintf("the agent of %s, an instance of the formula %s, starts a role step's role on %s alone, and asked to run the role %s %s: give the item %s", it.ID, it.Formula, it.ID, name, onItem(item), it.ID)}
	}
	list, _, err := formulas.Load(repo.Root, repo.Manifest.Layout["design"])
	if err != nil {
		return "", err
	}
	f, err := formulas.Find(list, it.Formula)
	if err != nil {
		return "", fmt.Errorf("%s is an instance of the formula %s, whose file does not load: %w; fix it, as flai check reports it", it.ID, it.Formula, err)
	}
	var named []string
	for _, s := range f.Steps {
		if s.Kind() != formulas.KindRole {
			continue
		}
		if s.Run.Role == name {
			return serve.FormulaTrigger(f.Name, it.ID), nil
		}
		if !slices.Contains(named, s.Run.Role) {
			named = append(named, s.Run.Role)
		}
	}
	if len(named) == 0 {
		return "", &serve.Refused{Why: fmt.Sprintf("the formula %s, which %s is an instance of, has no role step, so its agent starts no role, and not %s: %s", f.Name, it.ID, name, operators)}
	}
	return "", &serve.Refused{Why: fmt.Sprintf("the formula %s, which %s is an instance of, names the role %s in its steps and not %s: its agent starts those alone, and %s", f.Name, it.ID, strings.Join(named, ", "), name, operators)}
}

// onItem says where a role's run is asked for: on its item, or on no item.
func onItem(item string) string {
	if item == "" {
		return "on no item"
	}
	return "on " + item
}

// definedRoleAction is the host action the file of name, a role the project
// at root defines, names, or "" for a built-in role's name, a name no file
// that loads gives, and a project whose role files cannot be read.
func definedRoleAction(root, name string) string {
	for _, builtin := range roles.BuiltinActions {
		if builtin == name {
			return ""
		}
	}
	repo, err := workitem.Open(root)
	if err != nil {
		return ""
	}
	list, _, err := roles.Load(repo.Root, repo.Manifest.Layout["design"])
	if err != nil {
		return ""
	}
	for _, r := range list {
		if r.Name == name {
			return r.Action
		}
	}
	return ""
}

// loadRoles reads the project's roles and says each role file that does not
// load on stderr.
func (a *app) loadRoles() ([]roles.Role, error) {
	repo, err := a.project()
	if err != nil {
		return nil, err
	}
	list, problems, err := roles.Load(repo.Root, repo.Manifest.Layout["design"])
	if err != nil {
		return nil, err
	}
	for _, p := range problems {
		fmt.Fprintf(a.errOut, "%s (left out)\n", p)
	}
	if list == nil {
		list = []roles.Role{}
	}
	return list, nil
}

// printRole prints what flai role list says of a role.
func (a *app) printRole(r roles.Role) {
	fmt.Fprintf(a.out, "%s  %s  %s\n", r.Name, r.Role, r.Path)
	var triggers []string
	for _, t := range r.Trigger {
		triggers = append(triggers, roleTrigger(t))
	}
	item := "none"
	if len(r.Item) > 0 {
		item = strings.Join(r.Item, ", ")
	}
	for _, f := range [][2]string{
		{"action", r.Action}, {"trigger", digestList(triggers)}, {"item", item}, {"topics", digestList(r.Topics)}, {"charge", r.Charge},
		{"tools", roleCalls(r.Tools)}, {"commands", roleCalls(r.Commands)}, {"files", digestList(r.Files)},
	} {
		fmt.Fprintf(a.out, "  %-12s %s\n", f[0]+":", f[1])
	}
	if len(r.Permissions) == 0 {
		fmt.Fprintf(a.out, "  %-12s -\n", "permissions:")
	} else {
		fmt.Fprintf(a.out, "  %s\n", "permissions:")
	}
	for _, p := range r.Permissions {
		fmt.Fprintf(a.out, "    %s → %s\n", p.Setting, roleCalls(append(append([]roles.Call{}, p.Tools...), p.Commands...)))
	}
}

// roleTrigger is a trigger entry as flai role prints it: asked, event
// <event>, or schedule <setting>, then (when <condition>). Another word, such
// as continuous, which flai check reports (ADR-0209 §1), is printed as it is.
func roleTrigger(t roles.Trigger) string {
	s := t.Kind
	switch t.Kind {
	case roles.TriggerEvent:
		s += " " + t.Event
	case roles.TriggerSchedule:
		s += " " + t.Schedule
	}
	if t.When != "" {
		s += " (when " + t.When + ")"
	}
	return s
}

// roleCalls are call entries as flai role prints them: call, or call (when
// condition), comma separated, - for none, as digestList prints a list.
func roleCalls(calls []roles.Call) string {
	out := make([]string, len(calls))
	for i, c := range calls {
		out[i] = c.Call
		if c.When != "" {
			out[i] += " (when " + c.When + ")"
		}
	}
	return digestList(out)
}
