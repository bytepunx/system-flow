package verify

import (
	"io/fs"
	"path"
	"slices"
	"strings"

	"github.com/bytepunx/system-flow/flai/internal/manifest"
)

// The placeholders a tier's command may hold as whole arguments.
const (
	packagesArg = manifest.PlaceholderPackages
	filesArg    = manifest.PlaceholderFiles
)

// Selected is a tier chosen for a run, with the argv it runs.
type Selected struct {
	Tier Tier
	Argv []string
}

// Select chooses the tiers, in list order, that the root-relative paths
// select, each with its placeholders filled; with all, every tier, AllOnly
// ones too. A tier is selected by a path that one of its patterns matches
// and none of its ! patterns does, and is left out when a placeholder in its
// command comes to nothing, as {packages} does for files no Go package
// holds. fsys is the checkout, read for the directories that hold .go files.
func Select(fsys fs.FS, tiers []Tier, paths []string, all bool) []Selected {
	var out []Selected
	for _, t := range tiers {
		if all {
			out = append(out, Selected{Tier: t, Argv: allArgv(t)})
			continue
		}
		if t.AllOnly {
			continue
		}
		matched := matching(t.Paths, paths)
		if len(matched) == 0 {
			continue
		}
		dir := cleanDir(t.Dir)
		if argv, ok := fill(t.Command, files(dir, matched), packages(fsys, dir, matched)); ok {
			out = append(out, Selected{Tier: t, Argv: argv})
		}
	}
	return out
}

// allArgv is the argv a tier runs when every tier runs.
func allArgv(t Tier) []string {
	if len(t.AllCommand) > 0 {
		return slices.Clone(t.AllCommand)
	}
	argv, _ := fill(t.Command, []string{"."}, []string{"./..."})
	return argv
}

// matching are the paths that a pattern matches and no ! pattern does.
func matching(patterns, paths []string) []string {
	var in, out manifest.Claims
	for _, p := range patterns {
		if rest, ok := strings.CutPrefix(p, "!"); ok {
			out.Shared = append(out.Shared, rest)
		} else {
			in.Shared = append(in.Shared, p)
		}
	}
	var got []string
	for _, p := range paths {
		if _, ok := in.Match(p); !ok {
			continue
		}
		if _, ok := out.Match(p); ok {
			continue
		}
		got = append(got, p)
	}
	return got
}

// fill is argv with each placeholder replaced by its list, and whether
// every placeholder it holds came to something.
func fill(argv, files, packages []string) ([]string, bool) {
	var out []string
	for _, a := range argv {
		switch a {
		case filesArg:
			if len(files) == 0 {
				return nil, false
			}
			out = append(out, files...)
		case packagesArg:
			if len(packages) == 0 {
				return nil, false
			}
			out = append(out, packages...)
		default:
			out = append(out, a)
		}
	}
	return out, true
}

// files are the paths below dir, relative to it, sorted.
func files(dir string, paths []string) []string {
	var out []string
	for _, p := range paths {
		if rel, ok := below(dir, p); ok && rel != "." {
			out = append(out, rel)
		}
	}
	slices.Sort(out)
	return slices.Compact(out)
}

// packages are the Go packages that hold the paths below dir, as ./x/y or
// ., relative to it, sorted: for each path, the nearest directory at or
// above it within dir that holds a .go file. A directory the go tool ignores
// (testdata, or a name beginning with . or _) belongs to the package above.
func packages(fsys fs.FS, dir string, paths []string) []string {
	var out []string
	for _, p := range paths {
		rel, ok := below(dir, p)
		if !ok {
			continue
		}
		if pkg, ok := nearestPackage(fsys, dir, path.Dir(rel)); ok {
			out = append(out, pkg)
		}
	}
	slices.Sort(out)
	return slices.Compact(out)
}

// nearestPackage is the directory at or above d, relative to dir, that holds
// a .go file, as the go tool names it.
func nearestPackage(fsys fs.FS, dir, d string) (string, bool) {
	if d != "." {
		segs := strings.Split(d, "/")
		for i, s := range segs {
			if s == "testdata" || strings.HasPrefix(s, ".") || strings.HasPrefix(s, "_") {
				d = path.Join(append([]string{"."}, segs[:i]...)...)
				break
			}
		}
	}
	for {
		if holdsGo(fsys, path.Join(dir, d)) {
			if d == "." {
				return ".", true
			}
			return "./" + d, true
		}
		if d == "." {
			return "", false
		}
		d = path.Dir(d)
	}
}

// holdsGo reports whether the directory d of fsys holds a .go file.
func holdsGo(fsys fs.FS, d string) bool {
	if d == "" {
		d = "."
	}
	entries, err := fs.ReadDir(fsys, d)
	if err != nil {
		return false
	}
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".go") {
			return true
		}
	}
	return false
}

// cleanDir is a tier's Dir as paths are compared with it: "" for the root.
func cleanDir(dir string) string {
	d := path.Clean(strings.TrimSpace(dir))
	if d == "." || d == "/" {
		return ""
	}
	return strings.TrimPrefix(d, "./")
}

// below is p relative to dir, and whether p is dir or below it; dir "" is
// the root, which holds every path.
func below(dir, p string) (string, bool) {
	if dir == "" {
		return p, true
	}
	if p == dir {
		return ".", true
	}
	rest, ok := strings.CutPrefix(p, dir+"/")
	return rest, ok
}
