---
id: T-1019
type: task
nature: remediation
title: Close I-0090 saying what fixed it
status: cancelled
parent: S-0292
owner: alex
created: 2026-10-06T11:34:29Z
updated: 2026-10-06T11:35:51Z
transitions:
  - to: cancelled
    at: 2026-10-06T11:35:51Z
    by: agent-S-0292
stream: S-0292
tags: [issues]
touches: [design/issues/I-0090-flai-serve-s-orchestrator-test-that-waits-for-the-stopped-run-s-activity-fails-under-the-close-out-s-full-parallel-test-run.md, design/issues/summary.md]
after: [T-1018]
---
# T-1019 Close I-0090 saying what fixed it

## Work

In the story's worktree, run `flai issue close I-0090 --reason "<the cause and the fix>"`, naming the cause T-1017 found and the change T-1018 made. It waits for T-1018 because the reason must name a fix that is built and passes.

## Done when

- I-0090's status is closed, with a reason naming the cause and the fix.
- `design/issues/summary.md` no longer lists I-0090 as open.
- `flai check --strict` is clean for the story.
- The second criterion of S-0292 is ticked.

## Notes

Planned by planner-S-0292. Run it in the worktree, not the main checkout, so that the issue's change lands on the story branch.
- 2026-10-06T11:35:51Z: moved to cancelled: T-1016, written before the planner's draft arrived, already reproduced, fixed, and closes I-0090 in one task
