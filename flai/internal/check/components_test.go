package check

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// S-0437: a sub-project's components key that is not valid is an error on
// the projects key, naming the sub-project's index and field; a valid key,
// and none, are not reported.
func TestComponentsKeysAreChecked(t *testing.T) {
	const web = "projects:\n  - name: web\n    path: web\n    kind: sveltekit\n"
	for _, extra := range []string{
		web,
		web + "    components:\n      include: [\"src/lib/**/*.svelte\"]\n      setup: src/app.css\n",
	} {
		repo := planningProject(t, extra)
		if err := os.MkdirAll(filepath.Join(repo.Root, "web", "src"), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(repo.Root, "web", "src", "app.css"), nil, 0o644); err != nil {
			t.Fatal(err)
		}
		if got := findings(t, repo, "manifest.components"); len(got) != 0 {
			t.Errorf("valid components %q: %+v", extra, got)
		}
	}
	repo := planningProject(t, web+"    components:\n      setup: src/missing.css\n")
	got := findings(t, repo, "manifest.components")
	if len(got) != 1 {
		t.Fatalf("want one error, the setup: %+v", got)
	}
	if f := got[0]; f.Level != Error || f.Path != "system-flow.yaml" || f.Line != 8 || !strings.Contains(f.Message, `projects[0].components.setup (sub-project "web") "src/missing.css" names no file`) {
		t.Errorf("finding: %+v, want an error on line 8 naming projects[0].components.setup", f)
	}
}
