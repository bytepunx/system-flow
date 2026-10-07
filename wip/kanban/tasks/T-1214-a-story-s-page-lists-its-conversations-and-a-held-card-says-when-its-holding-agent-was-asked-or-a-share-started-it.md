---
id: T-1214
type: task
nature: feature
title: A story's page lists its conversations, and a held card says when its holding agent was asked or a share started it
status: backlog
parent: S-0336
owner: alex
created: 2026-10-07T20:17:06Z
updated: 2026-10-07T20:17:06Z
transitions: []
stream: S-0336
tags: [flaiover]
touches: ["flaiover/src/routes/items/[id]/+page.svelte", "flaiover/src/routes/items/[id]/item.svelte.test.ts", flaiover/src/lib/components/BoardCard.svelte, flaiover/src/lib/components/BoardCard.svelte.test.ts]
after: [T-1213]
---
# T-1214 A story's page lists its conversations, and a held card says when its holding agent was asked or a share started it

## Work

Put messages where the operator already looks. It waits for T-1213, whose `Messages.svelte` it reuses.

- A story's page shows its conversations with `Messages.svelte`, open first, under its threads.
- A held card whose holding agent was asked (S-0334) says so and links the conversation; a card of a story started on a share names the shared paths.
- The inbox badge and the inbox view stay as they are: a test shows a message does not change their count.

## Done when

- Tests cover the story page with and without conversations, the held card asked, the card started on a share, and the inbox count.
- `flai test` passes on the changed paths.
- `flai check --strict` passes.

## Notes
