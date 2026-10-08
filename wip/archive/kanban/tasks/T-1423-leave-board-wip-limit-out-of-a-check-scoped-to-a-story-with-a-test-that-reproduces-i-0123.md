---
id: T-1423
type: task
nature: improvement
title: Leave board.wip-limit out of a check scoped to a story, with a test that reproduces I-0123
status: done
parent: S-0348
owner: alex
created: 2026-10-08T08:59:00Z
updated: 2026-10-08T09:27:52Z
transitions:
  - to: ready
    at: 2026-10-08T09:19:03Z
    by: agent-S-0348
  - to: in-progress
    at: 2026-10-08T09:20:09Z
    by: agent-S-0348
  - to: done
    at: 2026-10-08T09:27:52Z
    by: agent-S-0348
stream: S-0348
tags: [flai, check]
touches: [flai/internal/check/scope.go, flai/internal/check/scope_test.go, flai/cmd/check.go, flai/cmd/check_test.go, docs/users/flai-reference.md]
after: [T-1422]
usage:
  source: log
  seconds: 463
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 25
      output: 6348
      cache_read: 1388719
      cache_write: 57694
      cost: 0.8058
---
# T-1423 Leave board.wip-limit out of a check scoped to a story, with a test that reproduces I-0123

## Work

In `ScopeToStory` (`flai/internal/check/scope.go`), drop a `board.wip-limit` finding with `res.drop`, beside `item.archive` and the findings on another open story's narrative, as the ADR of T-1422 decides. Extend its doc comment with the ADR and I-0123.

In `flai/internal/check/scope_test.go`, reproduce I-0123: a repository whose board has more stories in progress than its limit, checked with `--story` on one of them, keeps no `board.wip-limit` finding and counts none as outside. Assert that the same repository checked without `--story` still warns on it and fails `--strict`.

In `flai/cmd/check.go`, update `recordOutside`'s comment and the `--record-issues` text of the command's long help to say a `board.wip-limit` never reaches an issue. Add a case to `flai/cmd/check_test.go` that runs `flai check --story S-nnnn --record-issues` over such a board and finds no issue opened. Regenerate `docs/users/flai-reference.md` with `make flai-reference`.

It waits for T-1422 so that the comments and the help cite the ADR's number.

## Done when

- `flai check --story S-nnnn --record-issues` over a board with in-progress over its limit prints no `board.wip-limit` and records no issue, shown by the new tests.
- `flai check --strict` without `--story` still fails on it.
- `flai test` on the changed paths passes.

## Notes
