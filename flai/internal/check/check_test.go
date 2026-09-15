package check

import (
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
	for _, f := range res.Findings {
		t.Errorf("%s:%d: %s: %s: %s", f.Path, f.Line, f.Level, f.Rule, f.Message)
	}
}
