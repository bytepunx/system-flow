---
id: T-3269
type: task
nature: improvement
title: Record in an ADR whether a close-out leaves out item.blocked-closed outside the story
status: backlog
parent: S-0761
owner: alex
created: 2026-11-13T07:29:36Z
updated: 2026-11-13T07:29:36Z
transitions: []
stream: S-0761
tags: [flai, check, adr]
touches: [design/adrs/README.md]
---
# T-3269 Record in an ADR whether a close-out leaves out item.blocked-closed outside the story

## Work

Write and accept an ADR with `flai adr new --commit`. It refines ADR-0085 and follows ADR-0137 (`threads.archived`) and ADR-0256 (`narrative.missing`). The proposed decision: an `item.blocked-closed` finding outside the story is left out of a check scoped to a story (`flai check --story`), so a close-out neither notes it nor records it in an issue. One on the story itself or its tasks stays inside.

Its context is I-0347 and I-0354. Thirty-seven close-outs between 2026-11-11 and 2026-11-12 found `item.blocked-closed` on the archived S-0678, which none of them touched. S-0723 stopped flai's own moves from making new ones: a move to `done` or `cancelled` closes an open interval. But no story branch can clear one on another item: only `flai unblock` on that item, in the main checkout, does. S-0723 closed I-0347 on its branch, and its own close-out recorded the same finding again as I-0354. The consequences name what still reports it: `flai check` in the main checkout, and the story's own close-out for its own item and tasks.

The alternative is repair alone: unblock S-0678 and change nothing in the check, since only a hand edit of front matter can now leave an interval open. Decide it instead if it fits better, and say why in the ADR. If it is decided, say so on the story's thread before the code and docs tasks are worked, since their work then changes or they are cancelled.

Add the ADR's file to this task's and the story's touches once its title is set. It waits for nothing in this story.

## Done when

- The ADR is accepted under `design/adrs/`, indexed in `design/adrs/README.md`, with its topics set.
- It says whether an `item.blocked-closed` outside the story is left out, why, and what still reports it.

## Notes
