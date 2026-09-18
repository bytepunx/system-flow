package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/bytepunx/system-flow/flai/internal/release"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// acceptOptions are shared by `flai accept` and `flai move <story> done`:
// there is one way for a story to become done, and it is acceptance (S-0046).
type acceptOptions struct {
	by, deliver, remote string
	trailers            []string
	noRelease, noPush   bool
	noPublish, dryRun   bool
}

// acceptResult is what acceptance did, or under dryRun what it would do.
type acceptResult struct {
	ID        string        `json:"id"`
	Status    string        `json:"status"`
	DryRun    bool          `json:"dry_run,omitempty"`
	Resumed   bool          `json:"resumed,omitempty"` // the item was already done but never archived
	Branch    string        `json:"branch,omitempty"`
	Merged    bool          `json:"merged"`
	Archived  int           `json:"archived"`
	Plan      *release.Plan `json:"plan"`
	Tags      []string      `json:"tags"`
	Pushed    bool          `json:"pushed"`
	PushError string        `json:"push_error,omitempty"` // accepted locally; the push did not happen
	Published []string      `json:"published"`
	Blockers  []string      `json:"blockers,omitempty"` // what would stop acceptance before it changes anything
}

func newAcceptCmd(a *app) *cobra.Command {
	var o acceptOptions
	c := &cobra.Command{
		Use:   "accept <id>",
		Short: "Accept an item: merge its branch, move to done, archive, commit, release, push",
		Long: `The operator's acceptance step as one command, per
design/conventions/work-management.md and git.md:

  0. rebase the story branch and fast-forward it into the main branch
  1. move the item to done (its rules apply: children closed, criteria checked)
  2. flai archive for the item and its children and narrative
  3. compute the release; bump the template version file and changelog if the
     template is a component in the plan
  4. git commit the work item, archive, and version changes
  5. create the release tags on that commit
  6. push the branch and the tags

flai move <story> done from review runs exactly this. An item that is already
done but was never archived (an older flai, a hand edit) is completed from
step 0 without a second transition. --dry-run shows the plan and changes
nothing.`,
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
			it, err := repo.Get(args[0])
			if err != nil {
				return err
			}
			res, err := a.acceptItem(repo, it, o)
			if err != nil {
				return err
			}
			return a.printAccept(res)
		},
	}
	addAcceptFlags(c, &o)
	c.Flags().BoolVar(&o.dryRun, "dry-run", false, "print the release plan and change nothing")
	return c
}

// addAcceptFlags registers the flags acceptance shares with move.
func addAcceptFlags(c *cobra.Command, o *acceptOptions) {
	if c.Flags().Lookup("by") == nil {
		c.Flags().StringVar(&o.by, "by", "", "who accepted (default: config author)")
	}
	c.Flags().StringVar(&o.deliver, "deliver", "", "component the item delivers to, when its tags do not say")
	c.Flags().StringVar(&o.remote, "remote", "origin", "git remote to push to")
	c.Flags().StringArrayVar(&o.trailers, "trailer", nil, "line appended to the commit message (repeatable)")
	c.Flags().BoolVar(&o.noRelease, "no-release", false, "accept without computing or creating a release")
	c.Flags().BoolVar(&o.noPush, "no-push", false, "do not push the commit and tags (implies --no-publish)")
	c.Flags().BoolVar(&o.noPublish, "no-publish", false, "do not push template components to their publish remote")
}

