---
id: ADR-0123
title: "A close-out records no markdown finding on another open story's narrative, and a check scoped to a story leaves every such finding out"
status: proposed
date: 2026-10-08
supersedes: []
superseded_by: []
refines: [ADR-0085, ADR-0122]
---

# ADR-0123 A close-out records no markdown finding on another open story's narrative, and a check scoped to a story leaves every such finding out

## Context

[ADR-0085](0085-a-close-out-s-flai-check-reports-findings-outside-the-story-as-notes-and.md) scoped a close-out's `flai check` to its story and recorded each rule's findings outside the story in one issue per rule. [ADR-0122](0122-a-close-out-records-no-item-archive-and-a-check-scoped-to-a-story-leaves-out.md) took `item.archive` out of that recording, as [ADR-0115](0115-a-close-out-records-no-wip-overlap-and-notes-only-an-overlap-naming-its-story.md) had taken out a `wip.overlap` between two other stories: neither is anything the closing story's agent can fix.

I-0096, ``flai check finds `markdown.MD038` outside the story at close-out``, has one instance. Its cause is a second finding of that kind.

1. **Another story's hand-edited narrative.** On 2026-10-06, S-0227's close-out found `markdown.MD038` in `wip/agents/S-0229.md`, the narrative of S-0229, then in progress: a code span ending in a space, in its `## Decisions`. A narrative's `## Context`, `## Decisions`, and `## Open questions` are edited by hand (`tooling.md`), so flai's wip lint guard ([ADR-0061](0061-flai-lints-the-markdown-it-writes-in-wip-with-its-own-implementation-of-the.md)), which `flai stream log` and `flai stream state` run, never saw it.
2. **Fixed by its own story, not by the one that recorded it.** S-0229's agent fixed it fifteen minutes later. Its own close-out checks its own narrative, since a check scoped to a story counts the story's narrative as inside it, and `--strict` fails on a markdown warning there. S-0227's agent could not fix it, because the narrative is S-0229's, yet its close-out marked the finding outside and recorded it in I-0096.

A lint finding on an open story's narrative is that story's work in progress, and its close-out already stops on it. The issue only counted a close-out that ran while another story's narrative was mid-edit.

The question is what a close-out does with a markdown finding on another story's narrative.

## Decision

A close-out records no markdown finding on the narrative of another open story: a check scoped to a story leaves out every `markdown.*` finding on the narrative of another story that is neither done nor cancelled, as ADR-0122 does with an `item.archive` outside the story.

1. **A check scoped to a story leaves it out.** `flai check --story S-nnnn` drops from its result every finding whose rule starts with `markdown.` and whose path is the narrative, `wip/agents/<id>.md`, of another story the repository holds whose status is neither `done` nor `cancelled`, and takes back the counts it added. That story's own close-out, whose scoped check counts its narrative as inside, finds it and stops on it.
2. **Every other markdown finding outside the story is still recorded.** One on a thread, a task, a conversation, a narrative under `wip/archive/agents`, or the narrative of a done or cancelled story is a note recorded in an issue, as ADR-0085 says. No later close-out checks those as its own, and such a finding shows a gap in flai's own wip lint, as I-0056, I-0070, and I-0072 did.
3. **Without `--story` every finding is reported**, so the main checkout's `flai check` and CI still warn.

This refines ADR-0085, whose rule that every finding outside the story is a note recorded in an issue now leaves out a markdown finding on another open story's narrative, and ADR-0122, whose treatment of `item.archive` it extends to a third kind of finding.

## Consequences

- A close-out no longer records an issue because another story's agent left its narrative failing the markdown lint for a while, and I-0096 can be closed: its cause no longer occurs.
- The closing story's agent no longer sees, at its close-out, a finding on a narrative it does not own. The owning story's close-out still stops on it, and the main checkout's check still warns.
- A markdown finding on another story's narrative that outlives that story, because the story was accepted or cancelled with it, is still recorded, by the next close-out.
- An open issue for a markdown rule recorded by an older flai stays as it is until it is closed.

## Alternatives considered

- **A flai command that writes a narrative's hand-edited sections through the lint guard.** It would keep the finding out of `wip/` at its source, but it changes `tooling.md`, which lets agents edit those sections by hand, and it is a story of its own.
- **Leave out every finding on another story's narrative.** It would also cover `narrative.state`, which S-0323 remediates on its own; this decision keeps to the markdown rules I-0096 recorded.
- **Leave out every markdown finding outside the story.** A finding on a thread, a task, or a closed story's narrative would then reach no issue, though no story's close-out finds it and it shows a gap in flai's wip lint.
