---
id: T-1995
type: task
nature: remediation
title: stats.get and flai stats are measured again on this repository, server-performance.md records the regression and its fix, and I-0155 is closed
status: backlog
parent: S-0471
owner: alex
created: 2026-10-21T12:08:25Z
updated: 2026-10-21T12:08:25Z
transitions: []
stream: S-0471
tags: [flai, docs]
touches: [design/system/server-performance.md, design/issues/I-0155-the-chart-data-takes-about-4-s-warm-against-10-ms-measured-on-2026-09-29-almost-all-of-it-in-metrics-replayholds.md, design/issues/summary.md]
after: [T-1993]
---
# T-1995 stats.get and flai stats are measured again on this repository, server-performance.md records the regression and its fix, and I-0155 is closed

## Work

Show the fix on the real repository and close the issue. It waits for T-1993 because it measures the fixed replay.

- Measure `stats.get` with `since=30d`, `type=story`, `bucket=day` cold and warm, and `flai stats --json`, on this repository with the tree's `bin/flai`, as `design/system/server-performance.md` § How it is measured and S-0427's baseline did. Never time a figure you did not observe.
- In `design/system/server-performance.md`, add a section after `## After the stories` that records the regression: 9.1 to 10.6 ms warm on 2026-09-29, about 4,022 ms warm on 2026-10-16 with 67 percent under `metrics.replayHolds`, its cause, the fix, and the figures measured now. Add S-0471 to `## Stories` and bump `updated`.
- Close the issue with `flai issue close I-0155 --reason "<what fixed it>"`, which writes the issue and `design/issues/summary.md`.
- Tick the story's criteria that the tests and the measurement verify, with `flai criteria tick`.

## Done when

- `design/system/server-performance.md` holds the before and after figures, and the markdown lint passes on it.
- I-0155 is closed with its reason, and `flai check --strict` is clean.

## Notes
