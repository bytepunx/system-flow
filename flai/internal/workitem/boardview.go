package workitem

import (
	"fmt"
	"sort"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/workflow"
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
	// Archived marks an accepted card shown only because its release is
	// pending (S-0087): the dashboard leaves it out while the clone's release
	// tags lag the remote's, when pending cannot be told from published
	// (S-0174).
	Archived bool `json:"archived,omitempty"`
	// Held is why a story in the pull's source is not started or offered
	// (S-0128), and on a hold on overlap alone, whose agent flai asked about
	// it (S-0338).
	Held *Hold `json:"held,omitempty"`
	// Shared is, for a story in the pull's source or pulled and not yet
	// closed, each share in force on which it is the held story: the paths
	// the story that held it shared, the split, and the conversation
	// (S-0338).
	Shared []Share `json:"shared,omitempty"`
	// Tasks is a story's task plan in counts, when it has tasks (S-0176).
	Tasks *TaskSummary `json:"tasks,omitempty"`
	// Draft marks a story an agent wrote that the operator has not finalized,
	// so the board can say so on its card (S-0201).
	Draft bool `json:"draft,omitempty"`
	// UndoesMain is, for a story whose record lists any, the paths its
	// branch still undoes of the main branch's change, sorted (ADR-0155,
	// S-0388). NewBoardView leaves it empty; the host's board.get fills it.
	UndoesMain []string `json:"undoes_main,omitempty"`
}

// Lane is one column of the board as its workflow defines it (ADR-0301 §2):
// what a reader of the view needs to draw it and judge it without a list of
// states of its own.
type Lane struct {
	Name string `json:"name"`
	// Title is the state's heading, its name when the workflow gives none.
	Title string `json:"title,omitempty"`
	Kind  string `json:"kind"`
	// WIPLimit is the column's limit on stories, 0 for none.
	WIPLimit int `json:"wip_limit,omitempty"`
	// TakesLimit says the column may carry a WIP limit, so one may be set.
	TakesLimit bool `json:"takes_limit,omitempty"`
	// Ordered says the column reads in the pull order.
	Ordered bool `json:"ordered,omitempty"`
}

// BoardView is the board with its cards: one reading of the repository for
// the command line and for agents over MCP (S-0058).
type BoardView struct {
	// Lanes are the board's columns in the workflow's order, each with its
	// kind and WIP limit (ADR-0301 §2).
	Lanes     []Lane                 `json:"lanes"`
	Columns   map[string][]BoardCard `json:"columns"`
	WIPLimits map[string]int         `json:"wip_limits"`
	Order     []string               `json:"order"`
	Breaches  []string               `json:"breaches"`
	// Unpublished names the accepted items no release has covered yet, by ID:
	// what publishing would send to the remote, which is the operator's to
	// do or to ask for (ADR-0067).
	Unpublished []string `json:"unpublished,omitempty"`
	// Counts are stories per column, which is what the limits apply to.
	Counts map[string]int `json:"-"`
	// flow is the workflow the view was laid out by; the zero value, as in
	// a view built by hand, reads as the baseline.
	flow workflow.Workflow
}

// boardFlow is the workflow the view was laid out by, the baseline when none
// was given.
func (v BoardView) boardFlow() workflow.Workflow {
	if len(v.flow.States) == 0 {
		return baselineFlow
	}
	return v.flow
}

