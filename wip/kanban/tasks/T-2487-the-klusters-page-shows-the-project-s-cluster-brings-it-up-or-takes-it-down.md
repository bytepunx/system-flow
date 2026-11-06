---
id: T-2487
type: task
nature: feature
title: The Klusters page shows the project's cluster, brings it up or takes it down through the dashboard's API routes, and says when the host action is off or the block is missing
status: backlog
parent: S-0580
owner: alex
created: 2026-11-06T21:45:32Z
updated: 2026-11-06T21:45:32Z
transitions: []
stream: S-0580
tags: [kluster, flaiover]
touches: [flaiover/src/lib/server/agent.ts, flaiover/src/lib/server/repo.ts, flaiover/src/routes/api/kluster/+server.ts, flaiover/src/routes/api/kluster/up/+server.ts, flaiover/src/routes/api/kluster/down/+server.ts, flaiover/src/routes/api/kluster/kluster.test.ts, flaiover/src/routes/ops/klusters/+page.svelte, flaiover/src/routes/ops/klusters/klusters.svelte.test.ts]
after: [T-2485, T-2486]
---
# T-2487 The Klusters page shows the project's cluster, brings it up or takes it down through the dashboard's API routes, and says when the host action is off or the block is missing

## Work

It waits for T-2485, whose host API methods the routes call, and T-2486, whose rail entry the page answers.

- `flaiover/src/lib/server/agent.ts`: allow `kluster.status`, `kluster.up`, and `kluster.down`, beside the components methods.
- `flaiover/src/lib/server/repo.ts`: `klusterStatus`, `klusterUp`, and `klusterDown`, as `componentsStatus`, `componentsStart`, and `componentsStop` are, with writes through `write`.
- `flaiover/src/routes/api/kluster/+server.ts` (GET, the status), `up/+server.ts` and `down/+server.ts` (POST), as `routes/api/components` does, keeping `?project=`.
- `flaiover/src/routes/ops/klusters/+page.svelte`: the cluster's name, profile, provider, whether it runs, and its age. Up and Down buttons, with Down asking for confirmation first. While one runs, the page says so, polls the status until it ends, and links the log. flai's refusal or error text is shown verbatim. With the action off it says so and links the Features page, and with no block it names the block and links `design/system/kluster.md`, each in place of the buttons. Build it from the design system's tokens and the existing components.

Test the routes in `kluster.test.ts` and the page in `klusters.svelte.test.ts`: the status, each write, a refusal, a running action, the action off, and no block.

## Done when

- The page shows and manages the cluster as above, from the rail's Klusters entry.
- Every case above is tested, and `flai test` on the eight paths passes, svelte-check included.

## Notes
