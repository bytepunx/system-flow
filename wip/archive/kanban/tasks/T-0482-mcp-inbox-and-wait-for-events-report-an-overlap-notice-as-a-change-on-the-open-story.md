---
id: T-0482
type: task
nature: feature
title: MCP inbox and wait_for_events report an overlap notice as a change on the open story
status: done
parent: S-0132
owner: alex
created: 2026-09-26T18:17:26Z
updated: 2026-09-26T18:21:58Z
transitions:
  - to: ready
    at: 2026-09-26T18:17:41Z
    by: agent-S-0132
  - to: in-progress
    at: 2026-09-26T18:19:59Z
    by: agent-S-0132
  - to: done
    at: 2026-09-26T18:21:58Z
    by: agent-S-0132
stream: S-0132
tags: []
touches: [flai/internal/mcpserver]
---
# T-0482 MCP inbox and wait_for_events report an overlap notice as a change on the open story

## Work

`catchUp` reads the overlap notices beside the edit notices and reports each as a change on the open story (kind `overlapped`, cause the accepted story), with a summary that says which paths changed and to run `flai stream sync` and the tests. `wait_for_events` wakes on the overlap log as it does on the edit log, in one project and in folder mode.

## Done when

- Tests show `inbox` and `wait_for_events` report the notice, and not twice.
- `go test -short ./internal/mcpserver/...` passes.

## Notes
