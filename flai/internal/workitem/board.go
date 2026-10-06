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

// Save writes board.md.
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
