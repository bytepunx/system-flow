package verify

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/atomicfile"
	"github.com/bytepunx/system-flow/flai/internal/check"
	"github.com/bytepunx/system-flow/flai/internal/execx"
	"github.com/bytepunx/system-flow/flai/internal/gitver"
	"github.com/bytepunx/system-flow/flai/internal/storygit"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// The steps Verify runs before the tiers, in order, cheapest first.
const (
	StepRebase    = "rebase"
	StepSync      = "sync"
	StepNarrative = "narrative"
	StepCheck     = "check"
)

// StoryEnv names the variable a story's tiers run with set to the story, so
// that the scripts and tests that check the project scope themselves to it
// (ADR-0085).
const StoryEnv = "CLOSE_OUT_STORY"

// narrativeSections are the narrative's sections that must be written
// before a story goes to review.
var narrativeSections = []string{"Current state", "Next steps"}

// StoryOptions are what Verify needs to answer for a story.
type StoryOptions struct {
	// Story is the story's ID, as typed: S-12 is S-0012.
	Story string
	// Project is the project, opened at any of its checkouts.
	Project *workitem.Repo
	// Worktree is the story's worktree; "" is where flai stream open puts it.
	Worktree string
	// Base is the main branch the story's branch must contain and its
	// changes are counted against; "" is the branch the main checkout has.
	Base string
	// Tiers are the project's tiers, cheapest first, as the worktree's
	// manifest declares them.
	Tiers []Tier
	// Git runs git; nil is execx.System.
	Git execx.Runner
	// Fresh runs every selected tier, whatever the story's last report
	// would let the run reuse.
	Fresh bool
	// RunOptions are the tiers' run's; FS is the worktree, read for the
	// tiers' packages, and Now times every step and stamps the report.
	RunOptions
}

// Report is the answer to Verify: pass, or the step that failed and why.
type Report struct {
	Story string `json:"story"`
	// Commit is the worktree's HEAD when it was verified.
	Commit string `json:"commit"`
	// Base is the main branch the story was verified against, and
	// BaseCommit the commit it had then; "" when git could not read it.
	Base       string `json:"base"`
	BaseCommit string `json:"base_commit,omitempty"`
	// MergeBase is where the branch left Base, git merge-base HEAD Base;
	// "" when git could not tell.
	MergeBase string `json:"merge_base,omitempty"`
	// Patch is the stable patch id of the branch's own change, what it
	// changed since MergeBase; "" when the worktree had anything
	// uncommitted, the branch changed nothing, or git could not tell
	// (ADR-0156).
	Patch string `json:"patch,omitempty"`
	// Fingerprint hashes what selects and runs the tiers: the project's
	// tiers and the selected ones' commands. A later run with the same
	// Fingerprint reuses the tiers this one passed when it is at the same
	// Commit and BaseCommit (S-0341), or, rebased, has the same Patch, but
	// for the tiers that what Base changed between the two MergeBases sends
	// back (ADR-0156).
	Fingerprint string    `json:"fingerprint,omitempty"`
	RanAt       time.Time `json:"ran_at"`
	// DurationMS and Duration are how long the whole run took.
	DurationMS int64  `json:"duration_ms"`
	Duration   string `json:"duration,omitempty"`
	Passed     bool   `json:"passed"`
	// StoppedAt names the step that failed; "" when every step passed.
	StoppedAt string `json:"stopped_at,omitempty"`
	// Paths are the root-relative files the branch changed against Base,
	// which the tiers were selected for.
	Paths []string `json:"paths,omitempty"`
	Steps []Step   `json:"steps"`
	// Notes are the check's findings outside the story, which do not fail
	// it (S-0249); flai verify --record-issues records them.
	Notes []Note `json:"notes,omitempty"`
}

