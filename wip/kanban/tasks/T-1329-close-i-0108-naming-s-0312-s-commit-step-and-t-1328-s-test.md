---
id: T-1329
type: task
nature: remediation
title: Close I-0108, naming S-0312's commit step and T-1328's test
status: backlog
parent: S-0322
owner: alex
created: 2026-10-08T04:32:45Z
updated: 2026-10-08T04:32:45Z
transitions: []
stream: S-0322
tags: [flai]
touches: [design/issues/I-0108-flai-task-done-commits-every-uncommitted-file-in-the-story-worktree-so-closing-one-task-of-a-layer-sweeps-its-siblings-work-into-its-commit.md, design/issues/summary.md]
after: [T-1328]
---
# T-1329 Close I-0108, naming S-0312's commit step and T-1328's test

## Work

From the story's worktree, run `flai issue close I-0108 --reason "<what fixed it>"`. The reason names the ADR S-0312 wrote, its commit step in `flai/internal/taskdone/taskdone.go`, and the test T-1328 added or named. `flai issue close` rewrites `design/issues/summary.md`.

It waits for T-1328, so the issue closes on a fix a test shows.

## Done when

- I-0108 is closed with a reason naming what fixed it, and `design/issues/summary.md` no longer lists it.
- `flai check --strict` and the markdown lint are clean.

## Notes

Drafted by the planner for S-0322.
