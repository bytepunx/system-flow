// Package context builds the context pack an agent working a story is primed
// with (ADR-0047). This part selects the conventions: each file in read order
// with the sections whose topics miss the story's left out, and a list of
// what was left out so the agent knows the rule exists. design.go, rank.go,
// and items.go select the design, tech, and ADRs and catalog the rest.
package context

import (
	"fmt"
	"path/filepath"
	"slices"
	"strings"
	"unicode/utf8"

	"github.com/bytepunx/system-flow/flai/internal/conventions"
	"github.com/bytepunx/system-flow/flai/internal/topics"
)

// AdditionsHeading is the heading kept in every convention whatever its
// topics, with the baseline marker, so the file's shape stays visible.
const AdditionsHeading = "Project additions"

// Piece is one section of a convention, kept or left out.
type Piece struct {
	Heading string   `json:"heading,omitempty"` // empty for the text before the first heading
	Path    []string `json:"path"`              // headings from the outermost down to this one
	Line    int      `json:"line"`              // 1-based line of the heading in the file
	Topics  []string `json:"topics"`            // the section's effective topics
	outer   bool     // no ancestor is left out as well
}

// Convention is one convention file as the story reads it.
type Convention struct {
	Path    string   `json:"path"` // relative to the repository root
	Name    string   `json:"name"`
	Title   string   `json:"title,omitempty"`
	Order   int      `json:"order"`
	Topics  []string `json:"topics"` // the file's; [all] when it declares none
	Kept    []Piece  `json:"kept"`
	LeftOut []Piece  `json:"left_out"`
	Text    string   `json:"text"` // the file with what is left out taken out
	title   string   // the level-1 heading, dropped from left-out labels
}

// Size is how much a pack prints.
type Size struct {
	Bytes int `json:"bytes"`
	Lines int `json:"lines"`
}

// Pack is the context pack for one story.
type Pack struct {
	Story       string              `json:"story"`
	Title       string              `json:"title"`
	Topics      []topics.StoryTopic `json:"topics"`
	Size        Size                `json:"size"` // of Body
	README      *Convention         `json:"readme,omitempty"`
	Conventions []Convention        `json:"conventions"`
	Issues      string              `json:"-"` // the open-issues table, empty when none is open
	OpenIssues  int                 `json:"open_issues"`
	Omitted     []string            `json:"left_out"` // one line per outermost section left out
	Items       []Item              `json:"items"`    // the design, tech, and ADRs chosen, in the order chosen
	Catalog     Catalog             `json:"catalog"`
}

// Catalog is every design, tech, and ADR document the pack does not print
// whole.
type Catalog struct {
	NotLoaded []Entry `json:"not_loaded"`
	InPart    []Entry `json:"in_part"`
}

// Build selects the conventions of set for a story with the given topics.
// root is the repository root, to name files relative to it; issues is the
// open-issues table printed after the conventions.
func Build(root, story, title string, storyTopics []topics.StoryTopic, set *conventions.Set, issues string) (*Pack, error) {
	names := topics.Names(storyTopics)
	p := &Pack{Story: story, Title: title, Topics: storyTopics, Issues: issues, Conventions: []Convention{}, Omitted: []string{}, Items: []Item{},
		Catalog: Catalog{NotLoaded: []Entry{}, InPart: []Entry{}}}
	if set.README != "" {
		rel, err := filepath.Rel(root, filepath.Join(set.Dir, "README.md"))
		if err != nil {
			return nil, err
		}
		c, err := Filter(filepath.ToSlash(rel), set.README, names)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", rel, err)
		}
		c.Order = 0
		p.README = c
	}
	for _, f := range set.Files {
		c, err := Filter(filepath.ToSlash(f.Path), f.Raw, names)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", f.Path, err)
		}
		c.Title, c.Order = f.Title, f.Order
		p.Conventions = append(p.Conventions, *c)
	}
	if issues != "" {
		p.OpenIssues = len(strings.Split(strings.TrimSpace(issues), "\n")) - 2
	}
	for _, c := range p.files() {
		p.Omitted = append(p.Omitted, c.Omitted()...)
	}
	p.size()
	return p, nil
}

func (p *Pack) size() {
	body := p.Body()
	p.Size = Size{Bytes: len(body), Lines: strings.Count(body, "\n")}
}

