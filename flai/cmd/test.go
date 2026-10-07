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
		Use:   "test [path|package]...",
		Short: "Run the test and lint tiers for paths or packages and answer pass or the first findings",
		Long: `Run the project's test and lint tiers, the tests list in system-flow.yaml, for
the files and folders given, in the checkout of the working directory: a
story's worktree or the main checkout. A Go package is given as its folder.
With no argument the files are those the checkout changed against the main
branch, committed or not. --all runs every tier, all_only ones too, for the
whole checkout.

The tiers whose paths select a file run in the manifest's order, cheapest
first, and the run stops at the first that fails. The answer is pass, or the
failing tier's first findings, never its whole output: each as path:line
name: message, at most --max across the run, with how many were left out.
A manifest without tests runs scripts/test.sh as its one tier. Each tier runs
with FLAI_ROLE=verify whatever role ran flai test, so a test that runs flai's
writes is not refused as the orchestrator (S-0311).

The exit status is 0 on pass, 1 when a tier fails, and 2 when flai test could
not answer: an argument outside the checkout, tiers the manifest does not
declare validly, or a run stopped before it finished. With --json the answer
is the result, {"passed": ..., "paths": [...], "tiers": [...]}, each tier with
its state, command, duration, and findings.`,
		Example: `  flai test flai/internal/manifest
  flai test flaiover/src/lib/board.ts --json
  flai test                # what this checkout changed against the main branch
  flai test --all --max 10`,
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
			rel, err := checkoutArgs(root, cwd, args)
			if err != nil {
				return &exitError{code: exitTestUnusable, msg: err.Error()}
			}
			tiers, err := verify.CheckoutTiers(root)
			if err != nil {
				return &exitError{code: exitTestUnusable, msg: err.Error()}
			}
			opts := verify.Options{Root: root, Tiers: tiers, Args: rel, All: all, Git: a.runner, RunOptions: verify.RunOptions{Max: maxFindings}}
			if !all && len(rel) == 0 {
				if opts.Base, err = a.checkoutBase(root); err != nil {
					return &exitError{code: exitTestUnusable, msg: err.Error()}
				}
			}
			ctx, stop := signal.NotifyContext(cmd.Context(), os.Interrupt, syscall.SIGTERM)
			defer stop()
			res, err := verify.Test(ctx, opts)
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
	c.Flags().BoolVar(&all, "all", false, "run every tier for the whole checkout, all_only ones too, whatever the arguments")
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
			return nil, fmt.Errorf("%s is not a file or folder; name one in the checkout at %s, relative to the working directory", arg, root)
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
