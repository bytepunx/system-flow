// Package manifest reads system-flow.yaml, the file that marks a conforming
// project. See design/system/project-manifest.md and ADR-0011.
package manifest

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"

	"github.com/goccy/go-yaml"

	"github.com/bytepunx/system-flow/flai/internal/buildinfo"
)

// File is the manifest file name at a project root.
const File = "system-flow.yaml"

// ErrNotFound means no manifest exists in the directory or any parent.
var ErrNotFound = errors.New("no " + File + " found in this directory or any parent")

// Manifest is the schema of system-flow.yaml.
type Manifest struct {
	Version     int               `yaml:"version"`
	Name        string            `yaml:"name"`
	Key         string            `yaml:"key"`
	Description string            `yaml:"description"`
	Owner       string            `yaml:"owner"`
	Repo        string            `yaml:"repo"`
	Template    Template          `yaml:"template"`
	Layout      map[string]string `yaml:"layout"`
	Projects    []Project         `yaml:"projects"`
	Dashboard   Dashboard         `yaml:"dashboard"`
	// Checks are the commands a story in review is checked with, in the
	// story's worktree, in order (S-0082). Used only when the host's own
	// configuration names none: that is the operator's choice of where to
	// name them, not a merge of both.
	Checks []NamedCommand `yaml:"checks,omitempty" json:"checks,omitempty"`
	// Tests are the project's test tiers, cheapest first, that flai test runs
	// for the paths each selects (S-0273). Nil, with no tests key, means the
	// default TestTiers gives; an empty list means none.
	Tests []TestTier `yaml:"tests,omitempty" json:"tests,omitempty"`
	// Agent is the project's default agent, copied into every story created
	// while it is set (S-0103). flai agent sets it.
	Agent *Agent `yaml:"agent,omitempty" json:"agent,omitempty"`
	// Providers override, by name, the api, base_url, and models of the
	// host's providers an agent may name (ADR-0129); key_env is the host's
	// alone, and ProviderErrors refuses it here.
	Providers map[string]Provider `yaml:"providers,omitempty" json:"providers,omitempty"`
	// Prime is how flai prime --story builds this project's context packs
	// (ADR-0049).
	Prime Prime `yaml:"prime,omitempty" json:"prime,omitzero"`
	// Issues is how flai check treats the project's open issues (S-0198).
	Issues Issues `yaml:"issues,omitempty" json:"issues,omitzero"`
	// Planning is the units the planner's numbers are in (S-0199), and the
	// planner's agent (S-0208).
	Planning Planning `yaml:"planning,omitempty" json:"planning,omitzero"`
	// Orchestration is the policy that orders the ready column and the
	// policy that says when to release (S-0217).
	Orchestration Orchestration `yaml:"orchestration,omitempty" json:"orchestration,omitzero"`
	// Analysis is the analyzer's agent and schedule (S-0223).
	Analysis Analysis `yaml:"analysis,omitempty" json:"analysis,omitzero"`
	// Digest is the curator's agent (ADR-0196), schedule, and compaction
	// window (ADR-0199).
	Digest Digest `yaml:"digest,omitempty" json:"digest,omitzero"`
	// Threads is when flai serve reminds the operator of a thread that
	// awaits them (ADR-0205).
	Threads Threads `yaml:"threads,omitempty" json:"threads,omitzero"`
	// Roles are the operator's settings for each role a role file defines, by
	// the role's name (ADR-0207); flai guard reads a role's permissions from
	// it at each call (S-0416). The orchestrator's stay in Orchestration.
	Roles map[string]RoleSettings `yaml:"roles,omitempty" json:"roles,omitempty"`
	// Acceptance is how a story in review is accepted: on this machine, or
	// through a pull request on the project's forge (S-0392). Load refuses a
	// value Acceptance.Errors finds wrong.
	Acceptance Acceptance `yaml:"acceptance,omitempty" json:"acceptance,omitzero"`
	// MergeQueue is the merge queue's say about its batches (S-0399,
	// ADR-0173). Load refuses a value MergeQueue.Errors finds wrong.
	MergeQueue MergeQueue `yaml:"merge_queue,omitempty" json:"merge_queue,omitzero"`
	// Claims is how stories' touches claim paths: the shared paths whose
	// overlaps hold no story (S-0295, ADR-0096). A pattern that is not valid
	// does not stop the load; Claims.Errors names it, and it frees nothing.
	Claims Claims `yaml:"claims,omitempty" json:"claims,omitzero"`
	// Design is where the project's design system is kept (S-0381,
	// ADR-0146). Load refuses a path in it that is not inside the repository.
	Design Design `yaml:"design,omitempty" json:"design,omitzero"`
	// Flai is what the project asks of the flai that reads it (S-0181).
	Flai Requirement `yaml:"flai,omitempty" json:"flai,omitzero"`
	// Kluster names the project's local Kubernetes cluster and the commands
	// that deploy into it (S-0577, ADR-0254); nil when the manifest has no
	// kluster: block, and then flai runs no cluster for the project. A block
	// that is not valid does not stop the load; KlusterErrors names what is
	// wrong, as the check finding manifest.kluster.
	Kluster *Kluster `yaml:"kluster,omitempty" json:"kluster,omitempty"`
}

