---
id: S-0846
type: story
nature: improvement
title: flai serve's auto-upgrade restart interrupts the close-out jobs it was running, which start again from scratch
status: ready
owner: alex
created: 2026-11-17T08:32:57Z
updated: 2026-11-18T09:20:10Z
transitions:
  - to: ready
    at: 2026-11-18T09:20:04Z
    by: orchestrator
tags: []
topics: [serve, close-out, release]
touches: [design/system/flai-cli.md, flai/internal/serve/upgrade.go, flai/internal/serve/upgrade_test.go, flai/internal/serve/verifyjob.go, flai/internal/serve/verifyjob_test.go, flai/internal/serve/resolvejob.go, flai/internal/serve/resolvejob_test.go, flai/internal/serve/serve.go, docs/operators/index.md, docs/users/flai.md, design/issues/I-0388-flai-serve-s-auto-upgrade-restart-interrupts-the-close-out-jobs-it-was-running.md, design/issues/summary.md]
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
      seconds: 92
      estimated: true
      models:
        - model: claude-haiku-5-5
          input: 46
          output: 1656
          cache_read: 4161099
          cache_write: 66411
          cost: 0.2888
        - model: claude-opus-5-5
          input: 87
          output: 9661
          cache_read: 7310582
          cache_write: 129317
          cost: 2.7491
    - kind: orchestrator
      seconds: 99
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 20
          output: 99
          cache_read: 1284106
          cache_write: 179384
          cost: 0.4582
    - kind: analyzer
      seconds: 99
      estimated: true
      models:
        - model: claude-fable-5-1
          input: 47
          output: 6584
          cache_read: 349360
          cache_write: 34306
          cost: 1.1031
        - model: claude-haiku-5-5
          input: 36
          output: 14843
          cache_read: 2398931
          cache_write: 47318
          cost: 0.1557
cost_of_delay:
  inputs:
    time_lost_per_cycle: 2h
    by: flai
    at: 2026-11-17T08:32:57Z
  value: 300
  by: planner-S-0846
  at: 2026-11-17T12:34:08Z
forecast:
  duration: 1h30m
  delivery: 2026-11-18T14:05:00Z
  basis: "Its own forecast of 1h30m; 5th in the pull order with an in-progress limit of 3, behind S-0807, S-0810, S-0839, S-0843, S-0847, S-0808 and S-0844."
  by: flai
  at: 2026-11-18T09:20:10Z
finalized:
  by: orchestrator
  at: 2026-11-17T12:38:55Z
---
# S-0846 flai serve's auto-upgrade restart interrupts the close-out jobs it was running, which start again from scratch

## Goal

This story remediates [I-0388](../../../design/issues/I-0388-flai-serve-s-auto-upgrade-restart-interrupts-the-close-out-jobs-it-was-running.md), "flai serve's auto-upgrade restart interrupts the close-out jobs it was running, which start again from scratch", which the analysis in [2026-11-16-all.md](../../../design/analysis/2026-11-16-all.md) found. The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0388 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0388 is closed with `flai issue close I-0388 --reason` saying what fixed it

## Tasks
- T-3583 Decide in an ADR that the auto-upgrade drains the close-out jobs before it restarts and that serve asks again once a job a restart interrupted
- T-3584 Hold the auto-upgrade while a close-out or rebase and resolve job runs, take no new job meanwhile, and upgrade once the bound passes
- T-3585 Ask again once, at serve's start, a close-out or rebase and resolve job a restart left interrupted
- T-3586 Document the drain and the ask again for operators and users, and close I-0388

## Notes

The analyzer found I-0388: the analysis the goal links gives its evidence.

Cost of delay inputs set by flai from I-0388. time_lost_per_cycle 2h carried over from I-0388's Impact section. Its Impact time_lost_per_cycle was taken rather than the 8m derived from its cost and count.

Strategic usage carried over from I-0388: analyzer 1m39s, 1.2588 USD.

### Planning

The cause, from the code: `upgrade.go` `look` waits only for dashboard writes (`hostapi.WritesRunning`) before it asks the host to upgrade (ADR-0158 §5). The host's restart SIGTERMs serve, whose verify pass stops each job's process group (`runCloseOut`, `terminateThenKill`) and leaves its job file running; the next serve's `verifyPass.sweep` ends it `interrupted` and asks nothing again. Rebase and resolve jobs share that pass. The proposed fix, decided in T-3583's ADR before any code: drain (no upgrade while a job runs, no new job while an upgrade is due, bounded at 45 minutes) and ask again once at start what a restart still interrupted.

Touches, all files; no folder touch:

| Touch | Source |
|-------|--------|
| `flai/internal/serve/upgrade.go`, `upgrade_test.go` | layout: the upgrade's `look` and its tests |
| `flai/internal/serve/verifyjob.go`, `verifyjob_test.go` | layout: the verify pass, `startJob`, `sweep` |
| `flai/internal/serve/resolvejob.go`, `resolvejob_test.go` | layout: the rebase and resolve job, which shares the pass |
| `flai/internal/serve/serve.go` | layout: where the verify pass is built and wired |
| `design/system/flai-cli.md` | design: § Commands › The close-out job |
| `docs/operators/index.md`, `docs/users/flai.md` | design: the close-out job and auto-upgrade sections |
| `design/issues/I-0388-…md`, `design/issues/summary.md` | criterion 2, `flai issue close` |

The new ADR's file is not named: its number is not known until `flai adr new`, and `flai task done` widens T-3583's touches with it. `flai touches suggest` listed co-change paths none of which this change needs (`docs/users/flai-reference.md` changes only with a flag, and no flag is added; `design/adrs/README.md` is generated).

No `after`: the story reads no output of another open story.

Forecast: `flai forecast` gave 26m; raised to 1h30m, reason in the forecast's basis. Cost of delay: 300 USD a week, as `flai cod` works it from the inputs flai carried over from I-0388 (2h lost per 168h cycle at 150 USD an hour); it stands, since the issue counts 11 interrupted close-outs in W46 alone.
