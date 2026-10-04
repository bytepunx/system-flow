---
id: T-0807
type: task
nature: improvement
title: flai show and the item page show the expected cost and the strategic usage
status: in-progress
parent: S-0225
owner: alex
created: 2026-10-04T04:07:29Z
updated: 2026-10-04T04:21:11Z
transitions:
  - to: ready
    at: 2026-10-04T04:07:53Z
    by: agent-S-0225
  - to: in-progress
    at: 2026-10-04T04:21:11Z
    by: agent-S-0225
stream: S-0225
tags: []
touches: [flai/cmd/show.go, flaiover/src/routes/items, flaiover/src/lib/usage.ts, docs/users/flai.md, docs/users/flaiover.md, design/system/flaiover-dashboard.md]
after: [T-0806]
---
# T-0807 flai show and the item page show the expected cost and the strategic usage

## Work

`flai show`, in text and `--json`, and the dashboard's item page show the expected cost beside the forecast or the estimate, marked as an estimate, and an item's strategic usage beside its agents' usage, per kind. The figures come from what the stats task defines; the page reads them through the API it already uses for the item.

Waits for the stats task, whose expected cost and strategic figures it shows, and whose `docs/users/flai.md` it also edits.

## Done when

- [ ] `flai show` prints the expected cost beside the forecast or estimate, and the strategic usage per kind beside the agents' usage, in text and JSON, with tests
- [ ] The item page shows both, the expected cost marked as an estimate, with a test where the page has them
- [ ] `docs/users/flai.md`, `docs/users/flaiover.md`, and `design/system/flaiover-dashboard.md` say so
- [ ] The tests for what changed pass

## Notes
