---
id: TH-0694
title: "Plan for S-0523: widen plan_backlog_stories so the orchestrator sets cost of delay inputs on an epic or an open story, five tasks in four layers, needs a cost of delay input (recommend time lost 30m per cycle)"
anchor:
  path: wip/kanban/stories/S-0523-the-orchestrator-cannot-set-a-cost-of-delay-input-the-operator-decided-on-an-epic-or-a-story-past-the-backlog.md
  item: S-0523
status: open
participants: [planner-S-0523, orchestrator]
created: 2026-10-28T14:34:12Z
updated: 2026-10-28T14:34:37Z
---

# TH-0694 Plan for S-0523: widen plan_backlog_stories so the orchestrator sets cost of delay inputs on an epic or an open story, five tasks in four layers, needs a cost of delay input (recommend time lost 30m per cycle)

On wip/kanban/stories/S-0523-the-orchestrator-cannot-set-a-cost-of-delay-input-the-operator-decided-on-an-epic-or-a-story-past-the-backlog.md.

## Entries

### 2026-10-28T14:34:12Z planner-S-0523
S-0523 is planned: 27 touches, a 1h15m forecast, and five tasks in four layers. Two things need an answer, the cost of delay input first.

## 1. Cost of delay input (the operator's)

S-0523 has no cost of delay inputs and no epic, so `flai cod` refuses it a value.

- **Recommended:** `time_lost_per_cycle: 30m`, about 75 USD a week. Each input the operator decides on an epic, or on a story already started, goes back to them to run by hand (TH-0574, TH-0579), and while it waits the planners value stories by hand and `flai promote --drafts` lists the epic's drafts incomplete: E-0025's ten stories waited on one 4h input.
- Alternative: 15m, about 37.50 USD a week, if the hand edits are the whole cost.
- Alternative: 1h, about 150 USD a week, if the stalled finalizations count as well.

Under `plan_backlog_stories` the orchestrator may set it: S-0523 is a backlog story without inputs or epic. I write the value on my next run, or now if the input lands while I wait.

## 2. The solution proposed from I-0187's instances

Two gates stop the orchestrator. `flai guard` admits `item_edit` with `cost_of_delay` for a story ID alone (`setsCostInputs`, `setsCostInputsOnly` in `flai/internal/guard/guard.go`), so an epic is refused whatever the permissions. flai's own rule, `Repo.OrchestratorSetsCostOfDelay` in `flai/internal/workitem/promotable.go`, refuses a story not in the backlog, one that has inputs, and one whose epic has.

- **Recommended:** a new ADR refining ADR-0119. Under `plan_backlog_stories` the orchestrator may add or change the cost of delay inputs of an epic, or of a story in any state but done or cancelled, archived items left out; still inputs only, never a value, never a removal, nothing else in the call. The convention says when: a figure the operator decided on a thread (an answer or a confirmed recommendation), or the planner's recommendation on a story planner's thread for a backlog story without any, as today. No new permission: the operator has already handed over this decision by turning `plan_backlog_stories` on, and its settings text says the risk now reaches every open epic and story.
- Alternative: a new permission `set_cost_of_delay`, off by default, for the epic and started-story inputs alone; adds a manifest key, a settings toggle, and docs rows.
- Alternative: let it also set a value the operator names, as on TH-0574 ("Set COD to 190"); left out so the value stays what `flai cod` works out from the inputs.

The tasks are written for the recommended answer. If you choose an alternative, replan the story and I rewrite them.

## Tasks and layers

| Layer | Task | After | What |
|-------|------|-------|------|
| 1 | T-2241 | none | The ADR refining ADR-0119, with `flai adr new` |
| 2 | T-2242 | T-2241 | `Repo.OrchestratorSetsCostOfDelay` admits an epic and an open story, inputs only; the MCP and CLI tests reproduce the refusals first |
| 2 | T-2243 | T-2241 | `flai guard` admits an epic ID as it does a story's, in `item_edit` and `flai edit`; refusal texts reworded |
| 3 | T-2244 | T-2242, T-2243 | Design, convention and template copy, orchestrator prompt, settings text, generated reference, guides, template changelog |
| 4 | T-2245 | T-2244 | `flai issue close I-0187`, criteria ticked |

T-2242 and T-2243 run together; they share no path.

## Figures

- **Forecast** 1h15m, raised from flai's 44m: four serial layers, pinned tests rewritten in three packages, fourteen documents in T-2244. Delivery 2026-10-29T04:43Z, flai's 04:12Z moved by the 31m added.
- **Touches**: 27, none declared before; every one from the issue, the design, and the code. One folder touch kept, `design/adrs`, for the ADR whose number is allotted when written; S-0523's Planning note records it.

## Assumptions

- The dashboard's cost of delay write (`hostapi/writes.go`) runs `flai edit --by owner` without `FLAI_ROLE`, so it is not gated and needs no change.
- `.claude/agents/orchestrator.md` and its template copy say "giving the story cost of delay inputs", which stays true; they are left out of the touches because a write under `.claude/` needs your permission prompt. The story's agent may widen the touches and ask here if it judges the sentence worth changing.
- `docs/users/flai-reference.md` is regenerated by `make flai-reference` from `flai/internal/manifest/settings.go`, so T-2244 changes the Go text and regenerates it.
- The value stays the planner's: a story given inputs is listed again by `flai plan --candidates` until the value is written, as today.

### 2026-10-28T14:34:37Z orchestrator (recommendation)
**Recommended for the operator:** on question 2, take the alternative, a new permission `set_cost_of_delay`, off by default, rather than widening `plan_backlog_stories`. The planner then rewrites T-2242 to T-2244 for it.

## What I set (question 1)

S-0523's cost of delay input is now `time_lost_per_cycle: 30m`, your recommended figure. S-0523 is a backlog story with no inputs and no epic, which ADR-0119 lets me set. Nothing in the thread or the story argues for 15m or 1h. Write the value on your next run.

## The plan fits

- **Tasks:** five tasks in four layers. T-2241, the ADR, comes first. T-2242 (flai's rule) and T-2243 (the guard) share no path. T-2244 documents them, and T-2245 closes I-0187.
- **Figures:** 27 touches plus `design/adrs` for the ADR, forecast 1h15m over flai's 44m with the reason stated.

## Why question 2 is yours

This story widens what the orchestrator, me, may write. ADR-0119 § Decision limits my cost of delay writes to a backlog story. The plan extends them to every open epic and story, including ones whose inputs you already set. I should not approve an extension of my own powers. It is a decision not yet recorded, and it touches money.

## Why I recommend a permission of its own

- **Consent stays explicit.** You turned `plan_backlog_stories` on for backlog stories. A new switch lets you choose the wider reach on its own, and turn it off without losing story planning.
- **Changing your inputs is new.** Today flai refuses me inputs once a story or its epic has any, and the strategic-agents convention says "never overwrite the operator's inputs". The recommended answer would let me change inputs you set. A separate switch makes that a choice you make.
- **The cost is small:** a manifest key, a settings toggle, and docs rows, which T-2244 already touches.

**The planner's recommendation**, widening `plan_backlog_stories` with no new key, is simpler. If you prefer it, confirm the planner's option instead.

The thread stays open for your answer.

Source: design/adrs/0119-with-plan-backlog-stories-on-the-orchestrator-asks-the-planner-for-each-backlog.md § Decision
