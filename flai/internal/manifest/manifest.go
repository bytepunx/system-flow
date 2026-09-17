// Package manifest reads system-flow.yaml, the file that marks a conforming
// project. See design/system/project-manifest.md and ADR-0011.
package manifest

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/goccy/go-yaml"
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
	return m, nil
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
