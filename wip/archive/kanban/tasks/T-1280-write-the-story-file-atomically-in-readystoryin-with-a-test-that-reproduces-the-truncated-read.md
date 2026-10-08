---
id: T-1280
type: task
nature: remediation
title: Write the story file atomically in readyStoryIn, with a test that reproduces the truncated read
status: done
parent: S-0310
owner: alex
created: 2026-10-07T23:30:25Z
updated: 2026-10-07T23:39:42Z
transitions:
  - to: ready
    at: 2026-10-07T23:37:01Z
    by: agent-S-0310
  - to: in-progress
    at: 2026-10-07T23:37:01Z
    by: agent-S-0310
  - to: done
    at: 2026-10-07T23:39:42Z
    by: agent-S-0310
stream: S-0310
tags: [flai, mcp, tests]
touches: [flai/internal/mcpserver/folder_test.go]
after: [T-1279]
usage:
  source: log
  seconds: 161
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 4
      output: 1185
      cache_read: 298275
      cache_write: 9643
      cost: 0.1605
---
# T-1280 Write the story file atomically in readyStoryIn, with a test that reproduces the truncated read

## Work

Remove the cause T-1279 confirmed. If it confirmed the planner's hypothesis:

1. In `readyStoryIn` in `flai/internal/mcpserver/folder_test.go`, write the story with the criterion through `atomicfile.WriteFile` (`flai/internal/atomicfile`), as `workitem.Repo` writes items, instead of `os.WriteFile`. A reader then sees the old file or the new one, never a truncated one.
2. Add the reproduction where it fits: a test that writes a story file the old way, truncated and then filled, while a `wait_for_work` call is held, and shows the call failing on it, so the cause stays documented. If that cannot be made deterministic, record why in the narrative instead.
3. Run `TestWaitForWorkAcrossAFolder` 300 times under load, as T-1279 did, and record the result in the narrative's `## Decisions`.

If T-1279 found another cause, rewrite this task's Work and Done when for that cause before working it, and say why in `## Decisions`.

A held wait failing on a file being written by hand is a server behaviour (`folderWaits.hold` and the per-project hold return the first read error). Changing it is not in this story's scope; the plan's thread proposes it as a story of its own.

This task waits for T-1279: the fix depends on the cause it confirms, and both change `folder_test.go`.

## Done when

- `readyStoryIn` writes no story file non-atomically.
- `TestWaitForWorkAcrossAFolder` passes 300 runs under load, recorded in the narrative.
- `flai test flai/internal/mcpserver/folder_test.go` passes.

## Notes
