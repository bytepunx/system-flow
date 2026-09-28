package upgrade

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/lock"
	"github.com/bytepunx/system-flow/flai/internal/template"
)

const mini = "../template/testdata/mini"

func vars() map[string]any {
	return map[string]any{"project_name": "p", "project_key": "p", "description": "d"}
}

// render creates a project from the fixture and returns its root and lock.
func render(t *testing.T) (string, *lock.Lock) {
	t.Helper()
	m, err := template.LoadManifest(mini)
	if err != nil {
		t.Fatal(err)
	}
	root := t.TempDir()
	src := template.Source{Repo: mini, Local: true, Dir: mini}
	res, err := template.Render(m, mini, root, template.Options{Vars: vars(), Source: src})
	if err != nil {
		t.Fatal(err)
	}
	lk := &lock.Lock{Template: lock.Template{Version: m.Version}, Files: res.Hashes}
	if err := lock.Save(root, lk); err != nil {
		t.Fatal(err)
	}
	back, err := lock.Load(root)
	if err != nil || back == nil || len(back.Files) != len(res.Hashes) || back.Files["README.md"] != res.Hashes["README.md"] {
		t.Fatalf("lock round trip: %+v %v", back, err)
	}
	return root, back
}

// newer copies the fixture and changes it: bumped version, changed README
// template, changed CLAUDE baseline, changed plain.txt, a new file.
func newer(t *testing.T) string {
	t.Helper()
	dst := t.TempDir()
	_ = filepath.WalkDir(mini, func(path string, d os.DirEntry, err error) error {
		rel, _ := filepath.Rel(mini, path)
		if d.IsDir() {
			return os.MkdirAll(filepath.Join(dst, rel), 0o755)
		}
		b, _ := os.ReadFile(path)
		return os.WriteFile(filepath.Join(dst, rel), b, 0o644)
	})
	w := func(rel, body string) { _ = os.WriteFile(filepath.Join(dst, rel), []byte(body), 0o644) }
	y, _ := os.ReadFile(filepath.Join(dst, "template.yaml"))
	w("template.yaml", strings.Replace(string(y), "version: 9.9.9", "version: 10.0.0", 1))
	w("root/README.md.tmpl", "# {{ .project_name }} v2\n")
	w("root/CLAUDE.md.tmpl", "# {{ .project_name }}\n\nBaseline line two.\n\n<!-- system-flow:end-of-baseline -->\n\n## Project additions\n")
	w("root/plain.txt", "verbatim v2\n")
	w("root/new.txt", "brand new\n")
	return dst
}

func TestComputeClassesAndApply(t *testing.T) {
	root, lk := render(t)
	// project edits: CLAUDE additions below the marker, plain.txt changed
	claude := filepath.Join(root, "CLAUDE.md")
	b, _ := os.ReadFile(claude)
	_ = os.WriteFile(claude, append(b, []byte("- my rule\n")...), 0o644)
	_ = os.WriteFile(filepath.Join(root, "plain.txt"), []byte("edited by project\n"), 0o644)

	nt := newer(t)
	m, _ := template.LoadManifest(nt)
	src := template.Source{Repo: nt, Local: true, Dir: nt}
	plan, err := Compute(root, m, src, template.Options{Vars: vars(), Source: src}, lk, "9.9.9")
	if err != nil {
		t.Fatal(err)
	}
	classes := map[string]string{}
	for _, c := range plan.Changes {
		classes[c.Path] = c.Class
	}
	want := map[string]string{"README.md": Replace, "CLAUDE.md": Merge, "plain.txt": Conflict, "new.txt": Add, "wip/README.md": Same, "design/system/overview.md": Same}
	for p, cls := range want {
		if classes[p] != cls {
			t.Errorf("%s: got %s want %s", p, classes[p], cls)
		}
	}
	if plan.NoLock || len(plan.Conflicts()) != 1 {
		t.Errorf("plan: %+v", plan)
	}
	written, err := Apply(root, plan, Policy{"plain.txt": "keep"})
	if err != nil {
		t.Fatal(err)
	}
	if len(written) != 3 {
		t.Errorf("written: %v", written)
	}
	got, _ := os.ReadFile(claude)
	if !strings.Contains(string(got), "Baseline line two.") || !strings.Contains(string(got), "- my rule") || strings.Contains(string(got), "line one") {
		t.Errorf("merge:\n%s", got)
	}
	if b, _ := os.ReadFile(filepath.Join(root, "plain.txt")); string(b) != "edited by project\n" {
		t.Error("kept conflict was overwritten")
	}
	if b, _ := os.ReadFile(filepath.Join(root, "README.md")); string(b) != "# p v2\n" {
		t.Errorf("replace: %q", b)
	}
	if _, err := os.Stat(filepath.Join(root, "new.txt")); err != nil {
		t.Error("add missing")
	}
	nl := NewLock(plan, src, m.Version, time.Now())
	if nl.Files["plain.txt"] != lock.Hash([]byte("verbatim v2\n")) {
		t.Error("lock must record the template's hash for kept conflicts")
	}
	// replace-all policy
	written, _ = Apply(root, plan, Policy{"plain.txt": "replace"})
	if b, _ := os.ReadFile(filepath.Join(root, "plain.txt")); string(b) != "verbatim v2\n" || len(written) != 4 {
		t.Errorf("replace policy: %q %v", b, written)
	}
}

