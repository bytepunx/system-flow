package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/bytepunx/system-flow/flai/internal/conventions"
	"github.com/bytepunx/system-flow/flai/internal/issues"
)

func newPrimeCmd(a *app) *cobra.Command {
	var cat bool
	c := &cobra.Command{
		Use:   "prime",
		Short: "Print the conventions an agent reads at session start, in order",
		Long: `List design/conventions in read order (README first, then by order) so an
agent, a hook, or a script can load them in one call. --cat prints the
content of each file with a header instead of the paths.`,
		Example: `  flai prime
  flai prime --cat
  flai prime --json`,
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
	return c
}
