package roles

import (
	"embed"
	"fmt"
	"path"
)

// baselineFS holds the template's baseline role files, byte copies of
// template/root/design/workflows/roles/<name>.md, so that flai holds a
// built-in role to its rules in a project that lacks the file (ADR-0207).
//
//go:embed baseline/*.md
var baselineFS embed.FS

// baselineNames are the baseline roles' names, in ADR-0207's order.
var baselineNames = []string{"planner", "orchestrator", "analyzer", "curator"}

// Baseline is the baseline roles flai carries, each parsed as if it were
// <design>/workflows/roles/<name>.md.
func Baseline(design string) []Role {
	out := make([]Role, 0, len(baselineNames))
	for _, name := range baselineNames {
		data, err := baselineFS.ReadFile("baseline/" + name + ".md")
		if err != nil {
			panic(fmt.Sprintf("the baseline role file %s.md is not embedded: %v", name, err))
		}
		r, p := Parse(path.Join(Dir(design), name+".md"), string(data))
		if p != nil {
			panic("the embedded baseline role file does not parse: " + p.String())
		}
		out = append(out, r)
	}
	return out
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
