package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/bytepunx/system-flow/flai/internal/roles"
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
roles.file, roles.grant, or roles.definition.`,
	}
	c.AddCommand(newRoleListCmd(a), newRoleShowCmd(a))
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
