package verify

import (
	"context"
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
	// RunOptions are the tiers' run's; FS is the worktree, read for the
	// tiers' packages, and Now times every step and stamps the report.
	RunOptions
}

// Report is the answer to Verify: pass, or the step that failed and why.
type Report struct {
	Story string `json:"story"`
	// Commit is the worktree's HEAD when it was verified.
	Commit string `json:"commit"`
	// Base is the main branch the story was verified against.
	Base  string    `json:"base"`
	RanAt time.Time `json:"ran_at"`
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
	// Note is what a step that passed says about what it passed over, as
	// sync names the commits of Base it did not need (ADR-0135); it is the
	// step's own, not a check note, and nothing records it as an issue.
	Note string `json:"note,omitempty"`
	// Omitted counts the findings past the run's cap.
	Omitted int `json:"omitted,omitempty"`
	// Command, Dir, and ExitCode are a tier's, as TierResult has them.
	Command  []string `json:"command,omitempty"`
	Dir      string   `json:"dir,omitempty"`
	ExitCode *int     `json:"exit_code,omitempty"`
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
// step that fails; the steps after it are not reached. It commits nothing
// and records no issue. The report is stored as the story's last (SaveReport)
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
	// commit and paths are the worktree's HEAD and what it changed.
	commit string
	paths  []string
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
	return Report{Story: v.story, Commit: v.commit, Base: v.base, RanAt: start.UTC(), Passed: true, Paths: v.paths, Steps: []Step{}}
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
	rep := v.newReport(start)
	v.runSteps(&rep, append(v.syncSteps(),
		storyStep{StepNarrative, v.narrative},
		storyStep{StepCheck, func() []Finding { return v.check(&rep) }},
	))
	selected := SelectStory(v.opts.FS, v.opts.Tiers, v.paths)
	var runErr error
	if rep.Passed {
		ro := v.opts.RunOptions
		ro.Env = append(append([]string{}, ro.Env...), StoryEnv+"="+v.story)
		var res Result
		res, runErr = Run(ctx, v.worktree, selected, ro)
		for _, tr := range res.Tiers {
			rep.Steps = append(rep.Steps, tierStep(tr))
			if tr.State == Failed && rep.Passed {
				rep.Passed, rep.StoppedAt = false, tr.Name
			}
		}
	} else {
		for _, s := range selected {
			rep.Steps = append(rep.Steps, Step{Name: s.Tier.Name, Tier: true, State: NotReached, Command: s.Argv, Dir: s.Tier.Dir})
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
	root := mainRoot(v.opts.Project)
	wip, err := filepath.Rel(root, v.opts.Project.WipDir())
	if err != nil {
		return []Finding{{Name: StepSync, Message: fmt.Sprintf("cannot tell which paths are under the wip folder %s: %v; check the layout in the project's manifest", v.opts.Project.WipDir(), err)}}
	}
	wip = filepath.ToSlash(wip) + "/"
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
	switch {
	case len(outside) > 0:
		return []Finding{{Name: StepSync, Message: fmt.Sprintf("%s, which change paths outside %s: %s; %s", behind, wip, few(outside, 3), fix)}}
	case len(shared) > 0:
		return []Finding{{Name: StepSync, Message: fmt.Sprintf("%s, which change %s paths the branch changes too: %s; %s", behind, wip, few(shared, 3), fix)}}
	}
	v.note = fmt.Sprintf("passed over %s of %s that change only %s paths the branch does not change: %s", plural(len(commits), "commit"), v.base, wip, few(commits, 10))
	return nil
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
	}
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
