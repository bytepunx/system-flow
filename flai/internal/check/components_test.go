package check

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// S-0437: a sub-project's components key that is not valid is an error on
// the projects key, naming the sub-project's index and field; a valid key,
// and none, are not reported.
func TestComponentsKeysAreChecked(t *testing.T) {
	const web = "projects:\n  - name: web\n    path: web\n    kind: sveltekit\n"
	for _, extra := range []string{
		web,
		web + "    components:\n      include: [\"src/lib/**/*.svelte\"]\n      setup: src/app.css\n",
	} {
		repo := planningProject(t, extra)
		if err := os.MkdirAll(filepath.Join(repo.Root, "web", "src"), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(repo.Root, "web", "src", "app.css"), nil, 0o644); err != nil {
			t.Fatal(err)
		}
		if got := findings(t, repo, "manifest.components"); len(got) != 0 {
			t.Errorf("valid components %q: %+v", extra, got)
		}
	}
	repo := planningProject(t, web+"    components:\n      setup: src/missing.css\n")
	got := findings(t, repo, "manifest.components")
	if len(got) != 1 {
		t.Fatalf("want one error, the setup: %+v", got)
	}
	if f := got[0]; f.Level != Error || f.Path != "system-flow.yaml" || f.Line != 8 || !strings.Contains(f.Message, `projects[0].components.setup (sub-project "web") "src/missing.css" names no file`) {
		t.Errorf("finding: %+v, want an error on line 8 naming projects[0].components.setup", f)
	}
}

// fakeNode answers as node reading props would, without running it: props
// holds each component's props as JSON, errs a component's parse error, and
// answer, when set, the whole answer.
type fakeNode struct {
	missing bool
	props   map[string]string
	errs    map[string]string
	answer  string
	looked  int
	dirs    []string
	asked   [][]string
}

func (f *fakeNode) LookPath(name string) (string, error) {
	f.looked++
	if f.missing {
		return "", exec.ErrNotFound
	}
	return "/usr/bin/" + name, nil
}

func (f *fakeNode) Run(_ context.Context, dir, _, input string, _ ...string) ([]byte, []byte, error) {
	f.dirs = append(f.dirs, dir)
	var paths []string
	if err := json.Unmarshal([]byte(input), &paths); err != nil {
		return nil, nil, err
	}
	f.asked = append(f.asked, paths)
	if f.answer != "" {
		return []byte(f.answer), nil, nil
	}
	var b strings.Builder
	b.WriteString(`{"components":[`)
	for i, p := range paths {
		if i > 0 {
			b.WriteString(",")
		}
		props := f.props[p]
		if props == "" {
			props = "[]"
		}
		fmt.Fprintf(&b, `{"path":%q,"props":%s,"error":%q}`, p, props, f.errs[p])
	}
	b.WriteString("]}")
	return []byte(b.String()), nil, nil
}

// withNode has the check run f for node until the test ends.
func withNode(t *testing.T, f *fakeNode) {
	t.Helper()
	was := nodeRunner
	nodeRunner = f
	t.Cleanup(func() { nodeRunner = was })
}

// examplesProject is a project with a SvelteKit sub-project web, whose
// components are src/lib/Card.svelte and src/lib/Panel.svelte, a Go
// sub-project api, and files, slash paths under design/components to their
// content.
func examplesProject(t *testing.T, projects string, files map[string]string) *workitem.Repo {
	t.Helper()
	repo := planningProject(t, projects)
	for name, content := range map[string]string{"web/src/lib/Card.svelte": "", "web/src/lib/Panel.svelte": ""} {
		writeFile(t, filepath.Join(repo.Root, filepath.FromSlash(name)), content)
	}
	for name, content := range files {
		writeFile(t, filepath.Join(repo.Root, "design", "components", filepath.FromSlash(name)), content)
	}
	return repo
}

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

const (
	webAndAPI = "projects:\n  - name: web\n    path: web\n    kind: sveltekit\n  - name: api\n    path: api\n    kind: go\n"
	cardFile  = "web/src/lib/Card.examples.json"
	cardPath  = "design/components/web/src/lib/Card.examples.json"
	cardProps = `[{"name":"title","kind":"string","required":true},{"name":"open","kind":"boolean","required":false},{"name":"onclick","kind":"function","required":true}]`
)

