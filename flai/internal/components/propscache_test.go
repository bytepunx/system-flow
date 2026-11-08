package components

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"slices"
	"sync"
	"testing"
	"time"
)

// countingNode answers as the script would, without running it: each file
// on its input gets one string prop named after its sub-project, and a file
// named in errs that error. A file named in deps is stamped, and the answer
// gives the stamps of the dependency files, taken as the script takes them.
// It records each start's sub-project folder and files.
type countingNode struct {
	mu     sync.Mutex
	deps   []string
	errs   map[string]string
	starts []countedStart
}

type countedStart struct {
	dir   string
	files []string
}

func (f *countingNode) LookPath(name string) (string, error) { return "/usr/bin/" + name, nil }

func (f *countingNode) Run(_ context.Context, dir, _, input string, _ ...string) ([]byte, []byte, error) {
	var files []string
	if err := json.Unmarshal([]byte(input), &files); err != nil {
		return nil, nil, err
	}
	f.mu.Lock()
	f.starts = append(f.starts, countedStart{dir: dir, files: files})
	f.mu.Unlock()
	var a struct {
		Components []answered `json:"components"`
		Stamps     [][2]any   `json:"stamps"`
	}
	stamped := false
	for _, file := range files {
		c := answered{Path: file, Props: []Prop{{Name: filepath.Base(dir), Kind: KindString}}}
		if e, ok := f.errs[file]; ok {
			c.Props, c.Error = []Prop{}, e
		}
		if slices.Contains(f.deps, file) {
			c.Stamped, stamped = true, true
		}
		a.Components = append(a.Components, c)
	}
	if stamped {
		for _, d := range f.deps {
			if filepath.IsAbs(d) {
				s, _ := stampOf("file", d)
				a.Stamps = append(a.Stamps, [2]any{[2]string{"file", d}, s})
			}
		}
	}
	out, err := json.Marshal(a)
	return out, nil, err
}

// started is how many times node started, and forgets them.
func (f *countingNode) started() []countedStart {
	f.mu.Lock()
	defer f.mu.Unlock()
	s := f.starts
	f.starts = nil
	return s
}

