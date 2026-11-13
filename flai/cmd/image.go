package cmd

import (
	"context"
	"fmt"
	"io"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/spf13/cobra"

	"github.com/bytepunx/system-flow/flai/images"
	"github.com/bytepunx/system-flow/flai/internal/agentimage"
	"github.com/bytepunx/system-flow/flai/internal/buildinfo"
	"github.com/bytepunx/system-flow/flai/internal/harness"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// flai image is the agent images a project's container agents run in
// (ADR-0295). It is not flai agent, the project's default agent.
func newImageCmd(a *app) *cobra.Command {
	c := &cobra.Command{
		Use:   "image",
		Short: "The agent images a project's container agents run in",
		Long: `An agent on the container harness runs in an image: its inner harness's, the
host's own when flai serve agent harness container --image <inner>=<image>
names one, else flai's default, and the project's layer over it when
system-flow.yaml names agent.image.layer (ADR-0295). flai image status says
which images the project's agents need and whether each is on this host and
current for the project's flai. flai image ensure makes each harness image
that is missing or stale, pulling it by the digest its release's signed list
names, else building it on this host, then builds the project's layer over its
harness image when the layer is missing or stale or its harness image was just
made, with agent.image.args as build arguments and BASE set to the harness
image; flai image build builds them all without the pull, and flai image build
--layer builds the layer alone.`,
		Args: cobra.NoArgs,
	}
	c.AddCommand(&cobra.Command{
		Use:   "status",
		Short: "Say which agent images the project's agents need, and whether each is on this host and current for its flai",
		Long: `flai image status lists the images the agents of the project in the working
directory need: those of its default agent, its roles' agents, and its open
stories' agents (ready, in progress, in review) on the container harness, and
the project's layer when system-flow.yaml names agent.image.layer. For each it
asks docker whether the image is on this host and reads its labels:

  missing   the image is not on this host
  current   its labels are those it must carry: its flai label the project's
            flai (flai.command's when the manifest names one, else this
            flai), and for a layer its recipe, over its Dockerfile, and its
            base, the harness image now under it
  stale     a label is other than it must be; the line says which and why
  own       the operator's own image (--image), which is never judged

It exits non-zero when any image is missing or stale, after the list.`,
		Example: `  flai image status
  flai image status --json`,
		Args: cobra.NoArgs,
		RunE: func(*cobra.Command, []string) error { return a.imageStatus() },
	})
	c.AddCommand(&cobra.Command{
		Use:   "ensure",
		Short: "Make each agent image the project's agents need that is missing or stale: pull a harness image by its signed digest, else build it on this host, then build the project's layer over it",
		Long: `flai image ensure makes each harness image the agents of the project in the
working directory need, as flai image status lists them, that is missing or
stale for the project's flai. It takes one of three ways, and says which:

  pulled     by the index digest the signed digest list of the flai release
             names (flai-agent-<inner>_<version>.digests on the release
             flai v<version>), verified with the release key; a list or an
             image that does not verify is refused, and nothing is built in
             its place
  built      on this host, when the registry or the release cannot be reached
             or the pull fails, from the Dockerfiles this flai carries and the
             release's Linux flai and pilot, checked against its signed
             checksums.txt
  tree       for a development flai, from the development tree's own
             flai/images/bedrock/Dockerfile and flai/images/<inner>/Dockerfile
             and source, never pulled for

Each harness image is tagged with the name a start runs it by and with
flai-agent-<inner>:<version>. A current image is left as it is, and the
operator's own image (--image) is never touched.

After the harness images it builds the project's layer, when system-flow.yaml
names agent.image.layer, over the harness image of its inner harness, the
operator's own when --image names one, else flai's default: when the layer is
missing or stale, or when that harness image was just pulled or built. It
builds it from the layer's Dockerfile with the main checkout as the context,
BASE set to the harness image and each of agent.image.args passed as a build
argument, tags it flai-agent-<key>-<inner>:local, and labels it with its
recipe, over the Dockerfile and the arguments, and the ID of the harness image
under it, so that flai image status reports it stale when the Dockerfile, an
argument, or its base changes. A layer whose harness image is not on this host
is not built: it fails naming the harness image.

Every docker step's output streams here and to a log file under
.flai-cache/agent-images/ in the main checkout, which the first and last lines
name. A step that writes nothing for 10 minutes is stopped; a step that fails
on the network (a connection reset, a timeout, a DNS or TLS failure, a
registry 5xx) is tried twice more, after 30 s and then 2 min. After it, the
images' status is printed as flai image status prints it, and it exits
non-zero when a harness image or a layer flai makes is still missing or stale.`,
		Example: `  flai image ensure`,
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return a.imageMake(cmd.Context(), agentimage.ModeEnsure)
		},
	})
	var layerOnly bool
	build := &cobra.Command{
		Use:   "build",
		Short: "Build the agent images on this host, current or not, without the pull, or with --layer the project's layer alone",
		Long: `flai image build builds on this host, without the pull, each harness image of
flai's default the agents of the project in the working directory need, or,
when they need none, the default image of each inner harness the host names
no image of its own for (claude-code's), current or not, then the project's
layer over its harness image when system-flow.yaml names agent.image.layer.

A released flai builds the harness images from the Dockerfiles it carries and
its release's Linux flai and pilot, checked against the release's signed
checksums.txt; a development flai builds them from the development tree's own
flai/images/bedrock/Dockerfile and flai/images/<inner>/Dockerfile and source,
when the main checkout holds them. It builds the layer as flai image ensure
does, with agent.image.args as build arguments and BASE set to the harness
image, labelled with its recipe and base. It tags, logs, stops, and retries as
flai image ensure does, prints the images' status after it, and exits
non-zero when an image it makes is still missing or stale.

With --layer it builds only the project's layer, one per inner harness the
agents need, current or not, over the harness image already on this host, and
pulls or builds no harness image, as scripts/agent-image.sh --layer does. It
refuses when system-flow.yaml names no agent.image.layer, and when the harness
image is not on this host, naming it and flai image ensure, which makes it.`,
		Example: `  flai image build
  flai image build --layer`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if layerOnly {
				return a.imageMake(cmd.Context(), agentimage.ModeLayer)
			}
			return a.imageMake(cmd.Context(), agentimage.ModeBuild)
		},
	}
	build.Flags().BoolVar(&layerOnly, "layer", false, "build only the project's layer (agent.image.layer), current or not, over the harness image already on this host, and no harness image")
	c.AddCommand(build)
	return c
}