// Step is what one step of Verify did.
type Step struct {
	Name string `json:"name"`
	// Tier says the step ran a test or lint tier, named Name.
	Tier       bool      `json:"tier,omitempty"`
	State      State     `json:"state"`
	DurationMS int64     `json:"duration_ms"`
	Duration   string    `json:"duration,omitempty"`
	Findings   []Finding `json:"findings,omitempty"`
	// Note is what a step says about why it ran as it did: what a step
	// that passed passed over, as sync names the commits of Base it did not
	// need (ADR-0135), or, after a rebase, that a tier was reused across it
	// or that the paths Base brought in sent it back (ADR-0156). It is the
	// step's own, not a check note, and nothing records it as an issue.
	Note string `json:"note,omitempty"`
	// Omitted counts the findings past the run's cap.
	Omitted int `json:"omitted,omitempty"`
	// Command, Dir, and ExitCode are a tier's, as TierResult has them.
	Command  []string `json:"command,omitempty"`
	Dir      string   `json:"dir,omitempty"`
	ExitCode *int     `json:"exit_code,omitempty"`
	// ReusedFrom is when the run that ran a Reused tier began.
	ReusedFrom *time.Time `json:"reused_from,omitempty"`
	// CoveredBy names the tier that covers a skipped tier, as TierResult
	// has it (S-0342).
	CoveredBy string `json:"covered_by,omitempty"`
}

// Reused counts the report's reused tiers and answers the earliest run they
// come from; 0 and the zero time when it ran every tier it reached.
func (r Report) Reused() (int, time.Time) {
	n, from := 0, time.Time{}
	for _, s := range r.Steps {
		if s.State != Reused {
			continue
		}
		n++
		if s.ReusedFrom != nil && (from.IsZero() || s.ReusedFrom.Before(from)) {
			from = *s.ReusedFrom
		}
	}
	return n, from
}

// Note is a check finding outside the story.
type Note struct {
	Rule  string `json:"rule"`
	Level string `json:"level"`
	// Path is relative to the worktree, or for wip/ to the main checkout.
	Path    string `json:"path,omitempty"`
	Line    int    `json:"line,omitempty"`
	Message string `json:"message"`
}

// Verify runs a story's close-out checks in its worktree and answers one
// report: no rebase is unfinished, the branch contains Base but for commits
// that change only wip paths the branch does not change, the narrative's
// Current state and Next steps are written, flai check --strict scoped to the
// story passes, and then each tier the branch's changes select (SelectStory)
// passes, run with CLOSE_OUT_STORY set to the story. It stops at the first
// step that fails; the steps after it are not reached, but for a tier that
// another selected tier covers, which is skipped either way. When the story's
// last report has the same head, base commit, and fingerprint, and Fresh is
// not set, the tiers before the first it shows failed or not reached are not
// run again but reported Reused, with the time of the run that ran them; the
// steps before the tiers always run (S-0341). When the head or base commit
// differ, as after a clean rebase, they are reused only when the branch's own
// change has the same patch id, and only up to the first tier that the paths
// the main branch brought in select, which runs again with every tier after
// it (ADR-0156). It commits nothing and records
// no issue. The report, whole, is stored as the story's last (SaveReport)
// whether it passes or not. It returns an error, with no report, when the
// story, its worktree, or its changes cannot be read, and, with the report
// so far, when ctx ends during a tier or the report cannot be stored.
func Verify(ctx context.Context, opts StoryOptions) (Report, error) {
	v, err := newStoryRun(opts)
	if err != nil {
		return Report{}, err
	}
	rep, runErr := v.run(ctx)
	if err := SaveReport(opts.Project, rep); err != nil {
		return rep, errors.Join(runErr, err)
	}
	return rep, runErr
}

// SyncOnly runs only Verify's rebase and sync steps for a story, in its
// worktree, and answers their report, which it does not store, so the
// story's last report stays the last full run's: what flai verify
// --sync-only answers (ADR-0135). It returns an error, with no report, when
// the story, its worktree, or its changes cannot be read.
func SyncOnly(opts StoryOptions) (Report, error) {
	v, err := newStoryRun(opts)
	if err != nil {
		return Report{}, err
	}
	start := v.opts.Now()
	rep := v.newReport(start)
	v.runSteps(&rep, v.syncSteps())
	rep.DurationMS, rep.Duration = took(v.opts.Now().Sub(start))
	return rep, nil
}

