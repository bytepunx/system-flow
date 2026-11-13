package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/agentimage"
	"github.com/bytepunx/system-flow/flai/internal/buildinfo"
	"github.com/bytepunx/system-flow/flai/internal/harness"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// imageDocker answers docker image inspect for the images it holds, by
// reference, and "No such image" for any other; with noDocker, docker is
// not on PATH. It answers nothing else.
type imageDocker struct {
	images   map[string]agentimage.Info
	noDocker bool
	asked    []string
}

func (d *imageDocker) Run(_, name string, args ...string) (string, error) {
	if name != "docker" || len(args) != 5 || args[0] != "image" || args[1] != "inspect" {
		return "", fmt.Errorf("imageDocker does not answer %s %s", name, strings.Join(args, " "))
	}
	ref := args[4]
	d.asked = append(d.asked, ref)
	info, ok := d.images[ref]
	if !ok {
		return "Error: No such image: " + ref, errors.New("exit status 1")
	}
	data, err := json.Marshal(map[string]any{"Id": info.ID, "Config": map[string]any{"Labels": info.Labels}})
	return string(data), err
}

func (d *imageDocker) RunInput(dir, name, _ string, args ...string) (string, error) {
	return d.Run(dir, name, args...)
}

func (d *imageDocker) LookPath(name string) (string, error) {
	if name == "docker" && d.noDocker {
		return "", errors.New("not found")
	}
	return "/usr/bin/" + name, nil
}

const (
	imageLayerFile = "images/agent/Dockerfile"
	imageLayerText = "FROM flai-agent-claude-code:local\nRUN true\n"
	imageLayerRef  = "flai-agent-t-claude-code:local"
	imageBaseID    = "sha256:base"
)

var imageHarnessRef = harness.DefaultImages[harness.ClaudeCode]

// imageProject is a project whose manifest adds agentYAML, an agent: block,
// with its layer's Dockerfile written, and the path of a config file of its
// own for the command to read.
func imageProject(t *testing.T, agentYAML string) (root, cfgPath string) {
	t.Helper()
	root = tempProject(t)
	m := "version: 1\nname: t\nkey: t\nlayout:\n  design: design\n  docs: docs\n  wip: wip\n" + agentYAML
	if err := os.WriteFile(filepath.Join(root, "system-flow.yaml"), []byte(m), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(root, filepath.Dir(imageLayerFile)), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, imageLayerFile), []byte(imageLayerText), 0o644); err != nil {
		t.Fatal(err)
	}
	return root, filepath.Join(t.TempDir(), "config.json")
}

const imageLayerYAML = "agent:\n  image:\n    layer: " + imageLayerFile + "\n"

// currentImages are the claude-code harness image and the project's layer
// over it, both labelled as the running flai must find them.
func currentImages() map[string]agentimage.Info {
	return map[string]agentimage.Info{
		imageHarnessRef: {ID: imageBaseID, Labels: map[string]string{agentimage.LabelFlai: buildinfo.Get().Version}},
		imageLayerRef: {ID: "sha256:layer", Labels: map[string]string{
			agentimage.LabelFlai:   buildinfo.Get().Version,
			agentimage.LabelRecipe: agentimage.Recipe([]byte(imageLayerText), "", nil),
			agentimage.LabelBase:   imageBaseID,
			agentimage.LabelCommit: "abc123",
		}},
	}
}

// runImageStatus runs flai image status with args in root, with the config
// file cfgPath and docker d, and returns its output, its exit code as run
// gives it, and its error.
func runImageStatus(t *testing.T, root, cfgPath string, d *imageDocker, args ...string) (string, int, error) {
	t.Helper()
	var out, errOut bytes.Buffer
	a := &app{configPath: cfgPath, cwd: root, runner: d, out: &out, errOut: &errOut}
	c := newRootCmdWith(a)
	c.SetArgs(append([]string{"image", "status"}, args...))
	err := c.Execute()
	code := 0
	var ee *exitError
	switch {
	case errors.As(err, &ee):
		code = ee.code
	case err != nil:
		code = 1
	}
	return out.String(), code, err
}

