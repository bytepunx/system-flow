---
id: T-1425
type: task
nature: improvement
title: Close I-0123 saying what fixed it
status: backlog
parent: S-0348
owner: alex
created: 2026-10-08T08:59:12Z
updated: 2026-10-08T08:59:12Z
transitions: []
stream: S-0348
tags: [issues]
touches: [design/issues/I-0123-flai-check-finds-board-wip-limit-outside-the-story-at-close-out.md, design/issues/summary.md]
after: [T-1423, T-1424]
---
# T-1425 Close I-0123 saying what fixed it

## Work

In the story's worktree, run `flai issue close I-0123 --reason "<reason>"`, the reason naming the ADR of T-1422 and the change of T-1423: a check scoped to a story leaves `board.wip-limit` out, so a close-out no longer notes it or records it. The command regenerates `design/issues/summary.md`.

It waits for T-1423 and T-1424, so that it closes the issue only once the fix and its documents are in.

## Done when

- I-0123 is closed with a reason naming the fix, and `design/issues/summary.md` no longer lists it.
- `flai check --strict` passes on the story.

## Notes
