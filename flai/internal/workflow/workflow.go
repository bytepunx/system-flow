// Package workflow holds the board as data (ADR-0301): a workflow file's
// states, each with a kind that says what it means to flai, the transitions
// between them with who may make each, the rules it passes, and the effects
// flai gives it, and the timers measured over every item's transitions.
// It reads the file and answers lookups on it; it imports nothing of the work
// items it describes.
package workflow

import (
	"errors"
	"fmt"
	"regexp"
	"slices"
	"strings"
	"time"

	"github.com/goccy/go-yaml"
	"github.com/goccy/go-yaml/ast"
	"github.com/goccy/go-yaml/parser"
	"github.com/goccy/go-yaml/token"
)

// File is the workflow file's path inside the design folder.
const File = "workflows/board.md"

// CodeFile is the code flai reports a workflow file that does not parse
// under (ADR-0301 §1).
const CodeFile = "workflow.file"

// State kinds: what a state means to flai (ADR-0301 §2).
const (
	KindQueue     = "queue"
	KindActive    = "active"
	KindWaiting   = "waiting"
	KindAccepted  = "accepted"
	KindCancelled = "cancelled"
)

// Kinds are the state kinds, in ADR-0301 §2's order.
var Kinds = []string{KindQueue, KindActive, KindWaiting, KindAccepted, KindCancelled}

// Item types a state, a transition, or a timer takes (ADR-0003).
const (
	TypeEpic  = "epic"
	TypeStory = "story"
	TypeTask  = "task"
)

// Types are the item types, epic, story, and task.
var Types = []string{TypeEpic, TypeStory, TypeTask}

// Who may make a transition, as its by lists them (ADR-0301 §3).
const (
	ByAnyone     = "anyone"
	ByOperator   = "operator"
	ByAgent      = "agent"
	ByFlai       = "flai"
	ByAcceptance = "acceptance"
	ByCascade    = "cascade"
)

// Bys are the words a transition's by takes.
var Bys = []string{ByAnyone, ByOperator, ByAgent, ByFlai, ByAcceptance, ByCascade}

// Effects: what flai does once, at a move (ADR-0301 §3).
const (
	EffectCommit   = "commit"
	EffectPull     = "pull"
	EffectAccept   = "accept"
	EffectSendBack = "send-back"
	EffectCancel   = "cancel"
	EffectReopen   = "reopen"
)

// Effects are the effects a transition may give.
var Effects = []string{EffectCommit, EffectPull, EffectAccept, EffectSendBack, EffectCancel, EffectReopen}

// Severities a rule entry may set (ADR-0301 §4).
const (
	SeverityRefuse = "refuse"
	SeverityWarn   = "warn"
)

// Severities are the severities a rule entry may set.
var Severities = []string{SeverityRefuse, SeverityWarn}

// Which completed items a timer's aggregates take (ADR-0301 §5).
const (
	CountsAccepted = "accepted"
	CountsClosed   = "closed"
)

// Statuses are the documentation standard's statuses a workflow file takes.
var Statuses = []string{"draft", "active", "deprecated"}

// Workflow is a workflow file read: its states in board order, its
// transitions in file order, and its timers.
type Workflow struct {
	Path        string       `json:"path"`
	Title       string       `json:"title"`
	Updated     string       `json:"updated"`
	Status      string       `json:"status"`
	States      []State      `json:"states"`
	Transitions []Transition `json:"transitions"`
	Timers      []Timer      `json:"timers"`
}

// State is one column of the board.
type State struct {
	Name    string   `yaml:"name" json:"name"`
	Kind    string   `yaml:"kind" json:"kind"`
	Title   string   `yaml:"title" json:"title"`
	Types   []string `yaml:"types" json:"types"`
	Initial bool     `yaml:"initial" json:"initial"`
	Order   bool     `yaml:"order" json:"order"`
	Limit   bool     `yaml:"limit" json:"limit"`
	Aliases []string `yaml:"aliases" json:"aliases"`
	// Line is the file line the state's entry starts on.
	Line int `yaml:"-" json:"line"`
}

// Transition is one move between states.
type Transition struct {
	From    []string  `json:"from"`
	To      string    `json:"to"`
	Types   []string  `json:"types"`
	By      []string  `json:"by"`
	Rules   []RuleRef `json:"rules"`
	Effects []string  `json:"effects"`
	// Line is the file line the transition's entry starts on.
	Line int `json:"line"`
}

