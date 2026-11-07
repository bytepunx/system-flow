---
id: T-2611
type: task
nature: research
title: Draft a performance epic with one draft story per actionable cause, each with the cause's numbers and a proposed fix, and link them from server-performance.md
status: done
parent: S-0596
owner: alex
created: 2026-11-07T22:11:23Z
updated: 2026-11-07T22:43:24Z
transitions:
  - to: ready
    at: 2026-11-07T22:38:35Z
    by: agent-S-0596
  - to: in-progress
    at: 2026-11-07T22:38:36Z
    by: agent-S-0596
  - to: done
    at: 2026-11-07T22:43:24Z
    by: agent-S-0596
stream: S-0596
tags: [cli]
touches: [design/system/server-performance.md]
after: [T-2610]
usage:
  source: log
  seconds: 288
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 107
      output: 31082
      cache_read: 7031522
      cache_write: 133386
      cost: 2.863
---
# T-2611 Draft a performance epic with one draft story per actionable cause, each with the cause's numbers and a proposed fix, and link them from server-performance.md

## Work

Make an epic with `flai epic new` and, under it, one draft story per actionable cause T-2610 records, each with a goal carrying the cause's measurement, acceptance criteria that say what is measured after, and a proposed fix in its notes. Make none for a cause that would need a change to how flai keeps state. Then link the epic and its stories from the section's story table in `design/system/server-performance.md`. Waits for T-2610, whose section names the causes and whose table it fills in.

## Done when

- [ ] The epic and its draft stories exist, pass `flai check`, and the section's table links each.

## Notes
