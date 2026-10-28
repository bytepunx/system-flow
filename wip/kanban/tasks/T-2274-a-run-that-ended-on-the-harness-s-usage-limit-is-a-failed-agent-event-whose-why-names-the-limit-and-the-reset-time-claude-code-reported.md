---
id: T-2274
type: task
nature: improvement
title: A run that ended on the harness's usage limit is a failed agent event whose why names the limit and the reset time Claude Code reported
status: backlog
parent: S-0532
owner: alex
created: 2026-10-28T21:09:11Z
updated: 2026-10-28T21:09:11Z
transitions: []
stream: S-0532
tags: [flai, serve, events, usage]
touches: [flai/internal/serve/agent_events.go, flai/internal/serve/agent_events_test.go, flai/internal/serve/limit.go, flai/internal/serve/limit_test.go]
after: [T-2271]
---
# T-2274 A run that ended on the harness's usage limit is a failed agent event whose why names the limit and the reset time Claude Code reported

## Work

When the serve helper writes a `failed` event, it reads the run's log for the usage-limit refusal with the reader S-0514 adds (`flai/internal/usage` reads Claude Code's `rate_limit_event` and the `result` record with `api_error: usage_limit_reached` and its reset time; `flai/internal/serve/limit.go` reads a run's limit from its log). When the run ended on it, `why` names the limit as Claude Code reports it, such as `seven_day`, and the reset time, in UTC, beside the exit: for example `usage limit seven_day reached, resets 2026-10-28T08:00:00Z (exit 1)`. A limit with no reset time says so. Any other failure keeps the run's own `why`.

The helper is the one place every run kind's event is written, so this task covers the story's agent and every role at once and changes no run path. It waits for T-2271, which writes the helper, and needs S-0514's reader on main, which the story names in `after`. It runs beside the story's agent and roles tasks.

## Done when

- [ ] A test builds a run log from a recorded Claude Code usage-limit refusal and reads back one `failed` event whose `why` names the limit and the reset time
- [ ] A refusal without a reset time gives a `why` that names the limit and says the reset is unknown
- [ ] A run that failed for another reason keeps its own `why`
- [ ] `flai test` passes on the changed paths

## Notes

Planned by planner-S-0532. If S-0514's reader is not on main when this task starts, the plan is wrong: say so in the narrative's `## Decisions` and on the plan's thread rather than build a second reader.
