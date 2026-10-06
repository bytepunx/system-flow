package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/goccy/go-yaml"

	"github.com/bytepunx/system-flow/flai/internal/config"
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

// releasedTemplate makes a git template, as a file:// URL, with the release
// tags v1.0.18 and v1.0.60, a newer untagged commit on main (1.0.61), and a
// branch edge at 2.0.0. Each version's team.txt says which it is.
func releasedTemplate(t *testing.T) string {
	t.Helper()
	if testing.Short() {
		t.Skip("integration: builds a git template")
	}
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}
	tpl := copyTemplate(t)
	gitIn(t, tpl, "init", "-q", "-b", "main")
	commit := func(version, tag string) {
		forkVersion(t, tpl, version, "", "v"+version+"\n")
		gitIn(t, tpl, "add", "-A")
		gitIn(t, tpl, "commit", "-q", "-m", version)
		if tag != "" {
			gitIn(t, tpl, "tag", tag)
		}
	}
	commit("1.0.18", "v1.0.18")
	commit("1.0.60", "v1.0.60")
	commit("1.0.61", "")
	gitIn(t, tpl, "checkout", "-q", "-b", "edge")
	commit("2.0.0", "")
	gitIn(t, tpl, "checkout", "-q", "main")
	return "file://" + tpl
}

// madeAt checks the project in dest was rendered from version, and that
// system-flow.yaml and the lock record the repo, the ref, and the version.
func madeAt(t *testing.T, dest, url, ref, version string) {
	t.Helper()
	if b, _ := os.ReadFile(filepath.Join(dest, "team.txt")); string(b) != "v"+version+"\n" {
		t.Errorf("team.txt = %q, want the render of %s", b, version)
	}
	m, err := manifest.Load(filepath.Join(dest, manifest.File))
	if err != nil || m.Template.Repo != url || m.Template.Ref != ref || m.Template.Version != version {
		t.Errorf("system-flow.yaml must record %s at %s from %s: %+v %v", ref, version, url, m.Template, err)
	}
	lk, err := lock.Load(dest)
	if err != nil || lk == nil || lk.Template.Repo != url || lk.Template.Ref != ref || lk.Template.Version != version {
		t.Errorf("the lock must record %s at %s from %s: %+v %v", ref, version, url, lk, err)
	}
}

// With no --ref, a new project is made at the template's newest release
// when the configured ref follows releases; --ref, as a tag or as the
// version it spells, wins; a configured ref naming another branch is used
// as given (ADR-0103).
func TestNewMakesTheProjectAtTheNewestRelease(t *testing.T) {
	url := releasedTemplate(t)
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	if _, errOut, code := runIn(t, ".", "template", "use", url, "--ref", "main"); code != 0 {
		t.Fatalf("template use: %s", errOut)
	}

	dest := filepath.Join(t.TempDir(), "newest")
	out, errOut, code := runIn(t, ".", "new", dest, "--defaults", "--no-git")
	if code != 0 {
		t.Fatalf("new: %s", errOut)
	}
	if !strings.Contains(out, "(1.0.60)") || !strings.Contains(out, "at v1.0.60, the template's newest release") {
		t.Errorf("new does not say it made the newest release:\n%s", out)
	}
	madeAt(t, dest, url, "v1.0.60", "1.0.60")

	// --template with no --ref follows releases too, and --json says so.
	dest = filepath.Join(t.TempDir(), "json")
	out, errOut, code = runIn(t, ".", "new", dest, "--template", url, "--defaults", "--no-git", "--json")
	var res struct {
		Ref     string `json:"ref"`
		Version string `json:"version"`
		Newest  bool   `json:"newest_release"`
	}
	if code != 0 || json.Unmarshal([]byte(out), &res) != nil || res.Ref != "v1.0.60" || res.Version != "1.0.60" || !res.Newest {
		t.Errorf("new --template --json: %d %+v %s %s", code, res, out, errOut)
	}

	for _, ref := range []string{"v1.0.18", "1.0.18"} {
		dest := filepath.Join(t.TempDir(), "old")
		out, errOut, code := runIn(t, ".", "new", dest, "--ref", ref, "--defaults", "--no-git")
		if code != 0 {
			t.Fatalf("new --ref %s: %s", ref, errOut)
		}
		if strings.Contains(out, "newest release") {
			t.Errorf("new --ref %s calls it the newest release:\n%s", ref, out)
		}
		madeAt(t, dest, url, "v1.0.18", "1.0.18")
	}

	if _, errOut, code := runIn(t, ".", "template", "use", url, "--ref", "edge"); code != 0 {
		t.Fatalf("template use --ref edge: %s", errOut)
	}
	dest = filepath.Join(t.TempDir(), "edge")
	if _, errOut, code := runIn(t, ".", "new", dest, "--defaults", "--no-git"); code != 0 {
		t.Fatalf("new on edge: %s", errOut)
	}
	madeAt(t, dest, url, "edge", "2.0.0")
}

