// Package check validates a conforming repository against the standard:
// work items, narratives, the board, and documentation front matter. It is
// the reference validator that flaiover mirrors.
package check

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"slices"
	"sort"
	"strings"
	"time"

	"github.com/goccy/go-yaml"

	"github.com/bytepunx/system-flow/flai/internal/adrindex"
	"github.com/bytepunx/system-flow/flai/internal/analysis"
	"github.com/bytepunx/system-flow/flai/internal/components"
	"github.com/bytepunx/system-flow/flai/internal/conflictmark"
	ctxpack "github.com/bytepunx/system-flow/flai/internal/context"
	"github.com/bytepunx/system-flow/flai/internal/conventions"
	"github.com/bytepunx/system-flow/flai/internal/digest"
	"github.com/bytepunx/system-flow/flai/internal/experiment"
	"github.com/bytepunx/system-flow/flai/internal/formulas"
	"github.com/bytepunx/system-flow/flai/internal/guard"
	"github.com/bytepunx/system-flow/flai/internal/issues"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/mdlint"
	"github.com/bytepunx/system-flow/flai/internal/messages"
	"github.com/bytepunx/system-flow/flai/internal/roles"
	"github.com/bytepunx/system-flow/flai/internal/threads"
	"github.com/bytepunx/system-flow/flai/internal/topics"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// Levels.
const (
	Error   = "error"
	Warning = "warning"
)

// Finding is one rule breach.
type Finding struct {
	Level   string `json:"level"`
	Rule    string `json:"rule"`
	Path    string `json:"path"`
	Line    int    `json:"line"`
	Message string `json:"message"`
	// Outside marks a finding outside the story a run is scoped to
	// (ScopeToStory): a note the run passes over, whatever its level.
	Outside bool `json:"outside,omitempty"`
	// Stories are the two stories a wip.overlap names, lower ID first, so a
	// run scoped to a story keeps it only when it names that story (ADR-0115),
	// or the one story a touches.folder or a wip.after-hub is on, so a run
	// scoped to another story leaves it out (ADR-0188).
	Stories []string `json:"stories,omitempty"`
	// advisory marks a warning --strict passes over (advise).
	advisory bool
	// unchecked marks a finding only the host clears: a components.examples
	// note that examples were left unchecked for want of a tool on the host or
	// in the worktree, node or a sub-project's svelte (ADR-0171), or an
	// agent.image finding that the layer's image went unchecked for want of
	// docker or is not built (ADR-0255, I-0264). A check scoped to a story
	// leaves it out unless the story changes its path or one of inside.
	unchecked bool
	// inside are the paths, relative to the project root, besides Path, whose
	// change brings an unchecked finding inside a story (ADR-0255).
	inside []string
}

// Result is a whole run.
type Result struct {
	Findings []Finding `json:"findings"`
	Errors   int       `json:"errors"`
	Warnings int       `json:"warnings"`
	// Advisory counts the warnings --strict passes over, which Warnings
	// includes: the review column over its limit, which only acceptance
	// clears (S-0243). Scoped to a story, it counts only those inside it;
	// Outside counts the rest.
	Advisory int `json:"advisory"`
	// Outside counts the findings outside the story the run is scoped to,
	// errors and warnings, which Errors and Warnings include (S-0249).
	Outside int `json:"outside,omitempty"`
	Items   int `json:"items"`
	// outsideErrors and outsideWarnings split Outside by level, for OK.
	outsideErrors, outsideWarnings int
}

// OK reports whether the run passed: no errors, and when strict no warnings
// but the advisory ones, leaving out the findings outside the story the run
// is scoped to.
func (r *Result) OK(strict bool) bool {
	return r.Errors == r.outsideErrors && (!strict || r.Warnings-r.outsideWarnings == r.Advisory)
}

type checker struct {
	repo  *workitem.Repo
	now   time.Time
	items []*workitem.Item
	byID  map[string]*workitem.Item
	res   *Result
	vocab map[string]bool // topics a story can have (ADR-0047)
}

// Run executes every rule against the repo.
func Run(repo *workitem.Repo, now time.Time) (*Result, error) {
	items, err := repo.List(true)
	if err != nil {
		return nil, err
	}
	c := &checker{repo: repo, now: now.UTC(), items: items, byID: map[string]*workitem.Item{}, res: &Result{Items: len(items)}}
	c.layout()
	c.instructions()
	c.workItems()
	c.narratives()
	c.activities()
	c.overlap()
	c.componentTag()
	c.after()
	c.touchesFolder()
	c.afterHub()
	c.afterReason()
	c.unaccepted()
	c.epicLags()
	c.board()
	c.documentation()
	c.adrIndex()
	c.conventions()
	c.roles()
	c.formulas()
	c.issues()
	c.experiments()
	c.analysis()
	c.digest()
	c.designTokens()
	c.componentExamples()
	c.threads()
	c.messages()
	c.markdown()
	c.conflictMarkers()
	sortFindings(c.res.Findings)
	return c.res, nil
}

// Overlaps returns the wip.overlap findings of Run among items, in Run's
// order, without the other rules: the designer's inbox lists them (S-0158).
func Overlaps(repo *workitem.Repo, items []*workitem.Item) []Finding {
	c := &checker{repo: repo, items: items, res: &Result{}}
	c.overlap()
	sortFindings(c.res.Findings)
	return c.res.Findings
}

func sortFindings(fs []Finding) {
	sort.SliceStable(fs, func(i, j int) bool {
		a, b := fs[i], fs[j]
		if a.Path != b.Path {
			return a.Path < b.Path
		}
		return a.Line < b.Line
	})
}

func (c *checker) add(level, rule, path string, line int, format string, args ...any) {
	rel := path
	if r, err := filepath.Rel(c.repo.Root, path); err == nil {
		rel = r
	}
	c.res.Findings = append(c.res.Findings, Finding{Level: level, Rule: rule, Path: rel, Line: line, Message: fmt.Sprintf(format, args...)})
	if level == Error {
		c.res.Errors++
	} else {
		c.res.Warnings++
	}
}

// advise adds a warning that --strict passes over: one no story's agent can
// clear, which would otherwise stop every close-out (S-0243), or one every
// story passes through, such as a narrative not yet written (narrativeState).
func (c *checker) advise(rule, path string, line int, format string, args ...any) {
	c.add(Warning, rule, path, line, format, args...)
	c.res.Findings[len(c.res.Findings)-1].advisory = true
	c.res.Advisory++
}

// adviseUnchecked is advise for a components.examples note that examples on
// path were left unchecked for want of node or a sub-project's svelte, marked
// unchecked (ADR-0171).
func (c *checker) adviseUnchecked(path, format string, args ...any) {
	c.advise("components.examples", path, 1, format, args...)
	c.res.Findings[len(c.res.Findings)-1].unchecked = true
}

// keyLine returns the 1-based line of a top-level front matter key.
func keyLine(path, key string) int {
	f, err := os.Open(path)
	if err != nil {
		return 1
	}
	defer func() { _ = f.Close() }()
	sc := bufio.NewScanner(f)
	n := 0
	for sc.Scan() {
		n++
		if n > 1 && sc.Text() == "---" {
			break
		}
		if strings.HasPrefix(sc.Text(), key+":") {
			return n
		}
	}
	return 1
}

// headingLine returns the line of a "## Heading" in the body, or 1.
func headingLine(path, heading string) int {
	f, err := os.Open(path)
	if err != nil {
		return 1
	}
	defer func() { _ = f.Close() }()
	sc := bufio.NewScanner(f)
	n := 0
	for sc.Scan() {
		n++
		if strings.TrimSpace(sc.Text()) == heading {
			return n
		}
	}
	return 1
}

func (c *checker) layout() {
	m := c.repo.Manifest
	// API responses name the project by this key (ADR-0024). A warning now;
	// an error after the template minor that ships the requirement.
	if strings.TrimSpace(m.Key) == "" {
		mf := filepath.Join(c.repo.Root, "system-flow.yaml")
		c.add(Warning, "manifest.key", mf, keyLine(mf, "name"), "system-flow.yaml has no key; add a short one (for example the project's initials), which the dashboard's API and a hub use to name this project")
	}
	if err := m.Agent.Validate(); err != nil {
		mf := filepath.Join(c.repo.Root, "system-flow.yaml")
		c.add(Error, "manifest.agent", mf, keyLine(mf, "agent"), "%s; flai agent set replaces it", err)
	}
	if _, err := m.Issues.StoryAfterDuration(); err != nil {
		mf := filepath.Join(c.repo.Root, "system-flow.yaml")
		c.add(Error, "manifest.issues", mf, keyLine(mf, "issues"), "%v", err)
	}
	for _, msg := range m.Planning.Errors() {
		mf := filepath.Join(c.repo.Root, "system-flow.yaml")
		c.add(Error, "manifest.planning", mf, keyLine(mf, "planning"), "%s", msg)
	}
	for _, msg := range m.Orchestration.Errors() {
		mf := filepath.Join(c.repo.Root, "system-flow.yaml")
		c.add(Error, "manifest.orchestration", mf, keyLine(mf, "orchestration"), "%s", msg)
	}
	for _, msg := range m.RoleErrors() {
		mf := filepath.Join(c.repo.Root, "system-flow.yaml")
		c.add(Error, "manifest.roles", mf, keyLine(mf, "roles"), "%s", msg)
	}
	for _, msg := range m.Analysis.Errors() {
		mf := filepath.Join(c.repo.Root, "system-flow.yaml")
		c.add(Error, "manifest.analysis", mf, keyLine(mf, "analysis"), "%s", msg)
	}
	for _, msg := range m.Digest.Errors() {
		mf := filepath.Join(c.repo.Root, "system-flow.yaml")
		c.add(Error, "manifest.digest", mf, keyLine(mf, "digest"), "%s", msg)
	}
	for _, msg := range m.Threads.Errors() {
		mf := filepath.Join(c.repo.Root, "system-flow.yaml")
		c.add(Error, "manifest.threads", mf, keyLine(mf, "threads"), "%s", msg)
	}
	for _, msg := range m.TestErrors() {
		mf := filepath.Join(c.repo.Root, "system-flow.yaml")
		c.add(Error, "manifest.tests", mf, keyLine(mf, "tests"), "%s", msg)
	}
	for _, msg := range m.TestWarnings() {
		mf := filepath.Join(c.repo.Root, "system-flow.yaml")
		c.add(Warning, "manifest.tests-timeout", mf, keyLine(mf, "tests"), "%s", msg)
	}
	for _, msg := range m.ProjectErrors(os.DirFS(c.repo.Root)) {
		mf := filepath.Join(c.repo.Root, "system-flow.yaml")
		c.add(Error, "manifest.projects", mf, keyLine(mf, "projects"), "%s", msg)
	}
	for _, msg := range m.ComponentErrors(os.DirFS(c.repo.Root)) {
		mf := filepath.Join(c.repo.Root, "system-flow.yaml")
		c.add(Error, "manifest.components", mf, keyLine(mf, "projects"), "%s", msg)
	}
	for _, msg := range m.InstallErrors() {
		mf := filepath.Join(c.repo.Root, "system-flow.yaml")
		c.add(Error, "manifest.install", mf, keyLine(mf, "projects"), "%s", msg)
	}
	for _, msg := range m.ProviderErrors() {
		mf := filepath.Join(c.repo.Root, "system-flow.yaml")
		c.add(Error, "manifest.providers", mf, keyLine(mf, "providers"), "%s", msg)
	}
	// the kluster: block (S-0577, ADR-0254): each finding names its key
	for _, msg := range m.KlusterErrors(os.DirFS(c.repo.Root)) {
		mf := filepath.Join(c.repo.Root, "system-flow.yaml")
		c.add(Error, "manifest.kluster", mf, keyLine(mf, "kluster"), "%s", msg)
	}
	c.sharedPatterns()
	if m.Template.Version != "" {
		if _, err := os.Stat(filepath.Join(c.repo.Root, "system-flow.lock.yaml")); err != nil {
			c.add(Warning, "layout.lock", filepath.Join(c.repo.Root, "system-flow.yaml"), keyLine(filepath.Join(c.repo.Root, "system-flow.yaml"), "template"), "template %s is recorded but there is no system-flow.lock.yaml; run flai upgrade --relock", m.Template.Version)
		}
	}
	for key, subs := range map[string][]string{
		"design": {"adrs", "system", "tech", "conventions"},
		"docs":   nil,
		"wip":    {"kanban", "agents", "archive"},
	} {
		dir := m.Dir(c.repo.Root, key)
		if st, err := os.Stat(dir); err != nil || !st.IsDir() {
			c.add(Error, "layout.folder", c.repo.Root, 1, "layout.%s folder %q does not exist", key, m.Layout[key])
			continue
		}
		for _, sub := range subs {
			if st, err := os.Stat(filepath.Join(dir, sub)); err != nil || !st.IsDir() {
				c.add(Warning, "layout.subfolder", dir, 1, "%s/%s is missing", m.Layout[key], sub)
			}
		}
	}
	for _, sub := range []string{"epics", "stories", "tasks"} {
		if st, err := os.Stat(filepath.Join(c.repo.KanbanDir(), sub)); err != nil || !st.IsDir() {
			c.add(Warning, "layout.subfolder", c.repo.KanbanDir(), 1, "kanban/%s is missing", sub)
		}
	}
	c.cacheIgnored()
}

