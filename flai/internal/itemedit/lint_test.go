package itemedit

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/docedit"
	"github.com/bytepunx/system-flow/flai/internal/itemnew"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// S-0179: with a markdownlint configuration, flai check lints wip, so an
// edit or a new item whose body the lint rejects is refused with the rule
// and the line, and nothing is written.
func TestBodiesTheLintRejectsAreRefused(t *testing.T) {
	root := t.TempDir()
	_ = os.WriteFile(filepath.Join(root, "system-flow.yaml"), []byte("version: 1\nname: t\nkey: t\nlayout:\n  design: design\n  docs: docs\n  wip: wip\n"), 0o644)
	_ = os.WriteFile(filepath.Join(root, ".markdownlint.yaml"), []byte("default: true\nMD013: false\nMD022:\n  lines_below: 0\nMD024:\n  siblings_only: true\nMD025:\n  front_matter_title: \"\"\nMD032: false\nMD041: false\n"), 0o644)
	for _, d := range []string{"wip/kanban/epics", "wip/kanban/stories", "wip/kanban/tasks", "wip/agents", "wip/archive"} {
		_ = os.MkdirAll(filepath.Join(root, d), 0o755)
	}
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	now := time.Date(2026, 10, 1, 9, 0, 0, 0, time.UTC)
	it, err := repo.Create(workitem.NewOptions{Type: workitem.Story, Title: "Lint me", Nature: "feature", Owner: "alex", Now: now})
	if err != nil {
		t.Fatal(err)
	}
	was, _ := os.ReadFile(it.Path)

	body := "## Goal\n\nDo it.\n\n**Not a heading**\n\n## Acceptance criteria\n- [ ] Done\n\n## Tasks\n\n## Notes\n"
	_, err = Apply(repo, nil, it.ID, Change{Body: &body}, Options{By: "alex", NoCommit: true, Now: now})
	r, ok := docedit.IsRefused(err)
	if !ok {
		t.Fatalf("want a refusal, got %v", err)
	}
	if len(r.Findings) != 1 || r.Findings[0].Rule != "markdown.MD036" || r.Findings[0].Line != 19 {
		t.Fatalf("findings: %+v", r.Findings)
	}
	if now, _ := os.ReadFile(it.Path); string(now) != string(was) {
		t.Error("a refused edit leaves the file as it was")
	}

	body = strings.Replace(body, "**Not a heading**", "Not a heading.", 1)
	if _, err := Apply(repo, nil, it.ID, Change{Body: &body}, Options{By: "alex", NoCommit: true, Now: now}); err != nil {
		t.Fatalf("a clean body is kept: %v", err)
	}

	_, err = itemnew.Create(repo, nil, itemnew.Options{New: workitem.NewOptions{Type: workitem.Story, Title: "Another", Nature: "feature", Owner: "alex", Now: now,
		Body: "## Goal\n\nTrailing space. \n\n## Acceptance criteria\n- [ ] Done\n\n## Tasks\n\n## Notes\n"}})
	r, ok = docedit.IsRefused(err)
	if !ok || len(r.Findings) != 1 || r.Findings[0].Rule != "markdown.MD009" {
		t.Fatalf("a new item with a finding is refused: %v %+v", err, r)
	}
	if m, _ := filepath.Glob(filepath.Join(root, "wip/kanban/stories", "S-0002-*")); len(m) != 0 {
		t.Errorf("nothing was created: %v", m)
	}
}
