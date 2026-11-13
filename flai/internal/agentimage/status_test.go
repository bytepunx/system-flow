package agentimage

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/harness"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
)

const (
	dockerfile = "images/agent/Dockerfile"
	hostImage  = "host-agent:1"
	layerImage = "flai-agent-t-claude-code:local"
)

var defaultImage = harness.DefaultImages[harness.ClaudeCode]

// layered is a manifest with key t whose agent.image.layer names dockerfile.
func layered() manifest.Manifest {
	return manifest.Manifest{Key: "t", Agent: &manifest.Agent{Image: &manifest.AgentImage{Layer: dockerfile}}}
}

// S-0735: For answers what flai serve's containerImageHost answers for each
// case of TestContainerImageHostPutsTheLayerImageInACopyOfTheHostsImages,
// and that the answer is the layer's, with the name's source.
func TestForDecidesTheImageAnInnerHarnessRunsIn(t *testing.T) {
	host := map[string]string{harness.ClaudeCode: hostImage}
	bare := manifest.Manifest{Key: "t"}
	for _, c := range []struct {
		name  string
		m     manifest.Manifest
		host  map[string]string
		inner string
		want  Image
	}{
		{"host image", bare, host, harness.ClaudeCode, Image{Kind: KindHarness, Inner: harness.ClaudeCode, Ref: hostImage, Source: SourceHost}},
		{"default", bare, nil, harness.ClaudeCode, Image{Kind: KindHarness, Inner: harness.ClaudeCode, Ref: defaultImage, Source: SourceDefault}},
		{"layer over the host's", layered(), host, harness.ClaudeCode, Image{Kind: KindLayer, Inner: harness.ClaudeCode, Ref: layerImage, Source: SourceLayer, Dockerfile: dockerfile}},
		{"layer over the default", layered(), nil, harness.ClaudeCode, Image{Kind: KindLayer, Inner: harness.ClaudeCode, Ref: layerImage, Source: SourceLayer, Dockerfile: dockerfile}},
		{"no harness image", layered(), nil, harness.Command, Image{Kind: KindHarness, Inner: harness.Command}},
		{"no default image but the host's", layered(), map[string]string{harness.Command: "mine:1"}, harness.Command, Image{Kind: KindHarness, Inner: harness.Command, Ref: "mine:1", Source: SourceHost}},
	} {
		t.Run(c.name, func(t *testing.T) {
			if got := For(c.m, c.host, c.inner); got != c.want {
				t.Errorf("For = %+v, want %+v", got, c.want)
			}
		})
	}
	if len(host) != 1 || host[harness.ClaudeCode] != hostImage {
		t.Errorf("the host's images were changed: %v", host)
	}
}

// S-0735: Needed lists each image the agents need once, by inner harness,
// the harness image before the layer, and the claude-code harness image and
// layer under a layer even when no agent runs in a container.
func TestNeededListsEachImageTheAgentsNeedOnce(t *testing.T) {
	container := &manifest.Agent{Harness: harness.Container}
	other := &manifest.Agent{Harness: harness.Container, Config: map[string]string{"inner": "other"}}
	native := &manifest.Agent{Harness: harness.ClaudeCode}
	host := map[string]string{"other": "other:1"}
	harnessImg := Image{Kind: KindHarness, Inner: harness.ClaudeCode, Ref: defaultImage, Source: SourceDefault}
	layerImg := Image{Kind: KindLayer, Inner: harness.ClaudeCode, Ref: layerImage, Source: SourceLayer, Dockerfile: dockerfile}
	otherImg := Image{Kind: KindHarness, Inner: "other", Ref: "other:1", Source: SourceHost}
	for _, c := range []struct {
		name   string
		m      manifest.Manifest
		agents []*manifest.Agent
		want   []Image
	}{
		{"no container agent, no layer", manifest.Manifest{Key: "t"}, []*manifest.Agent{native, nil}, nil},
		{"container agents, no layer", manifest.Manifest{Key: "t"}, []*manifest.Agent{container, other, container, native}, []Image{harnessImg, otherImg}},
		{"layer, no container agent", layered(), []*manifest.Agent{native}, []Image{harnessImg, layerImg}},
		{"layer and container agents", layered(), []*manifest.Agent{other, container, container}, []Image{harnessImg, layerImg, otherImg}},
		{"an inner harness with no image", manifest.Manifest{Key: "t"}, []*manifest.Agent{{Harness: harness.Container, Config: map[string]string{"inner": "bare"}}}, nil},
	} {
		t.Run(c.name, func(t *testing.T) {
			if got := Needed(c.m, c.agents, host); !reflect.DeepEqual(got, c.want) {
				t.Errorf("Needed =\n%+v\nwant\n%+v", got, c.want)
			}
		})
	}
}

