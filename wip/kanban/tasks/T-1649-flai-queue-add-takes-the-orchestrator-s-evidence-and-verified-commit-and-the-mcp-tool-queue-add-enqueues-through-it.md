---
id: T-1649
type: task
nature: feature
title: flai queue add takes the orchestrator's evidence and verified commit, and the MCP tool queue_add enqueues through it
status: backlog
parent: S-0403
owner: alex
created: 2026-10-16T05:01:21Z
updated: 2026-10-16T05:01:21Z
transitions: []
stream: S-0403
tags: [cli]
touches: [flai/cmd/queue.go, flai/cmd/queue_test.go, flai/internal/mcpserver/queue.go, flai/internal/mcpserver/queue_test.go, flai/internal/mcpserver/folder.go]
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
