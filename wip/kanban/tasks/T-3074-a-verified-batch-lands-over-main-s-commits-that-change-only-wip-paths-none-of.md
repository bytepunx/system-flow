---
id: T-3074
type: task
nature: improvement
title: A verified batch lands over main's commits that change only wip paths none of its stories change
status: done
parent: S-0675
owner: alex
created: 2026-11-10T10:36:36Z
updated: 2026-11-10T20:10:00Z
transitions:
  - to: ready
    at: 2026-11-10T19:54:43Z
    by: system-flow
  - to: in-progress
    at: 2026-11-10T19:54:43Z
    by: system-flow
  - to: done
    at: 2026-11-10T20:10:00Z
    by: agent-S-0675
stream: S-0675
tags: [flai, merge-queue]
touches: [flai/internal/mergequeue/land.go, flai/internal/mergequeue/land_test.go, flai/cmd/accept_queue_test.go]
after: [T-3073]
usage:
  source: log
  seconds: 917
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 71
      output: 18905
      cache_read: 7259489
      cache_write: 175710
      cost: 2.9347
---
# T-3074 A verified batch lands over main's commits that change only wip paths none of its stories change

## Work

Build T-3073's first decision in `flai/internal/mergequeue/land.go`. It waits for T-3073 because it builds the rule that ADR decides.

- `DropIfMainMoved` keeps a verified batch when main's head is not `expectedMain(b)` but descends from it, and the commits between change only paths under the manifest's `wip` folder that no story of the batch changes. Read the paths with `verify.BetweenPaths` (or `git diff --name-only`), and check descent with `git merge-base --is-ancestor`. The `wip` prefix comes from the manifest, as `verify.storyRun.wipPrefix` reads it. Each story's own changes are what its verified head changes from the batch's base. Otherwise it drops the batch as now, and its reason names the commits or paths that dropped it.
- `Landable`'s main check, "main moved other than by the batch's landing", passes over the same commits between two landings.
- Keep the rest of `Landable` as it is: a story branch that moved since the batch verified it still refuses.

Reproduce I-0287 first, in `land_test.go` beside `TestMainMovingDropsAVerifiedBatchWhichIsBuiltAgain` and `TestMainMovingBetweenLandingsDropsTheBatch`. Make the test fail on today's code: build and verify a one-story batch, then commit a change to a `wip/` file the story does not change, as flai's `chore: replan forecasts after reordered` does. Then make it pass. In `flai/cmd/accept_queue_test.go`, show that `flai accept --batch` lands that story onto the moved main with a clean rebase. Its `mergeStoryBranch` already rebases before its `--ff-only` merge.

## Done when

- A test reproduces I-0287: a verified batch with a `wip`-only commit on main after it is not dropped, and it lands with no second verification. It fails before the change and passes after.
- Tests show that a batch is still dropped when main gains a commit outside `wip/`, a commit to a `wip` path a story of the batch changes, or a head that does not descend from the expected commit.
- A test shows the same pass-over between two landings of one batch.
- `flai test` passes on the changed files.

## Notes