// S-0735: Agents holds the project's agent and the agents of the built-in
// roles and the defined ones, leaving out those that set nothing.
func TestAgentsAreTheManifestsAgents(t *testing.T) {
	if got := Agents(manifest.Manifest{}); len(got) != 0 {
		t.Errorf("a manifest with no agents gives %+v", got)
	}
	m := manifest.Manifest{Roles: map[string]manifest.RoleSettings{"reviewer": {Agent: &manifest.Agent{Harness: harness.Container}}}}
	m.Planning.Agent = &manifest.Agent{Harness: harness.Container, Config: map[string]string{"inner": "planner-inner"}}
	got := Agents(m)
	if len(got) != 2 || got[0].Config["inner"] != "planner-inner" || got[1].Harness != harness.Container {
		t.Errorf("Agents = %+v, want the planner's and the reviewer's", got)
	}
	m.Agent = &manifest.Agent{Harness: harness.ClaudeCode}
	if got := Agents(m); len(got) != 6 || got[0] != m.Agent {
		t.Errorf("Agents = %+v, want the project's and five roles'", got)
	}
}

// S-0735, ADR-0295 §2: the recipe is the SHA-256 of the Dockerfile, the
// target, and the arguments sorted by name, BASE and FLAI_VERSION left out.
func TestRecipeIgnoresBaseAndFlaiVersionAndSortsTheArguments(t *testing.T) {
	file := []byte("FROM x\n")
	sum := func(s string) string {
		h := sha256.Sum256([]byte(s))
		return "sha256:" + hex.EncodeToString(h[:])
	}
	if got, want := Recipe(file, "", nil), sum("FROM x\n\n--target=\n"); got != want {
		t.Errorf("Recipe = %s, want %s", got, want)
	}
	args := map[string]string{"Z": "1", "A": "2", "BASE": "b", "FLAI_VERSION": "v"}
	if got, want := Recipe(file, "release", args), sum("FROM x\n\n--target=release\nA=2\nZ=1\n"); got != want {
		t.Errorf("Recipe = %s, want %s", got, want)
	}
	if Recipe(file, "", map[string]string{"BASE": "other"}) != Recipe(file, "", nil) {
		t.Error("BASE changed the recipe")
	}
	if Recipe(file, "", map[string]string{"GO": "1.24"}) == Recipe(file, "", nil) {
		t.Error("an argument did not change the recipe")
	}
}

// fakeImages is an Inspector over images by name, counting its calls.
type fakeImages struct {
	images map[string]Info
	calls  int
}

func (f *fakeImages) inspect(ref string) (Info, bool, error) {
	f.calls++
	info, ok := f.images[ref]
	return info, ok, nil
}

// counter is a flai callback answering version, counting its calls.
type counter struct {
	version string
	calls   int
}

func (c *counter) flai() (string, error) {
	c.calls++
	return c.version, nil
}

// labels are the labels flai, recipe, and base name.
func labels(flai, recipe, base string) map[string]string {
	l := map[string]string{}
	for k, v := range map[string]string{LabelFlai: flai, LabelRecipe: recipe, LabelBase: base} {
		if v != "" {
			l[k] = v
		}
	}
	return l
}

