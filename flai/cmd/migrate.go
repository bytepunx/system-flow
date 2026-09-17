package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func newMigrateCmd(a *app) *cobra.Command {
	c := &cobra.Command{
		Use:   "migrate",
		Short: "One-off migrations of a repository to the current standard",
	}
	c.AddCommand(newMigrateIDsCmd(a))
	return c
}

func newMigrateIDsCmd(a *app) *cobra.Command {
	var dryRun bool
	c := &cobra.Command{
		Use:   "ids",
		Short: "Widen work item IDs to four digits (E-001 becomes E-0001)",
		Long: `Renames every item and narrative in kanban and archive to a four-digit ID
and rewrites every reference in the design, docs, and wip folders, the root's
markdown and yaml files, and each project's root markdown files. Files move
with git mv inside a git repository. Run with --dry-run first.`,
		Example: `  flai migrate ids --dry-run
  flai migrate ids`,
		RunE: func(cmd *cobra.Command, args []string) error {
			repo, err := a.project()
			if err != nil {
				return err
			}
			plan, err := repo.PlanIDMigration()
			if err != nil {
				return err
			}
			if !dryRun && (len(plan.Renames) > 0 || len(plan.Rewrites) > 0) {
				if err := repo.ApplyIDMigration(plan, a.mover(repo.Root)); err != nil {
					return err
				}
				if err := a.refreshIndex(repo); err != nil {
					return err
				}
			}
			if a.jsonOut {
				return a.printJSON(map[string]any{"ids": plan.Map, "renames": plan.Renames, "rewrites": plan.Rewrites, "dry_run": dryRun})
			}
			if len(plan.Renames) == 0 && len(plan.Rewrites) == 0 {
				fmt.Fprintln(a.out, "nothing to migrate: every ID already has four digits")
				return nil
			}
			verb := "migrated"
			if dryRun {
				verb = "would migrate"
			}
			fmt.Fprintf(a.out, "%s %d item(s), %d file move(s), %d file(s) rewritten\n", verb, len(plan.Map), len(plan.Renames), len(plan.Rewrites))
			for _, mv := range plan.Renames {
				fmt.Fprintf(a.out, "  mv %s -> %s\n", mv.From, mv.To)
			}
			for _, f := range plan.Rewrites {
				fmt.Fprintf(a.out, "  rewrite %s\n", f)
			}
			return nil
		},
	}
	c.Flags().BoolVar(&dryRun, "dry-run", false, "print the plan without changing anything")
	return c
}

// mover returns git mv for tracked files when root is inside a git work
// tree, and os.Rename for untracked files and repositories without git.
func (a *app) mover(root string) func(from, to string) error {
	if _, err := a.runner.Run(root, "git", "rev-parse", "--is-inside-work-tree"); err != nil {
		return os.Rename
	}
	return func(from, to string) error {
		if _, err := a.runner.Run(root, "git", "ls-files", "--error-unmatch", from); err != nil {
			return os.Rename(from, to)
		}
		_, err := a.runner.Run(root, "git", "mv", from, to)
		return err
	}
}
