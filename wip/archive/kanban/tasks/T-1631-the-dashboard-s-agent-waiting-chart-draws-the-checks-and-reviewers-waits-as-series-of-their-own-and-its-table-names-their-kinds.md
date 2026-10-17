---
id: T-1631
type: task
nature: feature
title: The dashboard's agent-waiting chart draws the checks and reviewers waits as series of their own, and its table names their kinds
status: done
parent: S-0396
owner: alex
created: 2026-10-16T04:39:14Z
updated: 2026-10-17T17:54:25Z
transitions:
  - to: ready
    at: 2026-10-17T17:36:08Z
    by: agent-S-0396
  - to: in-progress
    at: 2026-10-17T17:47:15Z
    by: agent-S-0396
  - to: done
    at: 2026-10-17T17:54:25Z
    by: agent-S-0396
stream: S-0396
tags: [dashboard]
touches: [flaiover/src/lib/viz/charts.ts, flaiover/src/lib/viz/charts.test.ts, flaiover/src/lib/components/WaitTable.svelte, flaiover/src/lib/components/WaitTable.svelte.test.ts]
after: [T-1629]
usage:
  source: log
  seconds: 430
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 53
      output: 21610
      cache_read: 2536042
      cache_write: 99303
      cost: 1.543
---
# T-1631 The dashboard's agent-waiting chart draws the checks and reviewers waits as series of their own, and its table names their kinds

## Work

- In `flaiover/src/lib/viz/charts.ts`, extend the `WaitWeek` and `LongestWait` types with `checks`, `reviewers`, and the two kinds. Add two series to `agentWaiting()` beside threads and review, coloured from the palette and the design system's tokens.
- Leave a week without them as it draws today.
- In `WaitTable.svelte`, name the kinds `checks` and `reviewers`, and link the story.
- Extend `charts.test.ts` and `WaitTable.svelte.test.ts`.
- Waits for T-1629: it draws that task's JSON.

## Done when

The chart draws the two series from a fixture with them and is unchanged without them, the table names the kinds, and the flaiover tests pass.

## Notes