// acceptItem runs the acceptance flow for a story or an epic.
func (a *app) acceptItem(repo *workitem.Repo, it *workitem.Item, o acceptOptions) (*acceptResult, error) {
	// Without git (a project that is not a repository) acceptance is the
	// transition and the archive; with git it also merges, commits, tags,
	// and pushes.
	useGit := a.inGitWorkTree(repo.MainRoot)
	if it.Type == workitem.Task {
		return nil, fmt.Errorf("%s is a task; accept its story instead", it.ID)
	}
	res := &acceptResult{ID: it.ID, Status: it.Status, DryRun: o.dryRun, Tags: []string{}, Published: []string{}}
	switch {
	case it.Status == workitem.Done && !it.Archived:
		// Done without acceptance: finish the job rather than refuse it.
		res.Resumed = true
	case it.Closed():
		return nil, fmt.Errorf("%s is already %s", it.ID, it.Status)
	case it.Type == workitem.Story && it.Status != workitem.Review:
		return nil, fmt.Errorf("%s is %s; a story is accepted from review", it.ID, it.Status)
	}
	if dirty := a.dirtyOutsideWip(repo); useGit && len(dirty) > 0 && !a.yes {
		return nil, fmt.Errorf("working tree has uncommitted changes outside wip (%s); commit or stash them so the acceptance commit holds only acceptance, or pass --yes to include them", strings.Join(dirty, ", "))
	}
	// Preflight: everything that would fail midway is checked before the
	// first change, so a story is never left half accepted.
	if useGit {
		if _, err := a.runner.Run(repo.MainRoot, "git", "var", "GIT_COMMITTER_IDENT"); err != nil {
			res.Blockers = append(res.Blockers, "git has no committer identity here; set user.name and user.email (git config), or restart the dashboard with flai dashboard so it passes yours into the container")
		}
	} else if _, err := os.Stat(filepath.Join(repo.MainRoot, ".git")); err == nil {
		res.Blockers = append(res.Blockers, "this is a git repository but git cannot be run here; accept from a shell with flai accept "+it.ID)
	}
	if len(res.Blockers) > 0 && !o.dryRun {
		return nil, fmt.Errorf("%s cannot be accepted yet: %s", it.ID, strings.Join(res.Blockers, "; "))
	}
	hasBranch := it.Type == workitem.Story && useGit && a.branchExists(repo.MainRoot, storyBranch(it.ID))
	if hasBranch {
		res.Branch = storyBranch(it.ID)
	}

	// 0. bring the story branch into the main branch (ADR-0019)
	planRoot := ""
	if o.dryRun {
		// The branch is not merged in a dry run, so plan from the worktree,
		// whose history already holds the story's commits.
		if wt := repo.WorktreePath(it.ID); hasBranch {
			if _, err := os.Stat(wt); err == nil {
				planRoot = wt
			}
		}
	} else if it.Type == workitem.Story {
		merged, err := a.mergeStoryBranch(repo, it.ID)
		if err != nil {
			return nil, err
		}
		res.Merged = merged
	}
	if !o.noRelease && useGit {
		plan, err := a.planReleaseAt(repo, it, o.deliver, planRoot)
		if err != nil {
			return nil, err
		}
		res.Plan = plan
		if !a.jsonOut {
			a.printPlan(plan)
		}
	}
	if o.dryRun {
		return res, nil
	}

	// 1. done (epics: walk through review if needed)
	items, err := repo.List(false)
	if err != nil {
		return nil, err
	}
	board, err := repo.LoadBoard()
	if err != nil {
		return nil, err
	}
	if !res.Resumed {
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
			if _, err := repo.Move(it, st, workitem.MoveOptions{By: orDefault(o.by, a.author()), Now: a.now(), Items: items, Board: board}); err != nil {
				return nil, err
			}
		}
		if err := repo.Save(it); err != nil {
			return nil, err
		}
		if it.Type == workitem.Story {
			if err := board.Save(a.now().Format("2006-01-02")); err != nil {
				return nil, err
			}
		}
	}
	res.Status = it.Status
	// 2. archive
	items, _ = repo.List(false)
	ap, err := repo.PlanArchive(items, []string{it.ID})
	if err != nil {
		return nil, err
	}
	if err := repo.Archive(ap); err != nil {
		return nil, err
	}
	res.Archived = len(ap.Items)
	if err := a.refreshIndex(repo); err != nil {
		return nil, err
	}
	if !useGit {
		return res, nil
	}
	plan := res.Plan
	// 3. version files
	if plan != nil && plan.Skipped == "" {
		if err := release.Apply(a.runner, repo.Root, plan, a.now()); err != nil {
			return nil, err
		}
	}
	// 4. commit
	if _, err := a.runner.Run(repo.Root, "git", "add", "-A"); err != nil {
		return nil, err
	}
	msg := fmt.Sprintf("chore: [%s] accept and archive", it.ID)
	if plan != nil && plan.Skipped == "" {
		var parts []string
		for _, st := range plan.Steps {
			parts = append(parts, fmt.Sprintf("%s %s", st.Component.Name, st.To))
		}
		msg += "; release " + strings.Join(parts, ", ")
	}
	for _, t := range o.trailers {
		msg += "\n\n" + t
	}
	if _, err := a.runner.Run(repo.Root, "git", "commit", "-q", "-m", msg); err != nil {
		return nil, err
	}
	// 5. tags
	if plan != nil && plan.Skipped == "" {
		if res.Tags, err = release.Tag(a.runner, repo.Root, plan); err != nil {
			return nil, err
		}
	}
	// 6. push
	// The acceptance is complete and committed by now; a push that cannot
	// happen (no credentials in the dashboard container, no network) is
	// reported, not treated as a failed acceptance. Pushing later is idempotent.
	remote := orDefault(o.remote, "origin")
	if !o.noPush {
		if _, err := a.runner.Run(repo.Root, "git", "remote", "get-url", remote); err == nil {
			refs := append([]string{"HEAD"}, res.Tags...)
			if _, err := a.runner.Run(repo.Root, "git", append([]string{"push", "-q", remote}, refs...)...); err != nil {
				res.PushError = firstLine(err.Error())
				a.logger().Warn("accepted locally but not pushed", "component", "git", "item", it.ID, "detail", "run: git push "+remote+" "+strings.Join(refs, " "))
			} else {
				res.Pushed = true
			}
		}
	}
	// 7. publish template components at their new version
	if plan != nil && plan.Skipped == "" && !o.noPublish && !o.noPush && res.PushError == "" {
		for _, st := range plan.Steps {
			if st.Component.Kind != "template" {
				continue
			}
			dir := filepath.Join(repo.Root, st.Component.Path)
			pub, err := a.publishTemplate(dir, "", "", true, false, false)
			if err != nil {
				return nil, fmt.Errorf("template published locally at %s but pushing to its remote failed: %w", st.To, err)
			}
			if !pub.Nothing {
				res.Published = append(res.Published, fmt.Sprintf("%s %s (%s)", pub.Remote, pub.Tag, pub.Commit))
			}
		}
	}
	return res, nil
}

