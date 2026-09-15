package workitem

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/goccy/go-yaml"
)

// BoardFile is wip/kanban/board.md.
const BoardFile = "board.md"

// Board is the parsed board.md: WIP limits, pull order, and the markdown body.
type Board struct {
	Title     string         `yaml:"title"`
	Updated   string         `yaml:"updated"`
	Status    string         `yaml:"status"`
	WIPLimits map[string]int `yaml:"wip_limits"`
	Order     []string       `yaml:"order"`

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
	for _, k := range []string{Ready, InProgress, Review} {
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
	sb.WriteString("---\n")
	sb.WriteString(b.Body)
	if err := os.MkdirAll(filepath.Dir(b.Path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(b.Path, []byte(sb.String()), 0o644)
}

// RemoveFromOrder drops id from the pull order.
func (b *Board) RemoveFromOrder(id string) {
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

func orDefault(s, d string) string {
	if s == "" {
		return d
	}
	return s
}
