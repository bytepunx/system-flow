---
id: T-0646
type: task
nature: remediation
title: flai records the variables a project was rendered with, and flai upgrade renders with them
status: done
parent: S-0185
owner: arobson
created: 2026-10-01T08:49:41Z
updated: 2026-10-01T08:57:03Z
transitions:
  - to: ready
    at: 2026-10-01T08:49:50Z
    by: agent-S-0185
  - to: in-progress
    at: 2026-10-01T08:50:11Z
    by: agent-S-0185
  - to: done
    at: 2026-10-01T08:57:03Z
    by: agent-S-0185
stream: S-0185
tags: []
touches: [flai/cmd/upgrade.go, flai/cmd/new.go, flai/cmd/upgrade_test.go, flai/internal/lock, flai/internal/upgrade]
usage:
  source: log
  seconds: 412
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 108
      output: 740
      cache_read: 6531355
      cache_write: 88194
      cost: 2.6185
---
# T-0646 flai records the variables a project was rendered with, and flai upgrade renders with them

## Work

- `system-flow.lock.yaml` records every template variable's value under `vars`; `flai new` writes them, and `flai upgrade` and `--relock` write the values they rendered with.
- `flai upgrade` resolves each template variable, in order: `--var`; the manifest's own field for the five standard ones (`name`, `key`, `description`, `owner`, `repo`); the lock's recorded value; the template's default. It refuses an unknown `--var`, and a `--var` for a variable the manifest holds, saying where to change it.
- A required variable still empty is named, and the upgrade refuses before rendering or writing anything.

## Done when

- Tests in `flai/cmd/upgrade_test.go` cover an upgrade of a fork whose template adds a variable: with the value recorded by `flai new`, with the default of a variable added since, with `--var`, and the refusal of a required variable with no value.
- `make test` and lint pass.

## Notes
