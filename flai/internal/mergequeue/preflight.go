package mergequeue

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/buildinfo"
	"github.com/bytepunx/system-flow/flai/internal/check"
	"github.com/bytepunx/system-flow/flai/internal/execx"
	"github.com/bytepunx/system-flow/flai/internal/release"
	"github.com/bytepunx/system-flow/flai/internal/storygit"
	"github.com/bytepunx/system-flow/flai/internal/verify"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// The landing's preflight (ADR-0234): before it builds a batch, the merge
// queue runs on the main checkout what flai accept --batch would refuse a
// landing for there, the lint of the tiers that fix over the uncommitted
// files the landing commits, and the error-level findings of flai check, and
// holds itself on what it finds until a look finds main clean. It changes no
// file of the main checkout: the operator's uncommitted files are theirs.

// PreflightWorktreeName is the folder, under the main checkout's
// .flai-cache/worktrees, that the preflight confirms a lint finding in: a
// scratch worktree of the main branch's head, removed once it has looked.
const PreflightWorktreeName = "merge-queue-preflight"

// PreflightCheckTier is the Tier of a preflight finding that flai check made
// rather than a test tier.
const PreflightCheckTier = "flai-check"

// PreflightReport is what the landing's preflight found: the findings that
// hold the queue, and why the running flai may not judge them.
type PreflightReport struct {
	// Findings hold the queue; none when main is clean.
	Findings []verify.FileFinding
	// Behind is why the findings may be the running flai's rather than the
	// files' (ADR-0311); nil when there is no such doubt.
	Behind *FlaiBehind
}

// Preflighter runs the landing's preflight on repo's main checkout and
// answers its report, with no findings when main is clean, or an error when
// it could not run.
type Preflighter func(ctx context.Context, repo *workitem.Repo) (PreflightReport, error)

// MainPreflight is the Preflighter flai serve runs: Preflight with git run
// through r, flai check judged at now by this flai, and what it could not
// clean up or settle logged to log.
func MainPreflight(r execx.Runner, now func() time.Time, log *slog.Logger) Preflighter {
	return func(ctx context.Context, repo *workitem.Repo) (PreflightReport, error) {
		return Preflight(ctx, r, repo, now(), buildinfo.Get(), log)
	}
}

// Preflight runs the landing's preflight on repo's main checkout, reading
// git through r and judging flai check at now (ADR-0234), and answers the
// findings that hold the queue, each to be corrected in verify.CheckoutMain:
// first what the lint of the tiers that carry a fix command still finds,
// after their fix, in the main checkout's uncommitted files but the rest,
// those outside the shared set of wip, design, and the shared paths, which a
// landing leaves as they are (ADR-0178, ADR-0248), as flai accept --batch
// lints them; then each error-level finding of flai check on the main
// checkout, unscoped. A warning holds nothing, since main always has some
// (D-0338). When the check finds an error and running, the flai that runs
// it, is behind main's flai sub-project, main is checked again with the
// manifest's flai.command, and only what that check finds counts; where it
// cannot settle them, the report says why beside the findings (ADR-0311).
// A lint finding in place is confirmed in a scratch worktree of
// the main branch's head, PreflightWorktreeName, with those files copied in
// and the fix run there, and only one that stands there counts; the
// worktree is removed after. It changes no file of the main checkout. An
// error is a preflight that could not run.
func Preflight(ctx context.Context, r execx.Runner, repo *workitem.Repo, now time.Time, running buildinfo.Info, log *slog.Logger) (PreflightReport, error) {
	if log == nil {
		log = slog.New(slog.DiscardHandler)
	}
	d, err := storygit.DirtyMain(r, repo)
	if err != nil {
		return PreflightReport{}, fmt.Errorf("run the merge queue's preflight: %w", err)
	}
	left := d.Rest
	lint, err := preflightLint(ctx, r, repo, left, log)
	if err != nil {
		return PreflightReport{}, err
	}
	checked, behind, err := preflightCheck(r, repo, now, running, left, log)
	if err != nil {
		return PreflightReport{}, err
	}
	return PreflightReport{Findings: append(lint, checked...), Behind: behind}, nil
}

