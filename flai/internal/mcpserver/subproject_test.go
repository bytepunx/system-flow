package mcpserver

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/gittest"
	"github.com/bytepunx/system-flow/flai/internal/guard"
)

// subprojectManifest is the fixture's manifest with two sub-projects: cli,
// whose folder is there, with tags and an install key, and web, whose folder
// is missing, with a components key and no tags.
const subprojectManifest = "version: 1\nname: t\nkey: t\nlayout:\n  design: design\n  docs: docs\n  wip: wip\nprojects:\n  # the sub-projects\n  - name: cli\n    path: flai\n    kind: go\n    tags: [backend]\n    install:\n      command: [go, mod, download]\n      paths: [go.sum]\n  - name: web\n    path: flaiover\n    kind: sveltekit\n    components:\n      setup: src/app.css\n"

// subprojectFixture is setup in a session flai serve started when served is
// guard.StartedByServe, with subprojectManifest written after the server
// opened the project and cli's folder made, so that the tools must read the
// manifest as it is now.
func subprojectFixture(t *testing.T, served string) *fixture {
	t.Helper()
	for _, v := range []string{"FLAI_ROLE", "FLAI_STORY"} {
		t.Setenv(v, "")
	}
	t.Setenv(guard.StartedByEnv, served)
	f := setup(t)
	if err := os.WriteFile(f.manifest(), []byte(subprojectManifest), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(f.repo.Root, "flai"), 0o755); err != nil {
		t.Fatal(err)
	}
	return f
}

// subprojectsOut decodes a subprojects answer.
func subprojectsOut(t *testing.T, out map[string]any) SubprojectsOut {
	t.Helper()
	var got SubprojectsOut
	data, _ := json.Marshal(out)
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatal(err)
	}
	return got
}

// listed is a subprojects answer, one "name path kind tags folder" line per
// sub-project.
func listed(ps []Subproject) string {
	var lines []string
	for _, p := range ps {
		lines = append(lines, strings.Join([]string{p.Name, p.Path, p.Kind, "[" + strings.Join(p.Tags, ",") + "]", map[bool]string{true: "folder", false: "missing"}[p.Folder]}, " "))
	}
	return strings.Join(lines, "\n")
}

// S-0491: subprojects lists the manifest's sub-projects as flai subproject
// list --json prints them, in a session flai serve started too: tags [] for
// none, folder for whether the path is a folder, and components and install
// as the manifest has them, null when it leaves them out.
func TestSubprojectsListsTheManifestsSubprojects(t *testing.T) {
	f := subprojectFixture(t, guard.StartedByServe)
	out, failed := f.call(t, guard.Subprojects, map[string]any{})
	if failed != "" {
		t.Fatal(failed)
	}
	got := subprojectsOut(t, out)
	if want := "cli flai go [backend] folder\nweb flaiover sveltekit [] missing"; listed(got.Projects) != want || got.Story != "" {
		t.Errorf("subprojects:\n%s\nwant:\n%s\n(story %q)", listed(got.Projects), want, got.Story)
	}
	cli, web := got.Projects[0], got.Projects[1]
	if cli.Install == nil || strings.Join(cli.Install.Command, " ") != "go mod download" || cli.Components != nil {
		t.Errorf("cli's keys: %+v", cli)
	}
	if web.Components == nil || web.Components.Setup != "src/app.css" || web.Install != nil {
		t.Errorf("web's keys: %+v", web)
	}
	entries := out["projects"].([]any)
	if e := entries[1].(map[string]any); e["install"] != nil || !reflect.DeepEqual(e["tags"], []any{}) {
		t.Errorf("web as JSON: install null and tags [], got %v", e)
	}
	if _, ok := entries[1].(map[string]any)["install"]; !ok {
		t.Error("install should be given as null, not left out, as flai subproject list --json gives it")
	}
}

