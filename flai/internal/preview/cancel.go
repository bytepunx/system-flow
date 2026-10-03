// Package preview works out what a command would do without doing it: a
// cancellation's cascade, what blocks an acceptance, what a push would push,
// and what a publish would release. Each command's --dry-run prints it, its
// real run starts from it, and flai serve answers the dashboard's previews
// with it in its own process (S-0159).
package preview

import (
	"os"
	"path/filepath"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/execx"
	"github.com/bytepunx/system-flow/flai/internal/storygit"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// A cancellation takes everything open under the item with it (ADR-0028),
// so flai move says what that is before it happens, and afterwards names
// what it left for a person: git is never touched.

// LeftBehind is what a cancelled story still has on disk and in git.
type LeftBehind struct {
	ID        string `json:"id"`
	Narrative string `json:"narrative,omitempty"`
	Branch    string `json:"branch,omitempty"`
	Worktree  string `json:"worktree,omitempty"`
}

// Cancellation is the JSON of a move to cancelled, or under DryRun of what
// it would do.
type Cancellation struct {
	ID         string              `json:"id"`
	Status     string              `json:"status"`
	DryRun     bool                `json:"dry_run,omitempty"`
	Cancelled  []workitem.Cascaded `json:"cancelled"`
	LeftBehind []LeftBehind        `json:"left_behind"`
	Warnings   []string            `json:"warnings"`
	// Followed is the walk a cancelled story's epic takes with it: to review
	// at most, when its other stories are all in review or done (S-0200).
	Followed *workitem.Followed `json:"followed,omitempty"`
}

// Cancel is what moving it to cancelled would do: the item's own rules,
// tried on a copy so that a refusal is an error before anything changes,
// everything open under it that would go with it, and what would be left
// on disk and in git.
func Cancel(r execx.Runner, repo *workitem.Repo, it *workitem.Item, by, reason string, now time.Time) (*Cancellation, error) {
	items, err := repo.List(false)
	if err != nil {
		return nil, err
	}
	probe := *it
	probe.Transitions = append([]workitem.Transition(nil), it.Transitions...)
	if _, err := repo.Move(&probe, workitem.Cancelled, workitem.MoveOptions{By: by, Reason: reason, Now: now, Items: items}); err != nil {
		return nil, err
	}
	res := &Cancellation{ID: it.ID, Status: it.Status, DryRun: true, Cancelled: []workitem.Cascaded{}, Warnings: []string{}}
	for _, c := range workitem.CancelPlan(items, it.ID) {
		res.Cancelled = append(res.Cancelled, workitem.Cascaded{ID: c.ID, Type: c.Type, Title: c.Title, From: c.Status})
	}
	if it.Type == workitem.Story {
		if _, res.Followed, err = EpicWalk(repo, &probe, it.Status, by, now, false); err != nil {
			return nil, err
		}
	}
	res.LeftBehind = Left(r, repo, it, res.Cancelled)
	return res, nil
}

// Left looks at each story in the cancellation for a narrative, a story
// branch, and a worktree. Nothing is removed: unmerged work is the
// operator's to keep or delete.
func Left(r execx.Runner, repo *workitem.Repo, it *workitem.Item, with []workitem.Cascaded) []LeftBehind {
	ids := []string{}
	if it.Type == workitem.Story {
		ids = append(ids, it.ID)
	}
	for _, c := range with {
		if c.Type == workitem.Story {
			ids = append(ids, c.ID)
		}
	}
	git := storygit.InWorkTree(r, repo.MainRoot)
	out := []LeftBehind{}
	for _, id := range ids {
		l := LeftBehind{ID: id}
		if _, err := os.Stat(repo.NarrativePath(id)); err == nil {
			l.Narrative = rel(repo.MainRoot, repo.NarrativePath(id))
		}
		if git && storygit.BranchExists(r, repo.MainRoot, storygit.Branch(id)) {
			l.Branch = storygit.Branch(id)
		}
		if _, err := os.Stat(repo.WorktreePath(id)); err == nil {
			l.Worktree = rel(repo.MainRoot, repo.WorktreePath(id))
		}
		if l.Narrative != "" || l.Branch != "" || l.Worktree != "" {
			out = append(out, l)
		}
	}
	return out
}

// rel shows a path relative to root when possible, as the commands do.
func rel(root, p string) string {
	if r, err := filepath.Rel(root, p); err == nil {
		return r
	}
	return p
}
