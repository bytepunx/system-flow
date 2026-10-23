package components

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
	"testing"
	"time"
)

// fakeNode answers as node would without running it.
type fakeNode struct {
	missing bool
	stdout  string
	stderr  string
	err     error
	block   bool
	calls   []fakeCall
}

type fakeCall struct {
	dir, input string
	args       []string
}

func (f *fakeNode) LookPath(name string) (string, error) {
	if f.missing {
		return "", exec.ErrNotFound
	}
	return "/usr/bin/" + name, nil
}

func (f *fakeNode) Run(ctx context.Context, dir, name, input string, args ...string) ([]byte, []byte, error) {
	f.calls = append(f.calls, fakeCall{dir: dir, input: input, args: append([]string{name}, args...)})
	if f.block {
		<-ctx.Done()
		return nil, nil, ctx.Err()
	}
	return []byte(f.stdout), []byte(f.stderr), f.err
}

func comps(project string, files ...string) []Component {
	out := make([]Component, len(files))
	for i, f := range files {
		out[i] = Component{Project: project, Path: f, Name: strings.TrimSuffix(filepath.Base(f), ".svelte"), Props: []Prop{}}
	}
	return out
}

func TestReadPropsMapsNodesAnswerOntoEachComponent(t *testing.T) {
	f := &fakeNode{stdout: `{"components":[
		{"path":"src/lib/A.svelte","props":[
			{"name":"size","kind":"union","type":"'sm' | 'md'","default":"'md'","required":false,"literals":["\"sm\"","\"md\""]},
			{"name":"label","kind":"string","type":"string","default":"","required":true}]},
		{"path":"src/lib/Bad.svelte","props":[],"error":"line 2, column 12: Unexpected token"}]}`}
	in := comps("web", "src/lib/A.svelte", "src/lib/Bad.svelte", "src/lib/Lost.svelte")
	got, err := ReadProps(context.Background(), f, "/repo", in)
	if err != nil {
		t.Fatal(err)
	}
	a := got[0]
	want := []Prop{
		{Name: "size", Kind: KindUnion, Type: "'sm' | 'md'", Default: "'md'", Literals: []string{`"sm"`, `"md"`}},
		{Name: "label", Kind: KindString, Type: "string", Required: true},
	}
	if len(a.Props) != len(want) {
		t.Fatalf("props %+v, want %+v", a.Props, want)
	}
	for i := range want {
		p, w := a.Props[i], want[i]
		if p.Name != w.Name || p.Kind != w.Kind || p.Type != w.Type || p.Default != w.Default || p.Required != w.Required || !slices.Equal(p.Literals, w.Literals) {
			t.Errorf("prop %d %+v, want %+v", i, p, w)
		}
	}
	if b := got[1]; len(b.Props) != 0 || b.Error != "line 2, column 12: Unexpected token" {
		t.Errorf("unparsable component %+v, want no props and its parse error", b)
	}
	if l := got[2]; l.Error != "node gave no answer for it" {
		t.Errorf("unanswered component %+v, want an error saying so", l)
	}
	if in[0].Props == nil || len(in[0].Props) != 0 {
		t.Error("ReadProps changed the components it was given")
	}

	call := f.calls[0]
	if call.dir != filepath.Join("/repo", "web") {
		t.Errorf("ran in %q, want the sub-project's folder", call.dir)
	}
	if len(call.args) != 6 || call.args[0] != "node" || call.args[1] != "--input-type=module" || call.args[2] != "-e" || call.args[3] != script {
		t.Errorf("ran %q, want node --input-type=module -e <script> <cache> <hash>", call.args[:min(3, len(call.args))])
	} else if call.args[4] != filepath.Join("/repo", ".flai-cache", "component-props", "web") || call.args[5] != scriptHash || len(scriptHash) != 64 {
		t.Errorf("ran the script with %q, want the sub-project's cache folder and the script's hash", call.args[4:])
	}
	var files []string
	if err := json.Unmarshal([]byte(call.input), &files); err != nil || !slices.Equal(files, []string{"src/lib/A.svelte", "src/lib/Bad.svelte", "src/lib/Lost.svelte"}) {
		t.Errorf("stdin %q (%v), want the files as a JSON array", call.input, err)
	}
}

