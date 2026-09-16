package conventions

import (
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

func rules(fs []Finding) map[string]int {
	m := map[string]int{}
	for _, f := range fs {
		m[f.Rule+":"+f.Level]++
	}
	return m
}

func TestGoodSet(t *testing.T) {
	repo, err := workitem.Open("../metrics/testdata/good")
	if err != nil {
		t.Fatal(err)
	}
	set, errs, err := Load(repo)
	if err != nil {
		t.Fatal(err)
	}
	if set.Missing || len(set.Files) != 2 || set.Files[0].Name != "session-start.md" || set.Files[1].Order != 20 || set.README == "" {
		t.Fatalf("set: %+v", set)
	}
	if f := set.Validate(errs); len(f) != 0 {
		t.Errorf("unexpected findings: %+v", f)
	}
}

func TestBadSet(t *testing.T) {
	repo, _ := workitem.Open("../check/testdata/bad")
	set, errs, err := Load(repo)
	if err != nil {
		t.Fatal(err)
	}
	got := rules(set.Validate(errs))
	for rule, want := range map[string]int{
		"conventions.marker:error":       2, // twice in session-start, none in dup-order
		"conventions.marker:warning":     1, // long.md has no project additions
		"conventions.order:error":        1, // dup-order reuses 10
		"conventions.front-matter:error": 1, // audience human
		"conventions.length:warning":     1,
		"conventions.index:error":        3, // dup-order and long not listed, missing.md linked
		"conventions.index:warning":      1, // session-start listed twice
	} {
		if got[rule] != want {
			t.Errorf("%s: got %d want %d (all: %v)", rule, got[rule], want, got)
		}
	}
}

func TestMissingFolder(t *testing.T) {
	repo := &workitem.Repo{Root: t.TempDir()}
	repo.Manifest.Layout = map[string]string{"design": "design"}
	set, errs, err := Load(repo)
	if err != nil || !set.Missing {
		t.Fatalf("expected missing: %+v %v", set, err)
	}
	f := set.Validate(errs)
	if len(f) != 1 || f[0].Rule != "conventions.missing" || f[0].Level != "warning" {
		t.Errorf("missing folder findings: %+v", f)
	}
}
