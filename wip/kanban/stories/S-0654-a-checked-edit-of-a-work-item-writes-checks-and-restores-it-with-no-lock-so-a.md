---
id: S-0654
type: story
nature: improvement
title: A checked edit of a work item writes, checks, and restores it with no lock, so a reader in between sees an edit the check refuses
status: in-progress
owner: alex
created: 2026-11-10T05:45:23Z
updated: 2026-11-10T07:34:29Z
transitions:
  - to: ready
    at: 2026-11-10T07:12:12Z
    by: orchestrator
  - to: in-progress
    at: 2026-11-10T07:34:29Z
    by: agent-S-0654
tags: [flai, itemedit, locking]
topics: [work-items, claims]
touches: [flai/internal/itemedit/itemedit.go, flai/internal/itemedit/lock_test.go, flai/internal/itemedit/forecasts.go, flai/internal/itemedit/forecasts_test.go, flai/internal/workitem/itemlock.go, design/adrs, design/adrs/README.md, design/system/flai-cli.md, design/system/workflow.md, docs/users/flai.md, design/issues/I-0328-a-checked-edit-of-a-work-item-writes-checks-and-restores-it-with-no-lock-so-a.md, design/issues/summary.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 576
  estimated: true
  priced_by: estimate
  turns:
    - day: 2026-11-10
      ceremony: 1
      hand_edits: 2
      work: 14
  models:
    - model: claude-opus-5-5
      input: 90
      output: 628
      cache_read: 8983687
      cache_write: 581545
      cost: 3.5191
  strategic:
    - kind: orchestrator
      seconds: 143
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 40
          output: 289
          cache_read: 2478483
          cache_write: 161902
          cost: 0.7099
cost_of_delay:
  inputs:
    time_lost_per_cycle: 1m
    by: orchestrator
    at: 2026-11-10T07:04:53Z
  value: 2.5
  by: planner-S-0654
  at: 2026-11-10T07:05:18Z
forecast:
  duration: 1h30m
  delivery: 2026-11-10T09:51:00Z
  basis: "Its own forecast of 1h30m; 5th in the pull order with an in-progress limit of 5, behind S-0628, S-0636, S-0638, S-0647, S-0649, S-0650, S-0655, S-0651 and S-0653."
  by: flai
  at: 2026-11-10T07:12:18Z
finalized:
  by: orchestrator
  at: 2026-11-10T07:12:08Z
---
# S-0654 A checked edit of a work item writes, checks, and restores it with no lock, so a reader in between sees an edit the check refuses

## Goal

This story remediates [I-0328](../../../design/issues/I-0328-a-checked-edit-of-a-work-item-writes-checks-and-restores-it-with-no-lock-so-a.md), "A checked edit of a work item writes, checks, and restores it with no lock, so a reader in between sees an edit the check refuses". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0328 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0328 is closed with `flai issue close I-0328 --reason` saying what fixed it

## Tasks
- T-2844 itemedit.Apply holds the checked-write lock from its first save until it keeps or restores the edit
- T-2847 itemedit.ApplyForecasts holds the checked-write lock from its save until it keeps or restores the forecasts
- T-2848 Record in an ADR refining ADR-0272, the design, and the user docs that a checked edit holds the checked-write lock
- T-2850 Close I-0328 naming the lock the checked edit now holds

## Notes

### Planning

Proposed fix, from I-0328's one instance and ADR-0272: every checked write of an existing work item holds the checked-write lock (`(*workitem.Repo).LockCheckedWrites`, `.flai-cache/items-checked.lock`) from its first save until the check keeps the write or the restore takes it out, as `itemnew.Create` does for a creation. flai serve's look (`readyChecked` in `flai/internal/serve/agents.go`) already reads under that lock, so it then never judges a claim by touches a refused edit wrote. Two writers do this today:

- `itemedit.Apply` (`flai edit`, `flai criteria`, the MCP tools `item_edit` and `criteria_tick`): it saves the item, its parents, the narrative, and `wip/agents/index.md`, runs `check.RunItemWrite`, and calls `undo.restore` on a refusal.
- `itemedit.ApplyForecasts` (flai serve's replan): it saves many stories with `SaveAll`, checks them, and calls `Repo.Restore` on each one refused. A forecast does not change a claim, but it is the same window, and holding the lock keeps one rule for every checked write.

`flai touches` (`flai/cmd/touches.go`) saves with no check and no restore, so it has no such window, though I-0328 names it. It is left as it is.

Touches, and where each came from:

| Touch | Source | Why |
|-------|--------|-----|
| `flai/internal/itemedit/itemedit.go` | layout | `Apply`'s write, check, and restore |
| `flai/internal/itemedit/lock_test.go` | layout | the restore tests live here; the new lock test joins them |
| `flai/internal/itemedit/forecasts.go` | layout | `ApplyForecasts`' `SaveAll`, check, and `Restore` |
| `flai/internal/itemedit/forecasts_test.go` | layout | its test |
| `flai/internal/workitem/itemlock.go` | layout | `LockCheckedWrites`' comment says a checked creation holds it; it will say a checked edit does too |
| `design/adrs/` | design | folder kept: the new ADR refining ADR-0272 gets its number and slug from `flai adr new` when it is written, so no task can name its file yet |
| `design/adrs/README.md` | design | the index `flai adr new` writes again |
| `design/system/flai-cli.md` | design | the `flai edit` row describes the checked edit |
| `design/system/workflow.md` | design | **Held** says what the look reads under the lock |
| `docs/users/flai.md` | design | the `flai edit` section describes the checked edit |
| `design/issues/I-0328-…` and `design/issues/summary.md` | declared by the goal | `flai issue close I-0328` writes both |

`flai touches suggest S-0654` listed nothing: the story declared no touches.

No `after`: the story reads nothing another open story writes. ADR-0272's lock is on main since S-0633.

Forecast: 1h30m, adjusted from `flai forecast`'s 24m (a median of 99 s per unit of size over 146 done improvement stories, times size 14). S-0633 added the same lock for a creation. It was forecast at 40m and took 3h48m of agent time. This story has two writers to lock, a lock order to keep, and a new ADR. flai's replan keeps the duration and sets the delivery by the pull order.

Cost of delay: 2.50 USD a week, as `flai cod` works it out from the operator's input on TH-0975, `time_lost_per_cycle: 1m`, at 150 USD an hour over one cycle a week. It stands unadjusted: I-0328 has no instance yet.
