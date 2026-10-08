---
id: T-1325
type: task
nature: feature
title: The Messages component shows a conversation's share and its split
status: backlog
parent: S-0338
owner: alex
created: 2026-10-08T04:32:02Z
updated: 2026-10-08T04:32:02Z
transitions: []
stream: S-0338
tags: [flaiover]
touches: [flaiover/src/lib/components/Messages.svelte, flaiover/src/lib/components/Messages.svelte.test.ts]
---
# T-1325 The Messages component shows a conversation's share and its split

## Work

Show the share where the conversation is read. It waits for nothing in this story: S-0336 built `Messages.svelte`, and S-0334 put the share in the conversation the host reads serve. It shares no path with the board task, so the two run together.

- `Messages.svelte` shows a conversation's share: the paths, the split in its holder's words, who shared, and when, and whether the share has ended.
- The Messages view and a story's page both render `Messages.svelte`, so neither page changes.

## Done when

- Component tests cover a conversation with a share, one whose share ended, and one without.
- `flai test` passes on the changed paths.
- `flai check --strict` passes.

## Notes
