---
id: T-1552
type: task
nature: feature
title: The host method design.get takes a story, and design.proposals lists the open stories whose branches change the design system
status: in-progress
parent: S-0384
owner: alex
created: 2026-10-13T17:48:15Z
updated: 2026-10-15T03:23:17Z
transitions:
  - to: ready
    at: 2026-10-15T03:17:15Z
    by: system-flow
  - to: in-progress
    at: 2026-10-15T03:23:17Z
    by: system-flow
stream: S-0384
tags: [cli]
touches: [flai/internal/hostapi/design.go, flai/internal/hostapi/design_test.go]
after: [T-1548]
---
# T-1552 The host method design.get takes a story, and design.proposals lists the open stories whose branches change the design system

## Work

In `flai/internal/hostapi/design.go`, `design.get` takes an optional `story` and answers `ShowStory` for it, and the new read `design.proposals` answers `Proposals`. Tests cover both with and without a proposal, and a refused story. It waits for T-1548, whose functions it answers.

## Done when

- `go test ./internal/hostapi/` passes, the contract tests included.

## Notes

- Layer 2, beside T-1551.