// imageExec runs each docker step of flai image ensure and build; tests
// that replace it restore it and do not run in parallel.
var imageExec agentimage.Exec = agentimage.SystemExec

// imageMake is flai image ensure, with agentimage.ModeEnsure, flai image
// build, with agentimage.ModeBuild, and flai image build --layer, with
// agentimage.ModeLayer: it makes the agent images of the project in the
// working directory as mode asks, the harness images and then the project's
// layer over them, writing what it does to the output and to a log under the
// main checkout's .flai-cache, prints their status after it, and fails when
// one it makes is still missing or stale.
func (a *app) imageMake(ctx context.Context, mode agentimage.Mode) (err error) {
	if ctx == nil {
		ctx = context.Background()
	}
	repo, err := a.project()
	if err != nil {
		return err
	}
	cfg, _, err := a.loadConfig()
	if err != nil {
		return err
	}
	if mode == agentimage.ModeLayer && repo.Manifest.Layer() == "" {
		return fmt.Errorf("build the project's layer image: %s names no agent.image.layer, so there is no layer to build; name its Dockerfile under agent.image.layer, or run flai image build for the harness images", manifest.File)
	}
	root := mainRootOf(repo)
	agents, err := projectAgents(repo)
	if err != nil {
		return err
	}
	host := a.harnessHost(cfg.Agent.ForProject(root), harness.Container)
	program := host.Program
	if program == "" {
		program = "docker"
	}
	needed := agentimage.Needed(repo.Manifest, agents, host.Images)
	if mode == agentimage.ModeBuild {
		needed = buildImages(needed, host.Images)
	}

	logPath := agentimage.LogPath(root, a.now())
	if err := os.MkdirAll(filepath.Dir(logPath), 0o755); err != nil {
		return fmt.Errorf("make the folder of the agent image log %s: %w; check that the project's .flai-cache can be written", logPath, err)
	}
	logFile, err := os.Create(logPath)
	if err != nil {
		return fmt.Errorf("open the agent image log %s: %w; check that the project's .flai-cache can be written", logPath, err)
	}
	w := io.MultiWriter(a.out, logFile)
	_, _ = fmt.Fprintf(w, "agent image: log %s\n", logPath)
	defer func() {
		if err != nil {
			_, _ = fmt.Fprintf(logFile, "agent image: failed: %v\n", err)
		}
		_, _ = fmt.Fprintf(w, "agent image: log %s\n", logPath)
		if cerr := logFile.Close(); cerr != nil && err == nil {
			err = fmt.Errorf("close the agent image log %s: %w; check that the project's .flai-cache can be written", logPath, cerr)
		}
	}()

	if len(needed) == 0 {
		printImageStatus(w, nil)
		return nil
	}
	_, after, err := a.ensureImages(ctx, root, repo.Manifest, program, mode, needed, w)
	if after != nil {
		printImageStatus(w, after)
	}
	if err != nil {
		return err
	}
	bad := len(imagesLeft(mode, after))
	if bad == 0 {
		return nil
	}
	noun := "agent images flai makes are"
	if bad == 1 {
		noun = "agent image flai makes is"
	}
	return &exitError{code: 1, msg: fmt.Sprintf("%d %s still missing or stale after %s: read why in the log %s, then run it again", bad, noun, imageCommand(mode), logPath)}
}

