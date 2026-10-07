package storygit

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"slices"
	"sort"
	"strings"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/execx"
	"github.com/bytepunx/system-flow/flai/internal/gitver"
	"github.com/bytepunx/system-flow/flai/internal/messages"
	"github.com/bytepunx/system-flow/flai/internal/threads"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// What flai stream sync reads of a story's worktree before and after it
// rebases (ADR-0069): agents never rebase by hand, so sync refuses a worktree
// it could leave in a state that loses work, and names what stopped it.

// Uncommitted lists the paths in the worktree at dir with uncommitted
// changes, tracked and untracked, as git status --porcelain names them (an
// untracked folder is one path).
func Uncommitted(r execx.Runner, dir string) ([]string, error) {
	st, err := r.Run(dir, "git", "status", "--porcelain")
	if err != nil {
		return nil, err
	}
	return workitem.PorcelainPaths(st), nil
}

// RebaseInProgress reports whether a rebase has stopped in the worktree at
// dir and waits to be continued or aborted.
func RebaseInProgress(r execx.Runner, dir string) bool {
	for _, state := range []string{"rebase-merge", "rebase-apply"} {
		p, err := r.Run(dir, "git", "rev-parse", "--git-path", state)
		if err != nil || p == "" {
			continue
		}
		if !filepath.IsAbs(p) {
			p = filepath.Join(dir, p)
		}
		if _, err := os.Stat(p); err == nil {
			return true
		}
	}
	return false
}

// ContinueRebase continues the rebase stopped in the worktree at dir once
// its conflicts are resolved and staged, keeping each commit's message with
// no editor (ADR-0098). It returns git's error when the rebase stops again,
// which RebaseInProgress then reports waiting.
func ContinueRebase(r execx.Runner, dir string) error {
	_, err := r.Run(dir, "git", "-c", "core.editor=true", "rebase", "--continue")
	return err
}

// Conflicts lists the paths left unmerged in the worktree at dir, which a
// stopped rebase waits for an agent to resolve and git add.
func Conflicts(r execx.Runner, dir string) []string {
	out, err := r.Run(dir, "git", "diff", "--name-only", "--diff-filter=U")
	if err != nil {
		return nil
	}
	var paths []string
	for _, f := range strings.Split(out, "\n") {
		if f = strings.TrimSpace(f); f != "" {
			paths = append(paths, f)
		}
	}
	return paths
}

// Why a sync did not rebase a story's branch, as SyncResult.Stopped says it.
const (
	// StopUncommitted refuses a worktree with uncommitted changes, touching
	// nothing.
	StopUncommitted = "uncommitted"
	// StopRebaseInProgress refuses a worktree where a rebase already waits,
	// touching nothing.
	StopRebaseInProgress = "rebase-in-progress"
	// StopConflicts is this sync's rebase, stopped on conflicts and left in
	// progress for the agent.
	StopConflicts = "conflicts"
)

// GeneratedFile is a committed file flai writes from others, so a rebase
// stopped on it alone is resolved by writing it again (ADR-0098).
type GeneratedFile struct {
	Path string // relative to the repository root, as git names it
	// Regenerate writes the file again in the story's worktree as it is at
	// the stop.
	Regenerate func() error
}

// SyncOptions is what Sync works with.
type SyncOptions struct {
	Runner execx.Runner
	Repo   *workitem.Repo
	Story  *workitem.Item
	// Now dates the conflict messages Sync writes and the conversations and
	// old conflict threads it closes.
	Now time.Time
	// Generated is every generated file: a stop on them alone is continued,
	// and the trial merge leaves them out of a pair's conflicts.
	Generated []GeneratedFile
	// Again is the command a stop tells the agent to run once it is
	// resolved; flai stream sync <story> when empty.
	Again string
	// Log takes what Sync does and what it could not; nil discards it.
	Log *slog.Logger
}

