package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/bytepunx/system-flow/flai/internal/pending"
)

func newPushCmd(a *app) *cobra.Command {
	var onlyPending, dryRun bool
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
				return fmt.Errorf("%s and %s have diverged: %d commit(s) here and %d there. Fetch and merge first (git fetch %s && git merge %s), then run this again; flai never forces a push", u.Branch, u.Upstream, u.Commits, u.Behind, u.Remote, u.Upstream)
			}
			if reason, ok := result["reason"]; ok {
				if a.jsonOut {
					return a.printJSON(result)
				}
				fmt.Fprintln(a.out, reason)
				return nil
			}
			result["unpushed"] = u
			what := fmt.Sprintf("%s to %s: %s%s", strings.Join(u.Acceptances, ", "), u.Remote, u.Branch, tagsNote(u.Tags))
			if dryRun {
				result["dry_run"] = true
				if a.jsonOut {
					return a.printJSON(result)
				}
				fmt.Fprintf(a.out, "would push %s\ndry run: nothing pushed\n", what)
				return nil
			}
			if _, err := a.runner.Run(root, "git", append([]string{"push", u.Remote}, u.Refs()...)...); err != nil {
				return fmt.Errorf("the push failed and nothing was forced: %s. If the remote moved, fetch and merge, then run this again", firstLine(err.Error()))
			}
			result["pushed"] = true
			if a.jsonOut {
				return a.printJSON(result)
			}
			fmt.Fprintf(a.out, "pushed %s\n", what)
			return nil
		},
	}
	c.Flags().BoolVar(&onlyPending, "pending", false, "push the branch and tags of acceptances that were made and not pushed")
	c.Flags().BoolVar(&dryRun, "dry-run", false, "say what would be pushed")
	return c
}