// storyRun is one Verify call with its options resolved.
type storyRun struct {
	opts     StoryOptions
	story    string
	worktree string
	base     string
	max      int
	// commit and paths are the worktree's HEAD and what it changed, and
	// baseCommit the commit base has, "" when git cannot read it.
	commit     string
	baseCommit string
	paths      []string
	// mergeBase and patch are where the branch left base and the patch id
	// of its own change, as Report has them; Verify reads them.
	mergeBase string
	patch     string
	// note is what the step running says about what it passed over, which
	// runSteps puts on its result.
	note string
}

// newStoryRun resolves opts: the story's canonical ID, its worktree, the
// main branch, and the defaults; and reads the commit the worktree has and
// the paths it changed against the main branch.
func newStoryRun(opts StoryOptions) (*storyRun, error) {
	if opts.Project == nil {
		return nil, errors.New("verify a story: no project given; open the project first")
	}
	it, err := opts.Project.Get(opts.Story)
	if err != nil {
		return nil, fmt.Errorf("verify %s: %w; name a story of this project", opts.Story, err)
	}
	if it.Type != workitem.Story {
		return nil, fmt.Errorf("verify %s: it is a %s; name a story (S-nnnn)", it.ID, it.Type)
	}
	if opts.Git == nil {
		opts.Git = execx.System{}
	}
	if opts.Now == nil {
		opts.Now = time.Now
	}
	v := &storyRun{opts: opts, story: it.ID, worktree: opts.Worktree, base: opts.Base, max: opts.Max}
	if v.worktree == "" {
		v.worktree = opts.Project.WorktreePath(it.ID)
	}
	if st, err := os.Stat(v.worktree); err != nil || !st.IsDir() {
		return nil, fmt.Errorf("verify %s: it has no worktree at %s; open it with flai stream open %s", it.ID, v.worktree, it.ID)
	}
	if v.base == "" {
		if v.base, err = storygit.MainBranch(opts.Git, mainRoot(opts.Project)); err != nil {
			return nil, fmt.Errorf("verify %s: find the main branch its branch must contain: %w", it.ID, err)
		}
	}
	if v.max <= 0 {
		v.max = DefaultMax
	}
	if v.opts.FS == nil {
		v.opts.FS = os.DirFS(v.worktree)
	}
	commit, err := opts.Git.Run(v.worktree, "git", "rev-parse", "HEAD")
	if err != nil {
		return nil, fmt.Errorf("verify %s: read the commit its worktree %s has checked out: %w", v.story, v.worktree, err)
	}
	v.commit = strings.TrimSpace(commit)
	// a base git cannot read fails the sync step, and leaves nothing to resume
	if base, err := opts.Git.Run(v.worktree, "git", "rev-parse", "--verify", "--quiet", v.base+"^{commit}"); err == nil {
		v.baseCommit = strings.TrimSpace(base)
	}
	if v.paths, err = ChangedPaths(opts.Git, v.worktree, v.base); err != nil {
		return nil, fmt.Errorf("verify %s: %w", v.story, err)
	}
	return v, nil
}

// storyStep is a step of Verify before the tiers: its name, and what finds
// why it fails, none when it passes.
type storyStep struct {
	name string
	run  func() []Finding
}

// syncSteps are the steps that say whether the branch is ready to verify:
// no rebase left unfinished, and the main branch contained.
func (v *storyRun) syncSteps() []storyStep {
	return []storyStep{{StepRebase, v.rebase}, {StepSync, v.sync}}
}

// newReport is a report begun at start that has passed so far.
func (v *storyRun) newReport(start time.Time) Report {
	return Report{
		Story: v.story, Commit: v.commit, Base: v.base, BaseCommit: v.baseCommit, MergeBase: v.mergeBase, Patch: v.patch,
		RanAt: start.UTC(), Passed: true, Paths: v.paths, Steps: []Step{},
	}
}

