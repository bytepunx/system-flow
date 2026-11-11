// Package config reads and writes the flai user configuration file.
//
// The file lives at ~/.flai/config.json by default (ADR-0008). It is plain
// JSON with a fixed schema; unknown keys are rejected so typos surface early.
package config

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"maps"
	"net"
	"net/url"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"unicode"

	"github.com/bytepunx/system-flow/flai/internal/manifest"
)

// EnvVar overrides the config file path when set.
const EnvVar = "FLAI_CONFIG"

// CacheEnvVar overrides the default cache_dir written on first run, so
// scripts and tests can keep template clones out of the operator's home.
const CacheEnvVar = "FLAI_CACHE_DIR"

// Config is the schema of ~/.flai/config.json.
type Config struct {
	Template  Template  `json:"template"`
	Dashboard Dashboard `json:"dashboard"`
	CacheDir  string    `json:"cache_dir"`
	Author    string    `json:"author"`
	Worktrees Worktrees `json:"worktrees"`
	// HostActions are the actions flai serve may perform on this host when a
	// dashboard asks (ADR-0029): the action's name, and for which projects,
	// each the absolute path of a main checkout, or "*" for every project. It
	// is empty by default, which means none, and it is not among Keys: flai
	// serve enable and disable manage it, on the host, and nothing a
	// dashboard can ask for reads or writes this file.
	HostActions map[string][]string `json:"host_actions,omitempty"`
	// Agent is what flai serve starts when a story becomes ready and nobody
	// is attending the project, once the agent action is enabled (S-0079).
	// There is no default command: an operator writes one, as an argument
	// list that is run as it stands, never through a shell. Like HostActions
	// it is not among Keys; flai serve agent manages it.
	Agent AgentStart `json:"agent,omitzero"`
	// Forges are the forges flai calls for a project that accepts by pull
	// request (S-0392): for each, by name, the name of the environment
	// variable that holds its token, as Agent.Providers name a key's (ADR-0129).
	// It never holds a token; flai serve forge manages it, so it is not among
	// Keys.
	Forges map[string]ForgeHost `json:"forges,omitempty"`
	// Checks are the commands a story in review is checked with, once the
	// checks action is enabled (S-0082). When Commands is empty the
	// manifest's own checks: are used instead; this is not among Keys, and
	// flai serve checks manages it.
	Checks ChecksConfig `json:"checks,omitzero"`
	// Components is the host's say about the component dev servers flai
	// serve runs once the components action is enabled (S-0438); not among
	// Keys, flai serve components manages it.
	Components ComponentsConfig `json:"components,omitzero"`
	// ImportRoots are the folders flai serve looks in for repositories that
	// are not system-flow projects yet, which the board then offers to import
	// (S-0098). Naming a folder is the operator's consent to importing, and
	// testing, what is in it. Absolute paths; flai serve import manages it,
	// so it is not among Keys.
	ImportRoots []string `json:"import_roots,omitempty"`
	// Notify is the channel flai serve reminds the operator through of a
	// thread that has awaited them past the manifest's threads.remind_after
	// (ADR-0205): a webhook, a shell command, or both. Unset, a reminder is
	// only logged. flai serve notify manages it, on the host, so it is not
	// among Keys.
	Notify Notify `json:"notify,omitzero"`
	// MCP is the host's say about the HTTP MCP server flai host starts for
	// each project it serves (S-0506). flai serve mcp manages it, on the
	// host, so it is not among Keys.
	MCP MCPHost `json:"mcp,omitzero"`
}

// DefaultMCPBind is the address flai host starts each project's MCP server
// on when mcp.bind is unset: this machine only.
const DefaultMCPBind = "127.0.0.1"

// MCPHost is the host's say about each project's HTTP MCP server.
type MCPHost struct {
	// Bind is the address, an IP address or a host name without a port, that
	// flai host starts each project's flai mcp http on, keeping the port it
	// picks for the project; DefaultMCPBind when empty. An address a
	// container reaches lets an agent in one call the server; the bearer
	// token and the Origin refusal guard it wherever it listens.
	Bind string `json:"bind,omitempty"`
}

// MCPBind is the address flai host starts each project's MCP server on:
// mcp.bind, or DefaultMCPBind when it is unset.
func (c Config) MCPBind() string {
	if c.MCP.Bind == "" {
		return DefaultMCPBind
	}
	return c.MCP.Bind
}

// WithMCPBind returns the configuration with mcp.bind set to addr, refusing
// one that is not an IP address or a host name, or that names a port.
func (c Config) WithMCPBind(addr string) (Config, error) {
	addr = strings.TrimSpace(addr)
	if addr == "" {
		return c, errors.New("mcp.bind needs an address, such as 127.0.0.1 or 0.0.0.0; flai serve mcp clear goes back to " + DefaultMCPBind)
	}
	if err := checkMCPBind(addr); err != nil {
		return c, err
	}
	c.MCP.Bind = addr
	return c, nil
}

// WithoutMCPBind returns the configuration with mcp.bind unset, so that flai
// host starts each project's MCP server on DefaultMCPBind.
func (c Config) WithoutMCPBind() Config {
	c.MCP.Bind = ""
	return c
}

// checkMCPBind refuses an mcp.bind that is neither an IP address nor a host
// name, or that names a port; empty is unset and passes.
func checkMCPBind(addr string) error {
	if addr == "" {
		return nil
	}
	if why := badBindAddress(addr); why != "" {
		return fmt.Errorf("mcp.bind %q %s; give an IP address or a host name without a port, such as 127.0.0.1, 0.0.0.0, or 172.17.0.1, with flai serve mcp bind <address>", addr, why)
	}
	return nil
}

