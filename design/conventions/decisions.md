---
title: Decisions
updated: 2026-09-15
audience: agent
order: 40
status: active
---

# Decisions

What counts as a decision, where each kind is recorded, and what never happens to a recorded decision.

## Rules

- A decision is any choice that a later reader could reasonably have made differently: a technology, a structure, a rule, a naming, a scope cut, a deferral.
- Record every decision the moment it is made, in the narrative's `## Decisions` with a one-line reason. Do not wait for the end of the story.
- Architectural decisions get an ADR in `design/adrs`: anything that changes structure, technology, contracts between parts, or a rule the tooling enforces. Copy `0000-template.md`, take the next number, one decision per file.
- A decision that changes how the system currently is also changes the living design in `design/system` or `design/tech`, in the same change, with a link to the ADR.
- A decision that only affects one story is recorded in that story's `## Notes` and the narrative. It does not need an ADR but if it changes the system behavior, must also be captured in the `design/system` or `design/tech` folders.
- Never edit an accepted ADR except to set `superseded_by`. To change a decision, write a new ADR that supersedes the old one in whole or in part, then update the living design.
- Never reverse a recorded decision silently. If the decision looks wrong, say so in `## Open questions` with a proposed alternative, and keep following the decision until the operator or a new ADR changes it.
- Refinements that do not reverse a decision (a default value, an added case, a clarified rule) are living-design edits, not ADRs. Say in the narrative that it is a refinement.
- When the operator decides something in conversation, write it down before doing anything else. A decision that exists only in the transcript does not exist.
- Do not present a question as a decision or a decision as a question. Decide what is yours to decide; ask what is the operator's.

## When in doubt

- If you are unsure whether it needs an ADR, it is at least a living-design edit plus a narrative entry; the ADR can follow if the operator wants it.
- If two recorded decisions conflict, the newer one wins and the older gets `superseded_by`.

<!-- system-flow:end-of-baseline -->

## Project additions
- ADR numbering continues from the index in `design/adrs/README.md`; next is 0015.
- Changes to conventions land in `template/root/design/conventions/` first and are copied here above the marker; project rules go below the marker in `design/conventions/` only.
- Metrics definitions in `design/system/metrics.md` are the contract between flai and flaiover; change them only with an ADR.
