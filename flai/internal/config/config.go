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
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
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
	// Checks are the commands a story in review is checked with, once the
	// checks action is enabled (S-0082). When Commands is empty the
	// manifest's own checks: are used instead; this is not among Keys, and
	// flai serve checks manages it.
	Checks ChecksConfig `json:"checks,omitzero"`
	// ImportRoots are the folders flai serve looks in for repositories that
	// are not system-flow projects yet, which the board then offers to import
	// (S-0098). Naming a folder is the operator's consent to importing, and
	// testing, what is in it. Absolute paths; flai serve import manages it,
	// so it is not among Keys.
	ImportRoots []string `json:"import_roots,omitempty"`
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
	// AutoRestarts is how many times flai serve restarts a story's agent on
	// its own after it ends with its story in progress (ADR-0108). Nil means
	// DefaultAutoRestarts, 0 turns the restart off, and a negative count is
	// refused; read it through AutoRestartLimit.
	AutoRestarts *int `json:"auto_restarts,omitempty"`
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

// checkAutoRestarts refuses a negative count of automatic restarts.
func checkAutoRestarts(n int) error {
	if n < 0 {
		return fmt.Errorf("agent.auto_restarts is %d: give how many times flai serve restarts a story's agent on its own, 0 or more, with flai serve agent set --auto-restarts <n> (0 turns it off; unset is %d)", n, DefaultAutoRestarts)
	}
	return nil
}

// HarnessHost is the program a harness is and the operator's arguments to it.
type HarnessHost struct {
	Program string `json:"program,omitempty"`
	// Args replace the adapter's defaults when set, even to an empty list;
	// nil leaves them.
	Args *[]string `json:"args,omitempty"`
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
	// PushKey is an SSH private key on this host that the container may push
	// acceptances with (ADR-0026). Empty, the default, gives it nothing. It is
	// a fact about one machine, so it lives here and never in the manifest.
	PushKey string `json:"push_key"`
	// PushKnownHosts names the file the remote's host keys are taken from
	// instead of this user's and the system's known_hosts.
	PushKnownHosts string `json:"push_known_hosts"`
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

// Parse decodes JSON into a Config, rejecting unknown keys and a negative
// agent.auto_restarts.
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
	return cfg, nil
}

// Save writes cfg to path atomically, creating parent directories.
func Save(path string, cfg Config) error {
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
		"dashboard.image", "dashboard.tag", "dashboard.port", "dashboard.bind", "dashboard.no_restart", "dashboard.push_key", "dashboard.push_known_hosts",
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
