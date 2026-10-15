package check

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/designsystem"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// fullTokens is designsystem's token file with every role and step, which
// gives no problem.
func fullTokens(t *testing.T) string {
	t.Helper()
	data, err := os.ReadFile("../designsystem/testdata/full.tokens.json")
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}

// writeTokens writes body to the repository path name under repo's root.
func writeTokens(t *testing.T, repo *workitem.Repo, name, body string) {
	t.Helper()
	path := filepath.Join(repo.Root, filepath.FromSlash(name))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}

// lineOf is the 1-based line of body's first line holding s.
func lineOf(t *testing.T, body, s string) int {
	t.Helper()
	for i, l := range strings.Split(body, "\n") {
		if strings.Contains(l, s) {
			return i + 1
		}
	}
	t.Fatalf("no line holds %q", s)
	return 0
}

// wantTokenFinding fails unless got has an error at name and line whose
// message starts with prefix.
func wantTokenFinding(t *testing.T, got []Finding, name string, line int, prefix string) {
	t.Helper()
	for _, f := range got {
		if f.Level == Error && f.Path == filepath.FromSlash(name) && f.Line == line && strings.HasPrefix(f.Message, prefix) {
			return
		}
	}
	t.Errorf("want an error at %s:%d starting %q in %+v", name, line, prefix, got)
}

// ADR-0148: a project with no token file passes, and one whose token file has
// every role and step has no design.tokens finding.
func TestDesignTokensNoneAndValid(t *testing.T) {
	repo := planningProject(t, "")
	if got := findings(t, repo, "design.tokens"); len(got) != 0 {
		t.Errorf("no token file: %+v, want no finding", got)
	}
	writeTokens(t, repo, manifest.DefaultDesignTokens, fullTokens(t))
	if got := findings(t, repo, "design.tokens"); len(got) != 0 {
		t.Errorf("a valid token file: %+v, want no finding", got)
	}
}

// ADR-0148: a colour space flai does not read is an error at the token's line
// on the token file, its message naming the token's path first.
func TestDesignTokensTokenProblem(t *testing.T) {
	repo := planningProject(t, "")
	good := `"brand-sun": { "$value": { "colorSpace": "srgb", "components": [0.961, 0.718, 0], "hex": "#f5b700" } }`
	body := fullTokens(t)
	if !strings.Contains(body, good) {
		t.Fatalf("the full token file has no %s", good)
	}
	body = strings.Replace(body, good, `"brand-sun": { "$value": { "colorSpace": "display-p3", "components": [0.961, 0.718, 0] } }`, 1)
	writeTokens(t, repo, manifest.DefaultDesignTokens, body)
	got := findings(t, repo, "design.tokens")
	wantTokenFinding(t, got, manifest.DefaultDesignTokens, lineOf(t, body, `"brand-sun"`), "color.brand-sun: ")
	for _, f := range got {
		if !strings.HasPrefix(f.Message, "color.brand-sun: ") {
			t.Errorf("finding %+v, want only brand-sun's", f)
		}
	}
}

// ADR-0148: a resolver document with a modifier other than theme is an error
// at the modifier's line, named by its path.
func TestDesignTokensResolverProblem(t *testing.T) {
	repo := planningProject(t, "")
	body := `{
  "version": "2025.10",
  "sets": { "tokens": { "sources": [
` + fullTokens(t) + `  ] } },
  "modifiers": {
    "theme": { "contexts": { "light": [], "dark": [] }, "default": "light" },
    "density": { "contexts": { "compact": [] } }
  },
  "resolutionOrder": [ { "$ref": "#/sets/tokens" }, { "$ref": "#/modifiers/theme" } ]
}`
	writeTokens(t, repo, manifest.DefaultDesignTokens, body)
	got := findings(t, repo, "design.tokens")
	wantTokenFinding(t, got, manifest.DefaultDesignTokens, lineOf(t, body, `"density"`), "modifiers.density: ")
}