func TestReadPropsRunsNodeOncePerSubProject(t *testing.T) {
	f := &fakeNode{stdout: `{"components":[]}`}
	in := slices.Concat(comps("web", "A.svelte"), comps("docs", "B.svelte"), comps("web", "C.svelte"))
	if _, err := ReadProps(context.Background(), f, "/repo", in); err != nil {
		t.Fatal(err)
	}
	if len(f.calls) != 2 || f.calls[0].input != `["A.svelte","C.svelte"]` || f.calls[1].input != `["B.svelte"]` {
		t.Errorf("calls %+v, want web's two files, then docs' one", f.calls)
	}
}

func TestReadPropsWithNothingToReadNeedsNoNode(t *testing.T) {
	f := &fakeNode{missing: true}
	got, err := ReadProps(context.Background(), f, "/repo", nil)
	if err != nil || len(got) != 0 {
		t.Errorf("got %v, %v; want nothing and no error", got, err)
	}
}

func TestReadPropsWithoutNodeSaysToInstallIt(t *testing.T) {
	_, err := ReadProps(context.Background(), &fakeNode{missing: true}, "/repo", comps("web", "A.svelte"))
	if !errors.Is(err, ErrNodeMissing) || !strings.Contains(err.Error(), "install Node.js") {
		t.Errorf("err %v, want ErrNodeMissing saying to install Node.js", err)
	}
}

func TestReadPropsReportsWhatStopsASubProject(t *testing.T) {
	tests := []struct {
		name string
		f    *fakeNode
		is   error
		want []string
	}{
		{"svelte missing", &fakeNode{stdout: `{"error":{"code":"svelte-missing","message":"cannot resolve svelte"}}`},
			ErrSvelteMissing, []string{"svelte is not installed in web", "run npm install", "in web"}},
		{"svelte old", &fakeNode{stdout: `{"error":{"code":"svelte-old","message":"4.2.0"}}`},
			nil, []string{"web has svelte 4.2.0", "upgrade svelte"}},
		{"node fails", &fakeNode{err: errors.New("exit status 1"), stderr: "SyntaxError: bad\n"},
			nil, []string{"node failed", "web", "exit status 1", "SyntaxError: bad"}},
		{"not JSON", &fakeNode{stdout: "Warning: hello"},
			nil, []string{"not the JSON flai expects", "web"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ReadProps(context.Background(), tt.f, "/repo", comps("web", "A.svelte"))
			if err == nil {
				t.Fatal("no error")
			}
			if tt.is != nil && !errors.Is(err, tt.is) {
				t.Errorf("err %v, want it to wrap %v", err, tt.is)
			}
			for _, w := range tt.want {
				if !strings.Contains(err.Error(), w) {
					t.Errorf("err %q, want it to say %q", err, w)
				}
			}
		})
	}
}

func TestReadPropsStopsNodeThatTakesTooLong(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()
	_, err := ReadProps(ctx, &fakeNode{block: true}, "/repo", comps("web", "A.svelte"))
	if err == nil || !strings.Contains(err.Error(), "did not finish") || !strings.Contains(err.Error(), "web") {
		t.Errorf("err %v, want one saying node did not finish in web", err)
	}
}

// needNode skips the test when node is not on PATH.
func needNode(t *testing.T) {
	t.Helper()
	if _, err := exec.LookPath("node"); err != nil {
		t.Skip("node is not on PATH")
	}
}

