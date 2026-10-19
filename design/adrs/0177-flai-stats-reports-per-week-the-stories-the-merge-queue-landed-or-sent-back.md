---
id: ADR-0177
title: "flai stats reports per week the stories the merge queue landed or sent back with their queue wait, the batches it built with their size, and its bisections, from an event log the queue appends under .flai-cache"
status: accepted
date: 2026-10-19
supersedes: []
superseded_by: []
refines: [ADR-0081, ADR-0138]
topics: [cli, dashboard]
---

# ADR-0177 flai stats reports per week the stories the merge queue landed or sent back with their queue wait, the batches it built with their size, and its bisections, from an event log the queue appends under .flai-cache

## Context

The merge queue in `flai serve` batches enqueued stories, verifies each batch once, lands a green one, and bisects a red one to the story it sends back ([ADR-0173](0173-a-merge-queue-in-flai-serve-with-no-model-runs-behind-the-merge-queue-host.md), [ADR-0176](0176-the-merge-queue-lands-each-story-of-a-passing-batch-through-flai-accept-batch.md)). ADR-0173 gives S-0404 the queue wait per story, the batch size, and the bisections per week. [ADR-0081](0081-flai-stats-reports-forecast-error-cost-of-delay-waiting-holds-touches-drift-and.md) has `flai stats` report what the dashboard charts, so that the dashboard works out no figures of its own. [ADR-0138](0138-flai-stats-reports-per-week-the-conversations-between-agents-and-how-they-ended.md) counted conflicts per week from an append-only log under `.flai-cache`. `design/system/metrics.md`, the contract between `flai stats` and the dashboard, changes only with an ADR.

The queue's state, `.flai-cache/merge-queue/state.json`, holds what is enqueued now, the batch under verification, and each story's last outcome. Each write replaces it. So it says nothing of a week gone: when a story entered, how long it waited, how many stories a batch held, or how often a batch was bisected.

## Decision

**The merge queue appends each event it acts on to `.flai-cache/merge-queue/events.jsonl`, and `flai stats` reports from it, under `queue` and per ISO week of the window, the stories the queue landed and sent back with their queue wait, the batches it built with their size, and its bisections, with each item's queue wait in `items[]`.** It reports them as `metrics.md` § Merge queue defines them, refining ADR-0081 and ADR-0138.

- **The log.** One JSON line per event, beside `state.json` in the main checkout. Its events are `enqueued`, `batch_built`, `verified`, `bisected`, `landed`, `sent_back`, and `dequeued`. Each line has `at` and `event`, and as the event needs, `batch`, `stories`, `by`, `passed`, `left_out`, and `reason`.
- **Who writes it.** `mergequeue.Update`, the state's only writer. It compares the state before and after each write and appends a line per event the write made. So flai serve, `flai queue add` and `remove`, the host API, and `flai accept`'s dequeue are all covered, and no caller can forget a line.
- **Kept as the conflict log is.** Append-only and never trimmed, so a window of any length can be counted. A last line a write was cut off in is skipped when read and dropped by the next append, so it never joins the next line. A failed write is logged as a warning and the write to the queue goes on: the log measures, it does not gate.
- **Queue wait.** For a landed or sent-back story, from its latest `enqueued` line before the outcome to the outcome. Counted in the week of the outcome, and summed per item over the waits that end in the window.
- **Batches.** The `batch_built` lines, once a batch has left the building and rebasing steps, so its stories are the ones it verifies. A batch whose every story was left out writes none.
- **Bisections.** The `bisected` lines, one per failed batch whose bisection began.
- **Absent without a log.** With no log, or one that cannot be read, `queue` is absent; an unreadable log is warned of on stderr.
- **Not counted.** Hold and resume of the queue, and a dropped batch. A dropped batch's stories stay enqueued, so their wait goes on, and the next batch counts again.

## Consequences

- The dashboard's Queue chart reads `flai stats --json` alone, and matches it.
- Figures start with the release that writes the log. Weeks before it show nothing, which is no evidence either way. A story enqueued before that release and landed after it has no wait.
- The log grows by a line per event and is never trimmed. A write that fails loses that event from the stats; the warning says so.
- A dropped batch and its rebuild count as two batches, since each was built. The waits of its stories run on through both.
- A story enqueued again after it was sent back waits from the new enqueueing, not the first.

## Alternatives considered

- **Read the figures from `state.json`.** It keeps only the present and each story's last outcome, so a week gone cannot be counted.
- **Have each caller append its own line.** flai serve, two `flai queue` commands, the host API, and `flai accept` all write the state. A caller that forgets a line undercounts silently. Diffing in `Update` covers every writer, present and future.
- **Count from the board's transitions.** A landing moves a story to done, but so does an acceptance by hand, and nothing on the board records enqueueing, batches, or bisections.
- **Count a dropped batch apart, or not at all.** It was built and its stories rebased, which is the cost the batch size measures. Counting it as built keeps the count true to the work done.
- **Wait from the first enqueueing.** A story sent back goes to in-progress and is worked again. Its time there is not queue time, so the wait starts at its latest enqueueing.
