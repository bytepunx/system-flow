---
id: S-0530
type: story
nature: improvement
title: flai test on flaiover paths runs no svelte-check, so a type error a task's tests passed over first shows in the close-out
status: done
owner: alex
created: 2026-10-28T12:54:31Z
updated: 2026-10-28T21:29:13Z
transitions:
  - to: ready
    at: 2026-10-28T20:48:59Z
    by: orchestrator
  - to: in-progress
    at: 2026-10-28T20:49:08Z
    by: agent-S-0530
  - to: review
    at: 2026-10-28T21:07:50Z
    by: agent-S-0530
  - to: done
    at: 2026-10-28T21:29:13Z
    by: orchestrator
tags: [flaiover, testing, scripts]
topics: [testing]
touches: [system-flow.yaml, scripts/flaiover-check.sh, scripts/flaiover-lint.sh, scripts/README.md, flai/internal/verify/select_test.go, design/conventions/tooling.md, design/system/devex.md, design/issues/I-0206-flai-test-on-flaiover-paths-runs-no-svelte-check-so-a-type-error-a-task-s-tests-passed-over-first-shows-in-the-close-out.md, design/issues/summary.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 1135
  priced_by: harness
  turns:
    - day: 2026-10-28
      ceremony: 2
      test_runs: 1
      hand_edits: 1
      work: 19
  models:
    - model: claude-opus-5-5
      input: 106
      output: 22411
      cache_read: 4471222
      cache_write: 230980
      cost: 2.9888
  strategic:
    - kind: orchestrator
      seconds: 403
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 114
          output: 1951
          cache_read: 13969560
          cache_write: 34675
          cost: 3.0805
        - model: claude-sonnet-5-5
          input: 4
          output: 40
          cache_read: 12921
          cache_write: 29652
          cost: 0.0294
cost_of_delay:
  inputs:
    time_lost_per_cycle: 10m
    by: flai
    at: 2026-10-28T12:54:31Z
  value: 25
  by: planner-S-0530
  at: 2026-10-28T18:21:32Z
forecast:
  duration: 25m
  delivery: 2026-10-28T21:19:00Z
  basis: "Its own forecast of 25m; 1st in the pull order with an in-progress limit of 5, behind S-0416."
  by: flai
  at: 2026-10-28T20:49:02Z
finalized:
  by: alex
  at: 2026-10-28T12:56:09Z
---
# S-0530 flai test on flaiover paths runs no svelte-check, so a type error a task's tests passed over first shows in the close-out

## Goal

This story remediates [I-0206](../../../design/issues/I-0206-flai-test-on-flaiover-paths-runs-no-svelte-check-so-a-type-error-a-task-s-tests-passed-over-first-shows-in-the-close-out.md), "flai test on flaiover paths runs no svelte-check, so a type error a task's tests passed over first shows in the close-out". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [x] The cause I-0206 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0206 is closed with `flai issue close I-0206 --reason` saying what fixed it

## Tasks
- T-2265 flai test runs svelte-check on the flaiover project when a changed TypeScript, JavaScript, or Svelte file under flaiover selects its new tier
- T-2266 tooling.md and devex.md say flai test runs svelte-check for a changed flaiover file, and I-0206 is closed naming the svelte-check tier

## Notes

Cost of delay inputs set by flai from I-0206. time_lost_per_cycle 10m: 10m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-24T22:32:29Z, 3.6 days before this story; under one cycle counts as one).

### Planning

Proposed fix, from the I-0206 instance (S-0434's T-1845, TS2783): `system-flow.yaml` gains a `svelte-check` tier outside `--all`, after `vitest`. It is selected by a changed `.ts`, `.js`, or `.svelte` file under `flaiover/src`, a root config file of `flaiover`, or its own script. It runs `scripts/flaiover-check.sh`, which runs `pnpm check` on the whole project, since svelte-check takes no file list (D-0244). The `flaiover` tier covers it, so `--all` and the close-out run svelte-check once. This follows S-0319, which added `flaiover-lint` for I-0105 in the same way.

Touches, all files and no folders:

| Touch | Source |
|-------|--------|
| `system-flow.yaml` | design: the `tests` tiers (tooling.md project additions) |
| `scripts/flaiover-check.sh` | layout: new script beside `flaiover-lint.sh` |
| `scripts/flaiover-lint.sh` | layout: its comment says svelte-check stays in `flaiover-test.sh` |
| `scripts/README.md` | layout: one row per script |
| `flai/internal/verify/select_test.go` | layout: `TestRepositoryManifestLintsChangedFlaioverFiles` asserts the selected tiers (D-0244) |
| `design/conventions/tooling.md` | design: project additions describe the tiers |
| `design/system/devex.md` | design: the close-out row names where svelte-check runs |
| `design/issues/I-0206-…md`, `design/issues/summary.md` | criterion 2: `flai issue close` |

`flai touches suggest` listed co-change files: `design/system/flai-cli.md`, `docs/users/flai.md`, and other issues. None is taken, because no flai command or flag changes.

`system-flow.yaml` is also in S-0502's touches, so S-0530 is held while S-0502 is in progress. No `after` is needed, since S-0530 needs nothing from S-0502 on main.

Figures:

- Forecast: `flai forecast` gave 16m (size 11 at 86 s per unit). It is raised to 25m because S-0319, the like change, took 1267 s (21m), and T-2265 also times the new tier. The delivery moves the same 9m, to 2026-10-29T08:40Z.
- Cost of delay: `flai cod` gives 25.00 USD a week, from the 10m `time_lost_per_cycle` input at 150 USD an hour. It stands as given.

### Accepted by the orchestrator

- Verified: 01fa007fde9671f4ac4a42f64493e3f3a9f56e87
- Batch: 20261028-210833
- Enqueued verified: 01fa007fde9671f4ac4a42f64493e3f3a9f56e87
- At: 2026-10-28T21:29:13Z

Verdict: meets all criteria
- 1: system-flow.yaml, scripts/flaiover-check.sh, scripts/flaiover-lint.sh, scripts/README.md, flai/internal/verify/select_test.go, design/conventions/tooling.md, design/system/devex.md
- 2: design/issues/I-0206-flai-test-on-flaiover-paths-runs-no-svelte-check-so-a-type-error-a-task-s-tests-passed-over-first-shows-in-the-close-out.md, design/issues/summary.md
