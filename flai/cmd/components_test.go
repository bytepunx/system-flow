package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/components"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
)

// fakeNode answers for node as props.mjs would: each component path it is
// given gets the props in props, keyed by path, or the error in errs; code,
// when set, is the error the whole answer carries instead.
type fakeNode struct {
	noNode bool
	code   string
	props  map[string][]components.Prop
	errs   map[string]string
	dirs   []string
}

func (f *fakeNode) LookPath(name string) (string, error) {
	if f.noNode {
		return "", errors.New("executable file not found in $PATH")
	}
	return "/usr/bin/" + name, nil
}

func (f *fakeNode) Run(_ context.Context, dir, _, input string, _ ...string) ([]byte, []byte, error) {
	f.dirs = append(f.dirs, dir)
	if f.code != "" {
		out, err := json.Marshal(map[string]any{"error": map[string]string{"code": f.code, "message": "cannot resolve svelte/compiler"}})
		return out, nil, err
	}
	var paths []string
	if err := json.Unmarshal([]byte(input), &paths); err != nil {
		return nil, nil, err
	}
	type entry struct {
		Path  string            `json:"path"`
		Props []components.Prop `json:"props"`
		Error string            `json:"error,omitempty"`
	}
	answer := struct {
		Components []entry `json:"components"`
	}{}
	for _, p := range paths {
		answer.Components = append(answer.Components, entry{Path: p, Props: f.props[p], Error: f.errs[p]})
	}
	out, err := json.Marshal(answer)
	return out, nil, err
}

// useNode makes flai components list run f in place of node for the test.
func useNode(t *testing.T, f *fakeNode) {
	t.Helper()
	old := componentsRunner
	t.Cleanup(func() { componentsRunner = old })
	componentsRunner = f
}

// componentsManifest is the projects block componentsFixture adds: a
// SvelteKit app with the default components key, a Go sub-project, and a
// Svelte package with a key of its own.
const componentsManifest = `# the sub-projects
projects:
  - name: ui   # the dashboard
    path: ui
    kind: sveltekit
  - name: cli
    path: cli
    kind: go
  - name: kit
    path: packages/kit
    kind: svelte
    components:
      include: [src/**/*.svelte]
      exclude: [src/internal/**]
      setup: src/app.css
# the end
`

