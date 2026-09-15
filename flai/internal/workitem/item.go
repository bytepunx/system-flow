// Package workitem reads, validates, transitions, and writes the epics,
// stories, tasks, and narratives under wip/. See design/system/work-hierarchy.md,
// workflow.md, and agent-narrative.md.
package workitem

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/goccy/go-yaml"
)

// Item types.
const (
	Epic  = "epic"
	Story = "story"
	Task  = "task"
)

// States.
const (
	Backlog    = "backlog"
	Ready      = "ready"
	InProgress = "in-progress"
	Review     = "review"
	Done       = "done"
	Cancelled  = "cancelled"
)

// Types, States, and Natures are the closed lists from ADR-0003 and ADR-0004.
var (
	Types   = []string{Epic, Story, Task}
	States  = []string{Backlog, Ready, InProgress, Review, Done, Cancelled}
	Natures = []string{"feature", "improvement", "remediation", "research", "experiment"}
)

// TimeFormat is the timestamp format used in front matter.
const TimeFormat = "2006-01-02T15:04:05Z"

// Transition is one state change.
type Transition struct {
	To string `yaml:"to" json:"to"`
	At string `yaml:"at" json:"at"`
	By string `yaml:"by" json:"by"`
}

// Block is one blocked interval; Until is empty while open.
type Block struct {
	From   string `yaml:"from" json:"from"`
	Until  string `yaml:"until" json:"until"`
	Reason string `yaml:"reason" json:"reason"`
}

// Item is a work item: front matter plus the markdown body.
type Item struct {
	ID          string       `yaml:"id" json:"id"`
	Type        string       `yaml:"type" json:"type"`
	Nature      string       `yaml:"nature" json:"nature"`
	Title       string       `yaml:"title" json:"title"`
	Status      string       `yaml:"status" json:"status"`
	Parent      string       `yaml:"parent" json:"parent"`
	Owner       string       `yaml:"owner" json:"owner"`
	Created     string       `yaml:"created" json:"created"`
	Updated     string       `yaml:"updated" json:"updated"`
	Transitions []Transition `yaml:"transitions" json:"transitions"`
	Blocked     []Block      `yaml:"blocked" json:"blocked"`
	Estimate    string       `yaml:"estimate" json:"estimate"`
	Stream      string       `yaml:"stream" json:"stream"`
	Tags        []string     `yaml:"tags" json:"tags"`

	Path     string `yaml:"-" json:"path"`     // file on disk
	Archived bool   `yaml:"-" json:"archived"` // lives under wip/archive
	Body     string `yaml:"-" json:"body"`     // markdown after the front matter
}

// ErrNoFrontMatter is returned for files without a leading --- block.
var ErrNoFrontMatter = errors.New("no front matter")

// SplitFrontMatter separates a markdown document into its YAML block and body.
func SplitFrontMatter(doc string) (fm, body string, err error) {
	if !strings.HasPrefix(doc, "---\n") {
		return "", "", ErrNoFrontMatter
	}
	rest := doc[4:]
	idx := strings.Index(rest, "\n---\n")
	if idx < 0 {
		if strings.HasSuffix(rest, "\n---") {
			return rest[:len(rest)-4], "", nil
		}
		return "", "", errors.New("unterminated front matter")
	}
	return rest[:idx+1], rest[idx+5:], nil
}

// ParseItem decodes a work item document.
func ParseItem(doc string) (*Item, error) {
	fm, body, err := SplitFrontMatter(doc)
	if err != nil {
		return nil, err
	}
	var it Item
	if err := yaml.UnmarshalWithOptions([]byte(fm), &it, yaml.Strict()); err != nil {
		return nil, err
	}
	it.Body = body
	return &it, nil
}

// ReadItem loads an item from disk.
func ReadItem(path string) (*Item, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	it, err := ParseItem(string(data))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", path, err)
	}
	it.Path = path
	return it, nil
}

