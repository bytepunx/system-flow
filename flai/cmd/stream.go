package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

func newStreamCmd(a *app) *cobra.Command {
	c := &cobra.Command{
		Use:   "stream",
		Short: "Open and append to agent narratives in wip/agents",
		Long: `A stream is the narrative for one story. Set FLAI_AGENT and FLAI_SESSION
so entries record who wrote them.`,
	}
	c.AddCommand(newStreamOpenCmd(a), newStreamLogCmd(a))
	return c
}

func newStreamOpenCmd(a *app) *cobra.Command {
	return &cobra.Command{
		Use:   "open <story-id>",
		Short: "Create the narrative for a story from the template",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			repo, err := a.project()
			if err != nil {
				return err
			}
			story, err := repo.Get(args[0])
			if err != nil {
				return err
			}
			agent, session := agentIdentity()
			n, err := repo.OpenStream(story, workitem.StreamOptions{Agent: agent, Session: session, Now: a.now()})
			if err != nil {
				return err
			}
			if err := a.refreshIndex(repo); err != nil {
				return err
			}
			if a.jsonOut {
				return a.printJSON(map[string]string{"stream": n.Stream, "path": n.Path})
			}
			fmt.Fprintf(a.out, "opened %s\n", relPath(repo.Root, n.Path))
			return nil
		},
	}
}

func newStreamLogCmd(a *app) *cobra.Command {
	return &cobra.Command{
		Use:   "log <story-id> \"<entry>\"",
		Short: "Append a timestamped entry to a story's narrative log",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			repo, err := a.project()
			if err != nil {
				return err
			}
			agent, session := agentIdentity()
			n, err := repo.LogStream(args[0], args[1], workitem.StreamOptions{Agent: agent, Session: session, Now: a.now()})
			if err != nil {
				return err
			}
			if err := a.refreshIndex(repo); err != nil {
				return err
			}
			if a.jsonOut {
				return a.printJSON(map[string]string{"stream": n.Stream, "updated": n.Updated})
			}
			fmt.Fprintf(a.out, "logged to %s at %s\n", n.Stream, n.Updated)
			return nil
		},
	}
}