// NewBoardView lays the active items out by column, the columns and their
// kinds, limits, and orders the board's workflow's (ADR-0301 §2). Stories
// only unless all. A story in the pull's source that holds judges held is
// marked so (S-0128): Repo.Holds of the same items, which reads the
// manifest's sub-projects and shared paths (ADR-0096), the shares in force,
// the asks about a hold (S-0338), and the repository's trial merge, by which
// an overlap with a started story holds only when it conflicts (Repo.Trial,
// ADR-0230); nil judges with none of them. A story in the pull's source, or
// in an active or waiting state, carries the shares in force on which it is
// the held story (S-0338). A story with tasks carries its task plan in
// counts (S-0176). Acceptance archives a story the moment it merges
// (S-0087); pendingPublish names the stories a release has not yet covered,
// so an accepted, archived story stays on its column until it is published,
// instead of vanishing the instant it is accepted, before anyone has had the
// chance to see it there; the view lists them as unpublished too.
func NewBoardView(items []*Item, board *Board, now time.Time, all bool, pendingPublish map[string]bool, holds *Holds) BoardView {
	w := board.Flow()
	v := BoardView{Columns: map[string][]BoardCard{}, WIPLimits: board.WIPLimits, Order: board.Order, Counts: map[string]int{}, flow: w}
	for _, s := range w.States {
		lane := Lane{Name: s.Name, Title: s.Title, Kind: s.Kind, TakesLimit: s.Limit, Ordered: s.Order}
		if s.Limit {
			lane.WIPLimit = board.WIPLimits[s.Name]
		}
		v.Lanes = append(v.Lanes, lane)
	}
	source := pullSource(w)
	if holds == nil {
		holds = NewHolds(items, nil)
	}
	for id, ok := range pendingPublish {
		if ok {
			v.Unpublished = append(v.Unpublished, id)
		}
	}
	sort.Strings(v.Unpublished)
	// A parent is looked up among every item given, archived ones included.
	titles := map[string]string{}
	tasks := map[string][]*Item{} // by story, for its plan
	for _, it := range items {
		titles[it.ID] = it.Title
		if it.Type == Task {
			tasks[it.Parent] = append(tasks[it.Parent], it)
		}
	}
	for _, it := range items {
		kind := w.Kind(it.Status)
		archivedButPending := it.Archived && kind == workflow.KindAccepted && pendingPublish[it.ID]
		if (it.Archived && !archivedButPending) || (!all && it.Type != Story) {
			continue
		}
		age := now.Sub(it.EnteredAt())
		card := BoardCard{
			ID: it.ID, Type: it.Type, Title: it.Title, Nature: it.Nature, Parent: it.Parent,
			ParentTitle: titles[it.Parent], Status: it.Status, EnteredAt: it.EnteredAt().UTC().Format(TimeFormat),
			Blocked: it.IsBlocked(), Age: HumanDuration(age), AgeSecs: int64(age.Seconds()),
			Touches: it.Touches, Archived: it.Archived, Draft: it.Type == Story && it.Draft,
		}
		if it.Type == Story && it.Status == source && !it.Archived {
			card.Held = holds.Of(it)
		}
		if it.Type == Story && !it.Archived && (it.Status == source || underwayKind(kind)) {
			card.Shared = holds.SharesWith(it.ID)
		}
		if it.Type == Story {
			if p := PlanOf(tasks[it.ID], it.ID); p != nil {
				card.Tasks = p.Summary()
			}
		}
		v.Columns[it.Status] = append(v.Columns[it.Status], card)
		if it.Type == Story {
			v.Counts[it.Status]++
		}
	}
	// The ordered states read in pull order: the stories take the sequence, in
	// the places stories already hold, and epics and tasks stay where they are.
	for _, st := range orderedStates(w) {
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
	for _, st := range stateNames(limitedStates(w)) {
		if limit, ok := board.WIPLimits[st]; ok && limit > 0 && v.Counts[st] > limit {
			v.Breaches = append(v.Breaches, fmt.Sprintf("%s has %d stories, limit %d", st, v.Counts[st], limit))
		}
	}
	return v
}

// ReadyInPullOrder returns the stories in the pull's source, those named in
// the pull order first and in that order, the rest after them by ID. The
// view's columns are already in that sequence.
func (v BoardView) ReadyInPullOrder() []BoardCard {
	var out []BoardCard
	for _, c := range v.Columns[pullSource(v.boardFlow())] {
		if c.Type == Story {
			out = append(out, c)
		}
	}
	return out
}

// FirstPullable is the first story in the pull's source, in pull order,
// that is not held, or nil when every one is (S-0128).
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

// CanPull reports whether a story may be pulled now: every active state's
// limit leaves room for one more, and every waiting state is under its limit
// (S-0243, ADR-0301 §2).
func (v BoardView) CanPull() bool { return v.PullHold() == "" }

// PullHold says why no story may be pulled now, or "" when one may: an
// active state's limit is full, or a waiting state is full and waits on the
// operator's acceptance (S-0243, I-0007, ADR-0301 §2), a full active state
// first, the first in board order. Stories already started still finish;
// only the pull stops.
func (v BoardView) PullHold() string {
	if st, n, limit, full := v.fullOfKind(workflow.KindActive); full {
		return fmt.Sprintf("the %s limit is full (%d of %d)", st, n, limit)
	}
	return v.ReviewHold()
}

// ReviewHold says that a waiting state is at or over its limit, so that no
// story is pulled until one is accepted or sent back, or "" when every
// waiting state has room or no limit (S-0243). flai serve, which counts the
// agents it started against the active limit itself, asks it on its own.
func (v BoardView) ReviewHold() string {
	if st, n, limit, full := v.fullOfKind(workflow.KindWaiting); full {
		return fmt.Sprintf("%s is full (%d of %d): accept or send back a story", st, n, limit)
	}
	return ""
}

// fullOfKind is the first state of kind, in board order, that takes a limit
// and holds at least as many stories as it allows: its name, its count, and
// its limit.
func (v BoardView) fullOfKind(kind string) (state string, n, limit int, full bool) {
	for _, s := range limitedStates(v.boardFlow()) {
		if s.Kind != kind {
			continue
		}
		if l, ok := v.WIPLimits[s.Name]; ok && l > 0 && v.Counts[s.Name] >= l {
			return s.Name, v.Counts[s.Name], l, true
		}
	}
	return "", 0, 0, false
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
