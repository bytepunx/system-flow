---
id: T-1678
type: task
nature: feature
title: The workflow, the coordination design, and the CLI docs describe what the queue settles in the shared paths and what it reports outside them
status: done
parent: S-0401
owner: alex
created: 2026-10-16T05:03:38Z
updated: 2026-10-19T22:24:03Z
transitions:
  - to: ready
    at: 2026-10-19T22:20:05Z
    by: agent-S-0401
  - to: in-progress
    at: 2026-10-19T22:20:06Z
    by: agent-S-0401
  - to: done
    at: 2026-10-19T22:24:02Z
    by: agent-S-0401
stream: S-0401
tags: [cli]
touches: [design/system/workflow.md, design/system/agent-coordination.md, design/system/flai-cli.md, docs/users/flai.md, flai/cmd/stream_sync.go, flai/cmd/stream_sync_test.go]
after: [T-1670]
usage:
  source: log
  seconds: 236
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 76
      output: 39773
      cache_read: 4973322
      cache_write: 152557
      cost: 2.7075
---
# T-1678 The workflow, the coordination design, and the CLI docs describe what the queue settles in the shared paths and what it reports outside them

## Work

- `design/system/workflow.md` § Branches and collisions: the queue's merge steps, the markdown merge, a story kept waiting, and the message or send-back outside the shared paths.
- `design/system/agent-coordination.md` § Messages between stories: the queue as a sender of conflict messages.
- `design/system/flai-cli.md` and `docs/users/flai.md`: the same in `flai queue`'s output.
- Waits for T-1670, whose behaviour it describes.

## Done when

The documents say the above, and the markdown lint and `flai check --strict` pass.

## Notes
