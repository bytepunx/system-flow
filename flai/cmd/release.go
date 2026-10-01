package cmd

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/bytepunx/system-flow/flai/internal/execx"
	"github.com/bytepunx/system-flow/flai/internal/pending"
	"github.com/bytepunx/system-flow/flai/internal/preview"
	"github.com/bytepunx/system-flow/flai/internal/release"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// planRelease computes the release for an item in the repo.
func (a *app) planRelease(repo *workitem.Repo, it *workitem.Item, deliver string) (*release.Plan, error) {
	return a.planReleaseAt(repo, it, deliver, "")
}

// planReleaseAt computes the plan from the history at root (default: the
// repository root). A dry run of acceptance plans from the story worktree,
// whose branch already holds the commits that the merge will bring to main.
func (a *app) planReleaseAt(repo *workitem.Repo, it *workitem.Item, deliver, root string) (*release.Plan, error) {
	if err := execx.Require(a.runner, "git", "Releases are git tags; install git."); err != nil {
		return nil, err
	}
	var parent *workitem.Item
	if it.Parent != "" {
		parent, _ = repo.Get(it.Parent)
	}
	if root == "" {
		root = repo.Root
	}
	return release.Compute(a.runner, root, repo.Manifest, it, parent, deliver)
}

func (a *app) printPlan(plan *release.Plan) {
	if plan.Skipped != "" {
		fmt.Fprintf(a.out, "%s %s: no release (%s)\n", plan.Item, plan.Title, plan.Skipped)
		for _, u := range plan.Unreleased {
			fmt.Fprintf(a.out, "  %-10s lands on main without a release  [%d files]\n", u.Component, len(u.Files))
		}
		return
	}
	fmt.Fprintf(a.out, "%s %s: %s release from %d commit(s)\n", plan.Item, plan.Title, plan.Level, len(plan.Commits))
	for _, st := range plan.Steps {
		kind := "incidental"
		if st.Delivered {
			kind = "delivered"
		}
		target := st.Tag
		if target == "" {
			target = st.Version + " (version file + changelog)"
		}
		fmt.Fprintf(a.out, "  %-10s %s -> %s  %-6s %s  [%d files]\n", st.Component.Name, st.From, st.To, st.Level, kind, len(st.Files))
		fmt.Fprintf(a.out, "             %s\n", target)
	}
	if len(plan.Untouched) > 0 {
		fmt.Fprintf(a.out, "  %d touched files outside any component (design, docs, wip)\n", len(plan.Untouched))
	}
}

