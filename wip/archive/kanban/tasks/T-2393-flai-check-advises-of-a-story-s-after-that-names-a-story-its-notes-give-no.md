---
id: T-2393
type: task
nature: improvement
title: flai check advises of a story's after that names a story its Notes give no reason for
status: done
parent: S-0551
owner: alex
created: 2026-10-31T07:43:11Z
updated: 2026-10-31T10:05:28Z
transitions:
  - to: ready
    at: 2026-10-31T08:42:43Z
    by: agent-S-0551
  - to: in-progress
    at: 2026-10-31T08:42:45Z
    by: agent-S-0551
  - to: done
    at: 2026-10-31T08:57:10Z
    by: agent-S-0551
stream: S-0551
tags: [flai, go]
touches: [flai/internal/check/check.go, flai/internal/check/check_test.go, flai/internal/check/after_reason_test.go, flai/internal/check/scope.go, flai/internal/check/scope_test.go, flai/internal/itemnew/itemnew.go, flai/internal/itemedit/itemedit.go, flai/internal/itemnew/itemnew_test.go, flai/internal/storygit/trial.go, flai/internal/storygit/trial_test.go]
usage:
  source: log
  seconds: 865
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 30
      output: 9459
      cache_read: 1860123
      cache_write: 68222
      cost: 0.9415
---
# T-2393 flai check advises of a story's after that names a story its Notes give no reason for

## Work

In `flai/internal/check/check.go`, add an advisory finding, such as `wip.after-reason`, on an open story whose `after` names a story that its `## Notes` do not name. Make it advisory with `c.advise`, as `afterHub` is, so that `--strict` passes over it and it stops no close-out. Leave tasks out: a task says why it waits in its `## Work`. Read the Notes as `PlanningNote` in `flai/internal/workitem/hold.go` does, without changing that file.

It waits for nothing: the rule is the story's goal, and the ADR task records it at the same time.

## Done when

- Tests in `check_test.go`, next to `TestAfterHub`, show a story whose Notes name its `after` story raising nothing, one whose Notes do not raising the advisory, a done or cancelled story raising nothing, and a task's `after` raising nothing.
- `flai test flai/internal/check` passes.

## Notes

The finding names the story and each `after` entry without a reason.
