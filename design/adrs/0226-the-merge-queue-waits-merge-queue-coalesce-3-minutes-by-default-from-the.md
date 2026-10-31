---
id: ADR-0226
title: "The merge queue waits merge_queue.coalesce, 3 minutes by default, from the earliest enqueue among a batch's stories before it builds a batch of fewer than merge_queue.batch_max, and shows the window it waits on"
status: accepted
date: 2026-10-31
supersedes: []
superseded_by: []
refines: [ADR-0173]
topics: [orchestration, flai, cli, dashboard]
---

# ADR-0226 The merge queue waits merge_queue.coalesce, 3 minutes by default, from the earliest enqueue among a batch's stories before it builds a batch of fewer than merge_queue.batch_max, and shows the window it waits on

## Context

`mergequeue.RunBatch` built a batch as soon as one story was enqueued and eligible ([ADR-0173](0173-a-merge-queue-in-flai-serve-with-no-model-runs-behind-the-merge-queue-host.md)). `flai serve`'s queue pass looks every 30 seconds and at once when the work items change, so a story enqueued alone was batched alone. Stories reach review minutes apart, so the next one waited for that batch to finish. From 2026-10-15 to 2026-10-29 the queue built 111 batches of mean size 1.0 to land 49 stories (S-0552's evidence, from `flai stats`). Each batch is one verification of about ten minutes here. A batch of several stories costs one verification while it passes.

## Decision

The merge queue waits `merge_queue.coalesce`, 3 minutes by default, from the earliest enqueue among a batch's stories before it builds a batch of fewer than `merge_queue.batch_max`. It shows the window it waits on.

- **The key.** `merge_queue.coalesce` in `system-flow.yaml`, a Go duration such as `3m` or `90s`. Unset, it is 3 minutes. `0` builds each batch at once, as before. `manifest.Load` refuses a value that is not a duration, and a negative one, naming the key. `flai manifest set` writes it. It is not a front-matter field, so it raises no `flai.minimum`; an older flai ignores it.
- **The window.** It opens when the earliest of the stories the next batch would pick was enqueued, and ends `merge_queue.coalesce` later. `RunBatch` builds no batch while those stories are fewer than `merge_queue.batch_max` and the window has not ended. With `merge_queue.batch_max` stories, the batch starts at once.
- **No new window.** A story keeps the time it was enqueued. A story enqueued while a batch ran, or left enqueued after a batch was dropped or landed, opens no new window. Its window may have ended already, and then the next batch is built at once.
- **When it starts.** The pass looks every 30 seconds, so a batch starts within 30 seconds of its window's end.
- **Shown.** The queue's state carries `window`, `{ends, stories, batch_max}`, while the queue is not held, no batch is under way or waits to land, and a window has not ended. `flai queue --json`, `flai board --json`, the MCP `board`, the host read `board.get`, and the dashboard's merge queue panel read it. `flai queue` prints `batch fills until HH:MM:SSZ: <n> of <batch_max> stories (<stories>)`. The window is computed by the same function `RunBatch` waits on.

This refines ADR-0173: the queue still batches at most `merge_queue.batch_max` enqueued, eligible stories, and now waits for the window before it builds a smaller batch.

## Consequences

- Stories enqueued within the window of one another are verified together, once.
- A story enqueued alone lands up to `merge_queue.coalesce` and 30 seconds later than it did.
- The wait is bounded. The window runs from the earliest enqueue and does not grow with each new one.
- A larger batch fails when any of its stories fails, and its bisection takes more steps. `flai stats` counts those steps ([ADR-0227](0227-the-merge-queue-logs-each-bisection-step-s-verification-as-a-step-verified.md)).
- Whether batches grew is measured with `flai stats`'s mean batch size in the fortnight after the release.

## Alternatives considered

- **A window that restarts at each enqueue.** Under a steady flow of stories no batch would start.
- **A window per story, from its own enqueue.** A story enqueued after the first would hold back a batch whose first story had waited long enough.
- **Always wait for `merge_queue.batch_max` stories.** A story enqueued alone would never land.
- **Lower `merge_queue.batch_max`, or look less often.** Neither puts stories enqueued close together into one batch.
