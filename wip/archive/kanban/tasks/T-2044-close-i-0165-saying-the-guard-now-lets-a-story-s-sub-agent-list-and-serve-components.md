---
id: T-2044
type: task
nature: feature
title: Close I-0165 saying the guard now lets a story's sub-agent list and serve components
status: cancelled
parent: S-0477
owner: alex
created: 2026-10-22T21:00:11Z
updated: 2026-10-22T21:13:37Z
transitions:
  - to: cancelled
    at: 2026-10-22T21:13:37Z
    by: agent-S-0477
stream: S-0477
tags: []
touches: [design/issues/I-0165-flai-guard-refuses-a-story-s-sub-agent-flai-components-list-and-flai-components-serve-so-a-headless-pass-of-the-harness-falls-to-the-story-s-agent.md, design/issues/summary.md]
after: [T-2043]
---

# T-2044 Close I-0165 saying the guard now lets a story's sub-agent list and serve components

## Work

Waits for T-2043: the issue is closed only once the guard lets a sub-agent list and serve. Close I-0165 with `flai issue close I-0165 --reason ... --commit` from the worktree, naming the fix and its ADR.

## Done when

- I-0165 is closed with a reason naming S-0477's guard change, committed on `story/S-0477`.

## Notes
- 2026-10-22T21:13:37Z: moved to cancelled: Duplicates the planner's T-2036, which closes I-0165 after T-2032 and T-2034
