---
id: T-1764
type: task
nature: feature
title: "A RolePage component and the /workflow/[role] route show any role's action, run and stream, activity log, runs, and its Run, Stop, and Start"
status: backlog
parent: S-0420
owner: alex
created: 2026-10-16T05:46:46Z
updated: 2026-10-16T05:50:22Z
transitions: []
stream: S-0420
tags: [dashboard]
touches: [flaiover/src/lib/components/RolePage.svelte, flaiover/src/lib/components/RolePage.svelte.test.ts, flaiover/src/lib/roles.ts, flaiover/src/lib/roles.test.ts, "flaiover/src/routes/workflow/[role]/+page.svelte", "flaiover/src/routes/api/roles/[name]/+server.ts", "flaiover/src/routes/api/roles/[name]/roles.test.ts", flaiover/src/lib/agent.ts, flaiover/src/lib/server/agent.ts]
after: [T-1762]
---
# T-1764 A RolePage component and the /workflow/[role] route show any role's action, run and stream, activity log, runs, and its Run, Stop, and Start

## Work

- Build `RolePage.svelte` from what the Planner, Orchestrator, and Analyzer pages share today: the action's state, the run under way with `AgentStream`, the activity log, and the runs with their cost and outcome. Add **Run** (with an item picker when the role takes one), **Stop**, and **Start**, offered as the role's trigger allows.
- `roles.ts` types the `roles` entry of `agent.status`, and `agent.ts`'s `Role` type takes any name.
- `/workflow/[role]` renders the page for a defined role. `/api/roles/[name]` posts `role.run`, `role.stop`, and `role.start`.
- Write component and API tests from a fixture `agent.status` with a role `reviewer`.

Waits for the status task, whose fields it reads.

## Done when

`pnpm vitest` on the new files passes, and `flai test` on `flaiover/src` is clean.

## Notes
