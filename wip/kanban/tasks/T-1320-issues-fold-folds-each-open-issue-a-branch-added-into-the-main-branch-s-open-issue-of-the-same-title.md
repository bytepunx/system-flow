---
id: T-1320
type: task
nature: remediation
title: issues.Fold folds each open issue a branch added into the main branch's open issue of the same title
status: done
parent: S-0326
owner: alex
created: 2026-10-08T00:33:44Z
updated: 2026-10-08T06:37:46Z
transitions:
  - to: ready
    at: 2026-10-08T06:35:37Z
    by: agent-S-0326
  - to: in-progress
    at: 2026-10-08T06:35:37Z
    by: agent-S-0326
  - to: done
    at: 2026-10-08T06:37:46Z
    by: agent-S-0326
stream: S-0326
tags: [cli, go]
touches: [flai/internal/issues/fold.go, flai/internal/issues/fold_test.go]
after: [T-1332]
usage:
  source: log
  seconds: 129
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 28
      output: 138
      cache_read: 1109279
      cache_write: 65092
      cost: 0.524
---
# T-1320 issues.Fold folds each open issue a branch added into the main branch's open issue of the same title

## Work

Write the fold T-1319's ADR decides as a function of the issues package, in a new `flai/internal/issues/fold.go`, apart from any git: given a story's checkout and the names of the issue files the main branch has, it finds each open issue the branch added whose title an open issue on the main branch has too, moves its instances into that issue as `issues.Bump` would (count, cost, `last_reported`, `updated`), deletes the branch's file, and reports what it folded into what, so that the caller can commit it and regenerate the summary.

It waits for T-1319 because the ADR decides what the fold keeps and what it does to a story made from the folded issue. It runs beside T-1322, which touches documents only.

## Done when

- A unit test in `flai/internal/issues/fold_test.go` folds a branch's issue into the main branch's issue of the same title: the instances, count, and cost add up, and the branch's file is gone
- It leaves alone an issue the branch added under a title the main branch has no open issue for, an issue the main branch has too, and a closed one
- `flai test flai/internal/issues` passes

## Notes