func TestReadPropsWithRealNodeAndNoSvelteSaysToInstallIt(t *testing.T) {
	needNode(t)
	root := t.TempDir()
	tree(t, root, "web/src/lib/A.svelte")
	if err := os.WriteFile(filepath.Join(root, "web", "package.json"), []byte(`{"name":"web"}`), 0o644); err != nil {
		t.Fatal(err)
	}
	_, err := List(context.Background(), System{}, root, "web", nil, nil)
	if !errors.Is(err, ErrSvelteMissing) || !strings.Contains(err.Error(), "npm install") {
		t.Errorf("err %v, want ErrSvelteMissing saying to run npm install", err)
	}
}

// installed is the package name of flaiover's node_modules in the nearest
// folder at or above this package's that has one, or "".
func installed(t *testing.T, name string) string {
	t.Helper()
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	for {
		p := filepath.Join(dir, "flaiover", "node_modules", name)
		if _, err := os.Stat(filepath.Join(p, "package.json")); err == nil {
			real, err := filepath.EvalSymlinks(p)
			if err != nil {
				t.Fatal(err)
			}
			return real
		}
		up := filepath.Dir(dir)
		if up == dir {
			return ""
		}
		dir = up
	}
}

// fixture is testdata/<dir> copied under a temporary root as sub-project
// web, with each of packages in its node_modules linked to flaiover's.
func fixture(t *testing.T, dir string, packages ...string) string {
	t.Helper()
	if testing.Short() {
		t.Skip("integration: reads flaiover's installed packages from the monorepo")
	}
	needNode(t)
	root := t.TempDir()
	web := filepath.Join(root, "web")
	if err := os.CopyFS(web, os.DirFS(filepath.Join("testdata", dir))); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(web, "node_modules"), 0o755); err != nil {
		t.Fatal(err)
	}
	for _, name := range packages {
		p := installed(t, name)
		if p == "" {
			t.Skipf("no flaiover/node_modules/%s above this package; run npm install in flaiover", name)
		}
		if err := os.Symlink(p, filepath.Join(web, "node_modules", name)); err != nil {
			t.Fatal(err)
		}
	}
	return root
}

func TestListReadsEachKindFromFixtureComponents(t *testing.T) {
	root := fixture(t, "project", "svelte")
	got, err := List(context.Background(), System{}, root, "web", nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	byPath := map[string]Component{}
	for _, c := range got {
		byPath[c.Path] = c
	}
	wantPaths := []string{"src/lib/Alias.svelte", "src/lib/Broken.svelte", "src/lib/Kinds.svelte", "src/lib/Static.svelte", "src/lib/Untyped.svelte", "src/lib/Whole.svelte", "src/lib/forms/Field.svelte"}
	if !slices.Equal(paths(got), wantPaths) {
		t.Fatalf("paths %v, want %v", paths(got), wantPaths)
	}

	tests := []struct {
		path string
		want []Prop
	}{
		{"src/lib/Kinds.svelte", []Prop{
			{Name: "label", Kind: KindString, Type: "string", Required: true},
			{Name: "count", Kind: KindNumber, Type: "number", Default: "3"},
			{Name: "open", Kind: KindBoolean, Type: "boolean", Default: "false"},
			{Name: "size", Kind: KindUnion, Type: "'sm' | 'md' | 'lg' | undefined", Default: "'md'", Literals: []string{`"sm"`, `"md"`, `"lg"`}},
			{Name: "tags", Kind: KindArray, Type: "string[]", Required: true},
			{Name: "ids", Kind: KindArray, Type: "Array<number>", Default: "[]"},
			{Name: "options", Kind: KindObject, Type: "{ dense: boolean }", Required: true},
			{Name: "meta", Kind: KindObject, Type: "Record<string, unknown>", Required: true},
			{Name: "onpick", Kind: KindFunction, Type: "(id: number) => void", Required: true},
			{Name: "children", Kind: KindSnippet, Type: "Snippet<[string]>"},
			{Name: "theme", Kind: KindUnknown, Type: "Theme", Required: true},
			{Name: "ratio", Kind: KindNumber, Type: "number | null"},
		}},
		{"src/lib/forms/Field.svelte", []Prop{
			{Name: "value", Kind: KindString, Type: "string", Default: "''"},
			{Name: "level", Kind: KindUnion, Type: "1 | 2 | 3", Required: true, Literals: []string{"1", "2", "3"}},
			{Name: "id", Kind: KindString, Type: "string", Required: true},
			{Name: "validate", Kind: KindFunction, Type: "(value: string): boolean", Required: true},
		}},
		{"src/lib/Alias.svelte", []Prop{
			{Name: "variant", Kind: KindUnion, Type: "Variant", Default: "'primary'", Literals: []string{`"primary"`, `"secondary"`}},
			{Name: "disabled", Kind: KindBoolean, Type: "boolean"},
		}},
		{"src/lib/Whole.svelte", []Prop{
			{Name: "title", Kind: KindString, Type: "string", Required: true},
			{Name: "subtitle", Kind: KindString, Type: "string"},
		}},
		{"src/lib/Untyped.svelte", []Prop{
			{Name: "name", Kind: KindString, Default: "'x'"},
			{Name: "step", Kind: KindNumber, Default: "1"},
			{Name: "label", Kind: KindUnknown, Required: true},
		}},
		{"src/lib/Static.svelte", []Prop{}},
	}
	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			c := byPath[tt.path]
			if c.Error != "" {
				t.Fatalf("error %q, want none", c.Error)
			}
			if len(c.Props) != len(tt.want) {
				t.Fatalf("props %+v, want %+v", c.Props, tt.want)
			}
			for i, w := range tt.want {
				p := c.Props[i]
				if p.Name != w.Name || p.Kind != w.Kind || p.Type != w.Type || p.Default != w.Default || p.Required != w.Required || !slices.Equal(p.Literals, w.Literals) {
					t.Errorf("prop %d %+v, want %+v", i, p, w)
				}
			}
		})
	}

	b := byPath["src/lib/Broken.svelte"]
	if len(b.Props) != 0 || !strings.Contains(b.Error, "line 2") || !strings.Contains(b.Error, "Unexpected token") {
		t.Errorf("unparsable component %+v, want no props and its parse error with its line", b)
	}
}

