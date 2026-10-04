---
id: T-0808
type: task
nature: improvement
title: The cost charts draw strategic cost as its own series
status: ready
parent: S-0225
owner: alex
created: 2026-10-04T04:07:30Z
updated: 2026-10-04T04:07:54Z
transitions:
  - to: ready
    at: 2026-10-04T04:07:54Z
    by: agent-S-0225
stream: S-0225
tags: []
touches: [flaiover/src/lib/viz, design/system/metrics.md, design/system/flaiover-dashboard.md, docs/users/flaiover.md]
after: [T-0806, T-0807]
---
# T-0808 The cost charts draw strategic cost as its own series

## Work

The dashboard's cost charts draw strategic cost as its own series: `$ / Item` adds a strategic segment to each item's bar, beside the models' segments, and `$ / bucket` adds a strategic series to each bucket's stack; the per-model charts (`Avg. Cost / Model`, `Tokens / $`, and the rest) stay the agents' alone. The charts table in `design/system/metrics.md`, `design/system/flaiover-dashboard.md`, and `docs/users/flaiover.md` say so.

Waits for the item page task, which edits the same dashboard documents, and for the stats task, whose JSON it reads.

## Done when

- [ ] `$ / Item` and `$ / bucket` show a strategic series when items carry strategic usage, marked as estimated, and none otherwise, with tests
- [ ] The per-model charts are unchanged by strategic usage, with a test
- [ ] The flaiover fixture test still agrees with `flai stats --json`
- [ ] The documents say so, and the dashboard's tests and check pass

## Notes