// componentsFixture is a scratch project with componentsManifest's
// sub-projects and their files, and its manifest's contents.
func componentsFixture(t *testing.T) (string, string) {
	t.Helper()
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	root := tempProject(t)
	file := filepath.Join(root, manifest.File)
	data, err := os.ReadFile(file)
	if err != nil {
		t.Fatal(err)
	}
	data = append(data, componentsManifest...)
	if err := os.WriteFile(file, data, 0o644); err != nil {
		t.Fatal(err)
	}
	for _, f := range []string{
		"ui/src/lib/Button.svelte", "ui/src/lib/forms/Field.svelte", "ui/src/routes/+page.svelte", "ui/src/app.css",
		"cli/main.go",
		"packages/kit/src/Card.svelte", "packages/kit/src/internal/Hidden.svelte", "packages/kit/src/app.css",
	} {
		p := filepath.Join(root, filepath.FromSlash(f))
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte("<p>x</p>\n"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return root, string(data)
}

// componentsNode is a fake node with props for the fixture's components:
// Field.svelte does not parse.
func componentsNode() *fakeNode {
	return &fakeNode{
		props: map[string][]components.Prop{
			"src/lib/Button.svelte": {
				{Name: "label", Kind: components.KindString, Type: "string", Required: true},
				{Name: "size", Kind: components.KindUnion, Type: "'sm' | 'lg'", Default: "'sm'", Literals: []string{`"sm"`, `"lg"`}},
				{Name: "onclick", Kind: components.KindFunction, Type: "() => void"},
			},
			"src/Card.svelte": {
				{Name: "items", Kind: components.KindArray, Type: "{\n  id: string;\n}[]", Default: "[]"},
			},
		},
		errs: map[string]string{"src/lib/forms/Field.svelte": "Unexpected token (3:4)"},
	}
}

// S-0437: flai components list prints each sub-project in manifest order,
// a Svelte one's components with their props aligned, a parse error in place
// of a component's props, and any other kind as unsupported.
func TestComponentsListText(t *testing.T) {
	root, _ := componentsFixture(t)
	node := componentsNode()
	useNode(t, node)
	out, errOut, code := runIn(t, root, "components", "list")
	if code != 0 {
		t.Fatalf("exit %d: %s %s", code, out, errOut)
	}
	want := `ui (ui, sveltekit): 2 components
  src/lib/Button.svelte
    label    string           string       required
    size     union "sm"|"lg"  'sm' | 'lg'  = 'sm'
    onclick  function         () => void
  src/lib/forms/Field.svelte
    error: Unexpected token (3:4)

cli (cli): unsupported kind go

kit (packages/kit, svelte, setup src/app.css): 1 component
  src/Card.svelte
    items  array  { id: string; }[]  = []
`
	if out != want {
		t.Errorf("output:\n%s\nwant:\n%s", out, want)
	}
	if len(node.dirs) != 2 || node.dirs[0] != filepath.Join(root, "ui") || node.dirs[1] != filepath.Join(root, "packages", "kit") {
		t.Errorf("node runs in each Svelte sub-project's folder: %v", node.dirs)
	}
}

// S-0437: --json answers {"projects": [...]}, each with its kind, whether it
// is supported, its setup, and its components with their props.
func TestComponentsListJSON(t *testing.T) {
	root, _ := componentsFixture(t)
	useNode(t, componentsNode())
	out, errOut, code := runIn(t, root, "components", "list", "--json")
	if code != 0 {
		t.Fatalf("exit %d: %s %s", code, out, errOut)
	}
	var got struct {
		Projects []struct {
			Name       string                 `json:"name"`
			Path       string                 `json:"path"`
			Kind       string                 `json:"kind"`
			Supported  bool                   `json:"supported"`
			Setup      *string                `json:"setup"`
			Components []components.Component `json:"components"`
		} `json:"projects"`
	}
	if err := json.Unmarshal([]byte(out), &got); err != nil {
		t.Fatalf("not JSON: %v\n%s", err, out)
	}
	if len(got.Projects) != 3 {
		t.Fatalf("projects: %+v", got.Projects)
	}
	ui, cli, kit := got.Projects[0], got.Projects[1], got.Projects[2]
	if ui.Name != "ui" || !ui.Supported || ui.Setup != nil || len(ui.Components) != 2 {
		t.Errorf("ui: %+v", ui)
	}
	button := ui.Components[0]
	if button.Project != "ui" || button.Path != "src/lib/Button.svelte" || button.Name != "Button" || len(button.Props) != 3 {
		t.Errorf("Button: %+v", button)
	} else if size := button.Props[1]; size.Kind != components.KindUnion || size.Default != "'sm'" || size.Required || len(size.Literals) != 2 {
		t.Errorf("size: %+v", size)
	}
	if field := ui.Components[1]; field.Error != "Unexpected token (3:4)" || field.Props == nil || len(field.Props) != 0 {
		t.Errorf("Field: %+v", field)
	}
	if cli.Kind != "go" || cli.Supported || cli.Components == nil || len(cli.Components) != 0 {
		t.Errorf("cli: %+v", cli)
	}
	if kit.Path != "packages/kit" || kit.Setup == nil || *kit.Setup != "src/app.css" || len(kit.Components) != 1 || kit.Components[0].Project != "packages/kit" {
		t.Errorf("kit: %+v", kit)
	}
	if !strings.Contains(out, `"components": [`) || strings.Contains(out, `"error": ""`) {
		t.Errorf("an empty list is [], an empty error left out:\n%s", out)
	}
}

// S-0437: --project lists the one sub-project it names, and an unknown name
// is an error that lists the names there are.
func TestComponentsListProject(t *testing.T) {
	root, _ := componentsFixture(t)
	node := componentsNode()
	useNode(t, node)
	out, errOut, code := runIn(t, root, "components", "list", "--project", "kit")
	if code != 0 || !strings.HasPrefix(out, "kit (packages/kit, svelte, setup src/app.css): 1 component\n") || strings.Contains(out, "ui (") || strings.Contains(out, "cli (") {
		t.Errorf("--project kit: %d %q %s", code, out, errOut)
	}
	out, errOut, code = runIn(t, root, "components", "list", "--project", "cli", "--json")
	if code != 0 || !strings.Contains(out, `"supported": false`) || strings.Contains(out, `"name": "ui"`) {
		t.Errorf("--project cli: %d %q %s", code, out, errOut)
	}
	if len(node.dirs) != 1 {
		t.Errorf("node ran only for kit: %v", node.dirs)
	}
	out, errOut, code = runIn(t, root, "components", "list", "--project", "web")
	if code != 1 || !strings.Contains(errOut, `system-flow.yaml has no sub-project \"web\"; name one of ui, cli, kit`) {
		t.Errorf("unknown --project: %d %q %s", code, out, errOut)
	}
}

// S-0437: without node, or without svelte installed in a Svelte
// sub-project, list fails with what to install.
func TestComponentsListMissingTools(t *testing.T) {
	root, _ := componentsFixture(t)
	useNode(t, &fakeNode{noNode: true})
	out, errOut, code := runIn(t, root, "components", "list")
	if code != 1 || out != "" || !strings.Contains(errOut, "node is not on PATH") || !strings.Contains(errOut, "install Node.js") {
		t.Errorf("no node: %d %q %s", code, out, errOut)
	}
	useNode(t, &fakeNode{code: "svelte-missing"})
	out, errOut, code = runIn(t, root, "components", "list", "--project", "kit")
	if code != 1 || out != "" || !strings.Contains(errOut, "svelte is not installed in packages/kit") || !strings.Contains(errOut, "run npm install") {
		t.Errorf("no svelte: %d %q %s", code, out, errOut)
	}
	// A Go sub-project alone needs neither.
	out, errOut, code = runIn(t, root, "components", "list", "--project", "cli")
	if code != 0 || out != "cli (cli): unsupported kind go\n" {
		t.Errorf("--project cli: %d %q %s", code, out, errOut)
	}
}

// S-0437: config writes a sub-project's key in place, keeping the comments
// and the order; a later flag replaces its field alone; --unset removes the
// key.
func TestComponentsConfigWritesMergesUnsets(t *testing.T) {
	root, before := componentsFixture(t)
	out, errOut, code := runIn(t, root, "components", "config", "ui", "--include", "src/lib/**/*.svelte", "--include", "src/widgets/*.svelte", "--setup", "src/app.css")
	if code != 0 {
		t.Fatalf("exit %d: %s %s", code, out, errOut)
	}
	written := strings.Replace(before, "    kind: sveltekit\n",
		"    kind: sveltekit\n    components:\n      include: [\"src/lib/**/*.svelte\", \"src/widgets/*.svelte\"]\n      setup: src/app.css\n", 1)
	if got := readManifest(t, root); got != written {
		t.Errorf("manifest:\n%s\nwant:\n%s", got, written)
	}
	if want := "set sub-project ui's components\n  include  src/lib/**/*.svelte, src/widgets/*.svelte\n  exclude  none\n  setup    src/app.css\n"; out != want {
		t.Errorf("output %q, want %q", out, want)
	}

	out, errOut, code = runIn(t, root, "components", "config", "ui", "--exclude", "src/lib/internal/**", "--json")
	var got struct {
		Project    string               `json:"project"`
		Components *manifest.Components `json:"components"`
	}
	if code != 0 || json.Unmarshal([]byte(out), &got) != nil {
		t.Fatalf("--exclude --json: %d %q %s", code, out, errOut)
	}
	if got.Project != "ui" || got.Components == nil || len(got.Components.Include) != 2 || got.Components.Setup != "src/app.css" || len(got.Components.Exclude) != 1 {
		t.Errorf("--exclude keeps include and setup: %+v", got.Components)
	}
	merged := strings.Replace(written, "      setup: src/app.css\n", "      exclude: [\"src/lib/internal/**\"]\n      setup: src/app.css\n", 1)
	if m := readManifest(t, root); m != merged {
		t.Errorf("manifest:\n%s\nwant:\n%s", m, merged)
	}

	out, errOut, code = runIn(t, root, "components", "config", "kit", "--setup", "")
	if code != 0 || !strings.Contains(out, "  setup    none\n") {
		t.Errorf(`--setup "": %d %q %s`, code, out, errOut)
	}
	if m := readManifest(t, root); strings.Contains(m, "setup: src/app.css\n# the end") || !strings.Contains(m, "      exclude: [\"src/internal/**\"]\n# the end") {
		t.Errorf(`--setup "" removes kit's setup alone:\n%s`, m)
	}

	out, errOut, code = runIn(t, root, "components", "config", "ui", "--unset")
	if code != 0 || !strings.HasPrefix(out, "unset sub-project ui's components") {
		t.Errorf("--unset: %d %q %s", code, out, errOut)
	}
	out, errOut, code = runIn(t, root, "components", "config", "kit", "--unset", "--json")
	if code != 0 || strings.TrimSpace(out) != "{\n  \"components\": null,\n  \"project\": \"kit\"\n}" {
		t.Errorf("--unset --json: %d %q %s", code, out, errOut)
	}
	unset := strings.Replace(before, "    components:\n      include: [src/**/*.svelte]\n      exclude: [src/internal/**]\n      setup: src/app.css\n", "", 1)
	if m := readManifest(t, root); m != unset {
		t.Errorf("after both --unset:\n%s\nwant:\n%s", m, unset)
	}
}

// S-0437: config refuses a glob that does not parse, a setup that names no
// file, a sub-project of an unsupported kind, an unknown sub-project, and
// flags that do not go together, and writes nothing.
func TestComponentsConfigRefuses(t *testing.T) {
	root, before := componentsFixture(t)
	for _, tc := range []struct {
		name   string
		args   []string
		out    string
		errOut string
	}{
		{"bad glob", []string{"ui", "--include", "src/[lib"}, `projects[0].components.include: (sub-project "ui") pattern "src/[lib"`, "rule: "},
		{"missing setup", []string{"kit", "--setup", "src/nope.css"}, `"src/nope.css" names no file under packages/kit`, "rule: "},
		{"unsupported kind", []string{"cli", "--include", "src/**/*.svelte"}, `whose components flai does not read`, "rule: "},
		{"unknown sub-project", []string{"web", "--setup", "src/app.css"}, "", `has no sub-project \"web\"; name one of ui, cli, kit`},
		{"no flag", []string{"ui"}, "", "nothing to write for sub-project"},
		{"unset with another", []string{"ui", "--unset", "--setup", "src/app.css"}, "", "takes no other flag"},
	} {
		out, errOut, code := runIn(t, root, append([]string{"components", "config"}, tc.args...)...)
		if code != 1 || !strings.Contains(out, tc.out) || !strings.Contains(errOut, tc.errOut) {
			t.Errorf("%s: %d %q %s", tc.name, code, out, errOut)
		}
		if m := readManifest(t, root); m != before {
			t.Errorf("%s changed the manifest:\n%s", tc.name, m)
		}
	}
	out, _, code := runIn(t, root, "components", "config", "ui", "--include", "/abs/*.svelte", "--json")
	var got struct {
		Refused []manifest.Problem `json:"refused"`
	}
	if code != 1 || json.Unmarshal([]byte(out), &got) != nil || len(got.Refused) != 1 || !strings.Contains(got.Refused[0].Reason, "absolute path") {
		t.Errorf("--json refusal: %d %q", code, out)
	}
}