// Requirement is the oldest flai that may read the project.
type Requirement struct {
	// Minimum is a flai release, X.Y.Z: one that knows every front-matter
	// field the project's items carry. A publish raises it to the newest
	// flai release whose front-matter fields changed, no newer than the flai
	// doing the publish, so it names only a release that can be installed
	// (release.RaiseMinimum, I-0107).
	Minimum string `yaml:"minimum,omitempty" json:"minimum,omitempty"`
	// Command is the repository path of a program that runs the flai the
	// tree builds, such as scripts/flai.sh: task done's check step runs it in
	// a story's worktree in place of the running flai's own check, so the
	// branch's rules check the branch (ADR-0201). Only a project that builds
	// flai sets it.
	Command string `yaml:"command,omitempty" json:"command,omitempty"`
}

// TooOldError is a manifest whose minimum flai is newer than the running one.
type TooOldError struct {
	Path    string
	Minimum string
	Running string
}

func (e *TooOldError) Error() string {
	return fmt.Sprintf("%s needs flai %s or newer, and this is flai %s, which may not know the fields its items carry: upgrade it with %s, then run this again",
		e.Path, e.Minimum, e.Running, buildinfo.UpgradeCommand)
}

// Design is the project's say about where its design system is kept: one per
// project (ADR-0146).
type Design struct {
	// Tokens is the repository path of the project's design token file, a
	// DTCG file; empty means the default DesignTokens gives.
	Tokens string `yaml:"tokens,omitempty" json:"tokens,omitempty"`
	// Stylesheet is the repository path of the generated Tailwind v4
	// stylesheet flai design css --write writes the design tokens to, a file
	// under a sub-project (S-0385, ADR-0152, ADR-0159); empty means none, and
	// nothing is written or checked.
	Stylesheet string `yaml:"stylesheet,omitempty" json:"stylesheet,omitempty"`
}

// exampleDesignStylesheet is the stylesheet a refusal of design.stylesheet
// gives as one to write.
const exampleDesignStylesheet = "web/src/theme.css"

// DefaultDesignTokens is the design token file when design.tokens is not set
// and the design layout folder is design.
const DefaultDesignTokens = "design/" + designTokensUnder

// designTokensUnder is the default design token file under the design layout
// folder.
const designTokensUnder = "system/design-system.tokens.json"

// DesignTokens is the repository path of the design token file:
// design.tokens, or, when it is empty, design-system.tokens.json in the
// system folder of the design layout folder, DefaultDesignTokens under the
// default layout.
func (m Manifest) DesignTokens() string {
	if s := strings.TrimSpace(m.Design.Tokens); s != "" {
		return s
	}
	if d := m.Layout["design"]; d != "" {
		return path.Join(d, designTokensUnder)
	}
	return DefaultDesignTokens
}

// DesignStylesheet is the repository path of the generated design stylesheet,
// design.stylesheet, or "" when it is not set.
func (m Manifest) DesignStylesheet() string {
	return strings.TrimSpace(m.Design.Stylesheet)
}

