package cmd

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"

	"github.com/bytepunx/system-flow/flai/internal/kluster"
	"github.com/bytepunx/system-flow/flai/internal/storygit"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// newKlusterCmd is flai kluster: the project's local Kubernetes cluster, as
// the kluster: block of its system-flow.yaml names it (S-0577, ADR-0254).
func newKlusterCmd(a *app) *cobra.Command {
	c := &cobra.Command{
		Use:   "kluster",
		Short: "Run the project's kluster and deploy stories and main into it",
		Long: `A project names its local Kubernetes cluster in the kluster: block of its
system-flow.yaml: the kluster config file, made for bytepunx/kluster with the
k3d provider, and the commands that build, deploy into, and tear down a
namespace of it (design/system/kluster.md, ADR-0254). flai runs the kluster
command in the project's main checkout, with that config, so a story's
worktree acts on main's cluster.

flai kluster up brings the cluster up and keeps its kubeconfig in
.flai-cache/kluster/kubeconfig, never in ~/.kube/config; flai kluster down
takes it down; flai kluster status reports it.

flai kluster deploy S-nnnn builds a story's images in its worktree, loads
them into the cluster, and runs the deploy command in the story's own
namespace, s-nnnn; flai kluster teardown S-nnnn removes it; flai kluster
promote deploys main's head into the default namespace; flai kluster
deployments lists what each namespace holds. Each job logs its steps under
.flai-cache/kluster/logs and records its outcome under
.flai-cache/kluster/deployments, and one job runs at a time per project.

They work from your shell whatever the kluster host action says. flai guard
refuses every agent and sub-agent up, down, deploy, teardown, and promote,
which change shared infrastructure on your machine; status and deployments
read, and any agent may run them.

A project with no kluster: block is refused, saying to turn the kluster host
action on, from the dashboard's Features page or with flai serve enable
kluster, which writes the block and a kluster.yaml for it. A block with no
deploy brings the cluster up and down, and deploy and promote refuse it
until kluster.deploy is set. A host without kluster on PATH is refused with
how to install kluster.`,
		Example: `  flai kluster up
  flai kluster status
  flai kluster status --json
  flai kluster deploy S-0578
  flai kluster deployments
  flai kluster teardown S-0578 --yes
  flai kluster promote
  flai kluster down --yes`,
		Args: cobra.NoArgs,
	}
	c.AddCommand(newKlusterUpCmd(a), newKlusterDownCmd(a), newKlusterStatusCmd(a),
		newKlusterDeployCmd(a), newKlusterTeardownCmd(a), newKlusterPromoteCmd(a), newKlusterDeploymentsCmd(a))
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

// klusterStory is a story a deploy or a teardown acts on: its ID, its
// status, and its worktree, empty when it has none.
type klusterStory struct {
	id, status, worktree string
}

// klusterStoryOf finds story id and its worktree; the worktree is empty when
// the story has none.
func klusterStoryOf(repo *workitem.Repo, id string) (klusterStory, error) {
	it, err := repo.Get(id)
	if err != nil {
		return klusterStory{}, err
	}
	if it.Type != "story" {
		return klusterStory{}, fmt.Errorf("%s is a %s, and flai deploys stories into the kluster; give a story's ID, such as S-0578", it.ID, it.Type)
	}
	s := klusterStory{id: it.ID, status: it.Status}
	if wt := repo.WorktreePath(it.ID); dirExists(wt) {
		s.worktree = wt
	}
	return s, nil
}

// dirExists says whether path is a folder.
func dirExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

// progressOut is where a job's progress lines go: standard output, or, with
// --json, standard error, so that standard output holds only the JSON.
func (a *app) progressOut() io.Writer {
	if a.jsonOut {
		return a.errOut
	}
	return a.out
}

func newKlusterDeployCmd(a *app) *cobra.Command {
	return &cobra.Command{
		Use:   "deploy S-nnnn",
		Short: "Build a story's images, load them into the cluster, and deploy it into its namespace",
		Long: `Deploys a story in progress or in review into its own namespace, its ID in
lower case such as s-0578, from its worktree at the worktree's head:

  build   the block's build command, in the worktree, lists the images it
          built in the file FLAI_KLUSTER_IMAGES names (skipped with no build)
  load    each image is saved with docker save and imported into every node
          of the cluster with ctr
  deploy  the block's deploy command, in the worktree

Every command gets the namespace's variables, FLAI_STORY,
FLAI_KLUSTER_NAMESPACE, FLAI_KLUSTER_PREFIX, FLAI_KLUSTER_TAG, KUBECONFIG, and
the rest that design/system/kluster.md lists. Each step's output goes to
.flai-cache/kluster/logs/<namespace>/<step>.log, and the deployment is
recorded in .flai-cache/kluster/deployments/<namespace>.json: its commit,
tag, images, each step's outcome, times, and log. A step that fails stops the
deploy, names the step, its exit code, and its log, and exits 1. A deploy
again at a new commit updates the namespace and replaces its record.

It refuses a story with no worktree, a story not in progress or in review,
a cluster with no kubeconfig (run flai kluster up), and a second job while
one runs for the project. --json prints the record, and the step lines go to
standard error.`,
		Example: `  flai kluster deploy S-0578
  flai kluster deploy S-0578 --json`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			repo, err := a.project()
			if err != nil {
				return err
			}
			c, err := kluster.New(mainRootOf(repo), repo.Manifest, kluster.System{})
			if err != nil {
				return err
			}
			s, err := klusterStoryOf(repo, args[0])
			if err != nil {
				return err
			}
			if s.status != "in-progress" && s.status != "review" {
				return fmt.Errorf("%s is %s, and flai deploys only a story in progress or in review into its namespace; move it to in-progress or review first", s.id, s.status)
			}
			if s.worktree == "" {
				return fmt.Errorf("%s has no worktree at %s to build and deploy from; open it with flai stream open %s, then deploy again", s.id, repo.WorktreePath(s.id), s.id)
			}
			commit, err := a.runner.Run(s.worktree, "git", "rev-parse", "HEAD")
			if err != nil {
				return fmt.Errorf("read the head of %s's worktree %s: %w", s.id, s.worktree, err)
			}
			ctx, stop := signal.NotifyContext(cmd.Context(), os.Interrupt, syscall.SIGTERM)
			defer stop()
			r, err := c.Deploy(ctx, repo.Manifest, s.id, s.status, s.worktree, strings.TrimSpace(commit), a.progressOut())
			return a.reportRecord(r, err)
		},
	}
}

