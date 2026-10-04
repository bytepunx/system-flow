---
id: T-0799
type: task
nature: feature
title: A task the planner writes that fails flai check --strict or the markdown lint is refused and leaves nothing
status: done
parent: S-0255
owner: alex
created: 2026-10-04T04:03:35Z
updated: 2026-10-04T04:07:53Z
transitions:
  - to: ready
    at: 2026-10-04T04:04:27Z
    by: agent-S-0255
  - to: in-progress
    at: 2026-10-04T04:04:27Z
    by: agent-S-0255
  - to: done
    at: 2026-10-04T04:07:53Z
    by: agent-S-0255
stream: S-0255
tags: []
touches: [flai/internal/itemnew/itemnew_test.go]
usage:
  source: log
  seconds: 206
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 17
      output: 3258
      cache_read: 1410439
      cache_write: 22355
      cost: 0.5038
---

# T-0799 A task the planner writes that fails flai check --strict or the markdown lint is refused and leaves nothing

## Work

The planner's tasks must pass `flai check --strict` and the markdown lint before they are kept. `itemnew.Create` already runs the check before and after and refuses new findings, and `workitem.Create` refuses a lint finding, but only a story's lint refusal is tested. Add tests in `flai/internal/itemnew/itemnew_test.go` that a task with a body the lint refuses, and a task whose `after` names a task that does not exist or one of another story, is refused and leaves no file behind. Waits for nothing; shares no path with T-0797 or T-0798.

## Done when

- the tests pass and each fails when the refusal it pins is taken out
- `go test ./internal/itemnew/` passes

## Notes
