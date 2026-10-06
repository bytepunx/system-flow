package cmd

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/execx"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
)

func legacyRepo(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	w := func(rel, body string) {
		p := filepath.Join(root, rel)
		_ = os.MkdirAll(filepath.Dir(p), 0o755)
		_ = os.WriteFile(p, []byte(body), 0o644)
	}
	w("README.md", "# legacy\nkeep me\n")
	w("notes.md", "# notes\n")
	w("docs/guide.md", "# guide\n")
	w("adr/0001-first.md", "# first\n")
	w("services/api/go.mod", "module example.com/api\n")
	w("web/package.json", "{}")
	w("web/svelte.config.js", "export default {}")
	if _, err := exec.LookPath("git"); err == nil {
		r := execx.System{}
		for _, args := range [][]string{{"init", "-q", "-b", "main"}, {"config", "user.email", "t@t"}, {"config", "user.name", "t"}, {"add", "-A"}, {"commit", "-q", "-m", "legacy"}} {
			_, _ = r.Run(root, "git", args...)
		}
	}
	return root
}

func TestImportDryRunAndApply(t *testing.T) {
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	root := legacyRepo(t)
	out, errOut, code := runIn(t, ".", "import", root, "--template", "../../template", "--dry-run")
	if _, err := os.Stat(filepath.Join("..", "..", "template", "template.yaml")); err != nil {
		t.Skip("prototype template not present")
	}
	if code != 0 {
		t.Fatalf("dry run: %s", errOut)
	}
	for _, want := range []string{"move folder: adr/ -> design/adrs/", "sub-project: api (go) at services/api", "sub-project: web (sveltekit) at web", "markdown to place (1): notes.md", "docs (exists)", "dry run: nothing changed"} {
		if !strings.Contains(out, want) {
			t.Errorf("missing %q in:\n%s", want, out)
		}
	}
	if _, err := os.Stat(filepath.Join(root, "system-flow.yaml")); err == nil {
		t.Fatal("dry run wrote the manifest")
	}

	out, errOut, code = runIn(t, ".", "import", root, "--template", "../../template", "--yes", "--var", "description=Legacy")
	if code != 0 {
		t.Fatalf("apply: %s\n%s", errOut, out)
	}
	if b, _ := os.ReadFile(filepath.Join(root, "README.md")); string(b) != "# legacy\nkeep me\n" {
		t.Error("existing README was overwritten")
	}
	if _, err := os.Stat(filepath.Join(root, "design", "adrs", "0001-first.md")); err != nil {
		t.Error("adr folder not moved")
	}
	if _, err := os.Stat(filepath.Join(root, "notes.md")); err != nil {
		t.Error("loose markdown must be left in place without prompts")
	}
	if _, err := os.Stat(filepath.Join(root, "docs", "guide.md")); err != nil {
		t.Error("existing docs content lost")
	}
	if _, err := os.Stat(filepath.Join(root, "design", "conventions", "README.md")); err != nil {
		t.Error("template files not rendered")
	}
	m, err := manifest.Load(filepath.Join(root, "system-flow.yaml"))
	if err != nil || len(m.Projects) != 2 || m.Projects[0].Kind != "go" || m.Projects[1].Kind != "sveltekit" || m.Description != "Legacy" {
		t.Errorf("manifest: %+v %v", m, err)
	}
	if !strings.Contains(out, "moved adr/0001-first.md -> design/adrs/0001-first.md") || !strings.Contains(out, "2 sub-projects recorded") {
		t.Errorf("summary:\n%s", out)
	}
	if _, err := exec.LookPath("git"); err == nil {
		st, _ := (execx.System{}).Run(root, "git", "status", "--short")
		if !strings.Contains(st, "R  adr/0001-first.md -> design/adrs/0001-first.md") {
			t.Errorf("expected git mv:\n%s", st)
		}
	}
	_, errOut, code = runIn(t, ".", "import", root, "--template", "../../template", "--yes")
	if code == 0 || !strings.Contains(errOut, "already has system-flow.yaml") {
		t.Errorf("second import must refuse without --force: %d %s", code, errOut)
	}
	out, _, code = runIn(t, ".", "import", root, "--template", "../../template", "--dry-run", "--json")
	var v struct {
		Plan struct {
			Projects []any `json:"projects"`
		} `json:"plan"`
	}
	if code == 0 && json.Unmarshal([]byte(out), &v) == nil && len(v.Plan.Projects) != 2 {
		t.Errorf("json dry run: %s", out)
	}
}

