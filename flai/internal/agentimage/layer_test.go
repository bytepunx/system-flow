package agentimage

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/harness"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
)

const (
	layerHarness = "flai-agent-claude-code:local"
	harnessID    = "sha256:harness"
)

// layerRoot writes a project under a new folder whose layer Dockerfile,
// at dockerfile, holds data, and is its root.
func layerRoot(t *testing.T, data string) string {
	t.Helper()
	root := t.TempDir()
	path := filepath.Join(root, filepath.FromSlash(dockerfile))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(data), 0o644); err != nil {
		t.Fatal(err)
	}
	return root
}

// layerImg is the layer image of manifest key t with the build arguments args.
func layerImg(args map[string]string) Image {
	return Image{Kind: KindLayer, Inner: harness.ClaudeCode, Ref: layerImage, Source: SourceLayer, Dockerfile: dockerfile, Args: args}
}

// layerBuilder is a Builder whose docker is exec and on whose host the
// harness image layerHarness is harnessID, writing to out.
func layerBuilder(exec Exec, out *syncBuffer, sleep func(context.Context, time.Duration) error) Builder {
	inspect := &fakeImages{images: map[string]Info{layerHarness: {ID: harnessID}}}
	r := Runner{Exec: exec, Sleep: sleep}
	if out != nil {
		r.Out = out
	}
	return Builder{Run: r, Inspect: inspect.inspect}
}

// S-0745, ADR-0295 §2, §4: the layer is built in the project's root, from
// its Dockerfile with the root as the context, BASE set to the harness
// image, each argument by name, and labelled with its recipe and the
// harness image's ID, and the build says so before it runs.
func TestLayerBuildsOverTheHarnessImage(t *testing.T) {
	const file = "ARG BASE\nFROM ${BASE}\nARG A\nARG Z\n"
	root := layerRoot(t, file)
	d := &buildDocker{}
	out := &syncBuffer{}
	args := map[string]string{"Z": "1", "A": "2"}
	got, err := layerBuilder(d.exec, out, nil).Layer(context.Background(), root, layerImg(args), layerHarness)
	if err != nil {
		t.Fatalf("Layer: %v", err)
	}
	recipe := Recipe([]byte(file), "", args)
	want := Layered{Image: layerImg(args), Tags: []string{layerImage}, Base: layerHarness, BaseID: harnessID, Recipe: recipe}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Layer =\n%+v\nwant\n%+v", got, want)
	}
	wantArgs := []string{"build", "-f", filepath.FromSlash(dockerfile),
		"--build-arg", "BASE=" + layerHarness, "--build-arg", "A=2", "--build-arg", "Z=1",
		"--label", LabelRecipe + "=" + recipe, "--label", LabelBase + "=" + harnessID,
		"-t", layerImage, "."}
	if len(d.builds) != 1 || !reflect.DeepEqual(d.builds[0].args, wantArgs) || d.builds[0].dir != root {
		t.Fatalf("docker ran %+v, want one build in %s with\n%q", d.builds, root, wantArgs)
	}
	line := "agent image: building " + layerImage + " over " + layerHarness + " (" + harnessID + ") from " + dockerfile + " with the build arguments A=2, Z=1\n"
	if out.String() != line {
		t.Errorf("Out =\n%q\nwant\n%q", out.String(), line)
	}
}

// S-0745, ADR-0295 §4: a layer with no arguments is built with BASE alone,
// over its Dockerfile's own ARG defaults, and its recipe is that of no arguments.
func TestLayerWithNoArgumentsPassesBaseAlone(t *testing.T) {
	const file = "ARG BASE\nFROM ${BASE}\nARG GO_VERSION=1.26.0\n"
	root := layerRoot(t, file)
	d := &buildDocker{}
	out := &syncBuffer{}
	got, err := layerBuilder(d.exec, out, nil).Layer(context.Background(), root, layerImg(nil), layerHarness)
	if err != nil {
		t.Fatalf("Layer: %v", err)
	}
	if got.Recipe != Recipe([]byte(file), "", nil) {
		t.Errorf("the recipe is %s, not that of no arguments", got.Recipe)
	}
	wantArgs := []string{"build", "-f", filepath.FromSlash(dockerfile), "--build-arg", "BASE=" + layerHarness,
		"--label", LabelRecipe + "=" + got.Recipe, "--label", LabelBase + "=" + harnessID, "-t", layerImage, "."}
	if len(d.builds) != 1 || !reflect.DeepEqual(d.builds[0].args, wantArgs) {
		t.Errorf("docker ran %+v, want one build with\n%q", d.builds, wantArgs)
	}
	if !strings.Contains(out.String(), "with its Dockerfile's own ARG defaults") {
		t.Errorf("Out does not say the Dockerfile's defaults hold:\n%s", out.String())
	}
}