// DesignStylesheetProject is the sub-project the design stylesheet lies under,
// the entry of projects whose path holds it; false when design.stylesheet is
// not set or lies outside every sub-project.
func (m Manifest) DesignStylesheetProject() (Project, bool) {
	s := m.DesignStylesheet()
	if s == "" {
		return Project{}, false
	}
	s = path.Clean(s)
	for _, p := range m.Projects {
		if d := path.Clean(strings.TrimSpace(p.Path)); strings.HasPrefix(s, d+"/") {
			return p, true
		}
	}
	return Project{}, false
}

// DesignGuidance is the repository path of the design system's guidance,
// design-system.md in the system folder of the design layout folder,
// design/system/design-system.md under the default layout (ADR-0146).
func (m Manifest) DesignGuidance() string {
	d := m.Layout["design"]
	if d == "" {
		d = "design"
	}
	return path.Join(d, "system/design-system.md")
}

// Template records which template produced the project.
type Template struct {
	Repo    string `yaml:"repo"`
	Ref     string `yaml:"ref"`
	Version string `yaml:"version"`
	Applied string `yaml:"applied"`
}

// Project is a releasable component at the repo root: a code sub-project
// or the template. Tags are story tags that mean "this story delivers to
// this component" (for example cli for flai).
type Project struct {
	Name string   `yaml:"name" json:"name"`
	Path string   `yaml:"path" json:"path"`
	Kind string   `yaml:"kind" json:"kind"`
	Tags []string `yaml:"tags,omitempty" json:"tags,omitempty"`
	// Components is where a Svelte sub-project's components are found and
	// what the harness imports before them (S-0437); nil means the defaults
	// ComponentSettings gives. Only a kind SupportsComponents takes it.
	Components *Components `yaml:"components,omitempty" json:"components,omitempty"`
	// Install is how the sub-project installs its dependencies, for flai
	// accept to run in the main checkout after an acceptance that changed
	// them (S-0483, ADR-0185); nil means flai installs nothing for it.
	Install *Install `yaml:"install,omitempty" json:"install,omitempty"`
}

// Install is a sub-project's install command and the paths whose change
// needs it (S-0483, ADR-0185). Its globs have the syntax of components'
// globs, relative to the sub-project's folder.
type Install struct {
	// Command is the program and its arguments, run with no shell in the
	// sub-project's folder.
	Command []string `yaml:"command" json:"command"`
	// Paths are the globs of the files whose change needs Command run.
	Paths []string `yaml:"paths" json:"paths"`
}

// exampleInstallPath is the glob a refusal of install.paths gives as one to
// write.
const exampleInstallPath = "package.json"

// InstallErrors are what is wrong with the sub-projects' install keys, one
// sentence each, each beginning with the sub-project's index and field, such
// as projects[1].install.command: a command that is empty or has no program
// first, and paths that are empty or hold a glob that is not valid. None when
// every key is valid.
func (m Manifest) InstallErrors() []string {
	var errs []string
	for i, p := range m.Projects {
		in := p.Install
		if in == nil {
			continue
		}
		at := func(field string) string { return projectField(i, p, "install", field) }
		switch {
		case len(in.Command) == 0:
			errs = append(errs, at("command")+" is empty; write the program and its arguments as a list, run in the sub-project's folder, such as [pnpm, install]")
		case strings.TrimSpace(in.Command[0]) == "":
			errs = append(errs, at("command")+" has no program first; write the program, then its arguments, such as [pnpm, install]")
		}
		if len(in.Paths) == 0 {
			errs = append(errs, at("paths")+" is empty; write the globs of the files whose change needs the install, relative to the sub-project's folder, such as "+exampleInstallPath)
		}
		for _, g := range in.Paths {
			if r := projectGlobProblem(g, exampleInstallPath); r != "" {
				errs = append(errs, fmt.Sprintf("%s pattern %q %s", at("paths"), g, r))
			}
		}
	}
	return errs
}

