---
id: T-1407
type: task
nature: feature
title: LiteLLM's spend reader reads a run's cost from /spend/logs by key and window, or by session where the proxy keeps it
status: backlog
parent: S-0357
owner: alex
created: 2026-10-08T08:50:44Z
updated: 2026-10-08T08:50:44Z
transitions: []
stream: S-0357
tags: [cli]
touches: [flai/internal/usage/litellm.go, flai/internal/usage/litellm_test.go]
after: [T-1406]
---
# T-1407 LiteLLM's spend reader reads a run's cost from /spend/logs by key and window, or by session where the proxy keeps it

## Work

- `litellm.go`: a `Spend` reader that calls `/spend/logs` with the key from its variable (or the spend key's, T-1408), filtered to the run's window, and joins rows by the session ID where S-0356 found the proxy keeps it, by window alone otherwise, with the reason in the result. It sums `spend` and tokens per model.
- Tests serve recorded `/spend/logs` responses with `httptest`: a join by session, a join by window with another run's rows outside it, a 403 for a key that cannot read the log, and a timeout.

It waits for T-1406, whose `Spend` and `RunRef` it implements.

## Done when

- `flai test flai/internal/usage/` passes.

## Notes

Layer 2 of S-0357. Read `agent-adapters.md` § Checked against LiteLLM before writing the join.
