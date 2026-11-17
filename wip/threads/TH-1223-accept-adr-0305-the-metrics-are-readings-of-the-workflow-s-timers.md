---
id: TH-1223
title: "Accept ADR-0305: the metrics are readings of the workflow's timers"
anchor:
  path: wip/kanban/stories/S-0750-flai-stats-computes-lead-cycle-queue-and-review-time-from-the-baseline-workflow.md
  item: S-0750
status: resolved
participants: [agent-S-0750, orchestrator, alex]
created: 2026-11-15T20:32:16Z
updated: 2026-11-17T07:47:22Z
---

# TH-1223 Accept ADR-0305: the metrics are readings of the workflow's timers

On wip/kanban/stories/S-0750-flai-stats-computes-lead-cycle-queue-and-review-time-from-the-baseline-workflow.md.

## Entries

### 2026-11-15T20:32:16Z agent-S-0750
**Recommendation: accept ADR-0305 as written** (`design/adrs/0305-lead-cycle-queue-and-review-time-and-time-in-state-are-readings-of-the-workflow.md`, on `story/S-0750`). Criterion 1 needs an accepted ADR, and `metrics.md` changes only with one.

On the baseline it changes no figure. Lead, cycle, queue, and review time become the values of the baseline's `lead`, `cycle`, `queue`, and `review` timers (ADR-0301 §5). Time in state becomes a sum per state of the workflow.

## What is new beyond ADR-0301 §5

Each rule below keeps a figure as `flai stats` gives it today.

1. **A closed item's last interval is not counted.** An item in an `accepted` or `cancelled` state stops accruing time. Today's time in state does this. Counting it would make time in `done` grow without end.
2. **A bounded sum has no value for an item that first entered its states after the `until` point.** Today an item that started before it was ever ready has no queue time. Giving it 0 would change that.
3. **`completed` reads kinds.** It is the last transition while the item is in an `accepted` or `cancelled` state. A reopened item is open, so it has no lead or cycle time and counts in no aggregate until it closes (ADR-0055).
4. **`counts: accepted` is what keeps cancelled items out of the aggregates.** `closed` lets them in.
5. **The JSON field names stay** (`lead_time_seconds` and the rest) on every workflow, and the `time_in_state` keys are the workflow's state names. S-0752 reports a workflow's own timers beside them.

## Alternative

Accept with changes. Name the rule to change: a change to 1 or 2 changes past figures, which criterion 4 forbids on the fixtures.

## What happens next

The code tasks (T-3237, T-3241, T-3246) go on meanwhile, since the byte-for-byte criterion fixes the figures either way. Once you answer, I accept the ADR with `flai adr accept` and write `metrics.md` (T-3242).

### 2026-11-15T20:33:28Z orchestrator (recommendation)
**Recommended:** accept ADR-0305 as written. agent-S-0750 then accepts it with `flai adr accept` and writes `design/system/metrics.md` (T-3242).

**Why:**

- ADR-0301 §5 asks for this ADR. It says S-0750 "writes the ADR that changes the metrics contract". It also fixes the timers' forms, points, and the baseline's `lead`, `cycle`, `queue`, and `review`. ADR-0305 §2, §3, and §5 restate these unchanged.
- The five rules agent-S-0750 lists are new, but each one keeps a figure `flai stats` gives today. Criterion 4 asks for exactly that: the same bytes on `testdata/good` and `windows.golden.json`.
- Rule 1 (no last interval for a closed item) is how time in state works today. Without it, time in `done` grows without end.
- Rule 2 (a bounded sum has no value when the item entered its states only after the point) keeps ADR-0190's queue time. That rule gives no value to an item that started before it was ready.
- Rules 3 and 4 restate ADR-0055 and ADR-0301 §5's `counts: accepted` in terms of kinds.
- Rule 5 keeps the JSON field names, so flaiover needs no change for the baseline.

**Why you:** accepting an ADR records a decision, and `metrics.md` changes only with one. That call is yours.

**Alternative:** accept with changes, naming the rule. A change to rule 1 or 2 changes past figures, which criterion 4 forbids on the fixtures.

Source: design/adrs/0301-the-board-is-a-workflow-file-design-workflows-board-md-whose-front-matter-names.md § 5. Timers

### 2026-11-17T07:46:10Z alex
Confirmed the recommendation of 2026-11-15T20:33:28Z orchestrator.

Source: design/adrs/0301-the-board-is-a-workflow-file-design-workflows-board-md-whose-front-matter-names.md § 5. Timers

### 2026-11-17T07:47:22Z agent-S-0750
Resolved.
