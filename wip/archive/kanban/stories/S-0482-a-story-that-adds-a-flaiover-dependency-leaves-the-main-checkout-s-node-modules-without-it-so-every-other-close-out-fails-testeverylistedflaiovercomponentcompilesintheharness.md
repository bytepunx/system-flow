---
id: S-0482
type: story
nature: remediation
title: A story that adds a flaiover dependency leaves the main checkout's node_modules without it, so every other close-out fails TestEveryListedFlaioverComponentCompilesInTheHarness
status: done
owner: alex
created: 2026-10-21T12:00:09Z
updated: 2026-10-21T12:49:26Z
transitions:
  - to: ready
    at: 2026-10-21T12:14:19Z
    by: orchestrator
  - to: in-progress
    at: 2026-10-21T12:19:11Z
    by: agent-S-0482
  - to: review
    at: 2026-10-21T12:36:06Z
    by: agent-S-0482
  - to: done
    at: 2026-10-21T12:49:26Z
    by: orchestrator
tags: []
touches: [scripts/integration.sh, scripts/README.md, design/system/devex.md, flai/internal/components/integration_test.go, design/issues/I-0178-a-story-that-adds-a-flaiover-dependency-leaves-the-main-checkout-s-node-modules-without-it-so-every-other-close-out-fails-testeverylistedflaiovercomponentcompilesintheharness.md, design/issues/summary.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 1033
  priced_by: harness
  turns:
    - day: 2026-10-21
      ceremony: 2
      hand_edits: 1
      work: 35
  models:
    - model: claude-opus-5-5
      input: 78
      output: 16595
      cache_read: 5303964
      cache_write: 158946
      cost: 2.6646
  strategic:
    - kind: orchestrator
      seconds: 281
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 55
          output: 700
          cache_read: 12407783
          cache_write: 25506
          cost: 2.7453
        - model: claude-sonnet-5-5
          input: 6
          output: 44
          cache_read: 41888
          cache_write: 39953
          cost: 0.0563
cost_of_delay:
  inputs:
    time_lost_per_cycle: 20m
    by: flai
    at: 2026-10-21T12:00:09Z
  value: 50
  by: planner-S-0482
  at: 2026-10-21T12:09:42Z
forecast:
  duration: 45m
  delivery: 2026-10-22T06:13:00Z
  basis: "flai forecast's 16m (median 115 s per unit of size over 27 done medium remediation stories, size 8) raised to 45m because the close-out of a flai/** change runs the integration and smoke tiers in full, and the worktree's first flaiover install adds to it; delivery is flai's 05:44Z moved by the same 29m."
  by: planner-S-0482
  at: 2026-10-21T12:09:42Z
finalized:
  by: alex
  at: 2026-10-21T12:03:18Z
---
# S-0482 A story that adds a flaiover dependency leaves the main checkout's node_modules without it, so every other close-out fails TestEveryListedFlaioverComponentCompilesInTheHarness

## Goal

This story remediates [I-0178](../../../design/issues/I-0178-a-story-that-adds-a-flaiover-dependency-leaves-the-main-checkout-s-node-modules-without-it-so-every-other-close-out-fails-testeverylistedflaiovercomponentcompilesintheharness.md), "A story that adds a flaiover dependency leaves the main checkout's node_modules without it, so every other close-out fails TestEveryListedFlaioverComponentCompilesInTheHarness". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [x] The cause I-0178 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0178 is closed with `flai issue close I-0178 --reason` saying what fixed it

## Tasks
- T-1994 scripts/integration.sh installs flaiover's dependencies in a worktree before the Go run, as flaiover-unit.sh does
- T-1996 flaioverModules links another checkout's node_modules only when it was installed from this checkout's lockfile, with a test that reproduces I-0178
- T-1997 Verify a worktree's integration tier runs on its own flaiover install, and close I-0178

## Notes

Cost of delay inputs set by flai from I-0178. time_lost_per_cycle 20m: 20m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-21T08:42:54Z, 0.1 days before this story; under one cycle counts as one).

### Planning

Proposed fix, from I-0178's one instance: give every worktree its own flaiover install from its own lockfile, and stop the integration tests from borrowing another checkout's install made from a different lockfile.

- `scripts/integration.sh` runs `scripts/flaiover-install.sh --if-needed` in a worktree (`CACHE_ROOT` differs from `ROOT`) before `go test`, as `scripts/flaiover-unit.sh` already does. That covers story worktrees and the merge queue's batch worktree, and the main checkout stays uninstalled, as `design/system/component-library.md` says.
- `flaioverModules` in `flai/internal/components/integration_test.go` accepts a `node_modules` only when its `.pnpm/lock.yaml` matches the checkout's `flaiover/pnpm-lock.yaml`, and fails naming the install otherwise. A unit test of that choice reproduces I-0178.

Touches, all files, no folder touch:

| Touch | Source |
| --- | --- |
| `scripts/integration.sh` | design: the script the integration tier runs (`system-flow.yaml` `tests`), where the install is missing |
| `scripts/README.md` | co-change with `scripts/integration.sh` (3 of 9 commits); its row describes the script |
| `design/system/devex.md` | design: its Test tiers row says which scripts install flaiover's dependencies in a worktree |
| `flai/internal/components/integration_test.go` | layout: `flaioverModules`, which links the main checkout's `node_modules` |
| `design/issues/I-0178-…md` | criterion 2: `flai issue close` writes it |
| `design/issues/summary.md` | criterion 2: regenerated when the issue closes |

`flai touches suggest` declared nothing at first; run from the two code paths it listed `harness.mjs`, `component-library.md`, `docs/users/flai.md`, a testdata component, and `scripts/smoke.sh` as past co-changes. None is left out by accident: the fix changes no harness, no command, and no user-facing behaviour, and `component-library.md`'s "The main checkout is never installed into" stays true. `props_test.go`'s `installed()` walks up to the main checkout too, but reads only svelte's compiler, so it is left alone.

Figures:

- Forecast 45m, adjusted from flai's 16m (median 115 s per unit of size over 27 done medium remediation stories, size 8). The close-out of a `flai/**` change runs the integration and smoke tiers in full, and the worktree's first `pnpm install` adds to it. Delivery 2026-10-22T06:13Z is flai's 05:44Z moved by the same 29m; the story is 44th in the pull order.
- Cost of delay value 50 USD a week, as `flai cod` computes it from the operator's input: 20m lost per 168h cycle at 150 USD an hour. Kept: the issue's one instance cost 20m, and the input is flai's from it. It likely understates the harm, since main fails the test for every close-out on the host until someone installs by hand, but the input is the operator's to raise.

### Accepted by the orchestrator

- Verified: 0974f9105bbd8e474920a9872208b18603529e64
- Batch: 20261021-123703
- Enqueued verified: 0974f9105bbd8e474920a9872208b18603529e64
- At: 2026-10-21T12:49:26Z

Verdict: Pass. Both criteria are met, every changed file is in the story's touches, and nothing installs into the main checkout.
- 1: scripts/integration.sh, flai/internal/components/integration_test.go, scripts/README.md, design/system/devex.md
- 2: design/issues/I-0178-a-story-that-adds-a-flaiover-dependency-leaves-the-main-checkout-s-node-modules-without-it-so-every-other-close-out-fails-testeverylistedflaiovercomponentcompilesintheharness.md, design/issues/summary.md
