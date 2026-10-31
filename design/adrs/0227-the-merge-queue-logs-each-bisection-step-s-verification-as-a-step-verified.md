---
id: ADR-0227
title: "The merge queue logs each bisection step's verification as a step_verified event, and flai stats reports per week the queue's verifications and its verifications per story landed"
status: accepted
date: 2026-10-31
supersedes: []
superseded_by: []
refines: [ADR-0177]
topics: [cli, dashboard]
---

# ADR-0227 The merge queue logs each bisection step's verification as a step_verified event, and flai stats reports per week the queue's verifications and its verifications per story landed

## Context

The merge queue appends each event it acts on to `.flai-cache/merge-queue/events.jsonl`, and `flai stats` reports per week, from it, the stories landed and sent back, the batches with their size, and the bisections ([ADR-0177](0177-flai-stats-reports-per-week-the-stories-the-merge-queue-landed-or-sent-back.md)). The log has one `verified` line per batch. A bisection verifies main with a part of the batch at each step, and the log recorded none of those steps. So the log counted fewer verifications than the queue ran.

S-0552 lets batches fill ([ADR-0226](0226-the-merge-queue-waits-merge-queue-coalesce-3-minutes-by-default-from-the.md)) and lands the operator's acceptances through the queue ([ADR-0225](0225-while-the-merge-queue-batches-and-is-not-held-the-operator-s-acceptance-of-a.md)), so that one verification lands several stories. Whether it does needs a figure: verifications per story landed, bisection steps among them. `design/system/metrics.md`, the contract between `flai stats` and the dashboard, changes only with an ADR.

## Decision

The merge queue logs each bisection step's verification as a `step_verified` event, and `flai stats` reports per week the queue's verifications and its verifications per story landed.

- **The event.** `mergequeue.Update` appends a `step_verified` line once a bisection step has a result. The line has `batch`, `stories` (those the step verified, the batch's first ones), `passed`, and `step`, the step's number from 1. Its `at` is the result's finish time. A step a stopped run left without a result is verified again in place and logged once.
- **A kind of its own.** `verified` stays one line per batch, so ADR-0177's figures keep their meaning.
- **The figures.** Under `queue.weeks[]`, `verifications` counts the `verified` and `step_verified` lines, and `verifications_per_landed` is `verifications` over `landed`, 0 when nothing landed. `flai stats` prints them on the batches line as `verifications <n> · <x.x> per story landed`, the second only when a story landed. `metrics.md` § Merge queue defines them.
- **The clean half.** A story in a failed batch's clean half lands through `flai accept --batch` on the bisection's verification ([ADR-0176](0176-the-merge-queue-lands-each-story-of-a-passing-batch-through-flai-accept-batch.md)), with no verification of its own, and adds none.

This refines ADR-0177: the log gains `step_verified`, and `queue.weeks[]` gains `verifications` and `verifications_per_landed`.

## Consequences

- One figure says what a landing costs in verifications, and whether larger batches lower it.
- A bisection's steps count, so a failed batch costs what it ran.
- Weeks before the release that writes `step_verified` count each batch's verification only.
- The figure is per week: a batch verified late in one week and landed in the next counts its verification in the first and its stories in the second.
- The log grows by a line per bisection step.

## Alternatives considered

- **Log each step as `verified`.** `verified` would no longer be one a batch, and a reader of the log could not tell a batch from a step.
- **Count steps from `state.json`.** Each write replaces it, so a week gone cannot be counted, the reason ADR-0177 gives.
- **Count tier runs.** A verification stops at its first failing tier. The cost S-0552 lowers is the verification, a batch's or a step's.