// S-0443: an examples file that fits the format and its component's props
// has no finding, and node reads the props of only the components with
// examples, once for their sub-project.
func TestComponentExamplesThatFitHaveNoFinding(t *testing.T) {
	f := &fakeNode{props: map[string]string{"src/lib/Card.svelte": cardProps}}
	withNode(t, f)
	repo := examplesProject(t, webAndAPI, map[string]string{
		cardFile: "{\n  \"examples\": [\n    {\n      \"name\": \"Open\",\n      \"props\": { \"title\": \"A\", \"open\": true },\n      \"theme\": \"dark\",\n      \"width\": 375\n    }\n  ]\n}\n",
	})
	if got := findings(t, repo, "components.examples"); len(got) != 0 {
		t.Errorf("findings: %+v", got)
	}
	if len(f.asked) != 1 || !slices.Equal(f.asked[0], []string{"src/lib/Card.svelte"}) || f.dirs[0] != filepath.Join(repo.Root, "web") {
		t.Errorf("node asked %v in %v, want Card alone, once, in web", f.asked, f.dirs)
	}
}

// Each format problem and each example that does not fit its component's
// props is a warning on its line, naming the example; a file whose
// component's props could not be read is one, on line 1, with the error.
func TestComponentExamplesReportEachFindingOnItsLine(t *testing.T) {
	f := &fakeNode{
		props: map[string]string{"src/lib/Card.svelte": cardProps},
		errs:  map[string]string{"src/lib/Panel.svelte": "line 2, column 12: Unexpected token"},
	}
	withNode(t, f)
	repo := examplesProject(t, webAndAPI, map[string]string{
		cardFile: `{
  "examples": [
    {
      "name": "Wrong",
      "props": {
        "open": "yes",
        "colour": "red"
      },
      "theme": "blue"
    }
  ]
}
`,
		"web/src/lib/Panel.examples.json": `{"examples": [{"name": "A", "props": {}}]}`,
	})
	res, err := Run(repo, now)
	if err != nil {
		t.Fatal(err)
	}
	var got []Finding
	for _, fd := range res.Findings {
		if fd.Rule == "components.examples" {
			got = append(got, fd)
		}
	}
	want := []struct {
		path string
		line int
		text string
	}{
		{cardPath, 3, `example "Wrong": its theme "blue" is not "light" or "dark"`},
		{cardPath, 3, `example "Wrong": it leaves out prop "title", which its component requires; give it a string`},
		{cardPath, 6, `example "Wrong": prop "open" takes true or false, and "yes" does not fit`},
		{cardPath, 7, `example "Wrong": prop "colour" is not one its component takes; it takes title, open, onclick`},
		{"design/components/web/src/lib/Panel.examples.json", 1, "its examples are not checked, since the props of src/lib/Panel.svelte in web could not be read: line 2, column 12: Unexpected token"},
	}
	for _, w := range want {
		if !slices.ContainsFunc(got, func(fd Finding) bool {
			return filepath.ToSlash(fd.Path) == w.path && fd.Line == w.line && fd.Level == Warning && !fd.advisory && strings.Contains(fd.Message, w.text)
		}) {
			t.Errorf("no warning on %s:%d with %q in %+v", w.path, w.line, w.text, got)
		}
	}
	if len(got) != len(want) {
		t.Errorf("%d findings, want %d: %+v", len(got), len(want), got)
	}
	if len(f.asked) != 1 || len(f.asked[0]) != 2 {
		t.Errorf("node asked %v, want Card and Panel at once", f.asked)
	}
	if res.OK(true) {
		t.Error("the strict check passes over the warnings")
	}
}