// reportRecord prints a job's record, as text or JSON, whenever the job made
// one, and then returns the job's error.
func (a *app) reportRecord(r kluster.Record, err error) error {
	if r.Namespace != "" {
		if a.jsonOut {
			if perr := a.printJSON(r); perr != nil && err == nil {
				err = perr
			}
		} else {
			printKlusterRecord(a.out, r)
		}
	}
	return err
}

// printKlusterRecord writes r as a short block of lines.
func printKlusterRecord(w io.Writer, r kluster.Record) {
	what := r.Namespace
	if r.Story != "" {
		what = r.Story + " in " + r.Namespace
	}
	fmt.Fprintf(w, "%s: %s\n", what, r.Outcome)
	fmt.Fprintf(w, "  commit  %s\n", r.Commit)
	fmt.Fprintf(w, "  tag     %s\n", r.Tag)
	if len(r.Images) > 0 {
		fmt.Fprintf(w, "  images  %s\n", strings.Join(r.Images, ", "))
	}
	for _, st := range r.Steps {
		line := fmt.Sprintf("  %-8s %s", st.Name, st.Outcome)
		if st.Outcome == kluster.OutcomeFailed {
			line += fmt.Sprintf(", exit code %d", st.ExitCode)
		}
		if st.Log != "" {
			line += ", log " + st.Log
		}
		fmt.Fprintln(w, line)
	}
}

func newKlusterTeardownCmd(a *app) *cobra.Command {
	return &cobra.Command{
		Use:   "teardown S-nnnn",
		Short: "Remove a story's deployment: the teardown command, then its namespace",
		Long: `Tears a story's namespace down: runs the block's teardown command, when it
names one, with the namespace's variables at the commit its record names, in
the story's worktree or the main checkout when the worktree is gone; then
deletes the namespace with kubectl delete namespace --wait, which passes over
a namespace already gone. Both are logged to the namespace's teardown log.
Then the record, the namespace's state folder, and its logs are removed.

A story never deployed has no record: flai says so, asks nothing, and runs
nothing, and exits 0. Otherwise flai asks first on
a terminal, and --yes answers it; without a terminal it needs --yes. A
teardown that fails keeps the record, marked failed, and exits 1. --json
prints the record removed, or {"story": ..., "deployed": false}.`,
		Example: `  flai kluster teardown S-0578
  flai kluster teardown S-0578 --yes --json`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			repo, err := a.project()
			if err != nil {
				return err
			}
			c, err := kluster.New(mainRootOf(repo), repo.Manifest, kluster.System{})
			if err != nil {
				return err
			}
			s, err := klusterStoryOf(repo, args[0])
			if err != nil {
				return err
			}
			notDeployed := func() error {
				if a.jsonOut {
					return a.printJSON(map[string]any{"story": s.id, "deployed": false})
				}
				fmt.Fprintf(a.out, "%s has no deployment recorded; nothing to tear down\n", s.id)
				return nil
			}
			ns, err := kluster.StoryNamespace(s.id)
			if err != nil {
				return err
			}
			if _, ok, err := kluster.ReadRecord(c.Root, ns); err != nil {
				return err
			} else if !ok {
				return notDeployed()
			}
			if !a.yes {
				if !a.isTerminal() {
					return fmt.Errorf("flai kluster teardown deletes %s's namespace and what it made in the shared services, and asks first; pass --yes to tear it down without asking", s.id)
				}
				ok, err := a.ask(fmt.Sprintf("Tear down %s's namespace in the cluster %s?", s.id, c.Name))
				if err != nil {
					return err
				}
				if !ok {
					return fmt.Errorf("%s's namespace was left as it is", s.id)
				}
			}
			ctx, stop := signal.NotifyContext(cmd.Context(), os.Interrupt, syscall.SIGTERM)
			defer stop()
			r, err := c.Teardown(ctx, repo.Manifest, s.id, s.worktree, a.progressOut())
			if errors.Is(err, kluster.ErrNotDeployed) {
				return notDeployed()
			}
			if err != nil {
				return a.reportRecord(r, err)
			}
			if a.jsonOut {
				return a.printJSON(r)
			}
			fmt.Fprintf(a.out, "%s's namespace %s is torn down\n", s.id, r.Namespace)
			return nil
		},
	}
}

