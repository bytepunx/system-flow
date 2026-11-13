package cmd

import (
	"errors"
	"fmt"
	"io"
	"maps"
	"os"
	"os/exec"
	"os/signal"
	"slices"
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
.flai-cache/kluster/logs, which flai kluster log S-nnnn --step <step>
prints, and records its outcome under .flai-cache/kluster/deployments, and
one job runs at a time per project. flai kluster exec S-nnnn -- <command>
runs a command in the story's worktree with its namespace's variables and
the namespace's own kubeconfig, while the namespace has a record, and flai
kluster env S-nnnn prints those variables, as export lines or JSON.

They work from your shell whatever the kluster host action says. flai guard
refuses every agent and sub-agent up, down, deploy, teardown, and promote,
which change shared infrastructure on your machine; status and deployments
read, and any agent may run them. exec and log reach a story's namespace
without changing the cluster's shared infrastructure, and flai guard lets a
story's agent and its sub-agents run them for their own story alone
(ADR-0292 §2). A story's agent asks flai serve for its own story's deploy or
teardown instead, with flai kluster deploy or teardown S-nnnn --job, or the
MCP tools kluster_deploy and kluster_teardown, which run nothing themselves
(ADR-0292).

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
  flai kluster log S-0578 --step deploy
  flai kluster exec S-0578 -- kubectl get pods
  eval "$(flai kluster env S-0578)"
  flai kluster teardown S-0578 --yes
  flai kluster promote
  flai kluster down --yes`,
		Args: cobra.NoArgs,
	}
	c.AddCommand(newKlusterUpCmd(a), newKlusterDownCmd(a), newKlusterStatusCmd(a),
		newKlusterDeployCmd(a), newKlusterTeardownCmd(a), newKlusterPromoteCmd(a), newKlusterDeploymentsCmd(a), newKlusterLogCmd(a), newKlusterExecCmd(a), newKlusterEnvCmd(a))
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
	var job klusterJobFlags
	c := &cobra.Command{
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
standard error.

` + klusterJobHelp("deploy", "build, load, and deploy the story at its branch's head"),
		Example: `  flai kluster deploy S-0578
  flai kluster deploy S-0578 --json
  flai kluster deploy S-0578 --job
  flai kluster deploy S-0578 --job --wait 10m --json`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if job.job || job.wait != 0 {
				return a.askKlusterJob(cmd, kluster.KindDeploy, args[0], job)
			}
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
	addKlusterJobFlags(c, &job)
	return c
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
	var job klusterJobFlags
	c := &cobra.Command{
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
prints the record removed, or {"story": ..., "deployed": false}.

` + klusterJobHelp("teardown", "tear the story's namespace down") + ` A teardown
asked with --job needs no --yes: it asks nothing on the terminal.`,
		Example: `  flai kluster teardown S-0578
  flai kluster teardown S-0578 --yes --json
  flai kluster teardown S-0578 --job --wait 5m`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if job.job || job.wait != 0 {
				return a.askKlusterJob(cmd, kluster.KindTeardown, args[0], job)
			}
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
	addKlusterJobFlags(c, &job)
	return c
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
				fmt.Fprintf(tw, "%s\t%s\t%s\t%s\t%s\n", d.Namespace, story, short12(d.Commit), outcome, ended)
			}
			return tw.Flush()
		},
	}
}

func newKlusterLogCmd(a *app) *cobra.Command {
	var step string
	c := &cobra.Command{
		Use:   "log S-nnnn --step <step>",
		Short: "Print a step's log of a story's namespace, as the dashboard's live feed reads it",
		Long: `Prints the log of one step of the last job in a story's namespace, s-nnnn:
.flai-cache/kluster/logs/<namespace>/<step>.log, read from its start to its
end in the 64 KB parts the host method kluster.log answers the dashboard's
live feed, each ending on a whole character. A step is build, build-<name>
for a component, load, deploy, or teardown; flai kluster deployments lists
the steps each namespace's last job ran. A log a new job clears while it is
read is read again from its start.

It only reads files, so it needs no kluster: block and no flai serve, and
reads the logs of the checkout it runs in. A session with a story of its
own, FLAI_STORY, else the story in an
agent name of the form agent-S-nnnn, is refused another story's log, as
deploy --job is (ADR-0292). A step not of that form is refused, and a log
that is not there is said to be absent, exiting 1.

--json prints one object: namespace, step, log (its path from the root),
exists, offset (0), next (where the text ends), size, text (the whole log),
and outcome, the step's outcome in the namespace's record.`,
		Example: `  flai kluster log S-0578 --step deploy
  flai kluster log S-0578 --step build-api
  flai kluster log S-0578 --step deploy --json`,
		Args: cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			story := workitem.CanonicalID(args[0])
			ns, err := kluster.StoryNamespace(story)
			if err != nil {
				return err
			}
			if !kluster.IsStep(step) {
				return fmt.Errorf("%q is not a step; give --step one flai kluster deployments lists under %s's namespace, such as build, build-api, load, deploy, or teardown", step, story)
			}
			if own := sessionStory(); own != "" && workitem.CanonicalID(own) != story {
				return fmt.Errorf("this session works %s, and a story's agent reads its own story's kluster logs alone (ADR-0292), so nothing was read for %s; read %s's, or leave %s's namespace to its own agent", workitem.CanonicalID(own), story, workitem.CanonicalID(own), story)
			}
			repo, err := a.project()
			if err != nil {
				return err
			}
			l, err := readWholeLog(mainRootOf(repo), ns, step)
			if err != nil {
				return err
			}
			if a.jsonOut {
				if err := a.printJSON(l); err != nil {
					return err
				}
			} else {
				fmt.Fprint(a.out, l.Text)
			}
			if !l.Exists {
				return fmt.Errorf("%s's namespace %s has no %s log at %s: no job of it ran that step, or a teardown removed its logs; flai kluster deployments lists the steps its last job ran", story, ns, step, l.Log)
			}
			return nil
		},
	}
	c.Flags().StringVar(&step, "step", "", "the step whose log to print: build, build-<name>, load, deploy, or teardown")
	return c
}

func newKlusterExecCmd(a *app) *cobra.Command {
	return &cobra.Command{
		Use:   "exec S-nnnn -- <command> [args...]",
		Short: "Run a command in a story's worktree with its namespace's variables and kubeconfig",
		Long: `Runs a command in a story's worktree, .flai-cache/worktrees/S-nnnn of the
main checkout, with the variables of the story's namespace, s-nnnn, as its
deploy gives them (design/system/kluster.md, ADR-0254 §4): FLAI_STORY,
FLAI_KLUSTER_NAME, FLAI_KLUSTER_NAMESPACE, FLAI_KLUSTER_PREFIX,
FLAI_KLUSTER_COMMIT and FLAI_KLUSTER_TAG at the commit the namespace's record
names, FLAI_KLUSTER_STATE, and FLAI_KLUSTER_CONTEXT, but FLAI_KLUSTER_IMAGES
and FLAI_KLUSTER_COMPONENT, which only a build gets. KUBECONFIG is not flai's
kubeconfig for the cluster but the namespace's own,
.flai-cache/kluster/agents/s-nnnn.kubeconfig, which each deploy writes once
its deploy step passes, so the API server refuses the command anything
outside the namespace (ADR-0292 §2). The rest of the environment is this
shell's, with any KUBECONFIG and FLAI_KLUSTER_ variable of its own replaced
or dropped. Standard input, output, and error are the command's, and flai
exits with its exit code.

It runs while the namespace has a record under
.flai-cache/kluster/deployments, a failed one included, for diagnosis: a
deploy that failed before its kubeconfig was written is said to have none,
and the command still runs. A story with no record is refused, saying to
deploy it first with flai kluster deploy S-nnnn --job; a story with no
worktree is refused, saying to open it with flai stream open S-nnnn.

It reads the kluster: block for the cluster's name and the record for the
commit, and runs neither kluster nor flai serve, so it needs no kluster on
PATH. A session with a story of its own, FLAI_STORY, else the story in an
agent name of the form agent-S-nnnn, is refused another story's namespace,
as deploy --job is. The command goes after --, and with none it is refused.`,
		Example: `  flai kluster exec S-0578 -- kubectl get pods
  flai kluster exec S-0578 -- kubectl logs deploy/api --tail 50
  flai kluster exec S-0578 -- make integration-test
  flai kluster exec S-0578 -- sh -c 'echo $FLAI_KLUSTER_TAG'`,
		Args: func(cmd *cobra.Command, args []string) error {
			dash := cmd.ArgsLenAtDash()
			switch {
			case dash < 0 && len(args) == 0, dash == 0:
				return fmt.Errorf("give the story's ID, then -- and the command to run, such as flai kluster exec S-0578 -- kubectl get pods")
			case dash < 0 || dash == len(args):
				return fmt.Errorf("no command to run: give it after --, such as flai kluster exec %s -- kubectl get pods", workitem.CanonicalID(args[0]))
			case dash > 1:
				return fmt.Errorf("one story at a time, then -- and the command, such as flai kluster exec %s -- kubectl get pods", workitem.CanonicalID(args[0]))
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			e, err := a.klusterAgentEnv(args[0])
			if err != nil {
				return err
			}
			wt := e.repo.WorktreePath(e.story)
			if !dirExists(wt) {
				return fmt.Errorf("%s has no worktree at %s to run in; open it with flai stream open %s, then run again", e.story, wt, e.story)
			}
			a.warnNoAgentKubeconfig(e)
			return a.execIn(cmd, wt, klusterExecEnv(os.Environ(), e.vars), args[1:])
		},
	}
}

func newKlusterEnvCmd(a *app) *cobra.Command {
	return &cobra.Command{
		Use:   "env S-nnnn",
		Short: "Print a story's namespace's variables and kubeconfig, as export lines or JSON",
		Long: `Prints the variables of a story's namespace, s-nnnn, the ones flai kluster
exec runs a command with (design/system/kluster.md, ADR-0292 §2): FLAI_STORY,
FLAI_KLUSTER_NAME, FLAI_KLUSTER_NAMESPACE, FLAI_KLUSTER_PREFIX,
FLAI_KLUSTER_COMMIT and FLAI_KLUSTER_TAG at the commit the namespace's record
names, FLAI_KLUSTER_STATE, FLAI_KLUSTER_CONTEXT, and KUBECONFIG, the
namespace's own kubeconfig, .flai-cache/kluster/agents/s-nnnn.kubeconfig, not
flai's for the cluster, so the API server refuses a command run with it
anything outside the namespace. FLAI_KLUSTER_IMAGES and FLAI_KLUSTER_COMPONENT,
which only a build gets, are not among them.

It prints one line a variable, sorted by name, of the form
export NAME='value', each value quoted for sh and zsh, so that
eval "$(flai kluster env S-nnnn)" sets them in the shell, for a command of
the agent's own such as a request to a service in its namespace. flai guard
refuses a story's agent a line with $ around it, as one not written plainly,
so an agent writes them to a file and sources it on the same line,
flai kluster env S-nnnn > /tmp/s-nnnn.env && . /tmp/s-nnnn.env && <command>,
or runs the command with flai kluster exec. --json prints one object of each
variable's name and value.

It prints while the namespace has a record under
.flai-cache/kluster/deployments, a failed one included, for diagnosis: a
deploy that failed before its kubeconfig was written is said to have none on
standard error, and the variables are printed all the same. A story with no
record is refused, saying to deploy it first with flai kluster deploy S-nnnn
--job. It needs no worktree, runs neither kluster nor flai serve, and so
needs no kluster on PATH. A session with a story of its own, FLAI_STORY, else
the story in an agent name of the form agent-S-nnnn, is refused another
story's namespace, as exec is.`,
		Example: `  flai kluster env S-0578
  eval "$(flai kluster env S-0578)"
  flai kluster env S-0578 > /tmp/s-0578.env && . /tmp/s-0578.env && kubectl get pods
  flai kluster env S-0578 --json`,
		Args: cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			e, err := a.klusterAgentEnv(args[0])
			if err != nil {
				return err
			}
			a.warnNoAgentKubeconfig(e)
			vars := make(map[string]string, len(e.vars))
			for _, v := range e.vars {
				k, val, _ := strings.Cut(v, "=")
				vars[k] = val
			}
			if a.jsonOut {
				return a.printJSON(vars)
			}
			for _, k := range slices.Sorted(maps.Keys(vars)) {
				fmt.Fprintf(a.out, "export %s=%s\n", k, shellQuote(vars[k]))
			}
			return nil
		},
	}
}