// ADR-0146: design.tokens names the token file; a file at the default path is
// then not read. A file that does not parse is one error at its line, at no
// token path; one flai cannot read is one error on its path.
func TestDesignTokensManifestPath(t *testing.T) {
	const name = "brand/tokens.json"
	repo := planningProject(t, "design:\n  tokens: "+name+"\n")
	writeTokens(t, repo, manifest.DefaultDesignTokens, "{")
	writeTokens(t, repo, name, "{\n  \"color\": \n}")
	got := findings(t, repo, "design.tokens")
	if len(got) != 1 {
		t.Fatalf("want one error, the named file's: %+v", got)
	}
	wantTokenFinding(t, got, name, 3, "the token file does not parse: ")

	repo = planningProject(t, "")
	if err := os.MkdirAll(filepath.Join(repo.Root, filepath.FromSlash(manifest.DefaultDesignTokens)), 0o755); err != nil {
		t.Fatal(err)
	}
	got = findings(t, repo, "design.tokens")
	if len(got) != 1 {
		t.Fatalf("a token file that is a folder: %+v, want one error", got)
	}
	wantTokenFinding(t, got, manifest.DefaultDesignTokens, 1, "read the token file "+manifest.DefaultDesignTokens)
}

// stylesheetProject is an empty project with a web sub-project whose
// manifest names the design stylesheet web/src/theme.css, and a full token
// file at the default path.
func stylesheetProject(t *testing.T) *workitem.Repo {
	t.Helper()
	repo := planningProject(t, "projects:\n  - name: web\n    path: web\n    kind: sveltekit\ndesign:\n  stylesheet: web/src/theme.css\n")
	writeTokens(t, repo, manifest.DefaultDesignTokens, fullTokens(t))
	return repo
}

// inStep is the stylesheet flai design css --write writes for repo's token
// file at the default path.
func inStep(t *testing.T, repo *workitem.Repo) string {
	t.Helper()
	sys, err := designsystem.Load(os.DirFS(repo.Root), manifest.DefaultDesignTokens)
	if err != nil {
		t.Fatal(err)
	}
	css := designsystem.Stylesheet(sys, manifest.DefaultDesignTokens)
	if css == "" {
		t.Fatal("the full token file gives no stylesheet")
	}
	return css
}

// wantStale fails unless got is two design.stale warnings, one on
// web/src/theme.css and one on the token file, each at line 1 and holding
// want and the command.
func wantStale(t *testing.T, got []Finding, want string) {
	t.Helper()
	if len(got) != 2 {
		t.Fatalf("want two design.stale findings: %+v", got)
	}
	paths := map[string]bool{}
	for _, f := range got {
		if f.Level != Warning || f.Line != 1 || !strings.Contains(f.Message, want) || !strings.Contains(f.Message, "flai design css --write") {
			t.Errorf("finding %+v, want a warning at line 1 holding %q and the command", f, want)
		}
		paths[filepath.ToSlash(f.Path)] = true
	}
	if !paths["web/src/theme.css"] || !paths[manifest.DefaultDesignTokens] {
		t.Errorf("findings on %v, want the stylesheet and the token file", paths)
	}
}

// ADR-0152, ADR-0159: a stylesheet that is what flai design css --write
// writes, with LF or CRLF line endings, gives no design.stale finding.
func TestDesignStaleInStep(t *testing.T) {
	repo := stylesheetProject(t)
	css := inStep(t, repo)
	writeTokens(t, repo, "web/src/theme.css", css)
	if got := findings(t, repo, "design.stale"); len(got) != 0 {
		t.Errorf("in step: %+v, want no finding", got)
	}
	writeTokens(t, repo, "web/src/theme.css", strings.ReplaceAll(css, "\n", "\r\n"))
	if got := findings(t, repo, "design.stale"); len(got) != 0 {
		t.Errorf("in step with CRLF: %+v, want no finding", got)
	}
	if got := findings(t, repo, "design.stylesheet"); len(got) != 0 {
		t.Errorf("a stylesheet under a sub-project: %+v, want no design.stylesheet finding", got)
	}
}

