---
id: T-0522
type: task
nature: feature
title: A size budget on the pack, with sizes in the header
status: done
parent: S-0146
owner: alex
created: 2026-09-29T03:19:37Z
updated: 2026-09-29T03:33:54Z
transitions:
  - to: ready
    at: 2026-09-29T03:19:45Z
    by: agent-S-0146
  - to: in-progress
    at: 2026-09-29T03:19:46Z
    by: agent-S-0146
  - to: done
    at: 2026-09-29T03:33:54Z
    by: agent-S-0146
stream: S-0146
tags: []
touches: [flai/internal/context, flai/cmd/prime.go, flai/internal/manifest]
---
# T-0522 A size budget on the pack, with sizes in the header

## Work

A `--budget` flag on `flai prime --story` (default 80 KB) and a project default under `prime.budget` in `system-flow.yaml`. Every item gets a size; the header gives the budget, the pack size, and one line per item with its size, and says when the conventions or the named documents alone exceed the budget. When the conventions alone exceed it, the pack is the conventions and a catalog.

## Done when

- `--budget` and `prime.budget` parse sizes such as `80KB` and `81920`, with a behaviour test.
- The header lines are tested, including both exceeded cases.

## Notes
