---
id: T-2610
type: task
nature: research
title: "Record the exploration's findings in server-performance.md: the cost of each slow read and write at 596 stories, the serial I/O that could run in parallel, and what would need a change to how flai keeps state"
status: backlog
parent: S-0596
owner: alex
created: 2026-11-07T22:10:56Z
updated: 2026-11-07T22:12:56Z
transitions: []
stream: S-0596
tags: [cli]
touches: [design/system/server-performance.md]
after: [T-2607, T-2608, T-2609]
---
# T-2610 Record the exploration's findings in server-performance.md: the cost of each slow read and write at 596 stories, the serial I/O that could run in parallel, and what would need a change to how flai keeps state

## Work

Add a section to `design/system/server-performance.md` from the findings of T-2607, T-2608, and T-2609, which it waits for: what was measured on this repository and in the operator's serve log, where the time goes for the board, the writes, the charts, the close-out, and acceptance, which I/O runs serially where it could run in parallel, which causes are actionable within how flai keeps its state today, and which would need that to change. It is the only task that writes the document, so the three exploration tasks can run at once.

## Done when

- [ ] `design/system/server-performance.md` holds the section, with measurements, paths, and the actionable causes listed apart from those that need a change to how flai keeps state.

## Notes
