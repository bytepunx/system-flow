---
id: T-3365
type: task
nature: feature
title: The review page's Rebase and resolve button asks the story's agent, through accept.resolve, to take the recommended actions that hold its acceptance
status: done
parent: S-0788
owner: alex
created: 2026-11-14T18:08:52Z
updated: 2026-11-14T18:46:03Z
transitions:
  - to: ready
    at: 2026-11-14T18:09:11Z
    by: agent-S-0788
  - to: in-progress
    at: 2026-11-14T18:39:12Z
    by: agent-S-0788
  - to: done
    at: 2026-11-14T18:46:03Z
    by: agent-S-0788
stream: S-0788
tags: []
touches: [flaiover/src/lib/components/Review.svelte, flaiover/src/lib/components/Review.svelte.test.ts, "flaiover/src/routes/api/items/[id]/resolve/+server.ts", "flaiover/src/routes/api/items/[id]/resolve/resolve.test.ts", flaiover/src/lib/server/agent.ts, docs/users/flaiover.md, design/system/flaiover-dashboard.md]
after: [T-3364]
usage:
  source: log
  seconds: 411
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 58
      output: 19581
      cache_read: 3500638
      cache_write: 180515
      cost: 2.1724
---

# T-3365 The review page's Rebase and resolve button asks the story's agent, through accept.resolve, to take the recommended actions that hold its acceptance

## Work

Waits for T-3364, whose host method `accept.resolve` this calls and whose answer it shows. Add the route `POST /api/items/[id]/resolve`, which calls `accept.resolve`, and allow the method in `flaiover/src/lib/server/agent.ts`. On the review page, put a **Rebase and resolve** button in the box "This cannot be accepted from here yet", beside the blockers it lists, and beside a refused acceptance's or the merge queue's refusal. It sends the story back through the route, and the page then shows the send-back notice the acceptance's own send-back shows, with its thread. Document it in `docs/users/flaiover.md` under Reviewing a story and in `design/system/flaiover-dashboard.md` where the review page is described.

## Done when

- The button shows while a blocker holds Accept, and its press sends the story back and shows the notice, with a component test.
- The route calls `accept.resolve` and passes on its refusal, with a test.
- `flai test` passes on the changed paths.

## Notes