// decodeImageStatus is flai image status --json's output, which must be one
// JSON object.
func decodeImageStatus(t *testing.T, out string) imageStatusReport {
	t.Helper()
	var r imageStatusReport
	dec := json.NewDecoder(strings.NewReader(out))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&r); err != nil {
		t.Fatalf("flai image status --json printed no report: %v\n%s", err, out)
	}
	if dec.More() {
		t.Fatalf("flai image status --json printed more than one object:\n%s", out)
	}
	return r
}

// S-0735: with a layer and both its images current, flai image status
// exits 0 and lists the harness image, then the layer naming its
// Dockerfile; --json prints the flai version they hold, each image's
// status, and ok.
func TestImageStatusAllCurrent(t *testing.T) {
	t.Parallel()
	root, cfg := imageProject(t, imageLayerYAML)
	d := &imageDocker{images: currentImages()}
	out, code, err := runImageStatus(t, root, cfg, d, "--json")
	if code != 0 || err != nil {
		t.Fatalf("exit %d, %v\n%s", code, err, out)
	}
	r := decodeImageStatus(t, out)
	if r.Flai != buildinfo.Get().Version || !r.OK || len(r.Images) != 2 {
		t.Fatalf("report = %+v", r)
	}
	h, l := r.Images[0], r.Images[1]
	if h.Kind != agentimage.KindHarness || h.Ref != imageHarnessRef || h.Source != agentimage.SourceDefault || h.State != agentimage.StateCurrent || !h.Present || h.ID != imageBaseID {
		t.Errorf("harness image = %+v", h)
	}
	if l.Kind != agentimage.KindLayer || l.Ref != imageLayerRef || l.Dockerfile != imageLayerFile || l.State != agentimage.StateCurrent || l.Want.Base != imageBaseID || l.Want.Recipe != agentimage.Recipe([]byte(imageLayerText), "", nil) || l.Labels.Commit != "abc123" {
		t.Errorf("layer image = %+v", l)
	}

	out, code, err = runImageStatus(t, root, cfg, d)
	if code != 0 || err != nil {
		t.Fatalf("text: exit %d, %v\n%s", code, err, out)
	}
	want := "claude-code harness " + imageHarnessRef + " (default): current\n" +
		"claude-code layer " + imageLayerRef + " (layer from " + imageLayerFile + "): current (built from commit abc123)\n"
	if out != want {
		t.Errorf("text output =\n%s\nwant\n%s", out, want)
	}
}

// S-0735: an image built for another flai is stale: the line says why, the
// report is not ok, and the command exits non-zero after printing it, with
// --json still one JSON object on stdout.
func TestImageStatusStale(t *testing.T) {
	t.Parallel()
	root, cfg := imageProject(t, imageLayerYAML)
	images := currentImages()
	images[imageHarnessRef] = agentimage.Info{ID: imageBaseID, Labels: map[string]string{agentimage.LabelFlai: "1.86.12"}}
	d := &imageDocker{images: images}
	out, code, err := runImageStatus(t, root, cfg, d)
	if code == 0 || err == nil || !strings.Contains(err.Error(), "1 agent image is missing or stale") || !strings.Contains(err.Error(), "make them with flai image ensure") {
		t.Fatalf("exit %d, %v", code, err)
	}
	if want := "claude-code harness " + imageHarnessRef + " (default): stale: built for flai 1.86.12, not " + buildinfo.Get().Version + "\n"; !strings.HasPrefix(out, want) {
		t.Errorf("text output =\n%s\nwant it to start\n%s", out, want)
	}

	out, code, _ = runImageStatus(t, root, cfg, d, "--json")
	r := decodeImageStatus(t, out)
	if code == 0 || r.OK || r.Images[0].State != agentimage.StateStale || len(r.Images[0].Reasons) != 1 || r.Images[0].Reasons[0].Kind != agentimage.ReasonFlai {
		t.Errorf("exit %d, report = %+v", code, r)
	}
}

