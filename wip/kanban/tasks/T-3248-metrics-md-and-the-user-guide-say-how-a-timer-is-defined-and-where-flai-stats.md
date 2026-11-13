---
id: T-3248
type: task
nature: feature
title: metrics.md and the user guide say how a timer is defined and where flai stats reports it
status: backlog
parent: S-0752
owner: alex
created: 2026-11-12T16:09:27Z
updated: 2026-11-12T16:09:27Z
transitions: []
stream: S-0752
tags: [cli]
touches: [design/system/metrics.md, docs/users/flai.md, docs/users/flai-reference.md]
after: [T-3243]
---
# T-3248 metrics.md and the user guide say how a timer is defined and where flai stats reports it

## Work

Waits for T-3243: it documents the section that task names. Add to `design/system/metrics.md` how a workflow's own timer is computed and where `flai stats --json` gives it, and to `docs/users/flai.md` how an operator adds one and reads it. Regenerate `docs/users/flai-reference.md` with `make flai-reference` if the command's help changed.

## Done when

Both documents say how a timer is defined and read, and `TestReferenceIsCurrent` passes.

## Notes
