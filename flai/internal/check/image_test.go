package check

import (
	"errors"
	"fmt"
	"maps"
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
			want: [][]string{{image, layer, "not built", "flai image status", "flai image ensure", "flai image build --layer"}}, advisory: true},
		{name: "stale layer", layer: layer, dockerfile: true, images: map[string]agentimage.Info{harnessImage: current, image: {ID: "sha256:layer", Labels: map[string]string{agentimage.LabelFlai: "dev", agentimage.LabelRecipe: "sha256:other", agentimage.LabelBase: "sha256:harness"}}},
			want: [][]string{{image, layer, "is stale", "built from another Dockerfile", "flai image status", "rebuild it with `flai image ensure`"}}, advisory: true},
		{name: "stale harness image", layer: layer, dockerfile: true, images: map[string]agentimage.Info{harnessImage: {ID: "sha256:harness", Labels: map[string]string{agentimage.LabelFlai: "1.2.3"}}, image: layered("sha256:harness")},
			want: [][]string{{"harness image " + harnessImage, "is stale", "built for flai 1.2.3", "flai image status", "flai image ensure"}}, advisory: true},
		{name: "harness image not built", layer: layer, dockerfile: true, images: map[string]agentimage.Info{image: layered("sha256:harness")},
			want: [][]string{{"harness image " + harnessImage, "not built", "flai image status", "flai image ensure"}, {image, "is stale", "not on this host", "flai image status"}}, advisory: true},
		{name: "the host's own image", layer: layer, dockerfile: true, host: map[string]string{harness.ClaudeCode: "mine:1"}, images: map[string]agentimage.Info{"mine:1": {ID: "sha256:mine"}, image: layered("sha256:mine")}},
		{name: "Dockerfile missing", layer: layer, want: [][]string{{layer, "not a file", manifest.File, "flai image ensure"}}},
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

