package workitem

import (
	"fmt"
	"sort"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/pending"
)

// BoardCard is one item on the board as flai board --json and the MCP board
// tool print it.
type BoardCard struct {
	ID     string `json:"id"`
	Type   string `json:"type"`
	Title  string `json:"title"`
	Nature string `json:"nature"`
	Parent string `json:"parent,omitempty"`
	// ParentTitle, Status, and EnteredAt are what the dashboard's board shows
	// beyond the CLI's (S-0073): a card names its parent without a second
	// request, and its age is counted from EnteredAt in the browser.
	ParentTitle string   `json:"parent_title,omitempty"`
	Status      string   `json:"status"`
	EnteredAt   string   `json:"entered_at"`
	Blocked     bool     `json:"blocked"`
	Age         string   `json:"age_in_column"`
	AgeSecs     int64    `json:"age_in_column_seconds"`
	Touches     []string `json:"touches,omitempty"`
	// Held is why a ready story is not started or offered (S-0128).
	Held *Hold `json:"held,omitempty"`
}

// BoardView is the board with its cards: one reading of the repository for
// the command line and for agents over MCP (S-0058).
type BoardView struct {
	Columns   map[string][]BoardCard `json:"columns"`
	WIPLimits map[string]int         `json:"wip_limits"`
	Order     []string               `json:"order"`
	Breaches  []string               `json:"breaches"`
	// Unpushed is set by callers that can ask git: an acceptance in the main
	// checkout that its remote-tracking branch does not have yet (S-0063).
	Unpushed *pending.Unpushed `json:"unpushed,omitempty"`
	// Counts are stories per column, which is what the limits apply to.
	Counts map[string]int `json:"-"`
}

// NewBoardView lays the active items out by column. Stories only unless all.
// A ready story whose claim overlaps an open story's is marked held (S-0128);
// projects are the manifest's sub-projects, whose names a claim reads as paths.
// Acceptance archives a story the moment it merges (S-0087); pendingPublish
// names the stories a release has not yet covered, so a done, archived story
// still on the done column until it is published, instead of vanishing the
// instant it is accepted, before anyone has had the chance to see it there.
func NewBoardView(items []*Item, board *Board, now time.Time, all bool, pendingPublish map[string]bool, projects []manifest.Project) BoardView {
	v := BoardView{Columns: map[string][]BoardCard{}, WIPLimits: board.WIPLimits, Order: board.Order, Counts: map[string]int{}}
	holds := NewHolds(items, projects)
	// A parent is looked up among every item given, archived ones included.
	titles := map[string]string{}
	for _, it := range items {
		titles[it.ID] = it.Title
	}
	for _, it := range items {
		archivedButPending := it.Archived && it.Status == Done && pendingPublish[it.ID]
		if (it.Archived && !archivedButPending) || (!all && it.Type != Story) {
			continue
		}
		age := now.Sub(it.EnteredAt())
		card := BoardCard{
			ID: it.ID, Type: it.Type, Title: it.Title, Nature: it.Nature, Parent: it.Parent,
			ParentTitle: titles[it.Parent], Status: it.Status, EnteredAt: it.EnteredAt().UTC().Format(TimeFormat),
			Blocked: it.IsBlocked(), Age: HumanDuration(age), AgeSecs: int64(age.Seconds()),
			Touches: it.Touches,
		}
		if it.Type == Story && it.Status == Ready && !it.Archived {
			card.Held = holds.Of(it)
		}
		v.Columns[it.Status] = append(v.Columns[it.Status], card)
		if it.Type == Story {
			v.Counts[it.Status]++
		}
	}
	// Backlog and ready read in pull order: the stories take the sequence, in
	// the places stories already hold, and epics and tasks stay where they are.
	for _, st := range []string{Backlog, Ready} {
		rank := map[string]int{}
		for i, id := range PullSequence(board.Order, items, st) {
			rank[id] = i
		}
		cards := v.Columns[st]
		var slots []int
		var stories []BoardCard
		for i, c := range cards {
			if c.Type == Story {
				slots = append(slots, i)
				stories = append(stories, c)
			}
		}
		sort.SliceStable(stories, func(i, j int) bool { return rank[stories[i].ID] < rank[stories[j].ID] })
		for i, slot := range slots {
			cards[slot] = stories[i]
		}
	}
	for _, st := range States {
		if limit, ok := board.WIPLimits[st]; ok && limit > 0 && v.Counts[st] > limit {
			v.Breaches = append(v.Breaches, fmt.Sprintf("%s has %d stories, limit %d", st, v.Counts[st], limit))
		}
	}
	return v
}

// ReadyInPullOrder returns the ready stories, those named in the pull order
// first and in that order, the rest after them by ID. The view's columns are
// already in that sequence.
func (v BoardView) ReadyInPullOrder() []BoardCard {
	var out []BoardCard
	for _, c := range v.Columns[Ready] {
		if c.Type == Story {
			out = append(out, c)
		}
	}
	return out
}

// FirstPullable is the first ready story in pull order that is not held, or
// nil when every one is (S-0128).
func (v BoardView) FirstPullable() *BoardCard { return FirstClear(v.ReadyInPullOrder()) }

// FirstClear is the first of cards that is not held, or nil.
func FirstClear(cards []BoardCard) *BoardCard {
	for _, c := range cards {
		if c.Held == nil {
			return &c
		}
	}
	return nil
}

// CanPull reports whether the in-progress limit leaves room for one more story.
func (v BoardView) CanPull() bool {
	limit, ok := v.WIPLimits[InProgress]
	return !ok || limit <= 0 || v.Counts[InProgress] < limit
}

// HumanDuration is an age for people: now, 12m, 5h, 3d.
func HumanDuration(d time.Duration) string {
	switch {
	case d < time.Minute:
		return "now"
	case d < time.Hour:
		return fmt.Sprintf("%dm", int(d.Minutes()))
	case d < 48*time.Hour:
		return fmt.Sprintf("%dh", int(d.Hours()))
	default:
		return fmt.Sprintf("%dd", int(d.Hours()/24))
	}
}
