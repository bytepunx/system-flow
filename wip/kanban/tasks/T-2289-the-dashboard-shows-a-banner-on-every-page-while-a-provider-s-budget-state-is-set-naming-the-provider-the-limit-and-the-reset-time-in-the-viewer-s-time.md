---
id: T-2289
type: task
nature: improvement
title: The dashboard shows a banner on every page while a provider's budget state is set, naming the provider, the limit, and the reset time in the viewer's time
status: backlog
parent: S-0534
owner: alex
created: 2026-10-28T21:27:16Z
updated: 2026-10-28T21:27:16Z
transitions: []
stream: S-0534
tags: [flaiover, dashboard, budget]
touches: [flaiover/src/lib/server/agent.ts, flaiover/src/lib/server/agent.test.ts, flaiover/src/routes/api/agent/+server.ts, flaiover/src/routes/api/agent/agent.test.ts, flaiover/src/lib/hostflai.svelte.ts, flaiover/src/lib/components/BudgetBanner.svelte, flaiover/src/lib/components/BudgetBanner.svelte.test.ts, flaiover/src/routes/+layout.svelte, flaiover/src/routes/layout.svelte.test.ts]
after: [T-2287]
---
# T-2289 The dashboard shows a banner on every page while a provider's budget state is set, naming the provider, the limit, and the reset time in the viewer's time

## Work

Carry `agent.status`'s `budget` (T-2287) to the browser:

- `status()` in `flaiover/src/lib/server/agent.ts` reads it.
- `/api/agent` returns it beside the connection status.
- `hostflai.svelte.ts`, which polls that route, keeps it.

Add `BudgetBanner.svelte`. It shows one line per held provider, such as `claude-code is out of budget (weekly limit) until 09:00; flai serve holds every agent start until then`, with the time in the viewer's zone as the dashboard's other times are. While a probe runs, it says the probe is running. It shows nothing when `budget` is empty.

Place it in `+layout.svelte` beside `HostFlaiBanner` and `UnsignedBanner`, so it shows on every page while signed in. Build it from the design system's tokens and the generated stylesheet, as the other banners are. If `flai check` reports `components.examples` for the new component, add its examples file under `design/components/` and add that path to this task's touches.

This task waits for T-2287, which puts `budget` on `agent.status`.

## Done when

- [ ] The banner shows on every page with the provider, the limit, and `until` while `budget` is set, and is gone once it clears, with a component test and a layout test.
- [ ] `/api/agent` returns `budget`, tested.
- [ ] `flai test` passes on the changed paths, svelte-check included.

## Notes

Planned by planner-S-0534.
