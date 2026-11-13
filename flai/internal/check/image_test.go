package check

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/agentimage"
	"github.com/bytepunx/system-flow/flai/internal/config"
	"github.com/bytepunx/system-flow/flai/internal/harness"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// imageManifest is a manifest naming the layer images/agent/Dockerfile, with
// layer: on line 6.
const imageManifest = "name: sample\nkey: SMP\n# the agent\nagent:\n  image:\n    layer: images/agent/Dockerfile\nplanning:\n  layer: no\n"

// imageRepo is a project at a temporary root whose manifest names layer,
// with the Dockerfile written when dockerfile.
func imageRepo(t *testing.T, layer string, dockerfile bool) *workitem.Repo {
	t.Helper()
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, manifest.File), []byte(imageManifest), 0o644); err != nil {
		t.Fatal(err)
	}
	if dockerfile {
		if err := os.MkdirAll(filepath.Join(root, "images", "agent"), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(root, "images", "agent", "Dockerfile"), []byte(layerDockerfile), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	m := manifest.Manifest{Name: "sample", Key: "SMP"}
	if layer != "" {
		m.Agent = &manifest.Agent{Image: &manifest.AgentImage{Layer: layer}}
	}
	return &workitem.Repo{Root: root, Manifest: m}
}

// layerDockerfile is the Dockerfile imageRepo writes.
const layerDockerfile = "ARG BASE\nFROM ${BASE}\n"

// S-0510, S-0735, ADR-0295, ADR-0255: a manifest that names a layer gets an
// agent.image finding for each image its agents need that is not built on
// this host or is stale, each naming flai image status, and a note when
// docker cannot be asked; all of them advisory and unchecked, inside the
// manifest and the layer, since only the host clears them. A layer that is
// not a file is a plain warning. The host's own images are never judged.
func TestAgentImage(t *testing.T) {
	const layer = "images/agent/Dockerfile"
	const image = "flai-agent-smp-claude-code:local"
	harnessImage := harness.DefaultImages[harness.ClaudeCode]
	recipe := agentimage.Recipe([]byte(layerDockerfile), "", nil)
	current := agentimage.Info{ID: "sha256:harness", Labels: map[string]string{agentimage.LabelFlai: "dev"}}
	layered := func(base string) agentimage.Info {
		return agentimage.Info{ID: "sha256:layer", Labels: map[string]string{agentimage.LabelFlai: "dev", agentimage.LabelRecipe: recipe, agentimage.LabelBase: base}}
	}
	for _, c := range []struct {
		name       string
		layer      string
		dockerfile bool
		images     map[string]agentimage.Info // the images on this host
		host       map[string]string
		err        error
		want       [][]string // what each finding's message says, in order
		advisory   bool
	}{
		{name: "no layer"},
		{name: "current", layer: layer, dockerfile: true, images: map[string]agentimage.Info{harnessImage: current, image: layered("sha256:harness")}},
		{name: "not built", layer: layer, dockerfile: true, images: map[string]agentimage.Info{harnessImage: current},
			want: [][]string{{image, layer, "not built", "flai image status", "make agent-images", "scripts/agent-image.sh --layer"}}, advisory: true},
		{name: "stale layer", layer: layer, dockerfile: true, images: map[string]agentimage.Info{harnessImage: current, image: {ID: "sha256:layer", Labels: map[string]string{agentimage.LabelFlai: "dev", agentimage.LabelRecipe: "sha256:other", agentimage.LabelBase: "sha256:harness"}}},
			want: [][]string{{image, layer, "is stale", "built from another Dockerfile", "flai image status", "rebuild it with `make agent-images`"}}, advisory: true},
		{name: "stale harness image", layer: layer, dockerfile: true, images: map[string]agentimage.Info{harnessImage: {ID: "sha256:harness", Labels: map[string]string{agentimage.LabelFlai: "1.2.3"}}, image: layered("sha256:harness")},
			want: [][]string{{"harness image " + harnessImage, "is stale", "built for flai 1.2.3", "flai image status", "make agent-images"}}, advisory: true},
		{name: "harness image not built", layer: layer, dockerfile: true, images: map[string]agentimage.Info{image: layered("sha256:harness")},
			want: [][]string{{"harness image " + harnessImage, "not built", "flai image status", "make agent-images"}, {image, "is stale", "not on this host", "flai image status"}}, advisory: true},
		{name: "the host's own image", layer: layer, dockerfile: true, host: map[string]string{harness.ClaudeCode: "mine:1"}, images: map[string]agentimage.Info{"mine:1": {ID: "sha256:mine"}, image: layered("sha256:mine")}},
		{name: "Dockerfile missing", layer: layer, want: [][]string{{layer, "not a file", manifest.File, "make agent-images"}}},
		{name: "docker not on PATH", layer: layer, dockerfile: true, err: fmt.Errorf("inspect it: %w", ErrNoDocker), want: [][]string{{image, "not checked", "docker is not on PATH"}}, advisory: true},
		{name: "docker cannot say", layer: layer, dockerfile: true, err: errors.New("docker could not say whether the agent image x is on this host (Cannot connect to the Docker daemon)\nmore"), want: [][]string{{image, "not checked", "Cannot connect to the Docker daemon", "start Docker"}}, advisory: true},
	} {
		t.Run(c.name, func(t *testing.T) {
			repo := imageRepo(t, c.layer, c.dockerfile)
			var asked []string
			probe := ImageProbe{
				Inspect: func(ref string) (agentimage.Info, bool, error) {
					asked = append(asked, ref)
					info, present := c.images[ref]
					return info, present, c.err
				},
				HostImages: func(main string) map[string]string {
					if main != repo.Root {
						t.Errorf("the host's images were asked for %s, not the main checkout %s", main, repo.Root)
					}
					return c.host
				},
			}
			res := &Result{}
			AgentImage(res, repo, probe)
			if c.want == nil {
				if len(res.Findings) != 0 || res.Warnings != 0 || res.Advisory != 0 {
					t.Fatalf("want nothing, got %+v", res)
				}
				if c.layer == "" && len(asked) != 0 {
					t.Errorf("a manifest naming no layer asked docker for %v", asked)
				}
				if c.host != nil && slices.Contains(asked, harnessImage) {
					t.Errorf("the default harness image was asked of in place of the host's: %v", asked)
				}
				return
			}
			if len(res.Findings) != len(c.want) || res.Warnings != len(c.want) || res.Errors != 0 {
				t.Fatalf("want %d warnings, got %+v", len(c.want), res)
			}
			advisory := 0
			if c.advisory {
				advisory = len(c.want)
			}
			if res.Advisory != advisory {
				t.Errorf("advisory count %d, want %d", res.Advisory, advisory)
			}
			for i, f := range res.Findings {
				if f.Rule != "agent.image" || f.Level != Warning || f.Path != manifest.File || f.Line != 6 {
					t.Errorf("finding: %+v", f)
				}
				if f.advisory != c.advisory {
					t.Errorf("advisory = %v, want %v: %s", f.advisory, c.advisory, f.Message)
				}
				// Only the host clears an advisory one, so it is marked for a
				// scoped check to leave out unless the manifest or the layer
				// changed (ADR-0255).
				if f.unchecked != c.advisory || (c.advisory && !slices.Equal(f.inside, []string{manifest.File, layer})) {
					t.Errorf("unchecked = %v inside %v, want %v", f.unchecked, f.inside, c.advisory)
				}
				for _, w := range c.want[i] {
					if !strings.Contains(f.Message, w) {
						t.Errorf("message should mention %q: %s", w, f.Message)
					}
				}
				if strings.Contains(f.Message, "\n") {
					t.Errorf("message spans lines: %q", f.Message)
				}
			}
			if res.OK(true) == !c.advisory {
				t.Errorf("--strict passes = %v, want %v", res.OK(true), c.advisory)
			}
			if !c.dockerfile && len(asked) != 0 {
				t.Errorf("a missing Dockerfile still asked docker for %v", asked)
			}
		})
	}
}

// imageRunner answers LookPath and docker image inspect as a host does that
// has docker or not, and has the image or not, or whose daemon is down.
type imageRunner struct {
	noDocker, built, down bool
	ran                   []string
}

func (r *imageRunner) Run(_, name string, args ...string) (string, error) {
	r.ran = append(r.ran, name+" "+strings.Join(args, " "))
	switch {
	case r.down:
		return "Cannot connect to the Docker daemon at unix:///var/run/docker.sock", fmt.Errorf("docker image inspect: exit status 1")
	case !r.built:
		return "Error response from daemon: No such image: x", fmt.Errorf("docker image inspect: exit status 1")
	}
	return `{"Id":"sha256:abc","Config":{"Labels":{"` + agentimage.LabelFlai + `":"dev"}}}`, nil
}

func (r *imageRunner) RunInput(dir, name, _ string, args ...string) (string, error) {
	return r.Run(dir, name, args...)
}

func (r *imageRunner) LookPath(name string) (string, error) {
	if r.noDocker {
		return "", errors.New("not found")
	}
	return "/usr/bin/" + name, nil
}

// S-0735: DockerImage asks docker image inspect for the image and its
// labels, as flai serve and flai image status do, runs the project's flai
// through the same runner, and reads the host's images.
func TestDockerImage(t *testing.T) {
	for _, c := range []struct {
		name    string
		r       *imageRunner
		built   bool
		noDock  bool
		failure bool
	}{
		{name: "built", r: &imageRunner{built: true}, built: true},
		{name: "not built", r: &imageRunner{}},
		{name: "no docker", r: &imageRunner{noDocker: true}, noDock: true},
		{name: "daemon down", r: &imageRunner{built: true, down: true}, failure: true},
	} {
		t.Run(c.name, func(t *testing.T) {
			probe := DockerImage(c.r)
			if probe.Runner != c.r || probe.HostImages == nil {
				t.Errorf("the probe: %+v", probe)
			}
			info, built, err := probe.Inspect("img:local")
			if built != c.built || errors.Is(err, ErrNoDocker) != c.noDock || (err != nil && !errors.Is(err, ErrNoDocker)) != c.failure {
				t.Errorf("got %v, %v", built, err)
			}
			if c.built && (info.ID != "sha256:abc" || info.Labels[agentimage.LabelFlai] != "dev") {
				t.Errorf("info: %+v", info)
			}
			if c.noDock {
				if len(c.r.ran) != 0 {
					t.Errorf("ran %v without docker", c.r.ran)
				}
				return
			}
			if len(c.r.ran) != 1 || c.r.ran[0] != "docker image inspect --format {{json .}} img:local" {
				t.Errorf("ran %v", c.r.ran)
			}
		})
	}
}

// S-0735: HostImages reads the host's container images from the
// configuration file FLAI_CONFIG names, a project's own entry over the
// host-wide one, and names none for a file that is missing, which it does
// not write, or that does not parse.
func TestHostImages(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")
	t.Setenv(config.EnvVar, path)
	own, other := filepath.Join(dir, "own"), filepath.Join(dir, "other")
	if got := HostImages(own); got != nil {
		t.Errorf("a missing file names %v", got)
	}
	if _, err := os.Stat(path); err == nil {
		t.Error("a missing configuration file was written")
	}
	cfg := `{"agent": {"harnesses": {"container": {"images": {"claude-code": "wide:1"}}}, "projects": {"` + own + `": {"harnesses": {"container": {"images": {"claude-code": "own:1"}}}}}}}`
	if err := os.WriteFile(path, []byte(cfg), 0o644); err != nil {
		t.Fatal(err)
	}
	if got := HostImages(own); got[harness.ClaudeCode] != "own:1" || len(got) != 1 {
		t.Errorf("the project's own images: %v", got)
	}
	if got := HostImages(other); got[harness.ClaudeCode] != "wide:1" || len(got) != 1 {
		t.Errorf("another project's images: %v", got)
	}
	if err := os.WriteFile(path, []byte(`{"unknown": true}`), 0o644); err != nil {
		t.Fatal(err)
	}
	if got := HostImages(own); got != nil {
		t.Errorf("a file that does not parse names %v", got)
	}
}

func TestLayerLine(t *testing.T) {
	path := filepath.Join(t.TempDir(), manifest.File)
	for _, c := range []struct {
		yaml string
		line int
	}{
		{imageManifest, 6},
		{"name: x\nplanning:\n  layer: no\n", 1},
		{"agent:\n  harness: claude-code\n  layer: no\n", 1},
	} {
		if err := os.WriteFile(path, []byte(c.yaml), 0o644); err != nil {
			t.Fatal(err)
		}
		if got := layerLine(path); got != c.line {
			t.Errorf("layerLine(%q) = %d, want %d", c.yaml, got, c.line)
		}
	}
	if got := layerLine(filepath.Join(t.TempDir(), "absent.yaml")); got != 1 {
		t.Errorf("absent file: %d", got)
	}
}
