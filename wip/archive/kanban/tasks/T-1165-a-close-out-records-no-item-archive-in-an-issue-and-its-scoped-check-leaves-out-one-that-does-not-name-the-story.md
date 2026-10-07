---
id: T-1165
type: task
nature: improvement
title: A close-out records no item.archive in an issue, and its scoped check leaves out one that does not name the story
status: done
parent: S-0280
owner: alex
created: 2026-10-07T15:03:34Z
updated: 2026-10-07T23:29:23Z
transitions:
  - to: ready
    at: 2026-10-07T23:24:58Z
    by: agent-S-0280
  - to: in-progress
    at: 2026-10-07T23:24:58Z
    by: agent-S-0280
  - to: done
    at: 2026-10-07T23:29:23Z
    by: agent-S-0280
stream: S-0280
tags: [flai]
touches: [flai/internal/check/scope.go, flai/internal/check/scope_test.go, flai/cmd/check.go, flai/cmd/check_test.go, docs/users/flai-reference.md]
after: [T-1164]
usage:
  source: log
  seconds: 265
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 23
      output: 6013
      cache_read: 1702522
      cache_write: 50918
      cost: 0.8682
---
# T-1165 A close-out records no item.archive in an issue, and its scoped check leaves out one that does not name the story

## Work

Build the remedy T-1164's ADR decides. It waits for T-1164, because the ADR settles what is left out and what is recorded.

- In `ScopeToStory` (`flai/internal/check/scope.go`), leave out an `item.archive` that does not name the story, as it leaves out a `wip.overlap` that does not name it (ADR-0115 §3).
- In `recordOutside` (`flai/cmd/check.go`), record no `item.archive`, beside the `wip.overlap` exemption. Say so in the command's help, and regenerate `docs/users/flai-reference.md` with `make flai-reference`.
- Reproduce I-0078 in tests:
  - `flai/internal/check/scope_test.go`: a cancelled story left in `wip/kanban` beside the story checked; the scoped result holds no `item.archive`.
  - `flai/cmd/check_test.go`: `flai check --story S-nnnn --record-issues` with that fixture records no issue for `item.archive`.
- An unscoped `flai check` still warns `item.archive`, and an existing test shows it.

## Done when

- Both tests fail before the change and pass after it.
- `flai test` passes on the changed paths.
- `flai check --strict` passes.

## Notes
