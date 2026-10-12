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
// command (S-0367, TH-0413): the main checkout's uncommitted files are fixed
// before anything changes and again before the acceptance commit, and the
// story's branch was fixed by its rebase; what a fix leaves refuses the
// acceptance, naming each file, line, and error, before anything changes or
// before the merge.

// fixUncommitted runs the lint fix over what the checkout at root has not
// committed, logging at warn what it could not do, and answers those paths.
func (a *app) fixUncommitted(ctx context.Context, id, root string) []string {
	paths, err := verify.ChangedPaths(a.runner, root, "")
	if err != nil {
		a.logger().Warn("the lint fix of the acceptance's uncommitted files did not run", "component", "accept", "item", id, "err", err)
		return nil
	}
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

// refuseUnfixedUncommitted fixes the lint of what the main checkout at root
// has not committed and refuses id's acceptance for what the fix left, before
// anything else changes.
func (a *app) refuseUnfixedUncommitted(ctx context.Context, id, root string) error {
	paths := a.fixUncommitted(ctx, id, root)
	if len(paths) == 0 {
		return nil
	}
	return a.refuseUnfixed(ctx, id, root, paths, "the main checkout's uncommitted files", "in the main checkout")
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
	return a.refuseUnfixed(context.Background(), id, wt, paths, where, fixIn)
}

// refuseUnfixed runs the check of the tiers that fix over paths in the
// checkout at root and refuses id's acceptance when it fails; tiers that
// cannot be read are logged at warn and refuse nothing.
func (a *app) refuseUnfixed(ctx context.Context, id, root string, paths []string, where, fixIn string) error {
	res, err := verify.CheckoutUnfixed(ctx, root, paths)
	if err != nil {
		a.logger().Warn("the lint the fix left was not checked", "component", "accept", "item", id, "checkout", root, "err", err)
		return nil
	}
	if res.Passed {
		return nil
	}
	return unfixedError(id, where, fixIn, res)
}

// unfixedError is the refusal of id's acceptance for the findings of res's
// failing tiers in where, saying where to correct them and to accept again.
func unfixedError(id, where, fixIn string, res verify.Result) error {
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
	head := fmt.Sprintf("rule: %s cannot be accepted: %s hold lint the fix of the %s tier could not fix:", id, where, strings.Join(tiers, ", "))
	tail := fmt.Sprintf("correct them %s and run flai accept %s again", fixIn, id)
	return errors.New(strings.Join(append(append([]string{head}, lines...), tail), "\n"))
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
