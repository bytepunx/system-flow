package experiment

import (
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

var m = manifest.Manifest{Layout: map[string]string{"design": "notes", "docs": "docs", "wip": "wip"}}

func TestExpectedFollowsTheLayoutAndTheStorysName(t *testing.T) {
	it := &workitem.Item{ID: "S-0007", Type: workitem.Story, Nature: Nature, Title: "Try a cache"}
	if got := Expected(m, it); got != "notes/experiments/S-0007-try-a-cache.md" {
		t.Errorf("expected document: %q", got)
	}
	if !Needs(it) {
		t.Error("an experiment story needs its results")
	}
	for _, other := range []*workitem.Item{
		{ID: "S-0008", Type: workitem.Story, Nature: "research"},
		{ID: "E-0001", Type: workitem.Epic, Nature: Nature},
	} {
		if Needs(other) {
			t.Errorf("%s %s needs no results document", other.Type, other.Nature)
		}
	}
	if why := Missing(m, it, "main"); !strings.Contains(why, "notes/experiments/S-0007-try-a-cache.md") || !strings.Contains(why, "ADR-0066") {
		t.Errorf("the reason names the document: %q", why)
	}
}

func TestFindMatchesTheStorysIDWhateverTheSlug(t *testing.T) {
	cases := map[string]bool{
		"S-0007-try-a-cache.md": true,
		"S-0007-renamed.md":     true,
		"S-0007.md":             true,
		"S-00071-other.md":      false,
		"S-0007-notes.txt":      false,
		"README.md":             false,
	}
	for name, want := range cases {
		if got := Names(name, "S-0007"); got != want {
			t.Errorf("Names(%q) = %v, want %v", name, got, want)
		}
	}
	if got, ok := Find([]string{"README.md", "S-0006-x.md", "S-0007-y.md"}, "S-0007"); !ok || got != "S-0007-y.md" {
		t.Errorf("find: %q %v", got, ok)
	}
}
