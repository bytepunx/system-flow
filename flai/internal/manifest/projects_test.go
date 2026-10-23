package manifest

import (
	"errors"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/goccy/go-yaml"
)

// projectsFixture is a manifest with comments and a blank line around its
// sub-projects, one of them with components and install keys.
const projectsFixture = `version: 1
name: demo
layout:
  design: design
  docs: docs
  wip: wip
projects:                 # the releasable components
  - name: api             # the service
    path: api
    kind: go
    tags: [cli]
  # the app
  - name: web
    path: web             # its folder
    kind: sveltekit
    components:
      setup: src/app.css
    install:
      command: [pnpm, install]
      paths: [package.json]

dashboard:
  port: 4242
`

// projectsFiles are the folders the fixture's sub-projects and the ones the
// tests add are in.
var projectsFiles = fstest.MapFS{
	"api/main.go":          {Data: nil},
	"web/src/app.css":      {Data: nil},
	"apps/web/src/app.css": {Data: nil},
	"tools/package.json":   {Data: nil},
	"README.md":            {Data: nil},
}

// S-0491: ProjectErrors reports, by index, field, and name, a name that is
// empty or repeated, a path that is empty, the root, absolute, unclean,
// outside the repository, repeated, or no folder, and an empty kind; a kind
// flai does not know is not one.
func TestProjectErrors(t *testing.T) {
	entry := func(name, path, kind string) Project { return Project{Name: name, Path: path, Kind: kind} }
	if errs := (Manifest{Projects: []Project{entry("api", "api", "go"), entry("tools", "tools", "deno")}}).ProjectErrors(projectsFiles); len(errs) != 0 {
		t.Errorf("valid projects: %v", errs)
	}
	for _, c := range []struct {
		name     string
		projects []Project
		want     string
	}{
		{"an empty name", []Project{entry(" ", "api", "go")}, "projects[0].name is empty"},
		{"a repeated name", []Project{entry("api", "api", "go"), entry("api", "tools", "node")}, `projects[1].name (sub-project "api") repeats the name of projects[0]`},
		{"an empty path", []Project{entry("api", "", "go")}, `projects[0].path (sub-project "api") is empty; write the sub-project's folder, relative to the repository root, such as api`},
		{"the root", []Project{entry("api", ".", "go")}, `projects[0].path (sub-project "api") "." is the repository root`},
		{"an absolute path", []Project{entry("api", "/srv/api", "go")}, `projects[0].path (sub-project "api") "/srv/api" is an absolute path`},
		{"an unclean path", []Project{entry("api", "api/", "go")}, `"api/" is not a clean path; write it with single slashes, no . or .. segments, and none at the end, such as api`},
		{"a path outside", []Project{entry("api", "../api", "go")}, `"../api" leaves the repository`},
		{"a repeated path", []Project{entry("api", "api", "go"), entry("svc", "api", "go")}, `projects[1].path (sub-project "svc") "api" is also the path of projects[0]`},
		{"an empty kind", []Project{entry("api", "api", "")}, `projects[0].kind (sub-project "api") is empty; write what the sub-project is, such as go`},
		{"no folder", []Project{entry("lib", "lib", "go")}, `projects[0].path (sub-project "lib") "lib" names no folder in the repository`},
		{"a file", []Project{entry("readme", "README.md", "go")}, `"README.md" names no folder`},
	} {
		t.Run(c.name, func(t *testing.T) {
			errs := (Manifest{Projects: c.projects}).ProjectErrors(projectsFiles)
			if len(errs) != 1 || !strings.Contains(errs[0], c.want) {
				t.Errorf("got %q, want one error ...%s...", errs, c.want)
			}
		})
	}
	if errs := (Manifest{Projects: []Project{entry("lib", "lib", "go")}}).ProjectErrors(nil); len(errs) != 0 {
		t.Errorf("with no file system, a path is not looked for: %v", errs)
	}
}

