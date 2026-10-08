---
id: TH-0378
title: "S-0288 plan: reproduce and fix the stale manifest read in notify.test.ts, then close I-0085"
anchor:
  path: wip/kanban/stories/S-0288-flaiover-s-notify-test-ts-fails-now-and-then-under-the-full-vitest-run-because-project-info-reads-a-system-flow-yaml-with-no-version.md
  item: S-0288
status: resolved
participants: [planner-S-0288, alex]
created: 2026-10-08T08:42:51Z
updated: 2026-10-08T08:49:32Z
---

# TH-0378 S-0288 plan: reproduce and fix the stale manifest read in notify.test.ts, then close I-0085

On wip/kanban/stories/S-0288-flaiover-s-notify-test-ts-fails-now-and-then-under-the-full-vitest-run-because-project-info-reads-a-system-flow-yaml-with-no-version.md.

## Entries

### 2026-10-08T08:42:51Z planner-S-0288
S-0288 is planned with two tasks in two layers.

## Tasks

| Layer | Task | Touches | After |
|-------|------|---------|-------|
| 1 | T-1369 notify.test.ts reproduces the stale manifest read, and its setManifest writes atomically and tells the Repo | `flaiover/src/lib/server/notify.test.ts` | none |
| 2 | T-1370 Close I-0085 saying what fixed it | `design/issues/I-0085-…md`, `design/issues/summary.md` | T-1369 |

T-1370 waits because its close reason names the fix and the runs T-1369 makes.

## Proposed cause

- In "posts an entry that appears after it started…", `addQuestion` runs before `setManifest`. Its `repo.changed(...)` on a fresh `Repo` starts a background `project.info` read that `remember` caches.
- `setManifest` truncates and rewrites `system-flow.yaml` and never tells the `Repo`.
- A read during the truncation fails with "version is required" (S-0220). A read before the write caches a manifest with no `notify_url`, so `startNotifier` returns null at line 78 (S-0295, S-0336).

## Figures

- Forecast 20m, up from flai's 6m: the race needs a deterministic reproduction and five full flaiover vitest runs. Delivery 2026-10-08T09:39:00Z.
- Cost of delay 12.50 USD a week, from the inputs as they stand (5m per 168h cycle).

## Assumptions

1. The race is in the test's setup only. In production, flai on the host reports a manifest change with `repo.changed('system-flow.yaml')`, which drops the cached answer. So `repo.ts` and `notify.ts` stay unchanged unless T-1369's reproduction shows otherwise.
2. Five other flaiover tests rewrite `system-flow.yaml`. None has a recorded instance, so they are out of scope.
3. `design/issues` is in `claims.shared`, so the issue touches hold no other story. No folder touch is kept.

## For the operator, no answer needed to proceed

The cost of delay inputs predate two of I-0085's three instances. Three stopped close-outs in two days suggest about 15m lost a week, not 5m. If you agree, set `time_lost_per_cycle: 15m`. The value would then be 37.50 USD a week. I left the inputs as they are because they are yours.

### 2026-10-08T08:49:32Z alex
Resolved.
