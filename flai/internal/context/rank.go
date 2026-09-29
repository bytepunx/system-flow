package context

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/bytepunx/system-flow/flai/internal/search"
	"github.com/bytepunx/system-flow/flai/internal/topics"
)

// Ranked is how many ADRs and design sections the ranked step adds
// (ADR-0047; the count is from the S-0125 replay).
const Ranked = 5

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

// Rank adds the adrs ADRs and the sections design sections that rank
// highest by BM25 against query among those nothing has chosen, with the
// reason "rank n". Sections are indexed one by one, each down to the next
// heading of any level; an ADR ranks by its best section and loads whole.
// Superseded ADRs, and sections with nothing below their heading, are left
// out of the index.
func (s *Selection) Rank(query string, adrs, sections int) {
	type key struct {
		doc *Doc
		sec int
	}
	var (
		keys []key
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
			keys = append(keys, key{d, i})
		}
	}
	ix := search.Build(docs)
	nADR, nSec := 0, 0
	for _, h := range ix.Search(strings.Join(Terms(query), " "), true, ix.Size()) {
		if nADR >= adrs && nSec >= sections {
			break
		}
		n, _ := strconv.Atoi(h.Path)
		k := keys[n]
		switch {
		case k.doc.Kind == KindADR:
			if nADR < adrs && !s.Loaded(k.doc) {
				nADR++
				s.choose(k.doc, k.doc.all(), StepRanked, fmt.Sprintf("rank %d", nADR))
			}
		case nSec < sections:
			nSec++
			s.choose(k.doc, []int{k.sec}, StepRanked, fmt.Sprintf("rank %d", nSec))
		}
	}
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
