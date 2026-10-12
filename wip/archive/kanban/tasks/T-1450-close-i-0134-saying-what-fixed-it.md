---
id: T-1450
type: task
nature: improvement
title: Close I-0134 saying what fixed it
status: done
parent: S-0365
owner: alex
created: 2026-10-12T14:56:45Z
updated: 2026-10-12T15:06:14Z
transitions:
  - to: ready
    at: 2026-10-12T14:58:37Z
    by: agent-S-0365
  - to: in-progress
    at: 2026-10-12T15:06:09Z
    by: agent-S-0365
  - to: done
    at: 2026-10-12T15:06:14Z
    by: agent-S-0365
stream: S-0365
tags: [issues]
touches: [design/issues/I-0134-flai-check-finds-issues-no-story-outside-the-story-at-close-out.md, design/issues/summary.md]
after: [T-1448, T-1449]
usage:
  source: log
  seconds: 5
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 4
      output: 907
      cache_read: 184886
      cache_write: 8249
      cost: 0.1108
---
# T-1450 Close I-0134 saying what fixed it

## Work

In the story's worktree, run `flai issue close I-0134 --reason "<reason>"`. The reason names T-1447's ADR and T-1448's change: a check scoped to a story leaves `issues.no-story` outside the story out, so a close-out no longer notes it or records it. The command regenerates `design/issues/summary.md`.

It waits for T-1448 and T-1449, so that it closes the issue only once the fix and its documents are in.

## Done when

- I-0134 is closed with a reason naming the fix, and `design/issues/summary.md` no longer lists it.
- `flai check --strict` passes on the story.

## Notes