// RuleRef is a rule a transition lists, with the severity it sets, empty
// for the rule's default.
type RuleRef struct {
	Rule     string `json:"rule"`
	Severity string `json:"severity,omitempty"`
}

// Point is a moment of an item's history a timer measures from or to:
// created or completed (At), or the first or last transition (Edge) into a
// state or into any state of a kind.
type Point struct {
	At    string `json:"at,omitempty"`
	Edge  string `json:"edge,omitempty"`
	State string `json:"state,omitempty"`
	Kind  string `json:"kind,omitempty"`
}

// The points a timer names by a word, and the edges of the others.
const (
	PointCreated   = "created"
	PointCompleted = "completed"
	EdgeFirst      = "first"
	EdgeLast       = "last"
)

// Timer is a duration measured over every item's transitions: between two
// points, or the sum of the time in some states.
type Timer struct {
	Name   string   `json:"name"`
	Title  string   `json:"title"`
	From   *Point   `json:"from,omitempty"`
	To     *Point   `json:"to,omitempty"`
	Sum    []string `json:"sum"`
	Until  *Point   `json:"until,omitempty"`
	Types  []string `json:"types"`
	Counts string   `json:"counts"`
	// Line is the file line the timer's entry starts on.
	Line int `json:"line"`
}

// Problem is why a workflow file does not parse: the file, the 1-based line,
// the field, such as states[2].kind, and what to do.
type Problem struct {
	Path    string `json:"path"`
	Line    int    `json:"line"`
	Field   string `json:"field,omitempty"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

// String is the problem as "path:line: code: field: message".
func (p Problem) String() string {
	if p.Field == "" {
		return fmt.Sprintf("%s:%d: %s: %s", p.Path, p.Line, p.Code, p.Message)
	}
	return fmt.Sprintf("%s:%d: %s: %s: %s", p.Path, p.Line, p.Code, p.Field, p.Message)
}

// Error is the problem's String.
func (p *Problem) Error() string { return p.String() }

var kebab = regexp.MustCompile(`^[a-z0-9]+(-[a-z0-9]+)*$`)

// Parse reads the workflow file at path rel from its content, with the
// defaults filled in, or the first Problem that keeps it from parsing.
func Parse(rel, content string) (Workflow, *Problem) {
	fm, ok := frontMatter(content)
	if !ok {
		return Workflow{}, &Problem{Path: rel, Line: 1, Code: CodeFile, Message: "the file has no front matter; open it with a --- line and close the front matter with another, as ADR-0301 §1 lists it"}
	}
	// The newline puts each YAML line on its file line: the front matter
	// starts on the file's second.
	file, err := parser.ParseBytes([]byte("\n"+fm), 0)
	if err != nil {
		line, msg, _ := yamlError(err)
		return Workflow{}, &Problem{Path: rel, Line: line, Code: CodeFile, Message: fmt.Sprintf("front matter does not parse: %s; fix the YAML", msg)}
	}
	p := &reader{rel: rel, lines: map[string]int{}, at: map[token.Position]string{}, onLine: map[int]string{}}
	var raw rawFile
	if len(file.Docs) > 0 && file.Docs[0].Body != nil {
		node := file.Docs[0].Body
		if pairs(node) == nil {
			return Workflow{}, &Problem{Path: rel, Line: 2, Code: CodeFile, Message: "front matter is not a mapping of fields; give each field as key: value, as ADR-0301 §1 lists them"}
		}
		p.walk(node, "")
		if err := yaml.NodeToValue(node, &raw, yaml.Strict()); err != nil {
			line, msg, pos := yamlError(err)
			return Workflow{}, &Problem{Path: rel, Line: line, Field: p.fieldAt(pos), Code: CodeFile, Message: fmt.Sprintf("%s; give each field the shape ADR-0301 §1 gives it", msg)}
		}
	}
	w := Workflow{Path: rel, Title: raw.Title, Updated: raw.Updated, Status: raw.Status, States: raw.States, Transitions: raw.Transitions, Timers: raw.Timers}
	if prob := p.validate(&w); prob != nil {
		return Workflow{}, prob
	}
	w.fill()
	return w, nil
}

// frontMatter is the text between a file's opening --- line and the next.
func frontMatter(content string) (string, bool) {
	content = strings.ReplaceAll(content, "\r\n", "\n")
	rest, ok := strings.CutPrefix(content, "---\n")
	if !ok {
		return "", false
	}
	if strings.HasPrefix(rest, "---\n") || rest == "---" {
		return "", true
	}
	if i := strings.Index(rest, "\n---\n"); i >= 0 {
		return rest[:i+1], true
	}
	if strings.HasSuffix(rest, "\n---") {
		return rest[:len(rest)-3], true
	}
	return "", false
}

// rawFile is the front matter as YAML gives it.
type rawFile struct {
	Title       string       `yaml:"title"`
	Updated     string       `yaml:"updated"`
	Status      string       `yaml:"status"`
	States      []State      `yaml:"states"`
	Transitions []Transition `yaml:"transitions"`
	Timers      []Timer      `yaml:"timers"`
}

// UnmarshalYAML reads a state and the line its entry starts on.
func (s *State) UnmarshalYAML(n ast.Node) error {
	type fields State
	if pairs(n) == nil {
		return shapeError(n, "a state is a mapping, such as {name: ready, kind: queue}")
	}
	if err := yaml.NodeToValue(n, (*fields)(s), yaml.Strict()); err != nil {
		return err
	}
	s.Line = n.GetToken().Position.Line
	return nil
}

// UnmarshalYAML reads a transition, its from a state or a list of them.
func (t *Transition) UnmarshalYAML(n ast.Node) error {
	if pairs(n) == nil {
		return shapeError(n, "a transition is a mapping, such as {from: ready, to: in-progress}")
	}
	var f struct {
		From    names     `yaml:"from"`
		To      string    `yaml:"to"`
		Types   []string  `yaml:"types"`
		By      []string  `yaml:"by"`
		Rules   []RuleRef `yaml:"rules"`
		Effects []string  `yaml:"effects"`
	}
	if err := yaml.NodeToValue(n, &f, yaml.Strict()); err != nil {
		return err
	}
	*t = Transition{From: f.From, To: f.To, Types: f.Types, By: f.By, Rules: f.Rules, Effects: f.Effects, Line: n.GetToken().Position.Line}
	return nil
}

// UnmarshalYAML reads a rule entry: the rule's name, or {rule, severity}.
func (r *RuleRef) UnmarshalYAML(n ast.Node) error {
	if pairs(n) == nil {
		if err := yaml.NodeToValue(n, &r.Rule); err != nil || r.Rule == "" {
			return shapeError(n, "a rule entry is a rule's name or {rule: <name>, severity: refuse | warn}")
		}
		return nil
	}
	var f struct {
		Rule     string `yaml:"rule"`
		Severity string `yaml:"severity"`
	}
	if err := yaml.NodeToValue(n, &f, yaml.Strict()); err != nil {
		return err
	}
	if f.Rule == "" {
		return shapeError(n, "a rule entry's mapping needs rule, the rule's name, besides severity")
	}
	*r = RuleRef(f)
	return nil
}

// UnmarshalYAML reads a timer and the line its entry starts on.
func (t *Timer) UnmarshalYAML(n ast.Node) error {
	if pairs(n) == nil {
		return shapeError(n, "a timer is a mapping, such as {name: lead, title: Lead time, from: created, to: completed}")
	}
	var f struct {
		Name   string   `yaml:"name"`
		Title  string   `yaml:"title"`
		From   *Point   `yaml:"from"`
		To     *Point   `yaml:"to"`
		Sum    []string `yaml:"sum"`
		Until  *Point   `yaml:"until"`
		Types  []string `yaml:"types"`
		Counts string   `yaml:"counts"`
	}
	if err := yaml.NodeToValue(n, &f, yaml.Strict()); err != nil {
		return err
	}
	*t = Timer{Name: f.Name, Title: f.Title, From: f.From, To: f.To, Sum: f.Sum, Until: f.Until, Types: f.Types, Counts: f.Counts, Line: n.GetToken().Position.Line}
	return nil
}

const pointShape = "a point is created, completed, {first: <state>}, {last: <state>}, {first: {kind: <kind>}}, or {last: {kind: <kind>}}"

// UnmarshalYAML reads a point: created, completed, {first|last: <state>},
// or {first|last: {kind: <kind>}}.
func (p *Point) UnmarshalYAML(n ast.Node) error {
	ps := pairs(n)
	if ps == nil {
		var s string
		if err := yaml.NodeToValue(n, &s); err != nil || (s != PointCreated && s != PointCompleted) {
			return shapeError(n, pointShape)
		}
		*p = Point{At: s}
		return nil
	}
	if len(ps) != 1 {
		return shapeError(n, pointShape)
	}
	edge := ps[0].Key.GetToken().Value
	if edge != EdgeFirst && edge != EdgeLast {
		return shapeError(n, pointShape)
	}
	v := ps[0].Value
	if pairs(v) == nil {
		var s string
		if err := yaml.NodeToValue(v, &s); err != nil || s == "" {
			return shapeError(n, pointShape)
		}
		*p = Point{Edge: edge, State: s}
		return nil
	}
	var k struct {
		Kind string `yaml:"kind"`
	}
	if err := yaml.NodeToValue(v, &k, yaml.Strict()); err != nil || k.Kind == "" {
		return shapeError(n, pointShape)
	}
	*p = Point{Edge: edge, Kind: k.Kind}
	return nil
}

// String is the point as the file writes it.
func (p Point) String() string {
	switch {
	case p.At != "":
		return p.At
	case p.Kind != "":
		return fmt.Sprintf("{%s: {kind: %s}}", p.Edge, p.Kind)
	}
	return fmt.Sprintf("{%s: %s}", p.Edge, p.State)
}

// names is a state or a list of states.
type names []string

// UnmarshalYAML reads one name or a list of names.
func (l *names) UnmarshalYAML(n ast.Node) error {
	if n.Type() == ast.SequenceType {
		var s []string
		if err := yaml.NodeToValue(n, &s); err != nil {
			return shapeError(n, "from is a state or a list of states")
		}
		*l = s
		return nil
	}
	var s string
	if err := yaml.NodeToValue(n, &s); err != nil {
		return shapeError(n, "from is a state or a list of states")
	}
	*l = names{s}
	return nil
}

// reader keeps where each field of the front matter is, to name it in a
// problem.
type reader struct {
	rel    string
	lines  map[string]int
	at     map[token.Position]string
	onLine map[int]string
}

// walk records the line of every field and entry below n, by its path.
func (p *reader) walk(n ast.Node, prefix string) {
	record := func(path string, pos *token.Position) {
		if pos == nil {
			return
		}
		if _, ok := p.lines[path]; !ok {
			p.lines[path] = pos.Line
		}
		key := token.Position{Line: pos.Line, Column: pos.Column}
		if _, ok := p.at[key]; !ok {
			p.at[key] = path
		}
		if _, ok := p.onLine[pos.Line]; !ok {
			p.onLine[pos.Line] = path
		}
	}
	if ps := pairs(n); ps != nil {
		for _, kv := range ps {
			path := kv.Key.GetToken().Value
			if prefix != "" {
				path = prefix + "." + path
			}
			record(path, kv.Key.GetToken().Position)
			if kv.Value != nil && kv.Value.GetToken() != nil {
				record(path, kv.Value.GetToken().Position)
			}
			p.walk(kv.Value, path)
		}
		return
	}
	if seq, ok := n.(*ast.SequenceNode); ok {
		for i, v := range seq.Values {
			path := fmt.Sprintf("%s[%d]", prefix, i)
			if v.GetToken() != nil {
				record(path, v.GetToken().Position)
			}
			p.walk(v, path)
		}
	}
}

// fieldAt is the field whose key or value starts at pos, else the first
// recorded on its line.
func (p *reader) fieldAt(pos *token.Position) string {
	if pos == nil {
		return ""
	}
	if f, ok := p.at[token.Position{Line: pos.Line, Column: pos.Column}]; ok {
		return f
	}
	return p.onLine[pos.Line]
}

// line is the file line of the field at path, or of the nearest entry or
// field above it that the file has, or 1.
func (p *reader) line(path string) int {
	for path != "" {
		if n, ok := p.lines[path]; ok {
			return n
		}
		i := strings.LastIndexAny(path, ".[")
		if i < 0 {
			break
		}
		path = path[:i]
	}
	return 1
}

// fail is the problem at field, on its line.
func (p *reader) fail(field, format string, args ...any) *Problem {
	return &Problem{Path: p.rel, Line: p.line(field), Field: field, Code: CodeFile, Message: fmt.Sprintf(format, args...)}
}

// validate is the first value of w outside what ADR-0301 §1, §3, and §5
// allow, in file order.
func (p *reader) validate(w *Workflow) *Problem {
	if strings.TrimSpace(w.Title) == "" {
		return p.fail("title", "front matter needs title, the document's title")
	}
	if w.Updated == "" {
		return p.fail("updated", "front matter needs updated, the day of the last meaningful edit, as YYYY-MM-DD")
	}
	if _, err := time.Parse(time.DateOnly, w.Updated); err != nil {
		return p.fail("updated", "updated is %q, not a date; give the day of the last meaningful edit as YYYY-MM-DD", w.Updated)
	}
	if w.Status == "" {
		return p.fail("status", "front matter needs status, one of %s", strings.Join(Statuses, ", "))
	}
	if !slices.Contains(Statuses, w.Status) {
		return p.fail("status", "status %q is not one of %s", w.Status, strings.Join(Statuses, ", "))
	}
	if len(w.States) == 0 {
		return p.fail("states", "front matter needs states, the board's columns in board order, one entry or more (ADR-0301 §2)")
	}
	if prob := p.states(w); prob != nil {
		return prob
	}
	if len(w.Transitions) == 0 {
		return p.fail("transitions", "front matter needs transitions, the moves between states, one entry or more (ADR-0301 §3)")
	}
	if prob := p.transitions(w); prob != nil {
		return prob
	}
	return p.timers(w)
}

func (p *reader) states(w *Workflow) *Problem {
	seen := map[string]string{}
	for i, s := range w.States {
		at := fmt.Sprintf("states[%d]", i)
		if s.Name == "" {
			return p.fail(at+".name", "a state needs name, as status and transitions name it")
		}
		if !kebab.MatchString(s.Name) {
			return p.fail(at+".name", "state name %q is not kebab-case; use lowercase words joined by dashes", s.Name)
		}
		if prev, ok := seen[s.Name]; ok {
			return p.fail(at+".name", "state name %s is also %s; give each state a name of its own", s.Name, prev)
		}
		seen[s.Name] = "the name of " + at
		if s.Kind == "" {
			return p.fail(at+".kind", "state %s needs kind, one of %s (ADR-0301 §2)", s.Name, strings.Join(Kinds, ", "))
		}
		if !slices.Contains(Kinds, s.Kind) {
			return p.fail(at+".kind", "state %s's kind %q is not one of %s", s.Name, s.Kind, strings.Join(Kinds, ", "))
		}
		if prob := p.types(at+".types", s.Types); prob != nil {
			return prob
		}
		if s.Initial && s.Kind != KindQueue {
			return p.fail(at+".initial", "state %s is %s, and only a queue state is initial", s.Name, s.Kind)
		}
		if s.Order && s.Kind != KindQueue {
			return p.fail(at+".order", "state %s is %s, and only a queue state reads in the pull order", s.Name, s.Kind)
		}
	}
	for i, s := range w.States {
		for j, a := range s.Aliases {
			at := fmt.Sprintf("states[%d].aliases[%d]", i, j)
			if !kebab.MatchString(a) {
				return p.fail(at, "alias %q is not kebab-case; use lowercase words joined by dashes", a)
			}
			if prev, ok := seen[a]; ok {
				return p.fail(at, "alias %s of state %s is also %s; a name is never a current state's name, and never two states' alias (ADR-0301 §6)", a, s.Name, prev)
			}
			seen[a] = "an alias of state " + s.Name
		}
	}
	return nil
}

func (p *reader) transitions(w *Workflow) *Problem {
	type move struct{ from, to, typ string }
	seen := map[move]int{}
	for i, t := range w.Transitions {
		at := fmt.Sprintf("transitions[%d]", i)
		if len(t.From) == 0 {
			return p.fail(at+".from", "a transition needs from, the state or states it leaves")
		}
		for j, f := range t.From {
			field := at + ".from"
			if len(t.From) > 1 {
				field = fmt.Sprintf("%s.from[%d]", at, j)
			}
			if prob := p.stateOf(w, field, f); prob != nil {
				return prob
			}
		}
		if t.To == "" {
			return p.fail(at+".to", "a transition needs to, the state it enters")
		}
		if prob := p.stateOf(w, at+".to", t.To); prob != nil {
			return prob
		}
		if prob := p.types(at+".types", t.Types); prob != nil {
			return prob
		}
		for j, b := range t.By {
			if !slices.Contains(Bys, b) {
				return p.fail(fmt.Sprintf("%s.by[%d]", at, j), "by %q is not one of %s (ADR-0301 §3)", b, strings.Join(Bys, ", "))
			}
		}
		for j, r := range t.Rules {
			field := fmt.Sprintf("%s.rules[%d]", at, j)
			if !kebab.MatchString(r.Rule) {
				return p.fail(field, "rule %q is not kebab-case; name a rule flai implements (ADR-0301 §4)", r.Rule)
			}
			if r.Severity != "" && !slices.Contains(Severities, r.Severity) {
				return p.fail(field+".severity", "rule %s's severity %q is not one of %s", r.Rule, r.Severity, strings.Join(Severities, ", "))
			}
		}
		for j, e := range t.Effects {
			if !slices.Contains(Effects, e) {
				return p.fail(fmt.Sprintf("%s.effects[%d]", at, j), "effect %q is not one of %s (ADR-0301 §3)", e, strings.Join(Effects, ", "))
			}
		}
		types := t.Types
		if len(types) == 0 {
			types = Types
		}
		for _, f := range t.From {
			for _, typ := range types {
				m := move{f, t.To, typ}
				if prev, ok := seen[m]; ok {
					return p.fail(at, "transitions[%d] already moves a %s from %s to %s; no two transitions share a from, a to, and a type", prev, typ, f, t.To)
				}
				seen[m] = i
			}
		}
	}
	return nil
}

func (p *reader) timers(w *Workflow) *Problem {
	seen := map[string]int{}
	for i, t := range w.Timers {
		at := fmt.Sprintf("timers[%d]", i)
		if t.Name == "" {
			return p.fail(at+".name", "a timer needs name, its name in flai stats and the charts")
		}
		if !kebab.MatchString(t.Name) {
			return p.fail(at+".name", "timer name %q is not kebab-case; use lowercase words joined by dashes", t.Name)
		}
		if prev, ok := seen[t.Name]; ok {
			return p.fail(at+".name", "timer name %s is also timers[%d]'s; give each timer a name of its own", t.Name, prev)
		}
		seen[t.Name] = i
		if strings.TrimSpace(t.Title) == "" {
			return p.fail(at+".title", "timer %s needs title, its heading on the charts", t.Name)
		}
		points := t.From != nil || t.To != nil
		switch {
		case points && (t.From == nil || t.To == nil):
			missing := "to"
			if t.From == nil {
				missing = "from"
			}
			return p.fail(at+"."+missing, "timer %s gives one of from and to; give both, or sum in place of them (ADR-0301 §5)", t.Name)
		case points && t.Sum != nil:
			return p.fail(at+".sum", "timer %s gives from and to and sum; a timer is between two points or a sum, not both", t.Name)
		case !points && len(t.Sum) == 0:
			return p.fail(at, "timer %s needs from and to, or sum, the states it sums (ADR-0301 §5)", t.Name)
		case t.Until != nil && points:
			return p.fail(at+".until", "timer %s gives until without sum; until bounds a sum only", t.Name)
		}
		for _, pt := range []struct {
			name  string
			point *Point
		}{{"from", t.From}, {"to", t.To}, {"until", t.Until}} {
			if prob := p.point(w, at+"."+pt.name, pt.point); prob != nil {
				return prob
			}
		}
		for j, s := range t.Sum {
			if prob := p.stateOf(w, fmt.Sprintf("%s.sum[%d]", at, j), s); prob != nil {
				return prob
			}
		}
		if prob := p.types(at+".types", t.Types); prob != nil {
			return prob
		}
		if t.Counts != "" && t.Counts != CountsAccepted && t.Counts != CountsClosed {
			return p.fail(at+".counts", "timer %s's counts %q is not %s or %s", t.Name, t.Counts, CountsAccepted, CountsClosed)
		}
	}
	return nil
}

// point is the problem with a point naming a state the file lacks or an
// unknown kind.
func (p *reader) point(w *Workflow, field string, pt *Point) *Problem {
	switch {
	case pt == nil:
		return nil
	case pt.Kind != "" && !slices.Contains(Kinds, pt.Kind):
		return p.fail(field, "kind %q is not one of %s", pt.Kind, strings.Join(Kinds, ", "))
	case pt.State != "":
		return p.stateOf(w, field, pt.State)
	}
	return nil
}

// stateOf is the problem with name when it is no state of w.
func (p *reader) stateOf(w *Workflow, field, name string) *Problem {
	if slices.Contains(w.Names(), name) {
		return nil
	}
	return p.fail(field, "%q is no state of the file; name one of %s", name, strings.Join(w.Names(), ", "))
}

// types is the problem with an item type outside epic, story, and task, or
// named twice.
func (p *reader) types(field string, types []string) *Problem {
	for j, t := range types {
		if !slices.Contains(Types, t) {
			return p.fail(fmt.Sprintf("%s[%d]", field, j), "type %q is not one of %s", t, strings.Join(Types, ", "))
		}
		if slices.Index(types, t) != j {
			return p.fail(fmt.Sprintf("%s[%d]", field, j), "type %s is named twice", t)
		}
	}
	return nil
}

// fill gives the fields their defaults and every nil list an empty one.
func (w *Workflow) fill() {
	for i := range w.States {
		s := &w.States[i]
		if s.Title == "" {
			s.Title = s.Name
		}
		if len(s.Types) == 0 {
			s.Types = slices.Clone(Types)
		}
		if s.Aliases == nil {
			s.Aliases = []string{}
		}
	}
	for i := range w.Transitions {
		t := &w.Transitions[i]
		if len(t.Types) == 0 {
			t.Types = slices.Clone(Types)
		}
		if len(t.By) == 0 {
			t.By = []string{ByAnyone}
		}
		if t.Rules == nil {
			t.Rules = []RuleRef{}
		}
		if t.Effects == nil {
			t.Effects = []string{}
		}
	}
	if w.Timers == nil {
		w.Timers = []Timer{}
	}
	for i := range w.Timers {
		t := &w.Timers[i]
		if len(t.Types) == 0 {
			t.Types = slices.Clone(Types)
		}
		if t.Counts == "" {
			t.Counts = CountsAccepted
		}
		if t.Sum == nil {
			t.Sum = []string{}
		}
	}
}

// Names are the states' names in board order.
func (w Workflow) Names() []string {
	out := make([]string, len(w.States))
	for i, s := range w.States {
		out[i] = s.Name
	}
	return out
}

// State is the state named name, or the one that lists name among its
// aliases.
func (w Workflow) State(name string) (State, bool) {
	for _, s := range w.States {
		if s.Name == name || slices.Contains(s.Aliases, name) {
			return s, true
		}
	}
	return State{}, false
}

// Kind is the kind of the state named name, empty when there is none.
func (w Workflow) Kind(name string) string {
	s, _ := w.State(name)
	return s.Kind
}

// OfKind are the states of a kind, in board order.
func (w Workflow) OfKind(kind string) []State {
	var out []State
	for _, s := range w.States {
		if s.Kind == kind {
			out = append(out, s)
		}
	}
	return out
}

// Initial is the state items are created in.
func (w Workflow) Initial() (State, bool) {
	for _, s := range w.States {
		if s.Initial {
			return s, true
		}
	}
	return State{}, false
}

// canonical is the current name of the state name or an alias names.
func (w Workflow) canonical(name string) string {
	if s, ok := w.State(name); ok {
		return s.Name
	}
	return name
}

// From are the transitions an item of type typ may take out of state from,
// in file order.
func (w Workflow) From(typ, from string) []Transition {
	from = w.canonical(from)
	var out []Transition
	for _, t := range w.Transitions {
		if t.Takes(typ) && slices.Contains(t.From, from) {
			out = append(out, t)
		}
	}
	return out
}

// Find is the transition that moves an item of type typ from one state to
// the other.
func (w Workflow) Find(typ, from, to string) (Transition, bool) {
	to = w.canonical(to)
	for _, t := range w.From(typ, from) {
		if t.To == to {
			return t, true
		}
	}
	return Transition{}, false
}

// WithEffect is the first transition an item of type typ takes that gives
// effect.
func (w Workflow) WithEffect(typ, effect string) (Transition, bool) {
	for _, t := range w.Transitions {
		if t.Takes(typ) && t.HasEffect(effect) {
			return t, true
		}
	}
	return Transition{}, false
}

// Takes reports whether the state takes items of type typ.
func (s State) Takes(typ string) bool { return slices.Contains(s.Types, typ) }

// Takes reports whether an item of type typ may make the transition.
func (t Transition) Takes(typ string) bool { return slices.Contains(t.Types, typ) }

// HasEffect reports whether the transition gives effect.
func (t Transition) HasEffect(effect string) bool { return slices.Contains(t.Effects, effect) }

// Allows reports whether the transition's by lists by.
func (t Transition) Allows(by string) bool { return slices.Contains(t.By, by) }

// Who makes a move, as a Mover names it.
const (
	MoverOperator = "operator"
	MoverAgent    = "agent"
	MoverRole     = "role"
	MoverFlai     = "flai"
)

// Mover is who makes a move: the operator, a story's agent, a strategic
// role, or flai, and whether the move is an acceptance or a parent's
// cancellation taking the item along.
type Mover struct {
	Who        string
	Acceptance bool
	Cascade    bool
}

// Permits is nil when the transition's by lets m make it, else why not and
// what to do.
func (t Transition) Permits(m Mover) error {
	for _, b := range t.By {
		switch {
		case b == ByAnyone,
			b == ByOperator && m.Who == MoverOperator,
			b == ByAgent && m.Who == MoverAgent,
			b == ByFlai && m.Who == MoverFlai,
			b == ByAcceptance && m.Acceptance,
			b == ByCascade && m.Cascade:
			return nil
		}
	}
	move := fmt.Sprintf("from %s to %s", orList(t.From), t.To)
	switch {
	case len(t.By) == 1 && t.By[0] == ByAcceptance:
		return fmt.Errorf("only an acceptance moves %s %s: accept it with flai accept", typeList(t.Types), move)
	case len(t.By) == 1 && t.By[0] == ByCascade:
		return fmt.Errorf("the move %s is made only as a parent's cancellation takes the item along: cancel its parent", move)
	}
	who := make([]string, 0, len(t.By))
	for _, b := range t.By {
		who = append(who, byWords[b])
	}
	return fmt.Errorf("only %s may move %s %s: ask %s to make it", orList(who), typeList(t.Types), move, orList(who))
}

var byWords = map[string]string{
	ByOperator:   "the operator",
	ByAgent:      "the story's agent",
	ByFlai:       "flai serve",
	ByAcceptance: "an acceptance",
	ByCascade:    "a parent's cancellation",
}

// typeList is "an item" for every type, else the types as "an epic or a
// story".
func typeList(types []string) string {
	if len(types) == len(Types) {
		return "an item"
	}
	words := make([]string, len(types))
	for i, t := range types {
		if t == TypeEpic {
			words[i] = "an epic"
		} else {
			words[i] = "a " + t
		}
	}
	return orList(words)
}

// orList joins words as "a", "a or b", or "a, b, or c".
func orList(words []string) string {
	switch len(words) {
	case 0:
		return ""
	case 1:
		return words[0]
	case 2:
		return words[0] + " or " + words[1]
	}
	return strings.Join(words[:len(words)-1], ", ") + ", or " + words[len(words)-1]
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

// yamlError is the line, first message line, and position of a YAML error.
func yamlError(err error) (int, string, *token.Position) {
	var ye yaml.Error
	if errors.As(err, &ye) && ye.GetToken() != nil {
		pos := ye.GetToken().Position
		return pos.Line, firstLine(ye.GetMessage()), pos
	}
	return 1, firstLine(err.Error()), nil
}

func firstLine(s string) string {
	l, _, _ := strings.Cut(s, "\n")
	return strings.TrimSpace(l)
}

// shapeError is a YAML error saying msg at node n.
func shapeError(n ast.Node, msg string) error {
	return &yaml.SyntaxError{Message: msg, Token: n.GetToken()}
}
