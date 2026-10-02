---
id: T-0705
type: task
nature: remediation
title: The activity page keeps what it shows when a reload fails or flai is away
status: done
parent: S-0178
owner: arobson
created: 2026-10-02T16:46:44Z
updated: 2026-10-02T16:50:51Z
transitions:
  - to: ready
    at: 2026-10-02T16:47:08Z
    by: agent-S-0178
  - to: in-progress
    at: 2026-10-02T16:47:08Z
    by: agent-S-0178
  - to: done
    at: 2026-10-02T16:50:51Z
    by: agent-S-0178
stream: S-0178
tags: []
touches: [flaiover/src/routes/activity]
usage:
  source: log
  seconds: 223
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 78
      output: 31845
      cache_read: 2665332
      cache_write: 114981
      cost: 1.8918
---
# T-0705 The activity page keeps what it shows when a reload fails or flai is away

## Work

`routes/activity/+page.svelte` replaces the whole list with an error paragraph when a reload of `/api/activity` fails, so the page collapses to the top and stays there after the next good load. Keep the last streams shown and put the error above them; the bare error stays only when nothing was ever loaded.

`/api/host-agent` answers `{ enabled: false }`, with no `state`, when no flai is connected; the page then passes no agents, every stream window unmounts, and they all come back from the tail when flai returns. When an answer has no `state`, keep the agents last known.

Waits for nothing: no other task touches the route.

## Done when

- A test shows a failed reload after a good one keeps the cards and shows the error; it fails without the change.
- A test shows an answer without `state` after one with it keeps the stream windows.
- The route's tests pass.

## Notes