// shellQuote is v in single quotes for sh and zsh, each single quote in it
// written as a quote that closes the string, a backslash-escaped quote, and
// a quote that opens it again.
func shellQuote(v string) string {
	return "'" + strings.ReplaceAll(v, "'", `'\''`) + "'"
}

// klusterAgent is a story's namespace as flai kluster exec and env give it
// to the story's agent: the story, its namespace, the project, the cluster,
// and the namespace's variables at its record's commit.
type klusterAgent struct {
	story   string
	ns      string
	repo    *workitem.Repo
	cluster kluster.Cluster
	vars    []string
}

// klusterAgentEnv reads the namespace of the story arg names and its
// variables (Cluster.AgentEnv), refusing a session with a story of its own
// another story's namespace, and a namespace with no record.
func (a *app) klusterAgentEnv(arg string) (klusterAgent, error) {
	story := workitem.CanonicalID(arg)
	ns, err := kluster.StoryNamespace(story)
	if err != nil {
		return klusterAgent{}, err
	}
	if own := sessionStory(); own != "" && workitem.CanonicalID(own) != story {
		return klusterAgent{}, fmt.Errorf("this session works %s, and a story's agent runs commands against its own story's namespace alone (ADR-0292), so nothing was run for %s; run them for %s, or leave %s's namespace to its own agent", workitem.CanonicalID(own), story, workitem.CanonicalID(own), story)
	}
	repo, err := a.project()
	if err != nil {
		return klusterAgent{}, err
	}
	c, err := kluster.New(mainRootOf(repo), repo.Manifest, kluster.System{})
	if err != nil {
		return klusterAgent{}, err
	}
	vars, err := c.AgentEnv(story)
	if errors.Is(err, kluster.ErrNotDeployed) {
		return klusterAgent{}, fmt.Errorf("%s's namespace %s has no record under .flai-cache/kluster/deployments, so there are no variables or kubeconfig to run with; deploy it first with flai kluster deploy %s --job", story, ns, story)
	}
	if err != nil {
		return klusterAgent{}, err
	}
	return klusterAgent{story: story, ns: ns, repo: repo, cluster: c, vars: vars}, nil
}

