package context

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"unicode"

	"github.com/goccy/go-yaml"

	"github.com/bytepunx/system-flow/flai/internal/topics"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// Kinds of design document, by the folder under design/ they live in.
const (
	KindSystem = "system"
	KindTech   = "tech"
	KindADR    = "adr"
)

// Steps of the pack that select design documents, in the order they run.
const (
	StepTopics = "topics"
	StepLinked = "linked"
	StepRanked = "ranked"
)

// folders maps each kind to its folder under design/, in the order the
// pack walks them.
var folders = []struct{ kind, dir string }{{KindSystem, "system"}, {KindTech, "tech"}, {KindADR, "adrs"}}

// Doc is one design, tech, or ADR file, cut at its headings.
type Doc struct {
	Path         string           `json:"path"` // relative to the repository root
	Kind         string           `json:"kind"`
	Title        string           `json:"title"`
	ID           string           `json:"id,omitempty"` // ADR-nnnn for an ADR
	Refines      []string         `json:"refines,omitempty"`
	SupersededBy []string         `json:"superseded_by,omitempty"`
	Topics       []string         `json:"topics"` // the file's; empty when it declares none
	Sections     []topics.Section `json:"-"`
	Raw          string           `json:"-"`
	title        string           // the level-1 heading, dropped from labels
}

// LoadDocs reads design/system, design/tech, and design/adrs under design,
// leaving out each folder's README and the ADR template. root names paths.
func LoadDocs(root, design string) ([]*Doc, error) {
	var docs []*Doc
	for _, f := range folders {
		dir := filepath.Join(design, f.dir)
		entries, err := os.ReadDir(dir)
		if os.IsNotExist(err) {
			continue
		}
		if err != nil {
			return nil, err
		}
		for _, e := range entries {
			name := e.Name()
			if e.IsDir() || !strings.HasSuffix(name, ".md") || name == "README.md" || (f.kind == KindADR && strings.HasPrefix(name, "0000-")) {
				continue
			}
			raw, err := os.ReadFile(filepath.Join(dir, name))
			if err != nil {
				return nil, err
			}
			rel, err := filepath.Rel(root, filepath.Join(dir, name))
			if err != nil {
				return nil, err
			}
			d, err := ParseDoc(filepath.ToSlash(rel), f.kind, string(raw))
			if err != nil {
				return nil, fmt.Errorf("%s: %w", rel, err)
			}
			docs = append(docs, d)
		}
	}
	return docs, nil
}

var adrFile = regexp.MustCompile(`^(\d{4})-`)

// ParseDoc reads one document: its front matter, its topics, and its
// sections. A file without topics is selected only by links and ranking.
func ParseDoc(rel, kind, raw string) (*Doc, error) {
	parsed, err := topics.Parse(raw, nil)
	if err != nil {
		return nil, err
	}
	d := &Doc{Path: rel, Kind: kind, Topics: parsed.Topics, Sections: parsed.Sections, Raw: raw}
	if fm, _, err := workitem.SplitFrontMatter(raw); err == nil {
		var f struct {
			ID           string   `yaml:"id"`
			Title        string   `yaml:"title"`
			Refines      []string `yaml:"refines"`
			SupersededBy []string `yaml:"superseded_by"`
		}
		if err := yaml.Unmarshal([]byte(fm), &f); err != nil {
			return nil, fmt.Errorf("the front matter does not parse: %w", err)
		}
		d.ID, d.Title, d.Refines, d.SupersededBy = f.ID, f.Title, adrIDs(f.Refines), adrIDs(f.SupersededBy)
	}
	for _, s := range d.Sections {
		if s.Level == 1 {
			d.title = s.Heading
			break
		}
	}
	if d.Title == "" {
		d.Title = d.title
	}
	if d.Title == "" {
		d.Title = path.Base(rel)
	}
	if kind == KindADR && d.ID == "" {
		if m := adrFile.FindStringSubmatch(path.Base(rel)); m != nil {
			d.ID = "ADR-" + m[1]
		}
	}
	return d, nil
}

var adrIDRe = regexp.MustCompile(`\bADR-(\d{1,4})\b`)

// adrIDs normalises ADR references to ADR-nnnn, dropping anything else.
func adrIDs(in []string) []string {
	var out []string
	for _, s := range in {
		if m := adrIDRe.FindStringSubmatch(s); m != nil {
			out = append(out, adrID(m[1]))
		}
	}
	return out
}

func adrID(digits string) string {
	n, _ := strconv.Atoi(digits)
	return fmt.Sprintf("ADR-%04d", n)
}

// Label is a section's heading path without the document's title,
// "Commands › flai prime"; empty for the text before the first heading.
func (d *Doc) Label(sec int) string {
	p := d.Sections[sec].Path
	if len(p) > 1 && p[0] == d.title {
		p = p[1:]
	}
	return strings.Join(p, " › ")
}

