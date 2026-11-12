---
id: T-3143
type: task
nature: improvement
title: The design and the users guide say a closing move ends a block and flai unblock repairs a closed item, and I-0347 is closed
status: backlog
parent: S-0723
owner: alex
created: 2026-11-12T15:53:59Z
updated: 2026-11-12T15:53:59Z
transitions: []
stream: S-0723
tags: [docs]
touches: [design/system/workflow.md, docs/users/flai.md, design/issues/I-0347-flai-check-finds-item-blocked-closed-outside-the-story-at-close-out.md, design/issues/summary.md]
after: [T-3141, T-3142]
---
# T-3143 The design and the users guide say a closing move ends a block and flai unblock repairs a closed item, and I-0347 is closed

## Work

- `design/system/workflow.md` § Blocking: a move to `done` or `cancelled` closes an open interval at the move, with a warning (T-3141), and `flai unblock` on a closed item, archived or not, closes one left open at the item's close (T-3142). Name S-0723 and I-0347.
- `docs/users/flai.md` § Blocking: the same, for the operator, with the repair command for S-0678 written out: `flai unblock S-0678`, run in the main checkout once this story is accepted and the host's flai carries it.
- Close I-0347 with `flai issue close I-0347 --reason "..."`, naming T-3141's change to `Repo.Move` as the fix of the cause and `flai unblock S-0678` as the repair of the one record it left. That regenerates `design/issues/summary.md`.
- Write in the story's Notes, outside `### Planning`, that the operator runs `flai unblock S-0678` in the main checkout after acceptance, since the story's branch cannot write the main checkout's archive.

Waits for T-3141 and T-3142: it describes what they built, and closes the issue only once both are done.

## Done when

- Both documents say what T-3141 and T-3142 built, and the markdown lint passes on them.
- I-0347 is closed with the reason, and `design/issues/summary.md` no longer lists it as open.
- The story's Notes name the repair command for the operator.

## Notes
