---
id: T-1332
type: task
nature: remediation
title: issues.Merge merges an issue file both sides of a stopped rebase changed, keeping the instances of both
status: done
parent: S-0326
owner: alex
created: 2026-10-08T06:25:51Z
updated: 2026-10-08T06:35:27Z
transitions:
  - to: ready
    at: 2026-10-08T06:27:47Z
    by: agent-S-0326
  - to: in-progress
    at: 2026-10-08T06:27:47Z
    by: agent-S-0326
  - to: done
    at: 2026-10-08T06:35:27Z
    by: agent-S-0326
stream: S-0326
tags: [cli, go]
touches: [flai/internal/issues/merge.go, flai/internal/issues/merge_test.go]
after: [T-1319]
usage:
  source: log
  seconds: 460
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 25
      output: 10434
      cache_read: 1584589
      cache_write: 65080
      cost: 0.9077
---
# T-1332 issues.Merge merges an issue file both sides of a stopped rebase changed, keeping the instances of both

## Work

I-0112's second instance (S-0318): two branches bumped the same issue, so the rebase stopped on the issue file and the summary, and the agent merged the instances by hand. Write the merge T-1319's ADR decides in a new `flai/internal/issues/merge.go`, apart from any git: given the issue file's three versions at a stop (the base, the main branch's side, and the replayed commit's side), it writes one issue that holds the base's instances and the instances each side added, with count, average cost, `last_reported`, and `updated` added up from both sides; every other front matter field and the rest of the body is taken from the side that changed it. When both sides changed one of those differently, it merges nothing and says so, and the stop is left for the agent.

It waits for T-1319 because the ADR decides what the merge keeps. It runs beside T-1320 and T-1322, with no path in common.

## Done when

- A unit test in `flai/internal/issues/merge_test.go` merges two bumps of one issue: both instances are kept in timestamp order, the count and the average cost add up, and `last_reported` is the later
- A bump on one side and a close on the other merge into a closed issue holding the bump's instance
- Both sides changing the description or the same front matter field differently merges nothing
- `flai test flai/internal/issues` passes

## Notes
