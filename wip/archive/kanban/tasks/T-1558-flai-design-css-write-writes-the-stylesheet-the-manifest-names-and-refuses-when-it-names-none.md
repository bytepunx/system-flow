---
id: T-1558
type: task
nature: feature
title: flai design css --write writes the stylesheet the manifest names, and refuses when it names none
status: done
parent: S-0385
owner: alex
created: 2026-10-13T17:48:50Z
updated: 2026-10-15T02:48:33Z
transitions:
  - to: ready
    at: 2026-10-15T02:23:02Z
    by: agent-S-0385
  - to: in-progress
    at: 2026-10-15T02:30:47Z
    by: agent-S-0385
  - to: done
    at: 2026-10-15T02:48:33Z
    by: agent-S-0385
stream: S-0385
tags: [cli]
touches: [flai/cmd/design.go, flai/cmd/design_test.go, docs/users/flai-reference.md, docs/operators/settings.md, flai/cmd/reference_test.go]
after: [T-1556, T-1557]
usage:
  source: log
  seconds: 1066
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 49
      output: 293
      cache_read: 1715064
      cache_write: 74858
      cost: 0.8282
---
# T-1558 flai design css --write writes the stylesheet the manifest names, and refuses when it names none

## Work

In `flai/cmd/design.go`, `flai design css --write` writes `Stylesheet` to the manifest's `design.stylesheet`, creating its folder, and prints the path; it refuses, writing nothing, when the manifest names none, when there is no token file, or when the token file has problems. `flai design init` writes the stylesheet too when the manifest names one, since ADR-0150 has flai write it whenever it writes the token file. Tests cover each.

Regenerate `docs/users/flai-reference.md` and the flag index in `docs/operators/settings.md` with `make flai-reference` (`scripts/flai-reference.sh`), whose test fails when the help changes without them. It waits for T-1556 and T-1557, whose key and content it uses, and for T-1556's settings row, which shares `docs/operators/settings.md`.

## Done when

- `flai test` on `flai/cmd` passes with the new cases and the reference test.

## Notes

- Layer 2, beside T-1559.