// An empty required --var stops flai import before it moves a folder or
// writes a file (I-0041).
func TestImportRefusesAnEmptyRequiredVarBeforeMovingAnything(t *testing.T) {
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	root := legacyRepo(t)
	_, errOut, code := runIn(t, ".", "import", root, "--template", miniTemplate, "--yes", "--var", "project_name=")
	if code == 0 || !strings.Contains(errOut, "required variables not set: project_name") {
		t.Fatalf("expected the required error, got %d %s", code, errOut)
	}
	if _, err := os.Stat(filepath.Join(root, "adr", "0001-first.md")); err != nil {
		t.Errorf("a refused import moved adr/: %v", err)
	}
	if _, err := os.Stat(filepath.Join(root, "system-flow.yaml")); err == nil {
		t.Error("a refused import wrote the manifest")
	}
}

// With no --ref, an import applies the template's newest release when the
// configured ref follows releases; --ref, as a tag or as the version it
// spells, wins; a configured ref naming another branch is used as given
// (ADR-0103). An import writes no lock, so system-flow.yaml is the record.
func TestImportAppliesTheNewestRelease(t *testing.T) {
	url := releasedTemplate(t)
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	if _, errOut, code := runIn(t, ".", "template", "use", url, "--ref", "main"); code != 0 {
		t.Fatalf("template use: %s", errOut)
	}
	importedAt := func(root, ref, version string) {
		t.Helper()
		if b, _ := os.ReadFile(filepath.Join(root, "team.txt")); string(b) != "v"+version+"\n" {
			t.Errorf("team.txt = %q, want the render of %s", b, version)
		}
		m, err := manifest.Load(filepath.Join(root, manifest.File))
		if err != nil || m.Template.Repo != url || m.Template.Ref != ref || m.Template.Version != version {
			t.Errorf("system-flow.yaml must record %s at %s from %s: %+v %v", ref, version, url, m.Template, err)
		}
	}

	root := legacyRepo(t)
	out, errOut, code := runIn(t, ".", "import", root, "--yes")
	if code != 0 {
		t.Fatalf("import: %s\n%s", errOut, out)
	}
	if !strings.Contains(out, "template version: v1.0.60, the newest release") {
		t.Errorf("import does not say it applied the newest release:\n%s", out)
	}
	importedAt(root, "v1.0.60", "1.0.60")

	out, errOut, code = runIn(t, ".", "import", legacyRepo(t), "--template", url, "--dry-run", "--json")
	var v struct {
		Template struct {
			Ref     string `json:"ref"`
			Version string `json:"version"`
			Newest  bool   `json:"newest_release"`
		} `json:"template"`
	}
	if code != 0 || json.Unmarshal([]byte(out), &v) != nil || v.Template.Ref != "v1.0.60" || v.Template.Version != "1.0.60" || !v.Template.Newest {
		t.Errorf("import --dry-run --json: %d %+v %s %s", code, v, out, errOut)
	}

	for _, ref := range []string{"v1.0.18", "1.0.18"} {
		root := legacyRepo(t)
		out, errOut, code := runIn(t, ".", "import", root, "--ref", ref, "--yes")
		if code != 0 {
			t.Fatalf("import --ref %s: %s", ref, errOut)
		}
		if strings.Contains(out, "newest release") {
			t.Errorf("import --ref %s calls it the newest release:\n%s", ref, out)
		}
		importedAt(root, "v1.0.18", "1.0.18")
	}

	if _, errOut, code := runIn(t, ".", "template", "use", url, "--ref", "edge"); code != 0 {
		t.Fatalf("template use --ref edge: %s", errOut)
	}
	root = legacyRepo(t)
	if _, errOut, code := runIn(t, ".", "import", root, "--yes"); code != 0 {
		t.Fatalf("import on edge: %s", errOut)
	}
	importedAt(root, "edge", "2.0.0")
}