// preflightLint is what the lint of the tiers that fix leaves in the main
// checkout's uncommitted files that a landing commits, all but left: nothing
// when the lint passes in place, and otherwise what still fails once the fix
// ran on copies of them in the scratch worktree.
func preflightLint(ctx context.Context, r execx.Runner, repo *workitem.Repo, left []string, log *slog.Logger) ([]verify.FileFinding, error) {
	root := repo.MainRoot
	changed, err := verify.ChangedPaths(r, root, "")
	if err != nil {
		return nil, fmt.Errorf("run the merge queue's preflight: %w; check that git can read the main checkout %s", err, root)
	}
	paths := verify.WithoutLeft(changed, left)
	if len(paths) == 0 {
		return nil, nil
	}
	res, err := verify.CheckoutUnfixed(ctx, root, paths)
	if err != nil {
		return nil, fmt.Errorf("run the merge queue's preflight: lint the main checkout's uncommitted files: %w", err)
	}
	if res.Passed {
		return nil, nil
	}
	return confirmLint(ctx, r, repo, paths, log)
}

// confirmLint runs the fix and then the lint of the tiers that fix over
// paths copied from the main checkout into a scratch worktree of the main
// branch's head, and answers what still fails there, by the paths the main
// checkout has them at. The scratch worktree is removed after, whatever
// came of it; one it cannot remove is logged at warn.
func confirmLint(ctx context.Context, r execx.Runner, repo *workitem.Repo, paths []string, log *slog.Logger) ([]verify.FileFinding, error) {
	root := repo.MainRoot
	_, head, err := mainHead(r, repo, "run the merge queue's preflight")
	if err != nil {
		return nil, err
	}
	scratch := filepath.Join(repo.CacheDir(), "worktrees", PreflightWorktreeName)
	if err := removeScratch(r, root, scratch); err != nil {
		return nil, err
	}
	if err := os.MkdirAll(filepath.Dir(scratch), 0o755); err != nil {
		return nil, fmt.Errorf("run the merge queue's preflight: make %s: %w; check that the project's .flai-cache is writable", filepath.Dir(scratch), err)
	}
	if _, err := r.Run(root, "git", "worktree", "add", "--quiet", "--detach", scratch, head); err != nil {
		return nil, fmt.Errorf("run the merge queue's preflight: check the main branch's head %s out in %s: %w; check that git can add a worktree there", short(head), scratch, err)
	}
	defer func() {
		if err := removeScratch(r, root, scratch); err != nil {
			log.Warn("merge queue preflight worktree not removed; the next preflight removes it", "component", "mergequeue", "path", scratch, "err", err.Error())
		}
	}()
	for _, p := range paths {
		if err := copyInto(root, scratch, p); err != nil {
			return nil, fmt.Errorf("run the merge queue's preflight: copy %s into %s: %w; check that the project's .flai-cache is writable", p, scratch, err)
		}
	}
	if _, err := verify.CheckoutFix(ctx, scratch, paths, log); err != nil {
		return nil, fmt.Errorf("run the merge queue's preflight: fix the copies of the main checkout's uncommitted files: %w", err)
	}
	res, err := verify.CheckoutUnfixed(ctx, scratch, paths)
	if err != nil {
		return nil, fmt.Errorf("run the merge queue's preflight: lint the fixed copies of the main checkout's uncommitted files: %w", err)
	}
	if res.Passed {
		return nil, nil
	}
	return unfixedFindings(scratch, res), nil
}

