---
id: T-3347
type: task
nature: remediation
title: Record in an ADR whether a close-out leaves out story.unaccepted outside the story
status: done
parent: S-0781
owner: alex
created: 2026-11-13T22:53:27Z
updated: 2026-11-13T23:10:46Z
transitions:
  - to: ready
    at: 2026-11-13T23:00:11Z
    by: system-flow
  - to: in-progress
    at: 2026-11-13T23:00:11Z
    by: system-flow
  - to: done
    at: 2026-11-13T23:10:46Z
    by: agent-S-0781
stream: S-0781
tags: [flai, check, adr]
touches: [design/adrs/README.md, design/adrs/0299-a-close-out-records-no-story-unaccepted-and-a-check-scoped-to-a-story-leaves.md]
usage:
  source: log
  seconds: 635
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 11
      output: 3453
      cache_read: 575541
      cache_write: 25685
      cost: 0.347
---
# T-3347 Record in an ADR whether a close-out leaves out story.unaccepted outside the story

## Work

Write and accept an ADR with `flai adr new --commit`. It refines ADR-0085 and follows ADR-0137 (`threads.archived`), ADR-0256 (`narrative.missing`), and ADR-0296 (`item.blocked-closed`). The proposed decision: a `story.unaccepted` finding outside the story is left out of a check scoped to a story (`flai check --story`), so a close-out neither notes it nor records it in an issue. `flai check` in the main checkout still reports it.

Its context is I-0371's five instances, close-outs of S-0778, S-0730, S-0777, S-0780, and S-0779 on 2026-11-13. Each found the same four branches, `story/S-0001`, `story/S-0002`, `story/S-0003`, and `story/S-0005`, which S-0777's test run leaked into the main repository (see this story's first task). Neither of the rule's two cases can be cleared by a story branch. A story done but not archived is completed by `flai accept` in the main checkout. A branch left after acceptance is a ref of the main repository, which the operator merges or deletes. Inside the story the finding cannot arise: a story at its close-out is in progress, not done.

The alternative is the root fix alone: stop the leak and remove the four branches, and change nothing in the check. Then a real unaccepted story would still be noted by every close-out until the operator clears it. Decide that instead if it fits better, and say why in the ADR. If it is decided, say so on the story's thread before the code and docs tasks are worked, since they are then cancelled.

Add the ADR's file to this task's and the story's touches once its title is set. It waits for nothing in this story.

## Done when

- The ADR is accepted under `design/adrs/`, indexed in `design/adrs/README.md`, with its topics set.
- It says whether a `story.unaccepted` outside the story is left out, why, and what still reports it.

## Notes
