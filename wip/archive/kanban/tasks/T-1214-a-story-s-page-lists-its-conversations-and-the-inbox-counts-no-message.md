---
id: T-1214
type: task
nature: feature
title: A story's page lists its conversations, and the inbox counts no message
status: done
parent: S-0336
owner: alex
created: 2026-10-07T20:17:06Z
updated: 2026-10-08T06:31:18Z
transitions:
  - to: ready
    at: 2026-10-08T06:28:04Z
    by: agent-S-0336
  - to: in-progress
    at: 2026-10-08T06:28:05Z
    by: agent-S-0336
  - to: done
    at: 2026-10-08T06:31:18Z
    by: agent-S-0336
stream: S-0336
tags: [flaiover]
touches: ["flaiover/src/routes/items/[id]/+page.svelte", "flaiover/src/routes/items/[id]/item.svelte.test.ts", flai/internal/hostapi/hostapi_test.go]
after: [T-1213]
usage:
  source: log
  seconds: 193
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 30
      output: 9150
      cache_read: 1645356
      cache_write: 51417
      cost: 0.8359
---
# T-1214 A story's page lists its conversations, and the inbox counts no message

## Work

Put messages where the operator already looks. It waits for T-1213, whose `Messages.svelte` it reuses.

- A story's page shows its conversations with `Messages.svelte`, open first, under its threads.
- The inbox badge and the inbox view stay as they are: a test shows a message does not change their count.
- The held card is not this task's: it went to the story split from S-0336.

## Done when

- Tests cover the story page with and without conversations, and the inbox count.
- `flai test` passes on the changed paths.
- `flai check --strict` passes.

## Notes
