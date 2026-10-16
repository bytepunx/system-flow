---
id: T-1679
type: task
nature: feature
title: The workflow, the CLI and dashboard designs, and the operator and user docs describe landing on a dirty checkout and catching it up
status: backlog
parent: S-0402
owner: alex
created: 2026-10-16T05:03:42Z
updated: 2026-10-16T05:03:42Z
transitions: []
stream: S-0402
tags: [cli]
touches: [design/system/workflow.md, design/system/flai-cli.md, design/system/flaiover-dashboard.md, docs/operators/index.md, docs/users/flai.md]
after: [T-1671]
---
# T-1679 The workflow, the CLI and dashboard designs, and the operator and user docs describe landing on a dirty checkout and catching it up

## Work

- `design/system/workflow.md` (review to done): the queue lands on a dirty checkout, and `flai accept` by hand keeps its refusal.
- `docs/operators/index.md`: what the operator sees when their checkout is behind, and the command that catches it up.
- `design/system/flai-cli.md`, `design/system/flaiover-dashboard.md`, `docs/users/flai.md`: the same in `flai queue` and the panel.
- Waits for T-1671, whose behaviour it describes.

## Done when

The documents say the above, and the markdown lint and `flai check --strict` pass.

## Notes
