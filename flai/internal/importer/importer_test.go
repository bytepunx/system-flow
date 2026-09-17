package importer

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/execx"
)

// legacy builds a repository that predates the standard.
func legacy(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	w := func(rel, body string) {
		p := filepath.Join(root, rel)
		_ = os.MkdirAll(filepath.Dir(p), 0o755)
		_ = os.WriteFile(p, []byte(body), 0o644)
	}
	w("README.md", "# legacy\n")
	w("CHANGELOG.md", "# changes\n")
	w("notes.md", "# notes\n")
	w("docs/guide.md", "# guide\n")
	w("adr/0001-first.md", "# first\n")
	w("adr/nested/0002-second.md", "# second\n")
	w("services/api/go.mod", "module example.com/api\n")
	w("services/api/README.md", "# api\n")
	w("web/package.json", "{}")
	w("web/svelte.config.js", "export default {}")
	w("web/node_modules/pkg/README.md", "ignored")
	w("tools/cli/Cargo.toml", "[package]\n")
	w("misc/deep/er/too/deep.md", "# deep\n")
	return root
}

func TestScanAndPlan(t *testing.T) {
	root := legacy(t)
	a, err := Scan(root)
	if err != nil {
		t.Fatal(err)
	}
	if a.Git || a.Manifest {
		t.Errorf("git=%v manifest=%v", a.Git, a.Manifest)
	}
	if a.Existing["docs"] != "docs" || a.Existing["design"] != "" {
		t.Errorf("existing: %v", a.Existing)
	}
	if a.Candidates["adr"] != "design/adrs" {
		t.Errorf("candidates: %v", a.Candidates)
	}
	kinds := map[string]string{}
	for _, p := range a.Projects {
		kinds[p.Path] = p.Kind
	}
	if kinds["services/api"] != "go" || kinds["web"] != "sveltekit" || kinds["tools/cli"] != "rust" || len(a.Projects) != 3 {
		t.Errorf("projects: %+v", a.Projects)
	}
	want := []string{"misc/deep/er/too/deep.md", "notes.md"}
	if len(a.Markdown) != 2 || a.Markdown[0] != want[0] || a.Markdown[1] != want[1] {
		t.Errorf("markdown: %v (README, CHANGELOG, docs/, adr/, project files must be excluded)", a.Markdown)
	}
	plan := BuildPlan(a, Layout{"design": "design", "docs": "docs", "wip": "wip"})
	if len(plan.FolderMoves) != 1 || plan.FolderMoves[0].To != "design/adrs" {
		t.Errorf("folder moves: %+v", plan.FolderMoves)
	}
	creates := map[string]bool{}
	for _, c := range plan.Create {
		creates[c] = true
	}
	if !creates["design"] || !creates["wip/kanban/tasks"] || creates["docs"] || !creates["docs/users"] || !creates["scripts"] {
		t.Errorf("create: %v", plan.Create)
	}
	renamed := BuildPlan(a, Layout{"design": "architecture", "docs": "docs", "wip": "wip"})
	if renamed.FolderMoves[0].To != "architecture/adrs" || !contains(renamed.Create, "architecture/conventions") {
		t.Errorf("renamed layout: %+v", renamed)
	}
}

func TestMoveFolderContentsWithGit(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}
	root := legacy(t)
	r := execx.System{}
	for _, args := range [][]string{{"init", "-q", "-b", "main"}, {"config", "user.email", "t@t"}, {"config", "user.name", "t"}, {"add", "-A"}, {"commit", "-q", "-m", "legacy"}} {
		if _, err := r.Run(root, "git", args...); err != nil {
			t.Fatal(err)
		}
	}
	_ = os.WriteFile(filepath.Join(root, "adr", "untracked.md"), []byte("# new\n"), 0o644)
	_ = os.MkdirAll(filepath.Join(root, "design", "adrs"), 0o755)
	_ = os.WriteFile(filepath.Join(root, "design", "adrs", "0001-first.md"), []byte("# existing\n"), 0o644)
	moved, err := MoveFolderContents(r, root, "adr", "design/adrs", true)
	if err != nil {
		t.Fatal(err)
	}
	if len(moved) != 2 {
		t.Fatalf("moved %+v; the conflicting 0001 must be left in place", moved)
	}
	if b, _ := os.ReadFile(filepath.Join(root, "design", "adrs", "0001-first.md")); string(b) != "# existing\n" {
		t.Error("existing file was overwritten")
	}
	if _, err := os.Stat(filepath.Join(root, "adr", "0001-first.md")); err != nil {
		t.Error("conflicting source should remain")
	}
	out, _ := r.Run(root, "git", "status", "--short")
	if !containsStr(out, "R  adr/nested/0002-second.md -> design/adrs/nested/0002-second.md") {
		t.Errorf("tracked file should be git mv'd:\n%s", out)
	}
	if _, err := os.Stat(filepath.Join(root, "design", "adrs", "untracked.md")); err != nil {
		t.Error("untracked file should be renamed")
	}
}

func contains(list []string, s string) bool {
	for _, x := range list {
		if x == s {
			return true
		}
	}
	return false
}

func containsStr(s, sub string) bool {
	return len(sub) > 0 && len(s) >= len(sub) && (func() bool { return indexOf(s, sub) >= 0 })()
}

func indexOf(s, sub string) int {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return i
		}
	}
	return -1
}
