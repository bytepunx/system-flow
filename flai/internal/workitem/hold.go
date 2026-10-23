package workitem

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"slices"
	"sort"
	"strings"

	"github.com/bytepunx/system-flow/flai/internal/manifest"
)

// Holds on ready stories (S-0128, ADR-0046). A story's claim is its touches,
// each folder among them narrowed to the touches its tasks name inside it,
// and the touches of its open tasks outside them (ADR-0096). A ready story whose
// claim overlaps the claim of a story in progress is held, and so is one that
// names in after: a story that is not done (S-0130): flai serve does not start
// its agent and wait_for_work does not offer it. A story in review holds
// nothing (S-0295, ADR-0096): its branch is finished and synced, and the
// notice at acceptance tells an overlapping story what changed. An overlap
// that lies wholly inside a pattern of the manifest's claims.shared holds
// nothing either (ADR-0096): many stories change those paths in separate
// sections or new files. An overlap between a ready story and the story in
// progress that holds it, on paths that story's agent shared with it in their
// conversation, holds nothing while the share is in force (ADR-0134, share.go).
// A hold on overlap alone names each holding story whose agent flai asked
// about it in their open conversation since the ready story last entered
// ready (S-0338). The operator's own moves only warn.

// Hold reason codes.
const (
	HoldOverlap   = "overlap"    // its claim overlaps an open story's
	HoldNoTouches = "no-touches" // its claim or an open story's is empty
	HoldAfter     = "after"      // a story it names in after: is not done (S-0130)
)

// Hold is why a ready story waits: a reason code, and one line with the
// evidence and what clears it.
type Hold struct {
	Code   string `json:"code"`
	Reason string `json:"reason"`
	// Asked is, on a hold on overlap alone, each holding story whose agent
	// flai asked about the hold since the story last entered ready, in ID
	// order (S-0338).
	Asked []HoldAsk `json:"asked,omitempty"`
}

// HoldAsk is the newest message by flai that asked a holding story's agent
// about a hold (S-0338).
type HoldAsk struct {
	By           string `json:"by"`           // the story in progress whose agent was asked
	Conversation string `json:"conversation"` // the conversation it was asked in
	At           string `json:"at"`           // when, in TimeFormat
}

// Holds judges ready stories against the stories open at one reading of the
// repository.
type Holds struct {
	projects []manifest.Project
	shared   manifest.Claims    // an overlap wholly inside these holds nothing
	shares   []Share            // shares in force, IDs canonical (ADR-0134)
	asks     []Ask              // messages by flai asking about a hold, IDs canonical (S-0338)
	tasks    map[string][]*Item // tasks by story, cancelled ones left out
	open     []openClaim
	stories  map[string]*Item // every story given, archived ones included
	// lookup finds a story named in after: that was not given, such as an
	// archived one when only the active items were read.
	lookup func(id string) *Item
	memo   *HoldsMemo
}

// HoldsMemo remembers, for the Holds that share it, the claim of each story,
// whether the shared paths cover a path, and how one claim holds another, for
// a caller that judges the same stories again and again, such as the replay
// of the holds (I-0155). Its Holds must be given the same projects and shared
// paths, and its caller must Forget a story's claim when the story's tasks or
// their states change. The claims it gives are shared, not to be changed.
type HoldsMemo struct {
	claims  map[string][]memoClaim // by story ID, one per story item
	covered map[string]bool        // whether the shared paths cover a path
	pairs   map[claimPair]pairHold // how a ready story's claim is held by an open one's
	last    int                    // the number of the last claim worked out
}

type memoClaim struct {
	story *Item
	paths []string
	n     int // its number, from 1, a new one each time a claim is worked out
}

// claimPair is a ready story's claim and an open story's, by number, and how
// the open one is named.
type claimPair struct {
	held, holder int
	label        string
}

// pairHold is what holdBy says of a pair.
type pairHold struct {
	part, code string
}

