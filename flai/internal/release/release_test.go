package release

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/execx"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

func TestVersionsAndLevels(t *testing.T) {
	v, ok := ParseVersion("v1.2.3")
	if !ok || v.String() != "1.2.3" || v.Bumped(Major).String() != "2.0.0" || v.Bumped(Minor).String() != "1.3.0" || v.Bumped(Patch).String() != "1.2.4" {
		t.Errorf("version: %v %v", v, ok)
	}
	if _, ok := ParseVersion("flai/v1.0"); ok {
		t.Error("partial version accepted")
	}
	cases := map[*workitem.Item]string{
		{Type: workitem.Epic, Nature: "feature"}:      Major,
		{Type: workitem.Story, Nature: "feature"}:     Minor,
		{Type: workitem.Story, Nature: "remediation"}: Patch,
		{Type: workitem.Story, Nature: "improvement"}: Patch,
	}
	for it, want := range cases {
		if got, err := LevelFor(it); err != nil || got != want {
			t.Errorf("%s/%s: %s %v", it.Type, it.Nature, got, err)
		}
	}
	if _, err := LevelFor(&workitem.Item{ID: "S-9", Type: workitem.Story, Nature: "research"}); err == nil {
		t.Error("research must not release")
	}
}

func gitRepo(t *testing.T) (string, execx.Runner) {
	t.Helper()
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}
	root := t.TempDir()
	r := execx.System{}
	run := func(args ...string) {
		if out, err := r.Run(root, "git", args...); err != nil {
			t.Fatalf("git %v: %v %s", args, err, out)
		}
	}
	w := func(rel, body string) {
		p := filepath.Join(root, rel)
		_ = os.MkdirAll(filepath.Dir(p), 0o755)
		_ = os.WriteFile(p, []byte(body), 0o644)
	}
	run("init", "-q", "-b", "main")
	run("config", "user.email", "t@t")
	run("config", "user.name", "t")
	w("cli/main.go", "package main\n")
	w("web/app.js", "// web\n")
	w("tpl/template.yaml", "version: 1.0.0\nlayout: {}\n")
	w("tpl/CHANGELOG.md", "# Changelog\n\n## 1.0.0 - 2026-01-01\n\n- start\n")
	w("docs/a.md", "# a\n")
	run("add", "-A")
	run("commit", "-q", "-m", "chore: [S-001] init")
	run("tag", "-a", "cli/v0.2.0", "-m", "cli 0.2.0")
	run("tag", "-a", "cli/v0.10.0", "-m", "cli 0.10.0") // numeric sort matters
	w("cli/feature.go", "package main\n")
	w("web/x.js", "// touched incidentally\n")
	run("add", "-A")
	run("commit", "-q", "-m", "feat: [S-002] feature")
	w("docs/b.md", "# b\n")
	run("add", "-A")
	run("commit", "-q", "-m", "docs: [S-003] docs only")
	w("tpl/root/x.md", "x\n")
	run("add", "-A")
	run("commit", "-q", "-m", "feat: [S-004] template change")
	return root, r
}

var m = manifest.Manifest{Projects: []manifest.Project{
	{Name: "cli", Path: "cli", Kind: "go", Tags: []string{"command"}},
	{Name: "web", Path: "web", Kind: "sveltekit"},
	{Name: "tpl", Path: "tpl", Kind: "template"},
}}

func TestComputeAndApply(t *testing.T) {
	root, r := gitRepo(t)
	story := &workitem.Item{ID: "S-002", Type: workitem.Story, Nature: "feature", Title: "Feature", Tags: []string{"command"}}
	plan, err := Compute(r, root, m, story, nil, "")
	if err != nil {
		t.Fatal(err)
	}
	if len(plan.Steps) != 2 || plan.Steps[0].Component.Name != "cli" || !plan.Steps[0].Delivered || plan.Steps[0].Level != Minor || plan.Steps[0].To.String() != "0.11.0" || plan.Steps[0].Tag != "cli/v0.11.0" {
		t.Errorf("delivered step: %+v", plan.Steps)
	}
	if plan.Steps[1].Component.Name != "web" || plan.Steps[1].Delivered || plan.Steps[1].Level != Patch || plan.Steps[1].To.String() != "0.0.1" {
		t.Errorf("incidental step: %+v", plan.Steps[1])
	}
	// docs-only story releases nothing
	docs := &workitem.Item{ID: "S-003", Type: workitem.Story, Nature: "feature", Title: "Docs"}
	plan, err = Compute(r, root, m, docs, nil, "")
	if err != nil || plan.Skipped == "" || len(plan.Steps) != 0 || len(plan.Untouched) != 1 {
		t.Errorf("docs-only: %+v %v", plan, err)
	}
	// ambiguous: two components touched, no tag
	amb := &workitem.Item{ID: "S-002", Type: workitem.Story, Nature: "feature", Title: "Ambiguous"}
	if _, err := Compute(r, root, m, amb, nil, ""); err == nil || !strings.Contains(err.Error(), "--deliver") {
		t.Errorf("ambiguous should ask for --deliver: %v", err)
	}
	// parent tags decide; --deliver overrides
	parent := &workitem.Item{ID: "E-001", Type: workitem.Epic, Tags: []string{"web"}}
	plan, _ = Compute(r, root, m, amb, parent, "")
	if plan.Steps[0].Component.Name != "web" || plan.Steps[0].Level != Minor {
		t.Errorf("parent tag: %+v", plan.Steps)
	}
	plan, _ = Compute(r, root, m, amb, parent, "cli")
	if plan.Steps[0].Component.Name != "cli" {
		t.Errorf("--deliver: %+v", plan.Steps)
	}
	// template component: version file and changelog, single touched => delivered
	tplStory := &workitem.Item{ID: "S-004", Type: workitem.Story, Nature: "improvement", Title: "Template tweak"}
	plan, err = Compute(r, root, m, tplStory, nil, "")
	if err != nil || len(plan.Steps) != 1 || plan.Steps[0].Version != "tpl/template.yaml" || plan.Steps[0].To.String() != "1.0.1" || plan.Steps[0].Tag != "" {
		t.Fatalf("template plan: %+v %v", plan, err)
	}
	if err := Apply(r, root, plan, time.Date(2026, 9, 17, 0, 0, 0, 0, time.UTC)); err != nil {
		t.Fatal(err)
	}
	ty, _ := os.ReadFile(filepath.Join(root, "tpl", "template.yaml"))
	cl, _ := os.ReadFile(filepath.Join(root, "tpl", "CHANGELOG.md"))
	if !strings.Contains(string(ty), "version: 1.0.1") || !strings.HasPrefix(string(cl), "# Changelog\n\n## 1.0.1 - 2026-09-17\n\n- S-004 Template tweak (patch).\n\n## 1.0.0") {
		t.Errorf("apply:\n%s\n%s", ty, cl)
	}
	// epic: major on delivered
	epic := &workitem.Item{ID: "S-002", Type: workitem.Epic, Nature: "feature", Title: "Epic", Tags: []string{"cli"}}
	plan, _ = Compute(r, root, m, epic, nil, "")
	if plan.Steps[0].To.String() != "1.0.0" {
		t.Errorf("epic major: %+v", plan.Steps[0])
	}
	tags, err := Tag(r, root, plan)
	if err != nil || len(tags) != 2 || tags[0] != "cli/v1.0.0" {
		t.Errorf("tag: %v %v", tags, err)
	}
	out, _ := r.Run(root, "git", "tag", "-l", "cli/*")
	if !strings.Contains(out, "cli/v1.0.0") {
		t.Errorf("tags: %s", out)
	}
}
