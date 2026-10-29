---
id: T-2282
type: task
nature: improvement
title: "flai serve keeps a budget state per provider in budget.json in the serve folder: set from a usage limit, held, probed once, cleared, or extended, each change journaled"
status: done
parent: S-0534
owner: alex
created: 2026-10-28T21:25:24Z
updated: 2026-10-29T05:56:22Z
transitions:
  - to: ready
    at: 2026-10-29T05:48:34Z
    by: agent-S-0534
  - to: in-progress
    at: 2026-10-29T05:48:35Z
    by: agent-S-0534
  - to: done
    at: 2026-10-29T05:56:22Z
    by: agent-S-0534
stream: S-0534
tags: [flai, serve, budget]
touches: [flai/internal/serve/budget.go, flai/internal/serve/budget_test.go, flai/internal/serve/serve.go]
usage:
  source: log
  seconds: 467
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 50
      output: 274
      cache_read: 2100384
      cache_write: 125633
      cost: 0.7611
---
# T-2282 flai serve keeps a budget state per provider in budget.json in the serve folder: set from a usage limit, held, probed once, cleared, or extended, each change journaled

## Work

Add the budget state of ADR-0209 §5 as a file of its own in the serve folder, `budget.json`, beside `upgrade-asks.json` and `reminders.json`. Give `Dir` its path in `serve.go` and write it as `Dir.write` writes the others: a temp file, then a rename. Take one mutex across the projects one `flai serve` serves, as `remind.go` takes one for `reminders.json`.

Key the state by provider. A provider is the name `startOf`/`provide` records on the run (`agents.go`). A run with no provider is keyed by its harness, such as `claude-code`, since it then reaches the harness's own account. Per provider, keep:

- the limit named, such as `seven_day` or `five_hour`
- `since`
- `until`: the reset time the API reported, or else now plus the backoff
- `backoff`: `5m` at first, doubling on each refusal up to `1h`; unused while a reset time is known
- `probe`: the run that is probing, if any

Give it these operations, all with an injected clock:

- `Set(provider, limit, reset)`, from a run that ended on the limit
- `Held(provider, now)`: whether a start must wait, and until when
- `Probe(provider, run, now)`: lets exactly one start through once `until` has passed, and holds every other start while it runs
- `Ended(provider, run, onLimit)`: a probe that ends on anything but the limit clears the state; one that ends on it extends the state with the new reset time or the doubled backoff
- `Read`, for the MCP server and the host API, which only read

Each change records a journal entry through the launcher's `record` hook: set, extended, probe started, and cleared. Each entry carries the provider, the limit, and `until`. A reminder through the ADR-0205 channel is the notify task's.

This task waits for nothing. Every other task of the story uses this state.

## Done when

- [ ] A `Set` with a reset time holds the provider until that time. One without a reset holds it for `5m`, and each refused probe doubles the hold up to `1h`.
- [ ] `Held` answers false for another provider, and for a provider that was cleared.
- [ ] Once `until` has passed, `Probe` lets one run through and holds the rest until that run ends.
- [ ] A probe that ends clean clears the state, and one that ends on the limit extends it.
- [ ] The state survives a restart of `flai serve`, read back from `budget.json`.
- [ ] Each change appends one journal entry naming the provider and `until`.
- [ ] `flai test` passes on the changed paths.

## Notes

Planned by planner-S-0534. The plan's thread records two assumptions: the harness name keys a run with no provider, and the next waiting start is the probe.
