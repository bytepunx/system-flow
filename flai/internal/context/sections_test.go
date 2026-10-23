package context

import (
	"os"
	"path/filepath"
	"slices"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/workitem"
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

// TestSearchSectionsFindsADigestEntry: an entry under design/digest is a
// document of the design folder, so doc_search and flai doc search find it
// by a word of its title with no code of its own, and leave the folder's
// README out as they leave every README out (ADR-0191).
func TestSearchSectionsFindsADigestEntry(t *testing.T) {
	root := t.TempDir()
	entry := "design/digest/D-0001-the-operator-keeps-the-forge-token-in-the-host-s-environment-never-in-the.md"
	title := "The operator keeps the forge token in the host's environment, never in the repository"
	files := map[string]string{
		"system-flow.yaml": "version: 1\nname: t\nkey: t\nlayout:\n  design: design\n  docs: docs\n  wip: wip\n",
		entry: "---\nid: D-0001\ntitle: " + title + "\nkind: answer\nsources: [TH-0001]\ntopics: [git]\npaths: []\ndate: 2026-10-23\nstatus: active\nabsorbed_by: \"\"\nasked_again: []\ncited: []\n---\n\n" +
			"# D-0001 " + title + "\n\nAsked on TH-0001; the operator answered that the host holds it.\n",
		"design/digest/README.md": "# Digest\n\nOne file per fact a story learned, such as where the forge token lives.\n",
		"design/system/cli.md":    "---\ntitle: CLI\n---\n\n# CLI\n\nThe command line.\n",
	}
	for p, c := range files {
		abs := filepath.Join(root, filepath.FromSlash(p))
		if err := os.MkdirAll(filepath.Dir(abs), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(abs, []byte(c), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	for _, d := range []string{"wip/kanban/epics", "wip/kanban/stories", "wip/kanban/tasks", "wip/agents"} {
		if err := os.MkdirAll(filepath.Join(root, d), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	hits, err := SearchSections(repo, "forge", 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(hits) != 1 {
		t.Fatalf("hits %+v, want the entry alone", hits)
	}
	if h := hits[0]; h.Path != entry || h.Heading != "D-0001 "+title || h.Title != title || h.Line != 15 {
		t.Errorf("hit %+v", h)
	}
}
