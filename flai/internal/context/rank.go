package context

import (
	"fmt"
	"slices"
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

// Goal is a story's "## Goal" and "## Acceptance criteria" sections as they
// stand, the headings under them included, or its whole body when it has
// neither.
func Goal(body string) string {
	doc, err := topics.Parse(body, nil)
	if err != nil {
		return body
	}
	var b strings.Builder
	in := false
	for _, s := range doc.Sections {
		if s.Level > 0 && s.Level <= 2 {
			in = s.Level == 2 && (s.Heading == "Goal" || s.Heading == "Acceptance criteria")
		}
		if in {
			b.WriteString(s.Text)
		}
	}
	if b.Len() == 0 {
		return body
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
// best first, with the section index doc_search uses (ADR-0049). Sections
// are indexed one by one, each down to the next heading of any level; an ADR
// appears once, at its best section. ADRs loaded in any part, superseded
// ADRs, and sections with nothing below their heading are left out. A
// briefed document's sections are in: a brief is not its body.
func (s *Selection) Candidates(query string) []Candidate {
	var (
		keys []Candidate
		secs []search.Section
	)
	for _, d := range s.Docs {
		if d.Kind == KindADR && (s.Loaded(d) || s.Superseded(d)) {
			continue
		}
		for i, sec := range d.Cuts() {
			if s.owner[d][i] >= 0 {
				continue
			}
			secs = append(secs, sec)
			keys = append(keys, Candidate{d, i})
		}
	}
	var out []Candidate
	seen := map[*Doc]bool{}
	for _, n := range search.IndexSections(secs).Rank(query) {
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