// S-0491: EditProjects adds an entry after the last, to projects: [] as the
// template writes it, and where there is no projects key, laid out as the
// first entry is, keeping every other line and comment.
func TestEditProjectsAdds(t *testing.T) {
	tools := Project{Name: "tools", Path: "tools", Kind: "node", Tags: []string{"build", "ci"}}
	for _, c := range []struct {
		name, data string
		p          Project
		want       string
	}{
		{"after the last entry", projectsFixture, tools,
			strings.Replace(projectsFixture, "      paths: [package.json]\n", "      paths: [package.json]\n  - name: tools\n    path: tools\n    kind: node\n    tags: [build, ci]\n", 1)},
		{"to projects: []", "version: 1\nname: demo\nprojects: [] # none yet\ndashboard:\n  port: 4242\n", tools,
			"version: 1\nname: demo\nprojects: # none yet\n  - name: tools\n    path: tools\n    kind: node\n    tags: [build, ci]\ndashboard:\n  port: 4242\n"},
		{"with no projects key", "version: 1\nname: demo\n", Project{Name: "api", Path: "api", Kind: "go"},
			"version: 1\nname: demo\nprojects:\n  - name: api\n    path: api\n    kind: go\n"},
		{"at the first entry's indent", "version: 1\nprojects:\n-   name: api\n    path: api\n    kind: go\n", Project{Name: "tools", Path: "tools", Kind: "node", Install: &Install{Command: []string{"pnpm", "install"}, Paths: []string{"package.json"}}},
			"version: 1\nprojects:\n-   name: api\n    path: api\n    kind: go\n-   name: tools\n    path: tools\n    kind: node\n    install:\n      command: [pnpm, install]\n      paths: [package.json]\n"},
	} {
		t.Run(c.name, func(t *testing.T) {
			out, err := EditProjects([]byte(c.data), projectsFiles, ProjectOp{Action: ProjectAdd}, c.p)
			if err != nil {
				t.Fatal(err)
			}
			if string(out) != c.want {
				t.Errorf("got:\n%s\nwant:\n%s", out, c.want)
			}
			got := projectsOf(t, out)
			if last := got[len(got)-1]; !sameProject(last, c.p) {
				t.Errorf("reads back as %+v, want %+v", last, c.p)
			}
		})
	}
}

// S-0491: EditProjects sets an entry's keys where they are, keeping their
// comments and the components and install it leaves out, adds a key the
// entry lacks after the one before it, removes tags it gives none of, and
// renames an entry.
func TestEditProjectsSets(t *testing.T) {
	for _, c := range []struct {
		name string
		op   ProjectOp
		p    Project
		want string
	}{
		{"its path, a tag added", ProjectOp{Action: ProjectSet}, Project{Name: "web", Path: "apps/web", Kind: "sveltekit", Tags: []string{"dashboard"}},
			strings.Replace(projectsFixture, "    path: web             # its folder\n    kind: sveltekit\n", "    path: apps/web        # its folder\n    kind: sveltekit\n    tags: [dashboard]\n", 1)},
		{"renamed, its tags removed", ProjectOp{Action: ProjectSet, Name: "api"}, Project{Name: "server", Path: "api", Kind: "go"},
			strings.Replace(projectsFixture, "  - name: api             # the service\n    path: api\n    kind: go\n    tags: [cli]\n", "  - name: server          # the service\n    path: api\n    kind: go\n", 1)},
		{"its components and install replaced", ProjectOp{Action: ProjectSet}, Project{Name: "web", Path: "web", Kind: "sveltekit",
			Components: &Components{Exclude: []string{"src/lib/x.svelte"}}, Install: &Install{Command: []string{"npm", "ci"}, Paths: []string{"package-lock.json"}}},
			strings.Replace(projectsFixture, "      setup: src/app.css\n    install:\n      command: [pnpm, install]\n      paths: [package.json]\n", "      exclude: [src/lib/x.svelte]\n    install:\n      command: [npm, ci]\n      paths: [package-lock.json]\n", 1)},
	} {
		t.Run(c.name, func(t *testing.T) {
			out, err := EditProjects([]byte(projectsFixture), projectsFiles, c.op, c.p)
			if err != nil {
				t.Fatal(err)
			}
			if string(out) != c.want {
				t.Errorf("got:\n%s\nwant:\n%s", out, c.want)
			}
		})
	}
	out, err := EditProjects([]byte(projectsFixture), projectsFiles, ProjectOp{Action: ProjectSet}, Project{Name: "web", Path: "web", Kind: "svelte"})
	if err != nil {
		t.Fatal(err)
	}
	if web := loaded(t, out).Projects[1]; web.Kind != "svelte" || web.Components == nil || web.Components.Setup != "src/app.css" || web.Install == nil {
		t.Errorf("a set that leaves components and install out keeps them: %+v", web)
	}
}