// webFolder makes web's folder and its components' setup module under root,
// so that the manifest checks as flai check checks it, which an add or a
// set needs of every entry.
func webFolder(t *testing.T, root string) {
	t.Helper()
	for _, d := range []string{"flai", "flaiover/src"} {
		if err := os.MkdirAll(filepath.Join(root, d), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.WriteFile(filepath.Join(root, "flaiover/src/app.css"), nil, 0o644); err != nil {
		t.Fatal(err)
	}
}

// builtFlai builds flai from this tree for subproject_edit to run, in place
// of the test binary os.Executable names, and gives the run a flai
// configuration of its own.
func builtFlai(t *testing.T) {
	t.Helper()
	if _, err := exec.LookPath("go"); err != nil {
		t.Skip("go not installed, to build flai with")
	}
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}
	exe := filepath.Join(t.TempDir(), "flai")
	if out, err := exec.Command("go", "build", "-o", exe, "github.com/bytepunx/system-flow/flai").CombinedOutput(); err != nil {
		t.Fatalf("go build: %v\n%s", err, out)
	}
	was := flaiExecutable
	flaiExecutable = func() (string, error) { return exe, nil }
	t.Cleanup(func() { flaiExecutable = was })
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "config.json"))
	gittest.Identity(t, "t", "t@t")
}

// S-0491: subproject_edit runs flai subproject in the project's checkout,
// and is not refused in a session flai serve started, unlike
// shared_paths_edit: add appends the entry and creates its missing folder
// with a README.md, set renames it and clears its tags while keeping its
// other keys, and remove takes it out and leaves the folder, each committing
// nothing; a refusal writes nothing, neither the manifest nor a folder.
func TestSubprojectEditAddsSetsAndRemovesCommittingNothing(t *testing.T) {
	builtFlai(t)
	f := subprojectFixture(t, guard.StartedByServe)
	root := f.repo.Root
	webFolder(t, root)
	gitIn(t, root, "init", "-q", "-b", "main")
	gitIn(t, root, "add", "-A")
	gitIn(t, root, "commit", "-q", "-m", "project")
	head := gitIn(t, root, "rev-parse", "HEAD")

	out, failed := f.call(t, guard.SubprojectEdit, map[string]any{"op": "add", "name": "api", "kind": "go", "path": "services/api", "tags": []string{"backend", "http"}})
	if failed != "" {
		t.Fatal(failed)
	}
	if out["action"] != "add" || out["created"] != "services/api" || !reflect.DeepEqual(toStrings(out["written"]), []string{"system-flow.yaml", "services/api/README.md"}) || out["story"] != nil {
		t.Errorf("add: %v", out)
	}
	if p := out["project"].(map[string]any); p["name"] != "api" || p["path"] != "services/api" || p["kind"] != "go" || strings.Join(toStrings(p["tags"]), ",") != "backend,http" {
		t.Errorf("add's entry: %v", p)
	}
	if readme, err := os.ReadFile(filepath.Join(root, "services/api/README.md")); err != nil || !strings.Contains(string(readme), "# api") {
		t.Errorf("the folder's README.md: %q, %v", readme, err)
	}
	if h := gitIn(t, root, "rev-parse", "HEAD"); h != head {
		t.Errorf("add committed %s", h)
	}
	if st := gitIn(t, root, "status", "--porcelain"); !strings.Contains(st, "M system-flow.yaml") || !strings.Contains(st, "?? services/") {
		t.Errorf("add should leave its files uncommitted:\n%s", st)
	}
	listedOut, _ := f.call(t, guard.Subprojects, map[string]any{})
	if l := listed(subprojectsOut(t, listedOut).Projects); !strings.HasSuffix(l, "\napi services/api go [backend,http] folder") {
		t.Errorf("subprojects after the add:\n%s", l)
	}

	out, failed = f.call(t, guard.SubprojectEdit, map[string]any{"op": "set", "name": "cli", "new_name": "tool", "tags": []string{}})
	if failed != "" {
		t.Fatal(failed)
	}
	p := out["project"].(map[string]any)
	if out["action"] != "set" || out["from"] != "cli" || p["name"] != "tool" || p["path"] != "flai" || p["tags"] != nil || p["install"] == nil || out["created"] != nil {
		t.Errorf("set: %v", out)
	}

	out, failed = f.call(t, guard.SubprojectEdit, map[string]any{"op": "remove", "name": "web"})
	if failed != "" {
		t.Fatal(failed)
	}
	if out["action"] != "remove" || out["project"].(map[string]any)["path"] != "flaiover" || !reflect.DeepEqual(toStrings(out["written"]), []string{"system-flow.yaml"}) {
		t.Errorf("remove: %v", out)
	}
	data, _ := os.ReadFile(f.manifest())
	for _, want := range []string{"# the sub-projects", "  - name: tool\n    path: flai\n    kind: go\n    install:", "  - name: api\n    path: services/api\n    kind: go\n"} {
		if !strings.Contains(string(data), want) {
			t.Errorf("the manifest should hold %q:\n%s", want, data)
		}
	}
	if strings.Contains(string(data), "name: web") || strings.Contains(string(data), "backend]") {
		t.Errorf("the manifest after remove and set:\n%s", data)
	}
	if h := gitIn(t, root, "rev-parse", "HEAD"); h != head {
		t.Errorf("set or remove committed %s", h)
	}

	before := string(data)
	for _, c := range []struct {
		args map[string]any
		says string
	}{
		{map[string]any{"op": "add", "name": "tool", "kind": "go", "path": "other"}, "name"},
		{map[string]any{"op": "add", "name": "abs", "kind": "go", "path": "/abs"}, "path"},
		{map[string]any{"op": "add", "name": "nokind", "path": "nokind"}, "kind"},
		{map[string]any{"op": "set", "name": "nosuch", "kind": "go"}, "nosuch"},
		{map[string]any{"op": "set", "name": "tool"}, "give new_name, path, kind, or tags"},
		{map[string]any{"op": "remove", "name": "tool", "path": "flai"}, "remove takes name alone"},
		{map[string]any{"op": "add", "name": "x", "kind": "go", "tags": []string{"a,b"}}, "holds a comma"},
		{map[string]any{"op": "move", "name": "tool"}, "give add, set, or remove"},
		{map[string]any{"op": "add", "name": " ", "kind": "go"}, "give name"},
	} {
		_, failed := f.call(t, guard.SubprojectEdit, c.args)
		if !strings.Contains(failed, c.says) {
			t.Errorf("%v: %q, want it to say %q", c.args, failed, c.says)
		}
		if data, _ := os.ReadFile(f.manifest()); string(data) != before {
			t.Errorf("%v changed the manifest:\n%s", c.args, data)
		}
	}
	for _, folder := range []string{"other", "nokind", "x"} {
		if _, err := os.Stat(filepath.Join(root, folder)); err == nil {
			t.Errorf("a refusal created the folder %s", folder)
		}
	}
}

