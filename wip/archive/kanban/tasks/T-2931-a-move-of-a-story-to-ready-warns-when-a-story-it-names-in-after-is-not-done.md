---
id: T-2931
type: task
nature: improvement
title: A move of a story to ready warns when a story it names in after is not done
status: done
parent: S-0673
owner: alex
created: 2026-11-10T08:49:18Z
updated: 2026-11-10T10:20:24Z
transitions:
  - to: ready
    at: 2026-11-10T09:53:41Z
    by: agent-S-0673
  - to: in-progress
    at: 2026-11-10T09:53:43Z
    by: agent-S-0673
  - to: done
    at: 2026-11-10T10:20:24Z
    by: agent-S-0673
stream: S-0673
tags: [flai, workitem]
touches: [flai/internal/workitem/rules.go, flai/internal/workitem/rules_test.go, flai/cmd/move.go, docs/users/flai-reference.md, design/system/workflow.md]
after: [T-2928]
usage:
  source: log
  seconds: 1601
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 41
      output: 7639
      cache_read: 5690892
      cache_write: 279801
      cost: 2.8894
---
# T-2931 A move of a story to ready warns when a story it names in after is not done

## Work

In `Repo.Move` (`flai/internal/workitem/rules.go`), when a story moves to ready and a story it names in `after` is not done, add a warning naming each such story and saying the story is held until it is done, in the words the hold's reason uses (`held (after): waits for …`). Warn; do not refuse. The move to in-progress's hold warning stays as it is.

Add a test in `rules_test.go`: a story whose `after` names a story in progress gets the warning on its move to ready, and one whose `after` names only done stories gets none.

Say it in `flai move`'s help in `flai/cmd/move.go` beside "WIP limit breaches warn", regenerate `docs/users/flai-reference.md` with `make flai-reference`, and say it in `design/system/workflow.md` where holds and the move to ready are described.

It waits for T-2928, whose ADR decides the warning. It shares no path with T-2930, so the two run together.

## Done when

- The new test passes, and `flai test` on the changed paths passes, `TestReferenceIsCurrent` among them.
- `flai move` to ready of a story with an open `after` prints the warning.

## Notes
