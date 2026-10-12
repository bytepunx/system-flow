---
id: T-1452
type: task
nature: improvement
title: Leave threads.archived outside the story out of a check scoped to a story, with a test that reproduces I-0135
status: in-progress
parent: S-0366
owner: alex
created: 2026-10-12T14:59:41Z
updated: 2026-10-12T15:27:27Z
transitions:
  - to: ready
    at: 2026-10-12T15:27:26Z
    by: agent-S-0366
  - to: in-progress
    at: 2026-10-12T15:27:27Z
    by: agent-S-0366
stream: S-0366
tags: [flai, check]
touches: [flai/internal/check/scope.go, flai/internal/check/scope_test.go, flai/cmd/check.go, flai/cmd/check_test.go, docs/users/flai-reference.md]
after: [T-1451]
---
# T-1452 Leave threads.archived outside the story out of a check scoped to a story, with a test that reproduces I-0135

## Work

In `ScopeToStory` (`flai/internal/check/scope.go`), add `threads.archived` to the rules dropped with `res.drop` when `!inside(f.Path)`, beside `item.archive`, `board.wip-limit`, and S-0365's `issues.no-story`, as T-1451's ADR decides. If the ADR decides otherwise, build what it decides instead. Extend the function's doc comment with the ADR and I-0135. The finding comes from `flai/internal/check/check.go`; leave its production alone unless the ADR changes it.

In `flai/internal/check/scope_test.go`, reproduce I-0135. Set up a repository with an archived story and an open thread anchored on it, as TH-0389 was on S-0338, and a story in progress that does not change it. Checked with `--story` on that story, it keeps no `threads.archived` and counts none as outside. A thread anchored on the story itself keeps the finding. Checked without `--story`, it still warns.

In `flai/cmd/check.go`, update `recordOutside`'s comment and the `--record-issues` text of the command's long help to say a `threads.archived` outside the story never reaches an issue. Its help names `threads.archived` as the example of an issue title; change the example to a rule that is still recorded. Add a case to `flai/cmd/check_test.go` that runs `flai check --story S-nnnn --record-issues` over such a thread and finds no issue opened or bumped. Regenerate `docs/users/flai-reference.md` with `make flai-reference`.

It waits for T-1451 so that the comments and the help cite the ADR's number.

## Done when

- `flai check --story S-nnnn --record-issues` over an open thread on an archived item outside the story prints no `threads.archived` and records no issue, shown by the new tests.
- `flai check` without `--story` still warns `threads.archived` on it.
- `flai test` on the changed paths passes.

## Notes
