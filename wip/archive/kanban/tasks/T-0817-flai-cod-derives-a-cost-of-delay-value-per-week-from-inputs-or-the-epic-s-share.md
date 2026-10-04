---
id: T-0817
type: task
nature: feature
title: flai cod derives a cost of delay value per week from inputs or the epic's share
status: done
parent: S-0210
owner: alex
created: 2026-10-04T19:45:59Z
updated: 2026-10-04T20:06:51Z
transitions:
  - to: ready
    at: 2026-10-04T19:59:36Z
    by: agent-S-0210
  - to: in-progress
    at: 2026-10-04T19:59:36Z
    by: agent-S-0210
  - to: done
    at: 2026-10-04T20:06:51Z
    by: agent-S-0210
stream: S-0210
tags: [flai]
touches: [flai/internal/planning/cod.go, flai/internal/planning/cod_test.go, flai/cmd/cod.go, flai/cmd/cod_test.go, flai/cmd/root.go]
after: [T-0816]
usage:
  source: log
  seconds: 435
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 45
      output: 17163
      cache_read: 2231886
      cache_write: 63855
      cost: 1.1667
---
# T-0817 flai cod derives a cost of delay value per week from inputs or the epic's share

## Work

Add `planning.CostOfDelay` and `flai cod <E-nnnn|S-nnnn> [--json]`. With inputs, the value per week is revenue plus penalty plus the hours of time lost per cycle times `planning.hour_rate` times the cycles in a week (604800 seconds over `planning.cycle`); time lost with no hour rate is refused, naming `planning.hour_rate`. A story with no inputs gets its epic's value (from the epic's inputs, else its value) apportioned over the epic's open stories without inputs of their own, by forecast duration: each story's `forecast.duration`, else the duration `planning.Forecast` computes. An epic with neither inputs nor value, or a story whose epic has none, is refused, saying the inputs are the operator's. Amounts are rounded to two decimals in `planning.currency`. The command only prints; the planner or a person writes the value with `flai edit` or `item_edit`.

Waits for T-0816: it uses `planning.Forecast` for durations and registers its command in `root.go` beside `flai forecast`.

## Done when

- [ ] `flai cod` prints the value, the currency, and a basis; `--json` adds the inputs or the epic's share and the stories it was apportioned over
- [ ] Fixture tests pin the value from inputs (revenue, penalty, time lost with a cycle other than a week), the apportioned share, and each refusal

## Notes