// NewHoldsMemo is an empty HoldsMemo.
func NewHoldsMemo() *HoldsMemo {
	return &HoldsMemo{claims: map[string][]memoClaim{}, covered: map[string]bool{}, pairs: map[claimPair]pairHold{}}
}

// Forget has m forget the claim of each story with ID id.
func (m *HoldsMemo) Forget(id string) {
	delete(m.claims, id)
}

// NewHolds is NewHolds with the Holds remembering what they work out in m,
// the claims of the stories in progress among items too.
func (m *HoldsMemo) NewHolds(items []*Item, projects []manifest.Project) *Holds {
	return newHolds(items, projects, m)
}

type openClaim struct {
	id    string
	label string // how the reason names its state: in progress, its agent started
	paths []string
	n     int // the claim's number in the memo, 0 without one
}

// NewHolds reads from items the claims of the stories in progress, the only
// stories that hold (ADR-0096); projects are the manifest's sub-projects,
// whose names and tags a claim reads as their paths.
func NewHolds(items []*Item, projects []manifest.Project) *Holds {
	return newHolds(items, projects, nil)
}

// newHolds is NewHolds with memo, if any.
func newHolds(items []*Item, projects []manifest.Project, memo *HoldsMemo) *Holds {
	h := &Holds{projects: projects, tasks: map[string][]*Item{}, stories: map[string]*Item{}, memo: memo}
	for _, it := range items {
		if !it.Archived && it.Type == Task && it.Status != Cancelled {
			h.tasks[it.Parent] = append(h.tasks[it.Parent], it)
		}
		if it.Type == Story {
			h.stories[it.ID] = it
		}
	}
	for _, it := range items {
		if !it.Archived && it.Type == Story && it.Status == InProgress {
			h.Open(it, "in progress")
		}
	}
	return h
}

// Holds judges ready stories against items, as NewHolds does, with the
// manifest's shared paths as SharedClaims reads them, and the shares in force
// of the open conversations (ADR-0134) and their asks about a hold (S-0338),
// as Conversations reads them, and finds a story named in after: or in a
// share in the archive when items do not hold it.
func (r *Repo) Holds(items []*Item) *Holds {
	h := NewHolds(items, r.Manifest.Projects).WithShared(r.SharedClaims())
	h.lookup = func(id string) *Item {
		if it, err := r.Get(id); err == nil && it.Type == Story {
			return it
		}
		return nil
	}
	shares, asks := r.Conversations()
	return h.WithShares(shares).WithAsks(asks)
}

// SharedClaims is the manifest's claims as system-flow.yaml has them now,
// read again rather than taken from the manifest r was opened with: flai mcp
// keeps one Repo for as long as it runs, and flai shared or shared_paths_edit
// may have changed the list since. The manifest r was opened with gives them
// when the file cannot be read.
func (r *Repo) SharedClaims() manifest.Claims {
	m, err := manifest.Load(filepath.Join(r.Root, manifest.File))
	if err != nil {
		return r.Manifest.Claims
	}
	return m.Claims
}

// WithShared has h treat an overlap wholly inside a pattern of claims as no
// overlap (ADR-0096), and returns h.
func (h *Holds) WithShared(claims manifest.Claims) *Holds {
	h.shared = claims
	return h
}

// Overlaps says whether paths a and b, such as entries of two claims, overlap
// where it counts: they cover a common path, as PathsOverlap says, and the
// narrower of the two, the deeper one or either when they are equal, does not
// lie wholly inside a shared path (ADR-0096).
func (h *Holds) Overlaps(a, b string) bool {
	if !PathsOverlap(a, b) {
		return false
	}
	return !h.covers(narrower(a, b))
}

// covers says whether p lies wholly inside a shared path, as Covers says.
func (h *Holds) covers(p string) bool {
	if h.memo == nil {
		_, ok := h.shared.Covers(p)
		return ok
	}
	ok, seen := h.memo.covered[p]
	if !seen {
		_, ok = h.shared.Covers(p)
		h.memo.covered[p] = ok
	}
	return ok
}

