package cmd

import (
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/spf13/cobra"

	"github.com/bytepunx/system-flow/flai/internal/config"
)

// flai serve agent image: how a project's agent images are made on this host
// (ADR-0295 point 7), the operator's build command in place of flai's own
// pull and build, and whether flai serve makes a missing or stale image by
// itself. Both are the project's own only: there is no --all-projects.
func newServeAgentImageCmd(a *app) *cobra.Command {
	var project string
	c := &cobra.Command{
		Use:   "image [--project <root>]",
		Short: "How the project's agent images are made: your build command in place of flai's, and whether flai serve makes them by itself",
		Long: `Prints how the agent images of the project in the working directory, or of
--project <root>, are made on this host (ADR-0295): the build command you
named for it, or flai's own pull and build, and whether flai serve makes a
missing or stale image by itself.

flai serve agent image command names your build command, which flai runs in
place of its own pull and build. flai serve agent image auto off turns the
automatic image job off for the project. Both are kept in the host's
configuration, as agent.projects.<root>.image, and are the project's own: a
project's manifest never names them, and there is no --all-projects.`,
		Example: `  flai serve agent image
  flai serve agent image command -- scripts/agent-image.sh
  flai serve agent image command --clear
  flai serve agent image auto off
  flai serve agent image auto on --project /home/me/src/other`,
		Args: cobra.NoArgs,
		RunE: func(*cobra.Command, []string) error {
			root, err := a.imageScope(project)
			if err != nil {
				return err
			}
			return a.showAgentImage(root, true, true)
		},
	}
	c.PersistentFlags().StringVar(&project, "project", "", "the main checkout of the project whose image settings these are, not the one in the working directory")
	var drop bool
	command := &cobra.Command{
		Use:   "command [--project <root>] [--clear] [-- <program> [args...]]",
		Short: "Name your build command for the project's agent images, which flai runs in place of its own pull and build",
		Long: `Names the build command of the project's agent images, as an argument list
after --. flai image ensure and flai serve's image job then run it in place
of flai's own pull and build: with your rights, in the project's main
checkout, never through a shell, with the FLAI_IMAGE_* variables saying what
the images must be (FLAI_IMAGE_VERSION, FLAI_IMAGE_INNER, the tags in
FLAI_IMAGE_BEDROCK, FLAI_IMAGE_HARNESS, and FLAI_IMAGE_LAYER, the recipes in
FLAI_IMAGE_RECIPE_*, and the images missing or stale in FLAI_IMAGE_STALE).
flai checks the images and their labels after it, and a non-zero exit, or an
image missing or with other labels, fails the run.

--clear removes it, so that flai's own pull and build make the images again.
With neither, it prints the command set, or says that flai's own pull and
build make them. It writes the project's own setting, of the working
directory or of --project <root>, kept in the host's configuration as
agent.projects.<root>.image.command (ADR-0295).`,
		Example: `  flai serve agent image command -- scripts/agent-image.sh
  flai serve agent image command --project /home/me/src/other -- /home/me/bin/build-images --push
  flai serve agent image command
  flai serve agent image command --clear`,
		Args: cobra.ArbitraryArgs,
		RunE: func(_ *cobra.Command, args []string) error {
			if drop && len(args) > 0 {
				return errors.New("give a program or --clear, not both")
			}
			if len(args) > 0 && strings.TrimSpace(args[0]) == "" {
				return errors.New("the program is empty")
			}
			root, err := a.imageScope(project)
			if err != nil {
				return err
			}
			if !drop && len(args) == 0 {
				return a.showAgentImage(root, true, false)
			}
			if err := a.setAgentImage(root, func(i *config.ProjectImage) { i.Command = slices.Clone(args) }); err != nil {
				return err
			}
			return a.showAgentImage(root, true, false)
		},
	}
	command.Flags().BoolVar(&drop, "clear", false, "remove the command, so that flai's own pull and build make the images again")
	auto := &cobra.Command{
		Use:   "auto [on|off] [--project <root>]",
		Short: "Whether flai serve makes the project's missing or stale agent images by itself; off leaves them to you",
		Long: `Turns the automatic image job on or off for the project, of the working
directory or of --project <root>. On, the default, flai serve makes a missing
or stale agent image by itself, with your build command when one is named,
and holds the starts that need it until it is made. Off, a start refuses a
missing or stale image, naming it and flai image ensure, which still runs by
hand. With no argument it prints on or off.

It is kept in the host's configuration as agent.projects.<root>.image.auto:
off is stored as false, and on removes it, the default (ADR-0295).`,
		Example: `  flai serve agent image auto off
  flai serve agent image auto on
  flai serve agent image auto`,
		Args:      cobra.MaximumNArgs(1),
		ValidArgs: []string{"on", "off"},
		RunE: func(_ *cobra.Command, args []string) error {
			var off *bool
			if len(args) == 1 {
				switch args[0] {
				case "on":
				case "off":
					off = new(bool)
				default:
					return fmt.Errorf("%q is neither on nor off", args[0])
				}
			}
			root, err := a.imageScope(project)
			if err != nil {
				return err
			}
			if len(args) == 1 {
				if err := a.setAgentImage(root, func(i *config.ProjectImage) { i.Auto = off }); err != nil {
					return err
				}
			}
			return a.showAgentImage(root, false, true)
		},
	}
	c.AddCommand(command, auto)
	return c
}

