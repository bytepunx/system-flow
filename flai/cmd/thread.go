package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/bytepunx/system-flow/flai/internal/threads"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// threadAuthor is --by, else the agent identity, else the config author.
func (a *app) threadAuthor(by string) string {
	if by != "" {
		return by
	}
	if agent, _ := agentIdentity(); agent != "" {
		return agent
	}
	return a.author()
}

func newThreadCmd(a *app) *cobra.Command {
	c := &cobra.Command{
		Use:   "thread",
		Short: "Threads between the designer and agents, anchored to documents and items (wip/threads)",
	}
	c.AddCommand(newThreadNewCmd(a), newThreadReplyCmd(a), newThreadConfirmCmd(a), newThreadResolveCmd(a), newThreadListCmd(a), newThreadShowCmd(a))
	return c
}

func (a *app) afterThread(repo *workitem.Repo, th *threads.Thread) error {
	if story := threads.StoryOf(repo, th); story != "" {
		return threads.MirrorNarrative(repo, story)
	}
	return nil
}

func newThreadNewCmd(a *app) *cobra.Command {
	var on, heading, by string
	c := &cobra.Command{
		Use:   "new \"<title>\" \"<text>\"",
		Short: "Open a thread on a document, a heading in it, or a work item",
		Example: `  flai thread new --on design/system/flaiover-dashboard.md --heading "Workbench (E-0006)" "Should threads mirror into docs?" "Question text"
  flai thread new --on S-0038 "Anchor format" "Do we need line anchors?"`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			repo, err := a.project()
			if err != nil {
				return err
			}
			th, err := threads.New(repo, threads.NewOptions{Title: args[0], Text: args[1], On: on, Heading: heading, Author: a.threadAuthor(by), Now: a.now()})
			if err != nil {
				return err
			}
			if err := a.afterThread(repo, th); err != nil {
				return err
			}
			if a.jsonOut {
				return a.printJSON(threadJSON(repo, th))
			}
			fmt.Fprintf(a.out, "%s %s\n  on %s\n  %s\n", th.ID, th.Title, describeAnchor(th.Anchor), relPath(repo.MainRoot, th.Path))
			return nil
		},
	}
	c.Flags().StringVar(&on, "on", "", "repository path or item ID the thread is about")
	c.Flags().StringVar(&heading, "heading", "", "heading in the document (must exist)")
	c.Flags().StringVar(&by, "by", "", "author (default: FLAI_AGENT, then config author)")
	_ = c.MarkFlagRequired("on")
	return c
}

func newThreadReplyCmd(a *app) *cobra.Command {
	var by, source string
	var recommend bool
	c := &cobra.Command{
		Use:   "reply <id> \"<text>\"",
		Short: "Add an entry; answered when someone other than the opener replies, unless it is a recommendation",
		Long: `Adds an entry to a thread. A reply from anyone but the opener marks the thread answered, and one from the opener marks it open.

With --recommend the entry is a recommendation (ADR-0090): its heading says (recommendation), the status stays as it was, and the thread awaits the operator, who makes it the answer with flai thread confirm or answers otherwise. A recommendation is refused on a resolved thread and from the thread's opener.

--source cites what the reply rests on: a file in the repository, or an item ID, and after # a heading in it. Both must exist. It is written as the entry's last line, Source: <path> § <heading>.`,
		Example: `  flai thread reply TH-0012 "Nine metres."
  flai thread reply TH-0012 --recommend --source "design/system/overview.md#Delivery sequence" "Build the CLI first, as the delivery sequence says."`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			repo, err := a.project()
			if err != nil {
				return err
			}
			src, err := threads.ParseSource(source)
			if err != nil {
				return err
			}
			th, err := threads.ReplyWith(repo, args[0], a.threadAuthor(by), args[1], a.now(), threads.Marks{Recommendation: recommend, Source: src})
			if err != nil {
				return err
			}
			if err := a.afterThread(repo, th); err != nil {
				return err
			}
			if a.jsonOut {
				return a.printJSON(threadJSON(repo, th))
			}
			fmt.Fprintf(a.out, "%s is %s (%d entries)\n", th.ID, th.Status, len(th.Entries()))
			if recommend {
				fmt.Fprintf(a.out, "  the recommendation awaits the operator: flai thread confirm %s\n", th.ID)
			}
			return nil
		},
	}
	c.Flags().StringVar(&by, "by", "", "author (default: FLAI_AGENT, then config author)")
	c.Flags().BoolVar(&recommend, "recommend", false, "post a recommendation: the status stays and the thread awaits the operator's confirmation (ADR-0090)")
	c.Flags().StringVar(&source, "source", "", "what the reply rests on, <path>[#<heading>]: a file in the repository and a heading in it")
	return c
}