// overlapsFor is Overlaps for mine, an entry of held's claim, and theirs, an
// entry of holder's: a pair whose narrower lies inside a path holder shared
// with held does not overlap either (ADR-0134).
func (h *Holds) overlapsFor(held, holder, mine, theirs string) bool {
	return h.Overlaps(mine, theirs) && !h.sharedBy(holder, held, narrower(mine, theirs))
}

// narrower is the deeper of two overlapping paths, the second when they are
// equal: one is the other or lies below it.
func narrower(a, b string) string {
	if len(strings.TrimSuffix(a, "/")) > len(strings.TrimSuffix(b, "/")) {
		return a
	}
	return b
}

// Open counts story as open from now on, named with label: flai serve's
// launcher opens a story whose agent it has started before the agent moves
// it to in progress, so that one look does not start two that overlap.
func (h *Holds) Open(story *Item, label string) {
	paths, n := h.claimOf(story)
	h.open = append(h.open, openClaim{id: story.ID, label: label, paths: paths, n: n})
	sort.SliceStable(h.open, func(i, j int) bool { return h.open[i].id < h.open[j].id })
}

// Claim is what story claims (ADR-0096), each entry read as a path first, a
// sub-project's when it names the sub-project or one of its tags, without
// duplicates. A story touch that holds touches of its tasks, done or open, is
// replaced by them; one that no task names inside, or that a task names whole,
// stays whole; and a touch of an open task outside every story touch is added.
// A done task's touch counts only where it narrows: the branch changed it.
func (h *Holds) Claim(story *Item) []string {
	paths, _ := h.claimOf(story)
	return paths
}

// claimOf is Claim and the claim's number in the memo, 0 without one.
func (h *Holds) claimOf(story *Item) ([]string, int) {
	if h.memo == nil {
		return h.claim(story), 0
	}
	for _, c := range h.memo.claims[story.ID] {
		if c.story == story {
			return c.paths, c.n
		}
	}
	h.memo.last++
	c := memoClaim{story: story, paths: h.claim(story), n: h.memo.last}
	h.memo.claims[story.ID] = append(h.memo.claims[story.ID], c)
	return c.paths, c.n
}

// claim is Claim worked out.
func (h *Holds) claim(story *Item) []string {
	var out []string
	seen := map[string]bool{}
	add := func(p string) {
		if p != "" && !seen[p] {
			seen[p] = true
			out = append(out, p)
		}
	}
	var mine, theirs, open []string // the story's touches, its tasks', its open tasks'
	for _, e := range story.Touches {
		mine = append(mine, h.path(e))
	}
	for _, t := range h.tasks[story.ID] {
		for _, e := range t.Touches {
			p := h.path(e)
			theirs = append(theirs, p)
			if !t.Closed() {
				open = append(open, p)
			}
		}
	}
	for _, s := range mine {
		var inside []string
		whole := false
		for _, p := range theirs {
			switch {
			case p == "" || !PathsOverlap(p, s):
			case strings.HasPrefix(p, s+"/"):
				inside = append(inside, p)
			default:
				whole = true // the task names the touch itself, or a folder around it
			}
		}
		if whole || len(inside) == 0 {
			add(s)
			continue
		}
		for _, p := range inside {
			add(p)
		}
	}
	for _, p := range open {
		if !within(p, mine) {
			add(p)
		}
	}
	return out
}

// within says whether p is one of entries or lies below one.
func within(p string, entries []string) bool {
	for _, e := range entries {
		if e != "" && (p == e || strings.HasPrefix(p, e+"/")) {
			return true
		}
	}
	return false
}

// path is a touches entry as it is compared: a component's name or tag
// becomes its path, and anything else is taken as it is.
func (h *Holds) path(entry string) string {
	return touchPath(h.projects, entry)
}

// touchPath is a touches entry read as a path, as Holds.path reads it with
// the manifest's sub-projects.
func touchPath(projects []manifest.Project, entry string) string {
	e := strings.TrimSuffix(strings.TrimPrefix(strings.TrimSpace(entry), "./"), "/")
	for _, p := range projects {
		if e == p.Name || contains(p.Tags, e) {
			return strings.TrimSuffix(p.Path, "/")
		}
	}
	return e
}

