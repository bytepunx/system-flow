---
id: T-0395
type: task
nature: improvement
title: The host page's process table gives State, Version, and Restarts the width and a narrow Process column
status: done
parent: S-0113
owner: alex
created: 2026-09-24T07:42:30Z
updated: 2026-09-24T07:48:59Z
transitions:
  - to: ready
    at: 2026-09-24T07:42:33Z
    by: agent-S-0113
  - to: in-progress
    at: 2026-09-24T07:42:33Z
    by: agent-S-0113
  - to: done
    at: 2026-09-24T07:48:59Z
    by: agent-S-0113
stream: S-0113
tags: []
touches: [flaiover/src/lib/components]
---
# T-0395 The host page's process table gives State, Version, and Restarts the width and a narrow Process column

## Work

- In `flaiover/src/lib/components/HostProcesses.svelte`, stop the Process column taking the table's spare width: keep it as narrow as its label, and wrap the MCP project list under it instead of letting it widen the column.
- Give State, Version, and Restarts the rest of the width, with generous horizontal padding between the columns so their values never sit close together.
- Keep the actions column right-aligned and on one line.
- Pin the layout in `HostProcesses.svelte.test.ts` and look at the page in a browser with and without the actions column.

## Done when

- [x] Process is narrow; State, Version, and Restarts share the remaining width with visible space between them, with the actions column shown and hidden.
- [x] `scripts/flaiover-test.sh` passes.

## Notes
