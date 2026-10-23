package cmd

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/manifest"
)

// subprojectProject is a scratch project whose manifest says projects: [],
// and the manifest's text.
func subprojectProject(t *testing.T) (string, string) {
	t.Helper()
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	t.Setenv("FLAI_AGENT", "tester")
	root := tempProject(t)
	file := filepath.Join(root, manifest.File)
	data, _ := os.ReadFile(file)
	data = append(data, "# the sub-projects\nprojects: []\n"...)
	if err := os.WriteFile(file, data, 0o644); err != nil {
		t.Fatal(err)
	}
	return root, string(data)
}

// S-0491: add to projects: [] writes the entry as a block, creates the
// folder with a README.md that names the sub-project and its kind, says what
// it wrote, and list shows the entry with its folder present.
func TestSubprojectAddCreatesTheFolderAndListShowsIt(t *testing.T) {
	root, before := subprojectProject(t)
	out := sharedRun(t, root, "subproject", "add", "web", "--kind", "sveltekit")
	if want := "added sub-project web (web, sveltekit)\ncreated folder web\nwrote system-flow.yaml\nwrote web/README.md\n"; out != want {
		t.Errorf("add says %q, want %q", out, want)
	}
	want := strings.Replace(before, "projects: []\n", "projects:\n  - name: web\n    path: web\n    kind: sveltekit\n", 1)
	if got := readManifest(t, root); got != want {
		t.Errorf("manifest:\n%s\nwant:\n%s", got, want)
	}
	readme, err := os.ReadFile(filepath.Join(root, "web", "README.md"))
	if err != nil {
		t.Fatalf("add created no README.md: %v", err)
	}
	if want := "# web\n\n`flai subproject add` created this folder for the sub-project web, of kind `sveltekit`.\n"; string(readme) != want {
		t.Errorf("README.md = %q, want %q", readme, want)
	}

	out = sharedRun(t, root, "subproject", "list")
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) != 2 || strings.Join(strings.Fields(lines[0]), " ") != "NAME PATH KIND TAGS FOLDER COMPONENTS INSTALL" ||
		strings.Join(strings.Fields(lines[1]), " ") != "web web sveltekit - present - -" {
		t.Errorf("list:\n%s", out)
	}
	var listed struct {
		Projects []listedSubproject `json:"projects"`
	}
	if err := json.Unmarshal([]byte(sharedRun(t, root, "subproject", "list", "--json")), &listed); err != nil {
		t.Fatal(err)
	}
	wantListed := []listedSubproject{{Name: "web", Path: "web", Kind: "sveltekit", Tags: []string{}, Folder: true}}
	if !reflect.DeepEqual(listed.Projects, wantListed) {
		t.Errorf("list --json = %+v, want %+v", listed.Projects, wantListed)
	}

	_ = os.RemoveAll(filepath.Join(root, "web"))
	if out := sharedRun(t, root, "subproject", "list"); !strings.Contains(out, "missing") {
		t.Errorf("list does not say the folder is missing:\n%s", out)
	}
}

// S-0491: list says so when there are no sub-projects, and --json gives [].
func TestSubprojectListEmpty(t *testing.T) {
	root, _ := subprojectProject(t)
	if out := sharedRun(t, root, "subproject", "list"); !strings.Contains(out, "no sub-projects") {
		t.Errorf("an empty list says so, got %q", out)
	}
	if out := sharedRun(t, root, "subproject", "list", "--json"); strings.Join(strings.Fields(out), "") != `{"projects":[]}` {
		t.Errorf("an empty list --json, got %q", out)
	}
}

// S-0491: add into a folder that is there writes no README.md, and --json
// gives the entry, the paths written, and no commit.
func TestSubprojectAddIntoAFolderThatIsThere(t *testing.T) {
	root, _ := subprojectProject(t)
	if err := os.MkdirAll(filepath.Join(root, "services", "api"), 0o755); err != nil {
		t.Fatal(err)
	}
	var got subprojectChange
	out := sharedRun(t, root, "subproject", "add", "api", "--path", "services/api", "--kind", "go", "--tags", "backend, http", "--json")
	if err := json.Unmarshal([]byte(out), &got); err != nil {
		t.Fatalf("%v: %s", err, out)
	}
	want := subprojectChange{Action: manifest.ProjectAdd, Project: manifest.Project{Name: "api", Path: "services/api", Kind: "go", Tags: []string{"backend", "http"}}, Written: []string{manifest.File}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("add --json = %+v, want %+v", got, want)
	}
	if _, err := os.Stat(filepath.Join(root, "services", "api", "README.md")); !os.IsNotExist(err) {
		t.Errorf("add wrote a README.md in a folder that was there: %v", err)
	}
	if m := readManifest(t, root); !strings.Contains(m, "  - name: api\n    path: services/api\n    kind: go\n    tags: [backend, http]\n") {
		t.Errorf("manifest:\n%s", m)
	}
}

