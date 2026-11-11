---
id: T-3113
type: task
nature: remediation
title: Close I-0343 with what fixed it
status: backlog
parent: S-0715
owner: alex
created: 2026-11-11T09:04:44Z
updated: 2026-11-11T09:04:44Z
transitions: []
stream: S-0715
tags: [issues]
touches: [design/issues/I-0343-testarestartstartswhatisreadyandnothingtwice-fails-its-temp-folder-s-cleanup.md, design/issues/summary.md]
after: [T-3112]
---
# T-3113 Close I-0343 with what fixed it

## Work

Close I-0343 with `flai issue close I-0343 --commit --reason "<reason>"` from the story's worktree. The reason names T-3112's commit, the cause it fixed (the lab's `settle` did not wait for the second launcher the restart test made), and the race-run counts T-3112 recorded. This task waits for T-3112, since the reason cites its commit and its counts.

## Done when

- I-0343's status is `closed`, and its body ends with the reason.
- `design/issues/summary.md` no longer lists I-0343 as open.
- `flai check --strict` passes on the story.

## Notes
