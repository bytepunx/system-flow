---
id: T-0984
type: task
nature: remediation
title: flai stream sync reports a conflict only on paths both stories changed, not on what main brought
status: done
parent: S-0251
owner: alex
created: 2026-10-05T05:50:45Z
updated: 2026-10-07T00:04:00Z
transitions:
  - to: ready
    at: 2026-10-07T00:00:07Z
    by: agent-S-0251
  - to: in-progress
    at: 2026-10-07T00:00:07Z
    by: agent-S-0251
  - to: done
    at: 2026-10-07T00:04:00Z
    by: agent-S-0251
stream: S-0251
tags: [flai, git]
touches: [flai/cmd/stream_sync.go, flai/cmd/stream_sync_test.go, flai/cmd/stream.go, docs/users/flai-reference.md]
usage:
  source: log
  seconds: 233
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 47
      output: 13806
      cache_read: 1753042
      cache_write: 80156
      cost: 1.1773
    - model: claude-sonnet-5-5
      input: 18
      output: 3356
      cache_read: 277169
      cache_write: 44481
      cost: 0.2002
---
# T-0984 flai stream sync reports a conflict only on paths both stories changed, not on what main brought

## Work

The cause, from I-0064's instances. `trialMerge` (`flai/cmd/stream_sync.go`) runs `git merge-tree --write-tree --name-only --no-messages ours theirs`, and git takes the merge base itself. The syncing story has just rebased onto main. The other story's branch has not, so the base is that branch's old fork point. The syncing side then carries everything main gained since that fork point, and a path where main and the stale branch disagree is reported against the syncing story. In S-0203's instance, five such paths came from S-0200's accepted changes, though S-0203 changed none of them (TH-0088). In S-0225's instance, `template/CHANGELOG.md` was blamed: main gained an entry after both bases, the other branch added its own, and S-0225 did not change the file (TH-0104, TH-0105).

Propose the fix in the narrative's `## Decisions` before building it. The recommended fix keeps a conflicted path only when both stories changed it since they left main. The syncing story's own paths are `git diff --name-only <main>...story/<id>`, and the other story's are `git diff --name-only <main>...story/<other>`. Both diffs are taken from each branch's merge base with main, so neither counts what main brought. A path conflicting between main and the other branch alone is left to that story's own sync to find. A trial merge left with no paths is clean. An existing conflict thread then resolves through `reportConflicts`, as it does today.

The alternative is to trial-merge the other branch onto main first, then merge that result with ours, using main as the base. It is more exact for a path both stories touch, but it costs a second `merge-tree` per story. Choose it if the intersection misses a case the tests show.

Write the reproduction tests first, with the `syncProject` and `openSyncStory` fixtures in `flai/cmd/stream_sync_test.go`.

Then say in `flai stream sync`'s help (`flai/cmd/stream.go`) that a conflict names only paths both branches changed, and regenerate the reference with `make flai-reference`.

It waits for no task.

## Done when

- A test reproduces S-0203's instance. Story A and story B branch from main. Main then gains a change to a path that B also changes, so the two conflict. A changes another path and syncs. A's sync reports B's branch clean and opens no thread. The test fails without the change.
- A test reproduces S-0225's instance. Main appends a line to a changelog after both branches' bases, and B appends its own line. A does not touch the file. A's sync reports B's branch clean.
- `TestSyncTrialMergesOtherOpenBranches` and `TestSyncConflictIsAThreadInEveryInbox` still pass: a path both stories change in conflicting ways is still reported and still opens a thread.
- `docs/users/flai-reference.md` matches the help text, and `scripts/flai-test.sh` passes.

## Notes
