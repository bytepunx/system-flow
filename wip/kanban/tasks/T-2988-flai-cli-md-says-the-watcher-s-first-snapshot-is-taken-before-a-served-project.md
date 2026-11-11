---
id: T-2988
type: task
nature: remediation
title: flai-cli.md says the watcher's first snapshot is taken before a served project's first looks
status: done
parent: S-0690
owner: alex
created: 2026-11-10T10:29:50Z
updated: 2026-11-11T02:33:19Z
transitions:
  - to: ready
    at: 2026-11-11T02:32:56Z
    by: agent-S-0690
  - to: in-progress
    at: 2026-11-11T02:33:00Z
    by: agent-S-0690
  - to: done
    at: 2026-11-11T02:33:19Z
    by: agent-S-0690
stream: S-0690
tags: [docs, serve]
touches: [design/system/flai-cli.md]
after: [T-2984]
usage:
  source: log
  seconds: 19
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 5
      output: 1640
      cache_read: 491315
      cache_write: 14036
      cost: 0.2434
---
# T-2988 flai-cli.md says the watcher's first snapshot is taken before a served project's first looks

## Work

It waits for T-2984, whose order of the first snapshot, the item snapshot, and the first looks it describes.

In `design/system/flai-cli.md`, the paragraph on the item snapshot (S-0606) says that serve marks the root stale once a minute "for a change the watcher never reports, such as one made before its first look". Rewrite that sentence: the watcher takes its first snapshot before the item snapshot and the project's first looks (S-0690), so a change made from then on is reported within a tick. The minute's refresh stays for a change no watched path shows. Bump the document's `updated` date.

## Done when

- `design/system/flai-cli.md` says when the watcher's first snapshot is taken, and no longer names a change before its first look as one the watcher misses.
- `flai test design/system/flai-cli.md` passes.

## Notes
