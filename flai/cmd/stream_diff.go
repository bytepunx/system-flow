package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/bytepunx/system-flow/flai/internal/storygit"
)

func newStreamDiffCmd(a *app) *cobra.Command {
	return &cobra.Command{
		Use:   "diff <story-id>",
		Short: "What a story branch changes: files and hunks against the main branch",
		Long: `The story's branch (ADR-0019) against its merge base with the main
branch, read with git in the main checkout, so it does not depend on the
worktree. The dashboard's review page shows this (S-0041).`,
		Example: `  flai stream diff S-0041
  flai stream diff S-0041 --json`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			repo, err := a.project()
			if err != nil {
				return err
			}
			it, err := repo.Get(args[0])
			if err != nil {
				return err
			}
			d, err := storygit.StoryDiff(a.runner, repo, it.ID)
			if err != nil {
				return err
			}
			if a.jsonOut {
				return a.printJSON(d)
			}
			fmt.Fprintf(a.out, "%s against %s: %d commit(s), %d file(s), +%d -%d\n", d.Branch, d.Base, d.Commits, len(d.Files), d.Additions, d.Deletions)
			for _, f := range d.Files {
				name := f.Path
				if f.OldPath != "" {
					name = f.OldPath + " -> " + f.Path
				}
				fmt.Fprintf(a.out, "  %-9s %-60s +%d -%d\n", f.Status, name, f.Additions, f.Deletions)
			}
			return nil
		},
	}
}
