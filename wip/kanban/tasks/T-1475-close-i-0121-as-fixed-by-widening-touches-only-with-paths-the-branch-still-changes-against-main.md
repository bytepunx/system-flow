---
id: T-1475
type: task
nature: remediation
title: Close I-0121 as fixed by widening touches only with paths the branch still changes against main
status: backlog
parent: S-0370
owner: alex
created: 2026-10-12T17:10:11Z
updated: 2026-10-12T17:10:11Z
transitions: []
stream: S-0370
tags: [cli]
touches: [design/issues/I-0121-flai-task-done-widens-touches-with-a-path-the-commit-reverted-to-main-re-growing-a-claim-the-story-narrowed.md, design/issues/summary.md]
after: [T-1473, T-1474]
---
# T-1475 Close I-0121 as fixed by widening touches only with paths the branch still changes against main

## Work

Close I-0121 in the story's worktree with `flai issue close I-0121 --reason "<reason>"`, the reason naming S-0370, the new ADR, and the fix: `flai task done` widens the task's and the story's touches only with the committed paths the story's branch still changes against main, so a path a commit reverts to main's version is not added back, with the test T-1473 added. `flai issue close` regenerates `design/issues/summary.md`.

Waits for T-1473 and T-1474: the issue is closed once the fix and its documentation are in.

## Done when

- I-0121 has status closed with that reason, and `design/issues/summary.md` no longer lists it.
- `flai check --strict` is clean.

## Notes

Drafted by the planner for S-0370. This meets the story's criterion 2.
