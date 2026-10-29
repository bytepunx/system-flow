---
id: T-2288
type: task
nature: improvement
title: The ADR-0205 channel gets one notice per budget state change, set, extended, or cleared, with the provider, the limit, and until, and none when no channel is set
status: done
parent: S-0534
owner: alex
created: 2026-10-28T21:26:53Z
updated: 2026-10-29T06:33:09Z
transitions:
  - to: ready
    at: 2026-10-29T05:56:39Z
    by: agent-S-0534
  - to: in-progress
    at: 2026-10-29T05:56:40Z
    by: agent-S-0534
  - to: done
    at: 2026-10-29T06:33:09Z
    by: agent-S-0534
stream: S-0534
tags: [flai, serve, notify, budget]
touches: [flai/internal/serve/notify.go, flai/internal/serve/notify_test.go, flai/internal/serve/budget.go, flai/internal/serve/budget_test.go, flai/internal/serve/serve.go]
after: [T-2282]
usage:
  source: log
  seconds: 2189
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 56
      output: 25383
      cache_read: 3783752
      cache_write: 127582
      cost: 1.9714
---
# T-2288 The ADR-0205 channel gets one notice per budget state change, set, extended, or cleared, with the provider, the limit, and until, and none when no channel is set

## Work

Send a notice through the channel ADR-0205 keeps in the host configuration (`notify.url` or `notify.command`, sent by `notify.go`'s `Send`) on each change of a provider's budget state: set, extended after a refused probe, and cleared. Send no notice when a probe starts.

The notice is JSON in the reminder's family. Give it a `kind` of `budget`, so a receiver can tell it from a thread reminder, and these fields:

- `provider`, `limit`, `since`, and `until`
- `state`: `set`, `extended`, or `cleared`
- `projects`: the projects the host serves
- `href`: the dashboard's root

It carries no token and no file content. Send one attempt, with the reminder's timeouts. It is a host-wide state, so it is sent once per change, not once per project. With no channel set, log it and warn once per run, as reminders do.

Send it from the state's change points in `budget.go` (T-2282), through a sender the launcher gives, so tests use a fake.

This task waits for T-2282. It runs beside the run-end, inbox, and host API tasks: none of them touches `notify.go` or `budget.go`.

## Done when

- [ ] Set, extended, and cleared each send exactly one notice with `kind: budget`, the provider, and `until`, and a probe's start sends none.
- [ ] Two projects under one serve folder produce one notice per change, not two.
- [ ] With no channel set, nothing is sent, and the warning is logged once.
- [ ] `flai test` passes on the changed paths.

## Notes

Planned by planner-S-0534.