func newReleaseCmd(a *app) *cobra.Command {
	var deliver string
	var dryRun, onlyPending bool
	c := &cobra.Command{
		Use:   "release <id> | --pending",
		Short: "Compute a release for one item, or publish everything accepted since the last release",
		Long: `Per design/conventions/git.md: the component the item delivers to gets the
delivery-type bump (epic major, feature story minor, remediation or
improvement patch); every other component its commits touched gets a patch.
Components come from system-flow.yaml projects; the delivered one from the
item's tags (a project name or one of its tags), the parent's tags, or
--deliver. Code components get an annotated tag <name>/vX.Y.Z on HEAD; the
template component gets its version file and changelog bumped (commit them).

flai accept never does this (S-0087): it only merges, archives, and commits.

flai release --pending computes one release per component, the highest
delivery level among everything accepted and unreleased for it since its last
tag, bumps and commits, tags, and pushes the branch and every tag together,
three tags to a push (I-0026). Run again after a partial failure: what already tagged or
pushed is not redone. Tags on commits the remote already has, such as an
acceptance pushed before it was released, are pushed on their own. flai push
--pending does the same computing, applying, and tagging before it pushes
only with the auto-publish host action enabled (S-0144); otherwise this is
where what has accumulated is released, when you choose.`,
		Example: `  flai release S-031 --dry-run
  flai release S-031 --deliver flai --apply
  flai release --pending --dry-run
  flai release --pending`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if onlyPending {
				if len(args) > 0 {
					return fmt.Errorf("--pending publishes everything accumulated, not one item; drop the id")
				}
				return a.publishPending(dryRun)
			}
			if len(args) != 1 {
				return fmt.Errorf("an item ID is required, or --pending to publish everything accumulated")
			}
			repo, err := a.project()
			if err != nil {
				return err
			}
			it, err := repo.Get(args[0])
			if err != nil {
				return err
			}
			plan, err := a.planRelease(repo, it, deliver)
			if err != nil {
				return err
			}
			apply, _ := cmd.Flags().GetBool("apply")
			if a.jsonOut && (dryRun || !apply) {
				return a.printJSON(plan)
			}
			a.printPlan(plan)
			if dryRun || !apply || plan.Skipped != "" {
				if !apply && !dryRun && plan.Skipped == "" {
					fmt.Fprintln(a.out, "(plan only; add --apply to create tags and bump version files, or flai release --pending to publish it with everything else accumulated)")
				}
				return nil
			}
			if err := release.Apply(a.runner, repo.Root, plan, a.now()); err != nil {
				return err
			}
			tags, err := release.Tag(a.runner, repo.Root, plan)
			if err != nil {
				return err
			}
			if a.jsonOut {
				return a.printJSON(map[string]any{"plan": plan, "tags": tags})
			}
			if len(tags) > 0 {
				fmt.Fprintf(a.out, "tagged %s on HEAD (push with git push origin %s)\n", strings.Join(tags, ", "), strings.Join(tags, " "))
			}
			return nil
		},
	}
	c.Flags().StringVar(&deliver, "deliver", "", "component the item delivers to, when its tags do not say")
	c.Flags().BoolVar(&dryRun, "dry-run", false, "print the plan only, or with --pending what would publish")
	c.Flags().Bool("apply", false, "create tags and bump version files")
	c.Flags().BoolVar(&onlyPending, "pending", false, "publish everything merged and unreleased since each component's last tag")
	return c
}

// computeApplyAndTagPending is the first half of flai release --pending
// (S-0087) and, since S-0094, of flai push --pending too: everything
// release.Pending finds since each component's last tag, applied (version
// files and one changelog entry, committed) and tagged on HEAD, before
// anything is pushed. flai accept computes none of this (it only merges,
// archives, and commits), and flai push does only with the auto-publish host
// action enabled (S-0144). TagPending is idempotent, so a
// partial failure here and a rerun does not retag what already tagged.
// Nothing is applied while the remote has a newer release tag than this
// clone, or cannot be asked whether it has (S-0174): exit 3, as for a remote
// that moved.
func (a *app) computeApplyAndTagPending(root string, repo *workitem.Repo) ([]*release.PendingPlan, []string, error) {
	if err := a.remoteTagsInStep(root, repo); err != nil {
		return nil, nil, err
	}
	plans, err := release.Pending(a.runner, root, repo.Manifest, repo)
	if err != nil {
		return nil, nil, err
	}
	for _, p := range plans {
		if err := release.ApplyPending(p, root, a.now()); err != nil {
			return nil, nil, err
		}
	}
	if len(plans) > 0 {
		if status, _ := a.runner.Run(root, "git", "status", "--porcelain"); strings.TrimSpace(status) != "" {
			if _, err := a.runner.Run(root, "git", "add", "-A"); err != nil {
				return nil, nil, err
			}
			if _, err := a.runner.Run(root, "git", "commit", "-q", "-m", "chore: publish "+summarizePlans(plans)); err != nil {
				return nil, nil, err
			}
		}
	}
	var tags []string
	for _, p := range plans {
		tag, err := release.TagPending(a.runner, root, p)
		if err != nil {
			return nil, nil, err
		}
		if tag != "" {
			tags = append(tags, tag)
		}
	}
	return plans, tags, nil
}

// remoteTagsInStep refuses a publish from a clone missing release tags its
// remote has, which would tag versions already published, and from one that
// cannot ask its remote, whose push would fail after tagging (S-0174). A
// clone with no remote publishes locally as before.
func (a *app) remoteTagsInStep(root string, repo *workitem.Repo) error {
	if refusal := release.CheckRemote(a.runner, root, repo.Manifest).Refusal(); refusal != "" {
		return &exitError{code: exitPushDiverged, msg: refusal}
	}
	return nil
}

