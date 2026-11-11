// Package formulas reads and checks a project's formulas (ADR-0238): one
// markdown file per formula under the design folder's workflows/formulas,
// whose front matter names the formula, its nature, its inputs, and its
// steps, each a flai command, a script, or a role with its prompt, and whose
// body says what the formula is for. flai formula run makes a story of one, a
// task a step; flai serve makes one on the schedule its file names (S-0422).
package formulas

import (
	"errors"
	"fmt"
	"math"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"time"

	"github.com/goccy/go-yaml"
	"github.com/goccy/go-yaml/ast"
	"github.com/goccy/go-yaml/parser"

	"github.com/bytepunx/system-flow/flai/internal/cron"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// Folder is the formula files' folder inside the design folder.
const Folder = "workflows/formulas"

// Index is the folder's README.md, which says what a formula file holds and
// is not a formula.
const Index = "README.md"

// CodeFile is the code flai check reports a formula file's problems under.
const CodeFile = "formulas.file"

// DefaultNature is the nature of a formula that names none.
const DefaultNature = "improvement"

// Statuses are a formula file's documentation statuses.
var Statuses = []string{"draft", "active", "deprecated"}

// ScriptsFolder is the folder a script step's script lives under, from the
// project root.
const ScriptsFolder = "scripts/"

// The kinds of thing a step runs.
const (
	KindFlai   = "flai"
	KindScript = "script"
	KindRole   = "role"
)

var (
	kebab     = regexp.MustCompile(`^[a-z][a-z0-9]*(?:-[a-z0-9]+)*$`)
	inputName = regexp.MustCompile(`^[a-z][a-z0-9_-]*$`)
	// use is an input's use, {{name}}, with optional spaces inside the
	// braces; it matches a name that is no input's too, so that Validate can
	// report it.
	use = regexp.MustCompile(`\{\{\s*([A-Za-z0-9_.-]+)\s*\}\}`)
)

// Formula is one formula file: its front matter, where it is, and its body.
type Formula struct {
	// Path is the file's path from the project root, slash separated.
	Path        string  `yaml:"-" json:"path"`
	Title       string  `yaml:"title" json:"title"`
	Updated     string  `yaml:"updated" json:"updated"`
	Status      string  `yaml:"status" json:"status"`
	Name        string  `yaml:"name" json:"name"`
	Description string  `yaml:"description" json:"description"`
	Nature      string  `yaml:"nature" json:"nature"`
	Inputs      []Input `yaml:"inputs" json:"inputs"`
	Steps       []Step  `yaml:"steps" json:"steps"`
	// Schedule is when flai serve makes the formula's instance while the
	// formula host action is on (S-0422): a five-field cron expression in
	// UTC or daily, as planning.schedule is; empty is no schedule.
	Schedule string `yaml:"schedule" json:"schedule,omitempty"`

	// Body is the text after the front matter.
	Body string `yaml:"-" json:"-"`
	// Lines are the file lines of the front matter's top-level keys; a key
	// not there is absent.
	Lines map[string]int `yaml:"-" json:"-"`
}

// Input is one of a formula's inputs: its name, what it is, and its default,
// nil for an input the run must be given.
type Input struct {
	Name        string  `yaml:"name" json:"name"`
	Description string  `yaml:"description" json:"description"`
	Default     *string `yaml:"default" json:"default,omitempty"`
	// Line is the entry's file line.
	Line int `yaml:"-" json:"-"`
}

// Scheduled is the formula's schedule parsed, as planning.schedule is: nil
// when it has none.
func (f Formula) Scheduled() (*cron.Schedule, error) {
	spec := strings.TrimSpace(f.Schedule)
	if spec == "" {
		return nil, nil
	}
	sched, err := cron.Parse(spec)
	if err != nil {
		return nil, err
	}
	return &sched, nil
}

// Required reports whether the input has no default, so a run must give it.
func (in Input) Required() bool { return in.Default == nil }

// Where a step runs (S-0423): the instance's worktree, as every step did
// before, or the project's main checkout, for a step such as flai release
// --pending that tags and pushes the checked-out branch.
const (
	InWorktree = "worktree"
	InMain     = "main"
)

// Ins are the values a step's in may take; empty is InWorktree.
var Ins = []string{InWorktree, InMain}

// Step is one of a formula's steps: its id, its title, what it runs, the
// steps it comes after, and where it runs.
type Step struct {
	ID    string   `yaml:"id" json:"id"`
	Title string   `yaml:"title" json:"title"`
	Run   Run      `yaml:"run" json:"run"`
	After []string `yaml:"after" json:"after"`
	// In is where the step runs: InWorktree, the instance's worktree, when
	// empty, or InMain, the project's main checkout (S-0423).
	In string `yaml:"in" json:"in,omitempty"`
	// Stop are the exits that stop the instance cleanly rather than fail the
	// step, each from 1 to 255, on a flai or script step only (S-0621,
	// ADR-0271 §1); none when it names no stop. UnmarshalYAML reads it.
	Stop []int `yaml:"-" json:"stop,omitempty"`
	// stopShape is set when the step's stop is not a list of integers, which
	// the step checks report.
	stopShape bool
	// Line is the entry's file line.
	Line int `yaml:"-" json:"-"`
	// Lines are the file lines of the entry's keys; a key not there is
	// absent.
	Lines map[string]int `yaml:"-" json:"-"`
}

// Kind is what the step runs: KindFlai, KindScript, or KindRole, or empty
// when its run gives none of them or more than one.
func (s Step) Kind() string { return s.Run.Kind() }

// InMain reports whether the step runs in the project's main checkout rather
// than the instance's worktree (S-0423).
func (s Step) InMain() bool { return s.In == InMain }

// Stops reports whether exit is one of the step's stop exits, which stop the
// instance cleanly (ADR-0271 §1).
func (s Step) Stops(exit int) bool { return exit != 0 && slices.Contains(s.Stop, exit) }

// Stop exits are from MinStop to MaxStop: 0 is the step's success, and an
// exit status is at most 255.
const (
	MinStop = 1
	MaxStop = 255
)

// Run is what a step runs: a flai command line without flai, a script's
// command line, or a role and the prompt it is given.
type Run struct {
	Flai   string `yaml:"flai" json:"flai,omitempty"`
	Script string `yaml:"script" json:"script,omitempty"`
	Role   string `yaml:"role" json:"role,omitempty"`
	Prompt string `yaml:"prompt" json:"prompt,omitempty"`
}

// Kind is KindFlai, KindScript, or KindRole for the one of flai, script, and
// role the run gives, or empty when it gives none or more than one.
func (r Run) Kind() string {
	kinds := r.kinds()
	if len(kinds) != 1 {
		return ""
	}
	return kinds[0]
}

// kinds are the kinds the run gives a value, in the order flai, script, role.
func (r Run) kinds() []string {
	var out []string
	for _, k := range []struct{ kind, value string }{{KindFlai, r.Flai}, {KindScript, r.Script}, {KindRole, r.Role}} {
		if strings.TrimSpace(k.value) != "" {
			out = append(out, k.kind)
		}
	}
	return out
}

// ScriptPath is a script run's first word, the script's path from the
// project root.
func (r Run) ScriptPath() string {
	if words := strings.Fields(r.Script); len(words) > 0 {
		return words[0]
	}
	return ""
}

// texts are the run's texts an input may be used in: the flai command line,
// the script's command line, and the role's prompt.
func (r Run) texts() []string { return []string{r.Flai, r.Script, r.Prompt} }

// Problem is one way a formula file falls short: the file, the 1-based line
// of the field it is about (1 for the file), the check code, and what to do.
type Problem struct {
	Path    string `json:"path"`
	Line    int    `json:"line"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

// String is the problem as "path:line: code: message".
func (p Problem) String() string {
	return fmt.Sprintf("%s:%d: %s: %s", p.Path, p.Line, p.Code, p.Message)
}

// Dir is the formula files' folder for the design folder design, relative to
// the project root, slash separated.
func Dir(design string) string { return path.Join(design, Folder) }

// Load reads every formula file but the index under <root>/<design>/
// workflows/formulas, sorted by path. A missing folder is no formulas. A file
// that does not parse is a Problem, and the others still load.
func Load(root, design string) ([]Formula, []Problem, error) {
	dir := Dir(design)
	entries, err := os.ReadDir(filepath.Join(root, filepath.FromSlash(dir)))
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil, nil
	}
	if err != nil {
		return nil, nil, fmt.Errorf("reading the formula files' folder %s: %w", dir, err)
	}
	var out []Formula
	var problems []Problem
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".md") || e.Name() == Index {
			continue
		}
		rel := path.Join(dir, e.Name())
		data, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(rel)))
		if err != nil {
			return nil, nil, fmt.Errorf("reading formula file %s: %w", rel, err)
		}
		f, p := Parse(rel, string(data))
		if p != nil {
			problems = append(problems, *p)
			continue
		}
		out = append(out, f)
	}
	return out, problems, nil
}

// Parse reads the formula file at path rel from its content, with the
// defaults filled in, or the Problem that keeps it from parsing.
func Parse(rel, content string) (Formula, *Problem) {
	fail := func(line int, format string, args ...any) (Formula, *Problem) {
		return Formula{}, &Problem{Path: rel, Line: line, Code: CodeFile, Message: fmt.Sprintf(format, args...)}
	}
	fm, body, err := workitem.SplitFrontMatter(content)
	if err != nil {
		return fail(1, "%v; a formula file opens with --- and its front matter, as ADR-0238 §2 lists it", err)
	}
	// The newline puts each YAML line on its file line: the front matter
	// starts on the file's second.
	file, err := parser.ParseBytes([]byte("\n"+fm), 0)
	if err != nil {
		line, msg := yamlError(err)
		return fail(line, "front matter does not parse: %s; fix the YAML", msg)
	}
	f := Formula{Lines: map[string]int{}}
	if len(file.Docs) > 0 && file.Docs[0].Body != nil {
		node := file.Docs[0].Body
		if pairs(node) == nil {
			return fail(2, "front matter is not a mapping of fields; give each field as key: value, as ADR-0238 §2 lists them")
		}
		if err := yaml.NodeToValue(node, &f); err != nil {
			line, msg := yamlError(err)
			return fail(line, "front matter does not parse: %s; give inputs as a list of {name, description, default} and steps as a list of {id, title, run, after}", msg)
		}
		for _, kv := range pairs(node) {
			key := kv.Key.GetToken()
			f.Lines[key.Value] = key.Position.Line
		}
	}
	f.Path, f.Body = rel, body
	f.fill()
	return f, nil
}

// fill gives the fields their defaults and every nil list an empty one.
func (f *Formula) fill() {
	if f.Nature == "" {
		f.Nature = DefaultNature
	}
	if f.Inputs == nil {
		f.Inputs = []Input{}
	}
	if f.Steps == nil {
		f.Steps = []Step{}
	}
	for i := range f.Steps {
		if f.Steps[i].After == nil {
			f.Steps[i].After = []string{}
		}
		if f.Steps[i].Lines == nil {
			f.Steps[i].Lines = map[string]int{}
		}
	}
}

// line is the file line of the top-level field key, or 1 when it is absent.
func (f Formula) line(key string) int {
	if n := f.Lines[key]; n > 0 {
		return n
	}
	return 1
}

// line is the file line of the step's key, or the step's own when it is
// absent, or 1 when that is unknown.
func (s Step) line(key string) int {
	if n := s.Lines[key]; n > 0 {
		return n
	}
	return max(s.Line, 1)
}

// UnmarshalYAML reads an inputs entry and the line it starts on.
func (in *Input) UnmarshalYAML(n ast.Node) error {
	if pairs(n) == nil {
		return shapeError(n, "an inputs entry is a mapping: {name: <name>, description: <text>, default: <value>}")
	}
	type fields Input
	if err := yaml.NodeToValue(n, (*fields)(in)); err != nil {
		return err
	}
	in.Line = n.GetToken().Position.Line
	return nil
}

// UnmarshalYAML reads a steps entry, the line it starts on, and its keys'.
func (s *Step) UnmarshalYAML(n ast.Node) error {
	if pairs(n) == nil {
		return shapeError(n, "a steps entry is a mapping: {id: <id>, title: <text>, run: {...}, after: [<id>...], in: worktree or main, stop: [<exit>...]}")
	}
	type fields Step
	if err := yaml.NodeToValue(n, (*fields)(s)); err != nil {
		return err
	}
	s.Line = n.GetToken().Position.Line
	s.Lines = map[string]int{}
	for _, kv := range pairs(n) {
		key := kv.Key.GetToken()
		s.Lines[key.Value] = key.Position.Line
		if key.Value == "stop" {
			s.Stop, s.stopShape = exits(kv.Value)
		}
	}
	return nil
}

// exits are the integers of a stop's sequence node, and whether the node is
// not a list of integers, as stop: [2] is (ADR-0271 §1). An integer too large
// for an int is not one.
func exits(n ast.Node) ([]int, bool) {
	seq, ok := n.(*ast.SequenceNode)
	if !ok {
		return nil, true
	}
	out := make([]int, 0, len(seq.Values))
	for _, v := range seq.Values {
		i, ok := v.(*ast.IntegerNode)
		if !ok {
			return nil, true
		}
		switch n := i.Value.(type) {
		case int64:
			if n < math.MinInt32 || n > math.MaxInt32 {
				return nil, true
			}
			out = append(out, int(n))
		case uint64:
			if n > math.MaxInt32 {
				return nil, true
			}
			out = append(out, int(n))
		default:
			return nil, true
		}
	}
	return out, false
}

// UnmarshalYAML reads a step's run, a mapping.
func (r *Run) UnmarshalYAML(n ast.Node) error {
	if pairs(n) == nil {
		return shapeError(n, "run is a mapping with one of flai: <command line>, script: <command line>, or role: <name> and prompt: <text>")
	}
	type fields Run
	return yaml.NodeToValue(n, (*fields)(r))
}

// pairs are a mapping node's key-value pairs, nil for any other node.
func pairs(n ast.Node) []*ast.MappingValueNode {
	switch m := n.(type) {
	case *ast.MappingNode:
		return m.Values
	case *ast.MappingValueNode:
		return []*ast.MappingValueNode{m}
	}
	return nil
}

// yamlError is the line and first message line of a YAML error.
func yamlError(err error) (int, string) {
	var ye yaml.Error
	if errors.As(err, &ye) && ye.GetToken() != nil {
		return ye.GetToken().Position.Line, firstLine(ye.GetMessage())
	}
	return 1, firstLine(err.Error())
}

func firstLine(s string) string {
	l, _, _ := strings.Cut(s, "\n")
	return strings.TrimSpace(l)
}

// shapeError is a YAML error saying msg at node n.
func shapeError(n ast.Node, msg string) error {
	return &yaml.SyntaxError{Message: msg, Token: n.GetToken()}
}

// Uses are the input names text uses as {{name}}, in the order they first
// appear, each once.
func Uses(text string) []string {
	var out []string
	for _, m := range use.FindAllStringSubmatch(text, -1) {
		if !slices.Contains(out, m[1]) {
			out = append(out, m[1])
		}
	}
	return out
}

// Fill is text with each {{name}} that values names replaced by its value;
// a use of a name values does not hold is left as it is.
func Fill(text string, values map[string]string) string {
	return use.ReplaceAllStringFunc(text, func(m string) string {
		if v, ok := values[use.FindStringSubmatch(m)[1]]; ok {
			return v
		}
		return m
	})
}

// Filled is the formula with its title, its description, its steps' titles,
// and their runs' texts filled with values.
func (f Formula) Filled(values map[string]string) Formula {
	out := f
	out.Title = Fill(f.Title, values)
	out.Description = Fill(f.Description, values)
	out.Steps = make([]Step, len(f.Steps))
	for i, s := range f.Steps {
		s.Title = Fill(s.Title, values)
		s.Run.Flai = Fill(s.Run.Flai, values)
		s.Run.Script = Fill(s.Run.Script, values)
		s.Run.Prompt = Fill(s.Run.Prompt, values)
		s.After = slices.Clone(s.After)
		s.Stop = slices.Clone(s.Stop)
		out.Steps[i] = s
	}
	return out
}

// Resolve is every input's value for a run given the values given: each
// given one, and the default of each not given. It refuses a given input the
// formula does not declare, and a required one not given.
func (f Formula) Resolve(given map[string]string) (map[string]string, error) {
	values := map[string]string{}
	var missing []string
	for _, in := range f.Inputs {
		switch v, ok := given[in.Name]; {
		case ok:
			values[in.Name] = v
		case in.Default != nil:
			values[in.Name] = *in.Default
		default:
			missing = append(missing, in.Name)
		}
	}
	var undeclared []string
	for name := range given {
		if !slices.ContainsFunc(f.Inputs, func(in Input) bool { return in.Name == name }) {
			undeclared = append(undeclared, name)
		}
	}
	slices.Sort(undeclared)
	if len(missing) == 0 && len(undeclared) == 0 {
		return values, nil
	}
	var says []string
	if len(undeclared) > 0 {
		says = append(says, fmt.Sprintf("formula %s declares no input %s", f.Name, strings.Join(undeclared, ", ")))
	}
	if len(missing) > 0 {
		says = append(says, fmt.Sprintf("formula %s needs input %s, which has no default", f.Name, strings.Join(missing, ", ")))
	}
	return nil, fmt.Errorf("%s; its inputs are %s: give each as --input <name>=<value>", strings.Join(says, ", and "), f.declared())
}

// declared lists the formula's inputs for a message: each name, marked
// "(required)" when it has no default, or "none".
func (f Formula) declared() string {
	if len(f.Inputs) == 0 {
		return "none"
	}
	names := make([]string, len(f.Inputs))
	for i, in := range f.Inputs {
		names[i] = in.Name
		if in.Required() {
			names[i] += " (required)"
		}
	}
	return strings.Join(names, ", ")
}

// Find is the formula named name among formulas, or an error that lists
// their names.
func Find(formulas []Formula, name string) (Formula, error) {
	names := make([]string, 0, len(formulas))
	for _, f := range formulas {
		if f.Name == name {
			return f, nil
		}
		names = append(names, f.Name)
	}
	if len(names) == 0 {
		return Formula{}, fmt.Errorf("no formula is named %s, and the project has none: add one as %s/%s.md under the design folder", name, Folder, name)
	}
	return Formula{}, fmt.Errorf("no formula is named %s; the formulas are %s (flai formula list)", name, strings.Join(names, ", "))
}

// AnalysisFormulaErrors are the ways analysis.formula naming name falls short
// of the formulas a project's files define (S-0424, S-0692), each a sentence
// beginning with the key: the formula is one a file defines, and, since
// analysis.schedule makes its instance with no inputs, its every input carries
// a default. list and problems are what Load gave for the design folder
// design. None for an empty name, nor for a formula whose file does not
// parse, which Load's problem on it reports. flai check reports each as
// manifest.analysis, and flai manifest set refuses each.
func AnalysisFormulaErrors(list []Formula, problems []Problem, design, name string) []string {
	if name == "" {
		return nil
	}
	file := path.Join(Dir(design), name+".md")
	for _, p := range problems {
		if p.Path == file {
			return nil
		}
	}
	f, err := Find(list, name)
	if err != nil {
		return []string{"analysis.formula: " + err.Error()}
	}
	var out []string
	for _, in := range f.Inputs {
		if in.Required() {
			out = append(out, fmt.Sprintf("analysis.formula names %s, whose input %s has no default, but analysis.schedule makes its instance with no inputs; give %s a default in %s, or name another formula", name, in.Name, in.Name, f.Path))
		}
	}
	return out
}

// Layers are the formula's step ids in layers: first the steps after no
// step, then each layer the steps whose after are all in earlier ones, ids in
// file order within a layer, as a story's plan layers its tasks. An after
// naming no step is passed over; a step on a cycle, or after one, is in no
// layer, and nor is a second step with an id already seen.
func Layers(f Formula) [][]string {
	const visiting, cyclic = -1, -2
	byID := index(f)
	depth := map[string]int{}
	var walk func(s Step) int
	walk = func(s Step) int {
		if d, ok := depth[s.ID]; ok {
			if d == visiting {
				return cyclic
			}
			return d
		}
		depth[s.ID] = visiting
		d := 0
		for _, a := range s.After {
			other, ok := byID[a]
			if !ok {
				continue
			}
			od := walk(other)
			if od == cyclic {
				d = cyclic
				break
			}
			d = max(d, od+1)
		}
		depth[s.ID] = d
		return d
	}
	out := [][]string{}
	placed := map[string]bool{}
	for _, s := range f.Steps {
		if placed[s.ID] {
			continue
		}
		placed[s.ID] = true
		d := walk(byID[s.ID])
		if d < 0 {
			continue
		}
		for len(out) <= d {
			out = append(out, nil)
		}
		out[d] = append(out[d], s.ID)
	}
	return out
}

// index is the formula's steps by id, the first of each id.
func index(f Formula) map[string]Step {
	byID := map[string]Step{}
	for _, s := range f.Steps {
		if _, ok := byID[s.ID]; !ok {
			byID[s.ID] = s
		}
	}
	return byID
}

// Names are what Validate checks a step's run against.
type Names struct {
	// Roles are the names of the roles the project's role files define, as
	// roles.Load reads them; none is no roles.
	Roles []string
	// Script reports whether a script exists at a path from the project
	// root; nil skips the check.
	Script func(path string) bool
}

// ScriptOnDisk reports whether a script exists as a file at a path from the
// project root root, for Names.Script.
func ScriptOnDisk(root string) func(path string) bool {
	return func(p string) bool {
		info, err := os.Stat(filepath.Join(root, filepath.FromSlash(p)))
		return err == nil && !info.IsDir()
	}
}

// Validate checks formulas against ADR-0238 and the names in n, and returns their problems sorted by path and line.
func Validate(formulas []Formula, n Names) []Problem {
	c := &checker{n: n}
	for _, f := range formulas {
		c.formula(f)
	}
	slices.SortStableFunc(c.out, func(a, b Problem) int {
		if a.Path != b.Path {
			return strings.Compare(a.Path, b.Path)
		}
		return a.Line - b.Line
	})
	return c.out
}

type checker struct {
	n   Names
	out []Problem
}

func (c *checker) add(f Formula, line int, format string, args ...any) {
	c.out = append(c.out, Problem{Path: f.Path, Line: line, Code: CodeFile, Message: fmt.Sprintf(format, args...)})
}

func (c *checker) formula(f Formula) {
	file := func(key, format string, args ...any) { c.add(f, f.line(key), format, args...) }
	for _, k := range []struct{ key, value string }{{"title", f.Title}, {"updated", f.Updated}, {"status", f.Status}, {"name", f.Name}, {"description", f.Description}} {
		if strings.TrimSpace(k.value) == "" {
			file(k.key, "front matter needs %s; add it (ADR-0238 §2)", k.key)
		}
	}
	if f.Updated != "" {
		if _, err := time.Parse(time.DateOnly, f.Updated); err != nil {
			file("updated", "updated is %q, not a date; give the day of the last meaningful edit as YYYY-MM-DD", f.Updated)
		}
	}
	if f.Status != "" && !slices.Contains(Statuses, f.Status) {
		file("status", "status is %q, which is none of %s", f.Status, strings.Join(Statuses, ", "))
	}
	stem := strings.TrimSuffix(path.Base(f.Path), ".md")
	switch {
	case f.Name == "":
	case !kebab.MatchString(f.Name):
		file("name", "name %q is not kebab-case; use lowercase words joined by dashes", f.Name)
	case f.Name != stem:
		file("name", "name is %s but the file is %s.md; make them the same", f.Name, stem)
	}
	if !slices.Contains(workitem.Natures, f.Nature) {
		file("nature", "nature is %q, which is none of %s", f.Nature, strings.Join(workitem.Natures, ", "))
	}
	declared := c.inputs(f)
	c.schedule(f)
	c.uses(f, f.line("title"), "title", f.Title, declared)
	c.uses(f, f.line("description"), "description", f.Description, declared)
	if len(f.Steps) == 0 {
		file("steps", "front matter needs steps, a list of {id, title, run, after}; give the formula at least one (ADR-0238 §2)")
	}
	ids := map[string]bool{}
	for _, s := range f.Steps {
		if s.ID != "" && ids[s.ID] {
			c.add(f, s.line("id"), "step id %s is also an earlier step's; give each step its own id", s.ID)
		}
		ids[s.ID] = true
	}
	for _, s := range f.Steps {
		c.step(f, s, ids, declared)
	}
	for _, cycle := range cycles(f) {
		first := index(f)[cycle[0]]
		c.add(f, first.line("after"), "steps %s are after each other in a cycle: %s; remove one of those after entries so that the steps have an order", strings.Join(members(cycle), ", "), strings.Join(cycle, " after "))
	}
}

// inputs checks the formula's inputs and returns the names it declares.
func (c *checker) inputs(f Formula) map[string]bool {
	declared := map[string]bool{}
	for _, in := range f.Inputs {
		switch {
		case in.Name == "":
			c.add(f, in.Line, "inputs entry needs name; give it one such as story-title")
		case !inputName.MatchString(in.Name):
			c.add(f, in.Line, "input name %q is not a lowercase word; start it with a letter and use a-z, 0-9, _, and -", in.Name)
		case declared[in.Name]:
			c.add(f, in.Line, "input %s is declared twice; keep one", in.Name)
		}
		if in.Name != "" && strings.TrimSpace(in.Description) == "" {
			c.add(f, in.Line, "input %s needs description: what the value is", in.Name)
		}
		declared[in.Name] = true
	}
	return declared
}

// schedule checks the formula's schedule: one that parses, on a formula whose
// every input carries a default, since a run its schedule starts is given no
// inputs and takes each input's default (S-0424); each input with none is
// named.
func (c *checker) schedule(f Formula) {
	if strings.TrimSpace(f.Schedule) == "" {
		return
	}
	if _, err := f.Scheduled(); err != nil {
		c.add(f, f.line("schedule"), "%v", err)
	}
	var names []string
	for _, in := range f.Inputs {
		if in.Required() {
			names = append(names, in.Name)
		}
	}
	if len(names) > 0 {
		has, it := "has", "it"
		if len(names) > 1 {
			has, it = "have", "each"
		}
		c.add(f, f.line("schedule"), "schedule starts the formula with no inputs, as nobody is there to give them, but %s %s no default; give %s a default, or remove the schedule", strings.Join(names, ", "), has, it)
	}
}

// uses checks that text, the formula's field what on line line, uses only
// inputs the formula declares.
func (c *checker) uses(f Formula, line int, what, text string, declared map[string]bool) {
	for _, name := range Uses(text) {
		if !declared[name] {
			c.add(f, line, "%s uses {{%s}}, which the formula's inputs do not declare; declare it under inputs or remove the use", what, name)
		}
	}
}

// step checks one step against the formula's step ids and inputs.
func (c *checker) step(f Formula, s Step, ids, declared map[string]bool) {
	at := func(key, format string, args ...any) { c.add(f, s.line(key), format, args...) }
	name := s.ID
	switch {
	case s.ID == "":
		name = fmt.Sprintf("on line %d", s.Line)
		at("id", "step needs id; give it a kebab-case id such as open-pr")
	case !kebab.MatchString(s.ID):
		at("id", "step id %q is not kebab-case; use lowercase words joined by dashes", s.ID)
	}
	if strings.TrimSpace(s.Title) == "" {
		at("title", "step %s needs title; say what the step does", name)
	}
	c.uses(f, s.line("title"), "step "+name+"'s title", s.Title, declared)
	c.run(f, s, name)
	for _, text := range s.Run.texts() {
		c.uses(f, s.line("run"), "step "+name+"'s run", text, declared)
	}
	for _, a := range s.After {
		if !ids[a] || a == "" {
			at("after", "step %s is after %q, which is no step's id; name a step of this formula", name, a)
		}
	}
	if s.In != "" && !slices.Contains(Ins, s.In) {
		at("in", "step %s's in is %q, which is none of %s; give %s, the instance's worktree and the default, or %s, the project's main checkout", name, s.In, strings.Join(Ins, ", "), InWorktree, InMain)
	}
	c.stop(f, s, name)
}

// stop checks a step's stop (ADR-0271 §1): a list of integers, each from
// MinStop to MaxStop and none twice, on a flai or script step.
func (c *checker) stop(f Formula, s Step, name string) {
	if s.Lines["stop"] == 0 {
		return
	}
	at := func(format string, args ...any) { c.add(f, s.line("stop"), format, args...) }
	if s.stopShape {
		at("step %s's stop is not a list of exit codes; give the exits that stop the run cleanly, such as stop: [2]", name)
		return
	}
	if s.Kind() == KindRole {
		at("step %s names stop, but a role step ends worked, asked, or failed with no exit code; give stop to a flai or script step only", name)
		return
	}
	seen := map[int]bool{}
	for _, e := range s.Stop {
		switch {
		case e < MinStop || e > MaxStop:
			at("step %s's stop lists exit %d; a stop exit is from %d to %d, as 0 is the step's success", name, e, MinStop, MaxStop)
		case seen[e]:
			at("step %s's stop lists exit %d twice; keep one", name, e)
		}
		seen[e] = true
	}
}

// run checks a step's run: exactly one kind, and that kind's fields.
func (c *checker) run(f Formula, s Step, name string) {
	at := func(format string, args ...any) { c.add(f, s.line("run"), format, args...) }
	r := s.Run
	if s.Lines["run"] == 0 {
		at("step %s needs run: one of flai: <command line>, script: <command line>, or role: <name> and prompt: <text> (ADR-0238 §3)", name)
		return
	}
	kinds := r.kinds()
	switch len(kinds) {
	case 0:
		at("step %s's run gives none of flai, script, or role a value; give exactly one", name)
		return
	case 1:
	default:
		at("step %s's run gives %s; give exactly one of them, and make the others steps of their own", name, strings.Join(kinds, " and "))
		return
	}
	switch kinds[0] {
	case KindScript:
		p := r.ScriptPath()
		clean := path.Clean(p)
		if path.IsAbs(p) || !strings.HasPrefix(clean, ScriptsFolder) || !filepath.IsLocal(filepath.FromSlash(clean)) {
			at("step %s's script %s is not under %s; a script step runs a project script, such as scripts/check.sh", name, p, ScriptsFolder)
			return
		}
		if c.n.Script != nil && !c.n.Script(clean) {
			at("step %s's script %s does not exist; add it, or name a script that does", name, clean)
		}
	case KindRole:
		switch {
		case slices.Contains(unstartedBuiltins, r.Role):
			at("step %s's role %s is a built-in role a role step cannot start; only the analyzer is started for a role step, through flai analyze's start", name, r.Role)
		case !slices.Contains(c.n.Roles, r.Role):
			at("step %s's role %s is no role file's; name one of %s", name, r.Role, roleList(c.n.Roles))
		}
		if strings.TrimSpace(r.Prompt) == "" {
			at("step %s's run names role %s but no prompt; give prompt: what the role is asked to do", name, r.Role)
		}
	}
}

// unstartedBuiltins are the built-in roles a role step cannot start: each
// keeps its own start, which no formula instance's agent may ask for. The
// analyzer, the fourth, is started for a role step through its own start
// (S-0424).
var unstartedBuiltins = []string{workitem.ActivityPlanner, workitem.ActivityOrchestrator, workitem.ActivityCurator}

// roleList is the role names for a message, or a note that there are none.
func roleList(names []string) string {
	if len(names) == 0 {
		return "the roles a file under workflows/roles defines, and the project has none"
	}
	return strings.Join(names, ", ")
}

// cycles are the cycles among the formula's steps' after, each as the path
// from a step back to itself, such as [a b c a] for a after b after c after
// a, each set of steps once, in the order a walk in file order meets them.
func cycles(f Formula) [][]string {
	const (
		unseen = iota
		open
		closed
	)
	byID := index(f)
	state := map[string]int{}
	var stack []string
	seen := map[string]bool{}
	var out [][]string
	var visit func(id string)
	visit = func(id string) {
		state[id] = open
		stack = append(stack, id)
		for _, a := range byID[id].After {
			if _, ok := byID[a]; !ok {
				continue
			}
			switch state[a] {
			case unseen:
				visit(a)
			case open:
				cycle := append(slices.Clone(stack[slices.Index(stack, a):]), a)
				key := strings.Join(members(cycle), " ")
				if !seen[key] {
					seen[key] = true
					out = append(out, cycle)
				}
			}
		}
		stack = stack[:len(stack)-1]
		state[id] = closed
	}
	for _, s := range f.Steps {
		if state[s.ID] == unseen {
			visit(s.ID)
		}
	}
	return out
}

// members are a cycle's steps, each once, sorted.
func members(cycle []string) []string {
	out := slices.Clone(cycle[:len(cycle)-1])
	slices.Sort(out)
	return out
}
