package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/bytepunx/system-flow/flai/internal/pending"
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
	// Uncommitted paths outside wip. The real run refuses them unless --yes
	// includes them in the acceptance commit; a dry run reports them so the
	// choice can be made before confirming (S-0051).
	Uncommitted []string `json:"uncommitted,omitempty"`
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
	if useGit {
		res.Uncommitted = a.dirtyOutsideWip(repo)
	}
	if dirty := res.Uncommitted; len(dirty) > 0 && !a.yes && !o.dryRun {
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
	// The workflow's own rules for done (open tasks, unticked criteria) are
	// checked on a copy before the branch is merged: a story that cannot be
	// done must not have its branch merged and its worktree removed first,
	// which is what happened until S-0041 found it from the review page.
	if !res.Resumed {
		if items, err := repo.List(false); err == nil {
			probe := *it
			probe.Transitions = append([]workitem.Transition(nil), it.Transitions...)
			for _, st := range stepsToDone(it.Status) {
				if _, err := repo.Move(&probe, st, workitem.MoveOptions{By: orDefault(o.by, a.author()), Now: a.now(), Items: items}); err != nil {
					res.Blockers = append(res.Blockers, strings.TrimPrefix(err.Error(), "rule: "))
					break
				}
			}
		}
	}
	// A nature that is not accepted onto main (an experiment, ADR-0025) is
	// refused here, before the merge, not by the release plan after it.
	releasable := true
	if !o.noRelease {
		if _, err := release.LevelFor(it); err != nil {
			releasable = false
			res.Blockers = append(res.Blockers, err.Error())
		}
	}
	// A story worktree git cannot open from here (I-0017): its links are
	// absolute host paths, and this process sees the repository somewhere
	// else. Say what to do instead of failing later with git's own error.
	worktreeReadable := true
	if wt := repo.WorktreePath(it.ID); useGit && it.Type == workitem.Story {
		if _, err := os.Stat(wt); err == nil {
			if _, err := a.runner.Run(wt, "git", "rev-parse", "--git-dir"); err != nil {
				worktreeReadable = false
				a.logger().Warn("story worktree cannot be opened by git", "component", "git", "worktree", relPath(repo.MainRoot, wt), "err", err)
				res.Blockers = append(res.Blockers, fmt.Sprintf("git cannot open the story worktree %s from here: the paths git keeps for it do not exist in this environment, which happens when the dashboard sees the repository at a different path than the host does. Accept from a shell on the host with flai accept %s, or stop the dashboard and start it with flai dashboard, which mounts the repository at its host path", relPath(repo.MainRoot, wt), it.ID))
			}
		}
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
		if merged {
			a.acceptStep(it, "merged", res.Branch+" rebased and fast-forwarded into the main branch")
		}
	}
	// An unreadable worktree blocks acceptance; its release cannot be planned either.
	// Nor can the release of a nature that is refused: the preview carries the blocker instead.
	if !o.noRelease && useGit && worktreeReadable && releasable {
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
		for _, st := range stepsToDone(it.Status) {
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
	a.acceptStep(it, "done", it.ID+" moved to done")
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
	a.acceptStep(it, "archived", fmt.Sprintf("%d item(s) and the narrative archived", res.Archived))
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
	a.acceptStep(it, "committed", firstLine(msg))
	// 5. tags
	if plan != nil && plan.Skipped == "" {
		if res.Tags, err = release.Tag(a.runner, repo.Root, plan); err != nil {
			return nil, err
		}
		a.acceptStep(it, "tagged", strings.Join(res.Tags, ", "))
	}
	// 6. push
	// The acceptance is complete and committed by now; a push that cannot
	// happen (no credentials in the dashboard container, no network) is
	// reported, not treated as a failed acceptance. Pushing later is idempotent.
	remote := orDefault(o.remote, "origin")
	if !o.noPush {
		if _, err := a.runner.Run(repo.Root, "git", "remote", "get-url", remote); err == nil {
			refs := append([]string{"HEAD"}, res.Tags...)
			// said before it starts: a push can take a while, and a dashboard
			// showing the steps should say what it is waiting for (S-0078)
			a.acceptStep(it, "pushing", "pushing to "+remote)
			var err error
			for _, batch := range pending.Batches("HEAD", res.Tags) {
				if _, err = a.runner.Run(repo.Root, "git", append([]string{"push", "-q", remote}, batch...)...); err != nil {
					break
				}
			}
			if err != nil {
				res.PushError = firstLine(err.Error())
				a.logger().Warn("accepted locally but not pushed", "component", "git", "item", it.ID, "detail", "run: git push "+remote+" "+strings.Join(refs, " "))
				a.acceptStep(it, "not-pushed", "accepted locally, not pushed ("+res.PushError+"); run: flai push --pending")
			} else {
				res.Pushed = true
				a.acceptStep(it, "pushed", "pushed to "+remote)
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

// stepsToDone is the walk from a state to done, one transition at a time.
func stepsToDone(status string) []string {
	path := []string{workitem.Ready, workitem.InProgress, workitem.Review, workitem.Done}
	if status == workitem.Backlog {
		return path
	}
	for i, st := range path {
		if st == status {
			return path[i+1:]
		}
	}
	return []string{workitem.Done}
}

// acceptStep logs one completed step of an acceptance as an info event with
// stable fields, so a client such as the dashboard can show progress while
// the command runs (S-0041). The message is fixed; what varies is in fields.
func (a *app) acceptStep(it *workitem.Item, step, detail string) {
	a.logger().Info("acceptance step", "component", "accept", "item", it.ID, "step", step, "detail", detail)
}

func (a *app) printAccept(res *acceptResult) error {
	if a.jsonOut {
		return a.printJSON(res)
	}
	if res.DryRun {
		for _, b := range res.Blockers {
			fmt.Fprintf(a.out, "blocked: %s\n", b)
		}
		if len(res.Uncommitted) > 0 {
			effect := "the real run refuses until they are committed or stashed, or --yes includes them"
			if a.yes {
				effect = "--yes includes them in the acceptance commit"
			}
			fmt.Fprintf(a.out, "uncommitted outside wip: %s; %s\n", strings.Join(res.Uncommitted, ", "), effect)
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
