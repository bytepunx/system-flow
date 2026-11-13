package agentimage

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/images"
)

const commandLayerData = "ARG BASE\nFROM ${BASE}\nARG GO\n"

// commandRecipes are the recipe labels flai's own release build gives the
// bedrock and claude-code harness images, from the Dockerfiles flai carries.
func commandRecipes(t *testing.T) (string, string) {
	t.Helper()
	bedrock, err := images.Dockerfile(images.Bedrock)
	if err != nil {
		t.Fatal(err)
	}
	inner, err := images.Dockerfile("claude-code")
	if err != nil {
		t.Fatal(err)
	}
	return Recipe(bedrock, targetPrebuilt, nil), Recipe(inner, "", nil)
}

// S-0759, ADR-0295 §7: the build command's variables name the version
// without a leading v, the inner harness, the tags, the recipes flai's own
// build would label, and the images to make, in the ADR's order.
func TestCommandImagesEnvNamesWhatTheImagesMustHave(t *testing.T) {
	t.Parallel()
	root := layerRoot(t, commandLayerData)
	layer := layerStatus(StateCurrent)
	c, err := commandImages(root, ModeEnsure, harnessStatus(StateMissing), &layer, BuildSource{Version: "v" + pullV})
	if err != nil {
		t.Fatal(err)
	}
	bedrock, inner := commandRecipes(t)
	want := []string{
		"FLAI_IMAGE_VERSION=" + pullV,
		"FLAI_IMAGE_INNER=claude-code",
		"FLAI_IMAGE_BEDROCK=flai-agent-bedrock:" + pullV,
		"FLAI_IMAGE_HARNESS=" + ensureRef,
		"FLAI_IMAGE_LAYER=" + ensureLayerRef,
		"FLAI_IMAGE_RECIPE_BEDROCK=" + bedrock,
		"FLAI_IMAGE_RECIPE_HARNESS=" + inner,
		"FLAI_IMAGE_RECIPE_LAYER=" + Recipe([]byte(commandLayerData), "", map[string]string{"GO": "1.25"}),
		"FLAI_IMAGE_STALE=bedrock harness layer",
	}
	if got := c.Env(); !reflect.DeepEqual(got, want) {
		t.Errorf("Env =\n%s\nwant\n%s", strings.Join(got, "\n"), strings.Join(want, "\n"))
	}
}

// S-0759: FLAI_IMAGE_STALE names, in ensure mode, a missing or stale
// harness image with bedrock under it and a layer missing, stale, or over a
// harness image to be made; in build mode every image; in layer mode the
// layer alone; and never the operator's own image, whose tag the layer
// builds over, with no bedrock and no harness recipe.
func TestCommandImagesStaleNamesWhatTheRunAsksMade(t *testing.T) {
	t.Parallel()
	root := layerRoot(t, commandLayerData)
	own := Status{Image: Image{Kind: KindHarness, Inner: "claude-code", Ref: "mine:1", Source: SourceHost}, Present: true, State: StateOwn}
	for _, tc := range []struct {
		name    string
		mode    Mode
		harness Status
		layer   State // "" for no layer
		want    []string
	}{
		{"all current", ModeEnsure, harnessStatus(StateCurrent), StateCurrent, nil},
		{"stale layer", ModeEnsure, harnessStatus(StateCurrent), StateStale, []string{PartLayer}},
		{"missing layer", ModeEnsure, harnessStatus(StateCurrent), StateMissing, []string{PartLayer}},
		{"stale harness, current layer", ModeEnsure, harnessStatus(StateStale), StateCurrent, []string{PartBedrock, PartHarness, PartLayer}},
		{"missing harness, no layer", ModeEnsure, harnessStatus(StateMissing), "", []string{PartBedrock, PartHarness}},
		{"build, all current", ModeBuild, harnessStatus(StateCurrent), StateCurrent, []string{PartBedrock, PartHarness, PartLayer}},
		{"layer mode, stale harness", ModeLayer, harnessStatus(StateStale), StateCurrent, []string{PartLayer}},
		{"layer mode, no layer", ModeLayer, harnessStatus(StateCurrent), "", nil},
		{"own harness, missing layer", ModeEnsure, own, StateMissing, []string{PartLayer}},
		{"own harness, build", ModeBuild, own, StateCurrent, []string{PartLayer}},
	} {
		var l *Status
		if tc.layer != "" {
			s := layerStatus(tc.layer)
			l = &s
		}
		c, err := commandImages(root, tc.mode, tc.harness, l, BuildSource{Version: pullV})
		if err != nil {
			t.Fatalf("%s: %v", tc.name, err)
		}
		if !reflect.DeepEqual(c.Stale, tc.want) {
			t.Errorf("%s: Stale = %v, want %v", tc.name, c.Stale, tc.want)
		}
		if tc.harness.Source == SourceHost && (c.Harness != "mine:1" || c.Bedrock != "" || c.RecipeBedrock != "" || c.RecipeHarness != "" || c.RecipeLayer == "") {
			t.Errorf("%s: images %+v", tc.name, c)
		}
		if tc.layer == "" && (c.Layer != "" || c.RecipeLayer != "") {
			t.Errorf("%s: a layer is named with none in the manifest: %+v", tc.name, c)
		}
	}
}

