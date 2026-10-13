---
id: T-1556
type: task
nature: feature
title: The manifest's design.stylesheet names the stylesheet the design system's tokens are written to
status: backlog
parent: S-0385
owner: alex
created: 2026-10-13T17:48:42Z
updated: 2026-10-13T17:48:42Z
transitions: []
stream: S-0385
tags: [cli]
touches: [flai/internal/manifest/manifest.go, flai/internal/manifest/manifest_test.go]
---
# T-1556 The manifest's design.stylesheet names the stylesheet the design system's tokens are written to

## Work

Add `design.stylesheet` to the manifest in `flai/internal/manifest/manifest.go`: a repository path, optional, and a helper that says which sub-project holds it, or none. Tests cover it set, unset, and outside every sub-project. It waits for nothing.

## Done when

- `go test ./internal/manifest/` passes with the three cases.

## Notes

- Layer 1, beside T-1557.
