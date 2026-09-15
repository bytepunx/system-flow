package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newArchiveCmd(a *app) *cobra.Command {
	var dryRun bool
	c := &cobra.Command{
		Use:   "archive [id...]",
		Short: "Move done and cancelled items and their narratives to wip/archive",
		Long: `Without IDs, every closed epic, every closed story with its tasks, and
every closed task whose story is no longer on the board is archived. With
IDs, each must be done or cancelled; a story brings its tasks and narrative.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			repo, err := a.project()
			if err != nil {
				return err
			}
			items, err := repo.List(false)
			if err != nil {
				return err
			}
			plan, err := repo.PlanArchive(items, args)
			if err != nil {
				return err
			}
			if !dryRun {
				if err := repo.Archive(plan); err != nil {
					return err
				}
				if err := a.refreshIndex(repo); err != nil {
					return err
				}
			}
			if a.jsonOut {
				ids := make([]string, len(plan.Items))
				for i, it := range plan.Items {
					ids[i] = it.ID
				}
				return a.printJSON(map[string]any{"archived": ids, "narratives": plan.Narratives, "dry_run": dryRun})
			}
			if len(plan.Items) == 0 {
				fmt.Fprintln(a.out, "nothing to archive")
				return nil
			}
			verb := "archived"
			if dryRun {
				verb = "would archive"
			}
			for _, it := range plan.Items {
				fmt.Fprintf(a.out, "%s %s %s\n", verb, it.ID, it.Title)
			}
			for _, n := range plan.Narratives {
				fmt.Fprintf(a.out, "%s narrative %s\n", verb, relPath(repo.Root, n))
			}
			return nil
		},
	}
	c.Flags().BoolVar(&dryRun, "dry-run", false, "print the plan without moving anything")
	return c
}
