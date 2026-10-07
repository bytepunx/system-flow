---
id: T-1216
type: task
nature: feature
title: An ADR and metrics.md define the coordination values flai stats reports per week
status: backlog
parent: S-0337
owner: alex
created: 2026-10-07T20:17:17Z
updated: 2026-10-07T20:17:17Z
transitions: []
stream: S-0337
tags: [flai]
touches: [design/adrs, design/system/metrics.md]
---
# T-1216 An ADR and metrics.md define the coordination values flai stats reports per week

## Work

`metrics.md` is the contract between `flai stats` and the dashboard and changes only with an ADR. It waits for nothing.

- Add § Coordination to `design/system/metrics.md`, defining per ISO week: conversations opened, closed by agreement (a share or a narrowed claim), escalated, and closed with their story; conflicts found by a trial merge while both stories were open; conflicts met by a rebase at `flai stream sync` or `flai accept`, with the paths; and stories started on a share, with the hold each would have waited, measured to when the holding story moved to review.
- Say where each comes from: the message files, a conflict log the next task writes, and the board's transitions.
- Record the decision in an ADR refining ADR-0081 and ADR-0113.

## Done when

- The ADR is accepted, and `metrics.md` defines every value with its source and its precision.
- `flai check --strict` passes.

## Notes
