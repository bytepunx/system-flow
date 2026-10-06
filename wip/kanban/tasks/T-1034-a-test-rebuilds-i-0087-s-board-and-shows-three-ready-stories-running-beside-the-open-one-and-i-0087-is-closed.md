---
id: T-1034
type: task
nature: remediation
title: A test rebuilds I-0087's board and shows three ready stories running beside the open one, and I-0087 is closed
status: backlog
parent: S-0295
owner: alex
created: 2026-10-06T12:17:02Z
updated: 2026-10-06T12:17:02Z
transitions: []
stream: S-0295
tags: [flai]
touches: [flai/internal/workitem/hold_i0087_test.go, design/issues/I-0087-one-in-progress-story-holds-every-ready-story-by-overlap-through-folder-wide-touches-and-docs-users-flai-md-so-the-board-runs-one-story-at-a-time-under-a-limit-of-three.md, design/issues/summary.md]
after: [T-1033]
---
# T-1034 A test rebuilds I-0087's board and shows three ready stories running beside the open one, and I-0087 is closed

## Work

Write `flai/internal/workitem/hold_i0087_test.go`, which rebuilds the board of 2026-10-06 from I-0087's instances, as story and task items in memory:

- an in-progress story claiming `flai/internal/mcpserver`, `flai/internal/harness`, and `docs/users`, whose tasks name files inside them;
- a story in review with folder-wide touches;
- ready stories whose touches meet those only through other files of the same folders, `docs/users/flai.md`, and `design/adrs`.

With this project's shared list from `system-flow.yaml`, it asserts that every ready story that overlaps no named file is clear, and that `FirstClear` offers them in pull order. With the old rules (the list empty, review holding, no narrowing), it asserts the full hold the issue describes, so the test shows the cause and the fix together.

Then close the issue with `flai issue close I-0087 --reason`, run in the story's worktree. The reason names the three rules, the ADR, and this test.

This task waits for T-1033, so that the issue closes on the finished story.

## Done when

- The test passes, and it fails when any one of the three rules is reverted.
- I-0087 is closed with its reason, and `design/issues/summary.md` lists it as closed.
- `scripts/flai-test.sh` and `flai check --strict` pass.

## Notes
