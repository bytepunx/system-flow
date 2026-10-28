// Package roles reads and checks a project's strategic roles (ADR-0207):
// one markdown file per role under the design folder's workflows/roles,
// whose front matter names the role's harness definition, host action,
// triggers, items, topics, and the calls flai guard allows it, and whose body
// is its definition's prompt.
package roles

import (
	"errors"
	"fmt"
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

	"github.com/bytepunx/system-flow/flai/internal/conventions"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// Folder is the role files' folder inside the design folder.
const Folder = "workflows/roles"

// Index is the folder's README.md, which says what a role file holds and is
// not a role.
const Index = "README.md"

// Marker separates a role file's baseline text from the project's additions,
// as in a convention (ADR-0207 §4).
const Marker = conventions.Marker

// The codes flai check reports a role file's problems under (ADR-0207).
const (
	CodeFile       = "roles.file"
	CodeGrant      = "roles.grant"
	CodeDefinition = "roles.definition"
)

// The trigger kinds a trigger entry takes (ADR-0207 §3).
const (
	TriggerAsked      = "asked"
	TriggerContinuous = "continuous"
	TriggerEvent      = "event"
	TriggerSchedule   = "schedule"
)

// Orchestrator is the orchestrator's role name: the one file whose
// permissions name orchestration.permissions.<p> (ADR-0207 §5).
const Orchestrator = "orchestrator"

var (
	// Statuses are a role file's documentation statuses.
	Statuses = []string{"draft", "active", "deprecated"}
	// Charges are how a role's usage is charged.
	Charges = []string{"item", "named", "issues", "none"}
	// ItemKinds are the item kinds a run takes.
	ItemKinds = []string{"epic", "story"}
	// Events are what an event trigger names.
	Events = []string{"accepted", "cancelled", "ready", "review", "edited", "reordered", "thread-opened"}
	// Schedules are the manifest settings a schedule trigger names besides
	// roles.<name>.schedule.
	Schedules = []string{"planning.schedule", "analysis.schedule", "digest.schedule"}
	// BuiltinActions are the built-in roles' host actions, each only its
	// role's, by role name.
	BuiltinActions = map[string]string{"plan": "planner", "orchestrate": Orchestrator, "analyze": "analyzer", "curate": "curator"}
	// ReservedActions are host actions that are no role's.
	ReservedActions = []string{"push", "forge", "auto-publish", "auto-approve", "auto-upgrade", "agent", "dashboard", "checks", "host", "settings", "components", "merge-queue"}
	// ReservedNames are the sub-agents' definitions, which no role takes.
	ReservedNames = []string{"explorer", "verifier"}
	// ReservedRoles are the verbs of the story's agent and its sub-agents.
	ReservedRoles = []string{conventions.RoleStory, conventions.RoleExplore, conventions.RoleVerify}
	// EditTools are the harness tools that edit a file, which a definition
	// lists only when files names a folder.
	EditTools = []string{"Edit", "MultiEdit", "Write", "NotebookEdit"}
)

// MCPPrefix starts a flai MCP tool's name in a definition's tools.
const MCPPrefix = "mcp__flai__"

var kebab = regexp.MustCompile(`^[a-z][a-z0-9]*(?:-[a-z0-9]+)*$`)

// Role is one role file: its front matter, where it is, and its body.
type Role struct {
	// Path is the file's path from the project root, slash separated.
	Path        string       `yaml:"-" json:"path"`
	Title       string       `yaml:"title" json:"title"`
	Updated     string       `yaml:"updated" json:"updated"`
	Status      string       `yaml:"status" json:"status"`
	Name        string       `yaml:"name" json:"name"`
	Role        string       `yaml:"role" json:"role"`
	Definition  Definition   `yaml:"definition" json:"definition"`
	Action      string       `yaml:"action" json:"action"`
	Trigger     []Trigger    `yaml:"trigger" json:"trigger"`
	Item        Kinds        `yaml:"item" json:"item"`
	Topics      []string     `yaml:"topics" json:"topics"`
	WorksStory  bool         `yaml:"works_story" json:"works_story"`
	Charge      string       `yaml:"charge" json:"charge"`
	Tools       []Call       `yaml:"tools" json:"tools"`
	Commands    []Call       `yaml:"commands" json:"commands"`
	Files       []string     `yaml:"files" json:"files"`
	Permissions []Permission `yaml:"permissions" json:"permissions"`

	// Body is the text after the front matter.
	Body string `yaml:"-" json:"-"`
	// Baseline is the body above the marker, the whole body when it has none.
	Baseline string `yaml:"-" json:"-"`
	// Additions is the body from the marker on, empty when it has none.
	Additions string `yaml:"-" json:"-"`
	// Lines are the file lines of the front matter's top-level keys, and of
	// definition's as definition.<key>; a key not there is absent.
	Lines map[string]int `yaml:"-" json:"-"`
}

// Definition is what flai generates the role's harness definition from
// (ADR-0207 §4).
type Definition struct {
	Description string   `yaml:"description" json:"description"`
	Tools       []string `yaml:"tools" json:"tools"`
	Model       string   `yaml:"model" json:"model"`
}

// Call is a call entry: the call's name, plainly or under a condition.
type Call struct {
	Call string `json:"call"`
	When string `json:"when,omitempty"`
	// Line is the entry's file line.
	Line int `json:"-"`
}

// Trigger is a trigger entry: asked, continuous, an event, or a schedule,
// the last two under an optional condition (ADR-0207 §3).
type Trigger struct {
	Kind     string `json:"kind"`
	Event    string `json:"event,omitempty"`
	Schedule string `json:"schedule,omitempty"`
	When     string `json:"when,omitempty"`
	// Line is the entry's file line.
	Line int `json:"-"`
}

// Permission is a permissions entry: the operator's setting and the calls
// it allows (ADR-0207 §5).
type Permission struct {
	Setting  string `yaml:"setting" json:"setting"`
	Tools    []Call `yaml:"tools" json:"tools"`
	Commands []Call `yaml:"commands" json:"commands"`
	// Line is the entry's file line.
	Line int `yaml:"-" json:"-"`
}

// Kinds are the item kinds a run takes; item: none is empty.
type Kinds []string

// Problem is one way a role file falls short: the file, the 1-based line of
// the field it is about (1 for the file), the check code, and what to do.
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

// Vocabulary is what Validate checks names against, which the caller takes
// from flai's MCP server and command tree. A nil Tools or Commands skips the
// check that an entry names a known MCP tool or command; a nil Reads skips
// the check that a definition's flai MCP tools are ones its rules grant.
type Vocabulary struct {
	// Tools are every flai MCP tool's name.
	Tools []string
	// Commands are every flai command path without "flai", such as "edit"
	// and "story new".
	Commands []string
	// Reads are the MCP tools every agent may call.
	Reads []string
}

// Dir is the role files' folder for the design folder design, relative to
// the project root, slash separated.
func Dir(design string) string { return path.Join(design, Folder) }

// Load reads every role file but the index under <root>/<design>/workflows/
// roles, sorted by path. A missing folder is no roles. A file that does not
// parse is a Problem, and the others still load.
func Load(root, design string) ([]Role, []Problem, error) {
	dir := Dir(design)
	entries, err := os.ReadDir(filepath.Join(root, filepath.FromSlash(dir)))
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil, nil
	}
	if err != nil {
		return nil, nil, fmt.Errorf("reading the role files' folder %s: %w", dir, err)
	}
	var out []Role
	var problems []Problem
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".md") || e.Name() == Index {
			continue
		}
		rel := path.Join(dir, e.Name())
		data, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(rel)))
		if err != nil {
			return nil, nil, fmt.Errorf("reading role file %s: %w", rel, err)
		}
		r, p := Parse(rel, string(data))
		if p != nil {
			problems = append(problems, *p)
			continue
		}
		out = append(out, r)
	}
	return out, problems, nil
}

