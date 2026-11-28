---
id: T-3483
type: task
nature: feature
title: The board's column for a state with runs says what it runs, the suites and the cluster steps, from what board.get gives
status: in-progress
parent: S-0814
owner: alex
created: 2026-11-16T23:51:19Z
updated: 2026-11-28T01:13:49Z
transitions:
  - to: ready
    at: 2026-11-28T01:13:49Z
    by: agent-S-0814
  - to: in-progress
    at: 2026-11-28T01:13:49Z
    by: agent-S-0814
stream: S-0814
tags: [flaiover]
touches: [flaiover/src/lib/server/board.ts, flaiover/src/lib/lanes.ts, flaiover/src/lib/lanes.test.ts, flaiover/src/routes/board/+page.svelte, flaiover/src/routes/board/board.svelte.test.ts]
after: [T-3476]
---
# T-3483 The board's column for a state with runs says what it runs, the suites and the cluster steps, from what board.get gives

## Work

Criterion 3 on the page. Read the verify, runs, and verify_at blocks T-3476 adds to `board.get` in `server/board.ts`. In the lane header of `board/+page.svelte`, for a column whose state has `runs:`, say what entering it runs: its suites, or every suite for `all`, and its cluster steps, such as `runs all suites, deploy, checks` on the baseline's review. Keep the words in a helper in `lanes.ts`, so the page and its test share them. A column with no `runs:` shows nothing new.

Waits for T-3476, whose answer this reads.

## Done when

- The baseline's review lane says it deploys and checks, a fixture's integration lane names its suites, and a lane with no `runs:` is unchanged.
- `lanes.test.ts` and `board.svelte.test.ts` cover these, and `flai test` on the changed paths passes.

## Notes

- Drafted by planner-S-0814. S-0753 makes the board draw its columns from flai; this task builds on that and adds no list of states of its own.
