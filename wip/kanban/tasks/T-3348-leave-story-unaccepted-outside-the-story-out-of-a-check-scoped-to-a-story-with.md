---
id: T-3348
type: task
nature: remediation
title: Leave story.unaccepted outside the story out of a check scoped to a story, with a test that reproduces I-0371
status: done
parent: S-0781
owner: alex
created: 2026-11-13T22:53:38Z
updated: 2026-11-13T23:27:56Z
transitions:
  - to: ready
    at: 2026-11-13T23:10:59Z
    by: system-flow
  - to: in-progress
    at: 2026-11-13T23:10:59Z
    by: system-flow
  - to: done
    at: 2026-11-13T23:27:56Z
    by: agent-S-0781
stream: S-0781
tags: [flai, check]
touches: [flai/internal/check/scope.go, flai/internal/check/scope_test.go, flai/cmd/check.go, flai/cmd/check_test.go, docs/users/flai-reference.md, flai/cmd/verify_test.go, flai/internal/verify/story_test.go]
after: [T-3347]
usage:
  source: log
  seconds: 1017
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 74
      output: 22246
      cache_read: 3707903
      cache_write: 165478
      cost: 2.2353
---
# T-3348 Leave story.unaccepted outside the story out of a check scoped to a story, with a test that reproduces I-0371

## Work

As the ADR of T-3347 decides. It waits for T-3347 for that decision and the ADR's number, which the comments name. If the ADR decides the root fix alone, ask on the story's thread to cancel this task.

In `ScopeToStories` (`flai/internal/check/scope.go`), add `story.unaccepted` to the rules dropped when outside the story, beside `item.blocked-closed`, and extend the comment above it with the reason, the ADR, and I-0371. The rule itself in `flai/internal/check/check_items.go` does not change. In `flai/cmd/check.go`, name it in the `--record-issues` help and the comment of `recordOutside`, as S-0761 did for `item.blocked-closed`. Regenerate `docs/users/flai-reference.md` with `make flai-reference`.

Reproduce I-0371 in the tests, following the `TestScopedCheckLeavesOut…` tests in `flai/internal/check/scope_test.go` and the `TestCheckRecordIssuesLeavesOut…` tests in `flai/cmd/check_test.go`. Make an archived done story whose `story/S-nnnn` branch still exists, and check that a check scoped to another story leaves the finding out and `--record-issues` records no issue. Also check that the unscoped check still reports it.

## Done when

- A check scoped to a story drops a `story.unaccepted` on another story, and the unscoped check still reports it.
- `flai check --story S-nnnn --record-issues` records no issue for it.
- The tests that show both fail without the change, and `flai test` passes on the changed files.

## Notes
