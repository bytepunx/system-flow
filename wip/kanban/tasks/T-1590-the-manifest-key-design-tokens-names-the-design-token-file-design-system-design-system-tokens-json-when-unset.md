---
id: T-1590
type: task
nature: feature
title: The manifest key design.tokens names the design token file, design/system/design-system.tokens.json when unset
status: done
parent: S-0381
owner: alex
created: 2026-10-14T23:52:52Z
updated: 2026-10-15T00:09:34Z
transitions:
  - to: ready
    at: 2026-10-14T23:57:51Z
    by: agent-S-0381
  - to: in-progress
    at: 2026-10-14T23:57:52Z
    by: agent-S-0381
  - to: done
    at: 2026-10-15T00:09:34Z
    by: agent-S-0381
stream: S-0381
tags: [cli]
touches: [flai/internal/manifest/manifest.go, flai/internal/manifest/manifest_test.go, design/system/project-manifest.md, docs/operators/settings.md]
usage:
  source: log
  seconds: 702
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 48
      output: 245
      cache_read: 1852740
      cache_write: 223631
      cost: 0.9588
---
# T-1590 The manifest key design.tokens names the design token file, design/system/design-system.tokens.json when unset

## Work

Add a `design` block to the manifest in `flai/internal/manifest/manifest.go` with the key `tokens`, the repository path of the project's design token file (ADR-0146), and a method that answers it or `design/system/design-system.tokens.json` when unset. Refuse an absolute path or one that leaves the repository, as other path keys do. S-0385 adds `design.stylesheet` to the same block later, so name the block for the design system, not for tokens alone. Describe the key in `design/system/project-manifest.md` and give it a row in `docs/operators/settings.md`, whose test fails on a manifest key with no row. It waits for nothing: the default is a constant the other tasks call.

## Done when

- `go test ./internal/manifest/` passes with the key unset, set, and refused.
- The settings index test passes with the new row.
- `flai check --strict` on this repository is unchanged.

## Notes

- Layer 1, beside T-1532 and the documentation task, with which it shares no path. Written by planner-S-0381.
