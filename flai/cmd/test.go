package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/spf13/cobra"

	"github.com/bytepunx/system-flow/flai/internal/buildinfo"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/verify"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// exitTestUnusable is flai test's exit status when it could not answer: a
// usage error, tiers the manifest does not declare validly, or a checkout it
// cannot read, kept apart from a failing tier's 1 (S-0273).
const exitTestUnusable = 2

// newTestCmd runs the project's test and lint tiers for paths or packages
// and answers pass or the first findings (S-0273).
func newTestCmd(a *app) *cobra.Command {
	var all bool
	var maxFindings int
	c := &cobra.Command{
		Use:   "test [path|package]... | test <tier>",
		Short: "Run the test and lint tiers for paths or packages, or one tier by name, and answer pass or the first findings",
		Long: `Run the project's test and lint tiers, the tests list in system-flow.yaml, for
the files and folders given, in the checkout of the working directory: a
story's worktree or the main checkout. A Go package is given as its folder.
With no argument the files are those the checkout changed against the main
branch, committed or not. --all runs every tier but the manual ones,
all_only ones too, for the whole checkout.

A lone argument that names a tier runs that tier alone for the whole
checkout, as --all runs it, with its all_command when it has one: all_only,
manual, or neither, and even where a file or folder has the same name. A
manual tier runs only so: no path selects it, and neither --all, flai
verify, nor the merge queue's batch runs it (S-0548). A tier's name with
paths, with another name, or with --all is refused.

The tiers whose paths select a file run in the manifest's order, cheapest
first, and the run stops at the first that fails. The answer is pass, or the
failing tier's first findings, never its whole output: each as path:line
name: message, at most --max across the run, with how many were left out.
A selected tier that another selected tier covers, by the covers key, is not
run: it is listed as "skipped <tier>: covered by <tier>" (S-0342).
A manifest without tests runs scripts/test.sh as its one tier. Each tier runs
with FLAI_ROLE=verify whatever role ran flai test, so a test that runs flai's
writes is not refused as the orchestrator (S-0311).

The exit status is 0 on pass, 1 when a tier fails, and 2 when flai test could
not answer: an argument outside the checkout, one that is neither a tier nor
a file or folder, a tier's name given with others, tiers the manifest does
not declare validly, or a run stopped before it finished. With --json the answer
is the result, {"passed": ..., "paths": [...], "tiers": [...]}, each tier with
its state, command, duration, and findings, and a skipped tier with
covered_by.`,
		Example: `  flai test flai/internal/manifest
  flai test flaiover/src/lib/board.ts --json
  flai test                # what this checkout changed against the main branch
  flai test --all --max 10
  flai test gateway        # the tier named gateway, alone, for the whole checkout`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if maxFindings < 1 {
				return &exitError{code: exitTestUnusable, msg: fmt.Sprintf("--max %d is not a count of findings; give 1 or more", maxFindings)}
			}
			cwd, err := a.workingDir()
			if err != nil {
				return &exitError{code: exitTestUnusable, msg: "read the working directory: " + err.Error()}
			}
			root, err := a.checkoutRoot(cwd)
			if err != nil {
				return &exitError{code: exitTestUnusable, msg: err.Error()}
			}
			tiers, err := verify.CheckoutTiersFor(root, a.runner, buildinfo.Version)
			if err != nil {
				return &exitError{code: exitTestUnusable, msg: err.Error()}
			}
			named, err := namedTier(tiers, args, all)
			if err != nil {
				return &exitError{code: exitTestUnusable, msg: err.Error()}
			}
			opts := verify.Options{Root: root, Tiers: tiers, All: all, Git: a.runner, RunOptions: verify.RunOptions{Max: maxFindings}}
			if named == nil {
				if opts.Args, err = checkoutArgs(root, cwd, args); err != nil {
					var missing missingArgError
					if len(args) == 1 && errors.As(err, &missing) {
						err = fmt.Errorf("%s is neither a tier nor a file or folder; name a tier, %s, or a file or folder in the checkout at %s, relative to the working directory", missing.arg, tierNameList(tiers), root)
					}
					return &exitError{code: exitTestUnusable, msg: err.Error()}
				}
				if !all && len(opts.Args) == 0 {
					if opts.Base, err = a.checkoutBase(root); err != nil {
						return &exitError{code: exitTestUnusable, msg: err.Error()}
					}
				}
			}
			ctx, stop := signal.NotifyContext(cmd.Context(), os.Interrupt, syscall.SIGTERM)
			defer stop()
			var res verify.Result
			if named != nil {
				res, err = verify.Run(ctx, root, []verify.Selected{*named}, opts.RunOptions)
			} else {
				res, err = verify.Test(ctx, opts)
			}
			stopped := err != nil && ctx.Err() != nil
			if err != nil && !stopped {
				return &exitError{code: exitTestUnusable, msg: err.Error()}
			}
			a.logger().Debug("tiers run", "component", "cmd", "root", root, "paths", len(res.Paths), "tiers", len(res.Tiers), "passed", res.Passed)
			if a.jsonOut {
				if err := a.printJSON(res); err != nil {
					return err
				}
			} else {
				fmt.Fprint(a.out, res.Text())
			}
			switch {
			case stopped:
				return &exitError{code: exitTestUnusable, msg: "the run was stopped before it finished: " + err.Error()}
			case !res.Passed:
				return &exitError{code: 1}
			}
			return nil
		},
	}
	c.SetFlagErrorFunc(func(_ *cobra.Command, err error) error {
		return &exitError{code: exitTestUnusable, msg: err.Error()}
	})
	c.Flags().BoolVar(&all, "all", false, "run every tier for the whole checkout, all_only ones too and manual ones not, whatever the paths")
	c.Flags().IntVar(&maxFindings, "max", verify.DefaultMax, "the most findings to report across the run")
	return c
}