// BranchCheck is how another open story's branch merges with this one.
type BranchCheck struct {
	Story     string   `json:"story"`
	Status    string   `json:"status"`
	Branch    string   `json:"branch"`
	Clean     bool     `json:"clean"`
	Conflicts []string `json:"conflicts"`
	// Conversation is the pair's conversation that tells of the conflict,
	// when the two conflict (ADR-0121).
	Conversation string `json:"conversation,omitempty"`
	// Thread is the old conflict thread of the pair this sync resolved, when
	// the two merge cleanly.
	Thread string `json:"thread,omitempty"`
}

// SyncResult is what a sync did to a story's branch and what its checks
// found. Synced is false and Stopped says why when the rebase was refused or
// stopped; the checks run only after a clean rebase.
type SyncResult struct {
	Story    string `json:"story"`
	Branch   string `json:"branch"`
	Base     string `json:"base"`
	Worktree string `json:"worktree"` // relative to the main checkout
	Synced   bool   `json:"synced"`
	Stopped  string `json:"stopped,omitempty"` // StopUncommitted, StopRebaseInProgress, or StopConflicts
	// Uncommitted is the worktree's uncommitted paths, when they refused it.
	Uncommitted []string `json:"uncommitted"`
	// Conflicts is the paths the waiting rebase stopped on.
	Conflicts []string `json:"conflicts"`
	Continue  string   `json:"continue,omitempty"`
	Abort     string   `json:"abort,omitempty"` // "" when there is nothing to abort
	// Regenerated is the generated files written again to continue the
	// rebase, once for each stop on them alone.
	Regenerated []string `json:"regenerated"`
	// Branches is the trial merge with every other open story's branch.
	Branches []BranchCheck `json:"branches"`
	// TrialMergeSkipped says why the trial merge did not run, when it did not.
	TrialMergeSkipped string `json:"trial_merge_skipped,omitempty"`
	// Outside is the paths the branch changed outside the story's claim.
	Outside []string `json:"outside_touches"`
}

// RebaseInProgress reports whether the result leaves a rebase waiting in the
// worktree, whether this sync stopped it or found it.
func (s SyncResult) RebaseInProgress() bool {
	return s.Stopped == StopRebaseInProgress || s.Stopped == StopConflicts
}

// Sync rebases the story's branch onto the main branch inside its worktree,
// without stashing (ADR-0069), as flai stream sync does. It refuses, touching
// nothing, a worktree with uncommitted changes or with a rebase already in
// progress. A stop whose conflicts are all generated files it resolves and
// continues (ADR-0098); on any other conflict the rebase is left in progress
// for the agent. A refusal or a stop is a result, with how to continue and
// abort; an error is a sync that could not be tried, or a rebase git failed
// without stopping. After a clean rebase it trial-merges the branch with
// every other open story's, telling each conflicting pair in its
// conversation (ADR-0121), and lists what the branch changed outside the
// story's claim; a check that fails then is logged, the sync stands, and the
// result holds what the checks found before it failed.
func Sync(o SyncOptions) (SyncResult, error) {
	if o.Log == nil {
		o.Log = slog.New(slog.DiscardHandler)
	}
	res, err := Rebase(o)
	if err != nil || !res.Synced {
		return res, err
	}
	if err := checkSync(o, &res); err != nil {
		o.Log.Warn("story branch checks failed after the rebase", "component", "git", "story", o.Story.ID, "err", err)
		return res, nil
	}
	if res.TrialMergeSkipped != "" {
		return res, nil
	}
	tellConflicts(o, &res)
	if err := resolveConflictThreads(o, &res); err != nil {
		o.Log.Warn("old conflict threads not resolved", "component", "threads", "story", o.Story.ID, "err", err)
	}
	return res, nil
}

