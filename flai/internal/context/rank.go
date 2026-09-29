package context

import (
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/bytepunx/system-flow/flai/internal/search"
	"github.com/bytepunx/system-flow/flai/internal/topics"
)

// Query is what a story is ranked against: its title, goal, and acceptance
// criteria, or its whole body when it has neither section.
func Query(title, body string) string {
	doc, err := topics.Parse(body, nil)
	if err != nil {
		return title + "\n" + body
	}
	var b strings.Builder
	b.WriteString(title + "\n")
	found := false
	for _, s := range doc.Sections {
		if s.Level == 2 && (s.Heading == "Goal" || s.Heading == "Acceptance criteria") {
			b.WriteString(s.Text)
			found = true
		}
	}
	if !found {
		b.WriteString(body)
	}
	return b.String()
}

// Candidate is what the ranked step may load: a design section, cut at its
// own heading, or an ADR by its best section.
type Candidate struct {
	doc *Doc
	sec int
}

// Candidates ranks by BM25 against query every section nothing has loaded,
// best first. Sections are indexed one by one, each down to the next heading
// of any level; an ADR appears once, at its best section. ADRs loaded in
// any part, superseded ADRs, and sections with nothing below their heading
// are left out. A briefed document's sections are in: a brief is not its
// body.
func (s *Selection) Candidates(query string) []Candidate {
	var (
		keys []Candidate
		docs []search.Doc
	)
	for _, d := range s.Docs {
		if d.Kind == KindADR && (s.Loaded(d) || s.Superseded(d)) {
			continue
		}
		for i, sec := range d.Sections {
			if s.owner[d][i] >= 0 {
				continue
			}
			body := sec.Text
			if sec.Level > 0 {
				_, body, _ = strings.Cut(body, "\n")
			}
			if strings.TrimSpace(body) == "" {
				continue
			}
			docs = append(docs, search.Doc{Path: strconv.Itoa(len(keys)), Kind: "doc", Title: d.Title, Headings: d.Label(i), Body: body, Scope: "design"})
			keys = append(keys, Candidate{d, i})
		}
	}
	ix := search.Build(docs)
	var out []Candidate
	seen := map[*Doc]bool{}
	for _, h := range ix.Search(strings.Join(Terms(query), " "), true, ix.Size()) {
		n, _ := strconv.Atoi(h.Path)
		c := keys[n]
		if c.doc.Kind == KindADR {
			if seen[c.doc] {
				continue
			}
			seen[c.doc] = true
		}
		out = append(out, c)
	}
	return out
}

// Rank loads a candidate with the reason "rank n": an ADR whole, or its
// best section when whole is false; a design section alone. It returns the
// undo that puts the selection back as it was.
func (s *Selection) Rank(c Candidate, n int, whole bool) (undo func()) {
	own, count := slices.Clone(s.owner[c.doc]), len(s.choices)
	secs := []int{c.sec}
	if whole && c.doc.Kind == KindADR {
		secs = c.doc.all()
	}
	s.choose(c.doc, secs, StepRanked, fmt.Sprintf("rank %d", n))
	return func() {
		s.owner[c.doc] = own
		s.choices = s.choices[:count]
	}
}

// RankSize is how much a candidate adds to the pack at least: the text it
// loads.
func (c Candidate) RankSize(whole bool) int {
	if whole && c.doc.Kind == KindADR {
		return len(c.doc.Raw)
	}
	return len(c.doc.Sections[c.sec].Text)
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
	for _, t := range search.Tokens(query) {
		if len([]rune(t)) > 1 && !stopwords[t] && !seen[t] {
			seen[t] = true
			out = append(out, t)
		}
	}
	return out
}
