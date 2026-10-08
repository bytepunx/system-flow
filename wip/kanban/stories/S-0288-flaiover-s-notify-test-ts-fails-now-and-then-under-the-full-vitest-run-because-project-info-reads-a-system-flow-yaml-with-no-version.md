---
id: S-0288
type: story
nature: remediation
title: flaiover's notify.test.ts fails now and then under the full vitest run because project.info reads a system-flow.yaml with no version
status: review
owner: alex
created: 2026-10-06T09:56:51Z
updated: 2026-10-08T09:08:25Z
transitions:
  - to: ready
    at: 2026-10-08T08:40:50Z
    by: alex
  - to: in-progress
    at: 2026-10-08T08:58:41Z
    by: agent-S-0288
  - to: review
    at: 2026-10-08T09:08:25Z
    by: agent-S-0288
tags: [flaiover, tests]
touches: [flaiover/src/lib/server/notify.test.ts, design/issues/I-0085-flaiover-s-notify-test-ts-fails-now-and-then-under-the-full-vitest-run-because-project-info-reads-a-system-flow-yaml-with-no-version.md, design/issues/summary.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 593
  turns:
    - day: 2026-10-08
      ceremony: 2
      hand_edits: 1
      work: 19
  models:
    - model: claude-opus-5-5
      input: 56
      output: 11808
      cache_read: 2824921
      cache_write: 180740
      cost: 2.1581
  strategic:
    - kind: planner
      seconds: 187
      estimated: true
      models:
        - model: claude-haiku-4-5-20251001
          input: 274
          output: 68
          cache_read: 1510375
          cache_write: 103716
          cost: 0.3211
        - model: claude-opus-5-5
          input: 144
          output: 21541
          cache_read: 8495734
          cache_write: 270812
          cost: 5.1771
    - kind: orchestrator
      seconds: 172
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 83
          output: 1558
          cache_read: 22604263
          cache_write: 16434
          cost: 5.5734
cost_of_delay:
  inputs:
    time_lost_per_cycle: 15m
    by: alex
    at: 2026-10-08T08:49:21Z
  value: 37.5
  by: planner-S-0288
  at: 2026-10-08T08:50:18Z
forecast:
  duration: 20m
  delivery: 2026-10-08T09:14:00Z
  basis: "Its own forecast of 20m; 1st in the pull order with an in-progress limit of 5, behind S-0232, S-0321, S-0322, S-0340 and S-0341."
  by: flai
  at: 2026-10-08T08:53:22Z
finalized:
  by: alex
  at: 2026-10-07T02:20:05Z
---
# S-0288 flaiover's notify.test.ts fails now and then under the full vitest run because project.info reads a system-flow.yaml with no version

## Goal

This story remediates [I-0085](../../../design/issues/I-0085-flaiover-s-notify-test-ts-fails-now-and-then-under-the-full-vitest-run-because-project-info-reads-a-system-flow-yaml-with-no-version.md), "flaiover's notify.test.ts fails now and then under the full vitest run because project.info reads a system-flow.yaml with no version". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [x] The cause I-0085 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0085 is closed with `flai issue close I-0085 --reason` saying what fixed it

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
| `flaiover/src/lib/server/notify.test.ts` | declared (layout) | The flaky test, its `setManifest` helper, and the reproducing test. |
| `design/issues/I-0085-…md` | declared (design) | Criterion 2 closes the issue. |
| `design/issues/summary.md` | declared (design) | `flai issue close` regenerates it. |

`flai touches suggest` lists co-change from the three declared paths, all of it from other issues, `flai-cli.md`, and user docs that this test-only fix does not change, so none is added. No folder touch is kept. `design/issues` is in the manifest's `claims.shared`, so the two issue paths hold no story.

Figures:

- Forecast 20m, adjusted from flai's 14m (median 160 s per unit of size over 5 small remediations, size 5): this one must show a race gone, which means a deterministic reproduction and five full flaiover vitest runs. Delivery 2026-10-08T09:10:00Z, flai's 09:04 moved by the same 6m.
- Cost of delay 37.50 USD a week, as `flai cod` works it out from the operator's inputs (15m lost per 168h cycle at 150 USD an hour). The operator raised the input from 5m to 15m on TH-0378, for I-0085's three stopped close-outs in two days; the value follows it unadjusted.
