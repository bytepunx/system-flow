package workflow

import (
	"maps"
	"slices"
	"testing"
)

// S-0806: Fields lists the keys of the front matter and of each mapping in
// it from the types the file decodes into, ADR-0312's keys among them, so
// that front-matter-fields.txt carries them and a release that adds them
// raises flai.minimum.
func TestFieldsAreTheDecodersKeys(t *testing.T) {
	f := Fields()
	names := slices.Sorted(maps.Keys(f))
	if want := []string{"workflow", "workflow.state", "workflow.state.runs", "workflow.timer", "workflow.timer.point", "workflow.timer.point.edge", "workflow.transition", "workflow.transition.rule", "workflow.transition.verify", "workflow.verify_at"}; !slices.Equal(names, want) {
		t.Errorf("mappings %v, want %v", names, want)
	}
	for name, want := range map[string][]string{
		"workflow":                   {"states", "status", "timers", "title", "transitions", "updated", "verify_at"},
		"workflow.state":             {"aliases", "initial", "kind", "limit", "name", "order", "runs", "title", "types"},
		"workflow.state.runs":        {"kluster", "suites"},
		"workflow.transition":        {"by", "effects", "from", "rules", "to", "types", "verify"},
		"workflow.transition.rule":   {"rule", "severity"},
		"workflow.transition.verify": {"severity", "suites"},
		"workflow.timer":             {"counts", "from", "name", "sum", "title", "to", "types", "until"},
		"workflow.timer.point":       {"first", "last"},
		"workflow.timer.point.edge":  {"kind"},
		"workflow.verify_at":         {"batch", "ci", "publish"},
	} {
		if !slices.Equal(f[name], want) {
			t.Errorf("%s: %v, want %v", name, f[name], want)
		}
	}
}