// unfixedFindings are the findings of res's failing tiers, run in the
// checkout at root, to be corrected in the main checkout, with one naming
// the tier alone for a failing tier that gives none, so that a failure
// whose output names no file still holds the queue.
func unfixedFindings(root string, res verify.Result) []verify.FileFinding {
	got := verify.FileFindings(root, verify.CheckoutMain, "", res)
	for _, t := range res.Tiers {
		if t.State != verify.Failed || slices.ContainsFunc(got, func(f verify.FileFinding) bool { return f.Tier == t.Name }) {
			continue
		}
		got = append(got, verify.FileFinding{Tier: t.Name, Checkout: verify.CheckoutMain, Message: "the " + t.Name + " tier fails on the main checkout's uncommitted files and names no file; run flai test on them in the main checkout"})
	}
	return got
}

// copyInto makes p, relative to the checkouts' roots, in the checkout at to
// what it is in the checkout at from: a copy of the regular file, or no file
// when from has none.
func copyInto(from, to, p string) error {
	src, dst := filepath.Join(from, filepath.FromSlash(p)), filepath.Join(to, filepath.FromSlash(p))
	info, err := os.Stat(src)
	if errors.Is(err, fs.ErrNotExist) {
		if err := os.Remove(dst); err != nil && !errors.Is(err, fs.ErrNotExist) {
			return err
		}
		return nil
	}
	if err != nil {
		return err
	}
	if !info.Mode().IsRegular() {
		return nil
	}
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}
	return os.WriteFile(dst, data, info.Mode().Perm())
}

// removeScratch removes the preflight's scratch worktree at path, whatever
// it holds, and prunes git's record of it.
func removeScratch(r execx.Runner, root, path string) error {
	if _, err := os.Stat(path); err == nil {
		if storygit.IsWorktree(r, path) {
			if _, err := r.Run(root, "git", "worktree", "remove", "--force", path); err != nil {
				return fmt.Errorf("remove the merge queue's preflight worktree %s: %w; remove it with git worktree remove --force %s", path, err, path)
			}
		} else if err := os.RemoveAll(path); err != nil {
			return fmt.Errorf("remove the merge queue's preflight folder %s: %w; remove it by hand", path, err)
		}
	}
	if _, err := r.Run(root, "git", "worktree", "prune"); err != nil {
		return fmt.Errorf("prune the merge queue's preflight worktree: %w; run git worktree prune in %s", err, root)
	}
	return nil
}

// preflightCheck is each error-level finding of flai check on repo's main
// checkout, unscoped, judged at now, but those on left, the uncommitted
// files outside the shared set a landing leaves as they are (ADR-0178,
// ADR-0248): neither the batch, built from main's commit, nor the landing
// sees them. When it finds one and running, the flai that ran it, is behind
// main's flai sub-project, the findings are those of the check the main
// checkout's flai.command makes (checkAgain); where that cannot settle
// them, they stand with the cause beside them (ADR-0311).
func preflightCheck(r execx.Runner, repo *workitem.Repo, now time.Time, running buildinfo.Info, left []string, log *slog.Logger) ([]verify.FileFinding, *FlaiBehind, error) {
	main := repo
	if repo.Root != repo.MainRoot {
		var err error
		if main, err = workitem.Open(repo.MainRoot); err != nil {
			return nil, nil, fmt.Errorf("run the merge queue's preflight: open the main checkout %s for flai check: %w", repo.MainRoot, err)
		}
		main.Workflow = repo.Workflow // the main checkout runs the queue's board
	}
	res, err := check.Run(main, now)
	if err != nil {
		return nil, nil, fmt.Errorf("run the merge queue's preflight: run flai check on the main checkout %s: %w; run flai check there and fix what it names", main.Root, err)
	}
	found := checkFindings(res, left)
	if len(found) == 0 {
		return nil, nil, nil
	}
	head, path, behind := runningBehind(r, main, running, log)
	if behind == nil {
		return found, nil, nil
	}
	again, why := checkAgain(r, main, head, path, left)
	if why != "" {
		return found, &FlaiBehind{Behind: *behind, Why: why}, nil
	}
	log.Info("merge queue preflight checked main again with a flai built from main, the running flai being behind it", "component", "mergequeue",
		"running", behind.Running, "behind_commits", behind.Commits, "command", main.Manifest.Flai.Command, "findings_before", len(found), "findings_after", len(again))
	return again, nil, nil
}

