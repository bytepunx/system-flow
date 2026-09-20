// Package search is full-text search over a project's Markdown, for the
// dashboard's search page (S-0074). It indexes what MiniSearch indexed in
// the dashboard before it: item ID, title, tags, headings, and body, with
// the same boosts, prefix and fuzzy matching, and BM25 ranking. It is small
// on purpose: a few thousand documents index in well under a second, and a
// search engine would be a dependency many times the size of the problem.
package search

import (
	"math"
	"sort"
	"strings"
	"unicode"
)

// Doc is one Markdown file as it is indexed.
type Doc struct {
	Path     string
	Kind     string // item or doc
	ItemID   string
	Title    string
	Tags     string
	Headings string
	Body     string
	Type     string
	Nature   string
	Status   string
	Scope    string // design, wip, or docs
}

// Hit is one result, with what the search page shows.
type Hit struct {
	Path    string  `json:"path"`
	Kind    string  `json:"kind"`
	ItemID  string  `json:"itemId,omitempty"`
	Title   string  `json:"title"`
	Scope   string  `json:"scope"`
	Status  string  `json:"status,omitempty"`
	Type    string  `json:"type,omitempty"`
	Nature  string  `json:"nature,omitempty"`
	Score   float64 `json:"score"`
	Snippet string  `json:"snippet"`
}

const (
	fItemID = iota
	fTitle
	fTags
	fHeadings
	fBody
	nFields
)

var boost = [nFields]float64{fItemID: 4, fTitle: 3, fTags: 1, fHeadings: 2, fBody: 1}

// MiniSearch's defaults: BM25+ and the weights of inexact matches.
const (
	bm25K        = 1.2
	bm25B        = 0.7
	bm25D        = 0.5
	prefixWeight = 0.375
	fuzzyWeight  = 0.45
	fuzziness    = 0.2
	maxFuzzy     = 6
)

type posting struct {
	doc   int
	field int
	tf    int
}

// Index answers queries over the documents it was built from.
type Index struct {
	docs     []Doc
	postings map[string][]posting
	vocab    []string         // sorted, for prefixes
	lens     [][nFields]int   // terms per field per document
	avg      [nFields]float64 // mean terms per field
	df       map[string][nFields]int
}

// Tokens are lower-cased runs of letters and digits.
func Tokens(s string) []string {
	return strings.FieldsFunc(strings.ToLower(s), func(r rune) bool { return !unicode.IsLetter(r) && !unicode.IsDigit(r) })
}

// Build indexes the documents.
func Build(docs []Doc) *Index {
	ix := &Index{docs: docs, postings: map[string][]posting{}, lens: make([][nFields]int, len(docs)), df: map[string][nFields]int{}}
	var total [nFields]int
	for d, doc := range docs {
		for f, text := range [nFields]string{doc.ItemID, doc.Title, doc.Tags, doc.Headings, doc.Body} {
			tf := map[string]int{}
			toks := Tokens(text)
			for _, t := range toks {
				tf[t]++
			}
			ix.lens[d][f] = len(toks)
			total[f] += len(toks)
			for t, n := range tf {
				ix.postings[t] = append(ix.postings[t], posting{doc: d, field: f, tf: n})
				c := ix.df[t]
				c[f]++
				ix.df[t] = c
			}
		}
	}
	for f := range total {
		if len(docs) > 0 {
			ix.avg[f] = float64(total[f]) / float64(len(docs))
		}
	}
	ix.vocab = make([]string, 0, len(ix.postings))
	for t := range ix.postings {
		ix.vocab = append(ix.vocab, t)
	}
	sort.Strings(ix.vocab)
	return ix
}

// Size is the number of documents indexed.
func (ix *Index) Size() int { return len(ix.docs) }

type match struct {
	term   string
	weight float64
}

// expand finds the indexed terms a query term matches: itself, the terms it
// is a prefix of, and the terms within a fifth of its length in edits.
func (ix *Index) expand(q string) []match {
	best := map[string]float64{}
	add := func(term string, w float64) {
		if w > best[term] {
			best[term] = w
		}
	}
	if _, ok := ix.postings[q]; ok {
		add(q, 1)
	}
	n := len([]rune(q))
	// A single character is not a prefix worth following: the "s" of S-0052
	// would bring in every word that starts with it.
	for i := sort.SearchStrings(ix.vocab, q); n > 1 && i < len(ix.vocab) && strings.HasPrefix(ix.vocab[i], q); i++ {
		if extra := len([]rune(ix.vocab[i])) - n; extra > 0 {
			add(ix.vocab[i], prefixWeight*float64(n)/(float64(n)+0.3*float64(extra)))
		}
	}
	if max := int(math.Min(maxFuzzy, math.Round(float64(n)*fuzziness))); max > 0 {
		qr := []rune(q)
		for _, t := range ix.vocab {
			tr := []rune(t)
			if d := len(tr) - n; d > max || -d > max || t == q {
				continue
			}
			if d := editDistance(qr, tr, max); d > 0 && d <= max {
				add(t, fuzzyWeight*float64(n)/(float64(n)+float64(d)))
			}
		}
	}
	out := make([]match, 0, len(best))
	for t, w := range best {
		out = append(out, match{t, w})
	}
	return out
}

