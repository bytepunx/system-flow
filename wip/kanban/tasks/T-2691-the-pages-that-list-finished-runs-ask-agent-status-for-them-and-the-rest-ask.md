---
id: T-2691
type: task
nature: improvement
title: The pages that list finished runs ask agent.status for them, and the rest ask for the open and live runs
status: done
parent: S-0611
owner: alex
created: 2026-11-07T23:56:20Z
updated: 2026-11-09T17:24:15Z
transitions:
  - to: ready
    at: 2026-11-09T17:00:56Z
    by: system-flow
  - to: in-progress
    at: 2026-11-09T17:17:22Z
    by: system-flow
  - to: done
    at: 2026-11-09T17:24:14Z
    by: agent-S-0611
stream: S-0611
tags: [dashboard]
touches: [flaiover/src/routes/api/planner/+server.ts, flaiover/src/routes/api/planner/planner.test.ts, flaiover/src/routes/api/orchestrator/+server.ts, flaiover/src/routes/api/orchestrator/orchestrator.test.ts, flaiover/src/routes/api/analyzer/+server.ts, flaiover/src/routes/api/analyzer/analyzer.test.ts, "flaiover/src/routes/api/roles/[name]/+server.ts", "flaiover/src/routes/api/roles/[name]/roles.test.ts", "flaiover/src/routes/api/items/[id]/plan/+server.ts", "flaiover/src/routes/api/items/[id]/plan/plan.test.ts", flaiover/src/routes/api/host-agent/host-agent.test.ts, design/system/flaiover-dashboard.md, flaiover/src/lib/components/StoryAgent.svelte, flaiover/src/lib/components/StoryAgent.svelte.test.ts, flaiover/src/routes/api/host-agent/+server.ts]
after: [T-2690]
usage:
  source: log
  seconds: 412
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 58
      output: 368
      cache_read: 10249337
      cache_write: 45317
      cost: 3.7575
---
# T-2691 The pages that list finished runs ask agent.status for them, and the rest ask for the open and live runs

## Work

Criterion 2, the dashboard side. Every flaiover caller of `agent.status` asks it with no parameter today.

- Ask with `{ finished: true }` where the page lists finished runs:
  - `/api/planner` (`flaiover/src/routes/api/planner/+server.ts`), every item's newest planner run;
  - `/api/orchestrator` and `/api/analyzer`, their past runs;
  - `/api/roles/[name]`, a role's runs;
  - `/api/items/[id]/plan`, the item's newest planner run, which may have ended.
- Leave the default on `/api/host-agent` (the board, the activity page, a story's agent panel, polled every 15 s), `/api/projects`, `/api/agent`, and the layout's roles: they show only what is open or live. Pin that in `host-agent.test.ts`.
- If the activity page shows ended runs of done stories, say in the story's `## Decisions` whether it asks with `finished` too, and add its route to this task's touches.

Waits for T-2690, which adds the `finished` parameter and decides what the default answer keeps.

## Done when

- [ ] Each route that lists finished runs asks `agent.status` with `finished: true`, and its test checks the params it sends.
- [ ] `host-agent.test.ts` checks `/api/host-agent` asks with no `finished`.
- [ ] `flai test` passes on the changed paths, the `flaiover-flai` tier included.

## Notes

- Planned by planner-S-0611.
