package template

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/execx"
)

const mini = "testdata/mini"

func TestLoadManifest(t *testing.T) {
	m, err := LoadManifest(mini)
	if err != nil {
		t.Fatal(err)
	}
	if m.Version != "9.9.9" || len(m.Variables) != 3 || m.Layout["design"] != "design" || m.Render.Suffix != ".tmpl" {
		t.Fatalf("unexpected manifest: %+v", m)
	}
	if _, err := LoadManifest(t.TempDir()); err == nil {
		t.Fatal("expected error for missing manifest")
	}
	bad := t.TempDir()
	_ = os.WriteFile(filepath.Join(bad, ManifestFile), []byte("version: 1\nnope: 2\n"), 0o644)
	if _, err := LoadManifest(bad); err == nil {
		t.Fatal("expected strict parse error for unknown key")
	}
}

func TestHelpers(t *testing.T) {
	if got := Initials("system-flow"); got != "sf" {
		t.Errorf("Initials = %q", got)
	}
	if got := Initials("My Big  Project 2"); got != "mbp2" {
		t.Errorf("Initials = %q", got)
	}
	if got := Slug(" Hello, World! "); got != "hello-world" {
		t.Errorf("Slug = %q", got)
	}
}

func TestEvalDefault(t *testing.T) {
	m, _ := LoadManifest(mini)
	data := Data(m, Options{Vars: map[string]any{"project_name": "system-flow"}})
	got, err := EvalDefault(m.Variables[1], data)
	if err != nil || got != "sf" {
		t.Fatalf("EvalDefault = %q, %v", got, err)
	}
	if got, _ := EvalDefault(m.Variables[2], data); got != "" {
		t.Fatalf("plain default = %q", got)
	}
}

func render(t *testing.T, opt Options) (string, Result) {
	t.Helper()
	m, err := LoadManifest(mini)
	if err != nil {
		t.Fatal(err)
	}
	dest := t.TempDir()
	res, err := Render(m, mini, dest, opt)
	if err != nil {
		t.Fatal(err)
	}
	return dest, res
}

func TestRenderRules(t *testing.T) {
	now := time.Date(2026, 9, 15, 18, 30, 0, 0, time.UTC)
	dest, res := render(t, Options{
		Vars:   map[string]any{"project_name": "Demo App", "project_key": "da", "description": "d"},
		Layout: map[string]string{"design": "arch"},
		Source: Source{Repo: "https://x/tpl.git", Ref: "v1"},
		Now:    now,
	})
	read := func(p string) string {
		b, err := os.ReadFile(filepath.Join(dest, p))
		if err != nil {
			t.Fatalf("missing %s: %v", p, err)
		}
		return string(b)
	}
	if got := read("README.md"); !strings.Contains(got, "# Demo App (da)") || !strings.Contains(got, "design=arch wip=wip today=2026-09-15 tpl=9.9.9") {
		t.Errorf("rendered README:\n%s", got)
	}
	if got := read("plain.txt"); got != "verbatim {{ not rendered }}\n" {
		t.Errorf("verbatim copy changed: %q", got)
	}
	if got := read("system-flow.yaml"); !strings.Contains(got, `repo: "https://x/tpl.git"`) || !strings.Contains(got, "applied: 2026-09-15T18:30:00Z") {
		t.Errorf("manifest:\n%s", got)
	}
	if _, err := os.Stat(filepath.Join(dest, "arch", "system", "overview.md")); err != nil {
		t.Error("layout rename not applied to first path segment")
	}
	if _, err := os.Stat(filepath.Join(dest, "design")); err == nil {
		t.Error("default layout folder should not exist after rename")
	}
	for _, p := range []string{"IGNORED.md", "notes.skip", "README.md.tmpl"} {
		if _, err := os.Stat(filepath.Join(dest, p)); err == nil {
			t.Errorf("%s should not be written", p)
		}
	}
	if st, _ := os.Stat(filepath.Join(dest, "tools", "run.sh")); st == nil || st.Mode().Perm()&0o100 == 0 {
		t.Error("exec bit not preserved")
	}
	if len(res.Written) != 6 || len(res.Skipped) != 0 {
		t.Errorf("written=%v skipped=%v", res.Written, res.Skipped)
	}
}

