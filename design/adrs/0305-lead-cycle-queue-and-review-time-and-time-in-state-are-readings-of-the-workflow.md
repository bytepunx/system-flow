---
id: ADR-0305
title: "Lead, cycle, queue, and review time and time in state are readings of the workflow's timers over an item's created and its transitions, and a closed item's last interval is not counted"
status: accepted
date: 2026-11-17
supersedes: []
superseded_by: []
refines: [ADR-0004, ADR-0055, ADR-0190, ADR-0301]
---

# ADR-0305 Lead, cycle, queue, and review time and time in state are readings of the workflow's timers over an item's created and its transitions, and a closed item's last interval is not counted

## Context

[metrics.md](../system/metrics.md) is the contract between `flai stats` and the dashboard, and changes only with an ADR. It defines lead, cycle, queue, and review time, and time in state, by state names: `created`, `committed` (the first move to `ready`), `started` (the first move to `in-progress`), `completed` (the last transition while `done` or `cancelled`, [ADR-0055](0055-a-story-moves-back-one-column-from-ready-in-progress-review-or-cancelled-and.md)), and the `state_intervals` between transitions ([ADR-0004](0004-workflow-states-and-transitions.md), [ADR-0190](0190-an-item-s-queue-time-is-the-seconds-it-spent-in-ready-before-it-started-so-a.md)).

[ADR-0301](0301-the-board-is-a-workflow-file-design-workflows-board-md-whose-front-matter-names.md) §5 makes a timer a reading of an item's transitions, gives its two forms and its points, and names the baseline's four timers, `lead`, `cycle`, `queue`, and `review`. It leaves the metrics contract as it was. This ADR changes the contract: the figures become readings of the workflow's timers, and this ADR says how a timer reads an item's history. On the baseline workflow every figure stays byte for byte as it is.

## Decision

**Lead, cycle, queue, and review time are the values of the workflow's timers `lead`, `cycle`, `queue`, and `review`, and time in state is a sum of the same kind for each state of the workflow. A timer reads an item's `created` and its recorded transitions as below.**

### 1. An item's history

An item's history is its `created`, then its transitions in order, each with the state it entered and its `at`. It spends one interval in each state it enters:

- from `created` to the first transition, in the workflow's initial state;
- from each transition to the next one, in the state the transition entered;
- from the last transition to now, in the state it entered, unless that state's kind is `accepted` or `cancelled`. A closed item's last interval is not counted. That is the open interval, the only one that grows with time.

### 2. Points

| Point | Is | None while |
|-------|----|------------|
| `created` | The item's `created`. | Never none. |
| `completed` | The `at` of the item's last transition, while the item's state is of kind `accepted` or `cancelled` ([ADR-0055](0055-a-story-moves-back-one-column-from-ready-in-progress-review-or-cancelled-and.md)). | The item is open, as is one reopened from a cancelled state. |
| `{first: <state>}`, `{last: <state>}` | The `at` of the first or last transition into the state. | The item never entered it. |
| `{first: {kind: <kind>}}`, `{last: {kind: <kind>}}` | The same, into any state of the kind. | The item never entered a state of the kind. |

### 3. The two forms

- **Between two points**, `from` and `to`: the time from `from` to `to`. It has no value while either point is none.
- **A sum**, `sum`: the length of every interval the item spent in a state the timer names, the open interval included (§1). It has no value for an item that never entered one of the states.
- **A bounded sum**, `sum` with `until`: only the intervals that end at or before the point. It has no value while the point is none, or when the item first entered one of the states after the point.

### 4. Which items an aggregate counts

A timer's per item value is given for every item its `types` names, open or closed. An aggregate over the window, its distribution, takes the items completed in the window that `counts` names: `accepted`, the default, those whose state is of kind `accepted`; `closed`, those of kind `cancelled` too. A cancelled item counts in the cancellation rate, never in an `accepted` timer's aggregate. A reopened item is open, so it has no value for a timer that ends at `completed` and counts in no aggregate until it closes again.

### 5. The figures

| Figure | Its timer on the baseline | Today's definition, which it equals |
|--------|---------------------------|-------------------------------------|
| Lead time | `lead`: `from: created`, `to: completed` | `completed - created` |
| Cycle time | `cycle`: `from: {first: in-progress}`, `to: completed` | `completed - started` |
| Queue time | `queue`: `sum: [ready]`, `until: {first: in-progress}` | The `ready` intervals that end at or before `started`, with a value only when `committed` and `started` exist and `started` is not before `committed` ([ADR-0190](0190-an-item-s-queue-time-is-the-seconds-it-spent-in-ready-before-it-started-so-a.md)) |
| Review time | `review`: `sum: [review]` | `wait_review_seconds`, the seconds in `review` with the open interval, for an item that ever entered `review` |
| Time in state | `sum: [<state>]`, for each state of the workflow, kept in `time_in_state` under the state's name | The sum of `state_intervals` per state |

Time in state and blocked time stay figures flai gives for every state and every item; they are not timers a workflow file names, as ADR-0301 §5 says. `committed` and `started` stay derived timestamps: the first transitions into the states the workflow's `commit` and `pull` transitions move to, `ready` and `in-progress` on the baseline.

### 6. Names in `flai stats`

`flai stats --json` keeps its field names on every workflow: `lead_time_seconds`, `cycle_time_seconds`, `queue_time_seconds`, `wait_review_seconds`, and the `time_in_state` keys, which are the workflow's state names. A workflow's timers of its own are reported beside these by S-0752.

## Consequences

- `flai stats` computes lead, cycle, queue, and review time by evaluating the workflow's timers, in `flai/internal/workflow`, over each item's history. On the baseline every figure, per item and aggregate, gives the same bytes as before.
- A state added between `in-progress` and `review` is inside cycle time, which runs from the first move to `in-progress` to `completed`. A state renamed keeps its intervals under its new name in time in state.
- A workflow that renames `ready`, `in-progress`, or `review` must name the new state in its `cycle`, `queue`, or `review` timer: `flai check` reports `workflow.timer` for a timer that names a state the file lacks.
- The rest of the metrics read a state's kind, not its name: work in progress is the items in an `active` or `waiting` state, completion is reaching an `accepted` state, and cancellation a `cancelled` one.
- The dashboard implements the same readings. It draws the timers flai gives, so it needs no change for the baseline.

## Alternatives considered

- **Keep the figures in code by state names, and the timers only as labels.** A workflow with a renamed or added state would then report figures its file does not define.
- **Count a closed item's last interval up to now.** Time in `done` would grow without end, and the time in state shares would change for every past item.
- **Give a bounded sum a value of zero when the item entered its states only after the point.** Queue time would then exist for an item that started before it was ever ready, where today it has none.
- **Make time in state a timer per state in the workflow file.** ADR-0301 §5 keeps it a figure flai gives for every state, so a file need not list one timer per state.
