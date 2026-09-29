---
id: T-0520
type: task
nature: feature
title: Write the ADR draft with the recommendation
status: done
parent: S-0145
owner: alex
created: 2026-09-29T02:12:50Z
updated: 2026-09-29T02:21:21Z
transitions:
  - to: ready
    at: 2026-09-29T02:20:52Z
    by: agent-S-0145
  - to: in-progress
    at: 2026-09-29T02:20:52Z
    by: agent-S-0145
  - to: done
    at: 2026-09-29T02:21:21Z
    by: agent-S-0145
stream: S-0145
tags: []
---

# T-0520 Write the ADR draft with the recommendation

## Work

Create the ADR with `flai adr new --status proposed --body-stdin`: context with the measurements, the recommended approach as one decision, its consequences, and every alternative with why it lost. The recommendation must be one that can cut the pack well below one tool result, and the ADR must say by how much on the measured story.

## Done when

A proposed ADR exists in `design/adrs`, is in the index, `flai check --strict` passes, and the ADR states the estimated pack size on S-0138 under the recommendation.

## Notes