// under is the section and every section below it, down to the next
// heading at its level or higher.
func (d *Doc) under(sec int) []int {
	out := []int{sec}
	level := d.Sections[sec].Level
	if level == 0 {
		return out
	}
	for i := sec + 1; i < len(d.Sections) && d.Sections[i].Level > level; i++ {
		out = append(out, i)
	}
	return out
}

// all is every section of the document.
func (d *Doc) all() []int {
	out := make([]int, len(d.Sections))
	for i := range out {
		out[i] = i
	}
	return out
}

// choice is one reason the pack loads something, and the reasons found for
// the same sections after it.
type choice struct {
	step   string
	reason string
	also   []string
}

// Selection is what the pack loads of the design documents: each section
// is owned by the first reason that chose it, so nothing prints twice, and
// later reasons for it are listed on that one.
type Selection struct {
	Docs    []*Doc
	byPath  map[string]*Doc
	byID    map[string]*Doc
	owner   map[*Doc][]int // the choice owning each section, -1 when none
	choices []*choice
}

// NewSelection starts a selection over docs with nothing chosen.
func NewSelection(docs []*Doc) *Selection {
	s := &Selection{Docs: docs, byPath: map[string]*Doc{}, byID: map[string]*Doc{}, owner: map[*Doc][]int{}}
	for _, d := range docs {
		s.byPath[d.Path] = d
		if d.ID != "" {
			s.byID[d.ID] = d
		}
		own := make([]int, len(d.Sections))
		for i := range own {
			own[i] = -1
		}
		s.owner[d] = own
	}
	return s
}

// choose gives sections of d to a reason. A section already chosen keeps
// its reason and lists this one.
func (s *Selection) choose(d *Doc, secs []int, step, reason string) {
	idx := -1
	for _, i := range secs {
		switch o := s.owner[d][i]; {
		case o < 0:
			if idx < 0 {
				s.choices = append(s.choices, &choice{step: step, reason: reason})
				idx = len(s.choices) - 1
			}
			s.owner[d][i] = idx
		case o != idx && s.choices[o].reason != reason && !slices.Contains(s.choices[o].also, reason):
			s.choices[o].also = append(s.choices[o].also, reason)
		}
	}
}

// chooseADR chooses a whole ADR, or what supersedes it when it is
// superseded, each replacement with the reason "supersedes ADR-nnnn".
func (s *Selection) chooseADR(d *Doc, step, reason string) {
	for _, r := range s.Replace(d) {
		if r.by == nil {
			s.choose(r.doc, r.doc.all(), step, reason)
		} else {
			s.choose(r.doc, r.doc.all(), step, "supersedes "+r.by.ID)
		}
	}
}

// replacement is an ADR in force and the one it directly superseded, nil
// when it is the ADR asked about.
type replacement struct{ doc, by *Doc }

// Replace follows superseded_by from an ADR to the ADRs in force: the ADR
// itself when nothing present supersedes it.
func (s *Selection) Replace(d *Doc) []replacement {
	var out []replacement
	seen := map[*Doc]bool{d: true}
	var walk func(cur *Doc)
	walk = func(cur *Doc) {
		for _, id := range cur.SupersededBy {
			next := s.byID[id]
			if next == nil || seen[next] {
				continue
			}
			seen[next] = true
			if len(next.SupersededBy) == 0 || !s.anyPresent(next.SupersededBy) {
				out = append(out, replacement{doc: next, by: cur})
				continue
			}
			walk(next)
		}
	}
	walk(d)
	if len(out) == 0 {
		return []replacement{{doc: d}}
	}
	return out
}

func (s *Selection) anyPresent(ids []string) bool {
	for _, id := range ids {
		if s.byID[id] != nil {
			return true
		}
	}
	return false
}

// Superseded says whether an ADR is superseded by one present.
func (s *Selection) Superseded(d *Doc) bool {
	return d.Kind == KindADR && s.anyPresent(d.SupersededBy)
}

// ByTopics chooses every section whose topics match the story's, with the
// reason "topics: " and the topics that matched. An ADR matched in any
// section is chosen whole, or what supersedes it.
func (s *Selection) ByTopics(story []string) {
	for _, d := range s.Docs {
		byReason := map[string][]int{}
		var order []string
		for i, sec := range d.Sections {
			hit := matched(sec.Topics, story)
			if len(hit) == 0 {
				continue
			}
			r := "topics: " + strings.Join(hit, ", ")
			if _, ok := byReason[r]; !ok {
				order = append(order, r)
			}
			byReason[r] = append(byReason[r], i)
		}
		for _, r := range order {
			if d.Kind == KindADR && s.Superseded(d) {
				s.chooseADR(d, StepTopics, r)
				break
			}
			s.choose(d, byReason[r], StepTopics, r)
		}
	}
}

// matched is the section's topics that select it for the story: all, or
// those the story has.
func matched(section, story []string) []string {
	var out []string
	for _, t := range section {
		if t == topics.All || slices.Contains(story, t) {
			out = append(out, t)
		}
	}
	return out
}

// Source is a work item whose body can link documents: the story, its
// epic, or one of its tasks.
type Source struct {
	ID   string // S-nnnn, E-nnnn, or T-nnnn
	Path string // the item's file, relative to the repository root
	Body string
}

