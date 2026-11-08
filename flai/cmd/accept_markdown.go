package cmd

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/bytepunx/system-flow/flai/internal/verify"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// What flai accept commits passes the lint of the tiers that carry a fix
// command (S-0367, TH-0413): the main checkout's uncommitted files it
// commits, the shared set and never the rest it leaves as it is (ADR-0248),
// are fixed before anything changes, and again before each commit where flai
// wrote them since; the story's branch was fixed by its rebase. What a fix
// leaves refuses the acceptance, naming each file, line, and error, before
// anything changes or before the merge. Each set is linted once: the fix's
// own exit says what it left, so no check runs after it, unless that exit
// cannot say (S-0609).

// lintedFiles are the main checkout's files whose lint the acceptance's
// first fix passed, by their root-relative path, each with the sha256 of its
// content then (S-0609). A later fix before a commit passes over a file
// whose content is the same, since the fix leaves a file that passes as it is.
type lintedFiles map[string][sha256.Size]byte

// unchanged reports whether the file p of the checkout at root still holds
// the content its lint passed with.
func (l lintedFiles) unchanged(root, p string) bool {
	sum, ok := l[p]
	if !ok {
		return false
	}
	data, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(p)))
	return err == nil && sha256.Sum256(data) == sum
}

// newLintedFiles records the content of each of the paths of the checkout at
// root that is a file, as lint passed it.
func newLintedFiles(root string, paths []string) lintedFiles {
	l := lintedFiles{}
	for _, p := range paths {
		if data, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(p))); err == nil {
			l[p] = sha256.Sum256(data)
		}
	}
	return l
}

// uncommittedToLint are what the checkout at root has not committed but the
// paths left, which the acceptance leaves as they are (S-0402, ADR-0248),
// logging at warn when they cannot be listed.
func (a *app) uncommittedToLint(id, root string, left []string) ([]string, bool) {
	paths, err := verify.ChangedPaths(a.runner, root, "")
	if err != nil {
		a.logger().Warn("the lint fix of the acceptance's uncommitted files did not run", "component", "accept", "item", id, "err", err)
		return nil, false
	}
	return verify.WithoutLeft(paths, left), true
}

// fixUncommitted runs the lint fix before a commit over what the checkout at
// root has not committed but the paths left, passing over the files linted
// holds unchanged, and logs at warn what it could not do. With none left to
// fix it runs nothing (S-0609).
func (a *app) fixUncommitted(ctx context.Context, id, root string, left []string, linted lintedFiles) {
	paths, ok := a.uncommittedToLint(id, root, left)
	if !ok {
		return
	}
	var changed []string
	for _, p := range paths {
		if !linted.unchanged(root, p) {
			changed = append(changed, p)
		}
	}
	if len(changed) == 0 {
		return
	}
	fixed, err := verify.CheckoutFix(ctx, root, changed, a.logger())
	if err != nil {
		a.logger().Warn("the lint fix of the acceptance's uncommitted files did not run", "component", "accept", "item", id, "err", err)
		return
	}
	if len(fixed) > 0 {
		a.logger().Info("lint fixed", "component", "accept", "item", id, "paths", strings.Join(fixed, ", "))
	}
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
// the fix left, before anything else changes. The fix's own exit says what it
// left; only when it cannot say is the check run after it (S-0609). It
// answers the files whose lint passed, for the fixes before the commits.
func (a *app) refuseUnfixedUncommitted(ctx context.Context, id, root string, left []string) (lintedFiles, error) {
	paths, ok := a.uncommittedToLint(id, root, left)
	if !ok || len(paths) == 0 {
		return nil, nil
	}
	at := unfixedAt{root: root, checkout: verify.CheckoutMain, where: "the main checkout's uncommitted files", fixIn: "in the main checkout"}
	out, err := verify.CheckoutFixRemaining(ctx, root, paths, a.logger())
	if err != nil {
		a.logger().Warn("the lint fix of the acceptance's uncommitted files did not run", "component", "accept", "item", id, "err", err)
		return nil, a.refuseUnfixed(ctx, id, paths, at)
	}
	if len(out.Changed) > 0 {
		a.logger().Info("lint fixed", "component", "accept", "item", id, "paths", strings.Join(out.Changed, ", "))
	}
	if err := a.refuseRemaining(ctx, id, paths, at, out); err != nil {
		return nil, err
	}
	return newLintedFiles(root, paths), nil
}

// refuseUnfixedBranch refuses story id's acceptance for the lint the fix
// after its rebase onto base left in the files its branch changes, before
// the merge: from what that fix's exit said it left, fixed, or, when it did
// not run or cannot say, by the check (S-0609).
func (a *app) refuseUnfixedBranch(repo *workitem.Repo, id, base string, fixed *verify.FixOutcome) error {
	wt := repo.WorktreePath(id)
	where := fmt.Sprintf("the files %s changes", storyBranch(id))
	fixIn := fmt.Sprintf("in %s, commit them on %s,", relPath(repo.MainRoot, wt), storyBranch(id))
	at := unfixedAt{root: wt, checkout: verify.CheckoutStory, story: id, where: where, fixIn: fixIn}
	if fixed != nil && fixed.Known() {
		return a.refuseRemaining(context.Background(), id, nil, at, *fixed)
	}
	paths, err := verify.ChangedPaths(a.runner, wt, base)
	if err != nil {
		a.logger().Warn("the lint of the story branch's files was not checked", "component", "accept", "item", id, "err", err)
		return nil
	}
	return a.refuseUnfixed(context.Background(), id, paths, at)
}

// refuseRemaining refuses id's acceptance for what the fix out says it left
// in at, or, when it cannot say, checks paths as refuseUnfixed does.
func (a *app) refuseRemaining(ctx context.Context, id string, paths []string, at unfixedAt, out verify.FixOutcome) error {
	if !out.Known() {
		return a.refuseUnfixed(ctx, id, paths, at)
	}
	if out.Remains.Passed {
		return nil
	}
	return unfixedError(id, at, out.Remains)
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
