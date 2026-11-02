package cmd

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/bytepunx/system-flow/flai/internal/verify"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// What flai accept commits passes the lint of the tiers that carry a fix
// command (S-0367, TH-0413): the main checkout's uncommitted files it
// commits, the shared set and never the rest it leaves as it is (ADR-0248),
// are fixed before anything changes and again before each commit, and the
// story's branch was fixed by its rebase; what a fix leaves refuses the
// acceptance, naming each file, line, and error, before anything changes or
// before the merge.

// fixUncommitted runs the lint fix over what the checkout at root has not
// committed but the paths left, which the acceptance leaves as they are
// (S-0402, ADR-0248), logging at warn what it could not do, and answers
// those paths.
func (a *app) fixUncommitted(ctx context.Context, id, root string, left []string) []string {
	paths, err := verify.ChangedPaths(a.runner, root, "")
	if err != nil {
		a.logger().Warn("the lint fix of the acceptance's uncommitted files did not run", "component", "accept", "item", id, "err", err)
		return nil
	}
	paths = verify.WithoutLeft(paths, left)
	fixed, err := verify.CheckoutFix(ctx, root, paths, a.logger())
	if err != nil {
		a.logger().Warn("the lint fix of the acceptance's uncommitted files did not run", "component", "accept", "item", id, "err", err)
		return paths
	}
	if len(fixed) > 0 {
		a.logger().Info("lint fixed", "component", "accept", "item", id, "paths", strings.Join(fixed, ", "))
	}
	return paths
}

// unfixedAt is where the lint an acceptance refuses for is to be corrected:
// the checkout the check ran in, which of the two it is, and how the
// refusal's text names it.
type unfixedAt struct {
	root     string // the checkout on disk
	checkout string // verify.CheckoutMain or verify.CheckoutStory
	story    string // the story whose branch it is, for verify.CheckoutStory
	where    string // what holds the lint, as the refusal says it
	fixIn    string // where to correct it, as the refusal says it
}

// refuseUnfixedUncommitted fixes the lint of what the main checkout at root
// has not committed, but the paths left, and refuses id's acceptance for what
// the fix left, before anything else changes.
func (a *app) refuseUnfixedUncommitted(ctx context.Context, id, root string, left []string) error {
	paths := a.fixUncommitted(ctx, id, root, left)
	if len(paths) == 0 {
		return nil
	}
	at := unfixedAt{root: root, checkout: verify.CheckoutMain, where: "the main checkout's uncommitted files", fixIn: "in the main checkout"}
	return a.refuseUnfixed(ctx, id, paths, at)
}

// refuseUnfixedBranch refuses story id's acceptance for the lint the fix
// after its rebase onto base left in the files its branch changes, before
// the merge.
func (a *app) refuseUnfixedBranch(repo *workitem.Repo, id, base string) error {
	wt := repo.WorktreePath(id)
	paths, err := verify.ChangedPaths(a.runner, wt, base)
	if err != nil {
		a.logger().Warn("the lint of the story branch's files was not checked", "component", "accept", "item", id, "err", err)
		return nil
	}
	where := fmt.Sprintf("the files %s changes", storyBranch(id))
	fixIn := fmt.Sprintf("in %s, commit them on %s,", relPath(repo.MainRoot, wt), storyBranch(id))
	at := unfixedAt{root: wt, checkout: verify.CheckoutStory, story: id, where: where, fixIn: fixIn}
	return a.refuseUnfixed(context.Background(), id, paths, at)
}

// refuseUnfixed runs the check of the tiers that fix over paths in the
// checkout at.root and refuses id's acceptance when it fails; tiers that
// cannot be read are logged at warn and refuse nothing.
func (a *app) refuseUnfixed(ctx context.Context, id string, paths []string, at unfixedAt) error {
	res, err := verify.CheckoutUnfixed(ctx, at.root, paths)
	if err != nil {
		a.logger().Warn("the lint the fix left was not checked", "component", "accept", "item", id, "checkout", at.root, "err", err)
		return nil
	}
	if res.Passed {
		return nil
	}
	return unfixedError(id, at, res)
}

// unfixedRefusal is the refusal of an acceptance for lint the fix left
// (ADR-0140): its text, and each finding as data for flai accept --json,
// which prints it as flai doc save prints its refusal (S-0378).
type unfixedRefusal struct {
	ID string `json:"id"`
	// Reason is the refusal's text, without its rule: prefix.
	Reason   string               `json:"reason"`
	Findings []verify.FileFinding `json:"findings"`
}

// Error is the refusal's text as a rule refusal.
func (e *unfixedRefusal) Error() string { return "rule: " + e.Reason }

// unfixedError is the refusal of id's acceptance for the findings of res's
// failing tiers in at, saying where to correct them and to accept again.
func unfixedError(id string, at unfixedAt, res verify.Result) *unfixedRefusal {
	var lines, tiers []string
	for _, t := range res.Tiers {
		if t.State != verify.Failed {
			continue
		}
		tiers = append(tiers, t.Name)
		for _, f := range t.Findings {
			lines = append(lines, "  "+findingLine(t.Name, f))
		}
		if t.Omitted > 0 {
			lines = append(lines, "  and "+plural(t.Omitted, "finding")+" more")
		}
	}
	head := fmt.Sprintf("%s cannot be accepted: %s hold lint the fix of the %s tier could not fix:", id, at.where, strings.Join(tiers, ", "))
	tail := fmt.Sprintf("correct them %s and run flai accept %s again", at.fixIn, id)
	findings := verify.FileFindings(at.root, at.checkout, at.story, res)
	return &unfixedRefusal{ID: id, Reason: strings.Join(append(append([]string{head}, lines...), tail), "\n"), Findings: findings}
}

// acceptRefused is err, or under --json, when err is a refusal for lint the
// fix left, that refusal printed as {"refused": {...}} with exit 4, as flai
// doc save prints its own, so the dashboard can open each file (S-0378).
func (a *app) acceptRefused(err error) error {
	var r *unfixedRefusal
	if !a.jsonOut || !errors.As(err, &r) {
		return err
	}
	_ = a.printJSON(map[string]any{"refused": r})
	return &exitError{code: exitDocRefused, msg: err.Error()}
}

// findingLine names one finding of tier as path:line: message, with the
// tier and the finding's own name, or as the tier and its message when the
// finding names no file.
func findingLine(tier string, f verify.Finding) string {
	name := tier
	if f.Name != "" && f.Name != tier {
		name += ", " + f.Name
	}
	msg := strings.ReplaceAll(strings.TrimSpace(f.Message), "\n", "\n    ")
	if f.Path == "" {
		return name + ": " + msg
	}
	loc := f.Path
	if f.Line > 0 {
		loc = fmt.Sprintf("%s:%d", f.Path, f.Line)
	}
	return fmt.Sprintf("%s: %s (%s)", loc, msg, name)
}
