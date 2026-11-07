package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"

	"github.com/bytepunx/system-flow/flai/internal/kluster"
)

// newKlusterCmd is flai kluster: the project's local Kubernetes cluster, as
// the kluster: block of its system-flow.yaml names it (S-0577, ADR-0254).
func newKlusterCmd(a *app) *cobra.Command {
	c := &cobra.Command{
		Use:   "kluster",
		Short: "Bring the project's kluster up, take it down, and report it",
		Long: `A project names its local Kubernetes cluster in the kluster: block of its
system-flow.yaml: the kluster config file, made for bytepunx/kluster with the
k3d provider, and the commands that deploy into it (design/system/kluster.md,
ADR-0254). flai runs the kluster command in the project's main checkout, with
that config, so a story's worktree acts on main's cluster.

flai kluster up brings the cluster up and keeps its kubeconfig in
.flai-cache/kluster/kubeconfig, never in ~/.kube/config; flai kluster down
takes it down; flai kluster status reports it. They work from your shell
whatever the kluster host action says. flai guard refuses every agent and
sub-agent up and down, which change shared infrastructure on your machine;
status reads, and any agent may run it.

A project with no kluster: block, or a host without kluster on PATH, is
refused with what to add or how to install kluster.`,
		Example: `  flai kluster up
  flai kluster status
  flai kluster status --json
  flai kluster down --yes`,
		Args: cobra.NoArgs,
	}
	c.AddCommand(newKlusterUpCmd(a), newKlusterDownCmd(a), newKlusterStatusCmd(a))
	return c
}

// cluster is the project's cluster, run from its main checkout.
func (a *app) cluster() (kluster.Cluster, error) {
	repo, err := a.project()
	if err != nil {
		return kluster.Cluster{}, err
	}
	return kluster.New(mainRootOf(repo), repo.Manifest, kluster.System{})
}

func newKlusterUpCmd(a *app) *cobra.Command {
	return &cobra.Command{
		Use:   "up",
		Short: "Bring the project's cluster up with kluster up and keep its kubeconfig",
		Long: `Runs kluster --config <config> up in the project's main checkout, with
kluster's output shown as it comes, then writes the cluster's kubeconfig to
.flai-cache/kluster/kubeconfig with kluster kubeconfig. A bring-up installs
the profile's shared services and takes minutes (11m31s for authstar).`,
		Example: `  flai kluster up`,
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			c, err := a.cluster()
			if err != nil {
				return err
			}
			ctx, stop := signal.NotifyContext(cmd.Context(), os.Interrupt, syscall.SIGTERM)
			defer stop()
			if err := c.Up(ctx, a.out); err != nil {
				return err
			}
			fmt.Fprintf(a.out, "%s is up; its kubeconfig is %s, context %s\n", c.Name, kluster.KubeconfigPath(c.Root), c.Context())
			return nil
		},
	}
}

func newKlusterDownCmd(a *app) *cobra.Command {
	return &cobra.Command{
		Use:   "down",
		Short: "Take the project's cluster down with kluster down",
		Long: `Runs kluster down for the project's cluster in its main checkout, without
kluster's own prompt: flai asks first on a terminal, and --yes answers it.
Taking the cluster down deletes every namespace in it, main's and each
story's, and the shared services' data.`,
		Example: `  flai kluster down
  flai kluster down --yes`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			c, err := a.cluster()
			if err != nil {
				return err
			}
			if !a.yes {
				if !a.isTerminal() {
					return fmt.Errorf("flai kluster down deletes the cluster %s and everything in it, and asks first; pass --yes to take it down without asking", c.Name)
				}
				ok, err := a.ask(fmt.Sprintf("Take the cluster %s down, with every namespace and the shared services' data in it?", c.Name))
				if err != nil {
					return err
				}
				if !ok {
					return fmt.Errorf("the cluster %s was left up", c.Name)
				}
			}
			ctx, stop := signal.NotifyContext(cmd.Context(), os.Interrupt, syscall.SIGTERM)
			defer stop()
			if err := c.Down(ctx, a.out); err != nil {
				return err
			}
			fmt.Fprintf(a.out, "%s is down\n", c.Name)
			return nil
		},
	}
}

func newKlusterStatusCmd(a *app) *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Report the project's cluster: its name, whether it runs, and its age",
		Long: `Reads kluster status and reports the project's cluster: its name, profile,
provider, and kube context, from the kluster config, and whether kluster
lists it, whether it runs, and its age. A cluster never brought up, or taken
down, is reported as not created; that is not a failure. --json prints one
object with name, profile, provider, context, exists, running, age, and
age_seconds. It only reads, so any agent may run it.`,
		Example: `  flai kluster status
  flai kluster status --json`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			c, err := a.cluster()
			if err != nil {
				return err
			}
			s, err := c.Status(cmd.Context())
			if err != nil {
				return err
			}
			if a.jsonOut {
				return a.printJSON(s)
			}
			printKlusterStatus(a, s)
			return nil
		},
	}
}

// printKlusterStatus writes s as a short block of lines.
func printKlusterStatus(a *app, s kluster.Status) {
	state := "not created"
	switch {
	case s.Exists && s.Running:
		state = "running"
	case s.Exists:
		state = "stopped"
	}
	fmt.Fprintf(a.out, "%s: %s\n", s.Name, state)
	if s.Profile != "" {
		fmt.Fprintf(a.out, "  profile   %s\n", s.Profile)
	}
	fmt.Fprintf(a.out, "  provider  %s\n", s.Provider)
	fmt.Fprintf(a.out, "  context   %s\n", s.Context)
	if s.Age != "" {
		fmt.Fprintf(a.out, "  age       %s\n", s.Age)
	}
}
