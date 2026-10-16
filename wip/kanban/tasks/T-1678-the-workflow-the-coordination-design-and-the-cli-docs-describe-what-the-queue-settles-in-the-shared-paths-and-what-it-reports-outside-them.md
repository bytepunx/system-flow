---
id: T-1678
type: task
nature: feature
title: The workflow, the coordination design, and the CLI docs describe what the queue settles in the shared paths and what it reports outside them
status: backlog
parent: S-0401
owner: alex
created: 2026-10-16T05:03:38Z
updated: 2026-10-16T05:03:38Z
transitions: []
stream: S-0401
tags: [cli]
touches: [design/system/workflow.md, design/system/agent-coordination.md, design/system/flai-cli.md, docs/users/flai.md]
after: [T-1670]
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
