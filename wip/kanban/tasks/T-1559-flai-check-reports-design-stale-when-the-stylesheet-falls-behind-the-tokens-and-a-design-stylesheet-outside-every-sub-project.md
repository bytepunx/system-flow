---
id: T-1559
type: task
nature: feature
title: flai check reports design.stale when the stylesheet falls behind the tokens, and a design.stylesheet outside every sub-project
status: ready
parent: S-0385
owner: alex
created: 2026-10-13T17:48:53Z
updated: 2026-10-15T02:23:02Z
transitions:
  - to: ready
    at: 2026-10-15T02:23:02Z
    by: agent-S-0385
stream: S-0385
tags: [cli]
touches: [flai/internal/check/design.go, flai/internal/check/design_test.go]
after: [T-1556, T-1557]
---
# T-1559 flai check reports design.stale when the stylesheet falls behind the tokens, and a design.stylesheet outside every sub-project

## Work

In `flai/internal/check/design.go`, add `design.stale`: a warning when the stylesheet `design.stylesheet` names is missing or differs from `Stylesheet`, on the stylesheet and on the token file both, so that a branch that changed only the token file stops at its close-out (ADR-0152), each naming `flai design css --write`. No finding when the manifest names no stylesheet, when there is no token file, or when the token file has problems (`design.tokens` reports those). Add an error, on the manifest, when `design.stylesheet` lies outside every sub-project. Tests cover each. It waits for T-1556 and T-1557, whose key and comparison it checks.

## Done when

- `flai test` on `flai/internal/check` passes with the new cases, and `flai check --strict` on this repository is unchanged.

## Notes

- Layer 2, beside T-1558.