// S-0735: a layer not on this host is missing, and the command exits
// non-zero after the list.
func TestImageStatusMissing(t *testing.T) {
	t.Parallel()
	root, cfg := imageProject(t, imageLayerYAML)
	images := currentImages()
	delete(images, imageLayerRef)
	out, code, err := runImageStatus(t, root, cfg, &imageDocker{images: images}, "--json")
	if code == 0 || err == nil || !strings.Contains(err.Error(), "1 agent image is missing or stale") {
		t.Fatalf("exit %d, %v", code, err)
	}
	r := decodeImageStatus(t, out)
	if r.OK || len(r.Images) != 2 || r.Images[0].State != agentimage.StateCurrent || r.Images[1].State != agentimage.StateMissing || r.Images[1].Present {
		t.Errorf("report = %+v", r)
	}
	out, _, _ = runImageStatus(t, root, cfg, &imageDocker{images: images})
	if want := "claude-code layer " + imageLayerRef + " (layer from " + imageLayerFile + "): missing\n"; !strings.HasSuffix(out, want) {
		t.Errorf("text output =\n%s\nwant it to end\n%s", out, want)
	}
}

// S-0735: the host's --image for the project, read as flai serve reads it
// (the project's own container entry over the host-wide one, ADR-0261), is
// the operator's own and not judged; an open story's agent on the container
// harness is what needs it.
func TestImageStatusOwnImage(t *testing.T) {
	t.Parallel()
	root, cfg := imageProject(t, "")
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	projects, _ := json.Marshal(map[string]any{mainRootOf(repo): map[string]any{"harnesses": map[string]any{"container": map[string]any{"images": map[string]string{harness.ClaudeCode: "mine:1"}}}}})
	conf := `{"agent": {"harnesses": {"container": {"images": {"claude-code": "elsewhere:1"}}}, "projects": ` + string(projects) + `}}`
	if err := os.WriteFile(cfg, []byte(conf), 0o600); err != nil {
		t.Fatal(err)
	}
	forecastStory(t, root, "S-0001", "ready", false, "agent:\n  harness: container\n")
	d := &imageDocker{images: map[string]agentimage.Info{"mine:1": {ID: "sha256:mine"}}}
	out, code, err := runImageStatus(t, root, cfg, d)
	if code != 0 || err != nil {
		t.Fatalf("exit %d, %v\n%s", code, err, out)
	}
	if want := "claude-code harness mine:1 (host): the operator's own (--image), not judged\n"; out != want {
		t.Errorf("text output =\n%s\nwant\n%s", out, want)
	}
	out, _, _ = runImageStatus(t, root, cfg, d, "--json")
	if r := decodeImageStatus(t, out); !r.OK || r.Flai != "" || len(r.Images) != 1 || r.Images[0].State != agentimage.StateOwn || r.Images[0].Source != agentimage.SourceHost {
		t.Errorf("report = %+v", r)
	}
}

// S-0735: a project whose agents run on no container and that names no
// layer needs no image, whatever its stories in the backlog or done name:
// one line says so, --json lists none as [], and docker is not asked.
func TestImageStatusNoImageNeeded(t *testing.T) {
	t.Parallel()
	root, cfg := imageProject(t, "")
	forecastStory(t, root, "S-0001", "backlog", false, "agent:\n  harness: container\n")
	forecastStory(t, root, "S-0002", "done", false, "agent:\n  harness: container\n")
	d := &imageDocker{noDocker: true}
	out, code, err := runImageStatus(t, root, cfg, d)
	if code != 0 || err != nil || !strings.HasPrefix(out, "the project's agents need no image: ") || strings.Count(out, "\n") != 1 {
		t.Fatalf("exit %d, %v\n%s", code, err, out)
	}
	out, code, _ = runImageStatus(t, root, cfg, d, "--json")
	if code != 0 || !strings.Contains(out, `"images": []`) {
		t.Errorf("exit %d, --json =\n%s", code, out)
	}
	if r := decodeImageStatus(t, out); !r.OK || len(r.Images) != 0 {
		t.Errorf("report = %+v", r)
	}
	if len(d.asked) != 0 {
		t.Errorf("docker was asked of %v", d.asked)
	}
}

