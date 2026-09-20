package cmd

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/bytepunx/system-flow/flai/internal/pending"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

func newPushCmd(a *app) *cobra.Command {
	var onlyPending, dryRun, publishToo bool
	c := &cobra.Command{
		Use:   "push --pending",
		Short: "Push an acceptance that was made and not pushed",
		Long: `An acceptance made where there is no git credential, such as the dashboard
container without a push key, is committed and tagged in the clone and not
pushed. Run this on the host, with your own credentials: when the main
checkout's branch is ahead of its remote-tracking branch and the commits
ahead include an acceptance, it pushes the branch and the tags on those
commits. It never forces. When the remote has commits this clone lacks it
refuses and says to fetch and merge first.

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
			u := pending.Detect(a.runner, root)
			result := map[string]any{"pushed": false}
			switch {
			case u == nil:
				result["reason"] = "nothing pending"
			case !u.Pending():
				result["reason"] = fmt.Sprintf("%s is ahead of %s by %d commit(s), none of them an acceptance; that is yours to push with git", u.Branch, u.Upstream, u.Commits)
			case u.Behind > 0:
				// exit 3: a caller that is not a person (the dashboard's push
				// action) tells a remote that moved from a failure of flai's
				return &exitError{code: exitPushDiverged, msg: fmt.Sprintf("conflict: %s and %s have diverged: %d commit(s) here and %d there. Fetch and merge first (git fetch %s && git merge %s), then run this again; flai never forces a push", u.Branch, u.Upstream, u.Commits, u.Behind, u.Remote, u.Upstream)}
			}
			if reason, ok := result["reason"]; ok {
				if a.jsonOut {
					return a.printJSON(result)
				}
				fmt.Fprintln(a.out, reason)
				return nil
			}
			result["unpushed"] = u
			// which templates those commits moved is asked before the push:
			// afterwards nothing is ahead and nothing says so
			var moved []string
			if publishToo {
				moved = a.templatesMoved(repo, root, u.Upstream)
			}
			what := fmt.Sprintf("%s to %s: %s%s", strings.Join(u.Acceptances, ", "), u.Remote, u.Branch, tagsNote(u.Tags))
			if len(moved) > 0 {
				what += ", then publish " + strings.Join(moved, ", ")
			}
			if dryRun {
				result["dry_run"] = true
				if a.jsonOut {
					return a.printJSON(result)
				}
				fmt.Fprintf(a.out, "would push %s\ndry run: nothing pushed\n", what)
				return nil
			}
			for _, refs := range pending.Batches(u.Branch, u.Tags) {
				if _, err := a.runner.Run(root, "git", append([]string{"push", u.Remote}, refs...)...); err != nil {
					return fmt.Errorf("the push failed and nothing was forced: %s. If the remote moved, fetch and merge, then run this again", firstLine(err.Error()))
				}
			}
			result["pushed"] = true
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
