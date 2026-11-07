package check

import (
	"os"
	"path/filepath"
	"testing"
)

// S-0406: a project whose root has neither AGENTS.md nor CLAUDE.md gets one
// instructions.missing warning on CLAUDE.md; either file alone, or both,
// clears it.
func TestInstructionsMissing(t *testing.T) {
	for _, tc := range []struct {
		name  string
		files []string
		want  int
	}{
		{"neither", nil, 1},
		{"only CLAUDE.md", []string{"CLAUDE.md"}, 0},
		{"only AGENTS.md", []string{"AGENTS.md"}, 0},
		{"both", []string{"AGENTS.md", "CLAUDE.md"}, 0},
	} {
		t.Run(tc.name, func(t *testing.T) {
			repo := planningProject(t, "")
			for _, name := range tc.files {
				if err := os.WriteFile(filepath.Join(repo.Root, name), []byte("# Map\n"), 0o644); err != nil {
					t.Fatal(err)
				}
			}
			got := findings(t, repo, "instructions.missing")
			if len(got) != tc.want {
				t.Fatalf("findings = %+v, want %d", got, tc.want)
			}
			for _, f := range got {
				if f.Level != Warning || f.Path != "CLAUDE.md" || f.Line != 1 {
					t.Errorf("finding = %+v, want a warning on CLAUDE.md:1", f)
				}
			}
		})
	}
}

// S-0406: a folder named CLAUDE.md is no map.
func TestInstructionsFolderIsNoMap(t *testing.T) {
	repo := planningProject(t, "")
	if err := os.Mkdir(filepath.Join(repo.Root, "CLAUDE.md"), 0o755); err != nil {
		t.Fatal(err)
	}
	if got := findings(t, repo, "instructions.missing"); len(got) != 1 {
		t.Errorf("findings = %+v, want one", got)
	}
}