// Design runs the steps that choose design, tech, and ADRs for a story, in
// the order ADR-0047 gives: by topics, linked from the story, its epic, and
// its tasks, one step further from both, then ranked against query.
func Design(docs []*Doc, storyTopics []string, sources []Source, query string) *Selection {
	s := NewSelection(docs)
	s.ByTopics(storyTopics)
	s.Linked(sources)
	s.Step(StepTopics, StepLinked)
	s.Rank(query, Ranked, Ranked)
	return s
}

// AddDesign puts what a selection chose, and its catalog, into the pack.
func (p *Pack) AddDesign(s *Selection) {
	p.Items = s.Items()
	p.Catalog.NotLoaded, p.Catalog.InPart = s.Catalog()
	p.size()
}

// files is the README, when there is one, then the conventions in read order.
func (p *Pack) files() []Convention {
	var out []Convention
	if p.README != nil {
		out = append(out, *p.README)
	}
	return append(out, p.Conventions...)
}

// Matches says whether a section with these topics is for a story with
// these: the section's include all or any of the story's.
func Matches(section, story []string) bool {
	for _, t := range section {
		if t == topics.All || slices.Contains(story, t) {
			return true
		}
	}
	return false
}

// Filter reads one convention and keeps what matches the story's topics:
// the front matter, every section whose effective topics match, the heading
// line of a left-out section that encloses a kept one, the baseline marker,
// and the "## Project additions" heading. A file without topics is [all].
func Filter(path, content string, story []string) (*Convention, error) {
	doc, err := topics.Parse(content, []string{topics.All})
	if err != nil {
		return nil, err
	}
	c := &Convention{Path: path, Name: filepath.Base(path), Topics: doc.Topics, Kept: []Piece{}, LeftOut: []Piece{}}
	var out strings.Builder
	lines := strings.SplitAfter(content, "\n")
	if len(doc.Sections) == 0 {
		c.Text = content
		return c, nil
	}
	out.WriteString(strings.Join(lines[:doc.Sections[0].Line-1], ""))

	type open struct {
		level   int
		line    string
		dropped bool
		written bool
	}
	var stack []*open
	ancestors := func() { // write the heading lines of left-out sections above a kept line
		for _, o := range stack[:len(stack)-1] {
			if !o.written {
				out.WriteString(o.line + "\n")
				o.written = true
			}
		}
	}
	for _, s := range doc.Sections {
		if s.Level == 1 && c.title == "" {
			c.title = s.Heading
		}
		first, rest, _ := strings.Cut(s.Text, "\n")
		for len(stack) > 0 && s.Level > 0 && stack[len(stack)-1].level >= s.Level {
			stack = stack[:len(stack)-1]
		}
		under := false
		for _, o := range stack {
			under = under || o.dropped
		}
		cur := &open{level: s.Level, line: first}
		if s.Level > 0 {
			stack = append(stack, cur)
		}
		piece := Piece{Heading: s.Heading, Path: s.Path, Line: s.Line, Topics: s.Topics, outer: !under}
		if Matches(s.Topics, story) {
			if s.Level > 0 {
				ancestors()
			}
			cur.written = true
			out.WriteString(s.Text)
			c.Kept = append(c.Kept, piece)
			continue
		}
		cur.dropped = true
		c.LeftOut = append(c.LeftOut, piece)
		if s.Level == 2 && s.Heading == AdditionsHeading {
			ancestors()
			cur.written = true
			out.WriteString(first + "\n\n")
		}
		for _, l := range strings.SplitAfter(rest, "\n") {
			if strings.TrimSpace(l) == conventions.Marker {
				if len(stack) > 0 {
					ancestors()
				}
				out.WriteString(strings.TrimRight(l, "\n") + "\n\n")
			}
		}
	}
	c.Text = out.String()
	return c, nil
}

// Omitted is one line per outermost section left out, with its topics:
// "code-quality.md § Go (cli)". A section left out with its parent is
// covered by the parent's line.
func (c Convention) Omitted() []string {
	var out []string
	for _, p := range c.LeftOut {
		if p.outer {
			out = append(out, fmt.Sprintf("%s § %s (%s)", c.Name, c.label(p), strings.Join(p.Topics, ", ")))
		}
	}
	return out
}

// label is a section's heading path without the file's title.
func (c Convention) label(p Piece) string {
	path := p.Path
	if len(path) > 1 && path[0] == c.title {
		path = path[1:]
	}
	if len(path) == 0 {
		return "(before the first heading)"
	}
	return strings.Join(path, " › ")
}

