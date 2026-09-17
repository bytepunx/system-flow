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
