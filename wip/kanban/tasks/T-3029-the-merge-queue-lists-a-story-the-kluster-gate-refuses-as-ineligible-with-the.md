---
id: T-3029
type: task
nature: remediation
title: The merge queue lists a story the kluster gate refuses as ineligible, with the reason kluster, instead of batching it
status: done
parent: S-0694
owner: alex
created: 2026-11-10T10:33:28Z
updated: 2026-11-11T14:23:58Z
transitions:
  - to: ready
    at: 2026-11-11T09:13:59Z
    by: agent-S-0694
  - to: in-progress
    at: 2026-11-11T09:13:59Z
    by: agent-S-0694
  - to: done
    at: 2026-11-11T14:23:58Z
    by: agent-S-0694
stream: S-0694
tags: [flai, kluster, merge-queue]
touches: [flai/internal/mergequeue/eligible.go, flai/internal/mergequeue/eligible_test.go, flai/internal/mergequeue/batch.go, flai/internal/mergequeue/batch_test.go, flai/internal/serve/mergequeue.go, flai/cmd/queue.go]
after: [T-3024]
usage:
  source: log
  seconds: 1796
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 98
      output: 471
      cache_read: 3743959
      cache_write: 144881
      cost: 1.4396
---
# T-3029 The merge queue lists a story the kluster gate refuses as ineligible, with the reason kluster, instead of batching it

## Work

The merge queue's acceptance is refused under ADR-0254 §5a, as the operator's is. Without this task the queue would batch a story whose namespace failed, verify the batch, and then drop it when `flai accept --batch` refuses the story (`Land()` in `land.go`).

- In `ComputeAt()` in `flai/internal/mergequeue/eligible.go`, give a story in review the reason `kluster` while the `kluster` host action is on and T-3024's `kluster.AcceptGate` refuses it at the story branch's head. The reason names the last job and its outcome, as the other reasons do, and the story stays out of every batch.
- A story the gate passes, or any story while the action is off, is judged as it is today.

This task waits for T-3024, whose gate it calls. It runs beside the send-back and the accept blocker, since no file overlaps.

## Done when

- `eligible_test.go` shows a story refused by the gate listed with the reason `kluster` and left out of the batch. It shows the same story eligible once the gate passes, and eligible with the action off.
- `flai test` passes on the two files.

## Notes
