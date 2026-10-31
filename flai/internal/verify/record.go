package verify

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/atomicfile"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// RecordPath is where the record keeps a story's verdict at a commit: in the
// project's cache, .flai-cache/verify/<story>/<commit>.json, beside the
// story's last report (ADR-0229).
func RecordPath(project *workitem.Repo, story, commit string) string {
	return filepath.Join(project.CacheDir(), "verify", story, commit+".json")
}

// Recorded answers the verdict the record holds for the story's worktree at
// its head commit, when the worktree is clean and the main branch is at the
// record's base commit or ahead of it only by commits the sync step passes
// over (ADR-0135); ok is false, with no error, otherwise.
func Recorded(o StoryOptions) (rep Report, ok bool, err error) {
	v, err := newStoryRun(o)
	if err != nil {
		return Report{}, false, err
	}
	return v.recorded()
}

// KeepAt keeps rep in the record as the story's verdict at commit as well,
// as a close-out job does for the commit its close-out made after a passing
// run.
func KeepAt(project *workitem.Repo, commit string, rep Report) error {
	if project == nil || rep.Story == "" || commit == "" {
		return fmt.Errorf("keep the verify verdict of %q at %q: name a project, a report's story, and a commit", rep.Story, commit)
	}
	rep.Commit, rep.FromRecord = commit, false
	path := RecordPath(project, rep.Story, commit)
	data, err := json.MarshalIndent(rep, "", "  ")
	if err != nil {
		return fmt.Errorf("keep the verify verdict of %s at %s: %w", rep.Story, shortHash(commit), err)
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("keep the verify verdict of %s at %s: make %s: %w; check that the project's .flai-cache is writable", rep.Story, shortHash(commit), filepath.Dir(path), err)
	}
	if err := atomicfile.WriteFile(path, append(data, '\n'), 0o644); err != nil {
		return fmt.Errorf("keep the verify verdict of %s at %s in %s: %w; check that the project's .flai-cache is writable", rep.Story, shortHash(commit), path, err)
	}
	return nil
}

// RecordText says that a report came from the record, when its run was, and
// at which commit, as "answered from the record of <RanAt> at <commit>; ran
// nothing"; "" for a report that ran.
func (r Report) RecordText() string {
	if !r.FromRecord {
		return ""
	}
	return fmt.Sprintf("answered from the record of %s at %s; ran nothing", r.RanAt.UTC().Format(time.RFC3339), shortHash(r.Commit))
}

// readRecord is the record's verdict of story at commit and true, or false
// with no error when it holds none.
func readRecord(project *workitem.Repo, story, commit string) (Report, bool, error) {
	path := RecordPath(project, story, commit)
	data, err := os.ReadFile(path)
	if errors.Is(err, fs.ErrNotExist) {
		return Report{}, false, nil
	}
	if err != nil {
		return Report{}, false, fmt.Errorf("read the verify verdict of %s at %s in %s: %w", story, shortHash(commit), path, err)
	}
	var rep Report
	if err := json.Unmarshal(data, &rep); err != nil {
		return Report{}, false, fmt.Errorf("read the verify verdict of %s at %s in %s: %w; delete it and verify again", story, shortHash(commit), path, err)
	}
	return rep, true, nil
}

// recorded is Recorded for the run's resolved options: the record's passing
// verdict at the worktree's head, marked FromRecord, when the worktree is
// clean, the record verified the same paths with the same tiers against the
// same main branch, ran every tier this run would not leave to the merge
// queue, and the main branch is at its base commit or ahead of it by commits
// the rebase and sync steps pass now. A verdict that did not pass is not
// answered: what it stopped at may lie outside the commit, as the narrative
// does, or be a flake, and a run again resumes at it (S-0341).
func (v *storyRun) recorded() (Report, bool, error) {
	if !v.clean || v.baseCommit == "" {
		return Report{}, false, nil
	}
	rec, ok, err := readRecord(v.opts.Project, v.story, v.commit)
	if err != nil || !ok {
		return Report{}, false, err
	}
	if !rec.Passed || rec.Commit != v.commit || rec.Base != v.base || rec.BaseCommit == "" || !slices.Equal(rec.Paths, v.paths) {
		return Report{}, false, nil
	}
	if _, _, fp := v.selection(); fp == "" || rec.Fingerprint != fp {
		return Report{}, false, nil
	}
	if !v.opts.LeaveToQueue && slices.ContainsFunc(rec.Steps, func(s Step) bool { return s.State == Queued }) {
		return Report{}, false, nil
	}
	if rec.BaseCommit != v.baseCommit && !v.baseAhead(rec.BaseCommit) {
		return Report{}, false, nil
	}
	now := v.newReport(v.opts.Now())
	v.runSteps(&now, v.syncSteps())
	if !now.Passed {
		return Report{}, false, nil
	}
	rec.FromRecord = true
	return rec, true, nil
}

// baseAhead says whether the main branch's commit has from among its
// ancestors. git exits with an error both when it has not and when it cannot
// tell, and either way the record's verdict does not cover this main branch.
func (v *storyRun) baseAhead(from string) bool {
	_, err := v.opts.Git.Run(v.worktree, "git", "merge-base", "--is-ancestor", from, v.baseCommit)
	return err == nil
}

// worktreeClean says whether the worktree has nothing uncommitted, untracked
// files included, as git status --porcelain shows it.
func (v *storyRun) worktreeClean() (bool, error) {
	out, err := v.opts.Git.Run(v.worktree, "git", "status", "--porcelain", "--untracked-files=all")
	if err != nil {
		return false, fmt.Errorf("list what %s has not committed: %w", v.worktree, err)
	}
	return strings.TrimSpace(out) == "", nil
}
