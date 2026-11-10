---
id: S-0675
type: story
nature: improvement
title: The merge queue builds one-story batches and runs three verifications per story landed
status: ready
owner: alex
created: 2026-11-10T08:34:53Z
updated: 2026-11-10T19:26:06Z
transitions:
  - to: ready
    at: 2026-11-10T10:57:32Z
    by: alex
tags: [flai, merge-queue]
topics: [orchestration, metrics]
touches: [design/adrs, design/adrs/README.md, flai/internal/mergequeue/land.go, flai/internal/mergequeue/land_test.go, flai/cmd/accept_queue_test.go, flai/internal/mergequeue/log.go, flai/internal/mergequeue/log_test.go, flai/internal/metrics/queue.go, flai/internal/metrics/queue_test.go, flai/cmd/stats.go, design/system/workflow.md, design/system/metrics.md, design/system/flai-cli.md, docs/users/flai.md, template/CHANGELOG.md, design/issues/I-0287-the-merge-queue-builds-one-story-batches-and-runs-three-verifications-per-story.md, design/issues/summary.md]
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
      seconds: 197
      estimated: true
      models:
        - model: claude-haiku-5-5
          input: 47
          output: 24954
          cache_read: 4068103
          cache_write: 80647
          cost: 0.2921
        - model: claude-opus-5-5
          input: 237
          output: 51423
          cache_read: 19053967
          cache_write: 337286
          cost: 7.7275
    - kind: analyzer
      seconds: 42
      estimated: true
      models:
        - model: claude-fable-5-1
          input: 32
          output: 3103
          cache_read: 238317
          cache_write: 16129
          cost: 0.5376
        - model: claude-opus-5-5
          input: 15
          output: 2841
          cache_read: 759050
          cache_write: 13854
          cost: 0.278
cost_of_delay:
  inputs:
    time_lost_per_cycle: 7h
    by: flai
    at: 2026-11-10T08:34:53Z
  value: 1050
  by: planner-S-0675
  at: 2026-11-10T10:35:58Z
forecast:
  duration: 1h30m
  delivery: 2026-11-10T22:13:00Z
  basis: "Its own forecast of 1h30m; 4th in the pull order with an in-progress limit of 3, behind S-0666, S-0667, S-0708, S-0688, S-0677 and S-0679."
  by: flai
  at: 2026-11-10T19:26:06Z
finalized:
  by: alex
  at: 2026-11-10T10:57:29Z
---
# S-0675 The merge queue builds one-story batches and runs three verifications per story landed

## Goal

This story remediates [I-0287](../../../design/issues/I-0287-the-merge-queue-builds-one-story-batches-and-runs-three-verifications-per-story.md), "The merge queue builds one-story batches and runs three verifications per story landed", which the analysis in [2026-11-07-all.md](../../../design/analysis/2026-11-07-all.md) found. The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0287 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0287 is closed with `flai issue close I-0287 --reason` saying what fixed it

## Tasks
- T-3073 Record in an ADR that a verified batch lands over main's wip-only commits and that every drop is logged and counted
- T-3074 A verified batch lands over main's commits that change only wip paths none of its stories change
- T-3075 The queue logs each dropped batch with its reason, and flai stats counts drops per week
- T-3076 The design and the user guide say which moves of main drop a verified batch and that flai stats counts drops
- T-3077 Close I-0287 with what fixed it

## Notes

The analyzer found I-0287: the analysis the goal links gives its evidence.

Cost of delay inputs set by flai from I-0287. time_lost_per_cycle 7h carried over from I-0287's Impact section.

Strategic usage carried over from I-0287: analyzer 42s, 0.8156 USD.

### Planning

Planned by planner-S-0675 on 2026-11-10.

#### The cause, from the queue's own log

`.flai-cache/merge-queue/events.jsonl` (2026-10-21 to 2026-11-02) holds 129 `batch_built` lines, 124 of one story and 5 of two. It also holds 94 `verified` lines that passed, and 37 of those batches never landed. A batch that passed and did not land was dropped after it had spent its whole verification. The log records no line for a drop, so drops show in `flai stats` only as `verifications` over `landed`.

`mergequeue.DropIfMainMoved` (`flai/internal/mergequeue/land.go`) drops a verified batch when the main branch's head is not exactly the commit the batch expects. The landing's own check in `Landable` (same file, "main moved other than by the batch's landing") does the same between two landings. Neither looks at what moved main. In the three latest drops sampled (batches 20261101-234541, 20261101-234945, and 20261102-051431), main moved by one commit, flai's `chore: replan forecasts after reordered`, which changed only 20 to 26 files under `wip/`. Each such commit cost a full verification of the `integration` and `smoke` tiers, about ten minutes, and the rebuilt batch verified the same story again.