// Rebase is Sync's rebase without the checks after it, as flai accept
// runs it before merging the branch.
func Rebase(o SyncOptions) (SyncResult, error) {
	if o.Log == nil {
		o.Log = slog.New(slog.DiscardHandler)
	}
	repo, id := o.Repo, o.Story.ID
	path := repo.WorktreePath(id)
	res := SyncResult{Story: id, Branch: Branch(id), Worktree: relPath(repo.MainRoot, path),
		Uncommitted: []string{}, Conflicts: []string{}, Regenerated: []string{}, Branches: []BranchCheck{}, Outside: []string{}}
	if _, err := os.Stat(path); err != nil {
		return res, fmt.Errorf("%s has no worktree at %s; open one with flai stream open %s", id, res.Worktree, id)
	}
	base, err := MainBranch(o.Runner, repo.MainRoot)
	if err != nil {
		return res, err
	}
	res.Base = base
	if RebaseInProgress(o.Runner, path) {
		return stopped(o, res, StopRebaseInProgress, Conflicts(o.Runner, path)), nil
	}
	dirty, err := Uncommitted(o.Runner, path)
	if err != nil {
		return res, err
	}
	if len(dirty) > 0 {
		res.Stopped, res.Uncommitted = StopUncommitted, dirty
		res.Continue = fmt.Sprintf("commit them on %s (or stash them), then run %s again", res.Branch, again(o))
		return res, nil
	}
	if _, err := o.Runner.Run(path, "git", "rebase", base); err != nil {
		if !RebaseInProgress(o.Runner, path) {
			return res, err
		}
		conflicts, waits, err := continueOverGenerated(o, path, &res)
		if err != nil {
			return res, err
		}
		if waits {
			return stopped(o, res, StopConflicts, conflicts), nil
		}
	}
	res.Synced = true
	return res, nil
}

// again is the command a stop tells the agent to run once it is resolved.
func again(o SyncOptions) string {
	if o.Again != "" {
		return o.Again
	}
	return "flai stream sync " + o.Story.ID
}

// stopped is res with the rebase waiting in the worktree, why, its
// conflicts, and how to continue or abort it.
func stopped(o SyncOptions, res SyncResult, why string, conflicts []string) SyncResult {
	res.Stopped = why
	if conflicts != nil {
		res.Conflicts = conflicts
	}
	res.Continue = fmt.Sprintf("in %s, resolve each conflicting path, git add it, and run git rebase --continue; then run %s again", res.Worktree, again(o))
	if len(conflicts) == 0 {
		res.Continue = fmt.Sprintf("in %s, git add the resolved paths and run git rebase --continue; then run %s again", res.Worktree, again(o))
	}
	res.Abort = fmt.Sprintf("in %s, run git rebase --abort, which puts %s back as it was before the sync", res.Worktree, res.Branch)
	return res
}

// continueOverGenerated continues the rebase stopped in the story's worktree
// at path for as long as every path of each stop is a generated file,
// writing each again as the worktree is at that stop, staging it, and
// continuing (ADR-0098); it adds each file it writes to res.Regenerated. It
// reports whether the rebase is left waiting for the agent, and that stop's
// conflicts.
func continueOverGenerated(o SyncOptions, path string, res *SyncResult) (conflicts []string, waits bool, err error) {
	id := o.Story.ID
	for RebaseInProgress(o.Runner, path) {
		conflicts = Conflicts(o.Runner, path)
		regen, ok := onlyGenerated(conflicts, o.Generated)
		if !ok {
			return conflicts, true, nil
		}
		for _, f := range regen {
			if err := f.Regenerate(); err != nil {
				o.Log.Warn("generated file not regenerated; the rebase is left for the agent", "component", "git", "story", id, "path", f.Path, "err", err)
				return conflicts, true, nil
			}
			res.Regenerated = append(res.Regenerated, f.Path)
			o.Log.Info("generated file regenerated to continue the rebase", "component", "git", "story", id, "path", f.Path)
		}
		if _, err := o.Runner.Run(path, "git", append([]string{"add", "--"}, conflicts...)...); err != nil {
			o.Log.Warn("generated files not staged; the rebase is left for the agent", "component", "git", "story", id, "paths", conflicts, "err", err)
			return conflicts, true, nil
		}
		if err := ContinueRebase(o.Runner, path); err != nil && !RebaseInProgress(o.Runner, path) {
			return nil, false, err
		}
	}
	return nil, false, nil
}

