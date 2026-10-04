---
id: T-0806
type: task
nature: improvement
title: flai stats reports strategic usage as its own figures, the mean cost per agent hour, and each item's expected cost
status: done
parent: S-0225
owner: alex
created: 2026-10-04T04:07:28Z
updated: 2026-10-04T04:21:10Z
transitions:
  - to: ready
    at: 2026-10-04T04:07:53Z
    by: agent-S-0225
  - to: in-progress
    at: 2026-10-04T04:12:44Z
    by: agent-S-0225
  - to: done
    at: 2026-10-04T04:21:10Z
    by: agent-S-0225
stream: S-0225
tags: []
touches: [flai/internal/metrics, flai/cmd/stats.go, design/system/metrics.md, docs/users/flai.md]
after: [T-0803]
usage:
  source: log
  seconds: 506
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 92
      output: 644
      cache_read: 5914883
      cache_write: 145627
      cost: 2.472
---
# T-0806 flai stats reports strategic usage as its own figures, the mean cost per agent hour, and each item's expected cost

## Work

`flai stats` reports strategic usage as figures of its own, apart from the agents': per item, `items[].usage.strategic`, one entry per kind with `tokens`, `cost`, `seconds`, and `estimated`; in totals, `usage.strategic` over the items of the report's type completed in the window; and in spend, the same under each type and each bucket, so that the dashboard can draw it as its own series. An item that carries only strategic usage is left out of every agent aggregate and per-model figure.

It also reports the project's mean cost per agent hour, the agents' cost over their hours across the stories measured from their logs, archived ones included, and for each item with a `forecast.duration`, or else an `estimate`, its expected cost at that rate, marked as an estimate. `flai stats` prints a strategic line with the totals. `design/system/metrics.md` defines every new value and `docs/users/flai.md` mentions them.

Waits for the usage task. Runs beside the serve task: no path in common.

## Done when

- [ ] `flai stats --json` carries `items[].usage.strategic`, `usage.strategic`, the strategic figures under `usage.spend`, `usage.cost_per_agent_hour`, and `items[].expected_cost`, as `metrics.md` defines them
- [ ] Per-model figures and agent totals are unchanged by strategic usage, with a test
- [ ] Tests pin the mean cost per agent hour and the expected cost from a forecast and from an estimate
- [ ] `go test ./internal/metrics/... ./cmd/...` passes
- [ ] `design/system/metrics.md` and `docs/users/flai.md` describe them

## Notes
