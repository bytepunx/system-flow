---
id: T-1638
type: task
nature: feature
title: "An ADR records the merge queue: what runs it, its host action and permission, eligibility, the batch's upper size, and that a story undoing main's change is never batched"
status: in-progress
parent: S-0398
owner: alex
created: 2026-10-16T05:00:37Z
updated: 2026-10-19T14:35:43Z
transitions:
  - to: ready
    at: 2026-10-19T14:35:43Z
    by: agent-S-0398
  - to: in-progress
    at: 2026-10-19T14:35:43Z
    by: agent-S-0398
stream: S-0398
tags: [cli]
touches: [design/adrs]
usage:
  source: log
  seconds: 3201
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 108
      output: 37345
      cache_read: 4873857
      cache_write: 243484
      cost: 3.0709
---
# T-1638 An ADR records the merge queue: what runs it, its host action and permission, eligibility, the batch's upper size, and that a story undoing main's change is never batched

## Work

- Write the ADR with `flai adr new`, from the answers on TH-0465 to questions 2 to 4: what runs the queue, the `merge-queue` host action and `accept_reviews`, enqueueing, eligibility, `merge_queue.batch_max`, and ADR-0155's stories kept out.
- Record that the queue runs under `acceptance.mode: local` only. Under pull-request acceptance the forge's required checks gate each story (ADR-0169). A batch's rebase would leave the forge state behind the branch's head, so every landing would stop on `forge_stale`. Link ADR-0165 and ADR-0169.
- Name the state file and its package, which the later stories of E-0022 write to.
- Name the stories of E-0022 that build each part.
- First in its layer: it waits for no task.

## Done when

The ADR is accepted, indexed in `design/adrs/README.md`, and `flai check --strict` passes.

## Notes