// warnNoAgentKubeconfig says on standard error that e's namespace has no
// kubeconfig of its own yet, when it has none.
func (a *app) warnNoAgentKubeconfig(e klusterAgent) {
	if kc := kluster.AgentKubeconfigPath(e.cluster.Root, e.story); !fileExists(kc) {
		fmt.Fprintf(a.errOut, "flai: %s's namespace %s has no kubeconfig at %s yet, as its deploy step has not passed; KUBECONFIG names it all the same\n", e.story, e.ns, kc)
	}
}

// klusterExecEnv is environ with vars in place of its own KUBECONFIG and
// FLAI_KLUSTER_ variables, and of any other it shares with vars.
func klusterExecEnv(environ, vars []string) []string {
	set := map[string]bool{"KUBECONFIG": true}
	for _, v := range vars {
		k, _, _ := strings.Cut(v, "=")
		set[k] = true
	}
	env := make([]string, 0, len(environ)+len(vars))
	for _, v := range environ {
		if k, _, _ := strings.Cut(v, "="); !set[k] && !strings.HasPrefix(k, "FLAI_KLUSTER_") {
			env = append(env, v)
		}
	}
	return append(env, vars...)
}

// fileExists says whether path is a file.
func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.Mode().IsRegular()
}

// execIn runs argv in dir with env, PWD dir, with the command's standard
// input, output, and error, and passes an interrupt or a termination on to it rather than
// ending flai first; a non-zero exit is an exitError of its code, 128 and the
// signal's number for a command a signal ended.
func (a *app) execIn(cmd *cobra.Command, dir string, env, argv []string) error {
	x := exec.Command(argv[0], argv[1:]...)
	x.Dir, x.Env = dir, append(slices.DeleteFunc(env, func(v string) bool { return strings.HasPrefix(v, "PWD=") }), "PWD="+dir)
	x.Stdin, x.Stdout, x.Stderr = cmd.InOrStdin(), a.out, a.errOut
	if err := x.Start(); err != nil {
		return fmt.Errorf("run %s in %s: %w", argv[0], dir, err)
	}
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(sigs)
	done := make(chan error, 1)
	go func() { done <- x.Wait() }()
	for {
		select {
		case s := <-sigs:
			_ = x.Process.Signal(s)
		case err := <-done:
			var ee *exec.ExitError
			if !errors.As(err, &ee) {
				if err != nil {
					return fmt.Errorf("run %s in %s: %w", argv[0], dir, err)
				}
				return nil
			}
			if ws, ok := ee.Sys().(syscall.WaitStatus); ok && ws.Signaled() {
				return &exitError{code: 128 + int(ws.Signal())}
			}
			return &exitError{code: ee.ExitCode()}
		}
	}
}

