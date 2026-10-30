package workitem

import (
	"fmt"
	"strings"

	"github.com/bytepunx/system-flow/flai/internal/manifest"
)

// OrchestratorPermits says why the orchestrator may not do what, on item, or
// nil when it may: the project's orchestration.permissions gives it
// permission. flai holds the orchestrator to its permissions itself, as flai
// guard does before the call, so that a call the guard does not see is held
// all the same (S-0219).
func (r *Repo) OrchestratorPermits(permission, what, item string) error {
	if r.Manifest.Orchestration.Permissions.Allows(permission) {
		return nil
	}
	return fmt.Errorf("the orchestrator %s only with orchestration.permissions.%s, which is off: ask the operator with thread_open on %s", what, permission, item)
}

// Promotable says why the orchestrator may not move it to ready, or nil when
// it may (S-0219): it is a story among the promotion candidates, as flai
// promote --candidates lists them, and the ready column is under its WIP
// limit, a limit of 0 or none being no limit, while the project gives it
// promote_to_ready. A story that is not a candidate is refused with every
// reason the candidates give for it.
func (r *Repo) Promotable(it *Item) error {
	if err := r.OrchestratorPermits(manifest.PermitPromoteToReady, "moves a story to ready", it.ID); err != nil {
		return err
	}
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

// OrchestratorSetsCostOfDelay says why the orchestrator may not give it cost
// of delay inputs with an edit, or nil when it may (S-0328, ADR-0119,
// I-0187): the project gives it plan_backlog_stories, it is an epic or a
// story that is open, done and cancelled left out, and the edit gives inputs
// alone. besides names, as the caller words them, what else the edit gives: a
// value, a removal, or another field; any refuses it. The orchestrator sets
// the inputs the operator decided and no value: a value stays the
// operator's.
func (r *Repo) OrchestratorSetsCostOfDelay(it *Item, besides []string) error {
	if err := r.OrchestratorPermits(manifest.PermitPlanBacklogStories, "sets cost of delay inputs", it.ID); err != nil {
		return err
	}
	if len(besides) > 0 {
		return fmt.Errorf("%s: the orchestrator sets cost of delay inputs and nothing else, and this edit gives %s too: it sets no value, removes nothing, and changes nothing else of an item", it.ID, strings.Join(besides, ", "))
	}
	switch {
	case it.Type != Epic && it.Type != Story:
		return fmt.Errorf("%s is %s: the orchestrator sets cost of delay inputs on an epic or a story, and no other item", it.ID, articled(it.Type))
	case it.Archived:
		return fmt.Errorf("%s is archived: the orchestrator sets cost of delay inputs on an open epic or story alone", it.ID)
	case it.Closed():
		return fmt.Errorf("%s is %s: the orchestrator sets cost of delay inputs on an open epic or story alone, and none on one done or cancelled", it.ID, it.Status)
	}
	return nil
}