// InstallMatches reports whether the sub-project's install command needs to
// run for changed, paths relative to the repository root: whether one of
// them lies under the sub-project's folder and a valid glob of install.paths
// matches it there. False when the sub-project has no install key.
func (p Project) InstallMatches(changed []string) bool {
	if p.Install == nil {
		return false
	}
	dir := path.Clean(strings.TrimSpace(p.Path)) + "/"
	for _, c := range changed {
		rel, ok := strings.CutPrefix(cleanEntry(c), dir)
		if !ok || rel == "" {
			continue
		}
		segs := strings.Split(rel, "/")
		for _, g := range p.Install.Paths {
			if projectGlobProblem(g, exampleInstallPath) == "" && patternMatches(g, segs, false) {
				return true
			}
		}
	}
	return false
}

// The kinds of sub-project whose components flai reads (S-0437).
const (
	// ProjectKindSvelteKit is a SvelteKit app.
	ProjectKindSvelteKit = "sveltekit"
	// ProjectKindSvelte is a Svelte project that is not a SvelteKit app, such
	// as a component package built with Vite.
	ProjectKindSvelte = "svelte"
)

// ComponentKinds are the kinds of sub-project whose components flai reads,
// in the order they are listed to the operator.
var ComponentKinds = []string{ProjectKindSvelteKit, ProjectKindSvelte}

// DefaultComponentInclude is the glob, relative to the sub-project's folder,
// that finds its components when components.include is not set: every
// .svelte file under src/lib, so that routes are left out.
const DefaultComponentInclude = "src/lib/**/*.svelte"

// exampleComponentSetup is the setup a refusal of components.setup gives as
// one to write.
const exampleComponentSetup = "src/app.css"

// Components is a sub-project's say about its components (S-0437). Its globs
// have the syntax of claims.shared's, relative to the sub-project's folder.
type Components struct {
	// Include are the globs of the sub-project's component files; empty means
	// DefaultComponentInclude.
	Include []string `yaml:"include,omitempty" json:"include,omitempty"`
	// Exclude are the globs of files Include finds that are not components to
	// list; empty means none.
	Exclude []string `yaml:"exclude,omitempty" json:"exclude,omitempty"`
	// Setup is a module of the sub-project, relative to its folder, that the
	// harness imports before any component, such as the app's stylesheet;
	// empty means none.
	Setup string `yaml:"setup,omitempty" json:"setup,omitempty"`
}

// SupportsComponents reports whether flai reads the components of a
// sub-project of kind: one of ComponentKinds.
func SupportsComponents(kind string) bool {
	return slices.Contains(ComponentKinds, kind)
}

// ComponentSettings are the sub-project's components settings in effect:
// its components key over the defaults, include DefaultComponentInclude,
// nothing excluded, and no setup. False, with none, when its kind is not one
// SupportsComponents takes. It does not say whether the settings are valid;
// ComponentErrors does.
func (p Project) ComponentSettings() (Components, bool) {
	if !SupportsComponents(p.Kind) {
		return Components{}, false
	}
	c := Components{Include: []string{DefaultComponentInclude}}
	if p.Components == nil {
		return c, true
	}
	if len(p.Components.Include) > 0 {
		c.Include = slices.Clone(p.Components.Include)
	}
	c.Exclude = slices.Clone(p.Components.Exclude)
	c.Setup = strings.TrimSpace(p.Components.Setup)
	return c, true
}

// ComponentErrors are what is wrong with the sub-projects' components keys,
// one sentence each, each beginning with the sub-project's index and field,
// such as projects[1].components.include: a key on a sub-project of a kind
// SupportsComponents does not take, a glob that is not valid, and a setup
// that is not a clean path inside the sub-project or, when fsys, the
// repository root, is not nil, names no file under the sub-project's folder.
// None when every key is valid.
func (m Manifest) ComponentErrors(fsys fs.FS) []string {
	var errs []string
	for i, p := range m.Projects {
		c := p.Components
		if c == nil {
			continue
		}
		at := func(field string) string { return projectField(i, p, "components", field) }
		if !SupportsComponents(p.Kind) {
			errs = append(errs, fmt.Sprintf("%s is set on a sub-project of kind %q, whose components flai does not read; remove it, or set kind to %s if the sub-project is one", at(""), p.Kind, strings.Join(ComponentKinds, " or ")))
			continue
		}
		for _, list := range []struct {
			field string
			globs []string
		}{{"include", c.Include}, {"exclude", c.Exclude}} {
			for _, g := range list.globs {
				if r := projectGlobProblem(g, DefaultComponentInclude); r != "" {
					errs = append(errs, fmt.Sprintf("%s pattern %q %s", at(list.field), g, r))
				}
			}
		}
		setup := strings.TrimSpace(c.Setup)
		if r := pathProblem(setup, "the sub-project", "the sub-project's folder", "a module", exampleComponentSetup, exampleComponentSetup); r != "" {
			errs = append(errs, fmt.Sprintf("%s %q %s", at("setup"), c.Setup, r))
		} else if setup != "" && fsys != nil && !isFile(fsys, path.Join(path.Clean(strings.TrimSpace(p.Path)), setup)) {
			errs = append(errs, fmt.Sprintf("%s %q names no file under %s; write the path of a module of the sub-project, relative to its folder, such as %s, or remove it", at("setup"), c.Setup, path.Clean(strings.TrimSpace(p.Path)), exampleComponentSetup))
		}
	}
	return errs
}

