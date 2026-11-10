---
id: T-2907
type: task
nature: improvement
title: A test reproduces the Activity page's missing Retry after a stop, and the narrative names its cause and the fix proposed
status: done
parent: S-0669
owner: alex
created: 2026-11-10T08:47:33Z
updated: 2026-11-10T09:48:03Z
transitions:
  - to: ready
    at: 2026-11-10T09:44:39Z
    by: agent-S-0669
  - to: in-progress
    at: 2026-11-10T09:44:40Z
    by: agent-S-0669
  - to: done
    at: 2026-11-10T09:48:03Z
    by: agent-S-0669
stream: S-0669
tags: [flaiover, agent, test]
touches: [flaiover/src/lib/components/ActivityView.svelte.test.ts]
usage:
  source: log
  seconds: 203
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 18
      output: 184
      cache_read: 2119711
      cache_write: 12366
      cost: 0.7847
---
# T-2907 A test reproduces the Activity page's missing Retry after a stop, and the narrative names its cause and the fix proposed

## Work

This task waits for none: it pins the cause the fix builds on.

The story's agent found the cause before writing this task's test (narrative `## Decisions`). The operator pressed Stop on S-0556's Activity page, as TH-0836 asked, and reported "the story's agent page" from there. `ActivityView.svelte` draws Stop under `stoppable()` and never draws Retry. After the stop its card reads "agent stopped by the operator" with nothing to press. The story page's `StoryAgent` and the card's menu both apply `retryable()` and offer Retry.

Add a case to `ActivityView.svelte.test.ts`, in the "stopping an agent (S-0170)" block. Give a stream for a story in progress whose agent is `failed`, `outcome: 'stopped'`, `stopped` set, exit 143. With the agent action on and the page writable, Retry is shown on its activity line. Pressing it posts `{ action: 'restart' }` to `/api/items/<id>/agent`. Give `ActivityView` the prop `writable` the fix adds. The case must fail on main, since there is no Retry there.

## Done when

- [ ] A test in `ActivityView.svelte.test.ts` fails on main because the Activity page offers no Retry for a stopped run of a story in progress.
- [ ] The narrative's `## Decisions` names the cause, the evidence for it, and the fix proposed.

## Notes
