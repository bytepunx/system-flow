---
id: T-1313
type: task
nature: improvement
title: An ADR refining ADR-0085 and S-0318's ADR records the remedy for I-0109, proposed from its instances
status: in-progress
parent: S-0323
owner: alex
created: 2026-10-08T00:26:26Z
updated: 2026-10-08T06:12:29Z
transitions:
  - to: ready
    at: 2026-10-08T06:12:27Z
    by: agent-S-0323
  - to: in-progress
    at: 2026-10-08T06:12:29Z
    by: agent-S-0323
stream: S-0323
tags: [flai]
touches: [design/adrs]
---
# T-1313 An ADR refining ADR-0085 and S-0318's ADR records the remedy for I-0109, proposed from its instances

## Work

Propose the remedy from I-0109's instances and record it in an ADR with `flai adr new`. It waits for no task: it decides what the others build and write.

- All three instances are the same: another story in progress, just started, whose narrative's `## Current state` and `## Next steps` still hold the template's placeholder (S-0265, S-0308, S-0298). `narrative.state` is advisory, so the close-out passed, but `--record-issues` recorded it, though only that story's agent can write its narrative, and that story's own close-out stops on it in its `narrative` step.
- The proposed remedy: a check scoped to a story leaves out a `narrative.state` finding on another story's narrative, with its counts taken back, as ADR-0122 does for `item.archive` and S-0318's ADR does for `markdown.*` on another open story's narrative. Unscoped `flai check` still warns it.
- Say what stays: every other finding on another story's narrative (`narrative.front-matter`, `narrative.stream`, `narrative.section`, `narrative.updated`) is still an outside note and still recorded, because it points at a hand edit or a flai defect rather than a narrative not yet written.
- Link ADR-0085, ADR-0115, ADR-0122, and S-0318's ADR, and refine the last.
- Record the alternative the orchestrator raised on TH-0346, leaving out every finding on another open story's narrative, and why it is not taken.

## Done when

- The ADR is accepted, names I-0109 and its instances, the rule left out, what is still recorded, and the alternative rejected.
- `flai check --strict` passes.

## Notes
