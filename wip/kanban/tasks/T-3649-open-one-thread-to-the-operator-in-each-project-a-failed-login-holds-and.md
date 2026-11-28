---
id: T-3649
type: task
nature: remediation
title: Open one thread to the operator in each project a failed login holds, and resolve it when the hold clears
status: backlog
parent: S-0855
owner: alex
created: 2026-11-27T23:59:56Z
updated: 2026-11-27T23:59:56Z
transitions: []
stream: S-0855
tags: [flai]
touches: [flai/internal/serve/budget_thread.go, flai/internal/serve/budget_thread_test.go, flai/internal/serve/budget.go]
after: [T-3647]
---
# T-3649 Open one thread to the operator in each project a failed login holds, and resolve it when the hold clears

## Work

Waits for T-3647, since it opens the thread from the login state's set and clear, which that task writes. It also changes `budget.go` after that task.

- Add `flai/internal/serve/budget_thread.go`. When a login state is set, open one thread to the operator in each served project whose agents run on the account. It says what failed, since when, what it holds, and how to renew the login. A state that is extended writes nothing new. When the state clears, resolve the thread with the reason. Follow the way `curate.go` (`raiseCurateLimit`) and `claudecheck.go` open their threads, and record the thread on the state so that it is opened only once.
- Call it from `launcher.budgetChanged` in `budget.go`, for the login kind only. A usage limit keeps its notice alone, as ADR-0209 has it.
- Test in `budget_thread_test.go`: the thread is opened once over several extensions, it is opened in each project the state holds, and it is resolved on a clear by a probe and on a clear by hand.

## Done when

- `flai test` on `flai/internal/serve` passes, with the new cases.
- No login state opens a second thread while the first is open.

## Notes
