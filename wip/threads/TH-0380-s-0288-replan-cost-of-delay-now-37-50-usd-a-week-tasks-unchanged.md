---
id: TH-0380
title: "S-0288 replan: cost of delay now 37.50 USD a week; tasks unchanged"
anchor:
  path: wip/kanban/stories/S-0288-flaiover-s-notify-test-ts-fails-now-and-then-under-the-full-vitest-run-because-project-info-reads-a-system-flow-yaml-with-no-version.md
  item: S-0288
status: resolved
participants: [planner-S-0288, orchestrator]
created: 2026-10-08T08:50:26Z
updated: 2026-10-08T08:50:32Z
---

# TH-0380 S-0288 replan: cost of delay now 37.50 USD a week; tasks unchanged

On wip/kanban/stories/S-0288-flaiover-s-notify-test-ts-fails-now-and-then-under-the-full-vitest-run-because-project-info-reads-a-system-flow-yaml-with-no-version.md.

## Entries

### 2026-10-08T08:50:26Z planner-S-0288
S-0288 was planned again after the cost of delay input went to 15m. The tasks are unchanged.

## Tasks

| Layer | Task | Touches | After |
|-------|------|---------|-------|
| 1 | T-1369 notify.test.ts reproduces the stale manifest read, and its setManifest writes atomically and tells the Repo | `flaiover/src/lib/server/notify.test.ts` | none |
| 2 | T-1370 Close I-0085 saying what fixed it | `design/issues/I-0085-…md`, `design/issues/summary.md` | T-1369 |

T-1370 waits for T-1369 because its close reason names the fix and the runs that showed it. Both tasks still fit the story. No task to split, merge, or drop.

## Figures

- Cost of delay: 37.50 USD a week, up from 12.50. This is `flai cod`'s figure for your 15m per 168h cycle at 150 USD an hour.
- Forecast: 20m, delivery 2026-10-08T09:10:00Z. flai now says 14m. I raised it for a deterministic reproduction and five full flaiover vitest runs.
- Touches: the three declared paths stay. The co-change list names other issues and docs that this test-only fix does not change, so none was added. There is no folder touch.

## Assumptions

1. The race is in the test's setup only, as TH-0378 said. `repo.ts` and `notify.ts` stay unchanged unless T-1369's reproduction shows otherwise.
2. The other flaiover tests that rewrite `system-flow.yaml` are out of scope. None has a recorded instance.

No answer is needed.

### 2026-10-08T08:50:31Z orchestrator
Replan approved. The two tasks still fit the story, the figures are as `flai cod` gives them for the operator's 15m input (37.50 USD/week), and the 20m forecast allows for the reproduction and the vitest runs. No answer was needed.

Source: wip/kanban/stories/S-0288-flaiover-s-notify-test-ts-fails-now-and-then-under-the-full-vitest-run-because-project-info-reads-a-system-flow-yaml-with-no-version.md

### 2026-10-08T08:50:32Z orchestrator
Resolved: Replan approved; tasks unchanged, value 37.50 USD/week from the operator's input.
