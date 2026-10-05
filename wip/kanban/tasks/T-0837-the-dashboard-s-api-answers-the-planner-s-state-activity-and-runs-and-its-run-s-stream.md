---
id: T-0837
type: task
nature: feature
title: The dashboard's API answers the planner's state, activity, and runs, and its run's stream
status: done
parent: S-0259
owner: alex
created: 2026-10-05T00:07:41Z
updated: 2026-10-05T00:12:03Z
transitions:
  - to: ready
    at: 2026-10-05T00:08:19Z
    by: agent-S-0259
  - to: in-progress
    at: 2026-10-05T00:08:20Z
    by: agent-S-0259
  - to: done
    at: 2026-10-05T00:12:03Z
    by: agent-S-0259
stream: S-0259
tags: []
touches: [flaiover/src/routes/api/planner, flaiover/src/routes/api/agent-stream, flaiover/src/lib/planner.ts, flaiover/src/lib/planner.test.ts]
usage:
  source: log
  seconds: 223
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 37
      output: 197
      cache_read: 971711
      cache_write: 57140
      cost: 0.4185
---
# T-0837 The dashboard's API answers the planner's state, activity, and runs, and its run's stream

## Work

What the Planner page reads, in the dashboard's server routes and a pure module the page and the routes share:

- `GET /api/planner`: `{ plan_enabled, activity, runs }`. `plan_enabled` is `project.info`'s `host_actions.plan`, false when flai cannot be asked; `activity` is flai's `activity.document` for `planner`; `runs` is every item's newest planner run from `agent.status`'s `plans`, an empty list when flai cannot be asked.
- `GET /api/agent-stream/[story]?plan`: with `plan`, flai's `agent.stream` with `plan` in place of `story`, for the planner run of the epic or story named.
- `flaiover/src/lib/planner.ts`: the activity document's types; its entries newest first; the runs newest first, each with the cost and seconds of the activity entries that ended while it ran (from its start to its end, or now while it runs); the current run, the newest that has not ended; and a run's outcome as a word.

It waits for no task: flai's answers are given by the tests, in the shapes T-0836's `## Work` names, so it runs beside T-0836.

## Done when

- [x] `/api/planner` and `/api/agent-stream/[story]?plan` answer as above, with route tests beside them
- [x] `planner.ts` has behaviour tests for the ordering, the cost of a run, the current run, and the outcome
- [x] `npx vitest run` passes for the files it changed

## Notes