// S-0735, ADR-0295 §2, §3: an image is missing when not on the host, the
// operator's own when the host names it, and otherwise current or stale by
// its flai, recipe, and base labels, each reason said.
func TestCheckJudgesEachImageByItsLabels(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "images", "agent"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, filepath.FromSlash(dockerfile)), []byte("FROM base\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	recipe := Recipe([]byte("FROM base\n"), "", nil)
	def := Image{Kind: KindHarness, Inner: harness.ClaudeCode, Ref: defaultImage, Source: SourceDefault}
	own := Image{Kind: KindHarness, Inner: harness.ClaudeCode, Ref: hostImage, Source: SourceHost}
	layer := Image{Kind: KindLayer, Inner: harness.ClaudeCode, Ref: layerImage, Source: SourceLayer, Dockerfile: dockerfile}
	type want struct {
		state   State
		reasons []ReasonKind
	}
	for _, c := range []struct {
		name    string
		flai    string
		images  []Image
		present map[string]Info
		want    []want
	}{
		{"missing", "1.2.3", []Image{def, layer}, nil, []want{{StateMissing, nil}, {StateMissing, nil}}},
		{"own image missing", "1.2.3", []Image{own}, nil, []want{{StateMissing, nil}}},
		{"own image never stale", "1.2.3", []Image{own, layer},
			map[string]Info{hostImage: {ID: "sha256:own"}, layerImage: {ID: "sha256:l", Labels: labels("1.2.3", recipe, "sha256:own")}},
			[]want{{StateOwn, nil}, {StateCurrent, nil}}},
		{"current release, v prefixes", "v1.2.3", []Image{def, layer},
			map[string]Info{defaultImage: {ID: "sha256:h", Labels: labels("1.2.3", "", "")}, layerImage: {ID: "sha256:l", Labels: labels("v1.2.3", recipe, "sha256:h")}},
			[]want{{StateCurrent, nil}, {StateCurrent, nil}}},
		{"another flai", "1.2.4", []Image{def},
			map[string]Info{defaultImage: {ID: "sha256:h", Labels: labels("1.2.3", "", "")}},
			[]want{{StateStale, []ReasonKind{ReasonFlai}}}},
		{"no flai label", "1.2.3", []Image{def},
			map[string]Info{defaultImage: {ID: "sha256:h"}},
			[]want{{StateStale, []ReasonKind{ReasonFlai}}}},
		{"development want, development label", "dev", []Image{def},
			map[string]Info{defaultImage: {ID: "sha256:h", Labels: labels("0.0.0-20261113-abc", "", "")}},
			[]want{{StateCurrent, nil}}},
		{"development want, release label", "", []Image{def},
			map[string]Info{defaultImage: {ID: "sha256:h", Labels: labels("1.2.3", "", "")}},
			[]want{{StateStale, []ReasonKind{ReasonFlai}}}},
		{"development want, no label", "dev", []Image{def},
			map[string]Info{defaultImage: {ID: "sha256:h"}},
			[]want{{StateStale, []ReasonKind{ReasonFlai}}}},
		{"another Dockerfile", "1.2.3", []Image{def, layer},
			map[string]Info{defaultImage: {ID: "sha256:h", Labels: labels("1.2.3", "", "")}, layerImage: {ID: "sha256:l", Labels: labels("1.2.3", "sha256:old", "sha256:h")}},
			[]want{{StateCurrent, nil}, {StateStale, []ReasonKind{ReasonRecipe}}}},
		{"another base", "1.2.3", []Image{def, layer},
			map[string]Info{defaultImage: {ID: "sha256:new", Labels: labels("1.2.3", "", "")}, layerImage: {ID: "sha256:l", Labels: labels("1.2.3", recipe, "sha256:old")}},
			[]want{{StateCurrent, nil}, {StateStale, []ReasonKind{ReasonBase}}}},
		{"harness image under the layer missing", "1.2.3", []Image{def, layer},
			map[string]Info{layerImage: {ID: "sha256:l", Labels: labels("1.2.3", recipe, "sha256:h")}},
			[]want{{StateMissing, nil}, {StateStale, []ReasonKind{ReasonBase}}}},
		{"no labels on a layer", "1.2.3", []Image{layer},
			map[string]Info{layerImage: {ID: "sha256:l"}},
			[]want{{StateStale, []ReasonKind{ReasonFlai, ReasonRecipe}}}},
	} {
		t.Run(c.name, func(t *testing.T) {
			f := &fakeImages{images: c.present}
			v := &counter{version: c.flai}
			got, err := Check(root, c.images, f.inspect, v.flai)
			if err != nil {
				t.Fatal(err)
			}
			if len(got) != len(c.want) || f.calls != len(c.images) || v.calls > 1 {
				t.Fatalf("Check = %+v after %d inspections and %d flai calls", got, f.calls, v.calls)
			}
			for i, w := range c.want {
				s := got[i]
				var kinds []ReasonKind
				for _, r := range s.Reasons {
					if r.Message == "" {
						t.Errorf("%s: a reason says nothing: %+v", s.Ref, s.Reasons)
					}
					kinds = append(kinds, r.Kind)
				}
				if s.Image != c.images[i] || s.State != w.state || !reflect.DeepEqual(kinds, w.reasons) {
					t.Errorf("status %d = %+v, want %s with %v", i, s, w.state, w.reasons)
				}
				if s.Kind == KindLayer && s.Present && s.Want.Recipe != recipe {
					t.Errorf("the layer's wanted recipe is %q, not %q", s.Want.Recipe, recipe)
				}
			}
		})
	}
}

