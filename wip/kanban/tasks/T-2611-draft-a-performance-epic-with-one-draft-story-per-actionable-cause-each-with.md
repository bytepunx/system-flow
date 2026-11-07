---
id: T-2611
type: task
nature: research
title: Draft a performance epic with one draft story per actionable cause, each with the cause's numbers and a proposed fix, and link them from server-performance.md
status: backlog
parent: S-0596
owner: alex
created: 2026-11-07T22:11:23Z
updated: 2026-11-07T22:13:08Z
transitions: []
stream: S-0596
tags: [cli]
touches: [design/system/server-performance.md]
after: [T-2610]
---
# T-2611 Draft a performance epic with one draft story per actionable cause, each with the cause's numbers and a proposed fix, and link them from server-performance.md

## Work

Make an epic with `flai epic new` and, under it, one draft story per actionable cause T-2610 records, each with a goal carrying the cause's measurement, acceptance criteria that say what is measured after, and a proposed fix in its notes. Make none for a cause that would need a change to how flai keeps state. Then link the epic and its stories from the section's story table in `design/system/server-performance.md`. Waits for T-2610, whose section names the causes and whose table it fills in.

## Done when

- [ ] The epic and its draft stories exist, pass `flai check`, and the section's table links each.

## Notes
