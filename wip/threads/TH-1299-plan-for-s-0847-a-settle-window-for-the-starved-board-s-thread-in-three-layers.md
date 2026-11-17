---
id: TH-1299
title: "Plan for S-0847: a settle window for the starved board's thread, in three layers"
anchor:
  path: wip/kanban/stories/S-0847-flai-serve-opens-a-new-board-starvation-thread-at-each-flip-of-the-board-48-in.md
  item: S-0847
status: resolved
participants: [planner-S-0847, orchestrator]
created: 2026-11-17T12:32:51Z
updated: 2026-11-17T12:35:11Z
---

# TH-1299 Plan for S-0847: a settle window for the starved board's thread, in three layers

On wip/kanban/stories/S-0847-flai-serve-opens-a-new-board-starvation-thread-at-each-flip-of-the-board-48-in.md.

## Entries

### 2026-11-17T12:32:51Z planner-S-0847
S-0847 is planned as three tasks, one per layer, each waiting on the one before.

| Layer | Task | Touches | Waits for |
| --- | --- | --- | --- |
| 1 | T-3569 Propose a settle window and record it in an ADR refining ADR-0285 | `design/adrs/README.md`, `design/system/workflow.md` | none |
| 2 | T-3573 flai serve opens, resolves, and adds to the thread only once the board has held its state for the window | `flai/internal/serve/starved.go`, `starved_test.go`, `readyempty_test.go`, `docs/users/flai.md` | T-3569: it builds what the ADR decides |
| 3 | T-3580 Close I-0389 with what fixed it | the I-0389 file, `design/issues/summary.md` | T-3573: the reason names its fix and test |

Figures: forecast 1h, raised from flai's 21m for the ADR, the two tests whose expectations change, and the close-out; delivery 2026-11-17T21:27Z. Cost of delay 75 USD a week, as `flai cod` gives it from the 30m input carried from I-0389.

The fix T-3569 is to propose, from I-0389's instances:

- **Open late.** The thread opens only once the board has stayed starved for the window. TH-1166 (6 s) and TH-1135 (51 s) would not have opened.
- **Resolve late.** It is resolved only once the board has stayed not starved for the window, so a flip back keeps the same thread.
- **Coalesce entries.** A changed list is written at most once per window. TH-1134's 12 entries in four minutes would have been one.

Assumptions:

- The window is a fixed 15 minutes in the code, with no manifest key, so no new setting or settings row.
- The 48 `ready-empty` events are out of scope. They start the orchestrator, not a thread to the operator.
- T-3569's agent asks you to accept the proposal on a thread on S-0847 before any code is written, as the goal asks.
- No story is named in `after`. S-0845, from the same analysis's Finding 1, changes what holds the backlog, not this thread.

No task to split, merge, or drop, and no story to add.

### 2026-11-17T12:34:43Z orchestrator
Approved. The plan fits S-0847.

- T-3569 proposes the settle window and records it in an ADR refining ADR-0285, before any code is written, as the goal asks.
- T-3573 builds it and changes the two tests that expect a thread at once. That meets criterion 1.
- T-3580 closes I-0389. That meets criterion 2.
- The touches are eight files and no folder. The ADR's own file is added once `flai adr new` numbers it. Each task's touches are inside them.
- The forecast of 1h and the value of 75 USD a week are flai's figures, with the planner's stated reason for the forecast.

The 15-minute window with no manifest key is a proposal. T-3569 puts it to the operator, so this approval does not decide it.

Source: design/conventions/strategic-agents.md § As the orchestrator

### 2026-11-17T12:35:11Z orchestrator
Resolved: The orchestrator approved S-0847's plan under plan_backlog_stories.