// runSteps runs steps in order onto rep; once one fails, those after it are
// not reached.
func (v *storyRun) runSteps(rep *Report, steps []storyStep) {
	for _, s := range steps {
		if !rep.Passed {
			rep.Steps = append(rep.Steps, Step{Name: s.name, State: NotReached})
			continue
		}
		began := v.opts.Now()
		v.note = ""
		found := s.run()
		step := Step{Name: s.name, State: Passed, Note: v.note}
		step.DurationMS, step.Duration = took(v.opts.Now().Sub(began))
		if len(found) > 0 {
			step.State, step.Note = Failed, ""
			n := min(v.max, len(found))
			step.Findings, step.Omitted = found[:n], len(found)-n
			rep.Passed, rep.StoppedAt = false, s.name
		}
		rep.Steps = append(rep.Steps, step)
	}
}

// run runs the steps in order and stops at the first that fails. It returns
// ctx's error, with the report, when ctx ends during a tier.
func (v *storyRun) run(ctx context.Context) (Report, error) {
	start := v.opts.Now()
	// what git cannot tell is left "", which reuses nothing across a
	// rebase; the steps below find what is wrong with the branch
	v.mergeBase, v.patch, _ = BranchPatch(v.opts.Git, v.worktree, v.base)
	rep := v.newReport(start)
	v.runSteps(&rep, append(v.syncSteps(),
		storyStep{StepNarrative, v.narrative},
		storyStep{StepCheck, func() []Finding { return v.check(&rep) }},
	))
	selected := SelectStory(v.opts.FS, v.opts.Tiers, v.paths)
	rep.Fingerprint = fingerprint(v.opts.Tiers, selected)
	var runErr error
	if rep.Passed {
		reused, sentBack := v.reusable(rep.Fingerprint, selected)
		rep.Steps = append(rep.Steps, reused...)
		ro := v.opts.RunOptions
		ro.Env = append(append([]string{}, ro.Env...), StoryEnv+"="+v.story)
		var res Result
		res, runErr = Run(ctx, v.worktree, selected[len(reused):], ro)
		for i, tr := range res.Tiers {
			step := tierStep(tr)
			if i == 0 {
				step.Note = sentBack
			}
			rep.Steps = append(rep.Steps, step)
			if tr.State == Failed && rep.Passed {
				rep.Passed, rep.StoppedAt = false, tr.Name
			}
		}
	} else {
		for _, s := range selected {
			step := Step{Name: s.Tier.Name, Tier: true, State: NotReached, Command: s.Argv, Dir: s.Tier.Dir}
			if s.CoveredBy != "" {
				step.State, step.CoveredBy = Skipped, s.CoveredBy
			}
			rep.Steps = append(rep.Steps, step)
		}
	}
	rep.DurationMS, rep.Duration = took(v.opts.Now().Sub(start))
	return rep, runErr
}

// rebase finds a rebase left unfinished in the worktree.
func (v *storyRun) rebase() []Finding {
	if !storygit.RebaseInProgress(v.opts.Git, v.worktree) {
		return nil
	}
	return []Finding{{Name: StepRebase, Message: "a rebase is in progress; finish it with git rebase --continue, or undo it with git rebase --abort, then verify again"}}
}