// S-0491: EditProjects removes an entry with its lines, keeping the comments
// before the next one, and leaves projects: [] when it removes the last.
func TestEditProjectsRemoves(t *testing.T) {
	for _, c := range []struct {
		name, data, project, want string
	}{
		{"the first", projectsFixture, "api", strings.Replace(projectsFixture, "  - name: api             # the service\n    path: api\n    kind: go\n    tags: [cli]\n", "", 1)},
		{"the last", projectsFixture, "web", strings.Replace(projectsFixture, "  - name: web\n    path: web             # its folder\n    kind: sveltekit\n    components:\n      setup: src/app.css\n    install:\n      command: [pnpm, install]\n      paths: [package.json]\n", "", 1)},
		{"the only one", "version: 1\nprojects: # the components\n  - name: api\n    path: api\n    kind: go\ndashboard:\n  port: 4242\n", "api", "version: 1\nprojects: [] # the components\ndashboard:\n  port: 4242\n"},
	} {
		t.Run(c.name, func(t *testing.T) {
			out, err := EditProjects([]byte(c.data), projectsFiles, ProjectOp{Action: ProjectRemove, Name: c.project}, Project{})
			if err != nil {
				t.Fatal(err)
			}
			if string(out) != c.want {
				t.Errorf("got:\n%s\nwant:\n%s", out, c.want)
			}
		})
	}
}