// readWholeLog reads namespace ns's step log under root from 0 to its end in
// kluster.ReadLog's parts, as one part holding the whole text; a log a new
// job clears while it is read is read again from 0.
func readWholeLog(root, ns, step string) (kluster.LogPart, error) {
	var text strings.Builder
	var offset int64
	for {
		part, err := kluster.ReadLog(root, ns, step, offset)
		if err != nil {
			return kluster.LogPart{}, err
		}
		if part.Reset {
			text.Reset()
		}
		text.WriteString(part.Text)
		if part.Text == "" || part.Next >= part.Size {
			part.Offset, part.Reset, part.Text = 0, false, text.String()
			return part, nil
		}
		offset = part.Next
	}
}

// klusterJobWait is how long flai kluster deploy and teardown --job wait for
// flai serve to take and answer the request, and how often --wait reads the
// job's state; zero fields are kluster's defaults, and tests shorten them.
var klusterJobWait = kluster.Wait{}

// klusterJobFlags are deploy's and teardown's --job and --wait.
type klusterJobFlags struct {
	job  bool
	wait time.Duration
}

// addKlusterJobFlags adds --job and --wait to deploy or teardown c.
func addKlusterJobFlags(c *cobra.Command, f *klusterJobFlags) {
	c.Flags().BoolVar(&f.job, "job", false, "ask flai serve for the job and print its id and the jobs ahead at once, running nothing here (the story's agent, for its own story)")
	c.Flags().DurationVar(&f.wait, "wait", 0, "with --job, wait at most this long for the job to end, such as 10m, and print its outcome")
}

