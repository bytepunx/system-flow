package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"slices"
	"strings"
	"syscall"

	"github.com/spf13/cobra"

	"github.com/bytepunx/system-flow/flai/internal/buildinfo"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/verify"
	"github.com/bytepunx/system-flow/flai/internal/workflow"
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
	var suite, at string
	var maxFindings int
	c := &cobra.Command{
		Use:   "test [--suite <suite>] [path|package]... | test <tier> | test --at <point> [--suite <suite>]",
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
paths, with another name, with --all, or with --suite is refused.

--suite <suite> runs only the tiers of that suite, a tier's suite key or,
without one, its name: those of it that the paths select, or with --all
every one of it but a manual one, which still runs only by name. A tier of
the suite that a tier of another suite covers runs, since the covering tier
does not. A suite no tier is in is refused, naming the suites there are
(ADR-0312).

--at <point> runs the tiers of the suites the workflow file's verify_at
names for the point, batch, publish, or ci, over all paths, as --all runs
them: every tier of those suites, all_only ones too, but a manual one. A
point given all runs every tier but the manual ones. A workflow file that
gives none of verify, runs, and verify_at reads as the baseline's, which
gives batch and ci all and publish nothing. A point the workflow gives no
suites is refused, as is a suite it names that no tier is in. With --suite,
only that suite's tiers run, and only when the point names it or all, so
that a CI job can run one suite of the point. --at with paths, with --all,
or with a tier's name is refused, since it runs over all paths already. The
answer's first line names the point and its suites (ADR-0312 §6).

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
a file or folder, a tier's name given with others, a suite no tier is in,
a point that is not one or that the workflow gives no suites, --at with
paths, --all, a tier's name, or a suite the point does not name, tiers the
manifest does not declare validly, or a run stopped before it finished.
With --json the answer is the result, {"passed": ..., "paths": [...], "tiers": [...]}, each tier with
its state, command, duration, and findings, and a skipped tier with
covered_by; with --at it adds at, the point, and suites, the suites that
ran, or ["all"].`,
		Example: `  flai test flai/internal/manifest
  flai test flaiover/src/lib/board.ts --json
  flai test                # what this checkout changed against the main branch
  flai test --all --max 10
  flai test --suite static flai/internal/verify
  flai test --suite integration --all
  flai test gateway        # the tier named gateway, alone, for the whole checkout
  flai test --at ci        # the suites the workflow's verify_at names for ci
  flai test --at ci --suite static`,
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
			var run *atRun
			if cmd.Flags().Changed("at") {
				if run, err = atSuites(root, tiers, at, suite, args, all); err != nil {
					return &exitError{code: exitTestUnusable, msg: err.Error()}
				}
				all = true
			}
			named, err := namedTier(tiers, args, all, suite)
			if err != nil {
				return &exitError{code: exitTestUnusable, msg: err.Error()}
			}
			if err := verify.CheckSuite(tiers, suite); err != nil {
				return &exitError{code: exitTestUnusable, msg: err.Error()}
			}
			opts := verify.Options{Root: root, Tiers: tiers, All: all, Suite: suite, Git: a.runner, RunOptions: verify.RunOptions{Max: maxFindings}}
			if run != nil {
				opts.Suites = run.suites
			}
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
			if run != nil {
				res.At, res.Suites = run.point, run.ran
			}
			a.logger().Debug("tiers run", "component", "cmd", "root", root, "at", res.At, "paths", len(res.Paths), "tiers", len(res.Tiers), "passed", res.Passed)
			if a.jsonOut {
				if err := a.printJSON(res); err != nil {
					return err
				}
			} else {
				if run != nil {
					fmt.Fprintf(a.out, "at %s: suites %s\n", run.point, strings.Join(run.ran, ", "))
				}
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
	c.Flags().StringVar(&suite, "suite", "", "run only the tiers of this suite, by a tier's suite or else its name: those the paths select, or with --all every one but the manual ones")
	c.Flags().StringVar(&at, "at", "", "run the tiers of the suites the workflow's verify_at names for this point, batch, publish, or ci, over all paths, as --all runs them")
	c.Flags().IntVar(&maxFindings, "max", verify.DefaultMax, "the most findings to report across the run")
	return c
}

// atRun is what flai test --at runs: the point, the suites the run selects,
// nil for every one, and the suites it ran as the answer names them.
type atRun struct {
	point  string
	suites []string
	ran    []string
}

// atSuites is the run of the suites the workflow of the checkout at root
// names for point, narrowed to suite when it is given (ADR-0312 §6). It
// refuses a point that is none of verify_at's, paths, a tier's name, or all
// beside it, a point the workflow gives no suites, a suite the point names
// that no tier is in, and a suite the point does not name.
func atSuites(root string, tiers []verify.Tier, point, suite string, args []string, all bool) (*atRun, error) {
	point = strings.TrimSpace(point)
	points := strings.Join(workflow.VerifyAtPoints, ", ")
	if !slices.Contains(workflow.VerifyAtPoints, point) {
		return nil, fmt.Errorf("--at %s is not a point verify_at gives suites for; give one of %s (ADR-0312 §6)", point, points)
	}
	if all {
		return nil, fmt.Errorf("--at %s runs every tier of its suites over all paths, as --all does; give --at alone, not with --all", point)
	}
	if len(args) > 0 {
		if _, ok := verify.SelectTier(tiers, args[0]); ok && len(args) == 1 {
			return nil, fmt.Errorf("%s names a tier, which runs alone, and --at %s runs the tiers of the suites the workflow names for %s; give the tier alone or --at, not both", strings.TrimSpace(args[0]), point, point)
		}
		return nil, fmt.Errorf("--at %s runs over all paths, so it takes no path; give --at alone, or the paths without it", point)
	}
	repo, err := workitem.Open(root)
	if err != nil {
		return nil, fmt.Errorf("open the project at %s to read its workflow: %w", root, err)
	}
	flow := repo.Flow()
	if testWorkflow != nil {
		flow = *testWorkflow
	}
	s := flow.SuitesAt(point)
	if s == nil {
		return nil, fmt.Errorf("run the tiers at %s: the workflow names no suites for %s; give verify_at: {%s: all}, or a list of suites, in %s (ADR-0312 §6)", point, point, point, workflow.Rel(repo.Manifest.Layout["design"]))
	}
	run := &atRun{point: point, suites: s.Names, ran: s.Names}
	if s.All {
		run.suites, run.ran = nil, []string{workflow.SuitesAll}
	} else if err := verify.CheckSuites(tiers, s.Names); err != nil {
		return nil, fmt.Errorf("%w (the workflow names it for %s)", err, point)
	}
	if suite = strings.TrimSpace(suite); suite != "" {
		if !s.All && !slices.Contains(s.Names, suite) {
			return nil, fmt.Errorf("--suite %s is not among the suites the workflow names for %s, %s; give one of them, or leave --suite out", suite, point, strings.Join(s.Names, ", "))
		}
		run.ran = []string{suite}
	}
	return run, nil
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
// refuses a tier's name given with --all, with --suite, with another name,
// or with paths.
func namedTier(tiers []verify.Tier, args []string, all bool, suite string) (*verify.Selected, error) {
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
	case strings.TrimSpace(suite) != "":
		return nil, fmt.Errorf("%s names a tier, which runs alone whatever its suite, and --suite %s runs the tiers of a suite; give the tier alone, or --suite with paths or --all, not both", names[0], strings.TrimSpace(suite))
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
