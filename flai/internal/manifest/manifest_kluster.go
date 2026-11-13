package manifest

import (
	"fmt"
	"io/fs"
	"maps"
	"math"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/goccy/go-yaml"
)

// Kluster is the manifest's kluster: block, the project's one local
// Kubernetes cluster made by bytepunx/kluster and the project's own commands
// that build, deploy into, and tear down a namespace of it
// (design/system/kluster.md, ADR-0254). Each command is an argument list,
// run as it stands, never through a shell.
type Kluster struct {
	// Config is the kluster config file, from the repository root. Required.
	Config string `yaml:"config" json:"config"`
	// Namespace is the default namespace, which main deploys into; empty
	// means the manifest's name as a DNS label (Manifest.KlusterNamespace).
	Namespace string `yaml:"namespace,omitempty" json:"namespace,omitempty"`
	// Build builds the project's images and lists their references in the
	// file FLAI_KLUSTER_IMAGES names; empty loads nothing.
	Build []string `yaml:"build,omitempty" json:"build,omitempty"`
	// Components are the project's components, each built by a command of
	// its own as a step of its own, in this order, in place of Build
	// (ADR-0263); a block names Build or Components, not both.
	Components []KlusterComponent `yaml:"components,omitempty" json:"components,omitempty"`
	// Deploy makes a namespace's settings in the shared services, applies the
	// project's manifests into it, and waits for their rollouts; empty means
	// the project's deploys are not set up yet, and flai deploys nothing.
	Deploy []string `yaml:"deploy,omitempty" json:"deploy,omitempty"`
	// Teardown removes what Deploy made outside the namespace; empty means
	// flai deletes only the namespace.
	Teardown []string `yaml:"teardown,omitempty" json:"teardown,omitempty"`
	// Timeout is the most one job may take, a Go duration; empty means
	// DefaultKlusterTimeout.
	Timeout string `yaml:"timeout,omitempty" json:"timeout,omitempty"`
	// AgentNamespaces is the most stories in progress that hold a namespace
	// at once from deploys their agents asked for, a whole number of 0 or
	// more, where 0 turns deploys before review off; stories in review do not
	// count (ADR-0292 §4). Nil means DefaultKlusterAgentNamespaces
	// (Manifest.KlusterAgentNamespaces).
	AgentNamespaces *int `yaml:"agent_namespaces,omitempty" json:"agent_namespaces,omitempty"`

	// problems are what the block, as written, says that its fields cannot
	// hold: a key that names nothing, or a value of the wrong kind, one
	// sentence each, for KlusterErrors.
	problems []string
}

// KlusterComponent is one component of the kluster: block: a name, a DNS
// label unique in the list, and the command that builds its images and lists
// their references in the file FLAI_KLUSTER_IMAGES names (ADR-0263).
type KlusterComponent struct {
	// Name names the component and its build step, build-<name>.
	Name string `yaml:"name" json:"name"`
	// Build is the argument list that builds the component, run as it stands.
	Build []string `yaml:"build" json:"build"`
}

// The keys of the kluster: block, in the order they are listed to the
// operator.
var klusterKeys = []string{"config", "namespace", "build", "components", "deploy", "teardown", "timeout", "agent_namespaces"}

// The keys of a component of the kluster: block.
var klusterComponentKeys = []string{"name", "build"}

// DefaultKlusterTimeout is the most one kluster job may take when
// kluster.timeout is not set.
const DefaultKlusterTimeout = 30 * time.Minute

// DefaultKlusterAgentNamespaces is the most stories in progress that hold a
// namespace at once when kluster.agent_namespaces is not set (ADR-0292 §4).
const DefaultKlusterAgentNamespaces = 2

// DefaultKlusterProvider is the provider a kluster config names when it names
// none, and the only one flai runs (ADR-0254 §1).
const DefaultKlusterProvider = "k3d"

