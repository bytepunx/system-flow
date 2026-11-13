---
id: T-3226
type: task
nature: feature
title: task done, the sync, item edits and claims, and the inbox tell open from closed by kind
status: backlog
parent: S-0746
owner: alex
created: 2026-11-12T16:07:48Z
updated: 2026-11-12T16:07:48Z
transitions: []
stream: S-0746
tags: [cli]
touches: [flai/internal/taskdone/taskdone.go, flai/internal/taskdone/taskdone_test.go, flai/internal/storygit/sync.go, flai/internal/itemedit/itemedit.go, flai/internal/itemedit/claim.go, flai/internal/itemedit/covers.go, flai/internal/inbox/inbox.go]
---
# T-3226 task done, the sync, item edits and claims, and the inbox tell open from closed by kind

## Work

Replace the references to a state's name in these files: `task done` closes a task along the task's transition to an accepted state, a story's claim counts its open tasks by kind, the sync and item edits tell an open item from a closed one by kind, and the inbox reports each move in the workflow's names. Test `task done` on a workflow with a state renamed.

## Done when

These files name no state, and `taskdone_test.go` and the existing tests pass on both workflows.

## Notes
