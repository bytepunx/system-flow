---
id: T-1107
type: task
nature: improvement
title: The MCP tool task_done closes a task with the same answer as flai task done --json
status: done
parent: S-0269
owner: alex
created: 2026-10-06T22:53:26Z
updated: 2026-10-07T02:58:33Z
transitions:
  - to: ready
    at: 2026-10-07T02:14:03Z
    by: agent-S-0269
  - to: in-progress
    at: 2026-10-07T02:14:03Z
    by: agent-S-0269
  - to: done
    at: 2026-10-07T02:58:33Z
    by: agent-S-0269
stream: S-0269
tags: [mcp]
touches: [flai/internal/mcpserver/task.go, flai/internal/mcpserver/task_test.go, flai/internal/mcpserver/folder.go, flai/internal/mcpserver/server.go, flai/internal/mcpserver/server_test.go, flai/internal/mcpserver/folder_test.go]
after: [T-1098]
usage:
  source: log
  seconds: 353
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 52
      output: 22467
      cache_read: 3024694
      cache_write: 88469
      cost: 1.5702
---
# T-1107 The MCP tool task_done closes a task with the same answer as flai task done --json

## Work

Add the operation over MCP.

- In `flai/internal/mcpserver/task.go`, add `TaskDoneIn` and `TaskDoneOut`, and `(*server).taskDone`. The input is the task, the message, and an optional log entry and project. The method calls `taskdone.Run` (T-1098) with the server's agent and runner, and answers its result as the output, field for field the same as `flai task done --json`.
- Register `task_done` in `flai/internal/mcpserver/folder.go` beside `criteria_tick`. Its description names the steps, their order, and the stop rule.
- In `flai/internal/mcpserver/server.go`, the server's instructions send the agent to `task_done` at every task transition, in place of commit, sync, move, log, touches, check, and inbox one by one.
- In `flai/internal/mcpserver/task_test.go`, cover the happy path, a refused sync, and a failed check. Test as well that the output is the same as the result `flai task done --json` prints for the same case.

This task waits for T-1098, whose `Run` it calls. It is layer 3, beside the CLI command, with which it shares no file.

## Done when

- [ ] `task_done` over MCP closes a task and answers the same result as `flai task done --json`.
- [ ] The server's instructions name it at the task transition.
- [ ] Tests cover the happy path, a refused sync, and a failed check, and `scripts/flai-test.sh` passes.

## Notes

Drafted by the planner.