// badBindAddress is why addr is not an IP address or a host name with no
// port, or "" when it is one.
func badBindAddress(addr string) string {
	if net.ParseIP(addr) != nil {
		return ""
	}
	if strings.HasPrefix(addr, "[") {
		return "is bracketed: write an IPv6 address without brackets or a port"
	}
	if strings.Contains(addr, ":") {
		return "names a port, or is no IP address: flai picks each project's port itself"
	}
	if len(addr) > 253 {
		return "is longer than a host name may be"
	}
	for _, label := range strings.Split(addr, ".") {
		if label == "" || len(label) > 63 || strings.HasPrefix(label, "-") || strings.HasSuffix(label, "-") {
			return "is not a host name"
		}
		for _, r := range label {
			if (r < 'a' || r > 'z') && (r < 'A' || r > 'Z') && (r < '0' || r > '9') && r != '-' {
				return "is not a host name"
			}
		}
	}
	return ""
}

// Notify is the host's reminder channel (ADR-0205). Either may be empty,
// and both empty means no channel.
type Notify struct {
	// URL is a webhook, http or https, that each reminder is POSTed to as
	// JSON. It may carry a secret in its path or query, so flai logs its host
	// alone.
	URL string `json:"url,omitempty"`
	// Command is run through the host's shell with each reminder's JSON on
	// its standard input.
	Command string `json:"command,omitempty"`
}

// WithNotify returns the configuration with its reminder channel set to n,
// refusing a URL that is not an absolute http or https URL and a command
// that is only white space. A channel that sets nothing removes the channel.
func (c Config) WithNotify(n Notify) (Config, error) {
	if err := checkNotify(n); err != nil {
		return c, err
	}
	c.Notify = n
	return c, nil
}

// WithoutNotify returns the configuration with no reminder channel, so that
// flai serve logs a reminder and sends it nowhere.
func (c Config) WithoutNotify() Config {
	c.Notify = Notify{}
	return c
}

// checkNotify refuses a notify.url that is not an absolute http or https URL
// with a host and a notify.command that is only white space, with every
// problem.
func checkNotify(n Notify) error {
	var errs []string
	if n.URL != "" {
		u, err := url.Parse(n.URL)
		switch {
		case err != nil:
			errs = append(errs, "notify.url is not a URL; write an http or https webhook, such as https://ntfy.sh/<topic>, with flai serve notify set --url <url>")
		case u.Scheme == "":
			errs = append(errs, "notify.url names no scheme; write an absolute http or https URL, such as https://ntfy.sh/<topic>, with flai serve notify set --url <url>")
		case u.Scheme != "http" && u.Scheme != "https":
			errs = append(errs, fmt.Sprintf("notify.url's scheme is %s, and flai posts reminders over http or https alone; write a URL such as https://ntfy.sh/<topic> with flai serve notify set --url <url>", u.Scheme))
		case u.Host == "":
			errs = append(errs, "notify.url names no host; write an absolute http or https URL, such as https://ntfy.sh/<topic>, with flai serve notify set --url <url>")
		}
	}
	if n.Command != "" && strings.TrimSpace(n.Command) == "" {
		errs = append(errs, "notify.command is empty; write a shell command that reads the reminder's JSON on its standard input with flai serve notify set --command <command>, or remove it with flai serve notify clear")
	}
	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "; "))
	}
	return nil
}

// AgentStart is the command flai serve starts an agent with.
type AgentStart struct {
	// Command is the program and its arguments. In an argument, {story} is
	// replaced by the story's ID and {root} by the project's directory;
	// nothing else is interpreted.
	Command []string `json:"command,omitempty"`
	// Name is the FLAI_AGENT the session works under; "agent" when empty.
	Name string `json:"name,omitempty"`
	// AttendedMinutes is retired (S-0116, ADR-0043): nobody attending holds
	// a ready story back any more. It is still read so that a configuration
	// that sets it loads, and nothing uses it.
	AttendedMinutes int `json:"attended_minutes,omitempty"`
	// Harnesses are the operator's say about each harness a story may name
	// (S-0104): the program that is run and the arguments that say what the
	// agent may do. A harness not here runs with its adapter's defaults.
	Harnesses map[string]HarnessHost `json:"harnesses,omitempty"`
	// Providers are where the model calls of a story's agent that names a
	// provider go (ADR-0129): each entry's API, base URL, the name of the
	// environment variable that holds its key, and the harness's model
	// aliases on it, every entry whole. It never holds a key; flai serve
	// agent provider manages it, and a project's manifest may override an
	// entry's api, base_url, and models.
	Providers map[string]manifest.Provider `json:"providers,omitempty"`
	// AutoRestarts is how many times flai serve restarts a story's agent on
	// its own after it ends with its story in progress (ADR-0108). Nil means
	// DefaultAutoRestarts, 0 turns the restart off, and a negative count is
	// refused; read it through AutoRestartLimit.
	AutoRestarts *int `json:"auto_restarts,omitempty"`
	// Projects are each project's own command, name, and harnesses, by the
	// absolute path of its main checkout, as host_actions are keyed
	// (ADR-0261): read over the host-wide Command, Name, and Harnesses above,
	// which stay the fallback for a project that sets none. Read a project's
	// settings through ForProject.
	Projects map[string]ProjectAgent `json:"projects,omitempty"`
}

// ProjectAgent is one project's own agent settings on this host, over the
// host-wide ones (ADR-0261).
type ProjectAgent struct {
	// Command is the project's command, the program and its arguments; nil
	// leaves the host-wide command.
	Command []string `json:"command,omitempty"`
	// Name is the FLAI_AGENT the project's sessions work under; empty leaves
	// the host-wide name.
	Name string `json:"name,omitempty"`
	// Harnesses are the project's entries for the harnesses it sets, each
	// replacing the host-wide entry of that harness whole.
	Harnesses map[string]HarnessHost `json:"harnesses,omitempty"`
}

// IsZero reports whether the entry sets nothing.
func (p ProjectAgent) IsZero() bool {
	return len(p.Command) == 0 && p.Name == "" && len(p.Harnesses) == 0
}

// Project is the project at root's own entry, empty when it sets nothing.
func (a AgentStart) Project(root string) ProjectAgent { return a.Projects[root] }

