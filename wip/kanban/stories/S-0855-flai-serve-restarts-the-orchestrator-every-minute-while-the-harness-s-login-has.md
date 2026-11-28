---
id: S-0855
type: story
nature: remediation
title: flai serve restarts the orchestrator every minute while the harness's login has expired, 666 runs in 14 hours, with no hold and no thread to the operator
status: backlog
owner: alex
created: 2026-11-27T22:55:27Z
updated: 2026-11-28T00:03:04Z
transitions: []
tags: [flai, flaiover]
topics: [serve, budget]
touches: [design/adrs, design/adrs/README.md, flai/internal/usage/event.go, flai/internal/usage/claudecode.go, flai/internal/usage/claudecode_test.go, flai/internal/serve/testdata/auth-expired-result.jsonl, flai/internal/serve/limit.go, flai/internal/serve/limit_test.go, flai/internal/serve/agents.go, flai/internal/serve/budget.go, flai/internal/serve/budget_test.go, flai/internal/serve/orchestrate.go, flai/internal/serve/orchestrate_test.go, flai/internal/serve/notify.go, flai/internal/serve/notify_test.go, flai/internal/serve/budget_thread.go, flai/internal/serve/budget_thread_test.go, flaiover/src/lib/hostflai.svelte.ts, flaiover/src/lib/components/BudgetBanner.svelte, flaiover/src/lib/components/BudgetBanner.svelte.test.ts, design/system/flai-cli.md, design/system/strategic-agents.md, docs/operators/index.md, docs/users/flai.md, docs/users/flaiover.md, design/issues/I-0412-flai-serve-restarts-the-orchestrator-every-minute-while-the-harness-s-login-has.md, design/issues/summary.md, flai/internal/serve/agents_restart.go]
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
    - kind: planner
      seconds: 45
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 37
          output: 8138
          cache_read: 1403770
          cache_write: 79843
          cost: 0.7748
    - kind: orchestrator
      seconds: 3
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 0
          output: 1
          cache_read: 17772
          cache_write: 340
          cost: 0.0058
    - kind: analyzer
      seconds: 85
      estimated: true
      models:
        - model: claude-fable-5-1
          input: 55
          output: 5532
          cache_read: 391176
          cache_write: 29514
          cost: 0.9652
        - model: claude-opus-5-5
          input: 29
          output: 4664
          cache_read: 1264948
          cache_write: 28114
          cost: 0.487
cost_of_delay:
  inputs:
    time_lost_per_cycle: 14h
    by: flai
    at: 2026-11-27T22:55:27Z
  value: 2100
  by: planner-S-0855
  at: 2026-11-27T23:59:15Z
forecast:
  duration: 4h
  delivery: 2026-11-28T12:16:00Z
  basis: "Its own forecast of 4h; 15th in the pull order with an in-progress limit of 3, behind S-0837, S-0844, S-0846, S-0835, S-0814, S-0852, S-0813, S-0822, S-0848, S-0853, S-0854, S-0824, S-0845, S-0241, S-0815, S-0849 and S-0851."
  by: flai
  at: 2026-11-28T00:03:04Z
finalized:
  by: alex
  at: 2026-11-28T00:01:55Z
---
# S-0855 flai serve restarts the orchestrator every minute while the harness's login has expired, 666 runs in 14 hours, with no hold and no thread to the operator

## Goal

This story remediates [I-0412](../../../design/issues/I-0412-flai-serve-restarts-the-orchestrator-every-minute-while-the-harness-s-login-has.md), "flai serve restarts the orchestrator every minute while the harness's login has expired, 666 runs in 14 hours, with no hold and no thread to the operator", which the analysis in [2026-11-27-all.md](../../../design/analysis/2026-11-27-all.md) found. The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0412 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0412 is closed with `flai issue close I-0412 --reason` saying what fixed it

## Tasks
- T-3643 Decide in an ADR how flai serve holds starts and tells the operator when a run ends on the harness's failed login
- T-3644 Read a claude-code result that failed to authenticate as a login limit
- T-3647 Hold every start on an account whose login failed, and test that the orchestrator is not started again each minute
- T-3649 Open one thread to the operator in each project a failed login holds, and resolve it when the hold clears
- T-3650 Make the dashboard's budget banner say the login failed and how to renew it
- T-3651 Document the login hold in the design and the docs, and close I-0412

## Notes

The analyzer found I-0412: the analysis the goal links gives its evidence.

Cost of delay inputs set by flai from I-0412. time_lost_per_cycle 14h carried over from I-0412's Impact section. Its Impact time_lost_per_cycle was taken rather than the 14h derived from its cost and count.

Strategic usage carried over from I-0412: analyzer 1m25s, 1.4522 USD.

### Planning

The cause, as read at planning: `flai/internal/serve/orchestrate.go` gives a failed run only the one-minute `orchestrateRetry` gate and puts its reasons back on the queue. `flai/internal/usage/claudecode.go` (`ended`) reads a result as a limit only on status 429 or after a rejecting `rate_limit_event`. So Claude Code's "Failed to authenticate: OAuth session expired" result sets no budget state ([ADR-0209](../../../design/adrs/0209-the-orchestrator-runs-for-a-reason-and-ends-started-by-flai-serve-on-the-events.md) §5), holds nothing, and tells no one. The plan proposes that a run ending on the harness's failed login feeds the account's budget state as a hold of its own kind, with the state's backoff, probe, and notice, and one thread to the operator. The first task, T-3643, decides that in an ADR before any code.

Touches, by source:

- Declared: none. `flai touches suggest` had no path to start from.
- Design: `design/adrs/README.md` and the folder `design/adrs`. The folder is kept because the ADR's number is the next free one when it is written: S-0852 claims 0322 already, and other planners are writing. T-3643 names a tentative file; `task_done` widens the story to the real one.
- Layout, from the code read at planning: `flai/internal/usage/event.go`, `claudecode.go`, and `claudecode_test.go` (detection); `flai/internal/serve/limit.go`, `agents.go`, `agents_restart.go`, `budget.go`, `orchestrate.go`, `notify.go`, and their tests (hold, restart gates, journal wording, notice); `flai/internal/serve/testdata/auth-expired-result.jsonl` (new fixture); `flai/internal/serve/budget_thread.go` and its test (new, the operator's thread); `flaiover/src/lib/hostflai.svelte.ts` and `flaiover/src/lib/components/BudgetBanner.svelte` with its test (the banner's wording).
- Design docs that describe the hold: `design/system/flai-cli.md` (the host's budget state), `design/system/strategic-agents.md` (the orchestrator's start and run-end tables), `docs/operators/index.md`, `docs/users/flai.md`, and `docs/users/flaiover.md`.
- The issue: `design/issues/I-0412-…md` and `design/issues/summary.md`, which `flai issue close` writes.

Tags: `flai` and `flaiover`, since the story delivers to both.

No `after`: the story reads nothing another open story writes.

Forecast: `flai forecast` said 5m from size 2 (2 criteria, 0 touches) over 25 small remediation stories. That band does not fit an ADR, two Go layers, a new thread path, a dashboard change, and five docs. The duration is set to 4h, the length of comparable budget-hold stories with an ADR and six tasks. The delivery is kept at flai's figure, which follows from the pull order.

Cost of delay: `flai cod` gives 2100 USD a week from the operator's input, 14h lost per 168h cycle at 150 USD an hour. The figure stands: the loop ran for 14h twice in a month, so one cycle a week is a fair rate.
