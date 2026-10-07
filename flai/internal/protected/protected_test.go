package protected

import (
	"path/filepath"
	"slices"
	"testing"
)

func TestPathMatchesEveryProtectedFolderAtTheTopAndNested(t *testing.T) {
	for _, f := range folders {
		for _, p := range []string{f + "/x", f + "/a/b.json", "template/root/" + f + "/x.md"} {
			if !Path(p) {
				t.Errorf("Path(%q) = false, want true: %s is a protected folder", p, f)
			}
		}
	}
	for _, p := range []string{".claude/agents/x.md", "template/root/.claude/agents/x.md", ".config/git/config", "home/.config/git/ignore"} {
		if !Path(p) {
			t.Errorf("Path(%q) = false, want true", p)
		}
	}
}

func TestPathMatchesEveryProtectedFileAtTheTopAndNested(t *testing.T) {
	for _, f := range files {
		for _, p := range []string{f, "template/root/" + f, "a/b/" + f} {
			if !Path(p) {
				t.Errorf("Path(%q) = false, want true: %s is a protected file", p, f)
			}
		}
	}
}

func TestPathLeavesOrdinaryPaths(t *testing.T) {
	for _, p := range []string{
		"", ".", "flai/main.go", "docs/x.md", ".claudeignore", "my.mcp.json", ".mcp.json.bak",
		".config/other", ".config/gitx/config", "git/.config", "git/.config/x", ".config/git",
		".claude", ".vscode", "flai/claude/settings.json", "vscode/settings.json",
	} {
		if Path(p) {
			t.Errorf("Path(%q) = true, want false", p)
		}
	}
}

func TestGitFindsAGitFolderOrFile(t *testing.T) {
	for p, want := range map[string]bool{
		".git": true, ".git/config": true, "sub/.git": true, "sub/.git/hooks/pre-commit": true,
		".gitignore": false, ".github/workflows/x.yml": false, "flai/main.go": false, "": false,
	} {
		if Git(p) != want {
			t.Errorf("Git(%q) = %v, want %v", p, !want, want)
		}
		if want && !Path(p) {
			t.Errorf("Path(%q) = false: a .git path is protected", p)
		}
	}
}

func TestPathTakesEitherSeparator(t *testing.T) {
	for _, p := range []string{filepath.Join("template", "root", ".mcp.json"), filepath.Join(".claude", "settings.json"), "./.vscode/settings.json", "a//.idea/x"} {
		if !Path(p) {
			t.Errorf("Path(%q) = false, want true", p)
		}
	}
}

func TestChangedKeepsTheProtectedInOrder(t *testing.T) {
	in := []string{"flai/main.go", ".mcp.json", "docs/x.md", "template/root/.claude/settings.json", ".claudeignore", ".vscode/settings.json"}
	want := []string{".mcp.json", "template/root/.claude/settings.json", ".vscode/settings.json"}
	if got := Changed(in); !slices.Equal(got, want) {
		t.Errorf("Changed = %q, want %q", got, want)
	}
	if got := Changed([]string{"flai/main.go"}); got != nil {
		t.Errorf("Changed of no protected path = %q, want nil", got)
	}
}
