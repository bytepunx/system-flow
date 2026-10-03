---
id: ADR-0074
title: "Work items carry planning data: a story's draft flag, an epic's or story's cost of delay, and a story's forecast, each saying who set it"
status: accepted
date: 2026-10-03
supersedes: []
superseded_by: []
---

# ADR-0074 Work items carry planning data: a story's draft flag, an epic's or story's cost of delay, and a story's forecast, each saying who set it

## Context

E-0016's planner, orchestrator, and analyzer need data the work items did not carry. They need to know whether a story is a draft an agent wrote, what delaying an item costs, and how long a story will take and when it will land. Nothing recorded who set a value, so a human's figure could not be told from a planner's. `estimate`, a Go duration, existed, and nothing wrote or read it. The designer decided the fields on 2026-10-02 for E-0016 (S-0199).

## Decision

Stories carry an optional `draft` flag and `forecast`, epics and stories an optional `cost_of_delay`, each block saying who set it and when, all set only through flai.

- `draft: true`, on a story an agent wrote (the planner, or the step that makes a story from an issue). It is set with `flai story new --draft`, MCP `item_new` `draft`, and `flai edit --draft` on a story in the backlog. `flai edit --no-draft` clears it, and so does `flai move <story> ready --yes`, which finalizes the story as it moves. A move of a draft story to ready is refused with "finalize it first", in `workitem.Move`, unless the move finalizes (`MoveOptions.Finalize`). The operator finalizes with `--yes`, and the dashboard with `item.move` `finalize`. No agent may finalize: MCP `item_move` refuses a draft to ready, and `item_edit` refuses `draft: false`, until the orchestrator's permission to finalize exists (S-0218). `flai check` warns (`story.draft`) on a draft story in ready or later.
- `cost_of_delay`, on epics and stories:
  - `inputs` the operator supplies, each optional: `revenue_per_week` and `penalty_per_week`, amounts in the project's currency, and `time_lost_per_cycle`, a Go duration.
  - `value`, the cost of delay per week that the planner derives.
  - `by` and `at`: who last changed the block, and when.

  A story with no inputs of its own may carry only a `value`, the share of its epic's value the planner apportions it.
- `forecast`, on stories: `duration` (agent wall clock, a Go duration), `delivery` (a UTC timestamp), `basis` (what it was computed from, in one sentence), `by`, and `at`. `estimate` stays the human's figure, and `forecast` is the planner's: the metrics compare each with the actual (S-0205).

Amounts are plain non-negative numbers in `planning.currency`, a manifest key with an ISO 4217 code (default USD). The manifest also sets `planning.hour_rate`, what an hour of work costs, which prices time lost (unset means unknown), and `planning.cycle`, the period time lost is counted over (a Go duration, default 168h).

Every field is set through flai, never by a command that owns state: `flai edit`, MCP `item_edit`, and hostapi `item.edit`, which runs `flai edit`. Each edit changes only the keys it names. An empty value removes a key, and `--clear-cost-of-delay` or `--clear-forecast` removes the block. A block that changes is stamped with the editor (`--by`, else `FLAI_AGENT`, else the config author) and the time. `Validate`, and through it `flai check` (`item.front-matter`), refuses a field on a type that does not carry it, a negative or non-finite amount, a duration that is not a Go duration longer than zero, a timestamp that is not UTC, a block without `by` and `at`, a cost of delay with neither inputs nor value, and a forecast with neither duration nor delivery. `flai check` reports a bad `planning` block as `manifest.planning`.

The fields and their keys are listed in `front-matter-fields.txt`, so publishing the flai release that carries them raises `flai.minimum` (S-0181). An older flai stops before reading the project.

## Consequences

- The planner, orchestrator, and analyzer have a place to write what they work out, and a reader can see who wrote it.
- A draft story cannot reach an agent before the operator has read it.
- One `by` and `at` per block records the last editor, not each key's. The keys say whose a figure is: the inputs are the operator's, and the value is the planner's.
- A JSON number cannot be empty, so MCP removes one amount by clearing the block and giving the keys to keep.
- The host's flai must be upgraded once the release carrying the fields is published.

## Alternatives considered

- **`by` and `at` on each key.** It is exact, but the front matter doubles in size, and the keys already say whose a figure is.
- **Amounts written with their currency, such as `1200 EUR`.** Mixed currencies could not be summed, and one project has one currency.
- **Reusing `estimate` for the planner's forecast.** The human's figure would be lost, and estimate against forecast is one of the comparisons the metrics want.
- **A draft as a state before backlog.** The state machine is shared by every type and drives the metrics. A flag on stories changes neither.
