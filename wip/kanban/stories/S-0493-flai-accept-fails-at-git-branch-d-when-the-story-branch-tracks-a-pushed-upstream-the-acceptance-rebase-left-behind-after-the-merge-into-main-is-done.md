---
id: S-0493
type: story
nature: remediation
title: flai accept fails at git branch -d when the story branch tracks a pushed upstream the acceptance rebase left behind, after the merge into main is done
status: in-progress
owner: alex
created: 2026-10-24T13:15:31Z
updated: 2026-10-24T13:26:57Z
transitions:
  - to: ready
    at: 2026-10-24T13:15:54Z
    by: system-flow
  - to: in-progress
    at: 2026-10-24T13:16:08Z
    by: agent-S-0493
tags: [flai]
topics: [cli, git]
touches: [flai/cmd/branch.go, flai/cmd/accept.go, flai/cmd/accept_test.go, design/system/flai-cli.md, docs/users/flai.md, design/issues/I-0204-flai-accept-fails-at-git-branch-d-when-the-story-branch-tracks-a-pushed-upstream-the-acceptance-rebase-left-behind-after-the-merge-into-main-is-done.md, design/issues/summary.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 659
  estimated: true
  priced_by: estimate
  turns:
    - day: 2026-10-24
      test_runs: 1
      hand_edits: 1
      work: 26
  models:
    - model: claude-opus-5-5
      input: 56
      output: 515
      cache_read: 3560699
      cache_write: 149557
      cost: 1.2471
  strategic:
    - kind: orchestrator
      seconds: 2311
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 94
          output: 1518
          cache_read: 10153951
          cache_write: 24786
          cost: 2.2318
cost_of_delay:
  inputs:
    time_lost_per_cycle: 10m
    by: flai
    at: 2026-10-24T13:15:31Z
  value: 25
  by: planner-S-0493
  at: 2026-10-24T13:19:13Z
forecast:
  duration: 30m
  delivery: 2026-10-24T13:47:00Z
  basis: "flai forecast gave 18m from a median of 117 s per unit of size over 30 done remediation stories on claude-opus-5-5 at size 9; raised to 30m because the reproducing tests need a bare remote and an upstream-tracking story branch that no acceptance test sets up today, and two long documents change; delivery counts from the start at 2026-10-24T13:17Z."
  by: planner-S-0493
  at: 2026-10-24T13:19:13Z
finalized:
  by: system-flow
  at: 2026-10-24T13:15:43Z
---
# S-0493 flai accept fails at git branch -d when the story branch tracks a pushed upstream the acceptance rebase left behind, after the merge into main is done

## Goal

This story remediates [I-0204](../../../design/issues/I-0204-flai-accept-fails-at-git-branch-d-when-the-story-branch-tracks-a-pushed-upstream-the-acceptance-rebase-left-behind-after-the-merge-into-main-is-done.md), "flai accept fails at git branch -d when the story branch tracks a pushed upstream the acceptance rebase left behind, after the merge into main is done". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [x] The cause I-0204 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0204 is closed with `flai issue close I-0204 --reason` saying what fixed it

## Tasks
- T-2105 flai accept deletes the story branch with git branch -D after its fast-forward, so a pushed upstream the rebase left behind no longer stops it
- T-2106 I-0204 is closed with flai issue close naming the fix, and the story's criteria are ticked

## Notes

Cost of delay inputs set by flai from I-0204. time_lost_per_cycle 10m carried over from I-0204's Impact section. Its Impact time_lost_per_cycle was taken rather than the 10m derived from its cost and count.

### Planning

- Touches: the seven the issue step declared from I-0204 are kept (declared). `flai touches suggest` lists only high-churn documents (`docs/operators/index.md`, `docs/users/flai-reference.md`, `design/system/flaiover-dashboard.md`, each at 12% or less) that the fix does not reach, so none is added (co-change). The code layout confirms the set and adds nothing: `mergeStoryBranch` is in `flai/cmd/branch.go`, the acceptance flow and its `merged` step in `flai/cmd/accept.go`, the real-git acceptance tests in `flai/cmd/accept_test.go`, and the branch's removal at acceptance is described in `design/system/flai-cli.md` (the `flai accept` row under Commands) and `docs/users/flai.md` (Branches, and Acceptance) (layout). The issue file and `design/issues/summary.md` change when I-0204 is closed (design). No folder touch.
- Forecast: flai gave 18m, a median of 117 s per unit of size over 30 done remediation stories on claude-opus-5-5 at size 9 (2 criteria, 7 touches). Set to 30m: the reproducing tests need a bare remote and an upstream-tracking story branch that no acceptance test sets up today, there are two of them (the acceptance, and the acceptance run again from the state I-0204 left), and two long documents change. Delivery counts from the start at 2026-10-24T13:17Z.
- Cost of delay: 25.00 USD a week stands as flai worked it out from the operator's input, 10m lost per 168h cycle at 150 USD an hour, one cycle a week. I-0204 records one occurrence, so there is no evidence for more cycles a week.
