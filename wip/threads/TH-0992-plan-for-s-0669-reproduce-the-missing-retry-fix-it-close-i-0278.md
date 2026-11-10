---
id: TH-0992
title: "Plan for S-0669: reproduce the missing Retry, fix it, close I-0278"
anchor:
  path: wip/kanban/stories/S-0669-after-a-stop-from-the-dashboard-the-story-s-agent-page-shows-the-agent-stopped.md
  item: S-0669
status: open
participants: [planner-S-0669]
created: 2026-11-10T08:48:26Z
updated: 2026-11-10T08:48:26Z
---

# TH-0992 Plan for S-0669: reproduce the missing Retry, fix it, close I-0278

On wip/kanban/stories/S-0669-after-a-stop-from-the-dashboard-the-story-s-agent-page-shows-the-agent-stopped.md.

## Entries

### 2026-11-10T08:48:26Z planner-S-0669
S-0669 is planned as three tasks in three layers. Each task waits for the one before it.

| Layer | Task | Touches |
|-------|------|---------|
| 1 | T-2907: a test reproduces the missing Retry, and the narrative names the cause and the fix proposed | `StoryAgent.svelte.test.ts`, `items/[id]/item.svelte.test.ts` |
| 2 | T-2909 (after T-2907): the story page offers Retry for a stopped agent whenever the card's menu does | `activity.ts`, `activity.test.ts`, `StoryAgent.svelte`, `items/[id]/+page.svelte`, `design/system/flaiover-dashboard.md` |
| 3 | T-2912 (after T-2909): I-0278 is closed with the cause and the fix | the I-0278 issue file |

Figures:

- Forecast: 45m, delivery 2026-11-10T11:49Z. flai gave 18m; I raised it because I-0278 was never reproduced.
- Cost of delay: 5.00 USD a week, from the 2m `time_lost_per_cycle` flai set from I-0278. It stands.
- After: none.

Assumptions:

- The page and the card's menu share `retryable()`. So the cause is one of the inputs that differ: the page's `item.status`, its `writable && !item.archived`, its own `/api/host-agent` answer, or the `retried` guard in `StoryAgent.svelte`. If it is none of these, the cause is in flai's `agent.status`, and T-2909's touches widen to `flai/internal/serve/agents.go`.
- The goal says to propose a fix before building it. T-2907 records the proposal in the narrative's `## Decisions`. It asks on this story only when the fix would change what ADR-0058 §4 says a stopped run offers.
- The story is a draft with no `after`, tagged `flaiover` and `agent`, with topics `dashboard` and `agents`.

No split, merge, or drop is proposed.
