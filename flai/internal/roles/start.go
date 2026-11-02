package roles

import (
	"fmt"
	"reflect"
	"regexp"
	"slices"
	"strings"
	"text/template"
)

// StartHeading is the heading in a role file's project additions, below the
// marker, that a project's additions to the start prompt go under (ADR-0250
// §1). Its text runs to the next heading of level three or above, or the end
// of the body.
const StartHeading = "### Start prompt"

// StartFields are the run's fields a role's start prompt is rendered over
// (ADR-0250 §2). The caller fills them; folding an input's whitespace to
// single spaces is the caller's job, as is giving each reason its form.
type StartFields struct {
	// Name is the run's FLAI_AGENT: the role's name, or <name>-<item>.
	Name string
	// Root is the project's directory, where the run works.
	Root string
	// Role is the role's verb, its file's role.
	Role string
	// RoleName is the role's name, its file's name.
	RoleName string
	// File is the role file's path from the project root.
	File string
	// Body is the definition's prompt, the body without the marker (Prompt).
	Body string
	// Item is the epic or story the run is on, empty for none.
	Item string
	// Kind is the item's kind, epic or story, empty for none.
	Kind string
	// AskedOn is the thread the operator wrote on since the item was planned, when the orchestrator asked for it again; empty otherwise.
	AskedOn string
	// Reasons are why the run was started, in the order given.
	Reasons []StartReason
	// Primed says whether every reason has a form, so that the run primes with them.
	Primed bool
	// PrimeForms are each reason's form once, in the order first given.
	PrimeForms []string
	// Queue says whether the merge queue runs for the project.
	Queue bool
	// Focus is what the analyzer looks for, bottlenecks, intent, or risk; empty for all three.
	Focus string
	// Since is the window's start since the role's last run; empty on a first run.
	Since string
	// Inputs are a defined role's run's inputs, in the order given, each value's whitespace folded.
	Inputs []StartInput
}

// StartReason is one reason a run was started.
type StartReason struct {
	// Line is the reason's line.
	Line string
	// Prime is the form flai prime --reason takes it in, empty when it has none.
	Prime string
}

// StartInput is one of a defined role's run's inputs.
type StartInput struct {
	// Name is the input's name.
	Name string
	// Value is the input's value, its whitespace folded to single spaces.
	Value string
}

// DefaultStart is the start prompt a role the project defines renders when
// its file has no start (ADR-0250 §1): what its file's body says, on its
// item, with its inputs, as flai said it before start prompts were data.
const DefaultStart = `You are {{.Name}}, the {{.RoleName}} role this project defines in {{.File}}, started by flai serve on this host because the operator asked for a run of it{{if .Item}} on {{.Item}}{{end}}, in the project at {{.Root}}. Your role file says what you are and what you do:

{{.Body}}

{{if .Item}}This run is on {{.Item}}: read it with the flai MCP tool item_get, and what it links with item_get and doc_get.{{else}}This run is on no item: it works the project as your role file says.{{end}} Prime your session with {{if .Item}}flai prime --role {{.Role}} --{{.Kind}} {{.Item}} (or the flai MCP tool prime with role {{.Role}} and {{.Kind}} {{.Item}}){{else}}flai prime --role {{.Role}} (or the flai MCP tool prime with role {{.Role}}){{end}}, which prints the conventions you work by{{if .Item}} and what {{.Item}} names whole{{end}}, and briefs the design your role's topics select. A brief is not the document: read the section that bears on your work with the flai MCP tool doc_get and its heading before relying on it, and find sections by their words with doc_search. Call the flai MCP tool inbox.

{{if .Inputs}}The run's inputs, in the order given:
{{range .Inputs}}
- {{.Name}}: {{.Value}}{{end}}{{else}}The run has no inputs.{{end}}

Work in the main checkout. flai guard holds every call you make to what your role file allows: never work around a refusal, by another tool, another command, or the shell. A refusal from flai guard or from flai ends that attempt.

When an input the operator owns is missing, do not guess past it: ask with the flai MCP tool thread_open on {{if .Item}}{{.Item}}{{else}}the item it concerns{{end}}, your recommended answer first, do what needs no answer meanwhile, and hold the flai MCP tool wait_for_events, again each time it returns, until the thread is answered; then go on.

End with a one-line summary of what you did: flai serve logs the run's activity in wip/agents/{{.RoleName}}.md with it, with the run's cost, so you need not log it yourself with activity_log.`

// StartTemplate is the start prompt template the role renders (ADR-0250
// §1): its file's start when it has one; else, for a built-in role, the
// start of the baseline file flai carries for it, read as if under the
// design folder design, even when that is empty; else DefaultStart.
func (r Role) StartTemplate(design string) string {
	switch {
	case r.Start != "":
		return r.Start
	case slices.Contains(baselineNames, r.Name):
		return baselineRole(design, r.Name).Start
	}
	return DefaultStart
}

