package agentimage

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/bytepunx/system-flow/flai/internal/buildinfo"
	"github.com/bytepunx/system-flow/flai/internal/execx"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
)

// ErrNoDocker is an Inspector's answer when the docker client is not on PATH.
var ErrNoDocker = errors.New("docker is not on PATH")

// Info is what docker says of an image on this host: its ID, its labels,
// and its RepoDigests, each <repository>@sha256:… of a registry it was
// pulled from or pushed to.
type Info struct {
	ID          string
	Labels      map[string]string
	RepoDigests []string
}

// Inspector reports whether the image ref is on this host and, when it is,
// its Info: ErrNoDocker when docker is not on PATH, another error when
// docker could not say.
type Inspector func(ref string) (Info, bool, error)

// Docker is the Inspector that asks program, the docker client, through r:
// one docker image inspect per image.
func Docker(r execx.Runner, program string) Inspector {
	return func(ref string) (Info, bool, error) {
		if _, err := r.LookPath(program); err != nil {
			return Info{}, false, fmt.Errorf("inspect the agent image %s: %s is not on PATH (%w): install Docker, or name its client with flai serve agent harness container --program <path>", ref, program, ErrNoDocker)
		}
		out, err := r.Run("", program, "image", "inspect", "--format", "{{json .}}", ref)
		if err != nil {
			// "Error: No such image" or "Error response from daemon: No such
			// image", by the client's version.
			if strings.Contains(strings.ToLower(out+" "+err.Error()), "no such image") {
				return Info{}, false, nil
			}
			said := firstLine(out)
			if said == "" {
				said = firstLine(err.Error())
			}
			return Info{}, false, fmt.Errorf("%s could not say whether the agent image %s is on this host (%s): start Docker, then ask again", program, ref, said)
		}
		var got struct {
			ID          string   `json:"Id"`
			RepoDigests []string `json:"RepoDigests"`
			Config      struct {
				Labels map[string]string `json:"Labels"`
			} `json:"Config"`
		}
		raw, ok := jsonObject(out)
		if ok {
			ok = json.Unmarshal(raw, &got) == nil && got.ID != ""
		}
		if !ok {
			return Info{}, false, fmt.Errorf("%s image inspect printed no image for the agent image %s (%s): check that %s is a docker client, then ask again", program, ref, firstLine(out), program)
		}
		return Info{ID: got.ID, Labels: got.Config.Labels, RepoDigests: got.RepoDigests}, true, nil
	}
}

// State is what an image on this host is to the flai that must run in it.
type State string

// The states of an agent image.
const (
	// StateMissing is an image not on this host.
	StateMissing State = "missing"
	// StateCurrent is an image whose labels are those it must carry.
	StateCurrent State = "current"
	// StateStale is an image with a label other than the one it must carry.
	StateStale State = "stale"
	// StateOwn is the operator's own image, the host's, which is never judged.
	StateOwn State = "own"
)

// ReasonKind is the label that makes an image stale.
type ReasonKind string

// The reasons an image is stale.
const (
	// ReasonFlai is an image built for another flai.
	ReasonFlai ReasonKind = "flai"
	// ReasonRecipe is an image built from another Dockerfile, target, or arguments.
	ReasonRecipe ReasonKind = "recipe"
	// ReasonBase is an image built over another image than the one now under it.
	ReasonBase ReasonKind = "base"
)

// Reason is one reason an image is stale, in a sentence.
type Reason struct {
	Kind    ReasonKind `json:"kind"`
	Message string     `json:"message"`
}

// Labels are the labels an image carries that tie it to what it was built from.
type Labels struct {
	Flai   string `json:"flai,omitempty"`
	Recipe string `json:"recipe,omitempty"`
	Base   string `json:"base,omitempty"`
	Commit string `json:"commit,omitempty"`
}

// Want is what an image's labels must hold to be current, "" for a label
// not compared.
type Want struct {
	Flai   string `json:"flai,omitempty"`
	Recipe string `json:"recipe,omitempty"`
	Base   string `json:"base,omitempty"`
}

// Status is an image, whether it is on this host, its ID and labels there,
// its state, why it is stale, and what its labels must hold.
type Status struct {
	Image
	Present bool     `json:"present"`
	ID      string   `json:"id,omitempty"`
	Labels  Labels   `json:"labels"`
	State   State    `json:"state"`
	Reasons []Reason `json:"reasons,omitempty"`
	Want    Want     `json:"want"`
}