// sharedPatterns reports each pattern of claims.shared that is not valid, on
// its line of the manifest: it frees nothing (ADR-0096).
func (c *checker) sharedPatterns() {
	errs := c.repo.SharedClaims().Errors() // as the file has them, like the lines
	if len(errs) == 0 {
		return
	}
	mf := filepath.Join(c.repo.Root, manifest.File)
	data, _ := os.ReadFile(mf) // without it, each is reported on line 1
	lines := manifest.SharedLines(data)
	for _, e := range errs {
		line := 1
		if e.Index < len(lines) {
			line = lines[e.Index]
		}
		c.add(Error, "manifest.claims", mf, line, "%s; until it is fixed it frees nothing", e)
	}
}

// cacheIgnored refuses a repository whose .flai-cache (dashboard token,
// worktrees, tool caches) could be committed (ADR-0018).
func (c *checker) cacheIgnored() {
	cache := filepath.Join(c.repo.Root, ".flai-cache")
	if st, err := os.Stat(cache); err != nil || !st.IsDir() {
		return
	}
	gi := filepath.Join(c.repo.Root, ".gitignore")
	data, err := os.ReadFile(gi)
	if err != nil {
		c.add(Error, "layout.gitignore", c.repo.Root, 1, ".flai-cache/ exists but there is no .gitignore; it holds the dashboard token and must be ignored")
		return
	}
	for _, line := range strings.Split(string(data), "\n") {
		switch strings.TrimSpace(line) {
		case ".flai-cache", ".flai-cache/", "/.flai-cache", "/.flai-cache/":
			return
		}
	}
	c.add(Error, "layout.gitignore", gi, 1, ".flai-cache/ is not ignored; it holds the dashboard token and must be")
}

var requiredHeadings = map[string][]string{
	workitem.Epic:  {"## Outcome", "## Stories", "## Notes"},
	workitem.Story: {"## Goal", "## Acceptance criteria", "## Tasks", "## Notes"},
	workitem.Task:  {"## Work", "## Done when", "## Notes"},
}

func (c *checker) workItems() {
	seen := map[string]string{}
	for _, it := range c.items {
		if prev, dup := seen[it.ID]; dup {
			c.add(Error, "item.duplicate-id", it.Path, keyLine(it.Path, "id"), "%s is also defined in %s", it.ID, prev)
		} else {
			seen[it.ID] = it.Path
			c.byID[it.ID] = it
		}
	}
	for _, it := range c.items {
		c.oneItem(it)
	}
}

func (c *checker) oneItem(it *workitem.Item) {
	p := it.Path
	var problems workitem.Problems
	if err := it.Validate(); errors.As(err, &problems) {
		for _, msg := range problems {
			key := strings.Fields(msg)[0]
			// transitions[0].at and cost_of_delay.value are on their block's key
			key = strings.SplitN(strings.SplitN(key, "[", 2)[0], ".", 2)[0]
			c.add(Error, "item.front-matter", p, keyLine(p, key), "%s", msg)
		}
	}
	c.unknownFields("item.unknown-field", p, it.Unknown)
	base := filepath.Base(p)
	if !strings.HasPrefix(base, it.ID+"-") && base != it.ID+".md" {
		c.add(Error, "item.filename", p, 1, "file name should start with %s-", it.ID)
	}
	if it.Parent != "" {
		parent, ok := c.byID[it.Parent]
		switch {
		case !ok:
			c.add(Error, "item.parent-missing", p, keyLine(p, "parent"), "parent %s does not exist", it.Parent)
		case parent.Type != map[string]string{workitem.Story: workitem.Epic, workitem.Task: workitem.Story}[it.Type]:
			c.add(Error, "item.parent-type", p, keyLine(p, "parent"), "parent %s is a %s", it.Parent, parent.Type)
		default:
			// Cancelling a parent cancels what is open under it (ADR-0028); a tree
			// cancelled by an older flai, or edited by hand, is caught here.
			if parent.Status == workitem.Cancelled && !it.Closed() {
				c.add(Error, "item.parent-cancelled", p, keyLine(p, "status"), "%s is %s under %s, which is cancelled; cancel it (flai move %s cancelled --reason \"%s cancelled\") or give it another parent", it.ID, it.Status, parent.ID, it.ID, parent.ID)
			}
			if !strings.Contains(parent.Body, "- "+it.ID+" ") && !strings.Contains(parent.Body, it.ID+"\n") {
				c.add(Warning, "item.parent-list", parent.Path, headingLine(parent.Path, map[string]string{workitem.Story: "## Stories", workitem.Task: "## Tasks"}[it.Type]), "%s does not list child %s", parent.ID, it.ID)
			}
		}
	}
	if it.Type == workitem.Task && it.Stream != "" && it.Stream != it.Parent {
		c.add(Warning, "item.stream", p, keyLine(p, "stream"), "stream %s differs from parent %s", it.Stream, it.Parent)
	}
	if it.Archived && !it.Closed() {
		c.add(Error, "item.archived-open", p, keyLine(p, "status"), "archived item is %s; only done or cancelled items belong in archive", it.Status)
	}
	if !it.Archived && it.Closed() && it.Type != workitem.Task {
		c.add(Warning, "item.archive", p, keyLine(p, "status"), "%s is %s; run flai archive", it.ID, it.Status)
	}
	c.history(it)
	children := workitem.Children(c.items, it.ID)
	// Tasks are written by the agent that pulls the story, once it is in
	// progress, and are required from review onwards (ADR-0021).
	if it.Type == workitem.Story && (it.Status == workitem.Review || it.Status == workitem.Done) && len(children) == 0 {
		c.add(Error, "story.tasks", p, keyLine(p, "status"), "a %s story needs at least one task", it.Status)
	}
	if it.Type == workitem.Story && it.Status != workitem.Backlog && it.Status != workitem.Cancelled && !hasCriteria(it.Body) {
		c.add(Error, "story.criteria", p, headingLine(p, "## Acceptance criteria"), "a %s story needs acceptance criteria with at least one checkbox", it.Status)
	}
	// A draft is finalized before it is ready (S-0199); one past backlog was
	// moved by an older flai or edited by hand. The archive is not edited and
	// is not reported.
	if it.Type == workitem.Story && it.Draft && !it.Archived && it.Status != workitem.Backlog && it.Status != workitem.Cancelled {
		c.add(Warning, "story.draft", p, keyLine(p, "draft"), "%s is a draft in %s; finalize it (flai edit %s --no-draft) or move it back to backlog", it.ID, it.Status, it.ID)
	}
	if it.Status == workitem.Done {
		for _, ch := range children {
			if !ch.Closed() {
				c.add(Error, "item.done-children", p, keyLine(p, "status"), "%s is done but %s is %s", it.ID, ch.ID, ch.Status)
			}
		}
		if it.Type == workitem.Story && hasUnchecked(it.Body) {
			c.add(Error, "story.unchecked", p, headingLine(p, "## Acceptance criteria"), "done story has unchecked acceptance criteria")
		}
	}
	if it.Closed() && it.IsBlocked() {
		c.add(Warning, "item.blocked-closed", p, keyLine(p, "blocked"), "%s item still has an open blocked interval", it.Status)
	}
	for _, h := range requiredHeadings[it.Type] {
		if !strings.Contains(it.Body, "\n"+h+"\n") && !strings.Contains(it.Body, "\n"+h+" ") {
			c.add(Warning, "item.heading", p, 1, "body is missing the %q section", h)
		}
	}
}

func (c *checker) history(it *workitem.Item) {
	p := it.Path
	created, errC := time.Parse(workitem.TimeFormat, it.Created)
	updated, errU := time.Parse(workitem.TimeFormat, it.Updated)
	prevState := workitem.Backlog
	prevAt := created
	for i, tr := range it.Transitions {
		at, err := time.Parse(workitem.TimeFormat, tr.At)
		if err != nil {
			continue // reported by Validate
		}
		if errC == nil && at.Before(prevAt) {
			c.add(Error, "item.chronology", p, keyLine(p, "transitions"), "transitions[%d] at %s is before %s", i, tr.At, prevAt.Format(workitem.TimeFormat))
		}
		prevAt = at
		// A parent's cancellation takes an item out of review; nothing else does (ADR-0028).
		cascaded := tr.To == workitem.Cancelled && prevState == workitem.Review && workitem.CancelledWith(c.byID, it, tr.At) != ""
		// The transitions flai move allows, back moves included (ADR-0055), from the one table.
		if tr.To == workitem.Done && prevState == workitem.InProgress && it.Type != workitem.Task {
			c.add(Error, "item.sequence", p, keyLine(p, "transitions"), "transitions[%d]: a %s must go through review before done", i, it.Type)
		} else if !workitem.Follows(it.Type, prevState, tr.To) && !cascaded {
			c.add(Error, "item.sequence", p, keyLine(p, "transitions"), "transitions[%d]: %s cannot follow %s", i, tr.To, prevState)
		}
		prevState = tr.To
	}
	if errC == nil && errU == nil {
		if updated.Before(created) {
			c.add(Error, "item.chronology", p, keyLine(p, "updated"), "updated is before created")
		}
		if len(it.Transitions) > 0 && updated.Before(prevAt) {
			c.add(Error, "item.chronology", p, keyLine(p, "updated"), "updated is before the last transition")
		}
	}
}

