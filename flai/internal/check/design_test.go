package check

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

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
