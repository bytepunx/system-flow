---
id: T-1914
type: task
nature: improvement
title: Record in an ADR and the design that task_done and story_start carry the story's threads in full and only a count of the rest, fitted to 40,000 bytes
status: done
parent: S-0450
owner: alex
created: 2026-10-20T19:05:11Z
updated: 2026-10-20T19:11:53Z
transitions:
  - to: ready
    at: 2026-10-20T19:08:49Z
    by: agent-S-0450
  - to: in-progress
    at: 2026-10-20T19:08:50Z
    by: agent-S-0450
  - to: done
    at: 2026-10-20T19:11:52Z
    by: agent-S-0450
stream: S-0450
tags: [flai]
touches: [design/adrs/README.md, design/system/agent-narrative.md, design/system/flai-cli.md, design/adrs/0179-task-done-and-story-start-carry-an-inbox-that-lists-the-story-s-threads-and.md]
usage:
  source: log
  seconds: 182
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 26
      output: 7872
      cache_read: 1485951
      cache_write: 74589
      cost: 0.9086
---
# T-1914 Record in an ADR and the design that task_done and story_start carry the story's threads in full and only a count of the rest, fitted to 40,000 bytes

## Work

Propose the fix from I-0167's four instances before building it, as the story's goal asks, and record it.

The instances: `story_start` answered 87,106 to 101,139 characters and `task_done` 83,000 to 84,000, with 31 to 33 threads awaiting the operator, none on the story. Each thread in the inbox carries its last entry and any pending recommendation whole (`flai/internal/inbox/inbox.go`, `ThreadSummary`). `story_start`'s `fit` (`flai/internal/mcpserver/story_start.go`) trims only the changes and the pack, so it cannot bring the threads under the limit; `task_done` and both CLI `--json` forms have no fit at all.

The proposal to record, unless the work shows a better one:

- The inbox that `task_done` and `story_start` carry lists in full only the threads on the story and its tasks, with the inbox's `story` filter, and gives the number of the other threads awaiting the agent in a new field, so that the agent calls `inbox` for them.
- The whole answer, MCP and CLI `--json` alike, is fitted to `context.PartLimit` (40,000 bytes, ADR-0104): changes are trimmed first, then the story's own threads' last entries, and what was left out is counted.

Write the ADR with `flai adr new` and accept it; it refines ADR-0104 and ADR-0107. Update `design/system/agent-narrative.md` § What an agent is told through MCP and the `flai story start`, `flai task done`, and MCP tool rows of `design/system/flai-cli.md`. The ADR's own file is named by `flai adr new`, so it is not in this task's touches; `flai task done` widens them with it.

This task waits for nothing: it is the first layer.

## Done when

- An accepted ADR states what the embedded inbox carries, the new count field, and the 40,000-byte fit, with I-0167's instances as its context.
- `design/system/agent-narrative.md` and `design/system/flai-cli.md` describe the same and link the ADR.
- `flai check --strict` and the markdown lint pass.

## Notes
