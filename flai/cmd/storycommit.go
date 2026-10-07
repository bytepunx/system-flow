package cmd

import (
	"errors"
	"fmt"
	"strings"

	"github.com/bytepunx/system-flow/flai/internal/adr"
	"github.com/bytepunx/system-flow/flai/internal/issues"
	"github.com/bytepunx/system-flow/flai/internal/itemedit"
	"github.com/bytepunx/system-flow/flai/internal/storygit"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// A command run with --commit in a story's worktree commits the files it
// wrote on the story's branch, with the story's commit prefix, and widens
// the story's touches to them, in one call (S-0275): flai issue new, bump,
// and close, and flai adr new.

// commitFlagHelp is the --commit flag's help, the same on every command
// that takes it.
const commitFlagHelp = "commit what was written on the story's branch, checked out here in its worktree, as docs: [S-nnnn] ..., and add it to the story's touches; refused, with nothing written, where no story's branch is checked out"

// storyCommit is a command's --commit: the story whose branch the files it
// writes are committed on, in the checkout at repo.Root.
type storyCommit struct {
	repo     *workitem.Repo
	story    string
	trailers []string
}

// storyCommitted is what --commit did: the commit, nil when none of the
// files changed, and what it added to the story's touches.
type storyCommitted struct {
	story   string
	commit  *storygit.Commit
	touches itemedit.Widened
}

// beginStoryCommit resolves the story --commit commits for, from the flag,
// FLAI_STORY, FLAI_AGENT, or the branch checked out, and checks before
// anything is written that the story's branch is checked out at the
// project's root with no rebase unfinished, so that a refused --commit
// writes nothing.
func (a *app) beginStoryCommit(repo *workitem.Repo, flag string, autocommit bool, trailers []string) (*storyCommit, error) {
	if autocommit {
		return nil, errors.New("--commit and --autocommit cannot be given together: --commit commits on the story's branch in its worktree, --autocommit in the main checkout; nothing was written; give one of them")
	}
	story := a.issueStory(repo, flag)
	if story == "" {
		return nil, fmt.Errorf("--commit commits on a story's branch, and no story resolves in %s (no --story, FLAI_STORY, FLAI_AGENT of the form agent-S-nnnn, or story/S-nnnn branch checked out); nothing was written; run the command in the story's worktree, which flai stream open S-nnnn opens, or leave --commit out", repo.Root)
	}
	id := workitem.CanonicalID(story)
	branch := storygit.Branch(id)
	head, err := a.runner.Run(repo.Root, "git", "rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		return nil, fmt.Errorf("read the branch checked out in %s to commit on %s: %w; nothing was written; run the command in %s's worktree, or leave --commit out", repo.Root, branch, err, id)
	}
	if head = strings.TrimSpace(head); head != branch {
		return nil, fmt.Errorf("%s has %s checked out, not %s: --commit commits on the story's branch, so nothing was written; run the command in %s's worktree (flai stream open %s opens it), or leave --commit out", repo.Root, head, branch, id, id)
	}
	if storygit.RebaseInProgress(a.runner, repo.Root) {
		return nil, fmt.Errorf("a rebase is unfinished in %s, so --commit cannot commit on %s and nothing was written: finish it with git rebase --continue, or undo it with git rebase --abort, then run the command again", repo.Root, branch)
	}
	return &storyCommit{repo: repo, story: id, trailers: trailers}, nil
}

// commitForStory commits paths, relative to the project's root, on the
// story's branch with subject, and widens the story's touches to what the
// commit changed. Overlap notices that cannot be sent are logged, not
// failed: the touches stand.
func (a *app) commitForStory(sc *storyCommit, subject string, paths []string) (storyCommitted, error) {
	done := storyCommitted{story: sc.story, touches: itemedit.Widened{Task: []string{}, Story: []string{}, Overlaps: []itemedit.Overlapping{}}}
	c, err := storygit.CommitPaths(storygit.CommitOptions{
		Runner: a.runner, Dir: sc.repo.Root, Story: sc.story, Paths: paths, Subject: subject, Trailers: sc.trailers,
	})
	if err != nil {
		return done, fmt.Errorf("%w (%s is written and not committed)", err, strings.Join(paths, ", "))
	}
	if c == nil {
		return done, nil
	}
	done.commit = c
	by, _ := agentIdentity()
	w, err := itemedit.WidenStory(itemedit.WidenOptions{Repo: sc.repo, Story: sc.story, Paths: c.Paths, By: by, Now: a.now})
	if err != nil {
		return done, fmt.Errorf("committed %s on %s, but widening %s's touches failed: %w; add %s with flai touches %s --add", short(c.Hash), storygit.Branch(sc.story), sc.story, err, strings.Join(c.Paths, ", "), sc.story)
	}
	done.touches = w
	if w.Untold != nil {
		// advisory, as flai touches has it: the touches stand
		a.logger().Warn("overlap notices not sent", "component", "touches", "item", sc.story, "err", w.Untold)
	}
	return done, nil
}

// print writes what --commit did after the command's own output: the
// commit, the paths added to the story's touches, and the stories in
// progress the story's claim grew into.
func (d storyCommitted) print(a *app) {
	if c := d.commit; c != nil {
		fmt.Fprintf(a.out, "commit: %s %s\n", short(c.Hash), c.Subject)
	} else {
		fmt.Fprintln(a.out, "commit: nothing to commit")
	}
	if len(d.touches.Story) == 0 {
		fmt.Fprintln(a.out, "touches: nothing to add")
	} else {
		fmt.Fprintf(a.out, "touches: added to %s: %s\n", d.story, strings.Join(d.touches.Story, ", "))
	}
	printOverlapping(a.out, d.touches.Overlaps)
}

// issueCommitted is an issue as flai issue new, bump, and close --commit
// --json give it: with the outcome when --report was given, the commit, and
// the paths added to the story's touches.
type issueCommitted struct {
	*issues.Issue
	Outcome      issues.Outcome   `json:"outcome,omitempty"`
	Commit       *storygit.Commit `json:"commit"`
	TouchesAdded []string         `json:"touches_added"`
}

// adrCommitted is an ADR as flai adr new --commit --json gives it: with the
// commit on the story's branch, which replaces the result's commit, and the
// paths added to the story's touches.
type adrCommitted struct {
	*adr.Result
	Commit       *storygit.Commit `json:"commit"`
	TouchesAdded []string         `json:"touches_added"`
}

// issueWritten is what an issue command wrote for --commit to commit: the
// issue's file and summary.md, or nothing when the issue's file was not
// written.
func issueWritten(repo *workitem.Repo, is *issues.Issue) ([]string, error) {
	changed, err := is.Changed(repo)
	if err != nil || len(changed) == 0 {
		return changed, err
	}
	return append(changed, issues.SummaryPath(repo)), nil
}
