---
id: ADR-0137
title: "A close-out records no threads.archived, and a check scoped to a story leaves out every threads.archived outside it"
status: proposed
date: 2026-10-12
supersedes: []
superseded_by: []
refines: [ADR-0085, ADR-0122, ADR-0133, ADR-0136]
---

# ADR-0137 A close-out records no threads.archived, and a check scoped to a story leaves out every threads.archived outside it

## Context

[ADR-0085](0085-a-close-out-s-flai-check-reports-findings-outside-the-story-as-notes-and.md) scoped a close-out's `flai check` to its story and recorded each rule's findings outside the story in one issue per rule. [ADR-0122](0122-a-close-out-records-no-item-archive-and-a-check-scoped-to-a-story-leaves-out.md), [ADR-0133](0133-a-close-out-records-no-board-wip-limit-and-a-check-scoped-to-a-story-leaves-out.md), and [ADR-0136](0136-a-close-out-records-no-issues-no-story-and-a-check-scoped-to-a-story-leaves-out.md) took out of that recording an `item.archive`, a `board.wip-limit`, and an `issues.no-story` outside the story: findings on the main checkout's state that no story branch clears.

I-0135, ``flai check finds `threads.archived` outside the story at close-out``, has three instances, from the close-outs of S-0342, S-0343, and S-0289 between 2026-10-08 and 2026-10-09. All three have the same cause.

1. **The finding names a thread, not the story.** Each is on TH-0389: "TH-0389 is open but S-0338 is archived; resolve it with flai thread resolve TH-0389". `threads.archived` warns on every open or answered thread whose anchor item is archived. None of the three stories touched S-0338 or TH-0389.
2. **Archiving could not resolve it.** The orchestrator opened TH-0389 on S-0338 at 2026-10-08T20:55:51Z, after S-0338's acceptance had archived it, so [ADR-0109](0109-archiving-an-item-resolves-the-threads-still-open-on-it.md)'s resolve-on-archive had already run. The thread was open for a reason: it awaited the operator.
3. **Only the operator or the thread's opener clears it.** Threads live in `wip/` on the main branch, written in the main checkout, so no story branch resolves one. The closing story's agent can fix nothing in its worktree, and none of the three did.
4. **Recording it adds nothing.** Each close-out bumped I-0135 with the same finding on the same thread, which counts how many close-outs ran while TH-0389 waited, not a cost any story could remove.

## Decision

A close-out records no `threads.archived` in an issue: a check scoped to a story leaves out every `threads.archived` outside the story, as ADR-0122, ADR-0133, and ADR-0136 do with an `item.archive`, a `board.wip-limit`, and an `issues.no-story`.

1. **A check scoped to a story leaves it out.** `flai check --story S-nnnn` drops from its result every finding of rule `threads.archived` whose path is not inside the story, and takes back the counts it added.
2. **One inside the story stays.** A `threads.archived` on a thread anchored on the story or one of its tasks, or on a thread file the story's branch changes or that is uncommitted in its worktree, is inside the story and is reported as before.
3. **`--record-issues` records no `threads.archived` outside the story.** What is left out of the result is not recorded.
4. **Without `--story` every `threads.archived` is reported**, so the main checkout's `flai check` still warns on each open thread on an archived item.

This refines ADR-0085, whose rule that every finding outside the story is a note recorded in an issue now leaves out `threads.archived`, and ADR-0122, ADR-0133, and ADR-0136, whose treatment it extends to a fourth rule.

## Consequences

- A close-out no longer bumps an issue because a thread it did not touch is open on an archived item, and I-0135 can be closed: its cause no longer occurs.
- The closing story's agent no longer sees, at its close-out, a finding only the operator or the thread's opener can act on. They still see it in the main checkout's `flai check`, and `inbox` and the dashboard list the open thread to whoever it awaits.
- An open issue for `threads.archived` recorded by an older flai stays as it is until it is closed.

## Alternatives considered

- **Stop `threads.archived` firing on a thread opened after its item was archived, in every check.** It fixes the cause of TH-0389's instances, but not a thread left open by an archive that stopped part way, which would still reach every close-out. It would also hide a thread opened on an archived item from the main checkout's check, where the operator should see it until it is resolved.
- **Keep it as a note but record none**, as [ADR-0115](0115-a-close-out-records-no-wip-overlap-and-notes-only-an-overlap-naming-its-story.md) does with a `wip.overlap` naming the story. Three close-outs saw it and none could resolve it: the note would only repeat what the thread's own reader already sees.
