---
id: T-1551
type: task
nature: feature
title: flai design show takes --story, and flai design proposals lists the open stories whose branches change the design system
status: in-progress
parent: S-0384
owner: alex
created: 2026-10-13T17:48:12Z
updated: 2026-10-15T03:23:16Z
transitions:
  - to: ready
    at: 2026-10-15T03:17:14Z
    by: system-flow
  - to: in-progress
    at: 2026-10-15T03:23:16Z
    by: system-flow
stream: S-0384
tags: [cli]
touches: [flai/cmd/design.go, flai/cmd/design_test.go]
after: [T-1548]
---
# T-1551 flai design show takes --story, and flai design proposals lists the open stories whose branches change the design system

## Work

In `flai/cmd/design.go`, `flai design show --story S-nnnn` prints `ShowStory`, and `flai design proposals` prints `Proposals`, each in text and with `--json`. Tests cover a story with a change, a closed story, and a story without a worktree. It waits for T-1548, whose functions it prints.

## Done when

- `go test ./cmd/ -run Design` passes with the new cases.

## Notes

- Layer 2, beside T-1552.