// unaccepted flags stories that are done without having been accepted: still
// in kanban, or with their story branch still present (S-0046). Done means
// accepted; flai accept <id> completes them.
func (c *checker) unaccepted() {
	for _, it := range c.items {
		if it.Type != workitem.Story || it.Status != workitem.Done {
			continue
		}
		switch {
		case !it.Archived:
			c.add(Warning, "story.unaccepted", it.Path, keyLine(it.Path, "status"), "%s is done but was never accepted (not archived, not released); run flai accept %s", it.ID, it.ID)
		case storyBranchExists(c.repo.MainRoot, it.ID):
			c.add(Warning, "story.unaccepted", it.Path, keyLine(it.Path, "status"), "%s is done but its branch story/%s was never merged; merge or delete it", it.ID, it.ID)
		}
	}
}

// boardRank orders the states along the board; cancelled has no place in it.
var boardRank = map[string]int{workitem.Backlog: 0, workitem.Ready: 1, workitem.InProgress: 2, workitem.Review: 3, workitem.Done: 4}

// epicLags warns about each open epic that its stories put further along the
// board than its status (S-0200), naming the moves that catch it up. It is
// advisory: the epic is the operator's to move, not a story's agent's, so it
// must not stop every close-out. An epic ahead of its stories is not warned.
func (c *checker) epicLags() {
	for _, it := range c.items {
		if it.Type != workitem.Epic || it.Archived || it.Closed() {
			continue
		}
		target, ok := workitem.EpicFollows(c.items, it.ID)
		if !ok || boardRank[target] <= boardRank[it.Status] {
			continue
		}
		why := "an epic follows its stories since S-0200, and this one was moved before that, or by hand"
		if target == workitem.Done {
			c.advise("epic.lags-stories", it.Path, keyLine(it.Path, "status"), "%s is %s but its stories are all done; accept it with flai accept %s, which moves it to done and archives it (%s)", it.ID, it.Status, it.ID, why)
			continue
		}
		var moves []string
		for _, st := range []string{workitem.Ready, workitem.InProgress, workitem.Review} {
			if boardRank[st] > boardRank[it.Status] && boardRank[st] <= boardRank[target] {
				moves = append(moves, st)
			}
		}
		moves[0] = "flai move " + it.ID + " " + moves[0]
		c.advise("epic.lags-stories", it.Path, keyLine(it.Path, "status"), "%s is %s but its stories put it in %s; catch it up with %s (%s)", it.ID, it.Status, target, strings.Join(moves, ", then "), why)
	}
}

// storyBranchExists looks for refs/heads/story/<id> without shelling out.
func storyBranchExists(mainRoot, id string) bool {
	if mainRoot == "" {
		return false
	}
	if _, err := os.Stat(filepath.Join(mainRoot, ".git", "refs", "heads", "story", id)); err == nil {
		return true
	}
	data, err := os.ReadFile(filepath.Join(mainRoot, ".git", "packed-refs"))
	return err == nil && strings.Contains(string(data), " refs/heads/story/"+id+"\n")
}

// overlap warns when the claims of two stories in progress cover the same
// path (ADR-0019); it is advisory, humans and agents coordinate. Each story
// is compared by its claim, as the pull hold reads it (ADR-0096 §3): its
// touches with each folder narrowed to its tasks' touches inside it, and its
// open tasks' touches outside them; a task is never compared on its own. An
// overlap that lies wholly inside a shared path is not reported (ADR-0096).
// A pair of stories is one warning, on the item file of the one with the
// lower ID, naming each overlapping entry of both claims, with both IDs in
// Stories (ADR-0115).
func (c *checker) overlap() {
	holds := workitem.NewHolds(c.items, c.repo.Manifest.Projects).WithShared(c.repo.SharedClaims())
	type open struct {
		it    *workitem.Item
		claim []string
	}
	var stories []open
	for _, it := range c.items {
		if !it.Archived && it.Type == workitem.Story && it.Status == workitem.InProgress {
			stories = append(stories, open{it, holds.Claim(it)})
		}
	}
	sort.SliceStable(stories, func(i, j int) bool { return stories[i].it.ID < stories[j].it.ID })
	for i, a := range stories {
		for _, b := range stories[i+1:] {
			mine, theirs := overlapping(holds, a.claim, b.claim), overlapping(holds, b.claim, a.claim)
			if len(mine) == 0 {
				continue
			}
			// advisory since ADR-0230: a ready story that overlaps a story in
			// progress starts on a clean trial merge, so two stories in progress
			// overlap by design, and the trial merge at sync settles what
			// conflicts between them; --strict on the main branch must not fail
			// on the state the board is in
			c.advise("wip.overlap", a.it.Path, keyLine(a.it.Path, "touches"), "%s touches %s, which %s (in progress) also touches as %s", a.it.ID, andList(mine), b.it.ID, andList(theirs))
			c.res.Findings[len(c.res.Findings)-1].Stories = []string{a.it.ID, b.it.ID}
		}
	}
}

// overlapping is the entries of claim that overlap an entry of other where it
// counts, as holds.Overlaps says, in claim's order.
func overlapping(holds *workitem.Holds, claim, other []string) []string {
	var out []string
	for _, p := range claim {
		if slices.ContainsFunc(other, func(o string) bool { return holds.Overlaps(p, o) }) {
			out = append(out, p)
		}
	}
	return out
}

// andList joins xs with ", " and " and " before the last.
func andList(xs []string) string {
	if len(xs) < 2 {
		return strings.Join(xs, "")
	}
	return strings.Join(xs[:len(xs)-1], ", ") + " and " + xs[len(xs)-1]
}

// componentTag warns on an open story that a release could not plan: its
// claim, its touches with each folder narrowed to the touches its tasks name
// inside it, done ones included, and its open tasks' touches outside them
// (ADR-0096), reaches two or more components, and no tag of
// its own or its epic's names one of them, so nothing says which it delivers
// to (I-0024). Research and experiments cut no release (ADR-0025,
// ADR-0066), so neither is warned.
func (c *checker) componentTag() {
	projects := c.repo.Manifest.Projects
	if len(projects) < 2 {
		return
	}
	holds := workitem.NewHolds(c.items, projects)
	for _, it := range c.items {
		if it.Archived || it.Type != workitem.Story || it.Closed() || it.Nature == "research" || it.Nature == "experiment" {
			continue
		}
		claim := holds.Claim(it)
		var reached []manifest.Project
		for _, p := range projects {
			for _, path := range claim {
				if workitem.PathsOverlap(path, p.Path) {
					reached = append(reached, p)
					break
				}
			}
		}
		if len(reached) < 2 {
			continue
		}
		tags := append([]string{}, it.Tags...)
		if parent := c.byID[it.Parent]; parent != nil {
			tags = append(tags, parent.Tags...)
		}
		named := false
		var names []string
		for _, p := range reached {
			names = append(names, p.Name)
			for _, t := range tags {
				named = named || t == p.Name || slices.Contains(p.Tags, t)
			}
		}
		if named {
			continue
		}
		// --tag replaces the tags, so the suggestion keeps the ones there
		suggest := append([]string{"<" + strings.Join(names, "|") + ">"}, it.Tags...)
		c.add(Warning, "story.component-tag", it.Path, keyLine(it.Path, "touches"), "%s touches %s, and no tag of its own or its epic's names one of them, so its release cannot tell which it delivers to; tag it with the one it delivers to: flai edit %s --tag %s", it.ID, strings.Join(names, " and "), it.ID, strings.Join(suggest, ","))
	}
}

// named is the item id names, in any zero padding, or one stored under a
// legacy three-digit ID as written; nil when there is none.
func (c *checker) named(id string) *workitem.Item {
	if it, ok := c.byID[workitem.CanonicalID(id)]; ok {
		return it
	}
	return c.byID[id]
}

// after reports an after: entry that names no story, the story itself, or
// a story that waits, through after:, for this one: each would hold the
// story until someone edits it (S-0130, ADR-0046). A task's after: names
// tasks of its own story, and is reported the same way under task.after,
// with an entry that names a task of another story (S-0176). The archive is
// not edited and is not reported.
func (c *checker) after() {
	for _, it := range c.items {
		if it.Archived || (it.Type != workitem.Story && it.Type != workitem.Task) {
			continue
		}
		p, rule := it.Path, it.Type+".after"
		for _, e := range it.After {
			switch other := c.named(e); {
			case other == nil:
				c.add(Error, rule, p, keyLine(p, "after"), "after names %s, which does not exist; fix it or clear it (flai edit %s --clear-after)", e, it.ID)
			case other.ID == it.ID:
				c.add(Error, rule, p, keyLine(p, "after"), "after names %s itself; a %s cannot wait for itself", it.ID, it.Type)
			case other.Type != it.Type:
				c.add(Error, rule, p, keyLine(p, "after"), "after names %s, a %s; a %s's after names %s", e, other.Type, it.Type, plural(it.Type))
			case it.Type == workitem.Task && other.Parent != it.Parent:
				c.add(Error, rule, p, keyLine(p, "after"), "after names %s, a task of %s; a task waits only for tasks of its own story, %s", other.ID, other.Parent, it.Parent)
			}
		}
		if cycle := c.afterCycle(it, c.named); cycle != nil {
			c.add(Error, rule, p, keyLine(p, "after"), "after forms a cycle, %s, so none of them would start; drop one of the entries", strings.Join(cycle, " waits for "))
		}
	}
}

// touchesFolder gives an advisory warning on each story in progress whose
// claim, its touches narrowed by its tasks' (ADR-0096), keeps a folder touch:
// an entry FolderTouch finds a folder of the checkout holding files, outside
// the shared paths (S-0468, ADR-0188). It names the folder, its file count,
// and the ready stories the folder holds (readyHeld), with the story in
// Stories, so a run scoped to another story leaves it out. It is advisory: a
// story keeps its folder touch until its tasks name files, and a planner may
// keep one with its note, so it informs rather than stopping a close-out.
func (c *checker) touchesFolder() {
	shared := c.repo.SharedClaims()
	var holds *workitem.Holds
	var ready []readyClaim
	for _, it := range c.items {
		if it.Archived || it.Type != workitem.Story || it.Status != workitem.InProgress {
			continue
		}
		if holds == nil {
			holds = c.repo.Holds(c.items)
			ready = readyClaims(holds, c.items)
		}
		for _, p := range holds.Claim(it) {
			n, ok := workitem.FolderTouch(c.repo.Root, c.repo.Manifest.Projects, shared, p)
			if !ok {
				continue
			}
			held := "and holds no ready story now"
			if ids := readyHeld(holds, ready, it.ID, p); len(ids) > 0 {
				held = "and holds the ready " + oneOrMany(len(ids), "story", "stories") + " " + andList(ids) + " through it"
			}
			c.advise("touches.folder", it.Path, keyLine(it.Path, "touches"), "%s's claim keeps the folder %s, which holds %d %s, %s; name the files its tasks change with flai touches, or say under ### Planning in %s's Notes why files may be added there (ADR-0188)",
				it.ID, p, n, oneOrMany(n, "file", "files"), held, it.ID)
			c.res.Findings[len(c.res.Findings)-1].Stories = []string{it.ID}
		}
	}
}