func TestNoLockMakesEveryDifferenceAConflict(t *testing.T) {
	root, _ := render(t)
	_ = os.Remove(filepath.Join(root, lock.File))
	nt := newer(t)
	m, _ := template.LoadManifest(nt)
	src := template.Source{Repo: nt, Local: true, Dir: nt}
	plan, err := Compute(root, m, src, template.Options{Vars: vars(), Source: src}, nil, "9.9.9")
	if err != nil {
		t.Fatal(err)
	}
	if !plan.NoLock || plan.Count(Replace) != 0 || plan.Count(Conflict) != 2 || plan.Count(Merge) != 1 || plan.Count(Add) != 1 {
		t.Errorf("%s", Describe(plan))
	}
	rl, err := Relock(root, m, src, template.Options{Vars: vars(), Source: src}, time.Now())
	if err != nil || len(rl.Files) < 6 || rl.Files["README.md"] != mustHash(t, filepath.Join(root, "README.md")) {
		t.Errorf("relock: %+v %v", rl, err)
	}
}

func mustHash(t *testing.T, p string) string {
	h, err := lock.HashFile(p)
	if err != nil {
		t.Fatal(err)
	}
	return h
}

// S-0134: a project made before conventions carried topics gets topics:
// [all] on every convention from flai upgrade, and keeps its additions.
func TestUpgradeBringsConventionTopics(t *testing.T) {
	const tpl = "../../../template"
	m, err := template.LoadManifest(tpl)
	if err != nil {
		t.Skip("the monorepo's template is not present")
	}
	src := template.Source{Repo: tpl, Local: true, Dir: tpl}
	v := vars()
	for _, k := range []string{"owner", "repo_url", "today"} {
		v[k] = ""
	}
	opts := template.Options{Vars: v, Source: src}
	root := t.TempDir()
	res, err := template.Render(m, tpl, root, opts)
	if err != nil {
		t.Fatal(err)
	}
	dir := filepath.Join(root, "design", "conventions")
	files, _ := filepath.Glob(filepath.Join(dir, "*.md"))
	if len(files) == 0 {
		t.Fatal("the template renders no conventions")
	}
	// as the project was: no topics, and a rule of its own below the marker
	for _, f := range files {
		b, _ := os.ReadFile(f)
		if !strings.Contains(string(b), "\ntopics: [all]\n") {
			t.Fatalf("the template's %s has no topics: [all]", filepath.Base(f))
		}
		old := strings.Replace(string(b), "topics: [all]\n", "", 1)
		if err := os.WriteFile(f, []byte(old), 0o644); err != nil {
			t.Fatal(err)
		}
		res.Hashes[filepath.ToSlash(filepath.Join("design", "conventions", filepath.Base(f)))] = lock.Hash([]byte(old))
		if filepath.Base(f) != "README.md" {
			_ = os.WriteFile(f, []byte(old+"- A rule of this project.\n"), 0o644)
		}
	}
	lk := &lock.Lock{Template: lock.Template{Version: "1.0.0"}, Files: res.Hashes}
	plan, err := Compute(root, m, src, opts, lk, "1.0.0")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := Apply(root, plan, Policy{}); err != nil {
		t.Fatal(err)
	}
	for _, f := range files {
		b, _ := os.ReadFile(f)
		if !strings.Contains(string(b), "\ntopics: [all]\n") {
			t.Errorf("%s has no topics after the upgrade:\n%s", filepath.Base(f), b)
		}
		if filepath.Base(f) != "README.md" && !strings.HasSuffix(string(b), "- A rule of this project.\n") {
			t.Errorf("%s lost the project's addition:\n%s", filepath.Base(f), b)
		}
	}
}

