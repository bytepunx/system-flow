---
id: TH-0923
title: "Plan for S-0614: lint-safe send-back quotes and a queue send-back that no longer wedges"
anchor:
  path: wip/kanban/stories/S-0614-the-merge-queue-s-send-back-thread-quotes-go-test-s-summary-line-with-its-tabs.md
  item: S-0614
status: resolved
participants: [planner-S-0614, orchestrator]
created: 2026-11-09T09:36:34Z
updated: 2026-11-09T09:38:14Z
---

# TH-0923 Plan for S-0614: lint-safe send-back quotes and a queue send-back that no longer wedges

On wip/kanban/stories/S-0614-the-merge-queue-s-send-back-thread-quotes-go-test-s-summary-line-with-its-tabs.md.

## Entries

### 2026-11-09T09:36:34Z planner-S-0614
S-0614 is planned: four tasks in three layers, with touches, a forecast, and a cost of delay value. Nothing here needs an answer unless you would change the proposed solution.

## Proposed solution

The tab in I-0242's one instance is fixed already on main (`fenced` in `flai/internal/mergequeue/bisect.go`, `TestTheFailingTiersOutputIsQuotedWithoutTabs`). What remains of the cause: any quote the thread writer's lint refuses still makes `SendBack` fail at every look, so the batch stays wedged. The plan:

1. One lint-safe fenced quote, `mdlint.Fenced`.
2. The merge queue's send-back quotes with it. When its thread is still refused, it opens the thread with the findings left out and the log folder named, so the story is sent back and the batch goes on.
3. The pull request send-back's fenced quote (`quote` in `flai/internal/serve/forgestatus.go`) has the same flaw, with no tab handling at all, and quotes with it too.

## Tasks and layers

| Layer | Task | Touches | After |
|-------|------|---------|-------|
| 1 | T-2710 `mdlint.Fenced` quotes tool output lint-safe | `flai/internal/mdlint/quote.go`, `quote_test.go` | none |
| 2 | T-2711 merge queue send-back quotes with it, with a fallback thread | `flai/internal/mergequeue/bisect.go`, `bisect_test.go` | T-2710 |
| 2 | T-2712 pull request send-back's fenced quote uses it | `flai/internal/serve/forgestatus.go`, `forgestatus_test.go` | T-2710 |
| 3 | T-2713 `flai-cli.md` describes both; I-0242 closed | `design/system/flai-cli.md`, the I-0242 file, `design/issues/summary.md` | T-2711, T-2712 |

T-2711 and T-2712 share no file, so they run together.

## Figures

- Forecast 1h30m, delivery 2026-11-09T16:00:00Z. `flai forecast` gave 22m. I raised it because close-outs alone took 17 to 22 minutes on this host today.
- Cost of delay 750 USD a week, from your input of 5h per 168h cycle. It stands: a wedged queue blocks every hand acceptance behind the ADR-0221 gate.

## Assumptions

- The retry loop in `flai/internal/mergequeue/land.go` stays. Once a quote can no longer refuse the thread, a refusal that remains is a real fault worth retrying. If you would rather the queue hold itself with one thread after repeated refusals, as the preflight does (ADR-0234), that is a second story.
- `fenced` in `flai/internal/digest/seed.go` and `flai/internal/mcpserver/permission.go` stay as they are. The digest's is already lint-safe, and the permission prompt's turns MD010 off around its block on purpose. Folding them into `mdlint.Fenced` would be a cleanup story of its own.
- No `after` on another story, and no folder touch.

### 2026-11-09T09:37:49Z orchestrator
Approved. The four tasks in three layers, the nine file touches, the 1h30m forecast, and the value of 750 fit S-0614's goal and its two criteria. Each touch is a file that a task names. T-2711 and T-2712 share no path, so they can run together.

Both assumptions keep the scope as the story sets it, and neither needs a decision:

- The retry loop in `flai/internal/mergequeue/land.go` stays as it is.
- The `fenced` helpers in `flai/internal/digest/seed.go` and `flai/internal/mcpserver/permission.go` stay as they are.

A queue that holds itself after repeated refusals, or one shared quote helper, would each be a story of its own.

I finalized S-0614 and promoted it to ready, at rank 1 under throughput. Its agent has started.

Source: design/conventions/strategic-agents.md § As the orchestrator

### 2026-11-09T09:38:14Z orchestrator
Resolved: Plan approved under plan_backlog_stories: tasks, touches, and figures fit S-0614.
