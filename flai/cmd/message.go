package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/bytepunx/system-flow/flai/internal/messages"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// fromHelp is the --from flag's help for message send and reply.
const fromHelp = "the story the message comes from (default: FLAI_STORY, else the story in FLAI_AGENT of the form agent-S-nnnn, else the story branch checked out here)"

func newMessageCmd(a *app) *cobra.Command {
	c := &cobra.Command{
		Use:   "message",
		Short: "Conversations between the agents of two open stories, apart from the operator's threads (wip/messages)",
		Long: `Conversations between the agents of two stories in progress or in review (ADR-0120). A message goes from one story to another, and the messages between the two make a conversation, one file in wip/messages, MS-nnnn-<slug>.md, in the main checkout.

A conversation awaits the story that did not write its last message. It reads as closed when its status is closed or either story is done, cancelled, or archived, and a closed conversation takes no reply: a new message starts a new one. flai accept and flai archive close the conversations of the stories they archive.

Messages are kept apart from threads: none appears in flai thread list, among the threads awaiting the operator, or in a narrative's Open questions. A question for the designer is still a thread.`,
	}
	c.AddCommand(newMessageSendCmd(a), newMessageReplyCmd(a), newMessageListCmd(a), newMessageShowCmd(a))
	return c
}

// messageStory is the story a message comes from: --from, else FLAI_STORY,
// else the story in FLAI_AGENT of the form agent-S-nnnn, else the story whose
// branch is checked out. None is an error, since a message needs a sender.
func (a *app) messageStory(repo *workitem.Repo, from string) (string, error) {
	if s := a.issueStory(repo, from); s != "" {
		return s, nil
	}
	return "", fmt.Errorf("no story to send from: give --from S-nnnn, set FLAI_STORY, run as an agent named agent-S-nnnn, or run in the story's worktree; nothing was written")
}

func newMessageSendCmd(a *app) *cobra.Command {
	var from, by string
	var about []string
	c := &cobra.Command{
		Use:   "send <S-nnnn> \"<text>\"",
		Short: "Start a conversation from one open story to another with its first message",
		Long: `Starts a conversation from the sender's story to the story <S-nnnn> with its first message, and prints its ID, title, the two stories, and its path. The title is the first line of the message.

The sender's story is --from, else FLAI_STORY, else the story in FLAI_AGENT of the form agent-S-nnnn, else the story branch checked out here; with none, the send is refused. Both stories must be in progress or in review, and not the same story: a send to or from a story in any other state, or archived, is refused with its state, and nothing is written.

--about names a repository path the conversation is about, a file or folder relative to the root that must exist; give it once per path. The author of the message is --by, else FLAI_AGENT, else the config author.`,
		Example: `  flai message send S-0331 "I am changing flai/cmd/root.go. Will you leave it to me until I push?" --from S-0330
  flai message send S-0331 "Who adds the docs row?" --about docs/operators/settings.md --json`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			repo, err := a.project()
			if err != nil {
				return err
			}
			sender, err := a.messageStory(repo, from)
			if err != nil {
				return err
			}
			c, err := messages.Send(repo, messages.SendOptions{From: sender, To: args[0], Author: a.threadAuthor(by), Text: args[1], About: about, Now: a.now()})
			if err != nil {
				return err
			}
			if a.jsonOut {
				return a.printJSON(messages.View(repo, c))
			}
			between := c.From + " → " + c.To
			if len(c.About) > 0 {
				between += ", about " + strings.Join(c.About, ", ")
			}
			fmt.Fprintf(a.out, "%s %s\n  %s\n  %s\n", c.ID, c.Title, between, relPath(repo.MainRoot, c.Path))
			return nil
		},
	}
	c.Flags().StringVar(&from, "from", "", fromHelp)
	c.Flags().StringArrayVar(&about, "about", nil, "repository path the conversation is about, relative to the root (repeatable)")
	c.Flags().StringVar(&by, "by", "", "author (default: FLAI_AGENT, then config author)")
	return c
}

