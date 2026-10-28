---
id: S-0514
type: story
nature: remediation
title: A curator run that ends on the harness's usage limit is recorded failed, logs no activity, and is not started again until the next acceptance or schedule
status: backlog
owner: alex
created: 2026-10-28T12:40:16Z
updated: 2026-10-28T22:54:04Z
transitions: []
tags: [flai, serve, curator, usage]
topics: [planning, digest]
touches: [flai/internal/usage/event.go, flai/internal/usage/claudecode.go, flai/internal/usage/claudecode_test.go, flai/internal/usage/log.go, flai/internal/usage/log_test.go, flai/internal/serve/agents.go, flai/internal/serve/curate.go, flai/internal/serve/curate_test.go, flai/internal/serve/curate_schedule.go, flai/internal/serve/curate_schedule_test.go, flai/internal/serve/limit.go, flai/internal/serve/limit_test.go, design/adrs/README.md, design/system/strategic-agents.md, docs/users/flai.md, docs/users/flai-reference.md, docs/operators/index.md, design/issues/I-0212-a-curator-run-that-ends-on-the-harness-s-usage-limit-is-recorded-failed-logs-no-activity-and-is-not-started-again-until-the-next-acceptance-or-schedule.md, design/issues/summary.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: sum
  seconds: 0
  models: []
  strategic:
    - kind: orchestrator
      seconds: 242
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 91
          output: 1404
          cache_read: 17778943
          cache_write: 31786
          cost: 3.9104
cost_of_delay:
  inputs:
    time_lost_per_cycle: 30m
    by: flai
    at: 2026-10-28T12:40:16Z
  value: 75
  by: planner-S-0514
  at: 2026-10-28T12:51:54Z
forecast:
  duration: 2h30m
  delivery: 2026-10-29T13:18:00Z
  basis: "Its own forecast of 2h30m; 23rd in the pull order with an in-progress limit of 5, behind S-0535, S-0532, S-0537, S-0534, S-0536, S-0525, S-0507, S-0241, S-0406, S-0407, S-0408, S-0409, S-0410, S-0411, S-0412, S-0424, S-0425, S-0506, S-0508, S-0509, S-0510, S-0511 and S-0512."
  by: flai
  at: 2026-10-28T22:54:04Z
finalized:
  by: alex
  at: 2026-10-28T12:56:56Z
---
# S-0514 A curator run that ends on the harness's usage limit is recorded failed, logs no activity, and is not started again until the next acceptance or schedule

## Goal

This story remediates [I-0212](../../../design/issues/I-0212-a-curator-run-that-ends-on-the-harness-s-usage-limit-is-recorded-failed-logs-no-activity-and-is-not-started-again-until-the-next-acceptance-or-schedule.md), "A curator run that ends on the harness's usage limit is recorded failed, logs no activity, and is not started again until the next acceptance or schedule". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0212 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0212 is closed with `flai issue close I-0212 --reason` saying what fixed it

## Tasks
- T-2198 The usage reader reads a run that ended on the harness's usage limit, with its reset time, from Claude Code's stream-json
- T-2199 An ADR records that a curator run that ends on the harness's usage limit is recorded as such, started again once the limit resets, and raised on a thread when that fails
- T-2200 curateEnded records a run that ended on the harness's usage limit as failed on the limit, with the reset time, from its log
- T-2201 The curate scheduler holds its triggers while the last run's usage limit has not reset, and starts the curator again once it has
- T-2202 flai serve raises on a thread a curator retry that ends on the usage limit again, or a limit whose reset time it cannot read
- T-2203 The design and the docs say how a curator run on the usage limit is recorded, retried, and raised, and I-0212 is closed

## Notes

Cost of delay inputs set by flai from I-0212. time_lost_per_cycle 30m: 30m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-28T12:03:27Z, 0 days before this story; under one cycle counts as one).

### Planning

Planned by planner-S-0514 on 2026-10-28. The plan and its assumptions are on the story's thread.

The solution proposed from I-0212's instance: the curator's log holds the limit as data. Claude Code's stream-json ends such a run with a `rate_limit_event` record (`rate_limit_info.status: rejected`, `rateLimitType: seven_day`, `resetsAt` as Unix seconds) and a `result` record with `is_error: true`, `api_error: usage_limit_reached`, `api_error_status: 429`, and the text "You've hit your weekly limit · resets Oct 28, 2am (America/Denver)". So the usage reader reads the limit and its reset time, `curateEnded` records a run that ended on it with that reason and the reset time, the curate scheduler holds the acceptances and the schedule that come meanwhile and starts one run once the limit resets, and a retry that ends on the limit again, or a limit whose reset cannot be read, is raised to the operator on a thread. A run on the limit spent nothing, so the activity document rightly gains no entry: the retry takes in from the last activity, as `since` already does.

Touches, and where each came from:

- `flai/internal/usage/event.go`, `claudecode.go`, `claudecode_test.go`, `log.go`, `log_test.go`: layout; the usage reader and the run's log are where the End's error and the reset time are read (`ccLine`, `end`, `Run.Result`). `claudecode.go` was a starting path of `flai touches suggest`.
- `flai/internal/serve/agents.go` (co-change, 27%): `AgentRun` gains the limit's fields.
- `flai/internal/serve/curate.go`, `curate_test.go`, `curate_schedule.go`, `curate_schedule_test.go`: starting paths of `flai touches suggest` and layout; `curateEnded` and the scheduler's `look` are what changes.
- `flai/internal/serve/limit.go`, `limit_test.go`: layout, new files for reading a run's limit from its log, apart from `activity.go`, which the suggestion listed by co-change but which this story does not change.
- `design/adrs/README.md`: design; the ADR task adds a row. The ADR file itself is not named: `flai adr new` allocates its number, `design/adrs` is in `claims.shared`, so no folder touch is kept, and `flai task done` widens the touches with the file once it is committed (ADR-0144).
- `design/system/strategic-agents.md` (co-change, 13%): the curator's run outcome table and the scheduler's section describe the new behaviour.
- `docs/users/flai.md`, `docs/users/flai-reference.md`, `docs/operators/index.md`: design, the three sections that say a refused or failed start is not tried again.
- `design/issues/I-0212-...md`, `design/issues/summary.md`: declared by the criteria; `flai issue close` writes both.
- Left out of the suggestion: `design/system/flai-cli.md`, `plan.go`, `serve.go`, `serve_actions.go`, `hostapi/writes.go` and their tests, `agent-narrative.md`, `settings.md`, `cmd/activity.go`, `analyze.go`, `usage.go`: co-change of past curator and activity stories, but nothing in them changes here; `flai curate`'s flags, the host API, and `digest.schedule` stay as they are.

No folder touch is kept.

Forecast: `flai forecast` gave 34m (median 95 s per unit of size over 38 done remediation stories on claude-opus-5-5, size 21: 2 criteria, 19 touches). Raised to 2h30m: the fix spans an ADR, a usage-reader extension with its own tests, a scheduler change with timed tests, a thread to the operator, and docs in three files, more than the median remediation. Delivery is flai's 2026-10-29T08:46Z moved by the added duration, 28th in the pull order with an in-progress limit of 5.

Cost of delay: `flai cod` gave 75.00 USD a week from the inputs flai set from I-0212 (30m of time lost per 168h cycle at 150 USD an hour, 1.00 cycles a week). The value stands as computed; the inputs are not changed.
