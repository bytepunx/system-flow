---
id: ADR-0296
title: "A close-out records no item.blocked-closed, and a check scoped to a story leaves out every item.blocked-closed outside it"
status: accepted
date: 2026-11-13
supersedes: []
superseded_by: []
refines: [ADR-0085]
---

# ADR-0296 A close-out records no item.blocked-closed, and a check scoped to a story leaves out every item.blocked-closed outside it

## Context

[ADR-0085](0085-a-close-out-s-flai-check-reports-findings-outside-the-story-as-notes-and.md) scoped a close-out's `flai check` to its story and recorded each rule's findings outside the story in one issue per rule. Later ADRs took out of that recording the findings the closing story's agent cannot clear, among them [ADR-0137](0137-a-close-out-records-no-threads-archived-and-a-check-scoped-to-a-story-leaves.md) a `threads.archived` outside the story and [ADR-0256](0256-a-close-out-records-no-narrative-missing-and-a-check-scoped-to-a-story-leaves.md) a `narrative.missing` on another story.

[I-0354](../issues/I-0354-flai-check-finds-item-blocked-closed-outside-the-story-at-close-out.md), ``flai check finds `item.blocked-closed` outside the story at close-out``, has four instances, from the close-outs of S-0723, S-0659, S-0762, and S-0733 between 2026-11-12 and 2026-11-13, and I-0347 before it had 36. All have the same cause.

1. **The finding names another item.** Each is on the archived S-0678: "done item still has an open blocked interval; close it at the item's close with flai unblock S-0678". `item.blocked-closed` warns on every done or cancelled item, archived ones included, whose last `blocked` interval has no `until`. S-0678 was blocked at 2026-11-11T05:26:16Z and accepted at 06:15:07Z with the interval still open. None of the closing stories touched it.
2. **Only `flai unblock` in the main checkout clears it.** Work items live in `wip/` on the main branch, written in the main checkout, so no story branch closes another item's interval. The closing story's agent can fix nothing in its worktree.
3. **S-0723 stopped flai making new ones, not the hand edit.** Since S-0723 a move to `done` or `cancelled` closes an open interval, and `flai unblock` on a closed item closes one left open at the item's close. A hand edit of front matter, or an older flai, can still leave one open, and each would reach every close-out until someone runs `flai unblock`.
4. **Recording it adds nothing.** Each close-out bumped the issue with the same finding on the same item, which counts how many close-outs ran while S-0678 waited, not a cost any closing story could remove.

S-0761 repaired S-0678 with `flai unblock S-0678` in the main checkout, which closed the interval at S-0678's close.

## Decision

A close-out records no `item.blocked-closed` in an issue: a check scoped to a story leaves out every `item.blocked-closed` outside the story, as ADR-0137 does with a `threads.archived` and ADR-0256 with a `narrative.missing`.

1. **A check scoped to a story leaves it out.** `flai check --story S-nnnn` drops from its result every finding of rule `item.blocked-closed` whose path is not inside the story, and takes back the counts it added. The finding's path is the closed item's file, so another item's is outside.
2. **One inside the story stays.** An `item.blocked-closed` on the story's own item file or one of its tasks', or on an item file its branch changes or that is uncommitted in its worktree, is inside the story and is reported as before.
3. **`--record-issues` records no `item.blocked-closed` outside the story.** What is left out of the result is not recorded.
4. **Without `--story` every `item.blocked-closed` is reported**, so the main checkout's `flai check` still warns on each closed item with an open interval, naming `flai unblock`.

This refines ADR-0085, whose rule that every finding outside the story is a note recorded in an issue now leaves out `item.blocked-closed`, and extends the treatment of ADR-0137 and ADR-0256 to one more finding no story branch clears.

## Consequences

- A close-out no longer bumps an issue because an item it did not touch was closed with an open blocked interval, and I-0354 can be closed: its cause no longer occurs.
- The closing story's agent no longer sees, at its close-out, a finding only `flai unblock` in the main checkout clears. The operator still sees it in the main checkout's `flai check`.
- An open issue for `item.blocked-closed` recorded by an older flai stays as it is until it is closed.

## Alternatives considered

- **Repair alone.** `flai unblock S-0678` clears every instance today, and S-0723 stops flai's own moves making new ones. But a hand edit of front matter, or a move by an older flai on another host, can still leave an interval open, and each would again reach every close-out, which can do nothing about it. The repair is done as well, so the main checkout's check is clean.
- **Keep it as a note but record none**, as [ADR-0115](0115-a-close-out-records-no-wip-overlap-and-notes-only-an-overlap-naming-its-story.md) does with a `wip.overlap` naming the story. The note would name another item that no close-out can act on, and the main checkout's check already reports it, where the operator sees it.
