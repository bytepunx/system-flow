package cmd

import (
	"fmt"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

func newPromoteCmd(a *app) *cobra.Command {
	var candidates, drafts bool
	var limit int
	c := &cobra.Command{
		Use:   "promote --candidates | --drafts",
		Short: "List the backlog stories that could go to ready, or the drafts and what each lacks",
		Long: `List the backlog stories that could go to ready, ordered by the project's
policy, and say why each other backlog story cannot; or, with --drafts, list
each draft story in the backlog and what it lacks to be finalized.

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

--drafts lists each draft story in the backlog, in the same policy order, as
complete or with each thing it lacks. A draft is complete when it has every
section of the project's story template, a goal, acceptance criteria with a
checkbox, at least one touch, a forecast duration and delivery, a cost of
delay value, and an epic that is open or none. Whether its criteria,
touches, and forecast describe the same work is for whoever finalizes it.

It writes nothing: moving a candidate to ready is flai move, and finalizing
a draft is flai edit --no-draft.`,
		Example: `  flai promote --candidates
  flai promote --candidates --limit 3 --json
  flai promote --drafts --json`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if candidates == drafts {
				return fmt.Errorf("give --candidates or --drafts: flai promote lists the backlog stories that could go to ready, or the drafts and what each lacks")
			}
			if drafts {
				if cmd.Flags().Changed("limit") {
					return fmt.Errorf("--limit caps the candidates; --drafts lists every draft")
				}
				return a.printDrafts()
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
	c.Flags().BoolVar(&drafts, "drafts", false, "list each draft story in the backlog, complete or with what it lacks to be finalized")
	c.Flags().IntVar(&limit, "limit", 0, "list at most this many candidates, 0 for all; those beyond it are not listed")
	return c
}

// printDrafts lists each draft story in the backlog by the project's policy,
// complete or with each thing it lacks.
func (a *app) printDrafts() error {
	repo, err := a.project()
	if err != nil {
		return err
	}
	got, err := repo.Drafts()
	if err != nil {
		return err
	}
	if a.jsonOut {
		return a.printJSON(got)
	}
	if len(got.Drafts) == 0 {
		fmt.Fprintln(a.out, "no draft stories in the backlog")
		return nil
	}
	fmt.Fprintf(a.out, "drafts by %s:\n", got.Policy)
	for _, d := range got.Drafts {
		state := "complete"
		if !d.Complete {
			state = "incomplete"
		}
		fmt.Fprintf(a.out, "  %d  %s  %s  %s\n", d.Position, d.ID, state, d.Title)
		for _, lack := range d.Lacks {
			fmt.Fprintf(a.out, "    - %s\n", lack)
		}
	}
	return nil
}
