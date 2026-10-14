---
id: T-1494
type: task
nature: improvement
title: flai verify reuses the tiers a clean rebase onto a newer main cannot reach, and runs again those main's new paths select
status: done
parent: S-0374
owner: alex
created: 2026-10-12T17:34:01Z
updated: 2026-10-14T22:03:08Z
transitions:
  - to: ready
    at: 2026-10-14T21:45:25Z
    by: agent-S-0374
  - to: in-progress
    at: 2026-10-14T21:45:25Z
    by: agent-S-0374
  - to: done
    at: 2026-10-14T22:03:08Z
    by: agent-S-0374
stream: S-0374
tags: [cli]
touches: [flai/internal/verify/story.go, flai/internal/verify/paths.go, flai/internal/verify/story_test.go, flai/internal/verify/paths_test.go, design/adrs/0156-a-close-out-that-main-outran-syncs-and-verifies-again-at-most-three-rounds-and.md, docs/users/flai-reference.md, flai/cmd/verify.go]
after: [T-1493]
usage:
  source: log
  seconds: 1063
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 79
      output: 33717
      cache_read: 5377977
      cache_write: 122632
      cost: 2.553
---
# T-1494 flai verify reuses the tiers a clean rebase onto a newer main cannot reach, and runs again those main's new paths select

## Work

Waits for T-1493, whose ADR decides the rule this builds.

- In `flai/internal/verify/story.go`, record on the `Report` what identifies the branch's own change: the `git patch-id --stable` of its diff from the merge base with the base. Widen `reusable` so that a report whose `Commit` or `BaseCommit` differ from this run's still lends its passed tiers when that patch id and the `Fingerprint` are the same, but for each tier that the paths main changed from the report's `BaseCommit` to this run's select. Those tiers, and every tier after the first that runs, run again as resume does now. A tier selected now that the report did not run, or a report without the patch id, reuses nothing.
- In `flai/internal/verify/paths.go`, add what lists the paths changed between two commits of the main branch (`git diff --no-renames --name-only <old>..<new>`), and what computes the branch's patch id.
- Select the tiers for main's paths through the same selection the story's own paths go through (`SelectStory` or `Select`, `{importers}` included), so that a Go package main changed reruns the Go tiers of a branch that imports it.
- The step's text says which tiers it reused across the rebase and which main's paths sent back, as the reused line does now.

Tests, in `flai/internal/verify/story_test.go` and `paths_test.go`, reproduce I-0136 in a real repository: a report passes every tier, main gains a commit outside wip that changes none of the branch's paths, the branch is rebased onto it, and `Verify` passes again reusing the tiers main's paths do not select and running the ones they do. Cover too: a branch whose own change differs after the rebase reuses nothing; `--fresh` reuses nothing; a wip-only main commit keeps ADR-0135's pass-over.

## Done when

- `flai verify` after a clean rebase onto a main that gained an acceptance reuses the tiers that acceptance's paths do not select, and runs the ones they do.
- The tests above pass with `flai test` on the changed paths.

## Notes

Layer 2: runs beside T-1495, which changes only the close-out scripts.
