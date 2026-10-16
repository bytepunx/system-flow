---
id: T-1621
type: task
nature: feature
title: The dashboard's story and review pages show a story's pull request, its required checks, and its reviews
status: backlog
parent: S-0394
owner: alex
created: 2026-10-16T04:37:08Z
updated: 2026-10-16T04:37:08Z
transitions: []
stream: S-0394
tags: [dashboard]
touches: [flaiover/src/lib/server/agent.ts, "flaiover/src/routes/api/items/[id]/forge/+server.ts", "flaiover/src/routes/api/items/[id]/forge/forge.test.ts", flaiover/src/lib/components/ForgeStatus.svelte, flaiover/src/lib/components/ForgeStatus.svelte.test.ts, "flaiover/src/routes/items/[id]/+page.svelte", "flaiover/src/routes/items/[id]/item.svelte.test.ts", "flaiover/src/routes/review/[id]/+page.svelte"]
after: [T-1619]
---
# T-1621 The dashboard's story and review pages show a story's pull request, its required checks, and its reviews

## Work

- Allow `forge.status` in `flaiover/src/lib/server/agent.ts`.
- Add `GET /api/items/[id]/forge`, answering it, with a test as the `checks` route has.
- Add `ForgeStatus.svelte`, built from the design system's tokens. It shows:
  - the link to the pull request
  - one row per required check, with its state and its log link
  - one row per review, with its author and state
- It shows nothing when there is no pull request.
- Put it on `items/[id]/+page.svelte` and `review/[id]/+page.svelte`, and test it and the item page.
- Waits for T-1619: it calls that host method.

## Done when

Both pages show the pull request, its checks, and its reviews for a story that has one, and nothing for one that has none, and the flaiover tests pass.

## Notes
