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
delivery-type bump (feature story minor, remediation or improvement patch;
an epic none of its own, its stories carry theirs, ADR-0078); every other
component its commits touched gets a patch.
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
acceptance pushed before it was released, are pushed on their own. This is
the one way accepted work reaches the remote, when you choose (ADR-0067);
flai push --pending, and its auto-publish host action (S-0144), are an
operator's tools outside the workflow.

What is pending is worked out from this clone's tags and branch, and flai
never fetches: publishing is git fetch, then flai release --pending. Before
planning, --pending asks the remote the branch tracks (else origin) for its
release tags (S-0174) and its head of that branch (ADR-0067). When it has a
newer <name>/vX.Y.Z than this clone, or commits on the branch this clone
lacks, nothing is planned, applied, committed, or tagged, and publishing is
refused (exit 3), naming what to run: git fetch --tags for the tags; for the
branch, git fetch, then git merge <remote>/<branch> or git rebase onto it,
then flai release --pending again. When the branch moves after that check
and the push is refused, the release tags the push did not send are deleted
here, since they no longer tag what will be published, and the exit is 3:
git fetch, rebase onto the remote branch or merge it (merge when some tags
already went), verify, and flai release --pending again tags again
(S-0242). When the remote cannot be reached,
--dry-run warns and shows the plan, and publishing is refused. A clone with
no remote publishes locally. An accepted item no plan can
cover, such as one touching two components with no tag saying which it
delivers to, is named with the reason (I-0024).`,
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
// clone (S-0174) or commits on its branch this clone lacks (ADR-0067), or
// cannot be asked whether it has: exit 3, as for a remote that moved.
func (a *app) computeApplyAndTagPending(root string, repo *workitem.Repo) ([]*release.PendingPlan, []string, error) {
	if err := a.remoteTagsInStep(root, repo); err != nil {
		return nil, nil, err
	}
	b, err := release.PendingBatch(a.runner, root, repo.Manifest, repo)
	if err != nil {
		return nil, nil, err
	}
	for _, u := range b.Unplanned {
		a.logger().Warn("accepted item left out of the release", "component", "release", "item", u.ID, "detail", fmt.Sprintf("%s is left out of the release: %s", u.ID, u.Reason))
	}
	plans := b.Plans
	for _, p := range plans {
		if err := release.ApplyPending(p, root, a.now()); err != nil {
			return nil, nil, err
		}
		raised, err := release.RaiseMinimum(a.runner, root, p)
		if err != nil {
			return nil, nil, err
		}
		if raised {
			a.logger().Warn("flai.minimum raised: this release changes the front-matter fields flai reads, so a flai below it stops on this project; upgrade the host's flai once the release's binaries are built", "component", "release", "minimum", p.To.String())
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
// remote has, which would tag versions already published (S-0174), from one
// missing commits on its remote branch, whose push would fail after tagging
// (ADR-0067), and from one that cannot ask its remote. A clone with no
// remote publishes locally as before.
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
		for _, u := range pub.Unplanned {
			fmt.Fprintf(a.out, "left out: %s %s: %s\n", u.ID, u.Title, u.Reason)
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
	pushes := u.Pushes()
	for i, batch := range pushes {
		// --atomic: a batch whose branch the remote refuses takes its tags
		// with it, where a plain push would leave them on the remote tagging
		// commits it never got (S-0242).
		if _, err := a.runner.Run(repo.Root, "git", append([]string{"push", "-q", "--atomic", u.Remote}, batch...)...); err != nil {
			result["push_error"] = firstLine(err.Error())
			// Only a push carrying commits can be refused because the
			// remote's branch moved; tags alone tag what it already has.
			if u.Commits > 0 {
				if now := release.CheckRemoteAfresh(a.runner, repo.Root, repo.Manifest); now != nil && now.Branch != nil {
					return a.remoteMovedUnderPush(repo, u.Remote, now.Branch, pushes[:i], pushes[i:], result)
				}
			}
			if a.jsonOut {
				return a.printJSON(result)
			}
			return fmt.Errorf("published and committed locally, but the push failed and nothing was forced: %s. Put that right, then run flai release --pending again; what already tagged is not redone", firstLine(err.Error()))
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

// remoteMovedUnderPush is a publish whose push the remote refused because
// its branch moved between the check before tagging and the push (S-0242,
// ADR-0067, TH-0070). The release tags of the batch refused and of every
// batch after it never reached the remote and tag commits that a rebase or
// merge will replace as what is published, so they are deleted here; the
// next flai release --pending tags again. Tags of the batches before it did
// reach the remote and are kept, which is why the advice is then to merge.
// Exit 3, as for the remote-moved refusal.
func (a *app) remoteMovedUnderPush(repo *workitem.Repo, remote string, moved *release.BranchLag, sent, unsent [][]string, result map[string]any) error {
	kept, deleted := []string{}, []string{}
	for _, batch := range sent {
		kept = append(kept, batch...)
	}
	for _, batch := range unsent {
		for _, tag := range batch {
			if !release.IsReleaseTag(repo.Manifest, tag) {
				continue
			}
			if _, err := a.runner.Run(repo.Root, "git", "tag", "-d", tag); err != nil {
				return fmt.Errorf("the push was refused because %s moved (at %s), and deleting the release tag %s here failed: %s. Delete it with git tag -d %s, then fetch, rebase onto %s or merge it, verify, and run flai release --pending again", moved.Upstream, short(moved.Head), tag, firstLine(err.Error()), tag, moved.Upstream)
			}
			deleted = append(deleted, tag)
		}
	}
	result["remote_moved"] = true
	result["deleted_tags"] = deleted
	result["kept_tags"] = kept
	if a.jsonOut {
		if err := a.printJSON(result); err != nil {
			return err
		}
		return &exitError{code: exitPushDiverged}
	}
	var msg strings.Builder
	fmt.Fprintf(&msg, "conflict: the push was refused because %s moved (now at %s)", moved.Upstream, short(moved.Head))
	if len(deleted) > 0 {
		fmt.Fprintf(&msg, "; deleted the release tags %s here, because they no longer tag what will be published", strings.Join(deleted, ", "))
	}
	if len(kept) > 0 {
		fmt.Fprintf(&msg, ". Tags %s already reached %s and are kept, here and on %s. Run git fetch %s, then merge %s rather than rebase onto it, so the commits they tag stay in history", strings.Join(kept, ", "), remote, remote, remote, moved.Upstream)
	} else {
		fmt.Fprintf(&msg, ". Run git fetch %s, then rebase onto %s or merge it", remote, moved.Upstream)
	}
	msg.WriteString("; resolve any conflicts, verify, and run flai release --pending again, which tags again")
	return &exitError{code: exitPushDiverged, msg: msg.String()}
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
