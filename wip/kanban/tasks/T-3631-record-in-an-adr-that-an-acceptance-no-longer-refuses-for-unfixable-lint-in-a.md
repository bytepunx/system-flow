---
id: T-3631
type: task
nature: improvement
title: Record in an ADR that an acceptance no longer refuses for unfixable lint in a shared-set file the story does not change, and that digest absorb and compact lint what they rewrite
status: backlog
parent: S-0853
owner: alex
created: 2026-11-27T23:58:28Z
updated: 2026-11-27T23:58:28Z
transitions: []
stream: S-0853
tags: [flai]
touches: [design/adrs/README.md]
---
# T-3631 Record in an ADR that an acceptance no longer refuses for unfixable lint in a shared-set file the story does not change, and that digest absorb and compact lint what they rewrite

## Work

I-0411's goal asks for a fix proposed from its instances before it is built. Write it as a new ADR with `flai adr new`, refining ADR-0140, ADR-0234, and ADR-0318, and regenerate `design/adrs/README.md`. The story's `### Planning` gives the proposal:

1. Every acceptance's preflight lint, by hand, by the orchestrator, from the dashboard, and as the merge queue's landing, refuses only for a file the story's change, `main...story/S-nnnn`, touches. A shared-set file the story does not touch whose lint the fix leaves failing is left out of the acceptance's commits and stays uncommitted, as ADR-0248 leaves "the rest". The acceptance names it in its output, at warn in its log, and in the preview, under a field of its own.
2. The merge queue's preflight holds the queue on such a file only when an enqueued story changes it.
3. `flai digest absorb` and `compact` lint what they rewrite, as `new` and `edit` do, and refuse a write that brings a finding.

Say in Consequences that ADR-0140's "Uncommitted files in the main checkout that fail the check block every acceptance" no longer holds, and that the left-out file reaches main with the next commit that includes it, such as the operator's or a later acceptance once it is fixed. Cite I-0411's instances: 20 refusals over 16 stories, and the five stories held 76 h to 82 h from 2026-11-18. If the code the later tasks read shows the proposal cannot hold, ask on a thread on S-0853 before going on.

It waits for no task: it is the first layer, and the three code tasks build what it decides.

## Done when

- A new ADR under `design/adrs/` records the three rules above, with Context, Decision, Consequences, and Alternatives considered, and `refines` names ADR-0140, ADR-0234, and ADR-0318.
- `design/adrs/README.md` lists it.
- `flai check --strict` and the markdown lint pass on the ADR.

## Notes
