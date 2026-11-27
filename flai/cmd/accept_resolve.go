package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"

	"github.com/bytepunx/system-flow/flai/internal/closeout"
	"github.com/bytepunx/system-flow/flai/internal/guard"
	"github.com/bytepunx/system-flow/flai/internal/preview"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// resolveRefusal is the refusal of flai accept --resolve for a story it asks
// flai serve nothing for, or one flai serve runs nothing for (ADR-0302,
// S-0804): its text, and for --json, which prints it as {"refused": {...}}
// and exits 4, the story's ID.
type resolveRefusal struct {
	ID     string `json:"id"`
	Reason string `json:"reason"`
}

// resolveAnswer is what flai accept --resolve prints with --json once flai
// serve answered with a job (S-0804): the story, the job's id, whether this
// request started it or one of the story was running already, where its
// output is kept, and what it does.
type resolveAnswer struct {
	ID      string `json:"id"`
	Job     string `json:"job"`
	Started bool   `json:"started"`
	Running bool   `json:"running"`
	Log     string `json:"log"`
	Reason  string `json:"reason"`
	// Held is a request flai serve keeps until it has restarted on a newer
	// flai, which then starts the job; Reason says so (ADR-0320).
	Held bool `json:"held,omitempty"`
}

// resolveJobWait is how long flai accept --resolve waits for flai serve to
// take its request and answer it: closeout.Wait's defaults, which tests
// shorten.
var resolveJobWait = closeout.Wait{}

// notWithResolve are the flags of an acceptance, which flai accept --resolve,
// accepting nothing, refuses beside it.
var notWithResolve = []string{"dry-run", "batch", "verified", "evidence", "keep-removal", "reason", "without-queue-tiers", "without-kluster", "trailer"}

