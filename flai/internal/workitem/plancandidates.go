package workitem

import (
	"fmt"
	"sort"
	"strings"
)

// Planning candidates (S-0219). An epic is a candidate for the planner when
// it is in the backlog with no story, for the planner to draft its stories,
// or when it is not done or cancelled and every story under it is done or
// cancelled with at least one done, for the planner to draft what its
// outcome still lacks. Stories in the archive count. An epic is not a draft
// (only a story can be), so no epic is left out for being one. A story is a
// candidate when it is in the backlog, not archived, and lacks a plan: it
// has no touches, no forecast duration, no cost of delay value, or no task
// that is not cancelled, archived tasks counting (S-0328, ADR-0119); a draft
// is planned as a finalized story is. The caller names the items to leave
// out, with why: those a planner runs for now, and those whose planner asked
// a question the operator has not answered, which flai serve's runs record
// and the work items do not.

// PlanCandidates are the epics and the stories the planner should plan, and
// those that would be but are left out, with why.
type PlanCandidates struct {
	Candidates []PlanCandidate `json:"candidates"`
	LeftOut    []PlanCandidate `json:"left_out"`
}

// PlanCandidate is an epic or a story with the reason it is a candidate, or
// the reason it is left out.
type PlanCandidate struct {
	ID string `json:"id"`
	// Type is epic or story (S-0328).
	Type   string `json:"type"`
	Title  string `json:"title"`
	Status string `json:"status"`
	Reason string `json:"reason"`
}

// PlanCandidatesOf judges the epics among items, archived stories included,
// in ID order, then the stories, archived tasks included, in ID order. An
// item in leave, by ID, is listed as left out with the reason leave gives,
// when it would otherwise be a candidate. Archived epics and stories are
// never judged.
func PlanCandidatesOf(items []*Item, leave map[string]string) PlanCandidates {
	out := PlanCandidates{Candidates: []PlanCandidate{}, LeftOut: []PlanCandidate{}}
	type tally struct{ stories, done, cancelled int }
	under := map[string]*tally{}
	tasked := map[string]bool{}
	for _, it := range items {
		if it.Type == Task && it.Status != Cancelled {
			tasked[it.Parent] = true
		}
		if it.Type != Story {
			continue
		}
		n := under[it.Parent]
		if n == nil {
			n = &tally{}
			under[it.Parent] = n
		}
		n.stories++
		switch it.Status {
		case Done:
			n.done++
		case Cancelled:
			n.cancelled++
		}
	}
	var epics, stories []*Item
	for _, it := range items {
		switch {
		case it.Archived:
		case it.Type == Epic:
			epics = append(epics, it)
		case it.Type == Story:
			stories = append(stories, it)
		}
	}
	sort.SliceStable(epics, func(i, j int) bool { return lessID(epics[i].ID, epics[j].ID) })
	sort.SliceStable(stories, func(i, j int) bool { return lessID(stories[i].ID, stories[j].ID) })
	judged := func(it *Item, reason string) {
		c := PlanCandidate{ID: it.ID, Type: it.Type, Title: it.Title, Status: it.Status, Reason: reason}
		if why, ok := leave[it.ID]; ok {
			c.Reason = why
			out.LeftOut = append(out.LeftOut, c)
			return
		}
		out.Candidates = append(out.Candidates, c)
	}
	for _, epic := range epics {
		n := under[epic.ID]
		if n == nil {
			n = &tally{}
		}
		var reason string
		switch {
		case epic.Status == Done || epic.Status == Cancelled:
			continue
		case n.stories == 0 && epic.Status == Backlog:
			reason = "in the backlog with no stories: the planner drafts them"
		case n.stories > 0 && n.done > 0 && n.done+n.cancelled == n.stories:
			reason = fmt.Sprintf("every story done or cancelled (%d done, %d cancelled) while the epic is %s: the planner drafts what its outcome still lacks", n.done, n.cancelled, epic.Status)
		default:
			continue
		}
		judged(epic, reason)
	}
	for _, story := range stories {
		if story.Status != Backlog {
			continue
		}
		if lacks := storyPlanLacks(story, tasked[story.ID]); len(lacks) > 0 {
			judged(story, "in the backlog without a plan: "+strings.Join(lacks, "; "))
		}
	}
	return out
}

// storyPlanLacks is each part of a plan story lacks (S-0328): touches, a
// forecast duration, a cost of delay value, as flai promote --drafts judges
// them, and a task that is not cancelled, which tasked says it has. None
// means it has a plan.
func storyPlanLacks(story *Item, tasked bool) []string {
	var lacks []string
	if len(story.Touches) == 0 {
		lacks = append(lacks, "no touches")
	}
	if _, ok := forecastHours(story); !ok {
		lacks = append(lacks, "no forecast duration")
	}
	if codValue(story) == nil {
		lacks = append(lacks, "no cost of delay value")
	}
	if !tasked {
		lacks = append(lacks, "no tasks")
	}
	return lacks
}