// S-0735: the reasons name what the image holds and what it must hold.
func TestStaleReasonsSayWhatTheImageHolds(t *testing.T) {
	def := Image{Kind: KindHarness, Inner: harness.ClaudeCode, Ref: defaultImage, Source: SourceDefault}
	f := &fakeImages{images: map[string]Info{defaultImage: {ID: "sha256:h", Labels: labels("1.86.12", "", "")}}}
	got, err := Check(t.TempDir(), []Image{def}, f.inspect, (&counter{version: "dev"}).flai)
	if err != nil || len(got) != 1 || len(got[0].Reasons) != 1 || got[0].Reasons[0].Message != "built for flai 1.86.12, not dev" {
		t.Errorf("Check = %+v, %v", got, err)
	}
	if got[0].Want.Flai != "dev" || got[0].Labels.Flai != "1.86.12" {
		t.Errorf("want and labels = %+v, %+v", got[0].Want, got[0].Labels)
	}
}

// S-0735: the flai callback is not called when no image is present or every
// present one is the operator's own, and its failure fails the check.
func TestCheckAsksTheFlaiVersionOnlyForAnImageItJudges(t *testing.T) {
	def := Image{Kind: KindHarness, Inner: harness.ClaudeCode, Ref: defaultImage, Source: SourceDefault}
	own := Image{Kind: KindHarness, Inner: "other", Ref: hostImage, Source: SourceHost}
	failing := func() (string, error) { return "", errors.New("no flai") }
	for name, present := range map[string]map[string]Info{"none present": nil, "only the own": {hostImage: {ID: "sha256:own"}}} {
		t.Run(name, func(t *testing.T) {
			if _, err := Check(t.TempDir(), []Image{def, own}, (&fakeImages{images: present}).inspect, failing); err != nil {
				t.Errorf("the flai version was asked: %v", err)
			}
		})
	}
	present := &fakeImages{images: map[string]Info{defaultImage: {ID: "sha256:h"}}}
	if _, err := Check(t.TempDir(), []Image{def}, present.inspect, failing); err == nil || !strings.Contains(err.Error(), "no flai") {
		t.Errorf("a failing flai callback did not fail the check: %v", err)
	}
	unreadable := Image{Kind: KindLayer, Inner: harness.ClaudeCode, Ref: layerImage, Source: SourceLayer, Dockerfile: dockerfile}
	layer := &fakeImages{images: map[string]Info{layerImage: {ID: "sha256:l"}}}
	if _, err := Check(t.TempDir(), []Image{unreadable}, layer.inspect, (&counter{version: "1.0.0"}).flai); err == nil || !strings.Contains(err.Error(), "Dockerfile") {
		t.Errorf("an unreadable Dockerfile did not fail the check: %v", err)
	}
}

// fakeRunner is an execx.Runner that answers each program's run with out
// and err, and finds on PATH only the programs in path.
type fakeRunner struct {
	path []string
	out  string
	err  error
	runs [][]string
}

func (f *fakeRunner) Run(dir, name string, args ...string) (string, error) {
	f.runs = append(f.runs, append([]string{dir, name}, args...))
	return f.out, f.err
}

func (f *fakeRunner) RunInput(dir, name, _ string, args ...string) (string, error) {
	return f.Run(dir, name, args...)
}

func (f *fakeRunner) LookPath(name string) (string, error) {
	for _, p := range f.path {
		if p == name {
			return "/usr/bin/" + name, nil
		}
	}
	return "", errors.New("not found")
}