func (a *app) printAccept(res *acceptResult) error {
	if a.jsonOut {
		return a.printJSON(res)
	}
	if res.DryRun {
		for _, b := range res.Blockers {
			fmt.Fprintf(a.out, "blocked: %s\n", b)
		}
		if res.Branch != "" {
			fmt.Fprintf(a.out, "would merge %s into the main branch and remove its worktree\n", res.Branch)
		}
		fmt.Fprintln(a.out, "dry run: nothing changed")
		return nil
	}
	verb := "accepted"
	if res.Resumed {
		verb = "completed the acceptance of"
	}
	fmt.Fprintf(a.out, "%s %s: done, %d items archived, committed", verb, res.ID, res.Archived)
	if res.Merged {
		fmt.Fprintf(a.out, ", %s merged and removed", res.Branch)
	}
	if len(res.Tags) > 0 {
		fmt.Fprintf(a.out, ", tagged %s", strings.Join(res.Tags, ", "))
	}
	if res.Pushed {
		fmt.Fprint(a.out, ", pushed")
	}
	if res.PushError != "" {
		fmt.Fprintf(a.out, ", NOT pushed (%s); push the commit and tags from a shell", res.PushError)
	}
	for _, p := range res.Published {
		fmt.Fprintf(a.out, ", template published to %s", p)
	}
	fmt.Fprintln(a.out)
	return nil
}

func firstLine(s string) string {
	for _, line := range strings.Split(s, "\n") {
		if line = strings.TrimSpace(line); line != "" {
			return line
		}
	}
	return s
}
