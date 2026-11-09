package issues

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/adrindex"
	"github.com/bytepunx/system-flow/flai/internal/atomicfile"
	"github.com/bytepunx/system-flow/flai/internal/execx"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/storygit"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// Generated is every committed file flai writes from others, so that a
// rebase of story's branch stopped on them alone is resolved by writing them
// again in its worktree (ADR-0098, ADR-0228): the summary, written from the
// issue files at now's time, and the ADR index, written from the ADR files.
func Generated(r *workitem.Repo, story string, now func() time.Time) []storygit.GeneratedFile {
	return generatedIn(r, RepoFor(r, story), execx.System{}, now)
}

// generatedIn is Generated written in wt, the project r checked out in the
// worktree whose rebase stopped, reading git there with run.
func generatedIn(r, wt *workitem.Repo, run execx.Runner, now func() time.Time) []storygit.GeneratedFile {
	index := adrindex.Path(r)
	return []storygit.GeneratedFile{{
		Path: SummaryPath(r),
		Regenerate: func() error {
			_, err := WriteSummary(wt, now())
			return err
		},
		From: func(p string) bool { return isIssuePath(r, p) },
	}, {
		Path:       index,
		Regenerate: func() error { return writeIndex(run, wt, index) },
		From:       func(p string) bool { return adrindex.Source(r, p) },
	}}
}

// writeIndex writes the ADR index p, as git names it, in wt from the ADR
// files there, as the worktree is at the stop (ADR-0228). An index with the
// markers of a stop on it is rendered from the stop's three versions, so
// that the text around the table one side changed is kept. A project with no
// index has none to write.
func writeIndex(run execx.Runner, wt *workitem.Repo, p string) error {
	file := filepath.Join(wt.Root, filepath.FromSlash(p))
	data, err := os.ReadFile(file)
	if errors.Is(err, fs.ErrNotExist) {
		return nil
	}
	if err != nil {
		return err
	}
	folder := os.DirFS(filepath.Dir(file))
	var text string
	if conflicted(string(data)) {
		base, ours, theirs, serr := storygit.Stages(run, wt.Root, p)
		if serr != nil {
			return serr
		}
		text, err = adrindex.RenderStopped(folder, base, ours, theirs)
	} else {
		text, err = adrindex.Render(folder, string(data))
	}
	if err != nil {
		return fmt.Errorf("write %s from the ADR files in %s: %w", p, wt.Root, err)
	}
	if text == string(data) {
		return nil
	}
	if err := atomicfile.WriteFile(file, []byte(text), 0o644); err != nil {
		return fmt.Errorf("write %s in %s: %w; check that it is writable and sync again", p, wt.Root, err)
	}
	return nil
}

// conflicted says whether text holds a conflict's opening marker at the
// start of a line, as a file a rebase stopped on does.
func conflicted(text string) bool {
	return strings.HasPrefix(text, "<<<<<<< ") || strings.Contains(text, "\n<<<<<<< ")
}

// SyncFiles is what a rebase of story's branch resolves itself in its
// worktree, run reading git there: the generated files, the summary and the
// ADR index, the issue files it merges by their instances, the other
// markdown files inside the manifest's claims.shared it merges by sections,
// and the fold of the issues the branch added into the base branch's open
// issues of the same title (ADR-0098, ADR-0126, ADR-0175, ADR-0228).
func SyncFiles(r *workitem.Repo, story string, run execx.Runner, now func() time.Time) storygit.Files {
	return FilesIn(r, RepoFor(r, story).Root, run, now)
}

// FilesIn is what a rebase stopped in the worktree at root, any checkout of
// r's repository, resolves itself there, as SyncFiles is for a story's
// worktree: the merge queue's batch worktree among them. A markdown file
// inside r's claims.shared that is not a generated file or an issue file is
// merged by sections with storygit.MergeMarkdown; one it does not settle
// leaves the stop waiting (ADR-0175).
func FilesIn(r *workitem.Repo, root string, run execx.Runner, now func() time.Time) storygit.Files {
	copied := *r
	copied.Root = root
	wt := &copied
	generated := generatedIn(r, wt, run, now)
	isGenerated := func(p string) bool {
		return slices.ContainsFunc(generated, func(g storygit.GeneratedFile) bool { return g.Path == p })
	}
	return storygit.Files{
		Generated: generated,
		Merged: storygit.MergedFiles{
			Match: func(p string) bool { return isIssuePath(r, p) },
			Merge: func(p string) error { return mergeStopped(run, wt, p, Merge) },
		},
		Shared: storygit.MergedFiles{
			Match: sharedMarkdown(r.SharedClaims(), func(p string) bool { return isIssuePath(r, p) || isGenerated(p) }),
			Merge: func(p string) error { return mergeStopped(run, wt, p, storygit.MergeMarkdown) },
		},
		Fold: func(base string) ([]storygit.Folded, []string, error) { return foldOnto(run, wt, base, now()) },
	}
}

// sharedMarkdown says whether a path, as git names it, is a markdown file
// inside claims that the other merge steps, covered, do not settle: the one
// a stop merges by sections (ADR-0175).
func sharedMarkdown(claims manifest.Claims, covered func(string) bool) func(string) bool {
	return func(p string) bool {
		if !strings.HasSuffix(p, ".md") || covered(p) {
			return false
		}
		_, ok := claims.Match(p)
		return ok
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

// mergeStopped writes the file p, as git names it, merged with merge in wt,
// the worktree whose rebase stopped on it, from the stop's three versions.
func mergeStopped(run execx.Runner, wt *workitem.Repo, p string, merge func(base, ours, theirs string) (string, error)) error {
	base, ours, theirs, err := storygit.Stages(run, wt.Root, p)
	if err != nil {
		return err
	}
	merged, err := merge(base, ours, theirs)
	if err != nil {
		return fmt.Errorf("merge %s: %w", p, err)
	}
	if err := atomicfile.WriteFile(filepath.Join(wt.Root, filepath.FromSlash(p)), []byte(merged), 0o644); err != nil {
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
