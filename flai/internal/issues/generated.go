package issues

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/execx"
	"github.com/bytepunx/system-flow/flai/internal/storygit"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// Generated is every committed file flai writes from others, so that a
// rebase of story's branch stopped on them alone is resolved by writing them
// again in its worktree (ADR-0098): the summary alone today, written from
// the issue files at now's time.
func Generated(r *workitem.Repo, story string, now func() time.Time) []storygit.GeneratedFile {
	wt := RepoFor(r, story)
	return []storygit.GeneratedFile{{
		Path: SummaryPath(r),
		Regenerate: func() error {
			_, err := WriteSummary(wt, now())
			return err
		},
	}}
}

// SyncFiles is what a rebase of story's branch resolves itself in its
// worktree, run reading git there: the generated files, the issue files it
// merges by their instances, and the fold of the issues the branch added
// into the base branch's open issues of the same title (ADR-0098, ADR-0126).
func SyncFiles(r *workitem.Repo, story string, run execx.Runner, now func() time.Time) storygit.Files {
	wt := RepoFor(r, story)
	return storygit.Files{
		Generated: Generated(r, story, now),
		Merged: storygit.MergedFiles{
			Match: func(p string) bool { return isIssuePath(r, p) },
			Merge: func(p string) error { return mergeStopped(run, wt, p) },
		},
		Fold: func(base string) ([]storygit.Folded, []string, error) { return foldOnto(run, wt, base, now()) },
	}
}

// isIssuePath reports whether p, relative to the repository root and
// slash-separated as git names it, is an issue file: I-nnnn-*.md directly in
// the issues folder, not summary.md or README.md.
func isIssuePath(r *workitem.Repo, p string) bool {
	return path.Dir(p) == issuesFolder(r) && fileIDPattern.MatchString(path.Base(p))
}

// issuesFolder is the issues folder relative to the repository root, as git
// names it: the folder summary.md is in.
func issuesFolder(r *workitem.Repo) string { return path.Dir(SummaryPath(r)) }

// mergeStopped writes the issue file p, as git names it, merged by its
// instances in wt, the worktree whose rebase stopped on it, from the stop's
// three versions.
func mergeStopped(run execx.Runner, wt *workitem.Repo, p string) error {
	base, ours, theirs, err := storygit.Stages(run, wt.Root, p)
	if err != nil {
		return err
	}
	merged, err := Merge(base, ours, theirs)
	if err != nil {
		return fmt.Errorf("merge %s: %w", p, err)
	}
	if err := os.WriteFile(filepath.Join(wt.Root, filepath.FromSlash(p)), []byte(merged), 0o644); err != nil {
		return fmt.Errorf("write %s merged in %s: %w; check that it is writable and sync again", p, wt.Root, err)
	}
	return nil
}

// foldOnto folds the open issues the branch checked out in wt added into
// the base branch's open issues of the same title, reading the names of the
// base branch's issue files with git.
func foldOnto(run execx.Runner, wt *workitem.Repo, base string, now time.Time) ([]storygit.Folded, []string, error) {
	folder := issuesFolder(wt)
	out, err := run.Run(wt.Root, "git", "ls-tree", "-z", "--name-only", base, "--", folder+"/")
	if err != nil {
		return nil, nil, fmt.Errorf("list the issue files on %s in %s: %w", base, wt.Root, err)
	}
	var names []string
	for _, n := range strings.Split(out, "\x00") {
		if n = strings.TrimSpace(n); n != "" {
			names = append(names, path.Base(n))
		}
	}
	folds, paths, err := Fold(wt, names, now)
	if err != nil {
		return nil, nil, err
	}
	folded := make([]storygit.Folded, len(folds))
	for i, f := range folds {
		folded[i] = storygit.Folded{From: f.From, Into: f.Into}
	}
	return folded, paths, nil
}
