---
id: T-2261
type: task
nature: remediation
title: An ADR refining ADR-0174 records that a hand acceptance while the queue can batch is refused until the queue tiers ran at the story's head
status: done
parent: S-0529
owner: alex
created: 2026-10-28T18:16:33Z
updated: 2026-10-30T00:23:21Z
transitions:
  - to: ready
    at: 2026-10-30T00:18:02Z
    by: agent-S-0529
  - to: in-progress
    at: 2026-10-30T00:18:03Z
    by: agent-S-0529
  - to: done
    at: 2026-10-30T00:23:21Z
    by: agent-S-0529
stream: S-0529
tags: [cli]
touches: [design/adrs/README.md, design/adrs/0221-while-the-merge-queue-can-batch-a-hand-acceptance-is-refused-until-the-queue.md]
usage:
  source: log
  seconds: 318
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 6
      output: 2010
      cache_read: 421131
      cache_write: 17712
      cost: 0.2298
---
# T-2261 An ADR refining ADR-0174 records that a hand acceptance while the queue can batch is refused until the queue tiers ran at the story's head

## Work

Write a new ADR with `flai adr new`, refining ADR-0174 (and ADR-0173), from I-0203's instance: S-0491 left `integration` to the queue, the operator accepted it by hand, main kept a stale `docs/users/flai-reference.md`, and the next batch sent S-0490 back for it.

Decision, as the plan's thread settles it (the planner's proposal until then):

- While the queue can batch (ADR-0174's three conditions), `flai accept` by hand, the operator's from the CLI, `flai move <S> done`, the board, or the dashboard, is refused for a story whose last `flai verify` report is not at its branch's head, did not pass, or left a tier marked `queue` that its paths select as `queued`. The refusal names the tiers and `flai verify <S-nnnn> --all`. It is a blocker in `--dry-run` and the dashboard's preview.
- One flag lands it anyway, and the acceptance records that the queue tiers did not run.
- `flai verify --all` reuses the tiers a run that left the queue tiers passed at the same patch, so it runs only the queue tiers.
- The orchestrator's acceptance and the queue's landing are unchanged.

Alternatives to record: running the queue tiers inside `flai accept` (about twelve minutes here, past the dashboard's ten-minute timeout on `accept.run` and `item.move`); enqueuing the operator's acceptance; a batch that checks main alone before it blames a story (a second defence, proposed apart on the plan's thread).

It waits for nothing: it records the decision the other tasks build.

## Done when

- The ADR is accepted, refines ADR-0174, and its consequences say ADR-0174's "`flai accept` by hand ... runs no test tier" no longer holds while the queue can batch.
- `design/adrs/README.md` lists it.
- `flai check --strict` passes.

## Notes