// S-0735: Docker reads an image's ID and labels from one docker image
// inspect, says a missing image is missing, and says when docker is not on
// PATH or could not say.
func TestDockerInspectsAnImage(t *testing.T) {
	inspected := `{"Id":"sha256:abc","Config":{"Labels":{"dev.system-flow.flai":"dev","dev.system-flow.commit":"1234"}}}`
	for _, c := range []struct {
		name        string
		r           *fakeRunner
		wantPresent bool
		wantErr     string
		noDocker    bool
	}{
		{"present", &fakeRunner{path: []string{"docker"}, out: inspected}, true, "", false},
		{"missing", &fakeRunner{path: []string{"docker"}, out: "Error response from daemon: No such image: x:1", err: errors.New("exit status 1")}, false, "", false},
		{"missing, older client", &fakeRunner{path: []string{"docker"}, err: errors.New("Error: No Such Image: x:1")}, false, "", false},
		{"not on PATH", &fakeRunner{}, false, "not on PATH", true},
		{"docker failing", &fakeRunner{path: []string{"docker"}, out: "Cannot connect to the Docker daemon\nmore", err: errors.New("exit status 1")}, false, "Cannot connect to the Docker daemon", false},
		{"no image printed", &fakeRunner{path: []string{"docker"}, out: "nonsense"}, false, "printed no image", false},
	} {
		t.Run(c.name, func(t *testing.T) {
			info, present, err := Docker(c.r, "docker")("x:1")
			if present != c.wantPresent || errors.Is(err, ErrNoDocker) != c.noDocker {
				t.Fatalf("Docker = %+v, %v, %v", info, present, err)
			}
			if c.wantErr == "" && err != nil || c.wantErr != "" && (err == nil || !strings.Contains(err.Error(), c.wantErr)) {
				t.Fatalf("Docker error = %v, want one saying %q", err, c.wantErr)
			}
			if c.noDocker && len(c.r.runs) != 0 {
				t.Errorf("docker ran though not on PATH: %v", c.r.runs)
			}
			if c.wantPresent {
				if info.ID != "sha256:abc" || info.Labels[LabelFlai] != "dev" || info.Labels[LabelCommit] != "1234" {
					t.Errorf("Info = %+v", info)
				}
				want := []string{"", "docker", "image", "inspect", "--format", "{{json .}}", "x:1"}
				if len(c.r.runs) != 1 || !reflect.DeepEqual(c.r.runs[0], want) {
					t.Errorf("runs = %v, want one %v", c.r.runs, want)
				}
			}
		})
	}
}

// S-0735, ADR-0295 §3: the project's flai is the version its flai.command
// prints with version --json, run in the root, else the running one.
func TestProjectFlaiIsTheCommandsVersionElseTheRunningOne(t *testing.T) {
	root := t.TempDir()
	r := &fakeRunner{}
	if got, err := ProjectFlai(root, manifest.Manifest{}, r, "1.2.3"); err != nil || got != "1.2.3" || len(r.runs) != 0 {
		t.Errorf("ProjectFlai with no command = %q, %v after %v", got, err, r.runs)
	}
	m := manifest.Manifest{Flai: manifest.Requirement{Command: "scripts/flai.sh"}}
	r = &fakeRunner{out: "building flai\n{\n  \"version\": \"dev\",\n  \"commit\": \"abc\"\n}"}
	got, err := ProjectFlai(root, m, r, "1.2.3")
	want := []string{root, filepath.Join(root, "scripts", "flai.sh"), "version", "--json"}
	if err != nil || got != "dev" || len(r.runs) != 1 || !reflect.DeepEqual(r.runs[0], want) {
		t.Errorf("ProjectFlai = %q, %v after %v; want dev after %v", got, err, r.runs, want)
	}
	for name, r := range map[string]*fakeRunner{
		"failing":    {out: "build failed", err: errors.New("exit status 2")},
		"no version": {out: "{\"commit\": \"abc\"}"},
	} {
		t.Run(name, func(t *testing.T) {
			if _, err := ProjectFlai(root, m, r, "1.2.3"); err == nil || !strings.Contains(err.Error(), "scripts/flai.sh version --json") {
				t.Errorf("ProjectFlai = %v, want an error naming the command", err)
			}
		})
	}
}
