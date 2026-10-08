---
id: T-1214
type: task
nature: feature
title: A story's page lists its conversations, and the inbox counts no message
status: backlog
parent: S-0336
owner: alex
created: 2026-10-07T20:17:06Z
updated: 2026-10-08T04:31:25Z
transitions: []
stream: S-0336
tags: [flaiover]
touches: ["flaiover/src/routes/items/[id]/+page.svelte", "flaiover/src/routes/items/[id]/item.svelte.test.ts"]
after: [T-1213]
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
