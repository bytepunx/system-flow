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
	packagesArg  = manifest.PlaceholderPackages
	filesArg     = manifest.PlaceholderFiles
	importersArg = manifest.PlaceholderImporters
)

// Selected is a tier chosen for a run, with the argv it runs.
type Selected struct {
	Tier Tier
	Argv []string
	// CoveredBy names the selected tier that covers this one, which Run
	// skips rather than runs; "" when it runs.
	CoveredBy string
}

// Select chooses the tiers, in list order, that the root-relative paths
// select, each with its placeholders filled; with all, every tier, AllOnly
// ones too. A tier is selected by a path that one of its patterns matches
// and none of its ! patterns does, and is left out when a placeholder in its
// command comes to nothing, as {packages} and {importers} do for files no Go
// package holds. fsys is the checkout, read for the directories that hold .go
// files and, for {importers}, for the imports of each.
// A selected tier that another selected tier covers is marked CoveredBy.
func Select(fsys fs.FS, tiers []Tier, paths []string, all bool) []Selected {
	return cover(tiers, selectTiers(fsys, tiers, paths, all))
}

// selectTiers is Select before the covered tiers are marked.
func selectTiers(fsys fs.FS, tiers []Tier, paths []string, all bool) []Selected {
	var out []Selected
	for _, t := range tiers {
		if all {
			out = append(out, Selected{Tier: t, Argv: allArgv(t)})
			continue
		}
		if t.AllOnly {
			continue
		}
		if s, ok := selectByPaths(fsys, t, paths); ok {
			out = append(out, s)
		}
	}
	return out
}

// SelectStory chooses the tiers, in list order, that a story's close-out
// runs for the paths its branch changed: those Select chooses without all,
// and the AllOnly tiers that the paths select or that have no Paths, which
// run as they do with all, each tier another of them covers marked as Select
// marks it. Verify runs them.
func SelectStory(fsys fs.FS, tiers []Tier, paths []string) []Selected {
	return cover(tiers, selectStoryTiers(fsys, tiers, paths))
}

// selectStoryTiers is SelectStory before the covered tiers are marked.
func selectStoryTiers(fsys fs.FS, tiers []Tier, paths []string) []Selected {
	var out []Selected
	for _, t := range tiers {
		if !t.AllOnly {
			if s, ok := selectByPaths(fsys, t, paths); ok {
				out = append(out, s)
			}
			continue
		}
		if len(t.Paths) == 0 || len(matching(t.Paths, paths)) > 0 {
			out = append(out, Selected{Tier: t, Argv: allArgv(t)})
		}
	}
	return out
}

// cover marks each selected tier that another selected tier covers, by
// naming it in Covers or through a chain of tiers that name the next, as
// covered by the nearest such tier that is not covered itself, so that it
// names the tier that runs in its place: nearest in links of the chain, then
// first in list order. A chain may pass through tiers that were not
// selected, as tiers holds them all; a tier that only unselected tiers cover
// is left to run. Names are matched trimmed, each the first tier of it.
func cover(tiers []Tier, selected []Selected) []Selected {
	first := map[string]int{}
	for i, t := range tiers {
		name := strings.TrimSpace(t.Name)
		if _, ok := first[name]; !ok {
			first[name] = i
		}
	}
	// coverers[j] are the tiers that name tier j in Covers, in list order
	coverers := make([][]int, len(tiers))
	for i, t := range tiers {
		for _, c := range t.Covers {
			if j, ok := first[strings.TrimSpace(c)]; ok && j != i {
				coverers[j] = append(coverers[j], i)
			}
		}
	}
	at := make([]int, len(selected))
	chosen := map[int]bool{}
	for n, s := range selected {
		at[n] = -1
		if i, ok := first[strings.TrimSpace(s.Tier.Name)]; ok {
			at[n], chosen[i] = i, true
		}
	}
	// above[n] are the tiers that cover selected[n], directly or through a
	// chain, nearest first
	above := make([][]int, len(selected))
	covered := map[int]bool{}
	for n, i := range at {
		if i < 0 {
			continue
		}
		seen := map[int]bool{i: true}
		for next := slices.Clone(coverers[i]); len(next) > 0; next = next[1:] {
			k := next[0]
			if seen[k] {
				continue
			}
			seen[k] = true
			above[n] = append(above[n], k)
			next = append(next, coverers[k]...)
			if chosen[k] {
				covered[i] = true
			}
		}
	}
	for n, i := range at {
		if !covered[i] {
			continue
		}
		for _, k := range above[n] {
			if chosen[k] && !covered[k] {
				selected[n].CoveredBy = tiers[k].Name
				break
			}
		}
	}
	return selected
}

// selectByPaths is the tier with its placeholders filled for the paths it
// selects, and whether it is selected: some path matches it and every
// placeholder in its command comes to something.
func selectByPaths(fsys fs.FS, t Tier, paths []string) (Selected, bool) {
	matched := matching(t.Paths, paths)
	if len(matched) == 0 {
		return Selected{}, false
	}
	dir := cleanDir(t.Dir)
	pkgs := packages(fsys, dir, matched)
	var imps []string
	if slices.Contains(t.Command, importersArg) {
		imps = importers(fsys, dir, pkgs)
	}
	argv, ok := fill(t.Command, map[string][]string{filesArg: files(dir, matched), packagesArg: pkgs, importersArg: imps})
	return Selected{Tier: t, Argv: argv}, ok
}

// allArgv is the argv a tier runs when every tier runs.
func allArgv(t Tier) []string {
	if len(t.AllCommand) > 0 {
		return slices.Clone(t.AllCommand)
	}
	argv, _ := fill(t.Command, map[string][]string{filesArg: {"."}, packagesArg: {"./..."}, importersArg: {"./..."}})
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

// fill is argv with each placeholder that lists holds replaced by its list,
// and whether every placeholder it holds came to something.
func fill(argv []string, lists map[string][]string) ([]string, bool) {
	var out []string
	for _, a := range argv {
		list, ok := lists[a]
		if !ok {
			out = append(out, a)
			continue
		}
		if len(list) == 0 {
			return nil, false
		}
		out = append(out, list...)
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
			if ignoredDir(s) {
				d = path.Join(append([]string{"."}, segs[:i]...)...)
				break
			}
		}
	}
	for {
		if holdsGo(fsys, path.Join(dir, d)) {
			return packageArg(d), true
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
