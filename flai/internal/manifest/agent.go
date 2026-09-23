package manifest

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"
)

// Agent is who works a story (S-0103, E-0008): the harness that runs it (a
// command line or an API, such as claude-code), the model it runs, and
// options for that harness. In system-flow.yaml it is the project's default,
// copied into every story created while it is set; in a story's front matter
// it is what flai will start the agent for that story with (S-0104). Which
// harnesses flai can start, and which config keys each understands, is up to
// its adapters; here each is only checked for shape.
type Agent struct {
	Harness string            `yaml:"harness,omitempty" json:"harness,omitempty"`
	Model   string            `yaml:"model,omitempty" json:"model,omitempty"`
	Config  map[string]string `yaml:"config,omitempty" json:"config,omitempty"`
}

var (
	harnessPattern   = regexp.MustCompile(`^[a-z][a-z0-9._-]{0,63}$`)
	modelPattern     = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._:/@+-]{0,127}$`)
	configKeyPattern = regexp.MustCompile(`^[a-z][a-z0-9_.-]{0,63}$`)
)

// IsZero reports whether nothing is set.
func (a *Agent) IsZero() bool {
	return a == nil || a.Harness == "" && a.Model == "" && len(a.Config) == 0
}

// Validate checks each part's shape: an identifier for the harness, a model ID
// for the model, and config keys that are identifiers with values on one line.
func (a *Agent) Validate() error {
	if a == nil {
		return nil
	}
	var errs []string
	if a.Harness != "" && !harnessPattern.MatchString(a.Harness) {
		errs = append(errs, fmt.Sprintf("agent harness %q is not a name such as claude-code (lower case letters, digits, . _ -)", a.Harness))
	}
	if a.Model != "" && !modelPattern.MatchString(a.Model) {
		errs = append(errs, fmt.Sprintf("agent model %q is not a model ID such as claude-opus-5-5", a.Model))
	}
	for _, k := range a.ConfigKeys() {
		if !configKeyPattern.MatchString(k) {
			errs = append(errs, fmt.Sprintf("agent config key %q is not a name such as effort or max_turns (lower case letters, digits, . _ -)", k))
		}
		if v := a.Config[k]; strings.ContainsAny(v, "\n\r") {
			errs = append(errs, fmt.Sprintf("agent config %s spans lines; a value is one line", k))
		}
	}
	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "; "))
	}
	return nil
}

// ConfigKeys are the config's keys in order.
func (a *Agent) ConfigKeys() []string {
	if a == nil {
		return nil
	}
	keys := make([]string, 0, len(a.Config))
	for k := range a.Config {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// With is a with over's harness and model where over sets them, and their
// configs merged key by key, over's winning. Nil when the result is empty.
func (a *Agent) With(over *Agent) *Agent {
	out := &Agent{}
	for _, x := range []*Agent{a, over} {
		if x == nil {
			continue
		}
		if x.Harness != "" {
			out.Harness = x.Harness
		}
		if x.Model != "" {
			out.Model = x.Model
		}
		for k, v := range x.Config {
			if out.Config == nil {
				out.Config = map[string]string{}
			}
			out.Config[k] = v
		}
	}
	if out.IsZero() {
		return nil
	}
	return out
}

// String is the agent in one line, for the terminal: claude-code, claude-opus-5-5, effort=high.
func (a *Agent) String() string {
	if a.IsZero() {
		return "none"
	}
	var parts []string
	if a.Harness != "" {
		parts = append(parts, a.Harness)
	}
	if a.Model != "" {
		parts = append(parts, a.Model)
	}
	for _, k := range a.ConfigKeys() {
		parts = append(parts, k+"="+a.Config[k])
	}
	return strings.Join(parts, ", ")
}

// ParseConfig reads key=value pairs, as flags give them, into a config.
func ParseConfig(pairs []string) (map[string]string, error) {
	if len(pairs) == 0 {
		return nil, nil
	}
	out := map[string]string{}
	for _, p := range pairs {
		k, v, ok := strings.Cut(p, "=")
		k = strings.TrimSpace(k)
		if !ok || k == "" {
			return nil, fmt.Errorf("agent config %q is not key=value", p)
		}
		out[k] = strings.TrimSpace(v)
	}
	return out, nil
}

// WriteAgent rewrites the top-level agent block of the manifest at path, or
// removes it when a is empty, leaving every other line as it was: flai owns
// this key, and the rest of the file is the operator's.
func WriteAgent(path string, a *Agent) error {
	if err := a.Validate(); err != nil {
		return err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	lines := strings.Split(strings.TrimRight(string(data), "\n"), "\n")
	var kept []string
	skipping := false
	for _, l := range lines {
		top := l != "" && l[0] != ' ' && l[0] != '\t' && l[0] != '#'
		if top {
			skipping = l == "agent:" || strings.HasPrefix(l, "agent:")
		}
		if !skipping {
			kept = append(kept, l)
		}
	}
	for len(kept) > 0 && strings.TrimSpace(kept[len(kept)-1]) == "" {
		kept = kept[:len(kept)-1]
	}
	out := strings.Join(kept, "\n") + "\n"
	if !a.IsZero() {
		out += AgentBlock(a, "")
	}
	return os.WriteFile(path, []byte(out), 0o644)
}

// AgentBlock renders an agent block at indent, as system-flow.yaml and a
// story's front matter write it.
func AgentBlock(a *Agent, indent string) string {
	var b strings.Builder
	b.WriteString(indent + "agent:\n")
	if a.Harness != "" {
		fmt.Fprintf(&b, "%s  harness: %s\n", indent, yamlScalar(a.Harness))
	}
	if a.Model != "" {
		fmt.Fprintf(&b, "%s  model: %s\n", indent, yamlScalar(a.Model))
	}
	if len(a.Config) > 0 {
		fmt.Fprintf(&b, "%s  config:\n", indent)
		for _, k := range a.ConfigKeys() {
			fmt.Fprintf(&b, "%s    %s: %s\n", indent, k, yamlScalar(a.Config[k]))
		}
	}
	return b.String()
}

var plainValue = regexp.MustCompile(`^[A-Za-z][A-Za-z0-9._/@+-]*$`)

// yamlScalar is a value plain when that reads back as the same string, and
// double-quoted otherwise.
func yamlScalar(s string) string {
	switch strings.ToLower(s) {
	case "true", "false", "null", "yes", "no", "on", "off", "~":
		return fmt.Sprintf("%q", s)
	}
	if plainValue.MatchString(s) {
		return s
	}
	return fmt.Sprintf("%q", s)
}
