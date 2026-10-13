---
id: T-1407
type: task
nature: feature
title: LiteLLM's spend reader reads a run's cost from /spend/logs by key and window, or by session where the proxy keeps it
status: done
parent: S-0357
owner: alex
created: 2026-10-08T08:50:44Z
updated: 2026-10-13T19:37:06Z
transitions:
  - to: ready
    at: 2026-10-13T19:30:58Z
    by: agent-S-0357
  - to: in-progress
    at: 2026-10-13T19:30:59Z
    by: agent-S-0357
  - to: done
    at: 2026-10-13T19:37:06Z
    by: agent-S-0357
stream: S-0357
tags: [cli]
touches: [flai/internal/usage/litellm.go, flai/internal/usage/litellm_test.go, flai/internal/usage/spend.go, flai/internal/usage/spend_test.go, flai/internal/usage/testdata/litellm]
after: [T-1406]
usage:
  source: log
  seconds: 367
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 59
      output: 30958
      cache_read: 3800980
      cache_write: 115592
      cost: 2.0745
---
# T-1407 LiteLLM's spend reader reads a run's cost from /spend/logs by key and window, or by session where the proxy keeps it

## Work

- `litellm.go`: `NewLiteLLM(baseURL, key, client)`, a `Spend` reader that calls the proxy's `GET /spend/logs` with the key as a bearer (S-0356: only the master key reads it; the provider's `spend_key_env`, T-1408, names that key's variable where the agent's own virtual key cannot), with `start_date` and `end_date` around the run's window.
  - Rows whose `session_id` is one of the run's `Sessions` (or whose `end_user` is, as S-0356 found on the pass-through) are the run's; `By` is `session`.
  - With no session match, or no sessions, the rows inside the window (by `startTime`) whose key (`api_key`, the SHA-256 of the virtual key the agent used, where the reader is told it) matches are the run's; `By` is `window`. Without the agent's key hash, a window join is refused as `ErrNoSpend`: another story's rows may fall inside it.
  - `Cost` sums `spend`; `Calls` counts rows. No row is `ErrNoSpend`.
  - The constructor takes the agent's key hash as an option, or a field of `RunRef` (add `KeyHash string` to `RunRef` in `spend.go` if needed; it is a hash, never the key).
  - Remove the placeholder T-1406 left in `For`, so `For` returns LiteLLM's reader for any host but OpenRouter's and Anthropic's.
- Tests serve recorded `/spend/logs` responses with `httptest`, shaped as `gateway_test.go`'s `logLiteLLMSpend` reads them: a join by session, a join by window and key hash with another run's rows outside it, no rows, a 401 for a key that cannot read the log, and a timeout.

It waits for T-1406, whose `Spend` and `RunRef` it implements.

## Done when

- `flai test flai/internal/usage/` passes.

## Notes

Layer 2 of S-0357. Read `design/system/agent-adapters.md` § Checked against LiteLLM before writing the join. Meets the rest of criterion 1.