// checkoutRoot is the root of the checkout holding dir: git's top level, or,
// outside git, the folder of the project's system-flow.yaml.
func (a *app) checkoutRoot(dir string) (string, error) {
	if out, err := a.runner.Run(dir, "git", "rev-parse", "--show-toplevel"); err == nil && strings.TrimSpace(out) != "" {
		return filepath.Clean(strings.TrimSpace(out)), nil
	}
	path, err := manifest.Find(dir)
	if err != nil {
		return "", fmt.Errorf("%s is in no checkout of a project: %w; run flai test in a story's worktree or the main checkout", dir, err)
	}
	return filepath.Dir(path), nil
}

// checkoutBase is the main branch, as the main checkout of root has it
// checked out, that the checkout's changes are counted against.
func (a *app) checkoutBase(root string) (string, error) {
	repo, err := workitem.Open(root)
	if err != nil {
		return "", fmt.Errorf("open the project at %s: %w", root, err)
	}
	base, err := a.mainBranch(mainRootOf(repo))
	if err != nil {
		return "", fmt.Errorf("find the main branch the changes are counted against, in %s: %w; name the paths to test instead, or --all", mainRootOf(repo), err)
	}
	return base, nil
}

// namedTier is the tier the arguments name, selected alone for the whole
// checkout, when the one argument is a tier's name, trimmed, even where a
// path has the same name; nil when no argument names a tier (S-0548). It
// refuses a tier's name given with --all, with another name, or with paths.
func namedTier(tiers []verify.Tier, args []string, all bool) (*verify.Selected, error) {
	var names []string
	var first verify.Selected
	for _, arg := range args {
		if s, ok := verify.SelectTier(tiers, arg); ok {
			if len(names) == 0 {
				first = s
			}
			names = append(names, strings.TrimSpace(arg))
		}
	}
	switch {
	case len(names) == 0:
		return nil, nil
	case all:
		return nil, fmt.Errorf("%s names a tier, and --all runs every tier but the manual ones; give the tier alone or --all, not both", names[0])
	case len(names) > 1:
		return nil, fmt.Errorf("%s each name a tier, and flai test runs one tier by name; give one of them alone", strings.Join(names, ", "))
	case len(args) > 1:
		return nil, fmt.Errorf("%s names a tier, which runs alone for the whole checkout; give it with no path, or the paths without it", names[0])
	}
	return &first, nil
}

// tierNameList names the tiers in list order, for a refusal that asks for
// one, or says the manifest has none.
func tierNameList(tiers []verify.Tier) string {
	if len(tiers) == 0 {
		return "of which the manifest declares none"
	}
	names := make([]string, len(tiers))
	for i, t := range tiers {
		names[i] = strings.TrimSpace(t.Name)
	}
	return "one of " + strings.Join(names, ", ")
}

// missingArgError is an argument to flai test that names no file or folder.
type missingArgError struct{ arg, root string }

func (e missingArgError) Error() string {
	return fmt.Sprintf("%s is not a file or folder; name one in the checkout at %s, relative to the working directory", e.arg, e.root)
}

// checkoutArgs are the arguments, files or folders relative to cwd or
// absolute, as paths relative to root with forward slashes. It refuses one
// that is not there or lies outside the checkout.
func checkoutArgs(root, cwd string, args []string) ([]string, error) {
	realRoot, err := filepath.EvalSymlinks(root)
	if err != nil {
		return nil, fmt.Errorf("read the checkout at %s: %w", root, err)
	}
	out := make([]string, 0, len(args))
	for _, arg := range args {
		p := arg
		if !filepath.IsAbs(p) {
			p = filepath.Join(cwd, p)
		}
		real, err := filepath.EvalSymlinks(p)
		if errors.Is(err, os.ErrNotExist) {
			return nil, missingArgError{arg: arg, root: root}
		}
		if err != nil {
			return nil, fmt.Errorf("read %s: %w", arg, err)
		}
		rel, err := filepath.Rel(realRoot, real)
		if err != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
			return nil, fmt.Errorf("%s is outside the checkout at %s; name a file or folder inside it", arg, root)
		}
		out = append(out, filepath.ToSlash(rel))
	}
	return out, nil
}
