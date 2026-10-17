package components

// A component's examples (S-0443, ADR-0164): named sets of its props, kept
// in one JSON file per component under the design folder's components
// folder. ParseExamples reads a file's content and says what in it does not
// fit the format; Validate says which examples do not fit the component's
// props. Neither does any I/O.

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"math"
	"path"
	"slices"
	"strings"

	"github.com/bytepunx/system-flow/flai/internal/manifest"
)

// ExamplesSuffix ends the name of every examples file, in place of the
// component's .svelte.
const ExamplesSuffix = ".examples.json"

// examplesUnder is the folder under the design folder that holds the
// examples files.
const examplesUnder = "components"

// ExamplesFolder is the repository path of the folder that holds the
// examples files: components in the design layout folder,
// design/components under the default layout.
func ExamplesFolder(m manifest.Manifest) string {
	d := m.Layout["design"]
	if d == "" {
		d = "design"
	}
	return path.Join(d, examplesUnder)
}

// ExamplesPath is the repository path of the examples file of the component
// at component, its slash path in the sub-project, of the sub-project named
// project in the manifest: ExamplesFolder, the sub-project's name, then the
// component's path with .svelte replaced by ExamplesSuffix. The file need
// not exist.
func ExamplesPath(m manifest.Manifest, project, component string) string {
	return path.Join(ExamplesFolder(m), project, strings.TrimSuffix(component, ".svelte")+ExamplesSuffix)
}

// ExamplesComponent is the sub-project's name and the component's path in it
// that rel, an examples file's slash path relative to ExamplesFolder, names:
// ExamplesPath the other way round. False when rel does not end in
// ExamplesSuffix or has no folder for its sub-project.
func ExamplesComponent(rel string) (project, component string, ok bool) {
	project, rest, found := strings.Cut(path.Clean(rel), "/")
	if !found || project == "" || project == "." || project == ".." || !strings.HasSuffix(rest, ExamplesSuffix) {
		return "", "", false
	}
	return project, strings.TrimSuffix(rest, ExamplesSuffix) + ".svelte", true
}

// Example is one named set of a component's props.
type Example struct {
	// Name names it, unique in its file.
	Name string `json:"name"`
	// Props are the props the page sends, prop name to JSON value; nil when
	// the example gives none or they are not an object.
	Props map[string]json.RawMessage `json:"props"`
	// Theme is the one theme the page shows, light or dark; empty for the
	// theme or themes it shows.
	Theme string `json:"theme,omitempty"`
	// Width is the frame's width in CSS pixels; 0 for the width it shows.
	Width int `json:"width,omitempty"`
	// Index is its place in the file's examples, from 1.
	Index int `json:"-"`
	// Line is the line its object opens on.
	Line int `json:"-"`
	// PropLines are the line of each prop's key.
	PropLines map[string]int `json:"-"`
}

// Problem is one thing in an examples file that does not fit the format or
// the component's props.
type Problem struct {
	// Line is where: the prop's key, the example's opening brace, a syntax
	// error, or a key of the file's own.
	Line int
	// Index is the example's place in the file, from 1; 0 for a problem of
	// the file's own.
	Index int
	// Example is the example's name, empty when it has none.
	Example string
	// Message says what is wrong and what to do.
	Message string
}

// String is the problem as a sentence that names its example.
func (p Problem) String() string {
	switch {
	case p.Index == 0:
		return p.Message
	case p.Example != "":
		return fmt.Sprintf("example %q: %s", p.Example, p.Message)
	default:
		return fmt.Sprintf("example %d: %s", p.Index, p.Message)
	}
}

// Examples reads the examples file name in fsys and parses it as
// ParseExamples does. An absent file has no examples and is no error; an
// error is for a file that cannot be read.
func Examples(fsys fs.FS, name string) ([]Example, []Problem, error) {
	data, err := fs.ReadFile(fsys, name)
	switch {
	case errors.Is(err, fs.ErrNotExist):
		return []Example{}, nil, nil
	case err != nil:
		return nil, nil, fmt.Errorf("cannot read the examples file %s: %w; check that it is a readable file", name, err)
	}
	ex, problems := ParseExamples(data)
	return ex, problems, nil
}

