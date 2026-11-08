---
id: T-2687
type: task
nature: improvement
title: Measure a warm stats.get at 30 and 365 days before and after on this repository, and record them with the kept replay in the design
status: backlog
parent: S-0608
owner: alex
created: 2026-11-07T23:51:26Z
updated: 2026-11-07T23:51:26Z
transitions: []
stream: S-0608
tags: [cli]
touches: [design/system/server-performance.md, design/system/flai-cli.md]
after: [T-2672, T-2675, T-2683]
---
# T-2687 Measure a warm stats.get at 30 and 365 days before and after on this repository, and record them with the kept replay in the design

## Work

Measure a warm `stats.get` on this repository at windows of 1, 7, 30, 90, and 365 days, with `flai hostapi --timing stats.get` as `design/system/server-performance.md` measured cause 18, on the main branch's flai for before and the story branch's for after. Record the before and after per window under cause 18 in `design/system/server-performance.md`, with the kept replay's own time. Say in `design/system/flai-cli.md`, where its Internal structure describes the item store and the stats reads, that the store counts a generation and that flai serve keeps the hold replay and the drift until the generation or the commits' key moves.

Waits for T-2672, T-2675, and T-2683, so that it measures every change of the story.

## Done when

- A warm 365-day `stats.get` takes under 1.2 times a warm 30-day one, and `server-performance.md` records both, before and after.
- `flai-cli.md` describes the generation and the kept replay.
- `flai check --strict` and the markdown lint pass on both files.

## Notes
