// Package agentimage says which images a project's container agents run in,
// where each image's name comes from, and whether the image on this host is
// the one the project's flai must run them in (ADR-0295): missing, current,
// stale by its labels, or the operator's own, which is never judged.
//
// It decides the image an inner harness runs in as flai serve's
// containerImageHost does, and runs docker only through the Inspector a
// caller gives it.
package agentimage

import (
	"crypto/sha256"
	"encoding/hex"
	"maps"
	"slices"

	"github.com/bytepunx/system-flow/flai/internal/harness"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
)

// The labels that tie an agent image to what it was built from (ADR-0295 §2, §3).
const (
	// LabelFlai holds the version the image's flai reports.
	LabelFlai = "dev.system-flow.flai"
	// LabelRecipe holds the image's recipe, as Recipe computes it.
	LabelRecipe = "dev.system-flow.recipe"
	// LabelBase holds the ID of the image the image was built over.
	LabelBase = "dev.system-flow.base"
	// LabelCommit holds the commit of the tree a development image was built from.
	LabelCommit = "dev.system-flow.commit"
)

// Kind is what an image is to the agents: an inner harness's or the project's layer.
type Kind string

// The kinds of agent image.
const (
	// KindHarness is an inner harness's image.
	KindHarness Kind = "harness"
	// KindLayer is the project's layer, agent.image.layer, built over a harness image.
	KindLayer Kind = "layer"
)

// Source is where an image's name came from.
type Source string

// The sources of an image's name.
const (
	// SourceLayer is the manifest's agent.image.layer.
	SourceLayer Source = "layer"
	// SourceHost is the host's flai serve agent harness container --image <inner>=<image>.
	SourceHost Source = "host"
	// SourceDefault is harness.DefaultImages.
	SourceDefault Source = "default"
)

// Image is one image the agents of an inner harness need: its kind, the
// inner harness, its name, where the name came from, and for a layer the
// Dockerfile, from the repository root, it is built from.
type Image struct {
	Kind       Kind   `json:"kind"`
	Inner      string `json:"inner"`
	Ref        string `json:"ref"`
	Source     Source `json:"source"`
	Dockerfile string `json:"dockerfile,omitempty"`
}

// Harness is the image the inner harness inner runs in under no layer: the
// host's image for it, from hostImages, else its default. Its Ref is "" for
// an inner harness with neither.
func Harness(hostImages map[string]string, inner string) Image {
	img := Image{Kind: KindHarness, Inner: inner}
	switch {
	case hostImages[inner] != "":
		img.Ref, img.Source = hostImages[inner], SourceHost
	case harness.DefaultImages[inner] != "":
		img.Ref, img.Source = harness.DefaultImages[inner], SourceDefault
	}
	return img
}

// For is the image an agent whose inner harness is inner runs in, for a
// project with manifest m on a host whose container harness names
// hostImages: the project's layer image (Kind KindLayer) when m names
// agent.image.layer and inner has a default image, which the layer is built
// over, else Harness. Its Ref is "" for an inner harness with no image.
func For(m manifest.Manifest, hostImages map[string]string, inner string) Image {
	if layer := m.Layer(); layer != "" && harness.DefaultImages[inner] != "" {
		return Image{Kind: KindLayer, Inner: inner, Ref: m.LayerImage(inner), Source: SourceLayer, Dockerfile: layer}
	}
	return Harness(hostImages, inner)
}

// Needed is every image the agents in agents need, for a project with
// manifest m on a host whose container harness names hostImages: for each
// inner harness an agent on the container harness names, and claude-code
// when m names a layer, its harness image, then its layer image when it has
// one. They come by inner harness, harness before layer, each once; an
// inner harness with no image has none.
func Needed(m manifest.Manifest, agents []*manifest.Agent, hostImages map[string]string) []Image {
	inners := map[string]bool{}
	for _, a := range agents {
		if a != nil && a.Harness == harness.Container {
			inners[innerOf(a)] = true
		}
	}
	if m.Layer() != "" {
		inners[harness.ClaudeCode] = true
	}
	var out []Image
	for _, inner := range slices.Sorted(maps.Keys(inners)) {
		if h := Harness(hostImages, inner); h.Ref != "" {
			out = append(out, h)
		}
		if img := For(m, hostImages, inner); img.Kind == KindLayer {
			out = append(out, img)
		}
	}
	return out
}

// innerOf is the inner harness agent's config.inner names, claude-code when
// it names none, as flai serve reads it.
func innerOf(a *manifest.Agent) string {
	if inner := a.Config["inner"]; inner != "" {
		return inner
	}
	return harness.ClaudeCode
}

// builtinRoles are the roles every project has, whose agents Agents lists.
var builtinRoles = []string{"planner", "orchestrator", "analyzer", "curator"}

// Agents is the agents manifest m starts: the project's agent, the built-in
// roles' (the planner's, the orchestrator's, the analyzer's, the curator's),
// and each role's m.Roles names, by name; an agent that sets nothing is left
// out. A caller adds the open stories' own agents to give Needed.
func Agents(m manifest.Manifest) []*manifest.Agent {
	var out []*manifest.Agent
	if m.Agent != nil {
		out = append(out, m.Agent)
	}
	for _, name := range append(slices.Clone(builtinRoles), slices.Sorted(maps.Keys(m.Roles))...) {
		if a := m.RoleAgent(name); a != nil {
			out = append(out, a)
		}
	}
	return out
}

// Recipe is the recipe label an image built from dockerfile at target with
// the build arguments args must carry (ADR-0295 §2): "sha256:" and the
// lowercase hex SHA-256 of the Dockerfile's bytes, then "\n--target=" and
// target, then "\n", then for each argument by name in byte order, leaving
// out BASE and FLAI_VERSION, its name, "=", its value, and "\n". With no
// target and no arguments, a shell reproduces it with
// { cat Dockerfile; printf '\n--target=\n'; } | sha256sum.
func Recipe(dockerfile []byte, target string, args map[string]string) string {
	h := sha256.New()
	h.Write(dockerfile)
	h.Write([]byte("\n--target=" + target + "\n"))
	for _, name := range slices.Sorted(maps.Keys(args)) {
		if name == "BASE" || name == "FLAI_VERSION" {
			continue
		}
		h.Write([]byte(name + "=" + args[name] + "\n"))
	}
	return "sha256:" + hex.EncodeToString(h.Sum(nil))
}