// Parse reads the role file at path rel from its content, with the defaults
// filled in, or the Problem that keeps it from parsing.
func Parse(rel, content string) (Role, *Problem) {
	fail := func(line int, format string, args ...any) (Role, *Problem) {
		return Role{}, &Problem{Path: rel, Line: line, Code: CodeFile, Message: fmt.Sprintf(format, args...)}
	}
	fm, body, err := workitem.SplitFrontMatter(content)
	if err != nil {
		return fail(1, "%v; a role file opens with --- and its front matter, as ADR-0207 §2 lists it", err)
	}
	// The newline puts each YAML line on its file line: the front matter
	// starts on the file's second.
	file, err := parser.ParseBytes([]byte("\n"+fm), 0)
	if err != nil {
		line, msg := yamlError(err)
		return fail(line, "front matter does not parse: %s; fix the YAML", msg)
	}
	r := Role{Lines: map[string]int{}}
	if len(file.Docs) > 0 && file.Docs[0].Body != nil {
		node := file.Docs[0].Body
		if pairs(node) == nil {
			return fail(2, "front matter is not a mapping of fields; give each field as key: value, as ADR-0207 §2 lists them")
		}
		if err := yaml.NodeToValue(node, &r, yaml.Strict()); err != nil {
			line, msg := yamlError(err)
			return fail(line, "front matter does not parse: %s; give each field the shape ADR-0207 §2 gives it", msg)
		}
		for _, kv := range pairs(node) {
			key := kv.Key.GetToken()
			r.Lines[key.Value] = key.Position.Line
			if key.Value == "definition" {
				for _, sub := range pairs(kv.Value) {
					t := sub.Key.GetToken()
					r.Lines["definition."+t.Value] = t.Position.Line
				}
			}
		}
	}
	r.Path, r.Body, r.Baseline = rel, body, body
	if i := strings.Index(body, Marker); i >= 0 {
		r.Baseline, r.Additions = body[:i], body[i:]
	}
	r.fill()
	return r, nil
}

