---
id: T-1390
type: task
nature: improvement
title: metrics.md defines an empty wake and a turn on the neutral events, and flai-cli.md describes capabilities and readers
status: backlog
parent: S-0352
owner: alex
created: 2026-10-08T08:46:07Z
updated: 2026-10-08T08:46:07Z
transitions: []
stream: S-0352
tags: [cli]
touches: [design/system/metrics.md, design/system/flai-cli.md]
after: [T-1389]
---
# T-1390 metrics.md defines an empty wake and a turn on the neutral events, and flai-cli.md describes capabilities and readers

## Work

- `design/system/metrics.md`, the rows for an empty wake and a turn: state each on the neutral events (a `ToolCall` of `wait_for_events` with no `ParentCall` whose `ToolResult` timed out with nothing; a `Call` of the story's own agent that calls a tool, gathered by message ID), then give Claude Code's stream-json fields as that reader's mapping. No figure changes; cite ADR-0130.
- `design/system/flai-cli.md`, Internal structure: `Capabilities` and `Reader` on the adapter, the event kinds, and that `stream.go` still reads stream-json for display.

It waits for the code task of the story's second layer, so that it describes the events as built.

## Done when

- The two documents match the code.
- `flai test design/system/metrics.md design/system/flai-cli.md` passes.

## Notes

Layer 3 of S-0352.