// sameJSON reports whether got is the JSON value want is.
func sameJSON(t *testing.T, got json.RawMessage, want string) bool {
	t.Helper()
	if len(got) == 0 {
		return false
	}
	var g, w any
	if err := json.Unmarshal(got, &g); err != nil {
		t.Fatalf("skeleton %s is not JSON: %v", got, err)
	}
	if err := json.Unmarshal([]byte(want), &w); err != nil {
		t.Fatalf("want %s is not JSON: %v", want, err)
	}
	return reflect.DeepEqual(g, w)
}

// skeletons are a component's props' skeletons by name, each as JSON, with
// "" for a prop that has none.
func skeletons(t *testing.T, got []Component, path string) map[string]string {
	t.Helper()
	for _, c := range got {
		if c.Path != path {
			continue
		}
		if c.Error != "" {
			t.Fatalf("%s: error %q, want none", path, c.Error)
		}
		out := map[string]string{}
		for _, p := range c.Props {
			out[p.Name] = string(p.Skeleton)
		}
		return out
	}
	t.Fatalf("no component %s in %v", path, paths(got))
	return nil
}

func checkSkeletons(t *testing.T, got map[string]string, want map[string]string) {
	t.Helper()
	if len(got) != len(want) {
		t.Errorf("props %v, want %v", got, want)
	}
	for name, w := range want {
		g, ok := got[name]
		switch {
		case !ok:
			t.Errorf("no prop %s", name)
		case w == "" && g != "":
			t.Errorf("%s: skeleton %s, want none", name, g)
		case w != "" && !sameJSON(t, json.RawMessage(g), w):
			t.Errorf("%s: skeleton %s, want %s", name, g, w)
		}
	}
}

const storySkeleton = `{"id": "", "points": 0, "done": false, "size": "S", "owner": "", "due": "", "tags": [""], "pair": ["", 0], "meta": {}}`

