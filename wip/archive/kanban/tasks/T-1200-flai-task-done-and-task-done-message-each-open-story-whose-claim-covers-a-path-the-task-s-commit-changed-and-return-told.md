---
id: T-1200
type: task
nature: feature
title: flai task done and task_done message each open story whose claim covers a path the task's commit changed, and return told
status: done
parent: S-0333
owner: alex
created: 2026-10-07T20:15:38Z
updated: 2026-10-07T22:27:18Z
transitions:
  - to: ready
    at: 2026-10-07T22:13:26Z
    by: agent-S-0333
  - to: in-progress
    at: 2026-10-07T22:13:27Z
    by: agent-S-0333
  - to: done
    at: 2026-10-07T22:27:18Z
    by: agent-S-0333
stream: S-0333
tags: [flai]
touches: [flai/internal/taskdone/taskdone.go, flai/internal/taskdone/taskdone_test.go, flai/cmd/task_done.go, flai/cmd/task_done_test.go, flai/internal/mcpserver/task.go, flai/internal/mcpserver/task_test.go, flai/internal/messages/messages.go, flai/internal/messages/messages_test.go, docs/users/flai-reference.md]
after: [T-1199]
usage:
  source: log
  seconds: 831
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 95
      output: 31839
      cache_read: 5753748
      cache_write: 153848
      cost: 2.7822
---
# T-1200 flai task done and task_done message each open story whose claim covers a path the task's commit changed, and return told

## Work

Tell the other open stories at each task's close. It waits for T-1199, whose function says whom to tell.

- After the commit, take the paths it changed and, for each story T-1199's function returns, open or add to the pair's conversation, `about` those paths, naming the task, the commit, and its subject.
- A commit that changes nothing another open story claims sends nothing; a message that cannot be written is logged as a warning and the close goes on.
- `flai task done --json` and the MCP tool `task_done` return `told`: each story told and its paths. `flai task done` prints them.

## Done when

- Tests cover a story told, a shared path told, an empty claim told of every path, nobody told, and a write failure that does not fail the close.
- `flai test` passes on the changed paths.
- `flai check --strict` passes.

## Notes
