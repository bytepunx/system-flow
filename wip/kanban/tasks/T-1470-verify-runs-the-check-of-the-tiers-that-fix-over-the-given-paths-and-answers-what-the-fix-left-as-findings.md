---
id: T-1470
type: task
nature: improvement
title: verify runs the check of the tiers that fix over the given paths and answers what the fix left as findings
status: done
parent: S-0367
owner: alex
created: 2026-10-12T17:07:45Z
updated: 2026-10-12T17:09:53Z
transitions:
  - to: ready
    at: 2026-10-12T17:09:18Z
    by: agent-S-0367
  - to: in-progress
    at: 2026-10-12T17:09:19Z
    by: agent-S-0367
  - to: done
    at: 2026-10-12T17:09:53Z
    by: agent-S-0367
stream: S-0367
tags: []
touches: [flai/internal/verify/fix.go, flai/internal/verify/fix_test.go]
after: [T-1458]
usage:
  source: log
  seconds: 34
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 4
      output: 24
      cache_read: 463917
      cache_write: 4012
      cost: 0.2113
---
# T-1470 verify runs the check of the tiers that fix over the given paths and answers what the fix left as findings

## Work

- `verify.CheckoutUnfixed(ctx, root, paths)`: select, among the checkout's tiers that carry a `fix_command`, those the paths select, run their `command` over those paths as `flai test` does, and answer the `Result`, with the failing tier's findings (file, line, message).
- The operator decided on TH-0413 that markdown the fix cannot fix blocks an acceptance with its files, lines, and errors; this is the check `flai accept` runs.
- Waits for T-1458: both build `verify`, which `taskdone` imports.

## Done when

- [ ] A tier with a fix command whose command fails over the paths answers its findings; a tier without one is not run, and no selected path passes, under tests.

## Notes
