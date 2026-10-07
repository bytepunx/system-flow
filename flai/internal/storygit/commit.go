package storygit

import (
	"errors"
	"fmt"
	"strings"

	"github.com/bytepunx/system-flow/flai/internal/execx"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// A command run with --commit in a story's worktree commits the files it
// wrote, and only those, on the story's branch with the story's commit
// prefix (S-0275), so that writing an issue or an ADR from the worktree is
// one call. Whatever else the worktree holds uncommitted stays as it was.

// DefaultCommitType is the conventional-commit type CommitPaths uses when
// none is given: the files flai's commands write are design records.
const DefaultCommitType = "docs"

// CommitOptions is what CommitPaths commits, where, and how its message reads.
type CommitOptions struct {
	Runner execx.Runner
	// Dir is the story's worktree, where the story's branch is checked out.
	Dir   string
	Story string // the story's ID, in any zero padding
	// Paths are the files the command wrote, relative to Dir.
	Paths []string
	// Type is the conventional-commit type; DefaultCommitType when empty.
	Type string
	// Subject is what changed, without the type and the story's prefix.
	Subject string
	// Trailers end the message after a blank line, one a line, such as
	// "Co-Authored-By: ...".
	Trailers []string
}

// Commit is a commit CommitPaths made.
type Commit struct {
	Hash    string   `json:"hash"`
	Subject string   `json:"subject"`
	Paths   []string `json:"paths"` // what it changed, relative to Dir as git names them
}

// CommitPaths commits the paths that changed among o.Paths on the story's
// branch in o.Dir as "<type>: [S-nnnn] <subject>" with the trailers, leaving
// every other change in the worktree, staged or not, as it was. It returns
// nil when none of the paths changed, committing nothing. It refuses a
// checkout whose branch is not the story's, and one where a rebase is
// unfinished.
func CommitPaths(o CommitOptions) (*Commit, error) {
	if o.Runner == nil || o.Dir == "" || o.Story == "" {
		return nil, errors.New("storygit.CommitPaths needs a runner, the story's worktree, and the story's ID")
	}
	subject := strings.TrimSpace(o.Subject)
	if subject == "" || strings.Contains(subject, "\n") {
		return nil, fmt.Errorf("the commit for %s needs a subject of one line, saying what changed; got %q", o.Story, o.Subject)
	}
	id := workitem.CanonicalID(o.Story)
	branch := Branch(id)
	head, err := o.Runner.Run(o.Dir, "git", "rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		return nil, fmt.Errorf("read the branch checked out in %s to commit for %s: %w", o.Dir, id, err)
	}
	if head = strings.TrimSpace(head); head != branch {
		return nil, fmt.Errorf("%s has %s checked out, not %s: --commit commits on the story's branch, so run the command in %s's worktree (flai stream open %s opens it), or use --autocommit to commit in the main checkout", o.Dir, head, branch, id, id)
	}
	if RebaseInProgress(o.Runner, o.Dir) {
		waiting := ""
		if c := Conflicts(o.Runner, o.Dir); len(c) > 0 {
			waiting = " on " + strings.Join(c, ", ")
		}
		return nil, fmt.Errorf("a rebase is unfinished in %s%s: resolve each conflicting path, git add it, and run git rebase --continue there, or undo it with git rebase --abort; then run the command again", o.Dir, waiting)
	}
	if len(o.Paths) == 0 {
		return nil, nil
	}
	spec := append([]string{"--"}, o.Paths...)
	if _, err := RunPastIndexLock(o.Runner, o.Dir, append([]string{"add", "-A"}, spec...)...); err != nil {
		return nil, fmt.Errorf("stage %s in %s: %w; check that each path exists or is tracked, then run the command again", strings.Join(o.Paths, ", "), o.Dir, err)
	}
	staged, err := o.Runner.Run(o.Dir, "git", append([]string{"-c", "core.quotePath=false", "diff", "--cached", "--name-only", "--no-renames", "--relative"}, spec...)...)
	if err != nil {
		return nil, fmt.Errorf("list the changes staged in %s: %w", o.Dir, err)
	}
	paths := nonEmptyLines(staged)
	if len(paths) == 0 {
		return nil, nil
	}
	typ := strings.TrimSpace(o.Type)
	if typ == "" {
		typ = DefaultCommitType
	}
	subject = fmt.Sprintf("%s: [%s] %s", typ, id, subject)
	message := subject
	if len(o.Trailers) > 0 {
		message += "\n\n" + strings.Join(o.Trailers, "\n")
	}
	// With paths, git commit takes only theirs, leaving the rest of the index
	// as it was.
	if _, err := RunPastIndexLock(o.Runner, o.Dir, append([]string{"commit", "-q", "-m", message, "--"}, paths...)...); err != nil {
		return nil, fmt.Errorf("commit %s on %s in %s: %w; fix what git says and run the command again", strings.Join(paths, ", "), branch, o.Dir, err)
	}
	hash, err := o.Runner.Run(o.Dir, "git", "rev-parse", "HEAD")
	if err != nil {
		return nil, fmt.Errorf("read the commit just made on %s in %s: %w", branch, o.Dir, err)
	}
	return &Commit{Hash: strings.TrimSpace(hash), Subject: subject, Paths: paths}, nil
}

// nonEmptyLines is git's output, one path a line, without empty lines.
func nonEmptyLines(out string) []string {
	var paths []string
	for _, l := range strings.Split(out, "\n") {
		if l = strings.TrimSpace(l); l != "" {
			paths = append(paths, l)
		}
	}
	return paths
}