// S-0735: with docker not on PATH, flai image status fails, naming it, and
// prints no list.
func TestImageStatusNoDocker(t *testing.T) {
	t.Parallel()
	root, cfg := imageProject(t, imageLayerYAML)
	out, code, err := runImageStatus(t, root, cfg, &imageDocker{noDocker: true}, "--json")
	if code == 0 || err == nil || !errors.Is(err, agentimage.ErrNoDocker) || !strings.Contains(err.Error(), "docker is not on PATH") {
		t.Fatalf("exit %d, %v", code, err)
	}
	if out != "" {
		t.Errorf("printed %q", out)
	}
}

// imageSteps records each docker step flai image ensure or build runs, with
// the folder it runs in, and writes a line of output for it.
type imageSteps struct {
	mu   sync.Mutex
	dirs []string
	args [][]string
}

func (s *imageSteps) exec(_ context.Context, dir string, _ []string, out io.Writer, program string, args ...string) error {
	s.mu.Lock()
	s.dirs = append(s.dirs, dir)
	s.args = append(s.args, append([]string{program}, args...))
	s.mu.Unlock()
	_, err := fmt.Fprintf(out, "step output %d\n", len(s.args))
	return err
}

// useImageExec makes s run flai image ensure's and build's docker steps
// until the test ends; a test that calls it must not run in parallel.
func useImageExec(t *testing.T, s *imageSteps) {
	t.Helper()
	was := imageExec
	imageExec = s.exec
	t.Cleanup(func() { imageExec = was })
}

// developmentFlai skips a test that needs this flai to be a development
// build, as a test binary is unless its version is stamped.
func developmentFlai(t *testing.T) string {
	t.Helper()
	v := buildinfo.Get().Version
	if !buildinfo.Development(v) {
		t.Skipf("this flai is %s, not a development build", v)
	}
	return v
}

// runImageMake runs flai image with args in root, with the config file
// cfgPath, docker d, and the clock at now, and returns its output, its exit
// code as run gives it, and its error.
func runImageMake(t *testing.T, root, cfgPath string, d *imageDocker, now time.Time, args ...string) (string, int, error) {
	t.Helper()
	var out, errOut bytes.Buffer
	a := &app{configPath: cfgPath, cwd: root, runner: d, out: &out, errOut: &errOut, clock: func() time.Time { return now }}
	c := newRootCmdWith(a)
	c.SetArgs(append([]string{"image"}, args...))
	err := c.Execute()
	code := 0
	var ee *exitError
	switch {
	case errors.As(err, &ee):
		code = ee.code
	case err != nil:
		code = 1
	}
	return out.String(), code, err
}