// ParseExamples is the examples in data, an examples file's content, each
// that is an object, in the file's order, and what in it does not fit the
// format: JSON that does not parse, a key the format does not name, a name
// that is missing, empty, or used twice, props that are missing or not an
// object, a theme that is not light or dark, and a width that is not a whole
// number, 1 or more.
func ParseExamples(data []byte) ([]Example, []Problem) {
	p := &parser{data: data}
	if !json.Valid(data) {
		var syntax *json.SyntaxError
		var whole any
		if err := json.Unmarshal(data, &whole); errors.As(err, &syntax) {
			p.file(p.lineAt(int(syntax.Offset)-1), "the file is not valid JSON: %s; fix the JSON there", syntax.Error())
		} else {
			p.file(1, "the file is not valid JSON: %v; fix the JSON", err)
		}
		return []Example{}, p.problems
	}
	start := len(data) - len(bytes.TrimLeft(data, " \t\r\n"))
	if kind := jsonKind(data[start:]); kind != "an object" {
		p.file(p.lineAt(start), `the file is %s, not an object holding examples; write {"examples": [...]}`, kind)
		return []Example{}, p.problems
	}
	examples := []Example{}
	found := false
	for _, m := range p.members(data[start:], start) {
		if m.key != "examples" {
			p.file(p.lineAt(m.keyAt), "the file's key %q is not one an examples file has; it holds examples only, so remove it", m.key)
			continue
		}
		found = true
		if kind := jsonKind(m.value); kind != "an array" {
			p.file(p.lineAt(m.keyAt), "examples is %s, not an array of examples; write it as [...]", kind)
			continue
		}
		examples = append(examples, p.examples(m.value, m.valueAt)...)
	}
	if !found {
		p.file(p.lineAt(start), `the file has no examples; write {"examples": []} to hold none`)
	}
	return examples, p.problems
}

// parser gathers the problems of one file's content as it reads it.
type parser struct {
	data     []byte
	problems []Problem
}

// member is one key of a JSON object and its value, with the offsets of
// both in the file.
type member struct {
	key            string
	keyAt, valueAt int
	value          json.RawMessage
}

// file adds a problem of the file's own at line.
func (p *parser) file(line int, format string, args ...any) {
	p.problems = append(p.problems, Problem{Line: line, Message: fmt.Sprintf(format, args...)})
}

// example adds a problem of ex at its opening brace.
func (p *parser) example(ex Example, format string, args ...any) {
	p.problems = append(p.problems, Problem{Line: ex.Line, Index: ex.Index, Example: ex.Name, Message: fmt.Sprintf(format, args...)})
}

// lineAt is the line, from 1, of the byte at offset in the file.
func (p *parser) lineAt(offset int) int {
	offset = min(max(offset, 0), len(p.data))
	return 1 + bytes.Count(p.data[:offset], []byte("\n"))
}

// members are the keys of raw, a valid JSON object found at offset at in the
// file, in their order.
func (p *parser) members(raw []byte, at int) []member {
	dec := json.NewDecoder(bytes.NewReader(raw))
	if _, err := dec.Token(); err != nil {
		return nil
	}
	var out []member
	for dec.More() {
		tok, err := dec.Token()
		if err != nil {
			return out
		}
		key, _ := tok.(string)
		keyAt := at + int(dec.InputOffset()) - 1
		var value json.RawMessage
		if err := dec.Decode(&value); err != nil {
			return out
		}
		out = append(out, member{key: key, keyAt: keyAt, valueAt: at + int(dec.InputOffset()) - len(value), value: value})
	}
	return out
}

