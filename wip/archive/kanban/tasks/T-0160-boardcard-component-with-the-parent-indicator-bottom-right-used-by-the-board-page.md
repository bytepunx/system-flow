---
id: T-0160
type: task
nature: improvement
title: BoardCard component with the parent indicator bottom right, used by the board page
status: done
parent: S-0048
owner: alex
created: 2026-09-18T20:41:13Z
updated: 2026-09-18T20:44:10Z
transitions:
  - to: ready
    at: 2026-09-18T20:42:03Z
    by: alex
  - to: in-progress
    at: 2026-09-18T20:42:03Z
    by: alex
  - to: done
    at: 2026-09-18T20:44:10Z
    by: alex
stream: S-0048
tags: []
touches: [flaiover/src/lib/components, flaiover/src/routes/board]
---

# T-0160 BoardCard component with the parent indicator bottom right, used by the board page

## Work
Move the card markup from `src/routes/board/+page.svelte` into `src/lib/components/BoardCard.svelte` with the same link, drag handlers, and classes, taking the card, `draggable`, `dragging`, and the drag callbacks as props. In the bottom row keep nature, type, and BLOCKED on the left and add the parent ID on the right (`ml-auto`, monospace, muted), with `title` and `aria-label` of "<ID> <parent title>" (the ID alone when the title is missing). Render it for stories and tasks that have a parent; nothing for epics or items without one, with no empty element left behind. It is a `span`, not a link, because the card is one `<a>`. Component test beside it, in the style of `AcceptConfirm.svelte.test.ts`: story with a parent, story without one, task, epic, a missing parent title, and that the card is still a single link to the item and draggable.

## Done when
- The component test covers the six cases and passes
- The board page renders cards through the component and its behaviour is unchanged: drag, link, age, blocked flag
- flaiover lint and type check pass

## Notes
