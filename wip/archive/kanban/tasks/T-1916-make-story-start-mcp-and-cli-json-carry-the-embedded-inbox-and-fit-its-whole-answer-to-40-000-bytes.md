---
id: T-1916
type: task
nature: improvement
title: Make story_start, MCP and CLI --json, carry the embedded inbox and fit its whole answer to 40,000 bytes
status: done
parent: S-0450
owner: alex
created: 2026-10-20T19:05:30Z
updated: 2026-10-20T19:51:49Z
transitions:
  - to: ready
    at: 2026-10-20T19:15:54Z
    by: agent-S-0450
  - to: in-progress
    at: 2026-10-20T19:15:55Z
    by: agent-S-0450
  - to: done
    at: 2026-10-20T19:51:49Z
    by: agent-S-0450
stream: S-0450
tags: [flai]
touches: [flai/internal/storystart/start.go, flai/internal/storystart/start_test.go, flai/internal/mcpserver/story_start.go, flai/internal/mcpserver/story_start_test.go, flai/cmd/story_start.go, flai/cmd/story_start_test.go, design/adrs/0179-task-done-and-story-start-carry-an-inbox-that-lists-the-story-s-threads-and.md, flai/internal/inbox/inbox.go]
after: [T-1915]
usage:
  source: log
  seconds: 2154
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 76
      output: 22705
      cache_read: 4286101
      cache_write: 215145
      cost: 2.6207
---
# T-1916 Make story_start, MCP and CLI --json, carry the embedded inbox and fit its whole answer to 40,000 bytes

## Work

- `flai/internal/storystart/start.go`: read the inbox with the story set, through T-1915's embedded read.
- `flai/internal/mcpserver/story_start.go`: make `fit` use T-1915's inbox fit after the pack is cut to its header, so that threads can no longer keep the answer over `context.PartLimit`. Update `storyStartDescription` to name the count of the other threads and that `inbox` lists them.
- `flai/cmd/story_start.go`: fit `--json` the same way, so that the help's claim that it answers as the MCP tool does holds for size too; print the count of the other threads in the text form.
- Tests: in `flai/internal/mcpserver/story_start_test.go` and `flai/cmd/story_start_test.go`, reproduce S-0399's instance, a project with 32 threads awaiting the agent and none on the story, and check that each answer is at most 40,000 bytes, lists no thread of another item, and counts them; keep a thread on the story listed.

It waits for T-1915, whose read and fit it calls. It runs beside T-1917: the two share no path.

## Done when

- `story_start` and `flai story start --json` answer at most 40,000 bytes of JSON in the reproduction, with the story's own threads listed and the others counted.
- The existing story_start tests pass, changed only where the answer's shape changed.
- `flai test` on the changed paths passes.

## Notes
