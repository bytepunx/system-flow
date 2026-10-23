package check

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// S-0491: a sub-project entry that is not valid, such as a name another
// entry has, is an error on the projects key, naming the entry's index and
// field; valid entries are not reported, and a path that names no folder is.
func TestProjectEntriesAreChecked(t *testing.T) {
	const two = "projects:\n  - name: web\n    path: web\n    kind: sveltekit\n  - name: api\n    path: api\n    kind: go\n"
	repo := planningProject(t, two)
	for _, d := range []string{"web", "api"} {
		if err := os.MkdirAll(filepath.Join(repo.Root, d), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	if got := findings(t, repo, "manifest.projects"); len(got) != 0 {
		t.Errorf("valid projects: %+v", got)
	}
	repo = planningProject(t, strings.Replace(two, "name: api", "name: web", 1))
	for _, d := range []string{"web", "api"} {
		if err := os.MkdirAll(filepath.Join(repo.Root, d), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	got := findings(t, repo, "manifest.projects")
	if len(got) != 1 {
		t.Fatalf("want one error, the repeated name: %+v", got)
	}
	if f := got[0]; f.Level != Error || f.Path != "system-flow.yaml" || f.Line != 8 || !strings.Contains(f.Message, `projects[1].name (sub-project "web") repeats the name of projects[0]`) {
		t.Errorf("finding: %+v, want an error on line 8 naming projects[1].name", f)
	}
	got = findings(t, planningProject(t, two), "manifest.projects")
	if len(got) != 2 || !strings.Contains(got[0].Message, `projects[0].path (sub-project "web") "web" names no folder`) {
		t.Errorf("want each missing folder reported: %+v", got)
	}
}
