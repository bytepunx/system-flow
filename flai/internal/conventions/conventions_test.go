package conventions

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
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

// ADR-0068: a convention's roles are every agent that reads it; each role
// flai knows is accepted and any other is a warning.
func TestRoles(t *testing.T) {
	root := t.TempDir()
	repo := &workitem.Repo{Root: root}
	repo.Manifest.Layout = map[string]string{"design": "design"}
	dir := Dir(repo)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	write := func(name, s string) {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(s), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	write("README.md", "# Conventions\n\n- [safety.md](safety.md)\n- [git.md](git.md)\n- [plan.md](plan.md)\n")
	file := func(title string, order int, roles string) string {
		return fmt.Sprintf("---\ntitle: %s\nupdated: 2026-10-01\naudience: agent\norder: %d\nstatus: active\n%s---\n\n# %s\n\n%s\n\n## Project additions\n", title, order, roles, title, Marker)
	}
	write("safety.md", file("Safety", 80, "roles: [explore, verify]\n"))
	write("git.md", file("Git", 70, "roles: [verify, write]\n"))
	write("plan.md", file("Plan", 90, "roles: [story, orchestrator, planner, analyzer]\n"))
	set, errs, err := Load(repo)
	if err != nil {
		t.Fatal(err)
	}
	if len(set.Files) != 3 || strings.Join(set.Files[0].Roles, ",") != "verify,write" || strings.Join(set.Files[1].Roles, ",") != "explore,verify" {
		t.Fatalf("roles: %+v", set.Files)
	}
	got := set.Validate(errs)
	if len(got) != 1 || got[0].Rule != "conventions.roles" || got[0].Level != "warning" || !strings.Contains(got[0].Message, `"write"`) {
		t.Errorf("findings: %+v", got)
	}
}

// ADR-0068: a convention without roles is read by every agent, one with
// roles by the agents they list.
func TestReadBy(t *testing.T) {
	for _, c := range []struct {
		roles []string
		role  string
		want  bool
	}{
		{nil, RoleStory, true},
		{nil, RoleVerify, true},
		{[]string{RoleStory}, RoleStory, true},
		{[]string{RoleStory}, RoleExplore, false},
		{[]string{RoleVerify}, RoleStory, false},
		{[]string{RoleExplore, RoleVerify}, RoleVerify, true},
	} {
		if got := (File{Roles: c.roles}).ReadBy(c.role); got != c.want {
			t.Errorf("roles %v read by %s: %v, want %v", c.roles, c.role, got, c.want)
		}
	}
}