// readyClaim is a ready story held by something, as touchesFolder reads it:
// the stories that hold it on overlap alone with the paths they hold it on,
// which no share clears, or, held otherwise, its claim.
type readyClaim struct {
	id      string
	overlap []workitem.Overlap
	claim   []string
}

// readyClaims is the ready stories among items that holds holds, in ID
// order.
func readyClaims(holds *workitem.Holds, items []*workitem.Item) []readyClaim {
	var out []readyClaim
	for _, it := range items {
		if it.Archived || it.Type != workitem.Story || it.Status != workitem.Ready {
			continue
		}
		if over := holds.OverlapsBy(it); over != nil {
			out = append(out, readyClaim{id: it.ID, overlap: over})
		} else if holds.Of(it) != nil {
			out = append(out, readyClaim{id: it.ID, claim: holds.Claim(it)})
		}
	}
	sort.SliceStable(out, func(i, j int) bool { return out[i].id < out[j].id })
	return out
}

// readyHeld is the IDs of the ready stories the folder, an entry of holder's
// claim, holds: one held on overlap alone when holder holds it on a path in
// the folder, and one held otherwise, on after: or by a story with no
// touches, when its claim overlaps the folder where it counts.
func readyHeld(holds *workitem.Holds, ready []readyClaim, holder, folder string) []string {
	var out []string
	for _, r := range ready {
		in := slices.ContainsFunc(r.overlap, func(o workitem.Overlap) bool {
			return o.By == holder && slices.ContainsFunc(o.Paths, func(p string) bool { return workitem.PathsOverlap(p, folder) })
		})
		in = in || slices.ContainsFunc(r.claim, func(p string) bool { return holds.Overlaps(p, folder) })
		if in {
			out = append(out, r.id)
		}
	}
	return out
}

// oneOrMany is one when n is 1, else many.
func oneOrMany(n int, one, many string) string {
	if n == 1 {
		return one
	}
	return many
}

// afterHubAt is how many open stories must name a story in after: for it to
// be a hub (ADR-0188).
const afterHubAt = 3

// afterHub gives an advisory warning on each story not done or cancelled
// that afterHubAt or more stories in backlog, ready, or in progress name in
// after:, naming them, with the hub in Stories, so a run scoped to another
// story leaves it out (S-0468, ADR-0188). An archived story is done or
// cancelled, so it is neither a hub nor one that waits. It is advisory: only
// the stories that name the hub can drop it, which its own agent cannot.
func (c *checker) afterHub() {
	waiting := map[string][]string{}
	for _, it := range c.items {
		if it.Archived || it.Type != workitem.Story || (it.Status != workitem.Backlog && it.Status != workitem.Ready && it.Status != workitem.InProgress) {
			continue
		}
		seen := map[string]bool{}
		for _, e := range it.After {
			hub := c.named(e)
			if hub == nil || hub.Type != workitem.Story || hub.ID == it.ID || seen[hub.ID] {
				continue
			}
			seen[hub.ID] = true
			waiting[hub.ID] = append(waiting[hub.ID], it.ID)
		}
	}
	hubs := make([]string, 0, len(waiting))
	for id, by := range waiting {
		if len(by) >= afterHubAt {
			hubs = append(hubs, id)
		}
	}
	slices.Sort(hubs)
	for _, id := range hubs {
		hub := c.byID[id]
		if hub.Archived || hub.Closed() {
			continue
		}
		by := waiting[id]
		slices.Sort(by)
		c.advise("wip.after-hub", hub.Path, 1, "%s is named in after by %d open stories, %s, and each of them waits until it is accepted; after names only a story whose accepted change the story needs on main, and the pull order sets an order of preference, so drop %s from the after of each that does not need its change (flai edit <id> --after) (ADR-0188)",
			id, len(by), andList(by), id)
		c.res.Findings[len(c.res.Findings)-1].Stories = []string{id}
	}
}

// storyIDWord is a story ID as a whole word, in any form CanonicalID
// accepts: S-42, S-0042, s42.
var storyIDWord = regexp.MustCompile(`(?i)\bS-?\d+\b`)

// afterReason gives an advisory warning on each open story for each story
// in its after: that its "## Notes" section does not name, as a whole word
// in any padding, with the story in Stories, so a run scoped to another
// story leaves it out (S-0551, the after rule). An entry other rules report,
// one naming nothing, itself, or no story, is passed over, and so is a
// task's after:. It is advisory, as wip.after-hub is.
func (c *checker) afterReason() {
	for _, it := range c.items {
		if it.Archived || it.Type != workitem.Story || it.Closed() || len(it.After) == 0 {
			continue
		}
		named := map[string]bool{}
		if notes, ok := workitem.NarrativeSection(it.Body, "Notes"); ok {
			for _, w := range storyIDWord.FindAllString(notes.Text, -1) {
				named[workitem.CanonicalID(w)] = true
			}
		}
		seen := map[string]bool{}
		for _, e := range it.After {
			other := c.named(e)
			if other == nil || other.Type != workitem.Story || other.ID == it.ID || seen[other.ID] || named[other.ID] {
				continue
			}
			seen[other.ID] = true
			c.advise("wip.after-reason", it.Path, keyLine(it.Path, "after"), "%s names %s in after with no reason in its Notes: say there what %s reads from %s, or drop it from after (the after rule)",
				it.ID, other.ID, it.ID, other.ID)
			c.res.Findings[len(c.res.Findings)-1].Stories = []string{it.ID}
		}
	}
}

// plural is an item type's plural: stories, tasks.
func plural(typ string) string {
	if typ == workitem.Story {
		return "stories"
	}
	return typ + "s"
}

// afterCycle is the path from item back to itself through after:, when
// there is one and item has the lowest ID on it, so that a cycle is
// reported once. It follows only entries that after: may name: stories from
// a story, tasks of the same story from a task.
func (c *checker) afterCycle(item *workitem.Item, named func(string) *workitem.Item) []string {
	seen := map[string]bool{}
	var walk func(it *workitem.Item, path []string) []string
	walk = func(it *workitem.Item, path []string) []string {
		for _, e := range it.After {
			next := named(e)
			if next == nil || next.Type != item.Type || (item.Type == workitem.Task && next.Parent != item.Parent) || next.ID == it.ID {
				continue
			}
			if next.ID == item.ID {
				return append(path, next.ID)
			}
			if seen[next.ID] || next.ID < item.ID {
				continue
			}
			seen[next.ID] = true
			if found := walk(next, append(path, next.ID)); found != nil {
				return found
			}
		}
		return nil
	}
	return walk(item, []string{item.ID})
}

func (c *checker) narratives() {
	dir := c.repo.AgentsDir()
	active := map[string]bool{}
	for _, it := range c.items {
		if it.Type == workitem.Story && !it.Archived && (it.Status == workitem.InProgress || it.Status == workitem.Review) {
			if _, err := os.Stat(c.repo.NarrativePath(it.ID)); err != nil {
				c.add(Error, "narrative.missing", it.Path, keyLine(it.Path, "status"), "%s is %s but has no narrative; run flai stream open %s", it.ID, it.Status, it.ID)
			}
		}
	}
	matches, _ := filepath.Glob(filepath.Join(dir, "S-*.md"))
	for _, path := range matches {
		want := strings.TrimSuffix(filepath.Base(path), ".md")
		n, err := workitem.ReadNarrative(path)
		if err != nil {
			// The index lists it as unreadable (I-0129), so this is its one
			// finding, not an index that lists a stream with no narrative too.
			active[want] = true
			c.add(Error, "narrative.front-matter", path, 1, "%v", err)
			continue
		}
		if n.Stream != want {
			c.add(Error, "narrative.stream", path, keyLine(path, "stream"), "stream %q does not match file name %s", n.Stream, want)
		}
		active[want] = true
		story, ok := c.byID[want]
		switch {
		case !ok:
			c.add(Warning, "narrative.orphan", path, 1, "no story %s exists", want)
		case story.Archived:
			c.add(Warning, "narrative.stale", path, 1, "story %s is archived; move this narrative to archive/agents", want)
		}
		if _, err := time.Parse(workitem.TimeFormat, n.Updated); err != nil {
			c.add(Warning, "narrative.updated", path, keyLine(path, "updated"), "updated %q is not a UTC timestamp", n.Updated)
		}
		for _, sec := range workitem.NarrativeSections {
			if !strings.Contains(n.Body, "\n## "+sec+"\n") {
				c.add(Warning, "narrative.section", path, 1, "missing %q section", "## "+sec)
			}
		}
		if ok && !story.Archived && (story.Status == workitem.InProgress || story.Status == workitem.Review) {
			c.narrativeState(path, story)
		}
	}
	indexPath := filepath.Join(dir, "index.md")
	data, err := os.ReadFile(indexPath)
	if err != nil {
		if len(active) > 0 {
			c.add(Warning, "narrative.index", dir, 1, "index.md is missing; any flai command that touches items regenerates it")
		}
		return
	}
	for id := range active {
		if !strings.Contains(string(data), "["+id+"]") {
			c.add(Warning, "narrative.index", indexPath, 1, "index.md does not list %s; run flai stream log or flai move to regenerate", id)
		}
	}
	for _, id := range regexp.MustCompile(`\[(S-\d+)\]\(`).FindAllStringSubmatch(string(data), -1) {
		if !active[id[1]] {
			c.add(Warning, "narrative.index", indexPath, 1, "index.md lists %s which has no narrative in agents/", id[1])
		}
	}
}

// listItem is a line that opens a list item with something in it: a bullet
// or a number such as "1.", then text.
var listItem = regexp.MustCompile(`^\s*(?:[-*]|[0-9]+\.)\s+\S`)

// narrativeState flags the narrative at path of story, in progress or in
// review, whose ## Current state or ## Next steps is not written, by the test
// the close-out's narrative step applies (Section.Written, as flai verify's
// Unwritten reads it), or whose ## Next steps is not a list. Each finding is
// on the section's heading and advisory: a narrative is unwritten from flai
// stream open until its agent first writes it, and the close-out stops a
// story that reaches it so; failing --strict too would stop flai task done
// and the main branch's check in that window. A missing section is
// narrative.section's.
func (c *checker) narrativeState(path string, story *workitem.Item) {
	data, err := os.ReadFile(path)
	if err != nil {
		return // narrative.front-matter has read it already
	}
	text := string(data)
	if s, ok := workitem.NarrativeSection(text, workitem.CurrentState); ok && !s.Written() {
		c.advise("narrative.state", path, s.Line, "%s is %s and its narrative's ## Current state is empty or the template's placeholder; write what is true now with flai stream state %s --current, or the MCP tool stream_state", story.ID, story.Status, story.ID)
	}
	s, ok := workitem.NarrativeSection(text, workitem.NextSteps)
	switch {
	case !ok:
	case !s.Written():
		c.advise("narrative.state", path, s.Line, "%s is %s and its narrative's ## Next steps is empty or the template's placeholder; write them as a list, the very next action first, with flai stream state %s --next, or the MCP tool stream_state", story.ID, story.Status, story.ID)
	case !slices.ContainsFunc(strings.Split(s.Text, "\n"), listItem.MatchString):
		c.advise("narrative.state", path, s.Line, "%s's narrative's ## Next steps is not a list; write them as one, the very next action first, with flai stream state %s --next, or the MCP tool stream_state", story.ID, story.ID)
	}
}

