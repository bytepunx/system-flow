package components

import (
	"encoding/json"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/bytepunx/system-flow/flai/internal/manifest"
)

// S-0443: a component's examples file is under the design folder's
// components folder, at its sub-project's name and its path with .svelte made
// .examples.json, and the file names its component back.
func TestExamplesPathMirrorsTheComponent(t *testing.T) {
	m := manifest.Manifest{Layout: map[string]string{"design": "design"}}
	if got := ExamplesPath(m, "flaiover", "src/lib/components/Card.svelte"); got != "design/components/flaiover/src/lib/components/Card.examples.json" {
		t.Errorf("ExamplesPath = %q", got)
	}
	if got := ExamplesPath(manifest.Manifest{Layout: map[string]string{"design": "docs/design"}}, "web", "Card.svelte"); got != "docs/design/components/web/Card.examples.json" {
		t.Errorf("ExamplesPath under another design folder = %q", got)
	}
	if got := ExamplesFolder(manifest.Manifest{}); got != "design/components" {
		t.Errorf("ExamplesFolder with no layout = %q", got)
	}
	project, comp, ok := ExamplesComponent("flaiover/src/lib/components/Card.examples.json")
	if !ok || project != "flaiover" || comp != "src/lib/components/Card.svelte" {
		t.Errorf("ExamplesComponent = %q %q %v", project, comp, ok)
	}
	for _, rel := range []string{"Card.examples.json", "web/Card.json", "../web/Card.examples.json"} {
		if _, _, ok := ExamplesComponent(rel); ok {
			t.Errorf("ExamplesComponent(%q) names a component", rel)
		}
	}
}

const validExamples = `{
  "examples": [
    {
      "name": "A long title",
      "props": {
        "title": "Quarterly report for the north region",
        "open": true
      },
      "theme": "dark",
      "width": 375
    },
    {
      "name": "Empty",
      "props": {}
    }
  ]
}
`

// A file that fits the format gives each example with its line and each
// prop's line, and no problem.
func TestParseExamplesReadsAValidFile(t *testing.T) {
	got, problems := ParseExamples([]byte(validExamples))
	if len(problems) != 0 {
		t.Fatalf("problems: %+v", problems)
	}
	if len(got) != 2 {
		t.Fatalf("examples: %+v", got)
	}
	a, b := got[0], got[1]
	if a.Name != "A long title" || a.Index != 1 || a.Line != 3 || a.Theme != "dark" || a.Width != 375 {
		t.Errorf("first: %+v", a)
	}
	if string(a.Props["title"]) != `"Quarterly report for the north region"` || string(a.Props["open"]) != "true" {
		t.Errorf("first's props: %v", a.Props)
	}
	if a.PropLines["title"] != 6 || a.PropLines["open"] != 7 {
		t.Errorf("first's prop lines: %v", a.PropLines)
	}
	if b.Name != "Empty" || b.Index != 2 || b.Line != 12 || b.Props == nil || len(b.Props) != 0 || b.Theme != "" || b.Width != 0 {
		t.Errorf("second: %+v", b)
	}
	out, err := json.Marshal(b)
	if err != nil || string(out) != `{"name":"Empty","props":{}}` {
		t.Errorf("second as JSON: %s %v", out, err)
	}
}

// problemAt is the problem on line whose sentence holds want, or false.
func problemAt(problems []Problem, line int, want string) bool {
	for _, p := range problems {
		if p.Line == line && strings.Contains(p.String(), want) {
			return true
		}
	}
	return false
}

// Each way a file can miss the format is a problem at its line: the
// example's opening brace for one about the example, the key for one of the
// file's own.
func TestParseExamplesReportsTheFormat(t *testing.T) {
	data := `{
  "examples": [
    { "name": "", "props": {} },
    { "name": "Twice", "props": [] },
    { "name": "Twice", "props": {}, "theme": "blue", "width": 0 },
    { "props": {}, "colour": 1, "width": 2.5 },
    { "name": 7 },
    "a string"
  ],
  "extra": true
}`
	got, problems := ParseExamples([]byte(data))
	for _, w := range []struct {
		line int
		want string
	}{
		{3, "example 1: its name is empty"},
		{4, `example "Twice": its props are an array, not an object`},
		{5, `example "Twice": its name is example 2's too, on line 4`},
		{5, `example "Twice": its theme "blue" is not "light" or "dark"`},
		{5, `example "Twice": its width 0 is not a whole number`},
		{6, `example 4: key "colour" is not one an example has`},
		{6, "example 4: its width 2.5 is not a whole number"},
		{6, "example 4: it has no name"},
		{7, "example 5: its name is a number, not a string"},
		{7, "example 5: it has no props"},
		{8, "example 6: is a string, not an object"},
		{10, `the file's key "extra" is not one an examples file has`},
	} {
		if !problemAt(problems, w.line, w.want) {
			t.Errorf("no problem on line %d with %q in %+v", w.line, w.want, problems)
		}
	}
	if len(problems) != 12 {
		t.Errorf("%d problems, want 12: %+v", len(problems), problems)
	}
	if len(got) != 5 || got[1].Props != nil {
		t.Errorf("examples: %+v, want the five objects, the second without props", got)
	}
}