// klusterJobHelp is the part of deploy's or teardown's help that says what
// --job, --wait, and --json with them do; what is the job, in words.
func klusterJobHelp(kind, what string) string {
	return fmt.Sprintf(`With --job it runs nothing here: it asks flai serve to %s
as a job (ADR-0292), and is how a story's agent asks for its own story alone:
FLAI_STORY, else the story in an agent name of the form agent-S-nnnn. A
session with a story of its own is refused another before anything is asked;
with neither, as in the operator's shell, the story given is asked for. It
writes a request under .flai-cache/kluster/requests, waits at most ten
seconds for flai serve to take it, and prints the answer: the job's id, its
kind, and the jobs ahead of it, and whether this request queued it or found
a %s job of the story already queued or running, which queues nothing
again; or, for a story in review whose deploy passed at its head, the
record's commit and outcome, with nothing queued. A refusal prints flai
serve's reason and exits 1, as while the kluster host action is off, with no
kluster.deploy set, over kluster.agent_namespaces, with uncommitted changes in
the worktree, or for a story not in progress. A request no flai serve takes
is withdrawn, and flai says flai serve is not running and exits 1.

--wait <duration> (default 0, with --job only) then waits that long for the
job to end and prints its outcome (passed, failed, refused, or interrupted),
its commit, each step with its outcome, exit code, and log, and its error;
it exits 1 unless the job passed. A job that outlasts the wait is said to be
still queued or running, and flai exits 0: read it later with flai kluster
deployments, or hold wait_for_events for its kluster event, which reports
the outcome, the failed step, its exit code, and the log. Interrupting the
wait exits 1 and leaves the job to flai serve.

--json with --job prints one object: answer, flai serve's answer (job, kind,
story, ahead, started, recorded, commit, outcome, log, and reason); and with
--wait, job, the job's state (id, kind, story, by, queued_at, started_at,
ended_at, commit, outcome, failed_step, exit_code, error, steps, and log),
and ended, whether it ended within the wait.`, what, kind)
}