// ForProject is the agent settings that apply to the project at root
// (ADR-0261): its own command and name where it sets them, else the
// host-wide ones; each harness entry it sets in place of the host-wide
// entry of that harness, whole; and the host-wide providers and restart
// limit. The result carries no Projects. The maps are copied.
func (a AgentStart) ForProject(root string) AgentStart {
	p := a.Projects[root]
	out := a
	out.Projects = nil
	if len(p.Command) > 0 {
		out.Command = slices.Clone(p.Command)
	}
	if p.Name != "" {
		out.Name = p.Name
	}
	if len(p.Harnesses) > 0 {
		next := maps.Clone(a.Harnesses)
		if next == nil {
			next = map[string]HarnessHost{}
		}
		maps.Copy(next, p.Harnesses)
		out.Harnesses = next
	}
	return out
}

// WithProject returns the settings with the project at root's entry set to
// p, refusing a root that is not an absolute path (the host-wide settings
// are AgentStart's own fields, not an entry), a command whose program is empty, and a harness entry checkHarnesses
// refuses; an entry that sets nothing is removed, and a harness entry that
// sets nothing is dropped from it. The maps are copied.
func (a AgentStart) WithProject(root string, p ProjectAgent) (AgentStart, error) {
	if len(p.Harnesses) > 0 {
		hs := map[string]HarnessHost{}
		for name, h := range p.Harnesses {
			if !h.IsZero() {
				hs[name] = h
			}
		}
		if len(hs) == 0 {
			hs = nil
		}
		p.Harnesses = hs
	}
	if err := checkProjects(map[string]ProjectAgent{root: p}); err != nil {
		return a, err
	}
	next := maps.Clone(a.Projects)
	if next == nil {
		next = map[string]ProjectAgent{}
	}
	if p.IsZero() {
		delete(next, root)
	} else {
		next[root] = p
	}
	if len(next) == 0 {
		next = nil
	}
	a.Projects = next
	return a, nil
}

// checkProjects refuses a project entry keyed by anything but an absolute
// path, one whose command has an empty program, and one whose harness entries
// checkHarnesses refuses, each naming the project.
func checkProjects(projects map[string]ProjectAgent) error {
	var errs []string
	for _, root := range slices.Sorted(maps.Keys(projects)) {
		p := projects[root]
		key := fmt.Sprintf("agent.projects[%q]", root)
		if !filepath.IsAbs(root) {
			errs = append(errs, fmt.Sprintf("%s is keyed by %q, which is not the absolute path of a project's main checkout: set a project's own agent settings with flai serve agent set, clear, or harness in the project, or --project <root>", key, root))
			continue
		}
		if len(p.Command) > 0 && strings.TrimSpace(p.Command[0]) == "" {
			errs = append(errs, fmt.Sprintf("%s.command has an empty program: set it with flai serve agent set --project %s -- <program> [args...], or clear it with flai serve agent clear --project %s", key, root, root))
		}
		if err := checkHarnesses(p.Harnesses); err != nil {
			errs = append(errs, key+": "+err.Error())
		}
	}
	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "; "))
	}
	return nil
}

// DefaultAutoRestarts is how many times flai serve restarts a story's agent
// on its own while agent.auto_restarts is unset (ADR-0108).
const DefaultAutoRestarts = 2

// AutoRestartLimit is how many times flai serve restarts a story's agent on
// its own: AutoRestarts, or DefaultAutoRestarts while it is unset.
func (a AgentStart) AutoRestartLimit() int {
	if a.AutoRestarts == nil {
		return DefaultAutoRestarts
	}
	return *a.AutoRestarts
}

// WithAutoRestarts returns the settings with AutoRestarts set to n, refusing
// a negative n.
func (a AgentStart) WithAutoRestarts(n int) (AgentStart, error) {
	if err := checkAutoRestarts(n); err != nil {
		return a, err
	}
	a.AutoRestarts = &n
	return a, nil
}

// WithProvider returns the settings with the provider name set to p, refusing
// a name that is not one and an entry manifest.HostProviderProblems refuses,
// such as one without its api, base_url, or key_env. The map is copied.
func (a AgentStart) WithProvider(name string, p manifest.Provider) (AgentStart, error) {
	if err := checkProviders(map[string]manifest.Provider{name: p}); err != nil {
		return a, err
	}
	next := maps.Clone(a.Providers)
	if next == nil {
		next = map[string]manifest.Provider{}
	}
	next[name] = p
	a.Providers = next
	return a, nil
}

// WithoutProvider returns the settings without the provider name. The map is
// copied.
func (a AgentStart) WithoutProvider(name string) AgentStart {
	next := maps.Clone(a.Providers)
	delete(next, name)
	if len(next) == 0 {
		next = nil
	}
	a.Providers = next
	return a
}

// checkProviders refuses providers manifest.HostProviderProblems finds wrong,
// with every problem, each naming the setting that fixes it.
func checkProviders(providers map[string]manifest.Provider) error {
	if errs := manifest.HostProviderProblems(providers); len(errs) > 0 {
		return errors.New(strings.Join(errs, "; "))
	}
	return nil
}

// checkAutoRestarts refuses a negative count of automatic restarts.
func checkAutoRestarts(n int) error {
	if n < 0 {
		return fmt.Errorf("agent.auto_restarts is %d: give how many times flai serve restarts a story's agent on its own, 0 or more, with flai serve agent set --auto-restarts <n> (0 turns it off; unset is %d)", n, DefaultAutoRestarts)
	}
	return nil
}