// sync finds that the story's branch lacks commits of the main branch,
// unless what they changed since the branch left it is only paths under the
// manifest's wip folder that the branch does not change; then it passes
// over them, and its note names them (ADR-0135).
func (v *storyRun) sync() []Finding {
	paths, commits, err := BaseChanges(v.opts.Git, v.worktree, v.base)
	if err != nil {
		return []Finding{{Name: StepSync, Message: fmt.Sprintf("cannot tell whether the branch contains %s: %v; check that %s is a branch of the repository", v.base, err, v.base)}}
	}
	if len(commits) == 0 {
		return nil
	}
	wip, err := v.wipPrefix()
	if err != nil {
		return []Finding{{Name: StepSync, Message: err.Error()}}
	}
	var outside, shared []string
	for _, p := range paths {
		switch {
		case !strings.HasPrefix(p, wip):
			outside = append(outside, p)
		case slices.Contains(v.paths, p):
			shared = append(shared, p)
		}
	}
	behind := fmt.Sprintf("the branch does not contain %s, %s behind it", v.base, plural(len(commits), "commit"))
	fix := fmt.Sprintf("run flai stream sync %s, resolve what it reports, and verify again", v.story)
	change := "change"
	if len(commits) == 1 {
		change = "changes"
	}
	switch {
	case len(outside) > 0:
		return []Finding{{Name: StepSync, Message: fmt.Sprintf("%s, which %s paths outside %s: %s; %s", behind, change, wip, few(outside, 3), fix)}}
	case len(shared) > 0:
		return []Finding{{Name: StepSync, Message: fmt.Sprintf("%s, which %s %s paths the branch changes too: %s; %s", behind, change, wip, few(shared, 3), fix)}}
	}
	v.note = fmt.Sprintf("passed over %s of %s that %s only %s paths the branch does not change: %s", plural(len(commits), "commit"), v.base, change, wip, few(commits, 10))
	return nil
}

// wipPrefix is the manifest's wip folder relative to the project's root,
// with a trailing slash, as git names the paths under it.
func (v *storyRun) wipPrefix() (string, error) {
	wip, err := filepath.Rel(mainRoot(v.opts.Project), v.opts.Project.WipDir())
	if err != nil {
		return "", fmt.Errorf("cannot tell which paths are under the wip folder %s: %w; check the layout in the project's manifest", v.opts.Project.WipDir(), err)
	}
	return filepath.ToSlash(wip) + "/", nil
}

// plural is n and the noun, plural but for one.
func plural(n int, noun string) string {
	if n == 1 {
		return "1 " + noun
	}
	return strconv.Itoa(n) + " " + noun + "s"
}

// few lists the first limit of items and says how many more there are.
func few(items []string, limit int) string {
	if len(items) <= limit {
		return strings.Join(items, ", ")
	}
	return fmt.Sprintf("%s and %d more", strings.Join(items[:limit], ", "), len(items)-limit)
}

// narrative finds the narrative's sections that are not written, in the
// main checkout, where the narrative lives.
func (v *storyRun) narrative() []Finding {
	path := v.opts.Project.NarrativePath(v.story)
	rel := projectPath(v.opts.Project, v.worktree, path)
	data, err := os.ReadFile(path)
	if err != nil {
		return []Finding{{Name: StepNarrative, Path: rel, Message: fmt.Sprintf("no narrative for %s: %v; open it with flai stream open %s", v.story, err, v.story)}}
	}
	found := Unwritten(string(data), narrativeSections)
	for i := range found {
		found[i].Path = rel
	}
	return found
}

// check runs flai check --strict scoped to the story, in process, on the
// worktree: a finding inside the story at error level, or at warning level
// but those --strict passes over, fails it, and one outside is put in rep's
// notes (S-0249).
func (v *storyRun) check(rep *Report) []Finding {
	repo, err := workitem.Open(v.worktree)
	if err != nil {
		return []Finding{{Name: StepCheck, Message: fmt.Sprintf("open the project at the worktree %s: %v", v.worktree, err)}}
	}
	repo.Git, repo.TemplateDir = v.opts.Git, v.opts.Project.TemplateDir
	res, err := check.Run(repo, v.opts.Now())
	if err != nil {
		return []Finding{{Name: StepCheck, Message: fmt.Sprintf("flai check could not run: %v", err)}}
	}
	// The clone may need a newer git than the one installed (ADR-0022).
	if installed, err := gitver.Installed(v.opts.Git); err == nil {
		check.GitCompat(res, repo, installed)
	}
	changed, err := storygit.StoryChanges(v.opts.Git, repo, v.story)
	if err != nil {
		return []Finding{{Name: StepCheck, Message: fmt.Sprintf("scope the check to %s: %v", v.story, err)}}
	}
	if err := check.ScopeToStory(res, repo, v.story, changed); err != nil {
		return []Finding{{Name: StepCheck, Message: err.Error()}}
	}
	var inside []Finding
	for _, f := range res.Findings {
		p := ""
		if f.Path != "" {
			p = projectPath(repo, v.worktree, f.Path)
		}
		if f.Outside {
			rep.Notes = append(rep.Notes, Note{Rule: f.Rule, Level: f.Level, Path: p, Line: f.Line, Message: f.Message})
			continue
		}
		inside = append(inside, Finding{Name: f.Rule, Path: p, Line: f.Line, Message: f.Level + ": " + f.Message})
	}
	if res.OK(true) {
		return nil
	}
	return inside
}

