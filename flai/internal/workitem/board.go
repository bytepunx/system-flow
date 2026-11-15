package workitem

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/goccy/go-yaml"

	"github.com/bytepunx/system-flow/flai/internal/atomicfile"
	"github.com/bytepunx/system-flow/flai/internal/workflow"
)

// BoardFile is wip/kanban/board.md.
const BoardFile = "board.md"

// Board is the parsed board.md: WIP limits, pull order, who placed a story
// in it by hand, and the markdown body.
type Board struct {
	Title     string         `yaml:"title"`
	Updated   string         `yaml:"updated"`
	Status    string         `yaml:"status"`
	WIPLimits map[string]int `yaml:"wip_limits"`
	Order     []string       `yaml:"order"`
	// Placed is, by story, the last placement flai order made of it
	// (S-0219). A story's entry goes when it leaves its column.
	Placed map[string]Placed `yaml:"placed"`

	Body string `yaml:"-"`
	Path string `yaml:"-"`
}

// DefaultWIPLimits are the workflow.md defaults.
var DefaultWIPLimits = map[string]int{Ready: 5, InProgress: 2, Review: 3}

// LoadBoard reads board.md, returning defaults when it does not exist.
func (r *Repo) LoadBoard() (*Board, error) {
	path := filepath.Join(r.KanbanDir(), BoardFile)
	b := &Board{Path: path, Title: "Board", Status: "active", WIPLimits: map[string]int{}}
	for k, v := range DefaultWIPLimits {
		b.WIPLimits[k] = v
	}
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		b.Body = "\n# Board\n"
		return b, nil
	}
	if err != nil {
		return nil, err
	}
	fm, body, err := SplitFrontMatter(string(data))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", path, err)
	}
	if err := yaml.Unmarshal([]byte(fm), b); err != nil {
		return nil, fmt.Errorf("%s: %w", path, err)
	}
	b.Body = body
	b.Path = path
	if b.WIPLimits == nil {
		b.WIPLimits = map[string]int{}
	}
	return b, nil
}

// LockBoard takes the one lock every write of board.md is made under, an
// operating system lock on .flai-cache/board.lock, so that two writers do not
// each read the board, change their copy, and write it, the second over the
// first (I-0375). It is apart from the items lock and the checked-write lock:
// it is taken before the items lock, and never while the items lock or the
// checked-write lock is held (ADR-0306). A holder reads the board after it
// takes the lock, never before.
func (r *Repo) LockBoard() (unlock func(), err error) {
	return lockFile(filepath.Join(r.CacheDir(), "board.lock"), itemsLockWait)
}

// UpdateBoard is a write of board.md under the board lock (ADR-0306): it
// takes the lock, reads the active items and the board inside it, lets change
// make its change to that board, normalises the pull order against those
// items, and saves the board, dated today, before it lets the lock go. An
// error from change writes nothing and is returned as it is, so a rule
// refusal stays one.
func (r *Repo) UpdateBoard(today string, change func(b *Board, items []*Item) error) error {
	unlock, err := r.LockBoard()
	if err != nil {
		return err
	}
	defer unlock()
	items, err := r.List(false)
	if err != nil {
		return err
	}
	b, err := r.LoadBoard()
	if err != nil {
		return err
	}
	if err := change(b, items); err != nil {
		return err
	}
	b.NormalizeOrder(r.Flow(), items, r.unparsedStories())
	return b.Save(today)
}

// NormalizeOrder drops from the pull order each name that is not an active
// story in a state flow orders, ready and backlog on the baseline, and its
// placement with it, and every other placement of a name that is not such a
// story, and keeps the other names in their order (ADR-0306 §4).
// A name in unparsed, a story whose file the list left out because it did
// not read or parse, is kept, since its state cannot be told. A name the
// order repeats is kept once, where it first comes.
func (b *Board) NormalizeOrder(flow workflow.Workflow, items []*Item, unparsed map[string]bool) {
	ordered := map[string]bool{}
	for _, it := range items {
		if it.Archived || it.Type != Story {
			continue
		}
		if s, ok := flow.State(it.Status); ok && s.Order {
			ordered[it.ID] = true
		}
	}
	seen := map[string]bool{}
	out := make([]string, 0, len(b.Order))
	for _, id := range b.Order {
		if seen[id] {
			continue
		}
		seen[id] = true
		if ordered[id] || unparsed[CanonicalID(id)] {
			out = append(out, id)
			continue
		}
		delete(b.Placed, id)
	}
	b.Order = out
	// A placement a rewrite of the list has already dropped the name of,
	// as Place drops a stale one, goes too.
	for id := range b.Placed {
		if !ordered[id] && !unparsed[CanonicalID(id)] {
			delete(b.Placed, id)
		}
	}
}

