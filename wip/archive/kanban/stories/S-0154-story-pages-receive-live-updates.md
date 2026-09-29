---
id: S-0154
type: story
nature: feature
title: Story pages receive live updates
status: done
parent: E-0013
owner: alex
created: 2026-09-29T05:55:11Z
updated: 2026-09-29T19:04:33Z
transitions:
  - to: ready
    at: 2026-09-29T06:04:14Z
    by: alex
  - to: in-progress
    at: 2026-09-29T07:06:53Z
    by: agent-S-0154
  - to: review
    at: 2026-09-29T07:21:41Z
    by: agent-S-0154
  - to: done
    at: 2026-09-29T19:04:33Z
    by: alex
tags: [dashboard, cli]
topics: [client-side, server-side]
touches: [flaiover/src, flai/cmd, flai/internal/serve, design/system/dashboard-host-channel.md, design/system/flaiover-dashboard.md, docs/users/flaiover.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 925
  models:
    - model: claude-opus-5-5
      input: 228
      output: 64112
      cache_read: 16357976
      cache_write: 213666
      cost: 6.2641
---
# S-0154 Story pages receive live updates

## Goal

Story pages currently go stale once they're opened. This means that the operator must refresh them to catch updates on agent status or updates to threads once they respond.

## Acceptance criteria
- [x] The story gets updates as they occur and the page updates to reflect the changes
- [x] The active thread pane needs to show that the agent is working on processing the operator response with a simple animation of some kind.

## Tasks
- T-0553 flai serve tells the dashboard when a story's agent starts or ends
- T-0554 The story page reloads its item, threads, and agent panel on every change, over one shared event stream
- T-0555 The thread pane shows the agent working on the operator's reply with an animation
- T-0556 Document live story pages and the agent notification; all tiers and lint pass

## Notes

- Verified by behaviour tests (`flaiover/src/routes/items/[id]/item.svelte.test.ts`, `Threads.svelte.test.ts`, `StoryAgent.svelte.test.ts`, `events.test.ts`, the server-side `agent`, `sse`, and `project-events` tests, and `flai/internal/serve/serve_test.go`), not in a browser against a running dashboard: the running dashboard and `flai serve` are the installed releases, not this branch.
- The animation is three bouncing dots in the working-dot green on the agent's side of the thread, with "<agent> is working on your reply", shown while the agent runs or is being started again on that thread; it is still for a reader who asks for less motion.
- An agent starting or ending changes no file, so `flai serve` now sends the channel notification `agent`; with an older flai the page falls back to the 15 s re-ask StoryAgent already had.