// Check is the status of each of images, for the project at root, asking
// inspect of each and flai for the version they must hold (ProjectFlai),
// which it calls once at most, and only for an image on this host that is
// not the operator's own. An image is missing when not on this host, the
// operator's own when present and named by the host, and otherwise stale
// when a label is not what it must hold (ADR-0295 §2, §3): its flai label,
// a development version when flai is one, else the same version, a leading
// v aside; a layer's recipe label, Recipe over its Dockerfile; and a
// layer's base label, the ID of the harness image of its inner harness among
// images, which is stale when that image is missing. A harness image's
// recipe and base are not compared yet.
func Check(root string, images []Image, inspect Inspector, flai func() (string, error)) ([]Status, error) {
	out := make([]Status, 0, len(images))
	for _, img := range images {
		info, present, err := inspect(img.Ref)
		if err != nil {
			return nil, err
		}
		s := Status{Image: img, Present: present, State: StateMissing}
		if present {
			s.ID = info.ID
			s.Labels = Labels{Flai: info.Labels[LabelFlai], Recipe: info.Labels[LabelRecipe], Base: info.Labels[LabelBase], Commit: info.Labels[LabelCommit]}
		}
		out = append(out, s)
	}
	want, asked := "", false
	for i := range out {
		s := &out[i]
		switch {
		case !s.Present:
			continue
		case s.Source == SourceHost:
			s.State = StateOwn
			continue
		}
		if !asked {
			v, err := flai()
			if err != nil {
				return nil, fmt.Errorf("tell the flai version the agent images must hold: %w", err)
			}
			want, asked = v, true
		}
		s.Want.Flai = want
		if !flaiMatches(want, s.Labels.Flai) {
			s.stale(ReasonFlai, LabelFlai, s.Labels.Flai, "built for flai "+s.Labels.Flai+", not "+shown(want))
		}
		if s.Kind == KindLayer {
			if err := s.judgeLayer(root, out); err != nil {
				return nil, err
			}
		}
		s.State = StateCurrent
		if len(s.Reasons) > 0 {
			s.State = StateStale
		}
	}
	return out, nil
}

// judgeLayer compares the layer image s's recipe label with its Dockerfile's
// recipe under root, and its base label with the ID of the harness image of
// its inner harness among all.
func (s *Status) judgeLayer(root string, all []Status) error {
	path := filepath.Join(root, filepath.FromSlash(s.Dockerfile))
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read %s, the Dockerfile agent.image.layer names, for the recipe its image %s must carry: %w; fix agent.image.layer in %s or add the Dockerfile", path, s.Ref, err, manifest.File)
	}
	s.Want.Recipe = Recipe(data, "", nil)
	if s.Labels.Recipe != s.Want.Recipe {
		s.stale(ReasonRecipe, LabelRecipe, s.Labels.Recipe, "built from another Dockerfile than "+s.Dockerfile+": its recipe is "+s.Labels.Recipe+", not "+s.Want.Recipe)
	}
	for _, h := range all {
		if h.Kind != KindHarness || h.Inner != s.Inner {
			continue
		}
		if !h.Present {
			s.Reasons = append(s.Reasons, Reason{Kind: ReasonBase, Message: "built over " + shown(s.Labels.Base) + ", but the harness image " + h.Ref + " under it is not on this host"})
			return nil
		}
		s.Want.Base = h.ID
		if s.Labels.Base != h.ID {
			s.stale(ReasonBase, LabelBase, s.Labels.Base, "built over "+s.Labels.Base+", not "+h.ID+", the harness image "+h.Ref+" now on this host")
		}
		return nil
	}
	return nil
}

// stale adds the reason kind to s, saying msg, or, when have, s's label
// for kind, is empty, that s carries no label.
func (s *Status) stale(kind ReasonKind, label, have, msg string) {
	if have == "" {
		msg = "carries no " + label + " label, so what it was built from is not known"
	}
	s.Reasons = append(s.Reasons, Reason{Kind: kind, Message: msg})
}

// flaiMatches reports whether an image whose flai label is label holds the
// flai want: a development version for a development want, else the same
// version with any leading v dropped. No label matches nothing.
func flaiMatches(want, label string) bool {
	switch {
	case label == "":
		return false
	case buildinfo.Development(want):
		return buildinfo.Development(label)
	}
	return strings.TrimPrefix(label, "v") == strings.TrimPrefix(want, "v")
}

// shown is v, or "none" when it is empty.
func shown(v string) string {
	if v == "" {
		return "none"
	}
	return v
}

// ProjectFlai is the flai version the agent images of the project at root,
// with manifest m, must hold (ADR-0295 §3): the version its flai.command
// prints with version --json, run through r in root, when m names one, else
// running, the running flai's.
func ProjectFlai(root string, m manifest.Manifest, r execx.Runner, running string) (string, error) {
	command := m.Flai.Command
	if command == "" {
		return running, nil
	}
	path := filepath.Join(root, filepath.FromSlash(command))
	out, err := r.Run(root, path, "version", "--json")
	if err != nil {
		return "", fmt.Errorf("run %s version --json in %s, the flai.command %s names, for the flai version the agent images must hold (%s): fix the tree's flai so that it builds and runs, then ask again", command, root, manifest.File, firstLine(err.Error()))
	}
	var v struct {
		Version string `json:"version"`
	}
	raw, ok := jsonObject(out)
	if ok {
		ok = json.Unmarshal(raw, &v) == nil && v.Version != ""
	}
	if !ok {
		return "", fmt.Errorf("%s version --json in %s, the flai.command %s names, printed no version (%s): fix the tree's flai so that it prints its version, then ask again", command, root, manifest.File, firstLine(out))
	}
	return v.Version, nil
}

// jsonObject is the first JSON object in out, a command's combined output,
// that starts a line, and whether there is one.
func jsonObject(out string) (json.RawMessage, bool) {
	for offset := 0; offset < len(out); {
		line, _, _ := strings.Cut(out[offset:], "\n")
		if strings.HasPrefix(strings.TrimSpace(line), "{") {
			var raw json.RawMessage
			if err := json.NewDecoder(strings.NewReader(out[offset:])).Decode(&raw); err == nil {
				return raw, true
			}
		}
		offset += len(line) + 1
	}
	return nil, false
}

// firstLine is s up to its first newline, trimmed.
func firstLine(s string) string {
	line, _, _ := strings.Cut(strings.TrimSpace(s), "\n")
	return strings.TrimSpace(line)
}
