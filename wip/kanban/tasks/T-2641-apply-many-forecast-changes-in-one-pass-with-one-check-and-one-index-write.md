---
id: T-2641
type: task
nature: improvement
title: Apply many forecast changes in one pass with one check and one index write
status: done
parent: S-0599
owner: alex
created: 2026-11-07T23:48:43Z
updated: 2026-11-08T00:18:20Z
transitions:
  - to: ready
    at: 2026-11-08T00:06:00Z
    by: system-flow
  - to: in-progress
    at: 2026-11-08T00:06:07Z
    by: system-flow
  - to: done
    at: 2026-11-08T00:18:20Z
    by: agent-S-0599
stream: S-0599
tags: [cli]
touches: [flai/internal/itemedit/forecasts.go, flai/internal/itemedit/forecasts_test.go]
after: [T-2628, T-2632]
usage:
  source: log
  seconds: 733
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 26
      output: 118
      cache_read: 3836192
      cache_write: 31649
      cost: 1.3903
---
# T-2641 Apply many forecast changes in one pass with one check and one index write

## Work

Add `itemedit.ApplyForecasts(repo, changes, opt)` in a new `flai/internal/itemedit/forecasts.go`. Each change names an item and a `ForecastEdit`. In one pass:

1. Read each item, refuse one that `editable` refuses, and apply its edit with `editForecast` and the stamps `Apply` gives (`by`, `at`, `updated`). Skip an item whose forecast does not change.
2. Write them all with `Repo.SaveAll` from T-2628, under one hold of the items lock.
3. Write `wip/agents/index.md` once, from one `repo.List`.
4. Run `check.Items` from T-2632 once over the written items. Put back, with `Repo.Restore`, each item with an error on its own file, as `Apply` does.
5. Record one `Notice` per item written, as `Apply` does, so `flai serve`'s news sees each edit.

It returns the items written, the files to commit, and a warning per item not written. It commits nothing.

Waits for T-2628 and T-2632: it calls `SaveAll` and `check.Items`.

## Done when

- A test shows `ApplyForecasts` writes the same forecast, stamps, and body as `Apply` does, change for change, on the same items.
- A test shows the index is written once and the check runs once, whatever the number of items.
- A test shows a forecast another writer changed meanwhile is refused for that item alone, and a change to another part of it is merged (ADR-0184).
- `flai test flai/internal/itemedit` passes.

## Notes