// HarnessHost is the program a harness is, the operator's arguments to it,
// and what the operator lifts on it.
type HarnessHost struct {
	Program string `json:"program,omitempty"`
	// Args replace the adapter's defaults when set, even to an empty list;
	// nil leaves them.
	Args *[]string `json:"args,omitempty"`
	// Guard GuardNone is the operator taking the risk that a story's
	// sub-agents run without flai guard on a harness that cannot run it
	// (ADR-0131); empty leaves such a story refused. No other value is taken.
	Guard string `json:"guard,omitempty"`
	// DenyProtected is the operator's say that their arguments for the
	// harness already deny the protected paths in the harness's own
	// vocabulary, which flai cannot read (ADR-0131).
	DenyProtected bool `json:"deny_protected,omitempty"`
	// Images are the container harness's image for each inner harness, by
	// the inner harness's name; one not here runs in DefaultImage's (S-0508).
	Images map[string]string `json:"images,omitempty"`
	// Network is the Docker network the container harness's agents join;
	// empty is Docker's default bridge, the host reached through Docker's
	// host gateway.
	Network string `json:"network,omitempty"`
	// Login is the absolute path of the folder holding the operator's Claude
	// login, which the container harness mounts read-only.
	Login string `json:"login,omitempty"`
	// Outside are the folders outside the project that a story's agent and
	// its sub-agents may read and write, each absolute and cleaned, as
	// OutsideFolders checks them: passed as --add-dir by the claude-code
	// harness and mounted by the container harness (ADR-0289). Unlike
	// Images, Network, and Login, every harness takes them.
	Outside []string `json:"outside,omitempty"`
}

// IsZero reports whether the entry sets nothing.
func (h HarnessHost) IsZero() bool {
	return h.Program == "" && h.Args == nil && h.Guard == "" && !h.DenyProtected && len(h.Images) == 0 && h.Network == "" && h.Login == "" && len(h.Outside) == 0
}

// Image is the container harness's image for the inner harness: the one
// Images sets, else DefaultImage's.
func (h HarnessHost) Image(inner string) string {
	if ref := h.Images[inner]; ref != "" {
		return ref
	}
	return DefaultImage(inner)
}

// ShownImages are the container harness's image for every inner harness it
// runs, the defaults filled in, as flai shows them.
func (h HarnessHost) ShownImages() map[string]string {
	out := map[string]string{}
	for _, inner := range containerInners {
		out[inner] = h.Image(inner)
	}
	return out
}

// GuardNone is the one value HarnessHost.Guard takes: the harness's story
// sub-agents run unguarded.
const GuardNone = "none"

// ContainerHarness is the harness that runs an inner harness in a Docker
// container, harness.Container; config cannot import harness, which imports
// config through template.
const ContainerHarness = "container"

// containerInners are the harnesses flai runs inside the container, the
// container adapter's inner harnesses among harness.Names, kept here for the
// same reason as ContainerHarness.
var containerInners = []string{"claude-code"}

// ContainerInners are the harnesses flai runs inside the container.
func ContainerInners() []string { return slices.Clone(containerInners) }

// DefaultImage is the image the container harness runs the inner harness in
// while Images names none: the one make agent-images builds.
func DefaultImage(inner string) string { return "flai-agent-" + inner + ":local" }

// WithHarness returns the settings with the harness name's entry set to h,
// refusing an entry checkHarnesses refuses; an entry that sets nothing is
// removed. The map is copied.
func (a AgentStart) WithHarness(name string, h HarnessHost) (AgentStart, error) {
	if err := checkHarnesses(map[string]HarnessHost{name: h}); err != nil {
		return a, err
	}
	next := maps.Clone(a.Harnesses)
	if next == nil {
		next = map[string]HarnessHost{}
	}
	if h.IsZero() {
		delete(next, name)
	} else {
		next[name] = h
	}
	if len(next) == 0 {
		next = nil
	}
	a.Harnesses = next
	return a, nil
}

// checkHarnesses refuses a harness entry whose guard checkGuard refuses,
// whose container settings checkContainer refuses, or whose outside folders
// checkOutside refuses.
func checkHarnesses(harnesses map[string]HarnessHost) error {
	var errs []string
	for _, name := range slices.Sorted(maps.Keys(harnesses)) {
		if err := checkGuard(name, harnesses[name].Guard); err != nil {
			errs = append(errs, err.Error())
		}
		errs = append(errs, checkContainer(name, harnesses[name])...)
		errs = append(errs, checkOutside(name, harnesses[name].Outside)...)
	}
	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "; "))
	}
	return nil
}

// checkContainer refuses the container settings of the harness name's entry
// h: any on a harness other than ContainerHarness, an image for a harness
// flai does not run inside the container, an image reference that is empty,
// holds a space, or starts with "-", a network that is no Docker network
// name, and a login path that is not absolute.
func checkContainer(name string, h HarnessHost) []string {
	if name != ContainerHarness {
		if len(h.Images) > 0 || h.Network != "" || h.Login != "" {
			return []string{fmt.Sprintf("agent.harnesses.%s sets images, network, or login, which only the %s harness takes: remove them from agent.harnesses.%s in the configuration file; flai serve agent harness %s --image, --network, and --login set the container's", name, ContainerHarness, name, ContainerHarness)}
		}
		return nil
	}
	var errs []string
	for _, inner := range slices.Sorted(maps.Keys(h.Images)) {
		ref := h.Images[inner]
		key := fmt.Sprintf("agent.harnesses.%s.images.%s", name, inner)
		switch {
		case !slices.Contains(containerInners, inner):
			errs = append(errs, fmt.Sprintf("%s: flai runs no harness %q inside the container, only %s; set an image with flai serve agent harness %s --image <inner>=<image>, or --reset", key, inner, strings.Join(containerInners, ", "), name))
		case ref == "":
			errs = append(errs, fmt.Sprintf("%s is empty: give the image reference, as --image %s=%s, or leave it unset for that default", key, inner, DefaultImage(inner)))
		case strings.HasPrefix(ref, "-") || strings.ContainsFunc(ref, unicode.IsSpace):
			errs = append(errs, fmt.Sprintf("%s is %q, which is no image reference: give one such as %s with flai serve agent harness %s --image %s=<image>", key, ref, DefaultImage(inner), name, inner))
		}
	}
	if h.Network != "" && !dockerName.MatchString(h.Network) {
		errs = append(errs, fmt.Sprintf("agent.harnesses.%s.network is %q, which is no Docker network name: name one docker network ls lists with flai serve agent harness %s --network <name>, or clear it with --network \"\" for Docker's default bridge", name, h.Network, name))
	}
	if h.Login != "" && !filepath.IsAbs(h.Login) {
		errs = append(errs, fmt.Sprintf("agent.harnesses.%s.login is %q, which is not an absolute path: set the folder holding your Claude login with flai serve agent harness %s --login ~/.claude", name, h.Login, name))
	}
	return errs
}

