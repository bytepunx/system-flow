---
id: T-1917
type: task
nature: improvement
title: Make task_done, MCP and CLI --json, carry the embedded inbox and fit its whole answer to 40,000 bytes
status: done
parent: S-0450
owner: alex
created: 2026-10-20T19:05:36Z
updated: 2026-10-20T19:51:59Z
transitions:
  - to: ready
    at: 2026-10-20T19:15:56Z
    by: agent-S-0450
  - to: in-progress
    at: 2026-10-20T19:15:56Z
    by: agent-S-0450
  - to: done
    at: 2026-10-20T19:51:59Z
    by: agent-S-0450
stream: S-0450
tags: [flai]
touches: [flai/internal/taskdone/taskdone.go, flai/internal/taskdone/taskdone_test.go, flai/internal/mcpserver/task.go, flai/internal/mcpserver/task_test.go, flai/cmd/task_done.go, flai/cmd/task_done_test.go]
after: [T-1915]
usage:
  source: log
  seconds: 2163
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 89
      output: 26515
      cache_read: 5005364
      cache_write: 251249
      cost: 3.0604
---
# T-1917 Make task_done, MCP and CLI --json, carry the embedded inbox and fit its whole answer to 40,000 bytes

## Work

- `flai/internal/taskdone/taskdone.go`: in the inbox step, read the inbox with the story set, through T-1915's embedded read, and fit the result to `context.PartLimit` with T-1915's fit, so the MCP tool and `flai task done --json` stay equal, as `flai/internal/mcpserver/task_test.go`'s parity test requires.
- `flai/internal/mcpserver/task.go`: update `taskDoneDescription` to name the count of the other threads and that `inbox` lists them; check the outdated-flai fallback, which closes the task with the installed flai (ADR-0160), still passes its answer through.
- `flai/cmd/task_done.go`: print the count of the other threads in `printTaskDoneInbox`.
- Tests: in `flai/internal/taskdone/taskdone_test.go` and `flai/internal/mcpserver/task_test.go`, reproduce S-0398's instance, 33 threads awaiting the agent and none on the story, and check that the answer is at most 40,000 bytes, still gives `stopped`, `error`, and `left`, lists no thread of another item, and counts them.

It waits for T-1915, whose read and fit it calls. It runs beside T-1916: the two share no path.

## Done when

- `task_done` and `flai task done --json` answer at most 40,000 bytes of JSON in the reproduction, equal to each other, with the story's own threads listed and the others counted.
- The existing task_done tests pass, changed only where the answer's shape changed.
- `flai test` on the changed paths passes.

## Notes
