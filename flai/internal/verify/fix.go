package verify

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/fs"
	"log/slog"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"

	"github.com/bytepunx/system-flow/flai/internal/conventions"
)

// fixLogLines is how many of a failing fix command's last lines of output
// its warning carries.
const fixLogLines = 5

// Fix runs, in list order, the fix command of each tier that some of the
// root-relative paths select, over the files among them that exist, with its
// placeholders filled as Select fills Command's, and answers the paths whose
// content the fixes changed or that they removed, sorted. A fix command that
// cannot start or exits non-zero is logged at warn and the next tier's runs;
// a tier whose
// placeholder comes to nothing is skipped. When ctx ends no further tier
// starts. A nil log discards.
func Fix(ctx context.Context, root string, tiers []Tier, paths []string, proc Proc, log *slog.Logger) []string {
	if log == nil {
		log = slog.New(slog.DiscardHandler)
	}
	changed, _, _ := fixTiers(ctx, root, tiers, paths, proc, func(t Tier, argv []string, out output) {
		if out.err == nil && out.exit == 0 {
			return
		}
		warnFix(log, t, argv, out)
	})
	return changed
}

// FixOutcome is what a fix did and what its commands' exits say remains of
// the lint it fixes (S-0609).
type FixOutcome struct {
	// Changed are the root-relative paths the fixes changed or removed,
	// sorted.
	Changed []string
	// Remains is, as Unfixed answers it, what the fix commands' own runs say
	// they left: a tier whose fix exited 0 passed, one that exited non-zero
	// failed with the findings its output gives in the tier's format, capped
	// across the tiers at DefaultMax. It does not pass while Unknown names a
	// tier.
	Remains Result
	// Unknown names the tiers whose fix's run does not say what remains: its
	// command could not start, ctx ended before or while it ran, or its
	// output, in a format other than plain, gives no finding of its own.
	// Their check is to be run, as CheckoutUnfixed runs it.
	Unknown []string
}

// Known reports whether the fix's runs say what remains of every tier they
// ran, so Remains stands in for a check after the fix.
func (o FixOutcome) Known() bool { return len(o.Unknown) == 0 }

// FixRemaining runs the fix commands as Fix does, each once, and answers
// what they changed and what their exits say remains, so that a caller needs
// no check after the fix where Known holds. A fix command that cannot start
// is logged at warn; one that exits non-zero is lint that remains, not a
// failure of the fix, and is not logged. A tier that is not run, its paths
// or its placeholder coming to nothing, leaves nothing.
func FixRemaining(ctx context.Context, root string, tiers []Tier, paths []string, proc Proc, log *slog.Logger) FixOutcome {
	if log == nil {
		log = slog.New(slog.DiscardHandler)
	}
	type ran struct {
		tier Tier
		argv []string
		out  output
	}
	var runs []ran
	changed, existing, skipped := fixTiers(ctx, root, tiers, paths, proc, func(t Tier, argv []string, out output) {
		if out.err != nil {
			warnFix(log, t, argv, out)
		}
		runs = append(runs, ran{tier: t, argv: argv, out: out})
	})
	o := FixOutcome{Changed: changed, Remains: Result{Passed: true, Paths: existing}, Unknown: skipped}
	e := env{root: root, fsys: os.DirFS(root)}
	budget := DefaultMax
	for _, r := range runs {
		tr := TierResult{Name: r.tier.Name, Command: r.argv, Dir: r.tier.Dir, State: Passed}
		if r.out.err == nil {
			exit := r.out.exit
			tr.ExitCode = &exit
		}
		switch {
		case r.out.err != nil:
			tr.State = NotReached
			o.Unknown = append(o.Unknown, r.tier.Name)
		case r.out.exit != 0:
			e.dir = cleanDir(r.tier.Dir)
			all := parse(r.tier, r.out, e)
			if !ownFindings(r.tier, r.out, all) {
				tr.State = NotReached
				o.Unknown = append(o.Unknown, r.tier.Name)
				break
			}
			tr.State = Failed
			o.Remains.Passed = false
			n := min(budget, len(all))
			tr.Findings, tr.Omitted = all[:n], len(all)-n
			budget -= n
		}
		o.Remains.Tiers = append(o.Remains.Tiers, tr)
	}
	if !o.Known() {
		o.Remains.Passed = false
	}
	return o
}

// ownFindings reports whether the findings parse gave a failing tier's
// output are its format's own: always for a plain tier, whose finding is its
// last lines; for another format, unless parse fell back to the plain one.
func ownFindings(t Tier, out output, found []Finding) bool {
	if t.Format == "" || t.Format == FormatPlain {
		return true
	}
	return len(found) != 1 || found[0] != plain(t.Name, out)
}

// CheckoutFixRemaining runs FixRemaining over the paths of the checkout at
// root with the tiers its manifest declares and the commands as processes.
// It returns an error only when the tiers cannot be read.
func CheckoutFixRemaining(ctx context.Context, root string, paths []string, log *slog.Logger) (FixOutcome, error) {
	tiers, err := CheckoutTiers(root)
	if err != nil {
		return FixOutcome{}, err
	}
	return FixRemaining(ctx, root, tiers, paths, OS{}, log), nil
}

