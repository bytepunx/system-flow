package cmd

import (
	"fmt"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

func newPromoteCmd(a *app) *cobra.Command {
	var candidates bool
	var limit int
	c := &cobra.Command{
		Use:   "promote --candidates",
		Short: "List the backlog stories that could go to ready",
		Long: `List the backlog stories that could go to ready, ordered by the project's
policy, and say why each other backlog story cannot.

A backlog story is a candidate when it is not a draft, meets the definition of
ready (a goal, acceptance criteria with a checkbox, and an epic that is not
cancelled), would not be held if it were ready (it declares touches that
overlap no story in progress or in review, and every story it names in after:
is done), and has a forecast duration and a cost of delay value.

The candidates are ordered by orchestration.policy in system-flow.yaml, fifo
when it is not set, as flai order --by orders the ready column, and each is
printed with the figure it was ordered by. Every other backlog story is
printed with each reason it is not a candidate.

--limit caps the candidates listed; those beyond it are not listed at all.
It writes nothing: moving a candidate to ready is flai move.`,
		Example: `  flai promote --candidates
  flai promote --candidates --limit 3 --json`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if !candidates {
				return fmt.Errorf("give --candidates: flai promote lists the backlog stories that could go to ready")
			}
			if limit < 0 {
				return fmt.Errorf("--limit is a number of candidates, 0 for all; got %d", limit)
			}
			repo, err := a.project()
			if err != nil {
				return err
			}
			got, err := repo.PromotionCandidates(limit)
			if err != nil {
				return err
			}
			if a.jsonOut {
				return a.printJSON(got)
			}
			if len(got.Candidates) == 0 && len(got.Others) == 0 {
				fmt.Fprintln(a.out, "no backlog stories")
				return nil
			}
			w := tabwriter.NewWriter(a.out, 0, 0, 2, ' ', 0)
			if len(got.Candidates) == 0 {
				fmt.Fprintf(w, "no candidates by %s\n", got.Policy)
			} else {
				fmt.Fprintf(w, "candidates by %s:\n", got.Policy)
			}
			for _, r := range got.Candidates {
				fmt.Fprintf(w, "  %d\t%s\t%s\t%s\n", r.Position, r.ID, r.Text, r.Title)
			}
			if err := w.Flush(); err != nil {
				return err
			}
			if len(got.Others) > 0 {
				fmt.Fprintln(a.out, "not candidates:")
			}
			for _, o := range got.Others {
				fmt.Fprintf(a.out, "  %s  %s\n", o.ID, o.Title)
				for _, reason := range o.Reasons {
					fmt.Fprintf(a.out, "    - %s\n", reason)
				}
			}
			return nil
		},
	}
	c.Flags().BoolVar(&candidates, "candidates", false, "list the backlog stories that could go to ready, and why each other one cannot")
	c.Flags().IntVar(&limit, "limit", 0, "list at most this many candidates, 0 for all; those beyond it are not listed")
	return c
}
