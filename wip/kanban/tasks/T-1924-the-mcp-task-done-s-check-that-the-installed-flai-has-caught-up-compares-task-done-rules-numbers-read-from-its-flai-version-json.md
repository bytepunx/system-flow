---
id: T-1924
type: task
nature: improvement
title: The MCP task_done's check that the installed flai has caught up compares task-done rules numbers, read from its flai version --json
status: done
parent: S-0445
owner: alex
created: 2026-10-20T20:42:17Z
updated: 2026-10-20T21:05:41Z
transitions:
  - to: ready
    at: 2026-10-20T20:54:31Z
    by: agent-S-0445
  - to: in-progress
    at: 2026-10-20T20:54:31Z
    by: agent-S-0445
  - to: done
    at: 2026-10-20T21:05:41Z
    by: agent-S-0445
stream: S-0445
tags: [flai, mcp, task-done]
touches: [flai/internal/mcpserver/task.go, flai/internal/mcpserver/task_test.go, flai/internal/mcpserver/installed.go, flai/internal/mcpserver/installed_test.go]
after: [T-1923]
usage:
  source: log
  seconds: 670
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 31
      output: 8952
      cache_read: 2011062
      cache_write: 43645
      cost: 0.8795
---
# T-1924 The MCP task_done's check that the installed flai has caught up compares task-done rules numbers, read from its flai version --json

## Work

- In `flai/internal/mcpserver/installed.go`, read `task_done_rules` from the installed flai's `version --json` beside `version`.
- In `closeWithInstalled` in `flai/internal/mcpserver/task.go`, decide "caught up" as T-1923's refusal does: by rules numbers when the newest tag and the installed flai both carry one, else by `BehindRelease` as today.
- Fix the gap the explorer found: an installed flai whose version is not a release (such as `dev`) counts as caught up only as ADR-0160 point 2 says; cover it with a test.
- Tests in `installed_test.go` and `task_test.go`: an installed flai with the newest rules number closes the task; one with a lower number refuses; one that prints no number falls back to the minor rule.

It waits for T-1923 because it reuses the comparison and the `Outdated` field T-1923 adds.

## Done when

- An MCP server whose own flai is refused closes the task through the installed flai whenever the installed flai's rules number is not below the newest tag's.
- The tests pass with `flai test` on the four paths.

## Notes

In practice the running server and the installed flai now rarely differ in rules number, so this path runs only on a release that raised it.
