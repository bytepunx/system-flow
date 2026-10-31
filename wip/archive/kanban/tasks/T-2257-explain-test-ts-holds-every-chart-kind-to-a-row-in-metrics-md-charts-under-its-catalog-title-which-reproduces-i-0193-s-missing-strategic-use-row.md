---
id: T-2257
type: task
nature: remediation
title: "explain.test.ts holds every chart kind to a row in metrics.md § Charts under its catalog title, which reproduces I-0193's missing Strategic Use row"
status: done
parent: S-0527
owner: alex
created: 2026-10-28T15:01:15Z
updated: 2026-10-30T00:15:20Z
transitions:
  - to: ready
    at: 2026-10-30T00:13:49Z
    by: agent-S-0527
  - to: in-progress
    at: 2026-10-30T00:13:50Z
    by: agent-S-0527
  - to: done
    at: 2026-10-30T00:15:20Z
    by: agent-S-0527
stream: S-0527
tags: [flaiover, test]
touches: [flaiover/src/lib/viz/explain.test.ts]
after: [T-2256]
usage:
  source: log
  seconds: 90
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 13
      output: 2830
      cache_read: 1017548
      cache_write: 31262
      cost: 0.5103
---
# T-2257 explain.test.ts holds every chart kind to a row in metrics.md § Charts under its catalog title, which reproduces I-0193's missing Strategic Use row

## Work

Add to `flaiover/src/lib/viz/explain.test.ts`, which already reads `design/system/metrics.md` from the repository to check each explainer's anchor, a test that every kind in `KINDS` (`flaiover/src/lib/viz/catalog.ts`) has a row in the § Charts table whose first cell is the kind's title from the catalog's title map, trimmed, as `TITLES` or whatever name the catalog gives it. Parse the table from the lines between the `## Charts` heading and the next heading of any level that start with `|`, the header and its rule left out, and take each row's first cell. Name it after I-0193: before T-2256 it fails on Agent Waiting, Strategic Cost, and Strategic Use, which the table lacked, which is the test that reproduces the issue; after T-2256 it passes. Cover both the row's presence and that a kind's title is matched whole, so that `Tokens / Day` does not satisfy `Tokens / Min`. Do not move the fixture reading into a helper file: the test stays one file, as S-0485 left it.

It waits for T-2256 because it fails until the rows exist, and `flai task done` runs the flaiover tier on the files it commits.

Run it with `flai test flaiover/src/lib/viz/explain.test.ts`, and the whole flaiover vitest tier once, so that no other test reads the § Charts rows by position.

## Done when

- `explain.test.ts` has a test that fails when a kind of `KINDS` has no § Charts row under its catalog title, and it passes on the branch with T-2256's rows.
- The test's name cites I-0193.
- `flai test` on the file is green, prettier and eslint included, and the vitest tier as a whole is green.

## Notes

Drafted by planner-S-0527. The chart titles live in `flaiover/src/lib/viz/catalog.ts`; the test reads them and changes nothing there, so that file is not a touch.
