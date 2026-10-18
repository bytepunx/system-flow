package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/hostapi"
	"github.com/bytepunx/system-flow/flai/internal/preview"
	"github.com/bytepunx/system-flow/flai/internal/serve"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// Accepting through the remote (S-0397, ADR-0170). Under acceptance.mode:
// pull-request, flai accept merges no story branch on this machine. Once the
// forge gate passes (ADR-0169), it refuses while local main has unpublished
// commits outside the wip folder, merges the story's pull request on the
// forge by acceptance.merge_method, in this process when the forge's token
// can be read here and through flai serve on the host otherwise, fetches
// origin's main and brings the main branch up to it, rebasing main's
// unpublished wip commits onto it, and removes the story's worktree and
// branch. Done, archive, and commit then go on as under local acceptance.

// dashboardAgent is FLAI_AGENT in the flai accept the dashboard's host API
// runs in flai serve's environment: it merges with the token there and never
// hands the merge back to flai serve, which would have no token either.
const dashboardAgent = "flaiover"

var (
	// acceptNewForge makes the forge's adapter for a merge flai accept makes
	// in its own process; tests put a fake forge here.
	acceptNewForge serve.NewForge = newForge
	// acceptMergeWait is how long flai accept waits for flai serve to take
	// and answer a merge it hands over; zero fields are RequestMerge's
	// defaults, and tests shorten them.
	acceptMergeWait serve.MergeWait
)

// remoteRoute is where a merge through the remote is made: here, when the
// forge's token can be read in this process, or by flai serve.
type remoteRoute struct {
	token    serve.ForgeToken
	direct   bool  // the token can be read here
	tokenErr error // why it cannot, when it cannot
}

// remoteRoute reads, from the configuration and the environment and calling
// nothing, whether this process can merge repo's pull requests itself.
func (a *app) remoteRoute(repo *workitem.Repo) remoteRoute {
	token := a.serveForgeToken(repo.Manifest.Acceptance.ForgeOrDefault())
	_, err := token.Read()
	return remoteRoute{token: token, direct: err == nil, tokenErr: err}
}

// forgeActionOn reports whether the forge host action is on for the project
// at root in this process's configuration.
func (a *app) forgeActionOn(root string) bool {
	cfg, _, err := a.loadConfig()
	return err == nil && cfg.ActionEnabled(hostapi.ActionForge, root)
}

// acceptThroughRemote is step 0 of a story's acceptance under pull-request
// acceptance, in place of mergeStoryBranch (ADR-0170). It refuses, before
// anything changes, what would stop it; merges the story's pull request on
// the forge, or finds it merged already; brings the main branch up to
// origin's main; and removes the story's worktree and local branch. It
// answers the paths the main branch changed, for the stories still open. A
// merge the forge refuses changes nothing; a step that fails after the merge
// says so, and a second flai accept finds the merge and makes no other.
func (a *app) acceptThroughRemote(ctx context.Context, repo *workitem.Repo, it *workitem.Item, res *preview.Acceptance) ([]string, error) {
	plan, root, id := res.Remote, mainRootOf(repo), it.ID
	main, err := a.mainBranch(root)
	if err != nil {
		return nil, fmt.Errorf("read the main checkout's branch, into which %s's pull request is merged: %w", id, err)
	}
	route := a.remoteRoute(repo)
	plan.ThroughServe = !route.direct
	if err := a.remotePreflight(repo, id, main, route); err != nil {
		return nil, err
	}
	answer, err := a.mergeRemote(ctx, repo, id, res.By, plan, route)
	if err != nil {
		return nil, err
	}
	if !answer.Merged() {
		return nil, fmt.Errorf("%s cannot be accepted, and nothing changed here: %s", id, answer.Reason)
	}
	plan.Outcome, plan.SHA = string(answer.Outcome), answer.SHA
	plan.BranchDeleted, plan.BranchReason = answer.BranchDeleted, answer.BranchReason
	if answer.PullRequest > 0 {
		plan.PullRequest = answer.PullRequest
	}
	if answer.URL != "" {
		res.PullRequest = answer.URL
	}
	a.logger().Info("pull request merged through the remote", "component", "accept", "item", id, "pull_request", plan.PullRequest,
		"outcome", plan.Outcome, "method", string(answer.Method), "through_serve", plan.ThroughServe, "branch_deleted", plan.BranchDeleted)
	if !answer.BranchDeleted && answer.BranchReason != "" {
		a.logger().Warn("story branch not deleted on origin", "component", "accept", "item", id, "reason", answer.BranchReason)
	}
	before := a.headOf(root)
	merged := fmt.Sprintf("%s's pull request #%d is merged on the forge", id, plan.PullRequest)
	if err := a.catchUpMain(repo, id, main, answer.SHA, merged); err != nil {
		return nil, err
	}
	if err := a.removeStoryCheckout(repo, id, merged); err != nil {
		return nil, err
	}
	res.Merged, res.Branch = true, storyBranch(id)
	a.acceptStep(it, "merged", remoteMergedLine(plan)+", and the main branch brought up to origin's")
	var changed []string
	if before != "" {
		if changed, err = a.changedSince(root, before); err != nil {
			a.logger().Warn("overlap notices not sent", "component", "accept", "item", id, "err", err)
		}
	}
	return changed, nil
}

