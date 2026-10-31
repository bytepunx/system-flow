---
id: S-0552
type: story
nature: improvement
title: Acceptances land only through the merge queue, and the queue waits a window to fill a batch, so one verification lands several stories and an acceptance by hand drops no batch
status: in-progress
parent: E-0031
owner: alex
created: 2026-10-30T01:52:31Z
updated: 2026-10-31T08:24:35Z
transitions:
  - to: ready
    at: 2026-10-31T07:23:26Z
    by: alex
  - to: in-progress
    at: 2026-10-31T07:26:13Z
    by: agent-S-0552
tags: [flai]
topics: [orchestration, testing, cli]
touches: [flai/internal/mergequeue, flai/cmd, design/system/flai-cli.md, docs/users/flai.md, docs/operators, design/system/project-manifest.md, flai/internal/manifest/manifest.go, flai/internal/manifest/manifest_test.go, flai/internal/manifest/settings.go, flai/internal/manifest/settings_test.go, docs/users/flai-reference.md, design/adrs/0225-while-the-merge-queue-batches-and-is-not-held-the-operator-s-acceptance-of-a.md, design/adrs/0226-the-merge-queue-waits-merge-queue-coalesce-3-minutes-by-default-from-the.md, design/adrs/0227-the-merge-queue-logs-each-bisection-step-s-verification-as-a-step-verified.md, design/adrs/README.md, design/system/flaiover-dashboard.md, design/system/workflow.md, flai/internal/metrics/queue.go, flai/internal/metrics/queue_test.go, design/system/metrics.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 3464
  estimated: true
  priced_by: estimate
  turns:
    - day: 2026-10-31
      ceremony: 3
      hand_edits: 1
      work: 53
  models:
    - model: claude-haiku-5-5
      input: 94
      output: 398
      cache_read: 5581074
      cache_write: 205101
      cost: 0.4179
    - model: claude-opus-5-5
      input: 530
      output: 3179
      cache_read: 29543743
      cache_write: 1556033
      cost: 10.8151
---
# S-0552 Acceptances land only through the merge queue, and the queue waits a window to fill a batch, so one verification lands several stories and an acceptance by hand drops no batch

## Goal

In the fortnight the merge queue ran 111 batches of mean size 1.0 to land 49 stories: every landing paid a full verification of about ten minutes, and an acceptance by the operator outside the queue dropped the batch under way and, as I-0203 records, landed without the queue tiers. This story lets batches fill: an enqueue waits a window, `merge_queue.coalesce` in the manifest with a default of a few minutes, during which further enqueues join the batch, and a batch under way is joined by the next rather than dropped; and a story the operator accepts on the board or with `flai move` is enqueued rather than landed, so every landing is verified once with its batch. The dashboard's queue page shows the window and the batch filling.

## Acceptance criteria
- [ ] An enqueue starts a window of `merge_queue.coalesce` during which later enqueues join the same batch, with the default, the key, and its check in the manifest schema and reference
- [ ] Accepting a story on the board or with `flai move` enqueues it; landing happens only through the queue, and the operator's acceptance of a story that is in a batch under way does not drop the batch
- [ ] A batch that fails bisects as today, and a story that passed in a failed batch's clean half lands without another verification of its own
- [ ] `flai stats` reports batch size and verifications per landed story, and the fortnight after landing shows a mean batch size above one, recorded in the close-out

## Tasks

The planner drafts them.
- T-2370 The merge queue waits merge_queue.coalesce after the first pending enqueue before it builds a batch, and shows the window
- T-2371 An operator's acceptance on the board or with flai move enqueues the story while the queue batches, and drops no batch
- T-2372 flai stats reports verifications per landed story, counting bisection steps, and a clean half lands without a verification of its own
- T-2373 The board's merge queue panel shows the coalesce window and the batch filling
- T-2374 ADRs and the design, user, and operator docs say acceptances land only through the queue, after its coalesce window, and stats count verifications per landing

## Notes

Evidence: `flai stats` merge-queue figures 2026-10-15 to 2026-10-29. ADR-0174 governs the queue; this story refines it with an ADR where the landing path changes.
