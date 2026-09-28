package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	ctxpack "github.com/bytepunx/system-flow/flai/internal/context"
	"github.com/bytepunx/system-flow/flai/internal/conventions"
	"github.com/bytepunx/system-flow/flai/internal/issues"
	"github.com/bytepunx/system-flow/flai/internal/topics"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

func newPrimeCmd(a *app) *cobra.Command {
	var (
		cat   bool
		story string
	)
	c := &cobra.Command{
		Use:   "prime",
		Short: "Print the conventions an agent reads at session start, in order",
		Long: `List design/conventions in read order (README first, then by order) so an
agent, a hook, or a script can load them in one call. --cat prints the
content of each file with a header instead of the paths.

--story S-nnnn prints the context pack for an agent working that story
(ADR-0047): a header naming the story, its topics and where each came from,
and the size of what follows; then every convention as --cat prints it, with
the sections whose topics include neither all nor one of the story's left
out (the front matter, the baseline marker, and the Project additions heading
always stay); the open issues; and one line per section left out, with its
topics. With every section at [all] the text below the header is --cat's.`,
		Example: `  flai prime
  flai prime --cat
  flai prime --json
  flai prime --story S-0136
  flai prime --story S-0136 --json`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			repo, err := a.project()
			if err != nil {
				return err
			}
			set, errs, err := conventions.Load(repo)
			if err != nil {
				return err
			}
			if set.Missing {
				return fmt.Errorf("no conventions folder at %s; render it from the template or run flai upgrade", relPath(repo.Root, set.Dir))
			}
			readme := filepath.Join(set.Dir, "README.md")
			var paths []string
			if set.README != "" {
				paths = append(paths, readme)
			}
			for _, f := range set.Files {
				paths = append(paths, filepath.Join(repo.Root, f.Path))
			}
			for rel, e := range errs {
				a.logger().Warn("convention file unreadable", "component", "conventions", "path", rel, "err", e.Error())
			}
			list, _ := issues.List(repo)
			table := issues.SummaryTable(list)
			if story != "" {
				return a.primeStory(repo, set, story, table)
			}
			if a.jsonOut {
				return a.printJSON(map[string]any{"dir": relPath(repo.Root, set.Dir), "files": set.Files, "readme": set.README != "", "open_issues": len(strings.Split(strings.TrimSpace(table), "\n")) - 2})
			}
			if !cat {
				for _, p := range paths {
					fmt.Fprintln(a.out, relPath(repo.Root, p))
				}
				if table != "" {
					fmt.Fprintln(a.out, relPath(repo.Root, filepath.Join(issues.Dir(repo), issues.SummaryFile)))
				}
				return nil
			}
			for i, p := range paths {
				data, err := os.ReadFile(p)
				if err != nil {
					return err
				}
				if i > 0 {
					fmt.Fprintln(a.out)
				}
				rel := relPath(repo.Root, p)
				fmt.Fprintf(a.out, "%s\n%s\n\n", rel, strings.Repeat("=", len(rel)))
				fmt.Fprint(a.out, string(data))
			}
			if table != "" {
				fmt.Fprintf(a.out, "\nopen issues\n===========\n\n%s", table)
			}
			return nil
		},
	}
	c.Flags().BoolVar(&cat, "cat", false, "print file contents instead of paths")
	c.Flags().StringVar(&story, "story", "", "print the context pack for this story: the convention sections its topics select, and what was left out")
	return c
}

// primeStory prints the context pack for an active story (ADR-0047, S-0136).
func (a *app) primeStory(repo *workitem.Repo, set *conventions.Set, id, table string) error {
	it, err := repo.Get(id)
	if err != nil {
		return fmt.Errorf("flai prime --story %s: %w; give the ID of an active story", id, err)
	}
	if it.Type != workitem.Story {
		return fmt.Errorf("flai prime --story %s: %s is %s, not a story", id, it.ID, it.Type)
	}
	if it.Archived {
		return fmt.Errorf("flai prime --story %s: %s is archived; a pack is for a story still being worked", id, it.ID)
	}
	storyTopics, err := topics.ForStory(repo, it.ID)
	if err != nil {
		return err
	}
	pack, err := ctxpack.Build(repo.Root, it.ID, it.Title, storyTopics, set, table)
	if err != nil {
		return err
	}
	if a.jsonOut {
		return a.printJSON(pack)
	}
	fmt.Fprint(a.out, pack.Header())
	fmt.Fprint(a.out, pack.Body())
	return nil
}