func (c *checker) board() {
	b, err := c.repo.LoadBoard()
	if err != nil {
		c.add(Error, "board.front-matter", filepath.Join(c.repo.KanbanDir(), workitem.BoardFile), 1, "%v", err)
		return
	}
	counts := map[string]int{}
	for _, it := range c.items {
		if it.Type == workitem.Story && !it.Archived {
			counts[it.Status]++
		}
	}
	for _, st := range []string{workitem.Ready, workitem.InProgress, workitem.Review} {
		limit := b.WIPLimits[st]
		if limit <= 0 || counts[st] <= limit {
			continue
		}
		line := keyLine(b.Path, "wip_limits")
		if st == workitem.Review {
			// Review over its limit waits on the operator's acceptance, so it
			// warns without failing --strict; ready and in-progress do (S-0243).
			c.advise("board.wip-limit", b.Path, line, "%d stories in %s, limit %d", counts[st], st, limit)
		} else {
			c.add(Warning, "board.wip-limit", b.Path, line, "%d stories in %s, limit %d", counts[st], st, limit)
		}
	}
	for _, id := range b.Order {
		it, ok := c.byID[id]
		if !ok || it.Archived {
			c.add(Warning, "board.order", b.Path, keyLine(b.Path, "order"), "order lists %s which is not on the board", id)
		} else if it.Status != workitem.Ready && it.Status != workitem.Backlog {
			c.add(Warning, "board.order", b.Path, keyLine(b.Path, "order"), "order lists %s which is %s; only ready and backlog stories belong in the pull order", id, it.Status)
		}
	}
}

type docFront struct {
	Title   string `yaml:"title"`
	Updated string `yaml:"updated"`
	ID      string `yaml:"id"`
	Status  string `yaml:"status"`
	Date    string `yaml:"date"`
}

// ADRStatuses are the statuses an ADR may carry (ADR-0223).
var ADRStatuses = []string{"proposed", "accepted", "superseded", "deprecated", "rejected"}

var adrName = regexp.MustCompile(`^(\d{4})-.*\.md$`)

// adrIndex keeps design/adrs/README.md what its ADR files render
// (ADR-0228): an index that differs is one warning, on its first line that
// differs, naming the rows that differ, are missing, or link no file. A
// project without the folder or the index has nothing to keep, and an ADR
// whose front matter does not read is doc.front-matter's to report.
func (c *checker) adrIndex() {
	dir := filepath.Join(c.repo.Manifest.Dir(c.repo.Root, "design"), "adrs")
	index := filepath.Join(dir, "README.md")
	data, err := os.ReadFile(index)
	if err != nil {
		return
	}
	rendered, err := adrindex.Render(os.DirFS(dir), string(data))
	if err != nil {
		return // the unreadable ADR is doc.front-matter's to report
	}
	d := adrindex.Diff(string(data), rendered)
	if d.Line == 0 {
		return
	}
	var what []string
	for _, p := range []struct {
		nums        []string
		one, plural string
	}{
		{d.Changed, "row %s differs", "rows %s differ"},
		{d.Missing, "row %s is missing", "rows %s are missing"},
		{d.Extra, "row %s links no ADR file", "rows %s link no ADR file"},
	} {
		switch len(p.nums) {
		case 0:
		case 1:
			what = append(what, fmt.Sprintf(p.one, p.nums[0]))
		default:
			what = append(what, fmt.Sprintf(p.plural, andList(p.nums)))
		}
	}
	if len(what) == 0 {
		what = []string{"the order or layout of its table differs"}
	}
	c.add(Warning, "adr.index", index, d.Line, "the ADR index differs from what its ADR files render: %s; run flai check --fix or flai adr index, which rewrite it", andList(what))
}

func (c *checker) documentation() {
	for _, key := range []string{"design", "docs"} {
		root := c.repo.Manifest.Dir(c.repo.Root, key)
		_ = filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
			if err != nil {
				return nil //nolint:nilerr // an unreadable entry is skipped, the walk continues
			}
			if d.IsDir() {
				if (d.Name() == conventions.Folder || d.Name() == issues.Folder || (key == "design" && (d.Name() == experiment.Folder || d.Name() == analysis.Folder || d.Name() == digest.Folder))) && filepath.Dir(path) == root {
					return filepath.SkipDir // validated by their own rules
				}
				return nil
			}
			if !strings.HasSuffix(path, ".md") || d.Name() == "README.md" {
				return nil
			}
			data, err := os.ReadFile(path)
			if err != nil {
				return nil //nolint:nilerr // unreadable files are skipped, not fatal
			}
			fm, _, err := workitem.SplitFrontMatter(string(data))
			if err != nil {
				c.add(Warning, "doc.front-matter", path, 1, "no front matter; every design and docs file needs title and updated")
				return nil //nolint:nilerr // recorded as a finding
			}
			var f docFront
			if err := yaml.Unmarshal([]byte(fm), &f); err != nil {
				c.add(Warning, "doc.front-matter", path, 1, "front matter does not parse: %v", err)
				return nil //nolint:nilerr // recorded as a finding
			}
			isADR := filepath.Base(filepath.Dir(path)) == "adrs" && adrName.MatchString(d.Name())
			if key == "design" {
				c.docTopics(path, string(data), isADR)
			}
			if isADR {
				num := adrName.FindStringSubmatch(d.Name())[1]
				if f.ID != "ADR-"+num {
					c.add(Warning, "adr.id", path, keyLine(path, "id"), "id should be ADR-%s to match the file name", num)
				}
				if f.Status == "" || f.Date == "" || f.Title == "" {
					c.add(Warning, "adr.front-matter", path, 1, "ADRs need id, title, status, and date")
				}
				// The template, 0000, carries the status template.
				if f.Status != "" && num != "0000" && !slices.Contains(ADRStatuses, f.Status) {
					c.add(Warning, "adr.status", path, keyLine(path, "status"), "status %q is not one of %s (ADR-0223)", f.Status, strings.Join(ADRStatuses, ", "))
				}
				if _, ok := ctxpack.DecisionSentence(string(data)); !ok {
					c.add(Warning, "adr.decision", path, headingLine(path, "## Decision"), "## Decision does not open with a sentence stating the decision; a context pack briefs the ADR by that sentence (ADR-0049)")
				}
				return nil
			}
			if f.Title == "" {
				c.add(Warning, "doc.title", path, 1, "front matter needs a title")
			}
			if f.Updated == "" {
				c.add(Warning, "doc.updated", path, 1, "front matter needs an updated date")
			}
			return nil
		})
	}
}

// Code is the topic of every sub-project that is not the template (ADR-0047).
const Code = topics.Code

// Vocabulary is every topic a story can have: all, code, each sub-project's
// name, tags, and kind (ADR-0047), and the topics that the stories and epics
// among items declare (S-0135); and the words a pack primed for the
// orchestrator's reasons reads on a heading beside all, the kinds of reason
// and the built-in strategic roles' topics (ADR-0217).
func Vocabulary(repo *workitem.Repo, items []*workitem.Item) map[string]bool {
	v := map[string]bool{topics.All: true, Code: true}
	for _, w := range ctxpack.SectionWords() {
		v[w] = true
	}
	for _, it := range items {
		for _, t := range it.Topics {
			v[t] = true
		}
	}
	for _, p := range repo.Manifest.Projects {
		for _, w := range append([]string{p.Name, p.Kind}, p.Tags...) {
			if w != "" {
				v[w] = true
			}
		}
	}
	return v
}

// docTopics checks the topics of a design file: design/system and
// design/tech files declare them, and every topic on the file or a heading
// is one a story can have. ADRs may leave them out.
func (c *checker) docTopics(path, content string, isADR bool) {
	rel, _ := filepath.Rel(c.repo.Manifest.Dir(c.repo.Root, "design"), path)
	folder := strings.SplitN(filepath.ToSlash(rel), "/", 2)[0]
	if !isADR && folder != "system" && folder != "tech" {
		return
	}
	doc, err := topics.Parse(content, nil)
	if err != nil {
		c.add(Warning, "doc.topics", path, keyLine(path, "topics"), "%v", err)
		return
	}
	if !doc.Declared && !isADR {
		c.add(Warning, "doc.topics", path, 1, "no topics; say which stories it is for, such as topics: [all] or a sub-project's name, tag, or kind (ADR-0047)")
	}
	c.topicWords(path, doc)
}

// topicWords warns on each topic of a document that no story can have.
func (c *checker) topicWords(path string, doc *topics.Doc) {
	if c.vocab == nil {
		c.vocab = Vocabulary(c.repo, c.items)
	}
	unknown := func(line int, where string, list []string) {
		for _, t := range list {
			if !c.vocab[t] {
				c.add(Warning, "doc.topic", path, line, "topic %q %s is used by nothing: not all, code, a sub-project's name, tag, or kind in system-flow.yaml, nor a story's or epic's topics (%s)", t, where, strings.Join(c.words(), ", "))
			}
		}
	}
	if doc.Declared {
		unknown(keyLine(path, "topics"), "in the front matter", doc.Topics)
	}
	for _, s := range doc.Sections {
		unknown(s.Line, "on heading \""+s.Heading+"\"", s.Own)
	}
}

func (c *checker) words() []string {
	var out []string
	for w := range c.vocab {
		out = append(out, w)
	}
	sort.Strings(out)
	return out
}

var (
	criteriaHeading = regexp.MustCompile(`(?m)^## Acceptance criteria\s*$`)
	checkbox        = regexp.MustCompile(`(?m)^\s*- \[[ xX]\] \S`)
	unchecked       = regexp.MustCompile(`(?m)^\s*- \[ \] \S`)
)

func criteriaSection(body string) string {
	loc := criteriaHeading.FindStringIndex(body)
	if loc == nil {
		return ""
	}
	rest := body[loc[1]:]
	if i := strings.Index(rest, "\n## "); i >= 0 {
		rest = rest[:i]
	}
	return rest
}

func hasCriteria(body string) bool  { return checkbox.MatchString(criteriaSection(body)) }
func hasUnchecked(body string) bool { return unchecked.MatchString(criteriaSection(body)) }

func (c *checker) conventions() {
	set, errs, err := conventions.Load(c.repo)
	if err != nil {
		c.add(Error, "conventions.read", conventions.Dir(c.repo), 1, "%v", err)
		return
	}
	for _, f := range set.Validate(errs, c.definedRoles()...) {
		path := filepath.Join(c.repo.Root, f.Path)
		line := 1
		switch f.Rule {
		case "conventions.front-matter", "conventions.order":
			line = keyLine(path, "order")
			if strings.Contains(f.Message, "audience") {
				line = keyLine(path, "audience")
			}
		}
		c.add(f.Level, f.Rule, path, line, "%s", f.Message)
	}
	for _, f := range set.Files {
		if _, bad := errs[f.Path]; bad {
			continue
		}
		path := filepath.Join(c.repo.Root, f.Path)
		doc, err := topics.Parse(f.Raw, []string{topics.All})
		if err != nil {
			c.add(Warning, "doc.topics", path, keyLine(path, "topics"), "%v", err)
			continue
		}
		c.topicWords(path, doc)
	}
}

// definedRoles are the roles the project's role files define (ADR-0207),
// which a convention's roles may list. A file that does not parse defines
// none; c.roles reports it.
func (c *checker) definedRoles() []string {
	list, _, err := roles.Load(c.repo.Root, c.repo.Manifest.Layout["design"])
	if err != nil {
		return nil
	}
	var out []string
	for _, r := range list {
		if r.Role != "" {
			out = append(out, r.Role)
		}
	}
	return out
}

