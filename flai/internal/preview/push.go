package preview

import (
	"errors"
	"fmt"

	"github.com/bytepunx/system-flow/flai/internal/execx"
	"github.com/bytepunx/system-flow/flai/internal/pending"
	"github.com/bytepunx/system-flow/flai/internal/release"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// Pushing is what flai push --pending finds before it pushes.
type Pushing struct {
	// Result is what --json prints: pushed, the release auto-publish
	// computes, and then either why nothing is pushed, the tags made
	// locally with nothing to push them to, or what is unpushed.
	Result map[string]any
	// Plans are the releases auto-publish computed first (S-0094).
	Plans []*release.PendingPlan
	// Unpushed is what the push sends; nil when nothing is to be pushed,
	// and Result says why.
	Unpushed *pending.Unpushed
}

// Diverged is a remote with commits this clone lacks: nothing is pushed,
// and fetching and merging is the way on.
type Diverged struct{ Message string }

func (d *Diverged) Error() string { return d.Message }

// IsDiverged reports whether err is a remote that has moved.
func IsDiverged(err error) bool {
	var d *Diverged
	return errors.As(err, &d)
}

// Push decides what a push of pending acceptances sends from the main
// checkout at root: the branch and the release tags of acceptances ahead
// of the remote-tracking branch, and tags, the tags auto-publish made just
// now. plans are the releases auto-publish computed, if it did. A remote
// that has moved is a *Diverged.
func Push(r execx.Runner, root string, plans []*release.PendingPlan, tags []string) (*Pushing, error) {
	u := pending.Detect(r, root)
	if u == nil {
		u = pending.TagsOnly(r, root, tags)
	}
	p := &Pushing{Result: map[string]any{"pushed": false}, Plans: plans}
	if len(plans) > 0 {
		p.Result["release"] = plans
	}
	switch {
	case u == nil && len(tags) > 0:
		// tagged and committed locally, but nothing ahead of a remote to
		// push to (no upstream configured)
		p.Result["tags"] = tags
		return p, nil
	case u == nil:
		p.Result["reason"] = "nothing pending"
		return p, nil
	case !u.Pending() && len(tags) == 0:
		p.Result["reason"] = fmt.Sprintf("%s is ahead of %s by %d commit(s), none of them an acceptance; that is yours to push with git", u.Branch, u.Upstream, u.Commits)
		return p, nil
	case u.Behind > 0:
		return nil, &Diverged{Message: fmt.Sprintf("conflict: %s and %s have diverged: %d commit(s) here and %d there. Fetch and merge first (git fetch %s && git merge %s), then run this again; flai never forces a push", u.Branch, u.Upstream, u.Commits, u.Behind, u.Remote, u.Upstream)}
	}
	p.Unpushed = u
	p.Result["unpushed"] = u
	return p, nil
}

// PushDryRun is what flai push --pending --dry-run says: with auto-publish
// on, the release it would tag first, and what it would push. Nothing is
// applied, tagged, or pushed. With auto-publish on, a clone whose release
// tags lag the remote's, or that cannot ask it, is refused as the push
// itself would be (S-0174).
func PushDryRun(r execx.Runner, root string, repo *workitem.Repo, autoPublish bool) (*Pushing, error) {
	var plans []*release.PendingPlan
	if autoPublish {
		if refusal := release.CheckRemote(r, root, repo.Manifest).Refusal(); refusal != "" {
			return nil, &Diverged{Message: refusal}
		}
		var err error
		if plans, err = release.Pending(r, root, repo.Manifest, repo); err != nil {
			return nil, err
		}
	}
	p, err := Push(r, root, plans, nil)
	if err != nil {
		return nil, err
	}
	if p.Unpushed != nil {
		p.Result["dry_run"] = true
	}
	return p, nil
}

// Publishing is what flai release --pending --dry-run prints: everything
// release.Pending would release (S-0087).
type Publishing struct {
	Plans  []*release.PendingPlan `json:"plans"`
	DryRun bool                   `json:"dry_run"`
	// Remote is set when the remote has something to say: a component whose
	// tags here lag the remote's (S-0174), or a remote branch with commits
	// this clone lacks (ADR-0067), when Plans is empty, since a plan built on
	// the stale tag would publish what is already published, and one built
	// on the stale branch could not be pushed; or a remote that could not be
	// asked, when Plans is what this clone's tags alone say.
	Remote *release.RemoteTags `json:"remote,omitempty"`
	// Unplanned are the accepted items no plan covers, with why (I-0024).
	Unplanned []release.Unplanned `json:"unplanned,omitempty"`
}

// Publish is what flai release --pending would release, without changing
// anything.
func Publish(r execx.Runner, repo *workitem.Repo) (*Publishing, error) {
	if err := execx.Require(r, "git", "Releases are git tags; install git."); err != nil {
		return nil, err
	}
	remote := release.CheckRemote(r, repo.Root, repo.Manifest)
	if remote.Lagging() {
		return &Publishing{DryRun: true, Remote: remote}, nil
	}
	b, err := release.PendingBatch(r, repo.Root, repo.Manifest, repo)
	if err != nil {
		return nil, err
	}
	return &Publishing{Plans: b.Plans, DryRun: true, Remote: remote, Unplanned: b.Unplanned}, nil
}