// elements are the values of raw, a valid JSON array found at offset at in
// the file, each with its offset, in their order.
func (p *parser) elements(raw []byte, at int) []member {
	dec := json.NewDecoder(bytes.NewReader(raw))
	if _, err := dec.Token(); err != nil {
		return nil
	}
	var out []member
	for dec.More() {
		var value json.RawMessage
		if err := dec.Decode(&value); err != nil {
			return out
		}
		out = append(out, member{valueAt: at + int(dec.InputOffset()) - len(value), value: value})
	}
	return out
}

// examples are the examples of raw, the examples array found at offset at,
// each that is an object, with the problems of each and of names used twice.
func (p *parser) examples(raw []byte, at int) []Example {
	out := []Example{}
	named := map[string]Example{}
	for i, el := range p.elements(raw, at) {
		ex := Example{Index: i + 1, Line: p.lineAt(el.valueAt)}
		if kind := jsonKind(el.value); kind != "an object" {
			p.example(ex, "is %s, not an object with a name and props", kind)
			continue
		}
		ex = p.oneExample(ex, el)
		if first, ok := named[ex.Name]; ok && ex.Name != "" {
			p.example(ex, "its name is example %d's too, on line %d; give each example a name of its own", first.Index, first.Line)
		} else if ex.Name != "" {
			named[ex.Name] = ex
		}
		out = append(out, ex)
	}
	return out
}

// oneExample is ex read from el, an example's object, with its problems.
func (p *parser) oneExample(ex Example, el member) Example {
	members := p.members(el.value, el.valueAt)
	// The name first, so that every problem of the example names it.
	for _, m := range members {
		if m.key == "name" {
			_ = json.Unmarshal(m.value, &ex.Name)
		}
	}
	var hasName, hasProps bool
	for _, m := range members {
		switch m.key {
		case "name":
			hasName = true
			switch {
			case jsonKind(m.value) != "a string":
				p.example(ex, "its name is %s, not a string; give it a name", jsonKind(m.value))
			case strings.TrimSpace(ex.Name) == "":
				p.example(ex, "its name is empty; give it a name, unique in the file")
			}
		case "props":
			hasProps = true
			if kind := jsonKind(m.value); kind != "an object" {
				p.example(ex, "its props are %s, not an object of prop name to value; write {} for none", kind)
				ex.Props = nil
				continue
			}
			ex.Props, ex.PropLines = map[string]json.RawMessage{}, map[string]int{}
			for _, prop := range p.members(m.value, m.valueAt) {
				ex.Props[prop.key] = prop.value
				ex.PropLines[prop.key] = p.lineAt(prop.keyAt)
			}
		case "theme":
			var theme string
			if err := json.Unmarshal(m.value, &theme); err != nil || theme != "light" && theme != "dark" {
				p.example(ex, `its theme %s is not "light" or "dark"; give one of them, or leave theme out for the theme the page shows`, m.value)
				continue
			}
			ex.Theme = theme
		case "width":
			var width float64
			if err := json.Unmarshal(m.value, &width); err != nil || width < 1 || width != math.Trunc(width) || width > math.MaxInt32 {
				p.example(ex, "its width %s is not a whole number of CSS pixels, 1 or more; give one, such as 375, or leave width out for the width the page shows", m.value)
				continue
			}
			ex.Width = int(width)
		default:
			p.example(ex, "key %q is not one an example has; it has name, props, theme, and width only, so remove it", m.key)
		}
	}
	if !hasName {
		p.example(ex, "it has no name; give it one, unique in the file")
	}
	if !hasProps {
		p.example(ex, "it has no props; write props, {} for none")
	}
	return ex
}