// onlyGenerated returns the generated files among conflicts, and whether
// they are all of them; a stop with no conflicts is not one of these.
func onlyGenerated(conflicts []string, files []GeneratedFile) ([]GeneratedFile, bool) {
	var regen []GeneratedFile
	for _, c := range conflicts {
		i := slices.IndexFunc(files, func(f GeneratedFile) bool { return f.Path == c })
		if i < 0 {
			return nil, false
		}
		regen = append(regen, files[i])
	}
	return regen, len(regen) > 0
}

// relPath is p relative to root, or p when it is not below it.
func relPath(root, p string) string {
	if rel, err := filepath.Rel(root, p); err == nil {
		return rel
	}
	return p
}

// The checks a sync runs once its rebase is clean (S-0131, ADR-0046): a
// trial merge with every other open story's branch, which writes nothing to
// any worktree, to catch the overlaps the declared touches missed while both
// stories are still open, and the paths the branch changed outside the
// story's claim, so that it is widened. A conflict is a message flai writes
// from the syncing story to the other in the pair's conversation, which then
// awaits the other story's agent; the two agents settle it between them, and
// either escalates to the operator only when they do not agree (ADR-0121).
// Conflict threads, which earlier syncs wrote, are only resolved now. The
// generated files are left out of a pair's conflicts: the rebase writes them
// again when it stops on them
// alone, so a pair whose only conflict is one merges cleanly as far as the
// check is concerned (ADR-0098). A pair's conflicts are the paths both
// stories changed since each left the main branch: the trial merge's base is
// where the two branches meet, which is the other branch's old base when it
// is stale, so it also stops where main has since changed what the other
// did. That is the other story's own rebase to settle, not a conflict
// between the two (I-0064, S-0251).

// checkSync sets res.Outside to the paths the story's branch changed since
// res.Base outside its claim, and res.Branches to the trial merge of the
// branch with the branch of every other story in progress or in review that
// has one, in ID order.
func checkSync(o SyncOptions, res *SyncResult) error {
	repo, story, r := o.Repo, o.Story, o.Runner
	items, err := repo.List(false)
	if err != nil {
		return err
	}
	mine, err := branchChanges(r, repo.MainRoot, res.Base, res.Branch)
	if err != nil {
		return err
	}
	res.Outside = outsideClaim(repo, story, items, mine)
	var others []*workitem.Item
	for _, it := range items {
		if it.Type == workitem.Story && it.ID != story.ID && (it.Status == workitem.InProgress || it.Status == workitem.Review) && BranchExists(r, repo.MainRoot, Branch(it.ID)) {
			others = append(others, it)
		}
	}
	if len(others) == 0 {
		return nil
	}
	sort.Slice(others, func(i, j int) bool { return others[i].ID < others[j].ID })
	if v, err := gitver.Installed(r); err != nil || !v.AtLeast(gitver.MergeTree) {
		have := "unreadable"
		if err == nil {
			have = v.String()
		}
		res.TrialMergeSkipped = fmt.Sprintf("git %s is older than %s, the first with merge-tree --write-tree", have, gitver.MergeTree)
		o.Log.Warn("git is too old to trial-merge story branches, skipping", "component", "git", "git", have, "needs", gitver.MergeTree.String())
		return nil
	}
	generated := make([]string, 0, len(o.Generated))
	for _, f := range o.Generated {
		generated = append(generated, f.Path)
	}
	for _, other := range others {
		conflicts, err := TrialMerge(r, repo.MainRoot, res.Branch, Branch(other.ID))
		if err != nil {
			return fmt.Errorf("trial merge of %s with %s: %w", res.Branch, Branch(other.ID), err)
		}
		theirs, err := branchChanges(r, repo.MainRoot, res.Base, Branch(other.ID))
		if err != nil {
			return err
		}
		conflicts = changedByBoth(withoutGenerated(conflicts, generated), mine, theirs)
		res.Branches = append(res.Branches, BranchCheck{Story: other.ID, Status: other.Status, Branch: Branch(other.ID), Clean: len(conflicts) == 0, Conflicts: conflicts})
	}
	return nil
}

