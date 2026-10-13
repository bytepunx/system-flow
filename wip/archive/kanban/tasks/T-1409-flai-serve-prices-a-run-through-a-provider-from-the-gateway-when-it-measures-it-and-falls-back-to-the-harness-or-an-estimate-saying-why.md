---
id: T-1409
type: task
nature: feature
title: flai serve prices a run through a provider from the gateway when it measures it, and falls back to the harness or an estimate, saying why
status: done
parent: S-0357
owner: alex
created: 2026-10-08T08:51:03Z
updated: 2026-10-13T20:06:42Z
transitions:
  - to: ready
    at: 2026-10-13T19:37:28Z
    by: agent-S-0357
  - to: in-progress
    at: 2026-10-13T19:37:28Z
    by: agent-S-0357
  - to: done
    at: 2026-10-13T20:06:42Z
    by: agent-S-0357
stream: S-0357
tags: [cli]
touches: [flai/internal/serve/usage.go, flai/internal/serve/usage_test.go, flai/internal/serve/agents.go, flai/internal/usage/log.go, flai/internal/usage/log_test.go, flai/cmd/serve_usage.go]
after: [T-1405, T-1407, T-1408]
usage:
  source: log
  seconds: 1754
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 173
      output: 91213
      cache_read: 11198886
      cache_write: 340571
      cost: 6.1121
---
# T-1409 flai serve prices a run through a provider from the gateway when it measures it, and falls back to the harness or an estimate, saying why

## Work

- A spend record per run, beside its log in flai serve's agents folder (`<log>.spend.json`): the provider's name, and the key's usage read at start where the reader is a `usage.KeyReader`; then, once a measurement prices the run, its `cost`, `by`, and `priced_by`. It holds variable names at most, never a key.
- `agents.go`: when a story's agent starts with a provider (`provide`), write the record; where the reader is a `KeyReader`, read the key's usage first, with a short timeout, and a failed read is logged (`warn`, no key) and leaves the record without it.
- `usage/log.go`: each `Run` gains the harness session IDs its log holds (`Sessions`), so a run can be looked up by session.
- `serve/usage.go`, when it measures a story: read each log apart, so each run is priced on its own:
  - a run whose record names a provider: resolve the provider as at start (host providers and the manifest's overrides), pick the reader with `usage.For`, its key from the variable `spend_key_env` names, else `key_env`, and ask it for the run, with a timeout of about 10 seconds. The run's priced record is used without asking again.
  - The gateway's cost replaces the run's cost; it is apportioned over the run's models by their share of the harness's cost, or of their tokens when that is nothing; tokens stay the harness's; `priced_by` is `gateway`; the harness's own cost is kept out of the total.
  - A reader that fails, or answers `ErrNoSpend`, leaves the run priced by the harness, or estimated, logs `warn` with the story, the provider, and the reason once per measurement, and the next measurement asks again.
  - A run with no provider: `harness` when its log's ends reported every call's cost, `estimate` when any was estimated.
  - The story's usage is the sum of its runs, `priced_by` the least certain (T-1405). Each task's share is scaled per model by the story's ratio of priced to harness cost, and carries the story's `priced_by`.
- Tests with a fake reader in place of `usage.For` (a package variable or a field the tests set): a gateway-priced run, a failed reader falling back and asking again at the next measurement, a run with no provider, a story summing one of each to `estimate`, and the start record written without a key.

It waits for T-1405 (`priced_by`), T-1407 (both readers, after T-1406), and T-1408 (the spend key).

## Done when

- `flai test flai/internal/serve/ flai/internal/usage/` passes.

## Notes

Layer 3 of S-0357. Meets criteria 2 and 3 and the rest of criterion 4.
