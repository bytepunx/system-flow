package mergequeue

import (
	"slices"
	"time"
)

// What a verified batch comes to (S-0400, ADR-0176): each story's heads, the
// bisection of a batch that failed, and the landing of the stories it lets
// land, kept on the batch in the queue's state.

// The outcomes a batch's end records: a story the queue landed, and the story
// a red batch was bisected to, sent back to in-progress.
const (
	OutcomeLanded  = "landed"
	OutcomeFailing = "failing"
)

// Head is a batched story's heads: Verified, the head its last passing verify
// was at, of which the batch rebased a copy, and Rebased, that copy's head in
// the batch, the commit the batch verified it at.
type Head struct {
	Story    string `json:"story"`
	Verified string `json:"verified"`
	Rebased  string `json:"rebased"`
}

// Bisect is the bisection of a batch that failed: each step, in the order it
// ran; Good, the most stories, from the first, known to pass with main (0,
// main itself, at first); Bad, the fewest known to fail (the whole batch at
// first); and, once found, Failing, the first story with which the batch
// fails, with the failing commit and the result there.
type Bisect struct {
	Steps   []BisectStep `json:"steps"`
	Good    int          `json:"good"`
	Bad     int          `json:"bad"`
	Failing string       `json:"failing,omitempty"`
	Commit  string       `json:"commit,omitempty"`
	Result  *Result      `json:"result,omitempty"`
	// SentBack is whether the failing story was sent back.
	SentBack bool `json:"sent_back,omitempty"`
}

// BisectStep is one step of a bisection: the stories, the batch's first, that
// it verified with main, the commit it verified them at, when it started, in
// UTC, and its result, nil while it runs.
type BisectStep struct {
	Stories   []string  `json:"stories"`
	Commit    string    `json:"commit"`
	StartedAt time.Time `json:"started_at"`
	Result    *Result   `json:"result,omitempty"`
}

// Found reports whether the bisection has found the failing story.
func (b *Bisect) Found() bool { return b != nil && b.Failing != "" }

// Landing is the landing of a verified batch: the stories it lands, in queue
// order, the whole batch when it passed and the stories before the failing one
// when it was bisected; those landed so far; and Head, the main branch's
// commit the next landing must find, the batch's base at first and then the
// head the last landing left.
type Landing struct {
	Stories []string `json:"stories"`
	Landed  []string `json:"landed"`
	Head    string   `json:"head"`
}

// Next is the next story to land, and false when every story has landed.
func (l *Landing) Next() (string, bool) {
	if l == nil {
		return "", false
	}
	for _, s := range l.Stories {
		if !slices.ContainsFunc(l.Landed, func(x string) bool { return sameStory(x, s) }) {
			return s, true
		}
	}
	return "", false
}

// HeadOf is story's heads in the batch, and false when the batch kept none
// for it.
func (b *Batch) HeadOf(story string) (Head, bool) {
	for _, h := range b.Heads {
		if sameStory(h.Story, story) {
			return h, true
		}
	}
	return Head{}, false
}
