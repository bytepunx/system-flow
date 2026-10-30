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
// of delay inputs with an edit, or nil when it may. With set_cost_of_delay
// (S-0523, ADR-0224) it adds or changes the inputs of an epic or a story that
// is open, done and cancelled left out, whatever inputs it or its epic has.
// With plan_backlog_stories alone (S-0328, ADR-0119) it gives a story in the
// backlog its first inputs, while neither it nor its epic has any. Either way
// the edit gives inputs alone: besides names, as the caller words them, what
// else the edit gives, a value, a removal, or another field, and any refuses
// it. A value stays the operator's.
func (r *Repo) OrchestratorSetsCostOfDelay(it *Item, besides []string) error {
	p := r.Manifest.Orchestration.Permissions
	wide, narrow := p.Allows(manifest.PermitSetCostOfDelay), p.Allows(manifest.PermitPlanBacklogStories)
	if !wide && !narrow {
		return fmt.Errorf("the orchestrator sets cost of delay inputs only with orchestration.permissions.%s, or on a backlog story's first inputs with orchestration.permissions.%s, which are off: ask the operator with thread_open on %s", manifest.PermitSetCostOfDelay, manifest.PermitPlanBacklogStories, it.ID)
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
	case wide:
		return nil
	}
	// the rest needs set_cost_of_delay, which is off
	rest := fmt.Sprintf("with %s the orchestrator gives a backlog story its first cost of delay inputs alone; the rest needs orchestration.permissions.%s, which is off", manifest.PermitPlanBacklogStories, manifest.PermitSetCostOfDelay)
	switch {
	case it.Type == Epic:
		return fmt.Errorf("%s is an epic: %s", it.ID, rest)
	case it.Status != Backlog:
		return fmt.Errorf("%s is %s, not in the backlog: %s", it.ID, it.Status, rest)
	case hasCostInputs(it):
		return fmt.Errorf("%s has cost of delay inputs already: %s", it.ID, rest)
	}
	epic, err := r.parentOf(it, nil)
	if err != nil {
		return err
	}
	if epic != nil && hasCostInputs(epic) {
		return fmt.Errorf("%s's epic %s has cost of delay inputs: %s", it.ID, epic.ID, rest)
	}
	return nil
}

// hasCostInputs says whether it has a cost of delay input.
func hasCostInputs(it *Item) bool {
	return it.CostOfDelay != nil && !it.CostOfDelay.Inputs.IsZero()
}
