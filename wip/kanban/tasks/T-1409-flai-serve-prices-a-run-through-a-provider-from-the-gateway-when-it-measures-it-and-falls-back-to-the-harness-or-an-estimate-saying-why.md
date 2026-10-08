---
id: T-1409
type: task
nature: feature
title: flai serve prices a run through a provider from the gateway when it measures it, and falls back to the harness or an estimate, saying why
status: backlog
parent: S-0357
owner: alex
created: 2026-10-08T08:51:03Z
updated: 2026-10-08T08:51:03Z
transitions: []
stream: S-0357
tags: [cli]
touches: [flai/internal/serve/usage.go, flai/internal/serve/usage_test.go]
after: [T-1405, T-1407, T-1408]
---
# T-1409 flai serve prices a run through a provider from the gateway when it measures it, and falls back to the harness or an estimate, saying why

## Work

- When `serve/usage.go` measures a story's runs, for each run started with a provider, resolve the provider as at start, pick the reader with `usage.For`, and ask it for the run's cost with a short timeout; use the spend key's variable when the entry names one.
- The gateway's figures replace the run's cost per model, tokens stay the harness's, and `priced_by` is `gateway`. With no provider, `harness` when the log gave a `result`, else `estimate`. A reader that fails or returns nothing leaves the harness's or the estimate, logs the reason once per run in `flai serve`'s log, and the next measurement asks again.
- The run's provider name is read from its start record, so a run started before the provider changed is priced by the gateway it ran through.
- Tests with a fake reader: a gateway-priced run, a failed reader falling back, a run with no provider, and a story summing one of each to `estimate`.

It waits for T-1405 (`priced_by`), T-1407 (both readers, after T-1406), and T-1408 (the spend key).

## Done when

- `flai test flai/internal/serve/` passes.

## Notes

Layer 3 of S-0357.
