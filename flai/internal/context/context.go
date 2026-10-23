// Package context builds the context pack an agent working a story is primed
// with (ADR-0047), fitted to a size budget (ADR-0049). This part selects the
// conventions: each file in read order with the sections whose topics miss
// the story's left out, and a list of what was left out so the agent knows
// the rule exists; and it fits the pack to its budget. design.go, brief.go,
// rank.go, and items.go select, brief, and rank the design, tech, and ADRs
// and catalog the rest.
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
	Size    int      `json:"size"` // bytes of Text
	// Chunk and Chunks number a convention too large for any part of a pack,
	// which Parts splits into chunks of its text (ADR-0104); zero otherwise.
	Chunk  int    `json:"chunk,omitempty"`
	Chunks int    `json:"chunks,omitempty"`
	title  string // the level-1 heading, dropped from left-out labels
}

// Size is how much a pack prints, its header included.
type Size struct {
	Bytes int `json:"bytes"`
	Lines int `json:"lines"`
}

// What a pack's Exceeded names: the part that alone does not fit the
// budget.
const (
	ExceededConventions = "conventions"
	ExceededNamed       = "named"
	ExceededBriefs      = "briefs"
)

// Pack is the context pack for one story, or for a strategic agent.
type Pack struct {
	Story       string              `json:"story,omitempty"` // the story worked, empty for a strategic agent's pack
	Item        string              `json:"item,omitempty"`  // the epic or story the planner plans
	Title       string              `json:"title"`           // the story's or the item's
	Role        string              `json:"role,omitempty"`  // a sub-agent's role (ADR-0059) or a strategic agent's (E-0016); empty for the story's agent
	PartNumber  int                 `json:"part,omitempty"`  // which part of the pack this is, on a part Parts cut (ADR-0104); zero on the whole pack
	PartCount   int                 `json:"parts,omitempty"` // how many parts Parts cut the pack into; zero on the whole pack
	Goal        string              `json:"goal,omitempty"`  // a role pack's story goal and acceptance criteria`
	Topics      []topics.StoryTopic `json:"topics"`
	Budget      int                 `json:"budget"`             // bytes the pack is fitted to
	Exceeded    string              `json:"exceeded,omitempty"` // conventions, named, or briefs: the part that takes the pack over the budget
	Size        Size                `json:"size"`               // of Header and Body together
	README      *Convention         `json:"readme,omitempty"`
	Conventions []Convention        `json:"conventions"`
	Issues      string              `json:"-"` // the open-issues table, empty when none is open
	OpenIssues  int                 `json:"open_issues"`
	Omitted     []string            `json:"left_out"` // one line per outermost section left out
	Items       []Item              `json:"items"`    // the design, tech, ADRs, and digest entries chosen, in the order chosen
	Catalog     Catalog             `json:"catalog"`
	// BriefsLeftOut counts the briefs a role pack had no room for.
	BriefsLeftOut int `json:"briefs_left_out,omitempty"`
}

// Catalog is every design, tech, and ADR document the pack does not print
// whole.
type Catalog struct {
	NotLoaded []Entry `json:"not_loaded"`
	InPart    []Entry `json:"in_part"`
}

