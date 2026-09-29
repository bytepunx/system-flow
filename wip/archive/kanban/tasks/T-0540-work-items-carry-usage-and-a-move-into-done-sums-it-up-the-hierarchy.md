---
id: T-0540
type: task
nature: feature
title: Work items carry usage, and a move into done sums it up the hierarchy
status: done
parent: S-0143
owner: alex
created: 2026-09-29T06:08:15Z
updated: 2026-09-29T06:15:48Z
transitions:
  - to: ready
    at: 2026-09-29T06:11:04Z
    by: agent-S-0143
  - to: in-progress
    at: 2026-09-29T06:11:04Z
    by: agent-S-0143
  - to: done
    at: 2026-09-29T06:15:48Z
    by: agent-S-0143
stream: S-0143
tags: []
touches: [flai/internal/workitem, flai/cmd, flai/internal/mcpserver, design/system/work-hierarchy.md, template]
---
# T-0540 Work items carry usage, and a move into done sums it up the hierarchy

## Work

- `usage` in the front matter of epics, stories, and tasks: `source` (`log` when measured from an agent's log, `sum` when summed from children), `seconds` of agent work, `estimated` when any cost in it is, and `models`, one entry per model with `input`, `output`, `cache_read`, `cache_write` tokens and `cost` in US dollars.
- A move into done, by `flai move`, the MCP `item_move`, or `flai accept`, recomputes each ancestor's usage: an epic's is the sum of its stories', archived ones included; a story's, unless an agent's log measured it, the sum of its tasks'.
- `flai check` validates the block; `flai show` prints it.
- Schema in `design/system/work-hierarchy.md` and the template's copy.

## Done when

- Tests show the cascade from a task to its story and epic, a measured story kept as measured, and a refused move writing nothing.
- `make test`, lint, and `flai check --strict` pass.

## Notes
