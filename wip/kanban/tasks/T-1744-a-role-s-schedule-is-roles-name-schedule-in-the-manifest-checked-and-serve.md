---
id: T-1744
type: task
nature: feature
title: "A role's schedule is roles.<name>.schedule in the manifest, checked, and serve starts the role each time it comes round through cadence"
status: done
parent: S-0418
owner: alex
created: 2026-10-16T05:45:31Z
updated: 2026-10-29T17:11:48Z
transitions:
  - to: ready
    at: 2026-10-29T16:52:48Z
    by: agent-S-0418
  - to: in-progress
    at: 2026-10-29T16:52:49Z
    by: agent-S-0418
  - to: done
    at: 2026-10-29T17:11:48Z
    by: agent-S-0418
stream: S-0418
tags: [cli]
touches: [flai/internal/serve/role_schedule.go, flai/internal/serve/role_schedule_test.go, flai/internal/manifest/manifest.go, flai/internal/manifest/manifest_test.go, flai/internal/manifest/settings.go, flai/internal/manifest/settings_test.go, docs/operators/settings.md]
after: [T-1748]
usage:
  source: log
  seconds: 1139
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 65
      output: 26502
      cache_read: 4914041
      cache_write: 128635
      cost: 2.2947
---
# T-1744 A role's schedule is roles.<name>.schedule in the manifest, checked, and serve starts the role each time it comes round through cadence

## Work

Give a defined role a schedule and start it when the schedule comes round.

- Add `schedule` to the manifest's role settings, `roles.<name>.schedule`: a five-field cron expression in UTC or `daily`, empty by default, parsed as `analysis.schedule` is (`parseSchedule`).
- `RoleErrors` reports a value that does not parse or never comes round, as the other schedule settings are reported, and a schedule set on a built-in role, whose schedule is its own block's. `flai check`'s `roles.file` finding for an unknown schedule setting on a trigger stays as S-0536 has it.
- A `roleScheduler` per project, in `role_schedule.go`, with a `look(ctx)`: for each role a role file defines but the built-in ones whose trigger is `{schedule: roles.<name>.schedule}`, a `cadence` (in `replan.go`, left where it is). Each time it comes round with the role's action on, start the role through `startRole` (T-1748) with the trigger `schedule <spec>`. A refused start is logged and not tried again before the schedule next comes round; a start the budget holds is tried again at each look until it starts, as the analysis scheduler's is. Times missed between two looks come round once; while the action is off it keeps up with the schedule.
- A role that takes an item has no item to run a schedule on: report it as a refusal, logged, rather than guess one.
- Test it with a fake clock, and the manifest's parsing and its error.

Waits for T-1748, whose `startRole` it calls. Runs beside T-1743.

## Done when

`flai test` passes on the paths it changed. `analysis_schedule_test.go` is unchanged.

## Notes
