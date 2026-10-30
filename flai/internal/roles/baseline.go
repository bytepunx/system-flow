package roles

import (
	"embed"
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"slices"
)

// baselineFS holds the template's baseline role files, byte copies of
// template/root/design/workflows/roles/<name>.md, so that flai holds a
// built-in role to its rules in a project that lacks the file (ADR-0207).
//
//go:embed baseline/*.md
var baselineFS embed.FS

// CodeBaseline is the code flai check warns under when the project lacks the
// file of a role flai carries a baseline for, and so runs the baseline.
const CodeBaseline = "roles.baseline"

// baselineNames are the baseline roles' names, in ADR-0207's order.
var baselineNames = []string{"planner", "orchestrator", "analyzer", "curator"}

// Baseline is the baseline roles flai carries, each parsed as if it were
// <design>/workflows/roles/<name>.md.
func Baseline(design string) []Role {
	out := make([]Role, 0, len(baselineNames))
	for _, name := range baselineNames {
		out = append(out, baselineRole(design, name))
	}
	return out
}

// baselineRole is the embedded baseline role file name.md, parsed as if it
// were <design>/workflows/roles/<name>.md.
func baselineRole(design, name string) Role {
	data, err := baselineFS.ReadFile("baseline/" + name + ".md")
	if err != nil {
		panic(fmt.Sprintf("the baseline role file %s.md is not embedded: %v", name, err))
	}
	r, p := Parse(path.Join(Dir(design), name+".md"), string(data))
	if p != nil {
		panic("the embedded baseline role file does not parse: " + p.String())
	}
	return r
}

// Builtin is the role named name that flai runs in the project at root: its
// own <design>/workflows/roles/<name>.md when that file exists and parses,
// else the baseline flai carries, with fromBaseline saying which. A file
// that does not parse is flai check's roles.file error, not this one's.
func Builtin(root, design, name string) (r Role, fromBaseline bool, err error) {
	rel := path.Join(Dir(design), name+".md")
	data, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(rel)))
	switch {
	case err == nil:
		if r, p := Parse(rel, string(data)); p == nil {
			return r, false, nil
		}
	case !errors.Is(err, os.ErrNotExist):
		return Role{}, false, fmt.Errorf("reading role file %s: %w", rel, err)
	}
	if !slices.Contains(baselineNames, name) {
		return Role{}, false, fmt.Errorf("role %s: %s is missing or does not parse, and flai carries no baseline for it; add the file, or run flai check to see what is wrong with it", name, rel)
	}
	return baselineRole(design, name), true, nil
}

// BaselineFor is the baseline role whose role, its verb, is role, such as
// plan, and whether there is one.
func BaselineFor(design, role string) (Role, bool) {
	for _, r := range Baseline(design) {
		if r.Role == role {
			return r, true
		}
	}
	return Role{}, false
}