// editDistance is Levenshtein distance, given up on past max.
func editDistance(a, b []rune, max int) int {
	prev := make([]int, len(b)+1)
	cur := make([]int, len(b)+1)
	for j := range prev {
		prev[j] = j
	}
	for i := 1; i <= len(a); i++ {
		cur[0] = i
		low := cur[0]
		for j := 1; j <= len(b); j++ {
			cost := 1
			if a[i-1] == b[j-1] {
				cost = 0
			}
			cur[j] = min(prev[j]+1, cur[j-1]+1, prev[j-1]+cost)
			if cur[j] < low {
				low = cur[j]
			}
		}
		if low > max {
			return max + 1
		}
		prev, cur = cur, prev
	}
	return prev[len(b)]
}

// Search ranks the documents for a query. Documentation under docs is left
// out unless includeDocs. Terms are combined with OR, so a document that
// matches more of them scores higher.
func (ix *Index) Search(query string, includeDocs bool, limit int) []Hit {
	terms := Tokens(query)
	if len(terms) == 0 || len(ix.docs) == 0 {
		return []Hit{}
	}
	if limit <= 0 {
		limit = 30
	}
	n := float64(len(ix.docs))
	scores := map[int]float64{}
	matched := map[int]map[string]bool{}
	seen := map[string]bool{}
	for _, q := range terms {
		if seen[q] {
			continue
		}
		seen[q] = true
		for _, m := range ix.expand(q) {
			df := ix.df[m.term]
			for _, p := range ix.postings[m.term] {
				if !includeDocs && ix.docs[p.doc].Scope == "docs" {
					continue
				}
				idf := math.Log(1 + (n-float64(df[p.field])+0.5)/(float64(df[p.field])+0.5))
				norm := 1.0
				if ix.avg[p.field] > 0 {
					norm = 1 - bm25B + bm25B*float64(ix.lens[p.doc][p.field])/ix.avg[p.field]
				}
				tf := float64(p.tf)
				scores[p.doc] += m.weight * boost[p.field] * idf * (bm25D + tf*(bm25K+1)/(tf+bm25K*norm))
				if matched[p.doc] == nil {
					matched[p.doc] = map[string]bool{}
				}
				matched[p.doc][m.term] = true
			}
		}
	}
	order := make([]int, 0, len(scores))
	for d := range scores {
		order = append(order, d)
	}
	sort.Slice(order, func(i, j int) bool {
		if scores[order[i]] != scores[order[j]] {
			return scores[order[i]] > scores[order[j]]
		}
		return ix.docs[order[i]].Path < ix.docs[order[j]].Path
	})
	if len(order) > limit {
		order = order[:limit]
	}
	hits := make([]Hit, 0, len(order))
	for _, d := range order {
		doc := ix.docs[d]
		found := make([]string, 0, len(matched[d]))
		for t := range matched[d] {
			found = append(found, t)
		}
		hits = append(hits, Hit{Path: doc.Path, Kind: doc.Kind, ItemID: doc.ItemID, Title: doc.Title, Scope: doc.Scope, Status: doc.Status,
			Type: doc.Type, Nature: doc.Nature, Score: math.Round(scores[d]*100) / 100, Snippet: Snippet(doc.Body, found, 160)})
	}
	return hits
}

// Snippet is a short window of the body around the first matched term.
func Snippet(body string, terms []string, width int) string {
	text := []rune(strings.Join(strings.Fields(body), " "))
	lower := []rune(strings.ToLower(string(text)))
	at := -1
	for _, t := range terms {
		if i := indexRunes(lower, []rune(strings.ToLower(t))); i >= 0 && (at < 0 || i < at) {
			at = i
		}
	}
	if at < 0 {
		if len(text) > width {
			return string(text[:width]) + "…"
		}
		return string(text)
	}
	start := max(0, at-width/3)
	end := min(len(text), start+width)
	out := string(text[start:end])
	if start > 0 {
		out = "…" + out
	}
	if end < len(text) {
		out += "…"
	}
	return out
}

func indexRunes(haystack, needle []rune) int {
	if len(needle) == 0 || len(needle) > len(haystack) {
		return -1
	}
	for i := 0; i+len(needle) <= len(haystack); i++ {
		j := 0
		for j < len(needle) && haystack[i+j] == needle[j] {
			j++
		}
		if j == len(needle) {
			return i
		}
	}
	return -1
}
