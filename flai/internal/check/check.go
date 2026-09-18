// Package check validates a conforming repository against the standard:
// work items, narratives, the board, and documentation front matter. It is
// the reference validator that flaiover mirrors.
package check

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/goccy/go-yaml"

	"github.com/bytepunx/system-flow/flai/internal/conventions"
	"github.com/bytepunx/system-flow/flai/internal/issues"
	"github.com/bytepunx/system-flow/flai/internal/threads"
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
}

// Result is a whole run.
type Result struct {
	Findings []Finding `json:"findings"`
	Errors   int       `json:"errors"`
	Warnings int       `json:"warnings"`
	Items    int       `json:"items"`
}

// OK reports whether the run passed: no errors, and no warnings when strict.
func (r *Result) OK(strict bool) bool {
	return r.Errors == 0 && (!strict || r.Warnings == 0)
}

type checker struct {
	repo  *workitem.Repo
	now   time.Time
	items []*workitem.Item
	byID  map[string]*workitem.Item
	res   *Result
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
	c.overlap()
	c.board()
	c.documentation()
	c.conventions()
	c.issues()
	c.threads()
	sort.SliceStable(c.res.Findings, func(i, j int) bool {
		a, b := c.res.Findings[i], c.res.Findings[j]
		if a.Path != b.Path {
			return a.Path < b.Path
		}
		return a.Line < b.Line
	})
	return c.res, nil
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

var allowedNext = map[string][]string{
	workitem.Backlog:    {workitem.Ready, workitem.Cancelled},
	workitem.Ready:      {workitem.InProgress, workitem.Cancelled},
	workitem.InProgress: {workitem.Review, workitem.Cancelled, workitem.Done},
	workitem.Review:     {workitem.Done, workitem.InProgress},
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
	if err := it.Validate(); err != nil {
		for _, msg := range strings.Split(err.Error(), "; ") {
			key := strings.Fields(msg)[0]
			key = strings.SplitN(key, "[", 2)[0]
			c.add(Error, "item.front-matter", p, keyLine(p, key), "%s", msg)
		}
	}
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
	if it.Type == workitem.Story && it.Status != workitem.Backlog && it.Status != workitem.Cancelled {
		if len(children) == 0 {
			c.add(Error, "story.tasks", p, keyLine(p, "status"), "a %s story needs at least one task", it.Status)
		}
		if !hasCriteria(it.Body) {
			c.add(Error, "story.criteria", p, headingLine(p, "## Acceptance criteria"), "a %s story needs acceptance criteria with at least one checkbox", it.Status)
		}
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
		if !contains(allowedNext[prevState], tr.To) {
			c.add(Error, "item.sequence", p, keyLine(p, "transitions"), "transitions[%d]: %s cannot follow %s", i, tr.To, prevState)
		} else if tr.To == workitem.Done && prevState == workitem.InProgress && it.Type != workitem.Task {
			c.add(Error, "item.sequence", p, keyLine(p, "transitions"), "transitions[%d]: a %s must go through review before done", i, it.Type)
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

var narrativeSections = []string{"## Context", "## Current state", "## Next steps", "## Decisions", "## Open questions", "## Log"}

// overlap warns when two in-progress items declare touches that cover the
// same path (ADR-0019); it is advisory, humans and agents coordinate.
func (c *checker) overlap() {
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
			if pathsOverlap(a.path, b.path) {
				c.add(Warning, "wip.overlap", a.it.Path, keyLine(a.it.Path, "touches"), "%s touches %s, which %s (in progress) also touches as %s", a.it.ID, a.path, b.it.ID, b.path)
			}
		}
	}
}

func pathsOverlap(a, b string) bool {
	a, b = strings.TrimSuffix(a, "/"), strings.TrimSuffix(b, "/")
	return a == b || strings.HasPrefix(a, b+"/") || strings.HasPrefix(b, a+"/")
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
		for _, sec := range narrativeSections {
			if !strings.Contains(n.Body, "\n"+sec+"\n") {
				c.add(Warning, "narrative.section", path, 1, "missing %q section", sec)
			}
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
		if limit := b.WIPLimits[st]; limit > 0 && counts[st] > limit {
			c.add(Warning, "board.wip-limit", b.Path, keyLine(b.Path, "wip_limits"), "%d stories in %s, limit %d", counts[st], st, limit)
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

func (c *checker) documentation() {
	for _, key := range []string{"design", "docs"} {
		root := c.repo.Manifest.Dir(c.repo.Root, key)
		_ = filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
			if err != nil {
				return nil //nolint:nilerr // an unreadable entry is skipped, the walk continues
			}
			if d.IsDir() {
				if (d.Name() == conventions.Folder || d.Name() == issues.Folder) && filepath.Dir(path) == root {
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
			if isADR {
				num := adrName.FindStringSubmatch(d.Name())[1]
				if f.ID != "ADR-"+num {
					c.add(Warning, "adr.id", path, keyLine(path, "id"), "id should be ADR-%s to match the file name", num)
				}
				if f.Status == "" || f.Date == "" || f.Title == "" {
					c.add(Warning, "adr.front-matter", path, 1, "ADRs need id, title, status, and date")
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

func contains(list []string, s string) bool {
	for _, x := range list {
		if x == s {
			return true
		}
	}
	return false
}

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
		if base := filepath.Base(th.Path); !strings.HasPrefix(base, th.ID+"-") {
			c.add(Error, "threads.filename", th.Path, 1, "file name should start with %s-", th.ID)
		}
		abs := filepath.Join(c.repo.Root, filepath.FromSlash(th.Anchor.Path))
		data, err := os.ReadFile(abs)
		if err != nil {
			if d2, err2 := os.ReadFile(filepath.Join(c.repo.MainRoot, filepath.FromSlash(th.Anchor.Path))); err2 == nil {
				data, err = d2, nil
			}
		}
		if err != nil {
			c.add(Error, "threads.anchor", th.Path, keyLine(th.Path, "anchor"), "anchor %s does not exist", th.Anchor.Path)
		} else if th.Anchor.Heading != "" && !threads.HasHeading(string(data), th.Anchor.Heading) {
			c.add(Warning, "threads.heading", th.Path, keyLine(th.Path, "anchor"), "heading %q is no longer in %s", th.Anchor.Heading, th.Anchor.Path)
		}
		if th.Anchor.Item != "" {
			it, ok := byID[th.Anchor.Item]
			switch {
			case !ok:
				c.add(Error, "threads.anchor", th.Path, keyLine(th.Path, "anchor"), "item %s does not exist", th.Anchor.Item)
			case it.Archived && th.Open():
				c.add(Warning, "threads.archived", th.Path, keyLine(th.Path, "status"), "%s is %s but %s is archived; resolve it or move it", th.ID, th.Status, th.Anchor.Item)
			}
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
		if base := filepath.Base(is.Path); !strings.HasPrefix(base, is.ID+"-") {
			c.add(Error, "issues.filename", is.Path, 1, "file name should start with %s-", is.ID)
		}
	}
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
