---
id: T-1677
type: task
nature: feature
title: The workflow, the strategic agents' design, and the CLI, dashboard, and operator docs describe landing and bisection
status: backlog
parent: S-0400
owner: alex
created: 2026-10-16T05:03:35Z
updated: 2026-10-16T05:03:35Z
transitions: []
stream: S-0400
tags: [cli]
touches: [design/system/workflow.md, design/system/strategic-agents.md, design/system/flai-cli.md, design/system/flaiover-dashboard.md, docs/users/flai.md, docs/operators/index.md]
after: [T-1642, T-1669, T-1660]
---
# T-1677 The workflow, the strategic agents' design, and the CLI, dashboard, and operator docs describe landing and bisection

## Work

- `design/system/workflow.md` § Transitions and who makes them: review to done by the queue's landing; review to in-progress by `flai serve` for the queue's send-back.
- `design/system/strategic-agents.md` § Accepting a story: the batch's verification in place of the head's.
- `design/system/flai-cli.md`, `docs/users/flai.md`, and `docs/operators/index.md`: landing, bisection, the send-back's thread, and a batch dropped when main moves.
- `design/system/flaiover-dashboard.md`: the panel's steps and outcomes.
- Waits for T-1642, T-1669, and T-1660, whose decision and behaviour it describes.

## Done when

The documents say the above, and the markdown lint and `flai check --strict` pass.

## Notes
