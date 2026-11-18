---
id: S-0840
type: story
nature: remediation
title: "flai kluster up cannot refresh a stale kubeconfig: kluster up refuses an existing cluster before flai writes it"
status: done
owner: alex
created: 2026-11-17T08:32:51Z
updated: 2026-11-18T08:13:56Z
transitions:
  - to: ready
    at: 2026-11-17T12:35:54Z
    by: orchestrator
  - to: in-progress
    at: 2026-11-17T13:22:14Z
    by: agent-S-0840
  - to: review
    at: 2026-11-18T08:12:26Z
    by: agent-S-0840
  - to: done
    at: 2026-11-18T08:13:56Z
    by: orchestrator
blocked:
  - from: 2026-11-17T13:24:02Z
    until: 2026-11-18T07:11:06Z
    reason: awaiting the operator on TH-1305
  - from: 2026-11-18T07:12:17Z
    until: 2026-11-18T07:45:55Z
    reason: awaiting the operator on TH-1305
tags: [flai, kluster]
topics: [kluster]
touches: [flai/internal/kluster/kluster.go, flai/internal/kluster/kluster_test.go, flai/cmd/kluster.go, flai/cmd/kluster_test.go, docs/users/flai-reference.md, design/system/kluster.md, design/tech/kluster.md, design/issues/I-0367-flai-kluster-up-cannot-refresh-a-stale-kubeconfig-kluster-up-refuses-an.md, design/issues/summary.md, docs/users/flai.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 1746
  priced_by: harness
  turns:
    - day: 2026-11-17
      ceremony: 1
      work: 12
    - day: 2026-11-18
      ceremony: 1
      work: 44
  models:
    - model: claude-opus-5-5
      input: 132
      output: 29622
      cache_read: 16353827
      cache_write: 555413
      cost: 8.307
  strategic:
    - kind: planner
      seconds: 33
      estimated: true
      models:
        - model: claude-haiku-5-5
          input: 38
          output: 121
          cache_read: 2079043
          cache_write: 58508
          cost: 0.146
        - model: claude-opus-5-5
          input: 62
          output: 4569
          cache_read: 5232191
          cache_write: 69067
          cost: 1.705
    - kind: orchestrator
      seconds: 207
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 67
          output: 624
          cache_read: 2475903
          cache_write: 361477
          cost: 0.8915
        - model: claude-sonnet-5-5
          input: 4
          output: 40
          cache_read: 14309
          cache_write: 31633
          cost: 0.0371
cost_of_delay:
  inputs:
    time_lost_per_cycle: 30m
    by: flai
    at: 2026-11-17T08:32:51Z
  value: 75
  by: planner-S-0840
  at: 2026-11-17T12:28:53Z
forecast:
  duration: 40m
  delivery: 2026-11-17T14:25:00Z
  basis: "Its own forecast of 40m; 2nd in the pull order with an in-progress limit of 3, behind S-0806, S-0829, S-0834 and S-0811."
  by: flai
  at: 2026-11-17T12:57:33Z
finalized:
  by: orchestrator
  at: 2026-11-17T12:35:33Z
---
# S-0840 flai kluster up cannot refresh a stale kubeconfig: kluster up refuses an existing cluster before flai writes it

## Goal

This story remediates [I-0367](../../../design/issues/I-0367-flai-kluster-up-cannot-refresh-a-stale-kubeconfig-kluster-up-refuses-an.md), "flai kluster up cannot refresh a stale kubeconfig: kluster up refuses an existing cluster before flai writes it". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [x] The cause I-0367 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0367 is closed with `flai issue close I-0367 --reason` saying what fixed it

## Tasks
- T-3564 Propose the fix for I-0367 and record it in the kluster design
- T-3566 Cluster.Up writes the kubeconfig again for a cluster that runs, with a test that reproduces I-0367
- T-3570 flai kluster up says when it wrote the kubeconfig of a cluster that was up, its reference is regenerated, and I-0367 is closed

## Notes

Cost of delay inputs set by flai from I-0367. time_lost_per_cycle 30m: 30m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-11-13T20:28:29Z, 3.5 days before this story; under one cycle counts as one).

### Planning

The cause is in `Cluster.Up` in `flai/internal/kluster/kluster.go`. It runs `kluster up` and returns on its failure, so `kluster kubeconfig --name <name> --output` never runs while the cluster exists. The fix proposed here: `Up` reads `kluster status` first (`Cluster.Status`). When the cluster is listed and running, `Up` skips `kluster up` and writes the kubeconfig again. When the cluster is listed but stopped, `Up` refuses and names the cause. The flai serve job behind the dashboard's **Up** calls the same `Cluster.Up` (`flai/internal/hostapi/kluster.go`), so it gets the fix with no change. The first task puts this proposal to the operator before any code is written, as the goal asks.

Touches, and where each came from:

| Touch | Source | Why |
|-------|--------|-----|
| `flai/internal/kluster/kluster.go` | design, layout | `Cluster.Up` is the cause |
| `flai/internal/kluster/kluster_test.go` | co-change, layout | The test that reproduces I-0367 with the fake kluster |
| `flai/cmd/kluster.go` | co-change, design | `flai kluster up`'s help and its last line |
| `flai/cmd/kluster_test.go` | co-change (8 of 15 commits) | `TestKlusterUp` and a case for an existing cluster |
| `docs/users/flai-reference.md` | co-change (8 of 15 commits) | Generated from the command's help by `make flai-reference` |
| `design/system/kluster.md` | design | § The project's cluster and the `flai kluster up` row of § Commands |
| `design/tech/kluster.md` | design | Its Calls row: `up` now runs `status` first |

Left out: `docs/operators/settings.md` (co-change, 6 of 15), since the story adds no setting. `flai/internal/kluster/deploy.go` and `agentconfig.go` (co-change) are left out too: their refusals say to write the kubeconfig again with `flai kluster up`, which is true once this story lands. The flai serve and MCP files co-change lists are left out because the story does not change them. The touches name no folders.

Forecast: `flai forecast` gave 23m: the median of 147 s per unit of size over 70 done medium remediation stories on claude-opus-5-5, times size 9. I raised it to 40m. The story has three tasks in three layers, and a change under `flai/` makes the close-out run the integration and smoke tiers, about 10m as S-0828's did. I moved the delivery later by the same 17m.

Cost of delay: `flai cod` gives 75 USD a week: 30m lost per 168h cycle at 150 USD an hour. The value stands as flai worked it out from the inputs, which are the operator's.

`after`: none. No story writes output this one reads.

### Accepted by the orchestrator

- Verified: e67f75e76594142403f86e99177d21267c889509
- At: 2026-11-18T08:13:56Z

Verdict: pass
- 1: flai/internal/kluster/kluster.go, flai/internal/kluster/kluster_test.go, flai/cmd/kluster.go, flai/cmd/kluster_test.go, design/system/kluster.md, design/tech/kluster.md, docs/users/flai.md, docs/users/flai-reference.md
- 2: design/issues/I-0367-flai-kluster-up-cannot-refresh-a-stale-kubeconfig-kluster-up-refuses-an.md, design/issues/summary.md