// pin is the value scripts/install-tools.sh pins for name, as
// scripts/agent-image.sh's pinned reads it.
func pin(t *testing.T, tools []byte, name string) string {
	t.Helper()
	m := regexp.MustCompile(`(?m)^` + name + `="\$\{` + name + `:-(.*)\}"`).FindSubmatch(tools)
	if m == nil {
		t.Fatalf("scripts/install-tools.sh pins no %s", name)
	}
	return string(m[1])
}

// sortedPairs is a docker build's arguments with each flag and its value
// as one string, sorted, and its last argument, the context.
func sortedPairs(t *testing.T, args []string) ([]string, string) {
	t.Helper()
	if len(args) < 2 || args[0] != "build" || len(args)%2 != 0 {
		t.Fatalf("not a docker build with flags in pairs and a context: %q", args)
	}
	var pairs []string
	for i := 1; i < len(args)-1; i += 2 {
		pairs = append(pairs, args[i]+" "+args[i+1])
	}
	slices.Sort(pairs)
	return pairs, args[len(args)-1]
}

// S-0745: for this repository's layer and the arguments scripts/agent-image.sh
// passes it, from flai/go.mod and scripts/install-tools.sh's pins, flai's
// build passes the script's docker arguments, plus the recipe and base
// labels, in the repository's root.
func TestLayerGivesTheArgumentsScriptsAgentImageGives(t *testing.T) {
	root, err := filepath.Abs(filepath.Join("..", "..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	gomod, err := os.ReadFile(filepath.Join(root, "flai", "go.mod"))
	if err != nil {
		t.Fatal(err)
	}
	goVersion := regexp.MustCompile(`(?m)^go (\S+)$`).FindSubmatch(gomod)
	tools, err := os.ReadFile(filepath.Join(root, "scripts", "install-tools.sh"))
	if err != nil {
		t.Fatal(err)
	}
	if goVersion == nil {
		t.Fatal("flai/go.mod has no go line")
	}
	args := map[string]string{
		"GO_VERSION":            string(goVersion[1]),
		"GOLANGCI_LINT_VERSION": pin(t, tools, "GOLANGCI_LINT_VERSION"),
		"SHELLCHECK_VERSION":    pin(t, tools, "SHELLCHECK_VERSION"),
		"PNPM_VERSION":          pin(t, tools, "PNPM_VERSION"),
		"GORELEASER_VERSION":    pin(t, tools, "GORELEASER_VERSION"),
	}
	const layer = "flai/images/system-flow/Dockerfile"
	m := manifest.Manifest{Key: "sf", Agent: &manifest.Agent{Image: &manifest.AgentImage{Layer: layer, Args: args}}}
	img := For(m, nil, harness.ClaudeCode)
	d := &buildDocker{}
	got, err := layerBuilder(d.exec, nil, nil).Layer(context.Background(), root, img, layerHarness)
	if err != nil {
		t.Fatalf("Layer: %v", err)
	}
	if len(d.builds) != 1 || d.builds[0].dir != root {
		t.Fatalf("docker ran %+v, want one build in %s", d.builds, root)
	}
	// build_layer's docker build, as scripts/agent-image.sh runs it from the root.
	script := []string{"build", "-f", layer, "-t", "flai-agent-sf-claude-code:local", "--build-arg", "BASE=" + layerHarness,
		"--build-arg", "GO_VERSION=" + args["GO_VERSION"], "--build-arg", "GOLANGCI_LINT_VERSION=" + args["GOLANGCI_LINT_VERSION"],
		"--build-arg", "SHELLCHECK_VERSION=" + args["SHELLCHECK_VERSION"], "--build-arg", "PNPM_VERSION=" + args["PNPM_VERSION"],
		"--build-arg", "GORELEASER_VERSION=" + args["GORELEASER_VERSION"], "."}
	want := slices.Concat(script[:len(script)-1], []string{"--label", LabelRecipe + "=" + got.Recipe, "--label", LabelBase + "=" + harnessID, "."})
	gotPairs, gotContext := sortedPairs(t, d.builds[0].args)
	wantPairs, wantContext := sortedPairs(t, want)
	if !reflect.DeepEqual(gotPairs, wantPairs) || gotContext != wantContext {
		t.Errorf("flai built with\n%q\nwant the script's arguments and the labels\n%q", d.builds[0].args, want)
	}
	data, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(layer)))
	if err != nil {
		t.Fatal(err)
	}
	if got.Recipe != Recipe(data, "", args) {
		t.Errorf("the recipe %s is not that of %s and the arguments", got.Recipe, layer)
	}
}

