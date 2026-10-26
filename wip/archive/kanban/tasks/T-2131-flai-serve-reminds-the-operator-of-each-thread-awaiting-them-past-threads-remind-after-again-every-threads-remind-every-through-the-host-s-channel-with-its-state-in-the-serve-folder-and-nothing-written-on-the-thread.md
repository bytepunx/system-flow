---
id: T-2131
type: task
nature: improvement
title: flai serve reminds the operator of each thread awaiting them past threads.remind_after, again every threads.remind_every, through the host's channel, with its state in the serve folder and nothing written on the thread
status: done
parent: S-0498
owner: alex
created: 2026-10-26T04:15:01Z
updated: 2026-10-26T05:05:36Z
transitions:
  - to: ready
    at: 2026-10-26T04:57:08Z
    by: agent-S-0498
  - to: in-progress
    at: 2026-10-26T04:57:09Z
    by: agent-S-0498
  - to: done
    at: 2026-10-26T05:05:36Z
    by: agent-S-0498
stream: S-0498
tags: [flai]
touches: [flai/internal/serve/remind.go, flai/internal/serve/remind_test.go, flai/internal/serve/serve.go, flai/internal/serve/serve_test.go, flai/cmd/serve.go]
after: [T-2128, T-2130]
usage:
  source: log
  seconds: 507
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 108
      output: 36191
      cache_read: 6185216
      cache_write: 184535
      cost: 3.0249
---
# T-2131 flai serve reminds the operator of each thread awaiting them past threads.remind_after, again every threads.remind_every, through the host's channel, with its state in the serve folder and nothing written on the thread

## Work

A `reminder` in `flai/internal/serve/remind.go`, built as `upgrader` in `upgrade.go` is, with a `look(ctx)` the per-project loop in `serve.go` calls on each minute's tick and on the watcher's change, wired where the analysis, curate, and upgrade schedulers are (`serve.go` around the per-project loop), with its tests:

- **Due.** Read the project's manifest each look (the operator may change it while `flai serve` runs). With `threads.remind_after` unset, do nothing. Otherwise list the threads with `threads.List`, keep those `inbox.Summarize` would mark awaiting the operator (a pending recommendation, or a last entry by neither the project's owner nor its story's agent; reuse the inbox package's rule rather than copying it) whose anchor is a work item that is not done or cancelled or a document, and take the time of the thread's last entry (its creation when it has none). A thread is due when now minus that time is at least `remind_after` and it was not reminded since that entry, or at least `remind_every` has passed since its last reminder.
- **Send.** Build the `Reminder` of T-2130 (the dashboard link is `/items/<item>?thread=<TH>` for a thread on an item, `/threads?thread=<TH>` otherwise, on the dashboard URL the manifest names when it does) and call `Send`. With no channel set, log the reminder at info and, once per project per run, warn that no channel is set.
- **State.** `reminders.json` in the project's serve folder, beside `upgrade-asks.json`: per thread, the time of the entry the reminder was for and when it was last sent. A restart repeats nothing; a thread resolved or no longer awaiting the operator is dropped from the file at the next look. Nothing is written in the repository: no thread entry, no narrative line, so no wait ends (ADR-0091) and no status changes (ADR-0090).
- **Log.** Each reminder sent, and each send that failed, to the serve log with the thread, the item, and the hours awaited.

Tests in `remind_test.go` with a fake clock and a recording channel: a thread with a pending recommendation older than `remind_after` is reminded once, then again after `remind_every`, not before; an answer or a confirmation on it stops the reminders and a new question starts the count over; a thread awaiting another agent, on a done item, or younger than the threshold is not reminded; an unset `remind_after` sends nothing; the thread file is byte-identical before and after a reminder; and a reminder's state survives a new `reminder` built over the same folder. This is the test that reproduces I-0151: a recommendation pending for 25 hours, as TH-0436's was, reaches the operator at `remind_after` instead of waiting unseen.

Waits for T-2128 (the settings it reads) and T-2130 (the channel it sends to).

## Done when

- With `threads.remind_after: 4h` in a project's manifest and a channel set, a thread that awaits the operator longer than 4 h is reminded at the next look, and again every `remind_every`, until they answer or confirm; the thread file is unchanged by it.
- The eight tests above pass under `flai test flai/internal/serve/remind.go flai/internal/serve/serve.go`, and the existing serve tests are unchanged.

## Notes

Layer 4, with T-2132 (docs and design), which shares no path with it. T-2133 (closing the issue) waits for both.
