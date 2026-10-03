// Package manifest reads system-flow.yaml, the file that marks a conforming
// project. See design/system/project-manifest.md and ADR-0011.
package manifest

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

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
	// Agent is the project's default agent, copied into every story created
	// while it is set (S-0103). flai agent sets it.
	Agent *Agent `yaml:"agent,omitempty" json:"agent,omitempty"`
	// Prime is how flai prime --story builds this project's context packs
	// (ADR-0049).
	Prime Prime `yaml:"prime,omitempty" json:"prime,omitzero"`
	// Issues is how flai check treats the project's open issues (S-0198).
	Issues Issues `yaml:"issues,omitempty" json:"issues,omitzero"`
	// Flai is what the project asks of the flai that reads it (S-0181).
	Flai Requirement `yaml:"flai,omitempty" json:"flai,omitzero"`
}

// Requirement is the oldest flai that may read the project.
type Requirement struct {
	// Minimum is a flai release, X.Y.Z: one that knows every front-matter
	// field the project's items carry. Publishing a flai release whose
	// front-matter fields changed raises it (release.RaiseMinimum).
	Minimum string `yaml:"minimum,omitempty" json:"minimum,omitempty"`
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

// Prime is the project's say about its context packs.
type Prime struct {
	// Budget is the size a story's context pack fits, such as 80KB or
	// 81920 (bytes); empty means flai's default.
	Budget string `yaml:"budget,omitempty" json:"budget,omitempty"`
}

// Issues is the project's say about its open issues.
type Issues struct {
	// StoryAfter is how long an issue may stay open with no open story
	// linking it before flai check warns, a Go duration such as 168h; empty
	// means DefaultStoryAfter, and 0 turns the warning off.
	StoryAfter string `yaml:"story_after,omitempty" json:"story_after,omitempty"`
}

// DefaultStoryAfter is how long an issue may stay open with no open story
// linking it when issues.story_after is not set: 7 days.
const DefaultStoryAfter = 168 * time.Hour

// StoryAfterDuration is issues.story_after as a duration: DefaultStoryAfter
// when it is empty, and zero when the warning is off.
func (i Issues) StoryAfterDuration() (time.Duration, error) {
	s := strings.TrimSpace(i.StoryAfter)
	if s == "" {
		return DefaultStoryAfter, nil
	}
	d, err := time.ParseDuration(s)
	if err != nil {
		return 0, fmt.Errorf("issues.story_after %q is not a duration; write one such as 168h or 24h, or 0 to turn the warning off", i.StoryAfter)
	}
	if d < 0 {
		return 0, fmt.Errorf("issues.story_after %q is negative; write a duration such as 168h or 24h, or 0 to turn the warning off", i.StoryAfter)
	}
	return d, nil
}

// NamedCommand is one command by name: an argument list, run as it stands,
// never through a shell. In an argument, {story} is replaced by the story's
// ID and {root} by the directory it runs in; nothing else is interpreted.
type NamedCommand struct {
	Name    string   `yaml:"name" json:"name"`
	Command []string `yaml:"command" json:"command"`
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
}

// Dashboard is how the flaiover image runs for this project.
type Dashboard struct {
	Image string `yaml:"image"`
	Tag   string `yaml:"tag"`
	Port  int    `yaml:"port"`
	Bind  string `yaml:"bind"` // host address to publish on; default 0.0.0.0
	// Autocommit commits documents saved from the dashboard; nil means true (ADR-0023).
	Autocommit *bool `yaml:"autocommit,omitempty"`
	// NotifyURL, when set, is where the dashboard's server posts new inbox
	// entries (S-0042). flai itself never calls it.
	NotifyURL string `yaml:"notify_url,omitempty"`
}

// Autocommit reports whether documents saved from the dashboard are committed.
func (m Manifest) Autocommit() bool {
	return m.Dashboard.Autocommit == nil || *m.Dashboard.Autocommit
}

// Load parses the manifest at path.
func Load(path string) (Manifest, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Manifest{}, err
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
