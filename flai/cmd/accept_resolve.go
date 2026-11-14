package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/bytepunx/system-flow/flai/internal/guard"
	"github.com/bytepunx/system-flow/flai/internal/preview"
	"github.com/bytepunx/system-flow/flai/internal/storygit"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// resolveRefusal is the refusal of flai accept --resolve for a story it does
// not send back (ADR-0302): its text, and for --json, which prints it as
// {"refused": {...}} and exits 4, the story's ID.
type resolveRefusal struct {
	ID     string `json:"id"`
	Reason string `json:"reason"`
}

// notWithResolve are the flags of an acceptance, which flai accept --resolve,
// accepting nothing, refuses beside it.
var notWithResolve = []string{"dry-run", "batch", "verified", "evidence", "keep-removal", "reason", "without-queue-tiers", "without-kluster", "trailer"}

// resolveAccept is flai accept <id> --resolve (ADR-0302): the operator's
// request that the agent of story it, in review, rebase its branch onto the
// main branch and resolve what holds its acceptance. It accepts nothing. It
// runs the preview flai accept --dry-run runs, and when that lists a blocker
// sends the story back to its agent through the acceptance's send back
// (ADR-0248), which flai serve starts the agent on; a story not in review,
// or one nothing holds, is refused.
func (a *app) resolveAccept(cmd *cobra.Command, repo *workitem.Repo, it *workitem.Item, o acceptOptions) error {
	if err := a.resolveFlagsRefused(cmd, it.ID, o); err != nil {
		return err
	}
	only := "flai accept --resolve sends back to its agent only a story in review whose acceptance something holds"
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
		return a.resolveRefused(it.ID, fmt.Sprintf("nothing holds %s's acceptance, so its agent has nothing to resolve: accept it with flai accept %s", it.ID, it.ID))
	}
	base, err := a.mainBranch(repo.MainRoot)
	if err != nil {
		return fmt.Errorf("send %s back to its agent to rebase onto the main branch: cannot tell the main branch: %w", it.ID, err)
	}
	sb, err := a.sendBackBlocked(repo, it, by, base, res.Blockers)
	if err != nil {
		return err
	}
	if a.jsonOut {
		return a.printJSON(sb)
	}
	_, err = fmt.Fprintln(a.out, sb.Reason)
	return err
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
		return fmt.Errorf("rule: --resolve is the operator's request that %s's agent resolve what holds its acceptance, not the orchestrator's: the operator asks for it with flai accept %s --resolve (ADR-0302)", id, id)
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

// sendBackBlocked sends story it, in review, back to its agent for by, who
// asked it to rebase onto base and resolve blockers (ADR-0302): it moves it
// to in-progress as flai with a reason naming the first blocker, and opens
// the send back's thread on it listing each, or writes on the one still
// open, for flai serve to start its agent on (ADR-0248).
func (a *app) sendBackBlocked(repo *workitem.Repo, it *workitem.Item, by, base string, blockers []string) (*sentBack, error) {
	branch, n := storyBranch(it.ID), plural(len(blockers), "blocker")
	reason := fmt.Sprintf("%s asked its agent to rebase %s onto %s and resolve what holds its acceptance, %s: %s; its agent syncs it with flai stream sync %s, takes each blocker's recommended action, closes out again, and moves it to review",
		by, branch, base, n, shortBlockers(blockers), it.ID)
	if _, err := repo.Transition(it, workitem.InProgress, storygit.SentBackAuthor, reason, a.now()); err != nil {
		return nil, fmt.Errorf("send %s back to its agent: cannot move it to in-progress: %w; nothing changed", it.ID, err)
	}
	title := storygit.SentBackPrefix + fmt.Sprintf("%s cannot be accepted until its agent rebases it onto %s and resolves %s, as %s asked", branch, base, n, by)
	th, err := a.sendBackThread(repo, it.ID, title, resolveText(it.ID, branch, base, by, blockers))
	if err != nil {
		return nil, fmt.Errorf("%s was sent back to in-progress, but its thread could not be opened: %w; tell its agent with flai thread new --on %s", it.ID, err, it.ID)
	}
	a.logger().Info("story sent back to resolve its acceptance", "component", "accept", "item", it.ID, "by", by, "branch", branch, "blockers", len(blockers), "thread", th.ID)
	return &sentBack{
		ID: it.ID, SentBack: true, Branch: branch, Base: base, Conflicts: []string{}, Thread: th.ID, Blockers: blockers,
		Reason: fmt.Sprintf("%s was sent back to in-progress: %s asked its agent to rebase %s onto %s and resolve %s, on %s (ADR-0302)", it.ID, by, branch, base, n, th.ID),
	}, nil
}

// shortBlockers is the first of blockers, its first line at most 160
// characters, and how many more there are, for the move's reason.
func shortBlockers(blockers []string) string {
	first := []rune(firstLine(blockers[0]))
	s := string(first)
	if len(first) > 160 {
		s = strings.TrimSpace(string(first[:160])) + "..."
	}
	if len(blockers) > 1 {
		s += fmt.Sprintf(", and %d more", len(blockers)-1)
	}
	return s
}

// resolveText is what the send back's thread tells story's agent when by
// asked it to rebase branch onto base and resolve blockers: each blocker as
// the preview words it, and how the agent resolves them and comes back to
// review.
func resolveText(story, branch, base, by string, blockers []string) string {
	var b strings.Builder
	fmt.Fprintf(&b, "%s asked %s's agent to rebase `%s` onto `%s` and resolve what holds its acceptance, which `flai accept %s --dry-run` lists as:\n\n", by, story, branch, base, story)
	for _, bl := range blockers {
		fmt.Fprintf(&b, "- %s\n", strings.ReplaceAll(strings.TrimSpace(bl), "\n", "\n  "))
	}
	fmt.Fprintf(&b, "\nflai moved %s back to in-progress. "+
		"This thread is for %s's agent, which resolves them; the operator does not. "+
		"In its worktree, run `flai stream sync %s`, resolve each path it stops on, `git add` it, and run `flai stream sync %s` again, which continues the rebase; "+
		"take each blocker's recommended action above; close out again, with the MCP tool `verify` or `flai verify %s --job`; "+
		"and move %s to review again. flai serve resolves this thread once %s is back in review.",
		story, story, story, story, story, story, story)
	return b.String()
}
