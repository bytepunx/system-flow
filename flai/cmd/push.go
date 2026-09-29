package cmd

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/bytepunx/system-flow/flai/internal/hostapi"
	"github.com/bytepunx/system-flow/flai/internal/preview"
	"github.com/bytepunx/system-flow/flai/internal/release"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

func newPushCmd(a *app) *cobra.Command {
	var onlyPending, dryRun, publishToo bool
	c := &cobra.Command{
		Use:   "push --pending",
		Short: "Push an acceptance that was made and not pushed",
		Long: `An acceptance made where nothing could push it, such as one from the
dashboard with the push host action off, is committed in the main checkout
and not pushed. Run this on the host, with your own credentials: when the
main checkout's branch is ahead of its remote-tracking branch and the
commits ahead include an acceptance or a release tagged here, it pushes the
branch and the tags together. It never forces. When the remote has commits
this clone lacks it refuses and says to fetch and merge first.

It releases nothing of its own accord: what is accepted waits, unreleased,
to be published together by flai release --pending or the board's Publish
(S-0144, ADR-0032). With the auto-publish host action enabled for the project
(flai serve enable auto-publish), it first computes the release of everything
accepted and unreleased since each component's last tag, applies the
version bump, commits it, and tags it, so every push publishes (S-0094).

--publish also publishes each template component whose version those
commits moved, as flai template push --tag does, after the push and never
forced. It is what the dashboard's push action runs (flai serve enable push).

flai board and the dashboard say when there is something pending; agents are
told by the MCP inbox.`,
		Example: `  flai push --pending
  flai push --pending --dry-run`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if !onlyPending {
				return fmt.Errorf("flai push pushes pending acceptances only; say so with --pending, or use git push for anything else")
			}
			repo, err := a.project()
			if err != nil {
				return err
			}
			root := mainRootOf(repo)

			// Tag whatever release has accumulated before deciding what to
			// push (S-0094), only when the operator enabled auto-publish: off, the
			// accepted work waits to be released together (S-0144). A dry run
			// only previews the plan; nothing is applied or tagged.
			cfg, _, err := a.loadConfig()
			if err != nil {
				return err
			}
			autoPublish := cfg.ActionEnabled(hostapi.ActionAutoPublish, root)
			var p *preview.Pushing
			var newTags []string
			switch {
			case dryRun:
				p, err = preview.PushDryRun(a.runner, root, repo, autoPublish)
			case autoPublish:
				var plans []*release.PendingPlan
				if plans, newTags, err = a.computeApplyAndTagPending(root, repo); err == nil {
					p, err = preview.Push(a.runner, root, plans, newTags)
				}
			default:
				p, err = preview.Push(a.runner, root, nil, nil)
			}
			if preview.IsDiverged(err) {
				// exit 3: a caller that is not a person (the dashboard's push
				// action) tells a remote that moved from a failure of flai's
				return &exitError{code: exitPushDiverged, msg: err.Error()}
			}
			if err != nil {
				return err
			}
			result, plans, u := p.Result, p.Plans, p.Unpushed
			if u == nil {
				if a.jsonOut {
					return a.printJSON(result)
				}
				if reason, ok := result["reason"]; ok {
					fmt.Fprintln(a.out, reason)
					return nil
				}
				// tagged and committed locally, but nothing ahead of a remote
				// to push to (no upstream configured): say what was done.
				fmt.Fprintf(a.out, "tagged locally: %s (no upstream to push to)\n", strings.Join(newTags, ", "))
				return nil
			}
			// which templates those commits moved is asked before the push:
			// afterwards nothing is ahead and nothing says so
			var moved []string
			if publishToo {
				moved = a.templatesMoved(repo, root, u.Upstream)
			}
			what := fmt.Sprintf("%s to %s: %s%s", strings.Join(u.Acceptances, ", "), u.Remote, u.Branch, tagsNote(u.Tags))
			if u.Commits == 0 {
				what = fmt.Sprintf("tags %s to %s", strings.Join(u.Tags, ", "), u.Remote)
			}
			if len(moved) > 0 {
				what += ", then publish " + strings.Join(moved, ", ")
			}
			if dryRun {
				if a.jsonOut {
					return a.printJSON(result)
				}
				if len(plans) > 0 {
					fmt.Fprintln(a.out, "would also tag, first:")
					for _, pl := range plans {
						a.printPendingPlan(pl)
					}
				}
				fmt.Fprintf(a.out, "would push %s\ndry run: nothing pushed\n", what)
				return nil
			}
			for _, refs := range u.Pushes() {
				if _, err := a.runner.Run(root, "git", append([]string{"push", u.Remote}, refs...)...); err != nil {
					return fmt.Errorf("the push failed and nothing was forced: %s. If the remote moved, fetch and merge, then run this again", firstLine(err.Error()))
				}
			}
			result["pushed"] = true
			result["tags"] = newTags
			published := []string{}
			for _, path := range moved {
				pub, err := a.publishTemplate(filepath.Join(root, path), "", "", true, false, false)
				if err != nil {
					return fmt.Errorf("pushed %s, but publishing %s failed and nothing was forced: %s. Run flai template push %s --tag when it is put right", what, path, firstLine(err.Error()), path)
				}
				if !pub.Nothing {
					published = append(published, fmt.Sprintf("%s %s (%s)", pub.Remote, pub.Tag, pub.Commit))
				}
			}
			result["published"] = published
			if len(published) > 0 {
				what += "; published " + strings.Join(published, ", ")
			}
			if a.jsonOut {
				return a.printJSON(result)
			}
			fmt.Fprintf(a.out, "pushed %s\n", what)
			return nil
		},
	}
	c.Flags().BoolVar(&onlyPending, "pending", false, "push the branch and tags of acceptances that were made and not pushed")
	c.Flags().BoolVar(&dryRun, "dry-run", false, "say what would be pushed")
	c.Flags().BoolVar(&publishToo, "publish", false, "also publish template components whose version the pushed commits moved")
	return c
}

// exitPushDiverged is flai push's exit code when the remote has commits this
// clone lacks: nothing was pushed, and fetching and merging is the way on.
const exitPushDiverged = 3

// templatesMoved lists the template components whose template.yaml differs
// between the remote-tracking branch and the checkout: the ones an unpushed
// acceptance released.
func (a *app) templatesMoved(repo *workitem.Repo, root, upstream string) []string {
	var moved []string
	for _, p := range repo.Manifest.Projects {
		if p.Kind != "template" {
			continue
		}
		file := filepath.ToSlash(filepath.Join(p.Path, "template.yaml"))
		if _, err := a.runner.Run(root, "git", "diff", "--quiet", upstream, "HEAD", "--", file); err != nil {
			moved = append(moved, p.Path)
		}
	}
	return moved
}