// runningBehind is main's head, the flai sub-project's path, and how far
// running is behind that path at main's head in the main checkout of repo
// (release.FlaiBehind), nil when it is not; a look git cannot make is
// logged at warn and taken as not behind, so the findings stand as found.
func runningBehind(r execx.Runner, repo *workitem.Repo, running buildinfo.Info, log *slog.Logger) (string, string, *release.Behind) {
	_, head, err := mainHead(r, repo, "see whether the running flai is behind main")
	if err != nil {
		log.Warn("merge queue preflight could not tell whether the running flai is behind main; its findings stand", "component", "mergequeue", "err", err.Error())
		return "", "", nil
	}
	path := flaiPath(repo)
	behind, err := release.FlaiBehind(r, repo.MainRoot, head, path, running)
	if err != nil {
		log.Warn("merge queue preflight could not tell whether the running flai is behind main; its findings stand", "component", "mergequeue", "err", err.Error())
		return "", "", nil
	}
	return head, path, behind
}

// flaiPath is the path of repo's sub-project named flai, slash-separated,
// or "" for the whole tree when it has none (ADR-0311).
func flaiPath(repo *workitem.Repo) string {
	for _, p := range repo.Manifest.Projects {
		if p.Name == "flai" {
			return filepath.ToSlash(p.Path)
		}
	}
	return ""
}

// checkAgain runs the manifest's flai.command of repo, the main checkout,
// as check --json there, unscoped, and answers its error-level findings but
// those on left, when version --json names a flai that is not behind head
// under path, so that a wrapper that fell back to the running flai cannot
// confirm its findings (ADR-0311). Otherwise it answers why the check does
// not count: the manifest names no flai.command, the command gave no
// result, or its flai is behind main too or did not say its version.
func checkAgain(r execx.Runner, repo *workitem.Repo, head, path string, left []string) ([]verify.FileFinding, string) {
	command := strings.TrimSpace(repo.Manifest.Flai.Command)
	if command == "" {
		return nil, "the manifest names no flai.command, a program that runs the flai main builds, so the preflight could not check main with a flai built from main"
	}
	root := repo.MainRoot
	program := filepath.Join(root, filepath.FromSlash(command))
	out, err := r.Run(root, program, "check", "--json")
	res, ok := check.ResultFromOutput(out) // check --json exits non-zero on an error it reports
	if !ok {
		return nil, fmt.Sprintf("%s check --json, the manifest's flai.command, gave no check result in %s (%s)", command, root, noResultCause(out, err))
	}
	out, err = r.Run(root, program, "version", "--json")
	var v struct {
		Version string `json:"version"`
		Commit  string `json:"commit"`
	}
	raw, ok := check.JSONObject(out)
	if err != nil || !ok || json.Unmarshal(raw, &v) != nil || v.Version == "" {
		return nil, fmt.Sprintf("%s version --json, the manifest's flai.command, did not say its version (%s), so its check cannot be told from the running flai's", command, noResultCause(out, err))
	}
	behind, err := release.FlaiBehind(r, root, head, path, buildinfo.Info{Version: v.Version, Commit: v.Commit})
	if err != nil {
		return nil, fmt.Sprintf("whether flai %s, which %s runs, is behind main could not be told: %s", v.Version, command, firstLine(err.Error()))
	}
	if behind != nil {
		return nil, fmt.Sprintf("flai %s, which %s runs, is behind main too: %s", v.Version, command, behind.Message)
	}
	return checkFindings(res, left), ""
}

// noResultCause is the first line of err when a command failed, else the
// first line of out, what it printed, or "it printed nothing".
func noResultCause(out string, err error) string {
	switch {
	case err != nil:
		return firstLine(err.Error())
	case strings.TrimSpace(out) != "":
		return "it printed " + firstLine(out)
	}
	return "it printed nothing"
}