// imageScope is the project whose image settings flai serve agent image
// reads and writes: the project of --project, or of the working directory,
// by its main checkout. Outside a project, with no --project, it refuses,
// naming the flag; the settings are each project's own (ADR-0295).
func (a *app) imageScope(project string) (string, error) {
	if project != "" {
		return a.agentScope(project, false)
	}
	repo, err := a.project()
	if err != nil {
		return "", fmt.Errorf("%w; or name a project with --project <root>: the image settings are each project's own", err)
	}
	return mainRootOf(repo), nil
}

// setAgentImage changes the project at root's image settings with change and
// saves the host configuration; an entry left setting nothing is removed.
func (a *app) setAgentImage(root string, change func(*config.ProjectImage)) error {
	cfg, path, err := a.loadConfig()
	if err != nil {
		return err
	}
	p := cfg.Agent.Project(root)
	change(&p.Image)
	if cfg.Agent, err = cfg.Agent.WithProject(root, p); err != nil {
		return err
	}
	return config.Save(path, cfg)
}

// agentImageShown is the project at root's image settings as show, the
// settings page, and flai serve agent image --json give them: the command,
// empty when flai's own pull and build make the images, and whether flai
// serve makes them by itself, with auto_set saying whether the operator
// turned it off (ADR-0295).
func agentImageShown(agent config.AgentStart, root string) map[string]any {
	return map[string]any{"command": nonNil(agent.ImageCommand(root)), "auto": agent.AutoImage(root), "auto_set": agent.Project(root).Image.Auto != nil}
}

// showAgentImage prints the project at root's image settings, the command,
// the automatic job, or both.
func (a *app) showAgentImage(root string, command, auto bool) error {
	cfg, _, err := a.loadConfig()
	if err != nil {
		return err
	}
	if a.jsonOut {
		shown := agentImageShown(cfg.Agent, root)
		shown["project"] = root
		return a.printJSON(shown)
	}
	if command {
		printImageCommand(a, cfg.Agent, root)
	}
	if auto {
		printImageAuto(a, cfg.Agent, root)
	}
	return nil
}

// printImageCommand prints the project at root's build command, or that
// flai's own pull and build make its images.
func printImageCommand(a *app, agent config.AgentStart, root string) {
	if cmd := agent.ImageCommand(root); len(cmd) > 0 {
		fmt.Fprintf(a.out, "image command: %s (this project)\n  run in place of flai's own pull and build, as you, in the project's main checkout, never through a shell, with the FLAI_IMAGE_* variables; --clear brings flai's back\n", quoteArgs(cmd))
		return
	}
	fmt.Fprintln(a.out, "image command: none, so flai's own pull and build make the agent images; flai serve agent image command -- <program> [args...] names one of yours")
}

// printImageAuto prints whether flai serve makes the project at root's
// missing or stale images by itself.
func printImageAuto(a *app, agent config.AgentStart, root string) {
	if agent.AutoImage(root) {
		fmt.Fprintln(a.out, "image auto: on (default), so flai serve makes a missing or stale agent image by itself; flai serve agent image auto off leaves them to you")
		return
	}
	fmt.Fprintln(a.out, "image auto: off, so a start refuses a missing or stale agent image; flai image ensure makes them, and flai serve agent image auto on turns the job back on")
}
