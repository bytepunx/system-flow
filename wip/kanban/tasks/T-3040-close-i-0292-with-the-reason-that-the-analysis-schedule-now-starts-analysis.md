---
id: T-3040
type: task
nature: improvement
title: Close I-0292 with the reason that the analysis schedule now starts analysis-cadence
status: backlog
parent: S-0692
owner: alex
created: 2026-11-10T10:34:10Z
updated: 2026-11-10T10:34:10Z
transitions: []
stream: S-0692
tags: [issues]
touches: [design/issues/I-0292-no-analysis-schedule-or-analysis-formula-is-set-so-the-analyzer-ran-once-in.md, design/issues/summary.md]
after: [T-3034, T-3037]
---
# T-3040 Close I-0292 with the reason that the analysis schedule now starts analysis-cadence

## Work

Waits for T-3034 and T-3037: the issue is closed on what fixed it, the schedule set in the manifest and the documents that say how it is set.

- In the story's worktree, run `flai issue close I-0292 --commit --reason "<reason>"`, or the MCP tool `issue_close` with `commit`, the reason naming the `analysis.schedule` set, `analysis.formula: analysis-cadence`, and that `flai manifest set` now writes the formula (T-3031).
- Tick the story's criteria you verified with `flai criteria tick`.

## Done when

- I-0292 is closed with that reason, and `design/issues/summary.md` no longer lists it.

## Notes
