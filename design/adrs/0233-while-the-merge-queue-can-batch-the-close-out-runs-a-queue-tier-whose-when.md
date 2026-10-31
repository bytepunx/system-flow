---
id: ADR-0233
title: "While the merge queue can batch, the close-out runs a queue tier whose when matches the story's claim or changed paths, and the batch runs it again"
status: accepted
date: 2026-10-31
supersedes: []
superseded_by: []
refines: [ADR-0174, ADR-0221]
---

# ADR-0233 While the merge queue can batch, the close-out runs a queue tier whose when matches the story's claim or changed paths, and the batch runs it again

## Context

[ADR-0174](0174-while-the-merge-queue-can-batch-flai-verify-and-the-close-out-leave-the-test.md) decided that while the merge queue can batch, `flai verify` and the close-out leave every tier the manifest marks `queue: true` to the batch. In this repository those are `integration` and `smoke`. In the fortnight from 2026-10-15 to 2026-10-29, 24 of 104 queue verifications failed and 14 stories were sent back, 22 percent of landings (E-0031). Each send-back cost the story's agent a restart, a fix, another close-out, a re-enqueue, and another verification of about ten minutes. I-0221 records a change that broke a caller package's tests and passed its close-out; I-0203 a stale generated file that failed every later batch at `integration`. Both would have been found before review had the close-out run the queue tiers.

[ADR-0221](0221-while-the-merge-queue-can-batch-a-hand-acceptance-is-refused-until-the-queue.md) refuses a hand acceptance until the queue tiers ran at the story's head, and lets `flai verify --all` after a close-out at the same patch run only those tiers.

## Decision

While the merge queue can batch, the close-out runs a selected `queue` tier whose `when` matches the story's claim or a path its branch changes, and the batch runs it again.

- A test tier with `queue: true` may carry `when`, a list of path patterns in the language `paths` uses.
- `flai verify`, the MCP `verify` tool, and the close-out run such a tier when a pattern of its `when` matches one of the story's touches, as its claim holds them (the story's and its open tasks'), or a path its branch changes. They leave it to the batch otherwise, as ADR-0174 says. The report names the pattern and the touch or path that chose it.
- A tier without `when` is left to the batch as before.
- The batch's verification never leaves a tier out, so it runs the tier again as the landing check, on the merge of the batch.
- A `queue` tier the close-out ran at the story's head counts for ADR-0221's hand acceptance, and `flai verify --all` at the same patch reuses it.
- `when` on a tier without `queue: true` is a `manifest.tests` error, as is an empty or invalid pattern.

## Consequences

- A story whose claim reaches what a queue tier tests finds that tier's failures in its own close-out, before review, rather than in a batch that sends it back.
- The close-out takes longer by the queue tiers' time, about twelve minutes in this repository for `integration` and `smoke`. A tier that runs covers what it covers, so `integration` skips `go-test` again. Running the close-out as a `flai serve` job (S-0546) keeps that time off a tool call's timeout.
- The batch repeats the tiers on the merged result, so a conflict between two stories' changes is still caught at landing.

## Alternatives considered

- Run every `queue` tier in the close-out, with no `when`: simplest, but a story that changes only documentation would wait for the whole suite.
- Match `when` against the touches alone: a story whose touches are narrower than its change would skip the tier.
- Drop the batch's run of a tier the close-out ran: the batch's run is the only one on the merge of several stories.
