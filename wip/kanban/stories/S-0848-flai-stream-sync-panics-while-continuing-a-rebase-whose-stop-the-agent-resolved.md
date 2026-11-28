---
id: S-0848
type: story
nature: remediation
title: flai stream sync panics while continuing a rebase whose stop the agent resolved, and leaves the rebase stopped at the next commit
status: ready
owner: alex
created: 2026-11-21T17:20:23Z
updated: 2026-11-28T00:01:28Z
transitions:
  - to: ready
    at: 2026-11-28T00:01:13Z
    by: orchestrator
tags: []
topics: [git, sync]
touches: [flai/internal/storygit/mdmerge.go, flai/internal/storygit/mdmerge_test.go, flai/internal/storygit/sync.go, flai/internal/storygit/sync_test.go, design/issues/I-0407-flai-stream-sync-panics-while-continuing-a-rebase-whose-stop-the-agent-resolved.md, design/issues/summary.md]
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
      seconds: 45
      estimated: true
      models:
        - model: claude-haiku-5-5
          input: 2
          output: 3
          cache_read: 131146
          cache_write: 362
          cost: 0.009
        - model: claude-opus-5-5
          input: 61
          output: 1107
          cache_read: 3714378
          cache_write: 60380
          cost: 1.0516
    - kind: orchestrator
      seconds: 10
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 1
          output: 2
          cache_read: 83340
          cache_write: 1445
          cost: 0.027
cost_of_delay:
  inputs:
    time_lost_per_cycle: 5m
    by: flai
    at: 2026-11-21T17:20:23Z
  value: 12.5
  by: planner-S-0848
  at: 2026-11-27T23:57:18Z
forecast:
  duration: 1h30m
  delivery: 2026-11-28T03:55:00Z
  basis: "Its own forecast of 1h30m; 5th in the pull order with an in-progress limit of 3, behind S-0837, S-0844, S-0846, S-0835, S-0814, S-0813 and S-0822."
  by: flai
  at: 2026-11-28T00:01:28Z
finalized:
  by: orchestrator
  at: 2026-11-28T00:00:02Z
---
# S-0848 flai stream sync panics while continuing a rebase whose stop the agent resolved, and leaves the rebase stopped at the next commit

## Goal

This story remediates [I-0407](../../../design/issues/I-0407-flai-stream-sync-panics-while-continuing-a-rebase-whose-stop-the-agent-resolved.md), "flai stream sync panics while continuing a rebase whose stop the agent resolved, and leaves the rebase stopped at the next commit". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0407 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0407 is closed with `flai issue close I-0407 --reason` saying what fixed it

## Tasks
- T-3615 Find the input that panics the section merge at S-0809's second stop and fix it, with a test that reproduces it
- T-3616 A resolver that panics at a stop leaves the stop waiting with its paths, so flai stream sync reports them instead of crashing
- T-3617 Close I-0407 saying what fixed it

## Notes

Cost of delay inputs set by flai from I-0407. time_lost_per_cycle 5m: 5m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-11-18T13:08:06Z, 3.2 days before this story; under one cycle counts as one).

### Planning

Proposed cause, from I-0407's instance. The continue of the agent's resolved stop (ADR-0294, `resume` in `flai/internal/storygit/sync.go`) went on and stopped at the next commit, on `design/system/flai-cli.md` and `design/system/workflow.md`. Both are markdown under the design folder, so `settleStops` and `ContinueOverResolved` merged them by sections with `storygit.MergeMarkdown` (ADR-0276). The panic most likely came from that merge, under the `Rebase` frame the trace showed. The word merge S-0817 added to `flai/internal/storygit/mdmerge.go` (ADR-0314) landed two days before the instance. The panic stopped the sync before `stopped` could report the stop, which is why the agent saw a crash and not the two paths. The fix has two parts: T-3615 finds and fixes the input that panics, and T-3616 makes any resolver's panic leave the stop waiting with its paths.

Touches, each a file:

- `flai/internal/storygit/mdmerge.go` and `mdmerge_test.go`: design (ADR-0276, ADR-0314) and layout, the section merge that the second stop ran; T-3615.
- `flai/internal/storygit/sync.go` and `sync_test.go`: layout, `resume`, `settleStops`, and `ContinueOverResolved`, from the trace's `storygit.Sync` frame; T-3616.
- `design/issues/I-0407-…md` and `design/issues/summary.md`: from the criterion that closes I-0407; T-3617.
- `flai touches suggest` lists `design/system/flai-cli.md`, `docs/users/flai.md`, `flai/cmd/stream.go`, and `flai/cmd/stream_sync_test.go` as co-changes. They are left out because the fix changes no documented behaviour: ADR-0175 and ADR-0276 already say that a merge that does not settle a file leaves the stop waiting. The story's agent adds them if its fix changes what `flai stream sync` prints.

Topics `git` and `sync` were added for the rebase and the sync this story changes.

No `after`: the story reads nothing that an open story writes.

Forecast: `flai forecast` gave 21m, from a median of 151 s per unit over size 8. It is raised to 1h30m, because the panic's first frame was not recorded. T-3615 must rebuild the stop's three versions from git history, or search for a panicking input, before it can fix anything. That search is the bulk of the work and is not in the size. The delivery moves later by the same 69m, from `flai forecast`'s 2026-11-28T05:03Z to 06:15Z, since the story is 9th in the pull order.

Cost of delay: `flai cod` gives 12.50 USD a week from the operator's input of 5m lost per 168h cycle at 150 USD an hour. It stands: one instance, and the second sync recovered it.
