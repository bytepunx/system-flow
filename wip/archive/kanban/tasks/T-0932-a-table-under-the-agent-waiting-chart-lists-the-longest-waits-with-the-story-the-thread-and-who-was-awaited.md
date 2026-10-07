---
id: T-0932
type: task
nature: feature
title: A table under the agent-waiting chart lists the longest waits with the story, the thread, and who was awaited
status: done
parent: S-0215
owner: alex
created: 2026-10-05T05:45:14Z
updated: 2026-10-07T09:05:35Z
transitions:
  - to: ready
    at: 2026-10-07T08:54:32Z
    by: agent-S-0215
  - to: in-progress
    at: 2026-10-07T09:02:04Z
    by: agent-S-0215
  - to: done
    at: 2026-10-07T09:05:35Z
    by: agent-S-0215
stream: S-0215
tags: [dashboard]
touches: [flaiover/src/lib/components/WaitTable.svelte, flaiover/src/lib/components/WaitTable.svelte.test.ts, flaiover/src/routes/charts]
after: [T-0923, T-0927]
usage:
  source: log
  seconds: 211
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 44
      output: 15026
      cache_read: 2008880
      cache_write: 70319
      cost: 1.1272
---
# T-0932 A table under the agent-waiting chart lists the longest waits with the story, the thread, and who was awaited

## Work

- Write `flaiover/src/lib/components/WaitTable.svelte` on the pattern of `SpendTable.svelte`. Its rows are `waiting.longest[]`: the story, linked to `/items/<id>`; the kind; the thread, linked to it, or `review`; started and ended, or `open`; the hours waited; and who was awaited, or `—` while the wait is open.
- Render it under the chart in `flaiover/src/routes/charts/[kind]/+page.svelte` when the kind is `agent-waiting`, and say so when the window holds no wait.
- Check the page against `flai stats --json` on this repository over the same window: the weeks' hours and the table's rows match to the second.

It waits for T-0927 for the `waiting` types and the chart kind, and for T-0923 so that `/api/stats` carries `longest[]` to check against.

## Done when

- The table renders under `/charts/agent-waiting` and nowhere else.
- Tests in `WaitTable.svelte.test.ts` cover a thread row, a review row, an open wait, and an empty window.
- The page's figures match `flai stats --json` for the same window.
- The dashboard's lint, type check, and unit tests pass.

## Notes