// checkOutside refuses an outside folder of the harness name's entry that is
// not an absolute path, as a hand-edited file may hold; OutsideFolders, which
// flai serve agent harness runs, checks the rest against the disk (ADR-0289).
func checkOutside(name string, dirs []string) []string {
	var errs []string
	for _, dir := range dirs {
		if !filepath.IsAbs(dir) {
			errs = append(errs, fmt.Sprintf("agent.harnesses.%s.outside holds %q, which is not an absolute path: set the folders with flai serve agent harness %s --outside <dir>, given absolute or under ~/, or clear them with --outside \"\"", name, dir, name))
		}
	}
	return errs
}

// OutsideFolders is the folders dirs, given as flai serve agent harness
// name's --outside, each absolute and cleaned, ~/ expanded against home, in
// the order given and each once: the folders outside the project that a
// story's agent and its sub-agents may read and write (ADR-0289). It refuses
// a folder that is not absolute or under ~/, one that does not exist or is a
// file, home itself, one that is, is under, or holds home's .flai, and one
// that holds a project at roots or that a project at roots holds, every
// refusal naming the folder and the flag that fixes it. A symbolic link is
// checked as itself and as what it points at.
func OutsideFolders(name string, dirs []string, home string, roots []string) ([]string, error) {
	var out, errs []string
	fix := fmt.Sprintf("name another folder with flai serve agent harness %s --outside <dir>, or clear them with --outside \"\"", name)
	flaiDir := filepath.Join(home, ".flai")
	for _, given := range dirs {
		dir := given
		if dir == "~" || strings.HasPrefix(dir, "~/") {
			dir = filepath.Join(home, strings.TrimPrefix(dir, "~"))
		}
		if !filepath.IsAbs(dir) {
			errs = append(errs, fmt.Sprintf("--outside %q is not an absolute path: give the folder absolute or under ~/, as flai serve agent harness %s --outside ~/src/other", given, name))
			continue
		}
		dir = filepath.Clean(dir)
		if info, err := os.Stat(dir); err != nil || !info.IsDir() {
			errs = append(errs, fmt.Sprintf("--outside %q: %s is not an existing folder; %s", given, dir, fix))
			continue
		}
		forms := pathForms(dir)
		var problem string
		switch {
		case overlaps(forms, pathForms(home), false):
			problem = "is your home folder, which would open every file you own to the agent"
		case overlaps(forms, pathForms(flaiDir), true):
			problem = fmt.Sprintf("is, is under, or holds %s, flai's own configuration and state, which an agent that wrote could widen its own permissions (ADR-0220)", flaiDir)
		default:
			for _, root := range roots {
				rootForms := pathForms(root)
				if anyHolds(forms, rootForms) {
					problem = fmt.Sprintf("holds the project at %s, whose main checkout only flai writes", root)
					break
				}
				if anyHolds(rootForms, forms) {
					problem = fmt.Sprintf("is inside the project at %s, which a story's agent reaches through its worktree already", root)
					break
				}
			}
		}
		if problem != "" {
			errs = append(errs, fmt.Sprintf("--outside %q: %s %s; %s", given, dir, problem, fix))
			continue
		}
		if !slices.Contains(out, dir) {
			out = append(out, dir)
		}
	}
	if len(errs) > 0 {
		return nil, errors.New(strings.Join(errs, "; "))
	}
	return out, nil
}

// pathForms are the absolute path p as given and, when it differs, the path
// its symbolic links resolve to.
func pathForms(p string) []string {
	forms := []string{filepath.Clean(p)}
	if resolved, err := filepath.EvalSymlinks(p); err == nil && resolved != forms[0] {
		forms = append(forms, resolved)
	}
	return forms
}

// overlaps reports whether a form of a is a form of b, or, with nested,
// whether one holds the other.
func overlaps(a, b []string, nested bool) bool {
	if nested {
		return anyHolds(a, b) || anyHolds(b, a)
	}
	for _, x := range a {
		if slices.Contains(b, x) {
			return true
		}
	}
	return false
}

// anyHolds reports whether a form of parent is, or holds, a form of child.
func anyHolds(parents, children []string) bool {
	for _, p := range parents {
		for _, c := range children {
			if rel, err := filepath.Rel(p, c); err == nil && rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
				return true
			}
		}
	}
	return false
}

