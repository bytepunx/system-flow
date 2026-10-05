---
id: T-0889
type: task
nature: feature
title: The dashboard shows the orchestrator's run on the Activity page and the orchestrate action on the Settings page
status: backlog
parent: S-0218
owner: alex
created: 2026-10-05T04:46:38Z
updated: 2026-10-05T04:47:21Z
transitions: []
stream: S-0218
tags: [dashboard]
touches: [flaiover/src/lib/activity.ts, flaiover/src/lib/server/agent.ts, flaiover/src/lib/server/agent.test.ts, flaiover/src/routes/activity/+page.svelte, flaiover/src/routes/activity/activity.svelte.test.ts, design/system/flaiover-dashboard.md, docs/users/flaiover.md]
after: [T-0883]
---
# T-0889 The dashboard shows the orchestrator's run on the Activity page and the orchestrate action on the Settings page

## Work

Make the orchestrator visible where the dashboard already shows agent runs, as S-0208 did for the planner:

- `flaiover/src/lib/activity.ts` and `flaiover/src/lib/server/agent.ts`: `agent.status` carries the run under `orchestrator`, as T-0883 writes it; type it beside `plans`, a run with no story and no item.
- `flaiover/src/routes/activity/+page.svelte`: the Activity page lists the orchestrator's run with its stream, as it lists a planner's, and hears its start and end from the `agent` notification.
- The Settings page lists `orchestrate` with the other host actions, from `hostapi.Actions`; add it by hand only if that list is not read from flai.
- `design/system/flaiover-dashboard.md` and `docs/users/flaiover.md` say so.

The orchestrator's own page, with its activity log, its decisions, and Stop, is S-0228's, which waits for this story. Its decisions and refusals are in `wip/agents/orchestrator.md`, which `activity.document` already answers for kind `orchestrator` (S-0206).

This task waits for T-0883, which puts the run in `agent.status`. It runs with T-0892, whose paths it does not share.

## Done when

- the Activity page shows a running orchestrator and its stream, and a component test pins it
- `agent.ts` reads a status with an `orchestrator` run and one without, and a test pins both
- the Settings page lists `orchestrate`
- `flaiover-dashboard.md` and `flaiover.md` describe both
- the dashboard's lint, check, and tests pass

## Notes