// tierStep is a tier's result as a step of Verify.
func tierStep(tr TierResult) Step {
	return Step{
		Name: tr.Name, Tier: true, State: tr.State, DurationMS: tr.DurationMS, Duration: tr.Duration,
		Findings: tr.Findings, Omitted: tr.Omitted, Command: tr.Command, Dir: tr.Dir, ExitCode: tr.ExitCode,
		CoveredBy: tr.CoveredBy,
	}
}

// reusable is the steps of the selected tiers that the story's last report
// lets this run reuse, in order: those before the first it shows failed or
// not reached, each Reused with the time of the run that ran it, but for a
// tier another selected tier covers, which stays Skipped. It is none with
// Fresh, with no last report, or when the report's fingerprint or tiers
// differ from this run's; every tier, when it passed them all. At the
// report's head and base commit that is all (S-0341). At another, it reuses
// only when the branch's own change, its patch id, is the same, and stops
// at the first tier that the paths the main branch brought in between the
// two merge bases select, but for wip paths the branch does not change
// (ADR-0135), as the branch's own select its tiers, or that
// covers one they select; then it answers, as the note for that tier, why
// it runs again, and the first tier it reuses says it was reused across the
// main branch's new commits (ADR-0156).
func (v *storyRun) reusable(fp string, selected []Selected) (steps []Step, sentBackNote string) {
	if v.opts.Fresh || v.baseCommit == "" || fp == "" {
		return nil, ""
	}
	last, ok, err := LastReport(v.opts.Project, v.story)
	if err != nil || !ok || last.Fingerprint != fp {
		return nil, ""
	}
	moved := last.Commit != v.commit || last.BaseCommit != v.baseCommit
	var mainPaths []string
	sentBack := map[string]bool{}
	if moved {
		if v.patch == "" || last.Patch != v.patch || v.mergeBase == "" || last.MergeBase == "" {
			return nil, ""
		}
		if last.MergeBase != v.mergeBase {
			// what main brought in cannot be told, so nothing is reused
			if mainPaths, err = BetweenPaths(v.opts.Git, v.worktree, last.MergeBase, v.mergeBase); err != nil {
				return nil, ""
			}
			if mainPaths, err = v.outsideWip(mainPaths); err != nil {
				return nil, ""
			}
		}
		sentBack = v.sentBack(selected, mainPaths)
	}
	var tiers []Step
	for _, s := range last.Steps {
		if s.Tier {
			tiers = append(tiers, s)
		}
	}
	if len(tiers) != len(selected) {
		return nil, ""
	}
	var out []Step
	for i, s := range tiers {
		if s.Name != selected[i].Tier.Name {
			break
		}
		s.Note = ""
		if sentBack[s.Name] {
			sentBackNote = v.sentBackNote(selected[i:], sentBack, mainPaths)
			break
		}
		// A tier another selected tier covers ran nothing, so it neither
		// stops the reuse nor is reused: it stays Skipped (S-0342).
		if s.State == Skipped && selected[i].CoveredBy != "" {
			s.CoveredBy = selected[i].CoveredBy
			out = append(out, s)
			continue
		}
		if s.State != Passed && s.State != Reused {
			break
		}
		from := last.RanAt.UTC()
		if s.State == Reused && s.ReusedFrom != nil {
			from = s.ReusedFrom.UTC()
		}
		s.State, s.ReusedFrom = Reused, &from
		out = append(out, s)
	}
	if moved {
		for i := range out {
			if out[i].State == Reused {
				out[i].Note = v.reusedNote(last, mainPaths)
				break
			}
		}
	}
	return out, sentBackNote
}