// startSource says where the role's start prompt template comes from, for
// an error.
func (r Role) startSource() string {
	switch {
	case r.Start != "":
		return "the start field of " + r.Path
	case slices.Contains(baselineNames, r.Name):
		return "the baseline start prompt flai carries for " + r.Name
	}
	return "flai's default start prompt"
}

// RenderStart is the role's start prompt for a run with the fields f: its
// StartTemplate rendered over f, then, when the file has a Start prompt
// section in its additions, that section rendered over f after a blank
// line, trimmed of the whitespace around it (ADR-0250 §1, §2).
func (r Role) RenderStart(design string, f StartFields) (string, error) {
	out, err := renderStart("start", r.StartTemplate(design), f)
	if err != nil {
		return "", fmt.Errorf("rendering the role %s's start prompt from %s: %w; flai check reports what is wrong with it, and ADR-0250 §2 lists the fields it may name", r.Name, r.startSource(), err)
	}
	if r.StartAdditions != "" {
		more, err := renderStart("start additions", r.StartAdditions, f)
		if err != nil {
			return "", fmt.Errorf("rendering the %s section of %s: %w; flai check reports what is wrong with it, and ADR-0250 §2 lists the fields it may name", StartHeading, r.Path, err)
		}
		out = strings.TrimSpace(out + "\n\n" + more)
	}
	return out, nil
}

// parseStart parses text as the start prompt template named name.
func parseStart(name, text string) (*template.Template, error) {
	return template.New(name).Option("missingkey=error").Parse(text)
}

// renderStart renders the template text, named name, over f, trimmed.
func renderStart(name, text string, f StartFields) (string, error) {
	t, err := parseStart(name, text)
	if err != nil {
		return "", err
	}
	var b strings.Builder
	if err := t.Execute(&b, f); err != nil {
		return "", err
	}
	return strings.TrimSpace(b.String()), nil
}

// sampleStart is a run with every field set: each string non-empty, each
// list one element long, and each flag true.
var sampleStart = StartFields{
	Name: "reviewer-S-0001", Root: "/project", Role: "review", RoleName: "reviewer", File: "design/workflows/roles/reviewer.md",
	Body: "You are the reviewer.", Item: "S-0001", Kind: "story", AskedOn: "TH-0001",
	Reasons: []StartReason{{Line: "S-0001 moved to review", Prime: "review:S-0001"}}, Primed: true, PrimeForms: []string{"review:S-0001"},
	Queue: true, Focus: "risk", Since: "2026-01-01", Inputs: []StartInput{{Name: "scope", Value: "the flai CLI"}},
}

// ValidateStart checks the start prompt template text (ADR-0250 §2): that
// it parses, and that it renders over a run with every field empty and over
// one with every field set, so that a field outside StartFields is found
// in any branch. Its error says which failed, with the template's own.
func ValidateStart(text string) error {
	if _, err := parseStart("start", text); err != nil {
		return fmt.Errorf("does not parse: %w", err)
	}
	for _, f := range []StartFields{{}, sampleStart} {
		if _, err := renderStart("start", text, f); err != nil {
			return fmt.Errorf("does not render over a run's fields: %w", err)
		}
	}
	return nil
}

// StartFieldNames are the fields a start prompt may name, as .Name, .Root,
// and so on, in ADR-0250 §2's order.
func StartFieldNames() []string {
	t := reflect.TypeFor[StartFields]()
	out := make([]string, t.NumField())
	for i := range out {
		out[i] = "." + t.Field(i).Name
	}
	return out
}

// headingUpTo3 is a markdown heading of level one to three.
var headingUpTo3 = regexp.MustCompile(`^#{1,3}(?:[ \t]|$)`)

// cutStart finds the Start prompt sections in text, a role file's additions
// with LF line ends: it returns text without them, each heading and the
// text under it up to the next heading of level one to three or the end;
// their text, trimmed and joined by a blank line; and the 0-based line in
// text of the first heading, -1 when there is none. A heading inside a
// fenced code block is no heading. Text with no section is returned as it
// is.
func cutStart(text string) (rest, section string, line int) {
	lines := strings.Split(text, "\n")
	var kept, sections []string
	line = -1
	fenced, in := false, false
	var cur []string
	flush := func() {
		if s := strings.TrimSpace(strings.Join(cur, "\n")); s != "" {
			sections = append(sections, s)
		}
		cur = nil
	}
	for i, l := range lines {
		trimmed := strings.TrimSpace(l)
		heading := !fenced && headingUpTo3.MatchString(trimmed)
		if strings.HasPrefix(trimmed, "```") || strings.HasPrefix(trimmed, "~~~") {
			fenced = !fenced
		}
		switch {
		case heading && trimmed == StartHeading:
			if in {
				flush()
			}
			in = true
			if line < 0 {
				line = i
			}
			continue
		case heading && in:
			flush()
			in = false
		}
		if in {
			cur = append(cur, l)
			continue
		}
		kept = append(kept, l)
	}
	if line < 0 {
		return text, "", -1
	}
	if in {
		flush()
	}
	return strings.Join(kept, "\n"), strings.Join(sections, "\n\n"), line
}
