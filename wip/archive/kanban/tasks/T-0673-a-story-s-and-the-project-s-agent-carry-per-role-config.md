---
id: T-0673
type: task
nature: improvement
title: A story's and the project's agent carry per-role config
status: done
parent: S-0189
owner: arobson
created: 2026-10-01T10:53:31Z
updated: 2026-10-01T11:03:25Z
transitions:
  - to: ready
    at: 2026-10-01T10:53:48Z
    by: agent-S-0189
  - to: in-progress
    at: 2026-10-01T11:03:25Z
    by: agent-S-0189
  - to: done
    at: 2026-10-01T11:03:25Z
    by: agent-S-0189
stream: S-0189
tags: []
touches: [flai/internal/manifest, flai/internal/workitem, flai/cmd, flai/internal/mcpserver, flai/internal/hostapi, flai/internal/itemedit, flai/internal/serve]
usage:
  source: log
  seconds: 0
  models: []
---
# T-0673 A story's and the project's agent carry per-role config

## Work

Add `roles` to `manifest.Agent`: a map from a role name (`explore`, `verify`, open to more) to the same shape, harness, model, and config. Parse, validate, render (`AgentBlock`), compare (`Same`), merge (`With`, role by role), and print it. Carry it wherever an agent is set: `flai story new` and `flai agent`, the MCP `item_new`, and the dashboard's edits through hostapi. Document it in `work-hierarchy.md`, `flai-cli.md`, and the user and operator docs.

## Done when

- A story or system-flow.yaml with `agent.roles` round-trips, and one without reads and writes as before.
- A new story copies the default's roles, and flags set them.
- Tests cover parse, validate, merge, compare, and render.

## Notes