// dockerName is what Docker takes as a network's name, and the container:
// form that joins another container's network.
var dockerName = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9_.:-]*$`)

// checkGuard refuses a harness's guard that is neither GuardNone nor empty.
func checkGuard(name, guard string) error {
	if guard == "" || guard == GuardNone {
		return nil
	}
	return fmt.Errorf("agent.harnesses.%s.guard is %q: the only guard flai takes is %q, which lets a story's sub-agents run unguarded on a harness that cannot run flai guard (ADR-0131); set it with flai serve agent harness %s --guard none, or clear it with --guard \"\"", name, guard, GuardNone, name)
}

// ForgeHost is the host's say about one forge.
type ForgeHost struct {
	// TokenEnv is the name of the environment variable that holds the
	// forge's token; the forge's default, such as GITHUB_TOKEN, when empty.
	TokenEnv string `json:"token_env,omitempty"`
}

// DefaultGitHubTokenEnv is the variable GitHub's token is read from while
// the host names none.
const DefaultGitHubTokenEnv = "GITHUB_TOKEN"

// defaultTokenEnvs are the forges flai calls, each with the variable its
// token is read from while the host names none.
var defaultTokenEnvs = map[string]string{manifest.ForgeGitHub: DefaultGitHubTokenEnv}

// KnownForges are the forges flai calls, in the order they are listed.
func KnownForges() []string {
	return slices.Sorted(maps.Keys(defaultTokenEnvs))
}

// DefaultForgeTokenEnv is the variable forge's token is read from while the
// host names none, or "" for a forge flai does not call.
func DefaultForgeTokenEnv(forge string) string {
	return defaultTokenEnvs[forge]
}

// ForgeTokenEnv is the name of the variable forge's token is read from: the
// host's forges entry, or the forge's default; "" for a forge flai does not
// call.
func (c Config) ForgeTokenEnv(forge string) string {
	if _, ok := defaultTokenEnvs[forge]; !ok {
		return ""
	}
	if env := c.Forges[forge].TokenEnv; env != "" {
		return env
	}
	return defaultTokenEnvs[forge]
}

// ForgeTokenFunc returns a function that reads forge's token from the
// variable ForgeTokenEnv names, through getenv, each time it is called and
// never before, so that the token is read only when flai calls the forge.
// The function refuses a forge flai does not call and a variable that is
// unset or empty, naming the variable and the setting that renames it,
// never a value.
func (c Config) ForgeTokenFunc(forge string, getenv func(string) string) func() (string, error) {
	env := c.ForgeTokenEnv(forge)
	return func() (string, error) {
		if env == "" {
			return "", fmt.Errorf("flai calls no forge %q; it calls %s", forge, strings.Join(KnownForges(), ", "))
		}
		token := strings.TrimSpace(getenv(env))
		if token == "" {
			return "", fmt.Errorf("the %s token is read from the variable %s, which is not set in flai's environment: export it there, or name the variable that holds it with flai serve forge %s --token-env <NAME>", forge, env, forge)
		}
		return token, nil
	}
}

// WithForge returns the configuration with forge's entry set to f, refusing
// a forge flai does not call and a token_env that is not a variable's name;
// an entry that sets nothing is removed. The map is copied.
func (c Config) WithForge(forge string, f ForgeHost) (Config, error) {
	if err := checkForges(map[string]ForgeHost{forge: f}); err != nil {
		return c, err
	}
	next := maps.Clone(c.Forges)
	if next == nil {
		next = map[string]ForgeHost{}
	}
	if f == (ForgeHost{}) {
		delete(next, forge)
	} else {
		next[forge] = f
	}
	if len(next) == 0 {
		next = nil
	}
	c.Forges = next
	return c, nil
}

// WithoutForge returns the configuration without forge's entry, so that its
// token is read from the forge's default variable again. The map is copied.
func (c Config) WithoutForge(forge string) Config {
	next := maps.Clone(c.Forges)
	delete(next, forge)
	if len(next) == 0 {
		next = nil
	}
	c.Forges = next
	return c
}

// checkForges refuses a forges entry for a forge flai does not call and a
// token_env that is not a variable's name, with every problem.
func checkForges(forges map[string]ForgeHost) error {
	var errs []string
	for _, name := range slices.Sorted(maps.Keys(forges)) {
		if _, ok := defaultTokenEnvs[name]; !ok {
			errs = append(errs, fmt.Sprintf("forges.%s: flai calls no forge %q; it calls %s", name, name, strings.Join(KnownForges(), ", ")))
			continue
		}
		// The value is not quoted: one that is not a variable's name may be
		// the token itself, pasted in its place.
		if env := forges[name].TokenEnv; env != "" && !manifest.IsEnvName(env) {
			errs = append(errs, fmt.Sprintf("forges.%s.token_env is not the name of an environment variable; write one such as %s (letters, digits, _, not beginning with a digit), never the token, with flai serve forge %s --token-env <NAME>", name, defaultTokenEnvs[name], name))
		}
	}
	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "; "))
	}
	return nil
}

// ChecksConfig is the operator's say about running checks for a story in
// review, on this host.
type ChecksConfig struct {
	// Commands are the named checks, run in order, in the story's worktree.
	// Empty means: use the manifest's checks: instead.
	Commands []NamedCommand `json:"commands,omitempty"`
	// TimeoutMinutes bounds one run of every command together; 15 when zero.
	TimeoutMinutes int `json:"timeout_minutes,omitempty"`
}

// ComponentsConfig is the operator's say about the component dev servers
// flai serve runs, on this host (ADR-0163).
type ComponentsConfig struct {
	// Idle is how long a dev server may go without a components.status
	// before flai serve stops it, as a Go duration such as 30m; 30 minutes
	// when empty.
	Idle string `json:"idle,omitempty"`
	// MaxServers is how many component dev servers flai serve runs at once
	// per project, main's and every story's together; 3 when zero.
	MaxServers int `json:"max_servers,omitempty"`
}

// NamedCommand is one command by name: an argument list, run as it stands,
// never through a shell. In an argument, {story} is replaced by the story's
// ID and {root} by the directory it runs in; nothing else is interpreted.
type NamedCommand struct {
	Name    string   `json:"name"`
	Command []string `json:"command"`
}

// AllProjects stands for every project in HostActions.
const AllProjects = "*"

// ActionEnabled reports whether a host action is enabled for a project.
func (c Config) ActionEnabled(action, root string) bool {
	for _, r := range c.HostActions[action] {
		if r == AllProjects || r == root {
			return true
		}
	}
	return false
}

// WithAction returns the configuration with an action enabled or disabled
// for a project (or AllProjects). Disabling for AllProjects disables it
// everywhere: no project keeps it.
func (c Config) WithAction(action, root string, on bool) Config {
	next := map[string][]string{}
	for a, roots := range c.HostActions {
		next[a] = append([]string{}, roots...)
	}
	var kept []string
	for _, r := range next[action] {
		if r != root && (on || root != AllProjects) {
			kept = append(kept, r)
		}
	}
	if on {
		kept = append(kept, root)
	}
	if len(kept) == 0 {
		delete(next, action)
	} else {
		next[action] = kept
	}
	if len(next) == 0 {
		next = nil
	}
	c.HostActions = next
	return c
}

// Worktrees is how flai stream open creates a story's worktree.
type Worktrees struct {
	// RelativePaths links worktrees with relative paths (git 2.48 or newer).
	// Off by default and never turned on by flai: creating one sets a
	// repository extension that older git refuses for the whole clone (ADR-0022).
	RelativePaths bool `json:"relative_paths"`
}

// Template is where flai fetches the project template from.
type Template struct {
	Repo string `json:"repo"`
	Ref  string `json:"ref"`
}

// Dashboard is how flai runs the flaiover image.
type Dashboard struct {
	Image string `json:"image"`
	Tag   string `json:"tag"`
	Port  int    `json:"port"`
	Bind  string `json:"bind"` // host address the port is published on; 0.0.0.0 for every interface
	// NoRestart stops flai host from restarting a dashboard that is gone or
	// not answering (S-0184). Off by default: the host restarts it.
	NoRestart bool `json:"no_restart"`
	// AllowUnsigned lets flai dial a dashboard that is no signed release and
	// lets the dashboard accept a flai that is none (ADR-0070). Off by default.
	AllowUnsigned bool `json:"allow_unsigned"`
	// PushKey is an SSH private key on this host that the container may push
	// acceptances with (ADR-0026). Empty, the default, gives it nothing. It is
	// a fact about one machine, so it lives here and never in the manifest.
	PushKey string `json:"push_key"`
	// PushKnownHosts names the file the remote's host keys are taken from
	// instead of this user's and the system's known_hosts.
	PushKnownHosts string `json:"push_known_hosts"`
	// AllowedHosts is a comma-separated list of host names, each with an
	// optional port, that the dashboard accepts in a browser's Origin beside
	// the address it is published on, such as a tunnel's or a domain's (S-0489).
	// flai starts the container with it, normalized (AllowedHostList), as
	// FLAIOVER_ALLOWED_HOSTS. Empty, the default, adds none. It is a fact about
	// one machine's shared container, so it lives here and never in the manifest.
	AllowedHosts string `json:"allowed_hosts"`
}

// AllowedHostList is raw, dashboard.allowed_hosts as written, as flai passes
// it on: each comma-separated entry trimmed and lower-cased, the empty ones
// dropped. It is nil when raw names none.
func AllowedHostList(raw string) []string {
	var out []string
	for _, h := range strings.Split(raw, ",") {
		if h = strings.ToLower(strings.TrimSpace(h)); h != "" {
			out = append(out, h)
		}
	}
	return out
}

// NormalizeAllowedHosts is AllowedHostList(raw) joined with commas, the value
// of FLAIOVER_ALLOWED_HOSTS: "" when raw names none.
func NormalizeAllowedHosts(raw string) string {
	return strings.Join(AllowedHostList(raw), ",")
}

// checkAllowedHosts refuses a dashboard.allowed_hosts entry that is not a
// host name with an optional port, naming each.
func checkAllowedHosts(raw string) error {
	var errs []string
	for _, h := range AllowedHostList(raw) {
		if why := badAllowedHost(h); why != "" {
			errs = append(errs, fmt.Sprintf("%q %s", h, why))
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("dashboard.allowed_hosts: %s; give host names, each with an optional :port, separated by commas, such as flow.example.com, board.example.net:8443", strings.Join(errs, "; "))
	}
	return nil
}

// badAllowedHost is why h, an entry AllowedHostList gave, is not a host name
// with an optional port, or "" when it is one. An IPv6 address is bracketed,
// as an Origin writes it.
func badAllowedHost(h string) string {
	if strings.ContainsAny(h, " \t/@?#") {
		return "is not a host name: write no scheme, path, or spaces"
	}
	var name, port string
	var hasPort bool
	if strings.HasPrefix(h, "[") {
		end := strings.Index(h, "]")
		if end < 0 {
			return "opens a bracket it does not close"
		}
		rest := h[end+1:]
		name = h[:end+1]
		if rest != "" && !strings.HasPrefix(rest, ":") {
			return "has something after its address that is no :port"
		}
		if name == "[]" {
			return "has no address in its brackets"
		}
		port, hasPort = strings.CutPrefix(rest, ":")
	} else {
		if strings.Count(h, ":") > 1 {
			return "is an IPv6 address without brackets"
		}
		name, port, hasPort = strings.Cut(h, ":")
	}
	if name == "" {
		return "has no host name"
	}
	if hasPort {
		if n, err := strconv.Atoi(port); err != nil || n < 1 || n > 65535 {
			return "has a port that is not a number from 1 to 65535"
		}
	}
	return ""
}

// Default returns the configuration written on first run.
func Default() Config {
	return Config{
		Template:  Template{Repo: "https://github.com/bytepunx/system-flow-template", Ref: "main"},
		Dashboard: Dashboard{Image: "ghcr.io/bytepunx/flaiover", Tag: "latest", Port: 4242, Bind: "0.0.0.0"},
		CacheDir:  orDefault(os.Getenv(CacheEnvVar), "~/.flai/cache"),
		Author:    currentUser(),
	}
}

// DefaultPath is ~/.flai/config.json, or a relative fallback if the home
// directory cannot be determined.
func DefaultPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return filepath.Join(".flai", "config.json")
	}
	return filepath.Join(home, ".flai", "config.json")
}

// ResolvePath picks the config path: the explicit override, then FLAI_CONFIG,
// then the default.
func ResolvePath(override string) string {
	if override != "" {
		return ExpandHome(override)
	}
	if env := os.Getenv(EnvVar); env != "" {
		return ExpandHome(env)
	}
	return DefaultPath()
}

// Load reads the config at path. If the file does not exist it is created
// with defaults and created reports true.
func Load(path string) (cfg Config, created bool, err error) {
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		cfg = Default()
		if err := Save(path, cfg); err != nil {
			return Config{}, false, err
		}
		return cfg, true, nil
	}
	if err != nil {
		return Config{}, false, fmt.Errorf("read config: %w", err)
	}
	cfg, err = Parse(data)
	if err != nil {
		return Config{}, false, fmt.Errorf("parse %s: %w", path, err)
	}
	return cfg, false, nil
}

// Parse decodes JSON into a Config, rejecting unknown keys, a negative
// agent.auto_restarts, an agent.harnesses entry checkHarnesses refuses, an
// agent.projects entry checkProjects refuses, an agent.providers entry that
// is not whole or sets a value flai refuses, a forges entry checkForges
// refuses, a dashboard.allowed_hosts entry that is not a host name with an
// optional port, a notify channel checkNotify refuses, and an mcp.bind that is no IP address or host name or names a
// port.
func Parse(data []byte) (Config, error) {
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	var cfg Config
	if err := dec.Decode(&cfg); err != nil {
		return Config{}, err
	}
	if n := cfg.Agent.AutoRestarts; n != nil {
		if err := checkAutoRestarts(*n); err != nil {
			return Config{}, err
		}
	}
	if err := checkHarnesses(cfg.Agent.Harnesses); err != nil {
		return Config{}, err
	}
	if err := checkProjects(cfg.Agent.Projects); err != nil {
		return Config{}, err
	}
	if err := checkProviders(cfg.Agent.Providers); err != nil {
		return Config{}, err
	}
	if err := checkForges(cfg.Forges); err != nil {
		return Config{}, err
	}
	if err := checkAllowedHosts(cfg.Dashboard.AllowedHosts); err != nil {
		return Config{}, err
	}
	if err := checkNotify(cfg.Notify); err != nil {
		return Config{}, err
	}
	if err := checkMCPBind(cfg.MCP.Bind); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

// Save writes cfg to path atomically, creating parent directories. It
// refuses an agent.harnesses entry checkHarnesses refuses and an
// agent.projects entry checkProjects refuses, which Load would refuse.
func Save(path string, cfg Config) error {
	if err := checkHarnesses(cfg.Agent.Harnesses); err != nil {
		return err
	}
	if err := checkProjects(cfg.Agent.Projects); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("create config dir: %w", err)
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	tmp, err := os.CreateTemp(filepath.Dir(path), ".config-*.json")
	if err != nil {
		return fmt.Errorf("write config: %w", err)
	}
	tmpName := tmp.Name()
	if _, err := tmp.Write(data); err != nil {
		_ = tmp.Close()
		_ = os.Remove(tmpName)
		return fmt.Errorf("write config: %w", err)
	}
	if err := tmp.Close(); err != nil {
		_ = os.Remove(tmpName)
		return fmt.Errorf("write config: %w", err)
	}
	if err := os.Chmod(tmpName, 0o644); err != nil {
		_ = os.Remove(tmpName)
		return fmt.Errorf("write config: %w", err)
	}
	if err := os.Rename(tmpName, path); err != nil {
		_ = os.Remove(tmpName)
		return fmt.Errorf("write config: %w", err)
	}
	return nil
}

// Keys lists every settable dotted key in schema order.
func Keys() []string {
	return []string{
		"template.repo", "template.ref",
		"dashboard.image", "dashboard.tag", "dashboard.port", "dashboard.bind", "dashboard.no_restart", "dashboard.allow_unsigned", "dashboard.push_key", "dashboard.push_known_hosts", "dashboard.allowed_hosts",
		"cache_dir", "author",
		"worktrees.relative_paths",
	}
}

// Get returns the value at a dotted key. Nested objects are returned as
// map[string]any.
func (c Config) Get(key string) (any, error) {
	m := c.toMap()
	var cur any = m
	for _, part := range strings.Split(key, ".") {
		obj, ok := cur.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("unknown key %q", key)
		}
		cur, ok = obj[part]
		if !ok {
			return nil, fmt.Errorf("unknown key %q", key)
		}
	}
	return cur, nil
}

// Set assigns raw (a string from the command line) to a dotted key,
// converting it to the field's type. It returns the updated Config.
func (c Config) Set(key, raw string) (Config, error) {
	m := c.toMap()
	parts := strings.Split(key, ".")
	cur := m
	for _, part := range parts[:len(parts)-1] {
		next, ok := cur[part].(map[string]any)
		if !ok {
			return Config{}, fmt.Errorf("unknown key %q", key)
		}
		cur = next
	}
	last := parts[len(parts)-1]
	existing, ok := cur[last]
	if !ok {
		return Config{}, fmt.Errorf("unknown key %q", key)
	}
	switch existing.(type) {
	case float64: // json numbers
		n, err := strconv.Atoi(raw)
		if err != nil {
			return Config{}, fmt.Errorf("%s expects an integer, got %q", key, raw)
		}
		cur[last] = n
	case bool:
		b, err := strconv.ParseBool(raw)
		if err != nil {
			return Config{}, fmt.Errorf("%s expects true or false, got %q", key, raw)
		}
		cur[last] = b
	case string:
		cur[last] = raw
	default:
		return Config{}, fmt.Errorf("%s is not a settable value", key)
	}
	data, err := json.Marshal(m)
	if err != nil {
		return Config{}, err
	}
	return Parse(data)
}

// ExpandHome replaces a leading ~ with the user's home directory.
func ExpandHome(p string) string {
	if p == "~" || strings.HasPrefix(p, "~/") {
		if home, err := os.UserHomeDir(); err == nil {
			return filepath.Join(home, strings.TrimPrefix(p, "~"))
		}
	}
	return p
}

func (c Config) toMap() map[string]any {
	data, _ := json.Marshal(c)
	var m map[string]any
	_ = json.Unmarshal(data, &m)
	return m
}

func orDefault(s, d string) string {
	if s == "" {
		return d
	}
	return s
}

func currentUser() string {
	if u, err := user.Current(); err == nil && u.Username != "" {
		return u.Username
	}
	return os.Getenv("USER")
}
