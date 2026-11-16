---
id: T-3394
type: task
nature: remediation
title: "Write ADR-0308 refining ADR-0291: while the merge queue batches and is not held, the preview lists verify_behind_main only for an acceptance that lands by hand"
status: done
parent: S-0796
owner: alex
created: 2026-11-16T00:50:46Z
updated: 2026-11-16T01:01:46Z
transitions:
  - to: ready
    at: 2026-11-16T00:53:16Z
    by: agent-S-0796
  - to: in-progress
    at: 2026-11-16T00:53:16Z
    by: agent-S-0796
  - to: done
    at: 2026-11-16T01:01:46Z
    by: agent-S-0796
stream: S-0796
tags: [adr, acceptance, merge-queue]
touches: [design/adrs/0308-while-the-merge-queue-batches-and-is-not-held-the-acceptance-preview-lists.md, design/adrs/README.md]
usage:
  source: log
  seconds: 510
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 9
      output: 1546
      cache_read: 854188
      cache_write: 16485
      cost: 0.3013
---
# T-3394 Write ADR-0308 refining ADR-0291: while the merge queue batches and is not held, the preview lists verify_behind_main only for an acceptance that lands by hand

## Work

Write the ADR the story's goal asks for with `flai adr new`, titled "While the merge queue batches and is not held, the acceptance preview lists verify_behind_main only for an acceptance that lands by hand, and reports the stale tiers the batch will run", with `refines: [ADR-0291]`, and accept it.

- Context: I-0381. ADR-0291 §1 says the gate meets only what lands by hand, but `verifyBehindMainGate` in `flai/internal/preview/accept.go` skips only for the queue's landing (`opts.Batch`). So `flai accept --dry-run` and the dashboard's `accept.preview` list the blocker for a story the acceptance would enqueue, and the review page offers Rebase and resolve in place of Accept (S-0742, S-0747, S-0739 on 2026-11-15).
- Decision: the operator's preview leaves `verify_behind_main` out of its blockers while the queue batches, is not held, and the acceptance would enqueue, as `queuesInstead` in `flai/cmd/accept.go` decides. It still reports the stale tiers, with a sentence saying the batch runs them. A held queue, `--without-queue-tiers`, `--without-kluster`, `--keep-removal`, and the orchestrator's acceptance keep the gate, since each lands by hand.
- Consequences: the review page offers Accept, which enqueues; the batch runs every tier over the main it lands on.

The number is the next free one when the task runs; if another story took 0308, use the next and let `flai task done` correct the touches. This task waits for none: the decision is in the story already, and the code tasks can be written beside it.

## Done when

- The ADR is accepted under `design/adrs/`, refines ADR-0291, and names the conditions that keep the gate.
- `design/adrs/README.md` lists it, and `flai check --strict` and the markdown lint pass on both files.

## Notes
