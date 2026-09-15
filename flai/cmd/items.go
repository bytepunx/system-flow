package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

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
	var tags []string
	parentFlag := map[string]string{workitem.Story: "epic", workitem.Task: "story"}[typ]
	c := &cobra.Command{
		Use:   "new \"<title>\"",
		Short: fmt.Sprintf("Create a %s from the item template", typ),
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			repo, err := a.project()
			if err != nil {
				return err
			}
			it, err := repo.Create(workitem.NewOptions{
				Type: typ, Title: args[0], Nature: nature, Parent: parent,
				Owner: orDefault(owner, a.author()), Tags: tags, Now: a.now(),
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
	if parentFlag != "" {
		c.Flags().StringVar(&parent, parentFlag, "", "parent "+parentFlag+" ID")
		_ = c.MarkFlagRequired(parentFlag)
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
