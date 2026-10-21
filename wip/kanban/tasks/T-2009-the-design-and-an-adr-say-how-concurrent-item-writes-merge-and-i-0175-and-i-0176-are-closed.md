---
id: T-2009
type: task
nature: remediation
title: The design and an ADR say how concurrent item writes merge, and I-0175 and I-0176 are closed
status: done
parent: S-0479
owner: alex
created: 2026-10-21T12:13:03Z
updated: 2026-10-21T12:48:08Z
transitions:
  - to: ready
    at: 2026-10-21T12:47:07Z
    by: system-flow
  - to: in-progress
    at: 2026-10-21T12:47:08Z
    by: system-flow
  - to: done
    at: 2026-10-21T12:48:08Z
    by: agent-S-0479
stream: S-0479
tags: []
touches: [design/system/work-hierarchy.md, design/adrs, design/issues/I-0175-a-story-s-in-progress-transition-from-flai-story-start-was-lost-when-a-concurrent-write-rewrote-the-item-file-leaving-it-ready-while-its-agent-worked-it.md, design/issues/summary.md, design/issues/I-0176-a-usage-write-to-a-story-overwrites-its-move-to-in-progress-made-meanwhile-so-a-started-story-reads-as-ready-again.md]
after: [T-2008]
---
# T-2009 The design and an ADR say how concurrent item writes merge, and I-0175 and I-0176 are closed

## Work

Waits for T-2008, so that the documents describe what was built. Record the decision as an ADR with `flai adr new --commit`: item writes go through one lock and merge a stale copy by top-level key, refusing a key both changed. Say in `design/system/work-hierarchy.md` how concurrent writes to an item behave, linking the ADR. Close I-0175 with `flai issue close I-0175 --reason` naming the fix, and I-0176, the same lost update from a usage write, which the operator merged into this story on TH-0577.

## Done when

- The ADR is in `design/adrs` and its index, and `design/system/work-hierarchy.md` links it.
- I-0175 and I-0176 are closed and `design/issues/summary.md` no longer lists them.

## Notes
