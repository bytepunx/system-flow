---
id: T-1672
type: task
nature: feature
title: The strategic agents' design, the CLI docs, and the template's changelog describe the orchestrator's enqueueing
status: backlog
parent: S-0403
owner: alex
created: 2026-10-16T05:03:07Z
updated: 2026-10-16T05:03:07Z
transitions: []
stream: S-0403
tags: [cli]
touches: [design/system/strategic-agents.md, design/system/flai-cli.md, docs/users/flai.md, docs/users/flai-reference.md, template/template.yaml, template/CHANGELOG.md]
after: [T-1650, T-1663]
---
# T-1672 The strategic agents' design, the CLI docs, and the template's changelog describe the orchestrator's enqueueing

## Work

- `design/system/strategic-agents.md` § Accepting a story: with the queue on, the orchestrator enqueues, and the queue lands as the orchestrator with the evidence.
- `design/system/flai-cli.md` and `docs/users/flai.md`: `flai queue add`'s orchestrator flags and `queue_add`. Regenerate `flai-reference.md` with `make flai-reference`.
- Bump the template's version in `template.yaml`, and name the convention's change and the flai it needs in `template/CHANGELOG.md`.
- Waits for T-1650 and T-1663, whose behaviour it describes.

## Done when

The documents say the above, the reference is regenerated, and the markdown lint and `flai check --strict` pass.

## Notes
