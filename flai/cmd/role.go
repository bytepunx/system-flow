package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"

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
the project defines, behind the host action its file names.`,
	}
	c.AddCommand(newRoleListCmd(a), newRoleShowCmd(a), newRoleRunCmd(a))
	return c
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
wip/agents/<name>.md, which charges nothing. It is no story's agent: the
in-progress limit does not count it, and it holds no story back.

It is behind the host action its file names, off until you enable it (flai
serve enable <action>). It refuses, and says why, for a built-in role
(flai plan, flai analyze, and flai curate start those) and a name no role
file gives, while the action is off for the project, for an item when its
file takes none, without one when it takes some, for an item of a kind it
does not take or that is archived, done, or cancelled, while a run of the
role on the same item, or on none, runs, and when its agent names no harness
and no command is set on the host. The host API's role.run runs this.`,
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
// serve is running to settle it, as flai plan does.
func (a *app) roleRun(name, item string, inputs []harness.Input) error {
	repo, err := a.project()
	if err != nil {
		return err
	}
	o, e := a.planOn(repo)
	run, err := serve.RunRole(context.Background(), o, e, name, item, inputs)
	var no *serve.Refused
	if errors.As(err, &no) {
		return fmt.Errorf("rule: %s", no.Why)
	}
	if err != nil {
		return err
	}
	if _, running := o.Dir.ReadStatus(a.now()); !running {
		fmt.Fprintln(a.errOut, "flai serve is not running (flai serve start, or flai host start): the role runs, and its end is recorded and its activity logged once flai serve is")
	}
	if a.jsonOut {
		return a.printJSON(map[string]any{"role": run.RoleName, "item": run.Item, "action": run.Action, "agent": run.Agent, "harness": run.Harness, "command": run.Command,
			"pid": run.PID, "log": run.Log, "session": run.Session, "started": run.Started, "trigger": run.Trigger})
	}
	on := "on no item"
	if run.Item != "" {
		on = "on " + run.Item
	}
	fmt.Fprintf(a.out, "started %s (%s) to run the role %s %s as %s (pid %d); log %s\n", run.Command, run.Harness, run.RoleName, on, run.Agent, run.PID, run.Log)
	return nil
}

// mcpRoleRun starts the role name the project defines on item, or on none,
// with inputs, for flai mcp's tool role_run (S-0417), as flai role run does,
// under the host action its file names, and journals who asked and what came
// of it under that action. No agent flai serve started may ask: a role is
// the operator's to start. A name the project defines no role by, a built-in
// role's among them, names no action, and its refusal is journalled nowhere,
// as role.run's is not.
func (a *app) mcpRoleRun(ctx context.Context, root, name, item string, inputs []harness.Input, by string) (mcpserver.RoleStarted, error) {
	if os.Getenv("FLAI_STARTED_BY") == "flai-serve" {
		return mcpserver.RoleStarted{}, fmt.Errorf("an agent flai serve started does not start a role: running the role %s is the operator's to ask for, from the dashboard or with flai role run %s on the host", name, name)
	}
	repo, err := workitem.Open(root)
	if err != nil {
		return mcpserver.RoleStarted{}, err
	}
	o, e := a.planOn(repo)
	run, err := serve.RunRole(ctx, o, e, name, item, inputs)
	asked := fmt.Sprintf("%s asked to run the role %s", by, name)
	if item != "" {
		asked += " on " + item
	}
	if action := definedRoleAction(e.Root, name); action != "" {
		entry := hostapi.Entry{At: a.now().UTC().Format(time.RFC3339), Action: action, Method: "mcp.role_run", Project: e.Key, Root: e.Root, By: by, Outcome: "done"}
		var no *serve.Refused
		switch {
		case errors.As(err, &no) && !o.Agent(e.Root).RoleActions[action]:
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
	return mcpserver.RoleStarted{Role: run.RoleName, Item: run.Item, Action: run.Action, Agent: run.Agent, Harness: run.Harness, Command: run.Command, PID: run.PID,
		Log: run.Log, Session: run.Session, Started: run.Started, Trigger: run.Trigger}, nil
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

// roleTrigger is a trigger entry as flai role prints it: asked, continuous,
// event <event>, or schedule <setting>, then (when <condition>).
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
