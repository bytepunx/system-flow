---
id: T-2198
type: task
nature: remediation
title: The usage reader reads a run that ended on the harness's usage limit, with its reset time, from Claude Code's stream-json
status: done
parent: S-0514
owner: alex
created: 2026-10-28T12:52:46Z
updated: 2026-10-29T01:02:03Z
transitions:
  - to: ready
    at: 2026-10-29T00:58:05Z
    by: agent-S-0514
  - to: in-progress
    at: 2026-10-29T00:58:06Z
    by: agent-S-0514
  - to: done
    at: 2026-10-29T01:02:03Z
    by: agent-S-0514
stream: S-0514
tags: [flai, usage]
touches: [flai/internal/usage/event.go, flai/internal/usage/claudecode.go, flai/internal/usage/claudecode_test.go, flai/internal/usage/log.go, flai/internal/usage/log_test.go, design/adrs/0214-a-curator-run-that-ends-on-the-harness-s-usage-limit-is-recorded-failed-with.md]
usage:
  source: log
  seconds: 235
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 46
      output: 214
      cache_read: 1663693
      cache_write: 94329
      cost: 0.5984
---
# T-2198 The usage reader reads a run that ended on the harness's usage limit, with its reset time, from Claude Code's stream-json

## Work

Give `usage.Event` a neutral mark of a run that ended on the harness's usage limit, so that every harness adapter can carry it (ADR-0130): on an `End`, the error the harness reports, such as `usage_limit_reached` with its HTTP status, and the time the limit resets, a `time.Time`, zero when the harness names none.

In `flai/internal/usage/claudecode.go`, read it from Claude Code's stream-json as the incident's log holds it (I-0212, `serve/agents/sf-curator-20261026T160352Z.log`):

- a `rate_limit_event` record, `{"type":"rate_limit_event","rate_limit_info":{"status":"rejected","resetsAt":1793174400,"rateLimitType":"seven_day",...}}`, where `resetsAt` is Unix seconds; one with `status: allowed` marks nothing;
- the `result` record, with `is_error: true`, `api_error: "usage_limit_reached"`, `api_error_status: 429`, and `result` holding the text "You've hit your weekly limit · resets Oct 28, 2am (America/Denver)". Add these fields to `ccLine`, top level, apart from the nested `is_error` of a tool result, and add `rate_limit_event` to `claudeCodeRecords` so the line is read at all.

In `flai/internal/usage/log.go`, give `Run` the limit as the run's log left it, set from the newest `End` and the reset time the `rate_limit_event` before it gave, so that `serve` can read it from the log without reading the stream again.

Tests: `claudecode_test.go` reads the two records as the incident wrote them and gives the End's error, status, and reset time, `2026-10-28T08:00:00Z`; a `result` with `is_error: false` gives none; `log_test.go` reads a log that ends on the limit into a `Run` that carries it, and one that does not into a `Run` that carries none. Run `flai test` on the paths.

## Done when

- `usage.Event` and `usage.Run` carry a run's end on the usage limit, with the harness's error name, its status, and the reset time, and carry none for a run that ended otherwise.
- The Claude Code reader reads the `rate_limit_event` and the `result` records of I-0212's log into them, with tests that use the records as the incident wrote them.
- `flai test` on the task's paths passes, gofmt and vet included.

## Notes

Layer 1 of S-0514's plan; waits for nothing. `resetsAt` 1793174400 is 2026-10-28T08:00:00Z, which is "Oct 28, 2am (America/Denver)": the Unix seconds need no time zone.
