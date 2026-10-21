---
id: T-1951
type: task
nature: improvement
title: Document the merge queue's kept batch logs and the findings its send-back thread quotes, for builders, users, and operators
status: done
parent: S-0457
owner: alex
created: 2026-10-21T06:10:26Z
updated: 2026-10-21T06:56:44Z
transitions:
  - to: ready
    at: 2026-10-21T06:54:58Z
    by: agent-S-0457
  - to: in-progress
    at: 2026-10-21T06:54:59Z
    by: agent-S-0457
  - to: done
    at: 2026-10-21T06:56:44Z
    by: agent-S-0457
stream: S-0457
tags: [flai, docs]
touches: [design/system/flai-cli.md, docs/users/flai.md, docs/operators/index.md, docs/operators/runbooks/backup.md]
after: [T-1950]
usage:
  source: log
  seconds: 105
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 29
      output: 9735
      cache_read: 1268493
      cache_write: 63954
      cost: 0.8251
---
# T-1951 Document the merge queue's kept batch logs and the findings its send-back thread quotes, for builders, users, and operators

## Work

- In `design/system/flai-cli.md`, under `### The merge queue pass`, say in the bullets on the state and on send-back that a result keeps the failing tier's findings in full, within their cap, beside its short summary, that each batch and bisection step keeps its tiers' output under `.flai-cache/merge-queue/batches/<batch id>/`, the newest 20 kept, and that the send-back thread quotes the findings and names the folder.
- In `docs/users/flai.md`, under `### The merge queue` › `#### Landing and sending back`, say what the send-back thread quotes and where the batch's whole output is, beside the paragraph naming `.flai-cache/merge-queue/events.jsonl`.
- In `docs/operators/index.md`, under the merge-queue host action's `Where to look` and `What it does with a red batch`, name the batch log folder.
- In `docs/operators/runbooks/backup.md`, under `## What exists, and what is worth keeping`, add a row for `.flai-cache/merge-queue/batches/`: not worth keeping, since it is diagnostic and pruned.

It waits for T-1950, so that it documents the thread's final text, the cap, and the folder as they were built.

## Done when

- The four documents state the kept logs and the thread's quote as built, with their `updated` dates bumped where they carry one.
- The markdown lint passes on them.

## Notes

S-0453, in review, also changes `design/system/flai-cli.md` and `docs/users/flai.md`, in other sections: sync before editing them.