// fill gives the fields their defaults and every nil list an empty one.
func (r *Role) fill() {
	if r.Definition.Model == "" {
		r.Definition.Model = "inherit"
	}
	if r.Charge == "" {
		r.Charge = "none"
	}
	if r.Item == nil {
		r.Item = Kinds{}
	}
	for _, l := range []*[]string{&r.Topics, &r.Files, &r.Definition.Tools} {
		if *l == nil {
			*l = []string{}
		}
	}
	for _, l := range []*[]Call{&r.Tools, &r.Commands} {
		if *l == nil {
			*l = []Call{}
		}
	}
	if r.Trigger == nil {
		r.Trigger = []Trigger{}
	}
	if r.Permissions == nil {
		r.Permissions = []Permission{}
	}
	for i := range r.Permissions {
		p := &r.Permissions[i]
		if p.Tools == nil {
			p.Tools = []Call{}
		}
		if p.Commands == nil {
			p.Commands = []Call{}
		}
	}
}

// line is the file line of the top-level field key, or 1 when it is absent.
func (r Role) line(key string) int {
	if n := r.Lines[key]; n > 0 {
		return n
	}
	return 1
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

// UnmarshalYAML reads a call entry: the call's name, or {call, when}.
func (c *Call) UnmarshalYAML(n ast.Node) error {
	c.Line = n.GetToken().Position.Line
	if pairs(n) == nil {
		if err := yaml.NodeToValue(n, &c.Call); err != nil || c.Call == "" {
			return shapeError(n, "a call entry is a call's name or {call: <name>, when: <condition>}")
		}
		return nil
	}
	var f struct {
		Call string `yaml:"call"`
		When string `yaml:"when"`
	}
	if err := yaml.NodeToValue(n, &f, yaml.Strict()); err != nil {
		return err
	}
	if f.Call == "" {
		return shapeError(n, "a call entry's mapping needs call, the call's name, besides when")
	}
	c.Call, c.When = f.Call, f.When
	return nil
}

// UnmarshalYAML reads a trigger entry: asked, continuous, {event, when}, or
// {schedule, when}.
func (t *Trigger) UnmarshalYAML(n ast.Node) error {
	t.Line = n.GetToken().Position.Line
	if pairs(n) == nil {
		if err := yaml.NodeToValue(n, &t.Kind); err != nil || t.Kind == "" {
			return shapeError(n, "a trigger entry is asked, continuous, {event: <event>, when: <condition>}, or {schedule: <setting>, when: <condition>}")
		}
		return nil
	}
	var f struct {
		Event    string `yaml:"event"`
		Schedule string `yaml:"schedule"`
		When     string `yaml:"when"`
	}
	if err := yaml.NodeToValue(n, &f, yaml.Strict()); err != nil {
		return err
	}
	switch {
	case f.Event != "" && f.Schedule != "":
		return shapeError(n, "a trigger entry names an event or a schedule, not both; make them two entries")
	case f.Event != "":
		t.Kind, t.Event = TriggerEvent, f.Event
	case f.Schedule != "":
		t.Kind, t.Schedule = TriggerSchedule, f.Schedule
	default:
		return shapeError(n, "a trigger entry's mapping needs event or schedule")
	}
	t.When = f.When
	return nil
}

// UnmarshalYAML reads a permissions entry and the line it starts on.
func (p *Permission) UnmarshalYAML(n ast.Node) error {
	type fields Permission
	if err := yaml.NodeToValue(n, (*fields)(p), yaml.Strict()); err != nil {
		return err
	}
	p.Line = n.GetToken().Position.Line
	return nil
}

// UnmarshalYAML reads item: none as no kinds, a kind as one, or a list.
func (k *Kinds) UnmarshalYAML(n ast.Node) error {
	if n.Type() == ast.SequenceType {
		var l []string
		if err := yaml.NodeToValue(n, &l); err != nil {
			return shapeError(n, "item is none or a list of epic and story")
		}
		*k = l
		return nil
	}
	var s string
	if err := yaml.NodeToValue(n, &s); err != nil {
		return shapeError(n, "item is none or a list of epic and story")
	}
	if s == "none" {
		*k = Kinds{}
	} else {
		*k = Kinds{s}
	}
	return nil
}

// Validate checks roles, read from the design folder design, against
// ADR-0207 and the names in v, and returns their problems sorted by path
// and line.
func Validate(roles []Role, design string, v Vocabulary) []Problem {
	c := &checker{design: strings.TrimSuffix(path.Clean(design), "/"), v: v}
	names, verbs := map[string]string{}, map[string]string{}
	for _, r := range roles {
		c.role(r)
		if prev, ok := names[r.Name]; ok && r.Name != "" {
			c.add(r, CodeFile, r.line("name"), "name %s is also %s's; give each role its own name", r.Name, prev)
		} else {
			names[r.Name] = r.Path
		}
		if prev, ok := verbs[r.Role]; ok && r.Role != "" {
			c.add(r, CodeFile, r.line("role"), "role %s is also %s's; give each role its own verb", r.Role, prev)
		} else {
			verbs[r.Role] = r.Path
		}
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
	design string
	v      Vocabulary
	out    []Problem
}

func (c *checker) add(r Role, code string, line int, format string, args ...any) {
	c.out = append(c.out, Problem{Path: r.Path, Line: line, Code: code, Message: fmt.Sprintf(format, args...)})
}

func (c *checker) role(r Role) {
	file := func(key, format string, args ...any) { c.add(r, CodeFile, r.line(key), format, args...) }
	for _, k := range []struct{ key, value string }{{"title", r.Title}, {"updated", r.Updated}, {"status", r.Status}, {"name", r.Name}, {"role", r.Role}, {"action", r.Action}} {
		if strings.TrimSpace(k.value) == "" {
			file(k.key, "front matter needs %s; add it (ADR-0207 §2)", k.key)
		}
	}
	if r.Updated != "" {
		if _, err := time.Parse(time.DateOnly, r.Updated); err != nil {
			file("updated", "updated is %q, not a date; give the day of the last meaningful edit as YYYY-MM-DD", r.Updated)
		}
	}
	if r.Status != "" && !slices.Contains(Statuses, r.Status) {
		file("status", "status is %q, which is none of %s", r.Status, strings.Join(Statuses, ", "))
	}
	stem := strings.TrimSuffix(path.Base(r.Path), ".md")
	switch {
	case r.Name == "":
	case !kebab.MatchString(r.Name):
		file("name", "name %q is not kebab-case; use lowercase words joined by dashes", r.Name)
	case r.Name != stem:
		file("name", "name is %s but the file is %s.md; make them the same", r.Name, stem)
	case slices.Contains(ReservedNames, r.Name):
		file("name", "name %s is a sub-agent's definition; give the role another name", r.Name)
	}
	switch {
	case r.Role == "":
	case !kebab.MatchString(r.Role):
		file("role", "role %q is not kebab-case; use a lowercase verb, such as plan", r.Role)
	case slices.Contains(ReservedRoles, r.Role):
		file("role", "role %s is the story's agent's or a sub-agent's; give the role another verb", r.Role)
	}
	c.action(r)
	c.triggers(r)
	if r.Lines["item"] == 0 {
		file("item", "front matter needs item; give none or a list of epic and story (ADR-0207 §2)")
	}
	for _, k := range r.Item {
		if !slices.Contains(ItemKinds, k) {
			file("item", "item kind %q is none of %s; give none or a list of them", k, strings.Join(ItemKinds, ", "))
		}
	}
	for _, t := range r.Topics {
		if !workitem.ValidTopic(t) {
			file("topics", "topic %q is not a word; a topic is letters, digits, dot, dash, or underscore, such as planning", t)
		}
	}
	if !slices.Contains(Charges, r.Charge) {
		file("charge", "charge is %q, which is none of %s", r.Charge, strings.Join(Charges, ", "))
	}
	for _, e := range r.Tools {
		c.call(r, e, "tools", "")
	}
	for _, e := range r.Commands {
		c.call(r, e, "commands", "")
	}
	for _, f := range r.Files {
		c.folder(r, f)
	}
	for _, p := range r.Permissions {
		c.permission(r, p)
	}
	switch n := strings.Count(r.Body, Marker); n {
	case 0:
		c.add(r, CodeFile, 1, "the body has no baseline marker %s; end the baseline text with it, followed by ## Project additions (ADR-0207 §4)", Marker)
	case 1:
	default:
		c.add(r, CodeFile, 1, "the baseline marker appears %d times; keep exactly one", n)
	}
	c.definition(r)
}

func (c *checker) action(r Role) {
	switch a := r.Action; {
	case a == "":
	case !kebab.MatchString(a):
		c.add(r, CodeFile, r.line("action"), "action %q is not kebab-case; use lowercase words joined by dashes", a)
	case slices.Contains(ReservedActions, a):
		c.add(r, CodeFile, r.line("action"), "action %s is a host action that is no role's; give the role an action of its own", a)
	case BuiltinActions[a] != "" && BuiltinActions[a] != r.Name:
		c.add(r, CodeFile, r.line("action"), "action %s is the %s's; give the role an action of its own", a, BuiltinActions[a])
	}
}

func (c *checker) triggers(r Role) {
	if len(r.Trigger) == 0 {
		c.add(r, CodeFile, r.line("trigger"), "front matter needs trigger; give asked, continuous, an {event}, or a {schedule} (ADR-0207 §3)")
		return
	}
	schedules := append(slices.Clone(Schedules), "roles."+r.Name+".schedule")
	for _, t := range r.Trigger {
		at := func(format string, args ...any) { c.add(r, CodeFile, t.Line, format, args...) }
		switch t.Kind {
		case TriggerAsked:
		case TriggerContinuous:
			if len(r.Trigger) > 1 {
				at("continuous is a role's only trigger; remove the others, or continuous")
			}
		case TriggerEvent:
			if !slices.Contains(Events, t.Event) {
				at("event %q is none of %s", t.Event, strings.Join(Events, ", "))
			}
		case TriggerSchedule:
			if !slices.Contains(schedules, t.Schedule) {
				at("schedule %q is none of %s; name the manifest setting that holds the schedule", t.Schedule, strings.Join(schedules, ", "))
			}
		default:
			at("trigger %q is none of asked, continuous, {event: <event>}, or {schedule: <setting>}", t.Kind)
		}
		if t.When == "" {
			continue
		}
		cond, ok := Condition(t.When)
		switch {
		case !ok:
			at("condition %q is not one flai implements; use one of %s, or ask for a new one", t.When, strings.Join(ConditionNames(), ", "))
		case !cond.OnTrigger():
			at("condition %s qualifies a call, not a trigger; remove it from the trigger", t.When)
		case !cond.QualifiesTrigger(t):
			at("condition %s does not qualify this trigger; %s", t.When, triggerTargets(cond))
		}
	}
}

// triggerTargets says which trigger entries a trigger condition qualifies.
func triggerTargets(cond Cond) string {
	if cond.Schedule {
		return "it qualifies a {schedule} entry"
	}
	return "it qualifies an {event} entry for " + strings.Join(cond.Events, ", ")
}

// call checks one call entry in field (tools or commands), inside the
// permission p, or outside any when p is empty.
func (c *checker) call(r Role, e Call, field, p string) {
	at := func(code, format string, args ...any) { c.add(r, code, e.Line, format, args...) }
	command := strings.HasPrefix(e.Call, "flai ")
	switch {
	case strings.HasPrefix(e.Call, "git "):
		at(CodeGrant, "%s names git, which no role file grants: every role has git's reads, and no setting allows its writes; remove it", e.Call)
		return
	case field == "tools" && command:
		at(CodeFile, "%s is a flai command; list it under commands", e.Call)
		return
	case field == "commands" && !command:
		at(CodeFile, "%s is not a flai command, flai <command>, flai <command> <subcommand>, or flai <command> *; an MCP tool goes under tools", e.Call)
		return
	case command && !c.knownCommand(e.Call):
		at(CodeFile, "%s is not a flai command; flai --help lists them", e.Call)
		return
	case !command && c.v.Tools != nil && !slices.Contains(c.v.Tools, e.Call):
		at(CodeFile, "%s is not one of flai's MCP tools; name one flai mcp serves", e.Call)
		return
	}
	if e.When != "" {
		cond, ok := Condition(e.When)
		switch {
		case !ok:
			at(CodeFile, "condition %q is not one flai implements; use one of %s, or ask for a new one", e.When, strings.Join(ConditionNames(), ", "))
			return
		case cond.OnTrigger():
			at(CodeFile, "condition %s qualifies a trigger, not a call; remove it from %s", e.When, e.Call)
			return
		case !cond.Qualifies(e.Call):
			at(CodeFile, "condition %s does not qualify %s; it qualifies %s", e.When, e.Call, qualifies(cond))
			return
		}
	}
	c.grant(r, e, p)
}

// qualifies says what a call condition qualifies.
func qualifies(cond Cond) string {
	if cond.Files {
		return "each files entry, and no call"
	}
	return strings.Join(cond.Calls, ", ")
}

// knownCommand reports whether a command entry names a flai command, or a
// command with subcommands for flai <command> *.
func (c *checker) knownCommand(call string) bool {
	words := strings.Fields(strings.TrimPrefix(call, "flai "))
	if len(words) == 0 || slices.Contains(words[:len(words)-1], "*") {
		return false
	}
	wild := words[len(words)-1] == "*"
	if wild {
		words = words[:len(words)-1]
		if len(words) == 0 {
			return false
		}
	}
	if c.v.Commands == nil {
		return true
	}
	cmd := strings.Join(words, " ")
	if !wild {
		return slices.Contains(c.v.Commands, cmd)
	}
	return slices.ContainsFunc(c.v.Commands, func(s string) bool { return strings.HasPrefix(s, cmd+" ") })
}

// grant checks that a call entry inside permission p, or outside any when p
// is empty, grants nothing only the operator's settings, or no setting,
// allows.
func (c *checker) grant(r Role, e Call, p string) {
	at := func(format string, args ...any) { c.add(r, CodeGrant, e.Line, format, args...) }
	name := e.Call
	if e.When != "" {
		name = e.Call + " when " + e.When
	}
	if e.When != "" {
		cond, _ := Condition(e.When)
		if cond.Permission != "" && cond.Permission != p {
			where := "outside any permission"
			if p != "" {
				where = "under " + p
			}
			at("%s is gated by %s but is named %s; move it into the permissions entry whose setting is %s", name, cond.Permission, where, settingFor(r, cond.Permission))
			return
		}
	} else {
		for _, t := range PlainGated {
			if covers(e.Call, t) && !GrantsPlain(t, p) {
				at("%s named plainly allows %s, which only the operator's settings allow; %s", e.Call, t, howToGrant(r, t))
				return
			}
		}
	}
	if f, ok := Forbids(e.Call, e.When, r.WorksStory); ok {
		if f.StoryAgent {
			at("%s, a story's agent's call, while works_story is false; remove it, or set works_story: true", reaches(name, f))
			return
		}
		at("%s, which no setting allows (ADR-0207 §5); remove it", reaches(name, f))
	}
}

// reaches says that the entry name allows the forbidden call f, or its
// form: "flai queue * reaches flai queue hold", "flai push is never granted".
func reaches(name string, f Forbidden) string {
	target := f.Call
	if f.Form != "" {
		target += " " + f.Form
	}
	if name == f.Call && f.Form == "" {
		return name + " is never granted"
	}
	return name + " reaches " + target
}

// howToGrant says how a role file names a plain-gated call: under a
// condition no permission gates, or inside the permission that allows it.
func howToGrant(r Role, call string) string {
	var ways []string
	for _, cond := range Conditions {
		if cond.Permission == "" && cond.Qualifies(call) {
			ways = append(ways, "under "+cond.Name)
		}
	}
	for _, g := range Gated {
		if g.Call != call {
			continue
		}
		how := "plainly"
		if g.Condition != "" {
			how = "under " + g.Condition
		}
		ways = append(ways, how+" in the permissions entry for "+settingFor(r, g.Permission))
	}
	return "name it " + strings.Join(ways, ", or ")
}

// settingFor is the manifest setting a role's file names for permission p.
func settingFor(r Role, p string) string {
	if r.Name == Orchestrator {
		return "orchestration.permissions." + p
	}
	return "roles." + r.Name + ".permissions." + p
}

// folder checks one files entry.
func (c *checker) folder(r Role, f string) {
	at := func(format string, args ...any) { c.add(r, CodeFile, r.line("files"), format, args...) }
	if !strings.HasSuffix(f, "/") {
		at("files entry %q is not a folder; end it with /", f)
		return
	}
	resolved := f
	if rest, ok := strings.CutPrefix(f, "{design}/"); ok {
		resolved = c.design + "/" + rest
	}
	if strings.Contains(resolved, "{") || path.IsAbs(resolved) || !filepath.IsLocal(filepath.FromSlash(resolved)) || slices.Contains(strings.Split(resolved, "/"), "..") {
		at("files entry %q is not a folder path from the project root; give one such as {design}/analysis/", f)
		return
	}
	clean := path.Clean(resolved)
	if clean == "." {
		at("files entry %q is the project root; name the folder under it the role edits", f)
		return
	}
	for _, guarded := range []string{".git", ".claude", c.design + "/workflows"} {
		if clean == guarded || strings.HasPrefix(clean, guarded+"/") || strings.HasPrefix(guarded, clean+"/") {
			at("files entry %q reaches %s/, which no role edits, so that no role widens itself or touches git or the harness; name a folder outside it", f, guarded)
			return
		}
	}
}

// permission checks one permissions entry and the calls under it.
func (c *checker) permission(r Role, p Permission) {
	at := func(code, format string, args ...any) { c.add(r, code, p.Line, format, args...) }
	if len(p.Tools) == 0 && len(p.Commands) == 0 {
		at(CodeFile, "permissions entry %s allows nothing; give it tools, commands, or both", p.Setting)
	}
	perm, owner, ok := parseSetting(p.Setting)
	switch {
	case p.Setting == "":
		at(CodeFile, "permissions entry needs setting, the operator's switch, such as %s", settingFor(r, Publish))
		return
	case !ok:
		at(CodeFile, "setting %q is neither orchestration.permissions.<p> nor roles.<name>.permissions.<p>; name %s", p.Setting, settingFor(r, "<p>"))
		return
	case !slices.Contains(Permissions, perm):
		at(CodeFile, "permission %q is none of %s", perm, strings.Join(Permissions, ", "))
		return
	case owner != r.Name:
		at(CodeGrant, "setting %s is not this role's: a role names only its own switches, so name %s", p.Setting, settingFor(r, perm))
	}
	for _, e := range p.Tools {
		c.call(r, e, "tools", perm)
	}
	for _, e := range p.Commands {
		c.call(r, e, "commands", perm)
	}
}

// parseSetting reads a permission setting's <p> and the role it belongs to:
// the orchestrator for orchestration.permissions.<p>, <name> for
// roles.<name>.permissions.<p>.
func parseSetting(s string) (perm, owner string, ok bool) {
	if p, found := strings.CutPrefix(s, "orchestration.permissions."); found {
		return p, Orchestrator, p != ""
	}
	rest, found := strings.CutPrefix(s, "roles.")
	if !found {
		return "", "", false
	}
	name, p, found := strings.Cut(rest, ".permissions.")
	if !found || name == "" || p == "" {
		return "", "", false
	}
	if name == Orchestrator {
		// The orchestrator's switches are orchestration.permissions.
		return p, "", true
	}
	return p, name, true
}

// definition checks the harness definition's fields against the role's
// rules (ADR-0207 §4).
func (c *checker) definition(r Role) {
	d := r.Definition
	at := func(key, format string, args ...any) { c.add(r, CodeDefinition, r.line(key), format, args...) }
	if r.Lines["definition"] == 0 {
		at("definition", "front matter needs definition, with description and tools, which flai generates the role's harness definition from (ADR-0207 §4)")
		return
	}
	if strings.TrimSpace(d.Description) == "" {
		at("definition", "definition needs description: when the harness should use the agent")
	}
	if len(d.Tools) == 0 {
		at("definition", "definition needs tools: the harness tools the session has, flai's as mcp__flai__<tool>")
	}
	toolsAt := "definition.tools"
	if r.Lines[toolsAt] == 0 {
		toolsAt = "definition"
	}
	granted := map[string]bool{}
	for _, e := range r.Tools {
		granted[e.Call] = true
	}
	for _, p := range r.Permissions {
		for _, e := range p.Tools {
			granted[e.Call] = true
		}
	}
	for _, t := range d.Tools {
		if tool, ok := strings.CutPrefix(t, MCPPrefix); ok && c.v.Reads != nil && !slices.Contains(c.v.Reads, tool) && !granted[tool] {
			at(toolsAt, "definition.tools lists %s, which the role's rules refuse: %s is neither a read every agent has nor in tools or a permission's tools; remove it, or grant it", t, tool)
		}
		if slices.Contains(EditTools, t) && len(r.Files) == 0 {
			at(toolsAt, "definition.tools lists %s, but files names no folder the role may edit; remove it, or name the folder in files", t)
		}
	}
}
