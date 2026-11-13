package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

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
current for the project's flai.`,
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
	return c
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
	images := agentimage.Needed(repo.Manifest, agents, host.Images)
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
	statuses, err := agentimage.Check(root, images, agentimage.Docker(a.runner, program), flai)
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
		printImageStatus(a, statuses)
	}
	if bad == 0 {
		return nil
	}
	noun := "agent images are"
	if bad == 1 {
		noun = "agent image is"
	}
	return &exitError{code: 1, msg: fmt.Sprintf("%d %s missing or stale: build them with make agent-images in the system-flow repository, or name another with flai serve agent harness container --image <inner>=<image>", bad, noun)}
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

// printImageStatus prints one line per image in statuses, or, with none,
// that the project's agents need no image.
func printImageStatus(a *app, statuses []agentimage.Status) {
	if len(statuses) == 0 {
		_, _ = fmt.Fprintf(a.out, "the project's agents need no image: no agent runs on the container harness and %s names no agent.image.layer\n", manifest.File)
		return
	}
	for _, s := range statuses {
		from := string(s.Source)
		if s.Dockerfile != "" {
			from += " from " + s.Dockerfile
		}
		_, _ = fmt.Fprintf(a.out, "%s %s %s (%s): %s\n", s.Inner, s.Kind, s.Ref, from, imageState(s))
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
