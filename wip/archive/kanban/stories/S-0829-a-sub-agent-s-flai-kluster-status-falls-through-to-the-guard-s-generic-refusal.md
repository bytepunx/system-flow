---
id: S-0829
type: story
nature: remediation
title: A sub-agent's flai kluster status falls through to the guard's generic refusal, against kluster.md's rule that status is a read every agent may run
status: done
owner: alex
created: 2026-11-17T08:17:21Z
updated: 2026-11-17T13:23:36Z
transitions:
  - to: ready
    at: 2026-11-17T12:35:54Z
    by: orchestrator
  - to: in-progress
    at: 2026-11-17T12:56:17Z
    by: agent-S-0829
  - to: review
    at: 2026-11-17T13:22:06Z
    by: agent-S-0829
  - to: done
    at: 2026-11-17T13:23:36Z
    by: orchestrator
tags: []
topics: [guard, kluster]
touches: [flai/internal/guard/guard_tables.go, flai/internal/guard/kluster_test.go, design/system/flai-cli.md, design/system/kluster.md, design/issues/I-0355-a-sub-agent-s-flai-kluster-status-falls-through-to-the-guard-s-generic-refusal.md, design/issues/summary.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 1561
  priced_by: harness
  turns:
    - day: 2026-11-17
      ceremony: 1
      work: 31
  models:
    - model: claude-opus-5-5
      input: 70
      output: 11308
      cache_read: 6845590
      cache_write: 251828
      cost: 3.6102
  strategic:
    - kind: planner
      seconds: 181
      estimated: true
      models:
        - model: claude-haiku-5-5
          input: 232
          output: 847
          cache_read: 8653564
          cache_write: 549702
          cost: 0.6288
        - model: claude-opus-5-5
          input: 454
          output: 17970
          cache_read: 29493521
          cache_write: 1735506
          cost: 9.596
    - kind: orchestrator
      seconds: 276
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 38
          output: 316
          cache_read: 1659292
          cache_write: 167459
          cost: 0.5721
        - model: claude-sonnet-5-5
          input: 12
          output: 114
          cache_read: 140246
          cache_write: 36521
          cost: 0.1427
cost_of_delay:
  inputs:
    time_lost_per_cycle: 5m
    by: flai
    at: 2026-11-17T08:17:21Z
  value: 12.5
  by: planner-S-0829
  at: 2026-11-17T12:28:28Z
forecast:
  duration: 40m
  delivery: 2026-11-17T13:58:00Z
  basis: "Its own forecast of 40m; 2nd in the pull order with an in-progress limit of 3, behind S-0806, S-0834, S-0841 and S-0811."
  by: flai
  at: 2026-11-17T12:55:13Z
finalized:
  by: orchestrator
  at: 2026-11-17T12:35:32Z
---
# S-0829 A sub-agent's flai kluster status falls through to the guard's generic refusal, against kluster.md's rule that status is a read every agent may run

## Goal

This story remediates [I-0355](../../../design/issues/I-0355-a-sub-agent-s-flai-kluster-status-falls-through-to-the-guard-s-generic-refusal.md), "A sub-agent's flai kluster status falls through to the guard's generic refusal, against kluster.md's rule that status is a read every agent may run". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [x] The cause I-0355 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0355 is closed with `flai issue close I-0355 --reason` saying what fixed it

## Tasks
- T-3559 The guard lets every sub-agent and strategic role run flai kluster status and deployments
- T-3560 The design says flai kluster status and deployments are on the guard's sub-agent reads
- T-3562 Close I-0355 with what fixed it

## Notes

Cost of delay inputs set by flai from I-0355. time_lost_per_cycle 5m: 5m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-11-13T07:33:55Z, 4 days before this story; under one cycle counts as one).

### Planning

Proposed fix, from the instance and the code: a sub-agent's `flai kluster status` passes `klusterWrite` and `klusterLine` in `Decide`, then reaches `asSubAgent`, whose `reads` looks the command up in `cliReads` (`flai/internal/guard/guard_tables.go`). `cliReads` has no `kluster` entry, so the call falls to the generic refusal. Add `"kluster": {"status", "deployments"}` there: `design/system/kluster.md` § Reads says both are open to every agent, and the guard's own `reachesServed` refusal already tells served sessions so. A strategic role held by its file takes the same subcommands, since `subcommandsOf` in `roles.go` reads `cliReads`. `TestFlaiKlusterWritesAreTheOperators` in `kluster_test.go` checks only that a read is not refused as a write, which is why the gap passed; the new test asserts the read is allowed.

Touches, and where each came from:

| Touch | Source | Why |
|-------|--------|-----|
| `flai/internal/guard/guard_tables.go` | layout, design | `cliReads`, the sub-agent read table the call falls through |
| `flai/internal/guard/kluster_test.go` | layout, co-change | the kluster guard tests; `flai touches suggest` from `guard.go` names `guard_test.go` at 49%, but the kluster cases live here |
| `design/system/flai-cli.md` | co-change (15%), design | the `flai guard` row lists every flai command a sub-agent may run |
| `design/system/kluster.md` | design | the rule I-0355 cites; its guard history gains S-0829 |
| `design/issues/I-0355-…md`, `design/issues/summary.md` | criteria | `flai issue close` writes both |

Left out: `flai/internal/guard/guard.go`, whose comment already says status reads; `roles.go`, which reads `cliReads` and needs no change; `docs/users/flai-reference.md`, since no command's help changes; `flai/cmd/guard.go` and `guard_test.go`, the other co-change hits, which this fix does not reach. No folder touch is kept. No `after`: the fix reads nothing another open story writes.

Forecast: `flai forecast` gave 6m with no touches, and 20m with these six (147 s per unit of size over 70 done remediation stories, size 8). Set to 40m: S-0680, the same kind of fix (a missing guard read form, its docs, and closing the issue in three tasks), took 2178 s, and a close-out here alone takes about 10m. The delivery, 2026-11-17T18:49:00Z, is kept as `flai forecast` worked it out from the pull order.

Cost of delay: `flai cod` gives 12.50 USD a week from the operator's input, 5m lost per 168h cycle at 150 USD an hour. Kept as it stands: one instance, and a sub-agent can still work around it by asking the story's agent.

### Accepted by the orchestrator

- Verified: 209451e0ad8f2eaeff9e920e4f9b4427858cd629
- At: 2026-11-17T13:23:36Z

Verdict: pass at 209451e0ad8f2eaeff9e920e4f9b4427858cd629
- 1: flai/internal/guard/guard_tables.go, flai/internal/guard/kluster_test.go, design/system/kluster.md, design/system/flai-cli.md
- 2: design/issues/I-0355-a-sub-agent-s-flai-kluster-status-falls-through-to-the-guard-s-generic-refusal.md, design/issues/summary.md
