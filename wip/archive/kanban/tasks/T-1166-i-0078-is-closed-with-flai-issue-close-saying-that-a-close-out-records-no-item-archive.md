---
id: T-1166
type: task
nature: improvement
title: I-0078 is closed with flai issue close, saying that a close-out records no item.archive
status: done
parent: S-0280
owner: alex
created: 2026-10-07T15:03:46Z
updated: 2026-10-07T23:30:02Z
transitions:
  - to: ready
    at: 2026-10-07T23:29:38Z
    by: agent-S-0280
  - to: in-progress
    at: 2026-10-07T23:29:39Z
    by: agent-S-0280
  - to: done
    at: 2026-10-07T23:30:02Z
    by: agent-S-0280
stream: S-0280
tags: [flai]
touches: [design/issues/I-0078-flai-check-finds-item-archive-outside-the-story-at-close-out.md, design/issues/summary.md]
after: [T-1165]
usage:
  source: log
  seconds: 23
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 1
      output: 297
      cache_read: 84007
      cache_write: 2512
      cost: 0.0428
---
# T-1166 I-0078 is closed with flai issue close, saying that a close-out records no item.archive

## Work

Close I-0078 with `flai issue close I-0078 --reason`, in the story's worktree. It waits for T-1165, because the issue closes only once the fix is built and tested.

The reason names:

- the cause: S-0250, cancelled and left unarchived, which every close-out found outside the story;
- the fix: a close-out records no `item.archive` and leaves out one that does not name the story;
- the ADR T-1164 wrote, and the tests T-1165 added.

It shares no path with T-1167 and can run beside it.

## Done when

- I-0078's status is closed, with the reason above.
- `design/issues/summary.md` no longer lists it as open.
- `flai check --strict` passes.

## Notes