func newKlusterPromoteCmd(a *app) *cobra.Command {
	return &cobra.Command{
		Use:   "promote",
		Short: "Deploy main's head into the default namespace",
		Long: `Deploys the head of the main branch into the default namespace,
kluster.namespace or the project's name: it makes, or moves, a detached
worktree of that head at .flai-cache/kluster/main, then builds, loads, and
deploys there as flai kluster deploy does for a story, with FLAI_STORY and
FLAI_KLUSTER_PREFIX empty and the tag main-<commit>. The outcome is recorded
under the default namespace. A step that fails names the step, its exit code,
and its log, and exits 1. --json prints the record, and the step lines go to
standard error.`,
		Example: `  flai kluster promote
  flai kluster promote --json`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			repo, err := a.project()
			if err != nil {
				return err
			}
			root := mainRootOf(repo)
			c, err := kluster.New(root, repo.Manifest, kluster.System{})
			if err != nil {
				return err
			}
			branch, err := storygit.MainBranch(a.runner, root)
			if err != nil {
				return fmt.Errorf("find the main branch of %s to promote: %w", root, err)
			}
			ctx, stop := signal.NotifyContext(cmd.Context(), os.Interrupt, syscall.SIGTERM)
			defer stop()
			r, err := c.Promote(ctx, repo.Manifest, branch, a.progressOut())
			return a.reportRecord(r, err)
		},
	}
}

func newKlusterDeploymentsCmd(a *app) *cobra.Command {
	return &cobra.Command{
		Use:   "deployments",
		Short: "List each namespace's deployment, the default namespace's first",
		Long: `Lists the deployments flai recorded under .flai-cache/kluster/deployments:
the default namespace's first, then each story's by its number. Each row
gives the namespace, the story, the commit, the outcome, and when it ended. A
namespace a job acts on now shows that job; a record a job left running
when it ended without its outcome shows interrupted. --json prints them as an
array of records, each with default and, while a job runs on it, job. It only
reads, so any agent may run it.`,
		Example: `  flai kluster deployments
  flai kluster deployments --json`,
		Args: cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			repo, err := a.project()
			if err != nil {
				return err
			}
			root := mainRootOf(repo)
			if _, err := kluster.New(root, repo.Manifest, kluster.System{}); err != nil {
				return err
			}
			list, err := kluster.Deployments(root, repo.Manifest.KlusterNamespace())
			if err != nil {
				return err
			}
			if a.jsonOut {
				return a.printJSON(list)
			}
			if len(list) == 0 {
				fmt.Fprintln(a.out, "no deployments")
				return nil
			}
			tw := tabwriter.NewWriter(a.out, 0, 0, 2, ' ', 0)
			fmt.Fprintln(tw, "NAMESPACE\tSTORY\tCOMMIT\tOUTCOME\tENDED")
			for _, d := range list {
				story := d.Story
				if d.Default {
					story = "(main)"
				}
				outcome := d.Outcome
				if d.Job != nil {
					outcome = d.Job.Kind + " running"
				}
				ended := "-"
				if !d.EndedAt.IsZero() {
					ended = d.EndedAt.UTC().Format(time.RFC3339)
				}
				commit := d.Commit
				if len(commit) > 12 {
					commit = commit[:12]
				}
				fmt.Fprintf(tw, "%s\t%s\t%s\t%s\t%s\n", d.Namespace, story, commit, outcome, ended)
			}
			return tw.Flush()
		},
	}
}
