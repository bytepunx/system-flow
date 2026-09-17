package cmd

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/bytepunx/system-flow/flai/internal/execx"
	"github.com/bytepunx/system-flow/flai/internal/release"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

func newAcceptCmd(a *app) *cobra.Command {
	var by, deliver, remote string
	var trailers []string
	var noRelease, noPush, noPublish, dryRun bool
	c := &cobra.Command{
		Use:   "accept <id>",
		Short: "Accept an item: move to done, archive, commit, release, push",
		Long: `The operator's acceptance step as one command, per
design/conventions/work-management.md and git.md:

  1. move the item to done (its rules apply: children closed, criteria checked)
  2. flai archive for the item and its children and narrative
  3. compute the release; bump the template version file and changelog if the
     template is a component in the plan
  4. git commit the work item, archive, and version changes
  5. create the release tags on that commit
  6. push the branch and the tags

--dry-run shows the release plan and stops before step 1.`,
		Example: `  flai accept S-031 --by alex
  flai accept E-002 --by alex --deliver flai
  flai accept S-016 --by alex --no-release      # docs-only story
  flai accept S-031 --by alex --no-push`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			repo, err := a.project()
			if err != nil {
				return err
			}
			if err := execx.Require(a.runner, "git", "Acceptance commits and tags; install git."); err != nil {
				return err
			}
			it, err := repo.Get(args[0])
			if err != nil {
				return err
			}
			if it.Type == workitem.Task {
				return fmt.Errorf("%s is a task; accept its story instead", it.ID)
			}
			if it.Closed() {
				return fmt.Errorf("%s is already %s", it.ID, it.Status)
			}
			if it.Type == workitem.Story && it.Status != workitem.Review {
				return fmt.Errorf("%s is %s; a story is accepted from review", it.ID, it.Status)
			}
			if dirty := a.dirtyOutsideWip(repo); len(dirty) > 0 && !a.yes {
				return fmt.Errorf("working tree has uncommitted changes outside wip (%s); commit or stash them so the acceptance commit holds only acceptance, or pass --yes to include them", strings.Join(dirty, ", "))
			}
			// 0. bring the story branch into the main branch (ADR-0019)
			merged := false
			if it.Type == workitem.Story {
				if dryRun {
					if a.branchExists(repo.MainRoot, storyBranch(it.ID)) {
						fmt.Fprintf(a.out, "would merge %s into the main branch and remove its worktree\n", storyBranch(it.ID))
					}
				} else if merged, err = a.mergeStoryBranch(repo, it.ID); err != nil {
					return err
				}
			}
			var plan *release.Plan
			if !noRelease {
				plan, err = a.planRelease(repo, it, deliver)
				if err != nil {
					return err
				}
				a.printPlan(plan)
			}
			if dryRun {
				fmt.Fprintln(a.out, "dry run: nothing changed")
				return nil
			}

			// 1. done (epics: walk through review if needed)
			items, err := repo.List(false)
			if err != nil {
				return err
			}
			board, err := repo.LoadBoard()
			if err != nil {
				return err
			}
			// epics are accepted from wherever they are: walk the state machine to done
			path := []string{workitem.Ready, workitem.InProgress, workitem.Review, workitem.Done}
			steps := []string{workitem.Done}
			for i, st := range path {
				if st == it.Status {
					steps = path[i+1:]
				}
			}
			if it.Status == workitem.Backlog {
				steps = path
			}
			for _, st := range steps {
				if _, err := repo.Move(it, st, workitem.MoveOptions{By: orDefault(by, a.author()), Now: a.now(), Items: items, Board: board}); err != nil {
					return err
				}
			}
			if err := repo.Save(it); err != nil {
				return err
			}
			if it.Type == workitem.Story {
				if err := board.Save(a.now().Format("2006-01-02")); err != nil {
					return err
				}
			}
			// 2. archive
			items, _ = repo.List(false)
			ap, err := repo.PlanArchive(items, []string{it.ID})
			if err != nil {
				return err
			}
			if err := repo.Archive(ap); err != nil {
				return err
			}
			if err := a.refreshIndex(repo); err != nil {
				return err
			}
			// 3. version files
			if plan != nil && plan.Skipped == "" {
				if err := release.Apply(a.runner, repo.Root, plan, a.now()); err != nil {
					return err
				}
			}
			// 4. commit
			if _, err := a.runner.Run(repo.Root, "git", "add", "-A"); err != nil {
				return err
			}
			msg := fmt.Sprintf("chore: [%s] accept and archive", it.ID)
			if plan != nil && plan.Skipped == "" {
				var parts []string
				for _, st := range plan.Steps {
					parts = append(parts, fmt.Sprintf("%s %s", st.Component.Name, st.To))
				}
				msg += "; release " + strings.Join(parts, ", ")
			}
			for _, t := range trailers {
				msg += "\n\n" + t
			}
			if _, err := a.runner.Run(repo.Root, "git", "commit", "-q", "-m", msg); err != nil {
				return err
			}
			// 5. tags
			var tags []string
			if plan != nil && plan.Skipped == "" {
				if tags, err = release.Tag(a.runner, repo.Root, plan); err != nil {
					return err
				}
			}
			// 6. push
			pushed := false
			if !noPush {
				if _, err := a.runner.Run(repo.Root, "git", "remote", "get-url", remote); err == nil {
					if _, err := a.runner.Run(repo.Root, "git", "push", "-q", remote, "HEAD"); err != nil {
						return err
					}
					for _, t := range tags {
						if _, err := a.runner.Run(repo.Root, "git", "push", "-q", remote, t); err != nil {
							return err
						}
					}
					pushed = true
				}
			}
			// 7. publish template components at their new version
			var published []string
			if plan != nil && plan.Skipped == "" && !noPublish && !noPush {
				for _, st := range plan.Steps {
					if st.Component.Kind != "template" {
						continue
					}
					dir := filepath.Join(repo.Root, st.Component.Path)
					res, err := a.publishTemplate(dir, "", "", true, false, false)
					if err != nil {
						return fmt.Errorf("template published locally at %s but pushing to its remote failed: %w", st.To, err)
					}
					if !res.Nothing {
						published = append(published, fmt.Sprintf("%s %s (%s)", res.Remote, res.Tag, res.Commit))
					}
				}
			}
			if a.jsonOut {
				return a.printJSON(map[string]any{"id": it.ID, "archived": len(ap.Items), "merged": merged, "plan": plan, "tags": tags, "pushed": pushed, "published": published})
			}
			fmt.Fprintf(a.out, "accepted %s: done, %d items archived, committed", it.ID, len(ap.Items))
			if merged {
				fmt.Fprintf(a.out, ", %s merged and removed", storyBranch(it.ID))
			}
			if len(tags) > 0 {
				fmt.Fprintf(a.out, ", tagged %s", strings.Join(tags, ", "))
			}
			if pushed {
				fmt.Fprint(a.out, ", pushed")
			}
			for _, p := range published {
				fmt.Fprintf(a.out, ", template published to %s", p)
			}
			fmt.Fprintln(a.out)
			return nil
		},
	}
	c.Flags().StringVar(&by, "by", "", "who accepted (default: config author)")
	c.Flags().StringVar(&deliver, "deliver", "", "component the item delivers to, when its tags do not say")
	c.Flags().StringVar(&remote, "remote", "origin", "git remote to push to")
	c.Flags().StringArrayVar(&trailers, "trailer", nil, "line appended to the commit message (repeatable)")
	c.Flags().BoolVar(&noRelease, "no-release", false, "accept without computing or creating a release")
	c.Flags().BoolVar(&noPush, "no-push", false, "do not push the commit and tags (implies --no-publish)")
	c.Flags().BoolVar(&noPublish, "no-publish", false, "do not push template components to their publish remote")
	c.Flags().BoolVar(&dryRun, "dry-run", false, "print the release plan and change nothing")
	return c
}
