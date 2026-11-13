---
id: T-3230
type: task
nature: feature
title: writes_describe.go describes each move from its transition, in the workflow's names and with the reason it asks for
status: backlog
parent: S-0747
owner: alex
created: 2026-11-12T16:08:13Z
updated: 2026-11-12T16:08:13Z
transitions: []
stream: S-0747
tags: [cli]
touches: [flai/internal/hostapi/writes_describe.go, flai/internal/hostapi/writes_test.go]
after: [T-3228]
---
# T-3230 writes_describe.go describes each move from its transition, in the workflow's names and with the reason it asks for

## Work

Waits for T-3228: both change `writes_test.go`, and the descriptions name the moves that task allows. Replace the 31 references to a state's name in `writes_describe.go` with descriptions built from the workflow's transitions: each move's from and to in the workflow's names, who may make it, and the reason it asks for. With the baseline, every description reads as today.

## Done when

`writes_describe.go` names no state, and its descriptions match today's word for word on the baseline.

## Notes
