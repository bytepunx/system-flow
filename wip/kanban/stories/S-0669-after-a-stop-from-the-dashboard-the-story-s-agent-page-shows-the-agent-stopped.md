---
id: S-0669
type: story
nature: improvement
title: After a stop from the dashboard, the story's agent page shows the agent stopped by the operator with no Retry, while the card's menu offers Retry
status: backlog
owner: alex
created: 2026-11-10T08:34:44Z
updated: 2026-11-10T08:48:12Z
transitions: []
tags: [flaiover, agent]
topics: [dashboard, agents]
touches: [flaiover/src/lib/activity.ts, flaiover/src/lib/activity.test.ts, flaiover/src/lib/components/StoryAgent.svelte, flaiover/src/lib/components/StoryAgent.svelte.test.ts, "flaiover/src/routes/items/[id]/+page.svelte", "flaiover/src/routes/items/[id]/item.svelte.test.ts", design/system/flaiover-dashboard.md, design/issues/I-0278-after-a-stop-from-the-dashboard-the-story-s-agent-page-shows-the-agent-stopped.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
draft: true
cost_of_delay:
  inputs:
    time_lost_per_cycle: 2m
    by: flai
    at: 2026-11-10T08:34:44Z
  value: 5
  by: planner-S-0669
  at: 2026-11-10T08:48:12Z
forecast:
  duration: 45m
  delivery: 2026-11-10T11:49:00Z
  basis: flai forecast's 18m (107 s a unit over 30 medium improvement stories, size 10), raised to 45m because I-0278 was never reproduced and T-2907 must find the cause before T-2909 fixes it
  by: planner-S-0669
  at: 2026-11-10T08:48:12Z
---
# S-0669 After a stop from the dashboard, the story's agent page shows the agent stopped by the operator with no Retry, while the card's menu offers Retry

## Goal

This story remediates [I-0278](../../../design/issues/I-0278-after-a-stop-from-the-dashboard-the-story-s-agent-page-shows-the-agent-stopped.md), "After a stop from the dashboard, the story's agent page shows the agent stopped by the operator with no Retry, while the card's menu offers Retry". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0278 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0278 is closed with `flai issue close I-0278 --reason` saying what fixed it

## Tasks
- T-2907 A test reproduces the story page's missing Retry after a stop, and the narrative names its cause and the fix proposed
- T-2909 The story page offers Retry for a stopped agent whenever the card's menu does
- T-2912 I-0278 is closed with the cause and the fix

## Notes

Cost of delay inputs set by flai from I-0278. time_lost_per_cycle 2m: 2m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-11-07T06:28:23Z, 3.1 days before this story; under one cycle counts as one).

### Planning

Touches. The story declared none. `flai touches suggest` was run from the two files I-0278 names, `flaiover/src/lib/components/StoryAgent.svelte` and `flaiover/src/lib/activity.ts` (27 of 3463 commits changed them).

| Touch | Source |
|-------|--------|
| `flaiover/src/lib/activity.ts` | design: I-0278 names `retryable()` here |
| `flaiover/src/lib/components/StoryAgent.svelte` | design: I-0278 names `canRetry` here |
| `flaiover/src/lib/components/StoryAgent.svelte.test.ts` | co-change, 56% |
| `flaiover/src/lib/activity.test.ts` | co-change, 37% |
| `design/system/flaiover-dashboard.md` | co-change, 26%; its "Stopping an agent" and "Each story's agent" paragraphs describe Retry |
| `flaiover/src/routes/items/[id]/+page.svelte` | co-change, 26%; it passes the page's `status` and `writable` to `StoryAgent` |
| `flaiover/src/routes/items/[id]/item.svelte.test.ts` | co-change, 15%; layout: the page's own test, for a cause in its props |
| `design/issues/I-0278-after-a-stop-from-the-dashboard-the-story-s-agent-page-shows-the-agent-stopped.md` | criterion 2 closes it |

All are files; no folder touch is kept. `flai/internal/serve/agents.go` (co-change, 15%) is left out. It matters only if T-2907 finds that `agent.status` answers the page differently from the board, and `flai task done` then widens the touches.

Forecast. `flai forecast` gave 2m with no touches (size 2), then 18m with the eight touches (size 10). It is raised to 45m: the size counts criteria and touches, not a diagnosis, and I-0278 says the cause was never reproduced. Delivery is flai's 11:22Z, 11th in the pull order, moved by the 27m added.

Cost of delay. `flai cod` gives 5.00 USD a week: 2m lost per 168h cycle at 150 USD an hour. It stands. The single 2m occurrence is the operator's input, set by flai from I-0278.

After. None: the story reads nothing another open story writes.
