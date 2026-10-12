---
id: ADR-0136
title: "A close-out records no issues.no-story, and a check scoped to a story leaves out every issues.no-story outside it"
status: proposed
date: 2026-10-12
supersedes: []
superseded_by: []
refines: [ADR-0085, ADR-0122, ADR-0133]
---

# ADR-0136 A close-out records no issues.no-story, and a check scoped to a story leaves out every issues.no-story outside it

## Context

[ADR-0085](0085-a-close-out-s-flai-check-reports-findings-outside-the-story-as-notes-and.md) scoped a close-out's `flai check` to its story and recorded each rule's findings outside the story in one issue per rule. [ADR-0122](0122-a-close-out-records-no-item-archive-and-a-check-scoped-to-a-story-leaves-out.md) and [ADR-0133](0133-a-close-out-records-no-board-wip-limit-and-a-check-scoped-to-a-story-leaves-out.md) took out of that recording an `item.archive` and a `board.wip-limit` outside the story: findings on the main checkout's state that no story branch clears.

I-0134, ``flai check finds `issues.no-story` outside the story at close-out``, has four instances, from the close-outs of S-0343, S-0289, S-0342, and S-0361 between 2026-10-08 and 2026-10-12. All four have the same cause.

1. **The finding names an issue, not the story.** Each is on I-0052's file: "I-0052 has been open since 2026-10-01 with no open story linking it: make one with flai issue story I-0052". `issues.no-story` warns on every open issue older than `issues.story_after` that no open story links. None of the four stories touched I-0052.
2. **Only the operator clears it.** Making a story for an issue is the operator's choice ([continuous-improvement.md](../conventions/continuous-improvement.md)): at acceptance, a story's review page lists every open issue no open story links, and the operator checks those that become backlog stories. The closing story's agent can fix nothing in its worktree, and none of the four did.
3. **Recording it adds nothing.** Each close-out bumped I-0134 with the same finding on the same issue, which counts how many close-outs ran while I-0052 waited, not a cost any story could remove.

## Decision

A close-out records no `issues.no-story` in an issue: a check scoped to a story leaves out every `issues.no-story` outside the story, as ADR-0122 and ADR-0133 do with an `item.archive` and a `board.wip-limit`.

1. **A check scoped to a story leaves it out.** `flai check --story S-nnnn` drops from its result every finding of rule `issues.no-story` whose path is not inside the story, and takes back the counts it added.
2. **One inside the story stays.** An `issues.no-story` on an issue file the story's branch changes, or that is uncommitted in its worktree, is inside the story and is reported as before.
3. **`--record-issues` records no `issues.no-story` outside the story.** What is left out of the result is not recorded.
4. **Without `--story` every `issues.no-story` is reported**, so the main checkout's `flai check` still warns on each open issue with no story, and the rule to make one with `flai issue story` when it does stands.

This refines ADR-0085, whose rule that every finding outside the story is a note recorded in an issue now leaves out `issues.no-story`, and ADR-0122 and ADR-0133, whose treatment it extends to a third rule.

## Consequences

- A close-out no longer bumps an issue because an issue it did not touch has no story, and I-0134 can be closed: its cause no longer occurs.
- The closing story's agent no longer sees, at its close-out, a finding only the operator can act on. The operator still sees it in the main checkout's `flai check`, and on each story's review page, which lists every open issue no open story links.
- An open issue for `issues.no-story` recorded by an older flai stays as it is until it is closed.

## Alternatives considered

- **Keep it as a note but record none**, as [ADR-0115](0115-a-close-out-records-no-wip-overlap-and-notes-only-an-overlap-naming-its-story.md) does with a `wip.overlap` naming the story. It would keep the nudge to make a story in front of each story's agent, but four close-outs saw it and none made the story, which is the operator's choice at acceptance.
- **Leave out only an issue that some story in progress or review recorded.** The finding would still fire on any issue no story has touched, as I-0052 had not.
