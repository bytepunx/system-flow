package kluster

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/bytepunx/system-flow/flai/internal/manifest"
)

// ScaffoldConfig is the kluster config file Init writes, at the repository
// root, since the layout allows no new top-level folder.
const ScaffoldConfig = "kluster.yaml"

// ScaffoldProfile is the profile the scaffold's config names: kluster's own
// default and its smallest documented profile (cert-manager, SPIRE, and
// traefik-tls). kluster v1.3.1 has no profile that installs nothing, and gives
// a config that names none this one anyway, so the scaffold writes it out.
const ScaffoldProfile = "spire"

// Scaffold is what Init wrote: the repository paths, system-flow.yaml first,
// then the config file when it was absent; none when the manifest had a
// kluster: block already.
type Scaffold struct {
	// Wrote are the paths written, from the repository root.
	Wrote []string `json:"wrote"`
	// Config is the config file the block names.
	Config string `json:"config,omitempty"`
	// Kept is set when the config file was there already and was kept.
	Kept bool `json:"kept,omitempty"`
}

// ScaffoldName is the cluster name the scaffold's config gives a project
// named name: its name as a DNS label with -local after it, or local when
// the name makes no DNS label.
func ScaffoldName(name string) string {
	label := manifest.DNSLabel(name)
	if label == "" {
		return "local"
	}
	if len(label) > 63-len("-local") {
		label = label[:63-len("-local")]
		for label != "" && label[len(label)-1] == '-' {
			label = label[:len(label)-1]
		}
	}
	return label + "-local"
}

// ScaffoldText is the scaffold's kluster config for a project named name.
func ScaffoldText(name string) string {
	return fmt.Sprintf(`# The project's local Kubernetes cluster, which flai brings up with kluster
# (design/system/kluster.md). Change the profile, or add addons, from the
# dashboard's Klusters page before the first bring-up.
name: %s
provider: %s
profile: %s
`, ScaffoldName(name), manifest.DefaultKlusterProvider, ScaffoldProfile)
}

// Init gives the project at root, whose manifest is m, a kluster: block that
// names only config: kluster.yaml, and writes kluster.yaml when it is absent,
// keeping it when it is there. The block names no deploy, so flai deploys
// nothing until kluster.deploy is set. system-flow.yaml is edited a line at a
// time, so every other line stays as it was. A manifest with a block already
// is left as it is, and Init writes nothing. It commits nothing: the caller
// commits Wrote.
func Init(root string, m manifest.Manifest) (Scaffold, error) {
	if m.Kluster != nil {
		return Scaffold{Config: m.Kluster.Config}, nil
	}
	path := filepath.Join(root, manifest.File)
	data, err := os.ReadFile(path)
	if err != nil {
		return Scaffold{}, fmt.Errorf("read %s to write its kluster: block: %w", path, err)
	}
	out, err := manifest.EditSettings(data, []manifest.Assignment{{Key: "kluster.config", Value: ScaffoldConfig}}, nil)
	if err != nil {
		return Scaffold{}, fmt.Errorf("write the kluster: block into %s: %w; fix what it names, then turn the kluster host action on again", path, err)
	}
	cfg := filepath.Join(root, ScaffoldConfig)
	s := Scaffold{Config: ScaffoldConfig, Wrote: []string{manifest.File}}
	switch _, err := os.Stat(cfg); {
	case err == nil:
		s.Kept = true
	case errors.Is(err, fs.ErrNotExist):
		if err := os.WriteFile(cfg, []byte(ScaffoldText(m.Name)), 0o644); err != nil {
			return Scaffold{}, fmt.Errorf("write the kluster config %s: %w", cfg, err)
		}
		s.Wrote = append(s.Wrote, ScaffoldConfig)
	default:
		return Scaffold{}, fmt.Errorf("look for the kluster config %s: %w", cfg, err)
	}
	if err := os.WriteFile(path, out, 0o644); err != nil {
		return Scaffold{}, fmt.Errorf("write the kluster: block into %s: %w", path, err)
	}
	return s, nil
}
