---
id: T-0944
type: task
nature: improvement
title: The story the issue step makes from an issue carries the issue's strategic usage, summed up to its epic
status: done
parent: S-0227
owner: alex
created: 2026-10-05T05:45:42Z
updated: 2026-10-06T21:42:48Z
transitions:
  - to: ready
    at: 2026-10-06T21:35:20Z
    by: agent-S-0227
  - to: in-progress
    at: 2026-10-06T21:35:20Z
    by: agent-S-0227
  - to: done
    at: 2026-10-06T21:42:48Z
    by: agent-S-0227
stream: S-0227
tags: [flai]
touches: [flai/internal/issues/stories.go, flai/internal/issues/stories_test.go, flai/cmd/issue.go, flai/cmd/issue_test.go, flai/internal/mcpserver/issues.go, flai/internal/mcpserver/issues_test.go]
after: [T-0940]
usage:
  source: log
  seconds: 448
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 45
      output: 17096
      cache_read: 2409470
      cache_write: 70669
      cost: 1.2463
---
# T-0944 The story the issue step makes from an issue carries the issue's strategic usage, summed up to its epic

## Work

Carry an issue's strategic usage to the story `flai issue story` and the MCP tool `issue_story` make from it. The dashboard's `issue.story` runs the CLI, so it follows. It waits for T-0940, which gives issues their usage. It shares no path with T-0942 and runs beside it.

- `issues.ForStory` (`flai/internal/issues/stories.go`) puts the issue's `strategic` entries on the `StoryDraft`.
- `flai/cmd/issue.go` and `flai/internal/mcpserver/issues.go` charge them, once the story is made, to the story and every item above it, with `workitem.ChargeStrategic` per kind, as ADR-0083 sums a planner's charge.
- The issue keeps its own figures. The story's Notes say what was carried, beside how each cost of delay input was set.
- An issue with no usage makes a story with none, as today.

## Done when

- Tests pin a story made from an issue with an `analyzer` entry, through the CLI and through MCP: the story and its epic each carry that entry, and the issue keeps it.
- A story made from an issue without usage has no `usage`: a test pins it.
- `scripts/flai-test.sh` passes.

## Notes
