---
id: T-0741
type: task
nature: feature
title: flai edit, flai story new --draft, and flai move --yes set and clear the planning fields
status: done
parent: S-0199
owner: alex
created: 2026-10-03T05:41:12Z
updated: 2026-10-03T06:01:03Z
transitions:
  - to: ready
    at: 2026-10-03T05:48:08Z
    by: agent-S-0199
  - to: in-progress
    at: 2026-10-03T05:48:08Z
    by: agent-S-0199
  - to: done
    at: 2026-10-03T06:01:03Z
    by: agent-S-0199
stream: S-0199
tags: []
touches: [flai/internal/itemedit, flai/internal/itemnew, flai/cmd]
after: [T-0739]
usage:
  source: log
  seconds: 775
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 105
      output: 38916
      cache_read: 6292250
      cache_write: 158737
      cost: 2.9735
---
# T-0741 flai edit, flai story new --draft, and flai move --yes set and clear the planning fields

## Work

`itemedit.Change` gains `Draft`, a cost of delay edit (each input and the value: nil leaves it, empty removes it), a forecast edit (`duration`, `delivery`, `basis`), and clears for both blocks; a change to a block stamps its `by` (the edit's `By`) and `at`. `flai edit` takes `--draft`/`--no-draft`, `--revenue-per-week`, `--penalty-per-week`, `--time-lost-per-cycle`, `--cost-of-delay-value`, `--clear-cost-of-delay`, `--forecast-duration`, `--forecast-delivery`, `--forecast-basis`, `--clear-forecast`. `flai story new --draft` (through `itemnew` and `workitem.NewOptions`) makes a draft. `flai move <story> ready` on a draft is refused unless `--yes`, which finalizes it. `flai show` prints the fields. Waits for T-0739; shares no path with T-0740, so the two can run together.

## Done when

Each flag sets and clears its field with the right `by` and `at`, `story new --draft` writes `draft: true`, and the move is refused without `--yes` and finalizes with it, each with a test.

## Notes
