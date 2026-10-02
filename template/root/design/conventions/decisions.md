---
title: Decisions
updated: 2026-10-02
audience: agent
order: 40
status: active
topics: [all]
roles: [story]
---

# Decisions

What counts as a decision, where each kind is recorded, and what never happens to a recorded decision.

## Rules

- A decision is any choice that someone could reasonably have made differently: a technology, a structure, a rule, a naming, a scope cut, a deferral.
- Record decisions immediately in the narrative's `## Decisions` with a one-line reason.
- Architecture Decision Records (ADRs):
  - in `design/adrs`
  - changes to structure, technology, contracts between parts, or a tool-enforced rule
  - record with `flai adr new "<the decision, as a sentence>"`:
    - calculates the next number from the files present
    - names the file
    - writes the front matter
    - adds the row to the index
    - sets `superseded_by` on any ADR it supersedes (`--supersedes`, `--refines`, `--status accepted`)
    - sets the body on standard input with `--body-stdin`
    - do not copy the template by hand or look the number up
    - One decision per file
- Never edit an accepted ADR except to set `superseded_by`, or its `topics` with `flai adr topics`.
- An ADR that changes the current system also changes the living design in `design/system` or `design/tech`, in the same change, with a link to the ADR.
- Decisions contained to a single story's scope are recorded in `## Notes` and the narrative.
- Stories that change the system behavior must also be captured in the `design/system` or `design/tech` folders.
- Never silently reverse recorded decisions:
  - if a decision looks wrong, propose alternatives in `## Open questions`
  - keep following the decision until the operator changes it
- Refinements that change a default value, an added case, a clarified rule:
  - are living-design edits, not ADRs
  - capture refinements in the narrative
- Record operator decisions immediately, whether made in a thread or in conversation, before doing anything else. A decision that exists only in a transcript does not exist.
- Do not present questions and decisions in conflated ways.

## When in doubt

- If you are unsure whether it needs an ADR, it is at least a living-design edit plus a narrative entry; the ADR can follow if the operator wants it.
- If two recorded decisions conflict, newer wins and older gets `superseded_by`.

<!-- system-flow:end-of-baseline -->

## Project additions
