// Package check validates a conforming repository against the standard:
// work items, narratives, the board, and documentation front matter. It is
// the reference validator that flaiover mirrors.
package check

import (
	"bufio"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"sort"
	"strings"
	"time"

	"github.com/goccy/go-yaml"

	"github.com/bytepunx/system-flow/flai/internal/analysis"
	"github.com/bytepunx/system-flow/flai/internal/conflictmark"
	ctxpack "github.com/bytepunx/system-flow/flai/internal/context"
	"github.com/bytepunx/system-flow/flai/internal/conventions"
	"github.com/bytepunx/system-flow/flai/internal/experiment"
	"github.com/bytepunx/system-flow/flai/internal/issues"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/mdlint"
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
	// advisory marks a warning --strict passes over (advise).
	advisory bool
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
	c.workItems()
	c.narratives()
	c.activities()
	c.overlap()
	c.componentTag()
	c.after()
	c.unaccepted()
	c.epicLags()
	c.board()
	c.documentation()
	c.adrIndex()
	c.conventions()
	c.issues()
	c.experiments()
	c.analysis()
	c.threads()
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
	for _, msg := range m.Analysis.Errors() {
		mf := filepath.Join(c.repo.Root, "system-flow.yaml")
		c.add(Error, "manifest.analysis", mf, keyLine(mf, "analysis"), "%s", msg)
	}
	for _, msg := range m.TestErrors() {
		mf := filepath.Join(c.repo.Root, "system-flow.yaml")
		c.add(Error, "manifest.tests", mf, keyLine(mf, "tests"), "%s", msg)
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

// overlap warns when two in-progress items declare touches that cover the
// same path (ADR-0019); it is advisory, humans and agents coordinate. An
// overlap that lies wholly inside a shared path is not reported (ADR-0096).
func (c *checker) overlap() {
	holds := workitem.NewHolds(nil, c.repo.Manifest.Projects).WithShared(c.repo.SharedClaims())
	type owner struct {
		it   *workitem.Item
		path string
	}
	var owners []owner
	for _, it := range c.items {
		if it.Archived || it.Type == workitem.Epic || it.Status != workitem.InProgress {
			continue
		}
		for _, p := range it.Touches {
			owners = append(owners, owner{it, p})
		}
	}
	for i := 0; i < len(owners); i++ {
		for j := i + 1; j < len(owners); j++ {
			a, b := owners[i], owners[j]
			if a.it.ID == b.it.ID || a.it.Parent == b.it.ID || b.it.Parent == a.it.ID {
				continue
			}
			if holds.Overlaps(a.path, b.path) {
				c.add(Warning, "wip.overlap", a.it.Path, keyLine(a.it.Path, "touches"), "%s touches %s, which %s (in progress) also touches as %s", a.it.ID, a.path, b.it.ID, b.path)
			}
		}
	}
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

// after reports an after: entry that names no story, the story itself, or
// a story that waits, through after:, for this one: each would hold the
// story until someone edits it (S-0130, ADR-0046). A task's after: names
// tasks of its own story, and is reported the same way under task.after,
// with an entry that names a task of another story (S-0176). The archive is
// not edited and is not reported.
func (c *checker) after() {
	named := func(e string) *workitem.Item {
		if it, ok := c.byID[workitem.CanonicalID(e)]; ok {
			return it
		}
		return c.byID[e]
	}
	for _, it := range c.items {
		if it.Archived || (it.Type != workitem.Story && it.Type != workitem.Task) {
			continue
		}
		p, rule := it.Path, it.Type+".after"
		for _, e := range it.After {
			switch other := named(e); {
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
		if cycle := c.afterCycle(it, named); cycle != nil {
			c.add(Error, rule, p, keyLine(p, "after"), "after forms a cycle, %s, so none of them would start; drop one of the entries", strings.Join(cycle, " waits for "))
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
		n, err := workitem.ReadNarrative(path)
		if err != nil {
			c.add(Error, "narrative.front-matter", path, 1, "%v", err)
			continue
		}
		want := strings.TrimSuffix(filepath.Base(path), ".md")
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

var adrName = regexp.MustCompile(`^(\d{4})-.*\.md$`)

var adrIndexRow = regexp.MustCompile(`^\| *\[\d{4}\]\(([^)]+\.md)\)`)

// adrIndex keeps design/adrs/README.md honest (S-0060): every ADR file has a
// row, and every row has its file. flai adr new writes both; this is for the
// ADR someone made by hand. The template, 0000, has no row.
func (c *checker) adrIndex() {
	dir := filepath.Join(c.repo.Manifest.Dir(c.repo.Root, "design"), "adrs")
	index := filepath.Join(dir, "README.md")
	data, err := os.ReadFile(index)
	if err != nil {
		return // a project without an ADR index has nothing to keep honest
	}
	rows := map[string]int{}
	for i, l := range strings.Split(string(data), "\n") {
		if m := adrIndexRow.FindStringSubmatch(l); m != nil {
			rows[m[1]] = i + 1
		}
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return
	}
	present := map[string]bool{}
	for _, e := range entries {
		m := adrName.FindStringSubmatch(e.Name())
		if e.IsDir() || m == nil || m[1] == "0000" {
			continue
		}
		present[e.Name()] = true
		if _, ok := rows[e.Name()]; !ok {
			indexRel := index
			if r, err := filepath.Rel(c.repo.Root, index); err == nil {
				indexRel = filepath.ToSlash(r)
			}
			c.add(Warning, "adr.index", filepath.Join(dir, e.Name()), 1, "no row in %s; add one, or record ADRs with flai adr new, which does", indexRel)
		}
	}
	for name, line := range rows {
		if !present[name] && name != "0000-template.md" {
			c.add(Warning, "adr.index", index, line, "the row links to %s, which is not there", name)
		}
	}
}

func (c *checker) documentation() {
	for _, key := range []string{"design", "docs"} {
		root := c.repo.Manifest.Dir(c.repo.Root, key)
		_ = filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
			if err != nil {
				return nil //nolint:nilerr // an unreadable entry is skipped, the walk continues
			}
			if d.IsDir() {
				if (d.Name() == conventions.Folder || d.Name() == issues.Folder || (key == "design" && (d.Name() == experiment.Folder || d.Name() == analysis.Folder))) && filepath.Dir(path) == root {
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
// among items declare (S-0135).
func Vocabulary(repo *workitem.Repo, items []*workitem.Item) map[string]bool {
	v := map[string]bool{topics.All: true, Code: true}
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
	for _, f := range set.Validate(errs) {
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
	summaryPath := filepath.Join(dir, issues.SummaryFile)
	data, err := os.ReadFile(summaryPath)
	if err != nil {
		if len(list) > 0 {
			c.add(Warning, "issues.summary", dir, 1, "summary.md is missing; run flai issue summary")
		}
		return
	}
	text := string(data)
	for _, is := range list {
		linked := strings.Contains(text, "["+is.ID+"]")
		switch {
		case is.Status == "open" && !linked:
			c.add(Warning, "issues.summary", summaryPath, 1, "open issue %s is not in summary.md; run flai issue summary", is.ID)
		case is.Status != "open" && linked:
			c.add(Warning, "issues.summary", summaryPath, 1, "closed issue %s is still in summary.md; run flai issue summary", is.ID)
		}
	}
}

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
		switch it := c.byID[workitem.CanonicalID(f.Story)]; {
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
// them, and the commands that write refuse a change that brings one. A
// project without a markdownlint configuration is not linted.
func (c *checker) markdown() {
	root := c.repo.MainRoot
	if root == "" {
		root = c.repo.Root
	}
	cfg, err := mdlint.Load(root)
	if err != nil {
		c.add(Warning, "markdown.config", root, 1, "the markdownlint configuration cannot be read, so wip is not linted: %v", err)
		return
	}
	if cfg == nil {
		return
	}
	_ = filepath.WalkDir(c.repo.WipDir(), func(path string, e fs.DirEntry, err error) error {
		if err != nil || e.IsDir() || !strings.HasSuffix(path, ".md") {
			return nil //nolint:nilerr // an unreadable entry is skipped, as the other rules skip what they cannot read
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return nil //nolint:nilerr // as above
		}
		for _, f := range cfg.Lint(string(data)) {
			c.add(Warning, "markdown."+f.Rule, path, f.Line, "%s", f)
		}
		return nil
	})
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