func newThreadConfirmCmd(a *app) *cobra.Command {
	var by string
	c := &cobra.Command{
		Use:   "confirm <id>",
		Short: "Make the pending recommendation the answer: an entry citing its source, and the thread answered",
		Long: `Makes the thread's pending recommendation the answer (ADR-0090). It adds an entry by you, "Confirmed the recommendation of <time> <author>.", citing the recommendation's source, and marks the thread answered.

A recommendation is pending while the thread is not resolved and it is the newest entry by someone other than the opener. Confirm is refused when none is pending, and to the recommendation's own author.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			repo, err := a.project()
			if err != nil {
				return err
			}
			th, err := threads.Confirm(repo, args[0], a.threadAuthor(by), a.now())
			if err != nil {
				return err
			}
			if err := a.afterThread(repo, th); err != nil {
				return err
			}
			if a.jsonOut {
				return a.printJSON(threadJSON(repo, th))
			}
			fmt.Fprintf(a.out, "%s is %s (%d entries)\n", th.ID, th.Status, len(th.Entries()))
			return nil
		},
	}
	c.Flags().StringVar(&by, "by", "", "who confirms it (default: FLAI_AGENT, then config author)")
	return c
}

func newThreadResolveCmd(a *app) *cobra.Command {
	var by, reason string
	c := &cobra.Command{
		Use:   "resolve <id>",
		Short: "Close a thread",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			repo, err := a.project()
			if err != nil {
				return err
			}
			th, err := threads.Resolve(repo, args[0], a.threadAuthor(by), reason, a.now())
			if err != nil {
				return err
			}
			if err := a.afterThread(repo, th); err != nil {
				return err
			}
			if a.jsonOut {
				return a.printJSON(threadJSON(repo, th))
			}
			fmt.Fprintf(a.out, "%s resolved\n", th.ID)
			return nil
		},
	}
	c.Flags().StringVar(&by, "by", "", "who resolves it (default: FLAI_AGENT, then config author)")
	c.Flags().StringVar(&reason, "reason", "", "what settled it")
	return c
}

func newThreadListCmd(a *app) *cobra.Command {
	var on string
	var all bool
	c := &cobra.Command{
		Use:   "list",
		Short: "List threads; unresolved by default",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			repo, err := a.project()
			if err != nil {
				return err
			}
			var list []*threads.Thread
			if on != "" {
				list, err = threads.For(repo, on)
			} else {
				list, err = threads.List(repo)
			}
			if err != nil {
				return err
			}
			var out []*threads.Thread
			for _, th := range list {
				if all || th.Open() {
					out = append(out, th)
				}
			}
			if a.jsonOut {
				rows := make([]map[string]any, 0, len(out))
				for _, th := range out {
					rows = append(rows, threadJSON(repo, th))
				}
				return a.printJSON(rows)
			}
			if len(out) == 0 {
				fmt.Fprintln(a.out, "no threads")
				return nil
			}
			for _, th := range out {
				fmt.Fprintf(a.out, "%-8s %-9s %-10s %s  [%s]\n", th.ID, th.Status, th.Updated[:10], th.Title, describeAnchor(th.Anchor))
			}
			return nil
		},
	}
	c.Flags().StringVar(&on, "on", "", "only threads on this path or item")
	c.Flags().BoolVar(&all, "all", false, "include resolved threads")
	return c
}

func newThreadShowCmd(a *app) *cobra.Command {
	return &cobra.Command{
		Use:   "show <id>",
		Short: "Print a thread with its entries",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			repo, err := a.project()
			if err != nil {
				return err
			}
			th, err := threads.Get(repo, args[0])
			if err != nil {
				return err
			}
			if a.jsonOut {
				return a.printJSON(threadJSON(repo, th))
			}
			fmt.Fprintf(a.out, "%s %s\n  %s · on %s · %s\n", th.ID, th.Title, th.Status, describeAnchor(th.Anchor), strings.Join(th.Participants, ", "))
			if p := th.PendingRecommendation(); p != nil {
				fmt.Fprintf(a.out, "  the recommendation of %s %s awaits the operator: flai thread confirm %s\n", p.At, p.Author, th.ID)
			}
			for _, e := range th.Entries() {
				mark := ""
				if e.Recommendation {
					mark = " (recommendation)"
				}
				fmt.Fprintf(a.out, "\n%s %s%s\n%s\n", e.At, e.Author, mark, e.Text)
			}
			return nil
		},
	}
}

func describeAnchor(an threads.Anchor) string {
	s := an.Path
	if an.Item != "" {
		s = an.Item + " (" + an.Path + ")"
	}
	if an.Heading != "" {
		s += " § " + an.Heading
	}
	return s
}

func threadJSON(repo *workitem.Repo, th *threads.Thread) map[string]any {
	return threads.View(repo, th)
}
