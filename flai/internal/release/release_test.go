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
	// ADR-0025: research is accepted without a release; experiment is refused.
	if got, err := LevelFor(&workitem.Item{ID: "S-9", Type: workitem.Story, Nature: "research"}); err != nil || got != None {
		t.Errorf("research is accepted with no release: %q %v", got, err)
	}
	_, err := LevelFor(&workitem.Item{ID: "S-9", Type: workitem.Story, Nature: "experiment"})
	if err == nil || !strings.Contains(err.Error(), "S-9 is an experiment") || strings.Contains(err.Error(), "research") {
		t.Errorf("experiment stays on its branch, and the message names only experiment: %v", err)
	}
}

// ADR-0025: a research story plans no release whatever it touched, and the
// plan names the components whose files land on main unreleased.
func TestComputeResearch(t *testing.T) {
	root, r := gitRepo(t)
	// S-003 touched only docs
	findings := &workitem.Item{ID: "S-003", Type: workitem.Story, Nature: "research", Title: "Findings"}
	plan, err := Compute(r, root, m, findings, nil, "")
	if err != nil {
		t.Fatal(err)
	}
	if plan.Level != None || len(plan.Steps) != 0 || len(plan.Unreleased) != 0 || !strings.Contains(plan.Skipped, "S-003 is research") || !strings.Contains(plan.Skipped, "ADR-0025") {
		t.Errorf("docs-only research: %+v", plan)
	}
	if len(plan.Commits) != 1 || len(plan.Untouched) != 1 {
		t.Errorf("the plan still says what the story changed: %+v", plan)
	}
	// S-002 touched two components and carries no tag: for a feature that is
	// an error asking for --deliver; for research there is nothing to deliver.
	code := &workitem.Item{ID: "S-002", Type: workitem.Story, Nature: "research", Title: "Spike with code"}
	plan, err = Compute(r, root, m, code, nil, "")
	if err != nil {
		t.Fatalf("research that touched components is not refused: %v", err)
	}
	if len(plan.Steps) != 0 || plan.Skipped == "" {
		t.Errorf("no step, so no tag and no version bump: %+v", plan)
	}
	if len(plan.Unreleased) != 2 || plan.Unreleased[0].Component != "cli" || plan.Unreleased[1].Component != "web" || len(plan.Unreleased[0].Files) != 1 || plan.Unreleased[0].Files[0] != "cli/feature.go" {
		t.Errorf("components landing unreleased, in manifest order: %+v", plan.Unreleased)
	}
	// a tag or --deliver changes nothing: research never delivers to a component
	code.Tags = []string{"command"}
	plan, _ = Compute(r, root, m, code, nil, "web")
	if len(plan.Steps) != 0 || len(plan.Unreleased) != 2 {
		t.Errorf("tags and --deliver do not make research release: %+v", plan)
	}
	// applying and tagging a research plan changes nothing
	if err := Apply(r, root, plan, time.Date(2026, 9, 19, 0, 0, 0, 0, time.UTC)); err != nil {
		t.Fatal(err)
	}
	tags, err := Tag(r, root, plan)
	if err != nil || len(tags) != 0 {
		t.Errorf("no tags for research: %v %v", tags, err)
	}
	out, _ := r.Run(root, "git", "status", "--porcelain")
	if strings.TrimSpace(out) != "" {
		t.Errorf("no version file or changelog changed: %q", out)
	}
	// experiment is refused by Compute too
	if _, err := Compute(r, root, m, &workitem.Item{ID: "S-003", Type: workitem.Story, Nature: "experiment"}, nil, ""); err == nil {
		t.Error("experiment must be refused")
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

// I-0016, S-0047: a tag delivers only to a component the commits touched.
func TestDeliveryGoesToATouchedComponent(t *testing.T) {
	root, r := gitRepo(t)
	run := func(args ...string) {
		if out, err := r.Run(root, "git", args...); err != nil {
			t.Fatalf("git %v: %v %s", args, err, out)
		}
	}
	w := func(rel string) {
		p := filepath.Join(root, rel)
		_ = os.MkdirAll(filepath.Dir(p), 0o755)
		_ = os.WriteFile(p, []byte("x\n"), 0o644)
	}
	// S-010: pure CLI work
	w("cli/only.go")
	run("add", "-A")
	run("commit", "-q", "-m", "feat: [S-010] cli only")
	// S-011: mostly web, a little cli
	w("web/a.js")
	w("web/b.js")
	w("web/c.js")
	w("cli/little.go")
	run("add", "-A")
	run("commit", "-q", "-m", "feat: [S-011] mostly web")
	// S-012: one file in each
	w("cli/tie.go")
	w("web/tie.js")
	run("add", "-A")
	run("commit", "-q", "-m", "feat: [S-012] a tie")

	steps := func(plan *Plan) (out []string) {
		for _, st := range plan.Steps {
			kind := "incidental"
			if st.Delivered {
				kind = "delivered"
			}
			out = append(out, st.Component.Name+" "+st.Level+" "+kind)
		}
		return
	}
	story := func(id string, tags ...string) *workitem.Item {
		return &workitem.Item{ID: id, Type: workitem.Story, Nature: "feature", Title: id, Tags: tags}
	}

	// the case that went wrong: [web, command] on CLI-only work
	plan, err := Compute(r, root, m, story("S-010", "web", "command"), nil, "")
	if err != nil {
		t.Fatal(err)
	}
	if got := steps(plan); len(got) != 1 || got[0] != "cli minor delivered" {
		t.Errorf("an untouched component gets no release at all, the touched one delivers: %v", got)
	}
	// several tagged and touched: the most touched files wins, whatever the tag order
	for _, tags := range [][]string{{"command", "web"}, {"web", "command"}} {
		plan, _ = Compute(r, root, m, story("S-011", tags...), nil, "")
		if got := steps(plan); len(got) != 2 || got[0] != "web minor delivered" || got[1] != "cli patch incidental" {
			t.Errorf("tags %v: %v", tags, got)
		}
	}
	// a tie goes to the earlier tag
	plan, _ = Compute(r, root, m, story("S-012", "web", "cli"), nil, "")
	if got := steps(plan); got[0] != "web minor delivered" {
		t.Errorf("tie, web first: %v", got)
	}
	plan, _ = Compute(r, root, m, story("S-012", "cli", "web"), nil, "")
	if got := steps(plan); got[0] != "cli minor delivered" {
		t.Errorf("tie, cli first: %v", got)
	}
	// the story's tags name only an untouched component: the epic's tags are asked next
	epic := &workitem.Item{ID: "E-001", Type: workitem.Epic, Tags: []string{"command"}}
	plan, err = Compute(r, root, m, story("S-010", "tpl"), epic, "")
	if err != nil || steps(plan)[0] != "cli minor delivered" || len(plan.Steps) != 1 {
		t.Errorf("epic tags: %v %v", steps(plan), err)
	}
	// no tag names a touched component, one component touched: it delivers
	plan, err = Compute(r, root, m, story("S-010", "tpl"), nil, "")
	if err != nil || len(plan.Steps) != 1 || steps(plan)[0] != "cli minor delivered" {
		t.Errorf("the only touched component: %v %v", steps(plan), err)
	}
	// no tag names a touched component, two touched: still the operator's call
	if _, err := Compute(r, root, m, story("S-011", "tpl"), nil, ""); err == nil || !strings.Contains(err.Error(), "--deliver") {
		t.Errorf("ambiguous still asks: %v", err)
	}
	// --deliver overrides everything, even towards an untouched component
	plan, _ = Compute(r, root, m, story("S-010", "command"), nil, "web")
	if got := steps(plan); got[0] != "web minor delivered" {
		t.Errorf("--deliver: %v", got)
	}
	// tagged, and touching no component at all: nothing to release
	plan, err = Compute(r, root, m, story("S-003", "command"), nil, "")
	if err != nil || plan.Skipped == "" || len(plan.Steps) != 0 {
		t.Errorf("a tagged docs-only story releases nothing: %+v %v", plan, err)
	}
}
