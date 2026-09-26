---
id: T-0470
type: task
nature: feature
title: Show a held story's yellow dot and a short hold line on its board card
status: in-progress
parent: S-0129
owner: alex
created: 2026-09-26T17:54:07Z
updated: 2026-09-26T17:54:12Z
transitions:
  - to: ready
    at: 2026-09-26T17:54:12Z
    by: agent-S-0129
  - to: in-progress
    at: 2026-09-26T17:54:12Z
    by: agent-S-0129
stream: S-0129
tags: []
touches: [flaiover/src/lib]
---
# T-0470 Show a held story's yellow dot and a short hold line on its board card

## Work

- Give `StoryActivity` in `flaiover/src/lib/activity.ts` the `hold` flai's `agent.status` carries (`{code, reason}`), a helper that reads the story IDs a hold waits for from the reason's `starts when` clause, and a short line (`held (overlap): S-0128`).
- Show a held story's activity whether or not the `agent` action is on, since `agent.status` reports holds either way; `activityLine` says the hold's reason rather than "agent waiting ()".
- `BoardCard.svelte` prints the short line in the details row after `BLOCKED`.
- Tests in `activity.test.ts` and `BoardCard.svelte.test.ts`.

## Done when

- A held ready story's card shows a yellow dot and the short line, for a story that never had an agent too, and `BLOCKED` precedes it on a blocked card; tests pass.

## Notes
