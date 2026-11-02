package roles

import (
	"reflect"
	"slices"
	"strings"
	"testing"
)

func definedRole() Role {
	return Role{
		Path:       "design/workflows/roles/reviewer.md",
		Name:       "reviewer",
		Definition: Definition{Description: "Reviews a story: against the conventions.", Tools: []string{"Read", "mcp__flai__item_get"}, Model: "sonnet"},
		Body:       "\nYou are the reviewer.\n\n" + Marker + "\n\n## Project additions\n\n- Name the file.\n",
	}
}

// ADR-0207 §4: the definition is the file's name, description, tools, and
// model, and its body without the marker's line, the project's additions
// kept, as its prompt.
func TestGenerateTakesTheRoleFilesDefinitionAndBody(t *testing.T) {
	g, err := Generate(definedRole())
	if err != nil {
		t.Fatal(err)
	}
	want := Generated{Name: "reviewer", Description: "Reviews a story: against the conventions.", Tools: []string{"Read", "mcp__flai__item_get"}, Model: "sonnet",
		Prompt: "You are the reviewer.\n\n## Project additions\n\n- Name the file."}
	if !reflect.DeepEqual(g, want) {
		t.Errorf("Generate = %+v, want %+v", g, want)
	}
	if strings.Contains(g.Prompt, Marker) {
		t.Errorf("the marker is in the prompt: %q", g.Prompt)
	}
	r := definedRole()
	r.Definition.Model = ""
	if g, err := Generate(r); err != nil || g.Model != Inherit {
		t.Errorf("no model: %+v, %v, want %s", g, err, Inherit)
	}
}

// ADR-0207 §4: the agent's model wins over the file's, and a definition
// whose model is inherit leaves it out, so that the session's holds.
func TestAgentPutsTheAgentsModelOverTheFiles(t *testing.T) {
	g, err := Generate(definedRole())
	if err != nil {
		t.Fatal(err)
	}
	for _, c := range []struct{ file, agent, want string }{
		{"sonnet", "claude-opus-5-5", "claude-opus-5-5"},
		{"sonnet", "", "sonnet"},
		{Inherit, "haiku", "haiku"},
		{Inherit, "", ""},
		{"sonnet", Inherit, ""},
	} {
		g.Model = c.file
		a := g.Agent(c.agent)
		model, ok := a["model"]
		if c.want == "" && ok || c.want != "" && model != c.want {
			t.Errorf("file %q, agent %q: model %v, want %q", c.file, c.agent, model, c.want)
		}
		if a["description"] != g.Description || a["prompt"] != g.Prompt || !slices.Equal(a["tools"].([]string), g.Tools) {
			t.Errorf("Agent = %v", a)
		}
	}
	a := g.Agent("")
	a["tools"].([]string)[0] = "Write"
	if g.Tools[0] != "Read" {
		t.Errorf("the --agents tools share the definition's: %v", g.Tools)
	}
}

// ADR-0236: each baseline role's definition, written as a Claude Code agent
// file, reads back as itself, its tools on one comma-separated line.
func TestAgentFileRoundTripsEachBaselineRole(t *testing.T) {
	for _, r := range append(Baseline("design"), definedRole()) {
		g, err := Generate(r)
		if err != nil {
			t.Fatalf("%s: %v", r.Name, err)
		}
		file := g.File()
		if !strings.HasPrefix(string(file), "---\nname: "+r.Name+"\ndescription: ") || !strings.Contains(string(file), "\ntools: "+strings.Join(r.Definition.Tools, ", ")+"\n") {
			t.Errorf("%s: the file's front matter:\n%s", r.Name, file)
		}
		back, err := ParseAgentFile(file)
		if err != nil {
			t.Fatalf("%s: %v\n%s", r.Name, err, file)
		}
		if !reflect.DeepEqual(back, g) {
			t.Errorf("%s: read back %+v, want %+v", r.Name, back, g)
		}
	}
}

// An agent file's tools are a comma-separated line or a list, its model
// inherit when it names none, and its prompt its body trimmed; a file with
// no front matter, or tools of another shape, is refused.
func TestParseAgentFileReadsEitherToolsForm(t *testing.T) {
	g, err := ParseAgentFile([]byte("---\r\nname: x\r\ndescription: d\r\ntools: [Read, Bash]\r\n---\r\n\r\nBody.\r\n"))
	want := Generated{Name: "x", Description: "d", Tools: []string{"Read", "Bash"}, Model: Inherit, Prompt: "Body."}
	if err != nil || !reflect.DeepEqual(g, want) {
		t.Errorf("ParseAgentFile = %+v, %v, want %+v", g, err, want)
	}
	if g, err := ParseAgentFile([]byte("---\nname: x\ndescription: d\ntools: Read ,, Bash\nmodel: haiku\n---\n")); err != nil || !slices.Equal(g.Tools, []string{"Read", "Bash"}) || g.Model != "haiku" || g.Prompt != "" {
		t.Errorf("a comma line: %+v, %v", g, err)
	}
	for _, bad := range []string{"name: x\n", "---\nname: x\n", "---\ntools: {a: b}\n---\n", "---\ntools: [1]\n---\n", "---\nname: [\n---\n"} {
		if _, err := ParseAgentFile([]byte(bad)); err == nil {
			t.Errorf("accepted %q", bad)
		}
	}
}

// ADR-0207 §4: a role file whose definition has no description or no tools,
// or whose body is empty, is refused, naming the file and flai check.
func TestGenerateRefusesWhatTheDefinitionLacks(t *testing.T) {
	for want, mod := range map[string]func(*Role){
		"the role reviewer runs as the agent reviewer, which flai generates from its file design/workflows/roles/reviewer.md, and the file has no definition.description; flai check reports what it lacks": func(r *Role) { r.Definition.Description = " " },
		"and the file has no definition.tools; flai check": func(r *Role) { r.Definition.Tools = nil },
		"and the file has no body; flai check":             func(r *Role) { r.Body = "\n" + Marker + "\n\n" },
		"and the file has no definition.description and no definition.tools and no body": func(r *Role) {
			r.Definition = Definition{}
			r.Body = ""
		},
	} {
		r := definedRole()
		mod(&r)
		if _, err := Generate(r); err == nil || !strings.Contains(err.Error(), want) {
			t.Errorf("%v, want %q", err, want)
		}
	}
}

// ADR-0207 §4: a role file's body is its definition's prompt without the
// marker's line and the blank lines after it; a body with no marker is
// itself, trimmed.
func TestPromptIsTheBodyWithoutTheMarker(t *testing.T) {
	for _, c := range []struct{ body, want string }{
		{"\nBaseline.\n\n" + Marker + "\n\n## Project additions\n\nMore.\n", "Baseline.\n\n## Project additions\n\nMore."},
		{"\r\nBaseline.\r\n" + Marker + "\r\n## Project additions\r\n", "Baseline.\n## Project additions"},
		{"\nNo marker.\n\n", "No marker."},
	} {
		if got := Prompt(c.body); got != c.want {
			t.Errorf("Prompt(%q) = %q, want %q", c.body, got, c.want)
		}
	}
}
