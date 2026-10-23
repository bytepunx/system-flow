package digest

import (
	"regexp"
	"slices"
	"strings"

	"github.com/bytepunx/system-flow/flai/internal/threads"
)

// Candidate is one thing past work learned that could become a digest entry:
// an operator's answer or overruled recommendation on a thread, a decision in
// an archived narrative, or a closed issue's cause and fix. Harvest lists
// them and Seed writes the answers, recommendations, and causes.
type Candidate struct {
	// Kind is one of Kinds.
	Kind string `json:"kind"`
	// Source is the item the fact comes from: S-nnnn, TH-nnnn, or I-nnnn.
	Source string `json:"source"`
	// Ref says where in the source: the narrative's path for a decision, the
	// entry's time for an answer or a recommendation, empty for a cause.
	Ref string `json:"ref,omitempty"`
	// Title is the fact, as one sentence, as an entry's title would be.
	Title string `json:"title"`
	// Topics and Paths are those of the item the fact comes from: a thread's
	// anchored item, a narrative's story, an issue's closing story.
	Topics []string `json:"topics"`
	Paths  []string `json:"paths"`
	// At is when the fact was learned, YYYY-MM-DDTHH:MM:SSZ where the source
	// has a time, else YYYY-MM-DD; Date is its day.
	At   string `json:"at"`
	Date string `json:"date"`
	// Question is what a thread asked, its opener's first entry; empty for a
	// decision or a cause.
	Question string `json:"question,omitempty"`
	// Text is the answer, the decision, or the issue's cause.
	Text string `json:"text"`
	// Overruled is the recommendation an overruled recommendation's answer
	// did not confirm; empty otherwise.
	Overruled string `json:"overruled,omitempty"`
	// Fix is a cause's fix; empty otherwise.
	Fix string `json:"fix,omitempty"`
	// CitedBy is the entries whose sources already name Source, filled by
	// Harvest.
	CitedBy []string `json:"cited_by"`
}

// ItemLookup gives the topics and paths of a work item by ID, archived or
// not: a story's topics with its epic's and tasks', and its touches. An item
// it cannot find has none.
type ItemLookup func(id string) (topics, paths []string)

var (
	// permissionTitle is the title permission_prompt gives the thread it
	// opens to ask the operator for a write.
	permissionTitle = regexp.MustCompile(`^Allow \S+ .+\?$`)
	// confirmed is the line threads.Confirm writes, naming the confirmed
	// recommendation's time and author.
	confirmed = regexp.MustCompile(`(?m)^Confirmed the recommendation of (\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z) (.+?)\.\s*$`)
	// listMarker is a markdown list item's marker at the start of a line.
	listMarker = regexp.MustCompile(`^(?:[-*+]|\d+[.)])\s+`)
)

// Answers is the operator's answers on a thread, in the order they were
// given: each entry by owner, the manifest's owner, that is not a
// resolution, on a thread someone else opened. A confirmation answers with
// the recommendation it confirms; an owner's reply to a recommendation still
// pending overrules it and is a recommendation candidate carrying both
// texts. A thread permission_prompt opened gives none. lookup gives the
// topics and paths of the item the thread is anchored to; a nil lookup, or
// an anchor that names no item, gives none.
func Answers(th *threads.Thread, owner string, lookup ItemLookup) []Candidate {
	if th == nil || owner == "" || th.Opener() == owner || permissionTitle.MatchString(th.Title) {
		return nil
	}
	entries := th.Entries()
	if len(entries) == 0 {
		return nil
	}
	var topics, paths []string
	if lookup != nil && th.Anchor.Item != "" {
		topics, paths = lookup(th.Anchor.Item)
	}
	question := strings.TrimSpace(entries[0].Text)
	var out []Candidate
	var recommendations []threads.Entry
	var pending *threads.Entry
	for _, e := range entries {
		if e.Author != owner {
			switch {
			case e.Confirms():
				pending = nil
			case e.Recommendation:
				recommendations = append(recommendations, e)
				pending = &e
			}
			continue
		}
		text := withoutResolution(e.Text)
		if text == "" {
			continue
		}
		c := Candidate{
			Kind:     "answer",
			Source:   threads.CanonicalID(th.ID),
			Ref:      e.At,
			Topics:   listOrEmpty(topics),
			Paths:    listOrEmpty(paths),
			At:       e.At,
			Date:     e.At[:10],
			Question: question,
			Text:     text,
			CitedBy:  []string{},
		}
		switch {
		case e.Confirms():
			r := confirmedRecommendation(e.Text, recommendations)
			pending = nil
			if r == nil {
				continue
			}
			c.Text = strings.TrimSpace(r.Text)
		case pending != nil:
			c.Kind = "recommendation"
			c.Overruled = strings.TrimSpace(pending.Text)
			pending = nil
		}
		c.Title = answerTitle(th.Title, c.Text)
		out = append(out, c)
	}
	return out
}

// withoutResolution is an entry's text without the note threads.Resolve
// writes, "Resolved." or "Resolved: <reason>", whether it is the whole entry
// or its last paragraph, joined to a reply made in the same second; empty
// when the note is all there is.
func withoutResolution(text string) string {
	text = strings.TrimSpace(text)
	if resolution(text) {
		return ""
	}
	if i := strings.LastIndex(text, "\n\n"); i >= 0 && resolution(text[i+2:]) {
		return strings.TrimSpace(text[:i])
	}
	return text
}

// resolution reports whether a paragraph is the note threads.Resolve writes.
func resolution(p string) bool {
	return p == "Resolved." || strings.HasPrefix(p, "Resolved: ")
}

// confirmedRecommendation is the recommendation a confirmation names by its
// time and author, or, when none matches, the latest recommendation before
// it; nil when there is none.
func confirmedRecommendation(text string, recommendations []threads.Entry) *threads.Entry {
	if m := confirmed.FindStringSubmatch(text); m != nil {
		for i := len(recommendations) - 1; i >= 0; i-- {
			if r := recommendations[i]; r.At == m[1] && r.Author == m[2] {
				return &recommendations[i]
			}
		}
	}
	if len(recommendations) == 0 {
		return nil
	}
	return &recommendations[len(recommendations)-1]
}

// answerTitle is a thread's question and an answer's first sentence as one
// sentence, "<title>: <first sentence>", cut on a word near titleLimit. The
// sentence is taken from the answer's first line, without a list marker or
// bold markers.
func answerTitle(title, answer string) string {
	line, _, _ := strings.Cut(strings.TrimSpace(answer), "\n")
	line = listMarker.ReplaceAllString(strings.ReplaceAll(line, "**", ""), "")
	return clipOnWord(strings.TrimRight(strings.TrimSpace(title), "?.") + ": " + firstSentence(line))
}

// listOrEmpty is a copy of list, empty rather than nil so a candidate prints
// it as [].
func listOrEmpty(list []string) []string {
	if list == nil {
		return []string{}
	}
	return slices.Clone(list)
}