// JSON that does not parse is one problem at the line of its syntax error;
// a file that is not an object, or has no examples, is one at its start.
func TestParseExamplesReportsJSONAndShape(t *testing.T) {
	for _, tc := range []struct {
		data string
		line int
		want string
	}{
		{"{\n  \"examples\": [\n    { \"name\": \"A\", }\n  ]\n}", 3, "is not valid JSON: invalid character '}'"},
		{"{\n  \"examples\": [\n", 2, "is not valid JSON: unexpected end of JSON input"},
		{"", 1, "is not valid JSON"},
		{"{} {}", 1, "is not valid JSON"},
		{"\n[]", 2, "the file is an array, not an object holding examples"},
		{"{\n}", 1, "the file has no examples"},
		{`{"examples": {}}`, 1, "examples is an object, not an array"},
	} {
		got, problems := ParseExamples([]byte(tc.data))
		if len(problems) != 1 || !problemAt(problems, tc.line, tc.want) {
			t.Errorf("%q: problems %+v, want one on line %d with %q", tc.data, problems, tc.line, tc.want)
		}
		if got == nil || len(got) != 0 {
			t.Errorf("%q: examples %+v, want none", tc.data, got)
		}
	}
}

// Examples reads a file from a file system; an absent one has none.
func TestExamplesReadsAFile(t *testing.T) {
	fsys := fstest.MapFS{"design/components/web/Card.examples.json": {Data: []byte(validExamples)}}
	got, problems, err := Examples(fsys, "design/components/web/Card.examples.json")
	if err != nil || len(problems) != 0 || len(got) != 2 {
		t.Errorf("Examples = %+v, %+v, %v", got, problems, err)
	}
	got, problems, err = Examples(fsys, "design/components/web/Missing.examples.json")
	if err != nil || len(problems) != 0 || got == nil || len(got) != 0 {
		t.Errorf("Examples of an absent file = %+v, %+v, %v", got, problems, err)
	}
}

// example is an example of props, each a JSON value, with each prop on a
// line of its own after its opening brace on line 1.
func example(name string, props ...string) Example {
	ex := Example{Name: name, Index: 1, Line: 1, Props: map[string]json.RawMessage{}, PropLines: map[string]int{}}
	for i := 0; i+1 < len(props); i += 2 {
		ex.Props[props[i]] = json.RawMessage(props[i+1])
		ex.PropLines[props[i]] = 2 + i/2
	}
	return ex
}

var cardProps = []Prop{
	{Name: "title", Kind: KindString, Required: true},
	{Name: "count", Kind: KindNumber},
	{Name: "open", Kind: KindBoolean},
	{Name: "size", Kind: KindUnion, Literals: []string{`"sm"`, `"md"`, "2", "true"}},
	{Name: "tags", Kind: KindArray},
	{Name: "item", Kind: KindObject},
	{Name: "children", Kind: KindSnippet},
	{Name: "onclick", Kind: KindFunction, Required: true},
	{Name: "data", Kind: KindUnknown, Required: true},
}

// A value fits its prop's kind by ADR-0164's table, null fitting all but a
// function; a prop the component does not take, a function given a value,
// and a required prop left out are problems too.
func TestValidateChecksPropsByKind(t *testing.T) {
	for _, tc := range []struct {
		prop, value string
		fits        bool
	}{
		{"title", `"a"`, true}, {"title", "1", false}, {"title", "null", true},
		{"children", `"text"`, true}, {"children", "{}", false},
		{"count", "1.5", true}, {"count", `"1"`, false},
		{"open", "false", true}, {"open", `"yes"`, false},
		{"size", `"md"`, true}, {"size", "2", true}, {"size", "2.0", true}, {"size", "true", true},
		{"size", `"lg"`, false}, {"size", "false", false}, {"size", `["sm"]`, false}, {"size", "null", true},
		{"tags", "[]", true}, {"tags", "{}", false},
		{"item", `{"id": 1}`, true}, {"item", "[]", false},
		{"data", `[1, {"a": null}]`, true},
		{"onclick", "null", false},
	} {
		ex := example("E", "title", `"t"`, "data", "1", tc.prop, tc.value)
		got := Validate([]Example{ex}, cardProps)
		if tc.fits && len(got) != 0 {
			t.Errorf("%s = %s: %+v, want it to fit", tc.prop, tc.value, got)
		}
		if !tc.fits && (len(got) != 1 || got[0].Line != ex.PropLines[tc.prop] || !strings.Contains(got[0].String(), `example "E": prop "`+tc.prop+`"`)) {
			t.Errorf("%s = %s: %+v, want one problem on its line %d", tc.prop, tc.value, got, ex.PropLines[tc.prop])
		}
	}
	ex := example("Broken", "size", `"lg"`, "colour", `"red"`, "onclick", `"f"`)
	got := Validate([]Example{ex, {Name: "No props", Index: 2, Line: 9}}, cardProps)
	for _, w := range []struct {
		line int
		want string
	}{
		{1, `it leaves out prop "title", which its component requires; give it a string`},
		{1, `it leaves out prop "data"`},
		{2, `prop "size" takes one of "sm", "md", 2, true, and "lg" does not fit`},
		{3, `prop "colour" is not one its component takes; it takes title, count, open`},
		{4, `prop "onclick" is a function, which the page gives a logger; leave it out`},
	} {
		if !problemAt(got, w.line, `example "Broken": `+w.want) {
			t.Errorf("no problem on line %d with %q in %+v", w.line, w.want, got)
		}
	}
	if len(got) != 5 || got[0].Line != 1 || got[4].Line != 4 {
		t.Errorf("%d problems in line order, want 5: %+v", len(got), got)
	}
	if got := Validate([]Example{example("None")}, nil); len(got) != 0 {
		t.Errorf("a component without props, an example without props: %+v", got)
	}
	if got := Validate([]Example{example("Extra", "x", "1")}, nil); len(got) != 1 || !strings.Contains(got[0].Message, "it takes none") {
		t.Errorf("a prop of a component without props: %+v", got)
	}
}