// RolesVocabulary is what flai check holds role files' names to (ADR-0207):
// the reads every agent has, from the guard. Tools and Commands are nil, so
// the check of an entry's MCP tool or command name is skipped: the MCP
// server and the command tree import this package, and a check that read
// them in one caller and not another would differ where it runs. A test in
// cmd holds the shipped role files' names to both.
func RolesVocabulary() roles.Vocabulary {
	return roles.Vocabulary{Reads: guard.MCPReads}
}

// roles validates the design folder's role files (ADR-0207): each problem
// roles.Load or roles.Validate finds is an error under its code, on its line,
// and each built-in role file missing is an advisory roles.baseline warning.
func (c *checker) roles() {
	design := c.repo.Manifest.Layout["design"]
	list, problems, err := roles.Load(c.repo.Root, design)
	if err != nil {
		c.add(Error, roles.CodeFile, filepath.Join(c.repo.Root, filepath.FromSlash(roles.Dir(design))), 1, "%v", err)
		return
	}
	c.missingBaselineRoles(design)
	c.agentCopies(design)
	for _, p := range append(problems, roles.Validate(list, design, RolesVocabulary())...) {
		c.add(Error, p.Code, filepath.Join(c.repo.Root, filepath.FromSlash(p.Path)), p.Line, "%s", p.Message)
	}
	if len(problems) > 0 {
		return // a file that does not parse defines no role, so the manifest's roles cannot be held to them
	}
	mf := filepath.Join(c.repo.Root, manifest.File)
	for _, msg := range c.repo.Manifest.UndeclaredRoleSettings(declaredPermissions(list)) {
		c.add(Warning, "manifest.roles", mf, keyLine(mf, "roles"), "%s", msg)
	}
}

// formulas validates the design folder's formula files (ADR-0238 §6): each
// problem formulas.Load or formulas.Validate finds is a formulas.file error on
// its line. A role step is held to the names of the roles the role files
// define, and a script step to the scripts on disk. A project with no formula
// files' folder has none to report.
func (c *checker) formulas() {
	design := c.repo.Manifest.Layout["design"]
	list, problems, err := formulas.Load(c.repo.Root, design)
	if err != nil {
		c.add(Error, formulas.CodeFile, filepath.Join(c.repo.Root, filepath.FromSlash(formulas.Dir(design))), 1, "%v", err)
		return
	}
	names := formulas.Names{Roles: c.roleNames(), Script: formulas.ScriptOnDisk(c.repo.Root)}
	for _, p := range append(problems, formulas.Validate(list, names)...) {
		c.add(Error, p.Code, filepath.Join(c.repo.Root, filepath.FromSlash(p.Path)), p.Line, "%s", p.Message)
	}
	c.analysisFormula(list, problems)
}

// analysisFormula holds analysis.formula to the formula files (S-0424): the
// formula it names is one a file defines, and analysis.schedule makes its
// instance with no inputs, so its every input carries a default. Either
// falls short as a manifest.analysis error. A formula whose file does not
// parse is left to the formulas.file error on it.
func (c *checker) analysisFormula(list []formulas.Formula, problems []formulas.Problem) {
	name := c.repo.Manifest.Analysis.Formula
	if name == "" {
		return
	}
	file := path.Join(formulas.Dir(c.repo.Manifest.Layout["design"]), name+".md")
	for _, p := range problems {
		if p.Path == file {
			return
		}
	}
	mf := filepath.Join(c.repo.Root, manifest.File)
	f, err := formulas.Find(list, name)
	if err != nil {
		c.add(Error, "manifest.analysis", mf, keyLine(mf, "analysis"), "analysis.formula: %v", err)
		return
	}
	for _, in := range f.Inputs {
		if in.Required() {
			c.add(Error, "manifest.analysis", mf, keyLine(mf, "analysis"), "analysis.formula names %s, whose input %s has no default, but analysis.schedule makes its instance with no inputs; give %s a default in %s, or name another formula", name, in.Name, in.Name, f.Path)
		}
	}
}

// roleNames are the names of the roles the project's role files define
// (ADR-0207), which a formula's role step names (ADR-0238 §3), as definedRoles
// are their verbs. A file that does not parse defines none; c.roles reports it.
func (c *checker) roleNames() []string {
	list, _, err := roles.Load(c.repo.Root, c.repo.Manifest.Layout["design"])
	if err != nil {
		return nil
	}
	var out []string
	for _, r := range list {
		if r.Name != "" {
			out = append(out, r.Name)
		}
	}
	return out
}

// missingBaselineRoles advises roles.baseline for each role flai carries a
// baseline for whose file the project lacks: flai runs the baseline in its
// place, and flai upgrade adds the file (S-0419).
func (c *checker) missingBaselineRoles(design string) {
	for _, r := range roles.Baseline(design) {
		path := filepath.Join(c.repo.Root, filepath.FromSlash(r.Path))
		if _, err := os.Stat(path); !errors.Is(err, os.ErrNotExist) {
			continue
		}
		c.advise(roles.CodeBaseline, path, 1, "%s is missing: flai runs the %s from the baseline it carries; run flai upgrade to add the file from the template", r.Path, r.Name)
	}
}

// agentCopies warns roles.definition for each built-in role whose
// .claude/agents/<name>.md exists and is neither the link ADR-0236 names,
// ../../.flai-cache/agents/<name>.md, whether or not its target exists, nor
// a file whose description, tools, model, and prompt equal the definition
// flai generates from the role's file, naming the first that differs. A
// missing file is no finding, and neither is a role whose own definition
// cannot be generated: the role file's findings cover it. A warning, since
// only a session started by hand reads the file, so --strict fails on it.
func (c *checker) agentCopies(design string) {
	for _, base := range roles.Baseline(design) {
		name := base.Name
		rel := filepath.Join(".claude", "agents", name+".md")
		path := filepath.Join(c.repo.Root, rel)
		link := "../../.flai-cache/agents/" + name + ".md"
		fix := fmt.Sprintf("run flai role definition --write, then ln -sf %s %s (ADR-0236)", link, filepath.ToSlash(rel))
		info, err := os.Lstat(path)
		switch {
		case errors.Is(err, os.ErrNotExist):
			continue
		case err != nil:
			c.add(Warning, roles.CodeDefinition, path, 1, "%s cannot be read: %v; it should be the link %s: %s", filepath.ToSlash(rel), err, link, fix)
			continue
		case info.Mode()&os.ModeSymlink != 0:
			target, err := os.Readlink(path)
			if err == nil && target == link {
				continue
			}
			if err != nil {
				target = fmt.Sprintf("a target that cannot be read (%v)", err)
			}
			c.add(Warning, roles.CodeDefinition, path, 1, "%s is a symbolic link to %s; it should be the link %s, to the definition flai generates for the %s: %s", filepath.ToSlash(rel), target, link, name, fix)
			continue
		case !info.Mode().IsRegular():
			c.add(Warning, roles.CodeDefinition, path, 1, "%s is neither a file nor a symbolic link; it should be the link %s: %s", filepath.ToSlash(rel), link, fix)
			continue
		}
		r, fromBaseline, err := roles.Builtin(c.repo.Root, design, name)
		if err != nil {
			continue
		}
		if fromBaseline {
			if _, err := os.Stat(filepath.Join(c.repo.Root, filepath.FromSlash(r.Path))); err == nil {
				continue // the role file does not parse, which its roles.file error reports
			}
		}
		want, err := roles.Generate(r)
		if err != nil {
			continue // the role file's own roles.definition errors report it
		}
		source := "its role file " + r.Path
		if fromBaseline {
			source = "the baseline flai carries for " + r.Path
		}
		data, err := os.ReadFile(path)
		if err != nil {
			c.add(Warning, roles.CodeDefinition, path, 1, "%s cannot be read: %v; replace it with the link to the definition flai generates for the %s: %s", filepath.ToSlash(rel), err, name, fix)
			continue
		}
		got, err := roles.ParseAgentFile(data)
		if err != nil {
			c.add(Warning, roles.CodeDefinition, path, 1, "%s is not a Claude Code agent file: %v; replace it with the link to the definition flai generates from %s: %s", filepath.ToSlash(rel), err, source, fix)
			continue
		}
		field, line := agentDiff(got, want, path)
		if field == "" {
			continue
		}
		c.add(Warning, roles.CodeDefinition, path, line, "%s is a copy whose %s differs from the definition flai generates for the %s from %s; replace it with the link to that definition: %s", filepath.ToSlash(rel), field, name, source, fix)
	}
}

// agentDiff is the first of description, tools, model, and prompt in which
// the agent file at path, read as got, differs from the definition want, and
// its line there: the key's in the front matter, or 1 for the prompt. An
// empty field means none differs.
func agentDiff(got, want roles.Generated, path string) (string, int) {
	model := func(m string) string {
		if m == "" {
			return roles.Inherit
		}
		return m
	}
	switch {
	case strings.TrimSpace(got.Description) != strings.TrimSpace(want.Description):
		return "description", keyLine(path, "description")
	case !slices.Equal(got.Tools, want.Tools):
		return "tools", keyLine(path, "tools")
	case model(got.Model) != model(want.Model):
		return "model", keyLine(path, "model")
	case strings.TrimSpace(got.Prompt) != strings.TrimSpace(want.Prompt):
		return "prompt", 1
	}
	return "", 0
}

// declaredPermissions are, by the name of each role the files define, the
// permissions its entries name as roles.<name>.permissions.<p>.
func declaredPermissions(list []roles.Role) map[string][]string {
	out := map[string][]string{}
	for _, r := range list {
		perms := out[r.Name]
		for _, p := range r.Permissions {
			if perm, owner, ok := roles.ParseSetting(p.Setting); ok && owner == r.Name && strings.HasPrefix(p.Setting, "roles.") {
				perms = append(perms, perm)
			}
		}
		out[r.Name] = perms
	}
	return out
}

// unknownFields reports each front-matter field this flai does not know: the
// listing paths read past them (S-0181), and check stays strict.
func (c *checker) unknownFields(rule, path string, fields []workitem.Field) {
	for i, msg := range workitem.UnknownFieldErrors(fields) {
		c.add(Error, rule, path, keyLine(path, fields[i].Name), "%s", msg)
	}
}

