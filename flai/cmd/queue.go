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
	c := &cobra.Command{
		Use:   "queue",
		Short: "List the stories in review the merge queue would batch, and why each other would not",
		Long: `List the stories in review that the merge queue would batch, in the order
they last entered review: the last flai verify passed at the head of the
story's branch, every acceptance criterion is ticked, no thread on the story
or its tasks is open, and no path is recorded as undoing main's change
(ADR-0155). Each other story in review follows, with every reason it is not
eligible. Then the batch under verification, its step, and its result, and
each story's last outcome, from the queue's state in
.flai-cache/merge-queue; first, when the queue is held, who held it and when.

The listing verifies, merges, and moves nothing (ADR-0173). flai queue add
and remove enqueue a story and take it out, and flai queue hold and resume
stop the queue building batches and let it again; flai accept takes the
story it accepts out of the queue. Under acceptance.mode pull-request the
queue does not run, because the forge's required checks gate each story
(ADR-0169), and no story is listed.`,
		Example: `  flai queue
  flai queue --json
  flai queue add S-0042
  flai queue hold`,
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
	c.AddCommand(newQueueAddCmd(a), newQueueRemoveCmd(a), newQueueHoldCmd(a), newQueueResumeCmd(a))
	return c
}

func newQueueAddCmd(a *app) *cobra.Command {
	return &cobra.Command{
		Use:   "add <story>",
		Short: "Enqueue a story in review that is eligible for the merge queue",
		Long: `Enqueue a story in review in the merge queue, recorded with who enqueued it
(FLAI_AGENT, else the config's author) and when. The queue lands each story
as the one who enqueued it (ADR-0173).

Refused, and nothing changes, for a story not in review, for one not
eligible, naming every reason flai queue gives, and under acceptance.mode
pull-request, where the queue does not run. A story enqueued already keeps
its place and who enqueued it, and nothing changes.`,
		Example: `  flai queue add S-0042
  flai queue add S-0042 --json`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			repo, err := a.project()
			if err != nil {
				return err
			}
			ch, err := mergequeue.Add(a.runner, repo, args[0], a.movedBy(""), a.now())
			if err != nil {
				return err
			}
			return a.printQueueChange(ch)
		},
	}
}

func newQueueRemoveCmd(a *app) *cobra.Command {
	return &cobra.Command{
		Use:   "remove <story>",
		Short: "Take a story out of the merge queue",
		Long: `Take a story out of the merge queue. The batch under verification that holds
it is marked dropped, so that the queue does not land a story no one has
enqueued. Refused for a story not enqueued.`,
		Example: `  flai queue remove S-0042`,
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			repo, err := a.project()
			if err != nil {
				return err
			}
			ch, err := mergequeue.Remove(repo, args[0], a.movedBy(""), a.now())
			if err != nil {
				return err
			}
			return a.printQueueChange(ch)
		},
	}
}

func newQueueHoldCmd(a *app) *cobra.Command {
	return &cobra.Command{
		Use:   "hold",
		Short: "Hold the merge queue: it builds no batch until resumed",
		Long: `Hold the merge queue, recorded with who held it (FLAI_AGENT, else the
config's author) and when: it builds no batch until flai queue resume. The
stories enqueued keep their places. A queue held already keeps who held it
first, and nothing changes.`,
		Example: `  flai queue hold`,
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			repo, err := a.project()
			if err != nil {
				return err
			}
			ch, err := mergequeue.Hold(repo, a.movedBy(""), a.now())
			if err != nil {
				return err
			}
			return a.printQueueChange(ch)
		},
	}
}

func newQueueResumeCmd(a *app) *cobra.Command {
	return &cobra.Command{
		Use:   "resume",
		Short: "Resume a held merge queue",
		Long: `Resume a held merge queue, so that it builds batches again. Who resumed it
and when are printed and logged. A queue not held is left as it is.`,
		Example: `  flai queue resume`,
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			repo, err := a.project()
			if err != nil {
				return err
			}
			ch, err := mergequeue.Resume(repo, a.movedBy(""), a.now())
			if err != nil {
				return err
			}
			return a.printQueueChange(ch)
		},
	}
}