// keptFixture is a project with a .flai-cache folder and two sub-projects,
// web and docs, each with svelte installed: web with A.svelte, which depends
// on types.ts, and B.svelte, docs with C.svelte. It gives the root and the
// components.
func keptFixture(t *testing.T) (string, []Component) {
	t.Helper()
	root := t.TempDir()
	files := map[string]string{
		".flai-cache/.keep":                     "",
		"web/package.json":                      `{"name": "web"}`,
		"web/node_modules/svelte/package.json":  `{"name": "svelte", "version": "5.1.0"}`,
		"web/src/A.svelte":                      "<script>let { web } = $props();</script>",
		"web/src/B.svelte":                      "<script>let { web } = $props();</script>",
		"web/src/types.ts":                      "export type T = { id: string };",
		"docs/package.json":                     `{"name": "docs"}`,
		"docs/node_modules/svelte/package.json": `{"name": "svelte", "version": "5.2.0"}`,
		"docs/src/C.svelte":                     "<script>let { docs } = $props();</script>",
	}
	for name, body := range files {
		p := filepath.Join(root, filepath.FromSlash(name))
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return root, slices.Concat(comps("web", "src/A.svelte", "src/B.svelte"), comps("docs", "src/C.svelte"))
}

// touch rewrites path with body and a modification time a minute on.
func touch(t *testing.T, path, body string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	later := time.Now().Add(time.Minute)
	if err := os.Chtimes(path, later, later); err != nil {
		t.Fatal(err)
	}
}

func readAll(t *testing.T, f *countingNode, root string, in []Component) []Component {
	t.Helper()
	got, err := ReadProps(context.Background(), f, root, in)
	if err != nil {
		t.Fatal(err)
	}
	return got
}

func TestReadPropsStartsNodeOnlyForAComponentWhoseKeyChanged(t *testing.T) {
	root, in := keptFixture(t)
	types := filepath.Join(root, "web", "src", "types.ts")
	f := &countingNode{deps: []string{"src/A.svelte", types}}

	first := readAll(t, f, root, in)
	if s := f.started(); len(s) != 2 {
		t.Fatalf("first read started node %d times, want once per sub-project", len(s))
	}
	if _, err := os.Stat(KeptPropsPath(root, "web")); err != nil {
		t.Fatalf("no kept read for web after the first read: %v", err)
	}

	second := readAll(t, f, root, in)
	if s := f.started(); len(s) != 0 {
		t.Fatalf("a read with nothing changed started node for %+v, want none", s)
	}
	for i := range first {
		if first[i].Path != second[i].Path || !slices.EqualFunc(first[i].Props, second[i].Props, func(a, b Prop) bool { return a.Name == b.Name && a.Kind == b.Kind }) {
			t.Errorf("kept %+v, want what node answered, %+v", second[i], first[i])
		}
	}

	webDir, docsDir := filepath.Join(root, "web"), filepath.Join(root, "docs")
	tests := []struct {
		name   string
		change func()
		want   []countedStart
	}{
		{"a component's file", func() {
			touch(t, filepath.Join(webDir, "src", "B.svelte"), "<script>let { web, more } = $props();</script>")
		},
			[]countedStart{{dir: webDir, files: []string{"src/B.svelte"}}}},
		{"a file its skeletons read", func() { touch(t, types, "export type T = { id: string; title: string };") },
			[]countedStart{{dir: webDir, files: []string{"src/A.svelte"}}}},
		{"the svelte version", func() {
			touch(t, filepath.Join(docsDir, "node_modules", "svelte", "package.json"), `{"name": "svelte", "version": "5.3.0"}`)
		}, []countedStart{{dir: docsDir, files: []string{"src/C.svelte"}}}},
		{"the sub-project's package.json", func() {
			touch(t, filepath.Join(webDir, "package.json"), `{"name": "web", "devDependencies": {"@sveltejs/kit": "^2"}}`)
		},
			[]countedStart{{dir: webDir, files: []string{"src/A.svelte", "src/B.svelte"}}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.change()
			readAll(t, f, root, in)
			if s := f.started(); !slices.EqualFunc(s, tt.want, func(a, b countedStart) bool { return a.dir == b.dir && slices.Equal(a.files, b.files) }) {
				t.Errorf("after a change to %s node started %+v, want %+v", tt.name, s, tt.want)
			}
			readAll(t, f, root, in)
			if s := f.started(); len(s) != 0 {
				t.Errorf("the read after that started node %+v, want none", s)
			}
		})
	}
}

func TestReadPropsKeepsAComponentsError(t *testing.T) {
	root, in := keptFixture(t)
	f := &countingNode{errs: map[string]string{"src/B.svelte": "line 1, column 2: Unexpected token"}}
	readAll(t, f, root, in)
	f.started()
	got := readAll(t, f, root, in)
	if s := f.started(); len(s) != 0 {
		t.Fatalf("started node %+v, want none", s)
	}
	if got[1].Error != "line 1, column 2: Unexpected token" || got[1].Props == nil || len(got[1].Props) != 0 {
		t.Errorf("B %+v, want its kept error and no props", got[1])
	}
}

func TestReadPropsKeepsNothingWithoutAFlaiCacheFolder(t *testing.T) {
	root, in := keptFixture(t)
	if err := os.RemoveAll(filepath.Join(root, ".flai-cache")); err != nil {
		t.Fatal(err)
	}
	f := &countingNode{}
	readAll(t, f, root, in)
	readAll(t, f, root, in)
	if s := f.started(); len(s) != 4 {
		t.Errorf("node started %d times over two reads, want twice each read", len(s))
	}
	if _, err := os.Stat(filepath.Join(root, ".flai-cache")); !errors.Is(err, os.ErrNotExist) {
		t.Errorf("a .flai-cache folder was made (%v), want none (D-0402)", err)
	}
}

func TestReadPropsKeepsNothingWithoutSvelte(t *testing.T) {
	root, in := keptFixture(t)
	if err := os.RemoveAll(filepath.Join(root, "web", "node_modules")); err != nil {
		t.Fatal(err)
	}
	f := &countingNode{}
	readAll(t, f, root, in)
	f.started()
	readAll(t, f, root, in)
	if s := f.started(); len(s) != 1 || s[0].dir != filepath.Join(root, "web") {
		t.Errorf("node started %+v, want for web alone, whose svelte is not found", s)
	}
}

func TestReadPropsTakesAKeptFileThatDoesNotParseAsAMiss(t *testing.T) {
	root, in := keptFixture(t)
	f := &countingNode{}
	readAll(t, f, root, in)
	f.started()
	if err := os.WriteFile(KeptPropsPath(root, "web"), []byte(`{"version": 1, "components": [`), 0o644); err != nil {
		t.Fatal(err)
	}
	readAll(t, f, root, in)
	if s := f.started(); len(s) != 1 || len(s[0].files) != 2 {
		t.Errorf("node started %+v, want web's two components read again", s)
	}
	readAll(t, f, root, in)
	if s := f.started(); len(s) != 0 {
		t.Errorf("node started %+v once the kept file was written again, want none", s)
	}
}

func TestReadPropsAtOnceKeepsEveryAnswer(t *testing.T) {
	root, in := keptFixture(t)
	f := &countingNode{}
	var wg sync.WaitGroup
	for range 4 {
		wg.Go(func() {
			if _, err := ReadProps(context.Background(), f, root, in); err != nil {
				t.Error(err)
			}
		})
	}
	wg.Wait()
	f.started()
	readAll(t, f, root, in)
	if s := f.started(); len(s) != 0 {
		t.Errorf("node started %+v after reads at once, want none", s)
	}
}

func TestStampOfGivesTheScriptsStamps(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "a.ts")
	if err := os.WriteFile(file, []byte("abc"), 0o644); err != nil {
		t.Fatal(err)
	}
	at := time.Unix(1792953279, 129979280)
	if err := os.Chtimes(file, at, at); err != nil {
		t.Fatal(err)
	}
	tests := []struct{ kind, path, want string }{
		{"file", file, "f:1792953279129979280:3"},
		{"file", filepath.Join(dir, "none.ts"), "-"},
		{"dir", dir, "d"},
		{"dir", file, "f"},
		{"real", filepath.Join(dir, "none.ts"), "-"},
	}
	for _, tt := range tests {
		if got, ok := stampOf(tt.kind, tt.path); !ok || got != tt.want {
			t.Errorf("stampOf(%s, %s) = %q, %v; want %q", tt.kind, tt.path, got, ok, tt.want)
		}
	}
	if got, _ := stampOf("list", dir); len(got) < 3 || got[:2] != "d:" {
		t.Errorf("stampOf(list, dir) = %q, want d:<nanoseconds>", got)
	}
	if _, ok := stampOf("other", dir); ok {
		t.Error("stampOf took a kind the script does not give")
	}
}

