---
id: T-2909
type: task
nature: improvement
title: The Activity page offers Retry for a stopped or failed agent whenever the card's menu does
status: done
parent: S-0669
owner: alex
created: 2026-11-10T08:47:44Z
updated: 2026-11-10T09:53:21Z
transitions:
  - to: ready
    at: 2026-11-10T09:48:30Z
    by: agent-S-0669
  - to: in-progress
    at: 2026-11-10T09:48:31Z
    by: agent-S-0669
  - to: done
    at: 2026-11-10T09:53:21Z
    by: agent-S-0669
stream: S-0669
tags: [flaiover, agent]
touches: [flaiover/src/lib/components/ActivityView.svelte, flaiover/src/routes/activity/+page.svelte, flaiover/src/lib/activity.test.ts, design/system/flaiover-dashboard.md, flaiover/src/lib/components/ActivityView.svelte.test.ts]
after: [T-2907]
usage:
  source: log
  seconds: 290
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 30
      output: 249
      cache_read: 3954011
      cache_write: 26146
      cost: 1.4649
---
# T-2909 The Activity page offers Retry for a stopped or failed agent whenever the card's menu does

## Work

This task waits for T-2907: it builds the fix proposed there, against the test written there. Keep `retryable()` in `flaiover/src/lib/activity.ts` the one rule for the story page, the card's menu (S-0202), and now the Activity page.

- `ActivityView.svelte` takes a `writable` prop. A narrated card whose agent `retryable(a, canStop, s.status, writable)` shows **Retry** on its activity line, where Stop stands for a running agent. Retry posts `{ action: 'restart' }` to `/api/items/<id>/agent`. It hides once pressed for that run, as `StoryAgent` does, until flai refuses. A refusal shows flai's reason under the line. Then it tells the page with `onstopped`'s sibling so the page asks `/api/host-agent` again.
- `routes/activity/+page.svelte` asks `/api/board` for `writable`, once, as the item page does, and passes it.
- `activity.test.ts` gains a case: a stopped run of a story in progress (`outcome: 'stopped'`, `stopped` set, exit 143) gets Retry from both `retryable` and `agentAction`.
- `design/system/flaiover-dashboard.md`: the "Stopping an agent (S-0170)" paragraph says the Activity page offers Retry for a failed or stopped agent under the same rule.

## Done when

- [ ] The test T-2907 wrote passes.
- [ ] `activity.test.ts` checks that a stopped run of a story in progress gets Retry from both `retryable` and `agentAction`.
- [ ] `flai test` passes on the paths the task changed.

## Notes
