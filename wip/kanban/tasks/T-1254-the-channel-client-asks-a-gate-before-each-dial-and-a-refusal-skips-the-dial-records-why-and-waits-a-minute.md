---
id: T-1254
type: task
nature: feature
title: The channel client asks a gate before each dial, and a refusal skips the dial, records why, and waits a minute
status: done
parent: S-0238
owner: alex
created: 2026-10-07T23:05:40Z
updated: 2026-10-12T20:06:44Z
transitions:
  - to: ready
    at: 2026-10-12T20:02:16Z
    by: system-flow
  - to: in-progress
    at: 2026-10-12T20:02:17Z
    by: system-flow
  - to: done
    at: 2026-10-12T20:06:44Z
    by: agent-S-0238
stream: S-0238
tags: [cli]
touches: [flai/internal/channel/channel.go, flai/internal/channel/channel_test.go]
usage:
  source: log
  seconds: 267
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 23
      output: 147
      cache_read: 763124
      cache_write: 74139
      cost: 0.3805
---
# T-1254 The channel client asks a gate before each dial, and a refusal skips the dial, records why, and waits a minute

## Work

`channel.Client.Run` dials in `serveOnce` and reconnects on its own. `flai/internal/serve` cannot run a check before each dial without a hook in `flai/internal/channel/channel.go`.

- Add an optional gate to `channel.Client`, such as `BeforeDial func(ctx context.Context) error`. `Run` calls it before every attempt: the first dial and every reconnection.
- When the gate refuses, `Run` does not dial. It records the reason in `State`, waits `HeldBackoff` (a minute) as it does on `4409`, and asks the gate again.
- `Run` logs one `error` event per distinct refusal, not one per attempt.
- A client with no gate behaves as today.
- This task waits for no task and shares no path with T-1253.

S-0237 changes this file for `4403` and its backoff. Use its refusal state and backoff where it has landed, so the two refusals read alike.

## Done when

- `channel_test.go` covers four cases: a gate that passes dials, a gate that refuses does not dial and sets `LastError`, the refusal is logged once, and the gate is asked again on the next attempt and on a reconnection.
- `flai test flai/internal/channel` passes.

## Notes
