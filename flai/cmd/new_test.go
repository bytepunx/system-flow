package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/goccy/go-yaml"

	"github.com/bytepunx/system-flow/flai/internal/lock"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/prompt"
	"github.com/bytepunx/system-flow/flai/internal/template"
)

const miniTemplate = "../internal/template/testdata/mini"

func TestNewWithDefaultsAndVars(t *testing.T) {
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	dest := filepath.Join(t.TempDir(), "demo-app")
	out, errOut, code := runCLI(t, "new", dest, "--template", miniTemplate, "--defaults", "--var", "description=Hello", "--no-git")
	if code != 0 {
		t.Fatalf("exit %d: %s%s", code, out, errOut)
	}
	if !strings.Contains(out, "7 files written") {
		t.Errorf("summary: %s", out)
	}
	readme, _ := os.ReadFile(filepath.Join(dest, "README.md"))
	if !strings.Contains(string(readme), "# demo-app (da)") || !strings.Contains(string(readme), "Hello") {
		t.Errorf("project_name should default to dir name and key to initials:\n%s", readme)
	}
	m, err := manifest.Load(filepath.Join(dest, manifest.File))
	if err != nil || m.Name != "demo-app" || m.Template.Version != "9.9.9" {
		t.Fatalf("manifest: %+v %v", m, err)
	}
	if _, err := os.Stat(filepath.Join(dest, ".git")); err == nil {
		t.Error("--no-git should skip git init")
	}

	// Second run keeps existing files, reports skips, and --json is parseable.
	out, _, code = runCLI(t, "new", dest, "--template", miniTemplate, "--defaults", "--no-git", "--json")
	if code != 0 {
		t.Fatalf("rerun exit %d", code)
	}
	var res struct {
		Written []string `json:"written"`
		Skipped []string `json:"skipped"`
	}
	if err := json.Unmarshal([]byte(out), &res); err != nil {
		t.Fatalf("json: %v\n%s", err, out)
	}
	if len(res.Written) != 0 || len(res.Skipped) != 7 {
		t.Errorf("rerun written=%v skipped=%v", res.Written, res.Skipped)
	}
}

func TestNewLayoutRenameAndErrors(t *testing.T) {
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	dest := filepath.Join(t.TempDir(), "p")
	_, errOut, code := runCLI(t, "new", dest, "--template", miniTemplate, "--defaults", "--no-git", "--layout", "design=arch")
	if code != 0 {
		t.Fatalf("exit %d: %s", code, errOut)
	}
	if _, err := os.Stat(filepath.Join(dest, "arch", "system", "overview.md")); err != nil {
		t.Error("layout rename not applied")
	}
	m, _ := manifest.Load(filepath.Join(dest, manifest.File))
	if m.Layout["design"] != "arch" {
		t.Errorf("manifest layout not rewritten: %v", m.Layout)
	}

	_, errOut, code = runCLI(t, "new", filepath.Join(t.TempDir(), "q"), "--template", miniTemplate, "--defaults", "--no-git", "--layout", "nope=x")
	if code == 0 || !strings.Contains(errOut, "unknown layout key") {
		t.Errorf("bad layout key: %d %s", code, errOut)
	}
	_, errOut, code = runCLI(t, "new", filepath.Join(t.TempDir(), "q"), "--template", miniTemplate, "--defaults", "--no-git", "--var", "bogus=1")
	if code == 0 || !strings.Contains(errOut, "unknown variable") {
		t.Errorf("bad var: %d %s", code, errOut)
	}
	_, errOut, code = runCLI(t, "new", filepath.Join(t.TempDir(), "q"), "--template", t.TempDir(), "--defaults", "--no-git")
	if code == 0 || !strings.Contains(errOut, "not a template") {
		t.Errorf("non-template dir: %d %s", code, errOut)
	}
}

func TestNewRequiresVarsWhenNotInteractive(t *testing.T) {
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	tpl := t.TempDir()
	_ = os.MkdirAll(filepath.Join(tpl, "root"), 0o755)
	_ = os.WriteFile(filepath.Join(tpl, "template.yaml"), []byte("version: 1\nvariables:\n  - name: team\n    required: true\nlayout: {design: design, docs: docs, wip: wip}\n"), 0o644)
	_ = os.WriteFile(filepath.Join(tpl, "root", "system-flow.yaml.tmpl"), []byte("version: 1\nname: x\nlayout: {design: design, docs: docs, wip: wip}\n"), 0o644)
	_, errOut, code := runCLI(t, "new", filepath.Join(t.TempDir(), "p"), "--template", tpl, "--no-git")
	if code == 0 || !strings.Contains(errOut, "required variables not set: team") {
		t.Fatalf("expected required error, got %d %s", code, errOut)
	}
	_, _, code = runCLI(t, "new", filepath.Join(t.TempDir(), "p"), "--template", tpl, "--no-git", "--var", "team=core")
	if code != 0 {
		t.Fatal("var should satisfy requirement")
	}
}

// An empty value given with --var for a required variable is refused before
// anything is written, as an empty default is (I-0041).
func TestNewRefusesAnEmptyRequiredVarBeforeWriting(t *testing.T) {
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	dest := filepath.Join(t.TempDir(), "p")
	_, errOut, code := runCLI(t, "new", dest, "--template", miniTemplate, "--defaults", "--no-git", "--var", "project_name=")
	if code == 0 || !strings.Contains(errOut, "required variables not set: project_name") {
		t.Fatalf("expected the required error, got %d %s", code, errOut)
	}
	if _, err := os.Stat(dest); !os.IsNotExist(err) {
		t.Errorf("a refused flai new left %s behind: %v", dest, err)
	}
}

