---
id: ADR-0256
title: "A close-out records no narrative.missing, and a check scoped to a story leaves out every narrative.missing outside it"
status: accepted
date: 2026-11-07
supersedes: []
superseded_by: []
refines: [ADR-0085]
topics: [cli, conventions, template]
---

# ADR-0256 A close-out records no narrative.missing, and a check scoped to a story leaves out every narrative.missing outside it

## Context

[ADR-0085](0085-a-close-out-s-flai-check-reports-findings-outside-the-story-as-notes-and.md) scoped a close-out's `flai check` to its story and recorded each rule's findings outside the story in one issue per rule. Later ADRs took out of that recording the findings the closing story's agent cannot clear: [ADR-0125](0125-a-close-out-records-no-narrative-state-finding-on-another-open-story-s.md) a `narrative.state` on another open story's narrative, [ADR-0137](0137-a-close-out-records-no-threads-archived-and-a-check-scoped-to-a-story-leaves.md) a `threads.archived` outside the story, and [ADR-0255](0255-a-check-scoped-to-a-story-leaves-out-the-agent-image-findings-only-the-host.md) the `agent.image` findings only the host clears.

[I-0235](../issues/I-0235-flai-check-finds-narrative-missing-outside-the-story-at-close-out.md), ``flai check finds `narrative.missing` outside the story at close-out``, has five instances, from the close-outs of S-0548, S-0547, S-0550, S-0546, and S-0551 between 2026-10-31 and 2026-11-01. All five have the same cause.

1. **The finding names another story.** Each is on S-0549's item file: "S-0549 is in-progress but has no narrative; run flai stream open S-0549". `narrative.missing` is an error on every story in progress or in review whose narrative, `wip/agents/<id>.md`, does not exist. None of the five stories touched S-0549.
2. **A move opens no narrative.** The operator moved S-0549 to in progress at 2026-10-31T09:47:53Z. That move went through `flai move`, which changes the status and nothing else. Only `flai story start` and `flai stream open` open a narrative, through `OpenStream`.
3. **Only that story's agent or the operator can clear it.** Narratives live in `wip/` on the main branch, written in the main checkout, so no other story's branch opens one. The closing story's agent can fix nothing in its worktree, and none of the five did.
4. **Recording it adds nothing.** Each close-out bumped I-0235 with the same finding on the same story. That counts how many close-outs ran while S-0549 had no narrative, not a cost any closing story could remove.

## Decision

A close-out records no `narrative.missing` in an issue: a check scoped to a story leaves out every `narrative.missing` outside the story, as ADR-0137 does with a `threads.archived` and ADR-0125 with a `narrative.state` on another open story's narrative.

1. **A check scoped to a story leaves it out.** `flai check --story S-nnnn` drops from its result every finding of rule `narrative.missing` whose path is not inside the story, and takes back the counts it added. The finding's path is the item file of the story that lacks a narrative, so another story's is outside.
2. **One on the story itself stays.** A `narrative.missing` on the story's own item file is inside the story and is reported as before, an error, so the story's own close-out stops on it.
3. **`--record-issues` records no `narrative.missing` outside the story.** What is left out of the result is not recorded.
4. **Without `--story` every `narrative.missing` is reported**, so the main checkout's `flai check` still fails on each story in progress or in review with no narrative.

This refines ADR-0085, whose rule that every finding outside the story is a note recorded in an issue now leaves out `narrative.missing`, and extends the treatment of ADR-0125, ADR-0137, and ADR-0255 to one more finding no story branch clears.

## Consequences

- A close-out no longer bumps an issue because another story was moved to in progress without its narrative, and I-0235 can be closed: its cause no longer occurs.
- The closing story's agent no longer sees, at its close-out, a finding only the other story's agent or the operator can act on, with `flai stream open` on that story.
- Two places still report it: `flai check` in the main checkout, and the close-out of the story that lacks the narrative, whose own scoped check keeps it inside.
- An open issue for `narrative.missing` recorded by an older flai stays as it is until it is closed.

## Alternatives considered

- **Have a move to in progress open the narrative**, from `flai move`, the dashboard, and `item_move`, as `flai story start` does. It removes the cause, but it changes what a move does for every caller. A move would open a branch and a worktree for a story nobody has started, and an operator's move back and forth would leave them behind. `flai story start` and `flai stream open` already do that on purpose, and a story's agent starts that way. It would also leave a narrative missing for any other reason, such as one deleted by hand, still reaching every close-out.
- **Keep it as a note but record none**, as [ADR-0115](0115-a-close-out-records-no-wip-overlap-and-notes-only-an-overlap-naming-its-story.md) does with a `wip.overlap` naming the story. The note would name another story that no close-out can act on, and the main checkout's check already fails on it, where the operator sees it.
