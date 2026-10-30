package roles

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// templateRoles is the template's role files' folder, relative to this
// package.
const templateRoles = "../../../template/root/design/workflows/roles"

// S-0416, S-0419: flai carries the four baseline role files, each parsing
// with no problem at its path under the design folder, byte copies of the
// template's, front matter and body, since the body is the role's prompt;
// BaselineFor finds each by its verb.
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
		if string(embedded) != string(shipped) {
			t.Errorf("%s: the embedded file is not the template's, front matter and body; run: cp template/root/design/workflows/roles/%s.md flai/internal/roles/baseline/%s.md", name, name, name)
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

// S-0419: a project with no role files' folder runs the planner, the
// orchestrator, and the analyzer from the baseline flai carries; a project
// with its own planner.md runs that file, its additions below the marker
// with it; a role flai carries no baseline for and the project lacks is an
// error naming the file.
func TestBuiltinIsTheProjectsFileElseTheBaseline(t *testing.T) {
	bare := t.TempDir()
	for _, name := range []string{"planner", "orchestrator", "analyzer"} {
		r, fromBaseline, err := Builtin(bare, "plans", name)
		if err != nil || !fromBaseline || r.Name != name || r.Path != "plans/workflows/roles/"+name+".md" || r.Body == "" {
			t.Errorf("%s with no roles folder: %q at %q, from baseline %v, %v", name, r.Name, r.Path, fromBaseline, err)
		}
	}

	own := t.TempDir()
	shipped, err := os.ReadFile(filepath.Join(templateRoles, "planner.md"))
	if err != nil {
		t.Fatal(err)
	}
	const addition = "- Plan nothing on Fridays.\n"
	dir := filepath.Join(own, "plans", "workflows", "roles")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "planner.md"), append(shipped, addition...), 0o644); err != nil {
		t.Fatal(err)
	}
	r, fromBaseline, err := Builtin(own, "plans", "planner")
	if err != nil || fromBaseline || r.Name != "planner" {
		t.Fatalf("the project's planner.md: %q, from baseline %v, %v", r.Name, fromBaseline, err)
	}
	if !strings.HasPrefix(r.Additions, Marker) || !strings.HasSuffix(r.Additions, addition) || !strings.HasSuffix(r.Body, addition) {
		t.Errorf("the project's addition is not in its planner: additions %q", r.Additions)
	}
	if r, fromBaseline, err := Builtin(own, "plans", "orchestrator"); err != nil || !fromBaseline || r.Name != "orchestrator" {
		t.Errorf("the orchestrator the project lacks: %q, from baseline %v, %v", r.Name, fromBaseline, err)
	}

	if _, _, err := Builtin(bare, "plans", "reviewer"); err == nil || !strings.Contains(err.Error(), "plans/workflows/roles/reviewer.md") {
		t.Errorf("a role with neither file nor baseline: %v", err)
	}
}