// S-0759, ADR-0295 §3: a development tree's recipes are its own Dockerfiles'
// at bedrock's tree target; a development flai with no tree gives none.
func TestCommandImagesRecipesOfADevelopmentFlai(t *testing.T) {
	t.Parallel()
	tree := writeTree(t, "FROM scratch AS tree\n", "ARG BASE\nFROM ${BASE}\n")
	c, err := commandImages(tree, ModeEnsure, harnessStatus(StateMissing), nil, BuildSource{Version: "dev", Tree: tree})
	if err != nil {
		t.Fatal(err)
	}
	if c.Version != "dev" || c.Bedrock != "flai-agent-bedrock:dev" || c.RecipeBedrock != Recipe([]byte("FROM scratch AS tree\n"), targetTree, nil) || c.RecipeHarness != Recipe([]byte("ARG BASE\nFROM ${BASE}\n"), "", nil) {
		t.Errorf("tree: %+v", c)
	}
	c, err = commandImages(t.TempDir(), ModeEnsure, harnessStatus(StateMissing), nil, BuildSource{Version: "dev"})
	if err != nil {
		t.Fatal(err)
	}
	if c.RecipeBedrock != "" || c.RecipeHarness != "" || c.Bedrock != "flai-agent-bedrock:dev" {
		t.Errorf("no tree: %+v", c)
	}
}

// S-0759, ADR-0295 §7: the command runs in the root with no shell, argv[0]
// as the program, its environment the inherited one and the FLAI_IMAGE_*
// variables, its output streamed to the Runner's Out after a line naming it.
func TestRunCommandRunsTheProgramInTheRootWithTheVariables(t *testing.T) {
	if _, err := exec.LookPath("sh"); err != nil {
		t.Skip("sh is not on PATH")
	}
	t.Setenv("FLAI_COMMAND_TEST_INHERITED", "kept")
	root := t.TempDir()
	envFile := filepath.Join(t.TempDir(), "env")
	c := CommandImages{Version: pullV, Inner: "claude-code", Bedrock: BedrockTag(pullV), Harness: ensureRef, RecipeBedrock: "sha256:b", RecipeHarness: "sha256:h", Stale: []string{PartBedrock, PartHarness}}
	out := &syncBuffer{}
	argv := []string{"sh", "-c", `env > "$1"; pwd >> "$1"; echo "made $FLAI_IMAGE_HARNESS"`, "make images", envFile}
	if err := RunCommand(context.Background(), Runner{Out: out}, root, argv, c); err != nil {
		t.Fatalf("RunCommand: %v\n%s", err, out.String())
	}
	data, err := os.ReadFile(envFile)
	if err != nil {
		t.Fatal(err)
	}
	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	resolved, _ := filepath.EvalSymlinks(root)
	if pwd := lines[len(lines)-1]; pwd != root && pwd != resolved {
		t.Errorf("ran in %s, not the root %s", pwd, root)
	}
	for _, want := range append(c.Env(), "FLAI_COMMAND_TEST_INHERITED=kept") {
		if !strings.Contains("\n"+string(data), "\n"+want+"\n") {
			t.Errorf("the command's environment lacks %q:\n%s", want, data)
		}
	}
	text := out.String()
	for _, want := range []string{
		`agent image: running the operator's build command sh -c "env > \"$1\"; pwd >> \"$1\"; echo \"made $FLAI_IMAGE_HARNESS\"" "make images" ` + envFile + " in " + root + " for the images of claude-code",
		"FLAI_IMAGE_STALE=\"bedrock harness\"",
		"FLAI_IMAGE_LAYER=\"\"",
		"made " + ensureRef + "\n",
	} {
		if !strings.Contains(text, want) {
			t.Errorf("output lacks %q:\n%s", want, text)
		}
	}
}

