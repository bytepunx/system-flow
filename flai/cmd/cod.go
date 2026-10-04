package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/bytepunx/system-flow/flai/internal/planning"
)

// cod prints an epic's or story's cost of delay per week, from its inputs
// or its epic's share (S-0210); it writes nothing, so the planner records
// what it prints with flai edit.
func newCodCmd(a *app) *cobra.Command {
	return &cobra.Command{
		Use:   "cod <E-nnnn|S-nnnn>",
		Short: "Print an epic's or story's cost of delay per week, from its inputs or its epic's share",
		Long: `Works out what each week of waiting for an epic or a story costs, in
planning.currency. From the item's own cost of delay inputs, the value per week
is revenue_per_week plus penalty_per_week plus the hours of time_lost_per_cycle
times planning.hour_rate times the cycles in a week (a week over
planning.cycle, 168h by default); time lost needs planning.hour_rate. A story
without inputs of its own takes a share of its epic's value, worked out from
the epic's inputs or, without them, the value recorded on it: the share its
forecast duration is of the durations of the epic's open stories without
inputs, each its own forecast.duration or one worked out from history as flai
forecast does. An epic without inputs, or a story with neither inputs nor an
epic that has inputs or a value, is refused: the inputs are the operator's to
give, with flai edit --revenue-per-week, --penalty-per-week, and
--time-lost-per-cycle. Amounts are rounded to two decimals. flai cod writes
nothing; record the result with flai edit --cost-of-delay-value.`,
		Example: `  flai cod E-0016
  flai cod S-0210 --json`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			repo, err := a.project()
			if err != nil {
				return err
			}
			items, err := repo.List(true)
			if err != nil {
				return err
			}
			res, err := planning.CostOfDelay(items, repo.Manifest.Planning, args[0])
			if err != nil {
				return err
			}
			if a.jsonOut {
				return a.printJSON(res)
			}
			fmt.Fprintf(a.out, "%s cost of delay %.2f %s a week\n%s\n", res.ID, res.Value, res.Currency, res.Basis)
			return nil
		},
	}
}
