---
id: T-2075
type: task
nature: remediation
title: itemedit.Apply refuses a new finding only inside the edited item's story's scope, so criteria_tick is no longer refused on another story's narrative.state or an untracked task file, with tests that reproduce I-0182's two instances
status: done
parent: S-0486
owner: alex
created: 2026-10-23T17:50:36Z
updated: 2026-10-23T23:56:42Z
transitions:
  - to: ready
    at: 2026-10-23T23:45:15Z
    by: agent-S-0486
  - to: in-progress
    at: 2026-10-23T23:45:16Z
    by: agent-S-0486
  - to: done
    at: 2026-10-23T23:56:42Z
    by: agent-S-0486
stream: S-0486
tags: [flai]
touches: [flai/internal/itemedit/itemedit.go, flai/internal/itemedit/criteria_test.go, flai/internal/mcpserver/criteria_test.go]
after: [T-2074]
usage:
  source: log
  seconds: 686
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 64
      output: 18866
      cache_read: 3701645
      cache_write: 118658
      cost: 1.8677
---
# T-2075 itemedit.Apply refuses a new finding only inside the edited item's story's scope, so criteria_tick is no longer refused on another story's narrative.state or an untracked task file, with tests that reproduce I-0182's two instances

## Work

Change the check in `itemedit.Apply` (`flai/internal/itemedit/itemedit.go`, the block after the second `check.Run`) as the ADR of T-2074 decides:

- When the edited item is a story, or a task whose parent is a story, scope the `after` result with `check.ScopeToStories` to that story, giving it the paths the edit wrote (`u.paths(repo)`, relative to the repository) as the changed paths. Build the `known` set from `before` as now. A new finding the scope dropped is gone from the result and refuses nothing; a new finding marked `Outside` goes to `res.Warnings`; a new finding inside the scope, or an error on the item's own file, refuses the edit with the message as now, and `wip.overlap` stays a warning.
- When the edited item is an epic, or a task with no story to resolve, a new finding on a path the edit did not write goes to `res.Warnings` instead of refusing.
- Keep the undo and the `RefusedError` shape unchanged: `flai criteria tick`, `criteria_tick`, `flai edit`, `item_edit`, and the host API print its `Findings` as they do.

Tests, in `flai/internal/itemedit/criteria_test.go` and `flai/internal/mcpserver/criteria_test.go`, reproduce both instances of I-0182 by writing the other story's file between the two runs, with a hook on the check or by preparing the repository so the first run sees it absent and the second present: a `narrative.state` on the narrative of another story in progress whose `## Current state` is the template's placeholder, and a `markdown.MD012` in an untracked task file of another story. Each test asserts the tick is written and the finding is not among the refusal's findings; a third asserts an error the tick itself brings on the item's own file still refuses it. Run `flai test` on the paths changed.

Waits for T-2074, which fixes the rule the code applies.

## Done when

`go test ./internal/itemedit/... ./internal/mcpserver/...` passes with the three new tests, `flai criteria tick` on a story ticks while another open story's narrative is unwritten and while another story's task file lints wrong, and `flai check --strict` and `flai test` on the changed paths report nothing.

## Notes

`flai/internal/itemedit/itemedit.go` is also a touch of S-0468, in review: when it is accepted while this story is in progress, `flai stream sync` brings its change, and the `overlapped` notice says to sync and test again.
