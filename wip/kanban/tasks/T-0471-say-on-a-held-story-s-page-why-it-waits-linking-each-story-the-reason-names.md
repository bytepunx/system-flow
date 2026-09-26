---
id: T-0471
type: task
nature: feature
title: Say on a held story's page why it waits, linking each story the reason names
status: ready
parent: S-0129
owner: alex
created: 2026-09-26T17:54:07Z
updated: 2026-09-26T17:54:12Z
transitions:
  - to: ready
    at: 2026-09-26T17:54:12Z
    by: agent-S-0129
stream: S-0129
tags: []
touches: [flaiover/src/lib/components/StoryAgent.svelte]
---
# T-0471 Say on a held story's page why it waits, linking each story the reason names

## Work

- `StoryAgent.svelte` shows the hold's full reason with every story ID in it linked to that story's page, and no empty "started" line for a story that never had an agent.
- It keeps **Start agent** for a held story that never had an agent, since `flai serve agent start` overrides a hold (ADR-0046).
- The page's header shows `HELD` after `BLOCKED`.
- Tests in `StoryAgent.svelte.test.ts`.

## Done when

- The page shows the full reason with links, clears on the next change or poll without a reload, and tests pass.

## Notes
