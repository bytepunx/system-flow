package cmd

import (
	"context"
	"errors"
	"fmt"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
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
// committed but the paths left, which the acceptance leaves as they are
// (S-0402), logging at warn what it could not do, and answers those paths.
func (a *app) fixUncommitted(ctx context.Context, id, root string, left []string) []string {
	paths, err := verify.ChangedPaths(a.runner, root, "")
	if err != nil {
		a.logger().Warn("the lint fix of the acceptance's uncommitted files did not run", "component", "accept", "item", id, "err", err)
		return nil
	}
	paths = withoutLeft(paths, left)
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

// The checkouts a finding of unfixed lint is to be corrected in (S-0378): the
// main checkout's uncommitted files, or the story branch's in its worktree.
const (
	checkoutMain  = "main"
	checkoutStory = "story"
)

// unfixedAt is where the lint an acceptance refuses for is to be corrected:
// the checkout the check ran in, which of the two it is, and how the
// refusal's text names it.
type unfixedAt struct {
	root     string // the checkout on disk
	checkout string // checkoutMain or checkoutStory
	story    string // the story whose branch it is, for checkoutStory
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
	at := unfixedAt{root: root, checkout: checkoutMain, where: "the main checkout's uncommitted files", fixIn: "in the main checkout"}
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
	at := unfixedAt{root: wt, checkout: checkoutStory, story: id, where: where, fixIn: fixIn}
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
	Reason   string           `json:"reason"`
	Findings []unfixedFinding `json:"findings"`
}

// Error is the refusal's text as a rule refusal.
func (e *unfixedRefusal) Error() string { return "rule: " + e.Reason }

// unfixedFinding is one finding of unfixed lint and where to correct it.
type unfixedFinding struct {
	Tier string `json:"tier"`
	// Checkout is checkoutMain or checkoutStory; Story names the story whose
	// branch holds the file, for checkoutStory.
	Checkout string `json:"checkout"`
	Story    string `json:"story,omitempty"`
	// Path is relative to the repository's root.
	Path   string `json:"path"`
	Line   int    `json:"line,omitempty"`
	Column int    `json:"column,omitempty"`
	// Level is the tool's severity when it gives one: error or warning.
	Level string `json:"level,omitempty"`
	// Rule is the rule's ID, MD009, or the check's or linter's name; Alias
	// is the rule's first name, no-trailing-spaces, when the tool gives one.
	Rule    string `json:"rule,omitempty"`
	Alias   string `json:"alias,omitempty"`
	Message string `json:"message"`
}

// unfixedError is the refusal of id's acceptance for the findings of res's
// failing tiers in at, saying where to correct them and to accept again.
func unfixedError(id string, at unfixedAt, res verify.Result) *unfixedRefusal {
	var lines, tiers []string
	findings := []unfixedFinding{}
	for _, t := range res.Tiers {
		if t.State != verify.Failed {
			continue
		}
		tiers = append(tiers, t.Name)
		for _, f := range t.Findings {
			lines = append(lines, "  "+findingLine(t.Name, f))
			findings = append(findings, unfixedFindings(at, t, f)...)
		}
		if t.Omitted > 0 {
			lines = append(lines, "  and "+plural(t.Omitted, "finding")+" more")
		}
	}
	head := fmt.Sprintf("%s cannot be accepted: %s hold lint the fix of the %s tier could not fix:", id, at.where, strings.Join(tiers, ", "))
	tail := fmt.Sprintf("correct them %s and run flai accept %s again", at.fixIn, id)
	return &unfixedRefusal{ID: id, Reason: strings.Join(append(append([]string{head}, lines...), tail), "\n"), Findings: findings}
}

// markdownlintLine is a line of markdownlint-cli2's default output:
// path:line[:column] [error|warning] MDnnn[/alias...] message.
var markdownlintLine = regexp.MustCompile(`^(.+?):(\d+)(?::(\d+))?\s+(?:(error|warning)\s+)?(MD\d{3})((?:/[A-Za-z0-9_-]+)*)\s+(.+)$`)

// unfixedFindings are the findings f of tier t gives as data: f itself when
// it names a file, or each markdownlint-cli2 line of a plain tier's output,
// none when it holds no such line.
func unfixedFindings(at unfixedAt, t verify.TierResult, f verify.Finding) []unfixedFinding {
	base := unfixedFinding{Tier: t.Name, Checkout: at.checkout, Story: at.story}
	if f.Path != "" {
		u := base
		u.Path, u.Line, u.Message = f.Path, f.Line, strings.TrimSpace(f.Message)
		if f.Name != t.Name {
			u.Rule = f.Name
		}
		return []unfixedFinding{u}
	}
	var got []unfixedFinding
	for _, line := range strings.Split(f.Message, "\n") {
		m := markdownlintLine.FindStringSubmatch(strings.TrimSpace(line))
		if m == nil {
			continue
		}
		u := base
		u.Path = checkoutRel(at.root, t.Dir, m[1])
		u.Line, _ = strconv.Atoi(m[2])
		u.Column, _ = strconv.Atoi(m[3])
		u.Level, u.Rule, u.Message = m[4], m[5], m[7]
		if aliases := strings.TrimPrefix(m[6], "/"); aliases != "" {
			u.Alias, _, _ = strings.Cut(aliases, "/")
		}
		got = append(got, u)
	}
	return got
}

// checkoutRel is a path a tool run in dir of the checkout at root printed,
// relative to root: an absolute one below root made relative, and a relative
// one joined to dir.
func checkoutRel(root, dir, p string) string {
	if filepath.IsAbs(p) {
		if r, err := filepath.Rel(root, p); err == nil && r != ".." && !strings.HasPrefix(r, ".."+string(filepath.Separator)) {
			return filepath.ToSlash(r)
		}
		return filepath.ToSlash(p)
	}
	return path.Join(filepath.ToSlash(dir), filepath.ToSlash(p))
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
