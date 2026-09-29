package context

import (
	"os"
	"path/filepath"
	"slices"
	"testing"
)

func TestLoadSectionsCutsDesignAndDocsAsThePackDoes(t *testing.T) {
	root := t.TempDir()
	files := map[string]string{
		"design/system/cli.md":         "---\ntitle: The CLI\n---\n\n# The CLI\n\nIntro.\n\n## Commands\n\n### flai prime\n\nPrints the pack.\n",
		"design/system/README.md":      "# Index\n",
		"design/adrs/0000-template.md": "# Template\n",
		"design/adrs/0001-first.md":    "---\nid: ADR-0001\ntitle: First\n---\n\n# ADR-0001 First\n\n## Decision\n\nWe decide.\n",
		"design/conventions/git.md":    "---\ntitle: Git\n---\n\n# Git\n\n## Rules\n\n- Commit.\n",
		"docs/users/flai.md":           "# flai\n\n## Install\n\nRun it.\n",
		"docs/.hidden/secret.md":       "# Hidden\n",
		"design/system/notes.txt":      "not markdown",
	}
	for p, c := range files {
		abs := filepath.Join(root, p)
		if err := os.MkdirAll(filepath.Dir(abs), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(abs, []byte(c), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	secs, err := LoadSections(root, filepath.Join(root, "design"), filepath.Join(root, "docs"), filepath.Join(root, "missing"))
	if err != nil {
		t.Fatal(err)
	}
	var got []string
	for _, s := range secs {
		got = append(got, s.Path+" § "+s.Heading)
	}
	want := []string{
		"design/adrs/0001-first.md § ADR-0001 First", "design/adrs/0001-first.md § Decision",
		"design/conventions/git.md § Git", "design/conventions/git.md § Rules",
		"design/system/cli.md § The CLI", "design/system/cli.md § Commands", "design/system/cli.md § Commands › flai prime",
		"docs/users/flai.md § flai", "docs/users/flai.md § Install",
	}
	if !slices.Equal(got, want) {
		t.Errorf("sections\n got %q\nwant %q", got, want)
	}
	if s := secs[6]; s.Title != "The CLI" || s.Line != 11 || s.Level != 3 || s.Text != "### flai prime\n\nPrints the pack.\n" {
		t.Errorf("section %+v", s)
	}
}
