---
id: ADR-0122
title: "A close-out records no item.archive, and a check scoped to a story leaves out every item.archive outside it"
status: proposed
date: 2026-10-07
supersedes: []
superseded_by: []
refines: [ADR-0085, ADR-0115]
---

# ADR-0122 A close-out records no item.archive, and a check scoped to a story leaves out every item.archive outside it

## Context

[ADR-0085](0085-a-close-out-s-flai-check-reports-findings-outside-the-story-as-notes-and.md) scoped a close-out's `flai check` to its story and recorded each rule's findings outside the story in one issue per rule. [ADR-0115](0115-a-close-out-records-no-wip-overlap-and-notes-only-an-overlap-naming-its-story.md) took `wip.overlap` out of that recording, because an overlap is the state of the board, reported elsewhere to the agents who can act on it.

I-0078, ``flai check finds `item.archive` outside the story at close-out``, was bumped 38 times between 2026-10-05 and 2026-10-07. Its instances show one cause.

1. **One lingering item, recorded by every close-out.** Every instance names the same file, S-0250's, with "S-0250 is cancelled; run flai archive". The operator cancelled S-0250 from backlog on 2026-10-05. Cancelling does not archive ([ADR-0028](0028-cancelling-an-item-cancels-everything-open-under-it.md) leaves the narrative in place, and [ADR-0055](0055-a-story-moves-back-one-column-from-ready-in-progress-review-or-cancelled-and.md) lets a cancelled story move back), so the file stayed in `wip/kanban/stories` until commit 1f7b7296 archived it by hand on 2026-10-07. Each close-out in those two days recorded it again. None came after.
2. **Never the closing story's.** `item.archive` warns on a done or cancelled epic or story that is not archived, and never on a task. The closing story is in progress, and so is its epic, so the finding never names either. No story branch makes it: `wip/` is written in the main checkout, and only the operator's `flai archive` there, or an acceptance, clears it. A story's agent can fix nothing in its worktree.

An unarchived item is the state of the main checkout, as an overlap is the state of the board. The main checkout's own `flai check` warns on it, where the operator can act. The issue only counted how many close-outs ran while it lingered.

The question is what a close-out does with an `item.archive`.

## Decision

A close-out records no `item.archive` in an issue: a check scoped to a story leaves out every `item.archive` outside the story, as ADR-0115 does with a `wip.overlap` that does not name the story.

1. **A check scoped to a story leaves out every `item.archive` outside it.** `flai check --story S-nnnn` drops from its result every `item.archive` that is not on the story's own item file, and takes back the counts it added. Since `item.archive` never warns on a story in progress, a scoped check, as a close-out runs it, reports none. Without `--story` every `item.archive` is reported, as before, so the main checkout's check and CI still warn.
2. **`--record-issues` records no `item.archive`.** What is left out of the result is not recorded. The findings of every other rule outside the story are recorded as ADR-0085 says.

This refines ADR-0085, whose rule that every finding outside the story is a note recorded in an issue now leaves out `item.archive`, and ADR-0115, whose treatment of a `wip.overlap` between two other stories it extends to a second rule.

## Consequences

- A close-out no longer bumps an issue because a done or cancelled item waits to be archived, and I-0078 can be closed: its cause no longer occurs.
- The closing story's agent no longer sees, at its close-out, an item it cannot archive. The operator sees it in the main checkout's `flai check` and in CI, as before.
- Nothing counts how long items linger unarchived. A cancelled item that is never archived is the operator's to archive, with `flai archive`; whether cancelling should archive is a question of its own, not this one.
- An open issue for `item.archive` recorded by an older flai stays as it is until it is closed.

## Alternatives considered

- **Keep `item.archive` as a note but record none.** The closing story's agent would still see a finding it cannot act on at every close-out, in every story, until the operator archives the item.
- **Archive on cancel.** It would end the lingering at its source, but it reverses ADR-0055's move back from cancelled and ADR-0028's narrative left in place, and it would not cover a done item an acceptance failed to archive.
- **Record it once, not per story.** An issue that one lingering item opens and no story bumps again still records the state of the main checkout, not friction, and its remediation would be to run `flai archive`.
