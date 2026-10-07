---
id: T-1101
type: task
nature: improvement
title: flai stats counts each story's empty wakes from its agents' run logs, as text and --json
status: in-progress
parent: S-0272
owner: alex
created: 2026-10-06T22:53:15Z
updated: 2026-10-07T00:49:24Z
transitions:
  - to: ready
    at: 2026-10-07T00:49:24Z
    by: agent-S-0272
  - to: in-progress
    at: 2026-10-07T00:49:24Z
    by: agent-S-0272
stream: S-0272
tags: [cli, metrics]
touches: [flai/internal/usage/log.go, flai/internal/usage/log_test.go, flai/internal/usage/usage.go, flai/internal/usage/usage_test.go, flai/internal/workitem/usage.go, flai/internal/workitem/usage_test.go, flai/internal/workitem/front-matter-fields.txt, flai/internal/metrics/usage.go, flai/internal/metrics/usage_test.go, flai/internal/metrics/waiting.go, flai/internal/metrics/waiting_test.go, flai/cmd/stats.go, flai/cmd/check_stats_test.go]
after: [T-1097]
---
# T-1101 flai stats counts each story's empty wakes from its agents' run logs, as text and --json

## Work

Criterion 4, the code, as T-1097's ADR defines it.

- **`flai/internal/usage/log.go`:** while it reads a run log, pair each main-agent `mcp__flai__wait_for_events` tool use with its tool result. Count one empty wake for a result with `timed_out` true, no events and no changed paths. Sub-agents' calls, those with a `parent_tool_use_id`, do not count. Add the count to what a run's usage gives.
- **`flai/internal/metrics/waiting.go`:** sum the count per story and over the window, beside the thread and review waits.
- **`flai/cmd/stats.go`:** print the count in the waiting section and in `--json`.

Tests:

- `log_test.go`: a log with an empty wake, a wake with events, an `end` answer, and a sub-agent's call
- `waiting_test.go`: the per-story and window sums
- `check_stats_test.go`: the text and JSON output

It waits for T-1097, whose ADR defines what counts. The two share no path.

## Done when

- `flai stats` and `flai stats --json` report empty wakes per story and for the window, as the ADR defines them.
- The new tests pass with `scripts/flai-test.sh`.

## Notes

Drafted by the planner. S-0293 classifies every turn of a run, empty wakes among them, and touches `usage/log.go` and `cmd/stats.go` too. Whichever is worked second builds on the other's parsing rather than adding a second one.