// unparsedStories is the IDs, by their file names, of the active story files
// the latest list left out because they could not be read or did not parse
// (Unreadable): the names NormalizeOrder keeps.
func (r *Repo) unparsedStories() map[string]bool {
	dir := filepath.Clean(r.ItemDir(Story, false))
	out := map[string]bool{}
	for _, u := range r.Unreadable() {
		if filepath.Clean(filepath.Dir(u.Path)) != dir {
			continue
		}
		if id := fileID.FindString(filepath.Base(u.Path)); id != "" {
			out[CanonicalID(id)] = true
		}
	}
	return out
}

// Save writes board.md as it is. It is the write a holder of the board lock
// makes, through UpdateBoard or a move holding LockBoard, after it has read
// the board inside the lock (ADR-0306); it takes no lock and does not
// normalise the pull order by itself.
func (b *Board) Save(today string) error {
	var sb strings.Builder
	sb.WriteString("---\n")
	fmt.Fprintf(&sb, "title: %s\n", Scalar(b.Title))
	fmt.Fprintf(&sb, "updated: %s\n", today)
	fmt.Fprintf(&sb, "status: %s\n", orDefault(b.Status, "active"))
	sb.WriteString("wip_limits:\n")
	for _, k := range LimitedColumns {
		if v, ok := b.WIPLimits[k]; ok {
			fmt.Fprintf(&sb, "  %s: %d\n", k, v)
		}
	}
	if len(b.Order) == 0 {
		sb.WriteString("order: []\n")
	} else {
		sb.WriteString("order:\n")
		for _, id := range b.Order {
			fmt.Fprintf(&sb, "  - %s\n", id)
		}
	}
	if len(b.Placed) > 0 {
		sb.WriteString("placed:\n")
		ids := make([]string, 0, len(b.Placed))
		for id := range b.Placed {
			ids = append(ids, id)
		}
		sort.Slice(ids, func(i, j int) bool { return lessID(ids[i], ids[j]) })
		for _, id := range ids {
			p := b.Placed[id]
			fmt.Fprintf(&sb, "  %s:\n    by: %s\n    at: %s\n", id, Scalar(p.By), p.At)
		}
	}
	sb.WriteString("---\n")
	sb.WriteString(b.Body)
	if err := os.MkdirAll(filepath.Dir(b.Path), 0o755); err != nil {
		return err
	}
	return atomicfile.WriteFile(b.Path, []byte(sb.String()), 0o644)
}

// LimitedColumns are the columns that carry a WIP limit (workflow.md).
var LimitedColumns = []string{Ready, InProgress, Review}

// SetLimit sets a column's WIP limit; 0 means none (S-0167).
func (b *Board) SetLimit(column string, n int) error {
	if !contains(LimitedColumns, column) {
		return fmt.Errorf("%s has no WIP limit; only %s do", column, strings.Join(LimitedColumns, ", "))
	}
	if n < 0 {
		return fmt.Errorf("a WIP limit is a whole number, 0 for none; got %d", n)
	}
	b.WIPLimits[column] = n
	return nil
}

// RemoveFromOrder drops id from the pull order, and its placement with it: a
// story flai move takes out of its column is no longer where it was placed.
func (b *Board) RemoveFromOrder(id string) {
	delete(b.Placed, id)
	out := b.Order[:0]
	for _, x := range b.Order {
		if x != id {
			out = append(out, x)
		}
	}
	b.Order = out
}

// AppendToOrder adds id to the end of the pull order if absent.
func (b *Board) AppendToOrder(id string) {
	for _, x := range b.Order {
		if x == id {
			return
		}
	}
	b.Order = append(b.Order, id)
}

// Placed is who placed a story in the pull order by hand, with flai order or
// the dashboard's drag, which runs it, and when.
type Placed struct {
	By string `yaml:"by" json:"by"`
	// At is in TimeFormat.
	At string `yaml:"at" json:"at"`
}

// RecordPlacement records that by placed id at now, over any earlier
// placement of it.
func (b *Board) RecordPlacement(id, by string, now time.Time) {
	if b.Placed == nil {
		b.Placed = map[string]Placed{}
	}
	b.Placed[id] = Placed{By: by, At: now.UTC().Format(TimeFormat)}
}

// IsOrchestrator reports whether a placer is the orchestrator, by the name
// flai serve runs it under (ADR-0087).
func IsOrchestrator(by string) bool { return by == ActivityOrchestrator }

// HandPlaced is id's placement when someone other than the orchestrator made
// it within keep of now. There is none when keep is not positive.
func (b *Board) HandPlaced(id string, keep time.Duration, now time.Time) (Placed, bool) {
	p, ok := b.Placed[id]
	if !ok || keep <= 0 || IsOrchestrator(p.By) {
		return Placed{}, false
	}
	at, err := time.Parse(TimeFormat, p.At)
	if err != nil || at.Before(now.Add(-keep)) {
		return Placed{}, false
	}
	return p, true
}

func orDefault(s, d string) string {
	if s == "" {
		return d
	}
	return s
}