// outsideWip leaves out of paths the main branch brought in those under the
// wip folder that the branch does not change: no tier's result reads them
// (ADR-0135), so they send no tier back, as they hold no sync back.
func (v *storyRun) outsideWip(paths []string) ([]string, error) {
	wip, err := v.wipPrefix()
	if err != nil {
		return nil, err
	}
	var out []string
	for _, p := range paths {
		if !strings.HasPrefix(p, wip) || slices.Contains(v.paths, p) {
			out = append(out, p)
		}
	}
	return out, nil
}

// sentBack are the names of the selected tiers that the paths the main
// branch brought in reach: those the paths select as the branch's own select
// its tiers (SelectStory, {importers} included), and the tier that covers
// each of them in this run, which runs in its place (ADR-0156).
func (v *storyRun) sentBack(selected []Selected, mainPaths []string) map[string]bool {
	out := map[string]bool{}
	if len(mainPaths) == 0 {
		return out
	}
	for _, s := range SelectStory(v.opts.FS, v.opts.Tiers, mainPaths) {
		out[s.Tier.Name] = true
	}
	for _, s := range selected {
		if out[s.Tier.Name] && s.CoveredBy != "" {
			out[s.CoveredBy] = true
		}
	}
	return out
}

// sentBackNote says that the tiers from rest on run again, and which of
// them the main branch's new paths select.
func (v *storyRun) sentBackNote(rest []Selected, sentBack map[string]bool, mainPaths []string) string {
	var names []string
	for _, s := range rest {
		if sentBack[s.Tier.Name] {
			names = append(names, s.Tier.Name)
		}
	}
	return fmt.Sprintf("run again from here: the %s %s brought in since the last run %s %s: %s",
		plural(len(mainPaths), "path"), v.base, selectVerb(len(mainPaths)), strings.Join(names, ", "), few(mainPaths, 3))
}

// reusedNote says that a tier was reused across the main branch's new
// commits, and why it could be.
func (v *storyRun) reusedNote(last Report, mainPaths []string) string {
	if last.MergeBase == v.mergeBase {
		return fmt.Sprintf("reused though %s moved on: the branch's own change is the same as at %s, and it was not rebased", v.base, shortHash(last.Commit))
	}
	if len(mainPaths) == 0 {
		return fmt.Sprintf("reused across the rebase onto %s: the branch's own change is the same as at %s, and %s brought in only wip paths the branch does not change",
			v.base, shortHash(last.Commit), v.base)
	}
	return fmt.Sprintf("reused across the rebase onto %s: the branch's own change is the same as at %s, and the %s %s brought in %s none of the tiers reused",
		v.base, shortHash(last.Commit), plural(len(mainPaths), "path"), v.base, selectVerb(len(mainPaths)))
}

// selectVerb is "select" agreeing with n paths.
func selectVerb(n int) string {
	if n == 1 {
		return "selects"
	}
	return "select"
}

// shortHash is a commit's first seven characters, as git abbreviates it.
func shortHash(commit string) string {
	if len(commit) > 7 {
		return commit[:7]
	}
	return commit
}

// fingerprint hashes the project's tiers as declared and the selected ones'
// names, directories, and commands: a change to any changes what a story's
// tiers would run.
func fingerprint(tiers []Tier, selected []Selected) string {
	type sel struct {
		Name string   `json:"name"`
		Dir  string   `json:"dir"`
		Argv []string `json:"argv"`
	}
	in := struct {
		Tiers    []Tier `json:"tiers"`
		Selected []sel  `json:"selected"`
	}{Tiers: tiers, Selected: []sel{}}
	for _, s := range selected {
		in.Selected = append(in.Selected, sel{Name: s.Tier.Name, Dir: s.Tier.Dir, Argv: s.Argv})
	}
	data, err := json.Marshal(in)
	if err != nil {
		// nothing to compare, so nothing is reused
		return ""
	}
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}

// took is d as a report gives it: whole milliseconds, and as a person
// reads it.
func took(d time.Duration) (int64, string) {
	return d.Round(time.Millisecond).Milliseconds(), human(d)
}

