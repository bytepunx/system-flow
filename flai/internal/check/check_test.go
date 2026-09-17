package check

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

var now = time.Date(2026, 9, 1, 12, 0, 0, 0, time.UTC)

func TestGoodFixtureIsClean(t *testing.T) {
	repo, err := workitem.Open("../metrics/testdata/good")
	if err != nil {
		t.Fatal(err)
	}
	res, err := Run(repo, now)
	if err != nil {
		t.Fatal(err)
	}
	for _, f := range res.Findings {
		t.Errorf("unexpected finding: %+v", f)
	}
	if !res.OK(true) || res.Items != 10 {
		t.Fatalf("result: %+v", res)
	}
}

func TestBadFixtureFindings(t *testing.T) {
	repo, err := workitem.Open("testdata/bad")
	if err != nil {
		t.Fatal(err)
	}
	res, err := Run(repo, now)
	if err != nil {
		t.Fatal(err)
	}
	want := map[string]string{ // rule -> substring of path
		"story.tasks":         "S-001",
		"story.criteria":      "S-001",
		"narrative.missing":   "S-001",
		"item.parent-missing": "T-001",
		"item.stream":         "T-001",
		"item.done-children":  "E-001",
		"item.archive":        "E-001",
		"item.filename":       "S-003-wrong-name",
		"item.chronology":     "S-003-wrong-name",
		"item.front-matter":   "S-003-wrong-name",
		"item.sequence":       "S-003-wrong-name",
		"item.parent-list":    "E-001",
		"board.order":         "board.md",
		"narrative.stream":    "S-007",
		"narrative.section":   "S-007",
		"narrative.index":     "index.md",
		"doc.front-matter":    "nofm.md",
		"doc.updated":         "docs/users/index.md",
		"adr.id":              "0001-first.md",
	}
	got := map[string][]string{}
	for _, f := range res.Findings {
		got[f.Rule] = append(got[f.Rule], f.Path)
		if f.Line < 1 {
			t.Errorf("finding without line: %+v", f)
		}
	}
	for rule, path := range want {
		found := false
		for _, p := range got[rule] {
			if strings.Contains(p, path) {
				found = true
			}
		}
		if !found {
			t.Errorf("expected %s on %s, got %v", rule, path, got[rule])
		}
	}
	if res.OK(false) {
		t.Error("bad fixture must fail")
	}
	// line numbers point at the offending key
	for _, f := range res.Findings {
		if f.Rule == "item.parent-missing" && f.Line != 7 {
			t.Errorf("parent-missing should point at the parent line (7), got %d", f.Line)
		}
	}
}

func TestMonorepoIsClean(t *testing.T) {
	if testing.Short() {
		t.Skip("integration: reads the monorepo")
	}
	repo, err := workitem.Open("../../..")
	if err != nil {
		t.Skip("monorepo not present")
	}
	res, err := Run(repo, time.Now())
	if err != nil {
		t.Fatal(err)
	}
	// Warnings (a done story awaiting archive, a WIP breach) are board state,
	// not code defects; the strict check job gates them. Errors fail here.
	for _, f := range res.Findings {
		if f.Level == Error {
			t.Errorf("%s:%d: %s: %s: %s", f.Path, f.Line, f.Level, f.Rule, f.Message)
		} else {
			t.Logf("warning: %s:%d: %s: %s", f.Path, f.Line, f.Rule, f.Message)
		}
	}
}

func TestCacheMustBeIgnored(t *testing.T) {
	root := t.TempDir()
	_ = os.WriteFile(filepath.Join(root, "system-flow.yaml"), []byte("version: 1\nname: t\nkey: t\nlayout:\n  design: design\n  docs: docs\n  wip: wip\n"), 0o644)
	for _, d := range []string{"design/adrs", "design/system", "design/tech", "design/conventions", "docs", "wip/kanban/epics", "wip/kanban/stories", "wip/kanban/tasks", "wip/agents", "wip/archive", ".flai-cache"} {
		_ = os.MkdirAll(filepath.Join(root, d), 0o755)
	}
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	has := func(rule string) bool {
		res, err := Run(repo, now)
		if err != nil {
			t.Fatal(err)
		}
		for _, f := range res.Findings {
			if f.Rule == rule {
				return true
			}
		}
		return false
	}
	if !has("layout.gitignore") {
		t.Error("expected layout.gitignore without a .gitignore")
	}
	_ = os.WriteFile(filepath.Join(root, ".gitignore"), []byte("bin/\n"), 0o644)
	if !has("layout.gitignore") {
		t.Error("expected layout.gitignore when .flai-cache is not listed")
	}
	_ = os.WriteFile(filepath.Join(root, ".gitignore"), []byte("bin/\n.flai-cache/\n"), 0o644)
	if has("layout.gitignore") {
		t.Error("ignored cache should pass")
	}
}

func TestTouchesOverlap(t *testing.T) {
	root := t.TempDir()
	_ = os.WriteFile(filepath.Join(root, "system-flow.yaml"), []byte("version: 1\nname: t\nkey: t\nlayout:\n  design: design\n  docs: docs\n  wip: wip\n"), 0o644)
	for _, d := range []string{"design/adrs", "design/system", "design/tech", "design/conventions", "docs", "wip/kanban/epics", "wip/kanban/stories", "wip/kanban/tasks", "wip/agents", "wip/archive"} {
		_ = os.MkdirAll(filepath.Join(root, d), 0o755)
	}
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	e := mustItem(t, repo, workitem.Epic, "E", "")
	s1 := mustItem(t, repo, workitem.Story, "One", e.ID)
	s2 := mustItem(t, repo, workitem.Story, "Two", e.ID)
	for _, s := range []*workitem.Item{s1, s2} {
		s.Status = workitem.InProgress
		s.Touches = []string{"flaiover/src/lib"}
		if err := repo.Save(s); err != nil {
			t.Fatal(err)
		}
	}
	s2.Touches = []string{"flaiover/src/lib/server/auth.ts"}
	_ = repo.Save(s2)
	res, err := Run(repo, now)
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for _, f := range res.Findings {
		if f.Rule == "wip.overlap" && strings.Contains(f.Message, "S-0001 touches flaiover/src/lib") {
			found = true
		}
	}
	if !found {
		t.Errorf("expected wip.overlap: %+v", res.Findings)
	}
	s2.Touches = []string{"docs"}
	_ = repo.Save(s2)
	res, _ = Run(repo, now)
	for _, f := range res.Findings {
		if f.Rule == "wip.overlap" {
			t.Errorf("disjoint paths must not overlap: %+v", f)
		}
	}
}

func mustItem(t *testing.T, r *workitem.Repo, typ, title, parent string) *workitem.Item {
	t.Helper()
	it, err := r.Create(workitem.NewOptions{Type: typ, Title: title, Parent: parent, Owner: "t", Now: now})
	if err != nil {
		t.Fatal(err)
	}
	return it
}
