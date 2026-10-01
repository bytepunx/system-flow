package manifest

import (
	"errors"
	"fmt"
	"maps"
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
// its adapters; here each is only checked for shape. Roles are the agents
// for the parts of a story's work its agent hands to sub-agents (S-0189):
// explore and verify, by role name, each a harness, a model, and config of
// its own; the fields above are the story's own.
type Agent struct {
	Harness string            `yaml:"harness,omitempty" json:"harness,omitempty"`
	Model   string            `yaml:"model,omitempty" json:"model,omitempty"`
	Config  map[string]string `yaml:"config,omitempty" json:"config,omitempty"`
	Roles   map[string]Role   `yaml:"roles,omitempty" json:"roles,omitempty"`
}

// Role is the agent for one role of a story's work other than the story's
// own: what a sub-agent in that role runs (S-0189).
type Role struct {
	Harness string            `yaml:"harness,omitempty" json:"harness,omitempty"`
	Model   string            `yaml:"model,omitempty" json:"model,omitempty"`
	Config  map[string]string `yaml:"config,omitempty" json:"config,omitempty"`
}

// The roles flai's harnesses know; a role is any name of this shape, and an
// adapter refuses one it cannot start.
const (
	RoleExplore = "explore"
	RoleVerify  = "verify"
)

var (
	harnessPattern   = regexp.MustCompile(`^[a-z][a-z0-9._-]{0,63}$`)
	modelPattern     = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._:/@+-]{0,127}$`)
	configKeyPattern = regexp.MustCompile(`^[a-z][a-z0-9_.-]{0,63}$`)
	rolePattern      = regexp.MustCompile(`^[a-z][a-z0-9_-]{0,31}$`)
)

// IsZero reports whether nothing is set.
func (a *Agent) IsZero() bool {
	return a == nil || a.Harness == "" && a.Model == "" && len(a.Config) == 0 && len(a.Roles) == 0
}

// ConfigKeys are the role's config keys in order.
func (r Role) ConfigKeys() []string {
	return sortedKeys(r.Config)
}

// IsZero reports whether the role sets nothing.
func (r Role) IsZero() bool {
	return r.Harness == "" && r.Model == "" && len(r.Config) == 0
}

// Validate checks each part's shape: an identifier for the harness, a model ID
// for the model, config keys that are identifiers with values on one line,
// and role names that are identifiers, each role setting something and
// checked the same way.
func (a *Agent) Validate() error {
	if a == nil {
		return nil
	}
	errs := validate("agent", a.Harness, a.Model, a.Config)
	for _, n := range a.RoleNames() {
		r := a.Roles[n]
		switch {
		case !rolePattern.MatchString(n):
			errs = append(errs, fmt.Sprintf("agent role %q is not a name such as explore or verify (lower case letters, digits, _ -)", n))
		case r.IsZero():
			errs = append(errs, fmt.Sprintf("agent role %s sets nothing; give it a harness, a model, or config, or leave it out", n))
		default:
			errs = append(errs, validate("agent role "+n, r.Harness, r.Model, r.Config)...)
		}
	}
	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "; "))
	}
	return nil
}

func validate(what, harness, model string, config map[string]string) []string {
	var errs []string
	if harness != "" && !harnessPattern.MatchString(harness) {
		errs = append(errs, fmt.Sprintf("%s harness %q is not a name such as claude-code (lower case letters, digits, . _ -)", what, harness))
	}
	if model != "" && !modelPattern.MatchString(model) {
		errs = append(errs, fmt.Sprintf("%s model %q is not a model ID such as claude-opus-5-5", what, model))
	}
	for _, k := range sortedKeys(config) {
		if !configKeyPattern.MatchString(k) {
			errs = append(errs, fmt.Sprintf("%s config key %q is not a name such as effort or max_turns (lower case letters, digits, . _ -)", what, k))
		}
		if v := config[k]; strings.ContainsAny(v, "\n\r") {
			errs = append(errs, fmt.Sprintf("%s config %s spans lines; a value is one line", what, k))
		}
	}
	return errs
}

// Same reports whether a and b say the same thing, roles included; nothing
// set is the same as nil.
func (a *Agent) Same(b *Agent) bool {
	if a.IsZero() || b.IsZero() {
		return a.IsZero() && b.IsZero()
	}
	if a.Harness != b.Harness || a.Model != b.Model || !sameConfig(a.Config, b.Config) || len(a.Roles) != len(b.Roles) {
		return false
	}
	for n, r := range a.Roles {
		s, ok := b.Roles[n]
		if !ok || r.Harness != s.Harness || r.Model != s.Model || !sameConfig(r.Config, s.Config) {
			return false
		}
	}
	return true
}

func sameConfig(a, b map[string]string) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if w, ok := b[k]; !ok || w != v {
			return false
		}
	}
	return true
}

// Clone is a copy of a that shares no map with it; nil for nil.
func (a *Agent) Clone() *Agent {
	if a == nil {
		return nil
	}
	c := *a
	c.Config = maps.Clone(a.Config)
	if a.Roles != nil {
		c.Roles = make(map[string]Role, len(a.Roles))
		for n, r := range a.Roles {
			r.Config = maps.Clone(r.Config)
			c.Roles[n] = r
		}
	}
	return &c
}

// RoleNames are the roles' names in order.
func (a *Agent) RoleNames() []string {
	if a == nil {
		return nil
	}
	return sortedKeys(a.Roles)
}

func sortedKeys[V any](m map[string]V) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// ConfigKeys are the config's keys in order.
func (a *Agent) ConfigKeys() []string {
	if a == nil {
		return nil
	}
	return sortedKeys(a.Config)
}

// With is a with over's harness and model where over sets them, and their
// configs merged key by key, over's winning; each role is merged the same
// way, role by role. Nil when the result is empty.
func (a *Agent) With(over *Agent) *Agent {
	out := &Agent{}
	for _, x := range []*Agent{a, over} {
		if x == nil {
			continue
		}
		out.Harness, out.Model, out.Config = merge(out.Harness, out.Model, out.Config, x.Harness, x.Model, x.Config)
		for n, r := range x.Roles {
			if out.Roles == nil {
				out.Roles = map[string]Role{}
			}
			o := out.Roles[n]
			o.Harness, o.Model, o.Config = merge(o.Harness, o.Model, o.Config, r.Harness, r.Model, r.Config)
			out.Roles[n] = o
		}
	}
	if out.IsZero() {
		return nil
	}
	return out
}

// merge lays harness, model, and config over the ones before, into a new
// config map so that neither side's is shared.
func merge(harness, model string, config map[string]string, h, m string, c map[string]string) (string, string, map[string]string) {
	if h != "" {
		harness = h
	}
	if m != "" {
		model = m
	}
	var out map[string]string
	for _, src := range []map[string]string{config, c} {
		for k, v := range src {
			if out == nil {
				out = map[string]string{}
			}
			out[k] = v
		}
	}
	return harness, model, out
}

// String is the agent in one line, for the terminal: claude-code,
// claude-opus-5-5, effort=high; explore: haiku; verify: sonnet.
func (a *Agent) String() string {
	if a.IsZero() {
		return "none"
	}
	parts := []string{line(a.Harness, a.Model, a.Config)}
	if parts[0] == "" {
		parts[0] = "default"
	}
	for _, n := range a.RoleNames() {
		r := a.Roles[n]
		parts = append(parts, n+": "+line(r.Harness, r.Model, r.Config))
	}
	return strings.Join(parts, "; ")
}

func line(harness, model string, config map[string]string) string {
	var parts []string
	if harness != "" {
		parts = append(parts, harness)
	}
	if model != "" {
		parts = append(parts, model)
	}
	for _, k := range sortedKeys(config) {
		parts = append(parts, k+"="+config[k])
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

// ParseRoles reads roles as flags give them: role=harness pairs, role=model
// pairs, and role.key=value config. Nil when none is given.
func ParseRoles(harnesses, models, configs []string) (map[string]Role, error) {
	out := map[string]Role{}
	set := func(flag string, pairs []string, put func(r *Role, v string)) error {
		for _, p := range pairs {
			n, v, ok := strings.Cut(p, "=")
			n = strings.TrimSpace(n)
			if !ok || n == "" {
				return fmt.Errorf("agent %s %q is not role=value, such as verify=sonnet", flag, p)
			}
			r := out[n]
			put(&r, strings.TrimSpace(v))
			out[n] = r
		}
		return nil
	}
	if err := set("role harness", harnesses, func(r *Role, v string) { r.Harness = v }); err != nil {
		return nil, err
	}
	if err := set("role model", models, func(r *Role, v string) { r.Model = v }); err != nil {
		return nil, err
	}
	for _, p := range configs {
		nk, v, ok := strings.Cut(p, "=")
		n, k, dotted := strings.Cut(strings.TrimSpace(nk), ".")
		if !ok || !dotted || n == "" || k == "" {
			return nil, fmt.Errorf("agent role config %q is not role.key=value, such as verify.effort=medium", p)
		}
		r := out[n]
		if r.Config == nil {
			r.Config = map[string]string{}
		}
		r.Config[k] = strings.TrimSpace(v)
		out[n] = r
	}
	if len(out) == 0 {
		return nil, nil
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
	fields(&b, indent+"  ", a.Harness, a.Model, a.Config)
	if len(a.Roles) > 0 {
		fmt.Fprintf(&b, "%s  roles:\n", indent)
		for _, n := range a.RoleNames() {
			r := a.Roles[n]
			fmt.Fprintf(&b, "%s    %s:\n", indent, n)
			fields(&b, indent+"      ", r.Harness, r.Model, r.Config)
		}
	}
	return b.String()
}

func fields(b *strings.Builder, indent, harness, model string, config map[string]string) {
	if harness != "" {
		fmt.Fprintf(b, "%sharness: %s\n", indent, yamlScalar(harness))
	}
	if model != "" {
		fmt.Fprintf(b, "%smodel: %s\n", indent, yamlScalar(model))
	}
	if len(config) > 0 {
		fmt.Fprintf(b, "%sconfig:\n", indent)
		for _, k := range sortedKeys(config) {
			fmt.Fprintf(b, "%s  %s: %s\n", indent, k, yamlScalar(config[k]))
		}
	}
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
