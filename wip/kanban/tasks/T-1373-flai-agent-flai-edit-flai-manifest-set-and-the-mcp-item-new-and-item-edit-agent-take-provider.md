---
id: T-1373
type: task
nature: feature
title: flai agent, flai edit, flai manifest set, and the MCP item_new and item_edit agent take provider
status: done
parent: S-0349
owner: alex
created: 2026-10-08T08:43:11Z
updated: 2026-10-12T16:29:17Z
transitions:
  - to: ready
    at: 2026-10-12T16:15:47Z
    by: agent-S-0349
  - to: in-progress
    at: 2026-10-12T16:15:47Z
    by: agent-S-0349
  - to: done
    at: 2026-10-12T16:29:17Z
    by: agent-S-0349
stream: S-0349
tags: [cli]
touches: [flai/cmd/agent.go, flai/cmd/agent_test.go, flai/cmd/edit.go, flai/cmd/edit_test.go, flai/internal/itemedit/itemedit.go, flai/internal/manifest/settings.go, flai/internal/manifest/settings_test.go, flai/internal/mcpserver/items_write.go, flai/internal/mcpserver/items_write_test.go, docs/operators/settings.md, docs/users/flai-reference.md, flai/cmd/items.go, flai/internal/hostapi/writes.go, flai/internal/hostapi/settings.go, flai/internal/hostapi/writes_test.go, flai/cmd/workitems_test.go, flai/internal/mcpserver/folder.go, flai/cmd/serve_actions.go]
after: [T-1372]
usage:
  source: log
  seconds: 810
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 92
      output: 634
      cache_read: 5102602
      cache_write: 127014
      cost: 2.3552
---
# T-1373 flai agent, flai edit, flai manifest set, and the MCP item_new and item_edit agent take provider

## Work

- `flai agent --provider <name>` sets the project's default; `flai edit <S-nnnn> --provider <name>` sets a story's, through `itemedit`, and `--clear-agent` clears it with the rest.
- `flai manifest set` takes `agent.provider`, `planning.agent.provider`, `orchestration.agent.provider`, and `analysis.agent.provider` (`manifest/settings.go`).
- The MCP `item_new` and `item_edit` agent schema in `items_write.go` gains `provider`.
- The flag rows in `docs/operators/settings.md` for the new flags.
- A provider name is checked here by its pattern only: whether the host has the entry is the start's to refuse, since a story may be started on another host.

It waits for T-1372 because both change `docs/operators/settings.md`, and for T-1371 through it.

## Done when

- Tests in `agent_test.go`, `edit_test.go`, `settings_test.go`, and `items_write_test.go` set, print, and clear a provider.
- `flai test flai/cmd/ flai/internal/manifest/ flai/internal/mcpserver/ flai/internal/itemedit/` passes.

## Notes

Layer 3 of S-0349.
