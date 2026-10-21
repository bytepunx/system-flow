---
id: T-1921
type: task
nature: improvement
title: flai carries a task-done rules number in buildinfo, and flai version --json prints it
status: done
parent: S-0445
owner: alex
created: 2026-10-20T20:41:50Z
updated: 2026-10-20T20:52:03Z
transitions:
  - to: ready
    at: 2026-10-20T20:45:44Z
    by: agent-S-0445
  - to: in-progress
    at: 2026-10-20T20:45:44Z
    by: agent-S-0445
  - to: done
    at: 2026-10-20T20:52:03Z
    by: agent-S-0445
stream: S-0445
tags: [flai]
touches: [flai/internal/buildinfo/buildinfo.go, flai/internal/buildinfo/buildinfo_test.go, flai/internal/buildinfo/task-done-rules.txt, flai/cmd/version.go, flai/cmd/version_test.go]
after: [T-1920]
usage:
  source: log
  seconds: 379
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 15
      output: 4263
      cache_read: 1183152
      cache_write: 31066
      cost: 0.5491
---
# T-1921 flai carries a task-done rules number in buildinfo, and flai version --json prints it

## Work

- Add `TaskDoneRules`, an integer constant, to `flai/internal/buildinfo/buildinfo.go`, with a comment saying when to bump it, as the ADR of T-1920 decides. Start it at 1. `buildinfo` is a leaf package, so `release`, `taskdone`, and `mcpserver` can all read it without an import cycle (`taskdone` already imports `release`).
- Add `task_done_rules` to `versionReport` in `flai/cmd/version.go`, so `flai version --json` prints it beside `version`.
- Test that the JSON carries the number, in a new `flai/cmd/version_test.go`, and that the constant is positive, in `buildinfo_test.go`.

It waits for T-1920 because the ADR fixes the name, the field, and the bump rule.

## Done when

- `flai version --json` prints `task_done_rules` with the constant's value.
- The tests pass with `flai test` on the four paths.

## Notes

ADR-0160's caught-up check (T-1924) reads this field from the installed flai.