// threads validates wip/threads (ADR-0020): schema, file names, anchors
// that exist, headings that are present, and open threads on archived items.
func (c *checker) threads() {
	dir := threads.Dir(c.repo)
	if st, err := os.Stat(dir); err != nil || !st.IsDir() {
		return // optional until the first thread
	}
	list, err := threads.List(c.repo)
	if err != nil {
		c.add(Error, "threads.front-matter", dir, 1, "%v", err)
		return
	}
	byID := map[string]*workitem.Item{}
	for _, it := range c.items {
		byID[it.ID] = it
	}
	seen := map[string]string{}
	for _, th := range list {
		if err := th.Validate(); err != nil {
			c.add(Error, "threads.front-matter", th.Path, keyLine(th.Path, "id"), "%v", err)
			continue
		}
		if prev, dup := seen[th.ID]; dup {
			c.add(Error, "threads.duplicate-id", th.Path, keyLine(th.Path, "id"), "%s is also defined in %s", th.ID, prev)
		}
		seen[th.ID] = th.Path
		c.unknownFields("threads.unknown-field", th.Path, th.Unknown)
		if base := filepath.Base(th.Path); !strings.HasPrefix(base, th.ID+"-") {
			c.add(Error, "threads.filename", th.Path, 1, "file name should start with %s-", th.ID)
		}
		// An item anchor follows the item: archiving moves the file, the
		// thread stays attached. A plain path anchor must exist as written.
		anchorPath := th.Anchor.Path
		if th.Anchor.Item != "" {
			it, ok := byID[th.Anchor.Item]
			switch {
			case !ok:
				c.add(Error, "threads.anchor", th.Path, keyLine(th.Path, "anchor"), "item %s does not exist", th.Anchor.Item)
				continue
			case it.Archived && th.Open():
				c.add(Warning, "threads.archived", th.Path, keyLine(th.Path, "status"), "%s is %s but %s is archived; resolve it with flai thread resolve %s", th.ID, th.Status, th.Anchor.Item, th.ID)
			}
			anchorPath = it.Path
		}
		abs := anchorPath
		if !filepath.IsAbs(abs) {
			abs = filepath.Join(c.repo.Root, filepath.FromSlash(anchorPath))
		}
		data, err := os.ReadFile(abs)
		if err != nil && !filepath.IsAbs(anchorPath) {
			if d2, err2 := os.ReadFile(filepath.Join(c.repo.MainRoot, filepath.FromSlash(anchorPath))); err2 == nil {
				data, err = d2, nil
			}
		}
		if err != nil {
			c.add(Error, "threads.anchor", th.Path, keyLine(th.Path, "anchor"), "anchor %s does not exist", th.Anchor.Path)
		} else if th.Anchor.Heading != "" && !threads.HasHeading(string(data), th.Anchor.Heading) {
			c.add(Warning, "threads.heading", th.Path, keyLine(th.Path, "anchor"), "heading %q is no longer in %s", th.Anchor.Heading, th.Anchor.Path)
		}
		if len(th.Entries()) == 0 {
			c.add(Warning, "threads.entries", th.Path, 1, "%s has no dated entries", th.ID)
		}
	}
}

// messages validates wip/messages (ADR-0120): each conversation's front
// matter, unique IDs, file names, that both its stories exist, its dated
// entries, and an open conversation of an archived story. Its markdown is
// the wip walk's (markdown). Every markdown file but README.md is read on its
// own, so that one that does not parse is reported and the rest still are.
func (c *checker) messages() {
	dir := messages.Dir(c.repo)
	entries, err := os.ReadDir(dir)
	if err != nil {
		return // optional until the first message
	}
	seen := map[string]string{}
	for _, e := range entries {
		name := e.Name()
		if e.IsDir() || !strings.HasSuffix(name, ".md") || name == "README.md" {
			continue
		}
		path := filepath.Join(dir, name)
		data, err := os.ReadFile(path)
		if err != nil {
			c.add(Error, "messages.front-matter", path, 1, "%v", err)
			continue
		}
		conv, err := messages.Parse(string(data))
		if err == nil {
			err = conv.Validate()
		}
		if err != nil {
			c.add(Error, "messages.front-matter", path, keyLine(path, "id"), "%v", err)
			continue
		}
		if prev, dup := seen[conv.ID]; dup {
			c.add(Error, "messages.duplicate-id", path, keyLine(path, "id"), "%s is also defined in %s", conv.ID, prev)
		}
		seen[conv.ID] = path
		if rel, err := filepath.Rel(c.repo.Root, path); err == nil {
			seen[conv.ID] = rel
		}
		c.unknownFields("messages.unknown-field", path, conv.Unknown)
		if !strings.HasPrefix(name, conv.ID+"-") {
			c.add(Error, "messages.filename", path, 1, "file name should start with %s-", conv.ID)
		}
		archived := ""
		for _, side := range []struct{ field, story string }{{"from", conv.From}, {"to", conv.To}} {
			it := c.named(side.story)
			switch {
			case it == nil:
				c.add(Error, "messages.story", path, keyLine(path, side.field), "%s names story %s, which does not exist", side.field, side.story)
			case it.Archived && archived == "":
				archived = it.ID
			}
		}
		// flai accept and flai archive close it; one an older flai left open
		// is closed by hand, as they would have closed it
		if archived != "" && conv.Status == messages.StatusOpen {
			c.add(Warning, "messages.closed", path, keyLine(path, "status"), "%s is open but %s is archived, so it reads as closed; close it as flai archive would: set its status to closed and add a dated entry \"Closed: %s was archived\" under ## Entries", conv.ID, archived, archived)
		}
		if len(conv.Entries()) == 0 {
			c.add(Warning, "messages.entries", path, 1, "%s has no dated entries", conv.ID)
		}
	}
}

func (c *checker) issues() {
	dir := issues.Dir(c.repo)
	if st, err := os.Stat(dir); err != nil || !st.IsDir() {
		return // optional until the first issue is recorded
	}
	list, err := issues.List(c.repo)
	if err != nil {
		c.add(Error, "issues.front-matter", dir, 1, "%v", err)
		return
	}
	seen := map[string]string{}
	for _, is := range list {
		if err := is.Validate(); err != nil {
			c.add(Error, "issues.front-matter", is.Path, keyLine(is.Path, "id"), "%v", err)
		}
		if prev, dup := seen[is.ID]; dup {
			c.add(Error, "issues.duplicate-id", is.Path, keyLine(is.Path, "id"), "%s is also defined in %s", is.ID, prev)
		}
		seen[is.ID] = is.Path
		c.unknownFields("issues.unknown-field", is.Path, is.Unknown)
		if base := filepath.Base(is.Path); !strings.HasPrefix(base, is.ID+"-") {
			c.add(Error, "issues.filename", is.Path, 1, "file name should start with %s-", is.ID)
		}
	}
	c.issuesNoStory(list)
	if len(list) == 0 {
		return // a project that has never had an issue keeps the template's summary
	}
	summaryPath := filepath.Join(dir, issues.SummaryFile)
	data, err := os.ReadFile(summaryPath)
	if err != nil {
		c.add(Warning, "issues.summary", dir, 1, "summary.md is missing; run flai check --fix or flai issue summary, which write it")
		return
	}
	if summaryUpdated.ReplaceAllString(string(data), "") != summaryUpdated.ReplaceAllString(issues.Summary(list, c.now), "") {
		c.add(Warning, "issues.summary", summaryPath, 1, "summary.md differs from what its issue files render; run flai check --fix or flai issue summary, which rewrite it")
	}
}

// summaryUpdated is summary.md's updated line, which says when it was
// written rather than what it holds, so the issues.summary comparison leaves
// it out.
var summaryUpdated = regexp.MustCompile(`(?m)^updated: .*$`)

// issuesNoStory warns about each open issue older than issues.story_after
// that no open story links (S-0198). An invalid story_after is layout's
// finding, and turns this rule off, as 0 does.
func (c *checker) issuesNoStory(list []*issues.Issue) {
	after, err := c.repo.Manifest.Issues.StoryAfterDuration()
	if err != nil || after == 0 {
		return
	}
	setting := "issues.story_after is " + strings.TrimSpace(c.repo.Manifest.Issues.StoryAfter)
	if strings.TrimSpace(c.repo.Manifest.Issues.StoryAfter) == "" {
		setting = "issues.story_after is unset, so " + strings.TrimSuffix(manifest.DefaultStoryAfter.String(), "0m0s")
	}
	for _, is := range issues.NoStory(list, c.items) {
		first, err := time.Parse(workitem.TimeFormat, is.FirstReported)
		if err != nil || c.now.Sub(first) <= after {
			continue // an unreadable first_reported is issues.front-matter's finding
		}
		c.add(Warning, "issues.no-story", is.Path, keyLine(is.Path, "first_reported"),
			"%s has been open since %s with no open story linking it (%s): make one with flai issue story %s",
			is.ID, first.Format(time.DateOnly), setting, is.ID)
	}
}

// experiments validates design/experiments (ADR-0066): each results document's
// front matter, its file name, its sections, and that it names an experiment
// story that exists.
func (c *checker) experiments() {
	dir := filepath.Join(c.repo.Root, filepath.FromSlash(experiment.Dir(c.repo.Manifest)))
	entries, err := os.ReadDir(dir)
	if err != nil {
		return // optional until the first experiment records its results
	}
	seen := map[string]string{}
	for _, e := range entries {
		name := e.Name()
		if e.IsDir() || !strings.HasSuffix(name, ".md") || experiment.Skipped(name) {
			continue
		}
		path := filepath.Join(dir, name)
		data, err := os.ReadFile(path)
		if err != nil {
			c.add(Error, "experiments.front-matter", path, 1, "%v", err)
			continue
		}
		f, problems := experiment.Validate(name, string(data))
		for _, p := range problems {
			line := 1
			switch {
			case p.Field != "":
				line = keyLine(path, p.Field)
			case p.Heading != "":
				line = headingLine(path, p.Heading)
			}
			c.add(Error, "experiments.document", path, line, "%s", p.Message)
		}
		if f.Story == "" {
			continue
		}
		if prev, dup := seen[f.Story]; dup {
			c.add(Error, "experiments.duplicate", path, keyLine(path, "story"), "%s already has its results in %s", f.Story, prev)
		}
		seen[f.Story] = path
		switch it := c.named(f.Story); {
		case it == nil:
			c.add(Error, "experiments.story", path, keyLine(path, "story"), "story %s does not exist", f.Story)
		case it.Type != workitem.Story || it.Nature != experiment.Nature:
			c.add(Warning, "experiments.story", path, keyLine(path, "story"), "%s is a %s %s, not an experiment story; results of anything else belong elsewhere", f.Story, it.Nature, it.Type)
		}
	}
}

// analysis validates design/analysis (S-0223): each report's file name and
// front matter, and that the folder's README.md lists it, as adrIndex keeps
// the ADRs' index, with no row linking to a report that is not there.
func (c *checker) analysis() {
	dir := filepath.Join(c.repo.Root, filepath.FromSlash(analysis.Dir(c.repo.Manifest)))
	entries, err := os.ReadDir(dir)
	if err != nil {
		return // optional until the analyzer writes its first report
	}
	index := filepath.Join(dir, analysis.Index)
	listed, _ := os.ReadFile(index) // without an index, every report is unlisted, and said to be
	present := map[string]bool{}
	for _, e := range entries {
		name := e.Name()
		if e.IsDir() || !strings.HasSuffix(name, ".md") || analysis.Skipped(name) {
			continue
		}
		present[name] = true
		path := filepath.Join(dir, name)
		data, err := os.ReadFile(path)
		if err != nil {
			c.add(Error, "analysis.front-matter", path, 1, "%v", err)
			continue
		}
		_, problems := analysis.Validate(name, string(data))
		for _, p := range problems {
			line := 1
			if p.Field != "" {
				line = keyLine(path, p.Field)
			}
			c.add(Error, "analysis.report", path, line, "%s", p.Message)
		}
		if !analysis.Lists(string(listed), name) {
			c.add(Warning, "analysis.index", path, 1, "%s/%s does not list it; add a row linking to %s", analysis.Dir(c.repo.Manifest), analysis.Index, name)
		}
	}
	for i, l := range strings.Split(string(listed), "\n") {
		for _, m := range analysisLink.FindAllStringSubmatch(l, -1) {
			if !present[m[1]] {
				c.add(Warning, "analysis.index", index, i+1, "the row links to %s, which is not there", m[1])
			}
		}
	}
}