func newMessageReplyCmd(a *app) *cobra.Command {
	var from, by string
	c := &cobra.Command{
		Use:   "reply <MS-nnnn> \"<text>\"",
		Short: "Add a message to a conversation from one of its two stories",
		Long: `Adds a message to the conversation <MS-nnnn> from one of its two stories, and prints whom it now awaits: the other story.

The replying story is --from, else FLAI_STORY, else the story in FLAI_AGENT of the form agent-S-nnnn, else the story branch checked out here. A reply is refused from a story that is not one of the two, from a story no longer in progress or in review, and on a conversation that reads as closed, which takes none; nothing is written. The author is --by, else FLAI_AGENT, else the config author.`,
		Example: `  flai message reply MS-0004 "Yes, it is yours until you push." --from S-0331`,
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			repo, err := a.project()
			if err != nil {
				return err
			}
			story, err := a.messageStory(repo, from)
			if err != nil {
				return err
			}
			c, err := messages.Reply(repo, args[0], story, a.threadAuthor(by), args[1], a.now())
			if err != nil {
				return err
			}
			if a.jsonOut {
				return a.printJSON(messages.View(repo, c))
			}
			fmt.Fprintf(a.out, "%s awaits %s (%d entries)\n", c.ID, c.Awaiting(), len(c.Entries()))
			return nil
		},
	}
	c.Flags().StringVar(&from, "from", "", fromHelp)
	c.Flags().StringVar(&by, "by", "", "author (default: FLAI_AGENT, then config author)")
	return c
}

func newMessageListCmd(a *app) *cobra.Command {
	var story string
	var all bool
	c := &cobra.Command{
		Use:   "list",
		Short: "List conversations; those still open by default",
		Long: `Lists conversations in ID order, one a line: its ID, the story that started it and the one it went to, the story it awaits, the time of its last message, and its title.

By default only the conversations that do not read as closed are listed, those of every story. --story lists only the conversations the story is one of the two of, in any padding; it has no default, so that the operator sees every story's. --all adds the closed ones, each with why it is closed in place of the story it awaits.`,
		Example: `  flai message list --story S-0330
  flai message list --all --json`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			repo, err := a.project()
			if err != nil {
				return err
			}
			var list []*messages.Conversation
			if story != "" {
				list, err = messages.For(repo, story)
			} else {
				list, err = messages.List(repo)
			}
			if err != nil {
				return err
			}
			type row struct {
				c     *messages.Conversation
				state string
			}
			var rows []row
			for _, c := range list {
				closed, why := c.Closed(repo)
				switch {
				case !closed:
					rows = append(rows, row{c, "awaiting " + c.Awaiting()})
				case all:
					rows = append(rows, row{c, "closed: " + why})
				}
			}
			if a.jsonOut {
				views := make([]map[string]any, 0, len(rows))
				for _, r := range rows {
					views = append(views, messages.View(repo, r.c))
				}
				return a.printJSON(views)
			}
			if len(rows) == 0 {
				fmt.Fprintln(a.out, noConversations(story, all))
				return nil
			}
			for _, r := range rows {
				fmt.Fprintf(a.out, "%-8s %s → %s  %-15s  %s  %s\n", r.c.ID, r.c.From, r.c.To, r.state, lastAt(r.c), r.c.Title)
			}
			return nil
		},
	}
	c.Flags().StringVar(&story, "story", "", "only the conversations this story is one of the two of")
	c.Flags().BoolVar(&all, "all", false, "include the conversations that read as closed")
	return c
}

// noConversations says plainly that the list is empty, and of what.
func noConversations(story string, all bool) string {
	s := "no open conversations"
	if all {
		s = "no conversations"
	}
	if story != "" {
		s += " for " + workitem.CanonicalID(story)
	}
	return s
}

// lastAt is the time of a conversation's last entry, else its updated time.
func lastAt(c *messages.Conversation) string {
	if es := c.Entries(); len(es) > 0 {
		return es[len(es)-1].At
	}
	return c.Updated
}

func newMessageShowCmd(a *app) *cobra.Command {
	return &cobra.Command{
		Use:   "show <MS-nnnn>",
		Short: "Print a conversation with every message",
		Long:  `Prints the conversation <MS-nnnn>: its ID and title; the two stories, its status, and the story it awaits, or why it reads as closed; the paths it is about; and every entry, with its time, author, story, and text. The ID may be given in any padding, such as ms-4.`,
		Example: `  flai message show MS-0004
  flai message show ms-4 --json`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			repo, err := a.project()
			if err != nil {
				return err
			}
			c, err := messages.Get(repo, args[0])
			if err != nil {
				return err
			}
			if a.jsonOut {
				return a.printJSON(messages.View(repo, c))
			}
			state := c.Status + " · awaiting " + c.Awaiting()
			if closed, why := c.Closed(repo); closed {
				state = "closed: " + why
			}
			fmt.Fprintf(a.out, "%s %s\n  between %s and %s · %s\n", c.ID, c.Title, c.From, c.To, state)
			if len(c.About) > 0 {
				fmt.Fprintf(a.out, "  about %s\n", strings.Join(c.About, ", "))
			}
			for _, e := range c.Entries() {
				fmt.Fprintf(a.out, "\n%s\n%s\n", strings.TrimSpace(e.At+" "+e.Author+" "+e.Story), e.Text)
			}
			return nil
		},
	}
}