// projectField is a field of a sub-project's key as ComponentErrors and
// InstallErrors name it, such as projects[1].components.include (sub-project
// "web"), or the key itself when field is empty.
func projectField(i int, p Project, key, field string) string {
	s := fmt.Sprintf("projects[%d].%s", i, key)
	if field != "" {
		s += "." + field
	}
	if name := strings.TrimSpace(p.Name); name != "" {
		s += fmt.Sprintf(" (sub-project %q)", name)
	}
	return s
}

// projectGlobProblem says what is wrong with a glob of a sub-project's key,
// such as components.include or install.paths, or "" when it is valid: a path
// relative to the sub-project's folder and inside it, with the syntax of
// claims.shared's. example is a glob the reason gives as one to write.
func projectGlobProblem(g, example string) string {
	switch {
	case strings.TrimSpace(g) == "":
		return "is empty; write a glob relative to the sub-project's folder, such as " + example + ", or remove it"
	case isAbsolute(g):
		return "is an absolute path; write it relative to the sub-project's folder, such as " + example
	}
	for _, seg := range strings.Split(g, "/") {
		switch seg {
		case "..":
			return "has a .. segment, which leaves the sub-project; write a glob inside its folder, relative to it, such as " + example
		case ".":
			return "has a . segment; write the glob without it, such as " + example + " for ./" + example
		}
	}
	return patternProblem(g)
}

// Dashboard is how the flaiover image runs for this project.
type Dashboard struct {
	Image string `yaml:"image"`
	Tag   string `yaml:"tag"`
	Port  int    `yaml:"port"`
	Bind  string `yaml:"bind"` // host address to publish on; default 0.0.0.0
	// Autocommit commits documents saved from the dashboard; nil means true (ADR-0023).
	Autocommit *bool `yaml:"autocommit,omitempty"`
	// AllowUnsigned, when set, is this project's dashboard.allow_unsigned over
	// the host configuration's (ADR-0070); nil leaves the host's.
	AllowUnsigned *bool `yaml:"allow_unsigned,omitempty"`
	// NotifyURL, when set, is where the dashboard's server posts new inbox
	// entries (S-0042). flai itself never calls it.
	NotifyURL string `yaml:"notify_url,omitempty"`
}

// Autocommit reports whether documents saved from the dashboard are committed.
func (m Manifest) Autocommit() bool {
	return m.Dashboard.Autocommit == nil || *m.Dashboard.Autocommit
}

// AllowUnsignedOr reports whether unsigned dashboards and flais are allowed
// for this project, given the host configuration's dashboard.allow_unsigned.
func (d Dashboard) AllowUnsignedOr(host bool) bool {
	if d.AllowUnsigned != nil {
		return *d.AllowUnsigned
	}
	return host
}

