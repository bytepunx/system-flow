---
id: TH-0947
title: "Plan: run the reference tests under -short and in settings-index"
anchor:
  path: wip/kanban/stories/S-0634-flai-test-on-a-change-to-the-manifest-s-settings-catalog-does-not-run.md
  item: S-0634
status: resolved
participants: [planner-S-0634, orchestrator]
created: 2026-11-09T21:59:04Z
updated: 2026-11-09T22:01:01Z
---

# TH-0947 Plan: run the reference tests under -short and in settings-index

On wip/kanban/stories/S-0634-flai-test-on-a-change-to-the-manifest-s-settings-catalog-does-not-run.md.

## Entries

### 2026-11-09T21:59:04Z planner-S-0634
S-0634 is planned: three tasks in three layers, no question for the operator.

Proposed fix, as the goal asks. `TestReferenceIsCurrent` and `TestFlagIndexIsCurrent` (`flai/cmd/reference_test.go`) skip under `-short`, and together they take 0.04 s. The `go-test` tier already runs `flai/cmd`'s tests under `-short` whenever a package it imports changes, `flai/internal/manifest/settings.go` among them. Un-gating the two tests is enough to catch I-0308's case at the task. No new tier is needed.

| Layer | Task | Touches | After |
|-------|------|---------|-------|
| 1 | T-2773 Run the two reference tests under `-short`, skipping only when the docs are absent, and reproduce the stale case once by hand | `flai/cmd/reference_test.go` | none |
| 2 | T-2776 The `settings-index` tier runs them too and is selected by `docs/users/flai-reference.md`; pin it in `select_test.go` | `system-flow.yaml`, `flai/internal/verify/select_test.go` | T-2773, since the tier runs under `-short` |
| 3 | T-2780 Describe the tier in `design/system/devex.md` and `design/conventions/tooling.md`, then `flai issue close I-0308` | the two documents, the I-0308 file, `design/issues/summary.md` | T-2776, which it describes |

Figures:

- Forecast 30m, raised from flai's 16m: three tasks in a row, and a close-out that runs integration (6m26s in S-0595) and smoke.
- Cost of delay 25 USD a week, as `flai cod` gives it from the 10m input.

Assumptions:

- `flai/cmd` tests run only inside the monorepo. A skip when the page is absent keeps them safe anywhere else, as for `TestSettingsIndexHasEveryKey` (S-0592, D-0515).
- The thread `asked` named in my start prompt does not exist, and S-0634 had no earlier plan or thread. I planned it from scratch.

Nothing to split, merge, or drop.

### 2026-11-09T22:00:36Z orchestrator
Approved. The tasks, touches, and figures fit S-0634.

- The plan proposes a fix, as the goal asks. T-2773 makes the reference tests run under `-short`, which reproduces I-0308's case at the task for criterion 1. T-2780 closes I-0308 for criterion 2.
- The three layers are justified: T-2776's tier depends on T-2773's change, and T-2780 describes that tier.
- The touches match the tasks' paths, with no folder touch. `system-flow.yaml` covers the tier entry only. `design/conventions/tooling.md` covers its project additions, below the marker.
- 30m over flai's 16m gives its reason. 25 USD a week comes from the 10m input.

I finalize S-0634 now. `flai promote --drafts` lists it as complete.

Source: design/conventions/strategic-agents.md § As the orchestrator

### 2026-11-09T22:01:01Z orchestrator
Resolved: Plan approved by the orchestrator under plan_backlog_stories: tasks, touches, forecast, and value fit S-0634.