// notCounted are the folders whose files FolderTouch does not count, at any
// depth: version control, flai's cache with its worktrees, and installed
// dependencies, none of them files a story changes.
var notCounted = map[string]bool{".git": true, ".flai-cache": true, "node_modules": true}

// FolderTouch says whether a touches entry is a folder touch (ADR-0188), and
// how many files the folder holds when it is. The entry is read as a claim
// reads it: a sub-project's name or tag among projects is its path. It is a
// folder touch when it names a folder of the tree at root holding at least
// one file, counted at any depth. A path absent from the tree is a file to be
// made, and a folder lying wholly inside a pattern of shared, by ADR-0096's
// rule, holds nothing: neither is a folder touch.
func FolderTouch(root string, projects []manifest.Project, shared manifest.Claims, entry string) (files int, folder bool) {
	p := touchPath(projects, entry)
	if p == "" || strings.HasPrefix(p, "/") || slices.Contains(strings.Split(p, "/"), "..") {
		return 0, false
	}
	dir := filepath.Join(root, filepath.FromSlash(p))
	if fi, err := os.Lstat(dir); err != nil || !fi.IsDir() {
		return 0, false
	}
	if sharedFolder(shared, p) {
		return 0, false
	}
	_ = filepath.WalkDir(dir, func(_ string, d fs.DirEntry, err error) error {
		switch {
		case err != nil, d.IsDir() && notCounted[d.Name()]:
			// a folder that cannot be read is passed over, as one not counted
			return filepath.SkipDir
		case !d.IsDir():
			files++
		}
		return nil
	})
	return files, files > 0
}

// sharedFolder says whether the folder p lies wholly inside one pattern of
// shared. Covers reads an entry whose last segment has an extension as a
// file, so for a folder named like one, such as conf.d, the pattern must
// also cover a folder below it: docs/* matches docs/conf.d, but not all of it.
func sharedFolder(shared manifest.Claims, p string) bool {
	if _, ok := shared.Covers(p); !ok {
		return false
	}
	base := path.Base(p)
	if i := strings.LastIndex(base, "."); i <= 0 || i == len(base)-1 {
		return true // Covers read it as a folder already
	}
	for _, pattern := range shared.Shared {
		one := manifest.Claims{Shared: []string{pattern}}
		if _, ok := one.Covers(p); !ok {
			continue
		}
		if _, ok := one.Covers(p + "/x"); ok {
			return true
		}
	}
	return false
}

// FolderTouchError refuses a task's folder touch (ADR-0188). Its message is
// the one every write path gives.
type FolderTouchError struct {
	ID     string // the task
	Touch  string // the entry as given
	Folder string // the folder it names, a sub-project's path for its name or tag
	Files  int    // how many files the folder holds
}

// Error names the folder, its file count, and what to do instead.
func (e *FolderTouchError) Error() string {
	what := e.Touch + " is a folder"
	if e.Folder != e.Touch {
		what = e.Touch + " is the folder " + e.Folder + ","
	}
	return fmt.Sprintf("%s's touch %s holding %d %s: a task takes no folder touch (ADR-0188); name the files the task changes, or leave the folder to the story's touches with a ### Planning note",
		e.ID, what, e.Files, oneOrMany(e.Files, "file", "files"))
}

// CheckTaskTouches refuses the task id a folder touch (ADR-0188): the first
// entry of touches that before does not hold and that FolderTouch finds a
// folder of r's checkout, as a *FolderTouchError, or nil. before is the
// task's touches before the write, nil for a new task: an entry the task
// already had stays, so that a task written before the rule can still be
// edited, and a path made a folder since it was named is not refused.
func (r *Repo) CheckTaskTouches(id string, before, touches []string) error {
	had := map[string]bool{}
	for _, e := range before {
		had[touchPath(nil, e)] = true
	}
	var shared *manifest.Claims
	for _, e := range touches {
		if had[touchPath(nil, e)] {
			continue
		}
		if shared == nil {
			c := r.SharedClaims()
			shared = &c
		}
		if n, ok := FolderTouch(r.Root, r.Manifest.Projects, *shared, e); ok {
			return &FolderTouchError{ID: id, Touch: e, Folder: touchPath(r.Manifest.Projects, e), Files: n}
		}
	}
	return nil
}

