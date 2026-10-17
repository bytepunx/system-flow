---
id: T-1621
type: task
nature: feature
title: The dashboard's story and review pages show a story's pull request, its required checks, and its reviews
status: done
parent: S-0394
owner: alex
created: 2026-10-16T04:37:08Z
updated: 2026-10-17T16:49:15Z
transitions:
  - to: ready
    at: 2026-10-17T16:31:01Z
    by: agent-S-0394
  - to: in-progress
    at: 2026-10-17T16:31:02Z
    by: agent-S-0394
  - to: done
    at: 2026-10-17T16:49:15Z
    by: agent-S-0394
stream: S-0394
tags: [dashboard]
touches: ["flaiover/src/routes/api/items/[id]/forge/+server.ts", "flaiover/src/routes/api/items/[id]/forge/forge.test.ts", flaiover/src/lib/components/ForgeStatus.svelte, flaiover/src/lib/components/ForgeStatus.svelte.test.ts, "flaiover/src/routes/items/[id]/+page.svelte", "flaiover/src/routes/items/[id]/item.svelte.test.ts", flaiover/src/lib/components/Review.svelte, flaiover/src/lib/components/Review.svelte.test.ts, design/components/flaiover/src/lib/components/ForgeStatus.examples.json]
after: [T-1619]
usage:
  source: log
  seconds: 1093
  estimated: true
  priced_by: harness
  models:
    - model: claude-haiku-5-5
      input: 124
      output: 47706
      cache_read: 8429901
      cache_write: 214196
      cost: 0.5703
    - model: claude-opus-5-5
      input: 80
      output: 33707
      cache_read: 4825940
      cache_write: 145785
      cost: 2.4647
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
