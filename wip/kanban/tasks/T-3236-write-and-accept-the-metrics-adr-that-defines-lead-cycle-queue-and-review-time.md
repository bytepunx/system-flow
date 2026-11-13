---
id: T-3236
type: task
nature: feature
title: Write and accept the metrics ADR that defines lead, cycle, queue, and review time and time in state as readings of the workflow's timers
status: backlog
parent: S-0750
owner: alex
created: 2026-11-12T16:08:33Z
updated: 2026-11-12T16:08:33Z
transitions: []
stream: S-0750
tags: [cli]
touches: [design/adrs/README.md]
---
# T-3236 Write and accept the metrics ADR that defines lead, cycle, queue, and review time and time in state as readings of the workflow's timers

## Work

Write the ADR with `flai adr new`: each of today's definitions in `design/system/metrics.md` § Derived timestamps and § Per item durations restated as a timer of the baseline, its start and end points or the states it sums, how a reopened or cancelled item counts (ADR-0055, ADR-0190), and that the baseline's figures stay byte for byte. Ask the operator on a thread on the story for what is new, with your recommendation, and accept it once answered.

## Done when

The ADR is accepted and the index lists it.

## Notes