var analysisLink = regexp.MustCompile(`\]\((?:\./)?(\d{4}-\d{2}-\d{2}-[a-z]+\.md)\)`)

// digest validates the design folder's digest (ADR-0191): each entry's name
// and front matter, what its sources and absorbed_by name, as digest.entry,
// and text with the shape of a secret, as digest.secret, each on its line.
func (c *checker) digest() {
	dir := filepath.Join(c.repo.Root, filepath.FromSlash(digest.Dir(c.repo.Manifest)))
	entries, err := os.ReadDir(dir)
	if err != nil {
		return // optional until the first entry is written
	}
	env := digest.EnvFor(c.repo)
	for _, e := range entries {
		name := e.Name()
		if e.IsDir() || !strings.HasSuffix(name, ".md") || digest.Skipped(name) {
			continue
		}
		path := filepath.Join(dir, name)
		data, err := os.ReadFile(path)
		if err != nil {
			c.add(Error, digest.RuleEntry, path, 1, "%v", err)
			continue
		}
		_, problems := digest.Validate(name, string(data), env)
		for _, p := range problems {
			c.add(Error, p.Rule, path, p.Line, "%s", p.Message)
		}
	}
}

// LintFindings are a refused write's markdown lint findings as flai check
// reports them, for the commands that refuse with check findings.
func LintFindings(e *mdlint.Error) []Finding {
	out := make([]Finding, len(e.Findings))
	for i, f := range e.Findings {
		out[i] = Finding{Level: Warning, Rule: "markdown." + f.Rule, Path: e.Path, Line: f.Line, Message: f.String()}
	}
	return out
}

// markdown lints every markdown file under the wip folder with the project's
// markdownlint configuration (S-0179). flai writes them in the main
// checkout, where no story's lint runs (ADR-0019), so they reach main unlinted
// unless flai check says so first. Findings are warnings: --strict fails on
// them, and the commands that write refuse a change that brings one. The
// configuration is the main checkout's, or the worktree's when the main
// checkout's root holds none, as in an agent container that mounts only its
// .git, wip, and system-flow.yaml (ADR-0253); in the main checkout the two
// are the same. A project without a markdownlint configuration is not
// linted.
func (c *checker) markdown() {
	root := c.repo.MainRoot
	if root == "" {
		root = c.repo.Root
	}
	cfg, err := mdlint.Load(root)
	if err == nil && cfg == nil && root != c.repo.Root {
		root = c.repo.Root
		cfg, err = mdlint.Load(root)
	}
	if err != nil {
		c.add(Warning, "markdown.config", root, 1, "the markdownlint configuration cannot be read, so wip is not linted: %v", err)
		return
	}
	if cfg == nil {
		return
	}
	for _, ff := range cfg.LintFolder(c.repo.WipDir()) {
		for _, f := range ff.Findings {
			c.add(Warning, "markdown."+f.Rule, ff.Path, f.Line, "%s", f)
		}
	}
}

// nodeRunner runs node for the components.examples check; tests replace it.
var nodeRunner components.Runner = components.System{}

// componentExamples reports components.examples (ADR-0164): each examples
// file under the design folder's components folder that does not fit the
// format, names a component that does not exist, or has examples that do not
// fit the component's props, read with node once for each sub-project with
// examples. Warnings, at the problem's line, or line 1 for a file without its
// component or whose component's props could not be read. A project without
// examples files is not checked. One with no sub-project whose components
// flai reads, or a host without node, gets one note, which --strict passes
// over, and so does a sub-project whose svelte is not installed, on its
// examples folder (ADR-0171); without node or svelte the format and the
// components are still checked. The notes for node and svelte are marked
// unchecked.
func (c *checker) componentExamples() {
	folder := filepath.Join(c.repo.Root, filepath.FromSlash(components.ExamplesFolder(c.repo.Manifest)))
	var files []string
	_ = filepath.WalkDir(folder, func(p string, e fs.DirEntry, err error) error {
		if err == nil && e.Type().IsRegular() && strings.HasSuffix(e.Name(), components.ExamplesSuffix) {
			files = append(files, p)
		}
		return nil //nolint:nilerr // an absent folder has no examples, and an unreadable entry is skipped
	})
	if len(files) == 0 {
		return
	}
	shown := components.ExamplesFolder(c.repo.Manifest)
	projects := map[string]manifest.Project{}
	for _, p := range c.repo.Manifest.Projects {
		if manifest.SupportsComponents(p.Kind) {
			projects[p.Name] = p
		}
	}
	if len(projects) == 0 {
		c.advise("components.examples", folder, 1, "the %d examples files under %s are not checked: system-flow.yaml names no sub-project of kind %s, whose components they are of; add the sub-project under projects with its kind, or remove the files", len(files), shown, strings.Join(manifest.ComponentKinds, " or "))
		return
	}
	_, nodeErr := nodeRunner.LookPath("node")
	// Each sub-project's components, found once, and the files of those with
	// examples, in the order they were found.
	found := map[string]map[string]components.Component{}
	var withExamples []string
	byProject := map[string][]exampleFile{}
	for _, path := range files {
		f, ok := c.exampleFile(path, folder, projects, found)
		if !ok {
			continue
		}
		if _, ok := byProject[f.project]; !ok {
			withExamples = append(withExamples, f.project)
		}
		byProject[f.project] = append(byProject[f.project], f)
	}
	if nodeErr != nil {
		c.adviseUnchecked(folder, "the props of the examples under %s are not checked against their components': flai reads a component's props with node, which is not on PATH; install Node.js (https://nodejs.org) to check them", shown)
		return
	}
	for _, name := range withExamples {
		c.exampleProps(projects[name], byProject[name], folder, shown)
	}
}

// exampleFile is one examples file whose format was checked and whose
// component exists.
type exampleFile struct {
	path     string
	project  string
	comp     components.Component
	examples []components.Example
}

// exampleFile checks the examples file at path, under folder: its format,
// then that its component is one of a sub-project of projects, found once
// into found. False when its component does not exist or its sub-project's
// components could not be found, each reported.
func (c *checker) exampleFile(path, folder string, projects map[string]manifest.Project, found map[string]map[string]components.Component) (exampleFile, bool) {
	const rule = "components.examples"
	data, err := os.ReadFile(path)
	if err != nil {
		c.add(Warning, rule, path, 1, "cannot read the examples file: %v; check that it is a readable file", err)
		return exampleFile{}, false
	}
	examples, problems := components.ParseExamples(data)
	for _, p := range problems {
		c.add(Warning, rule, path, p.Line, "%s", p)
	}
	rel, err := filepath.Rel(folder, path)
	if err != nil {
		return exampleFile{}, false
	}
	name, compPath, ok := components.ExamplesComponent(filepath.ToSlash(rel))
	if !ok {
		c.add(Warning, rule, path, 1, "the file is in no sub-project's folder, so it is no component's; move it to the folder of the sub-project's name, at the component's path with .svelte made %s", components.ExamplesSuffix)
		return exampleFile{}, false
	}
	p, ok := projects[name]
	if !ok {
		c.add(Warning, rule, path, 1, "system-flow.yaml names no sub-project %q of kind %s, so the file's component %s does not exist; move the file to the folder of its sub-project's name, or remove it", name, strings.Join(manifest.ComponentKinds, " or "), compPath)
		return exampleFile{}, false
	}
	comps, ok := found[name]
	if !ok {
		s, _ := p.ComponentSettings()
		list, err := components.Find(c.repo.Root, p.Path, s.Include, s.Exclude)
		if err != nil {
			c.add(Warning, rule, path, 1, "its examples are not checked: %v", err)
			return exampleFile{}, false
		}
		comps = make(map[string]components.Component, len(list))
		for _, comp := range list {
			comps[comp.Path] = comp
		}
		found[name] = comps
	}
	comp, ok := comps[compPath]
	if !ok {
		c.add(Warning, rule, path, 1, "sub-project %q has no component %s, so the file is no component's; move it with its component, or remove it", name, compPath)
		return exampleFile{}, false
	}
	return exampleFile{path: path, project: name, comp: comp, examples: examples}, true
}

// exampleProps reads the props of files' components, all of project, with
// node once, and reports each example that does not fit its component's;
// a file whose component's props could not be read is reported once, with
// the reader's error. When svelte is not installed in project, the reader
// stops for all its files, and the project's folder under folder, the
// examples folder, shown as shown, gets one note instead, marked unchecked
// (ADR-0171).
func (c *checker) exampleProps(project manifest.Project, files []exampleFile, folder, shown string) {
	const rule = "components.examples"
	comps := make([]components.Component, len(files))
	for i, f := range files {
		comps[i] = f.comp
	}
	read, err := components.ReadProps(context.Background(), nodeRunner, c.repo.Root, comps)
	if errors.Is(err, components.ErrSvelteMissing) {
		c.adviseUnchecked(filepath.Join(folder, project.Name), "the props of the examples under %s/%s are not checked against their components': %v", shown, project.Name, err)
		return
	}
	for i, f := range files {
		switch {
		case err != nil:
			c.add(Warning, rule, f.path, 1, "its examples are not checked, since the props of %s in %s could not be read: %v", f.comp.Path, project.Name, err)
		case read[i].Error != "":
			c.add(Warning, rule, f.path, 1, "its examples are not checked, since the props of %s in %s could not be read: %s", f.comp.Path, project.Name, read[i].Error)
		default:
			for _, p := range components.Validate(f.examples, read[i].Props) {
				c.add(Warning, rule, f.path, p.Line, "%s", p)
			}
		}
	}
}

// conflictMarkers reports each line that opens or closes a merge conflict in
// the markdown under the design, docs, and wip folders and at the root: a
// conflict git-added unresolved reached main in a design document and nothing
// caught it (S-0253, I-0066). Errors, so CI's check fails without --strict.
func (c *checker) conflictMarkers() {
	seen := map[string]bool{}
	scan := func(path string) {
		if seen[path] {
			return
		}
		seen[path] = true
		data, err := os.ReadFile(path)
		if err != nil {
			return // unreadable files are skipped, as the other rules skip them
		}
		for _, line := range conflictmark.Lines(string(data)) {
			c.add(Error, "markdown.conflict-marker", path, line, "a merge conflict marker: resolve the conflict, keeping what both sides meant, and remove the markers")
		}
	}
	root := filepath.Clean(c.repo.Root)
	for _, dir := range []string{c.repo.Manifest.Dir(c.repo.Root, "design"), c.repo.Manifest.Dir(c.repo.Root, "docs"), c.repo.WipDir()} {
		if filepath.Clean(dir) == root {
			continue // a folder laid out at the root: its markdown is the root's, below
		}
		_ = filepath.WalkDir(dir, func(path string, e fs.DirEntry, err error) error {
			if err != nil || e.IsDir() || !strings.HasSuffix(path, ".md") {
				return nil //nolint:nilerr // an unreadable entry is skipped, the walk continues
			}
			scan(path)
			return nil
		})
	}
	entries, err := os.ReadDir(root)
	if err != nil {
		return
	}
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".md") {
			scan(filepath.Join(root, e.Name()))
		}
	}
}
