package cmd

import (
	"fmt"
	"io"
	"strings"

	"github.com/spf13/cobra"

	"github.com/bytepunx/system-flow/flai/internal/docedit"
	"github.com/bytepunx/system-flow/flai/internal/itemnew"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

func newItemCmd(a *app, typ string) *cobra.Command {
	c := &cobra.Command{
		Use:   typ,
		Short: fmt.Sprintf("Create and inspect %ss", typ),
	}
	c.AddCommand(newItemNewCmd(a, typ))
	return c
}

func newItemNewCmd(a *app, typ string) *cobra.Command {
	var nature, owner, parent string
	var tags, touches, trailers []string
	var bodyStdin, autocommit, printBody bool
	parentFlag := map[string]string{workitem.Story: "epic", workitem.Task: "story"}[typ]
	c := &cobra.Command{
		Use:   "new \"<title>\"",
		Short: fmt.Sprintf("Create a %s from the item template", typ),
		Long: fmt.Sprintf(`Create a %s from the project's item template with the next free ID,
linked into its parent.

With --body-stdin the body below the item's heading is read from standard
input instead of the template's empty sections, and the creation is one step
that happens or does not: flai check runs with the new item in place, and if
it reports anything the item introduces, the item is removed, its parent is
restored, and the findings are printed (exit 4). --autocommit commits the new
item and its parent on their own, unless the project sets
dashboard.autocommit: false. Nothing is pushed. --print-body prints the body
the template gives, for a form or a script to start from, and creates nothing.`, typ),
		Args: func(cmd *cobra.Command, args []string) error {
			if printBody {
				return cobra.NoArgs(cmd, args)
			}
			return cobra.ExactArgs(1)(cmd, args)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			repo, err := a.project()
			if err != nil {
				return err
			}
			if printBody {
				body, err := repo.TemplateBody(typ)
				if err != nil {
					return err
				}
				if a.jsonOut {
					return a.printJSON(map[string]string{"type": typ, "body": body})
				}
				fmt.Fprint(a.out, body)
				return nil
			}
			opt := workitem.NewOptions{
				Type: typ, Title: args[0], Nature: nature, Parent: parent,
				Owner: orDefault(owner, a.author()), Tags: tags, Touches: touches, Now: a.now(),
			}
			if bodyStdin || autocommit {
				if bodyStdin {
					data, err := io.ReadAll(cmd.InOrStdin())
					if err != nil {
						return err
					}
					if strings.TrimSpace(string(data)) == "" {
						return fmt.Errorf("--body-stdin was given and standard input is empty; write the body, or leave the flag out for the template's empty sections")
					}
					opt.Body = string(data)
				}
				res, err := itemnew.Create(repo, a.runner, itemnew.Options{New: opt, Autocommit: autocommit, Trailers: trailers})
				if r, ok := docedit.IsRefused(err); ok {
					if a.jsonOut {
						_ = a.printJSON(map[string]any{"refused": r})
					} else {
						for _, f := range r.Findings {
							fmt.Fprintf(a.out, "%s:%d: %s: %s: %s\n", f.Path, f.Line, f.Level, f.Rule, f.Message)
						}
					}
					return &exitError{code: exitDocRefused, msg: r.Error()}
				}
				if err != nil {
					return err
				}
				if a.jsonOut {
					return a.printJSON(res)
				}
				fmt.Fprintf(a.out, "%s %s\n  %s\n", res.Item.ID, res.Item.Title, res.Path)
				switch {
				case res.Committed:
					fmt.Fprintf(a.out, "  committed %s\n", res.Commit)
				case res.CommitError != "":
					fmt.Fprintf(a.out, "  NOT committed (%s)\n", firstLine(res.CommitError))
				}
				return nil
			}
			it, err := repo.Create(workitem.NewOptions{
				Type: typ, Title: args[0], Nature: nature, Parent: parent,
				Owner: orDefault(owner, a.author()), Tags: tags, Touches: touches, Now: a.now(),
			})
			if err != nil {
				return err
			}
			if a.jsonOut {
				return a.printJSON(it)
			}
			fmt.Fprintf(a.out, "%s %s\n  %s\n", it.ID, it.Title, relPath(repo.Root, it.Path))
			return nil
		},
	}
	c.Flags().StringVar(&nature, "nature", "feature", "one of "+strings.Join(workitem.Natures, ", "))
	c.Flags().StringVar(&owner, "owner", "", "owner (default: config author)")
	c.Flags().StringSliceVar(&tags, "tag", nil, "tag (repeatable or comma separated)")
	c.Flags().StringSliceVar(&touches, "touches", nil, "paths or components this work changes (repeatable or comma separated)")
	c.Flags().BoolVar(&bodyStdin, "body-stdin", false, "read the body below the heading from standard input; checked before it is kept")
	c.Flags().BoolVar(&autocommit, "autocommit", false, "commit the new item and its parent on their own, unless dashboard.autocommit is false")
	c.Flags().StringArrayVar(&trailers, "trailer", nil, "trailer line for the commit (repeatable)")
	c.Flags().BoolVar(&printBody, "print-body", false, "print the body the template gives this type and create nothing")
	if parentFlag != "" {
		help := "parent " + parentFlag + " ID"
		if typ == workitem.Story {
			help += " (optional: a story need not belong to one, S-0092)"
		}
		c.Flags().StringVar(&parent, parentFlag, "", help)
		if typ != workitem.Story {
			c.PreRunE = func(cmd *cobra.Command, args []string) error {
				if !printBody && parent == "" {
					return fmt.Errorf("required flag \"%s\" not set", parentFlag)
				}
				return nil
			}
		}
	}
	return c
}

func newShowCmd(a *app) *cobra.Command {
	return &cobra.Command{
		Use:   "show <id>",
		Short: "Print one work item with its children and history",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			repo, err := a.project()
			if err != nil {
				return err
			}
			it, err := repo.Get(args[0])
			if err != nil {
				return err
			}
			items, err := repo.List(true)
			if err != nil {
				return err
			}
			children := workitem.Children(items, it.ID)
			if a.jsonOut {
				return a.printJSON(map[string]any{"item": it, "children": children})
			}
			fmt.Fprintf(a.out, "%s %s\n", it.ID, it.Title)
			fmt.Fprintf(a.out, "  %s · %s · %s", it.Type, it.Nature, it.Status)
			if it.IsBlocked() {
				fmt.Fprint(a.out, " · BLOCKED")
			}
			if it.Archived {
				fmt.Fprint(a.out, " · archived")
			}
			fmt.Fprintln(a.out)
			if it.Parent != "" {
				fmt.Fprintf(a.out, "  parent: %s\n", it.Parent)
			}
			fmt.Fprintf(a.out, "  owner: %s · created %s · updated %s\n", it.Owner, it.Created, it.Updated)
			if len(it.Tags) > 0 {
				fmt.Fprintf(a.out, "  tags: %s\n", strings.Join(it.Tags, ", "))
			}
			fmt.Fprintf(a.out, "  file: %s\n", relPath(repo.Root, it.Path))
			if len(it.Transitions) > 0 {
				fmt.Fprintln(a.out, "  history:")
				for _, tr := range it.Transitions {
					fmt.Fprintf(a.out, "    %s  %-12s by %s\n", tr.At, tr.To, tr.By)
				}
			}
			for _, b := range it.Blocked {
				until := orDefault(b.Until, "open")
				fmt.Fprintf(a.out, "  blocked: %s to %s: %s\n", b.From, until, b.Reason)
			}
			if len(children) > 0 {
				fmt.Fprintln(a.out, "  children:")
				for _, c := range children {
					fmt.Fprintf(a.out, "    %s  %-12s %s\n", c.ID, c.Status, c.Title)
				}
			}
			return nil
		},
	}
}
