---
id: ADR-0133
title: "A close-out records no board.wip-limit, and a check scoped to a story leaves out every board.wip-limit outside it"
status: proposed
date: 2026-10-08
supersedes: []
superseded_by: []
refines: [ADR-0085, ADR-0122]
---

# ADR-0133 A close-out records no board.wip-limit, and a check scoped to a story leaves out every board.wip-limit outside it

## Context

[ADR-0085](0085-a-close-out-s-flai-check-reports-findings-outside-the-story-as-notes-and.md) scoped a close-out's `flai check` to its story and recorded each rule's findings outside the story in one issue per rule. [ADR-0122](0122-a-close-out-records-no-item-archive-and-a-check-scoped-to-a-story-leaves-out.md), [ADR-0123](0123-a-close-out-records-no-markdown-finding-on-another-open-story-s-narrative-and-a.md), and [ADR-0125](0125-a-close-out-records-no-narrative-state-finding-on-another-open-story-s.md) took out of that recording the findings the closing story's agent cannot fix.

I-0123, ``flai check finds `board.wip-limit` outside the story at close-out``, has two instances, from the close-outs of S-0321 and S-0339 on 2026-10-08. Both have the same cause.

1. **The board's counts are the main checkout's state.** Each finding is on `wip/kanban/board.md`: "6 stories in in-progress, limit 5". `board.wip-limit` counts the stories in each column of the main checkout's `wip/`, which no story branch writes. A start the operator made past the limit, or a story sent back, takes a column over it.
2. **No story branch clears it, and it names no story.** Only the operator's acceptances, cancellations, or a raised limit clear it. The closing story's agent can fix nothing in its worktree, and moving its own story to review, which would lower the count, comes after the close-out.

[ADR-0073](0073-a-full-review-holds-the-pull-and-flai-check-strict-passes-over-review-over-its.md) already made review over its limit advisory, so that a close-out does not stop on the operator's queue. Ready and in-progress over their limits still warn, and a close-out marked them outside and recorded them in I-0123.

## Decision

A close-out records no `board.wip-limit` in an issue: a check scoped to a story leaves out every `board.wip-limit` outside the story, as ADR-0122 does with an `item.archive`.

1. **A check scoped to a story leaves it out.** `flai check --story S-nnnn` drops from its result every finding of rule `board.wip-limit` whose path is not inside the story, and takes back the counts it added, advisory included. No story branch changes `wip/kanban/board.md`, so a scoped check, as a close-out runs it, reports none.
2. **`--record-issues` records no `board.wip-limit`.** What is left out of the result is not recorded.
3. **Without `--story` every `board.wip-limit` is reported**, so the main checkout's `flai check --strict` and CI still fail on ready or in-progress over its limit, and warn on review over its limit, as ADR-0073 says.

This refines ADR-0085, whose rule that every finding outside the story is a note recorded in an issue now leaves out `board.wip-limit`, and ADR-0122, whose treatment of an `item.archive` it extends to a second rule about the main checkout's state.

## Consequences

- A close-out no longer bumps an issue because a column of the board is over its limit, and I-0123 can be closed: its cause no longer occurs.
- The closing story's agent no longer sees, at its close-out, a breach it cannot clear. The operator sees it on the board, in `flai board`, and in the main checkout's `flai check`, as before.
- Nothing counts how often close-outs ran while a column was over its limit. A breach is recorded where it happens: `flai serve`, `inbox`, and `wait_for_work` hold the pull while in-progress or review is full.
- An open issue for `board.wip-limit` recorded by an older flai stays as it is until it is closed.

## Alternatives considered

- **Keep it as a note but record none.** The closing story's agent would still see a finding it cannot act on at every close-out while the column stays over its limit.
- **Make in-progress over its limit advisory, as review is.** It would change what the main checkout's `flai check --strict` and CI report, which is the operator's view of the breach, not only the close-out's.
- **Count the closing story out of in-progress.** The story is still in progress during its close-out, and the finding would still fire whenever two other stories go past the limit.
