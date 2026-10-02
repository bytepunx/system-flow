package check

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

const results = `---
title: Try a cache
updated: 2026-09-01
status: active
story: %s
---

# Try a cache

## Hypothesis
A cache halves the time.

## Success measure
Median under 50 ms.

## What was done
Added one.

## Results
48 ms.

## Recommendation
Adopt it.
`

// ADR-0066: flai check validates each results document under
// design/experiments, and leaves the README and the template alone.
func TestExperimentRules(t *testing.T) {
	root := t.TempDir()
	_ = os.WriteFile(filepath.Join(root, "system-flow.yaml"), []byte("version: 1\nname: t\nkey: t\nlayout:\n  design: design\n  docs: docs\n  wip: wip\n"), 0o644)
	for _, d := range []string{"design/experiments", "docs", "wip/kanban/epics", "wip/kanban/stories", "wip/kanban/tasks", "wip/agents", "wip/archive"} {
		_ = os.MkdirAll(filepath.Join(root, d), 0o755)
	}
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	epic := mustItem(t, repo, workitem.Epic, "E", "")
	for _, nature := range []string{"experiment", "experiment", "research"} {
		if _, err := repo.Create(workitem.NewOptions{Type: workitem.Story, Title: "S", Parent: epic.ID, Nature: nature, Owner: "t", Now: now}); err != nil {
			t.Fatal(err)
		}
	}
	w := func(name, body string) {
		_ = os.WriteFile(filepath.Join(root, "design/experiments", name), []byte(body), 0o644)
	}
	w("README.md", "# experiments\n")
	w("template.md", "---\ntitle: X\n---\n")
	w("S-0001-good.md", strings.Replace(results, "%s", "S-0001", 1))
	w("S-0002-bare.md", "---\ntitle: Bare\n---\n\n# Bare\n\n## Hypothesis\n\n## Recommendation\nKeep thinking.\n")
	w("S-0003-research.md", strings.Replace(results, "%s", "S-0003", 1))
	w("S-0009-gone.md", strings.Replace(results, "%s", "S-0009", 1))
	w("S-0005-mismatch.md", strings.Replace(results, "%s", "S-0001", 1))
	w("notes.md", strings.Replace(results, "%s", "S-0002", 1))

	res, err := Run(repo, now)
	if err != nil {
		t.Fatal(err)
	}
	got := map[string][]string{}
	for _, f := range res.Findings {
		if strings.HasPrefix(f.Rule, "experiments.") || strings.HasPrefix(f.Rule, "doc.") {
			got[filepath.Base(f.Path)] = append(got[filepath.Base(f.Path)], f.Rule+": "+f.Message)
		}
	}
	for _, clean := range []string{"S-0001-good.md", "README.md", "template.md"} {
		if len(got[clean]) != 0 {
			t.Errorf("%s has no findings: %q", clean, got[clean])
		}
	}
	for file, wants := range map[string][]string{
		"S-0002-bare.md":     {"front matter needs updated", "front matter needs status", "front matter needs story", "## Hypothesis is empty", "no ## Success measure section", "no ## What was done section", "no ## Results section", "says none of adopt, adapt, drop"},
		"S-0003-research.md": {"experiments.story: S-0003 is a research story, not an experiment story"},
		"S-0009-gone.md":     {"experiments.story: story S-0009 does not exist"},
		"S-0005-mismatch.md": {"story is S-0001 but the file is named for S-0005", "experiments.duplicate: S-0001 already has its results"},
		"notes.md":           {"named <S-nnnn>-<slug>.md"},
	} {
		all := strings.Join(got[file], "\n")
		for _, want := range wants {
			if !strings.Contains(all, want) {
				t.Errorf("%s: want %q in\n%s", file, want, all)
			}
		}
	}
}