// listMarker is what a line of a section may begin with and still be
// empty: blanks, and a bullet or a number such as "1.", as close-out.sh
// reads it.
var listMarker = regexp.MustCompile(`^\s*(?:[-*]|[0-9]+\.)?\s*`)

// Unwritten finds each of the sections, by heading text, that a narrative
// lacks or leaves empty: under its "## " heading, up to the next, no line
// holds more than blanks and a bare list marker. A finding's line is the
// heading's, 0 when it is missing; its path is left for the caller.
func Unwritten(narrative string, sections []string) []Finding {
	lines := strings.Split(strings.ReplaceAll(narrative, "\r\n", "\n"), "\n")
	var out []Finding
	for _, sec := range sections {
		heading := "## " + sec
		at, written := -1, false
		for i, line := range lines {
			if at < 0 {
				if line == heading {
					at = i
				}
				continue
			}
			if strings.HasPrefix(line, "## ") {
				break
			}
			if listMarker.ReplaceAllString(line, "") != "" {
				written = true
				break
			}
		}
		switch {
		case at < 0:
			out = append(out, Finding{Name: StepNarrative, Message: fmt.Sprintf("the narrative has no %s section; add it and write it, then verify again", heading)})
		case !written:
			out = append(out, Finding{Name: StepNarrative, Line: at + 1, Message: fmt.Sprintf("%s is empty; write it, then verify again", heading)})
		}
	}
	return out
}

// projectPath is p relative to the worktree, or, for what lives in the main
// checkout such as wip/, to it; p itself when it is below neither.
func projectPath(repo *workitem.Repo, worktree, p string) string {
	if !filepath.IsAbs(p) {
		p = filepath.Join(worktree, p)
	}
	for _, root := range []string{worktree, mainRoot(repo)} {
		if rel, err := filepath.Rel(root, p); err == nil && rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
			return filepath.ToSlash(rel)
		}
	}
	return filepath.ToSlash(p)
}

// mainRoot is the project's main checkout.
func mainRoot(repo *workitem.Repo) string {
	if repo.MainRoot != "" {
		return repo.MainRoot
	}
	return repo.Root
}

// ReportPath is where a story's last report is stored: in the project's
// cache, .flai-cache/verify/<story>.json.
func ReportPath(project *workitem.Repo, story string) string {
	return filepath.Join(project.CacheDir(), "verify", story+".json")
}

// SaveReport stores rep as its story's last report, replacing the one
// before it in one step, so that a reader sees one whole report or the
// other.
func SaveReport(project *workitem.Repo, rep Report) error {
	path := ReportPath(project, rep.Story)
	data, err := json.MarshalIndent(rep, "", "  ")
	if err != nil {
		return fmt.Errorf("store the verify report of %s: %w", rep.Story, err)
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("store the verify report of %s: make %s: %w; check that the project's .flai-cache is writable", rep.Story, filepath.Dir(path), err)
	}
	if err := atomicfile.WriteFile(path, append(data, '\n'), 0o644); err != nil {
		return fmt.Errorf("store the verify report of %s at %s: %w; check that the project's .flai-cache is writable", rep.Story, path, err)
	}
	return nil
}

// LastReport is the story's last report and true, or false with no error
// when the story has none. The story may be given as typed, S-12 for
// S-0012, when the project has it.
func LastReport(project *workitem.Repo, story string) (Report, bool, error) {
	if it, err := project.Get(story); err == nil {
		story = it.ID
	}
	path := ReportPath(project, story)
	data, err := os.ReadFile(path)
	if errors.Is(err, fs.ErrNotExist) {
		return Report{}, false, nil
	}
	if err != nil {
		return Report{}, false, fmt.Errorf("read the last verify report of %s at %s: %w", story, path, err)
	}
	var rep Report
	if err := json.Unmarshal(data, &rep); err != nil {
		return Report{}, false, fmt.Errorf("read the last verify report of %s at %s: %w; delete it and verify again", story, path, err)
	}
	return rep, true, nil
}
