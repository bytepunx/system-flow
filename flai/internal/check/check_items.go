package check

import (
	"errors"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"sort"
	"strings"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

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
	if c.write != nil {
		c.write.pick(c)
	}
	for _, it := range c.items {
		if c.write == nil || c.write.items[it.ID] {
			c.oneItem(it)
		}
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
		if it.Type == workitem.Story && !it.Archived && (it.Status == workitem.InProgress || it.Status == workitem.Review) && c.write.story(it.ID) {
			if _, err := os.Stat(c.repo.NarrativePath(it.ID)); err != nil {
				c.add(Error, "narrative.missing", it.Path, keyLine(it.Path, "status"), "%s is %s but has no narrative; run flai stream open %s", it.ID, it.Status, it.ID)
			}
		}
	}
	matches, _ := filepath.Glob(filepath.Join(dir, "S-*.md"))
	for _, path := range matches {
		want := strings.TrimSuffix(filepath.Base(path), ".md")
		if !c.write.story(want) {
			// a write's run reads only its stories' narratives; the index
			// still counts this one as a stream (RunItemWrite)
			active[want] = true
			continue
		}
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
	if c.write != nil && !c.write.paths[filepath.Clean(indexPath)] {
		return // a write that did not write the index brings no finding on it
	}
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