// withoutGenerated is conflicts less the generated files, which sync and
// acceptance write again when a rebase stops on them alone, so no agent
// settles them (ADR-0098).
func withoutGenerated(conflicts, generated []string) []string {
	out := []string{}
	for _, p := range conflicts {
		if !coveredBy(p, generated) {
			out = append(out, p)
		}
	}
	return out
}

// changedByBoth is conflicts less the paths only one of the two branches
// changed, mine and theirs: those conflict with what main brought, not with
// each other (I-0064).
func changedByBoth(conflicts, mine, theirs []string) []string {
	out := []string{}
	for _, p := range conflicts {
		if slices.Contains(mine, p) && slices.Contains(theirs, p) {
			out = append(out, p)
		}
	}
	return out
}

// branchChanges is the paths branch changed since it left base, from the
// merge base of the two.
func branchChanges(r execx.Runner, root, base, branch string) ([]string, error) {
	diff, err := r.Run(root, "git", "diff", "--name-only", base+"..."+branch)
	if err != nil {
		return nil, err
	}
	out := []string{}
	for _, p := range strings.Split(diff, "\n") {
		if p = strings.TrimSpace(p); p != "" {
			out = append(out, p)
		}
	}
	return out, nil
}

// outsideClaim is the paths of changed, what story's branch changed since it
// left the main branch, that no entry of its claim (ADR-0046, ADR-0096: its
// touches, a folder narrowed to its tasks' touches inside it, and its open
// tasks', a component as its path) covers. The wip folder is flai's and is
// left out.
func outsideClaim(repo *workitem.Repo, story *workitem.Item, items []*workitem.Item, changed []string) []string {
	claim := workitem.NewHolds(items, repo.Manifest.Projects).Claim(story)
	wip := strings.TrimSuffix(repo.Manifest.Layout["wip"], "/")
	if wip == "" {
		wip = "wip"
	}
	out := []string{}
	for _, p := range changed {
		if coveredBy(p, []string{wip}) || coveredBy(p, claim) {
			continue
		}
		out = append(out, p)
	}
	return out
}

// coveredBy says whether path is one of entries or lies below one.
func coveredBy(path string, entries []string) bool {
	for _, e := range entries {
		if path == e || strings.HasPrefix(path, e+"/") {
			return true
		}
	}
	return false
}

// TrialMerge merges two branches in git's object store only and returns the
// paths that conflict, none when they merge cleanly. merge-tree exits 1 on
// conflicts, printing the tree and then one conflicted path per line.
func TrialMerge(r execx.Runner, root, ours, theirs string) ([]string, error) {
	out, err := r.Run(root, "git", "merge-tree", "--write-tree", "--name-only", "--no-messages", ours, theirs)
	if err == nil {
		return []string{}, nil
	}
	var exit *exec.ExitError
	if !errors.As(err, &exit) || exit.ExitCode() != 1 {
		return nil, err
	}
	lines := strings.Split(out, "\n")
	seen := map[string]bool{}
	paths := []string{}
	for _, l := range lines[1:] {
		l = strings.TrimSpace(l)
		if l == "" {
			break
		}
		if !seen[l] {
			seen[l] = true
			paths = append(paths, l)
		}
	}
	sort.Strings(paths)
	return paths, nil
}

// ConflictAuthor writes the conflict messages, closes the conversations
// they are in, and resolves the old conflict threads.
const ConflictAuthor = "flai"

var (
	conflictTitlePattern   = regexp.MustCompile(`^(S-\d+) and (S-\d+) conflict when merged$`)
	conflictMessagePattern = regexp.MustCompile(`^` + regexp.QuoteMeta(Prefix) + `S-\d+ and ` + regexp.QuoteMeta(Prefix) + `S-\d+ conflict when merged\.`)
)

// IsConflictTitle reports whether title is a conflict thread's, as syncs
// before ADR-0121 opened them.
func IsConflictTitle(title string) bool { return conflictTitlePattern.MatchString(title) }

// ConflictTitle names the pair's old conflict thread, the same whichever
// story synced.
func ConflictTitle(a, b string) string {
	if b < a {
		a, b = b, a
	}
	return a + " and " + b + " conflict when merged"
}

