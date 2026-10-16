---
id: T-1651
type: task
nature: feature
title: An ADR adds the queue wait per story, the batch size, and the bisections to metrics.md
status: backlog
parent: S-0404
owner: alex
created: 2026-10-16T05:01:28Z
updated: 2026-10-16T05:01:28Z
transitions: []
stream: S-0404
tags: [cli]
touches: [design/adrs, design/system/metrics.md]
---
# T-1651 An ADR adds the queue wait per story, the batch size, and the bisections to metrics.md

## Work

- Write the ADR with `flai adr new`, as ADR-0138 added coordination.
- Define in `metrics.md`, in a section of its own:
  - the queue wait per story: from enqueueing to landing or send-back
  - the batch size
  - the bisections per batch
  - their weekly aggregates under `queue`
  - the log they are read from.
- First in its layer: it waits for no task.

## Done when

The ADR is accepted and indexed, `metrics.md` defines each value, and `flai check --strict` passes.

## Notes
