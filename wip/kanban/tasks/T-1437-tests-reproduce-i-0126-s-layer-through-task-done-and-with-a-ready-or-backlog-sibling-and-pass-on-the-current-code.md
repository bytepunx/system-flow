---
id: T-1437
type: task
nature: remediation
title: Tests reproduce I-0126's layer through task_done and with a ready or backlog sibling, and pass on the current code
status: backlog
parent: S-0362
owner: alex
created: 2026-10-12T14:20:46Z
updated: 2026-10-12T14:20:46Z
transitions: []
stream: S-0362
tags: [cli]
touches: [flai/internal/taskdone/taskdone_test.go, flai/internal/mcpserver/task_test.go]
---
# T-1437 Tests reproduce I-0126's layer through task_done and with a ready or backlog sibling, and pass on the current code

## Work

I-0126's one instance, S-0347's T-1377 and T-1379, ran on the installed flai 1.39.13. S-0322 implemented ADR-0128 and was published in flai 1.39.14, an hour before the instance. The current `commit` in `flai/internal/taskdone/taskdone.go` leaves a path that only another open task covers. Existing tests already cover two in-progress tasks closing apart: `TestRunClosesTwoTasksOfOneStoryApart` in `taskdone_test.go`, and `TestTaskDoneLeavesAnotherOpenTasksPaths` in `mcpserver/task_test.go`. They leave three gaps. Fill them so a regression shows:

- In `flai/internal/mcpserver/task_test.go`, extend `TestTaskDoneLeavesAnotherOpenTasksPaths`, or add a test beside it, so that it reproduces I-0126's shape through the MCP tool. Use two in-progress tasks with disjoint touches: one has code files, the other has a convention, a template copy of it, and a changelog. Assert that the closing task's touches and the story's touches are not widened with the sibling's paths, and that the sibling's paths stay uncommitted and are listed in `left`.
- In `flai/internal/taskdone/taskdone_test.go`, add a table test with siblings in `ready` and in `backlog`. Each sibling's paths are left. Add a `done` sibling and a `cancelled` sibling too: their paths are committed with the closing task, because neither is open.
- Run the new tests on the current code. If one fails, fix `commit` or `coverage` in `taskdone.go` within ADR-0128. `task_done` then widens this task's touches to it. Otherwise change no code.

Waits for no task: it is the first layer.

## Done when

- The new and extended tests pass under `flai test flai/internal/taskdone flai/internal/mcpserver`.
- The MCP test fails if `commit` is changed back to staging every changed path, as 1.39.13 did. Check this once by hand, then revert the change.

## Notes

Drafted by the planner for S-0362. The guard against an outdated MCP flai, the cause behind this instance, is S-0363's (I-0127). Do not build it here.