func TestListGivesEachJSONPropASkeletonOfItsResolvedType(t *testing.T) {
	root := fixture(t, "typed", "svelte", "typescript")
	got, err := List(context.Background(), System{}, root, "web", nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	checkSkeletons(t, skeletons(t, got, "src/lib/Board.svelte"), map[string]string{
		"label":   "",
		"story":   storySkeleton,
		"stories": "[" + storySkeleton + "]",
		"tree":    `{"name": "", "children": [{}], "parent": {}}`,
		"status":  `"open"`,
		"deep":    `{"a": {"b": {"c": {"d": {"e": {"f": {}}}}}}}`,
		"missing": `null`,
		"ghost":   `[]`,
		"extra":   `{}`,
		"local":   `{"at": "", "ok": false}`,
	})
	checkSkeletons(t, skeletons(t, got, "src/lib/Imported.svelte"), map[string]string{
		"story": storySkeleton,
		"count": `0`,
	})
}

func TestListMapsLibInASvelteKitSubProjectWhoseGeneratedTsconfigIsMissing(t *testing.T) {
	root := fixture(t, "typed", "svelte", "typescript")
	web := filepath.Join(root, "web")
	files := map[string]string{
		"tsconfig.json": `{"extends": "./.svelte-kit/tsconfig.json", "compilerOptions": {"strict": true}}`,
		"package.json":  `{"name": "kit", "private": true, "type": "module", "devDependencies": {"@sveltejs/kit": "^2.0.0"}}`,
	}
	for name, body := range files {
		if err := os.WriteFile(filepath.Join(web, name), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	got, err := List(context.Background(), System{}, root, "web", nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	if sk := skeletons(t, got, "src/lib/Imported.svelte"); !sameJSON(t, json.RawMessage(sk["story"]), storySkeleton) {
		t.Errorf("story: skeleton %s, want %s", sk["story"], storySkeleton)
	}
}

// cacheFiles are the files of the props cache of sub-project web under root.
func cacheFiles(t *testing.T, root string) []string {
	t.Helper()
	files, err := filepath.Glob(filepath.Join(PropsCacheDir(root, "web"), "*.json"))
	if err != nil {
		t.Fatal(err)
	}
	return files
}

// editCache sets the skeleton the cache file holds for prop of component.
func editCache(t *testing.T, file, component, prop, skeleton string) {
	t.Helper()
	data, err := os.ReadFile(file)
	if err != nil {
		t.Fatal(err)
	}
	var cache map[string]json.RawMessage
	if err := json.Unmarshal(data, &cache); err != nil {
		t.Fatal(err)
	}
	var sk map[string]map[string]json.RawMessage
	if err := json.Unmarshal(cache["skeletons"], &sk); err != nil {
		t.Fatal(err)
	}
	sk[component][prop] = json.RawMessage(skeleton)
	if cache["skeletons"], err = json.Marshal(sk); err != nil {
		t.Fatal(err)
	}
	if data, err = json.Marshal(cache); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(file, data, 0o644); err != nil {
		t.Fatal(err)
	}
}

// writeLater writes body to path with a modification time a minute on, so
// that the change is seen however coarse the file system's clock.
func writeLater(t *testing.T, path, body string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	later := time.Now().Add(time.Minute)
	for _, p := range []string{path, filepath.Dir(path)} {
		if err := os.Chtimes(p, later, later); err != nil {
			t.Fatal(err)
		}
	}
}

// I-0166: every run built a TypeScript program over everything the
// components import, about two seconds for flaiover, though nothing had
// changed since the last.
func TestListTakesSkeletonsFromTheCacheUntilAFileTypescriptReadChanges(t *testing.T) {
	root := fixture(t, "typed", "svelte", "typescript")
	if _, err := List(context.Background(), System{}, root, "web", nil, nil); err != nil {
		t.Fatal(err)
	}
	files := cacheFiles(t, root)
	if len(files) != 1 {
		t.Fatalf("cache files %v after one run, want one", files)
	}
	// A skeleton no type gives shows that the next run answers from the cache.
	editCache(t, files[0], "src/lib/Imported.svelte", "story", `"from the cache"`)
	got, err := List(context.Background(), System{}, root, "web", nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	if sk := skeletons(t, got, "src/lib/Imported.svelte")["story"]; sk != `"from the cache"` {
		t.Errorf("story: skeleton %s on a run with nothing changed, want the cache's", sk)
	}

	types := filepath.Join(root, "web", "src", "lib", "types.ts")
	body, err := os.ReadFile(types)
	if err != nil {
		t.Fatal(err)
	}
	writeLater(t, types, strings.Replace(string(body), "\tid: string;\n", "\tid: string;\n\ttitle: string;\n", 1))
	got, err = List(context.Background(), System{}, root, "web", nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	want := `{"id": "", "title": "", "points": 0, "done": false, "size": "S", "owner": "", "due": "", "tags": [""], "pair": ["", 0], "meta": {}}`
	if sk := skeletons(t, got, "src/lib/Imported.svelte")["story"]; !sameJSON(t, json.RawMessage(sk), want) {
		t.Errorf("story: skeleton %s after its imported type changed, want %s", sk, want)
	}
}

func TestListResolvesAModuleMadeAfterTheCacheMissedIt(t *testing.T) {
	root := fixture(t, "typed", "svelte", "typescript")
	got, err := List(context.Background(), System{}, root, "web", nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	if sk := skeletons(t, got, "src/lib/Board.svelte")["missing"]; sk != "null" {
		t.Fatalf("missing: skeleton %s, want null while $lib/missing does not exist", sk)
	}
	writeLater(t, filepath.Join(root, "web", "src", "lib", "missing.ts"), "export type Missing = { n: number };\n")
	if got, err = List(context.Background(), System{}, root, "web", nil, nil); err != nil {
		t.Fatal(err)
	}
	if sk := skeletons(t, got, "src/lib/Board.svelte")["missing"]; !sameJSON(t, json.RawMessage(sk), `{"n": 0}`) {
		t.Errorf("missing: skeleton %s once $lib/missing exists, want {\"n\": 0}", sk)
	}
}

func TestListKeepsACacheForEachSetOfComponents(t *testing.T) {
	root := fixture(t, "typed", "svelte", "typescript")
	all, err := Find(root, "web", nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	one := slices.DeleteFunc(slices.Clone(all), func(c Component) bool { return c.Path != "src/lib/Imported.svelte" })
	for _, comps := range [][]Component{all, one} {
		if _, err := ReadProps(context.Background(), System{}, root, comps); err != nil {
			t.Fatal(err)
		}
	}
	if files := cacheFiles(t, root); len(files) != 2 {
		t.Errorf("cache files %v after runs over two sets of components, want one for each", files)
	}
}

func TestListWithoutTypescriptGivesEachJSONPropItsKindsFallback(t *testing.T) {
	root := fixture(t, "project", "svelte")
	got, err := List(context.Background(), System{}, root, "web", nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	checkSkeletons(t, skeletons(t, got, "src/lib/Kinds.svelte"), map[string]string{
		"label": "", "count": "", "open": "", "size": "", "onpick": "", "children": "", "ratio": "",
		"tags": `[]`, "ids": `[]`, "options": `{}`, "meta": `{}`, "theme": `null`,
	})
	checkSkeletons(t, skeletons(t, got, "src/lib/Untyped.svelte"), map[string]string{
		"name": "", "step": "", "label": `null`,
	})
	for _, c := range got {
		if c.Path != "src/lib/Kinds.svelte" {
			continue
		}
		for _, p := range c.Props {
			if p.Name != "theme" {
				continue
			}
			out, err := json.Marshal(p)
			if err != nil {
				t.Fatal(err)
			}
			if !strings.Contains(string(out), `"skeleton":null`) {
				t.Errorf("theme marshals as %s, want its null skeleton kept", out)
			}
		}
	}
}
