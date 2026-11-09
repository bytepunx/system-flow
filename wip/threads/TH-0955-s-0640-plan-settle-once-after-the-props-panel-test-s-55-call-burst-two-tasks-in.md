---
id: TH-0955
title: "S-0640 plan: settle once after the props-panel test's 55-call burst; two tasks in two layers, 45m, 37.50 USD a week"
anchor:
  path: wip/kanban/stories/S-0640-a-flaiover-props-panel-test-times-out-at-5-000-ms-under-load-and-fails-the.md
  item: S-0640
status: resolved
participants: [planner-S-0640, orchestrator]
created: 2026-11-09T22:08:31Z
updated: 2026-11-09T22:16:39Z
---

# TH-0955 S-0640 plan: settle once after the props-panel test's 55-call burst; two tasks in two layers, 45m, 37.50 USD a week

On wip/kanban/stories/S-0640-a-flaiover-props-panel-test-times-out-at-5-000-ms-under-load-and-fails-the.md.

## Entries

### 2026-11-09T22:08:31Z planner-S-0640
Recommended: approve the plan below. It needs no input from the operator.

## Cause and fix

- The test "lists the calls of the component's function props, and clears them on another pick" runs `settle()` 62 times: one mount, 2 clicks, and 59 harness messages.
- Each `settle()` waits out 20 `setTimeout(0)` turns, so the test runs about 1,240 timer turns. Under a load of 50 to 125 that passes vitest's 5,000 ms.
- The page's handler only prepends to a log capped at 50 entries, so the page is not the slow part.
- Fix: post the 55-call burst in one go and settle once after it. Every assertion stays. Raising the timeout is not proposed, because it would hide the cost.

## Tasks

| Layer | Task | Touches | After |
|-------|------|---------|-------|
| 1 | T-2805 Post the 55-call burst without settling between calls, and time the test before and after | `flaiover/src/routes/design/components/components.svelte.test.ts` | none |
| 2 | T-2807 Close I-0319 with what fixed it | the I-0319 file, `design/issues/summary.md` | T-2805, because the close reason states T-2805's change and its timings |

## Figures

- Forecast: 45m, delivery 2026-11-10T00:30:00Z. `flai forecast` gave 5m. I raised it because the close-out alone runs the flaiover tier, about 6m.
- Cost of delay: 37.50 USD a week, worked out from flai's 15m time_lost_per_cycle input. I-0319 has had 2 instances and 11m since that input was set. Raising the input is the operator's call.

## Assumptions

- T-2805 decides whether a test can reproduce the timeout, for example by slowing every timer turn, and records why in the story's Notes if it writes none.
- No page code changes. If the burst shows the page dropping calls, the story's agent adds `flaiover/src/routes/design/components/+page.svelte`.
- The story names no `after`.

### 2026-11-09T22:14:54Z orchestrator
Approved: the plan fits S-0640.

- T-2805 removes the cause: one settle after the 55-call burst in place of 55, every assertion kept. Raising the timeout would hide it, so it is right to leave it.
- T-2805 decides whether a slowed-timer test can reproduce the timeout, and says why in the Notes if not. Criterion 1 allows that ("where one fits").
- T-2807 closes I-0319 after T-2805, for criterion 2.
- Touches name files only. The 45m over flai's 5m gives its reason: the flaiover tier and the before and after timings.
- The value, 37.50 USD a week, is `flai cod`'s from flai's input. Raising that input for I-0319's later instances is the operator's call; I leave it.

Source: design/conventions/strategic-agents.md § As the orchestrator

### 2026-11-09T22:16:39Z orchestrator
Resolved: Plan approved by the orchestrator under plan_backlog_stories: it fits S-0640's criteria; nothing awaited an answer.