// UnmarshalYAML reads the block leniently, so that a misspelt key or a value
// of the wrong kind is a manifest.kluster finding rather than a manifest
// that does not load.
func (k *Kluster) UnmarshalYAML(unmarshal func(any) error) error {
	var raw any
	if err := unmarshal(&raw); err != nil {
		return err
	}
	*k = Kluster{}
	block, ok := raw.(map[string]any)
	if !ok {
		k.problems = append(k.problems, "kluster is not a block of keys; write it with config: under it, and deploy: once the project deploys, as design/system/kluster.md shows")
		return nil
	}
	text := func(key string, into *string) {
		switch v := block[key].(type) {
		case nil:
		case string:
			*into = v
		default:
			k.problems = append(k.problems, fmt.Sprintf("kluster.%s %v is not text; write it as a string", key, v))
		}
	}
	argv := func(field, example string, value any, into *[]string) {
		switch v := value.(type) {
		case nil:
		case []any:
			for i, a := range v {
				s, ok := a.(string)
				if !ok {
					k.problems = append(k.problems, fmt.Sprintf("%s argument %d, %v, is not text; write every argument as a string", field, i+1, a))
					s = fmt.Sprint(a)
				}
				*into = append(*into, s)
			}
		default:
			k.problems = append(k.problems, fmt.Sprintf("%s is not an argument list; write the program, then its arguments, in brackets, such as [%s]", field, example))
		}
	}
	command := func(key string, into *[]string) {
		argv("kluster."+key, "scripts/kluster-"+key+".sh", block[key], into)
	}
	text("config", &k.Config)
	text("namespace", &k.Namespace)
	command("build", &k.Build)
	k.components(block["components"], argv)
	command("deploy", &k.Deploy)
	command("teardown", &k.Teardown)
	switch v := block["timeout"].(type) {
	case nil:
	case string:
		k.Timeout = v
	default:
		k.problems = append(k.problems, fmt.Sprintf("kluster.timeout %v is not a duration; write one with its unit, such as 30m or 1h", v))
	}
	k.agentNamespaces(block["agent_namespaces"])
	for _, key := range slices.Sorted(maps.Keys(block)) {
		if !slices.Contains(klusterKeys, key) {
			k.problems = append(k.problems, fmt.Sprintf("kluster has no key %q; write one of %s, or remove it", key, strings.Join(klusterKeys, ", ")))
		}
	}
	return nil
}

// components reads kluster.components leniently into k.Components, adding a
// problem for each part not of its kind; argv reads a command as
// UnmarshalYAML's own commands are read.
func (k *Kluster) components(value any, argv func(field, example string, value any, into *[]string)) {
	if value == nil {
		return
	}
	list, ok := value.([]any)
	if !ok {
		k.problems = append(k.problems, "kluster.components is not a list; write one item per component, each with name: and build:, such as - {name: api, build: [scripts/build.sh, api]}")
		return
	}
	k.Components = []KlusterComponent{}
	for i, item := range list {
		field := fmt.Sprintf("kluster.components[%d]", i+1)
		block, ok := item.(map[string]any)
		if !ok {
			k.problems = append(k.problems, fmt.Sprintf("%s is not a block of keys; write it with name: and build:, such as {name: api, build: [scripts/build.sh, api]}", field))
			continue
		}
		var c KlusterComponent
		switch v := block["name"].(type) {
		case nil:
		case string:
			c.Name = v
		default:
			k.problems = append(k.problems, fmt.Sprintf("%s.name %v is not text; write it as a string", field, v))
		}
		argv(field+".build", "scripts/build.sh", block["build"], &c.Build)
		for _, key := range slices.Sorted(maps.Keys(block)) {
			if !slices.Contains(klusterComponentKeys, key) {
				k.problems = append(k.problems, fmt.Sprintf("%s has no key %q; write one of %s, or remove it", field, key, strings.Join(klusterComponentKeys, ", ")))
			}
		}
		k.Components = append(k.Components, c)
	}
}

// agentNamespacesAdvice is what each kluster.agent_namespaces finding asks
// for.
var agentNamespacesAdvice = fmt.Sprintf("write the most stories in progress that may hold a namespace at once, such as %d, 0 to turn deploys before review off, or remove it for %d", DefaultKlusterAgentNamespaces, DefaultKlusterAgentNamespaces)

