---
id: T-2283
type: task
nature: improvement
title: A story's agent or a role's run that ends on the usage limit sets its provider's budget state and is not counted a failed run against auto_restarts or orchestrateRetry
status: backlog
parent: S-0534
owner: alex
created: 2026-10-28T21:25:49Z
updated: 2026-10-28T21:25:49Z
transitions: []
stream: S-0534
tags: [flai, serve, budget]
touches: [flai/internal/serve/agents.go, flai/internal/serve/agents_test.go, flai/internal/serve/plan.go, flai/internal/serve/plan_test.go, flai/internal/serve/orchestrate.go, flai/internal/serve/orchestrate_test.go, flai/internal/serve/analyze.go, flai/internal/serve/analyze_test.go, flai/internal/serve/curate.go, flai/internal/serve/curate_test.go, flai/internal/serve/testdata/usage-limit-result.jsonl]
after: [T-2282]
---
# T-2283 A story's agent or a role's run that ends on the usage limit sets its provider's budget state and is not counted a failed run against auto_restarts or orchestrateRetry

## Work

Feed the budget state from every run's end:

- a story's agent: `judgeRun` and `endRun` in `agents.go`, and orphans in `settleOrphans`
- the planner: `planEnded` in `plan.go`
- the orchestrator: `orchestrateEnded` in `orchestrate.go`
- the analyzer: `analyzeEnded` in `analyze.go`
- the curator: `curateEnded` in `curate.go`

Read whether a run ended on the usage limit, and its reset time, with S-0514's reader (`flai/internal/usage`, `flai/internal/serve/limit.go`). S-0534 waits for S-0532, which waits for S-0514, so the reader is on main. Do not write a second one. On the limit, call the state's `Set` with the run's provider (T-2282). When the run was the probe, call `Ended` instead.

A run that ended on the limit is recorded with that outcome, not as an ordinary failure:

- `restartEnded` and `autoRestart` neither restart it nor count it against `agent.auto_restarts`, so the story keeps its count.
- `orchestrateRetry`'s gate does not start the orchestrator a minute later.
- The curator's own hold from S-0514 gives way to this state, with one hold rather than two.

S-0532 already makes such a run a `failed` agent event whose `why` names the limit; leave that as it is.

Reproduce the refusal from a recorded Claude Code result line, kept as `testdata/usage-limit-result.jsonl`. Take it from S-0514's tests if it kept one there, with its `rate_limit_event` and its `result` with `api_error: usage_limit_reached`. Then test one story's agent and one role through it.

This task waits for T-2282's state.

## Done when

- [ ] A story's agent and each of the four roles, ending on the recorded refusal, set the provider's state with the reported reset time, or with the backoff when the line has none.
- [ ] Such a story's agent is not restarted and its `auto_restarts` count is unchanged. The orchestrator is not started again after `orchestrateRetry`.
- [ ] A probe run that ends clean clears the state, and one that ends on the limit extends it.
- [ ] A run that fails for any other reason is judged and restarted as before.
- [ ] `flai test` passes on the changed paths.

## Notes

Planned by planner-S-0534. Where the reader stands is on the plan's thread.