// Build selects the conventions of set for a story with the given topics,
// for a pack of budget bytes. root is the repository root, to name files
// relative to it; issues is the open-issues table printed after the
// conventions.
func Build(root, story, title string, storyTopics []topics.StoryTopic, set *conventions.Set, issues string, budget int) (*Pack, error) {
	names := topics.Names(storyTopics)
	p := &Pack{Story: story, Title: title, Topics: storyTopics, Budget: budget, Issues: issues, Conventions: []Convention{}, Omitted: []string{}, Items: []Item{},
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
		c.Order, c.Size = 0, len(c.Text)
		p.README = c
	}
	for _, f := range set.Files {
		c, err := Filter(filepath.ToSlash(f.Path), f.Raw, names)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", f.Path, err)
		}
		c.Title, c.Order, c.Size = f.Title, f.Order, len(c.Text)
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

// size measures the pack as printed, header and body. The header states
// the size, so it is measured until the number it prints is its own.
func (p *Pack) size() {
	body := p.Body()
	lines := strings.Count(body, "\n")
	for range 4 {
		h := p.Header()
		next := Size{Bytes: len(h) + len(body), Lines: strings.Count(h, "\n") + lines}
		if next == p.Size {
			return
		}
		p.Size = next
	}
}

// Design runs the steps that choose design, tech, and ADRs for a story, in
// the order ADR-0049 gives: what the story, its epic, and its tasks name, to
// load whole, save a large document named only by its path written out,
// which is briefed when it is over briefOver bytes (ADR-0050); the design
// system ds, when its guidance's topics reach the story, to load whole
// likewise (ADR-0152), nil for none; what its topics select, to brief; and
// the ADRs one step from them all, to brief. The ranked step runs in
// AddDesign, against the budget.
func Design(docs []*Doc, storyTopics []string, sources []Source, briefOver int, ds *DesignSystem) *Selection {
	s := NewSelection(docs)
	s.BriefOver = briefOver
	s.Linked(sources)
	s.DesignSystem(ds, storyTopics)
	s.ByTopics(storyTopics)
	s.Step()
	return s
}

// AddDesign puts what a selection chose into the pack, fills what the
// budget leaves with the sections that rank highest against query, and adds
// the catalog. When the conventions alone exceed the budget the pack is
// the conventions and a catalog of every document; when the conventions and
// what the story names exceed it, everything named is still loaded; when
// the briefs take it over, every brief is still printed (TH-0032). In the
// last two nothing is ranked. Which part took it over is in Exceeded.
func (p *Pack) AddDesign(s *Selection, query string) { p.addDesign(s, query, nil) }

// addDesign is AddDesign with the digest step (ADR-0194) after the briefs
// and before the ranked sections: the digest items, best first, as
// digestFor gives them, loaded as addDigest loads them. Like the ranked
// step, it loads nothing when Exceeded is set.
func (p *Pack) addDesign(s *Selection, query string, digest []Item) {
	p.size()
	if p.Size.Bytes > p.Budget {
		p.Exceeded = ExceededConventions
		p.Catalog.NotLoaded, p.Catalog.InPart = NewSelection(s.Docs).Catalog()
		p.size()
		return
	}
	p.Items = s.loaded(StepNamed)
	p.size()
	if p.Size.Bytes > p.Budget {
		p.Exceeded = ExceededNamed
	}
	p.take(s, nil)
	if p.Exceeded == "" && p.Size.Bytes > p.Budget {
		p.Exceeded = ExceededBriefs
	}
	if p.Exceeded != "" {
		return
	}
	digested := p.addDigest(digest)
	n := 0
	for _, c := range s.Candidates(query) {
		wholes := []bool{false}
		if c.doc.Kind == KindADR {
			wholes = []bool{true, false}
		}
		for _, whole := range wholes {
			if p.Size.Bytes+c.RankSize(whole) > p.Budget {
				continue
			}
			items, catalog, size := p.Items, p.Catalog, p.Size
			undo := s.Rank(c, n+1, whole)
			p.take(s, digested)
			if p.Size.Bytes <= p.Budget {
				n++
				break
			}
			undo()
			p.Items, p.Catalog, p.Size = items, catalog, size
		}
	}
}

// AddBriefs puts a role pack's briefs into it (ADR-0059): each brief the
// selection made, in the order Items gives them, then the digest briefs
// digest gives, best first (ADR-0194), each that fits what the budget
// leaves, and a count of those that do not. Nothing is ranked and there is
// no catalog: a sub-agent finds the rest with doc_search. When the
// conventions and the story's goal alone exceed the budget, no brief is
// added and Exceeded says conventions.
func (p *Pack) AddBriefs(s *Selection, digest []Item) {
	p.size()
	over := p.Size.Bytes > p.Budget
	if over {
		p.Exceeded = ExceededConventions
	}
	var briefs []Item
	for _, it := range s.Items() {
		if it.Step == StepBriefed {
			briefs = append(briefs, it)
		}
	}
	for _, it := range append(briefs, digest...) {
		if over {
			p.BriefsLeftOut++
			continue
		}
		p.Items = append(p.Items, it)
		p.size()
		if p.Size.Bytes > p.Budget {
			p.Items = p.Items[:len(p.Items)-1]
			p.BriefsLeftOut++
		}
	}
	p.size()
	// The header's line counting what was left out can itself take the pack
	// over; then the last brief makes room for it.
	for !over && p.Size.Bytes > p.Budget && len(p.Items) > 0 {
		p.Items = p.Items[:len(p.Items)-1]
		p.BriefsLeftOut++
		p.size()
	}
}

// take puts everything a selection chose, and its catalog, into the pack,
// with the digest items digested before what is ranked.
func (p *Pack) take(s *Selection, digested []Item) {
	items := s.Items()
	at := slices.IndexFunc(items, func(it Item) bool { return it.Step == StepRanked })
	if at < 0 {
		at = len(items)
	}
	p.Items = slices.Insert(items, at, digested...)
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
// table, each named item headed with its path, heading path, and reason,
// the brief of each design and tech file, the decisions of the ADRs
// briefed, the digest entries, each ranked item, the catalog, then what was
// left out. With nothing left out and no design documents it is flai prime
// --cat.
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
	if p.Goal != "" {
		fmt.Fprintf(&b, "\nstory\n=====\n\n%s", p.Goal)
	}
	var last string
	for _, it := range p.Items {
		if h := p.groupHead()[group(it)]; group(it) != last && h != "" {
			b.WriteString(h)
		}
		last = group(it)
		b.WriteString(p.printed(it))
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

// Groups of items that print under one heading: design briefs, ADR briefs,
// and digest entries. Other loaded items head themselves.
const (
	groupBriefs    = "briefs"
	groupDecisions = "decisions"
	groupDigest    = "digest"
)

var groupHead = map[string]string{
	groupBriefs:    "\nbriefs\n======\n\nThe design and tech files the story, its epic, or its tasks name by a path written out and that are too large to load whole (reason: named in <ID>), the design system's guidance and token file when they are too large to load whole (reason: design system: topics <topics>), then those the story's topics select, each as its title, size, reason, first paragraph, and outline. A brief is not the document: when one bears on the story, read the section with the MCP doc_get and its heading, or flai doc show <path> --heading \"<heading>\", or the whole file, before relying on it or changing what it describes. doc_search, or flai doc search, finds sections by their words.\n",
	groupDecisions: "\ndecisions\n=========\n\nThe ADRs the pack reached, and those named by a path written out and too large to load whole, each by its decision sentence. When one bears on the story, read it with the MCP doc_get, or its decision alone with heading Decision (flai doc show <path> --heading Decision), before relying on it.\n\n",
	groupDigest:    "\ndigest\n======\n\nThe digest entries whose paths meet the story's claim (reason: path: <claimed path>) or whose topics the story shares (reason: topics: <topics>), ranked against the story's title, goal, and criteria, within an eighth of the budget (ADR-0194): each with its body, or as one line when it is compacted. They are what past stories decided and the operator answered, not the design. When you rely on one, cite its D-nnnn in your narrative or thread; check it against the design before changing what it describes; and ask the operator only what no entry answers. Read one whole with the MCP doc_get and its path, or flai doc show <path>.\n",
}

// roleHead is groupHead for a sub-agent's pack, which briefs everything.
var roleHead = map[string]string{
	groupBriefs:    "\nbriefs\n======\n\nThe design and tech files the story, its epic, or its tasks name, the design system's guidance and token file when its topics reach the story, then those the story's topics select, each as its title, size, reason, first paragraph, and outline. A sub-agent's pack holds no bodies: when one bears on your question, read the section that does with the MCP doc_get and its heading, or flai doc show <path> --heading \"<heading>\". doc_search, or flai doc search, finds sections by their words.\n",
	groupDecisions: "\ndecisions\n=========\n\nThe ADRs the story names and those the pack reached, each by its decision sentence. When one bears on your question, read it with the MCP doc_get.\n\n",
	groupDigest:    "\ndigest\n======\n\nThe digest entries the story's claim or topics select, ranked against the story (ADR-0194), each as one line: its ID, title, path, size, and reason. They are what past stories decided and the operator answered, not the design. When one bears on your question, read it with the MCP doc_get and its path, and name its D-nnnn in your report.\n",
}

// strategicHead is groupHead for a strategic agent's pack, which has no
// story.
var strategicHead = map[string]string{
	groupBriefs:    "\nbriefs\n======\n\nThe design and tech files the item planned or its stories name by a path written out and that are too large to load whole (reason: named in <ID>), then those the pack's topics select, each as its title, size, reason, first paragraph, and outline. A brief is not the document: when one bears on the work, read the section with the MCP doc_get and its heading, or flai doc show <path> --heading \"<heading>\", or the whole file, before relying on it or changing what it describes. doc_search, or flai doc search, finds sections by their words.\n",
	groupDecisions: "\ndecisions\n=========\n\nThe ADRs the pack reached, and those named by a path written out and too large to load whole, each by its decision sentence. When one bears on the work, read it with the MCP doc_get, or its decision alone with heading Decision (flai doc show <path> --heading Decision), before relying on it.\n\n",
	groupDigest:    "\ndigest\n======\n\nThe answers and recommendations of the digest that rank highest against each thread awaiting the operator (reason: thread TH-nnnn; rank n), within an eighth of the budget (ADR-0194): each with its body, or as one line when it is compacted. They are what the operator answered and confirmed before, not the design. When one answers a thread, cite its D-nnnn and give its path as a source on thread_reply; ask the operator only what no entry answers. Read one whole with the MCP doc_get and its path, or flai doc show <path>.\n",
}

// groupHead is the heading each group prints under in this pack.
func (p *Pack) groupHead() map[string]string {
	switch {
	case p.subAgent():
		return roleHead
	case p.strategic():
		return strategicHead
	}
	return groupHead
}

// printed is an item as this pack prints it: as Printed, save that a role
// pack does not say a document it names is briefed for its size, because
// it briefs everything.
func (p *Pack) printed(it Item) string {
	if p.subAgent() {
		return it.print(false)
	}
	return it.Printed()
}

// namedBrief follows the line naming a brief of a document the story named
// by its path written out (ADR-0050).
const namedBrief = "The story names this file; it is briefed for its size. Decide from the brief whether the work needs its body, and read the file, or the sections you will change, before relying on it or changing it.\n"

// group is the heading an item prints under, empty for one loaded that is
// not a digest entry.
func group(it Item) string {
	switch {
	case it.Step == StepDigest:
		return groupDigest
	case it.Step != StepBriefed:
		return ""
	case it.ID != "":
		return groupDecisions
	}
	return groupBriefs
}

// Printed is the item as the pack prints it: an ADR brief as one line, a
// design brief as a line naming it followed by its text, a digest entry
// that loads no body as one line, all under their group's heading; what is
// loaded headed with its path, heading path, and reason, a digest entry
// under the digest heading. A brief names its first reason, without the
// heading it was found under, and counts the rest; what is loaded names
// them all.
func (it Item) Printed() string { return it.print(true) }

// print is Printed, with the line that says a named document is briefed for
// its size when sized is true.
func (it Item) print(sized bool) string {
	switch group(it) {
	case groupDigest:
		if it.Text == "" {
			return it.digestBrief()
		}
	case groupDecisions:
		return fmt.Sprintf("- %s %s (%s; %s): %s\n", it.ID, it.Title, it.Path, it.briefReason(), it.Text)
	case groupBriefs:
		head := fmt.Sprintf("\n%s: %s (%d bytes; %s)\n", it.Path, it.Title, it.Whole, it.briefReason())
		switch {
		case !sized:
		case strings.HasPrefix(it.Reason, NamedIn):
			head += namedBrief
		case strings.HasPrefix(it.Reason, DesignSystemReason):
			head += designSystemBrief
		}
		if it.Text == "" {
			return head
		}
		return head + "\n" + it.Text
	}
	head := it.Path
	if it.Label != "" {
		head += " § " + it.Label
	}
	reason := it.Reason
	if len(it.Also) > 0 {
		reason += "; also " + strings.Join(it.Also, "; ")
	}
	return fmt.Sprintf("\n%s\n%s\nreason: %s\n\n%s", head, strings.Repeat("=", utf8.RuneCountInString(head)), reason, it.Text)
}

// briefReason is a brief's first reason without the heading it was found
// under, and how many more there are.
func (it Item) briefReason() string {
	r, _, _ := strings.Cut(it.Reason, " § ")
	switch n := len(it.Also); n {
	case 0:
		return r
	case 1:
		return r + ", and one more"
	default:
		return fmt.Sprintf("%s, and %d more", r, n)
	}
}

// writeCatalog lists the design documents not printed whole.
func (p *Pack) writeCatalog(b *strings.Builder) {
	c := p.Catalog
	if len(c.NotLoaded) == 0 && len(c.InPart) == 0 {
		return
	}
	b.WriteString("\ncatalog\n=======\n\nRead any of these, or one section with its heading, with the MCP doc_get or flai doc show <path> --heading \"<heading>\"; find sections by their words with doc_search or flai doc search.\n")
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

// Header names the story, or the item planned, its topics and where each
// came from, the size of the pack against its budget, whether a part alone
// exceeds it, and the size of each thing the pack prints.
func (p *Pack) Header() string {
	var b strings.Builder
	subject := p.Story
	switch {
	case p.Item != "":
		subject = p.Item
	case p.strategic():
		subject = "Project"
	}
	head := subject + " context pack"
	if p.Role != "" {
		head += " for the " + p.Role + " role"
	}
	fmt.Fprintf(&b, "%s\n%s\n\n", head, strings.Repeat("=", len(head)))
	if p.strategic() {
		b.WriteString(p.strategicLines())
	} else {
		fmt.Fprintf(&b, "story: %s %s\n", p.Story, p.Title)
	}
	if p.subAgent() {
		fmt.Fprintf(&b, "role: %s, a sub-agent of the story's agent (ADR-0059): the conventions the role reads, the story's goal and criteria, and briefs while the budget has room; doc_search and doc_get find the rest. Report back to the agent that started you; never move, edit, or create an item, and never write to a thread.\n", p.Role)
	}
	b.WriteString("topics:\n")
	width := 0
	for _, t := range p.Topics {
		width = max(width, len(t.Topic))
	}
	for _, t := range p.Topics {
		fmt.Fprintf(&b, "  %-*s  %s\n", width, t.Topic, t.Summary())
	}
	fmt.Fprintf(&b, "size: %d bytes, %d lines, this header included; budget %d bytes\n", p.Size.Bytes, p.Size.Lines, p.Budget)
	switch {
	case p.subAgent() && p.Exceeded == ExceededConventions:
		b.WriteString("over budget: the role's conventions and the story's goal alone exceed it, so no brief is printed; narrowing the conventions' roles or topics makes room (ADR-0059)\n")
	case p.Exceeded == ExceededConventions:
		b.WriteString("over budget: the conventions alone exceed it, so the pack is the conventions and a catalog; narrowing the conventions' topics makes room (ADR-0049)\n")
	case p.Exceeded == ExceededNamed && p.Item != "":
		fmt.Fprintf(&b, "over budget: the conventions and what %s names exceed it, so nothing is ranked (ADR-0049)\n", p.Item)
	case p.Exceeded == ExceededNamed:
		b.WriteString("over budget: the conventions and what the story, its epic, and its tasks name exceed it, so nothing is ranked; a story that names more than the budget holds is a story to split (ADR-0049)\n")
	case p.Exceeded == ExceededBriefs:
		b.WriteString("over budget: the conventions, what is named, and the briefs exceed it, so nothing is ranked; every brief is kept, and narrowing the conventions' topics makes room (ADR-0049, TH-0032)\n")
	}
	b.WriteString("contents, in bytes:\n")
	row := func(n int, what, name string) { fmt.Fprintf(&b, "  %6d  %-10s  %s\n", n, what, name) }
	for _, c := range p.files() {
		row(c.Size, "convention", c.Path)
	}
	if p.Issues != "" {
		row(len(p.Issues), "issues", fmt.Sprintf("%d open", p.OpenIssues))
	}
	if p.Goal != "" {
		row(len(p.Goal), "story", "goal and acceptance criteria")
	}
	decisions, adrs := len(p.groupHead()[groupDecisions]), 0
	digested, entries, whole := len(p.groupHead()[groupDigest]), 0, 0
	for _, it := range p.Items {
		switch group(it) {
		case groupDecisions:
			decisions += len(p.printed(it))
			adrs++
		case groupDigest:
			digested += len(p.printed(it))
			entries++
			if it.Text != "" {
				whole++
			}
		}
	}
	for _, it := range p.Items {
		switch {
		case group(it) == groupDigest:
			if entries > 0 {
				row(digested, "digest", digestRow(entries, whole))
				entries = 0
			}
		case group(it) != groupDecisions:
			name := it.Path
			if it.Label != "" {
				name += " § " + it.Label
			}
			row(len(p.printed(it)), it.Step, name)
		case adrs == 1:
			row(decisions, "briefed", "1 ADR by its decision sentence, on one line")
		case adrs > 1:
			row(decisions, "briefed", fmt.Sprintf("%d ADRs by their decision sentences, one line each", adrs))
			adrs = 0
		}
	}
	if n, m := len(p.Catalog.NotLoaded), len(p.Catalog.InPart); n+m > 0 {
		var c strings.Builder
		p.writeCatalog(&c)
		row(c.Len(), "catalog", fmt.Sprintf("%d documents not loaded, %d loaded in part", n, m))
	}
	if p.BriefsLeftOut > 0 {
		fmt.Fprintf(&b, "  %6s  %-10s  %d briefs with no room in the budget, not printed\n", "", "left out", p.BriefsLeftOut)
	}
	if len(p.Omitted) > 0 {
		n := 0
		for _, l := range p.Omitted {
			n += len(l) + 3
		}
		row(n, "left out", fmt.Sprintf("%d convention sections, listed at the end", len(p.Omitted)))
	}
	b.WriteString("\n")
	return b.String()
}
