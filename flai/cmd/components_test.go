package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"slices"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/components"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
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

// S-0443: each component carries its examples file's examples and its path,
// given whether or not the file exists; a file that does not fit the format
// gives the examples that parsed; the text counts a component's examples.
func TestComponentsListExamples(t *testing.T) {
	root, _ := componentsFixture(t)
	useNode(t, componentsNode())
	writeExamples := func(rel, data string) {
		t.Helper()
		p := filepath.Join(root, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte(data), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	writeExamples("design/components/ui/src/lib/Button.examples.json", `{
  "examples": [
    {"name": "Small", "props": {"label": "Save", "size": "sm"}},
    {"name": "Large, dark, narrow", "props": {"label": "Delete", "size": "lg"}, "theme": "dark", "width": 375}
  ]
}
`)
	// kit is the sub-project's name and packages/kit its path: the file is
	// found by the name. Its second example's props are not an object.
	writeExamples("design/components/kit/src/Card.examples.json", `{"examples": [{"name": "Empty", "props": {"items": []}, "colour": "red"}, {"name": "Broken", "props": "x"}]}`)

	out, errOut, code := runIn(t, root, "components", "list", "--json")
	if code != 0 {
		t.Fatalf("exit %d: %s %s", code, out, errOut)
	}
	var got struct {
		Projects []struct {
			Name       string `json:"name"`
			Components []struct {
				Path     string `json:"path"`
				Examples []struct {
					Name  string                     `json:"name"`
					Props map[string]json.RawMessage `json:"props"`
					Theme string                     `json:"theme"`
					Width int                        `json:"width"`
				} `json:"examples"`
				ExamplesPath *string `json:"examples_path"`
			} `json:"components"`
		} `json:"projects"`
	}
	if err := json.Unmarshal([]byte(out), &got); err != nil {
		t.Fatalf("not JSON: %v\n%s", err, out)
	}
	ui, kit := got.Projects[0], got.Projects[2]
	button, field, card := ui.Components[0], ui.Components[1], kit.Components[0]
	if button.ExamplesPath == nil || *button.ExamplesPath != "design/components/ui/src/lib/Button.examples.json" || len(button.Examples) != 2 {
		t.Fatalf("Button: %+v", button)
	}
	if small := button.Examples[0]; small.Name != "Small" || string(small.Props["label"]) != `"Save"` || small.Theme != "" || small.Width != 0 {
		t.Errorf("Small: %+v", small)
	}
	if large := button.Examples[1]; large.Name != "Large, dark, narrow" || string(large.Props["size"]) != `"lg"` || large.Theme != "dark" || large.Width != 375 {
		t.Errorf("Large: %+v", large)
	}
	if field.ExamplesPath == nil || *field.ExamplesPath != "design/components/ui/src/lib/forms/Field.examples.json" || field.Examples == nil || len(field.Examples) != 0 {
		t.Errorf("a component without a file has examples [] and the path its first goes to: %+v", field)
	}
	if card.ExamplesPath == nil || *card.ExamplesPath != "design/components/kit/src/Card.examples.json" || len(card.Examples) != 2 ||
		card.Examples[0].Name != "Empty" || string(card.Examples[0].Props["items"]) != "[]" || card.Examples[1].Props == nil || len(card.Examples[1].Props) != 0 {
		t.Errorf("a file that does not fit the format gives what parsed, props {} where they are not an object: %+v", card)
	}
	if strings.Contains(out, `"examples": null`) || strings.Contains(out, `"props": null`) || strings.Contains(out, `"theme": ""`) || strings.Contains(out, `"width": 0`) {
		t.Errorf("examples and props are never null, an absent theme and width left out:\n%s", out)
	}

	out, errOut, code = runIn(t, root, "components", "list")
	if code != 0 {
		t.Fatalf("exit %d: %s %s", code, out, errOut)
	}
	for _, line := range []string{"\n  src/lib/Button.svelte (2 examples)\n", "\n  src/lib/forms/Field.svelte\n", "\n  src/Card.svelte (2 examples)\n"} {
		if !strings.Contains(out, line) {
			t.Errorf("text has %q:\n%s", line, out)
		}
	}

	// A file that cannot be read fails the list, saying which.
	if err := os.Remove(filepath.Join(root, "design/components/kit/src/Card.examples.json")); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(root, "design/components/kit/src/Card.examples.json"), 0o755); err != nil {
		t.Fatal(err)
	}
	out, errOut, code = runIn(t, root, "components", "list")
	if code != 1 || !strings.Contains(errOut, "cannot read the examples file design/components/kit/src/Card.examples.json") {
		t.Errorf("an unreadable examples file: %d %q %s", code, out, errOut)
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

// TestHelperComponentsVite is the fake vite flai components serve's tests
// start: the test binary run again by the script fakeComponentsVite writes,
// with FLAI_FAKE_COMPONENTS_VITE naming what it does: serve answers 200 at
// the harness's path on the --port it is given, and crash does so and exits
// soon after. It writes its pid and its arguments, one per line, into
// FLAI_FAKE_COMPONENTS_VITE_DIR first. Run as a test, it does nothing.
func TestHelperComponentsVite(t *testing.T) {
	mode := os.Getenv("FLAI_FAKE_COMPONENTS_VITE")
	if mode == "" {
		return
	}
	args := os.Args
	if i := slices.Index(args, "--"); i >= 0 {
		args = args[i+1:]
	}
	state := os.Getenv("FLAI_FAKE_COMPONENTS_VITE_DIR")
	_ = os.WriteFile(filepath.Join(state, "args"), []byte(strings.Join(args, "\n")), 0o600)
	_ = os.WriteFile(filepath.Join(state, "pid"), []byte(strconv.Itoa(os.Getpid())), 0o600)
	port := ""
	if i := slices.Index(args, "--port"); i >= 0 && i+1 < len(args) {
		port = args[i+1]
	}
	if mode == "crash" {
		go func() {
			time.Sleep(300 * time.Millisecond)
			fmt.Fprintln(os.Stderr, "the fake vite crashed")
			os.Exit(2)
		}()
	}
	mux := http.NewServeMux()
	mux.HandleFunc(components.HarnessPath, func(w http.ResponseWriter, _ *http.Request) {
		_, _ = io.WriteString(w, "harness")
	})
	fmt.Println("VITE ready")
	err := http.ListenAndServe("127.0.0.1:"+port, mux)
	fmt.Fprintln(os.Stderr, "cannot serve:", err)
	os.Exit(1)
}

// fakeComponentsVite installs, as the node_modules/.bin/vite of the
// sub-project at sub under root, a script that runs this test binary as
// TestHelperComponentsVite in mode, and returns the folder it writes into.
func fakeComponentsVite(t *testing.T, root, sub, mode string) string {
	t.Helper()
	if runtime.GOOS == "windows" {
		t.Skip("the fake vite is a POSIX shell script")
	}
	exe, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	bin := filepath.Join(root, filepath.FromSlash(sub), "node_modules", ".bin")
	if err := os.MkdirAll(bin, 0o755); err != nil {
		t.Fatal(err)
	}
	script := fmt.Sprintf("#!/bin/sh\nexec '%s' '-test.run=^TestHelperComponentsVite$' -- \"$@\"\n", exe)
	if err := os.WriteFile(filepath.Join(bin, "vite"), []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	state := t.TempDir()
	t.Setenv("FLAI_FAKE_COMPONENTS_VITE", mode)
	t.Setenv("FLAI_FAKE_COMPONENTS_VITE_DIR", state)
	return state
}

// lockedBuffer is a bytes.Buffer safe to write from the command and vite's
// output while the test reads it.
type lockedBuffer struct {
	mu sync.Mutex
	b  bytes.Buffer
}

func (l *lockedBuffer) Write(p []byte) (int, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.b.Write(p)
}

func (l *lockedBuffer) String() string {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.b.String()
}

// served is a flai components serve running in the background.
type served struct {
	out, errOut *lockedBuffer
	cancel      context.CancelFunc
	done        chan error
}

// serveComponentsIn runs flai with args in root in the background, with a
// context the test cancels as an interrupt would.
func serveComponentsIn(t *testing.T, root string, args ...string) *served {
	t.Helper()
	t.Setenv("FLAI_CACHE_DIR", filepath.Join(t.TempDir(), "cache"))
	ctx, cancel := context.WithCancel(context.Background())
	s := &served{out: &lockedBuffer{}, errOut: &lockedBuffer{}, cancel: cancel, done: make(chan error, 1)}
	a := &app{out: s.out, errOut: s.errOut, cwd: root}
	cmd := newRootCmdWith(a)
	cmd.SetArgs(args)
	go func() { s.done <- cmd.ExecuteContext(ctx) }()
	t.Cleanup(func() {
		cancel()
		select {
		case <-s.done:
		case <-time.After(15 * time.Second):
			t.Error("flai components serve did not stop")
		}
	})
	return s
}

// waitOutput waits for the command to have written what ok accepts on
// standard output, and fails the test when it ends or takes too long.
func (s *served) waitOutput(t *testing.T, ok func(string) bool) string {
	t.Helper()
	deadline := time.Now().Add(15 * time.Second)
	for {
		out := s.out.String()
		if ok(out) {
			return out
		}
		select {
		case err := <-s.done:
			t.Fatalf("flai components serve ended: %v\n%s\n%s", err, out, s.errOut.String())
		default:
		}
		if time.Now().After(deadline) {
			t.Fatalf("flai components serve printed only %q\n%s", out, s.errOut.String())
		}
		time.Sleep(20 * time.Millisecond)
	}
}

// stop cancels the command's context and returns its error.
func (s *served) stop(t *testing.T) error {
	t.Helper()
	s.cancel()
	select {
	case err := <-s.done:
		s.done <- err
		return err
	case <-time.After(15 * time.Second):
		t.Fatal("flai components serve did not stop on interrupt")
		return nil
	}
}

// fakeVitePID is the pid the fake vite wrote into state.
func fakeVitePID(t *testing.T, state string) int {
	t.Helper()
	b, err := os.ReadFile(filepath.Join(state, "pid"))
	if err != nil {
		t.Fatalf("the fake vite wrote no pid: %v", err)
	}
	pid, err := strconv.Atoi(string(b))
	if err != nil {
		t.Fatal(err)
	}
	return pid
}

// processGone reports whether no process answers a signal 0 at pid.
func processGone(pid int) bool {
	p, err := os.FindProcess(pid)
	if err != nil {
		return true
	}
	defer func() { _ = p.Release() }()
	return p.Signal(syscall.Signal(0)) != nil
}

// oneSvelteFixture is a scratch project whose manifest names a SvelteKit
// app, ui, and a Go sub-project, cli, so that ui is the one supported.
func oneSvelteFixture(t *testing.T) string {
	t.Helper()
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	root := tempProject(t)
	file := filepath.Join(root, manifest.File)
	data, err := os.ReadFile(file)
	if err != nil {
		t.Fatal(err)
	}
	data = append(data, "projects:\n  - name: ui\n    path: ui\n    kind: sveltekit\n    components:\n      setup: src/app.css\n  - name: cli\n    path: cli\n    kind: go\n"...)
	if err := os.WriteFile(file, data, 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(root, "ui", "src"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "ui", "src", "app.css"), []byte("p {}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	return root
}

var harnessURL = regexp.MustCompile(`^http://127\.0\.0\.1:(\d+)/__flai/components/$`)

// S-0438: serve picks the one Svelte sub-project, writes its harness config,
// runs its own vite on it, prints the harness's URL once it answers, and
// stops vite when interrupted.
func TestComponentsServePrintsTheURLAndStopsOnInterrupt(t *testing.T) {
	root := oneSvelteFixture(t)
	state := fakeComponentsVite(t, root, "ui", "serve")
	s := serveComponentsIn(t, root, "components", "serve")
	out := s.waitOutput(t, func(out string) bool { return strings.Count(out, "\n") >= 2 })
	lines := strings.Split(strings.TrimSuffix(out, "\n"), "\n")
	pid := fakeVitePID(t, state)
	if len(lines) != 2 || !harnessURL.MatchString(lines[0]) || lines[1] != fmt.Sprintf("  sub-project ui (ui), Vite pid %d; Ctrl-C stops it", pid) {
		t.Fatalf("output %q", out)
	}
	resp, err := http.Get(lines[0])
	if err != nil {
		t.Fatal(err)
	}
	_ = resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("GET %s: %d", lines[0], resp.StatusCode)
	}
	config := filepath.Join(root, ".flai-cache", "components", "t", "ui", components.HarnessConfig)
	written, err := os.ReadFile(config)
	if err != nil || !strings.Contains(string(written), `"setup":"src/app.css"`) {
		t.Errorf("the harness config at %s: %v\n%s", config, err, written)
	}
	args, _ := os.ReadFile(filepath.Join(state, "args"))
	if !strings.HasPrefix(string(args), "--config\n"+config+"\n--host\n127.0.0.1\n") {
		t.Errorf("vite ran with\n%s", args)
	}
	if !strings.Contains(s.errOut.String(), "VITE ready") {
		t.Errorf("vite's output is not on standard error: %q", s.errOut.String())
	}
	if err := s.stop(t); err != nil {
		t.Errorf("an interrupt is not an error: %v", err)
	}
	if !processGone(pid) {
		if p, err := os.FindProcess(pid); err == nil {
			_ = p.Kill()
		}
		t.Errorf("the fake vite, pid %d, is still running after the interrupt", pid)
	}
}

// S-0438: --project names the sub-project among several, and --json prints
// {project, url, port, pid}.
func TestComponentsServeProjectJSON(t *testing.T) {
	root, _ := componentsFixture(t)
	state := fakeComponentsVite(t, root, "packages/kit", "serve")
	s := serveComponentsIn(t, root, "components", "serve", "--project", "kit", "--json")
	var got struct {
		Project string `json:"project"`
		URL     string `json:"url"`
		Port    int    `json:"port"`
		PID     int    `json:"pid"`
	}
	s.waitOutput(t, func(out string) bool { return json.Unmarshal([]byte(out), &got) == nil })
	m := harnessURL.FindStringSubmatch(got.URL)
	if got.Project != "kit" || m == nil || m[1] != strconv.Itoa(got.Port) || got.PID != fakeVitePID(t, state) {
		t.Errorf("--json: %+v", got)
	}
	if _, err := os.Stat(filepath.Join(root, ".flai-cache", "components", "t", "kit", components.HarnessConfig)); err != nil {
		t.Errorf("kit's harness config: %v", err)
	}
	if err := s.stop(t); err != nil {
		t.Errorf("stop: %v", err)
	}
	if !processGone(got.PID) {
		t.Errorf("the fake vite, pid %d, is still running", got.PID)
	}
}

// S-0438: a dev server that exits on its own ends the command with why.
func TestComponentsServeEndsWhenViteExits(t *testing.T) {
	root := oneSvelteFixture(t)
	fakeComponentsVite(t, root, "ui", "crash")
	s := serveComponentsIn(t, root, "components", "serve")
	select {
	case err := <-s.done:
		s.done <- err
		if err == nil || !strings.Contains(err.Error(), "exited: exit status 2") || !strings.Contains(err.Error(), "the fake vite crashed") {
			t.Errorf("error %v, want vite's exit with its last lines", err)
		}
	case <-time.After(15 * time.Second):
		t.Fatal("flai components serve did not end when vite exited")
	}
}

// S-0438: serve refuses an unknown, an unsupported, or an ambiguous
// sub-project, and one without vite installed, and starts nothing.
func TestComponentsServeRefuses(t *testing.T) {
	root, _ := componentsFixture(t)
	for _, tc := range []struct {
		name   string
		args   []string
		errOut string
	}{
		{"unknown", []string{"--project", "web"}, `system-flow.yaml has no sub-project \"web\"; name one of ui, cli, kit`},
		{"unsupported", []string{"--project", "cli"}, `sub-project \"cli\" is of kind go`},
		{"ambiguous", nil, "has 2 sub-projects of kind sveltekit or svelte (ui, kit)"},
		{"no vite", []string{"--project", "ui"}, "cannot serve sub-project ui's components: vite is not installed in " + filepath.Join(root, "ui")},
	} {
		out, errOut, code := runIn(t, root, append([]string{"components", "serve"}, tc.args...)...)
		if code != 1 || out != "" || !strings.Contains(errOut, tc.errOut) {
			t.Errorf("%s: %d %q %s", tc.name, code, out, errOut)
		}
	}
	if _, err := os.Stat(filepath.Join(root, ".flai-cache", "components", "t", "cli")); err == nil {
		t.Error("an unsupported sub-project got a harness")
	}
}

// copyInstall is a package manager that installs a sub-project's
// dependencies by copying the script at vite into the folder it runs in,
// and remembers each folder it ran in.
type copyInstall struct {
	vite string
	dirs []string
}

func (c *copyInstall) LookPath(name string) (string, error) { return "/usr/bin/" + name, nil }

func (c *copyInstall) Run(_ context.Context, dir, _, _ string, _ ...string) ([]byte, []byte, error) {
	c.dirs = append(c.dirs, dir)
	script, err := os.ReadFile(c.vite)
	if err != nil {
		return nil, nil, err
	}
	bin := filepath.Join(dir, "node_modules", ".bin")
	if err := os.MkdirAll(bin, 0o755); err != nil {
		return nil, nil, err
	}
	return nil, nil, os.WriteFile(filepath.Join(bin, "vite"), script, 0o755)
}

// useInstaller makes flai components serve --story run i in place of the
// package manager for the test.
func useInstaller(t *testing.T, i components.Runner) {
	t.Helper()
	old := componentsInstaller
	t.Cleanup(func() { componentsInstaller = old })
	componentsInstaller = i
}

// storyWorktrees gives the project at root an epic and stories: S-0001 in
// progress and S-0002 ready, each with a worktree that holds the main
// checkout's manifest and ui with a pnpm lock file and no dependencies,
// S-0003 in progress without one, and a task, T-0001.
func storyWorktrees(t *testing.T, root string) {
	t.Helper()
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	now := time.Date(2026, 10, 16, 8, 0, 0, 0, time.UTC)
	if _, err := repo.Create(workitem.NewOptions{Type: workitem.Epic, Title: "Yard", Owner: "olive", Now: now}); err != nil {
		t.Fatal(err)
	}
	m, err := os.ReadFile(filepath.Join(root, manifest.File))
	if err != nil {
		t.Fatal(err)
	}
	for _, st := range []struct {
		states   []string
		worktree bool
	}{
		{[]string{workitem.Ready, workitem.InProgress}, true},
		{[]string{workitem.Ready}, true},
		{[]string{workitem.Ready, workitem.InProgress}, false},
	} {
		it, err := repo.Create(workitem.NewOptions{Type: workitem.Story, Title: "Crane", Parent: "E-0001", Owner: "olive", Now: now})
		if err != nil {
			t.Fatal(err)
		}
		body, _ := os.ReadFile(it.Path)
		_ = os.WriteFile(it.Path, []byte(strings.Replace(string(body), "## Acceptance criteria\n- [ ]\n", "## Acceptance criteria\n- [x] works\n", 1)), 0o644)
		it, _ = repo.Get(it.ID)
		for i, to := range st.states {
			if _, err := repo.Transition(it, to, "olive", "", now.Add(time.Duration(i+1)*time.Minute)); err != nil {
				t.Fatalf("%s to %s: %v", it.ID, to, err)
			}
		}
		if !st.worktree {
			continue
		}
		wt := repo.WorktreePath(it.ID)
		for rel, data := range map[string]string{manifest.File: string(m), "ui/src/app.css": "p {}\n", "ui/pnpm-lock.yaml": "lockfileVersion: '9.0'\n"} {
			p := filepath.Join(wt, filepath.FromSlash(rel))
			if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(p, []byte(data), 0o644); err != nil {
				t.Fatal(err)
			}
		}
	}
	if _, err := repo.Create(workitem.NewOptions{Type: workitem.Task, Title: "Rig", Parent: "S-0001", Owner: "olive", Now: now}); err != nil {
		t.Fatal(err)
	}
}

// S-0442: --story serves the sub-project in the story's worktree, after
// installing its dependencies there by its lock file, and --json names the
// story.
func TestComponentsServeStoryRunsInItsWorktree(t *testing.T) {
	root := oneSvelteFixture(t)
	storyWorktrees(t, root)
	state := fakeComponentsVite(t, root, "ui", "serve")
	install := &copyInstall{vite: filepath.Join(root, "ui", "node_modules", ".bin", "vite")}
	useInstaller(t, install)
	wt := filepath.Join(root, ".flai-cache", "worktrees", "S-0001")
	s := serveComponentsIn(t, root, "components", "serve", "--story", "S-1", "--json")
	var got struct {
		Project string `json:"project"`
		Story   string `json:"story"`
		URL     string `json:"url"`
		PID     int    `json:"pid"`
	}
	s.waitOutput(t, func(out string) bool { return json.Unmarshal([]byte(out), &got) == nil })
	if got.Project != "ui" || got.Story != "S-0001" || !harnessURL.MatchString(got.URL) || got.PID != fakeVitePID(t, state) {
		t.Errorf("--story --json: %+v", got)
	}
	if len(install.dirs) != 1 || install.dirs[0] != filepath.Join(wt, "ui") {
		t.Errorf("the worktree's dependencies are installed in it: %v", install.dirs)
	}
	config := filepath.Join(wt, ".flai-cache", "components", "t", "ui", components.HarnessConfig)
	if args, _ := os.ReadFile(filepath.Join(state, "args")); !strings.HasPrefix(string(args), "--config\n"+config+"\n") {
		t.Errorf("vite ran on the worktree's harness config %s:\n%s", config, args)
	}
	if err := s.stop(t); err != nil {
		t.Errorf("stop: %v", err)
	}
	if !processGone(got.PID) {
		t.Errorf("the fake vite, pid %d, is still running", got.PID)
	}
}

// S-0442, ADR-0145: --story refuses what is not a story in progress or in
// review with a worktree, saying how to open one, and starts nothing.
func TestComponentsServeStoryRefuses(t *testing.T) {
	root := oneSvelteFixture(t)
	storyWorktrees(t, root)
	install := &copyInstall{}
	useInstaller(t, install)
	for _, tc := range []struct{ story, errOut string }{
		{"S-0099", "S-0099 not found; name a story in progress or in review with a worktree"},
		{"T-0001", "T-0001 is a task; name a story (S-nnnn) in progress or in review"},
		{"S-0002", "S-0002 is ready: flai serves a story's components only while it is in progress or in review"},
		{"S-0003", "S-0003 has no worktree; open it with flai stream open S-0003"},
	} {
		out, errOut, code := runIn(t, root, "components", "serve", "--story", tc.story)
		if code != 1 || out != "" || !strings.Contains(errOut, tc.errOut) {
			t.Errorf("--story %s: %d %q %s", tc.story, code, out, errOut)
		}
	}
	if len(install.dirs) != 0 {
		t.Errorf("a refused story installs nothing: %v", install.dirs)
	}
}
