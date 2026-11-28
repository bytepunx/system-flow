---
id: T-3647
type: task
nature: remediation
title: Hold every start on an account whose login failed, and test that the orchestrator is not started again each minute
status: backlog
parent: S-0855
owner: alex
created: 2026-11-27T23:59:48Z
updated: 2026-11-28T00:00:13Z
transitions: []
stream: S-0855
tags: [flai]
touches: [flai/internal/serve/limit.go, flai/internal/serve/limit_test.go, flai/internal/serve/agents.go, flai/internal/serve/agents_restart.go, flai/internal/serve/budget.go, flai/internal/serve/budget_test.go, flai/internal/serve/orchestrate.go, flai/internal/serve/orchestrate_test.go, flai/internal/serve/notify.go, flai/internal/serve/notify_test.go]
after: [T-3644]
---
# T-3647 Hold every start on an account whose login failed, and test that the orchestrator is not started again each minute

## Work

Waits for T-3644, since it reads the login limit that task returns.

- Feed a login limit into the account's budget state as the ADR decides: `Dir.budgetEnded` and the `Budget` record in `flai/internal/serve/budget.go`, with the kind kept on the state.
- In `limit.go`, make `limitWhy` name the failed login in the run's why. The journal line then says `failed: the harness's login failed` and not `ended (exit 1)`.
- In `agents.go`, mark the run so that `look` in `orchestrate.go` waits for the budget state and not `orchestrateRetry`. Check the story agent's restart gates in `agents_restart.go` the same way: a story's agent ended on a failed login should spend no automatic restart. Widen the touches with the file if it changes.
- In `notify.go`, add the login kind's wording to `budgetNoticeTitles` and to the `BudgetNotice`.
- Reproduce I-0412 in `orchestrate_test.go` with the fixture from T-3644. An orchestrator run that ends on the failed login, with its reason still queued, starts no second run on the next look a minute later. One probe goes through after the backoff. A probe that ends cleanly clears the hold, and the queued reason starts the orchestrator.
- Add cases to `budget_test.go`, `limit_test.go`, and `notify_test.go` for the set, extended, and cleared login state.

## Done when

- The reproducing test fails on the code before this task and passes after it.
- `flai test` on `flai/internal/serve` passes.

## Notes
