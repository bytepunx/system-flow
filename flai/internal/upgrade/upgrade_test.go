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
