package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/bytepunx/system-flow/flai/internal/buildinfo"
	"github.com/bytepunx/system-flow/flai/internal/verify"
	"github.com/bytepunx/system-flow/flai/internal/workflow"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// exitReleaseSuiteFailed is flai release --pending's exit status when a tier
// of the suites the workflow file's verify_at.publish names fails, or they
// cannot be run, at the head being published (ADR-0312 §6): 5, apart from 2,
// nothing pending, 3, a remote that moved, and 4, the push host action off,
// so that the publish formula's step stops there with nothing published.
const exitReleaseSuiteFailed = 5

// publishSuitesRun is what verify_at.publish ran before a publish: the
// suites the file names, all for every suite, the head they ran at, the
// folder holding each tier's whole output, and each tier's outcome with its
// suite.
type publishSuitesRun struct {
	Suites []string      `json:"suites"`
	Head   string        `json:"head"`
	Log    string        `json:"log"`
	Passed bool          `json:"passed"`
	Tiers  []publishTier `json:"tiers"`
}

// publishTier is one tier's outcome in a publishSuitesRun, with its suite.
type publishTier struct {
	Suite string `json:"suite"`
	verify.TierResult
}

// publishSuitesPreview is what verify_at.publish would run, as flai release
// --pending --dry-run names it: the suites the file names and the tiers they
// select, or why the tiers cannot be read.
type publishSuitesPreview struct {
	Suites   []string    `json:"suites"`
	WouldRun []suiteTier `json:"would_run"`
	Error    string      `json:"error,omitempty"`
}

// suiteTier is a tier by name with its suite.
type suiteTier struct {
	Suite string `json:"suite"`
	Tier  string `json:"tier"`
}

// publishSuites is the suites the workflow file's verify_at.publish names,
// nil when it names none, as the baseline does, and as a file that gives
// none of ADR-0312's keys reads (UsesSuiteKeys, D-0669): then no tier runs.
func publishSuites(repo *workitem.Repo) *workflow.Suites {
	return repo.Flow().VerifyAt.Publish
}

// suiteNames is the suites s names as its file gives them: all, or the list.
func suiteNames(s *workflow.Suites) []string {
	if s.All {
		return []string{workflow.SuitesAll}
	}
	return s.Names
}

// selectPublishTiers is the tiers each suite s names selects at root, in the
// order the file names the suites, one group per suite, or one group of
// every suite for all: every tier of it, all_only ones too and manual ones
// not, as flai test --suite <suite> --all selects them. A suite no tier is
// in is refused before anything runs.
func selectPublishTiers(root string, tiers []verify.Tier, s *workflow.Suites) ([][]verify.Selected, error) {
	fsys := os.DirFS(root)
	if s.All {
		return [][]verify.Selected{verify.SelectSuite(fsys, tiers, nil, true, "")}, nil
	}
	groups := make([][]verify.Selected, 0, len(s.Names))
	for _, name := range s.Names {
		if err := verify.CheckSuite(tiers, name); err != nil {
			return nil, err
		}
		groups = append(groups, verify.SelectSuite(fsys, tiers, nil, true, name))
	}
	return groups, nil
}

// previewPublishSuites is what verify_at.publish would run at root before a
// publish, nil when it names no suite; it runs nothing.
func (a *app) previewPublishSuites(repo *workitem.Repo) *publishSuitesPreview {
	s := publishSuites(repo)
	if s == nil {
		return nil
	}
	p := &publishSuitesPreview{Suites: suiteNames(s), WouldRun: []suiteTier{}}
	tiers, err := verify.CheckoutTiersFor(repo.Root, a.runner, buildinfo.Version)
	if err != nil {
		p.Error = err.Error()
		return p
	}
	groups, err := selectPublishTiers(repo.Root, tiers, s)
	if err != nil {
		p.Error = err.Error()
		return p
	}
	for _, g := range groups {
		for _, sel := range g {
			if sel.CoveredBy == "" {
				p.WouldRun = append(p.WouldRun, suiteTier{Suite: sel.Tier.SuiteOrName(), Tier: sel.Tier.Name})
			}
		}
	}
	return p
}

// printPublishSuitesPreview says what verify_at.publish would run.
func (a *app) printPublishSuitesPreview(p *publishSuitesPreview) {
	if p == nil {
		return
	}
	if p.Error != "" {
		fmt.Fprintf(a.out, "verify_at.publish would run the suites %s before tagging, and would refuse the publish: %s\n", strings.Join(p.Suites, ", "), p.Error)
		return
	}
	names := make([]string, len(p.WouldRun))
	for i, t := range p.WouldRun {
		names[i] = fmt.Sprintf("%s (%s)", t.Tier, t.Suite)
	}
	tiers := "no tier"
	if len(names) > 0 {
		tiers = "the tiers " + strings.Join(names, ", ")
	}
	fmt.Fprintf(a.out, "verify_at.publish would run the suites %s before tagging: %s\n", strings.Join(p.Suites, ", "), tiers)
}

