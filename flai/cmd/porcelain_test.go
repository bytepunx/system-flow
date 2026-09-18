package cmd

import (
	"reflect"
	"testing"
)

func TestPorcelainPaths(t *testing.T) {
	// The first line has lost its leading space, as the runner trims output.
	out := "M wip/agents/index.md\n M wip/kanban/board.md\n?? flai/cmd/new.go\nA  docs/a.md\nRM old.md -> docs/new.md\n M \"with space.md\""
	want := []string{"wip/agents/index.md", "wip/kanban/board.md", "flai/cmd/new.go", "docs/a.md", "docs/new.md", "with space.md"}
	if got := porcelainPaths(out); !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
	if got := porcelainPaths(""); len(got) != 0 {
		t.Errorf("empty: %v", got)
	}
}