func TestTemplateCommands(t *testing.T) {
	cfgPath := filepath.Join(t.TempDir(), "cfg.json")
	t.Setenv("FLAI_CONFIG", cfgPath)
	out, _, code := runCLI(t, "template", "use", miniTemplate, "--ref", "x")
	if code != 0 || !strings.Contains(out, "template.repo = "+miniTemplate) {
		t.Fatalf("use: %d %s", code, out)
	}
	out, _, code = runCLI(t, "template", "show")
	if code != 0 || !strings.Contains(out, "version:  9.9.9") || !strings.Contains(out, "project_name") {
		t.Fatalf("show: %d %s", code, out)
	}
	out, _, code = runCLI(t, "template", "update")
	if code != 0 || !strings.Contains(out, "local directory") {
		t.Fatalf("update: %d %s", code, out)
	}
	out, _, code = runCLI(t, "template", "show", "--json")
	if code != 0 {
		t.Fatal("show --json")
	}
	var v map[string]any
	if err := json.Unmarshal([]byte(out), &v); err != nil || v["local"] != true {
		t.Fatalf("show json: %v %v", err, v)
	}
}

// TestRenderPrototypeTemplate renders the real ./template from the monorepo
// and checks the result is a valid conforming project. Skipped when the
// prototype is not present (for example when flai is built standalone).
func TestRenderPrototypeTemplate(t *testing.T) {
	if testing.Short() {
		t.Skip("integration: renders the monorepo template")
	}
	proto := filepath.Join("..", "..", "template")
	if _, err := os.Stat(filepath.Join(proto, "template.yaml")); err != nil {
		t.Skip("prototype template not present")
	}
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	dest := filepath.Join(t.TempDir(), "sample")
	_, errOut, code := runCLI(t, "new", dest, "--template", proto, "--defaults", "--no-git",
		"--var", "description=Sample project", "--var", "owner=qa")
	if code != 0 {
		t.Fatalf("render prototype: %s", errOut)
	}
	m, err := manifest.Load(filepath.Join(dest, manifest.File))
	if err != nil {
		t.Fatal(err)
	}
	if m.Name != "sample" || m.Key != "s" || m.Template.Version == "" || m.Dashboard.Port != 4242 {
		t.Errorf("manifest: %+v", m)
	}
	for _, p := range []string{"CLAUDE.md", "README.md", ".editorconfig", ".gitignore", "Makefile",
		"design/README.md", "design/adrs/0000-template.md", "design/system/overview.md", "design/tech/README.md",
		"docs/users/index.md", "wip/kanban/board.md", "wip/agents/index.md", ".github/workflows/system-flow-check.yml"} {
		if _, err := os.Stat(filepath.Join(dest, p)); err != nil {
			t.Errorf("missing %s", p)
		}
	}
	// The lock records the template's topics on each convention (S-0134).
	if lk, err := lock.Load(dest); err != nil || lk == nil || strings.Join(lk.Topics["design/conventions/git.md"], ",") != "all" {
		t.Errorf("lock topics: %+v %v", lk, err)
	}
	// No template syntax left behind, and every front matter block parses.
	err = filepath.WalkDir(dest, func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		b, _ := os.ReadFile(path)
		s := string(b)
		if strings.HasSuffix(path, ".tmpl") || strings.Contains(s, "{{") {
			t.Errorf("unrendered template content in %s", path)
		}
		if strings.HasSuffix(path, ".md") && strings.HasPrefix(s, "---\n") {
			fm := strings.SplitN(s, "---\n", 3)[1]
			var v map[string]any
			if err := yaml.Unmarshal([]byte(fm), &v); err != nil {
				t.Errorf("front matter in %s: %v", path, err)
			}
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	claude, _ := os.ReadFile(filepath.Join(dest, "CLAUDE.md"))
	if !strings.Contains(string(claude), "# sample") || !strings.Contains(string(claude), "system-flow:end-of-baseline") {
		t.Error("CLAUDE.md baseline not rendered as expected")
	}
}

func TestNewAsksForEachVariableAtATerminal(t *testing.T) {
	var out strings.Builder
	tty := true
	a := &app{out: &out, stdinIsTerminal: &tty,
		prompter: prompt.New(strings.NewReader("\n\nOwner\nnot a url\nhttps://example.org/o/r\n"), &out)}
	m := template.Manifest{Variables: []template.Variable{
		{Name: "project_name", Prompt: "Project name"},
		{Name: "owner", Required: true},
		{Name: repoURLVar, Prompt: "Repository URL"},
	}}
	vars, err := a.collectVars(m, newOptions{}, "demo")
	if err != nil {
		t.Fatalf("%v\n%s", err, &out)
	}
	if vars["project_name"] != "demo" || vars["owner"] != "Owner" || vars[repoURLVar] != "https://example.org/o/r" {
		t.Errorf("vars: %v", vars)
	}
	s := out.String()
	for _, want := range []string{"Project name [demo]: ", "  owner is required\n", "Repository URL: "} {
		if !strings.Contains(s, want) {
			t.Errorf("the prompts do not say %q:\n%s", want, s)
		}
	}
	if n := strings.Count(s, "Repository URL: "); n != 2 {
		t.Errorf("asked for the URL %d times, want 2 (one refused):\n%s", n, s)
	}
}
