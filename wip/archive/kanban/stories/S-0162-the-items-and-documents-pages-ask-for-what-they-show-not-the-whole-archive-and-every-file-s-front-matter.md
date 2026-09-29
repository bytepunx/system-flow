---
id: S-0162
type: story
nature: improvement
title: The items and documents pages ask for what they show, not the whole archive and every file's front matter
status: done
parent: E-0012
owner: alex
created: 2026-09-29T07:00:31Z
updated: 2026-09-29T20:28:44Z
transitions:
  - to: ready
    at: 2026-09-29T19:20:21Z
    by: alex
  - to: in-progress
    at: 2026-09-29T20:13:22Z
    by: agent-S-0162
  - to: review
    at: 2026-09-29T20:27:58Z
    by: agent-S-0162
  - to: done
    at: 2026-09-29T20:28:44Z
    by: alex
tags: [dashboard]
topics: [server-side, back-end]
touches: [flaiover/src, flai/internal/hostapi, design/issues/I-0027-work-items-written-in-the-main-checkout-reach-ci-without-a-markdown-lint.md, design/issues/summary.md, design/system/flai-cli.md, design/system/flaiover-dashboard.md, design/system/server-performance.md, docs/operators/index.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 909
  models:
    - model: claude-opus-5-5
      input: 198
      output: 53185
      cache_read: 14972728
      cache_write: 199387
      cost: 5.6541
---
# S-0162 The items and documents pages ask for what they show, not the whole archive and every file's front matter

## Goal

Two answers are large and rebuilt on every ask: `items.list` with the archive and bodies is 1.28 MB, and `docs.tree` with the front matter of every file under `design`, `docs`, and `wip` is 772 KB after a 127 ms walk. Cause 7 of `design/system/server-performance.md`.

## Acceptance criteria
- [x] The items page's first answer is under 200 KB on this repository: bodies are asked for when an item is opened, and the archive when it is shown.
- [x] `docs.tree` is under 200 KB, or `flai serve` keeps it until a file under the three folders changes, and a warm answer takes under 10 ms.
- [x] The readiness probe, `/_ready`, asks for no more than the manifest to know flai answers.

## Tasks
- T-0576 docs.tree answers titles without front matter and keeps them while files are unchanged
- T-0577 items.count answers how many items are active and archived
- T-0578 The dashboard asks for items without bodies and the archive only where it is shown
- T-0579 Record cause 7's remeasurement and the changed contracts

## Notes

Measured by S-0152: `items.list` 122 ms and 1,278,709 bytes; `docs.tree` 135 ms and 772,454 bytes. `/_ready` asks `repo().items()`, the whole archive with bodies.

Verified by S-0162 on 2026-09-29, in one process through the method table `flai serve` answers with (`design/system/server-performance.md`, cause 7):

- The overview's first answers are `items.list` for the active items without bodies, 8.4 KB, and `items.count`, 28 B, where it asked 1.36 MB. The charts ask for the epics with the archive, 8.4 KB. A body is asked for with `item.get` when its item is opened.
- `docs.tree` is 320 KB without front matter, still over 200 KB, so flai keeps each file's title while the file is unchanged: warm 6.1 to 8.2 ms over eleven calls, cold 129 to 131 ms.
- `/_ready` asks `project.info` alone; `src/routes/_ready/ready.test.ts` pins it.
- Not measured in the running dashboard: that needs an image built from this story in place of the operator's container.
