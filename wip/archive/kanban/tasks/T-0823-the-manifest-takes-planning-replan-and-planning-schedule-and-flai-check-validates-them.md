---
id: T-0823
type: task
nature: improvement
title: The manifest takes planning.replan and planning.schedule, and flai check validates them
status: done
parent: S-0211
owner: alex
created: 2026-10-04T20:42:09Z
updated: 2026-10-04T20:55:56Z
transitions:
  - to: ready
    at: 2026-10-04T20:42:48Z
    by: agent-S-0211
  - to: in-progress
    at: 2026-10-04T20:49:38Z
    by: agent-S-0211
  - to: done
    at: 2026-10-04T20:55:56Z
    by: agent-S-0211
stream: S-0211
tags: [flai]
touches: [flai/internal/manifest, design/system/project-manifest.md, docs/operators/settings.md]
after: [T-0820]
usage:
  source: log
  seconds: 378
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 45
      output: 17414
      cache_read: 2796764
      cache_write: 70684
      cost: 1.3317
---
# T-0823 The manifest takes planning.replan and planning.schedule, and flai check validates them

## Work

`manifest.Planning` gains `Replan` (`never`, `deterministic`, `agent`; unset is `deterministic`) and `Schedule` (a cron expression or `daily`, parsed with `flai/internal/cron`; unset is none), with accessors that give the default and the parsed schedule, and `flai check` reports a bad value as `manifest.planning`. `design/system/project-manifest.md` shows both keys in the example and describes them; `docs/operators/settings.md` gets a row for each, which the settings doc test requires. It waits for T-0820, whose parser validates the schedule.

## Done when

- [ ] Both keys parse, default as ADR-0084 says, and bad values are `manifest.planning` errors
- [ ] The settings doc test and `flai/internal/manifest` tests pass

## Notes
