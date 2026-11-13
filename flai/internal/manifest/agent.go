package manifest

import (
	"errors"
	"fmt"
	"maps"
	"os"
	"path"
	"regexp"
	"sort"
	"strings"

	"github.com/goccy/go-yaml"
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
// its own; the fields above are the story's own. Provider names the host's
// provider the model's calls go to (ADR-0129); empty leaves the harness's
// own default. Image is the project's own image for its agents (S-0510): it
// belongs to the project's agent alone, and With never carries it into a
// story's agent or a strategic role's.
type Agent struct {
	Harness  string            `yaml:"harness,omitempty" json:"harness,omitempty"`
	Model    string            `yaml:"model,omitempty" json:"model,omitempty"`
	Provider string            `yaml:"provider,omitempty" json:"provider,omitempty"`
	Config   map[string]string `yaml:"config,omitempty" json:"config,omitempty"`
	Roles    map[string]Role   `yaml:"roles,omitempty" json:"roles,omitempty"`
	Image    *AgentImage       `yaml:"image,omitempty" json:"image,omitempty"`
}

// AgentImage is the project's say about the image its agents run in
// (S-0510, ADR-0203's one layer a project may add): Layer is a Dockerfile,
// relative to the repository root, built over the harness's image with the
// tools the project's test tiers need. Args are literal values for the
// layer's build arguments, each over the default its ARG gives in the
// Dockerfile, which keeps the pins (ADR-0295); flai sets BASE itself.
type AgentImage struct {
	Layer string            `yaml:"layer,omitempty" json:"layer,omitempty"`
	Args  map[string]string `yaml:"args,omitempty" json:"args,omitempty"`
}

// exampleLayer is the Dockerfile a problem with agent.image.layer suggests.
const exampleLayer = "images/agent/Dockerfile"

// BaseArg is the layer's build argument flai sets to the harness image it
// builds the layer over, which agent.image.args may not name.
const BaseArg = "BASE"

// argPattern is a Dockerfile ARG name.
var argPattern = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)

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
	return a == nil || a.Harness == "" && a.Model == "" && a.Provider == "" && len(a.Config) == 0 && len(a.Roles) == 0 && a.Image == nil
}

// ConfigKeys are the role's config keys in order.
func (r Role) ConfigKeys() []string {
	return sortedKeys(r.Config)
}

// IsZero reports whether the role sets nothing.
func (r Role) IsZero() bool {
	return r.Harness == "" && r.Model == "" && len(r.Config) == 0
}

// Validate checks each part's shape: an identifier for the harness and the
// provider, a model ID for the model, config keys that are identifiers with
// values on one line, and role names that are identifiers, each role setting
// something and checked the same way.
func (a *Agent) Validate() error {
	if errs := a.Problems(); len(errs) > 0 {
		return errors.New(strings.Join(errs, "; "))
	}
	return nil
}

// Problems are what Validate finds, one phrase each, for a caller that reports
// each on its own.
func (a *Agent) Problems() []string {
	return a.problems("agent")
}

// problems are what Validate finds, one phrase each, each naming the block
// as what: agent, or planning.agent.
func (a *Agent) problems(what string) []string {
	if a == nil {
		return nil
	}
	errs := validate(what, a.Harness, a.Model, a.Config)
	if a.Provider != "" && !harnessPattern.MatchString(a.Provider) {
		errs = append(errs, fmt.Sprintf("%s provider %q is not a name such as litellm or openrouter (lower case letters, digits, . _ -)", what, a.Provider))
	}
	for _, n := range a.RoleNames() {
		r := a.Roles[n]
		switch {
		case !rolePattern.MatchString(n):
			errs = append(errs, fmt.Sprintf("%s role %q is not a name such as explore or verify (lower case letters, digits, _ -)", what, n))
		case r.IsZero():
			errs = append(errs, fmt.Sprintf("%s role %s sets nothing; give it a harness, a model, or config, or leave it out", what, n))
		default:
			errs = append(errs, validate(what+" role "+n, r.Harness, r.Model, r.Config)...)
		}
	}
	return append(errs, a.Image.problems(what)...)
}

