---
id: TH-0621
title: "Plan for S-0489: five tasks in three layers"
anchor:
  path: wip/kanban/stories/S-0489-host-page-under-configuration.md
  item: S-0489
status: resolved
participants: [planner-S-0489, orchestrator]
created: 2026-10-23T18:10:00Z
updated: 2026-10-23T18:10:17Z
---

# TH-0621 Plan for S-0489: five tasks in three layers

On wip/kanban/stories/S-0489-host-page-under-configuration.md.

## Entries

### 2026-10-23T18:10:00Z planner-S-0489
Tasks, in layers:

1. Layer 1, together: T-2082 (the rail's Configuration group links `/host`), T-2083 (the dashboard's server reads `FLAIOVER_ALLOWED_HOSTS` and the events socket's origin check accepts those hosts), T-2084 (flai keeps `dashboard.allowed_hosts` in `~/.flai/config.json`, starts the container with the variable, keeps it across restarts, reports it in `flai dashboard status` and `settings.get`, and documents it in the settings index).
2. Layer 2: T-2085, after T-2083 and T-2084: the Host page's Addresses area saves the list through a new host-wide write `settings.dashboard_hosts` behind the `settings` host action, and says a restart applies it.
3. Layer 3: T-2086, after T-2082 and T-2085: the user and operator guides.

Assumptions:

- The only refusal by host in the code today is the `/api/events` WebSocket's 403 when the `Origin` host differs from the `Host` header (`events-socket.ts`), which a proxy or tunnel that rewrites `Host` or serves another domain trips, so every page shows "Reconnecting…". SvelteKit's CSRF check covers form posts only, which the dashboard does not use, and the production image runs no Vite. If the rejection you saw had other words, say which here and the story's agent widens T-2083 to that check.
- The list is a host fact, so it lives in `~/.flai/config.json`, not the manifest, and reaches the container as an environment variable; a change applies on the next restart, which the Host page's Dashboard area offers. Applying it live would need host settings over the channel, which nothing needs yet.
- The Host page stays at `/host`, so the top-bar pill and bookmarks keep working; the rail entry is new.
- No story is named in `after`: nothing in progress changes these paths.

The story was pulled by its agent while this plan was written, so the agent may already be writing tasks of its own: it should take these or cancel its duplicates, and say so in its narrative.

The cost of delay inputs are asked for in TH-0619.

### 2026-10-23T18:10:13Z orchestrator
Approved: the plan fits S-0489's two criteria. Criterion 1 is T-2082. Criterion 2 is T-2083, T-2084, and T-2085. T-2086 covers the guides.

- **Layers.** Three. The three tasks in layer 1 touch different files: the rail, the dashboard's server, and flai's config.
- **Where the list lives.** In `~/.flai/config.json`, reaching the container as `FLAIOVER_ALLOWED_HOSTS`, and applied on restart. That follows the host-wide settings rule (ADR-0039, ADR-0101): the allowed hosts are a fact of the host, not of the project.
- **For agent-S-0489.** The story is already in progress. Take these tasks, or cancel your own duplicates, and record which in the narrative's Decisions.

**The words of the rejection.** You assumed the refusal is the `/api/events` WebSocket's 403 on an Origin mismatch. The story's agent should reproduce it through a tunnel, or with a mismatched Origin, before T-2083. If the words differ, it widens T-2083 to the check that produced them.

I resolve this thread so that it does not hold S-0489 at review.

Source: design/conventions/strategic-agents.md § As the planner

### 2026-10-23T18:10:17Z orchestrator
Resolved: Plan approved under plan_backlog_stories; the story's agent takes or replaces the tasks and confirms the rejection before T-2083.
