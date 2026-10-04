---
id: ADR-0084
title: "flai serve plans again on its own behind the plan host action: on an edit, when work ahead completes or the order changes as planning.replan allows, and on planning.schedule, coalescing runs and logging each one's trigger"
status: accepted
date: 2026-10-04
supersedes: []
superseded_by: []
refines: [ADR-0079, ADR-0082]
topics: [cli, dashboard, planning]
---

# ADR-0084 flai serve plans again on its own behind the plan host action: on an edit, when work ahead completes or the order changes as planning.replan allows, and on planning.schedule, coalescing runs and logging each one's trigger

## Context

ADR-0082 started the planner only on the operator's word: "Nothing starts it on a move or a timer." A forecast and a cost of delay value go stale, though. A story's words or the operator's inputs change; a story ahead is accepted or cancelled, or the pull order changes, so every delivery after it moves. S-0211 asks flai serve to plan again when it should, within what the operator allows, and to say in the planner's activity document what started each run.

## Decision

With the `plan` host action on, `flai serve` plans again on three triggers, each bounded by the manifest, and logs what triggered every planner run.

1. **The `plan` host action gates every trigger.** Off, flai serve plans nothing on its own and writes no forecast, as before. The operator's own asking (`flai plan`, the dashboard's Plan, the MCP tool `plan`) is unchanged.
2. **An edit.** When someone other than a planner edits a backlog or ready story's goal, criteria, or touches (an `edited` change), and the story has a forecast or a cost of delay value older than the edit, the planner is queued for it. A change to the cost of delay counts when it changed the inputs and the value is older than them; a value written alone does not. A story never planned, with neither figure, is not planned on an edit: replanning keeps planning fresh, and planning a story first is the operator's to ask.
3. **Work ahead completes or the order changes.** When a story is accepted or cancelled, or the pull order changes, `planning.replan` decides: `never` does nothing; `deterministic`, the default, plays the board out again for every ready and backlog story that has a forecast, as `flai forecast` does with no agent, keeping the story's own duration and writing the new delivery and basis where the delivery moved, all in one commit by `flai`; `agent` does the same and queues the planner for each story whose delivery moved.
4. **A schedule.** `planning.schedule`, a five-field cron expression in UTC or `daily` (00:00 UTC), queues the planner for every ready story each time it comes round. Unset, there is no schedule. flai parses the expression itself; `flai check` reports one it cannot parse as `manifest.planning`.
5. **Coalesced.** The queue holds a story once: a trigger for a story already queued is added to its entry, not queued again. flai serve runs the queue one planner at a time per project, through the same checks as `serve.Plan`, and drops an entry the checks refuse. A run the operator asks for is not queued.
6. **The trigger is logged.** The run records its trigger in `serve/agents.json`, and its activity entry in `wip/agents/planner.md` gains a `- Trigger:` line: `asked`, `edited <fields> by <who>`, `accepted <ID>`, `cancelled <ID>`, `reordered`, or `schedule <expression>`, separated by semicolons when coalesced. Entries without one, older entries and those an agent logs with `activity_log`, stay valid.
7. **The settings page shows the triggers**, read from the manifest and the host action: whether edits replan, the replan policy, and the schedule with its next run. They are set in `system-flow.yaml` by hand, as the other `planning` keys.

## Consequences

- With `plan` on, an operator who edits a planned story, accepts one, or reorders the board gets forecasts that follow, at no agent cost under the default `deterministic`.
- `agent` and a schedule spend money without the operator asking each time. Both are off until the operator sets them in the manifest.
- flai serve writes and commits items in the main checkout on its own for the first time: the forecasts of the deterministic replan, stamped `by: flai`.
- A flai older than this one refuses a `planner.md` holding a `- Trigger:` line, as its parse is strict; publishing this flai raises `flai.minimum` as for any change to what flai writes.
- flai serve remembers what it has seen from its start: an edit, a move, or a scheduled time that passed while it was down is not acted on.

## Alternatives considered

- **A planner run on every trigger, with no `replan` policy.** Every acceptance would start a planner for each story behind it; the deterministic replay moves the delivery at no cost, and the policy lets the operator pay for judgement when they want it.
- **Overwriting the forecast with what `flai forecast` prints.** It works the duration out from history, which would undo the duration the planner adjusted with a reason. Only the delivery and the basis move.
- **A cron library.** `robfig/cron` parses more than a five-field expression needs; a small parser has no dependency to keep.
- **The trigger in the summary line.** It would parse with an older flai, but the dashboard could not tell a trigger from what the planner said.
- **Remembering what was seen across restarts.** A schedule or a later edit catches a story up; state kept for it would be one more file to reconcile.