// Of is why story is held, or nil when it is not. Every story that holds it
// is named, in ID order. A story held both by after: and by an overlap is
// held (after), and its reason says both. A hold on overlap alone names each
// holding story whose agent was asked about it (S-0338).
func (h *Holds) Of(story *Item) *Hold {
	hold := h.of(story)
	if hold != nil && hold.Code == HoldOverlap && len(h.asks) > 0 {
		hold.Asked = h.asked(story, h.overlapsBy(story, hold))
	}
	return hold
}

// of is Of without the asks.
func (h *Holds) of(story *Item) *Hold {
	after, overlap := h.after(story), h.overlap(story)
	switch {
	case after == nil:
		return overlap
	case overlap == nil:
		return after
	}
	return &Hold{Code: HoldAfter, Reason: after.Reason + "; also " + overlap.Reason}
}

// after is why story waits for the stories it names in after:, or nil when
// each of them is done. One that was cancelled keeps the hold, and so does
// one that does not exist: flai check reports it, and the operator decides.
func (h *Holds) after(story *Item) *Hold {
	var waits, ids, notes []string
	seen := map[string]bool{}
	for _, e := range story.After {
		id := CanonicalID(e)
		if id == story.ID || seen[id] {
			continue
		}
		seen[id] = true
		it := h.story(id)
		switch {
		case it == nil:
			waits = append(waits, id+" (no such story)")
			notes = append(notes, id+" names no story, so fix after:")
		case it.Status == Done:
			continue
		case it.Status == Cancelled:
			waits = append(waits, id+" (cancelled)")
			notes = append(notes, id+" was cancelled, so drop it from after: if "+story.ID+" no longer needs it")
		default:
			waits = append(waits, id+" ("+stateLabel(it.Status)+")")
		}
		ids = append(ids, id)
	}
	if len(waits) == 0 {
		return nil
	}
	reason := fmt.Sprintf("held (after): waits for %s; starts when %s %s done", and(waits), and(ids), oneOrMany(len(ids), "is", "are"))
	if len(notes) > 0 {
		reason += "; " + strings.Join(notes, "; ")
	}
	return &Hold{Code: HoldAfter, Reason: reason}
}

// story finds a story by ID among those given, else through lookup.
func (h *Holds) story(id string) *Item {
	if it, ok := h.stories[id]; ok {
		return it
	}
	if h.lookup == nil {
		return nil
	}
	it := h.lookup(id)
	h.stories[id] = it
	return it
}

// stateLabel is how a reason names a story's state.
func stateLabel(status string) string {
	switch status {
	case Backlog:
		return "in backlog"
	case InProgress:
		return "in progress"
	case Review:
		return "in review"
	}
	return status
}

// overlap is why story's claim is held by the open stories', or nil.
func (h *Holds) overlap(story *Item) *Hold {
	claim, n := h.claimOf(story)
	var others []openClaim
	for _, o := range h.open {
		if o.id != story.ID {
			others = append(others, o)
		}
	}
	if len(others) == 0 {
		return nil
	}
	if len(claim) == 0 {
		var named []string
		for _, o := range others {
			named = append(named, o.id+" ("+o.label+")")
		}
		return &Hold{Code: HoldNoTouches, Reason: fmt.Sprintf("held (no-touches): declares no touches, so it may change what %s %s; starts when it declares touches that overlap no story's in progress, or when %s", and(named), oneOrMany(len(others), "changes", "change"), clears(others))}
	}
	var code string
	var parts []string
	var by []openClaim
	for _, o := range others {
		part, c := h.heldBy(story.ID, claim, n, o)
		if part == "" {
			continue
		}
		if code == "" {
			code = c
		}
		parts = append(parts, part)
		by = append(by, o)
	}
	if len(parts) == 0 {
		return nil
	}
	return &Hold{Code: code, Reason: fmt.Sprintf("held (%s): %s; starts when %s", code, strings.Join(parts, "; "), clears(by))}
}

