---
id: T-1989
type: task
nature: remediation
title: Close I-0158 with the reason the pinned run times fixed it
status: cancelled
parent: S-0473
owner: alex
created: 2026-10-21T12:07:05Z
updated: 2026-10-21T12:11:07Z
transitions:
  - to: cancelled
    at: 2026-10-21T12:11:07Z
    by: agent-S-0473
stream: S-0473
tags: [issues]
touches: [design/issues/I-0158-testasettledreplystartstheagentwaitingonitonce-in-flai-internal-serve-fails-now-and-then-in-the-close-out-s-integration-tier-on-a-branch-that-does-not-touch-serve.md, design/issues/summary.md]
after: [T-1986]
---
# T-1989 Close I-0158 with the reason the pinned run times fixed it

## Work

Close I-0158 from the story's worktree with `flai issue close I-0158 --commit --reason "<reason>"`, which writes the issue's file and `design/issues/summary.md` and commits them on the story branch. The reason names the cause T-1986 confirmed and the change that removed it. For the planner's reading, that is the two `putEnded` runs straddling a second boundary, so the message read as new to the other story, and the test now pins the other run to the message's second.

It waits for T-1986 (`after`): the reason can only name a fix that is in place and passed its repeated run.

## Done when

- I-0158's status is `closed` and its body ends with the reason.
- `design/issues/summary.md` no longer lists I-0158.
- Both are committed on the story branch.

## Notes

Drafted by the planner. The second acceptance criterion of S-0473.
- 2026-10-21T12:11:07Z: moved to cancelled: duplicate of T-1992, which did the fix and closed I-0158
