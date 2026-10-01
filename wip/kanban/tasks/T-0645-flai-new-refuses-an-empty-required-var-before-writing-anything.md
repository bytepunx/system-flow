---
id: T-0645
type: task
nature: remediation
title: flai new refuses an empty required --var before writing anything
status: done
parent: S-0185
owner: arobson
created: 2026-10-01T08:49:41Z
updated: 2026-10-01T08:50:11Z
transitions:
  - to: ready
    at: 2026-10-01T08:49:50Z
    by: agent-S-0185
  - to: in-progress
    at: 2026-10-01T08:49:50Z
    by: agent-S-0185
  - to: done
    at: 2026-10-01T08:50:11Z
    by: agent-S-0185
stream: S-0185
tags: []
touches: [flai/cmd/new.go, flai/cmd/new_test.go]
usage:
  source: log
  seconds: 21
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 6
      output: 48
      cache_read: 339701
      cache_write: 5999
      cost: 0.1368
---
# T-0645 flai new refuses an empty required --var before writing anything

## Work

In `collectVars` (`flai/cmd/new.go`), apply the required check to a value given with `--var` as to a default or a prompted value, so `flai new` and `flai import` stop before anything is written and name the variable.

## Done when

- `flai new --var project_name=` exits non-zero naming `project_name`, and the target directory has no files.
- A test in `flai/cmd/new_test.go` covers it and fails without the fix.

## Notes
