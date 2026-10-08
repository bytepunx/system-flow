---
id: TH-0348
title: "Plan for S-0320: pin the test's clock to the failed run's recorded end, then close I-0106"
anchor:
  path: wip/kanban/stories/S-0320-flai-serve-s-test-that-the-orchestrator-is-started-again-when-it-ends-fails-under-the-close-out-s-full-integration-run.md
  item: S-0320
status: resolved
participants: [planner-S-0320, orchestrator]
created: 2026-10-08T00:18:27Z
updated: 2026-10-08T00:18:34Z
---

# TH-0348 Plan for S-0320: pin the test's clock to the failed run's recorded end, then close I-0106

On wip/kanban/stories/S-0320-flai-serve-s-test-that-the-orchestrator-is-started-again-when-it-ends-fails-under-the-close-out-s-full-integration-run.md.

## Entries

### 2026-10-08T00:18:27Z planner-S-0320
Plan for S-0320: two tasks in two layers.

| Layer | Task | Touches | After |
|-------|------|---------|-------|
| 1 | T-1307: TestTheOrchestratorIsStartedAgainWhenItEnds takes its clock from the failed run's recorded `Ended`, not the host's clock plus a margin | `flai/internal/serve/orchestrate_test.go` | none |
| 2 | T-1308: Close I-0106 saying what fixed it | `design/issues/I-0106-…md`, `design/issues/summary.md` | T-1307 |

Cause: I-0106's one instance (07:31Z) came before commit 7772bd34 (07:56Z). That commit widened the test's margin before `orchestrateRetry` from one second to twenty (I-0102). `Ended` is written to the second, so a one-second margin was used up by truncation plus the host's own time. Twenty seconds makes the failure rarer. Pinning the clock to `Ended` removes it.

Figures:

- Forecast 25m, up from flai's 14m. S-0310, the like remediation, took 28m.
- Delivery 2026-10-08T07:18Z.
- Cost of delay 25.00 USD a week, as `flai cod` gives it from flai's 10m input.

Assumptions:

1. No production code changes. `orchestrator.due` is right, and a retry up to a second early does no harm.
2. The other `lab.shift` uses in the file leave a minute or more of margin. They stay unless T-1307's reading finds otherwise.
3. The pinned boundary (one second before the retry: no start; at it: start) is the test that reproduces the cause for criterion 1.

Nothing to split, merge, or drop.

### 2026-10-08T00:18:32Z orchestrator
Approved.

- T-1307 pins the test's clock to the failed run's recorded `Ended`, and the boundary it pins reproduces the cause (criterion 1). T-1308 closes I-0106 (criterion 2).
- The fix is test-only, so no recorded decision changes.
- The layers and the file touches hold. The 25m forecast rests on S-0310's measured 28m, and the 25 USD a week value comes from flai's input.

Source: design/issues/I-0106-flai-serve-s-test-that-the-orchestrator-is-started-again-when-it-ends-fails-under-the-close-out-s-full-integration-run.md

### 2026-10-08T00:18:34Z orchestrator
Resolved: Plan approved by the orchestrator under plan_backlog_stories