// S-0759: a command that exits non-zero fails once, with no try again even
// when its output names a network failure, naming the command and what it
// said; the fake records it ran argv[0] with argv[1:], in the root, with
// the variables.
func TestRunCommandFailsOnANonZeroExitWithNoTryAgain(t *testing.T) {
	t.Parallel()
	docker := &fakeDocker{t: t, tries: []try{{out: "pulling\nread: connection reset by peer\n", err: errExit}}}
	waits := &sleeps{}
	c := CommandImages{Version: pullV, Inner: "claude-code", Stale: []string{PartLayer}}
	err := RunCommand(context.Background(), Runner{Exec: docker.exec, Sleep: waits.sleep}, "/project", []string{"/opt/make-images", "--fast"}, c)
	var failed *CommandError
	if !errors.As(err, &failed) || !errors.Is(err, errExit) {
		t.Fatalf("RunCommand = %v, want a *CommandError over the exit", err)
	}
	if failed.Said != "read: connection reset by peer" || failed.Inner != "claude-code" {
		t.Errorf("CommandError = %+v", failed)
	}
	if msg := err.Error(); !strings.Contains(msg, "the operator's build command /opt/make-images --fast failed for the agent images of claude-code (read: connection reset by peer; exit status 1)") || !strings.Contains(msg, "flai serve agent image command --clear") {
		t.Errorf("message %q", msg)
	}
	if docker.ran() != 1 || len(waits.waits) != 0 {
		t.Errorf("ran %d times and waited %v, want once and no wait", docker.ran(), waits.waits)
	}
	if docker.program != "/opt/make-images" || docker.dir != "/project" || !reflect.DeepEqual(docker.args, []string{"--fast"}) || !reflect.DeepEqual(docker.env, c.Env()) {
		t.Errorf("ran %s %v in %s with %v", docker.program, docker.args, docker.dir, docker.env)
	}
}

// S-0759, ADR-0295 §5: a command silent for the stall limit is stopped with
// a *StallError naming it and the limit.
func TestRunCommandStopsASilentCommand(t *testing.T) {
	t.Parallel()
	docker := &fakeDocker{t: t, tries: []try{{out: "starting\n", block: true}}}
	err := RunCommand(context.Background(), Runner{Exec: docker.exec, Stall: 50 * time.Millisecond}, "/project", []string{"make-images"}, CommandImages{Inner: "claude-code"})
	var stall *StallError
	if !errors.As(err, &stall) || stall.Step != "the build command make-images" || stall.Limit != 50*time.Millisecond {
		t.Fatalf("RunCommand = %v, want a *StallError naming the command and 50ms", err)
	}
	if !strings.Contains(err.Error(), "the operator's build command") || docker.ran() != 1 {
		t.Errorf("message %q, ran %d times", err, docker.ran())
	}
}

// S-0759: a command that names no program is refused before anything runs.
func TestRunCommandRefusesNoProgram(t *testing.T) {
	t.Parallel()
	for _, argv := range [][]string{nil, {" ", "x"}} {
		docker := &fakeDocker{t: t}
		if err := RunCommand(context.Background(), Runner{Exec: docker.exec}, "/project", argv, CommandImages{Inner: "claude-code"}); err == nil || !strings.Contains(err.Error(), "names no program") || docker.ran() != 0 {
			t.Errorf("%q: %v, ran %d", argv, err, docker.ran())
		}
	}
}

func TestCommandWords(t *testing.T) {
	t.Parallel()
	got := CommandWords([]string{"scripts/agent-image.sh", "--layer", "", "a b", `say "hi"`, "it's", "tab\there"})
	want := `scripts/agent-image.sh --layer "" "a b" "say \"hi\"" "it's" "tab\there"`
	if got != want {
		t.Errorf("CommandWords = %s, want %s", got, want)
	}
}
