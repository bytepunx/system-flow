---
id: T-3431
type: task
nature: feature
title: The review page follows a Rebase and resolve job in a collapsible activity stream to its end and shows whether the story stayed in review or was sent back and why
status: done
parent: S-0804
owner: alex
created: 2026-11-16T23:19:49Z
updated: 2026-11-17T01:06:16Z
transitions:
  - to: ready
    at: 2026-11-16T23:20:38Z
    by: agent-S-0804
  - to: in-progress
    at: 2026-11-17T00:52:30Z
    by: agent-S-0804
  - to: done
    at: 2026-11-17T01:06:16Z
    by: agent-S-0804
stream: S-0804
tags: []
touches: [flaiover/src/lib/components/Review.svelte, flaiover/src/lib/components/Review.svelte.test.ts, flaiover/src/lib/review.ts, flaiover/src/lib/review.test.ts, "flaiover/src/routes/api/items/[id]/resolve/+server.ts", "flaiover/src/routes/api/items/[id]/resolve/resolve.test.ts", design/system/flaiover-dashboard.md]
after: [T-3430]
usage:
  source: log
  seconds: 826
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 112
      output: 47466
      cache_read: 8563397
      cache_write: 321778
      cost: 4.4577
---
# T-3431 The review page follows a Rebase and resolve job in a collapsible activity stream to its end and shows whether the story stayed in review or was sent back and why

## Work

Waits for T-3430, whose `accept.resolve` answer and `resolve.status` method this reads.

1. `flaiover/src/routes/api/items/[id]/resolve/+server.ts`: POST keeps calling `accept.resolve`, now answering the job (`{id, job, started, running, log, reason}`) or `{refused}`; add GET, which calls `resolve.status` with `{id}`.
2. `flaiover/src/lib/review.ts`: types and pure functions for the job: its kind, step, tier, steps, elapsed time, and its end (stayed in review, sent back with its thread, conflicts, and reason, or interrupted), with tests in `review.test.ts`.
3. `Review.svelte`: after Rebase and resolve answers a job, and on load when `resolve.status` says a resolve job of the story runs, show it in a collapsible activity stream like an agent's activity: its current step and tier, the elapsed time ticking, and its output as it is written, through `LogFeed` reading the job's log with `log.read` (GET `/api/log`, as `LogViewer.svelte` reads one). Poll GET resolve until the job ends, without a reload; then show whether the story stayed in review or was sent back and why, with the thread link and the conflicting paths or the step, and load the item and the preview again. The old sent-back block stays for a refusal-free answer that has no job, which no longer comes. Build from the existing classes and tokens only.
4. `design/system/flaiover-dashboard.md`: the Review view's Rebase and resolve paragraph, the resolve route's row (GET and POST), and `resolve.status` in the host-method list; the design-routes test reads it.
5. Tests in `Review.svelte.test.ts` and `resolve.test.ts`: a running job's stream, its end in each outcome, and the GET route.

## Done when

- `flai test` passes on the changed paths: flaiover-lint, vitest, svelte-check, design-routes.

## Notes
