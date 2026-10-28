---
id: T-1721
type: task
nature: feature
title: flai/internal/roles reads and validates a project's role files into roles with their fields, conditions, and definitions
status: done
parent: S-0414
owner: alex
created: 2026-10-16T05:42:59Z
updated: 2026-10-28T13:32:12Z
transitions:
  - to: ready
    at: 2026-10-28T13:18:47Z
    by: agent-S-0414
  - to: in-progress
    at: 2026-10-28T13:18:48Z
    by: agent-S-0414
  - to: done
    at: 2026-10-28T13:32:12Z
    by: agent-S-0414
stream: S-0414
tags: [cli]
touches: [flai/internal/roles/roles.go, flai/internal/roles/roles_test.go, flai/internal/roles/conditions.go, flai/internal/roles/conditions_test.go]
after: [T-1720]
usage:
  source: log
  seconds: 804
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 50
      output: 23945
      cache_read: 2820279
      cache_write: 118781
      cost: 1.7169
---
# T-1721 flai/internal/roles reads and validates a project's role files into roles with their fields, conditions, and definitions

## Work

Add the package `flai/internal/roles`.

- `Load(project)` reads every file under the roles folder into `Role` values. Each value holds the ADR's fields, its body, and the text below its marker.
- `Validate` returns the problems the check reports: a field missing, an unknown trigger, item kind, or condition, and a call that only an operator's permission may allow, granted outside any permission.
- `conditions.go` lists the named conditions the ADR settles. Each condition carries what it allows, so the guard (S-0416) can decide from it.
- Tests cover a valid file, each problem, and a project with no roles folder.

Waits for T-1720, the ADR that fixes the fields.

## Done when

`go test ./internal/roles/...` passes and `flai test` on the package is clean.

## Notes

- After TH-0471's answers, the roles folder is `design/workflows/roles` (answer 2). The role file holds the fields its harness definition is generated from (answer 6). `Validate` checks those fields as T-1720's ADR lists them.
