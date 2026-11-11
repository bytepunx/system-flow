// Package check validates a conforming repository against the standard:
// work items, narratives, the board, and documentation front matter. It is
// the reference validator that flaiover mirrors.
package check

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"sync"
	"time"

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
	// children are each item's children by parent ID, in items' order, as
	// workitem.Children gives them; indexItems fills it, and with it nil
	// childrenOf scans items (S-0677).
	children map[string][]*workitem.Item
	res      *Result
	vocab    map[string]bool // topics a story can have (ADR-0047)
	// unread is the item files the list left out because they did not read
	// or parse (I-0344); indexItems reports them and fills it when nil, and
	// the rules that look an item up by ID pass over one it names.
	unread *unreadSet
	// write narrows the rules that read files to what a work item write
	// touches (RunItemWrite); nil runs them over the whole repository.
	write *writeScope
}

// rule is one of the rules Run runs; first marks indexItems, which fills the
// checker's byID and children that the rules after it read, so it runs before
// them, alone. The per-item rules of workItems, oneItem, then run on the pool
// in chunks of items, merged at first's place (S-0677).
type rule struct {
	run   func(*checker)
	first bool
}

// sequence is every rule Run runs, in the order their findings are merged:
// the order Run ran them in, one after another, before S-0601. sortFindings
// is stable, so findings at one path and line keep this order.
var sequence = []rule{
	{run: (*checker).layout},
	{run: (*checker).instructions},
	{run: (*checker).indexItems, first: true},
	{run: (*checker).narratives},
	{run: (*checker).activities},
	{run: (*checker).overlap},
	{run: (*checker).componentTag},
	{run: (*checker).after},
	{run: (*checker).touchesFolder},
	{run: (*checker).afterHub},
	{run: (*checker).afterReason},
	{run: (*checker).unaccepted},
	{run: (*checker).epicLags},
	{run: (*checker).board},
	{run: (*checker).documentation},
	{run: (*checker).adrIndex},
	{run: (*checker).conventions},
	{run: (*checker).roles},
	{run: (*checker).formulas},
	{run: (*checker).issues},
	{run: (*checker).experiments},
	{run: (*checker).analysis},
	{run: (*checker).digest},
	{run: (*checker).designTokens},
	{run: (*checker).componentExamples},
	{run: (*checker).threads},
	{run: (*checker).messages},
	{run: (*checker).markdown},
	{run: (*checker).conflictMarkers},
}

// Run executes every rule against the repo: indexItems first and alone, then
// the rest at once on at most GOMAXPROCS workers, each into a Result of its
// own, merged in sequence's order (S-0601). workItems' per-item rules run on
// the same workers in chunks of itemChunk items, merged in item order at
// indexItems' place, so the findings and counts are those of workItems run
// whole (S-0677). The rules only read, but for what each adds to its own
// Result; the topic vocabulary they share is worked out before.
//
// An item file the list leaves out because it does not read or parse is one
// item.unreadable finding on its path, and every other item is checked as
// before (I-0344): the rules that look an item up pass over the ID it names,
// and the markdown lint's findings on it are left out.
func Run(repo *workitem.Repo, now time.Time) (*Result, error) {
	items, err := repo.List(true)
	if err != nil {
		return nil, err
	}
	c := &checker{repo: repo, now: now.UTC(), items: items, byID: map[string]*workitem.Item{}, children: map[string][]*workitem.Item{}, res: &Result{Items: len(items)},
		unread: newUnreadSet(repo.Unreadable())}
	c.vocab = Vocabulary(repo, items)
	// slots holds each rule's Results in merge order: one for a rule, and for
	// the first the index's then one per chunk of items.
	type job struct {
		slot, part int
		run        func(*checker)
	}
	slots := make([][]*Result, len(sequence))
	var jobs []job
	for i, r := range sequence {
		if !r.first {
			slots[i] = make([]*Result, 1)
			jobs = append(jobs, job{slot: i, run: r.run})
			continue
		}
		slots[i] = []*Result{c.runRule(r.run)}
		for _, chunk := range chunks(c.items, itemChunk) {
			slots[i] = append(slots[i], nil)
			jobs = append(jobs, job{slot: i, part: len(slots[i]) - 1, run: func(c *checker) { c.oneItems(chunk) }})
		}
	}
	next := make(chan job)
	var wg sync.WaitGroup
	for range min(runtime.GOMAXPROCS(0), len(jobs)) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := range next {
				slots[j.slot][j.part] = c.runRule(j.run)
			}
		}()
	}
	for _, j := range jobs {
		next <- j
	}
	close(next)
	wg.Wait()
	for _, rs := range slots {
		for _, r := range rs {
			c.res.absorb(r)
		}
	}
	c.res.dropLintOnUnread(c.repo, c.unread)
	sortFindings(c.res.Findings)
	return c.res, nil
}

// dropLintOnUnread leaves out the markdown lint's findings on the item files
// in unread, and takes back their counts: a file with no front matter breaks
// the lint's rules too, and its item.unreadable says what to fix (I-0344). A
// merge conflict marker stays, since it is what keeps the file from parsing.
func (r *Result) dropLintOnUnread(repo *workitem.Repo, unread *unreadSet) {
	if unread.empty() {
		return
	}
	kept := r.Findings[:0]
	for _, f := range r.Findings {
		path := f.Path
		if !filepath.IsAbs(path) {
			path = filepath.Join(repo.Root, path)
		}
		if strings.HasPrefix(f.Rule, "markdown.MD") && unread.paths[filepath.Clean(path)] {
			r.drop(f)
			continue
		}
		kept = append(kept, f)
	}
	clear(r.Findings[len(kept):])
	r.Findings = kept
}

// itemChunk is how many items one job of Run's per-item rules checks: small
// enough that the chunks spread over the workers beside the slower rules.
var itemChunk = 128

// chunks cuts items into runs of at most n, in order.
func chunks(items []*workitem.Item, n int) [][]*workitem.Item {
	var out [][]*workitem.Item
	for len(items) > n {
		out = append(out, items[:n:n])
		items = items[n:]
	}
	if len(items) > 0 {
		out = append(out, items)
	}
	return out
}

// runRule runs one rule on a copy of c with a Result of its own, so that what
// add, advise, and the rules that set Stories write by the last finding's
// index stays in that rule's findings, and returns that Result. The copy
// shares c's items, byID, children, and vocab, which no rule but indexItems
// writes.
func (c *checker) runRule(run func(*checker)) *Result {
	rc := *c
	rc.res = &Result{}
	run(&rc)
	return rc.res
}

// absorb adds o's findings after r's and o's counts to r's.
func (r *Result) absorb(o *Result) {
	r.Findings = append(r.Findings, o.Findings...)
	r.Errors += o.Errors
	r.Warnings += o.Warnings
	r.Advisory += o.Advisory
	r.Outside += o.Outside
	r.outsideErrors += o.outsideErrors
	r.outsideWarnings += o.outsideWarnings
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