// S-0134: the merge keeps the topics a project set on a convention, and
// takes the template's when the project has none or has what the template
// last gave it.
func TestMergeKeepsTheProjectsTopics(t *testing.T) {
	const marker = "<!-- system-flow:end-of-baseline -->"
	file := func(topics, body string) []byte {
		fm := "---\ntitle: X\n"
		if topics != "" {
			fm += "topics: [" + topics + "]\n"
		}
		return []byte(fm + "order: 10\n---\n\n# X\n\n" + body + "\n\n" + marker + "\n\n## Project additions\n- mine\n")
	}
	tpl := file("code", "v2")
	cases := []struct {
		name, project string
		was           []string
		want          string
	}{
		{"unchanged since applied", "all", []string{"all"}, "code"},
		{"narrowed by the project", "cli, go", []string{"all"}, "cli, go"},
		{"nothing recorded", "cli", nil, "cli"},
		{"no topics in the project", "", []string{"all"}, "code"},
		{"the template's already", "code", []string{"all"}, "code"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := string(merge(tpl, file(c.project, "v1"), c.was))
			want := string(file(c.want, "v2"))
			if got != want {
				t.Errorf("got:\n%s\nwant:\n%s", got, want)
			}
		})
	}
	if Topics([]byte("---\ntopics: [code]\n---\n\n# no marker\n")) != nil {
		t.Error("a file without the marker has no topics to record")
	}
}

func TestLockRoundTripsTopics(t *testing.T) {
	root := t.TempDir()
	lk := &lock.Lock{Files: map[string]string{"a.md": "h"}, Topics: map[string][]string{"a.md": {"cli", "go"}, "b.md": {"all"}}}
	if err := lock.Save(root, lk); err != nil {
		t.Fatal(err)
	}
	back, err := lock.Load(root)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Join(back.Topics["a.md"], ",") != "cli,go" || strings.Join(back.Topics["b.md"], ",") != "all" {
		t.Errorf("topics: %v", back.Topics)
	}
}

// S-0134: a project that narrows a convention's topics keeps them through
// flai upgrade, and a convention it left alone takes the template's.
func TestUpgradeKeepsNarrowedConventionTopics(t *testing.T) {
	const tpl = "../../../template"
	m, err := template.LoadManifest(tpl)
	if err != nil {
		t.Skip("the monorepo's template is not present")
	}
	src := template.Source{Repo: tpl, Local: true, Dir: tpl}
	v := vars()
	for _, k := range []string{"owner", "repo_url", "today"} {
		v[k] = ""
	}
	opts := template.Options{Vars: v, Source: src}
	root := t.TempDir()
	res, err := template.Render(m, tpl, root, opts)
	if err != nil {
		t.Fatal(err)
	}
	lk := &lock.Lock{Template: lock.Template{Version: m.Version}, Files: res.Hashes, Topics: TopicsOf(root, res.Written)}
	narrowed := filepath.Join(root, "design", "conventions", "code-quality.md")
	left := filepath.Join(root, "design", "conventions", "git.md")
	if strings.Join(lk.Topics["design/conventions/code-quality.md"], ",") != "all" {
		t.Fatalf("the lock does not record the template's topics: %v", lk.Topics)
	}
	edit := func(p, from, to string) {
		b, _ := os.ReadFile(p)
		if !strings.Contains(string(b), from) {
			t.Fatalf("%s has no %q", p, from)
		}
		_ = os.WriteFile(p, []byte(strings.Replace(string(b), from, to, 1)), 0o644)
	}
	edit(narrowed, "topics: [all]", "topics: [code]")
	// as if the template had given git.md other topics before
	edit(left, "topics: [all]", "topics: [old]")
	lk.Topics["design/conventions/git.md"] = []string{"old"}
	plan, err := Compute(root, m, src, opts, lk, m.Version)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := Apply(root, plan, Policy{}); err != nil {
		t.Fatal(err)
	}
	if b, _ := os.ReadFile(narrowed); !strings.Contains(string(b), "\ntopics: [code]\n") {
		t.Errorf("the project's topics were not kept:\n%s", b)
	}
	if b, _ := os.ReadFile(left); !strings.Contains(string(b), "\ntopics: [all]\n") {
		t.Errorf("the template's topics were not taken:\n%s", b)
	}
	nl := NewLock(plan, src, m.Version, time.Now())
	if strings.Join(nl.Topics["design/conventions/code-quality.md"], ",") != "all" {
		t.Errorf("the new lock must record the template's topics, not the project's: %v", nl.Topics["design/conventions/code-quality.md"])
	}
	rl, err := Relock(root, m, src, opts, time.Now())
	if err != nil || strings.Join(rl.Topics["design/conventions/code-quality.md"], ",") != "all" {
		t.Errorf("relock must record the template's topics: %v %v", rl.Topics, err)
	}
}