// agentNamespaces reads kluster.agent_namespaces leniently into
// k.AgentNamespaces: a value that is not a whole number, such as 1.5 or text,
// or one too large to count, adds a problem instead. Whether it is 0 or
// more is KlusterErrors' to say.
func (k *Kluster) agentNamespaces(value any) {
	var n int
	switch v := value.(type) {
	case nil:
		return
	case int:
		n = v
	case int64:
		n = int(v)
	case uint64:
		if v > math.MaxInt {
			k.problems = append(k.problems, fmt.Sprintf("kluster.agent_namespaces %d is too large to count; %s", v, agentNamespacesAdvice))
			return
		}
		n = int(v)
	case float64:
		switch {
		case v != math.Trunc(v):
			k.problems = append(k.problems, fmt.Sprintf("kluster.agent_namespaces %s is not a whole number; %s", strconv.FormatFloat(v, 'f', -1, 64), agentNamespacesAdvice))
			return
		case v < math.MinInt || v > math.MaxInt:
			k.problems = append(k.problems, fmt.Sprintf("kluster.agent_namespaces %s is too large to count; %s", strconv.FormatFloat(v, 'g', -1, 64), agentNamespacesAdvice))
			return
		}
		n = int(v)
	case string:
		k.problems = append(k.problems, fmt.Sprintf("kluster.agent_namespaces %q is not a whole number; %s", v, agentNamespacesAdvice))
		return
	default:
		k.problems = append(k.problems, fmt.Sprintf("kluster.agent_namespaces %v is not a whole number; %s", v, agentNamespacesAdvice))
		return
	}
	k.AgentNamespaces = &n
}

// KlusterAgentNamespaces is the most stories in progress that may hold a
// kluster namespace at once from deploys their agents asked for
// (ADR-0292 §4): kluster.agent_namespaces, or DefaultKlusterAgentNamespaces
// when it is not set or not valid, which KlusterErrors reports; 0 means
// deploys before review are off. Zero when the manifest has no kluster:
// block, as there is no cluster to deploy into.
func (m Manifest) KlusterAgentNamespaces() int {
	if m.Kluster == nil {
		return 0
	}
	if n := m.Kluster.AgentNamespaces; n != nil && *n >= 0 {
		return *n
	}
	return DefaultKlusterAgentNamespaces
}

// TimeoutDuration is kluster.timeout as a duration, DefaultKlusterTimeout
// when it is empty.
func (k Kluster) TimeoutDuration() (time.Duration, error) {
	s := strings.TrimSpace(k.Timeout)
	if s == "" {
		return DefaultKlusterTimeout, nil
	}
	d, err := time.ParseDuration(s)
	if err != nil {
		return 0, fmt.Errorf("kluster.timeout %q is not a duration; write the most one job may take, such as 30m or 1h, or remove it for 30m", k.Timeout)
	}
	if d <= 0 {
		return 0, fmt.Errorf("kluster.timeout %q is not above zero; write the most one job may take, such as 30m or 1h, or remove it for 30m", k.Timeout)
	}
	return d, nil
}

// KlusterTimeout is the most one kluster job may take: kluster.timeout, or
// DefaultKlusterTimeout when it is not set or not valid, which KlusterErrors
// reports. Zero when the manifest has no kluster: block.
func (m Manifest) KlusterTimeout() time.Duration {
	if m.Kluster == nil {
		return 0
	}
	d, err := m.Kluster.TimeoutDuration()
	if err != nil {
		return DefaultKlusterTimeout
	}
	return d
}

// KlusterNamespace is the default namespace main deploys into:
// kluster.namespace, or the manifest's name as a DNS label when it is not
// set. Empty when the manifest has no kluster: block, or when its name makes
// no DNS label, which KlusterErrors reports.
func (m Manifest) KlusterNamespace() string {
	if m.Kluster == nil {
		return ""
	}
	if ns := strings.TrimSpace(m.Kluster.Namespace); ns != "" {
		return ns
	}
	return DNSLabel(m.Name)
}