// fixTiers runs, in list order, the fix command of each tier that some of
// the existing files among paths select, as Fix describes, handing each run
// to ran, and answers the paths the fixes changed or removed, sorted, the
// existing files, and the tiers it would have run but for ctx ending.
func fixTiers(ctx context.Context, root string, tiers []Tier, paths []string, proc Proc, ran func(Tier, []string, output)) (changed, existing, skipped []string) {
	if proc == nil {
		proc = OS{}
	}
	fsys := os.DirFS(root)
	existing = regularFiles(fsys, paths)
	before := map[string][]byte{}
	for _, t := range tiers {
		if len(t.FixCommand) == 0 {
			continue
		}
		matched := matching(t.Paths, existing)
		if len(matched) == 0 {
			continue
		}
		argv, ok := fixArgv(fsys, t, matched)
		if !ok {
			continue
		}
		if ctx.Err() != nil {
			skipped = append(skipped, t.Name)
			continue
		}
		for _, p := range matched {
			if _, seen := before[p]; !seen {
				if data, err := fs.ReadFile(fsys, p); err == nil {
					before[p] = data
				}
			}
		}
		ran(t, argv, runFix(ctx, root, t, argv, proc))
	}
	for p, was := range before {
		if now, err := fs.ReadFile(fsys, p); err != nil || !bytes.Equal(now, was) {
			changed = append(changed, p)
		}
	}
	slices.Sort(changed)
	return changed, existing, skipped
}

// CheckoutFix runs Fix over the paths of the checkout at root with the tiers
// its manifest declares and the commands as processes. It returns an error
// only when the tiers cannot be read.
func CheckoutFix(ctx context.Context, root string, paths []string, log *slog.Logger) ([]string, error) {
	tiers, err := CheckoutTiers(root)
	if err != nil {
		return nil, err
	}
	return Fix(ctx, root, tiers, paths, OS{}, log), nil
}

// regularFiles are the paths, cleaned, that name regular files in fsys,
// sorted and each once.
func regularFiles(fsys fs.FS, paths []string) []string {
	var out []string
	for _, p := range paths {
		p = path.Clean(filepath.ToSlash(strings.TrimSpace(p)))
		if !fs.ValidPath(p) || p == "." {
			continue
		}
		if info, err := fs.Stat(fsys, p); err == nil && info.Mode().IsRegular() {
			out = append(out, p)
		}
	}
	slices.Sort(out)
	return slices.Compact(out)
}

// fixArgv is the tier's fix command with its placeholders filled for the
// matched paths, relative to its Dir, and whether every placeholder came to
// something.
func fixArgv(fsys fs.FS, t Tier, matched []string) ([]string, bool) {
	dir := cleanDir(t.Dir)
	pkgs := packages(fsys, dir, matched)
	var imps []string
	if slices.Contains(t.FixCommand, importersArg) {
		imps = importers(fsys, dir, pkgs)
	}
	return fill(t.FixCommand, map[string][]string{filesArg: files(dir, matched), packagesArg: pkgs, importersArg: imps})
}

// runFix runs one tier's fix command in its Dir below root, as Run runs a
// tier's command, and answers what it did; a command that ran on as ctx
// ended answers that it was stopped.
func runFix(ctx context.Context, root string, t Tier, argv []string, proc Proc) output {
	var stdout, stderr bytes.Buffer
	var combined lockedBuffer
	dir := filepath.Join(root, filepath.FromSlash(cleanDir(t.Dir)))
	env := []string{roleEnv + "=" + conventions.RoleVerify}
	exit, err := proc.Run(ctx, dir, argv, env, io.MultiWriter(&stdout, &combined), io.MultiWriter(&stderr, &combined))
	if err == nil && ctx.Err() != nil {
		err = fmt.Errorf("stopped: %w", ctx.Err())
	}
	return output{stdout: stdout.String(), stderr: stderr.String(), combined: combined.String(), exit: exit, err: err}
}

// warnFix logs at warn a tier's fix command that could not start or exited
// non-zero.
func warnFix(log *slog.Logger, t Tier, argv []string, out output) {
	args := []any{"component", "verify", "tier", t.Name, "command", strings.Join(argv, " "), "exit_code", out.exit, "output", lastLines(out.combined, fixLogLines)}
	if out.err != nil {
		args = append(args, "err", out.err.Error())
	}
	log.Warn("a tier's fix command failed; what it changed before it stopped is kept", args...)
}

// Unfixed runs, as Run runs them, the commands of the tiers that carry a fix
// command and that the existing files among the root-relative paths select,
// over those files, and answers what their fix left: a pass, or the first
// tier that fails with its findings (S-0367, TH-0413). flai accept refuses on
// a failure. With no such tier selected it passes.
func Unfixed(ctx context.Context, root string, tiers []Tier, paths []string, opts RunOptions) (Result, error) {
	var fixing []Tier
	for _, t := range tiers {
		if len(t.FixCommand) > 0 {
			fixing = append(fixing, t)
		}
	}
	fsys := os.DirFS(root)
	existing := regularFiles(fsys, paths)
	res, err := Run(ctx, root, Select(fsys, fixing, existing, false), opts)
	res.Paths = existing
	return res, err
}

// CheckoutUnfixed runs Unfixed over the paths of the checkout at root with
// the tiers its manifest declares and the commands as processes. It returns
// an error when the tiers cannot be read or ctx ends.
func CheckoutUnfixed(ctx context.Context, root string, paths []string) (Result, error) {
	tiers, err := CheckoutTiers(root)
	if err != nil {
		return Result{}, err
	}
	return Unfixed(ctx, root, tiers, paths, RunOptions{})
}