// ensureImages makes needed, agent images of the project at root whose
// manifest is m, as mode asks, through the docker client program, writing
// what it does to w: it judges them against the project's flai, makes those
// mode makes, and judges them again. It is what it did with each, their
// statuses after, nil when they could not be judged again, and the first
// error of the version, the judgements, and the making.
func (a *app) ensureImages(ctx context.Context, root string, m manifest.Manifest, program string, mode agentimage.Mode, needed []agentimage.Image, w io.Writer) ([]agentimage.Outcome, []agentimage.Status, error) {
	version, err := agentimage.ProjectFlai(root, m, a.runner, buildinfo.Get().Version)
	if err != nil {
		return nil, nil, err
	}
	flai := func() (string, error) { return version, nil }
	inspect := agentimage.Docker(a.runner, program)
	statuses, err := agentimage.Check(root, needed, inspect, flai)
	if err != nil {
		return nil, nil, err
	}
	src := a.imageSource(root, version)
	runner := agentimage.Runner{Program: program, Exec: imageExec, Out: w}
	releases := a.dashboardReleases()
	releases.Token = a.githubToken()
	puller := agentimage.Puller{Releases: releases, CacheDir: string(a.serveDir()), Run: runner, Inspect: inspect}
	builder := agentimage.Builder{Releases: releases, Run: runner, Inspect: inspect, Arch: agentimage.EngineArch(a.runner, program)}
	ensurer := agentimage.Ensurer{Pull: puller.Pull, Build: builder.Build, Layer: builder.Layer, Root: root, Out: w}
	outcomes, ensureErr := ensurer.Ensure(ctx, mode, statuses, src)
	after, checkErr := agentimage.Check(root, needed, inspect, flai)
	if ensureErr != nil {
		return outcomes, after, ensureErr
	}
	return outcomes, after, checkErr
}

