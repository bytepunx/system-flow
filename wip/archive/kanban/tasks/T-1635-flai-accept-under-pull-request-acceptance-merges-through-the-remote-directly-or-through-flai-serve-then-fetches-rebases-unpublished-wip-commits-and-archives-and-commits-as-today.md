---
id: T-1635
type: task
nature: feature
title: flai accept under pull-request acceptance merges through the remote, directly or through flai serve, then fetches, rebases unpublished wip commits, and archives and commits as today
status: done
parent: S-0397
owner: alex
created: 2026-10-16T04:40:07Z
updated: 2026-10-18T03:06:42Z
transitions:
  - to: ready
    at: 2026-10-18T02:42:38Z
    by: agent-S-0397
  - to: in-progress
    at: 2026-10-18T02:42:38Z
    by: agent-S-0397
  - to: done
    at: 2026-10-18T03:06:42Z
    by: agent-S-0397
stream: S-0397
tags: [cli]
touches: [flai/cmd/accept.go, flai/cmd/accept_remote.go, flai/cmd/accept_remote_test.go, flai/cmd/branch.go, flai/internal/preview/accept.go, flai/internal/preview/accept_test.go, flai/cmd/accept_forge_test.go, flai/cmd/forge_test.go]
after: [T-1634]
usage:
  source: log
  seconds: 1444
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 114
      output: 52066
      cache_read: 7385442
      cache_write: 289611
      cost: 4.1503
---
# T-1635 flai accept under pull-request acceptance merges through the remote, directly or through flai serve, then fetches, rebases unpublished wip commits, and archives and commits as today

## Work

- In `flai/cmd/accept.go`, under `acceptance.mode: pull-request`, take a new path in `accept_remote.go` instead of `mergeStoryBranch`. In order, after S-0395's gate:
  1. **Preflight, before anything changes.**
     - Refuse while the `forge` host action is off for the project, naming `flai serve enable forge`.
     - Refuse when origin is missing.
     - `git fetch origin <main>`. Refuse, naming each commit, when local main has unpublished commits outside the wip folder.
  2. **Merge.**
     - With the forge token's variable readable in this process, call T-1634's `serve.MergeStory` here.
     - Otherwise ask flai serve with `serve.RequestMerge` and wait.
     - A dashboard child, `FLAI_AGENT=flaiover`, never hands off: without the token it refuses, naming the variable.
     - Pass the head the gate passed at. A refusal stops the acceptance with the forge's reason, and nothing local has changed.
  3. **Local main.**
     - `git fetch origin <main>`.
     - With no unpublished commits, main's new head is `origin/<main>`. Otherwise, rebase the unpublished wip-only commits onto it in a temporary detached worktree under `.flai-cache`, then remove that worktree.
     - In the main checkout, `git reset --keep <new head>`, which keeps uncommitted wip edits.
     - A rebase conflict refuses, naming the paths and saying that the pull request is merged and a second `flai accept` makes no second merge.
  4. **Story branch.** Remove the story's worktree as `mergeStoryBranch` does, and its branch with `git branch -D`.
  5. **Finish as today.** Done, archive, and commit, with the same acceptance commit subject and archive as under `local`.
- **Gate.** In `flai/internal/preview/accept.go`, S-0395's gate passes over a kept status whose `State` is merged.
- **Dry run.** Under pull-request acceptance it says either `would merge pull request #N through the remote by <method>` (through flai serve when the token is not at hand), or that the pull request is merged on the forge already and would be accepted with no second merge. It reads the kept state and calls nothing. `--json` carries the same.
- **Tests.** In `accept_remote_test.go`, with a scratch bare origin and a fake forge or `pullsStandIn`:
  - a merge by each method
  - already merged
  - refused by the forge, with nothing local changed
  - unpublished commits outside wip refused before the merge
  - unpublished wip commits rebased, with uncommitted wip edits kept
  - the host action off
  - the hand-off answered by a fake merge pass
  - the dry run
  - a clean run whose acceptance commit subject and archive match a local acceptance's
- **Waits for T-1634:** it calls `MergeStory` and `RequestMerge`.

## Done when

`flai accept` under pull-request acceptance passes each case above against the scratch remote. Its dry run says what it would do. The cmd and preview tests pass.

## Notes
