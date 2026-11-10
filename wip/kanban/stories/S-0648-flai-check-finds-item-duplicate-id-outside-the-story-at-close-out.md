---
id: S-0648
type: story
nature: improvement
title: "flai check finds `item.duplicate-id` outside the story at close-out"
status: review
owner: alex
created: 2026-11-10T05:12:06Z
updated: 2026-11-10T08:43:34Z
transitions:
  - to: ready
    at: 2026-11-10T07:44:17Z
    by: alex
  - to: in-progress
    at: 2026-11-10T08:08:31Z
    by: agent-S-0648
  - to: review
    at: 2026-11-10T08:43:34Z
    by: agent-S-0648
tags: [flai, workitem]
topics: [code, archive, usage]
touches: [flai/internal/workitem/itemlock.go, flai/internal/workitem/itemlock_test.go, flai/internal/workitem/archive.go, flai/internal/workitem/archive_test.go, flai/internal/serve/usage.go, flai/internal/serve/usage_test.go, wip/kanban/tasks/T-2760-fetch-the-remote-s-branch-and-tags.md, design/system/work-hierarchy.md, design/issues/I-0323-flai-check-finds-item-duplicate-id-outside-the-story-at-close-out.md, design/issues/summary.md, design/adrs/0277-repo-save-refuses-to-write-an-active-work-item-whose-file-was-archived-since-it.md, design/adrs/README.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 2114
  priced_by: harness
  turns:
    - day: 2026-11-10
      hand_edits: 3
      work: 67
  models:
    - model: claude-opus-5-5
      input: 252
      output: 61610
      cache_read: 22716962
      cache_write: 580966
      cost: 9.6459
  strategic:
    - kind: planner
      seconds: 419
      estimated: true
      models:
        - model: claude-haiku-5-5
          input: 66
          output: 26729
          cache_read: 2684738
          cache_write: 127559
          cost: 0.1598
        - model: claude-opus-5-5
          input: 76
          output: 19981
          cache_read: 5168314
          cache_write: 158018
          cost: 2.6977
    - kind: orchestrator
      seconds: 299
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 75
          output: 559
          cache_read: 2640155
          cache_write: 340991
          cost: 0.8007
cost_of_delay:
  inputs:
    time_lost_per_cycle: 2m
    by: orchestrator
    at: 2026-11-10T06:33:20Z
  value: 5
  by: planner-S-0648
  at: 2026-11-10T06:36:35Z
forecast:
  duration: 1h
  delivery: 2026-11-10T09:12:00Z
  basis: "Its own forecast of 1h; 2nd in the pull order with an in-progress limit of 5, behind S-0636, S-0637, S-0638, S-0651, S-0654 and S-0652."
  by: flai
  at: 2026-11-10T07:57:16Z
finalized:
  by: alex
  at: 2026-11-10T07:44:12Z
---
# S-0648 flai check finds `item.duplicate-id` outside the story at close-out

## Goal

This story remediates [I-0323](../../../design/issues/I-0323-flai-check-finds-item-duplicate-id-outside-the-story-at-close-out.md), "flai check finds `item.duplicate-id` outside the story at close-out". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [x] The cause I-0323 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0323 is closed with `flai issue close I-0323 --reason` saying what fixed it

## Tasks
- T-2827 Repo.Save refuses to recreate an item file moved since it was read, and Archive moves under the items lock
- T-2828 Measure writes a usage to the item where it now is when the item was archived since it was read
- T-2829 Remove the stray open copy of T-2760 that the archive race left in wip/kanban/tasks
- T-2830 Close I-0323 with the fix as its reason

## Notes

### Planning

#### The cause, from the instances

Every instance is the same file: `wip/kanban/tasks/T-2760-fetch-the-remote-s-branch-and-tags.md`, a task of S-0642. S-0642's acceptance commit, 436e8f612, adds the archived copy and also modifies the open copy, where its other tasks are renamed cleanly. The two copies differ only in `usage`: the open copy is `priced_by: harness`, cost 0.3414, and the archived copy `priced_by: estimate`, cost 0.2198. Both carry the same `updated`, so the open copy was written by a usage measurement, not by a move.

Read from the code by the explorer, not reproduced:

- `serve.Measure`'s `put` (`flai/internal/serve/usage.go`) reads the task and then saves its usage with `Repo.Save`.
- `Archive` (`flai/internal/workitem/archive.go`) renames the file into `wip/archive` without the items lock, so the rename can land between that read and that save.
- `saveLocked` and `mergeSaved` (`flai/internal/workitem/itemlock.go`) then find no file at the item's path. They take that as "nothing saved since", and write the whole item there, which recreates the file.

#### Proposed fix

1. `Repo.Save` never recreates an item file that was there when the item was read and is gone now. It refuses with an error that says the item moved, and `Archive` takes the items lock for its moves.
2. `Measure` reads the item again where it now is and writes its usage there, its archived copy included, so the harness's figure is not lost.
3. The stray open copy of T-2760 is removed. Its archived copy stays as the acceptance wrote it, since `wip/archive` is not edited.

#### Touches

| Touch | Source |
|-------|--------|
| `flai/internal/workitem/itemlock.go` | layout: `saveLocked`, `mergeSaved` |
| `flai/internal/workitem/itemlock_test.go` | co-change (4 of 14 commits) |
| `flai/internal/workitem/archive.go` | layout: `Archive`, `moveFile` |
| `flai/internal/workitem/archive_test.go` | layout: the archive tests |
| `flai/internal/serve/usage.go` | layout: `Measure`, `put` |
| `flai/internal/serve/usage_test.go` | co-change (5 of 14 commits) |
| `wip/kanban/tasks/T-2760-fetch-the-remote-s-branch-and-tags.md` | the issue's instances: the stray copy |
| `design/system/work-hierarchy.md` | design: the ADR-0184 rule on concurrent writers |
| `design/issues/I-0323-flai-check-finds-item-duplicate-id-outside-the-story-at-close-out.md`, `design/issues/summary.md` | criterion 2 |

`flai touches suggest` also listed `flai/internal/workitem/repo.go`, `flai/internal/serve/agents.go`, `flai/internal/usage/log.go`, and `design/system/flai-cli.md`. They are left out: the fix needs none of them. `flai/cmd/accept.go` calls `Archive` after a `Save`. It is touched only if taking the lock inside `Archive` would deadlock there, and `flai task done` widens the touches if so. No folder touch, no `after`.

#### Figures

- Forecast: `flai forecast` gave 22m from size 12 (2 criteria, 10 touches; 107 s per unit over 30 medium improvement stories). It is raised to 1h for four tasks, two race tests, and a close-out whose integration tier alone took 13m to 17m on S-0640, S-0625, and S-0630. The delivery is flai's 07:48Z moved by the 38m added, 08:26Z.
- Cost of delay: 5.00 USD a week from `flai cod`, from the input `time_lost_per_cycle: 2m` the orchestrator set on TH-0967. That was the planner's recommendation: the note stops no close-out and costs its reading and the issue churn. The value stands as `flai cod` gives it.