// ConflictText is the message that tells a pair of its conflicting paths,
// the same whichever story synced, so that a sync that finds nothing new
// writes nothing.
func ConflictText(a, b string, paths []string) string {
	if b < a {
		a, b = b, a
	}
	var list strings.Builder
	for _, p := range paths {
		fmt.Fprintf(&list, "- `%s`\n", p)
	}
	return fmt.Sprintf("%s and %s conflict when merged.\n\n"+
		"A trial merge of the two at flai stream sync conflicts in:\n\n%s\n"+
		"Whichever of %s and %s is accepted second will stop on these paths when it rebases. "+
		"Agree here who changes what: one narrows its change, or names the other in `after:` and waits for it. "+
		"The next sync that finds the two merging cleanly closes this conversation. "+
		"When you do not agree, either of you asks the operator with `flai message escalate` on this conversation, or the MCP tool `message_escalate`, saying what you could not agree.",
		Branch(a), Branch(b), list.String(), a, b)
}

// tellConflicts tells each pair of the story and another open story whose
// branches conflict in the pair's conversation (ADR-0121): a message from
// the story to the other about the conflicting paths, unless the newest
// conflict message there names the same paths. It closes the pair's open
// conversations whose newest message by flai tells of a conflict once the
// two merge cleanly or the other story is no longer in progress or in
// review. It sets each conflicting branch's Conversation. A conversation
// that cannot be read or written is logged and the others are still told.
func tellConflicts(o SyncOptions, res *SyncResult) {
	repo, story := o.Repo, o.Story
	convs, err := openConversations(repo, story.ID)
	if err != nil {
		o.Log.Warn("conflict messages not written", "component", "messages", "story", story.ID, "err", err)
		return
	}
	checked := map[string]bool{}
	for i, b := range res.Branches {
		other := workitem.CanonicalID(b.Story)
		checked[other] = true
		pair := between(convs, other)
		if b.Clean {
			closeConflicts(o, pair, cleanReason(res.Branch, b.Branch, story.ID))
			continue
		}
		id, err := tellConflict(o, b, pair)
		if err != nil {
			o.Log.Warn("conflict message not written", "component", "messages", "story", story.ID, "item", b.Story, "err", err)
			continue
		}
		res.Branches[i].Conversation = id
	}
	// A pair's conversation whose other story is open no longer, sent back
	// out of in progress or review: what is left of the conflict is this
	// story's own rebase to settle. A story accepted, cancelled, or archived
	// already closes its conversations (ADR-0120).
	for _, c := range convs {
		other := workitem.CanonicalID(c.Other(story.ID))
		if checked[other] {
			continue
		}
		it, err := repo.Get(other)
		if err != nil || it.Status == workitem.InProgress || it.Status == workitem.Review {
			continue
		}
		closeConflicts(o, []*messages.Conversation{c}, goneReason(other, it.Status, story.ID))
	}
}

// tellConflict adds the conflict message for b to the pair's conversation,
// or starts one, unless the newest conflict message in pair, the pair's open
// conversations, already says it. It returns the conversation that tells it.
func tellConflict(o SyncOptions, b BranchCheck, pair []*messages.Conversation) (string, error) {
	text := ConflictText(o.Story.ID, b.Story, b.Conflicts)
	for i := len(pair) - 1; i >= 0; i-- {
		if told, _ := toldConflict(pair[i]); told != "" {
			if told == text {
				return pair[i].ID, nil
			}
			break
		}
	}
	c, err := messages.Notify(o.Repo, messages.SendOptions{From: o.Story.ID, To: b.Story, Author: ConflictAuthor, Text: text, About: b.Conflicts, Now: o.Now})
	if err != nil {
		return "", err
	}
	return c.ID, nil
}

// closeConflicts closes each of convs whose newest message by flai tells of
// a conflict, with reason; one that cannot be closed is logged.
func closeConflicts(o SyncOptions, convs []*messages.Conversation, reason string) {
	for _, c := range convs {
		if _, newest := toldConflict(c); !newest {
			continue
		}
		if _, err := messages.Close(o.Repo, c.ID, ConflictAuthor, reason, o.Now); err != nil {
			o.Log.Warn("conflict conversation not closed", "component", "messages", "story", o.Story.ID, "conversation", c.ID, "err", err)
		}
	}
}