// publishPending is flai release --pending (S-0087): everything
// computeApplyAndTagPending finds, applied and tagged, pushed together,
// resumable on partial failure.
func (a *app) publishPending(dryRun bool) error {
	if err := execx.Require(a.runner, "git", "Releases are git tags; install git."); err != nil {
		return err
	}
	repo, err := a.project()
	if err != nil {
		return err
	}
	if dryRun {
		pub, err := preview.Publish(a.runner, repo)
		if err != nil {
			return err
		}
		if a.jsonOut {
			return a.printJSON(pub)
		}
		if pub.Remote.Lagging() {
			fmt.Fprintf(a.out, "nothing can publish: %s\n", pub.Remote.Message)
			return nil
		}
		if pub.Remote != nil {
			fmt.Fprintf(a.out, "warning: %s\n", pub.Remote.Message)
		}
		if len(pub.Plans) == 0 {
			fmt.Fprintln(a.out, "nothing pending")
			return nil
		}
		for _, p := range pub.Plans {
			a.printPendingPlan(p)
		}
		fmt.Fprintln(a.out, "dry run: nothing changed")
		return nil
	}
	plans, tags, err := a.computeApplyAndTagPending(repo.Root, repo)
	if err != nil {
		return err
	}
	u := pending.Detect(a.runner, repo.Root)
	if u == nil {
		u = pending.TagsOnly(a.runner, repo.Root, tags)
	}
	result := map[string]any{"plans": plans, "tags": tags, "pushed": false}
	if u == nil {
		if len(plans) == 0 {
			result["reason"] = "nothing pending"
			if a.jsonOut {
				return a.printJSON(result)
			}
			fmt.Fprintln(a.out, "nothing pending")
			return nil
		}
		// tagged and committed locally, but nothing ahead of a remote to push
		// to (no upstream configured): report what was done, not an error.
		if a.jsonOut {
			return a.printJSON(result)
		}
		fmt.Fprintf(a.out, "published locally: %s (no upstream to push to)\n", summarizePlans(plans))
		return nil
	}
	for _, batch := range u.Pushes() {
		if _, err := a.runner.Run(repo.Root, "git", append([]string{"push", "-q", u.Remote}, batch...)...); err != nil {
			result["push_error"] = firstLine(err.Error())
			if a.jsonOut {
				return a.printJSON(result)
			}
			return fmt.Errorf("published and committed locally, but the push failed and nothing was forced: %s. Fetch and merge if the remote moved, then run flai release --pending again; what already tagged is not redone", firstLine(err.Error()))
		}
	}
	result["pushed"] = true
	var published []string
	for _, p := range plans {
		if p.Component.Kind != "template" {
			continue
		}
		pub, err := a.publishTemplate(filepath.Join(repo.Root, p.Component.Path), "", "", true, false, false)
		if err != nil {
			return fmt.Errorf("published and pushed %s, but publishing %s to its own remote failed: %w. Run flai template push %s --tag when it is put right", summarizePlans(plans), p.Component.Path, err, p.Component.Path)
		}
		if !pub.Nothing {
			published = append(published, fmt.Sprintf("%s %s (%s)", pub.Remote, pub.Tag, pub.Commit))
		}
	}
	result["published"] = published
	if a.jsonOut {
		return a.printJSON(result)
	}
	fmt.Fprintf(a.out, "published %s, pushed to %s", summarizePlans(plans), u.Remote)
	if len(published) > 0 {
		fmt.Fprintf(a.out, "; published %s", strings.Join(published, ", "))
	}
	fmt.Fprintln(a.out)
	return nil
}

func summarizePlans(plans []*release.PendingPlan) string {
	parts := make([]string, len(plans))
	for i, p := range plans {
		parts[i] = fmt.Sprintf("%s %s", p.Component.Name, p.To)
	}
	return strings.Join(parts, ", ")
}

func (a *app) printPendingPlan(p *release.PendingPlan) {
	fmt.Fprintf(a.out, "%-10s %s -> %s  %-6s [%d item(s), %d file(s)]\n", p.Component.Name, p.From, p.To, p.Level, len(p.Items), len(p.Files))
	for _, it := range p.Items {
		fmt.Fprintf(a.out, "  %-10s %-6s %s\n", it.ID, it.Level, it.Title)
	}
}
