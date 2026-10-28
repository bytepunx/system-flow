---
id: T-2200
type: task
nature: remediation
title: curateEnded records a run that ended on the harness's usage limit as failed on the limit, with the reset time, from its log
status: backlog
parent: S-0514
owner: alex
created: 2026-10-28T12:53:20Z
updated: 2026-10-28T12:54:20Z
transitions: []
stream: S-0514
tags: [flai, serve, curator]
touches: [flai/internal/serve/agents.go, flai/internal/serve/curate.go, flai/internal/serve/curate_test.go, flai/internal/serve/limit.go, flai/internal/serve/limit_test.go]
after: [T-2198, T-2199]
---
# T-2200 curateEnded records a run that ended on the harness's usage limit as failed on the limit, with the reset time, from its log

## Work

Waits for T-2198, whose `usage.Run` carries the limit, and T-2199, the ADR it follows.

In `flai/internal/serve/limit.go`, add a helper that reads a run's end on the usage limit from the run's log (`run.Log`, read with the usage package's log reader as `readActivity` and `activityMeasure.newest` in `activity.go` read it): the harness's error name, its status, and the reset time, or nothing when the log does not end on a limit or cannot be read, which is warned of. Keep it apart from `activity.go`, which this story does not change.

In `flai/internal/serve/agents.go`, give `AgentRun` two fields, documented as the others are: the time the limit resets, RFC 3339, set on a run that ended on the limit (empty when the harness named none), and a mark that the run is a retry of one that ended on the limit, for T-2201 to set. Reuse `LimitThread` for the thread T-2202 opens.

In `flai/internal/serve/curate.go`, have `curateEnded` read the limit before it judges the outcome: a run whose log ends on the limit is `failed` with `why` such as `ended on the harness's usage limit (usage_limit_reached, 429); it resets at 2026-10-28T08:00:00Z`, and the reset time on its record, whether its exit was seen (exit 1, as in I-0212) or not (a run `flai curate` started, settled by `settleOrphans`, which today is `worked`). The log line `curator ended` names the limit and the reset time. An exit other than 0 on a run whose log names no limit stays `ended (exit N)`. The activity is still logged through `logRunEnd`, which logs nothing for a run that spent nothing: say so in the comment, and that `since` keeps pointing at the last activity.

Tests in `curate_test.go`, with the `curateLab` stub and a stream whose last lines are the `rate_limit_event` and `result` records of I-0212 (a helper beside `curatorRuns`): a run that exits 1 on the limit is recorded `failed` with the why and the reset time, and the activity document gains no entry; a run whose process is gone and whose log ends on the limit is settled `failed` on the limit at a look, not `worked`; `TestACuratorRunThatFailsIsRecordedAsFailed` still holds for an exit without a limit. `limit_test.go` covers the helper on a log with the limit, one without, and a missing log.

## Done when

- A curator run whose log ends on the usage limit is recorded `failed` with a why naming the limit and the reset time, and the reset time on its record, with or without an exit seen, and logs no activity.
- A run that fails without a limit, or works, is recorded as before.
- `flai test` on the task's paths passes.

## Notes

Layer 2 of S-0514's plan. It shares `curate.go` with T-2201 and T-2202, which come in later layers.