// toldConflict is the text of the newest conflict message flai wrote in c,
// "" when none, and whether it is the newest message flai wrote there.
func toldConflict(c *messages.Conversation) (text string, newest bool) {
	es := c.Entries()
	newest = true
	for i := len(es) - 1; i >= 0; i-- {
		e := es[i]
		if e.Author != ConflictAuthor || e.Story == "" {
			continue
		}
		if conflictMessagePattern.MatchString(e.Text) {
			return e.Text, newest
		}
		newest = false
	}
	return "", false
}

// openConversations is the story's conversations stored open that read as
// open, in ID order.
func openConversations(repo *workitem.Repo, story string) ([]*messages.Conversation, error) {
	all, err := messages.For(repo, story)
	if err != nil {
		return nil, err
	}
	var out []*messages.Conversation
	for _, c := range all {
		if closed, _ := c.Closed(repo); c.Status == messages.StatusOpen && !closed {
			out = append(out, c)
		}
	}
	return out, nil
}

// between is those of convs that name the other story, in their order.
func between(convs []*messages.Conversation, other string) []*messages.Conversation {
	var out []*messages.Conversation
	for _, c := range convs {
		if c.Names(other) {
			out = append(out, c)
		}
	}
	return out
}

// cleanReason says that two branches merge cleanly at the story's sync.
func cleanReason(mine, theirs, story string) string {
	if theirs < mine {
		mine, theirs = theirs, mine
	}
	return fmt.Sprintf("%s and %s merge cleanly at the sync of %s", mine, theirs, story)
}

// goneReason says that the other story, in state, is open no longer at the
// story's sync.
func goneReason(other, state, story string) string {
	return fmt.Sprintf("%s is %s, no longer open, at the sync of %s", other, state, story)
}

// resolveConflictThreads resolves the conflict threads syncs opened before
// ADR-0121, one per pair, once the pair merges cleanly or the other story is
// no longer open, setting the clean branch's Thread, and mirrors the
// narratives of the stories they are on. It opens none and adds to none.
func resolveConflictThreads(o SyncOptions, res *SyncResult) error {
	repo, story, now := o.Repo, o.Story, o.Now
	all, err := threads.List(repo)
	if err != nil {
		return err
	}
	open := map[string]*threads.Thread{}
	for _, th := range all {
		if th.Open() && IsConflictTitle(th.Title) {
			open[th.Title] = th
		}
	}
	var changed []*threads.Thread
	for i, b := range res.Branches {
		title := ConflictTitle(story.ID, b.Story)
		th := open[title]
		delete(open, title)
		if !b.Clean || th == nil {
			continue
		}
		if th, err = threads.Resolve(repo, th.ID, ConflictAuthor, cleanReason(res.Branch, b.Branch, story.ID), now); err != nil {
			return err
		}
		res.Branches[i].Thread = th.ID
		changed = append(changed, th)
	}
	// A pair's thread whose other story is no longer open: its branch is
	// merged or dropped, and what is left of the conflict is this story's
	// own rebase to settle.
	for title, th := range open {
		m := conflictTitlePattern.FindStringSubmatch(title)
		other := ""
		switch story.ID {
		case m[1]:
			other = m[2]
		case m[2]:
			other = m[1]
		default:
			continue
		}
		state := "gone"
		if it, err := repo.Get(other); err == nil {
			if it.Status == workitem.InProgress || it.Status == workitem.Review {
				continue
			}
			state = it.Status
		}
		if th, err = threads.Resolve(repo, th.ID, ConflictAuthor, goneReason(other, state, story.ID), now); err != nil {
			return err
		}
		changed = append(changed, th)
	}
	for _, th := range changed {
		if s := threads.StoryOf(repo, th); s != "" {
			if err := threads.MirrorNarrative(repo, s); err != nil {
				return err
			}
		}
	}
	return nil
}