// firstLine is s up to its first newline, trimmed.
func firstLine(s string) string {
	line, _, _ := strings.Cut(strings.TrimSpace(s), "\n")
	return strings.TrimSpace(line)
}

// checkFindings are res's error-level findings as preflight findings, in
// their order, each path slash-separated, but those on a path left names or
// holds.
func checkFindings(res *check.Result, left []string) []verify.FileFinding {
	var got []verify.FileFinding
	for _, f := range res.Findings {
		if f.Level != check.Error || verify.IsLeft(filepath.ToSlash(f.Path), left) {
			continue
		}
		got = append(got, verify.FileFinding{Tier: PreflightCheckTier, Checkout: verify.CheckoutMain, Path: filepath.ToSlash(f.Path), Line: f.Line,
			Level: f.Level, Rule: f.Rule, Message: f.Message})
	}
	return got
}

// preflightLook runs d.Preflight when a batch is about to be built, and
// otherwise only while the preflight holds the queue, so that a clean main
// clears the hold, and reports whether the queue may build: false while
// the preflight finds anything. Findings record the hold, or refresh it,
// keeping when it began and its thread; a clean run clears it. A preflight
// that could not run is an error, and leaves the hold as it was.
func preflightLook(ctx context.Context, d Deps, repo *workitem.Repo, building bool) (bool, error) {
	st, err := ReadState(repo.CacheDir())
	if err != nil {
		return false, err
	}
	if !building && !st.PreflightHeld() {
		return true, nil
	}
	report, err := d.Preflight(ctx, repo)
	if err != nil {
		return false, err
	}
	findings := report.Findings
	head := ""
	if len(findings) > 0 {
		if _, head, err = mainHead(d.Runner, repo, "record the merge queue's preflight hold"); err != nil {
			return false, err
		}
	}
	var prev *PreflightHold
	now := d.Now().UTC()
	if _, err := Update(repo.CacheDir(), func(s *State) error {
		prev = s.Preflight
		if len(findings) == 0 {
			s.Preflight = nil
			return nil
		}
		hold := &PreflightHold{Since: now, Main: head, Findings: findings, FlaiBehind: report.Behind}
		if prev != nil {
			hold.Since, hold.Thread = prev.Since, prev.Thread
		}
		s.Preflight = hold
		return nil
	}); err != nil {
		return false, err
	}
	switch {
	case len(findings) == 0 && prev != nil:
		d.Log.Info("merge queue preflight hold cleared", "component", "mergequeue", "held_since", prev.Since, "main", short(prev.Main))
	case len(findings) > 0 && prev == nil:
		attrs := []any{"component", "mergequeue", "main", short(head), "files", findingFiles(findings), "rules", findingRules(findings)}
		if b := report.Behind; b != nil {
			attrs = append(attrs, "flai_behind", b.Message, "why", b.Why)
		}
		d.Log.Warn("merge queue held by its preflight; it builds no batch until main is clean", attrs...)
	case len(findings) > 0 && !slices.Equal(prev.Findings, findings):
		d.Log.Info("merge queue preflight findings changed; the queue stays held", "component", "mergequeue", "main", short(head),
			"files", findingFiles(findings), "rules", findingRules(findings))
	}
	return len(findings) == 0, nil
}

// findingFiles are the files findings name, sorted, each once.
func findingFiles(findings []verify.FileFinding) []string {
	return distinct(findings, func(f verify.FileFinding) string { return f.Path })
}

// findingRules are the rules findings name, sorted, each once.
func findingRules(findings []verify.FileFinding) []string {
	return distinct(findings, func(f verify.FileFinding) string { return f.Rule })
}

// distinct are the non-empty values of field over findings, sorted, each
// once.
func distinct(findings []verify.FileFinding, field func(verify.FileFinding) string) []string {
	var out []string
	for _, f := range findings {
		if v := strings.TrimSpace(field(f)); v != "" {
			out = append(out, v)
		}
	}
	slices.Sort(out)
	return slices.Compact(out)
}
