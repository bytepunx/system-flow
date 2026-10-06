package workitem

import (
	"fmt"
	"strings"
)

// Promotable says why the orchestrator may not move it to ready, or nil when
// it may (S-0219): it is a story among the promotion candidates, as flai
// promote --candidates lists them, and the ready column is under its WIP
// limit, a limit of 0 or none being no limit. A story that is not a
// candidate is refused with every reason the candidates give for it.
func (r *Repo) Promotable(it *Item) error {
	if it.Type != Story {
		return fmt.Errorf("%s is %s: the orchestrator moves a story to ready, a candidate of flai promote --candidates, and no other item", it.ID, articled(it.Type))
	}
	items, err := r.List(false)
	if err != nil {
		return err
	}
	board, err := r.LoadBoard()
	if err != nil {
		return err
	}
	got, err := PromotionCandidates(items, board.Order, r.Holds(items), r.Manifest.Orchestration.PolicyOrDefault(), 0)
	if err != nil {
		return err
	}
	candidate := false
	for _, c := range got.Candidates {
		candidate = candidate || c.ID == it.ID
	}
	if !candidate {
		for _, o := range got.Others {
			if o.ID == it.ID {
				return fmt.Errorf("%s is not a candidate to go to ready (flai promote --candidates): %s", it.ID, strings.Join(o.Reasons, "; "))
			}
		}
		return fmt.Errorf("%s is %s, not in the backlog: the orchestrator moves a candidate of flai promote --candidates to ready, and no other story", it.ID, it.Status)
	}
	if limit := board.WIPLimits[Ready]; limit > 0 {
		n := 0
		for _, x := range items {
			if x.Type == Story && x.Status == Ready && !x.Archived {
				n++
			}
		}
		if n >= limit {
			return fmt.Errorf("the ready column is at its WIP limit (%d of %d): the orchestrator moves no story to ready until one leaves it", n, limit)
		}
	}
	return nil
}