// heldBy is holdBy, remembered in the memo for claim, numbered n, and o's
// while no share is in force.
func (h *Holds) heldBy(held string, claim []string, n int, o openClaim) (part, code string) {
	if n == 0 || o.n == 0 || len(h.shares) > 0 {
		return h.holdBy(held, claim, o)
	}
	k := claimPair{held: n, holder: o.n, label: o.label}
	if r, ok := h.memo.pairs[k]; ok {
		return r.part, r.code
	}
	part, code = h.holdBy(held, claim, o)
	h.memo.pairs[k] = pairHold{part: part, code: code}
	return part, code
}

// holdBy says how claim, held's, overlaps o's, if it does, naming the first
// pair that holds: a pair whose overlap lies wholly inside a shared path, or
// inside a path o shared with held, does not, so the pair named is one still
// held.
func (h *Holds) holdBy(held string, claim []string, o openClaim) (part, code string) {
	who := o.id + " (" + o.label + ")"
	if len(o.paths) == 0 {
		return who + " declares no touches, so it may change anything", HoldNoTouches
	}
	mine, theirs, ok := h.firstHolding(held, claim, o)
	if !ok {
		return "", ""
	}
	switch {
	case mine == theirs:
		return fmt.Sprintf("touches %s, which %s touches too", mine, who), HoldOverlap
	case strings.HasPrefix(mine, theirs+"/"):
		return fmt.Sprintf("touches %s, inside %s which %s touches", mine, theirs, who), HoldOverlap
	default:
		return fmt.Sprintf("touches %s, which holds %s that %s touches", mine, theirs, who), HoldOverlap
	}
}

// holding is every pair of an entry of claim, held's, and one of o's that
// holds, in claim order, then o's.
func (h *Holds) holding(held string, claim []string, o openClaim) [][2]string {
	var out [][2]string
	for _, mine := range claim {
		for _, theirs := range o.paths {
			if h.overlapsFor(held, o.id, mine, theirs) {
				out = append(out, [2]string{mine, theirs})
			}
		}
	}
	return out
}

// firstHolding is the first pair holding gives, mine and theirs, and whether
// there is one.
func (h *Holds) firstHolding(held string, claim []string, o openClaim) (string, string, bool) {
	for _, mine := range claim {
		for _, theirs := range o.paths {
			if h.overlapsFor(held, o.id, mine, theirs) {
				return mine, theirs, true
			}
		}
	}
	return "", "", false
}

// clears says what ends a hold by these stories: each leaving in progress.
func clears(by []openClaim) string {
	ids := make([]string, len(by))
	for i, o := range by {
		ids[i] = o.id
	}
	if len(ids) == 1 {
		return ids[0] + " moves to review, is cancelled, or is sent back"
	}
	return and(ids) + " move to review, are cancelled, or are sent back"
}

// PathsOverlap says whether two touches entries cover a common path: equal,
// or one a /-bounded prefix of the other (ADR-0019).
func PathsOverlap(a, b string) bool {
	a, b = strings.TrimSuffix(a, "/"), strings.TrimSuffix(b, "/")
	if len(a) > len(b) {
		a, b = b, a
	}
	return strings.HasPrefix(b, a) && (len(a) == len(b) || b[len(a)] == '/')
}

func and(xs []string) string {
	switch len(xs) {
	case 0:
		return ""
	case 1:
		return xs[0]
	default:
		return strings.Join(xs[:len(xs)-1], ", ") + " and " + xs[len(xs)-1]
	}
}

func oneOrMany(n int, one, many string) string {
	if n == 1 {
		return one
	}
	return many
}