// Load parses the manifest at path.
func Load(path string) (Manifest, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Manifest{}, err
	}
	if p := batchMaxProblem(data); p != "" {
		return Manifest{}, fmt.Errorf("%s: %s", path, p)
	}
	var m Manifest
	if err := yaml.Unmarshal(data, &m); err != nil {
		return Manifest{}, fmt.Errorf("parse %s: %w", path, err)
	}
	if m.Version == 0 {
		return Manifest{}, fmt.Errorf("%s: version is required", path)
	}
	if m.Name == "" {
		return Manifest{}, fmt.Errorf("%s: name is required", path)
	}
	for _, k := range []string{"design", "docs", "wip"} {
		if m.Layout[k] == "" {
			return Manifest{}, fmt.Errorf("%s: layout.%s is required", path, k)
		}
	}
	if r := repoPathProblem(strings.TrimSpace(m.Design.Tokens), "a file", DefaultDesignTokens, DefaultDesignTokens); r != "" {
		return Manifest{}, fmt.Errorf("%s: design.tokens %q %s", path, m.Design.Tokens, r)
	}
	if r := repoPathProblem(m.DesignStylesheet(), "a file", exampleDesignStylesheet, exampleDesignStylesheet); r != "" {
		return Manifest{}, fmt.Errorf("%s: design.stylesheet %q %s", path, m.Design.Stylesheet, r)
	}
	if errs := m.Acceptance.Errors(); len(errs) > 0 {
		return Manifest{}, fmt.Errorf("%s: %s", path, strings.Join(errs, "; "))
	}
	if errs := m.MergeQueue.Errors(); len(errs) > 0 {
		return Manifest{}, fmt.Errorf("%s: %s", path, strings.Join(errs, "; "))
	}
	if r := repoPathProblem(strings.TrimSpace(m.Flai.Command), "a program", "scripts/flai.sh", "scripts/flai.sh"); r != "" {
		return Manifest{}, fmt.Errorf("%s: flai.command %q %s", path, m.Flai.Command, r)
	}
	if min := m.Flai.Minimum; min != "" {
		if _, ok := buildinfo.Semver(min); !ok {
			return Manifest{}, fmt.Errorf("%s: flai.minimum %q is not a release version like 1.27.0", path, min)
		}
		if buildinfo.Below(buildinfo.Version, min) {
			return Manifest{}, &TooOldError{Path: path, Minimum: min, Running: buildinfo.Version}
		}
	}
	return m, nil
}

// SetMinimum writes flai.minimum into the manifest at path, keeping the rest
// of the file as written: the minimum line under a flai: block is replaced
// or added, and a manifest with no flai: block gets one at its end.
func SetMinimum(path, version string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	line := "  minimum: " + version
	lines := strings.Split(strings.TrimRight(string(data), "\n"), "\n")
	block := -1
	for i, l := range lines {
		if strings.TrimRight(l, " ") == "flai:" {
			block = i
			break
		}
	}
	if block < 0 {
		lines = append(lines, "flai:", line)
	} else {
		at := -1
		for i := block + 1; i < len(lines) && (strings.HasPrefix(lines[i], " ") || lines[i] == ""); i++ {
			if strings.HasPrefix(strings.TrimSpace(lines[i]), "minimum:") {
				at = i
				break
			}
		}
		if at >= 0 {
			lines[at] = line
		} else {
			lines = append(lines[:block+1], append([]string{line}, lines[block+1:]...)...)
		}
	}
	return os.WriteFile(path, []byte(strings.Join(lines, "\n")+"\n"), 0o644)
}

// Find walks up from start looking for the manifest and returns its path.
func Find(start string) (string, error) {
	dir, err := filepath.Abs(start)
	if err != nil {
		return "", err
	}
	for {
		p := filepath.Join(dir, File)
		if _, err := os.Stat(p); err == nil {
			return p, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", ErrNotFound
		}
		dir = parent
	}
}

// Dir returns the folder for a layout key, resolved against the project root.
func (m Manifest) Dir(root, key string) string {
	return filepath.Join(root, m.Layout[key])
}

// Layer is the Dockerfile agent.image.layer names, from the repository root, or "".
func (m Manifest) Layer() string {
	if m.Agent == nil || m.Agent.Image == nil {
		return ""
	}
	return m.Agent.Image.Layer
}

// LayerImage is the project's layer image for the inner harness, flai-agent-<key>-<inner>:local, or "".
func (m Manifest) LayerImage(inner string) string {
	if m.Layer() == "" {
		return ""
	}
	return "flai-agent-" + strings.ToLower(m.Key) + "-" + inner + ":local"
}