// A file whose sub-project or component does not exist is a warning on line
// 1, and its examples' props are not read; JSON that does not parse is a
// warning on its line, as well.
func TestComponentExamplesWithoutTheirComponent(t *testing.T) {
	f := &fakeNode{}
	withNode(t, f)
	repo := examplesProject(t, webAndAPI, map[string]string{
		"api/src/lib/Card.examples.json":  `{"examples": []}`,
		"shop/src/lib/Card.examples.json": `{"examples": []}`,
		"web/src/lib/Gone.examples.json":  "{\n  \"examples\": [\n    {\n  ]\n}\n",
		"Loose.examples.json":             `{"examples": []}`,
	})
	got := findings(t, repo, "components.examples")
	for _, w := range []struct {
		path string
		line int
		text string
	}{
		{"design/components/api/src/lib/Card.examples.json", 1, `system-flow.yaml names no sub-project "api" of kind sveltekit or svelte, so the file's component src/lib/Card.svelte does not exist`},
		{"design/components/shop/src/lib/Card.examples.json", 1, `names no sub-project "shop"`},
		{"design/components/web/src/lib/Gone.examples.json", 1, `sub-project "web" has no component src/lib/Gone.svelte`},
		{"design/components/web/src/lib/Gone.examples.json", 4, "the file is not valid JSON"},
		{"design/components/Loose.examples.json", 1, "the file is in no sub-project's folder"},
	} {
		if !slices.ContainsFunc(got, func(fd Finding) bool {
			return filepath.ToSlash(fd.Path) == w.path && fd.Line == w.line && fd.Level == Warning && strings.Contains(fd.Message, w.text)
		}) {
			t.Errorf("no warning on %s:%d with %q in %+v", w.path, w.line, w.text, got)
		}
	}
	if len(got) != 5 {
		t.Errorf("%d findings, want 5: %+v", len(got), got)
	}
	if len(f.asked) != 0 {
		t.Errorf("node asked %v, want nothing, since no file has its component", f.asked)
	}
}

// A sub-project whose props node cannot read at all, for a reason other than
// svelte not installed, such as a svelte older than 5, has each of its files
// reported once, with the error.
func TestComponentExamplesWhosePropsCannotBeRead(t *testing.T) {
	withNode(t, &fakeNode{answer: `{"error":{"code":"svelte-old","message":"4.2.0"}}`})
	repo := examplesProject(t, webAndAPI, map[string]string{
		cardFile:                          `{"examples": [{"name": "A", "props": {}}, {"name": "B", "props": {}}]}`,
		"web/src/lib/Panel.examples.json": `{"examples": [{"name": "A", "props": {}}]}`,
	})
	got := findings(t, repo, "components.examples")
	if len(got) != 2 {
		t.Fatalf("findings: %+v, want one for each file", got)
	}
	for _, fd := range got {
		if fd.Line != 1 || fd.Level != Warning || fd.advisory || fd.unchecked || !strings.Contains(fd.Message, "could not be read: web has svelte 4.2.0") {
			t.Errorf("finding: %+v", fd)
		}
	}
}

// S-0446, I-0161: a sub-project whose svelte is not installed gets one note
// on its examples folder, which --strict passes over, not a warning on each
// of its files; each file's format and component are checked still.
func TestComponentExamplesWithoutSvelteGetOneNote(t *testing.T) {
	const missing = `{"error":{"code":"svelte-missing","message":"cannot resolve svelte"}}`
	withNode(t, &fakeNode{answer: missing})
	repo := examplesProject(t, webAndAPI, map[string]string{
		cardFile:                          `{"examples": [{"name": "A", "props": {}}, {"name": "B", "props": {}}]}`,
		"web/src/lib/Panel.examples.json": `{"examples": [{"name": "A", "props": {}}]}`,
	})
	res, err := Run(repo, now)
	if err != nil {
		t.Fatal(err)
	}
	var got []Finding
	for _, fd := range res.Findings {
		if fd.Rule == "components.examples" {
			got = append(got, fd)
		}
	}
	if len(got) != 1 {
		t.Fatalf("findings: %+v, want one note", got)
	}
	if note := got[0]; !note.advisory || !note.unchecked || note.Line != 1 || filepath.ToSlash(note.Path) != "design/components/web" ||
		!strings.Contains(note.Message, "the props of the examples under design/components/web are not checked against their components': svelte is not installed in web") ||
		!strings.Contains(note.Message, "run npm install (or its package manager's install) in web") {
		t.Errorf("note: %+v", note)
	}
	if res.Advisory < 1 {
		t.Errorf("advisory %d, want the note counted", res.Advisory)
	}

	withNode(t, &fakeNode{answer: missing})
	repo = examplesProject(t, webAndAPI, map[string]string{
		cardFile:                          `{"examples": [{"name": "A", "props": {}}]}`,
		"web/src/lib/Panel.examples.json": "{\"examples\": [\n  {\"name\": \"A\", \"props\": {}, \"width\": -1}\n]}",
		"web/src/lib/Gone.examples.json":  `{"examples": []}`,
	})
	got = findings(t, repo, "components.examples")
	if len(got) != 3 {
		t.Fatalf("findings: %+v, want the note, the width, and the missing component", got)
	}
	for _, w := range []struct {
		path     string
		line     int
		advisory bool
		text     string
	}{
		{"design/components/web", 1, true, "svelte is not installed in web"},
		{"design/components/web/src/lib/Panel.examples.json", 2, false, `example "A": its width -1 is not a whole number`},
		{"design/components/web/src/lib/Gone.examples.json", 1, false, `sub-project "web" has no component src/lib/Gone.svelte`},
	} {
		if !slices.ContainsFunc(got, func(fd Finding) bool {
			return filepath.ToSlash(fd.Path) == w.path && fd.Line == w.line && fd.advisory == w.advisory && strings.Contains(fd.Message, w.text)
		}) {
			t.Errorf("no finding on %s:%d (advisory %v) with %q in %+v", w.path, w.line, w.advisory, w.text, got)
		}
	}
}