`flai accept --batch` rebases the story's branch onto main before its `git merge --ff-only` (`mergeStoryBranch` in `flai/cmd/branch.go`), so a landing onto a main that gained only `wip` commits no story of the batch changes is clean. What the batch verified outside `wip/` is unchanged. `flai verify`'s sync step passes over the same kind of commit for the same reason (ADR-0135, `verify.storyRun.sync`).

#### The proposed fix

- A verified batch is not dropped when main moved from the commit the batch expects only by commits that descend from it and that change only paths under the manifest's `wip` folder that no story of the batch changes. The batch lands on the new main without being verified again. The same rule applies between two landings. A commit outside `wip/`, a `wip` path a story of the batch changes, or a main that does not descend from the expected commit still drops the batch, as today.
- Every drop appends a `dropped` line to the queue's event log, with the batch, its stories, and the drop's reason. `flai stats` counts `dropped` per week in `queue.weeks[]`, so the next analysis can tell drops from bisections.
- One ADR records both decisions. It refines ADR-0173, ADR-0177, and ADR-0227. The metrics change needs the ADR, since `metrics.md` changes only with one.

One-story batches are not changed: `merge_queue.coalesce` (ADR-0226) already joins the stories enqueued within 3 minutes of each other, and stories reached review hours apart (review mean 16,227 s in W44). The waste I-0287 measures is the 2.95 verifications per landed story, and dropped batches are most of it. The plan's thread asks whether the operator wants the coalesce window revisited too.

#### Touches

| Touch | Source |
|-------|--------|
| `design/adrs` | Design: the new ADR. Kept as a folder because the ADR's file name is not known until `flai adr new` writes it, and no task can name it yet |
| `design/adrs/README.md` | Design: the index `flai adr new` regenerates |
| `flai/internal/mergequeue/land.go`, `land_test.go` | Layout: `DropIfMainMoved`, `Landable`, and their tests (`TestMainMovingDropsAVerifiedBatchWhichIsBuiltAgain`, `TestMainMovingBetweenLandingsDropsTheBatch`) |
| `flai/cmd/accept_queue_test.go` | Layout: `flai accept --batch` landing onto a main that gained a `wip`-only commit |
| `flai/internal/mergequeue/log.go`, `log_test.go` | Layout: the event log's kinds and the events derived from a state change |
| `flai/internal/metrics/queue.go`, `queue_test.go` | Layout: `queue.weeks[]` in `flai stats` |
| `design/system/workflow.md` | Design: § What waits on an acceptance, the "Main moving" bullet |
| `design/system/metrics.md` | Design: § Merge queue, the `queue.weeks[]` rows |
| `design/system/flai-cli.md` | Design: § The merge queue pass |
| `docs/users/flai.md` | Design: § The merge queue and the `flai stats` queue fields |
| `template/CHANGELOG.md` | Layout: the changelog this repository keeps for released behaviour |
| `design/issues/I-0287-…md`, `design/issues/summary.md` | Declared by criterion 2: `flai issue close` |

The story declared no touches, so `flai touches suggest` was run on the touches above. It listed only hub files that change with most commits: `docs/operators/index.md` (10%), `docs/users/flai-reference.md` (9%), and `design/system/flaiover-dashboard.md` (7%). None is added: the fix adds no setting, flag, or dashboard field. The story's agent widens the touches if `docs/operators/index.md` § The merge-queue host action turns out to describe the drop.

#### Figures

- Forecast: `flai forecast` gave 31m, 103 s per unit of size over 155 done large improvement stories, times size 18 (2 criteria, 16 touches). Raised to 1h30m. The story writes an ADR and real-git tests of a landing. Its close-out also runs the `integration` and `smoke` tiers, and here one run of each takes over ten minutes. Delivery moves by the hour added, to 2026-11-10T15:47:00Z.
- Cost of delay: `flai cod` gave 1050 USD a week, from 7h of time lost per 168h cycle at 150 USD an hour. That stands. The 7h input is the operator's, carried from I-0287's Impact. It matches the 37 verifications beyond one per landed story in W44, at about 12 minutes each.
- No `after`: the story reads nothing another open story writes. S-0698 also changes `flai/internal/mergequeue` (how the queue reads heads), and a trial merge judges that overlap.
