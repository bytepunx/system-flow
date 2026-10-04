package cmd

import (
	"fmt"
	"io"
	"slices"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/bytepunx/system-flow/flai/internal/itemedit"
	"github.com/bytepunx/system-flow/flai/internal/planning"
	"github.com/bytepunx/system-flow/flai/internal/storygit"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// touches sets the advisory list of paths or components a story or task is
// working on (ADR-0019); flai check warns when in-progress items overlap.
func newTouchesCmd(a *app) *cobra.Command {
	var clear bool
	c := &cobra.Command{
		Use:   "touches <id> [path-or-component...]",
		Short: "Set what a story or task is working on; flai check warns on overlap",
		Example: `  flai touches S-0037 flai/internal/workitem flaiover/src/routes/docs
  flai touches T-0121 --clear
  flai touches S-0037            # show`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			repo, err := a.project()
			if err != nil {
				return err
			}
			it, err := repo.Get(args[0])
			if err != nil {
				return err
			}
			if it.Type == workitem.Epic {
				return fmt.Errorf("%s is an epic; touches belong to stories and tasks", it.ID)
			}
			if clear || len(args) > 1 {
				before := slices.Clone(it.Touches)
				it.Touches = nil
				for _, p := range args[1:] {
					p = strings.TrimSuffix(strings.TrimSpace(p), "/")
					if p != "" {
						it.Touches = append(it.Touches, p)
					}
				}
				it.Updated = a.now().UTC().Format(workitem.TimeFormat)
				if err := repo.Save(it); err != nil {
					return err
				}
				// an edit notice, as flai edit leaves, so that agents hear of it
				// and flai serve plans the story again (S-0211)
				if !slices.Equal(before, it.Touches) {
					by, _ := agentIdentity()
					if cfg, _, err := a.loadConfig(); err == nil && by == "agent" && cfg.Author != "" {
						by = cfg.Author
					}
					itemedit.Record(repo, itemedit.Notice{At: it.Updated, By: by, ID: it.ID, Type: it.Type, Title: it.Title, Changed: []string{"touches"}})
				}
			}
			if a.jsonOut {
				return a.printJSON(map[string]any{"id": it.ID, "touches": it.Touches})
			}
			if len(it.Touches) == 0 {
				fmt.Fprintf(a.out, "%s touches nothing\n", it.ID)
				return nil
			}
			fmt.Fprintf(a.out, "%s touches %s\n", it.ID, strings.Join(it.Touches, ", "))
			return nil
		},
	}
	c.Flags().BoolVar(&clear, "clear", false, "remove the list")
	c.AddCommand(newTouchesSuggestCmd(a))
	return c
}

// touchesSuggestion is what flai touches suggest prints with --json.
type touchesSuggestion struct {
	ID    string   `json:"id"`
	Seeds []string `json:"seeds"`
	planning.Suggestions
}

// newTouchesSuggestCmd lists the files most often changed together with a
// story's touches and the paths given, from the main branch's history
// (S-0210). It changes nothing.
func newTouchesSuggestCmd(a *app) *cobra.Command {
	var minCount, limit int
	c := &cobra.Command{
		Use:   "suggest <S-nnnn> [path...]",
		Short: "List files often changed together with a story's touches, from git history",
		Long: `Seeds are the story's touches, its own and those of its tasks not cancelled,
and the paths given. Every commit on the main branch that changed a file under a
seed counts the other files it changed; flai's own bookkeeping commits and the
wip folder are left out. Nothing is written.`,
		Example: `  flai touches suggest S-0037
  flai touches suggest S-0037 flai/internal/workitem --min 3 --limit 10`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			repo, err := a.project()
			if err != nil {
				return err
			}
			it, err := repo.Get(args[0])
			if err != nil {
				return err
			}
			if it.Type != workitem.Story {
				return fmt.Errorf("%s is not a story; run flai touches suggest on a story: flai touches suggest S-nnnn <path>", it.ID)
			}
			seeds, err := storySeeds(repo, it, args[1:])
			if err != nil {
				return err
			}
			if len(seeds) == 0 {
				return fmt.Errorf("%s touches nothing and no path was given; name the paths to start from: flai touches suggest %s <path>", it.ID, it.ID)
			}
			commits, err := storygit.CommitFiles(a.runner, repo)
			if err != nil {
				return err
			}
			res := touchesSuggestion{ID: it.ID, Seeds: seeds, Suggestions: planning.CoChanged(commits, seeds, minCount, limit)}
			if a.jsonOut {
				return a.printJSON(res)
			}
			printTouchesSuggestion(a.out, res)
			return nil
		},
	}
	c.Flags().IntVar(&minCount, "min", 2, "only files changed in at least this many of the seed commits")
	c.Flags().IntVar(&limit, "limit", 20, "at most this many files; 0 for every one")
	return c
}

// storySeeds is a story's touches, its own and those of its tasks not
// cancelled, and the paths given, each read as a claim reads it.
func storySeeds(repo *workitem.Repo, story *workitem.Item, paths []string) ([]string, error) {
	items, err := repo.List(story.Archived)
	if err != nil {
		return nil, fmt.Errorf("cannot read the tasks of %s: %w; run flai check to see what is wrong", story.ID, err)
	}
	touches := slices.Clone(story.Touches)
	for _, t := range items {
		if t.Type == workitem.Task && t.Status != workitem.Cancelled && workitem.CanonicalID(t.Parent) == workitem.CanonicalID(story.ID) {
			touches = append(touches, t.Touches...)
		}
	}
	touches = append(touches, paths...)
	return workitem.NewHolds(nil, repo.Manifest.Projects).Claim(&workitem.Item{Touches: touches}), nil
}

// printTouchesSuggestion writes a header naming the seeds, then one line per
// file: its path, its count, and its share of the seed commits.
func printTouchesSuggestion(w io.Writer, res touchesSuggestion) {
	fmt.Fprintf(w, "%s from %s: %d of %d commits changed them\n", res.ID, strings.Join(res.Seeds, ", "), res.SeedCommits, res.Commits)
	list := res.Suggestions.Suggestions
	switch {
	case res.SeedCommits == 0:
		fmt.Fprintln(w, "no commit on the main branch changed them, so there is nothing to suggest")
		return
	case len(list) == 0:
		fmt.Fprintln(w, "no other file was changed with them often enough; lower --min to see more")
		return
	}
	// the first count is the largest: the list is ordered by count
	width, counts := 0, len(strconv.Itoa(list[0].Count))
	for _, s := range list {
		width = max(width, len(s.Path))
	}
	for _, s := range list {
		fmt.Fprintf(w, "%-*s  %*d  %3.0f%%\n", width, s.Path, counts, s.Count, 100*s.Share)
	}
}
