---
id: T-1323
type: task
nature: remediation
title: I-0112 is closed with flai issue close, saying that a sync folds a branch's duplicate issue
status: backlog
parent: S-0326
owner: alex
created: 2026-10-08T00:34:04Z
updated: 2026-10-08T00:34:04Z
transitions: []
stream: S-0326
tags: [issues]
touches: [design/issues/I-0112-flai-verify-record-issues-opens-an-issue-on-the-story-branch-that-another-branch-opened-under-the-same-title-meanwhile.md, design/issues/summary.md]
after: [T-1321, T-1322]
---
# T-1323 I-0112 is closed with flai issue close, saying that a sync folds a branch's duplicate issue

## Work

Close I-0112 in the story's worktree with `flai issue close I-0112 --reason`, naming the ADR T-1319 recorded and saying that a sync, a task's close, and an acceptance now fold an open issue the branch added into the main branch's open issue of the same title. `flai issue close` regenerates `design/issues/summary.md`.

It waits for T-1321 and T-1322 because the reason states what they built and documented.

## Done when

- I-0112's status is closed and its reason names the ADR and the fold
- `design/issues/summary.md` no longer lists I-0112
- `flai check --strict` passes

## Notes
