package workitem

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/mdlint"
)

const projectLint = "default: true\nMD013: false\nMD022:\n  lines_below: 0\nMD024:\n  siblings_only: true\nMD025:\n  front_matter_title: \"\"\nMD032: false\nMD041: false\n"

func lintProject(t *testing.T) *Repo {
	t.Helper()
	root := t.TempDir()
	_ = os.WriteFile(filepath.Join(root, "system-flow.yaml"), []byte("version: 1\nname: t\nkey: t\nlayout:\n  design: design\n  docs: docs\n  wip: wip\n"), 0o644)
	_ = os.WriteFile(filepath.Join(root, ".markdownlint.yaml"), []byte(projectLint), 0o644)
	for _, d := range []string{"wip/kanban/epics", "wip/kanban/stories", "wip/kanban/tasks", "wip/agents", "wip/archive"} {
		_ = os.MkdirAll(filepath.Join(root, d), 0o755)
	}
	r, err := Open(root)
	if err != nil {
		t.Fatal(err)
	}
	return r
}

func TestCleanTitle(t *testing.T) {
	for in, want := range map[string]string{
		"Fix the bug.":           "Fix the bug",
		"  Spaced   out  ":       "Spaced out",
		"Why does it fail?":      "Why does it fail?",
		"Ends badly!!":           "Ends badly",
		"Lists: a, b, and c;":    "Lists: a, b, and c",
		"全角。":                    "全角",
		"...":                    "",
		"S-0179 keeps its ID.  ": "S-0179 keeps its ID",
	} {
		if got := CleanTitle(in); got != want {
			t.Errorf("CleanTitle(%q) = %q, want %q", in, got, want)
		}
	}
}

// S-0179: what Create and LogStream write cannot bring markdown the
// project's lint rejects, and a title loses the punctuation MD026 rejects.
func TestWritesTheLintRejectsAreRefused(t *testing.T) {
	r := lintProject(t)
	now := time.Date(2026, 10, 1, 9, 0, 0, 0, time.UTC)
	it, err := r.Create(NewOptions{Type: Story, Title: "Ends with a period.", Owner: "alex", Now: now})
	if err != nil {
		t.Fatal(err)
	}
	data, _ := os.ReadFile(it.Path)
	if it.Title != "Ends with a period" || !strings.Contains(string(data), "\n# S-0001 Ends with a period\n") || !strings.Contains(string(data), "title: Ends with a period\n") {
		t.Errorf("title written with its period:\n%s", data)
	}

	_, err = r.Create(NewOptions{Type: Story, Title: "Bad body", Owner: "alex", Now: now, Body: "## Goal\n\nText\n## Acceptance criteria\n- [ ] x\n"})
	var le *mdlint.Error
	if !errors.As(err, &le) || len(le.Findings) != 1 || le.Findings[0].Rule != "MD022" || !strings.Contains(err.Error(), "wip/kanban/stories/S-0002-bad-body.md") {
		t.Fatalf("a body the lint rejects is refused with the rule and line: %v", err)
	}
	if m, _ := filepath.Glob(filepath.Join(r.WipDir(), "kanban/stories/S-0002-*")); len(m) != 0 {
		t.Errorf("nothing is written: %v", m)
	}

	if _, err := r.OpenStream(it, StreamOptions{Agent: "a", Now: now}); err != nil {
		t.Fatal(err)
	}
	if _, err := r.LogStream(it.ID, "Started.", StreamOptions{Now: now}); err != nil {
		t.Fatalf("a clean entry is logged: %v", err)
	}
	was, _ := os.ReadFile(r.NarrativePath(it.ID))
	_, err = r.LogStream(it.ID, "Did a thing.\n\n**Looks like a heading**", StreamOptions{Now: now.Add(time.Second)})
	if !errors.As(err, &le) || le.Findings[0].Rule != "MD036" {
		t.Fatalf("an entry the lint rejects is refused: %v", err)
	}
	if got, _ := os.ReadFile(r.NarrativePath(it.ID)); string(got) != string(was) {
		t.Error("a refused entry leaves the narrative as it was")
	}

	// without a configuration nothing is linted
	_ = os.Remove(filepath.Join(r.Root, ".markdownlint.yaml"))
	if _, err := r.LogStream(it.ID, "Did a thing.\n\n**Looks like a heading**", StreamOptions{Now: now.Add(time.Second)}); err != nil {
		t.Errorf("no configuration, no lint: %v", err)
	}
}