// ref is a document, or part of one, that a text links or names.
type ref struct {
	doc  *Doc
	secs []int // nil for the whole document
}

var (
	linkRe   = regexp.MustCompile(`\]\(\s*<?([^)\s>]+)>?(?:\s+"[^"]*")?\s*\)`)
	schemeRe = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9+.-]*:`)
)

// refs finds what a text at from (a repository path) links or names:
// markdown links to a document, with a fragment naming a heading narrowing
// it to that section; a document's repository path written out; and ADR IDs.
func (s *Selection) refs(text, from string) []ref {
	var out []ref
	for _, m := range linkRe.FindAllStringSubmatch(text, -1) {
		target := m[1]
		if schemeRe.MatchString(target) {
			continue
		}
		file, frag, _ := strings.Cut(target, "#")
		if file == "" {
			continue
		}
		d := s.resolve(path.Clean(path.Join(path.Dir(from), file)))
		if d == nil {
			continue
		}
		r := ref{doc: d}
		if frag != "" {
			for i, sec := range d.Sections {
				if sec.Level > 0 && Slug(sec.Heading) == strings.ToLower(frag) {
					r.secs = d.under(i)
					break
				}
			}
		}
		out = append(out, r)
	}
	plain := linkRe.ReplaceAllString(text, "]()")
	for _, d := range s.Docs {
		if strings.Contains(plain, d.Path) {
			out = append(out, ref{doc: d})
		}
	}
	for _, m := range adrIDRe.FindAllStringSubmatch(text, -1) {
		if d := s.byID[adrID(m[1])]; d != nil {
			out = append(out, ref{doc: d})
		}
	}
	return out
}

// resolve finds the document a cleaned link target names. A target that is
// not a document path as it stands is matched by its tail, so the links of
// an archived item, written from where it was before it moved, still resolve.
func (s *Selection) resolve(target string) *Doc {
	if d := s.byPath[target]; d != nil {
		return d
	}
	for _, d := range s.Docs {
		if strings.HasSuffix(target, "/"+d.Path) {
			return d
		}
	}
	return nil
}

// Slug is the anchor a heading gets: lower case, spaces to hyphens, and
// everything but letters, digits, hyphens, and underscores dropped.
func Slug(heading string) string {
	var b strings.Builder
	for _, r := range strings.ToLower(heading) {
		switch {
		case unicode.IsLetter(r) || unicode.IsDigit(r) || r == '-' || r == '_':
			b.WriteRune(r)
		case r == ' ':
			b.WriteRune('-')
		}
	}
	return b.String()
}

// Linked chooses what the story, its epic, and its tasks link or name, with
// the reason "linked from <ID>". A superseded ADR gives way to what
// supersedes it.
func (s *Selection) Linked(sources []Source) {
	for _, src := range sources {
		for _, r := range s.refs(src.Body, src.Path) {
			s.chooseRef(r, StepLinked, "linked from "+src.ID)
		}
	}
}

func (s *Selection) chooseRef(r ref, step, reason string) {
	switch {
	case r.doc.Kind == KindADR:
		s.chooseADR(r.doc, step, reason)
	case r.secs != nil:
		s.choose(r.doc, r.secs, step, reason)
	default:
		s.choose(r.doc, r.doc.all(), step, reason)
	}
}

// Step goes one step further from what the given steps chose: the ADRs a
// chosen section links or names ("linked from design/system/x.md §
// Heading") and the ADRs a chosen ADR refines ("refined by ADR-nnnn"). What
// the step adds is not followed again.
func (s *Selection) Step(from ...string) {
	type at struct {
		doc  *Doc
		secs []int
	}
	var chosen []at
	for _, d := range s.Docs {
		var secs []int
		for i, o := range s.owner[d] {
			if o >= 0 && slices.Contains(from, s.choices[o].step) {
				secs = append(secs, i)
			}
		}
		if len(secs) > 0 {
			chosen = append(chosen, at{d, secs})
		}
	}
	for _, c := range chosen {
		if c.doc.Kind == KindADR {
			for _, id := range c.doc.Refines {
				if r := s.byID[id]; r != nil && r != c.doc {
					s.chooseADR(r, StepLinked, "refined by "+c.doc.ID)
				}
			}
		}
		for _, i := range c.secs {
			for _, r := range s.refs(c.doc.Sections[i].Text, c.doc.Path) {
				if r.doc.Kind != KindADR || r.doc == c.doc {
					continue
				}
				from := c.doc.Path
				if l := c.doc.Label(i); l != "" {
					from += " § " + l
				}
				s.chooseADR(r.doc, StepLinked, "linked from "+from)
			}
		}
	}
}

// Loaded says whether any section of d is chosen.
func (s *Selection) Loaded(d *Doc) bool {
	return slices.ContainsFunc(s.owner[d], func(o int) bool { return o >= 0 })
}

// Whole says whether every section of d is chosen.
func (s *Selection) Whole(d *Doc) bool {
	return !slices.Contains(s.owner[d], -1)
}
