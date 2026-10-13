---
id: T-1558
type: task
nature: feature
title: flai design css --write writes the stylesheet the manifest names, and refuses when it names none
status: backlog
parent: S-0385
owner: alex
created: 2026-10-13T17:48:50Z
updated: 2026-10-13T17:48:50Z
transitions: []
stream: S-0385
tags: [cli]
touches: [flai/cmd/design.go, flai/cmd/design_test.go]
after: [T-1556, T-1557]
---
# T-1558 flai design css --write writes the stylesheet the manifest names, and refuses when it names none

## Work

In `flai/cmd/design.go`, `flai design css --write` writes `Stylesheet` to the manifest's `design.stylesheet`, creating its folder, and prints the path; it refuses, writing nothing, when the manifest names none or the design system has a token that does not parse. Tests cover each. It waits for T-1556 and T-1557, whose key and content it uses.

## Done when

- `go test ./cmd/ -run Design` passes with the new cases.

## Notes

- Layer 2, beside T-1559.
