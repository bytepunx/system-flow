---
id: T-1438
type: task
nature: remediation
title: Close I-0126 as fixed by S-0322 in flai 1.39.14, its instance having run on 1.39.13
status: backlog
parent: S-0362
owner: alex
created: 2026-10-12T14:20:55Z
updated: 2026-10-12T14:20:55Z
transitions: []
stream: S-0362
tags: [cli]
touches: [design/issues/I-0126-flai-task-done-commits-the-changed-paths-of-another-open-task-of-the-same-layer.md, design/issues/summary.md]
after: [T-1437]
---
# T-1438 Close I-0126 as fixed by S-0322 in flai 1.39.14, its instance having run on 1.39.13

## Work

In the story's worktree, run `flai issue close I-0126 --reason "<reason>"`. The command writes under the checkout it runs in. The reason says:

- I-0126's instance closed T-1377 at 2026-10-08T10:35:55Z through an MCP server running the installed flai 1.39.13.
- S-0322 implemented ADR-0128 and was published in flai 1.39.14 at 09:30Z the same day, so 1.39.13 still committed every changed path.
- T-1437's tests show that the current `flai task done` and MCP `task_done` leave a sibling's paths and keep them out of the closing task's touches. Name the tests.
- Name any fix T-1437 had to make in `taskdone.go`.
- The guard against an MCP flai older than the project's history is S-0363's, from I-0127.

Waits for T-1437: the reason names its tests and what they found.

## Done when

- I-0126's status is closed, and its reason names S-0322, flai 1.39.14, the tests, and S-0363.
- `design/issues/summary.md` no longer lists I-0126.
- The markdown lint passes on both files.

## Notes

Drafted by the planner for S-0362.