// S-0491: set changes the keys given and keeps the rest, components among
// them; renames; and creates a folder that a new path names, leaving the old
// one alone.
func TestSubprojectSet(t *testing.T) {
	root, _ := subprojectProject(t)
	sharedRun(t, root, "subproject", "add", "web", "--kind", "sveltekit", "--tags", "ui")
	sharedRun(t, root, "components", "config", "web", "--include", "src/**/*.svelte")

	out := sharedRun(t, root, "subproject", "set", "web", "--tags", "ui,app")
	if out != "set sub-project web (web, sveltekit)\nwrote system-flow.yaml\n" {
		t.Errorf("set says %q", out)
	}
	m := readManifest(t, root)
	if !strings.Contains(m, "    tags: [ui, app]\n    components:\n      include: [\"src/**/*.svelte\"]\n") {
		t.Errorf("set --tags changes tags alone:\n%s", m)
	}

	var got subprojectChange
	out = sharedRun(t, root, "subproject", "set", "web", "--name", "site", "--path", "apps/site", "--json")
	if err := json.Unmarshal([]byte(out), &got); err != nil {
		t.Fatalf("%v: %s", err, out)
	}
	want := subprojectChange{
		Action:  manifest.ProjectSet,
		Project: manifest.Project{Name: "site", Path: "apps/site", Kind: "sveltekit", Tags: []string{"ui", "app"}, Components: &manifest.Components{Include: []string{"src/**/*.svelte"}}},
		From:    "web", Created: "apps/site", Written: []string{manifest.File, "apps/site/README.md"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("set --json =\n%+v\nwant\n%+v", got, want)
	}
	if !isFolder(filepath.Join(root, "apps", "site")) || !isFolder(filepath.Join(root, "web")) {
		t.Error("set --path creates the new folder and leaves the old one")
	}
	m = readManifest(t, root)
	if !strings.Contains(m, "  - name: site\n    path: apps/site\n    kind: sveltekit\n    tags: [ui, app]\n    components:\n") {
		t.Errorf("set renamed and moved:\n%s", m)
	}

	if out := sharedRun(t, root, "subproject", "set", "site", "--tags", ""); !strings.HasPrefix(out, "set sub-project site") {
		t.Errorf("set --tags \"\" says %q", out)
	}
	if m := readManifest(t, root); strings.Contains(m, "tags:") {
		t.Errorf("set --tags \"\" removes the tags:\n%s", m)
	}
}

// S-0491: remove takes the entry out, leaves the folder, and says so;
// removing the last leaves projects: [].
func TestSubprojectRemove(t *testing.T) {
	root, before := subprojectProject(t)
	sharedRun(t, root, "subproject", "add", "web", "--kind", "sveltekit")
	out := sharedRun(t, root, "subproject", "remove", "web")
	if out != "removed sub-project web; its folder web is left as it is\nwrote system-flow.yaml\n" {
		t.Errorf("remove says %q", out)
	}
	if got := readManifest(t, root); got != before {
		t.Errorf("manifest after removing the last:\n%s\nwant:\n%s", got, before)
	}
	if _, err := os.Stat(filepath.Join(root, "web", "README.md")); err != nil {
		t.Errorf("remove left the folder alone: %v", err)
	}
}

// S-0491: an add, set, or remove that is refused writes nothing, neither the
// manifest nor a folder, and says why.
func TestSubprojectRefusals(t *testing.T) {
	root, _ := subprojectProject(t)
	sharedRun(t, root, "subproject", "add", "web", "--kind", "sveltekit")
	if err := os.WriteFile(filepath.Join(root, "notes"), []byte("a file\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	outside := "flai-subproject-outside-" + filepath.Base(root)
	for _, tc := range []struct {
		name        string
		args        []string
		out, errOut string
	}{
		{"duplicate name", []string{"add", "web", "--kind", "go"}, "", "already; give the new one another name"},
		{"path leaving the repository", []string{"add", "out", "--path", "../" + outside, "--kind", "go"}, "projects[1].path", "rule: "},
		{"path that is a file", []string{"add", "notes", "--kind", "go"}, "names no folder", "rule: "},
		{"path under a file", []string{"add", "deep", "--path", "notes/deep", "--kind", "go"}, "names no folder", "rule: "},
		{"empty kind", []string{"add", "api", "--kind", ""}, "kind", "rule: "},
		{"no kind", []string{"add", "api"}, "", "required flag"},
		{"set unknown", []string{"set", "api", "--kind", "go"}, "", "; name one of web"},
		{"set nothing", []string{"set", "web"}, "", "nothing to change"},
		{"set to another's path", []string{"set", "web", "--path", "notes"}, "names no folder", "rule: "},
		{"remove unknown", []string{"remove", "api"}, "", "no sub-project"},
	} {
		before := readManifest(t, root)
		out, errOut, code := runIn(t, root, append([]string{"subproject"}, tc.args...)...)
		if code != 1 || !strings.Contains(out, tc.out) || !strings.Contains(errOut, tc.errOut) {
			t.Errorf("%s: exit %d, out %q, err %q", tc.name, code, out, errOut)
		}
		if m := readManifest(t, root); m != before {
			t.Errorf("%s changed the manifest:\n%s", tc.name, m)
		}
	}
	for _, p := range []string{"api", "out", "deep", filepath.Join("..", outside)} {
		if _, err := os.Stat(filepath.Join(root, p)); !os.IsNotExist(err) {
			t.Errorf("a refused edit created %s: %v", p, err)
		}
	}

	// A refusal found only once the folder is made removes the folder.
	file := filepath.Join(root, manifest.File)
	data, _ := os.ReadFile(file)
	_ = os.WriteFile(file, []byte(strings.Replace(string(data), "projects:\n", "projects:\n  - name: gone\n    path: gone\n    kind: go\n", 1)), 0o644)
	before := readManifest(t, root)
	out, _, code := runIn(t, root, "subproject", "add", "api", "--path", "apps/api", "--kind", "go", "--json")
	var refused struct {
		Refused []manifest.Problem `json:"refused"`
	}
	if code != 1 || json.Unmarshal([]byte(out), &refused) != nil || len(refused.Refused) != 1 || !strings.Contains(refused.Refused[0].Reason, "names no folder") {
		t.Errorf("another entry's missing folder: exit %d, %q", code, out)
	}
	if _, err := os.Stat(filepath.Join(root, "apps")); !os.IsNotExist(err) {
		t.Errorf("the refused add left the folder it made: %v", err)
	}
	if m := readManifest(t, root); m != before {
		t.Errorf("the refused add changed the manifest:\n%s", m)
	}
	// remove is not stopped by another entry's missing folder.
	sharedRun(t, root, "subproject", "remove", "web")
}

// S-0491: with --autocommit, add commits the manifest and the README.md it
// wrote in one commit, with the trailer given, and nothing else.
func TestSubprojectAddAutocommitsTheManifestAndTheReadme(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}
	root, _ := subprojectProject(t)
	gitIn(t, root, "init", "-q", "-b", "main")
	gitIn(t, root, "config", "user.email", "t@t")
	gitIn(t, root, "config", "user.name", "t")
	gitIn(t, root, "add", "-A")
	gitIn(t, root, "commit", "-q", "-m", "start")
	_ = os.WriteFile(filepath.Join(root, "notes.txt"), []byte("not mine\n"), 0o644)

	var got subprojectChange
	out := sharedRun(t, root, "subproject", "add", "web", "--kind", "sveltekit", "--autocommit", "--trailer", "Co-Authored-By: someone <s@s>", "--json")
	if err := json.Unmarshal([]byte(out), &got); err != nil {
		t.Fatalf("%v: %s", err, out)
	}
	if head := gitIn(t, root, "rev-parse", "--short", "HEAD"); got.Commit == "" || got.Commit != head {
		t.Errorf("commit = %q, HEAD %q", got.Commit, head)
	}
	msg := gitIn(t, root, "log", "-1", "--format=%B")
	if !strings.Contains(msg, "chore: add sub-project web") || !strings.Contains(msg, "Co-Authored-By: someone <s@s>") {
		t.Errorf("commit message: %q", msg)
	}
	if files := gitIn(t, root, "show", "--name-only", "--format=", "HEAD"); files != manifest.File+"\nweb/README.md" {
		t.Errorf("the commit holds the manifest and the README.md alone: %q", files)
	}
	if st := gitIn(t, root, "status", "--porcelain"); st != "?? notes.txt" {
		t.Errorf("status after the commit: %q", st)
	}
}
