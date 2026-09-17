package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// touches sets the advisory list of paths or components a story or task is
// working on (ADR-0019); flai check warns when in-progress items overlap.
func newTouchesCmd(a *app) *cobra.Command {
	var clear bool
	c := &cobra.Command{
		Use:   "touches <id> [path-or-component...]",
		Short: "Set what a story or task is working on; flai check warns on overlap",
		Example: `  flai touches S-0037 flai/internal/workitem flaiover/src/routes/docs
  flai touches T-0121 --clear
  flai touches S-0037            # show`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			repo, err := a.project()
			if err != nil {
				return err
			}
			it, err := repo.Get(args[0])
			if err != nil {
				return err
			}
			if it.Type == workitem.Epic {
				return fmt.Errorf("%s is an epic; touches belong to stories and tasks", it.ID)
			}
			if clear || len(args) > 1 {
				it.Touches = nil
				for _, p := range args[1:] {
					p = strings.TrimSuffix(strings.TrimSpace(p), "/")
					if p != "" {
						it.Touches = append(it.Touches, p)
					}
				}
				it.Updated = a.now().UTC().Format(workitem.TimeFormat)
				if err := repo.Save(it); err != nil {
					return err
				}
			}
			if a.jsonOut {
				return a.printJSON(map[string]any{"id": it.ID, "touches": it.Touches})
			}
			if len(it.Touches) == 0 {
				fmt.Fprintf(a.out, "%s touches nothing\n", it.ID)
				return nil
			}
			fmt.Fprintf(a.out, "%s touches %s\n", it.ID, strings.Join(it.Touches, ", "))
			return nil
		},
	}
	c.Flags().BoolVar(&clear, "clear", false, "remove the list")
	return c
}