// imageEnsure is the work of flai serve's image job (serve.ImageEnsure,
// ADR-0295 §5): flai image ensure for images, of the project at root, its
// main checkout, writing what it does and the images' status after to log.
// It fails naming each image it makes that is still missing or stale after it.
func (a *app) imageEnsure(ctx context.Context, root string, images []agentimage.Image, log io.Writer) ([]agentimage.Outcome, error) {
	m, err := manifest.Load(filepath.Join(root, manifest.File))
	if err != nil {
		return nil, fmt.Errorf("read the %s of the project at %s for its agent images: %w; mend it (flai check says what is wrong), then start the agent again", manifest.File, root, err)
	}
	cfg, _, err := a.loadConfig()
	if err != nil {
		return nil, err
	}
	program := a.harnessHost(cfg.Agent.ForProject(root), harness.Container).Program
	if program == "" {
		program = "docker"
	}
	outcomes, after, err := a.ensureImages(ctx, root, m, program, agentimage.ModeEnsure, images, log)
	if after != nil {
		printImageStatus(log, after)
	}
	if err != nil {
		return outcomes, err
	}
	if left := imagesLeft(agentimage.ModeEnsure, after); len(left) > 0 {
		return outcomes, fmt.Errorf("flai image ensure left the agent images %s for %s: read why in the image job's log, mend it, then run flai image ensure in the project", strings.Join(left, ", "), root)
	}
	return outcomes, nil
}

// imagesLeft are the images of after that mode makes and that are still
// missing or stale, each as its ref and its state.
func imagesLeft(mode agentimage.Mode, after []agentimage.Status) []string {
	var left []string
	for _, s := range after {
		if imageMade(mode, s) && (s.State == agentimage.StateMissing || s.State == agentimage.StateStale) {
			left = append(left, s.Ref+" "+string(s.State))
		}
	}
	return left
}

// imageMade reports whether flai image makes s in mode: the project's layer
// always, and a harness image of flai's default unless mode builds the layer alone.
func imageMade(mode agentimage.Mode, s agentimage.Status) bool {
	if s.Kind == agentimage.KindLayer {
		return true
	}
	return mode != agentimage.ModeLayer && s.Kind == agentimage.KindHarness && s.Source == agentimage.SourceDefault
}

// imageCommand is the command that runs mode: flai image ensure, flai image
// build, or flai image build --layer.
func imageCommand(mode agentimage.Mode) string {
	if mode == agentimage.ModeLayer {
		return "flai image build --layer"
	}
	return "flai image " + string(mode)
}

// buildImages are the images flai image build makes of needed, what the
// project's agents need: needed itself when it holds a harness image of
// flai's default, else needed and the harness image of each inner harness
// of harness.DefaultImages that hostImages names no image of its own for.
func buildImages(needed []agentimage.Image, hostImages map[string]string) []agentimage.Image {
	for _, img := range needed {
		if img.Kind == agentimage.KindHarness && img.Source == agentimage.SourceDefault {
			return needed
		}
	}
	out := slices.Clone(needed)
	for _, inner := range slices.Sorted(maps.Keys(harness.DefaultImages)) {
		if hostImages[inner] == "" {
			out = append(out, agentimage.Harness(hostImages, inner))
		}
	}
	return out
}

// imageSource is what the harness images of flai version are made from:
// for a development version, the development tree at root when it holds
// bedrock's and claude-code's Dockerfiles, with its commit and that
// commit's date as git gives them ("" when git cannot); otherwise the
// release version.
func (a *app) imageSource(root, version string) agentimage.BuildSource {
	src := agentimage.BuildSource{Version: version}
	if !buildinfo.Development(version) {
		return src
	}
	for _, image := range []string{images.Bedrock, harness.ClaudeCode} {
		if _, err := os.Stat(filepath.Join(root, filepath.FromSlash(images.Path(image)))); err != nil {
			return src
		}
	}
	src.Tree = root
	if out, err := a.runner.Run(root, "git", "rev-parse", "--short", "HEAD"); err == nil {
		src.Commit = strings.TrimSpace(out)
	}
	if out, err := a.runner.Run(root, "git", "log", "-1", "--format=%cI"); err == nil {
		src.Date = strings.TrimSpace(out)
	}
	return src
}

// imageStatusReport is what flai image status --json prints: the flai
// version the images must hold ("" when no image on this host was judged),
// each image's status, and whether none is missing or stale.
type imageStatusReport struct {
	Flai   string              `json:"flai"`
	Images []agentimage.Status `json:"images"`
	OK     bool                `json:"ok"`
}

