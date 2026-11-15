package verify

import (
	"io/fs"
	"path"
	"slices"
	"strings"
)

// OverlappingPaths are the paths of mainPaths that overlap changed, a story's
// change set, as ADR-0303 defines it, in mainPaths' order: those changed
// holds, those in a Go package of changed's, main's Go files that import one,
// and those in a package a Go file of changed's imports.
func OverlappingPaths(mainFS, storyFS fs.FS, changed, mainPaths []string) []string {
	if len(changed) == 0 || len(mainPaths) == 0 {
		return nil
	}
	main, story := newRepoView(mainFS), newRepoView(storyFS)
	set := newChangeSet(changed, story, main)
	seen := map[string]bool{}
	var out []string
	for _, p := range mainPaths {
		if seen[p] {
			continue
		}
		seen[p] = true
		if set.overlaps(p, main, story) {
			out = append(out, p)
		}
	}
	return out
}

// changeSet is a story's change set as overlap reads it.
type changeSet struct {
	// paths are the paths the story changes.
	paths map[string]bool
	// pkgs are the Go packages of those paths.
	pkgs map[string]bool
	// imported are the Go packages the story's changed Go files import.
	imported map[string]bool
}

// newChangeSet reads changed in story, falling back to main for a path the
// story deleted.
func newChangeSet(changed []string, story, main *repoView) changeSet {
	set := changeSet{paths: map[string]bool{}, pkgs: map[string]bool{}, imported: map[string]bool{}}
	for _, c := range changed {
		set.paths[c] = true
		if pkg, ok := packageOf(c, story, main); ok {
			set.pkgs[pkg] = true
		}
		for _, imp := range story.imports(c) {
			set.imported[imp] = true
		}
	}
	return set
}

// overlaps reports whether main's path p overlaps the change set.
func (s changeSet) overlaps(p string, main, story *repoView) bool {
	if s.paths[p] {
		return true
	}
	if pkg, ok := packageOf(p, main, story); ok && (s.pkgs[pkg] || s.imported[pkg]) {
		return true
	}
	return slices.ContainsFunc(main.imports(p), func(imp string) bool { return s.pkgs[imp] })
}

// packageOf is p's Go package, read in the checkout that holds p, first before
// second, and in the other when that one finds none.
func packageOf(p string, first, second *repoView) (string, bool) {
	if !first.holds(p) && second.holds(p) {
		first, second = second, first
	}
	if pkg, ok := first.pkg(p); ok {
		return pkg, true
	}
	return second.pkg(p)
}

// repoView is a checkout, rooted at the repository's root, with what overlap has
// read of it cached so that each file is read once.
type repoView struct {
	fsys    fs.FS
	modules map[string]module
	goDirs  map[string]bool
	imps    map[string][]string
}

// module is the Go module a directory lies in.
type module struct {
	// root is the directory of its go.mod.
	root string
	// path is its module path, "" when go.mod names none.
	path string
	// ok is false for a directory outside every module.
	ok bool
}

func newRepoView(fsys fs.FS) *repoView {
	return &repoView{fsys: fsys, modules: map[string]module{}, goDirs: map[string]bool{}, imps: map[string][]string{}}
}

// holds reports whether p names a file or directory of the checkout.
func (v *repoView) holds(p string) bool {
	_, err := fs.Stat(v.fsys, p)
	return err == nil
}

// module is the module of the nearest go.mod at or above dir.
func (v *repoView) module(dir string) module {
	if m, ok := v.modules[dir]; ok {
		return m
	}
	var m module
	if isFile(v.fsys, path.Join(dir, "go.mod")) {
		m = module{root: dir, path: modulePath(env{dir: dir, fsys: v.fsys}), ok: true}
	} else if parent := path.Dir(dir); parent != dir {
		m = v.module(parent)
	}
	v.modules[dir] = m
	return m
}

// pkg is the Go package p is in: the nearest directory from p's own up to
// its module's root that holds a .go file the go tool reads.
func (v *repoView) pkg(p string) (string, bool) {
	dir := path.Dir(p)
	m := v.module(dir)
	if !m.ok {
		return "", false
	}
	for d := dir; ; d = path.Dir(d) {
		if v.goDir(d) && !ignoredBelow(m.root, d) {
			return d, true
		}
		if d == m.root {
			return "", false
		}
	}
}

// goDir reports whether dir holds a .go file.
func (v *repoView) goDir(dir string) bool {
	if has, ok := v.goDirs[dir]; ok {
		return has
	}
	entries, _ := fs.ReadDir(v.fsys, dir)
	has := slices.ContainsFunc(entries, func(e fs.DirEntry) bool {
		return !e.IsDir() && strings.HasSuffix(e.Name(), ".go")
	})
	v.goDirs[dir] = has
	return has
}

// imports are the packages of its own module that the Go file p imports,
// each relative to the repository's root; none for any other file.
func (v *repoView) imports(p string) []string {
	if !strings.HasSuffix(p, ".go") {
		return nil
	}
	if imps, ok := v.imps[p]; ok {
		return imps
	}
	var imps []string
	if m := v.module(path.Dir(p)); m.ok && m.path != "" {
		for _, rel := range moduleImports(v.fsys, p, m.path) {
			imps = append(imps, path.Join(m.root, rel))
		}
	}
	v.imps[p] = imps
	return imps
}

// ignoredBelow reports whether the go tool ignores dir as a package of the
// module rooted at root: a directory between them it does not look in.
func ignoredBelow(root, dir string) bool {
	rel := relDir(root, dir)
	return rel != "." && slices.ContainsFunc(strings.Split(rel, "/"), ignoredDir)
}