// ADR-0152 § Decision 2: a stylesheet that differs from the tokens is a
// warning on the stylesheet and one on the token file, so that a branch that
// changed only the token file stops at its close-out.
func TestDesignStaleDiffers(t *testing.T) {
	repo := stylesheetProject(t)
	writeTokens(t, repo, "web/src/theme.css", inStep(t, repo)+"/* hand edit */\n")
	wantStale(t, findings(t, repo, "design.stale"), "does not match the tokens in "+manifest.DefaultDesignTokens)
}

// ADR-0152: a stylesheet the manifest names that is not there is the same
// pair of warnings, saying it is missing.
func TestDesignStaleMissing(t *testing.T) {
	repo := stylesheetProject(t)
	wantStale(t, findings(t, repo, "design.stale"), "web/src/theme.css for the tokens in "+manifest.DefaultDesignTokens+" is missing")
}

// A stylesheet flai cannot read is one error on it.
func TestDesignStaleUnreadable(t *testing.T) {
	repo := stylesheetProject(t)
	if err := os.MkdirAll(filepath.Join(repo.Root, "web", "src", "theme.css"), 0o755); err != nil {
		t.Fatal(err)
	}
	got := findings(t, repo, "design.stale")
	if len(got) != 1 || got[0].Level != Error || filepath.ToSlash(got[0].Path) != "web/src/theme.css" {
		t.Errorf("a stylesheet that is a folder: %+v, want one error on it", got)
	}
}

// ADR-0152: no design.stale finding when the manifest names no stylesheet,
// when there is no token file, or when the token file has problems, which
// design.tokens reports alone.
func TestDesignStaleNotChecked(t *testing.T) {
	repo := planningProject(t, "projects:\n  - name: web\n    path: web\n    kind: sveltekit\n")
	writeTokens(t, repo, manifest.DefaultDesignTokens, fullTokens(t))
	if got := findings(t, repo, "design.stale"); len(got) != 0 {
		t.Errorf("no stylesheet named: %+v, want no finding", got)
	}

	repo = planningProject(t, "projects:\n  - name: web\n    path: web\n    kind: sveltekit\ndesign:\n  stylesheet: web/src/theme.css\n")
	if got := findings(t, repo, "design.stale"); len(got) != 0 {
		t.Errorf("no token file: %+v, want no finding", got)
	}

	body := strings.Replace(fullTokens(t), `"colorSpace": "srgb"`, `"colorSpace": "display-p3"`, 1)
	writeTokens(t, repo, manifest.DefaultDesignTokens, body)
	if got := findings(t, repo, "design.tokens"); len(got) == 0 {
		t.Fatal("the token file has no problem to report")
	}
	if got := findings(t, repo, "design.stale"); len(got) != 0 {
		t.Errorf("a token file with problems: %+v, want no design.stale finding", got)
	}
}

// ADR-0152: a stylesheet outside every sub-project is an error on the
// manifest at the key's line, and is not read.
func TestDesignStylesheetOutsideProjects(t *testing.T) {
	extra := "projects:\n  - name: web\n    path: web\n    kind: sveltekit\ndesign:\n  tokens: " + manifest.DefaultDesignTokens + "\n  stylesheet: site/theme.css\n"
	repo := planningProject(t, extra)
	writeTokens(t, repo, manifest.DefaultDesignTokens, fullTokens(t))
	got := findings(t, repo, "design.stylesheet")
	body, err := os.ReadFile(filepath.Join(repo.Root, manifest.File))
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got[0].Level != Error || got[0].Path != manifest.File || got[0].Line != lineOf(t, string(body), "stylesheet:") || !strings.Contains(got[0].Message, `"site/theme.css" lies outside every sub-project`) {
		t.Errorf("a stylesheet outside every sub-project: %+v, want one error on the manifest at its key", got)
	}
	if got := findings(t, repo, "design.stale"); len(got) != 0 {
		t.Errorf("a stylesheet outside every sub-project: %+v, want no design.stale finding", got)
	}
}