func TestRenderNoOverwriteWithoutForce(t *testing.T) {
	m, _ := LoadManifest(mini)
	dest := t.TempDir()
	_ = os.WriteFile(filepath.Join(dest, "README.md"), []byte("mine"), 0o644)
	opt := Options{Vars: map[string]any{"project_name": "p", "project_key": "p", "description": ""}}
	res, err := Render(m, mini, dest, opt)
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Skipped) != 1 || res.Skipped[0] != "README.md" {
		t.Fatalf("skipped = %v", res.Skipped)
	}
	if b, _ := os.ReadFile(filepath.Join(dest, "README.md")); string(b) != "mine" {
		t.Fatal("existing file was overwritten")
	}
	opt.Force = true
	if _, err := Render(m, mini, dest, opt); err != nil {
		t.Fatal(err)
	}
	if b, _ := os.ReadFile(filepath.Join(dest, "README.md")); string(b) == "mine" {
		t.Fatal("force did not overwrite")
	}
}

func TestRenderMissingVariableFails(t *testing.T) {
	m, _ := LoadManifest(mini)
	_, err := Render(m, mini, t.TempDir(), Options{Vars: map[string]any{"project_name": "p"}})
	if err == nil || !strings.Contains(err.Error(), "project_key") {
		t.Fatalf("expected missing key error, got %v", err)
	}
}

func TestResolveLocalAndRemote(t *testing.T) {
	s, err := Resolve(mini, "ignored", "/cache")
	if err != nil || !s.Local || !filepath.IsAbs(s.Dir) || s.Ref != "" {
		t.Fatalf("local resolve: %+v %v", s, err)
	}
	if !s.Cached() {
		t.Fatal("local source is always cached")
	}
	if err := s.Ensure(execx.System{}, false); err != nil {
		t.Fatal(err)
	}
	r, err := Resolve("https://example.com/t.git", "main", "/cache")
	if err != nil || r.Local || !strings.HasPrefix(r.Dir, filepath.Join("/cache", "templates")) || r.Cached() {
		t.Fatalf("remote resolve: %+v %v", r, err)
	}
	r2, _ := Resolve("https://example.com/t.git", "other", "/cache")
	if r2.Dir == r.Dir {
		t.Fatal("different refs must cache separately")
	}
	if _, err := Resolve("", "", "/cache"); err == nil {
		t.Fatal("empty repo should error")
	}
}

func TestEnsureClonesWithGit(t *testing.T) {
	if testing.Short() {
		t.Skip("integration: needs real git")
	}
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}
	run := func(dir string, args ...string) {
		t.Helper()
		if out, err := (execx.System{}).Run(dir, "git", args...); err != nil {
			t.Fatalf("git %v: %v\n%s", args, err, out)
		}
	}
	// Build a template repo with a branch.
	src := t.TempDir()
	run(src, "init", "-q", "-b", "main")
	run(src, "config", "user.email", "t@t")
	run(src, "config", "user.name", "t")
	_ = os.WriteFile(filepath.Join(src, ManifestFile), []byte("version: 1.0.0\n"), 0o644)
	_ = os.MkdirAll(filepath.Join(src, RootDir), 0o755)
	_ = os.WriteFile(filepath.Join(src, RootDir, "a.txt"), []byte("main\n"), 0o644)
	run(src, "add", "-A")
	run(src, "commit", "-q", "-m", "init")
	run(src, "checkout", "-q", "-b", "feature")
	_ = os.WriteFile(filepath.Join(src, RootDir, "a.txt"), []byte("feature\n"), 0o644)
	run(src, "commit", "-q", "-am", "feature")
	run(src, "checkout", "-q", "main")

	cache := t.TempDir()
	s, _ := Resolve("file://"+src, "feature", cache)
	if s.Local {
		t.Fatal("expected git source")
	}
	if err := s.Ensure(execx.System{}, false); err != nil {
		t.Fatal(err)
	}
	if b, _ := os.ReadFile(filepath.Join(s.Dir, RootDir, "a.txt")); string(b) != "feature\n" {
		t.Fatalf("wrong ref checked out: %q", b)
	}
	if !s.Cached() {
		t.Fatal("should be cached after clone")
	}
	// A commit hash ref goes through the fallback path.
	hash, _ := (execx.System{}).Run(src, "git", "rev-parse", "main")
	h, _ := Resolve("file://"+src, hash, cache)
	if err := h.Ensure(execx.System{}, false); err != nil {
		t.Fatal(err)
	}
	if b, _ := os.ReadFile(filepath.Join(h.Dir, RootDir, "a.txt")); string(b) != "main\n" {
		t.Fatalf("hash ref: %q", b)
	}
	bad, _ := Resolve("file://"+src, "nope", cache)
	if err := bad.Ensure(execx.System{}, false); err == nil {
		t.Fatal("missing ref should fail")
	}
}
