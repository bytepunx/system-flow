---
id: T-0921
type: task
nature: feature
title: The dashboard's API answers the orchestrator's and the analyzer's state, activity, and runs, and runs the analyzer and stops and starts the orchestrator
status: done
parent: S-0228
owner: alex
created: 2026-10-05T05:44:54Z
updated: 2026-10-07T00:23:18Z
transitions:
  - to: ready
    at: 2026-10-07T00:03:28Z
    by: agent-S-0228
  - to: in-progress
    at: 2026-10-07T00:03:29Z
    by: agent-S-0228
  - to: done
    at: 2026-10-07T00:23:18Z
    by: agent-S-0228
stream: S-0228
tags: [dashboard]
touches: [flaiover/src/routes/api/orchestrator, flaiover/src/routes/api/analyzer, flaiover/src/routes/api/agent-stream, flaiover/src/lib/strategic.ts, flaiover/src/lib/strategic.test.ts, flaiover/src/lib/planner.ts, flaiover/src/lib/planner.test.ts, flaiover/src/lib/server/agent.ts]
usage:
  source: log
  seconds: 1189
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 104
      output: 658
      cache_read: 4717939
      cache_write: 139809
      cost: 2.1526
---
# T-0921 The dashboard's API answers the orchestrator's and the analyzer's state, activity, and runs, and runs the analyzer and stops and starts the orchestrator

## Work

What the two pages read and write, in the dashboard's server routes and a pure module they share, built as S-0259's `/api/planner` and `planner.ts` were:

- `flaiover/src/lib/strategic.ts`: what `planner.ts` does for the planner, for any strategic agent: the activity document's types, its entries newest first, runs newest first, each with the cost and seconds of the activity entries that ended while it ran, the current run, and a run's outcome as a word. `planner.ts` keeps its exports and uses it.
- `GET /api/orchestrator`: `{ enabled, held, activity, run, runs }`, from `project.info`'s `host_actions.orchestrate`, `agent.status`'s `orchestrator` run and its `held`, and `activity.document` for `orchestrator`. Its decisions are the document's entries, each summary carrying its reason.
- `GET /api/analyzer`: `{ enabled, activity, runs }`, from `host_actions.analyze`, the analyzer runs in `agent.status`, and `activity.document` for `analyzer`.
- `POST /api/orchestrator` with `{ action: 'stop' | 'start' }`: flai's `orchestrate.stop` or `orchestrate.start`. `POST /api/analyzer` with `{ focus }`: flai's `analyze.run`, the focus `bottlenecks`, `intent`, or `risk`, or none. Each says what flai answered, and why when the action is off.
- `GET /api/agent-stream/[story]?role=orchestrate|analyze`: flai's `agent.stream` with `role`.
- `flaiover/src/lib/server/agent.ts` types the analyzer runs and the orchestrator's `held` in `agent.status`.

Every answer falls back to off and empty when flai cannot be asked, as `/api/planner`'s does.

It waits for no task: flai's answers are given by the tests, in the shapes T-0917's `## Work` names, so it runs beside T-0917, whose paths it does not share.

## Done when

- [ ] `/api/orchestrator`, `/api/analyzer`, and `/api/agent-stream/[story]?role` answer as above, refuse a focus, action, or role that is not one, and have route tests beside them
- [ ] `strategic.ts` has behaviour tests for the ordering, the cost of a run, the current run, and the outcome, and `planner.ts`'s tests still pass
- [ ] `npx vitest run` passes for the files it changed

## Notes
