---
id: T-3109
type: task
nature: remediation
title: stream_state, wait_for_events, and inbox go on past a work item file that does not parse, with tests
status: done
parent: S-0716
owner: alex
created: 2026-11-11T08:56:52Z
updated: 2026-11-11T15:08:59Z
transitions:
  - to: ready
    at: 2026-11-11T14:55:08Z
    by: agent-S-0716
  - to: in-progress
    at: 2026-11-11T14:55:08Z
    by: agent-S-0716
  - to: done
    at: 2026-11-11T15:08:59Z
    by: agent-S-0716
stream: S-0716
tags: [flai, mcp]
touches: [flai/internal/workitem/narrative_test.go, flai/internal/mcpserver/inbox_test.go, flai/internal/mcpserver/stream_test.go, flai/internal/mcpserver/folder_test.go]
after: [T-3107]
usage:
  source: log
  seconds: 831
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 54
      output: 14218
      cache_read: 2678046
      cache_write: 107660
      cost: 1.4652
---
# T-3109 stream_state, wait_for_events, and inbox go on past a work item file that does not parse, with tests

## Work

In I-0344's instance, `wait_for_events` failed three times, and `stream_state` wrote the narrative but not `wip/agents/index.md`. Both read `Repo.List`: `catchUp` in `mcpserver/cursor.go` through `listItems`, `inbox.go`, and `WriteStreamState` in `workitem/narrative.go` before `WriteIndex`. Once T-3107 lands, they no longer fail. Pin that with tests that put a `PLACEHOLDER` file with no front matter in the tasks folder:

- `wait_for_events` and `inbox` answer with the other items' changes and no error (`mcpserver/inbox_test.go`).
- `stream_state` writes the narrative and `wip/agents/index.md` (`mcpserver/stream_test.go`).
- `WriteIndex` leaves the unreadable item out of the status column and goes on (`workitem/narrative_test.go`). Change `narrative.go` only if `WriteIndex` or `WriteStreamState` still fails.

Read the callers of `Repo.List` in `flai/internal/storystart`, `taskdone`, `itemedit`, and `cmd/board.go` for any that depend on the list failing. Treat a bad file in a thread, issue, or message folder as out of scope, and name any you find in the narrative.

This task waits for T-3107, whose loader change these tests pin. It runs beside T-3108, and the two touch no common file.

## Done when

- The three tests pass with a `PLACEHOLDER` task file present. Each one failed before T-3107.
- `flai test` on the changed files passes.

## Notes
