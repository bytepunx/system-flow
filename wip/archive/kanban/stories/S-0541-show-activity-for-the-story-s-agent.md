---
id: S-0541
type: story
nature: improvement
title: "Show Activity for the Story’s Agent"
status: done
owner: alex
created: 2026-10-29T14:13:31Z
updated: 2026-10-29T15:12:10Z
transitions:
  - to: ready
    at: 2026-10-29T14:13:37Z
    by: alex
  - to: in-progress
    at: 2026-10-29T14:13:47Z
    by: agent-S-0541
  - to: review
    at: 2026-10-29T14:20:34Z
    by: agent-S-0541
  - to: done
    at: 2026-10-29T15:12:10Z
    by: orchestrator
tags: [dashboard]
touches: [flaiover/src, flaiover/src/lib/components/StoryAgent.svelte, flaiover/src/lib/components/StoryAgent.svelte.test.ts, design/system/flaiover-dashboard.md, docs/users/flaiover.md]
agent:
  harness: claude-code
  model: claude-sonnet-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 417
  priced_by: harness
  turns:
    - day: 2026-10-29
      ceremony: 1
      work: 27
  models:
    - model: claude-sonnet-5-5
      input: 62
      output: 9180
      cache_read: 3013621
      cache_write: 132828
      cost: 0.9246
  strategic:
    - kind: planner
      seconds: 147
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 54
          output: 12470
          cache_read: 3824687
          cache_write: 180104
          cost: 2.4554
    - kind: orchestrator
      seconds: 2360
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 206
          output: 3809
          cache_read: 70244192
          cache_write: 41183
          cost: 15.4992
        - model: claude-sonnet-5-5
          input: 6
          output: 52
          cache_read: 50830
          cache_write: 39170
          cost: 0.0621
cost_of_delay:
  inputs:
    time_lost_per_cycle: 5m
    by: alex
    at: 2026-10-29T14:13:31Z
  value: 12.5
  by: planner-S-0541
  at: 2026-10-29T14:15:58Z
forecast:
  duration: 30m
  delivery: 2026-10-29T14:45:00Z
  basis: "flai forecast's 3m (80 s per unit of size over 13 small improvement stories, size 2) raised to 30m: two tasks, a component change whose 22 tests mock api for every URL so AgentStream's reads must be answered, a docs task, and the close-out's vitest and svelte-check tiers; counted from the start at 14:14Z."
  by: planner-S-0541
  at: 2026-10-29T14:15:58Z
---
# S-0541 Show Activity for the Story’s Agent

## Goal

Currently, the only way to see the activity feed for a story’s agent is to go to the Activity page and scroll to it. The activity for a specific story should also be visible in the Story’s Agent sub-section in the same form as it’s presented on the Activity page today.

## Acceptance criteria
- [x] The agent panel’s activity component for a story is now visible from the Story’s Agent sub-section

## Tasks
- T-2331 Show the story agent's stream in the Agent panel
- T-2332 Say in the dashboard's design and user guide that a story's Agent pane shows its agent's stream

## Notes

### Planning

The "Story’s Agent sub-section" is the story page's Agent pane, `flaiover/src/lib/components/StoryAgent.svelte`, mounted in `flaiover/src/routes/items/[id]/+page.svelte`. The "activity component" is `AgentStream.svelte` (S-0142), which `ActivityView.svelte` renders on the Activity page for each story whose agent has started, open while it works or waits.

Touches:

| Touch | Source | Why |
|-------|--------|-----|
| `flaiover/src` | declared | The operator's touch, kept. A folder: the story's tasks narrow it in the claim to the files they name (ADR-0096), so it holds no other story beyond them. |
| `flaiover/src/lib/components/StoryAgent.svelte` | layout | Where the Agent pane is drawn; T-2331 mounts `AgentStream` in it. |
| `flaiover/src/lib/components/StoryAgent.svelte.test.ts` | layout, co-change | Its tests; S-0335 and S-0329 changed the two together. Its `api` mock answers every URL, so AgentStream's reads need answering there. |
| `design/system/flaiover-dashboard.md` | design, co-change | The `/items/:id` row and the "An agent's stream (S-0142)" bullet describe the pane and the stream; 20% of `flaiover/src` commits changed it. |
| `docs/users/flaiover.md` | design, co-change | The story page's Agent row; 16% of `flaiover/src` commits changed it. S-0536's `docs` folder touch overlaps it; S-0536 names no change to this file. |

Not predicted: `flaiover/src/routes/items/[id]/+page.svelte` and `AgentStream.svelte`, which need no change when `StoryAgent` passes the story, its run's start, and whether it is live; and `design/components/`, since no examples file exists for `StoryAgent` and its props stay the same. `flai touches suggest`'s flai paths (`flai/internal/hostapi`, `flai/cmd/serve_actions.go`) are left out: the stream's endpoint, `/api/agent-stream/[story]`, exists already.

Forecast: `flai forecast` gave 3m, 80 s per unit of size over 13 done small improvement stories times size 2. Raised to 30m, since the story now has two tasks and five touches, the component's 22 tests mock `api` for every URL, and the close-out runs vitest and svelte-check. Delivery 14:45Z, counted from the start at 14:14Z.

Cost of delay: `flai cod` gives 12.50 USD a week from the operator's input, 5m lost per 168h cycle at 150 USD an hour. The value stands as computed.

### Accepted by the orchestrator

- Verified: a46038a14cc198daa78b1e1bb9043a1daaef5368
- Batch: 20261029-150902
- Enqueued verified: 40084d55e4a5e3668625a8c07634915b4d7766af
- At: 2026-10-29T15:12:10Z

Verdict: S-0541 meets its criterion at 40084d55; flai verify passed every step (flaiover tier included), and all four changed files are under the story's touches.
- 1: flaiover/src/lib/components/StoryAgent.svelte, flaiover/src/lib/components/StoryAgent.svelte.test.ts, design/system/flaiover-dashboard.md, docs/users/flaiover.md