// writeTreeDockerfiles writes bedrock's and claude-code's Dockerfiles where
// a development tree at root holds them.
func writeTreeDockerfiles(t *testing.T, root string) {
	t.Helper()
	for _, image := range []string{"bedrock", "claude-code"} {
		path := filepath.Join(root, "flai", "images", image, "Dockerfile")
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte("FROM scratch\n# "+image+"\n"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
}

// In a development tree, flai image build builds bedrock and then the
// claude-code harness image from the tree's own Dockerfiles, with the
// tree as the context, tags the harness image with the name a start runs it
// by and flai-agent-claude-code:<version>, streams each step's output to the
// output and to the log under .flai-cache/agent-images/ it names first and
// last, and exits 0 once the image is current.
func TestImageBuildInADevelopmentTree(t *testing.T) {
	version := developmentFlai(t)
	root, cfg := imageProject(t, "")
	writeTreeDockerfiles(t, root)
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	main := mainRootOf(repo)
	steps := &imageSteps{}
	useImageExec(t, steps)
	d := &imageDocker{images: map[string]agentimage.Info{
		agentimage.BedrockTag(version): {ID: "sha256:bedrock"},
		imageHarnessRef:                {ID: imageBaseID, Labels: map[string]string{agentimage.LabelFlai: version}},
	}}
	now := time.Date(2026, 11, 13, 10, 0, 0, 0, time.UTC)
	out, code, err := runImageMake(t, root, cfg, d, now, "build")
	if code != 0 || err != nil {
		t.Fatalf("exit %d, %v\n%s", code, err, out)
	}
	if len(steps.args) != 2 {
		t.Fatalf("docker steps %v", steps.args)
	}
	bedrockFile := filepath.Join(main, "flai", "images", "bedrock", "Dockerfile")
	innerFile := filepath.Join(main, "flai", "images", "claude-code", "Dockerfile")
	if got := steps.args[0][:6]; !reflect.DeepEqual(got, []string{"docker", "build", "-f", bedrockFile, "--target", "tree"}) {
		t.Errorf("bedrock step %v", steps.args[0])
	}
	if got := steps.args[1][:4]; !reflect.DeepEqual(got, []string{"docker", "build", "-f", innerFile}) {
		t.Errorf("harness step %v", steps.args[1])
	}
	harnessArgs := strings.Join(steps.args[1], " ")
	if !strings.Contains(harnessArgs, "-t "+imageHarnessRef+" -t flai-agent-claude-code:"+version+" .") || !strings.Contains(harnessArgs, "BASE="+agentimage.BedrockTag(version)) {
		t.Errorf("harness step %v", steps.args[1])
	}
	for i, dir := range steps.dirs {
		if dir != main {
			t.Errorf("step %d ran in %s, not the tree %s", i, dir, main)
		}
	}
	logPath := agentimage.LogPath(main, now)
	if !strings.HasPrefix(logPath, filepath.Join(main, ".flai-cache", "agent-images")+string(filepath.Separator)) {
		t.Errorf("log %s is not under .flai-cache/agent-images", logPath)
	}
	named := "agent image: log " + logPath + "\n"
	if !strings.HasPrefix(out, named) || !strings.HasSuffix(out, named) {
		t.Errorf("output does not name the log first and last:\n%s", out)
	}
	for _, want := range []string{"step output 1\n", "step output 2\n", "built from the development tree " + main, "claude-code harness " + imageHarnessRef + " (default): current\n"} {
		if !strings.Contains(out, want) {
			t.Errorf("output lacks %q:\n%s", want, out)
		}
	}
	logged, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatal(err)
	}
	if string(logged) != out {
		t.Errorf("log =\n%s\nwant the output\n%s", logged, out)
	}
}

// With a development flai and a main checkout that holds no Dockerfiles of
// its own, flai image ensure has nothing to pull or build from: it fails
// naming why and the Dockerfiles a tree must hold, runs no docker step,
// and still prints the images' status and names its log.
func TestImageEnsureRefusesADevelopmentFlaiWithNoTree(t *testing.T) {
	t.Parallel()
	developmentFlai(t)
	root, cfg := imageProject(t, imageLayerYAML)
	now := time.Date(2026, 11, 13, 11, 0, 0, 0, time.UTC)
	out, code, err := runImageMake(t, root, cfg, &imageDocker{}, now, "ensure")
	if code == 0 || err == nil || !strings.Contains(err.Error(), "a development flai has no release") || !strings.Contains(err.Error(), "flai/images/bedrock/Dockerfile and flai/images/claude-code/Dockerfile") {
		t.Fatalf("exit %d, %v\n%s", code, err, out)
	}
	for _, want := range []string{"claude-code harness " + imageHarnessRef + " (default): missing\n", "agent image: log "} {
		if !strings.Contains(out, want) {
			t.Errorf("output lacks %q:\n%s", want, out)
		}
	}
	if strings.Contains(out, "step output") {
		t.Errorf("a docker step ran:\n%s", out)
	}
}
