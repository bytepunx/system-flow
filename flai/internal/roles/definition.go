package roles

import (
	"fmt"
	"slices"
	"strings"

	"github.com/goccy/go-yaml"

	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// Inherit is a definition's model when it takes the session's.
const Inherit = "inherit"

// Generated is a role's harness definition, as flai generates it from the
// role's file (ADR-0207 §4, ADR-0236).
type Generated struct {
	Name        string
	Description string
	Tools       []string
	// Model is the file's definition.model, Inherit when it names none.
	Model string
	// Prompt is the file's body without the marker's line, trimmed.
	Prompt string
}

// Generate is the harness definition flai generates from the role file r:
// its name, its definition's description, tools, and model, and its body,
// the project's additions included, as the prompt. A file whose definition
// has no description or no tools, or whose body is empty, is refused,
// naming flai check, which reports it.
func Generate(r Role) (Generated, error) {
	d := r.Definition
	g := Generated{Name: r.Name, Description: d.Description, Tools: slices.Clone(d.Tools), Model: d.Model, Prompt: Prompt(r.Body)}
	if g.Model == "" {
		g.Model = Inherit
	}
	var lacks []string
	if strings.TrimSpace(d.Description) == "" {
		lacks = append(lacks, "definition.description")
	}
	if len(d.Tools) == 0 {
		lacks = append(lacks, "definition.tools")
	}
	if g.Prompt == "" {
		lacks = append(lacks, "body")
	}
	if len(lacks) > 0 {
		return Generated{}, fmt.Errorf("the role %s runs as the agent %s, which flai generates from its file %s, and the file has no %s; flai check reports what it lacks", r.Name, r.Name, r.Path, strings.Join(lacks, " and no "))
	}
	return g, nil
}

// Agent is the definition as Claude Code's --agents takes it: its
// description, tools, and prompt, and its model, model when that is not
// empty and the definition's otherwise, left out when it is Inherit, so
// that the session's holds.
func (g Generated) Agent(model string) map[string]any {
	agent := map[string]any{"description": g.Description, "tools": slices.Clone(g.Tools), "prompt": g.Prompt}
	if model == "" {
		model = g.Model
	}
	if model != "" && model != Inherit {
		agent["model"] = model
	}
	return agent
}

// File is the definition as a Claude Code agent file: front matter with its
// name, description, tools on one comma-separated line, and model, then its
// prompt (ADR-0236).
func (g Generated) File() []byte {
	model := g.Model
	if model == "" {
		model = Inherit
	}
	var b strings.Builder
	fmt.Fprintf(&b, "---\nname: %s\ndescription: %s\ntools: %s\nmodel: %s\n---\n\n%s\n",
		workitem.Scalar(g.Name), workitem.Scalar(g.Description), strings.Join(g.Tools, ", "), workitem.Scalar(model), g.Prompt)
	return []byte(b.String())
}

// ParseAgentFile reads a Claude Code agent file, front matter and a prompt,
// as a Generated: its tools a comma-separated line or a list, its model
// Inherit when it names none, and its body, trimmed, as the prompt.
func ParseAgentFile(data []byte) (Generated, error) {
	text := strings.ReplaceAll(string(data), "\r\n", "\n")
	fm, body, err := workitem.SplitFrontMatter(text)
	if err != nil {
		return Generated{}, fmt.Errorf("reading an agent file: %w; it opens with --- and its front matter, and closes it with ---", err)
	}
	var f struct {
		Name        string `yaml:"name"`
		Description string `yaml:"description"`
		Tools       any    `yaml:"tools"`
		Model       string `yaml:"model"`
	}
	if err := yaml.Unmarshal([]byte(fm), &f); err != nil {
		return Generated{}, fmt.Errorf("reading an agent file's front matter: %w; fix the YAML", err)
	}
	tools, err := toolList(f.Tools)
	if err != nil {
		return Generated{}, err
	}
	g := Generated{Name: f.Name, Description: f.Description, Tools: tools, Model: f.Model, Prompt: strings.TrimSpace(body)}
	if g.Model == "" {
		g.Model = Inherit
	}
	return g, nil
}

// toolList is an agent file's tools, a comma-separated line or a list of
// names, as a list; none is an empty one.
func toolList(v any) ([]string, error) {
	out := []string{}
	switch t := v.(type) {
	case nil:
	case string:
		for _, name := range strings.Split(t, ",") {
			if name = strings.TrimSpace(name); name != "" {
				out = append(out, name)
			}
		}
	case []any:
		for _, e := range t {
			name, ok := e.(string)
			if !ok {
				return nil, fmt.Errorf("an agent file's tools list holds %v, which is not a tool's name; list the tools' names", e)
			}
			out = append(out, name)
		}
	default:
		return nil, fmt.Errorf("an agent file's tools are %v, neither a comma-separated line nor a list; give them as one", v)
	}
	return out, nil
}

// Prompt is a role file's body as its definition's prompt (ADR-0207 §4):
// the baseline text and the project's additions, without the marker's line
// and the blank lines that followed it, and without the additions' Start
// prompt section, which is the start prompt's (ADR-0250 §1), trimmed.
func Prompt(body string) string {
	body = strings.ReplaceAll(body, "\r\n", "\n")
	if i := strings.Index(body, Marker); i >= 0 {
		rest, _, _ := cutStart(body[i:])
		body = body[:i] + rest
	}
	var out []string
	skip := false
	for _, l := range strings.Split(body, "\n") {
		blank := strings.TrimSpace(l) == ""
		switch {
		case strings.TrimSpace(l) == Marker:
			skip = len(out) > 0 && strings.TrimSpace(out[len(out)-1]) == ""
			continue
		case skip && blank:
			continue
		}
		skip = false
		out = append(out, l)
	}
	return strings.TrimSpace(strings.Join(out, "\n"))
}
