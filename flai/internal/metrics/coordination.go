package metrics

import (
	"math"
	"regexp"
	"strings"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/messages"
	"github.com/bytepunx/system-flow/flai/internal/storygit"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// Coordination is, per ISO week of the window, the conversations between
// stories and how they ended, the conflicts a trial merge found against those
// a rebase met at sync or at acceptance, and the stories started on a share
// with the hold each would have waited (S-0337, ADR-0138).
type Coordination struct {
	// Logged is whether the conflict log could be read: true when it is
	// missing, false when it could not be read.
	Logged bool               `json:"logged"`
	Weeks  []CoordinationWeek `json:"weeks"`
}

// CoordinationWeek is the coordination counted in one ISO week.
type CoordinationWeek struct {
	Week          string             `json:"week"`
	Start         string             `json:"start"`
	Conversations ConversationCounts `json:"conversations"`
	// Conflicts is nil when the conflict log could not be read.
	Conflicts *ConflictCounts `json:"conflicts,omitempty"`
	Shares    ShareCounts     `json:"shares"`
}

// ConversationCounts is how many conversations opened, first agreed, first
// escalated, and closed with a story in a week.
type ConversationCounts struct {
	Opened          int `json:"opened"`
	Agreed          int `json:"agreed"`
	Escalated       int `json:"escalated"`
	ClosedWithStory int `json:"closed_with_story"`
}

// ConflictCounts is the conflict log's lines in a week, by where the conflict
// was met, and the paths they name.
type ConflictCounts struct {
	Found          int `json:"found"`
	FoundPaths     int `json:"found_paths"`
	MetSync        int `json:"met_sync"`
	MetSyncPaths   int `json:"met_sync_paths"`
	MetAccept      int `json:"met_accept"`
	MetAcceptPaths int `json:"met_accept_paths"`
}

// ShareCounts is the stories started on a share in a week and the hold they
// would have waited without it, in whole seconds.
type ShareCounts struct {
	Started     int     `json:"started"`
	HeldSeconds float64 `json:"held_seconds"`
}

// escalatedPrefix starts the entry flai message escalate adds.
const escalatedPrefix = "Asked the operator on"

// The closing entries' reasons that count: two branches merging cleanly at a
// sync, as agreement, and a story accepted, archived, or no longer open, as
// closed with its story.
var (
	closedClean     = regexp.MustCompile(`^Closed: \S+ and \S+ merge cleanly at the sync of S-\d{3,}$`)
	closedWithStory = regexp.MustCompile(`^Closed: (?:S-\d{3,} was|S-\d{3,} and S-\d{3,} were) (?:accepted|archived)$|^Closed: S-\d{3,} is \S+, no longer open, at the sync of S-\d{3,}$`)
)

// coordination counts, in the ISO weeks from the one that holds start to the
// one that holds now, the events of convs, of the conflict log's lines when
// logged, and of the stories among all started on a share, each in the week
// that holds it and only from start to now.
func coordination(all []*workitem.Item, convs []*messages.Conversation, conflicts []storygit.Conflict, logged bool, start, now time.Time) Coordination {
	out := Coordination{Logged: logged, Weeks: []CoordinationWeek{}}
	at := map[string]int{}
	for i, m := 0, monday(start); !m.After(now); i, m = i+1, m.AddDate(0, 0, 7) {
		y, w := m.ISOWeek()
		at[weekKey(y, w)] = i
		wk := CoordinationWeek{Week: weekKey(y, w), Start: m.Format("2006-01-02")}
		if logged {
			wk.Conflicts = &ConflictCounts{}
		}
		out.Weeks = append(out.Weeks, wk)
	}
	week := func(t time.Time) *CoordinationWeek {
		if t.IsZero() || t.Before(start) || t.After(now) {
			return nil
		}
		y, w := t.ISOWeek()
		i, ok := at[weekKey(y, w)]
		if !ok {
			return nil
		}
		return &out.Weeks[i]
	}
	stories := map[string]*workitem.Item{}
	for _, it := range all {
		if it.Type == workitem.Story {
			stories[workitem.CanonicalID(it.ID)] = it
		}
	}
	for _, c := range convs {
		ev := conversationEvents(c, stories)
		if wk := week(ev.opened); wk != nil {
			wk.Conversations.Opened++
		}
		if wk := week(ev.agreed); wk != nil {
			wk.Conversations.Agreed++
		}
		if wk := week(ev.escalated); wk != nil {
			wk.Conversations.Escalated++
		}
		if wk := week(ev.closedWithStory); wk != nil {
			wk.Conversations.ClosedWithStory++
		}
	}
	if logged {
		for _, l := range conflicts {
			wk := week(l.At.UTC())
			if wk == nil {
				continue
			}
			switch l.Where {
			case storygit.WhereTrial:
				wk.Conflicts.Found++
				wk.Conflicts.FoundPaths += len(l.Paths)
			case storygit.WhereSync:
				wk.Conflicts.MetSync++
				wk.Conflicts.MetSyncPaths += len(l.Paths)
			case storygit.WhereAccept:
				wk.Conflicts.MetAccept++
				wk.Conflicts.MetAcceptPaths += len(l.Paths)
			}
		}
	}
	for _, s := range startedOnShare(stories, convs, now) {
		if wk := week(s.started); wk != nil {
			wk.Shares.Started++
			wk.Shares.HeldSeconds += s.held
		}
	}
	return out
}

// convEvents is when a conversation opened, first agreed, was first
// escalated, and closed with one of its stories; zero for what did not
// happen.
type convEvents struct {
	opened, agreed, escalated, closedWithStory time.Time
}

// conversationEvents reads c's events from its entries and shares. A
// conversation that reads as closed, with no closing entry, is closed with
// its story at the first done or cancelled transition of either story.
func conversationEvents(c *messages.Conversation, stories map[string]*workitem.Item) convEvents {
	var ev convEvents
	entries := c.Entries()
	if len(entries) > 0 {
		ev.opened = parseTime(entries[0].At)
	}
	for _, sh := range c.Shares {
		ev.agreed = earliest(ev.agreed, parseTime(sh.At))
	}
	closing := false
	for _, e := range entries {
		if ev.escalated.IsZero() && e.Story != "" && strings.HasPrefix(e.Text, escalatedPrefix) {
			ev.escalated = parseTime(e.At)
		}
		if e.Story != "" || closing {
			continue
		}
		closing = true
		reason := strings.TrimSpace(strings.SplitN(e.Text, "\n", 2)[0])
		switch {
		case closedClean.MatchString(reason):
			ev.agreed = earliest(ev.agreed, parseTime(e.At))
		case closedWithStory.MatchString(reason):
			ev.closedWithStory = parseTime(e.At)
		}
	}
	if closing || !readsClosed(c, stories) {
		return ev
	}
	for _, id := range []string{c.From, c.To} {
		if it := stories[workitem.CanonicalID(id)]; it != nil {
			ev.closedWithStory = earliest(ev.closedWithStory, firstLeft(it))
		}
	}
	return ev
}

// readsClosed reports whether c reads as closed, as messages.Closed reads it
// from the items: its status is closed, or either story is done, cancelled,
// archived, or not among them.
func readsClosed(c *messages.Conversation, stories map[string]*workitem.Item) bool {
	if c.Status == messages.StatusClosed {
		return true
	}
	for _, id := range []string{c.From, c.To} {
		it := stories[workitem.CanonicalID(id)]
		if it == nil || it.Archived || it.Status == workitem.Done || it.Status == workitem.Cancelled {
			return true
		}
	}
	return false
}

// firstLeft is the time of the story's first move to done or cancelled, zero
// when it has none.
func firstLeft(it *workitem.Item) time.Time {
	for _, tr := range it.Transitions {
		if tr.To == workitem.Done || tr.To == workitem.Cancelled {
			return parseTime(tr.At)
		}
	}
	return time.Time{}
}

// shareStart is a story started on a share: its first move to in-progress
// and the seconds it would have waited without the share.
type shareStart struct {
	started time.Time
	held    float64
}

// startedOnShare is every story whose first move to in-progress fell while a
// share named it as held: the share made at or before the move, by a holder
// in progress at it. Its hold runs from the move to the holder's first move
// to review, done, or cancelled after it, or to now when there is none; with
// two such shares, the longer.
func startedOnShare(stories map[string]*workitem.Item, convs []*messages.Conversation, now time.Time) []shareStart {
	byHeld := map[string][]workitem.Share{}
	for _, c := range convs {
		for _, sh := range c.Shares {
			held := workitem.CanonicalID(sh.Held)
			byHeld[held] = append(byHeld[held], sh)
		}
	}
	var out []shareStart
	for id, shares := range byHeld {
		it := stories[id]
		if it == nil {
			continue
		}
		moved := it.FirstAt(workitem.InProgress)
		if moved.IsZero() {
			continue
		}
		var s *shareStart
		for _, sh := range shares {
			shared, holder := parseTime(sh.At), stories[workitem.CanonicalID(sh.Holder)]
			if shared.IsZero() || shared.After(moved) || holder == nil || stateAt(holder, moved) != workitem.InProgress {
				continue
			}
			held := math.Round(holdEnd(holder, moved, now).Sub(moved).Seconds())
			if s == nil {
				s = &shareStart{started: moved}
			}
			s.held = math.Max(s.held, held)
		}
		if s != nil {
			out = append(out, *s)
		}
	}
	return out
}

// holdEnd is the holder's first move to review, done, or cancelled after
// from, or now when it has none.
func holdEnd(holder *workitem.Item, from, now time.Time) time.Time {
	for _, tr := range holder.Transitions {
		t := parseTime(tr.At)
		if t.After(from) && (tr.To == workitem.Review || tr.To == workitem.Done || tr.To == workitem.Cancelled) {
			return t
		}
	}
	return now
}

// parseTime is a timestamp in workitem.TimeFormat, zero when it does not
// parse.
func parseTime(s string) time.Time {
	t, err := time.Parse(workitem.TimeFormat, s)
	if err != nil {
		return time.Time{}
	}
	return t
}

// earliest is the earlier of a and b, a zero time counting as none.
func earliest(a, b time.Time) time.Time {
	if a.IsZero() || (!b.IsZero() && b.Before(a)) {
		return b
	}
	return a
}
