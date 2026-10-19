package cmd

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/bytepunx/system-flow/flai/internal/mergequeue"
)

func newQueueCmd(a *app) *cobra.Command {
	return &cobra.Command{
		Use:   "queue",
		Short: "List the stories in review the merge queue would batch, and why each other would not",
		Long: `List the stories in review that the merge queue would batch, in the order
they last entered review: the last flai verify passed at the head of the
story's branch, every acceptance criterion is ticked, no thread on the story
or its tasks is open, and no path is recorded as undoing main's change
(ADR-0155). Each other story in review follows, with every reason it is not
eligible. Then the batch under verification and its step, and each story's
last outcome, from the queue's state in .flai-cache/merge-queue.

Read-only: it verifies, merges, and moves nothing (ADR-0173). Under
acceptance.mode pull-request the queue does not run, because the forge's
required checks gate each story (ADR-0169), and no story is listed.`,
		Example: `  flai queue
  flai queue --json`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			repo, err := a.project()
			if err != nil {
				return err
			}
			q, err := mergequeue.Compute(a.runner, repo)
			if err != nil {
				return err
			}
			if a.jsonOut {
				return a.printJSON(q)
			}
			printQueue(a.out, q)
			return nil
		},
	}
}

// printQueue writes q as text: the eligible stories in order, each other
// story in review with its reasons, the batch, and the outcomes; or, when
// the queue is off, the one line saying why.
func printQueue(w io.Writer, q mergequeue.Queue) {
	if q.Off {
		fmt.Fprintf(w, "merge queue off under acceptance.mode %s: %s\n", q.Mode, q.OffReason)
		return
	}
	if len(q.Eligible) == 0 && len(q.Ineligible) == 0 {
		fmt.Fprintln(w, "no story is in review")
	} else {
		fmt.Fprintf(w, "eligible (%d)\n", len(q.Eligible))
		if len(q.Eligible) == 0 {
			fmt.Fprintln(w, "  -")
		}
		for i, e := range q.Eligible {
			fmt.Fprintf(w, "  %d. %s\n", i+1, queueLine(e))
		}
		fmt.Fprintf(w, "not eligible (%d)\n", len(q.Ineligible))
		if len(q.Ineligible) == 0 {
			fmt.Fprintln(w, "  -")
		}
		for _, e := range q.Ineligible {
			fmt.Fprintf(w, "  %s\n", queueLine(e))
			for _, r := range e.Reasons {
				fmt.Fprintf(w, "       %s: %s\n", r.Kind, r.Text)
			}
		}
	}
	if q.Batch == nil {
		fmt.Fprintln(w, "batch: none")
	} else {
		fmt.Fprintf(w, "batch: %s, step %s\n", strings.Join(q.Batch.Stories, ", "), q.Batch.Step)
	}
	if len(q.Outcomes) > 0 {
		fmt.Fprintln(w, "outcomes")
		for _, o := range q.Outcomes {
			line := fmt.Sprintf("  %-6s %s %s", o.Story, o.Outcome, o.At.UTC().Format(time.RFC3339))
			if o.Detail != "" {
				line += ": " + o.Detail
			}
			fmt.Fprintln(w, line)
		}
	}
}

// queueLine is a story in review on one line: its ID, title, when it last
// entered review, and whether it is enqueued.
func queueLine(e mergequeue.Entry) string {
	line := fmt.Sprintf("%-6s %-46s in review since %s", e.Story, truncate(e.Title, 46), e.EnteredReview.UTC().Format(time.RFC3339))
	if e.Enqueued {
		line += " enqueued"
	}
	return line
}
