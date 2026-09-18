---
id: T-0146
type: task
nature: improvement
title: "flai move: allow ready and in-progress without tasks, refuse review without tasks"
status: done
parent: S-0049
owner: alex
created: 2026-09-18T18:08:29Z
updated: 2026-09-18T18:11:59Z
transitions:
  - to: ready
    at: 2026-09-18T18:11:03Z
    by: alex
  - to: in-progress
    at: 2026-09-18T18:11:04Z
    by: alex
  - to: done
    at: 2026-09-18T18:11:59Z
    by: alex
stream: S-0049
tags: []
touches: [flai/internal/workitem]
---

# T-0146 flai move: allow ready and in-progress without tasks, refuse review without tasks

## Work
In `flai/internal/workitem/rules.go`, remove the task requirement from the `Ready` case and add a `Review` case that refuses a story with no child tasks, with a message naming the rule and the command to fix it. Keep the acceptance criteria rule on ready unchanged. In `workitem_test.go`, replace the "needs at least one task" ready assertion with: ready without tasks succeeds, in-progress without tasks succeeds, review without tasks is refused, review with a task succeeds, ready without criteria is still refused.

## Done when
- The new and changed cases in `workitem_test.go` pass and fail without the change
- `go test -race ./internal/workitem/...` passes
- The MCP `item_move` and `flai accept` paths are confirmed to go through `Repo.Move`, recorded in the narrative

## Notes
