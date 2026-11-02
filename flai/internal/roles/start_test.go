package roles

import (
	"slices"
	"strings"
	"testing"
)

// ADR-0250 §2: each field a start prompt may name renders its value, set or
// empty, and the rendered text is trimmed.
func TestRenderStartRendersEachField(t *testing.T) {
	cases := []struct {
		field, tmpl string
		set         StartFields
		want, empty string
	}{
		{".Name", "{{.Name}}", StartFields{Name: "planner-S-0001"}, "planner-S-0001", ""},
		{".Root", "{{.Root}}", StartFields{Root: "/home/op/project"}, "/home/op/project", ""},
		{".Role", "{{.Role}}", StartFields{Role: "plan"}, "plan", ""},
		{".RoleName", "{{.RoleName}}", StartFields{RoleName: "planner"}, "planner", ""},
		{".File", "{{.File}}", StartFields{File: "design/workflows/roles/planner.md"}, "design/workflows/roles/planner.md", ""},
		{".Body", "{{.Body}}", StartFields{Body: "You are the planner."}, "You are the planner.", ""},
		{".Item", "{{.Item}}", StartFields{Item: "S-0001"}, "S-0001", ""},
		{".Kind", "{{.Kind}}", StartFields{Kind: "story"}, "story", ""},
		{".AskedOn", "{{.AskedOn}}", StartFields{AskedOn: "TH-0007"}, "TH-0007", ""},
		{".Reasons", "{{range .Reasons}}{{.Line}}|{{.Prime}};{{end}}", StartFields{Reasons: []StartReason{{Line: "S-0001 moved to review", Prime: "review:S-0001"}, {Line: "the schedule came round"}}}, "S-0001 moved to review|review:S-0001;the schedule came round|;", ""},
		{".Primed", "{{if .Primed}}primed{{else}}not primed{{end}}", StartFields{Primed: true}, "primed", "not primed"},
		{".PrimeForms", "{{range .PrimeForms}}{{.}};{{end}}", StartFields{PrimeForms: []string{"review:S-0001", "ready-empty"}}, "review:S-0001;ready-empty;", ""},
		{".Queue", "{{if .Queue}}queue{{else}}no queue{{end}}", StartFields{Queue: true}, "queue", "no queue"},
		{".Focus", "{{.Focus}}", StartFields{Focus: "risk"}, "risk", ""},
		{".Since", "{{.Since}}", StartFields{Since: "2026-11-01"}, "2026-11-01", ""},
		{".Inputs", "{{range .Inputs}}{{.Name}}={{.Value}};{{end}}", StartFields{Inputs: []StartInput{{Name: "scope", Value: "the flai CLI"}}}, "scope=the flai CLI;", ""},
	}
	if len(cases) != len(StartFieldNames()) {
		t.Fatalf("%d cases for %d fields %v", len(cases), len(StartFieldNames()), StartFieldNames())
	}
	for i, c := range cases {
		if StartFieldNames()[i] != c.field {
			t.Errorf("field %d is %s, want %s in ADR-0250 §2's order", i, StartFieldNames()[i], c.field)
		}
		r := Role{Name: "reviewer", Path: "design/workflows/roles/reviewer.md", Start: "\n  " + c.tmpl + "\n\n"}
		if err := ValidateStart(r.Start); err != nil {
			t.Errorf("%s: ValidateStart = %v", c.field, err)
		}
		if got, err := r.RenderStart("design", c.set); err != nil || got != c.want {
			t.Errorf("%s set: %q, %v, want %q", c.field, got, err, c.want)
		}
		if got, err := r.RenderStart("design", StartFields{}); err != nil || got != c.empty {
			t.Errorf("%s empty: %q, %v, want %q", c.field, got, err, c.empty)
		}
	}
}

