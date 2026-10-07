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
	var clear, add, remove bool
	c := &cobra.Command{
		Use:   "touches <id> [path-or-component...]",
		Short: "Set what a story or task is working on; flai check warns on overlap",
		Long: `Paths or components given alone replace the item's touches: name every one
it keeps. --add adds those given to the list and --remove takes them out of it,
leaving the rest; --clear empties it. With no path, the touches are shown.`,
		Example: `  flai touches S-0037 flai/internal/workitem flaiover/src/routes/docs   # replace
  flai touches S-0037 --add docs/users/flai.md
  flai touches S-0037 --remove flaiover/src/routes/docs
  flai touches T-0121 --clear
  flai touches S-0037            # show`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if (add || remove) && len(args) == 1 {
				return fmt.Errorf("--add and --remove need the paths to change: flai touches %s --add <path>", args[0])
			}
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
			var told []itemedit.Overlapping
			if clear || len(args) > 1 {
				before := slices.Clone(it.Touches)
				watch := itemedit.WatchClaim(repo, it.ID)
				given, err := workitem.CleanTouches(args[1:])
				if err != nil {
					return err
				}
				if it.Touches, err = changedTouches(it, given, add, remove); err != nil {
					return err
				}
				it.Updated = a.now().UTC().Format(workitem.TimeFormat)
				if err := repo.Save(it); err != nil {
					return err
				}
				// an edit notice, as flai edit leaves, so that agents hear of it
				// and flai serve plans the story again (S-0211)
				if !slices.Equal(before, it.Touches) {
					by := a.writer()
					itemedit.Record(repo, itemedit.Notice{At: it.Updated, By: by, ID: it.ID, Type: it.Type, Title: it.Title, Changed: []string{"touches"}})
					told = a.grown(watch, it.ID, by)
				}
			}
			if a.jsonOut {
				out := map[string]any{"id": it.ID, "touches": it.Touches}
				if len(told) > 0 {
					out["overlaps"] = told
				}
				return a.printJSON(out)
			}
			if len(it.Touches) == 0 {
				fmt.Fprintf(a.out, "%s touches nothing\n", it.ID)
				return nil
			}
			fmt.Fprintf(a.out, "%s touches %s\n", it.ID, strings.Join(it.Touches, ", "))
			printOverlapping(a.out, told)
			return nil
		},
	}
	c.Flags().BoolVar(&clear, "clear", false, "remove the list")
	c.Flags().BoolVar(&add, "add", false, "add the paths given to the list rather than replace it")
	c.Flags().BoolVar(&remove, "remove", false, "take the paths given out of the list rather than replace it")
	c.MarkFlagsMutuallyExclusive("clear", "add", "remove")
	c.AddCommand(newTouchesSuggestCmd(a))
	return c
}

// changedTouches is the item's touches after the paths given: those paths alone, the
// list with them added, or the list with them taken out (I-0067). Removing a
// path the item does not touch is refused, so that a mistyped one is not
// taken for done.
func changedTouches(it *workitem.Item, given []string, add, remove bool) ([]string, error) {
	switch {
	case add:
		out := slices.Clone(it.Touches)
		for _, p := range given {
			if !slices.Contains(out, p) {
				out = append(out, p)
			}
		}
		return out, nil
	case remove:
		for _, p := range given {
			if !slices.Contains(it.Touches, p) {
				return nil, fmt.Errorf("%s does not touch %s; it touches %s", it.ID, p, touchesList(it.Touches))
			}
		}
		var out []string
		for _, p := range it.Touches {
			if !slices.Contains(given, p) {
				out = append(out, p)
			}
		}
		return out, nil
	}
	return given, nil
}

// touchesList is touches joined for a message, or nothing when there are none.
func touchesList(touches []string) string {
	if len(touches) == 0 {
		return "nothing"
	}
	return strings.Join(touches, ", ")
}

// writer is who writes an item, as agents are told: FLAI_AGENT, else the
// config's author.
func (a *app) writer() string {
	by, _ := agentIdentity()
	if cfg, _, err := a.loadConfig(); err == nil && by == "agent" && cfg.Author != "" {
		by = cfg.Author
	}
	return by
}

// grown is what the write to id grew its story's claim into, told to both
// stories as by's (I-0059). A failure is logged, and the write stands: the
// report is advisory.
func (a *app) grown(w *itemedit.ClaimWatch, id, by string) []itemedit.Overlapping {
	told, err := w.Grown(by, a.now())
	if err != nil {
		a.logger().Warn("overlap notices not sent", "component", "touches", "item", id, "err", err)
	}
	return told
}

// printOverlapping writes one line for each story in progress whose claim a
// write grew its story's claim into.
func printOverlapping(w io.Writer, told []itemedit.Overlapping) {
	for _, o := range told {
		fmt.Fprintf(w, "  overlaps %s %s (in progress) on %s: both stories are told; message its agent (flai message send) before you change them\n", o.Story, o.Title, strings.Join(o.Paths, ", "))
	}
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
