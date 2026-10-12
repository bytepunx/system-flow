package verify

import (
	"bytes"
	"context"
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
	if proc == nil {
		proc = OS{}
	}
	if log == nil {
		log = slog.New(slog.DiscardHandler)
	}
	fsys := os.DirFS(root)
	existing := regularFiles(fsys, paths)
	before := map[string][]byte{}
	for _, t := range tiers {
		if len(t.FixCommand) == 0 || ctx.Err() != nil {
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
		for _, p := range matched {
			if _, seen := before[p]; !seen {
				if data, err := fs.ReadFile(fsys, p); err == nil {
					before[p] = data
				}
			}
		}
		runFix(ctx, root, t, argv, proc, log)
	}
	var changed []string
	for p, was := range before {
		if now, err := fs.ReadFile(fsys, p); err != nil || !bytes.Equal(now, was) {
			changed = append(changed, p)
		}
	}
	slices.Sort(changed)
	return changed
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
// tier's command, and logs at warn when it cannot start or exits non-zero.
func runFix(ctx context.Context, root string, t Tier, argv []string, proc Proc, log *slog.Logger) {
	var combined lockedBuffer
	dir := filepath.Join(root, filepath.FromSlash(cleanDir(t.Dir)))
	env := []string{roleEnv + "=" + conventions.RoleVerify}
	exit, err := proc.Run(ctx, dir, argv, env, &combined, &combined)
	if err == nil && exit == 0 {
		return
	}
	args := []any{"component", "verify", "tier", t.Name, "command", strings.Join(argv, " "), "exit_code", exit, "output", lastLines(combined.String(), fixLogLines)}
	if err != nil {
		args = append(args, "err", err.Error())
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