// Validate is what in examples does not fit props, their component's: a prop
// the component does not take, a required prop that is not a function and
// that an example leaves out, and a value that does not fit its prop's kind,
// by example and then line. null fits every kind but a function's, which the
// page gives a logger. An example whose props are nil is passed over.
func Validate(examples []Example, props []Prop) []Problem {
	byName := make(map[string]Prop, len(props))
	names := make([]string, len(props))
	for i, pr := range props {
		byName[pr.Name] = pr
		names[i] = pr.Name
	}
	takes := "none"
	if len(names) > 0 {
		takes = strings.Join(names, ", ")
	}
	var out []Problem
	for _, ex := range examples {
		if ex.Props == nil {
			continue
		}
		var mine []Problem
		at := func(line int, format string, args ...any) {
			mine = append(mine, Problem{Line: line, Index: ex.Index, Example: ex.Name, Message: fmt.Sprintf(format, args...)})
		}
		for name, value := range ex.Props {
			line := ex.Line
			if l, ok := ex.PropLines[name]; ok {
				line = l
			}
			pr, ok := byName[name]
			switch {
			case !ok:
				at(line, "prop %q is not one its component takes; it takes %s", name, takes)
			case pr.Kind == KindFunction:
				at(line, "prop %q is a function, which the page gives a logger; leave it out", name)
			case !fits(pr, value):
				at(line, "prop %q takes %s, and %s does not fit; give %s", name, wants(pr), shown(value), wants(pr))
			}
		}
		for _, pr := range props {
			if _, given := ex.Props[pr.Name]; !given && pr.Required && pr.Kind != KindFunction {
				at(ex.Line, "it leaves out prop %q, which its component requires; give it %s", pr.Name, wants(pr))
			}
		}
		slices.SortStableFunc(mine, func(a, b Problem) int {
			if a.Line != b.Line {
				return a.Line - b.Line
			}
			return strings.Compare(a.Message, b.Message)
		})
		out = append(out, mine...)
	}
	return out
}

// fits reports whether value, a valid JSON value, fits pr's kind as the
// panel's control sends it.
func fits(pr Prop, value json.RawMessage) bool {
	kind := jsonKind(value)
	if kind == "null" {
		return pr.Kind != KindFunction
	}
	switch pr.Kind {
	case KindString, KindSnippet:
		return kind == "a string"
	case KindNumber:
		return kind == "a number"
	case KindBoolean:
		return kind == "a boolean"
	case KindArray:
		return kind == "an array"
	case KindObject:
		return kind == "an object"
	case KindFunction:
		return false
	case KindUnion:
		var v any
		if kind == "an array" || kind == "an object" || json.Unmarshal(value, &v) != nil {
			return false
		}
		for _, lit := range pr.Literals {
			var l any
			if json.Unmarshal([]byte(lit), &l) == nil && l == v {
				return true
			}
		}
		return false
	default:
		return true
	}
}

// wants is what fits pr's kind, as a phrase: a string, one of "sm", "md".
func wants(pr Prop) string {
	switch pr.Kind {
	case KindString:
		return "a string"
	case KindSnippet:
		return "a string, the snippet's text"
	case KindNumber:
		return "a number"
	case KindBoolean:
		return "true or false"
	case KindArray:
		return "an array"
	case KindObject:
		return "an object"
	case KindUnion:
		return "one of " + strings.Join(pr.Literals, ", ")
	default:
		return "any value"
	}
}

// shown is value as a message shows it: a scalar as written, cut short when
// long, an array or object by its kind.
func shown(value json.RawMessage) string {
	switch kind := jsonKind(value); kind {
	case "an array", "an object":
		return kind
	}
	const most = 40
	if r := []rune(string(value)); len(r) > most {
		return string(r[:most-3]) + "..."
	}
	return string(value)
}

// jsonKind is the kind of value, a valid JSON value, as a phrase: a string,
// a number, a boolean, null, an array, or an object.
func jsonKind(value []byte) string {
	value = bytes.TrimSpace(value)
	if len(value) == 0 {
		return "nothing"
	}
	switch value[0] {
	case '"':
		return "a string"
	case '{':
		return "an object"
	case '[':
		return "an array"
	case 't', 'f':
		return "a boolean"
	case 'n':
		return "null"
	default:
		return "a number"
	}
}
