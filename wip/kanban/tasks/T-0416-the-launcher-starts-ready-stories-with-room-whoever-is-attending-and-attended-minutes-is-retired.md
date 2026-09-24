---
id: T-0416
type: task
nature: remediation
title: The launcher starts ready stories with room whoever is attending, and attended_minutes is retired
status: done
parent: S-0116
owner: alex
created: 2026-09-24T09:00:32Z
updated: 2026-09-24T09:06:36Z
transitions:
  - to: ready
    at: 2026-09-24T09:01:30Z
    by: agent-S-0116
  - to: in-progress
    at: 2026-09-24T09:01:30Z
    by: agent-S-0116
  - to: done
    at: 2026-09-24T09:06:36Z
    by: agent-S-0116
stream: S-0116
tags: []
touches: [flai/internal/serve, flai/cmd, flai/internal/hostapi, flai/internal/config, flaiover/src]
---
# T-0416 The launcher starts ready stories with room whoever is attending, and attended_minutes is retired

## Work

Remove the attendance hold from `look` in `flai/internal/serve/agents.go`: `attended`, `ownSigns`, and the hold memory go. The limit alone holds a story back. Retire `attended_minutes`:
- `--attended-minutes` stays accepted, marked deprecated, and does nothing;
- the config field is no longer read;
- the hostapi settings write and the dashboard's settings panel drop it.

## Done when

- A behavior test shows a story with room is started at the next look while an MCP cursor and a narrative were just written.
- The tests of the hold are removed or rewritten.
- `make test`, lint, and the flaiover tests pass.

## Notes
