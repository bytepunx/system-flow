---
id: T-2491
type: task
nature: feature
title: The Deployments page lists main's namespace and each story's with its commit, outcome, logs, and demo access, and deploys, tears down, or promotes through the dashboard's API routes
status: done
parent: S-0581
owner: alex
created: 2026-11-06T21:46:17Z
updated: 2026-11-07T03:03:13Z
transitions:
  - to: ready
    at: 2026-11-07T02:55:37Z
    by: agent-S-0581
  - to: in-progress
    at: 2026-11-07T02:55:38Z
    by: agent-S-0581
  - to: done
    at: 2026-11-07T03:03:13Z
    by: agent-S-0581
stream: S-0581
tags: [kluster, flaiover]
touches: [flaiover/src/lib/server/agent.ts, flaiover/src/lib/server/repo.ts, flaiover/src/routes/api/kluster/deployments/+server.ts, flaiover/src/routes/api/kluster/deploy/+server.ts, flaiover/src/routes/api/kluster/teardown/+server.ts, flaiover/src/routes/api/kluster/promote/+server.ts, flaiover/src/routes/api/kluster/deployments.test.ts, flaiover/src/routes/ops/deployments/+page.svelte, flaiover/src/routes/ops/deployments/deployments.svelte.test.ts]
after: [T-2489, T-2490]
usage:
  source: log
  seconds: 455
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 67
      output: 34980
      cache_read: 9356272
      cache_write: 152987
      cost: 3.6451
---
# T-2491 The Deployments page lists main's namespace and each story's with its commit, outcome, logs, and demo access, and deploys, tears down, or promotes through the dashboard's API routes

## Work

It waits for T-2489, whose host API methods the routes call, and T-2490, whose rail entry the page answers.

- `flaiover/src/lib/server/agent.ts`: allow `kluster.deployments`, `kluster.deploy`, `kluster.teardown`, and `kluster.promote`.
- `flaiover/src/lib/server/repo.ts`: a function for each, as S-0580's `klusterStatus` and `klusterUp` are.
- `flaiover/src/routes/api/kluster/deployments/+server.ts` (GET), and `deploy`, `teardown`, and `promote` (POST, the first two with a story), keeping `?project=`.
- `flaiover/src/routes/ops/deployments/+page.svelte`: the default namespace first, then each story's row. Each row has the story's ID and title linking its item page, the namespace, the commit, when it was deployed, the outcome, a link to each step's log, and how to reach it for a demo. Give each row the actions that fit it: Deploy or Redeploy for a story in review not deployed at its head, Teardown with confirmation first, and Promote on main's row. While a deploy runs, the row shows it and the page polls until it ends. flai's refusal or error text is shown verbatim. With the action off, or no block, the page says which, as the Klusters page does. Build it from the design system's tokens and the existing components.

Test the routes in `deployments.test.ts` and the page in `deployments.svelte.test.ts`: the list, each write, a refusal, a running deploy, an empty cluster, and the action off.

## Done when

- The page lists and manages the deployments as above, from the rail's Deployments entry.
- Every case above is tested, and `flai test` on the nine paths passes, svelte-check included.

## Notes