// S-0491: for a session whose story has a worktree, both tools work in that
// worktree, where the story's branch carries the change for flai task done
// to commit, and say the story; the main checkout's manifest is left alone.
func TestSubprojectEditWorksInTheStorysWorktree(t *testing.T) {
	builtFlai(t)
	f := subprojectFixture(t, guard.StartedByServe)
	t.Setenv("FLAI_STORY", f.story.ID)
	wt := f.repo.WorktreePath(f.story.ID)
	webFolder(t, wt)
	branch := strings.Replace(subprojectManifest, "  - name: web\n", "  - name: site\n", 1)
	if err := os.WriteFile(filepath.Join(wt, "system-flow.yaml"), []byte(branch), 0o644); err != nil {
		t.Fatal(err)
	}
	out, failed := f.call(t, guard.SubprojectEdit, map[string]any{"op": "add", "name": "docs-site", "kind": "node"})
	if failed != "" {
		t.Fatal(failed)
	}
	rel, _ := filepath.Rel(f.repo.Root, wt)
	at := filepath.ToSlash(rel) + "/"
	if out["story"] != f.story.ID || out["created"] != "docs-site" || !reflect.DeepEqual(toStrings(out["written"]), []string{at + "system-flow.yaml", at + "docs-site/README.md"}) {
		t.Errorf("add in the worktree: %v", out)
	}
	if _, err := os.Stat(filepath.Join(wt, "docs-site", "README.md")); err != nil {
		t.Errorf("the folder should be made in the worktree: %v", err)
	}
	if _, err := os.Stat(filepath.Join(f.repo.Root, "docs-site")); err == nil {
		t.Error("the main checkout should be left alone while the story has a worktree")
	}
	if data, _ := os.ReadFile(f.manifest()); string(data) != subprojectManifest {
		t.Errorf("the main checkout's manifest changed:\n%s", data)
	}
	listedOut, failed := f.call(t, guard.Subprojects, map[string]any{})
	if failed != "" {
		t.Fatal(failed)
	}
	got := subprojectsOut(t, listedOut)
	if l := listed(got.Projects); got.Story != f.story.ID || l != "cli flai go [backend] folder\nsite flaiover sveltekit [] folder\ndocs-site docs-site node [] folder" {
		t.Errorf("subprojects in the worktree (story %q):\n%s", got.Story, l)
	}
}
