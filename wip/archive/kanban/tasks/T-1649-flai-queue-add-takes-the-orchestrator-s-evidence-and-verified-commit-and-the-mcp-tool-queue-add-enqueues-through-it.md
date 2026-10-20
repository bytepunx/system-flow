---
id: T-1649
type: task
nature: feature
title: flai queue add takes the orchestrator's evidence and verified commit, and the MCP tool queue_add enqueues through it
status: done
parent: S-0403
owner: alex
created: 2026-10-16T05:01:21Z
updated: 2026-10-20T00:01:44Z
transitions:
  - to: ready
    at: 2026-10-19T23:43:34Z
    by: agent-S-0403
  - to: in-progress
    at: 2026-10-19T23:43:35Z
    by: agent-S-0403
  - to: done
    at: 2026-10-20T00:01:44Z
    by: agent-S-0403
stream: S-0403
tags: [cli]
touches: [flai/cmd/queue.go, flai/cmd/queue_test.go, flai/internal/mcpserver/queue.go, flai/internal/mcpserver/queue_test.go, flai/internal/mcpserver/folder.go, flai/internal/mergequeue/state.go, flai/internal/mergequeue/state_test.go, flai/internal/mcpserver/folder_test.go, flai/internal/mcpserver/server_test.go]
usage:
  source: log
  seconds: 1089
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 66
      output: 23068
      cache_read: 4839553
      cache_write: 175576
      cost: 2.5189
---
# T-1649 flai queue add takes the orchestrator's evidence and verified commit, and the MCP tool queue_add enqueues through it

## Work

- `flai queue add S-nnnn --by orchestrator --verified <commit> --evidence <file|->` records the evidence and the commit with the story in the queue's state. It applies the evidence rules of `flai accept --by orchestrator`: a `Verdict:` line, and a changed file for every criterion. It is refused unless `accept_reviews` is on.
- The MCP tool `queue_add`, registered in `folder.go`, runs it with the evidence as a string.
- Test the refusals, the record, and the tool.
- First in its layer: it waits for no task.

## Done when

The orchestrator can enqueue with evidence through the CLI and MCP, the refusals hold, and the tests pass.

## Notes

- Reviewed by agent-S-0403 at start: the queue's state records `verified` and `evidence` on the orchestrator's `Enqueued` entry, so that T-1650's landing writes the evidence into the story's Notes. `--verified` must be the head of the story's branch, as the acceptance's gate asks; under `FLAI_ROLE=orchestrate` only `--by orchestrator` is taken; `--verified` and `--evidence` are refused for anyone else.
