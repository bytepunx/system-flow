package context

import (
	"cmp"
	"slices"
	"strconv"
	"strings"

	"github.com/bytepunx/system-flow/flai/internal/digest"
	"github.com/bytepunx/system-flow/flai/internal/search"
	"github.com/bytepunx/system-flow/flai/internal/topics"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// DigestShare is the digest step's share of a pack's budget: its entries
// together print within Budget/DigestShare bytes, 10 KB of 80 KB
// (ADR-0194).
const DigestShare = 8

// The statuses of the digest entries a pack selects: an active entry loads
// its body, a compacted one its title only. An absorbed entry is never
// selected; the document that absorbed it says it now.
const (
	digestActive    = "active"
	digestCompacted = "compacted"
)

// pick is a digest entry a pack selected, with what selected it and how it
// scores against what the pack is ranked against.
type pick struct {
	entry  digest.Entry
	shared []string // the pack's topics the entry has, all left out
	path   string   // the entry of the claim one of its paths meets, empty when none
	score  float64  // BM25 against the query, 0 when it has no word in common with it
}

// selectDigest is the entries a story with these topics and this claim
// selects (ADR-0194): those active or compacted that share a topic other
// than all with the story, or that have a path equal to an entry of the
// claim, a folder holding it, or held by it. They are in the order given.
func selectDigest(entries []digest.Entry, packTopics, claim []string) []pick {
	var out []pick
	for _, e := range entries {
		if e.Status != digestActive && e.Status != digestCompacted {
			continue
		}
		p := pick{entry: e, path: claimMet(e.Paths, claim)}
		for _, t := range e.Topics {
			if t != topics.All && slices.Contains(packTopics, t) && !slices.Contains(p.shared, t) {
				p.shared = append(p.shared, t)
			}
		}
		if p.path != "" || len(p.shared) > 0 {
			out = append(out, p)
		}
	}
	return out
}

// claimMet is the first entry of claim that one of paths meets, being equal
// to it, a folder holding it, or held by it; empty when none does.
func claimMet(paths, claim []string) string {
	for _, c := range claim {
		for _, p := range paths {
			if c != "" && p != "" && workitem.PathsOverlap(c, p) {
				return c
			}
		}
	}
	return ""
}

// scoreDigest sets each pick's score: BM25 against query over the entry's
// title and body, with the index and the query's terms the ranked step uses.
// A pick with no word in common with the query keeps a score of 0.
func scoreDigest(picks []pick, query string) {
	terms := search.Terms(query)
	if len(picks) == 0 || len(terms) == 0 {
		return
	}
	docs := make([]search.Doc, len(picks))
	for i, p := range picks {
		docs[i] = search.Doc{Path: strconv.Itoa(i), Kind: "doc", Title: p.entry.Title, Body: p.entry.Body, Scope: "design"}
	}
	for _, h := range search.Build(docs).Search(strings.Join(terms, " "), true, len(docs)) {
		if n, err := strconv.Atoi(h.Path); err == nil {
			picks[n].score = h.Score
		}
	}
}

// rankDigest is picks in the order a story's pack tries them (ADR-0194):
// those whose paths meet the claim first, then those its topics selected;
// within each by score, then the newest date, then by ID. A pick that meets
// no claimed path and has no word in common with the query is left out.
func rankDigest(picks []pick) []pick {
	out := slices.DeleteFunc(slices.Clone(picks), func(p pick) bool { return p.path == "" && p.score == 0 })
	slices.SortStableFunc(out, func(a, b pick) int {
		if (a.path == "") != (b.path == "") {
			if a.path != "" {
				return -1
			}
			return 1
		}
		return cmp.Or(cmp.Compare(b.score, a.score), cmp.Compare(b.entry.Date, a.entry.Date), cmp.Compare(a.entry.ID, b.entry.ID))
	})
	return out
}

// reason is what selected the pick, at rank n: "path: <the claimed path it
// meets>; rank n" when its paths meet the claim, else "topics: <the topics
// it shares>; rank n". A compacted entry's ends "; compacted".
func (p pick) reason(n int) string {
	r := "topics: " + strings.Join(p.shared, ", ")
	if p.path != "" {
		r = "path: " + p.path
	}
	r += "; rank " + strconv.Itoa(n)
	if p.entry.Status == digestCompacted {
		r += "; " + digestCompacted
	}
	return r
}

// digestFor is the digest items a story's pack tries, best first, as
// rankDigest orders the entries packTopics and claim select against query,
// each with its reason. With whole, an active entry with a body loads it;
// otherwise, and for a compacted entry, each is a brief.
func digestFor(entries []digest.Entry, packTopics, claim []string, query string, whole bool) []Item {
	picks := selectDigest(entries, packTopics, claim)
	scoreDigest(picks, query)
	var out []Item
	for i, p := range rankDigest(picks) {
		out = append(out, digestItem(p.entry, p.reason(i+1), whole))
	}
	return out
}

// digestItem is entry e as a pack holds it, with step digest and reason.
// With whole, an active entry with a body loads it under its "# D-nnnn
// <title>" heading as Text; otherwise the item is a brief, with no Text,
// printed as one line, and Whole is the size of what reading it adds, 0
// for an entry with no body.
func digestItem(e digest.Entry, reason string, whole bool) Item {
	it := Item{Path: e.Path, ID: e.ID, Title: e.Title, Step: StepDigest, Reason: reason}
	if e.Body == "" {
		return it
	}
	text := "# " + e.ID + " " + e.Title + "\n\n" + e.Body + "\n"
	if whole && e.Status == digestActive {
		it.Text, it.Size = text, len(text)
		return it
	}
	it.Whole = len(text)
	return it
}

// addDigest is the digest step of a story's pack (ADR-0194), after the
// briefs and before anything is ranked: it adds items in order while they
// together print within the step's share of the budget and the pack within
// the budget, skipping one that does not fit and trying the next. It
// returns what it added.
func (p *Pack) addDigest(items []Item) []Item {
	share, used := p.Budget/DigestShare, 0
	var added []Item
	for _, it := range items {
		n := len(p.printed(it))
		if used+n > share || p.Size.Bytes+n > p.Budget {
			continue
		}
		size := p.Size
		p.Items = append(p.Items, it)
		p.size()
		if p.Size.Bytes > p.Budget {
			p.Items, p.Size = p.Items[:len(p.Items)-1], size
			continue
		}
		used += n
		added = append(added, it)
	}
	return added
}

// digestBrief prints a digest item that loads no body: its ID, title,
// path, the size of its body when it has one, and its reason, on one line.
func (it Item) digestBrief() string {
	if it.Whole == 0 {
		return "\n- " + it.ID + " " + it.Title + " (" + it.Path + "; " + it.Reason + ")\n"
	}
	return "\n- " + it.ID + " " + it.Title + " (" + it.Path + ", " + strconv.Itoa(it.Whole) + " bytes; " + it.Reason + ")\n"
}

// digestRow is the header's contents row naming n digest entries, whole of
// them with their bodies.
func digestRow(n, whole int) string {
	entries := strconv.Itoa(n) + " digest entries"
	switch {
	case n == 1 && whole == 1:
		return "1 digest entry with its body"
	case n == 1:
		return "1 digest entry as a one-line brief"
	case whole == n:
		return entries + " with their bodies"
	case whole == 0:
		return entries + " as one-line briefs"
	}
	return entries + ", " + strconv.Itoa(whole) + " with their bodies, the rest as one-line briefs"
}
