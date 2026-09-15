// Package template resolves, caches, and renders system-flow template
// repositories. See design/system/template.md and ADR-0005.
package template

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/goccy/go-yaml"
)

// ManifestFile is the template manifest name at the template root.
const ManifestFile = "template.yaml"

// RootDir is the folder inside a template that is rendered into a project.
const RootDir = "root"

// Manifest is the schema of template.yaml.
type Manifest struct {
	Version     string            `yaml:"version"`
	MinFlai     string            `yaml:"min_flai"`
	Description string            `yaml:"description"`
	Variables   []Variable        `yaml:"variables"`
	Layout      map[string]string `yaml:"layout"`
	Render      RenderRules       `yaml:"render"`
	Dashboard   map[string]any    `yaml:"dashboard"`
	Items       map[string]string `yaml:"items"`
	Projects    []ProjectKind     `yaml:"projects"`
}

// Variable is a value the user supplies when rendering.
type Variable struct {
	Name     string `yaml:"name"`
	Prompt   string `yaml:"prompt"`
	Default  string `yaml:"default"`
	Required bool   `yaml:"required"`
}

// RenderRules controls which files are templated and which are skipped.
type RenderRules struct {
	Suffix string   `yaml:"suffix"`
	Ignore []string `yaml:"ignore"`
}

// ProjectKind is an optional sub-project skeleton.
type ProjectKind struct {
	Kind string `yaml:"kind"`
	Path string `yaml:"path"`
}

// LoadManifest reads and validates template.yaml in dir.
func LoadManifest(dir string) (Manifest, error) {
	data, err := os.ReadFile(filepath.Join(dir, ManifestFile))
	if err != nil {
		return Manifest{}, fmt.Errorf("not a template: %w", err)
	}
	var m Manifest
	if err := yaml.UnmarshalWithOptions(data, &m, yaml.Strict()); err != nil {
		return Manifest{}, fmt.Errorf("parse %s: %w", ManifestFile, err)
	}
	if m.Version == "" {
		return Manifest{}, fmt.Errorf("%s: version is required", ManifestFile)
	}
	if m.Render.Suffix == "" {
		m.Render.Suffix = ".tmpl"
	}
	if m.Layout == nil {
		m.Layout = map[string]string{}
	}
	for _, v := range m.Variables {
		if v.Name == "" {
			return Manifest{}, fmt.Errorf("%s: every variable needs a name", ManifestFile)
		}
	}
	if st, err := os.Stat(filepath.Join(dir, RootDir)); err != nil || !st.IsDir() {
		return Manifest{}, fmt.Errorf("template has no %s/ directory", RootDir)
	}
	return m, nil
}

// LayoutKeys returns the layout keys in a stable order.
func (m Manifest) LayoutKeys() []string {
	keys := make([]string, 0, len(m.Layout))
	for k := range m.Layout {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