// printQueueChange logs a change that wrote the queue's state and writes it
// as JSON under --json, else as a line saying what changed, or that nothing
// did and why.
func (a *app) printQueueChange(ch mergequeue.Change) error {
	if ch.Changed {
		a.logger().Info("merge queue changed", "component", "mergequeue", "action", ch.Action, "story", ch.Story, "by", ch.By)
	}
	if a.jsonOut {
		return a.printJSON(ch)
	}
	at := ch.At.Format(time.RFC3339)
	st := ch.State
	switch ch.Action {
	case mergequeue.ActionAdd:
		if e, ok := st.EnqueuedAs(ch.Story); ok && !ch.Changed {
			fmt.Fprintf(a.out, "%s is enqueued already, by %s at %s; nothing changed\n", ch.Story, e.By, e.At.UTC().Format(time.RFC3339))
			return nil
		}
		fmt.Fprintf(a.out, "%s enqueued by %s at %s, %d in the queue\n", ch.Story, ch.By, at, len(st.Enqueued))
	case mergequeue.ActionRemove:
		line := fmt.Sprintf("%s removed from the merge queue by %s at %s", ch.Story, ch.By, at)
		if st.HoldsStory(ch.Story) && st.Batch.Dropped != nil {
			line += "; the batch that holds it is dropped"
		}
		fmt.Fprintln(a.out, line)
	case mergequeue.ActionHold:
		if !ch.Changed {
			fmt.Fprintf(a.out, "merge queue held already, by %s at %s; nothing changed\n", st.HeldBy, st.HeldAt.UTC().Format(time.RFC3339))
			return nil
		}
		fmt.Fprintf(a.out, "merge queue held by %s at %s: it builds no batch until flai queue resume\n", ch.By, at)
	case mergequeue.ActionResume:
		if !ch.Changed {
			fmt.Fprintln(a.out, "merge queue is not held; nothing changed")
			return nil
		}
		fmt.Fprintf(a.out, "merge queue resumed by %s at %s\n", ch.By, at)
	}
	return nil
}

// printQueue writes q as text: the eligible stories in order, each other
// story in review with its reasons, the batch, and the outcomes; or, when
// the queue is off, the one line saying why.
func printQueue(w io.Writer, q mergequeue.Queue) {
	if q.Off {
		fmt.Fprintf(w, "merge queue off under acceptance.mode %s: %s\n", q.Mode, q.OffReason)
		return
	}
	if q.Held {
		fmt.Fprintf(w, "held by %s at %s: no batch is built until flai queue resume\n", q.HeldBy, q.HeldAt.UTC().Format(time.RFC3339))
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
	printBatch(w, q.Batch)
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

// printBatch writes the batch under verification and its step, and, once
// it has them, its result and why it was dropped.
func printBatch(w io.Writer, b *mergequeue.Batch) {
	if b == nil {
		fmt.Fprintln(w, "batch: none")
		return
	}
	fmt.Fprintf(w, "batch: %s, step %s\n", strings.Join(b.Stories, ", "), b.Step)
	if r := b.Result; r != nil {
		line := "  result: passed"
		if !r.Passed {
			line = "  result: failed"
			if r.StoppedAt != "" {
				line += " at " + r.StoppedAt
			}
		}
		line += ", finished " + r.FinishedAt.UTC().Format(time.RFC3339)
		if r.Summary != "" {
			line += ": " + r.Summary
		}
		fmt.Fprintln(w, line)
	}
	if d := b.Dropped; d != nil {
		fmt.Fprintf(w, "  dropped by %s at %s: %s\n", d.By, d.At.UTC().Format(time.RFC3339), d.Reason)
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
