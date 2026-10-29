---
id: T-2301
type: task
nature: feature
title: orchestration.coalesce, orchestration.run_timeout, and orchestration.schedule are manifest keys flai check validates and flai manifest set sets, with ADR-0209's defaults
status: done
parent: S-0536
owner: alex
created: 2026-10-28T21:45:17Z
updated: 2026-10-29T07:24:58Z
transitions:
  - to: ready
    at: 2026-10-29T07:21:04Z
    by: agent-S-0536
  - to: in-progress
    at: 2026-10-29T07:21:05Z
    by: agent-S-0536
  - to: done
    at: 2026-10-29T07:24:58Z
    by: agent-S-0536
stream: S-0536
tags: [flai, manifest, orchestrator]
touches: [flai/internal/manifest/manifest.go, flai/internal/manifest/manifest_test.go, flai/internal/manifest/settings.go, flai/internal/manifest/settings_test.go, docs/operators/settings.md]
usage:
  source: log
  seconds: 233
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 53
      output: 20002
      cache_read: 3873896
      cache_write: 123124
      cost: 1.8708
---
# T-2301 orchestration.coalesce, orchestration.run_timeout, and orchestration.schedule are manifest keys flai check validates and flai manifest set sets, with ADR-0209's defaults

## Work

ADR-0209 §1, §2, and §4, criterion 4. Add three keys to `manifest.Orchestration` in `flai/internal/manifest/manifest.go`:

- `coalesce`, a Go duration, `30s` when unset (`DefaultCoalesce`): how long after the first queued reason flai serve starts the orchestrator's run.
- `run_timeout`, a Go duration, `20m` when unset (`DefaultRunTimeout`): how long a run may hold before flai serve stops it.
- `schedule`, a five-field cron expression in UTC or `daily`, as `digest.schedule` and `planning.schedule` are, unset by default: a periodic look.

Give each an `OrDefault` accessor. Report a value that is not a duration, a negative or zero duration, or a cron expression that does not parse as a `manifest.orchestration` finding, as the release policy's errors are. Add each to the settings list in `flai/internal/manifest/settings.go`, beside `orchestration.policy`, so that `flai manifest set` sets it and the dashboard's settings page shows it, with its default. Add a row for each to `docs/operators/settings.md`, which the settings test requires.

Waits for nothing: the first layer. The trigger task names `orchestration.schedule` in the role file and the reason queue reads all three, so both wait for this.

## Done when

- Tests in `manifest_test.go` cover each default, a valid value, and each invalid value as a `manifest.orchestration` finding.
- A test in `settings_test.go` sets each with the settings list and reads it back.
- `docs/operators/settings.md` has a row for each key, and the settings index test passes.
- `flai test` passes on the changed paths.

## Notes
