package workitem

import (
	"fmt"
	"sort"
)

// Planning candidates (S-0219). An epic is a candidate for the planner when
// it is in the backlog with no story, for the planner to draft its stories,
// or when it is not done or cancelled and every story under it is done or
// cancelled with at least one done, for the planner to draft what its
// outcome still lacks. Stories in the archive count. An epic is not a draft
// (only a story can be), so no epic is left out for being one. The caller
// names the epics to leave out, with why: those a planner runs for now, and
// those whose planner asked a question the operator has not answered, which
// flai serve's runs record and the work items do not.

// PlanCandidates are the epics the planner should plan, and those that would
// be but are left out, with why.
type PlanCandidates struct {
	Candidates []PlanCandidate `json:"candidates"`
	LeftOut    []PlanCandidate `json:"left_out"`
}

// PlanCandidate is an epic with the reason it is a candidate, or the reason
// it is left out.
type PlanCandidate struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Status string `json:"status"`
	Reason string `json:"reason"`
}

// PlanCandidatesOf judges the epics among items, archived stories included,
// in ID order. An epic in leave, by ID, is listed as left out with the
// reason leave gives, when it would otherwise be a candidate. Archived
// epics are never judged.
func PlanCandidatesOf(items []*Item, leave map[string]string) PlanCandidates {
	out := PlanCandidates{Candidates: []PlanCandidate{}, LeftOut: []PlanCandidate{}}
	type tally struct{ stories, done, cancelled int }
	under := map[string]*tally{}
	for _, it := range items {
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
	var epics []*Item
	for _, it := range items {
		if it.Type == Epic && !it.Archived {
			epics = append(epics, it)
		}
	}
	sort.SliceStable(epics, func(i, j int) bool { return lessID(epics[i].ID, epics[j].ID) })
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
		c := PlanCandidate{ID: epic.ID, Title: epic.Title, Status: epic.Status, Reason: reason}
		if why, ok := leave[epic.ID]; ok {
			c.Reason = why
			out.LeftOut = append(out.LeftOut, c)
			continue
		}
		out.Candidates = append(out.Candidates, c)
	}
	return out
}

// PlanCandidates reads every item, the archive's included, and judges the
// epics, leaving out those in leave.
func (r *Repo) PlanCandidates(leave map[string]string) (PlanCandidates, error) {
	items, err := r.List(true)
	if err != nil {
		return PlanCandidates{}, err
	}
	return PlanCandidatesOf(items, leave), nil
}