// klusterJobResult is what flai kluster deploy and teardown --job print with
// --json: flai serve's answer and, with --wait, the job's state and whether
// it ended within the wait.
type klusterJobResult struct {
	Answer kluster.Answer    `json:"answer"`
	Job    *kluster.AskedJob `json:"job,omitempty"`
	Ended  *bool             `json:"ended,omitempty"`
}

// sessionStory is the story of this session, when it is a story's agent's:
// FLAI_STORY, else the story in FLAI_AGENT of the form agent-S-nnnn; "" with
// neither.
func sessionStory() string {
	return recordingStory("", os.Getenv, func() string { return "" })
}

// askKlusterJob asks flai serve for the kluster job of kind for story, in
// place of running it here (ADR-0292 §1), and prints its answer; with
// f.wait above 0, it then waits for the job to end and prints its outcome.
func (a *app) askKlusterJob(cmd *cobra.Command, kind, story string, f klusterJobFlags) error {
	switch {
	case !f.job:
		return fmt.Errorf("--wait waits for the job flai serve runs, so it goes with --job; add --job, or drop --wait")
	case f.wait < 0:
		return fmt.Errorf("--wait %s is not a time to wait; give 0 or more, such as 10m", f.wait)
	}
	story = workitem.CanonicalID(story)
	if own := sessionStory(); own != "" && workitem.CanonicalID(own) != story {
		return fmt.Errorf("this session works %s, and a story's agent asks flai serve for its own story's kluster %s alone (ADR-0292 §1), so nothing was asked for %s; ask for %s, or leave %s's namespace to its own agent", workitem.CanonicalID(own), kind, story, workitem.CanonicalID(own), story)
	}
	repo, err := a.project()
	if err != nil {
		return err
	}
	root := mainRootOf(repo)
	ctx, stop := signal.NotifyContext(cmd.Context(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	ans, err := kluster.Ask(ctx, root, kluster.Request{Story: story, Kind: kind, By: a.writer()}, klusterJobWait)
	if err != nil {
		return err
	}
	a.logger().Debug("kluster job asked", "component", "cmd", "story", ans.Story, "kind", kind, "job_id", ans.Job, "started", ans.Started, "recorded", ans.Recorded, "refused", ans.Refused())
	res := klusterJobResult{Answer: ans}
	var waitErr error
	// A story in review may be answered with the review's own deploy, which
	// keeps no job state to wait on: its record is read with flai kluster
	// deployments instead.
	if f.wait > 0 && ans.Job != "" && !ans.Refused() && !ans.Recorded && kluster.HasAskedJob(root, ans.Story, ans.Job) {
		j, ended, err := kluster.Await(ctx, root, ans.Story, ans.Job, f.wait, klusterJobWait.Poll)
		waitErr = err
		if j.ID != "" {
			res.Job, res.Ended = &j, &ended
		}
	}
	if a.jsonOut {
		if err := a.printJSON(res); err != nil {
			return err
		}
	} else {
		a.printKlusterJob(root, res, f.wait)
	}
	switch {
	case waitErr != nil:
		return waitErr
	case ans.Refused() && a.jsonOut:
		return &exitError{code: 1}
	case ans.Refused():
		return errors.New(ans.Reason)
	case ans.Recorded && ans.Outcome != kluster.OutcomePassed:
		return &exitError{code: 1}
	case res.Ended != nil && *res.Ended && res.Job.Outcome != kluster.OutcomePassed:
		return &exitError{code: 1}
	}
	return nil
}

// printKlusterJob writes flai serve's answer to a kluster request and, once
// --wait read it, the job's state, as lines; a refusal is left to the error.
func (a *app) printKlusterJob(root string, res klusterJobResult, wait time.Duration) {
	ans := res.Answer
	switch {
	case ans.Refused():
		return
	case ans.Recorded:
		fmt.Fprintf(a.out, "%s's %s is answered from its namespace's record, which %s at %s; nothing was queued\n", ans.Story, ans.Kind, ans.Outcome, short12(ans.Commit))
		if ans.Log != "" {
			fmt.Fprintf(a.out, "  log     %s\n", ans.Log)
		}
		return
	case ans.Started:
		fmt.Fprintf(a.out, "kluster %s job %s queued for %s, %s ahead of it\n", ans.Kind, ans.Job, ans.Story, jobsAhead(ans.Ahead))
	default:
		fmt.Fprintf(a.out, "kluster %s job %s of %s was already %s, so nothing more was queued; %s ahead of it\n", ans.Kind, ans.Job, ans.Story, foundState(root, ans), jobsAhead(ans.Ahead))
	}
	if res.Job == nil || !*res.Ended {
		lead := "Wait for"
		if res.Job != nil {
			lead = fmt.Sprintf("It is still %s after %s: wait for", jobState(*res.Job), wait)
		}
		fmt.Fprintf(a.out, "%s its kluster event (wait_for_events), or read it later with flai kluster deployments", lead)
		if ans.Log != "" {
			fmt.Fprintf(a.out, "; its output goes to %s", ans.Log)
		}
		fmt.Fprintln(a.out)
		return
	}
	printAskedJob(a.out, *res.Job)
}

// printAskedJob writes an ended kluster job as a short block of lines: its
// outcome, its commit, each step, and why it failed.
func printAskedJob(w io.Writer, j kluster.AskedJob) {
	fmt.Fprintf(w, "%s's %s job %s: %s\n", j.Story, j.Kind, j.ID, j.Outcome)
	if j.Commit != "" {
		fmt.Fprintf(w, "  commit  %s\n", short12(j.Commit))
	}
	for _, st := range j.Steps {
		line := fmt.Sprintf("  %-8s %s", st.Name, st.Outcome)
		if st.Outcome == kluster.OutcomeFailed {
			line += fmt.Sprintf(", exit code %d", st.ExitCode)
		}
		if st.Log != "" {
			line += ", log " + st.Log
		}
		fmt.Fprintln(w, line)
	}
	if j.FailedStep != "" {
		line := "  failed at " + j.FailedStep
		if j.ExitCode != nil {
			line += fmt.Sprintf(", exit code %d", *j.ExitCode)
		}
		fmt.Fprintln(w, line)
	}
	if j.Error != "" {
		fmt.Fprintf(w, "  error   %s\n", j.Error)
	}
	if j.Log != "" {
		fmt.Fprintf(w, "  log     %s\n", j.Log)
	}
}

// foundState says whether the job ans found of its story waits or runs, as
// its state file has it; "queued or running" when that cannot be read.
func foundState(root string, ans kluster.Answer) string {
	if j, ok, err := kluster.ReadAskedJob(root, ans.Story); err == nil && ok && j.ID == ans.Job && !j.Ended() {
		return jobState(j)
	}
	return "queued or running"
}

// jobState says in a word whether j waits or runs.
func jobState(j kluster.AskedJob) string {
	if j.Running() {
		return "running"
	}
	return "queued"
}

// jobsAhead is n jobs, in words.
func jobsAhead(n int) string {
	switch n {
	case 0:
		return "no job"
	case 1:
		return "1 job"
	}
	return fmt.Sprintf("%d jobs", n)
}

// short12 is a commit's first twelve characters, as kluster tags name it.
func short12(commit string) string {
	if len(commit) > 12 {
		return commit[:12]
	}
	return commit
}