// S-0745: a layer whose harness image is not on this host is refused,
// naming it and saying to ensure it first, and nothing is built; so is one
// whose Dockerfile cannot be read, naming it and the manifest's setting.
func TestLayerRefusesWhatItCannotBuildOver(t *testing.T) {
	root := layerRoot(t, "FROM x\n")
	for _, c := range []struct {
		name string
		root string
		img  Image
		base string
		want []string
	}{
		{"harness image missing", root, layerImg(nil), "flai-agent-claude-code:other", []string{"flai-agent-claude-code:other", "not on this host", "flai image ensure"}},
		{"Dockerfile missing", t.TempDir(), layerImg(nil), layerHarness, []string{dockerfile, "agent.image.layer", manifest.File}},
		{"no harness image named", root, layerImg(nil), "", []string{layerImage, "no harness image"}},
		{"a harness image", root, Image{Kind: KindHarness, Inner: harness.ClaudeCode, Ref: defaultImage, Source: SourceDefault}, layerHarness, []string{defaultImage, "not the project's layer"}},
	} {
		t.Run(c.name, func(t *testing.T) {
			d := &buildDocker{}
			_, err := layerBuilder(d.exec, nil, nil).Layer(context.Background(), c.root, c.img, c.base)
			if err == nil {
				t.Fatal("Layer built")
			}
			for _, w := range c.want {
				if !strings.Contains(err.Error(), filepath.FromSlash(w)) {
					t.Errorf("Layer = %v, want an error saying %q", err, w)
				}
			}
			if len(d.builds) != 0 {
				t.Errorf("docker ran %d builds", len(d.builds))
			}
		})
	}
	failing := Builder{Run: Runner{Exec: (&buildDocker{}).exec}, Inspect: func(string) (Info, bool, error) { return Info{}, false, errors.New("docker is down") }}
	if _, err := failing.Layer(context.Background(), root, layerImg(nil), layerHarness); err == nil || !strings.Contains(err.Error(), "docker is down") || !strings.Contains(err.Error(), layerHarness) {
		t.Errorf("Layer with a failing inspect = %v, want its error naming %s", err, layerHarness)
	}
}

// S-0745, ADR-0295 §5: the layer's build runs through the Runner, so a
// build that fails on a connection reset is tried again, and passes.
func TestLayerTriesANetworkResetAgain(t *testing.T) {
	root := layerRoot(t, "FROM x\n")
	docker := &fakeDocker{t: t, tries: []try{{out: networkOutputs["connection reset"], err: errExit}, {out: "done\n"}}}
	waits := &sleeps{}
	if _, err := layerBuilder(docker.exec, nil, waits.sleep).Layer(context.Background(), root, layerImg(nil), layerHarness); err != nil {
		t.Fatalf("Layer = %v, want a pass on the second try", err)
	}
	if want := []time.Duration{30 * time.Second}; docker.ran() != 2 || !reflect.DeepEqual(waits.waits, want) || docker.dir != root {
		t.Errorf("docker ran %d times in %s, waiting %v, want twice in %s waiting %v", docker.ran(), docker.dir, waits.waits, root, want)
	}
}
