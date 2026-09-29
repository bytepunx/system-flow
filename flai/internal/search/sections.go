package search

import (
	"strconv"
	"strings"
)

// Section is one section of a markdown document as the section index holds
// it: its heading line and what follows, down to the next heading of any
// level (ADR-0049). The context pack's ranked step and doc_search cut
// documents the same way, so what one finds the other can load.
type Section struct {
	Path    string // the document's repository path
	Title   string // the document's title
	Heading string // the heading path below the title, "Commands › flai prime"; empty before the first heading
	Line    int    // 1-based line of the heading in the file
	Level   int    // 1 to 6; 0 for the text before the first heading
	Text    string // the heading line and what follows it
}

// Body is the section's text without its heading line.
func (s Section) Body() string {
	if s.Level == 0 {
		return s.Text
	}
	_, body, _ := strings.Cut(s.Text, "\n")
	return body
}

// MaxSectionHits is the most results a section search returns.
const MaxSectionHits = 20

// SectionHit is one section a search found, with what an agent needs to
// decide whether to fetch it.
type SectionHit struct {
	Path    string  `json:"path"`
	Title   string  `json:"title"`
	Heading string  `json:"heading" jsonschema:"the heading path below the document's title; empty for the text before the first heading"`
	Line    int     `json:"line"`
	Lines   string  `json:"lines" jsonschema:"the first lines of the section below its heading"`
	Size    int     `json:"size" jsonschema:"bytes in the section, heading included"`
	Score   float64 `json:"score"`
}

// Sections is a BM25 index over sections, each indexed on its own.
type Sections struct {
	secs []Section
	at   []int // the section each indexed document is
	ix   *Index
}

// IndexSections indexes the sections that have anything below their
// heading.
func IndexSections(secs []Section) *Sections {
	x := &Sections{secs: secs}
	var docs []Doc
	for i, s := range secs {
		body := s.Body()
		if strings.TrimSpace(body) == "" {
			continue
		}
		docs = append(docs, Doc{Path: strconv.Itoa(len(x.at)), Kind: "doc", Title: s.Title, Headings: s.Heading, Body: body, Scope: "design"})
		x.at = append(x.at, i)
	}
	x.ix = Build(docs)
	return x
}

type ranked struct {
	sec   int
	score float64
}

func (x *Sections) rank(query string, limit int) []ranked {
	if limit <= 0 {
		limit = x.ix.Size()
	}
	hits := x.ix.Search(strings.Join(Terms(query), " "), true, limit)
	out := make([]ranked, len(hits))
	for i, h := range hits {
		n, _ := strconv.Atoi(h.Path)
		out[i] = ranked{x.at[n], h.Score}
	}
	return out
}

// Rank is every section that matches the query's terms, best first, as
// indexes into the sections given to IndexSections.
func (x *Sections) Rank(query string) []int {
	r := x.rank(query, 0)
	out := make([]int, len(r))
	for i, h := range r {
		out[i] = h.sec
	}
	return out
}

// Search is the sections that rank highest against the query, at most limit
// and never more than MaxSectionHits.
func (x *Sections) Search(query string, limit int) []SectionHit {
	if limit <= 0 || limit > MaxSectionHits {
		limit = MaxSectionHits
	}
	out := []SectionHit{}
	for _, h := range x.rank(query, limit) {
		s := x.secs[h.sec]
		out = append(out, SectionHit{Path: s.Path, Title: s.Title, Heading: s.Heading, Line: s.Line, Lines: FirstLines(s.Body(), 3, 160), Size: len(s.Text), Score: h.score})
	}
	return out
}

// FirstLines is the first n lines of text that are not blank, each cut at
// width characters.
func FirstLines(text string, n, width int) string {
	var out []string
	for _, line := range strings.Split(text, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if r := []rune(line); len(r) > width {
			line = string(r[:width]) + "…"
		}
		out = append(out, line)
		if len(out) == n {
			break
		}
	}
	return strings.Join(out, "\n")
}

// stopwords are the words too common to rank by: with terms combined by OR,
// "the" and "a" would match every section.
var stopwords = func() map[string]bool {
	m := map[string]bool{}
	for _, w := range strings.Fields(`a about above after again all also an and any are as at be because been
		before being below between both but by can could did do does doing down each few for from further had has
		have having here how if in into is it its itself just more most no nor not now of off on once only or other
		our out over own same should so some such than that the their them then there these they this those through
		to too under until up very was we were what when where which while who whom why will with would you your`) {
		m[w] = true
	}
	return m
}()

// Terms is a query's words without stopwords and single characters, each
// once, in the order they first appear.
func Terms(query string) []string {
	var out []string
	seen := map[string]bool{}
	for _, t := range Tokens(query) {
		if len([]rune(t)) > 1 && !stopwords[t] && !seen[t] {
			seen[t] = true
			out = append(out, t)
		}
	}
	return out
}
