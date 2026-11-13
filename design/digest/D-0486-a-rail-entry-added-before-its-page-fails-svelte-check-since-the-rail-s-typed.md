---
id: D-0486
title: A rail entry added before its page fails svelte-check, since the rail's typed resolve() needs the route to exist, so the task that adds it adds a placeholder +page.svelte that the page's task replaces
kind: decision
sources: [S-0581, T-2490, S-0757, T-3184]
topics: [flaiover, sveltekit, dashboard, planning, testing]
paths: [flaiover/src/lib/sitemenu.ts, flaiover/src/routes/ops/deployments/+page.svelte, flaiover/src/lib/components/StoryAgent.svelte, flaiover/src/routes/ops/images/+page.svelte]
date: 2026-11-07
status: active
absorbed_by: ""
asked_again: []
cited: [S-0757]
---

# D-0486 A rail entry added before its page fails svelte-check, since the rail's typed resolve() needs the route to exist, so the task that adds it adds a placeholder +page.svelte that the page's task replaces

S-0581 added the rail's Ops entry Deployments at `/ops/deployments` (T-2490) in a layer before the page itself (T-2491). Its narrative, `wip/archive/agents/S-0581.md`, records under Decisions:

> T-2490 adds a placeholder `flaiover/src/routes/ops/deployments/+page.svelte`, added to its touches: the rail's typed `resolve('/ops/deployments')` fails svelte-check without the route. T-2491 replaces it.

A rail entry is a link the page `resolve()`s (`design/system/flaiover-dashboard.md`), and SvelteKit types `resolve()` from the routes that exist. So a planner or a story's agent who puts a rail entry ahead of its page gives that task the page's `+page.svelte` path as well, and lets the page's task replace it. The other choice is to put the rail entry in the same task as its page.

S-0757 met it again, and the rule holds for any component's link, not only the rail's. Its plan took the other choice for the rail: T-3185 owns `sitemenu.ts` with the `/ops/images` page. But `StoryAgent.svelte` links to that page from an earlier layer, so its narrative, `wip/archive/agents/S-0757.md`, records:

> T-3184 added a placeholder /ops/images page that T-3185 replaced, so that StoryAgent's typed resolve('/ops/images') type-checks (D-0486).
