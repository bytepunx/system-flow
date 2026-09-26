---
id: T-0474
type: task
nature: feature
title: "Set and clear after: with flai edit, MCP item_edit, and the host API"
status: done
parent: S-0130
owner: alex
created: 2026-09-26T17:55:02Z
updated: 2026-09-26T18:02:25Z
transitions:
  - to: ready
    at: 2026-09-26T17:55:12Z
    by: agent-S-0130
  - to: in-progress
    at: 2026-09-26T17:59:29Z
    by: agent-S-0130
  - to: done
    at: 2026-09-26T18:02:25Z
    by: agent-S-0130
stream: S-0130
tags: []
touches: [flai/internal/itemedit, flai/cmd, flai/internal/mcpserver, flai/internal/hostapi]
---
# T-0474 Set and clear after: with flai edit, MCP item_edit, and the host API

## Work

- `itemedit.Change` gains `After` and `ClearAfter`; the view carries `after`; a change names `after` in what changed, so agents are told.
- Validate the entries (story IDs, canonical form, not the story itself); `flai check` refuses a cycle through the edit's own check.
- `flai edit --after S-1,S-2` and `--clear-after`; MCP `item_edit` gains `after` and `clear_after`; the host API's item edit, which the dashboard's editor calls, takes `after`.

## Done when

- Tests show each of the three setting and clearing `after`, and a cycle or unknown story refused.
- `make test` and lint pass.

## Notes
