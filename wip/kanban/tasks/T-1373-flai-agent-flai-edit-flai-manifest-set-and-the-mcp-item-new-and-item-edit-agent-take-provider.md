---
id: T-1373
type: task
nature: feature
title: flai agent, flai edit, flai manifest set, and the MCP item_new and item_edit agent take provider
status: backlog
parent: S-0349
owner: alex
created: 2026-10-08T08:43:11Z
updated: 2026-10-08T08:54:16Z
transitions: []
stream: S-0349
tags: [cli]
touches: [flai/cmd/agent.go, flai/cmd/agent_test.go, flai/cmd/edit.go, flai/cmd/edit_test.go, flai/internal/itemedit/itemedit.go, flai/internal/manifest/settings.go, flai/internal/manifest/settings_test.go, flai/internal/mcpserver/items_write.go, flai/internal/mcpserver/items_write_test.go, docs/operators/settings.md, docs/users/flai-reference.md]
after: [T-1372]
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