// Without node, or without a sub-project whose components flai reads, a
// project with examples gets one note that --strict passes over, saying what
// was not checked; without node, the format is checked still.
func TestComponentExamplesNoteWhatIsNotChecked(t *testing.T) {
	f := &fakeNode{missing: true}
	withNode(t, f)
	repo := examplesProject(t, webAndAPI, map[string]string{
		cardFile:                          `{"examples": [{"name": "A", "props": {"title": 1}}]}`,
		"web/src/lib/Panel.examples.json": "{\"examples\": [\n  {\"name\": \"A\", \"props\": {}, \"width\": -1}\n]}",
	})
	got := findings(t, repo, "components.examples")
	if len(got) != 2 {
		t.Fatalf("findings: %+v, want the note and the width", got)
	}
	note, width := got[0], got[1]
	if !note.advisory || !note.unchecked || note.Line != 1 || filepath.ToSlash(note.Path) != "design/components" || !strings.Contains(note.Message, "are not checked against their components': flai reads a component's props with node, which is not on PATH; install Node.js") {
		t.Errorf("note: %+v", note)
	}
	if width.advisory || width.Line != 2 || !strings.Contains(width.Message, `example "A": its width -1 is not a whole number`) {
		t.Errorf("width: %+v", width)
	}
	if len(f.asked) != 0 {
		t.Errorf("node asked %v without node", f.asked)
	}

	repo = examplesProject(t, "projects:\n  - name: api\n    path: api\n    kind: go\n", map[string]string{cardFile: `{"examples": "broken"}`})
	res, err := Run(repo, now)
	if err != nil {
		t.Fatal(err)
	}
	var notes []Finding
	for _, fd := range res.Findings {
		if fd.Rule == "components.examples" {
			notes = append(notes, fd)
		}
	}
	if len(notes) != 1 || !notes[0].advisory || notes[0].unchecked || !strings.Contains(notes[0].Message, "the 1 examples files under design/components are not checked: system-flow.yaml names no sub-project of kind sveltekit or svelte") {
		t.Errorf("findings without a Svelte sub-project: %+v, want one note", notes)
	}
	if res.Advisory < 1 {
		t.Errorf("advisory %d, want the note counted", res.Advisory)
	}
}

// A project without examples files is not checked: no finding, and node is
// not looked for.
func TestNoComponentExamplesNoCheck(t *testing.T) {
	f := &fakeNode{missing: true}
	withNode(t, f)
	repo := examplesProject(t, webAndAPI, map[string]string{"README.md": "# Components\n", "web/src/lib/Card.json": "{}"})
	if got := findings(t, repo, "components.examples"); len(got) != 0 {
		t.Errorf("findings: %+v", got)
	}
	if f.looked != 0 {
		t.Errorf("node looked for %d times", f.looked)
	}
}