// ADR-0250 §1: a project's Start prompt section is rendered over the same
// fields and appended to the rendered start after a blank line, and left
// out of the definition's prompt.
func TestRenderStartAppendsTheProjectsAdditions(t *testing.T) {
	body := "\nYou are the reviewer.\n\n" + Marker + "\n\n## Project additions\n\n- Name the file.\n\n### Start prompt\n\nAlso read {{.Item}}'s narrative.\n\n### Later\n\nMore.\n"
	r, p := Parse("design/workflows/roles/reviewer.md", "---\nname: reviewer\nstart: |\n  Review {{.Item}}.\n---\n"+body)
	if p != nil {
		t.Fatal(p)
	}
	if r.StartAdditions != "Also read {{.Item}}'s narrative." || r.StartAdditionsLine != 15 {
		t.Errorf("additions %q at line %d, want line 15", r.StartAdditions, r.StartAdditionsLine)
	}
	got, err := r.RenderStart("design", StartFields{Item: "S-0042"})
	if want := "Review S-0042.\n\nAlso read S-0042's narrative."; err != nil || got != want {
		t.Errorf("RenderStart = %q, %v, want %q", got, err, want)
	}
	if want := "You are the reviewer.\n\n## Project additions\n\n- Name the file.\n\n### Later\n\nMore."; Prompt(r.Body) != want {
		t.Errorf("Prompt = %q, want %q", Prompt(r.Body), want)
	}
}

// ADR-0250 §1: a role the project defines whose file has no start renders
// flai's default, which says what flai said before start prompts were data
// (harness.definedPrompt), with no item and no inputs, and on a story with
// inputs; a built-in role with no start takes the baseline's.
func TestADefinedRoleWithNoStartRendersTheDefault(t *testing.T) {
	r := Role{Name: "reviewer", Path: "design/workflows/roles/reviewer.md", Role: "review"}
	if r.StartTemplate("design") != DefaultStart {
		t.Fatalf("StartTemplate = %q, want DefaultStart", r.StartTemplate("design"))
	}
	f := StartFields{Name: "reviewer", Root: "/home/op/project", Role: "review", RoleName: "reviewer", File: r.Path, Body: "You are the reviewer.\n\nRead the diff."}
	got, err := r.RenderStart("design", f)
	want := "You are reviewer, the reviewer role this project defines in design/workflows/roles/reviewer.md, started by flai serve on this host because the operator asked for a run of it, in the project at /home/op/project. Your role file says what you are and what you do:\n\nYou are the reviewer.\n\nRead the diff.\n\nThis run is on no item: it works the project as your role file says. Prime your session with flai prime --role review (or the flai MCP tool prime with role review), which prints the conventions you work by, and briefs the design your role's topics select. A brief is not the document: read the section that bears on your work with the flai MCP tool doc_get and its heading before relying on it, and find sections by their words with doc_search. Call the flai MCP tool inbox.\n\nThe run has no inputs.\n\nWork in the main checkout. flai guard holds every call you make to what your role file allows: never work around a refusal, by another tool, another command, or the shell. A refusal from flai guard or from flai ends that attempt.\n\nWhen an input the operator owns is missing, do not guess past it: ask with the flai MCP tool thread_open on the item it concerns, your recommended answer first, do what needs no answer meanwhile, and hold the flai MCP tool wait_for_events, again each time it returns, until the thread is answered; then go on.\n\nEnd with a one-line summary of what you did: flai serve logs the run's activity in wip/agents/reviewer.md with it, with the run's cost, so you need not log it yourself with activity_log."
	if err != nil || got != want {
		t.Errorf("no item, no inputs:\n%q, %v\nwant\n%q", got, err, want)
	}

	f.Name, f.Item, f.Kind = "reviewer-S-0042", "S-0042", "story"
	f.Inputs = []StartInput{{Name: "scope", Value: "the flai CLI"}, {Name: "depth", Value: "deep"}}
	got, err = r.RenderStart("design", f)
	want = "You are reviewer-S-0042, the reviewer role this project defines in design/workflows/roles/reviewer.md, started by flai serve on this host because the operator asked for a run of it on S-0042, in the project at /home/op/project. Your role file says what you are and what you do:\n\nYou are the reviewer.\n\nRead the diff.\n\nThis run is on S-0042: read it with the flai MCP tool item_get, and what it links with item_get and doc_get. Prime your session with flai prime --role review --story S-0042 (or the flai MCP tool prime with role review and story S-0042), which prints the conventions you work by and what S-0042 names whole, and briefs the design your role's topics select. A brief is not the document: read the section that bears on your work with the flai MCP tool doc_get and its heading before relying on it, and find sections by their words with doc_search. Call the flai MCP tool inbox.\n\nThe run's inputs, in the order given:\n\n- scope: the flai CLI\n- depth: deep\n\nWork in the main checkout. flai guard holds every call you make to what your role file allows: never work around a refusal, by another tool, another command, or the shell. A refusal from flai guard or from flai ends that attempt.\n\nWhen an input the operator owns is missing, do not guess past it: ask with the flai MCP tool thread_open on S-0042, your recommended answer first, do what needs no answer meanwhile, and hold the flai MCP tool wait_for_events, again each time it returns, until the thread is answered; then go on.\n\nEnd with a one-line summary of what you did: flai serve logs the run's activity in wip/agents/reviewer.md with it, with the run's cost, so you need not log it yourself with activity_log."
	if err != nil || got != want {
		t.Errorf("a story with two inputs:\n%q, %v\nwant\n%q", got, err, want)
	}
	if err := ValidateStart(DefaultStart); err != nil {
		t.Errorf("ValidateStart(DefaultStart) = %v", err)
	}

	for _, name := range baselineNames {
		b := Role{Name: name}
		if got := b.StartTemplate("design"); got != baselineRole("design", name).Start {
			t.Errorf("%s: StartTemplate = %q, want the baseline's", name, got)
		}
	}
	r.Start = "Hello {{.Name}}."
	if got, err := r.RenderStart("design", f); err != nil || got != "Hello reviewer-S-0042." {
		t.Errorf("the file's start: %q, %v", got, err)
	}
}

