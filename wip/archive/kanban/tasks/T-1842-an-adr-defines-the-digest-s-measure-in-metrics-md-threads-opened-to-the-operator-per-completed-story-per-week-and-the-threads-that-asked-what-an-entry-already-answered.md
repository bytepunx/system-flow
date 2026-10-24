---
id: T-1842
type: task
nature: feature
title: "An ADR defines the digest's measure in metrics.md: threads opened to the operator per completed story per week, and the threads that asked what an entry already answered"
status: done
parent: S-0434
owner: alex
created: 2026-10-16T06:08:12Z
updated: 2026-10-24T21:41:31Z
transitions:
  - to: ready
    at: 2026-10-24T21:39:50Z
    by: system-flow
  - to: in-progress
    at: 2026-10-24T21:39:50Z
    by: system-flow
  - to: done
    at: 2026-10-24T21:41:31Z
    by: agent-S-0434
stream: S-0434
tags: [cli]
touches: [design/adrs, design/adrs/README.md, design/system/metrics.md]
usage:
  source: log
  seconds: 101
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 28
      output: 8883
      cache_read: 1487610
      cache_write: 73917
      cost: 0.894
---
# T-1842 An ADR defines the digest's measure in metrics.md: threads opened to the operator per completed story per week, and the threads that asked what an entry already answered

## Work

Write the ADR with `flai adr new`, refining ADR-0081 and ADR-0091, and add the definitions to `design/system/metrics.md` under a section of their own, with its chart under § Charts:

- A thread opened to the operator: opened by an author other than the operator (the manifest's `owner`) and other than flai's own notices (the upgrade thread, `permission_prompt`'s threads).
- The week a thread counts in: its first entry's ISO week.
- The stories completed in a week, as `waiting.weeks[]` counts them.
- `per_story`: threads over completed stories, absent with none.
- A repeat: a thread an entry lists in `asked_again` whose `date` is before the thread's first entry.
- `digest.entries`: the counts by state.
- The chart: threads per completed story per week, the repeats stacked.

First layer: the code and chart tasks build what this defines.

## Done when

The ADR is accepted and listed in `design/adrs/README.md`, `metrics.md` carries the definitions, and `flai check --strict` reports nothing on them.

## Notes

`metrics.md` is the contract between `flai stats` and the dashboard and changes only with an ADR (CLAUDE.md), so both land in this task.
