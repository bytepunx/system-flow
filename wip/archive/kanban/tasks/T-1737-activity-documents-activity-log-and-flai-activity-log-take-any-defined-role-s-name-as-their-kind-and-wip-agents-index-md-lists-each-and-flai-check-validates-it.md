---
id: T-1737
type: task
nature: feature
title: Activity documents, activity_log, and flai activity log take any defined role's name as their kind, and wip/agents/index.md lists each and flai check validates it
status: done
parent: S-0417
owner: alex
created: 2026-10-16T05:44:51Z
updated: 2026-10-29T03:49:28Z
transitions:
  - to: ready
    at: 2026-10-29T03:07:03Z
    by: agent-S-0417
  - to: in-progress
    at: 2026-10-29T03:07:04Z
    by: agent-S-0417
  - to: done
    at: 2026-10-29T03:49:28Z
    by: agent-S-0417
stream: S-0417
tags: [cli]
touches: [flai/internal/workitem/activity.go, flai/internal/workitem/activity_test.go, flai/internal/workitem/narrative.go, flai/internal/check/activity.go, flai/internal/check/activity_test.go, flai/internal/mcpserver/activity.go, flai/internal/mcpserver/activity_test.go, flai/cmd/activity.go, flai/cmd/activity_test.go, flai/internal/serve/activity.go, flai/internal/serve/agents.go, flai/internal/hostapi/people.go]
usage:
  source: log
  seconds: 2544
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 83
      output: 28770
      cache_read: 6561127
      cache_write: 256328
      cost: 3.2648
---
# T-1737 Activity documents, activity_log, and flai activity log take any defined role's name as their kind, and wip/agents/index.md lists each and flai check validates it

## Work

`ActivityKinds` in `flai/internal/workitem/activity.go` is a fixed list of three. Make an activity kind any role a project defines, read through `roles.Load`, keeping `planner`, `orchestrator`, and `analyzer`.

- `IsActivityDocument`, `ActivityPath`, `AppendActivity`, and `WriteIndex` in `narrative.go` take a defined role's document.
- The check of activity front matter covers it.
- `activity_log` and `flai activity log` accept the role's name and refuse any other kind, naming the defined ones.
- Add tests with a fixture role.

First layer: the serve path (T-1738) logs through it.

## Done when

`go test ./internal/workitem/... ./internal/check/... ./internal/mcpserver/... ./cmd/...` passes. The existing activity tests are unchanged.

## Notes
