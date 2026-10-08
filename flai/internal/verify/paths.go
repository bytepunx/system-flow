package verify

import (
	"fmt"
	"io/fs"
	"path"
	"path/filepath"
	"slices"
	"strconv"
	"strings"

	"github.com/bytepunx/system-flow/flai/internal/execx"
)

// skipDirs are the folders a folder argument does not reach into: installed
// dependencies, git's own, and flai's cache, which holds other worktrees.
var skipDirs = []string{"node_modules", ".git", ".flai-cache"}

// Resolve turns arguments, files or folders relative to the checkout's root,
// into the files they name, sorted: a file stands for itself and a folder
// for every file below it. fsys is the checkout.
func Resolve(fsys fs.FS, args []string) ([]string, error) {
	var out []string
	for _, arg := range args {
		p := path.Clean(filepath.ToSlash(strings.TrimSpace(arg)))
		if !fs.ValidPath(p) {
			return nil, fmt.Errorf("%s is not a path inside the checkout; name a file or folder relative to its root", arg)
		}
		info, err := fs.Stat(fsys, p)
		if err != nil {
			return nil, fmt.Errorf("%s is not a file or folder in the checkout; name one relative to its root", arg)
		}
		if !info.IsDir() {
			out = append(out, p)
			continue
		}
		err = fs.WalkDir(fsys, p, func(q string, d fs.DirEntry, err error) error {
			switch {
			case err != nil:
				return err
			case d.IsDir() && q != p && slices.Contains(skipDirs, d.Name()):
				return fs.SkipDir
			case !d.IsDir():
				out = append(out, q)
			}
			return nil
		})
		if err != nil {
			return nil, fmt.Errorf("list the files in %s: %w", arg, err)
		}
	}
	slices.Sort(out)
	return slices.Compact(out), nil
}

// ChangedPaths are the files the checkout at dir changed against base,
// relative to its root and sorted: those its commits changed since it left
// base, and those changed, staged, or untracked and not yet committed. A
// file deleted is not among them. With base "" only the uncommitted count.
func ChangedPaths(r execx.Runner, dir, base string) ([]string, error) {
	set := map[string]bool{}
	if base != "" {
		out, err := r.Run(dir, "git", "diff", "--name-only", "--diff-filter=d", base+"...HEAD")
		if err != nil {
			return nil, fmt.Errorf("list what %s changed against %s; check that %s is a branch of its repository: %w", dir, base, base, err)
		}
		for _, line := range strings.Split(out, "\n") {
			if p := unquote(line); p != "" {
				set[p] = true
			}
		}
	}
	out, err := r.Run(dir, "git", "status", "--porcelain", "--untracked-files=all")
	if err != nil {
		return nil, fmt.Errorf("list what %s has not committed: %w", dir, err)
	}
	for _, line := range strings.Split(out, "\n") {
		xy, p, ok := statusEntry(line)
		if !ok {
			continue
		}
		if strings.Contains(xy, "D") {
			delete(set, p)
		} else {
			set[p] = true
		}
	}
	paths := make([]string, 0, len(set))
	for p := range set {
		paths = append(paths, p)
	}
	slices.Sort(paths)
	return paths, nil
}

// BaseChanges are what base changed since the checkout at dir left it: the
// paths changed from their merge base to base, sorted, a rename as both of
// its paths and a deletion among them, and the short hashes of the commits
// of base that dir lacks, newest first (ADR-0135).
func BaseChanges(r execx.Runner, dir, base string) (paths, commits []string, err error) {
	out, err := r.Run(dir, "git", "diff", "--no-renames", "--name-only", "HEAD..."+base)
	if err != nil {
		return nil, nil, fmt.Errorf("list what %s changed since the checkout left it: %w", base, err)
	}
	for _, line := range strings.Split(out, "\n") {
		if p := unquote(line); p != "" {
			paths = append(paths, p)
		}
	}
	slices.Sort(paths)
	out, err = r.Run(dir, "git", "rev-list", "--abbrev-commit", "HEAD.."+base)
	if err != nil {
		return nil, nil, fmt.Errorf("list the commits of %s the checkout lacks: %w", base, err)
	}
	commits = strings.Fields(out)
	return slices.Compact(paths), commits, nil
}

// statusEntry is a git status --porcelain line's two status letters and the
// path it is now at. The runner trims its output, which takes the leading
// space from the first line's " M path", so a line whose third character is
// not the space before the path is read as having lost it.
func statusEntry(line string) (xy, p string, ok bool) {
	switch {
	case len(line) >= 4 && line[2] == ' ':
		xy, p = line[:2], line[3:]
	case len(line) >= 3 && line[1] == ' ':
		xy, p = " "+line[:1], line[2:]
	default:
		return "", "", false
	}
	if _, to, renamed := strings.Cut(p, " -> "); renamed {
		p = to
	}
	p = unquote(p)
	return xy, p, p != ""
}

// unquote is a path as git wrote it, without the quotes and escapes it puts
// around a name with unusual characters.
func unquote(p string) string {
	p = strings.TrimSpace(p)
	if strings.HasPrefix(p, `"`) {
		if u, err := strconv.Unquote(p); err == nil {
			return u
		}
	}
	return p
}