// Validate checks the closed lists and structural rules that do not need
// other items.
func (it *Item) Validate() error {
	var errs []string
	if !idPattern.MatchString(it.ID) {
		errs = append(errs, fmt.Sprintf("id %q must look like E-001, S-001, or T-001", it.ID))
	}
	if !contains(Types, it.Type) {
		errs = append(errs, fmt.Sprintf("type %q must be one of %s", it.Type, strings.Join(Types, ", ")))
	} else if it.ID != "" && !strings.HasPrefix(it.ID, strings.ToUpper(it.Type[:1])+"-") {
		errs = append(errs, fmt.Sprintf("id %s does not match type %s", it.ID, it.Type))
	}
	if !contains(Natures, it.Nature) {
		errs = append(errs, fmt.Sprintf("nature %q must be one of %s", it.Nature, strings.Join(Natures, ", ")))
	}
	if !contains(States, it.Status) {
		errs = append(errs, fmt.Sprintf("status %q must be one of %s", it.Status, strings.Join(States, ", ")))
	}
	if it.Title == "" {
		errs = append(errs, "title is required")
	}
	if it.Type == Epic && it.Parent != "" {
		errs = append(errs, "epics have no parent")
	}
	if it.Type != Epic && it.Parent == "" {
		errs = append(errs, "parent is required for stories and tasks")
	}
	if it.Type == Story && it.Parent != "" && !strings.HasPrefix(it.Parent, "E-") {
		errs = append(errs, "a story's parent must be an epic")
	}
	if it.Type == Task && it.Parent != "" && !strings.HasPrefix(it.Parent, "S-") {
		errs = append(errs, "a task's parent must be a story")
	}
	for _, name := range []string{"created", "updated"} {
		v := it.Created
		if name == "updated" {
			v = it.Updated
		}
		if _, err := time.Parse(TimeFormat, v); err != nil {
			errs = append(errs, fmt.Sprintf("%s %q is not a UTC timestamp like 2026-09-15T16:10:00Z", name, v))
		}
	}
	for i, tr := range it.Transitions {
		if !contains(States, tr.To) || tr.To == Backlog {
			errs = append(errs, fmt.Sprintf("transitions[%d].to %q is not a state that can be entered", i, tr.To))
		}
		if _, err := time.Parse(TimeFormat, tr.At); err != nil {
			errs = append(errs, fmt.Sprintf("transitions[%d].at %q is not a UTC timestamp", i, tr.At))
		}
	}
	if want := it.lastState(); it.Status != want {
		errs = append(errs, fmt.Sprintf("status %q does not match the last transition (%s)", it.Status, want))
	}
	open := 0
	for i, b := range it.Blocked {
		if _, err := time.Parse(TimeFormat, b.From); err != nil {
			errs = append(errs, fmt.Sprintf("blocked[%d].from is not a UTC timestamp", i))
		}
		if b.Until == "" {
			open++
		} else if _, err := time.Parse(TimeFormat, b.Until); err != nil {
			errs = append(errs, fmt.Sprintf("blocked[%d].until is not a UTC timestamp", i))
		}
		if b.Reason == "" {
			errs = append(errs, fmt.Sprintf("blocked[%d].reason is required", i))
		}
	}
	if open > 1 {
		errs = append(errs, "more than one open blocked interval")
	}
	if it.Estimate != "" {
		if _, err := time.ParseDuration(it.Estimate); err != nil {
			errs = append(errs, fmt.Sprintf("estimate %q is not a Go duration like 4h or 90m", it.Estimate))
		}
	}
	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "; "))
	}
	return nil
}

func (it *Item) lastState() string {
	if len(it.Transitions) == 0 {
		return Backlog
	}
	return it.Transitions[len(it.Transitions)-1].To
}

// EnteredAt returns when the item entered its current state.
func (it *Item) EnteredAt() time.Time {
	if len(it.Transitions) == 0 {
		t, _ := time.Parse(TimeFormat, it.Created)
		return t
	}
	t, _ := time.Parse(TimeFormat, it.Transitions[len(it.Transitions)-1].At)
	return t
}

// FirstAt returns the time of the first transition into state, or zero.
func (it *Item) FirstAt(state string) time.Time {
	for _, tr := range it.Transitions {
		if tr.To == state {
			t, _ := time.Parse(TimeFormat, tr.At)
			return t
		}
	}
	return time.Time{}
}

