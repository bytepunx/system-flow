---
id: T-1559
type: task
nature: feature
title: flai check reports design.stale when the stylesheet falls behind the tokens, and a design.stylesheet outside every sub-project
status: backlog
parent: S-0385
owner: alex
created: 2026-10-13T17:48:53Z
updated: 2026-10-13T17:48:53Z
transitions: []
stream: S-0385
tags: [cli]
touches: [flai/internal/check/design.go, flai/internal/check/design_test.go]
after: [T-1556, T-1557]
---
# T-1559 flai check reports design.stale when the stylesheet falls behind the tokens, and a design.stylesheet outside every sub-project

## Work

In `flai/internal/check/design.go`, add `design.stale`, a warning on the stylesheet when `Stale` says it differs, naming `flai design css --write`, and an error when `design.stylesheet` lies outside every sub-project. A manifest naming no stylesheet has no finding. Tests cover each. It waits for T-1556 and T-1557, whose key and comparison it checks.

## Done when

- `go test ./internal/check/` passes with the new cases, and `flai check --strict` on this repository is unchanged.

## Notes

- Layer 2, beside T-1558.