// resolveAccept is flai accept <id> --resolve (ADR-0302, S-0804): the
// operator's request that flai serve rebase story it, in review, onto the
// main branch and close it out again, with no agent. It accepts nothing,
// moves nothing, and opens no thread. It runs the preview flai accept
// --dry-run runs, and when that lists a blocker asks flai serve for a job of
// kind closeout.KindResolve and prints its answer at once; the job sends the
// story back to its agent only when the rebase stops on conflicts or a step
// of the close-out fails. A story not in review, one nothing holds, and a
// request flai serve does not take or refuses are refused.
func (a *app) resolveAccept(cmd *cobra.Command, repo *workitem.Repo, it *workitem.Item, o acceptOptions) error {
	if err := a.resolveFlagsRefused(cmd, it.ID, o); err != nil {
		return err
	}
	only := "flai accept --resolve asks flai serve to rebase and close out only a story in review whose acceptance something holds"
	switch {
	case it.Type != workitem.Story:
		return a.resolveRefused(it.ID, fmt.Sprintf("%s is not a story: %s", it.ID, only))
	case it.Status != workitem.Review:
		return a.resolveRefused(it.ID, fmt.Sprintf("%s is in %s, not in review: %s", it.ID, it.Status, only))
	}
	by := orDefault(o.by, a.author())
	opts, _, err := a.acceptance(repo, it, acceptOptions{by: o.by, dryRun: true})
	if err != nil {
		return err
	}
	res, err := preview.AcceptWith(a.runner, repo, it, by, a.now(), a.logger(), opts)
	if err != nil {
		return err
	}
	if len(res.Blockers) == 0 {
		return a.resolveRefused(it.ID, fmt.Sprintf("nothing holds %s's acceptance, so there is nothing to rebase and resolve: accept it with flai accept %s", it.ID, it.ID))
	}
	ctx, stop := signal.NotifyContext(cmd.Context(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	ans, err := closeout.Ask(ctx, repo.CacheDir(), closeout.Request{Story: it.ID, Kind: closeout.KindResolve, By: by}, resolveJobWait)
	switch {
	case errors.Is(err, closeout.ErrServeNotTaking):
		return a.resolveRefused(it.ID, fmt.Sprintf("flai serve did not take the request for %s's rebase and resolve, so nothing runs and %s stays in review: start flai serve for this project with flai host start serve, and ask again", it.ID, it.ID))
	case err != nil:
		return fmt.Errorf("ask flai serve for %s's rebase and resolve: %w", it.ID, err)
	case ans.Refused():
		return a.resolveRefused(it.ID, fmt.Sprintf("flai serve runs no rebase and resolve for %s: %s", it.ID, ans.Reason))
	case ans.Held:
		a.logger().Info("rebase and resolve held", "component", "accept", "item", it.ID, "by", by, "blockers", len(res.Blockers))
		out := resolveAnswer{ID: it.ID, Held: true, Reason: ans.Reason}
		if a.jsonOut {
			return a.printJSON(out)
		}
		_, err = fmt.Fprintln(a.out, out.Reason)
		return err
	case ans.Recorded || ans.Job == "":
		return a.resolveRefused(it.ID, fmt.Sprintf("flai serve answered the request for %s's rebase and resolve with no job, as a flai serve older than S-0804 does, so nothing runs: upgrade flai with flai self-upgrade, restart flai serve with flai host restart serve, and ask again", it.ID))
	}
	a.logger().Info("rebase and resolve asked", "component", "accept", "item", it.ID, "by", by, "job_id", ans.Job, "started", ans.Started, "blockers", len(res.Blockers))
	out := resolveAnswer{ID: it.ID, Job: ans.Job, Started: ans.Started, Running: ans.Running(), Log: ans.Log, Reason: resolveReason(repo.CacheDir(), it.ID, ans)}
	if a.jsonOut {
		return a.printJSON(out)
	}
	_, err = fmt.Fprintln(a.out, out.Reason)
	return err
}

// resolveReason says what flai serve's answer ans to the request for
// story's rebase and resolve means: the job it started and what that job
// does, or the job of the story that was running already, a rebase and
// resolve or a close-out, read from the project's cache dir, and that
// nothing more starts; and where the job's output is.
func resolveReason(cacheDir, story string, ans closeout.Answer) string {
	does := fmt.Sprintf("flai serve rebases %s onto the main branch and closes it out at the rebased head, and %s stays in review, ready to accept, unless the rebase stops on conflicts or a step of the close-out fails, when the job sends it back to its agent", storyBranch(story), story)
	if !ans.Running() {
		return fmt.Sprintf("rebase and resolve job %s started for %s: %s; its output is in %s", ans.Job, story, does, ans.Log)
	}
	if j, ok, err := closeout.ReadJob(cacheDir, story); err == nil && ok && j.ID == ans.Job && j.Kind != closeout.KindResolve {
		return fmt.Sprintf("close-out job %s is running for %s, so no rebase and resolve starts: ask again once it ends; its output is in %s", ans.Job, story, ans.Log)
	}
	return fmt.Sprintf("rebase and resolve job %s is already running for %s, so nothing more starts: %s; its output is in %s", ans.Job, story, does, ans.Log)
}

// resolveFlagsRefused refuses flai accept --resolve for story id beside
// what makes it anything but the operator's request: an acceptance's flags,
// --yes, the orchestrator's --by, and the orchestrator's role.
func (a *app) resolveFlagsRefused(cmd *cobra.Command, id string, o acceptOptions) error {
	for _, name := range notWithResolve {
		if cmd.Flags().Changed(name) {
			return fmt.Errorf("rule: --resolve accepts nothing, so --%s has nothing to do beside it: give one or the other (ADR-0302)", name)
		}
	}
	switch {
	case a.yes:
		return fmt.Errorf("rule: --resolve accepts nothing and commits nothing, so --yes has nothing to include: leave it out (ADR-0302)")
	case workitem.IsOrchestrator(o.by) || os.Getenv("FLAI_ROLE") == guard.RoleOrchestrate:
		return fmt.Errorf("rule: --resolve is the operator's request that flai serve rebase and resolve %s, not the orchestrator's: the operator asks for it with flai accept %s --resolve (ADR-0302)", id, id)
	}
	return nil
}

// resolveRefused is the refusal reason of flai accept --resolve for story
// id: under --json printed as {"refused": {"id", "reason"}}, and exit 4
// either way.
func (a *app) resolveRefused(id, reason string) error {
	if a.jsonOut {
		_ = a.printJSON(map[string]any{"refused": resolveRefusal{ID: id, Reason: reason}})
	}
	return &exitError{code: exitDocRefused, msg: reason}
}