// remotePreflight refuses, before anything changes, a merge through the
// remote that cannot be made or must not be: the forge host action off for
// a merge made here, the dashboard's acceptance with no token, the story's
// worktree with uncommitted changes, no origin, and local main with
// unpublished commits outside the wip folder, each named. A merge handed to
// flai serve is gated on the host action there, by flai serve's own
// configuration, and answered refused while it is off.
func (a *app) remotePreflight(repo *workitem.Repo, id, main string, route remoteRoute) error {
	root := mainRootOf(repo)
	switch {
	case route.direct && !a.forgeActionOn(root):
		return fmt.Errorf("the forge host action is off for this project, so flai accept merges no pull request and %s is not accepted; nothing changed: turn it on with flai serve enable forge, in a shell on the host, then accept again (ADR-0170)", id)
	case !route.direct && os.Getenv("FLAI_AGENT") == dashboardAgent:
		return fmt.Errorf("%s's pull request cannot be merged from the dashboard, so nothing changed: %w. The dashboard's acceptance runs in flai serve's environment, so flai serve has no token to merge with either: set %s where flai serve starts, restart it, and accept again (ADR-0170)", id, route.tokenErr, route.token.Env)
	}
	if wt := repo.WorktreePath(id); exists(wt) {
		if st, _ := a.runner.Run(wt, "git", "status", "--porcelain"); strings.TrimSpace(st) != "" {
			return fmt.Errorf("the worktree %s has uncommitted changes, which %s's pull request does not carry; commit them on %s and accept once flai serve has pushed them and read the checks again, or discard them, before accepting", relPath(root, wt), id, storyBranch(id))
		}
	}
	if out, err := a.runner.Run(root, "git", "remote", "get-url", "origin"); err != nil {
		return fmt.Errorf("the project has no origin remote, so flai accept can neither merge %s's pull request through it nor bring %s up to origin's, and nothing changed (%s): add the forge's repository with git remote add origin <url>, then accept again", id, main, gitSaid(out, err))
	}
	if out, err := a.fetchMain(root, main); err != nil {
		return fmt.Errorf("fetch %s from origin before merging %s's pull request, so nothing changed (%s): check that origin is reachable from here and has %s, then accept again", main, id, gitSaid(out, err), main)
	}
	commits, err := a.unpublishedOutsideWip(repo, main)
	if err != nil {
		return fmt.Errorf("list the commits on %s not yet on origin's %s before merging %s's pull request, so nothing changed: %w", main, main, id, err)
	}
	if len(commits) > 0 {
		return fmt.Errorf("%s has %s not on origin's %s that change paths outside %s/: %s. Under pull-request acceptance code reaches origin's %s only through a pull request, so flai accept merges nothing while they wait and nothing changed: publish them first with flai release --pending, or move them onto a story's branch so that its pull request carries them, then accept %s again (ADR-0170)",
			main, plural(len(commits), "commit"), main, wipDir(repo), strings.Join(commits, "; "), main, id)
	}
	return nil
}

// mergeRemote merges story id's pull request at plan's head and by its
// method: in this process with the token when route can read it, or by
// asking flai serve and waiting for its answer. A refusal is the answer's,
// not an error; an error is a hand-off flai serve did not answer, and
// changes nothing here. An interrupt withdraws a request flai serve has not
// taken.
func (a *app) mergeRemote(ctx context.Context, repo *workitem.Repo, id, by string, plan *preview.RemoteMerge, route remoteRoute) (serve.MergeAnswer, error) {
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt)
	defer stop()
	if route.direct {
		a.logger().Info("merging the pull request on the forge", "component", "accept", "item", id, "pull_request", plan.PullRequest, "method", string(plan.Method))
		return serve.MergeStory(ctx, serve.ForgeMerge{Repo: repo, Story: id, Head: plan.Head, Method: plan.Method, Git: a.runner, Token: route.token, NewForge: acceptNewForge, Now: a.now}), nil
	}
	a.logger().Info("merge handed to flai serve", "component", "accept", "item", id, "pull_request", plan.PullRequest, "method", string(plan.Method), "token_env", route.token.Env)
	answer, err := serve.RequestMerge(ctx, repo, serve.MergeRequest{Story: id, Head: plan.Head, Method: plan.Method, By: by, At: a.now().UTC().Format(time.RFC3339)}, acceptMergeWait)
	if err != nil {
		return answer, fmt.Errorf("%s cannot be accepted, and nothing changed here: %w", id, err)
	}
	return answer, nil
}

