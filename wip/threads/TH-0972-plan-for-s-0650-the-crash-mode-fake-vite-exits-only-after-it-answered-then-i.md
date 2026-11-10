---
id: TH-0972
title: "Plan for S-0650: the crash-mode fake vite exits only after it answered, then I-0325 closes"
anchor:
  path: wip/kanban/stories/S-0650-testcomponentsserveendswhenviteexits-fails-under-load-the-fake-vite-exits.md
  item: S-0650
status: open
participants: [planner-S-0650]
created: 2026-11-10T06:59:02Z
updated: 2026-11-10T06:59:02Z
---

# TH-0972 Plan for S-0650: the crash-mode fake vite exits only after it answered, then I-0325 closes

On wip/kanban/stories/S-0650-testcomponentsserveendswhenviteexits-fails-under-load-the-fake-vite-exits.md.

## Entries

### 2026-11-10T06:59:02Z planner-S-0650
Plan for S-0650, two tasks in two layers.

Cause: in `crash` mode the fake vite in `flai/cmd/components_test.go` exits on a 300 ms timer, whether or not flai's readiness poll (every 150 ms) has had an answer. Under load no poll lands in time, so `StartDevServer` correctly reports "exited before it answered", and the test, which wants "exited: exit status 2", fails. The product code is right. Only the test's fake is wrong.

| Layer | Task | Touches | after |
|-------|------|---------|-------|
| 1 | T-2836 The crash-mode fake vite exits only after it has answered, and a slow-starting fake reproduces the early exit | `flai/cmd/components_test.go` | none |
| 2 | T-2837 Close I-0325 with what fixed it | I-0325's file, `design/issues/summary.md` | T-2836 |

Figures: forecast 35m (flai forecast's 5m raised for the integration and smoke close-out a `flai/` change runs), delivery 2026-11-10T08:34Z; cost of delay value 20 USD a week, as `flai cod` gives it from flai's 8m input.

Assumptions:

- The fix is test-only. `flai/cmd/components.go` and `flai/internal/components/devserver.go` stay as they are. The exit-before-answer path is already covered by `TestStartDevServerReturnsTheLastLinesOfAViteThatExitsEarly`.
- A start delay before the fake listens reproduces I-0325 without load under the old timer, and the delay stays in the test.
- No `after` on S-0647 (I-0311, the same kind of race): it changes a different file and this story reads none of its output.

No earlier plan or operator entry on S-0650 was found. This is its first plan.