// The script's stamps and stampOf's agree, so that a kept read of a typed
// component stays kept: a second List with nothing changed starts no node.
func TestListWithRealNodeKeepsATypedComponentsRead(t *testing.T) {
	root := fixture(t, "typed", "svelte", "typescript")
	if err := os.Mkdir(filepath.Join(root, ".flai-cache"), 0o755); err != nil {
		t.Fatal(err)
	}
	r := &countingSystem{}
	if _, err := List(context.Background(), r, root, "web", nil, nil); err != nil {
		t.Fatal(err)
	}
	if _, err := List(context.Background(), r, root, "web", nil, nil); err != nil {
		t.Fatal(err)
	}
	if r.starts != 1 {
		t.Fatalf("node started %d times over two Lists with nothing changed, want once", r.starts)
	}
	types := filepath.Join(root, "web", "src", "lib", "types.ts")
	body, err := os.ReadFile(types)
	if err != nil {
		t.Fatal(err)
	}
	writeLater(t, types, string(body)+"\nexport type Extra = { n: number };\n")
	if _, err := List(context.Background(), r, root, "web", nil, nil); err != nil {
		t.Fatal(err)
	}
	if r.starts != 2 {
		t.Errorf("node started %d times in all once an imported type's file changed, want twice", r.starts)
	}
}

// countingSystem is System, counting node's starts.
type countingSystem struct {
	System
	starts int
}

func (c *countingSystem) Run(ctx context.Context, dir, name, input string, args ...string) ([]byte, []byte, error) {
	c.starts++
	return c.System.Run(ctx, dir, name, input, args...)
}