// DNSLabel is name made a DNS label (RFC 1123): lower case, each run of
// characters other than letters and digits made one hyphen, with none at
// either end, and at most 63 characters. Empty when name holds no letter or
// digit.
func DNSLabel(name string) string {
	var b strings.Builder
	hyphen := false
	for _, r := range strings.ToLower(name) {
		if r >= 'a' && r <= 'z' || r >= '0' && r <= '9' {
			if hyphen && b.Len() > 0 {
				b.WriteByte('-')
			}
			b.WriteRune(r)
			hyphen = false
			continue
		}
		hyphen = true
	}
	s := b.String()
	if len(s) > 63 {
		s = strings.TrimRight(s[:63], "-")
	}
	return s
}

var (
	dnsLabel       = regexp.MustCompile(`^[a-z0-9]([-a-z0-9]*[a-z0-9])?$`)
	storyNamespace = regexp.MustCompile(`^s-[0-9]+$`)
)

// KlusterConfig is what flai reads of a kluster config file: the cluster's
// name, its provider, and the profile it installs.
type KlusterConfig struct {
	// Name is the cluster's name; its kube context is k3d-<name>.
	Name string `yaml:"name" json:"name"`
	// Provider is the cluster's provider, DefaultKlusterProvider when the
	// file names none.
	Provider string `yaml:"provider" json:"provider"`
	// Profile is the profile kluster up activates, the shared
	// infrastructure; empty when the file names none.
	Profile string `yaml:"profile" json:"profile"`
}

// ReadKlusterConfig reads the kluster config file at name in fsys, such as
// the kluster.config of a manifest in os.DirFS of the repository root.
func ReadKlusterConfig(fsys fs.FS, name string) (KlusterConfig, error) {
	data, err := fs.ReadFile(fsys, name)
	if err != nil {
		return KlusterConfig{}, err
	}
	var c KlusterConfig
	if err := yaml.Unmarshal(data, &c); err != nil {
		return KlusterConfig{}, fmt.Errorf("parse %s: %w", name, err)
	}
	c.Name, c.Provider, c.Profile = strings.TrimSpace(c.Name), strings.TrimSpace(c.Provider), strings.TrimSpace(c.Profile)
	if c.Provider == "" {
		c.Provider = DefaultKlusterProvider
	}
	return c, nil
}

// KlusterErrors are what is wrong with the kluster: block, one sentence each
// beginning with the key it is about; none when the block is valid or the
// manifest has none. fsys is the repository, where the config file is read;
// nil checks the config's path but not the file. A block with no deploy is
// valid: its cluster can come up and go down, and nothing is deployed until
// kluster.deploy is set.
func (m Manifest) KlusterErrors(fsys fs.FS) []string {
	k := m.Kluster
	if k == nil {
		return nil
	}
	errs := slices.Clone(k.problems)
	switch cfg := strings.TrimSpace(k.Config); {
	case cfg == "":
		errs = append(errs, "kluster.config is required; write the kluster config file, from the repository root, such as integration/kluster.yaml")
	case repoPathProblem(cfg, "a file", "integration/kluster.yaml", "integration/kluster.yaml") != "":
		errs = append(errs, fmt.Sprintf("kluster.config %q %s", k.Config, repoPathProblem(cfg, "a file", "integration/kluster.yaml", "integration/kluster.yaml")))
	case fsys == nil:
		// no repository to read the config file in: its path alone is checked
	case !isFile(fsys, cfg):
		errs = append(errs, fmt.Sprintf("kluster.config %q does not exist; write the path of the kluster config file, from the repository root, or add the file", k.Config))
	default:
		c, err := ReadKlusterConfig(fsys, cfg)
		switch {
		case err != nil:
			errs = append(errs, fmt.Sprintf("kluster.config %q does not read as a kluster config: %v", k.Config, err))
		case c.Name == "":
			errs = append(errs, fmt.Sprintf("kluster.config %q names no cluster; give the file a name:, which flai runs the cluster by", k.Config))
		case c.Provider != DefaultKlusterProvider:
			errs = append(errs, fmt.Sprintf("kluster.config %q names the provider %q; flai runs only %s, so write provider: %s or remove it", k.Config, c.Provider, DefaultKlusterProvider, DefaultKlusterProvider))
		}
	}
	if ns := strings.TrimSpace(k.Namespace); ns != "" {
		errs = append(errs, namespaceProblems("kluster.namespace", fmt.Sprintf("%q", k.Namespace), ns)...)
	} else if ns := DNSLabel(m.Name); ns == "" {
		errs = append(errs, fmt.Sprintf("kluster.namespace is not set, and the manifest's name %q makes no DNS label; set kluster.namespace, such as app", m.Name))
	} else {
		errs = append(errs, namespaceProblems("kluster.namespace", fmt.Sprintf("is not set, and the manifest's name makes %q, which", ns), ns)...)
	}
	for _, c := range []struct {
		key  string
		args []string
	}{{"build", k.Build}, {"deploy", k.Deploy}, {"teardown", k.Teardown}} {
		if len(c.args) > 0 && strings.TrimSpace(c.args[0]) == "" {
			errs = append(errs, fmt.Sprintf("kluster.%s has no program first; write the program, then its arguments, such as [scripts/kluster-%s.sh]", c.key, c.key))
		}
	}
	errs = append(errs, k.componentErrors()...)
	if _, err := k.TimeoutDuration(); err != nil {
		errs = append(errs, err.Error())
	}
	if n := k.AgentNamespaces; n != nil && *n < 0 {
		errs = append(errs, fmt.Sprintf("kluster.agent_namespaces %d is not a whole number of 0 or more; %s", *n, agentNamespacesAdvice))
	}
	return errs
}