// catchUpMain brings the main branch up to origin's main after the merge
// (ADR-0170): it fetches main, checks that the merge commit sha is on it,
// rebases the main branch's unpublished commits, which are only under the
// wip folder, onto it in a temporary worktree, and moves the main branch
// there with git reset --keep, which keeps the main checkout's uncommitted
// edits. merged says what the forge did, for an error that comes after it.
func (a *app) catchUpMain(repo *workitem.Repo, id, main, sha, merged string) error {
	root, origin := mainRootOf(repo), "refs/remotes/origin/"+main
	again := fmt.Sprintf("then run flai accept %s again, which finds the merge and makes no second one", id)
	if out, err := a.fetchMain(root, main); err != nil {
		return fmt.Errorf("%s, but fetching %s from origin failed, so the main branch is left where it was (%s): check that origin is reachable from here, %s", merged, main, gitSaid(out, err), again)
	}
	if sha != "" {
		if _, err := a.runner.Run(root, "git", "merge-base", "--is-ancestor", sha, origin); err != nil {
			return fmt.Errorf("%s at %s, but origin's %s does not hold that commit after the fetch, so the main branch is left where it was: check that origin is the forge's repository the pull request is in, %s", merged, short(sha), main, again)
		}
	}
	ahead, err := a.runner.Run(root, "git", "rev-list", "--count", origin+"..refs/heads/"+main)
	if err != nil {
		return fmt.Errorf("%s, but flai could not count %s's commits not on origin's %s, so the main branch is left where it was (%s); %s", merged, main, main, gitSaid(ahead, err), again)
	}
	target := origin
	if strings.TrimSpace(ahead) != "0" {
		if target, err = a.rebaseUnpublished(repo, main); err != nil {
			return fmt.Errorf("%s, but %w, so the main branch is left where it was and nothing else changed: bring %s onto origin's %s yourself (commit or move aside the main checkout's uncommitted edits, then git rebase origin/%s and resolve the conflicts), %s", merged, err, main, main, main, again)
		}
	}
	if out, err := a.runner.Run(root, "git", "reset", "--quiet", "--keep", target); err != nil {
		return fmt.Errorf("%s, but git reset --keep would not bring %s up to origin's (%s): the main checkout has an uncommitted edit to a file the merge changed, which git names; commit it or move it aside, %s", merged, main, gitSaid(out, err), again)
	}
	return nil
}

// rebaseUnpublished rebases the main branch's commits not on origin's main
// onto it in a temporary detached worktree under the project's cache, since
// git rebase in the main checkout would refuse its uncommitted edits and the
// stash is shared by every worktree (ADR-0170). It answers the rebased head,
// and leaves the main branch where it was; on a conflict it aborts, naming
// the paths.
func (a *app) rebaseUnpublished(repo *workitem.Repo, main string) (string, error) {
	root := mainRootOf(repo)
	if err := os.MkdirAll(repo.CacheDir(), 0o755); err != nil {
		return "", fmt.Errorf("the temporary worktree to rebase %s's unpublished commits in could not be made under %s (%w)", main, relPath(root, repo.CacheDir()), err)
	}
	tmp, err := os.MkdirTemp(repo.CacheDir(), "accept-main-")
	if err != nil {
		return "", fmt.Errorf("the temporary worktree to rebase %s's unpublished commits in could not be made under %s (%w)", main, relPath(root, repo.CacheDir()), err)
	}
	defer func() { _ = os.RemoveAll(tmp) }()
	wt := filepath.Join(tmp, "main")
	if out, err := a.runner.Run(root, "git", "worktree", "add", "--quiet", "--detach", wt, "refs/heads/"+main); err != nil {
		return "", fmt.Errorf("git could not check %s out in a temporary worktree to rebase its unpublished commits (%s)", main, gitSaid(out, err))
	}
	defer func() {
		if out, err := a.runner.Run(root, "git", "worktree", "remove", "--force", wt); err != nil {
			a.logger().Warn("temporary worktree not removed", "component", "accept", "worktree", relPath(root, wt), "err", gitSaid(out, err))
		}
	}()
	if out, err := a.runner.Run(wt, "git", "-c", "core.editor=true", "rebase", "--quiet", "refs/remotes/origin/"+main); err != nil {
		conflicts, _ := a.runner.Run(wt, "git", "diff", "--name-only", "--diff-filter=U")
		if abort, err := a.runner.Run(wt, "git", "rebase", "--abort"); err != nil {
			a.logger().Warn("rebase in the temporary worktree not aborted", "component", "accept", "worktree", relPath(root, wt), "err", gitSaid(abort, err))
		}
		if paths := strings.Fields(conflicts); len(paths) > 0 {
			return "", fmt.Errorf("%s's unpublished commits under %s/ conflict with origin's %s in %s", main, wipDir(repo), main, workitem.Shorten(paths, 10))
		}
		return "", fmt.Errorf("%s's unpublished commits did not rebase onto origin's %s (%s)", main, main, gitSaid(out, err))
	}
	head, err := a.runner.Run(wt, "git", "rev-parse", "HEAD")
	if err != nil {
		return "", fmt.Errorf("the head of %s's commits rebased onto origin's %s could not be read (%s)", main, main, gitSaid(head, err))
	}
	return strings.TrimSpace(head), nil
}