// runPublishSuites runs the suites verify_at.publish names at the head of
// the checkout at repo.Root, before anything of a publish is applied
// (ADR-0312 §6), each tier's whole output kept under
// .flai-cache/publish/<UTC time>/. It answers nil when the file names no
// suite, and the run otherwise. A tier that fails, tiers that cannot be read,
// a suite no tier is in, or a run stopped is an exitError of
// exitReleaseSuiteFailed whose message names the suite, the tier, and the
// log, with the run so far.
func (a *app) runPublishSuites(repo *workitem.Repo) (*publishSuitesRun, error) {
	s := publishSuites(repo)
	if s == nil {
		return nil, nil
	}
	root := repo.Root
	file := repo.Flow().Path
	run := &publishSuitesRun{Suites: suiteNames(s), Tiers: []publishTier{}}
	head, err := a.runner.Run(root, "git", "rev-parse", "HEAD")
	if err != nil {
		return run, &exitError{code: exitReleaseSuiteFailed, msg: fmt.Sprintf("publish refused: read the head of %s to run verify_at.publish's suites %s at: %s; nothing was applied, committed, tagged, or pushed", root, strings.Join(run.Suites, ", "), firstLine(err.Error()))}
	}
	run.Head = strings.TrimSpace(head)
	refuse := func(why string) error {
		return &exitError{code: exitReleaseSuiteFailed, msg: fmt.Sprintf("publish refused: verify_at.publish in %s names the suites %s, which cannot be run at %s: %s; nothing was applied, committed, tagged, or pushed", file, strings.Join(run.Suites, ", "), short(run.Head), why)}
	}
	tiers, err := verify.CheckoutTiersFor(root, a.runner, buildinfo.Version)
	if err != nil {
		return run, refuse(err.Error())
	}
	groups, err := selectPublishTiers(root, tiers, s)
	if err != nil {
		return run, refuse(err.Error())
	}
	run.Log = filepath.Join(root, ".flai-cache", "publish", a.now().UTC().Format("20060102T150405Z"))
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	run.Passed = true
	for _, g := range groups {
		res, err := verify.Run(ctx, root, g, verify.RunOptions{LogDir: run.Log})
		for i, tr := range res.Tiers {
			run.Tiers = append(run.Tiers, publishTier{Suite: g[i].Tier.SuiteOrName(), TierResult: tr})
		}
		a.printSuiteResult(g, res)
		if err != nil {
			run.Passed = false
			return run, &exitError{code: exitReleaseSuiteFailed, msg: fmt.Sprintf("publish refused: the run of verify_at.publish's suites %s at %s was stopped before it finished: %v; nothing was applied, committed, tagged, or pushed, and what ran is in %s", strings.Join(run.Suites, ", "), short(run.Head), err, run.Log)}
		}
		if !res.Passed {
			run.Passed = false
			return run, a.suiteFailed(run, file)
		}
	}
	if !a.jsonOut {
		fmt.Fprintf(a.out, "verify_at.publish: the suites %s passed at %s; their whole output is in %s\n", strings.Join(run.Suites, ", "), short(run.Head), run.Log)
	}
	return run, nil
}

// printSuiteResult prints, in text, a group of verify_at.publish's tiers
// and their outcome, as flai test prints them.
func (a *app) printSuiteResult(g []verify.Selected, res verify.Result) {
	if a.jsonOut || len(g) == 0 {
		return
	}
	fmt.Fprintf(a.out, "verify_at.publish: suite %s\n", g[0].Tier.SuiteOrName())
	fmt.Fprint(a.out, res.Text())
}

// suiteFailed is the refusal of a publish whose run's tier failed, naming
// the tier, its suite, the head, and the log of its whole output.
func (a *app) suiteFailed(run *publishSuitesRun, file string) error {
	tier, suite, log := "", "", run.Log
	for _, t := range run.Tiers {
		if t.State == verify.Failed {
			tier, suite = t.Name, t.Suite
			if t.Log != "" {
				log = t.Log
			}
			break
		}
	}
	return &exitError{code: exitReleaseSuiteFailed, msg: fmt.Sprintf("publish refused: tier %s of suite %s failed at %s, which verify_at.publish in %s runs before a release, so nothing was applied, committed, tagged, or pushed; its whole output is in %s: fix it on the main branch, then run flai release --pending again", tier, suite, short(run.Head), file, log)}
}
