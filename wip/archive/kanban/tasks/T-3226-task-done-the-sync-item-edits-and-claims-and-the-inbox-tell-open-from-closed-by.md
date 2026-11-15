---
id: T-3226
type: task
nature: feature
title: task done, the sync, item edits and claims, and the inbox tell open from closed by kind
status: done
parent: S-0746
owner: alex
created: 2026-11-12T16:07:48Z
updated: 2026-11-15T00:46:34Z
transitions:
  - to: ready
    at: 2026-11-14T23:56:36Z
    by: agent-S-0746
  - to: in-progress
    at: 2026-11-14T23:56:36Z
    by: agent-S-0746
  - to: done
    at: 2026-11-15T00:25:16Z
    by: agent-S-0746
stream: S-0746
tags: [cli]
touches: [flai/internal/taskdone/taskdone.go, flai/internal/taskdone/taskdone_test.go, flai/internal/storygit/sync.go, flai/internal/itemedit/itemedit.go, flai/internal/itemedit/claim.go, flai/internal/itemedit/covers.go, flai/internal/inbox/inbox.go, flai/internal/inbox/inbox_test.go, flai/internal/itemedit/claim_test.go, flai/internal/itemedit/covers_test.go, flai/internal/storygit/sync_test.go, flai/internal/inbox/starved.go, flai/internal/itemedit/widen.go]
usage:
  source: log
  seconds: 1710
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 165
      output: 51630
      cache_read: 6359550
      cache_write: 272474
      cost: 3.8624
---
# T-3226 task done, the sync, item edits and claims, and the inbox tell open from closed by kind

## Work

Replace the references to a state's name in these files: `task done` closes a task along the task's transition to an accepted state, a story's claim counts its open tasks by kind, the sync and item edits tell an open item from a closed one by kind, and the inbox reports each move in the workflow's names. Test `task done` on a workflow with a state renamed.

## Done when

These files name no state, and `taskdone_test.go` and the existing tests pass on both workflows.

## Notes