// removeStoryCheckout removes story id's worktree, as mergeStoryBranch does,
// and its local branch with git branch -D: a squash or rebase merge leaves it
// unmerged in git's eyes (ADR-0170). merged says what the forge did, for an
// error.
func (a *app) removeStoryCheckout(repo *workitem.Repo, id, merged string) error {
	root, branch, path := mainRootOf(repo), storyBranch(id), repo.WorktreePath(id)
	again := fmt.Sprintf("then run flai accept %s again, which finds the merge and makes no second one", id)
	if exists(path) {
		if out, err := a.runner.Run(root, "git", "worktree", "remove", path); err != nil {
			return fmt.Errorf("%s and the main branch is up to origin's, but git could not remove the worktree %s (%s): remove it with git worktree remove %s, %s", merged, relPath(root, path), gitSaid(out, err), relPath(root, path), again)
		}
	}
	if a.branchExists(root, branch) {
		if out, err := a.runner.Run(root, "git", "branch", "-D", branch); err != nil {
			return fmt.Errorf("%s and the main branch is up to origin's, but git could not delete %s (%s): delete it with git branch -D %s, %s", merged, branch, gitSaid(out, err), branch, again)
		}
	}
	a.logger().Info("story branch removed", "component", "git", "branch", branch)
	return nil
}

// fetchMain fetches origin's main into its remote-tracking ref, whatever the
// remote's fetch configuration, answering git's output.
func (a *app) fetchMain(root, main string) (string, error) {
	return a.runner.Run(root, "git", "fetch", "--quiet", "origin", "+refs/heads/"+main+":refs/remotes/origin/"+main)
}

// unpublishedOutsideWip are the main branch's commits not on origin's main
// that change a path outside the wip folder, each as its short hash and
// subject, newest first.
func (a *app) unpublishedOutsideWip(repo *workitem.Repo, main string) ([]string, error) {
	root := mainRootOf(repo)
	out, err := a.runner.Run(root, "git", "log", "--format=%h %s", "refs/remotes/origin/"+main+"..refs/heads/"+main, "--", ".", ":(top,exclude)"+wipDir(repo))
	if err != nil {
		return nil, fmt.Errorf("%s", gitSaid(out, err))
	}
	var commits []string
	for _, line := range strings.Split(out, "\n") {
		if line = strings.TrimSpace(line); line != "" {
			commits = append(commits, line)
		}
	}
	return commits, nil
}

// wipDir is the project's wip folder, relative to its root.
func wipDir(repo *workitem.Repo) string {
	if d := strings.Trim(repo.Manifest.Layout["wip"], "/"); d != "" {
		return d
	}
	return "wip"
}

// exists reports whether path exists.
func exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// remotePlanLine is the dry run's line for a merge through the remote: what
// it would merge, by which method, and through whom, or that the pull
// request is merged already, and what follows.
func remotePlanLine(res *preview.Acceptance) string {
	r := res.Remote
	request := res.ID + "'s pull request"
	if r.PullRequest > 0 {
		request = fmt.Sprintf("pull request #%d", r.PullRequest)
	}
	rest := "fetch main and bring the main branch up to origin's"
	if res.Branch != "" {
		rest += fmt.Sprintf(", and remove %s and its worktree", res.Branch)
	}
	if r.AlreadyMerged {
		return fmt.Sprintf("%s is merged on the forge already: would accept %s with no second merge, %s", request, res.ID, rest)
	}
	through := ""
	if r.ThroughServe {
		through = " through flai serve"
	}
	return fmt.Sprintf("would merge %s through the remote by %s%s, then %s", request, r.Method, through, rest)
}

// remoteMergedLine says what the merge through the remote did.
func remoteMergedLine(r *preview.RemoteMerge) string {
	if r.Outcome == string(serve.MergeAlready) {
		return fmt.Sprintf("pull request #%d found merged on the forge already", r.PullRequest)
	}
	return fmt.Sprintf("pull request #%d merged on the forge by %s", r.PullRequest, r.Method)
}
