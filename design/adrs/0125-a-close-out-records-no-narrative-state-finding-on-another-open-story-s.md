---
id: ADR-0125
title: "A close-out records no narrative.state finding on another open story's narrative, and a check scoped to a story leaves every such finding out"
status: proposed
date: 2026-10-08
supersedes: []
superseded_by: []
refines: [ADR-0085, ADR-0123]
---

# ADR-0125 A close-out records no narrative.state finding on another open story's narrative, and a check scoped to a story leaves every such finding out

## Context

[ADR-0085](0085-a-close-out-s-flai-check-reports-findings-outside-the-story-as-notes-and.md) scoped a close-out's `flai check` to its story and recorded each rule's findings outside the story in one issue per rule. [ADR-0122](0122-a-close-out-records-no-item-archive-and-a-check-scoped-to-a-story-leaves-out.md) and [ADR-0123](0123-a-close-out-records-no-markdown-finding-on-another-open-story-s-narrative-and-a.md) took out of that recording two kinds of finding the closing story's agent cannot fix: an `item.archive` outside the story, and a markdown finding on another open story's narrative. ADR-0123 left `narrative.state` to S-0323.

I-0109 and its duplicate I-0111, ``flai check finds `narrative.state` outside the story at close-out``, have eight instances between them, from five closing stories. Every one has the same cause.

1. **Another story just started.** Each finding is on the narrative of a story in progress, such as S-0265, S-0308, S-0317, or S-0324, whose `## Current state` and `## Next steps` still hold the template's placeholder. `flai stream open` writes them so, and they stay so until that story's agent first runs `flai stream state`.
2. **Only that story's agent can write them, and its close-out stops on them.** The `narrative` step of `flai verify` stops a story whose own narrative is unwritten. The closing story's agent could not write the other story's narrative, yet its close-out marked the findings outside and recorded them in I-0109 or I-0111.

The rule is advisory, so `--strict` never failed on it; the cost is the issue, which only counted close-outs that ran while another story had just started.

## Decision

A close-out records no `narrative.state` finding on the narrative of another open story: a check scoped to a story leaves out every `narrative.state` finding on the narrative of another story that is neither done nor cancelled, as ADR-0123 does with a markdown finding there.

1. **A check scoped to a story leaves it out.** `flai check --story S-nnnn` drops from its result every finding of rule `narrative.state` whose path is the narrative, `wip/agents/<id>.md`, of another story the repository holds whose status is neither `done` nor `cancelled`, and takes back the counts it added. That story's own close-out stops on it in the `narrative` step, and its own scoped check counts its narrative as inside.
2. **Every other finding on another story's narrative is still a note recorded in an issue**, as ADR-0085 says, but the markdown findings ADR-0123 leaves out. `narrative.state` fires only on a story in progress or in review, so this leaves out every one outside the story.
3. **Without `--story` every finding is reported**, so the main checkout's `flai check` still warns.

This refines ADR-0085, whose rule that every finding outside the story is a note recorded in an issue now also leaves out a `narrative.state` finding on another open story's narrative, and ADR-0123, whose treatment of markdown findings there it extends to a second rule.

## Consequences

- A close-out no longer records an issue because another story had just started when it ran, and I-0109 and I-0111 can be closed: their cause no longer occurs.
- The closing story's agent no longer sees, at its close-out, a note on a narrative it does not own. The owning story's close-out still stops on it, and the main checkout's check still warns.
- An open issue for `narrative.state` recorded by an older flai stays as it is until it is closed.

## Alternatives considered

- **Leave out every finding on another open story's narrative.** It would also hide `narrative.section` and `narrative.front-matter`, which flai writes and which show a defect in flai rather than a story's work in progress; no instance asks for it.
- **Have `flai stream open` write a first Current state and Next steps.** The placeholder is what tells the story's own close-out that its agent never wrote them; filling it in would defeat the `narrative` step.
- **Leave `narrative.state` unrecorded everywhere outside the story, open or not.** It fires only on a story in progress or in review, so this is the same set; naming the open stories keeps one rule for both ADR-0123's findings and this one.
