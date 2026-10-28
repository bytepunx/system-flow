---
id: T-2251
type: task
nature: remediation
title: A checked creation refuses only a new finding inside the new item's story, with a test that reproduces I-0191
status: done
parent: S-0525
owner: alex
created: 2026-10-28T14:52:28Z
updated: 2026-10-28T22:56:55Z
transitions:
  - to: ready
    at: 2026-10-28T22:53:04Z
    by: agent-S-0525
  - to: in-progress
    at: 2026-10-28T22:53:05Z
    by: agent-S-0525
  - to: done
    at: 2026-10-28T22:56:55Z
    by: agent-S-0525
stream: S-0525
tags: [flai, check]
touches: [flai/internal/itemnew/itemnew.go, flai/internal/itemnew/itemnew_test.go, flai/internal/itemedit/itemedit.go, flai/internal/check/scope.go]
usage:
  source: log
  seconds: 230
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 15
      output: 3645
      cache_read: 1075188
      cache_write: 33547
      cost: 0.5564
---
# T-2251 A checked creation refuses only a new finding inside the new item's story, with a test that reproduces I-0191

## Work

First the reproduction: a test in `flai/internal/itemnew/itemnew_test.go` on a project with no `.gitignore` and no `.flai-cache` (the `project` helper writes a `.gitignore`; this test must not) whose first checked creation, a story with a body, is refused with the finding `layout.gitignore`, as I-0191 saw. It fails on main today.

Then the fix in `itemnew.Create` (`flai/internal/itemnew/itemnew.go`): scope the check run with the new item in place as `itemedit.Apply` scopes a checked edit (ADR-0197): the new story's own scope, a new task's parent story, and for an epic the paths the creation wrote. A new finding the scope marks outside, such as `layout.gitignore` on the repository root, goes into `Result.Warnings` and the creation stands; a new finding inside the scope, or an error on the item's own file, still refuses it and undoes the write; `wip.overlap`, `touches.folder`, and `wip.after-hub` stay warnings as today. Reuse `scope` from `flai/internal/itemedit/itemedit.go`: export it, or move it beside `check.ScopeToStories` in `flai/internal/check/scope.go`, whichever keeps the import graph acyclic (`itemedit` imports `itemnew` only in a test). The paths the creation wrote are the item's file and the parent's file it relinks.

Keep the lock where ADR-0184 puts it and write no `.gitignore`: the remedy is the write's check, not the lock. The reproduction passes once the fix is in, and the existing `itemnew` tests, which refuse a story missing a template section, still refuse: that finding is on the item's own file.

## Done when

- A test in `flai/internal/itemnew/itemnew_test.go` creates a story with a body as the first write of a project with no `.gitignore` and no `.flai-cache`, and the creation stands with `layout.gitignore` among its warnings; the test fails on main before the fix and passes after it.
- A test shows a checked creation is still refused for a new finding inside its scope, such as a body missing a template section.
- `flai test flai/internal/itemnew flai/internal/itemedit flai/internal/check` passes.

## Notes

Reproduced by the planner on 2026-10-28 on a temp copy of `flai/internal/metrics/testdata/good`: `flai epic new --body-stdin` as the first write refuses itself; `flai edit` on the same copy warns and stands.
