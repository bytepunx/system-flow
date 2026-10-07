---
id: T-1166
type: task
nature: improvement
title: I-0078 is closed with flai issue close, saying that a close-out records no item.archive
status: backlog
parent: S-0280
owner: alex
created: 2026-10-07T15:03:46Z
updated: 2026-10-07T15:04:31Z
transitions: []
stream: S-0280
tags: [flai]
touches: [design/issues/I-0078-flai-check-finds-item-archive-outside-the-story-at-close-out.md, design/issues/summary.md]
after: [T-1165]
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