// componentErrors are what is wrong with kluster.components, one sentence
// each beginning with the key: a block that names build as well, a
// component with no name, a repeated name or one that is not a DNS label,
// and a component whose build is missing or has no program first.
func (k *Kluster) componentErrors() []string {
	if len(k.Components) == 0 {
		return nil
	}
	var errs []string
	if len(k.Build) > 0 {
		errs = append(errs, "kluster.components and kluster.build are both set; a block builds with one command or one per component, so remove one of them")
	}
	seen := map[string]int{}
	for i, c := range k.Components {
		field := fmt.Sprintf("kluster.components[%d]", i+1)
		switch name := strings.TrimSpace(c.Name); {
		case name == "":
			errs = append(errs, field+" has no name; give it one, a DNS label such as api, which names its build step build-<name>")
		case len(name) > 63 || !dnsLabel.MatchString(name):
			errs = append(errs, fmt.Sprintf("%s.name %q is not a DNS label; write lower-case letters, digits, and hyphens, at most 63, beginning and ending with a letter or digit, such as api", field, c.Name))
		case seen[name] > 0:
			errs = append(errs, fmt.Sprintf("%s.name %q is the name of kluster.components[%d] too; give each component a name of its own", field, c.Name, seen[name]))
		default:
			seen[name] = i + 1
		}
		switch {
		case len(c.Build) == 0:
			errs = append(errs, field+".build is required; write the command that builds the component, the program then its arguments, such as [scripts/build.sh, api]")
		case strings.TrimSpace(c.Build[0]) == "":
			errs = append(errs, field+".build has no program first; write the program, then its arguments, such as [scripts/build.sh, api]")
		}
	}
	return errs
}

// namespaceProblems say why ns, the default namespace, cannot be one, each
// sentence beginning with key and then said, which names the value.
func namespaceProblems(key, said, ns string) []string {
	switch {
	case len(ns) > 63 || !dnsLabel.MatchString(ns):
		return []string{fmt.Sprintf("%s %s is not a DNS label; write lower-case letters, digits, and hyphens, at most 63, beginning and ending with a letter or digit, such as app", key, said)}
	case storyNamespace.MatchString(ns):
		return []string{fmt.Sprintf("%s %s has a story namespace's form, s- and digits, which each story in review deploys into; write another, such as app", key, said)}
	}
	return nil
}