// ADR-0250 §2: a start that does not parse, or names a field outside the
// table, in any branch, is refused, naming what failed and where.
func TestValidateStartRefusesWhatDoesNotRender(t *testing.T) {
	for _, c := range []struct{ tmpl, want string }{
		{"Plan {{.Item}.", "does not parse"},
		{"Plan {{upper .Item}}.", "does not parse"},
		{"Plan {{.Story}}.", "does not render over a run's fields"},
		{"{{if .Item}}{{.Story}}{{end}}", "can't evaluate field Story"},
		{"{{range .Inputs}}{{.Key}}{{end}}", "can't evaluate field Key"},
		{"{{range .Reasons}}{{.Reason}}{{end}}", "can't evaluate field Reason"},
	} {
		if err := ValidateStart(c.tmpl); err == nil || !strings.Contains(err.Error(), c.want) {
			t.Errorf("ValidateStart(%q) = %v, want %q", c.tmpl, err, c.want)
		}
	}
	r := Role{Name: "reviewer", Path: "design/workflows/roles/reviewer.md", Start: "{{.Story}}"}
	if _, err := r.RenderStart("design", StartFields{}); err == nil || !strings.Contains(err.Error(), "the start field of design/workflows/roles/reviewer.md") {
		t.Errorf("RenderStart = %v, want the file named", err)
	}
	r = Role{Name: "reviewer", Path: "design/workflows/roles/reviewer.md", Start: "Hi.", StartAdditions: "{{.Story}}"}
	if _, err := r.RenderStart("design", StartFields{}); err == nil || !strings.Contains(err.Error(), "the ### Start prompt section of design/workflows/roles/reviewer.md") {
		t.Errorf("RenderStart = %v, want the section named", err)
	}
}

// ADR-0250 §1: the Start prompt section runs to the next heading of level
// three or above, skipping fenced code; text with none is left as it is.
func TestCutStartFindsTheSection(t *testing.T) {
	text := Marker + "\n\n## Project additions\n\n### Start prompt\n\nOne.\n```\n# not a heading\n```\n#### Deeper\n\nTwo.\n\n## After\n\nKept.\n\n### Start prompt\nThree.\n"
	rest, section, line := cutStart(text)
	if want := Marker + "\n\n## Project additions\n\n## After\n\nKept.\n"; rest != want {
		t.Errorf("rest = %q, want %q", rest, want)
	}
	if want := "One.\n```\n# not a heading\n```\n#### Deeper\n\nTwo.\n\nThree."; section != want || line != 4 {
		t.Errorf("section %q at %d, want %q at 4", section, line, want)
	}
	for _, none := range []string{"", Marker + "\n\n## Project additions\n\n### Start prompts\n", "\n\n- ### Start prompt\n"} {
		if rest, section, line := cutStart(none); rest != none || section != "" || line != -1 {
			t.Errorf("cutStart(%q) = %q, %q, %d", none, rest, section, line)
		}
	}
	if !slices.Contains(StartFieldNames(), ".AskedOn") {
		t.Errorf("StartFieldNames = %v", StartFieldNames())
	}
}
