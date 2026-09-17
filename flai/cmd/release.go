package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/bytepunx/system-flow/flai/internal/execx"
	"github.com/bytepunx/system-flow/flai/internal/release"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// planRelease computes the release for an item in the repo.
func (a *app) planRelease(repo *workitem.Repo, it *workitem.Item, deliver string) (*release.Plan, error) {
	if err := execx.Require(a.runner, "git", "Releases are git tags; install git."); err != nil {
		return nil, err
	}
	var parent *workitem.Item
	if it.Parent != "" {
		parent, _ = repo.Get(it.Parent)
	}
	return release.Compute(a.runner, repo.Root, repo.Manifest, it, parent, deliver)
}

func (a *app) printPlan(plan *release.Plan) {
	if plan.Skipped != "" {
		fmt.Fprintf(a.out, "%s %s: no release (%s)\n", plan.Item, plan.Title, plan.Skipped)
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
	var dryRun bool
	c := &cobra.Command{
		Use:   "release <id>",
		Short: "Compute (and with --apply create) the semver release for an accepted item",
		Long: `Per design/conventions/git.md: the component the item delivers to gets the
delivery-type bump (epic major, feature story minor, remediation or
improvement patch); every other component its commits touched gets a patch.
Components come from system-flow.yaml projects; the delivered one from the
item's tags (a project name or one of its tags), the parent's tags, or
--deliver. Code components get an annotated tag <name>/vX.Y.Z on HEAD; the
template component gets its version file and changelog bumped (commit them).

flai accept does all of this after moving the item to done.`,
		Example: `  flai release S-031 --dry-run
  flai release S-031 --deliver flai --apply`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
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
					fmt.Fprintln(a.out, "(plan only; add --apply to create tags and bump version files, or use flai accept)")
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
	c.Flags().BoolVar(&dryRun, "dry-run", false, "print the plan only")
	c.Flags().Bool("apply", false, "create tags and bump version files")
	return c
}
