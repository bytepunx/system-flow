---
id: S-0251
type: story
nature: improvement
title: flai stream sync's trial merge blames a story for conflicts between main and another story's stale branch
status: backlog
owner: alex
created: 2026-10-03T18:33:16Z
updated: 2026-10-06T23:17:28Z
transitions: []
tags: []
topics: [cli]
touches: [flai/cmd/stream_sync.go, flai/cmd/stream_sync_test.go, flai/cmd/stream.go, docs/users/flai-reference.md, docs/users/flai.md, design/system/flai-cli.md, design/issues/I-0064-flai-stream-sync-s-trial-merge-blames-a-story-for-conflicts-between-main-and-another-story-s-stale-branch.md, design/issues/summary.md]
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
      seconds: 1
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 4
          output: 33
          cache_read: 147786
          cache_write: 9614
          cost: 0.0861
cost_of_delay:
  inputs:
    time_lost_per_cycle: 5m
    by: planner-S-0251
    at: 2026-10-05T05:52:35Z
  value: 12.5
  by: planner-S-0251
  at: 2026-10-05T05:52:47Z
forecast:
  duration: 30m
  delivery: 2026-10-07T08:08:00Z
  basis: "Its own forecast of 30m; 22nd in the pull order with an in-progress limit of 3, behind S-0261, S-0300, S-0301, S-0228, S-0269, S-0270, S-0302, S-0271, S-0212, S-0213, S-0214, S-0215, S-0216, S-0232, S-0233, S-0234, S-0235, S-0236, S-0237, S-0238, S-0239, S-0241, S-0245 and S-0246."
  by: flai
  at: 2026-10-06T23:17:28Z
---
# S-0251 flai stream sync's trial merge blames a story for conflicts between main and another story's stale branch

## Goal

This story remediates [I-0064](../../../design/issues/I-0064-flai-stream-sync-s-trial-merge-blames-a-story-for-conflicts-between-main-and-another-story-s-stale-branch.md), "flai stream sync's trial merge blames a story for conflicts between main and another story's stale branch". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0064 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0064 is closed with `flai issue close I-0064 --reason` saying what fixed it

## Tasks
- T-0984 flai stream sync reports a conflict only on paths both stories changed, not on what main brought
- T-0985 The design and the users' guide say how sync's trial merge finds a conflict, and I-0064 is closed

## Notes

### Planning

Touches, and where each came from:

- `flai/cmd/stream_sync.go`: layout. `checkSync` and `trialMerge` live here, and `trialMerge` runs the two-argument `git merge-tree` whose own merge base causes I-0064.
- `flai/cmd/stream_sync_test.go`: co-change. It changed in 3 of 3 commits with `stream_sync.go` (`flai touches suggest`), and the reproduction tests go here.
- `flai/cmd/stream.go`: co-change, 3 of 3. It holds `flai stream sync`'s help text, which describes the trial merge.
- `docs/users/flai-reference.md`: layout. It is generated from the help text.
- `design/system/flai-cli.md`: design. Its `flai stream sync` row describes `checkSync`'s `git merge-tree`.
- `docs/users/flai.md`: design. Its **Conflicts** bullet of `flai stream sync` describes the trial merge.
- `design/issues/I-0064-…md` and `design/issues/summary.md`: from the criteria. `flai issue close` writes both.

Not touched: ADR-0046 and `design/system/agent-coordination.md` describe the safety net only as a trial merge, and that stays true.

Topic `cli` was added, the topic of `flai-cli.md` and `design/tech/git.md`.

Forecast: 30m, delivery 2026-10-05T23:30Z. `flai forecast` gave 15m: 86 s per unit of size times size 10, in the medium band. I raised it to 30m by comparison with S-0197 and S-0252. Both are git-fixture stories in the same kind of code, and each took about 31m (1846 s and 1877 s). This story also has to propose its fix before building it. Delivery moves 15m later to match: the story is 29th in the pull order.

Cost of delay: 12.50 USD a week, as `flai cod` works it out. The operator gave `time_lost_per_cycle: 5m` on TH-0150, the cost recorded on I-0064. That is 5m per 168h cycle at 150 USD an hour. The value stands unadjusted: the input is the operator's, and I have no evidence beyond I-0064's two instances.