// problems are what is wrong with the image under the block what: an image
// only the project's agent names, which names a layer that is a clean path
// inside the repository on one line, and args only beside a layer, each an
// ARG name other than BASE with a value on one line.
func (i *AgentImage) problems(what string) []string {
	switch {
	case i == nil:
		return nil
	case what != "agent":
		return []string{fmt.Sprintf("%s image is the project's alone; name the layer under agent.image in %s, or leave it out here", what, File)}
	case i.Layer == "" && len(i.Args) == 0:
		return []string{fmt.Sprintf("%s image sets nothing; name a Dockerfile with layer, such as %s, or leave image out", what, exampleLayer)}
	case i.Layer == "":
		return []string{fmt.Sprintf("%s image args name no layer to build; name its Dockerfile with %s.image.layer, such as %s, or leave %s.image.args out", what, what, exampleLayer, what)}
	}
	return append(i.layerProblems(what), i.argProblems(what)...)
}

// layerProblems are what is wrong with the layer's path under the block what.
func (i *AgentImage) layerProblems(what string) []string {
	switch {
	case strings.ContainsAny(i.Layer, "\n\r"):
		return []string{fmt.Sprintf("%s image layer spans lines; a path is one line", what)}
	case strings.Contains(i.Layer, `\`):
		return []string{fmt.Sprintf("%s image layer %q is not a clean path; write it with forward slashes, such as %s", what, i.Layer, exampleLayer)}
	}
	if r := repoPathProblem(i.Layer, "a Dockerfile", exampleLayer, path.Clean(i.Layer)); r != "" {
		return []string{fmt.Sprintf("%s image layer %q %s", what, i.Layer, r)}
	}
	return nil
}

// argProblems are what is wrong with the args under the block what, by name.
func (i *AgentImage) argProblems(what string) []string {
	var errs []string
	for _, n := range sortedKeys(i.Args) {
		switch {
		case !argPattern.MatchString(n):
			errs = append(errs, fmt.Sprintf("%s image args name %q is not a Dockerfile ARG name such as GO_VERSION (letters, digits, _, not starting with a digit); write it in %s.image.args as the layer's ARG names it", what, n, what))
		case n == BaseArg:
			errs = append(errs, fmt.Sprintf("%s image args set %s, which flai sets to the harness image it builds the layer over; leave %s out of %s.image.args", what, BaseArg, BaseArg, what))
		case strings.ContainsAny(i.Args[n], "\n\r"):
			errs = append(errs, fmt.Sprintf("%s image args %s spans lines; write its value in %s.image.args as one literal line", what, n, what))
		}
	}
	return errs
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
	if a.Harness != b.Harness || a.Model != b.Model || a.Provider != b.Provider || !sameConfig(a.Config, b.Config) || len(a.Roles) != len(b.Roles) || !sameImage(a.Image, b.Image) {
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

func sameImage(a, b *AgentImage) bool {
	if a == nil || b == nil {
		return a == nil && b == nil
	}
	return a.Layer == b.Layer && maps.Equal(a.Args, b.Args)
}

// clone is a copy of i that shares no map with it; nil for nil.
func (i *AgentImage) clone() *AgentImage {
	if i == nil {
		return nil
	}
	c := *i
	c.Args = maps.Clone(i.Args)
	return &c
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
	c.Image = a.Image.clone()
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

// With is a with over's harness, model, and provider where over sets them,
// and their configs merged key by key, over's winning; each role is merged
// the same way, role by role. Neither side's image is carried: it is the
// project's alone. Nil when the result is empty.
func (a *Agent) With(over *Agent) *Agent {
	out := &Agent{}
	for _, x := range []*Agent{a, over} {
		if x == nil {
			continue
		}
		out.Harness, out.Model, out.Config = merge(out.Harness, out.Model, out.Config, x.Harness, x.Model, x.Config)
		if x.Provider != "" {
			out.Provider = x.Provider
		}
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

// BuiltinHarness is the harness of flai's built-in default agent
// (ADR-0261): the claude-code adapter's name.
const BuiltinHarness = "claude-code"

// Builtin is flai's built-in default agent (ADR-0261): harness claude-code,
// no model, so the harness runs its own default model, and no config. An
// agent that names no harness, under neither a story, a role, the project's
// default, nor the host's command, is started with it. It is never written
// into a story or a manifest.
func Builtin() *Agent {
	return &Agent{Harness: BuiltinHarness}
}

// OrBuiltin is the agent itself when it names a harness, and otherwise flai's
// built-in default with what a does say, model, provider, config, and roles,
// laid over it, and a's image kept (ADR-0261). a is not changed.
func (a *Agent) OrBuiltin() *Agent {
	if a != nil && a.Harness != "" {
		return a
	}
	out := Builtin().With(a)
	if a != nil {
		out.Image = a.Image.clone()
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
// claude-opus-5-5, provider litellm, effort=high; explore: haiku; verify:
// sonnet.
func (a *Agent) String() string {
	if a.IsZero() {
		return "none"
	}
	provider := ""
	if a.Provider != "" {
		provider = "provider " + a.Provider
	}
	parts := []string{line(a.Config, a.Harness, a.Model, provider)}
	if parts[0] == "" {
		parts[0] = "default"
	}
	for _, n := range a.RoleNames() {
		r := a.Roles[n]
		parts = append(parts, n+": "+line(r.Config, r.Harness, r.Model))
	}
	return strings.Join(parts, "; ")
}

// line is the words that are set, then config's pairs, comma separated.
func line(config map[string]string, words ...string) string {
	var parts []string
	for _, w := range words {
		if w != "" {
			parts = append(parts, w)
		}
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
// this key, and the rest of the file is the operator's. An image a does not
// name is kept as the file names it, since setting or clearing the agent that
// works the stories leaves the project's image alone (S-0510).
func WriteAgent(path string, a *Agent) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	if a == nil || a.Image == nil {
		var was struct {
			Agent *Agent `yaml:"agent"`
		}
		if err := yaml.Unmarshal(data, &was); err != nil {
			return fmt.Errorf("read the agent image %s names: %w; fix its YAML by hand, then run this again", path, err)
		}
		if was.Agent != nil && was.Agent.Image != nil {
			a = a.Clone()
			if a == nil {
				a = &Agent{}
			}
			a.Image = was.Agent.Image
		}
	}
	if err := a.Validate(); err != nil {
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
	fields(&b, indent+"  ", a.Harness, a.Model, a.Provider, a.Config)
	if len(a.Roles) > 0 {
		fmt.Fprintf(&b, "%s  roles:\n", indent)
		for _, n := range a.RoleNames() {
			r := a.Roles[n]
			fmt.Fprintf(&b, "%s    %s:\n", indent, n)
			fields(&b, indent+"      ", r.Harness, r.Model, "", r.Config)
		}
	}
	if a.Image != nil {
		fmt.Fprintf(&b, "%s  image:\n", indent)
		if a.Image.Layer != "" {
			fmt.Fprintf(&b, "%s    layer: %s\n", indent, yamlScalar(a.Image.Layer))
		}
		if len(a.Image.Args) > 0 {
			fmt.Fprintf(&b, "%s    args:\n", indent)
			for _, n := range sortedKeys(a.Image.Args) {
				fmt.Fprintf(&b, "%s      %s: %s\n", indent, n, yamlScalar(a.Image.Args[n]))
			}
		}
	}
	return b.String()
}

func fields(b *strings.Builder, indent, harness, model, provider string, config map[string]string) {
	if harness != "" {
		fmt.Fprintf(b, "%sharness: %s\n", indent, yamlScalar(harness))
	}
	if model != "" {
		fmt.Fprintf(b, "%smodel: %s\n", indent, yamlScalar(model))
	}
	if provider != "" {
		fmt.Fprintf(b, "%sprovider: %s\n", indent, yamlScalar(provider))
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
