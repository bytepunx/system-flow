---
id: T-1583
type: task
nature: improvement
title: The ADR index in design/adrs/README.md merges row by row from a stopped rebase's three versions, keeping each side's rows in number order
status: done
parent: S-0306
owner: alex
created: 2026-10-14T23:04:00Z
updated: 2026-10-14T23:07:28Z
transitions:
  - to: ready
    at: 2026-10-14T23:04:25Z
    by: agent-S-0306
  - to: in-progress
    at: 2026-10-14T23:04:26Z
    by: agent-S-0306
  - to: done
    at: 2026-10-14T23:07:28Z
    by: agent-S-0306
stream: S-0306
tags: []
touches: [flai/internal/adrindex/adrindex.go, flai/internal/adrindex/adrindex_test.go, design/adrs/0157-flai-stream-sync-and-flai-accept-merge-design-adrs-readme-md-s-index-row-by-row.md]
usage:
  source: log
  seconds: 179
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 29
      output: 189
      cache_read: 1114415
      cache_write: 63514
      cost: 0.5436
---
# T-1583 The ADR index in design/adrs/README.md merges row by row from a stopped rebase's three versions, keeping each side's rows in number order

## Work

A new leaf package `flai/internal/adrindex` (it imports no flai package that imports `check`, so `issues` can use it without a cycle: `adr` imports `check`, which imports `issues`). It has:

- `Path(r *workitem.Repo) string`: the index's path relative to the repository root, slash-separated as git names it, `<layout.design>/adrs/README.md`, as `issues.SummaryPath` does for the summary.
- `Merge(base, ours, theirs string) (string, error)`: the index merged three ways. Each version splits into the text above the table, the table's header and separator, its rows keyed by the file each links to, and the text below it. The text above and below and the header each take the side that changed it, or the one both agree on; when both changed one differently, it is an error naming the part, and the stop stays the agent's. Each row the same way, an absent row counting as a value, so a row both sides add is kept, a row one side changes is that side's, and a row one side removes is gone. The rows come out sorted by ADR number, then file name. A version with no table is an error.

Waits for nothing: it is pure, and the wiring (the next layer) waits for it.

## Done when

- Behavior tests cover: both sides append a row (both kept, in number order), one side notes a row (superseded by) while the other appends, both sides change one row differently (error), both change the text above the table differently (error), one side changes it (taken), no table (error).
- `flai test flai/internal/adrindex` passes.

## Notes
