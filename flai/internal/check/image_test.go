package check

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

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
		if err := os.WriteFile(filepath.Join(root, "images", "agent", "Dockerfile"), []byte("ARG BASE\nFROM ${BASE}\n"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	m := manifest.Manifest{Name: "sample", Key: "SMP"}
	if layer != "" {
		m.Agent = &manifest.Agent{Image: &manifest.AgentImage{Layer: layer}}
	}
	return &workitem.Repo{Root: root, Manifest: m}
}

func TestAgentImage(t *testing.T) {
	const layer = "images/agent/Dockerfile"
	const image = "flai-agent-smp-claude-code:local"
	for _, c := range []struct {
		name       string
		layer      string
		dockerfile bool
		built      bool
		err        error
		want       []string // what the one finding's message says; none for no finding
		advisory   bool
	}{
		{name: "no layer", built: false},
		{name: "built", layer: layer, dockerfile: true, built: true},
		{name: "not built", layer: layer, dockerfile: true, want: []string{image, layer, "not built", "make agent-images", "scripts/agent-image.sh --layer"}, advisory: true},
		{name: "Dockerfile missing", layer: layer, built: true, want: []string{layer, "not a file", manifest.File, "make agent-images"}},
		{name: "docker not on PATH", layer: layer, dockerfile: true, err: ErrNoDocker, want: []string{image, "not checked", "docker is not on PATH"}, advisory: true},
		{name: "docker cannot say", layer: layer, dockerfile: true, err: errors.New("docker image inspect: exit status 1\nCannot connect to the Docker daemon"), want: []string{image, "not checked", "exit status 1", "start Docker"}, advisory: true},
	} {
		t.Run(c.name, func(t *testing.T) {
			repo := imageRepo(t, c.layer, c.dockerfile)
			var asked []string
			probe := func(img string) (bool, error) {
				asked = append(asked, img)
				return c.built, c.err
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
				return
			}
			if len(res.Findings) != 1 || res.Warnings != 1 || res.Errors != 0 {
				t.Fatalf("want one warning, got %+v", res)
			}
			f := res.Findings[0]
			if f.Rule != "agent.image" || f.Level != Warning || f.Path != manifest.File || f.Line != 6 {
				t.Errorf("finding: %+v", f)
			}
			if f.advisory != c.advisory || (res.Advisory == 1) != c.advisory {
				t.Errorf("advisory = %v (count %d), want %v", f.advisory, res.Advisory, c.advisory)
			}
			// Only the host clears an advisory one, so it is marked for a
			// scoped check to leave out unless the manifest or the layer
			// changed (ADR-0255).
			if f.unchecked != c.advisory || (c.advisory && !slices.Equal(f.inside, []string{manifest.File, layer})) {
				t.Errorf("unchecked = %v inside %v, want %v", f.unchecked, f.inside, c.advisory)
			}
			if res.OK(true) == !c.advisory {
				t.Errorf("--strict passes = %v, want %v", res.OK(true), c.advisory)
			}
			for _, w := range c.want {
				if !strings.Contains(f.Message, w) {
					t.Errorf("message should mention %q: %s", w, f.Message)
				}
			}
			if strings.Contains(f.Message, "\n") {
				t.Errorf("message spans lines: %q", f.Message)
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
	return "sha256:abc", nil
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
			built, err := DockerImage(c.r)("img:local")
			if built != c.built || errors.Is(err, ErrNoDocker) != c.noDock || (err != nil && !errors.Is(err, ErrNoDocker)) != c.failure {
				t.Errorf("got %v, %v", built, err)
			}
			if c.noDock {
				if len(c.r.ran) != 0 {
					t.Errorf("ran %v without docker", c.r.ran)
				}
				return
			}
			if len(c.r.ran) != 1 || c.r.ran[0] != "docker image inspect --format {{.Id}} img:local" {
				t.Errorf("ran %v", c.r.ran)
			}
		})
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
