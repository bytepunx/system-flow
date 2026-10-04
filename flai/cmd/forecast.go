package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/bytepunx/system-flow/flai/internal/planning"
)

// forecast prints how long a story will take and when it will be delivered,
// worked out from done stories and the pull order (S-0210); it writes
// nothing, so the planner records what it prints with flai edit.
func newForecastCmd(a *app) *cobra.Command {
	return &cobra.Command{
		Use:   "forecast <S-nnnn>",
		Short: "Print a story's forecast duration and delivery, worked out from history and the pull order",
		Long: `Works out how long a story will take from the done stories with usage: the
median agent seconds per unit of size (acceptance criteria plus touches) over
the done stories that share its nature, model, and size band, falling back to
nature and band, then nature, then all, until at least three match; with fewer
than three on every rung it uses planning.default_duration. The delivery plays
out the pull order from now: the stories in progress hold their places in
progress, then ready and backlog stories start in order as the in-progress
limit leaves room, each after the stories in its after. A story ahead with a
forecast.duration is played out with it. flai forecast writes nothing; record
the result with flai edit --forecast-duration, --forecast-delivery, and
--forecast-basis.`,
		Example: `  flai forecast S-0210
  flai forecast S-0210 --json`,
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
			board, err := repo.LoadBoard()
			if err != nil {
				return fmt.Errorf("cannot read the board's pull order and in-progress limit: %w; restore board.md from git or run flai check to see what is wrong", err)
			}
			res, err := planning.Forecast(items, board, repo.Manifest.Planning, args[0], a.now())
			if err != nil {
				return err
			}
			if a.jsonOut {
				return a.printJSON(res)
			}
			fmt.Fprintf(a.out, "%s forecast %s, delivery %s\n%s\n", res.ID, res.Duration, res.Delivery, res.Basis)
			return nil
		},
	}
}
