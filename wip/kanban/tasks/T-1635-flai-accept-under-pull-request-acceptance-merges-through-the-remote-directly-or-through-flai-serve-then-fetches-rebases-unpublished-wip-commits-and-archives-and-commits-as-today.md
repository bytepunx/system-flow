---
id: T-1635
type: task
nature: feature
title: flai accept under pull-request acceptance merges through the remote, directly or through flai serve, then fetches, rebases unpublished wip commits, and archives and commits as today
status: backlog
parent: S-0397
owner: alex
created: 2026-10-16T04:40:07Z
updated: 2026-10-16T04:40:07Z
transitions: []
stream: S-0397
tags: [cli]
touches: [flai/cmd/accept.go, flai/cmd/accept_remote.go, flai/cmd/accept_remote_test.go]
after: [T-1634]
---
# T-1635 flai accept under pull-request acceptance merges through the remote, directly or through flai serve, then fetches, rebases unpublished wip commits, and archives and commits as today

## Work

- In `flai/cmd/accept.go`, under `acceptance.mode: pull-request`, take the path in a new `accept_remote.go` in place of `mergeStoryBranch`. In order:
  1. Refuse, naming the commits, when local main has unpublished commits outside `wip/`.
  2. Merge: with the token's variable set, call T-1634's function in this process; without it, ask flai serve and wait.
  3. `git fetch`, then rebase local main's unpublished `wip/` commits onto `origin`'s main, and check that the story's commits are there.
  4. Remove the story's local branch and worktree as `mergeStoryBranch` does.
  5. Continue with done, archive, and commit as today.
- Make the dry run name the method, or say "already merged".
- Test it in `accept_remote_test.go` with a scratch bare remote and the fake forge:
  - each method
  - already merged
  - refused by the forge, with nothing local changed
  - unpublished non-`wip/` commits
  - a clean run whose acceptance commit subject and archive match a local acceptance's
- Waits for T-1634: it calls that merge.

## Done when

`flai accept` under pull-request acceptance passes each case above against the scratch remote, its dry run says what it would do, and the tests pass.

## Notes