// imageStatus prints the status of each image the agents of the project in
// the working directory need, and fails when any is missing or stale.
func (a *app) imageStatus() error {
	repo, err := a.project()
	if err != nil {
		return err
	}
	cfg, _, err := a.loadConfig()
	if err != nil {
		return err
	}
	root := mainRootOf(repo)
	agents, err := projectAgents(repo)
	if err != nil {
		return err
	}
	// the container harness as flai serve starts it for this project
	// (ADR-0261): its images as the operator names them, before the defaults,
	// which agentimage.Harness fills in, and its docker client
	host := a.harnessHost(cfg.Agent.ForProject(root), harness.Container)
	needed := agentimage.Needed(repo.Manifest, agents, host.Images)
	program := host.Program
	if program == "" {
		program = "docker"
	}
	var want string
	flai := func() (string, error) {
		v, err := agentimage.ProjectFlai(root, repo.Manifest, a.runner, buildinfo.Get().Version)
		want = v
		return v, err
	}
	statuses, err := agentimage.Check(root, needed, agentimage.Docker(a.runner, program), flai)
	if err != nil {
		return err
	}
	bad := 0
	for _, s := range statuses {
		if s.State == agentimage.StateMissing || s.State == agentimage.StateStale {
			bad++
		}
	}
	if a.jsonOut {
		if err := a.printJSON(imageStatusReport{Flai: want, Images: statuses, OK: bad == 0}); err != nil {
			return err
		}
	} else {
		printImageStatus(a.out, statuses)
	}
	if bad == 0 {
		return nil
	}
	noun := "agent images are"
	if bad == 1 {
		noun = "agent image is"
	}
	return &exitError{code: 1, msg: fmt.Sprintf("%d %s missing or stale: make them with flai image ensure, or name another with flai serve agent harness container --image <inner>=<image>", bad, noun)}
}

// projectAgents is every agent the project's manifest starts, then the
// agent of each of its open stories, ready, in progress, or in review, as
// flai serve starts it (AgentOf).
func projectAgents(repo *workitem.Repo) ([]*manifest.Agent, error) {
	agents := agentimage.Agents(repo.Manifest)
	items, err := repo.List(false)
	if err != nil {
		return nil, fmt.Errorf("list the project's stories for the agents that need an image: %w", err)
	}
	for _, it := range items {
		if it.Type != workitem.Story {
			continue
		}
		switch it.Status {
		case workitem.Ready, workitem.InProgress, workitem.Review:
			agents = append(agents, repo.AgentOf(it))
		}
	}
	return agents, nil
}

// printImageStatus prints to w one line per image in statuses, or, with
// none, that the project's agents need no image.
func printImageStatus(w io.Writer, statuses []agentimage.Status) {
	if len(statuses) == 0 {
		_, _ = fmt.Fprintf(w, "the project's agents need no image: no agent runs on the container harness and %s names no agent.image.layer\n", manifest.File)
		return
	}
	for _, s := range statuses {
		from := string(s.Source)
		if s.Dockerfile != "" {
			from += " from " + s.Dockerfile
		}
		_, _ = fmt.Fprintf(w, "%s %s %s (%s): %s\n", s.Inner, s.Kind, s.Ref, from, imageState(s))
	}
}

// imageState is what s's line says of its state: why it is stale, and the
// commit a development image was built from when its labels name one.
func imageState(s agentimage.Status) string {
	var said string
	switch s.State {
	case agentimage.StateOwn:
		said = "the operator's own (--image), not judged"
	case agentimage.StateStale:
		reasons := make([]string, 0, len(s.Reasons))
		for _, r := range s.Reasons {
			reasons = append(reasons, r.Message)
		}
		said = "stale: " + strings.Join(reasons, "; ")
	default:
		said = string(s.State)
	}
	if s.Labels.Commit != "" {
		said += " (built from commit " + s.Labels.Commit + ")"
	}
	return said
}
