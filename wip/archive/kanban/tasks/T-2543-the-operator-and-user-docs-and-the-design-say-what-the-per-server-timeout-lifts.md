---
id: T-2543
type: task
nature: remediation
title: The operator and user docs and the design say what the per-server timeout lifts, and that an agent image must be rebuilt to carry it
status: done
parent: S-0583
owner: alex
created: 2026-11-07T07:06:38Z
updated: 2026-11-07T07:31:59Z
transitions:
  - to: ready
    at: 2026-11-07T07:26:44Z
    by: system-flow
  - to: in-progress
    at: 2026-11-07T07:26:45Z
    by: system-flow
  - to: done
    at: 2026-11-07T07:31:59Z
    by: agent-S-0583
stream: S-0583
tags: [docs]
touches: [docs/operators/index.md, docs/operators/runbooks/agent-container.md, docs/users/flai.md, design/system/agent-adapters.md, design/system/flai-cli.md]
after: [T-2541]
usage:
  source: log
  seconds: 314
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 33
      output: 8323
      cache_read: 4296300
      cache_write: 63791
      cost: 1.5121
  strategic:
    - kind: orchestrator
      seconds: 2
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 0
          output: 6
          cache_read: 18267
          cache_write: 781
          cost: 0.0047
---
# T-2543 The operator and user docs and the design say what the per-server timeout lifts, and that an agent image must be rebuilt to carry it

## Work

This task waits for T-2541, whose ADR and figure these documents link and state. It shares no path with the code task, so the two run together.

Change what each document says about the 60-second limit:

- `docs/operators/index.md`, the `Long requests` row: a call that does work now runs up to the figure when the MCP file names the `timeout`, as flai's own does, and a hand-written file without it stays at 60 s. Keep the 45-second holds and the proxy warning.
- `docs/users/flai.md`, the paragraph on how long a call is held over HTTP: the same, briefly, linking the ADR.
- `docs/operators/runbooks/agent-container.md`, the `mcp.json` item: the file names the `timeout`. Say that pilot in the image writes this file, so an image built before this change keeps the 60-second limit until `make agent-images` rebuilds it.
- `design/system/agent-adapters.md`: the container harness's MCP and the hold figures name the `timeout` and the ADR.
- `design/system/flai-cli.md`: the `flai mcp http` row's long-call text, and the printed sample if it quotes it.

Change only those passages, and bump each document's `updated`.

## Done when

- No document says that a `task_done` or `story_start` slower than 60 seconds is still cut short for an agent in a container, except for a file without the `timeout` or an image not yet rebuilt.
- Each changed document links the ADR.
- The markdown lint passes on the changed files with `flai test`.

## Notes