// S-0735: a check in a story's worktree judges the layer by the main
// checkout's Dockerfile, as flai serve's start and flai image status do, so
// a worktree whose Dockerfile differs from main's still finds main's image
// current.
func TestAgentImageJudgesTheLayerFromTheMainCheckout(t *testing.T) {
	const image = "flai-agent-smp-claude-code:local"
	main := imageRepo(t, "images/agent/Dockerfile", true)
	repo := imageRepo(t, "images/agent/Dockerfile", true)
	repo.MainRoot = main.Root
	if err := os.WriteFile(filepath.Join(repo.Root, "images", "agent", "Dockerfile"), []byte(layerDockerfile+"RUN true\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	recipe := agentimage.Recipe([]byte(layerDockerfile), "", nil)
	images := map[string]agentimage.Info{
		harness.DefaultImages[harness.ClaudeCode]: {ID: "sha256:harness", Labels: map[string]string{agentimage.LabelFlai: "dev"}},
		image: {ID: "sha256:layer", Labels: map[string]string{agentimage.LabelFlai: "dev", agentimage.LabelRecipe: recipe, agentimage.LabelBase: "sha256:harness"}},
	}
	probe := ImageProbe{Inspect: func(ref string) (agentimage.Info, bool, error) {
		info, ok := images[ref]
		return info, ok, nil
	}}
	res := &Result{}
	AgentImage(res, repo, probe)
	if len(res.Findings) != 0 {
		t.Errorf("the layer built from the main checkout's Dockerfile should be current: %+v", res.Findings)
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

// S-0745, ADR-0295: each name agent.image.args gives that no ARG of the
// layer declares, in any stage, is a plain agent.image warning on its line of
// the manifest, which a story can fix, so it is neither advisory nor left out
// of a scoped check; a layer with no args, or with args it declares, gets
// none.
func TestAgentImageArgs(t *testing.T) {
	const layer = "images/agent/Dockerfile"
	// GO_VERSION's entry is on line 8, PNPM_VERSION's on line 9, NODE's on 10.
	const mf = "name: sample\nkey: SMP\nagent:\n  image:\n    layer: images/agent/Dockerfile\n    # the pins\n    args:\n      GO_VERSION: 1.26.0\n      PNPM_VERSION: 10.34.5\n      NODE: \"22\"\n"
	const dockerfile = "# syntax=docker/dockerfile:1\nARG BASE\nFROM ${BASE} AS build\narg GO_VERSION=1.25.0 \\\n# the pnpm pin\n    PNPM_VERSION=\"10 34\"\nFROM build\nARG NODE\n"
	for _, c := range []struct {
		name       string
		args       map[string]string
		dockerfile string
		want       map[string]int // the undeclared names, and the line of each
	}{
		{name: "no args", dockerfile: "ARG BASE\nFROM ${BASE}\n"},
		{name: "all declared", args: map[string]string{"GO_VERSION": "1.26.0", "PNPM_VERSION": "10.34.5", "NODE": "22"}, dockerfile: dockerfile},
		{name: "none declared", args: map[string]string{"GO_VERSION": "1.26.0", "PNPM_VERSION": "10.34.5", "NODE": "22"}, dockerfile: "ARG BASE\nFROM ${BASE}\n# ARG NODE\nRUN echo GO_VERSION\n", want: map[string]int{"GO_VERSION": 8, "PNPM_VERSION": 9, "NODE": 10}},
		{name: "an arg not in the manifest's file", args: map[string]string{"RUST": "1"}, dockerfile: dockerfile, want: map[string]int{"RUST": 5}},
	} {
		t.Run(c.name, func(t *testing.T) {
			root := t.TempDir()
			if err := os.MkdirAll(filepath.Join(root, "images", "agent"), 0o755); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(filepath.Join(root, manifest.File), []byte(mf), 0o644); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(filepath.Join(root, "images", "agent", "Dockerfile"), []byte(c.dockerfile), 0o644); err != nil {
				t.Fatal(err)
			}
			repo := &workitem.Repo{Root: root, Manifest: manifest.Manifest{Name: "sample", Key: "SMP", Agent: &manifest.Agent{Image: &manifest.AgentImage{Layer: layer, Args: c.args}}}}
			res := &Result{}
			// Docker is not asked here: its note is advisory, and set aside.
			AgentImage(res, repo, ImageProbe{Inspect: func(string) (agentimage.Info, bool, error) {
				return agentimage.Info{}, false, ErrNoDocker
			}})
			got := map[string]int{}
			var order []string
			for _, f := range res.Findings {
				if f.advisory {
					continue
				}
				if f.Rule != "agent.image" || f.Level != Warning || f.Path != manifest.File || f.unchecked || f.inside != nil {
					t.Errorf("finding: %+v", f)
				}
				name, _, _ := strings.Cut(strings.TrimPrefix(f.Message, "agent.image.args names "), ",")
				for _, w := range []string{"no ARG of the layer " + layer, "add `ARG " + name + "`", "leave " + name + " out of agent.image.args"} {
					if !strings.Contains(f.Message, w) {
						t.Errorf("message should mention %q: %s", w, f.Message)
					}
				}
				got[name] = f.Line
				order = append(order, name)
			}
			if !maps.Equal(got, c.want) {
				t.Errorf("undeclared args on lines %v, want %v", got, c.want)
			}
			if !slices.IsSorted(order) {
				t.Errorf("findings not by name: %v", order)
			}
			if res.Warnings != len(res.Findings) || res.Advisory != len(res.Findings)-len(c.want) {
				t.Errorf("counts: %+v", res)
			}
		})
	}
}

// S-0745: dockerfileArgs reads the names every ARG declares.
func TestDockerfileArgs(t *testing.T) {
	for _, c := range []struct {
		dockerfile string
		want       []string
	}{
		{"", nil},
		{"ARG BASE\nFROM ${BASE}\n", []string{"BASE"}},
		{"FROM x\nArg A=1 B C=\"x y\"\n", []string{"A", "B", "C"}},
		{"ARG A=1 \\\n  # comment\n  B=2\r\nRUN echo C=3\n", []string{"A", "B"}},
		{"# ARG A\nRUN true # ARG B\nENV D=1\n", nil},
		{"FROM a AS one\nARG ONE\nFROM b\narg TWO='a b'\n", []string{"ONE", "TWO"}},
	} {
		got := slices.Sorted(maps.Keys(dockerfileArgs([]byte(c.dockerfile))))
		if !slices.Equal(got, c.want) {
			t.Errorf("dockerfileArgs(%q) = %v, want %v", c.dockerfile, got, c.want)
		}
	}
}

// S-0745, ADR-0295: flai check refuses a malformed agent.image.args as a
// manifest.agent error, and passes a good one.
func TestManifestAgentImageArgs(t *testing.T) {
	for _, c := range []struct {
		name string
		args string
		want string
	}{
		{"good", "      GO_VERSION: 1.26.0\n", ""},
		{"a bad name", "      GO-VERSION: 1.26.0\n", `"GO-VERSION" is not a Dockerfile ARG name`},
		{"BASE", "      BASE: debian\n", "leave BASE out of agent.image.args"},
		{"a value on two lines", "      GO_VERSION: |\n        1.26.0\n        RUN x\n", "GO_VERSION spans lines"},
	} {
		t.Run(c.name, func(t *testing.T) {
			repo := planningProject(t, "agent:\n  image:\n    layer: Dockerfile\n    args:\n"+c.args)
			got := findings(t, repo, "manifest.agent")
			if c.want == "" {
				if len(got) != 0 {
					t.Errorf("good args refused: %+v", got)
				}
				return
			}
			if len(got) != 1 || got[0].Level != Error || !strings.Contains(got[0].Message, c.want) || got[0].Line != 8 {
				t.Errorf("want one error on line 8 saying %q, got %+v", c.want, got)
			}
		})
	}
	repo := planningProject(t, "agent:\n  image:\n    args:\n      GO_VERSION: 1.26.0\n")
	if got := findings(t, repo, "manifest.agent"); len(got) != 1 || !strings.Contains(got[0].Message, "args name no layer to build") {
		t.Errorf("args without a layer: %+v", got)
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
	// S-0745: an arg's line is its entry under args:, else the layer's.
	if err := os.WriteFile(path, []byte("agent:\n  image:\n    args:\n      layer: x\n      GO: 1\n    layer: Dockerfile\nother:\n  GO: 2\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	for name, want := range map[string]int{"GO": 5, "NODE": 6} {
		if got := argLine(path, name); got != want {
			t.Errorf("argLine(%s) = %d, want %d", name, got, want)
		}
	}
	if got := layerLine(path); got != 6 {
		t.Errorf("layerLine under args: = %d, want 6", got)
	}
	if got := layerLine(filepath.Join(t.TempDir(), "absent.yaml")); got != 1 {
		t.Errorf("absent file: %d", got)
	}
}
