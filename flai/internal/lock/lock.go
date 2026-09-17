// Package lock reads and writes system-flow.lock.yaml, the record of what
// the template rendered (ADR-0015).
package lock

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/goccy/go-yaml"
)

// File is the lock file name at the project root.
const File = "system-flow.lock.yaml"

// Lock is the schema.
type Lock struct {
	Template Template          `yaml:"template"`
	Files    map[string]string `yaml:"files"` // project-relative path -> sha256
}

// Template is the source that was applied.
type Template struct {
	Repo    string `yaml:"repo"`
	Ref     string `yaml:"ref"`
	Version string `yaml:"version"`
	Applied string `yaml:"applied"`
}

// Hash returns the sha256 hex of content.
func Hash(content []byte) string {
	sum := sha256.Sum256(content)
	return hex.EncodeToString(sum[:])
}

// HashFile hashes a file on disk.
func HashFile(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return Hash(data), nil
}

// Load reads the lock at root; a missing file returns (nil, nil).
func Load(root string) (*Lock, error) {
	data, err := os.ReadFile(filepath.Join(root, File))
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var l Lock
	if err := yaml.Unmarshal(data, &l); err != nil {
		return nil, fmt.Errorf("%s: %w", File, err)
	}
	if l.Files == nil {
		l.Files = map[string]string{}
	}
	return &l, nil
}

// Save writes the lock deterministically (sorted paths).
func Save(root string, l *Lock) error {
	var b strings.Builder
	b.WriteString("# Written by flai. Hashes of what the template rendered, used by flai upgrade (ADR-0015).\n")
	fmt.Fprintf(&b, "template:\n  repo: %q\n  ref: %q\n  version: %q\n  applied: %s\n", l.Template.Repo, l.Template.Ref, l.Template.Version, l.Template.Applied)
	if len(l.Files) == 0 {
		b.WriteString("files: {}\n")
	} else {
		b.WriteString("files:\n")
		paths := make([]string, 0, len(l.Files))
		for p := range l.Files {
			paths = append(paths, p)
		}
		sort.Strings(paths)
		for _, p := range paths {
			fmt.Fprintf(&b, "  %q: %s\n", p, l.Files[p])
		}
	}
	return os.WriteFile(filepath.Join(root, File), []byte(b.String()), 0o644)
}
