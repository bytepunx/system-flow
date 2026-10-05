package workitem

import (
	"regexp"
	"strings"
)

// Promotion candidates (S-0217). A backlog story is a candidate to go to
// ready when it is not a draft, meets the definition of ready, would not be
// held if it were ready, and carries a forecast duration and a cost of delay
// value. The candidates are ordered by the project's policy; every other
// backlog story is listed with each reason it is not one.

// Candidates are the backlog stories that could go to ready, in the order a
// policy gives them, and the others with why not.
type Candidates struct {
	Policy     string    `json:"policy"`
	Candidates []Ranked  `json:"candidates"`
	Others     []Refused `json:"others"`
}

// Refused is a backlog story that is not a candidate, with every reason.
type Refused struct {
	ID      string   `json:"id"`
	Title   string   `json:"title"`
	Reasons []string `json:"reasons"`
	// Held is the hold it would have in ready, when it would have one.
	Held *Hold `json:"held,omitempty"`
}

// PromotionCandidates judges the backlog stories among items, taken in their
// pull order by order, and orders the candidates by policy. A limit above
// zero caps the candidates listed; those beyond it are not listed at all.
func PromotionCandidates(items []*Item, order []string, holds *Holds, policy string, limit int) (Candidates, error) {
	out := Candidates{Policy: policy, Candidates: []Ranked{}, Others: []Refused{}}
	byID := map[string]*Item{}
	for _, it := range items {
		if !it.Archived {
			byID[it.ID] = it
		}
	}
	var ok []*Item
	for _, id := range PullSequence(order, items, Backlog) {
		it := byID[id]
		reasons, held := promotionRefusals(it, byID, holds)
		if len(reasons) == 0 {
			ok = append(ok, it)
			continue
		}
		out.Others = append(out.Others, Refused{ID: it.ID, Title: it.Title, Reasons: reasons, Held: held})
	}
	ranked, err := OrderByPolicy(ok, policy)
	if err != nil {
		return Candidates{}, err
	}
	if limit > 0 && len(ranked) > limit {
		ranked = ranked[:limit]
	}
	out.Candidates = append(out.Candidates, ranked...)
	return out, nil
}

// PromotionCandidates reads the active items, the board's pull order, and
// the holds, and judges the backlog by orchestration.policy.
func (r *Repo) PromotionCandidates(limit int) (Candidates, error) {
	items, err := r.List(false)
	if err != nil {
		return Candidates{}, err
	}
	board, err := r.LoadBoard()
	if err != nil {
		return Candidates{}, err
	}
	return PromotionCandidates(items, board.Order, r.Holds(items), r.Manifest.Orchestration.PolicyOrDefault(), limit)
}

// promotionRefusals is every reason story is not a candidate, and the hold
// it would have in ready.
func promotionRefusals(story *Item, byID map[string]*Item, holds *Holds) ([]string, *Hold) {
	var reasons []string
	if story.Draft {
		reasons = append(reasons, "draft: finalize it first")
	}
	if !hasGoal(story.Body) {
		reasons = append(reasons, "no goal")
	}
	if !hasCriteria(story.Body) {
		reasons = append(reasons, "no acceptance criteria with a checkbox")
	}
	if p, ok := byID[story.Parent]; ok && p.Status == Cancelled {
		reasons = append(reasons, "its epic "+p.ID+" is cancelled")
	}
	held := holds.Of(story)
	if held != nil {
		reasons = append(reasons, held.Reason)
	}
	if _, ok := forecastHours(story); !ok {
		reasons = append(reasons, "no forecast duration")
	}
	if codValue(story) == nil {
		reasons = append(reasons, "no cost of delay value")
	}
	return reasons, held
}

var goalHeading = regexp.MustCompile(`(?m)^## Goal\s*$`)

// hasGoal says whether a body's "## Goal" section has any words.
func hasGoal(body string) bool {
	loc := goalHeading.FindStringIndex(body)
	if loc == nil {
		return false
	}
	rest := body[loc[1]:]
	if i := strings.Index(rest, "\n## "); i >= 0 {
		rest = rest[:i]
	}
	return strings.TrimSpace(rest) != ""
}