// IsBlocked reports whether a blocked interval is open.
func (it *Item) IsBlocked() bool {
	for _, b := range it.Blocked {
		if b.Until == "" {
			return true
		}
	}
	return false
}

// Closed reports whether the item is done or cancelled.
func (it *Item) Closed() bool {
	return it.Status == Done || it.Status == Cancelled
}

// Marshal renders the item back to a markdown document. The front matter is
// written by hand so that timestamps stay unquoted and sequences indented,
// matching the hand-written style in the standard.
func (it *Item) Marshal() string {
	var b strings.Builder
	b.WriteString("---\n")
	fmt.Fprintf(&b, "id: %s\n", it.ID)
	fmt.Fprintf(&b, "type: %s\n", it.Type)
	fmt.Fprintf(&b, "nature: %s\n", it.Nature)
	fmt.Fprintf(&b, "title: %s\n", Scalar(it.Title))
	fmt.Fprintf(&b, "status: %s\n", it.Status)
	if it.Parent != "" {
		fmt.Fprintf(&b, "parent: %s\n", it.Parent)
	}
	fmt.Fprintf(&b, "owner: %s\n", Scalar(it.Owner))
	fmt.Fprintf(&b, "created: %s\n", it.Created)
	fmt.Fprintf(&b, "updated: %s\n", it.Updated)
	if len(it.Transitions) == 0 {
		b.WriteString("transitions: []\n")
	} else {
		b.WriteString("transitions:\n")
		for _, tr := range it.Transitions {
			fmt.Fprintf(&b, "  - to: %s\n    at: %s\n    by: %s\n", tr.To, tr.At, Scalar(tr.By))
		}
	}
	if len(it.Blocked) > 0 {
		b.WriteString("blocked:\n")
		for _, bl := range it.Blocked {
			fmt.Fprintf(&b, "  - from: %s\n", bl.From)
			if bl.Until != "" {
				fmt.Fprintf(&b, "    until: %s\n", bl.Until)
			}
			fmt.Fprintf(&b, "    reason: %s\n", Scalar(bl.Reason))
		}
	}
	if it.Estimate != "" {
		fmt.Fprintf(&b, "estimate: %s\n", it.Estimate)
	}
	if it.Stream != "" {
		fmt.Fprintf(&b, "stream: %s\n", it.Stream)
	}
	fmt.Fprintf(&b, "tags: %s\n", FlowList(it.Tags))
	b.WriteString("---\n")
	b.WriteString(it.Body)
	return b.String()
}

var (
	idPattern    = regexp.MustCompile(`^[EST]-\d{3,}$`)
	plainPattern = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9 _./,()'!?&+@-]*$`)
	reserved     = map[string]bool{"true": true, "false": true, "null": true, "yes": true, "no": true, "on": true, "off": true, "~": true}
)

// Scalar renders a string as a YAML scalar, plain when safe and double-quoted
// otherwise.
func Scalar(s string) string {
	if s == "" {
		return `""`
	}
	if plainPattern.MatchString(s) && !strings.Contains(s, ": ") && !strings.Contains(s, " #") &&
		!reserved[strings.ToLower(s)] && !looksNumeric(s) && !strings.HasSuffix(s, ":") {
		return s
	}
	return strconv.Quote(s)
}

// FlowList renders a string list as [a, b].
func FlowList(items []string) string {
	if len(items) == 0 {
		return "[]"
	}
	parts := make([]string, len(items))
	for i, s := range items {
		parts[i] = Scalar(s)
	}
	return "[" + strings.Join(parts, ", ") + "]"
}

func looksNumeric(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

func contains(list []string, s string) bool {
	for _, x := range list {
		if x == s {
			return true
		}
	}
	return false
}

// TypeOfID returns the item type implied by an ID prefix.
func TypeOfID(id string) string {
	switch {
	case strings.HasPrefix(id, "E-"):
		return Epic
	case strings.HasPrefix(id, "S-"):
		return Story
	case strings.HasPrefix(id, "T-"):
		return Task
	}
	return ""
}
