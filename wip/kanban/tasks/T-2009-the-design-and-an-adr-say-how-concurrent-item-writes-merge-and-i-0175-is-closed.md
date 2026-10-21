---
id: T-2009
type: task
nature: remediation
title: The design and an ADR say how concurrent item writes merge, and I-0175 is closed
status: backlog
parent: S-0479
owner: alex
created: 2026-10-21T12:13:03Z
updated: 2026-10-21T12:13:33Z
transitions: []
stream: S-0479
tags: []
touches: [design/system/work-hierarchy.md, design/adrs, design/issues/I-0175-a-story-s-in-progress-transition-from-flai-story-start-was-lost-when-a-concurrent-write-rewrote-the-item-file-leaving-it-ready-while-its-agent-worked-it.md, design/issues/summary.md]
after: [T-2008]
---
# T-2009 The design and an ADR say how concurrent item writes merge, and I-0175 is closed

## Work

Waits for T-2008, so that the documents describe what was built. Record the decision as an ADR with `flai adr new --commit`: item writes go through one lock and merge a stale copy by top-level key, refusing a key both changed. Say in `design/system/work-hierarchy.md` how concurrent writes to an item behave, linking the ADR. Close I-0175 with `flai issue close I-0175 --reason` naming the fix.

## Done when

- The ADR is in `design/adrs` and its index, and `design/system/work-hierarchy.md` links it.
- I-0175 is closed and `design/issues/summary.md` no longer lists it.

## Notes
