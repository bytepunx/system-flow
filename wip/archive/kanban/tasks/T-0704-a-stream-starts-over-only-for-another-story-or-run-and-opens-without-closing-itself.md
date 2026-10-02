---
id: T-0704
type: task
nature: remediation
title: A stream starts over only for another story or run, and opens without closing itself
status: done
parent: S-0178
owner: arobson
created: 2026-10-02T16:46:44Z
updated: 2026-10-02T16:50:50Z
transitions:
  - to: ready
    at: 2026-10-02T16:47:08Z
    by: agent-S-0178
  - to: in-progress
    at: 2026-10-02T16:47:08Z
    by: agent-S-0178
  - to: done
    at: 2026-10-02T16:50:50Z
    by: agent-S-0178
stream: S-0178
tags: []
touches: [flaiover/src/lib/components/AgentStream.svelte, flaiover/src/lib/components/AgentStream.svelte.test.ts]
usage:
  source: log
  seconds: 222
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 78
      output: 31845
      cache_read: 2665332
      cache_write: 114981
      cost: 1.8918
---
# T-0704 A stream starts over only for another story or run, and opens without closing itself

## Work

`AgentStream.svelte`'s effect re-runs whenever the activity page reloads its streams or its agents, because it tracks the `story` and `started` prop getters, which are rebuilt from new objects on every load (`void started`, and `void step()` reading `story` synchronously through `read()` and `api()`). Each re-run empties the box and re-reads the tail, so the page shortens and the browser pulls the window up, and the box is put back at its end even when the reader had scrolled away. Make the effect depend only on the values of `story` and `started` (a `$derived` of each, the rest under `untrack`), so equal values given again change nothing.

The `<details>` follows `open={live(a)}`, so it closes by itself when the agent ends, which shortens the page under the reader too. Open it when the `open` prop becomes true, and never close it by itself: closing is the reader's.

Waits for nothing: no other task touches these files.

## Done when

- A test passes the same `story` and `started` again through new objects (as ActivityView does after a reload) and finds no read from the tail and the entries kept; it fails without the change.
- A test shows a box scrolled away from its end is not moved by an append, and one at its end follows it.
- A test shows the stream opening when `open` becomes true and staying open when it becomes false.
- The component's tests pass.

## Notes
