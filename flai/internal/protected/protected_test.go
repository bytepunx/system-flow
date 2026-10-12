package protected

import (
	"path/filepath"
	"slices"
	"testing"
)

// every returns what of every harness's list, harness by harness.
func every(what func(List) []string) []string {
	var out []string
	for _, h := range Harnesses() {
		out = append(out, what(For(h))...)
	}
	return out
}

func TestPathMatchesEveryProtectedFolderAtTheTopAndNested(t *testing.T) {
	for _, f := range every(func(l List) []string { return l.Folders }) {
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
	for _, f := range every(func(l List) []string { return l.Files }) {
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

func TestForPinsClaudeCodesList(t *testing.T) {
	l := For("claude-code")
	folders := []string{".git", ".config/git", ".vscode", ".idea", ".husky", ".cargo", ".devcontainer", ".yarn", ".mvn", ".claude"}
	files := []string{
		".gitconfig", ".gitmodules",
		".bashrc", ".bash_profile", ".bash_login", ".bash_aliases", ".bash_logout", ".zshrc", ".zprofile", ".zshenv", ".zlogin", ".zlogout", ".profile", ".envrc",
		".npmrc", ".yarnrc", ".yarnrc.yml", ".pnp.cjs", ".pnp.loader.mjs", ".pnpmfile.cjs", "bunfig.toml", ".bunfig.toml",
		".bazelrc", ".bazelversion", ".bazeliskrc",
		".pre-commit-config.yaml", "lefthook.yml", "lefthook.yaml", ".lefthook.yml", ".lefthook.yaml",
		"gradle-wrapper.properties", "maven-wrapper.properties",
		".devcontainer.json",
		".ripgreprc", "pyrightconfig.json",
		".mcp.json", ".claude.json",
	}
	if !slices.Equal(l.Folders, folders) {
		t.Errorf("For(claude-code).Folders = %q, want %q", l.Folders, folders)
	}
	if !slices.Equal(l.Files, files) {
		t.Errorf("For(claude-code).Files = %q, want %q", l.Files, files)
	}
}

func TestForReturnsACopy(t *testing.T) {
	l := For(ClaudeCode)
	l.Folders[0], l.Files[0] = "x", "y"
	if again := For(ClaudeCode); again.Folders[0] != ".git" || again.Files[0] != ".gitconfig" {
		t.Errorf("For(%s) after changing a copy = %q, %q: the copy changed the list", ClaudeCode, again.Folders[0], again.Files[0])
	}
}

func TestHarnessesIsClaudeCode(t *testing.T) {
	if got := Harnesses(); !slices.Equal(got, []string{"claude-code"}) {
		t.Errorf("Harnesses() = %q, want [claude-code]", got)
	}
}

func TestForAnUnknownHarnessIsEmpty(t *testing.T) {
	if l := For("codex"); len(l.Folders) != 0 || len(l.Files) != 0 {
		t.Errorf("For(codex) = %+v, want an empty list", l)
	}
}

func TestPathAnswersForTheUnion(t *testing.T) {
	saved := lists
	t.Cleanup(func() { lists = saved })
	lists = map[string]List{
		ClaudeCode: saved[ClaudeCode],
		"other":    {Folders: []string{".other"}, Files: []string{"OTHER.md"}},
	}
	for _, p := range []string{".other/x", "a/.other/b.json", "OTHER.md", "sub/OTHER.md", ".claude/x.md", ".mcp.json"} {
		if !Path(p) {
			t.Errorf("Path(%q) = false, want true: it is in one harness's list", p)
		}
	}
	if got, want := Changed([]string{"flai/main.go", "OTHER.md", ".mcp.json"}), []string{"OTHER.md", ".mcp.json"}; !slices.Equal(got, want) {
		t.Errorf("Changed = %q, want %q", got, want)
	}
	if Path("flai/main.go") || Path(".other") {
		t.Errorf("Path of a path in no harness's list = true, want false")
	}
}
