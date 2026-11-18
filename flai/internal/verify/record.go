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
	"github.com/bytepunx/system-flow/flai/internal/verifyrecord"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// RecordPath is where the record keeps a story's verdict at a commit: in the
// project's cache, .flai-cache/verify/<story>/<commit>.json, beside the
// story's last report (ADR-0229).
func RecordPath(project *workitem.Repo, story, commit string) string {
	return verifyrecord.CommitPath(project.CacheDir(), story, commit)
}

// Recorded answers the verdict the record holds for the story's worktree at
// its head commit, when the worktree is clean and the main branch is at the
// record's base commit or ahead of it only by commits the sync step passes
// over (ADR-0135); ok is false, with no error, otherwise. Its SuiteRuns are the
// suites the verdict covers (ADR-0312 §7).
func Recorded(o StoryOptions) (rep Report, ok bool, err error) {
	v, err := newStoryRun(o)
	if err != nil {
		return Report{}, false, err
	}
	return v.recorded()
}

// KeepAt keeps rep in the record as the story's verdict at commit as well,
// as a close-out job does for the commit its close-out made after a passing
// run: the suites rep's own run recorded are kept as verified at exactly
// commit, and those a last report keeps from earlier runs are left out. A
// verdict the record held at commit before keeps the suites rep did not run,
// with their tier steps, so that a run of fewer suites erases none another
// run recorded there (ADR-0312 §7); one that does not read is replaced.
func KeepAt(project *workitem.Repo, commit string, rep Report) error {
	if project == nil || rep.Story == "" || commit == "" {
		return fmt.Errorf("keep the verify verdict of %q at %q: name a project, a report's story, and a commit", rep.Story, commit)
	}
	var own []verifyrecord.Suite
	for _, s := range rep.SuiteRuns {
		if s.Commit == rep.Commit && s.RanAt.Equal(rep.RanAt) {
			s.Commit, s.Clean = commit, true
			own = append(own, s)
		}
	}
	rep.Commit, rep.FromRecord, rep.SuiteRuns = commit, false, own
	if before, ok, err := readRecord(project, rep.Story, commit); err == nil && ok {
		rep = mergedRecord(before, rep)
	}
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

// mergedRecord is rep kept over before, the verdict the record held at the
// same commit: the suites before recorded that rep did not run keep their
// entries, as verifyrecord.Merge keeps them, and their tier steps, in place
// of rep's step of the same tier or after rep's steps. The steps before the
// tiers are rep's. It passed when rep did and every suite kept passed; a
// kept suite that did not pass stops it at the first kept step that failed,
// else where before stopped.
func mergedRecord(before, rep Report) Report {
	suites, kept := verifyrecord.Merge(before.SuiteRuns, rep.SuiteRuns)
	rep.SuiteRuns = suites
	keptTier, keptPassed := map[string]bool{}, true
	for _, s := range suites {
		if !slices.Contains(kept, s.Name) {
			continue
		}
		keptPassed = keptPassed && s.Passed
		for _, t := range s.Tiers {
			keptTier[t.Name] = true
		}
	}
	stoppedAt, failed := before.StoppedAt, false
	rep.Steps = slices.Clone(rep.Steps)
	for _, o := range before.Steps {
		if !o.Tier || !keptTier[o.Name] {
			continue
		}
		if o.State == Failed && !failed {
			stoppedAt, failed = o.Name, true
		}
		if i := slices.IndexFunc(rep.Steps, func(s Step) bool { return s.Tier && s.Name == o.Name }); i >= 0 {
			rep.Steps[i] = o
			continue
		}
		rep.Steps = append(rep.Steps, o)
	}
	if rep.Passed && !keptPassed {
		rep.Passed, rep.StoppedAt = false, stoppedAt
	}
	return rep
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
// clean, the record verified the same paths against the same main branch,
// every step it ran before the tiers passed, and the record covers every
// suite of the tiers this run selects (covers). A verdict whose step before
// the tiers did not pass is not answered: what it stopped at may lie outside
// the commit, as the narrative does, and a run again resumes at it (S-0341).
// The answer holds only the suites this run selects and their tier steps,
// which all passed.
func (v *storyRun) recorded() (Report, bool, error) {
	if !v.clean || v.baseCommit == "" {
		return Report{}, false, nil
	}
	rec, ok, err := readRecord(v.opts.Project, v.story, v.commit)
	if err != nil || !ok {
		return Report{}, false, err
	}
	if rec.Commit != v.commit || rec.Base != v.base || !slices.Equal(rec.Paths, v.paths) ||
		slices.ContainsFunc(rec.Steps, func(s Step) bool { return !s.Tier && s.State != Passed }) {
		return Report{}, false, nil
	}
	selected, queued, _ := v.selection()
	listed := inListOrder(v.opts.Tiers, selected, queued)
	need := suiteNames(listed)
	if !v.covers(rec.SuiteRuns, need, listed) {
		return Report{}, false, nil
	}
	now := v.newReport(v.opts.Now())
	v.runSteps(&now, v.syncSteps())
	if !now.Passed {
		return Report{}, false, nil
	}
	rec = rec.ofSuites(need)
	rec.FromRecord = true
	return rec, true, nil
}

// covers says the record's suites hold, for each suite of need, an entry
// that passed at the run's head against the same main branch, with the
// fingerprint of the run's listed tiers of it, the same tiers in the same
// order, none left to the merge queue unless the run leaves them too, and a
// base commit the main branch is at or ahead of by commits the sync step
// passes over. A report stored before suites were recorded holds none.
func (v *storyRun) covers(suites []verifyrecord.Suite, need []string, listed []Selected) bool {
	ahead := map[string]bool{}
	for _, name := range need {
		i := slices.IndexFunc(suites, func(s verifyrecord.Suite) bool { return s.Name == name })
		if i < 0 {
			return false
		}
		s, of := suites[i], ofSuite(listed, name)
		if !s.Passed || s.Commit != v.commit || s.Base != v.base || s.BaseCommit == "" || s.Fingerprint == "" ||
			s.Fingerprint != fingerprint(v.opts.Tiers, of) || !slices.Equal(s.TierNames(), tierNames(of)) {
			return false
		}
		if !v.opts.LeaveToQueue && slices.ContainsFunc(s.Tiers, func(t verifyrecord.Tier) bool { return t.State == verifyrecord.Queued }) {
			return false
		}
		if s.BaseCommit == v.baseCommit {
			continue
		}
		if _, asked := ahead[s.BaseCommit]; !asked {
			ahead[s.BaseCommit] = v.baseAhead(s.BaseCommit)
		}
		if !ahead[s.BaseCommit] {
			return false
		}
	}
	return true
}

// ofSuites is the report with only the suites named and their tier steps,
// and the steps before the tiers, as passed: what a verdict covering those
// suites answers.
func (r Report) ofSuites(names []string) Report {
	tiers := map[string]bool{}
	suites := make([]verifyrecord.Suite, 0, len(names))
	for _, s := range r.SuiteRuns {
		if !slices.Contains(names, s.Name) {
			continue
		}
		suites = append(suites, s)
		for _, t := range s.Tiers {
			tiers[t.Name] = true
		}
	}
	steps := make([]Step, 0, len(r.Steps))
	for _, s := range r.Steps {
		if !s.Tier || tiers[s.Name] {
			steps = append(steps, s)
		}
	}
	r.Steps, r.SuiteRuns, r.Passed, r.StoppedAt = steps, suites, true, ""
	return r
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
