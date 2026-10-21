---
id: T-1930
type: task
nature: remediation
title: Make the MCP inbox, wait_for_work, the folder inbox, and the thread tools apply the rule, with the orchestrator listing the threads awaiting the operator by default
status: done
parent: S-0453
owner: alex
created: 2026-10-20T22:20:41Z
updated: 2026-10-21T01:19:59Z
transitions:
  - to: ready
    at: 2026-10-21T00:55:02Z
    by: system-flow
  - to: in-progress
    at: 2026-10-21T00:55:03Z
    by: system-flow
  - to: done
    at: 2026-10-21T01:19:59Z
    by: system-flow
stream: S-0453
tags: [flai]
touches: [flai/internal/mcpserver/server.go, flai/internal/mcpserver/work.go, flai/internal/mcpserver/folder.go, flai/internal/mcpserver/messages.go, flai/internal/mcpserver/server_test.go, flai/internal/mcpserver/work_test.go, flai/internal/mcpserver/folder_test.go, flai/internal/mcpserver/orchestrator_threads_test.go]
after: [T-1929]
usage:
  source: log
  seconds: 1496
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 51
      output: 18112
      cache_read: 2822299
      cache_write: 127629
      cost: 1.6465
---
# T-1930 Make the MCP inbox, wait_for_work, the folder inbox, and the thread tools apply the rule, with the orchestrator listing the threads awaiting the operator by default

## Work

In `flai/internal/mcpserver`, pass the agent's own story to the inbox and to every `summary` call, from `recordingStory` (`FLAI_STORY`, else `agent-S-nnnn`, else the branch), so the server has one notion of its own story:

- `inbox` in `server.go`: give `inbox.Read` the own story, and, when `s.role` is `orchestrate`, the option that lists the threads awaiting the operator without `all: true`.
- `wait_for_work` in `work.go`: answer `thread` only for a thread `awaiting: "you"` by the rule, so an agent with no story in progress is not woken for a thread awaiting the operator.
- The folder inbox in `folder.go`, and `thread_get`, `thread_open`, `thread_reply`, and the callers in `messages.go`, through `summary`.
- The descriptions: `inboxDescription` in `folder.go`, the server instructions and `wait_for_work`'s in `server.go`, and the `WorkOut` schema in `work.go` name `"operator"` beside `"you"` and `"other"`.

Tests: add one with open threads awaiting the operator, none on the agent's story, that finds no `thread` answer from `wait_for_work` and `awaiting_you` 0 from `inbox`, while a thread the operator answered on the agent's story still wakes it. Make the orchestrator's helper in `orchestrator_threads_test.go` call `inbox` without `all: true` and still find the threads awaiting the operator, as `awaiting: "operator"`. Update the tests in `server_test.go`, `work_test.go`, and `folder_test.go` that assume every thread awaits the reader.

It waits for T-1929, whose `Summarize` and `Read` it calls.

## Done when

- `inbox`, `wait_for_work`, the folder inbox, and the thread tools give `awaiting` by the rule.
- The orchestrator finds the threads awaiting the operator in `inbox` without `all: true`.
- `flai test flai/internal/mcpserver` passes.
- The operator's view is unchanged (criterion 4): `flai/internal/hostapi/people.go` and `flai/cmd/thread.go` are not changed, and `flai test flai/internal/hostapi flai/cmd` passes.

## Notes

`story_start` and `task_done` are T-1932's; the orchestrator's prompt is T-1933's.