// S-0491: EditProjects refuses, with the field and the reason, what
// ProjectErrors, ComponentErrors, and InstallErrors find in the manifest as
// it would be; it refuses a name added twice or set or removed that is not
// there, an entry it cannot edit, and one that would not read back; and
// SetProjects then writes nothing.
func TestEditProjectsRefuses(t *testing.T) {
	for _, c := range []struct {
		name          string
		op            ProjectOp
		p             Project
		field, reason string
	}{
		{"a repeated path", ProjectOp{Action: ProjectAdd}, Project{Name: "svc", Path: "api", Kind: "go"}, "projects[2].path", `(sub-project "svc") "api" is also the path of projects[0]`},
		{"an empty kind", ProjectOp{Action: ProjectAdd}, Project{Name: "tools", Path: "tools"}, "projects[2].kind", "is empty"},
		{"no folder", ProjectOp{Action: ProjectAdd}, Project{Name: "lib", Path: "lib", Kind: "go"}, "projects[2].path", "names no folder"},
		{"a rename to a name taken", ProjectOp{Action: ProjectSet, Name: "web"}, Project{Name: "api", Path: "web", Kind: "sveltekit"}, "projects[1].name", "repeats the name of projects[0]"},
		{"components on a go sub-project", ProjectOp{Action: ProjectSet}, Project{Name: "api", Path: "api", Kind: "go", Components: &Components{Setup: "main.go"}}, "projects[0].components", `is set on a sub-project of kind "go"`},
		{"an install with no paths", ProjectOp{Action: ProjectSet}, Project{Name: "api", Path: "api", Kind: "go", Install: &Install{Command: []string{"go", "mod", "download"}}}, "projects[0].install.paths", "is empty"},
	} {
		t.Run(c.name, func(t *testing.T) {
			_, err := EditProjects([]byte(projectsFixture), projectsFiles, c.op, c.p)
			var refused *RefusedError
			if !errors.As(err, &refused) {
				t.Fatalf("want a refusal, got %v", err)
			}
			i := slices.IndexFunc(refused.Problems, func(p Problem) bool { return p.Field == c.field })
			if i < 0 || !strings.Contains(refused.Problems[i].Reason, c.reason) {
				t.Errorf("want %s: ...%s..., got %+v", c.field, c.reason, refused.Problems)
			}
		})
	}
	api := Project{Name: "api", Path: "api", Kind: "go"}
	for _, c := range []struct {
		name, data string
		op         ProjectOp
		p          Project
		want       string
	}{
		{"a name added twice", projectsFixture, ProjectOp{Action: ProjectAdd}, api, `system-flow.yaml has a sub-project "api" already`},
		{"a set of a name not there", projectsFixture, ProjectOp{Action: ProjectSet, Name: "nope"}, api, `system-flow.yaml has no sub-project "nope"; name one of api, web`},
		{"a remove of a name not there", projectsFixture, ProjectOp{Action: ProjectRemove, Name: "nope"}, Project{}, `system-flow.yaml has no sub-project "nope"`},
		{"a remove from none", "version: 1\nprojects: []\n", ProjectOp{Action: ProjectRemove}, api, "names none under projects"},
		{"no name", projectsFixture, ProjectOp{Action: ProjectRemove}, Project{}, "give the name of the sub-project to remove"},
		{"an unknown action", projectsFixture, ProjectOp{Action: "move"}, api, `cannot edit projects with "move"`},
		{"projects on one line", "version: 1\nprojects: [{name: web, path: web, kind: sveltekit}]\n", ProjectOp{Action: ProjectAdd}, api, "projects is written on one line"},
		{"an entry on one line", "version: 1\nprojects:\n  - {name: api, path: api, kind: go}\n", ProjectOp{Action: ProjectSet}, api, "its entry under projects is written on one line"},
		{"an entry that begins with tags", "version: 1\nprojects:\n  - tags: [cli]\n    name: api\n    path: api\n    kind: go\n", ProjectOp{Action: ProjectSet}, api, "its entry under projects begins with tags"},
		{"a key it cannot find", "version: 1\nprojects:\n  - \"name\": api\n    path: api\n    kind: go\n", ProjectOp{Action: ProjectSet}, Project{Name: "api", Path: "api", Kind: "node"}, "so that system-flow.yaml reads it back as given"},
	} {
		t.Run(c.name, func(t *testing.T) {
			if _, err := EditProjects([]byte(c.data), projectsFiles, c.op, c.p); err == nil || !strings.Contains(err.Error(), c.want) {
				t.Errorf("want ...%s..., got %v", c.want, err)
			}
		})
	}
	root := t.TempDir()
	file := filepath.Join(root, File)
	if err := os.WriteFile(file, []byte(projectsFixture), 0o644); err != nil {
		t.Fatal(err)
	}
	tools := Project{Name: "tools", Path: "tools", Kind: "node"}
	if err := SetProjects(file, os.DirFS(root), ProjectOp{Action: ProjectAdd}, tools); err == nil {
		t.Error("SetProjects added a sub-project whose folder is not there")
	}
	if data, _ := os.ReadFile(file); string(data) != projectsFixture {
		t.Errorf("the file changed on a refusal:\n%s", data)
	}
	for _, dir := range []string{"api", "web/src", "tools"} {
		if err := os.MkdirAll(filepath.Join(root, filepath.FromSlash(dir)), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.WriteFile(filepath.Join(root, "web", "src", "app.css"), nil, 0o644); err != nil {
		t.Fatal(err)
	}
	if err := SetProjects(file, os.DirFS(root), ProjectOp{Action: ProjectAdd}, tools); err != nil {
		t.Fatal(err)
	}
	if m := loaded(t, mustRead(t, file)); len(m.Projects) != 3 || !sameProject(m.Projects[2], tools) {
		t.Errorf("SetProjects wrote %+v", m.Projects)
	}
}

// S-0491: flai manifest set refuses a manifest whose projects ProjectErrors
// finds wrong, as flai check reports it.
func TestEditSettingsRefusesBadProjects(t *testing.T) {
	data := "version: 1\nname: demo\nprojects:\n  - name: web\n    path: web\n    kind: go\n  - name: web\n    path: web2\n    kind: go\n"
	_, err := EditSettings([]byte(data), []Assignment{{Key: "planning.replan", Value: "agent"}}, nil)
	var refused *RefusedError
	if !errors.As(err, &refused) {
		t.Fatalf("want a refusal, got %v", err)
	}
	if i := slices.IndexFunc(refused.Problems, func(p Problem) bool { return p.Field == "projects[1].name" }); i < 0 || !strings.Contains(refused.Problems[i].Reason, "repeats the name of projects[0]") {
		t.Errorf("want projects[1].name repeated, got %+v", refused.Problems)
	}
}

// projectsOf are the sub-projects of the manifest data, read without the
// validation Load does.
func projectsOf(t *testing.T, data []byte) []Project {
	t.Helper()
	var m Manifest
	if err := yaml.Unmarshal(data, &m); err != nil {
		t.Fatal(err)
	}
	return m.Projects
}

func mustRead(t *testing.T, file string) []byte {
	t.Helper()
	data, err := os.ReadFile(file)
	if err != nil {
		t.Fatal(err)
	}
	return data
}
