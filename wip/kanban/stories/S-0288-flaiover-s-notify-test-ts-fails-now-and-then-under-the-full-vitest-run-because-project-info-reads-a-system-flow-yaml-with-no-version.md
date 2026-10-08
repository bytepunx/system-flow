---
id: S-0288
type: story
nature: remediation
title: flaiover's notify.test.ts fails now and then under the full vitest run because project.info reads a system-flow.yaml with no version
status: ready
owner: alex
created: 2026-10-06T09:56:51Z
updated: 2026-10-08T08:49:21Z
transitions:
  - to: ready
    at: 2026-10-08T08:40:50Z
    by: alex
tags: [flaiover, tests]
touches: [flaiover/src/lib/server/notify.test.ts, design/issues/I-0085-flaiover-s-notify-test-ts-fails-now-and-then-under-the-full-vitest-run-because-project-info-reads-a-system-flow-yaml-with-no-version.md, design/issues/summary.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: sum
  seconds: 0
  models: []
  strategic:
    - kind: planner
      seconds: 125
      estimated: true
      models:
        - model: claude-haiku-4-5-20251001
          input: 88
          output: 27
          cache_read: 563577
          cache_write: 40634
          cost: 0.1201
        - model: claude-opus-5-5
          input: 114
          output: 16583
          cache_read: 6069445
          cache_write: 168164
          cost: 3.4444
    - kind: orchestrator
      seconds: 55
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 17
          output: 291
          cache_read: 4347064
          cache_write: 3097
          cost: 1.0718
cost_of_delay:
  inputs:
    time_lost_per_cycle: 15m
    by: alex
    at: 2026-10-08T08:49:21Z
  value: 12.5
  by: planner-S-0288
  at: 2026-10-08T08:42:17Z
forecast:
  duration: 20m
  delivery: 2026-10-08T09:03:00Z
  basis: "Its own forecast of 20m; 1st in the pull order with an in-progress limit of 5, behind S-0232, S-0321, S-0322, S-0340 and S-0341."
  by: flai
  at: 2026-10-08T08:42:40Z
finalized:
  by: alex
  at: 2026-10-07T02:20:05Z
---
# S-0288 flaiover's notify.test.ts fails now and then under the full vitest run because project.info reads a system-flow.yaml with no version

## Goal

This story remediates [I-0085](../../../design/issues/I-0085-flaiover-s-notify-test-ts-fails-now-and-then-under-the-full-vitest-run-because-project-info-reads-a-system-flow-yaml-with-no-version.md), "flaiover's notify.test.ts fails now and then under the full vitest run because project.info reads a system-flow.yaml with no version". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0085 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0085 is closed with `flai issue close I-0085 --reason` saying what fixed it

## Tasks
- T-1369 notify.test.ts reproduces the stale manifest read, and its setManifest writes atomically and tells the Repo
- T-1370 Close I-0085 saying what fixed it

## Notes

Cost of delay inputs set by flai from I-0085. time_lost_per_cycle 5m: 5m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-06T07:08:02Z, 0.1 days before this story; under one cycle counts as one).

### Planning

Proposed cause, read from `flaiover/src/lib/server/notify.test.ts` and `flaiover/src/lib/server/repo.ts`:

- The test "posts an entry that appears after it started…" calls `addQuestion` before `setManifest`. `addQuestion` calls `repo.changed(...)` on a fresh `Repo`, whose layout is not yet known, so `changed` starts a background `manifest()` read (`project.info`) and caches its promise under `project`.
- `setManifest` then rewrites `system-flow.yaml` with `writeFile`, which truncates it first, and never tells the `Repo` the manifest changed.
- If the background read lands while the file is empty, `project.info` fails with "version is required" (the S-0220 instance). If it lands before the write, the cached manifest has no `notify_url`, `startNotifier` returns null, and line 78 fails with "expected null not to be null" (the S-0295 and S-0336 instances).
- In production flai on the host reports a manifest change through `repo.changed('system-flow.yaml')` (`repo.ts`), which drops the cached answer, so the race is in the test's setup, not in `Repo`.

Touches:

| Touch | Source | Why |
|-------|--------|-----|
| `flaiover/src/lib/server/notify.test.ts` | layout | The flaky test, its `setManifest` helper, and the reproducing test. |
| `design/issues/I-0085-…md` | design | Criterion 2 closes the issue. |
| `design/issues/summary.md` | design | `flai issue close` regenerates it. |

`flai touches suggest` listed nothing: the story declared no touches. No folder touch is kept. `design/issues` is in the manifest's `claims.shared`, so the two issue paths hold no story.

Figures:

- Forecast 20m, adjusted from flai's 6m: the 6m is the median of small remediations, but this one must show a race gone, which means a deterministic reproduction and several full flaiover vitest runs. Delivery moved by the same 14m to 2026-10-08T09:39:00Z.
- Cost of delay 12.50 USD a week, as `flai cod` worked it out from the inputs (5m lost per 168h cycle at 150 USD an hour). It stands, though I-0085 now has three instances in two days, each a stopped close-out: the inputs were set before the second and third, so they understate it. The inputs are the operator's; the plan thread says so.