// A git template with no release tags is used at the configured ref.
func TestNewWithoutReleaseTagsUsesTheConfiguredRef(t *testing.T) {
	if testing.Short() {
		t.Skip("integration: builds a git template")
	}
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	tpl := copyTemplate(t)
	forkVersion(t, tpl, "1.2.3", "", "v1.2.3\n")
	gitIn(t, tpl, "init", "-q", "-b", "main")
	gitIn(t, tpl, "add", "-A")
	gitIn(t, tpl, "commit", "-q", "-m", "untagged")
	url := "file://" + tpl
	if _, errOut, code := runIn(t, ".", "template", "use", url, "--ref", "main"); code != 0 {
		t.Fatalf("template use: %s", errOut)
	}
	dest := filepath.Join(t.TempDir(), "p")
	out, errOut, code := runIn(t, ".", "new", dest, "--defaults", "--no-git")
	if code != 0 || strings.Contains(out, "newest release") {
		t.Fatalf("new: %d %s %s", code, out, errOut)
	}
	madeAt(t, dest, url, "main", "1.2.3")
}

// staleRunner answers for a cached clone of a branch whose fetch fails, as
// offline: git is installed, the clone is on main, and fetch errors.
type staleRunner struct{ calls []string }

func (r *staleRunner) Run(dir, name string, args ...string) (string, error) {
	r.calls = append(r.calls, name+" "+strings.Join(args, " "))
	switch {
	case name == "git" && args[0] == "rev-parse":
		return "main\n", nil
	case name == "git" && args[0] == "fetch":
		return "", fmt.Errorf("fatal: unable to access the remote: network unreachable")
	}
	return "", fmt.Errorf("unexpected %s %v", name, args)
}

func (r *staleRunner) RunInput(dir, name, _ string, args ...string) (string, error) {
	return r.Run(dir, name, args...)
}

func (r *staleRunner) LookPath(name string) (string, error) { return "/usr/bin/" + name, nil }

// A cached branch clone that cannot be fetched again is used, with a
// warning, rather than failing the command (ADR-0103).
func TestResolveTemplateWarnsAndUsesAStaleClone(t *testing.T) {
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	cache := t.TempDir()
	t.Setenv(config.CacheEnvVar, cache)
	repo := "https://example.invalid/template.git"
	src, err := template.Resolve(repo, "main", cache)
	if err != nil {
		t.Fatal(err)
	}
	clone := copyTemplate(t)
	if err := os.MkdirAll(filepath.Dir(src.Dir), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Rename(clone, src.Dir); err != nil {
		t.Fatal(err)
	}
	var out, errOut strings.Builder
	r := &staleRunner{}
	a := &app{out: &out, errOut: &errOut, runner: r}
	got, m, err := a.resolveTemplate(repo, "main", false)
	if err != nil {
		t.Fatalf("a stale clone must not fail: %v", err)
	}
	if got.Dir != src.Dir || m.Version != "9.9.9" {
		t.Errorf("resolved %+v at %s, want the clone in %s", got, m.Version, src.Dir)
	}
	if !strings.Contains(errOut.String(), "using the cached template clone") || !strings.Contains(errOut.String(), "network unreachable") {
		t.Errorf("no warning naming the failed fetch:\n%s", errOut.String())
	}
	if !strings.Contains(strings.Join(r.calls, "\n"), "git fetch") {
		t.Errorf("the clone was not fetched again: %v", r.calls)
	}
}
