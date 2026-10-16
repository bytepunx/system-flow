// Package components catalogues the components of a project's Svelte
// sub-projects (S-0437, E-0026): Find walks a sub-project for its .svelte
// files by include and exclude globs, and ReadProps reads each one's props
// with props.mjs, a script flai embeds and runs with the project's own node
// and svelte/compiler, so that the props are read as the project's compiler
// reads them. No component is ever run.
//
// The package takes its globs as arguments and does not read the manifest;
// the caller passes a sub-project's components key, or nothing for the
// defaults.
package components

import (
	"fmt"
	"io/fs"
	"path"
	"path/filepath"
	"slices"
	"strings"
)

// DefaultInclude is the include glob used when none is given: the
// sub-project's library, which leaves its routes out.
var DefaultInclude = []string{"src/lib/**/*.svelte"}

// Component is one .svelte file of a sub-project and the props it takes.
type Component struct {
	// Project is the sub-project's path relative to the project root, in
	// slash form.
	Project string `json:"project"`
	// Path is the file's path relative to the sub-project, in slash form.
	Path string `json:"path"`
	// Name is the file's base name without .svelte.
	Name string `json:"name"`
	// Props are the props it takes, in the order its $props() destructures
	// them, then the rest its props type declares; empty when it has no
	// $props() or they could not be read.
	Props []Prop `json:"props"`
	// Error says why its props could not be read, such as a parse error;
	// empty when they were.
	Error string `json:"error,omitempty"`
}

// Find is the components of the sub-project at project, relative to root:
// each .svelte file whose sub-project-relative path matches an include glob
// and no exclude glob, sorted by path, without props. No include means
// DefaultInclude. node_modules and folders whose names begin with a dot are
// skipped. A glob is relative to the sub-project, as CheckGlob describes.
func Find(root, project string, include, exclude []string) ([]Component, error) {
	if len(include) == 0 {
		include = DefaultInclude
	}
	for _, g := range slices.Concat(include, exclude) {
		if err := CheckGlob(g); err != nil {
			return nil, fmt.Errorf("cannot find the components of %s: %w", project, err)
		}
	}
	dir := filepath.Join(root, filepath.FromSlash(project))
	slashProject := filepath.ToSlash(filepath.Clean(project))
	var out []Component
	err := filepath.WalkDir(dir, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if p != dir && (d.Name() == "node_modules" || strings.HasPrefix(d.Name(), ".")) {
				return filepath.SkipDir
			}
			return nil
		}
		if !d.Type().IsRegular() || !strings.HasSuffix(d.Name(), ".svelte") {
			return nil
		}
		rel, err := filepath.Rel(dir, p)
		if err != nil {
			return err
		}
		rel = filepath.ToSlash(rel)
		if !matchAny(include, rel) || matchAny(exclude, rel) {
			return nil
		}
		out = append(out, Component{
			Project: slashProject,
			Path:    rel,
			Name:    strings.TrimSuffix(d.Name(), ".svelte"),
			Props:   []Prop{},
		})
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("cannot walk sub-project %s for its components: %w; check its path in the manifest", project, err)
	}
	slices.SortFunc(out, func(a, b Component) int { return strings.Compare(a.Path, b.Path) })
	return out, nil
}

// CheckGlob says what is wrong with a components glob, or nil when it is
// valid: a path relative to the sub-project, segments separated by /, each a
// path.Match pattern or ** alone, which stands for zero or more whole
// segments. A glob with no glob character (* ? [ \) matches the path itself
// and everything below it, as claims.shared's patterns do.
func CheckGlob(g string) error {
	switch {
	case strings.TrimSpace(g) == "":
		return fmt.Errorf("glob %q is empty; write a path relative to the sub-project, such as src/lib/**/*.svelte", g)
	case strings.HasPrefix(g, "/") || strings.HasPrefix(g, `\`) || len(g) > 1 && g[1] == ':':
		return fmt.Errorf("glob %q is an absolute path; write it relative to the sub-project, such as src/lib/**/*.svelte", g)
	}
	for _, seg := range strings.Split(g, "/") {
		switch {
		case seg == "" || seg == "." || seg == "..":
			return fmt.Errorf("glob %q has an empty, . or .. segment; write it with single slashes, inside the sub-project, such as src/lib/**/*.svelte", g)
		case seg != "**" && strings.Contains(seg, "**"):
			return fmt.Errorf("glob %q has ** inside a segment; make ** a segment of its own, such as src/**/x.svelte, or use * within one", g)
		}
		if _, err := path.Match(seg, ""); err != nil {
			return fmt.Errorf("glob %q does not parse, such as an unclosed [ or a trailing \\; fix it", g)
		}
	}
	return nil
}

// matchAny reports whether any glob matches rel, a slash path relative to
// the sub-project. The globs are valid.
func matchAny(globs []string, rel string) bool {
	segs := strings.Split(rel, "/")
	for _, g := range globs {
		gs := strings.Split(g, "/")
		if !strings.ContainsAny(g, `*?[\`) {
			gs = append(gs, "**")
		}
		if matchSegments(gs, segs) {
			return true
		}
	}
	return false
}

// matchSegments reports whether the glob's segments match the path's: **
// matches zero or more whole segments, any other segment one by path.Match.
func matchSegments(g, s []string) bool {
	if len(g) == 0 {
		return len(s) == 0
	}
	if g[0] == "**" {
		return matchSegments(g[1:], s) || len(s) > 0 && matchSegments(g, s[1:])
	}
	if len(s) == 0 {
		return false
	}
	ok, _ := path.Match(g[0], s[0])
	return ok && matchSegments(g[1:], s[1:])
}
