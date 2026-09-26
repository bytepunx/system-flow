package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"

	"github.com/spf13/cobra"

	"github.com/bytepunx/system-flow/flai/internal/config"
	"github.com/bytepunx/system-flow/flai/internal/serve"
)

// flai serve import names the folders flai serve looks in for repositories
// the board offers to import (S-0098).
func newServeImportCmd(a *app) *cobra.Command {
	c := &cobra.Command{
		Use:   "import",
		Short: "The folders whose git repositories the board offers to import into system-flow",
		Long: `flai serve looks through each folder named here, and the folders in it to
three levels, for git repositories with no system-flow.yaml, and offers each to
the dashboards the projects it serves connect to: picking one there asks
whether to import it. Importing runs flai import --commit in the repository, as
you and on this host: it writes the standard's files, runs the repository's own
tests, and commits the import when they pass. Naming a folder is your say that
this may happen to what is in it; nothing else enables it.

A repository there that has a system-flow.yaml already, imported on the
command line say, is served as a project (S-0117), as a registered one is,
for as long as it is there and without being written to the registry;
flai serve status lists those it does not serve, and why.`,
		Example: `  flai serve import add ~/git
  flai serve import list
  flai serve import remove ~/git`,
	}
	change := func(add bool) func(*cobra.Command, []string) error {
		return func(_ *cobra.Command, args []string) error {
			dir, err := filepath.Abs(args[0])
			if err != nil {
				return err
			}
			if add {
				if info, err := os.Stat(dir); err != nil || !info.IsDir() {
					return fmt.Errorf("%s is not a folder", dir)
				}
			}
			cfg, path, err := a.loadConfig()
			if err != nil {
				return err
			}
			had := slices.Contains(cfg.ImportRoots, dir)
			switch {
			case add && !had:
				cfg.ImportRoots = append(cfg.ImportRoots, dir)
			case !add && had:
				cfg.ImportRoots = slices.DeleteFunc(cfg.ImportRoots, func(r string) bool { return r == dir })
			}
			if err := config.Save(path, cfg); err != nil {
				return err
			}
			if a.jsonOut {
				return a.printJSON(map[string]any{"roots": cfg.ImportRoots, "config": path})
			}
			if add {
				fmt.Fprintf(a.out, "%s: its git repositories without system-flow.yaml are offered for import, and those with one are served\n  flai serve looks again within 30 seconds; flai serve import list shows what it finds\n", dir)
			} else {
				fmt.Fprintf(a.out, "%s: no longer looked in\n", dir)
			}
			return nil
		}
	}
	c.AddCommand(
		&cobra.Command{Use: "add <folder>", Short: "Offer the repositories in a folder for import", Args: cobra.ExactArgs(1), RunE: change(true)},
		&cobra.Command{Use: "remove <folder>", Short: "Stop offering the repositories in a folder", Args: cobra.ExactArgs(1), RunE: change(false)},
		&cobra.Command{Use: "list", Short: "The folders named, and the repositories found in them", Args: cobra.NoArgs, RunE: func(*cobra.Command, []string) error {
			cfg, _, err := a.loadConfig()
			if err != nil {
				return err
			}
			served, taken := map[string]bool{}, map[string]bool{}
			if projects, err := a.serveDir().Projects(); err == nil {
				for _, p := range projects {
					served[p.Root], taken[p.Key] = true, true
				}
			}
			found := serve.FindCandidates(cfg.ImportRoots, served, taken)
			if found == nil {
				found = []serve.Candidate{}
			}
			roots := cfg.ImportRoots
			if roots == nil {
				roots = []string{}
			}
			if a.jsonOut {
				return a.printJSON(map[string]any{"roots": roots, "candidates": found})
			}
			if len(roots) == 0 {
				fmt.Fprintln(a.out, "no folders named; flai serve import add <folder> names one")
				return nil
			}
			for _, r := range roots {
				fmt.Fprintf(a.out, "looked in: %s\n", r)
			}
			if len(found) == 0 {
				fmt.Fprintln(a.out, "no repositories to import there")
			}
			for _, c := range found {
				fmt.Fprintf(a.out, "  %s  %s\n", c.Name, c.Root)
			}
			return nil
		}},
	)
	return c
}

// importRoots are the folders named for import, read again at every look so
// that flai serve import add takes effect with no restart.
func (a *app) importRoots() []string {
	cfg, _, err := a.loadConfig()
	if err != nil {
		return nil
	}
	return cfg.ImportRoots
}
