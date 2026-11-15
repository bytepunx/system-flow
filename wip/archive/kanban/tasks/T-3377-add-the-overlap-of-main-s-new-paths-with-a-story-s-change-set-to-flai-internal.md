---
id: T-3377
type: task
nature: remediation
title: Add the overlap of main's new paths with a story's change set to flai/internal/verify
status: done
parent: S-0792
owner: alex
created: 2026-11-15T07:43:22Z
updated: 2026-11-15T16:56:01Z
transitions:
  - to: ready
    at: 2026-11-15T16:45:41Z
    by: agent-S-0792
  - to: in-progress
    at: 2026-11-15T16:45:41Z
    by: agent-S-0792
  - to: done
    at: 2026-11-15T16:56:01Z
    by: agent-S-0792
stream: S-0792
tags: [flai, verify]
touches: [flai/internal/verify/overlap.go, flai/internal/verify/overlap_test.go]
after: [T-3376]
usage:
  source: log
  seconds: 620
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 18
      output: 5741
      cache_read: 1386826
      cache_write: 62635
      cost: 0.8405
  strategic:
    - kind: orchestrator
      seconds: 5
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 0
          output: 6
          cache_read: 21939
          cache_write: 603
          cost: 0.0068
---
# T-3377 Add the overlap of main's new paths with a story's change set to flai/internal/verify

## Work

In `flai/internal/verify/overlap.go`, add a function that takes main's new paths (as `MainPathsSince` lists them), the story's paths (`Report.Paths`), and the checkout, and returns the main paths that overlap the story's change set by T-3376's rule:

- a main path the story's paths hold; or
- a main path in a Go package that a package of the story's paths imports, directly or through others. Reuse the `importers` closure in `importers.go`: a main path's package overlaps when a story path's package is among its importers.

A non-Go path overlaps only by being one of the story's paths. Keep the order `MainPathsSince` gives. Return an error when the module's imports cannot be read, so the gate can stay closed on it as ADR-0291 does.

Test it in `overlap_test.go` with a module fixture: the same file; the I-0229 case (main changes `internal/roles`, the story changes a test in `internal/guard`, which imports it) overlaps; a sibling package with no import between them does not; a markdown file the story does not change does not; one it changes does.

It waits for T-3376, whose ADR fixes the rule it implements.

## Done when

- The function returns exactly the overlapping paths for each case above, and its tests pass with `flai test flai/internal/verify/overlap.go flai/internal/verify/overlap_test.go`.

## Notes
