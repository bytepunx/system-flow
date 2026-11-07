---
id: T-2509
type: task
nature: remediation
title: "Reproduce I-0180 in Playwright: two long-named projects at 390px leave the page no wider than the window"
status: done
parent: S-0557
owner: alex
created: 2026-11-07T03:46:11Z
updated: 2026-11-07T04:00:28Z
transitions:
  - to: ready
    at: 2026-11-07T03:52:47Z
    by: agent-S-0557
  - to: in-progress
    at: 2026-11-07T03:52:51Z
    by: agent-S-0557
  - to: done
    at: 2026-11-07T04:00:28Z
    by: agent-S-0557
stream: S-0557
tags: [flaiover, test]
touches: [flaiover/e2e/project-switcher.e2e.ts]
after: [T-2507]
usage:
  source: log
  seconds: 457
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 27
      output: 3889
      cache_read: 2761738
      cache_write: 47445
      cost: 1.0038
---
# T-2509 Reproduce I-0180 in Playwright: two long-named projects at 390px leave the page no wider than the window

## Work

Add `flaiover/e2e/project-switcher.e2e.ts`. Sign in as `flaiover/e2e/layout.e2e.ts` does, and answer `/api/projects` in the browser with `context.route`, as `flaiover/e2e/board-layout.e2e.ts` answers `/api/board`, with two projects whose names are as long as I-0180's ("E2E board (unsigned allowed: flai and dashboard)"), one connected and one not, so the "not connected: why?" link shows too. At 390px and at 1440px, on `/` and `/board`, check that the page does not scroll sideways (`scrollWidth - innerWidth` at most 0) and that the select's box ends inside the window. Route interception keeps the second project off the server, so the other test files still render with none. Waits for T-2507: its done condition is that the test passes against the fixed switcher.

## Done when

- The test fails on the switcher as it was before T-2507, the page wider than the window at 390px, and passes with T-2507's change; say both runs in this task's Notes.
- `scripts/flaiover-e2e.sh e2e/project-switcher.e2e.ts e2e/layout.e2e.ts e2e/board-layout.e2e.ts` passes.

## Notes

Runs on 2026-11-07, in S-0557's worktree:

| Switcher | Command | Result |
|----------|---------|--------|
| main's, before T-2507 (`git show main:` over the file, restored after) | `scripts/flaiover-e2e.sh e2e/project-switcher.e2e.ts` | 2 failed, 2 passed, exit 1: both 390px tests, the select's box ending at 548px past a 390px window; the 1440px tests pass |
| T-2507's | `scripts/flaiover-e2e.sh e2e/project-switcher.e2e.ts e2e/layout.e2e.ts e2e/board-layout.e2e.ts` | 22 passed, exit 0 |

The test was written by a task sub-agent and reviewed by the story's agent. The current project is chosen with `?project=e2e-board`, and `/api/board` is answered with an empty board, so `/board` renders without a project the server knows.
