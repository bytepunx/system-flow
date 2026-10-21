---
id: T-1932
type: task
nature: remediation
title: Make story_start and task_done, MCP and CLI, carry only the threads awaiting the story's agent by the rule, with tests of the 2026-10-20 measurement
status: done
parent: S-0453
owner: alex
created: 2026-10-20T22:20:55Z
updated: 2026-10-21T01:20:03Z
transitions:
  - to: ready
    at: 2026-10-21T00:55:04Z
    by: system-flow
  - to: in-progress
    at: 2026-10-21T00:55:04Z
    by: system-flow
  - to: done
    at: 2026-10-21T01:20:03Z
    by: system-flow
stream: S-0453
tags: [flai]
touches: [flai/internal/mcpserver/story_start_test.go, flai/internal/mcpserver/task_test.go, flai/internal/storystart/start_test.go, flai/internal/taskdone/taskdone_test.go, flai/cmd/task_done.go, flai/cmd/task_done_test.go, flai/cmd/story_start_test.go]
after: [T-1929]
usage:
  source: log
  seconds: 1499
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 101
      output: 36570
      cache_read: 5421321
      cache_write: 259668
      cost: 3.277
---
# T-1932 Make story_start and task_done, MCP and CLI, carry only the threads awaiting the story's agent by the rule, with tests of the 2026-10-20 measurement

## Work

`storystart.Start` and `taskdone` already call `inbox.Read` with `Own` set to the story, so T-1929's rule reaches their answers. Check that it does, and change their code only where it does not.

- In `flai/cmd/task_done.go`, the printed inbox (also used by `flai story start`) prints `awaiting` per thread: print the threads awaiting the agent only, and say nothing of the operator's.
- Add a test to each of `story_start` and `task_done`, at the package and at the MCP tool, with open threads awaiting the operator and none on the story: the answer lists no thread, `awaiting_you` 0, and `threads_omitted` 0. A thread on the story, or a task of it, whose last entry is the operator's is listed, awaiting the agent.
- Update the tests that expect every open thread to await the story's agent, such as `awaiting_you` 10 in `flai/internal/mcpserver/story_start_test.go`.

It waits for T-1929, whose rule it tests. It runs beside T-1930: no path in common.

## Done when

- The new tests pass, with the updated ones, under `flai test` on the packages and `flai/cmd`.
- `flai story start` and `flai task done` print no thread awaiting the operator.

## Notes

`flai/internal/storystart/start.go` and `flai/internal/taskdone/taskdone.go` are left out of the touches: the rule should reach them through `inbox.Read`. `flai task done` widens the touches if one changes.