// Body is the pack as flai prime --story prints it below its header: every
// convention under the header flai prime --cat gives it, the open-issues
// table, each design item headed with its path, heading path, and reason,
// the catalog, then what was left out. With nothing left out and no design
// documents it is flai prime --cat.
func (p *Pack) Body() string {
	var b strings.Builder
	for i, c := range p.files() {
		if i > 0 {
			b.WriteString("\n")
		}
		fmt.Fprintf(&b, "%s\n%s\n\n%s", c.Path, strings.Repeat("=", len(c.Path)), c.Text)
	}
	if p.Issues != "" {
		fmt.Fprintf(&b, "\nopen issues\n===========\n\n%s", p.Issues)
	}
	for _, it := range p.Items {
		head := it.Path
		if it.Label != "" {
			head += " § " + it.Label
		}
		reason := it.Reason
		if len(it.Also) > 0 {
			reason += "; also " + strings.Join(it.Also, "; ")
		}
		fmt.Fprintf(&b, "\n%s\n%s\nreason: %s\n\n%s", head, strings.Repeat("=", utf8.RuneCountInString(head)), reason, it.Text)
	}
	p.writeCatalog(&b)
	if len(p.Omitted) > 0 {
		b.WriteString("\nleft out\n========\n\n")
		for _, l := range p.Omitted {
			fmt.Fprintf(&b, "- %s\n", l)
		}
	}
	return b.String()
}

// writeCatalog lists the design documents not printed whole.
func (p *Pack) writeCatalog(b *strings.Builder) {
	c := p.Catalog
	if len(c.NotLoaded) == 0 && len(c.InPart) == 0 {
		return
	}
	b.WriteString("\ncatalog\n=======\n\nRead any of these with flai doc show <path> or the MCP doc_get.\n")
	if len(c.NotLoaded) > 0 {
		b.WriteString("\nNot loaded:\n\n")
		for _, e := range c.NotLoaded {
			fmt.Fprintf(b, "- %s: %s", e.Path, e.Title)
			if len(e.SupersededBy) > 0 {
				fmt.Fprintf(b, " (superseded by %s)", strings.Join(e.SupersededBy, ", "))
			}
			b.WriteString("\n")
		}
	}
	if len(c.InPart) > 0 {
		b.WriteString("\nLoaded in part, the sections above marked loaded:\n\n")
		for _, e := range c.InPart {
			fmt.Fprintf(b, "- %s: %s\n", e.Path, e.Title)
			for _, o := range e.Outline {
				mark := ""
				if o.Loaded {
					mark = " (loaded)"
				}
				fmt.Fprintf(b, "%s- %s%s\n", strings.Repeat("  ", o.Level), o.Heading, mark)
			}
		}
	}
}

// Header names the story, its topics and where each came from, and the size
// of the body that follows it.
func (p *Pack) Header() string {
	var b strings.Builder
	head := p.Story + " context pack"
	fmt.Fprintf(&b, "%s\n%s\n\n", head, strings.Repeat("=", len(head)))
	fmt.Fprintf(&b, "story: %s %s\n", p.Story, p.Title)
	b.WriteString("topics:\n")
	width := 0
	for _, t := range p.Topics {
		width = max(width, len(t.Topic))
	}
	for _, t := range p.Topics {
		fmt.Fprintf(&b, "  %-*s  %s\n", width, t.Topic, t.Summary())
	}
	if len(p.Items) > 0 {
		steps := map[string]int{}
		bytes := 0
		for _, it := range p.Items {
			steps[it.Step]++
			bytes += it.Size
		}
		fmt.Fprintf(&b, "design: %d items, %d bytes (by topics %d, linked %d, ranked %d)\n", len(p.Items), bytes, steps[StepTopics], steps[StepLinked], steps[StepRanked])
	}
	if n, m := len(p.Catalog.NotLoaded), len(p.Catalog.InPart); n+m > 0 {
		fmt.Fprintf(&b, "catalog: %d documents not loaded, %d loaded in part\n", n, m)
	}
	if len(p.Omitted) > 0 {
		fmt.Fprintf(&b, "left out: %d sections, listed at the end\n", len(p.Omitted))
	}
	fmt.Fprintf(&b, "size: %d bytes, %d lines below this header\n\n", p.Size.Bytes, p.Size.Lines)
	return b.String()
}
