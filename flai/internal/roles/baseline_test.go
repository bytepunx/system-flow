package roles

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// templateRoles is the template's role files' folder, relative to this
// package.
const templateRoles = "../../../template/root/design/workflows/roles"

// S-0416: flai carries the four baseline role files, each parsing with no
// problem at its path under the design folder, with the front matter the
// template ships; BaselineFor finds each by its verb.
func TestBaselineIsTheTemplatesRoles(t *testing.T) {
	got := Baseline("plans")
	if len(got) != len(baselineNames) {
		t.Fatalf("%d baseline roles, want %d", len(got), len(baselineNames))
	}
	for i, name := range baselineNames {
		r := got[i]
		if r.Name != name || r.Path != "plans/workflows/roles/"+name+".md" {
			t.Errorf("role %d: name %q at %q", i, r.Name, r.Path)
		}
		if problems := Validate([]Role{r}, "plans", Vocabulary{}); len(problems) != 0 {
			t.Errorf("%s: %v", name, problems)
		}
		embedded, err := baselineFS.ReadFile("baseline/" + name + ".md")
		if err != nil {
			t.Fatal(err)
		}
		shipped, err := os.ReadFile(filepath.Join(templateRoles, name+".md"))
		if err != nil {
			t.Fatal(err)
		}
		if a, b := frontMatter(t, string(embedded)), frontMatter(t, string(shipped)); a != b {
			t.Errorf("%s: the embedded front matter is not the template's; copy %s over baseline/%s.md:\n%s\nwant:\n%s", name, filepath.Join(templateRoles, name+".md"), name, a, b)
		}
		found, ok := BaselineFor("plans", r.Role)
		if !ok || found.Name != name {
			t.Errorf("BaselineFor(%q): %q, %v", r.Role, found.Name, ok)
		}
	}
	if _, ok := BaselineFor("plans", "review"); ok {
		t.Error("BaselineFor found a role flai does not carry")
	}
}

// frontMatter is the text between a file's opening and closing ---.
func frontMatter(t *testing.T, content string) string {
	t.Helper()
	fm, _, err := workitem.SplitFrontMatter(content)
	if err != nil {
		t.Fatal(err)
	}
	return fm
}
